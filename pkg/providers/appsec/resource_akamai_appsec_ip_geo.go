package appsec

import (
	"context"
	"errors"
	"fmt"
	"strconv"

	"github.com/akamai/AkamaiOPEN-edgegrid-golang/v8/pkg/appsec"
	"github.com/akamai/terraform-provider-akamai/v6/pkg/common/tf"
	"github.com/akamai/terraform-provider-akamai/v6/pkg/meta"
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
				Type:     schema.TypeList,
				Optional: true,
				Elem: &schema.Schema{
					Type:             schema.TypeString,
					ValidateDiagFunc: validateEmptyElementsInList,
				},
				Description: "List of IDs of geographic network list to be blocked",
			},
			"ip_network_lists": {
				Type:     schema.TypeList,
				Optional: true,
				Elem: &schema.Schema{
					Type:             schema.TypeString,
					ValidateDiagFunc: validateEmptyElementsInList,
				},
				Description: "List of IDs of IP network list to be blocked",
			},
			"asn_network_lists": {
				Type:     schema.TypeList,
				Optional: true,
				Elem: &schema.Schema{
					Type:             schema.TypeString,
					ValidateDiagFunc: validateEmptyElementsInList,
				},
				Description: "List of IDs of ASN network list to be blocked",
			},
			"exception_ip_network_lists": {
				Type:     schema.TypeList,
				Optional: true,
				Elem: &schema.Schema{
					Type:             schema.TypeString,
					ValidateDiagFunc: validateEmptyElementsInList,
				},
				Description: "List of IDs of network list that are always allowed",
			},
			"ukraine_geo_control_action": {
				Type:             schema.TypeString,
				Optional:         true,
				DiffSuppressFunc: suppressDiffUkraineGeoControlAction,
				Description:      "Action set for Ukraine geo control",
			},
		},
	}
}

func ipGeoNetworkListsFromStringList(strings []interface{}) *appsec.IPGeoNetworkLists {
	items := make([]string, 0, len(strings))
	for _, item := range strings {
		items = append(items, item.(string))
	}
	return &appsec.IPGeoNetworkLists{
		NetworkList: items,
	}
}

func ipControlsFromAllowLists(exceptionIPLists []interface{}) *appsec.IPGeoIPControls {
	if len(exceptionIPLists) > 0 {
		return &appsec.IPGeoIPControls{
			AllowedIPNetworkLists: ipGeoNetworkListsFromStringList(exceptionIPLists),
		}
	}
	return nil
}

func geoControlsFromBlockLists(blockedGeoLists []interface{}) *appsec.IPGeoGeoControls {
	if len(blockedGeoLists) > 0 {
		return &appsec.IPGeoGeoControls{
			BlockedIPNetworkLists: ipGeoNetworkListsFromStringList(blockedGeoLists),
		}
	}
	return nil
}

func asnControlsFromBlockLists(blockedASNLists []interface{}) *appsec.IPGeoASNControls {
	if len(blockedASNLists) > 0 {
		return &appsec.IPGeoASNControls{
			BlockedIPNetworkLists: ipGeoNetworkListsFromStringList(blockedASNLists),
		}
	}
	return nil
}

