package appsec

import (
	"context"

	"github.com/akamai/AkamaiOPEN-edgegrid-golang/v2/pkg/appsec"
	"github.com/stretchr/testify/mock"
)

type mockappsec struct {
	mock.Mock
}

func (p *mockappsec) GetConfigurations(ctx context.Context, params appsec.GetConfigurationsRequest) (*appsec.GetConfigurationsResponse, error) {
	args := p.Called(ctx, params)

	if args.Get(0) == nil {
		return nil, args.Error(1)
	}

	return args.Get(0).(*appsec.GetConfigurationsResponse), args.Error(1)
}

func (p *mockappsec) GetConfigurationVersions(ctx context.Context, params appsec.GetConfigurationVersionsRequest) (*appsec.GetConfigurationVersionsResponse, error) {
	args := p.Called(ctx, params)

	if args.Get(0) == nil {
		return nil, args.Error(1)
	}

	return args.Get(0).(*appsec.GetConfigurationVersionsResponse), args.Error(1)
}

func (p *mockappsec) CreateActivations(ctx context.Context, params appsec.CreateActivationsRequest, acknowledgeWarnings bool) (*appsec.CreateActivationsResponse, error) {
	args := p.Called(ctx, params)

	if args.Get(0) == nil {
		return nil, args.Error(1)
	}

	return args.Get(0).(*appsec.CreateActivationsResponse), args.Error(1)
}

func (p *mockappsec) GetActivations(ctx context.Context, params appsec.GetActivationsRequest) (*appsec.GetActivationsResponse, error) {
	args := p.Called(ctx, params)

	if args.Get(0) == nil {
		return nil, args.Error(1)
	}

	return args.Get(0).(*appsec.GetActivationsResponse), args.Error(1)
}

func (p *mockappsec) RemoveActivations(ctx context.Context, params appsec.RemoveActivationsRequest) (*appsec.RemoveActivationsResponse, error) {
	args := p.Called(ctx, params)

	if args.Get(0) == nil {
		return nil, args.Error(1)
	}

	return args.Get(0).(*appsec.RemoveActivationsResponse), args.Error(1)
}

func (p *mockappsec) CreateConfigurationClone(ctx context.Context, params appsec.CreateConfigurationCloneRequest) (*appsec.CreateConfigurationCloneResponse, error) {
	args := p.Called(ctx, params)

	if args.Get(0) == nil {
		return nil, args.Error(1)
	}

	return args.Get(0).(*appsec.CreateConfigurationCloneResponse), args.Error(1)
}

func (p *mockappsec) GetConfigurationClone(ctx context.Context, params appsec.GetConfigurationCloneRequest) (*appsec.GetConfigurationCloneResponse, error) {
	args := p.Called(ctx, params)

	if args.Get(0) == nil {
		return nil, args.Error(1)
	}

	return args.Get(0).(*appsec.GetConfigurationCloneResponse), args.Error(1)
}

func (p *mockappsec) CreateCustomRule(ctx context.Context, params appsec.CreateCustomRuleRequest) (*appsec.CreateCustomRuleResponse, error) {
	args := p.Called(ctx, params)

	if args.Get(0) == nil {
		return nil, args.Error(1)
	}

	return args.Get(0).(*appsec.CreateCustomRuleResponse), args.Error(1)
}

func (p *mockappsec) RemoveCustomRule(ctx context.Context, params appsec.RemoveCustomRuleRequest) (*appsec.RemoveCustomRuleResponse, error) {
	args := p.Called(ctx, params)

	if args.Get(0) == nil {
		return nil, args.Error(1)
	}

	return args.Get(0).(*appsec.RemoveCustomRuleResponse), args.Error(1)
}

func (p *mockappsec) UpdateCustomRule(ctx context.Context, params appsec.UpdateCustomRuleRequest) (*appsec.UpdateCustomRuleResponse, error) {
	args := p.Called(ctx, params)

	if args.Get(0) == nil {
		return nil, args.Error(1)
	}

	return args.Get(0).(*appsec.UpdateCustomRuleResponse), args.Error(1)
}

func (p *mockappsec) CreateMatchTarget(ctx context.Context, params appsec.CreateMatchTargetRequest) (*appsec.CreateMatchTargetResponse, error) {
	args := p.Called(ctx, params)

	if args.Get(0) == nil {
		return nil, args.Error(1)
	}

	return args.Get(0).(*appsec.CreateMatchTargetResponse), args.Error(1)
}

func (p *mockappsec) RemoveMatchTarget(ctx context.Context, params appsec.RemoveMatchTargetRequest) (*appsec.RemoveMatchTargetResponse, error) {
	args := p.Called(ctx, params)

	if args.Get(0) == nil {
		return nil, args.Error(1)
	}

	return args.Get(0).(*appsec.RemoveMatchTargetResponse), args.Error(1)
}

