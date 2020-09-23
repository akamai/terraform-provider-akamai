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
func resourceCustomRuleAction() *schema.Resource {
	return &schema.Resource{
		Create: resourceCustomRuleActionUpdate,
		Read:   resourceCustomRuleActionRead,
		Update: resourceCustomRuleActionUpdate,
		Delete: resourceCustomRuleActionDelete,
		Importer: &schema.ResourceImporter{
			State: schema.ImportStatePassthrough,
		},
		Schema: map[string]*schema.Schema{
			"config_id": {
				Type:     schema.TypeInt,
				Required: true,
			},
			"version": &schema.Schema{
				Type:     schema.TypeInt,
				Required: true,
			},
			"custom_rule_action": {
				Type:     schema.TypeString,
				Required: true,
				ValidateFunc: validation.StringInSlice([]string{
					Alert,
					Deny,
					None,
				}, false),
			},
			"policy_id": {
				Type:     schema.TypeString,
				Required: true,
			},
			"rule_id": {
				Type:     schema.TypeInt,
				Required: true,
			},
		},
	}
}

func resourceCustomRuleActionRead(d *schema.ResourceData, meta interface{}) error {
	CorrelationID := "[APPSEC][resourceCustomRuleActionRead-" + tools.CreateNonce() + "]"
	edge.PrintfCorrelation("[DEBUG]", CorrelationID, "  Read CustomRuleAction")

	customruleaction := appsec.NewCustomRuleActionResponse()

	configid := d.Get("config_id").(int)
	version := d.Get("version").(int)
	policyid := d.Get("policy_id").(string)
	ruleid := d.Get("rule_id").(int)

	err := customruleaction.GetCustomRuleAction(configid, version, policyid, CorrelationID)
	if err != nil {
		edge.PrintfCorrelation("[DEBUG]", CorrelationID, fmt.Sprintf("Error  %v\n", err))
		return nil
	}

	d.Set("rule_id", ruleid)
	d.SetId(strconv.Itoa(ruleid))
	return nil
}

func resourceCustomRuleActionDelete(d *schema.ResourceData, meta interface{}) error {
	CorrelationID := "[APPSEC][resourceCustomRuleActionDelete-" + tools.CreateNonce() + "]"
	edge.PrintfCorrelation("[DEBUG]", CorrelationID, "  Deleting CustomRuleAction")

	customruleaction := appsec.NewCustomRuleActionResponse()
	customruleactionpost := appsec.NewCustomRuleActionPost()

	configid := d.Get("config_id").(int)
	version := d.Get("version").(int)
	policyid := d.Get("policy_id").(string)
	ruleid := d.Get("rule_id").(int)
	customruleactionpost.Action = d.Get("custom_rule_action").(string)

	err := customruleaction.UpdateCustomRuleAction(configid, version, policyid, ruleid, customruleactionpost, CorrelationID)
	if err != nil {
		edge.PrintfCorrelation("[DEBUG]", CorrelationID, fmt.Sprintf("Error  %v\n", err))
		return err
	}

	return nil
}

func resourceCustomRuleActionUpdate(d *schema.ResourceData, meta interface{}) error {
	CorrelationID := "[APPSEC][resourceCustomRuleActionUpdate-" + tools.CreateNonce() + "]"
	edge.PrintfCorrelation("[DEBUG]", CorrelationID, "  Updating CustomRuleAction")

	customruleaction := appsec.NewCustomRuleActionResponse()

	customruleactionpost := appsec.NewCustomRuleActionPost()

	configid := d.Get("config_id").(int)
	version := d.Get("version").(int)
	policyid := d.Get("policy_id").(string)
	ruleid := d.Get("rule_id").(int)
	customruleactionpost.Action = d.Get("custom_rule_action").(string)

	err := customruleaction.UpdateCustomRuleAction(configid, version, policyid, ruleid, customruleactionpost, CorrelationID)
	if err != nil {
		edge.PrintfCorrelation("[DEBUG]", CorrelationID, fmt.Sprintf("Error  %v\n", err))
		return err
	}

	return resourceCustomRuleActionRead(d, meta)

}
