package dns

import (
	"errors"
	"regexp"
	"strings"
	"testing"

	"github.com/akamai/AkamaiOPEN-edgegrid-golang/v7/pkg/dns"
	"github.com/akamai/terraform-provider-akamai/v5/pkg/common/testutils"
	"github.com/hashicorp/terraform-plugin-testing/helper/resource"
	"github.com/stretchr/testify/mock"
)

func TestDataSourceAuthoritiesSet_basic(t *testing.T) {
	t.Run("basic", func(t *testing.T) {
		client := &dns.Mock{}

		dataSourceName := "data.akamai_authorities_set.test"
		outputName := "authorities"

		authorities := []string{"ns1.exampleterraform.io", "ns2.exampleterraform.io"}

		client.On("GetNameServerRecordList",
			mock.Anything, // ctx is irrelevant for this test
			mock.AnythingOfType("string"),
		).Return(authorities, nil)

		useClient(client, func() {
			resource.UnitTest(t, resource.TestCase{
				ProviderFactories: testAccProviders,
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
		})

		client.AssertExpectations(t)
	})

	t.Run("missing contract", func(t *testing.T) {
		client := &dns.Mock{}

		useClient(client, func() {
			resource.UnitTest(t, resource.TestCase{
				ProviderFactories: testAccProviders,
				Steps: []resource.TestStep{
					{
						Config:      testutils.LoadFixtureString(t, "testdata/TestDataSetAuthorities/missing_contract.tf"),
						ExpectError: regexp.MustCompile(`Missing required argument`),
					},
				},
			})
		})

		client.AssertExpectations(t)
	})

	t.Run("lookup error", func(t *testing.T) {
		client := &dns.Mock{}

		client.On("GetNameServerRecordList",
			mock.Anything, // ctx is irrelevant for this test
			mock.AnythingOfType("string"),
		).Return(nil, errors.New("invalid contract"))

		useClient(client, func() {
			resource.UnitTest(t, resource.TestCase{
				ProviderFactories: testAccProviders,
				Steps: []resource.TestStep{
					{
						Config:      testutils.LoadFixtureString(t, "testdata/TestDataSetAuthorities/basic.tf"),
						ExpectError: regexp.MustCompile(`invalid contract`),
					},
				},
			})
		})

		client.AssertExpectations(t)
	})
}
