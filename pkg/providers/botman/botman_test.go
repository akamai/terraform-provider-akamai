package botman

import (
	"context"

	"github.com/akamai/AkamaiOPEN-edgegrid-golang/v3/pkg/botman"
	"github.com/stretchr/testify/mock"
)

type mockbotman struct {
	mock.Mock
}

func (p *mockbotman) GetAkamaiBotCategoryList(ctx context.Context, params botman.GetAkamaiBotCategoryListRequest) (*botman.GetAkamaiBotCategoryListResponse, error) {
	args := p.Called(ctx, params)
	if args.Error(1) != nil {
		return nil, args.Error(1)
	}

	return args.Get(0).(*botman.GetAkamaiBotCategoryListResponse), nil
}

func (p *mockbotman) GetAkamaiBotCategoryActionList(ctx context.Context, params botman.GetAkamaiBotCategoryActionListRequest) (*botman.GetAkamaiBotCategoryActionListResponse, error) {
	args := p.Called(ctx, params)
	if args.Error(1) != nil {
		return nil, args.Error(1)
	}

	return args.Get(0).(*botman.GetAkamaiBotCategoryActionListResponse), nil
}
func (p *mockbotman) GetAkamaiBotCategoryAction(ctx context.Context, params botman.GetAkamaiBotCategoryActionRequest) (map[string]interface{}, error) {
	args := p.Called(ctx, params)
	if args.Error(1) != nil {
		return nil, args.Error(1)
	}

	return args.Get(0).(map[string]interface{}), nil
}
func (p *mockbotman) UpdateAkamaiBotCategoryAction(ctx context.Context, params botman.UpdateAkamaiBotCategoryActionRequest) (map[string]interface{}, error) {
	args := p.Called(ctx, params)
	if args.Error(1) != nil {
		return nil, args.Error(1)
	}

	return args.Get(0).(map[string]interface{}), nil
}

