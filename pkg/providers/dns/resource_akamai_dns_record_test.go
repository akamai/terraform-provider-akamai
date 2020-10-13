package dns

import (
	"context"
	"net/http"
	"testing"

	dns "github.com/akamai/AkamaiOPEN-edgegrid-golang/v2/pkg/configdns"
	"github.com/akamai/AkamaiOPEN-edgegrid-golang/v2/pkg/session"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/resource"
	"github.com/stretchr/testify/mock"
)

func TestResDnsRecord(t *testing.T) {
	parseRData := dns.Client(session.Must(session.New())).ParseRData

	rec := &dns.RecordBody{
		Name:       "exampleterraform.io",
		RecordType: "A",
		TTL:        300,
		Target:     []string{"10.0.0.2", "10.0.0.3"},
		Active:     true,
	}

	parsedData := parseRData(context.Background(), "A", rec.Target)

	// This test peforms a full life-cycle (CRUD) test
	t.Run("create record", func(t *testing.T) {
		client := &mockdns{}

		getCall := client.On("GetRecord",
			mock.Anything, // ctx is irrelevant for this test
			"exampleterraform.io",
			"exampleterraform.io",
			"A",
		).Return(nil, &dns.Error{
			StatusCode: http.StatusNotFound,
		})

		client.On("CreateRecord",
			mock.Anything, // ctx is irrelevant for this test
			mock.AnythingOfType("*dns.RecordBody"),
			"exampleterraform.io",
			mock.Anything,
		).Return(nil).Run(func(args mock.Arguments) {
			getCall.ReturnArguments = mock.Arguments{rec, nil}
		})

		client.On("ProcessRdata",
			mock.Anything, // ctx is irrelevant for this test
			mock.AnythingOfType("[]string"),
			"A",
		).Return(rec.Target, nil)

		parseCall := client.On("ParseRData",
			mock.Anything,
			"A",
			mock.AnythingOfType("[]string"),
		).Return(parsedData)

		client.On("UpdateRecord",
			mock.Anything, // ctx is irrelevant for this test
			mock.AnythingOfType("*dns.RecordBody"),
			"exampleterraform.io",
			mock.Anything,
		).Return(nil).Run(func(args mock.Arguments) {

			rec.Target = []string{"10.0.0.4", "10.0.0.5"}

			parsedData = parseRData(context.Background(), "A", rec.Target)
			parseCall.ReturnArguments = mock.Arguments{parsedData}
		})

		client.On("DeleteRecord",
			mock.Anything, // ctx is irrelevant for this test
			mock.AnythingOfType("*dns.RecordBody"),
			"exampleterraform.io",
			mock.Anything,
		).Return(nil)

		dataSourceName := "akamai_dns_record.a_record"

		useClient(client, func() {
			resource.UnitTest(t, resource.TestCase{
				PreCheck:  func() { testAccPreCheck(t) },
				Providers: testAccProviders,
				Steps: []resource.TestStep{
					{
						Config: loadFixtureString("testdata/TestResDnsRecord/create_basic.tf"),
						Check: resource.ComposeTestCheckFunc(
							resource.TestCheckResourceAttr(dataSourceName, "zone", "exampleterraform.io"),
							resource.TestCheckResourceAttr(dataSourceName, "name", "exampleterraform.io"),
							resource.TestCheckResourceAttr(dataSourceName, "recordtype", "A"),
						),
					},
					{
						Config: loadFixtureString("testdata/TestResDnsRecord/update_basic.tf"),
						Check: resource.ComposeTestCheckFunc(
							resource.TestCheckResourceAttr(dataSourceName, "zone", "exampleterraform.io"),
							resource.TestCheckResourceAttr(dataSourceName, "name", "exampleterraform.io"),
							resource.TestCheckResourceAttr(dataSourceName, "recordtype", "A"),
						),
					},
				},
			})
		})

		client.AssertExpectations(t)
	})

	// This example tests attempting to create an A record that already exists on the server
	// It is not a full lifecycle test as the "update" occurs in the create method
	t.Run("update existing record", func(t *testing.T) {
		client := &mockdns{}

		client.On("GetRecord",
			mock.Anything, // ctx is irrelevant for this test
			"exampleterraform.io",
			"exampleterraform.io",
			"A",
		).Return(rec, nil)

		client.On("ProcessRdata",
			mock.Anything, // ctx is irrelevant for this test
			mock.AnythingOfType("[]string"),
			"A",
		).Return(rec.Target, nil)

		client.On("UpdateRecord",
			mock.Anything, // ctx is irrelevant for this test
			mock.AnythingOfType("*dns.RecordBody"),
			"exampleterraform.io",
			mock.Anything,
		).Return(nil)

		client.On("ParseRData",
			mock.Anything,
			"A",
			mock.AnythingOfType("[]string"),
		).Return(parsedData)

		client.On("DeleteRecord",
			mock.Anything, // ctx is irrelevant for this test
			mock.AnythingOfType("*dns.RecordBody"),
			"exampleterraform.io",
			mock.Anything,
		).Return(nil)

		dataSourceName := "akamai_dns_record.a_record"

		useClient(client, func() {
			resource.UnitTest(t, resource.TestCase{
				PreCheck:  func() { testAccPreCheck(t) },
				Providers: testAccProviders,
				Steps: []resource.TestStep{
					{
						// use the update config because the rec value was changed in the previous example
						Config: loadFixtureString("testdata/TestResDnsRecord/update_basic.tf"),
						Check: resource.ComposeTestCheckFunc(
							resource.TestCheckResourceAttr(dataSourceName, "zone", "exampleterraform.io"),
							resource.TestCheckResourceAttr(dataSourceName, "name", "exampleterraform.io"),
							resource.TestCheckResourceAttr(dataSourceName, "recordtype", "A"),
						),
					},
				},
			})
		})

		client.AssertExpectations(t)
	})

	// This test does an "update" by returning empty rdata which forces a new record overrite
	// It is not a full lifecycle test as the "update" occurs in the create method
	t.Run("save record", func(t *testing.T) {
		client := &mockdns{}

		client.On("GetRecord",
			mock.Anything, // ctx is irrelevant for this test
			"exampleterraform.io",
			"exampleterraform.io",
			"A",
		).Return(rec, nil)

		// return empty rdata to trigger the "save" codepath
		procCall := client.On("ProcessRdata",
			mock.Anything, // ctx is irrelevant for this test
			rec.Target,
			"A",
		).Return([]string{}, nil)

		client.On("CreateRecord",
			mock.Anything, // ctx is irrelevant for this test
			mock.AnythingOfType("*dns.RecordBody"),
			"exampleterraform.io",
			mock.Anything,
		).Return(nil).Run(func(mock.Arguments) {
			procCall.ReturnArguments = mock.Arguments{rec.Target, nil}
		})

		client.On("ParseRData",
			mock.Anything,
			"A",
			mock.AnythingOfType("[]string"),
		).Return(parsedData)

		client.On("DeleteRecord",
			mock.Anything, // ctx is irrelevant for this test
			mock.AnythingOfType("*dns.RecordBody"),
			"exampleterraform.io",
			mock.Anything,
		).Return(nil)

		dataSourceName := "akamai_dns_record.a_record"

		useClient(client, func() {
			resource.UnitTest(t, resource.TestCase{
				PreCheck:  func() { testAccPreCheck(t) },
				Providers: testAccProviders,
				Steps: []resource.TestStep{
					{
						Config: loadFixtureString("testdata/TestResDnsRecord/update_basic.tf"),
						Check: resource.ComposeTestCheckFunc(
							resource.TestCheckResourceAttr(dataSourceName, "zone", "exampleterraform.io"),
							resource.TestCheckResourceAttr(dataSourceName, "name", "exampleterraform.io"),
							resource.TestCheckResourceAttr(dataSourceName, "recordtype", "A"),
						),
					},
				},
			})
		})

		client.AssertExpectations(t)
	})

	soaRec := &dns.RecordBody{
		RecordType: "SOA",
		Name:       "exampleterraform.io",
		Target:     []string{"ns1.exampleterraform.io root@exampleterraform.io 123456789 3600 600 3600 3600"},
		TTL:        300,
	}

	t.Run("create soa record", func(t *testing.T) {
		client := &mockdns{}

		getCall := client.On("GetRecord",
			mock.Anything, // ctx is irrelevant for this test
			"exampleterraform.io",
			"@",
			"SOA",
		).Return(nil, &dns.Error{
			StatusCode: http.StatusNotFound,
		})

		client.On("CreateRecord",
			mock.Anything, // ctx is irrelevant for this test
			mock.AnythingOfType("*dns.RecordBody"),
			"exampleterraform.io",
			mock.Anything,
		).Return(nil).Run(func(mock.Arguments) {
			getCall.ReturnArguments = mock.Arguments{soaRec, nil}
		})

		client.On("ProcessRdata",
			mock.Anything, // ctx is irrelevant for this test
			mock.AnythingOfType("[]string"),
			"SOA",
		).Return(soaRec.Target, nil)

		client.On("ParseRData",
			mock.Anything,
			"SOA",
			mock.AnythingOfType("[]string"),
		).Return(parseRData(context.Background(), "SOA", soaRec.Target))

		client.On("DeleteRecord",
			mock.Anything, // ctx is irrelevant for this test
			mock.AnythingOfType("*dns.RecordBody"),
			"exampleterraform.io",
			mock.Anything,
		).Return(nil)

		dataSourceName := "akamai_dns_record.soa_record"

		useClient(client, func() {
			resource.UnitTest(t, resource.TestCase{
				PreCheck:  func() { testAccPreCheck(t) },
				Providers: testAccProviders,
				Steps: []resource.TestStep{
					{
						Config: loadFixtureString("testdata/TestResDnsRecord/create_soa.tf"),
						Check: resource.ComposeTestCheckFunc(
							resource.TestCheckResourceAttr(dataSourceName, "zone", "exampleterraform.io"),
							resource.TestCheckResourceAttr(dataSourceName, "recordtype", "SOA"),
						),
					},
				},
			})
		})

		client.AssertExpectations(t)
	})

	t.Run("update soa record", func(t *testing.T) {
		client := &mockdns{}

		client.On("GetRecord",
			mock.Anything, // ctx is irrelevant for this test
			"exampleterraform.io",
			"@",
			"SOA",
		).Return(soaRec, nil)

		client.On("ProcessRdata",
			mock.Anything, // ctx is irrelevant for this test
			mock.AnythingOfType("[]string"),
			"SOA",
		).Return(soaRec.Target, nil)

		client.On("UpdateRecord",
			mock.Anything, // ctx is irrelevant for this test
			mock.AnythingOfType("*dns.RecordBody"),
			"exampleterraform.io",
			mock.Anything,
		).Return(nil)

		client.On("ParseRData",
			mock.Anything,
			"SOA",
			mock.AnythingOfType("[]string"),
		).Return(parseRData(context.Background(), "SOA", soaRec.Target))

		client.On("DeleteRecord",
			mock.Anything, // ctx is irrelevant for this test
			mock.AnythingOfType("*dns.RecordBody"),
			"exampleterraform.io",
			mock.Anything,
		).Return(nil)

		dataSourceName := "akamai_dns_record.soa_record"

		useClient(client, func() {
			resource.UnitTest(t, resource.TestCase{
				PreCheck:  func() { testAccPreCheck(t) },
				Providers: testAccProviders,
				Steps: []resource.TestStep{
					{
						Config: loadFixtureString("testdata/TestResDnsRecord/create_soa.tf"),
						Check: resource.ComposeTestCheckFunc(
							resource.TestCheckResourceAttr(dataSourceName, "zone", "exampleterraform.io"),
							resource.TestCheckResourceAttr(dataSourceName, "recordtype", "SOA"),
						),
					},
				},
			})
		})

		client.AssertExpectations(t)
	})
}
