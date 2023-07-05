// Package main allows to run Golang code for Akamai Terraform Provider
package main

import (
	"context"
	"flag"
	"log"

	"github.com/akamai/terraform-provider-akamai/v5/pkg/akamai"
	// Load the providers
	_ "github.com/akamai/terraform-provider-akamai/v5/pkg/providers"
	"github.com/akamai/terraform-provider-akamai/v5/pkg/providers/registry"
	"github.com/hashicorp/go-hclog"
	"github.com/hashicorp/terraform-plugin-framework/providerserver"
	"github.com/hashicorp/terraform-plugin-go/tfprotov5"
	"github.com/hashicorp/terraform-plugin-go/tfprotov5/tf5server"
	"github.com/hashicorp/terraform-plugin-mux/tf5muxserver"
)

func main() {
	var debugMode bool
	flag.BoolVar(&debugMode, "debug", false, "set to true to run the provider with support for debuggers like delve")
	flag.Parse()

	// We set this to trace because logs are passed via grpc to the terraform server
	// Anything lower and we risk losing those values to the ether
	hclog.Default().SetLevel(hclog.Trace)

	providers := []func() tfprotov5.ProviderServer{
		akamai.NewPluginProvider(registry.PluginSubproviders()...)().GRPCProvider,
		providerserver.NewProtocol5(
			akamai.NewFrameworkProvider(registry.FrameworkSubproviders()...)(),
		),
	}

	muxServer, err := tf5muxserver.NewMuxServer(context.Background(), providers...)
	if err != nil {
		log.Fatal(err)
	}

	var serveOpts []tf5server.ServeOpt

	if debugMode {
		serveOpts = append(serveOpts, tf5server.WithManagedDebug())
	}

	if err = tf5server.Serve(akamai.ProviderRegistryPath, muxServer.ProviderServer, serveOpts...); err != nil {
		log.Fatal(err)
	}
}
