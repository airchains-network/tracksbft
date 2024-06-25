package tracks

//fmt.Println("Ethereum Transaction Details:")
//fmt.Println("From:", sender)
//fmt.Println("To:", recipient)
//fmt.Println("Amount:", amount)
//fmt.Println("Gas:", gas)
//fmt.Println("ethereumTxHash:", ethereumTxHash)
//fmt.Println("txHash", txHash)

//type TransactionSecond struct {
//	To                string
//	From              string
//	Amount            string
//	FromBalances      string
//	ToBalances        string
//	TransactionHash   string
//	Messages          string
//	TransactionNonces string
//	AccountNonces     string
//}

type EthTransaction struct {
	From        string
	To          string
	FromCosmos  string
	ToCosmos    string
	Amount      string
	Gas         string
	TxHash      string
	EthTxHash   string
	ToBalance   string
	FromBalance string
	Nonce       string
}
