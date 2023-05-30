// Package config contains set of tools which allow to configure application
package config

import (
	frameworkSchema "github.com/hashicorp/terraform-plugin-framework/provider/schema"
	pluginSchema "github.com/hashicorp/terraform-plugin-sdk/v2/helper/schema"
)

// PluginOptions returns edgegrid config schema for terraform-plugin-sdk
func PluginOptions() *pluginSchema.Resource {
	return &pluginSchema.Resource{
		Schema: map[string]*pluginSchema.Schema{
			"host": {
				Type:     pluginSchema.TypeString,
				Required: true,
			},
			"access_token": {
				Type:     pluginSchema.TypeString,
				Required: true,
			},
			"client_token": {
				Type:     pluginSchema.TypeString,
				Required: true,
			},
			"client_secret": {
				Type:     pluginSchema.TypeString,
				Required: true,
			},
			"max_body": {
				Type:     pluginSchema.TypeInt,
				Optional: true,
			},
			"account_key": {
				Type:     pluginSchema.TypeString,
				Optional: true,
			},
		},
	}
}

// FrameworkOptions returns edgegrid config schema for terraform-plugin-framework
func FrameworkOptions() frameworkSchema.SetNestedBlock {
	return frameworkSchema.SetNestedBlock{
		NestedObject: frameworkSchema.NestedBlockObject{
			Attributes: map[string]frameworkSchema.Attribute{
				"host": frameworkSchema.StringAttribute{
					Required: true,
				},
				"access_token": frameworkSchema.StringAttribute{
					Required: true,
				},
				"client_token": frameworkSchema.StringAttribute{
					Required: true,
				},
				"client_secret": frameworkSchema.StringAttribute{
					Required: true,
				},
				"max_body": frameworkSchema.Int64Attribute{
					Optional: true,
				},
				"account_key": frameworkSchema.StringAttribute{
					Optional: true,
				},
			},
		},
	}
}