func (p *mockappsec) CreateRatePolicy(ctx context.Context, params appsec.CreateRatePolicyRequest) (*appsec.CreateRatePolicyResponse, error) {
	args := p.Called(ctx, params)

	if args.Get(0) == nil {
		return nil, args.Error(1)
	}

	return args.Get(0).(*appsec.CreateRatePolicyResponse), args.Error(1)
}

func (p *mockappsec) UpdateRatePolicy(ctx context.Context, params appsec.UpdateRatePolicyRequest) (*appsec.UpdateRatePolicyResponse, error) {
	args := p.Called(ctx, params)

	if args.Get(0) == nil {
		return nil, args.Error(1)
	}

	return args.Get(0).(*appsec.UpdateRatePolicyResponse), args.Error(1)
}

func (p *mockappsec) GetRatePolicy(ctx context.Context, params appsec.GetRatePolicyRequest) (*appsec.GetRatePolicyResponse, error) {
	args := p.Called(ctx, params)

	if args.Get(0) == nil {
		return nil, args.Error(1)
	}

	return args.Get(0).(*appsec.GetRatePolicyResponse), args.Error(1)
}

func (p *mockappsec) RemoveRatePolicy(ctx context.Context, params appsec.RemoveRatePolicyRequest) (*appsec.RemoveRatePolicyResponse, error) {
	args := p.Called(ctx, params)

	if args.Get(0) == nil {
		return nil, args.Error(1)
	}

	return args.Get(0).(*appsec.RemoveRatePolicyResponse), args.Error(1)
}

func (p *mockappsec) CreateRatePolicies(ctx context.Context, params appsec.CreateRatePolicyRequest) (*appsec.CreateRatePolicyResponse, error) {
	args := p.Called(ctx, params)

	if args.Get(0) == nil {
		return nil, args.Error(1)
	}

	return args.Get(0).(*appsec.CreateRatePolicyResponse), args.Error(1)
}

func (p *mockappsec) GetRatePolicies(ctx context.Context, params appsec.GetRatePoliciesRequest) (*appsec.GetRatePoliciesResponse, error) {
	args := p.Called(ctx, params)

	if args.Get(0) == nil {
		return nil, args.Error(1)
	}

	return args.Get(0).(*appsec.GetRatePoliciesResponse), args.Error(1)
}

func (p *mockappsec) GetRatePolicyAction(ctx context.Context, params appsec.GetRatePolicyActionRequest) (*appsec.GetRatePolicyActionResponse, error) {
	args := p.Called(ctx, params)

	if args.Get(0) == nil {
		return nil, args.Error(1)
	}

	return args.Get(0).(*appsec.GetRatePolicyActionResponse), args.Error(1)
}

func (p *mockappsec) UpdateRatePolicyAction(ctx context.Context, params appsec.UpdateRatePolicyActionRequest) (*appsec.UpdateRatePolicyActionResponse, error) {
	args := p.Called(ctx, params)

	if args.Get(0) == nil {
		return nil, args.Error(1)
	}

	return args.Get(0).(*appsec.UpdateRatePolicyActionResponse), args.Error(1)
}

func (p *mockappsec) GetRatePolicyActions(ctx context.Context, params appsec.GetRatePolicyActionsRequest) (*appsec.GetRatePolicyActionsResponse, error) {
	args := p.Called(ctx, params)

	if args.Get(0) == nil {
		return nil, args.Error(1)
	}

	return args.Get(0).(*appsec.GetRatePolicyActionsResponse), args.Error(1)
}

func (p *mockappsec) CreateSecurityPolicyClone(ctx context.Context, params appsec.CreateSecurityPolicyCloneRequest) (*appsec.CreateSecurityPolicyCloneResponse, error) {
	args := p.Called(ctx, params)

	if args.Get(0) == nil {
		return nil, args.Error(1)
	}

	return args.Get(0).(*appsec.CreateSecurityPolicyCloneResponse), args.Error(1)
}

func (p *mockappsec) GetSecurityPolicyClone(ctx context.Context, params appsec.GetSecurityPolicyCloneRequest) (*appsec.GetSecurityPolicyCloneResponse, error) {
	args := p.Called(ctx, params)

	if args.Get(0) == nil {
		return nil, args.Error(1)
	}

	return args.Get(0).(*appsec.GetSecurityPolicyCloneResponse), args.Error(1)
}
func (p *mockappsec) GetSecurityPolicyClones(ctx context.Context, params appsec.GetSecurityPolicyClonesRequest) (*appsec.GetSecurityPolicyClonesResponse, error) {
	args := p.Called(ctx, params)

	if args.Get(0) == nil {
		return nil, args.Error(1)
	}

	return args.Get(0).(*appsec.GetSecurityPolicyClonesResponse), args.Error(1)
}

