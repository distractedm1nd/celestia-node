package services

import (
	"context"
	"go.uber.org/fx"

	"github.com/celestiaorg/celestia-node/das"
	"github.com/celestiaorg/celestia-node/fraud"
	"github.com/celestiaorg/celestia-node/header"
	"github.com/celestiaorg/celestia-node/libs/fxutil"
	nodeheader "github.com/celestiaorg/celestia-node/node/header"
	"github.com/celestiaorg/celestia-node/service/share"
	"github.com/ipfs/go-datastore"
)

// DASer constructs a new Data Availability Sampler.
func DASer(
	ctx context.Context,
	lc fx.Lifecycle,
	avail share.Availability,
	sub header.Subscriber,
	hstore header.Store,
	ds datastore.Batching,
	fservice fraud.Service,
) *das.DASer {
	das := das.NewDASer(avail, sub, hstore, ds, fservice)
	lifecycleCtx := fxutil.WithLifecycle(ctx, lc)
	lc.Append(fx.Hook{
		OnStart: func(startContext context.Context) error {
			return nodeheader.FraudLifecycle(startContext, lifecycleCtx, fraud.BadEncoding, fservice, das.Start, das.Stop)
		},
		OnStop: das.Stop,
	})

	return das
}

// Metrics enables metrics for services.
func Metrics() fx.Option {
	return fx.Options(
		fx.Invoke(header.MonitorHead),
		// add more monitoring here
	)
}
