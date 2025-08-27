package cloudlets

import (
	"fmt"
	"regexp"
	"testing"

	"github.com/akamai/AkamaiOPEN-edgegrid-golang/v12/pkg/cloudlets"
	v3 "github.com/akamai/AkamaiOPEN-edgegrid-golang/v12/pkg/cloudlets/v3"
	"github.com/akamai/terraform-provider-akamai/v9/pkg/common/ptr"
	"github.com/akamai/terraform-provider-akamai/v9/pkg/common/testutils"
	"github.com/hashicorp/terraform-plugin-testing/helper/resource"
	"github.com/hashicorp/terraform-plugin-testing/terraform"
	"github.com/jinzhu/copier"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
	"golang.org/x/exp/slices"
)

func TestResourcePolicyV2(t *testing.T) {

	type policyAttributes struct {
		name, version, matchRulesPath string
		description                   string
		timeouts                      string
	}

	var (
		expectCreatePolicy = func(client *cloudlets.Mock, policyID int64, policyName string, matchRules cloudlets.MatchRules, description string) (*cloudlets.Policy, *cloudlets.PolicyVersion) {
			policy := &cloudlets.Policy{
				PolicyID:     policyID,
				GroupID:      123,
				Name:         policyName,
				CloudletID:   0,
				CloudletCode: "ER",
			}
			version := &cloudlets.PolicyVersion{
				Location:        "/version/1",
				PolicyID:        policyID,
				Version:         1,
				Description:     description,
				MatchRules:      matchRules,
				MatchRuleFormat: "1.0",
				Warnings: []cloudlets.Warning{
					{
						Detail:      "test warning details",
						JSONPointer: "/matchRules/1/matches/0",
						Title:       "test warning",
						Type:        "test type",
					},
				},
			}
			client.On("CreatePolicy", testutils.MockContext, cloudlets.CreatePolicyRequest{
				Name:       policyName,
				CloudletID: 0,
				GroupID:    123,
			}).Return(policy, nil).Once()
			if matchRules == nil && description == "" {
				return policy, version
			}
			if matchRules != nil {
				client.On("UpdatePolicyVersion", testutils.MockContext, cloudlets.UpdatePolicyVersionRequest{
					UpdatePolicyVersion: cloudlets.UpdatePolicyVersion{
						Description: description,
						MatchRules:  matchRules,
					},
					PolicyID: policyID,
					Version:  1,
				}).Return(version, nil).Once()
			} else {
				client.On("UpdatePolicyVersion", testutils.MockContext, cloudlets.UpdatePolicyVersionRequest{
					UpdatePolicyVersion: cloudlets.UpdatePolicyVersion{
						Description: description,
						MatchRules:  make(cloudlets.MatchRules, 0),
					},
					PolicyID: policyID,
					Version:  1,
				}).Return(version, nil).Once()
			}

			return policy, version
		}

		expectListPolicyVersions = func(client *cloudlets.Mock, policyId int64, versions []cloudlets.PolicyVersion, times int) {
			client.On("ListPolicyVersions", testutils.MockContext, cloudlets.ListPolicyVersionsRequest{
				PolicyID: policyId,
				Offset:   0,
				PageSize: ptr.To(1000),
			}).Return(versions, nil).Times(times)
		}

		expectReadPolicy = func(client *cloudlets.Mock, policy *cloudlets.Policy, versions []cloudlets.PolicyVersion, times int) {
			expectListPolicyVersions(client, policy.PolicyID, versions, times)
			var latestVersion cloudlets.PolicyVersion
			for _, version := range versions {
				if latestVersion.Version < version.Version {
					latestVersion = version
				}
			}
			client.On("GetPolicy", testutils.MockContext, cloudlets.GetPolicyRequest{PolicyID: policy.PolicyID}).Return(policy, nil).Times(times)
			if len(versions) > 0 {
				var versionWithoutWarnings cloudlets.PolicyVersion
				err := copier.CopyWithOption(&versionWithoutWarnings, latestVersion, copier.Option{DeepCopy: true})
				require.NoError(t, err)
				versionWithoutWarnings.Warnings = []cloudlets.Warning{}
				versionWithoutWarnings.MatchRules = latestVersion.MatchRules
				client.On("GetPolicyVersion", testutils.MockContext, cloudlets.GetPolicyVersionRequest{
					PolicyID: policy.PolicyID,
					Version:  latestVersion.Version,
				}).Return(&versionWithoutWarnings, nil).Times(times)
			}
		}

		expectUpdatePolicy = func(client *cloudlets.Mock, policy *cloudlets.Policy, updatedName string) *cloudlets.Policy {
			var policyUpdate cloudlets.Policy
			err := copier.CopyWithOption(&policyUpdate, policy, copier.Option{DeepCopy: true})
			require.NoError(t, err)
			policyUpdate.Name = updatedName
			client.On("UpdatePolicy", testutils.MockContext, cloudlets.UpdatePolicyRequest{
				UpdatePolicy: cloudlets.UpdatePolicy{
					Name:    updatedName,
					GroupID: 123,
				},
				PolicyID: policyUpdate.PolicyID,
			}).Return(&policyUpdate, nil).Once()
			return &policyUpdate
		}

		expectCreatePolicyVersion = func(client *cloudlets.Mock, policyID int64, version *cloudlets.PolicyVersion, newMatchRules cloudlets.MatchRules) *cloudlets.PolicyVersion {
			var activatedVersion cloudlets.PolicyVersion
			err := copier.CopyWithOption(&activatedVersion, version, copier.Option{DeepCopy: true})
			require.NoError(t, err)
			activatedVersion.Activations = []cloudlets.PolicyActivation{{Network: "PROD"}}
			client.On("GetPolicyVersion", testutils.MockContext, cloudlets.GetPolicyVersionRequest{
				PolicyID:  policyID,
				Version:   version.Version,
				OmitRules: true,
			}).Return(&activatedVersion, nil).Once()
			var versionUpdate cloudlets.PolicyVersion
			err = copier.CopyWithOption(&versionUpdate, activatedVersion, copier.Option{DeepCopy: true})
			require.NoError(t, err)
			versionUpdate.MatchRules = newMatchRules
			versionUpdate.Version = version.Version + 1
			client.On("CreatePolicyVersion", testutils.MockContext, cloudlets.CreatePolicyVersionRequest{
				CreatePolicyVersion: cloudlets.CreatePolicyVersion{
					MatchRuleFormat: "1.0",
					Description:     "test policy description",
					MatchRules:      newMatchRules,
				},
				PolicyID: policyID,
			}).Return(&versionUpdate, nil).Once()
			return &versionUpdate
		}

		expectUpdatePolicyVersion = func(client *cloudlets.Mock, policyID int64, version *cloudlets.PolicyVersion, newMatchRules cloudlets.MatchRules, newDescription string) *cloudlets.PolicyVersion {
			client.On("GetPolicyVersion", testutils.MockContext, cloudlets.GetPolicyVersionRequest{
				PolicyID:  policyID,
				Version:   version.Version,
				OmitRules: true,
			}).Return(version, nil).Once()
			var versionUpdate cloudlets.PolicyVersion
			err := copier.CopyWithOption(&versionUpdate, version, copier.Option{DeepCopy: true})
			require.NoError(t, err)
			versionUpdate.MatchRules = newMatchRules
			versionUpdate.Description = newDescription
			client.On("UpdatePolicyVersion", testutils.MockContext, cloudlets.UpdatePolicyVersionRequest{
				UpdatePolicyVersion: cloudlets.UpdatePolicyVersion{
					MatchRuleFormat: "1.0",
					Description:     newDescription,
					MatchRules:      newMatchRules,
				},
				PolicyID: policyID,
				Version:  version.Version,
			}).Return(&versionUpdate, nil).Once()
			return &versionUpdate
		}

		expectRemovePolicy = func(client *cloudlets.Mock, policyID int64, numVersions, numDeleteRetries int) {
			var versionList []cloudlets.PolicyVersion
			for i := 1; i <= numVersions; i++ {
				versionList = slices.Insert(versionList, 0, cloudlets.PolicyVersion{PolicyID: policyID, Version: int64(i)})
			}
			expectListPolicyVersions(client, policyID, versionList, 1)
			for _, ver := range versionList {
				client.On("DeletePolicyVersion", testutils.MockContext, cloudlets.DeletePolicyVersionRequest{
					PolicyID: ver.PolicyID,
					Version:  ver.Version,
				}).Return(nil).Once()
			}

			pendingError := &cloudlets.Error{Detail: "Unable to delete policy because an activation for this policy is still pending"}
			if numDeleteRetries != 0 {
				client.On("RemovePolicy", testutils.MockContext, cloudlets.RemovePolicyRequest{PolicyID: policyID}).Return(pendingError).Times(numDeleteRetries)
			}
			client.On("RemovePolicy", testutils.MockContext, cloudlets.RemovePolicyRequest{PolicyID: policyID}).Return(nil).Once()
		}

		expectImportPolicy = func(client *cloudlets.Mock, policyID int64, policyName string) {
			client.On("ListPolicies", testutils.MockContext, cloudlets.ListPoliciesRequest{PageSize: ptr.To(1000), Offset: 0}).Return([]cloudlets.Policy{
				{
					PolicyID: policyID, Name: policyName,
				},
			}, nil).Once()
		}

		checkPolicyAttributes = func(attrs policyAttributes) resource.TestCheckFunc {
			var matchRulesPath string
			if attrs.matchRulesPath != "" {
				matchRulesPath = testutils.LoadFixtureString(t, attrs.matchRulesPath)
			}
			checkFunc := []resource.TestCheckFunc{
				resource.TestCheckResourceAttr("akamai_cloudlets_policy.policy", "id", "2"),
				resource.TestCheckResourceAttr("akamai_cloudlets_policy.policy", "cloudlet_code", "ER"),
				resource.TestCheckResourceAttr("akamai_cloudlets_policy.policy", "cloudlet_id", "0"),
				resource.TestCheckResourceAttr("akamai_cloudlets_policy.policy", "group_id", "123"),
				resource.TestCheckResourceAttr("akamai_cloudlets_policy.policy", "description", attrs.description),
				resource.TestCheckResourceAttr("akamai_cloudlets_policy.policy", "match_rule_format", "1.0"),
				resource.TestCheckResourceAttr("akamai_cloudlets_policy.policy", "name", attrs.name),
				resource.TestCheckResourceAttr("akamai_cloudlets_policy.policy", "version", attrs.version),
				resource.TestCheckResourceAttr("akamai_cloudlets_policy.policy", "match_rules", matchRulesPath),
				resource.TestCheckResourceAttr("akamai_cloudlets_policy.policy", "is_shared", "false"),
			}

			if attrs.timeouts != "" {
				checkFunc = append(checkFunc,
					resource.TestCheckResourceAttr("akamai_cloudlets_policy.policy", "timeouts.#", "1"),
					resource.TestCheckResourceAttr("akamai_cloudlets_policy.policy", "timeouts.0.default", attrs.timeouts),
				)
			} else {
				checkFunc = append(checkFunc,
					resource.TestCheckResourceAttr("akamai_cloudlets_policy.policy", "timeouts.#", "0"),
				)
			}

			return resource.ComposeAggregateTestCheckFunc(checkFunc...)
		}

		expectSupressDescriptionChange = func(client *cloudlets.Mock, policyID, policyVersion int64, version *cloudlets.PolicyVersion, times int) {
			client.On("GetPolicyVersion", testutils.MockContext, cloudlets.GetPolicyVersionRequest{
				PolicyID:  policyID,
				Version:   policyVersion,
				OmitRules: true,
			}).Return(version, nil).Times(times)
		}

		commonMatchRules = cloudlets.MatchRules{
			&cloudlets.MatchRuleER{
				Name:                     "r1",
				Type:                     "erMatchRule",
				UseRelativeURL:           "copy_scheme_hostname",
				StatusCode:               301,
				RedirectURL:              "/ddd",
				MatchURL:                 "abc.com",
				UseIncomingSchemeAndHost: true,
			},
			&cloudlets.MatchRuleER{
				Name: "r3",
				Type: "erMatchRule",
				Matches: []cloudlets.MatchCriteriaER{
					{
						MatchType:     "hostname",
						MatchValue:    "3333.dom",
						MatchOperator: "equals",
						CaseSensitive: true,
					},
				},
				UseRelativeURL:           "copy_scheme_hostname",
				StatusCode:               307,
				RedirectURL:              "/abc/sss",
				UseIncomingSchemeAndHost: true,
			},
		}
	)

	t.Run("policy lifecycle with create new version", func(t *testing.T) {
		testDir := "testdata/TestResPolicy/lifecycle"

		client := new(cloudlets.Mock)
		policy, version := expectCreatePolicy(client, 2, "test_policy", commonMatchRules, "test policy description")
		policyVersions := []cloudlets.PolicyVersion{*version}
		expectReadPolicy(client, policy, policyVersions, 3)
		policy = expectUpdatePolicy(client, policy, "test_policy_updated")
		version = expectCreatePolicyVersion(client, policy.PolicyID, version, commonMatchRules[:1])
		policyVersions = slices.Insert(policyVersions, 0, *version)
		expectReadPolicy(client, policy, policyVersions, 2)
		expectRemovePolicy(client, policy.PolicyID, 2, 0)

		useClient(client, func() {
			resource.UnitTest(t, resource.TestCase{
				ProtoV6ProviderFactories: testutils.NewProtoV6ProviderFactory(NewSubprovider()),
				Steps: []resource.TestStep{
					{
						Config: testutils.LoadFixtureStringf(t, "%s/policy_create.tf", testDir),
						Check: checkPolicyAttributes(policyAttributes{
							name:           "test_policy",
							version:        "1",
							matchRulesPath: fmt.Sprintf("%s/match_rules/match_rules_create.json", testDir),
							description:    "test policy description",
						}),
					},
					{
						Config: testutils.LoadFixtureStringf(t, "%s/policy_update.tf", testDir),
						Check: checkPolicyAttributes(policyAttributes{
							name:           "test_policy_updated",
							version:        "2",
							matchRulesPath: fmt.Sprintf("%s/match_rules/match_rules_update.json", testDir),
							description:    "test policy description",
						}),
					},
				},
			})
		})
		client.AssertExpectations(t)
	})

	t.Run("policy lifecycle detects new version drift", func(t *testing.T) {
		testDir := "testdata/TestResPolicy/lifecycle_with_drift"

		client := new(cloudlets.Mock)
		policy, version := expectCreatePolicy(client, 2, "test_policy", commonMatchRules, "test policy description")
		policyVersions := []cloudlets.PolicyVersion{*version}
		expectReadPolicy(client, policy, policyVersions, 1)
		// CustomDiff call to check if the policy version is active
		expectSupressDescriptionChange(client, 2, 2, version, 1)
		// new version which causes drift
		versionWithDrift := &cloudlets.PolicyVersion{
			Location:        "/version/2",
			PolicyID:        2,
			Version:         2,
			Description:     "new description after drift",
			MatchRules:      commonMatchRules,
			MatchRuleFormat: "1.0",
			Warnings: []cloudlets.Warning{
				{
					Detail:      "test warning details",
					JSONPointer: "/matchRules/1/matches/0",
					Title:       "test warning",
					Type:        "test type",
				},
			},
		}
		policyVersions = slices.Insert(policyVersions, 0, *versionWithDrift)
		expectReadPolicy(client, policy, policyVersions, 1)
		expectRemovePolicy(client, policy.PolicyID, 2, 0)

		useClient(client, func() {
			resource.UnitTest(t, resource.TestCase{
				ProtoV6ProviderFactories: testutils.NewProtoV6ProviderFactory(NewSubprovider()),
				Steps: []resource.TestStep{
					{
						Config: testutils.LoadFixtureStringf(t, "%s/policy_create.tf", testDir),
						Check: checkPolicyAttributes(policyAttributes{
							name:           "test_policy",
							version:        "1",
							matchRulesPath: fmt.Sprintf("%s/match_rules/match_rules_create.json", testDir),
							description:    "test policy description",
						}),
						ExpectNonEmptyPlan: true,
					},
				},
			})
		})
		client.AssertExpectations(t)
	})

	t.Run("policy lifecycle and delete retries", func(t *testing.T) {
		testDir := "testdata/TestResPolicy/lifecycle"

		client := new(cloudlets.Mock)
		policy, version := expectCreatePolicy(client, 2, "test_policy", commonMatchRules, "test policy description")
		expectReadPolicy(client, policy, []cloudlets.PolicyVersion{*version}, 2)
		expectRemovePolicy(client, policy.PolicyID, 1, 1)

		useClient(client, func() {
			resource.UnitTest(t, resource.TestCase{
				ProtoV6ProviderFactories: testutils.NewProtoV6ProviderFactory(NewSubprovider()),
				Steps: []resource.TestStep{
					{
						Config: testutils.LoadFixtureStringf(t, "%s/policy_create.tf", testDir),
						Check: checkPolicyAttributes(policyAttributes{
							name:           "test_policy",
							version:        "1",
							matchRulesPath: fmt.Sprintf("%s/match_rules/match_rules_create.json", testDir),
							description:    "test policy description",
						}),
					},
				},
			})
		})
		client.AssertExpectations(t)
	})

	t.Run("policy lifecycle with update existing version", func(t *testing.T) {
		testDir := "testdata/TestResPolicy/lifecycle"

		client := new(cloudlets.Mock)
		policy, version := expectCreatePolicy(client, 2, "test_policy", commonMatchRules, "test policy description")
		policyVersions := []cloudlets.PolicyVersion{*version}
		expectReadPolicy(client, policy, policyVersions, 3)
		policy = expectUpdatePolicy(client, policy, "test_policy_updated")
		version = expectUpdatePolicyVersion(client, policy.PolicyID, version, commonMatchRules[:1], "test policy description")
		expectReadPolicy(client, policy, []cloudlets.PolicyVersion{*version}, 2)
		expectRemovePolicy(client, policy.PolicyID, 1, 0)

		useClient(client, func() {
			resource.UnitTest(t, resource.TestCase{
				ProtoV6ProviderFactories: testutils.NewProtoV6ProviderFactory(NewSubprovider()),
				Steps: []resource.TestStep{
					{
						Config: testutils.LoadFixtureStringf(t, "%s/policy_create.tf", testDir),
						Check: checkPolicyAttributes(policyAttributes{
							name:           "test_policy",
							version:        "1",
							matchRulesPath: fmt.Sprintf("%s/match_rules/match_rules_create.json", testDir),
							description:    "test policy description",
						}),
					},
					{
						Config: testutils.LoadFixtureStringf(t, "%s/policy_update.tf", testDir),
						Check: checkPolicyAttributes(policyAttributes{
							name:           "test_policy_updated",
							version:        "1",
							matchRulesPath: fmt.Sprintf("%s/match_rules/match_rules_update.json", testDir),
							description:    "test policy description",
						}),
					},
				},
			})
		})
		client.AssertExpectations(t)
	})

	t.Run("update only policy", func(t *testing.T) {
		testDir := "testdata/TestResPolicy/lifecycle_policy_update"

		client := new(cloudlets.Mock)
		policy, version := expectCreatePolicy(client, 2, "test_policy", commonMatchRules, "test policy description")
		policyVersions := []cloudlets.PolicyVersion{*version}
		expectReadPolicy(client, policy, policyVersions, 3)
		policy = expectUpdatePolicy(client, policy, "test_policy_updated")
		policyVersions = slices.Insert(policyVersions, 0, *version)
		expectReadPolicy(client, policy, policyVersions, 2)
		expectRemovePolicy(client, policy.PolicyID, 1, 0)

		useClient(client, func() {
			resource.UnitTest(t, resource.TestCase{
				ProtoV6ProviderFactories: testutils.NewProtoV6ProviderFactory(NewSubprovider()),
				Steps: []resource.TestStep{
					{
						Config: testutils.LoadFixtureStringf(t, "%s/policy_create.tf", testDir),
						Check: checkPolicyAttributes(policyAttributes{
							name:           "test_policy",
							version:        "1",
							matchRulesPath: fmt.Sprintf("%s/match_rules/match_rules.json", testDir),
							description:    "test policy description",
						}),
					},
					{
						Config: testutils.LoadFixtureStringf(t, "%s/policy_update.tf", testDir),
						Check: checkPolicyAttributes(policyAttributes{
							name:           "test_policy_updated",
							version:        "1",
							matchRulesPath: fmt.Sprintf("%s/match_rules/match_rules.json", testDir),
							description:    "test policy description",
						}),
					},
				},
			})
		})
		client.AssertExpectations(t)
	})

	t.Run("update only version", func(t *testing.T) {
		testDir := "testdata/TestResPolicy/lifecycle_version_update"

		client := new(cloudlets.Mock)
		policy, version := expectCreatePolicy(client, 2, "test_policy", commonMatchRules, "test policy description")
		policyVersions := []cloudlets.PolicyVersion{*version}
		expectReadPolicy(client, policy, policyVersions, 3)
		version = expectUpdatePolicyVersion(client, policy.PolicyID, version, commonMatchRules[:1], "test policy description")
		policyVersions = []cloudlets.PolicyVersion{*version}
		expectReadPolicy(client, policy, policyVersions, 2)
		expectRemovePolicy(client, policy.PolicyID, 1, 0)

		useClient(client, func() {
			resource.UnitTest(t, resource.TestCase{
				ProtoV6ProviderFactories: testutils.NewProtoV6ProviderFactory(NewSubprovider()),
				Steps: []resource.TestStep{
					{
						Config: testutils.LoadFixtureStringf(t, "%s/policy_create.tf", testDir),
						Check: checkPolicyAttributes(policyAttributes{
							name:           "test_policy",
							version:        "1",
							matchRulesPath: fmt.Sprintf("%s/match_rules/match_rules_create.json", testDir),
							description:    "test policy description",
						}),
					},
					{
						Config: testutils.LoadFixtureStringf(t, "%s/policy_update.tf", testDir),
						Check: checkPolicyAttributes(policyAttributes{
							name:           "test_policy",
							version:        "1",
							matchRulesPath: fmt.Sprintf("%s/match_rules/match_rules_update.json", testDir),
							description:    "test policy description",
						}),
					},
				},
			})
		})
		client.AssertExpectations(t)
	})

	t.Run("update only description for not active policy - expect no new version", func(t *testing.T) {
		testDir := "testdata/TestResPolicy/lifecycle_description_update"

		client := new(cloudlets.Mock)
		policy, version := expectCreatePolicy(client, 2, "test_policy", commonMatchRules, "test policy description")
		policyVersions := []cloudlets.PolicyVersion{*version}
		expectReadPolicy(client, policy, policyVersions, 3)
		// CustomDiff calls to check if the policy version is active
		expectSupressDescriptionChange(client, 2, 1, version, 2)
		version = expectUpdatePolicyVersion(client, policy.PolicyID, version, commonMatchRules, "test policy description - updated")
		policyVersions = []cloudlets.PolicyVersion{*version}
		expectReadPolicy(client, policy, policyVersions, 2)
		expectRemovePolicy(client, policy.PolicyID, 1, 0)

		useClient(client, func() {
			resource.UnitTest(t, resource.TestCase{
				ProtoV6ProviderFactories: testutils.NewProtoV6ProviderFactory(NewSubprovider()),
				Steps: []resource.TestStep{
					{
						Config: testutils.LoadFixtureStringf(t, "%s/policy_create.tf", testDir),
						Check: checkPolicyAttributes(policyAttributes{
							name:           "test_policy",
							version:        "1",
							matchRulesPath: fmt.Sprintf("%s/match_rules/match_rules.json", testDir),
							description:    "test policy description",
						}),
					},
					{
						Config: testutils.LoadFixtureStringf(t, "%s/policy_update.tf", testDir),
						Check: checkPolicyAttributes(policyAttributes{
							name:           "test_policy",
							version:        "1",
							matchRulesPath: fmt.Sprintf("%s/match_rules/match_rules.json", testDir),
							description:    "test policy description - updated",
						}),
					},
				},
			})
		})
		client.AssertExpectations(t)
	})

	t.Run("update description for active policy version - suppress diff", func(t *testing.T) {
		testDir := "testdata/TestResPolicy/lifecycle_description_update"

		client := new(cloudlets.Mock)
		policy, version := expectCreatePolicy(client, 2, "test_policy", commonMatchRules, "test policy description")
		policyVersions := []cloudlets.PolicyVersion{*version}
		expectReadPolicy(client, policy, policyVersions, 3)
		version.Activations = []cloudlets.PolicyActivation{{Network: "STAGING"}}
		// CustomDiff calls to check if the policy version is active
		expectSupressDescriptionChange(client, 2, 1, version, 3)
		expectReadPolicy(client, policy, policyVersions, 1)
		expectRemovePolicy(client, policy.PolicyID, 1, 0)

		useClient(client, func() {
			resource.UnitTest(t, resource.TestCase{
				ProtoV6ProviderFactories: testutils.NewProtoV6ProviderFactory(NewSubprovider()),
				Steps: []resource.TestStep{
					{
						Config: testutils.LoadFixtureStringf(t, "%s/policy_create.tf", testDir),
						Check: checkPolicyAttributes(policyAttributes{
							name:           "test_policy",
							version:        "1",
							matchRulesPath: fmt.Sprintf("%s/match_rules/match_rules.json", testDir),
							description:    "test policy description",
						}),
					},
					{
						Config: testutils.LoadFixtureStringf(t, "%s/policy_update.tf", testDir),
						Check: checkPolicyAttributes(policyAttributes{
							name:           "test_policy",
							version:        "1",
							matchRulesPath: fmt.Sprintf("%s/match_rules/match_rules.json", testDir),
							description:    "test policy description",
						}),
					},
				},
			})
		})
		client.AssertExpectations(t)
	})

	t.Run("warnings creating and updating version", func(t *testing.T) {
		testDir := "testdata/TestResPolicy/lifecycle_version_update"

		client := new(cloudlets.Mock)
		policy, version := expectCreatePolicy(client, 2, "test_policy", commonMatchRules, "test policy description")
		policyVersions := []cloudlets.PolicyVersion{*version}
		expectReadPolicy(client, policy, policyVersions, 3)
		version = expectUpdatePolicyVersion(client, policy.PolicyID, version, commonMatchRules[:1], "test policy description")
		// update existing version in slice by deleting old policyVersions and defining new one
		policyVersions = []cloudlets.PolicyVersion{*version}
		expectReadPolicy(client, policy, policyVersions, 4)
		expectRemovePolicy(client, policy.PolicyID, 1, 0)

		warningsJSON, err := warningsToJSON(version.Warnings)
		require.NoError(t, err)

		checkWarnings := resource.ComposeAggregateTestCheckFunc(
			resource.TestCheckResourceAttr("akamai_cloudlets_policy.policy", "warnings", string(warningsJSON)),
			resource.TestMatchOutput("policy_output", regexp.MustCompile("test warning")),
		)

		useClient(client, func() {
			resource.UnitTest(t, resource.TestCase{
				ProtoV6ProviderFactories: testutils.NewProtoV6ProviderFactory(NewSubprovider()),
				Steps: []resource.TestStep{
					{
						Config: testutils.LoadFixtureStringf(t, "%s/policy_create.tf", testDir),
						Check:  checkWarnings,
					},
					{
						Config: testutils.LoadFixtureStringf(t, "%s/policy_update.tf", testDir),
						Check:  checkWarnings,
					},
					{
						Config: testutils.LoadFixtureStringf(t, "%s/policy_update.tf", testDir),
						Check:  checkWarnings,
					},
				},
			})
		})
		client.AssertExpectations(t)
	})

	t.Run("remove match rules from version", func(t *testing.T) {
		testDir := "testdata/TestResPolicy/lifecycle_remove_match_rules"

		client := new(cloudlets.Mock)
		policy, version := expectCreatePolicy(client, 2, "test_policy", commonMatchRules, "test policy description")
		policyVersions := []cloudlets.PolicyVersion{*version}
		expectReadPolicy(client, policy, policyVersions, 3)
		version = expectUpdatePolicyVersion(client, policy.PolicyID, version, cloudlets.MatchRules{}, "test policy description")
		// update existing version in slice by deleting old policyVersions and defining new one
		policyVersions = []cloudlets.PolicyVersion{*version}
		expectReadPolicy(client, policy, policyVersions, 2)
		expectRemovePolicy(client, policy.PolicyID, 1, 0)

		useClient(client, func() {
			resource.UnitTest(t, resource.TestCase{
				ProtoV6ProviderFactories: testutils.NewProtoV6ProviderFactory(NewSubprovider()),
				Steps: []resource.TestStep{
					{
						Config: testutils.LoadFixtureStringf(t, "%s/policy_create.tf", testDir),
						Check: checkPolicyAttributes(policyAttributes{
							name:           "test_policy",
							version:        "1",
							matchRulesPath: fmt.Sprintf("%s/match_rules/match_rules_create.json", testDir),
							description:    "test policy description",
						}),
					},
					{
						Config: testutils.LoadFixtureStringf(t, "%s/policy_update.tf", testDir),
						Check: checkPolicyAttributes(policyAttributes{
							name:           "test_policy",
							version:        "1",
							matchRulesPath: "",
							description:    "test policy description",
						}),
					},
				},
			})
		})
		client.AssertExpectations(t)
	})

	t.Run("create policy without match rules or description", func(t *testing.T) {
		testDir := "testdata/TestResPolicy/create_no_match_rules_no_description"

		client := new(cloudlets.Mock)
		policy, version := expectCreatePolicy(client, 2, "test_policy", nil, "")
		policyVersions := []cloudlets.PolicyVersion{*version}
		expectReadPolicy(client, policy, policyVersions, 2)
		expectRemovePolicy(client, 2, 1, 0)
		useClient(client, func() {
			resource.UnitTest(t, resource.TestCase{
				ProtoV6ProviderFactories: testutils.NewProtoV6ProviderFactory(NewSubprovider()),
				Steps: []resource.TestStep{
					{
						Config: testutils.LoadFixtureStringf(t, "%s/policy_create.tf", testDir),
						Check: checkPolicyAttributes(policyAttributes{
							name:           "test_policy",
							version:        "1",
							matchRulesPath: "",
							description:    "",
						}),
					},
				},
			})
		})
		client.AssertExpectations(t)
	})

	t.Run("create policy without match rules with description", func(t *testing.T) {
		testDir := "testdata/TestResPolicy/create_no_match_rules"

		client := new(cloudlets.Mock)
		policy, version := expectCreatePolicy(client, 2, "test_policy", nil, "test policy description")
		policyVersions := []cloudlets.PolicyVersion{*version}
		expectReadPolicy(client, policy, policyVersions, 2)
		expectRemovePolicy(client, 2, 1, 0)
		useClient(client, func() {
			resource.UnitTest(t, resource.TestCase{
				ProtoV6ProviderFactories: testutils.NewProtoV6ProviderFactory(NewSubprovider()),
				Steps: []resource.TestStep{
					{
						Config: testutils.LoadFixtureStringf(t, "%s/policy_create.tf", testDir),
						Check: checkPolicyAttributes(policyAttributes{
							name:           "test_policy",
							version:        "1",
							matchRulesPath: "",
							description:    "test policy description",
						}),
					},
				},
			})
		})
		client.AssertExpectations(t)
	})

	t.Run("create policy with timeout", func(t *testing.T) {
		testDir := "testdata/TestResPolicy/timeouts"

		matchRules := cloudlets.MatchRules{
			&cloudlets.MatchRuleER{
				Name:                     "r1",
				Type:                     "erMatchRule",
				UseRelativeURL:           "copy_scheme_hostname",
				StatusCode:               301,
				RedirectURL:              "/ddd",
				MatchURL:                 "abc.com",
				UseIncomingSchemeAndHost: true,
			},
		}

		client := new(cloudlets.Mock)
		policy, version := expectCreatePolicy(client, 2, "test_policy", matchRules, "test policy description")
		policyVersions := []cloudlets.PolicyVersion{*version}
		expectReadPolicy(client, policy, policyVersions, 2)
		expectRemovePolicy(client, 2, 1, 0)
		useClient(client, func() {
			resource.UnitTest(t, resource.TestCase{
				ProtoV6ProviderFactories: testutils.NewProtoV6ProviderFactory(NewSubprovider()),
				Steps: []resource.TestStep{
					{
						Config: testutils.LoadFixtureStringf(t, "%s/policy_create.tf", testDir),
						Check: checkPolicyAttributes(policyAttributes{
							name:           "test_policy",
							version:        "1",
							matchRulesPath: fmt.Sprintf("%s/match_rules/match_rules.json", testDir),
							timeouts:       "4h",
							description:    "test policy description",
						}),
					},
				},
			})
		})
		client.AssertExpectations(t)
	})

	t.Run("error creating policy", func(t *testing.T) {
		testDir := "testdata/TestResPolicy/lifecycle"

		client := new(cloudlets.Mock)
		client.On("CreatePolicy", testutils.MockContext, cloudlets.CreatePolicyRequest{
			Name:       "test_policy",
			CloudletID: 0,
			GroupID:    123,
		}).Return(nil, fmt.Errorf("oops"))

		useClient(client, func() {
			resource.UnitTest(t, resource.TestCase{
				ProtoV6ProviderFactories: testutils.NewProtoV6ProviderFactory(NewSubprovider()),
				Steps: []resource.TestStep{
					{
						Config:      testutils.LoadFixtureStringf(t, "%s/policy_create.tf", testDir),
						ExpectError: regexp.MustCompile("oops"),
					},
				},
			})
		})
		client.AssertExpectations(t)
	})

	t.Run("error creating policy version", func(t *testing.T) {
		testDir := "testdata/TestResPolicy/lifecycle"

		policy := &cloudlets.Policy{
			PolicyID:     2,
			GroupID:      123,
			Name:         "test_policy",
			CloudletID:   0,
			CloudletCode: "ER",
		}

		expectErrorCreatingVersion := func(client *cloudlets.Mock) {
			client.On("CreatePolicy", testutils.MockContext, cloudlets.CreatePolicyRequest{
				Name:       "test_policy",
				CloudletID: 0,
				GroupID:    123,
			}).Return(policy, nil)
			client.On("UpdatePolicyVersion", testutils.MockContext, cloudlets.UpdatePolicyVersionRequest{
				UpdatePolicyVersion: cloudlets.UpdatePolicyVersion{
					Description: "test policy description",
					MatchRules:  commonMatchRules,
				},
				PolicyID: 2,
				Version:  1,
			}).Return(nil, fmt.Errorf("UpdatePolicyVersionError"))
		}

		testCases := []struct {
			Expectations  func(client *cloudlets.Mock)
			ExpectedError *regexp.Regexp
		}{
			{
				Expectations: func(client *cloudlets.Mock) {
					expectErrorCreatingVersion(client)
					expectReadPolicy(client, policy, []cloudlets.PolicyVersion{{
						PolicyID: 2,
						Version:  1,
					}}, 1)
					expectRemovePolicy(client, 2, 1, 0)

				},
				ExpectedError: regexp.MustCompile("UpdatePolicyVersionError"),
			},
			{
				Expectations: func(client *cloudlets.Mock) {
					expectErrorCreatingVersion(client)
					expectListPolicyVersions(client, policy.PolicyID, []cloudlets.PolicyVersion{{
						PolicyID: 2,
						Version:  1,
					}}, 1)
					client.On("GetPolicy", testutils.MockContext, cloudlets.GetPolicyRequest{PolicyID: policy.PolicyID}).Return(nil, fmt.Errorf("GetPolicyError"))
					expectRemovePolicy(client, 2, 1, 0)
				},
				ExpectedError: regexp.MustCompile("(?s)GetPolicyError.*UpdatePolicyVersionError"),
			},
		}

		for i := range testCases {
			client := new(cloudlets.Mock)
			testCases[i].Expectations(client)
			useClient(client, func() {
				resource.UnitTest(t, resource.TestCase{
					ProtoV6ProviderFactories: testutils.NewProtoV6ProviderFactory(NewSubprovider()),
					Steps: []resource.TestStep{
						{
							Config:      testutils.LoadFixtureStringf(t, "%s/policy_create.tf", testDir),
							ExpectError: testCases[i].ExpectedError,
						},
					},
				})
			})
			client.AssertExpectations(t)
		}
	})

	t.Run("error fetching policy", func(t *testing.T) {
		testDir := "testdata/TestResPolicy/create_no_match_rules"

		client := new(cloudlets.Mock)
		policy, version := expectCreatePolicy(client, 2, "test_policy", nil, "test policy description")
		policyVersions := []cloudlets.PolicyVersion{*version}
		expectListPolicyVersions(client, policy.PolicyID, policyVersions, 1)
		client.On("GetPolicy", testutils.MockContext, cloudlets.GetPolicyRequest{PolicyID: policy.PolicyID}).Return(nil, fmt.Errorf("oops"))
		expectRemovePolicy(client, 2, 1, 0)
		useClient(client, func() {
			resource.UnitTest(t, resource.TestCase{
				ProtoV6ProviderFactories: testutils.NewProtoV6ProviderFactory(NewSubprovider()),
				Steps: []resource.TestStep{
					{
						Config:      testutils.LoadFixtureStringf(t, "%s/policy_create.tf", testDir),
						ExpectError: regexp.MustCompile("oops"),
					},
				},
			})
		})
		client.AssertExpectations(t)
	})

	t.Run("error fetching policy version", func(t *testing.T) {
		testDir := "testdata/TestResPolicy/create_no_match_rules"

		client := new(cloudlets.Mock)
		policy, version := expectCreatePolicy(client, 2, "test_policy", nil, "test policy description")
		policyVersions := []cloudlets.PolicyVersion{*version}
		expectListPolicyVersions(client, policy.PolicyID, policyVersions, 1)
		client.On("GetPolicy", testutils.MockContext, cloudlets.GetPolicyRequest{PolicyID: policy.PolicyID}).Return(policy, nil)
		client.On("GetPolicyVersion", testutils.MockContext, cloudlets.GetPolicyVersionRequest{
			PolicyID: policy.PolicyID,
			Version:  version.Version,
		}).Return(nil, fmt.Errorf("oops"))
		expectRemovePolicy(client, 2, 1, 0)
		useClient(client, func() {
			resource.UnitTest(t, resource.TestCase{
				ProtoV6ProviderFactories: testutils.NewProtoV6ProviderFactory(NewSubprovider()),
				Steps: []resource.TestStep{
					{
						Config:      testutils.LoadFixtureStringf(t, "%s/policy_create.tf", testDir),
						ExpectError: regexp.MustCompile("oops"),
					},
				},
			})
		})
		client.AssertExpectations(t)
	})

	t.Run("error updating policy", func(t *testing.T) {
		testDir := "testdata/TestResPolicy/lifecycle_policy_update"

		client := new(cloudlets.Mock)
		policy, version := expectCreatePolicy(client, 2, "test_policy", commonMatchRules, "test policy description")
		policyVersions := []cloudlets.PolicyVersion{*version}
		expectReadPolicy(client, policy, policyVersions, 3)
		client.On("UpdatePolicy", testutils.MockContext, cloudlets.UpdatePolicyRequest{
			UpdatePolicy: cloudlets.UpdatePolicy{
				Name:    "test_policy_updated",
				GroupID: 123,
			},
			PolicyID: policy.PolicyID,
		}).Return(nil, fmt.Errorf("oops")).Once()
		expectRemovePolicy(client, policy.PolicyID, 1, 0)

		useClient(client, func() {
			resource.UnitTest(t, resource.TestCase{
				ProtoV6ProviderFactories: testutils.NewProtoV6ProviderFactory(NewSubprovider()),
				Steps: []resource.TestStep{
					{
						Config: testutils.LoadFixtureStringf(t, "%s/policy_create.tf", testDir),
					},
					{
						Config:      testutils.LoadFixtureStringf(t, "%s/policy_update.tf", testDir),
						ExpectError: regexp.MustCompile("oops"),
					},
				},
			})
		})
		client.AssertExpectations(t)
	})

	t.Run("error updating version", func(t *testing.T) {
		testDir := "testdata/TestResPolicy/lifecycle_version_update"

		expectErrorUpdatingVersion := func(client *cloudlets.Mock, expectReadPolicyTimes int) (policy *cloudlets.Policy) {
			policy, version := expectCreatePolicy(client, 2, "test_policy", commonMatchRules, "test policy description")
			expectReadPolicy(client, policy, []cloudlets.PolicyVersion{*version}, expectReadPolicyTimes)
			client.On("GetPolicyVersion", testutils.MockContext, cloudlets.GetPolicyVersionRequest{
				PolicyID:  policy.PolicyID,
				Version:   version.Version,
				OmitRules: true,
			}).Return(version, nil).Once()
			client.On("UpdatePolicyVersion", testutils.MockContext, cloudlets.UpdatePolicyVersionRequest{
				UpdatePolicyVersion: cloudlets.UpdatePolicyVersion{
					Description:     "test policy description",
					MatchRuleFormat: "1.0",
					MatchRules:      commonMatchRules[:1],
				},
				PolicyID: policy.PolicyID,
				Version:  version.Version,
			}).Return(nil, fmt.Errorf("UpdatePolicyVersionError")).Once()
			expectRemovePolicy(client, policy.PolicyID, 1, 0)
			return
		}

		testCases := []struct {
			Expectations  func(client *cloudlets.Mock)
			ExpectedError *regexp.Regexp
		}{
			{
				Expectations: func(client *cloudlets.Mock) {
					expectErrorUpdatingVersion(client, 4)
				},
				ExpectedError: regexp.MustCompile("UpdatePolicyVersionError"),
			},
			{
				Expectations: func(client *cloudlets.Mock) {
					policy := expectErrorUpdatingVersion(client, 3)
					client.On("GetPolicy", testutils.MockContext, cloudlets.GetPolicyRequest{PolicyID: policy.PolicyID}).Return(nil, fmt.Errorf("GetPolicyError"))
					client.On("ListPolicyVersions", testutils.MockContext, cloudlets.ListPolicyVersionsRequest{PolicyID: policy.PolicyID,
						Offset:   0,
						PageSize: ptr.To(1000),
					}).Return([]cloudlets.PolicyVersion{{
						PolicyID: 2,
						Version:  1,
					},
					}, nil)
				},
				ExpectedError: regexp.MustCompile("(?s)GetPolicyError.*UpdatePolicyVersionError"),
			},
		}

		for i := range testCases {
			client := new(cloudlets.Mock)
			testCases[i].Expectations(client)
			useClient(client, func() {
				resource.UnitTest(t, resource.TestCase{
					ProtoV6ProviderFactories: testutils.NewProtoV6ProviderFactory(NewSubprovider()),
					Steps: []resource.TestStep{
						{
							Config: testutils.LoadFixtureStringf(t, "%s/policy_create.tf", testDir),
						},
						{
							Config:      testutils.LoadFixtureStringf(t, "%s/policy_update.tf", testDir),
							ExpectError: testCases[i].ExpectedError,
						},
					},
				})
			})
			client.AssertExpectations(t)
		}
	})

	t.Run("invalid group id passed", func(t *testing.T) {
		testDir := "testdata/TestResPolicy/invalid_group_id"
		client := new(cloudlets.Mock)

		useClient(client, func() {
			resource.UnitTest(t, resource.TestCase{
				ProtoV6ProviderFactories: testutils.NewProtoV6ProviderFactory(NewSubprovider()),
				Steps: []resource.TestStep{
					{
						Config:      testutils.LoadFixtureStringf(t, "%s/policy_create.tf", testDir),
						ExpectError: regexp.MustCompile("invalid group_id provided"),
					},
				},
			})
		})
		client.AssertExpectations(t)
	})

	t.Run("import policy", func(t *testing.T) {
		testDir := "testdata/TestResPolicy/import"
		client := new(cloudlets.Mock)

		policy, version := expectCreatePolicy(client, 2, "test_policy", commonMatchRules, "test policy description")
		policyVersions := []cloudlets.PolicyVersion{*version}
		expectReadPolicy(client, policy, policyVersions, 3)
		expectImportPolicy(client, 2, "test_policy")
		expectRemovePolicy(client, 2, 1, 0)

		useClient(client, func() {
			resource.UnitTest(t, resource.TestCase{
				ProtoV6ProviderFactories: testutils.NewProtoV6ProviderFactory(NewSubprovider()),
				Steps: []resource.TestStep{
					{
						Config: testutils.LoadFixtureStringf(t, "%s/policy_create.tf", testDir),
					},
					{
						ImportState:             true,
						ImportStateId:           "test_policy",
						ResourceName:            "akamai_cloudlets_policy.policy",
						ImportStateVerify:       true,
						ImportStateVerifyIgnore: []string{"warnings"},
					},
				},
			})
		})
		client.AssertExpectations(t)
	})

	t.Run("import policy - test checkForV2Policy()", func(t *testing.T) {
		testDir := "testdata/TestResPolicy/import"
		client := new(cloudlets.Mock)
		policyID := int64(2)

		policy, version := expectCreatePolicy(client, policyID, "test_policy", commonMatchRules, "test policy description")
		policyVersions := []cloudlets.PolicyVersion{*version}
		expectReadPolicy(client, policy, policyVersions, 3)
		// custom import mocks
		// mock that 1000 policies are returned, desired not found
		client.On("ListPolicies", testutils.MockContext, cloudlets.ListPoliciesRequest{
			PageSize: ptr.To(1000),
			Offset:   0,
		}).Return(make([]cloudlets.Policy, 1000), nil).Once()
		// mock that desired policy is on the next page
		client.On("ListPolicies", testutils.MockContext, cloudlets.ListPoliciesRequest{
			PageSize: ptr.To(1000),
			Offset:   1,
		}).Return([]cloudlets.Policy{
			{
				PolicyID: policyID,
				Name:     "test_policy",
			},
		}, nil).Once()
		expectRemovePolicy(client, 2, 1, 0)

		useClient(client, func() {
			resource.UnitTest(t, resource.TestCase{
				ProtoV6ProviderFactories: testutils.NewProtoV6ProviderFactory(NewSubprovider()),
				Steps: []resource.TestStep{
					{
						Config: testutils.LoadFixtureStringf(t, "%s/policy_create.tf", testDir),
					},
					{
						ImportState:             true,
						ImportStateId:           "test_policy",
						ResourceName:            "akamai_cloudlets_policy.policy",
						ImportStateVerify:       true,
						ImportStateVerifyIgnore: []string{"warnings"},
					},
				},
			})
		})
		client.AssertExpectations(t)
	})

	t.Run("error importing policy not found", func(t *testing.T) {
		testDir := "testdata/TestResPolicy/import_no_match_rules"
		clientV2 := new(cloudlets.Mock)
		clientV3 := new(v3.Mock)

		policy, version := expectCreatePolicy(clientV2, 2, "test_policy", nil, "test policy description")
		policyVersions := []cloudlets.PolicyVersion{*version}
		expectReadPolicy(clientV2, policy, policyVersions, 2)
		clientV2.On("ListPolicies", testutils.MockContext, cloudlets.ListPoliciesRequest{PageSize: ptr.To(1000), Offset: 0}).
			Return([]cloudlets.Policy{}, nil).Once()
		clientV3.On("ListPolicies", testutils.MockContext, v3.ListPoliciesRequest{
			Page: 0,
			Size: 1000,
		}).Return(&v3.ListPoliciesResponse{
			Content: []v3.Policy{},
		}, nil).Once()
		expectRemovePolicy(clientV2, 2, 1, 0)

		useClientV2AndV3(clientV2, clientV3, func() {
			resource.UnitTest(t, resource.TestCase{
				ProtoV6ProviderFactories: testutils.NewProtoV6ProviderFactory(NewSubprovider()),
				Steps: []resource.TestStep{
					{
						Config: testutils.LoadFixtureStringf(t, "%s/policy_create.tf", testDir),
					},
					{
						ImportState:   true,
						ImportStateId: "not_existing_test_policy",
						ResourceName:  "akamai_cloudlets_policy.policy",
						ExpectError:   regexp.MustCompile("policy 'not_existing_test_policy' does not exist"),
					},
				},
			})
		})
		clientV2.AssertExpectations(t)
		clientV3.AssertExpectations(t)
	})

	t.Run("importing policy no version found", func(t *testing.T) {
		testDir := "testdata/TestResPolicy/import_no_version"
		client := new(cloudlets.Mock)

		policy, _ := expectCreatePolicy(client, 2, "test_policy", nil, "")
		expectReadPolicy(client, policy, nil, 2)
		expectImportPolicy(client, 2, "test_policy")
		//read after import
		expectReadPolicy(client, policy, nil, 1)
		expectRemovePolicy(client, 2, 1, 0)

		useClient(client, func() {
			resource.UnitTest(t, resource.TestCase{
				ProtoV6ProviderFactories: testutils.NewProtoV6ProviderFactory(NewSubprovider()),
				Steps: []resource.TestStep{
					{
						Config: testutils.LoadFixtureStringf(t, "%s/policy_create.tf", testDir),
					},
					{
						ImportState:   true,
						ImportStateId: "test_policy",
						ResourceName:  "akamai_cloudlets_policy.policy",
						Check: checkPolicyAttributes(policyAttributes{
							name:           "test_policy",
							version:        "",
							matchRulesPath: "",
							description:    "",
						}),
					},
				},
			})
		})
		client.AssertExpectations(t)
	})

	t.Run("error importing policy name cannot be empty", func(t *testing.T) {
		testDir := "testdata/TestResPolicy/import_no_match_rules"
		client := new(cloudlets.Mock)

		policy, version := expectCreatePolicy(client, 2, "test_policy", nil, "test policy description")
		policyVersions := []cloudlets.PolicyVersion{*version}
		expectReadPolicy(client, policy, policyVersions, 2)
		expectRemovePolicy(client, 2, 1, 0)

		useClient(client, func() {
			resource.UnitTest(t, resource.TestCase{
				ProtoV6ProviderFactories: testutils.NewProtoV6ProviderFactory(NewSubprovider()),
				Steps: []resource.TestStep{
					{
						Config: testutils.LoadFixtureStringf(t, "%s/policy_create.tf", testDir),
					},
					{
						ImportState: true,
						ImportStateIdFunc: func(_ *terraform.State) (string, error) {
							return "", nil
						},
						ResourceName: "akamai_cloudlets_policy.policy",
						ExpectError:  regexp.MustCompile("policy name cannot be empty"),
					},
				},
			})
		})
		client.AssertExpectations(t)
	})
}

