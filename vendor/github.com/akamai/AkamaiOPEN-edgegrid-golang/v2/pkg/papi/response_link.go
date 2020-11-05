// Package papi provides access to the Akamai Property APIs
package papi

import (
	"errors"
	"net/url"
	"strings"
)

var (
	// ErrInvalidResponseLink is returned when there was an error while fetching ID from location response object
	ErrInvalidResponseLink = errors.New("response link URL is invalid")
)

// ResponseLinkParse parse the link and returns the id
func ResponseLinkParse(link string) (string, error) {
	locURL, err := url.Parse(string(link))
	if err != nil {
		return "", err
	}
	pathSplit := strings.Split(locURL.Path, "/")
	return pathSplit[len(pathSplit)-1], nil
}
