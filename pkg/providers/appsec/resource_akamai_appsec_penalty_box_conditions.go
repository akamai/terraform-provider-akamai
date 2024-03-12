package appsec

import (
	"context"
	"encoding/json"
	"fmt"
	"strconv"

	"github.com/akamai/AkamaiOPEN-edgegrid-golang/v7/pkg/appsec"
	"github.com/akamai/terraform-provider-akamai/v5/pkg/common/tf"
	"github.com/akamai/terraform-provider-akamai/v5/pkg/meta"
	"github.com/hashicorp/terraform-plugin-sdk/v2/diag"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/customdiff"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/schema"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/validation"
)

// appsec v1
//
// https://techdocs.akamai.com/application-security/reference/api
func resourcePenaltyBoxConditions() *schema.Resource {
	return &schema.Resource{
		CreateContext: resourcePenaltyBoxConditionsCreate,
		ReadContext:   resourcePenaltyBoxConditionsRead,
		UpdateContext: resourcePenaltyBoxConditionsUpdate,
		DeleteContext: resourcePenaltyBoxConditionsDelete,
		CustomizeDiff: customdiff.All(
			VerifyIDUnchanged,
		),
		Importer: &schema.ResourceImporter{
			StateContext: schema.ImportStatePassthroughContext,
		},
		Schema: map[string]*schema.Schema{
			"config_id": {
				Type:        schema.TypeInt,
				Required:    true,
				Description: "Unique identifier of the security configuration",
			},
			"security_policy_id": {
				Type:        schema.TypeString,
				Required:    true,
				Description: "Unique identifier of the security policy",
			},
			"penalty_box_conditions": {
				Type:             schema.TypeString,
				Required:         true,
				ValidateDiagFunc: validation.ToDiagFunc(validation.StringIsJSON),
				DiffSuppressFunc: suppressEquivalentPenaltyBoxConditionsDiffs,
				Description:      "Describes the conditions and the operator to be applied for penalty box",
			},
		},
	}
}

func resourcePenaltyBoxConditionsCreate(ctx context.Context, d *schema.ResourceData, m interface{}) diag.Diagnostics {
	meta := meta.Must(m)
	client := inst.Client(meta)
	logger := meta.Log("APPSEC", "resourcePenaltyBoxConditionsCreate")
	logger.Debugf("in resourcePenaltyBoxConditionsCreate")

	configID, err := tf.GetIntValue("config_id", d)
	if err != nil {
		return diag.FromErr(err)
	}
	policyID, err := tf.GetStringValue("security_policy_id", d)
	if err != nil {
		return diag.FromErr(err)
	}
	version, err := getModifiableConfigVersion(ctx, configID, "penaltyBoxConditions", m)
	if err != nil {
		return diag.FromErr(err)
	}
	jsonPostPayload := d.Get("penalty_box_conditions")
	conditionsJSON := appsec.PenaltyBoxConditionsPayload{}
	err = json.Unmarshal([]byte(jsonPostPayload.(string)), &conditionsJSON)
	if err != nil {
		return diag.FromErr(err)
	}

	updatePenaltyBoxConditions := appsec.UpdatePenaltyBoxConditionsRequest{
		ConfigID:          configID,
		Version:           version,
		PolicyID:          policyID,
		ConditionsPayload: conditionsJSON,
	}

	_, err = client.UpdatePenaltyBoxConditions(ctx, updatePenaltyBoxConditions)
	if err != nil {
		logger.Errorf("calling 'createPenaltyBoxConditions': %s", err.Error())
		return diag.FromErr(err)
	}

	d.SetId(fmt.Sprintf("%d:%s", updatePenaltyBoxConditions.ConfigID, updatePenaltyBoxConditions.PolicyID))

	return resourcePenaltyBoxConditionsRead(ctx, d, m)
}

