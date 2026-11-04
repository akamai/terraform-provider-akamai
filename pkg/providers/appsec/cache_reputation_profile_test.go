package appsec

import (
	"context"
	"fmt"
	"testing"

	"github.com/akamai/AkamaiOPEN-edgegrid-golang/v15/pkg/appsec"
	akalog "github.com/akamai/AkamaiOPEN-edgegrid-golang/v15/pkg/log"
	"github.com/akamai/terraform-provider-akamai/v12/internal/edgegrid"
	"github.com/akamai/terraform-provider-akamai/v12/pkg/cache"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/mock"
	"github.com/stretchr/testify/require"
)

func TestFilterReputationProfile(t *testing.T) {
	logger := akalog.NOPLogger()

	profiles := buildReputationProfilesResponse(t, `{
		"reputationProfiles": [
			{"id": 31, "name": "profile-X", "context": "DOSATCK", "threshold": 5},
			{"id": 32, "name": "profile-Y", "context": "WEBATCK", "threshold": 3}
		]
	}`)

	t.Run("found", func(t *testing.T) {
		result, err := filterReputationProfile(profiles, 31, logger)
		require.NoError(t, err)
		require.NotNil(t, result)
		assert.Equal(t, "profile-X", result.Name)
		assert.Equal(t, "DOSATCK", result.Context)
		assert.Equal(t, 5, result.Threshold)
	})

	t.Run("not_found", func(t *testing.T) {
		result, err := filterReputationProfile(profiles, 99, logger)
		assert.Nil(t, result)
		assert.ErrorIs(t, err, ErrReputationProfileNotFound)
	})
}

func TestGetReputationProfile_CacheMiss(t *testing.T) {
	clearCache()
	defer clearCache()

	client := edgegrid.NewTestClient()
	configID, version, profileID := int64(7001), int64(1), int64(31)

	listResp := buildReputationProfilesResponse(t, `{"reputationProfiles":[{"id":31,"name":"profile-X","context":"DOSATCK","threshold":5}]}`)
	client.APPSEC.On("GetReputationProfiles", mock.Anything, appsec.GetReputationProfilesRequest{
		ConfigID: configID, ConfigVersion: version,
	}).Return(listResp, nil).Once()

	result, err := getReputationProfile(context.Background(), configID, version, profileID, newTestMeta(client))
	require.NoError(t, err)
	assert.Equal(t, "profile-X", result.Name)
	assert.Equal(t, "DOSATCK", result.Context)

	cached := &appsec.GetReputationProfilesResponse{}
	assert.NoError(t, cache.Get(cache.BucketName(SubproviderName), fmt.Sprintf("getReputationProfiles:%d:%d", configID, version), cached))
	client.APPSEC.AssertExpectations(t)
}

func TestGetReputationProfile_CacheHit(t *testing.T) {
	clearCache()
	defer clearCache()

	client := edgegrid.NewTestClient()
	configID, version, profileID := int64(7002), int64(2), int64(32)

	cached := buildReputationProfilesResponse(t, `{"reputationProfiles":[{"id":32,"name":"cached-profile","context":"WEBATCK","threshold":3}]}`)
	require.NoError(t, cache.Set(cache.BucketName(SubproviderName), fmt.Sprintf("getReputationProfiles:%d:%d", configID, version), cached))

	result, err := getReputationProfile(context.Background(), configID, version, profileID, newTestMeta(client))
	require.NoError(t, err)
	assert.Equal(t, "cached-profile", result.Name)
	client.APPSEC.AssertExpectations(t)
}

func TestGetReputationProfile_CacheDisabled(t *testing.T) {
	clearCache()
	defer clearCache()
	cache.Enable(false)

	client := edgegrid.NewTestClient()
	configID, version, profileID := int64(7003), int64(3), int64(33)

	client.APPSEC.On("GetReputationProfile", mock.Anything, appsec.GetReputationProfileRequest{
		ConfigID: configID, ConfigVersion: version, ReputationProfileID: profileID,
	}).Return(&appsec.GetReputationProfileResponse{Name: "direct-profile", Context: "DOSATCK"}, nil).Once()

	result, err := getReputationProfile(context.Background(), configID, version, profileID, newTestMeta(client))
	require.NoError(t, err)
	assert.Equal(t, "direct-profile", result.Name)
	client.APPSEC.AssertExpectations(t)
}

func TestGetReputationProfile_InvalidationForcesFreshFetch(t *testing.T) {
	clearCache()
	defer clearCache()

	client := edgegrid.NewTestClient()
	configID, version, profileID := int64(7005), int64(5), int64(55)

	stale := buildReputationProfilesResponse(t, `{"reputationProfiles":[{"id":55,"name":"stale-profile","context":"DOSATCK","threshold":2}]}`)
	require.NoError(t, cache.Set(cache.BucketName(SubproviderName), fmt.Sprintf("getReputationProfiles:%d:%d", configID, version), stale))

	invalidateReputationProfileCache(configID, version, akalog.NOPLogger())

	fresh := buildReputationProfilesResponse(t, `{"reputationProfiles":[{"id":55,"name":"fresh-profile","context":"DOSATCK","threshold":10}]}`)
	client.APPSEC.On("GetReputationProfiles", mock.Anything, appsec.GetReputationProfilesRequest{
		ConfigID: configID, ConfigVersion: version,
	}).Return(fresh, nil).Once()

	result, err := getReputationProfile(context.Background(), configID, version, profileID, newTestMeta(client))
	require.NoError(t, err)
	assert.Equal(t, "fresh-profile", result.Name)
	assert.Equal(t, 10, result.Threshold)
	client.APPSEC.AssertExpectations(t)
}
