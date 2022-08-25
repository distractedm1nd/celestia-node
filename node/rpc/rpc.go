package rpc

import (
	"go.uber.org/fx"

	"github.com/celestiaorg/celestia-node/das"
	"github.com/celestiaorg/celestia-node/service/header"
	"github.com/celestiaorg/celestia-node/service/rpc"
	"github.com/celestiaorg/celestia-node/service/share"
	"github.com/celestiaorg/celestia-node/service/state"
)

// Server constructs a new RPC Server from the given Config.
// TODO @renaynay @Wondertan: this component is meant to be removed on implementation
//  of https://github.com/celestiaorg/celestia-node/pull/506.
func Server(cfg rpc.Config) func(lc fx.Lifecycle) *rpc.Server {
	return func(lc fx.Lifecycle) *rpc.Server {
		serv := rpc.NewServer(cfg)
		lc.Append(fx.Hook{
			OnStart: serv.Start,
			OnStop:  serv.Stop,
		})
		return serv
	}
}

// Handler constructs a new RPC Handler from the given services.
func Handler(
	state *state.Service,
	share *share.Service,
	header *header.Service,
	serv *rpc.Server,
	daser *das.DASer,
) {
	handler := rpc.NewHandler(state, share, header, daser)
	handler.RegisterEndpoints(serv)
	handler.RegisterMiddleware(serv)
}
