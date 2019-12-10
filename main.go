package main

import (
	"github.com/edglynes/terraform-providers/terraform-provider-akamai/akamai"
	"github.com/hashicorp/terraform/plugin"
)

func main() {
	plugin.Serve(&plugin.ServeOpts{
		ProviderFunc: akamai.Provider,
	})
}
