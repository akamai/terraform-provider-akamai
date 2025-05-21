package property

import (
	"github.com/akamai/AkamaiOPEN-edgegrid-golang/v11/pkg/papi"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/schema"
)

func resourcePropertyActivationV0() *schema.Resource {
	return &schema.Resource{
		Schema: map[string]*schema.Schema{
			"property_id": {
				Type:     schema.TypeString,
				Required: true,
			},
			"activation_id": {
				Type:     schema.TypeString,
				Optional: true,
				Computed: true,
			},
			"errors": {
				Type:     schema.TypeString,
				Computed: true,
			},
			"warnings": {
				Type:     schema.TypeString,
				Computed: true,
			},
			"rule_errors": {
				Type:     schema.TypeList,
				Optional: true,
				Computed: true,
				Elem:     papiError(),
			},
			"auto_acknowledge_rule_warnings": {
				Type:     schema.TypeBool,
				Optional: true,
				Default:  false,
			},
			"version": {
				Type:     schema.TypeInt,
				Required: true,
			},
			"network": {
				Type:     schema.TypeString,
				Optional: true,
				Default:  papi.ActivationNetworkStaging,
			},
			"contact": {
				Type:     schema.TypeSet,
				Required: true,
				Elem:     &schema.Schema{Type: schema.TypeString},
			},
			"status": {
				Type:     schema.TypeString,
				Computed: true,
			},
			"note": {
				Type:     schema.TypeString,
				Optional: true,
			},
			"compliance_record": {
				Type:     schema.TypeList,
				Optional: true,
				MaxItems: 1,
				Elem:     complianceRecordSchema,
			},
		},
		Timeouts: &schema.ResourceTimeout{
			Default: &PropertyResourceTimeout,
		},
	}
}
