package cloudaccess

import "errors"

var (
	// ErrCloudAccessKey is returned when operation on cloud access key fails
	ErrCloudAccessKey = errors.New("cloud access key")
	// ErrCloudAccessKeys is returned when operation on cloud access keys fails
	ErrCloudAccessKeys = errors.New("cloud access keys")
	// ErrCloudAccessKeyVersions is returned when operation on cloud access key versions fails
	ErrCloudAccessKeyVersions = errors.New("cloud access key versions")
	// ErrCloudAccessKeyProperties is returned when operation on cloud access key properties fails
	ErrCloudAccessKeyProperties = errors.New("cloud access key properties")
)
