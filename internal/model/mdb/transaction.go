package mdb

import (
	"time"
)

const (
	TransactionStatusInit       = 0 + iota // 已确认
	TransactionStatusSuccess               // 已成功
	TransactionStatusFail                  // 已失败
	TransactionStatusCollecting            // 归集中
	TransactionStatusCollected             // 已归集完成
)

// Transaction 区块链交易信息
type Transaction struct {
	ID              uint       `json:"id" gorm:"primaryKey"`
	OutTradeNo      string     `json:"out_trade_no" gorm:"type:varchar(128);uniqueIndex:out_trade_nos;"` // 外部订单号
	UserId          uint       `json:"user_id" gorm:"uniqueIndex:out_trade_nos;"`                        // 用户ID
	Hash            string     `json:"hash" gorm:"type:varchar(192);"`                                   // 交易哈希
	Chain           string     `json:"chain" gorm:"type:varchar(32);"`                                   // 链简称
	FromAddress     string     `json:"from_address"`                                                     // 发送方地址
	ToAddress       string     `json:"to_address"`                                                       // 接收方地址
	ContractAddress string     `json:"contract_address"`                                                 // 合约地址
	Amount          float64    `json:"amount"`                                                           // 交易金额
	Status          int        `json:"status"`                                                           // 交易状态（0-待确认 1-已确认 2-失败）
	BlockHash       string     `json:"block_hash"`                                                       // 所属区块哈希
	BlockHeight     uint64     `json:"block_height"`                                                     // 所属区块高度
	CreatedAt       *time.Time `json:"created_at"`
	UpdatedAt       *time.Time `json:"updated_at"`
	ConfirmedAt     *time.Time `json:"confirmed_at"` // 确认时间
	// Fee           float64   `json:"fee"`           // 交易手续费
	CallbackUrl string `json:"callback_url"` // 回调URL
}
