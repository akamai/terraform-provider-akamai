package botman

import (
	"context"
	"errors"
	"fmt"
	"sync"

	"github.com/akamai/AkamaiOPEN-edgegrid-golang/v6/pkg/botman"
	"github.com/akamai/terraform-provider-akamai/v4/pkg/cache"
	akameta "github.com/akamai/terraform-provider-akamai/v4/pkg/meta"
	"github.com/apex/log"
)

var (
	botDetectionActionMutex      sync.Mutex
	customBotCategoryActionMutex sync.Mutex
	akamaiBotCategoryActionMutex sync.Mutex
	transactionalEndpointMutex   sync.Mutex
)

// getBotDetectionAction reads from the cache if present, or makes a getAll call to fetch all Bot Detection Actions for a security policy, stores in the cache and filters the required Bot Detection Action using ID.
func getBotDetectionAction(ctx context.Context, request botman.GetBotDetectionActionRequest, m interface{}) (map[string]interface{}, error) {
	meta := akameta.Must(m)
	client := inst.Client(meta)
	logger := meta.Log("BotMan", "getBotDetectionAction")

	cacheKey := fmt.Sprintf("%s:%d:%d:%s", "getBotDetectionAction", request.ConfigID, request.Version, request.SecurityPolicyID)
	botDetectionActions := &botman.GetBotDetectionActionListResponse{}
	err := cache.Get(inst, cacheKey, botDetectionActions)
	// if cache is disabled use GetBotDetectionAction to fetch one action at a time
	if errors.Is(err, cache.ErrDisabled) {
		return client.GetBotDetectionAction(ctx, request)
	}
	if err == nil {
		return filterBotDetectionAction(botDetectionActions, request, logger)
	}

	botDetectionActionMutex.Lock()
	defer func() {
		logger.Debugf("Unlocking mutex")
		botDetectionActionMutex.Unlock()
	}()

	err = cache.Get(inst, cacheKey, botDetectionActions)
	if err == nil {
		return filterBotDetectionAction(botDetectionActions, request, logger)
	}

	if !errors.Is(err, cache.ErrEntryNotFound) && !errors.Is(err, cache.ErrDisabled) {
		logger.Errorf("error reading from cache: %s", err.Error())
		return nil, err
	}

	botDetectionActions, err = client.GetBotDetectionActionList(ctx, botman.GetBotDetectionActionListRequest{
		ConfigID:         request.ConfigID,
		Version:          request.Version,
		SecurityPolicyID: request.SecurityPolicyID,
	})
	if err != nil {
		logger.Errorf("calling 'GetBotDetectionActionList': %s", err.Error())
		return nil, err
	}

	err = cache.Set(inst, cacheKey, botDetectionActions)
	if err != nil && !errors.Is(err, cache.ErrDisabled) {
		logger.Errorf("error caching botDetectionActions into cache: %s", err.Error())
		return nil, err
	}

	return filterBotDetectionAction(botDetectionActions, request, logger)
}

func filterBotDetectionAction(botDetectionActions *botman.GetBotDetectionActionListResponse, request botman.GetBotDetectionActionRequest, logger log.Interface) (map[string]interface{}, error) {
	for _, action := range botDetectionActions.Actions {
		if action["detectionId"].(string) == request.DetectionID {
			logger.Debugf("Found bot detection action for config %d version %d security policy %s bot detection %s", request.ConfigID, request.Version, request.SecurityPolicyID, request.DetectionID)
			return action, nil
		}
	}
	return nil, fmt.Errorf("BotDetectionAction with id: %s does not exist", request.DetectionID)
}

