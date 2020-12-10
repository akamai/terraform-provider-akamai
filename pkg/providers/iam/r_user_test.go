package iam

import (
	"testing"

	"github.com/akamai/terraform-provider-akamai/v2/pkg/test"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/resource"
)

func TestResUser(t *testing.T) {
	test.TODO(t, "Not implemented")
	prov := provider{}

	resource.UnitTest(t, resource.TestCase{
		ProviderFactories: prov.ProviderFactories(),
		Steps: []resource.TestStep{{
			Config: test.Fixture("%s/step0.tf", t.Name()),
		}},
	})
}
