package providers

import (
	"github.com/akamai/terraform-provider-akamai/v2/pkg/akamai"
	"github.com/akamai/terraform-provider-akamai/v2/pkg/providers/dns"
	"github.com/akamai/terraform-provider-akamai/v2/pkg/providers/gtm"
	"github.com/akamai/terraform-provider-akamai/v2/pkg/providers/property"
)

var (
	// AllProviders is the "registry" of all providers
	AllProviders = []akamai.Subprovider{
		property.Subprovider(),
		gtm.Subprovider(),
		dns.Subprovider(),
	}
)
