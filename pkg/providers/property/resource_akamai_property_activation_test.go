package property

import (
	"fmt"
	"regexp"
	"testing"

	"github.com/stretchr/testify/mock"

	"github.com/akamai/AkamaiOPEN-edgegrid-golang/v2/pkg/papi"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/resource"
)

type papiCall struct {
	methodName   string
	papiResponse interface{}
	papiRequest  interface{}
	error        error
	stubOnce     bool
}

func mockPAPIClient(callsToMock []papiCall) *mockpapi {
	client := &mockpapi{}
	for _, call := range callsToMock {
		var request interface{}
		request = mock.Anything
		if call.papiRequest != nil {
			request = call.papiRequest
		}
		stub := client.On(call.methodName, AnyCTX, request).Return(call.papiResponse, call.error)
		if call.stubOnce {
			stub.Once()
		}
	}

	return client
}

func TestResourcePropertyActivationCreate(t *testing.T) {
	t.Run("check schema property activation - OK", func(t *testing.T) {

		client := mockPAPIClient([]papiCall{
			{
				methodName: "GetRuleTree",
				papiResponse: &papi.GetRuleTreeResponse{
					Response: papi.Response{
						Errors:   make([]*papi.Error, 0),
						Warnings: []*papi.Error{{Title: "some warning"}},
					},
				},
				error:    nil,
				stubOnce: false,
			},
			{
				methodName: "GetActivations",
				papiResponse: &papi.GetActivationsResponse{
					Activations: papi.ActivationsItems{Items: []*papi.Activation{{
						AccountID:       "act_1-6JHGX",
						ActivationID:    "atv_activation1",
						ActivationType:  "ACTIVATE",
						GroupID:         "grp_91533",
						PropertyName:    "test",
						PropertyID:      "prp_test",
						PropertyVersion: 1,
						Network:         "STAGING",
						Status:          "ACTIVE",
						SubmitDate:      "2020-10-28T15:04:05Z",
					}}}},
				error:    nil,
				stubOnce: true,
			},
			{
				methodName: "GetActivations",
				papiResponse: &papi.GetActivationsResponse{
					Activations: papi.ActivationsItems{Items: []*papi.Activation{{
						AccountID:       "act_1-6JHGX",
						ActivationID:    "atv_deactivation1",
						ActivationType:  "DEACTIVATE",
						GroupID:         "grp_91533",
						PropertyName:    "test",
						PropertyID:      "prp_test",
						PropertyVersion: 1,
						Network:         "STAGING",
						Status:          "ACTIVE",
						SubmitDate:      "2020-10-28T15:05:05Z",
					}}}},
				error:    nil,
				stubOnce: true,
			},
			{
				methodName: "GetActivations",
				papiResponse: &papi.GetActivationsResponse{
					Activations: papi.ActivationsItems{Items: []*papi.Activation{{
						AccountID:       "act_1-6JHGX",
						ActivationID:    "atv_delete1",
						ActivationType:  "DEACTIVATE",
						GroupID:         "grp_91533",
						PropertyName:    "test",
						PropertyID:      "prp_test",
						PropertyVersion: 1,
						Network:         "STAGING",
						Status:          "ACTIVE",
						SubmitDate:      "2020-10-28T15:06:05Z",
					}}}},
				error:    nil,
				stubOnce: true,
			},
		})
		useClient(client, func() {
			resource.UnitTest(t, resource.TestCase{
				IsUnitTest: true,
				Providers:  testAccProviders,
				Steps: []resource.TestStep{
					{
						Config: loadFixtureString("testdata/TestPropertyActivation/ok/resource_property_activation.tf"),
						Check: resource.ComposeAggregateTestCheckFunc(
							resource.TestCheckResourceAttr("akamai_property_activation.test", "id", "prp_test:STAGING"),
							resource.TestCheckResourceAttr("akamai_property_activation.test", "property_id", "prp_test"),
							resource.TestCheckResourceAttr("akamai_property_activation.test", "network", "STAGING"),
							resource.TestCheckResourceAttr("akamai_property_activation.test", "version", "1"),
							resource.TestCheckResourceAttr("akamai_property_activation.test", "auto_acknowledge_rule_warnings", "true"),
							resource.TestCheckNoResourceAttr("akamai_property_activation.test", "warnings"),
							resource.TestCheckNoResourceAttr("akamai_property_activation.test", "rule_warnings"),
							resource.TestCheckNoResourceAttr("akamai_property_activation.test", "errors"),
							resource.TestCheckResourceAttr("akamai_property_activation.test", "activation_id", "atv_activation1"),
							resource.TestCheckResourceAttr("akamai_property_activation.test", "status", "ACTIVE"),
							resource.TestCheckResourceAttr("akamai_property_activation.test", "note", "property activation note for creating"),
						),
					},
				},
			})
		})

		client.AssertExpectations(t)
	})

	t.Run("check schema property activation - papi error", func(t *testing.T) {

		client := mockPAPIClient([]papiCall{
			{
				methodName:   "GetRuleTree",
				papiResponse: nil,
				error:        fmt.Errorf("failed to create request"),
				stubOnce:     false,
			},
		})
		useClient(client, func() {
			resource.UnitTest(t, resource.TestCase{
				IsUnitTest: true,
				Providers:  testAccProviders,
				Steps: []resource.TestStep{
					{
						Config:      loadFixtureString("testdata/TestPropertyActivation/ok/resource_property_activation.tf"),
						ExpectError: regexp.MustCompile("failed to create request"),
					},
				},
			})
		})

		client.AssertExpectations(t)
	})

	t.Run("check schema property activation - no property id nor property", func(t *testing.T) {

		client := mockPAPIClient([]papiCall{})
		useClient(client, func() {
			resource.UnitTest(t, resource.TestCase{
				IsUnitTest: true,
				Providers:  testAccProviders,
				Steps: []resource.TestStep{
					{
						Config:      loadFixtureString("testdata/TestPropertyActivation/no_propertyId/resource_property_activation.tf"),
						ExpectError: regexp.MustCompile("one of `property,property_id` must be specified"),
					},
				},
			})
		})

		client.AssertExpectations(t)
	})

	t.Run("check schema property activation - no contact", func(t *testing.T) {

		client := mockPAPIClient([]papiCall{})
		useClient(client, func() {
			resource.UnitTest(t, resource.TestCase{
				IsUnitTest: true,
				Providers:  testAccProviders,
				Steps: []resource.TestStep{
					{
						Config:      loadFixtureString("testdata/TestPropertyActivation/no_contact/resource_property_activation.tf"),
						ExpectError: regexp.MustCompile("Missing required argument"),
					},
				},
			})
		})

		client.AssertExpectations(t)
	})

	t.Run("schema with `property` instead of `property_id`", func(t *testing.T) {

		client := mockPAPIClient([]papiCall{
			{
				methodName: "GetRuleTree",
				papiResponse: &papi.GetRuleTreeResponse{
					Response: papi.Response{Errors: make([]*papi.Error, 0)},
				},
				error:    nil,
				stubOnce: false,
			},
			{
				methodName: "GetActivations",
				papiResponse: &papi.GetActivationsResponse{
					Activations: papi.ActivationsItems{Items: []*papi.Activation{{
						AccountID:       "act_1-6JHGX",
						ActivationID:    "atv_activation1",
						ActivationType:  "ACTIVATE",
						GroupID:         "grp_91533",
						PropertyName:    "test",
						PropertyID:      "prp_test",
						PropertyVersion: 1,
						Network:         "STAGING",
						Status:          "ACTIVE",
						SubmitDate:      "2020-10-28T15:04:05Z",
					}}}},
				error:    nil,
				stubOnce: true,
			},
			{
				methodName: "GetActivations",
				papiResponse: &papi.GetActivationsResponse{
					Activations: papi.ActivationsItems{Items: []*papi.Activation{{
						AccountID:       "act_1-6JHGX",
						ActivationID:    "atv_deactivation1",
						ActivationType:  "DEACTIVATE",
						GroupID:         "grp_91533",
						PropertyName:    "test",
						PropertyID:      "prp_test",
						PropertyVersion: 1,
						Network:         "STAGING",
						Status:          "ACTIVE",
						SubmitDate:      "2020-10-28T15:05:05Z",
					}}}},
				error:    nil,
				stubOnce: true,
			},
			{
				methodName: "GetActivations",
				papiResponse: &papi.GetActivationsResponse{
					Activations: papi.ActivationsItems{Items: []*papi.Activation{{
						AccountID:       "act_1-6JHGX",
						ActivationID:    "atv_delete1",
						ActivationType:  "DEACTIVATE",
						GroupID:         "grp_91533",
						PropertyName:    "test",
						PropertyID:      "prp_test",
						PropertyVersion: 1,
						Network:         "STAGING",
						Status:          "ACTIVE",
						SubmitDate:      "2020-10-28T15:06:05Z",
					}}}},
				error:    nil,
				stubOnce: true,
			},
		})
		useClient(client, func() {
			resource.UnitTest(t, resource.TestCase{
				IsUnitTest: true,
				Providers:  testAccProviders,
				Steps: []resource.TestStep{
					{
						Config: loadFixtureString("testdata/TestPropertyActivation/ok/resource_property_activation_deprecated_arg.tf"),
						Check: resource.ComposeAggregateTestCheckFunc(
							resource.TestCheckResourceAttr("akamai_property_activation.test", "id", "prp_test:STAGING"),
							resource.TestCheckResourceAttr("akamai_property_activation.test", "property_id", "prp_test"),
							resource.TestCheckResourceAttr("akamai_property_activation.test", "property", "prp_test"),
							resource.TestCheckResourceAttr("akamai_property_activation.test", "network", "STAGING"),
							resource.TestCheckResourceAttr("akamai_property_activation.test", "version", "1"),
							resource.TestCheckNoResourceAttr("akamai_property_activation.test", "warnings"),
							resource.TestCheckNoResourceAttr("akamai_property_activation.test", "errors"),
							resource.TestCheckResourceAttr("akamai_property_activation.test", "activation_id", "atv_activation1"),
							resource.TestCheckResourceAttr("akamai_property_activation.test", "status", "ACTIVE"),
						),
					},
				},
			})
		})

		client.AssertExpectations(t)
	})

	t.Run("check schema property activation update", func(t *testing.T) {

		client := mockPAPIClient([]papiCall{
			{
				methodName: "GetRuleTree",
				papiResponse: &papi.GetRuleTreeResponse{
					Response: papi.Response{
						Errors:   make([]*papi.Error, 0),
						Warnings: []*papi.Error{{Title: "some warning"}},
					},
				},
				error:    nil,
				stubOnce: false,
			},
			{
				methodName: "GetActivations",
				papiRequest: papi.GetActivationsRequest{
					PropertyID: "prp_test",
				},
				papiResponse: &papi.GetActivationsResponse{
					Activations: papi.ActivationsItems{Items: []*papi.Activation{{
						AccountID:       "act_1-6JHGX",
						ActivationID:    "atv_activation1",
						ActivationType:  "ACTIVATE",
						GroupID:         "grp_91533",
						PropertyName:    "test",
						PropertyID:      "prp_test",
						PropertyVersion: 1,
						Network:         "STAGING",
						Status:          "ACTIVE",
						SubmitDate:      "2020-10-28T15:04:05Z",
					}}}},
				error:    nil,
				stubOnce: true,
			},
			{
				methodName: "GetActivations",
				papiResponse: &papi.GetActivationsResponse{
					Activations: papi.ActivationsItems{Items: []*papi.Activation{{
						AccountID:       "act_1-6JHGX",
						ActivationID:    "atv_deactivation1",
						ActivationType:  "DEACTIVATE",
						GroupID:         "grp_91533",
						PropertyName:    "test",
						PropertyID:      "prp_test",
						PropertyVersion: 2,
						Network:         "STAGING",
						Status:          "ACTIVE",
						SubmitDate:      "2020-10-28T15:05:05Z",
					}}}},
				error:    nil,
				stubOnce: false,
			},
			{
				methodName: "GetActivations",
				papiResponse: &papi.GetActivationsResponse{
					Activations: papi.ActivationsItems{Items: []*papi.Activation{{
						AccountID:       "act_1-6JHGX",
						ActivationID:    "atv_delete1",
						ActivationType:  "DEACTIVATE",
						GroupID:         "grp_91533",
						PropertyName:    "test",
						PropertyID:      "prp_test",
						PropertyVersion: 1,
						Network:         "STAGING",
						Status:          "ACTIVE",
						SubmitDate:      "2020-10-28T15:06:05Z",
					}}}},
				error:    nil,
				stubOnce: false,
			},
			{
				methodName: "CreateActivation",
				papiRequest: papi.CreateActivationRequest{
					PropertyID: "prp_test",
					Activation: papi.Activation{
						ActivationType:         papi.ActivationTypeActivate,
						AcknowledgeAllWarnings: true,
						PropertyVersion:        2,
						Network:                "STAGING",
						NotifyEmails:           []string{"user@example.com"},
						Note:                   "property activation note for updating",
					},
				},
				papiResponse: &papi.CreateActivationResponse{
					ActivationID: "atv_update",
				},
				stubOnce: true,
			},
			{
				methodName: "GetActivation",
				papiRequest: papi.GetActivationRequest{
					PropertyID:   "prp_test",
					ActivationID: "atv_update",
				},
				papiResponse: &papi.GetActivationResponse{
					GetActivationsResponse: papi.GetActivationsResponse{},
					Activation: &papi.Activation{
						ActivationID:    "atv_update",
						PropertyID:      "prp_test",
						PropertyVersion: 2,
						Network:         "STAGING",
						Status:          papi.ActivationStatusActive,
					},
				},
			},
		})
		useClient(client, func() {
			resource.UnitTest(t, resource.TestCase{
				IsUnitTest: true,
				Providers:  testAccProviders,
				Steps: []resource.TestStep{
					{
						Config: loadFixtureString("testdata/TestPropertyActivation/ok/resource_property_activation.tf"),
						Check: resource.ComposeAggregateTestCheckFunc(
							resource.TestCheckResourceAttr("akamai_property_activation.test", "id", "prp_test:STAGING"),
							resource.TestCheckResourceAttr("akamai_property_activation.test", "property_id", "prp_test"),
							resource.TestCheckResourceAttr("akamai_property_activation.test", "network", "STAGING"),
							resource.TestCheckResourceAttr("akamai_property_activation.test", "version", "1"),
							resource.TestCheckNoResourceAttr("akamai_property_activation.test", "warnings"),
							resource.TestCheckNoResourceAttr("akamai_property_activation.test", "errors"),
							resource.TestCheckNoResourceAttr("akamai_property_activation.test", "rule_errors"),
							resource.TestCheckResourceAttr("akamai_property_activation.test", "auto_acknowledge_rule_warnings", "true"),
							resource.TestCheckNoResourceAttr("akamai_property_activation.test", "rule_warnings"),
							resource.TestCheckResourceAttr("akamai_property_activation.test", "activation_id", "atv_activation1"),
							resource.TestCheckResourceAttr("akamai_property_activation.test", "status", "ACTIVE"),
							resource.TestCheckResourceAttr("akamai_property_activation.test", "note", "property activation note for creating"),
						),
					},
					{
						Config: loadFixtureString("testdata/TestPropertyActivation/ok/resource_property_activation_update.tf"),
						Check: resource.ComposeAggregateTestCheckFunc(
							resource.TestCheckResourceAttr("akamai_property_activation.test", "id", "prp_test:STAGING"),
							resource.TestCheckResourceAttr("akamai_property_activation.test", "property_id", "prp_test"),
							resource.TestCheckResourceAttr("akamai_property_activation.test", "network", "STAGING"),
							resource.TestCheckResourceAttr("akamai_property_activation.test", "version", "2"),
							resource.TestCheckResourceAttr("akamai_property_activation.test", "activation_id", "atv_update"),
							resource.TestCheckResourceAttr("akamai_property_activation.test", "status", "ACTIVE"),
							resource.TestCheckResourceAttr("akamai_property_activation.test", "note", "property activation note for updating"),
						),
					},
				},
			})
		})

		client.AssertExpectations(t)
	})

	t.Run("check schema property activation with rule errors", func(t *testing.T) {

		client := mockPAPIClient([]papiCall{
			{
				methodName: "GetRuleTree",
				papiResponse: &papi.GetRuleTreeResponse{
					Response: papi.Response{
						Errors: []*papi.Error{
							{
								Title: "some error",
							},
						},
					},
				},
				error:    nil,
				stubOnce: false,
			},
		})
		useClient(client, func() {
			resource.UnitTest(t, resource.TestCase{
				IsUnitTest: true,
				Providers:  testAccProviders,
				Steps: []resource.TestStep{
					{
						Config:      loadFixtureString("testdata/TestPropertyActivation/ok/resource_property_activation.tf"),
						ExpectError: regexp.MustCompile("activation cannot continue due to rule errors"),
					},
				},
			})
		})

		client.AssertExpectations(t)
	})
}

