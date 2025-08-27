package accountprotection

import (
	"context"
	"errors"
	"fmt"
	"sync"

	apr "github.com/akamai/AkamaiOPEN-edgegrid-golang/v12/pkg/accountprotection"
	"github.com/akamai/AkamaiOPEN-edgegrid-golang/v12/pkg/log"
	"github.com/akamai/terraform-provider-akamai/v9/pkg/cache"
	akameta "github.com/akamai/terraform-provider-akamai/v9/pkg/meta"
)

var (
	protectedOperationMutex sync.Mutex
)

// getProtectedOperations reads from the cache if present, or makes a getAll call to fetch all Protected Operations for a security policy, stores in the cache and filters the required Transactional Endpoint using ID.
func getProtectedOperations(ctx context.Context, request apr.GetProtectedOperationByIDRequest, m interface{}) (map[string]interface{}, error) {
	meta := akameta.Must(m)
	client := inst.Client(meta)
	logger := meta.Log("accountprotection", "cache-getProtectedOperations")

	cacheKey := fmt.Sprintf("%s:%d:%d:%s", "aprGetProtectedOperations", request.ConfigID, request.Version, request.SecurityPolicyID)
	protectedOperations := &apr.ListProtectedOperationsResponse{}
	err := cache.Get(cache.BucketName(SubproviderName), cacheKey, protectedOperations)
	if err == nil {
		return filterProtectedOperation(protectedOperations, request, logger)
	}
	// if cache is disabled use GetProtectedOperationByID to fetch one action at a time
	if errors.Is(err, cache.ErrDisabled) {
		response, err := client.GetProtectedOperationByID(ctx, request)
		if err != nil {
			logger.Errorf("error getting response from GetProtectedOperationByID : %s", err.Error())
			return nil, err
		}
		return filterProtectedOperation(response, request, logger)
	}

	if !errors.Is(err, cache.ErrEntryNotFound) {
		logger.Errorf("error reading from cache: %s", err.Error())
		return nil, err
	}

	protectedOperationMutex.Lock()
	defer func() {
		logger.Debugf("Unlocking mutex")
		protectedOperationMutex.Unlock()
	}()

	protectedOperations, err = client.ListProtectedOperations(ctx, apr.ListProtectedOperationsRequest{
		ConfigID:         request.ConfigID,
		Version:          request.Version,
		SecurityPolicyID: request.SecurityPolicyID,
	})
	if err != nil {
		logger.Errorf("calling 'ListProtectedOperations': %s", err.Error())
		return nil, err
	}

	err = cache.Set(cache.BucketName(SubproviderName), cacheKey, protectedOperations)
	if err != nil && !errors.Is(err, cache.ErrDisabled) {
		logger.Errorf("error caching protectedOperations into cache: %s", err.Error())
		return nil, err
	}

	return filterProtectedOperation(protectedOperations, request, logger)
}

func filterProtectedOperation(protectedOperations *apr.ListProtectedOperationsResponse, request apr.GetProtectedOperationByIDRequest, logger log.Interface) (map[string]interface{}, error) {
	for _, action := range protectedOperations.Operations {
		if action["operationId"].(string) == request.OperationID {
			logger.Debugf("Found apr protected operation for config %d version %d security policy %s operation %s", request.ConfigID, request.Version, request.SecurityPolicyID, request.OperationID)
			return action, nil
		}
	}
	return nil, fmt.Errorf("ProtectedOperation with id: %s does not exist", request.OperationID)
}
