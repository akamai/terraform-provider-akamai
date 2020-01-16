package apiendpoints

import (
	"fmt"

	"github.com/akamai/AkamaiOPEN-edgegrid-golang/client-v1"
)

type Activation struct {
	Networks               []string `json:"networks"`
	NotificationRecipients []string `json:"notificationRecipients"`
	Notes                  string   `json:"notes"`
}

type ActivateEndpointOptions struct {
	APIEndPointId int
	VersionNumber int
}

func ActivateEndpoint(options *ActivateEndpointOptions, activation *Activation) (*Activation, error) {
	req, err := client.NewJSONRequest(
		Config,
		"POST",
		fmt.Sprintf(
			"/api-definitions/v2/endpoints/%d/versions/%d/activate",
			options.APIEndPointId,
			options.VersionNumber,
		),
		activation,
	)

	if err != nil {
		return nil, err
	}

	res, err := client.Do(Config, req)

	if client.IsError(res) {
		return nil, client.NewAPIError(res)
	}

	return activation, nil
}

func DeactivateEndpoint(options *ActivateEndpointOptions, activation *Activation) (*Activation, error) {
	req, err := client.NewJSONRequest(
		Config,
		"DELETE",
		fmt.Sprintf(
			"/api-definitions/v2/endpoints/%d/versions/%d/deactivate",
			options.APIEndPointId,
			options.VersionNumber,
		),
		activation,
	)

	if err != nil {
		return nil, err
	}

	res, err := client.Do(Config, req)

	if client.IsError(res) {
		return nil, client.NewAPIError(res)
	}

	return activation, nil
}

func IsActive(endpoint *Endpoint, network string) bool {
	if network == "production" {
		if endpoint.ProductionStatus == StatusPending || endpoint.ProductionStatus == StatusActive {
			return true
		}
	}

	if network == "staging" {
		if endpoint.StagingStatus == StatusPending || endpoint.StagingStatus == StatusActive {
			return true
		}
	}

	return false
}
