package appsec

import (
	"context"
	"fmt"
	"strconv"

	"github.com/akamai/AkamaiOPEN-edgegrid-golang/v10/pkg/appsec"
	"github.com/akamai/terraform-provider-akamai/v7/pkg/common/id"
	"github.com/akamai/terraform-provider-akamai/v7/pkg/common/tf"
	"github.com/akamai/terraform-provider-akamai/v7/pkg/meta"
	"github.com/hashicorp/terraform-plugin-sdk/v2/diag"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/customdiff"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/schema"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/validation"
)

// appsec v1
//
// https://techdocs.akamai.com/application-security/reference/api
func resourceEvalPenaltyBox() *schema.Resource {
	return &schema.Resource{
		CreateContext: resourceEvalPenaltyBoxCreate,
		ReadContext:   resourceEvalPenaltyBoxRead,
		UpdateContext: resourceEvalPenaltyBoxUpdate,
		DeleteContext: resourceEvalPenaltyBoxDelete,
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
			"penalty_box_protection": {
				Type:        schema.TypeBool,
				Required:    true,
				Description: "Whether to enable the penalty box for the specified security policy",
			},
			"penalty_box_action": {
				Type:        schema.TypeString,
				Required:    true,
				Description: "Action applied to requests from clients in the penalty box",
				ValidateDiagFunc: validation.ToDiagFunc(validation.StringInSlice([]string{
					string(appsec.ActionTypeDeny),
					string(appsec.ActionTypeAlert),
					string(appsec.ActionTypeNone),
				}, false)),
			},
		},
	}
}

func resourceEvalPenaltyBoxCreate(ctx context.Context, d *schema.ResourceData, m interface{}) diag.Diagnostics {
	meta := meta.Must(m)
	client := inst.Client(meta)
	logger := meta.Log("APPSEC", "resourceEvalPenaltyBoxCreate")
	logger.Debugf("in resourceEvalPenaltyBoxCreate")

	configID, err := tf.GetIntValue("config_id", d)
	if err != nil {
		return diag.FromErr(err)
	}
	version, err := getModifiableConfigVersion(ctx, configID, "evalPenaltyBox", m)
	if err != nil {
		return diag.FromErr(err)
	}
	policyID, err := tf.GetStringValue("security_policy_id", d)
	if err != nil {
		return diag.FromErr(err)
	}
	penaltyboxprotection, err := tf.GetBoolValue("penalty_box_protection", d)
	if err != nil {
		return diag.FromErr(err)
	}
	penaltyboxaction, err := tf.GetStringValue("penalty_box_action", d)
	if err != nil {
		return diag.FromErr(err)
	}

	createPenaltyBox := appsec.UpdatePenaltyBoxRequest{
		ConfigID:             configID,
		Version:              version,
		PolicyID:             policyID,
		PenaltyBoxProtection: penaltyboxprotection,
		Action:               penaltyboxaction,
	}

	_, err = client.UpdateEvalPenaltyBox(ctx, createPenaltyBox)
	if err != nil {
		logger.Errorf("calling 'createEvalPenaltyBox': %s", err.Error())
		return diag.FromErr(err)
	}

	d.SetId(fmt.Sprintf("%d:%s", createPenaltyBox.ConfigID, createPenaltyBox.PolicyID))

	return resourceEvalPenaltyBoxRead(ctx, d, m)
}

