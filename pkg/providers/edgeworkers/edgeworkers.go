//go:build all || edgeworkers
// +build all edgeworkers

package edgeworkers

import "github.com/akamai/terraform-provider-akamai/v4/pkg/providers/registry"

func init() {
	registry.RegisterProvider(Subprovider())
}
