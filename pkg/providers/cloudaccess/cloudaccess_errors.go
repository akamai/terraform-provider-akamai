package cloudaccess

import "errors"

var (
	// ErrCloudAccessKey is returned when operation on cloud access key fails
	ErrCloudAccessKey = errors.New("cloud access key")
	// ErrCloudAccessKeyProperties is returned when operation on cloud access key properties fails
	ErrCloudAccessKeyProperties = errors.New("cloud access key properties")
)
