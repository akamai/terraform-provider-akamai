package gtm

import (
	"testing"

	"github.com/akamai/AkamaiOPEN-edgegrid-golang/v9/pkg/gtm"
	"github.com/akamai/terraform-provider-akamai/v6/pkg/common/testutils"
	"github.com/hashicorp/terraform-plugin-testing/helper/resource"
)

func TestAccDataSourceGTMDefaultDatacenter_basic(t *testing.T) {
	t.Run("basic", func(t *testing.T) {
		client := &gtm.Mock{}

		dc := gtm.Datacenter{
			DatacenterID: 1000,
		}

		client.On("CreateMapsDefaultDatacenter",
			testutils.MockContext,
			"testdomain.net",
		).Return(&dc, nil)

		dataSourceName := "data.akamai_gtm_default_datacenter.test"

		useClient(client, func() {
			resource.UnitTest(t, resource.TestCase{
				ProtoV6ProviderFactories: testutils.NewProtoV6ProviderFactory(NewSubprovider()),
				Steps: []resource.TestStep{
					{
						Config: testutils.LoadFixtureString(t, "testdata/TestDataDefaultDatacenter/basic.tf"),
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