func (p *mockbotman) GetAkamaiDefinedBotList(ctx context.Context, params botman.GetAkamaiDefinedBotListRequest) (*botman.GetAkamaiDefinedBotListResponse, error) {
	args := p.Called(ctx, params)
	if args.Error(1) != nil {
		return nil, args.Error(1)
	}

	return args.Get(0).(*botman.GetAkamaiDefinedBotListResponse), nil
}
func (p *mockbotman) GetBotAnalyticsCookie(ctx context.Context, params botman.GetBotAnalyticsCookieRequest) (map[string]interface{}, error) {
	args := p.Called(ctx, params)
	if args.Error(1) != nil {
		return nil, args.Error(1)
	}
	return args.Get(0).(map[string]interface{}), nil
}
func (p *mockbotman) UpdateBotAnalyticsCookie(ctx context.Context, params botman.UpdateBotAnalyticsCookieRequest) (map[string]interface{}, error) {
	args := p.Called(ctx, params)
	if args.Error(1) != nil {
		return nil, args.Error(1)
	}
	return args.Get(0).(map[string]interface{}), nil
}
func (p *mockbotman) GetBotAnalyticsCookieValues(ctx context.Context) (map[string]interface{}, error) {
	args := p.Called(ctx)
	if args.Error(1) != nil {
		return nil, args.Error(1)
	}
	return args.Get(0).(map[string]interface{}), nil
}
func (p *mockbotman) GetBotCategoryException(ctx context.Context, params botman.GetBotCategoryExceptionRequest) (map[string]interface{}, error) {
	args := p.Called(ctx, params)
	if args.Error(1) != nil {
		return nil, args.Error(1)
	}
	return args.Get(0).(map[string]interface{}), nil
}
func (p *mockbotman) UpdateBotCategoryException(ctx context.Context, params botman.UpdateBotCategoryExceptionRequest) (map[string]interface{}, error) {
	args := p.Called(ctx, params)
	if args.Error(1) != nil {
		return nil, args.Error(1)
	}
	return args.Get(0).(map[string]interface{}), nil
}
func (p *mockbotman) GetBotDetectionActionList(ctx context.Context, params botman.GetBotDetectionActionListRequest) (*botman.GetBotDetectionActionListResponse, error) {
	args := p.Called(ctx, params)
	if args.Error(1) != nil {
		return nil, args.Error(1)
	}
	return args.Get(0).(*botman.GetBotDetectionActionListResponse), nil
}
func (p *mockbotman) GetBotDetectionAction(ctx context.Context, params botman.GetBotDetectionActionRequest) (map[string]interface{}, error) {
	args := p.Called(ctx, params)
	if args.Error(1) != nil {
		return nil, args.Error(1)
	}
	return args.Get(0).(map[string]interface{}), nil
}
func (p *mockbotman) UpdateBotDetectionAction(ctx context.Context, params botman.UpdateBotDetectionActionRequest) (map[string]interface{}, error) {
	args := p.Called(ctx, params)
	if args.Error(1) != nil {
		return nil, args.Error(1)
	}
	return args.Get(0).(map[string]interface{}), nil
}
func (p *mockbotman) GetBotEndpointCoverageReport(ctx context.Context, params botman.GetBotEndpointCoverageReportRequest) (*botman.GetBotEndpointCoverageReportResponse, error) {
	args := p.Called(ctx, params)
	if args.Error(1) != nil {
		return nil, args.Error(1)
	}
	return args.Get(0).(*botman.GetBotEndpointCoverageReportResponse), nil
}
func (p *mockbotman) GetBotManagementSetting(ctx context.Context, params botman.GetBotManagementSettingRequest) (map[string]interface{}, error) {
	args := p.Called(ctx, params)
	if args.Error(1) != nil {
		return nil, args.Error(1)
	}
	return args.Get(0).(map[string]interface{}), nil
}
func (p *mockbotman) UpdateBotManagementSetting(ctx context.Context, params botman.UpdateBotManagementSettingRequest) (map[string]interface{}, error) {
	args := p.Called(ctx, params)
	if args.Error(1) != nil {
		return nil, args.Error(1)
	}
	return args.Get(0).(map[string]interface{}), nil
}
func (p *mockbotman) GetChallengeActionList(ctx context.Context, params botman.GetChallengeActionListRequest) (*botman.GetChallengeActionListResponse, error) {
	args := p.Called(ctx, params)
	if args.Error(1) != nil {
		return nil, args.Error(1)
	}
	return args.Get(0).(*botman.GetChallengeActionListResponse), nil
}
func (p *mockbotman) GetChallengeAction(ctx context.Context, params botman.GetChallengeActionRequest) (map[string]interface{}, error) {
	args := p.Called(ctx, params)
	if args.Error(1) != nil {
		return nil, args.Error(1)
	}
	return args.Get(0).(map[string]interface{}), nil
}
func (p *mockbotman) CreateChallengeAction(ctx context.Context, params botman.CreateChallengeActionRequest) (map[string]interface{}, error) {
	args := p.Called(ctx, params)
	if args.Error(1) != nil {
		return nil, args.Error(1)
	}
	return args.Get(0).(map[string]interface{}), nil
}
func (p *mockbotman) UpdateChallengeAction(ctx context.Context, params botman.UpdateChallengeActionRequest) (map[string]interface{}, error) {
	args := p.Called(ctx, params)
	if args.Error(1) != nil {
		return nil, args.Error(1)
	}
	return args.Get(0).(map[string]interface{}), nil
}
func (p *mockbotman) RemoveChallengeAction(ctx context.Context, params botman.RemoveChallengeActionRequest) error {
	args := p.Called(ctx, params)
	return args.Error(0)
}
func (p *mockbotman) UpdateGoogleReCaptchaSecretKey(ctx context.Context, params botman.UpdateGoogleReCaptchaSecretKeyRequest) error {
	args := p.Called(ctx, params)
	return args.Error(0)
}
func (p *mockbotman) GetChallengeInterceptionRules(ctx context.Context, params botman.GetChallengeInterceptionRulesRequest) (map[string]interface{}, error) {
	args := p.Called(ctx, params)
	if args.Error(1) != nil {
		return nil, args.Error(1)
	}
	return args.Get(0).(map[string]interface{}), nil
}
func (p *mockbotman) UpdateChallengeInterceptionRules(ctx context.Context, params botman.UpdateChallengeInterceptionRulesRequest) (map[string]interface{}, error) {
	args := p.Called(ctx, params)
	if args.Error(1) != nil {
		return nil, args.Error(1)
	}
	return args.Get(0).(map[string]interface{}), nil
}
func (p *mockbotman) GetClientSideSecurity(ctx context.Context, params botman.GetClientSideSecurityRequest) (map[string]interface{}, error) {
	args := p.Called(ctx, params)
	if args.Error(1) != nil {
		return nil, args.Error(1)
	}
	return args.Get(0).(map[string]interface{}), nil
}
func (p *mockbotman) UpdateClientSideSecurity(ctx context.Context, params botman.UpdateClientSideSecurityRequest) (map[string]interface{}, error) {
	args := p.Called(ctx, params)
	if args.Error(1) != nil {
		return nil, args.Error(1)
	}
	return args.Get(0).(map[string]interface{}), nil
}
func (p *mockbotman) GetConditionalActionList(ctx context.Context, params botman.GetConditionalActionListRequest) (*botman.GetConditionalActionListResponse, error) {
	args := p.Called(ctx, params)
	if args.Error(1) != nil {
		return nil, args.Error(1)
	}
	return args.Get(0).(*botman.GetConditionalActionListResponse), nil
}
func (p *mockbotman) GetConditionalAction(ctx context.Context, params botman.GetConditionalActionRequest) (map[string]interface{}, error) {
	args := p.Called(ctx, params)
	if args.Error(1) != nil {
		return nil, args.Error(1)
	}
	return args.Get(0).(map[string]interface{}), nil
}
func (p *mockbotman) CreateConditionalAction(ctx context.Context, params botman.CreateConditionalActionRequest) (map[string]interface{}, error) {
	args := p.Called(ctx, params)
	if args.Error(1) != nil {
		return nil, args.Error(1)
	}
	return args.Get(0).(map[string]interface{}), nil
}
func (p *mockbotman) UpdateConditionalAction(ctx context.Context, params botman.UpdateConditionalActionRequest) (map[string]interface{}, error) {
	args := p.Called(ctx, params)
	if args.Error(1) != nil {
		return nil, args.Error(1)
	}
	return args.Get(0).(map[string]interface{}), nil
}
func (p *mockbotman) RemoveConditionalAction(ctx context.Context, params botman.RemoveConditionalActionRequest) error {
	args := p.Called(ctx, params)
	return args.Error(0)
}
func (p *mockbotman) GetCustomBotCategoryList(ctx context.Context, params botman.GetCustomBotCategoryListRequest) (*botman.GetCustomBotCategoryListResponse, error) {
	args := p.Called(ctx, params)
	if args.Error(1) != nil {
		return nil, args.Error(1)
	}
	return args.Get(0).(*botman.GetCustomBotCategoryListResponse), nil
}
func (p *mockbotman) GetCustomBotCategory(ctx context.Context, params botman.GetCustomBotCategoryRequest) (map[string]interface{}, error) {
	args := p.Called(ctx, params)
	if args.Error(1) != nil {
		return nil, args.Error(1)
	}
	return args.Get(0).(map[string]interface{}), nil
}
func (p *mockbotman) CreateCustomBotCategory(ctx context.Context, params botman.CreateCustomBotCategoryRequest) (map[string]interface{}, error) {
	args := p.Called(ctx, params)
	if args.Error(1) != nil {
		return nil, args.Error(1)
	}
	return args.Get(0).(map[string]interface{}), nil
}
func (p *mockbotman) UpdateCustomBotCategory(ctx context.Context, params botman.UpdateCustomBotCategoryRequest) (map[string]interface{}, error) {
	args := p.Called(ctx, params)
	if args.Error(1) != nil {
		return nil, args.Error(1)
	}
	return args.Get(0).(map[string]interface{}), nil
}
func (p *mockbotman) RemoveCustomBotCategory(ctx context.Context, params botman.RemoveCustomBotCategoryRequest) error {
	args := p.Called(ctx, params)
	return args.Error(0)
}
func (p *mockbotman) GetCustomBotCategoryActionList(ctx context.Context, params botman.GetCustomBotCategoryActionListRequest) (*botman.GetCustomBotCategoryActionListResponse, error) {
	args := p.Called(ctx, params)
	if args.Error(1) != nil {
		return nil, args.Error(1)
	}
	return args.Get(0).(*botman.GetCustomBotCategoryActionListResponse), nil
}
func (p *mockbotman) GetCustomBotCategoryAction(ctx context.Context, params botman.GetCustomBotCategoryActionRequest) (map[string]interface{}, error) {
	args := p.Called(ctx, params)
	if args.Error(1) != nil {
		return nil, args.Error(1)
	}
	return args.Get(0).(map[string]interface{}), nil
}
func (p *mockbotman) UpdateCustomBotCategoryAction(ctx context.Context, params botman.UpdateCustomBotCategoryActionRequest) (map[string]interface{}, error) {
	args := p.Called(ctx, params)
	if args.Error(1) != nil {
		return nil, args.Error(1)
	}
	return args.Get(0).(map[string]interface{}), nil
}
func (p *mockbotman) GetCustomBotCategorySequence(ctx context.Context, params botman.GetCustomBotCategorySequenceRequest) (*botman.CustomBotCategorySequenceResponse, error) {
	args := p.Called(ctx, params)
	if args.Error(1) != nil {
		return nil, args.Error(1)
	}
	return args.Get(0).(*botman.CustomBotCategorySequenceResponse), nil
}
func (p *mockbotman) UpdateCustomBotCategorySequence(ctx context.Context, params botman.UpdateCustomBotCategorySequenceRequest) (*botman.CustomBotCategorySequenceResponse, error) {
	args := p.Called(ctx, params)
	if args.Error(1) != nil {
		return nil, args.Error(1)
	}
	return args.Get(0).(*botman.CustomBotCategorySequenceResponse), nil
}
func (p *mockbotman) GetCustomClientList(ctx context.Context, params botman.GetCustomClientListRequest) (*botman.GetCustomClientListResponse, error) {
	args := p.Called(ctx, params)
	if args.Error(1) != nil {
		return nil, args.Error(1)
	}
	return args.Get(0).(*botman.GetCustomClientListResponse), nil
}
func (p *mockbotman) GetCustomClient(ctx context.Context, params botman.GetCustomClientRequest) (map[string]interface{}, error) {
	args := p.Called(ctx, params)
	if args.Error(1) != nil {
		return nil, args.Error(1)
	}
	return args.Get(0).(map[string]interface{}), nil
}
func (p *mockbotman) CreateCustomClient(ctx context.Context, params botman.CreateCustomClientRequest) (map[string]interface{}, error) {
	args := p.Called(ctx, params)
	if args.Error(1) != nil {
		return nil, args.Error(1)
	}
	return args.Get(0).(map[string]interface{}), nil
}
func (p *mockbotman) UpdateCustomClient(ctx context.Context, params botman.UpdateCustomClientRequest) (map[string]interface{}, error) {
	args := p.Called(ctx, params)
	if args.Error(1) != nil {
		return nil, args.Error(1)
	}
	return args.Get(0).(map[string]interface{}), nil
}
func (p *mockbotman) RemoveCustomClient(ctx context.Context, params botman.RemoveCustomClientRequest) error {
	args := p.Called(ctx, params)
	return args.Error(0)
}
func (p *mockbotman) GetCustomDefinedBotList(ctx context.Context, params botman.GetCustomDefinedBotListRequest) (*botman.GetCustomDefinedBotListResponse, error) {
	args := p.Called(ctx, params)
	if args.Error(1) != nil {
		return nil, args.Error(1)
	}
	return args.Get(0).(*botman.GetCustomDefinedBotListResponse), nil
}
func (p *mockbotman) GetCustomDefinedBot(ctx context.Context, params botman.GetCustomDefinedBotRequest) (map[string]interface{}, error) {
	args := p.Called(ctx, params)
	if args.Error(1) != nil {
		return nil, args.Error(1)
	}
	return args.Get(0).(map[string]interface{}), nil
}
func (p *mockbotman) CreateCustomDefinedBot(ctx context.Context, params botman.CreateCustomDefinedBotRequest) (map[string]interface{}, error) {
	args := p.Called(ctx, params)
	if args.Error(1) != nil {
		return nil, args.Error(1)
	}
	return args.Get(0).(map[string]interface{}), nil
}
func (p *mockbotman) UpdateCustomDefinedBot(ctx context.Context, params botman.UpdateCustomDefinedBotRequest) (map[string]interface{}, error) {
	args := p.Called(ctx, params)
	if args.Error(1) != nil {
		return nil, args.Error(1)
	}
	return args.Get(0).(map[string]interface{}), nil
}
func (p *mockbotman) RemoveCustomDefinedBot(ctx context.Context, params botman.RemoveCustomDefinedBotRequest) error {
	args := p.Called(ctx, params)
	return args.Error(0)
}
func (p *mockbotman) GetCustomDenyActionList(ctx context.Context, params botman.GetCustomDenyActionListRequest) (*botman.GetCustomDenyActionListResponse, error) {
	args := p.Called(ctx, params)
	if args.Error(1) != nil {
		return nil, args.Error(1)
	}
	return args.Get(0).(*botman.GetCustomDenyActionListResponse), nil
}
func (p *mockbotman) GetCustomDenyAction(ctx context.Context, params botman.GetCustomDenyActionRequest) (map[string]interface{}, error) {
	args := p.Called(ctx, params)
	if args.Error(1) != nil {
		return nil, args.Error(1)
	}
	return args.Get(0).(map[string]interface{}), nil
}
func (p *mockbotman) CreateCustomDenyAction(ctx context.Context, params botman.CreateCustomDenyActionRequest) (map[string]interface{}, error) {
	args := p.Called(ctx, params)
	if args.Error(1) != nil {
		return nil, args.Error(1)
	}
	return args.Get(0).(map[string]interface{}), nil
}
func (p *mockbotman) UpdateCustomDenyAction(ctx context.Context, params botman.UpdateCustomDenyActionRequest) (map[string]interface{}, error) {
	args := p.Called(ctx, params)
	if args.Error(1) != nil {
		return nil, args.Error(1)
	}
	return args.Get(0).(map[string]interface{}), nil
}
func (p *mockbotman) RemoveCustomDenyAction(ctx context.Context, params botman.RemoveCustomDenyActionRequest) error {
	args := p.Called(ctx, params)
	return args.Error(0)
}

