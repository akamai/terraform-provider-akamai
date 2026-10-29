package appsec

import (
	"context"
	"encoding/json"
	"fmt"
	"strconv"

	"github.com/akamai/AkamaiOPEN-edgegrid-golang/v15/pkg/appsec"
	"github.com/akamai/terraform-provider-akamai/v12/pkg/common/id"
	"github.com/akamai/terraform-provider-akamai/v12/pkg/common/tf"
	"github.com/akamai/terraform-provider-akamai/v12/pkg/meta"
	"github.com/hashicorp/terraform-plugin-sdk/v2/diag"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/customdiff"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/schema"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/validation"
)

// appsec v1
//
// https://techdocs.akamai.com/application-security/reference/api
func resourceReputationProfile() *schema.Resource {
	return &schema.Resource{
		CreateContext: resourceReputationProfileCreate,
		ReadContext:   resourceReputationProfileRead,
		UpdateContext: resourceReputationProfileUpdate,
		DeleteContext: resourceReputationProfileDelete,
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
			"reputation_profile": {
				Type:             schema.TypeString,
				Required:         true,
				ValidateDiagFunc: validation.ToDiagFunc(validation.StringIsJSON),
				DiffSuppressFunc: suppressEquivalentReputationProfileDiffs,
				Description:      "JSON-formatted definition of the reputation profile",
			},
			"reputation_profile_id": {
				Type:        schema.TypeInt,
				Computed:    true,
				Description: "Unique identifier of the reputation profile",
			},
		},
	}
}

func resourceReputationProfileCreate(ctx context.Context, d *schema.ResourceData, m interface{}) diag.Diagnostics {
	meta := meta.Must(m)
	client := meta.Client().GetAPPSEC()
	logger := meta.Log("APPSEC", "resourceReputationProfileCreate")
	logger.Debug("in resourceReputationProfileCreate")

	configID, err := tf.GetIntValueAsInt64("config_id", d)
	if err != nil {
		return diag.FromErr(err)
	}
	version, err := getModifiableConfigVersion(ctx, configID, "reputationProfile", client)
	if err != nil {
		return diag.FromErr(err)
	}

	defer invalidateReputationProfileCache(configID, version, logger)

	jsonpostpayload, err := tf.GetStringValue("reputation_profile", d)
	if err != nil {
		return diag.FromErr(err)
	}
	jsonPayloadRaw := []byte(jsonpostpayload)
	rawJSON := (json.RawMessage)(jsonPayloadRaw)

	createReputationProfile := appsec.CreateReputationProfileRequest{
		ConfigID:       configID,
		ConfigVersion:  version,
		JsonPayloadRaw: rawJSON,
	}

	response, err := client.CreateReputationProfile(ctx, createReputationProfile)
	if err != nil {
		logger.Errorf("calling 'CreateReputationProfile': %s", err.Error())
		return diag.FromErr(err)
	}

	d.SetId(fmt.Sprintf("%d:%d", createReputationProfile.ConfigID, response.ID))

	// Drop the cached list before reading back, otherwise the read would not see the new profile.
	invalidateReputationProfileCache(configID, version, logger)
	return resourceReputationProfileRead(ctx, d, m)
}

func resourceReputationProfileRead(ctx context.Context, d *schema.ResourceData, m interface{}) diag.Diagnostics {
	meta := meta.Must(m)
	client := meta.Client().GetAPPSEC()
	logger := meta.Log("APPSEC", "resourceReputationProfileRead")
	logger.Debug("in resourceReputationProfileRead")

	iDParts, err := id.Split(d.Id(), 2, "configID:reputationProfileID")
	if err != nil {
		return diag.FromErr(err)
	}

	configID, err := strconv.ParseInt(iDParts[0], 10, 64)
	if err != nil {
		return diag.FromErr(err)
	}
	reputationProfileID, err := strconv.ParseInt(iDParts[1], 10, 64)
	if err != nil {
		return diag.FromErr(err)
	}

	version, err := getLatestConfigVersion(ctx, configID, client)
	if err != nil {
		return diag.FromErr(err)
	}

	// Use cache helper: fetches all reputation profiles in one API call on first read,
	// then filters locally; subsequent reads for the same configID+version hit cache.
	reputationProfileResponse, err := getReputationProfile(ctx, configID, version, reputationProfileID, m)
	if err != nil {
		logger.Errorf("calling 'getReputationProfile': %s", err.Error())
		return diag.FromErr(err)
	}

	if err := d.Set("config_id", configID); err != nil {
		return diag.Errorf("%s: %s", tf.ErrValueSet, err.Error())
	}
	if err := d.Set("reputation_profile_id", reputationProfileID); err != nil {
		return diag.Errorf("%s: %s", tf.ErrValueSet, err.Error())
	}
	normalized, err := marshalNormalize(reputationProfileResponse)
	if err != nil {
		return diag.FromErr(err)
	}
	if err := d.Set("reputation_profile", normalized); err != nil {
		return diag.Errorf("%s: %s", tf.ErrValueSet, err.Error())
	}

	return nil
}

