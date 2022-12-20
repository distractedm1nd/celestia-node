package cache

import (
	"testing"

	ds "github.com/ipfs/go-datastore"
	dssync "github.com/ipfs/go-datastore/sync"
	mdutils "github.com/ipfs/go-merkledag/test"

	"github.com/celestiaorg/celestia-node/share"
	"github.com/celestiaorg/celestia-node/share/availability/full"
	"github.com/celestiaorg/celestia-node/share/availability/light"
	availability_test "github.com/celestiaorg/celestia-node/share/availability/test"
)

// RandLightLocalServiceWithSquare is the same as light.RandServiceWithSquare, except
// the share.Availability is wrapped with cache availability.
func RandLightLocalServiceWithSquare(t *testing.T, n int) (share.Availability, *share.Root) {
	bServ := mdutils.Bserv()
	store := dssync.MutexWrap(ds.NewMapDatastore())
	avail := NewShareAvailability(
		light.TestAvailability(bServ),
		store,
	)
	return avail, availability_test.RandFillBS(t, n, bServ)
}

// RandFullLocalServiceWithSquare is the same as full.RandServiceWithSquare, except
// the share.Availability is wrapped with cache availability.
func RandFullLocalServiceWithSquare(t *testing.T, n int) (share.Availability, *share.Root) {
	bServ := mdutils.Bserv()
	store := dssync.MutexWrap(ds.NewMapDatastore())
	avail := NewShareAvailability(
		full.TestAvailability(bServ),
		store,
	)
	return avail, availability_test.RandFillBS(t, n, bServ)
}
