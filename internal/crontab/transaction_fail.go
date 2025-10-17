package crontab

import (
	"DigitalCurrency/internal/logger"
	"DigitalCurrency/internal/model/cache"
	"DigitalCurrency/internal/model/dao"
	"DigitalCurrency/internal/model/mdb"
	"DigitalCurrency/internal/util"
	"context"
	"time"

	"go.uber.org/zap"
)

func init() {
	crontabs = append(crontabs, Crontab{
		Rule: "*/5 * * * *",
		Fun:  TransactionFail,
	})

}

func TransactionFail(ctx context.Context) {
	//获取超时30分钟的交易,设置为失败,并清空关联的缓存cache
	var list []mdb.Transaction
	dao.Mdb.Where("status = ? AND created_at < ?", mdb.TransactionStatusInit, time.Now().Add(-30*time.Minute)).Find(&list)
	for _, transaction := range list {
		cache.TransactionCacheDelete(&transaction)

		logger.Logger.Info("交易超时,已更新缓存", zap.Any("transaction", transaction))

		var wallet mdb.Wallet
		dao.Mdb.Where("address = ?", transaction.ToAddress).First(&wallet)
		if wallet.Status == mdb.WalletStatusUsed {
			wallet.Status = mdb.WalletStatusActive
			wallet.UpdatedAt = util.Now()
			dao.Mdb.Save(&wallet)
			logger.Logger.Info("钱包已恢复有效", zap.Any("wallet", wallet))
		}

		transaction.Status = mdb.TransactionStatusFail
		transaction.UpdatedAt = util.Now()
		dao.Mdb.Save(&transaction)
	}

	//获取超时30分钟的使用中的钱包,恢复有效
	var listWallet []mdb.Wallet
	dao.Mdb.Where("status = ? AND updated_at < ?", mdb.WalletStatusUsed, time.Now().Add(-30*time.Minute)).Find(&listWallet)
	for _, wallet := range listWallet {
		wallet.Status = mdb.WalletStatusActive
		wallet.UpdatedAt = util.Now()
		dao.Mdb.Save(&wallet)
		logger.Logger.Info("钱包已恢复有效", zap.Any("wallet", wallet))
	}

}
