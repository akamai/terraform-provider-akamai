package appsec

import (
	"context"
	"errors"
	"fmt"
	"sync"

	"github.com/akamai/AkamaiOPEN-edgegrid-golang/v12/pkg/appsec"
	akalog "github.com/akamai/AkamaiOPEN-edgegrid-golang/v12/pkg/log"
	"github.com/akamai/terraform-provider-akamai/v9/pkg/cache"
	akameta "github.com/akamai/terraform-provider-akamai/v9/pkg/meta"
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
	// and the API client obtained from m. Log messages are written to m's logger. A
	// mutex prevents calls made by multiple resources from creating unnecessary clones.
	GetModifiableConfigVersion = getModifiableConfigVersion
	// GetLatestConfigVersion returns the latest version number of the given security
	// configuration. API calls are made using the supplied context and the API client
	// obtained from m. Log messages are written to m's logger.
	GetLatestConfigVersion = getLatestConfigVersion
)

// getModifiableConfigVersion returns the number of the latest editable version
// of the given security configuration. If the most recent version is not editable
// (because it is active or was previously active in staging or production) a new
// version is cloned and the new version's number is returned. API calls are made
// using the supplied context and the API client obtained from m. Log messages are
// written to m's logger. A mutex prevents calls made by multiple resources from
// creating unnecessary clones.
func getModifiableConfigVersion(ctx context.Context, configID int, resource string, m interface{}) (int, error) {
	meta := akameta.Must(m)
	client := inst.Client(meta)
	logger := meta.Log("APPSEC", "getModifiableConfigVersion")

	// If the version info is in the cache, return it immediately.
	cacheKey := fmt.Sprintf("%s:%d", "getModifiableConfigVersion", configID)
	configuration := &appsec.GetConfigurationResponse{}
	if err := cache.Get(cache.BucketName(SubproviderName), cacheKey, configuration); err == nil {
		logger.Debugf("Resource %s returning modifiable version %d from cache", resource, configuration.LatestVersion)
		return configuration.LatestVersion, nil
	}

	logger.Debugf("Resource %s requesting mutex lock", resource)
	configCloneMutex.Lock()
	defer func() {
		logger.Debugf("Resource %s releasing mutex lock", resource)
		configCloneMutex.Unlock()
	}()

	// If the version info is in the cache, return it immediately.
	err := cache.Get(cache.BucketName(SubproviderName), cacheKey, configuration)
	if err == nil {
		logger.Debugf("Resource %s returning modifiable version %d from cache", resource, configuration.LatestVersion)
		return configuration.LatestVersion, nil
	}
	// Any error response other than 'not found' or 'cache disabled' is a problem.
	if !errors.Is(err, cache.ErrEntryNotFound) && !errors.Is(err, cache.ErrDisabled) {
		logger.Errorf("error reading from cache: %s", err.Error())
		return 0, err
	}

	// Check whether the latest version is active in staging or production
	logger.Debugf("Resource %s calling GetConfigurations", resource)
	configuration, err = client.GetConfiguration(ctx, appsec.GetConfigurationRequest{
		ConfigID: configID,
	})
	if err != nil {
		logger.Errorf("error calling 'getConfiguration': %s", err.Error())
		return 0, err
	}
	latestVersion := configuration.LatestVersion
	stagingVersion := configuration.StagingVersion
	productionVersion := configuration.ProductionVersion

	// Check if the latest version is modifiable (not currently active and not previously active)
	isModifiable, reason := checkIfVersionIsModifiable(ctx, client, configID, latestVersion, stagingVersion, productionVersion, logger)

	if isModifiable {
		// Latest version is modifiable, cache and return it
		if err := cache.Set(cache.BucketName(SubproviderName), cacheKey, configuration); err != nil {
			if !errors.Is(err, cache.ErrDisabled) {
				logger.Errorf("unable to set latestVersion %d into cache")
			}
		}
		logger.Debugf("Resource %s caching and returning latestVersion %d (staging version %d, production version %d)",
			resource, latestVersion, stagingVersion, productionVersion)
		return latestVersion, nil
	}

	// Latest version is not modifiable, need to clone a new version
	logger.Debugf("Resource %s cloning configuration version %d (%s)", resource, latestVersion, reason)
	ccr, err := client.CreateConfigurationVersionClone(ctx, appsec.CreateConfigurationVersionCloneRequest{
		ConfigID:          configID,
		CreateFromVersion: latestVersion,
	})
	if err != nil {
		logger.Errorf("error calling 'createConfigurationVersionClone': %s", err.Error())
		return 0, err
	}

	configuration.LatestVersion = ccr.Version
	if err := cache.Set(cache.BucketName(SubproviderName), cacheKey, configuration); err != nil && !errors.Is(err, cache.ErrDisabled) {
		logger.Errorf("unable to set latestVersion %d into cache: %s", err.Error())
	}

	logger.Debugf("Resource %s caching and returning new cloned version %d as modifiable version", ccr.Version)
	return ccr.Version, nil
}

