package appsec

import (
	"context"
	"errors"
	"fmt"
	"strconv"

	"github.com/akamai/AkamaiOPEN-edgegrid-golang/v2/pkg/appsec"
	"github.com/akamai/terraform-provider-akamai/v2/pkg/akamai"
	"github.com/akamai/terraform-provider-akamai/v2/pkg/tools"
	"github.com/hashicorp/terraform-plugin-sdk/v2/diag"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/customdiff"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/schema"
)

// appsec v1
//
// https://developer.akamai.com/api/cloud_security/application_security/v1.html
func resourceIPGeoProtection() *schema.Resource {
	return &schema.Resource{
		CreateContext: resourceIPGeoProtectionCreate,
		ReadContext:   resourceIPGeoProtectionRead,
		UpdateContext: resourceIPGeoProtectionUpdate,
		DeleteContext: resourceIPGeoProtectionDelete,
		CustomizeDiff: customdiff.All(
			VerifyIDUnchanged,
		),
		Importer: &schema.ResourceImporter{
			State: schema.ImportStatePassthrough,
		},
		Schema: map[string]*schema.Schema{
			"config_id": {
				Type:     schema.TypeInt,
				Required: true,
			},
			"security_policy_id": {
				Type:     schema.TypeString,
				Required: true,
			},
			"enabled": {
				Type:     schema.TypeBool,
				Required: true,
			},
			"output_text": {
				Type:        schema.TypeString,
				Computed:    true,
				Description: "Text Export representation",
			},
		},
	}
}

func resourceIPGeoProtectionCreate(ctx context.Context, d *schema.ResourceData, m interface{}) diag.Diagnostics {
	meta := akamai.Meta(m)
	client := inst.Client(meta)
	logger := meta.Log("APPSEC", "resourceIPGeoProtectionCreate")
	logger.Debugf("!!! in resourceIPGeoProtectionCreate")

	configid, err := tools.GetIntValue("config_id", d)
	if err != nil && !errors.Is(err, tools.ErrNotFound) {
		return diag.FromErr(err)
	}
	version := getModifiableConfigVersion(ctx, configid, "ipgeoProtection", m)
	policyid, err := tools.GetStringValue("security_policy_id", d)
	if err != nil && !errors.Is(err, tools.ErrNotFound) {
		return diag.FromErr(err)
	}
	applynetworkcontrols, err := tools.GetBoolValue("enabled", d)
	if err != nil && !errors.Is(err, tools.ErrNotFound) {
		return diag.FromErr(err)
	}

	createnetworkProtection := appsec.UpdateIPGeoProtectionRequest{}
	createnetworkProtection.ConfigID = configid
	createnetworkProtection.Version = version
	createnetworkProtection.PolicyID = policyid
	createnetworkProtection.ApplyNetworkLayerControls = applynetworkcontrols

	_, erru := client.UpdateIPGeoProtection(ctx, createnetworkProtection)
	if erru != nil {
		logger.Errorf("calling 'createnetworkProtection': %s", erru.Error())
		return diag.FromErr(erru)
	}

	d.SetId(fmt.Sprintf("%d:%s", createnetworkProtection.ConfigID, createnetworkProtection.PolicyID))

	return resourceIPGeoProtectionRead(ctx, d, m)
}

