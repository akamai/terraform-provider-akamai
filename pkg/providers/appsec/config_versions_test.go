// Package appsec provides unit tests for configuration version utility functions.
// This file contains comprehensive tests for the getLatestConfigVersion function,
// covering cache behavior, API interactions, error handling, and edge cases.
package appsec

import (
	"context"
	"errors"
	"testing"
	"time"

	"github.com/akamai/AkamaiOPEN-edgegrid-golang/v14/pkg/appsec"
	"github.com/akamai/terraform-provider-akamai/v11/internal/edgegrid"
	"github.com/akamai/terraform-provider-akamai/v11/pkg/cache"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/mock"
	"github.com/stretchr/testify/require"
)

// clearCache clears all cache entries for testing
func clearCache() {
	// Disable and re-enable cache to clear all entries
	cache.Enable(false)
	cache.Enable(true)
}

func TestGetLatestConfigVersion_CacheHit(t *testing.T) {
	// Clear cache before test
	clearCache()
	defer clearCache()

	// Setup
	client := edgegrid.NewTestClient()
	configID := 12345
	expectedVersion := 7

	configuration := &appsec.GetConfigurationResponse{
		ID:            configID,
		LatestVersion: expectedVersion,
	}
	cacheKey := latestVersionCacheKey(configID)
	err := cache.Set(cache.BucketName(SubproviderName), cacheKey, configuration)
	require.NoError(t, err)

	// No API calls should be made since we hit cache
	// Note: We don't add any expectations to the client mock

	ctx := context.Background()

	// Call the function under test
	result, err := getLatestConfigVersion(ctx, configID, client.APPSEC)

	// Assertions
	assert.NoError(t, err)
	assert.Equal(t, expectedVersion, result)

	// Verify no API calls were made
	client.APPSEC.AssertExpectations(t)
}

func TestGetLatestConfigVersion_CacheMiss_APISuccess(t *testing.T) {
	// Clear cache before test
	clearCache()
	defer clearCache()

	// Setup
	client := edgegrid.NewTestClient()
	configID := 12346
	expectedVersion := 9

	// Setup API mock
	getConfigResponse := appsec.GetConfigurationResponse{
		ID:            configID,
		LatestVersion: expectedVersion,
	}
	client.APPSEC.On("GetConfiguration",
		mock.Anything,
		appsec.GetConfigurationRequest{ConfigID: configID},
	).Return(&getConfigResponse, nil).Once()

	ctx := context.Background()

	// Call the function under test
	result, err := getLatestConfigVersion(ctx, configID, client.APPSEC)

	// Assertions
	assert.NoError(t, err)
	assert.Equal(t, expectedVersion, result)

	// Verify value was cached
	cachedConfig := &appsec.GetConfigurationResponse{}
	cacheKey := latestVersionCacheKey(configID)
	err = cache.Get(cache.BucketName(SubproviderName), cacheKey, cachedConfig)
	assert.NoError(t, err)
	assert.Equal(t, configID, cachedConfig.ID)
	assert.Equal(t, expectedVersion, cachedConfig.LatestVersion)

	// Verify API call was made
	client.APPSEC.AssertExpectations(t)
}

func TestGetLatestConfigVersion_APIError(t *testing.T) {
	// Clear cache before test
	clearCache()
	defer clearCache()

	// Setup
	client := edgegrid.NewTestClient()
	configID := 12347

	// Setup API mock to return error
	client.APPSEC.On("GetConfiguration",
		mock.Anything,
		appsec.GetConfigurationRequest{ConfigID: configID},
	).Return(nil, errors.New("API error")).Once()

	ctx := context.Background()

	// Call the function under test
	result, err := getLatestConfigVersion(ctx, configID, client.APPSEC)

	// Assertions
	assert.Error(t, err)
	assert.Contains(t, err.Error(), "API error")
	assert.Equal(t, 0, result)

	// Verify API call was made
	client.APPSEC.AssertExpectations(t)
}

func TestGetLatestConfigVersion_CacheDisabled(t *testing.T) {
	// Clear cache before test
	clearCache()
	defer clearCache()

	// Setup
	client := edgegrid.NewTestClient()
	configID := 12348
	expectedVersion := 5

	// Disable cache to test fallback behavior
	cache.Enable(false)

	// Setup API mock
	getConfigResponse := appsec.GetConfigurationResponse{
		ID:            configID,
		LatestVersion: expectedVersion,
	}
	client.APPSEC.On("GetConfiguration",
		mock.Anything,
		appsec.GetConfigurationRequest{ConfigID: configID},
	).Return(&getConfigResponse, nil).Once()

	ctx := context.Background()

	// Call the function under test
	result, err := getLatestConfigVersion(ctx, configID, client.APPSEC)

	// Assertions
	assert.NoError(t, err)
	assert.Equal(t, expectedVersion, result)

	// Verify API call was made
	client.APPSEC.AssertExpectations(t)
}

