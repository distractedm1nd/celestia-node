package share

import (
	disc "github.com/celestiaorg/celestia-node/share/availability/discovery"
	"github.com/celestiaorg/celestia-node/share/getters"
	"github.com/celestiaorg/celestia-node/share/p2p/shrexeds"
)

func WithShrexClientMetrics(c *shrexeds.Client) error {
	return c.WithMetrics()
}

func WithShrexServerMetrics(s *shrexeds.Server) error {
	return s.WithMetrics()
}

func WithShrexGetterMetrics(sg *getters.ShrexGetter) error {
	return sg.WithMetrics()
}

func WithDiscoveryMetrics(d *disc.Discovery) error {
	return d.WithMetrics()
}
