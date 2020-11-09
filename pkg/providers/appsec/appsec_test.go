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

func (p *mockappsec) GetKRSRuleActions(ctx context.Context, params appsec.GetKRSRuleActionsRequest) (*appsec.GetKRSRuleActionsResponse, error) {
	args := p.Called(ctx, params)

	if args.Get(0) == nil {
		return nil, args.Error(1)
	}

	return args.Get(0).(*appsec.GetKRSRuleActionsResponse), args.Error(1)
}

func (p *mockappsec) GetKRSRuleAction(ctx context.Context, params appsec.GetKRSRuleActionRequest) (*appsec.GetKRSRuleActionResponse, error) {
	args := p.Called(ctx, params)

	if args.Get(0) == nil {
		return nil, args.Error(1)
	}

	return args.Get(0).(*appsec.GetKRSRuleActionResponse), args.Error(1)
}

func (p *mockappsec) UpdateKRSRuleAction(ctx context.Context, params appsec.UpdateKRSRuleActionRequest) (*appsec.UpdateKRSRuleActionResponse, error) {
	args := p.Called(ctx, params)

	if args.Get(0) == nil {
		return nil, args.Error(1)
	}

	return args.Get(0).(*appsec.UpdateKRSRuleActionResponse), args.Error(1)
}

func (p *mockappsec) CreateWAFAttackGroupAction(ctx context.Context, params appsec.CreateWAFAttackGroupActionRequest) (*appsec.CreateWAFAttackGroupActionResponse, error) {
	args := p.Called(ctx, params)

	if args.Get(0) == nil {
		return nil, args.Error(1)
	}

	return args.Get(0).(*appsec.CreateWAFAttackGroupActionResponse), args.Error(1)
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

func (p *mockappsec) GetWAFAttackGroupAction(ctx context.Context, params appsec.GetWAFAttackGroupActionRequest) (*appsec.GetWAFAttackGroupActionResponse, error) {
	args := p.Called(ctx, params)

	if args.Get(0) == nil {
		return nil, args.Error(1)
	}

	return args.Get(0).(*appsec.GetWAFAttackGroupActionResponse), args.Error(1)
}

func (p *mockappsec) UpdateWAFAttackGroupAction(ctx context.Context, params appsec.UpdateWAFAttackGroupActionRequest) (*appsec.UpdateWAFAttackGroupActionResponse, error) {
	args := p.Called(ctx, params)

	if args.Get(0) == nil {
		return nil, args.Error(1)
	}

	return args.Get(0).(*appsec.UpdateWAFAttackGroupActionResponse), args.Error(1)
}

func (p *mockappsec) RemoveWAFAttackGroupAction(ctx context.Context, params appsec.RemoveWAFAttackGroupActionRequest) (*appsec.RemoveWAFAttackGroupActionResponse, error) {
	args := p.Called(ctx, params)

	if args.Get(0) == nil {
		return nil, args.Error(1)
	}

	return args.Get(0).(*appsec.RemoveWAFAttackGroupActionResponse), args.Error(1)
}

func (p *mockappsec) GetWAFAttackGroupActions(ctx context.Context, params appsec.GetWAFAttackGroupActionsRequest) (*appsec.GetWAFAttackGroupActionsResponse, error) {
	args := p.Called(ctx, params)

	if args.Get(0) == nil {
		return nil, args.Error(1)
	}

	return args.Get(0).(*appsec.GetWAFAttackGroupActionsResponse), args.Error(1)
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