func TestNoteFieldCannotBeAddedAfterActivationIsCompleted(t *testing.T) {
	t.Run("Note field cannot be added after activation is completed", func(t *testing.T) {

		client := mockPAPIClient([]papiCall{
			{
				methodName: "GetRuleTree",
				papiResponse: &papi.GetRuleTreeResponse{
					Response: papi.Response{
						Errors:   make([]*papi.Error, 0),
						Warnings: []*papi.Error{{Title: "some warning"}},
					},
				},
				error:    nil,
				stubOnce: false,
			},
			{
				methodName:   "GetActivations",
				papiResponse: getActivations(),
				error:        nil,
				stubOnce:     false,
			},
			{
				methodName:  "CreateActivation",
				papiRequest: createActivation(),
				papiResponse: &papi.CreateActivationResponse{
					ActivationID: "atv_activation1",
				},
				stubOnce: true,
			},
			{
				methodName: "GetActivation",
				papiRequest: papi.GetActivationRequest{
					PropertyID:   "prp_test",
					ActivationID: "atv_activation1",
				},
				papiResponse: getActivation(),
			},
		})
		useClient(client, func() {
			resource.UnitTest(t, resource.TestCase{
				IsUnitTest: true,
				Providers:  testAccProviders,
				Steps: []resource.TestStep{
					{
						Config: loadFixtureString("testdata/TestPropertyActivation/update_note_field/note_field_doesnt_exist/resource_property_activation.tf"),
					},
					{
						Config:      loadFixtureString("testdata/TestPropertyActivation/update_note_field/note_field_doesnt_exist/resource_property_activation_update.tf"),
						ExpectError: regexp.MustCompile("cannot update activation attribute note after creation"),
					},
				},
			})
		})
		client.AssertExpectations(t)
	})
}

