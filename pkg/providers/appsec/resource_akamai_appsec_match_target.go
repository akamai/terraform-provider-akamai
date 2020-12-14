package appsec

import (
	"context"
	"encoding/json"
	"errors"
	"strconv"

	"github.com/akamai/AkamaiOPEN-edgegrid-golang/v2/pkg/appsec"
	"github.com/akamai/terraform-provider-akamai/v2/pkg/akamai"
	"github.com/akamai/terraform-provider-akamai/v2/pkg/tools"
	"github.com/hashicorp/terraform-plugin-sdk/v2/diag"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/schema"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/validation"
)

// appsec v1
//
// https://developer.akamai.com/api/cloud_security/application_security/v1.html
func resourceMatchTarget() *schema.Resource {
	return &schema.Resource{
		CreateContext: resourceMatchTargetCreate,
		ReadContext:   resourceMatchTargetRead,
		UpdateContext: resourceMatchTargetUpdate,
		DeleteContext: resourceMatchTargetDelete,
		Schema: map[string]*schema.Schema{
			"config_id": {
				Type:     schema.TypeInt,
				Required: true,
			},
			"version": {
				Type:     schema.TypeInt,
				Required: true,
			},
			"match_target": {
				Type:             schema.TypeString,
				Required:         true,
				ValidateFunc:     validation.StringIsJSON,
				DiffSuppressFunc: suppressEquivalentJSONDiffs,
			},
			"match_target_id": {
				Type:     schema.TypeInt,
				Computed: true,
			},
		},
	}
}

func resourceMatchTargetCreate(ctx context.Context, d *schema.ResourceData, m interface{}) diag.Diagnostics {
	meta := akamai.Meta(m)
	client := inst.Client(meta)
	logger := meta.Log("APPSEC", "resourceMatchTargetCreate")

	createMatchTarget := appsec.CreateMatchTargetRequest{}

	jsonpostpayload := d.Get("match_target")

	json.Unmarshal([]byte(jsonpostpayload.(string)), &createMatchTarget)

	postresp, err := client.CreateMatchTarget(ctx, createMatchTarget)
	if err != nil {
		logger.Errorf("calling 'createMatchTarget': %s", err.Error())
		return diag.FromErr(err)
	}

	jsonBody, err := json.Marshal(postresp)
	if err != nil {
		return diag.FromErr(err)
	}
	d.Set("match_target", string(jsonBody))

	d.Set("match_target_id", postresp.TargetID)

	d.SetId(strconv.Itoa(postresp.TargetID))

	return resourceMatchTargetRead(ctx, d, m)
}

func resourceMatchTargetUpdate(ctx context.Context, d *schema.ResourceData, m interface{}) diag.Diagnostics {
	meta := akamai.Meta(m)
	client := inst.Client(meta)
	logger := meta.Log("APPSEC", "resourceMatchTargetUpdate")

	updateMatchTarget := appsec.UpdateMatchTargetRequest{}

	jsonpostpayload := d.Get("match_target")

	json.Unmarshal([]byte(jsonpostpayload.(string)), &updateMatchTarget)

	targetID, errconv := strconv.Atoi(d.Id())

	if errconv != nil {
		return diag.FromErr(errconv)
	}
	updateMatchTarget.TargetID = targetID

	jsonBody, err := json.Marshal(updateMatchTarget)
	if err != nil {
		return diag.FromErr(err)
	}
	d.Set("match_target", string(jsonBody))

	resp, err := client.UpdateMatchTarget(ctx, updateMatchTarget)
	if err != nil {
		logger.Errorf("calling 'updateMatchTarget': %s", err.Error())
		return diag.FromErr(err)
	}
	jsonBody, err = json.Marshal(resp)
	if err != nil {
		return diag.FromErr(err)
	}
	d.Set("match_target", string(jsonBody))
	return resourceMatchTargetRead(ctx, d, m)
}

func resourceMatchTargetDelete(ctx context.Context, d *schema.ResourceData, m interface{}) diag.Diagnostics {
	meta := akamai.Meta(m)
	client := inst.Client(meta)
	logger := meta.Log("APPSEC", "resourceMatchTargetRemove")

	removeMatchTarget := appsec.RemoveMatchTargetRequest{}

	configid, err := tools.GetIntValue("config_id", d)
	if err != nil && !errors.Is(err, tools.ErrNotFound) {
		return diag.FromErr(err)
	}
	removeMatchTarget.ConfigID = configid

	version, err := tools.GetIntValue("version", d)
	if err != nil && !errors.Is(err, tools.ErrNotFound) {
		return diag.FromErr(err)
	}
	removeMatchTarget.ConfigVersion = version

	targetID, errconv := strconv.Atoi(d.Id())

	if errconv != nil {
		return diag.FromErr(errconv)
	}
	removeMatchTarget.TargetID = targetID

	_, errd := client.RemoveMatchTarget(ctx, removeMatchTarget)
	if errd != nil {
		logger.Errorf("calling 'removeMatchTarget': %s", errd.Error())
		return diag.FromErr(errd)
	}

	d.SetId("")

	return nil
}

func resourceMatchTargetRead(ctx context.Context, d *schema.ResourceData, m interface{}) diag.Diagnostics {
	meta := akamai.Meta(m)
	client := inst.Client(meta)
	logger := meta.Log("APPSEC", "resourceMatchTargetRead")

	getMatchTarget := appsec.GetMatchTargetRequest{}

	configid, err := tools.GetIntValue("config_id", d)
	if err != nil && !errors.Is(err, tools.ErrNotFound) {
		return diag.FromErr(err)
	}
	getMatchTarget.ConfigID = configid

	version, err := tools.GetIntValue("version", d)
	if err != nil && !errors.Is(err, tools.ErrNotFound) {
		return diag.FromErr(err)
	}
	getMatchTarget.ConfigVersion = version

	targetID, errconv := strconv.Atoi(d.Id())

	if errconv != nil {
		return diag.FromErr(errconv)
	}
	getMatchTarget.TargetID = targetID

	matchtarget, err := client.GetMatchTarget(ctx, getMatchTarget)
	if err != nil {
		logger.Errorf("calling 'getMatchTarget': %s", err.Error())
		return diag.FromErr(err)
	}

	jsonBody, err := json.Marshal(matchtarget)
	if err != nil {
		return diag.FromErr(err)
	}
	d.Set("match_target", string(jsonBody))

	d.Set("match_target_id", matchtarget.TargetID)
	d.SetId(strconv.Itoa(matchtarget.TargetID))

	return nil
}
