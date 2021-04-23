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
func resourceCustomDeny() *schema.Resource {
	return &schema.Resource{
		CreateContext: resourceCustomDenyCreate,
		ReadContext:   resourceCustomDenyRead,
		UpdateContext: resourceCustomDenyUpdate,
		DeleteContext: resourceCustomDenyDelete,
		Importer: &schema.ResourceImporter{
			State: schema.ImportStatePassthrough,
		},
		CustomizeDiff: customdiff.All(
			VerifyIdUnchanged,
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
				ValidateFunc:     validation.StringIsJSON,
				DiffSuppressFunc: suppressCustomDenyJsonDiffs,
			},
		},
	}
}

func resourceCustomDenyCreate(ctx context.Context, d *schema.ResourceData, m interface{}) diag.Diagnostics {
	meta := akamai.Meta(m)
	client := inst.Client(meta)
	logger := meta.Log("APPSEC", "resourceCustomDenyCreate")
	logger.Debugf("!!! in resourceCustomDenyCreate")

	configid, err := tools.GetIntValue("config_id", d)
	if err != nil {
		return diag.FromErr(err)
	}
	version := getModifiableConfigVersion(ctx, configid, "customDeny", m)
	jsonpostpayload := d.Get("custom_deny")
	jsonPayloadRaw := []byte(jsonpostpayload.(string))
	rawJSON := (json.RawMessage)(jsonPayloadRaw)

	createCustomDeny := appsec.CreateCustomDenyRequest{}
	createCustomDeny.ConfigID = configid
	createCustomDeny.Version = version
	createCustomDeny.JsonPayloadRaw = rawJSON

	postresp, errc := client.CreateCustomDeny(ctx, createCustomDeny)
	if errc != nil {
		logger.Errorf("calling 'createCustomDeny': %s", errc.Error())
		return diag.FromErr(errc)
	}

	d.SetId(fmt.Sprintf("%d:%s", configid, postresp.ID))

	return resourceCustomDenyRead(ctx, d, m)
}

func resourceCustomDenyRead(ctx context.Context, d *schema.ResourceData, m interface{}) diag.Diagnostics {
	meta := akamai.Meta(m)
	client := inst.Client(meta)
	logger := meta.Log("APPSEC", "resourceCustomDenyRead")
	logger.Debugf("!!! in resourceCustomDenyRead")

	idParts, err := splitID(d.Id(), 2, "configid:customdenyid")
	if err != nil {
		return diag.FromErr(err)
	}

	configid, err := strconv.Atoi(idParts[0])
	if err != nil {
		return diag.FromErr(err)
	}
	customdenyid := idParts[1]

	getCustomDeny := appsec.GetCustomDenyRequest{}
	getCustomDeny.ConfigID = configid
	getCustomDeny.Version = getLatestConfigVersion(ctx, configid, m)
	getCustomDeny.ID = customdenyid

	customdeny, err := client.GetCustomDeny(ctx, getCustomDeny)
	if err != nil {
		logger.Errorf("calling 'getCustomDeny': %s", err.Error())
		return diag.FromErr(err)
	}

	if err := d.Set("config_id", configid); err != nil {
		return diag.FromErr(fmt.Errorf("%w: %s", tools.ErrValueSet, err.Error()))
	}
	if err := d.Set("custom_deny_id", customdenyid); err != nil {
		return diag.FromErr(fmt.Errorf("%w: %s", tools.ErrValueSet, err.Error()))
	}
	jsonBody, err := json.Marshal(customdeny)
	if err != nil {
		return diag.FromErr(err)
	}
	if err := d.Set("custom_deny", string(jsonBody)); err != nil {
		return diag.FromErr(fmt.Errorf("%w: %s", tools.ErrValueSet, err.Error()))
	}

	return nil
}

func resourceCustomDenyUpdate(ctx context.Context, d *schema.ResourceData, m interface{}) diag.Diagnostics {
	meta := akamai.Meta(m)
	client := inst.Client(meta)
	logger := meta.Log("APPSEC", "resourceCustomDenyUpdate")
	logger.Debugf("!!! in resourceCustomDenyUpdate")

	idParts, err := splitID(d.Id(), 2, "configid:customdenyid")
	if err != nil {
		return diag.FromErr(err)
	}

	configid, err := strconv.Atoi(idParts[0])
	if err != nil {
		return diag.FromErr(err)
	}
	customdenyid := idParts[1]

	jsonpostpayload := d.Get("custom_deny")
	jsonPayloadRaw := []byte(jsonpostpayload.(string))
	rawJSON := (json.RawMessage)(jsonPayloadRaw)

	updateCustomDeny := appsec.UpdateCustomDenyRequest{}
	updateCustomDeny.ConfigID = configid
	updateCustomDeny.Version = getModifiableConfigVersion(ctx, configid, "customDeny", m)
	updateCustomDeny.ID = customdenyid
	updateCustomDeny.JsonPayloadRaw = rawJSON

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
	logger := meta.Log("APPSEC", "resourceCustomDenyRemove")
	logger.Debugf("!!! in resourceCustomDenyDelete")

	idParts, err := splitID(d.Id(), 2, "configid:customdenyid")
	if err != nil {
		return diag.FromErr(err)
	}

	configid, err := strconv.Atoi(idParts[0])
	if err != nil {
		return diag.FromErr(err)
	}
	customdenyid := idParts[1]

	removeCustomDeny := appsec.RemoveCustomDenyRequest{}
	removeCustomDeny.ConfigID = configid
	removeCustomDeny.Version = getModifiableConfigVersion(ctx, configid, "customDeny", m)
	removeCustomDeny.ID = customdenyid

	_, err = client.RemoveCustomDeny(ctx, removeCustomDeny)
	if err != nil {
		logger.Errorf("calling 'removeCustomDeny': %s", err.Error())
		return diag.FromErr(err)
	}

	d.SetId("")

	return nil
}