func (p *mockappsec) GetSecurityPolicy(ctx context.Context, params appsec.GetSecurityPolicyRequest) (*appsec.GetSecurityPolicyResponse, error) {
	args := p.Called(ctx, params)

	if args.Get(0) == nil {
		return nil, args.Error(1)
	}

	return args.Get(0).(*appsec.GetSecurityPolicyResponse), args.Error(1)
}

func (p *mockappsec) CreateSecurityPolicy(ctx context.Context, params appsec.CreateSecurityPolicyRequest) (*appsec.CreateSecurityPolicyResponse, error) {
	args := p.Called(ctx, params)

	if args.Get(0) == nil {
		return nil, args.Error(1)
	}

	return args.Get(0).(*appsec.CreateSecurityPolicyResponse), args.Error(1)
}

func (p *mockappsec) UpdateSecurityPolicy(ctx context.Context, params appsec.UpdateSecurityPolicyRequest) (*appsec.UpdateSecurityPolicyResponse, error) {
	args := p.Called(ctx, params)

	if args.Get(0) == nil {
		return nil, args.Error(1)
	}

	return args.Get(0).(*appsec.UpdateSecurityPolicyResponse), args.Error(1)
}

func (p *mockappsec) RemoveSecurityPolicy(ctx context.Context, params appsec.RemoveSecurityPolicyRequest) (*appsec.RemoveSecurityPolicyResponse, error) {
	args := p.Called(ctx, params)

	if args.Get(0) == nil {
		return nil, args.Error(1)
	}

	return args.Get(0).(*appsec.RemoveSecurityPolicyResponse), args.Error(1)
}

func (p *mockappsec) GetCustomRule(ctx context.Context, params appsec.GetCustomRuleRequest) (*appsec.GetCustomRuleResponse, error) {
	args := p.Called(ctx, params)

	if args.Get(0) == nil {
		return nil, args.Error(1)
	}

	return args.Get(0).(*appsec.GetCustomRuleResponse), args.Error(1)
}

func (p *mockappsec) GetCustomRules(ctx context.Context, params appsec.GetCustomRulesRequest) (*appsec.GetCustomRulesResponse, error) {
	args := p.Called(ctx, params)

	if args.Get(0) == nil {
		return nil, args.Error(1)
	}

	return args.Get(0).(*appsec.GetCustomRulesResponse), args.Error(1)
}

func (p *mockappsec) GetCustomRuleAction(ctx context.Context, params appsec.GetCustomRuleActionRequest) (*appsec.GetCustomRuleActionResponse, error) {
	args := p.Called(ctx, params)

	if args.Get(0) == nil {
		return nil, args.Error(1)
	}

	return args.Get(0).(*appsec.GetCustomRuleActionResponse), args.Error(1)
}

func (p *mockappsec) UpdateCustomRuleAction(ctx context.Context, params appsec.UpdateCustomRuleActionRequest) (*appsec.UpdateCustomRuleActionResponse, error) {
	args := p.Called(ctx, params)

	if args.Get(0) == nil {
		return nil, args.Error(1)
	}

	return args.Get(0).(*appsec.UpdateCustomRuleActionResponse), args.Error(1)
}

func (p *mockappsec) GetCustomRuleActions(ctx context.Context, params appsec.GetCustomRuleActionsRequest) (*appsec.GetCustomRuleActionsResponse, error) {
	args := p.Called(ctx, params)

	if args.Get(0) == nil {
		return nil, args.Error(1)
	}

	return args.Get(0).(*appsec.GetCustomRuleActionsResponse), args.Error(1)
}

func (p *mockappsec) GetExportConfigurations(ctx context.Context, params appsec.GetExportConfigurationsRequest) (*appsec.GetExportConfigurationsResponse, error) {
	args := p.Called(ctx, params)

	if args.Get(0) == nil {
		return nil, args.Error(1)
	}

	return args.Get(0).(*appsec.GetExportConfigurationsResponse), args.Error(1)
}

func (p *mockappsec) GetMatchTarget(ctx context.Context, params appsec.GetMatchTargetRequest) (*appsec.GetMatchTargetResponse, error) {
	args := p.Called(ctx, params)

	if args.Get(0) == nil {
		return nil, args.Error(1)
	}

	return args.Get(0).(*appsec.GetMatchTargetResponse), args.Error(1)
}

func (p *mockappsec) UpdateMatchTarget(ctx context.Context, params appsec.UpdateMatchTargetRequest) (*appsec.UpdateMatchTargetResponse, error) {
	args := p.Called(ctx, params)

	if args.Get(0) == nil {
		return nil, args.Error(1)
	}

	return args.Get(0).(*appsec.UpdateMatchTargetResponse), args.Error(1)
}