func (p *mockbotman) GetRecategorizedAkamaiDefinedBotList(ctx context.Context, params botman.GetRecategorizedAkamaiDefinedBotListRequest) (*botman.GetRecategorizedAkamaiDefinedBotListResponse, error) {
	args := p.Called(ctx, params)
	if args.Error(1) != nil {
		return nil, args.Error(1)
	}
	return args.Get(0).(*botman.GetRecategorizedAkamaiDefinedBotListResponse), nil
}
func (p *mockbotman) GetRecategorizedAkamaiDefinedBot(ctx context.Context, params botman.GetRecategorizedAkamaiDefinedBotRequest) (*botman.RecategorizedAkamaiDefinedBotResponse, error) {
	args := p.Called(ctx, params)
	if args.Error(1) != nil {
		return nil, args.Error(1)
	}
	return args.Get(0).(*botman.RecategorizedAkamaiDefinedBotResponse), nil
}
func (p *mockbotman) CreateRecategorizedAkamaiDefinedBot(ctx context.Context, params botman.CreateRecategorizedAkamaiDefinedBotRequest) (*botman.RecategorizedAkamaiDefinedBotResponse, error) {
	args := p.Called(ctx, params)
	if args.Error(1) != nil {
		return nil, args.Error(1)
	}
	return args.Get(0).(*botman.RecategorizedAkamaiDefinedBotResponse), nil
}
func (p *mockbotman) UpdateRecategorizedAkamaiDefinedBot(ctx context.Context, params botman.UpdateRecategorizedAkamaiDefinedBotRequest) (*botman.RecategorizedAkamaiDefinedBotResponse, error) {
	args := p.Called(ctx, params)
	if args.Error(1) != nil {
		return nil, args.Error(1)
	}
	return args.Get(0).(*botman.RecategorizedAkamaiDefinedBotResponse), nil
}
func (p *mockbotman) RemoveRecategorizedAkamaiDefinedBot(ctx context.Context, params botman.RemoveRecategorizedAkamaiDefinedBotRequest) error {
	args := p.Called(ctx, params)
	return args.Error(0)
}
func (p *mockbotman) GetResponseActionList(ctx context.Context, params botman.GetResponseActionListRequest) (*botman.GetResponseActionListResponse, error) {
	args := p.Called(ctx, params)
	if args.Error(1) != nil {
		return nil, args.Error(1)
	}
	return args.Get(0).(*botman.GetResponseActionListResponse), nil
}

