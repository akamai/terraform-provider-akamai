//go:build all || cps
// +build all cps

package cps

import "github.com/akamai/terraform-provider-akamai/v4/pkg/providers/registry"

func init() {
	registry.RegisterProvider(Subprovider())
}
