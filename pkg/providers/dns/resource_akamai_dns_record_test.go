package dns

import (
	"context"
	"fmt"
	"net/http"
	"regexp"
	"strings"
	"testing"

	"github.com/akamai/AkamaiOPEN-edgegrid-golang/v13/pkg/dns"
	"github.com/akamai/AkamaiOPEN-edgegrid-golang/v13/pkg/ptr"
	"github.com/akamai/AkamaiOPEN-edgegrid-golang/v13/pkg/session"
	"github.com/akamai/terraform-provider-akamai/v10/internal/edgegrid"
	"github.com/akamai/terraform-provider-akamai/v10/pkg/common/test"
	"github.com/akamai/terraform-provider-akamai/v10/pkg/common/testutils"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/schema"
	"github.com/hashicorp/terraform-plugin-testing/helper/resource"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/mock"
	"github.com/stretchr/testify/require"
)

func TestResDnsRecord(t *testing.T) {
	t.Parallel()
	dnsClient := dns.Client(session.Must(session.New()))

	notFound := &dns.Error{
		StatusCode: http.StatusNotFound,
	}

	// This test performs a full life-cycle (CRUD) test
	t.Run("lifecycle test", func(t *testing.T) {
		t.Parallel()
		client := edgegrid.NewTestClient()

		// read
		client.DNS.On("GetRecord",
			testutils.MockContext,
			mock.AnythingOfType("dns.GetRecordRequest"),
		).Return(nil, notFound).Once()

		// create
		client.DNS.On("CreateRecord",
			testutils.MockContext,
			mock.AnythingOfType("dns.CreateRecordRequest"),
		).Return(nil).Once()

		// read
		client.DNS.On("GetRecord",
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

		client.DNS.On("ParseRData",
			testutils.MockContext,
			mock.AnythingOfType("string"),
			mock.AnythingOfType("[]string"),
		).Return(retCreate).Times(3)

		client.DNS.On("ProcessRdata",
			testutils.MockContext,
			mock.AnythingOfType("[]string"),
			mock.AnythingOfType("string"),
		).Return([]string{"A"}, nil).Times(4)

		client.DNS.On("GetRecord",
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
		client.DNS.On("UpdateRecord",
			testutils.MockContext,
			mock.AnythingOfType("dns.UpdateRecordRequest"),
		).Return(nil).Once()

		// read
		client.DNS.On("GetRecord",
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

		client.DNS.On("ParseRData",
			testutils.MockContext,
			mock.AnythingOfType("string"),
			mock.AnythingOfType("[]string"),
		).Return(retUpdate).Times(2)

		client.DNS.On("ProcessRdata",
			testutils.MockContext,
			mock.AnythingOfType("[]string"),
			mock.AnythingOfType("string"),
		).Return([]string{"A"}, nil).Times(2)

		// delete
		client.DNS.On("DeleteRecord",
			testutils.MockContext,
			mock.AnythingOfType("dns.DeleteRecordRequest"),
		).Return(nil).Once()

		dataSourceName := "akamai_dns_record.a_record"

		resource.UnitTest(t, resource.TestCase{
			ProtoV6ProviderFactories: testutils.NewTestProtoV6SDKProviderFactory(client, newSubproviderWithConfig(testSubproviderConfig())),
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

		client.DNS.AssertExpectations(t)
	})

	t.Run("TXT record test", func(t *testing.T) {
		t.Parallel()
		client := edgegrid.NewTestClient()

		target1 := "\"Hel\\\\lo\\\"world\""
		target2 := "\"extralongtargetwhichis\" \"intwoseparateparts\""

		client.DNS.On("GetRecord",
			testutils.MockContext,
			mock.AnythingOfType("dns.GetRecordRequest"),
		).Return(nil, notFound).Once()

		client.DNS.On("CreateRecord",
			testutils.MockContext,
			mock.AnythingOfType("dns.CreateRecordRequest"),
		).Return(nil)

		client.DNS.On("GetRecord",
			testutils.MockContext,
			mock.AnythingOfType("dns.GetRecordRequest"),
		).Return(&dns.GetRecordResponse{
			Name:       "exampleterraform.io",
			RecordType: "TXT",
			TTL:        300,
			Active:     false,
			Target:     []string{target1, target2},
		}, nil).Once()

		client.DNS.On("ParseRData",
			testutils.MockContext,
			"TXT",
			[]string{target1, target2},
		).Return(map[string]interface{}{
			"target": []string{target1, target2},
		}).Times(2)

		client.DNS.On("ProcessRdata",
			testutils.MockContext,
			[]string{target1, target2},
			"TXT",
		).Return([]string{target1, target2}).Times(2)

		client.DNS.On("GetRecord",
			testutils.MockContext,
			mock.AnythingOfType("dns.GetRecordRequest"),
		).Return(&dns.GetRecordResponse{
			Name:       "exampleterraform.io",
			RecordType: "TXT",
			TTL:        300,
			Active:     false,
			Target:     []string{target1, target2},
		}, nil).Once()

		client.DNS.On("DeleteRecord",
			testutils.MockContext,
			mock.AnythingOfType("dns.DeleteRecordRequest"),
		).Return(nil)

		resourceName := "akamai_dns_record.txt_record"

		resource.UnitTest(t, resource.TestCase{
			ProtoV6ProviderFactories: testutils.NewTestProtoV6SDKProviderFactory(client, newSubproviderWithConfig(testSubproviderConfig())),
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

		client.DNS.AssertExpectations(t)
	})

	t.Run("TXT record test - keep order of elements in 255 characters long targets", func(t *testing.T) {
		t.Parallel()

		target1 := strings.Repeat("Z", 255)
		target2 := strings.Repeat("A", 255)
		target3 := strings.Repeat("K", 255)

		normalizedTarget1 := fmt.Sprintf("%q", target1)
		normalizedTarget2 := fmt.Sprintf("%q", target2)
		normalizedTarget3 := fmt.Sprintf("%q", target3)

		client := edgegrid.NewTestClient()

		client.DNS.On("GetRecord",
			testutils.MockContext,
			dns.GetRecordRequest{Zone: "exampleterraform.io", Name: "exampleterraform.io", RecordType: "TXT"},
		).Return(nil, notFound).Once()

		client.DNS.On("CreateRecord",
			testutils.MockContext,
			dns.CreateRecordRequest{
				Record: &dns.RecordBody{
					Name:       "exampleterraform.io",
					RecordType: "TXT",
					TTL:        ptr.To(300),
					Active:     false,
					Target:     []string{normalizedTarget1, normalizedTarget2, normalizedTarget3},
				},
				Zone:    "exampleterraform.io",
				RecLock: []bool{false},
			},
		).Return(nil)

		client.DNS.On("GetRecord",
			testutils.MockContext,
			dns.GetRecordRequest{Zone: "exampleterraform.io", Name: "exampleterraform.io", RecordType: "TXT"},
		).Return(&dns.GetRecordResponse{
			Name:       "exampleterraform.io",
			RecordType: "TXT",
			TTL:        300,
			Active:     false,
			Target:     []string{normalizedTarget1, normalizedTarget2, normalizedTarget3},
		}, nil).Once()

		client.DNS.On("ParseRData",
			testutils.MockContext,
			"TXT",
			[]string{normalizedTarget1, normalizedTarget2, normalizedTarget3},
		).Return(map[string]interface{}{
			"target": []string{normalizedTarget1, normalizedTarget2, normalizedTarget3},
		}).Times(2)

		client.DNS.On("ProcessRdata",
			testutils.MockContext,
			[]string{normalizedTarget1, normalizedTarget2, normalizedTarget3},
			"TXT",
		).Return([]string{normalizedTarget1, normalizedTarget2, normalizedTarget3}).Times(2)

		client.DNS.On("GetRecord",
			testutils.MockContext,
			dns.GetRecordRequest{Zone: "exampleterraform.io", Name: "exampleterraform.io", RecordType: "TXT"},
		).Return(&dns.GetRecordResponse{
			Name:       "exampleterraform.io",
			RecordType: "TXT",
			TTL:        300,
			Active:     false,
			Target:     []string{normalizedTarget1, normalizedTarget2, normalizedTarget3},
		}, nil).Once()

		client.DNS.On("DeleteRecord",
			testutils.MockContext,
			dns.DeleteRecordRequest{Zone: "exampleterraform.io", Name: "exampleterraform.io", RecordType: "TXT", RecLock: []bool{false}},
		).Return(nil)

		resourceName := "akamai_dns_record.txt_record"

		resource.UnitTest(t, resource.TestCase{
			ProtoV6ProviderFactories: testutils.NewTestProtoV6SDKProviderFactory(client, newSubproviderWithConfig(testSubproviderConfig())),
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

		client.DNS.AssertExpectations(t)
	})

	t.Run("TXT record test - update target", func(t *testing.T) {

		target := "\"v=spf1 mx include:spf.domain.com include:spf.protection.outlook.com -all\""
		name := "infrastructure.domain.net"

		normalizedTarget := fmt.Sprintf("%q", target)

		client := &dns.Mock{}

		client.On("GetRecord",
			testutils.MockContext,
			dns.GetRecordRequest{Zone: name, Name: name, RecordType: "TXT"},
		).Return(nil, notFound).Once()

		client.On("CreateRecord",
			testutils.MockContext,
			dns.CreateRecordRequest{
				Record: &dns.RecordBody{
					Name:       name,
					RecordType: "TXT",
					TTL:        300,
					Active:     false,
					Target:     []string{normalizedTarget},
				},
				Zone:    name,
				RecLock: []bool{false},
			},
		).Return(nil)

		client.On("GetRecord",
			testutils.MockContext,
			dns.GetRecordRequest{Zone: name, Name: name, RecordType: "TXT"},
		).Return(&dns.GetRecordResponse{
			Name:       name,
			RecordType: "TXT",
			TTL:        300,
			Active:     false,
			Target:     []string{normalizedTarget},
		}, nil).Once()

		client.On("ParseRData",
			testutils.MockContext,
			"TXT",
			[]string{normalizedTarget},
		).Return(map[string]interface{}{
			"target": []string{normalizedTarget},
		}).Times(2)

		client.On("ProcessRdata",
			testutils.MockContext,
			[]string{normalizedTarget},
			"TXT",
		).Return([]string{normalizedTarget}).Times(2)

		client.On("GetRecord",
			testutils.MockContext,
			dns.GetRecordRequest{Zone: name, Name: name, RecordType: "TXT"},
		).Return(&dns.GetRecordResponse{
			Name:       name,
			RecordType: "TXT",
			TTL:        300,
			Active:     false,
			Target:     []string{normalizedTarget},
		}, nil).Once()

		client.On("DeleteRecord",
			testutils.MockContext,
			dns.DeleteRecordRequest{Zone: name, Name: name, RecordType: "TXT", RecLock: []bool{false}},
		).Return(nil)

		resourceName := "akamai_dns_record.txt_record"

		useClient(client, func() {
			resource.UnitTest(t, resource.TestCase{
				ProtoV6ProviderFactories: testutils.NewProtoV6ProviderFactory(NewSubprovider()),
				Steps: []resource.TestStep{
					{
						Config: testutils.LoadFixtureString(t, "testdata/TestResDnsRecord/quotation_marks/create.tf"),
						Check: resource.ComposeTestCheckFunc(
							resource.TestCheckResourceAttr(resourceName, "recordtype", "TXT"),
							resource.TestCheckResourceAttr(resourceName, "target.#", "1"),
							resource.TestCheckResourceAttr(resourceName, "target.0", "\"v=spf1 mx include:spf.domain.com include:spf.protection.outlook.com -all\""),
						),
					},
					{
						Config: testutils.LoadFixtureString(t, "testdata/TestResDnsRecord/quotation_marks/update.tf"),
						Check: resource.ComposeTestCheckFunc(
							resource.TestCheckResourceAttr(resourceName, "recordtype", "TXT"),
							resource.TestCheckResourceAttr(resourceName, "target.#", "1"),
							resource.TestCheckResourceAttr(resourceName, "target.0", "v=spf1 mx include:spf.domain.com include:spf.protection.outlook.com -all"),
						),
					},
				},
			})
		})

		client.AssertExpectations(t)
	})

	t.Run("SRV record with default values", func(t *testing.T) {
		t.Parallel()
		client := edgegrid.NewTestClient()

		targetBig := "10 60 5060 big.example.com."
		targetSmall := "10 60 5060 small.example.com."
		targetTiny := "10 60 5060 tiny.example.com."

		client.DNS.On("GetRecord",
			testutils.MockContext,
			dns.GetRecordRequest{Zone: "origin.org", Name: "origin.example.org", RecordType: "SRV"},
		).Return(nil, notFound).Once()

		client.DNS.On("CreateRecord",
			testutils.MockContext,
			dns.CreateRecordRequest{
				Record: &dns.RecordBody{
					Name:       "origin.example.org",
					RecordType: "SRV",
					TTL:        ptr.To(300),
					Active:     false,
					Target:     []string{targetBig, targetSmall, targetTiny},
				},
				Zone:    "origin.org",
				RecLock: []bool{false},
			},
		).Return(nil)

		client.DNS.On("GetRecord",
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

		client.DNS.On("ParseRData",
			testutils.MockContext,
			"SRV",
			[]string{targetBig, targetSmall, targetTiny},
		).Return(
			c.ParseRData(context.Background(), "SRV", []string{targetBig, targetSmall, targetTiny}),
		).Times(2)

		client.DNS.On("ProcessRdata",
			testutils.MockContext,
			[]string{targetBig, targetSmall, targetTiny},
			"SRV",
		).Return([]string{targetBig, targetSmall, targetTiny}).Times(2)

		client.DNS.On("GetRecord",
			testutils.MockContext,
			dns.GetRecordRequest{Zone: "origin.org", Name: "origin.example.org", RecordType: "SRV"},
		).Return(&dns.GetRecordResponse{
			Name:       "origin.example.org",
			RecordType: "SRV",
			TTL:        300,
			Active:     false,
			Target:     []string{targetBig, targetSmall, targetTiny},
		}, nil).Once()

		client.DNS.On("DeleteRecord",
			testutils.MockContext,
			dns.DeleteRecordRequest{Zone: "origin.org", Name: "origin.example.org", RecordType: "SRV", RecLock: []bool{false}},
		).Return(nil)

		resourceName := "akamai_dns_record.srv_record"

		resource.UnitTest(t, resource.TestCase{
			ProtoV6ProviderFactories: testutils.NewTestProtoV6SDKProviderFactory(client, newSubproviderWithConfig(testSubproviderConfig())),
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

		client.DNS.AssertExpectations(t)
	})
	t.Run("SRV record without default values", func(t *testing.T) {
		t.Parallel()
		client := edgegrid.NewTestClient()

		targetBig := "10 60 5060 big.example.com."
		targetSmall := "10 40 5060 small.example.com."
		targetTiny := "20 100 5060 tiny.example.com."

		client.DNS.On("GetRecord",
			testutils.MockContext,
			dns.GetRecordRequest{Zone: "origin.org", Name: "origin.example.org", RecordType: "SRV"},
		).Return(nil, notFound).Once()

		client.DNS.On("CreateRecord",
			testutils.MockContext,
			dns.CreateRecordRequest{
				Record: &dns.RecordBody{
					Name:       "origin.example.org",
					RecordType: "SRV",
					TTL:        ptr.To(300),
					Active:     false,
					Target:     []string{targetBig, targetSmall, targetTiny},
				},
				Zone:    "origin.org",
				RecLock: []bool{false},
			},
		).Return(nil)

		client.DNS.On("GetRecord",
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

		client.DNS.On("ParseRData",
			testutils.MockContext,
			"SRV",
			[]string{targetBig, targetSmall, targetTiny},
		).Return(
			c.ParseRData(context.Background(), "SRV", []string{targetBig, targetSmall, targetTiny}),
		).Times(2)

		client.DNS.On("ProcessRdata",
			testutils.MockContext,
			[]string{targetBig, targetSmall, targetTiny},
			"SRV",
		).Return([]string{targetBig, targetSmall, targetTiny}).Times(2)

		client.DNS.On("GetRecord",
			testutils.MockContext,
			dns.GetRecordRequest{Zone: "origin.org", Name: "origin.example.org", RecordType: "SRV"},
		).Return(&dns.GetRecordResponse{
			Name:       "origin.example.org",
			RecordType: "SRV",
			TTL:        300,
			Active:     false,
			Target:     []string{targetBig, targetSmall, targetTiny},
		}, nil).Once()

		client.DNS.On("DeleteRecord",
			testutils.MockContext,
			dns.DeleteRecordRequest{Zone: "origin.org", Name: "origin.example.org", RecordType: "SRV", RecLock: []bool{false}},
		).Return(nil)

		resourceName := "akamai_dns_record.srv_record"

		resource.UnitTest(t, resource.TestCase{
			ProtoV6ProviderFactories: testutils.NewTestProtoV6SDKProviderFactory(client, newSubproviderWithConfig(testSubproviderConfig())),
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

		client.DNS.AssertExpectations(t)
	})
	t.Run("SRV record with invalid mixed values", func(t *testing.T) {
		t.Parallel()
		client := edgegrid.NewTestClient()

		resource.UnitTest(t, resource.TestCase{
			ProtoV6ProviderFactories: testutils.NewTestProtoV6SDKProviderFactory(client, newSubproviderWithConfig(testSubproviderConfig())),
			Steps: []resource.TestStep{
				{
					Config:      testutils.LoadFixtureString(t, "testdata/TestResDnsRecord/srv/create_basic_srv_mix_invalid.tf"),
					ExpectError: regexp.MustCompile("target should consist of only simple or complete items"),
				},
			},
		})

		client.DNS.AssertExpectations(t)
	})

	t.Run("AAAA record with valid IPv6 addresses", func(t *testing.T) {
		t.Parallel()
		client := edgegrid.NewTestClient()

		target := []string{"2001:db8::68", "::ffff:192.0.2.1"}
		targetSent := []string{"0000:0000:0000:0000:0000:ffff:c000:0201", "2001:0db8:0000:0000:0000:0000:0000:0068"}
		targetReceived := []string{"2001:db8:0:0:0:0:0:68", "::ffff:192.0.2.1"}
		client.DNS.On("GetRecord",
			testutils.MockContext,
			dns.GetRecordRequest{Zone: "exampleterraform.io", Name: "exampleterraform.io", RecordType: "AAAA"},
		).Return(nil, notFound).Once()

		client.DNS.On("CreateRecord",
			testutils.MockContext,
			dns.CreateRecordRequest{
				Record: &dns.RecordBody{
					Name:       "exampleterraform.io",
					RecordType: "AAAA",
					TTL:        ptr.To(300),
					Active:     false,
					Target:     targetSent,
				},
				Zone:    "exampleterraform.io",
				RecLock: []bool{false},
			},
		).Return(nil)

		client.DNS.On("GetRecord",
			testutils.MockContext,
			dns.GetRecordRequest{Zone: "exampleterraform.io", Name: "exampleterraform.io", RecordType: "AAAA"},
		).Return(&dns.GetRecordResponse{
			Name:       "exampleterraform.io",
			RecordType: "AAAA",
			TTL:        300,
			Active:     false,
			Target:     targetReceived,
		}, nil).Once()

		client.DNS.On("ParseRData",
			testutils.MockContext,
			"AAAA",
			targetReceived,
		).Return(map[string]interface{}{
			"target": target,
		}).Times(2)

		client.DNS.On("ProcessRdata",
			testutils.MockContext,
			targetReceived,
			"AAAA",
		).Return(target).Times(2)

		client.DNS.On("GetRecord",
			testutils.MockContext,
			dns.GetRecordRequest{Zone: "exampleterraform.io", Name: "exampleterraform.io", RecordType: "AAAA"},
		).Return(&dns.GetRecordResponse{
			Name:       "exampleterraform.io",
			RecordType: "AAAA",
			TTL:        300,
			Active:     false,
			Target:     targetReceived,
		}, nil).Once()

		client.DNS.On("DeleteRecord",
			testutils.MockContext,
			dns.DeleteRecordRequest{Zone: "exampleterraform.io", Name: "exampleterraform.io", RecordType: "AAAA", RecLock: []bool{false}},
		).Return(nil)

		resourceName := "akamai_dns_record.aaaa_record"

		resource.UnitTest(t, resource.TestCase{
			ProtoV6ProviderFactories: testutils.NewTestProtoV6SDKProviderFactory(client, newSubproviderWithConfig(testSubproviderConfig())),
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

		client.DNS.AssertExpectations(t)
	})

	t.Run("A record with ttl set to 0", func(t *testing.T) {
		t.Parallel()
		client := edgegrid.NewTestClient()

		// read
		client.DNS.On("GetRecord",
			testutils.MockContext,
			dns.GetRecordRequest{
				Zone:       "exampleterraform.io",
				Name:       "exampleterraform.io",
				RecordType: "A",
			},
		).Return(nil, notFound).Once()

		// create
		client.DNS.On("CreateRecord",
			testutils.MockContext,
			dns.CreateRecordRequest{
				Record: &dns.RecordBody{
					Name:       "exampleterraform.io",
					RecordType: "A",
					TTL:        ptr.To(0),
					Active:     false,
					Target:     []string{"10.0.0.2", "10.0.0.3"},
				},
				Zone:    "exampleterraform.io",
				RecLock: []bool{false},
			},
		).Return(nil).Once()

		// read
		client.DNS.On("GetRecord",
			testutils.MockContext,
			mock.AnythingOfType("dns.GetRecordRequest"),
		).Return(&dns.GetRecordResponse{
			Name:       "exampleterraform.io",
			RecordType: "A",
			TTL:        0,
			Active:     false,
			Target:     []string{"10.0.0.2", "10.0.0.3"},
		}, nil).Once()

		retCreate := dnsClient.ParseRData(context.Background(), "A", []string{"10.0.0.2", "10.0.0.3"})

		client.DNS.On("ParseRData",
			testutils.MockContext,
			"A",
			[]string{"10.0.0.2", "10.0.0.3"},
		).Return(retCreate).Times(3)

		client.DNS.On("ProcessRdata",
			testutils.MockContext,
			[]string{"10.0.0.2", "10.0.0.3"},
			"A",
		).Return([]string{"A"}, nil).Times(4)

		client.DNS.On("GetRecord",
			testutils.MockContext,
			dns.GetRecordRequest{
				Zone:       "exampleterraform.io",
				Name:       "exampleterraform.io",
				RecordType: "A",
			},
		).Return(&dns.GetRecordResponse{
			Name:       "exampleterraform.io",
			RecordType: "exampleterraform.io",
			TTL:        0,
			Active:     false,
			Target:     []string{"10.0.0.2", "10.0.0.3"},
		}, nil).Times(3)

		// update
		client.DNS.On("UpdateRecord",
			testutils.MockContext,
			dns.UpdateRecordRequest{
				Record: &dns.RecordBody{
					Name:       "exampleterraform.io",
					RecordType: "A",
					TTL:        ptr.To(0),
					Active:     false,
					Target:     []string{"10.0.0.4", "10.0.0.5"},
				},
				Zone:    "exampleterraform.io",
				RecLock: []bool{false},
			},
		).Return(nil).Once()

		// read
		client.DNS.On("GetRecord",
			testutils.MockContext,
			dns.GetRecordRequest{
				Zone:       "exampleterraform.io",
				Name:       "exampleterraform.io",
				RecordType: "A",
			},
		).Return(&dns.GetRecordResponse{
			Name:       "exampleterraform.io",
			RecordType: "A",
			TTL:        0,
			Active:     false,
			Target:     []string{"10.0.0.4", "10.0.0.5"},
		}, nil).Times(2)

		retUpdate := dnsClient.ParseRData(context.Background(), "A", []string{"10.0.0.4", "10.0.0.5"})

		client.DNS.On("ParseRData",
			testutils.MockContext,
			"A",
			[]string{"10.0.0.4", "10.0.0.5"},
		).Return(retUpdate).Times(2)

		client.DNS.On("ProcessRdata",
			testutils.MockContext,
			[]string{"10.0.0.4", "10.0.0.5"},
			"A",
		).Return([]string{"A"}, nil).Times(2)

		// delete
		client.DNS.On("DeleteRecord",
			testutils.MockContext,
			dns.DeleteRecordRequest{
				Zone:       "exampleterraform.io",
				Name:       "exampleterraform.io",
				RecordType: "A",
				RecLock:    []bool{false},
			},
		).Return(nil).Once()

		resourceName := "akamai_dns_record.a_record"

		resource.UnitTest(t, resource.TestCase{
			ProtoV6ProviderFactories: testutils.NewTestProtoV6SDKProviderFactory(client, newSubproviderWithConfig(testSubproviderConfig())),
			Steps: []resource.TestStep{
				{
					Config: testutils.LoadFixtureString(t, "testdata/TestResDnsRecord/create_with_ttl_0.tf"),
					Check: resource.ComposeTestCheckFunc(
						resource.TestCheckResourceAttr(resourceName, "zone", "exampleterraform.io"),
						resource.TestCheckResourceAttr(resourceName, "name", "exampleterraform.io"),
						resource.TestCheckResourceAttr(resourceName, "recordtype", "A"),
						resource.TestCheckResourceAttr(resourceName, "ttl", "0"),
						resource.TestCheckNoResourceAttr(resourceName, "active"),
						resource.TestCheckResourceAttr(resourceName, "target.#", "2"),
						resource.TestCheckResourceAttr(resourceName, "target.0", "10.0.0.2"),
						resource.TestCheckResourceAttr(resourceName, "target.1", "10.0.0.3"),
					),
				},
				{
					Config: testutils.LoadFixtureString(t, "testdata/TestResDnsRecord/update_with_ttl_0.tf"),
					Check: resource.ComposeTestCheckFunc(
						resource.TestCheckResourceAttr(resourceName, "zone", "exampleterraform.io"),
						resource.TestCheckResourceAttr(resourceName, "name", "exampleterraform.io"),
						resource.TestCheckResourceAttr(resourceName, "recordtype", "A"),
						resource.TestCheckResourceAttr(resourceName, "ttl", "0"),
						resource.TestCheckNoResourceAttr(resourceName, "active"),
						resource.TestCheckResourceAttr(resourceName, "target.#", "2"),
						resource.TestCheckResourceAttr(resourceName, "target.0", "10.0.0.4"),
						resource.TestCheckResourceAttr(resourceName, "target.1", "10.0.0.5"),
					),
				},
			},
		})

		client.DNS.AssertExpectations(t)
	})

	t.Run("AAAA record with invalid IPv6 address", func(t *testing.T) {
		t.Parallel()
		client := edgegrid.NewTestClient()

		resource.UnitTest(t, resource.TestCase{
			ProtoV6ProviderFactories: testutils.NewTestProtoV6SDKProviderFactory(client, newSubproviderWithConfig(testSubproviderConfig())),
			Steps: []resource.TestStep{
				{
					Config:      testutils.LoadFixtureString(t, "testdata/TestResDnsRecord/aaaa/create_invalid_aaaa.tf"),
					ExpectError: regexp.MustCompile("target '1111:2222:3333:4444:55555:6666:7777:8888' is not a valid address"),
				},
			},
		})

		client.DNS.AssertExpectations(t)
	})
	t.Run("AAAA record with IP4 address - invalid", func(t *testing.T) {
		t.Parallel()
		client := edgegrid.NewTestClient()

		resource.UnitTest(t, resource.TestCase{
			ProtoV6ProviderFactories: testutils.NewTestProtoV6SDKProviderFactory(client, newSubproviderWithConfig(testSubproviderConfig())),
			Steps: []resource.TestStep{
				{
					Config:      testutils.LoadFixtureString(t, "testdata/TestResDnsRecord/aaaa/create_invalid_ipv4.tf"),
					ExpectError: regexp.MustCompile("target '18.244.102.124' is not a valid IPv6 or IPv4-mapped IPv6 address"),
				},
			},
		})

		client.DNS.AssertExpectations(t)
	})

	t.Run("expect error - empty zone", func(t *testing.T) {
		t.Parallel()
		client := edgegrid.NewTestClient()

		resource.UnitTest(t, resource.TestCase{
			ProtoV6ProviderFactories: testutils.NewTestProtoV6SDKProviderFactory(client, newSubproviderWithConfig(testSubproviderConfig())),
			Steps: []resource.TestStep{
				{
					Config:      testutils.LoadFixtureString(t, "testdata/TestResDnsRecord/validation/empty_zone.tf"),
					ExpectError: regexp.MustCompile("Error: zone must not be empty"),
				},
			},
		})

		client.DNS.AssertExpectations(t)
	})

	t.Run("expect error - empty name", func(t *testing.T) {
		t.Parallel()
		client := edgegrid.NewTestClient()

		resource.UnitTest(t, resource.TestCase{
			ProtoV6ProviderFactories: testutils.NewTestProtoV6SDKProviderFactory(client, newSubproviderWithConfig(testSubproviderConfig())),
			Steps: []resource.TestStep{
				{
					Config:      testutils.LoadFixtureString(t, "testdata/TestResDnsRecord/validation/empty_name.tf"),
					ExpectError: regexp.MustCompile("configuration argument name must be set"),
				},
			},
		})

		client.DNS.AssertExpectations(t)
	})

	t.Run("expect error - empty record type", func(t *testing.T) {
		t.Parallel()
		client := edgegrid.NewTestClient()

		resource.UnitTest(t, resource.TestCase{
			ProtoV6ProviderFactories: testutils.NewTestProtoV6SDKProviderFactory(client, newSubproviderWithConfig(testSubproviderConfig())),
			Steps: []resource.TestStep{
				{
					Config:      testutils.LoadFixtureString(t, "testdata/TestResDnsRecord/validation/empty_record_type.tf"),
					ExpectError: regexp.MustCompile(`The argument "recordtype" is required, but no definition was found.`),
				},
			},
		})

		client.DNS.AssertExpectations(t)
	})

	t.Run("expect error - wrong record type", func(t *testing.T) {
		t.Parallel()
		client := edgegrid.NewTestClient()

		resource.UnitTest(t, resource.TestCase{
			ProtoV6ProviderFactories: testutils.NewTestProtoV6SDKProviderFactory(client, newSubproviderWithConfig(testSubproviderConfig())),
			Steps: []resource.TestStep{
				{
					Config: testutils.LoadFixtureString(t, "testdata/TestResDnsRecord/validation/wrong_record_type.tf"),
					ExpectError: regexp.MustCompile(`Error: expected recordtype to be one of \["A" "AAAA" "CNAME" ` +
						`"LOC" "NS" "PTR" "SPF" "TXT" "AFSDB" "DNSKEY" "DS" "HINFO" "MX" "NAPTR" "NSEC3" ` +
						`"NSEC3PARAM" "RP" "RRSIG" "SRV" "SSHFP" "SOA" "AKAMAICDN" "AKAMAITLC" "CAA" "CERT" ` +
						`"TLSA" "SVCB" "HTTPS"\], got WRONG_TYPE`),
				},
			},
		})

		client.DNS.AssertExpectations(t)
	})

	t.Run("expect error - empty ttl", func(t *testing.T) {
		t.Parallel()
		client := edgegrid.NewTestClient()

		resource.UnitTest(t, resource.TestCase{
			ProtoV6ProviderFactories: testutils.NewTestProtoV6SDKProviderFactory(client, newSubproviderWithConfig(testSubproviderConfig())),
			Steps: []resource.TestStep{
				{
					Config:      testutils.LoadFixtureString(t, "testdata/TestResDnsRecord/validation/empty_ttl.tf"),
					ExpectError: regexp.MustCompile("The argument \"ttl\" is required, but no definition was found."),
				},
			},
		})

		client.DNS.AssertExpectations(t)
	})

	t.Run("expect error - negative ttl", func(t *testing.T) {
		t.Parallel()
		client := edgegrid.NewTestClient()

		resource.UnitTest(t, resource.TestCase{
			ProtoV6ProviderFactories: testutils.NewTestProtoV6SDKProviderFactory(client, newSubproviderWithConfig(testSubproviderConfig())),
			Steps: []resource.TestStep{
				{
					Config:      testutils.LoadFixtureString(t, "testdata/TestResDnsRecord/validation/negative_ttl.tf"),
					ExpectError: regexp.MustCompile(`Error: expected ttl to be at least \(0\), got -1`),
				},
			},
		})

		client.DNS.AssertExpectations(t)
	})
}

func TestDnsRecordImport(t *testing.T) {
	t.Parallel()
	baseImportChecker := test.NewImportChecker().
		CheckEqual("recordtype", "A").
		CheckEqual("name", "exampleterraform.io").
		CheckEqual("zone", "exampleterraform.io").
		CheckEqual("target.#", "2").
		CheckEqual("target.0", "10.0.0.2").
		CheckEqual("target.1", "10.0.0.3").
		CheckEqual("ttl", "300").
		CheckMissing("active")

	tests := map[string]struct {
		zone        string
		recordName  string
		recordType  string
		init        func(d *dns.Mock)
		expectError *regexp.Regexp
		stateCheck  resource.ImportStateCheckFunc
	}{
		"basic import": {
			zone:       "exampleterraform.io",
			recordName: "exampleterraform.io",
			recordType: "A",
			init: func(d *dns.Mock) {
				d.On("GetRecord", testutils.MockContext, dns.GetRecordRequest{
					Zone:       "exampleterraform.io",
					Name:       "exampleterraform.io",
					RecordType: "A",
				}).Return(&dns.GetRecordResponse{
					Name:       "exampleterraform.io",
					RecordType: "A",
					TTL:        300,
					Active:     false,
					Target:     []string{"10.0.0.2", "10.0.0.3"},
				}, nil).Times(2)

				d.On("ProcessRdata",
					testutils.MockContext,
					[]string{"10.0.0.2", "10.0.0.3"},
					"A",
				).Return([]string{"A"}, nil).Times(2)

				d.On("ParseRData",
					testutils.MockContext,
					"A",
					[]string{"10.0.0.2", "10.0.0.3"},
				).Return(map[string]interface{}{
					"target": []string{"10.0.0.2", "10.0.0.3"},
				}).Times(2)
			},
			stateCheck: baseImportChecker.
				Build(),
		},
		"import with ttl 0": {
			zone:       "exampleterraform.io",
			recordName: "exampleterraform.io",
			recordType: "A",
			init: func(d *dns.Mock) {
				d.On("GetRecord", testutils.MockContext, dns.GetRecordRequest{
					Zone:       "exampleterraform.io",
					Name:       "exampleterraform.io",
					RecordType: "A",
				}).Return(&dns.GetRecordResponse{
					Name:       "exampleterraform.io",
					RecordType: "A",
					TTL:        0,
					Active:     false,
					Target:     []string{"10.0.0.2", "10.0.0.3"},
				}, nil).Times(2)

				d.On("ProcessRdata",
					testutils.MockContext,
					[]string{"10.0.0.2", "10.0.0.3"},
					"A",
				).Return([]string{"A"}, nil).Times(2)

				d.On("ParseRData",
					testutils.MockContext,
					"A",
					[]string{"10.0.0.2", "10.0.0.3"},
				).Return(map[string]interface{}{
					"target": []string{"10.0.0.2", "10.0.0.3"},
				}).Times(2)
			},
			stateCheck: baseImportChecker.
				CheckEqual("ttl", "0").
				Build(),
		},
		"expect error - no zone": {
			zone:        "",
			recordName:  "exampleterraform.io",
			recordType:  "A",
			expectError: regexp.MustCompile("Error: invalid ID for Zone Import: #exampleterraform.io#A"),
		},
		"expect error - no record name": {
			zone:        "exampleterraform.io",
			recordName:  "",
			recordType:  "A",
			expectError: regexp.MustCompile("Error: invalid ID for Zone Import: exampleterraform.io##A"),
		},
		"expect error - no record type": {
			zone:        "exampleterraform.io",
			recordName:  "exampleterraform.io",
			recordType:  "",
			expectError: regexp.MustCompile("Error: invalid ID for Zone Import: exampleterraform.io#exampleterraform.io#"),
		},
		"expect error - wrong ID, record not found": {
			zone:       "wrong-zone.com",
			recordName: "wrong-record",
			recordType: "WRONG",
			init: func(d *dns.Mock) {
				d.On("GetRecord", testutils.MockContext, dns.GetRecordRequest{
					Zone:       "wrong-zone.com",
					Name:       "wrong-record",
					RecordType: "WRONG",
				}).Return(nil, &dns.Error{
					StatusCode: http.StatusNotFound,
				}).Once()
			},
			expectError: regexp.MustCompile("Error: record not found"),
		},
	}

	for name, tc := range tests {
		t.Run(name, func(t *testing.T) {
			t.Parallel()
			client := edgegrid.NewTestClient()
			if tc.init != nil {
				tc.init(client.DNS)
			}
			resource.UnitTest(t, resource.TestCase{
				ProtoV6ProviderFactories: testutils.NewTestProtoV6SDKProviderFactory(client, newSubproviderWithConfig(testSubproviderConfig())),
				Steps: []resource.TestStep{
					{
						ImportStateCheck: tc.stateCheck,
						ImportStateId:    fmt.Sprintf("%s#%s#%s", tc.zone, tc.recordName, tc.recordType),
						ImportState:      true,
						ResourceName:     "akamai_dns_record.a_record",
						Config:           testutils.LoadFixtureString(t, "testdata/TestResDnsRecord/import_with_ttl_0.tf"),
						ExpectError:      tc.expectError,
					},
				},
			})
			client.DNS.AssertExpectations(t)
		})
	}
}

func TestMXRecord(t *testing.T) {
	t.Parallel()
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
		response := &dns.GetRecordResponse{
			Name:       createdRecord.Name,
			RecordType: createdRecord.RecordType,
			TTL:        *createdRecord.TTL,
			Active:     createdRecord.Active,
			Target:     createdRecord.Target,
		}
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
		response := &dns.GetRecordResponse{
			Name:       previousRecord.Name,
			RecordType: previousRecord.RecordType,
			TTL:        *previousRecord.TTL,
			Active:     previousRecord.Active,
			Target:     previousRecord.Target,
		}
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
			TTL:        ptr.To(300),
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
			TTL:        ptr.To(300),
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
			TTL:        ptr.To(300),
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
			t.Parallel()
			client := edgegrid.NewTestClient()
			defaultInit(client.DNS, test.createTargets, test.updateTargets, test.deleteTargets)
			resource.UnitTest(t, resource.TestCase{
				ProtoV6ProviderFactories: testutils.NewTestProtoV6SDKProviderFactory(client, newSubproviderWithConfig(testSubproviderConfig())),
				IsUnitTest:               true,
				Steps:                    test.steps,
			})
			client.DNS.AssertExpectations(t)
		})
	}
}

func TestTargetDiffSuppress(t *testing.T) {
	t.Parallel()
	t.Run("target is computed and recordType is AAAA", func(t *testing.T) {
		t.Parallel()
		config := schema.TestResourceDataRaw(t, getResourceDNSRecordSchema(), map[string]interface{}{"recordtype": "AAAA"})
		assert.False(t, dnsRecordTargetSuppress("target.#", "0", "", config))
	})
}

func TestResolveTxtRecordTargets(t *testing.T) {
	t.Parallel()
	denormalized := []string{"onetwo", "\"one\" \"two\""}
	normalized := []string{"\"onetwo\"", "\"one\" \"two\"", "\"one\" \"two\""}
	expected := []string{"onetwo", "\"one\" \"two\"", "\"one\" \"two\""}

	res, err := resolveTxtRecordTargets(denormalized, normalized)
	require.NoError(t, err)

	assert.Equal(t, expected, res)
}

func TestResolveTargets(t *testing.T) {
	t.Parallel()
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
			t.Parallel()
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
