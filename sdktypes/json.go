package sdktypes

import (
	"database/sql/driver"
	"encoding/json"
	"errors"
	"fmt"
)

type JSON json.RawMessage

func (v JSON) Value() (driver.Value, error) {
	if len(v) == 0 {
		return nil, nil
	}
	return json.RawMessage(v).MarshalJSON()
}

func (v *JSON) Scan(value interface{}) error {
	bytes, ok := value.([]byte)
	if !ok {
		return errors.New(fmt.Sprint("Failed to unmarshal JSONB value:", value))
	}
	result := json.RawMessage{}
	err := json.Unmarshal(bytes, &result)
	*v = JSON(result)
	return err
}

func Marshal(data any) (JSON, error) {
	j, err := json.Marshal(data)
	if err != nil {
		return nil, err
	}
	return JSON(j), nil
}

func Unmarshal(j JSON, data any) error {
	return json.Unmarshal(j, data)
}
