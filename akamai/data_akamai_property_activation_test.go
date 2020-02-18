package akamai

import (
	"fmt"
	"testing"

	"github.com/hashicorp/terraform/helper/resource"
)

func TestAccDataSourcePropertyActivationComplete(t *testing.T) {
	dataSourceName := "akamai_property_activation.test"
	propertyName := "example-test"
	group := "group"

	config := testAccDataSourcePropertyActivationComplete(propertyName, group)
	resource.ParallelTest(t, resource.TestCase{
		PreCheck:     func() { testAccPreCheck(t) },
		Providers:    testAccProviders,
		CheckDestroy: testAccCheckAkamaiPropertyActivationDestroy,
		Steps: []resource.TestStep{
			{
				Config: config,
				Check: resource.ComposeTestCheckFunc(
					resource.TestCheckResourceAttrSet(dataSourceName, "property"),
					resource.TestCheckResourceAttrSet(dataSourceName, "id"),
					resource.TestCheckResourceAttr(dataSourceName, "version", "1"),
				),
			},
		},
	})
}

func testAccDataSourcePropertyActivationComplete(propertyName string, group string) string {
	t := testAccPropertyActivationComplete(propertyName, group)
	return fmt.Sprintf(`
%s

data "akamai_property_activation" "test" {
    property = "${akamai_property.property.name}"
}
`, t)
}
