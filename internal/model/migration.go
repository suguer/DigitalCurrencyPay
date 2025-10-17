package model

import (
	"DigitalCurrency/internal/config"
	"DigitalCurrency/internal/model/dao"
	"DigitalCurrency/internal/model/mdb"
	"DigitalCurrency/internal/util"
	"fmt"
)

func Migration() {
	err := dao.Mdb.AutoMigrate(
		&mdb.Transaction{},
		&mdb.Wallet{},
		&mdb.Deposit{},
		&mdb.User{},
		// &mdb.Notify{},
		&mdb.Configuration{},
	)
	if err != nil {
		fmt.Printf("err: %v\n", err)
	}
	insertDefaultData()
}

func insertDefaultData() {
	var UserCount int64
	dao.Mdb.Model(&mdb.User{}).Count(&UserCount)
	Secret := util.GenerateUUID()
	if UserCount == 0 {
		var user = mdb.User{
			Username: "admin",
			Password: util.Md5BySalt("admin", config.Conf.Salt),
			Secret:   Secret,
		}
		dao.Mdb.Create(&user)
	}

	var ConfigurationCount int64
	dao.Mdb.Model(&mdb.Configuration{}).Count(&ConfigurationCount)
	if ConfigurationCount == 0 {
		var configurations = []mdb.Configuration{
			{
				Key:    "mode",
				Value:  "single", // 单用户还是多用户分销模式
				Remark: "单用户模式下，所有用户的交易都由该用户处理",
			},
			{
				Key:    "ownerAddress", // 所有者地址
				Value:  "",
				Remark: "所有者地址，用于归集接收所有钱包的金额",
			},
			{
				Key:    "ownerPrivateKey", //所有者密钥
				Value:  "",
				Remark: "所有者密钥，用于签名交易",
			},
			{
				Key:    "tronMinFreeNet", // 最小免费网络资源
				Value:  "400",
				Remark: "波场最小免费网络资源，用于交易",
			},
			{
				Key:    "tronMinEnergy", // 最小能量资源
				Value:  "140000",
				Remark: "波场最小能量资源，用于交易",
			},
		}
		dao.Mdb.CreateInBatches(configurations, 100)
		defaultMessage := fmt.Sprintf("账号:admin\n密码:admin\nSecret:%s\n", Secret)
		//保存到 config.Conf.Storage.Path/defaultMessage.txt
		err := util.WriteFile(fmt.Sprintf("%s/defaultMessage.txt", config.Conf.Storage.Path), defaultMessage)
		if err != nil {
			fmt.Printf("err: %v\n", err)
		}

	}
}
