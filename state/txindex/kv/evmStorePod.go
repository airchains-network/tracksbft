package kv

import (
	"encoding/json"
	"fmt"
	"strconv"

	abci "github.com/airchains-network/tracksbft/abci/types"
	txindex "github.com/airchains-network/tracksbft/state/txindex"
	tracksTypes "github.com/airchains-network/tracksbft/types/tracks"
)

func extractAttribute(attributes []abci.EventAttribute, key string) string {
	for _, attr := range attributes {
		if string(attr.Key) == key {
			return string(attr.Value)
		}
	}
	return ""
}

func StorePod(txi *TxIndex, b *txindex.Batch) error {

	storeBatch := txi.store.NewBatch()
	defer storeBatch.Close()
	if len(b.Ops) == 0 {
		fmt.Println("no operations provided")
		return nil
	}

	var (
		blockTxCount = 0
		ethTxArray   []tracksTypes.EthTransaction
	)

	// Retrieve current counts from the database
	currentTxCount, err := RetrieveTxCount(txi)
	if err != nil {
		fmt.Println("Error retrieving transaction count:", err)
		return err
	}

	currentPodCount, err := RetrievePodCount(txi)
	if err != nil {
		fmt.Println("Error retrieving pod count:", err)
		return err
	}

	// Initialize a slice to hold transactions for the current pod
	var currentPodTxs [][]byte
	currentPodTxs, err = GetPod(txi, currentPodCount)
	if err != nil && err.Error() != "pod not found" {
		fmt.Println("Error retrieving the latest pod:", err)
		return err
	}

	for _, result := range b.Ops {
		var ethereumTxHash, txHash, module, recipient, sender, amount, gas, recipientCosmos, senderCosmos, fromBalance, toBalance string
		events := result.Result.Events

		isEthTx := false
		for _, event := range events {
			switch event.Type {
			case "message":
				for _, attribute := range event.Attributes {
					if string(attribute.Key) == "action" {
						if string(attribute.Value) == "/ethermint.evm.v1.MsgEthereumTx" {
							isEthTx = true
							blockTxCount++
						}
					}
				}
			}
		}

		if isEthTx == true {

			for _, event := range events {

				switch event.Type {

				case "ethereum_tx":
					if len(event.Attributes) == 6 {
						ethereumTxHash = extractAttribute(event.Attributes, "ethereumTxHash")
						txHash = extractAttribute(event.Attributes, "txHash")
						recipient = extractAttribute(event.Attributes, "recipient")
						amount = extractAttribute(event.Attributes, "amount")
						gas = extractAttribute(event.Attributes, "txGasUsed")
					}
				case "transfer":
					recipientCosmos = extractAttribute(event.Attributes, "recipient")
					senderCosmos = extractAttribute(event.Attributes, "sender")
				case "message":
					module = extractAttribute(event.Attributes, "module")
					if module == "evm" {
						sender = extractAttribute(event.Attributes, "sender")
					}
				}
			}

			previousBlockHeight := result.Height - 1
			fromBalance, err = GetBalanceAtHeight(sender, previousBlockHeight)
			if err != nil {
				fmt.Println("Error checking sender balance:", err)
				return err
			}

			toBalance, err = GetBalanceAtHeight(recipient, previousBlockHeight)
			if err != nil {
				fmt.Println("Error checking recipient balance:", err)
				return err
			}

			// get nonce
			fromNonce, err := GetNonce(txi, sender)
			if err != nil {
				fmt.Println("Error retrieving nonce:", err)
				return err
			}
			fromNonce++
			err = SetNonce(txi, sender, fromNonce)
			if err != nil {
				fmt.Println("Error setting nonce:", err)
				return err
			}

			ethTx := tracksTypes.EthTransaction{
				From:        sender,
				To:          recipient,
				FromCosmos:  senderCosmos,
				ToCosmos:    recipientCosmos,
				Amount:      amount,
				Gas:         gas,
				TxHash:      txHash,
				EthTxHash:   ethereumTxHash,
				Nonce:       strconv.FormatUint(fromNonce, 10),
				FromBalance: fromBalance,
				ToBalance:   toBalance,
			}

			//fmt.Println("Ethereum Transaction Details:", ethTx)
			ethTxArray = append(ethTxArray, ethTx)
		}
	}

	fmt.Println("currentTxCount=", currentTxCount, ", currentPodCount=", currentPodCount, ", blockTxCount=", len(ethTxArray), ", currentPodTxCount=", len(currentPodTxs))

	for _, ethTx := range ethTxArray {

		// Serialize the transaction and add it to the current pod
		serializedTx, err := serializeEthTransaction(ethTx)
		if err != nil {
			fmt.Println("Error serializing Ethereum transaction:", err)
			return err
		}

		currentPodTxs = append(currentPodTxs, serializedTx)
		currentTxCount++

		// Check if the current pod has reached the maximum size
		if len(currentPodTxs) == transactionPodSize {

			// Store the full pod
			err = SetPod(txi, currentPodCount, currentPodTxs)
			if err != nil {
				fmt.Println("Error storing pod:", err)
				return err
			}

			// Increment the pod count and reset the current pod
			currentPodCount++
			currentPodTxs = nil

			// Increment the pod count in the database
			err = IncrementPodCount(txi)
			if err != nil {
				fmt.Println("Error incrementing pod count:", err)
				return err
			}
		}
	}

	// Store any remaining transactions that didn't fill a full pod
	if len(currentPodTxs) > 0 {
		err = SetPod(txi, currentPodCount, currentPodTxs)
		if err != nil {
			fmt.Println("Error storing the final pod:", err)
			return err
		}
	}

	// Update the transaction count in the database
	err = SetTxCount(txi, currentTxCount)
	if err != nil {
		fmt.Println("Error setting transaction count:", err)
		return err
	}

	return nil

}

func serializeEthTransaction(ethTx tracksTypes.EthTransaction) ([]byte, error) {
	txByte, err := json.Marshal(ethTx)
	if err != nil {
		return nil, err
	}
	return txByte, nil

}
