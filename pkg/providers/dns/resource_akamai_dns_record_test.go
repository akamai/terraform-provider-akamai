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

func TestDnsRecordCreate(t *testing.T) {
	parseRData := dns.Client(session.Must(session.New())).ParseRData

	rec := &dns.RecordBody{
		Name:       "exampleterraform.io",
		RecordType: "A",
		TTL:        300,
		Target:     []string{"10.0.0.2", "10.0.0.3"},
		Active:     true,
	}

	parsedData := parseRData(context.Background(), "A", []string{"10.0.0.2", "10.0.0.3"})

	t.Run("create record", func(t *testing.T) {
		client := &mockdns{}

		client.On("GetRecord",
			mock.Anything, // ctx is irrelevant for this test
			"exampleterraform.io",
			"exampleterraform.io",
			"A",
		).Return(nil, &dns.Error{
			StatusCode: http.StatusNotFound,
		}).Once().Run(func(mock.Arguments) {
			client.On("GetRecord",
				mock.Anything, // ctx is irrelevant for this test
				"exampleterraform.io",
				"exampleterraform.io",
				"A",
			).Return(rec, nil)

			client.On("ProcessRdata",
				mock.Anything, // ctx is irrelevant for this test
				[]string{"10.0.0.2", "10.0.0.3"},
				"A",
			).Return([]string{"10.0.0.2", "10.0.0.3"}, nil)
		})

		client.On("CreateRecord",
			mock.Anything, // ctx is irrelevant for this test
			mock.AnythingOfType("*dns.RecordBody"),
			"exampleterraform.io",
			mock.Anything,
		).Return(nil)

		client.On("ParseRData",
			mock.Anything,
			"A",
			[]string{"10.0.0.2", "10.0.0.3"},
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
						Config: loadFixtureString("testdata/TestResDnsRecord/create_basic.tf"),
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

	t.Run("update record", func(t *testing.T) {
		client := &mockdns{}

		client.On("GetRecord",
			mock.Anything, // ctx is irrelevant for this test
			"exampleterraform.io",
			"exampleterraform.io",
			"A",
		).Return(rec, nil)

		client.On("ProcessRdata",
			mock.Anything, // ctx is irrelevant for this test
			[]string{"10.0.0.2", "10.0.0.3"},
			"A",
		).Return([]string{"10.0.0.2", "10.0.0.3"}, nil)

		client.On("UpdateRecord",
			mock.Anything, // ctx is irrelevant for this test
			mock.AnythingOfType("*dns.RecordBody"),
			"exampleterraform.io",
			mock.Anything,
		).Return(nil)

		client.On("ParseRData",
			mock.Anything,
			"A",
			[]string{"10.0.0.2", "10.0.0.3"},
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
						Config: loadFixtureString("testdata/TestResDnsRecord/create_basic.tf"),
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

	t.Run("save record", func(t *testing.T) {
		client := &mockdns{}

		client.On("GetRecord",
			mock.Anything, // ctx is irrelevant for this test
			"exampleterraform.io",
			"exampleterraform.io",
			"A",
		).Return(rec, nil)

		// return empty rdata to trigger the "save" codepath
		client.On("ProcessRdata",
			mock.Anything, // ctx is irrelevant for this test
			[]string{"10.0.0.2", "10.0.0.3"},
			"A",
		).Return([]string{}, nil).Once().Run(func(mock.Arguments) {
			// return valid rdata so save succeeds
			client.On("ProcessRdata",
				mock.Anything, // ctx is irrelevant for this test
				[]string{"10.0.0.2", "10.0.0.3"},
				"A",
			).Return([]string{"10.0.0.2", "10.0.0.3"}, nil)
		})

		client.On("CreateRecord",
			mock.Anything, // ctx is irrelevant for this test
			mock.AnythingOfType("*dns.RecordBody"),
			"exampleterraform.io",
			mock.Anything,
		).Return(nil)

		client.On("ParseRData",
			mock.Anything,
			"A",
			[]string{"10.0.0.2", "10.0.0.3"},
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
						Config: loadFixtureString("testdata/TestResDnsRecord/create_basic.tf"),
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

	t.Run("create soa record", func(t *testing.T) {
		client := &mockdns{}

		count := 0

		client.On("GetRecord",
			mock.Anything, // ctx is irrelevant for this test
			"exampleterraform.io",
			"@",
			"SOA",
		).Return(nil, &dns.Error{
			StatusCode: http.StatusNotFound,
		}).Twice().Run(func(mock.Arguments) {
			if count < 1 {
				count++
				return
			}
			client.On("GetRecord",
				mock.Anything, // ctx is irrelevant for this test
				"exampleterraform.io",
				"@",
				"SOA",
			).Return(&dns.RecordBody{
				RecordType: "SOA",
				Name:       "exampleterraform.io",
				Target:     []string{"ns1.exampleterraform.io root@exampleterraform.io 123456789 3600 600 3600 3600"},
				TTL:        300,
			}, nil)

			client.On("ProcessRdata",
				mock.Anything, // ctx is irrelevant for this test
				[]string{"ns1.exampleterraform.io root@exampleterraform.io 123456789 3600 600 3600 3600"},
				"SOA",
			).Return([]string{"ns1.exampleterraform.io root@exampleterraform.io 123456789 3600 600 3600 3600"}, nil)
		})

		client.On("CreateRecord",
			mock.Anything, // ctx is irrelevant for this test
			mock.AnythingOfType("*dns.RecordBody"),
			"exampleterraform.io",
			mock.Anything,
		).Return(nil)

		client.On("ParseRData",
			mock.Anything,
			"SOA",
			[]string{"ns1.exampleterraform.io root@exampleterraform.io 123456789 3600 600 3600 3600"},
		).Return(parseRData(context.Background(), "SOA", []string{
			"ns1.exampleterraform.io root@exampleterraform.io 123456789 3600 600 3600 3600",
		}))

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
						Config: loadFixtureString("testdata/TestResDnsRecord/create_basic_soa.tf"),
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
		).Return(&dns.RecordBody{
			RecordType: "SOA",
			Name:       "exampleterraform.io",
			Target:     []string{"ns1.exampleterraform.io root@exampleterraform.io 123456789 3600 600 3600 3600"},
			TTL:        300,
		}, nil)

		client.On("ProcessRdata",
			mock.Anything, // ctx is irrelevant for this test
			[]string{"ns1.exampleterraform.io root@exampleterraform.io 123456789 3600 600 3600 3600"},
			"SOA",
		).Return([]string{"ns1.exampleterraform.io root@exampleterraform.io 123456789 3600 600 3600 3600"}, nil)

		client.On("UpdateRecord",
			mock.Anything, // ctx is irrelevant for this test
			mock.AnythingOfType("*dns.RecordBody"),
			"exampleterraform.io",
			mock.Anything,
		).Return(nil)

		client.On("ParseRData",
			mock.Anything,
			"SOA",
			[]string{"ns1.exampleterraform.io root@exampleterraform.io 123456789 3600 600 3600 3600"},
		).Return(parseRData(context.Background(), "SOA", []string{
			"ns1.exampleterraform.io root@exampleterraform.io 123456789 3600 600 3600 3600",
		}))

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
						Config: loadFixtureString("testdata/TestResDnsRecord/create_basic_soa.tf"),
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
