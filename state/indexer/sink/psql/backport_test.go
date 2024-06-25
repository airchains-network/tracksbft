package psql

import (
	"github.com/airchains-network/tracksbft/state/indexer"
	"github.com/airchains-network/tracksbft/state/txindex"
)

var (
	_ indexer.BlockIndexer = BackportBlockIndexer{}
	_ txindex.TxIndexer    = BackportTxIndexer{}
)
