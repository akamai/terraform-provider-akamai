package property

import (
	"fmt"
	"regexp"
	"strconv"
	"testing"

	"github.com/akamai/AkamaiOPEN-edgegrid-golang/v3/pkg/papi"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/resource"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/mock"
)

func TestDataPropertyIncludeActivation(t *testing.T) {
	tests := map[string]struct {
		attrs      includeActivationTestAttributes
		init       func(*testing.T, *papi.Mock, includeActivationTestAttributes)
		configPath string
		error      *regexp.Regexp
	}{
		"happy path - latest active activation - multiple includes - STAGING": {
			attrs: includeActivationTestAttributes{
				contractID:          contractForTests,
				groupID:             groupForTests,
				includeID:           includeForTests,
				network:             stagingNetwork,
				activationsResponse: *createIncludeActivationsResponse(accountForTests, contractForTests, groupForTests, includeActivationsForTests),
			},
			init: func(t *testing.T, m *papi.Mock, attrs includeActivationTestAttributes) {
				mockListIncludeActivation(m, attrs, 5)
			},
			configPath: "testdata/TestDataPropertyIncludeActivation/valid_staging.tf",
		},
		"happy path - latest active activation - multiple includes - PRODUCTION": {
			attrs: includeActivationTestAttributes{
				contractID:          contractForTests,
				groupID:             groupForTests,
				includeID:           includeForTests,
				network:             productionNetwork,
				activationsResponse: *createIncludeActivationsResponse(accountForTests, contractForTests, groupForTests, includeActivationsForTests),
			},
			init: func(t *testing.T, m *papi.Mock, attrs includeActivationTestAttributes) {
				mockListIncludeActivation(m, attrs, 5)
			},
			configPath: "testdata/TestDataPropertyIncludeActivation/valid_production.tf",
		},
		"no latest activation for provided network": {
			attrs: includeActivationTestAttributes{
				contractID:          contractForTests,
				groupID:             groupForTests,
				includeID:           includeForTests,
				network:             productionNetwork,
				activationsResponse: *createIncludeActivationsResponse(accountForTests, contractForTests, groupForTests, includeActivationsForTests[:2]),
			},
			init: func(t *testing.T, m *papi.Mock, attrs includeActivationTestAttributes) {
				mockListIncludeActivation(m, attrs, 5)
			},
			configPath: "testdata/TestDataPropertyIncludeActivation/no_activation_for_given_network.tf",
		},
		"latest activation of type `DEACTIVATE` - no latest activation": {
			attrs: includeActivationTestAttributes{
				contractID:          contractForTests,
				groupID:             groupForTests,
				includeID:           includeForTests,
				network:             stagingNetwork,
				activationsResponse: *createIncludeActivationsResponse(accountForTests, contractForTests, groupForTests, includeActivationsWithLatestDeactivate),
			},
			init: func(t *testing.T, m *papi.Mock, attrs includeActivationTestAttributes) {
				mockListIncludeActivation(m, attrs, 5)
			},
			configPath: "testdata/TestDataPropertyIncludeActivation/valid_staging.tf",
		},
		"no `ACTIVE` activations - no latest activation": {
			attrs: includeActivationTestAttributes{
				contractID:          contractForTests,
				groupID:             groupForTests,
				includeID:           includeForTests,
				network:             productionNetwork,
				activationsResponse: *createIncludeActivationsResponse(accountForTests, contractForTests, groupForTests, includeActivationsForTests[:2]),
			},
			init: func(t *testing.T, m *papi.Mock, attrs includeActivationTestAttributes) {
				mockListIncludeActivation(m, attrs, 5)
			},
			configPath: "testdata/TestDataPropertyIncludeActivation/valid_production.tf",
		},
		"no emails in include activation": {
			attrs: includeActivationTestAttributes{
				contractID: contractForTests,
				groupID:    groupForTests,
				includeID:  includeForTests,
				network:    stagingNetwork,
				activationsResponse: *createIncludeActivationsResponse(accountForTests, contractForTests, groupForTests, []papi.IncludeActivation{
					createIncludeActivation(includeActivationData{
						network:             productionNetwork,
						activationType:      activationTypeActivate,
						activationID:        "1",
						status:              activationStatusActive,
						note:                "Note 1",
						updateDate:          "2022-11-12T14:15:27Z",
						includeID:           "1",
						includeName:         "Name 1",
						includeType:         includeTypeCommonSettings,
						includeVersion:      1,
						includeActivationID: "1",
						emails:              nil,
					}),
				}),
			},
			init: func(t *testing.T, m *papi.Mock, attrs includeActivationTestAttributes) {
				mockListIncludeActivation(m, attrs, 5)
			},
			configPath: "testdata/TestDataPropertyIncludeActivation/valid_staging.tf",
		},
		"required attribute missing - contract_id": {
			attrs:      includeActivationTestAttributes{},
			init:       func(t *testing.T, m *papi.Mock, attrs includeActivationTestAttributes) {},
			configPath: "testdata/TestDataPropertyIncludeActivation/no_contract_id.tf",
			error:      regexp.MustCompile("Error: Missing required argument"),
		},
		"required attribute missing - group_id": {
			attrs:      includeActivationTestAttributes{},
			init:       func(t *testing.T, m *papi.Mock, attrs includeActivationTestAttributes) {},
			configPath: "testdata/TestDataPropertyIncludeActivation/no_group_id.tf",
			error:      regexp.MustCompile("Error: Missing required argument"),
		},
		"required attribute missing - include_id": {
			attrs:      includeActivationTestAttributes{},
			init:       func(t *testing.T, m *papi.Mock, attrs includeActivationTestAttributes) {},
			configPath: "testdata/TestDataPropertyIncludeActivation/no_include_id.tf",
			error:      regexp.MustCompile("Error: Missing required argument"),
		},
		"required attribute missing - network": {
			attrs:      includeActivationTestAttributes{},
			init:       func(t *testing.T, m *papi.Mock, attrs includeActivationTestAttributes) {},
			configPath: "testdata/TestDataPropertyIncludeActivation/no_network.tf",
			error:      regexp.MustCompile("Error: Missing required argument"),
		},
		"ListIncludeActivations - API error": {
			attrs: includeActivationTestAttributes{
				contractID:          contractForTests,
				groupID:             groupForTests,
				includeID:           includeForTests,
				network:             stagingNetwork,
				activationsResponse: *createIncludeActivationsResponse(accountForTests, contractForTests, groupForTests, includeActivationsForTests),
			},
			init: func(t *testing.T, m *papi.Mock, attrs includeActivationTestAttributes) {
				m.On("ListIncludeActivations", mock.Anything, papi.ListIncludeActivationsRequest{
					IncludeID:  attrs.includeID,
					ContractID: attrs.contractID,
					GroupID:    attrs.groupID,
				}).Return(nil, fmt.Errorf("could not list include activations"))
			},
			configPath: "testdata/TestDataPropertyIncludeActivation/valid_staging.tf",
			error:      regexp.MustCompile("could not list include activations"),
		},
	}

	for name, test := range tests {
		t.Run(name, func(t *testing.T) {
			client := &papi.Mock{}
			test.init(t, client, test.attrs)
			useClient(client, nil, func() {
				resource.UnitTest(t, resource.TestCase{
					Providers:  testAccProviders,
					IsUnitTest: true,
					Steps: []resource.TestStep{{
						Config:      loadFixtureString(test.configPath),
						Check:       checkPropertyIncludeActivationAttrs(test.attrs),
						ExpectError: test.error,
					}},
				})
			})
			client.AssertExpectations(t)
		})
	}
}