func TestGetLatestConfigVersion_DoubleCallCachingBehavior(t *testing.T) {
	// Clear cache before test
	clearCache()
	defer clearCache()

	// This test verifies that subsequent calls use the cache
	client := edgegrid.NewTestClient()
	configID := 12349
	expectedVersion := 10

	// Setup API mock - should only be called once
	getConfigResponse := appsec.GetConfigurationResponse{
		ID:            configID,
		LatestVersion: expectedVersion,
	}
	client.APPSEC.On("GetConfiguration",
		mock.Anything,
		appsec.GetConfigurationRequest{ConfigID: configID},
	).Return(&getConfigResponse, nil).Once()

	ctx := context.Background()

	// First call should hit API and cache the result
	result1, err1 := getLatestConfigVersion(ctx, configID, client.APPSEC)
	assert.NoError(t, err1)
	assert.Equal(t, expectedVersion, result1)

	// Second call should hit cache (no additional API call)
	result2, err2 := getLatestConfigVersion(ctx, configID, client.APPSEC)
	assert.NoError(t, err2)
	assert.Equal(t, expectedVersion, result2)

	// Verify API was called exactly once
	client.APPSEC.AssertExpectations(t)
}

func TestGetLatestConfigVersion_ConcurrentCalls(t *testing.T) {
	// Clear cache before test
	clearCache()
	defer clearCache()

	// This test verifies that the mutex prevents race conditions
	client := edgegrid.NewTestClient()
	configID := 12350
	expectedVersion := 15

	// Setup API mock - should be called only once due to mutex and caching
	getConfigResponse := appsec.GetConfigurationResponse{
		ID:            configID,
		LatestVersion: expectedVersion,
	}
	client.APPSEC.On("GetConfiguration",
		mock.Anything,
		appsec.GetConfigurationRequest{ConfigID: configID},
	).Return(&getConfigResponse, nil).Once()

	ctx := context.Background()

	// Run multiple goroutines to test concurrency
	results := make(chan int, 3)
	errs := make(chan error, 3)

	for i := 0; i < 3; i++ {
		go func() {
			result, err := getLatestConfigVersion(ctx, configID, client.APPSEC)
			results <- result
			errs <- err
		}()
	}

	// Collect results
	for i := 0; i < 3; i++ {
		result := <-results
		err := <-errs

		assert.NoError(t, err)
		assert.Equal(t, expectedVersion, result)
	}

	// Verify that the API was called at most once due to caching and mutex
	client.APPSEC.AssertExpectations(t)
}

func TestGetLatestConfigVersion_InvalidConfigID(t *testing.T) {
	// Clear cache before test
	clearCache()
	defer clearCache()

	// Test with an invalid config ID
	client := edgegrid.NewTestClient()
	configID := -1

	// Setup API mock to return error for invalid config ID
	client.APPSEC.On("GetConfiguration",
		mock.Anything,
		appsec.GetConfigurationRequest{ConfigID: configID},
	).Return(nil, errors.New("invalid config ID")).Once()

	ctx := context.Background()

	// Call the function under test
	result, err := getLatestConfigVersion(ctx, configID, client.APPSEC)

	// Assertions
	assert.Error(t, err)
	assert.Contains(t, err.Error(), "invalid config ID")
	assert.Equal(t, 0, result)

	// Verify API call was made
	client.APPSEC.AssertExpectations(t)
}

func TestGetActiveConfigVersions_SharedCacheSingleAPICall(t *testing.T) {
	clearCache()
	defer clearCache()

	client := edgegrid.NewTestClient()
	configID := 40002
	response := appsec.GetConfigurationResponse{
		ID: configID, LatestVersion: 5, StagingVersion: 3, ProductionVersion: 2,
	}

	// Expect ONE API call — getActiveConfigVersions reuses the cache (the fix)
	client.APPSEC.On("GetConfiguration", mock.Anything,
		appsec.GetConfigurationRequest{ConfigID: configID}).Return(&response, nil).Once()

	ctx := context.Background()
	_, _ = getLatestConfigVersion(ctx, configID, client.APPSEC)
	staging, production, err := getActiveConfigVersions(ctx, configID, client.APPSEC)

	assert.NoError(t, err)
	assert.Equal(t, 3, staging)
	assert.Equal(t, 2, production)
	client.APPSEC.AssertExpectations(t)
}

