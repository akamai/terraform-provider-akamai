package files

import (
	"testing"

	"github.com/akamai/terraform-provider-akamai/v8/pkg/common/ptr"
	"github.com/stretchr/testify/assert"
)

func TestIsSymlink(t *testing.T) {
	tests := map[string]struct {
		path          string
		isSymlink     bool
		expectedError *string
	}{
		"regular file": {
			path:      "testdata/lorem-ipsum.txt",
			isSymlink: false,
		},
		"symlink to file": {
			path:      "testdata/link.txt",
			isSymlink: true,
		},
		"directory": {
			path:      "testdata/somedir",
			isSymlink: false,
		},
		"symlink to directory": {
			path:      "testdata/dirlink",
			isSymlink: true,
		},
		"error if no such file": {
			path:          "testdata/notexisting.txt",
			expectedError: ptr.To("testdata/notexisting.txt: no such file or directory"),
		},
	}

	for name, test := range tests {
		t.Run(name, func(t *testing.T) {
			isLink, err := IsSymlink(test.path)
			if test.expectedError != nil {
				assert.ErrorContains(t, err, *test.expectedError)
			} else {
				assert.NoError(t, err)
			}
			assert.Equal(t, test.isSymlink, isLink)
		})
	}
}

func TestIsSymlinkToDir(t *testing.T) {
	tests := map[string]struct {
		path          string
		isDirSymlink  bool
		expectedError *string
	}{
		"regular file": {
			path:         "testdata/lorem-ipsum.txt",
			isDirSymlink: false,
		},
		"symlink to file": {
			path:         "testdata/link.txt",
			isDirSymlink: false,
		},
		"directory": {
			path:         "testdata/somedir",
			isDirSymlink: false,
		},
		"symlink to directory": {
			path:         "testdata/dirlink",
			isDirSymlink: true,
		},
		"error if no such file": {
			path:          "testdata/notexisting.txt",
			expectedError: ptr.To("testdata/notexisting.txt: no such file or directory"),
		},
	}

	for name, test := range tests {
		t.Run(name, func(t *testing.T) {
			isLink, err := IsSymlinkToDir(test.path)
			if test.expectedError != nil {
				assert.ErrorContains(t, err, *test.expectedError)
			} else {
				assert.NoError(t, err)
			}
			assert.Equal(t, test.isDirSymlink, isLink)
		})
	}
}
