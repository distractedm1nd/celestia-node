package full

import (
	"testing"
	"time"

	"github.com/ipfs/go-blockservice"
	mdutils "github.com/ipfs/go-merkledag/test"
	routinghelpers "github.com/libp2p/go-libp2p-routing-helpers"
	"github.com/libp2p/go-libp2p/p2p/discovery/routing"

	"github.com/celestiaorg/celestia-node/share"
	"github.com/celestiaorg/celestia-node/share/availability/discovery"
	availability_test "github.com/celestiaorg/celestia-node/share/availability/test"
	"github.com/celestiaorg/celestia-node/share/service"
)

// RandServiceWithSquare provides a service.ShareService filled with 'n' NMT
// trees of 'n' random shares, essentially storing a whole square.
func RandServiceWithSquare(t *testing.T, n int) (*ShareAvailability, *share.Root) {
	bServ := mdutils.Bserv()
	return TestAvailability(bServ), availability_test.RandFillBS(t, n, bServ)
}

// RandNode creates a Full Node filled with a random block of the given size.
func RandNode(dn *availability_test.TestDagNet, squareSize int) (*availability_test.TestNode, *share.Root) {
	nd := Node(dn)
	return nd, availability_test.RandFillBS(dn.T, squareSize, nd.BlockService)
}

// Node creates a new empty Full Node.
func Node(dn *availability_test.TestDagNet) *availability_test.TestNode {
	nd := dn.NewTestNode()
	nd.ShareService = service.NewShareService(nd.BlockService)
	nd.Availability = TestAvailability(nd.BlockService)
	return nd
}

func TestAvailability(bServ blockservice.BlockService) *ShareAvailability {
	disc := discovery.NewDiscovery(nil, routing.NewRoutingDiscovery(routinghelpers.Null{}), 0, time.Second, time.Second)
	return NewShareAvailability(bServ, disc)
}
