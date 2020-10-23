package appsec

import (
	"context"
	"errors"
	"strconv"

	v2 "github.com/akamai/AkamaiOPEN-edgegrid-golang/v2/pkg/appsec"
	"github.com/akamai/terraform-provider-akamai/v2/pkg/akamai"
	"github.com/akamai/terraform-provider-akamai/v2/pkg/tools"

	"github.com/hashicorp/terraform-plugin-sdk/v2/diag"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/schema"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/validation"
)

// appsec v1
//
// https://developer.akamai.com/api/cloud_security/application_security/v1.html
func resourceRatePolicyAction() *schema.Resource {
	return &schema.Resource{
		CreateContext: resourceRatePolicyActionUpdate,
		ReadContext:   resourceRatePolicyActionRead,
		UpdateContext: resourceRatePolicyActionUpdate,
		DeleteContext: resourceRatePolicyActionDelete,
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

func resourceRatePolicyActionRead(ctx context.Context, d *schema.ResourceData, m interface{}) diag.Diagnostics {
	meta := akamai.Meta(m)
	client := inst.Client(meta)
	logger := meta.Log("APPSEC", "resourceRatePolicyActionRead")

	getRatePolicyAction := v2.GetRatePolicyActionRequest{}

	configid, err := tools.GetIntValue("config_id", d)
	if err != nil && !errors.Is(err, tools.ErrNotFound) {
		return diag.FromErr(err)
	}
	getRatePolicyAction.ConfigID = configid

	version, err := tools.GetIntValue("version", d)
	if err != nil && !errors.Is(err, tools.ErrNotFound) {
		return diag.FromErr(err)
	}
	getRatePolicyAction.Version = version

	policyid, err := tools.GetStringValue("policy_id", d)
	if err != nil && !errors.Is(err, tools.ErrNotFound) {
		return diag.FromErr(err)
	}
	getRatePolicyAction.PolicyID = policyid

	ratepolicyid, err := tools.GetIntValue("rate_policy_id", d)
	if err != nil && !errors.Is(err, tools.ErrNotFound) {
		return diag.FromErr(err)
	}
	getRatePolicyAction.ID = ratepolicyid

	ratepolicyaction, err := client.GetRatePolicyAction(ctx, getRatePolicyAction)
	if err != nil {
		logger.Warnf("calling 'getRatePolicyAction': %s", err.Error())
	}

	for _, configval := range ratepolicyaction.RatePolicyActions {
		d.SetId(strconv.Itoa(configval.ID))
	}

	return nil
}

func resourceRatePolicyActionDelete(ctx context.Context, d *schema.ResourceData, m interface{}) diag.Diagnostics {
	meta := akamai.Meta(m)
	client := inst.Client(meta)
	logger := meta.Log("APPSEC", "resourceRatePolicyActionRemove")

	updateRatePolicyAction := v2.UpdateRatePolicyActionRequest{}

	configid, err := tools.GetIntValue("config_id", d)
	if err != nil && !errors.Is(err, tools.ErrNotFound) {
		return diag.FromErr(err)
	}
	updateRatePolicyAction.ConfigID = configid

	version, err := tools.GetIntValue("version", d)
	if err != nil && !errors.Is(err, tools.ErrNotFound) {
		return diag.FromErr(err)
	}
	updateRatePolicyAction.Version = version

	policyid, err := tools.GetStringValue("policy_id", d)
	if err != nil && !errors.Is(err, tools.ErrNotFound) {
		return diag.FromErr(err)
	}
	updateRatePolicyAction.PolicyID = policyid

	ratepolicyid, err := tools.GetIntValue("rate_policy_id", d)
	if err != nil && !errors.Is(err, tools.ErrNotFound) {
		return diag.FromErr(err)
	}
	updateRatePolicyAction.ID = ratepolicyid

	updateRatePolicyAction.Ipv4Action = "none"
	updateRatePolicyAction.Ipv6Action = "none"

	_, erru := client.UpdateRatePolicyAction(ctx, updateRatePolicyAction)
	if erru != nil {
		logger.Warnf("calling 'removeRatePolicyAction': %s", erru.Error())
	}

	d.SetId("")

	return nil
}

func resourceRatePolicyActionUpdate(ctx context.Context, d *schema.ResourceData, m interface{}) diag.Diagnostics {
	meta := akamai.Meta(m)
	client := inst.Client(meta)
	logger := meta.Log("APPSEC", "resourceRatePolicyActionUpdate")

	updateRatePolicyAction := v2.UpdateRatePolicyActionRequest{}

	configid, err := tools.GetIntValue("config_id", d)
	if err != nil && !errors.Is(err, tools.ErrNotFound) {
		return diag.FromErr(err)
	}
	updateRatePolicyAction.ConfigID = configid

	version, err := tools.GetIntValue("version", d)
	if err != nil && !errors.Is(err, tools.ErrNotFound) {
		return diag.FromErr(err)
	}
	updateRatePolicyAction.Version = version

	policyid, err := tools.GetStringValue("policy_id", d)
	if err != nil && !errors.Is(err, tools.ErrNotFound) {
		return diag.FromErr(err)
	}
	updateRatePolicyAction.PolicyID = policyid

	ratepolicyid, err := tools.GetIntValue("rate_policy_id", d)
	if err != nil && !errors.Is(err, tools.ErrNotFound) {
		return diag.FromErr(err)
	}
	updateRatePolicyAction.ID = ratepolicyid

	ipv4action, err := tools.GetStringValue("ipv4_action", d)
	if err != nil && !errors.Is(err, tools.ErrNotFound) {
		return diag.FromErr(err)
	}
	updateRatePolicyAction.Ipv4Action = ipv4action

	ipv6action, err := tools.GetStringValue("ipv6_action", d)
	if err != nil && !errors.Is(err, tools.ErrNotFound) {
		return diag.FromErr(err)
	}
	updateRatePolicyAction.Ipv6Action = ipv6action

	_, erru := client.UpdateRatePolicyAction(ctx, updateRatePolicyAction)
	if erru != nil {
		logger.Warnf("calling 'updateRatePolicyAction': %s", erru.Error())
	}

	return resourceRatePolicyActionRead(ctx, d, m)
}