func TestGetActiveConfigVersions_CacheMiss_APISuccess(t *testing.T) {
	clearCache()
	defer clearCache()

	client := edgegrid.NewTestClient()
	configID := 40003
	response := appsec.GetConfigurationResponse{
		ID: configID, LatestVersion: 7, StagingVersion: 5, ProductionVersion: 4,
	}
	client.APPSEC.On("GetConfiguration", mock.Anything,
		appsec.GetConfigurationRequest{ConfigID: configID}).Return(&response, nil).Once()

	ctx := context.Background()
	staging, production, err := getActiveConfigVersions(ctx, configID, client.APPSEC)

	assert.NoError(t, err)
	assert.Equal(t, 5, staging)
	assert.Equal(t, 4, production)
	client.APPSEC.AssertExpectations(t)
}

func TestGetActiveConfigVersions_CacheHit_NoAPICall(t *testing.T) {
	clearCache()
	defer clearCache()

	configID := 40004
	cacheKey := latestVersionCacheKey(configID)
	configuration := &appsec.GetConfigurationResponse{
		ID: configID, LatestVersion: 7, StagingVersion: 5, ProductionVersion: 4,
	}
	require.NoError(t, cache.Set(cache.BucketName(SubproviderName), cacheKey, configuration))

	client := edgegrid.NewTestClient()

	ctx := context.Background()
	staging, production, err := getActiveConfigVersions(ctx, configID, client.APPSEC)

	assert.NoError(t, err)
	assert.Equal(t, 5, staging)
	assert.Equal(t, 4, production)
	client.APPSEC.AssertExpectations(t)
}

func TestGetActiveConfigVersions_APIError(t *testing.T) {
	clearCache()
	defer clearCache()

	client := edgegrid.NewTestClient()
	configID := 40005
	client.APPSEC.On("GetConfiguration", mock.Anything,
		appsec.GetConfigurationRequest{ConfigID: configID}).
		Return(nil, errors.New("API error")).Once()

	ctx := context.Background()
	staging, production, err := getActiveConfigVersions(ctx, configID, client.APPSEC)

	assert.Error(t, err)
	assert.Contains(t, err.Error(), "API error")
	assert.Equal(t, 0, staging)
	assert.Equal(t, 0, production)
	client.APPSEC.AssertExpectations(t)
}

func TestInvalidateConfigCache_ClearsKeysAndTriggersFetch(t *testing.T) {
	clearCache()
	defer clearCache()

	configID := 40007
	ctx := context.Background()

	// Pre-populate both cache keys
	latestKey := latestVersionCacheKey(configID)
	modifiableKey := modifiableVersionCacheKey(configID)
	configuration := &appsec.GetConfigurationResponse{
		ID: configID, LatestVersion: 5, StagingVersion: 3, ProductionVersion: 2,
	}
	require.NoError(t, cache.Set(cache.BucketName(SubproviderName), latestKey, configuration))
	require.NoError(t, cache.Set(cache.BucketName(SubproviderName), modifiableKey, configuration))

	// Invalidate
	invalidateConfigCache(configID)

	// Both keys must be gone
	require.ErrorIs(t, cache.Get(cache.BucketName(SubproviderName), latestKey, &appsec.GetConfigurationResponse{}), cache.ErrEntryNotFound)
	require.ErrorIs(t, cache.Get(cache.BucketName(SubproviderName), modifiableKey, &appsec.GetConfigurationResponse{}), cache.ErrEntryNotFound)

	// A subsequent getActiveConfigVersions call must re-hit the API
	client := edgegrid.NewTestClient()
	freshResponse := appsec.GetConfigurationResponse{
		ID: configID, LatestVersion: 6, StagingVersion: 4, ProductionVersion: 3,
	}
	client.APPSEC.On("GetConfiguration", mock.Anything,
		appsec.GetConfigurationRequest{ConfigID: configID}).Return(&freshResponse, nil).Once()

	staging, production, err := getActiveConfigVersions(ctx, configID, client.APPSEC)
	assert.NoError(t, err)
	assert.Equal(t, 4, staging)
	assert.Equal(t, 3, production)
	client.APPSEC.AssertExpectations(t)
}

func TestInvalidateConfigCache_EmptyCacheNoError(_ *testing.T) {
	clearCache()
	defer clearCache()

	// Invalidating when nothing is cached must not panic or error
	invalidateConfigCache(40008)
}