func (p *mockappsec) GetMatchTargets(ctx context.Context, params appsec.GetMatchTargetsRequest) (*appsec.GetMatchTargetsResponse, error) {
	args := p.Called(ctx, params)

	if args.Get(0) == nil {
		return nil, args.Error(1)
	}

	return args.Get(0).(*appsec.GetMatchTargetsResponse), args.Error(1)
}

func (p *mockappsec) GetMatchTargetSequence(ctx context.Context, params appsec.GetMatchTargetSequenceRequest) (*appsec.GetMatchTargetSequenceResponse, error) {
	args := p.Called(ctx, params)

	if args.Get(0) == nil {
		return nil, args.Error(1)
	}

	return args.Get(0).(*appsec.GetMatchTargetSequenceResponse), args.Error(1)
}

func (p *mockappsec) UpdateMatchTargetSequence(ctx context.Context, params appsec.UpdateMatchTargetSequenceRequest) (*appsec.UpdateMatchTargetSequenceResponse, error) {
	args := p.Called(ctx, params)

	if args.Get(0) == nil {
		return nil, args.Error(1)
	}

	return args.Get(0).(*appsec.UpdateMatchTargetSequenceResponse), args.Error(1)
}

func (p *mockappsec) GetMatchTargetSequences(ctx context.Context, params appsec.GetMatchTargetSequencesRequest) (*appsec.GetMatchTargetSequencesResponse, error) {
	args := p.Called(ctx, params)

	if args.Get(0) == nil {
		return nil, args.Error(1)
	}

	return args.Get(0).(*appsec.GetMatchTargetSequencesResponse), args.Error(1)
}

func (p *mockappsec) GetPenaltyBox(ctx context.Context, params appsec.GetPenaltyBoxRequest) (*appsec.GetPenaltyBoxResponse, error) {
	args := p.Called(ctx, params)

	if args.Get(0) == nil {
		return nil, args.Error(1)
	}

	return args.Get(0).(*appsec.GetPenaltyBoxResponse), args.Error(1)
}

func (p *mockappsec) UpdatePenaltyBox(ctx context.Context, params appsec.UpdatePenaltyBoxRequest) (*appsec.UpdatePenaltyBoxResponse, error) {
	args := p.Called(ctx, params)

	if args.Get(0) == nil {
		return nil, args.Error(1)
	}

	return args.Get(0).(*appsec.UpdatePenaltyBoxResponse), args.Error(1)
}

func (p *mockappsec) GetPenaltyBoxes(ctx context.Context, params appsec.GetPenaltyBoxesRequest) (*appsec.GetPenaltyBoxesResponse, error) {
	args := p.Called(ctx, params)

	if args.Get(0) == nil {
		return nil, args.Error(1)
	}

	return args.Get(0).(*appsec.GetPenaltyBoxesResponse), args.Error(1)
}

func (p *mockappsec) GetSecurityPolicies(ctx context.Context, params appsec.GetSecurityPoliciesRequest) (*appsec.GetSecurityPoliciesResponse, error) {
	args := p.Called(ctx, params)

	if args.Get(0) == nil {
		return nil, args.Error(1)
	}

	return args.Get(0).(*appsec.GetSecurityPoliciesResponse), args.Error(1)
}

func (p *mockappsec) GetSelectableHostnames(ctx context.Context, params appsec.GetSelectableHostnamesRequest) (*appsec.GetSelectableHostnamesResponse, error) {
	args := p.Called(ctx, params)

	if args.Get(0) == nil {
		return nil, args.Error(1)
	}

	return args.Get(0).(*appsec.GetSelectableHostnamesResponse), args.Error(1)
}

func (p *mockappsec) GetSelectedHostname(ctx context.Context, params appsec.GetSelectedHostnameRequest) (*appsec.GetSelectedHostnameResponse, error) {
	args := p.Called(ctx, params)

	if args.Get(0) == nil {
		return nil, args.Error(1)
	}

	return args.Get(0).(*appsec.GetSelectedHostnameResponse), args.Error(1)
}

func (p *mockappsec) UpdateSelectedHostname(ctx context.Context, params appsec.UpdateSelectedHostnameRequest) (*appsec.UpdateSelectedHostnameResponse, error) {
	args := p.Called(ctx, params)

	if args.Get(0) == nil {
		return nil, args.Error(1)
	}

	return args.Get(0).(*appsec.UpdateSelectedHostnameResponse), args.Error(1)
}

