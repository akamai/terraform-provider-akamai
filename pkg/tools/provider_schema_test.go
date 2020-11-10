package tools

import (
	"errors"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/schema"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/mock"
	"github.com/stretchr/testify/require"
	"testing"
)

type mocked struct {
	mock.Mock
}

func (m *mocked) GetOk(key string) (interface{}, bool) {
	args := m.Called(key)
	return args.Get(0), args.Bool(1)
}

func TestGetStringValue(t *testing.T) {
	tests := map[string]struct {
		key       string
		init      func(*mocked)
		expected  string
		withError error
	}{
		"string value found": {
			key: "key",
			init: func(m *mocked) {
				m.On("GetOk", "key").Return("value", true).Once()
			},
			expected: "value",
		},
		"string value not found": {
			key: "key",
			init: func(m *mocked) {
				m.On("GetOk", "key").Return("", false).Once()
			},
			withError: ErrNotFound,
		},
		"empty key passed": {
			key:       "",
			init:      func(m *mocked) {},
			withError: ErrEmptyKey,
		},
		"value is of invalid type": {
			key: "key",
			init: func(m *mocked) {
				m.On("GetOk", "key").Return(1, true).Once()
			},
			withError: ErrInvalidType,
		},
	}
	for name, test := range tests {
		t.Run(name, func(t *testing.T) {
			m := &mocked{}
			test.init(m)
			res, err := GetStringValue(test.key, m)
			m.AssertExpectations(t)
			if test.withError != nil {
				assert.True(t, errors.Is(err, test.withError), "want: %s; got: %s", test.withError, err)
				return
			}
			require.NoError(t, err)
			assert.Equal(t, test.expected, res)
		})
	}
}

func TestGetIntValue(t *testing.T) {
	tests := map[string]struct {
		key       string
		init      func(*mocked)
		expected  int
		withError error
	}{
		"int value found": {
			key: "key",
			init: func(m *mocked) {
				m.On("GetOk", "key").Return(1, true).Once()
			},
			expected: 1,
		},
		"int value not found": {
			key: "key",
			init: func(m *mocked) {
				m.On("GetOk", "key").Return(0, false).Once()
			},
			withError: ErrNotFound,
		},
		"empty key passed": {
			key:       "",
			init:      func(m *mocked) {},
			withError: ErrEmptyKey,
		},
		"value is of invalid type": {
			key: "key",
			init: func(m *mocked) {
				m.On("GetOk", "key").Return("value", true).Once()
			},
			withError: ErrInvalidType,
		},
	}
	for name, test := range tests {
		t.Run(name, func(t *testing.T) {
			m := &mocked{}
			test.init(m)
			res, err := GetIntValue(test.key, m)
			m.AssertExpectations(t)
			if test.withError != nil {
				assert.True(t, errors.Is(err, test.withError), "want: %s; got: %s", test.withError, err)
				return
			}
			assert.Equal(t, test.expected, res)
		})
	}
}

func TestGetBoolValue(t *testing.T) {
	tests := map[string]struct {
		key       string
		init      func(*mocked)
		expected  bool
		withError error
	}{
		"bool value found": {
			key: "key",
			init: func(m *mocked) {
				m.On("GetOk", "key").Return(true, true).Once()
			},
			expected: true,
		},
		"bool value not found": {
			key: "key",
			init: func(m *mocked) {
				m.On("GetOk", "key").Return(false, false).Once()
			},
			expected: false,
		},
		"empty key passed": {
			key:       "",
			init:      func(m *mocked) {},
			withError: ErrEmptyKey,
		},
		"value is of invalid type": {
			key: "key",
			init: func(m *mocked) {
				m.On("GetOk", "key").Return(1, true).Once()
			},
			withError: ErrInvalidType,
		},
	}
	for name, test := range tests {
		t.Run(name, func(t *testing.T) {
			m := &mocked{}
			test.init(m)
			res, err := GetBoolValue(test.key, m)
			m.AssertExpectations(t)
			if test.withError != nil {
				assert.True(t, errors.Is(err, test.withError), "want: %s; got: %s", test.withError, err)
				return
			}
			assert.Equal(t, test.expected, res)
		})
	}
}

