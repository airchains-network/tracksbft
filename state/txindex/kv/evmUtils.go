package kv

import (
	"encoding/json"
	"fmt"
	"strconv"
)

const (
	CounterTxsKey      = "countTxs"
	CounterPodsKey     = "countPods"
	transactionPodSize = 25
	rawPodPrefix       = "raw_pod_"
)

func InitiateDatabaseForPods(txi *TxIndex) error {

	fmt.Println("lets test what is this means")
	err := txi.store.Set([]byte(CounterTxsKey), []byte("0"))
	if err != nil {
		fmt.Println("Error initializing countTxs:", err)
		return err
	}
	err = txi.store.Set([]byte(CounterPodsKey), []byte("1"))
	if err != nil {
		fmt.Println("Error initializing countPods:", err)
		return err
	}
	fmt.Println("Initialized database for pods")
	return nil
}

func IncrementTxCount(txi *TxIndex) error {
	count, err := RetrieveTxCount(txi)
	if err != nil {
		return err
	}
	count++
	err = txi.store.Set([]byte(CounterTxsKey), []byte(strconv.Itoa(count)))
	if err != nil {
		fmt.Println("Error incrementing countTxs:", err)
		return err
	}
	return nil
}

func RetrieveTxCount(txStore *TxIndex) (int, error) {
	store := txStore.store
	byteRes, err := store.Get([]byte(CounterTxsKey))
	if err != nil {
		return 0, err
	}
	if byteRes == nil {
		return 0, fmt.Errorf("database not initiated")
	}
	strCount := string(byteRes)
	count, err := strconv.Atoi(strCount)
	if err != nil {
		return 0, fmt.Errorf("error converting count from string to int: %w", err)
	}
	return count, nil
}

func SetTxCount(txStore *TxIndex, count int) error {
	store := txStore.store
	err := store.Set([]byte(CounterTxsKey), []byte(strconv.Itoa(count)))
	if err != nil {
		return err
	}
	return nil
}

func IncrementPodCount(txi *TxIndex) error {
	count, err := RetrievePodCount(txi)
	if err != nil {
		return err
	}
	count++
	err = txi.store.Set([]byte(CounterPodsKey), []byte(strconv.Itoa(count)))
	if err != nil {
		fmt.Println("Error incrementing countPods:", err)
		return err
	}
	return nil
}

func RetrievePodCount(txStore *TxIndex) (int, error) {
	store := txStore.store
	byteRes, err := store.Get([]byte(CounterPodsKey))
	if err != nil {
		return 0, err
	}
	if byteRes == nil {
		return 0, fmt.Errorf("database not initiated")
	}
	strCount := string(byteRes)
	count, err := strconv.Atoi(strCount)
	if err != nil {
		return 0, fmt.Errorf("error converting count from string to int: %w", err)
	}
	return count, nil
}

// SetPod stores a pod of transactions as a single byte slice
func SetPod(txStore *TxIndex, podNumber int, pod [][]byte) error {
	store := txStore.store
	key := []byte(fmt.Sprintf("%s%d", rawPodPrefix, podNumber))

	// Serialize the pod
	podData, err := json.Marshal(pod)
	if err != nil {
		return fmt.Errorf("error serializing pod: %w", err)
	}

	// Store the serialized pod
	err = store.Set(key, podData)
	if err != nil {
		return fmt.Errorf("error storing pod: %w", err)
	}

	return nil
}

// GetPod retrieves and deserializes a pod of transactions
func GetPod(txStore *TxIndex, podNumber int) ([][]byte, error) {
	store := txStore.store
	key := []byte(fmt.Sprintf("%s%d", rawPodPrefix, podNumber))

	// Retrieve the serialized pod
	byteRes, err := store.Get(key)
	if err != nil {
		return nil, fmt.Errorf("error retrieving pod: %w", err)
	}
	if byteRes == nil {
		return nil, fmt.Errorf("pod not found")
	}

	// Deserialize the pod
	var pod [][]byte
	err = json.Unmarshal(byteRes, &pod)
	if err != nil {
		return nil, fmt.Errorf("error deserializing pod: %w", err)
	}

	return pod, nil
}

// GetNonce retrieves the nonce for a given address. if not exists 0
func GetNonce(txStore *TxIndex, address string) (uint64, error) {
	store := txStore.store
	key := []byte("nonce_" + address)

	// Retrieve the nonce
	byteRes, err := store.Get(key)
	if err != nil {
		return 0, fmt.Errorf("error retrieving nonce: %w", err)
	}
	if byteRes == nil {
		err := SetNonce(txStore, address, 0)
		if err != nil {
			return 0, err
		}
		return 0, nil
	}

	// Deserialize the nonce
	nonce, err := strconv.ParseUint(string(byteRes), 10, 64)
	if err != nil {
		return 0, fmt.Errorf("error deserializing nonce: %w", err)
	}

	return nonce, nil
}

// SetNonce stores the nonce for a given address
func SetNonce(txStore *TxIndex, address string, nonce uint64) error {
	store := txStore.store
	key := []byte("nonce_" + address)

	// Serialize the nonce
	nonceData := []byte(strconv.FormatUint(nonce, 10))

	// Store the serialized nonce
	err := store.Set(key, nonceData)
	if err != nil {
		return fmt.Errorf("error storing nonce: %w", err)
	}

	return nil
}
