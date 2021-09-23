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
				Type:     schema.TypeString,
				Required: true,
			},
			"version": {
				Type:     schema.TypeInt,
				Optional: true,
			},
			"description": {
				Type:     schema.TypeString,
				Computed: true,
			},
			"akamaized": {
				Type:     schema.TypeBool,
				Computed: true,
			},
			"checksum": {
				Type:     schema.TypeString,
				Computed: true,
			},
			"type": {
				Type:     schema.TypeString,
				Computed: true,
			},
			"balancing_type": {
				Type:     schema.TypeString,
				Computed: true,
			},
			"created_by": {
				Type:     schema.TypeString,
				Computed: true,
			},
			"created_date": {
				Type:     schema.TypeString,
				Computed: true,
			},
			"deleted": {
				Type:     schema.TypeBool,
				Computed: true,
			},
			"immutable": {
				Type:     schema.TypeBool,
				Computed: true,
			},
			"last_modified_by": {
				Type:     schema.TypeString,
				Computed: true,
			},
			"last_modified_date": {
				Type:     schema.TypeString,
				Computed: true,
			},
			"warnings": {
				Type:     schema.TypeString,
				Computed: true,
			},
			"data_centers": {
				Type:     schema.TypeSet,
				Computed: true,
				Elem: &schema.Resource{
					Schema: map[string]*schema.Schema{
						"city": {
							Type:     schema.TypeString,
							Computed: true,
						},
						"cloud_server_host_header_override": {
							Type:     schema.TypeBool,
							Computed: true,
						},
						"cloud_service": {
							Type:     schema.TypeBool,
							Computed: true,
						},
						"continent": {
							Type:     schema.TypeString,
							Computed: true,
						},
						"country": {
							Type:     schema.TypeString,
							Computed: true,
						},
						"hostname": {
							Type:     schema.TypeString,
							Computed: true,
						},
						"latitude": {
							Type:     schema.TypeFloat,
							Computed: true,
						},
						"liveness_hosts": {
							Type:     schema.TypeSet,
							Computed: true,
							Elem:     &schema.Schema{Type: schema.TypeString},
						},
						"longitude": {
							Type:     schema.TypeFloat,
							Computed: true,
						},
						"origin_id": {
							Type:     schema.TypeString,
							Computed: true,
						},
						"percent": {
							Type:     schema.TypeFloat,
							Computed: true,
						},
						"state_or_province": {
							Type:     schema.TypeString,
							Computed: true,
						},
					},
				},
			},
			"liveness_settings": {
				Type:     schema.TypeSet,
				Computed: true,
				Elem: &schema.Resource{
					Schema: map[string]*schema.Schema{
						"host_header": {
							Type:     schema.TypeString,
							Computed: true,
						},
						"additional_headers": {
							Type:     schema.TypeMap,
							Computed: true,
							Elem:     &schema.Schema{Type: schema.TypeString},
						},
						"interval": {
							Type:     schema.TypeInt,
							Computed: true,
						},
						"path": {
							Type:     schema.TypeString,
							Computed: true,
						},
						"peer_certificate_verification": {
							Type:     schema.TypeBool,
							Computed: true,
						},
						"port": {
							Type:     schema.TypeInt,
							Computed: true,
						},
						"protocol": {
							Type:     schema.TypeString,
							Computed: true,
						},
						"request_string": {
							Type:     schema.TypeString,
							Computed: true,
						},
						"response_string": {
							Type:     schema.TypeString,
							Computed: true,
						},
						"status_3xx_failure": {
							Type:     schema.TypeBool,
							Computed: true,
						},
						"status_4xx_failure": {
							Type:     schema.TypeBool,
							Computed: true,
						},
						"status_5xx_failure": {
							Type:     schema.TypeBool,
							Computed: true,
						},
						"timeout": {
							Type:     schema.TypeFloat,
							Computed: true,
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
		version, err = getTheLatestVersion(ctx, originID, client)
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
		return diag.Errorf("specified load balancer version is deleted: version = %d", version)
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
		"akamaized":          origin.Akamaized,
		"checksum":           origin.Checksum,
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

func getTheLatestVersion(ctx context.Context, originID string, client cloudlets.Cloudlets) (int64, error) {
	listLoadBalancerVersionsRequest := cloudlets.ListLoadBalancerVersionsRequest{
		OriginID: originID,
	}

	versions, err := client.ListLoadBalancerVersions(ctx, listLoadBalancerVersionsRequest)
	if err != nil {
		return 0, err
	}
	if len(versions) == 0 {
		return 0, fmt.Errorf("there is no any load balancer version")
	}

	var theLatestVersion int64
	for _, version := range versions {
		if version.Version > theLatestVersion {
			theLatestVersion = version.Version
		}
	}
	return theLatestVersion, nil
}
