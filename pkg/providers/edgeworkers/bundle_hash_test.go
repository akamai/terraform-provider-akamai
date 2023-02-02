package edgeworkers

import (
	"archive/tar"
	"bytes"
	"compress/gzip"
	"io"
	"testing"

	"github.com/akamai/AkamaiOPEN-edgegrid-golang/v4/pkg/edgeworkers"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

func TestGetSHA256FromBundle(t *testing.T) {
	tests := map[string]struct {
		firstBundlePath  string
		secondBundlePath string
		expectDiff       bool
	}{
		"no diff in bundles first test": {
			firstBundlePath:  bundlePathForCreate,
			secondBundlePath: bundlePathForCreate,
			expectDiff:       false,
		},
		"no diff in bundles second test": {
			firstBundlePath:  bundlePathForUpdate,
			secondBundlePath: bundlePathForUpdate,
			expectDiff:       false,
		},
		"compare two diff bundles": {
			firstBundlePath:  bundlePathForCreate,
			secondBundlePath: bundlePathForUpdate,
			expectDiff:       true,
		},
	}
	for testName, test := range tests {
		t.Run(testName, func(t *testing.T) {
			firstArrayOfBytes, err := convertLocalBundleFileIntoBytes(test.firstBundlePath)
			require.NoError(t, err)
			secondArrayOfBytes, err := convertLocalBundleFileIntoBytes(test.secondBundlePath)
			require.NoError(t, err)
			firstBundleShaHash, err := getSHAFromBundle(&edgeworkers.Bundle{Reader: bytes.NewBuffer(firstArrayOfBytes)})
			require.NoError(t, err)
			secondBundleShaHash, err := getSHAFromBundle(&edgeworkers.Bundle{Reader: bytes.NewBuffer(secondArrayOfBytes)})
			require.NoError(t, err)
			if test.expectDiff {
				assert.NotEqual(t, firstBundleShaHash, secondBundleShaHash)
			} else {
				assert.Equal(t, firstBundleShaHash, secondBundleShaHash)
			}
		})
	}

	t.Run("hash should be same when file order changes", func(t *testing.T) {
		bundleOrder1 := prepareBundleWithFiles(t, []bundleFile{
			{
				Name:    "file1",
				Content: []byte("content1"),
			},
			{
				Name:    "file2",
				Content: []byte("content2"),
			},
		})
		bundleOrder2 := prepareBundleWithFiles(t, []bundleFile{
			{
				Name:    "file2",
				Content: []byte("content2"),
			},
			{
				Name:    "file1",
				Content: []byte("content1"),
			},
		})

		firstBundleShaHash, err := getSHAFromBundle(&edgeworkers.Bundle{Reader: bundleOrder1})
		require.NoError(t, err)
		secondBundleShaHash, err := getSHAFromBundle(&edgeworkers.Bundle{Reader: bundleOrder2})
		require.NoError(t, err)

		assert.Equal(t, firstBundleShaHash, secondBundleShaHash)
	})

	t.Run("hash should be different when file name changes", func(t *testing.T) {
		bundleOrder1 := prepareBundleWithFiles(t, []bundleFile{
			{
				Name:    "fileA",
				Content: []byte("content1"),
			},
			{
				Name:    "fileB",
				Content: []byte("content2"),
			},
		})
		bundleOrder2 := prepareBundleWithFiles(t, []bundleFile{
			{
				Name:    "fileA",
				Content: []byte("content1"),
			},
			{
				Name:    "fileC",
				Content: []byte("content2"),
			},
		})

		firstBundleShaHash, err := getSHAFromBundle(&edgeworkers.Bundle{Reader: bundleOrder1})
		require.NoError(t, err)
		secondBundleShaHash, err := getSHAFromBundle(&edgeworkers.Bundle{Reader: bundleOrder2})
		require.NoError(t, err)

		assert.NotEqual(t, firstBundleShaHash, secondBundleShaHash)
	})

}

func prepareBundleWithFiles(t *testing.T, files []bundleFile) io.Reader {
	bundlebuf := &bytes.Buffer{}
	gw := gzip.NewWriter(bundlebuf)
	defer func() {
		err := gw.Close()
		assert.NoError(t, err)
	}()
	tw := tar.NewWriter(gw)
	defer func() {
		err := tw.Close()
		assert.NoError(t, err)
	}()

	for _, v := range files {
		err := tw.WriteHeader(&tar.Header{
			Name: v.Name,
			Size: int64(len(v.Content)),
		})
		assert.NoError(t, err)
		_, err = tw.Write(v.Content)
		assert.NoError(t, err)
	}

	return bundlebuf
}
