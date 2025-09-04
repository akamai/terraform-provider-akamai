package appsec

import (
	"context"
	"errors"
	"fmt"
	"strconv"

	"github.com/akamai/AkamaiOPEN-edgegrid-golang/v12/pkg/appsec"
	"github.com/akamai/terraform-provider-akamai/v9/pkg/common/id"
	"github.com/akamai/terraform-provider-akamai/v9/pkg/common/tf"
	"github.com/akamai/terraform-provider-akamai/v9/pkg/meta"
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
			"geo_controls": {
				Type:     schema.TypeList,
				MaxItems: 1,
				Optional: true,
				Elem: &schema.Resource{
					Schema: map[string]*schema.Schema{
						"geo_network_lists": {
							Type: schema.TypeSet,
							Elem: &schema.Schema{
								Type:             schema.TypeString,
								ValidateDiagFunc: validateEmptyElementsInList,
							},
							Required:    true,
							MinItems:    1,
							Description: "List of IDs of geographic network list to be blocked.",
						},
						"action": {
							Type:        schema.TypeString,
							Optional:    true,
							Default:     "deny",
							Description: "Action set for GEO Controls.",
						},
					},
				},
				Description: "An Object containing List of geographic network lists to be blocked with specified action",
			},
			"ip_controls": {
				Type:     schema.TypeList,
				MaxItems: 1,
				Optional: true,
				Elem: &schema.Resource{
					Schema: map[string]*schema.Schema{
						"ip_network_lists": {
							Type: schema.TypeSet,
							Elem: &schema.Schema{
								Type:             schema.TypeString,
								ValidateDiagFunc: validateEmptyElementsInList,
							},
							Required:    true,
							MinItems:    1,
							Description: "List of IDs of IP network list to be blocked.",
						},
						"action": {
							Type:        schema.TypeString,
							Optional:    true,
							Default:     "deny",
							Description: "Action set for IP Controls.",
						},
					},
				},
				Description: "An Object containing List of IP network lists to be blocked with specified action",
			},
			"asn_controls": {
				Type:     schema.TypeList,
				MaxItems: 1,
				Optional: true,
				Elem: &schema.Resource{
					Schema: map[string]*schema.Schema{
						"asn_network_lists": {
							Type: schema.TypeSet,
							Elem: &schema.Schema{
								Type:             schema.TypeString,
								ValidateDiagFunc: validateEmptyElementsInList,
							},
							Required:    true,
							MinItems:    1,
							Description: "List of IDs of ASN network list to be blocked.",
						},
						"action": {
							Type:        schema.TypeString,
							Optional:    true,
							Default:     "deny",
							Description: "Action set for ASN Controls",
						},
					},
				},
				Description: "An Object containing List of ASN network lists to be blocked with specified action",
			},
			"exception_ip_network_lists": {
				Type:     schema.TypeSet,
				Optional: true,
				Elem: &schema.Schema{
					Type:             schema.TypeString,
					ValidateDiagFunc: validateEmptyElementsInList,
				},
				Description: "List of unique identifiers of ip_network_lists allowed through the firewall.",
			},
			"ukraine_geo_control_action": {
				Type:             schema.TypeString,
				Optional:         true,
				Computed:         true,
				DiffSuppressFunc: suppressDiffUkraineGeoControlAction,
				Description:      "Action set for Ukraine geo control",
			},
			"block_action": {
				Type:        schema.TypeString,
				Optional:    true,
				Description: "Specifies the action set for BLOCK Mode blocking all the traffic except from lists identified in exception_ip_network_lists",
			},
		},
	}
}

func ipGeoNetworkListsFromStringList(strings []interface{}, action string) *appsec.IPGeoNetworkLists {
	items := make([]string, 0, len(strings))
	for _, item := range strings {
		items = append(items, item.(string))
	}
	return &appsec.IPGeoNetworkLists{
		NetworkList: items,
		Action:      action,
	}
}

