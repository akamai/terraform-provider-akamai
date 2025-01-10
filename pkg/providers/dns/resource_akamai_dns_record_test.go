package dns

import (
	"context"
	"fmt"
	"net/http"
	"regexp"
	"strings"
	"testing"

	"github.com/akamai/AkamaiOPEN-edgegrid-golang/v9/pkg/dns"
	"github.com/akamai/AkamaiOPEN-edgegrid-golang/v9/pkg/session"
	"github.com/akamai/terraform-provider-akamai/v6/pkg/common/testutils"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/schema"
	"github.com/hashicorp/terraform-plugin-testing/helper/resource"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/mock"
	"github.com/stretchr/testify/require"
)

func TestResDnsRecord(t *testing.T) {
	dnsClient := dns.Client(session.Must(session.New()))

	notFound := &dns.Error{
		StatusCode: http.StatusNotFound,
	}

	// This test performs a full life-cycle (CRUD) test
	t.Run("lifecycle test", func(t *testing.T) {
		client := &dns.Mock{}

		// read
		client.On("GetRecord",
			testutils.MockContext,
			mock.AnythingOfType("dns.GetRecordRequest"),
		).Return(nil, notFound).Once()

		// create
		client.On("CreateRecord",
			testutils.MockContext,
			mock.AnythingOfType("dns.CreateRecordRequest"),
		).Return(nil).Once()

		// read
		client.On("GetRecord",
			testutils.MockContext,
			mock.AnythingOfType("dns.GetRecordRequest"),
		).Return(&dns.GetRecordResponse{
			Name:       "",
			RecordType: "",
			TTL:        0,
			Active:     false,
			Target:     nil,
		}, nil).Once()

		retCreate := dnsClient.ParseRData(context.Background(), "A", []string{"10.0.0.2", "10.0.0.3"})

		client.On("ParseRData",
			testutils.MockContext,
			mock.AnythingOfType("string"),
			mock.AnythingOfType("[]string"),
		).Return(retCreate).Times(3)

		client.On("ProcessRdata",
			testutils.MockContext,
			mock.AnythingOfType("[]string"),
			mock.AnythingOfType("string"),
		).Return([]string{"A"}, nil).Times(4)

		client.On("GetRecord",
			testutils.MockContext,
			mock.AnythingOfType("dns.GetRecordRequest"),
		).Return(&dns.GetRecordResponse{
			Name:       "",
			RecordType: "",
			TTL:        0,
			Active:     false,
			Target:     nil,
		}, nil).Times(3)

		// update
		client.On("UpdateRecord",
			testutils.MockContext,
			mock.AnythingOfType("dns.UpdateRecordRequest"),
		).Return(nil).Once()

		// read
		client.On("GetRecord",
			testutils.MockContext,
			mock.AnythingOfType("dns.GetRecordRequest"),
		).Return(&dns.GetRecordResponse{
			Name:       "",
			RecordType: "",
			TTL:        0,
			Active:     false,
			Target:     nil,
		}, nil).Times(2)

		retUpdate := dnsClient.ParseRData(context.Background(), "A", []string{"10.0.0.4", "10.0.0.5"})

		client.On("ParseRData",
			testutils.MockContext,
			mock.AnythingOfType("string"),
			mock.AnythingOfType("[]string"),
		).Return(retUpdate).Times(2)

		client.On("ProcessRdata",
			testutils.MockContext,
			mock.AnythingOfType("[]string"),
			mock.AnythingOfType("string"),
		).Return([]string{"A"}, nil).Times(2)

		// delete
		client.On("DeleteRecord",
			testutils.MockContext,
			mock.AnythingOfType("dns.DeleteRecordRequest"),
		).Return(nil).Once()

		dataSourceName := "akamai_dns_record.a_record"

		useClient(client, func() {
			resource.UnitTest(t, resource.TestCase{
				ProtoV6ProviderFactories: testutils.NewProtoV6ProviderFactory(NewSubprovider()),
				Steps: []resource.TestStep{
					{
						Config: testutils.LoadFixtureString(t, "testdata/TestResDnsRecord/create_basic.tf"),
						Check: resource.ComposeTestCheckFunc(
							resource.TestCheckResourceAttr(dataSourceName, "recordtype", "A"),
						),
					},
					{
						Config: testutils.LoadFixtureString(t, "testdata/TestResDnsRecord/update_basic.tf"),
						Check: resource.ComposeTestCheckFunc(
							resource.TestCheckResourceAttr(dataSourceName, "recordtype", "A"),
						),
					},
				},
			})
		})

		client.AssertExpectations(t)
	})

	t.Run("TXT record test", func(t *testing.T) {
		client := &dns.Mock{}

		target1 := "\"Hel\\\\lo\\\"world\""
		target2 := "\"extralongtargetwhichis\" \"intwoseparateparts\""

		client.On("GetRecord",
			testutils.MockContext,
			mock.AnythingOfType("dns.GetRecordRequest"),
		).Return(nil, notFound).Once()

		client.On("CreateRecord",
			testutils.MockContext,
			mock.AnythingOfType("dns.CreateRecordRequest"),
		).Return(nil)

		client.On("GetRecord",
			testutils.MockContext,
			mock.AnythingOfType("dns.GetRecordRequest"),
		).Return(&dns.GetRecordResponse{
			Name:       "exampleterraform.io",
			RecordType: "TXT",
			TTL:        300,
			Active:     false,
			Target:     []string{target1, target2},
		}, nil).Once()

		client.On("ParseRData",
			testutils.MockContext,
			"TXT",
			[]string{target1, target2},
		).Return(map[string]interface{}{
			"target": []string{target1, target2},
		}).Times(2)

		client.On("ProcessRdata",
			testutils.MockContext,
			[]string{target1, target2},
			"TXT",
		).Return([]string{target1, target2}).Times(2)

		client.On("GetRecord",
			testutils.MockContext,
			mock.AnythingOfType("dns.GetRecordRequest"),
		).Return(&dns.GetRecordResponse{
			Name:       "exampleterraform.io",
			RecordType: "TXT",
			TTL:        300,
			Active:     false,
			Target:     []string{target1, target2},
		}, nil).Once()

		client.On("DeleteRecord",
			testutils.MockContext,
			mock.AnythingOfType("dns.DeleteRecordRequest"),
		).Return(nil)

		resourceName := "akamai_dns_record.txt_record"

		useClient(client, func() {
			resource.UnitTest(t, resource.TestCase{
				ProtoV6ProviderFactories: testutils.NewProtoV6ProviderFactory(NewSubprovider()),
				Steps: []resource.TestStep{
					{
						Config: testutils.LoadFixtureString(t, "testdata/TestResDnsRecord/create_basic_txt.tf"),
						Check: resource.ComposeTestCheckFunc(
							resource.TestCheckResourceAttr(resourceName, "recordtype", "TXT"),
							resource.TestCheckResourceAttr(resourceName, "target.#", "2"),
							resource.TestCheckResourceAttr(resourceName, "target.0", "Hel\\lo\"world"),
							resource.TestCheckResourceAttr(resourceName, "target.1", "\"extralongtargetwhichis\" \"intwoseparateparts\""),
						),
					},
				},
			})
		})

		client.AssertExpectations(t)
	})

	t.Run("TXT record test - keep order of elements in 255 characters long targets", func(t *testing.T) {

		target1 := strings.Repeat("Z", 255)
		target2 := strings.Repeat("A", 255)
		target3 := strings.Repeat("K", 255)

		normalizedTarget1 := fmt.Sprintf("%q", target1)
		normalizedTarget2 := fmt.Sprintf("%q", target2)
		normalizedTarget3 := fmt.Sprintf("%q", target3)

		client := &dns.Mock{}

		client.On("GetRecord",
			testutils.MockContext,
			dns.GetRecordRequest{Zone: "exampleterraform.io", Name: "exampleterraform.io", RecordType: "TXT"},
		).Return(nil, notFound).Once()

		client.On("CreateRecord",
			testutils.MockContext,
			dns.CreateRecordRequest{
				Record: &dns.RecordBody{
					Name:       "exampleterraform.io",
					RecordType: "TXT",
					TTL:        300,
					Active:     false,
					Target:     []string{normalizedTarget1, normalizedTarget2, normalizedTarget3},
				},
				Zone:    "exampleterraform.io",
				RecLock: []bool{false},
			},
		).Return(nil)

		client.On("GetRecord",
			testutils.MockContext,
			dns.GetRecordRequest{Zone: "exampleterraform.io", Name: "exampleterraform.io", RecordType: "TXT"},
		).Return(&dns.GetRecordResponse{
			Name:       "exampleterraform.io",
			RecordType: "TXT",
			TTL:        300,
			Active:     false,
			Target:     []string{normalizedTarget1, normalizedTarget2, normalizedTarget3},
		}, nil).Once()

		client.On("ParseRData",
			testutils.MockContext,
			"TXT",
			[]string{normalizedTarget1, normalizedTarget2, normalizedTarget3},
		).Return(map[string]interface{}{
			"target": []string{normalizedTarget1, normalizedTarget2, normalizedTarget3},
		}).Times(2)

		client.On("ProcessRdata",
			testutils.MockContext,
			[]string{normalizedTarget1, normalizedTarget2, normalizedTarget3},
			"TXT",
		).Return([]string{normalizedTarget1, normalizedTarget2, normalizedTarget3}).Times(2)

		client.On("GetRecord",
			testutils.MockContext,
			dns.GetRecordRequest{Zone: "exampleterraform.io", Name: "exampleterraform.io", RecordType: "TXT"},
		).Return(&dns.GetRecordResponse{
			Name:       "exampleterraform.io",
			RecordType: "TXT",
			TTL:        300,
			Active:     false,
			Target:     []string{normalizedTarget1, normalizedTarget2, normalizedTarget3},
		}, nil).Once()

		client.On("DeleteRecord",
			testutils.MockContext,
			dns.DeleteRecordRequest{Zone: "exampleterraform.io", Name: "exampleterraform.io", RecordType: "TXT", RecLock: []bool{false}},
		).Return(nil)

		resourceName := "akamai_dns_record.txt_record"

		useClient(client, func() {
			resource.UnitTest(t, resource.TestCase{
				ProtoV6ProviderFactories: testutils.NewProtoV6ProviderFactory(NewSubprovider()),
				Steps: []resource.TestStep{
					{
						Config: testutils.LoadFixtureString(t, "testdata/TestResDnsRecord/create_long_txt.tf"),
						Check: resource.ComposeTestCheckFunc(
							resource.TestCheckResourceAttr(resourceName, "recordtype", "TXT"),
							resource.TestCheckResourceAttr(resourceName, "target.#", "3"),
							resource.TestCheckResourceAttr(resourceName, "target.0", target1),
							resource.TestCheckResourceAttr(resourceName, "target.1", target2),
							resource.TestCheckResourceAttr(resourceName, "target.2", target3),
						),
					},
				},
			})
		})

		client.AssertExpectations(t)
	})

	t.Run("SRV record with default values", func(t *testing.T) {
		client := &dns.Mock{}

		targetBig := "10 60 5060 big.example.com."
		targetSmall := "10 60 5060 small.example.com."
		targetTiny := "10 60 5060 tiny.example.com."

		client.On("GetRecord",
			testutils.MockContext,
			dns.GetRecordRequest{Zone: "origin.org", Name: "origin.example.org", RecordType: "SRV"},
		).Return(nil, notFound).Once()

		client.On("CreateRecord",
			testutils.MockContext,
			dns.CreateRecordRequest{
				Record: &dns.RecordBody{
					Name:       "origin.example.org",
					RecordType: "SRV",
					TTL:        300,
					Active:     false,
					Target:     []string{targetBig, targetSmall, targetTiny},
				},
				Zone:    "origin.org",
				RecLock: []bool{false},
			},
		).Return(nil)

		client.On("GetRecord",
			testutils.MockContext,
			dns.GetRecordRequest{Zone: "origin.org", Name: "origin.example.org", RecordType: "SRV"},
		).Return(&dns.GetRecordResponse{
			Name:       "origin.example.org",
			RecordType: "SRV",
			TTL:        300,
			Active:     false,
			Target:     []string{targetBig, targetSmall, targetTiny},
		}, nil).Once()

		c := dns.Client(session.Must(session.New()))

		client.On("ParseRData",
			testutils.MockContext,
			"SRV",
			[]string{targetBig, targetSmall, targetTiny},
		).Return(
			c.ParseRData(context.Background(), "SRV", []string{targetBig, targetSmall, targetTiny}),
		).Times(2)

		client.On("ProcessRdata",
			testutils.MockContext,
			[]string{targetBig, targetSmall, targetTiny},
			"SRV",
		).Return([]string{targetBig, targetSmall, targetTiny}).Times(2)

		client.On("GetRecord",
			testutils.MockContext,
			dns.GetRecordRequest{Zone: "origin.org", Name: "origin.example.org", RecordType: "SRV"},
		).Return(&dns.GetRecordResponse{
			Name:       "origin.example.org",
			RecordType: "SRV",
			TTL:        300,
			Active:     false,
			Target:     []string{targetBig, targetSmall, targetTiny},
		}, nil).Once()

		client.On("DeleteRecord",
			testutils.MockContext,
			dns.DeleteRecordRequest{Zone: "origin.org", Name: "origin.example.org", RecordType: "SRV", RecLock: []bool{false}},
		).Return(nil)

		resourceName := "akamai_dns_record.srv_record"

		useClient(client, func() {
			resource.UnitTest(t, resource.TestCase{
				ProtoV6ProviderFactories: testutils.NewProtoV6ProviderFactory(NewSubprovider()),
				Steps: []resource.TestStep{
					{
						Config: testutils.LoadFixtureString(t, "testdata/TestResDnsRecord/srv/create_basic_srv_default.tf"),
						Check: resource.ComposeTestCheckFunc(
							resource.TestCheckResourceAttr(resourceName, "recordtype", "SRV"),
							resource.TestCheckResourceAttr(resourceName, "target.#", "3"),
							resource.TestCheckResourceAttr(resourceName, "target.0", "big.example.com."),
							resource.TestCheckResourceAttr(resourceName, "target.1", "small.example.com."),
							resource.TestCheckResourceAttr(resourceName, "target.2", "tiny.example.com."),
						),
					},
				},
			})
		})

		client.AssertExpectations(t)
	})
	t.Run("SRV record without default values", func(t *testing.T) {
		client := &dns.Mock{}

		targetBig := "10 60 5060 big.example.com."
		targetSmall := "10 40 5060 small.example.com."
		targetTiny := "20 100 5060 tiny.example.com."

		client.On("GetRecord",
			testutils.MockContext,
			dns.GetRecordRequest{Zone: "origin.org", Name: "origin.example.org", RecordType: "SRV"},
		).Return(nil, notFound).Once()

		client.On("CreateRecord",
			testutils.MockContext,
			dns.CreateRecordRequest{
				Record: &dns.RecordBody{
					Name:       "origin.example.org",
					RecordType: "SRV",
					TTL:        300,
					Active:     false,
					Target:     []string{targetBig, targetSmall, targetTiny},
				},
				Zone:    "origin.org",
				RecLock: []bool{false},
			},
		).Return(nil)

		client.On("GetRecord",
			testutils.MockContext,
			dns.GetRecordRequest{Zone: "origin.org", Name: "origin.example.org", RecordType: "SRV"},
		).Return(&dns.GetRecordResponse{
			Name:       "origin.example.org",
			RecordType: "SRV",
			TTL:        300,
			Active:     false,
			Target:     []string{targetBig, targetSmall, targetTiny},
		}, nil).Once()

		c := dns.Client(session.Must(session.New()))

		client.On("ParseRData",
			testutils.MockContext,
			"SRV",
			[]string{targetBig, targetSmall, targetTiny},
		).Return(
			c.ParseRData(context.Background(), "SRV", []string{targetBig, targetSmall, targetTiny}),
		).Times(2)

		client.On("ProcessRdata",
			testutils.MockContext,
			[]string{targetBig, targetSmall, targetTiny},
			"SRV",
		).Return([]string{targetBig, targetSmall, targetTiny}).Times(2)

		client.On("GetRecord",
			testutils.MockContext,
			dns.GetRecordRequest{Zone: "origin.org", Name: "origin.example.org", RecordType: "SRV"},
		).Return(&dns.GetRecordResponse{
			Name:       "origin.example.org",
			RecordType: "SRV",
			TTL:        300,
			Active:     false,
			Target:     []string{targetBig, targetSmall, targetTiny},
		}, nil).Once()

		client.On("DeleteRecord",
			testutils.MockContext,
			dns.DeleteRecordRequest{Zone: "origin.org", Name: "origin.example.org", RecordType: "SRV", RecLock: []bool{false}},
		).Return(nil)

		resourceName := "akamai_dns_record.srv_record"

		useClient(client, func() {
			resource.UnitTest(t, resource.TestCase{
				ProtoV6ProviderFactories: testutils.NewProtoV6ProviderFactory(NewSubprovider()),
				Steps: []resource.TestStep{
					{
						Config: testutils.LoadFixtureString(t, "testdata/TestResDnsRecord/srv/create_basic_srv_no_default.tf"),
						Check: resource.ComposeTestCheckFunc(
							resource.TestCheckResourceAttr(resourceName, "recordtype", "SRV"),
							resource.TestCheckResourceAttr(resourceName, "target.#", "3"),
							resource.TestCheckResourceAttr(resourceName, "target.0", "10 60 5060 big.example.com."),
							resource.TestCheckResourceAttr(resourceName, "target.1", "10 40 5060 small.example.com."),
							resource.TestCheckResourceAttr(resourceName, "target.2", "20 100 5060 tiny.example.com."),
						),
					},
				},
			})
		})

		client.AssertExpectations(t)
	})
	t.Run("SRV record with invalid mixed values", func(t *testing.T) {
		client := &dns.Mock{}

		useClient(client, func() {
			resource.UnitTest(t, resource.TestCase{
				ProtoV6ProviderFactories: testutils.NewProtoV6ProviderFactory(NewSubprovider()),
				Steps: []resource.TestStep{
					{
						Config:      testutils.LoadFixtureString(t, "testdata/TestResDnsRecord/srv/create_basic_srv_mix_invalid.tf"),
						ExpectError: regexp.MustCompile("target should consist of only simple or complete items"),
					},
				},
			})
		})

		client.AssertExpectations(t)
	})

	t.Run("AAAA record with valid IPv6 addresses", func(t *testing.T) {
		client := &dns.Mock{}

		target := []string{"2001:db8::68", "::ffff:192.0.2.1"}
		targetSent := []string{"0000:0000:0000:0000:0000:ffff:c000:0201", "2001:0db8:0000:0000:0000:0000:0000:0068"}
		targetReceived := []string{"2001:db8:0:0:0:0:0:68", "::ffff:192.0.2.1"}
		client.On("GetRecord",
			testutils.MockContext,
			dns.GetRecordRequest{Zone: "exampleterraform.io", Name: "exampleterraform.io", RecordType: "AAAA"},
		).Return(nil, notFound).Once()

		client.On("CreateRecord",
			testutils.MockContext,
			dns.CreateRecordRequest{
				Record: &dns.RecordBody{
					Name:       "exampleterraform.io",
					RecordType: "AAAA",
					TTL:        300,
					Active:     false,
					Target:     targetSent,
				},
				Zone:    "exampleterraform.io",
				RecLock: []bool{false},
			},
		).Return(nil)

		client.On("GetRecord",
			testutils.MockContext,
			dns.GetRecordRequest{Zone: "exampleterraform.io", Name: "exampleterraform.io", RecordType: "AAAA"},
		).Return(&dns.GetRecordResponse{
			Name:       "exampleterraform.io",
			RecordType: "AAAA",
			TTL:        300,
			Active:     false,
			Target:     targetReceived,
		}, nil).Once()

		client.On("ParseRData",
			testutils.MockContext,
			"AAAA",
			targetReceived,
		).Return(map[string]interface{}{
			"target": target,
		}).Times(2)

		client.On("ProcessRdata",
			testutils.MockContext,
			targetReceived,
			"AAAA",
		).Return(target).Times(2)

		client.On("GetRecord",
			testutils.MockContext,
			dns.GetRecordRequest{Zone: "exampleterraform.io", Name: "exampleterraform.io", RecordType: "AAAA"},
		).Return(&dns.GetRecordResponse{
			Name:       "exampleterraform.io",
			RecordType: "AAAA",
			TTL:        300,
			Active:     false,
			Target:     targetReceived,
		}, nil).Once()

		client.On("DeleteRecord",
			testutils.MockContext,
			dns.DeleteRecordRequest{Zone: "exampleterraform.io", Name: "exampleterraform.io", RecordType: "AAAA", RecLock: []bool{false}},
		).Return(nil)

		resourceName := "akamai_dns_record.aaaa_record"

		useClient(client, func() {
			resource.UnitTest(t, resource.TestCase{
				ProtoV6ProviderFactories: testutils.NewProtoV6ProviderFactory(NewSubprovider()),
				Steps: []resource.TestStep{
					{
						Config: testutils.LoadFixtureString(t, "testdata/TestResDnsRecord/aaaa/create_valid_aaaa.tf"),
						Check: resource.ComposeTestCheckFunc(
							resource.TestCheckResourceAttr(resourceName, "recordtype", "AAAA"),
							resource.TestCheckResourceAttr(resourceName, "target.#", "2"),
							resource.TestCheckResourceAttr(resourceName, "target.0", "2001:db8:0:0:0:0:0:68"),
							resource.TestCheckResourceAttr(resourceName, "target.1", "::ffff:192.0.2.1"),
						),
					},
				},
			})
		})

		client.AssertExpectations(t)
	})
	t.Run("AAAA record with invalid IPv6 address", func(t *testing.T) {
		client := &dns.Mock{}

		useClient(client, func() {
			resource.UnitTest(t, resource.TestCase{
				ProtoV6ProviderFactories: testutils.NewProtoV6ProviderFactory(NewSubprovider()),
				Steps: []resource.TestStep{
					{
						Config:      testutils.LoadFixtureString(t, "testdata/TestResDnsRecord/aaaa/create_invalid_aaaa.tf"),
						ExpectError: regexp.MustCompile("target '1111:2222:3333:4444:55555:6666:7777:8888' is not a valid address"),
					},
				},
			})
		})

		client.AssertExpectations(t)
	})
	t.Run("AAAA record with IP4 address - invalid", func(t *testing.T) {
		client := &dns.Mock{}

		useClient(client, func() {
			resource.UnitTest(t, resource.TestCase{
				ProtoV6ProviderFactories: testutils.NewProtoV6ProviderFactory(NewSubprovider()),
				Steps: []resource.TestStep{
					{
						Config:      testutils.LoadFixtureString(t, "testdata/TestResDnsRecord/aaaa/create_invalid_ipv4.tf"),
						ExpectError: regexp.MustCompile("target '18.244.102.124' is not a valid IPv6 or IPv4-mapped IPv6 address"),
					},
				},
			})
		})

		client.AssertExpectations(t)
	})

}

