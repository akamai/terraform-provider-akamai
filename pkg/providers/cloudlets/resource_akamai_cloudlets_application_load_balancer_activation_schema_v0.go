package cloudlets

import "github.com/hashicorp/terraform-plugin-sdk/v2/helper/schema"

func resourceCloudletsApplicationLoadBalancerActivationV0() *schema.Resource {
	return &schema.Resource{
		Schema: map[string]*schema.Schema{
			"origin_id": {
				Type:     schema.TypeString,
				Required: true,
				ForceNew: true,
			},
			"network": {
				Type:             schema.TypeString,
				Required:         true,
				ValidateDiagFunc: validateNetwork,
			},
			"version": {
				Type:     schema.TypeInt,
				Required: true,
			},
			"status": {
				Type:     schema.TypeString,
				Computed: true,
			},
		},
		Timeouts: &schema.ResourceTimeout{
			Default: &ApplicationLoadBalancerActivationResourceTimeout,
		},
	}
}
