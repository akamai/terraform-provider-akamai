package tools

import (
	"errors"
	"fmt"

	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/schema"
)

var (
	// ErrInvalidType is returned when given value is of invalid type (i.e. different type than requested)
	ErrInvalidType = errors.New("value must be of the specified type")
	// ErrNotFound is returned when value is not present on schema
	ErrNotFound = errors.New("value not found")
	// ErrValueSet is returned when setting property value returned an error
	ErrValueSet = errors.New("setting property value")
	// ErrValueSet is returned when setting property value returned an error
	ErrEmptyKey = errors.New("provided key cannot be empty")
)

type ResourceDataFetcher interface {
	GetOk(string) (interface{}, bool)
}

// GetStringValue fetches value with given key from ResourceData object and attempts type cast to string
//
// if value is not present on provided resource for key, ErrNotFound is returned
// if casting is not successful, ErrInvalidType is returned
func GetStringValue(key string, rd ResourceDataFetcher) (string, error) {
	if key == "" {
		return "", fmt.Errorf("%w: %s", ErrEmptyKey, key)
	}

	value, ok := rd.GetOk(key)
	if ok {
		str, ok := value.(string)
		if !ok {
			return "", fmt.Errorf("%w: %s, %q", ErrInvalidType, key, "string")
		}

		return str, nil
	}

	return "", fmt.Errorf("%w: %s", ErrNotFound, key)
}

// GetInterfaceArrayValue fetches value with given key from ResourceData object and attempts type cast to []interface{}
//
// if value is not present on provided resource for key, ErrNotFound is returned
// if casting is not successful, ErrInvalidType is returned
func GetInterfaceArrayValue(key string, rd ResourceDataFetcher) ([]interface{}, error) {
	if key == "" {
		return nil, fmt.Errorf("%w: %s", ErrEmptyKey, key)
	}

	value, ok := rd.GetOk(key)
	if ok {
		interf, ok := value.([]interface{})
		if !ok {
			return nil, fmt.Errorf("%w: %s, %q", ErrInvalidType, key, "[]interface{}")
		}

		return interf, nil
	}

	return nil, fmt.Errorf("%w: %s", ErrNotFound, key)
}

// GetIntValue fetches value with given key from ResourceData object and attempts type cast to int
//
// if value is not present on provided resource, ErrNotFound is returned
// if casting is not successful, ErrInvalidType is returned
func GetIntValue(key string, rd ResourceDataFetcher) (int, error) {
	if key == "" {
		return 0, fmt.Errorf("%w: %s", ErrEmptyKey, key)
	}
	value, ok := rd.GetOk(key)
	if !ok {
		return 0, fmt.Errorf("%w: %s", ErrNotFound, key)
	}
	var num int
	if num, ok = value.(int); !ok {
		return 0, fmt.Errorf("%w: %s, %q", ErrInvalidType, key, "int")
	}
	return num, nil
}

// GetFloat64Value fetches value with given key from ResourceData object and attempts type cast to float64
//
// if value is not present on provided resource, ErrNotFound is returned
// if casting is not successful, ErrInvalidType is returned
func GetFloat64Value(key string, rd ResourceDataFetcher) (float64, error) {
	if key == "" {
		return 0, fmt.Errorf("%w: %s", ErrEmptyKey, key)
	}
	value, ok := rd.GetOk(key)
	if !ok {
		return 0, fmt.Errorf("%w: %s", ErrNotFound, key)
	}
	var num float64
	if num, ok = value.(float64); !ok {
		return 0, fmt.Errorf("%w: %s, %q", ErrInvalidType, key, "float64")
	}
	return num, nil
}

// GetFloat32Value fetches value with given key from ResourceData object and attempts type cast to float64
//
// if value is not present on provided resource, ErrNotFound is returned
// if casting is not successful, ErrInvalidType is returned
func GetFloat32Value(key string, rd ResourceDataFetcher) (float32, error) {
	if key == "" {
		return 0, fmt.Errorf("%w: %s", ErrEmptyKey, key)
	}
	value, ok := rd.GetOk(key)
	if !ok {
		return 0, fmt.Errorf("%w: %s", ErrNotFound, key)
	}
	var num float32
	if num, ok = value.(float32); !ok {
		return 0, fmt.Errorf("%w: %s, %q", ErrInvalidType, key, "float32")
	}
	return num, nil
}

// GetBoolValue fetches value with given key from ResourceData object and attempts type cast to bool
//
// if value is not present on provided resource, ErrNotFound is returned
// if casting is not successful, ErrInvalidType is returned
func GetBoolValue(key string, rd ResourceDataFetcher) (bool, error) {
	if key == "" {
		return false, fmt.Errorf("%w: %s", ErrEmptyKey, key)
	}
	value, _ := rd.GetOk(key)
	var val bool
	var ok bool
	if val, ok = value.(bool); !ok {
		return false, fmt.Errorf("%w: %s, %q", ErrInvalidType, key, "bool")
	}
	return val, nil
}

// GetSetValue fetches value with given key from ResourceData object and attempts type cast to *schema.Set
//
// if value is not present on provided resource, ErrNotFound is returned
// if casting is not successful, ErrInvalidType is returned
func GetSetValue(key string, rd ResourceDataFetcher) (*schema.Set, error) {
	if key == "" {
		return nil, fmt.Errorf("%w: %s", ErrEmptyKey, key)
	}
	val := new(schema.Set)
	value, ok := rd.GetOk(key)
	if !ok {
		return val, fmt.Errorf("%w: %s", ErrNotFound, key)
	}
	if val, ok = value.(*schema.Set); !ok {
		return val, fmt.Errorf("%w: %s, %q", ErrInvalidType, key, "*schema.Set")
	}
	return val, nil
}

// GetListValue fetches value with given key from ResourceData object and attempts type cast to []interface{}
//
// if value is not present on provided resource, ErrNotFound is returned
// if casting is not successful, ErrInvalidType is returned
func GetListValue(key string, rd ResourceDataFetcher) ([]interface{}, error) {
	if key == "" {
		return nil, fmt.Errorf("%w: %s", ErrEmptyKey, key)
	}
	value, ok := rd.GetOk(key)
	val := make([]interface{}, 0)
	if !ok {
		return val, fmt.Errorf("%w: %s", ErrNotFound, key)
	}
	if val, ok = value.([]interface{}); !ok {
		return nil, fmt.Errorf("%w: %s, %q", ErrInvalidType, key, "[]interface{}")
	}
	return val, nil
}

// FindStringValues searches the ResourceData for the list of keys and returns the array of values
//
// if the value does not exist it is skipped
// if the value cannot be cast to string it is skipped
func FindStringValues(rd ResourceDataFetcher, keys ...string) []string {
	rval := make([]string, 0)

	for _, key := range keys {
		value, ok := rd.GetOk(key)
		if ok {
			str, ok := value.(string)
			if !ok {
				continue
			}

			rval = append(rval, str)
		}
	}

	return rval
}
