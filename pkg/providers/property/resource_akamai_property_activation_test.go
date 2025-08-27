package property

import (
	"fmt"
	"regexp"
	"testing"

	"github.com/akamai/AkamaiOPEN-edgegrid-golang/v12/pkg/papi"
	"github.com/akamai/terraform-provider-akamai/v9/pkg/common/test"
	"github.com/akamai/terraform-provider-akamai/v9/pkg/common/testutils"
	"github.com/hashicorp/terraform-plugin-testing/helper/resource"
	"github.com/stretchr/testify/mock"
)

func TestResourcePAPIPropertyActivation(t *testing.T) {
	baseChecker := test.NewStateChecker("akamai_property_activation.test").
		CheckEqual("id", "prp_test:STAGING").
		CheckEqual("property_id", "prp_test").
		CheckEqual("contact.#", "1").
		CheckEqual("contact.0", "user@example.com").
		CheckEqual("network", "STAGING").
		CheckEqual("version", "1").
		CheckEqual("auto_acknowledge_rule_warnings", "true").
		CheckEqual("activation_id", "atv_activation1").
		CheckEqual("status", "ACTIVE").
		CheckEqual("note", "property activation note for creating")

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
				expectGetActivations(m, "prp_test", generateActivationResponseMock("atv_activation1", "property activation note for creating", 1, papi.ActivationTypeActivate, "2020-10-28T15:04:05Z", []string{"user@example.com"}), nil).Once()

				// second step
				// read
				expectGetActivations(m, "prp_test", generateActivationResponseMock("atv_activation1", "property activation note for creating", 1, papi.ActivationTypeActivate, "2020-10-28T15:04:05Z", []string{"user@example.com"}), nil).Once()
				// update
				expectGetRuleTree(m, "prp_test", 2, ruleTreeResponseValid, nil).Once()
				expectGetActivations(m, "prp_test", generateActivationResponseMock("atv_activation1", "property activation note for creating", 1, papi.ActivationTypeActivate, "2020-10-28T15:04:05Z", []string{"user@example.com"}), nil).Once()
				expectGetPropertyVersion(m, "prp_test", "", "", 2, papi.VersionStatusInactive, "").Once()
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
					Config: testutils.LoadFixtureString(t, "./testdata/TestPropertyActivation/ok/resource_property_activation.tf"),
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
						resource.TestCheckResourceAttr("akamai_property_activation.test", "rule_errors.#", "0"),
						resource.TestCheckResourceAttr("akamai_property_activation.test", "auto_acknowledge_rule_warnings", "true"),
						resource.TestCheckResourceAttr("akamai_property_activation.test", "activation_id", "atv_activation1"),
						resource.TestCheckResourceAttr("akamai_property_activation.test", "status", "ACTIVE"),
						resource.TestCheckResourceAttr("akamai_property_activation.test", "note", "property activation note for creating"),
						resource.TestCheckResourceAttr("akamai_property_activation.test", "timeouts.#", "0"),
					),
				},
				{
					Config: testutils.LoadFixtureString(t, "./testdata/TestPropertyActivation/ok/resource_property_activation_update.tf"),
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
						resource.TestCheckResourceAttr("akamai_property_activation.test", "timeouts.#", "0"),
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
					testutils.MockContext,
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
				expectGetActivations(m, "prp_test", generateActivationResponseMock("atv_activation1", "property activation note for creating", 1, papi.ActivationTypeActivate, "2020-10-28T15:04:05Z", []string{"user@example.com"}), nil).Once()
				// delete
				expectGetActivations(m, "prp_test", generateActivationResponseMock("atv_activation1", "property activation note for creating", 1, papi.ActivationTypeActivate, "2020-10-28T15:04:05Z", []string{"user@example.com"}), nil).Once()
				m.On(
					"CreateActivation",
					testutils.MockContext,
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
					Config: testutils.LoadFixtureString(t, "./testdata/TestPropertyActivation/ok/resource_property_activation_with_compliance_record.tf"),
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
				expectGetActivations(m, "prp_test", generateActivationResponseMock("atv_activation1", "property activation note for creating", 1, papi.ActivationTypeActivate, "2020-10-28T15:04:05Z", []string{"user@example.com"}), nil).Once()
				// delete
				expectGetActivations(m, "prp_test", generateActivationResponseMock("atv_activation1", "property activation note for creating", 1, papi.ActivationTypeActivate, "2020-10-28T15:04:05Z", []string{"user@example.com"}), nil).Once()
				expectCreateActivation(m, "prp_test", papi.ActivationTypeDeactivate, 1, "STAGING",
					[]string{"user@example.com"}, "property activation note for creating", "atv_update", true, nil).Once()
				expectGetActivation(m, "prp_test", "atv_update", 1, "STAGING", papi.ActivationStatusActive, papi.ActivationTypeDeactivate, "property activation note for creating", []string{"user@example.com"}, nil).Once()

			},
			steps: []resource.TestStep{
				{
					Config: testutils.LoadFixtureString(t, "./testdata/TestPropertyActivation/ok/resource_property_activation.tf"),
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
		"activation with 500 error - OK": {
			init: func(m *papi.Mock) {
				// create
				expectGetRuleTree(m, "prp_test", 1, ruleTreeResponseValid, nil).Once()
				expectGetActivations(m, "prp_test", papi.GetActivationsResponse{}, nil).Once()

				expectCreateActivation500Err(m, "prp_test", papi.ActivationTypeActivate, 1, "STAGING",
					[]string{"user@example.com"}, "property activation note for creating", "atv_activation1", true, &papi.Error{StatusCode: 500})

				expectGetActivation(m, "prp_test", "atv_activation1", 1, "STAGING", papi.ActivationStatusActive, papi.ActivationTypeActivate, "property activation note for creating", []string{"user@example.com"}, nil).Once()
				// read
				expectGetActivations(m, "prp_test", generateActivationResponseMock("atv_activation1", "property activation note for creating", 1, papi.ActivationTypeActivate, "2020-10-28T15:04:05Z", []string{"user@example.com"}), nil).Once()
				// delete
				expectGetActivations(m, "prp_test", generateActivationResponseMock("atv_activation1", "property activation note for creating", 1, papi.ActivationTypeActivate, "2020-10-28T15:04:05Z", []string{"user@example.com"}), nil).Once()
				expectCreateActivation(m, "prp_test", papi.ActivationTypeDeactivate, 1, "STAGING",
					[]string{"user@example.com"}, "property activation note for creating", "atv_update", true, nil).Once()
				expectGetActivation(m, "prp_test", "atv_update", 1, "STAGING", papi.ActivationStatusActive, papi.ActivationTypeDeactivate, "property activation note for creating", []string{"user@example.com"}, nil).Once()

			},
			steps: []resource.TestStep{
				{
					Config: testutils.LoadFixtureString(t, "./testdata/TestPropertyActivation/ok/resource_property_activation.tf"),
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
				expectGetActivations(m, "prp_test", generateActivationResponseMock("atv_activation1", "", 1, papi.ActivationTypeActivate, "2020-10-28T15:04:05Z", []string{"user@example.com"}), nil).Once()
				// delete
				expectGetActivations(m, "prp_test", generateActivationResponseMock("atv_activation1", "", 1, papi.ActivationTypeActivate, "2020-10-28T15:04:05Z", []string{"user@example.com"}), nil).Once()
				expectCreateActivation(m, "prp_test", papi.ActivationTypeDeactivate, 1, "STAGING",
					[]string{"user@example.com"}, "", "atv_update", false, nil).Once()
				expectGetActivation(m, "prp_test", "atv_update", 1, "STAGING", papi.ActivationStatusActive, papi.ActivationTypeDeactivate, "", []string{"user@example.com"}, nil).Once()
			},
			steps: []resource.TestStep{
				{
					Config: testutils.LoadFixtureString(t, "testdata/TestPropertyActivation/ok/resource_property_activation_minimum_args.tf"),
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
				expectGetActivations(m, "prp_test", generateActivationResponseMock("atv_activation1", "", 1, papi.ActivationTypeActivate, "2020-10-28T15:04:05Z", []string{"user@example.com"}), nil).Once()

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
				expectGetActivations(m, "prp_test", generateActivationResponseMock("atv_activation1", "", 1, papi.ActivationTypeActivate, "2020-10-28T15:04:05Z", []string{"user@example.com"}), nil).Once()
				// delete
				expectGetActivations(m, "prp_test", generateActivationResponseMock("atv_activation1", "", 1, papi.ActivationTypeActivate, "2020-10-28T15:04:05Z", []string{"user@example.com"}), nil).Once()
				expectCreateActivation(m, "prp_test", papi.ActivationTypeDeactivate, 1, "STAGING", []string{"user@example.com"}, "", "atv_update", true, nil).Once()
				expectGetActivation(m, "prp_test", "atv_update", 1, "STAGING", papi.ActivationStatusActive, papi.ActivationTypeDeactivate, "", []string{"user@example.com"}, nil).Once()
			},
			steps: []resource.TestStep{
				{
					Config: testutils.LoadFixtureString(t, "./testdata/TestPropertyActivation/deactivated_in_other_source/resource_property_activation.tf"),
					Check:  resource.TestCheckResourceAttr("akamai_property_activation.test", "version", "1"),
				},
				{
					Config: testutils.LoadFixtureString(t, "./testdata/TestPropertyActivation/deactivated_in_other_source/resource_property_activation.tf"),
					Check:  resource.TestCheckResourceAttr("akamai_property_activation.test", "version", "1"),
				},
			},
		},
		"property activation when other version is already active": {
			init: func(m *papi.Mock) {
				// first step
				// create
				expectGetRuleTree(m, "prp_test", 1, ruleTreeResponseValid, nil).Once()

				firstActivationOnV1 := generateActivationItemMock("atv_activation1", note, 1, papi.ActivationTypeActivate, "2020-09-28T15:04:05Z", []string{"user@example.com"})
				activationOnV2 := generateActivationItemMock("atv_activation2", note, 2, papi.ActivationTypeActivate, "2020-10-28T15:04:05Z", []string{"user@example.com"})
				expectGetActivations(m, "prp_test", papi.GetActivationsResponse{Activations: papi.ActivationsItems{Items: []*papi.Activation{activationOnV2, firstActivationOnV1}}}, nil).Once()
				expectCreateActivation(m, "prp_test", papi.ActivationTypeActivate, 1, "STAGING", []string{"user@example.com"}, "", "atv_activation1", true, nil).Once()
				expectGetActivation(m, "prp_test", "atv_activation1", 1, "STAGING", papi.ActivationStatusActive, papi.ActivationTypeActivate, "", []string{"user@example.com"}, nil).Once()
				// read
				activationOnV1 := generateActivationItemMock("atv_activation1", note, 1, papi.ActivationTypeActivate, "2020-11-28T15:05:05Z", []string{"user@example.com"})
				allActivations := papi.GetActivationsResponse{Activations: papi.ActivationsItems{Items: []*papi.Activation{activationOnV1, activationOnV2, firstActivationOnV1}}}
				expectGetActivations(m, "prp_test", allActivations, nil).Once()

				// delete
				expectGetActivations(m, "prp_test", allActivations, nil).Once()
				expectCreateActivation(m, "prp_test", papi.ActivationTypeDeactivate, 1, "STAGING", []string{"user@example.com"}, "", "atv_update", true, nil).Once()
				expectGetActivation(m, "prp_test", "atv_update", 1, "STAGING", papi.ActivationStatusActive, papi.ActivationTypeDeactivate, "", []string{"user@example.com"}, nil).Once()
			},
			steps: []resource.TestStep{
				{
					Config: testutils.LoadFixtureString(t, "./testdata/TestPropertyActivation/deactivated_in_other_source/resource_property_activation.tf"),
					Check:  resource.TestCheckResourceAttr("akamai_property_activation.test", "version", "1"),
				},
			},
		},
		"activation with custom timeout - lifecycle": {
			init: func(m *papi.Mock) {
				// first step
				// create
				expectGetRuleTree(m, "prp_test", 1, ruleTreeResponseValid, nil).Once()
				expectGetActivations(m, "prp_test", papi.GetActivationsResponse{}, nil).Once()
				expectCreateActivation(m, "prp_test", papi.ActivationTypeActivate, 1, "STAGING",
					[]string{"user@example.com"}, "property activation note for creating", "atv_activation1", true, nil).Once()
				expectGetActivation(m, "prp_test", "atv_activation1", 1, "STAGING", papi.ActivationStatusActive, papi.ActivationTypeActivate, "property activation note for creating", []string{"user@example.com"}, nil).Once()
				// read
				expectGetActivations(m, "prp_test", generateActivationResponseMock("atv_activation1", "property activation note for creating", 1, papi.ActivationTypeActivate, "2020-10-28T15:04:05Z", []string{"user@example.com"}), nil).Once()

				// second step
				// read
				expectGetActivations(m, "prp_test", generateActivationResponseMock("atv_activation1", "property activation note for creating", 1, papi.ActivationTypeActivate, "2020-10-28T15:04:05Z", []string{"user@example.com"}), nil).Once()
				// no update for only timeout change

				// read
				expectGetActivations(m, "prp_test", generateActivationResponseMock("atv_activation1", "property activation note for creating", 1, papi.ActivationTypeActivate, "2020-10-28T15:04:05Z", []string{"user@example.com"}), nil).Once()

				//// delete
				// read
				expectGetActivations(m, "prp_test", generateActivationResponseMock("atv_activation1", "property activation note for creating", 1, papi.ActivationTypeActivate, "2020-10-28T15:04:05Z", []string{"user@example.com"}), nil).Once()
				expectCreateActivation(m, "prp_test", papi.ActivationTypeDeactivate, 1, "STAGING",
					[]string{"user@example.com"}, "property activation note for creating", "atv_activation1", true, nil).Once()
				expectGetActivation(m, "prp_test", "atv_activation1", 1, "STAGING", papi.ActivationStatusActive, papi.ActivationTypeDeactivate, "property activation note for updating", []string{"user@example.com"}, nil).Once()
			},
			steps: []resource.TestStep{
				{
					Config: testutils.LoadFixtureString(t, "./testdata/TestPropertyActivation/timeouts/resource_property_activation_with_timeout.tf"),
					Check: resource.ComposeAggregateTestCheckFunc(
						resource.TestCheckResourceAttr("akamai_property_activation.test", "id", "prp_test:STAGING"),
						resource.TestCheckResourceAttr("akamai_property_activation.test", "property_id", "prp_test"),
						resource.TestCheckResourceAttr("akamai_property_activation.test", "version", "1"),
						resource.TestCheckResourceAttr("akamai_property_activation.test", "status", "ACTIVE"),
						resource.TestCheckResourceAttr("akamai_property_activation.test", "timeouts.#", "1"),
						resource.TestCheckResourceAttr("akamai_property_activation.test", "timeouts.0.default", "2h1m"),
					),
				},
				{
					Config: testutils.LoadFixtureString(t, "./testdata/TestPropertyActivation/timeouts/resource_property_activation_with_timeout_update.tf"),
					Check: resource.ComposeAggregateTestCheckFunc(
						resource.TestCheckResourceAttr("akamai_property_activation.test", "id", "prp_test:STAGING"),
						resource.TestCheckResourceAttr("akamai_property_activation.test", "property_id", "prp_test"),
						resource.TestCheckResourceAttr("akamai_property_activation.test", "version", "1"),
						resource.TestCheckResourceAttr("akamai_property_activation.test", "status", "ACTIVE"),
						resource.TestCheckResourceAttr("akamai_property_activation.test", "timeouts.#", "1"),
						resource.TestCheckResourceAttr("akamai_property_activation.test", "timeouts.0.default", "2h2m"),
					),
				},
			},
		},
		"check schema property activation compliance record - error empty compliance_record block": {
			steps: []resource.TestStep{
				{
					Config:      testutils.LoadFixtureString(t, "./testdata/TestPropertyActivation/cr_validation/resource_property_activation_with_empty_cr.tf"),
					ExpectError: regexp.MustCompile("one of\n`compliance_record.0.noncompliance_reason_emergency,compliance_record.0.noncompliance_reason_no_production_traffic,compliance_record.0.noncompliance_reason_none,compliance_record.0.noncompliance_reason_other`\nmust be specified"),
				},
			},
		},
		"check schema property activation compliance record - error more than one cr type": {
			steps: []resource.TestStep{
				{
					Config:      testutils.LoadFixtureString(t, "./testdata/TestPropertyActivation/cr_validation/resource_property_activation_with_more_than_one_cr.tf"),
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
					Config:      testutils.LoadFixtureString(t, "./testdata/TestPropertyActivation/ok/resource_property_activation.tf"),
					ExpectError: regexp.MustCompile("failed to create request"),
				},
			},
		},
		"check schema property activation - no property id nor property": {
			steps: []resource.TestStep{
				{
					Config:      testutils.LoadFixtureString(t, "./testdata/TestPropertyActivation/no_propertyId/resource_property_activation.tf"),
					ExpectError: regexp.MustCompile("Missing required argument"),
				},
			},
		},
		"check schema property activation - no contact": {
			steps: []resource.TestStep{
				{
					Config:      testutils.LoadFixtureString(t, "./testdata/TestPropertyActivation/no_contact/resource_property_activation.tf"),
					ExpectError: regexp.MustCompile("Missing required argument"),
				},
			},
		},
		"check schema property activation - incorrect timeout duration": {
			steps: []resource.TestStep{
				{
					Config:      testutils.LoadFixtureString(t, "./testdata/TestPropertyActivation/timeouts/resource_property_activation_incorrect_timeout.tf"),
					ExpectError: regexp.MustCompile("provided incorrect duration"),
				},
			},
		},
		"timeout exceeded": {
			init: func(m *papi.Mock) {
				expectGetRuleTree(m, "prp_test", 1, ruleTreeResponseValid, nil).Once()
				expectGetActivations(m, "prp_test", papi.GetActivationsResponse{}, nil).Once()
				expectCreateActivation(m, "prp_test", papi.ActivationTypeActivate, 1, "STAGING",
					[]string{"user@example.com"}, "property activation note for creating", "atv_activation1", true, nil).Once()
				expectGetActivation(m, "prp_test", "atv_activation1", 1, "STAGING", papi.ActivationStatusPending, papi.ActivationTypeActivate, "property activation note for creating", []string{"user@example.com"}, nil)
			},
			steps: []resource.TestStep{
				{
					Config:      testutils.LoadFixtureString(t, "./testdata/TestPropertyActivation/timeouts/resource_property_activation_with_small_timeout.tf"),
					ExpectError: regexp.MustCompile("Timeout waiting for activation status"),
				},
			},
		},
		"check schema property activation with rule errors": {
			init: func(m *papi.Mock) {
				expectGetRuleTree(m, "prp_test", 1, ruleTreeResponseInvalid, nil).Once()
			},
			steps: []resource.TestStep{
				{
					Config:      testutils.LoadFixtureString(t, "./testdata/TestPropertyActivation/ok/resource_property_activation.tf"),
					ExpectError: regexp.MustCompile("activation cannot continue due to rule errors"),
				},
			},
		},
		"Note and contact fields change suppressed when other fields not changed": {
			init: func(m *papi.Mock) {
				// first step
				// create
				expectGetRuleTree(m, "prp_test", 1, ruleTreeResponseValid, nil).Once()
				expectGetActivations(m, "prp_test", papi.GetActivationsResponse{}, nil).Once()
				expectCreateActivation(m, "prp_test", papi.ActivationTypeActivate, 1, "STAGING",
					[]string{"user@example.com"}, "property activation note for creating", "atv_activation1", true, nil).Once()
				expectGetActivation(m, "prp_test", "atv_activation1", 1, "STAGING", papi.ActivationStatusActive, papi.ActivationTypeActivate, "property activation note for creating", []string{"user@example.com"}, nil).Once()
				// read
				expectGetActivations(m, "prp_test", generateActivationResponseMock("atv_activation1", "property activation note for creating", 1, papi.ActivationTypeActivate, "2020-10-28T15:04:05Z", []string{"user@example.com"}), nil).Once()

				// second step
				// read
				expectGetActivations(m, "prp_test", generateActivationResponseMock("atv_activation1", "property activation note for creating", 1, papi.ActivationTypeActivate, "2020-10-28T15:04:05Z", []string{"user@example.com"}), nil).Once()
				// update - try to update only note field, change suppressed
				expectGetActivations(m, "prp_test", generateActivationResponseMock("atv_activation1", "property activation note for creating", 1, papi.ActivationTypeActivate, "2020-10-28T15:04:05Z", []string{"user@example.com"}), nil).Once()
				// read
				// delete
				expectGetActivations(m, "prp_test", generateActivationResponseMock("atv_activation1", "property activation note for creating", 1, papi.ActivationTypeActivate, "2020-10-28T15:04:05Z", []string{"user@example.com"}), nil).Once()
				expectCreateActivation(m, "prp_test", papi.ActivationTypeDeactivate, 1, "STAGING",
					[]string{"user@example.com"}, "property activation note for creating", "atv_activation1", true, nil).Once()
				expectGetActivation(m, "prp_test", "atv_activation1", 1, "STAGING", papi.ActivationStatusActive, papi.ActivationTypeDeactivate, "property activation note for creating", []string{"user@example.com"}, nil).Once()
			},
			steps: []resource.TestStep{
				{
					Config: testutils.LoadFixtureString(t, "./testdata/TestPropertyActivation/diff_suppress/resource_property_activation.tf"),
					Check: resource.ComposeAggregateTestCheckFunc(
						resource.TestCheckResourceAttr("akamai_property_activation.test", "id", "prp_test:STAGING"),
						resource.TestCheckResourceAttr("akamai_property_activation.test", "property_id", "prp_test"),
						resource.TestCheckResourceAttr("akamai_property_activation.test", "contact.#", "1"),
						resource.TestCheckResourceAttrSet("akamai_property_activation.test", "contact.0"),
						resource.TestCheckResourceAttr("akamai_property_activation.test", "contact.0", "user@example.com"),
						resource.TestCheckResourceAttr("akamai_property_activation.test", "network", "STAGING"),
						resource.TestCheckResourceAttr("akamai_property_activation.test", "version", "1"),
						resource.TestCheckResourceAttr("akamai_property_activation.test", "auto_acknowledge_rule_warnings", "true"),
						resource.TestCheckResourceAttr("akamai_property_activation.test", "activation_id", "atv_activation1"),
						resource.TestCheckResourceAttr("akamai_property_activation.test", "status", "ACTIVE"),
						resource.TestCheckResourceAttr("akamai_property_activation.test", "note", "property activation note for creating"),
					),
				},
				{
					Config: testutils.LoadFixtureString(t, "./testdata/TestPropertyActivation/diff_suppress/resource_property_activation_update_suppressed.tf"),
					Check: resource.ComposeAggregateTestCheckFunc(
						resource.TestCheckResourceAttr("akamai_property_activation.test", "id", "prp_test:STAGING"),
						resource.TestCheckResourceAttr("akamai_property_activation.test", "property_id", "prp_test"),
						resource.TestCheckResourceAttr("akamai_property_activation.test", "contact.#", "1"),
						resource.TestCheckResourceAttrSet("akamai_property_activation.test", "contact.0"),
						resource.TestCheckResourceAttr("akamai_property_activation.test", "contact.0", "user@example.com"),
						resource.TestCheckResourceAttr("akamai_property_activation.test", "network", "STAGING"),
						resource.TestCheckResourceAttr("akamai_property_activation.test", "version", "1"),
						resource.TestCheckResourceAttr("akamai_property_activation.test", "auto_acknowledge_rule_warnings", "true"),
						resource.TestCheckResourceAttr("akamai_property_activation.test", "activation_id", "atv_activation1"),
						resource.TestCheckResourceAttr("akamai_property_activation.test", "status", "ACTIVE"),
						resource.TestCheckResourceAttr("akamai_property_activation.test", "note", "property activation note for creating"),
					),
				},
			},
		},
		"compliance_record updated when other fields not changed": {
			init: func(m *papi.Mock) {
				// first step
				// create
				expectGetRuleTree(m, "prp_test", 1, ruleTreeResponseValid, nil).Once()
				expectGetActivations(m, "prp_test", papi.GetActivationsResponse{}, nil).Once()
				expectCreateActivation(m, "prp_test", papi.ActivationTypeActivate, 1, "STAGING",
					[]string{"user@example.com"}, "property activation note for creating", "atv_activation1", true, nil).Once()
				expectGetActivation(m, "prp_test", "atv_activation1", 1, "STAGING", papi.ActivationStatusActive, papi.ActivationTypeActivate, "property activation note for creating", []string{"user@example.com"}, nil).Once()
				// read
				expectGetActivations(m, "prp_test", generateActivationResponseMock("atv_activation1", "property activation note for creating", 1, papi.ActivationTypeActivate, "2020-10-28T15:04:05Z", []string{"user@example.com"}), nil).Once()

				// second step
				// read
				expectGetActivations(m, "prp_test", generateActivationResponseMock("atv_activation1", "property activation note for creating", 1, papi.ActivationTypeActivate, "2020-10-28T15:04:05Z", []string{"user@example.com"}), nil).Once()

				// update
				expectGetActivations(m, "prp_test", generateActivationResponseMock("atv_activation1", "property activation note for creating", 1, papi.ActivationTypeActivate, "2020-10-28T15:04:05Z", []string{"user@example.com"}), nil).Once()

				// delete
				expectGetActivations(m, "prp_test", generateActivationResponseMock("atv_activation1", "property activation note for creating", 1, papi.ActivationTypeActivate, "2020-10-28T15:04:05Z", []string{"user@example.com"}), nil).Once()
				m.On(
					"CreateActivation",
					testutils.MockContext,
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
					Config: testutils.LoadFixtureString(t, "./testdata/TestPropertyActivation/ok/resource_property_activation.tf"),
					Check: resource.ComposeAggregateTestCheckFunc(
						resource.TestCheckResourceAttr("akamai_property_activation.test", "id", "prp_test:STAGING"),
						resource.TestCheckResourceAttr("akamai_property_activation.test", "property_id", "prp_test"),
						resource.TestCheckResourceAttr("akamai_property_activation.test", "contact.#", "1"),
						resource.TestCheckResourceAttrSet("akamai_property_activation.test", "contact.0"),
						resource.TestCheckResourceAttr("akamai_property_activation.test", "contact.0", "user@example.com"),
						resource.TestCheckResourceAttr("akamai_property_activation.test", "network", "STAGING"),
						resource.TestCheckResourceAttr("akamai_property_activation.test", "version", "1"),
						resource.TestCheckResourceAttr("akamai_property_activation.test", "auto_acknowledge_rule_warnings", "true"),
						resource.TestCheckResourceAttr("akamai_property_activation.test", "activation_id", "atv_activation1"),
						resource.TestCheckResourceAttr("akamai_property_activation.test", "status", "ACTIVE"),
						resource.TestCheckResourceAttr("akamai_property_activation.test", "note", "property activation note for creating"),
					),
				},
				{
					Config: testutils.LoadFixtureString(t, "./testdata/TestPropertyActivation/ok/resource_property_activation_with_compliance_record.tf"),
					Check: resource.ComposeAggregateTestCheckFunc(
						resource.TestCheckResourceAttr("akamai_property_activation.test", "id", "prp_test:STAGING"),
						resource.TestCheckResourceAttr("akamai_property_activation.test", "property_id", "prp_test"),
						resource.TestCheckResourceAttr("akamai_property_activation.test", "contact.#", "1"),
						resource.TestCheckResourceAttrSet("akamai_property_activation.test", "contact.0"),
						resource.TestCheckResourceAttr("akamai_property_activation.test", "contact.0", "user@example.com"),
						resource.TestCheckResourceAttr("akamai_property_activation.test", "network", "STAGING"),
						resource.TestCheckResourceAttr("akamai_property_activation.test", "version", "1"),
						resource.TestCheckResourceAttr("akamai_property_activation.test", "auto_acknowledge_rule_warnings", "true"),
						resource.TestCheckResourceAttr("akamai_property_activation.test", "activation_id", "atv_activation1"),
						resource.TestCheckResourceAttr("akamai_property_activation.test", "status", "ACTIVE"),
						resource.TestCheckResourceAttr("akamai_property_activation.test", "note", "property activation note for creating"),
						resource.TestCheckResourceAttr("akamai_property_activation.test", "compliance_record.0.noncompliance_reason_none.0.unit_tested", "true"),
					),
				},
			},
		},
		"Note field change not suppressed when other fields changed": {
			init: func(m *papi.Mock) {
				// first step
				// create
				expectGetRuleTree(m, "prp_test", 1, ruleTreeResponseValid, nil).Once()
				expectGetActivations(m, "prp_test", papi.GetActivationsResponse{}, nil).Once()
				expectCreateActivation(m, "prp_test", papi.ActivationTypeActivate, 1, "STAGING",
					[]string{"user@example.com"}, "property activation note for creating", "atv_activation1", true, nil).Once()
				expectGetActivation(m, "prp_test", "atv_activation1", 1, "STAGING", papi.ActivationStatusActive, papi.ActivationTypeActivate, "property activation note for creating", []string{"user@example.com"}, nil).Once()
				// read
				expectGetActivations(m, "prp_test", generateActivationResponseMock("atv_activation1", "property activation note for creating", 1, papi.ActivationTypeActivate, "2020-10-28T15:04:05Z", []string{"user@example.com"}), nil).Once()

				// second step
				// read
				expectGetActivations(m, "prp_test", generateActivationResponseMock("atv_activation1", "property activation note for creating", 1, papi.ActivationTypeActivate, "2020-10-28T15:04:05Z", []string{"user@example.com"}), nil).Once()
				// update - note field not suppressed update of contact field and version
				expectGetRuleTree(m, "prp_test", 2, ruleTreeResponseValid, nil).Once()
				expectGetActivations(m, "prp_test", generateActivationResponseMock("atv_activation1", "property activation note for creating", 1, papi.ActivationTypeActivate, "2020-10-28T15:04:05Z", []string{"user@example.com"}), nil).Once()
				expectGetPropertyVersion(m, "prp_test", "", "", 2, papi.VersionStatusInactive, "").Once()
				expectCreateActivation(m, "prp_test", papi.ActivationTypeActivate, 2, "STAGING",
					[]string{"user@example.com", "user2@example.com"}, "property activation note for updating", "atv_update", true, nil).Once()
				expectGetActivation(m, "prp_test", "atv_update", 2, "STAGING", papi.ActivationStatusActive, papi.ActivationTypeActivate, "property activation note for updating", []string{"user@example.com", "user2@example.com"}, nil).Once()
				// read
				expectGetActivations(m, "prp_test", activationsResponseWithNoteNoteSuppressed, nil).Once()
				// delete
				expectGetActivations(m, "prp_test", activationsResponseWithNoteNoteSuppressed, nil).Once()
				expectCreateActivation(m, "prp_test", papi.ActivationTypeDeactivate, 2, "STAGING",
					[]string{"user@example.com", "user2@example.com"}, "property activation note for updating", "atv_update", true, nil).Once()
				expectGetActivation(m, "prp_test", "atv_update", 2, "STAGING", papi.ActivationStatusActive, papi.ActivationTypeDeactivate, "property activation note for updating", []string{"user@example.com", "user2@example.com"}, nil).Once()
			},
			steps: []resource.TestStep{
				{
					Config: testutils.LoadFixtureString(t, "./testdata/TestPropertyActivation/diff_suppress/resource_property_activation.tf"),
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
						resource.TestCheckResourceAttr("akamai_property_activation.test", "rule_errors.#", "0"),
						resource.TestCheckResourceAttr("akamai_property_activation.test", "auto_acknowledge_rule_warnings", "true"),
						resource.TestCheckResourceAttr("akamai_property_activation.test", "activation_id", "atv_activation1"),
						resource.TestCheckResourceAttr("akamai_property_activation.test", "status", "ACTIVE"),
						resource.TestCheckResourceAttr("akamai_property_activation.test", "note", "property activation note for creating"),
					),
				},
				{
					Config: testutils.LoadFixtureString(t, "./testdata/TestPropertyActivation/diff_suppress/resource_property_activation_update.tf"),
					Check: resource.ComposeAggregateTestCheckFunc(
						resource.TestCheckResourceAttr("akamai_property_activation.test", "id", "prp_test:STAGING"),
						resource.TestCheckResourceAttr("akamai_property_activation.test", "contact.#", "2"),
						resource.TestCheckResourceAttrSet("akamai_property_activation.test", "contact.0"),
						resource.TestCheckResourceAttrSet("akamai_property_activation.test", "contact.1"),
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
		"property activation lifecycle - contact is correctly read in update when contact list is extended": {
			init: func(m *papi.Mock) {
				// first step
				// create
				expectGetRuleTree(m, "prp_test", 1, ruleTreeResponseValid, nil).Once()
				expectGetActivations(m, "prp_test", papi.GetActivationsResponse{}, nil).Once()
				expectCreateActivation(m, "prp_test", papi.ActivationTypeActivate, 1, "STAGING",
					[]string{"user@example.com"}, "property activation note for creating", "atv_activation1", true, nil).Once()
				expectGetActivation(m, "prp_test", "atv_activation1", 1, "STAGING", papi.ActivationStatusActive, papi.ActivationTypeActivate, "property activation note for creating", []string{"user@example.com"}, nil).Once()
				// read
				expectGetActivations(m, "prp_test", generateActivationResponseMock("atv_activation1", "property activation note for creating", 1, papi.ActivationTypeActivate, "2020-10-28T15:04:05Z", []string{"user@example.com"}), nil).Once()

				// second step
				// read
				expectGetActivations(m, "prp_test", generateActivationResponseMock("atv_activation1", "property activation note for creating", 1, papi.ActivationTypeActivate, "2020-10-28T15:04:05Z", []string{"user@example.com"}), nil).Once()
				// update
				expectGetRuleTree(m, "prp_test", 2, ruleTreeResponseValid, nil).Once()
				expectGetActivations(m, "prp_test", generateActivationResponseMock("atv_activation1", "property activation note for creating", 1, papi.ActivationTypeActivate, "2020-10-28T15:04:05Z", []string{"user@example.com"}), nil).Once()
				expectGetPropertyVersion(m, "prp_test", "", "", 2, papi.VersionStatusInactive, "").Once()
				expectCreateActivation(m, "prp_test", papi.ActivationTypeActivate, 2, "STAGING",
					[]string{"user@example.com", "user2@example.com"}, "property activation note for updating", "atv_update", true, nil).Once()
				expectGetActivation(m, "prp_test", "atv_update", 2, "STAGING", papi.ActivationStatusActive, papi.ActivationTypeActivate, "property activation note for updating", []string{"user@example.com", "user2@example.com"}, nil).Once()
				// read
				getActivationsAfterUpdateResponse := papi.GetActivationsResponse{
					Activations: papi.ActivationsItems{Items: []*papi.Activation{mockDeactivationForLifecycle, generateActivationItemMock("atv_update", "property activation note for updating", 2, papi.ActivationTypeActivate, "2020-10-28T15:04:05Z", []string{"user@example.com", "user2@example.com"})}},
				}
				expectGetActivations(m, "prp_test", getActivationsAfterUpdateResponse, nil).Once()
				// delete
				expectGetActivations(m, "prp_test", getActivationsAfterUpdateResponse, nil).Once()
				expectCreateActivation(m, "prp_test", papi.ActivationTypeDeactivate, 2, "STAGING",
					[]string{"user@example.com", "user2@example.com"}, "property activation note for updating", "atv_update", true, nil).Once()
				expectGetActivation(m, "prp_test", "atv_update", 2, "STAGING", papi.ActivationStatusActive, papi.ActivationTypeDeactivate, "property activation note for updating", []string{"user@example.com", "user2@example.com"}, nil).Once()
			},
			steps: []resource.TestStep{
				{
					Config: testutils.LoadFixtureString(t, "./testdata/TestPropertyActivation/ok/resource_property_activation.tf"),
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
						resource.TestCheckResourceAttr("akamai_property_activation.test", "rule_errors.#", "0"),
						resource.TestCheckResourceAttr("akamai_property_activation.test", "auto_acknowledge_rule_warnings", "true"),
						resource.TestCheckResourceAttr("akamai_property_activation.test", "activation_id", "atv_activation1"),
						resource.TestCheckResourceAttr("akamai_property_activation.test", "status", "ACTIVE"),
						resource.TestCheckResourceAttr("akamai_property_activation.test", "note", "property activation note for creating"),
						resource.TestCheckResourceAttr("akamai_property_activation.test", "timeouts.#", "0"),
					),
				},
				{
					Config: testutils.LoadFixtureString(t, "./testdata/TestPropertyActivation/ok/resource_property_activation_extend_contact.tf"),
					Check: resource.ComposeAggregateTestCheckFunc(
						resource.TestCheckResourceAttr("akamai_property_activation.test", "id", "prp_test:STAGING"),
						resource.TestCheckResourceAttr("akamai_property_activation.test", "contact.#", "2"),
						resource.TestCheckResourceAttrSet("akamai_property_activation.test", "contact.0"),
						resource.TestCheckResourceAttr("akamai_property_activation.test", "contact.0", "user2@example.com"),
						resource.TestCheckResourceAttr("akamai_property_activation.test", "contact.1", "user@example.com"),
						resource.TestCheckResourceAttr("akamai_property_activation.test", "property_id", "prp_test"),
						resource.TestCheckResourceAttr("akamai_property_activation.test", "network", "STAGING"),
						resource.TestCheckResourceAttr("akamai_property_activation.test", "version", "2"),
						resource.TestCheckResourceAttr("akamai_property_activation.test", "activation_id", "atv_update"),
						resource.TestCheckResourceAttr("akamai_property_activation.test", "status", "ACTIVE"),
						resource.TestCheckResourceAttr("akamai_property_activation.test", "note", "property activation note for updating"),
						resource.TestCheckResourceAttr("akamai_property_activation.test", "timeouts.#", "0"),
					),
				},
			},
		},
		"property activation lifecycle - contact is correctly read in update when contact list shrinks": {
			init: func(m *papi.Mock) {
				// first step
				// create
				expectGetRuleTree(m, "prp_test", 2, ruleTreeResponseValid, nil).Once()
				expectGetActivations(m, "prp_test", papi.GetActivationsResponse{}, nil).Once()
				expectCreateActivation(m, "prp_test", papi.ActivationTypeActivate, 2, "STAGING",
					[]string{"user@example.com", "user2@example.com"}, "property activation note for updating", "atv_update", true, nil).Once()
				expectGetActivation(m, "prp_test", "atv_update", 2, "STAGING", papi.ActivationStatusActive, papi.ActivationTypeActivate, "property activation note for updating", []string{"user@example.com", "user2@example.com"}, nil).Once()
				// read
				expectGetActivations(m, "prp_test", generateActivationResponseMock("atv_activation1", "property activation note for updating", 2, papi.ActivationTypeActivate, "2020-10-28T15:04:05Z", []string{"user@example.com", "user2@example.com"}), nil).Once()

				// second step
				// read
				expectGetActivations(m, "prp_test", generateActivationResponseMock("atv_activation1", "property activation note for updating", 2, papi.ActivationTypeActivate, "2020-10-28T15:04:05Z", []string{"user@example.com", "user2@example.com"}), nil).Once()
				// update
				expectGetRuleTree(m, "prp_test", 1, ruleTreeResponseValid, nil).Once()
				expectGetActivations(m, "prp_test", generateActivationResponseMock("atv_activation1", "property activation note for updating", 2, papi.ActivationTypeActivate, "2020-10-28T15:04:05Z", []string{"user@example.com", "user2@example.com"}), nil).Once()
				expectGetPropertyVersion(m, "prp_test", "", "", 1, papi.VersionStatusDeactivated, "").Once()
				expectCreateActivation(m, "prp_test", papi.ActivationTypeActivate, 1, "STAGING",
					[]string{"user@example.com"}, "property activation note for creating", "atv_activation1", true, nil).Once()
				expectGetActivation(m, "prp_test", "atv_activation1", 1, "STAGING", papi.ActivationStatusActive, papi.ActivationTypeActivate, "property activation note for creating", []string{"user@example.com"}, nil).Once()
				// read
				getActivationsAfterUpdateResponse := papi.GetActivationsResponse{
					Activations: papi.ActivationsItems{Items: []*papi.Activation{mockDeactivationForLifecycle, generateActivationItemMock("atv_activation1", "property activation note for creating", 1, papi.ActivationTypeActivate, "2020-10-28T15:04:05Z", []string{"user@example.com"})}},
				}
				expectGetActivations(m, "prp_test", getActivationsAfterUpdateResponse, nil).Once()
				// delete
				expectGetActivations(m, "prp_test", getActivationsAfterUpdateResponse, nil).Once()
				expectCreateActivation(m, "prp_test", papi.ActivationTypeDeactivate, 1, "STAGING",
					[]string{"user@example.com"}, "property activation note for creating", "atv_update", true, nil).Once()
				expectGetActivation(m, "prp_test", "atv_update", 1, "STAGING", papi.ActivationStatusActive, papi.ActivationTypeDeactivate, "property activation note for updating", []string{"user@example.com"}, nil).Once()
			},
			steps: []resource.TestStep{
				{
					Config: testutils.LoadFixtureString(t, "./testdata/TestPropertyActivation/ok/resource_property_activation_extend_contact.tf"),
					Check: resource.ComposeAggregateTestCheckFunc(
						resource.TestCheckResourceAttr("akamai_property_activation.test", "id", "prp_test:STAGING"),
						resource.TestCheckResourceAttr("akamai_property_activation.test", "contact.#", "2"),
						resource.TestCheckResourceAttrSet("akamai_property_activation.test", "contact.0"),
						resource.TestCheckResourceAttr("akamai_property_activation.test", "contact.0", "user2@example.com"),
						resource.TestCheckResourceAttr("akamai_property_activation.test", "contact.1", "user@example.com"),
						resource.TestCheckResourceAttr("akamai_property_activation.test", "property_id", "prp_test"),
						resource.TestCheckResourceAttr("akamai_property_activation.test", "network", "STAGING"),
						resource.TestCheckResourceAttr("akamai_property_activation.test", "version", "2"),
						resource.TestCheckResourceAttr("akamai_property_activation.test", "activation_id", "atv_update"),
						resource.TestCheckResourceAttr("akamai_property_activation.test", "status", "ACTIVE"),
						resource.TestCheckResourceAttr("akamai_property_activation.test", "note", "property activation note for updating"),
						resource.TestCheckResourceAttr("akamai_property_activation.test", "timeouts.#", "0"),
					),
				},
				{
					Config: testutils.LoadFixtureString(t, "./testdata/TestPropertyActivation/ok/resource_property_activation.tf"),
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
						resource.TestCheckResourceAttr("akamai_property_activation.test", "rule_errors.#", "0"),
						resource.TestCheckResourceAttr("akamai_property_activation.test", "auto_acknowledge_rule_warnings", "true"),
						resource.TestCheckResourceAttr("akamai_property_activation.test", "activation_id", "atv_activation1"),
						resource.TestCheckResourceAttr("akamai_property_activation.test", "status", "ACTIVE"),
						resource.TestCheckResourceAttr("akamai_property_activation.test", "note", "property activation note for creating"),
						resource.TestCheckResourceAttr("akamai_property_activation.test", "timeouts.#", "0"),
					),
				},
			},
		},
		"property activation lifecycle - contact is correctly read in update when contact is exchanged": {
			init: func(m *papi.Mock) {
				// first step
				// create
				expectGetRuleTree(m, "prp_test", 1, ruleTreeResponseValid, nil).Once()
				expectGetActivations(m, "prp_test", papi.GetActivationsResponse{}, nil).Once()
				expectCreateActivation(m, "prp_test", papi.ActivationTypeActivate, 1, "STAGING",
					[]string{"user@example.com"}, "property activation note for creating", "atv_activation1", true, nil).Once()
				expectGetActivation(m, "prp_test", "atv_activation1", 1, "STAGING", papi.ActivationStatusActive, papi.ActivationTypeActivate, "property activation note for creating", []string{"user@example.com"}, nil).Once()
				// read
				expectGetActivations(m, "prp_test", generateActivationResponseMock("atv_activation1", "property activation note for creating", 1, papi.ActivationTypeActivate, "2020-10-28T15:04:05Z", []string{"user@example.com"}), nil).Once()

				// second step
				// read
				expectGetActivations(m, "prp_test", generateActivationResponseMock("atv_activation1", "property activation note for creating", 1, papi.ActivationTypeActivate, "2020-10-28T15:04:05Z", []string{"user@example.com"}), nil).Once()
				// update
				expectGetRuleTree(m, "prp_test", 2, ruleTreeResponseValid, nil).Once()
				expectGetActivations(m, "prp_test", generateActivationResponseMock("atv_activation1", "property activation note for creating", 1, papi.ActivationTypeActivate, "2020-10-28T15:04:05Z", []string{"user@example.com"}), nil).Once()
				expectGetPropertyVersion(m, "prp_test", "", "", 2, papi.VersionStatusInactive, "").Once()
				expectCreateActivation(m, "prp_test", papi.ActivationTypeActivate, 2, "STAGING",
					[]string{"user2@example.com"}, "property activation note for updating", "atv_update", true, nil).Once()
				expectGetActivation(m, "prp_test", "atv_update", 2, "STAGING", papi.ActivationStatusActive, papi.ActivationTypeActivate, "property activation note for updating", []string{"user2@example.com"}, nil).Once()
				// read
				getActivationsAfterUpdateResponse := papi.GetActivationsResponse{
					Activations: papi.ActivationsItems{Items: []*papi.Activation{mockDeactivationForLifecycle, generateActivationItemMock("atv_update", "property activation note for updating", 2, papi.ActivationTypeActivate, "2020-10-28T15:04:05Z", []string{"user2@example.com"})}},
				}
				expectGetActivations(m, "prp_test", getActivationsAfterUpdateResponse, nil).Once()
				// delete
				expectGetActivations(m, "prp_test", getActivationsAfterUpdateResponse, nil).Once()
				expectCreateActivation(m, "prp_test", papi.ActivationTypeDeactivate, 2, "STAGING",
					[]string{"user2@example.com"}, "property activation note for updating", "atv_update", true, nil).Once()
				expectGetActivation(m, "prp_test", "atv_update", 2, "STAGING", papi.ActivationStatusActive, papi.ActivationTypeDeactivate, "property activation note for updating", []string{"user2@example.com"}, nil).Once()
			},
			steps: []resource.TestStep{
				{
					Config: testutils.LoadFixtureString(t, "./testdata/TestPropertyActivation/ok/resource_property_activation.tf"),
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
						resource.TestCheckResourceAttr("akamai_property_activation.test", "rule_errors.#", "0"),
						resource.TestCheckResourceAttr("akamai_property_activation.test", "auto_acknowledge_rule_warnings", "true"),
						resource.TestCheckResourceAttr("akamai_property_activation.test", "activation_id", "atv_activation1"),
						resource.TestCheckResourceAttr("akamai_property_activation.test", "status", "ACTIVE"),
						resource.TestCheckResourceAttr("akamai_property_activation.test", "note", "property activation note for creating"),
						resource.TestCheckResourceAttr("akamai_property_activation.test", "timeouts.#", "0"),
					),
				},
				{
					Config: testutils.LoadFixtureString(t, "./testdata/TestPropertyActivation/ok/resource_property_activation_exchange_contact.tf"),
					Check: resource.ComposeAggregateTestCheckFunc(
						resource.TestCheckResourceAttr("akamai_property_activation.test", "id", "prp_test:STAGING"),
						resource.TestCheckResourceAttr("akamai_property_activation.test", "contact.#", "1"),
						resource.TestCheckResourceAttrSet("akamai_property_activation.test", "contact.0"),
						resource.TestCheckResourceAttr("akamai_property_activation.test", "contact.0", "user2@example.com"),
						resource.TestCheckResourceAttr("akamai_property_activation.test", "property_id", "prp_test"),
						resource.TestCheckResourceAttr("akamai_property_activation.test", "network", "STAGING"),
						resource.TestCheckResourceAttr("akamai_property_activation.test", "version", "2"),
						resource.TestCheckResourceAttr("akamai_property_activation.test", "activation_id", "atv_update"),
						resource.TestCheckResourceAttr("akamai_property_activation.test", "status", "ACTIVE"),
						resource.TestCheckResourceAttr("akamai_property_activation.test", "note", "property activation note for updating"),
						resource.TestCheckResourceAttr("akamai_property_activation.test", "timeouts.#", "0"),
					),
				},
			},
		},
		"property activation lifecycle - contact is correctly read in update when many changes within contacts": {
			init: func(m *papi.Mock) {
				// first step
				// create
				expectGetRuleTree(m, "prp_test", 1, ruleTreeResponseValid, nil).Once()
				expectGetActivations(m, "prp_test", papi.GetActivationsResponse{}, nil).Once()
				expectCreateActivation(m, "prp_test", papi.ActivationTypeActivate, 1, "STAGING",
					[]string{"user3@example.com", "user@example.com", "user2@example.com"}, "property activation note for creating", "atv_activation1", true, nil).Once()
				expectGetActivation(m, "prp_test", "atv_activation1", 1, "STAGING", papi.ActivationStatusActive, papi.ActivationTypeActivate, "property activation note for creating", []string{"user3@example.com", "user@example.com", "user2@example.com"}, nil).Once()
				// read
				expectGetActivations(m, "prp_test", generateActivationResponseMock("atv_activation1", "property activation note for creating", 1, papi.ActivationTypeActivate, "2020-10-28T15:04:05Z", []string{"user3@example.com", "user@example.com", "user2@example.com"}), nil).Once()

				// second step
				// read
				expectGetActivations(m, "prp_test", generateActivationResponseMock("atv_activation1", "property activation note for creating", 1, papi.ActivationTypeActivate, "2020-10-28T15:04:05Z", []string{"user3@example.com", "user@example.com", "user2@example.com"}), nil).Once()
				// update
				expectGetRuleTree(m, "prp_test", 2, ruleTreeResponseValid, nil).Once()
				expectGetActivations(m, "prp_test", generateActivationResponseMock("atv_activation1", "property activation note for creating", 1, papi.ActivationTypeActivate, "2020-10-28T15:04:05Z", []string{"user3@example.com", "user@example.com", "user2@example.com"}), nil).Once()
				expectGetPropertyVersion(m, "prp_test", "", "", 2, papi.VersionStatusInactive, "").Once()
				expectCreateActivation(m, "prp_test", papi.ActivationTypeActivate, 2, "STAGING",
					[]string{"user3@example.com", "user5@example.com", "user2@example.com", "user4@example.com"}, "property activation note for updating", "atv_update", true, nil).Once()
				expectGetActivation(m, "prp_test", "atv_update", 2, "STAGING", papi.ActivationStatusActive, papi.ActivationTypeActivate, "property activation note for updating", []string{"user3@example.com", "user5@example.com", "user2@example.com", "user4@example.com"}, nil).Once()
				// read
				getActivationsAfterUpdateResponse := papi.GetActivationsResponse{
					Activations: papi.ActivationsItems{Items: []*papi.Activation{mockDeactivationForLifecycle, generateActivationItemMock("atv_update", "property activation note for updating", 2, papi.ActivationTypeActivate, "2020-10-28T15:04:05Z", []string{"user3@example.com", "user5@example.com", "user2@example.com", "user4@example.com"})}},
				}
				expectGetActivations(m, "prp_test", getActivationsAfterUpdateResponse, nil).Once()
				// delete
				expectGetActivations(m, "prp_test", getActivationsAfterUpdateResponse, nil).Once()
				expectCreateActivation(m, "prp_test", papi.ActivationTypeDeactivate, 2, "STAGING",
					[]string{"user3@example.com", "user5@example.com", "user2@example.com", "user4@example.com"}, "property activation note for updating", "atv_update", true, nil).Once()
				expectGetActivation(m, "prp_test", "atv_update", 2, "STAGING", papi.ActivationStatusActive, papi.ActivationTypeDeactivate, "property activation note for updating", []string{"user3@example.com", "user5@example.com", "user2@example.com", "user4@example.com"}, nil).Once()
			},
			steps: []resource.TestStep{
				{
					Config: testutils.LoadFixtureString(t, "./testdata/TestPropertyActivation/ok/resource_property_activation_complicated_contact.tf"),
					Check: resource.ComposeAggregateTestCheckFunc(
						resource.TestCheckResourceAttr("akamai_property_activation.test", "id", "prp_test:STAGING"),
						resource.TestCheckResourceAttr("akamai_property_activation.test", "property_id", "prp_test"),
						resource.TestCheckResourceAttr("akamai_property_activation.test", "contact.#", "3"),
						resource.TestCheckResourceAttr("akamai_property_activation.test", "contact.0", "user2@example.com"),
						resource.TestCheckResourceAttr("akamai_property_activation.test", "contact.1", "user3@example.com"),
						resource.TestCheckResourceAttr("akamai_property_activation.test", "contact.2", "user@example.com"),
						resource.TestCheckResourceAttr("akamai_property_activation.test", "network", "STAGING"),
						resource.TestCheckResourceAttr("akamai_property_activation.test", "version", "1"),
						resource.TestCheckResourceAttr("akamai_property_activation.test", "warnings", ""),
						resource.TestCheckResourceAttr("akamai_property_activation.test", "errors", ""),
						resource.TestCheckResourceAttr("akamai_property_activation.test", "rule_errors.#", "0"),
						resource.TestCheckResourceAttr("akamai_property_activation.test", "auto_acknowledge_rule_warnings", "true"),
						resource.TestCheckResourceAttr("akamai_property_activation.test", "activation_id", "atv_activation1"),
						resource.TestCheckResourceAttr("akamai_property_activation.test", "status", "ACTIVE"),
						resource.TestCheckResourceAttr("akamai_property_activation.test", "note", "property activation note for creating"),
						resource.TestCheckResourceAttr("akamai_property_activation.test", "timeouts.#", "0"),
					),
				},
				{
					Config: testutils.LoadFixtureString(t, "./testdata/TestPropertyActivation/ok/resource_property_activation_complicated_contact_update.tf"),
					Check: resource.ComposeAggregateTestCheckFunc(
						resource.TestCheckResourceAttr("akamai_property_activation.test", "id", "prp_test:STAGING"),
						resource.TestCheckResourceAttr("akamai_property_activation.test", "contact.#", "4"),
						resource.TestCheckResourceAttr("akamai_property_activation.test", "contact.0", "user2@example.com"),
						resource.TestCheckResourceAttr("akamai_property_activation.test", "contact.1", "user3@example.com"),
						resource.TestCheckResourceAttr("akamai_property_activation.test", "contact.2", "user4@example.com"),
						resource.TestCheckResourceAttr("akamai_property_activation.test", "contact.3", "user5@example.com"),
						resource.TestCheckResourceAttr("akamai_property_activation.test", "property_id", "prp_test"),
						resource.TestCheckResourceAttr("akamai_property_activation.test", "network", "STAGING"),
						resource.TestCheckResourceAttr("akamai_property_activation.test", "version", "2"),
						resource.TestCheckResourceAttr("akamai_property_activation.test", "activation_id", "atv_update"),
						resource.TestCheckResourceAttr("akamai_property_activation.test", "status", "ACTIVE"),
						resource.TestCheckResourceAttr("akamai_property_activation.test", "note", "property activation note for updating"),
						resource.TestCheckResourceAttr("akamai_property_activation.test", "timeouts.#", "0"),
					),
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
				expectGetActivations(m, "prp_test", generateActivationResponseMock("atv_activation1", "", 1, papi.ActivationTypeActivate, "2020-10-28T15:04:05Z", []string{"user@example.com"}), nil).Once()

				// second step
				// read
				expectGetActivations(m, "prp_test", generateActivationResponseMock("atv_activation1", "", 1, papi.ActivationTypeActivate, "2020-10-28T15:04:05Z", []string{"user@example.com"}), nil).Once()
				// update
				expectGetRuleTree(m, "prp_test", 2, ruleTreeResponseValid, nil).Once()
				expectGetActivations(m, "prp_test", generateActivationResponseMock("atv_activation1", "", 1, papi.ActivationTypeActivate, "2020-10-28T15:04:05Z", []string{"user@example.com"}), nil).Once()
				expectGetPropertyVersion(m, "prp_test", "", "", 2, papi.VersionStatusInactive, "").Once()
				// error on update
				m.On("CreateActivation", testutils.MockContext, papi.CreateActivationRequest{
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
				expectGetActivations(m, "prp_test", generateActivationResponseMock("atv_activation1", "", 1, papi.ActivationTypeActivate, "2020-10-28T15:04:05Z", []string{"user@example.com"}), nil).Once()
				expectCreateActivation(m, "prp_test", papi.ActivationTypeDeactivate, 2, "STAGING", []string{"user@example.com"}, "", "atv_update", true, nil).Once()
				expectGetActivation(m, "prp_test", "atv_update", 2, "STAGING", papi.ActivationStatusActive, papi.ActivationTypeDeactivate, "", []string{"user@example.com"}, nil).Once()

			},
			steps: []resource.TestStep{
				{
					Config: testutils.LoadFixtureString(t, "./testdata/TestPropertyActivation/500_on_activation/resource_property_activation.tf"),
					Check:  resource.TestCheckResourceAttr("akamai_property_activation.test", "version", "1"),
				},
				{
					Config:      testutils.LoadFixtureString(t, "./testdata/TestPropertyActivation/500_on_activation/resource_property_activation_update.tf"),
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
				expectGetActivations(m, "prp_test", generateActivationResponseMock("atv_activation1", "property activation note for importing", 1, papi.ActivationTypeActivate, "2020-10-28T15:04:05Z", []string{"user@example.com"}), nil).Once()
				// 2nd read for import
				expectGetActivations(m, "prp_test", generateActivationResponseMock("atv_activation1", "property activation note for importing", 1, papi.ActivationTypeActivate, "2020-10-28T15:04:05Z", []string{"user@example.com"}), nil).Once()
				// delete
				expectGetActivations(m, "prp_test", generateActivationResponseMock("atv_activation1", "property activation note for importing", 1, papi.ActivationTypeActivate, "2020-10-28T15:04:05Z", []string{"user@example.com"}), nil).Once()
				expectCreateActivation(m, "prp_test", papi.ActivationTypeDeactivate, 1, "STAGING",
					[]string{"user@example.com"}, "property activation note for importing", "atv_activation1", false, nil).Once()
			},
			steps: []resource.TestStep{
				{
					Config: testutils.LoadFixtureString(t, "./testdata/TestPropertyActivation/import/resource_property_activation_creation_for_import.tf"),
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
						resource.TestCheckResourceAttr("akamai_property_activation.test", "rule_errors.#", "0"),
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
		"property activation import - no differences": {
			init: func(m *papi.Mock) {
				// Create activation
				expectGetRuleTree(m, "prp_test", 1, ruleTreeResponseValid, nil).Once()
				expectGetActivations(m, "prp_test", papi.GetActivationsResponse{}, nil).Once()
				expectCreateActivation(m, "prp_test", papi.ActivationTypeActivate, 1, "STAGING",
					[]string{"user@example.com"}, "property activation note for importing", "atv_activation1", false, nil).Once()
				expectGetActivation(m, "prp_test", "atv_activation1", 1, "STAGING", papi.ActivationStatusActive, papi.ActivationTypeActivate, "property activation note for importing", []string{"user@example.com"}, nil).Twice()
				// Read activation details
				expectGetActivations(m, "prp_test", generateActivationResponseMock("atv_activation1", "property activation note for importing", 1, papi.ActivationTypeActivate, "2020-10-28T15:04:05Z", []string{"user@example.com"}), nil).Once()
				// 2nd read for import
				expectGetActivations(m, "prp_test", generateActivationResponseMock("atv_activation1", "property activation note for importing", 1, papi.ActivationTypeActivate, "2020-10-28T15:04:05Z", []string{"user@example.com"}), nil).Twice()
				// Prepare delete activation
				expectGetActivations(m, "prp_test", generateActivationResponseMock("atv_activation1", "property activation note for importing", 1, papi.ActivationTypeActivate, "2020-10-28T15:04:05Z", []string{"user@example.com"}), nil).Once()
				expectCreateActivation(m, "prp_test", papi.ActivationTypeDeactivate, 1, "STAGING",
					[]string{"user@example.com"}, "property activation note for importing", "atv_activation1", false, nil).Once()
			},
			steps: []resource.TestStep{
				{
					// Initial Configuration and Verification
					Config: testutils.LoadFixtureString(t, "./testdata/TestPropertyActivation/import/resource_property_activation_creation_for_import.tf"),
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
						resource.TestCheckResourceAttr("akamai_property_activation.test", "rule_errors.#", "0"),
						resource.TestCheckResourceAttr("akamai_property_activation.test", "errors", ""),
						resource.TestCheckResourceAttr("akamai_property_activation.test", "activation_id", "atv_activation1"),
						resource.TestCheckResourceAttr("akamai_property_activation.test", "status", "ACTIVE"),
						resource.TestCheckResourceAttr("akamai_property_activation.test", "note", "property activation note for importing"),
					),
				},
				{
					// Import Step
					ImportState:       true,
					ImportStateId:     "prp_test:STAGING",
					ResourceName:      "akamai_property_activation.test",
					ImportStateVerify: true,
				},
				{
					// Post-Import Diff Check (No-op Plan Verification)
					Config:   testutils.LoadFixtureString(t, "./testdata/TestPropertyActivation/import/resource_property_activation_creation_for_import.tf"),
					PlanOnly: true, // Ensures no changes are expected after import
					Check: resource.ComposeAggregateTestCheckFunc(
						resource.TestCheckResourceAttr("akamai_property_activation.test", "id", "prp_test:STAGING"),
						resource.TestCheckResourceAttr("akamai_property_activation.test", "status", "ACTIVE"),
						resource.TestCheckResourceAttr("akamai_property_activation.test", "activation_id", "atv_activation1"),
						resource.TestCheckResourceAttr("akamai_property_activation.test", "note", "property activation note for importing"),
						resource.TestCheckResourceAttr("akamai_property_activation.test", "rule_errors.#", "0"),
					),
				},
			},
		},
		"property activation update - resolve errors": {
			init: func(m *papi.Mock) {
				// Initial create activation call with rule validation
				expectGetRuleTree(m, "prp_test", 1, ruleTreeResponseValid, nil).Once()
				expectGetActivations(m, "prp_test", papi.GetActivationsResponse{}, nil).Once()
				expectCreateActivation(m, "prp_test", papi.ActivationTypeActivate, 1, "STAGING",
					[]string{"user@example.com"}, "property activation note for creating", "atv_activation1", true, nil).Once()
				expectGetActivation(m, "prp_test", "atv_activation1", 1, "STAGING", papi.ActivationStatusActive, papi.ActivationTypeActivate, "property activation note for creating", []string{"user@example.com"}, nil).Once()

				// Mock for initial read after create
				expectGetActivations(m, "prp_test", generateActivationResponseMock("atv_activation1", "property activation note for creating", 1, papi.ActivationTypeActivate, "2020-10-28T15:04:05Z", []string{"user@example.com"}), nil).Twice()

				// Handle activation with errors due to rule tree errors
				expectGetRuleTree(m, "prp_test", 2, ruleTreeResponseInvalid, nil).Once()

				expectGetActivations(m, "prp_test", generateActivationResponseMock("atv_activation1", "property activation note for creating", 1, papi.ActivationTypeActivate, "2020-10-28T15:04:05Z", []string{"user@example.com"}), nil).Once()

				// Retry after error resolution - Update
				expectGetRuleTree(m, "prp_test", 2, ruleTreeResponseValid, nil).Once()
				expectGetActivations(m, "prp_test", generateActivationResponseMock("atv_activation1", "property activation note for creating", 1, papi.ActivationTypeActivate, "2020-10-28T15:04:05Z", []string{"user@example.com"}), nil).Once()
				expectGetPropertyVersion(m, "prp_test", "", "", 2, papi.VersionStatusInactive, "").Once()
				expectCreateActivation(m, "prp_test", papi.ActivationTypeActivate, 2, "STAGING",
					[]string{"user@example.com", "user2@example.com"}, "property activation note for updating", "atv_update", true, nil).Once()
				expectGetActivation(m, "prp_test", "atv_update", 2, "STAGING", papi.ActivationStatusActive, papi.ActivationTypeActivate, "property activation note for updating", []string{"user@example.com", "user2@example.com"}, nil).Once()

				// read
				expectGetActivations(m, "prp_test", activationsResponseWithNoteNoteSuppressed, nil).Once()
				// delete
				expectGetActivations(m, "prp_test", activationsResponseWithNoteNoteSuppressed, nil).Once()
				expectCreateActivation(m, "prp_test", papi.ActivationTypeDeactivate, 2, "STAGING",
					[]string{"user@example.com", "user2@example.com"}, "property activation note for updating", "atv_update", true, nil).Once()
				expectGetActivation(m, "prp_test", "atv_update", 2, "STAGING", papi.ActivationStatusActive, papi.ActivationTypeDeactivate, "property activation note for updating", []string{"user@example.com", "user2@example.com"}, nil).Once()
			},

			steps: []resource.TestStep{
				{
					Config: testutils.LoadFixtureString(t, "./testdata/TestPropertyActivation/ok/resource_property_activation.tf"),
					Check: resource.ComposeAggregateTestCheckFunc(
						resource.TestCheckResourceAttr("akamai_property_activation.test", "id", "prp_test:STAGING"),
						resource.TestCheckResourceAttr("akamai_property_activation.test", "property_id", "prp_test"),
						resource.TestCheckResourceAttr("akamai_property_activation.test", "contact.#", "1"),
						resource.TestCheckResourceAttr("akamai_property_activation.test", "contact.0", "user@example.com"),
						resource.TestCheckResourceAttr("akamai_property_activation.test", "network", "STAGING"),
						resource.TestCheckResourceAttr("akamai_property_activation.test", "version", "1"),
						resource.TestCheckResourceAttr("akamai_property_activation.test", "note", "property activation note for creating"),
						resource.TestCheckResourceAttr("akamai_property_activation.test", "status", "ACTIVE"),
						resource.TestCheckResourceAttr("akamai_property_activation.test", "rule_errors.#", "0"),
						resource.TestCheckResourceAttrSet("akamai_property_activation.test", "activation_id"),
					),
				},
				{
					// Expect error due to rule validation issue
					Config:      testutils.LoadFixtureString(t, "./testdata/TestPropertyActivation/ok/resource_property_activation_update.tf"),
					ExpectError: regexp.MustCompile("activation cannot continue due to rule errors"),
				},
				{
					Config: testutils.LoadFixtureString(t, "./testdata/TestPropertyActivation/diff_suppress/resource_property_activation_update.tf"),
					Check: resource.ComposeAggregateTestCheckFunc(
						resource.TestCheckResourceAttr("akamai_property_activation.test", "id", "prp_test:STAGING"),
						resource.TestCheckResourceAttr("akamai_property_activation.test", "contact.#", "2"),
						resource.TestCheckResourceAttrSet("akamai_property_activation.test", "contact.0"),
						resource.TestCheckResourceAttrSet("akamai_property_activation.test", "contact.1"),
						resource.TestCheckResourceAttr("akamai_property_activation.test", "property_id", "prp_test"),
						resource.TestCheckResourceAttr("akamai_property_activation.test", "network", "STAGING"),
						resource.TestCheckResourceAttr("akamai_property_activation.test", "version", "2"),
						resource.TestCheckResourceAttr("akamai_property_activation.test", "activation_id", "atv_update"),
						resource.TestCheckResourceAttr("akamai_property_activation.test", "status", "ACTIVE"),
						resource.TestCheckResourceAttr("akamai_property_activation.test", "rule_errors.#", "0"),
						resource.TestCheckResourceAttr("akamai_property_activation.test", "note", "property activation note for updating"),
					),
				},
			},
		},
		"fail update of auto_acknowledge_rule_warnings updated if property version is active": {
			init: func(m *papi.Mock) {
				// first step
				// create
				expectGetRuleTree(m, "prp_test", 1, ruleTreeResponseValid, nil).Once()
				expectGetActivations(m, "prp_test", papi.GetActivationsResponse{}, nil).Once()
				expectCreateActivation(m, "prp_test", papi.ActivationTypeActivate, 1, "STAGING",
					[]string{"user@example.com"}, "property activation note for creating", "atv_activation1", true, nil).Once()
				expectGetActivation(m, "prp_test", "atv_activation1", 1, "STAGING", papi.ActivationStatusActive, papi.ActivationTypeActivate, "property activation note for creating", []string{"user@example.com"}, nil).Once()
				// read
				expectGetActivations(m, "prp_test", generateActivationResponseMock("atv_activation1", "property activation note for creating", 1, papi.ActivationTypeActivate, "2020-10-28T15:04:05Z", []string{"user@example.com"}), nil).Once()

				// second step
				expectGetRuleTree(m, "prp_test", 1, ruleTreeResponseValid, nil).Once()

				// property version is active
				expectGetPropertyVersion(m, "prp_test", "", "", 1, papi.VersionStatusActive, "").Once()

				// read
				expectGetActivations(m, "prp_test", generateActivationResponseMock("atv_activation1", "property activation note for creating", 1, papi.ActivationTypeActivate, "2020-10-28T15:04:05Z", []string{"user@example.com"}), nil).Once()

				// update
				expectGetActivations(m, "prp_test", generateActivationResponseMock("atv_activation1", "property activation note for creating", 1, papi.ActivationTypeActivate, "2020-10-28T15:04:05Z", []string{"user@example1.com"}), nil).Once()

				// delete
				expectGetActivations(m, "prp_test", generateActivationResponseMock("atv_activation1", "property activation note for creating", 1, papi.ActivationTypeActivate, "2020-10-28T15:04:05Z", []string{"user@example.com"}), nil).Once()
				expectCreateActivation(m, "prp_test", papi.ActivationTypeDeactivate, 1, "STAGING",
					[]string{"user@example.com"}, "property activation note for creating", "atv_activation1", false, nil).Once()

				expectGetActivation(m, "prp_test", "atv_activation1", 1, "STAGING", papi.ActivationStatusActive, papi.ActivationTypeDeactivate, "property activation note for creating", []string{"user@example.com"}, nil).Once()
			},
			steps: []resource.TestStep{
				{
					Config: testutils.LoadFixtureString(t, "./testdata/TestPropertyActivation/ok/resource_property_activation.tf"),
					Check:  baseChecker.Build(),
				},
				{
					Config:      testutils.LoadFixtureString(t, "./testdata/TestPropertyActivation/ok/resource_property_activation_with_changes.tf"),
					ExpectError: regexp.MustCompile(`Error: Cannot update 'auto_acknowledge_rule_warnings' field\(s\) while property version is ACTIVE\.\s*Deactivate the current version to update, or create a new property version activation\.`),
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
					ProtoV6ProviderFactories: testutils.NewProtoV6ProviderFactory(NewSubprovider()),
					IsUnitTest:               true,
					Steps:                    test.steps,
				})
			})
			client.AssertExpectations(t)
		})
	}
}

func generateActivationResponseMock(activationID string, note string, version int, activationType papi.ActivationType, date string, emails []string) papi.GetActivationsResponse {
	return papi.GetActivationsResponse{
		Activations: papi.ActivationsItems{Items: []*papi.Activation{generateActivationItemMock(activationID, note, version, activationType, date, emails)}},
	}
}

func generateActivationItemMock(activationID string, note string, version int, activationType papi.ActivationType, date string, emails []string) *papi.Activation {
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
		NotifyEmails:    emails,
		Note:            note,
	}
}

var (
	mockActivationsListForDeactivation      = []*papi.Activation{mockDeactivation, mockActivationCreation}
	mockActivationsListForLifecycle         = []*papi.Activation{mockDeactivationForLifecycle, mockActivationUpdate}
	mockActivationsListForNoteNotSuppressed = []*papi.Activation{mockDeactivationForLifecycle, mockActivationNoteNotSuppressed}
	mockDeactivation                        = generateActivationItemMock("atv_activation1", "property activation note for creating", 1, papi.ActivationTypeDeactivate, "2020-10-28T15:04:05Z", []string{"user@example.com"})
	mockDeactivationForLifecycle            = generateActivationItemMock("atv_activation1", "property activation note for creating", 1, papi.ActivationTypeDeactivate, "2020-10-28T14:04:05Z", []string{"user@example.com"})
	mockActivationCreation                  = generateActivationItemMock("atv_activation1", "property activation note for creating", 1, papi.ActivationTypeActivate, "2020-10-28T15:04:05Z", []string{"user@example.com"})
	mockActivationUpdate                    = generateActivationItemMock("atv_update", "property activation note for updating", 2, papi.ActivationTypeActivate, "2020-10-28T15:04:05Z", []string{"user@example.com"})
	mockActivationNoteNotSuppressed         = generateActivationItemMock("atv_update", "property activation note for updating", 2, papi.ActivationTypeActivate, "2020-10-28T15:04:05Z", []string{"user@example.com", "user2@example.com"})
	ruleTreeResponseValid                   = papi.GetRuleTreeResponse{
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
				testutils.MockContext,
				papi.GetRuleTreeRequest{PropertyID: propertyID, PropertyVersion: version, ValidateRules: true},
			).Return(nil, err)
		}
		return m.On(
			"GetRuleTree",
			testutils.MockContext,
			papi.GetRuleTreeRequest{PropertyID: propertyID, PropertyVersion: version, ValidateRules: true},
		).Return(&response, nil)
	}

	activationsResponseDeactivated = papi.GetActivationsResponse{
		Activations: papi.ActivationsItems{Items: mockActivationsListForDeactivation},
	}
	activationsResponseSecondVersionIsActive = papi.GetActivationsResponse{
		Activations: papi.ActivationsItems{Items: mockActivationsListForLifecycle},
	}
	activationsResponseWithNoteNoteSuppressed = papi.GetActivationsResponse{
		Activations: papi.ActivationsItems{Items: mockActivationsListForNoteNotSuppressed},
	}
	expectGetActivations = func(m *papi.Mock, propertyID string, response papi.GetActivationsResponse, err error) *mock.Call {
		if err != nil {
			return m.On(
				"GetActivations",
				testutils.MockContext,
				papi.GetActivationsRequest{PropertyID: propertyID},
			).Return(nil, err)
		}
		return m.On(
			"GetActivations",
			testutils.MockContext,
			papi.GetActivationsRequest{PropertyID: propertyID},
		).Return(&response, nil)
	}

	expectCreateActivation = func(m *papi.Mock, propertyID string, activationType papi.ActivationType, version int,
		network papi.ActivationNetwork, notify []string, note string, activationID string, acknowledgeAllWarnings bool, err error) *mock.Call {
		if err != nil {
			return m.On(
				"CreateActivation",
				testutils.MockContext,
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
			testutils.MockContext,
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

	expectCreateActivation500Err = func(m *papi.Mock, propertyID string, activationType papi.ActivationType, version int,
		network papi.ActivationNetwork, notify []string, note string, activationID string, acknowledgeAllWarnings bool, err error) {
		m.On(
			"CreateActivation",
			testutils.MockContext,
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
		).Return(nil, err).Once()

		m.On("GetActivations", testutils.MockContext, papi.GetActivationsRequest{
			PropertyID: propertyID,
		}).Return(&papi.GetActivationsResponse{
			Activations: papi.ActivationsItems{
				Items: []*papi.Activation{
					{
						ActivationID:           activationID,
						ActivationType:         activationType,
						AcknowledgeAllWarnings: acknowledgeAllWarnings,
						PropertyVersion:        version,
						Network:                network,
						NotifyEmails:           notify,
						Note:                   note,
						Status:                 papi.ActivationStatusActive,
						UpdateDate:             "2023-11-28T08:43:33Z",
					},
				},
			},
		}, nil).Once()

	}

	expectGetActivation = func(m *papi.Mock, propertyID string, activationID string, version int,
		network papi.ActivationNetwork, status papi.ActivationStatus, actType papi.ActivationType, note string, contact []string, err error) *mock.Call {
		if err != nil {
			return m.On(
				"GetActivation",
				testutils.MockContext,
				papi.GetActivationRequest{
					PropertyID:   propertyID,
					ActivationID: activationID,
				},
			).Return(nil, err)
		}
		return m.On(
			"GetActivation",
			testutils.MockContext,
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

	// Sets up an expected call to papi.GetPropertyVersion()
	expectGetPropertyVersion = func(client *papi.Mock, PropertyID, GroupID, ContractID string, Version int, StagStatus, ProdStatus papi.VersionStatus) *mock.Call {
		req := papi.GetPropertyVersionRequest{
			PropertyID:      PropertyID,
			GroupID:         GroupID,
			ContractID:      ContractID,
			PropertyVersion: Version,
		}

		res := papi.GetPropertyVersionsResponse{
			PropertyID: PropertyID,
			GroupID:    GroupID,
			ContractID: ContractID,
			Version: papi.PropertyVersionGetItem{
				StagingStatus:    StagStatus,
				ProductionStatus: ProdStatus,
			},
		}
		return client.On("GetPropertyVersion", testutils.MockContext, req).Return(&res, nil)
	}
)
