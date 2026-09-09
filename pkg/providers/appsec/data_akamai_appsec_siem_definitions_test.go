package appsec

import (
	"encoding/json"
	"testing"

	"github.com/akamai/AkamaiOPEN-edgegrid-golang/v14/pkg/appsec"
	"github.com/akamai/terraform-provider-akamai/v11/internal/edgegrid"
	"github.com/akamai/terraform-provider-akamai/v11/pkg/common/testutils"
	"github.com/hashicorp/terraform-plugin-testing/helper/resource"
	"github.com/stretchr/testify/require"
)

func TestAkamaiSiemDefinitions_data_basic(t *testing.T) {
	t.Parallel()
	t.Run("match by SiemDefinitions ID", func(t *testing.T) {
		t.Parallel()
		client := edgegrid.NewTestClient()

		getSiemDefinitionsResponse := appsec.GetSiemDefinitionsResponse{}
		err := json.Unmarshal(testutils.LoadFixtureBytes(t, "testdata/TestDSSiemDefinitions/SiemDefinitions.json"), &getSiemDefinitionsResponse)
		require.NoError(t, err)

		client.APPSEC.On("GetSiemDefinitions",
			testutils.MockContext,
			appsec.GetSiemDefinitionsRequest{ID: 0, SiemDefinitionName: "SIEM Version 01"},
		).Return(&getSiemDefinitionsResponse, nil)

		mockGetConfigurationVersionDefault(client.APPSEC)
		resource.UnitTest(t, resource.TestCase{
			ProtoV6ProviderFactories: testutils.NewTestProtoV6ProviderFactory(client, NewSubprovider()),
			Steps: []resource.TestStep{
				{
					Config: testutils.LoadFixtureString(t, "testdata/TestDSSiemDefinitions/match_by_id.tf"),
					Check: resource.ComposeAggregateTestCheckFunc(
						resource.TestCheckResourceAttr("data.akamai_appsec_siem_definitions.test", "id", "1"),
					),
				},
			},
		})

		client.APPSEC.AssertExpectations(t)
	})

}
