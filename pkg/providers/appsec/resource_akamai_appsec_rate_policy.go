package appsec

import (
	"context"
	"encoding/json"
	"errors"
	"fmt"
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
func resourceRatePolicy() *schema.Resource {
	return &schema.Resource{
		CreateContext: resourceRatePolicyCreate,
		ReadContext:   resourceRatePolicyRead,
		UpdateContext: resourceRatePolicyUpdate,
		DeleteContext: resourceRatePolicyDelete,
		Schema: map[string]*schema.Schema{
			"config_id": {
				Type:     schema.TypeInt,
				Required: true,
			},
			"version": {
				Type:     schema.TypeInt,
				Required: true,
			},
			"rate_policy": {
				Type:         schema.TypeString,
				Required:     true,
				ValidateFunc: validation.StringIsJSON,
			},
			"rate_policy_id": {
				Type:     schema.TypeInt,
				Computed: true,
			},
		},
	}
}

func resourceRatePolicyCreate(ctx context.Context, d *schema.ResourceData, m interface{}) diag.Diagnostics {
	meta := akamai.Meta(m)
	client := inst.Client(meta)
	logger := meta.Log("APPSEC", "resourceRatePolicyCreate")

	createRatePolicy := appsec.CreateRatePolicyRequest{}

	jsonpostpayload := d.Get("rate_policy")

	if err := json.Unmarshal([]byte(jsonpostpayload.(string)), &createRatePolicy); err != nil {
		return diag.FromErr(err)
	}

	configid, err := tools.GetIntValue("config_id", d)
	if err != nil && !errors.Is(err, tools.ErrNotFound) {
		return diag.FromErr(err)
	}
	createRatePolicy.ConfigID = configid

	version, err := tools.GetIntValue("version", d)
	if err != nil && !errors.Is(err, tools.ErrNotFound) {
		return diag.FromErr(err)
	}
	createRatePolicy.ConfigVersion = version

	ratepolicy, err := client.CreateRatePolicy(ctx, createRatePolicy)
	if err != nil {
		logger.Warnf("calling 'createRatePolicyAction': %s", err.Error())
		return diag.FromErr(err)
	}

	d.SetId(strconv.Itoa(ratepolicy.ID))

	return resourceRatePolicyRead(ctx, d, meta)
}

func resourceRatePolicyUpdate(ctx context.Context, d *schema.ResourceData, m interface{}) diag.Diagnostics {
	meta := akamai.Meta(m)
	client := inst.Client(meta)
	logger := meta.Log("APPSEC", "resourceRatePolicyUpdate")

	updateRatePolicy := appsec.UpdateRatePolicyRequest{}

	jsonpostpayload := d.Get("rate_policy")

	if err := json.Unmarshal([]byte(jsonpostpayload.(string)), &updateRatePolicy); err != nil {
		return diag.FromErr(err)
	}

	configid, err := tools.GetIntValue("config_id", d)
	if err != nil && !errors.Is(err, tools.ErrNotFound) {
		return diag.FromErr(err)
	}
	updateRatePolicy.ConfigID = configid

	version, err := tools.GetIntValue("version", d)
	if err != nil && !errors.Is(err, tools.ErrNotFound) {
		return diag.FromErr(err)
	}
	updateRatePolicy.ConfigVersion = version

	ratePolicyID, errconv := strconv.Atoi(d.Id())

	if errconv != nil {
		return diag.FromErr(errconv)
	}
	updateRatePolicy.RatePolicyID = ratePolicyID

	_, erru := client.UpdateRatePolicy(ctx, updateRatePolicy)
	if erru != nil {
		logger.Warnf("calling 'updateRatePolicyAction': %s", erru.Error())
		return diag.FromErr(erru)
	}

	return resourceRatePolicyRead(ctx, d, meta)
}

func resourceRatePolicyDelete(ctx context.Context, d *schema.ResourceData, m interface{}) diag.Diagnostics {
	meta := akamai.Meta(m)
	client := inst.Client(meta)
	logger := meta.Log("APPSEC", "resourceRatePolicyDelete")

	deleteRatePolicy := appsec.RemoveRatePolicyRequest{}

	configid, err := tools.GetIntValue("config_id", d)
	if err != nil && !errors.Is(err, tools.ErrNotFound) {
		return diag.FromErr(err)
	}
	deleteRatePolicy.ConfigID = configid

	version, err := tools.GetIntValue("version", d)
	if err != nil && !errors.Is(err, tools.ErrNotFound) {
		return diag.FromErr(err)
	}
	deleteRatePolicy.ConfigVersion = version

	ratePolicyID, errconv := strconv.Atoi(d.Id())

	if errconv != nil {
		return diag.FromErr(errconv)
	}
	deleteRatePolicy.RatePolicyID = ratePolicyID

	_, errd := client.RemoveRatePolicy(ctx, deleteRatePolicy)
	if errd != nil {
		logger.Warnf("calling 'removeRatePolicyAction': %s", errd.Error())
		return diag.FromErr(errd)
	}

	d.SetId("")

	return nil
}

func resourceRatePolicyRead(ctx context.Context, d *schema.ResourceData, m interface{}) diag.Diagnostics {
	meta := akamai.Meta(m)
	client := inst.Client(meta)
	logger := meta.Log("APPSEC", "resourceRatePolicyRead")

	getRatePolicy := appsec.GetRatePolicyRequest{}

	configid, err := tools.GetIntValue("config_id", d)
	if err != nil && !errors.Is(err, tools.ErrNotFound) {
		return diag.FromErr(err)
	}
	getRatePolicy.ConfigID = configid

	version, err := tools.GetIntValue("version", d)
	if err != nil && !errors.Is(err, tools.ErrNotFound) {
		return diag.FromErr(err)
	}
	getRatePolicy.ConfigVersion = version

	ratePolicyID, errconv := strconv.Atoi(d.Id())

	if errconv != nil {
		return diag.FromErr(errconv)
	}
	getRatePolicy.RatePolicyID = ratePolicyID

	ratepolicy, errd := client.GetRatePolicy(ctx, getRatePolicy)
	if errd != nil {
		logger.Warnf("calling 'getRatePolicyAction': %s", errd.Error())
		return diag.FromErr(errd)
	}

	if err := d.Set("rate_policy_id", ratepolicy.ID); err != nil {
		return diag.FromErr(fmt.Errorf("%w: %s", tools.ErrValueSet, err.Error()))
	}

	d.SetId(strconv.Itoa(ratepolicy.ID))

	return nil
}
