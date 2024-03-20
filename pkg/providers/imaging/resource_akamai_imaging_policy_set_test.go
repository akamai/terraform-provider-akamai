package imaging

import (
	"errors"
	"fmt"
	"regexp"
	"testing"

	"github.com/akamai/AkamaiOPEN-edgegrid-golang/v8/pkg/imaging"
	"github.com/akamai/terraform-provider-akamai/v6/pkg/common/testutils"
	"github.com/hashicorp/terraform-plugin-testing/helper/resource"
	"github.com/stretchr/testify/mock"
	"github.com/tj/assert"
)

func TestResourceImagingPolicySet(t *testing.T) {
	var (
		anError = errors.New("oops")

		expectPolicySetCreation = func(t *testing.T, client *imaging.Mock, contractID, name, region, mediaType string, policySet *imaging.PolicySet, createError error) {
			client.On("CreatePolicySet", mock.Anything, imaging.CreatePolicySetRequest{
				ContractID: contractID,
				CreatePolicySet: imaging.CreatePolicySet{
					Name:   name,
					Region: imaging.Region(region),
					Type:   imaging.MediaType(mediaType),
				},
			}).Return(policySet, createError).Once()
		}

		expectPolicySetRead = func(t *testing.T, client *imaging.Mock, contractID, policySetID string, policySet *imaging.PolicySet, getPolicyError error, times int) {
			client.On("GetPolicySet", mock.Anything, imaging.GetPolicySetRequest{
				PolicySetID: policySetID, ContractID: contractID,
			}).Return(policySet, getPolicyError).Times(times)
		}

		expectPolicySetUpdate = func(t *testing.T, client *imaging.Mock, contractID, policySetID, name, region string, updatePolicySetError error) {
			client.On("UpdatePolicySet", mock.Anything, imaging.UpdatePolicySetRequest{
				PolicySetID: policySetID,
				ContractID:  contractID,
				UpdatePolicySet: imaging.UpdatePolicySet{
					Name:   name,
					Region: imaging.Region(region),
				},
			}).Return(nil, updatePolicySetError).Once()
		}
		expectPolicySetDelete = func(t *testing.T, client *imaging.Mock, contractID, policySetID, network string, listPolicyResponse *imaging.ListPoliciesResponse, listPolicyError, deletePolicySetError error) {
			client.On("ListPolicies", mock.Anything, imaging.ListPoliciesRequest{
				Network:     imaging.PolicyNetworkProduction,
				ContractID:  contractID,
				PolicySetID: policySetID,
			}).Return(listPolicyResponse, listPolicyError).Once()

			if listPolicyError != nil {
				return
			}

			client.On("ListPolicies", mock.Anything, imaging.ListPoliciesRequest{
				Network:     imaging.PolicyNetworkStaging,
				ContractID:  contractID,
				PolicySetID: policySetID,
			}).Return(listPolicyResponse, listPolicyError).Once()

			client.On("DeletePolicySet", mock.Anything, imaging.DeletePolicySetRequest{
				PolicySetID: policySetID,
				ContractID:  contractID,
			}).Return(deletePolicySetError).Once()
		}
	)

	testDir := "testdata/TestResPolicySet"
	contractID, policySetID, policySetName, mediaType := "1-TEST", "testID", "test_policy_set", string(imaging.TypeImage)
	US, EMEA := "US", "EMEA"
	importStateID := fmt.Sprintf("%s:%s", policySetID, contractID)

	tests := map[string]struct {
		init  func(*imaging.Mock)
		steps []resource.TestStep
	}{
		"ok create": {
			init: func(m *imaging.Mock) {
				createdPolicySet := &imaging.PolicySet{Name: policySetName, ID: policySetID, Region: imaging.Region(EMEA), Type: mediaType}

				// create
				expectPolicySetCreation(t, m, contractID, policySetName, EMEA, mediaType, createdPolicySet, nil)

				expectPolicySetRead(t, m, contractID, policySetID, createdPolicySet, nil, 2)

				// delete
				expectPolicySetDelete(t, m, contractID, policySetID, "", &imaging.ListPoliciesResponse{
					ItemKind: "POLICY",
					Items: []imaging.PolicyOutput{
						&imaging.PolicyOutputImage{ID: ".auto"},
					},
					TotalItems: 1,
				}, nil, nil)
			},
			steps: []resource.TestStep{
				{
					Config: testutils.LoadFixtureString(t, "%s/lifecycle/create.tf", testDir),
					Check: resource.ComposeTestCheckFunc(
						resource.TestCheckResourceAttr("akamai_imaging_policy_set.test_image_set", "id", "testID"),
						resource.TestCheckResourceAttr("akamai_imaging_policy_set.test_image_set", "name", "test_policy_set"),
						resource.TestCheckResourceAttr("akamai_imaging_policy_set.test_image_set", "region", string(imaging.RegionEMEA)),
						resource.TestCheckResourceAttr("akamai_imaging_policy_set.test_image_set", "type", string(imaging.TypeImage)),
						resource.TestCheckResourceAttr("akamai_imaging_policy_set.test_image_set", "contract_id", contractID),
					),
				},
			},
		},
		"ok create with ctr_ prefix in configuration": {
			init: func(m *imaging.Mock) {
				createdPolicySet := &imaging.PolicySet{Name: policySetName, ID: policySetID, Region: imaging.Region(EMEA), Type: mediaType}

				// create
				expectPolicySetCreation(t, m, contractID, policySetName, EMEA, mediaType, createdPolicySet, nil)

				// read
				expectPolicySetRead(t, m, contractID, policySetID, createdPolicySet, nil, 2)

				// delete
				expectPolicySetDelete(t, m, contractID, policySetID, "", &imaging.ListPoliciesResponse{
					ItemKind: "POLICY",
					Items: []imaging.PolicyOutput{
						&imaging.PolicyOutputImage{ID: ".auto"},
					},
					TotalItems: 1,
				}, nil, nil)
			},
			steps: []resource.TestStep{
				{
					Config: testutils.LoadFixtureString(t, "%s/suppress_contract_prefix/create.tf", testDir),
					Check: resource.ComposeTestCheckFunc(
						resource.TestCheckResourceAttr("akamai_imaging_policy_set.test_image_set", "id", "testID"),
						resource.TestCheckResourceAttr("akamai_imaging_policy_set.test_image_set", "name", "test_policy_set"),
						resource.TestCheckResourceAttr("akamai_imaging_policy_set.test_image_set", "region", string(imaging.RegionEMEA)),
						resource.TestCheckResourceAttr("akamai_imaging_policy_set.test_image_set", "type", string(imaging.TypeImage)),
						resource.TestCheckResourceAttr("akamai_imaging_policy_set.test_image_set", "contract_id", "ctr_"+contractID),
					),
				},
			},
		},
		"nok create": {
			init: func(m *imaging.Mock) {
				createdPolicySet := &imaging.PolicySet{Name: policySetName, ID: policySetID, Region: imaging.Region(EMEA), Type: mediaType}

				// create
				expectPolicySetCreation(t, m, contractID, policySetName, EMEA, mediaType, createdPolicySet, anError)
			},
			steps: []resource.TestStep{
				{
					Config:      testutils.LoadFixtureString(t, "%s/lifecycle/create.tf", testDir),
					ExpectError: regexp.MustCompile(anError.Error()),
				},
			},
		},
		"nok get policy set post create": {
			init: func(m *imaging.Mock) {
				createdPolicySet := &imaging.PolicySet{Name: policySetName, ID: policySetID, Region: imaging.Region(EMEA), Type: mediaType}

				// create
				expectPolicySetCreation(t, m, contractID, policySetName, EMEA, mediaType, createdPolicySet, nil)

				// create -> read
				expectPolicySetRead(t, m, contractID, policySetID, nil, anError, 1)

				// delete
				expectPolicySetDelete(t, m, contractID, policySetID, "", &imaging.ListPoliciesResponse{
					ItemKind: "POLICY",
					Items: []imaging.PolicyOutput{
						&imaging.PolicyOutputImage{ID: ".auto"},
					},
					TotalItems: 1,
				}, nil, nil)
			},
			steps: []resource.TestStep{
				{
					Config:      testutils.LoadFixtureString(t, "%s/lifecycle/create.tf", testDir),
					ExpectError: regexp.MustCompile(anError.Error()),
				},
			},
		},
		"ok create update": {
			init: func(m *imaging.Mock) {
				createdPolicySet := &imaging.PolicySet{Name: policySetName, ID: policySetID, Region: imaging.Region(EMEA), Type: mediaType}
				updatedPolicySet := &imaging.PolicySet{Name: policySetName, ID: policySetID, Region: imaging.Region(US), Type: mediaType}

				// create
				expectPolicySetCreation(t, m, contractID, policySetName, EMEA, mediaType, createdPolicySet, nil)

				// create -> read, test -> read, refresh
				expectPolicySetRead(t, m, contractID, policySetID, createdPolicySet, nil, 3)

				// update
				expectPolicySetUpdate(t, m, contractID, policySetID, policySetName, US, nil)

				// update -> read
				expectPolicySetRead(t, m, contractID, policySetID, updatedPolicySet, nil, 2)

				// delete
				expectPolicySetDelete(t, m, contractID, policySetID, "", &imaging.ListPoliciesResponse{
					ItemKind: "POLICY",
					Items: []imaging.PolicyOutput{
						&imaging.PolicyOutputImage{ID: ".auto"},
					},
					TotalItems: 1,
				}, nil, nil)
			},
			steps: []resource.TestStep{
				{
					Config: testutils.LoadFixtureString(t, "%s/lifecycle/create.tf", testDir),
					Check: resource.ComposeTestCheckFunc(
						resource.TestCheckResourceAttr("akamai_imaging_policy_set.test_image_set", "id", "testID"),
						resource.TestCheckResourceAttr("akamai_imaging_policy_set.test_image_set", "name", "test_policy_set"),
						resource.TestCheckResourceAttr("akamai_imaging_policy_set.test_image_set", "region", string(imaging.RegionEMEA)),
						resource.TestCheckResourceAttr("akamai_imaging_policy_set.test_image_set", "type", string(imaging.TypeImage)),
						resource.TestCheckResourceAttr("akamai_imaging_policy_set.test_image_set", "contract_id", contractID),
					),
				},
				{
					Config: testutils.LoadFixtureString(t, "%s/lifecycle/update_region_us.tf", testDir),
					Check: resource.ComposeTestCheckFunc(
						resource.TestCheckResourceAttr("akamai_imaging_policy_set.test_image_set", "id", "testID"),
						resource.TestCheckResourceAttr("akamai_imaging_policy_set.test_image_set", "name", "test_policy_set"),
						resource.TestCheckResourceAttr("akamai_imaging_policy_set.test_image_set", "region", string(imaging.RegionUS)),
						resource.TestCheckResourceAttr("akamai_imaging_policy_set.test_image_set", "type", string(imaging.TypeImage)),
						resource.TestCheckResourceAttr("akamai_imaging_policy_set.test_image_set", "contract_id", contractID),
					),
				},
			},
		},
		"ok create, update with removed prefix": {
			init: func(m *imaging.Mock) {
				createdPolicySet := &imaging.PolicySet{Name: policySetName, ID: policySetID, Region: imaging.Region(EMEA), Type: mediaType}
				//updatedPolicySet := &imaging.PolicySet{Name: policySetName, ID: policySetID, Region: imaging.Region(US), Type: mediaType}

				// create
				expectPolicySetCreation(t, m, contractID, policySetName, EMEA, mediaType, createdPolicySet, nil)

				// create -> read, test -> read, refresh
				expectPolicySetRead(t, m, contractID, policySetID, createdPolicySet, nil, 3)

				// read after diff suppress
				expectPolicySetRead(t, m, contractID, policySetID, createdPolicySet, nil, 1)

				// delete
				expectPolicySetDelete(t, m, contractID, policySetID, "", &imaging.ListPoliciesResponse{
					ItemKind: "POLICY",
					Items: []imaging.PolicyOutput{
						&imaging.PolicyOutputImage{ID: ".auto"},
					},
					TotalItems: 1,
				}, nil, nil)
			},
			steps: []resource.TestStep{
				{
					Config: testutils.LoadFixtureString(t, "%s/suppress_contract_prefix/create.tf", testDir),
					Check: resource.ComposeTestCheckFunc(
						resource.TestCheckResourceAttr("akamai_imaging_policy_set.test_image_set", "id", "testID"),
						resource.TestCheckResourceAttr("akamai_imaging_policy_set.test_image_set", "name", "test_policy_set"),
						resource.TestCheckResourceAttr("akamai_imaging_policy_set.test_image_set", "region", string(imaging.RegionEMEA)),
						resource.TestCheckResourceAttr("akamai_imaging_policy_set.test_image_set", "type", string(imaging.TypeImage)),
						resource.TestCheckResourceAttr("akamai_imaging_policy_set.test_image_set", "contract_id", "ctr_"+contractID),
					),
				},
				{
					Config: testutils.LoadFixtureString(t, "%s/suppress_contract_prefix/update.tf", testDir),
					Check: resource.ComposeTestCheckFunc(
						resource.TestCheckResourceAttr("akamai_imaging_policy_set.test_image_set", "id", "testID"),
						resource.TestCheckResourceAttr("akamai_imaging_policy_set.test_image_set", "name", "test_policy_set"),
						resource.TestCheckResourceAttr("akamai_imaging_policy_set.test_image_set", "region", string(imaging.RegionEMEA)),
						resource.TestCheckResourceAttr("akamai_imaging_policy_set.test_image_set", "type", string(imaging.TypeImage)),
						resource.TestCheckResourceAttr("akamai_imaging_policy_set.test_image_set", "contract_id", "ctr_"+contractID),
					),
				},
			},
		},
		"test import": {
			init: func(m *imaging.Mock) {
				createdPolicySet := &imaging.PolicySet{Name: policySetName, ID: policySetID, Region: imaging.Region(EMEA), Type: mediaType}

				// create
				expectPolicySetCreation(t, m, contractID, policySetName, EMEA, mediaType, createdPolicySet, nil)

				expectPolicySetRead(t, m, contractID, policySetID, createdPolicySet, nil, 3)

				// delete
				expectPolicySetDelete(t, m, contractID, policySetID, "", &imaging.ListPoliciesResponse{
					ItemKind: "POLICY",
					Items: []imaging.PolicyOutput{
						&imaging.PolicyOutputImage{ID: ".auto"},
					},
					TotalItems: 1,
				}, nil, nil)
			},
			steps: []resource.TestStep{
				{
					Config: testutils.LoadFixtureString(t, "%s/lifecycle/create.tf", testDir),
				},
				{
					ImportState:       true,
					ImportStateId:     importStateID,
					ResourceName:      "akamai_imaging_policy_set.test_image_set",
					ImportStateVerify: true,
				},
			},
		},
		"test import - invalid ID": {
			init: func(m *imaging.Mock) {
				createdPolicySet := &imaging.PolicySet{Name: policySetName, ID: policySetID, Region: imaging.Region(EMEA), Type: mediaType}

				// create
				expectPolicySetCreation(t, m, contractID, policySetName, EMEA, mediaType, createdPolicySet, nil)

				expectPolicySetRead(t, m, contractID, policySetID, createdPolicySet, nil, 2)

				// delete
				expectPolicySetDelete(t, m, contractID, policySetID, "", &imaging.ListPoliciesResponse{
					ItemKind: "POLICY",
					Items: []imaging.PolicyOutput{
						&imaging.PolicyOutputImage{ID: ".auto"},
					},
					TotalItems: 1,
				}, nil, nil)
			},
			steps: []resource.TestStep{
				{
					Config: testutils.LoadFixtureString(t, "%s/lifecycle/create.tf", testDir),
				},
				{
					ImportState:       true,
					ImportStateId:     "DevExpTest",
					ResourceName:      "akamai_imaging_policy_set.test_image_set",
					ImportStateVerify: true,
					ExpectError:       regexp.MustCompile("colon-separated list of policy set ID and contract ID has to be supplied in import: DevExpTest"),
				},
			},
		},
	}

	for name, test := range tests {
		t.Run(name, func(t *testing.T) {
			client := &imaging.Mock{}
			test.init(client)
			useClient(client, func() {
				resource.UnitTest(t, resource.TestCase{
					ProtoV6ProviderFactories: testutils.NewProtoV6ProviderFactory(NewSubprovider()),
					Steps:                    test.steps,
				})
			})
			client.AssertExpectations(t)
		})
	}
}

