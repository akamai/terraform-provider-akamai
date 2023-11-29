// Package hash contains convenience functions for calculating hashes.
package hash

import (
	"bytes"
	"crypto/md5"
	"crypto/sha1"
	"encoding/gob"
	"encoding/hex"
	"fmt"
)

// GetSHAString returns a sha1 from the string
func GetSHAString(rdata string) string {
	h := sha1.New()
	h.Write([]byte(rdata))

	sha1hashtest := hex.EncodeToString(h.Sum(nil))
	return sha1hashtest
}

// GetMD5Sum calculates md5Sum for any given interface{}.
// Passing a nil pointer to GetMD5Sum will panic, as they cannot be transmitted by gob.
func GetMD5Sum(key interface{}) (string, error) {
	var buffer bytes.Buffer
	encoder := gob.NewEncoder(&buffer)
	err := encoder.Encode(key)
	if err != nil {
		return "", err
	}
	return fmt.Sprintf("%x", md5.Sum(buffer.Bytes())), nil
}