func resourceReputationProfileUpdate(ctx context.Context, d *schema.ResourceData, m interface{}) diag.Diagnostics {
	meta := meta.Must(m)
	client := meta.Client().GetAPPSEC()
	logger := meta.Log("APPSEC", "resourceReputationProfileUpdate")
	logger.Debug("in resourceReputationProfileUpdate")

	iDParts, err := id.Split(d.Id(), 2, "configID:reputationProfileID")
	if err != nil {
		return diag.FromErr(err)
	}

	configID, err := strconv.ParseInt(iDParts[0], 10, 64)
	if err != nil {
		return diag.FromErr(err)
	}
	reputationProfileID, err := strconv.ParseInt(iDParts[1], 10, 64)
	if err != nil {
		return diag.FromErr(err)
	}
	jsonpostpayload, err := tf.GetStringValue("reputation_profile", d)
	if err != nil {
		return diag.FromErr(err)
	}
	version, err := getModifiableConfigVersion(ctx, configID, "reputationProfile", client)
	if err != nil {
		return diag.FromErr(err)
	}

	jsonPayloadRaw := []byte(jsonpostpayload)
	rawJSON := (json.RawMessage)(jsonPayloadRaw)

	updateReputationProfile := appsec.UpdateReputationProfileRequest{
		ConfigID:            configID,
		ConfigVersion:       version,
		ReputationProfileID: reputationProfileID,
		JsonPayloadRaw:      rawJSON,
	}

	// Invalidate before the Read below, otherwise it would be served the stale cached list.
	invalidateReputationProfileCache(configID, version, logger)
	_, err = client.UpdateReputationProfile(ctx, updateReputationProfile)
	if err != nil {
		logger.Errorf("calling 'updateReputationProfile': %s", err.Error())
		return diag.FromErr(err)
	}

	return resourceReputationProfileRead(ctx, d, m)
}

func resourceReputationProfileDelete(ctx context.Context, d *schema.ResourceData, m interface{}) diag.Diagnostics {
	meta := meta.Must(m)
	client := meta.Client().GetAPPSEC()
	logger := meta.Log("APPSEC", "resourceReputationProfileDelete")
	logger.Debug("in resourceReputationProfileDelete")

	iDParts, err := id.Split(d.Id(), 2, "configID:reputationProfileID")
	if err != nil {
		return diag.FromErr(err)
	}

	configID, err := strconv.ParseInt(iDParts[0], 10, 64)
	if err != nil {
		return diag.FromErr(err)
	}
	version, err := getModifiableConfigVersion(ctx, configID, "reputationProfile", client)
	if err != nil {
		return diag.FromErr(err)
	}

	reputationProfileID, err := strconv.ParseInt(iDParts[1], 10, 64)
	if err != nil {
		return diag.FromErr(err)
	}

	defer invalidateReputationProfileCache(configID, version, logger)

	deleteReputationProfile := appsec.RemoveReputationProfileRequest{
		ConfigID:            configID,
		ConfigVersion:       version,
		ReputationProfileID: reputationProfileID,
	}

	_, err = client.RemoveReputationProfile(ctx, deleteReputationProfile)
	if err != nil {
		logger.Errorf("calling 'removeReputationProfile': %s", err.Error())
		return diag.FromErr(err)
	}
	return nil
}
