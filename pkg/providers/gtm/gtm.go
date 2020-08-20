// +build all gtm

package gtm

import "github.com/akamai/terraform-provider-akamai/v2/pkg/providers/registry"

func init() {
	registry.RegisterProvider(Subprovider())
}