func TestMXRecord(t *testing.T) {
	notFound := &dns.Error{
		StatusCode: http.StatusNotFound,
	}
	dnsClient := dns.Client(session.Must(session.New()))
	name, zone, mx := "exampleterraform.io", "exampleterraform.io", "MX"
	getRecordRequest := dns.GetRecordRequest{
		Zone:       zone,
		Name:       name,
		RecordType: mx}

	mockCreate := func(d *dns.Mock, _ dns.DNS, createdRecord *dns.RecordBody) {
		d.On("GetRecord", testutils.MockContext, getRecordRequest).
			Return(nil, notFound).Twice()
		d.On("CreateRecord", testutils.MockContext, dns.CreateRecordRequest{
			Record:  createdRecord,
			Zone:    zone,
			RecLock: []bool{false}}).
			Return(nil).Once()
	}
	mockRead := func(d *dns.Mock, realClient dns.DNS, createdRecord *dns.RecordBody) {
		response := (*dns.GetRecordResponse)(createdRecord)
		d.On("GetRecord", testutils.MockContext, getRecordRequest).
			Return(response, nil).Once()
		d.On("ProcessRdata", testutils.MockContext, createdRecord.Target, mx).
			Return(realClient.ProcessRdata(context.Background(), createdRecord.Target, mx)).Once()
		d.On("GetRecord", testutils.MockContext, getRecordRequest).
			Return(response, nil).Once()
		d.On("ParseRData", testutils.MockContext, mx, createdRecord.Target).
			Return(realClient.ParseRData(context.Background(), mx, createdRecord.Target)).Once()
		d.On("ProcessRdata", testutils.MockContext, createdRecord.Target, mx).
			Return(realClient.ProcessRdata(context.Background(), createdRecord.Target, mx)).Once()
	}
	mockUpdate := func(d *dns.Mock, realClient dns.DNS, previousRecord *dns.RecordBody, updatedRecord *dns.RecordBody) {
		response := (*dns.GetRecordResponse)(previousRecord)
		d.On("GetRecord", testutils.MockContext, getRecordRequest).
			Return(response, nil).Once()
		d.On("ProcessRdata", testutils.MockContext, previousRecord.Target, mx).
			Return(realClient.ProcessRdata(context.Background(), previousRecord.Target, mx)).Once()
		d.On("GetRecord", testutils.MockContext, getRecordRequest).
			Return(response, nil).Once()
		d.On("ProcessRdata", testutils.MockContext, previousRecord.Target, mx).
			Return(realClient.ProcessRdata(context.Background(), previousRecord.Target, mx)).Once()
		d.On("UpdateRecord", testutils.MockContext, dns.UpdateRecordRequest{
			Record:  updatedRecord,
			Zone:    zone,
			RecLock: []bool{false}}).
			Return(nil)
	}
	mockDelete := func(d *dns.Mock, createdRecord *dns.RecordBody) {
		d.On("DeleteRecord", testutils.MockContext, dns.DeleteRecordRequest{
			Zone:       zone,
			Name:       createdRecord.Name,
			RecordType: createdRecord.RecordType,
			RecLock:    []bool{false}}).
			Return(nil)
	}

	defaultInit := func(d *dns.Mock, createTargets, updateTargets, deleteTargets []string) {
		createdRecord := &dns.RecordBody{
			Name:       name,
			RecordType: mx,
			TTL:        300,
			Target:     createTargets,
		}
		mockCreate(d, dnsClient, createdRecord)
		//Read after create
		mockRead(d, dnsClient, createdRecord)
		//Second read
		mockRead(d, dnsClient, createdRecord)
		//Read before update
		mockRead(d, dnsClient, createdRecord)

		updatedRecord := &dns.RecordBody{
			Name:       name,
			RecordType: mx,
			TTL:        300,
			Target:     updateTargets,
		}
		mockUpdate(d, dnsClient, createdRecord, updatedRecord)
		//Read after update
		mockRead(d, dnsClient, updatedRecord)
		//Second read
		mockRead(d, dnsClient, updatedRecord)

		deleteRecord := &dns.RecordBody{
			Name:       name,
			RecordType: mx,
			TTL:        300,
			Target:     deleteTargets,
		}
		mockDelete(d, deleteRecord)
	}
	tests := map[string]struct {
		createTargets []string
		updateTargets []string
		deleteTargets []string
		steps         []resource.TestStep
	}{
		"priorities in targets": {
			createTargets: []string{
				"5 mx1.test.com.",
				"10 mx2.test.com.",
				"15 mx3.test.com.",
			},
			updateTargets: []string{
				"5 mx1.test.com.",
				"10 mx2.test.com.",
				"15 mx3.test.com.",
				"20 mx4.test.com.",
			},
			deleteTargets: []string{
				"5 mx1.test.com.",
				"10 mx2.test.com.",
				"15 mx3.test.com.",
				"20 mx4.test.com.",
			},
			steps: []resource.TestStep{
				{
					Config: testutils.LoadFixtureString(t, "testdata/TestResDnsRecordMX/create_target.tf"),
					Check: resource.ComposeTestCheckFunc(
						resource.TestCheckResourceAttr("akamai_dns_record.record", "name", "exampleterraform.io"),
						resource.TestCheckResourceAttr("akamai_dns_record.record", "zone", "exampleterraform.io"),
						resource.TestCheckResourceAttr("akamai_dns_record.record", "recordtype", "MX"),
						resource.TestCheckResourceAttr("akamai_dns_record.record", "ttl", "300"),
						resource.TestCheckResourceAttr("akamai_dns_record.record", "target.#", "3"),
						resource.TestCheckResourceAttr("akamai_dns_record.record", "target.0", "5 mx1.test.com."),
						resource.TestCheckResourceAttr("akamai_dns_record.record", "target.1", "10 mx2.test.com."),
						resource.TestCheckResourceAttr("akamai_dns_record.record", "target.2", "15 mx3.test.com."),
					),
				},
				{
					Config: testutils.LoadFixtureString(t, "testdata/TestResDnsRecordMX/update_target.tf"),
					Check: resource.ComposeTestCheckFunc(
						resource.TestCheckResourceAttr("akamai_dns_record.record", "name", "exampleterraform.io"),
						resource.TestCheckResourceAttr("akamai_dns_record.record", "zone", "exampleterraform.io"),
						resource.TestCheckResourceAttr("akamai_dns_record.record", "recordtype", "MX"),
						resource.TestCheckResourceAttr("akamai_dns_record.record", "ttl", "300"),
						resource.TestCheckResourceAttr("akamai_dns_record.record", "target.#", "4"),
						resource.TestCheckResourceAttr("akamai_dns_record.record", "target.0", "5 mx1.test.com."),
						resource.TestCheckResourceAttr("akamai_dns_record.record", "target.1", "10 mx2.test.com."),
						resource.TestCheckResourceAttr("akamai_dns_record.record", "target.2", "15 mx3.test.com."),
						resource.TestCheckResourceAttr("akamai_dns_record.record", "target.3", "20 mx4.test.com."),
					),
				},
			},
		},
		"priorities in priority": {
			createTargets: []string{
				"3 mx1.test.com.",
				"5 mx2.test.com.",
				"7 mx3.test.com.",
			},
			updateTargets: []string{
				"4 mx1.test.com.",
				"6 mx2.test.com.",
				"8 mx3.test.com.",
			},
			deleteTargets: []string{
				"mx1.test.com.",
				"mx2.test.com.",
				"mx3.test.com.",
			},
			steps: []resource.TestStep{
				{
					Config: testutils.LoadFixtureString(t, "testdata/TestResDnsRecordMX/create_priority.tf"),
					Check: resource.ComposeTestCheckFunc(
						resource.TestCheckResourceAttr("akamai_dns_record.record", "name", "exampleterraform.io"),
						resource.TestCheckResourceAttr("akamai_dns_record.record", "zone", "exampleterraform.io"),
						resource.TestCheckResourceAttr("akamai_dns_record.record", "recordtype", "MX"),
						resource.TestCheckResourceAttr("akamai_dns_record.record", "ttl", "300"),
						resource.TestCheckResourceAttr("akamai_dns_record.record", "target.#", "3"),
						resource.TestCheckResourceAttr("akamai_dns_record.record", "target.0", "mx1.test.com."),
						resource.TestCheckResourceAttr("akamai_dns_record.record", "target.1", "mx2.test.com."),
						resource.TestCheckResourceAttr("akamai_dns_record.record", "target.2", "mx3.test.com."),
						resource.TestCheckResourceAttr("akamai_dns_record.record", "priority", "3"),
						resource.TestCheckResourceAttr("akamai_dns_record.record", "priority_increment", "2"),
					),
				},
				{
					Config: testutils.LoadFixtureString(t, "testdata/TestResDnsRecordMX/update_priority.tf"),
					Check: resource.ComposeTestCheckFunc(
						resource.TestCheckResourceAttr("akamai_dns_record.record", "name", "exampleterraform.io"),
						resource.TestCheckResourceAttr("akamai_dns_record.record", "zone", "exampleterraform.io"),
						resource.TestCheckResourceAttr("akamai_dns_record.record", "recordtype", "MX"),
						resource.TestCheckResourceAttr("akamai_dns_record.record", "ttl", "300"),
						resource.TestCheckResourceAttr("akamai_dns_record.record", "target.#", "3"),
						resource.TestCheckResourceAttr("akamai_dns_record.record", "target.0", "mx1.test.com."),
						resource.TestCheckResourceAttr("akamai_dns_record.record", "target.1", "mx2.test.com."),
						resource.TestCheckResourceAttr("akamai_dns_record.record", "target.2", "mx3.test.com."),
						resource.TestCheckResourceAttr("akamai_dns_record.record", "priority", "4"),
						resource.TestCheckResourceAttr("akamai_dns_record.record", "priority_increment", "2"),
					),
				},
			},
		},
		"priorities in priority - update increment": {
			createTargets: []string{
				"3 mx1.test.com.",
				"5 mx2.test.com.",
				"7 mx3.test.com.",
			},
			updateTargets: []string{
				"3 mx1.test.com.",
				"6 mx2.test.com.",
				"9 mx3.test.com.",
			},
			deleteTargets: []string{
				"mx1.test.com.",
				"mx2.test.com.",
				"mx3.test.com.",
			},
			steps: []resource.TestStep{
				{
					Config: testutils.LoadFixtureString(t, "testdata/TestResDnsRecordMX/create_priority.tf"),
					Check: resource.ComposeTestCheckFunc(
						resource.TestCheckResourceAttr("akamai_dns_record.record", "name", "exampleterraform.io"),
						resource.TestCheckResourceAttr("akamai_dns_record.record", "zone", "exampleterraform.io"),
						resource.TestCheckResourceAttr("akamai_dns_record.record", "recordtype", "MX"),
						resource.TestCheckResourceAttr("akamai_dns_record.record", "ttl", "300"),
						resource.TestCheckResourceAttr("akamai_dns_record.record", "target.#", "3"),
						resource.TestCheckResourceAttr("akamai_dns_record.record", "target.0", "mx1.test.com."),
						resource.TestCheckResourceAttr("akamai_dns_record.record", "target.1", "mx2.test.com."),
						resource.TestCheckResourceAttr("akamai_dns_record.record", "target.2", "mx3.test.com."),
						resource.TestCheckResourceAttr("akamai_dns_record.record", "priority", "3"),
						resource.TestCheckResourceAttr("akamai_dns_record.record", "priority_increment", "2"),
					),
				},
				{
					Config: testutils.LoadFixtureString(t, "testdata/TestResDnsRecordMX/update_increment.tf"),
					Check: resource.ComposeTestCheckFunc(
						resource.TestCheckResourceAttr("akamai_dns_record.record", "name", "exampleterraform.io"),
						resource.TestCheckResourceAttr("akamai_dns_record.record", "zone", "exampleterraform.io"),
						resource.TestCheckResourceAttr("akamai_dns_record.record", "recordtype", "MX"),
						resource.TestCheckResourceAttr("akamai_dns_record.record", "ttl", "300"),
						resource.TestCheckResourceAttr("akamai_dns_record.record", "target.#", "3"),
						resource.TestCheckResourceAttr("akamai_dns_record.record", "target.0", "mx1.test.com."),
						resource.TestCheckResourceAttr("akamai_dns_record.record", "target.1", "mx2.test.com."),
						resource.TestCheckResourceAttr("akamai_dns_record.record", "target.2", "mx3.test.com."),
						resource.TestCheckResourceAttr("akamai_dns_record.record", "priority", "3"),
						resource.TestCheckResourceAttr("akamai_dns_record.record", "priority_increment", "3"),
					),
				},
			},
		},
		"change from targets to priority": {
			createTargets: []string{
				"5 mx1.test.com.",
				"10 mx2.test.com.",
				"15 mx3.test.com.",
			},
			updateTargets: []string{
				"4 mx1.test.com.",
				"6 mx2.test.com.",
				"8 mx3.test.com.",
			},
			deleteTargets: []string{
				"mx1.test.com.",
				"mx2.test.com.",
				"mx3.test.com.",
			},
			steps: []resource.TestStep{
				{
					Config: testutils.LoadFixtureString(t, "testdata/TestResDnsRecordMX/create_target.tf"),
					Check: resource.ComposeTestCheckFunc(
						resource.TestCheckResourceAttr("akamai_dns_record.record", "name", "exampleterraform.io"),
						resource.TestCheckResourceAttr("akamai_dns_record.record", "zone", "exampleterraform.io"),
						resource.TestCheckResourceAttr("akamai_dns_record.record", "recordtype", "MX"),
						resource.TestCheckResourceAttr("akamai_dns_record.record", "ttl", "300"),
						resource.TestCheckResourceAttr("akamai_dns_record.record", "target.#", "3"),
						resource.TestCheckResourceAttr("akamai_dns_record.record", "target.0", "5 mx1.test.com."),
						resource.TestCheckResourceAttr("akamai_dns_record.record", "target.1", "10 mx2.test.com."),
						resource.TestCheckResourceAttr("akamai_dns_record.record", "target.2", "15 mx3.test.com."),
					),
				},
				{
					Config: testutils.LoadFixtureString(t, "testdata/TestResDnsRecordMX/update_priority.tf"),
					Check: resource.ComposeTestCheckFunc(
						resource.TestCheckResourceAttr("akamai_dns_record.record", "name", "exampleterraform.io"),
						resource.TestCheckResourceAttr("akamai_dns_record.record", "zone", "exampleterraform.io"),
						resource.TestCheckResourceAttr("akamai_dns_record.record", "recordtype", "MX"),
						resource.TestCheckResourceAttr("akamai_dns_record.record", "ttl", "300"),
						resource.TestCheckResourceAttr("akamai_dns_record.record", "target.#", "3"),
						resource.TestCheckResourceAttr("akamai_dns_record.record", "target.0", "mx1.test.com."),
						resource.TestCheckResourceAttr("akamai_dns_record.record", "target.1", "mx2.test.com."),
						resource.TestCheckResourceAttr("akamai_dns_record.record", "target.2", "mx3.test.com."),
						resource.TestCheckResourceAttr("akamai_dns_record.record", "priority", "4"),
						resource.TestCheckResourceAttr("akamai_dns_record.record", "priority_increment", "2"),
					),
				},
			},
		},
		"change from priority to target": {
			createTargets: []string{
				"3 mx1.test.com.",
				"5 mx2.test.com.",
				"7 mx3.test.com.",
			},
			updateTargets: []string{
				"5 mx1.test.com.",
				"10 mx2.test.com.",
				"15 mx3.test.com.",
				"20 mx4.test.com.",
			},
			deleteTargets: []string{
				"5 mx1.test.com.",
				"10 mx2.test.com.",
				"15 mx3.test.com.",
				"20 mx4.test.com.",
			},
			steps: []resource.TestStep{
				{
					Config: testutils.LoadFixtureString(t, "testdata/TestResDnsRecordMX/create_priority.tf"),
					Check: resource.ComposeTestCheckFunc(
						resource.TestCheckResourceAttr("akamai_dns_record.record", "name", "exampleterraform.io"),
						resource.TestCheckResourceAttr("akamai_dns_record.record", "zone", "exampleterraform.io"),
						resource.TestCheckResourceAttr("akamai_dns_record.record", "recordtype", "MX"),
						resource.TestCheckResourceAttr("akamai_dns_record.record", "ttl", "300"),
						resource.TestCheckResourceAttr("akamai_dns_record.record", "target.#", "3"),
						resource.TestCheckResourceAttr("akamai_dns_record.record", "target.0", "mx1.test.com."),
						resource.TestCheckResourceAttr("akamai_dns_record.record", "target.1", "mx2.test.com."),
						resource.TestCheckResourceAttr("akamai_dns_record.record", "target.2", "mx3.test.com."),
						resource.TestCheckResourceAttr("akamai_dns_record.record", "priority", "3"),
						resource.TestCheckResourceAttr("akamai_dns_record.record", "priority_increment", "2"),
					),
				},
				{
					Config: testutils.LoadFixtureString(t, "testdata/TestResDnsRecordMX/update_target.tf"),
					Check: resource.ComposeTestCheckFunc(
						resource.TestCheckResourceAttr("akamai_dns_record.record", "name", "exampleterraform.io"),
						resource.TestCheckResourceAttr("akamai_dns_record.record", "zone", "exampleterraform.io"),
						resource.TestCheckResourceAttr("akamai_dns_record.record", "recordtype", "MX"),
						resource.TestCheckResourceAttr("akamai_dns_record.record", "ttl", "300"),
						resource.TestCheckResourceAttr("akamai_dns_record.record", "target.#", "4"),
						resource.TestCheckResourceAttr("akamai_dns_record.record", "target.0", "5 mx1.test.com."),
						resource.TestCheckResourceAttr("akamai_dns_record.record", "target.1", "10 mx2.test.com."),
						resource.TestCheckResourceAttr("akamai_dns_record.record", "target.2", "15 mx3.test.com."),
						resource.TestCheckResourceAttr("akamai_dns_record.record", "target.3", "20 mx4.test.com."),
					),
				},
			},
		},
	}

	for name, test := range tests {
		t.Run(name, func(t *testing.T) {
			client := new(dns.Mock)
			defaultInit(client, test.createTargets, test.updateTargets, test.deleteTargets)
			useClient(client, func() {
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

func TestTargetDiffSuppress(t *testing.T) {
	t.Run("target is computed and recordType is AAAA", func(t *testing.T) {
		config := schema.TestResourceDataRaw(t, getResourceDNSRecordSchema(), map[string]interface{}{"recordtype": "AAAA"})
		assert.False(t, dnsRecordTargetSuppress("target.#", "0", "", config))
	})
}

func TestResolveTxtRecordTargets(t *testing.T) {
	denormalized := []string{"onetwo", "\"one\" \"two\""}
	normalized := []string{"\"onetwo\"", "\"one\" \"two\"", "\"one\" \"two\""}
	expected := []string{"onetwo", "\"one\" \"two\"", "\"one\" \"two\""}

	res, err := resolveTxtRecordTargets(denormalized, normalized)
	require.NoError(t, err)

	assert.Equal(t, expected, res)
}

func TestResolveTargets(t *testing.T) {
	normalize := func(value string) (string, error) {
		if value == "error" {
			return "", fmt.Errorf("oops")
		}
		return strings.ToLower(value), nil
	}

	tests := map[string]struct {
		denormalized []string
		normalized   []string
		expected     []string
		withError    bool
	}{
		"replaces equal targets": {
			denormalized: []string{"a", "B", "C"},
			normalized:   []string{"a", "b", "c", "d"},
			expected:     []string{"a", "B", "C", "d"},
		},
		"preserves additional normalized targets": {
			denormalized: []string{"a", "b"},
			normalized:   []string{"a", "b", "c", "d"},
			expected:     []string{"a", "b", "c", "d"},
		},
		"does not append additional denormalized targets": {
			denormalized: []string{"a", "b", "C", "D"},
			normalized:   []string{"a", "b"},
			expected:     []string{"a", "b"},
		},
		"preserves denormalized targets when elements shift with normalized drift": {
			denormalized: []string{"a", "B", "C"},
			normalized:   []string{"a", "b", "bb", "c"},
			expected:     []string{"a", "B", "bb", "C"},
		},
		"preserves denormalized targets when order changes": {
			denormalized: []string{"a", "B", "C", "d"},
			normalized:   []string{"d", "c", "b", "a"},
			expected:     []string{"d", "C", "B", "a"},
		},
		"returns error when normalization failed": {
			denormalized: []string{"error"},
			normalized:   []string{"a"},
			withError:    true,
		},
	}

	for name, tc := range tests {
		t.Run(name, func(t *testing.T) {
			res, err := resolveTargets(tc.denormalized, tc.normalized, normalize)
			if tc.withError {
				assert.Error(t, err)
				return
			}

			require.NoError(t, err)
			assert.Equal(t, tc.expected, res)
		})
	}
}
