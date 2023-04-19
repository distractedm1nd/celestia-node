package share

import (
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
