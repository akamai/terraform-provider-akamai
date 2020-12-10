package property

import (
	"fmt"
	"regexp"
	"testing"

	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/resource"
	"github.com/stretchr/testify/mock"

	"github.com/akamai/AkamaiOPEN-edgegrid-golang/v2/pkg/papi"
)

func Test_readPropertyRuleFormats(t *testing.T) {
	tests := map[string]struct {
		givenTF            string
		init               func(*mockpapi)
		expectedAttributes map[string]string
		withError          *regexp.Regexp
	}{
		"get datasource property rule formats": {
			givenTF: "rule_formats.tf",
			init: func(m *mockpapi) {
				m.On("GetRuleFormats", mock.Anything).Return(&papi.GetRuleFormatsResponse{
					RuleFormats: papi.RuleFormatItems{
						Items: []string{
							"v2020-11-02",
							"v2020-03-04",
							"v2019-07-25",
							"v2018-09-12",
							"v2018-02-27",
							"v2017-06-19",
							"v2016-11-15",
							"v2015-08-17",
							"latest",
						},
					},
				}, nil).Once()
				m.On("GetRuleFormats", mock.Anything).Return(&papi.GetRuleFormatsResponse{
					RuleFormats: papi.RuleFormatItems{
						Items: []string{
							"v2020-11-02",
							"v2020-03-04",
							"v2019-07-25",
							"v2018-09-12",
							"v2018-02-27",
							"v2017-06-19",
							"v2016-11-15",
							"v2015-08-17",
							"latest",
						},
					},
				}, nil).Once()
			},
			expectedAttributes: map[string]string{},
		},
	}
	for name, test := range tests {
		t.Run(name, func(t *testing.T) {
			client := &mockpapi{}
			test.init(client)
			var checkFuncs []resource.TestCheckFunc
			for k, v := range test.expectedAttributes {
				checkFuncs = append(checkFuncs, resource.TestCheckResourceAttr("akamai_property_rules.getruleformats", k, v))
			}
			useClient(client, func() {
				resource.UnitTest(t, resource.TestCase{
					Providers: testAccProviders,
					Steps: []resource.TestStep{
						{
							ExpectNonEmptyPlan: true,
							Config:             loadFixtureString(fmt.Sprintf("testdata/TestDSPropertyRuleFormats/%s", test.givenTF)),
							Check:              resource.ComposeAggregateTestCheckFunc(checkFuncs...),
							ExpectError:        test.withError,
						},
					},
				})
			})
			client.AssertExpectations(t)
		})
	}
}
