package model

import (
	"math/big"
	"time"
)

type EthLog struct {
	TxHash          string
	Amount          uint
	FromAddress     string
	ToAddress       string
	ContractAddress string
	BlockNumber     uint64
	BlockTimestamp  *time.Time
	BlockHash       string
}
type EthGasFee struct {
	BaseFeePerGas        float64
	MaxFeePerGas         *big.Int
	MaxPriorityFeePerGas float64
	GasLimit             uint64
	RequiredAmount       float64
}

type TokenTxResponse struct {
	Status  string               `json:"status"`
	Message string               `json:"message"`
	Result  []TokenTxTransaction `json:"result"`
}
type TokenTxTransaction struct {
	BlockNumber       string `json:"blockNumber"`
	TimeStamp         string `json:"timeStamp"`
	Hash              string `json:"hash"`
	Nonce             string `json:"nonce"`
	BlockHash         string `json:"blockHash"`
	From              string `json:"from"`
	ContractAddress   string `json:"contractAddress"`
	To                string `json:"to"`
	Value             string `json:"value"`
	TokenName         string `json:"tokenName"`
	TokenSymbol       string `json:"tokenSymbol"`
	TokenDecimal      string `json:"tokenDecimal"`
	TransactionIndex  string `json:"transactionIndex"`
	Gas               string `json:"gas"`
	GasPrice          string `json:"gasPrice"`
	GasUsed           string `json:"gasUsed"`
	CumulativeGasUsed string `json:"cumulativeGasUsed"`
	Input             string `json:"input"`
	MethodId          string `json:"methodId"`
	FunctionName      string `json:"functionName"`
	Confirmations     string `json:"confirmations"`
}
