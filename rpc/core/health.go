package core

import (
	ctypes "github.com/airchains-network/tracksbft/rpc/core/types"
	rpctypes "github.com/airchains-network/tracksbft/rpc/jsonrpc/types"
)

// Health gets node health. Returns empty result (200 OK) on success, no
// response - in case of an error.
// More: https://docs.cometbft.com/v0.34/rpc/#/Info/health
func Health(ctx *rpctypes.Context) (*ctypes.ResultHealth, error) {
	return &ctypes.ResultHealth{}, nil
}