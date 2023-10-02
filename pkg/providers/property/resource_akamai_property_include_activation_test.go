package property

import (
	"fmt"
	"math/rand"
	"regexp"
	"strconv"
	"testing"
	"time"

	"github.com/akamai/AkamaiOPEN-edgegrid-golang/v7/pkg/papi"
	"github.com/akamai/terraform-provider-akamai/v5/pkg/common/testutils"
	"github.com/hashicorp/terraform-plugin-testing/helper/resource"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/mock"
)

var (
	includeID  = "inc_12345"
	contractID = "ctr_test_contract"
	groupID    = "grp_test_group"
	accountID  = "test_account"
	note       = "test activation"
	email      = "jbond@example.com"
	version    = 3
	testDir    = "testdata/TestResPropertyIncludeActivation"
)

func TestResourcePropertyIncludeActivation(t *testing.T) {

	// lower down the timeouts for testing purposes
	activationPollInterval = time.Microsecond
	getActivationInterval = time.Microsecond

	type attrs struct {
		includeID, contractID, groupID, network, note string
		version                                       int
		notifyEmails                                  []string
		autoAcknowledgeRuleWarnings                   bool
		timeout                                       string
	}

	type State struct {
		activations []papi.IncludeActivation
	}

	var (
		activateIncludeReq = func(network papi.ActivationNetwork, acknowledgeAllWarnings bool) papi.ActivateIncludeRequest {
			req := papi.ActivateIncludeRequest{
				IncludeID:              includeID,
				Version:                3,
				Network:                network,
				Note:                   note,
				NotifyEmails:           []string{email},
				AcknowledgeAllWarnings: acknowledgeAllWarnings,
			}
			if network == papi.ActivationNetworkProduction {
				req.ComplianceRecord = &papi.ComplianceRecordOther{}
			}
			return req
		}

		deactivateIncludeReq = func(network papi.ActivationNetwork, acknowledgeAllWarnings bool) papi.DeactivateIncludeRequest {
			req := papi.DeactivateIncludeRequest{
				IncludeID:              includeID,
				Version:                3,
				Network:                network,
				Note:                   note,
				NotifyEmails:           []string{email},
				AcknowledgeAllWarnings: acknowledgeAllWarnings,
			}
			if network == papi.ActivationNetworkProduction {
				req.ComplianceRecord = &papi.ComplianceRecordOther{}
			}
			return req
		}

		expectListIncludeActivations = func(client *papi.Mock, activations []papi.IncludeActivation) {
			client.On("ListIncludeActivations", mock.Anything, mock.Anything).
				Return(&papi.ListIncludeActivationsResponse{
					AccountID:  accountID,
					ContractID: contractID,
					GroupID:    groupID,
					Activations: papi.IncludeActivationsRes{
						Items: append([]papi.IncludeActivation(nil), activations...),
					},
				}, nil).Once()
		}

		expectGetIncludeActivation = func(client *papi.Mock, activation papi.IncludeActivation) *mock.Call {
			return client.On("GetIncludeActivation", mock.Anything, papi.GetIncludeActivationRequest{
				IncludeID:    includeID,
				ActivationID: activation.ActivationID,
			}).Return(&papi.GetIncludeActivationResponse{
				AccountID:  accountID,
				ContractID: contractID,
				GroupID:    groupID,
				Activation: activation,
			}, nil)
		}

		expectWaitPending = func(client *papi.Mock, state State, network papi.ActivationNetwork, Npending int) State {
			expectListIncludeActivations(client, state.activations)

			if len(state.activations) == 0 {
				// if there are no activations, wait logic wont do any other calls
				// so there is nothing more to mock -> return
				return state
			}

			for n := range state.activations {
				if state.activations[n].Network == network {
					if Npending > 0 && state.activations[n].Status == papi.ActivationStatusPending {
						expectGetIncludeActivation(client, state.activations[n]).Times(Npending - 1)
					}
					// Mutate state -> change activation state to active
					state.activations[n].Status = papi.ActivationStatusActive

					expectGetIncludeActivation(client, state.activations[n]).Once()

					return state
				}
			}

			// if not found, mock a call that returns an error
			client.On("GetIncludeActivation", mock.Anything, mock.Anything).
				Return(nil, fmt.Errorf("%w: %s", papi.ErrNotFound, papi.ErrGetIncludeActivation)).Once()

			return state
		}

		getRandomActID = func() string {
			return fmt.Sprintf("atv_%d", rand.Int()%10000)
		}

		getActivationBasedOnRequest = func(req papi.ActivateOrDeactivateIncludeRequest, activationType papi.ActivationType) papi.IncludeActivation {
			return papi.IncludeActivation{
				ActivationID:   getRandomActID(),
				Network:        req.Network,
				ActivationType: activationType,
				Status:         papi.ActivationStatusPending,
				SubmitDate:     "",
				UpdateDate:     time.Now().String(),
				NotifyEmails:   req.NotifyEmails,
				Note:           req.Note,
				IncludeID:      req.IncludeID,
				IncludeVersion: req.Version,
			}
		}

		getExpectedActivationBasedOnRequest = func(req papi.ActivateIncludeRequest) papi.IncludeActivation {
			return getActivationBasedOnRequest(papi.ActivateOrDeactivateIncludeRequest(req), papi.ActivationTypeActivate)
		}

		getActivationBasedOnDeactivationRequest = func(req papi.DeactivateIncludeRequest) papi.IncludeActivation {
			return getActivationBasedOnRequest(papi.ActivateOrDeactivateIncludeRequest(req), papi.ActivationTypeDeactivate)
		}

		expectActivateInclude = func(client *papi.Mock, state State, req papi.ActivateIncludeRequest, Nretries int) State {

			newIncludeActivation := getExpectedActivationBasedOnRequest(req)

			client.On("ActivateInclude", mock.Anything, req).
				Return(&papi.ActivationIncludeResponse{
					ActivationID: newIncludeActivation.ActivationID,
				}, nil).Once()

			if Nretries > 0 {
				// here we want to simulate some failing calls that may happen and upsert should just retry
				client.On("GetIncludeActivation", mock.Anything, papi.GetIncludeActivationRequest{
					IncludeID:    includeID,
					ActivationID: newIncludeActivation.ActivationID,
				}).Return(nil, fmt.Errorf("%w: %s", papi.ErrNotFound, papi.ErrGetIncludeActivation)).Times(Nretries)
			}

			// mutate state - add new activation
			state.activations = append([]papi.IncludeActivation{newIncludeActivation}, state.activations...)

			expectGetIncludeActivation(client, state.activations[0]).Once()

			return state
		}

		expectActivateIncludeWithFail = func(client *papi.Mock, state State, req papi.ActivateIncludeRequest) {
			client.On("ActivateInclude", mock.Anything, req).
				Return(nil, fmt.Errorf("oops")).Once()
		}

		expectAssertState = func(client *papi.Mock, state State) {
			expectListIncludeActivations(client, state.activations)
		}

		expectCreate = func(client *papi.Mock, state State, req papi.ActivateIncludeRequest) State {
			state = expectWaitPending(client, state, req.Network, 2)
			expectAssertState(client, state)
			state = expectActivateInclude(client, state, req, 2)
			state = expectWaitPending(client, state, req.Network, 2)
			return state
		}

		expectCreateWithFail = func(client *papi.Mock, state State, req papi.ActivateIncludeRequest) State {
			state = expectWaitPending(client, state, req.Network, 2)
			expectAssertState(client, state)
			expectActivateIncludeWithFail(client, state, req)
			return state
		}

		expectCreateOnlyReadState = func(client *papi.Mock, state State, req papi.ActivateIncludeRequest) State {
			state = expectWaitPending(client, state, req.Network, 2)
			expectAssertState(client, state)
			return state
		}

		expectRead = func(client *papi.Mock, state State, network papi.ActivationNetwork) {
			expectListIncludeActivations(client, state.activations)

			// read filters the list for latest active include in given network
			// find it and mock the call for it
			for _, a := range state.activations {
				if a.Network == network && a.Status == papi.ActivationStatusActive {
					expectGetIncludeActivation(client, a).Once()
					return
				}
			}
			// if not found, mock a call that returns an error
			client.On("GetIncludeActivation", mock.Anything, mock.Anything).
				Return(nil, fmt.Errorf("%w: %s", papi.ErrNotFound, papi.ErrGetIncludeActivation)).Once()
		}

		expectDectivateInclude = func(client *papi.Mock, state State, req papi.DeactivateIncludeRequest, Nretries int) State {
			newIncludeDeactivation := getActivationBasedOnDeactivationRequest(req)

			client.On("DeactivateInclude", mock.Anything, req).
				Return(&papi.DeactivationIncludeResponse{
					ActivationID: newIncludeDeactivation.ActivationID,
				}, nil).Once()

			if Nretries > 0 {
				// here we want to simulate some failing calls that may happen and upsert should just retry
				client.On("GetIncludeActivation", mock.Anything, papi.GetIncludeActivationRequest{
					IncludeID:    includeID,
					ActivationID: newIncludeDeactivation.ActivationID,
				}).Return(nil, fmt.Errorf("%w: %s", papi.ErrNotFound, papi.ErrGetIncludeActivation)).Times(Nretries)
			}

			// mutate state - add deactivation
			state.activations = append([]papi.IncludeActivation{newIncludeDeactivation}, state.activations...)

			expectGetIncludeActivation(client, newIncludeDeactivation).Once()

			return state
		}

		expectDelete = func(client *papi.Mock, state State, req papi.DeactivateIncludeRequest) State {
			state = expectWaitPending(client, state, req.Network, 2)
			expectAssertState(client, state)
			state = expectDectivateInclude(client, state, req, 2)
			state = expectWaitPending(client, state, req.Network, 2)
			return state
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
				resource.TestCheckResourceAttr("akamai_property_include_activation.activation", "auto_acknowledge_rule_warnings", strconv.FormatBool(attrs.autoAcknowledgeRuleWarnings)),
			}

			if attrs.timeout != "" {
				checks = append(checks, []resource.TestCheckFunc{
					resource.TestCheckResourceAttr("akamai_property_include_activation.activation", "timeouts.#", "1"),
					resource.TestCheckResourceAttr("akamai_property_include_activation.activation", "timeouts.0.default", attrs.timeout),
				}...)
			} else {
				checks = append(checks, resource.TestCheckResourceAttr("akamai_property_include_activation.activation", "timeouts.#", "0"))
			}

			return resource.ComposeAggregateTestCheckFunc(checks...)
		}
	)

	t.Run("create a new include activation lifecycle", func(t *testing.T) {
		client := new(papi.Mock)
		state := State{}

		// create
		actReq := activateIncludeReq("STAGING", false)
		state = expectCreate(client, state, actReq)

		// read
		expectRead(client, state, papi.ActivationNetworkStaging)

		// read
		expectRead(client, state, papi.ActivationNetworkStaging)

		// delete
		deactReq := deactivateIncludeReq("STAGING", false)
		_ = expectDelete(client, state, deactReq)

		useClient(client, nil, func() {
			resource.UnitTest(t, resource.TestCase{
				ProtoV6ProviderFactories: testutils.NewProtoV6ProviderFactory(NewSubprovider()),
				Steps: []resource.TestStep{
					{
						Config: testutils.LoadFixtureString(t, fmt.Sprintf("%s/property_include_activation.tf", testDir)),
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

	t.Run("include activation with timeout lifecycle", func(t *testing.T) {
		client := new(papi.Mock)
		state := State{}

		// create
		actReq := activateIncludeReq("STAGING", false)
		state = expectCreate(client, state, actReq)

		// read
		expectRead(client, state, papi.ActivationNetworkStaging)

		// update
		// no actual update on only timeout change
		expectRead(client, state, papi.ActivationNetworkStaging)
		expectRead(client, state, papi.ActivationNetworkStaging)

		// read
		expectRead(client, state, papi.ActivationNetworkStaging)

		// delete
		deactReq := deactivateIncludeReq("STAGING", false)
		_ = expectDelete(client, state, deactReq)

		useClient(client, nil, func() {
			resource.UnitTest(t, resource.TestCase{
				ProtoV5ProviderFactories: testAccProviders,
				Steps: []resource.TestStep{
					{
						Config: testutils.LoadFixtureString(t, fmt.Sprintf("%s/property_include_activation_with_timeout.tf", testDir)),
						Check: checkAttributes(attrs{
							includeID:    includeID,
							contractID:   contractID,
							groupID:      groupID,
							version:      version,
							network:      "STAGING",
							note:         note,
							notifyEmails: []string{email},
							timeout:      "2h1m",
						}),
					},
					{
						Config: testutils.LoadFixtureString(t, fmt.Sprintf("%s/property_include_activation_with_timeout_update.tf", testDir)),
						Check: checkAttributes(attrs{
							includeID:    includeID,
							contractID:   contractID,
							groupID:      groupID,
							version:      version,
							network:      "STAGING",
							note:         note,
							notifyEmails: []string{email},
							timeout:      "2h2m",
						}),
					},
				},
			})
		})
		client.AssertExpectations(t)
	})

	t.Run("update include activation lifecycle", func(t *testing.T) {
		client := new(papi.Mock)
		state := State{}

		// 1. first step

		// create
		actReq := activateIncludeReq("STAGING", false)
		state = expectCreate(client, state, actReq)

		// read
		expectRead(client, state, papi.ActivationNetworkStaging)

		// read
		expectRead(client, state, papi.ActivationNetworkStaging)

		// 2. second step - network ForceNew

		// read
		expectRead(client, state, papi.ActivationNetworkStaging)

		// delete
		deactReq := deactivateIncludeReq("STAGING", false)
		state = expectDelete(client, state, deactReq)

		// create
		actReq = activateIncludeReq("PRODUCTION", true)
		state = expectCreate(client, state, actReq)

		// read
		expectRead(client, state, papi.ActivationNetworkProduction)

		// read
		expectRead(client, state, papi.ActivationNetworkProduction)

		// delete
		deactReq = deactivateIncludeReq("PRODUCTION", true)
		_ = expectDelete(client, state, deactReq)

		useClient(client, nil, func() {
			resource.UnitTest(t, resource.TestCase{
				ProtoV6ProviderFactories: testutils.NewProtoV6ProviderFactory(NewSubprovider()),
				Steps: []resource.TestStep{
					{
						Config: testutils.LoadFixtureString(t, fmt.Sprintf("%s/property_include_activation.tf", testDir)),
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
						Config: testutils.LoadFixtureString(t, fmt.Sprintf("%s/property_include_activation_update.tf", testDir)),
						Check: checkAttributes(attrs{
							includeID:                   includeID,
							contractID:                  contractID,
							groupID:                     groupID,
							version:                     version,
							network:                     "PRODUCTION",
							note:                        note,
							notifyEmails:                []string{email},
							autoAcknowledgeRuleWarnings: true,
						}),
					},
				},
			})
		})

		client.AssertExpectations(t)
	})

	t.Run("wait for ongoing expected activation to finish", func(t *testing.T) {
		client := new(papi.Mock)
		state := State{}

		actReq := activateIncludeReq("STAGING", false)

		pendingIncludeActivation := getExpectedActivationBasedOnRequest(actReq)
		state.activations = append(state.activations, pendingIncludeActivation)

		// create
		state = expectCreateOnlyReadState(client, state, actReq)

		// read
		expectRead(client, state, papi.ActivationNetworkStaging)

		// read
		expectRead(client, state, papi.ActivationNetworkStaging)

		// delete
		deactReq := deactivateIncludeReq("STAGING", false)
		_ = expectDelete(client, state, deactReq)

		useClient(client, nil, func() {
			resource.UnitTest(t, resource.TestCase{
				ProtoV6ProviderFactories: testutils.NewProtoV6ProviderFactory(NewSubprovider()),
				Steps: []resource.TestStep{
					{
						Config: testutils.LoadFixtureString(t, fmt.Sprintf("%s/property_include_activation.tf", testDir)),
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

	t.Run("wait for ongoing unexpected activation to finish", func(t *testing.T) {
		client := new(papi.Mock)
		state := State{}

		actReq := activateIncludeReq("STAGING", false)
		unexpectedActivationReq := activateIncludeReq("PRODUCTION", false)
		unexpectedActivationReq.Version = 2

		assert.NotEqual(t, unexpectedActivationReq.Version, actReq.Version)

		pendingIncludeActivation := getExpectedActivationBasedOnRequest(unexpectedActivationReq)
		state.activations = append(state.activations, pendingIncludeActivation)

		// create
		state = expectCreate(client, state, actReq)

		// read
		expectRead(client, state, papi.ActivationNetworkStaging)

		// read
		expectRead(client, state, papi.ActivationNetworkStaging)

		// delete
		deactReq := deactivateIncludeReq("STAGING", false)
		_ = expectDelete(client, state, deactReq)

		useClient(client, nil, func() {
			resource.UnitTest(t, resource.TestCase{
				ProtoV6ProviderFactories: testutils.NewProtoV6ProviderFactory(NewSubprovider()),
				Steps: []resource.TestStep{
					{
						Config: testutils.LoadFixtureString(t, fmt.Sprintf("%s/property_include_activation.tf", testDir)),
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

	t.Run("missing compliance record error when network is PRODUCTION", func(t *testing.T) {
		client := new(papi.Mock)
		state := State{}

		// create
		req := papi.ActivateIncludeRequest{
			IncludeID:              includeID,
			Version:                3,
			Network:                "PRODUCTION",
			Note:                   note,
			NotifyEmails:           []string{email},
			AcknowledgeAllWarnings: false,
		}
		state = expectWaitPending(client, state, req.Network, 2)
		expectAssertState(client, state)
		newIncludeActivation := getExpectedActivationBasedOnRequest(req)

		client.On("ActivateInclude", mock.Anything, req).
			Return(&papi.ActivationIncludeResponse{
				ActivationID: newIncludeActivation.ActivationID,
			}, nil).Once()

		// GetIncludeActivation returns error about missing_compliance_record. TFP checks for that error and returns it
		client.On("GetIncludeActivation", mock.Anything, papi.GetIncludeActivationRequest{
			IncludeID:    includeID,
			ActivationID: newIncludeActivation.ActivationID,
		}).Return(nil, papi.ErrMissingComplianceRecord).Once()

		useClient(client, nil, func() {
			resource.UnitTest(t, resource.TestCase{
				ProtoV6ProviderFactories: testutils.NewProtoV6ProviderFactory(NewSubprovider()),
				Steps: []resource.TestStep{
					{
						Config:      testutils.LoadFixtureString(t, fmt.Sprintf("%s/no_compliance_record_on_production.tf", testDir)),
						ExpectError: regexp.MustCompile(`Error: for 'PRODUCTION' network, 'compliance_record' must be specified`),
					},
				},
			})
		})
		client.AssertExpectations(t)
	})

	t.Run("incorrect timeout format", func(t *testing.T) {

		useClient(client, nil, func() {
			resource.UnitTest(t, resource.TestCase{
				ProtoV5ProviderFactories: testAccProviders,
				Steps: []resource.TestStep{
					{
						Config:      testutils.LoadFixtureString(t, fmt.Sprintf("%s/property_include_activation_incorrect_timeout.tf", testDir)),
						ExpectError: regexp.MustCompile(`provided incorrect duration`),
					},
				},
			})
		})
	})

	t.Run("first create fails but second create works", func(t *testing.T) {
		client := new(papi.Mock)
		state := State{}

		actReq := activateIncludeReq("STAGING", false)

		// --- 1st step ---

		// create -> fail
		state = expectCreateWithFail(client, state, actReq)

		// --- 2nd step ---

		// create
		state = expectCreate(client, state, actReq)

		// read
		expectRead(client, state, papi.ActivationNetworkStaging)

		// read
		expectRead(client, state, papi.ActivationNetworkStaging)

		// delete
		deactReq := deactivateIncludeReq("STAGING", false)
		_ = expectDelete(client, state, deactReq)

		useClient(client, nil, func() {
			resource.UnitTest(t, resource.TestCase{
				ProtoV6ProviderFactories: testutils.NewProtoV6ProviderFactory(NewSubprovider()),
				Steps: []resource.TestStep{
					{
						Config: testutils.LoadFixtureString(t, fmt.Sprintf("%s/property_include_activation.tf", testDir)),
						Check: checkAttributes(attrs{
							includeID:    includeID,
							contractID:   contractID,
							groupID:      groupID,
							version:      version,
							network:      "STAGING",
							note:         note,
							notifyEmails: []string{email},
						}),
						ExpectError: regexp.MustCompile("oops"),
					},
					{
						Config: testutils.LoadFixtureString(t, fmt.Sprintf("%s/property_include_activation.tf", testDir)),
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

	t.Run("import", func(t *testing.T) {
		client := new(papi.Mock)
		state := State{}

		// create
		actReq := activateIncludeReq("STAGING", false)
		state = expectCreate(client, state, actReq)

		// read
		expectRead(client, state, papi.ActivationNetworkStaging)

		// read
		expectRead(client, state, papi.ActivationNetworkStaging)

		// read
		expectRead(client, state, papi.ActivationNetworkStaging)

		// delete
		deactReq := deactivateIncludeReq("STAGING", false)
		_ = expectDelete(client, state, deactReq)

		useClient(client, nil, func() {
			resource.UnitTest(t, resource.TestCase{
				ProtoV6ProviderFactories: testutils.NewProtoV6ProviderFactory(NewSubprovider()),
				Steps: []resource.TestStep{
					{
						Config: testutils.LoadFixtureString(t, fmt.Sprintf("%s/property_include_activation.tf", testDir)),
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

	t.Run("note filed change suppressed", func(t *testing.T) {
		client := new(papi.Mock)
		state := State{}

		// 1. first step

		// create
		actReq := activateIncludeReq("STAGING", false)
		state = expectCreate(client, state, actReq)

		// read
		expectRead(client, state, papi.ActivationNetworkStaging)

		// read
		expectRead(client, state, papi.ActivationNetworkStaging)

		// 2. second step - update only note field - change suppressed
		// delete
		deactReq := deactivateIncludeReq("STAGING", false)
		_ = expectDelete(client, state, deactReq)

		useClient(client, nil, func() {
			resource.UnitTest(t, resource.TestCase{
				ProtoV6ProviderFactories: testutils.NewProtoV6ProviderFactory(NewSubprovider()),
				Steps: []resource.TestStep{
					{
						Config: testutils.LoadFixtureString(t, fmt.Sprintf("%s/property_include_activation.tf", testDir)),
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

	t.Run("note and version filed change not suppressed", func(t *testing.T) {
		client := new(papi.Mock)
		state := State{}

		// 1. first step

		// create
		actReq := activateIncludeReq("STAGING", false)
		state = expectCreate(client, state, actReq)

		// read
		expectRead(client, state, papi.ActivationNetworkStaging)

		// read
		expectRead(client, state, papi.ActivationNetworkStaging)

		expectRead(client, state, papi.ActivationNetworkStaging)

		// 2. second step - update note and version with creation of new activation - note change not suppressed
		// create
		req := papi.ActivateIncludeRequest{
			IncludeID:    includeID,
			Version:      4,
			Network:      "STAGING",
			Note:         "not suppressed note field change",
			NotifyEmails: []string{email},
		}
		state = expectCreate(client, state, req)

		// read
		expectRead(client, state, papi.ActivationNetworkStaging)

		// read
		expectRead(client, state, papi.ActivationNetworkStaging)

		// delete
		deactReq := papi.DeactivateIncludeRequest{
			IncludeID:    includeID,
			Version:      4,
			Network:      "STAGING",
			Note:         "not suppressed note field change",
			NotifyEmails: []string{email},
		}
		_ = expectDelete(client, state, deactReq)

		useClient(client, nil, func() {
			resource.UnitTest(t, resource.TestCase{
				ProtoV6ProviderFactories: testutils.NewProtoV6ProviderFactory(NewSubprovider()),
				Steps: []resource.TestStep{
					{
						Config: testutils.LoadFixtureString(t, fmt.Sprintf("%s/property_include_activation.tf", testDir)),
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
						Config: testutils.LoadFixtureString(t, fmt.Sprintf("%s/property_include_update_note_not_suppressed.tf", testDir)),
						Check: checkAttributes(attrs{
							includeID:    includeID,
							contractID:   contractID,
							groupID:      groupID,
							version:      4,
							network:      "STAGING",
							note:         "not suppressed note field change",
							notifyEmails: []string{email},
						}),
					},
				},
			})
		})

		client.AssertExpectations(t)
	})
}

func TestReadTimeoutFromEnvOrDefault(t *testing.T) {
	tests := map[string]struct {
		envName      string
		envValue     string
		defaultValue time.Duration
		expect       time.Duration
	}{
		"no env value set": {
			envName:      "TEST_NAME",
			envValue:     "",
			defaultValue: time.Hour,
			expect:       time.Hour,
		},
		"correct env value 120 set": {
			envName:      "TEST_NAME",
			envValue:     "120",
			defaultValue: time.Hour,
			expect:       time.Hour * 2,
		},
		"correct env value 12 set": {
			envName:      "TEST_NAME",
			envValue:     "12",
			defaultValue: time.Hour,
			expect:       time.Minute * 12,
		},
		"incorrect env value set": {
			envName:      "TEST_NAME_2",
			envValue:     "not a number",
			defaultValue: time.Hour,
			expect:       time.Hour,
		},
	}
	for name, test := range tests {
		t.Run(name, func(t *testing.T) {
			t.Setenv(test.envName, test.envValue)
			result := readTimeoutFromEnvOrDefault(test.envName, test.defaultValue)
			assert.Equal(t, test.expect, *result)
		})
	}
}
