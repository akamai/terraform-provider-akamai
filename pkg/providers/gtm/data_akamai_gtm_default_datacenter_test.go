package gtm

import (
	"testing"

	"github.com/akamai/AkamaiOPEN-edgegrid-golang/v13/pkg/gtm"
	"github.com/akamai/terraform-provider-akamai/v11/internal/edgegrid"
	"github.com/akamai/terraform-provider-akamai/v11/pkg/common/testutils"
	"github.com/hashicorp/terraform-plugin-testing/helper/resource"
	"github.com/stretchr/testify/mock"
)

func TestAccDataSourceGTMDefaultDatacenter_basic(t *testing.T) {
	t.Parallel()
	t.Run("basic", func(t *testing.T) {
		t.Parallel()
		client := edgegrid.NewTestClient()

		dc := gtm.Datacenter{
			DatacenterID: 1000,
		}

		mockCreateMapsDefaultDatacenter(client.GTM, &dc, testutils.ThreeTimes)

		dataSourceName := "data.akamai_gtm_default_datacenter.test"

		resource.UnitTest(t, resource.TestCase{
			ProtoV6ProviderFactories: testutils.NewTestProtoV6SDKProviderFactory(client, NewSubprovider()),
			Steps: []resource.TestStep{
				{
					Config: testutils.LoadFixtureString(t, "testdata/TestDataDefaultDatacenter/basic.tf"),
					Check: resource.ComposeTestCheckFunc(
						resource.TestCheckResourceAttrSet(dataSourceName, "id"),
					),
				},
			},
		})

		client.GTM.AssertExpectations(t)
	})
}

func mockCreateMapsDefaultDatacenter(client *gtm.Mock, dc *gtm.Datacenter, times int) *mock.Call {
	return client.On("CreateMapsDefaultDatacenter",
		testutils.MockContext,
		"testdomain.net",
	).Return(dc, nil).Times(times)
}
