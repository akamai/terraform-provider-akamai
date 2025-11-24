package testutils

import (
	"context"
	"fmt"

	"github.com/akamai/terraform-provider-akamai/v9/internal/edgegrid"
	"github.com/akamai/terraform-provider-akamai/v9/pkg/akamai"
	"github.com/akamai/terraform-provider-akamai/v9/pkg/meta"
	"github.com/akamai/terraform-provider-akamai/v9/pkg/subprovider"
	"github.com/hashicorp/terraform-plugin-framework/providerserver"
	"github.com/hashicorp/terraform-plugin-go/tfprotov6"
	"github.com/hashicorp/terraform-plugin-mux/tf5to6server"
	"github.com/hashicorp/terraform-plugin-mux/tf6muxserver"
	"github.com/hashicorp/terraform-plugin-sdk/v2/diag"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/schema"
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

// NewTestProtoV6ProviderFactory uses provided subproviders and client to create provider factory for test purposes
func NewTestProtoV6ProviderFactory(client edgegrid.Client, subproviders ...subprovider.Subprovider) map[string]func() (tfprotov6.ProviderServer, error) {
	providerFunc := akamai.NewTestFrameworkProvider(client, subproviders...)
	return map[string]func() (tfprotov6.ProviderServer, error){
		"akamai": func() (tfprotov6.ProviderServer, error) {
			return providerserver.NewProtocol6(providerFunc())(), nil
		},
	}
}

// NewTestProtoV6SDKProviderFactory uses provided subproviders and client to create SDK provider factory for test purposes
func NewTestProtoV6SDKProviderFactory(client edgegrid.Client, subproviders ...subprovider.Subprovider) map[string]func() (tfprotov6.ProviderServer, error) {
	return map[string]func() (tfprotov6.ProviderServer, error){
		"akamai": func() (tfprotov6.ProviderServer, error) {
			provider := akamai.NewSDKProvider(subproviders...)()
			provider.ConfigureContextFunc = wrapConfigureContextFunc(provider.ConfigureContextFunc, client)
			v6Provider, err := tf5to6server.UpgradeServer(
				context.Background(),
				provider.GRPCProvider,
			)
			if err != nil {
				return nil, err
			}
			return v6Provider, nil
		},
	}
}

func wrapConfigureContextFunc(fn schema.ConfigureContextFunc, client edgegrid.Client) schema.ConfigureContextFunc {
	return func(ctx context.Context, d *schema.ResourceData) (any, diag.Diagnostics) {
		m, diags := fn(ctx, d)
		if diags.HasError() {
			return nil, diags
		}
		opMeta, ok := m.(*meta.OperationMeta)
		if !ok {
			return nil, diag.FromErr(fmt.Errorf("expected meta.OperationMeta, got %T", m))
		}
		opMeta.SetClient(client)
		return opMeta, diags
	}
}
