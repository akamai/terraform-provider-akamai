package dns

import (
	"errors"
	"regexp"
	"strings"
	"testing"

	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/resource"
	"github.com/stretchr/testify/mock"
)

func TestDataSourceAuthoritiesSet_basic(t *testing.T) {
	t.Run("basic", func(t *testing.T) {
		client := &mockdns{}

		dataSourceName := "data.akamai_authorities_set.test"
		outputName := "authorities"

		authorities := []string{"ns1.exampleterraform.io", "ns2.exampleterraform.io"}

		client.On("GetNameServerRecordList",
			mock.Anything, // ctx is irrelevant for this test
			mock.AnythingOfType("string"),
		).Return(authorities, nil)

		useClient(client, func() {
			resource.UnitTest(t, resource.TestCase{
				PreCheck:  func() { testAccPreCheck(t) },
				Providers: testAccProviders,
				Steps: []resource.TestStep{
					{
						Config: loadFixtureString("testdata/TestDataSetAuthorities/basic.tf"),
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
		client := &mockdns{}

		useClient(client, func() {
			resource.UnitTest(t, resource.TestCase{
				PreCheck:  func() { testAccPreCheck(t) },
				Providers: testAccProviders,
				Steps: []resource.TestStep{
					{
						Config:      loadFixtureString("testdata/TestDataSetAuthorities/missing_contract.tf"),
						ExpectError: regexp.MustCompile(`Missing required argument`),
					},
				},
			})
		})

		client.AssertExpectations(t)
	})

	t.Run("lookup error", func(t *testing.T) {
		client := &mockdns{}

		client.On("GetNameServerRecordList",
			mock.Anything, // ctx is irrelevant for this test
			mock.AnythingOfType("string"),
		).Return(nil, errors.New("invalid contract"))

		useClient(client, func() {
			resource.UnitTest(t, resource.TestCase{
				PreCheck:  func() { testAccPreCheck(t) },
				Providers: testAccProviders,
				Steps: []resource.TestStep{
					{
						Config:      loadFixtureString("testdata/TestDataSetAuthorities/basic.tf"),
						ExpectError: regexp.MustCompile(`invalid contract`),
					},
				},
			})
		})

		client.AssertExpectations(t)
	})
}
