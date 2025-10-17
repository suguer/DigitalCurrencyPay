package wallet

import (
	"DigitalCurrency/internal/model/dao"
	"DigitalCurrency/internal/model/mdb"
	"DigitalCurrency/internal/util"
	"errors"

	"gorm.io/gorm"
)

func Index(page int64, pageSize int64, condition map[string]any) ([]*mdb.Wallet, *mdb.Pagination, error) {
	var walletList []*mdb.Wallet
	var total int64
	query := dao.Mdb.Model(&mdb.Wallet{})

	err := query.Count(&total).Error
	if err != nil {
		return nil, nil, err
	}
	err = query.Offset(int((page - 1) * pageSize)).Limit(int(pageSize)).Order("id desc").Find(&walletList).Error
	if err != nil {
		return nil, nil, err
	}
	pagination := &mdb.Pagination{
		Current:  page,
		PageSize: pageSize,
		Total:    total,
	}
	return walletList, pagination, nil
}

func GetAvailableAddress(chain string) (wallet *mdb.Wallet, err error) {
	result := dao.Mdb.Where(" status = ?", mdb.WalletStatusActive).First(&wallet)

	if result.Error != nil && !errors.Is(result.Error, gorm.ErrRecordNotFound) {
		return nil, result.Error
	}
	if errors.Is(result.Error, gorm.ErrRecordNotFound) {
		// 生成新钱包
		wallet, err = CreateWallet()
		if err != nil {
			return nil, err
		}
	}
	wallet.Status = mdb.WalletStatusUsed
	wallet.LastAt = util.Now()
	dao.Mdb.Save(wallet)
	return wallet, nil

}

func CreateWallet() (*mdb.Wallet, error) {
	address, key, err := util.Generate()
	if err != nil {
		return nil, err
	}

	var wallet = &mdb.Wallet{
		Address:    address,
		PrivateKey: key,
		CreatedAt:  util.Now(),
		UpdatedAt:  util.Now(),
		LastAt:     util.Now(),
	}
	err = dao.Mdb.Create(wallet).Error
	return wallet, err
}

func Release(address string) error {
	var wallet mdb.Wallet
	err := dao.Mdb.Where("address = ?", address).First(&wallet).Error
	if err != nil {
		return err
	}
	wallet.Status = mdb.WalletStatusActive
	dao.Mdb.Save(&wallet)
	return nil
}

func InstanceByAddress(address string) (*mdb.Wallet, error) {
	var wallet mdb.Wallet
	err := dao.Mdb.Where("address = ?", address).First(&wallet).Error
	if err != nil {
		return nil, err
	}
	return &wallet, nil
}
