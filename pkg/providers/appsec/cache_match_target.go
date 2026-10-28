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
	matchTargetMutex sync.Mutex

	// ErrMatchTargetNotFound is returned when a match target ID is not present in the cached list.
	ErrMatchTargetNotFound = errors.New("match target not found")
)

// getMatchTarget reads from the cache if present, or makes a getAll call to fetch all match targets
// for a configuration version, stores in the cache and filters the required match target using ID.
func getMatchTarget(ctx context.Context, configID, version int64, targetID int, m interface{}) (*appsec.GetMatchTargetResponse, error) {
	meta := akameta.Must(m)
	client := meta.Client().GetAPPSEC()
	logger := meta.Log("APPSEC", "getMatchTarget")

	return getCached(
		&matchTargetMutex,
		fmt.Sprintf("getMatchTargets:%d:%d", configID, version),
		logger,
		// onDisabled: single-item fetch when caching is off.
		func() (*appsec.GetMatchTargetResponse, error) {
			return client.GetMatchTarget(ctx, appsec.GetMatchTargetRequest{
				ConfigID:      configID,
				ConfigVersion: version,
				TargetID:      targetID,
			})
		},
		// fetch: bulk-list call used to populate the cache on a miss.
		func() (*appsec.GetMatchTargetsResponse, error) {
			r, err := client.GetMatchTargets(ctx, appsec.GetMatchTargetsRequest{
				ConfigID:      configID,
				ConfigVersion: version,
			})
			if err != nil {
				logger.Errorf("calling 'GetMatchTargets': %s", err.Error())
			}
			return r, err
		},
		// filter: locate the requested target in the cached list.
		func(list *appsec.GetMatchTargetsResponse) (*appsec.GetMatchTargetResponse, error) {
			return filterMatchTarget(list, targetID, logger)
		},
	)
}

// filterMatchTarget searches both APITargets and WebsiteTargets for the given targetID.
// The list response splits targets into two slices by type; GetMatchTargetResponse is a
// unified struct covering both, so a JSON round-trip normalises whichever slice item is
// found into the single-GET response shape stored in Terraform state.
func filterMatchTarget(matchTargets *appsec.GetMatchTargetsResponse, targetID int, logger log.Interface) (*appsec.GetMatchTargetResponse, error) {
	for _, t := range matchTargets.MatchTargets.APITargets {
		if t.TargetID == targetID {
			logger.Debugf("Found API match target %d in cached list", targetID)
			return roundtripToMatchTargetResponse(t)
		}
	}
	for _, t := range matchTargets.MatchTargets.WebsiteTargets {
		if t.TargetID == targetID {
			logger.Debugf("Found website match target %d in cached list", targetID)
			return roundtripToMatchTargetResponse(t)
		}
	}
	return nil, fmt.Errorf("match target with id: %d does not exist: %w", targetID, ErrMatchTargetNotFound)
}

func roundtripToMatchTargetResponse(item interface{}) (*appsec.GetMatchTargetResponse, error) {
	b, err := json.Marshal(item)
	if err != nil {
		return nil, err
	}
	var result appsec.GetMatchTargetResponse
	if err = json.Unmarshal(b, &result); err != nil {
		return nil, err
	}
	return &result, nil
}

// invalidateMatchTargetCache removes the cached match-target list for the given config version.
// Call this after a match target is created, updated, or deleted so the next read fetches fresh data.
// The mutex is held so that an in-flight getMatchTarget call cannot overwrite the sentinel with stale data.
func invalidateMatchTargetCache(configID, version int64, logger log.Interface) {
	matchTargetMutex.Lock()
	defer matchTargetMutex.Unlock()
	if err := cache.Delete(cache.BucketName(SubproviderName), fmt.Sprintf("getMatchTargets:%d:%d", configID, version)); err != nil && !errors.Is(err, cache.ErrDisabled) {
		logger.Errorf("error invalidating match target cache: %s", err.Error())
	}
}
