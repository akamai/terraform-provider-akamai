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
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/validation"
)

// appsec v1
//
// https://techdocs.akamai.com/application-security/reference/api
func resourceIPGeo() *schema.Resource {
	return &schema.Resource{
		CreateContext: resourceIPGeoCreate,
		ReadContext:   resourceIPGeoRead,
		UpdateContext: resourceIPGeoUpdate,
		DeleteContext: resourceIPGeoDelete,
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
			"security_policy_id": {
				Type:        schema.TypeString,
				Required:    true,
				Description: "Unique identifier of the security policy",
			},
			"mode": {
				Type:     schema.TypeString,
				Required: true,
				ValidateDiagFunc: validation.ToDiagFunc(validation.StringInSlice([]string{
					Allow,
					Block,
				}, false)),
				Description: "Protection mode (block or allow)",
			},
			"geo_network_lists": {
				Type:        schema.TypeSet,
				Optional:    true,
				Elem:        &schema.Schema{Type: schema.TypeString},
				Description: "List of IDs of geographic network list to be blocked or allowed",
			},
			"ip_network_lists": {
				Type:        schema.TypeSet,
				Optional:    true,
				Elem:        &schema.Schema{Type: schema.TypeString},
				Description: "List of IDs of IP network list to be blocked or allowed",
			},
			"exception_ip_network_lists": {
				Type:        schema.TypeSet,
				Optional:    true,
				Elem:        &schema.Schema{Type: schema.TypeString},
				Description: "List of IDs of network list that are always allowed",
			},
		},
	}
}
func resourceIPGeoCreate(ctx context.Context, d *schema.ResourceData, m interface{}) diag.Diagnostics {
	meta := akamai.Meta(m)
	client := inst.Client(meta)
	logger := meta.Log("APPSEC", "resourceIPGeoCreate")
	logger.Debugf("in resourceIPGeoCreate")

	configID, err := tools.GetIntValue("config_id", d)
	if err != nil {
		return diag.FromErr(err)
	}
	version, err := getModifiableConfigVersion(ctx, configID, "ipgeo", m)
	if err != nil {
		return diag.FromErr(err)
	}
	policyID, err := tools.GetStringValue("security_policy_id", d)
	if err != nil {
		return diag.FromErr(err)
	}
	mode, err := tools.GetStringValue("mode", d)
	if err != nil {
		return diag.FromErr(err)
	}
	blockedGeoListsSet, err := tools.GetSetValue("geo_network_lists", d)
	if err != nil && !errors.Is(err, tools.ErrNotFound) {
		return diag.FromErr(err)
	}
	blockedIPListsSet, err := tools.GetSetValue("ip_network_lists", d)
	if err != nil && !errors.Is(err, tools.ErrNotFound) {
		return diag.FromErr(err)
	}
	exceptionIPListsSet, err := tools.GetSetValue("exception_ip_network_lists", d)
	if err != nil && !errors.Is(err, tools.ErrNotFound) {
		return diag.FromErr(err)
	}

	createIPGeo := appsec.UpdateIPGeoRequest{
		ConfigID: configID,
		Version:  version,
		PolicyID: policyID,
	}
	if mode == Allow {
		createIPGeo.Block = "blockAllTrafficExceptAllowedIPs"
	}
	if mode == Block {
		createIPGeo.Block = "blockSpecificIPGeo"
	}
	createIPGeo.GeoControls.BlockedIPNetworkLists.NetworkList = tools.SetToStringSlice(blockedGeoListsSet)
	createIPGeo.IPControls.BlockedIPNetworkLists.NetworkList = tools.SetToStringSlice(blockedIPListsSet)
	createIPGeo.IPControls.AllowedIPNetworkLists.NetworkList = tools.SetToStringSlice(exceptionIPListsSet)

	_, err = client.UpdateIPGeo(ctx, createIPGeo)
	if err != nil {
		logger.Errorf("calling 'createIPGeo': %s", err.Error())
		return diag.FromErr(err)
	}

	d.SetId(fmt.Sprintf("%d:%s", createIPGeo.ConfigID, createIPGeo.PolicyID))

	return resourceIPGeoRead(ctx, d, m)
}

