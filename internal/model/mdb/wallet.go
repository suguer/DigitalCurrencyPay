package mdb

import (
	"DigitalCurrency/internal/util"
	"time"

	"gorm.io/gorm"
)

var (
	encryptKey = []byte("ZmUbqPipJ0Pr7tGmkHDBazjpbKjFZc7S")
)

const (
	WalletStatusActive = iota + 1
	WalletStatusUsed
	WalletStatusStop
)

// Wallet 钱包信息
type Wallet struct {
	ID                uint       `json:"id" gorm:"primaryKey"`              // 钱包ID
	Address           string     `json:"address"  gorm:"uniqueIndex:addr;"` // 钱包地址
	Type              string     `json:"type"`
	PrivateKey        string     `json:"-"`                            // 密钥（加密存储）
	Status            int        `json:"status"`                       // 钱包状态（0-禁用 1-正常）
	CreatedAt         *time.Time `json:"created_at"`                   // 创建时间
	UpdatedAt         *time.Time `json:"updated_at"`                   // 更新时间
	LastAt            *time.Time `json:"last_at"`                      // 最后活跃时间
	PrivateKeyDecrypt string     `gorm:"-" json:"private_key_decrypt"` // 解密后的私钥
}

// BeforeSave GORM的保存前钩子，用于加密私钥
func (w *Wallet) BeforeCreate(tx *gorm.DB) error {
	if w.PrivateKey == "" {
		return nil
	}

	encrypted, err := util.AESEncrypt(encryptKey, w.PrivateKey)
	if err != nil {
		return err
	}

	w.PrivateKey = encrypted
	return nil
}

// AfterFind GORM的查询后钩子，用于解密私钥
func (w *Wallet) AfterFind(tx *gorm.DB) error {
	if w.PrivateKey == "" {
		return nil
	}

	decrypted, err := util.AESDecrypt(encryptKey, w.PrivateKey)
	if err != nil {
		return err
	}

	w.PrivateKeyDecrypt = decrypted
	return nil
}
