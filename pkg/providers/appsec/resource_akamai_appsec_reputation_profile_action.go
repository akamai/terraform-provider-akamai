package appsec

import (
	"context"
	"errors"
	"fmt"
	"strconv"

	"github.com/akamai/AkamaiOPEN-edgegrid-golang/v7/pkg/appsec"
	"github.com/akamai/terraform-provider-akamai/v5/pkg/common/tf"
	"github.com/akamai/terraform-provider-akamai/v5/pkg/meta"
	"github.com/hashicorp/terraform-plugin-sdk/v2/diag"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/customdiff"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/schema"
)

// appsec v1
//
// https://techdocs.akamai.com/application-security/reference/api
func resourceReputationProfileAction() *schema.Resource {
	return &schema.Resource{
		CreateContext: resourceReputationProfileActionCreate,
		ReadContext:   resourceReputationProfileActionRead,
		UpdateContext: resourceReputationProfileActionUpdate,
		DeleteContext: resourceReputationProfileActionDelete,
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
			"reputation_profile_id": {
				Type:        schema.TypeInt,
				Required:    true,
				Description: "Unique identifier of the reputation profile",
			},
			"action": {
				Type:             schema.TypeString,
				Required:         true,
				ValidateDiagFunc: ValidateActions,
				Description:      "Action to be taken when the reputation profile is triggered",
			},
		},
	}
}

func resourceReputationProfileActionCreate(ctx context.Context, d *schema.ResourceData, m interface{}) diag.Diagnostics {
	meta := meta.Must(m)
	client := inst.Client(meta)
	logger := meta.Log("APPSEC", "resourceReputationProfileActionCreate")
	logger.Debugf("in resourceReputationProfileActionCreate")

	configID, err := tf.GetIntValue("config_id", d)
	if err != nil {
		return diag.FromErr(err)
	}
	version, err := getModifiableConfigVersion(ctx, configID, "reputationProfileAction", m)
	if err != nil {
		return diag.FromErr(err)
	}
	policyID, err := tf.GetStringValue("security_policy_id", d)
	if err != nil {
		return diag.FromErr(err)
	}
	reputationProfileID, err := tf.GetIntValue("reputation_profile_id", d)
	if err != nil {
		return diag.FromErr(err)
	}
	action, err := tf.GetStringValue("action", d)
	if err != nil {
		return diag.FromErr(err)
	}

	createReputationProfileAction := appsec.UpdateReputationProfileActionRequest{
		ConfigID:            configID,
		Version:             version,
		PolicyID:            policyID,
		ReputationProfileID: reputationProfileID,
		Action:              action,
	}

	_, err = client.UpdateReputationProfileAction(ctx, createReputationProfileAction)
	if err != nil {
		logger.Errorf("calling 'createReputationProfileAction': %s", err.Error())
		return diag.FromErr(err)
	}

	d.SetId(fmt.Sprintf("%d:%s:%d", createReputationProfileAction.ConfigID, createReputationProfileAction.PolicyID, createReputationProfileAction.ReputationProfileID))

	return resourceReputationProfileActionRead(ctx, d, m)
}

func resourceReputationProfileActionRead(ctx context.Context, d *schema.ResourceData, m interface{}) diag.Diagnostics {
	meta := meta.Must(m)
	client := inst.Client(meta)
	logger := meta.Log("APPSEC", "resourceReputationProfileActionRead")
	logger.Debugf("in resourceReputationProfileActionRead")

	iDParts, err := splitID(d.Id(), 3, "configID:securityPolicyID:reputationProfileID")
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
	reputationProfileID, err := strconv.Atoi(iDParts[2])
	if err != nil {
		return diag.FromErr(err)
	}

	getReputationProfileAction := appsec.GetReputationProfileActionRequest{
		ConfigID:            configID,
		Version:             version,
		PolicyID:            policyID,
		ReputationProfileID: reputationProfileID,
	}

	resp, err := client.GetReputationProfileAction(ctx, getReputationProfileAction)
	if err != nil {
		logger.Errorf("calling 'getReputationProfileAction': %s", err.Error())
		return diag.FromErr(err)
	}

	if err := d.Set("config_id", getReputationProfileAction.ConfigID); err != nil {
		return diag.Errorf("%s: %s", tf.ErrValueSet, err.Error())
	}
	if err := d.Set("security_policy_id", getReputationProfileAction.PolicyID); err != nil {
		return diag.Errorf("%s: %s", tf.ErrValueSet, err.Error())
	}
	if err := d.Set("reputation_profile_id", getReputationProfileAction.ReputationProfileID); err != nil {
		return diag.Errorf("%s: %s", tf.ErrValueSet, err.Error())
	}
	if err := d.Set("action", resp.Action); err != nil {
		return diag.Errorf("%s: %s", tf.ErrValueSet, err.Error())
	}

	return nil
}

