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
	ratePolicyMutex sync.Mutex

	// ErrRatePolicyNotFound is returned when a rate policy ID is not present in the cached list.
	ErrRatePolicyNotFound = errors.New("rate policy not found")
)

// getRatePolicy reads from the cache if present, or makes a getAll call to fetch all rate policies
// for a configuration version, stores in the cache and filters the required rate policy using ID.
func getRatePolicy(ctx context.Context, configID, version int64, ratePolicyID int, m interface{}) (*appsec.GetRatePolicyResponse, error) {
	meta := akameta.Must(m)
	client := meta.Client().GetAPPSEC()
	logger := meta.Log("APPSEC", "getRatePolicy")

	return getCached(
		&ratePolicyMutex,
		fmt.Sprintf("getRatePolicies:%d:%d", configID, version),
		logger,
		// onDisabled: single-item fetch when caching is off.
		func() (*appsec.GetRatePolicyResponse, error) {
			return client.GetRatePolicy(ctx, appsec.GetRatePolicyRequest{
				ConfigID:      configID,
				ConfigVersion: version,
				RatePolicyID:  ratePolicyID,
			})
		},
		// fetch: bulk-list call used to populate the cache on a miss.
		func() (*appsec.GetRatePoliciesResponse, error) {
			r, err := client.GetRatePolicies(ctx, appsec.GetRatePoliciesRequest{
				ConfigID:      configID,
				ConfigVersion: version,
			})
			if err != nil {
				logger.Errorf("calling 'GetRatePolicies': %s", err.Error())
			}
			return r, err
		},
		// filter: locate the requested policy in the cached list.
		func(list *appsec.GetRatePoliciesResponse) (*appsec.GetRatePolicyResponse, error) {
			return filterRatePolicy(list, ratePolicyID, logger)
		},
	)
}

func filterRatePolicy(ratePolicies *appsec.GetRatePoliciesResponse, ratePolicyID int, logger log.Interface) (*appsec.GetRatePolicyResponse, error) {
	for _, p := range ratePolicies.RatePolicies {
		if p.ID == ratePolicyID {
			logger.Debugf("Found rate policy %d in cached list", ratePolicyID)
			// JSON round-trip converts the list-item type to GetRatePolicyResponse,
			// ensuring the JSON written to Terraform state is identical to what a
			// direct client.GetRatePolicy call returns (prevents spurious plan diffs).
			b, err := json.Marshal(p)
			if err != nil {
				return nil, err
			}
			var result appsec.GetRatePolicyResponse
			if err = json.Unmarshal(b, &result); err != nil {
				return nil, err
			}
			return &result, nil
		}
	}
	return nil, fmt.Errorf("rate policy with id: %d does not exist: %w", ratePolicyID, ErrRatePolicyNotFound)
}

// invalidateRatePolicyCache removes the cached rate-policy list for the given config version.
// Call this after a rate policy is created, updated, or deleted so the next read fetches fresh data.
// The mutex is held so that an in-flight getRatePolicy call cannot overwrite the sentinel with stale data.
func invalidateRatePolicyCache(configID, version int64, logger log.Interface) {
	ratePolicyMutex.Lock()
	defer ratePolicyMutex.Unlock()
	if err := cache.Delete(cache.BucketName(SubproviderName), fmt.Sprintf("getRatePolicies:%d:%d", configID, version)); err != nil && !errors.Is(err, cache.ErrDisabled) {
		logger.Errorf("error invalidating rate policy cache: %s", err.Error())
	}
}
