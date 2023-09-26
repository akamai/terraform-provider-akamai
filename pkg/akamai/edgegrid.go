package akamai

import (
	"errors"
	"fmt"

	"github.com/akamai/AkamaiOPEN-edgegrid-golang/v7/pkg/edgegrid"
)

// ErrWrongEdgeGridConfiguration is returned when the configuration could not be read
var ErrWrongEdgeGridConfiguration = errors.New("error reading Akamai EdgeGrid configuration")

// DefaultConfigFilePath is the default path for edgerc config file
var DefaultConfigFilePath = edgegrid.DefaultConfigFile

type configBearer struct {
	accessToken  string
	accountKey   string
	clientSecret string
	clientToken  string
	host         string
	maxBody      int
}

func (c configBearer) toEdgegridConfig() (*edgegrid.Config, error) {
	if !c.valid() {
		return nil, ErrWrongEdgeGridConfiguration
	}

	edgerc := &edgegrid.Config{
		AccessToken:  c.accessToken,
		AccountKey:   c.accountKey,
		ClientSecret: c.clientSecret,
		ClientToken:  c.clientToken,
		Host:         c.host,
		MaxBody:      c.maxBody,
	}

	if edgerc.MaxBody == 0 {
		edgerc.MaxBody = edgegrid.MaxBodySize
	}

	return edgerc, nil
}

func (c configBearer) valid() bool {
	return c.host != "" && c.accessToken != "" && c.clientSecret != "" && c.clientToken != ""
}

// newEdgegridConfig creates a new edgegrid.Config based on provided arguments.
//
// It evaluates possibility of creating the config in the following order:
//  1. Environmental variables
//  2. Config block
//  3. Edgerc file
//
// If edgerc path or section are not provided, it uses the edgegrid defaults.
func newEdgegridConfig(path, section string, config configBearer) (*edgegrid.Config, error) {
	envEdgerc := &edgegrid.Config{}
	err := envEdgerc.FromEnv(edgercSectionOrDefault(section))
	if err == nil {
		return validateEdgerc(envEdgerc)
	}

	configEdgerc, err := config.toEdgegridConfig()
	if err == nil {
		return validateEdgerc(configEdgerc)
	}

	fileEdgerc := &edgegrid.Config{}
	err = fileEdgerc.FromFile(edgercPathOrDefault(path), edgercSectionOrDefault(section))
	if err == nil {
		return validateEdgerc(fileEdgerc)
	}
	return nil, fmt.Errorf("%w: %s", ErrWrongEdgeGridConfiguration, err)
}

func validateEdgerc(edgerc *edgegrid.Config) (*edgegrid.Config, error) {
	if err := edgerc.Validate(); err != nil {
		return nil, fmt.Errorf("%w: %s", ErrWrongEdgeGridConfiguration, err)
	}

	return edgerc, nil
}

func edgercPathOrDefault(path string) string {
	if path == "" {
		return DefaultConfigFilePath
	}
	return path
}

func edgercSectionOrDefault(section string) string {
	if section == "" {
		return edgegrid.DefaultSection
	}
	return section
}
