package imaging

import (
	"errors"
	"fmt"
	"net/http"
	"regexp"
	"testing"

	"github.com/akamai/AkamaiOPEN-edgegrid-golang/v7/pkg/imaging"
	"github.com/akamai/terraform-provider-akamai/v5/pkg/common/testutils"
	"github.com/akamai/terraform-provider-akamai/v5/pkg/tools"
	"github.com/hashicorp/terraform-plugin-testing/helper/resource"
	"github.com/stretchr/testify/mock"
	"github.com/tj/assert"
)

func TestResourcePolicyImage(t *testing.T) {

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
		policyInput = imaging.PolicyInputImage{
			Breakpoints: &imaging.Breakpoints{
				Widths: []int{320, 640, 1024, 2048, 5000},
			},
			Output: &imaging.OutputImage{
				PerceptualQuality: &imaging.OutputImagePerceptualQualityVariableInline{
					Value: imaging.OutputImagePerceptualQualityPtr(imaging.OutputImagePerceptualQualityMediumHigh),
				},
			},
			Transformations: []imaging.TransformationType{
				&imaging.MaxColors{
					Colors: &imaging.IntegerVariableInline{
						Value: tools.IntPtr(2),
					},
					Transformation: imaging.MaxColorsTransformationMaxColors,
				},
			},
		}
		policyOutput = imaging.PolicyOutputImage{
			Breakpoints: &imaging.Breakpoints{
				Widths: []int{320, 640, 1024, 2048, 5000},
			},
			Output: &imaging.OutputImage{
				PerceptualQuality: &imaging.OutputImagePerceptualQualityVariableInline{
					Value: imaging.OutputImagePerceptualQualityPtr(imaging.OutputImagePerceptualQualityMediumHigh),
				},
			},
			Transformations: []imaging.TransformationType{
				&imaging.MaxColors{
					Colors: &imaging.IntegerVariableInline{
						Value: tools.IntPtr(2),
					},
					Transformation: imaging.MaxColorsTransformationMaxColors,
				},
			},
			Version: 1,
			Video:   tools.BoolPtr(false),
		}
		defaultAllowedFormats = []imaging.OutputImageAllowedFormats{
			imaging.OutputImageAllowedFormats("jpeg"),
			imaging.OutputImageAllowedFormats("webp"),
			imaging.OutputImageAllowedFormats("avif"),
			imaging.OutputImageAllowedFormats("png"),
			imaging.OutputImageAllowedFormats("gif"),
		}
		defaultForcedFormats = []imaging.OutputImageForcedFormats{
			imaging.OutputImageForcedFormats("jpeg"),
			imaging.OutputImageForcedFormats("webp"),
			imaging.OutputImageForcedFormats("avif"),
			imaging.OutputImageForcedFormats("png"),
			imaging.OutputImageForcedFormats("gif"),
		}
		perceptualQualityMediumHigh = &imaging.OutputImagePerceptualQualityVariableInline{
			Value: imaging.OutputImagePerceptualQualityPtr(imaging.OutputImagePerceptualQualityMediumHigh),
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
		defaultBreakpointsWidths = &imaging.Breakpoints{
			Widths: []int{320, 640, 1024, 2048, 5000},
		}
		policyInputDiff = imaging.PolicyInputImage{
			Breakpoints: defaultBreakpointsWidths,
			Hosts:       defaultHosts,
			Output: &imaging.OutputImage{
				AllowedFormats:    defaultAllowedFormats,
				ForcedFormats:     defaultForcedFormats,
				PerceptualQuality: perceptualQualityMediumHigh,
			},
			Variables: defaultVariables,
		}
		policyOutputDiff = imaging.PolicyOutputImage{
			Breakpoints: defaultBreakpointsWidths,
			Hosts:       defaultHosts,
			Output: &imaging.OutputImage{
				AllowedFormats:    defaultAllowedFormats,
				ForcedFormats:     defaultForcedFormats,
				PerceptualQuality: perceptualQualityMediumHigh,
			},
			Variables: defaultVariables,
			Version:   1,
			Video:     tools.BoolPtr(false),
		}

		expectUpsertPolicy = func(client *imaging.Mock, policyID, policySetID, contractID string, network imaging.PolicyNetwork, policy imaging.PolicyInput) {
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

		expectUpsertPolicyFailure = func(client *imaging.Mock, policyID, policySetID, contractID string, network imaging.PolicyNetwork, policy imaging.PolicyInput) {
			client.On("UpsertPolicy", mock.Anything, imaging.UpsertPolicyRequest{
				PolicyID:    policyID,
				Network:     network,
				ContractID:  contractID,
				PolicySetID: policySetID,
				PolicyInput: policy,
			}).Return(nil, errors.New("API error: Conflict (409)")).Once()
		}

		expectReadPolicy = func(client *imaging.Mock, policyID, policySetID, contractID string, network imaging.PolicyNetwork, policyOutput imaging.PolicyOutput, times int) {
			client.On("GetPolicy", mock.Anything, imaging.GetPolicyRequest{
				PolicyID:    policyID,
				Network:     network,
				ContractID:  contractID,
				PolicySetID: policySetID,
			}).Return(policyOutput, nil).Times(times)
		}

		expectDeletePolicy = func(client *imaging.Mock, policyID, policySetID, contractID string, network imaging.PolicyNetwork) {
			response := imaging.PolicyResponse{}
			client.On("DeletePolicy", mock.Anything, imaging.DeletePolicyRequest{
				PolicyID:    policyID,
				Network:     network,
				ContractID:  contractID,
				PolicySetID: policySetID,
			}).Return(&response, nil).Once()
		}

		expectUpsertPolicyWithError = func(client *imaging.Mock, policyID, policySetID, contractID string, network imaging.PolicyNetwork, policy imaging.PolicyInput, err error) {
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
				policyJSON = testutils.LoadFixtureString(t, attrs.policyPath)
			}
			f := resource.ComposeAggregateTestCheckFunc(
				resource.TestCheckResourceAttr("akamai_imaging_policy_image.policy", "id", fmt.Sprintf("%s:%s", attrs.policySetID, attrs.policyID)),
				resource.TestCheckResourceAttr("akamai_imaging_policy_image.policy", "policy_id", attrs.policyID),
				resource.TestCheckResourceAttr("akamai_imaging_policy_image.policy", "policyset_id", attrs.policySetID),
				resource.TestCheckResourceAttr("akamai_imaging_policy_image.policy", "contract_id", "test_contract"),
				resource.TestCheckResourceAttr("akamai_imaging_policy_image.policy", "version", attrs.version),
				resource.TestCheckResourceAttr("akamai_imaging_policy_image.policy", "activate_on_production", attrs.activateOnProduction),
			)
			if attrs.schema {
				if attrs.emptyPolicy {
					return resource.ComposeAggregateTestCheckFunc(
						f,
						resource.TestCheckResourceAttr("akamai_imaging_policy_image.policy", "policy.#", "1"),
					)
				}
				return resource.ComposeAggregateTestCheckFunc(
					f,
					resource.TestCheckResourceAttr("akamai_imaging_policy_image.policy", "policy.#", "1"),
					resource.TestCheckResourceAttr("akamai_imaging_policy_image.policy", "policy.0.output.0.perceptual_quality", "mediumHigh"),
				)
			}
			return resource.ComposeAggregateTestCheckFunc(
				f,
				resource.TestCheckResourceAttr("akamai_imaging_policy_image.policy", "json", policyJSON),
			)
		}
	)

	t.Run("regular policy create", func(t *testing.T) {
		testDir := "testdata/TestResPolicyImage/regular_policy"

		client := new(imaging.Mock)
		expectUpsertPolicy(client, "test_policy", "test_policy_set", "test_contract", imaging.PolicyNetworkStaging, &policyInput)
		expectReadPolicy(client, "test_policy", "test_policy_set", "test_contract", imaging.PolicyNetworkStaging, &policyOutput, 2)

		expectDeletePolicy(client, "test_policy", "test_policy_set", "test_contract", imaging.PolicyNetworkStaging)
		// it is faster to attempt to delete on production than checking if there is policy on production first
		expectDeletePolicy(client, "test_policy", "test_policy_set", "test_contract", imaging.PolicyNetworkProduction)

		useClient(client, func() {
			resource.UnitTest(t, resource.TestCase{
				ProtoV6ProviderFactories: testAccProviders,
				Steps: []resource.TestStep{
					{
						Config: testutils.LoadFixtureString(t, fmt.Sprintf("%s/policy_create.tf", testDir)),
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
		testDir := "testdata/TestResPolicyImage/regular_policy"

		client := new(imaging.Mock)
		expectUpsertPolicy(client, "test_policy", "test_policy_set", "test_contract", imaging.PolicyNetworkStaging, &policyInput)
		expectReadPolicy(client, "test_policy", "test_policy_set", "test_contract", imaging.PolicyNetworkStaging, &policyOutput, 3)

		// `activate_on_production` should not trigger Upsert for staging if the policy has not changed
		expectUpsertPolicy(client, "test_policy", "test_policy_set", "test_contract", imaging.PolicyNetworkProduction, &policyInput)
		expectReadPolicy(client, "test_policy", "test_policy_set", "test_contract", imaging.PolicyNetworkProduction, &policyOutput, 2)

		expectDeletePolicy(client, "test_policy", "test_policy_set", "test_contract", imaging.PolicyNetworkStaging)
		expectDeletePolicy(client, "test_policy", "test_policy_set", "test_contract", imaging.PolicyNetworkProduction)

		useClient(client, func() {
			resource.UnitTest(t, resource.TestCase{
				ProtoV6ProviderFactories: testAccProviders,
				Steps: []resource.TestStep{
					{
						Config: testutils.LoadFixtureString(t, fmt.Sprintf("%s/policy_create.tf", testDir)),
						Check: checkPolicyAttributes(policyAttributes{
							version:              "1",
							policyID:             "test_policy",
							policySetID:          "test_policy_set",
							activateOnProduction: "false",
							policyPath:           fmt.Sprintf("%s/policy/policy_create.json", testDir),
						}),
					},
					{
						Config: testutils.LoadFixtureString(t, fmt.Sprintf("%s/policy_update.tf", testDir)),
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
		testDir := "testdata/TestResPolicyImage/regular_policy_activate_same_time"

		client := new(imaging.Mock)
		expectReadPolicy(client, "test_policy", "test_policy_set", "test_contract", imaging.PolicyNetworkProduction, &policyOutput, 2)

		expectUpsertPolicy(client, "test_policy", "test_policy_set", "test_contract", imaging.PolicyNetworkStaging, &policyInput)
		expectUpsertPolicy(client, "test_policy", "test_policy_set", "test_contract", imaging.PolicyNetworkProduction, &policyInput)

		// update
		expectReadPolicy(client, "test_policy", "test_policy_set", "test_contract", imaging.PolicyNetworkProduction, &policyOutput, 1)

		policyInputV2 := getPolicyInputV2(policyInput)
		policyOutputV2 := getPolicyOutputV2(policyOutput)

		expectUpsertPolicy(client, "test_policy", "test_policy_set", "test_contract", imaging.PolicyNetworkStaging, &policyInputV2)
		expectUpsertPolicy(client, "test_policy", "test_policy_set", "test_contract", imaging.PolicyNetworkProduction, &policyInputV2)
		expectReadPolicy(client, "test_policy", "test_policy_set", "test_contract", imaging.PolicyNetworkProduction, &policyOutputV2, 2)

		expectDeletePolicy(client, "test_policy", "test_policy_set", "test_contract", imaging.PolicyNetworkStaging)
		expectDeletePolicy(client, "test_policy", "test_policy_set", "test_contract", imaging.PolicyNetworkProduction)

		useClient(client, func() {
			resource.UnitTest(t, resource.TestCase{
				ProtoV6ProviderFactories: testAccProviders,
				Steps: []resource.TestStep{
					{
						Config: testutils.LoadFixtureString(t, fmt.Sprintf("%s/policy_create.tf", testDir)),
						Check: checkPolicyAttributes(policyAttributes{
							version:              "1",
							policyID:             "test_policy",
							policySetID:          "test_policy_set",
							activateOnProduction: "true",
							policyPath:           fmt.Sprintf("%s/policy/policy_create.json", testDir),
						}),
					},
					{
						Config: testutils.LoadFixtureString(t, fmt.Sprintf("%s/policy_update.tf", testDir)),
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
		testDir := "testdata/TestResPolicyImage/regular_policy_activate_same_time"

		client := new(imaging.Mock)
		expectReadPolicy(client, "test_policy", "test_policy_set", "test_contract", imaging.PolicyNetworkProduction, &policyOutput, 2)

		expectUpsertPolicy(client, "test_policy", "test_policy_set", "test_contract", imaging.PolicyNetworkStaging, &policyInput)
		expectUpsertPolicy(client, "test_policy", "test_policy_set", "test_contract", imaging.PolicyNetworkProduction, &policyInput)

		expectReadPolicy(client, "test_policy", "test_policy_set", "test_contract", imaging.PolicyNetworkProduction, &policyOutput, 1)

		// update
		policyInputV2 := getPolicyInputV2(policyInput)

		expectUpsertPolicy(client, "test_policy", "test_policy_set", "test_contract", imaging.PolicyNetworkStaging, &policyInputV2)
		expectUpsertPolicyFailure(client, "test_policy", "test_policy_set", "test_contract", imaging.PolicyNetworkProduction, &policyInputV2)

		expectDeletePolicy(client, "test_policy", "test_policy_set", "test_contract", imaging.PolicyNetworkStaging)
		expectDeletePolicy(client, "test_policy", "test_policy_set", "test_contract", imaging.PolicyNetworkProduction)

		useClient(client, func() {
			resource.UnitTest(t, resource.TestCase{
				ProviderFactories: testAccProviders,
				Steps: []resource.TestStep{
					{
						Config: testutils.LoadFixtureString(t, fmt.Sprintf("%s/policy_create.tf", testDir)),
						Check: checkPolicyAttributes(policyAttributes{
							version:              "1",
							policyID:             "test_policy",
							policySetID:          "test_policy_set",
							activateOnProduction: "true",
							policyPath:           fmt.Sprintf("%s/policy/policy_create.json", testDir),
						}),
					},
					{
						Config:      testutils.LoadFixtureString(t, fmt.Sprintf("%s/policy_update.tf", testDir)),
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
		testDir := "testdata/TestResPolicyImage/change_policyset_id"

		client := new(imaging.Mock)
		expectUpsertPolicy(client, "test_policy", "test_policy_set", "test_contract", imaging.PolicyNetworkStaging, &policyInput)
		expectReadPolicy(client, "test_policy", "test_policy_set", "test_contract", imaging.PolicyNetworkStaging, &policyOutput, 3)

		// remove original policy
		expectDeletePolicy(client, "test_policy", "test_policy_set", "test_contract", imaging.PolicyNetworkStaging)
		expectDeletePolicy(client, "test_policy", "test_policy_set", "test_contract", imaging.PolicyNetworkProduction)

		// update
		expectUpsertPolicy(client, "test_policy", "test_policy_set_update", "test_contract", imaging.PolicyNetworkStaging, &policyInput)
		expectReadPolicy(client, "test_policy", "test_policy_set_update", "test_contract", imaging.PolicyNetworkStaging, &policyOutput, 2)

		// remove new policy
		expectDeletePolicy(client, "test_policy", "test_policy_set_update", "test_contract", imaging.PolicyNetworkStaging)
		expectDeletePolicy(client, "test_policy", "test_policy_set_update", "test_contract", imaging.PolicyNetworkProduction)

		useClient(client, func() {
			resource.UnitTest(t, resource.TestCase{
				ProtoV6ProviderFactories: testAccProviders,
				Steps: []resource.TestStep{
					{
						Config: testutils.LoadFixtureString(t, fmt.Sprintf("%s/policy_create.tf", testDir)),
						Check: checkPolicyAttributes(policyAttributes{
							version:              "1",
							policyID:             "test_policy",
							policySetID:          "test_policy_set",
							activateOnProduction: "false",
							policyPath:           fmt.Sprintf("%s/policy/policy_create.json", testDir),
						}),
					},
					{
						Config: testutils.LoadFixtureString(t, fmt.Sprintf("%s/policy_update.tf", testDir)),
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
		testDir := "testdata/TestResPolicyImage/regular_policy_update_staging"

		client := new(imaging.Mock)

		expectUpsertPolicy(client, "test_policy", "test_policy_set", "test_contract", imaging.PolicyNetworkStaging, &policyInput)
		expectReadPolicy(client, "test_policy", "test_policy_set", "test_contract", imaging.PolicyNetworkStaging, &policyOutput, 3)

		// `activate_on_production` should not trigger Upsert for staging if the policy has not changed
		expectUpsertPolicy(client, "test_policy", "test_policy_set", "test_contract", imaging.PolicyNetworkProduction, &policyInput)
		expectReadPolicy(client, "test_policy", "test_policy_set", "test_contract", imaging.PolicyNetworkProduction, &policyOutput, 3)

		policyInputV2 := getPolicyInputV2(policyInput)
		policyOutputV2 := getPolicyOutputV2(policyOutput)

		expectReadPolicy(client, "test_policy", "test_policy_set", "test_contract", imaging.PolicyNetworkStaging, &policyOutputV2, 2)
		expectUpsertPolicy(client, "test_policy", "test_policy_set", "test_contract", imaging.PolicyNetworkStaging, &policyInputV2)

		expectDeletePolicy(client, "test_policy", "test_policy_set", "test_contract", imaging.PolicyNetworkStaging)
		expectDeletePolicy(client, "test_policy", "test_policy_set", "test_contract", imaging.PolicyNetworkProduction)

		useClient(client, func() {
			resource.UnitTest(t, resource.TestCase{
				ProtoV6ProviderFactories: testAccProviders,
				Steps: []resource.TestStep{
					{
						Config: testutils.LoadFixtureString(t, fmt.Sprintf("%s/policy_create.tf", testDir)),
						Check: checkPolicyAttributes(policyAttributes{
							version:              "1",
							policyID:             "test_policy",
							policySetID:          "test_policy_set",
							activateOnProduction: "false",
							policyPath:           fmt.Sprintf("%s/policy/policy_create.json", testDir),
						}),
					},
					{
						Config: testutils.LoadFixtureString(t, fmt.Sprintf("%s/policy_update.tf", testDir)),
						Check: checkPolicyAttributes(policyAttributes{
							version:              "1",
							policyID:             "test_policy",
							policySetID:          "test_policy_set",
							activateOnProduction: "true",
							policyPath:           fmt.Sprintf("%s/policy/policy_update.json", testDir),
						}),
					},
					{
						Config: testutils.LoadFixtureString(t, fmt.Sprintf("%s/policy_update_staging.tf", testDir)),
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
		testDir := "testdata/TestResPolicyImage/auto_policy"

		client := new(imaging.Mock)
		expectUpsertPolicy(client, ".auto", "test_policy_set", "test_contract", imaging.PolicyNetworkStaging, &policyInput)
		expectReadPolicy(client, ".auto", "test_policy_set", "test_contract", imaging.PolicyNetworkStaging, &policyOutput, 3)

		// `activate_on_production` should not trigger Upsert for staging if the policy has not changed
		expectUpsertPolicy(client, ".auto", "test_policy_set", "test_contract", imaging.PolicyNetworkProduction, &policyInput)
		expectReadPolicy(client, ".auto", "test_policy_set", "test_contract", imaging.PolicyNetworkProduction, &policyOutput, 2)

		// .auto policy cannot be removed alone, only via removal of policy set

		useClient(client, func() {
			resource.UnitTest(t, resource.TestCase{
				ProtoV6ProviderFactories: testAccProviders,
				Steps: []resource.TestStep{
					{
						Config: testutils.LoadFixtureString(t, fmt.Sprintf("%s/policy_create.tf", testDir)),
						Check: checkPolicyAttributes(policyAttributes{
							version:              "1",
							policyID:             ".auto",
							policySetID:          "test_policy_set",
							activateOnProduction: "false",
							policyPath:           fmt.Sprintf("%s/policy/policy_create.json", testDir),
						}),
					},
					{
						Config: testutils.LoadFixtureString(t, fmt.Sprintf("%s/policy_update.tf", testDir)),
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
	t.Run("create policy, check diff in order of output.allowedFormats, output.forcedFormats, breakpoints.Widths, hosts, variables - no diff", func(t *testing.T) {
		testDir := "testdata/TestResPolicyImage/diff_suppress/fields"

		client := new(imaging.Mock)
		expectUpsertPolicy(client, "test_policy", "test_policy_set", "test_contract", imaging.PolicyNetworkStaging, &policyInputDiff)
		expectReadPolicy(client, "test_policy", "test_policy_set", "test_contract", imaging.PolicyNetworkStaging, &policyOutputDiff, 3)

		// remove original policy
		expectDeletePolicy(client, "test_policy", "test_policy_set", "test_contract", imaging.PolicyNetworkStaging)
		expectDeletePolicy(client, "test_policy", "test_policy_set", "test_contract", imaging.PolicyNetworkProduction)

		expectReadPolicy(client, "test_policy", "test_policy_set", "test_contract", imaging.PolicyNetworkStaging, &policyOutputDiff, 1)

		useClient(client, func() {
			resource.UnitTest(t, resource.TestCase{
				ProtoV6ProviderFactories: testAccProviders,
				Steps: []resource.TestStep{
					{
						Config: testutils.LoadFixtureString(t, fmt.Sprintf("%s/default.tf", testDir)),
						Check: checkPolicyAttributes(policyAttributes{
							version:              "1",
							policyID:             "test_policy",
							policySetID:          "test_policy_set",
							activateOnProduction: "false",
							policyPath:           fmt.Sprintf("%s/policy.json", testDir),
						}),
					},
					{
						Config:             testutils.LoadFixtureString(t, fmt.Sprintf("%s/diff_order.tf", testDir)),
						ExpectNonEmptyPlan: false,
					},
				},
			})
		})
		client.AssertExpectations(t)
	})
	t.Run("update rollout duration and ensure no diff", func(t *testing.T) {
		testDir := "testdata/TestResPolicyImage/regular_policy_update_rollout_duration"

		policyInputWithRollout := imaging.PolicyInputImage{
			RolloutDuration: tools.IntPtr(3600),
			Breakpoints: &imaging.Breakpoints{
				Widths: []int{320, 640, 1024, 2048, 5000},
			},
			Output: &imaging.OutputImage{
				PerceptualQuality: &imaging.OutputImagePerceptualQualityVariableInline{
					Value: imaging.OutputImagePerceptualQualityPtr(imaging.OutputImagePerceptualQualityMediumHigh),
				},
			},
			Transformations: []imaging.TransformationType{
				&imaging.MaxColors{
					Colors: &imaging.IntegerVariableInline{
						Value: tools.IntPtr(2),
					},
					Transformation: imaging.MaxColorsTransformationMaxColors,
				},
			},
		}

		client := new(imaging.Mock)
		expectUpsertPolicy(client, "test_policy", "test_policy_set", "test_contract", imaging.PolicyNetworkStaging, &policyInput)
		expectReadPolicy(client, "test_policy", "test_policy_set", "test_contract", imaging.PolicyNetworkStaging, &policyOutput, 2)

		expectUpsertPolicy(client, "test_policy", "test_policy_set", "test_contract", imaging.PolicyNetworkStaging, &policyInputWithRollout)
		expectReadPolicy(client, "test_policy", "test_policy_set", "test_contract", imaging.PolicyNetworkStaging, &policyOutput, 3)

		expectDeletePolicy(client, "test_policy", "test_policy_set", "test_contract", imaging.PolicyNetworkStaging)
		expectDeletePolicy(client, "test_policy", "test_policy_set", "test_contract", imaging.PolicyNetworkProduction)

		useClient(client, func() {
			resource.UnitTest(t, resource.TestCase{
				ProtoV6ProviderFactories: testAccProviders,
				Steps: []resource.TestStep{
					{
						Config: testutils.LoadFixtureString(t, fmt.Sprintf("%s/policy_create.tf", testDir)),
						Check: checkPolicyAttributes(policyAttributes{
							version:              "1",
							policyID:             "test_policy",
							policySetID:          "test_policy_set",
							activateOnProduction: "false",
							policyPath:           fmt.Sprintf("%s/policy/policy_create.json", testDir),
						}),
					},
					{
						Config: testutils.LoadFixtureString(t, fmt.Sprintf("%s/policy_update.tf", testDir)),
						Check: checkPolicyAttributes(policyAttributes{
							version:              "1",
							policyID:             "test_policy",
							policySetID:          "test_policy_set",
							activateOnProduction: "false",
							policyPath:           fmt.Sprintf("%s/policy/policy_update.json", testDir),
						}),
					},
				},
			})
		})
		client.AssertExpectations(t)
	})
	t.Run("update serve stale duration and ensure no diff", func(t *testing.T) {
		testDir := "testdata/TestResPolicyImage/regular_policy_update_serve_stale_duration"

		policyInputWithServeStale := imaging.PolicyInputImage{
			ServeStaleDuration: tools.IntPtr(3600),
			Breakpoints: &imaging.Breakpoints{
				Widths: []int{320, 640, 1024, 2048, 5000},
			},
			Output: &imaging.OutputImage{
				AllowPristineOnDownsize: tools.BoolPtr(true),
				PerceptualQuality: &imaging.OutputImagePerceptualQualityVariableInline{
					Value: imaging.OutputImagePerceptualQualityPtr(imaging.OutputImagePerceptualQualityMediumHigh),
				},
				PreferModernFormats: tools.BoolPtr(false),
			},
			Transformations: []imaging.TransformationType{
				&imaging.MaxColors{
					Colors: &imaging.IntegerVariableInline{
						Value: tools.IntPtr(2),
					},
					Transformation: imaging.MaxColorsTransformationMaxColors,
				},
			},
		}

		client := new(imaging.Mock)
		expectUpsertPolicy(client, "test_policy", "test_policy_set", "test_contract", imaging.PolicyNetworkStaging, &policyInput)
		expectReadPolicy(client, "test_policy", "test_policy_set", "test_contract", imaging.PolicyNetworkStaging, &policyOutput, 2)

		policyOutputAfterUpdate := imaging.PolicyOutputImage{
			Breakpoints: &imaging.Breakpoints{
				Widths: []int{320, 640, 1024, 2048, 5000},
			},
			Output: &imaging.OutputImage{
				AllowPristineOnDownsize: tools.BoolPtr(true),
				PerceptualQuality: &imaging.OutputImagePerceptualQualityVariableInline{
					Value: imaging.OutputImagePerceptualQualityPtr(imaging.OutputImagePerceptualQualityMediumHigh),
				},
				PreferModernFormats: tools.BoolPtr(false),
			},
			Transformations: []imaging.TransformationType{
				&imaging.MaxColors{
					Colors: &imaging.IntegerVariableInline{
						Value: tools.IntPtr(2),
					},
					Transformation: imaging.MaxColorsTransformationMaxColors,
				},
			},
			Version: 1,
			Video:   tools.BoolPtr(false),
		}

		expectUpsertPolicy(client, "test_policy", "test_policy_set", "test_contract", imaging.PolicyNetworkStaging, &policyInputWithServeStale)
		expectReadPolicy(client, "test_policy", "test_policy_set", "test_contract", imaging.PolicyNetworkStaging, &policyOutputAfterUpdate, 3)

		expectDeletePolicy(client, "test_policy", "test_policy_set", "test_contract", imaging.PolicyNetworkStaging)
		expectDeletePolicy(client, "test_policy", "test_policy_set", "test_contract", imaging.PolicyNetworkProduction)

		useClient(client, func() {
			resource.UnitTest(t, resource.TestCase{
				ProviderFactories: testAccProviders,
				Steps: []resource.TestStep{
					{
						Config: testutils.LoadFixtureString(t, fmt.Sprintf("%s/policy_create.tf", testDir)),
						Check: checkPolicyAttributes(policyAttributes{
							version:              "1",
							policyID:             "test_policy",
							policySetID:          "test_policy_set",
							activateOnProduction: "false",
							policyPath:           fmt.Sprintf("%s/policy/policy_create.json", testDir),
						}),
					},
					{
						Config: testutils.LoadFixtureString(t, fmt.Sprintf("%s/policy_update.tf", testDir)),
						Check: checkPolicyAttributes(policyAttributes{
							version:              "1",
							policyID:             "test_policy",
							policySetID:          "test_policy_set",
							activateOnProduction: "false",
							policyPath:           fmt.Sprintf("%s/policy/policy_update.json", testDir),
						}),
					},
				},
			})
		})
		client.AssertExpectations(t)
	})
	t.Run("import policy with activate_on_production=true", func(t *testing.T) {
		testDir := "testdata/TestResPolicyImage/regular_policy"

		client := new(imaging.Mock)
		expectUpsertPolicy(client, "test_policy", "test_policy_set", "test_contract", imaging.PolicyNetworkStaging, &policyInput)
		expectReadPolicy(client, "test_policy", "test_policy_set", "test_contract", imaging.PolicyNetworkStaging, &policyOutput, 3)

		expectUpsertPolicy(client, "test_policy", "test_policy_set", "test_contract", imaging.PolicyNetworkProduction, &policyInput)
		expectReadPolicy(client, "test_policy", "test_policy_set", "test_contract", imaging.PolicyNetworkProduction, &policyOutput, 3)

		expectReadPolicy(client, "test_policy", "test_policy_set", "test_contract", imaging.PolicyNetworkProduction, &policyOutput, 1)
		expectReadPolicy(client, "test_policy", "test_policy_set", "test_contract", imaging.PolicyNetworkStaging, &policyOutput, 1)

		expectDeletePolicy(client, "test_policy", "test_policy_set", "test_contract", imaging.PolicyNetworkStaging)
		expectDeletePolicy(client, "test_policy", "test_policy_set", "test_contract", imaging.PolicyNetworkProduction)

		useClient(client, func() {
			resource.UnitTest(t, resource.TestCase{
				ProtoV6ProviderFactories: testAccProviders,
				Steps: []resource.TestStep{
					{
						Config: testutils.LoadFixtureString(t, fmt.Sprintf("%s/policy_create.tf", testDir)),
					},
					{
						Config: testutils.LoadFixtureString(t, fmt.Sprintf("%s/policy_update.tf", testDir)),
					},
					{
						ImportState:       true,
						ImportStateId:     "test_policy:test_policy_set:test_contract",
						ResourceName:      "akamai_imaging_policy_image.policy",
						ImportStateVerify: true,
					},
				},
			})
		})
		client.AssertExpectations(t)
	})
	t.Run("import policy with activate_on_production=false", func(t *testing.T) {
		testDir := "testdata/TestResPolicyImage/regular_policy"
		policyOutputV2 := getPolicyOutputV2(policyOutput)

		client := new(imaging.Mock)
		expectUpsertPolicy(client, "test_policy", "test_policy_set", "test_contract", imaging.PolicyNetworkStaging, &policyInput)
		expectReadPolicy(client, "test_policy", "test_policy_set", "test_contract", imaging.PolicyNetworkStaging, &policyOutput, 3)
		expectReadPolicy(client, "test_policy", "test_policy_set", "test_contract", imaging.PolicyNetworkProduction, &policyOutputV2, 1)
		expectReadPolicy(client, "test_policy", "test_policy_set", "test_contract", imaging.PolicyNetworkStaging, &policyOutput, 1)

		expectDeletePolicy(client, "test_policy", "test_policy_set", "test_contract", imaging.PolicyNetworkStaging)
		expectDeletePolicy(client, "test_policy", "test_policy_set", "test_contract", imaging.PolicyNetworkProduction)

		useClient(client, func() {
			resource.UnitTest(t, resource.TestCase{
				ProtoV6ProviderFactories: testAccProviders,
				Steps: []resource.TestStep{
					{
						Config: testutils.LoadFixtureString(t, fmt.Sprintf("%s/policy_create.tf", testDir)),
					},
					{
						ImportState:       true,
						ImportStateId:     "test_policy:test_policy_set:test_contract",
						ResourceName:      "akamai_imaging_policy_image.policy",
						ImportStateVerify: true,
					},
				},
			})
		})
		client.AssertExpectations(t)
	})
	t.Run("import policy with activate_on_production=false and no policy on production", func(t *testing.T) {
		testDir := "testdata/TestResPolicyImage/regular_policy"

		client := new(imaging.Mock)
		expectUpsertPolicy(client, "test_policy", "test_policy_set", "test_contract", imaging.PolicyNetworkStaging, &policyInput)
		expectReadPolicy(client, "test_policy", "test_policy_set", "test_contract", imaging.PolicyNetworkStaging, &policyOutput, 3)
		client.On("GetPolicy", mock.Anything, imaging.GetPolicyRequest{
			PolicyID:    "test_policy",
			Network:     imaging.PolicyNetworkProduction,
			ContractID:  "test_contract",
			PolicySetID: "test_policy_set",
		}).Return(nil, fmt.Errorf("%s: %w", imaging.ErrGetPolicy, &imaging.Error{Status: http.StatusNotFound})).Once()
		expectReadPolicy(client, "test_policy", "test_policy_set", "test_contract", imaging.PolicyNetworkStaging, &policyOutput, 1)

		expectDeletePolicy(client, "test_policy", "test_policy_set", "test_contract", imaging.PolicyNetworkStaging)
		expectDeletePolicy(client, "test_policy", "test_policy_set", "test_contract", imaging.PolicyNetworkProduction)

		useClient(client, func() {
			resource.UnitTest(t, resource.TestCase{
				ProtoV6ProviderFactories: testAccProviders,
				Steps: []resource.TestStep{
					{
						Config: testutils.LoadFixtureString(t, fmt.Sprintf("%s/policy_create.tf", testDir)),
					},
					{
						ImportState:       true,
						ImportStateId:     "test_policy:test_policy_set:test_contract",
						ResourceName:      "akamai_imaging_policy_image.policy",
						ImportStateVerify: true,
					},
				},
			})
		})
		client.AssertExpectations(t)
	})
	t.Run("import policy with activate_on_production=false and no policy on production with rolloutDuration", func(t *testing.T) {
		testDir := "testdata/TestResPolicyImage/regular_policy_import"

		policyInput := policyInput
		policyInput.RolloutDuration = tools.IntPtr(3600)

		client := new(imaging.Mock)
		expectUpsertPolicy(client, "test_policy", "test_policy_set", "test_contract", imaging.PolicyNetworkStaging, &policyInput)
		expectReadPolicy(client, "test_policy", "test_policy_set", "test_contract", imaging.PolicyNetworkStaging, &policyOutput, 3)
		client.On("GetPolicy", mock.Anything, imaging.GetPolicyRequest{
			PolicyID:    "test_policy",
			Network:     imaging.PolicyNetworkProduction,
			ContractID:  "test_contract",
			PolicySetID: "test_policy_set",
		}).Return(nil, fmt.Errorf("%s: %w", imaging.ErrGetPolicy, &imaging.Error{Status: http.StatusNotFound})).Once()
		expectReadPolicy(client, "test_policy", "test_policy_set", "test_contract", imaging.PolicyNetworkStaging, &policyOutput, 1)
		expectReadPolicy(client, "test_policy", "test_policy_set", "test_contract", imaging.PolicyNetworkStaging, &policyOutput, 1)
		expectReadPolicy(client, "test_policy", "test_policy_set", "test_contract", imaging.PolicyNetworkStaging, &policyOutput, 1)

		expectDeletePolicy(client, "test_policy", "test_policy_set", "test_contract", imaging.PolicyNetworkStaging)
		expectDeletePolicy(client, "test_policy", "test_policy_set", "test_contract", imaging.PolicyNetworkProduction)

		useClient(client, func() {
			resource.UnitTest(t, resource.TestCase{
				ProtoV6ProviderFactories: testAccProviders,
				Steps: []resource.TestStep{
					{
						Config: testutils.LoadFixtureString(t, fmt.Sprintf("%s/policy_create.tf", testDir)),
					},
					{
						ImportState:   true,
						ImportStateId: "test_policy:test_policy_set:test_contract",
						ResourceName:  "akamai_imaging_policy_image.policy",
						// Current implementation is unable to handle correctly `rolloutDuration` during import.
						// It is recommended to not provide any value for that field before import.
						// `cli-terraform` will not set this field during export, assuming that it service will required default value.
						//ImportStateVerify: true,
					},
					{
						Config:   testutils.LoadFixtureString(t, fmt.Sprintf("%s/policy_create.tf", testDir)),
						PlanOnly: true,
					},
				},
			})
		})
		client.AssertExpectations(t)
	})
	t.Run("policy with invalid policy structure", func(t *testing.T) {
		testDir := "testdata/TestResPolicyImage/invalid_policy"

		client := new(imaging.Mock)
		useClient(client, func() {
			resource.UnitTest(t, resource.TestCase{
				ProtoV6ProviderFactories: testAccProviders,
				Steps: []resource.TestStep{
					{
						Config:      testutils.LoadFixtureString(t, fmt.Sprintf("%s/policy_create.tf", testDir)),
						ExpectError: regexp.MustCompile("\"json\" contains an invalid JSON: invalid character '6' looking for beginning of object key string"),
					},
				},
			})
		})
		client.AssertExpectations(t)
	})
	t.Run("policy with inconsistent policy structure", func(t *testing.T) {
		testDir := "testdata/TestResPolicyImage/inconsistent_policy"

		client := new(imaging.Mock)
		useClient(client, func() {
			resource.UnitTest(t, resource.TestCase{
				ProtoV6ProviderFactories: testAccProviders,
				Steps: []resource.TestStep{
					{
						Config:      testutils.LoadFixtureString(t, fmt.Sprintf("%s/policy_create.tf", testDir)),
						ExpectError: regexp.MustCompile("unmarshalling transformation list: unsupported transformation type: MaxColors3"),
					},
				},
			})
		})
		client.AssertExpectations(t)
	})
	t.Run("error when creating policy", func(t *testing.T) {
		testDir := "testdata/TestResPolicyImage/regular_policy"

		client := new(imaging.Mock)
		var withError = imaging.Error{
			Type:      "https://problems.luna.akamaiapis.net/image-policy-manager/IVM_1004",
			Title:     "Bad Request",
			Instance:  "52a21f40-9861-4d35-95d0-a603c85cb2ad",
			Status:    400,
			Detail:    "Policy fails to be properly created by AkaImaging: Unrecognized transformation type: MaxColors2",
			ProblemID: "52a21f40-9861-4d35-95d0-a603c85cb2ad",
		}
		expectUpsertPolicyWithError(client, "test_policy", "test_policy_set", "test_contract", imaging.PolicyNetworkStaging, &policyInput, &withError)

		useClient(client, func() {
			resource.UnitTest(t, resource.TestCase{
				ProtoV6ProviderFactories: testAccProviders,
				Steps: []resource.TestStep{
					{
						Config:      testutils.LoadFixtureString(t, fmt.Sprintf("%s/policy_create.tf", testDir)),
						ExpectError: regexp.MustCompile("\"detail\": \"Policy fails to be properly created by AkaImaging: Unrecognized transformation type: MaxColors2\","),
					},
				},
			})
		})
		client.AssertExpectations(t)
	})
	t.Run("invalid import id", func(t *testing.T) {
		testDir := "testdata/TestResPolicyImage/regular_policy"

		client := new(imaging.Mock)
		expectUpsertPolicy(client, "test_policy", "test_policy_set", "test_contract", imaging.PolicyNetworkStaging, &policyInput)
		expectReadPolicy(client, "test_policy", "test_policy_set", "test_contract", imaging.PolicyNetworkStaging, &policyOutput, 2)

		expectDeletePolicy(client, "test_policy", "test_policy_set", "test_contract", imaging.PolicyNetworkStaging)
		expectDeletePolicy(client, "test_policy", "test_policy_set", "test_contract", imaging.PolicyNetworkProduction)

		useClient(client, func() {
			resource.UnitTest(t, resource.TestCase{
				ProtoV6ProviderFactories: testAccProviders,
				Steps: []resource.TestStep{
					{
						Config: testutils.LoadFixtureString(t, fmt.Sprintf("%s/policy_create.tf", testDir)),
					},
					{
						ImportState:       true,
						ImportStateId:     "test_policy,test_policy_set",
						ResourceName:      "akamai_imaging_policy_image.policy",
						ImportStateVerify: true,
						ExpectError:       regexp.MustCompile("colon-separated list of policy ID, policy set ID and contract ID has to be supplied in import: test_policy,test_policy_set"),
					},
				},
			})
		})
		client.AssertExpectations(t)
	})
}

func getPolicyOutputV2(policyOutput imaging.PolicyOutputImage) imaging.PolicyOutputImage {
	var policyOutputV2 = policyOutput
	policyOutputV2.Transformations = []imaging.TransformationType{
		&imaging.MaxColors{
			Colors: &imaging.IntegerVariableInline{
				Value: tools.IntPtr(3),
			},
			Transformation: imaging.MaxColorsTransformationMaxColors,
		},
	}
	policyOutputV2.Version = 2
	return policyOutputV2
}

func getPolicyInputV2(policyInput imaging.PolicyInputImage) imaging.PolicyInputImage {
	var policyInputV2 = policyInput
	policyInputV2.Transformations = []imaging.TransformationType{
		&imaging.MaxColors{
			Colors: &imaging.IntegerVariableInline{
				Value: tools.IntPtr(3),
			},
			Transformation: imaging.MaxColorsTransformationMaxColors,
		},
	}
	return policyInputV2
}

func getPolicyOutputOrderV2(policyOutput imaging.PolicyOutputImage) imaging.PolicyOutputImage {
	var policyOutputV2 = policyOutput
	policyOutputV2.Transformations = []imaging.TransformationType{
		&imaging.Blur{
			Sigma:          &imaging.NumberVariableInline{Value: tools.Float64Ptr(5)},
			Transformation: imaging.BlurTransformationBlur,
		},
		&imaging.MaxColors{
			Colors: &imaging.IntegerVariableInline{
				Value: tools.IntPtr(4),
			},
			Transformation: imaging.MaxColorsTransformationMaxColors,
		},
	}
	policyOutputV2.Version = 2
	return policyOutputV2
}

func TestDiffSuppressPolicy(t *testing.T) {
	basePath := "testdata/TestResPolicyImage/diff_suppress"
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
		"policy diff in output.allowedFormats, output.forcedFormats, breakpoints.Widths, hosts, variables": {
			oldPath:  "fields/policy.json",
			newPath:  "fields/policy_diff_order.json",
			expected: true,
		},
	}

	for name, test := range tests {
		t.Run(name, func(t *testing.T) {
			oldJSON := testutils.LoadFixtureString(t, fmt.Sprintf("%s/%s", basePath, test.oldPath))
			newJSON := testutils.LoadFixtureString(t, fmt.Sprintf("%s/%s", basePath, test.newPath))
			res := diffSuppressPolicyImage("", oldJSON, newJSON, nil)
			assert.Equal(t, test.expected, res)
		})
	}
}
