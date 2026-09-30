package appsec

import (
	"context"
	"errors"
	"fmt"
	"sync"

	"github.com/akamai/AkamaiOPEN-edgegrid-golang/v14/pkg/appsec"
	"github.com/akamai/terraform-provider-akamai/v11/pkg/cache"
	"github.com/hashicorp/terraform-plugin-log/tflog"
)

// Utility functions for determining current and latest versions of a security
// configuration, and for identifying a modifiable (editable) version.

var (
	configCloneMutex   sync.Mutex
	latestVersionMutex sync.Mutex
	// GetModifiableConfigVersion returns the number of the latest editable version
	// of the given security configuration. If the most recent version is not editable
	// (because it is active in staging or production) a new version is cloned and the
	// new version's number is returned. API calls are made using the supplied context
	// and the passed API client.
	// A mutex prevents calls made by multiple resources from creating unnecessary clones.
	GetModifiableConfigVersion = getModifiableConfigVersion
	// GetLatestConfigVersion returns the latest version number of the given security
	// configuration. API calls are made using the supplied context
	// and the passed API client.
	GetLatestConfigVersion = getLatestConfigVersion
)

// modifiableVersionCacheKey returns the cache key used by getModifiableConfigVersion.
// This key is only written after verifying the stored version is modifiable, so
// cache hits on this key are always safe to return without re-checking modifiability.
func modifiableVersionCacheKey(configID int) string {
	return fmt.Sprintf("modifiable:%d", configID)
}

// latestVersionCacheKey returns the cache key used by getLatestConfigVersion.
func latestVersionCacheKey(configID int) string {
	return fmt.Sprintf("latest:%d", configID)
}

// invalidateConfigCache clears both cached versions for the given configID.
// Call this after activation or deactivation completes, since either can change which version is active.
func invalidateConfigCache(configID int) {
	_ = cache.Delete(cache.BucketName(SubproviderName), modifiableVersionCacheKey(configID))
	_ = cache.Delete(cache.BucketName(SubproviderName), latestVersionCacheKey(configID))
}

