// Package tools is where some legacy provider functions were dropped
package tools

import (
	"crypto/sha1"
	"encoding/hex"
	"time"

	"github.com/hashicorp/terraform-plugin-sdk/v2/diag"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/schema"
)

// GetSHAString returns a sha1 from the string
// TODO: utils should not exist, we should split this file into separate files, e.g. sha.go etc
func GetSHAString(rdata string) string {
	h := sha1.New()
	h.Write([]byte(rdata))

	sha1hashtest := hex.EncodeToString(h.Sum(nil))
	return sha1hashtest
}

//Convert schema.Set to a slice of strings
func SetToStringSlice(s *schema.Set) []string {
	list := make([]string, s.Len())
	for i, v := range s.List() {
		list[i] = v.(string)
	}
	return list
}

// MaxDuration returns the larger of x or y.
func MaxDuration(x, y time.Duration) time.Duration {
	if x < y {
		return y
	}
	return x
}

// DiagsWithErrors appends several errors to a diag.Diagnostics
func DiagsWithErrors(d diag.Diagnostics, errs ...error) diag.Diagnostics {
	for _, e := range errs {
		d = append(d, diag.FromErr(e)...)
	}
	return d
}
