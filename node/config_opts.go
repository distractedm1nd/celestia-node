package node

// WithRemoteCoreIP configures Node to connect to the given remote Core IP.
func WithRemoteCoreIP(ip string) Option {
	return func(sets *settings) {
		sets.cfg.Core.IP = ip
	}
}

// WithRemoteCorePort configures Node to connect to the given remote Core port.
func WithRemoteCorePort(port string) Option {
	return func(sets *settings) {
		sets.cfg.Core.RPCPort = port
	}
}

// WithGRPCPort configures Node to connect to given gRPC port
// for state-related queries.
func WithGRPCPort(port string) Option {
	return func(sets *settings) {
		sets.cfg.Core.GRPCPort = port
	}
}

// WithRPCPort configures Node to expose the given port for RPC
// queries.
func WithRPCPort(port string) Option {
	return func(sets *settings) {
		sets.cfg.RPC.Port = port
	}
}

// WithRPCAddress configures Node to listen on the given address for RPC
// queries.
func WithRPCAddress(addr string) Option {
	return func(sets *settings) {
		sets.cfg.RPC.Address = addr
	}
}

// WithConfig sets the entire custom config.
func WithConfig(custom *Config) Option {
	return func(sets *settings) {
		sets.cfg = custom
	}
}

// WithMutualPeers sets the `MutualPeers` field in the config.
func WithMutualPeers(addrs []string) Option {
	return func(sets *settings) {
		sets.cfg.P2P.MutualPeers = addrs
	}
}