func TestGetActiveConfigVersions_CacheDisabled_FallsBackToAPI(t *testing.T) {
	clearCache()
	defer clearCache()

	cache.Enable(false)

	client := edgegrid.NewTestClient()
	configID := 40006
	response := appsec.GetConfigurationResponse{
		ID: configID, LatestVersion: 7, StagingVersion: 5, ProductionVersion: 4,
	}
	client.APPSEC.On("GetConfiguration", mock.Anything,
		appsec.GetConfigurationRequest{ConfigID: configID}).Return(&response, nil).Once()

	ctx := context.Background()
	staging, production, err := getActiveConfigVersions(ctx, configID, client.APPSEC)

	assert.NoError(t, err)
	assert.Equal(t, 5, staging)
	assert.Equal(t, 4, production)
	client.APPSEC.AssertExpectations(t)
}

// Tests for getModifiableConfigVersion function

func TestGetModifiableConfigVersion_CloneUpdatesSharedCache(t *testing.T) {
	clearCache()
	defer clearCache()

	configID := 29999
	latestVersion := 5
	stagingVersion := 5 // active in staging — forces clone
	productionVersion := 2
	newClonedVersion := 6
	resource := "test_resource"

	client := edgegrid.NewTestClient()

	getConfigResponse := appsec.GetConfigurationResponse{
		ID:                configID,
		LatestVersion:     latestVersion,
		StagingVersion:    stagingVersion,
		ProductionVersion: productionVersion,
	}
	client.APPSEC.On("GetConfiguration", mock.Anything,
		appsec.GetConfigurationRequest{ConfigID: configID}).Return(&getConfigResponse, nil).Once()

	cloneResponse := appsec.CreateConfigurationVersionCloneResponse{
		ConfigID: configID,
		Version:  newClonedVersion,
	}
	client.APPSEC.On("CreateConfigurationVersionClone", mock.Anything,
		appsec.CreateConfigurationVersionCloneRequest{ConfigID: configID, CreateFromVersion: latestVersion},
	).Return(&cloneResponse, nil).Once()

	ctx := context.Background()
	result, err := getModifiableConfigVersion(ctx, configID, resource, client.APPSEC)
	require.NoError(t, err)
	assert.Equal(t, newClonedVersion, result)

	// getLatestConfigVersion must return the cloned version without an extra API call
	latest, err := getLatestConfigVersion(ctx, configID, client.APPSEC)
	assert.NoError(t, err)
	assert.Equal(t, newClonedVersion, latest)

	// getActiveConfigVersions must also reuse the shared cache — staging/production unchanged by the clone
	staging, production, err := getActiveConfigVersions(ctx, configID, client.APPSEC)
	assert.NoError(t, err)
	assert.Equal(t, productionVersion, production)
	assert.Equal(t, stagingVersion, staging)

	// API was called exactly once — clone path updated the shared cache
	client.APPSEC.AssertExpectations(t)
}

func TestGetModifiableConfigVersion_CacheHit(t *testing.T) {
	// Clear cache before test
	clearCache()
	defer clearCache()

	configID := 22345
	expectedVersion := 3
	resource := "test_resource"

	// Pre-populate cache with modifiable key (the key getModifiableConfigVersion reads/writes)
	cacheKey := modifiableVersionCacheKey(configID)
	configuration := &appsec.GetConfigurationResponse{
		ID:            configID,
		LatestVersion: expectedVersion,
	}

	err := cache.Set(cache.BucketName(SubproviderName), cacheKey, configuration)
	require.NoError(t, err)

	// No client mocks needed since we should hit cache
	client := edgegrid.NewTestClient()

	ctx := context.Background()

	// Call the function under test
	result, err := getModifiableConfigVersion(ctx, configID, resource, client.APPSEC)

	// Assertions
	assert.NoError(t, err)
	assert.Equal(t, expectedVersion, result)

	// Verify no API calls were made due to cache hit
	client.APPSEC.AssertExpectations(t)
}

