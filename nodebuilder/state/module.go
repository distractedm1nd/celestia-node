package state

import (
	"context"

	logging "github.com/ipfs/go-log/v2"
	"go.uber.org/fx"

	"github.com/celestiaorg/celestia-node/fraud"
	"github.com/celestiaorg/celestia-node/libs/fxutil"
	fraudServ "github.com/celestiaorg/celestia-node/nodebuilder/fraud"
	"github.com/celestiaorg/celestia-node/nodebuilder/node"
	"github.com/celestiaorg/celestia-node/state"
)

var log = logging.Logger("state-module")

// ConstructModule provides all components necessary to construct the
// state service.
func ConstructModule(tp node.Type, cfg *Config) fx.Option {
	// sanitize config values before constructing module
	cfgErr := cfg.Validate()

	baseComponents := fx.Options(
		fx.Supply(*cfg),
		fx.Error(cfgErr),
		fx.Provide(Keyring),
		fx.Provide(fx.Annotate(
			CoreAccessor,
			fx.OnStart(func(
				startCtx, ctx context.Context,
				lc fx.Lifecycle,
				fservice fraudServ.Module,
				ca *state.CoreAccessor,
			) error {
				lifecycleCtx := fxutil.WithLifecycle(ctx, lc)
				return fraudServ.Lifecycle(startCtx, lifecycleCtx, fraud.BadEncoding, fservice, ca.Start, ca.Stop)
			}),
			fx.OnStop(func(ctx context.Context, ca *state.CoreAccessor) error {
				return ca.Stop(ctx)
			}),
		)),
		// the module is needed for the handler
		fx.Provide(func(ca *state.CoreAccessor) Module {
			return ca
		}),
	)

	switch tp {
	case node.Light, node.Full, node.Bridge:
		return fx.Module(
			"state",
			baseComponents,
		)
	default:
		panic("invalid node type")
	}
}
