package appsec

import (
	"context"
	"encoding/json"
	"errors"
	"fmt"
	"sync"

	"github.com/akamai/AkamaiOPEN-edgegrid-golang/v14/pkg/appsec"
	"github.com/akamai/AkamaiOPEN-edgegrid-golang/v14/pkg/log"
	"github.com/akamai/terraform-provider-akamai/v11/pkg/cache"
	akameta "github.com/akamai/terraform-provider-akamai/v11/pkg/meta"
)

var (
	reputationProfMutex sync.Mutex

	// ErrReputationProfileNotFound is returned when a reputation profile ID is not present in the cached list.
	ErrReputationProfileNotFound = errors.New("reputation profile not found")
)

// getReputationProfile reads from the cache if present, or makes a getAll call to fetch all
// reputation profiles for a configuration version, stores in the cache and filters using ID.
func getReputationProfile(ctx context.Context, configID, version, reputationProfileID int, m interface{}) (*appsec.GetReputationProfileResponse, error) {
	meta := akameta.Must(m)
	client := meta.Client().GetAPPSEC()
	logger := meta.Log("APPSEC", "getReputationProfile")

	return getCached(
		&reputationProfMutex,
		fmt.Sprintf("getReputationProfiles:%d:%d", configID, version),
		logger,
		// onDisabled: single-item fetch when caching is off.
		func() (*appsec.GetReputationProfileResponse, error) {
			return client.GetReputationProfile(ctx, appsec.GetReputationProfileRequest{
				ConfigID:            configID,
				ConfigVersion:       version,
				ReputationProfileId: reputationProfileID,
			})
		},
		// fetch: bulk-list call used to populate the cache on a miss.
		func() (*appsec.GetReputationProfilesResponse, error) {
			r, err := client.GetReputationProfiles(ctx, appsec.GetReputationProfilesRequest{
				ConfigID:      configID,
				ConfigVersion: version,
			})
			if err != nil {
				logger.Errorf("calling 'GetReputationProfiles': %s", err.Error())
			}
			return r, err
		},
		// filter: locate the requested profile in the cached list.
		func(list *appsec.GetReputationProfilesResponse) (*appsec.GetReputationProfileResponse, error) {
			return filterReputationProfile(list, reputationProfileID, logger)
		},
	)
}

func filterReputationProfile(reputationProfiles *appsec.GetReputationProfilesResponse, reputationProfileID int, logger log.Interface) (*appsec.GetReputationProfileResponse, error) {
	for _, p := range reputationProfiles.ReputationProfiles {
		if p.ID == reputationProfileID {
			logger.Debugf("Found reputation profile %d in cached list", reputationProfileID)
			// JSON round-trip converts the list-item type to GetReputationProfileResponse
			b, err := json.Marshal(p)
			if err != nil {
				return nil, err
			}
			var result appsec.GetReputationProfileResponse
			if err = json.Unmarshal(b, &result); err != nil {
				return nil, err
			}
			return &result, nil
		}
	}
	return nil, fmt.Errorf("reputation profile with id: %d does not exist: %w", reputationProfileID, ErrReputationProfileNotFound)
}

// invalidateReputationProfileCache removes the cached reputation-profile list for the given config version.
// Call this after a reputation profile is created, updated, or deleted so the next read fetches fresh data.
// The mutex is held so that an in-flight getReputationProfile call cannot overwrite the sentinel with stale data.
func invalidateReputationProfileCache(configID, version int, logger log.Interface) {
	reputationProfMutex.Lock()
	defer reputationProfMutex.Unlock()
	if err := cache.Delete(cache.BucketName(SubproviderName), fmt.Sprintf("getReputationProfiles:%d:%d", configID, version)); err != nil && !errors.Is(err, cache.ErrDisabled) {
		logger.Errorf("error invalidating reputation profile cache: %s", err.Error())
	}
}
