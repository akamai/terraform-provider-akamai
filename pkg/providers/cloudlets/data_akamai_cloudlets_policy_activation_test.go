package cloudlets

import (
	"fmt"
	"net/http"
	"regexp"
	"testing"
	"time"

	"github.com/akamai/AkamaiOPEN-edgegrid-golang/v8/pkg/cloudlets"
	v3 "github.com/akamai/AkamaiOPEN-edgegrid-golang/v8/pkg/cloudlets/v3"
	"github.com/akamai/terraform-provider-akamai/v6/pkg/common/ptr"
	"github.com/akamai/terraform-provider-akamai/v6/pkg/common/testutils"
	"github.com/hashicorp/terraform-plugin-testing/helper/resource"
	"github.com/stretchr/testify/mock"
)

type testDataForNonSharedPolicyActivation struct {
	id          string
	policyID    int64
	version     int64
	groupID     int64
	name        string
	description string
	matchRules  cloudlets.MatchRules
	properties  []string
	network     cloudlets.PolicyActivationNetwork
}

type testDataForSharedPolicyActivation struct {
	id           string
	policyID     int64
	version      int64
	groupID      int64
	name         string
	cloudletType v3.CloudletType
	description  string
	matchRules   v3.MatchRules
	warnings     []v3.MatchRulesWarning
	activations  v3.CurrentActivations
}

func TestNonSharedPolicyActivationDataSource(t *testing.T) {
	t.Parallel()

	tests := map[string]struct {
		config      string
		data        testDataForNonSharedPolicyActivation
		expectError *regexp.Regexp
		init        func(*cloudlets.Mock, testDataForNonSharedPolicyActivation)
		check       resource.TestCheckFunc
	}{
		"no property id": {
			config:      "no_property_id.tf",
			expectError: regexp.MustCompile(`Error: Missing required argument`),
		},
		"no network": {
			config:      "no_network.tf",
			expectError: regexp.MustCompile(`Error: Missing required argument`),
		},
		"policy without activation": {
			config: "activation.tf",
			data: testDataForNonSharedPolicyActivation{
				id:          "akamai_cloudlets_shared_policy",
				policyID:    1,
				version:     2,
				groupID:     12,
				name:        "TestName",
				description: "Description",
				network:     cloudlets.PolicyActivationNetworkStaging,
				properties:  []string{"prp_0", "prp_1"},
			},
			init: func(m2 *cloudlets.Mock, data testDataForNonSharedPolicyActivation) {
				mockGetPolicyV2(m2, data, nil, 1)
				expectListPolicyActivations(m2, data.policyID, data.version, data.network, data.properties, cloudlets.PolicyActivationStatusActive, "", 0, nil).Times(1)
			},
			expectError: regexp.MustCompile(`(?s)policy activation read: cannot find any activation for the given policy '1'.+and network 'staging'`),
		},
		"policy with activation": {
			config: "activation.tf",
			data: testDataForNonSharedPolicyActivation{
				id:          "akamai_cloudlets_shared_policy",
				policyID:    1,
				version:     2,
				groupID:     12,
				name:        "TestName",
				description: "Description",
				network:     cloudlets.PolicyActivationNetworkStaging,
				properties:  []string{"prp_0", "prp_1"},
			},
			init: func(m2 *cloudlets.Mock, data testDataForNonSharedPolicyActivation) {
				activations := make([]cloudlets.PolicyActivation, len(data.properties))
				for _, p := range data.properties {
					activations = append(activations, cloudlets.PolicyActivation{APIVersion: "1.0", Network: data.network, PolicyInfo: cloudlets.PolicyInfo{
						PolicyID: data.policyID, Version: data.version, Status: cloudlets.PolicyActivationStatusInactive,
					}, PropertyInfo: cloudlets.PropertyInfo{Name: p}})
				}
				mockGetPolicyV2(m2, data, nil, 3)
				expectListPolicyActivations(m2, data.policyID, data.version, data.network, data.properties, cloudlets.PolicyActivationStatusActive, "", 1, nil).Times(3)
			},
			check: resource.ComposeAggregateTestCheckFunc(
				resource.TestCheckResourceAttr("data.akamai_cloudlets_policy_activation.test", "policy_id", "1"),
				resource.TestCheckResourceAttr("data.akamai_cloudlets_policy_activation.test", "network", "staging"),
				resource.TestCheckResourceAttr("data.akamai_cloudlets_policy_activation.test", "version", "2"),
				resource.TestCheckResourceAttr("data.akamai_cloudlets_policy_activation.test", "status", "active"),
			),
		},
		"api error": {
			config: "activation.tf",
			data: testDataForNonSharedPolicyActivation{
				policyID:    1,
				version:     2,
				groupID:     12,
				name:        "Name",
				description: "Description",
				matchRules: cloudlets.MatchRules{
					cloudlets.MatchRuleER{
						Name:           "Name",
						Type:           cloudlets.MatchRuleTypeER,
						Start:          1,
						End:            2,
						ID:             123,
						UseRelativeURL: "/url1",
						StatusCode:     200,
						RedirectURL:    "/url2",
						MatchURL:       "/url3",
					},
				},
			},
			expectError: regexp.MustCompile(`Error: Reading Policy Failed`),
			init: func(m2 *cloudlets.Mock, data testDataForNonSharedPolicyActivation) {
				mockGetPolicyV2WithError(m2, data.policyID, &cloudlets.Error{StatusCode: http.StatusNotFound}, 1)
			},
		},
	}

	for name, test := range tests {
		t.Run(name, func(t *testing.T) {
			clientV2 := &cloudlets.Mock{}
			if test.init != nil {
				test.init(clientV2, test.data)
			}
			useClient(clientV2, func() {
				resource.Test(t, resource.TestCase{
					IsUnitTest:               true,
					ProtoV6ProviderFactories: testutils.NewProtoV6ProviderFactory(NewSubprovider()),
					Steps: []resource.TestStep{{
						Config:      testutils.LoadFixtureString(t, fmt.Sprintf("testdata/TestDataCloudletsPolicyActivation/%s", test.config)),
						Check:       test.check,
						ExpectError: test.expectError,
					}},
				})
			})
			clientV2.AssertExpectations(t)
		})
	}
}

