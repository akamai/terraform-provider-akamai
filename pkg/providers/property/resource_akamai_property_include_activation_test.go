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
	"github.com/stretchr/testify/require"
)

var (
	includeID         = "inc_12345"
	contractID        = "ctr_test_contract"
	groupID           = "grp_test_group"
	accountID         = "test_account"
	note              = "test activation"
	email             = "jbond@example.com"
	networkStaging    = "STAGING"
	networkProduction = "PRODUCTION"
	version           = 3
	testDir           = "testdata/TestResPropertyIncludeActivation"
)

func TestResourcePropertyIncludeActivation(t *testing.T) {
	type attrs struct {
		includeID, contractID, groupID, network, note string
		version                                       int
		notifyEmails                                  []string
		autoAcknowledgeRuleWarnings                   bool
	}

	var (
		activateIncludeReq = func(network string, acknowledgeAllWarnings bool) papi.ActivateIncludeRequest {
			return papi.ActivateIncludeRequest{
				IncludeID:              includeID,
				Version:                3,
				Network:                papi.ActivationNetwork(network),
				Note:                   note,
				NotifyEmails:           []string{email},
				AcknowledgeAllWarnings: acknowledgeAllWarnings,
			}
		}

		activateIncludeRes = papi.ActivationIncludeResponse{
			ActivationID:   "temporary-activation-id",
			ActivationLink: "/papi/v1/includes/inc_12345/activations/temporary-activation-id",
		}

		expectActivateIncludeOnStaging = func(client *papi.Mock, network string, acknowledgeAllWarnings bool) *papi.ActivationIncludeResponse {
			activateIncludeReq := activateIncludeReq(network, acknowledgeAllWarnings)
			activateIncludeRes := activateIncludeRes
			client.On("ActivateInclude", mock.Anything, activateIncludeReq).Return(&activateIncludeRes, nil).Once()
			return &activateIncludeRes
		}

		expectActivateIncludeOnProduction = func(client *papi.Mock, network string, acknowledgeAllWarnings bool) *papi.ActivationIncludeResponse {
			activateIncludeReq := activateIncludeReq(network, acknowledgeAllWarnings)
			activateIncludeReq.ComplianceRecord = &papi.ComplianceRecordOther{}
			activateIncludeRes := activateIncludeRes
			client.On("ActivateInclude", mock.Anything, activateIncludeReq).Return(&activateIncludeRes, nil).Once()
			return &activateIncludeRes
		}

		getIncludeActivationReq = func(actID string) papi.GetIncludeActivationRequest {
			return papi.GetIncludeActivationRequest{
				IncludeID:    includeID,
				ActivationID: actID,
			}
		}

		expectGetTempIncludeActivation = func(client *papi.Mock, tempActID string, network papi.ActivationNetwork) {
			getIncludeActivationReq := getIncludeActivationReq(tempActID)
			getIncludeActivationRes := papi.GetIncludeActivationResponse{
				AccountID:  accountID,
				ContractID: contractID,
				GroupID:    groupID,
				Activation: papi.IncludeActivation{
					ActivationID:   tempActID,
					Network:        network,
					ActivationType: papi.ActivationTypeActivate,
					Status:         papi.ActivationStatusActive,
					SubmitDate:     "2022-10-27T10:21:40Z",
					UpdateDate:     "2022-10-27T10:22:54Z",
					Note:           note,
					NotifyEmails:   []string{email},
					IncludeID:      includeID,
					IncludeVersion: version,
				},
			}
			client.On("GetIncludeActivation", mock.Anything, getIncludeActivationReq).Return(&getIncludeActivationRes, nil).Once()
		}

		expectGetIncludeActivation = func(client *papi.Mock, actID string, network papi.ActivationNetwork) *papi.GetIncludeActivationResponse {
			getIncludeActivationReq := getIncludeActivationReq(actID)
			getIncludeActivationRes := papi.GetIncludeActivationResponse{
				AccountID:  accountID,
				ContractID: contractID,
				GroupID:    groupID,
				Activation: papi.IncludeActivation{
					ActivationID:   actID,
					Network:        network,
					ActivationType: papi.ActivationTypeActivate,
					Status:         papi.ActivationStatusActive,
					SubmitDate:     "2022-10-27T11:21:40Z",
					UpdateDate:     "2022-10-27T11:22:54Z",
					Note:           note,
					NotifyEmails:   []string{email},
					IncludeID:      includeID,
					IncludeVersion: version,
				},
			}
			client.On("GetIncludeActivation", mock.Anything, getIncludeActivationReq).Return(&getIncludeActivationRes, nil).Once()
			return &getIncludeActivationRes
		}

		expectGetIncludeActivationOnProduction = func(client *papi.Mock, actID string) *papi.GetIncludeActivationResponse {
			getIncludeActivationReq := getIncludeActivationReq(actID)
			getIncludeActivationRes := papi.GetIncludeActivationResponse{
				AccountID:  accountID,
				ContractID: contractID,
				GroupID:    groupID,
				Activation: papi.IncludeActivation{
					ActivationID:   actID,
					Network:        papi.ActivationNetworkProduction,
					ActivationType: papi.ActivationTypeActivate,
					Status:         papi.ActivationStatusActive,
					SubmitDate:     "2022-10-28T11:21:40Z",
					UpdateDate:     "2022-10-28T11:22:54Z",
					Note:           note,
					NotifyEmails:   []string{email},
					IncludeID:      includeID,
					IncludeVersion: version,
				},
			}
			client.On("GetIncludeActivation", mock.Anything, getIncludeActivationReq).Return(&getIncludeActivationRes, nil).Once()
			return &getIncludeActivationRes
		}

		listActivationsReq = papi.ListIncludeActivationsRequest{
			IncludeID:  includeID,
			GroupID:    groupID,
			ContractID: contractID,
		}

		activations = papi.IncludeActivationsRes{
			Items: []papi.IncludeActivation{
				{
					ActivationID:   "atv_12345",
					Network:        papi.ActivationNetworkStaging,
					ActivationType: papi.ActivationTypeActivate,
					Status:         papi.ActivationStatusActive,
					SubmitDate:     "2022-10-27T11:21:40Z",
					UpdateDate:     "2022-10-27T11:22:54Z",
					Note:           note,
					NotifyEmails:   []string{email},
					IncludeID:      includeID,
					IncludeVersion: 3,
				},
				{
					ActivationID:   "atv_12344",
					Network:        papi.ActivationNetworkStaging,
					ActivationType: papi.ActivationTypeActivate,
					Status:         papi.ActivationStatusActive,
					SubmitDate:     "2022-10-26T12:37:49Z",
					UpdateDate:     "2022-10-26T12:38:59Z",
					Note:           note,
					NotifyEmails:   []string{email},
					IncludeID:      includeID,
					IncludeVersion: 2,
				},
				{
					ActivationID:   "atv_12343",
					Network:        papi.ActivationNetworkStaging,
					ActivationType: papi.ActivationTypeActivate,
					Status:         papi.ActivationStatusActive,
					SubmitDate:     "2022-08-17T09:13:18Z",
					UpdateDate:     "2022-08-17T09:15:35Z",
					Note:           note,
					NotifyEmails:   []string{email},
					IncludeID:      includeID,
					IncludeVersion: 1,
				},
			},
		}

		expectListIncludeActivations = func(client *papi.Mock) *papi.ListIncludeActivationsResponse {
			activationsRes := papi.ListIncludeActivationsResponse{
				AccountID:   accountID,
				ContractID:  contractID,
				GroupID:     groupID,
				Activations: activations,
			}
			client.On("ListIncludeActivations", mock.Anything, listActivationsReq).Return(&activationsRes, nil).Once()
			return &activationsRes
		}

		expectListIncludeActivationsUpdate = func(client *papi.Mock) *papi.ListIncludeActivationsResponse {
			activations.Items = append(activations.Items, papi.IncludeActivation{
				ActivationID:   "atv_12346",
				Network:        papi.ActivationNetworkProduction,
				ActivationType: papi.ActivationTypeActivate,
				Status:         papi.ActivationStatusActive,
				SubmitDate:     "2022-10-28T11:21:40Z",
				UpdateDate:     "2022-10-28T11:22:54Z",
				Note:           note,
				NotifyEmails:   []string{email},
				IncludeID:      includeID,
				IncludeVersion: version,
			})
			activationsRes := papi.ListIncludeActivationsResponse{
				AccountID:   accountID,
				ContractID:  contractID,
				GroupID:     groupID,
				Activations: activations,
			}
			client.On("ListIncludeActivations", mock.Anything, listActivationsReq).Return(&activationsRes, nil).Once()
			return &activationsRes
		}

		expectDeactivateInclude = func(client *papi.Mock, network papi.ActivationNetwork, acknowledgedWarnings bool) *papi.DeactivationIncludeResponse {
			deactivateIncludeReq := papi.DeactivateIncludeRequest{
				IncludeID:              includeID,
				Version:                version,
				Network:                network,
				Note:                   note,
				NotifyEmails:           []string{email},
				AcknowledgeAllWarnings: acknowledgedWarnings,
			}
			if network == papi.ActivationNetworkProduction {
				deactivateIncludeReq.ComplianceRecord = &papi.ComplianceRecordOther{}
			}
			deactivateIncludeRes := papi.DeactivationIncludeResponse{
				ActivationID:   "temporary-deactivation-id",
				ActivationLink: "/papi/v1/includes/inc_12345/activations/temporary-deactivation-id",
			}
			client.On("DeactivateInclude", mock.Anything, deactivateIncludeReq).Return(&deactivateIncludeRes, nil).Once()
			return &deactivateIncludeRes
		}

		expectGetTempIncludeDeactivation = func(client *papi.Mock, tempDeactID string, network papi.ActivationNetwork) {
			getIncludeActivationReq := getIncludeActivationReq(tempDeactID)
			getIncludeActivationRes := papi.GetIncludeActivationResponse{
				AccountID:  accountID,
				ContractID: contractID,
				GroupID:    groupID,
				Activation: papi.IncludeActivation{
					ActivationID:   tempDeactID,
					Network:        network,
					ActivationType: papi.ActivationTypeDeactivate,
					Status:         papi.ActivationStatusActive,
					SubmitDate:     "2022-10-27T12:21:40Z",
					UpdateDate:     "2022-10-27T12:22:54Z",
					Note:           note,
					NotifyEmails:   []string{email},
					IncludeID:      includeID,
					IncludeVersion: version,
				},
			}
			client.On("GetIncludeActivation", mock.Anything, getIncludeActivationReq).Return(&getIncludeActivationRes, nil).Once()
		}

		checkAttributes = func(attrs attrs) resource.TestCheckFunc {
			checks := []resource.TestCheckFunc{
				resource.TestCheckResourceAttr("akamai_property_include_activation.activation", "include_id", attrs.includeID),
				resource.TestCheckResourceAttr("akamai_property_include_activation.activation", "contract_id", attrs.contractID),
				resource.TestCheckResourceAttr("akamai_property_include_activation.activation", "group_id", attrs.groupID),
				resource.TestCheckResourceAttr("akamai_property_include_activation.activation", "version", strconv.Itoa(attrs.version)),
				resource.TestCheckResourceAttr("akamai_property_include_activation.activation", "network", attrs.network),
				resource.TestCheckResourceAttr("akamai_property_include_activation.activation", "note", attrs.note),
				resource.TestCheckResourceAttr("akamai_property_include_activation.activation", "notify_emails.0", attrs.notifyEmails[0]),
			}

			return resource.ComposeAggregateTestCheckFunc(checks...)
		}
	)

	t.Run("create a new include activation lifecycle", func(t *testing.T) {
		client := new(papi.Mock)

		// create
		actResWithTempID := expectActivateIncludeOnStaging(client, networkStaging, false)
		expectGetTempIncludeActivation(client, actResWithTempID.ActivationID, papi.ActivationNetworkStaging)

		// read
		activations := expectListIncludeActivations(client)
		actID, err := getLatestIncludeActivationID(activations, networkStaging)
		require.NoError(t, err)
		expectGetIncludeActivation(client, actID, papi.ActivationNetworkStaging)

		// refresh
		activations = expectListIncludeActivations(client)
		actID, err = getLatestIncludeActivationID(activations, networkStaging)
		require.NoError(t, err)
		expectGetIncludeActivation(client, actID, papi.ActivationNetworkStaging)

		// destroy
		deactivation := expectDeactivateInclude(client, papi.ActivationNetworkStaging, false)
		expectGetTempIncludeDeactivation(client, deactivation.ActivationID, papi.ActivationNetworkStaging)

		useClient(client, nil, func() {
			resource.UnitTest(t, resource.TestCase{
				Providers: testAccProviders,
				Steps: []resource.TestStep{
					{
						Config: loadFixtureString(fmt.Sprintf("%s/property_include_activation.tf", testDir)),
						Check: checkAttributes(attrs{
							includeID:    includeID,
							contractID:   contractID,
							groupID:      groupID,
							version:      version,
							network:      "STAGING",
							note:         note,
							notifyEmails: []string{email},
						}),
					},
				},
			})
		})
		client.AssertExpectations(t)
	})

	t.Run("update include activation lifecycle", func(t *testing.T) {
		client := new(papi.Mock)

		// 1. first step
		// create
		actResWithTempIDStaging := expectActivateIncludeOnStaging(client, networkStaging, false)
		expectGetTempIncludeActivation(client, actResWithTempIDStaging.ActivationID, papi.ActivationNetworkStaging)

		// read
		activations := expectListIncludeActivations(client)
		actID, err := getLatestIncludeActivationID(activations, networkStaging)
		require.NoError(t, err)
		expectGetIncludeActivation(client, actID, papi.ActivationNetworkStaging)

		// refresh
		activations = expectListIncludeActivations(client)
		actID, err = getLatestIncludeActivationID(activations, networkStaging)
		require.NoError(t, err)
		expectGetIncludeActivation(client, actID, papi.ActivationNetworkStaging)

		// refresh
		activations = expectListIncludeActivations(client)
		actID, err = getLatestIncludeActivationID(activations, networkStaging)
		require.NoError(t, err)
		expectGetIncludeActivation(client, actID, papi.ActivationNetworkStaging)

		// destroy
		deactivation := expectDeactivateInclude(client, papi.ActivationNetworkStaging, false)
		expectGetTempIncludeDeactivation(client, deactivation.ActivationID, papi.ActivationNetworkStaging)

		// 2. second step - network ForceNew
		// create
		actResWithTempIDProduction := expectActivateIncludeOnProduction(client, networkProduction, true)
		expectGetTempIncludeActivation(client, actResWithTempIDProduction.ActivationID, papi.ActivationNetworkProduction)

		// read
		activations = expectListIncludeActivationsUpdate(client)
		actID, err = getLatestIncludeActivationID(activations, networkProduction)
		require.NoError(t, err)
		expectGetIncludeActivationOnProduction(client, actID)

		// destroy
		deactivation = expectDeactivateInclude(client, papi.ActivationNetworkProduction, true)
		expectGetTempIncludeDeactivation(client, deactivation.ActivationID, papi.ActivationNetworkProduction)

		// refresh
		activations = expectListIncludeActivationsUpdate(client)
		actID, err = getLatestIncludeActivationID(activations, networkProduction)
		require.NoError(t, err)
		expectGetIncludeActivationOnProduction(client, actID)

		useClient(client, nil, func() {
			resource.UnitTest(t, resource.TestCase{
				Providers: testAccProviders,
				Steps: []resource.TestStep{
					{
						Config: loadFixtureString(fmt.Sprintf("%s/property_include_activation.tf", testDir)),
						Check: checkAttributes(attrs{
							includeID:    includeID,
							contractID:   contractID,
							groupID:      groupID,
							version:      version,
							network:      "STAGING",
							note:         note,
							notifyEmails: []string{email},
						}),
					},
					{
						Config: loadFixtureString(fmt.Sprintf("%s/property_include_activation_update.tf", testDir)),
						Check: checkAttributes(attrs{
							includeID:    includeID,
							contractID:   contractID,
							groupID:      groupID,
							version:      version,
							network:      "PRODUCTION",
							note:         note,
							notifyEmails: []string{email},
						}),
					},
				},
			})
		})
		client.AssertExpectations(t)
	})

	t.Run("import", func(t *testing.T) {
		client := new(papi.Mock)

		// create
		actResWithTempIDStaging := expectActivateIncludeOnStaging(client, networkStaging, false)
		expectGetTempIncludeActivation(client, actResWithTempIDStaging.ActivationID, papi.ActivationNetworkStaging)

		// read
		activations := expectListIncludeActivations(client)
		actID, err := getLatestIncludeActivationID(activations, networkStaging)
		require.NoError(t, err)
		expectGetIncludeActivation(client, actID, papi.ActivationNetworkStaging)

		// refresh
		activations = expectListIncludeActivations(client)
		actID, err = getLatestIncludeActivationID(activations, networkStaging)
		require.NoError(t, err)
		expectGetIncludeActivation(client, actID, papi.ActivationNetworkStaging)

		// import
		activations = expectListIncludeActivations(client)
		actID, err = getLatestIncludeActivationID(activations, networkStaging)
		require.NoError(t, err)
		expectGetIncludeActivation(client, actID, papi.ActivationNetworkStaging)

		// destroy
		deactivation := expectDeactivateInclude(client, papi.ActivationNetworkStaging, false)
		expectGetTempIncludeDeactivation(client, deactivation.ActivationID, papi.ActivationNetworkStaging)

		useClient(client, nil, func() {
			resource.UnitTest(t, resource.TestCase{
				Providers: testAccProviders,
				Steps: []resource.TestStep{
					{
						Config: loadFixtureString(fmt.Sprintf("%s/property_include_activation.tf", testDir)),
					},
					{
						ImportState:       true,
						ImportStateVerify: true,
						ImportStateId:     "ctr_test_contract:grp_test_group:inc_12345:STAGING",
						ResourceName:      "akamai_property_include_activation.activation",
					},
				},
			})
		})
		client.AssertExpectations(t)
	})
}

func Test_addComplianceRecordByNetwork(t *testing.T) {
	_, err := addComplianceRecordByNetwork(networkProduction, "activate", []interface{}{}, papi.ActivateOrDeactivateIncludeRequest{})
	require.Error(t, err)
	assert.True(t, regexp.MustCompile("compliance_record field is required for 'PRODUCTION' network to activate include version").MatchString(err.Error()))
}