func resourceReputationProfileActionUpdate(ctx context.Context, d *schema.ResourceData, m interface{}) diag.Diagnostics {
	meta := meta.Must(m)
	client := inst.Client(meta)
	logger := meta.Log("APPSEC", "resourceReputationProfileActionUpdate")
	logger.Debugf("in resourceReputationProfileActionUpdate")

	iDParts, err := splitID(d.Id(), 3, "configID:securityPolicyID:reputationProfileID")
	if err != nil {
		return diag.FromErr(err)
	}
	configID, err := strconv.Atoi(iDParts[0])
	if err != nil {
		return diag.FromErr(err)
	}
	version, err := getModifiableConfigVersion(ctx, configID, "reputationProfileAction", m)
	if err != nil {
		return diag.FromErr(err)
	}
	policyID := iDParts[1]
	reputationProfileID, err := strconv.Atoi(iDParts[2])
	if err != nil {
		return diag.FromErr(err)
	}
	action, err := tf.GetStringValue("action", d)
	if err != nil && !errors.Is(err, tf.ErrNotFound) {
		return diag.FromErr(err)
	}

	updateReputationProfileAction := appsec.UpdateReputationProfileActionRequest{
		ConfigID:            configID,
		Version:             version,
		PolicyID:            policyID,
		ReputationProfileID: reputationProfileID,
		Action:              action,
	}

	_, err = client.UpdateReputationProfileAction(ctx, updateReputationProfileAction)
	if err != nil {
		logger.Errorf("calling 'updateReputationProfileAction': %s", err.Error())
		return diag.FromErr(err)
	}

	return resourceReputationProfileActionRead(ctx, d, m)
}

func resourceReputationProfileActionDelete(ctx context.Context, d *schema.ResourceData, m interface{}) diag.Diagnostics {
	meta := meta.Must(m)
	client := inst.Client(meta)
	logger := meta.Log("APPSEC", "resourceReputationProfileActionDelete")
	logger.Debugf("in resourceReputationProfileActionDelete")

	iDParts, err := splitID(d.Id(), 3, "configID:securityPolicyID:reputationProfileID")
	if err != nil {
		return diag.FromErr(err)
	}
	configID, err := strconv.Atoi(iDParts[0])
	if err != nil {
		return diag.FromErr(err)
	}
	version, err := getModifiableConfigVersion(ctx, configID, "reputationProfileAction", m)
	if err != nil {
		return diag.FromErr(err)
	}
	policyID := iDParts[1]
	reputationProfileID, err := strconv.Atoi(iDParts[2])
	if err != nil {
		return diag.FromErr(err)
	}

	removeReputationProfileAction := appsec.UpdateReputationProfileActionRequest{
		ConfigID:            configID,
		Version:             version,
		PolicyID:            policyID,
		ReputationProfileID: reputationProfileID,
		Action:              "none",
	}

	_, err = client.UpdateReputationProfileAction(ctx, removeReputationProfileAction)
	if err != nil {
		logger.Errorf("calling 'removeReputationProfileAction': %s", err.Error())
		return diag.FromErr(err)
	}
	return nil
}