func (p *mockbotman) GetTransactionalEndpointList(ctx context.Context, params botman.GetTransactionalEndpointListRequest) (*botman.GetTransactionalEndpointListResponse, error) {
	args := p.Called(ctx, params)
	if args.Error(1) != nil {
		return nil, args.Error(1)
	}
	return args.Get(0).(*botman.GetTransactionalEndpointListResponse), nil
}
func (p *mockbotman) GetTransactionalEndpoint(ctx context.Context, params botman.GetTransactionalEndpointRequest) (map[string]interface{}, error) {
	args := p.Called(ctx, params)
	if args.Error(1) != nil {
		return nil, args.Error(1)
	}
	return args.Get(0).(map[string]interface{}), nil
}
func (p *mockbotman) CreateTransactionalEndpoint(ctx context.Context, params botman.CreateTransactionalEndpointRequest) (map[string]interface{}, error) {
	args := p.Called(ctx, params)
	if args.Error(1) != nil {
		return nil, args.Error(1)
	}
	return args.Get(0).(map[string]interface{}), nil
}
func (p *mockbotman) UpdateTransactionalEndpoint(ctx context.Context, params botman.UpdateTransactionalEndpointRequest) (map[string]interface{}, error) {
	args := p.Called(ctx, params)
	if args.Error(1) != nil {
		return nil, args.Error(1)
	}
	return args.Get(0).(map[string]interface{}), nil
}
func (p *mockbotman) RemoveTransactionalEndpoint(ctx context.Context, params botman.RemoveTransactionalEndpointRequest) error {
	args := p.Called(ctx, params)
	return args.Error(0)
}
func (p *mockbotman) GetTransactionalEndpointProtection(ctx context.Context, params botman.GetTransactionalEndpointProtectionRequest) (map[string]interface{}, error) {
	args := p.Called(ctx, params)
	if args.Error(1) != nil {
		return nil, args.Error(1)
	}
	return args.Get(0).(map[string]interface{}), nil
}
func (p *mockbotman) UpdateTransactionalEndpointProtection(ctx context.Context, params botman.UpdateTransactionalEndpointProtectionRequest) (map[string]interface{}, error) {
	args := p.Called(ctx, params)
	if args.Error(1) != nil {
		return nil, args.Error(1)
	}
	return args.Get(0).(map[string]interface{}), nil
}