func TestGetModifiableConfigVersion_LatestVersionIsModifiable(t *testing.T) {
	// Clear cache before test
	clearCache()
	defer clearCache()

	configID := 23345
	latestVersion := 5
	stagingVersion := 3
	productionVersion := 2
	resource := "test_resource"

	// Setup mock client
	client := edgegrid.NewTestClient()

	getConfigResponse := appsec.GetConfigurationResponse{
		ID:                configID,
		LatestVersion:     latestVersion,
		StagingVersion:    stagingVersion,
		ProductionVersion: productionVersion,
	}

	// Mock the GetConfiguration call
	client.APPSEC.On("GetConfiguration",
		mock.Anything,
		appsec.GetConfigurationRequest{ConfigID: configID},
	).Return(&getConfigResponse, nil).Once()

	// Mock GetConfigurationVersion call for checkIfVersionWasPreviouslyActive
	getConfigVersionResponse := appsec.GetConfigurationVersionResponse{
		ConfigID:   configID,
		ConfigName: "Test Config",
		Version:    latestVersion,
		BasedOn:    1,
		CreateDate: time.Now(),
		CreatedBy:  "test@example.com",
		Production: appsec.EnvironmentStatus{
			Status: "Inactive",
			Time:   time.Now(),
		},
		Staging: appsec.EnvironmentStatus{
			Status: "Inactive",
			Time:   time.Now(),
		},
	}

	client.APPSEC.On("GetConfigurationVersion",
		mock.Anything,
		appsec.GetConfigurationVersionRequest{ConfigID: configID, Version: latestVersion},
	).Return(&getConfigVersionResponse, nil).Once()

	ctx := context.Background()

	// Call the function under test
	result, err := getModifiableConfigVersion(ctx, configID, resource, client.APPSEC)

	// Assertions
	assert.NoError(t, err)
	assert.Equal(t, latestVersion, result)

	client.APPSEC.AssertExpectations(t)
}

func TestGetModifiableConfigVersion_LatestVersionActiveInStaging(t *testing.T) {
	// Clear cache before test
	clearCache()
	defer clearCache()

	configID := 24345
	latestVersion := 5
	stagingVersion := 5 // Same as latest - active in staging
	productionVersion := 2
	newClonedVersion := 6
	resource := "test_resource"

	// Setup mock client
	client := edgegrid.NewTestClient()

	getConfigResponse := appsec.GetConfigurationResponse{
		ID:                configID,
		LatestVersion:     latestVersion,
		StagingVersion:    stagingVersion,
		ProductionVersion: productionVersion,
	}

	// Mock the GetConfiguration call
	client.APPSEC.On("GetConfiguration",
		mock.Anything,
		appsec.GetConfigurationRequest{ConfigID: configID},
	).Return(&getConfigResponse, nil).Once()

	// Mock the CreateConfigurationVersionClone call
	cloneResponse := appsec.CreateConfigurationVersionCloneResponse{
		ConfigID: configID,
		Version:  newClonedVersion,
	}

	client.APPSEC.On("CreateConfigurationVersionClone",
		mock.Anything,
		appsec.CreateConfigurationVersionCloneRequest{
			ConfigID:          configID,
			CreateFromVersion: latestVersion,
		},
	).Return(&cloneResponse, nil).Once()

	ctx := context.Background()

	// Call the function under test
	result, err := getModifiableConfigVersion(ctx, configID, resource, client.APPSEC)

	// Assertions
	assert.NoError(t, err)
	assert.Equal(t, newClonedVersion, result)

	client.APPSEC.AssertExpectations(t)
}

func TestGetModifiableConfigVersion_LatestVersionActiveInProduction(t *testing.T) {
	// Clear cache before test
	clearCache()
	defer clearCache()

	configID := 25345
	latestVersion := 5
	stagingVersion := 3
	productionVersion := 5 // Same as latest - active in production
	newClonedVersion := 6
	resource := "test_resource"

	// Setup mock client
	client := edgegrid.NewTestClient()

	getConfigResponse := appsec.GetConfigurationResponse{
		ID:                configID,
		LatestVersion:     latestVersion,
		StagingVersion:    stagingVersion,
		ProductionVersion: productionVersion,
	}

	// Mock the GetConfiguration call
	client.APPSEC.On("GetConfiguration",
		mock.Anything,
		appsec.GetConfigurationRequest{ConfigID: configID},
	).Return(&getConfigResponse, nil).Once()

	// Mock the CreateConfigurationVersionClone call
	cloneResponse := appsec.CreateConfigurationVersionCloneResponse{
		ConfigID: configID,
		Version:  newClonedVersion,
	}

	client.APPSEC.On("CreateConfigurationVersionClone",
		mock.Anything,
		appsec.CreateConfigurationVersionCloneRequest{
			ConfigID:          configID,
			CreateFromVersion: latestVersion,
		},
	).Return(&cloneResponse, nil).Once()

	ctx := context.Background()

	// Call the function under test
	result, err := getModifiableConfigVersion(ctx, configID, resource, client.APPSEC)

	// Assertions
	assert.NoError(t, err)
	assert.Equal(t, newClonedVersion, result)

	client.APPSEC.AssertExpectations(t)
}

