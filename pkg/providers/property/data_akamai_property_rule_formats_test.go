package property

import (
	"testing"

	"github.com/akamai/AkamaiOPEN-edgegrid-golang/v10/pkg/papi"
	"github.com/akamai/terraform-provider-akamai/v7/pkg/common/testutils"
	"github.com/hashicorp/terraform-plugin-testing/helper/resource"
)

func Test_readPropertyRuleFormats(t *testing.T) {
	t.Run("get datasource property rule formats", func(t *testing.T) {
		client := &papi.Mock{}
		ruleFormats := papi.RuleFormatItems{
			Items: []string{
				"latest",
				"v2015-08-08"}}

		client.On("GetRuleFormats",
			testutils.MockContext,
		).Return(&papi.GetRuleFormatsResponse{RuleFormats: ruleFormats}, nil)
		useClient(client, nil, func() {
			resource.UnitTest(t, resource.TestCase{
				ProtoV6ProviderFactories: testutils.NewProtoV6ProviderFactory(NewSubprovider()),
				Steps: []resource.TestStep{{
					Config: testutils.LoadFixtureString(t, "testdata/TestDSPropertyRuleFormats/rule_formats.tf"),
					Check: resource.ComposeAggregateTestCheckFunc(
						resource.TestCheckResourceAttr("data.akamai_property_rule_formats.akarulesformats", "id", "rule_format"),
						resource.TestCheckResourceAttr("data.akamai_property_rule_formats.akarulesformats", "rule_format.0", "latest"),
						resource.TestCheckResourceAttr("data.akamai_property_rule_formats.akarulesformats", "rule_format.1", "v2015-08-08"),
					),
				}},
			})
		})

		client.AssertExpectations(t)
	})
}