// getLatestConfigVersion returns the latest version number of the given security
// configuration. API calls are made using the supplied context and the API client
// obtained from m. Log messages are written to m's logger.
func getLatestConfigVersion(ctx context.Context, configID int, m interface{}) (int, error) {
	meta := akameta.Must(m)
	client := inst.Client(meta)
	logger := meta.Log("APPSEC", "getLatestConfigVersion")

	// Return the cached value if we have one
	cacheKey := fmt.Sprintf("%s:%d", "getLatestConfigVersion", configID)
	configuration := &appsec.GetConfigurationResponse{}
	if err := cache.Get(cache.BucketName(SubproviderName), cacheKey, configuration); err == nil {
		logger.Debugf("Found config %d, returning %d as its latest version", configuration.ID, configuration.LatestVersion)
		return configuration.LatestVersion, nil
	}

	// Wait for any prior call that might be populating the cache for us; if we obtain the lock, fetch the value ourselves
	latestVersionMutex.Lock()
	defer func() {
		logger.Debugf("Unlocking latest version mutex")
		latestVersionMutex.Unlock()
	}()

	err := cache.Get(cache.BucketName(SubproviderName), cacheKey, configuration)
	if err == nil {
		logger.Debugf("Found config %d, returning %d as its latest version", configuration.ID, configuration.LatestVersion)
		return configuration.LatestVersion, nil
	}
	// Any error response other than 'not found' or 'cache disabled' is a problem.
	if !errors.Is(err, cache.ErrEntryNotFound) && !errors.Is(err, cache.ErrDisabled) {
		logger.Errorf("error reading from cache: %s", err.Error())
		return 0, err
	}

	configuration, err = client.GetConfiguration(ctx, appsec.GetConfigurationRequest{ConfigID: configID})
	if err != nil {
		logger.Errorf("error calling GetConfiguration: %s", err.Error())
		return 0, err
	}
	if err := cache.Set(cache.BucketName(SubproviderName), cacheKey, configuration); err != nil && !errors.Is(err, cache.ErrDisabled) {
		logger.Errorf("error caching latestVersion into cache: %s", err.Error())
	}

	logger.Debugf("Caching and returning %d as latest version of config %s", configuration.LatestVersion, configuration.ID)
	return configuration.LatestVersion, nil
}

// getActiveConfigVersions returns the version numbers of the given security configuration
// active in staging and production respectively. API calls are made using the supplied
// context and the API client obtained from m. Log messages are written to m's logger.
func getActiveConfigVersions(ctx context.Context, configID int, m interface{}) (int, int, error) {
	meta := akameta.Must(m)
	client := inst.Client(meta)
	logger := meta.Log("APPSEC", "getActiveConfigVersions")

	logger.Debugf("getActiveConfigVersions calling GetConfigurations")
	configuration, err := client.GetConfiguration(ctx, appsec.GetConfigurationRequest{
		ConfigID: configID,
	})
	if err != nil {
		logger.Errorf("error calling getConfiguration: %s", err.Error())
		return 0, 0, err
	}

	logger.Debugf("Found config %d, returning %d, %d as staging & production versions",
		configuration.ID, configuration.StagingVersion, configuration.ProductionVersion)

	return configuration.StagingVersion, configuration.ProductionVersion, nil
}

// checkIfVersionIsModifiable checks if a version can be modified by checking:
// 1. If it's currently active in staging or production
// 2. If it was previously active (has "Deactivated" status)
// Returns true if modifiable, false if not, along with a reason string
func checkIfVersionIsModifiable(ctx context.Context, client appsec.APPSEC, configID, versionToCheck, stagingVersion, productionVersion int, logger akalog.Interface) (bool, string) {
	// First check if version is currently active
	if versionToCheck == stagingVersion || versionToCheck == productionVersion {
		return false, "version is active in staging or production"
	}

	// Version is not currently active, check if it was previously active
	getConfigVersionRequest := appsec.GetConfigurationVersionRequest{
		ConfigID: configID,
		Version:  versionToCheck,
	}

	configVersion, err := client.GetConfigurationVersion(ctx, getConfigVersionRequest)
	if err != nil {
		logger.Warnf("Could not get configuration version for previous activity check: %v", err)
		// If we can't check previous activity, assume it's modifiable to avoid blocking
		return true, ""
	}

	// Check if the specific version has "Deactivated" status in staging or production
	stagingStatus := configVersion.Staging.Status
	productionStatus := configVersion.Production.Status

	// Check for "Deactivated" status which indicates the version was previously active
	if stagingStatus == "Deactivated" || productionStatus == "Deactivated" {
		logger.Debugf("Version %d has 'Deactivated' status (staging: %s, production: %s) - was previously active",
			versionToCheck, stagingStatus, productionStatus)
		return false, "version was previously active but is now deactivated"
	}

	logger.Debugf("Version %d is modifiable (not currently active and not previously active)", versionToCheck)
	return true, ""
}