func TestGetSetValue(t *testing.T) {
	tests := map[string]struct {
		key       string
		init      func(*mocked)
		expected  *schema.Set
		withError error
	}{
		"set value found": {
			key: "key",
			init: func(m *mocked) {
				m.On("GetOk", "key").Return(&schema.Set{}, true).Once()
			},
			expected: &schema.Set{},
		},
		"string value not found": {
			key: "key",
			init: func(m *mocked) {
				m.On("GetOk", "key").Return(nil, false).Once()
			},
			withError: ErrNotFound,
		},
		"empty key passed": {
			key:       "",
			init:      func(m *mocked) {},
			withError: ErrEmptyKey,
		},
		"value is of invalid type": {
			key: "key",
			init: func(m *mocked) {
				m.On("GetOk", "key").Return(1, true).Once()
			},
			withError: ErrInvalidType,
		},
	}
	for name, test := range tests {
		t.Run(name, func(t *testing.T) {
			m := &mocked{}
			test.init(m)
			res, err := GetSetValue(test.key, m)
			m.AssertExpectations(t)
			if test.withError != nil {
				assert.True(t, errors.Is(err, test.withError), "want: %s; got: %s", test.withError, err)
				return
			}
			assert.Equal(t, test.expected, res)
		})
	}
}

func TestGetInterfaceArrayValue(t *testing.T) {
	tests := map[string]struct {
		key       string
		init      func(*mocked)
		expected  []interface{}
		withError error
	}{
		"[]interface{} value found": {
			key: "key",
			init: func(m *mocked) {
				m.On("GetOk", "key").Return(make([]interface{}, 1), true).Once()
			},
			expected: make([]interface{}, 1),
		},
		"[]interface{} value not found": {
			key: "key",
			init: func(m *mocked) {
				m.On("GetOk", "key").Return(nil, false).Once()
			},
			withError: ErrNotFound,
		},
		"empty key passed": {
			key:       "",
			init:      func(m *mocked) {},
			withError: ErrEmptyKey,
		},
		"value is of invalid type": {
			key: "key",
			init: func(m *mocked) {
				m.On("GetOk", "key").Return(1, true).Once()
			},
			withError: ErrInvalidType,
		},
	}

	for name, test := range tests {
		t.Run(name, func(t *testing.T) {
			m := &mocked{}
			test.init(m)
			res, err := GetInterfaceArrayValue(test.key, m)
			m.AssertExpectations(t)
			if test.withError != nil {
				assert.True(t, errors.Is(err, test.withError), "want: %s; got: %s", test.withError, err)
				return
			}
			assert.Equal(t, test.expected, res)
		})
	}
}

func TestGetFloat64Value(t *testing.T) {
	tests := map[string]struct {
		key       string
		init      func(*mocked)
		expected  float64
		withError error
	}{
		"float64 value found": {
			key: "key",
			init: func(m *mocked) {
				m.On("GetOk", "key").Return(float64(1), true).Once()
			},
			expected: 1,
		},
		"float64 value not found": {
			key: "key",
			init: func(m *mocked) {
				m.On("GetOk", "key").Return(float64(0), false).Once()
			},
			withError: ErrNotFound,
		},
		"empty key passed": {
			key:       "",
			init:      func(m *mocked) {},
			withError: ErrEmptyKey,
		},
		"value is of invalid type": {
			key: "key",
			init: func(m *mocked) {
				m.On("GetOk", "key").Return("not a float64", true).Once()
			},
			withError: ErrInvalidType,
		},
	}

	for name, test := range tests {
		t.Run(name, func(t *testing.T) {
			m := &mocked{}
			test.init(m)
			res, err := GetFloat64Value(test.key, m)
			m.AssertExpectations(t)
			if test.withError != nil {
				assert.True(t, errors.Is(err, test.withError), "want: %s; got: %s", test.withError, err)
				return
			}
			assert.Equal(t, test.expected, res)
		})
	}
}

func TestGetFloat32Value(t *testing.T) {
	tests := map[string]struct {
		key       string
		init      func(*mocked)
		expected  float32
		withError error
	}{
		"float32 value found": {
			key: "key",
			init: func(m *mocked) {
				m.On("GetOk", "key").Return(float32(1), true).Once()
			},
			expected: 1,
		},
		"float64 value not found": {
			key: "key",
			init: func(m *mocked) {
				m.On("GetOk", "key").Return(float32(0), false).Once()
			},
			withError: ErrNotFound,
		},
		"empty key passed": {
			key:       "",
			init:      func(m *mocked) {},
			withError: ErrEmptyKey,
		},
		"value is of invalid type": {
			key: "key",
			init: func(m *mocked) {
				m.On("GetOk", "key").Return("not a float32", true).Once()
			},
			withError: ErrInvalidType,
		},
	}

	for name, test := range tests {
		t.Run(name, func(t *testing.T) {
			m := &mocked{}
			test.init(m)
			res, err := GetFloat32Value(test.key, m)
			m.AssertExpectations(t)
			if test.withError != nil {
				assert.True(t, errors.Is(err, test.withError), "want: %s; got: %s", test.withError, err)
				return
			}
			assert.Equal(t, test.expected, res)
		})
	}
}

