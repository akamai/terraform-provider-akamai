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

	err = Delete(bucket, key)
	require.NoError(t, err)

	err = Get(bucket, key, &out)
	assert.ErrorIs(t, err, ErrEntryNotFound, "deleted entry should return ErrEntryNotFound")

	err = Set(bucket, key, object)
	require.NoError(t, err, "Set after Delete should succeed")

	err = Get(bucket, key, &out)
	require.NoError(t, err, "Get after re-Set should succeed")
	assert.Equal(t, object, out)

	Enable(false)

	err = Delete(bucket, key)
	assert.ErrorIs(t, err, ErrDisabled)

	err = Set(bucket, key, object)
	assert.ErrorIs(t, err, ErrDisabled)

	err = Get(bucket, key, nil)
	assert.ErrorIs(t, err, ErrDisabled)
}
