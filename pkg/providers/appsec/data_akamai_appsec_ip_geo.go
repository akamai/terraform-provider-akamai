package appsec

import (
	"context"
	"strconv"

	"github.com/akamai/AkamaiOPEN-edgegrid-golang/v12/pkg/appsec"
	"github.com/akamai/terraform-provider-akamai/v9/pkg/common/tf"
	"github.com/akamai/terraform-provider-akamai/v9/pkg/meta"
	"github.com/hashicorp/terraform-plugin-sdk/v2/diag"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/schema"
)

func dataSourceIPGeo() *schema.Resource {
	return &schema.Resource{
		ReadContext: dataSourceIPGeoRead,
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
				Type:        schema.TypeString,
				Computed:    true,
				Description: "IPGeo mode",
			},
			"asn_controls": {
				Type:     schema.TypeList,
				Computed: true,
				Elem: &schema.Resource{
					Schema: map[string]*schema.Schema{
						"asn_network_lists": {
							Type:        schema.TypeSet,
							Elem:        &schema.Schema{Type: schema.TypeString},
							Computed:    true,
							Description: "List of IDs of ASN network list to be blocked.",
						},
						"action": {
							Type:        schema.TypeString,
							Computed:    true,
							Description: "Action set for ASN Controls",
						},
					},
				},
			},
			"geo_controls": {
				Type:     schema.TypeList,
				Computed: true,
				Elem: &schema.Resource{
					Schema: map[string]*schema.Schema{
						"geo_network_lists": {
							Type:        schema.TypeSet,
							Elem:        &schema.Schema{Type: schema.TypeString},
							Computed:    true,
							Description: "List of IDs of geographic network list to be blocked.",
						},
						"action": {
							Type:        schema.TypeString,
							Computed:    true,
							Description: "Action set for GEO Controls.",
						},
					},
				},
			},
			"ip_controls": {
				Type:     schema.TypeList,
				Computed: true,
				Elem: &schema.Resource{
					Schema: map[string]*schema.Schema{
						"ip_network_lists": {
							Type:        schema.TypeSet,
							Elem:        &schema.Schema{Type: schema.TypeString},
							Computed:    true,
							Description: "List of IDs of IP network list to be blocked.",
						},
						"action": {
							Type:        schema.TypeString,
							Computed:    true,
							Description: "Action set for IP Controls.",
						},
					},
				},
			},
			"exception_ip_network_lists": {
				Type:        schema.TypeSet,
				Computed:    true,
				Elem:        &schema.Schema{Type: schema.TypeString},
				Description: "List of unique identifiers of network lists allowed through the firewall regardless of mode, asn_network_lists, geo_network_lists and ip_network_lists values",
			},
			"ukraine_geo_control_action": {
				Type:        schema.TypeString,
				Computed:    true,
				Description: "Action set for Ukraine geo control",
			},
			"block_action": {
				Type:        schema.TypeString,
				Computed:    true,
				Description: "Specifies the action set for BLOCK Mode blocking all the traffic except from lists identified in exception_ip_network_lists",
			},
			"output_text": {
				Type:        schema.TypeString,
				Computed:    true,
				Description: "Text representation",
			},
		},
	}
}

func dataSourceIPGeoRead(ctx context.Context, d *schema.ResourceData, m interface{}) diag.Diagnostics {
	meta := meta.Must(m)
	client := inst.Client(meta)
	logger := meta.Log("APPSEC", "dataSourceIPGeoRead")

	getIPGeo := appsec.GetIPGeoRequest{}

	configID, err := tf.GetIntValue("config_id", d)
	if err != nil {
		return diag.FromErr(err)
	}
	getIPGeo.ConfigID = configID

	if getIPGeo.Version, err = getLatestConfigVersion(ctx, configID, m); err != nil {
		return diag.FromErr(err)
	}

	policyID, err := tf.GetStringValue("security_policy_id", d)
	if err != nil {
		return diag.FromErr(err)
	}
	getIPGeo.PolicyID = policyID

	ipgeo, err := client.GetIPGeo(ctx, getIPGeo)
	if err != nil {
		logger.Errorf("calling 'getIPGeo': %s", err.Error())
		return diag.FromErr(err)
	}

	ots := OutputTemplates{}
	InitTemplates(ots)

	outputtext, err := RenderTemplates(ots, "IPGeoDS", ipgeo)
	if err != nil {
		return diag.FromErr(err)
	}
	if err := d.Set("output_text", outputtext); err != nil {
		return diag.Errorf("%s: %s", tf.ErrValueSet, err.Error())
	}

	res := setAttributes(d, ipgeo)
	if res != nil {
		return res
	}

	d.SetId(strconv.Itoa(getIPGeo.ConfigID))

	return nil
}

