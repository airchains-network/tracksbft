package kv

import (
	"context"
	"github.com/ethereum/go-ethereum/common"
	"github.com/ethereum/go-ethereum/ethclient"
	"log"
	"math/big"
	"time"
)

var (
	rpcUrl    = "http://0.0.0.0:8545" // Replace with your Infura project ID or your Ethereum node URL
	ethClient *ethclient.Client
	err       error
)

func init() {
	ethClient, err = ethclient.Dial(rpcUrl)
	if err != nil {
		log.Fatalf("Failed to connect to the Ethereum client: %v", err)
	}
}

func GetBalanceAtHeight(address string, height int64) (string, error) {
	account := common.HexToAddress(address)
	heightBigInt := big.NewInt(0).SetInt64(height) // *big.Int representation

	var balance *big.Int
	var err error
	for {
		balance, err = ethClient.BalanceAt(context.Background(), account, heightBigInt)
		if err == nil {
			break
		}

		// Sleep for a while before retrying to prevent overwhelming the service
		time.Sleep(time.Second * 5)
	}

	return balance.String(), nil
}