func Test_filterRemainingPolicies(t *testing.T) {
	tests := map[string]struct {
		input          *imaging.ListPoliciesResponse
		expectedOutput int
	}{
		"just 1 image policy .auto remaining policy": {
			input: &imaging.ListPoliciesResponse{
				ItemKind: "POLICY",
				Items: []imaging.PolicyOutput{
					&imaging.PolicyOutputImage{ID: ".auto"},
				},
				TotalItems: 0,
			},
			expectedOutput: 0,
		},
		"just 1 video policy .auto remaining policy": {
			input: &imaging.ListPoliciesResponse{
				ItemKind: "POLICY",
				Items: []imaging.PolicyOutput{
					&imaging.PolicyOutputVideo{ID: ".auto"},
				},
				TotalItems: 0,
			},
			expectedOutput: 0,
		},
		"2 video policies, one of them .auto": {
			input: &imaging.ListPoliciesResponse{
				ItemKind: "POLICY",
				Items: []imaging.PolicyOutput{
					&imaging.PolicyOutputVideo{ID: ".auto"},
					&imaging.PolicyOutputVideo{ID: "not-auto"},
				},
				TotalItems: 0,
			},
			expectedOutput: 1,
		},
	}
	for name, test := range tests {
		t.Run(name, func(t *testing.T) {
			assert.Equal(t, test.expectedOutput, filterRemainingPolicies(test.input))
		})
	}

}