func ipControlsFromAllowLists(exceptionIPLists []interface{}) *appsec.IPGeoIPControls {
	if len(exceptionIPLists) > 0 {
		return &appsec.IPGeoIPControls{
			AllowedIPNetworkLists: ipGeoNetworkListsFromStringList(exceptionIPLists, ""),
		}
	}
	return nil
}

func geoControlsFromBlockLists(blockedGeoLists []interface{}, geoAction string) *appsec.IPGeoGeoControls {
	if len(blockedGeoLists) > 0 {
		return &appsec.IPGeoGeoControls{
			BlockedIPNetworkLists: ipGeoNetworkListsFromStringList(blockedGeoLists, geoAction),
		}
	}
	return nil
}

func asnControlsFromBlockLists(blockedASNLists []interface{}, asnAction string) *appsec.IPGeoASNControls {
	if len(blockedASNLists) > 0 {
		return &appsec.IPGeoASNControls{
			BlockedIPNetworkLists: ipGeoNetworkListsFromStringList(blockedASNLists, asnAction),
		}
	}
	return nil
}

func ipControlsFromBlockAndAllowLists(blockedIPLists []interface{}, ipAction string, exceptionIPLists []interface{}) *appsec.IPGeoIPControls {
	if len(blockedIPLists) > 0 || len(exceptionIPLists) > 0 {
		ipControls := &appsec.IPGeoIPControls{}
		if len(blockedIPLists) > 0 {
			ipControls.BlockedIPNetworkLists = ipGeoNetworkListsFromStringList(blockedIPLists, ipAction)
		}
		if len(exceptionIPLists) > 0 {
			ipControls.AllowedIPNetworkLists = ipGeoNetworkListsFromStringList(exceptionIPLists, "")
		}
		return ipControls
	}
	return nil
}

func extractGeoControls(d *schema.ResourceData) (string, []interface{}, diag.Diagnostics) {
	geoControlsRaw := d.Get("geo_controls").([]interface{})

	var geoBlockAction string

	if len(geoControlsRaw) > 0 && geoControlsRaw[0] != nil {
		geoControl := geoControlsRaw[0].(map[string]interface{})

		// Extract action
		if v, ok := geoControl["action"].(string); ok {
			if v == "" {
				geoBlockAction = "deny" // fallback to default manually
			} else {
				geoBlockAction = v
			}
		}

		// Extract asn_network_lists as []interface{}
		if v, ok := geoControl["geo_network_lists"].(*schema.Set); ok {
			return geoBlockAction, v.List(), nil
		}
		return "", nil, diag.Errorf("%s: %s", tf.ErrValueSet, "Cannot parse geo_network_lists properly")
	}
	return "", nil, nil

}

func extractAsnControls(d *schema.ResourceData) (string, []interface{}, diag.Diagnostics) {
	asnControlsRaw := d.Get("asn_controls").([]interface{})

	var asnBlockAction string

	if len(asnControlsRaw) > 0 && asnControlsRaw[0] != nil {
		asnControl := asnControlsRaw[0].(map[string]interface{})

		// Extract action
		if v, ok := asnControl["action"].(string); ok {
			if v == "" {
				asnBlockAction = "deny" // fallback to default manually
			} else {
				asnBlockAction = v
			}
		}

		// Extract asn_network_lists as []interface{}
		if v, ok := asnControl["asn_network_lists"].(*schema.Set); ok {
			return asnBlockAction, v.List(), nil
		}
		return "", nil, diag.Errorf("%s: %s", tf.ErrValueSet, "Cannot parse asn_network_lists properly")
	}
	return "", nil, nil
}

func extractIPControls(d *schema.ResourceData) (string, []interface{}, diag.Diagnostics) {
	ipControlsRaw := d.Get("ip_controls").([]interface{})

	var ipBlockAction string

	if len(ipControlsRaw) > 0 && ipControlsRaw[0] != nil {
		ipControl := ipControlsRaw[0].(map[string]interface{})

		// Extract action
		if v, ok := ipControl["action"].(string); ok {
			if v == "" {
				ipBlockAction = "deny" // fallback to default manually
			} else {
				ipBlockAction = v
			}
		}

		// Extract asn_network_lists as []interface{}
		if v, ok := ipControl["ip_network_lists"].(*schema.Set); ok {
			return ipBlockAction, v.List(), nil
		}
		return "", nil, diag.Errorf("%s: %s", tf.ErrValueSet, "Cannot parse ip_network_lists properly")

	}
	return "", nil, nil
}

