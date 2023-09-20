package appsec

import (
	"encoding/json"
	"testing"

	"github.com/akamai/AkamaiOPEN-edgegrid-golang/v7/pkg/appsec"
	"github.com/akamai/terraform-provider-akamai/v5/pkg/common/testutils"
	"github.com/hashicorp/terraform-plugin-testing/helper/resource"
	"github.com/stretchr/testify/mock"
	"github.com/stretchr/testify/require"
)

func TestAkamaiSiemDefinitions_data_basic(t *testing.T) {
	t.Run("match by SiemDefinitions ID", func(t *testing.T) {
		client := &appsec.Mock{}

		getSiemDefinitionsResponse := appsec.GetSiemDefinitionsResponse{}
		err := json.Unmarshal(testutils.LoadFixtureBytes(t, "testdata/TestDSSiemDefinitions/SiemDefinitions.json"), &getSiemDefinitionsResponse)
		require.NoError(t, err)

		client.On("GetSiemDefinitions",
			mock.Anything,
			appsec.GetSiemDefinitionsRequest{ID: 0, SiemDefinitionName: "SIEM Version 01"},
		).Return(&getSiemDefinitionsResponse, nil)

		useClient(client, func() {
			resource.Test(t, resource.TestCase{
				IsUnitTest:               true,
				ProtoV6ProviderFactories: testAccProviders,
				Steps: []resource.TestStep{
					{
						Config: testutils.LoadFixtureString(t, "testdata/TestDSSiemDefinitions/match_by_id.tf"),
						Check: resource.ComposeAggregateTestCheckFunc(
							resource.TestCheckResourceAttr("data.akamai_appsec_siem_definitions.test", "id", "1"),
						),
					},
				},
			})
		})

		client.AssertExpectations(t)
	})

}
