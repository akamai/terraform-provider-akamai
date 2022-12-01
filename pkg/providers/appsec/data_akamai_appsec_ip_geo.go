package appsec

import (
	"context"
	"strconv"

	"github.com/akamai/AkamaiOPEN-edgegrid-golang/v3/pkg/appsec"
	"github.com/akamai/terraform-provider-akamai/v3/pkg/akamai"
	"github.com/akamai/terraform-provider-akamai/v3/pkg/tools"
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
			"geo_network_lists": {
				Type:        schema.TypeSet,
				Computed:    true,
				Elem:        &schema.Schema{Type: schema.TypeString},
				Description: "List of unique identifiers of available GEO network lists",
			},
			"ip_network_lists": {
				Type:        schema.TypeSet,
				Computed:    true,
				Elem:        &schema.Schema{Type: schema.TypeString},
				Description: "List of unique identifiers of available IP network lists",
			},
			"exception_ip_network_lists": {
				Type:        schema.TypeSet,
				Computed:    true,
				Elem:        &schema.Schema{Type: schema.TypeString},
				Description: "List of unique identifiers of network lists allowed through the firewall regardless of mode, geo_network_lists and ip_network_lists values",
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
	meta := akamai.Meta(m)
	client := inst.Client(meta)
	logger := meta.Log("APPSEC", "dataSourceIPGeoRead")

	getIPGeo := appsec.GetIPGeoRequest{}

	configID, err := tools.GetIntValue("config_id", d)
	if err != nil {
		return diag.FromErr(err)
	}
	getIPGeo.ConfigID = configID

	if getIPGeo.Version, err = getLatestConfigVersion(ctx, configID, m); err != nil {
		return diag.FromErr(err)
	}

	policyID, err := tools.GetStringValue("security_policy_id", d)
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
		return diag.Errorf("%s: %s", tools.ErrValueSet, err.Error())
	}

	if ipgeo.Block == "blockAllTrafficExceptAllowedIPs" {
		if err := d.Set("mode", Allow); err != nil {
			return diag.Errorf("%s: %s", tools.ErrValueSet, err.Error())
		}
		if ipgeo.IPControls != nil && ipgeo.IPControls.AllowedIPNetworkLists != nil && ipgeo.IPControls.AllowedIPNetworkLists.NetworkList != nil {
			if err := d.Set("exception_ip_network_lists", ipgeo.IPControls.AllowedIPNetworkLists.NetworkList); err != nil {
				return diag.Errorf("%s: %s", tools.ErrValueSet, err.Error())
			}
		} else {
			if err := d.Set("exception_ip_network_lists", []string{}); err != nil {
				return diag.Errorf("%s: %s", tools.ErrValueSet, err.Error())
			}
		}
	}
	if ipgeo.Block == "blockSpecificIPGeo" {
		if err := d.Set("mode", Block); err != nil {
			return diag.Errorf("%s: %s", tools.ErrValueSet, err.Error())
		}
		if ipgeo.GeoControls != nil && ipgeo.GeoControls.BlockedIPNetworkLists != nil && ipgeo.GeoControls.BlockedIPNetworkLists.NetworkList != nil {
			if err := d.Set("geo_network_lists", ipgeo.GeoControls.BlockedIPNetworkLists.NetworkList); err != nil {
				return diag.Errorf("%s: %s", tools.ErrValueSet, err.Error())
			}
		} else {
			if err := d.Set("geo_network_lists", []string{}); err != nil {
				return diag.Errorf("%s: %s", tools.ErrValueSet, err.Error())
			}
		}
		if ipgeo.IPControls != nil && ipgeo.IPControls.BlockedIPNetworkLists != nil && ipgeo.IPControls.BlockedIPNetworkLists.NetworkList != nil {
			if err := d.Set("ip_network_lists", ipgeo.IPControls.BlockedIPNetworkLists.NetworkList); err != nil {
				return diag.Errorf("%s: %s", tools.ErrValueSet, err.Error())
			}
		} else {
			if err := d.Set("ip_network_lists", []string{}); err != nil {
				return diag.Errorf("%s: %s", tools.ErrValueSet, err.Error())
			}
		}
	}

	d.SetId(strconv.Itoa(getIPGeo.ConfigID))

	return nil
}
