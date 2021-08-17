package dns

import (
	dns "github.com/akamai/AkamaiOPEN-edgegrid-golang/v2/pkg/configdns"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/resource"
	"github.com/stretchr/testify/mock"
	"net/http"
	"os"
	"testing"
)

func TestResDnsZone(t *testing.T) {
	zone := &dns.ZoneResponse{
		ContractID:      "ctr1",
		Zone:            "primaryexampleterraform.io",
		Type:            "primary",
		Comment:         "This is a test primary zone",
		SignAndServe:    false,
		ActivationState: "PENDING",
	}
	recordsetsResp := &dns.RecordSetResponse{Recordsets: make([]dns.Recordset, 2, 2)}

	// This test peforms a full life-cycle (CRUD) test
	t.Run("lifecycle test", func(t *testing.T) {
		client := &mockdns{}

		getCall := client.On("GetZone",
			mock.Anything, // ctx is irrelevant for this test
			zone.Zone,
		).Return(nil, &dns.Error{
			StatusCode: http.StatusNotFound,
		})

		client.On("CreateZone",
			mock.Anything, // ctx is irrelevant for this test
			mock.AnythingOfType("*dns.ZoneCreate"),
			mock.AnythingOfType("dns.ZoneQueryString"),
			true,
		).Return(nil).Run(func(args mock.Arguments) {
			getCall.ReturnArguments = mock.Arguments{zone, nil}
		})

		client.On("UpdateZone",
			mock.Anything, // ctx is irrelevant for this test
			mock.AnythingOfType("*dns.ZoneCreate"),
			mock.AnythingOfType("dns.ZoneQueryString"),
		).Return(nil).Run(func(args mock.Arguments) {
			zone.Comment = "This is an updated test primary zone"
		})

		client.On("SaveChangelist",
			mock.Anything, // ctx is irrelevant for this test
			mock.AnythingOfType("*dns.ZoneCreate"),
		).Return(nil)

		client.On("SubmitChangelist",
			mock.Anything, // ctx is irrelevant for this test
			mock.AnythingOfType("*dns.ZoneCreate"),
		).Return(nil)

		client.On("GetRecordsets",
			mock.Anything, // ctx is irrelevant for this test
			zone.Zone,
			mock.AnythingOfType("[]dns.RecordsetQueryArgs"),
		).Return(recordsetsResp, nil)

		dataSourceName := "akamai_dns_zone.primary_test_zone"

		// work around to skip Delete which fails intentionally
		os.Setenv("DNS_ZONE_SKIP_DELETE", "")
		defer os.Unsetenv("DNS_ZONE_SKIP_DELETE")
		useClient(client, func() {
			resource.UnitTest(t, resource.TestCase{
				PreCheck:  func() { testAccPreCheck(t) },
				Providers: testAccProviders,
				Steps: []resource.TestStep{
					{
						Config: loadFixtureString("testdata/TestResDnsZone/create_primary.tf"),
						Check: resource.ComposeTestCheckFunc(
							resource.TestCheckResourceAttr(dataSourceName, "zone", "primaryexampleterraform.io"),
						),
					},
					{
						Config: loadFixtureString("testdata/TestResDnsZone/update_primary.tf"),
						Check: resource.ComposeTestCheckFunc(
							resource.TestCheckResourceAttr(dataSourceName, "zone", "primaryexampleterraform.io"),
						),
					},
				},
			})
		})

		client.AssertExpectations(t)
	})
}
