package starfish_sdk

import (
	"database/sql/driver"
	"errors"
	"strconv"
	"time"
)

type Result[D any] struct {
	Code ICode
	Data D
}

func (r *Result[D]) IsOk() bool {
	return r.Code == nil
}

func (r *Result[D]) IsFaild() bool {

	return !r.IsOk()
}

type ID int64

func (v ID) Value() (driver.Value, error) {
	return int64(v), nil
}

func (v *ID) Scan(src interface{}) error {
	if src == nil {
		return nil
	}
	sv, ok := src.(int64)
	if !ok {
		return errors.New("类型转换错误")
	}
	*v = ID(sv)
	return nil
}

func (v ID) MarshalJSON() ([]byte, error) {
	return []byte(strconv.FormatInt(int64(v), 10)), nil
}

func (v *ID) UnmarshalJSON(src []byte) error {
	d, err := strconv.ParseInt(string(src), 10, 64)
	if err != nil {
		return err
	}
	*v = ID(d)
	return nil
}

type BaseModel struct {
	ID         ID         `gorm:"primarykey;type:bigint;not null;comment:ID;" json:"id"`
	CreateTime time.Time  `gorm:"->;type:timestamp;default:current_timestamp;not null;comment:创建时间;" json:"create_time"`
	UpdateTime time.Time  `gorm:"->;type:timestamp;default:current_timestamp on update current_timestamp;not null;comment:修改时间;" json:"update_time"`
	DeleteTime *time.Time `gorm:"->;type:timestamp;index;comment:删除时间;" json:"-"`
}

func NewBaseModel() BaseModel {
	return BaseModel{
		ID:         UUID().ID(),
		CreateTime: time.Time{},
		UpdateTime: time.Time{},
		DeleteTime: nil,
	}
}