func resourceEvalPenaltyBoxRead(ctx context.Context, d *schema.ResourceData, m interface{}) diag.Diagnostics {
	meta := meta.Must(m)
	client := inst.Client(meta)
	logger := meta.Log("APPSEC", "resourceEvalPenaltyBoxRead")
	logger.Debugf("in resourceEvalPenaltyBoxRead")

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

	getPenaltyBox := appsec.GetPenaltyBoxRequest{
		ConfigID: configID,
		Version:  version,
		PolicyID: policyID,
	}

	penaltybox, err := client.GetEvalPenaltyBox(ctx, getPenaltyBox)
	if err != nil {
		logger.Errorf("calling 'getPenaltyBox': %s", err.Error())
		return diag.FromErr(err)
	}

	if err := d.Set("config_id", getPenaltyBox.ConfigID); err != nil {
		return diag.Errorf("%s: %s", tf.ErrValueSet, err.Error())
	}
	if err := d.Set("security_policy_id", getPenaltyBox.PolicyID); err != nil {
		return diag.Errorf("%s: %s", tf.ErrValueSet, err.Error())
	}
	if err := d.Set("penalty_box_protection", penaltybox.PenaltyBoxProtection); err != nil {
		return diag.Errorf("%s: %s", tf.ErrValueSet, err.Error())
	}
	if err := d.Set("penalty_box_action", penaltybox.Action); err != nil {
		return diag.Errorf("%s: %s", tf.ErrValueSet, err.Error())
	}

	return nil
}

func resourceEvalPenaltyBoxUpdate(ctx context.Context, d *schema.ResourceData, m interface{}) diag.Diagnostics {
	meta := meta.Must(m)
	client := inst.Client(meta)
	logger := meta.Log("APPSEC", "resourceEvalPenaltyBoxUpdate")
	logger.Debugf("in resourceEvalPenaltyBoxUpdate")

	iDParts, err := id.Split(d.Id(), 2, "configID:securityPolicyID")
	if err != nil {
		return diag.FromErr(err)
	}
	configID, err := strconv.Atoi(iDParts[0])
	if err != nil {
		return diag.FromErr(err)
	}
	version, err := getModifiableConfigVersion(ctx, configID, "evalPenaltyBox", m)
	if err != nil {
		return diag.FromErr(err)
	}
	policyID := iDParts[1]
	penaltyboxprotection, err := tf.GetBoolValue("penalty_box_protection", d)
	if err != nil {
		return diag.FromErr(err)
	}
	penaltyboxaction, err := tf.GetStringValue("penalty_box_action", d)
	if err != nil {
		return diag.FromErr(err)
	}

	updatePenaltyBox := appsec.UpdatePenaltyBoxRequest{
		ConfigID:             configID,
		Version:              version,
		PolicyID:             policyID,
		PenaltyBoxProtection: penaltyboxprotection,
		Action:               penaltyboxaction,
	}

	_, err = client.UpdateEvalPenaltyBox(ctx, updatePenaltyBox)
	if err != nil {
		logger.Errorf("calling 'updateEvalPenaltyBox': %s", err.Error())
		return diag.FromErr(err)
	}

	return resourceEvalPenaltyBoxRead(ctx, d, m)
}

func resourceEvalPenaltyBoxDelete(ctx context.Context, d *schema.ResourceData, m interface{}) diag.Diagnostics {
	meta := meta.Must(m)
	client := inst.Client(meta)
	logger := meta.Log("APPSEC", "resourceEvalPenaltyBoxDelete")
	logger.Debugf("in resourceEvalPenaltyBoxDelete")

	iDParts, err := id.Split(d.Id(), 2, "configID:securityPolicyID")
	if err != nil {
		return diag.FromErr(err)
	}
	configID, err := strconv.Atoi(iDParts[0])
	if err != nil {
		return diag.FromErr(err)
	}
	version, err := getModifiableConfigVersion(ctx, configID, "evalPenaltyBox", m)
	if err != nil {
		return diag.FromErr(err)
	}
	policyID := iDParts[1]

	removePenaltyBox := appsec.UpdatePenaltyBoxRequest{
		ConfigID:             configID,
		Version:              version,
		PolicyID:             policyID,
		PenaltyBoxProtection: false,
		Action:               string(appsec.ActionTypeNone),
	}

	_, err = client.UpdateEvalPenaltyBox(ctx, removePenaltyBox)
	if err != nil {
		logger.Errorf("calling 'removeEvalPenaltyBox': %s", err.Error())
		return diag.FromErr(err)
	}
	d.SetId("")

	return nil
}
