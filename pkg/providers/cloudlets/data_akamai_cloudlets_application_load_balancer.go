package cloudlets

import (
	"context"
	"encoding/json"
	"errors"
	"fmt"

	"github.com/akamai/AkamaiOPEN-edgegrid-golang/v2/pkg/cloudlets"

	"github.com/akamai/terraform-provider-akamai/v2/pkg/akamai"
	"github.com/akamai/terraform-provider-akamai/v2/pkg/tools"

	"github.com/hashicorp/terraform-plugin-sdk/v2/diag"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/schema"
)

func dataSourceCloudletsApplicationLoadBalancer() *schema.Resource {
	return &schema.Resource{
		ReadContext: dataApplicationLoadBalancerRead,
		Schema: map[string]*schema.Schema{
			"origin_id": {
				Type:		 schema.TypeString,
				Required: 	 true,
				Description: "Describes Origin Id",
			},
			"version": {
				Type:     	 schema.TypeInt,
				Optional: 	 true,
				Description: "Describes load balancer configuration version",
			},
			"description": {
				Type:     	 schema.TypeString,
				Computed: 	 true,
				Description: "Describes load balancer configuration",
			},
			"type": {
				Type:     	 schema.TypeString,
				Computed: 	 true,
				Description: "Describes Origin type",
			},
			"balancing_type": {
				Type:     	 schema.TypeString,
				Computed: 	 true,
				Description: "Load balancer configuration type",
			},
			"created_by": {
				Type:     	 schema.TypeString,
				Computed: 	 true,
				Description: "Describes value which is set by the server at the time of creation and never subsequently changes",
			},
			"created_date": {
				Type:     	 schema.TypeString,
				Computed: 	 true,
				Description: "Describes the created date which is only set by the server the first time the load balancer version is created. All subsequent responses will contain the same value.",
			},
			"deleted": {
				Type:     	 schema.TypeBool,
				Computed: 	 true,
				Description: "Describes if load balancer version has been deleted",
			},
			"immutable": {
				Type:     	 schema.TypeBool,
				Computed: 	 true,
				Description: "Describes if the load balancer version is marked as immutable which means, that it has been activated",
			},
			"last_modified_by": {
				Type:     	 schema.TypeString,
				Computed: 	 true,
				Description: "Describes the last modification of load balancer configuration which is set by the server",
			},
			"last_modified_date": {
				Type:     	 schema.TypeString,
				Computed: 	 true,
				Description: "Describes when load balancer configuration has been modified for the last time",
			},
			"warnings": {
				Type:     	 schema.TypeString,
				Computed: 	 true,
				Description: "Describes warnings during activation of load balancer configuration",
			},
			"data_centers": {
				Type:     	 schema.TypeSet,
				Computed: 	 true,
				Description: "Describes list of data center configurations used for ALB load balancing. The total weight of all data centers in this list must be add to 100%.",
				Elem: &schema.Resource{
					Schema: map[string]*schema.Schema{
						"city": {
							Type:     	 schema.TypeString,
							Computed: 	 true,
							Description: "The name of the city where the data center is located.",
						},
						"cloud_server_host_header_override": {
							Type:     	 schema.TypeBool,
							Computed: 	 true,
							Description: "Describes if cloud server host header is overridden",
						},
						"cloud_service": {
							Type:     	 schema.TypeBool,
							Computed: 	 true,
							Description: "Describes if this datacenter is a cloud service",
						},
						"continent": {
							Type:     	 schema.TypeString,
							Computed: 	 true,
							Description: "Describes the two character ISO-3166 continent code for the data center's location",
						},
						"country": {
							Type:     	 schema.TypeString,
							Computed: 	 true,
							Description: "Describes the two character ISO-3166 country code for the data center's location",
						},
						"hostname": {
							Type:     	 schema.TypeString,
							Computed: 	 true,
							Description: "This should match the 'hostname' value defined for this datacenter in Property Manager",
						},
						"latitude": {
							Type:     	 schema.TypeFloat,
							Computed: 	 true,
							Description: "Describes latitude location where this data center is located",
						},
						"liveness_hosts": {
							Type:     	 schema.TypeSet,
							Computed: 	 true,
							Elem:     	 &schema.Schema{Type: schema.TypeString},
							Description: "Describes list of hosts which can be checked for data center liveness",
						},
						"longitude": {
							Type:    	 schema.TypeFloat,
							Computed: 	 true,
							Description: "Describes latitude location where this data center is located",
						},
						"origin_id": {
							Type:     	 schema.TypeString,
							Computed: 	 true,
							Description: "Describes the Cloudlets Origin Id corresponding to this data center",
						},
						"percent": {
							Type:     	 schema.TypeFloat,
							Computed: 	 true,
							Description: "Describes percent of traffic for this meta-origin overall should try to route to this data center",
						},
						"state_or_province": {
							Type:     	 schema.TypeString,
							Computed: 	 true,
							Description: "Describes name of the state or province where the data center is located",
						},
					},
				},
			},
			"liveness_settings": {
				Type:     	 schema.TypeSet,
				Computed: 	 true,
				Description: "The liveness settings are used to determine the health of each load balanced data center defined in the data center list",
				Elem: &schema.Resource{
					Schema: map[string]*schema.Schema{
						"host_header": {
							Type:     	 schema.TypeString,
							Computed: 	 true,
							Description: "Describes the Host header for the liveness HTTP request",
						},
						"additional_headers": {
							Type:     	 schema.TypeMap,
							Computed: 	 true,
							Elem:     	 &schema.Schema{Type: schema.TypeString},
							Description: "Describes the additional header for the leveness HTTP request",
						},
						"interval": {
							Type:     	 schema.TypeInt,
							Computed: 	 true,
							Description: "Describes how often the liveness test will be performed",
						},
						"path": {
							Type:     	 schema.TypeString,
							Computed:	 true,
							Description: "Describes the path that will be requested to test for liveness",
						},
						"peer_certificate_verification": {
							Type:     	 schema.TypeBool,
							Computed: 	 true,
							Description: "Describes whether or not to validate the origin certificate for an HTTPS request",
						},
						"port": {
							Type:     	 schema.TypeInt,
							Computed: 	 true,
							Description: "Describes the port that will be used for the HTTP request",
						},
						"protocol": {
							Type:     	 schema.TypeString,
							Computed: 	 true,
							Description: "Describes the protocol that will be used for the request. It can be HTTP, HTTPS, TCP, or TCPS",
						},
						"request_string": {
							Type:     	 schema.TypeString,
							Computed: 	 true,
							Description: "Describes the request which will be used for TCP(S) tests",
						},
						"response_string": {
							Type:     	 schema.TypeString,
							Computed: 	 true,
							Description: "Describes the response which will be used for TCP(S) tests",
						},
						"status_3xx_failure": {
							Type:     	 schema.TypeBool,
							Computed: 	 true,
							Description: "Describes whether HTTP status codes in the 3xx range are considered liveness failure",
						},
						"status_4xx_failure": {
							Type:     	 schema.TypeBool,
							Computed: 	 true,
							Description: "Describes whether HTTP status codes in the 4xx range are considered liveness failure",
						},
						"status_5xx_failure": {
							Type:     	 schema.TypeBool,
							Computed: 	 true,
							Description: "Describes whether HTTP status codes in the 5xx range are considered liveness failure",
						},
						"timeout": {
							Type:     	 schema.TypeFloat,
							Computed: 	 true,
							Description: "Describes the timeout in seconds for the HTTP request. After this period passes the test will be considered to have failed",
						},
					},
				},
			},
		},
	}
}

