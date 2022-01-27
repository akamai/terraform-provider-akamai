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
			VerifyIDUnchanged,
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
				ValidateDiagFunc: validation.ToDiagFunc(validation.StringIsJSON),
				DiffSuppressFunc: suppressEquivalentJSONDiffsGeneric,
			},
		},
	}
}

func resourceRatePolicyCreate(ctx context.Context, d *schema.ResourceData, m interface{}) diag.Diagnostics {
	meta := akamai.Meta(m)
	client := inst.Client(meta)
	logger := meta.Log("APPSEC", "resourceRatePolicyCreate")
	logger.Debugf("in resourceRatePolicyCreate")

	configID, err := tools.GetIntValue("config_id", d)
	if err != nil {
		return diag.FromErr(err)
	}
	version := getModifiableConfigVersion(ctx, configID, "ratePolicy", m)
	jsonpostpayload := d.Get("rate_policy")
	jsonPayloadRaw := []byte(jsonpostpayload.(string))
	rawJSON := (json.RawMessage)(jsonPayloadRaw)

	createRatePolicy := appsec.CreateRatePolicyRequest{
		ConfigID:       configID,
		ConfigVersion:  version,
		JsonPayloadRaw: rawJSON,
	}

	ratepolicy, err := client.CreateRatePolicy(ctx, createRatePolicy)
	if err != nil {
		logger.Warnf("calling 'createRatePolicy': %s", err.Error())
		return diag.FromErr(err)
	}

	d.SetId(fmt.Sprintf("%d:%d", createRatePolicy.ConfigID, ratepolicy.ID))

	return resourceRatePolicyRead(ctx, d, meta)
}

func resourceRatePolicyRead(ctx context.Context, d *schema.ResourceData, m interface{}) diag.Diagnostics {
	meta := akamai.Meta(m)
	client := inst.Client(meta)
	logger := meta.Log("APPSEC", "resourceRatePolicyRead")
	logger.Debugf("in resourceRatePolicyRead")

	iDParts, err := splitID(d.Id(), 2, "configID:ratePolicyID")
	if err != nil {
		return diag.FromErr(err)
	}

	configID, err := strconv.Atoi(iDParts[0])
	if err != nil {
		return diag.FromErr(err)
	}
	version := getLatestConfigVersion(ctx, configID, m)
	ratePolicyID, err := strconv.Atoi(iDParts[1])
	if err != nil {
		return diag.FromErr(err)
	}

	readRatePolicy := appsec.GetRatePolicyRequest{
		ConfigID:      configID,
		ConfigVersion: version,
		RatePolicyID:  ratePolicyID,
	}

	ratepolicy, err := client.GetRatePolicy(ctx, readRatePolicy)
	if err != nil {
		logger.Warnf("calling 'getRatePolicy': %s", err.Error())
		return diag.FromErr(err)
	}

	if err := d.Set("config_id", configID); err != nil {
		return diag.Errorf("%s: %s", tools.ErrValueSet, err.Error())
	}
	jsonBody, err := json.Marshal(ratepolicy)
	if err != nil {
		return diag.FromErr(err)
	}
	if err := d.Set("rate_policy_id", ratePolicyID); err != nil {
		return diag.Errorf("%s: %s", tools.ErrValueSet, err.Error())
	}
	if err := d.Set("rate_policy", string(jsonBody)); err != nil {
		return diag.Errorf("%s: %s", tools.ErrValueSet, err.Error())
	}

	return nil
}

func resourceRatePolicyUpdate(ctx context.Context, d *schema.ResourceData, m interface{}) diag.Diagnostics {
	meta := akamai.Meta(m)
	client := inst.Client(meta)
	logger := meta.Log("APPSEC", "resourceRatePolicyUpdate")
	logger.Debugf("in resourceRatePolicy`Update")

	iDParts, err := splitID(d.Id(), 2, "configID:ratePolicyID")
	if err != nil {
		return diag.FromErr(err)
	}

	configID, err := strconv.Atoi(iDParts[0])
	if err != nil {
		return diag.FromErr(err)
	}
	ratePolicyID, err := strconv.Atoi(iDParts[1])
	if err != nil {
		return diag.FromErr(err)
	}

	jsonpostpayload := d.Get("rate_policy")
	jsonPayloadRaw := []byte(jsonpostpayload.(string))
	rawJSON := (json.RawMessage)(jsonPayloadRaw)

	updateRatePolicy := appsec.UpdateRatePolicyRequest{
		ConfigID:       configID,
		ConfigVersion:  getModifiableConfigVersion(ctx, configID, "ratePolicy", m),
		RatePolicyID:   ratePolicyID,
		JsonPayloadRaw: rawJSON,
	}

	_, err = client.UpdateRatePolicy(ctx, updateRatePolicy)
	if err != nil {
		logger.Warnf("calling 'updateRatePolicy': %s", err.Error())
		return diag.FromErr(err)
	}

	return resourceRatePolicyRead(ctx, d, meta)
}

func resourceRatePolicyDelete(ctx context.Context, d *schema.ResourceData, m interface{}) diag.Diagnostics {
	meta := akamai.Meta(m)
	client := inst.Client(meta)
	logger := meta.Log("APPSEC", "resourceRatePolicyDelete")
	logger.Debugf("in resourceRatePolicyDelete")

	iDParts, err := splitID(d.Id(), 2, "configID:ratePolicyID")
	if err != nil {
		return diag.FromErr(err)
	}

	configID, err := strconv.Atoi(iDParts[0])
	if err != nil {
		return diag.FromErr(err)
	}
	version := getModifiableConfigVersion(ctx, configID, "ratePolicy", m)
	ratePolicyID, err := strconv.Atoi(iDParts[1])
	if err != nil {
		return diag.FromErr(err)
	}

	deleteRatePolicy := appsec.RemoveRatePolicyRequest{
		ConfigID:      configID,
		ConfigVersion: version,
		RatePolicyID:  ratePolicyID,
	}

	_, err = client.RemoveRatePolicy(ctx, deleteRatePolicy)
	if err != nil {
		logger.Warnf("calling 'removeRatePolicy': %s", err.Error())
		return diag.FromErr(err)
	}

	d.SetId("")

	return nil
}
