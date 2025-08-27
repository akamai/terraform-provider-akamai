package apidefinitions

import (
	"context"
	"errors"
	"fmt"
	"slices"
	"time"

	"github.com/akamai/AkamaiOPEN-edgegrid-golang/v12/pkg/apidefinitions"
	"github.com/akamai/terraform-provider-akamai/v9/pkg/common/date"
	"github.com/hashicorp/terraform-plugin-framework/diag"
	"github.com/hashicorp/terraform-plugin-log/tflog"
)

var (
	pollInterval    = time.Second * 30
	activationRetry = time.Second * 5
)

func startActivation(ctx context.Context, activationRequest apidefinitions.ActivateVersionRequest) error {

	activationRetry := activationRetry

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

func startDeactivation(ctx context.Context, deactivationRequest apidefinitions.DeactivateVersionRequest) error {

	deactivationRetry := activationRetry

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

func deactivateEndpoint(ctx context.Context, endpoint apidefinitions.EndpointDetail) diag.Diagnostics {
	return deactivateOnNetworks(ctx, endpoint, []apidefinitions.NetworkType{apidefinitions.ActivationNetworkStaging, apidefinitions.ActivationNetworkProduction})
}

func deactivateEndpointOnNetwork(ctx context.Context, endpointID int64, network apidefinitions.NetworkType) diag.Diagnostics {
	return deactivateEndpointOnNetworks(ctx, endpointID, []apidefinitions.NetworkType{network})
}

func deactivateEndpointOnNetworks(ctx context.Context, endpointID int64, scope []apidefinitions.NetworkType) diag.Diagnostics {
	var diags diag.Diagnostics

	endpoint, err := getEndpoint(ctx, endpointID)
	if err != nil {
		diags.AddError("Unable to read Endpoint", err.Error())
		return diags
	}

	return deactivateOnNetworks(ctx, *endpoint, scope)
}

func deactivateOnNetworks(ctx context.Context, endpoint apidefinitions.EndpointDetail, scope []apidefinitions.NetworkType) diag.Diagnostics {
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
		err := startDeactivation(ctx, request)
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
		err := startDeactivation(ctx, request)
		if err != nil {
			diags.AddError("Deactivation on Production Failed", err.Error())
			return diags
		}
		networkToDeactivate = append(networkToDeactivate, apidefinitions.ActivationNetworkProduction)
	}

	if len(networkToDeactivate) > 0 {
		pollDeactivation(ctx, endpoint.APIEndpointID, networkToDeactivate)
	}
	return nil
}

func getEndpoint(ctx context.Context, endpointID int64) (*apidefinitions.EndpointDetail, error) {
	endpoint, err := client.GetEndpoint(ctx, apidefinitions.GetEndpointRequest{APIEndpointID: endpointID})
	if err != nil {
		return nil, err
	}

	return (*apidefinitions.EndpointDetail)(endpoint), nil
}

func pollActivation(ctx context.Context, endpointID, version int64, network apidefinitions.NetworkType) (*apidefinitions.EndpointDetail, diag.Diagnostics) {
	var diags diag.Diagnostics

	for {
		select {
		case <-time.After(pollInterval):
			endpoint, err := getEndpoint(ctx, endpointID)
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

func pollDeactivation(ctx context.Context, endpointID int64, networkToDeactivate []apidefinitions.NetworkType) (*apidefinitions.EndpointDetail, diag.Diagnostics) {
	var diags diag.Diagnostics
	for {
		select {
		case <-time.After(pollInterval):
			endpoint, err := getEndpoint(ctx, endpointID)
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
