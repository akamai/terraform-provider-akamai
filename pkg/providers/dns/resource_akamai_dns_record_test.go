package dns

import (
	"context"
	"net/http"
	"testing"

	"github.com/akamai/AkamaiOPEN-edgegrid-golang/v5/pkg/dns"
	"github.com/akamai/AkamaiOPEN-edgegrid-golang/v5/pkg/session"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/resource"
	"github.com/stretchr/testify/mock"
)

func TestResDnsRecord(t *testing.T) {
	dnsClient := dns.Client(session.Must(session.New()))

	var rec *dns.RecordBody

	notFound := &dns.Error{
		StatusCode: http.StatusNotFound,
	}

	// This test peforms a full life-cycle (CRUD) test
	t.Run("lifecycle test", func(t *testing.T) {
		client := &dns.Mock{}

		getCall := client.On("GetRecord",
			mock.Anything, // ctx is irrelevant for this test
			mock.AnythingOfType("string"),
			mock.AnythingOfType("string"),
			mock.AnythingOfType("string"),
		).Return(nil, notFound)

		parseCall := client.On("ParseRData",
			mock.Anything,
			mock.AnythingOfType("string"),
			mock.AnythingOfType("[]string"),
		).Return(nil)

		procCall := client.On("ProcessRdata",
			mock.Anything, // ctx is irrelevant for this test
			mock.AnythingOfType("[]string"),
			mock.AnythingOfType("string"),
		).Return(nil, nil)

		updateArguments := func(args mock.Arguments) {
			rec = args.Get(1).(*dns.RecordBody)
			getCall.ReturnArguments = mock.Arguments{rec, nil}
			parseCall.ReturnArguments = mock.Arguments{
				dnsClient.ParseRData(context.Background(), rec.RecordType, rec.Target),
			}
			procCall.ReturnArguments = mock.Arguments{rec.Target, nil}
		}

		client.On("CreateRecord",
			mock.Anything, // ctx is irrelevant for this test
			mock.AnythingOfType("*dns.RecordBody"),
			mock.AnythingOfType("string"),
			mock.Anything,
		).Return(nil).Run(func(args mock.Arguments) {
			updateArguments(args)
		})

		client.On("UpdateRecord",
			mock.Anything, // ctx is irrelevant for this test
			mock.AnythingOfType("*dns.RecordBody"),
			mock.AnythingOfType("string"),
			mock.Anything,
		).Return(nil).Run(func(args mock.Arguments) {
			updateArguments(args)
		})

		client.On("DeleteRecord",
			mock.Anything, // ctx is irrelevant for this test
			mock.AnythingOfType("*dns.RecordBody"),
			mock.AnythingOfType("string"),
			mock.AnythingOfType("[]bool"),
		).Return(nil).Run(func(mock.Arguments) {
			getCall.ReturnArguments = mock.Arguments{nil, notFound}
		})

		dataSourceName := "akamai_dns_record.a_record"

		useClient(client, func() {
			resource.UnitTest(t, resource.TestCase{
				PreCheck:          func() { testAccPreCheck(t) },
				ProviderFactories: testAccProviders,
				Steps: []resource.TestStep{
					{
						Config: loadFixtureString("testdata/TestResDnsRecord/create_basic.tf"),
						Check: resource.ComposeTestCheckFunc(
							resource.TestCheckResourceAttr(dataSourceName, "recordtype", "A"),
						),
					},
					{
						Config: loadFixtureString("testdata/TestResDnsRecord/update_basic.tf"),
						Check: resource.ComposeTestCheckFunc(
							resource.TestCheckResourceAttr(dataSourceName, "recordtype", "A"),
						),
					},
				},
			})
		})

		client.AssertExpectations(t)
	})
}
