// Package main allows to run Golang code for Akamai Terraform Provider
package main

import (
	"context"
	"flag"
	"log"
	"log/slog"
	"os"

	akalog "github.com/akamai/AkamaiOPEN-edgegrid-golang/v10/pkg/log"
	"github.com/akamai/terraform-provider-akamai/v8/pkg/akamai"
	_ "github.com/akamai/terraform-provider-akamai/v8/pkg/providers" // Load the providers
	"github.com/akamai/terraform-provider-akamai/v8/pkg/providers/registry"
	"github.com/hashicorp/go-hclog"
	"github.com/hashicorp/terraform-plugin-framework/providerserver"
	"github.com/hashicorp/terraform-plugin-go/tfprotov6"
	"github.com/hashicorp/terraform-plugin-go/tfprotov6/tf6server"
	"github.com/hashicorp/terraform-plugin-mux/tf6muxserver"
)

func main() {
	var debugMode bool
	flag.BoolVar(&debugMode, "debug", false, "set to true to run the provider with support for debuggers like delve")
	flag.Parse()
	if debugMode {
		debugHandler := akalog.NewSlogHandler(os.Stderr, &slog.HandlerOptions{
			Level: slog.LevelDebug,
		})
		akalog.SetLogger(akalog.NewSlogAdapter(debugHandler))
	}

	// We set this to trace because logs are passed via grpc to the terraform server
	// Anything lower and we risk losing those values to the ether
	hclog.Default().SetLevel(hclog.Trace)
	sdkProviderV6, err := akamai.NewProtoV6SDKProvider(registry.Subproviders())
	if err != nil {
		log.Fatal(err)
	}

	providers := []func() tfprotov6.ProviderServer{
		sdkProviderV6,
		providerserver.NewProtocol6(
			akamai.NewFrameworkProvider(registry.Subproviders()...)(),
		),
	}

	muxServer, err := tf6muxserver.NewMuxServer(context.Background(), providers...)
	if err != nil {
		log.Fatal(err)
	}

	var serveOpts []tf6server.ServeOpt

	if debugMode {
		serveOpts = append(serveOpts, tf6server.WithManagedDebug())
	}

	if err = tf6server.Serve(akamai.ProviderRegistryPath, muxServer.ProviderServer, serveOpts...); err != nil {
		log.Fatal(err)
	}
}
