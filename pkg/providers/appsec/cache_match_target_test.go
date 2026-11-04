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

func TestFilterMatchTarget(t *testing.T) {
	logger := akalog.NOPLogger()

	targets := buildMatchTargetsResponse(t, `{
		"matchTargets": {
			"apiTargets": [
				{"targetId": 101, "type": "api", "sequence": 3}
			],
			"websiteTargets": [
				{"targetId": 202, "type": "website", "sequence": 2, "hostnames": ["example.com"]}
			]
		}
	}`)

	// sequence is part of the match_target JSON written to state, so it must survive the conversion
	// of a list item into the single-GET response for both API and website targets.
	t.Run("found_in_api_targets", func(t *testing.T) {
		result, err := filterMatchTarget(targets, 101, logger)
		require.NoError(t, err)
		require.NotNil(t, result)
		assert.Equal(t, "api", result.Type)
		assert.Equal(t, 3, result.Sequence)
	})

	t.Run("found_in_website_targets", func(t *testing.T) {
		result, err := filterMatchTarget(targets, 202, logger)
		require.NoError(t, err)
		require.NotNil(t, result)
		assert.Equal(t, "website", result.Type)
		assert.Equal(t, 2, result.Sequence)
	})

	t.Run("not_found_in_either", func(t *testing.T) {
		result, err := filterMatchTarget(targets, 999, logger)
		assert.Nil(t, result)
		assert.ErrorIs(t, err, ErrMatchTargetNotFound)
	})
}

func TestGetMatchTarget_CacheMiss_APITarget(t *testing.T) {
	clearCache()
	defer clearCache()

	client := edgegrid.NewTestClient()
	configID, version, targetID := int64(6001), int64(1), 101

	listResp := buildMatchTargetsResponse(t, `{
		"matchTargets": {
			"apiTargets": [{"targetId":101,"type":"api","sequence":1}],
			"websiteTargets": []
		}
	}`)
	client.APPSEC.On("GetMatchTargets", mock.Anything, appsec.GetMatchTargetsRequest{
		ConfigID: configID, ConfigVersion: version,
	}).Return(listResp, nil).Once()

	result, err := getMatchTarget(context.Background(), configID, version, targetID, newTestMeta(client))
	require.NoError(t, err)
	assert.Equal(t, "api", result.Type)

	cached := &appsec.GetMatchTargetsResponse{}
	assert.NoError(t, cache.Get(cache.BucketName(SubproviderName), fmt.Sprintf("getMatchTargets:%d:%d", configID, version), cached))
	client.APPSEC.AssertExpectations(t)
}

func TestGetMatchTarget_CacheMiss_WebsiteTarget(t *testing.T) {
	clearCache()
	defer clearCache()

	client := edgegrid.NewTestClient()
	configID, version, targetID := int64(6002), int64(2), 202

	listResp := buildMatchTargetsResponse(t, `{
		"matchTargets": {
			"apiTargets": [],
			"websiteTargets": [{"targetId":202,"type":"website","hostnames":["example.com"]}]
		}
	}`)
	client.APPSEC.On("GetMatchTargets", mock.Anything, appsec.GetMatchTargetsRequest{
		ConfigID: configID, ConfigVersion: version,
	}).Return(listResp, nil).Once()

	result, err := getMatchTarget(context.Background(), configID, version, targetID, newTestMeta(client))
	require.NoError(t, err)
	assert.Equal(t, "website", result.Type)
	client.APPSEC.AssertExpectations(t)
}

func TestGetMatchTarget_CacheHit(t *testing.T) {
	clearCache()
	defer clearCache()

	client := edgegrid.NewTestClient()
	configID, version, targetID := int64(6003), int64(3), 303

	cached := buildMatchTargetsResponse(t, `{
		"matchTargets": {
			"apiTargets": [{"targetId":303,"type":"api"}],
			"websiteTargets": []
		}
	}`)
	require.NoError(t, cache.Set(cache.BucketName(SubproviderName), fmt.Sprintf("getMatchTargets:%d:%d", configID, version), cached))

	result, err := getMatchTarget(context.Background(), configID, version, targetID, newTestMeta(client))
	require.NoError(t, err)
	assert.Equal(t, "api", result.Type)
	client.APPSEC.AssertExpectations(t)
}

func TestGetMatchTarget_CacheDisabled(t *testing.T) {
	clearCache()
	defer clearCache()
	cache.Enable(false)

	client := edgegrid.NewTestClient()
	configID, version, targetID := int64(6004), int64(4), 404

	client.APPSEC.On("GetMatchTarget", mock.Anything, appsec.GetMatchTargetRequest{
		ConfigID: configID, ConfigVersion: version, TargetID: targetID,
	}).Return(&appsec.GetMatchTargetResponse{Type: "website"}, nil).Once()

	result, err := getMatchTarget(context.Background(), configID, version, targetID, newTestMeta(client))
	require.NoError(t, err)
	assert.Equal(t, "website", result.Type)
	client.APPSEC.AssertExpectations(t)
}

func TestGetMatchTarget_InvalidationForcesFreshFetch(t *testing.T) {
	clearCache()
	defer clearCache()

	client := edgegrid.NewTestClient()
	configID, version, targetID := int64(6006), int64(6), 606

	stale := buildMatchTargetsResponse(t, `{"matchTargets":{"apiTargets":[{"targetId":606,"type":"api","apis":[{"id":1,"name":"stale-api"}]}],"websiteTargets":[]}}`)
	require.NoError(t, cache.Set(cache.BucketName(SubproviderName), fmt.Sprintf("getMatchTargets:%d:%d", configID, version), stale))

	invalidateMatchTargetCache(configID, version, akalog.NOPLogger())

	fresh := buildMatchTargetsResponse(t, `{"matchTargets":{"apiTargets":[{"targetId":606,"type":"api","apis":[{"id":2,"name":"fresh-api"}]}],"websiteTargets":[]}}`)
	client.APPSEC.On("GetMatchTargets", mock.Anything, appsec.GetMatchTargetsRequest{
		ConfigID: configID, ConfigVersion: version,
	}).Return(fresh, nil).Once()

	result, err := getMatchTarget(context.Background(), configID, version, targetID, newTestMeta(client))
	require.NoError(t, err)
	assert.Equal(t, "api", result.Type)
	require.Len(t, result.Apis, 1)
	assert.Equal(t, "fresh-api", result.Apis[0].Name)
	client.APPSEC.AssertExpectations(t)
}
