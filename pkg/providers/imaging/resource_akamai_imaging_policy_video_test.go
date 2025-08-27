package imaging

import (
	"errors"
	"fmt"
	"net/http"
	"regexp"
	"testing"

	"github.com/akamai/AkamaiOPEN-edgegrid-golang/v12/pkg/imaging"
	"github.com/akamai/terraform-provider-akamai/v8/pkg/common/ptr"
	"github.com/akamai/terraform-provider-akamai/v8/pkg/common/testutils"
	"github.com/hashicorp/terraform-plugin-testing/helper/resource"
	"github.com/stretchr/testify/assert"
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
			Video:   ptr.To(true),
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
			Video:       ptr.To(true),
		}

		expectUpsertPolicy = func(client *imaging.Mock, policyID, contractID, policySetID string, network imaging.PolicyNetwork, policy imaging.PolicyInput) {
			policyResponse := &imaging.PolicyResponse{
				OperationPerformed: "UPDATED",
				Description:        fmt.Sprintf("Policy %s updated.", policyID),
				ID:                 policyID,
			}
			client.On("UpsertPolicy", testutils.MockContext, imaging.UpsertPolicyRequest{
				PolicyID:    policyID,
				Network:     network,
				ContractID:  contractID,
				PolicySetID: policySetID,
				PolicyInput: policy,
			}).Return(policyResponse, nil).Once()
		}

		expectUpsertPolicyFailure = func(client *imaging.Mock, policyID, policySetID, contractID string, network imaging.PolicyNetwork, policy imaging.PolicyInput) {
			client.On("UpsertPolicy", testutils.MockContext, imaging.UpsertPolicyRequest{
				PolicyID:    policyID,
				Network:     network,
				ContractID:  contractID,
				PolicySetID: policySetID,
				PolicyInput: policy,
			}).Return(nil, errors.New("API error: Conflict (409)")).Once()
		}

		expectReadPolicy = func(client *imaging.Mock, policyID, contractID, policySetID string, network imaging.PolicyNetwork, policyOutput imaging.PolicyOutput, times int) {
			client.On("GetPolicy", testutils.MockContext, imaging.GetPolicyRequest{
				PolicyID:    policyID,
				Network:     network,
				ContractID:  contractID,
				PolicySetID: policySetID,
			}).Return(policyOutput, nil).Times(times)
		}

		expectDeletePolicy = func(client *imaging.Mock, policyID, contractID, policySetID string, network imaging.PolicyNetwork) {
			response := imaging.PolicyResponse{}
			client.On("DeletePolicy", testutils.MockContext, imaging.DeletePolicyRequest{
				PolicyID:    policyID,
				Network:     network,
				ContractID:  contractID,
				PolicySetID: policySetID,
			}).Return(&response, nil).Once()
		}

		expectUpsertPolicyWithError = func(client *imaging.Mock, policyID, contractID, policySetID string, network imaging.PolicyNetwork, policy imaging.PolicyInput, err error) {
			client.On("UpsertPolicy", testutils.MockContext, imaging.UpsertPolicyRequest{
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
				policyJSON = testutils.LoadFixtureString(t, attrs.policyPath)
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
		expectUpsertPolicy(client, "test_policy", "test_contract", "test_policy_set", imaging.PolicyNetworkStaging, &policyInput)
		expectReadPolicy(client, "test_policy", "test_contract", "test_policy_set", imaging.PolicyNetworkStaging, &policyOutput, 2)

		expectDeletePolicy(client, "test_policy", "test_contract", "test_policy_set", imaging.PolicyNetworkStaging)
		// it is faster to attempt to delete on production than checking if there is policy on production first
		expectDeletePolicy(client, "test_policy", "test_contract", "test_policy_set", imaging.PolicyNetworkProduction)

		useClient(client, func() {
			resource.UnitTest(t, resource.TestCase{
				ProtoV6ProviderFactories: testutils.NewProtoV6ProviderFactory(NewSubprovider()),
				Steps: []resource.TestStep{
					{
						Config: testutils.LoadFixtureStringf(t, "%s/policy_create.tf", testDir),
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
		expectUpsertPolicy(client, "test_policy", "test_contract", "test_policy_set", imaging.PolicyNetworkStaging, &policyInput)
		expectReadPolicy(client, "test_policy", "test_contract", "test_policy_set", imaging.PolicyNetworkStaging, &policyOutput, 3)

		// `activate_on_production` should not trigger Upsert for staging if the policy has not changed
		expectUpsertPolicy(client, "test_policy", "test_contract", "test_policy_set", imaging.PolicyNetworkProduction, &policyInput)
		expectReadPolicy(client, "test_policy", "test_contract", "test_policy_set", imaging.PolicyNetworkProduction, &policyOutput, 2)

		expectDeletePolicy(client, "test_policy", "test_contract", "test_policy_set", imaging.PolicyNetworkStaging)
		expectDeletePolicy(client, "test_policy", "test_contract", "test_policy_set", imaging.PolicyNetworkProduction)

		useClient(client, func() {
			resource.UnitTest(t, resource.TestCase{
				ProtoV6ProviderFactories: testutils.NewProtoV6ProviderFactory(NewSubprovider()),
				Steps: []resource.TestStep{
					{
						Config: testutils.LoadFixtureStringf(t, "%s/policy_create.tf", testDir),
						Check: checkPolicyAttributes(policyAttributes{
							version:              "1",
							policyID:             "test_policy",
							policySetID:          "test_policy_set",
							activateOnProduction: "false",
							policyPath:           fmt.Sprintf("%s/policy/policy_create.json", testDir),
						}),
					},
					{
						Config: testutils.LoadFixtureStringf(t, "%s/policy_update.tf", testDir),
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
		expectReadPolicy(client, "test_policy", "test_contract", "test_policy_set", imaging.PolicyNetworkProduction, &policyOutput, 2)

		expectUpsertPolicy(client, "test_policy", "test_contract", "test_policy_set", imaging.PolicyNetworkStaging, &policyInput)
		expectUpsertPolicy(client, "test_policy", "test_contract", "test_policy_set", imaging.PolicyNetworkProduction, &policyInput)

		// update
		expectReadPolicy(client, "test_policy", "test_contract", "test_policy_set", imaging.PolicyNetworkProduction, &policyOutput, 1)

		policyInputV2 := getPolicyInputVideoV2(policyInput)
		policyOutputV2 := getPolicyOutputVideoV2(policyOutput)

		expectReadPolicy(client, "test_policy", "test_contract", "test_policy_set", imaging.PolicyNetworkProduction, &policyOutputV2, 2)
		expectUpsertPolicy(client, "test_policy", "test_contract", "test_policy_set", imaging.PolicyNetworkStaging, &policyInputV2)
		expectUpsertPolicy(client, "test_policy", "test_contract", "test_policy_set", imaging.PolicyNetworkProduction, &policyInputV2)

		expectDeletePolicy(client, "test_policy", "test_contract", "test_policy_set", imaging.PolicyNetworkStaging)
		expectDeletePolicy(client, "test_policy", "test_contract", "test_policy_set", imaging.PolicyNetworkProduction)

		useClient(client, func() {
			resource.UnitTest(t, resource.TestCase{
				ProtoV6ProviderFactories: testutils.NewProtoV6ProviderFactory(NewSubprovider()),
				Steps: []resource.TestStep{
					{
						Config: testutils.LoadFixtureStringf(t, "%s/policy_create.tf", testDir),
						Check: checkPolicyAttributes(policyAttributes{
							version:              "1",
							policyID:             "test_policy",
							policySetID:          "test_policy_set",
							activateOnProduction: "true",
							policyPath:           fmt.Sprintf("%s/policy/policy_create.json", testDir),
						}),
					},
					{
						Config: testutils.LoadFixtureStringf(t, "%s/policy_update.tf", testDir),
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
	t.Run("regular policy create with activate_on_production=true, update immediately, fails on production", func(t *testing.T) {
		testDir := "testdata/TestResPolicyVideo/regular_policy_activate_same_time"

		client := new(imaging.Mock)
		expectReadPolicy(client, "test_policy", "test_contract", "test_policy_set", imaging.PolicyNetworkProduction, &policyOutput, 2)

		expectUpsertPolicy(client, "test_policy", "test_contract", "test_policy_set", imaging.PolicyNetworkStaging, &policyInput)
		expectUpsertPolicy(client, "test_policy", "test_contract", "test_policy_set", imaging.PolicyNetworkProduction, &policyInput)

		expectReadPolicy(client, "test_policy", "test_contract", "test_policy_set", imaging.PolicyNetworkProduction, &policyOutput, 1)

		// update
		policyInputV2 := getPolicyInputVideoV2(policyInput)

		expectUpsertPolicy(client, "test_policy", "test_contract", "test_policy_set", imaging.PolicyNetworkStaging, &policyInputV2)
		expectUpsertPolicyFailure(client, "test_policy", "test_policy_set", "test_contract", imaging.PolicyNetworkProduction, &policyInputV2)

		expectDeletePolicy(client, "test_policy", "test_contract", "test_policy_set", imaging.PolicyNetworkStaging)
		expectDeletePolicy(client, "test_policy", "test_contract", "test_policy_set", imaging.PolicyNetworkProduction)

		useClient(client, func() {
			resource.UnitTest(t, resource.TestCase{
				ProtoV6ProviderFactories: testutils.NewProtoV6ProviderFactory(NewSubprovider()),
				Steps: []resource.TestStep{
					{
						Config: testutils.LoadFixtureStringf(t, "%s/policy_create.tf", testDir),
						Check: checkPolicyAttributes(policyAttributes{
							version:              "1",
							policyID:             "test_policy",
							policySetID:          "test_policy_set",
							activateOnProduction: "true",
							policyPath:           fmt.Sprintf("%s/policy/policy_create.json", testDir),
						}),
					},
					{
						Config:      testutils.LoadFixtureStringf(t, "%s/policy_update.tf", testDir),
						ExpectError: regexp.MustCompile(`Error: API error: Conflict \(409\)`),
						Check: checkPolicyAttributes(policyAttributes{
							version:              "1",
							policyID:             "test_policy",
							policySetID:          "test_policy_set",
							activateOnProduction: "true",
							policyPath:           fmt.Sprintf("%s/policy/policy_create.json", testDir),
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
		expectUpsertPolicy(client, "test_policy", "test_contract", "test_policy_set", imaging.PolicyNetworkStaging, &policyInput)
		expectReadPolicy(client, "test_policy", "test_contract", "test_policy_set", imaging.PolicyNetworkStaging, &policyOutput, 3)

		// remove original policy
		expectDeletePolicy(client, "test_policy", "test_contract", "test_policy_set", imaging.PolicyNetworkStaging)
		expectDeletePolicy(client, "test_policy", "test_contract", "test_policy_set", imaging.PolicyNetworkProduction)

		// update
		expectUpsertPolicy(client, "test_policy", "test_contract", "test_policy_set_update", imaging.PolicyNetworkStaging, &policyInput)
		expectReadPolicy(client, "test_policy", "test_contract", "test_policy_set_update", imaging.PolicyNetworkStaging, &policyOutput, 2)

		// remove new policy
		expectDeletePolicy(client, "test_policy", "test_contract", "test_policy_set_update", imaging.PolicyNetworkStaging)
		expectDeletePolicy(client, "test_policy", "test_contract", "test_policy_set_update", imaging.PolicyNetworkProduction)

		useClient(client, func() {
			resource.UnitTest(t, resource.TestCase{
				ProtoV6ProviderFactories: testutils.NewProtoV6ProviderFactory(NewSubprovider()),
				Steps: []resource.TestStep{
					{
						Config: testutils.LoadFixtureStringf(t, "%s/policy_create.tf", testDir),
						Check: checkPolicyAttributes(policyAttributes{
							version:              "1",
							policyID:             "test_policy",
							policySetID:          "test_policy_set",
							activateOnProduction: "false",
							policyPath:           fmt.Sprintf("%s/policy/policy_create.json", testDir),
						}),
					},
					{
						Config: testutils.LoadFixtureStringf(t, "%s/policy_update.tf", testDir),
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
		expectUpsertPolicy(client, "test_policy", "test_contract", "test_policy_set", imaging.PolicyNetworkStaging, &policyInput)
		expectReadPolicy(client, "test_policy", "test_contract", "test_policy_set", imaging.PolicyNetworkStaging, &policyOutput, 3)

		// `activate_on_production` should not trigger Upsert for staging if the policy has not changed
		expectUpsertPolicy(client, "test_policy", "test_contract", "test_policy_set", imaging.PolicyNetworkProduction, &policyInput)
		expectReadPolicy(client, "test_policy", "test_contract", "test_policy_set", imaging.PolicyNetworkProduction, &policyOutput, 3)

		policyInputV2 := getPolicyInputVideoV2(policyInput)
		policyOutputV2 := getPolicyOutputVideoV2(policyOutput)

		expectReadPolicy(client, "test_policy", "test_contract", "test_policy_set", imaging.PolicyNetworkStaging, &policyOutputV2, 2)
		expectUpsertPolicy(client, "test_policy", "test_contract", "test_policy_set", imaging.PolicyNetworkStaging, &policyInputV2)

		expectDeletePolicy(client, "test_policy", "test_contract", "test_policy_set", imaging.PolicyNetworkStaging)
		expectDeletePolicy(client, "test_policy", "test_contract", "test_policy_set", imaging.PolicyNetworkProduction)

		useClient(client, func() {
			resource.UnitTest(t, resource.TestCase{
				ProtoV6ProviderFactories: testutils.NewProtoV6ProviderFactory(NewSubprovider()),
				Steps: []resource.TestStep{
					{
						Config: testutils.LoadFixtureStringf(t, "%s/policy_create.tf", testDir),
						Check: checkPolicyAttributes(policyAttributes{
							version:              "1",
							policyID:             "test_policy",
							policySetID:          "test_policy_set",
							activateOnProduction: "false",
							policyPath:           fmt.Sprintf("%s/policy/policy_create.json", testDir),
						}),
					},
					{
						Config: testutils.LoadFixtureStringf(t, "%s/policy_update.tf", testDir),
						Check: checkPolicyAttributes(policyAttributes{
							version:              "1",
							policyID:             "test_policy",
							policySetID:          "test_policy_set",
							activateOnProduction: "true",
							policyPath:           fmt.Sprintf("%s/policy/policy_update.json", testDir),
						}),
					},
					{
						Config: testutils.LoadFixtureStringf(t, "%s/policy_update_staging.tf", testDir),
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
		expectUpsertPolicy(client, ".auto", "test_contract", "test_policy_set", imaging.PolicyNetworkStaging, &policyInput)
		expectReadPolicy(client, ".auto", "test_contract", "test_policy_set", imaging.PolicyNetworkStaging, &policyOutput, 3)

		// `activate_on_production` should not trigger Upsert for staging if the policy has not changed
		expectUpsertPolicy(client, ".auto", "test_contract", "test_policy_set", imaging.PolicyNetworkProduction, &policyInput)
		expectReadPolicy(client, ".auto", "test_contract", "test_policy_set", imaging.PolicyNetworkProduction, &policyOutput, 2)

		// .auto policy cannot be removed alone, only via removal of policy set

		useClient(client, func() {
			resource.UnitTest(t, resource.TestCase{
				ProtoV6ProviderFactories: testutils.NewProtoV6ProviderFactory(NewSubprovider()),
				Steps: []resource.TestStep{
					{
						Config: testutils.LoadFixtureStringf(t, "%s/policy_create.tf", testDir),
						Check: checkPolicyAttributes(policyAttributes{
							version:              "1",
							policyID:             ".auto",
							policySetID:          "test_policy_set",
							activateOnProduction: "false",
							policyPath:           fmt.Sprintf("%s/policy/policy_create.json", testDir),
						}),
					},
					{
						Config: testutils.LoadFixtureStringf(t, "%s/policy_update.tf", testDir),
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
		expectUpsertPolicy(client, "test_policy", "test_contract", "test_policy_set", imaging.PolicyNetworkStaging, &policyInput)
		expectUpsertPolicy(client, "test_policy", "test_contract", "test_policy_set", imaging.PolicyNetworkProduction, &policyInput)
		expectReadPolicy(client, "test_policy", "test_contract", "test_policy_set", imaging.PolicyNetworkStaging, &policyOutput, 3)

		expectReadPolicy(client, "test_policy", "test_contract", "test_policy_set", imaging.PolicyNetworkProduction, &policyOutput, 4)
		expectReadPolicy(client, "test_policy", "test_contract", "test_policy_set", imaging.PolicyNetworkStaging, &policyOutput, 1)

		expectDeletePolicy(client, "test_policy", "test_contract", "test_policy_set", imaging.PolicyNetworkStaging)
		expectDeletePolicy(client, "test_policy", "test_contract", "test_policy_set", imaging.PolicyNetworkProduction)

		useClient(client, func() {
			resource.UnitTest(t, resource.TestCase{
				ProtoV6ProviderFactories: testutils.NewProtoV6ProviderFactory(NewSubprovider()),
				Steps: []resource.TestStep{
					{
						Config: testutils.LoadFixtureStringf(t, "%s/policy_create.tf", testDir),
					},
					{
						Config: testutils.LoadFixtureStringf(t, "%s/policy_update.tf", testDir),
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
		expectUpsertPolicy(client, "test_policy", "test_contract", "test_policy_set", imaging.PolicyNetworkStaging, &policyInputDiff)
		expectReadPolicy(client, "test_policy", "test_contract", "test_policy_set", imaging.PolicyNetworkStaging, &policyOutputDiff, 2)

		// remove original policy
		expectDeletePolicy(client, "test_policy", "test_contract", "test_policy_set", imaging.PolicyNetworkStaging)
		expectDeletePolicy(client, "test_policy", "test_contract", "test_policy_set", imaging.PolicyNetworkProduction)

		expectReadPolicy(client, "test_policy", "test_contract", "test_policy_set", imaging.PolicyNetworkStaging, &policyOutputDiff, 2)

		useClient(client, func() {
			resource.UnitTest(t, resource.TestCase{
				ProtoV6ProviderFactories: testutils.NewProtoV6ProviderFactory(NewSubprovider()),
				Steps: []resource.TestStep{
					{
						Config: testutils.LoadFixtureStringf(t, "%s/default.tf", testDir),
						Check: checkPolicyAttributes(policyAttributes{
							version:              "1",
							policyID:             "test_policy",
							policySetID:          "test_policy_set",
							activateOnProduction: "false",
							policyPath:           fmt.Sprintf("%s/policy.json", testDir),
						}),
					},
					{
						Config:             testutils.LoadFixtureStringf(t, "%s/diff_order.tf", testDir),
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
		expectUpsertPolicy(client, "test_policy", "test_contract", "test_policy_set", imaging.PolicyNetworkStaging, &policyInput)
		expectReadPolicy(client, "test_policy", "test_contract", "test_policy_set", imaging.PolicyNetworkStaging, &policyOutput, 3)
		expectReadPolicy(client, "test_policy", "test_contract", "test_policy_set", imaging.PolicyNetworkProduction, &policyOutputV2, 1)
		expectReadPolicy(client, "test_policy", "test_contract", "test_policy_set", imaging.PolicyNetworkStaging, &policyOutput, 1)
		expectDeletePolicy(client, "test_policy", "test_contract", "test_policy_set", imaging.PolicyNetworkStaging)
		expectDeletePolicy(client, "test_policy", "test_contract", "test_policy_set", imaging.PolicyNetworkProduction)

		useClient(client, func() {
			resource.UnitTest(t, resource.TestCase{
				ProtoV6ProviderFactories: testutils.NewProtoV6ProviderFactory(NewSubprovider()),
				Steps: []resource.TestStep{
					{
						Config: testutils.LoadFixtureStringf(t, "%s/policy_create.tf", testDir),
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
		expectUpsertPolicy(client, "test_policy", "test_contract", "test_policy_set", imaging.PolicyNetworkStaging, &policyInput)
		expectReadPolicy(client, "test_policy", "test_contract", "test_policy_set", imaging.PolicyNetworkStaging, &policyOutput, 3)
		client.On("GetPolicy", testutils.MockContext, imaging.GetPolicyRequest{
			PolicyID:    "test_policy",
			Network:     imaging.PolicyNetworkProduction,
			ContractID:  "test_contract",
			PolicySetID: "test_policy_set",
		}).Return(nil, fmt.Errorf("%s: %w", imaging.ErrGetPolicy, &imaging.Error{Status: http.StatusNotFound})).Once()
		expectReadPolicy(client, "test_policy", "test_contract", "test_policy_set", imaging.PolicyNetworkStaging, &policyOutput, 1)

		expectDeletePolicy(client, "test_policy", "test_contract", "test_policy_set", imaging.PolicyNetworkStaging)
		expectDeletePolicy(client, "test_policy", "test_contract", "test_policy_set", imaging.PolicyNetworkProduction)

		useClient(client, func() {
			resource.UnitTest(t, resource.TestCase{
				ProtoV6ProviderFactories: testutils.NewProtoV6ProviderFactory(NewSubprovider()),
				Steps: []resource.TestStep{
					{
						Config: testutils.LoadFixtureStringf(t, "%s/policy_create.tf", testDir),
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
				ProtoV6ProviderFactories: testutils.NewProtoV6ProviderFactory(NewSubprovider()),
				Steps: []resource.TestStep{
					{
						Config:      testutils.LoadFixtureStringf(t, "%s/policy_create.tf", testDir),
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
		expectUpsertPolicy(client, "test_policy", "test_contract", "test_policy_set", imaging.PolicyNetworkStaging, &policyInput)
		expectReadPolicy(client, "test_policy", "test_contract", "test_policy_set", imaging.PolicyNetworkStaging, &policyOutput, 2)

		expectDeletePolicy(client, "test_policy", "test_contract", "test_policy_set", imaging.PolicyNetworkStaging)
		expectDeletePolicy(client, "test_policy", "test_contract", "test_policy_set", imaging.PolicyNetworkProduction)

		useClient(client, func() {
			resource.UnitTest(t, resource.TestCase{
				ProtoV6ProviderFactories: testutils.NewProtoV6ProviderFactory(NewSubprovider()),
				Steps: []resource.TestStep{
					{
						Config: testutils.LoadFixtureStringf(t, "%s/policy_create.tf", testDir),
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
		expectUpsertPolicyWithError(client, "test_policy", "test_contract", "test_policy_set", imaging.PolicyNetworkStaging, &policyInput, &withError)

		useClient(client, func() {
			resource.UnitTest(t, resource.TestCase{
				ProtoV6ProviderFactories: testutils.NewProtoV6ProviderFactory(NewSubprovider()),
				Steps: []resource.TestStep{
					{
						Config:      testutils.LoadFixtureStringf(t, "%s/policy_create.tf", testDir),
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
		expectUpsertPolicy(client, "test_policy", "test_contract", "test_policy_set", imaging.PolicyNetworkStaging, &policyInput)
		expectReadPolicy(client, "test_policy", "test_contract", "test_policy_set", imaging.PolicyNetworkStaging, &policyOutput, 2)

		expectDeletePolicy(client, "test_policy", "test_contract", "test_policy_set", imaging.PolicyNetworkStaging)
		expectDeletePolicy(client, "test_policy", "test_contract", "test_policy_set", imaging.PolicyNetworkProduction)

		useClient(client, func() {
			resource.UnitTest(t, resource.TestCase{
				ProtoV6ProviderFactories: testutils.NewProtoV6ProviderFactory(NewSubprovider()),
				Steps: []resource.TestStep{
					{
						Config: testutils.LoadFixtureStringf(t, "%s/policy_create.tf", testDir),
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
			oldJSON := testutils.LoadFixtureStringf(t, "%s/%s", basePath, test.oldPath)
			newJSON := testutils.LoadFixtureStringf(t, "%s/%s", basePath, test.newPath)
			res := diffSuppressPolicyVideo("", oldJSON, newJSON, nil)
			assert.Equal(t, test.expected, res)
		})
	}
}
