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
func resourceReputationProfile() *schema.Resource {
	return &schema.Resource{
		CreateContext: resourceReputationProfileCreate,
		ReadContext:   resourceReputationProfileRead,
		UpdateContext: resourceReputationProfileUpdate,
		DeleteContext: resourceReputationProfileDelete,
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
			"reputation_profile": {
				Type:         schema.TypeString,
				Required:     true,
				ValidateFunc: validation.StringIsJSON,
			},
			"reputation_profile_id": {
				Type:     schema.TypeInt,
				Computed: true,
			},
		},
	}
}

func resourceReputationProfileCreate(ctx context.Context, d *schema.ResourceData, m interface{}) diag.Diagnostics {
	meta := akamai.Meta(m)
	client := inst.Client(meta)
	logger := meta.Log("APPSEC", "resourceReputationProfileCreate")

	createReputationProfile := appsec.CreateReputationProfileRequest{}

	jsonpostpayload := d.Get("reputation_profile")

	if err := json.Unmarshal([]byte(jsonpostpayload.(string)), &createReputationProfile); err != nil {
		return diag.FromErr(err)
	}

	configid, err := tools.GetIntValue("config_id", d)
	if err != nil && !errors.Is(err, tools.ErrNotFound) {
		return diag.FromErr(err)
	}
	createReputationProfile.ConfigID = configid

	configversion, err := tools.GetIntValue("version", d)
	if err != nil && !errors.Is(err, tools.ErrNotFound) {
		return diag.FromErr(err)
	}
	createReputationProfile.ConfigVersion = configversion

	postresp, errc := client.CreateReputationProfile(ctx, createReputationProfile)
	if errc != nil {
		logger.Errorf("calling 'createReputationProfile': %s", errc.Error())
		return diag.FromErr(errc)
	}

	d.SetId(strconv.Itoa(postresp.ID))

	return resourceReputationProfileRead(ctx, d, m)
}

func resourceReputationProfileUpdate(ctx context.Context, d *schema.ResourceData, m interface{}) diag.Diagnostics {
	meta := akamai.Meta(m)
	client := inst.Client(meta)
	logger := meta.Log("APPSEC", "resourceReputationProfileUpdate")

	updateReputationProfile := appsec.UpdateReputationProfileRequest{}

	jsonpostpayload := d.Get("reputation_profile")

	if err := json.Unmarshal([]byte(jsonpostpayload.(string)), &updateReputationProfile); err != nil {
		return diag.FromErr(err)
	}

	configid, err := tools.GetIntValue("config_id", d)
	if err != nil && !errors.Is(err, tools.ErrNotFound) {
		return diag.FromErr(err)
	}
	updateReputationProfile.ConfigID = configid

	configversion, err := tools.GetIntValue("version", d)
	if err != nil && !errors.Is(err, tools.ErrNotFound) {
		return diag.FromErr(err)
	}
	updateReputationProfile.ConfigVersion = configversion

	reputationProfileId, errconv := strconv.Atoi(d.Id())

	if errconv != nil {
		return diag.FromErr(errconv)
	}
	updateReputationProfile.ReputationProfileId = reputationProfileId

	_, erru := client.UpdateReputationProfile(ctx, updateReputationProfile)
	if erru != nil {
		logger.Errorf("calling 'updateReputationProfile': %s", erru.Error())
		return diag.FromErr(erru)
	}

	return resourceReputationProfileRead(ctx, d, m)
}

func resourceReputationProfileDelete(ctx context.Context, d *schema.ResourceData, m interface{}) diag.Diagnostics {
	meta := akamai.Meta(m)
	client := inst.Client(meta)
	logger := meta.Log("APPSEC", "resourceReputationProfileRemove")

	removeReputationProfile := appsec.RemoveReputationProfileRequest{}

	configid, err := tools.GetIntValue("config_id", d)
	if err != nil && !errors.Is(err, tools.ErrNotFound) {
		return diag.FromErr(err)
	}
	removeReputationProfile.ConfigID = configid

	configversion, err := tools.GetIntValue("version", d)
	if err != nil && !errors.Is(err, tools.ErrNotFound) {
		return diag.FromErr(err)
	}
	removeReputationProfile.ConfigVersion = configversion

	reputationProfileId, errconv := strconv.Atoi(d.Id())

	if errconv != nil {
		return diag.FromErr(errconv)
	}
	removeReputationProfile.ReputationProfileId = reputationProfileId

	_, errd := client.RemoveReputationProfile(ctx, removeReputationProfile)
	if errd != nil {
		logger.Errorf("calling 'removeReputationProfile': %s", errd.Error())
		return diag.FromErr(errd)
	}

	d.SetId("")

	return nil
}

func resourceReputationProfileRead(ctx context.Context, d *schema.ResourceData, m interface{}) diag.Diagnostics {
	meta := akamai.Meta(m)
	client := inst.Client(meta)
	logger := meta.Log("APPSEC", "resourceReputationProfileRead")

	getReputationProfile := appsec.GetReputationProfileRequest{}

	configid, err := tools.GetIntValue("config_id", d)
	if err != nil && !errors.Is(err, tools.ErrNotFound) {
		return diag.FromErr(err)
	}
	getReputationProfile.ConfigID = configid

	configversion, err := tools.GetIntValue("version", d)
	if err != nil && !errors.Is(err, tools.ErrNotFound) {
		return diag.FromErr(err)
	}
	getReputationProfile.ConfigVersion = configversion

	reputationProfileId, errconv := strconv.Atoi(d.Id())

	if errconv != nil {
		return diag.FromErr(errconv)
	}
	getReputationProfile.ReputationProfileId = reputationProfileId

	reputationprofile, err := client.GetReputationProfile(ctx, getReputationProfile)
	if err != nil {
		logger.Errorf("calling 'getReputationProfile': %s", err.Error())
		return diag.FromErr(err)
	}

	if err := d.Set("reputation_profile_id", reputationprofile.ID); err != nil {
		return diag.FromErr(fmt.Errorf("%w: %s", tools.ErrValueSet, err.Error()))
	}

	d.SetId(strconv.Itoa(reputationprofile.ID))

	return nil
}
