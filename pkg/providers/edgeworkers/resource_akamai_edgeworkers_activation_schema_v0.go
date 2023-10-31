package edgeworkers

import (
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/schema"
)

func resourceEdgeworkersActivationV0() *schema.Resource {
	return &schema.Resource{
		Schema: map[string]*schema.Schema{
			"edgeworker_id": {
				Type:     schema.TypeInt,
				Required: true,
				ForceNew: true,
			},
			"version": {
				Type:     schema.TypeString,
				Required: true,
			},
			"network": {
				Type:     schema.TypeString,
				Required: true,
			},
			"activation_id": {
				Type:     schema.TypeInt,
				Computed: true,
			},
		},
		Timeouts: &schema.ResourceTimeout{
			Delete:  &edgeworkersActivationResourceDeleteTimeout,
			Default: &edgeworkersActivationResourceDefaultTimeout,
		},
	}
}
