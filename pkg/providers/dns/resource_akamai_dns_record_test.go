package dns

import (
	"context"
	"fmt"
	"net/http"
	"regexp"
	"strings"
	"testing"

	"github.com/akamai/AkamaiOPEN-edgegrid-golang/v13/pkg/dns"
	akalog "github.com/akamai/AkamaiOPEN-edgegrid-golang/v13/pkg/log"
	"github.com/akamai/AkamaiOPEN-edgegrid-golang/v13/pkg/ptr"
	"github.com/akamai/AkamaiOPEN-edgegrid-golang/v13/pkg/session"
	"github.com/akamai/terraform-provider-akamai/v11/internal/edgegrid"
	"github.com/akamai/terraform-provider-akamai/v11/pkg/common/test"
	"github.com/akamai/terraform-provider-akamai/v11/pkg/common/testutils"
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

	t.Run("TXT record test - update single line target", func(t *testing.T) {
		t.Parallel()

		name := "infrastructure.domain.net"
		target := []string{"\"v=spf1 mx include:spf.domain.com include:spf.protection.outlook.com -all\""}
		normalizedUpdatedTarget := []string{"\"v=spf1\" \"mx\" \"include:spf.domain.com\" \"include:spf.protection.outlook.com\" \"-all\""}

		client := edgegrid.NewTestClient()

		// create.
		client.DNS.On("GetRecord",
			testutils.MockContext,
			dns.GetRecordRequest{Zone: name, Name: name, RecordType: "TXT"},
		).Return(nil, notFound).Once()

		client.DNS.On("CreateRecord",
			testutils.MockContext,
			dns.CreateRecordRequest{
				Record: &dns.RecordBody{
					Name:       name,
					RecordType: "TXT",
					TTL:        ptr.To(1800),
					Active:     false,
					Target:     target,
				},
				Zone:    name,
				RecLock: []bool{false},
			},
		).Return(nil)

		// read 3 times: 2 after create (read + check) + 1 before update (refresh)
		client.DNS.On("GetRecord",
			testutils.MockContext,
			dns.GetRecordRequest{Zone: name, Name: name, RecordType: "TXT"},
		).Return(&dns.GetRecordResponse{
			Name:       name,
			RecordType: "TXT",
			TTL:        300,
			Active:     false,
			Target:     target,
		}, nil).Times(4)

		client.DNS.On("ParseRData",
			testutils.MockContext,
			"TXT",
			target,
		).Return(map[string]interface{}{
			"target": target,
		}).Times(3)

		client.DNS.On("ProcessRdata",
			testutils.MockContext,
			target,
			"TXT",
		).Return(target).Times(4)

		// update.
		client.DNS.On("UpdateRecord",
			testutils.MockContext,
			dns.UpdateRecordRequest{
				Record: &dns.RecordBody{
					Name:       name,
					RecordType: "TXT",
					TTL:        ptr.To(1800),
					Active:     false,
					Target:     normalizedUpdatedTarget,
				},
				Zone:    name,
				RecLock: []bool{false},
			},
		).Return(nil)

		// read 2 times after update: post-update read + final check
		client.DNS.On("GetRecord",
			testutils.MockContext,
			dns.GetRecordRequest{Zone: name, Name: name, RecordType: "TXT"},
		).Return(&dns.GetRecordResponse{
			Name:       name,
			RecordType: "TXT",
			TTL:        300,
			Active:     false,
			Target:     normalizedUpdatedTarget,
		}, nil).Times(2)

		client.DNS.On("ParseRData",
			testutils.MockContext,
			"TXT",
			normalizedUpdatedTarget,
		).Return(map[string]interface{}{
			"target": normalizedUpdatedTarget,
		}).Times(2)

		client.DNS.On("ProcessRdata",
			testutils.MockContext,
			normalizedUpdatedTarget,
			"TXT",
		).Return(normalizedUpdatedTarget).Times(2)

		// delete.
		client.DNS.On("DeleteRecord",
			testutils.MockContext,
			dns.DeleteRecordRequest{Zone: name, Name: name, RecordType: "TXT", RecLock: []bool{false}},
		).Return(nil)

		resourceName := "akamai_dns_record.txt_record"

		resource.UnitTest(t, resource.TestCase{
			ProtoV6ProviderFactories: testutils.NewTestProtoV6SDKProviderFactory(client, newSubproviderWithConfig(testSubproviderConfig())),
			Steps: []resource.TestStep{
				{
					Config: testutils.LoadFixtureString(t, "testdata/TestResDnsRecord/quotation_marks/create_txt_sl.tf"),
					Check: test.NewStateChecker(resourceName).
						CheckEqual("recordtype", "TXT").
						CheckEqual("target.#", "1").
						CheckEqual("target.0", target[0]).Build(),
				},
				{
					Config: testutils.LoadFixtureString(t, "testdata/TestResDnsRecord/quotation_marks/update_txt_sl.tf"),
					Check: test.NewStateChecker(resourceName).
						CheckEqual("recordtype", "TXT").
						CheckEqual("target.#", "1").
						CheckEqual("target.0", "v=spf1 mx include:spf.domain.com include:spf.protection.outlook.com -all").Build(),
				},
			},
		})

		client.DNS.AssertExpectations(t)
	})

	t.Run("TXT record test - update multi line target", func(t *testing.T) {
		t.Parallel()

		name := "infrastructure.domain.net"
		target := []string{"\"v=spf1 mx include:spf.domain.com include:spf.protection.outlook.com -all\"", "\"foo\" \"foo\""}
		normalizedUpdatedTarget := []string{"\"v=spf1\" \"mx\" \"include:spf.domain.com\" \"include:spf.protection.outlook.com\" \"-all\"", "\"foo\" \"foo\""}

		client := edgegrid.NewTestClient()

		// create.
		client.DNS.On("GetRecord",
			testutils.MockContext,
			dns.GetRecordRequest{Zone: name, Name: name, RecordType: "TXT"},
		).Return(nil, notFound).Once()

		client.DNS.On("CreateRecord",
			testutils.MockContext,
			dns.CreateRecordRequest{
				Record: &dns.RecordBody{
					Name:       name,
					RecordType: "TXT",
					TTL:        ptr.To(1800),
					Active:     false,
					Target:     target,
				},
				Zone:    name,
				RecLock: []bool{false},
			},
		).Return(nil)

		// read 3 times: 2 after create (read + check) + 1 before update (refresh)
		client.DNS.On("GetRecord",
			testutils.MockContext,
			dns.GetRecordRequest{Zone: name, Name: name, RecordType: "TXT"},
		).Return(&dns.GetRecordResponse{
			Name:       name,
			RecordType: "TXT",
			TTL:        300,
			Active:     false,
			Target:     target,
		}, nil).Times(4)

		client.DNS.On("ParseRData",
			testutils.MockContext,
			"TXT",
			target,
		).Return(map[string]interface{}{
			"target": target,
		}).Times(3)

		client.DNS.On("ProcessRdata",
			testutils.MockContext,
			target,
			"TXT",
		).Return(target).Times(4)

		// update.
		client.DNS.On("UpdateRecord",
			testutils.MockContext,
			dns.UpdateRecordRequest{
				Record: &dns.RecordBody{
					Name:       name,
					RecordType: "TXT",
					TTL:        ptr.To(1800),
					Active:     false,
					Target:     normalizedUpdatedTarget,
				},
				Zone:    name,
				RecLock: []bool{false},
			},
		).Return(nil)

		// read 2 times after update: post-update read + final check
		client.DNS.On("GetRecord",
			testutils.MockContext,
			dns.GetRecordRequest{Zone: name, Name: name, RecordType: "TXT"},
		).Return(&dns.GetRecordResponse{
			Name:       name,
			RecordType: "TXT",
			TTL:        300,
			Active:     false,
			Target:     normalizedUpdatedTarget,
		}, nil).Times(2)

		client.DNS.On("ParseRData",
			testutils.MockContext,
			"TXT",
			normalizedUpdatedTarget,
		).Return(map[string]interface{}{
			"target": normalizedUpdatedTarget,
		}).Times(2)

		client.DNS.On("ProcessRdata",
			testutils.MockContext,
			normalizedUpdatedTarget,
			"TXT",
		).Return(normalizedUpdatedTarget).Times(2)

		// delete.
		client.DNS.On("DeleteRecord",
			testutils.MockContext,
			dns.DeleteRecordRequest{Zone: name, Name: name, RecordType: "TXT", RecLock: []bool{false}},
		).Return(nil)

		resourceName := "akamai_dns_record.txt_record"

		resource.UnitTest(t, resource.TestCase{
			ProtoV6ProviderFactories: testutils.NewTestProtoV6SDKProviderFactory(client, newSubproviderWithConfig(testSubproviderConfig())),
			Steps: []resource.TestStep{
				{
					Config: testutils.LoadFixtureString(t, "testdata/TestResDnsRecord/quotation_marks/create_txt_ml.tf"),
					Check: test.NewStateChecker(resourceName).
						CheckEqual("recordtype", "TXT").
						CheckEqual("target.#", "2").
						CheckEqual("target.0", target[0]).
						CheckEqual("target.1", "foo foo").Build(), // this value is normalized
				},
				{
					Config: testutils.LoadFixtureString(t, "testdata/TestResDnsRecord/quotation_marks/update_txt_ml.tf"),
					Check: test.NewStateChecker(resourceName).
						CheckEqual("recordtype", "TXT").
						CheckEqual("target.#", "2").
						CheckEqual("target.0", "v=spf1 mx include:spf.domain.com include:spf.protection.outlook.com -all").
						CheckEqual("target.1", "foo foo").Build(),
				},
			},
		})

		client.DNS.AssertExpectations(t)
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
					// After the logic fix the AAAA Read normalizes both sides to full
					// form before the SHA comparison.  When the addresses are identical
					// the early-return path is taken and the state keeps the config
					// values unchanged – i.e. the abbreviated forms typed by the user.
					Check: test.NewStateChecker(resourceName).
						CheckEqual("recordtype", "AAAA").
						CheckEqual("target.#", "2").
						CheckEqual("target.0", "2001:db8::68").
						CheckEqual("target.1", "::ffff:192.0.2.1").Build(),
				},
			},
		})

		client.DNS.AssertExpectations(t)
	})

	t.Run("AAAA record update with full IPv6 form", func(t *testing.T) {
		t.Parallel()
		client := edgegrid.NewTestClient()

		// Full-form addresses already in sorted order after FullIPv6 expansion
		createTargetSent := []string{
			"1000:0000:0000:0000:0000:0000:0000:0001",
			"1000:0000:0000:0000:0000:0000:0000:0002",
			"1000:0000:0000:0000:0000:0000:0000:0003",
			"1000:0000:0000:0000:0000:0000:0000:0004",
		}
		// The real Akamai API returns addresses in abbreviated form (leading zeros
		// dropped per group, no "::" compression used).
		targetReceived := []string{
			"1000:0:0:0:0:0:0:1",
			"1000:0:0:0:0:0:0:2",
			"1000:0:0:0:0:0:0:3",
			"1000:0:0:0:0:0:0:4",
		}

		// Sorted full form of update targets
		updateTargetSent := []string{
			"1000:0000:0000:0000:0000:0000:0000:0002",
			"1000:0000:0000:0000:0000:0000:0000:0004",
			"1000:0000:0000:0000:0000:0000:0000:0005",
			"1000:0000:0000:0000:0000:0000:0000:0007",
		}
		targetReceivedUpdated := []string{
			"1000:0:0:0:0:0:0:2",
			"1000:0:0:0:0:0:0:4",
			"1000:0:0:0:0:0:0:5",
			"1000:0:0:0:0:0:0:7",
		}

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
					Target:     createTargetSent,
				},
				Zone:    "exampleterraform.io",
				RecLock: []bool{false},
			},
		).Return(nil).Once()

		// read 4 times: 2 after create (read + check) + 2 before update (refresh)
		// The mock returns abbreviated form, matching what the real Akamai API returns.
		client.DNS.On("GetRecord",
			testutils.MockContext,
			dns.GetRecordRequest{Zone: "exampleterraform.io", Name: "exampleterraform.io", RecordType: "AAAA"},
		).Return(&dns.GetRecordResponse{
			Name:       "exampleterraform.io",
			RecordType: "AAAA",
			TTL:        300,
			Active:     false,
			Target:     targetReceived,
		}, nil).Times(4)

		// ParseRData / ProcessRdata receive the raw API response (abbreviated).
		// Their return values are not used for AAAA – the switch-case overrides them.
		client.DNS.On("ParseRData",
			testutils.MockContext,
			"AAAA",
			targetReceived,
		).Return(map[string]interface{}{
			"target": targetReceived,
		}).Times(3)

		client.DNS.On("ProcessRdata",
			testutils.MockContext,
			targetReceived,
			"AAAA",
		).Return(targetReceived).Times(4)

		// update
		client.DNS.On("UpdateRecord",
			testutils.MockContext,
			dns.UpdateRecordRequest{
				Record: &dns.RecordBody{
					Name:       "exampleterraform.io",
					RecordType: "AAAA",
					TTL:        ptr.To(300),
					Active:     false,
					Target:     updateTargetSent,
				},
				Zone:    "exampleterraform.io",
				RecLock: []bool{false},
			},
		).Return(nil).Once()

		// read 2 times after update: post-update read + final check
		client.DNS.On("GetRecord",
			testutils.MockContext,
			dns.GetRecordRequest{Zone: "exampleterraform.io", Name: "exampleterraform.io", RecordType: "AAAA"},
		).Return(&dns.GetRecordResponse{
			Name:       "exampleterraform.io",
			RecordType: "AAAA",
			TTL:        300,
			Active:     false,
			Target:     targetReceivedUpdated,
		}, nil).Times(2)

		client.DNS.On("ParseRData",
			testutils.MockContext,
			"AAAA",
			targetReceivedUpdated,
		).Return(map[string]interface{}{
			"target": targetReceivedUpdated,
		}).Times(2)

		client.DNS.On("ProcessRdata",
			testutils.MockContext,
			targetReceivedUpdated,
			"AAAA",
		).Return(targetReceivedUpdated).Times(2)

		client.DNS.On("DeleteRecord",
			testutils.MockContext,
			dns.DeleteRecordRequest{Zone: "exampleterraform.io", Name: "exampleterraform.io", RecordType: "AAAA", RecLock: []bool{false}},
		).Return(nil).Once()

		resourceName := "akamai_dns_record.aaaa_record_full"

		resource.UnitTest(t, resource.TestCase{
			ProtoV6ProviderFactories: testutils.NewTestProtoV6SDKProviderFactory(client, newSubproviderWithConfig(testSubproviderConfig())),
			Steps: []resource.TestStep{
				{
					Config: testutils.LoadFixtureString(t, "testdata/TestResDnsRecord/aaaa/create_full_form_aaaa.tf"),
					Check: test.NewStateChecker(resourceName).
						CheckEqual("recordtype", "AAAA").
						CheckEqual("target.#", "4").
						CheckEqual("target.0", "1000:0000:0000:0000:0000:0000:0000:0001").
						CheckEqual("target.1", "1000:0000:0000:0000:0000:0000:0000:0002").
						CheckEqual("target.2", "1000:0000:0000:0000:0000:0000:0000:0003").
						CheckEqual("target.3", "1000:0000:0000:0000:0000:0000:0000:0004").Build(),
				},
				{
					Config: testutils.LoadFixtureString(t, "testdata/TestResDnsRecord/aaaa/update_full_form_aaaa.tf"),
					Check: test.NewStateChecker(resourceName).
						CheckEqual("recordtype", "AAAA").
						CheckEqual("target.#", "4").
						CheckEqual("target.0", "1000:0000:0000:0000:0000:0000:0000:0005").
						CheckEqual("target.1", "1000:0000:0000:0000:0000:0000:0000:0002").
						CheckEqual("target.2", "1000:0000:0000:0000:0000:0000:0000:0007").
						CheckEqual("target.3", "1000:0000:0000:0000:0000:0000:0000:0004").Build(),
				},
			},
		})

		client.DNS.AssertExpectations(t)
	})

	t.Run("AAAA record update with abbreviated IPv6 form", func(t *testing.T) {
		t.Parallel()
		client := edgegrid.NewTestClient()

		// Full sorted form sent to the API (after FullIPv6 expansion + sort)
		createTargetSent := []string{
			"1000:0000:0000:0000:0000:0000:0000:0002",
			"1000:0000:0000:0000:0000:0000:0000:0004",
			"1000:0000:0000:0000:0000:0000:0000:0005",
			"1000:0000:0000:0000:0000:0000:0000:0007",
		}
		// The real Akamai API returns abbreviated form (leading zeros dropped per
		// group, no "::" compression used).
		targetReceived := []string{
			"1000:0:0:0:0:0:0:2",
			"1000:0:0:0:0:0:0:4",
			"1000:0:0:0:0:0:0:5",
			"1000:0:0:0:0:0:0:7",
		}

		// Full sorted form of update targets
		updateTargetSent := []string{
			"1000:0000:0000:0000:0000:0000:0000:0001",
			"1000:0000:0000:0000:0000:0000:0000:0002",
			"1000:0000:0000:0000:0000:0000:0000:0003",
			"1000:0000:0000:0000:0000:0000:0000:0004",
		}
		targetReceivedUpdated := []string{
			"1000:0:0:0:0:0:0:1",
			"1000:0:0:0:0:0:0:2",
			"1000:0:0:0:0:0:0:3",
			"1000:0:0:0:0:0:0:4",
		}

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
					Target:     createTargetSent,
				},
				Zone:    "exampleterraform.io",
				RecLock: []bool{false},
			},
		).Return(nil).Once()

		// read 4 times: 2 after create (read + check) + 2 before update (refresh)
		// The mock returns abbreviated form, matching what the real Akamai API returns.
		client.DNS.On("GetRecord",
			testutils.MockContext,
			dns.GetRecordRequest{Zone: "exampleterraform.io", Name: "exampleterraform.io", RecordType: "AAAA"},
		).Return(&dns.GetRecordResponse{
			Name:       "exampleterraform.io",
			RecordType: "AAAA",
			TTL:        300,
			Active:     false,
			Target:     targetReceived,
		}, nil).Times(4)

		// ParseRData / ProcessRdata receive the raw API response (abbreviated).
		// Their return values are not used for AAAA – the switch-case overrides them.
		client.DNS.On("ParseRData",
			testutils.MockContext,
			"AAAA",
			targetReceived,
		).Return(map[string]interface{}{
			"target": targetReceived,
		}).Times(3)

		client.DNS.On("ProcessRdata",
			testutils.MockContext,
			targetReceived,
			"AAAA",
		).Return(targetReceived).Times(4)

		// update
		client.DNS.On("UpdateRecord",
			testutils.MockContext,
			dns.UpdateRecordRequest{
				Record: &dns.RecordBody{
					Name:       "exampleterraform.io",
					RecordType: "AAAA",
					TTL:        ptr.To(300),
					Active:     false,
					Target:     updateTargetSent,
				},
				Zone:    "exampleterraform.io",
				RecLock: []bool{false},
			},
		).Return(nil).Once()

		// read 2 times after update: post-update read + final check
		client.DNS.On("GetRecord",
			testutils.MockContext,
			dns.GetRecordRequest{Zone: "exampleterraform.io", Name: "exampleterraform.io", RecordType: "AAAA"},
		).Return(&dns.GetRecordResponse{
			Name:       "exampleterraform.io",
			RecordType: "AAAA",
			TTL:        300,
			Active:     false,
			Target:     targetReceivedUpdated,
		}, nil).Times(2)

		client.DNS.On("ParseRData",
			testutils.MockContext,
			"AAAA",
			targetReceivedUpdated,
		).Return(map[string]interface{}{
			"target": targetReceivedUpdated,
		}).Times(2)

		client.DNS.On("ProcessRdata",
			testutils.MockContext,
			targetReceivedUpdated,
			"AAAA",
		).Return(targetReceivedUpdated).Times(2)

		client.DNS.On("DeleteRecord",
			testutils.MockContext,
			dns.DeleteRecordRequest{Zone: "exampleterraform.io", Name: "exampleterraform.io", RecordType: "AAAA", RecLock: []bool{false}},
		).Return(nil).Once()

		resourceName := "akamai_dns_record.aaaa_record_abbrev"

		resource.UnitTest(t, resource.TestCase{
			ProtoV6ProviderFactories: testutils.NewTestProtoV6SDKProviderFactory(client, newSubproviderWithConfig(testSubproviderConfig())),
			Steps: []resource.TestStep{
				{
					Config: testutils.LoadFixtureString(t, "testdata/TestResDnsRecord/aaaa/create_abbrev_form_aaaa.tf"),
					Check: test.NewStateChecker(resourceName).
						CheckEqual("recordtype", "AAAA").
						CheckEqual("target.#", "4").
						CheckEqual("target.0", "1000:0:0:0:0:0:0:5").
						CheckEqual("target.1", "1000:0:0:0:0:0:0:2").
						CheckEqual("target.2", "1000:0:0:0:0:0:0:7").
						CheckEqual("target.3", "1000:0:0:0:0:0:0:4").Build(),
				},
				{
					Config: testutils.LoadFixtureString(t, "testdata/TestResDnsRecord/aaaa/update_abbrev_form_aaaa.tf"),
					Check: test.NewStateChecker(resourceName).
						CheckEqual("recordtype", "AAAA").
						CheckEqual("target.#", "4").
						CheckEqual("target.0", "1000:0:0:0:0:0:0:1").
						CheckEqual("target.1", "1000:0:0:0:0:0:0:2").
						CheckEqual("target.2", "1000:0:0:0:0:0:0:3").
						CheckEqual("target.3", "1000:0:0:0:0:0:0:4").Build(),
				},
			},
		})

		client.DNS.AssertExpectations(t)
	})

	// CAA records: the Akamai API stores/returns values with surrounding
	// double-quotes (e.g. `0 issue "ca.example.net"`) while users typically
	// omit the quotes in their Terraform config.  Without normalisation in
	// the Read function, state would hold the quoted form, causing every plan
	// that also changes the list length to show the existing record as being
	// removed-and-re-added — exactly the "nasty diff" the user reported.
	t.Run("CAA record – API quoted form normalised to unquoted in state", func(t *testing.T) {
		t.Parallel()
		client := edgegrid.NewTestClient()

		// buildRecordsList always adds quotes, so the API receives quoted values.
		createTargetSent := []string{`0 issue "ca.example.net"`}

		// The Akamai API returns values with literal double-quotes.
		apiTarget1 := []string{`0 issue "ca.example.net"`}

		// After expanding to 3 records the update sends all three (quoted, sorted).
		updateTargetSent := []string{
			`0 iodef "https://example.com/iodef"`,
			`0 issue "ca.example.net"`,
			`0 issuewild "ca.example.net"`,
		}
		// API returns all three (with quotes).
		apiTarget3 := []string{
			`0 iodef "https://example.com/iodef"`,
			`0 issue "ca.example.net"`,
			`0 issuewild "ca.example.net"`,
		}

		// --- Step 1 mocks (create + 2 post-create reads) ---

		client.DNS.On("GetRecord",
			testutils.MockContext,
			dns.GetRecordRequest{Zone: "exampleterraform.io", Name: "exampleterraform.io", RecordType: "CAA"},
		).Return(nil, notFound).Once()

		client.DNS.On("CreateRecord",
			testutils.MockContext,
			dns.CreateRecordRequest{
				Record: &dns.RecordBody{
					Name:       "exampleterraform.io",
					RecordType: "CAA",
					TTL:        ptr.To(300),
					Active:     false,
					Target:     createTargetSent,
				},
				Zone:    "exampleterraform.io",
				RecLock: []bool{false},
			},
		).Return(nil).Once()

		// read 3 times: 2 post-create + 1 pre-update refresh
		client.DNS.On("GetRecord",
			testutils.MockContext,
			dns.GetRecordRequest{Zone: "exampleterraform.io", Name: "exampleterraform.io", RecordType: "CAA"},
		).Return(&dns.GetRecordResponse{
			Name:       "exampleterraform.io",
			RecordType: "CAA",
			TTL:        300,
			Active:     false,
			Target:     apiTarget1,
		}, nil).Times(3)

		// ParseRData / ProcessRdata are called for every successful GetRecord.
		// The CAA Read case overrides the results, so the return value is unused
		// except for satisfying the interface — anything non-nil is fine.
		client.DNS.On("ParseRData",
			testutils.MockContext,
			"CAA",
			apiTarget1,
		).Return(map[string]interface{}{
			"target": apiTarget1,
		}).Times(3)

		client.DNS.On("ProcessRdata",
			testutils.MockContext,
			apiTarget1,
			"CAA",
		).Return(apiTarget1).Times(3)

		// --- Step 2 mocks (update + 2 post-update reads) ---

		client.DNS.On("UpdateRecord",
			testutils.MockContext,
			dns.UpdateRecordRequest{
				Record: &dns.RecordBody{
					Name:       "exampleterraform.io",
					RecordType: "CAA",
					TTL:        ptr.To(300),
					Active:     false,
					Target:     updateTargetSent,
				},
				Zone:    "exampleterraform.io",
				RecLock: []bool{false},
			},
		).Return(nil).Once()

		client.DNS.On("GetRecord",
			testutils.MockContext,
			dns.GetRecordRequest{Zone: "exampleterraform.io", Name: "exampleterraform.io", RecordType: "CAA"},
		).Return(&dns.GetRecordResponse{
			Name:       "exampleterraform.io",
			RecordType: "CAA",
			TTL:        300,
			Active:     false,
			Target:     apiTarget3,
		}, nil).Times(3)

		client.DNS.On("ParseRData",
			testutils.MockContext,
			"CAA",
			apiTarget3,
		).Return(map[string]interface{}{
			"target": apiTarget3,
		}).Times(2)

		client.DNS.On("ProcessRdata",
			testutils.MockContext,
			apiTarget3,
			"CAA",
		).Return(apiTarget3).Times(3)

		client.DNS.On("DeleteRecord",
			testutils.MockContext,
			dns.DeleteRecordRequest{Zone: "exampleterraform.io", Name: "exampleterraform.io", RecordType: "CAA", RecLock: []bool{false}},
		).Return(nil).Once()

		resourceName := "akamai_dns_record.caa_record"

		resource.UnitTest(t, resource.TestCase{
			ProtoV6ProviderFactories: testutils.NewTestProtoV6SDKProviderFactory(client, newSubproviderWithConfig(testSubproviderConfig())),
			Steps: []resource.TestStep{
				{
					// After create the Read early-returns (API matches config logically),
					// so state preserves the exact config value — the quoted form
					// 0 issue "ca.example.net".
					Config: testutils.LoadFixtureString(t, "testdata/TestResDnsRecord/caa/create_caa_quoted.tf"),
					Check: test.NewStateChecker(resourceName).
						CheckEqual("recordtype", "CAA").
						CheckEqual("target.#", "1").
						CheckEqual("target.0", `0 issue "ca.example.net"`).Build(),
				},
				{
					// Expanding to 3 records: element 0 stays in its original quoted
					// form (early return preserves state); DiffSuppressFunc handles
					// format differences so only the 2 new records appear in the plan.
					Config: testutils.LoadFixtureString(t, "testdata/TestResDnsRecord/caa/update_caa_quoted.tf"),
					Check: test.NewStateChecker(resourceName).
						CheckEqual("recordtype", "CAA").
						CheckEqual("target.#", "3").
						CheckEqual("target.0", `0 issue "ca.example.net"`).
						CheckEqual("target.1", `0 issuewild "ca.example.net"`).
						CheckEqual("target.2", `0 iodef "https://example.com/iodef"`).Build(),
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

func TestDiffQuotedDNSRecordTXT(t *testing.T) {
	t.Parallel()
	logger := akalog.NOPLogger()

	spfValue := "v=spf1 mx include:spf.domain.com include:spf.protection.outlook.com -all"
	spfQuoted := "\"v=spf1 mx include:spf.domain.com include:spf.protection.outlook.com -all\""
	spfNormalizedChunks := "\"v=spf1\" \"mx\" \"include:spf.domain.com\" \"include:spf.protection.outlook.com\" \"-all\""

	tests := []struct {
		name           string
		recordType     string
		oldTargetList  []string
		newTargetList  []string
		oldVal         string
		newVal         string
		expectSuppress bool
	}{
		{
			name:           "TXT: removing outer quotes triggers diff",
			recordType:     RRTypeTxt,
			oldTargetList:  []string{spfQuoted},
			newTargetList:  []string{spfValue},
			oldVal:         spfQuoted,
			newVal:         spfValue,
			expectSuppress: false,
		},
		{
			name:           "TXT: same quoted value suppresses diff",
			recordType:     RRTypeTxt,
			oldTargetList:  []string{spfQuoted},
			newTargetList:  []string{spfQuoted},
			oldVal:         spfQuoted,
			newVal:         spfQuoted,
			expectSuppress: true,
		},
		{
			name:           "TXT: same unquoted value suppresses diff",
			recordType:     RRTypeTxt,
			oldTargetList:  []string{spfValue},
			newTargetList:  []string{spfValue},
			oldVal:         spfValue,
			newVal:         spfValue,
			expectSuppress: true,
		},
		{
			name:           "TXT: normalized chunked form vs normalized chunked form suppresses diff",
			recordType:     RRTypeTxt,
			oldTargetList:  []string{spfNormalizedChunks},
			newTargetList:  []string{spfNormalizedChunks},
			oldVal:         spfNormalizedChunks,
			newVal:         spfNormalizedChunks,
			expectSuppress: true,
		},
		{
			name:           "TXT: adding outer quotes triggers diff",
			recordType:     RRTypeTxt,
			oldTargetList:  []string{spfValue},
			newTargetList:  []string{spfQuoted},
			oldVal:         spfValue,
			newVal:         spfQuoted,
			expectSuppress: false,
		},
		{
			name:           "TXT: order change suppresses diff (set behavior)",
			recordType:     RRTypeTxt,
			oldTargetList:  []string{"\"one\"", "\"two\""},
			newTargetList:  []string{"\"two\"", "\"one\""},
			oldVal:         "\"two\"",
			newVal:         "\"one\"",
			expectSuppress: true,
		},
		{
			name:           "TXT: different values trigger diff",
			recordType:     RRTypeTxt,
			oldTargetList:  []string{spfQuoted},
			newTargetList:  []string{"\"v=spf1 include:other.com ~all\""},
			oldVal:         spfQuoted,
			newVal:         "\"v=spf1 include:other.com ~all\"",
			expectSuppress: false,
		},
		{
			name:           "TXT: different list lengths trigger diff",
			recordType:     RRTypeTxt,
			oldTargetList:  []string{spfQuoted},
			newTargetList:  []string{spfQuoted, "\"another\""},
			oldVal:         spfQuoted,
			newVal:         "\"another\"",
			expectSuppress: false,
		},
		{
			name:           "TXT: adding quotes should not trigger diff",
			recordType:     RRTypeTxt,
			oldTargetList:  []string{"abc"},
			newTargetList:  []string{"\"abc\""},
			oldVal:         "abc",
			newVal:         "\"abc\"",
			expectSuppress: true,
		},
		{
			name:           "TXT: removing quotes should not trigger diff",
			recordType:     RRTypeTxt,
			oldTargetList:  []string{"\"abc\""},
			newTargetList:  []string{"abc"},
			oldVal:         "\"abc\"",
			newVal:         "abc",
			expectSuppress: true,
		},
	}

	for _, tc := range tests {
		t.Run(tc.name, func(t *testing.T) {
			t.Parallel()
			result := diffQuotedDNSRecord(tc.oldTargetList, tc.newTargetList, tc.oldVal, tc.newVal, tc.recordType, logger)
			assert.Equal(t, tc.expectSuppress, result)
		})
	}
}

func TestDiffQuotedDNSRecordA(t *testing.T) {
	t.Parallel()
	logger := akalog.NOPLogger()

	tests := []struct {
		name           string
		oldTargetList  []string
		newTargetList  []string
		oldVal         string
		newVal         string
		expectSuppress bool
	}{
		{
			name:           "A: sliding-window update — old value present in new list at different index must not suppress new value",
			oldTargetList:  []string{"10.0.1.1", "10.0.1.2"},
			newTargetList:  []string{"10.0.1.2", "10.0.1.3"},
			oldVal:         "10.0.1.2",
			newVal:         "10.0.1.3",
			expectSuppress: false,
		},
		{
			name:           "A: identical lists suppress diff",
			oldTargetList:  []string{"10.0.1.1", "10.0.1.2"},
			newTargetList:  []string{"10.0.1.1", "10.0.1.2"},
			oldVal:         "10.0.1.1",
			newVal:         "10.0.1.1",
			expectSuppress: true,
		},
		{
			name:           "A: pure reorder suppresses diff",
			oldTargetList:  []string{"10.0.1.1", "10.0.1.2"},
			newTargetList:  []string{"10.0.1.2", "10.0.1.1"},
			oldVal:         "10.0.1.1",
			newVal:         "10.0.1.2",
			expectSuppress: true,
		},
		{
			name:           "A: all values replaced triggers diff",
			oldTargetList:  []string{"10.0.1.1", "10.0.1.2"},
			newTargetList:  []string{"10.0.1.3", "10.0.1.4"},
			oldVal:         "10.0.1.1",
			newVal:         "10.0.1.3",
			expectSuppress: false,
		},
		{
			name:           "A: different list lengths trigger diff",
			oldTargetList:  []string{"10.0.1.1"},
			newTargetList:  []string{"10.0.1.1", "10.0.1.2"},
			oldVal:         "10.0.1.1",
			newVal:         "10.0.1.2",
			expectSuppress: false,
		},
		// Per-element suppress: when some IPs genuinely change, an element whose
		// value is exactly unchanged should be suppressed so it does not appear
		// as plan noise alongside the real change.
		{
			// 10.0.1.2 is unchanged at position 1; 10.0.1.1→10.0.1.5 is genuine.
			name:           "A: unchanged element suppressed when another IP genuinely changes",
			oldTargetList:  []string{"10.0.1.1", "10.0.1.2", "10.0.1.3", "10.0.1.4"},
			newTargetList:  []string{"10.0.1.5", "10.0.1.2", "10.0.1.7", "10.0.1.4"},
			oldVal:         "10.0.1.2",
			newVal:         "10.0.1.2",
			expectSuppress: true,
		},
		{
			// 10.0.1.1→10.0.1.5 is genuine; must not be suppressed.
			name:           "A: genuinely changed IP not suppressed when another IP is unchanged",
			oldTargetList:  []string{"10.0.1.1", "10.0.1.2", "10.0.1.3", "10.0.1.4"},
			newTargetList:  []string{"10.0.1.5", "10.0.1.2", "10.0.1.7", "10.0.1.4"},
			oldVal:         "10.0.1.1",
			newVal:         "10.0.1.5",
			expectSuppress: false,
		},
	}

	for _, tc := range tests {
		t.Run(tc.name, func(t *testing.T) {
			t.Parallel()
			result := diffQuotedDNSRecord(tc.oldTargetList, tc.newTargetList, tc.oldVal, tc.newVal, RRTypeA, logger)
			assert.Equal(t, tc.expectSuppress, result)
		})
	}
}

func TestDiffQuotedDNSRecordAAAA(t *testing.T) {
	t.Parallel()
	logger := akalog.NOPLogger()

	tests := []struct {
		name           string
		oldTargetList  []string
		newTargetList  []string
		oldVal         string
		newVal         string
		expectSuppress bool
	}{
		{
			name:           "AAAA: sliding-window update — old address present in new list at different index must not suppress new address",
			oldTargetList:  []string{"2001:db8::1", "2001:db8::2"},
			newTargetList:  []string{"2001:db8::2", "2001:db8::3"},
			oldVal:         "2001:db8::2",
			newVal:         "2001:db8::3",
			expectSuppress: false,
		},
		{
			name:           "AAAA: abbreviated and expanded forms of the same address suppress diff",
			oldTargetList:  []string{"2001:db8::1"},
			newTargetList:  []string{"2001:0db8:0000:0000:0000:0000:0000:0001"},
			oldVal:         "2001:db8::1",
			newVal:         "2001:0db8:0000:0000:0000:0000:0000:0001",
			expectSuppress: true,
		},
		{
			name:           "AAAA: identical lists suppress diff",
			oldTargetList:  []string{"2001:db8::1", "2001:db8::2"},
			newTargetList:  []string{"2001:db8::1", "2001:db8::2"},
			oldVal:         "2001:db8::1",
			newVal:         "2001:db8::1",
			expectSuppress: true,
		},
		{
			name:           "AAAA: pure reorder suppresses diff",
			oldTargetList:  []string{"2001:db8::1", "2001:db8::2"},
			newTargetList:  []string{"2001:db8::2", "2001:db8::1"},
			oldVal:         "2001:db8::1",
			newVal:         "2001:db8::2",
			expectSuppress: true,
		},
		{
			name:           "AAAA: pure reorder with expanded form suppresses diff",
			oldTargetList:  []string{"2001:db8::1", "2001:db8::2"},
			newTargetList:  []string{"2001:db8::2", "2001:0db8:0000:0000:0000:0000:0000:0001"},
			oldVal:         "2001:db8::1",
			newVal:         "2001:db8::2",
			expectSuppress: true,
		},
		{
			name:           "AAAA: all addresses replaced triggers diff",
			oldTargetList:  []string{"2001:db8::1", "2001:db8::2"},
			newTargetList:  []string{"2001:db8::3", "2001:db8::4"},
			oldVal:         "2001:db8::1",
			newVal:         "2001:db8::3",
			expectSuppress: false,
		},
		{
			name:           "AAAA: different list lengths trigger diff",
			oldTargetList:  []string{"2001:db8::1"},
			newTargetList:  []string{"2001:db8::1", "2001:db8::2"},
			oldVal:         "2001:db8::1",
			newVal:         "2001:db8::2",
			expectSuppress: false,
		},
		// Per-element suppress: when some addresses genuinely change, elements
		// that are the same address in a different notation should be suppressed
		// so they don't appear as noise in the Terraform plan.
		{
			// The element under test (0:2 → 0000:0002) is unchanged; two other
			// addresses in the list genuinely changed (0:1→0:5, 0:3→0:7).
			name:           "AAAA: format-only element suppressed when other elements change (abbreviated→full)",
			oldTargetList:  []string{"1000:0:0:0:0:0:0:1", "1000:0:0:0:0:0:0:2", "1000:0:0:0:0:0:0:3", "1000:0:0:0:0:0:0:4"},
			newTargetList:  []string{"1000:0000:0000:0000:0000:0000:0000:0005", "1000:0000:0000:0000:0000:0000:0000:0002", "1000:0000:0000:0000:0000:0000:0000:0007", "1000:0000:0000:0000:0000:0000:0000:0004"},
			oldVal:         "1000:0:0:0:0:0:0:2",
			newVal:         "1000:0000:0000:0000:0000:0000:0000:0002",
			expectSuppress: true,
		},
		{
			// The element under test (0:1 → 0000:0005) is a genuine change and
			// must not be suppressed.
			name:           "AAAA: genuinely changed element not suppressed when other elements also change",
			oldTargetList:  []string{"1000:0:0:0:0:0:0:1", "1000:0:0:0:0:0:0:2", "1000:0:0:0:0:0:0:3", "1000:0:0:0:0:0:0:4"},
			newTargetList:  []string{"1000:0000:0000:0000:0000:0000:0000:0005", "1000:0000:0000:0000:0000:0000:0000:0002", "1000:0000:0000:0000:0000:0000:0000:0007", "1000:0000:0000:0000:0000:0000:0000:0004"},
			oldVal:         "1000:0:0:0:0:0:0:1",
			newVal:         "1000:0000:0000:0000:0000:0000:0000:0005",
			expectSuppress: false,
		},
		{
			// Format-only change on last element suppressed while first changes.
			name:           "AAAA: format-only last element suppressed when first element changes",
			oldTargetList:  []string{"1000:0:0:0:0:0:0:1", "1000:0:0:0:0:0:0:4"},
			newTargetList:  []string{"1000:0000:0000:0000:0000:0000:0000:0005", "1000:0000:0000:0000:0000:0000:0000:0004"},
			oldVal:         "1000:0:0:0:0:0:0:4",
			newVal:         "1000:0000:0000:0000:0000:0000:0000:0004",
			expectSuppress: true,
		},
	}

	for _, tc := range tests {
		t.Run(tc.name, func(t *testing.T) {
			t.Parallel()
			result := diffQuotedDNSRecord(tc.oldTargetList, tc.newTargetList, tc.oldVal, tc.newVal, RRTypeAaaa, logger)
			assert.Equal(t, tc.expectSuppress, result)
		})
	}
}

func TestDiffQuotedDNSRecordNS(t *testing.T) {
	t.Parallel()
	logger := akalog.NOPLogger()

	tests := []struct {
		name           string
		oldTargetList  []string
		newTargetList  []string
		oldVal         string
		newVal         string
		expectSuppress bool
	}{
		{
			name:           "NS: sliding-window update — old nameserver present in new list at different index must not suppress new nameserver",
			oldTargetList:  []string{"ns1.example.com.", "ns2.example.com."},
			newTargetList:  []string{"ns2.example.com.", "ns3.example.com."},
			oldVal:         "ns2.example.com.",
			newVal:         "ns3.example.com.",
			expectSuppress: false,
		},
		{
			name:           "NS: trailing dot vs no trailing dot suppresses diff",
			oldTargetList:  []string{"ns1.example.com."},
			newTargetList:  []string{"ns1.example.com"},
			oldVal:         "ns1.example.com.",
			newVal:         "ns1.example.com",
			expectSuppress: true,
		},
		{
			name:           "NS: identical lists suppress diff",
			oldTargetList:  []string{"ns1.example.com.", "ns2.example.com."},
			newTargetList:  []string{"ns1.example.com.", "ns2.example.com."},
			oldVal:         "ns1.example.com.",
			newVal:         "ns1.example.com.",
			expectSuppress: true,
		},
		{
			name:           "NS: pure reorder suppresses diff",
			oldTargetList:  []string{"ns1.example.com.", "ns2.example.com."},
			newTargetList:  []string{"ns2.example.com.", "ns1.example.com."},
			oldVal:         "ns1.example.com.",
			newVal:         "ns2.example.com.",
			expectSuppress: true,
		},
		{
			name:           "NS: different nameserver triggers diff",
			oldTargetList:  []string{"ns1.example.com."},
			newTargetList:  []string{"ns2.example.com."},
			oldVal:         "ns1.example.com.",
			newVal:         "ns2.example.com.",
			expectSuppress: false,
		},
		{
			name:           "NS: different list lengths trigger diff",
			oldTargetList:  []string{"ns1.example.com."},
			newTargetList:  []string{"ns1.example.com.", "ns2.example.com."},
			oldVal:         "ns1.example.com.",
			newVal:         "ns2.example.com.",
			expectSuppress: false,
		},
		// Per-element suppress: when some nameservers genuinely change, elements
		// whose only difference is a trailing dot should not appear as plan noise.
		{
			// ns1 is unchanged (just trailing-dot format), ns2→ns4 is genuine change.
			// Checking ns1: should be suppressed.
			name:           "NS: trailing-dot element suppressed when another NS genuinely changes",
			oldTargetList:  []string{"ns1.example.com.", "ns2.example.com.", "ns3.example.com."},
			newTargetList:  []string{"ns1.example.com", "ns4.example.com.", "ns3.example.com"},
			oldVal:         "ns1.example.com.",
			newVal:         "ns1.example.com",
			expectSuppress: true,
		},
		{
			// ns2→ns4 is genuine; must not be suppressed.
			name:           "NS: genuinely changed element not suppressed when another NS is format-only",
			oldTargetList:  []string{"ns1.example.com.", "ns2.example.com.", "ns3.example.com."},
			newTargetList:  []string{"ns1.example.com", "ns4.example.com.", "ns3.example.com"},
			oldVal:         "ns2.example.com.",
			newVal:         "ns4.example.com.",
			expectSuppress: false,
		},
	}

	for _, tc := range tests {
		t.Run(tc.name, func(t *testing.T) {
			t.Parallel()
			result := diffQuotedDNSRecord(tc.oldTargetList, tc.newTargetList, tc.oldVal, tc.newVal, RRTypeNs, logger)
			assert.Equal(t, tc.expectSuppress, result)
		})
	}
}

// TestDiffQuotedDNSRecordCNAME exercises the trailing-dot branch for CNAME
// (which shares the same code path as AFSDB, PTR, and SRV).
func TestDiffQuotedDNSRecordCNAME(t *testing.T) {
	t.Parallel()
	logger := akalog.NOPLogger()

	tests := []struct {
		name           string
		oldTargetList  []string
		newTargetList  []string
		oldVal         string
		newVal         string
		expectSuppress bool
	}{
		{
			name:           "CNAME: identical value suppresses diff",
			oldTargetList:  []string{"alias.example.com."},
			newTargetList:  []string{"alias.example.com."},
			oldVal:         "alias.example.com.",
			newVal:         "alias.example.com.",
			expectSuppress: true,
		},
		{
			name:           "CNAME: trailing dot vs no trailing dot suppresses diff",
			oldTargetList:  []string{"alias.example.com."},
			newTargetList:  []string{"alias.example.com"},
			oldVal:         "alias.example.com.",
			newVal:         "alias.example.com",
			expectSuppress: true,
		},
		{
			name:           "CNAME: different target triggers diff",
			oldTargetList:  []string{"alias1.example.com."},
			newTargetList:  []string{"alias2.example.com."},
			oldVal:         "alias1.example.com.",
			newVal:         "alias2.example.com.",
			expectSuppress: false,
		},
		// Per-element suppress: when some targets genuinely change, elements
		// whose only difference is a trailing dot must not appear as plan noise.
		{
			// alias1 is unchanged (trailing-dot only), alias2→alias4 is genuine.
			name:           "CNAME: trailing-dot element suppressed when another target genuinely changes",
			oldTargetList:  []string{"alias1.example.com.", "alias2.example.com."},
			newTargetList:  []string{"alias1.example.com", "alias4.example.com."},
			oldVal:         "alias1.example.com.",
			newVal:         "alias1.example.com",
			expectSuppress: true,
		},
		{
			// alias2→alias4 is genuine; must not be suppressed.
			name:           "CNAME: genuinely changed element not suppressed when another is format-only",
			oldTargetList:  []string{"alias1.example.com.", "alias2.example.com."},
			newTargetList:  []string{"alias1.example.com", "alias4.example.com."},
			oldVal:         "alias2.example.com.",
			newVal:         "alias4.example.com.",
			expectSuppress: false,
		},
	}

	for _, tc := range tests {
		t.Run(tc.name, func(t *testing.T) {
			t.Parallel()
			result := diffQuotedDNSRecord(tc.oldTargetList, tc.newTargetList, tc.oldVal, tc.newVal, RRTypeCname, logger)
			assert.Equal(t, tc.expectSuppress, result)
		})
	}
}

// TestDiffQuotedDNSRecordCAA verifies the two-step suppress logic for CAA
// records.  The only format variation for CAA is whether the value part carries
// surrounding double-quotes (e.g. `0 issue "letsencrypt.org"` vs
// `0 issue letsencrypt.org`).  The old membership-search approach had a
// sliding-window bug where a genuine change could be hidden when the old value
// appeared at a different index in the new list.
func TestDiffQuotedDNSRecordCAA(t *testing.T) {
	t.Parallel()
	logger := akalog.NOPLogger()

	tests := []struct {
		name           string
		oldTargetList  []string
		newTargetList  []string
		oldVal         string
		newVal         string
		expectSuppress bool
	}{
		{
			name:           "CAA: identical value suppresses diff",
			oldTargetList:  []string{`0 issue "letsencrypt.org"`},
			newTargetList:  []string{`0 issue "letsencrypt.org"`},
			oldVal:         `0 issue "letsencrypt.org"`,
			newVal:         `0 issue "letsencrypt.org"`,
			expectSuppress: true,
		},
		{
			name:           "CAA: quoted vs unquoted value suppresses diff (format-only)",
			oldTargetList:  []string{`0 issue letsencrypt.org`},
			newTargetList:  []string{`0 issue "letsencrypt.org"`},
			oldVal:         `0 issue letsencrypt.org`,
			newVal:         `0 issue "letsencrypt.org"`,
			expectSuppress: true,
		},
		{
			name:           "CAA: different value triggers diff",
			oldTargetList:  []string{`0 issue "letsencrypt.org"`},
			newTargetList:  []string{`0 issue "digicert.com"`},
			oldVal:         `0 issue "letsencrypt.org"`,
			newVal:         `0 issue "digicert.com"`,
			expectSuppress: false,
		},
		{
			name: "CAA: pure reorder suppresses diff",
			oldTargetList: []string{
				`0 issue "letsencrypt.org"`,
				`0 issuewild "letsencrypt.org"`,
			},
			newTargetList: []string{
				`0 issuewild "letsencrypt.org"`,
				`0 issue "letsencrypt.org"`,
			},
			oldVal:         `0 issue "letsencrypt.org"`,
			newVal:         `0 issuewild "letsencrypt.org"`,
			expectSuppress: true,
		},
		{
			name:           "CAA: different list lengths trigger diff",
			oldTargetList:  []string{`0 issue "letsencrypt.org"`},
			newTargetList:  []string{`0 issue "letsencrypt.org"`, `0 issuewild "letsencrypt.org"`},
			oldVal:         `0 issue "letsencrypt.org"`,
			newVal:         `0 issuewild "letsencrypt.org"`,
			expectSuppress: false,
		},
		// Sliding-window: old value "b" present in new list at a different index
		// must NOT suppress the genuinely-changed element that now holds "c".
		{
			name: "CAA: sliding-window update — old value at different index must not suppress genuine change",
			oldTargetList: []string{
				`0 issue "a.com"`,
				`0 issue "b.com"`,
			},
			newTargetList: []string{
				`0 issue "b.com"`,
				`0 issue "c.com"`,
			},
			// target.1 is genuinely changing b→c; old "b.com" is in new list at
			// index 0, but the per-element check must not suppress this.
			oldVal:         `0 issue "b.com"`,
			newVal:         `0 issue "c.com"`,
			expectSuppress: false,
		},
		// Per-element suppress: when some records genuinely change, elements
		// whose only difference is a quote format should not appear as plan noise.
		{
			// entry[0] is unchanged (just quote format), entry[1] is a genuine change.
			// Checking entry[0]: should be suppressed.
			name: "CAA: quote-format element suppressed when another CAA genuinely changes",
			oldTargetList: []string{
				`0 issue letsencrypt.org`,
				`0 issue "b.com"`,
				`0 issuewild "letsencrypt.org"`,
			},
			newTargetList: []string{
				`0 issue "letsencrypt.org"`,
				`0 issue "c.com"`,
				`0 issuewild "letsencrypt.org"`,
			},
			oldVal:         `0 issue letsencrypt.org`,
			newVal:         `0 issue "letsencrypt.org"`,
			expectSuppress: true,
		},
		{
			// entry[1] b→c is genuine; must not be suppressed.
			name: "CAA: genuinely changed element not suppressed when another CAA is format-only",
			oldTargetList: []string{
				`0 issue letsencrypt.org`,
				`0 issue "b.com"`,
				`0 issuewild "letsencrypt.org"`,
			},
			newTargetList: []string{
				`0 issue "letsencrypt.org"`,
				`0 issue "c.com"`,
				`0 issuewild "letsencrypt.org"`,
			},
			oldVal:         `0 issue "b.com"`,
			newVal:         `0 issue "c.com"`,
			expectSuppress: false,
		},
		// Backslash-quote form: the Akamai bind-zone renderer may produce
		// `0 issue \"ca.example.net\"` (literal `\` + `"` two-char sequence)
		// in the Go string even though the JSON layer uses standard `\"`.
		{
			name:           "CAA: backslash-quoted vs plain unquoted suppresses diff",
			oldTargetList:  []string{`0 issue \"ca.example.net\"`},
			newTargetList:  []string{`0 issue ca.example.net`},
			oldVal:         `0 issue \"ca.example.net\"`,
			newVal:         `0 issue ca.example.net`,
			expectSuppress: true,
		},
		{
			name:           "CAA: backslash-quoted vs plain-quoted suppresses diff",
			oldTargetList:  []string{`0 issue \"ca.example.net\"`},
			newTargetList:  []string{`0 issue "ca.example.net"`},
			oldVal:         `0 issue \"ca.example.net\"`,
			newVal:         `0 issue "ca.example.net"`,
			expectSuppress: true,
		},
		{
			// Different-length list: backslash-quoted old element vs plain new.
			name: "CAA: backslash-quoted element suppressed when list grows",
			oldTargetList: []string{
				`0 issue \"ca.example.net\"`,
			},
			newTargetList: []string{
				`0 issue "ca.example.net"`,
				`0 issuewild "ca.example.net"`,
			},
			oldVal:         `0 issue \"ca.example.net\"`,
			newVal:         `0 issue "ca.example.net"`,
			expectSuppress: true,
		},
	}

	for _, tc := range tests {
		t.Run(tc.name, func(t *testing.T) {
			t.Parallel()
			result := diffQuotedDNSRecord(tc.oldTargetList, tc.newTargetList, tc.oldVal, tc.newVal, RRTypeCaa, logger)
			assert.Equal(t, tc.expectSuppress, result)
		})
	}
}

func TestDiffQuotedDNSRecordSOA(t *testing.T) {
	t.Parallel()
	logger := akalog.NOPLogger()

	// SOA has no dedicated branch and falls through to the generic sorted-list
	// path. SOA records have exactly one target in practice, but the suppress
	// function still needs to correctly detect value changes.
	const (
		soaBase    = "ns1.example.com. hostmaster.example.com. 2024010101 3600 900 604800 300"
		soaUpdated = "ns1.example.com. hostmaster.example.com. 2024020101 3600 900 604800 300"
	)

	tests := []struct {
		name           string
		oldTargetList  []string
		newTargetList  []string
		oldVal         string
		newVal         string
		expectSuppress bool
	}{
		{
			name:           "SOA: identical value suppresses diff",
			oldTargetList:  []string{soaBase},
			newTargetList:  []string{soaBase},
			oldVal:         soaBase,
			newVal:         soaBase,
			expectSuppress: true,
		},
		{
			name:           "SOA: updated serial triggers diff",
			oldTargetList:  []string{soaBase},
			newTargetList:  []string{soaUpdated},
			oldVal:         soaBase,
			newVal:         soaUpdated,
			expectSuppress: false,
		},
	}

	for _, tc := range tests {
		t.Run(tc.name, func(t *testing.T) {
			t.Parallel()
			result := diffQuotedDNSRecord(tc.oldTargetList, tc.newTargetList, tc.oldVal, tc.newVal, RRTypeSoa, logger)
			assert.Equal(t, tc.expectSuppress, result)
		})
	}
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
