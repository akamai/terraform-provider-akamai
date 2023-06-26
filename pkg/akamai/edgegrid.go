package akamai

import (
	"errors"
	"fmt"

	"github.com/akamai/AkamaiOPEN-edgegrid-golang/v7/pkg/edgegrid"
)

// ErrWrongEdgeGridConfiguration is returned when the configuration could not be read
var ErrWrongEdgeGridConfiguration = errors.New("error reading Akamai EdgeGrid configuration")

func newEdgegridConfig(path, section string, config map[string]any) (*edgegrid.Config, error) {
	if (path != "" || section != "") && len(config) > 0 {
		return nil, fmt.Errorf("edgegrid cannot be simultaneously configured with file and config map") // should not happen as schema guarantees that
	}

	var edgerc *edgegrid.Config
	if len(config) > 0 {
		edgerc = &edgegrid.Config{
			Host:         config["host"].(string),
			AccessToken:  config["access_token"].(string),
			ClientToken:  config["client_token"].(string),
			ClientSecret: config["client_secret"].(string),
			MaxBody:      config["max_body"].(int),
			AccountKey:   config["account_key"].(string),
		}
	} else {
		edgerc = &edgegrid.Config{}
		err := edgerc.FromFile(edgercPathOrDefault(path), edgercSectionOrDefault(section))
		if err != nil {
			return nil, fmt.Errorf("%w: %s", ErrWrongEdgeGridConfiguration, err)
		}
	}

	if err := edgerc.Validate(); err != nil {
		return nil, fmt.Errorf("%w: %s", ErrWrongEdgeGridConfiguration, err)
	}

	return edgerc, nil
}

func edgercPathOrDefault(path string) string {
	if path == "" {
		return edgegrid.DefaultConfigFile
	}
	return path
}

func edgercSectionOrDefault(section string) string {
	if section == "" {
		return edgegrid.DefaultSection
	}
	return section
}
