package cloudlets

import (
	"fmt"
	"regexp"
	"testing"

	"github.com/akamai/AkamaiOPEN-edgegrid-golang/v2/pkg/cloudlets"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/resource"
	"github.com/jinzhu/copier"
	"github.com/stretchr/testify/mock"
	"github.com/stretchr/testify/require"
	"github.com/tj/assert"
)

func expectCreatePolicy(_ *testing.T, client *mockcloudlets, policyID int64, policyName string, matchRules cloudlets.MatchRules) (*cloudlets.Policy, *cloudlets.PolicyVersion) {
	policy := &cloudlets.Policy{
		PolicyID:     policyID,
		GroupID:      123,
		Name:         policyName,
		Description:  "test policy description",
		CloudletID:   0,
		CloudletCode: "ER",
	}
	version := &cloudlets.PolicyVersion{
		Location:        "/version/1",
		PolicyID:        policyID,
		Version:         1,
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
	client.On("CreatePolicy", mock.Anything, cloudlets.CreatePolicyRequest{
		Name:        policyName,
		CloudletID:  0,
		Description: "test policy description",
		GroupID:     123,
	}).Return(policy, nil).Once()
	if matchRules == nil {
		return policy, version
	}
	client.On("UpdatePolicyVersion", mock.Anything, cloudlets.UpdatePolicyVersionRequest{
		UpdatePolicyVersion: cloudlets.UpdatePolicyVersion{
			MatchRules: matchRules,
		},
		PolicyID: policyID,
		Version:  1,
	}).Return(version, nil).Once()
	return policy, version
}

func expectReadPolicy(_ *testing.T, client *mockcloudlets, policy *cloudlets.Policy, version *cloudlets.PolicyVersion, times int) {
	client.On("GetPolicy", mock.Anything, policy.PolicyID).Return(policy, nil).Times(times)
	client.On("GetPolicyVersion", mock.Anything, cloudlets.GetPolicyVersionRequest{
		PolicyID: policy.PolicyID,
		Version:  version.Version,
	}).Return(version, nil).Times(times)
}

func expectUpdatePolicy(t *testing.T, client *mockcloudlets, policy *cloudlets.Policy, updatedName string) *cloudlets.Policy {
	var policyUpdate cloudlets.Policy
	err := copier.CopyWithOption(&policyUpdate, policy, copier.Option{DeepCopy: true})
	require.NoError(t, err)
	policyUpdate.Name = updatedName
	client.On("UpdatePolicy", mock.Anything, cloudlets.UpdatePolicyRequest{
		UpdatePolicy: cloudlets.UpdatePolicy{
			Name:        updatedName,
			Description: "test policy description",
			GroupID:     123,
		},
		PolicyID: policyUpdate.PolicyID,
	}).Return(&policyUpdate, nil).Once()
	return &policyUpdate
}

func expectCreatePolicyVersion(t *testing.T, client *mockcloudlets, policyID int64, version *cloudlets.PolicyVersion, newMatchRules cloudlets.MatchRules) *cloudlets.PolicyVersion {
	var activatedVersion cloudlets.PolicyVersion
	err := copier.CopyWithOption(&activatedVersion, version, copier.Option{DeepCopy: true})
	require.NoError(t, err)
	activatedVersion.Activations = []*cloudlets.Activation{{Network: "PROD"}}
	client.On("GetPolicyVersion", mock.Anything, cloudlets.GetPolicyVersionRequest{
		PolicyID:  policyID,
		Version:   version.Version,
		OmitRules: true,
	}).Return(&activatedVersion, nil).Once()
	var versionUpdate cloudlets.PolicyVersion
	err = copier.CopyWithOption(&versionUpdate, activatedVersion, copier.Option{DeepCopy: true})
	require.NoError(t, err)
	versionUpdate.MatchRules = newMatchRules
	versionUpdate.Version = version.Version + 1
	client.On("CreatePolicyVersion", mock.Anything, cloudlets.CreatePolicyVersionRequest{
		CreatePolicyVersion: cloudlets.CreatePolicyVersion{
			MatchRuleFormat: "1.0",
			MatchRules:      newMatchRules,
		},
		PolicyID: policyID,
	}).Return(&versionUpdate, nil).Once()
	return &versionUpdate
}

func expectUpdatePolicyVersion(t *testing.T, client *mockcloudlets, policyID int64, version *cloudlets.PolicyVersion, newMatchRules cloudlets.MatchRules) *cloudlets.PolicyVersion {
	client.On("GetPolicyVersion", mock.Anything, cloudlets.GetPolicyVersionRequest{
		PolicyID:  policyID,
		Version:   version.Version,
		OmitRules: true,
	}).Return(version, nil).Once()
	var versionUpdate cloudlets.PolicyVersion
	err := copier.CopyWithOption(&versionUpdate, version, copier.Option{DeepCopy: true})
	require.NoError(t, err)
	versionUpdate.MatchRules = newMatchRules
	client.On("UpdatePolicyVersion", mock.Anything, cloudlets.UpdatePolicyVersionRequest{
		UpdatePolicyVersion: cloudlets.UpdatePolicyVersion{
			MatchRuleFormat: "1.0",
			MatchRules:      newMatchRules,
		},
		PolicyID: policyID,
		Version:  version.Version,
	}).Return(&versionUpdate, nil).Once()
	return &versionUpdate
}

func expectRemovePolicy(_ *testing.T, client *mockcloudlets, policyID int64, numVersions int) {
	var versionList []cloudlets.PolicyVersion
	for i := 1; i <= numVersions; i++ {
		versionList = append(versionList, cloudlets.PolicyVersion{PolicyID: policyID, Version: int64(i)})
	}
	client.On("ListPolicyVersions", mock.Anything, cloudlets.ListPolicyVersionsRequest{
		PolicyID: policyID,
	}).Return(versionList, nil).Once()
	for _, ver := range versionList {
		client.On("DeletePolicyVersion", mock.Anything, cloudlets.DeletePolicyVersionRequest{
			PolicyID: ver.PolicyID,
			Version:  ver.Version,
		}).Return(nil).Once()
	}
	client.On("RemovePolicy", mock.Anything, policyID).Return(nil).Once()
}

func TestResourcePolicy(t *testing.T) {

	type policyAttributes struct {
		name, version, matchRulesPath string
	}

	checkPolicyAttributes := func(attrs policyAttributes) resource.TestCheckFunc {
		var matchRulesPath string
		if attrs.matchRulesPath != "" {
			matchRulesPath = loadFixtureString(attrs.matchRulesPath)
		}
		return resource.ComposeAggregateTestCheckFunc(
			resource.TestCheckResourceAttr("akamai_cloudlets_policy.policy", "id", "2"),
			resource.TestCheckResourceAttr("akamai_cloudlets_policy.policy", "cloudlet_code", "ER"),
			resource.TestCheckResourceAttr("akamai_cloudlets_policy.policy", "cloudlet_id", "0"),
			resource.TestCheckResourceAttr("akamai_cloudlets_policy.policy", "group_id", "123"),
			resource.TestCheckResourceAttr("akamai_cloudlets_policy.policy", "description", "test policy description"),
			resource.TestCheckResourceAttr("akamai_cloudlets_policy.policy", "match_rule_format", "1.0"),
			resource.TestCheckResourceAttr("akamai_cloudlets_policy.policy", "name", attrs.name),
			resource.TestCheckResourceAttr("akamai_cloudlets_policy.policy", "version", attrs.version),
			resource.TestCheckResourceAttr("akamai_cloudlets_policy.policy", "match_rules", matchRulesPath),
		)
	}

	t.Run("policy lifecycle with create new version", func(t *testing.T) {
		testDir := "testdata/TestResPolicy/lifecycle"

		client := new(mockcloudlets)
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
		policy, version := expectCreatePolicy(t, client, 2, "test_policy", matchRules)
		expectReadPolicy(t, client, policy, version, 3)
		policy = expectUpdatePolicy(t, client, policy, "test_policy_updated")
		version = expectCreatePolicyVersion(t, client, policy.PolicyID, version, matchRules[:1])
		expectReadPolicy(t, client, policy, version, 2)
		expectRemovePolicy(t, client, policy.PolicyID, 2)

		useClient(client, func() {
			resource.UnitTest(t, resource.TestCase{
				Providers: testAccProviders,
				Steps: []resource.TestStep{
					{
						Config: loadFixtureString(fmt.Sprintf("%s/policy_create.tf", testDir)),
						Check: checkPolicyAttributes(policyAttributes{
							name:           "test_policy",
							version:        "1",
							matchRulesPath: fmt.Sprintf("%s/match_rules/match_rules_create.json", testDir),
						}),
					},
					{
						Config: loadFixtureString(fmt.Sprintf("%s/policy_update.tf", testDir)),
						Check: checkPolicyAttributes(policyAttributes{
							name:           "test_policy_updated",
							version:        "2",
							matchRulesPath: fmt.Sprintf("%s/match_rules/match_rules_update.json", testDir),
						}),
					},
				},
			})
		})
		client.AssertExpectations(t)
	})

	t.Run("policy lifecycle with update existing version", func(t *testing.T) {
		testDir := "testdata/TestResPolicy/lifecycle"

		client := new(mockcloudlets)
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
		policy, version := expectCreatePolicy(t, client, 2, "test_policy", matchRules)
		expectReadPolicy(t, client, policy, version, 3)
		policy = expectUpdatePolicy(t, client, policy, "test_policy_updated")
		version = expectUpdatePolicyVersion(t, client, policy.PolicyID, version, matchRules[:1])
		expectReadPolicy(t, client, policy, version, 2)
		expectRemovePolicy(t, client, policy.PolicyID, 1)

		useClient(client, func() {
			resource.UnitTest(t, resource.TestCase{
				Providers: testAccProviders,
				Steps: []resource.TestStep{
					{
						Config: loadFixtureString(fmt.Sprintf("%s/policy_create.tf", testDir)),
						Check: checkPolicyAttributes(policyAttributes{
							name:           "test_policy",
							version:        "1",
							matchRulesPath: fmt.Sprintf("%s/match_rules/match_rules_create.json", testDir),
						}),
					},
					{
						Config: loadFixtureString(fmt.Sprintf("%s/policy_update.tf", testDir)),
						Check: checkPolicyAttributes(policyAttributes{
							name:           "test_policy_updated",
							version:        "1",
							matchRulesPath: fmt.Sprintf("%s/match_rules/match_rules_update.json", testDir),
						}),
					},
				},
			})
		})
		client.AssertExpectations(t)
	})

	t.Run("update only policy", func(t *testing.T) {
		testDir := "testdata/TestResPolicy/lifecycle_policy_update"

		client := new(mockcloudlets)
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
		policy, version := expectCreatePolicy(t, client, 2, "test_policy", matchRules)
		expectReadPolicy(t, client, policy, version, 3)
		policy = expectUpdatePolicy(t, client, policy, "test_policy_updated")
		expectReadPolicy(t, client, policy, version, 2)
		expectRemovePolicy(t, client, policy.PolicyID, 1)

		useClient(client, func() {
			resource.UnitTest(t, resource.TestCase{
				Providers: testAccProviders,
				Steps: []resource.TestStep{
					{
						Config: loadFixtureString(fmt.Sprintf("%s/policy_create.tf", testDir)),
						Check: checkPolicyAttributes(policyAttributes{
							name:           "test_policy",
							version:        "1",
							matchRulesPath: fmt.Sprintf("%s/match_rules/match_rules.json", testDir),
						}),
					},
					{
						Config: loadFixtureString(fmt.Sprintf("%s/policy_update.tf", testDir)),
						Check: checkPolicyAttributes(policyAttributes{
							name:           "test_policy_updated",
							version:        "1",
							matchRulesPath: fmt.Sprintf("%s/match_rules/match_rules.json", testDir),
						}),
					},
				},
			})
		})
		client.AssertExpectations(t)
	})

	t.Run("update only version", func(t *testing.T) {
		testDir := "testdata/TestResPolicy/lifecycle_version_update"

		client := new(mockcloudlets)
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
		policy, version := expectCreatePolicy(t, client, 2, "test_policy", matchRules)
		expectReadPolicy(t, client, policy, version, 3)
		version = expectUpdatePolicyVersion(t, client, policy.PolicyID, version, matchRules[:1])
		expectReadPolicy(t, client, policy, version, 2)
		expectRemovePolicy(t, client, policy.PolicyID, 1)

		useClient(client, func() {
			resource.UnitTest(t, resource.TestCase{
				Providers: testAccProviders,
				Steps: []resource.TestStep{
					{
						Config: loadFixtureString(fmt.Sprintf("%s/policy_create.tf", testDir)),
						Check: checkPolicyAttributes(policyAttributes{
							name:           "test_policy",
							version:        "1",
							matchRulesPath: fmt.Sprintf("%s/match_rules/match_rules_create.json", testDir),
						}),
					},
					{
						Config: loadFixtureString(fmt.Sprintf("%s/policy_update.tf", testDir)),
						Check: checkPolicyAttributes(policyAttributes{
							name:           "test_policy",
							version:        "1",
							matchRulesPath: fmt.Sprintf("%s/match_rules/match_rules_update.json", testDir),
						}),
					},
				},
			})
		})
		client.AssertExpectations(t)
	})

	t.Run("remove match rules from version", func(t *testing.T) {
		testDir := "testdata/TestResPolicy/lifecycle_remove_match_rules"

		client := new(mockcloudlets)
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
		policy, version := expectCreatePolicy(t, client, 2, "test_policy", matchRules)
		expectReadPolicy(t, client, policy, version, 3)
		version = expectUpdatePolicyVersion(t, client, policy.PolicyID, version, cloudlets.MatchRules{})
		expectReadPolicy(t, client, policy, version, 2)
		expectRemovePolicy(t, client, policy.PolicyID, 1)

		useClient(client, func() {
			resource.UnitTest(t, resource.TestCase{
				Providers: testAccProviders,
				Steps: []resource.TestStep{
					{
						Config: loadFixtureString(fmt.Sprintf("%s/policy_create.tf", testDir)),
						Check: checkPolicyAttributes(policyAttributes{
							name:           "test_policy",
							version:        "1",
							matchRulesPath: fmt.Sprintf("%s/match_rules/match_rules_create.json", testDir),
						}),
					},
					{
						Config: loadFixtureString(fmt.Sprintf("%s/policy_update.tf", testDir)),
						Check: checkPolicyAttributes(policyAttributes{
							name:           "test_policy",
							version:        "1",
							matchRulesPath: "",
						}),
					},
				},
			})
		})
		client.AssertExpectations(t)
	})

	t.Run("create policy without match rules", func(t *testing.T) {
		testDir := "testdata/TestResPolicy/create_no_match_rules"

		client := new(mockcloudlets)
		policy, version := expectCreatePolicy(t, client, 2, "test_policy", nil)
		expectReadPolicy(t, client, policy, version, 2)
		expectRemovePolicy(t, client, 2, 1)
		useClient(client, func() {
			resource.UnitTest(t, resource.TestCase{
				Providers: testAccProviders,
				Steps: []resource.TestStep{
					{
						Config: loadFixtureString(fmt.Sprintf("%s/policy_create.tf", testDir)),
						Check: checkPolicyAttributes(policyAttributes{
							name:           "test_policy",
							version:        "1",
							matchRulesPath: "",
						}),
					},
				},
			})
		})
		client.AssertExpectations(t)
	})

	t.Run("error creating policy", func(t *testing.T) {
		testDir := "testdata/TestResPolicy/lifecycle"

		client := new(mockcloudlets)
		client.On("CreatePolicy", mock.Anything, cloudlets.CreatePolicyRequest{
			Name:        "test_policy",
			CloudletID:  0,
			GroupID:     123,
			Description: "test policy description",
		}).Return(nil, fmt.Errorf("oops"))

		useClient(client, func() {
			resource.UnitTest(t, resource.TestCase{
				Providers: testAccProviders,
				Steps: []resource.TestStep{
					{
						Config:      loadFixtureString(fmt.Sprintf("%s/policy_create.tf", testDir)),
						ExpectError: regexp.MustCompile("oops"),
					},
				},
			})
		})
		client.AssertExpectations(t)
	})

	t.Run("error creating policy version", func(t *testing.T) {
		testDir := "testdata/TestResPolicy/lifecycle"

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
		policy := &cloudlets.Policy{
			PolicyID:     2,
			GroupID:      123,
			Name:         "test_policy",
			Description:  "test policy description",
			CloudletID:   0,
			CloudletCode: "ER",
		}
		client := new(mockcloudlets)
		client.On("CreatePolicy", mock.Anything, cloudlets.CreatePolicyRequest{
			Name:        "test_policy",
			CloudletID:  0,
			GroupID:     123,
			Description: "test policy description",
		}).Return(policy, nil)
		client.On("UpdatePolicyVersion", mock.Anything, cloudlets.UpdatePolicyVersionRequest{
			UpdatePolicyVersion: cloudlets.UpdatePolicyVersion{
				MatchRules: matchRules,
			},
			PolicyID: 2,
			Version:  1,
		}).Return(nil, fmt.Errorf("oops"))
		expectRemovePolicy(t, client, 2, 0)

		client.On("CreatePolicy", mock.Anything, cloudlets.CreatePolicyRequest{
			Name:        "test_policy",
			CloudletID:  0,
			GroupID:     123,
			Description: "test policy description",
		}).Return(nil, fmt.Errorf("oops"))

		useClient(client, func() {
			resource.UnitTest(t, resource.TestCase{
				Providers: testAccProviders,
				Steps: []resource.TestStep{
					{
						Config:      loadFixtureString(fmt.Sprintf("%s/policy_create.tf", testDir)),
						ExpectError: regexp.MustCompile("oops"),
					},
				},
			})
		})
		client.AssertExpectations(t)
	})

	t.Run("error fetching policy", func(t *testing.T) {
		testDir := "testdata/TestResPolicy/create_no_match_rules"

		client := new(mockcloudlets)
		policy, _ := expectCreatePolicy(t, client, 2, "test_policy", nil)
		client.On("GetPolicy", mock.Anything, policy.PolicyID).Return(nil, fmt.Errorf("oops"))
		expectRemovePolicy(t, client, 2, 1)
		useClient(client, func() {
			resource.UnitTest(t, resource.TestCase{
				Providers: testAccProviders,
				Steps: []resource.TestStep{
					{
						Config:      loadFixtureString(fmt.Sprintf("%s/policy_create.tf", testDir)),
						ExpectError: regexp.MustCompile("oops"),
					},
				},
			})
		})
		client.AssertExpectations(t)
	})

	t.Run("error fetching policy version", func(t *testing.T) {
		testDir := "testdata/TestResPolicy/create_no_match_rules"

		client := new(mockcloudlets)
		policy, version := expectCreatePolicy(t, client, 2, "test_policy", nil)
		client.On("GetPolicy", mock.Anything, policy.PolicyID).Return(policy, nil)
		client.On("GetPolicyVersion", mock.Anything, cloudlets.GetPolicyVersionRequest{
			PolicyID: policy.PolicyID,
			Version:  version.Version,
		}).Return(nil, fmt.Errorf("oops"))
		expectRemovePolicy(t, client, 2, 1)
		useClient(client, func() {
			resource.UnitTest(t, resource.TestCase{
				Providers: testAccProviders,
				Steps: []resource.TestStep{
					{
						Config:      loadFixtureString(fmt.Sprintf("%s/policy_create.tf", testDir)),
						ExpectError: regexp.MustCompile("oops"),
					},
				},
			})
		})
		client.AssertExpectations(t)
	})

	t.Run("error updating policy", func(t *testing.T) {
		testDir := "testdata/TestResPolicy/lifecycle_policy_update"

		client := new(mockcloudlets)
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
		policy, version := expectCreatePolicy(t, client, 2, "test_policy", matchRules)
		expectReadPolicy(t, client, policy, version, 3)
		client.On("UpdatePolicy", mock.Anything, cloudlets.UpdatePolicyRequest{
			UpdatePolicy: cloudlets.UpdatePolicy{
				Name:        "test_policy_updated",
				Description: "test policy description",
				GroupID:     123,
			},
			PolicyID: policy.PolicyID,
		}).Return(nil, fmt.Errorf("oops")).Once()
		expectRemovePolicy(t, client, policy.PolicyID, 1)

		useClient(client, func() {
			resource.UnitTest(t, resource.TestCase{
				Providers: testAccProviders,
				Steps: []resource.TestStep{
					{
						Config: loadFixtureString(fmt.Sprintf("%s/policy_create.tf", testDir)),
					},
					{
						Config:      loadFixtureString(fmt.Sprintf("%s/policy_update.tf", testDir)),
						ExpectError: regexp.MustCompile("oops"),
					},
				},
			})
		})
		client.AssertExpectations(t)
	})

	t.Run("error updating version", func(t *testing.T) {
		testDir := "testdata/TestResPolicy/lifecycle_version_update"

		client := new(mockcloudlets)
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
		policy, version := expectCreatePolicy(t, client, 2, "test_policy", matchRules)
		expectReadPolicy(t, client, policy, version, 3)
		client.On("GetPolicyVersion", mock.Anything, cloudlets.GetPolicyVersionRequest{
			PolicyID:  policy.PolicyID,
			Version:   version.Version,
			OmitRules: true,
		}).Return(version, nil).Once()
		client.On("UpdatePolicyVersion", mock.Anything, cloudlets.UpdatePolicyVersionRequest{
			UpdatePolicyVersion: cloudlets.UpdatePolicyVersion{
				MatchRuleFormat: "1.0",
				MatchRules:      matchRules[:1],
			},
			PolicyID: policy.PolicyID,
			Version:  version.Version,
		}).Return(nil, fmt.Errorf("oops")).Once()
		expectRemovePolicy(t, client, policy.PolicyID, 1)

		useClient(client, func() {
			resource.UnitTest(t, resource.TestCase{
				Providers: testAccProviders,
				Steps: []resource.TestStep{
					{
						Config: loadFixtureString(fmt.Sprintf("%s/policy_create.tf", testDir)),
					},
					{
						Config:      loadFixtureString(fmt.Sprintf("%s/policy_update.tf", testDir)),
						ExpectError: regexp.MustCompile("oops"),
					},
				},
			})
		})
		client.AssertExpectations(t)
	})

	t.Run("invalid group id passed", func(t *testing.T) {
		testDir := "testdata/TestResPolicy/invalid_group_id"
		client := new(mockcloudlets)

		useClient(client, func() {
			resource.UnitTest(t, resource.TestCase{
				Providers: testAccProviders,
				Steps: []resource.TestStep{
					{
						Config:      loadFixtureString(fmt.Sprintf("%s/policy_create.tf", testDir)),
						ExpectError: regexp.MustCompile("invalid group_id provided"),
					},
				},
			})
		})
		client.AssertExpectations(t)
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
			oldJSON := loadFixtureString(fmt.Sprintf("%s/%s", basePath, test.oldPath))
			newJSON := loadFixtureString(fmt.Sprintf("%s/%s", basePath, test.newPath))
			res := diffSuppressMatchRules("", oldJSON, newJSON, nil)
			assert.Equal(t, test.expected, res)
		})
	}
}
