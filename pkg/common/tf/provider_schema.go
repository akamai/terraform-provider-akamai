package tf

import (
	"errors"
	"fmt"
	"strings"

	"github.com/hashicorp/go-cty/cty"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/schema"
)

var (
	// ErrInvalidType is returned when given value is of invalid type (i.e. different type than requested)
	ErrInvalidType = errors.New("value must be of the specified type")
	// ErrNotFound is returned when value is not present on schema
	ErrNotFound = errors.New("value not found")
	// ErrValueSet is returned when setting property value returned an error
	ErrValueSet = errors.New("setting property value")
	// ErrEmptyKey is returned when setting property value returned an error
	ErrEmptyKey = errors.New("provided key cannot be empty")
	// ErrEmptyPath is returned when path is empty
	ErrEmptyPath = errors.New("path cannot be empty")
)

// ResourceDataFetcher allows getting values from resource data.
type ResourceDataFetcher interface {
	GetOk(string) (any, bool)
}

// ResourceChangeFetcher allows getting changes to the resource data.
type ResourceChangeFetcher interface {
	GetChange(string) (any, any)
	HasChange(string) bool
}

// GetSchemaFieldNameFromPath returns schema field name from given path
//
// if len of path is zero it returns empty string and error
func GetSchemaFieldNameFromPath(path cty.Path) (string, error) {
	if len(path) == 0 {
		return "", ErrEmptyPath
	}

	attrStep, ok := path[len(path)-1].(cty.GetAttrStep)
	if !ok {
		return "", ErrInvalidType
	}

	return attrStep.Name, nil
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
	if value != nil && ok {
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
	if value != nil && ok {
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
	if value == nil || !ok {
		return 0, fmt.Errorf("%w: %s", ErrNotFound, key)
	}
	var num int
	if num, ok = value.(int); !ok {
		return 0, fmt.Errorf("%w: %s, %q", ErrInvalidType, key, "int")
	}
	return num, nil
}

// GetIntValueAsInt64 fetches value with given key from ResourceData object and attempts type cast to int, if succeed, it returns value as int64
//
// if value is not present on provided resource, ErrNotFound is returned
// if casting is not successful, ErrInvalidType is returned
func GetIntValueAsInt64(key string, rd ResourceDataFetcher) (int64, error) {
	num, err := GetIntValue(key, rd)
	if err != nil {
		return 0, err
	}
	return int64(num), nil
}

// GetInt64Value fetches value with given key from ResourceData object and attempts type cast to int64
//
// if value is not present on provided resource, ErrNotFound is returned
// if casting is not successful, ErrInvalidType is returned
func GetInt64Value(key string, rd ResourceDataFetcher) (int64, error) {
	if key == "" {
		return 0, fmt.Errorf("%w: %s", ErrEmptyKey, key)
	}
	value, ok := rd.GetOk(key)
	if value == nil || !ok {
		return 0, fmt.Errorf("%w: %s", ErrNotFound, key)
	}
	var num int64
	if num, ok = value.(int64); !ok {
		return 0, fmt.Errorf("%w: %s, %q", ErrInvalidType, key, "int64")
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
	if value == nil || !ok {
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
	if value == nil || !ok {
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
	if value == nil {
		return false, fmt.Errorf("%w: %s", ErrNotFound, key)
	}

	val, ok := value.(bool)
	if !ok {
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
	if value == nil || !ok {
		return val, fmt.Errorf("%w: %s", ErrNotFound, key)
	}
	if val, ok = value.(*schema.Set); !ok {
		return val, fmt.Errorf("%w: %s, %q", ErrInvalidType, key, "*schema.Set")
	}
	return val, nil
}

// GetRawSetValue fetches value with given key from RawConfig with fallback to ResourceData object and returns it as []any.
//
// It attempts to work around an issue that Set with diff suppress, returns in some cases only diff (rather than whole value).
func GetRawSetValue(key string, rd ResourceDataFetcher, rawConfig *RawConfig) ([]any, error) {
	setFunctionProvider, directErr := GetSetValue(key, rd)
	listValue, rawErr := GetListValue(key, rawConfig)

	if directErr == nil && rawErr == nil {
		// All data is available - using value from 'raw' and order from 'direct'
		return schema.NewSet(setFunctionProvider.F, listValue).List(), nil
	}

	if directErr != nil && rawErr == nil {
		// Only 'raw' is available - values should be correct, order may be slightly changed, but for set it shouldn't matter
		return listValue, nil
	}

	if directErr == nil {
		// Only 'direct' is available so falling back to it.
		// In some cases 'direct' can have incorrect values so we cannot always use it. No such use case was identified that 'raw' was not available and 'direct' hold incorrect values.
		return setFunctionProvider.List(), nil
	}

	// Both errors - unable to provide values
	return nil, fmt.Errorf("was not able to read %s value, %w: %s", key, directErr, rawErr)
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
	if value == nil || !ok {
		return val, fmt.Errorf("%w: %s", ErrNotFound, key)
	}
	if val, ok = value.([]interface{}); !ok {
		return nil, fmt.Errorf("%w: %s, %q", ErrInvalidType, key, "[]interface{}")
	}
	return val, nil
}

// GetTypedListValue fetches value with given key from ResourceData object and attempts to create []T
//
// if value is not present on provided resource, ErrNotFound is returned
// if casting is not successful, ErrInvalidType is returned
func GetTypedListValue[T any](key string, rd ResourceDataFetcher) ([]T, error) {
	values, err := GetListValue(key, rd)
	if err != nil {
		return []T{}, err
	}

	out := make([]T, 0, len(values))
	for _, val := range values {
		v, ok := val.(T)
		if !ok {
			var e T
			return []T{}, fmt.Errorf("%w: %s, %q", ErrInvalidType, key, fmt.Sprintf("%T", e))
		}
		out = append(out, v)
	}

	return out, nil
}

// GetMapValue fetches value with given key from ResourceData object and attempts type cast to map[string]interface{}
//
// if value is not present on provided resource, ErrNotFound is returned
// if casting is not successful, ErrInvalidType is returned
func GetMapValue(key string, rd ResourceDataFetcher) (map[string]interface{}, error) {
	if key == "" {
		return nil, fmt.Errorf("%w: %s", ErrEmptyKey, key)
	}
	val := make(map[string]interface{}, 0)
	value, ok := rd.GetOk(key)
	if value == nil || !ok {
		return val, fmt.Errorf("%w: %s", ErrNotFound, key)
	}
	if val, ok = value.(map[string]interface{}); !ok {
		return nil, fmt.Errorf("%w: %s, %q", ErrInvalidType, key, "map[string]interface{}")
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
		if value != nil && ok {
			str, ok := value.(string)
			if !ok {
				continue
			}

			rval = append(rval, str)
		}
	}

	return rval
}

// ResolveKeyStringState fetches value with given key (or fallbackKey) from ResourceData object
func ResolveKeyStringState(rd ResourceDataFetcher, key, fallbackKey string) (value string, err error) {
	value, err = GetStringValue(key, rd)
	if errors.Is(err, ErrNotFound) {
		value, err = GetStringValue(fallbackKey, rd)
	}
	if err != nil {
		return "", err
	}
	return value, nil
}

// StateNetwork changes the value of the input before storing it in state
func StateNetwork(i interface{}) string {
	val, ok := i.(string)
	if !ok {
		panic(fmt.Sprintf("value type is not a string: %T", i))
	}

	switch strings.ToLower(val) {
	case "production", "prod", "p":
		return "production"
	case "staging", "stag", "s":
		return "staging"
	}

	// this should never happen :-)
	return val
}

// RestoreOldValues reverts the value in schema of the given keys
func RestoreOldValues(rd *schema.ResourceData, keys []string) error {
	for _, key := range keys {
		if rd.HasChange(key) {
			oldVersion, _ := rd.GetChange(key)
			if err := rd.Set(key, oldVersion); err != nil {
				return err
			}
		}
	}
	return nil
}

// GetExactlyOneOf extracts exactly one value with given keys from ResourceData object
// if multiple values are present the function returns first one found
func GetExactlyOneOf(rd ResourceDataFetcher, keys []string) (foundKey string, value interface{}, err error) {
	for _, key := range keys {
		value, err := GetSetValue(key, rd)
		if errors.Is(err, ErrNotFound) {
			continue
		}

		if err != nil {
			return "", nil, err
		}

		return key, value, nil
	}
	return "", nil, ErrNotFound
}

// GetSetAsListValue retrieves a schema.TypeSet from Terraform resource data and returns its elements as a slice of interface{}.
// It performs key validation and type assertion, returning detailed errors for invalid or missing data.
func GetSetAsListValue(key string, rd ResourceDataFetcher) ([]interface{}, error) {
	if key == "" {
		return nil, fmt.Errorf("%w: %s", ErrEmptyKey, key)
	}
	value, ok := rd.GetOk(key)
	if value == nil || !ok {
		return nil, fmt.Errorf("%w: %s", ErrNotFound, key)
	}

	set, ok := value.(*schema.Set)
	if !ok {
		return nil, fmt.Errorf("%w: %s, expected *schema.Set", ErrInvalidType, key)
	}

	return set.List(), nil
}
