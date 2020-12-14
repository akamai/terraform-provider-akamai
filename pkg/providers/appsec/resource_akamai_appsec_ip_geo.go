package appsec

import (
	"context"
	"errors"
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
func resourceIPGeo() *schema.Resource {
	return &schema.Resource{
		CreateContext: resourceIPGeoUpdate,
		ReadContext:   resourceIPGeoRead,
		UpdateContext: resourceIPGeoUpdate,
		DeleteContext: resourceIPGeoDelete,
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
			"security_policy_id": {
				Type:     schema.TypeString,
				Required: true,
			},
			"mode": {
				Type:     schema.TypeString,
				Required: true,
				ValidateFunc: validation.StringInSlice([]string{
					Allow,
					Block,
				}, false),
			},
			"geo_network_lists": {
				Type:     schema.TypeSet,
				Optional: true,
				Elem:     &schema.Schema{Type: schema.TypeString},
			},
			"ip_network_lists": {
				Type:     schema.TypeSet,
				Optional: true,
				Elem:     &schema.Schema{Type: schema.TypeString},
			},
			"exception_ip_network_lists": {
				Type:     schema.TypeSet,
				Optional: true,
				Elem:     &schema.Schema{Type: schema.TypeString},
			},
			"output_text": {
				Type:        schema.TypeString,
				Computed:    true,
				Description: "Text Export representation",
			},
		},
	}
}

func resourceIPGeoRead(ctx context.Context, d *schema.ResourceData, m interface{}) diag.Diagnostics {
	meta := akamai.Meta(m)
	client := inst.Client(meta)
	logger := meta.Log("APPSEC", "resourceIPGeoRead")

	getIPGeo := appsec.GetIPGeoRequest{}

	configid, err := tools.GetIntValue("config_id", d)
	if err != nil && !errors.Is(err, tools.ErrNotFound) {
		return diag.FromErr(err)
	}
	getIPGeo.ConfigID = configid

	version, err := tools.GetIntValue("version", d)
	if err != nil && !errors.Is(err, tools.ErrNotFound) {
		return diag.FromErr(err)
	}
	getIPGeo.Version = version

	policyid, err := tools.GetStringValue("security_policy_id", d)
	if err != nil && !errors.Is(err, tools.ErrNotFound) {
		return diag.FromErr(err)
	}
	getIPGeo.PolicyID = policyid

	ipgeo, err := client.GetIPGeo(ctx, getIPGeo)
	if err != nil {
		logger.Errorf("calling 'getIPGeo': %s", err.Error())
		return diag.FromErr(err)
	}

	ots := OutputTemplates{}
	InitTemplates(ots)

	outputtext, err := RenderTemplates(ots, "IPGeoDS", ipgeo)
	if err == nil {
		d.Set("output_text", outputtext)
	}

	d.SetId(strconv.Itoa(getIPGeo.ConfigID))

	return nil
}

func resourceIPGeoDelete(ctx context.Context, d *schema.ResourceData, m interface{}) diag.Diagnostics {

	meta := akamai.Meta(m)
	client := inst.Client(meta)
	logger := meta.Log("APPSEC", "resourceIPGeoDelete")

	updatePolicyProtections := appsec.UpdateNetworkLayerProtectionRequest{}

	configid, err := tools.GetIntValue("config_id", d)
	if err != nil && !errors.Is(err, tools.ErrNotFound) {
		return diag.FromErr(err)
	}
	updatePolicyProtections.ConfigID = configid

	version, err := tools.GetIntValue("version", d)
	if err != nil && !errors.Is(err, tools.ErrNotFound) {
		return diag.FromErr(err)
	}
	updatePolicyProtections.Version = version

	policyid, err := tools.GetStringValue("security_policy_id", d)
	if err != nil && !errors.Is(err, tools.ErrNotFound) {
		return diag.FromErr(err)
	}
	updatePolicyProtections.PolicyID = policyid

	updatePolicyProtections.ApplyNetworkLayerControls = false

	logger.Errorf("calling 'resourceIPGeoDelete': %v", updatePolicyProtections)
	_, erru := client.UpdateNetworkLayerProtection(ctx, updatePolicyProtections)
	if erru != nil {
		logger.Errorf("calling 'resourceIPGeoDelete': %s", erru.Error())
		return diag.FromErr(erru)
	}
	d.SetId("")
	return nil
}

func resourceIPGeoUpdate(ctx context.Context, d *schema.ResourceData, m interface{}) diag.Diagnostics {
	meta := akamai.Meta(m)
	client := inst.Client(meta)
	logger := meta.Log("APPSEC", "resourceIPGeoUpdate")

	updateIPGeo := appsec.UpdateIPGeoRequest{}

	configid, err := tools.GetIntValue("config_id", d)
	if err != nil && !errors.Is(err, tools.ErrNotFound) {
		return diag.FromErr(err)
	}
	updateIPGeo.ConfigID = configid

	version, err := tools.GetIntValue("version", d)
	if err != nil && !errors.Is(err, tools.ErrNotFound) {
		return diag.FromErr(err)
	}
	updateIPGeo.Version = version

	policyid, err := tools.GetStringValue("security_policy_id", d)
	if err != nil && !errors.Is(err, tools.ErrNotFound) {
		return diag.FromErr(err)
	}
	updateIPGeo.PolicyID = policyid

	mode, err := tools.GetStringValue("mode", d)
	if err != nil && !errors.Is(err, tools.ErrNotFound) {
		return diag.FromErr(err)
	}

	if mode == Allow {
		updateIPGeo.Block = "blockAllTrafficExceptAllowedIPs"
	}

	if mode == Block {
		updateIPGeo.Block = "blockSpecificIPGeo"
	}

	updateIPGeo.GeoControls.BlockedIPNetworkLists.NetworkList = tools.SetToStringSlice(d.Get("geo_network_lists").(*schema.Set))
	updateIPGeo.IPControls.BlockedIPNetworkLists.NetworkList = tools.SetToStringSlice(d.Get("ip_network_lists").(*schema.Set))
	updateIPGeo.IPControls.AllowedIPNetworkLists.NetworkList = tools.SetToStringSlice(d.Get("exception_ip_network_lists").(*schema.Set))

	_, erru := client.UpdateIPGeo(ctx, updateIPGeo)
	if erru != nil {
		logger.Errorf("calling 'updateIPGeo': %s", erru.Error())
		return diag.FromErr(erru)
	}

	return resourceIPGeoRead(ctx, d, m)
}

const (
	Allow = "allow"
	Block = "block"
)
