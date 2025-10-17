package crontab

import (
	"DigitalCurrency/internal/blockchain"
	"DigitalCurrency/internal/config"
	"DigitalCurrency/internal/constant"
	"DigitalCurrency/internal/logger"
	"DigitalCurrency/internal/model/dao"
	"DigitalCurrency/internal/model/mdb"
	"DigitalCurrency/internal/service/configuration"
	"DigitalCurrency/internal/service/transaction"
	"DigitalCurrency/internal/service/wallet"
	"DigitalCurrency/internal/util"
	"context"
	"fmt"
	"time"

	"github.com/ethereum/go-ethereum/log"
)

func init() {
	crontabs = append(crontabs, Crontab{
		Rule: "10 00 * * *",
		Fun:  TronCollection,
	})

}
func TronCollection(ctx context.Context) {
	ownerAddress, err := configuration.Instance("ownerAddress")
	if err != nil || ownerAddress == "" {
		logger.ErrorLogger.Debug("TronCollection ownerAddress err: %v")
		return
	}
	ownerPrivateKey, _ := configuration.Instance("ownerPrivateKey", "")
	tronMinFreeNet, _ := configuration.Instance("ownerAddress", "400")
	tronMinEnergy, _ := configuration.Instance("ownerAddress", "35000")

	toAddr, _ := util.HexString2Address(ownerAddress)
	tron := blockchain.NewTron(context.Background(), &config.Conf.BlockChain.Tron)
	transactions, err := transaction.Collection(constant.ChainTron)
	if err != nil {
		return
	}
	fmt.Printf("transactions: %+v\n", transactions)
	// 检查地址是否激活
	for _, collection := range transactions {
		AccountResources, err := tron.GetAccountResources(collection.ToAddress)
		fmt.Printf("AccountResources: %+v\n", AccountResources)
		if err != nil {
			continue
		}
		if AccountResources.TotalEnergyLimit == 0 {
			if ownerPrivateKey == "" {
				continue
			} else {
				//转Trx 激活账号,等下一次循环归集
				tron.TransferTrx(ownerPrivateKey, ownerAddress, collection.ToAddress, 1)
				continue
			}
		}

		//判断能量和带宽是否足够
		FreeNet := AccountResources.FreeNetLimit + AccountResources.NetLimit - AccountResources.NetUsed - AccountResources.FreeNetUsed
		if FreeNet < util.StringToInt64(tronMinFreeNet) {
			continue
		}

		Energy := AccountResources.EnergyLimit - AccountResources.EnergyUsed
		// if Energy < 140000 {
		if Energy < util.StringToInt64(tronMinEnergy) {
			continue
		}
		timer := time.Now()
		//生成一条交易记录
		transaction := &mdb.Transaction{
			Chain:           constant.ChainTron,
			ContractAddress: collection.ContractAddress,
			Status:          mdb.TransactionStatusCollecting,
			Amount:          -collection.Amount,
			FromAddress:     collection.ToAddress,
			ToAddress:       ownerAddress,
			CreatedAt:       &timer,
			UpdatedAt:       &timer,
			OutTradeNo:      util.GenerateUUID(),
			UserId:          0,
		}
		dao.Mdb.Save(transaction)

		//合约转账
		walletInstance, err := wallet.InstanceByAddress(collection.ToAddress)
		if err != nil {
			log.Error("TronCollection wallet err: %v", err)
			continue
		}
		contractAddr, _ := util.HexString2Address(collection.ContractAddress)
		ownerAddr, _ := util.HexString2Address(collection.ToAddress)

		txId, err := tron.TransferTRC20(walletInstance.PrivateKeyDecrypt, ownerAddr, contractAddr, toAddr, int64(collection.Amount))
		fmt.Printf("txId: %v\n", txId)
		fmt.Printf("err: %v\n", err)
		if err != nil {
			continue
		}
		transaction.Hash = txId
		transaction.Status = mdb.TransactionStatusCollected
		dao.Mdb.Save(transaction)

	}
}
