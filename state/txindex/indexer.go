package txindex

import (
	"context"
	"errors"

	abci "github.com/airchains-network/tracksbft/abci/types"
	"github.com/airchains-network/tracksbft/libs/pubsub/query"
)

// XXX/TODO: These types should be moved to the indexer package.

//go:generate ../../scripts/mockery_generate.sh TxIndexer

// TxIndexer interface defines methods to index and search transactions.
type TxIndexer interface {
	AddPod(b *Batch, stationType string) error
	AddBatch(b *Batch) error
	GetbytedataFortracks(hash []byte) ([]byte, error)
	Index(result *abci.TxResult) error
	Get(hash []byte) (*abci.TxResult, error)
	Search(ctx context.Context, q *query.Query) ([]*abci.TxResult, error)
}

// Batch groups together multiple Index operations to be performed at the same time.
// NOTE: Batch is NOT thread-safe and must not be modified after starting its execution.
type Batch struct {
	Ops []*abci.TxResult
}

// NewBatch creates a new Batch.
func NewBatch(n int64) *Batch {
	return &Batch{
		Ops: make([]*abci.TxResult, n),
	}
}

// Add or update an entry for the given result.Index.
func (b *Batch) Add(result *abci.TxResult) error {
	b.Ops[result.Index] = result
	return nil
}

// Size returns the total number of operations inside the batch.
func (b *Batch) Size() int {
	return len(b.Ops)
}

// ErrorEmptyHash indicates empty hash
var ErrorEmptyHash = errors.New("transaction hash cannot be empty")
