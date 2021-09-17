package cloudlets

import (
	"regexp"
	"testing"

	"github.com/stretchr/testify/mock"

	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/resource"

	"github.com/akamai/AkamaiOPEN-edgegrid-golang/v2/pkg/cloudlets"
)

func TestDataCloudletsPolicy(t *testing.T) {
	tests := map[string]struct {
		configPath               string
		listPolicyVersionsReturn []cloudlets.PolicyVersion
		getPolicyReturn          cloudlets.Policy
		getPolicyVersionReturn   cloudlets.PolicyVersion
		checkFuncs               []resource.TestCheckFunc
		withError                *regexp.Regexp
	}{
		"validate basic schema": {
			configPath:               "testdata/TestDataCloudletsPolicy/policy.tf",
			listPolicyVersionsReturn: []cloudlets.PolicyVersion{{Version: 1}},
			getPolicyReturn: cloudlets.Policy{
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
			},
			getPolicyVersionReturn: cloudlets.PolicyVersion{
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
						Location:                 "dd.example.com",
						UseIncomingQueryString:   true,
						UseIncomingSchemeAndHost: false,
					},
				},
				MatchRuleFormat: "1.0",
				Deleted:         false,
				Warnings:        nil,
			},
			checkFuncs: []resource.TestCheckFunc{
				resource.TestCheckResourceAttr("data.akamai_cloudlets_policy.test", "group_id", "2345"),
				resource.TestCheckResourceAttr("data.akamai_cloudlets_policy.test", "name", "SomeName"),
				resource.TestCheckResourceAttr("data.akamai_cloudlets_policy.test", "description", "Fancy Description"),
				resource.TestCheckResourceAttr("data.akamai_cloudlets_policy.test", "version_description", "Example Description"),
				resource.TestCheckResourceAttr("data.akamai_cloudlets_policy.test", "cloudlet_id", "0"),
				resource.TestCheckResourceAttr("data.akamai_cloudlets_policy.test", "cloudlet_code", "ER"),
				resource.TestCheckResourceAttr("data.akamai_cloudlets_policy.test", "api_version", "2.0"),
				resource.TestCheckResourceAttr("data.akamai_cloudlets_policy.test", "revision_id", "4824132"),
				resource.TestCheckResourceAttr("data.akamai_cloudlets_policy.test", "rules_locked", "false"),
				resource.TestCheckResourceAttr("data.akamai_cloudlets_policy.test", "match_rules", loadFixtureString("testdata/TestDataCloudletsPolicy/rules/match_rules_out.json")),
				resource.TestCheckResourceAttr("data.akamai_cloudlets_policy.test", "match_rule_format", "1.0"),
				resource.TestCheckResourceAttr("data.akamai_cloudlets_policy.test", "warnings", "null"),
			},
		},
		"validate activations schema": {
			configPath:               "testdata/TestDataCloudletsPolicy/policy.tf",
			listPolicyVersionsReturn: []cloudlets.PolicyVersion{{Version: 3}, {Version: 2}, {Version: 1}},
			getPolicyReturn: cloudlets.Policy{
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
			},
			getPolicyVersionReturn: cloudlets.PolicyVersion{
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
				Activations: []*cloudlets.Activation{
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
						Location:                 "dd.example.com",
						UseIncomingQueryString:   true,
						UseIncomingSchemeAndHost: false,
					},
				},
				MatchRuleFormat: "1.0",
				Deleted:         false,
				Warnings:        nil,
			},
			checkFuncs: []resource.TestCheckFunc{
				resource.TestCheckResourceAttr("data.akamai_cloudlets_policy.test", "activations.0.api_version", "2.0"),
				resource.TestCheckResourceAttr("data.akamai_cloudlets_policy.test", "activations.0.network", "prod"),

				resource.TestCheckResourceAttr("data.akamai_cloudlets_policy.test", "activations.0.policy_info.0.policy_id", "1234"),
				resource.TestCheckResourceAttr("data.akamai_cloudlets_policy.test", "activations.0.policy_info.0.name", "policy_name_0"),
				resource.TestCheckResourceAttr("data.akamai_cloudlets_policy.test", "activations.0.policy_info.0.version", "3"),
				resource.TestCheckResourceAttr("data.akamai_cloudlets_policy.test", "activations.0.policy_info.0.status", "active"),
				resource.TestCheckResourceAttr("data.akamai_cloudlets_policy.test", "activations.0.policy_info.0.status_detail", ""),
				resource.TestCheckResourceAttr("data.akamai_cloudlets_policy.test", "activations.0.policy_info.0.activated_by", "jsmith"),
				resource.TestCheckResourceAttr("data.akamai_cloudlets_policy.test", "activations.0.policy_info.0.activation_date", "1607507783000"),

				resource.TestCheckResourceAttr("data.akamai_cloudlets_policy.test", "activations.0.property_info.0.name", "property_name_0"),
				resource.TestCheckResourceAttr("data.akamai_cloudlets_policy.test", "activations.0.property_info.0.version", "3"),
				resource.TestCheckResourceAttr("data.akamai_cloudlets_policy.test", "activations.0.property_info.0.group_id", "132"),
				resource.TestCheckResourceAttr("data.akamai_cloudlets_policy.test", "activations.0.property_info.0.status", "active"),
				resource.TestCheckResourceAttr("data.akamai_cloudlets_policy.test", "activations.0.property_info.0.activated_by", "jsmith"),
				resource.TestCheckResourceAttr("data.akamai_cloudlets_policy.test", "activations.0.property_info.0.activation_date", "1607507783812"),

				resource.TestCheckResourceAttr("data.akamai_cloudlets_policy.test", "activations.1.api_version", "2.0"),
				resource.TestCheckResourceAttr("data.akamai_cloudlets_policy.test", "activations.1.network", "stage"),

				resource.TestCheckResourceAttr("data.akamai_cloudlets_policy.test", "activations.1.policy_info.0.policy_id", "1234"),
				resource.TestCheckResourceAttr("data.akamai_cloudlets_policy.test", "activations.1.policy_info.0.name", "policy_name_1"),
				resource.TestCheckResourceAttr("data.akamai_cloudlets_policy.test", "activations.1.policy_info.0.version", "3"),
				resource.TestCheckResourceAttr("data.akamai_cloudlets_policy.test", "activations.1.policy_info.0.status", "active"),
				resource.TestCheckResourceAttr("data.akamai_cloudlets_policy.test", "activations.1.policy_info.0.status_detail", ""),
				resource.TestCheckResourceAttr("data.akamai_cloudlets_policy.test", "activations.1.policy_info.0.activated_by", "jsmith"),
				resource.TestCheckResourceAttr("data.akamai_cloudlets_policy.test", "activations.1.policy_info.0.activation_date", "1607507783001"),

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
			checkFuncs: []resource.TestCheckFunc{
				resource.TestCheckResourceAttr("data.akamai_cloudlets_policy.test", "version", "3"),
			},
		},
		"deleted policy version": {
			configPath: "testdata/TestDataCloudletsPolicy/policy_with_version.tf",
			getPolicyVersionReturn: cloudlets.PolicyVersion{
				Version: 3,
				Deleted: true,
			},
			withError: regexp.MustCompile("specified policy version is deleted"),
		},
		"deleted policy": {
			configPath: "testdata/TestDataCloudletsPolicy/policy_with_version.tf",
			getPolicyReturn: cloudlets.Policy{
				Deleted: true,
			},
			withError: regexp.MustCompile("specified policy is deleted"),
		},
	}

	for testName, test := range tests {
		t.Run(testName, func(t *testing.T) {
			client := mockcloudlets{}
			useClient(&client, func() {
				client.On("ListPolicyVersions", mock.Anything, mock.Anything).Return(test.listPolicyVersionsReturn, nil)
				client.On("GetPolicy", mock.Anything, mock.Anything).Return(&test.getPolicyReturn, nil)
				client.On("GetPolicyVersion", mock.Anything, mock.Anything).Return(&test.getPolicyVersionReturn, nil)
				resource.UnitTest(t, resource.TestCase{
					Providers: testAccProviders,
					Steps: []resource.TestStep{
						{
							Config: loadFixtureString(test.configPath),
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
