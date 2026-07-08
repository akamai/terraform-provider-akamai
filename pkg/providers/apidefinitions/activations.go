package apidefinitions

import (
	"context"
	"errors"
	"fmt"
	"slices"
	"time"

	"github.com/akamai/AkamaiOPEN-edgegrid-golang/v13/pkg/apidefinitions"
	"github.com/akamai/terraform-provider-akamai/v10/pkg/common/date"
	"github.com/hashicorp/terraform-plugin-framework/diag"
	"github.com/hashicorp/terraform-plugin-log/tflog"
)

func startActivation(ctx context.Context, client apidefinitions.APIDefinitions, activationRequest apidefinitions.ActivateVersionRequest, activationRetry time.Duration) error {

	for {
		tflog.Debug(ctx, "starting activation")
		_, err := client.ActivateVersion(ctx, activationRequest)

		if err == nil {
			return nil
		}

		if !isErrorRetryable(err) {
			return fmt.Errorf("activation failed: %s", err)
		}

		tflog.Debug(ctx, "retrying activation")

		select {
		case <-time.After(activationRetry):
			activationRetry = date.CapDuration(activationRetry*2, 5*time.Minute)
			continue

		case <-ctx.Done():
			return fmt.Errorf("activation context terminated: %w", ctx.Err())
		}
	}
}

func startDeactivation(ctx context.Context, client apidefinitions.APIDefinitions, deactivationRequest apidefinitions.DeactivateVersionRequest, deactivationRetry time.Duration) error {

	for {
		tflog.Debug(ctx, "starting deactivation")
		_, err := client.DeactivateVersion(ctx, deactivationRequest)

		if err == nil {
			return nil
		}

		if !isErrorRetryable(err) {
			return fmt.Errorf("deactivation failed: %s", err)
		}

		tflog.Debug(ctx, "retrying deactivation")

		select {
		case <-time.After(deactivationRetry):
			deactivationRetry = date.CapDuration(deactivationRetry*2, 5*time.Minute)
			continue

		case <-ctx.Done():
			return fmt.Errorf("deactivation context terminated: %w", ctx.Err())
		}
	}

}

func deactivateEndpoint(ctx context.Context, client apidefinitions.APIDefinitions, endpoint apidefinitions.EndpointDetail, deactivationRetry, pollInterval time.Duration) diag.Diagnostics {
	return deactivateOnNetworks(ctx, client, endpoint, []apidefinitions.NetworkType{apidefinitions.ActivationNetworkStaging, apidefinitions.ActivationNetworkProduction}, deactivationRetry, pollInterval)
}

func deactivateEndpointOnNetwork(ctx context.Context, client apidefinitions.APIDefinitions, endpointID int64, network apidefinitions.NetworkType, deactivationRetry, pollInterval time.Duration) diag.Diagnostics {
	return deactivateEndpointOnNetworks(ctx, client, endpointID, []apidefinitions.NetworkType{network}, deactivationRetry, pollInterval)
}

func deactivateEndpointOnNetworks(ctx context.Context, client apidefinitions.APIDefinitions, endpointID int64, scope []apidefinitions.NetworkType, deactivationRetry, pollInterval time.Duration) diag.Diagnostics {
	var diags diag.Diagnostics

	endpoint, err := getEndpoint(ctx, client, endpointID)
	if err != nil {
		diags.AddError("Unable to read Endpoint", err.Error())
		return diags
	}

	return deactivateOnNetworks(ctx, client, *endpoint, scope, deactivationRetry, pollInterval)
}

func deactivateOnNetworks(ctx context.Context, client apidefinitions.APIDefinitions, endpoint apidefinitions.EndpointDetail, scope []apidefinitions.NetworkType, deactivationRetry, pollInterval time.Duration) diag.Diagnostics {
	var diags diag.Diagnostics

	var networkToDeactivate []apidefinitions.NetworkType

	if isEndpointActive(endpoint.StagingVersion) && slices.Contains(scope, apidefinitions.ActivationNetworkStaging) {
		request := apidefinitions.DeactivateVersionRequest{
			APIEndpointID: endpoint.APIEndpointID,
			VersionNumber: *endpoint.StagingVersion.VersionNumber,
			Body: apidefinitions.ActivationRequestBody{
				Networks: []apidefinitions.NetworkType{apidefinitions.ActivationNetworkStaging},
			},
		}
		err := startDeactivation(ctx, client, request, deactivationRetry)
		if err != nil {
			diags.AddError("Deactivation on Staging Failed", err.Error())
			return diags
		}

		networkToDeactivate = append(networkToDeactivate, apidefinitions.ActivationNetworkStaging)
	}

	if isEndpointActive(endpoint.ProductionVersion) && slices.Contains(scope, apidefinitions.ActivationNetworkProduction) {
		request := apidefinitions.DeactivateVersionRequest{
			APIEndpointID: endpoint.APIEndpointID,
			VersionNumber: *endpoint.ProductionVersion.VersionNumber,
			Body: apidefinitions.ActivationRequestBody{
				Networks: []apidefinitions.NetworkType{apidefinitions.ActivationNetworkProduction},
			},
		}
		err := startDeactivation(ctx, client, request, deactivationRetry)
		if err != nil {
			diags.AddError("Deactivation on Production Failed", err.Error())
			return diags
		}
		networkToDeactivate = append(networkToDeactivate, apidefinitions.ActivationNetworkProduction)
	}

	if len(networkToDeactivate) > 0 {
		pollDeactivation(ctx, client, endpoint.APIEndpointID, networkToDeactivate, pollInterval)
	}
	return nil
}