func TestNoteFieldCannotBeUpdatedAfterActivationIsCompleted(t *testing.T) {
	t.Run("Note field cannot be updated after activation is completed", func(t *testing.T) {

		client := mockPAPIClient([]papiCall{
			{
				methodName: "GetRuleTree",
				papiResponse: &papi.GetRuleTreeResponse{
					Response: papi.Response{
						Errors:   make([]*papi.Error, 0),
						Warnings: []*papi.Error{{Title: "some warning"}},
					},
				},
				error:    nil,
				stubOnce: false,
			},
			{
				methodName:   "GetActivations",
				papiResponse: getActivations(),
				error:        nil,
				stubOnce:     false,
			},
			{
				methodName:  "CreateActivation",
				papiRequest: createActivation(),
				papiResponse: &papi.CreateActivationResponse{
					ActivationID: "atv_activation1",
				},
				stubOnce: true,
			},
			{
				methodName: "GetActivation",
				papiRequest: papi.GetActivationRequest{
					PropertyID:   "prp_test",
					ActivationID: "atv_activation1",
				},
				papiResponse: getActivation(),
			},
		})
		useClient(client, func() {
			resource.UnitTest(t, resource.TestCase{
				IsUnitTest: true,
				Providers:  testAccProviders,
				Steps: []resource.TestStep{
					{
						Config: loadFixtureString("testdata/TestPropertyActivation/update_note_field/note_field_doesnt_exist/resource_property_activation.tf"),
					},
					{
						Config:      loadFixtureString("testdata/TestPropertyActivation/update_note_field/note_field_exists/resource_property_activation_update.tf"),
						ExpectError: regexp.MustCompile("cannot update activation attribute note after creation"),
					},
				},
			})
		})

		client.AssertExpectations(t)
	})
}

