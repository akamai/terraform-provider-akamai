// Package providers allows to import list of sub-providers
package providers

import (
	// This is where providers are import so they can register themselves
	_ "github.com/akamai/terraform-provider-akamai/v5/pkg/providers/appsec"
	_ "github.com/akamai/terraform-provider-akamai/v5/pkg/providers/botman"
	_ "github.com/akamai/terraform-provider-akamai/v5/pkg/providers/clientlists"
	_ "github.com/akamai/terraform-provider-akamai/v5/pkg/providers/cloudlets"
	_ "github.com/akamai/terraform-provider-akamai/v5/pkg/providers/cps"
	_ "github.com/akamai/terraform-provider-akamai/v5/pkg/providers/datastream"
	_ "github.com/akamai/terraform-provider-akamai/v5/pkg/providers/dns"
	_ "github.com/akamai/terraform-provider-akamai/v5/pkg/providers/edgeworkers"
	_ "github.com/akamai/terraform-provider-akamai/v5/pkg/providers/gtm"
	_ "github.com/akamai/terraform-provider-akamai/v5/pkg/providers/iam"
	_ "github.com/akamai/terraform-provider-akamai/v5/pkg/providers/imaging"
	_ "github.com/akamai/terraform-provider-akamai/v5/pkg/providers/networklists"
	_ "github.com/akamai/terraform-provider-akamai/v5/pkg/providers/property"
)
