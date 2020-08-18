package main

import (
	"context"
	"flag"
	"github.com/akamai/terraform-provider-akamai/v2/pkg/providers/dns"
	"github.com/akamai/terraform-provider-akamai/v2/pkg/providers/gtm"
	"github.com/akamai/terraform-provider-akamai/v2/pkg/providers/property"
	"log"
	"os"

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

	prov := akamai.Provider(logger, property.Subprovider(), dns.Subprovider(), gtm.Subprovider())

	if debugMode {
		err := plugin.Debug(context.Background(), "registry.terraform.io/akamai/akamai",
			&plugin.ServeOpts{
				ProviderFunc: prov,
				Logger:       logger,
			})
		if err != nil {
			log.Println(err.Error())
		}
	} else {
		plugin.Serve(&plugin.ServeOpts{
			ProviderFunc: prov,
			Logger:       logger,
		})
	}
}
