package gtm

import (
	"fmt"
	"regexp"
	"testing"

	"github.com/akamai/AkamaiOPEN-edgegrid-golang/v9/pkg/gtm"
	"github.com/akamai/terraform-provider-akamai/v6/pkg/common/testutils"
	"github.com/hashicorp/terraform-plugin-testing/helper/resource"
	"github.com/stretchr/testify/mock"
)

func TestDataGTMDomains(t *testing.T) {
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
				mockListDomains(m, []gtm.DomainItem{
					{
						Name:                 "test1.terraformtesting.net",
						LastModified:         "2023-02-01T09:36:28.000+00:00",
						LastModifiedBy:       "test-user",
						ChangeID:             "ca4de6db-8d69-4980-8e2a-036b655f2e66",
						ActivationState:      "COMPLETE",
						ModificationComments: "Add AS Map New Map 1",
						SignAndServe:         false,
						Status:               "2023-02-01 09:47 GMT: Current configuration has been propagated to all GTM nameservers",
						AcgID:                "TestACGID-1",
						Links: []gtm.Link{{
							Rel:  "self",
							Href: "https://test-domain.net/config-gtm/v1/domains/test1.terraformtesting.net",
						},
						},
					},
					{
						Name:                 "test2.terraformtesting.net",
						LastModified:         "2023-12-21T08:34:31.463+00:00",
						LastModifiedBy:       "test-user",
						ChangeID:             "acca0158-398b-4a03-8886-81adc6328f56",
						ActivationState:      "COMPLETE",
						ModificationComments: "terraform test gtm domain",
						SignAndServe:         false,
						Status:               "2023-12-21 08:37 GMT: Current configuration has been propagated to all GTM nameservers",
						AcgID:                "TestACGID-1",
						Links: []gtm.Link{{
							Rel:  "self",
							Href: "https://test-domain.net/config-gtm/v1/domains/test2.terraformtesting.net",
						},
						},
					},
					{
						Name:                 "test3.terraformtesting.net",
						LastModified:         "2023-12-22T08:43:47.553+00:00",
						LastModifiedBy:       "test-user",
						ChangeID:             "abf5b76f-f9de-4404-bb2c-9d15e7b9ff5d",
						ActivationState:      "COMPLETE",
						ModificationComments: "terraform test gtm domain",
						SignAndServe:         false,
						Status:               "2023-12-22 08:46 GMT: Current configuration has been propagated to all GTM nameservers",
						AcgID:                "TestACGID-1",
						Links: []gtm.Link{{
							Rel:  "self",
							Href: "https://test-domain.net/config-gtm/v1/domains/test3.terraformtesting.net",
						},
						},
					},
				}, nil, 3)
			},
			expectedAttributes: map[string]string{
				"domains.0.name":           "test3.terraformtesting.net",
				"domains.0.sign_and_serve": "false",
				"domains.0.links.0.href":   "https://test-domain.net/config-gtm/v1/domains/test3.terraformtesting.net",
				"domains.0.links.0.rel":    "self",
				"domains.1.name":           "test2.terraformtesting.net",
				"domains.2.name":           "test1.terraformtesting.net",
			},
		},
		"no domains found": {
			givenTF: "valid.tf",
			init: func(m *gtm.Mock) {
				mockListDomains(m, []gtm.DomainItem{}, nil, 3)
			},
		},
		"error response from api": {
			givenTF: "valid.tf",
			init: func(m *gtm.Mock) {
				mockListDomains(m, nil, fmt.Errorf("oops"), 1)
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
			const datasourceName = "data.akamai_gtm_domains.domains"
			for k, v := range test.expectedAttributes {
				checkFuncs = append(checkFuncs, resource.TestCheckResourceAttr(datasourceName, k, v))
			}
			for _, v := range test.expectedMissingAttributes {
				checkFuncs = append(checkFuncs, resource.TestCheckNoResourceAttr(datasourceName, v))
			}
			useClient(client, func() {
				resource.Test(t, resource.TestCase{
					IsUnitTest:               true,
					ProtoV6ProviderFactories: testutils.NewProtoV6ProviderFactory(NewSubprovider()),
					Steps: []resource.TestStep{{
						Config:      testutils.LoadFixtureStringf(t, "testdata/TestDataGtmDomains/%s", test.givenTF),
						Check:       resource.ComposeAggregateTestCheckFunc(checkFuncs...),
						ExpectError: test.expectError,
					}},
				})
			})
			client.AssertExpectations(t)
		})
	}
}

func mockListDomains(client *gtm.Mock, resp []gtm.DomainItem, err error, times int) *mock.Call {
	return client.On("ListDomains", testutils.MockContext).Return(resp, err).Times(times)
}