func getActivation() *papi.GetActivationResponse {
	return &papi.GetActivationResponse{
		GetActivationsResponse: papi.GetActivationsResponse{},
		Activation: &papi.Activation{
			ActivationID:    "atv_activation1",
			PropertyID:      "prp_test",
			PropertyVersion: 1,
			Network:         "STAGING",
			Status:          papi.ActivationStatusActive,
		},
	}
}

func createActivation() papi.CreateActivationRequest {
	return papi.CreateActivationRequest{
		PropertyID: "prp_test",
		Activation: papi.Activation{
			ActivationType:         papi.ActivationTypeDeactivate,
			AcknowledgeAllWarnings: true,
			PropertyVersion:        1,
			Network:                "STAGING",
			NotifyEmails:           []string{"user@example.com"},
		},
	}
}

func getActivations() *papi.GetActivationsResponse {
	return &papi.GetActivationsResponse{
		Activations: papi.ActivationsItems{Items: []*papi.Activation{{
			AccountID:       "act_1-6JHGX",
			ActivationID:    "atv_activation1",
			ActivationType:  "ACTIVATE",
			GroupID:         "grp_91533",
			PropertyName:    "test",
			PropertyID:      "prp_test",
			PropertyVersion: 1,
			Network:         "STAGING",
			Status:          "ACTIVE",
			SubmitDate:      "2020-10-28T15:04:05Z",
		}}}}
}
