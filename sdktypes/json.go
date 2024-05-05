package sdktypes

import (
	"database/sql/driver"
	"encoding/json"
	"errors"
	"fmt"
)

type JSON struct {
	Raw json.RawMessage
}

func (v JSON) Value() (driver.Value, error) {
	if len(v.Raw) == 0 {
		return nil, nil
	}
	return v.Raw.MarshalJSON()
}

func (v *JSON) Scan(value interface{}) error {
	bytes, ok := value.([]byte)
	if !ok {
		return errors.New(fmt.Sprint("Failed to unmarshal JSONB value:", value))
	}
	err := json.Unmarshal(bytes, &v.Raw)
	return err
}

func (v JSON) MarshalJSON() ([]byte, error) {
	return v.Raw.MarshalJSON()
}

func (v *JSON) UnmarshalJSON(data []byte) error {
	return v.Raw.UnmarshalJSON(data)
}

func (v *JSON) Parse(data any) error {
	bytes, err := json.Marshal(data)
	if err != nil {
		return err
	}
	v.UnmarshalJSON(bytes)
	return nil
}

func (v *JSON) To(data any) error {
	return json.Unmarshal(v.Raw, data)
}
