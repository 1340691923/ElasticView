package hive_dirver

import (
	"context"
	"encoding/json"
	"fmt"
	"reflect"

	"gorm.io/gorm/schema"
)

func init() {
	schema.RegisterSerializer("map", MapSerializer{})
}

// MapSerializer map serializer
type MapSerializer struct{}

// Scan implements serializer interface
func (MapSerializer) Scan(ctx context.Context, field *schema.Field, dst reflect.Value, dbValue interface{}) (err error) {
	fieldValue := reflect.New(field.FieldType)

	if dbValue != nil {
		var bytes []byte
		switch v := dbValue.(type) {
		case []byte:
			bytes = v
		case string:
			bytes = []byte(v)
		default:
			return fmt.Errorf("failed to unmarshal JSONB value: %#v", dbValue)
		}

		if len(bytes) > 0 {
			err = json.Unmarshal(bytes, fieldValue.Interface())
		}
	}

	field.ReflectValueOf(ctx, dst).Set(fieldValue.Elem())
	return
}

// Value implements serializer interface
func (MapSerializer) Value(ctx context.Context, field *schema.Field, dst reflect.Value, fieldValue interface{}) (interface{}, error) {
	// TODO: not implemented
	// MAP('math',120,'english', 123)

	result, err := json.Marshal(fieldValue)
	if string(result) == "null" {
		if field.TagSettings["NOT NULL"] != "" {
			return "", nil
		}
		return nil, err
	}
	return string(result), err
}
