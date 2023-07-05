package collections

import (
	"testing"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

func TestAddMap(t *testing.T) {
	to := map[int]string{
		0: "a",
		2: "c",
	}

	from := map[int]string{
		1: "b",
		3: "d",
	}

	err := AddMap(to, from)
	require.NoError(t, err)

	assert.Contains(t, to, 1)
	assert.Equal(t, "b", to[1])

	assert.Contains(t, to, 3)
	assert.Equal(t, "d", to[3])
}

func TestAddMap_duplicate_key(t *testing.T) {
	to := map[int]string{
		0: "a",
		2: "c",
	}

	from := map[int]string{
		0: "b",
	}

	err := AddMap(to, from)
	assert.ErrorIs(t, err, ErrDuplicateKey)
}
