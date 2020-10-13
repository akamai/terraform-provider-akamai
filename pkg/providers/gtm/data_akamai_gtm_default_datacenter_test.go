package gtm

import (
	"testing"

	gtm "github.com/akamai/AkamaiOPEN-edgegrid-golang/v2/pkg/configgtm"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/resource"
	"github.com/stretchr/testify/mock"
)

func TestAccDataSourceGTMDefaultDatacenter_basic(t *testing.T) {
	t.Run("basic", func(t *testing.T) {
		client := &mockgtm{}

		dc := gtm.Datacenter{}
		client.On("CreateMapsDefaultDatacenter",
			mock.Anything, // ctx is irrelevant for this test
			mock.AnythingOfType("string"),
		).Return(&dc, nil)

		dataSourceName := "data.akamai_gtm_default_datacenter.test"

		useClient(client, func() {
			resource.UnitTest(t, resource.TestCase{
				PreCheck:  func() { testAccPreCheck(t) },
				Providers: testAccProviders,
				Steps: []resource.TestStep{
					{
						Config: loadFixtureString("testdata/TestDataDefaultDatacenter/basic.tf"),
						Check: resource.ComposeTestCheckFunc(
							resource.TestCheckResourceAttrSet(dataSourceName, "id"),
						),
					},
				},
			})
		})

		client.AssertExpectations(t)
	})
}
