package peers

import (
	"context"
	"errors"
	libhead "github.com/celestiaorg/celestia-node/libs/header"
	"github.com/celestiaorg/celestia-node/share/availability/discovery"
	logging "github.com/ipfs/go-log/v2"
	"github.com/libp2p/go-libp2p/core/peer"
	"sync"
	"sync/atomic"
	"time"

	pubsub "github.com/libp2p/go-libp2p-pubsub"

	"github.com/celestiaorg/celestia-node/header"
	"github.com/celestiaorg/celestia-node/share"
	"github.com/celestiaorg/celestia-node/share/p2p/shrexsub"
)

// TODO: Needs improvement:
//   - might need rework / adjustment: if node was behind head, DASer will start only after some header is sampled / sample timed out
//   - need to add cleanup, otherwise items store memory footpritn will grow over time
//   - note: recent block will be sampled without being restricted by DASer concurrency limit.
//     If DASer is busy with heavy tasks it could take a while until worker slot is available
//   - make params configurable
//   - add metrics

var (
	log         = logging.Logger("shrex/peers")
	syncTimeout = 60 * time.Second
)

// Manager keeps track of peers coming from shrex.Sub and discovery
type Manager struct {
	disc      discovery.Discovery
	headerSub libhead.Subscription[*header.ExtendedHeader]

	m         *sync.Mutex
	pools     map[hashStr]syncPool
	fullNodes pool

	cancel context.CancelFunc
	done   chan struct{}
}

type hashStr string

type syncPool struct {
	pool *pool

	isSynced  *atomic.Bool
	syncCh    chan struct{}
	syncTimer *time.Timer
}

func NewManager(
	shrexSub *shrexsub.PubSub,
	headerSub libhead.Subscription[*header.ExtendedHeader],
	discovery discovery.Discovery,
) (*Manager, error) {
	s := &Manager{
		disc:      discovery,
		headerSub: headerSub,
		m:         new(sync.Mutex),
		pools:     make(map[hashStr]syncPool),
		done:      make(chan struct{}),
	}

	discovery.WithOnPeersUpdate(
		func(peerID peer.ID, isAdded bool) {
			if isAdded {
				s.fullNodes.add(peerID)
				return
			}
			s.fullNodes.remove(peerID)
		})

	err := shrexSub.AddValidator(s.validator())
	return s, err
}

func (s *Manager) Start() {
	ctx, cancel := context.WithCancel(context.Background())
	s.cancel = cancel
	go s.disc.EnsurePeers(ctx)
	go s.subscribeHeader(ctx)
}

func (s *Manager) Stop(ctx context.Context) error {
	s.cancel()
	select {
	case <-s.done:
		return nil
	case <-ctx.Done():
		return ctx.Err()
	}
}

func (s *Manager) Get(ctx context.Context, h *header.ExtendedHeader) (peer.ID, error) {
	p := s.markSynced(hashStr(h.DataHash.String()))
	peer, ok := p.pool.tryGet()
	if ok {
		return peer, nil
	}

	// try fullnodes obtained from discovery
	peer, ok = s.fullNodes.tryGet()
	if ok {
		return peer, nil
	}

	// no peers available, wait for the first one
	select {
	case peer = <-p.pool.waitNext(ctx):
		return peer, nil
	case peer = <-s.fullNodes.waitNext(ctx):
		return peer, nil
	case <-ctx.Done():
		return "", ctx.Err()
	}
}

func (s *Manager) Remove(h *header.ExtendedHeader, ids ...peer.ID) {
	p := s.markSynced(hashStr(h.DataHash.String()))
	p.pool.remove(ids...)
}

func (s *Manager) subscribeHeader(ctx context.Context) {
	defer close(s.done)
	defer s.headerSub.Cancel()

	for {
		h, err := s.headerSub.NextHeader(ctx)
		if err != nil {
			if errors.Is(err, context.Canceled) {
				return
			}
			log.Errorw("get next header", "err", err)
			continue
		}

		s.markSynced(hashStr(h.DataHash.String()))
	}
}

func (s *Manager) markSynced(key hashStr) syncPool {
	s.m.Lock()
	defer s.m.Unlock()

	p, ok := s.pools[key]
	if !ok {
		// save as already synced
		p.pool = newPool()
		p.isSynced = new(atomic.Bool)
		p.isSynced.Store(true)
		s.pools[key] = p
		return p
	}

	if p.isSynced.CompareAndSwap(false, true) {
		// indicate that sync is done to unlock Validators.
		// if unable to stop the timer, the channel was already closed by afterfunc
		if p.syncTimer.Stop() {
			close(p.syncCh)
		}
	}
	return p
}

// Validator will block until header with given datahash received. This behavior opens an attack vector of multiple fake
// datahash spam, that will grow amount of hanging routines in node. To address this, validator should be reworked to be non-blocking,
// with retransmission being invoked upon header discovery and from another routine in sync manner.
func (s *Manager) validator() shrexsub.Validator {
	return func(ctx context.Context, peerID peer.ID, hash share.DataHash) pubsub.ValidationResult {
		p := s.addPool(hash)
		p.pool.add(peerID)

		if p.isSynced.Load() {
			// header is synced, allow the pubsub to retransmit the message by returning Accept
			return pubsub.ValidationAccept
		}

		// wait for header to sync within timeout.
		select {
		case <-p.syncCh:
			if p.isSynced.Load() {
				// headerSub found corresponding ExtendedHeader for dataHash,
				//retransmit the message by returning Accept
				return pubsub.ValidationAccept
			}

			// no corresponding header was received in time,
			//highly unlikely block with given datahash exist in chain, reject msg and punish the peer
			// TODO: maybe keep pool for fast rejects later?
			s.deletePool(hashStr(hash.String()))
			return pubsub.ValidationReject
		case <-ctx.Done():
			return pubsub.ValidationIgnore
		}
	}
}

func (s *Manager) addPool(dataHash share.DataHash) syncPool {
	s.m.Lock()
	defer s.m.Unlock()

	key := hashStr(dataHash.String())
	t, ok := s.pools[key]
	if !ok {
		syncCh := make(chan struct{})
		t = syncPool{
			pool:     newPool(),
			isSynced: new(atomic.Bool),
			syncCh:   syncCh,
			syncTimer: time.AfterFunc(syncTimeout, func() {
				close(syncCh)
			}),
		}
		s.pools[key] = t
	}
	return t
}

func (s *Manager) deletePool(key hashStr) {
	s.m.Lock()
	defer s.m.Unlock()
	delete(s.pools, key)
}
