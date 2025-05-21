package cloudlets

import (
	"fmt"
	"regexp"
	"testing"

	"github.com/akamai/AkamaiOPEN-edgegrid-golang/v11/pkg/cloudlets"
	"github.com/akamai/terraform-provider-akamai/v8/pkg/common/ptr"
	"github.com/akamai/terraform-provider-akamai/v8/pkg/common/testutils"
	"github.com/hashicorp/terraform-plugin-testing/helper/resource"
)

func TestDataCloudletsPolicy(t *testing.T) {
	getPolicyReturn := cloudlets.Policy{
		Location:         "/cloudlets/api/v2/policies/1234",
		PolicyID:         1234,
		GroupID:          2345,
		Name:             "SomeName",
		Description:      "Fancy Description",
		CreatedBy:        "jsmith",
		CreateDate:       1631190136928,
		LastModifiedBy:   "jsmith",
		LastModifiedDate: 1631190136928,
		CloudletID:       0,
		CloudletCode:     "ER",
		APIVersion:       "2.0",
		Deleted:          false,
	}

	getPolicyVersionReturn := cloudlets.PolicyVersion{
		Location:         "/cloudlets/api/v2/policies/1234/version/1",
		RevisionID:       4824132,
		PolicyID:         1234,
		Version:          1,
		Description:      "Example Description",
		CreatedBy:        "jsmith",
		CreateDate:       1631191583350,
		LastModifiedBy:   "jsmith",
		LastModifiedDate: 1631191583352,
		RulesLocked:      false,
		MatchRules: cloudlets.MatchRules{
			cloudlets.MatchRuleER{
				Name:                     "rule 2",
				Type:                     "erMatchRule",
				UseRelativeURL:           "none",
				StatusCode:               301,
				RedirectURL:              "ss.exmaple.com",
				MatchURL:                 "aa.exmaple.com",
				UseIncomingQueryString:   true,
				UseIncomingSchemeAndHost: false,
			},
		},
		MatchRuleFormat: "1.0",
		Deleted:         false,
		Warnings:        nil,
	}

	listPoliciesReturn1000 := make([]cloudlets.Policy, 1000)
	for i := 0; i < 1000; i++ {
		listPoliciesReturn1000[i] = cloudlets.Policy{
			Location:         fmt.Sprintf("/cloudlets/api/v2/policies/%d", i),
			PolicyID:         1000 + int64(i),
			GroupID:          3000 + int64(i),
			Name:             fmt.Sprintf("test_policy_%d", i),
			Description:      "Fancy Description",
			CreatedBy:        "jsmith",
			CreateDate:       1631190136928,
			LastModifiedBy:   "jsmith",
			LastModifiedDate: 1631190136928,
			Activations:      nil,
			CloudletID:       0,
			CloudletCode:     "ER",
			APIVersion:       "2.0",
			Deleted:          false,
		}
	}

	listPoliciesReturn100 := make([]cloudlets.Policy, 100)
	for i := 0; i < 100; i++ {
		listPoliciesReturn100[i] = cloudlets.Policy{
			Location:         fmt.Sprintf("/cloudlets/api/v2/policies/%d", i+1000),
			PolicyID:         2000 + int64(i),
			GroupID:          4000 + int64(i),
			Name:             fmt.Sprintf("test_policy_%d", i+1000),
			Description:      "Fancy Description",
			CreatedBy:        "jsmith",
			CreateDate:       1631190136928,
			LastModifiedBy:   "jsmith",
			LastModifiedDate: 1631190136928,
			Activations:      nil,
			CloudletID:       0,
			CloudletCode:     "ER",
			APIVersion:       "2.0",
			Deleted:          false,
		}
	}

	tests := map[string]struct {
		configPath string
		init       func(*cloudlets.Mock)
		checkFuncs []resource.TestCheckFunc
		withError  *regexp.Regexp
	}{
		"validate basic schema": {
			configPath: "testdata/TestDataCloudletsPolicy/policy.tf",
			init: func(m *cloudlets.Mock) {
				m.On("ListPolicyVersions", testutils.MockContext, cloudlets.ListPolicyVersionsRequest{
					PolicyID: 1234,
					PageSize: ptr.To(pageSize),
					Offset:   0,
				}).Return([]cloudlets.PolicyVersion{{Version: 1}}, nil).Times(3)

				m.On("GetPolicyVersion", testutils.MockContext, cloudlets.GetPolicyVersionRequest{
					PolicyID: 1234,
					Version:  1,
				}).Return(&getPolicyVersionReturn, nil).Times(3)

				m.On("GetPolicy", testutils.MockContext, cloudlets.GetPolicyRequest{
					PolicyID: 1234,
				}).Return(&getPolicyReturn, nil).Times(3)
			},
			checkFuncs: []resource.TestCheckFunc{
				resource.TestCheckResourceAttr("data.akamai_cloudlets_policy.test", "policy_id", "1234"),
				resource.TestCheckResourceAttr("data.akamai_cloudlets_policy.test", "version", "1"),
				resource.TestCheckResourceAttr("data.akamai_cloudlets_policy.test", "group_id", "2345"),
				resource.TestCheckResourceAttr("data.akamai_cloudlets_policy.test", "name", "SomeName"),
				resource.TestCheckResourceAttr("data.akamai_cloudlets_policy.test", "description", "Fancy Description"),
				resource.TestCheckResourceAttr("data.akamai_cloudlets_policy.test", "version_description", "Example Description"),
				resource.TestCheckResourceAttr("data.akamai_cloudlets_policy.test", "cloudlet_id", "0"),
				resource.TestCheckResourceAttr("data.akamai_cloudlets_policy.test", "cloudlet_code", "ER"),
				resource.TestCheckResourceAttr("data.akamai_cloudlets_policy.test", "api_version", "2.0"),
				resource.TestCheckResourceAttr("data.akamai_cloudlets_policy.test", "revision_id", "4824132"),
				resource.TestCheckResourceAttr("data.akamai_cloudlets_policy.test", "rules_locked", "false"),
				resource.TestCheckResourceAttr("data.akamai_cloudlets_policy.test", "match_rules", testutils.LoadFixtureString(t, "testdata/TestDataCloudletsPolicy/rules/match_rules_out.json")),
				resource.TestCheckResourceAttr("data.akamai_cloudlets_policy.test", "match_rule_format", "1.0"),
				resource.TestCheckResourceAttr("data.akamai_cloudlets_policy.test", "warnings", "null"),
			},
		},
		"validate activations schema": {
			configPath: "testdata/TestDataCloudletsPolicy/policy.tf",
			init: func(m *cloudlets.Mock) {
				getPolicy3VersionsReturn := cloudlets.PolicyVersion{
					Location:         "/cloudlets/api/v2/policies/1234/version/1",
					RevisionID:       4824132,
					PolicyID:         1234,
					Version:          3,
					Description:      "Example Description",
					CreatedBy:        "jsmith",
					CreateDate:       1631191583350,
					LastModifiedBy:   "jsmith",
					LastModifiedDate: 1631191583352,
					RulesLocked:      false,
					Activations: []cloudlets.PolicyActivation{
						{
							APIVersion: "2.0",
							Network:    "prod",
							PolicyInfo: cloudlets.PolicyInfo{
								PolicyID:       1234,
								Name:           "policy_name_0",
								Version:        3,
								Status:         "active",
								StatusDetail:   "",
								ActivatedBy:    "jsmith",
								ActivationDate: 1607507783000,
							},
							PropertyInfo: cloudlets.PropertyInfo{
								Name:           "property_name_0",
								Version:        3,
								GroupID:        132,
								Status:         "active",
								ActivatedBy:    "jsmith",
								ActivationDate: 1607507783812,
							},
						},
						{
							APIVersion: "2.0",
							Network:    "stage",
							PolicyInfo: cloudlets.PolicyInfo{
								PolicyID:       1234,
								Name:           "policy_name_1",
								Version:        3,
								Status:         "active",
								StatusDetail:   "",
								ActivatedBy:    "jsmith",
								ActivationDate: 1607507783001,
							},
							PropertyInfo: cloudlets.PropertyInfo{
								Name:           "property_name_1",
								Version:        4,
								GroupID:        133,
								Status:         "active",
								ActivatedBy:    "jsmith",
								ActivationDate: 1607507783813,
							},
						},
					},
					MatchRules: cloudlets.MatchRules{
						cloudlets.MatchRuleER{
							Name:                     "rule 2",
							Type:                     "erMatchRule",
							UseRelativeURL:           "none",
							StatusCode:               301,
							RedirectURL:              "ss.exmaple.com",
							MatchURL:                 "aa.exmaple.com",
							UseIncomingQueryString:   true,
							UseIncomingSchemeAndHost: false,
						},
					},
					MatchRuleFormat: "1.0",
					Deleted:         false,
					Warnings:        nil,
				}

				m.On("ListPolicyVersions", testutils.MockContext, cloudlets.ListPolicyVersionsRequest{
					PolicyID: 1234,
					Offset:   0,
					PageSize: ptr.To(pageSize),
				}).Return([]cloudlets.PolicyVersion{{Version: 3}, {Version: 2}, {Version: 1}}, nil).Times(3)

				m.On("GetPolicyVersion", testutils.MockContext, cloudlets.GetPolicyVersionRequest{
					PolicyID: 1234,
					Version:  3,
				}).Return(&getPolicy3VersionsReturn, nil).Times(3)

				m.On("GetPolicy", testutils.MockContext, cloudlets.GetPolicyRequest{
					PolicyID: 1234,
				}).Return(&getPolicyReturn, nil).Times(3)
			},
			checkFuncs: []resource.TestCheckFunc{
				resource.TestCheckResourceAttr("data.akamai_cloudlets_policy.test", "activations.#", "2"),
				resource.TestCheckResourceAttr("data.akamai_cloudlets_policy.test", "activations.0.api_version", "2.0"),
				resource.TestCheckResourceAttr("data.akamai_cloudlets_policy.test", "activations.0.network", "prod"),

				resource.TestCheckResourceAttr("data.akamai_cloudlets_policy.test", "activations.0.policy_info.#", "1"),
				resource.TestCheckResourceAttr("data.akamai_cloudlets_policy.test", "activations.0.policy_info.0.policy_id", "1234"),
				resource.TestCheckResourceAttr("data.akamai_cloudlets_policy.test", "activations.0.policy_info.0.name", "policy_name_0"),
				resource.TestCheckResourceAttr("data.akamai_cloudlets_policy.test", "activations.0.policy_info.0.version", "3"),
				resource.TestCheckResourceAttr("data.akamai_cloudlets_policy.test", "activations.0.policy_info.0.status", "active"),
				resource.TestCheckResourceAttr("data.akamai_cloudlets_policy.test", "activations.0.policy_info.0.status_detail", ""),
				resource.TestCheckResourceAttr("data.akamai_cloudlets_policy.test", "activations.0.policy_info.0.activated_by", "jsmith"),
				resource.TestCheckResourceAttr("data.akamai_cloudlets_policy.test", "activations.0.policy_info.0.activation_date", "1607507783000"),

				resource.TestCheckResourceAttr("data.akamai_cloudlets_policy.test", "activations.0.property_info.#", "1"),
				resource.TestCheckResourceAttr("data.akamai_cloudlets_policy.test", "activations.0.property_info.0.name", "property_name_0"),
				resource.TestCheckResourceAttr("data.akamai_cloudlets_policy.test", "activations.0.property_info.0.version", "3"),
				resource.TestCheckResourceAttr("data.akamai_cloudlets_policy.test", "activations.0.property_info.0.group_id", "132"),
				resource.TestCheckResourceAttr("data.akamai_cloudlets_policy.test", "activations.0.property_info.0.status", "active"),
				resource.TestCheckResourceAttr("data.akamai_cloudlets_policy.test", "activations.0.property_info.0.activated_by", "jsmith"),
				resource.TestCheckResourceAttr("data.akamai_cloudlets_policy.test", "activations.0.property_info.0.activation_date", "1607507783812"),

				resource.TestCheckResourceAttr("data.akamai_cloudlets_policy.test", "activations.1.api_version", "2.0"),
				resource.TestCheckResourceAttr("data.akamai_cloudlets_policy.test", "activations.1.network", "stage"),

				resource.TestCheckResourceAttr("data.akamai_cloudlets_policy.test", "activations.1.policy_info.#", "1"),
				resource.TestCheckResourceAttr("data.akamai_cloudlets_policy.test", "activations.1.policy_info.0.policy_id", "1234"),
				resource.TestCheckResourceAttr("data.akamai_cloudlets_policy.test", "activations.1.policy_info.0.name", "policy_name_1"),
				resource.TestCheckResourceAttr("data.akamai_cloudlets_policy.test", "activations.1.policy_info.0.version", "3"),
				resource.TestCheckResourceAttr("data.akamai_cloudlets_policy.test", "activations.1.policy_info.0.status", "active"),
				resource.TestCheckResourceAttr("data.akamai_cloudlets_policy.test", "activations.1.policy_info.0.status_detail", ""),
				resource.TestCheckResourceAttr("data.akamai_cloudlets_policy.test", "activations.1.policy_info.0.activated_by", "jsmith"),
				resource.TestCheckResourceAttr("data.akamai_cloudlets_policy.test", "activations.1.policy_info.0.activation_date", "1607507783001"),

				resource.TestCheckResourceAttr("data.akamai_cloudlets_policy.test", "activations.1.property_info.#", "1"),
				resource.TestCheckResourceAttr("data.akamai_cloudlets_policy.test", "activations.1.property_info.0.name", "property_name_1"),
				resource.TestCheckResourceAttr("data.akamai_cloudlets_policy.test", "activations.1.property_info.0.version", "4"),
				resource.TestCheckResourceAttr("data.akamai_cloudlets_policy.test", "activations.1.property_info.0.group_id", "133"),
				resource.TestCheckResourceAttr("data.akamai_cloudlets_policy.test", "activations.1.property_info.0.status", "active"),
				resource.TestCheckResourceAttr("data.akamai_cloudlets_policy.test", "activations.1.property_info.0.activated_by", "jsmith"),
				resource.TestCheckResourceAttr("data.akamai_cloudlets_policy.test", "activations.1.property_info.0.activation_date", "1607507783813"),
			},
		},
		"pass version in tf file": {
			configPath: "testdata/TestDataCloudletsPolicy/policy_with_version.tf",
			init: func(m *cloudlets.Mock) {
				m.On("GetPolicyVersion", testutils.MockContext, cloudlets.GetPolicyVersionRequest{
					PolicyID: 1234,
					Version:  3,
				}).Return(&getPolicyVersionReturn, nil).Times(3)

				m.On("GetPolicy", testutils.MockContext, cloudlets.GetPolicyRequest{
					PolicyID: 1234,
				}).Return(&getPolicyReturn, nil).Times(3)
			},
			checkFuncs: []resource.TestCheckFunc{
				resource.TestCheckResourceAttr("data.akamai_cloudlets_policy.test", "version", "1"),
			},
		},
		"pass name instead of policy ID in tf file": {
			configPath: "testdata/TestDataCloudletsPolicy/policy_only_with_name.tf",
			init: func(m *cloudlets.Mock) {
				m.On("ListPolicies", testutils.MockContext, cloudlets.ListPoliciesRequest{
					Offset:   0,
					PageSize: ptr.To(pageSize),
				}).Return([]cloudlets.Policy{
					{
						Location:         "/cloudlets/api/v2/policies/1234",
						PolicyID:         1234,
						GroupID:          2345,
						Name:             "test_policy",
						Description:      "Fancy Description",
						CreatedBy:        "jsmith",
						CreateDate:       1631190136928,
						LastModifiedBy:   "jsmith",
						LastModifiedDate: 1631190136928,
						Activations:      nil,
						CloudletID:       0,
						CloudletCode:     "ER",
						APIVersion:       "2.0",
						Deleted:          false,
					},
				}, nil).Times(3)

				m.On("ListPolicyVersions", testutils.MockContext, cloudlets.ListPolicyVersionsRequest{
					PolicyID: 1234,
					Offset:   0,
					PageSize: ptr.To(pageSize),
				}).Return([]cloudlets.PolicyVersion{{Version: 1}}, nil).Times(3)

				m.On("GetPolicyVersion", testutils.MockContext, cloudlets.GetPolicyVersionRequest{
					PolicyID:  1234,
					Version:   1,
					OmitRules: false,
				}).Return(&getPolicyVersionReturn, nil).Times(3)

				m.On("GetPolicy", testutils.MockContext, cloudlets.GetPolicyRequest{
					PolicyID: 1234,
				}).Return(&cloudlets.Policy{
					Location:         "/cloudlets/api/v2/policies/1234",
					PolicyID:         1234,
					GroupID:          2345,
					Name:             "test_policy",
					Description:      "Fancy Description",
					CreatedBy:        "jsmith",
					CreateDate:       1631190136928,
					LastModifiedBy:   "jsmith",
					LastModifiedDate: 1631190136928,
					CloudletID:       0,
					CloudletCode:     "ER",
					APIVersion:       "2.0",
					Deleted:          false,
				}, nil).Times(3)
			},
			checkFuncs: []resource.TestCheckFunc{
				resource.TestCheckResourceAttr("data.akamai_cloudlets_policy.test", "version", "1"),
				resource.TestCheckResourceAttr("data.akamai_cloudlets_policy.test", "group_id", "2345"),
				resource.TestCheckResourceAttr("data.akamai_cloudlets_policy.test", "name", "test_policy"),
				resource.TestCheckResourceAttr("data.akamai_cloudlets_policy.test", "policy_id", "1234"),
				resource.TestCheckResourceAttr("data.akamai_cloudlets_policy.test", "description", "Fancy Description"),
				resource.TestCheckResourceAttr("data.akamai_cloudlets_policy.test", "version_description", "Example Description"),
				resource.TestCheckResourceAttr("data.akamai_cloudlets_policy.test", "cloudlet_id", "0"),
				resource.TestCheckResourceAttr("data.akamai_cloudlets_policy.test", "cloudlet_code", "ER"),
				resource.TestCheckResourceAttr("data.akamai_cloudlets_policy.test", "api_version", "2.0"),
				resource.TestCheckResourceAttr("data.akamai_cloudlets_policy.test", "revision_id", "4824132"),
				resource.TestCheckResourceAttr("data.akamai_cloudlets_policy.test", "rules_locked", "false"),
				resource.TestCheckResourceAttr("data.akamai_cloudlets_policy.test", "match_rules", testutils.LoadFixtureString(t, "testdata/TestDataCloudletsPolicy/rules/match_rules_out.json")),
				resource.TestCheckResourceAttr("data.akamai_cloudlets_policy.test", "match_rule_format", "1.0"),
				resource.TestCheckResourceAttr("data.akamai_cloudlets_policy.test", "warnings", "null"),
			},
		},
		"policy name on the second page": {
			configPath: "testdata/TestDataCloudletsPolicy/policy_only_with_name.tf",
			init: func(m *cloudlets.Mock) {
				m.On("ListPolicies", testutils.MockContext, cloudlets.ListPoliciesRequest{
					Offset:   0,
					PageSize: ptr.To(pageSize),
				}).Return(listPoliciesReturn1000, nil).Times(3)

				listPoliciesReturn100[10].Name = "test_policy"
				m.On("ListPolicies", testutils.MockContext, cloudlets.ListPoliciesRequest{
					Offset:   1000,
					PageSize: ptr.To(pageSize),
				}).Return(listPoliciesReturn100, nil).Times(3)

				m.On("ListPolicyVersions", testutils.MockContext, cloudlets.ListPolicyVersionsRequest{
					PolicyID: 2010,
					Offset:   0,
					PageSize: ptr.To(pageSize),
				}).Return([]cloudlets.PolicyVersion{{Version: 1}}, nil).Times(3)

				m.On("GetPolicyVersion", testutils.MockContext, cloudlets.GetPolicyVersionRequest{
					PolicyID:  2010,
					Version:   1,
					OmitRules: false,
				}).Return(&cloudlets.PolicyVersion{
					Location:         "/cloudlets/api/v2/policies/2010",
					RevisionID:       2010,
					PolicyID:         2010,
					Version:          1,
					Description:      "Fancy Description",
					CreatedBy:        "jsmith",
					CreateDate:       1631190136928,
					LastModifiedBy:   "jsmith",
					LastModifiedDate: 1631190136928,
					RulesLocked:      false,
					Activations:      nil,
					MatchRules:       nil,
					MatchRuleFormat:  "",
					Deleted:          false,
					Warnings:         nil,
				}, nil).Times(3)

				m.On("GetPolicy", testutils.MockContext, cloudlets.GetPolicyRequest{
					PolicyID: 2010,
				}).Return(&cloudlets.Policy{
					Location:         "/cloudlets/api/v2/policies/2010",
					PolicyID:         2010,
					GroupID:          4010,
					Name:             "test_policy",
					Description:      "Fancy Description",
					CreatedBy:        "jsmith",
					CreateDate:       1631190136928,
					LastModifiedBy:   "jsmith",
					LastModifiedDate: 1631190136928,
					CloudletID:       0,
					CloudletCode:     "ER",
					APIVersion:       "2.0",
					Deleted:          false,
				}, nil).Times(3)
			},
		},
		"deleted policy version": {
			configPath: "testdata/TestDataCloudletsPolicy/policy_with_version.tf",
			init: func(m *cloudlets.Mock) {
				getPolicyVersionReturnDeleted := cloudlets.PolicyVersion{
					Version: 3,
					Deleted: true,
				}

				m.On("GetPolicyVersion", testutils.MockContext, cloudlets.GetPolicyVersionRequest{
					PolicyID:  1234,
					Version:   3,
					OmitRules: false,
				}).Return(&getPolicyVersionReturnDeleted, nil).Once()
			},
			withError: regexp.MustCompile("specified policy version is deleted"),
		},
		"deleted policy": {
			configPath: "testdata/TestDataCloudletsPolicy/policy.tf",
			init: func(m *cloudlets.Mock) {
				getPolicyReturnDeleted := cloudlets.Policy{
					Deleted: true,
				}
				m.On("ListPolicyVersions", testutils.MockContext, cloudlets.ListPolicyVersionsRequest{
					PolicyID: 1234,
					Offset:   0,
					PageSize: ptr.To(pageSize),
				}).Return([]cloudlets.PolicyVersion{}, nil).Once()

				m.On("GetPolicyVersion", testutils.MockContext, cloudlets.GetPolicyVersionRequest{
					PolicyID: 1234,
				}).Return(&getPolicyVersionReturn, nil).Once()

				m.On("GetPolicy", testutils.MockContext, cloudlets.GetPolicyRequest{
					PolicyID: 1234,
				}).Return(&getPolicyReturnDeleted, nil).Once()
			},
			withError: regexp.MustCompile("specified policy is deleted"),
		},
		"policy not found for the given name": {
			configPath: "testdata/TestDataCloudletsPolicy/policy_only_with_name.tf",
			init: func(m *cloudlets.Mock) {
				listPolicyReturn := []cloudlets.Policy{
					{
						Location:         "/cloudlets/api/v2/policies/1234",
						PolicyID:         1234,
						GroupID:          2345,
						Name:             "SomeName1",
						Description:      "Fancy Description",
						CreatedBy:        "jsmith",
						CreateDate:       1631190136928,
						LastModifiedBy:   "jsmith",
						LastModifiedDate: 1631190136928,
						CloudletID:       0,
						CloudletCode:     "ER",
						APIVersion:       "2.0",
						Deleted:          false,
					},
					{
						Location:         "/cloudlets/api/v2/policies/12354",
						PolicyID:         12345,
						GroupID:          23455,
						Name:             "SomeName2",
						Description:      "Fancy Description",
						CreatedBy:        "jsmith",
						CreateDate:       1631190136928,
						LastModifiedBy:   "jkowalski",
						LastModifiedDate: 1631190136928,
						CloudletID:       0,
						CloudletCode:     "ER",
						APIVersion:       "2.0",
						Deleted:          false,
					},
				}

				m.On("ListPolicies", testutils.MockContext, cloudlets.ListPoliciesRequest{
					Offset:   0,
					PageSize: ptr.To(pageSize),
				}).Return(listPolicyReturn, nil).Once()
			},
			withError: regexp.MustCompile("policy not found: test_policy"),
		},
		"config without policy ID and name": {
			configPath: "testdata/TestDataCloudletsPolicy/policy_without_name_and_policy_id.tf",
			withError:  regexp.MustCompile("one of `name,policy_id` must be specified"),
		},
		"config with both policy ID and name": {
			configPath: "testdata/TestDataCloudletsPolicy/policy_with_name_and_policy_id.tf",
			withError:  regexp.MustCompile("only one of `name,policy_id` can be specified, but `name,policy_id`\nwere specified"),
		},
		"no version for a policy": {
			configPath: "testdata/TestDataCloudletsPolicy/policy.tf",
			init: func(m *cloudlets.Mock) {
				m.On("ListPolicyVersions", testutils.MockContext, cloudlets.ListPolicyVersionsRequest{
					PolicyID: 1234,
					Offset:   0,
					PageSize: ptr.To(pageSize),
				}).Return([]cloudlets.PolicyVersion{}, nil).Times(3)

				m.On("GetPolicy", testutils.MockContext, cloudlets.GetPolicyRequest{
					PolicyID: 1234,
				}).Return(&getPolicyReturn, nil).Times(3)
			},
			checkFuncs: []resource.TestCheckFunc{
				resource.TestCheckResourceAttr("data.akamai_cloudlets_policy.test", "id", "1234"),
				resource.TestCheckNoResourceAttr("data.akamai_cloudlets_policy.test", "version"),
				resource.TestCheckResourceAttr("data.akamai_cloudlets_policy.test", "group_id", "2345"),
				resource.TestCheckResourceAttr("data.akamai_cloudlets_policy.test", "name", "SomeName"),
				resource.TestCheckResourceAttr("data.akamai_cloudlets_policy.test", "description", "Fancy Description"),
				resource.TestCheckNoResourceAttr("data.akamai_cloudlets_policy.test", "version_description"),
				resource.TestCheckResourceAttr("data.akamai_cloudlets_policy.test", "cloudlet_id", "0"),
				resource.TestCheckResourceAttr("data.akamai_cloudlets_policy.test", "cloudlet_code", "ER"),
				resource.TestCheckResourceAttr("data.akamai_cloudlets_policy.test", "api_version", "2.0"),
				resource.TestCheckNoResourceAttr("data.akamai_cloudlets_policy.test", "revision_id"),
				resource.TestCheckNoResourceAttr("data.akamai_cloudlets_policy.test", "rules_locked"),
				resource.TestCheckNoResourceAttr("data.akamai_cloudlets_policy.test", "match_rules"),
				resource.TestCheckNoResourceAttr("data.akamai_cloudlets_policy.test", "match_rule_format"),
				resource.TestCheckNoResourceAttr("data.akamai_cloudlets_policy.test", "warnings"),
			},
		},
	}

	for testName, test := range tests {
		t.Run(testName, func(t *testing.T) {
			client := &cloudlets.Mock{}
			if test.init != nil {
				test.init(client)
			}
			useClient(client, func() {
				resource.UnitTest(t, resource.TestCase{
					ProtoV6ProviderFactories: testutils.NewProtoV6ProviderFactory(NewSubprovider()),
					Steps: []resource.TestStep{
						{
							Config: testutils.LoadFixtureString(t, test.configPath),
							Check: resource.ComposeAggregateTestCheckFunc(
								test.checkFuncs...,
							),
							ExpectError: test.withError,
						},
					},
				})
			})
		})
	}
}
