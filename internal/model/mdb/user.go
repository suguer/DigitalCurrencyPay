package mdb

import "time"

type User struct {
	ID        uint       `json:"id" gorm:"primaryKey"`
	Username  string     `json:"username" gorm:"type:varchar(192);uniqueIndex;"` // 用户名
	Password  string     `json:"password" gorm:"type:varchar(192);"`             // 密码
	Secret    string     `json:"secret" gorm:"type:varchar(192);"`               // api key
	CreatedAt *time.Time `json:"created_at"`                                     // 创建时间
	UpdatedAt *time.Time `json:"updated_at"`                                     // 更新时间
}
