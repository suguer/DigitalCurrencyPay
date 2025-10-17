package model

type GetNowBlock struct {
	BlockID     string      `json:blockID`
	BlockHeader BlockHeader `json:"block_header"`
	Error       string      `json:"Error"`
}

type GetBlockByNum struct {
	BlockID      string         `json:"blockID"`
	BlockHeader  BlockHeader    `json:"block_header"`
	Transactions []Transactions `json:"transactions"`
}

type Transactions struct {
	Ret        []Ret    `json:"ret"`
	Signature  []string `json:"signature"`
	TxID       string   `json:"txID"`
	RawData    RawData  `json:"raw_data,omitempty"`
	RawDataHex string   `json:"raw_data_hex"`
}

type Ret struct {
	ContractRet string `json:"contractRet"`
}
type BlockHeader struct {
	RawData          RawData `json:"raw_data"`
	WitnessSignature string  `json:"witness_signature"`
}
type RawData struct {
	Number         int        `json:"number"`
	TxTrieRoot     string     `json:"txTrieRoot"`
	WitnessAddress string     `json:"witness_address"`
	ParentHash     string     `json:"parentHash"`
	Version        int        `json:"version"`
	Timestamp      int64      `json:"timestamp"`
	Contract       []Contract `json:"contract"`
}
type Contract struct {
	Parameter Parameter `json:"parameter"`
	Type      string    `json:"type"`
}
type Parameter struct {
	Value   Value  `json:"value"`
	TypeURL string `json:"type_url"`
}
type Value struct {
	Data            string `json:"data"`
	OwnerAddress    string `json:"owner_address"`
	ContractAddress string `json:"contract_address"`
	Amount          int64  `json:"amount"`
	ToAddress       string `json:"to_address"`
}

// TransactionResponse 表示 TRC20 交易的响应结构体
type TransactionResponse struct {
	Success      bool          `json:"success"` // 请求是否成功
	Transactions []Transaction `json:"data"`    // 交易数据列表
}

// Transaction 表示单个 TRC20 交易的结构体
type Transaction struct {
	Hash      string `json:"hash"`      // 交易哈希
	From      string `json:"from"`      // 发送方地址
	To        string `json:"to"`        // 接收方地址
	Value     string `json:"value"`     // 交易金额
	Timestamp int64  `json:"timestamp"` // 交易时间戳
	Block     int64  `json:"block"`     // 区块号
}
type AccountResource struct {
	FreeNetLimit      int64 `json:"freeNetLimit"`
	TotalNetLimit     int64 `json:"TotalNetLimit"`
	TotalNetWeight    int64 `json:"TotalNetWeight"`
	TotalEnergyLimit  int64 `json:"TotalEnergyLimit"`
	TotalEnergyWeight int64 `json:"TotalEnergyWeight"`
	EnergyLimit       int64 `json:"EnergyLimit"`
	EnergyUsed        int64 `json:"EnergyUsed"`
	NetLimit          int64 `json:"NetLimit"`
	NetUsed           int64 `json:"NetUsed"`
	FreeNetUsed       int64 `json:"freeNetUsed"`
}

type TronTokenTxResponse struct {
	Data    []TronTokenTxTransaction `json:"data"`
	Success bool                     `json:"success"`
}
type TronTokenTxTransaction struct {
	TransactionId  string `json:"transaction_id"`
	BlockTimestamp int64  `json:"block_timestamp"`
	From           string `json:"from"`
	To             string `json:"to"`
	Type           string `json:"type"`
	Value          string `json:"value"`
	TokenInfo      struct {
		Symbol   string `json:"symbol"`
		Address  string `json:"address"`
		Decimals int64  `json:"decimals"`
		Name     string `json:"name"`
	} `json:"token_info"`
}
