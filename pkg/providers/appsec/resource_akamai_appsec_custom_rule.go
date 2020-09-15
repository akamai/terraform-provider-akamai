package appsec

import (
	"encoding/json"
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
func resourceCustomRule() *schema.Resource {
	return &schema.Resource{
		Create: resourceCustomRuleCreate,
		Read:   resourceCustomRuleRead,
		Update: resourceCustomRuleUpdate,
		Delete: resourceCustomRuleDelete,
		Schema: map[string]*schema.Schema{
			"config_id": {
				Type:     schema.TypeInt,
				Required: true,
			},
			"rules": {
				Type:         schema.TypeString,
				Required:     true,
				ValidateFunc: validation.StringIsJSON,
			},
			"rule_id": {
				Type:     schema.TypeInt,
				Computed: true,
			},
		},
	}
}

func resourceCustomRuleCreate(d *schema.ResourceData, meta interface{}) error {
	CorrelationID := "[APPSEC][resourceCustomRuleCreate-" + tools.CreateNonce() + "]"
	edge.PrintfCorrelation("[DEBUG]", CorrelationID, " Creating CustomRule")

	customrule := appsec.NewCustomRuleResponse()

	configid := d.Get("config_id").(int)
	jsonpostpayload := d.Get("rules").(string)
	json.Unmarshal([]byte(jsonpostpayload), &customrule)

	err := customrule.SaveCustomRule(configid, CorrelationID)
	if err != nil {
		edge.PrintfCorrelation("[DEBUG]", CorrelationID, fmt.Sprintf("Error  %v\n", err))
		return err
	}

	d.Set("rule_id", customrule.ID)
	d.SetId(strconv.Itoa(customrule.ID))

	return resourceCustomRuleRead(d, meta)
}

func resourceCustomRuleUpdate(d *schema.ResourceData, meta interface{}) error {
	CorrelationID := "[APPSEC][resourceCustomRuleUpdate-" + tools.CreateNonce() + "]"
	edge.PrintfCorrelation("[DEBUG]", CorrelationID, " Updating CustomRule")

	customrule := appsec.NewCustomRuleResponse()

	configid := d.Get("config_id").(int)
	ruleid, _ := strconv.Atoi(d.Id())
	jsonpostpayload := d.Get("rules").(string)
	json.Unmarshal([]byte(jsonpostpayload), &customrule)

	err := customrule.UpdateCustomRule(configid, ruleid, CorrelationID)
	if err != nil {
		edge.PrintfCorrelation("[DEBUG]", CorrelationID, fmt.Sprintf("Error  %v\n", err))
		return nil
	}

	return resourceCustomRuleRead(d, meta)
}

func resourceCustomRuleDelete(d *schema.ResourceData, meta interface{}) error {
	CorrelationID := "[APPSEC][resourceCustomRuleDelete-" + tools.CreateNonce() + "]"
	edge.PrintfCorrelation("[DEBUG]", CorrelationID, "  Deleting CustomRule")

	customrule := appsec.NewCustomRuleResponse()

	configid := d.Get("config_id").(int)
	ruleid, _ := strconv.Atoi(d.Id())

	err := customrule.DeleteCustomRule(configid, ruleid, CorrelationID)
	if err != nil {
		edge.PrintfCorrelation("[DEBUG]", CorrelationID, fmt.Sprintf("Error  %v\n", err))
		return nil
	}

	d.SetId("")

	return nil
}

func resourceCustomRuleRead(d *schema.ResourceData, meta interface{}) error {
	CorrelationID := "[APPSEC][resourceCustomRuleRead-" + tools.CreateNonce() + "]"
	edge.PrintfCorrelation("[DEBUG]", CorrelationID, "  Read CustomRule")

	customrule := appsec.NewCustomRuleResponse()

	configid := d.Get("config_id").(int)
	ruleid, _ := strconv.Atoi(d.Id())

	err := customrule.GetCustomRule(configid, ruleid, CorrelationID)
	if err != nil {
		edge.PrintfCorrelation("[DEBUG]", CorrelationID, fmt.Sprintf("Error  %v\n", err))
		return err
	}

	d.Set("rule_id", ruleid)
	d.SetId(strconv.Itoa(ruleid))

	return nil
}
