package appsec

import (
	"context"
	"encoding/json"
	"fmt"
	"log"
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
func resourceCustomDeny() *schema.Resource {
	return &schema.Resource{
		CreateContext: resourceCustomDenyCreate,
		ReadContext:   resourceCustomDenyRead,
		UpdateContext: resourceCustomDenyUpdate,
		DeleteContext: resourceCustomDenyDelete,
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
				//ValidateFunc: ValidateConfigID,
			},
			"custom_deny_id": {
				Type:        schema.TypeString,
				Computed:    true,
				Description: "custom_deny_id",
			},
			"custom_deny": {
				Type:             schema.TypeString,
				Required:         true,
				ValidateDiagFunc: validation.ToDiagFunc(validation.StringIsJSON),
				DiffSuppressFunc: suppressCustomDenyJSONDiffs,
			},
		},
	}
}

func resourceCustomDenyCreate(ctx context.Context, d *schema.ResourceData, m interface{}) diag.Diagnostics {
	meta := akamai.Meta(m)
	client := inst.Client(meta)
	logger := meta.Log("APPSEC", "resourceCustomDenyCreate")
	logger.Debugf("in resourceCustomDenyCreate")

	configID, err := tools.GetIntValue("config_id", d)
	if err != nil {
		return diag.FromErr(err)
	}
	version, err := getModifiableConfigVersion(ctx, configID, "customDeny", m)
	if err != nil {
		return diag.FromErr(err)
	}
	jsonpostpayload := d.Get("custom_deny")
	jsonPayloadRaw := []byte(jsonpostpayload.(string))
	rawJSON := (json.RawMessage)(jsonPayloadRaw)

	createCustomDeny := appsec.CreateCustomDenyRequest{
		ConfigID:       configID,
		Version:        version,
		JsonPayloadRaw: rawJSON,
	}

	createCustomDenyResponse, err := client.CreateCustomDeny(ctx, createCustomDeny)
	if err != nil {
		logger.Errorf("calling 'createCustomDeny': %s", err.Error())
		return diag.FromErr(err)
	}
	for _, p := range createCustomDenyResponse.Parameters {
		name := p.Name
		val := p.Value
		log.Printf("%s = %s", string(name), string(val))
	}

	d.SetId(fmt.Sprintf("%d:%s", configID, createCustomDenyResponse.ID))

	return resourceCustomDenyRead(ctx, d, m)
}

func resourceCustomDenyRead(ctx context.Context, d *schema.ResourceData, m interface{}) diag.Diagnostics {
	meta := akamai.Meta(m)
	client := inst.Client(meta)
	logger := meta.Log("APPSEC", "resourceCustomDenyRead")
	logger.Debugf("in resourceCustomDenyRead")

	iDParts, err := splitID(d.Id(), 2, "configID:customDenyID")
	if err != nil {
		return diag.FromErr(err)
	}

	configID, err := strconv.Atoi(iDParts[0])
	if err != nil {
		return diag.FromErr(err)
	}
	customDenyID := iDParts[1]

	version, err := getLatestConfigVersion(ctx, configID, m)
	if err != nil {
		return diag.FromErr(err)
	}

	getCustomDeny := appsec.GetCustomDenyRequest{
		ConfigID: configID,
		Version:  version,
		ID:       customDenyID,
	}

	getCustomDenyResponse, err := client.GetCustomDeny(ctx, getCustomDeny)
	if err != nil {
		logger.Errorf("calling 'getCustomDeny': %s", err.Error())
		return diag.FromErr(err)
	}
	for _, p := range getCustomDenyResponse.Parameters {
		name := p.Name
		val := p.Value
		log.Printf("%s = %s", string(name), string(val))
	}

	if err := d.Set("config_id", configID); err != nil {
		return diag.Errorf("%s: %s", tools.ErrValueSet, err.Error())
	}
	if err := d.Set("custom_deny_id", customDenyID); err != nil {
		return diag.Errorf("%s: %s", tools.ErrValueSet, err.Error())
	}
	jsonBody, err := json.Marshal(getCustomDenyResponse)
	if err != nil {
		return diag.FromErr(err)
	}
	if err := d.Set("custom_deny", string(jsonBody)); err != nil {
		return diag.Errorf("%s: %s", tools.ErrValueSet, err.Error())
	}

	return nil
}

func resourceCustomDenyUpdate(ctx context.Context, d *schema.ResourceData, m interface{}) diag.Diagnostics {
	meta := akamai.Meta(m)
	client := inst.Client(meta)
	logger := meta.Log("APPSEC", "resourceCustomDenyUpdate")
	logger.Debugf("in resourceCustomDenyUpdate")

	iDParts, err := splitID(d.Id(), 2, "configID:customDenyID")
	if err != nil {
		return diag.FromErr(err)
	}

	configID, err := strconv.Atoi(iDParts[0])
	if err != nil {
		return diag.FromErr(err)
	}
	customDenyID := iDParts[1]

	jsonpostpayload := d.Get("custom_deny")
	jsonPayloadRaw := []byte(jsonpostpayload.(string))
	rawJSON := (json.RawMessage)(jsonPayloadRaw)

	version, err := getModifiableConfigVersion(ctx, configID, "customDeny", m)
	if err != nil {
		return diag.FromErr(err)
	}
	updateCustomDeny := appsec.UpdateCustomDenyRequest{
		ConfigID:       configID,
		Version:        version,
		ID:             customDenyID,
		JsonPayloadRaw: rawJSON,
	}

	_, err = client.UpdateCustomDeny(ctx, updateCustomDeny)
	if err != nil {
		logger.Errorf("calling 'updateCustomDeny': %s", err.Error())
		return diag.FromErr(err)
	}

	return resourceCustomDenyRead(ctx, d, m)
}

func resourceCustomDenyDelete(ctx context.Context, d *schema.ResourceData, m interface{}) diag.Diagnostics {
	meta := akamai.Meta(m)
	client := inst.Client(meta)
	logger := meta.Log("APPSEC", "resourceCustomDenyDelete")
	logger.Debugf("in resourceCustomDenyDelete")

	iDParts, err := splitID(d.Id(), 2, "configID:customDenyID")
	if err != nil {
		return diag.FromErr(err)
	}

	configID, err := strconv.Atoi(iDParts[0])
	if err != nil {
		return diag.FromErr(err)
	}
	customDenyID := iDParts[1]

	version, err := getModifiableConfigVersion(ctx, configID, "customDeny", m)
	if err != nil {
		return diag.FromErr(err)
	}
	removeCustomDeny := appsec.RemoveCustomDenyRequest{
		ConfigID: configID,
		Version:  version,
		ID:       customDenyID,
	}

	_, err = client.RemoveCustomDeny(ctx, removeCustomDeny)
	if err != nil {
		logger.Errorf("calling 'removeCustomDeny': %s", err.Error())
		return diag.FromErr(err)
	}

	d.SetId("")

	return nil
}
