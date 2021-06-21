package appsec

import (
	"context"
	"sync"

	"github.com/akamai/AkamaiOPEN-edgegrid-golang/v2/pkg/appsec"
	"github.com/akamai/terraform-provider-akamai/v2/pkg/akamai"
)

// Utility functions for determining current and latest versions of a security
// configuration, and for identifying a modifiable (editable) version.

var (
	configCloneMutex sync.Mutex
)

// getModifiableConfigVersion returns the number of the latest editable version
// of the given security configuration. If the most recent version is not editable
// (because it is active in staging or production) a new version is cloned and the
// new version's number is returned. API calls are made using the supplied context
// and the API client obtained from m. Log messages are written to m's logger. A
// mutex prevents calls made by multiple resources from creating unnecessary clones.
func getModifiableConfigVersion(ctx context.Context, configID int, resource string, m interface{}) int {
	meta := akamai.Meta(m)
	client := inst.Client(meta)
	logger := meta.Log("APPSEC", "getModifiableConfigVersion")

	logger.Debugf("Resource %s requesting mutex lock", resource)
	configCloneMutex.Lock()
	defer func() {
		logger.Debugf("Resource %s unlocking mutex", resource)
		configCloneMutex.Unlock()
	}()

	logger.Debugf("Resource %s calling GetConfigurations", resource)
	getConfigurationRequest := appsec.GetConfigurationRequest{}
	getConfigurationRequest.ConfigID = configID

	configuration, err := client.GetConfiguration(ctx, getConfigurationRequest)
	if err != nil {
		logger.Errorf("calling 'getConfiguration': %s", err.Error())
		return 0 // diag.FromErr(err)
	}

	var latestVersion, stagingVersion, productionVersion int

	latestVersion = configuration.LatestVersion
	stagingVersion = configuration.StagingVersion
	productionVersion = configuration.ProductionVersion

	if latestVersion != stagingVersion && latestVersion != productionVersion {
		logger.Debugf("Resource %s returning latestVersion %d - staging version %d, production version %d",
			resource, latestVersion, stagingVersion, productionVersion)
		return latestVersion
	}

	createConfigurationVersionClone := appsec.CreateConfigurationVersionCloneRequest{}
	createConfigurationVersionClone.ConfigID = configID
	createConfigurationVersionClone.CreateFromVersion = latestVersion

	logger.Debugf("Resource %s cloning configuration version %d", resource, latestVersion)
	ccr, err := client.CreateConfigurationVersionClone(ctx, createConfigurationVersionClone)
	if err != nil {
		logger.Errorf("calling 'createConfigurationVersionClone': %s", err.Error())
		return 0 // diag.FromErr(err)
	}

	logger.Debugf("Resource %s returning new latestVersion %d as modifiable version", resource, ccr.Version)
	return ccr.Version
}

// getLatestConfigVersion returns the latest version number of the given security
// configuration. API calls are made using the supplied context and the API client
// obtained from m. Log messages are written to m's logger.
func getLatestConfigVersion(ctx context.Context, configID int, m interface{}) int {
	meta := akamai.Meta(m)
	client := inst.Client(meta)
	logger := meta.Log("APPSEC", "getLatestConfigVersion")

	logger.Debugf("getLatestConfigVersion calling GetConfigurations")
	getConfigurationRequest := appsec.GetConfigurationRequest{}
	getConfigurationRequest.ConfigID = configID

	configuration, err := client.GetConfiguration(ctx, getConfigurationRequest)
	if err != nil {
		logger.Errorf("Did not find config with ID %d, returning 0", configID)
		logger.Errorf("calling 'getConfiguration': %s", err.Error())
		return 0
	}

	logger.Debugf("Found config %v, returning %d as its latest version", configuration.ID, configuration.LatestVersion)
	return configuration.LatestVersion
}

// getActiveConfigVersions returns the version numbers of the given security configuration
// active in staging and production respectively. API calls are made using the supplied
// context and the API client obtained from m. Log messages are written to m's logger.
func getActiveConfigVersions(ctx context.Context, configID int, m interface{}) (int, int) {
	meta := akamai.Meta(m)
	client := inst.Client(meta)
	logger := meta.Log("APPSEC", "getActiveConfigVersions")

	logger.Debugf("getActiveConfigVersions calling GetConfigurations")
	getConfigurationRequest := appsec.GetConfigurationRequest{}
	getConfigurationRequest.ConfigID = configID

	configuration, err := client.GetConfiguration(ctx, getConfigurationRequest)
	if err != nil {
		logger.Errorf("calling 'getConfiguration': %s", err.Error())
		return 0, 0
	}

	logger.Debugf("Found config %v, returning %d, %d as its staging & production versions",
		configuration.ID, configuration.StagingVersion, configuration.ProductionVersion)
	return configuration.StagingVersion, configuration.ProductionVersion
}
