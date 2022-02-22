package imaging

import (
	"fmt"
	"regexp"
	"testing"

	"github.com/akamai/AkamaiOPEN-edgegrid-golang/v2/pkg/cloudlets/tools"
	"github.com/akamai/AkamaiOPEN-edgegrid-golang/v2/pkg/imaging"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/resource"
	"github.com/stretchr/testify/mock"
	"github.com/tj/assert"
)

func TestResourcePolicy(t *testing.T) {

	type policyAttributes struct {
		version              string
		policyPath           string
		policyID             string
		policySetID          string
		activateOnProduction string
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
			Video:   false,
		}

		expectUpsertPolicy = func(_ *testing.T, client *mockimaging, policyID string, network imaging.PolicyNetwork, contractID string, policySetID string, policy imaging.PolicyInput) {
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

		expectReadPolicy = func(t *testing.T, client *mockimaging, policyID string, network imaging.PolicyNetwork, contractID string, policySetID string, policyOutput imaging.PolicyOutput, times int) {
			client.On("GetPolicy", mock.Anything, imaging.GetPolicyRequest{
				PolicyID:    policyID,
				Network:     network,
				ContractID:  contractID,
				PolicySetID: policySetID,
			}).Return(policyOutput, nil).Times(times)
		}

		expectDeletePolicy = func(_ *testing.T, client *mockimaging, policyID string, network imaging.PolicyNetwork, contractID string, policySetID string) {
			response := imaging.PolicyResponse{}
			client.On("DeletePolicy", mock.Anything, imaging.DeletePolicyRequest{
				PolicyID:    policyID,
				Network:     network,
				ContractID:  contractID,
				PolicySetID: policySetID,
			}).Return(&response, nil).Once()
		}

		expectUpsertPolicyWithError = func(_ *testing.T, client *mockimaging, policyID string, network imaging.PolicyNetwork, contractID string, policySetID string, policy imaging.PolicyInput, err error) {
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
			fs := []resource.TestCheckFunc{
				resource.TestCheckResourceAttr("akamai_imaging_policy_image.policy", "id", fmt.Sprintf("%s:%s", attrs.policySetID, attrs.policyID)),
				resource.TestCheckResourceAttr("akamai_imaging_policy_image.policy", "policy_id", attrs.policyID),
				resource.TestCheckResourceAttr("akamai_imaging_policy_image.policy", "policyset_id", attrs.policySetID),
				resource.TestCheckResourceAttr("akamai_imaging_policy_image.policy", "contract_id", "1YY1"),
				resource.TestCheckResourceAttr("akamai_imaging_policy_image.policy", "version", attrs.version),
				resource.TestCheckResourceAttr("akamai_imaging_policy_image.policy", "json", policyJSON),
			}
			if attrs.activateOnProduction != "" {
				fs = append(fs, resource.TestCheckResourceAttr("akamai_imaging_policy_image.policy", "activate_on_production", attrs.activateOnProduction))
			}
			return resource.ComposeAggregateTestCheckFunc(fs...)
		}
	)

	t.Run("regular policy create", func(t *testing.T) {
		testDir := "testdata/TestResPolicy/regular_policy"

		client := new(mockimaging)
		expectUpsertPolicy(t, client, "test_policy", imaging.PolicyNetworkStaging, "1YY1", "123", &policyInput)
		expectReadPolicy(t, client, "test_policy", imaging.PolicyNetworkStaging, "1YY1", "123", &policyOutput, 2)

		expectDeletePolicy(t, client, "test_policy", imaging.PolicyNetworkStaging, "1YY1", "123")
		// it is faster to attempt to delete on production than checking if there is policy on production first
		expectDeletePolicy(t, client, "test_policy", imaging.PolicyNetworkProduction, "1YY1", "123")

		useClient(client, func() {
			resource.UnitTest(t, resource.TestCase{
				Providers: testAccProviders,
				Steps: []resource.TestStep{
					{
						Config: loadFixtureString(fmt.Sprintf("%s/policy_create.tf", testDir)),
						Check: checkPolicyAttributes(policyAttributes{
							version:     "1",
							policyID:    "test_policy",
							policySetID: "123",
							policyPath:  fmt.Sprintf("%s/policy/policy_create.json", testDir),
						}),
					},
				},
			})
		})
		client.AssertExpectations(t)
	})
	t.Run("regular policy create and later activate on production", func(t *testing.T) {
		testDir := "testdata/TestResPolicy/regular_policy"

		client := new(mockimaging)
		expectUpsertPolicy(t, client, "test_policy", imaging.PolicyNetworkStaging, "1YY1", "123", &policyInput)
		expectReadPolicy(t, client, "test_policy", imaging.PolicyNetworkStaging, "1YY1", "123", &policyOutput, 5)

		// `activate_on_production` should not trigger Upsert for staging if the policy has not changed
		expectUpsertPolicy(t, client, "test_policy", imaging.PolicyNetworkProduction, "1YY1", "123", &policyInput)

		expectDeletePolicy(t, client, "test_policy", imaging.PolicyNetworkStaging, "1YY1", "123")
		expectDeletePolicy(t, client, "test_policy", imaging.PolicyNetworkProduction, "1YY1", "123")

		useClient(client, func() {
			resource.UnitTest(t, resource.TestCase{
				Providers: testAccProviders,
				Steps: []resource.TestStep{
					{
						Config: loadFixtureString(fmt.Sprintf("%s/policy_create.tf", testDir)),
						Check: checkPolicyAttributes(policyAttributes{
							version:     "1",
							policyID:    "test_policy",
							policySetID: "123",
							policyPath:  fmt.Sprintf("%s/policy/policy_create.json", testDir),
						}),
					},
					{
						Config: loadFixtureString(fmt.Sprintf("%s/policy_update.tf", testDir)),
						Check: checkPolicyAttributes(policyAttributes{
							version:              "1",
							policyID:             "test_policy",
							policySetID:          "123",
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
		testDir := "testdata/TestResPolicy/regular_policy_activate_same_time"

		client := new(mockimaging)
		expectReadPolicy(t, client, "test_policy", imaging.PolicyNetworkStaging, "1YY1", "123", &policyOutput, 2)

		expectUpsertPolicy(t, client, "test_policy", imaging.PolicyNetworkStaging, "1YY1", "123", &policyInput)
		expectUpsertPolicy(t, client, "test_policy", imaging.PolicyNetworkProduction, "1YY1", "123", &policyInput)

		// update
		expectReadPolicy(t, client, "test_policy", imaging.PolicyNetworkStaging, "1YY1", "123", &policyOutput, 1)

		policyInputV2 := getPolicyInputV2(policyInput)
		policyOutputV2 := getPolicyOutputV2(policyOutput)

		expectReadPolicy(t, client, "test_policy", imaging.PolicyNetworkStaging, "1YY1", "123", &policyOutputV2, 2)
		expectUpsertPolicy(t, client, "test_policy", imaging.PolicyNetworkStaging, "1YY1", "123", &policyInputV2)
		expectUpsertPolicy(t, client, "test_policy", imaging.PolicyNetworkProduction, "1YY1", "123", &policyInputV2)

		expectDeletePolicy(t, client, "test_policy", imaging.PolicyNetworkStaging, "1YY1", "123")
		expectDeletePolicy(t, client, "test_policy", imaging.PolicyNetworkProduction, "1YY1", "123")

		useClient(client, func() {
			resource.UnitTest(t, resource.TestCase{
				Providers: testAccProviders,
				Steps: []resource.TestStep{
					{
						Config: loadFixtureString(fmt.Sprintf("%s/policy_create.tf", testDir)),
						Check: checkPolicyAttributes(policyAttributes{
							version:              "1",
							policyID:             "test_policy",
							policySetID:          "123",
							activateOnProduction: "true",
							policyPath:           fmt.Sprintf("%s/policy/policy_create.json", testDir),
						}),
					},
					{
						Config: loadFixtureString(fmt.Sprintf("%s/policy_update.tf", testDir)),
						Check: checkPolicyAttributes(policyAttributes{
							version:              "2",
							policyID:             "test_policy",
							policySetID:          "123",
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
		testDir := "testdata/TestResPolicy/change_policyset_id"

		client := new(mockimaging)
		expectUpsertPolicy(t, client, "test_policy", imaging.PolicyNetworkStaging, "1YY1", "123", &policyInput)
		expectReadPolicy(t, client, "test_policy", imaging.PolicyNetworkStaging, "1YY1", "123", &policyOutput, 3)

		// remove original policy
		expectDeletePolicy(t, client, "test_policy", imaging.PolicyNetworkStaging, "1YY1", "123")
		expectDeletePolicy(t, client, "test_policy", imaging.PolicyNetworkProduction, "1YY1", "123")

		// update
		expectUpsertPolicy(t, client, "test_policy", imaging.PolicyNetworkStaging, "1YY1", "456", &policyInput)
		expectReadPolicy(t, client, "test_policy", imaging.PolicyNetworkStaging, "1YY1", "456", &policyOutput, 2)

		// remove new policy
		expectDeletePolicy(t, client, "test_policy", imaging.PolicyNetworkStaging, "1YY1", "456")
		expectDeletePolicy(t, client, "test_policy", imaging.PolicyNetworkProduction, "1YY1", "456")

		useClient(client, func() {
			resource.UnitTest(t, resource.TestCase{
				Providers: testAccProviders,
				Steps: []resource.TestStep{
					{
						Config: loadFixtureString(fmt.Sprintf("%s/policy_create.tf", testDir)),
						Check: checkPolicyAttributes(policyAttributes{
							version:     "1",
							policyID:    "test_policy",
							policySetID: "123",
							policyPath:  fmt.Sprintf("%s/policy/policy_create.json", testDir),
						}),
					},
					{
						Config: loadFixtureString(fmt.Sprintf("%s/policy_update.tf", testDir)),
						Check: checkPolicyAttributes(policyAttributes{
							version:     "1",
							policyID:    "test_policy",
							policySetID: "456",
							policyPath:  fmt.Sprintf("%s/policy/policy_update.json", testDir),
						}),
					},
				},
			})
		})
		client.AssertExpectations(t)
	})
	t.Run("regular policy create, later activate on production and later modify on staging only", func(t *testing.T) {
		testDir := "testdata/TestResPolicy/regular_policy_update_staging"

		client := new(mockimaging)
		expectUpsertPolicy(t, client, "test_policy", imaging.PolicyNetworkStaging, "1YY1", "123", &policyInput)
		expectReadPolicy(t, client, "test_policy", imaging.PolicyNetworkStaging, "1YY1", "123", &policyOutput, 6)

		// `activate_on_production` should not trigger Upsert for staging if the policy has not changed
		expectUpsertPolicy(t, client, "test_policy", imaging.PolicyNetworkProduction, "1YY1", "123", &policyInput)

		policyInputV2 := getPolicyInputV2(policyInput)
		policyOutputV2 := getPolicyOutputV2(policyOutput)

		expectReadPolicy(t, client, "test_policy", imaging.PolicyNetworkStaging, "1YY1", "123", &policyOutputV2, 2)
		expectUpsertPolicy(t, client, "test_policy", imaging.PolicyNetworkStaging, "1YY1", "123", &policyInputV2)

		expectDeletePolicy(t, client, "test_policy", imaging.PolicyNetworkStaging, "1YY1", "123")
		expectDeletePolicy(t, client, "test_policy", imaging.PolicyNetworkProduction, "1YY1", "123")

		useClient(client, func() {
			resource.UnitTest(t, resource.TestCase{
				Providers: testAccProviders,
				Steps: []resource.TestStep{
					{
						Config: loadFixtureString(fmt.Sprintf("%s/policy_create.tf", testDir)),
						Check: checkPolicyAttributes(policyAttributes{
							version:     "1",
							policyID:    "test_policy",
							policySetID: "123",
							policyPath:  fmt.Sprintf("%s/policy/policy_create.json", testDir),
						}),
					},
					{
						Config: loadFixtureString(fmt.Sprintf("%s/policy_update.tf", testDir)),
						Check: checkPolicyAttributes(policyAttributes{
							version:     "1",
							policyID:    "test_policy",
							policySetID: "123",
							policyPath:  fmt.Sprintf("%s/policy/policy_update.json", testDir),
						}),
					},
					{
						Config: loadFixtureString(fmt.Sprintf("%s/policy_update_staging.tf", testDir)),
						Check: checkPolicyAttributes(policyAttributes{
							version:     "2",
							policyID:    "test_policy",
							policySetID: "123",
							policyPath:  fmt.Sprintf("%s/policy/policy_update_staging.json", testDir),
						}),
					},
				},
			})
		})
		client.AssertExpectations(t)
	})
	t.Run("auto policy create and later activate on production, cannot delete", func(t *testing.T) {
		testDir := "testdata/TestResPolicy/auto_policy"

		client := new(mockimaging)
		expectUpsertPolicy(t, client, ".auto", imaging.PolicyNetworkStaging, "1YY1", "123", &policyInput)
		expectReadPolicy(t, client, ".auto", imaging.PolicyNetworkStaging, "1YY1", "123", &policyOutput, 5)

		// `activate_on_production` should not trigger Upsert for staging if the policy has not changed
		expectUpsertPolicy(t, client, ".auto", imaging.PolicyNetworkProduction, "1YY1", "123", &policyInput)

		// .auto policy cannot be removed alone, only via removal of policy set

		useClient(client, func() {
			resource.UnitTest(t, resource.TestCase{
				Providers: testAccProviders,
				Steps: []resource.TestStep{
					{
						Config: loadFixtureString(fmt.Sprintf("%s/policy_create.tf", testDir)),
						Check: checkPolicyAttributes(policyAttributes{
							version:     "1",
							policyID:    ".auto",
							policySetID: "123",
							policyPath:  fmt.Sprintf("%s/policy/policy_create.json", testDir),
						}),
					},
					{
						Config: loadFixtureString(fmt.Sprintf("%s/policy_update.tf", testDir)),
						Check: checkPolicyAttributes(policyAttributes{
							version:     "1",
							policyID:    ".auto",
							policySetID: "123",
							policyPath:  fmt.Sprintf("%s/policy/policy_update.json", testDir),
						}),
					},
				},
			})
		})
		client.AssertExpectations(t)
	})
	t.Run("policy with invalid policy structure", func(t *testing.T) {
		testDir := "testdata/TestResPolicy/invalid_policy"

		client := new(mockimaging)
		useClient(client, func() {
			resource.UnitTest(t, resource.TestCase{
				Providers: testAccProviders,
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
	t.Run("policy with inconsistent policy structure", func(t *testing.T) {
		testDir := "testdata/TestResPolicy/inconsistent_policy"

		client := new(mockimaging)
		useClient(client, func() {
			resource.UnitTest(t, resource.TestCase{
				Providers: testAccProviders,
				Steps: []resource.TestStep{
					{
						Config:      loadFixtureString(fmt.Sprintf("%s/policy_create.tf", testDir)),
						ExpectError: regexp.MustCompile("unmarshalling transformation list: unsupported transformation type: MaxColors3"),
					},
				},
			})
		})
		client.AssertExpectations(t)
	})
	t.Run("error when creating policy", func(t *testing.T) {
		testDir := "testdata/TestResPolicy/regular_policy"

		client := new(mockimaging)
		var withError = imaging.Error{
			Type:      "https://problems.luna.akamaiapis.net/image-policy-manager/IVM_1004",
			Title:     "Bad Request",
			Instance:  "52a21f40-9861-4d35-95d0-a603c85cb2ad",
			Status:    400,
			Detail:    "Policy fails to be properly created by AkaImaging: Unrecognized transformation type: MaxColors2",
			ProblemID: "52a21f40-9861-4d35-95d0-a603c85cb2ad",
		}
		expectUpsertPolicyWithError(t, client, "test_policy", imaging.PolicyNetworkStaging, "1YY1", "123", &policyInput, &withError)

		useClient(client, func() {
			resource.UnitTest(t, resource.TestCase{
				Providers: testAccProviders,
				Steps: []resource.TestStep{
					{
						Config:      loadFixtureString(fmt.Sprintf("%s/policy_create.tf", testDir)),
						ExpectError: regexp.MustCompile("\"detail\": \"Policy fails to be properly created by AkaImaging: Unrecognized transformation type: MaxColors2\","),
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

func TestDiffSuppressPolicy(t *testing.T) {
	basePath := "testdata/TestResPolicy/diff_suppress"
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
	}

	for name, test := range tests {
		t.Run(name, func(t *testing.T) {
			oldJSON := loadFixtureString(fmt.Sprintf("%s/%s", basePath, test.oldPath))
			newJSON := loadFixtureString(fmt.Sprintf("%s/%s", basePath, test.newPath))
			res := diffSuppressPolicy("", oldJSON, newJSON, nil)
			assert.Equal(t, test.expected, res)
		})
	}
}
