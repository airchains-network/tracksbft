package core

import (
	"encoding/json"
	"strconv"

	rpctypes "github.com/airchains-network/tracksbft/rpc/jsonrpc/types"
	tracksTypes "github.com/airchains-network/tracksbft/types/tracks"
)

func TracksGetPodCount(_ *rpctypes.Context) (int, error) {
	res, err := env.TxIndexer.GetbytedataFortracks([]byte("countPods"))
	if err != nil {
		return 0, err

	}
	podCount, err := strconv.Atoi(string(res)) // []byte("1")
	if err != nil {
		return 0, err
	}
	return podCount, nil
}

func TracksGetPodTxs(_ *rpctypes.Context, podNumber int) ([]tracksTypes.EthTransaction, error) {
	// podNumber to string
	//strPod := fmt.Sprintf("Pod %d", podNumber)
	env.Logger.Info("Tracks API Request", "req", "tracks_get_pod")

	byteRes, err := env.TxIndexer.GetbytedataFortracks([]byte("raw_pod_" + strconv.Itoa(podNumber)))
	if err != nil {
		return nil, err
	}

	// Deserialize the pod into []tracksTypes.EthTransaction
	var transactionsArrayByte [][]byte
	err = json.Unmarshal(byteRes, &transactionsArrayByte)
	if err != nil {
		return nil, err
	}

	var txArray []tracksTypes.EthTransaction
	for _, txByte := range transactionsArrayByte {
		var tx tracksTypes.EthTransaction
		err = json.Unmarshal(txByte, &tx)
		if err != nil {
			return nil, err
		}
		txArray = append(txArray, tx)
	}

	return txArray, nil
}
