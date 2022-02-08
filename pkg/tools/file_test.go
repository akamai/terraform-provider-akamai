package tools

import (
	"fmt"
	"io"
	"net/http"
	"os"
	"regexp"
	"testing"

	"github.com/stretchr/testify/mock"
	"github.com/stretchr/testify/require"
	"github.com/tj/assert"
)

type (
	mockfiledownloader struct {
		mock.Mock
	}
	mockreadcloser struct {
		mock.Mock
	}
)

func (c *mockreadcloser) Close() error {
	return c.Called().Error(0)
}

func (c *mockreadcloser) Read(p []byte) (int, error) {
	return c.Called(p).Get(0).(int), c.Called(p).Error(1)
}

func (d *mockfiledownloader) HTTPGet(url string) (*http.Response, error) {
	args := d.Called(url)
	if args.Get(0) == nil {
		return nil, args.Error(1)
	}
	return args.Get(0).(*http.Response), args.Error(1)
}

func (d *mockfiledownloader) CreateFile(name string) (*os.File, error) {
	args := d.Called(name)
	if args.Get(0) == nil {
		return nil, args.Error(1)
	}
	return args.Get(0).(*os.File), args.Error(1)
}

func (d *mockfiledownloader) Copy(writer io.Writer, reader io.Reader) (int64, error) {
	args := d.Called(writer, reader)
	return args.Get(0).(int64), args.Error(1)
}

func (d *mockfiledownloader) Close(file *os.File) error {
	return file.Close()
}

func TestDownloadFile(t *testing.T) {
	errorMsg := "an error"
	tests := map[string]struct {
		init      func(*mockfiledownloader)
		filepath  string
		url       string
		withError *regexp.Regexp
	}{
		"http GET error": {
			init: func(d *mockfiledownloader) {
				d.On("HTTPGet", "https://example.com").Return(
					nil, fmt.Errorf(errorMsg)).Once()
			},
			url:       "https://example.com",
			withError: regexp.MustCompile(errorMsg),
		},
		"os Create error": {
			init: func(d *mockfiledownloader) {
				body := &mockreadcloser{}
				body.On("Close").Return(nil).Once()
				d.On("HTTPGet", "https://example.com").Return(
					&http.Response{Body: body}, nil).Once()
				d.On("CreateFile", "./relative/path").Return(
					nil, fmt.Errorf(errorMsg)).Once()
			},
			url:       "https://example.com",
			filepath:  "./relative/path",
			withError: regexp.MustCompile(errorMsg),
		},
		"io Copy error": {
			init: func(d *mockfiledownloader) {
				body := &mockreadcloser{}
				body.On("Close").Return(nil).Once()
				d.On("HTTPGet", "https://example.com").Return(
					&http.Response{Body: body}, nil).Once()
				file := &os.File{}
				d.On("CreateFile", "./relative/path").Return(
					file, nil).Once()
				d.On("Copy", file, body).Return(int64(0), fmt.Errorf(errorMsg))
			},
			url:       "https://example.com",
			filepath:  "./relative/path",
			withError: regexp.MustCompile(errorMsg),
		},
		"resp Body Close error": {
			init: func(d *mockfiledownloader) {
				body := &mockreadcloser{}
				body.On("Close").Return(fmt.Errorf(errorMsg)).Once()
				d.On("HTTPGet", "https://example.com").Return(
					&http.Response{Body: body}, nil).Once()
				file := &os.File{}
				d.On("CreateFile", "./relative/path").Return(
					file, nil).Once()
				d.On("Copy", file, body).Return(int64(0), nil).Once()
				d.On("Close", file).Return(nil).Once()
			},
			url:       "https://example.com",
			filepath:  "./relative/path",
			withError: regexp.MustCompile(errorMsg),
		},
		"all fine": {
			init: func(d *mockfiledownloader) {
				body := &mockreadcloser{}
				body.On("Close").Return(nil)
				d.On("HTTPGet", "https://example.com").Return(
					&http.Response{Body: body}, nil).Once()
				file := &os.File{}
				d.On("CreateFile", "./relative/path").Return(
					file, nil).Once()
				d.On("Copy", file, body).Return(int64(0), nil)
			},
			url:      "https://example.com",
			filepath: "./relative/path",
		},
	}
	for name, test := range tests {
		t.Run(name, func(t *testing.T) {
			d := mockfiledownloader{}
			test.init(&d)
			err := DownloadFile(&d, test.filepath, test.url)
			if test.withError != nil {
				require.Error(t, err)
				assert.Regexp(t, test.withError, err.Error())
				return
			}
			require.NoError(t, err)
			d.AssertExpectations(t)
		})
	}
}
