package dns

import (
	"errors"
	"log"
	"regexp"
	"testing"

	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/resource"
	"github.com/hashicorp/terraform-plugin-sdk/v2/terraform"
	"github.com/stretchr/testify/mock"
)

func TestAccDataSourceAuthoritiesSet_basic(t *testing.T) {
	t.Run("basic", func(t *testing.T) {
		client := &mockdns{}

		client.On("GetNameServerRecordList",
			mock.Anything, // ctx is irrelevant for this test
			mock.AnythingOfType("string"),
		).Return([]string{}, nil)

		dataSourceName := "data.akamai_authorities_set.test"

		useClient(client, func() {
			resource.UnitTest(t, resource.TestCase{
				PreCheck:     func() { testAccPreCheck(t) },
				Providers:    testAccProviders,
				CheckDestroy: testAccCheckAuthoritiesSetDestroy,
				Steps: []resource.TestStep{
					{
						Config: loadFixtureString("testdata/TestDataSetAuthorities/basic.tf"),
						Check: resource.ComposeTestCheckFunc(
							resource.TestCheckResourceAttrSet(dataSourceName, "id"),
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
				PreCheck:     func() { testAccPreCheck(t) },
				Providers:    testAccProviders,
				CheckDestroy: testAccCheckAuthoritiesSetDestroy,
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
				PreCheck:     func() { testAccPreCheck(t) },
				Providers:    testAccProviders,
				CheckDestroy: testAccCheckAuthoritiesSetDestroy,
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

func testAccCheckAuthoritiesSetDestroy(*terraform.State) error {
	log.Printf("[Group] Searching for AuthoritiesSet Delete skipped ")

	return nil
}
