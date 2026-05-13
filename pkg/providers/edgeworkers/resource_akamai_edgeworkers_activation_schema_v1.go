package edgeworkers

import (
	"context"

	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/schema"
)

func resourceEdgeworkersActivationV1() *schema.Resource {
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
			"note": {
				Type:     schema.TypeString,
				Optional: true,
			},
			"timeouts": {
				Type:     schema.TypeList,
				Optional: true,
				MaxItems: 1,
				Elem: &schema.Resource{
					Schema: map[string]*schema.Schema{
						"default": {
							Type:     schema.TypeString,
							Optional: true,
						},
						"delete": {
							Type:     schema.TypeString,
							Optional: true,
						},
					},
				},
			},
		},
		Timeouts: &schema.ResourceTimeout{
			Delete:  &edgeworkersActivationResourceDeleteTimeout,
			Default: &edgeworkersActivationResourceDefaultTimeout,
		},
	}
}

func upgradeEdgeworkersActivationV1(_ context.Context, rawState map[string]interface{}, _ interface{}) (map[string]interface{}, error) {
	if _, ok := rawState["auto_pin"]; !ok {
		rawState["auto_pin"] = true
	}
	return rawState, nil
}
