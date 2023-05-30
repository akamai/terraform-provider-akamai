// Package config contains set of tools which allow to configure application
package config

import (
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/schema"
)

// Options initializes and returns terraform.Resource with credentials
func Options(_ string) *schema.Resource {
	return &schema.Resource{
		Schema: map[string]*schema.Schema{
			"host": {
				Type:     schema.TypeString,
				Required: true,
			},
			"access_token": {
				Type:     schema.TypeString,
				Required: true,
			},
			"client_token": {
				Type:     schema.TypeString,
				Required: true,
			},
			"client_secret": {
				Type:     schema.TypeString,
				Required: true,
			},
			"max_body": {
				Type:     schema.TypeInt,
				Optional: true,
			},
			"account_key": {
				Type:     schema.TypeString,
				Optional: true,
			},
		},
	}
}
