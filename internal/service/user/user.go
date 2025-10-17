package user

import (
	"DigitalCurrency/internal/config"
	"DigitalCurrency/internal/model/dao"
	"DigitalCurrency/internal/model/mdb"
	"DigitalCurrency/internal/util"
)

func Create(user *mdb.User) error {
	if user.Password == "" {
		user.Password = "123456"
	}
	user.Password = util.Md5BySalt(user.Password, config.Conf.Salt)
	return dao.Mdb.Create(user).Error
}

func Instance(id any) (*mdb.User, error) {
	var user mdb.User
	err := dao.Mdb.Where("id = ?", id).First(&user).Error
	if err != nil {
		return nil, err
	}
	return &user, nil
}

func Index(page int64, pageSize int64, condition map[string]any) ([]*mdb.User, *mdb.Pagination, error) {
	var list []*mdb.User
	var total int64
	query := dao.Mdb.Model(&mdb.User{})

	err := query.Count(&total).Error
	if err != nil {
		return nil, nil, err
	}
	err = query.Offset(int((page - 1) * pageSize)).Limit(int(pageSize)).Order("id desc").Find(&list).Error
	if err != nil {
		return nil, nil, err
	}
	pagination := &mdb.Pagination{
		Current:  page,
		PageSize: pageSize,
		Total:    total,
	}
	return list, pagination, nil
}

func Update(id uint, data map[string]any) error {
	return dao.Mdb.Model(&mdb.User{}).Where("id = ?", id).Updates(data).Error
}

func GetSecret(id any) (*mdb.User, error) {
	var user mdb.User
	err := dao.Mdb.Where("id = ?", id).Select("id", "secret").First(&user).Error
	if err != nil {
		return nil, err
	}
	return &user, nil
}
