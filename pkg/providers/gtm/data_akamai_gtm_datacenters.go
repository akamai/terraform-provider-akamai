package gtm

import (
	"context"

	"github.com/akamai/AkamaiOPEN-edgegrid-golang/v12/pkg/gtm"
	"github.com/akamai/AkamaiOPEN-edgegrid-golang/v12/pkg/session"
	"github.com/akamai/terraform-provider-akamai/v8/pkg/common/tf"
	"github.com/akamai/terraform-provider-akamai/v8/pkg/meta"
	"github.com/hashicorp/terraform-plugin-sdk/v2/diag"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/schema"
)

func dataSourceGTMDatacenters() *schema.Resource {
	return &schema.Resource{
		ReadContext: dataGTMDatacentersRead,
		Schema: map[string]*schema.Schema{
			"domain": {
				Type:        schema.TypeString,
				Required:    true,
				Description: "GTM domain name",
			},
			"datacenters": {
				Type:        schema.TypeList,
				Computed:    true,
				Description: "Contains information about the set of data centers assigned to this domain.",
				Elem: &schema.Resource{
					Schema: map[string]*schema.Schema{
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
							Description: "A two-letter ISO 3166 country code that specifies the country where the data center is located.",
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
				},
			},
		},
	}
}

func dataGTMDatacentersRead(ctx context.Context, d *schema.ResourceData, m interface{}) diag.Diagnostics {
	meta := meta.Must(m)
	logger := meta.Log("Akamai GTM", "dataGTMDatacentersRead")
	ctx = session.ContextWithOptions(ctx, session.WithContextLog(logger))
	client := Client(meta)

	logger.Debug("Fetching datacenters")

	domain, err := tf.GetStringValue("domain", d)
	if err != nil {
		return diag.FromErr(err)
	}

	datacenters, err := client.ListDatacenters(ctx, gtm.ListDatacentersRequest{DomainName: domain})
	if err != nil {
		return diag.FromErr(err)
	}

	datacentersAttrs := make([]interface{}, len(datacenters))
	for i, dc := range datacenters {
		dcAttrs := getDatacenterAttributes(&dc)
		dcAttrs["datacenter_id"] = dc.DatacenterID
		datacentersAttrs[i] = dcAttrs
	}

	attrs := map[string]interface{}{
		"datacenters": datacentersAttrs,
	}

	if err = tf.SetAttrs(d, attrs); err != nil {
		return diag.FromErr(err)
	}

	d.SetId(domain)

	return nil
}