func TestFilterActivations(t *testing.T) {
	t.Run("test filterActivation - STAGING network", func(t *testing.T) {
		stagingIncludes := filterIncludeActivationsByNetwork(includeActivationsForTests, stagingNetwork)
		assert.Equal(t, 2, len(stagingIncludes))
		assert.Equal(t, "1", stagingIncludes[0].IncludeID)
		assert.Equal(t, "2", stagingIncludes[1].IncludeID)
	})
	t.Run("test filterActivation - PRODUCTION network", func(t *testing.T) {
		productionIncludes := filterIncludeActivationsByNetwork(includeActivationsForTests, productionNetwork)
		assert.Equal(t, 2, len(productionIncludes))
		assert.Equal(t, "3", productionIncludes[0].IncludeID)
		assert.Equal(t, "4", productionIncludes[1].IncludeID)
	})
}

func TestFindLatestActivation(t *testing.T) {
	activationWithLatestDeactivate := append(includeActivationsForTests, includeActivationTypeDeactivate)
	noActiveStatusActivations := includeActivationsForTests[:2]
	manyActivations := append(includeActivationsForTests, createIncludeActivation(includeActivationData{
		network:             stagingNetwork,
		activationType:      activationTypeActivate,
		activationID:        "6",
		status:              activationStatusPendingCancellation,
		note:                "Note 6",
		updateDate:          "2022-11-01T14:15:27Z",
		includeID:           "6",
		includeName:         "Name 6",
		includeType:         includeTypeCommonSettings,
		includeVersion:      6,
		includeActivationID: "6",
	}))
	manyActivations = append(manyActivations, createIncludeActivation(includeActivationData{
		network:             productionNetwork,
		activationType:      activationTypeActivate,
		activationID:        "7",
		status:              activationStatusActive,
		note:                "Note 7",
		updateDate:          "2022-11-12T14:15:27Z",
		includeID:           "7",
		includeName:         "Name 7",
		includeType:         includeTypeCommonSettings,
		includeVersion:      7,
		includeActivationID: "7",
	}))
	manyActivations = append(manyActivations, includeActivationTypeDeactivate)

	tests := map[string]struct {
		activations            []papi.IncludeActivation
		actualLatestActivation *papi.IncludeActivation
	}{
		"many different activations": {
			activations:            manyActivations,
			actualLatestActivation: &manyActivations[0],
		},
		"oldest activation active": {
			activations:            includeActivationsForTests,
			actualLatestActivation: &includeActivationsForTests[2],
		},
		"activate and deactivate types with active status - the latest one `wins` (nil)": {
			activations:            includeActivationsBothTypesActive,
			actualLatestActivation: nil,
		},
		"no activations": {
			activations:            []papi.IncludeActivation{},
			actualLatestActivation: nil,
		},
		"latest activation of type DEACTIVATE": {
			activations:            activationWithLatestDeactivate,
			actualLatestActivation: nil,
		},
		"no latest activation - no ACTIVE status": {
			activations:            noActiveStatusActivations,
			actualLatestActivation: nil,
		},
	}

	for name, test := range tests {
		t.Run(name, func(t *testing.T) {
			latestActivation, _ := findLatestIncludeActivation(test.activations)
			if test.actualLatestActivation == nil {
				assert.Equal(t, (*papi.IncludeActivation)(nil), latestActivation)
			} else {
				assert.Equal(t, test.actualLatestActivation.IncludeID, latestActivation.IncludeID)
			}
		})
	}
}

