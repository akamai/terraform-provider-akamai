//go:build all || botman
// +build all botman

package botman

import "github.com/akamai/terraform-provider-akamai/v4/pkg/providers/registry"

func init() {
	registry.RegisterProvider(Subprovider())
}
