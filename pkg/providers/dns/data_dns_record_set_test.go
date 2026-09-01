package dns

import (
	"errors"
	"regexp"
	"testing"

	"github.com/akamai/AkamaiOPEN-edgegrid-golang/v13/pkg/dns"
	"github.com/akamai/terraform-provider-akamai/v11/internal/edgegrid"
	"github.com/akamai/terraform-provider-akamai/v11/pkg/common/testutils"
	"github.com/hashicorp/terraform-plugin-testing/helper/resource"
)

func TestDataSourceDNSRecordSet_basic(t *testing.T) {
	t.Parallel()

	t.Run("basic", func(t *testing.T) {
		t.Parallel()
		client := edgegrid.NewTestClient()

		dataSourceName := "data.akamai_dns_record_set.test"

		rdata := []string{"10.1.0.1", "10.2.0.1"}

		client.DNS.On("GetRdata",
			testutils.MockContext,
			dns.GetRdataRequest{Zone: "exampleterraform.io", Name: "www.exampleterraform.io", RecordType: "A"},
		).Return(rdata, nil)

		resource.UnitTest(t, resource.TestCase{
			ProtoV6ProviderFactories: testutils.NewTestProtoV6SDKProviderFactory(client, newSubproviderWithConfig(testSubproviderConfig())),
			Steps: []resource.TestStep{
				{
					Config: testutils.LoadFixtureString(t, "testdata/TestDataDnsRecordSet/basic.tf"),
					Check: resource.ComposeTestCheckFunc(
						// check the values set in dataSourceDNSRecordSetRead
						// rdata is an array that becomes rdata.0 and rdata.1 in tf state
						resource.TestCheckResourceAttr(dataSourceName, "rdata.0", "10.1.0.1"),
						resource.TestCheckResourceAttr(dataSourceName, "rdata.1", "10.2.0.1"),
						resource.TestCheckResourceAttr(dataSourceName, "id", "www.exampleterraform.io"),
					),
				},
			},
		})

		client.DNS.AssertExpectations(t)
	})

	t.Run("basic txt", func(t *testing.T) {
		t.Parallel()
		client := edgegrid.NewTestClient()

		dataSourceName := "data.akamai_dns_record_set.test"

		rdata := []string{"abc", "zxy", "hji"}

		client.DNS.On("GetRdata",
			testutils.MockContext,
			dns.GetRdataRequest{Zone: "exampleterraform.io", Name: "www.exampleterraform.io", RecordType: "TXT"},
		).Return(rdata, nil)

		resource.UnitTest(t, resource.TestCase{
			ProtoV6ProviderFactories: testutils.NewTestProtoV6SDKProviderFactory(client, newSubproviderWithConfig(testSubproviderConfig())),
			Steps: []resource.TestStep{
				{
					Config: testutils.LoadFixtureString(t, "testdata/TestDataDnsRecordSet/txt.tf"),
					Check: resource.ComposeTestCheckFunc(
						resource.TestCheckResourceAttr(dataSourceName, "rdata.0", "abc"),
						resource.TestCheckResourceAttr(dataSourceName, "rdata.1", "zxy"),
						resource.TestCheckResourceAttr(dataSourceName, "rdata.2", "hji"),
						resource.TestCheckResourceAttr(dataSourceName, "id", "www.exampleterraform.io"),
					),
				},
			},
		})

		client.DNS.AssertExpectations(t)
	})

	t.Run("error", func(t *testing.T) {
		t.Parallel()
		client := edgegrid.NewTestClient()

		client.DNS.On("GetRdata",
			testutils.MockContext,
			dns.GetRdataRequest{Zone: "exampleterraform.io", Name: "www.exampleterraform.io", RecordType: "A"},
		).Return(nil, errors.New("invalid zone"))

		resource.UnitTest(t, resource.TestCase{
			ProtoV6ProviderFactories: testutils.NewTestProtoV6SDKProviderFactory(client, newSubproviderWithConfig(testSubproviderConfig())),
			Steps: []resource.TestStep{
				{
					Config:      testutils.LoadFixtureString(t, "testdata/TestDataDnsRecordSet/basic.tf"),
					ExpectError: regexp.MustCompile(`invalid zone`),
				},
			},
		})

		client.DNS.AssertExpectations(t)
	})
}
