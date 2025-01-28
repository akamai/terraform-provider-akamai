package cloudlets

import (
	"context"
	"encoding/json"
	"errors"
	"fmt"

	"github.com/akamai/AkamaiOPEN-edgegrid-golang/v10/pkg/cloudlets"

	"github.com/akamai/terraform-provider-akamai/v6/pkg/common/tf"
	"github.com/akamai/terraform-provider-akamai/v6/pkg/meta"
	"github.com/hashicorp/terraform-plugin-sdk/v2/diag"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/schema"
)

func dataSourceCloudletsApplicationLoadBalancer() *schema.Resource {
	return &schema.Resource{
		ReadContext: dataApplicationLoadBalancerRead,
		Schema: map[string]*schema.Schema{
			"origin_id": {
				Type:        schema.TypeString,
				Required:    true,
				Description: "The conditional origin’s unique identifier",
			},
			"version": {
				Type:        schema.TypeInt,
				Optional:    true,
				Description: "The load balancer configuration version",
			},
			"description": {
				Type:        schema.TypeString,
				Computed:    true,
				Description: "The load balancer configuration description",
			},
			"type": {
				Type:        schema.TypeString,
				Computed:    true,
				Description: "The type of conditional origin",
			},
			"balancing_type": {
				Type:        schema.TypeString,
				Computed:    true,
				Description: "The type of load balancing being performed. Options include WEIGHTED and PERFORMANCE",
			},
			"created_by": {
				Type:        schema.TypeString,
				Computed:    true,
				Description: "The value which is set by the server at the time of creation and never subsequently changes",
			},
			"created_date": {
				Type:        schema.TypeString,
				Computed:    true,
				Description: "The created date which is only set by the server the first time the load balancer version is created. All subsequent responses will contain the same value.",
			},
			"deleted": {
				Type:        schema.TypeBool,
				Computed:    true,
				Description: "Describes if load balancer version has been deleted",
			},
			"immutable": {
				Type:        schema.TypeBool,
				Computed:    true,
				Description: "Denotes whether you can edit the load balancing version. The default setting for this member is false. It automatically becomes true when the load balancing version is activated for the first time.",
			},
			"last_modified_by": {
				Type:        schema.TypeString,
				Computed:    true,
				Description: "The last modification of load balancer configuration which is set by the server",
			},
			"last_modified_date": {
				Type:        schema.TypeString,
				Computed:    true,
				Description: "Describes when load balancer configuration has been modified for the last time",
			},
			"warnings": {
				Type:        schema.TypeString,
				Computed:    true,
				Description: "Describes warnings during activation of load balancer configuration",
			},
			"data_centers": {
				Type:        schema.TypeSet,
				Computed:    true,
				Description: "The object containing information on conditional origins being used as data centers for an Application Load Balancer implementation. Only Conditional Origins with an originType of CUSTOMER or NETSTORAGE can be used as data centers in an application load balancer configuration.",
				Elem: &schema.Resource{
					Schema: map[string]*schema.Schema{
						"city": {
							Type:        schema.TypeString,
							Computed:    true,
							Description: "The city in which the data center is located.",
						},
						"cloud_server_host_header_override": {
							Type:        schema.TypeBool,
							Computed:    true,
							Description: "Describes if cloud server host header is overridden",
						},
						"cloud_service": {
							Type:        schema.TypeBool,
							Computed:    true,
							Description: "Describes if this datacenter is a cloud service",
						},
						"continent": {
							Type:        schema.TypeString,
							Computed:    true,
							Description: "The continent on which the data center is located",
						},
						"country": {
							Type:        schema.TypeString,
							Computed:    true,
							Description: "The country in which the data center is located",
						},
						"hostname": {
							Type:        schema.TypeString,
							Computed:    true,
							Description: "This should match the 'hostname' value defined for this datacenter in Property Manager",
						},
						"latitude": {
							Type:        schema.TypeFloat,
							Computed:    true,
							Description: "The latitude value for the data center. This member supports six decimal places of precision.",
						},
						"liveness_hosts": {
							Type:        schema.TypeSet,
							Computed:    true,
							Elem:        &schema.Schema{Type: schema.TypeString},
							Description: "An array of strings that represent the origin servers used to poll the data centers in an application load balancer configuration. These servers support basic HTTP polling.",
						},
						"longitude": {
							Type:        schema.TypeFloat,
							Computed:    true,
							Description: "The longitude value for the data center. This member supports six decimal places of precision.",
						},
						"origin_id": {
							Type:        schema.TypeString,
							Computed:    true,
							Description: "The id of an origin that represents the data center. The conditional origin, which is defined in the Property Manager API, must have an originType of either CUSTOMER or NET_STORAGE",
						},
						"percent": {
							Type:        schema.TypeFloat,
							Computed:    true,
							Description: "The percent of traffic that is sent to the data center. The total for all data centers must equal 100%.",
						},
						"state_or_province": {
							Type:        schema.TypeString,
							Computed:    true,
							Description: "The state, province, or region where the data center is located",
						},
					},
				},
			},
			"liveness_settings": {
				Type:        schema.TypeSet,
				Computed:    true,
				Description: "The liveness settings are used to determine the health of each load balanced data center defined in the data center list",
				Elem: &schema.Resource{
					Schema: map[string]*schema.Schema{
						"host_header": {
							Type:        schema.TypeString,
							Computed:    true,
							Description: "The Host header for the liveness HTTP request",
						},
						"additional_headers": {
							Type:        schema.TypeMap,
							Computed:    true,
							Elem:        &schema.Schema{Type: schema.TypeString},
							Description: "Maps additional case-insensitive HTTP header names included to the liveness testing requests",
						},
						"interval": {
							Type:        schema.TypeInt,
							Computed:    true,
							Description: "Describes how often the liveness test will be performed. Optional defaults to 60 seconds, minimum is 10 seconds.",
						},
						"path": {
							Type:        schema.TypeString,
							Computed:    true,
							Description: "The path to the test object used for liveness testing. The function of the test object is to help determine whether the data center is functioning.",
						},
						"peer_certificate_verification": {
							Type:        schema.TypeBool,
							Computed:    true,
							Description: "Describes whether or not to validate the origin certificate for an HTTPS request",
						},
						"port": {
							Type:        schema.TypeInt,
							Computed:    true,
							Description: "The port for the test object. The default port is 80, which is standard for HTTP. Enter 443 if you are using HTTPS.",
						},
						"protocol": {
							Type:        schema.TypeString,
							Computed:    true,
							Description: "The protocol or scheme for the database, either HTTP or HTTPS.",
						},
						"request_string": {
							Type:        schema.TypeString,
							Computed:    true,
							Description: "The request which will be used for TCP(S) tests",
						},
						"response_string": {
							Type:        schema.TypeString,
							Computed:    true,
							Description: "The response which will be used for TCP(S) tests",
						},
						"status_3xx_failure": {
							Type:        schema.TypeBool,
							Computed:    true,
							Description: "Set to true to mark the liveness test as failed when the request returns a 3xx (redirection) status code.",
						},
						"status_4xx_failure": {
							Type:        schema.TypeBool,
							Computed:    true,
							Description: "Set to true to mark the liveness test as failed when the request returns a 4xx (client error) status code.",
						},
						"status_5xx_failure": {
							Type:        schema.TypeBool,
							Computed:    true,
							Description: "Set to true to mark the liveness test as failed when the request returns a 5xx (server error) status code.",
						},
						"timeout": {
							Type:        schema.TypeFloat,
							Computed:    true,
							Description: "The number of seconds the system waits before failing the liveness test. The default is 25 seconds.",
						},
					},
				},
			},
		},
	}
}