func TestSharedPolicyActivationDataSource(t *testing.T) {
	t.Parallel()

	tests := map[string]struct {
		config      string
		data        testDataForSharedPolicyActivation
		expectError *regexp.Regexp
		init        func(*cloudlets.Mock, *v3.Mock, testDataForSharedPolicyActivation)
		check       resource.TestCheckFunc
	}{
		"no property id": {
			config:      "no_property_id.tf",
			expectError: regexp.MustCompile(`Error: Missing required argument`),
		},
		"no network": {
			config:      "no_network.tf",
			expectError: regexp.MustCompile(`Error: Missing required argument`),
		},
		"no shared policy activation": {
			config: "activation.tf",
			data: testDataForSharedPolicyActivation{
				id:           "akamai_cloudlets_shared_policy",
				policyID:     1,
				version:      0,
				groupID:      12,
				name:         "TestName",
				cloudletType: v3.CloudletTypeAP,
				description:  "Description",
			},
			init: func(m2 *cloudlets.Mock, m3 *v3.Mock, data testDataForSharedPolicyActivation) {
				mockGetPolicyV2WithError(m2, data.policyID, &cloudlets.Error{StatusCode: http.StatusNotFound}, 1)
				mockGetPolicyV3(m3, data, 2)
			},
			expectError: regexp.MustCompile(`(?s)policy activation read: cannot find any activation for the given policy '1'.+and network 'staging'`),
		},
		"no shared policy": {
			config: "activation.tf",
			data: testDataForSharedPolicyActivation{
				id:           "akamai_cloudlets_shared_policy",
				policyID:     1,
				version:      0,
				groupID:      12,
				name:         "TestName",
				cloudletType: v3.CloudletTypeAP,
				description:  "Description",
			},
			init: func(m2 *cloudlets.Mock, m3 *v3.Mock, data testDataForSharedPolicyActivation) {
				mockGetPolicyV2WithError(m2, data.policyID, &cloudlets.Error{StatusCode: http.StatusNotFound}, 1)
				mockGetPolicyV3WithError(m3, data.policyID, &cloudlets.Error{StatusCode: http.StatusNotFound}, 1)
			},
			expectError: regexp.MustCompile(`Error: Reading Policy Failed`),
		},
		"shared policy activation on staging": {
			config: "activation.tf",
			data: testDataForSharedPolicyActivation{
				//id:           "akamai_cloudlets_shared_policy",
				policyID:     1,
				version:      2,
				groupID:      12,
				name:         "Name",
				cloudletType: v3.CloudletTypeAP,
				description:  "Description",
				matchRules: v3.MatchRules{
					v3.MatchRuleER{
						Name:           "Name",
						Type:           v3.MatchRuleTypeER,
						Start:          1,
						End:            2,
						ID:             123,
						UseRelativeURL: "/url1",
						StatusCode:     200,
						RedirectURL:    "/url2",
						MatchURL:       "/url3",
					},
				},
				warnings: []v3.MatchRulesWarning{
					{
						Detail:      "TestDetail",
						JSONPointer: "TestPointer",
						Title:       "TestTitle",
						Type:        "TestType",
					},
				},
				activations: v3.CurrentActivations{
					Production: v3.ActivationInfo{
						Effective: &v3.PolicyActivation{
							CreatedBy:            "TestUser",
							CreatedDate:          time.Date(2009, 11, 17, 20, 34, 58, 651387237, time.UTC),
							FinishDate:           &time.Time{},
							ID:                   111,
							Network:              v3.ProductionNetwork,
							Operation:            v3.OperationActivation,
							PolicyID:             1,
							Status:               v3.ActivationStatusSuccess,
							PolicyVersion:        2,
							PolicyVersionDeleted: false,
						},
						Latest: &v3.PolicyActivation{
							CreatedBy:            "TestUser",
							CreatedDate:          time.Date(2009, 11, 17, 20, 34, 58, 651387237, time.UTC),
							FinishDate:           &time.Time{},
							ID:                   112,
							Network:              v3.ProductionNetwork,
							Operation:            v3.OperationActivation,
							PolicyID:             1,
							Status:               v3.ActivationStatusSuccess,
							PolicyVersion:        2,
							PolicyVersionDeleted: false,
						},
					},
					Staging: v3.ActivationInfo{
						Effective: &v3.PolicyActivation{
							CreatedBy:            "TestUser",
							CreatedDate:          time.Date(2009, 11, 17, 20, 34, 58, 651387237, time.UTC),
							FinishDate:           &time.Time{},
							ID:                   113,
							Network:              v3.StagingNetwork,
							Operation:            v3.OperationActivation,
							PolicyID:             1,
							Status:               v3.ActivationStatusSuccess,
							PolicyVersion:        2,
							PolicyVersionDeleted: false,
						},
						Latest: &v3.PolicyActivation{
							CreatedBy:            "TestUser",
							CreatedDate:          time.Date(2009, 11, 17, 20, 34, 58, 651387237, time.UTC),
							FinishDate:           &time.Time{},
							ID:                   114,
							Network:              v3.StagingNetwork,
							Operation:            v3.OperationActivation,
							PolicyID:             1,
							Status:               v3.ActivationStatusSuccess,
							PolicyVersion:        2,
							PolicyVersionDeleted: false,
						},
					},
				},
			},
			check: resource.ComposeAggregateTestCheckFunc(
				resource.TestCheckResourceAttr("data.akamai_cloudlets_policy_activation.test", "policy_id", "1"),
				resource.TestCheckResourceAttr("data.akamai_cloudlets_policy_activation.test", "network", "staging"),
				resource.TestCheckResourceAttr("data.akamai_cloudlets_policy_activation.test", "version", "2"),
				resource.TestCheckResourceAttr("data.akamai_cloudlets_policy_activation.test", "status", "SUCCESS"),
			),
			init: func(m2 *cloudlets.Mock, m3 *v3.Mock, data testDataForSharedPolicyActivation) {
				mockGetPolicyV2WithError(m2, data.policyID, &cloudlets.Error{StatusCode: http.StatusNotFound}, 3)
				mockGetPolicyV3(m3, data, 6)
			},
		},
		"api error": {
			config: "activation.tf",
			data: testDataForSharedPolicyActivation{
				//id:           "akamai_cloudlets_shared_policy",
				policyID:     1,
				version:      2,
				groupID:      12,
				name:         "Name",
				cloudletType: v3.CloudletTypeAP,
				description:  "Description",
				matchRules: v3.MatchRules{
					v3.MatchRuleER{
						Name:           "Name",
						Type:           v3.MatchRuleTypeER,
						Start:          1,
						End:            2,
						ID:             123,
						UseRelativeURL: "/url1",
						StatusCode:     200,
						RedirectURL:    "/url2",
						MatchURL:       "/url3",
					},
				},
				warnings: []v3.MatchRulesWarning{
					{
						Detail:      "TestDetail",
						JSONPointer: "TestPointer",
						Title:       "TestTitle",
						Type:        "TestType",
					},
				},
				activations: v3.CurrentActivations{
					Production: v3.ActivationInfo{
						Effective: &v3.PolicyActivation{
							CreatedBy:            "TestUser",
							CreatedDate:          time.Date(2009, 11, 17, 20, 34, 58, 651387237, time.UTC),
							FinishDate:           &time.Time{},
							ID:                   111,
							Network:              v3.ProductionNetwork,
							Operation:            v3.OperationActivation,
							PolicyID:             1,
							Status:               v3.ActivationStatusSuccess,
							PolicyVersion:        2,
							PolicyVersionDeleted: false,
						},
						Latest: &v3.PolicyActivation{
							CreatedBy:            "TestUser",
							CreatedDate:          time.Date(2009, 11, 17, 20, 34, 58, 651387237, time.UTC),
							FinishDate:           &time.Time{},
							ID:                   112,
							Network:              v3.ProductionNetwork,
							Operation:            v3.OperationActivation,
							PolicyID:             1,
							Status:               v3.ActivationStatusSuccess,
							PolicyVersion:        2,
							PolicyVersionDeleted: false,
						},
					},
					Staging: v3.ActivationInfo{
						Effective: &v3.PolicyActivation{
							CreatedBy:            "TestUser",
							CreatedDate:          time.Date(2009, 11, 17, 20, 34, 58, 651387237, time.UTC),
							FinishDate:           &time.Time{},
							ID:                   113,
							Network:              v3.StagingNetwork,
							Operation:            v3.OperationActivation,
							PolicyID:             1,
							Status:               v3.ActivationStatusSuccess,
							PolicyVersion:        2,
							PolicyVersionDeleted: false,
						},
						Latest: &v3.PolicyActivation{
							CreatedBy:            "TestUser",
							CreatedDate:          time.Date(2009, 11, 17, 20, 34, 58, 651387237, time.UTC),
							FinishDate:           &time.Time{},
							ID:                   114,
							Network:              v3.StagingNetwork,
							Operation:            v3.OperationActivation,
							PolicyID:             1,
							Status:               v3.ActivationStatusSuccess,
							PolicyVersion:        2,
							PolicyVersionDeleted: false,
						},
					},
				},
			},
			expectError: regexp.MustCompile(`policy activation read: reading policy failed`),
			init: func(m2 *cloudlets.Mock, m3 *v3.Mock, data testDataForSharedPolicyActivation) {
				mockGetPolicyV2WithError(m2, data.policyID, &cloudlets.Error{StatusCode: http.StatusNotFound}, 1)
				mockGetPolicyV3(m3, data, 1)
				mockGetPolicyV3WithError(m3, data.policyID, &cloudlets.Error{StatusCode: http.StatusInternalServerError}, 1)
			},
		},
	}

	for name, test := range tests {
		t.Run(name, func(t *testing.T) {
			clientV2 := &cloudlets.Mock{}
			clientV3 := &v3.Mock{}
			if test.init != nil {
				test.init(clientV2, clientV3, test.data)
			}
			useClientV2AndV3(clientV2, clientV3, func() {
				resource.Test(t, resource.TestCase{
					IsUnitTest:               true,
					ProtoV6ProviderFactories: testutils.NewProtoV6ProviderFactory(NewSubprovider()),
					Steps: []resource.TestStep{{
						Config:      testutils.LoadFixtureString(t, fmt.Sprintf("testdata/TestDataCloudletsPolicyActivation/%s", test.config)),
						Check:       test.check,
						ExpectError: test.expectError,
					}},
				})
			})
			clientV2.AssertExpectations(t)
			clientV3.AssertExpectations(t)
		})
	}
}

