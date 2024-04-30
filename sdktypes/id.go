package sdktypes

import (
	"database/sql/driver"
	"errors"
	"fmt"
	"strconv"
)

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
	return []byte("\"" + strconv.FormatInt(int64(v), 10) + "\""), nil
}

func (v *ID) UnmarshalJSON(src []byte) error {
	s, err := unquoteIfQuoted(src)
	if err != nil {
		return err
	}
	d, err := strconv.ParseInt(s, 10, 64)
	if err != nil {
		return err
	}
	*v = ID(d)
	return nil
}

func (v ID) String() string {
	return strconv.FormatInt(int64(v), 10)
}

func unquoteIfQuoted(value interface{}) (string, error) {
	var bytes []byte
	switch v := value.(type) {
	case string:
		bytes = []byte(v)
	case []byte:
		bytes = v
	default:
		return "", fmt.Errorf("could not convert value '%+v' to byte array of type '%T'",
			value, value)
	}
	if len(bytes) > 2 && bytes[0] == '"' && bytes[len(bytes)-1] == '"' {
		bytes = bytes[1 : len(bytes)-1]
	}
	return string(bytes), nil
}
