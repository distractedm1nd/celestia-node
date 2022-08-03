package ipld

import (
	"context"
	"sync"

	"github.com/gammazero/workerpool"
	"github.com/ipfs/go-blockservice"
	"github.com/ipfs/go-cid"
	ipld "github.com/ipfs/go-ipld-format"

	"github.com/celestiaorg/celestia-node/ipld/plugin"
	"github.com/celestiaorg/nmt"
	"github.com/celestiaorg/nmt/namespace"
)

// TODO(@distractedm1nd) Should the namespace pool use NumWorkersLimit? This should probably be configurable.
// Worker pool responsible for the goroutines spawned by getLeavesByNamespace
var namespacePool = workerpool.New(NumWorkersLimit)

// GetSharesByNamespace walks the tree of a given root and returns its shares within the given namespace.ID.
func GetSharesByNamespace(
	ctx context.Context,
	bGetter blockservice.BlockGetter,
	root cid.Cid,
	nID namespace.ID,
	maxShares int,
) ([]Share, error) {
	leaves, err := getLeavesByNamespace(ctx, bGetter, root, nID, maxShares)
	if err != nil {
		return nil, err
	}

	shares := make([]Share, len(leaves))
	for i, leaf := range leaves {
		shares[i] = leafToShare(leaf)
	}

	return shares, nil
}

// wrappedWaitGroup is needed because waitgroups do not expose their internal counter,
// and we don't know in advance how many jobs we will have to wait for.
type wrappedWaitGroup struct {
	wg      sync.WaitGroup
	mu      sync.Mutex
	counter int
}

func (w *wrappedWaitGroup) Add(count int) {
	w.wg.Add(count)
	w.mu.Lock()
	w.counter += count
	w.mu.Unlock()
}

func (w *wrappedWaitGroup) Done() {
	w.wg.Done()
	w.mu.Lock()
	w.counter--
	w.mu.Unlock()
}

func (w *wrappedWaitGroup) Wait() {
	w.wg.Wait()
}

func (w *wrappedWaitGroup) IsWorking() bool {
	w.mu.Lock()
	defer w.mu.Unlock()
	return w.counter > 0
}

type fetchedBounds struct {
	mu      sync.Mutex
	lowest  int
	highest int
}

func (b *fetchedBounds) Update(index int) {
	b.mu.Lock()
	defer b.mu.Unlock()
	if index <= b.lowest {
		b.lowest = index
	}
	// not an `else if` because an element can be both the
	// lower and higher bound
	if index >= b.highest {
		b.highest = index
	}
}

// getLeavesByNamespace returns all the leaves from the given root with the given namespace.ID.
// If nothing is found, it returns both data and error as nil.
// The following implementation is based on `GetShares`.
func getLeavesByNamespace(
	ctx context.Context,
	bGetter blockservice.BlockGetter,
	root cid.Cid,
	nID namespace.ID,
	maxShares int,
) ([]ipld.Node, error) {
	err := SanityCheckNID(nID)
	if err != nil {
		return nil, err
	}

	// we don't know where in the tree the leaves in the namespace are,
	// so we keep track of the bounds to return the correct slice
	// maxShares acts as a sentinel to know if we find any leaves
	bounds := fetchedBounds{sync.Mutex{}, maxShares, 0}

	// buffer the jobs to avoid blocking, we only need as many
	// queued as the number of shares in the second-to-last layer
	jobs := make(chan *job, (maxShares+1)/2)
	jobs <- &job{id: root}

	// the wg counter cannot be pre-allocated either, it is incremented with each job
	wg := wrappedWaitGroup{sync.WaitGroup{}, sync.Mutex{}, 0}
	wg.Add(1)

	// we overallocate space for leaves since we do not know how many we will find
	// on the level above, the length of the Row is passed in as maxShares
	leaves := make([]ipld.Node, maxShares)
	leafLk := sync.Mutex{}

	// if the WaitGroup counter is above 0, we are still walking the tree
	for wg.IsWorking() {
		select {
		case j := <-jobs:
			namespacePool.Submit(func() {
				defer wg.Done()

				rootH := plugin.NamespacedSha256FromCID(j.id)
				if nID.Less(nmt.MinNamespace(rootH, nID.Size())) || !nID.LessOrEqual(nmt.MaxNamespace(rootH, nID.Size())) {
					return
				}

				nd, err := plugin.GetNode(ctx, bGetter, j.id)
				if err != nil {
					return
				}

				links := nd.Links()
				linkCount := len(links)
				// wg counter needs to be incremented **before** adding new jobs
				if linkCount > 1 {
					wg.Add(linkCount)
				} else if linkCount == 1 {
					// successfully fetched a leaf belonging to the namespace
					leafLk.Lock()
					leaves[j.pos] = nd
					leafLk.Unlock()
					// we found a leaf, so we update the bounds
					bounds.Update(j.pos)
					return
				}

				// this node has links in the namespace, so keep walking
				for i, lnk := range links {
					select {
					case jobs <- &job{
						id: lnk.Cid,
						// position represents the index in a flattened binary tree,
						// so we can return a slice of leaves in order
						pos: j.pos*2 + i,
					}:
					case <-ctx.Done():
						return
					}
				}
			})
		case <-ctx.Done():
			return nil, nil
		// default case is needed to prevent blocking
		default:
		}
	}

	// we didn't find any leaves, return nil
	if bounds.lowest == maxShares {
		return nil, nil
	}

	return leaves[bounds.lowest : bounds.highest+1], nil
}
