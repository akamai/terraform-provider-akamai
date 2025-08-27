package property

import (
	"fmt"
	"regexp"
	"testing"

	"github.com/akamai/AkamaiOPEN-edgegrid-golang/v12/pkg/hapi"
	"github.com/akamai/AkamaiOPEN-edgegrid-golang/v12/pkg/papi"
	"github.com/akamai/terraform-provider-akamai/v8/pkg/common/test"
	"github.com/akamai/terraform-provider-akamai/v8/pkg/common/testutils"
	"github.com/hashicorp/terraform-plugin-testing/helper/resource"
	"github.com/stretchr/testify/mock"
)

func TestDataPropertyHostnamesDiff(t *testing.T) {
	t.Parallel()
	baseChecker := test.NewStateChecker("data.akamai_property_hostnames_diff.diff").
		CheckEqual("property_id", "prp_1").
		CheckEqual("group_id", "grp_1").
		CheckEqual("contract_id", "ctr_1").
		CheckEqual("account_id", "act_1")

	hostnamesDiff3 := papi.HostnamesDiffResponseItems{
		Items: []papi.HostnameDiffItem{
			{
				CnameFrom:                      "www.example-stag1.com",
				ProductionCertProvisioningType: "",
				ProductionCnameTo:              "",
				ProductionCnameType:            "",
				ProductionEdgeHostnameID:       "",
				StagingCertProvisioningType:    "DEFAULT",
				StagingCnameTo:                 "www.example-stag1-test.com",
				StagingCnameType:               "EDGE_HOSTNAME",
				StagingEdgeHostnameID:          "ehn_1",
			},
			{
				CnameFrom:                      "www.example-prod1.com",
				ProductionCertProvisioningType: "DEFAULT",
				ProductionCnameTo:              "www.example-prod1-test.com",
				ProductionCnameType:            "EDGE_HOSTNAME",
				ProductionEdgeHostnameID:       "ehn_1",
				StagingCertProvisioningType:    "",
				StagingCnameTo:                 "",
				StagingCnameType:               "",
				StagingEdgeHostnameID:          "",
			},
			{
				CnameFrom:                      "www.example-prod2.com",
				ProductionCertProvisioningType: "DEFAULT",
				ProductionCnameTo:              "www.example-prod2-test.com",
				ProductionCnameType:            "EDGE_HOSTNAME",
				ProductionEdgeHostnameID:       "ehn_1",
				StagingCertProvisioningType:    "",
				StagingCnameTo:                 "",
				StagingCnameType:               "",
				StagingEdgeHostnameID:          "",
			},
		},
		TotalItems:       3,
		CurrentItemCount: 3,
	}

	hostnamesDiff999 := make([]papi.HostnameDiffItem, 999)
	for i := 0; i < 999; i++ {
		hostnamesDiff999[i] = papi.HostnameDiffItem{
			CnameFrom:                      fmt.Sprintf("www.example-prod%d.com", i),
			ProductionCertProvisioningType: "DEFAULT",
			ProductionCnameTo:              fmt.Sprintf("www.example-prod%d-test.com", i),
			ProductionCnameType:            "EDGE_HOSTNAME",
			ProductionEdgeHostnameID:       "ehn_1",
			StagingCertProvisioningType:    "",
			StagingCnameTo:                 "",
			StagingCnameType:               "",
			StagingEdgeHostnameID:          "",
		}
	}
	hostnamesDiff101 := make([]papi.HostnameDiffItem, 101)
	for i := 0; i < 101; i++ {
		hostnamesDiff101[i] = papi.HostnameDiffItem{
			CnameFrom:                      fmt.Sprintf("www.example-stag%d.com", i),
			ProductionCertProvisioningType: "",
			ProductionCnameTo:              "",
			ProductionCnameType:            "",
			ProductionEdgeHostnameID:       "",
			StagingCertProvisioningType:    "DEFAULT",
			StagingCnameTo:                 fmt.Sprintf("www.example-stag%d-test.com", i),
			StagingCnameType:               "EDGE_HOSTNAME",
			StagingEdgeHostnameID:          "ehn_1",
		}
	}

	tests := map[string]struct {
		init  func(*papi.Mock)
		steps []resource.TestStep
	}{
		"happy path": {
			init: func(m *papi.Mock) {
				mockGetActivePropertyHostnamesDiff(m, "prp_1", "grp_1", "ctr_1", 0, 999,
					3, hostnamesDiff3, nil)
			},
			steps: []resource.TestStep{
				{
					Config: testutils.LoadFixtureString(t, "testdata/TestDataPropertyHostnamesDiff/valid.tf"),
					Check: baseChecker.
						CheckEqual("hostnames.#", "3").
						CheckEqual("hostnames.0.cname_from", "www.example-stag1.com").
						CheckEqual("hostnames.0.staging_cname_type", "EDGE_HOSTNAME").
						CheckEqual("hostnames.0.production_cname_type", "").
						CheckEqual("hostnames.0.staging_edge_hostname_id", "ehn_1").
						CheckEqual("hostnames.0.production_edge_hostname_id", "").
						CheckEqual("hostnames.0.staging_cert_provisioning_type", "DEFAULT").
						CheckEqual("hostnames.0.production_cert_provisioning_type", "").
						CheckEqual("hostnames.0.staging_cname_to", "www.example-stag1-test.com").
						CheckEqual("hostnames.0.production_cname_to", "").
						CheckEqual("hostnames.1.cname_from", "www.example-prod1.com").
						CheckEqual("hostnames.1.staging_cname_type", "").
						CheckEqual("hostnames.1.production_cname_type", "EDGE_HOSTNAME").
						CheckEqual("hostnames.1.staging_edge_hostname_id", "").
						CheckEqual("hostnames.1.production_edge_hostname_id", "ehn_1").
						CheckEqual("hostnames.1.staging_cert_provisioning_type", "").
						CheckEqual("hostnames.1.production_cert_provisioning_type", "DEFAULT").
						CheckEqual("hostnames.1.staging_cname_to", "").
						CheckEqual("hostnames.1.production_cname_to", "www.example-prod1-test.com").
						CheckEqual("hostnames.2.cname_from", "www.example-prod2.com").
						CheckEqual("hostnames.2.staging_cname_type", "").
						CheckEqual("hostnames.2.production_cname_type", "EDGE_HOSTNAME").
						CheckEqual("hostnames.2.staging_edge_hostname_id", "").
						CheckEqual("hostnames.2.production_edge_hostname_id", "ehn_1").
						CheckEqual("hostnames.2.staging_cert_provisioning_type", "").
						CheckEqual("hostnames.2.production_cert_provisioning_type", "DEFAULT").
						CheckEqual("hostnames.2.staging_cname_to", "").
						CheckEqual("hostnames.2.production_cname_to", "www.example-prod2-test.com").
						Build(),
				},
			},
		},
		"happy path - no contract and group": {
			init: func(m *papi.Mock) {
				mockGetActivePropertyHostnamesDiff(m, "prp_1", "", "", 0, 999,
					3, hostnamesDiff3, nil)
			},
			steps: []resource.TestStep{
				{
					Config: testutils.LoadFixtureString(t, "testdata/TestDataPropertyHostnamesDiff/valid_no_contract_and_group.tf"),
					Check: baseChecker.
						CheckEqual("hostnames.#", "3").
						CheckEqual("hostnames.0.cname_from", "www.example-stag1.com").
						CheckEqual("hostnames.0.staging_cname_type", "EDGE_HOSTNAME").
						CheckEqual("hostnames.0.production_cname_type", "").
						CheckEqual("hostnames.0.staging_edge_hostname_id", "ehn_1").
						CheckEqual("hostnames.0.production_edge_hostname_id", "").
						CheckEqual("hostnames.0.staging_cert_provisioning_type", "DEFAULT").
						CheckEqual("hostnames.0.production_cert_provisioning_type", "").
						CheckEqual("hostnames.0.staging_cname_to", "www.example-stag1-test.com").
						CheckEqual("hostnames.0.production_cname_to", "").
						CheckEqual("hostnames.1.cname_from", "www.example-prod1.com").
						CheckEqual("hostnames.1.staging_cname_type", "").
						CheckEqual("hostnames.1.production_cname_type", "EDGE_HOSTNAME").
						CheckEqual("hostnames.1.staging_edge_hostname_id", "").
						CheckEqual("hostnames.1.production_edge_hostname_id", "ehn_1").
						CheckEqual("hostnames.1.staging_cert_provisioning_type", "").
						CheckEqual("hostnames.1.production_cert_provisioning_type", "DEFAULT").
						CheckEqual("hostnames.1.staging_cname_to", "").
						CheckEqual("hostnames.1.production_cname_to", "www.example-prod1-test.com").
						CheckEqual("hostnames.2.cname_from", "www.example-prod2.com").
						CheckEqual("hostnames.2.staging_cname_type", "").
						CheckEqual("hostnames.2.production_cname_type", "EDGE_HOSTNAME").
						CheckEqual("hostnames.2.staging_edge_hostname_id", "").
						CheckEqual("hostnames.2.production_edge_hostname_id", "ehn_1").
						CheckEqual("hostnames.2.staging_cert_provisioning_type", "").
						CheckEqual("hostnames.2.production_cert_provisioning_type", "DEFAULT").
						CheckEqual("hostnames.2.staging_cname_to", "").
						CheckEqual("hostnames.2.production_cname_to", "www.example-prod2-test.com").
						Build(),
				},
			},
		},
		"happy path - empty diff items": {
			init: func(m *papi.Mock) {
				mockGetActivePropertyHostnamesDiff(m, "prp_1", "grp_1", "ctr_1", 0, 999,
					3, papi.HostnamesDiffResponseItems{
						Items: []papi.HostnameDiffItem{},
					}, nil)
			},
			steps: []resource.TestStep{
				{
					Config: testutils.LoadFixtureString(t, "testdata/TestDataPropertyHostnamesDiff/valid.tf"),
					Check: baseChecker.
						CheckEqual("hostnames.#", "0").
						Build(),
				},
			},
		},
		"happy path - with paging": {
			init: func(m *papi.Mock) {
				mockGetActivePropertyHostnamesDiff(m, "prp_1", "grp_1", "ctr_1", 0, 999,
					3, papi.HostnamesDiffResponseItems{
						Items:            hostnamesDiff999,
						CurrentItemCount: 999,
						TotalItems:       1100,
					}, nil)
				mockGetActivePropertyHostnamesDiff(m, "prp_1", "grp_1", "ctr_1", 999, 999,
					3, papi.HostnamesDiffResponseItems{
						Items:            hostnamesDiff101,
						CurrentItemCount: 101,
						TotalItems:       1100,
					}, nil)
			},
			steps: []resource.TestStep{
				{
					Config: testutils.LoadFixtureString(t, "testdata/TestDataPropertyHostnamesDiff/valid.tf"),
					Check: baseChecker.
						CheckEqual("hostnames.#", "1100").
						CheckEqual("hostnames.0.cname_from", "www.example-prod0.com").
						CheckEqual("hostnames.0.staging_cname_type", "").
						CheckEqual("hostnames.0.production_cname_type", "EDGE_HOSTNAME").
						CheckEqual("hostnames.0.staging_edge_hostname_id", "").
						CheckEqual("hostnames.0.production_edge_hostname_id", "ehn_1").
						CheckEqual("hostnames.0.staging_cert_provisioning_type", "").
						CheckEqual("hostnames.0.production_cert_provisioning_type", "DEFAULT").
						CheckEqual("hostnames.0.staging_cname_to", "").
						CheckEqual("hostnames.0.production_cname_to", "www.example-prod0-test.com").
						CheckEqual("hostnames.998.cname_from", "www.example-prod998.com").
						CheckEqual("hostnames.998.production_cname_to", "www.example-prod998-test.com").
						CheckEqual("hostnames.999.cname_from", "www.example-stag0.com").
						CheckEqual("hostnames.999.staging_cname_to", "www.example-stag0-test.com").
						CheckEqual("hostnames.1099.cname_from", "www.example-stag100.com").
						CheckEqual("hostnames.1099.staging_cname_to", "www.example-stag100-test.com").
						Build(),
				},
			},
		},
		"happy path - with limit paging": {
			init: func(m *papi.Mock) {
				mockGetActivePropertyHostnamesDiff(m, "prp_1", "grp_1", "ctr_1", 0, 999,
					3, papi.HostnamesDiffResponseItems{
						Items:            hostnamesDiff999,
						CurrentItemCount: 999,
						TotalItems:       999,
					}, nil)
			},
			steps: []resource.TestStep{
				{
					Config: testutils.LoadFixtureString(t, "testdata/TestDataPropertyHostnamesDiff/valid.tf"),
					Check: baseChecker.
						CheckEqual("hostnames.#", "999").
						CheckEqual("hostnames.0.cname_from", "www.example-prod0.com").
						CheckEqual("hostnames.0.staging_cname_type", "").
						CheckEqual("hostnames.0.production_cname_type", "EDGE_HOSTNAME").
						CheckEqual("hostnames.0.staging_edge_hostname_id", "").
						CheckEqual("hostnames.0.production_edge_hostname_id", "ehn_1").
						CheckEqual("hostnames.0.staging_cert_provisioning_type", "").
						CheckEqual("hostnames.0.production_cert_provisioning_type", "DEFAULT").
						CheckEqual("hostnames.0.staging_cname_to", "").
						CheckEqual("hostnames.0.production_cname_to", "www.example-prod0-test.com").
						CheckEqual("hostnames.998.cname_from", "www.example-prod998.com").
						CheckEqual("hostnames.998.production_cname_to", "www.example-prod998-test.com").
						Build(),
				},
			},
		},
		"error response from api": {
			init: func(m *papi.Mock) {
				mockGetActivePropertyHostnamesDiff(m, "prp_1", "grp_1", "ctr_1", 0, 999,
					3, papi.HostnamesDiffResponseItems{}, fmt.Errorf("oops"))
			},
			steps: []resource.TestStep{
				{
					Config:      testutils.LoadFixtureString(t, "testdata/TestDataPropertyHostnamesDiff/valid.tf"),
					ExpectError: regexp.MustCompile("oops"),
				},
			},
		},
		"missing required argument property_id": {
			steps: []resource.TestStep{
				{
					Config:      testutils.LoadFixtureString(t, "testdata/TestDataPropertyHostnamesDiff/missing_property_id.tf"),
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

func mockGetActivePropertyHostnamesDiff(m *papi.Mock, propertyID, groupID, contractID string, offset, limit, times int,
	hostnamesResponse papi.HostnamesDiffResponseItems, err error) *mock.Call {
	if err != nil {
		return m.On("GetActivePropertyHostnamesDiff", testutils.MockContext, papi.GetActivePropertyHostnamesDiffRequest{
			ContractID: contractID,
			GroupID:    groupID,
			PropertyID: propertyID,
			Offset:     offset,
			Limit:      limit,
		}).Return(nil, err)
	}
	return m.On("GetActivePropertyHostnamesDiff", testutils.MockContext, papi.GetActivePropertyHostnamesDiffRequest{
		ContractID: contractID,
		GroupID:    groupID,
		PropertyID: propertyID,
		Offset:     offset,
		Limit:      limit,
	}).Return(&papi.GetActivePropertyHostnamesDiffResponse{
		AccountID:  "act_1",
		ContractID: "ctr_1",
		GroupID:    "grp_1",
		PropertyID: propertyID,
		Hostnames:  hostnamesResponse,
	}, nil).Times(times)
}
