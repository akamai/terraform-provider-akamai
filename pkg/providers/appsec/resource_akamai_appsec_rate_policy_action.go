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
				Type:         schema.TypeString,
				Required:     true,
				ValidateFunc: ValidateActions,
			},
			"ipv6_action": {
				Type:         schema.TypeString,
				Required:     true,
				ValidateFunc: ValidateActions},
		},
	}
}

func resourceRatePolicyActionCreate(ctx context.Context, d *schema.ResourceData, m interface{}) diag.Diagnostics {
	meta := akamai.Meta(m)
	client := inst.Client(meta)
	logger := meta.Log("APPSEC", "resourceRatePolicyActionUpdate")
	logger.Debugf("!!! in resourceRatePolicyActionCreate")

	configid, err := tools.GetIntValue("config_id", d)
	if err != nil {
		return diag.FromErr(err)
	}
	version := getModifiableConfigVersion(ctx, configid, "ratePolicyAction", m)
	securitypolicyid, err := tools.GetStringValue("security_policy_id", d)
	if err != nil {
		return diag.FromErr(err)
	}
	ratepolicyid, err := tools.GetIntValue("rate_policy_id", d)
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

	updateRatePolicyAction := appsec.UpdateRatePolicyActionRequest{}
	updateRatePolicyAction.ConfigID = configid
	updateRatePolicyAction.Version = version
	updateRatePolicyAction.PolicyID = securitypolicyid
	updateRatePolicyAction.RatePolicyID = ratepolicyid
	updateRatePolicyAction.Ipv4Action = ipv4action
	updateRatePolicyAction.Ipv6Action = ipv6action

	_, err = client.UpdateRatePolicyAction(ctx, updateRatePolicyAction)
	if err != nil {
		logger.Errorf("calling 'updateRatePolicyAction': %s", err.Error())
		return diag.FromErr(err)
	}

	d.SetId(fmt.Sprintf("%d:%s:%d", configid, securitypolicyid, ratepolicyid))

	return resourceRatePolicyActionRead(ctx, d, m)
}

func resourceRatePolicyActionRead(ctx context.Context, d *schema.ResourceData, m interface{}) diag.Diagnostics {
	meta := akamai.Meta(m)
	client := inst.Client(meta)
	logger := meta.Log("APPSEC", "resourceRatePolicyActionRead")
	logger.Debugf("!!! in resourceRatePolicyActionRead")

	idParts, err := splitID(d.Id(), 3, "configid:securitypolicyid:ratepolicyid")
	if err != nil {
		return diag.FromErr(err)
	}

	configid, err := strconv.Atoi(idParts[0])
	if err != nil {
		return diag.FromErr(err)
	}
	version := getLatestConfigVersion(ctx, configid, m)
	securitypolicyid := idParts[1]
	ratepolicyid, err := strconv.Atoi(idParts[2])
	if err != nil {
		return diag.FromErr(err)
	}

	readRatePolicyAction := appsec.GetRatePolicyActionRequest{}
	readRatePolicyAction.ConfigID = configid
	readRatePolicyAction.Version = version
	readRatePolicyAction.PolicyID = securitypolicyid
	readRatePolicyAction.ID = ratepolicyid

	ratepolicyaction, err := client.GetRatePolicyAction(ctx, readRatePolicyAction)
	if err != nil {
		logger.Errorf("calling 'getRatePolicyAction': %s", err.Error())
		return diag.FromErr(err)
	}

	if err := d.Set("config_id", configid); err != nil {
		return diag.FromErr(fmt.Errorf("%w: %s", tools.ErrValueSet, err.Error()))
	}
	if err := d.Set("security_policy_id", securitypolicyid); err != nil {
		return diag.FromErr(fmt.Errorf("%w: %s", tools.ErrValueSet, err.Error()))
	}
	if err := d.Set("rate_policy_id", ratepolicyid); err != nil {
		return diag.FromErr(fmt.Errorf("%w: %s", tools.ErrValueSet, err.Error()))
	}
	for _, action := range ratepolicyaction.RatePolicyActions {
		if action.ID == ratepolicyid {
			if err := d.Set("ipv4_action", action.Ipv4Action); err != nil {
				return diag.FromErr(fmt.Errorf("%w: %s", tools.ErrValueSet, err.Error()))
			}
			if err := d.Set("ipv6_action", action.Ipv6Action); err != nil {
				return diag.FromErr(fmt.Errorf("%w: %s", tools.ErrValueSet, err.Error()))
			}
		}
	}

	return nil
}

func resourceRatePolicyActionUpdate(ctx context.Context, d *schema.ResourceData, m interface{}) diag.Diagnostics {
	meta := akamai.Meta(m)
	client := inst.Client(meta)
	logger := meta.Log("APPSEC", "resourceRatePolicyActionUpdate")
	logger.Debugf("!!! in resourceRatePolicyActionUpdate")

	idParts, err := splitID(d.Id(), 3, "configid:securitypolicyid:ratepolicyid")
	if err != nil {
		return diag.FromErr(err)
	}
	configid, err := strconv.Atoi(idParts[0])
	if err != nil {
		return diag.FromErr(err)
	}
	version := getModifiableConfigVersion(ctx, configid, "ratePolicyAction", m)
	securitypolicyid := idParts[1]
	ratepolicyid, err := strconv.Atoi(idParts[2])
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

	updateRatePolicyAction := appsec.UpdateRatePolicyActionRequest{}
	updateRatePolicyAction.ConfigID = configid
	updateRatePolicyAction.Version = version
	updateRatePolicyAction.PolicyID = securitypolicyid
	updateRatePolicyAction.RatePolicyID = ratepolicyid
	updateRatePolicyAction.Ipv4Action = ipv4action
	updateRatePolicyAction.Ipv6Action = ipv6action

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
	logger.Debugf("!!! in resourceRatePolicyActionDelete")

	idParts, err := splitID(d.Id(), 3, "configid:securitypolicyid:ratepolicyid")
	if err != nil {
		return diag.FromErr(err)
	}
	configid, err := strconv.Atoi(idParts[0])
	if err != nil {
		return diag.FromErr(err)
	}
	version := getModifiableConfigVersion(ctx, configid, "ratePolicyAction", m)
	securitypolicyid := idParts[1]
	ratepolicyid, err := strconv.Atoi(idParts[2])
	if err != nil {
		return diag.FromErr(err)
	}

	deleteRatePolicyAction := appsec.UpdateRatePolicyActionRequest{}
	deleteRatePolicyAction.ConfigID = configid
	deleteRatePolicyAction.Version = version
	deleteRatePolicyAction.PolicyID = securitypolicyid
	deleteRatePolicyAction.RatePolicyID = ratepolicyid
	deleteRatePolicyAction.Ipv4Action = "none"
	deleteRatePolicyAction.Ipv6Action = "none"

	_, err = client.UpdateRatePolicyAction(ctx, deleteRatePolicyAction)
	if err != nil {
		logger.Errorf("calling 'removeRatePolicyAction': %s", err.Error())
		return diag.FromErr(err)
	}
	d.SetId("")

	return nil
}
