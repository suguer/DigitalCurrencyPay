package mdb

type Configuration struct {
	ID     uint   `json:"id" gorm:"primaryKey"`
	Key    string `json:"key" gorm:"type:varchar(128);uniqueIndex:keys;"` // 配置键
	Value  string `json:"value" gorm:"type:varchar(256);"`                // 配置值
	Status int    `json:"status"`                                         // 交易状态（0-待确认 1-已确认 2-失败）
	Remark string `json:"remark" gorm:"type:varchar(256);"`               // 备注
}