func resourceIPGeoProtectionRead(ctx context.Context, d *schema.ResourceData, m interface{}) diag.Diagnostics {
	meta := akamai.Meta(m)
	client := inst.Client(meta)
	logger := meta.Log("APPSEC", "resourceIPGeoProtectionRead")
	logger.Debugf("!!! in resourceIPGeoProtectionRead")

	idParts, err := splitID(d.Id(), 2, "configid:securitypolicyid")
	if err != nil {
		return diag.FromErr(err)
	}
	configid, err := strconv.Atoi(idParts[0])
	if err != nil {
		return diag.FromErr(err)
	}
	version := getLatestConfigVersion(ctx, configid, m)
	policyid := idParts[1]

	getNetworkProtection := appsec.GetIPGeoProtectionRequest{}
	getNetworkProtection.ConfigID = configid
	getNetworkProtection.Version = version
	getNetworkProtection.PolicyID = policyid

	networkprotection, err := client.GetIPGeoProtection(ctx, getNetworkProtection)
	if err != nil {
		logger.Errorf("calling 'getNetworkProtection': %s", err.Error())
		return diag.FromErr(err)
	}

	if err := d.Set("config_id", getNetworkProtection.ConfigID); err != nil {
		return diag.FromErr(fmt.Errorf("%w: %s", tools.ErrValueSet, err.Error()))
	}
	if err := d.Set("security_policy_id", getNetworkProtection.PolicyID); err != nil {
		return diag.FromErr(fmt.Errorf("%w: %s", tools.ErrValueSet, err.Error()))
	}
	if err := d.Set("enabled", networkprotection.ApplyNetworkLayerControls); err != nil {
		return diag.FromErr(fmt.Errorf("%w: %s", tools.ErrValueSet, err.Error()))
	}
	ots := OutputTemplates{}
	InitTemplates(ots)
	outputtext, err := RenderTemplates(ots, "networkProtectionDS", networkprotection)
	if err == nil {
		if err := d.Set("output_text", outputtext); err != nil {
			return diag.FromErr(fmt.Errorf("%w: %s", tools.ErrValueSet, err.Error()))
		}
	}

	return nil
}

func resourceIPGeoProtectionUpdate(ctx context.Context, d *schema.ResourceData, m interface{}) diag.Diagnostics {
	meta := akamai.Meta(m)
	client := inst.Client(meta)
	logger := meta.Log("APPSEC", "resourceIPGeoProtectionUpdate")
	logger.Debugf("!!! in resourceIPGeoProtectionUpdate")

	idParts, err := splitID(d.Id(), 2, "configid:securitypolicyid")
	if err != nil {
		return diag.FromErr(err)
	}
	configid, err := strconv.Atoi(idParts[0])
	if err != nil {
		return diag.FromErr(err)
	}
	version := getModifiableConfigVersion(ctx, configid, "networkProtection", m)
	policyid := idParts[1]
	applynetworkcontrols, err := tools.GetBoolValue("enabled", d)
	if err != nil && !errors.Is(err, tools.ErrNotFound) {
		return diag.FromErr(err)
	}

	updateNetworkProtection := appsec.UpdateIPGeoProtectionRequest{}
	updateNetworkProtection.ConfigID = configid
	updateNetworkProtection.Version = version
	updateNetworkProtection.PolicyID = policyid
	updateNetworkProtection.ApplyNetworkLayerControls = applynetworkcontrols

	_, erru := client.UpdateIPGeoProtection(ctx, updateNetworkProtection)
	if erru != nil {
		logger.Errorf("calling 'updateNetworkProtection': %s", erru.Error())
		return diag.FromErr(erru)
	}

	return resourceIPGeoProtectionRead(ctx, d, m)
}

func resourceIPGeoProtectionDelete(ctx context.Context, d *schema.ResourceData, m interface{}) diag.Diagnostics {
	meta := akamai.Meta(m)
	client := inst.Client(meta)
	logger := meta.Log("APPSEC", "resourceIPGeoProtectionDelete")

	logger.Debugf("!!! in resourceIPGeoProtectionDelete")

	idParts, err := splitID(d.Id(), 2, "configid:securitypolicyid")
	if err != nil {
		return diag.FromErr(err)
	}
	configid, err := strconv.Atoi(idParts[0])
	if err != nil {
		return diag.FromErr(err)
	}
	version := getModifiableConfigVersion(ctx, configid, "ipgeoProtection", m)
	policyid := idParts[1]

	removeIPGeoProtection := appsec.UpdateIPGeoProtectionRequest{}
	removeIPGeoProtection.ConfigID = configid

	removeIPGeoProtection.Version = version
	removeIPGeoProtection.PolicyID = policyid
	removeIPGeoProtection.ApplyNetworkLayerControls = false

	_, errd := client.UpdateIPGeoProtection(ctx, removeIPGeoProtection)
	if errd != nil {
		logger.Errorf("calling 'removeIPGeoProtection': %s", errd.Error())
		return diag.FromErr(errd)
	}
	d.SetId("")
	return nil
}