func (p *mockbotman) GetJavascriptInjection(ctx context.Context, params botman.GetJavascriptInjectionRequest) (map[string]interface{}, error) {
	args := p.Called(ctx, params)
	if args.Error(1) != nil {
		return nil, args.Error(1)
	}
	return args.Get(0).(map[string]interface{}), nil
}
func (p *mockbotman) UpdateJavascriptInjection(ctx context.Context, params botman.UpdateJavascriptInjectionRequest) (map[string]interface{}, error) {
	args := p.Called(ctx, params)
	if args.Error(1) != nil {
		return nil, args.Error(1)
	}
	return args.Get(0).(map[string]interface{}), nil
}
func (p *mockbotman) GetServeAlternateActionList(ctx context.Context, params botman.GetServeAlternateActionListRequest) (*botman.GetServeAlternateActionListResponse, error) {
	args := p.Called(ctx, params)
	if args.Error(1) != nil {
		return nil, args.Error(1)
	}
	return args.Get(0).(*botman.GetServeAlternateActionListResponse), nil
}
func (p *mockbotman) GetServeAlternateAction(ctx context.Context, params botman.GetServeAlternateActionRequest) (map[string]interface{}, error) {
	args := p.Called(ctx, params)
	if args.Error(1) != nil {
		return nil, args.Error(1)
	}
	return args.Get(0).(map[string]interface{}), nil
}
func (p *mockbotman) CreateServeAlternateAction(ctx context.Context, params botman.CreateServeAlternateActionRequest) (map[string]interface{}, error) {
	args := p.Called(ctx, params)
	if args.Error(1) != nil {
		return nil, args.Error(1)
	}
	return args.Get(0).(map[string]interface{}), nil
}
func (p *mockbotman) UpdateServeAlternateAction(ctx context.Context, params botman.UpdateServeAlternateActionRequest) (map[string]interface{}, error) {
	args := p.Called(ctx, params)
	if args.Error(1) != nil {
		return nil, args.Error(1)
	}
	return args.Get(0).(map[string]interface{}), nil
}
func (p *mockbotman) RemoveServeAlternateAction(ctx context.Context, params botman.RemoveServeAlternateActionRequest) error {
	args := p.Called(ctx, params)
	return args.Error(0)
}

func (p *mockbotman) GetBotDetectionList(ctx context.Context, params botman.GetBotDetectionListRequest) (*botman.GetBotDetectionListResponse, error) {
	args := p.Called(ctx, params)
	if args.Error(1) != nil {
		return nil, args.Error(1)
	}
	return args.Get(0).(*botman.GetBotDetectionListResponse), nil
}
