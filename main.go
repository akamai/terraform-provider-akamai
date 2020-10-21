package main

import (
	"context"
	"flag"

	// Load the providers
	_ "github.com/akamai/terraform-provider-akamai/v2/pkg/providers"

	"github.com/akamai/terraform-provider-akamai/v2/pkg/providers/registry"

	"github.com/akamai/terraform-provider-akamai/v2/pkg/akamai"
	"github.com/hashicorp/go-hclog"
	"github.com/hashicorp/terraform-plugin-sdk/v2/plugin"
)

func main() {
	var debugMode bool

	flag.BoolVar(&debugMode, "debug", false, "set to true to run the provider with support for debuggers like delve")
	flag.Parse()

	// We set this to trace because logs are passed via grpc to the terraform server
	// Anything lower and we risk losing those values to the ether
	hclog.Default().SetLevel(hclog.Trace)

	prov := akamai.Provider(
		registry.AllProviders()...,
	)

	if debugMode {
		err := plugin.Debug(context.Background(), akamai.ProviderRegistryPath,
			&plugin.ServeOpts{
				ProviderFunc: prov,
			})
		if err != nil {
			panic(err)
		}
	} else {
		plugin.Serve(&plugin.ServeOpts{
			ProviderFunc: prov,
		})
	}
}