func TestResourcePolicyV3(t *testing.T) {

	type policyAttributes struct {
		version, matchRulesPath string
		description             string
		timeouts                string
		groupID                 int64
	}

	var (
		expectCreatePolicy = func(client *v3.Mock, policyID int64, groupID int64, matchRules v3.MatchRules, description string) (*v3.Policy, *v3.PolicyVersion) {
			policy := &v3.Policy{
				ID:           policyID,
				GroupID:      groupID,
				Name:         "test_policy",
				CloudletType: "ER",
			}
			version := &v3.PolicyVersion{
				PolicyID:      policyID,
				PolicyVersion: 1,
				Description:   ptr.To(description),
				MatchRules:    matchRules,
				MatchRulesWarnings: []v3.MatchRulesWarning{
					{
						Detail:      "test warning details",
						JSONPointer: "/matchRules/1/matches/0",
						Title:       "test warning",
						Type:        "test type",
					},
				},
			}
			client.On("CreatePolicy", testutils.MockContext, v3.CreatePolicyRequest{
				Name:         "test_policy",
				CloudletType: v3.CloudletTypeER,
				GroupID:      groupID,
			}).Return(policy, nil).Once()
			if matchRules == nil && description == "" {
				return policy, nil
			}
			if matchRules != nil {
				client.On("CreatePolicyVersion", testutils.MockContext, v3.CreatePolicyVersionRequest{
					PolicyID: policyID,
					CreatePolicyVersion: v3.CreatePolicyVersion{
						MatchRules:  matchRules,
						Description: ptr.To(description),
					},
				}).Return(version, nil).Once()
			} else {
				client.On("CreatePolicyVersion", testutils.MockContext, v3.CreatePolicyVersionRequest{
					PolicyID: policyID,
					CreatePolicyVersion: v3.CreatePolicyVersion{
						MatchRules:  make(v3.MatchRules, 0),
						Description: ptr.To(description),
					},
				}).Return(version, nil).Once()
			}

			return policy, version
		}

		expectListPolicyVersions = func(client *v3.Mock, policyId int64, versions v3.ListPolicyVersions, times int) {
			client.On("ListPolicyVersions", testutils.MockContext, v3.ListPolicyVersionsRequest{
				PolicyID: policyId,
				Page:     0,
				Size:     1000,
			}).Return(&versions, nil).Times(times)
		}

		expectReadPolicy = func(client *v3.Mock, policy *v3.Policy, version *v3.PolicyVersion, times int) {
			if version != nil {
				var versions v3.ListPolicyVersions
				versions.PolicyVersions = slices.Insert(versions.PolicyVersions, 0, v3.ListPolicyVersionsItem{PolicyVersion: version.PolicyVersion, PolicyID: policy.ID})
				expectListPolicyVersions(client, policy.ID, versions, times)
			} else {
				expectListPolicyVersions(client, policy.ID, v3.ListPolicyVersions{PolicyVersions: make([]v3.ListPolicyVersionsItem, 0)}, times)
			}
			client.On("GetPolicy", testutils.MockContext, v3.GetPolicyRequest{PolicyID: policy.ID}).Return(policy, nil).Times(times)
			if version != nil {
				var versionWithoutWarnings v3.PolicyVersion
				err := copier.CopyWithOption(&versionWithoutWarnings, version, copier.Option{DeepCopy: true})
				require.NoError(t, err)
				versionWithoutWarnings.MatchRulesWarnings = []v3.MatchRulesWarning{}
				versionWithoutWarnings.MatchRules = version.MatchRules
				versionWithoutWarnings.Description = version.Description
				client.On("GetPolicyVersion", testutils.MockContext, v3.GetPolicyVersionRequest{
					PolicyID:      policy.ID,
					PolicyVersion: version.PolicyVersion,
				}).Return(&versionWithoutWarnings, nil).Times(times)
			}
		}

		expectUpdatePolicy = func(client *v3.Mock, policy *v3.Policy, updatedGroup int64) *v3.Policy {
			var policyUpdate v3.Policy
			err := copier.CopyWithOption(&policyUpdate, policy, copier.Option{DeepCopy: true})
			require.NoError(t, err)
			policyUpdate.GroupID = updatedGroup
			client.On("UpdatePolicy", testutils.MockContext, v3.UpdatePolicyRequest{
				Body: v3.UpdatePolicyRequestBody{
					GroupID: updatedGroup,
				},
				PolicyID: policyUpdate.ID,
			}).Return(&policyUpdate, nil).Once()
			return &policyUpdate
		}

		expectCreatePolicyVersion = func(client *v3.Mock, policyID int64, version *v3.PolicyVersion, newMatchRules v3.MatchRules) *v3.PolicyVersion {
			var activatedVersion v3.PolicyVersion
			err := copier.CopyWithOption(&activatedVersion, version, copier.Option{DeepCopy: true})
			require.NoError(t, err)
			activatedVersion.Immutable = true

			client.On("GetPolicyVersion", testutils.MockContext, v3.GetPolicyVersionRequest{
				PolicyID:      policyID,
				PolicyVersion: version.PolicyVersion,
			}).Return(&activatedVersion, nil).Once()

			var versionUpdate v3.PolicyVersion
			err = copier.CopyWithOption(&versionUpdate, activatedVersion, copier.Option{DeepCopy: true})
			require.NoError(t, err)
			versionUpdate.MatchRules = newMatchRules
			versionUpdate.PolicyVersion = version.PolicyVersion + 1
			versionUpdate.Immutable = false

			client.On("CreatePolicyVersion", testutils.MockContext, v3.CreatePolicyVersionRequest{
				CreatePolicyVersion: v3.CreatePolicyVersion{
					Description: ptr.To("test policy description"),
					MatchRules:  newMatchRules,
				},
				PolicyID: policyID,
			}).Return(&versionUpdate, nil).Once()
			return &versionUpdate
		}

		expectUpdatePolicyVersion = func(client *v3.Mock, policyID int64, version *v3.PolicyVersion, newMatchRules v3.MatchRules, newDescription string) *v3.PolicyVersion {
			client.On("GetPolicyVersion", testutils.MockContext, v3.GetPolicyVersionRequest{
				PolicyID:      policyID,
				PolicyVersion: version.PolicyVersion,
			}).Return(version, nil).Once()

			var versionUpdate v3.PolicyVersion
			err := copier.CopyWithOption(&versionUpdate, version, copier.Option{DeepCopy: true})
			require.NoError(t, err)
			versionUpdate.MatchRules = newMatchRules
			versionUpdate.Description = &newDescription
			client.On("UpdatePolicyVersion", testutils.MockContext, v3.UpdatePolicyVersionRequest{
				UpdatePolicyVersion: v3.UpdatePolicyVersion{
					Description: ptr.To(newDescription),
					MatchRules:  newMatchRules,
				},
				PolicyID:      policyID,
				PolicyVersion: version.PolicyVersion,
			}).Return(&versionUpdate, nil).Once()
			return &versionUpdate
		}

		expectRemovePolicy = func(client *v3.Mock, policyID int64) {
			client.On("GetPolicy", testutils.MockContext, v3.GetPolicyRequest{
				PolicyID: policyID,
			}).Return(&v3.Policy{
				CurrentActivations: v3.CurrentActivations{Production: v3.ActivationInfo{}, Staging: v3.ActivationInfo{}},
				ID:                 policyID,
			}, nil).Once()

			client.On("DeletePolicy", testutils.MockContext, v3.DeletePolicyRequest{PolicyID: policyID}).Return(nil).Once()
		}

		expectImportPolicy = func(clientV3 *v3.Mock, clientV2 *cloudlets.Mock, policyID int64, policyName string) {
			listPoliciesV2Resp := []cloudlets.Policy{
				{
					Name: "other-name",
				},
			}
			clientV2.On("ListPolicies", testutils.MockContext, cloudlets.ListPoliciesRequest{
				PageSize: ptr.To(1000),
				Offset:   0,
			}).Return(listPoliciesV2Resp, nil).Once()
			clientV3.On("ListPolicies", testutils.MockContext, v3.ListPoliciesRequest{
				Size: 1000, Page: 0,
			}).Return(&v3.ListPoliciesResponse{
				Content: []v3.Policy{
					{
						ID: policyID, Name: policyName,
					},
				},
			}, nil).Once()
		}

		expectSupressDescriptionChange = func(clientV3 *v3.Mock, policyID, policyVersion int64, version *v3.PolicyVersion, times int) {
			clientV3.On("GetPolicyVersion", testutils.MockContext, v3.GetPolicyVersionRequest{
				PolicyID:      policyID,
				PolicyVersion: policyVersion,
			}).Return(version, nil).Times(times)
		}

		checkPolicyAttributes = func(attrs policyAttributes) resource.TestCheckFunc {
			var matchRulesPath string
			if attrs.matchRulesPath != "" {
				matchRulesPath = testutils.LoadFixtureString(t, attrs.matchRulesPath)
			}
			checkFunc := []resource.TestCheckFunc{
				resource.TestCheckResourceAttr("akamai_cloudlets_policy.policy", "id", "2"),
				resource.TestCheckResourceAttr("akamai_cloudlets_policy.policy", "cloudlet_code", "ER"),
				resource.TestCheckResourceAttr("akamai_cloudlets_policy.policy", "group_id", fmt.Sprintf("%d", attrs.groupID)),
				resource.TestCheckResourceAttr("akamai_cloudlets_policy.policy", "description", attrs.description),
				resource.TestCheckResourceAttr("akamai_cloudlets_policy.policy", "name", "test_policy"),
				resource.TestCheckResourceAttr("akamai_cloudlets_policy.policy", "match_rules", matchRulesPath),
				resource.TestCheckResourceAttr("akamai_cloudlets_policy.policy", "is_shared", "true"),
			}

			if attrs.version != "" {
				checkFunc = append(checkFunc,
					resource.TestCheckResourceAttr("akamai_cloudlets_policy.policy", "version", attrs.version),
				)
			} else {
				checkFunc = append(checkFunc,
					resource.TestCheckNoResourceAttr("akamai_cloudlets_policy.policy", "version"),
				)
			}

			if attrs.timeouts != "" {
				checkFunc = append(checkFunc,
					resource.TestCheckResourceAttr("akamai_cloudlets_policy.policy", "timeouts.#", "1"),
					resource.TestCheckResourceAttr("akamai_cloudlets_policy.policy", "timeouts.0.default", attrs.timeouts),
				)
			} else {
				checkFunc = append(checkFunc,
					resource.TestCheckResourceAttr("akamai_cloudlets_policy.policy", "timeouts.#", "0"),
				)
			}

			return resource.ComposeAggregateTestCheckFunc(checkFunc...)
		}

		commonMatchRules = v3.MatchRules{
			&v3.MatchRuleER{
				Name:                     "r1",
				Type:                     "erMatchRule",
				UseRelativeURL:           "copy_scheme_hostname",
				StatusCode:               301,
				RedirectURL:              "/ddd",
				MatchURL:                 "abc.com",
				UseIncomingSchemeAndHost: true,
			},
			&v3.MatchRuleER{
				Name: "r3",
				Type: "erMatchRule",
				Matches: []v3.MatchCriteriaER{
					{
						MatchType:     "hostname",
						MatchValue:    "3333.dom",
						MatchOperator: "equals",
						CaseSensitive: true,
					},
				},
				UseRelativeURL:           "copy_scheme_hostname",
				StatusCode:               307,
				RedirectURL:              "/abc/sss",
				UseIncomingSchemeAndHost: true,
			},
		}
	)

	t.Run("policy v3 lifecycle with create new version", func(t *testing.T) {
		testDir := "testdata/TestResPolicyV3/lifecycle"

		client := new(v3.Mock)
		policy, version := expectCreatePolicy(client, 2, 123, commonMatchRules, "test policy description")
		expectReadPolicy(client, policy, version, 3)
		policy = expectUpdatePolicy(client, policy, 321)
		version = expectCreatePolicyVersion(client, policy.ID, version, commonMatchRules[:1])
		expectReadPolicy(client, policy, version, 2)
		expectRemovePolicy(client, policy.ID)

		useClientV3(client, func() {
			resource.UnitTest(t, resource.TestCase{
				ProtoV6ProviderFactories: testutils.NewProtoV6ProviderFactory(NewSubprovider()),
				Steps: []resource.TestStep{
					{
						Config: testutils.LoadFixtureStringf(t, "%s/policy_create.tf", testDir),
						Check: checkPolicyAttributes(policyAttributes{
							groupID:        123,
							version:        "1",
							matchRulesPath: fmt.Sprintf("%s/match_rules/match_rules_create.json", testDir),
							description:    "test policy description",
						}),
					},
					{
						Config: testutils.LoadFixtureStringf(t, "%s/policy_update.tf", testDir),
						Check: checkPolicyAttributes(policyAttributes{
							groupID:        321,
							version:        "2",
							matchRulesPath: fmt.Sprintf("%s/match_rules/match_rules_update.json", testDir),
							description:    "test policy description",
						}),
					},
				},
			})
		})
		client.AssertExpectations(t)
	})

	t.Run("policy v3 create policy and update with version drift", func(t *testing.T) {
		testDir := "testdata/TestResPolicyV3/lifecycle_with_drift"

		client := new(v3.Mock)
		policy, version := expectCreatePolicy(client, 2, 123, commonMatchRules, "test policy description")
		expectReadPolicy(client, policy, version, 1)
		version = &v3.PolicyVersion{
			PolicyID:      2,
			PolicyVersion: 2,
			Description:   ptr.To("new description after drift"),
			MatchRules:    commonMatchRules,
			MatchRulesWarnings: []v3.MatchRulesWarning{
				{
					Detail:      "test warning details",
					JSONPointer: "/matchRules/1/matches/0",
					Title:       "test warning",
					Type:        "test type",
				},
			},
		}
		expectReadPolicy(client, policy, version, 1)
		// CustomDiff call to check if the policy version is active
		expectSupressDescriptionChange(client, 2, 2, version, 1)
		expectRemovePolicy(client, policy.ID)

		useClientV3(client, func() {
			resource.UnitTest(t, resource.TestCase{
				ProtoV6ProviderFactories: testutils.NewProtoV6ProviderFactory(NewSubprovider()),
				Steps: []resource.TestStep{
					{
						Config: testutils.LoadFixtureStringf(t, "%s/policy_create.tf", testDir),
						Check: checkPolicyAttributes(policyAttributes{
							groupID:        123,
							version:        "1",
							matchRulesPath: fmt.Sprintf("%s/match_rules/match_rules_create.json", testDir),
							description:    "test policy description",
						}),
						ExpectNonEmptyPlan: true,
					},
				},
			})
		})
		client.AssertExpectations(t)
	})

	t.Run("policy V3 lifecycle, deactivation before delete", func(t *testing.T) {
		testDir := "testdata/TestResPolicyV3/lifecycle"

		client := new(v3.Mock)
		policy, version := expectCreatePolicy(client, 2, 123, commonMatchRules, "test policy description")
		expectReadPolicy(client, policy, version, 2)

		client.On("GetPolicy", testutils.MockContext, v3.GetPolicyRequest{
			PolicyID: policy.ID,
		}).Return(&v3.Policy{
			CurrentActivations: v3.CurrentActivations{Production: v3.ActivationInfo{}, Staging: v3.ActivationInfo{
				Effective: &v3.PolicyActivation{
					Network:       v3.StagingNetwork,
					Operation:     v3.OperationActivation,
					PolicyID:      policy.ID,
					PolicyVersion: version.PolicyVersion,
					Status:        v3.ActivationStatusSuccess,
				},
				Latest: &v3.PolicyActivation{
					Network:       v3.StagingNetwork,
					Operation:     v3.OperationActivation,
					PolicyID:      policy.ID,
					PolicyVersion: version.PolicyVersion,
					Status:        v3.ActivationStatusSuccess,
				},
			}},
			ID:      policy.ID,
			GroupID: 123,
		}, nil).Once()

		client.On("DeactivatePolicy", testutils.MockContext, v3.DeactivatePolicyRequest{
			PolicyID:      policy.ID,
			Network:       v3.StagingNetwork,
			PolicyVersion: version.PolicyVersion,
		}).Return(&v3.PolicyActivation{
			PolicyID:      policy.ID,
			PolicyVersion: version.PolicyVersion,
			Status:        v3.ActivationStatusInProgress,
		}, nil).Once()

		client.On("GetPolicy", testutils.MockContext, v3.GetPolicyRequest{PolicyID: policy.ID}).Return(&v3.Policy{
			CurrentActivations: v3.CurrentActivations{Production: v3.ActivationInfo{}, Staging: v3.ActivationInfo{
				Effective: &v3.PolicyActivation{
					Network:       v3.StagingNetwork,
					Operation:     v3.OperationDeactivation,
					PolicyID:      policy.ID,
					PolicyVersion: version.PolicyVersion,
					Status:        v3.ActivationStatusSuccess,
				},
				Latest: &v3.PolicyActivation{
					Network:       v3.StagingNetwork,
					Operation:     v3.OperationDeactivation,
					PolicyID:      policy.ID,
					PolicyVersion: version.PolicyVersion,
					Status:        v3.ActivationStatusSuccess,
				},
			}},
			ID:      policy.ID,
			GroupID: 123,
		}, nil).Once()

		client.On("DeletePolicy", testutils.MockContext, v3.DeletePolicyRequest{PolicyID: policy.ID}).Return(nil).Once()

		useClientV3(client, func() {
			resource.UnitTest(t, resource.TestCase{
				ProtoV6ProviderFactories: testutils.NewProtoV6ProviderFactory(NewSubprovider()),
				Steps: []resource.TestStep{
					{
						Config: testutils.LoadFixtureStringf(t, "%s/policy_create.tf", testDir),
						Check: checkPolicyAttributes(policyAttributes{
							groupID:        123,
							version:        "1",
							matchRulesPath: fmt.Sprintf("%s/match_rules/match_rules_create.json", testDir),
							description:    "test policy description",
						}),
					},
				},
			})
		})
		client.AssertExpectations(t)
	})

	t.Run("policy V3 lifecycle, in progress deactivation during delete", func(t *testing.T) {
		testDir := "testdata/TestResPolicyV3/lifecycle"

		client := new(v3.Mock)
		policy, version := expectCreatePolicy(client, 2, 123, commonMatchRules, "test policy description")
		expectReadPolicy(client, policy, version, 2)

		client.On("GetPolicy", testutils.MockContext, v3.GetPolicyRequest{
			PolicyID: policy.ID,
		}).Return(&v3.Policy{
			CurrentActivations: v3.CurrentActivations{Production: v3.ActivationInfo{}, Staging: v3.ActivationInfo{
				Effective: &v3.PolicyActivation{
					Network:       v3.StagingNetwork,
					Operation:     v3.OperationActivation,
					PolicyID:      policy.ID,
					PolicyVersion: version.PolicyVersion,
					Status:        v3.ActivationStatusSuccess,
				},
				Latest: &v3.PolicyActivation{
					Network:       v3.StagingNetwork,
					Operation:     v3.OperationDeactivation,
					PolicyID:      policy.ID,
					PolicyVersion: version.PolicyVersion,
					Status:        v3.ActivationStatusInProgress,
				},
			}},
			ID:      policy.ID,
			GroupID: 123,
		}, nil).Once()

		client.On("GetPolicy", testutils.MockContext, v3.GetPolicyRequest{
			PolicyID: policy.ID,
		}).Return(&v3.Policy{
			CurrentActivations: v3.CurrentActivations{Production: v3.ActivationInfo{}, Staging: v3.ActivationInfo{
				Effective: &v3.PolicyActivation{
					Network:       v3.StagingNetwork,
					Operation:     v3.OperationDeactivation,
					PolicyID:      policy.ID,
					PolicyVersion: version.PolicyVersion,
					Status:        v3.ActivationStatusSuccess,
				},
				Latest: &v3.PolicyActivation{
					Network:       v3.StagingNetwork,
					Operation:     v3.OperationDeactivation,
					PolicyID:      policy.ID,
					PolicyVersion: version.PolicyVersion,
					Status:        v3.ActivationStatusSuccess,
				},
			}},
			ID:      policy.ID,
			GroupID: 123,
		}, nil).Once()

		client.On("DeletePolicy", testutils.MockContext, v3.DeletePolicyRequest{PolicyID: policy.ID}).Return(nil).Once()

		useClientV3(client, func() {
			resource.UnitTest(t, resource.TestCase{
				ProtoV6ProviderFactories: testutils.NewProtoV6ProviderFactory(NewSubprovider()),
				Steps: []resource.TestStep{
					{
						Config: testutils.LoadFixtureStringf(t, "%s/policy_create.tf", testDir),
						Check: checkPolicyAttributes(policyAttributes{
							groupID:        123,
							version:        "1",
							matchRulesPath: fmt.Sprintf("%s/match_rules/match_rules_create.json", testDir),
							description:    "test policy description",
						}),
					},
				},
			})
		})
		client.AssertExpectations(t)
	})

	t.Run("policy v3 lifecycle with update existing version", func(t *testing.T) {
		testDir := "testdata/TestResPolicyV3/lifecycle"

		client := new(v3.Mock)
		policy, version := expectCreatePolicy(client, 2, 123, commonMatchRules, "test policy description")
		expectReadPolicy(client, policy, version, 3)
		policy = expectUpdatePolicy(client, policy, 321)
		version = expectUpdatePolicyVersion(client, policy.ID, version, commonMatchRules[:1], "test policy description")
		expectReadPolicy(client, policy, version, 2)
		expectRemovePolicy(client, policy.ID)

		useClientV3(client, func() {
			resource.UnitTest(t, resource.TestCase{
				ProtoV6ProviderFactories: testutils.NewProtoV6ProviderFactory(NewSubprovider()),
				Steps: []resource.TestStep{
					{
						Config: testutils.LoadFixtureStringf(t, "%s/policy_create.tf", testDir),
						Check: checkPolicyAttributes(policyAttributes{
							groupID:        123,
							version:        "1",
							matchRulesPath: fmt.Sprintf("%s/match_rules/match_rules_create.json", testDir),
							description:    "test policy description",
						}),
					},
					{
						Config: testutils.LoadFixtureStringf(t, "%s/policy_update.tf", testDir),
						Check: checkPolicyAttributes(policyAttributes{
							groupID:        321,
							version:        "1",
							matchRulesPath: fmt.Sprintf("%s/match_rules/match_rules_update.json", testDir),
							description:    "test policy description",
						}),
					},
				},
			})
		})
		client.AssertExpectations(t)
	})

	t.Run("update only policy v3", func(t *testing.T) {
		testDir := "testdata/TestResPolicyV3/lifecycle_policy_update"

		client := new(v3.Mock)
		policy, version := expectCreatePolicy(client, 2, 123, commonMatchRules, "test policy description")
		expectReadPolicy(client, policy, version, 3)
		policy = expectUpdatePolicy(client, policy, 321)
		expectReadPolicy(client, policy, version, 2)
		expectRemovePolicy(client, policy.ID)

		useClientV3(client, func() {
			resource.UnitTest(t, resource.TestCase{
				ProtoV6ProviderFactories: testutils.NewProtoV6ProviderFactory(NewSubprovider()),
				Steps: []resource.TestStep{
					{
						Config: testutils.LoadFixtureStringf(t, "%s/policy_create.tf", testDir),
						Check: checkPolicyAttributes(policyAttributes{
							groupID:        123,
							version:        "1",
							matchRulesPath: fmt.Sprintf("%s/match_rules/match_rules.json", testDir),
							description:    "test policy description",
						}),
					},
					{
						Config: testutils.LoadFixtureStringf(t, "%s/policy_update.tf", testDir),
						Check: checkPolicyAttributes(policyAttributes{
							groupID:        321,
							version:        "1",
							matchRulesPath: fmt.Sprintf("%s/match_rules/match_rules.json", testDir),
							description:    "test policy description",
						}),
					},
				},
			})
		})
		client.AssertExpectations(t)
	})

	t.Run("update only description for not active v3 policy - expect no new version", func(t *testing.T) {
		testDir := "testdata/TestResPolicyV3/lifecycle_description_update"

		client := new(v3.Mock)
		policy, version := expectCreatePolicy(client, 2, 123, commonMatchRules, "test policy description")
		expectReadPolicy(client, policy, version, 3)
		// CustomDiff calls to check if the policy version is active
		expectSupressDescriptionChange(client, policy.ID, version.PolicyVersion, version, 2)
		policyVersionUpdated := expectUpdatePolicyVersion(client, 2, version, commonMatchRules, "test policy description - updated")
		expectReadPolicy(client, policy, policyVersionUpdated, 2)
		expectRemovePolicy(client, policy.ID)

		useClientV3(client, func() {
			resource.UnitTest(t, resource.TestCase{
				ProtoV6ProviderFactories: testutils.NewProtoV6ProviderFactory(NewSubprovider()),
				Steps: []resource.TestStep{
					{
						Config: testutils.LoadFixtureStringf(t, "%s/policy_create.tf", testDir),
						Check: checkPolicyAttributes(policyAttributes{
							groupID:        123,
							version:        "1",
							matchRulesPath: fmt.Sprintf("%s/match_rules/match_rules.json", testDir),
							description:    "test policy description",
						}),
					},
					{
						Config: testutils.LoadFixtureStringf(t, "%s/policy_update.tf", testDir),
						Check: checkPolicyAttributes(policyAttributes{
							groupID:        123,
							version:        "1",
							matchRulesPath: fmt.Sprintf("%s/match_rules/match_rules.json", testDir),
							description:    "test policy description - updated",
						}),
					},
				},
			})
		})
		client.AssertExpectations(t)
	})

	t.Run("update only description for active v3 policy - expect no changes", func(t *testing.T) {
		testDir := "testdata/TestResPolicyV3/lifecycle_description_update"

		client := new(v3.Mock)
		policy, version := expectCreatePolicy(client, 2, 123, commonMatchRules, "test policy description")
		expectReadPolicy(client, policy, version, 3)
		version.Immutable = true
		// CustomDiff calls to check if the policy version is active
		expectSupressDescriptionChange(client, policy.ID, version.PolicyVersion, version, 3)
		expectReadPolicy(client, policy, version, 1)
		expectRemovePolicy(client, policy.ID)

		useClientV3(client, func() {
			resource.UnitTest(t, resource.TestCase{
				ProtoV6ProviderFactories: testutils.NewProtoV6ProviderFactory(NewSubprovider()),
				Steps: []resource.TestStep{
					{
						Config: testutils.LoadFixtureStringf(t, "%s/policy_create.tf", testDir),
						Check: checkPolicyAttributes(policyAttributes{
							groupID:        123,
							version:        "1",
							matchRulesPath: fmt.Sprintf("%s/match_rules/match_rules.json", testDir),
							description:    "test policy description",
						}),
					},
					{
						Config: testutils.LoadFixtureStringf(t, "%s/policy_update.tf", testDir),
						Check: checkPolicyAttributes(policyAttributes{
							groupID:        123,
							version:        "1",
							matchRulesPath: fmt.Sprintf("%s/match_rules/match_rules.json", testDir),
							description:    "test policy description",
						}),
					},
				},
			})
		})
		client.AssertExpectations(t)
	})

	t.Run("update only version for v3 policy", func(t *testing.T) {
		testDir := "testdata/TestResPolicyV3/lifecycle_version_update"

		client := new(v3.Mock)
		policy, version := expectCreatePolicy(client, 2, 123, commonMatchRules, "test policy description")
		expectReadPolicy(client, policy, version, 3)
		version = expectUpdatePolicyVersion(client, policy.ID, version, commonMatchRules[:1], "test policy description")
		expectReadPolicy(client, policy, version, 2)
		expectRemovePolicy(client, policy.ID)

		useClientV3(client, func() {
			resource.UnitTest(t, resource.TestCase{
				ProtoV6ProviderFactories: testutils.NewProtoV6ProviderFactory(NewSubprovider()),
				Steps: []resource.TestStep{
					{
						Config: testutils.LoadFixtureStringf(t, "%s/policy_create.tf", testDir),
						Check: checkPolicyAttributes(policyAttributes{
							groupID:        123,
							version:        "1",
							matchRulesPath: fmt.Sprintf("%s/match_rules/match_rules_create.json", testDir),
							description:    "test policy description",
						}),
					},
					{
						Config: testutils.LoadFixtureStringf(t, "%s/policy_update.tf", testDir),
						Check: checkPolicyAttributes(policyAttributes{
							groupID:        123,
							version:        "1",
							matchRulesPath: fmt.Sprintf("%s/match_rules/match_rules_update.json", testDir),
							description:    "test policy description",
						}),
					},
				},
			})
		})
		client.AssertExpectations(t)
	})

	t.Run("warnings creating and updating version for v3 policy", func(t *testing.T) {
		testDir := "testdata/TestResPolicyV3/lifecycle_version_update"

		client := new(v3.Mock)
		policy, version := expectCreatePolicy(client, 2, 123, commonMatchRules, "test policy description")
		expectReadPolicy(client, policy, version, 3)
		version = expectUpdatePolicyVersion(client, policy.ID, version, commonMatchRules[:1], "test policy description")
		expectReadPolicy(client, policy, version, 4)
		expectRemovePolicy(client, policy.ID)

		warningsJSON, err := warningsToJSON(version.MatchRulesWarnings)
		require.NoError(t, err)

		checkWarnings := resource.ComposeAggregateTestCheckFunc(
			resource.TestCheckResourceAttr("akamai_cloudlets_policy.policy", "warnings", string(warningsJSON)),
			resource.TestMatchOutput("policy_output", regexp.MustCompile("test warning")),
		)

		useClientV3(client, func() {
			resource.UnitTest(t, resource.TestCase{
				ProtoV6ProviderFactories: testutils.NewProtoV6ProviderFactory(NewSubprovider()),
				Steps: []resource.TestStep{
					{
						Config: testutils.LoadFixtureStringf(t, "%s/policy_create.tf", testDir),
						Check:  checkWarnings,
					},
					{
						Config: testutils.LoadFixtureStringf(t, "%s/policy_update.tf", testDir),
						Check:  checkWarnings,
					},
					{
						Config: testutils.LoadFixtureStringf(t, "%s/policy_update.tf", testDir),
						Check:  checkWarnings,
					},
				},
			})
		})
		client.AssertExpectations(t)
	})

	t.Run("remove match rules from version for v3 policy", func(t *testing.T) {
		testDir := "testdata/TestResPolicyV3/lifecycle_remove_match_rules"

		client := new(v3.Mock)
		policy, version := expectCreatePolicy(client, 2, 123, commonMatchRules, "test policy description")
		expectReadPolicy(client, policy, version, 3)
		version = expectUpdatePolicyVersion(client, policy.ID, version, v3.MatchRules{}, "test policy description")
		expectReadPolicy(client, policy, version, 2)
		expectRemovePolicy(client, policy.ID)

		useClientV3(client, func() {
			resource.UnitTest(t, resource.TestCase{
				ProtoV6ProviderFactories: testutils.NewProtoV6ProviderFactory(NewSubprovider()),
				Steps: []resource.TestStep{
					{
						Config: testutils.LoadFixtureStringf(t, "%s/policy_create.tf", testDir),
						Check: checkPolicyAttributes(policyAttributes{
							groupID:        123,
							version:        "1",
							matchRulesPath: fmt.Sprintf("%s/match_rules/match_rules_create.json", testDir),
							description:    "test policy description",
						}),
					},
					{
						Config: testutils.LoadFixtureStringf(t, "%s/policy_update.tf", testDir),
						Check: checkPolicyAttributes(policyAttributes{
							groupID:        123,
							version:        "1",
							matchRulesPath: "",
							description:    "test policy description",
						}),
					},
				},
			})
		})
		client.AssertExpectations(t)
	})

	t.Run("create v3 policy without match rules or description", func(t *testing.T) {
		testDir := "testdata/TestResPolicyV3/create_no_match_rules_no_description"

		client := new(v3.Mock)
		policy, version := expectCreatePolicy(client, 2, 123, nil, "")
		expectReadPolicy(client, policy, version, 2)
		expectRemovePolicy(client, 2)
		useClientV3(client, func() {
			resource.UnitTest(t, resource.TestCase{
				ProtoV6ProviderFactories: testutils.NewProtoV6ProviderFactory(NewSubprovider()),
				Steps: []resource.TestStep{
					{
						Config: testutils.LoadFixtureStringf(t, "%s/policy_create.tf", testDir),
						Check: checkPolicyAttributes(policyAttributes{
							groupID:        123,
							version:        "",
							matchRulesPath: "",
							description:    "",
						}),
					},
				},
			})
		})
		client.AssertExpectations(t)
	})

	t.Run("create v3 policy without version, update to create new version", func(t *testing.T) {
		testDir := "testdata/TestResPolicyV3/lifecycle_no_version"

		client := new(v3.Mock)
		policy, version := expectCreatePolicy(client, 2, 123, nil, "")
		expectReadPolicy(client, policy, version, 3)

		version = &v3.PolicyVersion{
			Description:   ptr.To("test policy description"),
			PolicyID:      2,
			PolicyVersion: 1,
		}
		client.On("CreatePolicyVersion", testutils.MockContext, v3.CreatePolicyVersionRequest{
			CreatePolicyVersion: v3.CreatePolicyVersion{
				Description: ptr.To("test policy description"),
				MatchRules:  make(v3.MatchRules, 0),
			},
			PolicyID: 2,
		}).Return(version, nil).Once()
		expectReadPolicy(client, policy, version, 2)
		expectRemovePolicy(client, 2)
		useClientV3(client, func() {
			resource.UnitTest(t, resource.TestCase{
				ProtoV6ProviderFactories: testutils.NewProtoV6ProviderFactory(NewSubprovider()),
				Steps: []resource.TestStep{
					{
						Config: testutils.LoadFixtureStringf(t, "%s/policy_create.tf", testDir),
						Check: checkPolicyAttributes(policyAttributes{
							groupID:        123,
							version:        "",
							matchRulesPath: "",
							description:    "",
						}),
					},
					{
						Config: testutils.LoadFixtureStringf(t, "%s/policy_update.tf", testDir),
						Check: checkPolicyAttributes(policyAttributes{
							groupID:        123,
							version:        "1",
							matchRulesPath: "",
							description:    "test policy description",
						}),
					},
				},
			})
		})
		client.AssertExpectations(t)
	})

	t.Run("create V3 policy without match rules with description", func(t *testing.T) {
		testDir := "testdata/TestResPolicyV3/create_no_match_rules"

		client := new(v3.Mock)
		policy, version := expectCreatePolicy(client, 2, 123, nil, "test policy description")
		expectReadPolicy(client, policy, version, 2)
		expectRemovePolicy(client, 2)
		useClientV3(client, func() {
			resource.UnitTest(t, resource.TestCase{
				ProtoV6ProviderFactories: testutils.NewProtoV6ProviderFactory(NewSubprovider()),
				Steps: []resource.TestStep{
					{
						Config: testutils.LoadFixtureStringf(t, "%s/policy_create.tf", testDir),
						Check: checkPolicyAttributes(policyAttributes{
							groupID:        123,
							version:        "1",
							matchRulesPath: "",
							description:    "test policy description",
						}),
					},
				},
			})
		})
		client.AssertExpectations(t)
	})

	t.Run("create V3 policy with timeout", func(t *testing.T) {
		testDir := "testdata/TestResPolicyV3/timeouts"

		matchRules := v3.MatchRules{
			&v3.MatchRuleER{
				Name:                     "r1",
				Type:                     "erMatchRule",
				UseRelativeURL:           "copy_scheme_hostname",
				StatusCode:               301,
				RedirectURL:              "/ddd",
				MatchURL:                 "abc.com",
				UseIncomingSchemeAndHost: true,
			},
		}

		client := new(v3.Mock)
		policy, version := expectCreatePolicy(client, 2, 123, matchRules, "test policy description")
		expectReadPolicy(client, policy, version, 2)
		expectRemovePolicy(client, 2)
		useClientV3(client, func() {
			resource.UnitTest(t, resource.TestCase{
				ProtoV6ProviderFactories: testutils.NewProtoV6ProviderFactory(NewSubprovider()),
				Steps: []resource.TestStep{
					{
						Config: testutils.LoadFixtureStringf(t, "%s/policy_create.tf", testDir),
						Check: checkPolicyAttributes(policyAttributes{
							groupID:        123,
							version:        "1",
							matchRulesPath: fmt.Sprintf("%s/match_rules/match_rules.json", testDir),
							timeouts:       "4h",
							description:    "test policy description",
						}),
					},
				},
			})
		})
		client.AssertExpectations(t)
	})

	t.Run("error creating V3 policy", func(t *testing.T) {
		testDir := "testdata/TestResPolicyV3/lifecycle"

		client := new(v3.Mock)
		client.On("CreatePolicy", testutils.MockContext, v3.CreatePolicyRequest{
			Name:         "test_policy",
			CloudletType: v3.CloudletTypeER,
			GroupID:      123,
		}).Return(nil, fmt.Errorf("oops")).Once()

		useClientV3(client, func() {
			resource.UnitTest(t, resource.TestCase{
				ProtoV6ProviderFactories: testutils.NewProtoV6ProviderFactory(NewSubprovider()),
				Steps: []resource.TestStep{
					{
						Config:      testutils.LoadFixtureStringf(t, "%s/policy_create.tf", testDir),
						ExpectError: regexp.MustCompile("oops"),
					},
				},
			})
		})
		client.AssertExpectations(t)
	})

	t.Run("error first update v3 policy version", func(t *testing.T) {
		testDir := "testdata/TestResPolicyV3/lifecycle"

		policy := &v3.Policy{
			ID:           2,
			GroupID:      123,
			Name:         "test_policy",
			CloudletType: "ER",
		}

		expectErrorCreatingVersion := func(client *v3.Mock) {
			client.On("CreatePolicy", testutils.MockContext, v3.CreatePolicyRequest{
				Name:         "test_policy",
				CloudletType: v3.CloudletTypeER,
				GroupID:      123,
			}).Return(policy, nil)
			client.On("CreatePolicyVersion", testutils.MockContext, v3.CreatePolicyVersionRequest{
				CreatePolicyVersion: v3.CreatePolicyVersion{
					Description: ptr.To("test policy description"),
					MatchRules:  commonMatchRules,
				},
				PolicyID: 2,
			}).Return(nil, fmt.Errorf("CreatePolicyVersionError"))
		}

		testCases := []struct {
			Expectations  func(client *v3.Mock)
			ExpectedError *regexp.Regexp
		}{
			{
				Expectations: func(client *v3.Mock) {
					expectErrorCreatingVersion(client)
					expectReadPolicy(client, policy, &v3.PolicyVersion{
						PolicyID:      2,
						PolicyVersion: 1,
					}, 1)
					expectRemovePolicy(client, 2)
				},
				ExpectedError: regexp.MustCompile("CreatePolicyVersionError"),
			},
			{
				Expectations: func(client *v3.Mock) {
					expectErrorCreatingVersion(client)
					expectListPolicyVersions(client, policy.ID, v3.ListPolicyVersions{
						PolicyVersions: []v3.ListPolicyVersionsItem{
							{PolicyVersion: 1, PolicyID: policy.ID},
						}}, 1)
					client.On("GetPolicy", testutils.MockContext, v3.GetPolicyRequest{PolicyID: policy.ID}).Return(nil, fmt.Errorf("GetPolicyError")).Once()
					expectRemovePolicy(client, 2)
				},
				ExpectedError: regexp.MustCompile("(?s)GetPolicyError.*CreatePolicyVersionError"),
			},
		}

		for i := range testCases {
			client := new(v3.Mock)
			testCases[i].Expectations(client)
			useClientV3(client, func() {
				resource.UnitTest(t, resource.TestCase{
					ProtoV6ProviderFactories: testutils.NewProtoV6ProviderFactory(NewSubprovider()),
					Steps: []resource.TestStep{
						{
							Config:      testutils.LoadFixtureStringf(t, "%s/policy_create.tf", testDir),
							ExpectError: testCases[i].ExpectedError,
						},
					},
				})
			})
			client.AssertExpectations(t)
		}
	})

	t.Run("error fetching V3 policy", func(t *testing.T) {
		testDir := "testdata/TestResPolicyV3/create_no_match_rules"

		client := new(v3.Mock)
		policy, version := expectCreatePolicy(client, 2, 123, nil, "test policy description")
		expectListPolicyVersions(client, policy.ID, v3.ListPolicyVersions{
			PolicyVersions: []v3.ListPolicyVersionsItem{
				{PolicyVersion: version.PolicyVersion, PolicyID: policy.ID},
			}}, 1)
		client.On("GetPolicy", testutils.MockContext, v3.GetPolicyRequest{PolicyID: policy.ID}).Return(nil, fmt.Errorf("oops")).Once()
		expectRemovePolicy(client, 2)
		useClientV3(client, func() {
			resource.UnitTest(t, resource.TestCase{
				ProtoV6ProviderFactories: testutils.NewProtoV6ProviderFactory(NewSubprovider()),
				Steps: []resource.TestStep{
					{
						Config:      testutils.LoadFixtureStringf(t, "%s/policy_create.tf", testDir),
						ExpectError: regexp.MustCompile("oops"),
					},
				},
			})
		})
		client.AssertExpectations(t)
	})

	t.Run("error fetching V3 policy version", func(t *testing.T) {
		testDir := "testdata/TestResPolicyV3/create_no_match_rules"

		client := new(v3.Mock)
		policy, version := expectCreatePolicy(client, 2, 123, nil, "test policy description")
		expectListPolicyVersions(client, policy.ID, v3.ListPolicyVersions{
			PolicyVersions: []v3.ListPolicyVersionsItem{
				{PolicyVersion: version.PolicyVersion, PolicyID: policy.ID},
			}}, 1)
		client.On("GetPolicy", testutils.MockContext, v3.GetPolicyRequest{PolicyID: policy.ID}).Return(policy, nil).Once()
		client.On("GetPolicyVersion", testutils.MockContext, v3.GetPolicyVersionRequest{
			PolicyID:      policy.ID,
			PolicyVersion: version.PolicyVersion,
		}).Return(nil, fmt.Errorf("oops")).Once()
		expectRemovePolicy(client, 2)
		useClientV3(client, func() {
			resource.UnitTest(t, resource.TestCase{
				ProtoV6ProviderFactories: testutils.NewProtoV6ProviderFactory(NewSubprovider()),
				Steps: []resource.TestStep{
					{
						Config:      testutils.LoadFixtureStringf(t, "%s/policy_create.tf", testDir),
						ExpectError: regexp.MustCompile("oops"),
					},
				},
			})
		})
		client.AssertExpectations(t)
	})

	t.Run("error updating V3 policy", func(t *testing.T) {
		testDir := "testdata/TestResPolicyV3/lifecycle_policy_update"

		client := new(v3.Mock)
		policy, version := expectCreatePolicy(client, 2, 123, commonMatchRules, "test policy description")
		expectReadPolicy(client, policy, version, 3)
		client.On("UpdatePolicy", testutils.MockContext, v3.UpdatePolicyRequest{
			Body: v3.UpdatePolicyRequestBody{
				GroupID: 321,
			},
			PolicyID: policy.ID,
		}).Return(nil, fmt.Errorf("oops")).Once()
		expectRemovePolicy(client, policy.ID)

		useClientV3(client, func() {
			resource.UnitTest(t, resource.TestCase{
				ProtoV6ProviderFactories: testutils.NewProtoV6ProviderFactory(NewSubprovider()),
				Steps: []resource.TestStep{
					{
						Config: testutils.LoadFixtureStringf(t, "%s/policy_create.tf", testDir),
					},
					{
						Config:      testutils.LoadFixtureStringf(t, "%s/policy_update.tf", testDir),
						ExpectError: regexp.MustCompile("oops"),
					},
				},
			})
		})
		client.AssertExpectations(t)
	})

	t.Run("error updating version in v3 policy", func(t *testing.T) {
		testDir := "testdata/TestResPolicyV3/lifecycle_version_update"

		expectErrorUpdatingVersion := func(client *v3.Mock, expectReadPolicyTimes int) (policy *v3.Policy) {
			policy, version := expectCreatePolicy(client, 2, 123, commonMatchRules, "test policy description")
			expectReadPolicy(client, policy, version, expectReadPolicyTimes)
			client.On("GetPolicyVersion", testutils.MockContext, v3.GetPolicyVersionRequest{
				PolicyID:      policy.ID,
				PolicyVersion: version.PolicyVersion,
			}).Return(version, nil).Once()
			client.On("UpdatePolicyVersion", testutils.MockContext, v3.UpdatePolicyVersionRequest{
				UpdatePolicyVersion: v3.UpdatePolicyVersion{
					Description: ptr.To("test policy description"),
					MatchRules:  commonMatchRules[:1],
				},
				PolicyID:      policy.ID,
				PolicyVersion: version.PolicyVersion,
			}).Return(nil, fmt.Errorf("UpdatePolicyVersionError")).Once()
			return
		}

		testCases := []struct {
			Expectations  func(client *v3.Mock)
			ExpectedError *regexp.Regexp
		}{
			{
				Expectations: func(client *v3.Mock) {
					policy := expectErrorUpdatingVersion(client, 4)
					expectRemovePolicy(client, policy.ID)
				},
				ExpectedError: regexp.MustCompile("UpdatePolicyVersionError"),
			},
			{
				Expectations: func(client *v3.Mock) {
					policy := expectErrorUpdatingVersion(client, 3)
					expectListPolicyVersions(client, policy.ID, v3.ListPolicyVersions{
						PolicyVersions: []v3.ListPolicyVersionsItem{
							{PolicyVersion: 1, PolicyID: policy.ID},
						}}, 1)
					client.On("GetPolicy", testutils.MockContext, v3.GetPolicyRequest{PolicyID: policy.ID}).Return(nil, fmt.Errorf("GetPolicyError")).Once()
					expectRemovePolicy(client, policy.ID)
				},
				ExpectedError: regexp.MustCompile("(?s)GetPolicyError.*UpdatePolicyVersionError"),
			},
		}

		for i := range testCases {
			client := new(v3.Mock)
			testCases[i].Expectations(client)
			useClientV3(client, func() {
				resource.UnitTest(t, resource.TestCase{
					ProtoV6ProviderFactories: testutils.NewProtoV6ProviderFactory(NewSubprovider()),
					Steps: []resource.TestStep{
						{
							Config: testutils.LoadFixtureStringf(t, "%s/policy_create.tf", testDir),
						},
						{
							Config:      testutils.LoadFixtureStringf(t, "%s/policy_update.tf", testDir),
							ExpectError: testCases[i].ExpectedError,
						},
					},
				})
			})
			client.AssertExpectations(t)
		}
	})

	t.Run("import policy v3", func(t *testing.T) {
		testDir := "testdata/TestResPolicyV3"
		clientV2 := new(cloudlets.Mock)
		clientV3 := new(v3.Mock)

		policy, version := expectCreatePolicy(clientV3, 2, 123, nil, "test policy description")
		expectReadPolicy(clientV3, policy, version, 3)
		expectImportPolicy(clientV3, clientV2, policy.ID, "test_policy")
		expectRemovePolicy(clientV3, policy.ID)

		useClientV2AndV3(clientV2, clientV3, func() {
			resource.UnitTest(t, resource.TestCase{
				ProtoV6ProviderFactories: testutils.NewProtoV6ProviderFactory(NewSubprovider()),
				Steps: []resource.TestStep{
					{
						Config: testutils.LoadFixtureStringf(t, "%s/create_no_match_rules/policy_create.tf", testDir),
						Check: checkPolicyAttributes(policyAttributes{
							groupID:     123,
							version:     "1",
							description: "test policy description",
						}),
					},
					{
						ImportState:             true,
						ImportStateId:           "test_policy",
						ResourceName:            "akamai_cloudlets_policy.policy",
						ImportStateVerify:       true,
						ImportStateVerifyIgnore: []string{"warnings"},
					},
				},
			})
		})
		clientV2.AssertExpectations(t)
		clientV3.AssertExpectations(t)
	})

	t.Run("import policy v3 without version", func(t *testing.T) {
		testDir := "testdata/TestResPolicyV3"
		clientV2 := new(cloudlets.Mock)
		clientV3 := new(v3.Mock)

		policy, version := expectCreatePolicy(clientV3, 2, 123, nil, "")
		expectReadPolicy(clientV3, policy, version, 3)
		expectImportPolicy(clientV3, clientV2, policy.ID, "test_policy")
		expectRemovePolicy(clientV3, policy.ID)

		useClientV2AndV3(clientV2, clientV3, func() {
			resource.UnitTest(t, resource.TestCase{
				ProtoV6ProviderFactories: testutils.NewProtoV6ProviderFactory(NewSubprovider()),
				Steps: []resource.TestStep{
					{
						Config: testutils.LoadFixtureStringf(t, "%s/create_no_match_rules_no_description/policy_create.tf", testDir),
						Check: checkPolicyAttributes(policyAttributes{
							groupID:     123,
							version:     "",
							description: "",
						}),
					},
					{
						ImportState:             true,
						ImportStateId:           "test_policy",
						ResourceName:            "akamai_cloudlets_policy.policy",
						ImportStateVerify:       true,
						ImportStateVerifyIgnore: []string{"warnings"},
					},
				},
			})
		})
		clientV2.AssertExpectations(t)
		clientV3.AssertExpectations(t)
	})

	t.Run("import policy v3 - no policy found", func(t *testing.T) {
		testDir := "testdata/TestResPolicyV3"
		clientV2 := new(cloudlets.Mock)
		clientV3 := new(v3.Mock)

		policy, version := expectCreatePolicy(clientV3, 2, 123, nil, "test policy description")
		expectReadPolicy(clientV3, policy, version, 2)
		clientV2.On("ListPolicies", testutils.MockContext, cloudlets.ListPoliciesRequest{
			PageSize: ptr.To(1000),
			Offset:   0,
		}).Return([]cloudlets.Policy{}, nil).Once()
		clientV3.On("ListPolicies", testutils.MockContext, v3.ListPoliciesRequest{
			Size: 1000, Page: 0,
		}).Return(&v3.ListPoliciesResponse{
			Content: []v3.Policy{},
		}, nil).Once()
		expectRemovePolicy(clientV3, policy.ID)

		useClientV2AndV3(clientV2, clientV3, func() {
			resource.UnitTest(t, resource.TestCase{
				ProtoV6ProviderFactories: testutils.NewProtoV6ProviderFactory(NewSubprovider()),
				Steps: []resource.TestStep{
					{
						Config: testutils.LoadFixtureStringf(t, "%s/create_no_match_rules/policy_create.tf", testDir),
						Check: checkPolicyAttributes(policyAttributes{
							groupID:     123,
							version:     "1",
							description: "test policy description",
						}),
					},
					{
						ImportState:             true,
						ImportStateId:           "test_policy",
						ResourceName:            "akamai_cloudlets_policy.policy",
						ImportStateVerify:       true,
						ImportStateVerifyIgnore: []string{"warnings"},
						ExpectError:             regexp.MustCompile("policy 'test_policy' does not exist"),
					},
				},
			})
		})
		clientV2.AssertExpectations(t)
		clientV3.AssertExpectations(t)
	})

	t.Run("import policy v3 - v2 api error, v3 policy found", func(t *testing.T) {
		testDir := "testdata/TestResPolicyV3"
		clientV2 := new(cloudlets.Mock)
		clientV3 := new(v3.Mock)

		policy, version := expectCreatePolicy(clientV3, 2, 123, nil, "test policy description")
		expectReadPolicy(clientV3, policy, version, 3)
		clientV2.On("ListPolicies", testutils.MockContext, cloudlets.ListPoliciesRequest{
			PageSize: ptr.To(1000),
			Offset:   0,
		}).Return(nil, fmt.Errorf("v2 api error")).Once()
		clientV3.On("ListPolicies", testutils.MockContext, v3.ListPoliciesRequest{
			Size: 1000, Page: 0,
		}).Return(&v3.ListPoliciesResponse{
			Content: []v3.Policy{
				{
					ID: policy.ID, Name: policy.Name,
				},
			},
		}, nil).Once()
		expectRemovePolicy(clientV3, policy.ID)

		useClientV2AndV3(clientV2, clientV3, func() {
			resource.UnitTest(t, resource.TestCase{
				ProtoV6ProviderFactories: testutils.NewProtoV6ProviderFactory(NewSubprovider()),
				Steps: []resource.TestStep{
					{
						Config: testutils.LoadFixtureStringf(t, "%s/create_no_match_rules/policy_create.tf", testDir),
						Check: checkPolicyAttributes(policyAttributes{
							groupID:     123,
							version:     "1",
							description: "test policy description",
						}),
					},
					{
						ImportState:             true,
						ImportStateId:           "test_policy",
						ResourceName:            "akamai_cloudlets_policy.policy",
						ImportStateVerify:       true,
						ImportStateVerifyIgnore: []string{"warnings"},
					},
				},
			})
		})
		clientV2.AssertExpectations(t)
		clientV3.AssertExpectations(t)
	})

	t.Run("import policy v3 - v2 and v3 api error", func(t *testing.T) {
		testDir := "testdata/TestResPolicyV3"
		clientV2 := new(cloudlets.Mock)
		clientV3 := new(v3.Mock)

		policy, version := expectCreatePolicy(clientV3, 2, 123, nil, "test policy description")
		expectReadPolicy(clientV3, policy, version, 2)
		clientV2.On("ListPolicies", testutils.MockContext, cloudlets.ListPoliciesRequest{
			PageSize: ptr.To(1000),
			Offset:   0,
		}).Return(nil, fmt.Errorf("v2 api error")).Once()
		clientV3.On("ListPolicies", testutils.MockContext, v3.ListPoliciesRequest{
			Size: 1000, Page: 0,
		}).Return(nil, fmt.Errorf("v3 api error")).Once()
		expectRemovePolicy(clientV3, policy.ID)

		useClientV2AndV3(clientV2, clientV3, func() {
			resource.UnitTest(t, resource.TestCase{
				ProtoV6ProviderFactories: testutils.NewProtoV6ProviderFactory(NewSubprovider()),
				Steps: []resource.TestStep{
					{
						Config: testutils.LoadFixtureStringf(t, "%s/create_no_match_rules/policy_create.tf", testDir),
						Check: checkPolicyAttributes(policyAttributes{
							groupID:     123,
							version:     "1",
							description: "test policy description",
						}),
					},
					{
						ImportState:             true,
						ImportStateId:           "test_policy",
						ResourceName:            "akamai_cloudlets_policy.policy",
						ImportStateVerify:       true,
						ImportStateVerifyIgnore: []string{"warnings"},
						ExpectError:             regexp.MustCompile("could not list V2 policies: v2 api error\ncould not list V3 policies: v3 api error"),
					},
				},
			})
		})
		clientV2.AssertExpectations(t)
		clientV3.AssertExpectations(t)
	})

	t.Run("import policy v3 - test checkForV3Policy()", func(t *testing.T) {
		testDir := "testdata/TestResPolicyV3"
		clientV2 := new(cloudlets.Mock)
		clientV3 := new(v3.Mock)
		policyID := int64(2)

		policy, version := expectCreatePolicy(clientV3, policyID, 123, nil, "test policy description")
		expectReadPolicy(clientV3, policy, version, 3)
		// custom import mocks
		listPoliciesV2Resp := []cloudlets.Policy{
			{
				Name: "other-name",
			},
		}
		clientV2.On("ListPolicies", testutils.MockContext, cloudlets.ListPoliciesRequest{
			PageSize: ptr.To(1000),
			Offset:   0,
		}).Return(listPoliciesV2Resp, nil).Once()
		// mock that 1000 policies are returned, desired one not found
		clientV3.On("ListPolicies", testutils.MockContext, v3.ListPoliciesRequest{
			Size: 1000, Page: 0,
		}).Return(&v3.ListPoliciesResponse{
			Content: make([]v3.Policy, 1000),
			Page: v3.Page{
				Number:        0,
				Size:          1000,
				TotalElements: 1001,
				TotalPages:    2,
			},
		}, nil).Once()
		// mock that desired policy is on the next page
		clientV3.On("ListPolicies", testutils.MockContext, v3.ListPoliciesRequest{
			Size: 1000, Page: 1,
		}).Return(&v3.ListPoliciesResponse{
			Content: []v3.Policy{
				{
					Name: "test_policy",
					ID:   policyID,
				},
			},
			Page: v3.Page{
				Number:        1,
				Size:          1000,
				TotalElements: 1001,
				TotalPages:    2,
			},
		}, nil).Once()
		expectRemovePolicy(clientV3, policy.ID)

		useClientV2AndV3(clientV2, clientV3, func() {
			resource.UnitTest(t, resource.TestCase{
				ProtoV6ProviderFactories: testutils.NewProtoV6ProviderFactory(NewSubprovider()),
				Steps: []resource.TestStep{
					{
						Config: testutils.LoadFixtureStringf(t, "%s/create_no_match_rules/policy_create.tf", testDir),
						Check: checkPolicyAttributes(policyAttributes{
							groupID:     123,
							version:     "1",
							description: "test policy description",
						}),
					},
					{
						ImportState:             true,
						ImportStateId:           "test_policy",
						ResourceName:            "akamai_cloudlets_policy.policy",
						ImportStateVerify:       true,
						ImportStateVerifyIgnore: []string{"warnings"},
					},
				},
			})
		})
		clientV2.AssertExpectations(t)
		clientV3.AssertExpectations(t)
	})
}