func (p *mockappsec) GetSelectedHostnames(ctx context.Context, params appsec.GetSelectedHostnamesRequest) (*appsec.GetSelectedHostnamesResponse, error) {
	args := p.Called(ctx, params)

	if args.Get(0) == nil {
		return nil, args.Error(1)
	}

	return args.Get(0).(*appsec.GetSelectedHostnamesResponse), args.Error(1)
}

func (p *mockappsec) GetSlowPostProtectionSetting(ctx context.Context, params appsec.GetSlowPostProtectionSettingRequest) (*appsec.GetSlowPostProtectionSettingResponse, error) {
	args := p.Called(ctx, params)

	if args.Get(0) == nil {
		return nil, args.Error(1)
	}

	return args.Get(0).(*appsec.GetSlowPostProtectionSettingResponse), args.Error(1)
}

func (p *mockappsec) GetSlowPostProtectionSettings(ctx context.Context, params appsec.GetSlowPostProtectionSettingsRequest) (*appsec.GetSlowPostProtectionSettingsResponse, error) {
	args := p.Called(ctx, params)

	if args.Get(0) == nil {
		return nil, args.Error(1)
	}

	return args.Get(0).(*appsec.GetSlowPostProtectionSettingsResponse), args.Error(1)
}

func (p *mockappsec) UpdateSlowPostProtectionSetting(ctx context.Context, params appsec.UpdateSlowPostProtectionSettingRequest) (*appsec.UpdateSlowPostProtectionSettingResponse, error) {
	args := p.Called(ctx, params)

	if args.Get(0) == nil {
		return nil, args.Error(1)
	}

	return args.Get(0).(*appsec.UpdateSlowPostProtectionSettingResponse), args.Error(1)
}

func (p *mockappsec) GetWAFMode(ctx context.Context, params appsec.GetWAFModeRequest) (*appsec.GetWAFModeResponse, error) {
	args := p.Called(ctx, params)

	if args.Get(0) == nil {
		return nil, args.Error(1)
	}

	return args.Get(0).(*appsec.GetWAFModeResponse), args.Error(1)
}

func (p *mockappsec) GetWAFModes(ctx context.Context, params appsec.GetWAFModesRequest) (*appsec.GetWAFModesResponse, error) {
	args := p.Called(ctx, params)

	if args.Get(0) == nil {
		return nil, args.Error(1)
	}

	return args.Get(0).(*appsec.GetWAFModesResponse), args.Error(1)
}

func (p *mockappsec) UpdateWAFMode(ctx context.Context, params appsec.UpdateWAFModeRequest) (*appsec.UpdateWAFModeResponse, error) {
	args := p.Called(ctx, params)

	if args.Get(0) == nil {
		return nil, args.Error(1)
	}

	return args.Get(0).(*appsec.UpdateWAFModeResponse), args.Error(1)
}

func (p *mockappsec) GetEval(ctx context.Context, params appsec.GetEvalRequest) (*appsec.GetEvalResponse, error) {
	args := p.Called(ctx, params)

	if args.Get(0) == nil {
		return nil, args.Error(1)
	}

	return args.Get(0).(*appsec.GetEvalResponse), args.Error(1)
}

func (p *mockappsec) GetEvals(ctx context.Context, params appsec.GetEvalsRequest) (*appsec.GetEvalsResponse, error) {
	args := p.Called(ctx, params)

	if args.Get(0) == nil {
		return nil, args.Error(1)
	}

	return args.Get(0).(*appsec.GetEvalsResponse), args.Error(1)
}

func (p *mockappsec) UpdateEval(ctx context.Context, params appsec.UpdateEvalRequest) (*appsec.UpdateEvalResponse, error) {
	args := p.Called(ctx, params)

	if args.Get(0) == nil {
		return nil, args.Error(1)
	}

	return args.Get(0).(*appsec.UpdateEvalResponse), args.Error(1)
}

func (p *mockappsec) GetWAFProtection(ctx context.Context, params appsec.GetWAFProtectionRequest) (*appsec.GetWAFProtectionResponse, error) {
	args := p.Called(ctx, params)

	if args.Get(0) == nil {
		return nil, args.Error(1)
	}

	return args.Get(0).(*appsec.GetWAFProtectionResponse), args.Error(1)
}

func (p *mockappsec) GetWAFProtections(ctx context.Context, params appsec.GetWAFProtectionsRequest) (*appsec.GetWAFProtectionsResponse, error) {
	args := p.Called(ctx, params)

	if args.Get(0) == nil {
		return nil, args.Error(1)
	}

	return args.Get(0).(*appsec.GetWAFProtectionsResponse), args.Error(1)
}

