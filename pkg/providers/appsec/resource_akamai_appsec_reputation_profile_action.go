package appsec

import (
	"context"
	"errors"
	"fmt"
	"strconv"

	"github.com/akamai/AkamaiOPEN-edgegrid-golang/v2/pkg/appsec"
	"github.com/akamai/terraform-provider-akamai/v2/pkg/akamai"
	"github.com/akamai/terraform-provider-akamai/v2/pkg/tools"
	"github.com/hashicorp/terraform-plugin-sdk/v2/diag"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/customdiff"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/schema"
)

// appsec v1
//
// https://developer.akamai.com/api/cloud_security/application_security/v1.html
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
				Type:     schema.TypeInt,
				Required: true,
			},
			"security_policy_id": {
				Type:     schema.TypeString,
				Required: true,
			},
			"reputation_profile_id": {
				Type:     schema.TypeInt,
				Required: true,
			},
			"action": {
				Type:             schema.TypeString,
				Required:         true,
				ValidateDiagFunc: ValidateActions,
			},
		},
	}
}

func resourceReputationProfileActionCreate(ctx context.Context, d *schema.ResourceData, m interface{}) diag.Diagnostics {
	meta := akamai.Meta(m)
	client := inst.Client(meta)
	logger := meta.Log("APPSEC", "resourceReputationProfileActionCreate")
	logger.Debugf("in resourceReputationProfileActionCreate")

	configID, err := tools.GetIntValue("config_id", d)
	if err != nil && !errors.Is(err, tools.ErrNotFound) {
		return diag.FromErr(err)
	}
	version := getModifiableConfigVersion(ctx, configID, "reputationProfileAction", m)
	policyID, err := tools.GetStringValue("security_policy_id", d)
	if err != nil && !errors.Is(err, tools.ErrNotFound) {
		return diag.FromErr(err)
	}
	reputationprofileid, err := tools.GetIntValue("reputation_profile_id", d)
	if err != nil && !errors.Is(err, tools.ErrNotFound) {
		return diag.FromErr(err)
	}
	action, err := tools.GetStringValue("action", d)
	if err != nil && !errors.Is(err, tools.ErrNotFound) {
		return diag.FromErr(err)
	}

	createReputationProfileAction := appsec.UpdateReputationProfileActionRequest{
		ConfigID:            configID,
		Version:             version,
		PolicyID:            policyID,
		ReputationProfileID: reputationprofileid,
		Action:              action,
	}

	_, erru := client.UpdateReputationProfileAction(ctx, createReputationProfileAction)
	if erru != nil {
		logger.Errorf("calling 'createReputationProfileAction': %s", erru.Error())
		return diag.FromErr(erru)
	}

	d.SetId(fmt.Sprintf("%d:%s:%d", createReputationProfileAction.ConfigID, createReputationProfileAction.PolicyID, createReputationProfileAction.ReputationProfileID))

	return resourceReputationProfileActionRead(ctx, d, m)
}

func resourceReputationProfileActionRead(ctx context.Context, d *schema.ResourceData, m interface{}) diag.Diagnostics {
	meta := akamai.Meta(m)
	client := inst.Client(meta)
	logger := meta.Log("APPSEC", "resourceReputationProfileActionRead")
	logger.Debugf("in resourceReputationProfileActionRead")

	idParts, err := splitID(d.Id(), 3, "configID:securityPolicyID:reputationprofileid")
	if err != nil {
		return diag.FromErr(err)
	}
	configID, err := strconv.Atoi(idParts[0])
	if err != nil {
		return diag.FromErr(err)
	}
	version := getLatestConfigVersion(ctx, configID, m)
	policyID := idParts[1]
	reputationprofileid, err := strconv.Atoi(idParts[2])
	if err != nil {
		return diag.FromErr(err)
	}

	getReputationProfileAction := appsec.GetReputationProfileActionRequest{
		ConfigID:            configID,
		Version:             version,
		PolicyID:            policyID,
		ReputationProfileID: reputationprofileid,
	}

	resp, errr := client.GetReputationProfileAction(ctx, getReputationProfileAction)
	if errr != nil {
		logger.Errorf("calling 'getReputationProfileAction': %s", errr.Error())
		return diag.FromErr(errr)
	}

	if err := d.Set("config_id", getReputationProfileAction.ConfigID); err != nil {
		return diag.Errorf("%s: %s", tools.ErrValueSet, err.Error())
	}
	if err := d.Set("security_policy_id", getReputationProfileAction.PolicyID); err != nil {
		return diag.Errorf("%s: %s", tools.ErrValueSet, err.Error())
	}
	if err := d.Set("reputation_profile_id", getReputationProfileAction.ReputationProfileID); err != nil {
		return diag.Errorf("%s: %s", tools.ErrValueSet, err.Error())
	}
	if err := d.Set("action", resp.Action); err != nil {
		return diag.Errorf("%s: %s", tools.ErrValueSet, err.Error())
	}

	return nil
}

