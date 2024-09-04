package iam

import "errors"

var (
	// ErrIAMListAllowedAPIs is returned when ListAllowedAPIs fails
	ErrIAMListAllowedAPIs = errors.New("IAM list allowed APIs failed")
	// ErrIAMGetCIDRBlock is returned when GetCIDRBlock fails
	ErrIAMGetCIDRBlock = errors.New("IAM get CIDR block failed")
	// ErrIAMListCIDRBlocks is returned when ListCIDRBlocks fails
	ErrIAMListCIDRBlocks = errors.New("IAM list CIDR blocks failed")
)
