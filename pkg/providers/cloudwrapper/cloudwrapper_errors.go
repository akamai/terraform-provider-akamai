package cloudwrapper

import "errors"

var (
	// ErrCloudWrapperLocation is returned when operation on cloudwrapper location fails
	ErrCloudWrapperLocation = errors.New("cloudwrapper location")
	// ErrCloudWrapperLocations is returned when operation on cloudwrapper locations fails
	ErrCloudWrapperLocations = errors.New("cloudwrapper locations")
)
