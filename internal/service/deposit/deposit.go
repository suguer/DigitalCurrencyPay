package deposit

import (
	"DigitalCurrency/internal/model/dao"
	"DigitalCurrency/internal/model/mdb"
	"DigitalCurrency/internal/util"
	"errors"

	"gorm.io/gorm"
)

func DepositList(userId uint) ([]*mdb.Deposit, error) {
	var depositList []*mdb.Deposit
	err := dao.Mdb.Where("user_id = ?", userId).Find(&depositList).Error
	if err != nil {
		return nil, err
	}
	return depositList, nil
}

func Increment(chain, contractAddress string, amount float64, userId uint) error {
	var depositInstance mdb.Deposit
	err := dao.Mdb.Where("chain = ? AND contract_address = ?", chain, contractAddress).First(&depositInstance).Error
	if err != nil && errors.Is(err, gorm.ErrRecordNotFound) {
		depositInstance = mdb.Deposit{
			Chain:           chain,
			ContractAddress: contractAddress,
			Amount:          0,
			UserId:          userId,
		}
		dao.Mdb.Create(&depositInstance)
	} else if err != nil {
		return err
	}
	depositInstance.Amount += amount
	dao.Mdb.Save(&depositInstance)
	return nil
}

func Withdraw(toAddress string, chain, contractAddress string, amount float64, userId uint) {
	newAmount := amount
	fee := 0.0
	// 生成2条Transaction,一条是提现记录,一条是手续费
	data := mdb.Transaction{
		Chain:           chain,
		ContractAddress: contractAddress,
		Status:          0,
		Amount:          -newAmount,
		CreatedAt:       util.Now(),
		OutTradeNo:      util.GenerateUUID(),
		UserId:          userId,
		ToAddress:       toAddress,
	}
	dao.Mdb.Create(&data)
	if fee > 0 {
		feeTransaction := mdb.Transaction{
			Chain:           chain,
			ContractAddress: contractAddress,
			Status:          1,
			Amount:          -fee,
			CreatedAt:       util.Now(),
			OutTradeNo:      util.GenerateUUID(),
			UserId:          userId,
			ToAddress:       "",
		}
		dao.Mdb.Create(&feeTransaction)
	}
}
