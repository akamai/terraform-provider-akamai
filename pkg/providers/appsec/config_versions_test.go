// Package appsec provides unit tests for configuration version utility functions.
// This file contains comprehensive tests for the getLatestConfigVersion function,
// covering cache behavior, API interactions, error handling, and edge cases.
package appsec

import (
	"context"
	"errors"
	"testing"
	"time"

	"github.com/akamai/AkamaiOPEN-edgegrid-golang/v12/pkg/appsec"
	akalog "github.com/akamai/AkamaiOPEN-edgegrid-golang/v12/pkg/log"
	"github.com/akamai/AkamaiOPEN-edgegrid-golang/v12/pkg/session"
	"github.com/akamai/terraform-provider-akamai/v9/pkg/cache"
	"github.com/akamai/terraform-provider-akamai/v9/pkg/log"
	akameta "github.com/akamai/terraform-provider-akamai/v9/pkg/meta"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/mock"
	"github.com/stretchr/testify/require"
)

// Simple mock implementation of the Meta interface for testing
type mockMeta struct{}

func (m *mockMeta) Log(args ...interface{}) akalog.Interface {
	// Return a logger that does nothing for testing
	return log.Get(args...)
}

func (m *mockMeta) OperationID() string {
	return "test-operation-id"
}

func (m *mockMeta) Session() session.Session {
	return nil
}

// Verify our mock implements the interface
var _ akameta.Meta = (*mockMeta)(nil)

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
	client := &appsec.Mock{}
	configID := 12345
	expectedVersion := 7

	configuration := &appsec.GetConfigurationResponse{
		ID:            configID,
		LatestVersion: expectedVersion,
	}
	cacheKey := "getLatestConfigVersion:12345"
	err := cache.Set(cache.BucketName(SubproviderName), cacheKey, configuration)
	require.NoError(t, err)

	// No API calls should be made since we hit cache
	// Note: We don't add any expectations to the client mock

	ctx := context.Background()

	useClient(client, func() {
		// Call the function under test
		result, err := getLatestConfigVersion(ctx, configID, &mockMeta{})

		// Assertions
		assert.NoError(t, err)
		assert.Equal(t, expectedVersion, result)
	})

	// Verify no API calls were made
	client.AssertExpectations(t)
}

func TestGetLatestConfigVersion_CacheMiss_APISuccess(t *testing.T) {
	// Clear cache before test
	clearCache()
	defer clearCache()

	// Setup
	client := &appsec.Mock{}
	configID := 12346
	expectedVersion := 9

	// Setup API mock
	getConfigResponse := appsec.GetConfigurationResponse{
		ID:            configID,
		LatestVersion: expectedVersion,
	}
	client.On("GetConfiguration",
		mock.Anything,
		appsec.GetConfigurationRequest{ConfigID: configID},
	).Return(&getConfigResponse, nil).Once()

	ctx := context.Background()

	useClient(client, func() {
		// Call the function under test
		result, err := getLatestConfigVersion(ctx, configID, &mockMeta{})

		// Assertions
		assert.NoError(t, err)
		assert.Equal(t, expectedVersion, result)

		// Verify value was cached
		cachedConfig := &appsec.GetConfigurationResponse{}
		cacheKey := "getLatestConfigVersion:12346"
		err = cache.Get(cache.BucketName(SubproviderName), cacheKey, cachedConfig)
		assert.NoError(t, err)
		assert.Equal(t, configID, cachedConfig.ID)
		assert.Equal(t, expectedVersion, cachedConfig.LatestVersion)
	})

	// Verify API call was made
	client.AssertExpectations(t)
}

func TestGetLatestConfigVersion_APIError(t *testing.T) {
	// Clear cache before test
	clearCache()
	defer clearCache()

	// Setup
	client := &appsec.Mock{}
	configID := 12347

	// Setup API mock to return error
	client.On("GetConfiguration",
		mock.Anything,
		appsec.GetConfigurationRequest{ConfigID: configID},
	).Return(nil, errors.New("API error")).Once()

	ctx := context.Background()

	useClient(client, func() {
		// Call the function under test
		result, err := getLatestConfigVersion(ctx, configID, &mockMeta{})

		// Assertions
		assert.Error(t, err)
		assert.Contains(t, err.Error(), "API error")
		assert.Equal(t, 0, result)
	})

	// Verify API call was made
	client.AssertExpectations(t)
}