func ipControlsFromBlockAndAllowLists(blockedIPLists []interface{}, exceptionIPLists []interface{}) *appsec.IPGeoIPControls {
	if len(blockedIPLists) > 0 || len(exceptionIPLists) > 0 {
		ipControls := &appsec.IPGeoIPControls{}
		if len(blockedIPLists) > 0 {
			ipControls.BlockedIPNetworkLists = ipGeoNetworkListsFromStringList(blockedIPLists)
		}
		if len(exceptionIPLists) > 0 {
			ipControls.AllowedIPNetworkLists = ipGeoNetworkListsFromStringList(exceptionIPLists)
		}
		return ipControls
	}
	return nil
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
	blockedGeoLists, err := tf.GetListValue("geo_network_lists", d)
	if err != nil && !errors.Is(err, tf.ErrNotFound) {
		return diag.FromErr(err)
	}
	blockedIPLists, err := tf.GetListValue("ip_network_lists", d)
	if err != nil && !errors.Is(err, tf.ErrNotFound) {
		return diag.FromErr(err)
	}
	blockedASNLists, err := tf.GetListValue("asn_network_lists", d)
	if err != nil && !errors.Is(err, tf.ErrNotFound) {
		return diag.FromErr(err)
	}
	exceptionIPLists, err := tf.GetListValue("exception_ip_network_lists", d)
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
		request.IPControls = ipControlsFromAllowLists(exceptionIPLists)
	}
	if mode == Block {
		request.Block = "blockSpecificIPGeo"
		request.GeoControls = geoControlsFromBlockLists(blockedGeoLists)
		request.ASNControls = asnControlsFromBlockLists(blockedASNLists)
		request.IPControls = ipControlsFromBlockAndAllowLists(blockedIPLists, exceptionIPLists)
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
	if ipgeo.GeoControls != nil && ipgeo.GeoControls.BlockedIPNetworkLists != nil && ipgeo.GeoControls.BlockedIPNetworkLists.NetworkList != nil {
		if err := d.Set("geo_network_lists", ipgeo.GeoControls.BlockedIPNetworkLists.NetworkList); err != nil {
			return diag.Errorf("%s: %s", tf.ErrValueSet, err.Error())
		}
	} else {
		if err := d.Set("geo_network_lists", make([]string, 0)); err != nil {
			return diag.Errorf("%s: %s", tf.ErrValueSet, err.Error())
		}
	}
	if ipgeo.IPControls != nil && ipgeo.IPControls.BlockedIPNetworkLists != nil && ipgeo.IPControls.BlockedIPNetworkLists.NetworkList != nil {
		if err := d.Set("ip_network_lists", ipgeo.IPControls.BlockedIPNetworkLists.NetworkList); err != nil {
			return diag.Errorf("%s: %s", tf.ErrValueSet, err.Error())
		}
	} else {
		if err := d.Set("ip_network_lists", make([]string, 0)); err != nil {
			return diag.Errorf("%s: %s", tf.ErrValueSet, err.Error())
		}
	}
	if ipgeo.ASNControls != nil && ipgeo.ASNControls.BlockedIPNetworkLists != nil && ipgeo.ASNControls.BlockedIPNetworkLists.NetworkList != nil {
		if err := d.Set("asn_network_lists", ipgeo.ASNControls.BlockedIPNetworkLists.NetworkList); err != nil {
			return diag.Errorf("%s: %s", tf.ErrValueSet, err.Error())
		}
	} else {
		if err := d.Set("asn_network_lists", make([]string, 0)); err != nil {
			return diag.Errorf("%s: %s", tf.ErrValueSet, err.Error())
		}
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
	mode, err := tf.GetStringValue("mode", d)
	if err != nil && !errors.Is(err, tf.ErrNotFound) {
		return diag.FromErr(err)
	}
	blockedASNLists, err := tf.GetListValue("asn_network_lists", d)
	if err != nil && !errors.Is(err, tf.ErrNotFound) {
		return diag.FromErr(err)
	}
	blockedGeoLists, err := tf.GetListValue("geo_network_lists", d)
	if err != nil && !errors.Is(err, tf.ErrNotFound) {
		return diag.FromErr(err)
	}
	blockedIPLists, err := tf.GetListValue("ip_network_lists", d)
	if err != nil && !errors.Is(err, tf.ErrNotFound) {
		return diag.FromErr(err)
	}
	exceptionIPLists, err := tf.GetListValue("exception_ip_network_lists", d)
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
		request.IPControls = ipControlsFromAllowLists(exceptionIPLists)
	}
	if mode == Block {
		request.Block = "blockSpecificIPGeo"
		request.ASNControls = asnControlsFromBlockLists(blockedASNLists)
		request.GeoControls = geoControlsFromBlockLists(blockedGeoLists)
		request.IPControls = ipControlsFromBlockAndAllowLists(blockedIPLists, exceptionIPLists)
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