func resourcePenaltyBoxConditionsRead(ctx context.Context, d *schema.ResourceData, m interface{}) diag.Diagnostics {
	meta := meta.Must(m)
	client := inst.Client(meta)
	logger := meta.Log("APPSEC", "resourcePenaltyBoxConditionsRead")

	iDParts, err := splitID(d.Id(), 2, "configID:securityPolicyID")
	if err != nil {
		return diag.FromErr(err)
	}
	configID, err := strconv.Atoi(iDParts[0])
	if err != nil {
		return diag.FromErr(err)
	}
	version, err := getLatestConfigVersion(ctx, configID, m)
	if err != nil {
		return diag.FromErr(err)
	}
	policyID := iDParts[1]

	getPenaltyBoxConditions := appsec.GetPenaltyBoxConditionsRequest{
		ConfigID: configID,
		Version:  version,
		PolicyID: policyID,
	}

	penaltyBoxConditions, err := client.GetPenaltyBoxConditions(ctx, getPenaltyBoxConditions)
	if err != nil {
		logger.Errorf("calling 'getPenaltyBoxConditions': %s", err.Error())
		return diag.FromErr(err)
	}

	if err := d.Set("config_id", configID); err != nil {
		return diag.Errorf("%s: %s", tf.ErrValueSet, err.Error())
	}
	if err := d.Set("security_policy_id", policyID); err != nil {
		return diag.Errorf("%s: %s", tf.ErrValueSet, err.Error())
	}
	jsonBody, err := json.Marshal(penaltyBoxConditions)
	if err != nil {
		return diag.FromErr(err)
	}
	if err := d.Set("penalty_box_conditions", string(jsonBody)); err != nil {
		return diag.Errorf("%s: %s", tf.ErrValueSet, err.Error())
	}

	return nil
}

func resourcePenaltyBoxConditionsUpdate(ctx context.Context, d *schema.ResourceData, m interface{}) diag.Diagnostics {
	meta := meta.Must(m)
	client := inst.Client(meta)
	logger := meta.Log("APPSEC", "resourcePenaltyBoxConditionsUpdate")
	logger.Debugf("in resourcePenaltyBoxConditionsUpdate")

	iDParts, err := splitID(d.Id(), 2, "configID:securityPolicyID")
	if err != nil {
		return diag.FromErr(err)
	}
	configID, err := strconv.Atoi(iDParts[0])
	if err != nil {
		return diag.FromErr(err)
	}
	version, err := getModifiableConfigVersion(ctx, configID, "penaltyBoxAction", m)
	if err != nil {
		return diag.FromErr(err)
	}
	policyID := iDParts[1]
	jsonPostPayload := d.Get("penalty_box_conditions")
	conditionsJSON := appsec.PenaltyBoxConditionsPayload{}
	err = json.Unmarshal([]byte(jsonPostPayload.(string)), &conditionsJSON)
	if err != nil {
		return diag.FromErr(err)
	}

	updatePenaltyBoxConditions := appsec.UpdatePenaltyBoxConditionsRequest{
		ConfigID:          configID,
		Version:           version,
		PolicyID:          policyID,
		ConditionsPayload: conditionsJSON,
	}

	_, err = client.UpdatePenaltyBoxConditions(ctx, updatePenaltyBoxConditions)
	if err != nil {
		logger.Errorf("calling 'updatePenaltyBox': %s", err.Error())
		return diag.FromErr(err)
	}

	return resourcePenaltyBoxConditionsRead(ctx, d, m)
}

func resourcePenaltyBoxConditionsDelete(ctx context.Context, d *schema.ResourceData, m interface{}) diag.Diagnostics {
	meta := meta.Must(m)
	client := inst.Client(meta)
	logger := meta.Log("APPSEC", "resourcePenaltyBoxConditionsDelete")
	logger.Debugf("in resourcePenaltyBoxConditionsDelete")

	iDParts, err := splitID(d.Id(), 2, "configID:securityPolicyID")
	if err != nil {
		return diag.FromErr(err)
	}
	configID, err := strconv.Atoi(iDParts[0])
	if err != nil {
		return diag.FromErr(err)
	}
	version, err := getModifiableConfigVersion(ctx, configID, "penaltyBoxAction", m)
	if err != nil {
		return diag.FromErr(err)
	}
	policyID := iDParts[1]

	conditionsJSON := appsec.PenaltyBoxConditionsPayload{
		ConditionOperator: "AND",
		Conditions:        &appsec.RuleConditions{},
	}

	updatePenaltyBoxConditions := appsec.UpdatePenaltyBoxConditionsRequest{
		ConfigID:          configID,
		Version:           version,
		PolicyID:          policyID,
		ConditionsPayload: conditionsJSON,
	}

	_, err = client.UpdatePenaltyBoxConditions(ctx, updatePenaltyBoxConditions)
	if err != nil {
		logger.Errorf("calling 'UpdatePenaltyBoxConditions': %s", err.Error())
		return diag.FromErr(err)
	}

	return nil

}