func dataApplicationLoadBalancerRead(ctx context.Context, d *schema.ResourceData, m interface{}) diag.Diagnostics {
	meta := akamai.Meta(m)
	client := inst.Client(meta)
	log := meta.Log("Cloudlets", "dataApplicationLoadBalancerRead")
	log.Debug("Reading Load Balancer")

	originID, err := tools.GetStringValue("origin_id", d)
	if err != nil {
		return diag.FromErr(err)
	}

	origin, err := client.GetOrigin(ctx, originID)

	if err != nil {
		return diag.FromErr(err)
	}

	var version int64
	if v, err := tools.GetIntValue("version", d); err != nil {
		if !errors.Is(err, tools.ErrNotFound) {
			return diag.FromErr(err)
		}
		version, err = getLatestVersionOfApplicationLoadBalancer(ctx, originID, client)
		if err != nil {
			return diag.FromErr(err)
		}
	} else {
		version = int64(v)
	}

	getLoadBalancerVersionRequest := cloudlets.GetLoadBalancerVersionRequest{
		OriginID:       originID,
		Version:        version,
		ShouldValidate: true,
	}

	loadBalancerVersion, err := client.GetLoadBalancerVersion(ctx, getLoadBalancerVersionRequest)
	if err != nil {
		return diag.FromErr(err)
	}

	if loadBalancerVersion.Deleted {
		return diag.Errorf("specified load balancer version is deleted: %d", version)
	}

	fields, err := getSchemaLoadBalancer(loadBalancerVersion, origin)
	if err != nil {
		return diag.FromErr(err)
	}

	err = tools.SetAttrs(d, fields)
	if err != nil {
		return diag.FromErr(err)
	}

	d.SetId(fmt.Sprintf("%s:%d", originID, version))

	return nil
}

