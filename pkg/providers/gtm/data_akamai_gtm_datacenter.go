package gtm

import (
	"context"
	"fmt"

	"github.com/akamai/AkamaiOPEN-edgegrid-golang/v10/pkg/gtm"
	"github.com/akamai/AkamaiOPEN-edgegrid-golang/v10/pkg/session"
	"github.com/akamai/terraform-provider-akamai/v6/pkg/common/tf"
	"github.com/akamai/terraform-provider-akamai/v6/pkg/meta"
	"github.com/hashicorp/terraform-plugin-sdk/v2/diag"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/schema"
)

func dataSourceGTMDatacenter() *schema.Resource {
	return &schema.Resource{
		ReadContext: dataGTMDatacenterRead,
		Schema: map[string]*schema.Schema{
			"domain": {
				Type:        schema.TypeString,
				Required:    true,
				Description: "GTM domain name.",
			},
			"datacenter_id": {
				Type:        schema.TypeInt,
				Required:    true,
				Description: "A unique identifier for an existing data center in the domain.",
			},
			"nickname": {
				Type:        schema.TypeString,
				Computed:    true,
				Description: "A descriptive label for the datacenter.",
			},
			"score_penalty": {
				Type:        schema.TypeInt,
				Computed:    true,
				Description: "Influences the score for a datacenter.",
			},
			"city": {
				Type:        schema.TypeString,
				Computed:    true,
				Description: "The name of the city where the data center is located.",
			},
			"state_or_province": {
				Type:        schema.TypeString,
				Computed:    true,
				Description: "Specifies a two-letter ISO 3166 country code for the state of province, where the data center is located.",
			},
			"country": {
				Type:        schema.TypeString,
				Computed:    true,
				Description: "A two-letter ISO 3166 contry code that specifies the country where the data center is located.",
			},
			"latitude": {
				Type:        schema.TypeFloat,
				Computed:    true,
				Description: "Specifies the geographic latitude of the data center's position.",
			},
			"longitude": {
				Type:        schema.TypeFloat,
				Computed:    true,
				Description: "Specifies the geographic longitude of the data center's position.",
			},
			"clone_of": {
				Type:        schema.TypeInt,
				Computed:    true,
				Description: "Identifies the data center's ID of which this data center is a clone.",
			},
			"virtual": {
				Type:        schema.TypeBool,
				Computed:    true,
				Description: "Indicates whether or not the data center is virtual or physical.",
			},
			"default_load_object": {
				Type:        schema.TypeList,
				Computed:    true,
				Description: "Specifies the load reporting interface between you and the GTM system.",
				Elem: &schema.Resource{
					Schema: map[string]*schema.Schema{
						"load_object": {
							Type:        schema.TypeString,
							Computed:    true,
							Description: "Specifies the load object that GTM requests.",
						},
						"load_object_port": {
							Type:        schema.TypeInt,
							Computed:    true,
							Description: "Specifies the TCP port to connect to when requesting the load object.",
						},
						"load_servers": {
							Type:        schema.TypeSet,
							Computed:    true,
							Description: "Specifies the list of servers to requests the load object from.",
							Elem: &schema.Schema{
								Type: schema.TypeString,
							},
						},
					},
				},
			},
			"continent": {
				Type:        schema.TypeString,
				Computed:    true,
				Description: "A two-letter code that specifies the continent where the data center maps to.",
			},
			"servermonitor_pool": {
				Type:        schema.TypeString,
				Computed:    true,
				Description: "The name of the pool from which servermonitors are drawn for liveness tests in this datacenter. If omitted (null), the domain-wide default is used. (If no domain-wide default is specified, the pool used is all servermonitors in the same continent as the datacenter.)",
			},
			"cloud_server_targeting": {
				Type:        schema.TypeBool,
				Computed:    true,
				Description: "Balances load between two or more servers in a cloud environment.",
			},
			"cloud_server_host_header_override": {
				Type:        schema.TypeBool,
				Computed:    true,
				Description: "Balances load between two or more servers in a cloud environment.",
			},
			"links": {
				Type:        schema.TypeList,
				Computed:    true,
				Description: "Provides a URL path that allows direct navigation to a data center.",
				Elem: &schema.Resource{
					Schema: map[string]*schema.Schema{
						"rel": {
							Type:        schema.TypeString,
							Computed:    true,
							Description: "Indicates the link relationship of the object.",
						},
						"href": {
							Type:        schema.TypeString,
							Computed:    true,
							Description: "A hypermedia link to the complete URL that uniquely defines a resource.",
						},
					},
				},
			},
		},
	}
}

func dataGTMDatacenterRead(ctx context.Context, d *schema.ResourceData, m interface{}) diag.Diagnostics {
	meta := meta.Must(m)
	logger := meta.Log("Akamai GTM", "dataGTMDatacenterRead")
	ctx = session.ContextWithOptions(ctx, session.WithContextLog(logger))
	client := Client(meta)

	logger.Debug("Fetching a datacenter")

	domain, err := tf.GetStringValue("domain", d)
	if err != nil {
		return diag.Errorf("could not get 'domain' attribute: %s", err)
	}

	dcID, err := tf.GetIntValue("datacenter_id", d)
	if err != nil {
		return diag.Errorf("could not get 'datacenter_id' attribute: %s", err)
	}

	dc, err := client.GetDatacenter(ctx, gtm.GetDatacenterRequest{
		DatacenterID: dcID,
		DomainName:   domain,
	})
	if err != nil {
		return diag.Errorf("could not get datacenter: %s", err)
	}

	if err = tf.SetAttrs(d, getDatacenterAttributes(dc)); err != nil {
		return diag.Errorf("could not set attributes: %s", err)
	}

	d.SetId(fmt.Sprintf("%s:%d", domain, dcID))

	return nil
}

func getDatacenterAttributes(dc *gtm.Datacenter) map[string]interface{} {
	var defaultLoadObject []map[string]interface{}
	if dc.DefaultLoadObject != nil {
		defaultLoadObject = []map[string]interface{}{
			{
				"load_object":      dc.DefaultLoadObject.LoadObject,
				"load_object_port": dc.DefaultLoadObject.LoadObjectPort,
				"load_servers":     dc.DefaultLoadObject.LoadServers,
			},
		}
	}

	links := make([]map[string]string, len(dc.Links))
	for i, link := range dc.Links {
		links[i] = map[string]string{
			"rel":  link.Rel,
			"href": link.Href,
		}
	}

	attrs := map[string]interface{}{
		"nickname":                          dc.Nickname,
		"score_penalty":                     dc.ScorePenalty,
		"city":                              dc.City,
		"state_or_province":                 dc.StateOrProvince,
		"country":                           dc.Country,
		"latitude":                          dc.Latitude,
		"longitude":                         dc.Longitude,
		"clone_of":                          dc.CloneOf,
		"virtual":                           dc.Virtual,
		"default_load_object":               defaultLoadObject,
		"continent":                         dc.Continent,
		"servermonitor_pool":                dc.ServermonitorPool,
		"cloud_server_targeting":            dc.CloudServerTargeting,
		"cloud_server_host_header_override": dc.CloudServerHostHeaderOverride,
		"links":                             links,
	}

	return attrs
}
