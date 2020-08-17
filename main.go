package main

import (
	"github.com/hashicorp/terraform-plugin-sdk/plugin"
	"github.com/terraform-providers/terraform-provider-akamai/pkg/providers"
	"github.com/terraform-providers/terraform-provider-akamai/pkg/providers/papi"
)

func main() {
	plugin.Serve(&plugin.ServeOpts{
		ProviderFunc: providers.Provider(papi.Provider{}),
	})
}