func resourceIPGeoCreate(ctx context.Context, d *schema.ResourceData, m interface{}) diag.Diagnostics {
	meta := meta.Must(m)
	client := inst.Client(meta)
	logger := meta.Log("APPSEC", "resourceIPGeoCreate")
	logger.Debugf("in resourceIPGeoCreate")

	configID, err := tf.GetIntValue("config_id", d)
	if err != nil {
		return diag.FromErr(err)
	}
	version, err := getModifiableConfigVersion(ctx, configID, "ipgeo", m)
	if err != nil {
		return diag.FromErr(err)
	}
	policyID, err := tf.GetStringValue("security_policy_id", d)
	if err != nil {
		return diag.FromErr(err)
	}
	mode, err := tf.GetStringValue("mode", d)
	if err != nil {
		return diag.FromErr(err)
	}

	asnBlockAction, blockedASNLists, diags := extractAsnControls(d)
	if diags != nil {
		return diags
	}

	geoBlockAction, blockedGeoLists, diags := extractGeoControls(d)
	if diags != nil {
		return diags
	}

	ipBlockAction, blockedIPLists, diags := extractIPControls(d)
	if diags != nil {
		return diags
	}

	exceptionIPLists, err := tf.GetSetAsListValue("exception_ip_network_lists", d)
	if err != nil && !errors.Is(err, tf.ErrNotFound) {
		return diag.FromErr(err)
	}

	ukraineGeoControlAction, err := tf.GetStringValue("ukraine_geo_control_action", d)
	if err != nil && !errors.Is(err, tf.ErrNotFound) {
		return diag.FromErr(err)
	}
	blockAction, err := tf.GetStringValue("block_action", d)
	if err != nil && !errors.Is(err, tf.ErrNotFound) {
		return diag.FromErr(err)
	}

	request := appsec.UpdateIPGeoRequest{
		ConfigID: configID,
		Version:  version,
		PolicyID: policyID,
	}

	if mode == Allow {
		request.Block = "blockAllTrafficExceptAllowedIPs"
		if blockAction == "" {
			request.BlockAllAction = "deny" // fallback to default manually
		} else {
			request.BlockAllAction = blockAction
		}
		request.IPControls = ipControlsFromAllowLists(exceptionIPLists)
	}
	if mode == Block {
		request.Block = "blockSpecificIPGeo"
		request.GeoControls = geoControlsFromBlockLists(blockedGeoLists, geoBlockAction)
		request.ASNControls = asnControlsFromBlockLists(blockedASNLists, asnBlockAction)
		request.IPControls = ipControlsFromBlockAndAllowLists(blockedIPLists, ipBlockAction, exceptionIPLists)
		if ukraineGeoControlAction != "" {
			request.UkraineGeoControls = &appsec.UkraineGeoControl{Action: ukraineGeoControlAction}
		}
	}

	_, err = client.UpdateIPGeo(ctx, request)
	if err != nil {
		logger.Errorf("calling 'createIPGeo': %s", err.Error())
		return diag.FromErr(err)
	}

	d.SetId(fmt.Sprintf("%d:%s", configID, policyID))

	return resourceIPGeoRead(ctx, d, m)
}