func TestGetLatestConfigVersion_CacheDisabled(t *testing.T) {
	// Clear cache before test
	clearCache()
	defer clearCache()

	// Setup
	client := &appsec.Mock{}
	configID := 12348
	expectedVersion := 5

	// Disable cache to test fallback behavior
	cache.Enable(false)

	// Setup API mock
	getConfigResponse := appsec.GetConfigurationResponse{
		ID:            configID,
		LatestVersion: expectedVersion,
	}
	client.On("GetConfiguration",
		mock.Anything,
		appsec.GetConfigurationRequest{ConfigID: configID},
	).Return(&getConfigResponse, nil).Once()

	ctx := context.Background()

	useClient(client, func() {
		// Call the function under test
		result, err := getLatestConfigVersion(ctx, configID, &mockMeta{})

		// Assertions
		assert.NoError(t, err)
		assert.Equal(t, expectedVersion, result)
	})

	// Verify API call was made
	client.AssertExpectations(t)
}

func TestGetLatestConfigVersion_DoubleCallCachingBehavior(t *testing.T) {
	// Clear cache before test
	clearCache()
	defer clearCache()

	// This test verifies that subsequent calls use the cache
	client := &appsec.Mock{}
	configID := 12349
	expectedVersion := 10

	// Setup API mock - should only be called once
	getConfigResponse := appsec.GetConfigurationResponse{
		ID:            configID,
		LatestVersion: expectedVersion,
	}
	client.On("GetConfiguration",
		mock.Anything,
		appsec.GetConfigurationRequest{ConfigID: configID},
	).Return(&getConfigResponse, nil).Once()

	ctx := context.Background()

	useClient(client, func() {
		// First call should hit API and cache the result
		result1, err1 := getLatestConfigVersion(ctx, configID, &mockMeta{})
		assert.NoError(t, err1)
		assert.Equal(t, expectedVersion, result1)

		// Second call should hit cache (no additional API call)
		result2, err2 := getLatestConfigVersion(ctx, configID, &mockMeta{})
		assert.NoError(t, err2)
		assert.Equal(t, expectedVersion, result2)
	})

	// Verify API was called exactly once
	client.AssertExpectations(t)
}

func TestGetLatestConfigVersion_ConcurrentCalls(t *testing.T) {
	// Clear cache before test
	clearCache()
	defer clearCache()

	// This test verifies that the mutex prevents race conditions
	client := &appsec.Mock{}
	configID := 12350
	expectedVersion := 15

	// Setup API mock - should be called only once due to mutex and caching
	getConfigResponse := appsec.GetConfigurationResponse{
		ID:            configID,
		LatestVersion: expectedVersion,
	}
	client.On("GetConfiguration",
		mock.Anything,
		appsec.GetConfigurationRequest{ConfigID: configID},
	).Return(&getConfigResponse, nil).Once()

	ctx := context.Background()

	useClient(client, func() {
		// Run multiple goroutines to test concurrency
		results := make(chan int, 3)
		errors := make(chan error, 3)

		for i := 0; i < 3; i++ {
			go func() {
				result, err := getLatestConfigVersion(ctx, configID, &mockMeta{})
				results <- result
				errors <- err
			}()
		}

		// Collect results
		for i := 0; i < 3; i++ {
			result := <-results
			err := <-errors

			assert.NoError(t, err)
			assert.Equal(t, expectedVersion, result)
		}
	})

	// Verify that the API was called at most once due to caching and mutex
	client.AssertExpectations(t)
}

