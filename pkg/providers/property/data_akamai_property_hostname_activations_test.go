package property

import (
	"fmt"
	"regexp"
	"testing"

	"github.com/akamai/AkamaiOPEN-edgegrid-golang/v12/pkg/hapi"
	"github.com/akamai/AkamaiOPEN-edgegrid-golang/v12/pkg/papi"
	internal "github.com/akamai/terraform-provider-akamai/v8/internal/test"
	"github.com/akamai/terraform-provider-akamai/v8/pkg/common/test"
	"github.com/akamai/terraform-provider-akamai/v8/pkg/common/testutils"
	"github.com/hashicorp/terraform-plugin-testing/helper/resource"
)

func TestDataPropertyHostnameActivations(t *testing.T) {
	t.Parallel()
	baseChecker := test.NewStateChecker("data.akamai_property_hostname_activations.activation").
		CheckEqual("property_id", "prp_1").
		CheckEqual("group_id", "grp_1").
		CheckEqual("contract_id", "ctr_1").
		CheckEqual("account_id", "act_1").
		CheckEqual("property_name", "my_property_1").
		CheckEqual("hostname_activations.0.activation_type", "ACTIVATE").
		CheckEqual("hostname_activations.0.hostname_activation_id", "44").
		CheckEqual("hostname_activations.0.network", "STAGING").
		CheckEqual("hostname_activations.0.status", "ACTIVE").
		CheckEqual("hostname_activations.0.submit_date", "2025-01-22T19:36:29Z").
		CheckEqual("hostname_activations.0.update_date", "2025-01-22T19:37:48Z").
		CheckEqual("hostname_activations.0.note", "   ").
		CheckEqual("hostname_activations.0.notify_emails.#", "1").
		CheckEqual("hostname_activations.0.notify_emails.0", "nomail@akamai.com")

	activations := papi.HostnameActivationsList{
		Items: []papi.HostnameActivationListItem{
			{
				ActivationType:       "ACTIVATE",
				HostnameActivationID: "44",
				PropertyName:         "my_property_1",
				PropertyID:           "prp_1",
				Network:              "STAGING",
				Status:               "ACTIVE",
				SubmitDate:           internal.NewTimeFromString(t, "2025-01-22T19:36:29Z"),
				UpdateDate:           internal.NewTimeFromString(t, "2025-01-22T19:37:48Z"),
				Note:                 "   ",
				NotifyEmails:         []string{"nomail@akamai.com"},
			},
			{
				ActivationType:       "ACTIVATE",
				HostnameActivationID: "55",
				PropertyName:         "my_property_1",
				PropertyID:           "prp_1",
				Network:              "PRODUCTION",
				Status:               "ACTIVE",
				SubmitDate:           internal.NewTimeFromString(t, "2025-01-22T19:36:29Z"),
				UpdateDate:           internal.NewTimeFromString(t, "2025-01-22T19:37:48Z"),
				Note:                 "   ",
				NotifyEmails:         []string{"nomail@akamai.com"},
			},
			{
				ActivationType:       "ACTIVATE",
				HostnameActivationID: "66",
				PropertyName:         "my_property_1",
				PropertyID:           "prp_1",
				Network:              "PRODUCTION",
				Status:               "ABORTED",
				SubmitDate:           internal.NewTimeFromString(t, "2025-01-22T19:36:29Z"),
				UpdateDate:           internal.NewTimeFromString(t, "2025-01-22T19:37:48Z"),
				Note:                 "   ",
				NotifyEmails:         []string{"nomail@akamai.com"},
			},
		},
		TotalItems:       3,
		CurrentItemCount: 3,
	}

	activations999 := make([]papi.HostnameActivationListItem, 999)
	for i := 0; i < 999; i++ {
		activations999[i] = papi.HostnameActivationListItem{
			ActivationType:       "ACTIVATE",
			HostnameActivationID: fmt.Sprintf("%d", i),
			PropertyName:         "my_property_1",
			PropertyID:           "prp_1",
			Network:              "PRODUCTION",
			Status:               "ACTIVE",
			SubmitDate:           internal.NewTimeFromString(t, "2025-01-22T19:36:29Z"),
			UpdateDate:           internal.NewTimeFromString(t, "2025-01-22T19:37:48Z"),
			Note:                 "   ",
			NotifyEmails:         []string{"nomail@akamai.com"},
		}
	}
	activations101 := make([]papi.HostnameActivationListItem, 101)
	for i := 0; i < 101; i++ {
		activations101[i] = papi.HostnameActivationListItem{
			ActivationType:       "ACTIVATE",
			HostnameActivationID: fmt.Sprintf("%d", i+999),
			PropertyName:         "my_property_1",
			PropertyID:           "prp_1",
			Network:              "PRODUCTION",
			Status:               "ACTIVE",
			SubmitDate:           internal.NewTimeFromString(t, "2025-01-22T19:36:29Z"),
			UpdateDate:           internal.NewTimeFromString(t, "2025-01-22T19:37:48Z"),
			Note:                 "   ",
			NotifyEmails:         []string{"nomail@akamai.com"},
		}
	}

	tests := map[string]struct {
		init  func(*papi.Mock)
		steps []resource.TestStep
	}{
		"happy path": {
			init: func(m *papi.Mock) {
				m.On("ListPropertyHostnameActivations", testutils.MockContext, papi.ListPropertyHostnameActivationsRequest{
					ContractID: "1",
					GroupID:    "1",
					PropertyID: "1",
					Offset:     0,
					Limit:      999,
				}).Return(&papi.ListPropertyHostnameActivationsResponse{
					AccountID:           "act_1",
					ContractID:          "ctr_1",
					GroupID:             "grp_1",
					HostnameActivations: activations,
				}, nil).Times(3)
			},
			steps: []resource.TestStep{
				{
					Config: testutils.LoadFixtureString(t, "testdata/TestDataPropertyHostnameActivations/valid.tf"),
					Check: baseChecker.
						CheckEqual("property_id", "1").
						CheckEqual("hostname_activations.#", "3").
						CheckEqual("hostname_activations.1.activation_type", "ACTIVATE").
						CheckEqual("hostname_activations.1.hostname_activation_id", "55").
						CheckEqual("hostname_activations.1.network", "PRODUCTION").
						CheckEqual("hostname_activations.1.status", "ACTIVE").
						CheckEqual("hostname_activations.1.submit_date", "2025-01-22T19:36:29Z").
						CheckEqual("hostname_activations.1.update_date", "2025-01-22T19:37:48Z").
						CheckEqual("hostname_activations.1.note", "   ").
						CheckEqual("hostname_activations.1.notify_emails.#", "1").
						CheckEqual("hostname_activations.1.notify_emails.0", "nomail@akamai.com").
						CheckEqual("hostname_activations.2.activation_type", "ACTIVATE").
						CheckEqual("hostname_activations.2.hostname_activation_id", "66").
						CheckEqual("hostname_activations.2.network", "PRODUCTION").
						CheckEqual("hostname_activations.2.status", "ABORTED").
						CheckEqual("hostname_activations.2.submit_date", "2025-01-22T19:36:29Z").
						CheckEqual("hostname_activations.2.update_date", "2025-01-22T19:37:48Z").
						CheckEqual("hostname_activations.2.note", "   ").
						CheckEqual("hostname_activations.2.notify_emails.#", "1").
						CheckEqual("hostname_activations.2.notify_emails.0", "nomail@akamai.com").
						Build(),
				},
			},
		},
		"happy path - no contract and group": {
			init: func(m *papi.Mock) {
				m.On("ListPropertyHostnameActivations", testutils.MockContext, papi.ListPropertyHostnameActivationsRequest{
					ContractID: "",
					GroupID:    "",
					PropertyID: "prp_1",
					Offset:     0,
					Limit:      999,
				}).Return(&papi.ListPropertyHostnameActivationsResponse{
					AccountID:           "act_1",
					ContractID:          "ctr_1",
					GroupID:             "grp_1",
					HostnameActivations: activations,
				}, nil).Times(3)
			},
			steps: []resource.TestStep{
				{
					Config: testutils.LoadFixtureString(t, "testdata/TestDataPropertyHostnameActivations/valid_no_contract_and_group.tf"),
					Check: baseChecker.
						CheckEqual("hostname_activations.#", "3").
						CheckEqual("hostname_activations.1.activation_type", "ACTIVATE").
						CheckEqual("hostname_activations.1.hostname_activation_id", "55").
						CheckEqual("hostname_activations.1.network", "PRODUCTION").
						CheckEqual("hostname_activations.1.status", "ACTIVE").
						CheckEqual("hostname_activations.1.submit_date", "2025-01-22T19:36:29Z").
						CheckEqual("hostname_activations.1.update_date", "2025-01-22T19:37:48Z").
						CheckEqual("hostname_activations.1.note", "   ").
						CheckEqual("hostname_activations.1.notify_emails.#", "1").
						CheckEqual("hostname_activations.1.notify_emails.0", "nomail@akamai.com").
						CheckEqual("hostname_activations.2.activation_type", "ACTIVATE").
						CheckEqual("hostname_activations.2.hostname_activation_id", "66").
						CheckEqual("hostname_activations.2.network", "PRODUCTION").
						CheckEqual("hostname_activations.2.status", "ABORTED").
						CheckEqual("hostname_activations.2.submit_date", "2025-01-22T19:36:29Z").
						CheckEqual("hostname_activations.2.update_date", "2025-01-22T19:37:48Z").
						CheckEqual("hostname_activations.2.note", "   ").
						CheckEqual("hostname_activations.2.notify_emails.#", "1").
						CheckEqual("hostname_activations.2.notify_emails.0", "nomail@akamai.com").
						Build(),
				},
			},
		},
		"happy path - only production network": {
			init: func(m *papi.Mock) {
				m.On("ListPropertyHostnameActivations", testutils.MockContext, papi.ListPropertyHostnameActivationsRequest{
					ContractID: "ctr_1",
					GroupID:    "grp_1",
					PropertyID: "prp_1",
					Offset:     0,
					Limit:      999,
				}).Return(&papi.ListPropertyHostnameActivationsResponse{
					AccountID:           "act_1",
					ContractID:          "ctr_1",
					GroupID:             "grp_1",
					HostnameActivations: activations,
				}, nil).Times(3)
			},
			steps: []resource.TestStep{
				{
					Config: testutils.LoadFixtureString(t, "testdata/TestDataPropertyHostnameActivations/valid_production.tf"),
					Check: baseChecker.
						CheckEqual("hostname_activations.#", "2").
						CheckEqual("hostname_activations.0.hostname_activation_id", "55").
						CheckEqual("hostname_activations.0.network", "PRODUCTION").
						CheckEqual("hostname_activations.1.activation_type", "ACTIVATE").
						CheckEqual("hostname_activations.1.hostname_activation_id", "66").
						CheckEqual("hostname_activations.1.network", "PRODUCTION").
						CheckEqual("hostname_activations.1.status", "ABORTED").
						CheckEqual("hostname_activations.1.submit_date", "2025-01-22T19:36:29Z").
						CheckEqual("hostname_activations.1.update_date", "2025-01-22T19:37:48Z").
						CheckEqual("hostname_activations.1.note", "   ").
						CheckEqual("hostname_activations.1.notify_emails.#", "1").
						CheckEqual("hostname_activations.1.notify_emails.0", "nomail@akamai.com").
						Build(),
				},
			},
		},
		"happy path - with paging": {
			init: func(m *papi.Mock) {
				m.On("ListPropertyHostnameActivations", testutils.MockContext, papi.ListPropertyHostnameActivationsRequest{
					ContractID: "ctr_1",
					GroupID:    "grp_1",
					PropertyID: "prp_1",
					Offset:     0,
					Limit:      999,
				}).Return(&papi.ListPropertyHostnameActivationsResponse{
					AccountID:  "act_1",
					ContractID: "ctr_1",
					GroupID:    "grp_1",
					HostnameActivations: papi.HostnameActivationsList{
						TotalItems:       1100,
						CurrentItemCount: 999,
						Items:            activations999,
					},
				}, nil).Times(3)
				m.On("ListPropertyHostnameActivations", testutils.MockContext, papi.ListPropertyHostnameActivationsRequest{
					ContractID: "ctr_1",
					GroupID:    "grp_1",
					PropertyID: "prp_1",
					Offset:     999,
					Limit:      999,
				}).Return(&papi.ListPropertyHostnameActivationsResponse{
					AccountID:  "act_1",
					ContractID: "ctr_1",
					GroupID:    "grp_1",
					HostnameActivations: papi.HostnameActivationsList{
						TotalItems:       1100,
						CurrentItemCount: 101,
						Items:            activations101,
					},
				}, nil).Times(3)
			},
			steps: []resource.TestStep{
				{
					Config: testutils.LoadFixtureString(t, "testdata/TestDataPropertyHostnameActivations/valid_production.tf"),
					Check: baseChecker.
						CheckEqual("hostname_activations.#", "1100").
						CheckEqual("hostname_activations.0.hostname_activation_id", "0").
						CheckEqual("hostname_activations.0.network", "PRODUCTION").
						CheckEqual("hostname_activations.1099.hostname_activation_id", "1099").
						Build(),
				},
			},
		},
		"happy path - with limit paging": {
			init: func(m *papi.Mock) {
				m.On("ListPropertyHostnameActivations", testutils.MockContext, papi.ListPropertyHostnameActivationsRequest{
					ContractID: "1",
					GroupID:    "1",
					PropertyID: "1",
					Offset:     0,
					Limit:      999,
				}).Return(&papi.ListPropertyHostnameActivationsResponse{
					AccountID:  "act_1",
					ContractID: "ctr_1",
					GroupID:    "grp_1",
					HostnameActivations: papi.HostnameActivationsList{
						TotalItems:       999,
						CurrentItemCount: 999,
						Items:            activations999,
					},
				}, nil).Times(3)
			},
			steps: []resource.TestStep{
				{
					Config: testutils.LoadFixtureString(t, "testdata/TestDataPropertyHostnameActivations/valid.tf"),
					Check: baseChecker.
						CheckEqual("property_id", "1").
						CheckEqual("hostname_activations.#", "999").
						CheckEqual("hostname_activations.0.hostname_activation_id", "0").
						CheckEqual("hostname_activations.0.network", "PRODUCTION").
						Build(),
				},
			},
		},
		"happy path - empty hostname activations list": {
			init: func(m *papi.Mock) {
				m.On("ListPropertyHostnameActivations", testutils.MockContext, papi.ListPropertyHostnameActivationsRequest{
					ContractID: "1",
					GroupID:    "1",
					PropertyID: "1",
					Offset:     0,
					Limit:      999,
				}).Return(&papi.ListPropertyHostnameActivationsResponse{
					AccountID:  "act_1",
					ContractID: "ctr_1",
					GroupID:    "grp_1",
					HostnameActivations: papi.HostnameActivationsList{
						TotalItems: 0,
						Items:      []papi.HostnameActivationListItem{},
					},
				}, nil).Times(3)
			},
			steps: []resource.TestStep{
				{
					Config: testutils.LoadFixtureString(t, "testdata/TestDataPropertyHostnameActivations/valid.tf"),
					Check: baseChecker.
						CheckEqual("property_id", "1").
						CheckEqual("hostname_activations.#", "0").
						CheckMissing("property_name").
						CheckMissing("hostname_activations.0.notify_emails.0").
						CheckMissing("hostname_activations.0.activation_type").
						CheckMissing("hostname_activations.0.note").
						CheckMissing("hostname_activations.0.status").
						CheckMissing("hostname_activations.0.network").
						CheckMissing("hostname_activations.0.submit_date").
						CheckMissing("hostname_activations.0.update_date").
						CheckMissing("hostname_activations.0.hostname_activation_id").
						CheckMissing("hostname_activations.0.notify_emails.#").
						Build(),
				},
			},
		},
		"error response from api": {
			init: func(m *papi.Mock) {
				m.On("ListPropertyHostnameActivations", testutils.MockContext, papi.ListPropertyHostnameActivationsRequest{
					ContractID: "ctr_1",
					GroupID:    "grp_1",
					PropertyID: "prp_1",
					Offset:     0,
					Limit:      999,
				}).Return(nil, fmt.Errorf("oops"))
			},
			steps: []resource.TestStep{
				{
					Config:      testutils.LoadFixtureString(t, "testdata/TestDataPropertyHostnameActivations/valid_production.tf"),
					ExpectError: regexp.MustCompile("oops"),
				},
			},
		},
		"missing required argument property_id": {
			steps: []resource.TestStep{
				{
					Config:      testutils.LoadFixtureString(t, "testdata/TestDataPropertyHostnameActivations/missing_property_id.tf"),
					ExpectError: regexp.MustCompile(`The argument "property_id" is required, but no definition was found`),
				},
			},
		},
	}

	for name, test := range tests {
		t.Run(name, func(t *testing.T) {
			t.Parallel()
			client := &papi.Mock{}
			hapiClient := &hapi.Mock{}

			if test.init != nil {
				test.init(client)
			}

			useClient(client, hapiClient, func() {
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
