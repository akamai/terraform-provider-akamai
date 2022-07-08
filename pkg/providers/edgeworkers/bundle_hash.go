package edgeworkers

import (
	"archive/tar"
	"compress/gzip"
	"crypto/sha256"
	"encoding/hex"
	"io"
	"sort"

	"github.com/akamai/AkamaiOPEN-edgegrid-golang/v2/pkg/edgeworkers"
)

type bundleFile struct {
	Name    string
	Content []byte
}

func hashTarFiles(tr *tar.Reader) ([]bundleFile, error) {
	var filesHashes = make([]bundleFile, 0)
	for {
		header, err := tr.Next()
		if err == io.EOF {
			break
		}
		if err != nil {
			return nil, err
		}
		hash := sha256.New()
		_, err = io.Copy(hash, tr)
		if err != nil && err != io.EOF {
			return nil, err
		}
		filesHashes = append(filesHashes, bundleFile{
			Name:    header.Name,
			Content: hash.Sum(nil),
		})
	}

	return filesHashes, nil
}

func sortBundleFilesByNames(files []bundleFile) {
	sort.Slice(files, func(i, j int) bool {
		return files[i].Name < files[j].Name
	})
}

func calculateHash(files []bundleFile) []byte {
	hash := sha256.New()
	for _, f := range files {
		hash.Write([]byte(f.Name))
		hash.Write(f.Content)
	}
	return hash.Sum(nil)
}

func getSHAFromBundle(bundleContent *edgeworkers.Bundle) (string, error) {
	gr, err := gzip.NewReader(bundleContent)
	if err != nil {
		return "", err
	}
	gr.Multistream(false)

	tr := tar.NewReader(gr)
	filesHashes, err := hashTarFiles(tr)
	if err != nil {
		return "", nil
	}

	sortBundleFilesByNames(filesHashes)
	sum := calculateHash(filesHashes)

	return hex.EncodeToString(sum), nil
}
