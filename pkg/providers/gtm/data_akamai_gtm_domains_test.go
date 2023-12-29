package gtm

import (
	"fmt"
	"regexp"
	"testing"

	"github.com/akamai/AkamaiOPEN-edgegrid-golang/v7/pkg/gtm"
	"github.com/akamai/terraform-provider-akamai/v5/pkg/common/testutils"
	"github.com/hashicorp/terraform-plugin-testing/helper/resource"
	"github.com/stretchr/testify/mock"
)

func TestDataGtmDomains(t *testing.T) {
	tests := map[string]struct {
		givenTF                   string
		init                      func(*gtm.Mock)
		expectedAttributes        map[string]string
		expectedMissingAttributes []string
		expectError               *regexp.Regexp
	}{
		"success - response is ok": {
			givenTF: "valid.tf",
			init: func(m *gtm.Mock) {
				m.On("ListDomains", mock.Anything).Return([]*gtm.DomainItem{
					{
						Name:                 "test.dev.exp.cli.terraform.import.akadns.net",
						LastModified:         "2023-02-01T09:36:28.000+00:00",
						LastModifiedBy:       "dl-terraform-dev+5",
						ChangeID:             "ca4de6db-8d69-4980-8e2a-036b655f2e66",
						ActivationState:      "COMPLETE",
						ModificationComments: "Add AS Map New Map 1",
						SignAndServe:         false,
						Status:               "2023-02-01 09:47 GMT: Current configuration has been propagated to all GTM nameservers",
						AcgID:                "G-29RS4N8",
						Links: []*gtm.Link{{
							Rel:  "self",
							Href: "https://akaa-ouijhfns55qwgfuc-knsod5nrjl2w2gmt.luna-dev.akamaiapis.net/config-gtm/v1/domains/test.dev.exp.cli.terraform.import.akadns.net",
						},
						},
					},
					{
						Name:                 "devexpautomatedtest_rsh7a1.devexp.terraformtesting",
						LastModified:         "2023-12-21T08:34:31.463+00:00",
						LastModifiedBy:       "dl-terraform-dev+5",
						ChangeID:             "acca0158-398b-4a03-8886-81adc6328f56",
						ActivationState:      "COMPLETE",
						ModificationComments: "terraform test gtm domain",
						SignAndServe:         false,
						Status:               "2023-12-21 08:37 GMT: Current configuration has been propagated to all GTM nameservers",
						AcgID:                "G-29RS4N8",
						Links: []*gtm.Link{{
							Rel:  "self",
							Href: "https://akaa-ouijhfns55qwgfuc-knsod5nrjl2w2gmt.luna-dev.akamaiapis.net/config-gtm/v1/domains/devexpautomatedtest_rsh7a1.devexp.terraformtesting",
						},
						},
					},
					{
						Name:                 "devexpautomatedtest_dx4dfc.devexp.terraformtesting",
						LastModified:         "2023-12-22T08:43:47.553+00:00",
						LastModifiedBy:       "dl-terraform-dev+5",
						ChangeID:             "abf5b76f-f9de-4404-bb2c-9d15e7b9ff5d",
						ActivationState:      "COMPLETE",
						ModificationComments: "terraform test gtm domain",
						SignAndServe:         false,
						Status:               "2023-12-22 08:46 GMT: Current configuration has been propagated to all GTM nameservers",
						AcgID:                "G-29RS4N8",
						Links: []*gtm.Link{{
							Rel:  "self",
							Href: "https://akaa-ouijhfns55qwgfuc-knsod5nrjl2w2gmt.luna-dev.akamaiapis.net/config-gtm/v1/domains/devexpautomatedtest_dx4dfc.devexp.terraformtesting",
						},
						},
					},
				}, nil)
			},
			expectedAttributes: map[string]string{
				"domains.0.name":           "devexpautomatedtest_dx4dfc.devexp.terraformtesting",
				"domains.0.sign_and_serve": "false",
				"domains.0.links.0.href":   "https://akaa-ouijhfns55qwgfuc-knsod5nrjl2w2gmt.luna-dev.akamaiapis.net/config-gtm/v1/domains/devexpautomatedtest_dx4dfc.devexp.terraformtesting",
				"domains.0.links.0.rel":    "self",
				"domains.1.name":           "devexpautomatedtest_rsh7a1.devexp.terraformtesting",
				"domains.2.name":           "test.dev.exp.cli.terraform.import.akadns.net",
			},
			expectedMissingAttributes: nil,
			expectError:               nil,
		},
		"no domains found": {
			givenTF: "valid.tf",
			init: func(m *gtm.Mock) {
				m.On("ListDomains", mock.Anything).Return([]*gtm.DomainItem{}, nil)
			},
			expectedAttributes:        nil,
			expectedMissingAttributes: nil,
			expectError:               nil,
		},
		"error response from api": {
			givenTF: "valid.tf",
			init: func(m *gtm.Mock) {
				m.On("ListDomains", mock.Anything).Return(nil, fmt.Errorf("oops"))
			},
			expectError: regexp.MustCompile("oops"),
		},
	}

	for name, test := range tests {
		t.Run(name, func(t *testing.T) {
			client := &gtm.Mock{}
			if test.init != nil {
				test.init(client)
			}
			var checkFuncs []resource.TestCheckFunc
			for k, v := range test.expectedAttributes {
				checkFuncs = append(checkFuncs, resource.TestCheckResourceAttr("data.akamai_gtm_domains.domains", k, v))
			}
			for _, v := range test.expectedMissingAttributes {
				checkFuncs = append(checkFuncs, resource.TestCheckNoResourceAttr("data.akamai_gtm_domains.domains", v))
			}
			useClient(client, func() {
				resource.Test(t, resource.TestCase{
					IsUnitTest:               true,
					ProtoV5ProviderFactories: testAccProvidersProtoV5,
					Steps: []resource.TestStep{{
						Config:      testutils.LoadFixtureString(t, fmt.Sprintf("testdata/TestDataGtmDomains/%s", test.givenTF)),
						Check:       resource.ComposeAggregateTestCheckFunc(checkFuncs...),
						ExpectError: test.expectError,
					}},
				})
			})
			client.AssertExpectations(t)
		})
	}
}