// getModifiableConfigVersion returns the number of the latest editable version
// of the given security configuration. If the most recent version is not editable
// (because it is active or was previously active in staging or production) a new
// version is cloned and the new version's number is returned. API calls are made
// using the supplied context and the passed API client.
// A mutex prevents calls made by multiple resources from creating unnecessary clones.
func getModifiableConfigVersion(ctx context.Context, configID int, resource string, client appsec.APPSEC) (int, error) {
	// If the version info is in the cache, return it immediately.
	cacheKey := modifiableVersionCacheKey(configID)
	configuration := &appsec.GetConfigurationResponse{}
	if err := cache.Get(cache.BucketName(SubproviderName), cacheKey, configuration); err == nil {
		tflog.Debug(ctx, "returning modifiable version from cache", map[string]any{"resource": resource, "version": configuration.LatestVersion})
		return configuration.LatestVersion, nil
	}

	tflog.Debug(ctx, "requesting mutex lock", map[string]any{"resource": resource})
	configCloneMutex.Lock()
	defer func() {
		tflog.Debug(ctx, "releasing mutex lock", map[string]any{"resource": resource})
		configCloneMutex.Unlock()
	}()

	// If the version info is in the cache, return it immediately.
	err := cache.Get(cache.BucketName(SubproviderName), cacheKey, configuration)
	if err == nil {
		tflog.Debug(ctx, "returning modifiable version from cache", map[string]any{"resource": resource, "version": configuration.LatestVersion})
		return configuration.LatestVersion, nil
	}
	// Any error response other than 'not found' or 'cache disabled' is a problem.
	if !errors.Is(err, cache.ErrEntryNotFound) && !errors.Is(err, cache.ErrDisabled) {
		tflog.Error(ctx, "error reading from cache", map[string]any{"error": err.Error()})
		return 0, err
	}

	// Check whether the latest version is active in staging or production
	tflog.Debug(ctx, "calling GetConfiguration", map[string]any{"resource": resource})
	configuration, err = client.GetConfiguration(ctx, appsec.GetConfigurationRequest{
		ConfigID: configID,
	})
	if err != nil {
		tflog.Error(ctx, "error calling GetConfiguration", map[string]any{"error": err.Error()})
		return 0, err
	}
	latestVersion := configuration.LatestVersion
	stagingVersion := configuration.StagingVersion
	productionVersion := configuration.ProductionVersion

	// Check if the latest version is modifiable (not currently active and not previously active)
	isModifiable, reason := checkIfVersionIsModifiable(ctx, client, configID, latestVersion, stagingVersion, productionVersion)

	if isModifiable {
		// Latest version is modifiable, cache and return it
		if err := cache.Set(cache.BucketName(SubproviderName), cacheKey, configuration); err != nil {
			if !errors.Is(err, cache.ErrDisabled) {
				tflog.Error(ctx, "unable to set latestVersion into cache", map[string]any{"latestVersion": latestVersion, "error": err.Error()})
			}
		}
		tflog.Debug(ctx, "caching and returning latest version", map[string]any{
			"resource": resource, "latestVersion": latestVersion,
			"stagingVersion": stagingVersion, "productionVersion": productionVersion,
		})
		return latestVersion, nil
	}

	// Latest version is not modifiable, need to clone a new version
	tflog.Debug(ctx, "cloning configuration version", map[string]any{"resource": resource, "latestVersion": latestVersion, "reason": reason})
	ccr, err := client.CreateConfigurationVersionClone(ctx, appsec.CreateConfigurationVersionCloneRequest{
		ConfigID:          configID,
		CreateFromVersion: latestVersion,
	})
	if err != nil {
		tflog.Error(ctx, "error calling CreateConfigurationVersionClone", map[string]any{"error": err.Error()})
		return 0, err
	}

	configuration.LatestVersion = ccr.Version
	if err := cache.Set(cache.BucketName(SubproviderName), cacheKey, configuration); err != nil && !errors.Is(err, cache.ErrDisabled) {
		tflog.Error(ctx, "unable to set latestVersion into cache", map[string]any{"error": err.Error(), "latestVersion": latestVersion})
	}
	// Also update the shared getLatestConfigVersion cache so getLatestConfigVersion and
	// getActiveConfigVersions immediately reflect the new LatestVersion without an extra API call.
	latestCacheKey := latestVersionCacheKey(configID)
	if err := cache.Set(cache.BucketName(SubproviderName), latestCacheKey, configuration); err != nil && !errors.Is(err, cache.ErrDisabled) {
		tflog.Error(ctx, "unable to update getLatestConfigVersion cache after clone", map[string]any{"error": err.Error(), "latestVersion": ccr.Version})
	}

	tflog.Debug(ctx, "caching and returning new cloned version as modifiable version", map[string]any{"resource": resource, "version": ccr.Version})
	return ccr.Version, nil
}

// getLatestConfigVersion returns the latest version number of the given security
// configuration. API calls are made using the supplied context and the passed API client.
func getLatestConfigVersion(ctx context.Context, configID int, client appsec.APPSEC) (int, error) {
	// Return the cached value if we have one
	cacheKey := latestVersionCacheKey(configID)
	configuration := &appsec.GetConfigurationResponse{}
	if err := cache.Get(cache.BucketName(SubproviderName), cacheKey, configuration); err == nil {
		tflog.Debug(ctx, "found config in cache, returning latest version", map[string]any{"configID": configuration.ID, "version": configuration.LatestVersion})
		return configuration.LatestVersion, nil
	}

	// Wait for any prior call that might be populating the cache for us; if we obtain the lock, fetch the value ourselves
	latestVersionMutex.Lock()
	defer func() {
		tflog.Debug(ctx, "unlocking latest version mutex")
		latestVersionMutex.Unlock()
	}()

	err := cache.Get(cache.BucketName(SubproviderName), cacheKey, configuration)
	if err == nil {
		tflog.Debug(ctx, "found config in cache, returning latest version", map[string]any{"configID": configuration.ID, "version": configuration.LatestVersion})
		return configuration.LatestVersion, nil
	}
	// Any error response other than 'not found' or 'cache disabled' is a problem.
	if !errors.Is(err, cache.ErrEntryNotFound) && !errors.Is(err, cache.ErrDisabled) {
		tflog.Error(ctx, "error reading from cache", map[string]any{"error": err.Error()})
		return 0, err
	}

	configuration, err = client.GetConfiguration(ctx, appsec.GetConfigurationRequest{ConfigID: configID})
	if err != nil {
		tflog.Error(ctx, "error calling GetConfiguration", map[string]any{"error": err.Error(), "configID": configID})
		return 0, err
	}
	if err := cache.Set(cache.BucketName(SubproviderName), cacheKey, configuration); err != nil && !errors.Is(err, cache.ErrDisabled) {
		tflog.Error(ctx, "error caching latestVersion into cache", map[string]any{"error": err.Error(), "latestVersion": configuration.LatestVersion})
	}

	tflog.Debug(ctx, "Caching and returning latest version of config", map[string]any{"configID": configID, "latestVersion": configuration.LatestVersion})
	return configuration.LatestVersion, nil
}

