package property

import (
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/schema"
)

func resourcePropertyIncludeActivationV0() *schema.Resource {
	return &schema.Resource{
		Schema: map[string]*schema.Schema{
			"include_id": {
				Type:     schema.TypeString,
				Required: true,
				ForceNew: true,
			},
			"contract_id": {
				Type:     schema.TypeString,
				Required: true,
				ForceNew: true,
			},
			"group_id": {
				Type:     schema.TypeString,
				Required: true,
				ForceNew: true,
			},
			"version": {
				Type:     schema.TypeInt,
				Required: true,
			},
			"network": {
				Type:     schema.TypeString,
				Required: true,
				ForceNew: true,
			},
			"notify_emails": {
				Type:     schema.TypeSet,
				Required: true,
				MinItems: 1,
				Elem: &schema.Schema{
					Type: schema.TypeString,
				},
			},
			"note": {
				Type:     schema.TypeString,
				Optional: true,
				Default:  "",
			},
			"auto_acknowledge_rule_warnings": {
				Type:     schema.TypeBool,
				Optional: true,
				Default:  false,
			},
			"validations": {
				Type:     schema.TypeString,
				Computed: true,
			},
			"compliance_record": {
				Type:     schema.TypeList,
				Optional: true,
				MaxItems: 1,
				Elem:     complianceRecordSchema,
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
					},
				},
			},
		},
		Timeouts: &schema.ResourceTimeout{
			Default: readTimeoutFromEnvOrDefault("AKAMAI_ACTIVATION_TIMEOUT", includeActivationTimeout),
		},
	}
}
