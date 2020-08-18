package tools

import (
	"crypto/sha1"
	"encoding/hex"
	"github.com/google/uuid"
	"log"
)

// TODO: utils should not exist, we should split this file into separate files, e.g. sha.go etc
func GetSHAString(rdata string) string {
	h := sha1.New()
	h.Write([]byte(rdata))

	sha1hashtest := hex.EncodeToString(h.Sum(nil))
	return sha1hashtest
}

func CreateNonce() string {
	uuid, err := uuid.NewRandom()
	if err != nil {
		log.Printf("[DEBUG] Generate Uuid failed %s", err)
		return ""
	}
	return uuid.String()
}