func dataApplicationLoadBalancerRead(ctx context.Context, d *schema.ResourceData, m interface{}) diag.Diagnostics {
	meta := meta.Must(m)
	client := Client(meta)
	log := meta.Log("Cloudlets", "dataApplicationLoadBalancerRead")
	log.Debug("Reading Load Balancer")

	originID, err := tf.GetStringValue("origin_id", d)
	if err != nil {
		return diag.FromErr(err)
	}

	origin, err := client.GetOrigin(ctx, cloudlets.GetOriginRequest{OriginID: originID})
	if err != nil {
		return diag.FromErr(err)
	}

	var version int64
	if v, err := tf.GetIntValue("version", d); err != nil {
		if !errors.Is(err, tf.ErrNotFound) {
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

	err = tf.SetAttrs(d, fields)
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

var errNoVersionForOrigin = errors.New("no load balancer version found for origin")

func getLatestVersionOfApplicationLoadBalancer(ctx context.Context, originID string, client cloudlets.Cloudlets) (int64, error) {
	listLoadBalancerVersionsRequest := cloudlets.ListLoadBalancerVersionsRequest{
		OriginID: originID,
	}

	versions, err := client.ListLoadBalancerVersions(ctx, listLoadBalancerVersionsRequest)
	if err != nil {
		return 0, err
	}
	if len(versions) == 0 {
		return 0, fmt.Errorf("%w: %s", errNoVersionForOrigin, originID)
	}

	var theLatestVersion int64
	for _, version := range versions {
		if version.Version > theLatestVersion {
			theLatestVersion = version.Version
		}
	}
	return theLatestVersion, nil
}