func (p *mockappsec) UpdateWAFProtection(ctx context.Context, params appsec.UpdateWAFProtectionRequest) (*appsec.UpdateWAFProtectionResponse, error) {
	args := p.Called(ctx, params)

	if args.Get(0) == nil {
		return nil, args.Error(1)
	}

	return args.Get(0).(*appsec.UpdateWAFProtectionResponse), args.Error(1)
}

func (p *mockappsec) GetPolicyProtections(ctx context.Context, params appsec.GetPolicyProtectionsRequest) (*appsec.GetPolicyProtectionsResponse, error) {
	args := p.Called(ctx, params)

	if args.Get(0) == nil {
		return nil, args.Error(1)
	}

	return args.Get(0).(*appsec.GetPolicyProtectionsResponse), args.Error(1)
}

func (p *mockappsec) UpdatePolicyProtections(ctx context.Context, params appsec.UpdatePolicyProtectionsRequest) (*appsec.UpdatePolicyProtectionsResponse, error) {
	args := p.Called(ctx, params)

	if args.Get(0) == nil {
		return nil, args.Error(1)
	}

	return args.Get(0).(*appsec.UpdatePolicyProtectionsResponse), args.Error(1)
}

func (p *mockappsec) GetRateProtection(ctx context.Context, params appsec.GetRateProtectionRequest) (*appsec.GetRateProtectionResponse, error) {
	args := p.Called(ctx, params)

	if args.Get(0) == nil {
		return nil, args.Error(1)
	}

	return args.Get(0).(*appsec.GetRateProtectionResponse), args.Error(1)
}

func (p *mockappsec) GetRateProtections(ctx context.Context, params appsec.GetRateProtectionsRequest) (*appsec.GetRateProtectionsResponse, error) {
	args := p.Called(ctx, params)

	if args.Get(0) == nil {
		return nil, args.Error(1)
	}

	return args.Get(0).(*appsec.GetRateProtectionsResponse), args.Error(1)
}

func (p *mockappsec) UpdateRateProtection(ctx context.Context, params appsec.UpdateRateProtectionRequest) (*appsec.UpdateRateProtectionResponse, error) {
	args := p.Called(ctx, params)

	if args.Get(0) == nil {
		return nil, args.Error(1)
	}

	return args.Get(0).(*appsec.UpdateRateProtectionResponse), args.Error(1)
}

func (p *mockappsec) GetRuleActions(ctx context.Context, params appsec.GetRuleActionsRequest) (*appsec.GetRuleActionsResponse, error) {
	args := p.Called(ctx, params)

	if args.Get(0) == nil {
		return nil, args.Error(1)
	}

	return args.Get(0).(*appsec.GetRuleActionsResponse), args.Error(1)
}

func (p *mockappsec) GetRuleAction(ctx context.Context, params appsec.GetRuleActionRequest) (*appsec.GetRuleActionResponse, error) {
	args := p.Called(ctx, params)

	if args.Get(0) == nil {
		return nil, args.Error(1)
	}

	return args.Get(0).(*appsec.GetRuleActionResponse), args.Error(1)
}

func (p *mockappsec) UpdateRuleAction(ctx context.Context, params appsec.UpdateRuleActionRequest) (*appsec.UpdateRuleActionResponse, error) {
	args := p.Called(ctx, params)

	if args.Get(0) == nil {
		return nil, args.Error(1)
	}

	return args.Get(0).(*appsec.UpdateRuleActionResponse), args.Error(1)
}

func (p *mockappsec) GetRuleConditionException(ctx context.Context, params appsec.GetRuleConditionExceptionRequest) (*appsec.GetRuleConditionExceptionResponse, error) {
	args := p.Called(ctx, params)

	if args.Get(0) == nil {
		return nil, args.Error(1)
	}

	return args.Get(0).(*appsec.GetRuleConditionExceptionResponse), args.Error(1)
}

func (p *mockappsec) GetRuleConditionExceptions(ctx context.Context, params appsec.GetRuleConditionExceptionsRequest) (*appsec.GetRuleConditionExceptionsResponse, error) {
	args := p.Called(ctx, params)

	if args.Get(0) == nil {
		return nil, args.Error(1)
	}

	return args.Get(0).(*appsec.GetRuleConditionExceptionsResponse), args.Error(1)
}

func (p *mockappsec) UpdateRuleConditionException(ctx context.Context, params appsec.UpdateRuleConditionExceptionRequest) (*appsec.UpdateRuleConditionExceptionResponse, error) {
	args := p.Called(ctx, params)

	if args.Get(0) == nil {
		return nil, args.Error(1)
	}

	return args.Get(0).(*appsec.UpdateRuleConditionExceptionResponse), args.Error(1)
}

