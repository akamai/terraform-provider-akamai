package gtm

import (
	"net/http"
	"testing"
	"log"

	gtm "github.com/akamai/AkamaiOPEN-edgegrid-golang/v2/pkg/configgtm"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/resource"
	"github.com/stretchr/testify/mock"
)

var gtmTestDomain = "gtm_terra_testdomain.akadns.net"
var contract = "1-2ABCDEF"
var group = "123ABC"

func TestGtmDomainCreate(t *testing.T) {

	dom := &gtm.Domain{
		Name:                    gtmTestDomain,
		Type:                    "weighted",
		LoadImbalancePercentage: 10,
	}

	t.Run("create domain", func(t *testing.T) {
		client := &mockgtm{}

		client.On("GetDomain",
			mock.Anything, // ctx is irrelevant for this test
			gtmTestDomain,
		).Return(nil, &gtm.Error{
			StatusCode: http.StatusNotFound,
		}).Once().Run(func(mock.Arguments) {
			client.On("GetDomain",
				mock.Anything, // ctx is irrelevant for this test
				gtmTestDomain,
			).Return(dom, nil)
		})
		dr := gtm.DomainResponse{}
		client.On("CreateDomain",
			mock.Anything, // ctx is irrelevant for this test
			mock.AnythingOfType("*gtm.Domain"),
			mock.AnythingOfType("map[string]string"),
		).Return(&dr, nil)

		rs := gtm.ResponseStatus{}
		client.On("DeleteDomain",
			mock.Anything, // ctx is irrelevant for this test
			mock.AnythingOfType("*gtm.Domain"),
		).Return(&rs, nil)

		dataSourceName := "akamai_gtm_domain.testdomain"

		useClient(client, func() {
			resource.UnitTest(t, resource.TestCase{
				PreCheck:  func() { testAccPreCheck(t) },
				Providers: testAccProviders,
				Steps: []resource.TestStep{
					{
						Config: loadFixtureString("testdata/TestResGtmDomain/create_basic.tf"),
						Check: resource.ComposeTestCheckFunc(
							resource.TestCheckResourceAttr(dataSourceName, "name", gtmTestDomain),
							resource.TestCheckResourceAttr(dataSourceName, "type", "weighted"),
							//resource.TestCheckResourceAttr(dataSourceName, "load_imbalance_percentage", 10),
						),
					},
				},
			})
		})

		client.AssertExpectations(t)
	})

	t.Run("update domain", func(t *testing.T) {
		client := &mockgtm{}

		client.On("GetDomain",
			mock.Anything, // ctx is irrelevant for this test
			gtmTestDomain,
		).Return(dom, nil)

		rs := gtm.ResponseStatus{}
		client.On("UpdateDomain",
			mock.Anything, // ctx is irrelevant for this test
			mock.AnythingOfType("*gtm.Domain"),
                        mock.AnythingOfType("map[string]string"),
		).Return(&rs, nil)

		rs = gtm.ResponseStatus{}
		client.On("DeleteDomain",
			mock.Anything, // ctx is irrelevant for this test
			mock.AnythingOfType("*gtm.Domain"),
		).Return(&rs, nil)

		dataSourceName := "akamai_gtm_domain.testdomain"

		useClient(client, func() {
			resource.UnitTest(t, resource.TestCase{
				PreCheck:  func() { testAccPreCheck(t) },
				Providers: testAccProviders,
				Steps: []resource.TestStep{
					{
						Config: loadFixtureString("testdata/TestResDnsDomain/create_basic.tf"),
						Check: resource.ComposeTestCheckFunc(
							resource.TestCheckResourceAttr(dataSourceName, "name", gtmTestDomain),
							resource.TestCheckResourceAttr(dataSourceName, "type", "weighted"),
							//resource.TestCheckResourceAttr(dataSourceName, "load_imbalance_percentage", 10),
						),
					},
				},
			})
		})

		client.AssertExpectations(t)
	})
}

// Sets a Hack flag so cn work with existing Domains (only Admin can Delete)
func testAccPreCheckTF(_ *testing.T) {

	// by definition, we are running acceptance tests. ;-)
	log.Printf("[DEBUG] [Akamai GTMV1] Setting HashiAcc true")
	HashiAcc = true

}
