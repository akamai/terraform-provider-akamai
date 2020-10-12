package appsec

import (
	"context"
	"strconv"

	v2 "github.com/akamai/AkamaiOPEN-edgegrid-golang/v2/pkg/appsec"
	"github.com/akamai/terraform-provider-akamai/v2/pkg/akamai"

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

	getRatePolicyAction.ConfigID = d.Get("config_id").(int)
	getRatePolicyAction.Version = d.Get("version").(int)
	getRatePolicyAction.PolicyID = d.Get("policy_id").(string)
	getRatePolicyAction.ID = d.Get("rate_policy_id").(int)

	ratepolicyaction, err := client.GetRatePolicyAction(ctx, getRatePolicyAction)
	if err != nil {
		logger.Warnf("calling 'getRatePolicyAction': %s", err.Error())
		return diag.FromErr(err)
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

	updateRatePolicyAction.ConfigID = d.Get("config_id").(int)
	updateRatePolicyAction.Version = d.Get("version").(int)
	updateRatePolicyAction.PolicyID = d.Get("policy_id").(string)
	updateRatePolicyAction.ID = d.Get("rate_policy_id").(int)
	updateRatePolicyAction.Ipv4Action = "none"
	updateRatePolicyAction.Ipv6Action = "none"

	_, err := client.UpdateRatePolicyAction(ctx, updateRatePolicyAction)
	if err != nil {
		logger.Warnf("calling 'removeRatePolicyAction': %s", err.Error())
		return diag.FromErr(err)
	}

	d.SetId("")

	return nil
}

func resourceRatePolicyActionUpdate(ctx context.Context, d *schema.ResourceData, m interface{}) diag.Diagnostics {
	meta := akamai.Meta(m)
	client := inst.Client(meta)
	logger := meta.Log("APPSEC", "resourceRatePolicyActionUpdate")

	updateRatePolicyAction := v2.UpdateRatePolicyActionRequest{}

	updateRatePolicyAction.ConfigID = d.Get("config_id").(int)
	updateRatePolicyAction.Version = d.Get("version").(int)
	updateRatePolicyAction.PolicyID = d.Get("policy_id").(string)
	updateRatePolicyAction.ID = d.Get("rate_policy_id").(int)
	updateRatePolicyAction.Ipv4Action = d.Get("ipv4_action").(string)
	updateRatePolicyAction.Ipv6Action = d.Get("ipv6_action").(string)

	_, erru := client.UpdateRatePolicyAction(ctx, updateRatePolicyAction)
	if erru != nil {
		logger.Warnf("calling 'updateRatePolicyAction': %s", erru.Error())
		return diag.FromErr(erru)
	}

	return resourceRatePolicyActionRead(ctx, d, m)
}

const (
	Alert = "alert"
	Deny  = "deny"
	None  = "none"
)
