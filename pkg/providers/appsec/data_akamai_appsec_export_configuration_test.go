package appsec

import (
	"encoding/json"
	"testing"

	"github.com/akamai/AkamaiOPEN-edgegrid-golang/v3/pkg/appsec"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/resource"
	"github.com/stretchr/testify/mock"
	"github.com/stretchr/testify/require"
)

func TestAkamaiExportConfiguration_data_basic(t *testing.T) {
	t.Run("Configuration Export Tests", func(t *testing.T) {
		client := &appsec.Mock{}

		getExportConfigurationResponse := appsec.GetExportConfigurationResponse{}
		err := json.Unmarshal(loadFixtureBytes("testdata/TestDSExportConfiguration/ExportConfiguration.json"), &getExportConfigurationResponse)
		require.NoError(t, err)

		client.On("GetExportConfiguration",
			mock.Anything,
			appsec.GetExportConfigurationRequest{ConfigID: 43253, Version: 7},
		).Return(&getExportConfigurationResponse, nil)

		expectedEvalGroups := "\n \n// terraform import akamai_appsec_eval_group.akamai_appsec_eval_group_AAAA_81230 43253:AAAA_81230:POLICY\nresource \"akamai_appsec_eval_group\" \"akamai_appsec_eval_group_AAAA_81230\" { \n  config_id = 43253\n  security_policy_id = \"AAAA_81230\" \n  attack_group = \"POLICY\" \n  attack_group_action = \"alert\"\n  condition_exception = <<-EOF\n {\"exception\":{\"specificHeaderCookieParamXmlOrJsonNames\":[{\"names\":[\"ASE-Manual-Active-COOKIES\"],\"selector\":\"REQUEST_COOKIES\",\"wildcard\":true}]}}  \n \n EOF \n \n}\n"

		useClient(client, func() {
			resource.Test(t, resource.TestCase{
				IsUnitTest: true,
				Providers:  testAccProviders,
				Steps: []resource.TestStep{
					{
						Config: loadFixtureString("testdata/TestDSExportConfiguration/match_by_id.tf"),
						Check: resource.ComposeAggregateTestCheckFunc(
							resource.TestCheckResourceAttr("data.akamai_appsec_export_configuration.test", "id", "43253"),
							resource.TestCheckResourceAttr("data.akamai_appsec_export_configuration.evalGroups", "output_text", expectedEvalGroups),
						),
					},
				},
			})
		})

		client.AssertExpectations(t)
	})

}
