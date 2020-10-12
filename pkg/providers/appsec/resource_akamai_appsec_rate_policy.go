package appsec

import (
	"context"
	"encoding/json"
	"errors"
	"strconv"

	edge "github.com/akamai/AkamaiOPEN-edgegrid-golang/edgegrid"
	v2 "github.com/akamai/AkamaiOPEN-edgegrid-golang/v2/pkg/appsec"
	"github.com/akamai/terraform-provider-akamai/v2/pkg/akamai"
	"github.com/akamai/terraform-provider-akamai/v2/pkg/tools"

	"github.com/hashicorp/terraform-plugin-sdk/v2/diag"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/schema"
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
			"json": {
				Type:     schema.TypeString,
				Required: true,
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

	createRatePolicy := v2.CreateRatePolicyRequest{}

	jsonpostpayload := d.Get("json")
	json.Unmarshal([]byte(jsonpostpayload.(string)), &createRatePolicy)

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
	CorrelationID := "[APPSEC][resourceRatePolicyUpdate-" + tools.CreateNonce() + "]"
	edge.PrintfCorrelation("[DEBUG]", CorrelationID, " Updating RatePolicy")

	updateRatePolicy := v2.UpdateRatePolicyRequest{}

	jsonpostpayload := d.Get("json")
	json.Unmarshal([]byte(jsonpostpayload.(string)), &updateRatePolicy)

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

	updateRatePolicy.RatePolicyID, _ = strconv.Atoi(d.Id())

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
	CorrelationID := "[APPSEC][resourceRatePolicyDelete-" + tools.CreateNonce() + "]"
	edge.PrintfCorrelation("[DEBUG]", CorrelationID, "  Deleting RatePolicy")

	deleteRatePolicy := v2.RemoveRatePolicyRequest{}

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

	deleteRatePolicy.RatePolicyID, _ = strconv.Atoi(d.Id())

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

	getRatePolicy := v2.GetRatePolicyRequest{}

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

	getRatePolicy.RatePolicyID, _ = strconv.Atoi(d.Id())

	ratepolicy, errd := client.GetRatePolicy(ctx, getRatePolicy)
	if errd != nil {
		logger.Warnf("calling 'getRatePolicyAction': %s", errd.Error())
		return diag.FromErr(errd)
	}

	d.SetId(strconv.Itoa(ratepolicy.ID))

	return nil
}
