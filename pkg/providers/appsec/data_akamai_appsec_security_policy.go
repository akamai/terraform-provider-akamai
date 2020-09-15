package appsec

import (
	"fmt"

	appsec "github.com/akamai/AkamaiOPEN-edgegrid-golang/appsec-v1"
	edge "github.com/akamai/AkamaiOPEN-edgegrid-golang/edgegrid"
	"github.com/akamai/terraform-provider-akamai/v2/pkg/tools"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/schema"
)

func dataSourceSecurityPolicy() *schema.Resource {
	return &schema.Resource{
		Read: dataSourceSecurityPolicyRead,
		Schema: map[string]*schema.Schema{
			"name": {
				Type:     schema.TypeString,
				Optional: true,
			},
			"config_id": {
				Type:     schema.TypeInt,
				Required: true,
			},
			"version": {
				Type:     schema.TypeInt,
				Required: true,
			},
			"policy_id": {
				Type:        schema.TypeString,
				Computed:    true,
				Description: "Policy ID",
			},
			"policy_list": {
				Type:        schema.TypeList,
				Computed:    true,
				Elem:        &schema.Schema{Type: schema.TypeString},
				Description: "Policy ID List",
			},
		},
	}
}

func dataSourceSecurityPolicyRead(d *schema.ResourceData, meta interface{}) error {
	CorrelationID := "[APPSEC][dataSourceSecurityPolicyRead-" + tools.CreateNonce() + "]"

	edge.PrintfCorrelation("[DEBUG]", CorrelationID, "  Read SecurityPolicy")

	securitypolicy := appsec.NewSecurityPolicyResponse()
	configName := d.Get("name").(string)
	securitypolicy.ConfigID = d.Get("config_id").(int)
	securitypolicy.Version = d.Get("version").(int)

	err := securitypolicy.GetSecurityPolicy(CorrelationID)
	if err != nil {
		edge.PrintfCorrelation("[DEBUG]", CorrelationID, fmt.Sprintf("Error  %v\n", err))
		return nil
	}

	edge.PrintfCorrelation("[DEBUG]", CorrelationID, fmt.Sprintf("SecurityPolicy   %v\n", securitypolicy))

	secpolicylist := make([]string, 0, len(securitypolicy.Policies))

	for _, configval := range securitypolicy.Policies {
		edge.PrintfCorrelation("[DEBUG]", CorrelationID, fmt.Sprintf("CONFIG value  %v\n", configval.PolicyID))
		secpolicylist = append(secpolicylist, configval.PolicyID)
		if configval.PolicyName == configName {
			d.Set("policy_id", configval.PolicyID)
		}
	}

	d.Set("policy_list", secpolicylist)
	d.SetId(fmt.Sprintf("%d:%d", securitypolicy.ConfigID, securitypolicy.Version))

	return nil
}
