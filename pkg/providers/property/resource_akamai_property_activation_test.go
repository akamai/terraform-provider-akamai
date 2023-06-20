package property

import (
	"fmt"
	"regexp"
	"testing"

	"github.com/stretchr/testify/mock"

	"github.com/akamai/AkamaiOPEN-edgegrid-golang/v6/pkg/papi"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/resource"
)

func TestResourcePAPIPropertyActivation(t *testing.T) {
	tests := map[string]struct {
		init  func(*papi.Mock)
		steps []resource.TestStep
	}{
		"property activation lifecycle - OK": {
			init: func(m *papi.Mock) {
				// first step
				// create
				expectGetRuleTree(m, "prp_test", 1, ruleTreeResponseValid, nil).Once()
				expectGetActivations(m, "prp_test", papi.GetActivationsResponse{}, nil).Once()
				expectCreateActivation(m, "prp_test", papi.ActivationTypeActivate, 1, "STAGING",
					[]string{"user@example.com"}, "property activation note for creating", "atv_activation1", true, nil).Once()
				expectGetActivation(m, "prp_test", "atv_activation1", 1, "STAGING", papi.ActivationStatusActive, papi.ActivationTypeActivate, "property activation note for creating", []string{"user@example.com"}, nil).Once()
				// read
				expectGetActivations(m, "prp_test", generateActivationResponseMock("atv_activation1", "property activation note for creating", 1, papi.ActivationTypeActivate, "2020-10-28T15:04:05Z"), nil).Once()

				// second step
				// read
				expectGetActivations(m, "prp_test", generateActivationResponseMock("atv_activation1", "property activation note for creating", 1, papi.ActivationTypeActivate, "2020-10-28T15:04:05Z"), nil).Once()
				// update
				expectGetRuleTree(m, "prp_test", 2, ruleTreeResponseValid, nil).Once()
				expectGetActivations(m, "prp_test", generateActivationResponseMock("atv_activation1", "property activation note for creating", 1, papi.ActivationTypeActivate, "2020-10-28T15:04:05Z"), nil).Once()
				ExpectGetPropertyVersion(m, "prp_test", "", "", 2, papi.VersionStatusInactive, "").Once()
				expectCreateActivation(m, "prp_test", papi.ActivationTypeActivate, 2, "STAGING",
					[]string{"user@example.com"}, "property activation note for updating", "atv_update", true, nil).Once()
				expectGetActivation(m, "prp_test", "atv_update", 2, "STAGING", papi.ActivationStatusActive, papi.ActivationTypeActivate, "property activation note for updating", []string{"user@example.com"}, nil).Once()
				// read
				expectGetActivations(m, "prp_test", activationsResponseSecondVersionIsActive, nil).Once()
				// delete
				expectGetActivations(m, "prp_test", activationsResponseSecondVersionIsActive, nil).Once()
				expectCreateActivation(m, "prp_test", papi.ActivationTypeDeactivate, 2, "STAGING",
					[]string{"user@example.com"}, "property activation note for updating", "atv_update", true, nil).Once()
				expectGetActivation(m, "prp_test", "atv_update", 2, "STAGING", papi.ActivationStatusActive, papi.ActivationTypeDeactivate, "property activation note for updating", []string{"user@example.com"}, nil).Once()
			},
			steps: []resource.TestStep{
				{
					Config: loadFixtureString("./testdata/TestPropertyActivation/ok/resource_property_activation.tf"),
					Check: resource.ComposeAggregateTestCheckFunc(
						resource.TestCheckResourceAttr("akamai_property_activation.test", "id", "prp_test:STAGING"),
						resource.TestCheckResourceAttr("akamai_property_activation.test", "property_id", "prp_test"),
						resource.TestCheckResourceAttr("akamai_property_activation.test", "contact.#", "1"),
						resource.TestCheckResourceAttrSet("akamai_property_activation.test", "contact.0"),
						resource.TestCheckResourceAttr("akamai_property_activation.test", "contact.0", "user@example.com"),
						resource.TestCheckResourceAttr("akamai_property_activation.test", "network", "STAGING"),
						resource.TestCheckResourceAttr("akamai_property_activation.test", "version", "1"),
						resource.TestCheckResourceAttr("akamai_property_activation.test", "warnings", ""),
						resource.TestCheckResourceAttr("akamai_property_activation.test", "errors", ""),
						resource.TestCheckNoResourceAttr("akamai_property_activation.test", "rule_errors"),
						resource.TestCheckResourceAttr("akamai_property_activation.test", "auto_acknowledge_rule_warnings", "true"),
						resource.TestCheckResourceAttr("akamai_property_activation.test", "activation_id", "atv_activation1"),
						resource.TestCheckResourceAttr("akamai_property_activation.test", "status", "ACTIVE"),
						resource.TestCheckResourceAttr("akamai_property_activation.test", "note", "property activation note for creating"),
					),
				},
				{
					Config: loadFixtureString("./testdata/TestPropertyActivation/ok/resource_property_activation_update.tf"),
					Check: resource.ComposeAggregateTestCheckFunc(
						resource.TestCheckResourceAttr("akamai_property_activation.test", "id", "prp_test:STAGING"),
						resource.TestCheckResourceAttr("akamai_property_activation.test", "contact.#", "1"),
						resource.TestCheckResourceAttrSet("akamai_property_activation.test", "contact.0"),
						resource.TestCheckResourceAttr("akamai_property_activation.test", "contact.0", "user@example.com"),
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
				expectGetActivation(m, "prp_test", "atv_activation1", 1, "STAGING", papi.ActivationStatusActive, papi.ActivationTypeActivate, "property activation note for creating", []string{"user@example.com"}, nil).Once()
				// read
				expectGetActivations(m, "prp_test", generateActivationResponseMock("atv_activation1", "property activation note for creating", 1, papi.ActivationTypeActivate, "2020-10-28T15:04:05Z"), nil).Once()
				// delete
				expectGetActivations(m, "prp_test", generateActivationResponseMock("atv_activation1", "property activation note for creating", 1, papi.ActivationTypeActivate, "2020-10-28T15:04:05Z"), nil).Once()
				m.On(
					"CreateActivation",
					mock.Anything,
					papi.CreateActivationRequest{
						PropertyID: "prp_test",
						Activation: papi.Activation{
							ActivationType:         papi.ActivationTypeDeactivate,
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
				expectGetActivation(m, "prp_test", "atv_activation1", 1, "STAGING", papi.ActivationStatusActive, papi.ActivationTypeDeactivate, "property activation note for creating", []string{"user@example.com"}, nil).Once()
			},
			steps: []resource.TestStep{
				{
					Config: loadFixtureString("./testdata/TestPropertyActivation/ok/resource_property_activation_with_compliance_record.tf"),
					Check: resource.ComposeAggregateTestCheckFunc(
						resource.TestCheckResourceAttr("akamai_property_activation.test", "id", "prp_test:STAGING"),
						resource.TestCheckResourceAttr("akamai_property_activation.test", "property_id", "prp_test"),
						resource.TestCheckResourceAttr("akamai_property_activation.test", "contact.#", "1"),
						resource.TestCheckResourceAttrSet("akamai_property_activation.test", "contact.0"),
						resource.TestCheckResourceAttr("akamai_property_activation.test", "contact.0", "user@example.com"),
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
				expectGetActivations(m, "prp_test", papi.GetActivationsResponse{}, nil).Once()
				expectCreateActivation(m, "prp_test", papi.ActivationTypeActivate, 1, "STAGING",
					[]string{"user@example.com"}, "property activation note for creating", "atv_activation1", true, nil).Once()
				expectGetActivation(m, "prp_test", "atv_activation1", 1, "STAGING", papi.ActivationStatusActive, papi.ActivationTypeActivate, "property activation note for creating", []string{"user@example.com"}, nil).Once()
				// read
				expectGetActivations(m, "prp_test", generateActivationResponseMock("atv_activation1", "property activation note for creating", 1, papi.ActivationTypeActivate, "2020-10-28T15:04:05Z"), nil).Once()
				// delete
				expectGetActivations(m, "prp_test", generateActivationResponseMock("atv_activation1", "property activation note for creating", 1, papi.ActivationTypeActivate, "2020-10-28T15:04:05Z"), nil).Once()
				expectCreateActivation(m, "prp_test", papi.ActivationTypeDeactivate, 1, "STAGING",
					[]string{"user@example.com"}, "property activation note for creating", "atv_update", true, nil).Once()
				expectGetActivation(m, "prp_test", "atv_update", 1, "STAGING", papi.ActivationStatusActive, papi.ActivationTypeDeactivate, "property activation note for creating", []string{"user@example.com"}, nil).Once()

			},
			steps: []resource.TestStep{
				{
					Config: loadFixtureString("./testdata/TestPropertyActivation/ok/resource_property_activation.tf"),
					Check: resource.ComposeAggregateTestCheckFunc(
						resource.TestCheckResourceAttr("akamai_property_activation.test", "id", "prp_test:STAGING"),
						resource.TestCheckResourceAttr("akamai_property_activation.test", "property_id", "prp_test"),
						resource.TestCheckResourceAttr("akamai_property_activation.test", "contact.#", "1"),
						resource.TestCheckResourceAttrSet("akamai_property_activation.test", "contact.0"),
						resource.TestCheckResourceAttr("akamai_property_activation.test", "contact.0", "user@example.com"),
						resource.TestCheckResourceAttr("akamai_property_activation.test", "network", "STAGING"),
						resource.TestCheckResourceAttr("akamai_property_activation.test", "version", "1"),
						resource.TestCheckResourceAttr("akamai_property_activation.test", "auto_acknowledge_rule_warnings", "true"),
						resource.TestCheckResourceAttr("akamai_property_activation.test", "warnings", ""),
						resource.TestCheckResourceAttr("akamai_property_activation.test", "errors", ""),
						resource.TestCheckResourceAttr("akamai_property_activation.test", "activation_id", "atv_activation1"),
						resource.TestCheckResourceAttr("akamai_property_activation.test", "status", "ACTIVE"),
						resource.TestCheckResourceAttr("akamai_property_activation.test", "note", "property activation note for creating"),
					),
				},
			},
		},
		"schema with minimum attributes - OK": {
			init: func(m *papi.Mock) {
				// create
				expectGetRuleTree(m, "prp_test", 1, ruleTreeResponseValid, nil).Once()
				expectGetActivations(m, "prp_test", papi.GetActivationsResponse{}, nil).Once()
				expectCreateActivation(m, "prp_test", papi.ActivationTypeActivate, 1, "STAGING",
					[]string{"user@example.com"}, "", "atv_activation1", false, nil).Once()
				expectGetActivation(m, "prp_test", "atv_activation1", 1, "STAGING", papi.ActivationStatusActive, papi.ActivationTypeActivate, "", []string{"user@example.com"}, nil).Once()
				// read
				expectGetActivations(m, "prp_test", generateActivationResponseMock("atv_activation1", "", 1, papi.ActivationTypeActivate, "2020-10-28T15:04:05Z"), nil).Once()
				// delete
				expectGetActivations(m, "prp_test", generateActivationResponseMock("atv_activation1", "", 1, papi.ActivationTypeActivate, "2020-10-28T15:04:05Z"), nil).Once()
				expectCreateActivation(m, "prp_test", papi.ActivationTypeDeactivate, 1, "STAGING",
					[]string{"user@example.com"}, "", "atv_update", false, nil).Once()
				expectGetActivation(m, "prp_test", "atv_update", 1, "STAGING", papi.ActivationStatusActive, papi.ActivationTypeDeactivate, "", []string{"user@example.com"}, nil).Once()
			},
			steps: []resource.TestStep{
				{
					Config: loadFixtureString("testdata/TestPropertyActivation/ok/resource_property_activation_minimum_args.tf"),
					Check: resource.ComposeAggregateTestCheckFunc(
						resource.TestCheckResourceAttr("akamai_property_activation.test", "id", "prp_test:STAGING"),
						resource.TestCheckResourceAttr("akamai_property_activation.test", "property_id", "prp_test"),
						resource.TestCheckResourceAttr("akamai_property_activation.test", "contact.#", "1"),
						resource.TestCheckResourceAttrSet("akamai_property_activation.test", "contact.0"),
						resource.TestCheckResourceAttr("akamai_property_activation.test", "contact.0", "user@example.com"),
						resource.TestCheckResourceAttr("akamai_property_activation.test", "network", "STAGING"),
						resource.TestCheckResourceAttr("akamai_property_activation.test", "version", "1"),
						resource.TestCheckResourceAttr("akamai_property_activation.test", "warnings", ""),
						resource.TestCheckResourceAttr("akamai_property_activation.test", "errors", ""),
						resource.TestCheckResourceAttr("akamai_property_activation.test", "activation_id", "atv_activation1"),
						resource.TestCheckResourceAttr("akamai_property_activation.test", "status", "ACTIVE"),
					),
				},
			},
		},
		"property activation is deactivated from other source (UI)": {
			init: func(m *papi.Mock) {
				// first step
				// create
				expectGetRuleTree(m, "prp_test", 1, ruleTreeResponseValid, nil).Once()
				expectGetActivations(m, "prp_test", papi.GetActivationsResponse{}, nil).Once()
				expectCreateActivation(m, "prp_test", papi.ActivationTypeActivate, 1, "STAGING", []string{"user@example.com"}, "", "atv_activation1", true, nil).Once()
				expectGetActivation(m, "prp_test", "atv_activation1", 1, "STAGING", papi.ActivationStatusActive, papi.ActivationTypeActivate, "", []string{"user@example.com"}, nil).Once()
				// read
				expectGetActivations(m, "prp_test", generateActivationResponseMock("atv_activation1", "", 1, papi.ActivationTypeActivate, "2020-10-28T15:04:05Z"), nil).Once()

				// second step
				// no changes in configuration, but it was deactivated in other source, for example on UI -> terraform cleans state and activate this version again
				// read
				expectGetActivations(m, "prp_test", activationsResponseDeactivated, nil).Once()
				// create
				expectGetRuleTree(m, "prp_test", 1, ruleTreeResponseValid, nil).Once()
				expectGetActivations(m, "prp_test", activationsResponseDeactivated, nil).Once()
				expectCreateActivation(m, "prp_test", papi.ActivationTypeActivate, 1, "STAGING", []string{"user@example.com"}, "", "atv_activation1", true, nil).Once()
				expectGetActivation(m, "prp_test", "atv_activation1", 1, "STAGING", papi.ActivationStatusActive, papi.ActivationTypeActivate, "", []string{"user@example.com"}, nil).Once()
				// read
				expectGetActivations(m, "prp_test", generateActivationResponseMock("atv_activation1", "", 1, papi.ActivationTypeActivate, "2020-10-28T15:04:05Z"), nil).Once()
				// delete
				expectGetActivations(m, "prp_test", generateActivationResponseMock("atv_activation1", "", 1, papi.ActivationTypeActivate, "2020-10-28T15:04:05Z"), nil).Once()
				expectCreateActivation(m, "prp_test", papi.ActivationTypeDeactivate, 1, "STAGING", []string{"user@example.com"}, "", "atv_update", true, nil).Once()
				expectGetActivation(m, "prp_test", "atv_update", 1, "STAGING", papi.ActivationStatusActive, papi.ActivationTypeDeactivate, "", []string{"user@example.com"}, nil).Once()
			},
			steps: []resource.TestStep{
				{
					Config: loadFixtureString("./testdata/TestPropertyActivation/deactivated_in_other_source/resource_property_activation.tf"),
					Check:  resource.TestCheckResourceAttr("akamai_property_activation.test", "version", "1"),
				},
				{
					Config: loadFixtureString("./testdata/TestPropertyActivation/deactivated_in_other_source/resource_property_activation.tf"),
					Check:  resource.TestCheckResourceAttr("akamai_property_activation.test", "version", "1"),
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
					ExpectError: regexp.MustCompile("Missing required argument"),
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
				expectGetActivations(m, "prp_test", generateActivationResponseMock("atv_activation1", "", 1, papi.ActivationTypeActivate, "2020-10-28T15:04:05Z"), nil).Once()
				// read twice
				expectGetActivations(m, "prp_test", generateActivationResponseMock("atv_activation1", "", 1, papi.ActivationTypeActivate, "2020-10-28T15:04:05Z"), nil).Twice()
				// update
				expectGetRuleTree(m, "prp_test", 1, ruleTreeResponseValid, nil).Once()
				expectGetActivations(m, "prp_test", generateActivationResponseMock("atv_activation1", "", 1, papi.ActivationTypeActivate, "2020-10-28T15:04:05Z"), nil).Once()
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
				expectGetActivations(m, "prp_test", generateActivationResponseMock("atv_activation1", "property activation note", 1, papi.ActivationTypeActivate, "2020-10-28T15:04:05Z"), nil).Once()
				// read twice
				expectGetActivations(m, "prp_test", generateActivationResponseMock("atv_activation1", "property activation note", 1, papi.ActivationTypeActivate, "2020-10-28T15:04:05Z"), nil).Twice()
				// update
				expectGetRuleTree(m, "prp_test", 1, ruleTreeResponseValid, nil).Once()
				expectGetActivations(m, "prp_test", generateActivationResponseMock("atv_activation1", "property activation note", 1, papi.ActivationTypeActivate, "2020-10-28T15:04:05Z"), nil).Once()
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
		"500 error on property activation update": {
			init: func(m *papi.Mock) {
				// first step
				// create
				expectGetRuleTree(m, "prp_test", 1, ruleTreeResponseValid, nil).Once()
				expectGetActivations(m, "prp_test", papi.GetActivationsResponse{}, nil).Once()
				expectCreateActivation(m, "prp_test", papi.ActivationTypeActivate, 1, "STAGING", []string{"user@example.com"}, "", "atv_activation1", true, nil).Once()
				expectGetActivation(m, "prp_test", "atv_activation1", 1, "STAGING", papi.ActivationStatusActive, papi.ActivationTypeActivate, "", []string{"user@example.com"}, nil).Once()
				// read
				expectGetActivations(m, "prp_test", generateActivationResponseMock("atv_activation1", "", 1, papi.ActivationTypeActivate, "2020-10-28T15:04:05Z"), nil).Once()

				// second step
				// read
				expectGetActivations(m, "prp_test", generateActivationResponseMock("atv_activation1", "", 1, papi.ActivationTypeActivate, "2020-10-28T15:04:05Z"), nil).Once()
				// update
				expectGetRuleTree(m, "prp_test", 2, ruleTreeResponseValid, nil).Once()
				expectGetActivations(m, "prp_test", generateActivationResponseMock("atv_activation1", "", 1, papi.ActivationTypeActivate, "2020-10-28T15:04:05Z"), nil).Once()
				ExpectGetPropertyVersion(m, "prp_test", "", "", 2, papi.VersionStatusInactive, "").Once()
				// error on update
				m.On("CreateActivation", AnyCTX, papi.CreateActivationRequest{
					PropertyID: "prp_test",
					Activation: papi.Activation{
						ActivationType:         papi.ActivationTypeActivate,
						Network:                papi.ActivationNetworkStaging,
						PropertyVersion:        2,
						NotifyEmails:           []string{"user@example.com"},
						AcknowledgeAllWarnings: true,
					},
				}).Return(nil, fmt.Errorf("some 500 error")).Once()
				// delete - terraform clean up after error is occurred
				expectGetActivations(m, "prp_test", generateActivationResponseMock("atv_activation1", "", 1, papi.ActivationTypeActivate, "2020-10-28T15:04:05Z"), nil).Once()
				expectCreateActivation(m, "prp_test", papi.ActivationTypeDeactivate, 2, "STAGING", []string{"user@example.com"}, "", "atv_update", true, nil).Once()
				expectGetActivation(m, "prp_test", "atv_update", 2, "STAGING", papi.ActivationStatusActive, papi.ActivationTypeDeactivate, "", []string{"user@example.com"}, nil).Once()

			},
			steps: []resource.TestStep{
				{
					Config: loadFixtureString("./testdata/TestPropertyActivation/500_on_activation/resource_property_activation.tf"),
					Check:  resource.TestCheckResourceAttr("akamai_property_activation.test", "version", "1"),
				},
				{
					Config:      loadFixtureString("./testdata/TestPropertyActivation/500_on_activation/resource_property_activation_update.tf"),
					Check:       resource.TestCheckResourceAttr("akamai_property_activation.test", "version", "1"),
					ExpectError: regexp.MustCompile("some 500 error"),
				},
			},
		},
		"property activation import": {
			init: func(m *papi.Mock) {
				// create
				expectGetRuleTree(m, "prp_test", 1, ruleTreeResponseValid, nil).Once()
				expectGetActivations(m, "prp_test", papi.GetActivationsResponse{}, nil).Once()
				expectCreateActivation(m, "prp_test", papi.ActivationTypeActivate, 1, "STAGING",
					[]string{"user@example.com"}, "property activation note for importing", "atv_activation1", false, nil).Once()
				expectGetActivation(m, "prp_test", "atv_activation1", 1, "STAGING", papi.ActivationStatusActive, papi.ActivationTypeActivate, "property activation note for importing", []string{"user@example.com"}, nil).Once()
				expectGetActivation(m, "prp_test", "atv_activation1", 1, "STAGING", papi.ActivationStatusActive, papi.ActivationTypeActivate, "property activation note for importing", []string{"user@example.com"}, nil).Once()
				// read
				expectGetActivations(m, "prp_test", generateActivationResponseMock("atv_activation1", "property activation note for importing", 1, papi.ActivationTypeActivate, "2020-10-28T15:04:05Z"), nil).Once()
				// 2nd read for import
				expectGetActivations(m, "prp_test", generateActivationResponseMock("atv_activation1", "property activation note for importing", 1, papi.ActivationTypeActivate, "2020-10-28T15:04:05Z"), nil).Once()
				// delete
				expectGetActivations(m, "prp_test", generateActivationResponseMock("atv_activation1", "property activation note for importing", 1, papi.ActivationTypeActivate, "2020-10-28T15:04:05Z"), nil).Once()
				expectCreateActivation(m, "prp_test", papi.ActivationTypeDeactivate, 1, "STAGING",
					[]string{"user@example.com"}, "property activation note for importing", "atv_activation1", false, nil).Once()
			},
			steps: []resource.TestStep{
				{
					Config: loadFixtureString("./testdata/TestPropertyActivation/import/resource_property_activation_creation_for_import.tf"),
					Check: resource.ComposeAggregateTestCheckFunc(
						resource.TestCheckResourceAttr("akamai_property_activation.test", "id", "prp_test:STAGING"),
						resource.TestCheckResourceAttr("akamai_property_activation.test", "property_id", "prp_test"),
						resource.TestCheckResourceAttr("akamai_property_activation.test", "contact.#", "1"),
						resource.TestCheckResourceAttrSet("akamai_property_activation.test", "contact.0"),
						resource.TestCheckResourceAttr("akamai_property_activation.test", "contact.0", "user@example.com"),
						resource.TestCheckResourceAttr("akamai_property_activation.test", "network", "STAGING"),
						resource.TestCheckResourceAttr("akamai_property_activation.test", "version", "1"),
						resource.TestCheckResourceAttr("akamai_property_activation.test", "auto_acknowledge_rule_warnings", "false"),
						resource.TestCheckResourceAttr("akamai_property_activation.test", "warnings", ""),
						resource.TestCheckNoResourceAttr("akamai_property_activation.test", "rule_warnings"),
						resource.TestCheckResourceAttr("akamai_property_activation.test", "errors", ""),
						resource.TestCheckResourceAttr("akamai_property_activation.test", "activation_id", "atv_activation1"),
						resource.TestCheckResourceAttr("akamai_property_activation.test", "status", "ACTIVE"),
						resource.TestCheckResourceAttr("akamai_property_activation.test", "note", "property activation note for importing"),
					),
				},
				{
					ImportState:       true,
					ImportStateId:     "prp_test:STAGING",
					ResourceName:      "akamai_property_activation.test",
					ImportStateVerify: true,
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
					ProtoV5ProviderFactories: testAccProviders,
					IsUnitTest:               true,
					Steps:                    test.steps,
				})
			})
			client.AssertExpectations(t)
		})
	}
}

func generateActivationResponseMock(activationID string, note string, version int, activationType papi.ActivationType, date string) papi.GetActivationsResponse {
	return papi.GetActivationsResponse{
		Activations: papi.ActivationsItems{Items: append([]*papi.Activation{}, generateActivationItemMock(activationID, note, version, activationType, date))},
	}
}

func generateActivationItemMock(activationID string, note string, version int, activationType papi.ActivationType, date string) *papi.Activation {
	return &papi.Activation{
		AccountID:       "act_1-6JHGX",
		ActivationID:    activationID,
		ActivationType:  activationType,
		GroupID:         "grp_91533",
		PropertyName:    "test",
		PropertyID:      "prp_test",
		PropertyVersion: version,
		Network:         "STAGING",
		Status:          "ACTIVE",
		SubmitDate:      date,
		UpdateDate:      date,
		NotifyEmails:    []string{"user@example.com"},
		Note:            note,
	}
}

var (
	mockActivationsListForDeactivation = append([]*papi.Activation{}, mockDeactivation, mockActivationCreation)
	mockActivationsListForLifecycle    = append([]*papi.Activation{}, mockDeactivationForLifecycle, mockActivationUpdate)
	mockDeactivation                   = generateActivationItemMock("atv_activation1", "property activation note for creating", 1, papi.ActivationTypeDeactivate, "2020-10-28T15:04:05Z")
	mockDeactivationForLifecycle       = generateActivationItemMock("atv_activation1", "property activation note for creating", 1, papi.ActivationTypeDeactivate, "2020-10-28T14:04:05Z")
	mockActivationCreation             = generateActivationItemMock("atv_activation1", "property activation note for creating", 1, papi.ActivationTypeActivate, "2020-10-28T15:04:05Z")
	mockActivationUpdate               = generateActivationItemMock("atv_update", "property activation note for updating", 2, papi.ActivationTypeActivate, "2020-10-28T15:04:05Z")
	ruleTreeResponseValid              = papi.GetRuleTreeResponse{
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

	activationsResponseDeactivated = papi.GetActivationsResponse{
		Activations: papi.ActivationsItems{Items: mockActivationsListForDeactivation},
	}
	activationsResponseSecondVersionIsActive = papi.GetActivationsResponse{
		Activations: papi.ActivationsItems{Items: mockActivationsListForLifecycle},
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
		network papi.ActivationNetwork, notify []string, note string, activationID string, acknowledgeAllWarnings bool, err error) *mock.Call {
		if err != nil {
			return m.On(
				"CreateActivation",
				mock.Anything,
				papi.CreateActivationRequest{
					PropertyID: propertyID,
					Activation: papi.Activation{
						ActivationType:         activationType,
						AcknowledgeAllWarnings: acknowledgeAllWarnings,
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
					AcknowledgeAllWarnings: acknowledgeAllWarnings,
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
		network papi.ActivationNetwork, status papi.ActivationStatus, actType papi.ActivationType, note string, contact []string, err error) *mock.Call {
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
				ActivationType:  actType,
				Note:            note,
				NotifyEmails:    contact,
			},
		}, nil)
	}
)