func TestGetListValue(t *testing.T) {
	tests := map[string]struct {
		key       string
		init      func(*mocked)
		expected  []interface{}
		withError error
	}{
		"set value found": {
			key: "key",
			init: func(m *mocked) {
				m.On("GetOk", "key").Return(make([]interface{}, 1), true).Once()
			},
			expected: make([]interface{}, 1),
		},
		"string value not found": {
			key: "key",
			init: func(m *mocked) {
				m.On("GetOk", "key").Return(nil, false).Once()
			},
			withError: ErrNotFound,
		},
		"empty key passed": {
			key:       "",
			init:      func(m *mocked) {},
			withError: ErrEmptyKey,
		},
		"value is of invalid type": {
			key: "key",
			init: func(m *mocked) {
				m.On("GetOk", "key").Return(1, true).Once()
			},
			withError: ErrInvalidType,
		},
	}
	for name, test := range tests {
		t.Run(name, func(t *testing.T) {
			m := &mocked{}
			test.init(m)
			res, err := GetListValue(test.key, m)
			m.AssertExpectations(t)
			if test.withError != nil {
				assert.True(t, errors.Is(err, test.withError), "want: %s; got: %s", test.withError, err)
				return
			}
			assert.Equal(t, test.expected, res)
		})
	}
}

func TestFindStringValues(t *testing.T) {
	tests := map[string]struct {
		key      string
		init     func(*mocked)
		expected []string
	}{
		"set value found": {
			key: "key",
			init: func(m *mocked) {
				m.On("GetOk", "key").Return("found", true).Once()
			},
			expected: []string{"found"},
		},
		"string value not found": {
			key: "key",
			init: func(m *mocked) {
				m.On("GetOk", "key").Return(nil, false).Once()
			},
			expected: []string{},
		},
		"value is of invalid type": {
			key: "key",
			init: func(m *mocked) {
				m.On("GetOk", "key").Return(1, true).Once()
			},
			expected: []string{},
		},
	}
	for name, test := range tests {
		t.Run(name, func(t *testing.T) {
			m := &mocked{}
			test.init(m)
			res := FindStringValues(m, test.key)
			m.AssertExpectations(t)
			assert.Equal(t, test.expected, res)
		})
	}
}

func TestResolveKeyState(t *testing.T) {
	const key = "key"
	const fallbackKey = "keyId"
	const value = "value"

	tests := map[string]struct {
		key       string
		fallback  string
		init      func(*mocked)
		expected  string
		withError error
	}{
		"key value found": {
			init: func(m *mocked) {
				m.On("GetOk", key).Return(value, true).Once()
			},
			expected: value,
		},
		"key not found, fall back value found": {
			init: func(m *mocked) {
				m.On("GetOk", key).Return(nil, false).Once()
				m.On("GetOk", fallbackKey).Return(value, true).Once()
			},
			expected: value,
		},
		"key not found, fall back not found": {
			init: func(m *mocked) {
				m.On("GetOk", mock.Anything).Return(nil, false)
			},
			expected:  "",
			withError: ErrNotFound,
		},
		"value type not string": {
			key: "other type",
			init: func(m *mocked) {
				m.On("GetOk", mock.Anything).Return(make([]string, 0), true).Once()
			},
			expected:  "",
			withError: ErrInvalidType,
		},
	}

	for name, test := range tests {
		t.Run(name, func(t *testing.T) {
			m := &mocked{}
			test.init(m)
			res, err := ResolveKeyStringState(key, fallbackKey, m)
			m.AssertExpectations(t)
			if test.withError != nil {
				assert.True(t, errors.Is(err, test.withError), "want: %s; got: %s", test.withError, err)
				return
			}
			assert.Equal(t, test.expected, res)
		})
	}
}