func TestGetModifiableConfigVersion_LatestVersionWasPreviouslyActive(t *testing.T) {
	// Clear cache before test
	clearCache()
	defer clearCache()

	configID := 26345
	latestVersion := 5
	stagingVersion := 3
	productionVersion := 2
	newClonedVersion := 6
	resource := "test_resource"

	// Setup mock client
	client := edgegrid.NewTestClient()

	getConfigResponse := appsec.GetConfigurationResponse{
		ID:                configID,
		LatestVersion:     latestVersion,
		StagingVersion:    stagingVersion,
		ProductionVersion: productionVersion,
	}

	// Mock the GetConfiguration call
	client.APPSEC.On("GetConfiguration",
		mock.Anything,
		appsec.GetConfigurationRequest{ConfigID: configID},
	).Return(&getConfigResponse, nil).Once()

	// Mock GetConfigurationVersion call for checkIfVersionWasPreviouslyActive
	// This simulates a version that was previously active but is now deactivated
	getConfigVersionResponse := appsec.GetConfigurationVersionResponse{
		ConfigID:   configID,
		ConfigName: "Test Config",
		Version:    latestVersion,
		BasedOn:    1,
		CreateDate: time.Now(),
		CreatedBy:  "test@example.com",
		Production: appsec.EnvironmentStatus{
			Status: "Inactive",
			Time:   time.Now(),
		},
		Staging: appsec.EnvironmentStatus{
			Status: "Deactivated", // Was previously active
			Time:   time.Now(),
		},
	}

	client.APPSEC.On("GetConfigurationVersion",
		mock.Anything,
		appsec.GetConfigurationVersionRequest{ConfigID: configID, Version: latestVersion},
	).Return(&getConfigVersionResponse, nil).Once()

	// Mock the CreateConfigurationVersionClone call
	cloneResponse := appsec.CreateConfigurationVersionCloneResponse{
		ConfigID: configID,
		Version:  newClonedVersion,
	}

	client.APPSEC.On("CreateConfigurationVersionClone",
		mock.Anything,
		appsec.CreateConfigurationVersionCloneRequest{
			ConfigID:          configID,
			CreateFromVersion: latestVersion,
		},
	).Return(&cloneResponse, nil).Once()

	ctx := context.Background()

	// Call the function under test
	result, err := getModifiableConfigVersion(ctx, configID, resource, client.APPSEC)

	// Assertions
	assert.NoError(t, err)
	assert.Equal(t, newClonedVersion, result)

	client.APPSEC.AssertExpectations(t)
}

func TestGetModifiableConfigVersion_GetConfigurationError(t *testing.T) {
	// Clear cache before test
	clearCache()
	defer clearCache()

	configID := 27345
	resource := "test_resource"
	expectedError := errors.New("API error")

	// Setup mock client
	client := edgegrid.NewTestClient()

	// Mock the GetConfiguration call to return error
	client.APPSEC.On("GetConfiguration",
		mock.Anything,
		appsec.GetConfigurationRequest{ConfigID: configID},
	).Return(nil, expectedError).Once()

	ctx := context.Background()

	// Call the function under test
	result, err := getModifiableConfigVersion(ctx, configID, resource, client.APPSEC)

	// Assertions
	assert.Error(t, err)
	assert.Contains(t, err.Error(), "API error")
	assert.Equal(t, 0, result)

	client.APPSEC.AssertExpectations(t)
}

func TestGetModifiableConfigVersion_CloneError(t *testing.T) {
	// Clear cache before test
	clearCache()
	defer clearCache()

	configID := 28345
	latestVersion := 5
	stagingVersion := 5 // Same as latest - active in staging
	productionVersion := 2
	resource := "test_resource"
	expectedError := errors.New("clone error")

	// Setup mock client
	client := edgegrid.NewTestClient()

	getConfigResponse := appsec.GetConfigurationResponse{
		ID:                configID,
		LatestVersion:     latestVersion,
		StagingVersion:    stagingVersion,
		ProductionVersion: productionVersion,
	}

	// Mock the GetConfiguration call
	client.APPSEC.On("GetConfiguration",
		mock.Anything,
		appsec.GetConfigurationRequest{ConfigID: configID},
	).Return(&getConfigResponse, nil).Once()

	// Mock the CreateConfigurationVersionClone call to return error
	client.APPSEC.On("CreateConfigurationVersionClone",
		mock.Anything,
		appsec.CreateConfigurationVersionCloneRequest{
			ConfigID:          configID,
			CreateFromVersion: latestVersion,
		},
	).Return(nil, expectedError).Once()

	ctx := context.Background()

	// Call the function under test
	result, err := getModifiableConfigVersion(ctx, configID, resource, client.APPSEC)

	// Assertions
	assert.Error(t, err)
	assert.Contains(t, err.Error(), "clone error")
	assert.Equal(t, 0, result)

	client.APPSEC.AssertExpectations(t)
}

