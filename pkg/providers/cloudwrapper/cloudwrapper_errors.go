package cloudwrapper

import "errors"

var (
	// ErrCloudWrapperProperties is returned when operation on cloudwrapper properties fails
	ErrCloudWrapperProperties = errors.New("cloudwrapper properties")
	// ErrCloudWrapperLocation is returned when operation on cloudwrapper location fails
	ErrCloudWrapperLocation = errors.New("cloudwrapper location")
)