func resourceReputationProfileActionUpdate(ctx context.Context, d *schema.ResourceData, m interface{}) diag.Diagnostics {
	meta := akamai.Meta(m)
	client := inst.Client(meta)
	logger := meta.Log("APPSEC", "resourceReputationProfileActionUpdate")
	logger.Debugf("in resourceReputationProfileActionUpdate")

	idParts, err := splitID(d.Id(), 3, "configID:securityPolicyID:reputationrofileid")
	if err != nil {
		return diag.FromErr(err)
	}
	configID, err := strconv.Atoi(idParts[0])
	if err != nil {
		return diag.FromErr(err)
	}
	version := getModifiableConfigVersion(ctx, configID, "reputationProfileAction", m)
	policyID := idParts[1]
	reputationprofileid, err := strconv.Atoi(idParts[2])
	if err != nil {
		return diag.FromErr(err)
	}
	action, err := tools.GetStringValue("action", d)
	if err != nil && !errors.Is(err, tools.ErrNotFound) {
		return diag.FromErr(err)
	}

	updateReputationProfileAction := appsec.UpdateReputationProfileActionRequest{
		ConfigID:            configID,
		Version:             version,
		PolicyID:            policyID,
		ReputationProfileID: reputationprofileid,
		Action:              action,
	}

	_, erru := client.UpdateReputationProfileAction(ctx, updateReputationProfileAction)
	if erru != nil {
		logger.Errorf("calling 'updateReputationProfileAction': %s", erru.Error())
		return diag.FromErr(erru)
	}

	return resourceReputationProfileActionRead(ctx, d, m)
}

func resourceReputationProfileActionDelete(ctx context.Context, d *schema.ResourceData, m interface{}) diag.Diagnostics {
	meta := akamai.Meta(m)
	client := inst.Client(meta)
	logger := meta.Log("APPSEC", "resourceReputationProfileActionDelete")
	logger.Debugf("in resourceReputationProfileActionDelete")

	idParts, err := splitID(d.Id(), 3, "configID:securityPolicyID:reputationprofileid")
	if err != nil {
		return diag.FromErr(err)
	}
	configID, err := strconv.Atoi(idParts[0])
	if err != nil {
		return diag.FromErr(err)
	}
	version := getModifiableConfigVersion(ctx, configID, "reputationProfileAction", m)
	policyID := idParts[1]
	reputationprofileid, err := strconv.Atoi(idParts[2])
	if err != nil {
		return diag.FromErr(err)
	}

	removeReputationProfileAction := appsec.UpdateReputationProfileActionRequest{
		ConfigID:            configID,
		Version:             version,
		PolicyID:            policyID,
		ReputationProfileID: reputationprofileid,
		Action:              "none",
	}

	_, errd := client.UpdateReputationProfileAction(ctx, removeReputationProfileAction)
	if errd != nil {
		logger.Errorf("calling 'removeReputationProfileAction': %s", errd.Error())
		return diag.FromErr(errd)
	}

	d.SetId("")

	return nil
}
