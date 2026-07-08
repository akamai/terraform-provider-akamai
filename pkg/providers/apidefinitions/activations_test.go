package apidefinitions

import (
	"context"
	"testing"
	"time"

	"github.com/akamai/AkamaiOPEN-edgegrid-golang/v13/pkg/apidefinitions"
	"github.com/akamai/terraform-provider-akamai/v10/internal/edgegrid"
	"github.com/stretchr/testify/assert"
)

func TestStartActivation_shouldRetry(t *testing.T) {
	t.Parallel()
	config := defaultSubproviderConfig()
	config.activation.activationRetry = 10 * time.Millisecond
	client := edgegrid.NewTestClient()
	mockActivateVersionFail(client.APIDefinitions)
	mockActivateVersion(client.APIDefinitions)
	err := startActivation(context.TODO(), client.APIDefinitions, apidefinitions.ActivateVersionRequest{
		Body: apidefinitions.ActivationRequestBody{
			Networks: []apidefinitions.NetworkType{apidefinitions.ActivationNetworkStaging},
		},
	}, config.activation.activationRetry)

	client.APIDefinitions.AssertExpectations(t)
	assert.Nil(t, err)
}

func TestStartDeactivation_shouldRetry(t *testing.T) {
	t.Parallel()
	config := defaultSubproviderConfig()
	config.activation.activationRetry = 10 * time.Millisecond
	client := edgegrid.NewTestClient()
	mockDeactivateVersionFail(client.APIDefinitions)
	mockDeactivateVersion(client.APIDefinitions, 1)
	err := startDeactivation(context.TODO(), client.APIDefinitions, apidefinitions.DeactivateVersionRequest{
		Body: apidefinitions.ActivationRequestBody{
			Networks: []apidefinitions.NetworkType{apidefinitions.ActivationNetworkStaging},
		},
	}, config.activation.activationRetry)

	client.APIDefinitions.AssertExpectations(t)
	assert.Nil(t, err)
}
