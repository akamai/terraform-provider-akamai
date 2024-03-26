package cloudlets

import (
	"github.com/akamai/terraform-provider-akamai/v6/pkg/common/tf"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/schema"
)

func resourceCloudletsPolicyActivationV0() *schema.Resource {
	return &schema.Resource{
		Schema: map[string]*schema.Schema{
			"status": {
				Type:     schema.TypeString,
				Computed: true,
			},
			"policy_id": {
				Type:     schema.TypeInt,
				Required: true,
				ForceNew: true,
			},
			"network": {
				Type:             schema.TypeString,
				Required:         true,
				ValidateDiagFunc: tf.ValidateNetwork,
			},
			"version": {
				Type:     schema.TypeInt,
				Required: true,
			},
			"associated_properties": {
				Type:     schema.TypeSet,
				Required: true,
				Elem:     &schema.Schema{Type: schema.TypeString},
				MinItems: 1,
			},
		},
		Timeouts: &schema.ResourceTimeout{
			Default: &PolicyActivationResourceTimeout,
		},
	}
}
