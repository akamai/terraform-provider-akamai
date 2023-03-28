package appsec

import (
	"context"
	"encoding/json"
	"fmt"
	"strconv"

	"github.com/akamai/AkamaiOPEN-edgegrid-golang/v5/pkg/appsec"
	"github.com/akamai/terraform-provider-akamai/v3/pkg/akamai"
	"github.com/akamai/terraform-provider-akamai/v3/pkg/tools"
	"github.com/hashicorp/terraform-plugin-sdk/v2/diag"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/customdiff"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/schema"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/validation"
)

// appsec v1
//
// https://techdocs.akamai.com/application-security/reference/api
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
				Type:        schema.TypeInt,
				Required:    true,
				Description: "Unique identifier of the security configuration",
			},
			"match_target": {
				Type:             schema.TypeString,
				Required:         true,
				ValidateDiagFunc: validation.ToDiagFunc(validation.StringIsJSON),
				DiffSuppressFunc: suppressEquivalentMatchTargetDiffs,
				Description:      "JSON-formatted definition of the match target",
			},
			"match_target_id": {
				Type:        schema.TypeInt,
				Computed:    true,
				Description: "Unique identifier of the match target",
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
	version, err := getModifiableConfigVersion(ctx, configID, "matchTarget", m)
	if err != nil {
		return diag.FromErr(err)
	}
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
	version, err := getLatestConfigVersion(ctx, configID, m)
	if err != nil {
		return diag.FromErr(err)
	}
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
	version, err := getModifiableConfigVersion(ctx, configID, "matchTarget", m)
	if err != nil {
		return diag.FromErr(err)
	}
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
	version, err := getModifiableConfigVersion(ctx, configID, "matchTarget", m)
	if err != nil {
		return diag.FromErr(err)
	}
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
	return nil
}