// getActiveConfigVersions returns the version numbers of the given security configuration
// active in staging and production respectively. API calls are made using the supplied context
// and the passed API client.
func getActiveConfigVersions(ctx context.Context, configID int, client appsec.APPSEC) (int, int, error) {
	// Reuse the same cache key and mutex as getLatestConfigVersion — same API call, same response struct.
	// This ensures both functions stay in sync: evicting one evicts the other.
	cacheKey := latestVersionCacheKey(configID)
	configuration := &appsec.GetConfigurationResponse{}
	if err := cache.Get(cache.BucketName(SubproviderName), cacheKey, configuration); err == nil {
		tflog.Debug(ctx, "found config in cache, returning active versions", map[string]any{
			"configID": configID, "stagingVersion": configuration.StagingVersion, "productionVersion": configuration.ProductionVersion,
		})
		return configuration.StagingVersion, configuration.ProductionVersion, nil
	}

	latestVersionMutex.Lock()
	defer latestVersionMutex.Unlock()

	// Double-check after acquiring the lock — another goroutine may have populated it.
	err := cache.Get(cache.BucketName(SubproviderName), cacheKey, configuration)
	if err == nil {
		return configuration.StagingVersion, configuration.ProductionVersion, nil
	}
	if !errors.Is(err, cache.ErrEntryNotFound) && !errors.Is(err, cache.ErrDisabled) {
		tflog.Error(ctx, "error reading from cache", map[string]any{"error": err.Error()})
		return 0, 0, err
	}

	tflog.Debug(ctx, "getActiveConfigVersions calling GetConfiguration", map[string]any{"configID": configID})
	configuration, err = client.GetConfiguration(ctx, appsec.GetConfigurationRequest{ConfigID: configID})
	if err != nil {
		tflog.Error(ctx, "error calling GetConfiguration", map[string]any{"error": err.Error()})
		return 0, 0, err
	}
	if err := cache.Set(cache.BucketName(SubproviderName), cacheKey, configuration); err != nil && !errors.Is(err, cache.ErrDisabled) {
		tflog.Error(ctx, "error caching config", map[string]any{"error": err.Error()})
	}
	tflog.Debug(ctx, "returning active versions", map[string]any{
		"configID": configID, "stagingVersion": configuration.StagingVersion, "productionVersion": configuration.ProductionVersion,
	})
	return configuration.StagingVersion, configuration.ProductionVersion, nil
}

// checkIfVersionIsModifiable checks if a version can be modified by checking:
// 1. If it's currently active in staging or production
// 2. If it was previously active (has "Deactivated" status)
// Returns true if modifiable, false if not, along with a reason string.
func checkIfVersionIsModifiable(ctx context.Context, client appsec.APPSEC, configID, versionToCheck, stagingVersion, productionVersion int) (bool, string) {
	// First check if version is currently active
	if versionToCheck == stagingVersion || versionToCheck == productionVersion {
		return false, "version is active in staging or production"
	}

	// Version is not currently active, check if it was previously active
	configVersion, err := client.GetConfigurationVersion(ctx, appsec.GetConfigurationVersionRequest{
		ConfigID: configID,
		Version:  versionToCheck,
	})
	if err != nil {
		tflog.Warn(ctx, "could not get configuration version for previous activity check", map[string]any{"error": err.Error()})
		// If we can't check previous activity, assume it's modifiable to avoid blocking
		return true, ""
	}

	// Check if the specific version has "Deactivated" status in staging or production
	// Check for "Deactivated" status which indicates the version was previously active
	if configVersion.Staging.Status == "Deactivated" || configVersion.Production.Status == "Deactivated" {
		return false, "version was previously active but is now deactivated"
	}

	tflog.Debug(ctx, fmt.Sprintf("version %d is modifiable (not currently active and not previously active)", versionToCheck))
	return true, ""
}