func (p *mockappsec) RemoveRuleConditionException(ctx context.Context, params appsec.RemoveRuleConditionExceptionRequest) (*appsec.RemoveRuleConditionExceptionResponse, error) {
	args := p.Called(ctx, params)

	if args.Get(0) == nil {
		return nil, args.Error(1)
	}

	return args.Get(0).(*appsec.RemoveRuleConditionExceptionResponse), args.Error(1)
}

func (p *mockappsec) CreateAttackGroupAction(ctx context.Context, params appsec.CreateAttackGroupActionRequest) (*appsec.CreateAttackGroupActionResponse, error) {
	args := p.Called(ctx, params)

	if args.Get(0) == nil {
		return nil, args.Error(1)
	}

	return args.Get(0).(*appsec.CreateAttackGroupActionResponse), args.Error(1)
}

func (p *mockappsec) GetAttackGroupConditionException(ctx context.Context, params appsec.GetAttackGroupConditionExceptionRequest) (*appsec.GetAttackGroupConditionExceptionResponse, error) {
	args := p.Called(ctx, params)

	if args.Get(0) == nil {
		return nil, args.Error(1)
	}

	return args.Get(0).(*appsec.GetAttackGroupConditionExceptionResponse), args.Error(1)
}

func (p *mockappsec) GetAttackGroupConditionExceptions(ctx context.Context, params appsec.GetAttackGroupConditionExceptionsRequest) (*appsec.GetAttackGroupConditionExceptionsResponse, error) {
	args := p.Called(ctx, params)

	if args.Get(0) == nil {
		return nil, args.Error(1)
	}

	return args.Get(0).(*appsec.GetAttackGroupConditionExceptionsResponse), args.Error(1)
}

func (p *mockappsec) UpdateAttackGroupConditionException(ctx context.Context, params appsec.UpdateAttackGroupConditionExceptionRequest) (*appsec.UpdateAttackGroupConditionExceptionResponse, error) {
	args := p.Called(ctx, params)

	if args.Get(0) == nil {
		return nil, args.Error(1)
	}

	return args.Get(0).(*appsec.UpdateAttackGroupConditionExceptionResponse), args.Error(1)
}

func (p *mockappsec) RemoveAttackGroupConditionException(ctx context.Context, params appsec.RemoveAttackGroupConditionExceptionRequest) (*appsec.RemoveAttackGroupConditionExceptionResponse, error) {
	args := p.Called(ctx, params)

	if args.Get(0) == nil {
		return nil, args.Error(1)
	}

	return args.Get(0).(*appsec.RemoveAttackGroupConditionExceptionResponse), args.Error(1)
}

func (p *mockappsec) GetAttackGroupAction(ctx context.Context, params appsec.GetAttackGroupActionRequest) (*appsec.GetAttackGroupActionResponse, error) {
	args := p.Called(ctx, params)

	if args.Get(0) == nil {
		return nil, args.Error(1)
	}

	return args.Get(0).(*appsec.GetAttackGroupActionResponse), args.Error(1)
}

func (p *mockappsec) UpdateAttackGroupAction(ctx context.Context, params appsec.UpdateAttackGroupActionRequest) (*appsec.UpdateAttackGroupActionResponse, error) {
	args := p.Called(ctx, params)

	if args.Get(0) == nil {
		return nil, args.Error(1)
	}

	return args.Get(0).(*appsec.UpdateAttackGroupActionResponse), args.Error(1)
}

func (p *mockappsec) RemoveAttackGroupAction(ctx context.Context, params appsec.RemoveAttackGroupActionRequest) (*appsec.RemoveAttackGroupActionResponse, error) {
	args := p.Called(ctx, params)

	if args.Get(0) == nil {
		return nil, args.Error(1)
	}

	return args.Get(0).(*appsec.RemoveAttackGroupActionResponse), args.Error(1)
}

func (p *mockappsec) GetAttackGroupActions(ctx context.Context, params appsec.GetAttackGroupActionsRequest) (*appsec.GetAttackGroupActionsResponse, error) {
	args := p.Called(ctx, params)

	if args.Get(0) == nil {
		return nil, args.Error(1)
	}

	return args.Get(0).(*appsec.GetAttackGroupActionsResponse), args.Error(1)
}

func (p *mockappsec) GetReputationProtections(ctx context.Context, params appsec.GetReputationProtectionsRequest) (*appsec.GetReputationProtectionsResponse, error) {
	args := p.Called(ctx, params)

	if args.Get(0) == nil {
		return nil, args.Error(1)
	}

	return args.Get(0).(*appsec.GetReputationProtectionsResponse), args.Error(1)
}

func (p *mockappsec) GetReputationProtection(ctx context.Context, params appsec.GetReputationProtectionRequest) (*appsec.GetReputationProtectionResponse, error) {
	args := p.Called(ctx, params)

	if args.Get(0) == nil {
		return nil, args.Error(1)
	}

	return args.Get(0).(*appsec.GetReputationProtectionResponse), args.Error(1)
}