// includeActivationAttributes contains data used for data source creation and ListIncludesActivation response
type includeActivationTestAttributes struct {
	contractID, groupID, includeID, network string
	activationsResponse                     papi.ListIncludeActivationsResponse
}

// includeActivationData holds data used for creation of include activation
type includeActivationData struct {
	network             string
	activationType      string
	activationID        string
	status              string
	note                string
	updateDate          string
	includeID           string
	includeName         string
	includeType         string
	includeVersion      int
	includeActivationID string
	emails              []string
}

var (
	productionNetwork                   = "PRODUCTION"
	stagingNetwork                      = "STAGING"
	activationTypeDeactivate            = "DEACTIVATE"
	activationTypeActivate              = "ACTIVATE"
	activationStatusActive              = "ACTIVE"
	activationStatusInactive            = "INACTIVE"
	activationStatusPending             = "PENDING"
	activationStatusPendingCancellation = "PENDING_CANCELLATION"
	includeTypeMicroServices            = "MICROSERVICES"
	includeTypeCommonSettings           = "COMMON_SETTINGS"
	includeForTests                     = "inc_1"
	accountForTests                     = "acc_1"

	// includeActivationTypeDeactivate is an include activation of type `DEACTIVATE` with `ACTIVE` status
	includeActivationTypeDeactivate = createIncludeActivation(includeActivationData{
		network:             productionNetwork,
		activationType:      activationTypeDeactivate,
		activationID:        "5",
		status:              activationStatusActive,
		note:                "Note 5",
		updateDate:          "2022-11-10T14:15:27Z",
		includeID:           "5",
		includeName:         "Name 5",
		includeType:         includeTypeCommonSettings,
		includeVersion:      5,
		includeActivationID: "5",
		emails:              nil,
	})

	// includeActivationsWithLatestDeactivate is a slice of include activations, with the latest activation being of type
	// `DEACTIVATE` with status `ACTIVE`
	includeActivationsWithLatestDeactivate = append(includeActivationsForTests, includeActivationTypeDeactivate)

	// includeActivationsBothTypesActive is a slice of two includes, one being of type `ACTIVATE` and one being of type
	// `DEACTIVATE` with both having status `ACTIVE`. The latest one is the one with the most recent updateDate.
	includeActivationsBothTypesActive = []papi.IncludeActivation{
		includeActivationTypeDeactivate,
		createIncludeActivation(includeActivationData{
			network:             stagingNetwork,
			activationType:      activationTypeActivate,
			activationID:        "2",
			status:              activationStatusActive,
			note:                "Note 2",
			updateDate:          "2022-11-09T14:15:27Z",
			includeID:           "2",
			includeName:         "Name 2",
			includeType:         includeTypeMicroServices,
			includeVersion:      2,
			includeActivationID: "2",
			emails:              nil,
		}),
	}

	// includeActivationsForTests is a slice of different include activations used in tests
	includeActivationsForTests = []papi.IncludeActivation{
		createIncludeActivation(includeActivationData{
			network:             stagingNetwork,
			activationType:      activationTypeActivate,
			activationID:        "1",
			status:              activationStatusInactive,
			note:                "Note 1",
			updateDate:          "2022-11-08T14:15:27Z",
			includeID:           "1",
			includeName:         "Name 1",
			includeType:         includeTypeMicroServices,
			includeVersion:      1,
			includeActivationID: "1",
		}),
		createIncludeActivation(includeActivationData{
			network:             stagingNetwork,
			activationType:      activationTypeActivate,
			activationID:        "2",
			status:              activationStatusInactive,
			note:                "Note 2",
			updateDate:          "2022-11-09T14:15:27Z",
			includeID:           "2",
			includeName:         "Name 2",
			includeType:         includeTypeMicroServices,
			includeVersion:      2,
			includeActivationID: "2",
		}),
		createIncludeActivation(includeActivationData{
			network:             productionNetwork,
			activationType:      activationTypeActivate,
			activationID:        "3",
			status:              activationStatusPending,
			note:                "Note 3",
			updateDate:          "2022-11-06T14:15:27Z",
			includeID:           "3",
			includeName:         "Name 3",
			includeType:         includeTypeMicroServices,
			includeVersion:      3,
			includeActivationID: "3",
		}),
		createIncludeActivation(includeActivationData{
			network:             productionNetwork,
			activationType:      activationTypeActivate,
			activationID:        "4",
			status:              activationStatusActive,
			note:                "Note 4",
			updateDate:          "2022-11-07T14:15:27Z",
			includeID:           "4",
			includeName:         "Name 4",
			includeType:         includeTypeMicroServices,
			includeVersion:      4,
			includeActivationID: "4",
		}),
	}

	// mockListIncludeActivation mocks ListIncludeActivation call with provided parameters
	mockListIncludeActivation = func(m *papi.Mock, attrs includeActivationTestAttributes, timesToRun int) {
		m.On("ListIncludeActivations", mock.Anything, papi.ListIncludeActivationsRequest{
			IncludeID:  attrs.includeID,
			ContractID: attrs.contractID,
			GroupID:    attrs.groupID,
		}).Return(&attrs.activationsResponse, nil).Times(timesToRun)
	}
)

