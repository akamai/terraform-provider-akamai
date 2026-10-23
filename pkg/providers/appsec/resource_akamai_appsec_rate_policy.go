package appsec

import (
	"context"
	"encoding/json"
	"errors"
	"fmt"
	"strconv"

	"github.com/akamai/AkamaiOPEN-edgegrid-golang/v14/pkg/appsec"
	"github.com/akamai/terraform-provider-akamai/v11/pkg/common/id"
	"github.com/akamai/terraform-provider-akamai/v11/pkg/common/tf"
	"github.com/akamai/terraform-provider-akamai/v11/pkg/meta"
	"github.com/hashicorp/terraform-plugin-sdk/v2/diag"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/customdiff"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/schema"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/validation"
)

// appsec v1
//
// https://techdocs.akamai.com/application-security/reference/api
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
				Type:        schema.TypeInt,
				Required:    true,
				Description: "Unique identifier of the security configuration",
			},
			"rate_policy": {
				Type:             schema.TypeString,
				Required:         true,
				ValidateDiagFunc: validation.ToDiagFunc(validation.StringIsJSON),
				DiffSuppressFunc: suppressRatePolicyDiffs,
				Description:      "JSON-formatted definition of the rate policy",
			},
			"rate_policy_id": {
				Type:        schema.TypeInt,
				Computed:    true,
				Description: "Unique identifier of the rate policy",
			},
		},
	}
}

func resourceRatePolicyCreate(ctx context.Context, d *schema.ResourceData, m interface{}) diag.Diagnostics {
	meta := meta.Must(m)
	client := meta.Client().GetAPPSEC()
	logger := meta.Log("APPSEC", "resourceRatePolicyCreate")
	logger.Debugf("in resourceRatePolicyCreate")

	configID, err := tf.GetIntValue("config_id", d)
	if err != nil {
		return diag.FromErr(err)
	}
	version, err := getModifiableConfigVersion(ctx, configID, "ratePolicy", client)
	if err != nil {
		return diag.FromErr(err)
	}

	defer invalidateRatePolicyCache(configID, version, logger)

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
		logger.Errorf("calling 'createRatePolicy': %s", err.Error())
		return diag.FromErr(err)
	}

	d.SetId(fmt.Sprintf("%d:%d", createRatePolicy.ConfigID, ratepolicy.ID))

	// Drop the cached list before reading back, otherwise the read would not see the new policy.
	invalidateRatePolicyCache(configID, version, logger)
	return resourceRatePolicyRead(ctx, d, m)
}

func resourceRatePolicyRead(ctx context.Context, d *schema.ResourceData, m interface{}) diag.Diagnostics {
	meta := meta.Must(m)
	client := meta.Client().GetAPPSEC()
	logger := meta.Log("APPSEC", "resourceRatePolicyRead")
	logger.Debugf("in resourceRatePolicyRead")

	iDParts, err := id.Split(d.Id(), 2, "configID:ratePolicyID")
	if err != nil {
		return diag.FromErr(err)
	}

	configID, err := strconv.Atoi(iDParts[0])
	if err != nil {
		return diag.FromErr(err)
	}
	version, err := getLatestConfigVersion(ctx, configID, client)
	if err != nil {
		return diag.FromErr(err)
	}
	ratePolicyID, err := strconv.Atoi(iDParts[1])
	if err != nil {
		return diag.FromErr(err)
	}

	// Use cache helper: fetches all rate policies in one API call on first read,
	// then filters locally; subsequent reads for the same configID+version hit cache.
	ratepolicy, err := getRatePolicy(ctx, configID, version, ratePolicyID, m)
	if errors.Is(err, ErrRatePolicyNotFound) {
		// The list endpoint can lag immediately after a Create or Update. Fall back to
		// the single-item fetch, which always reflects the just-written resource.
		ratepolicy, err = client.GetRatePolicy(ctx, appsec.GetRatePolicyRequest{
			ConfigID:      configID,
			ConfigVersion: version,
			RatePolicyID:  ratePolicyID,
		})
	}
	if err != nil {
		logger.Errorf("calling 'getRatePolicy': %s", err.Error())
		return diag.FromErr(err)
	}

	if err := d.Set("config_id", configID); err != nil {
		return diag.Errorf("%s: %s", tf.ErrValueSet, err.Error())
	}
	jsonBody, err := json.Marshal(ratepolicy)
	if err != nil {
		return diag.FromErr(err)
	}
	if err := d.Set("rate_policy_id", ratePolicyID); err != nil {
		return diag.Errorf("%s: %s", tf.ErrValueSet, err.Error())
	}
	if err := d.Set("rate_policy", string(jsonBody)); err != nil {
		return diag.Errorf("%s: %s", tf.ErrValueSet, err.Error())
	}

	return nil
}

func resourceRatePolicyUpdate(ctx context.Context, d *schema.ResourceData, m interface{}) diag.Diagnostics {
	meta := meta.Must(m)
	client := meta.Client().GetAPPSEC()
	logger := meta.Log("APPSEC", "resourceRatePolicyUpdate")
	logger.Debugf("in resourceRatePolicy`Update")

	iDParts, err := id.Split(d.Id(), 2, "configID:ratePolicyID")
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

	version, err := getModifiableConfigVersion(ctx, configID, "ratePolicy", client)
	if err != nil {
		return diag.FromErr(err)
	}

	defer invalidateRatePolicyCache(configID, version, logger)

	updateRatePolicy := appsec.UpdateRatePolicyRequest{
		ConfigID:       configID,
		ConfigVersion:  version,
		RatePolicyID:   ratePolicyID,
		JsonPayloadRaw: rawJSON,
	}

	if _, err := client.UpdateRatePolicy(ctx, updateRatePolicy); err != nil {
		logger.Errorf("calling 'updateRatePolicy': %s", err.Error())
		return diag.FromErr(err)
	}

	// Drop the cached list before reading back, otherwise the read would return the pre-update policy.
	invalidateRatePolicyCache(configID, version, logger)
	return resourceRatePolicyRead(ctx, d, m)
}

func resourceRatePolicyDelete(ctx context.Context, d *schema.ResourceData, m interface{}) diag.Diagnostics {
	meta := meta.Must(m)
	client := meta.Client().GetAPPSEC()
	logger := meta.Log("APPSEC", "resourceRatePolicyDelete")
	logger.Debugf("in resourceRatePolicyDelete")

	iDParts, err := id.Split(d.Id(), 2, "configID:ratePolicyID")
	if err != nil {
		return diag.FromErr(err)
	}

	configID, err := strconv.Atoi(iDParts[0])
	if err != nil {
		return diag.FromErr(err)
	}
	version, err := getModifiableConfigVersion(ctx, configID, "ratePolicy", client)
	if err != nil {
		return diag.FromErr(err)
	}

	ratePolicyID, err := strconv.Atoi(iDParts[1])
	if err != nil {
		return diag.FromErr(err)
	}

	defer invalidateRatePolicyCache(configID, version, logger)

	deleteRatePolicy := appsec.RemoveRatePolicyRequest{
		ConfigID:      configID,
		ConfigVersion: version,
		RatePolicyID:  ratePolicyID,
	}

	_, err = client.RemoveRatePolicy(ctx, deleteRatePolicy)
	if err != nil {
		logger.Errorf("calling 'removeRatePolicy': %s", err.Error())
		return diag.FromErr(err)
	}
	return nil
}
