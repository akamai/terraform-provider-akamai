package botman

import (
	"context"
	"errors"
	"fmt"
	"sync"

	"github.com/akamai/AkamaiOPEN-edgegrid-golang/v10/pkg/botman"
	"github.com/akamai/AkamaiOPEN-edgegrid-golang/v10/pkg/log"
	"github.com/akamai/terraform-provider-akamai/v8/pkg/cache"
	akameta "github.com/akamai/terraform-provider-akamai/v8/pkg/meta"
)

var (
	botDetectionActionMutex                       sync.Mutex
	customBotCategoryActionMutex                  sync.Mutex
	akamaiBotCategoryActionMutex                  sync.Mutex
	transactionalEndpointMutex                    sync.Mutex
	akamaiBotCategoryMutex                        sync.Mutex
	akamaiDefinedBotMutex                         sync.Mutex
	botDetectionMutex                             sync.Mutex
	contentProtectionRuleMutex                    sync.Mutex
	contentProtectionJavaScriptInjectionRuleMutex sync.Mutex
)

// getBotDetectionAction reads from the cache if present, or makes a getAll call to fetch all Bot Detection Actions for a security policy, stores in the cache and filters the required Bot Detection Action using ID.
func getBotDetectionAction(ctx context.Context, request botman.GetBotDetectionActionRequest, m interface{}) (map[string]interface{}, error) {
	meta := akameta.Must(m)
	client := inst.Client(meta)
	logger := meta.Log("BotMan", "getBotDetectionAction")

	cacheKey := fmt.Sprintf("%s:%d:%d:%s", "getBotDetectionAction", request.ConfigID, request.Version, request.SecurityPolicyID)
	botDetectionActions := &botman.GetBotDetectionActionListResponse{}
	err := cache.Get(cache.BucketName(SubproviderName), cacheKey, botDetectionActions)
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

	err = cache.Get(cache.BucketName(SubproviderName), cacheKey, botDetectionActions)
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

	err = cache.Set(cache.BucketName(SubproviderName), cacheKey, botDetectionActions)
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
	err := cache.Get(cache.BucketName(SubproviderName), cacheKey, customBotCategoryActions)
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

	err = cache.Get(cache.BucketName(SubproviderName), cacheKey, customBotCategoryActions)
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

	err = cache.Set(cache.BucketName(SubproviderName), cacheKey, customBotCategoryActions)
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
	err := cache.Get(cache.BucketName(SubproviderName), cacheKey, akamaiBotCategoryActions)
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

	err = cache.Get(cache.BucketName(SubproviderName), cacheKey, akamaiBotCategoryActions)
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

	err = cache.Set(cache.BucketName(SubproviderName), cacheKey, akamaiBotCategoryActions)
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
	err := cache.Get(cache.BucketName(SubproviderName), cacheKey, transactionalEndpoints)
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

	err = cache.Get(cache.BucketName(SubproviderName), cacheKey, transactionalEndpoints)
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

	err = cache.Set(cache.BucketName(SubproviderName), cacheKey, transactionalEndpoints)
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

func getAkamaiBotCategoryList(ctx context.Context, request botman.GetAkamaiBotCategoryListRequest, m interface{}) (*botman.GetAkamaiBotCategoryListResponse, error) {
	meta := akameta.Must(m)
	client := inst.Client(meta)
	logger := meta.Log("BotMan", "getAkamaiBotCategory")

	cacheKey := "getAkamaiBotCategory"
	akamaiBotCategoryList := &botman.GetAkamaiBotCategoryListResponse{}
	err := cache.Get(cache.BucketName(SubproviderName), cacheKey, akamaiBotCategoryList)
	// if cache is disabled make a direct all to GetAkamaiBotCategoryList
	if errors.Is(err, cache.ErrDisabled) {
		return client.GetAkamaiBotCategoryList(ctx, request)
	}
	if err == nil {
		return filterAkamaiBotCategoryList(akamaiBotCategoryList, request), nil
	}

	akamaiBotCategoryMutex.Lock()
	defer func() {
		logger.Debugf("Unlocking mutex")
		akamaiBotCategoryMutex.Unlock()
	}()

	err = cache.Get(cache.BucketName(SubproviderName), cacheKey, akamaiBotCategoryList)
	if err == nil {
		return filterAkamaiBotCategoryList(akamaiBotCategoryList, request), nil
	}

	if !errors.Is(err, cache.ErrEntryNotFound) && !errors.Is(err, cache.ErrDisabled) {
		logger.Errorf("error reading from cache: %s", err.Error())
		return nil, err
	}

	// fetch all akamaiBotCategoryList to store in cache and then filter based on request
	akamaiBotCategoryList, err = client.GetAkamaiBotCategoryList(ctx, botman.GetAkamaiBotCategoryListRequest{})
	if err != nil {
		logger.Errorf("calling 'GetAkamaiBotCategoryList': %s", err.Error())
		return nil, err
	}

	err = cache.Set(cache.BucketName(SubproviderName), cacheKey, akamaiBotCategoryList)
	if err != nil && !errors.Is(err, cache.ErrDisabled) {
		logger.Errorf("error caching akamaiBotCategoryList into cache: %s", err.Error())
		return nil, err
	}

	return filterAkamaiBotCategoryList(akamaiBotCategoryList, request), nil
}

func filterAkamaiBotCategoryList(akamaiBotCategoryList *botman.GetAkamaiBotCategoryListResponse, request botman.GetAkamaiBotCategoryListRequest) *botman.GetAkamaiBotCategoryListResponse {
	if request.CategoryName == "" {
		return akamaiBotCategoryList
	}
	var filteredResult botman.GetAkamaiBotCategoryListResponse
	for _, category := range akamaiBotCategoryList.Categories {
		if category["categoryName"].(string) == request.CategoryName {
			filteredResult.Categories = append(filteredResult.Categories, category)
			return &filteredResult
		}
	}
	return &filteredResult
}

func getAkamaiDefinedBotList(ctx context.Context, request botman.GetAkamaiDefinedBotListRequest, m interface{}) (*botman.GetAkamaiDefinedBotListResponse, error) {
	meta := akameta.Must(m)
	client := inst.Client(meta)
	logger := meta.Log("BotMan", "getAkamaiDefinedBot")

	cacheKey := "getAkamaiDefinedBot"
	akamaiDefinedBotList := &botman.GetAkamaiDefinedBotListResponse{}
	err := cache.Get(cache.BucketName(SubproviderName), cacheKey, akamaiDefinedBotList)
	// if cache is disabled make a direct all to GetAkamaiDefinedBotList
	if errors.Is(err, cache.ErrDisabled) {
		return client.GetAkamaiDefinedBotList(ctx, request)
	}
	if err == nil {
		return filterAkamaiDefinedBotList(akamaiDefinedBotList, request), nil
	}

	akamaiDefinedBotMutex.Lock()
	defer func() {
		logger.Debugf("Unlocking mutex")
		akamaiDefinedBotMutex.Unlock()
	}()

	err = cache.Get(cache.BucketName(SubproviderName), cacheKey, akamaiDefinedBotList)
	if err == nil {
		return filterAkamaiDefinedBotList(akamaiDefinedBotList, request), nil
	}

	if !errors.Is(err, cache.ErrEntryNotFound) && !errors.Is(err, cache.ErrDisabled) {
		logger.Errorf("error reading from cache: %s", err.Error())
		return nil, err
	}

	// fetch all akamaiDefinedBotList to store in cache and then filter based on request
	akamaiDefinedBotList, err = client.GetAkamaiDefinedBotList(ctx, botman.GetAkamaiDefinedBotListRequest{})
	if err != nil {
		logger.Errorf("calling 'GetAkamaiDefinedBotList': %s", err.Error())
		return nil, err
	}

	err = cache.Set(cache.BucketName(SubproviderName), cacheKey, akamaiDefinedBotList)
	if err != nil && !errors.Is(err, cache.ErrDisabled) {
		logger.Errorf("error caching akamaiDefinedBotList into cache: %s", err.Error())
		return nil, err
	}

	return filterAkamaiDefinedBotList(akamaiDefinedBotList, request), nil
}

func filterAkamaiDefinedBotList(akamaiDefinedBotList *botman.GetAkamaiDefinedBotListResponse, request botman.GetAkamaiDefinedBotListRequest) *botman.GetAkamaiDefinedBotListResponse {
	if request.BotName == "" {
		return akamaiDefinedBotList
	}
	var filteredResult botman.GetAkamaiDefinedBotListResponse
	for _, bot := range akamaiDefinedBotList.Bots {
		if bot["botName"].(string) == request.BotName {
			filteredResult.Bots = append(filteredResult.Bots, bot)
			return &filteredResult
		}
	}
	return &filteredResult
}
func getBotDetectionList(ctx context.Context, request botman.GetBotDetectionListRequest, m interface{}) (*botman.GetBotDetectionListResponse, error) {
	meta := akameta.Must(m)
	client := inst.Client(meta)
	logger := meta.Log("BotMan", "getBotDetection")

	cacheKey := "getBotDetection"
	botDetectionList := &botman.GetBotDetectionListResponse{}
	err := cache.Get(cache.BucketName(SubproviderName), cacheKey, botDetectionList)
	// if cache is disabled make a direct all to GetBotDetectionList
	if errors.Is(err, cache.ErrDisabled) {
		return client.GetBotDetectionList(ctx, request)
	}
	if err == nil {
		return filterBotDetectionList(botDetectionList, request), nil
	}

	botDetectionMutex.Lock()
	defer func() {
		logger.Debugf("Unlocking mutex")
		botDetectionMutex.Unlock()
	}()

	err = cache.Get(cache.BucketName(SubproviderName), cacheKey, botDetectionList)
	if err == nil {
		return filterBotDetectionList(botDetectionList, request), nil
	}

	if !errors.Is(err, cache.ErrEntryNotFound) && !errors.Is(err, cache.ErrDisabled) {
		logger.Errorf("error reading from cache: %s", err.Error())
		return nil, err
	}

	// fetch all botDetectionList to store in cache and then filter based on request
	botDetectionList, err = client.GetBotDetectionList(ctx, botman.GetBotDetectionListRequest{})
	if err != nil {
		logger.Errorf("calling 'GetBotDetectionList': %s", err.Error())
		return nil, err
	}

	err = cache.Set(cache.BucketName(SubproviderName), cacheKey, botDetectionList)
	if err != nil && !errors.Is(err, cache.ErrDisabled) {
		logger.Errorf("error caching botDetectionList into cache: %s", err.Error())
		return nil, err
	}

	return filterBotDetectionList(botDetectionList, request), nil
}

func filterBotDetectionList(botDetectionList *botman.GetBotDetectionListResponse, request botman.GetBotDetectionListRequest) *botman.GetBotDetectionListResponse {
	if request.DetectionName == "" {
		return botDetectionList
	}
	var filteredResult botman.GetBotDetectionListResponse
	for _, detection := range botDetectionList.Detections {
		if detection["detectionName"].(string) == request.DetectionName {
			filteredResult.Detections = append(filteredResult.Detections, detection)
			return &filteredResult
		}
	}
	return &filteredResult
}

// getContentProtectionRule reads from the cache if present, or makes a getAll call to fetch all Content Protection Rules for a security policy, stores in the cache and filters the required Content Protection Rule using ID.
func getContentProtectionRule(ctx context.Context, request botman.GetContentProtectionRuleRequest, m interface{}) (map[string]interface{}, error) {
	meta := akameta.Must(m)
	client := inst.Client(meta)
	logger := meta.Log("BotMan", "getContentProtectionRule")

	cacheKey := fmt.Sprintf("%s:%d:%d:%s", "getContentProtectionRule", request.ConfigID, request.Version, request.SecurityPolicyID)
	contentProtectionRules := &botman.GetContentProtectionRuleListResponse{}
	err := cache.Get(cache.BucketName(SubproviderName), cacheKey, contentProtectionRules)
	// if cache is disabled use GetTransactionalEndpoint to fetch one action at a time
	if errors.Is(err, cache.ErrDisabled) {
		return client.GetContentProtectionRule(ctx, request)
	}
	if err == nil {
		return filterContentProtectionRule(contentProtectionRules, request, logger)
	}

	contentProtectionRuleMutex.Lock()
	defer func() {
		logger.Debugf("Unlocking mutex")
		contentProtectionRuleMutex.Unlock()
	}()

	err = cache.Get(cache.BucketName(SubproviderName), cacheKey, contentProtectionRules)
	if err == nil {
		return filterContentProtectionRule(contentProtectionRules, request, logger)
	}

	if !errors.Is(err, cache.ErrEntryNotFound) && !errors.Is(err, cache.ErrDisabled) {
		logger.Errorf("error reading from cache: %s", err.Error())
		return nil, err
	}

	contentProtectionRules, err = client.GetContentProtectionRuleList(ctx, botman.GetContentProtectionRuleListRequest{
		ConfigID:         request.ConfigID,
		Version:          request.Version,
		SecurityPolicyID: request.SecurityPolicyID,
	})
	if err != nil {
		logger.Errorf("calling 'GetContentProtectionRuleListRequest': %s", err.Error())
		return nil, err
	}

	err = cache.Set(cache.BucketName(SubproviderName), cacheKey, contentProtectionRules)
	if err != nil && !errors.Is(err, cache.ErrDisabled) {
		logger.Errorf("error caching transactionalEndpoints into cache: %s", err.Error())
		return nil, err
	}

	return filterContentProtectionRule(contentProtectionRules, request, logger)
}

func filterContentProtectionRule(contentProtectionRules *botman.GetContentProtectionRuleListResponse, request botman.GetContentProtectionRuleRequest, logger log.Interface) (map[string]interface{}, error) {
	for _, rule := range contentProtectionRules.ContentProtectionRules {
		if rule["contentProtectionRuleId"].(string) == request.ContentProtectionRuleID {
			logger.Debugf("Found content protection rule for config %d version %d security policy %s contentProtectionRuleId %s", request.ConfigID, request.Version, request.SecurityPolicyID, request.ContentProtectionRuleID)
			return rule, nil
		}
	}
	return nil, fmt.Errorf("ContentProtectionRule with id: %s does not exist", request.ContentProtectionRuleID)
}

// getContentProtectionRule reads from the cache if present, or makes a getAll call to fetch all Content Protection Rules for a security policy, stores in the cache and filters the required Content Protection Rule using ID.
func getContentProtectionJavaScriptInjectionRule(ctx context.Context, request botman.GetContentProtectionJavaScriptInjectionRuleRequest, m interface{}) (map[string]interface{}, error) {
	meta := akameta.Must(m)
	client := inst.Client(meta)
	logger := meta.Log("BotMan", "getContentProtectionJavaScriptInjectionRule")

	cacheKey := fmt.Sprintf("%s:%d:%d:%s", "getContentProtectionJavaScriptInjectionRule", request.ConfigID, request.Version, request.SecurityPolicyID)
	contentProtectionJavaScriptInjectionRules := &botman.GetContentProtectionJavaScriptInjectionRuleListResponse{}
	err := cache.Get(cache.BucketName(SubproviderName), cacheKey, contentProtectionJavaScriptInjectionRules)
	// if cache is disabled use GetTransactionalEndpoint to fetch one action at a time
	if errors.Is(err, cache.ErrDisabled) {
		return client.GetContentProtectionJavaScriptInjectionRule(ctx, request)
	}
	if err == nil {
		return filterContentProtectionJavaScriptInjectionRule(contentProtectionJavaScriptInjectionRules, request, logger)
	}

	contentProtectionJavaScriptInjectionRuleMutex.Lock()
	defer func() {
		logger.Debugf("Unlocking mutex")
		contentProtectionJavaScriptInjectionRuleMutex.Unlock()
	}()

	err = cache.Get(cache.BucketName(SubproviderName), cacheKey, contentProtectionJavaScriptInjectionRules)
	if err == nil {
		return filterContentProtectionJavaScriptInjectionRule(contentProtectionJavaScriptInjectionRules, request, logger)
	}

	if !errors.Is(err, cache.ErrEntryNotFound) && !errors.Is(err, cache.ErrDisabled) {
		logger.Errorf("error reading from cache: %s", err.Error())
		return nil, err
	}

	contentProtectionJavaScriptInjectionRules, err = client.GetContentProtectionJavaScriptInjectionRuleList(ctx, botman.GetContentProtectionJavaScriptInjectionRuleListRequest{
		ConfigID:         request.ConfigID,
		Version:          request.Version,
		SecurityPolicyID: request.SecurityPolicyID,
	})
	if err != nil {
		logger.Errorf("calling 'GetContentProtectionJavaScriptInjectionRuleListRequest': %s", err.Error())
		return nil, err
	}

	err = cache.Set(cache.BucketName(SubproviderName), cacheKey, contentProtectionJavaScriptInjectionRules)
	if err != nil && !errors.Is(err, cache.ErrDisabled) {
		logger.Errorf("error caching transactionalEndpoints into cache: %s", err.Error())
		return nil, err
	}

	return filterContentProtectionJavaScriptInjectionRule(contentProtectionJavaScriptInjectionRules, request, logger)
}

func filterContentProtectionJavaScriptInjectionRule(contentProtectionJavaScriptInjectionRules *botman.GetContentProtectionJavaScriptInjectionRuleListResponse, request botman.GetContentProtectionJavaScriptInjectionRuleRequest, logger log.Interface) (map[string]interface{}, error) {
	for _, rule := range contentProtectionJavaScriptInjectionRules.ContentProtectionJavaScriptInjectionRules {
		if rule["contentProtectionJavaScriptInjectionRuleId"].(string) == request.ContentProtectionJavaScriptInjectionRuleID {
			logger.Debugf("Found content protection javascript injection rule for config %d version %d security policy %s contentProtectionJavaScriptInjectionRuleId %s", request.ConfigID, request.Version, request.SecurityPolicyID, request.ContentProtectionJavaScriptInjectionRuleID)
			return rule, nil
		}
	}
	return nil, fmt.Errorf("ContentProtectionJavaScriptInjectionRule with id: %s does not exist", request.ContentProtectionJavaScriptInjectionRuleID)
}
