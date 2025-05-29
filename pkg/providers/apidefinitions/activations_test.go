package apidefinitions

import (
	"context"
	"testing"

	"github.com/akamai/AkamaiOPEN-edgegrid-golang/v11/pkg/apidefinitions"
	"github.com/stretchr/testify/assert"
)

func TestStartActivation_shouldRetry(t *testing.T) {
	mock := &apidefinitions.Mock{}
	client = mock
	mockActivateVersionFail(mock)
	mockActivateVersion(mock)
	err := startActivation(context.TODO(), apidefinitions.ActivateVersionRequest{
		Body: apidefinitions.ActivationRequestBody{
			Networks: []apidefinitions.NetworkType{apidefinitions.ActivationNetworkStaging},
		},
	})

	assert.Nil(t, err)
}

func TestStartDeactivation_shouldRetry(t *testing.T) {
	mock := &apidefinitions.Mock{}
	client = mock
	mockDeactivateVersionFail(mock)
	mockDeactivateVersion(mock, 1)
	err := startDeactivation(context.TODO(), apidefinitions.DeactivateVersionRequest{
		Body: apidefinitions.ActivationRequestBody{
			Networks: []apidefinitions.NetworkType{apidefinitions.ActivationNetworkStaging},
		},
	})

	assert.Nil(t, err)
}
