package sdktypes

import "time"

type BaseModel struct {
	Id         ID        `gorm:"primarykey;type:bigint;not null;comment:ID;" json:"id"`
	CreateTime time.Time `gorm:"->;type:timestamp;default:current_timestamp;not null;comment:创建时间;" json:"createTime"`
	UpdateTime time.Time `gorm:"->;type:timestamp;default:current_timestamp on update current_timestamp;not null;comment:修改时间;" json:"updateTime"`
}

func NewBaseModel(id ID) BaseModel {
	return BaseModel{
		Id:         id,
		CreateTime: time.Time{},
		UpdateTime: time.Time{},
	}
}
