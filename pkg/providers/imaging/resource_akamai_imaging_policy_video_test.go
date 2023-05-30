package imaging

import (
	"fmt"
	"net/http"
	"regexp"
	"testing"

	"github.com/akamai/AkamaiOPEN-edgegrid-golang/v6/pkg/imaging"
	"github.com/akamai/terraform-provider-akamai/v4/pkg/tools"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/resource"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/mock"
)

func TestResourcePolicyVideo(t *testing.T) {

	type policyAttributes struct {
		version              string
		policyPath           string
		policyID             string
		policySetID          string
		activateOnProduction string
		schema               bool
		emptyPolicy          bool
	}

	var (
		policyInput = imaging.PolicyInputVideo{
			Breakpoints: &imaging.Breakpoints{
				Widths: []int{320, 640, 1024, 2048, 5000},
			},
			Output: &imaging.OutputVideo{
				PerceptualQuality: &imaging.OutputVideoPerceptualQualityVariableInline{
					Value: imaging.OutputVideoPerceptualQualityPtr(imaging.OutputVideoPerceptualQualityMediumHigh),
				},
			},
		}
		policyOutput = imaging.PolicyOutputVideo{
			Breakpoints: &imaging.Breakpoints{
				Widths: []int{320, 640, 1024, 2048, 5000},
			},
			Output: &imaging.OutputVideo{
				PerceptualQuality: &imaging.OutputVideoPerceptualQualityVariableInline{
					Value: imaging.OutputVideoPerceptualQualityPtr(imaging.OutputVideoPerceptualQualityMediumHigh),
				},
			},
			Version: 1,
			Video:   tools.BoolPtr(true),
		}
		defaultBreakpointsWidths = &imaging.Breakpoints{
			Widths: []int{320, 640, 1024, 2048, 5000},
		}
		defaultHosts     = []string{"test1", "test2", "test3"}
		defaultVariables = []imaging.Variable{
			{
				DefaultValue: "test1",
				Name:         "var1",
				Type:         "string",
			},
			{
				DefaultValue: "test2",
				Name:         "var2",
				Type:         "string",
			},
			{
				DefaultValue: "test3",
				Name:         "var3",
				Type:         "string",
			},
		}
		policyInputDiff = imaging.PolicyInputVideo{
			Breakpoints: defaultBreakpointsWidths,
			Hosts:       defaultHosts,
			Variables:   defaultVariables,
		}
		policyOutputDiff = imaging.PolicyOutputVideo{
			Breakpoints: defaultBreakpointsWidths,
			Hosts:       defaultHosts,
			Variables:   defaultVariables,
			Version:     1,
			Video:       tools.BoolPtr(true),
		}

		expectUpsertPolicy = func(_ *testing.T, client *imaging.Mock, policyID string, network imaging.PolicyNetwork, contractID string, policySetID string, policy imaging.PolicyInput) {
			policyResponse := &imaging.PolicyResponse{
				OperationPerformed: "UPDATED",
				Description:        fmt.Sprintf("Policy %s updated.", policyID),
				ID:                 policyID,
			}
			client.On("UpsertPolicy", mock.Anything, imaging.UpsertPolicyRequest{
				PolicyID:    policyID,
				Network:     network,
				ContractID:  contractID,
				PolicySetID: policySetID,
				PolicyInput: policy,
			}).Return(policyResponse, nil).Once()
		}

		expectReadPolicy = func(t *testing.T, client *imaging.Mock, policyID string, network imaging.PolicyNetwork, contractID string, policySetID string, policyOutput imaging.PolicyOutput, times int) {
			client.On("GetPolicy", mock.Anything, imaging.GetPolicyRequest{
				PolicyID:    policyID,
				Network:     network,
				ContractID:  contractID,
				PolicySetID: policySetID,
			}).Return(policyOutput, nil).Times(times)
		}

		expectDeletePolicy = func(_ *testing.T, client *imaging.Mock, policyID string, network imaging.PolicyNetwork, contractID string, policySetID string) {
			response := imaging.PolicyResponse{}
			client.On("DeletePolicy", mock.Anything, imaging.DeletePolicyRequest{
				PolicyID:    policyID,
				Network:     network,
				ContractID:  contractID,
				PolicySetID: policySetID,
			}).Return(&response, nil).Once()
		}

		expectUpsertPolicyWithError = func(_ *testing.T, client *imaging.Mock, policyID string, network imaging.PolicyNetwork, contractID string, policySetID string, policy imaging.PolicyInput, err error) {
			client.On("UpsertPolicy", mock.Anything, imaging.UpsertPolicyRequest{
				PolicyID:    policyID,
				Network:     network,
				ContractID:  contractID,
				PolicySetID: policySetID,
				PolicyInput: policy,
			}).Return(nil, err).Once()
		}

		checkPolicyAttributes = func(attrs policyAttributes) resource.TestCheckFunc {
			var policyJSON string
			if attrs.policyPath != "" {
				policyJSON = loadFixtureString(attrs.policyPath)
			}
			f := resource.ComposeAggregateTestCheckFunc(
				resource.TestCheckResourceAttr("akamai_imaging_policy_video.policy", "id", fmt.Sprintf("%s:%s", attrs.policySetID, attrs.policyID)),
				resource.TestCheckResourceAttr("akamai_imaging_policy_video.policy", "policy_id", attrs.policyID),
				resource.TestCheckResourceAttr("akamai_imaging_policy_video.policy", "policyset_id", attrs.policySetID),
				resource.TestCheckResourceAttr("akamai_imaging_policy_video.policy", "contract_id", "test_contract"),
				resource.TestCheckResourceAttr("akamai_imaging_policy_video.policy", "version", attrs.version),
				resource.TestCheckResourceAttr("akamai_imaging_policy_video.policy", "activate_on_production", attrs.activateOnProduction),
			)
			if attrs.schema {
				if attrs.emptyPolicy {
					return resource.ComposeAggregateTestCheckFunc(
						f,
						resource.TestCheckResourceAttr("akamai_imaging_policy_video.policy", "policy.#", "1"),
					)
				}
				return resource.ComposeAggregateTestCheckFunc(
					f,
					resource.TestCheckResourceAttr("akamai_imaging_policy_video.policy", "policy.#", "1"),
					resource.TestCheckResourceAttr("akamai_imaging_policy_video.policy", "policy.0.output.0.perceptual_quality", "mediumHigh"),
				)
			}
			return resource.ComposeAggregateTestCheckFunc(
				f,
				resource.TestCheckResourceAttr("akamai_imaging_policy_video.policy", "json", policyJSON),
			)
		}
	)

	t.Run("regular policy create", func(t *testing.T) {
		testDir := "testdata/TestResPolicyVideo/regular_policy"

		client := new(imaging.Mock)
		expectUpsertPolicy(t, client, "test_policy", imaging.PolicyNetworkStaging, "test_contract", "test_policy_set", &policyInput)
		expectReadPolicy(t, client, "test_policy", imaging.PolicyNetworkStaging, "test_contract", "test_policy_set", &policyOutput, 2)

		expectDeletePolicy(t, client, "test_policy", imaging.PolicyNetworkStaging, "test_contract", "test_policy_set")
		// it is faster to attempt to delete on production than checking if there is policy on production first
		expectDeletePolicy(t, client, "test_policy", imaging.PolicyNetworkProduction, "test_contract", "test_policy_set")

		useClient(client, func() {
			resource.UnitTest(t, resource.TestCase{
				ProviderFactories: testAccProviders,
				Steps: []resource.TestStep{
					{
						Config: loadFixtureString(fmt.Sprintf("%s/policy_create.tf", testDir)),
						Check: checkPolicyAttributes(policyAttributes{
							version:              "1",
							policyID:             "test_policy",
							policySetID:          "test_policy_set",
							activateOnProduction: "false",
							policyPath:           fmt.Sprintf("%s/policy/policy_create.json", testDir),
						}),
					},
				},
			})
		})
		client.AssertExpectations(t)
	})
	t.Run("regular policy create and later activate on production", func(t *testing.T) {
		testDir := "testdata/TestResPolicyVideo/regular_policy"

		client := new(imaging.Mock)
		expectUpsertPolicy(t, client, "test_policy", imaging.PolicyNetworkStaging, "test_contract", "test_policy_set", &policyInput)
		expectReadPolicy(t, client, "test_policy", imaging.PolicyNetworkStaging, "test_contract", "test_policy_set", &policyOutput, 5)

		// `activate_on_production` should not trigger Upsert for staging if the policy has not changed
		expectUpsertPolicy(t, client, "test_policy", imaging.PolicyNetworkProduction, "test_contract", "test_policy_set", &policyInput)

		expectDeletePolicy(t, client, "test_policy", imaging.PolicyNetworkStaging, "test_contract", "test_policy_set")
		expectDeletePolicy(t, client, "test_policy", imaging.PolicyNetworkProduction, "test_contract", "test_policy_set")

		useClient(client, func() {
			resource.UnitTest(t, resource.TestCase{
				ProviderFactories: testAccProviders,
				Steps: []resource.TestStep{
					{
						Config: loadFixtureString(fmt.Sprintf("%s/policy_create.tf", testDir)),
						Check: checkPolicyAttributes(policyAttributes{
							version:              "1",
							policyID:             "test_policy",
							policySetID:          "test_policy_set",
							activateOnProduction: "false",
							policyPath:           fmt.Sprintf("%s/policy/policy_create.json", testDir),
						}),
					},
					{
						Config: loadFixtureString(fmt.Sprintf("%s/policy_update.tf", testDir)),
						Check: checkPolicyAttributes(policyAttributes{
							version:              "1",
							policyID:             "test_policy",
							policySetID:          "test_policy_set",
							activateOnProduction: "true",
							policyPath:           fmt.Sprintf("%s/policy/policy_update.json", testDir),
						}),
					},
				},
			})
		})
		client.AssertExpectations(t)
	})
	t.Run("regular policy create and activate on production, later update both", func(t *testing.T) {
		testDir := "testdata/TestResPolicyVideo/regular_policy_activate_same_time"

		client := new(imaging.Mock)
		expectReadPolicy(t, client, "test_policy", imaging.PolicyNetworkStaging, "test_contract", "test_policy_set", &policyOutput, 2)

		expectUpsertPolicy(t, client, "test_policy", imaging.PolicyNetworkStaging, "test_contract", "test_policy_set", &policyInput)
		expectUpsertPolicy(t, client, "test_policy", imaging.PolicyNetworkProduction, "test_contract", "test_policy_set", &policyInput)

		// update
		expectReadPolicy(t, client, "test_policy", imaging.PolicyNetworkStaging, "test_contract", "test_policy_set", &policyOutput, 1)

		policyInputV2 := getPolicyInputVideoV2(policyInput)
		policyOutputV2 := getPolicyOutputVideoV2(policyOutput)

		expectReadPolicy(t, client, "test_policy", imaging.PolicyNetworkStaging, "test_contract", "test_policy_set", &policyOutputV2, 2)
		expectUpsertPolicy(t, client, "test_policy", imaging.PolicyNetworkStaging, "test_contract", "test_policy_set", &policyInputV2)
		expectUpsertPolicy(t, client, "test_policy", imaging.PolicyNetworkProduction, "test_contract", "test_policy_set", &policyInputV2)

		expectDeletePolicy(t, client, "test_policy", imaging.PolicyNetworkStaging, "test_contract", "test_policy_set")
		expectDeletePolicy(t, client, "test_policy", imaging.PolicyNetworkProduction, "test_contract", "test_policy_set")

		useClient(client, func() {
			resource.UnitTest(t, resource.TestCase{
				ProviderFactories: testAccProviders,
				Steps: []resource.TestStep{
					{
						Config: loadFixtureString(fmt.Sprintf("%s/policy_create.tf", testDir)),
						Check: checkPolicyAttributes(policyAttributes{
							version:              "1",
							policyID:             "test_policy",
							policySetID:          "test_policy_set",
							activateOnProduction: "true",
							policyPath:           fmt.Sprintf("%s/policy/policy_create.json", testDir),
						}),
					},
					{
						Config: loadFixtureString(fmt.Sprintf("%s/policy_update.tf", testDir)),
						Check: checkPolicyAttributes(policyAttributes{
							version:              "2",
							policyID:             "test_policy",
							policySetID:          "test_policy_set",
							activateOnProduction: "true",
							policyPath:           fmt.Sprintf("%s/policy/policy_update.json", testDir),
						}),
					},
				},
			})
		})
		client.AssertExpectations(t)
	})
	t.Run("regular policy create and later change policy set id (force new)", func(t *testing.T) {
		testDir := "testdata/TestResPolicyVideo/change_policyset_id"

		client := new(imaging.Mock)
		expectUpsertPolicy(t, client, "test_policy", imaging.PolicyNetworkStaging, "test_contract", "test_policy_set", &policyInput)
		expectReadPolicy(t, client, "test_policy", imaging.PolicyNetworkStaging, "test_contract", "test_policy_set", &policyOutput, 3)

		// remove original policy
		expectDeletePolicy(t, client, "test_policy", imaging.PolicyNetworkStaging, "test_contract", "test_policy_set")
		expectDeletePolicy(t, client, "test_policy", imaging.PolicyNetworkProduction, "test_contract", "test_policy_set")

		// update
		expectUpsertPolicy(t, client, "test_policy", imaging.PolicyNetworkStaging, "test_contract", "test_policy_set_update", &policyInput)
		expectReadPolicy(t, client, "test_policy", imaging.PolicyNetworkStaging, "test_contract", "test_policy_set_update", &policyOutput, 2)

		// remove new policy
		expectDeletePolicy(t, client, "test_policy", imaging.PolicyNetworkStaging, "test_contract", "test_policy_set_update")
		expectDeletePolicy(t, client, "test_policy", imaging.PolicyNetworkProduction, "test_contract", "test_policy_set_update")

		useClient(client, func() {
			resource.UnitTest(t, resource.TestCase{
				ProviderFactories: testAccProviders,
				Steps: []resource.TestStep{
					{
						Config: loadFixtureString(fmt.Sprintf("%s/policy_create.tf", testDir)),
						Check: checkPolicyAttributes(policyAttributes{
							version:              "1",
							policyID:             "test_policy",
							policySetID:          "test_policy_set",
							activateOnProduction: "false",
							policyPath:           fmt.Sprintf("%s/policy/policy_create.json", testDir),
						}),
					},
					{
						Config: loadFixtureString(fmt.Sprintf("%s/policy_update.tf", testDir)),
						Check: checkPolicyAttributes(policyAttributes{
							version:              "1",
							policyID:             "test_policy",
							policySetID:          "test_policy_set_update",
							activateOnProduction: "false",
							policyPath:           fmt.Sprintf("%s/policy/policy_update.json", testDir),
						}),
					},
				},
			})
		})
		client.AssertExpectations(t)
	})
	t.Run("regular policy create, later activate on production and later modify on staging only", func(t *testing.T) {
		testDir := "testdata/TestResPolicyVideo/regular_policy_update_staging"

		client := new(imaging.Mock)
		expectUpsertPolicy(t, client, "test_policy", imaging.PolicyNetworkStaging, "test_contract", "test_policy_set", &policyInput)
		expectReadPolicy(t, client, "test_policy", imaging.PolicyNetworkStaging, "test_contract", "test_policy_set", &policyOutput, 6)

		// `activate_on_production` should not trigger Upsert for staging if the policy has not changed
		expectUpsertPolicy(t, client, "test_policy", imaging.PolicyNetworkProduction, "test_contract", "test_policy_set", &policyInput)

		policyInputV2 := getPolicyInputVideoV2(policyInput)
		policyOutputV2 := getPolicyOutputVideoV2(policyOutput)

		expectReadPolicy(t, client, "test_policy", imaging.PolicyNetworkStaging, "test_contract", "test_policy_set", &policyOutputV2, 2)
		expectUpsertPolicy(t, client, "test_policy", imaging.PolicyNetworkStaging, "test_contract", "test_policy_set", &policyInputV2)

		expectDeletePolicy(t, client, "test_policy", imaging.PolicyNetworkStaging, "test_contract", "test_policy_set")
		expectDeletePolicy(t, client, "test_policy", imaging.PolicyNetworkProduction, "test_contract", "test_policy_set")

		useClient(client, func() {
			resource.UnitTest(t, resource.TestCase{
				ProviderFactories: testAccProviders,
				Steps: []resource.TestStep{
					{
						Config: loadFixtureString(fmt.Sprintf("%s/policy_create.tf", testDir)),
						Check: checkPolicyAttributes(policyAttributes{
							version:              "1",
							policyID:             "test_policy",
							policySetID:          "test_policy_set",
							activateOnProduction: "false",
							policyPath:           fmt.Sprintf("%s/policy/policy_create.json", testDir),
						}),
					},
					{
						Config: loadFixtureString(fmt.Sprintf("%s/policy_update.tf", testDir)),
						Check: checkPolicyAttributes(policyAttributes{
							version:              "1",
							policyID:             "test_policy",
							policySetID:          "test_policy_set",
							activateOnProduction: "true",
							policyPath:           fmt.Sprintf("%s/policy/policy_update.json", testDir),
						}),
					},
					{
						Config: loadFixtureString(fmt.Sprintf("%s/policy_update_staging.tf", testDir)),
						Check: checkPolicyAttributes(policyAttributes{
							version:              "2",
							policyID:             "test_policy",
							policySetID:          "test_policy_set",
							activateOnProduction: "false",
							policyPath:           fmt.Sprintf("%s/policy/policy_update_staging.json", testDir),
						}),
					},
				},
			})
		})
		client.AssertExpectations(t)
	})
	t.Run("auto policy create and later activate on production, cannot delete", func(t *testing.T) {
		testDir := "testdata/TestResPolicyVideo/auto_policy"

		client := new(imaging.Mock)
		expectUpsertPolicy(t, client, ".auto", imaging.PolicyNetworkStaging, "test_contract", "test_policy_set", &policyInput)
		expectReadPolicy(t, client, ".auto", imaging.PolicyNetworkStaging, "test_contract", "test_policy_set", &policyOutput, 5)

		// `activate_on_production` should not trigger Upsert for staging if the policy has not changed
		expectUpsertPolicy(t, client, ".auto", imaging.PolicyNetworkProduction, "test_contract", "test_policy_set", &policyInput)

		// .auto policy cannot be removed alone, only via removal of policy set

		useClient(client, func() {
			resource.UnitTest(t, resource.TestCase{
				ProviderFactories: testAccProviders,
				Steps: []resource.TestStep{
					{
						Config: loadFixtureString(fmt.Sprintf("%s/policy_create.tf", testDir)),
						Check: checkPolicyAttributes(policyAttributes{
							version:              "1",
							policyID:             ".auto",
							policySetID:          "test_policy_set",
							activateOnProduction: "false",
							policyPath:           fmt.Sprintf("%s/policy/policy_create.json", testDir),
						}),
					},
					{
						Config: loadFixtureString(fmt.Sprintf("%s/policy_update.tf", testDir)),
						Check: checkPolicyAttributes(policyAttributes{
							version:              "1",
							policyID:             ".auto",
							policySetID:          "test_policy_set",
							activateOnProduction: "true",
							policyPath:           fmt.Sprintf("%s/policy/policy_update.json", testDir),
						}),
					},
				},
			})
		})
		client.AssertExpectations(t)
	})
	t.Run("import policy with activate_on_production=true", func(t *testing.T) {
		testDir := "testdata/TestResPolicyVideo/regular_policy"

		client := new(imaging.Mock)
		expectUpsertPolicy(t, client, "test_policy", imaging.PolicyNetworkStaging, "test_contract", "test_policy_set", &policyInput)
		expectUpsertPolicy(t, client, "test_policy", imaging.PolicyNetworkProduction, "test_contract", "test_policy_set", &policyInput)
		expectReadPolicy(t, client, "test_policy", imaging.PolicyNetworkStaging, "test_contract", "test_policy_set", &policyOutput, 5)

		expectReadPolicy(t, client, "test_policy", imaging.PolicyNetworkProduction, "test_contract", "test_policy_set", &policyOutput, 1)
		expectReadPolicy(t, client, "test_policy", imaging.PolicyNetworkStaging, "test_contract", "test_policy_set", &policyOutput, 2)

		expectDeletePolicy(t, client, "test_policy", imaging.PolicyNetworkStaging, "test_contract", "test_policy_set")
		expectDeletePolicy(t, client, "test_policy", imaging.PolicyNetworkProduction, "test_contract", "test_policy_set")

		useClient(client, func() {
			resource.UnitTest(t, resource.TestCase{
				ProviderFactories: testAccProviders,
				Steps: []resource.TestStep{
					{
						Config: loadFixtureString(fmt.Sprintf("%s/policy_create.tf", testDir)),
					},
					{
						Config: loadFixtureString(fmt.Sprintf("%s/policy_update.tf", testDir)),
					},
					{
						ImportState:       true,
						ImportStateId:     "test_policy:test_policy_set:test_contract",
						ResourceName:      "akamai_imaging_policy_video.policy",
						ImportStateVerify: true,
					},
				},
			})
		})
		client.AssertExpectations(t)
	})
	t.Run("regular policy create, check diff for breakpoints.Widths, hosts, variables", func(t *testing.T) {
		testDir := "testdata/TestResPolicyVideo/diff_suppress/fields"

		client := new(imaging.Mock)
		expectUpsertPolicy(t, client, "test_policy", imaging.PolicyNetworkStaging, "test_contract", "test_policy_set", &policyInputDiff)
		expectReadPolicy(t, client, "test_policy", imaging.PolicyNetworkStaging, "test_contract", "test_policy_set", &policyOutputDiff, 2)

		// remove original policy
		expectDeletePolicy(t, client, "test_policy", imaging.PolicyNetworkStaging, "test_contract", "test_policy_set")
		expectDeletePolicy(t, client, "test_policy", imaging.PolicyNetworkProduction, "test_contract", "test_policy_set")

		expectReadPolicy(t, client, "test_policy", imaging.PolicyNetworkStaging, "test_contract", "test_policy_set", &policyOutputDiff, 2)

		useClient(client, func() {
			resource.UnitTest(t, resource.TestCase{
				ProviderFactories: testAccProviders,
				Steps: []resource.TestStep{
					{
						Config: loadFixtureString(fmt.Sprintf("%s/default.tf", testDir)),
						Check: checkPolicyAttributes(policyAttributes{
							version:              "1",
							policyID:             "test_policy",
							policySetID:          "test_policy_set",
							activateOnProduction: "false",
							policyPath:           fmt.Sprintf("%s/policy.json", testDir),
						}),
					},
					{
						Config:             loadFixtureString(fmt.Sprintf("%s/diff_order.tf", testDir)),
						ExpectNonEmptyPlan: false,
					},
				},
			})
		})
		client.AssertExpectations(t)
	})
	t.Run("import policy with activate_on_production=false", func(t *testing.T) {
		testDir := "testdata/TestResPolicyVideo/regular_policy"
		policyOutputV2 := getPolicyOutputVideoV2(policyOutput)

		client := new(imaging.Mock)
		expectUpsertPolicy(t, client, "test_policy", imaging.PolicyNetworkStaging, "test_contract", "test_policy_set", &policyInput)
		expectReadPolicy(t, client, "test_policy", imaging.PolicyNetworkStaging, "test_contract", "test_policy_set", &policyOutput, 3)
		expectReadPolicy(t, client, "test_policy", imaging.PolicyNetworkProduction, "test_contract", "test_policy_set", &policyOutputV2, 1)
		expectReadPolicy(t, client, "test_policy", imaging.PolicyNetworkStaging, "test_contract", "test_policy_set", &policyOutput, 1)

		expectDeletePolicy(t, client, "test_policy", imaging.PolicyNetworkStaging, "test_contract", "test_policy_set")
		expectDeletePolicy(t, client, "test_policy", imaging.PolicyNetworkProduction, "test_contract", "test_policy_set")

		useClient(client, func() {
			resource.UnitTest(t, resource.TestCase{
				ProviderFactories: testAccProviders,
				Steps: []resource.TestStep{
					{
						Config: loadFixtureString(fmt.Sprintf("%s/policy_create.tf", testDir)),
					},
					{
						ImportState:       true,
						ImportStateId:     "test_policy:test_policy_set:test_contract",
						ResourceName:      "akamai_imaging_policy_video.policy",
						ImportStateVerify: true,
					},
				},
			})
		})
		client.AssertExpectations(t)
	})
	t.Run("import policy with activate_on_production=false and no policy on production", func(t *testing.T) {
		testDir := "testdata/TestResPolicyVideo/regular_policy"

		client := new(imaging.Mock)
		expectUpsertPolicy(t, client, "test_policy", imaging.PolicyNetworkStaging, "test_contract", "test_policy_set", &policyInput)
		expectReadPolicy(t, client, "test_policy", imaging.PolicyNetworkStaging, "test_contract", "test_policy_set", &policyOutput, 3)
		client.On("GetPolicy", mock.Anything, imaging.GetPolicyRequest{
			PolicyID:    "test_policy",
			Network:     imaging.PolicyNetworkProduction,
			ContractID:  "test_contract",
			PolicySetID: "test_policy_set",
		}).Return(nil, fmt.Errorf("%s: %w", imaging.ErrGetPolicy, &imaging.Error{Status: http.StatusNotFound})).Once()
		expectReadPolicy(t, client, "test_policy", imaging.PolicyNetworkStaging, "test_contract", "test_policy_set", &policyOutput, 1)

		expectDeletePolicy(t, client, "test_policy", imaging.PolicyNetworkStaging, "test_contract", "test_policy_set")
		expectDeletePolicy(t, client, "test_policy", imaging.PolicyNetworkProduction, "test_contract", "test_policy_set")

		useClient(client, func() {
			resource.UnitTest(t, resource.TestCase{
				ProviderFactories: testAccProviders,
				Steps: []resource.TestStep{
					{
						Config: loadFixtureString(fmt.Sprintf("%s/policy_create.tf", testDir)),
					},
					{
						ImportState:       true,
						ImportStateId:     "test_policy:test_policy_set:test_contract",
						ResourceName:      "akamai_imaging_policy_video.policy",
						ImportStateVerify: true,
					},
				},
			})
		})
		client.AssertExpectations(t)
	})
	t.Run("policy with invalid policy structure", func(t *testing.T) {
		testDir := "testdata/TestResPolicyVideo/invalid_policy"

		client := new(imaging.Mock)
		useClient(client, func() {
			resource.UnitTest(t, resource.TestCase{
				ProviderFactories: testAccProviders,
				Steps: []resource.TestStep{
					{
						Config:      loadFixtureString(fmt.Sprintf("%s/policy_create.tf", testDir)),
						ExpectError: regexp.MustCompile("\"json\" contains an invalid JSON: invalid character '6' looking for beginning of object key string"),
					},
				},
			})
		})
		client.AssertExpectations(t)
	})
	t.Run("policy with unexpected transformation", func(t *testing.T) {
		testDir := "testdata/TestResPolicyVideo/invalid_field_transformation_policy"

		client := new(imaging.Mock)
		expectUpsertPolicy(t, client, "test_policy", imaging.PolicyNetworkStaging, "test_contract", "test_policy_set", &policyInput)
		expectReadPolicy(t, client, "test_policy", imaging.PolicyNetworkStaging, "test_contract", "test_policy_set", &policyOutput, 2)

		expectDeletePolicy(t, client, "test_policy", imaging.PolicyNetworkStaging, "test_contract", "test_policy_set")
		expectDeletePolicy(t, client, "test_policy", imaging.PolicyNetworkProduction, "test_contract", "test_policy_set")

		useClient(client, func() {
			resource.UnitTest(t, resource.TestCase{
				ProviderFactories: testAccProviders,
				Steps: []resource.TestStep{
					{
						Config: loadFixtureString(fmt.Sprintf("%s/policy_create.tf", testDir)),
						Check: checkPolicyAttributes(policyAttributes{
							version:              "1",
							policyID:             "test_policy",
							policySetID:          "test_policy_set",
							activateOnProduction: "false",
							policyPath:           fmt.Sprintf("%s/policy/policy_create.json", testDir),
						}),
					},
				},
			})
		})
		client.AssertExpectations(t)
	})
	t.Run("error when creating policy", func(t *testing.T) {
		testDir := "testdata/TestResPolicyVideo/regular_policy"

		client := new(imaging.Mock)
		var withError = imaging.Error{
			Type:      "/imaging/error-types/IVM_9000",
			Title:     "Bad Request",
			Instance:  "/imaging/error-instances/52a21f40-9861-4d35-95d0-a603c85cb2ad",
			Status:    400,
			Detail:    "Unable to parse element 'output' in JSON.",
			ProblemID: "52a21f40-9861-4d35-95d0-a603c85cb2ad",
		}
		expectUpsertPolicyWithError(t, client, "test_policy", imaging.PolicyNetworkStaging, "test_contract", "test_policy_set", &policyInput, &withError)

		useClient(client, func() {
			resource.UnitTest(t, resource.TestCase{
				ProviderFactories: testAccProviders,
				Steps: []resource.TestStep{
					{
						Config:      loadFixtureString(fmt.Sprintf("%s/policy_create.tf", testDir)),
						ExpectError: regexp.MustCompile("\"detail\": \"Unable to parse element 'output' in JSON.\","),
					},
				},
			})
		})
		client.AssertExpectations(t)
	})
	t.Run("invalid import id", func(t *testing.T) {
		testDir := "testdata/TestResPolicyVideo/regular_policy"

		client := new(imaging.Mock)
		expectUpsertPolicy(t, client, "test_policy", imaging.PolicyNetworkStaging, "test_contract", "test_policy_set", &policyInput)
		expectReadPolicy(t, client, "test_policy", imaging.PolicyNetworkStaging, "test_contract", "test_policy_set", &policyOutput, 2)

		expectDeletePolicy(t, client, "test_policy", imaging.PolicyNetworkStaging, "test_contract", "test_policy_set")
		expectDeletePolicy(t, client, "test_policy", imaging.PolicyNetworkProduction, "test_contract", "test_policy_set")

		useClient(client, func() {
			resource.UnitTest(t, resource.TestCase{
				ProviderFactories: testAccProviders,
				Steps: []resource.TestStep{
					{
						Config: loadFixtureString(fmt.Sprintf("%s/policy_create.tf", testDir)),
					},
					{
						ImportState:       true,
						ImportStateId:     "test_policy,test_policy_set",
						ResourceName:      "akamai_imaging_policy_video.policy",
						ImportStateVerify: true,
						ExpectError:       regexp.MustCompile("colon-separated list of policy ID, policy set ID and contract ID has to be supplied in import: test_policy,test_policy_set"),
					},
				},
			})
		})
		client.AssertExpectations(t)
	})
}