func resourceIPGeoRead(ctx context.Context, d *schema.ResourceData, m interface{}) diag.Diagnostics {
	meta := meta.Must(m)
	client := inst.Client(meta)
	logger := meta.Log("APPSEC", "resourceIPGeoRead")
	logger.Debugf("in resourceIPGeoRead")

	iDParts, err := id.Split(d.Id(), 2, "configID:securityPolicyID")
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
		return diag.Errorf("%s: %s", tf.ErrValueSet, err.Error())
	}
	if err := d.Set("security_policy_id", getIPGeo.PolicyID); err != nil {
		return diag.Errorf("%s: %s", tf.ErrValueSet, err.Error())
	}
	if ipgeo.Block == "blockAllTrafficExceptAllowedIPs" {
		if err := d.Set("mode", Allow); err != nil {
			return diag.Errorf("%s: %s", tf.ErrValueSet, err.Error())
		}
		if ipgeo.IPControls != nil && ipgeo.IPControls.AllowedIPNetworkLists != nil && ipgeo.IPControls.AllowedIPNetworkLists.NetworkList != nil {
			if err := d.Set("exception_ip_network_lists", ipgeo.IPControls.AllowedIPNetworkLists.NetworkList); err != nil {
				return diag.Errorf("%s: %s", tf.ErrValueSet, err.Error())
			}
		} else {
			if err := d.Set("exception_ip_network_lists", make([]string, 0)); err != nil {
				return diag.Errorf("%s: %s", tf.ErrValueSet, err.Error())
			}
		}
		if err := d.Set("block_action", ipgeo.BlockAllAction); err != nil {
			return diag.Errorf("%s: %s", tf.ErrValueSet, err.Error())
		}
	}
	if ipgeo.Block == "blockSpecificIPGeo" {
		diagnostics := readForBlockSpecificIPGeo(d, ipgeo)
		if diagnostics != nil {
			return diagnostics
		}
	}

	return nil
}

func readForBlockSpecificIPGeo(d *schema.ResourceData, ipgeo *appsec.GetIPGeoResponse) diag.Diagnostics {
	if err := d.Set("mode", Block); err != nil {
		return diag.Errorf("%s: %s", tf.ErrValueSet, err.Error())
	}
	asnControls := map[string]interface{}{
		"asn_network_lists": []string{},
		"action":            "",
	}

	if ipgeo.ASNControls != nil {
		if ipgeo.ASNControls.BlockedIPNetworkLists != nil && ipgeo.ASNControls.BlockedIPNetworkLists.NetworkList != nil {
			asnControls["asn_network_lists"] = ipgeo.ASNControls.BlockedIPNetworkLists.NetworkList
			asnControls["action"] = ipgeo.ASNControls.BlockedIPNetworkLists.Action
		}
	}
	if err := d.Set("asn_controls", []interface{}{asnControls}); err != nil {
		return diag.Errorf("%s: %s", tf.ErrValueSet, err.Error())
	}

	geoControls := map[string]interface{}{
		"geo_network_lists": []string{},
		"action":            "",
	}
	if ipgeo.GeoControls != nil {
		if ipgeo.GeoControls.BlockedIPNetworkLists != nil && ipgeo.GeoControls.BlockedIPNetworkLists.NetworkList != nil {
			geoControls["geo_network_lists"] = ipgeo.GeoControls.BlockedIPNetworkLists.NetworkList
			geoControls["action"] = ipgeo.GeoControls.BlockedIPNetworkLists.Action
		}
	}
	if err := d.Set("geo_controls", []interface{}{geoControls}); err != nil {
		return diag.Errorf("%s: %s", tf.ErrValueSet, err.Error())
	}

	ipControls := map[string]interface{}{
		"ip_network_lists": []string{},
		"action":           "",
	}
	if ipgeo.IPControls != nil {
		if ipgeo.IPControls.BlockedIPNetworkLists != nil && ipgeo.IPControls.BlockedIPNetworkLists.NetworkList != nil {
			ipControls["ip_network_lists"] = ipgeo.IPControls.BlockedIPNetworkLists.NetworkList
			ipControls["action"] = ipgeo.IPControls.BlockedIPNetworkLists.Action
		}
	}
	if err := d.Set("ip_controls", []interface{}{ipControls}); err != nil {
		return diag.Errorf("%s: %s", tf.ErrValueSet, err.Error())
	}

	if ipgeo.IPControls != nil && ipgeo.IPControls.AllowedIPNetworkLists != nil && ipgeo.IPControls.AllowedIPNetworkLists.NetworkList != nil {
		if err := d.Set("exception_ip_network_lists", ipgeo.IPControls.AllowedIPNetworkLists.NetworkList); err != nil {
			return diag.Errorf("%s: %s", tf.ErrValueSet, err.Error())
		}
	} else {
		if err := d.Set("exception_ip_network_lists", make([]string, 0)); err != nil {
			return diag.Errorf("%s: %s", tf.ErrValueSet, err.Error())
		}
	}
	if ipgeo.UkraineGeoControls != nil && ipgeo.UkraineGeoControls.Action != "" {
		if err := d.Set("ukraine_geo_control_action", ipgeo.UkraineGeoControls.Action); err != nil {
			return diag.Errorf("%s: %s", tf.ErrValueSet, err.Error())
		}
	}
	return nil
}

