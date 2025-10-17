package transaction

import (
	"DigitalCurrency/internal/blockchain"
	"DigitalCurrency/internal/constant"
	"DigitalCurrency/internal/model/cache"
	"DigitalCurrency/internal/model/dao"
	"DigitalCurrency/internal/model/mdb"
	"DigitalCurrency/internal/service/deposit"
	"DigitalCurrency/internal/service/wallet"
	"DigitalCurrency/internal/util"
	"context"
	"errors"
	"math"
	"sort"
	"strconv"
	"time"

	"github.com/fbsobreira/gotron-sdk/pkg/address"
	"github.com/shopspring/decimal"
)

func Create(tx *mdb.Transaction) (*mdb.Transaction, error) {
	timer := time.Now()
	data := mdb.Transaction{
		Chain:           tx.Chain,
		ContractAddress: tx.ContractAddress,
		Status:          0,
		Amount:          tx.Amount,
		CreatedAt:       &timer,
		UpdatedAt:       &timer,
		OutTradeNo:      tx.OutTradeNo,
		UserId:          tx.UserId,
		CallbackUrl:     tx.CallbackUrl,
	}
	contractAddr, err := constant.GetContractAddress(tx.Chain, tx.ContractAddress)
	if err == nil {
		data.ContractAddress = contractAddr
	} else if tx.Chain == constant.ChainTron || tx.Chain == constant.ChainTronShasta {
		contractAddr, err := address.Base58ToAddress(tx.ContractAddress)
		if err == nil {
			data.ContractAddress = contractAddr.Hex()
		}
	}

	//获取一个有效地址
	wallet, err := wallet.GetAvailableAddress(tx.Chain)
	if err != nil {
		return nil, err
	}
	data.ToAddress = wallet.Address
	result := dao.Mdb.Create(&data)
	cache.TransactionCacheSet(&data)
	return &data, result.Error
}

func Instance(Id uint) (*mdb.Transaction, error) {
	var tx mdb.Transaction
	err := dao.Mdb.Where("id = ?", Id).First(&tx).Error
	if err != nil {
		return nil, err
	}
	return &tx, nil
}

func InstanceByOutTradeNo(outTradeNo string, userId uint) (*mdb.Transaction, error) {
	var tx mdb.Transaction
	model := dao.Mdb.Where("out_trade_no = ? ", outTradeNo)
	if userId > 0 {
		model = model.Where("user_id = ?", userId)
	}
	err := model.First(&tx).Error
	if err != nil {
		return nil, err
	}
	return &tx, nil
}

/*
获取需要归集的钱包和总金额
*/
// CollectionItem 表示归集结果的单个项目
type CollectionItem struct {
	ToAddress       string  `json:"to_address"`       // 接收方地址
	ContractAddress string  `json:"contract_address"` // 合约地址
	Amount          float64 `json:"amount"`           // 总金额
}

func Collection(chain string) ([]CollectionItem, error) {
	var transactions *[]mdb.Transaction
	err := dao.Mdb.Where("chain = ? and status = ?", chain, mdb.TransactionStatusSuccess).Find(&transactions).Error
	if err != nil {
		return nil, err
	}

	// 使用临时map来累加金额
	tempMap := make(map[string]map[string]float64)

	for _, tx := range *transactions {
		// 如果地址不存在于结果集中，则创建新的映射
		if _, exists := tempMap[tx.ToAddress]; !exists {
			tempMap[tx.ToAddress] = make(map[string]float64)
		}

		// 累加该合约地址的金额
		tempMap[tx.ToAddress][tx.ContractAddress] += tx.Amount
	}

	// 将map转换为切片
	var results []CollectionItem
	for toAddress, contracts := range tempMap {
		for contractAddress, amount := range contracts {
			results = append(results, CollectionItem{
				ToAddress:       toAddress,
				ContractAddress: contractAddress,
				Amount:          amount,
			})
		}
	}

	// 按照 Amount 从大到小排序
	sort.Slice(results, func(i, j int) bool {
		return results[i].Amount > results[j].Amount
	})

	return results, nil
}

func TransactionRepair(out_trade_no string, user_id uint) error {
	var tx mdb.Transaction
	err := dao.Mdb.Where("out_trade_no = ?", out_trade_no).First(&tx).Error
	if err != nil {
		return err
	}
	if tx.Status != mdb.TransactionStatusInit && tx.Status != mdb.TransactionStatusFail {
		return errors.New("transaction status error")
	}
	client := blockchain.Factory(context.Background(), tx.Chain)
	txs, err := client.TokenTx(tx.ToAddress, tx.ContractAddress)
	if err != nil {
		return err
	}
	for _, item := range txs.Result {
		if item.Hash != "28f9901520c8180809cf6c529d3b715e7706c17e0b5a1186aeb312059c1c8583" {
			continue
		}
		if !util.MatchAddress(item.To, tx.ToAddress) {
			continue
		}
		TimeStamp, err := strconv.ParseInt(item.TimeStamp, 10, 64)
		if err != nil {
			continue
		}
		timestamp := time.Unix(TimeStamp, 0)
		if timestamp.Before(*tx.CreatedAt) || timestamp.After(tx.CreatedAt.Add(30*time.Minute)) {
			continue
		}
		precision, err := strconv.Atoi(item.TokenDecimal)
		if err != nil {
			continue
		}
		amount, err := strconv.ParseFloat(item.Value, 64)
		if err != nil {
			continue
		}
		amount = amount / math.Pow10(precision)
		scopeMin := decimal.NewFromFloat(tx.Amount * (1 - 0.001)).Round(2)
		scopeMax := decimal.NewFromFloat(tx.Amount * (1 + 0.001)).Round(2)

		if amount < scopeMin.InexactFloat64() || amount > scopeMax.InexactFloat64() {
			continue
		}
		if amount != tx.Amount {
			OutTradeNo := tx.OutTradeNo
			tx.OutTradeNo = "drop-" + tx.OutTradeNo
			dao.Mdb.Save(&tx)

			//创建新订单
			transaction, err := Create(&mdb.Transaction{
				OutTradeNo:      OutTradeNo,
				UserId:          tx.UserId,
				Hash:            item.Hash,
				Chain:           tx.Chain,
				FromAddress:     tx.FromAddress,
				ToAddress:       tx.ToAddress,
				ContractAddress: tx.ContractAddress,
				Amount:          amount,
				Status:          0,
			})
			if err != nil {
				return err
			}
			tx = *transaction
		}
		tx.Status = 1
		tx.Hash = item.Hash
		tx.ConfirmedAt = util.Now()
		tx.FromAddress = item.From
		dao.Mdb.Save(&tx)
		deposit.Increment(tx.Chain, tx.ContractAddress, tx.Amount, tx.UserId)
		return nil

	}
	return errors.New("not found")
}