func TestGetLatestConfigVersion_InvalidConfigID(t *testing.T) {
	// Clear cache before test
	clearCache()
	defer clearCache()

	// Test with an invalid config ID
	client := &appsec.Mock{}
	configID := -1

	// Setup API mock to return error for invalid config ID
	client.On("GetConfiguration",
		mock.Anything,
		appsec.GetConfigurationRequest{ConfigID: configID},
	).Return(nil, errors.New("invalid config ID")).Once()

	ctx := context.Background()

	useClient(client, func() {
		// Call the function under test
		result, err := getLatestConfigVersion(ctx, configID, &mockMeta{})

		// Assertions
		assert.Error(t, err)
		assert.Contains(t, err.Error(), "invalid config ID")
		assert.Equal(t, 0, result)
	})

	// Verify API call was made
	client.AssertExpectations(t)
}

// Tests for getModifiableConfigVersion function

func TestGetModifiableConfigVersion_CacheHit(t *testing.T) {
	// Clear cache before test
	clearCache()
	defer clearCache()

	configID := 22345
	expectedVersion := 3
	resource := "test_resource"

	// Pre-populate cache
	cacheKey := "getModifiableConfigVersion:22345"
	configuration := &appsec.GetConfigurationResponse{
		ID:            configID,
		LatestVersion: expectedVersion,
	}

	err := cache.Set(cache.BucketName(SubproviderName), cacheKey, configuration)
	require.NoError(t, err)

	// No client mocks needed since we should hit cache
	client := &appsec.Mock{}

	ctx := context.Background()

	useClient(client, func() {
		// Call the function under test
		result, err := getModifiableConfigVersion(ctx, configID, resource, &mockMeta{})

		// Assertions
		assert.NoError(t, err)
		assert.Equal(t, expectedVersion, result)
	})

	// Verify no API calls were made due to cache hit
	client.AssertExpectations(t)
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
	client := &appsec.Mock{}

	getConfigResponse := appsec.GetConfigurationResponse{
		ID:                configID,
		LatestVersion:     latestVersion,
		StagingVersion:    stagingVersion,
		ProductionVersion: productionVersion,
	}

	// Mock the GetConfiguration call
	client.On("GetConfiguration",
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
		Production: struct {
			Status string    `json:"status,omitempty"`
			Time   time.Time `json:"time,omitempty"`
		}{
			Status: "Inactive",
			Time:   time.Now(),
		},
		Staging: struct {
			Status string    `json:"status,omitempty"`
			Time   time.Time `json:"time,omitempty"`
		}{
			Status: "Inactive",
			Time:   time.Now(),
		},
	}

	client.On("GetConfigurationVersion",
		mock.Anything,
		appsec.GetConfigurationVersionRequest{ConfigID: configID, Version: latestVersion},
	).Return(&getConfigVersionResponse, nil).Once()

	ctx := context.Background()

	useClient(client, func() {
		// Call the function under test
		result, err := getModifiableConfigVersion(ctx, configID, resource, &mockMeta{})

		// Assertions
		assert.NoError(t, err)
		assert.Equal(t, latestVersion, result)
	})

	client.AssertExpectations(t)
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
	client := &appsec.Mock{}

	getConfigResponse := appsec.GetConfigurationResponse{
		ID:                configID,
		LatestVersion:     latestVersion,
		StagingVersion:    stagingVersion,
		ProductionVersion: productionVersion,
	}

	// Mock the GetConfiguration call
	client.On("GetConfiguration",
		mock.Anything,
		appsec.GetConfigurationRequest{ConfigID: configID},
	).Return(&getConfigResponse, nil).Once()

	// Mock the CreateConfigurationVersionClone call
	cloneResponse := appsec.CreateConfigurationVersionCloneResponse{
		ConfigID: configID,
		Version:  newClonedVersion,
	}

	client.On("CreateConfigurationVersionClone",
		mock.Anything,
		appsec.CreateConfigurationVersionCloneRequest{
			ConfigID:          configID,
			CreateFromVersion: latestVersion,
		},
	).Return(&cloneResponse, nil).Once()

	ctx := context.Background()

	useClient(client, func() {
		// Call the function under test
		result, err := getModifiableConfigVersion(ctx, configID, resource, &mockMeta{})

		// Assertions
		assert.NoError(t, err)
		assert.Equal(t, newClonedVersion, result)
	})

	client.AssertExpectations(t)
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
	client := &appsec.Mock{}

	getConfigResponse := appsec.GetConfigurationResponse{
		ID:                configID,
		LatestVersion:     latestVersion,
		StagingVersion:    stagingVersion,
		ProductionVersion: productionVersion,
	}

	// Mock the GetConfiguration call
	client.On("GetConfiguration",
		mock.Anything,
		appsec.GetConfigurationRequest{ConfigID: configID},
	).Return(&getConfigResponse, nil).Once()

	// Mock the CreateConfigurationVersionClone call
	cloneResponse := appsec.CreateConfigurationVersionCloneResponse{
		ConfigID: configID,
		Version:  newClonedVersion,
	}

	client.On("CreateConfigurationVersionClone",
		mock.Anything,
		appsec.CreateConfigurationVersionCloneRequest{
			ConfigID:          configID,
			CreateFromVersion: latestVersion,
		},
	).Return(&cloneResponse, nil).Once()

	ctx := context.Background()

	useClient(client, func() {
		// Call the function under test
		result, err := getModifiableConfigVersion(ctx, configID, resource, &mockMeta{})

		// Assertions
		assert.NoError(t, err)
		assert.Equal(t, newClonedVersion, result)
	})

	client.AssertExpectations(t)
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
	client := &appsec.Mock{}

	getConfigResponse := appsec.GetConfigurationResponse{
		ID:                configID,
		LatestVersion:     latestVersion,
		StagingVersion:    stagingVersion,
		ProductionVersion: productionVersion,
	}

	// Mock the GetConfiguration call
	client.On("GetConfiguration",
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
		Production: struct {
			Status string    `json:"status,omitempty"`
			Time   time.Time `json:"time,omitempty"`
		}{
			Status: "Inactive",
			Time:   time.Now(),
		},
		Staging: struct {
			Status string    `json:"status,omitempty"`
			Time   time.Time `json:"time,omitempty"`
		}{
			Status: "Deactivated", // Was previously active
			Time:   time.Now(),
		},
	}

	client.On("GetConfigurationVersion",
		mock.Anything,
		appsec.GetConfigurationVersionRequest{ConfigID: configID, Version: latestVersion},
	).Return(&getConfigVersionResponse, nil).Once()

	// Mock the CreateConfigurationVersionClone call
	cloneResponse := appsec.CreateConfigurationVersionCloneResponse{
		ConfigID: configID,
		Version:  newClonedVersion,
	}

	client.On("CreateConfigurationVersionClone",
		mock.Anything,
		appsec.CreateConfigurationVersionCloneRequest{
			ConfigID:          configID,
			CreateFromVersion: latestVersion,
		},
	).Return(&cloneResponse, nil).Once()

	ctx := context.Background()

	useClient(client, func() {
		// Call the function under test
		result, err := getModifiableConfigVersion(ctx, configID, resource, &mockMeta{})

		// Assertions
		assert.NoError(t, err)
		assert.Equal(t, newClonedVersion, result)
	})

	client.AssertExpectations(t)
}