func mockGetPolicyV2WithError(m *cloudlets.Mock, policyID int64, err error, times int) {
	m.On("GetPolicy", mock.Anything, cloudlets.GetPolicyRequest{
		PolicyID: policyID,
	}).Return(nil, err).Times(times)
	return
}

func mockGetPolicyV2(m *cloudlets.Mock, data testDataForNonSharedPolicyActivation, err error, times int) {
	if err != nil {
		m.On("GetPolicy", mock.Anything, cloudlets.GetPolicyRequest{
			PolicyID: data.policyID,
		}).Return(nil, err).Times(times)
		return
	}
	m.On("GetPolicy", mock.Anything, cloudlets.GetPolicyRequest{
		PolicyID: data.policyID,
	}).Return(&cloudlets.Policy{
		Description: data.description,
		GroupID:     data.groupID,
		PolicyID:    data.policyID,
		Name:        data.name,
	}, nil).Times(times)
}

func mockGetPolicyV3(m *v3.Mock, data testDataForSharedPolicyActivation, times int) {
	m.On("GetPolicy", mock.Anything, v3.GetPolicyRequest{
		PolicyID: data.policyID,
	}).Return(&v3.Policy{
		CloudletType:       data.cloudletType,
		CurrentActivations: data.activations,
		Description:        ptr.To(data.description),
		GroupID:            data.groupID,
		ID:                 data.policyID,
		Name:               data.name,
	}, nil).Times(times)
}

func mockGetPolicyV3WithError(m *v3.Mock, policyID int64, err error, times int) {
	m.On("GetPolicy", mock.Anything, v3.GetPolicyRequest{
		PolicyID: policyID,
	}).Return(nil, err).Times(times)
}
