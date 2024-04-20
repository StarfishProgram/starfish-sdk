package starfish_sdk

import (
	"database/sql/driver"
	"errors"
	"strconv"
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
