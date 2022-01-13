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
func resourceMatchTarget() *schema.Resource {
	return &schema.Resource{
		CreateContext: resourceMatchTargetCreate,
		ReadContext:   resourceMatchTargetRead,
		UpdateContext: resourceMatchTargetUpdate,
		DeleteContext: resourceMatchTargetDelete,
		CustomizeDiff: customdiff.All(
			VerifyIDUnchanged,
		),
		Importer: &schema.ResourceImporter{
			StateContext: schema.ImportStatePassthroughContext,
		},
		Schema: map[string]*schema.Schema{
			"config_id": {
				Type:     schema.TypeInt,
				Required: true,
			},
			"match_target_id": {
				Type:     schema.TypeInt,
				Computed: true,
			},
			"match_target": {
				Type:             schema.TypeString,
				Required:         true,
				ValidateDiagFunc: validation.ToDiagFunc(validation.StringIsJSON),
				DiffSuppressFunc: suppressEquivalentMatchTargetDiffs,
			},
		},
	}
}

func resourceMatchTargetCreate(ctx context.Context, d *schema.ResourceData, m interface{}) diag.Diagnostics {
	meta := akamai.Meta(m)
	client := inst.Client(meta)
	logger := meta.Log("APPSEC", "resourceMatchTargetCreate")
	logger.Debugf("in resourceMatchTargetCreate")

	configID, err := tools.GetIntValue("config_id", d)
	if err != nil {
		return diag.FromErr(err)
	}
	version := getModifiableConfigVersion(ctx, configID, "matchTarget", m)
	createMatchTarget := appsec.CreateMatchTargetRequest{}
	jsonpostpayload := d.Get("match_target")
	jsonPayloadRaw := []byte(jsonpostpayload.(string))
	rawJSON := (json.RawMessage)(jsonPayloadRaw)

	createMatchTarget.ConfigID = configID
	createMatchTarget.ConfigVersion = version
	createMatchTarget.JsonPayloadRaw = rawJSON

	postresp, err := client.CreateMatchTarget(ctx, createMatchTarget)
	if err != nil {
		logger.Errorf("calling 'createMatchTarget': %s", err.Error())
		return diag.FromErr(err)
	}

	d.SetId(fmt.Sprintf("%d:%d", createMatchTarget.ConfigID, postresp.TargetID))

	return resourceMatchTargetRead(ctx, d, m)
}

func resourceMatchTargetRead(ctx context.Context, d *schema.ResourceData, m interface{}) diag.Diagnostics {
	meta := akamai.Meta(m)
	client := inst.Client(meta)
	logger := meta.Log("APPSEC", "resourceMatchTargetRead")
	logger.Debugf("in resourceMatchTargetRead")

	iDParts, err := splitID(d.Id(), 2, "configID:matchTargetID")
	if err != nil {
		return diag.FromErr(err)
	}

	configID, err := strconv.Atoi(iDParts[0])
	if err != nil {
		return diag.FromErr(err)
	}
	version := getLatestConfigVersion(ctx, configID, m)
	targetID, err := strconv.Atoi(iDParts[1])
	if err != nil {
		return diag.FromErr(err)
	}

	getMatchTarget := appsec.GetMatchTargetRequest{
		ConfigID:      configID,
		ConfigVersion: version,
		TargetID:      targetID,
	}

	matchtarget, err := client.GetMatchTarget(ctx, getMatchTarget)
	if err != nil {
		logger.Errorf("calling 'getMatchTarget': %s", err.Error())
		return diag.FromErr(err)
	}

	jsonBody, err := json.Marshal(matchtarget)
	if err != nil {
		return diag.FromErr(err)
	}
	if err := d.Set("config_id", configID); err != nil {
		return diag.Errorf("%s: %s", tools.ErrValueSet, err.Error())
	}
	if err := d.Set("match_target", string(jsonBody)); err != nil {
		return diag.Errorf("%s: %s", tools.ErrValueSet, err.Error())
	}
	if err := d.Set("match_target_id", matchtarget.TargetID); err != nil {
		return diag.Errorf("%s: %s", tools.ErrValueSet, err.Error())
	}

	return nil
}

func resourceMatchTargetUpdate(ctx context.Context, d *schema.ResourceData, m interface{}) diag.Diagnostics {
	meta := akamai.Meta(m)
	client := inst.Client(meta)
	logger := meta.Log("APPSEC", "resourceMatchTargetUpdate")
	logger.Debugf("in resourceMatchTargetUpdate")

	iDParts, err := splitID(d.Id(), 2, "configID:matchTargetID")
	if err != nil {
		return diag.FromErr(err)
	}

	configID, err := strconv.Atoi(iDParts[0])
	if err != nil {
		return diag.FromErr(err)
	}
	version := getModifiableConfigVersion(ctx, configID, "matchTarget", m)
	targetID, err := strconv.Atoi(iDParts[1])
	if err != nil {
		return diag.FromErr(err)
	}
	jsonpostpayload := d.Get("match_target")
	jsonPayloadRaw := []byte(jsonpostpayload.(string))
	rawJSON := (json.RawMessage)(jsonPayloadRaw)

	updateMatchTarget := appsec.UpdateMatchTargetRequest{
		ConfigID:       configID,
		ConfigVersion:  version,
		TargetID:       targetID,
		JsonPayloadRaw: rawJSON,
	}

	_, err = client.UpdateMatchTarget(ctx, updateMatchTarget)
	if err != nil {
		logger.Errorf("calling 'updateMatchTarget': %s", err.Error())
		return diag.FromErr(err)
	}

	return resourceMatchTargetRead(ctx, d, m)
}

func resourceMatchTargetDelete(ctx context.Context, d *schema.ResourceData, m interface{}) diag.Diagnostics {
	meta := akamai.Meta(m)
	client := inst.Client(meta)
	logger := meta.Log("APPSEC", "resourceMatchTargetDelete")
	logger.Debugf("in resourceMatchTargetDelete")

	iDParts, err := splitID(d.Id(), 2, "configID:matchTargetID")
	if err != nil {
		return diag.FromErr(err)
	}

	configID, err := strconv.Atoi(iDParts[0])
	if err != nil {
		return diag.FromErr(err)
	}
	version := getModifiableConfigVersion(ctx, configID, "matchTarget", m)
	targetID, err := strconv.Atoi(iDParts[1])
	if err != nil {
		return diag.FromErr(err)
	}

	removeMatchTarget := appsec.RemoveMatchTargetRequest{
		ConfigID:      configID,
		ConfigVersion: version,
		TargetID:      targetID,
	}

	_, err = client.RemoveMatchTarget(ctx, removeMatchTarget)
	if err != nil {
		logger.Errorf("calling 'removeMatchTarget': %s", err.Error())
		return diag.FromErr(err)
	}

	d.SetId("")

	return nil
}
