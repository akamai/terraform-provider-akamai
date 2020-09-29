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
