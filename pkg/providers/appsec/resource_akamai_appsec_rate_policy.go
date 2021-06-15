package appsec

import (
	"context"
	"encoding/json"
	"fmt"
	"strconv"

	"github.com/akamai/AkamaiOPEN-edgegrid-golang/v2/pkg/appsec"
	"github.com/akamai/terraform-provider-akamai/v2/pkg/akamai"
	"github.com/akamai/terraform-provider-akamai/v2/pkg/tools"

	"github.com/hashicorp/terraform-plugin-sdk/v2/diag"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/customdiff"
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
		Importer: &schema.ResourceImporter{
			StateContext: schema.ImportStatePassthroughContext,
		},
		CustomizeDiff: customdiff.All(
			VerifyIdUnchanged,
		),
		Schema: map[string]*schema.Schema{
			"config_id": {
				Type:     schema.TypeInt,
				Required: true,
			},
			"rate_policy_id": {
				Type:     schema.TypeInt,
				Computed: true,
			},
			"rate_policy": {
				Type:             schema.TypeString,
				Required:         true,
				ValidateFunc:     validation.StringIsJSON,
				DiffSuppressFunc: suppressEquivalentJsonDiffsGeneric,
			},
		},
	}
}

func resourceRatePolicyCreate(ctx context.Context, d *schema.ResourceData, m interface{}) diag.Diagnostics {
	meta := akamai.Meta(m)
	client := inst.Client(meta)
	logger := meta.Log("APPSEC", "resourceRatePolicyCreate")
	logger.Debugf("!!! in resourceRatePolicyCreate")

	configid, err := tools.GetIntValue("config_id", d)
	if err != nil {
		return diag.FromErr(err)
	}
	version := getModifiableConfigVersion(ctx, configid, "ratePolicy", m)
	jsonpostpayload := d.Get("rate_policy")
	jsonPayloadRaw := []byte(jsonpostpayload.(string))
	rawJSON := (json.RawMessage)(jsonPayloadRaw)

	createRatePolicy := appsec.CreateRatePolicyRequest{}
	createRatePolicy.ConfigID = configid
	createRatePolicy.ConfigVersion = version
	createRatePolicy.JsonPayloadRaw = rawJSON

	ratepolicy, err := client.CreateRatePolicy(ctx, createRatePolicy)
	if err != nil {
		logger.Warnf("calling 'createRatePolicyAction': %s", err.Error())
		return diag.FromErr(err)
	}

	d.SetId(fmt.Sprintf("%d:%d", createRatePolicy.ConfigID, ratepolicy.ID))

	return resourceRatePolicyRead(ctx, d, meta)
}

func resourceRatePolicyRead(ctx context.Context, d *schema.ResourceData, m interface{}) diag.Diagnostics {
	meta := akamai.Meta(m)
	client := inst.Client(meta)
	logger := meta.Log("APPSEC", "resourceRatePolicyRead")
	logger.Debugf("!!! in resourceRatePolicyRead")

	idParts, err := splitID(d.Id(), 2, "configid:ratepolicyid")
	if err != nil {
		return diag.FromErr(err)
	}

	configid, err := strconv.Atoi(idParts[0])
	if err != nil {
		return diag.FromErr(err)
	}
	version := getLatestConfigVersion(ctx, configid, m)
	ratePolicyID, err := strconv.Atoi(idParts[1])
	if err != nil {
		return diag.FromErr(err)
	}

	readRatePolicy := appsec.GetRatePolicyRequest{}
	readRatePolicy.ConfigID = configid
	readRatePolicy.ConfigVersion = version
	readRatePolicy.RatePolicyID = ratePolicyID

	ratepolicy, err := client.GetRatePolicy(ctx, readRatePolicy)
	if err != nil {
		logger.Warnf("calling 'getRatePolicyAction': %s", err.Error())
		return diag.FromErr(err)
	}

	if err := d.Set("config_id", configid); err != nil {
		return diag.FromErr(fmt.Errorf("%w: %s", tools.ErrValueSet, err.Error()))
	}
	jsonBody, err := json.Marshal(ratepolicy)
	if err != nil {
		return diag.FromErr(err)
	}
	if err := d.Set("rate_policy_id", ratePolicyID); err != nil {
		return diag.FromErr(fmt.Errorf("%w: %s", tools.ErrValueSet, err.Error()))
	}
	if err := d.Set("rate_policy", string(jsonBody)); err != nil {
		return diag.FromErr(fmt.Errorf("%w: %s", tools.ErrValueSet, err.Error()))
	}

	return nil
}

func resourceRatePolicyUpdate(ctx context.Context, d *schema.ResourceData, m interface{}) diag.Diagnostics {
	meta := akamai.Meta(m)
	client := inst.Client(meta)
	logger := meta.Log("APPSEC", "resourceRatePolicyUpdate")
	logger.Debugf("!!! in resourceRatePolicy`Update")

	idParts, err := splitID(d.Id(), 2, "configid:ratepolicyid")
	if err != nil {
		return diag.FromErr(err)
	}

	configid, err := strconv.Atoi(idParts[0])
	if err != nil {
		return diag.FromErr(err)
	}
	ratePolicyID, err := strconv.Atoi(idParts[1])
	if err != nil {
		return diag.FromErr(err)
	}

	jsonpostpayload := d.Get("rate_policy")
	jsonPayloadRaw := []byte(jsonpostpayload.(string))
	rawJSON := (json.RawMessage)(jsonPayloadRaw)

	updateRatePolicy := appsec.UpdateRatePolicyRequest{}
	updateRatePolicy.ConfigID = configid
	updateRatePolicy.ConfigVersion = getModifiableConfigVersion(ctx, configid, "ratePolicy", m)
	updateRatePolicy.RatePolicyID = ratePolicyID
	updateRatePolicy.JsonPayloadRaw = rawJSON

	_, err = client.UpdateRatePolicy(ctx, updateRatePolicy)
	if err != nil {
		logger.Warnf("calling 'updateRatePolicyAction': %s", err.Error())
		return diag.FromErr(err)
	}

	return resourceRatePolicyRead(ctx, d, meta)
}

func resourceRatePolicyDelete(ctx context.Context, d *schema.ResourceData, m interface{}) diag.Diagnostics {
	meta := akamai.Meta(m)
	client := inst.Client(meta)
	logger := meta.Log("APPSEC", "resourceRatePolicyDelete")
	logger.Debugf("!!! in resourceRatePolicyDelete")

	idParts, err := splitID(d.Id(), 2, "configid:ratepolicyid")
	if err != nil {
		return diag.FromErr(err)
	}

	configid, err := strconv.Atoi(idParts[0])
	if err != nil {
		return diag.FromErr(err)
	}
	version := getModifiableConfigVersion(ctx, configid, "ratePolicy", m)
	ratePolicyID, err := strconv.Atoi(idParts[1])
	if err != nil {
		return diag.FromErr(err)
	}

	deleteRatePolicy := appsec.RemoveRatePolicyRequest{}
	deleteRatePolicy.ConfigID = configid
	deleteRatePolicy.ConfigVersion = version
	deleteRatePolicy.RatePolicyID = ratePolicyID

	_, err = client.RemoveRatePolicy(ctx, deleteRatePolicy)
	if err != nil {
		logger.Warnf("calling 'removeRatePolicyAction': %s", err.Error())
		return diag.FromErr(err)
	}

	d.SetId("")

	return nil
}
