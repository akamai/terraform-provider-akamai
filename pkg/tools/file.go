package tools

import (
	"io"
	"net/http"
	"os"
)

type (
	fileDownloader interface {
		// HTTPGet issues a GET to the specified URL.
		// It is basically a wrapper for http.Get()
		HTTPGet(string) (*http.Response, error)

		// CreateFile creates or truncates the named file.
		// It is basically a wrapper for os.Create()
		CreateFile(string) (*os.File, error)

		// Copy copies from src to dst until either EOF is reached on src or an error occurs.
		// It is basically a wrapper for io.Copy()
		Copy(io.Writer, io.Reader) (int64, error)

		// Close closes the File, rendering it unusable for I/O.
		// It is basically a wrapper for &os.File.Close()
		Close(*os.File) error
	}

	// FileDownloader represents object for file downloading
	FileDownloader struct{}
)

// DownloadFile downloads a file from the given URL and saves it under the given path
func DownloadFile(fileDownloader fileDownloader, filepath, url string) error {
	resp, err := fileDownloader.HTTPGet(url)
	if err != nil {
		return err
	}
	out, err := fileDownloader.CreateFile(filepath)
	if err != nil {
		return err
	}
	if _, err := fileDownloader.Copy(out, resp.Body); err != nil {
		return err
	}
	if err = resp.Body.Close(); err != nil {
		return err
	}
	defer func() {
		err = fileDownloader.Close(out)
	}()
	return err
}

func (d *FileDownloader) HTTPGet(url string) (*http.Response, error) {
	return http.Get(url)
}

func (d *FileDownloader) CreateFile(name string) (*os.File, error) {
	return os.Create(name)
}

func (d *FileDownloader) Copy(writer io.Writer, reader io.Reader) (int64, error) {
	return io.Copy(writer, reader)
}

func (d *FileDownloader) Close(file *os.File) error {
	return file.Close()
}