func getPolicyOutputVideoV2(policyOutput imaging.PolicyOutputVideo) imaging.PolicyOutputVideo {
	var policyOutputV2 = policyOutput
	policyOutputV2.Output = &imaging.OutputVideo{
		PerceptualQuality: &imaging.OutputVideoPerceptualQualityVariableInline{
			Value: imaging.OutputVideoPerceptualQualityPtr(imaging.OutputVideoPerceptualQualityMedium),
		},
	}
	policyOutputV2.Version = 2
	return policyOutputV2
}

func getPolicyInputVideoV2(policyInput imaging.PolicyInputVideo) imaging.PolicyInputVideo {
	var policyInputV2 = policyInput
	policyInputV2.Output = &imaging.OutputVideo{
		PerceptualQuality: &imaging.OutputVideoPerceptualQualityVariableInline{
			Value: imaging.OutputVideoPerceptualQualityPtr(imaging.OutputVideoPerceptualQualityMedium),
		},
	}
	return policyInputV2
}

func TestDiffSuppressPolicyVideo(t *testing.T) {
	basePath := "testdata/TestResPolicyVideo/diff_suppress"
	tests := map[string]struct {
		oldPath, newPath string
		expected         bool
	}{
		"identical JSON": {
			oldPath:  "policy.json",
			newPath:  "policy.json",
			expected: true,
		},
		"different formatting, same content": {
			oldPath:  "policy.json",
			newPath:  "different_format.json",
			expected: true,
		},
		"invalid 'old' json": {
			oldPath:  "invalid.json",
			newPath:  "policy.json",
			expected: false,
		},
		"invalid 'new' json": {
			oldPath:  "policy.json",
			newPath:  "invalid.json",
			expected: false,
		},
		"different order of breakpoints.Widths, hosts, variables - no diff": {
			oldPath:  "/fields/policy.json",
			newPath:  "/fields/policy_diff_order.json",
			expected: true,
		},
	}

	for name, test := range tests {
		t.Run(name, func(t *testing.T) {
			oldJSON := loadFixtureString(fmt.Sprintf("%s/%s", basePath, test.oldPath))
			newJSON := loadFixtureString(fmt.Sprintf("%s/%s", basePath, test.newPath))
			res := diffSuppressPolicyVideo("", oldJSON, newJSON, nil)
			assert.Equal(t, test.expected, res)
		})
	}
}