func getEndpoint(ctx context.Context, client apidefinitions.APIDefinitions, endpointID int64) (*apidefinitions.EndpointDetail, error) {
	endpoint, err := client.GetEndpoint(ctx, apidefinitions.GetEndpointRequest{APIEndpointID: endpointID})
	if err != nil {
		return nil, err
	}

	return (*apidefinitions.EndpointDetail)(endpoint), nil
}

func pollActivation(ctx context.Context, client apidefinitions.APIDefinitions, endpointID, version int64, network apidefinitions.NetworkType, pollInterval time.Duration) (*apidefinitions.EndpointDetail, diag.Diagnostics) {
	var diags diag.Diagnostics

	for {
		select {
		case <-time.After(pollInterval):
			endpoint, err := getEndpoint(ctx, client, endpointID)
			if err != nil {
				continue
			}
			if isActivationAccordingToTheState(version, network, *endpoint) {
				return endpoint, nil
			}
			if hasFailed(network, *endpoint) {
				diags.AddError("Activation Failed", fmt.Sprintf("Activation for version %v failed", version))
				return nil, diags
			}
		case <-ctx.Done():
			diags.AddError("activation context terminated: %w", ctx.Err().Error())
			return nil, diags
		}
	}
}

func pollDeactivation(ctx context.Context, client apidefinitions.APIDefinitions, endpointID int64, networkToDeactivate []apidefinitions.NetworkType, pollInterval time.Duration) (*apidefinitions.EndpointDetail, diag.Diagnostics) {
	var diags diag.Diagnostics
	for {
		select {
		case <-time.After(pollInterval):
			endpoint, err := getEndpoint(ctx, client, endpointID)
			if err != nil {
				continue
			}
			if isDeactivationFinished(*endpoint, networkToDeactivate) {
				return endpoint, nil
			}
			if hasDeactivationFailed(*endpoint, networkToDeactivate) {
				diags.AddError("Deactivation Failed", fmt.Sprintf("Deactivation for endpoint %d failed", endpointID))
				return nil, diags
			}
		case <-ctx.Done():
			diags.AddError("deactivation context terminated: %w", ctx.Err().Error())
			return nil, diags
		}
	}
}

func isDeactivationFinished(endpoint apidefinitions.EndpointDetail, networks []apidefinitions.NetworkType) bool {
	var state apidefinitions.VersionState

	if slices.Contains(networks, apidefinitions.ActivationNetworkStaging) && slices.Contains(networks, apidefinitions.ActivationNetworkProduction) {
		return isDeactivated(endpoint.StagingVersion) && isDeactivated(endpoint.ProductionVersion)
	} else if slices.Contains(networks, apidefinitions.ActivationNetworkStaging) {
		state = getStateOnNetwork(apidefinitions.ActivationNetworkStaging, endpoint)
	} else {
		state = getStateOnNetwork(apidefinitions.ActivationNetworkProduction, endpoint)
	}
	return isDeactivated(state)
}

func isDeactivated(state apidefinitions.VersionState) bool {
	return state.Status != nil && *state.Status == apidefinitions.ActivationStatusDeactivated
}

func hasDeactivationFailed(state apidefinitions.EndpointDetail, networks []apidefinitions.NetworkType) bool {
	for _, network := range networks {
		status := getStatusOnNetwork(network, state)
		if status != nil && *status == apidefinitions.ActivationStatusFailed {
			return true
		}
	}
	return false
}

func hasFailed(network apidefinitions.NetworkType, state apidefinitions.EndpointDetail) bool {
	status := getStatusOnNetwork(network, state)
	return status != nil && *status == apidefinitions.ActivationStatusFailed
}

func isErrorRetryable(err error) bool {
	var responseErr = &apidefinitions.Error{}
	if !errors.As(err, &responseErr) {
		return false
	}
	if responseErr.Status < 500 {
		return false
	}
	return true
}

func shouldActivate(versionStatus apidefinitions.VersionState, versionToActivate int64) bool {
	if versionStatus.VersionNumber == nil {
		return true
	}
	return *versionStatus.VersionNumber != versionToActivate || *versionStatus.Status != apidefinitions.ActivationStatusActive
}

func isVersionActive(status apidefinitions.VersionState, versionNumber int64) bool {
	if status.VersionNumber == nil {
		return false
	}
	return *status.VersionNumber == versionNumber && *status.Status == apidefinitions.ActivationStatusActive
}

func isEndpointActive(status apidefinitions.VersionState) bool {
	if status.VersionNumber == nil {
		return false
	}
	return *status.Status == apidefinitions.ActivationStatusActive
}

func isActivationAccordingToTheState(version int64, network apidefinitions.NetworkType, state apidefinitions.EndpointDetail) bool {
	if network == apidefinitions.ActivationNetworkStaging {
		return isVersionActive(state.StagingVersion, version)
	}
	return isVersionActive(state.ProductionVersion, version)
}

func getStateOnNetwork(network apidefinitions.NetworkType, endpoint apidefinitions.EndpointDetail) apidefinitions.VersionState {
	var activation = endpoint.ProductionVersion
	if network == apidefinitions.ActivationNetworkStaging {
		activation = endpoint.StagingVersion
	}
	return activation
}

func getStatusOnNetwork(network apidefinitions.NetworkType, endpoint apidefinitions.EndpointDetail) *apidefinitions.ActivationStatus {
	return getStateOnNetwork(network, endpoint).Status
}
