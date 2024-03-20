package cloudlets

import (
	"encoding/json"
	"fmt"
	"regexp"
	"strconv"
	"testing"
	"time"

	v3 "github.com/akamai/AkamaiOPEN-edgegrid-golang/v8/pkg/cloudlets/v3"
	"github.com/akamai/terraform-provider-akamai/v6/pkg/common/ptr"
	"github.com/akamai/terraform-provider-akamai/v6/pkg/common/testutils"
	"github.com/hashicorp/terraform-plugin-testing/helper/resource"
	"github.com/stretchr/testify/mock"
)

type testDataForSharedPolicy struct {
	id                 string
	policyID           int64
	version            int64
	versionDescription *string
	groupID            int64
	name               string
	cloudletType       v3.CloudletType
	description        string
	matchRules         v3.MatchRules
	warnings           []v3.MatchRulesWarning
	activations        v3.CurrentActivations
}

func TestSharedPolicyDataSource(t *testing.T) {
	t.Parallel()

	tests := map[string]struct {
		config      string
		data        testDataForSharedPolicy
		expectError *regexp.Regexp
		init        func(*v3.Mock, testDataForSharedPolicy)
	}{
		"success with version attribute - no activations": {
			config: "with_version.tf",
			data: testDataForSharedPolicy{
				id:                 "akamai_cloudlets_shared_policy",
				policyID:           1,
				version:            2,
				versionDescription: ptr.To("version 2 description"),
				groupID:            12,
				name:               "TestName",
				cloudletType:       v3.CloudletTypeAP,
				description:        "TestDescription",
				matchRules: v3.MatchRules{
					v3.MatchRuleER{
						Name:  "TestName",
						Type:  v3.MatchRuleTypeER,
						Start: 7,
						End:   8,
						ID:    789,
						Matches: []v3.MatchCriteriaER{
							{
								MatchType:        "TestType",
								MatchValue:       "TestValue",
								MatchOperator:    "TestOperator",
								CaseSensitive:    true,
								Negate:           true,
								CheckIPs:         "1.1.1.1",
								ObjectMatchValue: "1",
							},
						},
						MatchesAlways:            true,
						UseRelativeURL:           "/url1",
						StatusCode:               200,
						RedirectURL:              "/url2",
						MatchURL:                 "/url3",
						UseIncomingQueryString:   true,
						UseIncomingSchemeAndHost: true,
						Disabled:                 true,
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
			},
			init: func(m *v3.Mock, data testDataForSharedPolicy) {
				mockGetPolicy(m, data, 5)
				mockGetPolicyVersion(m, data, 5)
			},
		},
		"success with no version attribute - no activations and no match rules and warnings": {
			config: "no_version.tf",
			data: testDataForSharedPolicy{
				id:           "akamai_cloudlets_shared_policy",
				policyID:     1,
				version:      2,
				groupID:      12,
				name:         "TestName",
				cloudletType: v3.CloudletTypeAP,
				description:  "Description",
			},
			init: func(m *v3.Mock, data testDataForSharedPolicy) {
				mockGetPolicy(m, data, 5)
				mockListPolicyVersions(m, data, 2, 5)
				mockGetPolicyVersion(m, data, 5)
			},
		},
		"success with no version attribute - no shared policy versions": {
			config: "no_version.tf",
			data: testDataForSharedPolicy{
				id:           "akamai_cloudlets_shared_policy",
				policyID:     1,
				version:      0,
				groupID:      12,
				name:         "TestName",
				cloudletType: v3.CloudletTypeAP,
				description:  "Description",
			},
			init: func(m *v3.Mock, data testDataForSharedPolicy) {
				mockGetPolicy(m, data, 5)
				mockListPolicyVersions(m, data, 2, 5)
			},
		},
		"success with version attribute - all activations": {
			config: "with_version.tf",
			data: testDataForSharedPolicy{
				id:           "akamai_cloudlets_shared_policy",
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
			init: func(m *v3.Mock, data testDataForSharedPolicy) {
				mockGetPolicy(m, data, 5)
				mockGetPolicyVersion(m, data, 5)
			},
		},
		"expect error on ListPolicyVersions": {
			config: "no_version.tf",
			data: testDataForSharedPolicy{
				id:           "akamai_cloudlets_shared_policy",
				policyID:     1,
				version:      2,
				groupID:      12,
				name:         "TestName",
				cloudletType: v3.CloudletTypeAP,
				description:  "Description",
			},
			init: func(m *v3.Mock, data testDataForSharedPolicy) {
				mockGetPolicy(m, data, 1)
				m.On("ListPolicyVersions", mock.Anything, v3.ListPolicyVersionsRequest{
					PolicyID: data.policyID,
				}).Return(nil, fmt.Errorf("API error")).Once()
			},
			expectError: regexp.MustCompile(`Error: Reading Cloudlets Shared Policy Failed`),
		},
		"expect ErrNotFound error on GetPolicy": {
			config: "with_version.tf",
			data: testDataForSharedPolicy{
				policyID: 1,
			},
			init: func(m *v3.Mock, data testDataForSharedPolicy) {
				m.On("GetPolicy", mock.Anything, v3.GetPolicyRequest{
					PolicyID: data.policyID,
				}).Return(nil, fmt.Errorf("%s: %w: %s", v3.ErrGetPolicy, v3.ErrPolicyNotFound, "oops")).Once()
			},
			expectError: regexp.MustCompile(`Error: Policy does not exist or is not of 'SHARED' type`),
		},
		"expect error - missing required attribute": {
			config:      "no_policy_id.tf",
			data:        testDataForSharedPolicy{},
			init:        func(m *v3.Mock, data testDataForSharedPolicy) {},
			expectError: regexp.MustCompile(`The argument "policy_id" is required, but no definition was found.`),
		},
	}

	for name, test := range tests {
		t.Run(name, func(t *testing.T) {
			client := &v3.Mock{}
			if test.init != nil {
				test.init(client, test.data)
			}
			useClientV3(client, func() {
				resource.Test(t, resource.TestCase{
					IsUnitTest:               true,
					ProtoV6ProviderFactories: testutils.NewProtoV6ProviderFactory(NewSubprovider()),
					Steps: []resource.TestStep{{
						Config:      testutils.LoadFixtureString(t, fmt.Sprintf("testdata/TestDataCloudletsSharedPolicy/%s", test.config)),
						Check:       checkAttrsForSharedPolicy(test.data),
						ExpectError: test.expectError,
					}},
				})
			})
			client.AssertExpectations(t)
		})
	}
}

func checkAttrsForSharedPolicy(data testDataForSharedPolicy) resource.TestCheckFunc {
	var checkFuncs []resource.TestCheckFunc
	matchRules, err := json.Marshal(data.matchRules)
	if err != nil {
		panic(err)
	}
	matchWarnings, err := json.Marshal(data.warnings)
	if err != nil {
		panic(err)
	}
	if data.version == 0 {
		checkFuncs = append(checkFuncs,
			resource.TestCheckNoResourceAttr("data.akamai_cloudlets_shared_policy.test", "version"),
			resource.TestCheckNoResourceAttr("data.akamai_cloudlets_shared_policy.test", "match_rules"),
			resource.TestCheckNoResourceAttr("data.akamai_cloudlets_shared_policy.test", "warnings"))
	} else {
		if data.versionDescription != nil {
			checkFuncs = append(checkFuncs,
				resource.TestCheckResourceAttr("data.akamai_cloudlets_shared_policy.test", "version_description", *data.versionDescription))
		} else {
			checkFuncs = append(checkFuncs,
				resource.TestCheckNoResourceAttr("data.akamai_cloudlets_shared_policy.test", "version_description"))
		}
		checkFuncs = append(checkFuncs,
			resource.TestCheckResourceAttr("data.akamai_cloudlets_shared_policy.test", "version", strconv.FormatInt(data.version, 10)),
			resource.TestCheckResourceAttr("data.akamai_cloudlets_shared_policy.test", "match_rules", string(matchRules)),
			resource.TestCheckResourceAttr("data.akamai_cloudlets_shared_policy.test", "warnings", string(matchWarnings)))
	}

	checkFuncs = append(checkFuncs,
		resource.TestCheckResourceAttr("data.akamai_cloudlets_shared_policy.test", "name", data.name),
		resource.TestCheckResourceAttr("data.akamai_cloudlets_shared_policy.test", "id", data.id),
		resource.TestCheckResourceAttr("data.akamai_cloudlets_shared_policy.test", "group_id", strconv.FormatInt(data.groupID, 10)),
		resource.TestCheckResourceAttr("data.akamai_cloudlets_shared_policy.test", "description", data.description),
	)

	if data.activations.Production.Effective != nil {
		checkFuncs = append(checkFuncs, checkActivationAttributesForSharedPolicy("production", "effective", data.activations.Production.Effective))
	}
	if data.activations.Production.Latest != nil {
		checkFuncs = append(checkFuncs, checkActivationAttributesForSharedPolicy("production", "latest", data.activations.Production.Latest))
	}
	if data.activations.Staging.Effective != nil {
		checkFuncs = append(checkFuncs, checkActivationAttributesForSharedPolicy("staging", "effective", data.activations.Staging.Effective))
	}
	if data.activations.Staging.Latest != nil {
		checkFuncs = append(checkFuncs, checkActivationAttributesForSharedPolicy("staging", "latest", data.activations.Staging.Latest))
	}

	return resource.ComposeAggregateTestCheckFunc(checkFuncs...)
}

func checkActivationAttributesForSharedPolicy(actNetwork, actInfo string, actData *v3.PolicyActivation) resource.TestCheckFunc {
	return resource.ComposeAggregateTestCheckFunc(
		resource.TestCheckResourceAttr("data.akamai_cloudlets_shared_policy.test", fmt.Sprintf("activations.%s.%s.created_by", actNetwork, actInfo), actData.CreatedBy),
		resource.TestCheckResourceAttr("data.akamai_cloudlets_shared_policy.test", fmt.Sprintf("activations.%s.%s.created_date", actNetwork, actInfo), actData.CreatedDate.String()),
		resource.TestCheckResourceAttr("data.akamai_cloudlets_shared_policy.test", fmt.Sprintf("activations.%s.%s.network", actNetwork, actInfo), string(actData.Network)),
		resource.TestCheckResourceAttr("data.akamai_cloudlets_shared_policy.test", fmt.Sprintf("activations.%s.%s.status", actNetwork, actInfo), string(actData.Status)),
		resource.TestCheckResourceAttr("data.akamai_cloudlets_shared_policy.test", fmt.Sprintf("activations.%s.%s.policy_version", actNetwork, actInfo), strconv.FormatInt(actData.PolicyVersion, 10)),
		resource.TestCheckResourceAttr("data.akamai_cloudlets_shared_policy.test", fmt.Sprintf("activations.%s.%s.finish_date", actNetwork, actInfo), actData.FinishDate.String()),
		resource.TestCheckResourceAttr("data.akamai_cloudlets_shared_policy.test", fmt.Sprintf("activations.%s.%s.operation", actNetwork, actInfo), string(actData.Operation)),
		resource.TestCheckResourceAttr("data.akamai_cloudlets_shared_policy.test", fmt.Sprintf("activations.%s.%s.policy_id", actNetwork, actInfo), strconv.FormatInt(actData.PolicyID, 10)),
		resource.TestCheckResourceAttr("data.akamai_cloudlets_shared_policy.test", fmt.Sprintf("activations.%s.%s.activation_id", actNetwork, actInfo), strconv.FormatInt(actData.ID, 10)))
}

func mockGetPolicy(m *v3.Mock, data testDataForSharedPolicy, times int) {
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

func mockGetPolicyVersion(m *v3.Mock, data testDataForSharedPolicy, times int) {
	m.On("GetPolicyVersion", mock.Anything, v3.GetPolicyVersionRequest{
		PolicyID:      data.policyID,
		PolicyVersion: data.version,
	}).Return(&v3.PolicyVersion{
		PolicyID:           data.policyID,
		PolicyVersion:      data.version,
		Description:        data.versionDescription,
		ID:                 123,
		Immutable:          false,
		MatchRules:         data.matchRules,
		MatchRulesWarnings: data.warnings,
	}, nil).Times(times)
}

func createPolicyVersions(policyID int64, numberOfVersions, pageNumber int) *v3.ListPolicyVersions {
	var policyVersions v3.ListPolicyVersions
	for i := numberOfVersions; i > 0; i-- {
		policyVersions.PolicyVersions = append(policyVersions.PolicyVersions, v3.ListPolicyVersionsItem{
			Description:   ptr.To(fmt.Sprintf("Description%d", i)),
			ID:            int64(i),
			Immutable:     true,
			PolicyID:      policyID,
			PolicyVersion: int64(i),
		})
	}
	policyVersions.Page.Number = pageNumber

	return &policyVersions
}

func mockListPolicyVersions(m *v3.Mock, data testDataForSharedPolicy, numberOfActivations, times int) {
	if data.version == 0 {
		m.On("ListPolicyVersions", mock.Anything, v3.ListPolicyVersionsRequest{
			PolicyID: data.policyID,
		}).Return(&v3.ListPolicyVersions{
			PolicyVersions: []v3.ListPolicyVersionsItem{},
		}, nil).Times(times)
	} else {
		m.On("ListPolicyVersions", mock.Anything, v3.ListPolicyVersionsRequest{
			PolicyID: data.policyID,
		}).Return(createPolicyVersions(data.policyID, numberOfActivations, 0), nil).Times(times)
	}
}
