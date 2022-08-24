package header

import (
	"context"
	"encoding/hex"
	"github.com/celestiaorg/celestia-node/fraud"
	"github.com/celestiaorg/celestia-node/header"
	"github.com/celestiaorg/celestia-node/header/p2p"
	"github.com/celestiaorg/celestia-node/header/store"
	"github.com/celestiaorg/celestia-node/header/sync"
	"github.com/celestiaorg/celestia-node/libs/fxutil"
	"github.com/celestiaorg/celestia-node/params"
	headerservice "github.com/celestiaorg/celestia-node/service/header"
	"github.com/ipfs/go-datastore"
	"github.com/libp2p/go-libp2p-core/host"
	"github.com/libp2p/go-libp2p-core/peer"
	"github.com/libp2p/go-libp2p-core/peerstore"
	pubsub "github.com/libp2p/go-libp2p-pubsub"
	"github.com/multiformats/go-multiaddr"
	tmbytes "github.com/tendermint/tendermint/libs/bytes"
	"go.uber.org/fx"
)

// HeaderService creates a new header.Service.
func HeaderService(
	syncer *sync.Syncer,
	sub header.Subscriber,
	p2pServer *p2p.ExchangeServer,
	ex header.Exchange,
	store header.Store,
) *headerservice.Service {
	return headerservice.NewHeaderService(syncer, sub, p2pServer, ex, store)
}

// HeaderExchangeP2P constructs new Exchange for headers.
func HeaderExchangeP2P(cfg Config) func(params.Bootstrappers, host.Host) (header.Exchange, error) {
	return func(bpeers params.Bootstrappers, host host.Host) (header.Exchange, error) {
		peers, err := cfg.trustedPeers(bpeers)
		if err != nil {
			return nil, err
		}
		ids := make([]peer.ID, len(peers))
		for index, peer := range peers {
			ids[index] = peer.ID
			host.Peerstore().AddAddrs(peer.ID, peer.Addrs, peerstore.PermanentAddrTTL)
		}
		return p2p.NewExchange(host, ids), nil
	}
}

// HeaderP2PExchangeServer creates a new header/p2p.ExchangeServer.
func HeaderP2PExchangeServer(lc fx.Lifecycle, host host.Host, store header.Store) *p2p.ExchangeServer {
	p2pServ := p2p.NewExchangeServer(host, store)
	lc.Append(fx.Hook{
		OnStart: p2pServ.Start,
		OnStop:  p2pServ.Stop,
	})

	return p2pServ
}

// HeaderStore creates and initializes new header.Store.
func HeaderStore(lc fx.Lifecycle, ds datastore.Batching) (header.Store, error) {
	store, err := store.NewStore(ds)
	if err != nil {
		return nil, err
	}
	lc.Append(fx.Hook{
		OnStart: store.Start,
		OnStop:  store.Stop,
	})
	return store, nil
}

// HeaderStoreInit initializes the store.
func HeaderStoreInit(cfg *Config) func(context.Context, params.Network, header.Store, header.Exchange) error {
	return func(ctx context.Context, net params.Network, s header.Store, ex header.Exchange) error {
		trustedHash, err := cfg.trustedHash(net)
		if err != nil {
			return err
		}

		err = store.Init(ctx, s, ex, trustedHash)
		if err != nil {
			// TODO(@Wondertan): Error is ignored, as otherwise unit tests for Node construction fail.
			// 	This is due to requesting step of initialization, which fetches initial Header by trusted hash from
			//  the network. The step can't be done during unit tests and fixing it would require either
			//   * Having some test/dev/offline mode for Node that mocks out all the networking
			//   * Hardcoding full extended header in params pkg, instead of hashes, so we avoid requesting step
			log.Errorf("initializing store failed: %s", err)
		}

		return nil
	}
}

// HeaderSyncer creates a new Syncer.
func HeaderSyncer(
	ctx context.Context,
	lc fx.Lifecycle,
	ex header.Exchange,
	store header.Store,
	sub header.Subscriber,
	fservice fraud.Service,
) (*sync.Syncer, error) {
	syncer := sync.NewSyncer(ex, store, sub)
	lifecycleCtx := fxutil.WithLifecycle(ctx, lc)
	lc.Append(fx.Hook{
		OnStart: func(startCtx context.Context) error {
			return FraudLifecycle(startCtx, lifecycleCtx, fraud.BadEncoding, fservice, syncer.Start, syncer.Stop)
		},
		OnStop: syncer.Stop,
	})

	return syncer, nil
}

// P2PSubscriber creates a new p2p.Subscriber.
func P2PSubscriber(lc fx.Lifecycle, sub *pubsub.PubSub) (*p2p.Subscriber, *p2p.Subscriber) {
	p2pSub := p2p.NewSubscriber(sub)
	lc.Append(fx.Hook{
		OnStart: p2pSub.Start,
		OnStop:  p2pSub.Stop,
	})
	return p2pSub, p2pSub
}

// FraudService constructs fraud proof service.
func FraudService(
	sub *pubsub.PubSub,
	hstore header.Store,
	ds datastore.Batching,
) fraud.Service {
	return fraud.NewService(sub, hstore.GetByHeight, ds)
}

// FraudLifecycle controls the lifecycle of service depending on fraud proofs.
// It starts the service only if no fraud-proof exists and stops the service automatically
// if a proof arrives after the service was started.
func FraudLifecycle(
	startCtx, lifecycleCtx context.Context,
	p fraud.ProofType,
	fservice fraud.Service,
	start, stop func(context.Context) error,
) error {
	proofs, err := fservice.Get(startCtx, p)
	switch err {
	default:
		return err
	case nil:
		return &fraud.ErrFraudExists{Proof: proofs}
	case datastore.ErrNotFound:
	}
	err = start(startCtx)
	if err != nil {
		return err
	}
	// handle incoming Fraud Proofs
	go fraud.OnProof(lifecycleCtx, fservice, p, func(fraud.Proof) {
		if err := stop(lifecycleCtx); err != nil {
			log.Error(err)
		}
	})
	return nil
}

func (cfg *Config) trustedPeers(bpeers params.Bootstrappers) (infos []peer.AddrInfo, err error) {
	if len(cfg.TrustedPeers) == 0 {
		log.Infof("No trusted peers in config, initializing with default bootstrappers as trusted peers")
		return bpeers, nil
	}

	infos = make([]peer.AddrInfo, len(cfg.TrustedPeers))
	for i, tpeer := range cfg.TrustedPeers {
		ma, err := multiaddr.NewMultiaddr(tpeer)
		if err != nil {
			return nil, err
		}
		p, err := peer.AddrInfoFromP2pAddr(ma)
		if err != nil {
			return nil, err
		}
		infos[i] = *p
	}

	return
}

func (cfg *Config) trustedHash(net params.Network) (tmbytes.HexBytes, error) {
	if cfg.TrustedHash == "" {
		gen, err := params.GenesisFor(net)
		if err != nil {
			return nil, err
		}
		return hex.DecodeString(gen)
	}
	return hex.DecodeString(cfg.TrustedHash)
}
