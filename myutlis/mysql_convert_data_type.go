package myutils

import (
	"database/sql"
	"encoding/json"
	"reflect"
)

//INT64
type FieldInt sql.NullInt64

func (f *FieldInt) Scan(value interface{}) error {
	var i sql.NullInt64
	if err := i.Scan(value); err != nil {
		return err
	}

	if reflect.TypeOf(value) == nil {
		*f = FieldInt{
			Int64: 0,
			Valid: false,
		}
	} else {
		*f = FieldInt{
			Int64: i.Int64,
			Valid: true,
		}
	}
	return nil
}

func (f *FieldInt) GetValue() int64 {
	if f.Valid {
		return f.Int64
	} else {
		return 0
	}
}

func (f *FieldInt) GetInt32() int32 {
	return int32(f.GetValue())
}

func (f *FieldInt) GetInt() int {
	return int(f.GetValue())
}

func (f *FieldInt) MarshalJSON() ([]byte, error) {
	return json.Marshal(f.GetValue())
}

func (f *FieldInt) UnmarshalJSON(b []byte) error {
	err := json.Unmarshal(b, &f.Int64)
	f.Valid = err == nil
	return err
}

// STRING
type FieldString sql.NullString

func (f *FieldString) Scan(value interface{}) error {
	var s sql.NullString
	if err := s.Scan(value); err != nil {
		return err
	}

	if reflect.TypeOf(value) == nil {
		*f = FieldString{
			String: "",
			Valid:  false,
		}
	} else {
		*f = FieldString{
			String: s.String,
			Valid:  true,
		}
	}
	return nil
}

func (f *FieldString) GetValue() string {
	if f.Valid {
		return f.String
	} else {
		return ""
	}
}

func (f *FieldString) MarshalJSON() ([]byte, error) {
	return json.Marshal(f.GetValue())
}

func (f *FieldString) UnmarshalJSON(b []byte) error {
	err := json.Unmarshal(b, &f.String)
	f.Valid = err == nil
	return err
}

// BOOL
type FieldBoolean sql.NullBool

func (f *FieldBoolean) Scan(value interface{}) error {
	var b sql.NullBool
	if err := b.Scan(value); err != nil {
		return err
	}

	if reflect.TypeOf(value) == nil {
		*f = FieldBoolean{
			Bool:  false,
			Valid: false,
		}
	} else {
		*f = FieldBoolean{
			Bool:  b.Bool,
			Valid: true,
		}
	}
	return nil
}

func (f *FieldBoolean) GetValue() bool {
	if f.Valid {
		return f.Bool
	} else {
		return false
	}
}

func (f *FieldBoolean) MarshalJSON() ([]byte, error) {
	return json.Marshal(f.GetValue())
}

func (f *FieldBoolean) UnmarshalJSON(b []byte) error {
	err := json.Unmarshal(b, &f.Bool)
	f.Valid = err == nil
	return err
}

//FLOAT64
type FieldFloat sql.NullFloat64

func (f *FieldFloat) Scan(value interface{}) error {
	var fl sql.NullFloat64

	if err := fl.Scan(value); err != nil {
		return err
	}

	if reflect.TypeOf(value) == nil {
		*f = FieldFloat{
			Float64: 0,
			Valid:   false,
		}
	} else {
		*f = FieldFloat{
			Float64: fl.Float64,
			Valid:   true,
		}
	}

	return nil
}

func (f *FieldFloat) GetValue() float64 {
	if f.Valid {
		return f.Float64
	} else {
		return 0
	}
}

func (f *FieldFloat) GetFloat() float32 {
	return float32(f.GetValue())
}

func (f *FieldFloat) MarshalJSON() ([]byte, error) {
	return json.Marshal(f.GetValue())
}

func (f *FieldFloat) UnmarshalJSON(b []byte) error {
	err := json.Unmarshal(b, &f.Float64)
	f.Valid = err == nil
	return err
}
