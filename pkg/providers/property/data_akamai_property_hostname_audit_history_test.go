package property

import (
	"fmt"
	"regexp"
	"testing"

	"github.com/akamai/AkamaiOPEN-edgegrid-golang/v14/pkg/papi"
	"github.com/akamai/terraform-provider-akamai/v11/internal/edgegrid"
	"github.com/akamai/terraform-provider-akamai/v11/pkg/common/test"
	"github.com/akamai/terraform-provider-akamai/v11/pkg/common/testutils"
	"github.com/hashicorp/terraform-plugin-testing/helper/resource"
)

func TestDataPropertyHostnameAuditHistory(t *testing.T) {
	t.Parallel()

	tests := map[string]struct {
		init  func(*papi.Mock)
		steps []resource.TestStep
	}{
		"happy path - single history entry": {
			init: func(m *papi.Mock) {
				m.On("GetAuditHistory", testutils.MockContext, papi.GetAuditHistoryRequest{
					Hostname: "example.com",
				}).Return(&papi.GetAuditHistoryResponse{
					Hostname: "example.com",
					History: papi.HostnameHistory{
						Items: []papi.HostnameHistoryItem{
							{
								Action:               "ADD",
								CertProvisioningType: "CPS_MANAGED",
								CnameTo:              "example.com.edgekey.net",
								ContractID:           "ctr_123",
								EdgeHostnameID:       "ehn_123",
								GroupID:              "grp_123",
								Network:              "PRODUCTION",
								PropertyID:           "prp_123",
								Timestamp:            "2023-10-26T12:00:00Z",
								User:                 "user_123",
							},
						},
					},
				}, nil).Times(3)
			},
			steps: []resource.TestStep{
				{
					Config: testutils.LoadFixtureString(t, "testdata/TestDataPropertyHostnameAuditHistory/valid.tf"),
					Check: test.NewStateChecker("data.akamai_property_hostname_audit_history.test").
						CheckEqual("hostname", "example.com").
						CheckEqual("history.#", "1").
						CheckEqual("history.0.action", "ADD").
						CheckEqual("history.0.cert_provisioning_type", "CPS_MANAGED").
						CheckEqual("history.0.cname_to", "example.com.edgekey.net").
						CheckEqual("history.0.contract_id", "ctr_123").
						CheckEqual("history.0.edge_hostname_id", "ehn_123").
						CheckEqual("history.0.group_id", "grp_123").
						CheckEqual("history.0.network", "PRODUCTION").
						CheckEqual("history.0.property_id", "prp_123").
						CheckEqual("history.0.timestamp", "2023-10-26T12:00:00Z").
						CheckEqual("history.0.user", "user_123").
						Build(),
				},
			},
		},
		"happy path - multiple history entries": {
			init: func(m *papi.Mock) {
				m.On("GetAuditHistory", testutils.MockContext, papi.GetAuditHistoryRequest{
					Hostname: "example.com",
				}).Return(&papi.GetAuditHistoryResponse{
					Hostname: "example.com",
					History: papi.HostnameHistory{
						Items: []papi.HostnameHistoryItem{
							{
								Action:               "ADD",
								CertProvisioningType: "CPS_MANAGED",
								CnameTo:              "example.com.edgekey.net",
								ContractID:           "ctr_123",
								EdgeHostnameID:       "ehn_123",
								GroupID:              "grp_123",
								Network:              "STAGING",
								PropertyID:           "prp_123",
								Timestamp:            "2023-10-26T12:00:00Z",
								User:                 "user_123",
							},
							{
								Action:               "ACTIVATE",
								CertProvisioningType: "CPS_MANAGED",
								CnameTo:              "example.com.edgekey.net",
								ContractID:           "ctr_123",
								EdgeHostnameID:       "ehn_123",
								GroupID:              "grp_123",
								Network:              "STAGING",
								PropertyID:           "prp_123",
								Timestamp:            "2023-10-27T14:30:00Z",
								User:                 "user_123",
							},
							{
								Action:               "MODIFY",
								CertProvisioningType: "DEFAULT",
								CnameTo:              "example.com.edgesuite.net",
								ContractID:           "ctr_123",
								EdgeHostnameID:       "ehn_123",
								GroupID:              "grp_123",
								Network:              "PRODUCTION",
								PropertyID:           "prp_123",
								Timestamp:            "2023-10-28T16:45:00Z",
								User:                 "user_123",
							},
						},
					},
				}, nil).Times(3)
			},
			steps: []resource.TestStep{
				{
					Config: testutils.LoadFixtureString(t, "testdata/TestDataPropertyHostnameAuditHistory/valid.tf"),
					Check: test.NewStateChecker("data.akamai_property_hostname_audit_history.test").
						CheckEqual("hostname", "example.com").
						CheckEqual("history.#", "3").
						// First entry
						CheckEqual("history.0.action", "ADD").
						CheckEqual("history.0.cert_provisioning_type", "CPS_MANAGED").
						CheckEqual("history.0.cname_to", "example.com.edgekey.net").
						CheckEqual("history.0.contract_id", "ctr_123").
						CheckEqual("history.0.edge_hostname_id", "ehn_123").
						CheckEqual("history.0.group_id", "grp_123").
						CheckEqual("history.0.network", "STAGING").
						CheckEqual("history.0.property_id", "prp_123").
						CheckEqual("history.0.timestamp", "2023-10-26T12:00:00Z").
						CheckEqual("history.0.user", "user_123").
						// Second entry
						CheckEqual("history.1.action", "ACTIVATE").
						CheckEqual("history.1.cert_provisioning_type", "CPS_MANAGED").
						CheckEqual("history.1.cname_to", "example.com.edgekey.net").
						CheckEqual("history.1.contract_id", "ctr_123").
						CheckEqual("history.1.edge_hostname_id", "ehn_123").
						CheckEqual("history.1.group_id", "grp_123").
						CheckEqual("history.1.network", "STAGING").
						CheckEqual("history.1.property_id", "prp_123").
						CheckEqual("history.1.timestamp", "2023-10-27T14:30:00Z").
						CheckEqual("history.1.user", "user_123").
						// Third entry
						CheckEqual("history.2.action", "MODIFY").
						CheckEqual("history.2.cert_provisioning_type", "DEFAULT").
						CheckEqual("history.2.cname_to", "example.com.edgesuite.net").
						CheckEqual("history.2.contract_id", "ctr_123").
						CheckEqual("history.2.edge_hostname_id", "ehn_123").
						CheckEqual("history.2.group_id", "grp_123").
						CheckEqual("history.2.network", "PRODUCTION").
						CheckEqual("history.2.property_id", "prp_123").
						CheckEqual("history.2.timestamp", "2023-10-28T16:45:00Z").
						CheckEqual("history.2.user", "user_123").
						Build(),
				},
			},
		},
		"happy path - empty history": {
			init: func(m *papi.Mock) {
				m.On("GetAuditHistory", testutils.MockContext, papi.GetAuditHistoryRequest{
					Hostname: "example.com",
				}).Return(&papi.GetAuditHistoryResponse{
					Hostname: "example.com",
					History: papi.HostnameHistory{
						Items: []papi.HostnameHistoryItem{},
					},
				}, nil).Times(3)
			},
			steps: []resource.TestStep{
				{
					Config: testutils.LoadFixtureString(t, "testdata/TestDataPropertyHostnameAuditHistory/valid.tf"),
					Check: test.NewStateChecker("data.akamai_property_hostname_audit_history.test").
						CheckEqual("hostname", "example.com").
						CheckEqual("history.#", "0").
						Build(),
				},
			},
		},
		"error - API error": {
			init: func(m *papi.Mock) {
				m.On("GetAuditHistory", testutils.MockContext, papi.GetAuditHistoryRequest{
					Hostname: "example.com",
				}).Return(nil, fmt.Errorf("API error: could not fetch audit history")).Once()
			},
			steps: []resource.TestStep{
				{
					Config:      testutils.LoadFixtureString(t, "testdata/TestDataPropertyHostnameAuditHistory/valid.tf"),
					ExpectError: regexp.MustCompile("API error: could not fetch audit history"),
				},
			},
		},
		"error - missing required hostname attribute": {
			steps: []resource.TestStep{
				{
					Config:      testutils.LoadFixtureString(t, "testdata/TestDataPropertyHostnameAuditHistory/missing_hostname.tf"),
					ExpectError: regexp.MustCompile(`The argument "hostname" is required, but no definition was found`),
				},
			},
		},
		"error - empty hostname": {
			steps: []resource.TestStep{
				{
					Config:      testutils.LoadFixtureString(t, "testdata/TestDataPropertyHostnameAuditHistory/empty_hostname.tf"),
					ExpectError: regexp.MustCompile(`Attribute hostname string length must be between 1 and 253`),
				},
			},
		},
	}

	for name, test := range tests {
		t.Run(name, func(t *testing.T) {
			t.Parallel()
			client := edgegrid.NewTestClient()

			if test.init != nil {
				test.init(client.PAPI)
			}

			resource.UnitTest(t, resource.TestCase{
				ProtoV6ProviderFactories: testutils.NewTestProtoV6ProviderFactory(client, NewSubprovider()),
				IsUnitTest:               true,
				Steps:                    test.steps,
			})

			client.PAPI.AssertExpectations(t)
		})
	}
}
