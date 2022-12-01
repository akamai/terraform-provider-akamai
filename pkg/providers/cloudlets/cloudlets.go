//go:build all || cloudlets
// +build all cloudlets

package cloudlets

import "github.com/akamai/terraform-provider-akamai/v3/pkg/providers/registry"

func init() {
	registry.RegisterProvider(Subprovider())
}
