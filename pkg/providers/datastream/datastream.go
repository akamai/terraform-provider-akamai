//go:build all || datastream
// +build all datastream

package datastream

import "github.com/akamai/terraform-provider-akamai/v4/pkg/providers/registry"

func init() {
	registry.RegisterProvider(Subprovider())
}