func setAttributes(d *schema.ResourceData, ipgeo *appsec.GetIPGeoResponse) diag.Diagnostics {
	if ipgeo.Block == "blockAllTrafficExceptAllowedIPs" {
		if err := d.Set("mode", Allow); err != nil {
			return diag.Errorf("%s: %s", tf.ErrValueSet, err.Error())
		}

		if err := d.Set("block_action", ipgeo.BlockAllAction); err != nil {
			return diag.Errorf("%s: %s", tf.ErrValueSet, err.Error())
		}
	}

	if ipgeo.Block == "blockSpecificIPGeo" {
		if err := d.Set("mode", Block); err != nil {
			return diag.Errorf("%s: %s", tf.ErrValueSet, err.Error())
		}

		if ipgeo.ASNControls != nil && ipgeo.ASNControls.BlockedIPNetworkLists != nil && ipgeo.ASNControls.BlockedIPNetworkLists.NetworkList != nil {
			asnControls := map[string]interface{}{
				"asn_network_lists": nil,
				"action":            nil,
			}
			asnControls["asn_network_lists"] = ipgeo.ASNControls.BlockedIPNetworkLists.NetworkList
			asnControls["action"] = ipgeo.ASNControls.BlockedIPNetworkLists.Action
			if err := d.Set("asn_controls", []interface{}{asnControls}); err != nil {
				return diag.Errorf("%s: %s", tf.ErrValueSet, err.Error())
			}
		} else {
			// clear the field if the API returned nil.
			if err := d.Set("asn_controls", nil); err != nil {
				return diag.Errorf("%s: %s", tf.ErrValueSet, err.Error())
			}
		}

		if ipgeo.GeoControls != nil && ipgeo.GeoControls.BlockedIPNetworkLists != nil && ipgeo.GeoControls.BlockedIPNetworkLists.NetworkList != nil {
			geoControls := map[string]interface{}{
				"geo_network_lists": nil,
				"action":            nil,
			}
			geoControls["geo_network_lists"] = ipgeo.GeoControls.BlockedIPNetworkLists.NetworkList
			geoControls["action"] = ipgeo.GeoControls.BlockedIPNetworkLists.Action
			if err := d.Set("geo_controls", []interface{}{geoControls}); err != nil {
				return diag.Errorf("%s: %s", tf.ErrValueSet, err.Error())
			}
		} else {
			// clear the field if the API returned nil.
			if err := d.Set("geo_controls", nil); err != nil {
				return diag.Errorf("%s: %s", tf.ErrValueSet, err.Error())
			}
		}

		if ipgeo.IPControls != nil && ipgeo.IPControls.BlockedIPNetworkLists != nil && ipgeo.IPControls.BlockedIPNetworkLists.NetworkList != nil {
			ipControls := map[string]interface{}{
				"ip_network_lists": nil,
				"action":           nil,
			}
			ipControls["ip_network_lists"] = ipgeo.IPControls.BlockedIPNetworkLists.NetworkList
			ipControls["action"] = ipgeo.IPControls.BlockedIPNetworkLists.Action
			if err := d.Set("ip_controls", []interface{}{ipControls}); err != nil {
				return diag.Errorf("%s: %s", tf.ErrValueSet, err.Error())
			}
		} else {
			// clear the field if the API returned nil.
			if err := d.Set("ip_controls", nil); err != nil {
				return diag.Errorf("%s: %s", tf.ErrValueSet, err.Error())
			}
		}

		if ipgeo.UkraineGeoControls != nil && ipgeo.UkraineGeoControls.Action != "" {
			if err := d.Set("ukraine_geo_control_action", ipgeo.UkraineGeoControls.Action); err != nil {
				return diag.Errorf("%s: %s", tf.ErrValueSet, err.Error())
			}
		}
	}

	if ipgeo.IPControls != nil && ipgeo.IPControls.AllowedIPNetworkLists != nil && ipgeo.IPControls.AllowedIPNetworkLists.NetworkList != nil {
		if err := d.Set("exception_ip_network_lists", ipgeo.IPControls.AllowedIPNetworkLists.NetworkList); err != nil {
			return diag.Errorf("%s: %s", tf.ErrValueSet, err.Error())
		}
	} else {
		if err := d.Set("exception_ip_network_lists", []string{}); err != nil {
			return diag.Errorf("%s: %s", tf.ErrValueSet, err.Error())
		}
	}

	return nil
}
