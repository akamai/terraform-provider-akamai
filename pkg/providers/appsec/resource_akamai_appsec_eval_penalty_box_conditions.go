package appsec

import (
	"context"
	"encoding/json"
	"fmt"
	"strconv"

	"github.com/akamai/AkamaiOPEN-edgegrid-golang/v12/pkg/appsec"
	"github.com/akamai/terraform-provider-akamai/v9/pkg/common/id"
	"github.com/akamai/terraform-provider-akamai/v9/pkg/common/tf"
	"github.com/akamai/terraform-provider-akamai/v9/pkg/meta"
	"github.com/hashicorp/terraform-plugin-sdk/v2/diag"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/customdiff"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/schema"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/validation"
)

// appsec v1
//
// https://techdocs.akamai.com/application-security/reference/api
func resourceEvalPenaltyBoxConditions() *schema.Resource {
	return &schema.Resource{
		CreateContext: resourceEvalPenaltyBoxConditionsCreate,
		ReadContext:   resourceEvalPenaltyBoxConditionsRead,
		UpdateContext: resourceEvalPenaltyBoxConditionsUpdate,
		DeleteContext: resourceEvalPenaltyBoxConditionsDelete,
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
				Description:      "Description of evaluation penalty box conditions",
			},
		},
	}
}

func resourceEvalPenaltyBoxConditionsCreate(ctx context.Context, d *schema.ResourceData, m interface{}) diag.Diagnostics {
	must := meta.Must(m)
	client := inst.Client(must)
	logger := must.Log("APPSEC", "resourceEvalPenaltyBoxConditionsCreate")
	logger.Debugf("in resourceEvalPenaltyBoxConditionsCreate")

	configID, err := tf.GetIntValue("config_id", d)
	if err != nil {
		return diag.FromErr(err)
	}
	version, err := getModifiableConfigVersion(ctx, configID, "evalPenaltyBoxConditions", m)
	if err != nil {
		return diag.FromErr(err)
	}
	policyID, err := tf.GetStringValue("security_policy_id", d)
	if err != nil {
		return diag.FromErr(err)
	}

	jsonPostPayload := d.Get("penalty_box_conditions")
	conditionsPayload := appsec.PenaltyBoxConditionsPayload{}
	err = json.Unmarshal([]byte(jsonPostPayload.(string)), &conditionsPayload)
	if err != nil {
		return diag.FromErr(err)
	}

	createPenaltyBoxConditions := appsec.UpdatePenaltyBoxConditionsRequest{
		ConfigID:          configID,
		Version:           version,
		PolicyID:          policyID,
		ConditionsPayload: conditionsPayload,
	}

	_, err = client.UpdateEvalPenaltyBoxConditions(ctx, createPenaltyBoxConditions)
	if err != nil {
		logger.Errorf("calling 'createEvalPenaltyBoxConditions': %s", err.Error())
		return diag.FromErr(err)
	}

	d.SetId(fmt.Sprintf("%d:%s", createPenaltyBoxConditions.ConfigID, createPenaltyBoxConditions.PolicyID))

	return resourceEvalPenaltyBoxConditionsRead(ctx, d, m)
}

func resourceEvalPenaltyBoxConditionsRead(ctx context.Context, d *schema.ResourceData, m interface{}) diag.Diagnostics {
	must := meta.Must(m)
	client := inst.Client(must)
	logger := must.Log("APPSEC", "resourceEvalPenaltyBoxConditionsRead")
	logger.Debugf("in resourceEvalPenaltyBoxConditionsRead")

	iDParts, err := id.Split(d.Id(), 2, "configID:securityPolicyID")
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

	penaltyBoxConditions, err := client.GetEvalPenaltyBoxConditions(ctx, getPenaltyBoxConditions)
	if err != nil {
		logger.Errorf("calling 'GetEvalPenaltyBoxConditions': %s", err.Error())
		return diag.FromErr(err)
	}

	if err := d.Set("config_id", getPenaltyBoxConditions.ConfigID); err != nil {
		return diag.Errorf("%s: %s", tf.ErrValueSet, err.Error())
	}
	if err := d.Set("security_policy_id", getPenaltyBoxConditions.PolicyID); err != nil {
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

func resourceEvalPenaltyBoxConditionsUpdate(ctx context.Context, d *schema.ResourceData, m interface{}) diag.Diagnostics {
	must := meta.Must(m)
	client := inst.Client(must)
	logger := must.Log("APPSEC", "resourceEvalPenaltyBoxConditionsUpdate")
	logger.Debugf("in resourceEvalPenaltyBoxConditionsUpdate")

	iDParts, err := id.Split(d.Id(), 2, "configID:securityPolicyID")
	if err != nil {
		return diag.FromErr(err)
	}
	configID, err := strconv.Atoi(iDParts[0])
	if err != nil {
		return diag.FromErr(err)
	}
	version, err := getModifiableConfigVersion(ctx, configID, "evalPenaltyBoxConditions", m)
	if err != nil {
		return diag.FromErr(err)
	}
	policyID := iDParts[1]

	jsonPostPayload := d.Get("penalty_box_conditions")
	conditionsPayload := appsec.PenaltyBoxConditionsPayload{}
	err = json.Unmarshal([]byte(jsonPostPayload.(string)), &conditionsPayload)
	if err != nil {
		return diag.FromErr(err)
	}

	updatePenaltyBoxConditions := appsec.UpdatePenaltyBoxConditionsRequest{
		ConfigID:          configID,
		Version:           version,
		PolicyID:          policyID,
		ConditionsPayload: conditionsPayload,
	}

	_, err = client.UpdateEvalPenaltyBoxConditions(ctx, updatePenaltyBoxConditions)
	if err != nil {
		logger.Errorf("calling 'updateEvalPenaltyBoxConditions': %s", err.Error())
		return diag.FromErr(err)
	}

	return resourceEvalPenaltyBoxConditionsRead(ctx, d, m)
}

func resourceEvalPenaltyBoxConditionsDelete(ctx context.Context, d *schema.ResourceData, m interface{}) diag.Diagnostics {
	meta := meta.Must(m)
	client := inst.Client(meta)
	logger := meta.Log("APPSEC", "resourceEvalPenaltyBoxConditionsUpdate")
	logger.Debugf("in resourceEvalPenaltyBoxConditionsDelete")

	iDParts, err := id.Split(d.Id(), 2, "configID:securityPolicyID")
	if err != nil {
		return diag.FromErr(err)
	}
	configID, err := strconv.Atoi(iDParts[0])
	if err != nil {
		return diag.FromErr(err)
	}
	version, err := getModifiableConfigVersion(ctx, configID, "evalPenaltyBoxConditions", m)
	if err != nil {
		return diag.FromErr(err)
	}
	policyID := iDParts[1]

	conditionsPayload := appsec.PenaltyBoxConditionsPayload{
		ConditionOperator: "AND",
		Conditions:        &appsec.RuleConditions{},
	}

	removePenaltyBoxConditions := appsec.UpdatePenaltyBoxConditionsRequest{
		ConfigID:          configID,
		Version:           version,
		PolicyID:          policyID,
		ConditionsPayload: conditionsPayload,
	}

	_, err = client.UpdateEvalPenaltyBoxConditions(ctx, removePenaltyBoxConditions)
	if err != nil {
		logger.Errorf("calling 'removeEvalPenaltyBoxConditions': %s", err.Error())
		return diag.FromErr(err)
	}

	return nil
}
