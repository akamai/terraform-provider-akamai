package property

import (
	"testing"

	"github.com/akamai/AkamaiOPEN-edgegrid-golang/v13/pkg/papi"
	"github.com/akamai/terraform-provider-akamai/v11/internal/edgegrid"
	"github.com/akamai/terraform-provider-akamai/v11/pkg/common/testutils"
	"github.com/hashicorp/terraform-plugin-testing/helper/resource"
)

func Test_readPropertyRuleFormats(t *testing.T) {
	t.Parallel()
	t.Run("get datasource property rule formats", func(t *testing.T) {
		client := edgegrid.NewTestClient()
		ruleFormats := papi.RuleFormatItems{
			Items: []string{
				"latest",
				"v2015-08-08"}}

		client.PAPI.On("GetRuleFormats",
			testutils.MockContext,
		).Return(&papi.GetRuleFormatsResponse{RuleFormats: ruleFormats}, nil)
		resource.UnitTest(t, resource.TestCase{
			ProtoV6ProviderFactories: testutils.NewTestProtoV6SDKProviderFactory(client, NewSubprovider()),
			Steps: []resource.TestStep{{
				Config: testutils.LoadFixtureString(t, "testdata/TestDSPropertyRuleFormats/rule_formats.tf"),
				Check: resource.ComposeAggregateTestCheckFunc(
					resource.TestCheckResourceAttr("data.akamai_property_rule_formats.akarulesformats", "id", "rule_format"),
					resource.TestCheckResourceAttr("data.akamai_property_rule_formats.akarulesformats", "rule_format.0", "latest"),
					resource.TestCheckResourceAttr("data.akamai_property_rule_formats.akarulesformats", "rule_format.1", "v2015-08-08"),
				),
			}},
		})
		client.PAPI.AssertExpectations(t)
	})
}
