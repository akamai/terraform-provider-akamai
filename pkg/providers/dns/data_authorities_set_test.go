package dns

import (
	"errors"
	"regexp"
	"strings"
	"testing"

	"github.com/akamai/terraform-provider-akamai/v10/internal/edgegrid"
	"github.com/akamai/terraform-provider-akamai/v10/pkg/common/testutils"
	"github.com/hashicorp/terraform-plugin-testing/helper/resource"
	"github.com/stretchr/testify/mock"
)

func TestDataSourceAuthoritiesSet_basic(t *testing.T) {
	t.Parallel()
	t.Run("basic", func(t *testing.T) {
		t.Parallel()
		client := edgegrid.NewTestClient()

		dataSourceName := "data.akamai_authorities_set.test"
		outputName := "authorities"

		authorities := []string{"ns1.exampleterraform.io", "ns2.exampleterraform.io"}

		client.DNS.On("GetNameServerRecordList",
			testutils.MockContext,
			mock.AnythingOfType("dns.GetNameServerRecordListRequest"),
		).Return(authorities, nil)

		resource.UnitTest(t, resource.TestCase{
			ProtoV6ProviderFactories: testutils.NewTestProtoV6SDKProviderFactory(client, newSubproviderWithConfig(testSubproviderConfig())),
			Steps: []resource.TestStep{
				{
					Config: testutils.LoadFixtureString(t, "testdata/TestDataSetAuthorities/basic.tf"),
					Check: resource.ComposeTestCheckFunc(
						// check the values set in dataSourceAuthoritiesSetRead
						// authorities is an array that becomes authorities.0 and authorities.1 in tf state
						resource.TestCheckResourceAttrSet(dataSourceName, "authorities.0"),
						resource.TestCheckResourceAttrSet(dataSourceName, "authorities.1"),
						resource.TestCheckResourceAttrSet(dataSourceName, "id"),
						resource.TestCheckOutput(outputName, strings.Join(authorities, ",")),
					),
				},
			},
		})

		client.DNS.AssertExpectations(t)
	})

	t.Run("missing contract", func(t *testing.T) {
		t.Parallel()
		client := edgegrid.NewTestClient()

		resource.UnitTest(t, resource.TestCase{
			ProtoV6ProviderFactories: testutils.NewTestProtoV6SDKProviderFactory(client, newSubproviderWithConfig(testSubproviderConfig())),
			Steps: []resource.TestStep{
				{
					Config:      testutils.LoadFixtureString(t, "testdata/TestDataSetAuthorities/missing_contract.tf"),
					ExpectError: regexp.MustCompile(`Missing required argument`),
				},
			},
		})

		client.DNS.AssertExpectations(t)
	})

	t.Run("lookup error", func(t *testing.T) {
		t.Parallel()
		client := edgegrid.NewTestClient()

		client.DNS.On("GetNameServerRecordList",
			testutils.MockContext,
			mock.AnythingOfType("dns.GetNameServerRecordListRequest"),
		).Return(nil, errors.New("invalid contract"))

		resource.UnitTest(t, resource.TestCase{
			ProtoV6ProviderFactories: testutils.NewTestProtoV6SDKProviderFactory(client, newSubproviderWithConfig(testSubproviderConfig())),
			Steps: []resource.TestStep{
				{
					Config:      testutils.LoadFixtureString(t, "testdata/TestDataSetAuthorities/basic.tf"),
					ExpectError: regexp.MustCompile(`invalid contract`),
				},
			},
		})

		client.DNS.AssertExpectations(t)
	})
}