// createIncludeActivation creates include activation based on provided parameters
func createIncludeActivation(data includeActivationData) papi.IncludeActivation {
	return papi.IncludeActivation{
		ActivationID:        data.activationID,
		Network:             papi.ActivationNetwork(data.network),
		ActivationType:      papi.ActivationType(data.activationType),
		Status:              papi.ActivationStatus(data.status),
		SubmitDate:          data.updateDate,
		UpdateDate:          data.updateDate,
		Note:                data.note,
		NotifyEmails:        data.emails,
		IncludeID:           data.includeID,
		IncludeName:         data.includeName,
		IncludeType:         papi.IncludeType(data.includeType),
		IncludeVersion:      data.includeVersion,
		IncludeActivationID: data.includeActivationID,
	}
}

// createIncludeActivationsResponse creates response from ListIncludeActivations call with provided list of include activations
func createIncludeActivationsResponse(accountID, contractID, groupID string, includes []papi.IncludeActivation) *papi.ListIncludeActivationsResponse {
	return &papi.ListIncludeActivationsResponse{
		AccountID:  accountID,
		ContractID: contractID,
		GroupID:    groupID,
		Activations: papi.IncludeActivationsRes{
			Items: includes,
		},
	}
}

// checkPropertyIncludeActivationAttrs create check functions for a data source based on provided parameters
func checkPropertyIncludeActivationAttrs(data includeActivationTestAttributes) resource.TestCheckFunc {
	var testCheckFuncs []resource.TestCheckFunc
	dataSourceID := data.includeID + ":" + data.network
	filteredActivations := filterIncludeActivationsByNetwork(data.activationsResponse.Activations.Items, data.network)
	latestActivation, _ := findLatestIncludeActivation(filteredActivations)

	testCheckFuncs = append(testCheckFuncs, resource.TestCheckResourceAttr("data.akamai_property_include_activation.test", "id", dataSourceID))
	testCheckFuncs = append(testCheckFuncs, resource.TestCheckResourceAttr("data.akamai_property_include_activation.test", "contract_id", data.contractID))
	testCheckFuncs = append(testCheckFuncs, resource.TestCheckResourceAttr("data.akamai_property_include_activation.test", "group_id", data.groupID))
	testCheckFuncs = append(testCheckFuncs, resource.TestCheckResourceAttr("data.akamai_property_include_activation.test", "include_id", data.includeID))
	testCheckFuncs = append(testCheckFuncs, resource.TestCheckResourceAttr("data.akamai_property_include_activation.test", "network", data.network))

	var version, name, note string
	emailLength := "0"

	if latestActivation != nil {
		version = strconv.Itoa(latestActivation.IncludeVersion)
		name = latestActivation.IncludeName
		note = latestActivation.Note
		emailLength = strconv.Itoa(len(latestActivation.NotifyEmails))
		for i, email := range latestActivation.NotifyEmails {
			testCheckFuncs = append(testCheckFuncs, resource.TestCheckResourceAttr("data.akamai_property_include_activation.test", fmt.Sprintf("notify_emails.%d", i), email))
		}
	}

	testCheckFuncs = append(testCheckFuncs, resource.TestCheckResourceAttr("data.akamai_property_include_activation.test", "notify_emails.#", emailLength))
	testCheckFuncs = append(testCheckFuncs, resource.TestCheckResourceAttr("data.akamai_property_include_activation.test", "version", version))
	testCheckFuncs = append(testCheckFuncs, resource.TestCheckResourceAttr("data.akamai_property_include_activation.test", "name", name))
	testCheckFuncs = append(testCheckFuncs, resource.TestCheckResourceAttr("data.akamai_property_include_activation.test", "note", note))

	return resource.ComposeAggregateTestCheckFunc(testCheckFuncs...)
}