func TestGetModifiableConfigVersion_ConcurrentAccess(t *testing.T) {
	// Clear cache before test
	clearCache()
	defer clearCache()

	configID := 29345
	latestVersion := 5
	stagingVersion := 3
	productionVersion := 2
	resource := "test_resource"

	// Setup mock client
	client := edgegrid.NewTestClient()

	getConfigResponse := appsec.GetConfigurationResponse{
		ID:                configID,
		LatestVersion:     latestVersion,
		StagingVersion:    stagingVersion,
		ProductionVersion: productionVersion,
	}

	// Mock the GetConfiguration call
	client.APPSEC.On("GetConfiguration",
		mock.Anything,
		appsec.GetConfigurationRequest{ConfigID: configID},
	).Return(&getConfigResponse, nil).Once()

	// Mock GetConfigurationVersion call for checkIfVersionWasPreviouslyActive
	getConfigVersionResponse := appsec.GetConfigurationVersionResponse{
		ConfigID:   configID,
		ConfigName: "Test Config",
		Version:    latestVersion,
		BasedOn:    1,
		CreateDate: time.Now(),
		CreatedBy:  "test@example.com",
		Production: appsec.EnvironmentStatus{
			Status: "Inactive",
			Time:   time.Now(),
		},
		Staging: appsec.EnvironmentStatus{
			Status: "Inactive",
			Time:   time.Now(),
		},
	}

	client.APPSEC.On("GetConfigurationVersion",
		mock.Anything,
		appsec.GetConfigurationVersionRequest{ConfigID: configID, Version: latestVersion},
	).Return(&getConfigVersionResponse, nil).Once()

	ctx := context.Background()

	// Run multiple goroutines to test concurrency
	results := make(chan int, 3)
	errs := make(chan error, 3)

	for i := 0; i < 3; i++ {
		go func() {
			result, err := getModifiableConfigVersion(ctx, configID, resource, client.APPSEC)
			results <- result
			errs <- err
		}()
	}

	// Collect results
	for i := 0; i < 3; i++ {
		result := <-results
		err := <-errs
		assert.NoError(t, err)
		assert.Equal(t, latestVersion, result)
	}

	client.APPSEC.AssertExpectations(t)
}

func TestGetModifiableConfigVersion_GetConfigurationVersionError(t *testing.T) {
	// Clear cache before test
	clearCache()
	defer clearCache()

	configID := 30345
	latestVersion := 5
	stagingVersion := 3
	productionVersion := 2
	resource := "test_resource"

	// Setup mock client
	client := edgegrid.NewTestClient()

	getConfigResponse := appsec.GetConfigurationResponse{
		ID:                configID,
		LatestVersion:     latestVersion,
		StagingVersion:    stagingVersion,
		ProductionVersion: productionVersion,
	}

	// Mock the GetConfiguration call
	client.APPSEC.On("GetConfiguration",
		mock.Anything,
		appsec.GetConfigurationRequest{ConfigID: configID},
	).Return(&getConfigResponse, nil).Once()

	// Mock GetConfigurationVersion call to return error
	// In this case, the function should still work and return the latest version
	// because checkIfVersionWasPreviouslyActive handles errors gracefully
	client.APPSEC.On("GetConfigurationVersion",
		mock.Anything,
		appsec.GetConfigurationVersionRequest{ConfigID: configID, Version: latestVersion},
	).Return(nil, errors.New("version API error")).Once()

	ctx := context.Background()

	// Call the function under test
	result, err := getModifiableConfigVersion(ctx, configID, resource, client.APPSEC)

	// Assertions - should still succeed and return latest version
	// because checkIfVersionWasPreviouslyActive returns false on error
	assert.NoError(t, err)
	assert.Equal(t, latestVersion, result)

	client.APPSEC.AssertExpectations(t)
}

// TestCacheKeyDivergence_InvalidationClearsBothKeys verifies that invalidateConfigCache
// deletes both the modifiable and latest version keys, so stale diverged state cannot
// persist after an activation or deactivation completes.
func TestCacheKeyDivergence_InvalidationClearsBothKeys(t *testing.T) {
	clearCache()
	defer clearCache()

	configID := 99001
	staleModifiable := 5
	staleLatest := 3 // deliberately different — simulates diverged cache state

	// Seed both keys with stale, diverged values.
	staleModifiableConfig := &appsec.GetConfigurationResponse{ID: configID, LatestVersion: staleModifiable}
	staleLatestConfig := &appsec.GetConfigurationResponse{ID: configID, LatestVersion: staleLatest}
	require.NoError(t, cache.Set(cache.BucketName(SubproviderName), modifiableVersionCacheKey(configID), staleModifiableConfig))
	require.NoError(t, cache.Set(cache.BucketName(SubproviderName), latestVersionCacheKey(configID), staleLatestConfig))

	// Simulate activation completing.
	invalidateConfigCache(configID)

	// Both keys must now be absent — fresh fetches are required.
	out := &appsec.GetConfigurationResponse{}
	assert.ErrorIs(t, cache.Get(cache.BucketName(SubproviderName), modifiableVersionCacheKey(configID), out),
		cache.ErrEntryNotFound, "modifiable key should be invalidated")
	assert.ErrorIs(t, cache.Get(cache.BucketName(SubproviderName), latestVersionCacheKey(configID), out),
		cache.ErrEntryNotFound, "latest key should be invalidated")
}

