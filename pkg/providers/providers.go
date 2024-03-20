// Package providers allows to import list of sub-providers
package providers

import (
	// This is where providers are imported, so they can register themselves
	_ "github.com/akamai/terraform-provider-akamai/v6/pkg/providers/appsec"
	_ "github.com/akamai/terraform-provider-akamai/v6/pkg/providers/botman"
	_ "github.com/akamai/terraform-provider-akamai/v6/pkg/providers/clientlists"
	_ "github.com/akamai/terraform-provider-akamai/v6/pkg/providers/cloudlets"
	_ "github.com/akamai/terraform-provider-akamai/v6/pkg/providers/cloudwrapper"
	_ "github.com/akamai/terraform-provider-akamai/v6/pkg/providers/cps"
	_ "github.com/akamai/terraform-provider-akamai/v6/pkg/providers/datastream"
	_ "github.com/akamai/terraform-provider-akamai/v6/pkg/providers/dns"
	_ "github.com/akamai/terraform-provider-akamai/v6/pkg/providers/edgeworkers"
	_ "github.com/akamai/terraform-provider-akamai/v6/pkg/providers/gtm"
	_ "github.com/akamai/terraform-provider-akamai/v6/pkg/providers/iam"
	_ "github.com/akamai/terraform-provider-akamai/v6/pkg/providers/imaging"
	_ "github.com/akamai/terraform-provider-akamai/v6/pkg/providers/networklists"
	_ "github.com/akamai/terraform-provider-akamai/v6/pkg/providers/property"
)
