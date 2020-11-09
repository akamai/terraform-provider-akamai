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
func resourceWAFAttackGroupAction() *schema.Resource {
	return &schema.Resource{
		CreateContext: resourceWAFAttackGroupActionUpdate,
		ReadContext:   resourceWAFAttackGroupActionRead,
		UpdateContext: resourceWAFAttackGroupActionUpdate,
		DeleteContext: resourceWAFAttackGroupActionDelete,
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
			"group_id": {
				Type:     schema.TypeString,
				Required: true,
			},
			"action": {
				Type:     schema.TypeString,
				Required: true,
				ValidateFunc: validation.StringInSlice([]string{
					Alert,
					Deny,
					None,
				}, false),
			},
			"output_text": {
				Type:        schema.TypeString,
				Computed:    true,
				Description: "Text Export representation",
			},
		},
	}
}

func resourceWAFAttackGroupActionRead(ctx context.Context, d *schema.ResourceData, m interface{}) diag.Diagnostics {
	meta := akamai.Meta(m)
	client := inst.Client(meta)
	logger := meta.Log("APPSEC", "resourceWAFAttackGroupActionRead")

	getWAFAttackGroupAction := v2.GetWAFAttackGroupActionRequest{}

	configid, err := tools.GetIntValue("config_id", d)
	if err != nil && !errors.Is(err, tools.ErrNotFound) {
		return diag.FromErr(err)
	}
	getWAFAttackGroupAction.ConfigID = configid

	version, err := tools.GetIntValue("version", d)
	if err != nil && !errors.Is(err, tools.ErrNotFound) {
		return diag.FromErr(err)
	}
	getWAFAttackGroupAction.Version = version

	policyid, err := tools.GetStringValue("policy_id", d)
	if err != nil && !errors.Is(err, tools.ErrNotFound) {
		return diag.FromErr(err)
	}
	getWAFAttackGroupAction.PolicyID = policyid

	group, err := tools.GetStringValue("group_id", d)
	if err != nil && !errors.Is(err, tools.ErrNotFound) {
		return diag.FromErr(err)
	}
	getWAFAttackGroupAction.Group = group

	wafattackgroupaction, err := client.GetWAFAttackGroupAction(ctx, getWAFAttackGroupAction)
	if err != nil {
		logger.Errorf("calling 'getWAFAttackGroupAction': %s", err.Error())
		return diag.FromErr(err)
	}

	ots := OutputTemplates{}
	InitTemplates(ots)

	outputtext, err := RenderTemplates(ots, "WAFAttackGroupActionDS", wafattackgroupaction)
	if err == nil {
		d.Set("output_text", outputtext)
	}

	d.SetId(strconv.Itoa(getWAFAttackGroupAction.ConfigID))

	return nil
}

func resourceWAFAttackGroupActionDelete(ctx context.Context, d *schema.ResourceData, m interface{}) diag.Diagnostics {

	return schema.NoopContext(nil, d, m)
}

func resourceWAFAttackGroupActionUpdate(ctx context.Context, d *schema.ResourceData, m interface{}) diag.Diagnostics {
	meta := akamai.Meta(m)
	client := inst.Client(meta)
	logger := meta.Log("APPSEC", "resourceWAFAttackGroupActionUpdate")

	updateWAFAttackGroupAction := v2.UpdateWAFAttackGroupActionRequest{}

	configid, err := tools.GetIntValue("config_id", d)
	if err != nil && !errors.Is(err, tools.ErrNotFound) {
		return diag.FromErr(err)
	}
	updateWAFAttackGroupAction.ConfigID = configid

	version, err := tools.GetIntValue("version", d)
	if err != nil && !errors.Is(err, tools.ErrNotFound) {
		return diag.FromErr(err)
	}
	updateWAFAttackGroupAction.Version = version

	policyid, err := tools.GetStringValue("policy_id", d)
	if err != nil && !errors.Is(err, tools.ErrNotFound) {
		return diag.FromErr(err)
	}
	updateWAFAttackGroupAction.PolicyID = policyid

	group, err := tools.GetStringValue("group_id", d)
	if err != nil && !errors.Is(err, tools.ErrNotFound) {
		return diag.FromErr(err)
	}
	updateWAFAttackGroupAction.Group = group

	action, err := tools.GetStringValue("action", d)
	if err != nil && !errors.Is(err, tools.ErrNotFound) {
		return diag.FromErr(err)
	}
	updateWAFAttackGroupAction.Action = action
	//updateWAFAttackGroupAction..TargetID, _ = strconv.Atoi(d.Id())
	_, erru := client.UpdateWAFAttackGroupAction(ctx, updateWAFAttackGroupAction)
	if erru != nil {
		logger.Errorf("calling 'updateWAFAttackGroupAction': %s", erru.Error())
		return diag.FromErr(erru)
	}

	return resourceWAFAttackGroupActionRead(ctx, d, m)
}