// TestGetActiveConfigVersions_SharedCacheReverseDirection verifies that when
// getActiveConfigVersions populates latestVersionCacheKey first, a subsequent
// getLatestConfigVersion call reuses the cache without an extra API call.
func TestGetActiveConfigVersions_SharedCacheReverseDirection(t *testing.T) {
	clearCache()
	defer clearCache()

	client := edgegrid.NewTestClient()
	configID := 99002
	response := appsec.GetConfigurationResponse{
		ID: configID, LatestVersion: 8, StagingVersion: 6, ProductionVersion: 5,
	}

	// Expect ONE API call — getLatestConfigVersion reuses the cache populated by getActiveConfigVersions.
	client.APPSEC.On("GetConfiguration", mock.Anything,
		appsec.GetConfigurationRequest{ConfigID: configID}).Return(&response, nil).Once()

	ctx := context.Background()
	_, _, err := getActiveConfigVersions(ctx, configID, client.APPSEC)
	require.NoError(t, err)

	latest, err := getLatestConfigVersion(ctx, configID, client.APPSEC)
	assert.NoError(t, err)
	assert.Equal(t, 8, latest)
	client.APPSEC.AssertExpectations(t)
}

// TestInvalidateCacheForcesFreshFetch verifies that after invalidateConfigCache is called
// (e.g. following an activation), both functions re-fetch from the API instead of
// returning the stale cached value.
func TestInvalidateCacheForcesFreshFetch(t *testing.T) {
	clearCache()
	defer clearCache()

	configID := 99003
	cachedVersion := 5
	freshVersion := 6

	// Pre-populate both cache keys with a stale version.
	staleConfig := &appsec.GetConfigurationResponse{
		ID:            configID,
		LatestVersion: cachedVersion,
	}
	err := cache.Set(cache.BucketName(SubproviderName), modifiableVersionCacheKey(configID), staleConfig)
	require.NoError(t, err)
	err = cache.Set(cache.BucketName(SubproviderName), latestVersionCacheKey(configID), staleConfig)
	require.NoError(t, err)

	// Simulate an activation completing — this is what performActivation calls.
	invalidateConfigCache(configID)

	// Both functions should now miss the cache and re-fetch from the API.
	client := edgegrid.NewTestClient()
	freshConfig := appsec.GetConfigurationResponse{
		ID:            configID,
		LatestVersion: freshVersion,
	}
	getConfigVersionResponse := appsec.GetConfigurationVersionResponse{
		ConfigID:   configID,
		Version:    freshVersion,
		Staging:    appsec.EnvironmentStatus{Status: "Inactive"},
		Production: appsec.EnvironmentStatus{Status: "Inactive"},
	}
	client.APPSEC.On("GetConfiguration",
		mock.Anything,
		appsec.GetConfigurationRequest{ConfigID: configID},
	).Return(&freshConfig, nil).Times(2)
	client.APPSEC.On("GetConfigurationVersion",
		mock.Anything,
		appsec.GetConfigurationVersionRequest{ConfigID: configID, Version: freshVersion},
	).Return(&getConfigVersionResponse, nil).Once()

	ctx := context.Background()

	modifiable, err := getModifiableConfigVersion(ctx, configID, "test", client.APPSEC)
	require.NoError(t, err)
	assert.Equal(t, freshVersion, modifiable, "getModifiableConfigVersion should return fresh version after invalidation")

	// Invalidate again to force getLatestConfigVersion to also re-fetch
	// (getModifiableConfigVersion re-populated the cache with freshVersion).
	invalidateConfigCache(configID)

	latest, err := getLatestConfigVersion(ctx, configID, client.APPSEC)
	require.NoError(t, err)
	assert.Equal(t, freshVersion, latest, "getLatestConfigVersion should return fresh version after invalidation")

	client.APPSEC.AssertExpectations(t)
}