func (p *mockappsec) UpdateReputationProtection(ctx context.Context, params appsec.UpdateReputationProtectionRequest) (*appsec.UpdateReputationProtectionResponse, error) {
	args := p.Called(ctx, params)

	if args.Get(0) == nil {
		return nil, args.Error(1)
	}

	return args.Get(0).(*appsec.UpdateReputationProtectionResponse), args.Error(1)

}

func (p *mockappsec) GetSlowPostProtection(ctx context.Context, params appsec.GetSlowPostProtectionRequest) (*appsec.GetSlowPostProtectionResponse, error) {
	args := p.Called(ctx, params)

	if args.Get(0) == nil {
		return nil, args.Error(1)
	}

	return args.Get(0).(*appsec.GetSlowPostProtectionResponse), args.Error(1)
}

func (p *mockappsec) GetSlowPostProtections(ctx context.Context, params appsec.GetSlowPostProtectionsRequest) (*appsec.GetSlowPostProtectionsResponse, error) {
	args := p.Called(ctx, params)

	if args.Get(0) == nil {
		return nil, args.Error(1)
	}

	return args.Get(0).(*appsec.GetSlowPostProtectionsResponse), args.Error(1)
}

func (p *mockappsec) UpdateSlowPostProtection(ctx context.Context, params appsec.UpdateSlowPostProtectionRequest) (*appsec.UpdateSlowPostProtectionResponse, error) {
	args := p.Called(ctx, params)

	if args.Get(0) == nil {
		return nil, args.Error(1)
	}

	return args.Get(0).(*appsec.UpdateSlowPostProtectionResponse), args.Error(1)
}
func (p *mockappsec) GetEvalRuleAction(ctx context.Context, params appsec.GetEvalRuleActionRequest) (*appsec.GetEvalRuleActionResponse, error) {
	args := p.Called(ctx, params)

	if args.Get(0) == nil {
		return nil, args.Error(1)
	}

	return args.Get(0).(*appsec.GetEvalRuleActionResponse), args.Error(1)
}

func (p *mockappsec) GetEvalRuleActions(ctx context.Context, params appsec.GetEvalRuleActionsRequest) (*appsec.GetEvalRuleActionsResponse, error) {
	args := p.Called(ctx, params)

	if args.Get(0) == nil {
		return nil, args.Error(1)
	}

	return args.Get(0).(*appsec.GetEvalRuleActionsResponse), args.Error(1)
}

func (p *mockappsec) UpdateEvalRuleAction(ctx context.Context, params appsec.UpdateEvalRuleActionRequest) (*appsec.UpdateEvalRuleActionResponse, error) {
	args := p.Called(ctx, params)

	if args.Get(0) == nil {
		return nil, args.Error(1)
	}

	return args.Get(0).(*appsec.UpdateEvalRuleActionResponse), args.Error(1)
}

func (p *mockappsec) GetEvalRuleConditionException(ctx context.Context, params appsec.GetEvalRuleConditionExceptionRequest) (*appsec.GetEvalRuleConditionExceptionResponse, error) {
	args := p.Called(ctx, params)

	if args.Get(0) == nil {
		return nil, args.Error(1)
	}

	return args.Get(0).(*appsec.GetEvalRuleConditionExceptionResponse), args.Error(1)
}

func (p *mockappsec) GetEvalRuleConditionExceptions(ctx context.Context, params appsec.GetEvalRuleConditionExceptionsRequest) (*appsec.GetEvalRuleConditionExceptionsResponse, error) {
	args := p.Called(ctx, params)

	if args.Get(0) == nil {
		return nil, args.Error(1)
	}

	return args.Get(0).(*appsec.GetEvalRuleConditionExceptionsResponse), args.Error(1)
}

func (p *mockappsec) UpdateEvalRuleConditionException(ctx context.Context, params appsec.UpdateEvalRuleConditionExceptionRequest) (*appsec.UpdateEvalRuleConditionExceptionResponse, error) {
	args := p.Called(ctx, params)

	if args.Get(0) == nil {
		return nil, args.Error(1)
	}

	return args.Get(0).(*appsec.UpdateEvalRuleConditionExceptionResponse), args.Error(1)
}

func (p *mockappsec) RemoveEvalRuleConditionException(ctx context.Context, params appsec.RemoveEvalRuleConditionExceptionRequest) (*appsec.RemoveEvalRuleConditionExceptionResponse, error) {
	args := p.Called(ctx, params)

	if args.Get(0) == nil {
		return nil, args.Error(1)
	}

	return args.Get(0).(*appsec.RemoveEvalRuleConditionExceptionResponse), args.Error(1)
}
