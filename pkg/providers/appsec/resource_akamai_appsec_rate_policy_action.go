package appsec

import (
	"fmt"
	"strconv"

	appsec "github.com/akamai/AkamaiOPEN-edgegrid-golang/appsec-v1"
	edge "github.com/akamai/AkamaiOPEN-edgegrid-golang/edgegrid"
	"github.com/akamai/terraform-provider-akamai/v2/pkg/tools"

	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/schema"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/validation"
)

// appsec v1
//
// https://developer.akamai.com/api/cloud_security/application_security/v1.html
func resourceRatePolicyAction() *schema.Resource {
	return &schema.Resource{
		Create: resourceRatePolicyActionUpdate,
		Read:   resourceRatePolicyActionRead,
		Update: resourceRatePolicyActionUpdate,
		Delete: resourceRatePolicyActionDelete,
		Importer: &schema.ResourceImporter{
			State: schema.ImportStatePassthrough,
		},
		Schema: map[string]*schema.Schema{
			"config_id": {
				Type:     schema.TypeInt,
				Required: true,
			},
			"version": {
				Type:     schema.TypeInt,
				Required: true,
			},
			"policy_id": {
				Type:     schema.TypeString,
				Required: true,
			},
			"rate_policy_id": {
				Type:     schema.TypeInt,
				Required: true,
			},
			"ipv4_action": {
				Type:     schema.TypeString,
				Required: true,
				ValidateFunc: validation.StringInSlice([]string{
					Alert,
					Deny,
					None,
				}, false),
			},
			"ipv6_action": {
				Type:     schema.TypeString,
				Required: true,
				ValidateFunc: validation.StringInSlice([]string{
					Alert,
					Deny,
					None,
				}, false)},
		},
	}
}

func resourceRatePolicyActionRead(d *schema.ResourceData, meta interface{}) error {
	CorrelationID := "[APPSEC][resourceRatePolicyActionRead-" + tools.CreateNonce() + "]"
	edge.PrintfCorrelation("[DEBUG]", CorrelationID, "  Read RatePolicyAction")

	ratepolicyaction := appsec.NewRatePolicyActionResponse()

	configid := d.Get("config_id").(int)
	version := d.Get("version").(int)
	policyid := d.Get("policy_id").(string)

	err := ratepolicyaction.GetRatePolicyAction(configid, version, policyid, CorrelationID)
	if err != nil {
		edge.PrintfCorrelation("[DEBUG]", CorrelationID, fmt.Sprintf("Error  %v\n", err))
		return nil
	}

	for _, configval := range ratepolicyaction.RatePolicyActions {
		edge.PrintfCorrelation("[DEBUG]", CorrelationID, fmt.Sprintf("ratepolicyaction  configval %v\n", configval.ID))
		d.SetId(strconv.Itoa(configval.ID))
	}
	return nil
}

func resourceRatePolicyActionDelete(d *schema.ResourceData, meta interface{}) error {
	CorrelationID := "[APPSEC][resourceRatePolicyActionDelete-" + tools.CreateNonce() + "]"
	edge.PrintfCorrelation("[DEBUG]", CorrelationID, "  Deleting RatePolicyAction")
	ratepolicyaction := appsec.NewRatePolicyActionResponse()

	ratepolicyactionpost := appsec.NewRatePolicyActionPost()

	configid := d.Get("config_id").(int)
	version := d.Get("version").(int)
	policyid := d.Get("policy_id").(string)
	rate_policy_id := d.Get("rate_policy_id").(int)
	ratepolicyactionpost.Ipv4Action = "none"
	ratepolicyactionpost.Ipv6Action = "none"

	err := ratepolicyaction.UpdateRatePolicyAction(configid, version, policyid, rate_policy_id, ratepolicyactionpost, CorrelationID)
	if err != nil {
		edge.PrintfCorrelation("[DEBUG]", CorrelationID, fmt.Sprintf("Error  %v\n", err))
		return err
	}
	return nil
	//return schema.Noop(d, meta)
}

func resourceRatePolicyActionUpdate(d *schema.ResourceData, meta interface{}) error {
	CorrelationID := "[APPSEC][resourceRatePolicyActionUpdate-" + tools.CreateNonce() + "]"
	edge.PrintfCorrelation("[DEBUG]", CorrelationID, "  Updating RatePolicyAction")

	ratepolicyaction := appsec.NewRatePolicyActionResponse()

	ratepolicyactionpost := appsec.NewRatePolicyActionPost()

	configid := d.Get("config_id").(int)
	version := d.Get("version").(int)
	policyid := d.Get("policy_id").(string)
	rate_policy_id := d.Get("rate_policy_id").(int)
	ratepolicyactionpost.Ipv4Action = d.Get("ipv4_action").(string)
	ratepolicyactionpost.Ipv6Action = d.Get("ipv6_action").(string)

	err := ratepolicyaction.UpdateRatePolicyAction(configid, version, policyid, rate_policy_id, ratepolicyactionpost, CorrelationID)
	if err != nil {
		edge.PrintfCorrelation("[DEBUG]", CorrelationID, fmt.Sprintf("Error  %v\n", err))
		return err
	}

	return resourceRatePolicyActionRead(d, meta)

}

const (
	Alert = "alert"
	Deny  = "deny"
	None  = "none"
)
