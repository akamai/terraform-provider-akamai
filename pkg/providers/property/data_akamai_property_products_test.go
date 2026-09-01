package property

import (
	"fmt"
	"regexp"
	"testing"

	"github.com/akamai/AkamaiOPEN-edgegrid-golang/v13/pkg/papi"
	"github.com/akamai/terraform-provider-akamai/v11/internal/edgegrid"
	"github.com/akamai/terraform-provider-akamai/v11/pkg/common/testutils"
	"github.com/hashicorp/terraform-plugin-testing/helper/resource"
)

func TestVerifyProductsDataSourceSchema(t *testing.T) {
	t.Parallel()
	t.Run("akamai_property_products - test data source required contract", func(t *testing.T) {
		resource.UnitTest(t, resource.TestCase{
			ProtoV6ProviderFactories: testutils.NewProtoV6ProviderFactory(NewSubprovider()),
			IsUnitTest:               true,
			Steps: []resource.TestStep{{
				Config:      testConfig(""),
				ExpectError: regexp.MustCompile("The argument \"contract_id\" is required, but no definition was found"),
			}},
		})
	})
}

func TestOutputProductsDataSource(t *testing.T) {
	t.Parallel()
	t.Run("akamai_property_products - input OK - output OK", func(t *testing.T) {
		client := edgegrid.NewTestClient()
		client.PAPI.On("GetProducts", testutils.MockContext, testutils.MockContext).Return(&papi.GetProductsResponse{
			AccountID:  "act_anyAccount",
			ContractID: "ctr_AnyContract",
			Products: papi.ProductsItems{
				Items: []papi.ProductItem{{ProductName: "anyProduct", ProductID: "prd_anyProduct"}},
			},
		}, nil)

		resource.UnitTest(t, resource.TestCase{
			ProtoV6ProviderFactories: testutils.NewTestProtoV6SDKProviderFactory(client, NewSubprovider()),
			IsUnitTest:               true,
			Steps: []resource.TestStep{{
				Config: testConfig("contract_id = \"ctr_test\""),
				Check: resource.ComposeTestCheckFunc(
					resource.TestCheckOutput("product_name0", "anyProduct"),
					resource.TestCheckOutput("product_id0", "prd_anyProduct"),
				),
			}},
		})
		client.PAPI.AssertExpectations(t)
	})
}

func testConfig(contractIDConfig string) string {
	return fmt.Sprintf(`
	provider "akamai" {
		edgerc = "../../common/testutils/edgerc"
	}

	data "akamai_property_products" "example" { %s }

    output "product_name0" {
		value = "${data.akamai_property_products.example.products[0].product_name}"
	}

    output "product_id0" {
		value = "${data.akamai_property_products.example.products[0].product_id}"
	}
`, contractIDConfig)
}
