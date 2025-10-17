package configuration

import (
	"DigitalCurrency/internal/model/dao"
	"DigitalCurrency/internal/model/mdb"
)

func Index(page int64, pageSize int64, condition map[string]any) ([]*mdb.Configuration, *mdb.Pagination, error) {
	var configurationList []*mdb.Configuration
	var total int64
	query := dao.Mdb.Model(&mdb.Configuration{})

	err := query.Count(&total).Error
	if err != nil {
		return nil, nil, err
	}
	err = query.Offset(int((page - 1) * pageSize)).Limit(int(pageSize)).Order("id desc").Find(&configurationList).Error
	if err != nil {
		return nil, nil, err
	}
	pagination := &mdb.Pagination{
		Current:  page,
		PageSize: pageSize,
		Total:    total,
	}
	return configurationList, pagination, nil
}

func Instance(key string, defaultValue ...string) (value string, err error) {
	var configuration mdb.Configuration
	err = dao.Mdb.Where("key = ?", key).First(&configuration).Error
	if err != nil {
		if len(defaultValue) > 0 {
			value = defaultValue[0]
		}
		return
	}
	value = configuration.Value
	return
}

func Update(key, value string) error {
	return dao.Mdb.Model(&mdb.Configuration{}).Where("key = ?", key).Updates(map[string]any{"value": value}).Error
}
