package property

import (
	"fmt"
	"regexp"
	"testing"
	"time"

	"github.com/akamai/AkamaiOPEN-edgegrid-golang/v10/pkg/hapi"
	"github.com/akamai/AkamaiOPEN-edgegrid-golang/v10/pkg/papi"
	"github.com/akamai/terraform-provider-akamai/v8/pkg/common/test"
	"github.com/akamai/terraform-provider-akamai/v8/pkg/common/testutils"
	"github.com/hashicorp/terraform-plugin-testing/helper/resource"
	"github.com/stretchr/testify/mock"
)

func TestHostnameActivationDataSource(t *testing.T) {
	t.Parallel()
	workdir := "testdata/TestDataPropertyHostnameActivation"
	commonStateChecker := test.NewStateChecker("data.akamai_property_hostname_activation.activation").
		CheckEqual("property_id", "1").
		CheckEqual("hostname_activation_id", "1").
		CheckEqual("contract_id", "1").
		CheckEqual("group_id", "1").
		CheckEqual("account_id", "1").
		CheckEqual("activation_type", "ACTIVATE").
		CheckEqual("network", "STAGING").
		CheckEqual("note", "sample note").
		CheckEqual("notify_emails.0", "unknown@akamai.com").
		CheckEqual("notify_emails.1", "unknown2@akamai.com").
		CheckEqual("property_name", "testName").
		CheckEqual("submit_date", "2001-01-01T01:11:11Z").
		CheckEqual("update_date", "2022-02-02T02:22:22Z")

	tests := map[string]struct {
		init  func(*papi.Mock)
		steps []resource.TestStep
		error *regexp.Regexp
	}{
		"happy path - with all optionals": {
			init: func(m *papi.Mock) {
				mockGetPropertyHostnameActivation(m, "prp_1", "grp_1", "ctr_1", "atv_1", true, false, 3)
			},
			steps: []resource.TestStep{
				{
					Config: testutils.LoadFixtureStringf(t, "%s/full.tf", workdir),
					Check: commonStateChecker.
						CheckEqual("property_id", "prp_1").
						CheckEqual("include_hostnames", "true").
						CheckEqual("hostnames.0.edge_hostname_id", "ehn_1").
						CheckEqual("hostnames.0.cname_from", "test.buckets.1.1111.com.edgesuite.net").
						CheckEqual("hostnames.0.cname_to", "test.buckets.1.com.edgesuite.net").
						CheckEqual("hostnames.0.cert_provisioning_type", "CPS_MANAGED").
						CheckEqual("hostnames.0.action", "ADD").
						Build(),
				},
			},
		},
		"happy path - only required": {
			init: func(m *papi.Mock) {
				mockGetPropertyHostnameActivation(m, "1", "", "", "1", false, false, 3)
			},
			steps: []resource.TestStep{
				{
					Config: testutils.LoadFixtureStringf(t, "%s/base.tf", workdir),
					Check:  commonStateChecker.Build(),
				},
			},
		},
		"missing required argument property_id": {
			steps: []resource.TestStep{
				{
					Config:      testutils.LoadFixtureString(t, "testdata/TestDataPropertyHostnameActivation/missing_property_id.tf"),
					ExpectError: regexp.MustCompile(`The argument "property_id" is required, but no definition was found`),
				},
			},
		},
		"missing required argument hostname_activation_id": {
			steps: []resource.TestStep{
				{
					Config:      testutils.LoadFixtureString(t, "testdata/TestDataPropertyHostnameActivation/missing_hostname_activation_id.tf"),
					ExpectError: regexp.MustCompile(`The argument "hostname_activation_id" is required, but no definition was\s+found`),
				},
			},
		},
		"error API response": {
			init: func(m *papi.Mock) {
				mockGetPropertyHostnameActivation(m, "prp_1", "grp_1", "ctr_1", "atv_1", true, true, 1)
			},
			steps: []resource.TestStep{
				{
					Config:      testutils.LoadFixtureStringf(t, "%s/full.tf", workdir),
					ExpectError: regexp.MustCompile("oops"),
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

func mockGetPropertyHostnameActivation(m *papi.Mock, propertyID, groupID, contractID, hostnameActivationID string, includeHostnames, throwsError bool, times int) *mock.Call {
	var err interface{}
	response := &papi.GetPropertyHostnameActivationResponse{
		AccountID:  "1",
		ContractID: "1",
		GroupID:    "1",
		HostnameActivation: papi.HostnameActivationGetItem{
			ActivationType:       "ACTIVATE",
			HostnameActivationID: "1",
			PropertyName:         "testName",
			PropertyID:           "1",
			Network:              "STAGING",
			Status:               "ACTIVE",
			SubmitDate:           time.Date(2001, 1, 01, 1, 11, 11, 0, time.UTC),
			UpdateDate:           time.Date(2022, 2, 02, 2, 22, 22, 0, time.UTC),
			Note:                 "sample note",
			NotifyEmails:         []string{"unknown@akamai.com", "unknown2@akamai.com"},
		},
	}
	if includeHostnames {
		response.HostnameActivation.Hostnames = []papi.PropertyHostnameItem{
			{
				EdgeHostnameID:       "ehn_1",
				CertProvisioningType: "CPS_MANAGED",
				CnameFrom:            "test.buckets.1.1111.com.edgesuite.net",
				CnameTo:              "test.buckets.1.com.edgesuite.net",
				Action:               "ADD",
			},
		}
	}
	if throwsError {
		response = nil
		err = fmt.Errorf("oops")
	}
	return m.On("GetPropertyHostnameActivation", testutils.MockContext, papi.GetPropertyHostnameActivationRequest{
		PropertyID:           propertyID,
		GroupID:              groupID,
		ContractID:           contractID,
		HostnameActivationID: hostnameActivationID,
		IncludeHostnames:     includeHostnames,
	}).Return(response, err).Times(times)
}