func TestDiffSuppressMatchRules(t *testing.T) {
	basePath := "testdata/TestResPolicy/diff_suppress"
	tests := map[string]struct {
		oldPath, newPath string
		expected         bool
	}{
		"identical JSON": {
			oldPath:  "rules.json",
			newPath:  "rules.json",
			expected: true,
		},
		"different formatting, same content": {
			oldPath:  "rules.json",
			newPath:  "different_format.json",
			expected: true,
		},
		"difference in location and akaRuleId": {
			oldPath:  "with_location.json",
			newPath:  "rules.json",
			expected: true,
		},
		"invalid 'old' json": {
			oldPath:  "invalid.json",
			newPath:  "rules.json",
			expected: false,
		},
		"invalid 'new' json": {
			oldPath:  "rules.json",
			newPath:  "invalid.json",
			expected: false,
		},
	}

	for name, test := range tests {
		t.Run(name, func(t *testing.T) {
			oldJSON := testutils.LoadFixtureStringf(t, "%s/%s", basePath, test.oldPath)
			newJSON := testutils.LoadFixtureStringf(t, "%s/%s", basePath, test.newPath)
			res := diffSuppressMatchRules("", oldJSON, newJSON, nil)
			assert.Equal(t, test.expected, res)
		})
	}
}
