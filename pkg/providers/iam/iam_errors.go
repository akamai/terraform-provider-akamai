package iam

import "errors"

var (
	// ErrIAMListAllowedAPIs is returned when ListAllowedAPIs fails
	ErrIAMListAllowedAPIs = errors.New("IAM list allowed APIs failed")
	// ErrIAMGetCIDRBlock is returned when GetCIDRBlock fails
	ErrIAMGetCIDRBlock = errors.New("IAM get CIDR block failed")
	// ErrIAMListUsers is returned when ListUsers fails
	ErrIAMListUsers = errors.New("IAM list users failed")
	// ErrIAMGetUser is returned when GetUser fails
	ErrIAMGetUser = errors.New("IAM get user failed")
	// ErrIAMListCIDRBlocks is returned when ListCIDRBlocks fails
	ErrIAMListCIDRBlocks = errors.New("IAM list CIDR blocks failed")
)