// getCustomBotCategoryAction reads from the cache if present, or makes a getAll call to fetch all Custom Bot Category Actions for a security policy, stores in the cache and filters the required Custom Bot Category Action using ID.
func getCustomBotCategoryAction(ctx context.Context, request botman.GetCustomBotCategoryActionRequest, m interface{}) (map[string]interface{}, error) {
	meta := akameta.Must(m)
	client := inst.Client(meta)
	logger := meta.Log("BotMan", "getCustomBotCategoryAction")

	cacheKey := fmt.Sprintf("%s:%d:%d:%s", "getCustomBotCategoryAction", request.ConfigID, request.Version, request.SecurityPolicyID)
	customBotCategoryActions := &botman.GetCustomBotCategoryActionListResponse{}
	err := cache.Get(inst, cacheKey, customBotCategoryActions)
	// if cache is disabled use GetCustomBotCategoryAction to fetch one action at a time
	if errors.Is(err, cache.ErrDisabled) {
		return client.GetCustomBotCategoryAction(ctx, request)
	}
	if err == nil {
		return filterCustomBotCategoryAction(customBotCategoryActions, request, logger)
	}

	customBotCategoryActionMutex.Lock()
	defer func() {
		logger.Debugf("Unlocking mutex")
		customBotCategoryActionMutex.Unlock()
	}()

	err = cache.Get(inst, cacheKey, customBotCategoryActions)
	if err == nil {
		return filterCustomBotCategoryAction(customBotCategoryActions, request, logger)
	}

	if !errors.Is(err, cache.ErrEntryNotFound) && !errors.Is(err, cache.ErrDisabled) {
		logger.Errorf("error reading from cache: %s", err.Error())
		return nil, err
	}

	customBotCategoryActions, err = client.GetCustomBotCategoryActionList(ctx, botman.GetCustomBotCategoryActionListRequest{
		ConfigID:         request.ConfigID,
		Version:          request.Version,
		SecurityPolicyID: request.SecurityPolicyID,
	})
	if err != nil {
		logger.Errorf("calling 'GetCustomBotCategoryActionList': %s", err.Error())
		return nil, err
	}

	err = cache.Set(inst, cacheKey, customBotCategoryActions)
	if err != nil && !errors.Is(err, cache.ErrDisabled) {
		logger.Errorf("error caching customBotCategoryActions into cache: %s", err.Error())
		return nil, err
	}

	return filterCustomBotCategoryAction(customBotCategoryActions, request, logger)
}

func filterCustomBotCategoryAction(customBotCategoryActions *botman.GetCustomBotCategoryActionListResponse, request botman.GetCustomBotCategoryActionRequest, logger log.Interface) (map[string]interface{}, error) {
	for _, action := range customBotCategoryActions.Actions {
		if action["categoryId"].(string) == request.CategoryID {
			logger.Debugf("Found custom bot category action for config %d version %d security policy %s category %s", request.ConfigID, request.Version, request.SecurityPolicyID, request.CategoryID)
			return action, nil
		}
	}
	return nil, fmt.Errorf("CustomBotCategoryAction with id: %s does not exist", request.CategoryID)
}

// getAkamaiBotCategoryAction reads from the cache if present, or makes a getAll call to fetch all Akamai Bot Category Actions for a security policy, stores in the cache and filters the required Akamai Bot Category Action using ID.
func getAkamaiBotCategoryAction(ctx context.Context, request botman.GetAkamaiBotCategoryActionRequest, m interface{}) (map[string]interface{}, error) {
	meta := akameta.Must(m)
	client := inst.Client(meta)
	logger := meta.Log("BotMan", "getAkamaiBotCategoryAction")

	cacheKey := fmt.Sprintf("%s:%d:%d:%s", "getAkamaiBotCategoryAction", request.ConfigID, request.Version, request.SecurityPolicyID)
	akamaiBotCategoryActions := &botman.GetAkamaiBotCategoryActionListResponse{}
	err := cache.Get(inst, cacheKey, akamaiBotCategoryActions)
	// if cache is disabled use GetAkamaiBotCategoryAction to fetch one action at a time
	if errors.Is(err, cache.ErrDisabled) {
		return client.GetAkamaiBotCategoryAction(ctx, request)
	}
	if err == nil {
		return filterAkamaiBotCategoryAction(akamaiBotCategoryActions, request, logger)
	}

	akamaiBotCategoryActionMutex.Lock()
	defer func() {
		logger.Debugf("Unlocking mutex")
		akamaiBotCategoryActionMutex.Unlock()
	}()

	err = cache.Get(inst, cacheKey, akamaiBotCategoryActions)
	if err == nil {
		return filterAkamaiBotCategoryAction(akamaiBotCategoryActions, request, logger)
	}

	if !errors.Is(err, cache.ErrEntryNotFound) && !errors.Is(err, cache.ErrDisabled) {
		logger.Errorf("error reading from cache: %s", err.Error())
		return nil, err
	}

	akamaiBotCategoryActions, err = client.GetAkamaiBotCategoryActionList(ctx, botman.GetAkamaiBotCategoryActionListRequest{
		ConfigID:         request.ConfigID,
		Version:          request.Version,
		SecurityPolicyID: request.SecurityPolicyID,
	})
	if err != nil {
		logger.Errorf("calling 'GetAkamaiBotCategoryActionList': %s", err.Error())
		return nil, err
	}

	err = cache.Set(inst, cacheKey, akamaiBotCategoryActions)
	if err != nil && !errors.Is(err, cache.ErrDisabled) {
		logger.Errorf("error caching akamaiBotCategoryActions into cache: %s", err.Error())
		return nil, err
	}

	return filterAkamaiBotCategoryAction(akamaiBotCategoryActions, request, logger)
}

