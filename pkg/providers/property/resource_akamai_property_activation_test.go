package property

import (
	"fmt"
	"regexp"
	"testing"

	"github.com/stretchr/testify/mock"

	"github.com/akamai/AkamaiOPEN-edgegrid-golang/v5/pkg/papi"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/resource"
)

func TestResourcePAPIPropertyActivation(t *testing.T) {
	tests := map[string]struct {
		init  func(*papi.Mock)
		steps []resource.TestStep
	}{
		"property activation lifecycle - OK": {
			init: func(m *papi.Mock) {
				// create
				expectGetRuleTree(m, "prp_test", 1, ruleTreeResponseValid, nil).Once()
				expectGetActivations(m, "prp_test", papi.GetActivationsResponse{}, nil).Once()
				expectCreateActivation(m, "prp_test", papi.ActivationTypeActivate, 1, "STAGING",
					[]string{"user@example.com"}, "property activation note for creating", "atv_activation1", nil).Once()
				expectGetActivation(m, "prp_test", "atv_activation1", 1, "STAGING", papi.ActivationStatusActive, nil).Once()
				// read
				expectGetActivations(m, "prp_test", activationsResponseDeactivated, nil).Twice()
				// update
				expectGetRuleTree(m, "prp_test", 2, ruleTreeResponseValid, nil).Once()
				expectGetActivations(m, "prp_test", activationsResponseDeactivated, nil).Once()
				ExpectGetPropertyVersion(m, "prp_test", "", "", 2, papi.VersionStatusInactive, "").Once()
				expectCreateActivation(m, "prp_test", papi.ActivationTypeActivate, 2, "STAGING",
					[]string{"user@example.com"}, "property activation note for updating", "atv_update", nil).Once()
				expectGetActivation(m, "prp_test", "atv_update", 2, "STAGING", papi.ActivationStatusActive, nil).Once()
				// read and delete
				expectGetActivations(m, "prp_test", activationsResponseUpdated, nil).Twice()
			},
			steps: []resource.TestStep{
				{
					Config: loadFixtureString("./testdata/TestPropertyActivation/ok/resource_property_activation.tf"),
					Check: resource.ComposeAggregateTestCheckFunc(
						resource.TestCheckResourceAttr("akamai_property_activation.test", "id", "prp_test:STAGING"),
						resource.TestCheckResourceAttr("akamai_property_activation.test", "property_id", "prp_test"),
						resource.TestCheckResourceAttr("akamai_property_activation.test", "network", "STAGING"),
						resource.TestCheckResourceAttr("akamai_property_activation.test", "version", "1"),
						resource.TestCheckResourceAttr("akamai_property_activation.test", "warnings", ""),
						resource.TestCheckResourceAttr("akamai_property_activation.test", "errors", ""),
						resource.TestCheckNoResourceAttr("akamai_property_activation.test", "rule_errors"),
						resource.TestCheckResourceAttr("akamai_property_activation.test", "auto_acknowledge_rule_warnings", "true"),
						resource.TestCheckNoResourceAttr("akamai_property_activation.test", "rule_warnings"),
						resource.TestCheckResourceAttr("akamai_property_activation.test", "activation_id", "atv_activation1"),
						resource.TestCheckResourceAttr("akamai_property_activation.test", "status", "ACTIVE"),
						resource.TestCheckResourceAttr("akamai_property_activation.test", "note", "property activation note for creating"),
					),
				},
				{
					Config: loadFixtureString("./testdata/TestPropertyActivation/ok/resource_property_activation_update.tf"),
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
		},
		"check property activation with compliance record - OK": {
			init: func(m *papi.Mock) {
				// create
				expectGetRuleTree(m, "prp_test", 1, ruleTreeResponseValid, nil).Once()
				expectGetActivations(m, "prp_test", papi.GetActivationsResponse{}, nil).Once()
				// Create with compliance record
				m.On(
					"CreateActivation",
					mock.Anything,
					papi.CreateActivationRequest{
						PropertyID: "prp_test",
						Activation: papi.Activation{
							ActivationType:         papi.ActivationTypeActivate,
							AcknowledgeAllWarnings: true,
							PropertyVersion:        1,
							Network:                "STAGING",
							NotifyEmails:           []string{"user@example.com"},
							Note:                   "property activation note for creating",
							ComplianceRecord: &papi.ComplianceRecordNone{
								CustomerEmail:  "user@example.com",
								PeerReviewedBy: "user1@example.com",
								UnitTested:     true,
							},
						},
					},
				).Return(&papi.CreateActivationResponse{
					ActivationID: "atv_activation1",
				}, nil).Once()

				expectGetActivation(m, "prp_test", "atv_activation1", 1, "STAGING", papi.ActivationStatusActive, nil).Once()
				// read and delete
				expectGetActivations(m, "prp_test", activationsResponseDeactivated, nil).Twice()
			},
			steps: []resource.TestStep{
				{
					Config: loadFixtureString("./testdata/TestPropertyActivation/ok/resource_property_activation_with_compliance_record.tf"),
					Check: resource.ComposeAggregateTestCheckFunc(
						resource.TestCheckResourceAttr("akamai_property_activation.test", "id", "prp_test:STAGING"),
						resource.TestCheckResourceAttr("akamai_property_activation.test", "property_id", "prp_test"),
						resource.TestCheckResourceAttr("akamai_property_activation.test", "network", "STAGING"),
						resource.TestCheckResourceAttr("akamai_property_activation.test", "version", "1"),
						resource.TestCheckResourceAttr("akamai_property_activation.test", "auto_acknowledge_rule_warnings", "true"),
						resource.TestCheckResourceAttr("akamai_property_activation.test", "warnings", ""),
						resource.TestCheckResourceAttr("akamai_property_activation.test", "errors", ""),
						resource.TestCheckResourceAttr("akamai_property_activation.test", "activation_id", "atv_activation1"),
						resource.TestCheckResourceAttr("akamai_property_activation.test", "status", "ACTIVE"),
						resource.TestCheckResourceAttr("akamai_property_activation.test", "compliance_record.0.noncompliance_reason_none.0.unit_tested", "true"),
						resource.TestCheckResourceAttr("akamai_property_activation.test", "note", "property activation note for creating"),
					),
				},
			},
		},
		"check schema property activation - OK": {
			init: func(m *papi.Mock) {
				// create
				expectGetRuleTree(m, "prp_test", 1, ruleTreeResponseValid, nil).Once()
				expectGetActivations(m, "prp_test", activationsResponseActivated, nil).Once()
				// read and delete
				expectGetActivations(m, "prp_test", activationsResponseDeactivated, nil).Twice()
			},
			steps: []resource.TestStep{
				{
					Config: loadFixtureString("./testdata/TestPropertyActivation/ok/resource_property_activation.tf"),
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
		},
		"schema with `property` instead of `property_id` - OK": {
			init: func(m *papi.Mock) {
				// create
				expectGetRuleTree(m, "prp_test", 1, ruleTreeResponseValid, nil).Once()
				expectGetActivations(m, "prp_test", activationsResponseActivated, nil).Once()
				// read and delete
				expectGetActivations(m, "prp_test", activationsResponseDeactivated, nil).Twice()
			},
			steps: []resource.TestStep{
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
		},
		"check schema property activation compliance record - error empty compliance_record block": {
			steps: []resource.TestStep{
				{
					Config:      loadFixtureString("./testdata/TestPropertyActivation/cr_validation/resource_property_activation_with_empty_cr.tf"),
					ExpectError: regexp.MustCompile("one of\n`compliance_record.0.noncompliance_reason_emergency,compliance_record.0.noncompliance_reason_no_production_traffic,compliance_record.0.noncompliance_reason_none,compliance_record.0.noncompliance_reason_other`\nmust be specified"),
				},
			},
		},
		"check schema property activation compliance record - error more than one cr type": {
			steps: []resource.TestStep{
				{
					Config:      loadFixtureString("./testdata/TestPropertyActivation/cr_validation/resource_property_activation_with_more_than_one_cr.tf"),
					ExpectError: regexp.MustCompile("only one of\n`compliance_record.0.noncompliance_reason_emergency,compliance_record.0.noncompliance_reason_no_production_traffic,compliance_record.0.noncompliance_reason_none,compliance_record.0.noncompliance_reason_other`\ncan be specified"),
				},
			},
		},

		"check schema property activation - papi error": {
			init: func(m *papi.Mock) {
				expectGetRuleTree(m, "prp_test", 1, papi.GetRuleTreeResponse{}, fmt.Errorf("failed to create request")).Once()
			},
			steps: []resource.TestStep{
				{
					Config:      loadFixtureString("./testdata/TestPropertyActivation/ok/resource_property_activation.tf"),
					ExpectError: regexp.MustCompile("failed to create request"),
				},
			},
		},
		"check schema property activation - no property id nor property": {
			steps: []resource.TestStep{
				{
					Config:      loadFixtureString("./testdata/TestPropertyActivation/no_propertyId/resource_property_activation.tf"),
					ExpectError: regexp.MustCompile("one of `property,property_id` must be specified"),
				},
			},
		},
		"check schema property activation - no contact": {
			steps: []resource.TestStep{
				{
					Config:      loadFixtureString("./testdata/TestPropertyActivation/no_contact/resource_property_activation.tf"),
					ExpectError: regexp.MustCompile("Missing required argument"),
				},
			},
		},
		"check schema property activation with rule errors": {
			init: func(m *papi.Mock) {
				expectGetRuleTree(m, "prp_test", 1, ruleTreeResponseInvalid, nil).Once()
			},
			steps: []resource.TestStep{
				{
					Config:      loadFixtureString("./testdata/TestPropertyActivation/ok/resource_property_activation.tf"),
					ExpectError: regexp.MustCompile("activation cannot continue due to rule errors"),
				},
			},
		},
		"Note field cannot be added after activation is completed": {
			init: func(m *papi.Mock) {
				// create
				expectGetRuleTree(m, "prp_test", 1, ruleTreeResponseValid, nil).Once()
				expectGetActivations(m, "prp_test", activationsResponseActivated, nil).Once()
				// read twice
				expectGetActivations(m, "prp_test", activationsResponseActivated, nil).Twice()
				// update
				expectGetRuleTree(m, "prp_test", 1, ruleTreeResponseValid, nil).Once()
				expectGetActivations(m, "prp_test", activationsResponseActivated, nil).Once()
				ExpectGetPropertyVersion(m, "prp_test", "", "", 1, papi.VersionStatusActive, "").Once()
				// delete
				expectGetActivations(m, "prp_test", activationsResponseDeactivated, nil).Once()
			},
			steps: []resource.TestStep{
				{
					Config: loadFixtureString("testdata/TestPropertyActivation/update_note_field/note_field_doesnt_exist/resource_property_activation.tf"),
				},
				{
					Config:      loadFixtureString("testdata/TestPropertyActivation/update_note_field/note_field_doesnt_exist/resource_property_activation_update.tf"),
					ExpectError: regexp.MustCompile("cannot update activation attribute note after creation"),
				},
			},
		},
		"Note field cannot be updated after activation is completed": {
			init: func(m *papi.Mock) {
				// create
				expectGetRuleTree(m, "prp_test", 1, ruleTreeResponseValid, nil).Once()
				expectGetActivations(m, "prp_test", activationsResponseActivated, nil).Once()
				// read twice
				expectGetActivations(m, "prp_test", activationsResponseActivated, nil).Twice()
				// update
				expectGetRuleTree(m, "prp_test", 1, ruleTreeResponseValid, nil).Once()
				expectGetActivations(m, "prp_test", activationsResponseActivated, nil).Once()
				ExpectGetPropertyVersion(m, "prp_test", "", "", 1, papi.VersionStatusActive, "").Once()
				// delete
				expectGetActivations(m, "prp_test", activationsResponseDeactivated, nil).Once()
			},
			steps: []resource.TestStep{
				{
					Config: loadFixtureString("./testdata/TestPropertyActivation/update_note_field/note_field_exists/resource_property_activation.tf"),
				},
				{
					Config:      loadFixtureString("./testdata/TestPropertyActivation/update_note_field/note_field_exists/resource_property_activation_update.tf"),
					ExpectError: regexp.MustCompile("cannot update activation attribute note after creation"),
				},
			},
		},
	}

	for name, test := range tests {
		t.Run(name, func(t *testing.T) {
			client := &papi.Mock{}
			if test.init != nil {
				test.init(client)
			}
			useClient(client, nil, func() {
				resource.UnitTest(t, resource.TestCase{
					ProviderFactories: testAccProviders,
					IsUnitTest:        true,
					Steps:             test.steps,
				})
			})
			client.AssertExpectations(t)
		})
	}
}

var (
	ruleTreeResponseValid = papi.GetRuleTreeResponse{
		Response: papi.Response{
			Errors:   make([]*papi.Error, 0),
			Warnings: []*papi.Error{{Title: "some warning"}},
		},
	}
	ruleTreeResponseInvalid = papi.GetRuleTreeResponse{
		Response: papi.Response{
			Errors: []*papi.Error{
				{
					Title: "some error",
				},
			},
		},
	}
	expectGetRuleTree = func(m *papi.Mock, propertyID string, version int, response papi.GetRuleTreeResponse, err error) *mock.Call {
		if err != nil {
			return m.On(
				"GetRuleTree",
				mock.Anything,
				papi.GetRuleTreeRequest{PropertyID: propertyID, PropertyVersion: version, ValidateRules: true},
			).Return(nil, err)
		}
		return m.On(
			"GetRuleTree",
			mock.Anything,
			papi.GetRuleTreeRequest{PropertyID: propertyID, PropertyVersion: version, ValidateRules: true},
		).Return(&response, nil)
	}

	activationsResponseActivated = papi.GetActivationsResponse{
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
		}}},
	}
	activationsResponseDeactivated = papi.GetActivationsResponse{
		Activations: papi.ActivationsItems{Items: []*papi.Activation{{
			AccountID:       "act_1-6JHGX",
			ActivationID:    "atv_activation1",
			ActivationType:  "DEACTIVATE",
			GroupID:         "grp_91533",
			PropertyName:    "test",
			PropertyID:      "prp_test",
			PropertyVersion: 1,
			Network:         "STAGING",
			Status:          "ACTIVE",
			SubmitDate:      "2020-10-28T15:04:05Z",
		}}},
	}
	activationsResponseUpdated = papi.GetActivationsResponse{
		Activations: papi.ActivationsItems{Items: []*papi.Activation{
			{
				AccountID:       "act_1-6JHGX",
				ActivationID:    "atv_activation1",
				ActivationType:  "DEACTIVATE",
				GroupID:         "grp_91533",
				PropertyName:    "test",
				PropertyID:      "prp_test",
				PropertyVersion: 1,
				Network:         "STAGING",
				Status:          "ACTIVE",
				SubmitDate:      "2020-10-28T15:04:05Z",
			},
			{
				AccountID:       "act_1-6JHGX",
				ActivationID:    "atv_update",
				ActivationType:  "DEACTIVATE",
				GroupID:         "grp_91533",
				PropertyName:    "test",
				PropertyID:      "prp_test",
				PropertyVersion: 2,
				Network:         "STAGING",
				Status:          "ACTIVE",
				SubmitDate:      "2020-10-28T15:04:05Z",
			},
		}},
	}
	expectGetActivations = func(m *papi.Mock, propertyID string, response papi.GetActivationsResponse, err error) *mock.Call {
		if err != nil {
			return m.On(
				"GetActivations",
				mock.Anything,
				papi.GetActivationsRequest{PropertyID: propertyID},
			).Return(nil, err)
		}
		return m.On(
			"GetActivations",
			mock.Anything,
			papi.GetActivationsRequest{PropertyID: propertyID},
		).Return(&response, nil)
	}

	expectCreateActivation = func(m *papi.Mock, propertyID string, activationType papi.ActivationType, version int,
		network papi.ActivationNetwork, notify []string, note string, activationID string, err error) *mock.Call {
		if err != nil {
			return m.On(
				"CreateActivation",
				mock.Anything,
				papi.CreateActivationRequest{
					PropertyID: propertyID,
					Activation: papi.Activation{
						ActivationType:         activationType,
						AcknowledgeAllWarnings: true,
						PropertyVersion:        version,
						Network:                network,
						NotifyEmails:           notify,
						Note:                   note,
					},
				},
			).Return(nil, err)
		}
		return m.On(
			"CreateActivation",
			mock.Anything,
			papi.CreateActivationRequest{
				PropertyID: propertyID,
				Activation: papi.Activation{
					ActivationType:         activationType,
					AcknowledgeAllWarnings: true,
					PropertyVersion:        version,
					Network:                network,
					NotifyEmails:           notify,
					Note:                   note,
				},
			},
		).Return(&papi.CreateActivationResponse{
			ActivationID: activationID,
		}, nil)
	}

	expectGetActivation = func(m *papi.Mock, propertyID string, activationID string, version int,
		network papi.ActivationNetwork, status papi.ActivationStatus, err error) *mock.Call {
		if err != nil {
			return m.On(
				"GetActivation",
				mock.Anything,
				papi.GetActivationRequest{
					PropertyID:   propertyID,
					ActivationID: activationID,
				},
			).Return(nil, err)
		}
		return m.On(
			"GetActivation",
			mock.Anything,
			papi.GetActivationRequest{
				PropertyID:   propertyID,
				ActivationID: activationID,
			},
		).Return(&papi.GetActivationResponse{
			GetActivationsResponse: papi.GetActivationsResponse{},
			Activation: &papi.Activation{
				ActivationID:    activationID,
				PropertyID:      propertyID,
				PropertyVersion: version,
				Network:         network,
				Status:          status,
			},
		}, nil)
	}
)
