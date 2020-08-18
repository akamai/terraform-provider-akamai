package main

import (
	"context"
	"flag"
	"os"

	"github.com/akamai/terraform-provider-akamai/v2/pkg/providers"

	"github.com/akamai/terraform-provider-akamai/v2/pkg/akamai"
	"github.com/hashicorp/go-hclog"
	"github.com/hashicorp/terraform-plugin-sdk/v2/plugin"
)

func main() {
	var debugMode bool

	flag.BoolVar(&debugMode, "debug", false, "set to true to run the provider with support for debuggers like delve")
	flag.Parse()

	// init the standard logger here so we can pass it to the provider
	logger := hclog.New(&hclog.LoggerOptions{
		Level:      hclog.Trace,
		Output:     os.Stderr,
		JSONFormat: true,
	})

	prov := akamai.Provider(
		logger,
		providers.AllProviders...,
	)

	if debugMode {
		err := plugin.Debug(context.Background(), akamai.ProviderRegistryPath,
			&plugin.ServeOpts{
				ProviderFunc: prov,
				Logger:       logger,
			})
		if err != nil {
			panic(err)
		}
	} else {
		plugin.Serve(&plugin.ServeOpts{
			ProviderFunc: prov,
			Logger:       logger,
		})
	}
}
