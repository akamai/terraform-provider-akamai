package iam

import "errors"

var (
	// ErrIAMListAllowedAPIs is returned when ListAllowedAPIs fails
	ErrIAMListAllowedAPIs = errors.New("IAM list allowed APIs failed")
)
