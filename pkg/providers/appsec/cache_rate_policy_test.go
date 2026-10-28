package appsec

import (
	"context"
	"fmt"
	"testing"

	"github.com/akamai/AkamaiOPEN-edgegrid-golang/v14/pkg/appsec"
	akalog "github.com/akamai/AkamaiOPEN-edgegrid-golang/v14/pkg/log"
	"github.com/akamai/terraform-provider-akamai/v11/internal/edgegrid"
	"github.com/akamai/terraform-provider-akamai/v11/pkg/cache"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/mock"
	"github.com/stretchr/testify/require"
)

func TestFilterRatePolicy(t *testing.T) {
	logger := akalog.NOPLogger()

	policies := buildRatePoliciesResponse(t, `{
		"ratePolicies": [
			{"id": 11, "name": "policy-A", "type": "CONNECTION", "averageThreshold": 100},
			{"id": 22, "name": "policy-B", "type": "RATE",       "averageThreshold": 200}
		]
	}`)

	t.Run("found", func(t *testing.T) {
		result, err := filterRatePolicy(policies, 22, logger)
		require.NoError(t, err)
		require.NotNil(t, result)
		assert.Equal(t, "policy-B", result.Name)
		assert.Equal(t, "RATE", result.Type)
		assert.Equal(t, 200, result.AverageThreshold)
	})

	t.Run("not_found", func(t *testing.T) {
		result, err := filterRatePolicy(policies, 99, logger)
		assert.Nil(t, result)
		assert.ErrorIs(t, err, ErrRatePolicyNotFound)
	})
}

func TestGetRatePolicy_CacheMiss(t *testing.T) {
	clearCache()
	defer clearCache()

	client := edgegrid.NewTestClient()
	configID, version, policyID := int64(5001), int64(1), 11

	listResp := buildRatePoliciesResponse(t, `{"ratePolicies":[{"id":11,"name":"policy-A","type":"CONNECTION"}]}`)
	client.APPSEC.On("GetRatePolicies", mock.Anything, appsec.GetRatePoliciesRequest{
		ConfigID: configID, ConfigVersion: version,
	}).Return(listResp, nil).Once()

	result, err := getRatePolicy(context.Background(), configID, version, policyID, newTestMeta(client))
	require.NoError(t, err)
	assert.Equal(t, "policy-A", result.Name)

	cached := &appsec.GetRatePoliciesResponse{}
	assert.NoError(t, cache.Get(cache.BucketName(SubproviderName), fmt.Sprintf("getRatePolicies:%d:%d", configID, version), cached))
	client.APPSEC.AssertExpectations(t)
}

func TestGetRatePolicy_CacheHit(t *testing.T) {
	clearCache()
	defer clearCache()

	client := edgegrid.NewTestClient()
	configID, version, policyID := int64(5002), int64(2), 22

	cached := buildRatePoliciesResponse(t, `{"ratePolicies":[{"id":22,"name":"cached-policy","type":"RATE"}]}`)
	require.NoError(t, cache.Set(cache.BucketName(SubproviderName), fmt.Sprintf("getRatePolicies:%d:%d", configID, version), cached))

	// No API mock registered — any API call would cause the test to fail
	result, err := getRatePolicy(context.Background(), configID, version, policyID, newTestMeta(client))
	require.NoError(t, err)
	assert.Equal(t, "cached-policy", result.Name)
	client.APPSEC.AssertExpectations(t)
}

func TestGetRatePolicy_CacheDisabled(t *testing.T) {
	clearCache()
	defer clearCache()
	cache.Enable(false)

	client := edgegrid.NewTestClient()
	configID, version, policyID := int64(5004), int64(4), 44

	client.APPSEC.On("GetRatePolicy", mock.Anything, appsec.GetRatePolicyRequest{
		ConfigID: configID, ConfigVersion: version, RatePolicyID: policyID,
	}).Return(&appsec.GetRatePolicyResponse{Name: "direct-policy", Type: "RATE"}, nil).Once()

	result, err := getRatePolicy(context.Background(), configID, version, policyID, newTestMeta(client))
	require.NoError(t, err)
	assert.Equal(t, "direct-policy", result.Name)
	client.APPSEC.AssertExpectations(t)
}

func TestGetRatePolicy_InvalidationForcesFreshFetch(t *testing.T) {
	clearCache()
	defer clearCache()

	client := edgegrid.NewTestClient()
	configID, version, policyID := int64(5006), int64(6), 66

	stale := buildRatePoliciesResponse(t, `{"ratePolicies":[{"id":66,"name":"stale-policy"}]}`)
	require.NoError(t, cache.Set(cache.BucketName(SubproviderName), fmt.Sprintf("getRatePolicies:%d:%d", configID, version), stale))

	invalidateRatePolicyCache(configID, version, akalog.NOPLogger())

	fresh := buildRatePoliciesResponse(t, `{"ratePolicies":[{"id":66,"name":"fresh-policy"}]}`)
	client.APPSEC.On("GetRatePolicies", mock.Anything, appsec.GetRatePoliciesRequest{
		ConfigID: configID, ConfigVersion: version,
	}).Return(fresh, nil).Once()

	result, err := getRatePolicy(context.Background(), configID, version, policyID, newTestMeta(client))
	require.NoError(t, err)
	assert.Equal(t, "fresh-policy", result.Name)
	client.APPSEC.AssertExpectations(t)
}
