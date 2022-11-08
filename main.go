package main

import (
	"context"
	"flag"

	// Load the providers
	_ "github.com/akamai/terraform-provider-akamai/v3/pkg/providers"
	"github.com/akamai/terraform-provider-akamai/v3/pkg/providers/registry"
	goplugin "github.com/hashicorp/go-plugin"
	"github.com/hashicorp/terraform-plugin-go/tfprotov5"
	"github.com/hashicorp/terraform-plugin-go/tfprotov5/tf5server"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/schema"
	"google.golang.org/grpc"

	"github.com/akamai/terraform-provider-akamai/v3/pkg/akamai"
	"github.com/hashicorp/go-hclog"
	"github.com/hashicorp/terraform-plugin-sdk/v2/plugin"
)

// gRPC message limit of 64MB
const gRPCLimit = 64 << 20

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
		// modified implementation of plugin.Serve() method from terraform SDK
		// this is done in order to increase the max GRPC limit from 4MB to 64MB
		serveConfig := goplugin.ServeConfig{
			HandshakeConfig: plugin.Handshake,
			GRPCServer: func(opts []grpc.ServerOption) *grpc.Server {
				return grpc.NewServer(append(opts,
					grpc.MaxSendMsgSize(gRPCLimit),
					grpc.MaxRecvMsgSize(gRPCLimit))...)
			},
			VersionedPlugins: map[int]goplugin.PluginSet{
				5: {
					akamai.ProviderRegistryPath: &tf5server.GRPCProviderPlugin{
						GRPCProvider: func() tfprotov5.ProviderServer {
							return schema.NewGRPCProviderServer(prov())
						},
					},
				},
			},
		}
		goplugin.Serve(&serveConfig)
	}
}