func filterAkamaiBotCategoryAction(akamaiBotCategoryActions *botman.GetAkamaiBotCategoryActionListResponse, request botman.GetAkamaiBotCategoryActionRequest, logger log.Interface) (map[string]interface{}, error) {
	for _, action := range akamaiBotCategoryActions.Actions {
		if action["categoryId"].(string) == request.CategoryID {
			logger.Debugf("Found akamai bot category action for config %d version %d security policy %s category %s", request.ConfigID, request.Version, request.SecurityPolicyID, request.CategoryID)
			return action, nil
		}
	}
	return nil, fmt.Errorf("AkamaiBotCategoryAction with id: %s does not exist", request.CategoryID)
}

// getTransactionalEndpoint reads from the cache if present, or makes a getAll call to fetch all Transactional Endpoints for a security policy, stores in the cache and filters the required Transactional Endpoint using ID.
func getTransactionalEndpoint(ctx context.Context, request botman.GetTransactionalEndpointRequest, m interface{}) (map[string]interface{}, error) {
	meta := akameta.Must(m)
	client := inst.Client(meta)
	logger := meta.Log("BotMan", "getTransactionalEndpoint")

	cacheKey := fmt.Sprintf("%s:%d:%d:%s", "getTransactionalEndpoint", request.ConfigID, request.Version, request.SecurityPolicyID)
	transactionalEndpoints := &botman.GetTransactionalEndpointListResponse{}
	err := cache.Get(inst, cacheKey, transactionalEndpoints)
	// if cache is disabled use GetTransactionalEndpoint to fetch one action at a time
	if errors.Is(err, cache.ErrDisabled) {
		return client.GetTransactionalEndpoint(ctx, request)
	}
	if err == nil {
		return filterTransactionalEndpoint(transactionalEndpoints, request, logger)
	}

	transactionalEndpointMutex.Lock()
	defer func() {
		logger.Debugf("Unlocking mutex")
		transactionalEndpointMutex.Unlock()
	}()

	err = cache.Get(inst, cacheKey, transactionalEndpoints)
	if err == nil {
		return filterTransactionalEndpoint(transactionalEndpoints, request, logger)
	}

	if !errors.Is(err, cache.ErrEntryNotFound) && !errors.Is(err, cache.ErrDisabled) {
		logger.Errorf("error reading from cache: %s", err.Error())
		return nil, err
	}

	transactionalEndpoints, err = client.GetTransactionalEndpointList(ctx, botman.GetTransactionalEndpointListRequest{
		ConfigID:         request.ConfigID,
		Version:          request.Version,
		SecurityPolicyID: request.SecurityPolicyID,
	})
	if err != nil {
		logger.Errorf("calling 'GetTransactionalEndpointList': %s", err.Error())
		return nil, err
	}

	err = cache.Set(inst, cacheKey, transactionalEndpoints)
	if err != nil && !errors.Is(err, cache.ErrDisabled) {
		logger.Errorf("error caching transactionalEndpoints into cache: %s", err.Error())
		return nil, err
	}

	return filterTransactionalEndpoint(transactionalEndpoints, request, logger)
}

func filterTransactionalEndpoint(transactionalEndpoints *botman.GetTransactionalEndpointListResponse, request botman.GetTransactionalEndpointRequest, logger log.Interface) (map[string]interface{}, error) {
	for _, action := range transactionalEndpoints.Operations {
		if action["operationId"].(string) == request.OperationID {
			logger.Debugf("Found transactional endpoint for config %d version %d security policy %s operation %s", request.ConfigID, request.Version, request.SecurityPolicyID, request.OperationID)
			return action, nil
		}
	}
	return nil, fmt.Errorf("TransactionalEndpoint with id: %s does not exist", request.OperationID)
}
