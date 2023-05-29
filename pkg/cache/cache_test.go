package cache

import (
	"testing"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

type TestObject struct {
	ID string
}

func TestCache(t *testing.T) {
	bucket := BucketName("testBucket")
	key := "testKey"
	object := TestObject{"1234"}

	err := Set(bucket, key, object)
	assert.ErrorIs(t, err, ErrDisabled)

	err = Get(bucket, key, nil)
	assert.ErrorIs(t, err, ErrDisabled)

	Enable(true)

	err = Set(bucket, key, object)
	require.NoError(t, err)

	var out TestObject
	err = Get(bucket, key, &out)
	require.NoError(t, err)
	assert.Equal(t, object, out)

	err = Get(bucket, key+"5", &out)
	assert.ErrorIs(t, err, ErrEntryNotFound)

	Enable(false)

	err = Set(bucket, key, object)
	assert.ErrorIs(t, err, ErrDisabled)

	err = Get(bucket, key, nil)
	assert.ErrorIs(t, err, ErrDisabled)
}
