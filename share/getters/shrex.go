package getters

import (
	"context"
	"errors"
	"fmt"

	"github.com/celestiaorg/celestia-node/share"
	"github.com/celestiaorg/celestia-node/share/p2p"
	"github.com/celestiaorg/celestia-node/share/p2p/peers"
	"github.com/celestiaorg/celestia-node/share/p2p/shrexeds"
	"github.com/celestiaorg/celestia-node/share/p2p/shrexnd"
	"github.com/celestiaorg/celestia-node/share/p2p/shrexsub"

	"github.com/celestiaorg/nmt/namespace"
	"github.com/celestiaorg/rsmt2d"
)

var _ share.Getter = (*ShrexGetter)(nil)

// ShrexGetter is a share.Getter that uses the shrex/eds and shrex/nd protocol to retrieve shares.
type ShrexGetter struct {
	cancel context.CancelFunc

	edsClient *shrexeds.Client
	ndClient  *shrexnd.Client
	shrexSub  *shrexsub.PubSub

	peerManager *peers.Manager
}

func NewShrexGetter(
	edsClient *shrexeds.Client,
	ndClient *shrexnd.Client,
	shrexSub *shrexsub.PubSub,
	peerManager *peers.Manager,
) *ShrexGetter {
	return &ShrexGetter{
		edsClient:   edsClient,
		ndClient:    ndClient,
		shrexSub:    shrexSub,
		peerManager: peerManager,
	}
}

func (sg *ShrexGetter) Start(context.Context) error {
	ctx, cancel := context.WithCancel(context.Background())
	sg.cancel = cancel
	sg.peerManager.Start()
	go sg.listen(ctx)
	err := sg.shrexSub.AddValidator(sg.peerManager.Validate)
	return err
}

func (sg *ShrexGetter) Stop(ctx context.Context) error {
	defer sg.cancel()
	return sg.peerManager.Stop(ctx)
}

func (sg *ShrexGetter) listen(ctx context.Context) {
	sub, err := sg.shrexSub.Subscribe()
	if err != nil {
		panic(fmt.Errorf("couldn't start shrexsub subscription: %w", err))
	}
	defer sub.Cancel()

	for {
		dataHash, err := sub.Next(ctx)
		if err != nil {
			if err == context.Canceled {
				return
			}

			log.Errorw("failed to get next datahash", "err", err)
			continue
		}

		log.Debug("received datahash over shrexsub: ", dataHash.String())
	}
}

func (sg *ShrexGetter) GetShare(ctx context.Context, root *share.Root, row, col int) (share.Share, error) {
	return nil, errors.New("shrex-getter: GetShare is not supported")
}

func (sg *ShrexGetter) GetEDS(ctx context.Context, root *share.Root) (*rsmt2d.ExtendedDataSquare, error) {
	for {
		to, setStatus, err := sg.peerManager.GetPeer(ctx, root.Hash())
		if err != nil {
			return nil, err
		}

		eds, err := sg.edsClient.RequestEDS(ctx, root.Hash(), to)
		switch err {
		case nil:
		    setStatus(peers.ResultSuccess)
			err = sg.shrexSub.Broadcast(ctx, root.Hash())
			if err != nil {
				log.Errorf("couldn't rebroadcast datahash %s: %s", root.String(), err)
			}
			return eds, nil
		case p2p.ErrInvalidResponse:
			setStatus(peers.ResultPeerMisbehaved)
			if err != nil {
				return nil, err
			}
		case context.DeadlineExceeded, context.Canceled:
			setStatus(peers.ResultFail)
			return nil, ctx.Err()
		default:
			setStatus(peers.ResultFail)
		}
	}
}

func (sg *ShrexGetter) GetSharesByNamespace(
	ctx context.Context,
	root *share.Root,
	id namespace.ID,
) (share.NamespacedShares, error) {
	for {
		to, setStatus, err := sg.peerManager.GetPeer(ctx, root.Hash())
		if err != nil {
			return nil, err
		}

		eds, err := sg.ndClient.RequestND(ctx, root, id, to)
		switch err {
		case nil:
			setStatus(peers.ResultSuccess)
			err = sg.shrexSub.Broadcast(ctx, root.Hash())
			if err != nil {
				log.Errorf("couldn't rebroadcast datahash %s: %s", root.String(), err)
			}
			return eds, nil
		case p2p.ErrInvalidResponse:
			setStatus(peers.ResultPeerMisbehaved)
			if err != nil {
				return nil, err
			}
		case context.DeadlineExceeded, context.Canceled:
			setStatus(peers.ResultFail)
			return nil, ctx.Err()
		default:
			setStatus(peers.ResultFail)
		}
	}
}