func resourceIPGeoRead(ctx context.Context, d *schema.ResourceData, m interface{}) diag.Diagnostics {
	meta := akamai.Meta(m)
	client := inst.Client(meta)
	logger := meta.Log("APPSEC", "resourceIPGeoRead")
	logger.Debugf("in resourceIPGeoRead")

	iDParts, err := splitID(d.Id(), 2, "configID:securityPolicyID")
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
	policyID := iDParts[1]

	getIPGeo := appsec.GetIPGeoRequest{
		ConfigID: configID,
		Version:  version,
		PolicyID: policyID,
	}

	ipgeo, err := client.GetIPGeo(ctx, getIPGeo)
	if err != nil {
		logger.Errorf("calling 'getIPGeo': %s", err.Error())
		return diag.FromErr(err)
	}

	if err := d.Set("config_id", getIPGeo.ConfigID); err != nil {
		return diag.Errorf("%s: %s", tools.ErrValueSet, err.Error())
	}
	if err := d.Set("security_policy_id", getIPGeo.PolicyID); err != nil {
		return diag.Errorf("%s: %s", tools.ErrValueSet, err.Error())
	}
	if ipgeo.Block == "blockAllTrafficExceptAllowedIPs" {
		if err := d.Set("mode", Allow); err != nil {
			return diag.Errorf("%s: %s", tools.ErrValueSet, err.Error())
		}
	}
	if ipgeo.Block == "blockSpecificIPGeo" {
		if err := d.Set("mode", Block); err != nil {
			return diag.Errorf("%s: %s", tools.ErrValueSet, err.Error())
		}
	}
	if err := d.Set("geo_network_lists", ipgeo.GeoControls.BlockedIPNetworkLists.NetworkList); err != nil {
		return diag.Errorf("%s: %s", tools.ErrValueSet, err.Error())
	}
	if err := d.Set("ip_network_lists", ipgeo.IPControls.BlockedIPNetworkLists.NetworkList); err != nil {
		return diag.Errorf("%s: %s", tools.ErrValueSet, err.Error())
	}
	if err := d.Set("exception_ip_network_lists", ipgeo.IPControls.AllowedIPNetworkLists.NetworkList); err != nil {
		return diag.Errorf("%s: %s", tools.ErrValueSet, err.Error())
	}

	return nil
}

func resourceIPGeoUpdate(ctx context.Context, d *schema.ResourceData, m interface{}) diag.Diagnostics {
	meta := akamai.Meta(m)
	client := inst.Client(meta)
	logger := meta.Log("APPSEC", "resourceIPGeoUpdate")
	logger.Debugf("in resourceIPGeoUpdate")

	iDParts, err := splitID(d.Id(), 2, "configID:securityPolicyID")
	if err != nil {
		return diag.FromErr(err)
	}
	configID, err := strconv.Atoi(iDParts[0])
	if err != nil {
		return diag.FromErr(err)
	}
	version, err := getModifiableConfigVersion(ctx, configID, "ipgeo", m)
	if err != nil {
		return diag.FromErr(err)
	}
	policyID := iDParts[1]
	mode, err := tools.GetStringValue("mode", d)
	if err != nil && !errors.Is(err, tools.ErrNotFound) {
		return diag.FromErr(err)
	}
	blockedgeolists := tools.SetToStringSlice(d.Get("geo_network_lists").(*schema.Set))
	blockediplists := tools.SetToStringSlice(d.Get("ip_network_lists").(*schema.Set))
	exceptioniplists := tools.SetToStringSlice(d.Get("exception_ip_network_lists").(*schema.Set))

	updateIPGeo := appsec.UpdateIPGeoRequest{
		ConfigID: configID,
		Version:  version,
		PolicyID: policyID,
	}
	if mode == Allow {
		updateIPGeo.Block = "blockAllTrafficExceptAllowedIPs"
	}
	if mode == Block {
		updateIPGeo.Block = "blockSpecificIPGeo"
	}

	updateIPGeo.GeoControls.BlockedIPNetworkLists.NetworkList = blockedgeolists
	updateIPGeo.IPControls.BlockedIPNetworkLists.NetworkList = blockediplists
	updateIPGeo.IPControls.AllowedIPNetworkLists.NetworkList = exceptioniplists

	_, err = client.UpdateIPGeo(ctx, updateIPGeo)
	if err != nil {
		logger.Errorf("calling 'updateIPGeo': %s", err.Error())
		return diag.FromErr(err)
	}

	return resourceIPGeoRead(ctx, d, m)
}

func resourceIPGeoDelete(ctx context.Context, d *schema.ResourceData, m interface{}) diag.Diagnostics {
	meta := akamai.Meta(m)
	client := inst.Client(meta)
	logger := meta.Log("APPSEC", "resourceIPGeoDelete")
	logger.Debugf("in resourceIPGeoDelete")

	iDParts, err := splitID(d.Id(), 2, "configID:securityPolicyID")
	if err != nil {
		return diag.FromErr(err)
	}
	configID, err := strconv.Atoi(iDParts[0])
	if err != nil {
		return diag.FromErr(err)
	}
	version, err := getModifiableConfigVersion(ctx, configID, "ipgeo", m)
	if err != nil {
		return diag.FromErr(err)
	}
	policyID := iDParts[1]

	_, err = client.UpdateIPGeoProtection(ctx, appsec.UpdateIPGeoProtectionRequest{
		ConfigID:                  configID,
		Version:                   version,
		PolicyID:                  policyID,
		ApplyNetworkLayerControls: false,
	})
	if err != nil {
		logger.Errorf("calling UpdateIPGeoProtection: %s", err.Error())
		return diag.FromErr(err)
	}
	return nil
}

// Definition of constant variables
const (
	Allow = "allow"
	Block = "block"
)
