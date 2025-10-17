package mdb

import (
	"time"
)

type Deposit struct {
	ID              uint       `json:"id"`                                             // 钱包ID
	UserId          uint       `json:"user_id" gorm:"uniqueIndex:uid;"`                // 用户ID
	ContractAddress string     `json:"contract_address" gorm:"uniqueIndex:uid;"`       // 合约地址
	Chain           string     `json:"chain" gorm:"type:varchar(32);uniqueIndex:uid;"` // 链简称
	Amount          float64    `json:"amount"`                                         // 交易金额
	Status          int        `json:"status"`                                         // 钱包状态（0-禁用 1-正常）
	CreatedAt       *time.Time `json:"created_at"`                                     // 创建时间
	UpdatedAt       *time.Time `json:"updated_at"`                                     // 更新时间
}
