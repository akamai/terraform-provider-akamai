package gtm

import (
	"testing"

	"github.com/akamai/AkamaiOPEN-edgegrid-golang/v3/pkg/gtm"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/resource"
	"github.com/stretchr/testify/mock"
)

func TestAccDataSourceGTMDefaultDatacenter_basic(t *testing.T) {
	t.Run("basic", func(t *testing.T) {
		client := &mockgtm{}

		dc := gtm.Datacenter{
			DatacenterId: 1000,
		}

		client.On("CreateMapsDefaultDatacenter",
			mock.Anything, // ctx is irrelevant for this test
			mock.AnythingOfType("string"),
		).Return(&dc, nil)

		client.On("NewDatacenter",
			mock.Anything, // ctx is irrelevant for this test
		).Return(&dc)

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
