package testutils

import (
	"context"

	"github.com/akamai/terraform-provider-akamai/v5/pkg/akamai"
	"github.com/akamai/terraform-provider-akamai/v5/pkg/subprovider"
	"github.com/hashicorp/terraform-plugin-framework/providerserver"
	"github.com/hashicorp/terraform-plugin-go/tfprotov6"
	"github.com/hashicorp/terraform-plugin-mux/tf6muxserver"
)

// NewProtoV6ProviderFactory uses provided subprovider to create provider factory for test purposes
func NewProtoV6ProviderFactory(subproviders ...subprovider.Subprovider) map[string]func() (tfprotov6.ProviderServer, error) {
	return map[string]func() (tfprotov6.ProviderServer, error){
		"akamai": func() (tfprotov6.ProviderServer, error) {
			ctx := context.Background()

			sdkProviderV6, err := akamai.NewProtoV6SDKProvider(subproviders)
			if err != nil {
				return nil, err
			}

			providers := []func() tfprotov6.ProviderServer{
				sdkProviderV6,
				providerserver.NewProtocol6(
					akamai.NewFrameworkProvider(subproviders...)(),
				),
			}

			muxServer, err := tf6muxserver.NewMuxServer(ctx, providers...)
			if err != nil {
				return nil, err
			}

			return muxServer.ProviderServer(), nil
		},
	}
}
