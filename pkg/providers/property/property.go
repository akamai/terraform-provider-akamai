//go:build all || property
// +build all property

package property

import "github.com/akamai/terraform-provider-akamai/v3/pkg/providers/registry"

func init() {
	registry.RegisterProvider(Subprovider())
}