func getSchemaLoadBalancer(loadBalancerVersion *cloudlets.LoadBalancerVersion, origin *cloudlets.Origin) (map[string]interface{}, error) {
	warnings, err := json.MarshalIndent(loadBalancerVersion.Warnings, "", " ")
	if err != nil {
		return nil, fmt.Errorf("cannot marshal json %s", err)
	}
	return map[string]interface{}{
		"origin_id":          loadBalancerVersion.OriginID,
		"version":            loadBalancerVersion.Version,
		"description":        loadBalancerVersion.Description,
		"type":               origin.Type,
		"balancing_type":     loadBalancerVersion.BalancingType,
		"created_by":         loadBalancerVersion.CreatedBy,
		"created_date":       loadBalancerVersion.CreatedDate,
		"deleted":            loadBalancerVersion.Deleted,
		"immutable":          loadBalancerVersion.Immutable,
		"last_modified_by":   loadBalancerVersion.LastModifiedBy,
		"last_modified_date": loadBalancerVersion.LastModifiedDate,
		"warnings":           string(warnings),
		"data_centers":       getSchemaDataCenters(loadBalancerVersion.DataCenters),
		"liveness_settings":  getSchemaLivenessSettings(loadBalancerVersion.LivenessSettings),
	}, nil
}

func getSchemaDataCenters(centers []cloudlets.DataCenter) []map[string]interface{} {
	var dataCenters []map[string]interface{}

	for _, c := range centers {
		dataCenter := map[string]interface{}{
			"city":                              c.City,
			"cloud_server_host_header_override": c.CloudServerHostHeaderOverride,
			"cloud_service":                     c.CloudService,
			"continent":                         c.Continent,
			"country":                           c.Country,
			"hostname":                          c.Hostname,
			"latitude":                          c.Latitude,
			"liveness_hosts":                    c.LivenessHosts,
			"longitude":                         c.Longitude,
			"origin_id":                         c.OriginID,
			"percent":                           c.Percent,
			"state_or_province":                 c.StateOrProvince,
		}
		dataCenters = append(dataCenters, dataCenter)
	}
	return dataCenters
}

func getSchemaLivenessSettings(settings *cloudlets.LivenessSettings) []map[string]interface{} {
	if settings != nil {
		return []map[string]interface{}{
			{
				"host_header":                   settings.HostHeader,
				"additional_headers":            settings.AdditionalHeaders,
				"interval":                      settings.Interval,
				"path":                          settings.Path,
				"peer_certificate_verification": settings.PeerCertificateVerification,
				"port":                          settings.Port,
				"protocol":                      settings.Protocol,
				"request_string":                settings.RequestString,
				"response_string":               settings.ResponseString,
				"status_3xx_failure":            settings.Status3xxFailure,
				"status_4xx_failure":            settings.Status4xxFailure,
				"status_5xx_failure":            settings.Status5xxFailure,
				"timeout":                       settings.Timeout,
			},
		}
	}
	return nil
}

func getLatestVersionOfApplicationLoadBalancer(ctx context.Context, originID string, client cloudlets.Cloudlets) (int64, error) {
	listLoadBalancerVersionsRequest := cloudlets.ListLoadBalancerVersionsRequest{
		OriginID: originID,
	}

	versions, err := client.ListLoadBalancerVersions(ctx, listLoadBalancerVersionsRequest)
	if err != nil {
		return 0, err
	}
	if len(versions) == 0 {
		return 0, fmt.Errorf("no load balancer version found for given origin")
	}

	var theLatestVersion int64
	for _, version := range versions {
		if version.Version > theLatestVersion {
			theLatestVersion = version.Version
		}
	}
	return theLatestVersion, nil
}
