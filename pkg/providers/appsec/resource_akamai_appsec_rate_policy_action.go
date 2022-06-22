package appsec

import (
	"context"
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
func resourceRatePolicyAction() *schema.Resource {
	return &schema.Resource{
		CreateContext: resourceRatePolicyActionCreate,
		ReadContext:   resourceRatePolicyActionRead,
		UpdateContext: resourceRatePolicyActionUpdate,
		DeleteContext: resourceRatePolicyActionDelete,
		Importer: &schema.ResourceImporter{
			StateContext: schema.ImportStatePassthroughContext,
		},
		CustomizeDiff: customdiff.All(
			VerifyIDUnchanged,
		),
		Schema: map[string]*schema.Schema{
			"config_id": {
				Type:     schema.TypeInt,
				Required: true,
			},
			"security_policy_id": {
				Type:     schema.TypeString,
				Required: true,
			},
			"rate_policy_id": {
				Type:     schema.TypeInt,
				Required: true,
			},
			"ipv4_action": {
				Type:             schema.TypeString,
				Required:         true,
				ValidateDiagFunc: ValidateActions,
			},
			"ipv6_action": {
				Type:             schema.TypeString,
				Required:         true,
				ValidateDiagFunc: ValidateActions},
		},
	}
}

func resourceRatePolicyActionCreate(ctx context.Context, d *schema.ResourceData, m interface{}) diag.Diagnostics {
	meta := akamai.Meta(m)
	client := inst.Client(meta)
	logger := meta.Log("APPSEC", "resourceRatePolicyActionCreate")
	logger.Debugf("in resourceRatePolicyActionCreate")

	configID, err := tools.GetIntValue("config_id", d)
	if err != nil {
		return diag.FromErr(err)
	}
	version, err := getModifiableConfigVersion(ctx, configID, "ratePolicyAction", m)
	if err != nil {
		return diag.FromErr(err)
	}
	securityPolicyID, err := tools.GetStringValue("security_policy_id", d)
	if err != nil {
		return diag.FromErr(err)
	}
	ratePolicyID, err := tools.GetIntValue("rate_policy_id", d)
	if err != nil {
		return diag.FromErr(err)
	}
	ipv4action, err := tools.GetStringValue("ipv4_action", d)
	if err != nil {
		return diag.FromErr(err)
	}
	ipv6action, err := tools.GetStringValue("ipv6_action", d)
	if err != nil {
		return diag.FromErr(err)
	}

	updateRatePolicyAction := appsec.UpdateRatePolicyActionRequest{
		ConfigID:     configID,
		Version:      version,
		PolicyID:     securityPolicyID,
		RatePolicyID: ratePolicyID,
		Ipv4Action:   ipv4action,
		Ipv6Action:   ipv6action,
	}

	_, err = client.UpdateRatePolicyAction(ctx, updateRatePolicyAction)
	if err != nil {
		logger.Errorf("calling 'updateRatePolicyAction': %s", err.Error())
		return diag.FromErr(err)
	}

	d.SetId(fmt.Sprintf("%d:%s:%d", configID, securityPolicyID, ratePolicyID))

	return resourceRatePolicyActionRead(ctx, d, m)
}

func resourceRatePolicyActionRead(ctx context.Context, d *schema.ResourceData, m interface{}) diag.Diagnostics {
	meta := akamai.Meta(m)
	client := inst.Client(meta)
	logger := meta.Log("APPSEC", "resourceRatePolicyActionRead")
	logger.Debugf("in resourceRatePolicyActionRead")

	iDParts, err := splitID(d.Id(), 3, "configID:securityPolicyID:ratePolicyID")
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
	securityPolicyID := iDParts[1]
	ratePolicyID, err := strconv.Atoi(iDParts[2])
	if err != nil {
		return diag.FromErr(err)
	}

	getRatePolicyActionsRequest := appsec.GetRatePolicyActionsRequest{
		ConfigID:     configID,
		Version:      version,
		PolicyID:     securityPolicyID,
		RatePolicyID: ratePolicyID,
	}

	ratepolicyactions, err := client.GetRatePolicyActions(ctx, getRatePolicyActionsRequest)
	if err != nil {
		logger.Errorf("calling 'getRatePolicyActions': %s", err.Error())
		return diag.FromErr(err)
	}

	if err := d.Set("config_id", configID); err != nil {
		return diag.Errorf("%s: %s", tools.ErrValueSet, err.Error())
	}
	if err := d.Set("security_policy_id", securityPolicyID); err != nil {
		return diag.Errorf("%s: %s", tools.ErrValueSet, err.Error())
	}
	if err := d.Set("rate_policy_id", ratePolicyID); err != nil {
		return diag.Errorf("%s: %s", tools.ErrValueSet, err.Error())
	}
	for _, action := range ratepolicyactions.RatePolicyActions {
		if action.ID == ratePolicyID {
			if err := d.Set("ipv4_action", action.Ipv4Action); err != nil {
				return diag.Errorf("%s: %s", tools.ErrValueSet, err.Error())
			}
			if err := d.Set("ipv6_action", action.Ipv6Action); err != nil {
				return diag.Errorf("%s: %s", tools.ErrValueSet, err.Error())
			}
		}
	}

	return nil
}

func resourceRatePolicyActionUpdate(ctx context.Context, d *schema.ResourceData, m interface{}) diag.Diagnostics {
	meta := akamai.Meta(m)
	client := inst.Client(meta)
	logger := meta.Log("APPSEC", "resourceRatePolicyActionUpdate")
	logger.Debugf("in resourceRatePolicyActionUpdate")

	iDParts, err := splitID(d.Id(), 3, "configID:securityPolicyID:ratePolicyID")
	if err != nil {
		return diag.FromErr(err)
	}
	configID, err := strconv.Atoi(iDParts[0])
	if err != nil {
		return diag.FromErr(err)
	}
	version, err := getModifiableConfigVersion(ctx, configID, "ratePolicyAction", m)
	if err != nil {
		return diag.FromErr(err)
	}
	securityPolicyID := iDParts[1]
	ratePolicyID, err := strconv.Atoi(iDParts[2])
	if err != nil {
		return diag.FromErr(err)
	}
	ipv4action, err := tools.GetStringValue("ipv4_action", d)
	if err != nil {
		return diag.FromErr(err)
	}
	ipv6action, err := tools.GetStringValue("ipv6_action", d)
	if err != nil {
		return diag.FromErr(err)
	}

	updateRatePolicyAction := appsec.UpdateRatePolicyActionRequest{
		ConfigID:     configID,
		Version:      version,
		PolicyID:     securityPolicyID,
		RatePolicyID: ratePolicyID,
		Ipv4Action:   ipv4action,
		Ipv6Action:   ipv6action,
	}

	_, err = client.UpdateRatePolicyAction(ctx, updateRatePolicyAction)
	if err != nil {
		logger.Errorf("calling 'updateRatePolicyAction': %s", err.Error())
		return diag.FromErr(err)
	}

	return resourceRatePolicyActionRead(ctx, d, m)
}

func resourceRatePolicyActionDelete(ctx context.Context, d *schema.ResourceData, m interface{}) diag.Diagnostics {
	meta := akamai.Meta(m)
	client := inst.Client(meta)
	logger := meta.Log("APPSEC", "resourceRatePolicyActionDelete")
	logger.Debugf("in resourceRatePolicyActionDelete")

	iDParts, err := splitID(d.Id(), 3, "configID:securityPolicyID:ratePolicyID")
	if err != nil {
		return diag.FromErr(err)
	}
	configID, err := strconv.Atoi(iDParts[0])
	if err != nil {
		return diag.FromErr(err)
	}
	version, err := getModifiableConfigVersion(ctx, configID, "ratePolicyAction", m)
	if err != nil {
		return diag.FromErr(err)
	}
	securityPolicyID := iDParts[1]
	ratePolicyID, err := strconv.Atoi(iDParts[2])
	if err != nil {
		return diag.FromErr(err)
	}

	deleteRatePolicyAction := appsec.UpdateRatePolicyActionRequest{
		ConfigID:     configID,
		Version:      version,
		PolicyID:     securityPolicyID,
		RatePolicyID: ratePolicyID,
		Ipv4Action:   "none",
		Ipv6Action:   "none",
	}

	_, err = client.UpdateRatePolicyAction(ctx, deleteRatePolicyAction)
	if err != nil {
		logger.Errorf("calling 'removeRatePolicyAction': %s", err.Error())
		return diag.FromErr(err)
	}
	return nil
}