func TestGetModifiableConfigVersion_GetConfigurationError(t *testing.T) {
	// Clear cache before test
	clearCache()
	defer clearCache()

	configID := 27345
	resource := "test_resource"
	expectedError := errors.New("API error")

	// Setup mock client
	client := &appsec.Mock{}

	// Mock the GetConfiguration call to return error
	client.On("GetConfiguration",
		mock.Anything,
		appsec.GetConfigurationRequest{ConfigID: configID},
	).Return(nil, expectedError).Once()

	ctx := context.Background()

	useClient(client, func() {
		// Call the function under test
		result, err := getModifiableConfigVersion(ctx, configID, resource, &mockMeta{})

		// Assertions
		assert.Error(t, err)
		assert.Contains(t, err.Error(), "API error")
		assert.Equal(t, 0, result)
	})

	client.AssertExpectations(t)
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
	client := &appsec.Mock{}

	getConfigResponse := appsec.GetConfigurationResponse{
		ID:                configID,
		LatestVersion:     latestVersion,
		StagingVersion:    stagingVersion,
		ProductionVersion: productionVersion,
	}

	// Mock the GetConfiguration call
	client.On("GetConfiguration",
		mock.Anything,
		appsec.GetConfigurationRequest{ConfigID: configID},
	).Return(&getConfigResponse, nil).Once()

	// Mock the CreateConfigurationVersionClone call to return error
	client.On("CreateConfigurationVersionClone",
		mock.Anything,
		appsec.CreateConfigurationVersionCloneRequest{
			ConfigID:          configID,
			CreateFromVersion: latestVersion,
		},
	).Return(nil, expectedError).Once()

	ctx := context.Background()

	useClient(client, func() {
		// Call the function under test
		result, err := getModifiableConfigVersion(ctx, configID, resource, &mockMeta{})

		// Assertions
		assert.Error(t, err)
		assert.Contains(t, err.Error(), "clone error")
		assert.Equal(t, 0, result)
	})

	client.AssertExpectations(t)
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
	client := &appsec.Mock{}

	getConfigResponse := appsec.GetConfigurationResponse{
		ID:                configID,
		LatestVersion:     latestVersion,
		StagingVersion:    stagingVersion,
		ProductionVersion: productionVersion,
	}

	// Mock the GetConfiguration call
	client.On("GetConfiguration",
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
		Production: struct {
			Status string    `json:"status,omitempty"`
			Time   time.Time `json:"time,omitempty"`
		}{
			Status: "Inactive",
			Time:   time.Now(),
		},
		Staging: struct {
			Status string    `json:"status,omitempty"`
			Time   time.Time `json:"time,omitempty"`
		}{
			Status: "Inactive",
			Time:   time.Now(),
		},
	}

	client.On("GetConfigurationVersion",
		mock.Anything,
		appsec.GetConfigurationVersionRequest{ConfigID: configID, Version: latestVersion},
	).Return(&getConfigVersionResponse, nil).Once()

	ctx := context.Background()

	useClient(client, func() {
		// Run multiple goroutines to test concurrency
		results := make(chan int, 3)
		errors := make(chan error, 3)

		for i := 0; i < 3; i++ {
			go func() {
				result, err := getModifiableConfigVersion(ctx, configID, resource, &mockMeta{})
				results <- result
				errors <- err
			}()
		}

		// Collect results
		for i := 0; i < 3; i++ {
			result := <-results
			err := <-errors
			assert.NoError(t, err)
			assert.Equal(t, latestVersion, result)
		}
	})

	client.AssertExpectations(t)
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
	client := &appsec.Mock{}

	getConfigResponse := appsec.GetConfigurationResponse{
		ID:                configID,
		LatestVersion:     latestVersion,
		StagingVersion:    stagingVersion,
		ProductionVersion: productionVersion,
	}

	// Mock the GetConfiguration call
	client.On("GetConfiguration",
		mock.Anything,
		appsec.GetConfigurationRequest{ConfigID: configID},
	).Return(&getConfigResponse, nil).Once()

	// Mock GetConfigurationVersion call to return error
	// In this case, the function should still work and return the latest version
	// because checkIfVersionWasPreviouslyActive handles errors gracefully
	client.On("GetConfigurationVersion",
		mock.Anything,
		appsec.GetConfigurationVersionRequest{ConfigID: configID, Version: latestVersion},
	).Return(nil, errors.New("version API error")).Once()

	ctx := context.Background()

	useClient(client, func() {
		// Call the function under test
		result, err := getModifiableConfigVersion(ctx, configID, resource, &mockMeta{})

		// Assertions - should still succeed and return latest version
		// because checkIfVersionWasPreviouslyActive returns false on error
		assert.NoError(t, err)
		assert.Equal(t, latestVersion, result)
	})

	client.AssertExpectations(t)
}
