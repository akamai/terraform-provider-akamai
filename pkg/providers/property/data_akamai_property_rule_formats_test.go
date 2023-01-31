package property

import (
	"testing"

	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/resource"
	"github.com/stretchr/testify/mock"

	"github.com/akamai/AkamaiOPEN-edgegrid-golang/v4/pkg/papi"
)

func Test_readPropertyRuleFormats(t *testing.T) {
	t.Run("get datasource property rule formats", func(t *testing.T) {
		client := &papi.Mock{}
		ruleFormats := papi.RuleFormatItems{
			Items: []string{
				"latest",
				"v2015-08-08"}}

		client.On("GetRuleFormats",
			mock.Anything,
		).Return(&papi.GetRuleFormatsResponse{RuleFormats: ruleFormats}, nil)
		useClient(client, nil, func() {
			resource.UnitTest(t, resource.TestCase{
				Providers: testAccProviders,
				Steps: []resource.TestStep{{
					Config: loadFixtureString("testdata/TestDSPropertyRuleFormats/rule_formats.tf"),
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