func resourceIPGeoUpdate(ctx context.Context, d *schema.ResourceData, m interface{}) diag.Diagnostics {
	meta := meta.Must(m)
	client := inst.Client(meta)
	logger := meta.Log("APPSEC", "resourceIPGeoUpdate")
	logger.Debugf("in resourceIPGeoUpdate")

	iDParts, err := id.Split(d.Id(), 2, "configID:securityPolicyID")
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
	mode, err := tf.GetStringValue("mode", d)
	if err != nil && !errors.Is(err, tf.ErrNotFound) {
		return diag.FromErr(err)
	}

	asnBlockAction, blockedASNLists, diags := extractAsnControls(d)
	if diags != nil {
		return diags
	}

	geoBlockAction, blockedGeoLists, diags := extractGeoControls(d)
	if diags != nil {
		return diags
	}

	ipBlockAction, blockedIPLists, diags := extractIPControls(d)
	if diags != nil {
		return diags
	}

	exceptionIPLists, err := tf.GetSetAsListValue("exception_ip_network_lists", d)
	if err != nil && !errors.Is(err, tf.ErrNotFound) {
		return diag.FromErr(err)
	}

	blockAction, err := tf.GetStringValue("block_action", d)
	if err != nil && !errors.Is(err, tf.ErrNotFound) {
		return diag.FromErr(err)
	}

	ukraineGeoControlAction, err := tf.GetStringValue("ukraine_geo_control_action", d)
	if err != nil && !errors.Is(err, tf.ErrNotFound) {
		return diag.FromErr(err)
	}

	request := appsec.UpdateIPGeoRequest{
		ConfigID: configID,
		Version:  version,
		PolicyID: policyID,
	}
	if mode == Allow {
		request.Block = "blockAllTrafficExceptAllowedIPs"
		if blockAction == "" {
			request.BlockAllAction = "deny" // fallback to default manually
		} else {
			request.BlockAllAction = blockAction
		}
		request.IPControls = ipControlsFromAllowLists(exceptionIPLists)
	}
	if mode == Block {
		request.Block = "blockSpecificIPGeo"
		request.ASNControls = asnControlsFromBlockLists(blockedASNLists, asnBlockAction)
		request.GeoControls = geoControlsFromBlockLists(blockedGeoLists, geoBlockAction)
		request.IPControls = ipControlsFromBlockAndAllowLists(blockedIPLists, ipBlockAction, exceptionIPLists)
		if ukraineGeoControlAction != "" {
			request.UkraineGeoControls = &appsec.UkraineGeoControl{Action: ukraineGeoControlAction}
		}
	}

	_, err = client.UpdateIPGeo(ctx, request)
	if err != nil {
		logger.Errorf("calling 'updateIPGeo': %s", err.Error())
		return diag.FromErr(err)
	}

	return resourceIPGeoRead(ctx, d, m)
}

func resourceIPGeoDelete(ctx context.Context, d *schema.ResourceData, m interface{}) diag.Diagnostics {
	meta := meta.Must(m)
	client := inst.Client(meta)
	logger := meta.Log("APPSEC", "resourceIPGeoDelete")
	logger.Debugf("in resourceIPGeoDelete")

	iDParts, err := id.Split(d.Id(), 2, "configID:securityPolicyID")
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
