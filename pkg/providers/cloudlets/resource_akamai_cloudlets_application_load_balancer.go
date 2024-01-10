package cloudlets

import (
	"context"
	"encoding/json"
	"errors"
	"fmt"
	"strconv"
	"strings"

	"github.com/akamai/AkamaiOPEN-edgegrid-golang/v7/pkg/cloudlets"
	"github.com/akamai/AkamaiOPEN-edgegrid-golang/v7/pkg/session"
	"github.com/akamai/terraform-provider-akamai/v5/pkg/common/tf"
	"github.com/akamai/terraform-provider-akamai/v5/pkg/meta"
	ozzo "github.com/go-ozzo/ozzo-validation/v4"
	"github.com/go-ozzo/ozzo-validation/v4/is"
	"github.com/hashicorp/go-cty/cty"
	"github.com/hashicorp/terraform-plugin-sdk/v2/diag"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/customdiff"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/schema"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/validation"
)

func resourceCloudletsApplicationLoadBalancer() *schema.Resource {
	return &schema.Resource{
		CustomizeDiff: customdiff.All(
			EnforceVersionChange,
			ensureTotalPercentageSum,
		),
		CreateContext: resourceALBCreate,
		ReadContext:   resourceALBRead,
		UpdateContext: resourceALBUpdate,
		DeleteContext: resourceALBDelete,
		Schema: map[string]*schema.Schema{
			"origin_id": {
				Type:        schema.TypeString,
				Required:    true,
				ForceNew:    true,
				Description: "The conditional origin's unique identifier",
			},
			"description": {
				Type:        schema.TypeString,
				Optional:    true,
				Description: "The load balancer configuration version description",
			},
			"origin_description": {
				Type:        schema.TypeString,
				Optional:    true,
				Description: "The load balancer configuration description",
			},
			"balancing_type": {
				Type:     schema.TypeString,
				Optional: true,
				ValidateDiagFunc: validation.ToDiagFunc(validation.StringInSlice([]string{
					"WEIGHTED",
					"PERFORMANCE",
				}, false)),
				Description: "The type of load balancing being performed. Options include WEIGHTED and PERFORMANCE",
			},
			"data_centers": {
				Type:        schema.TypeSet,
				Required:    true,
				Description: "The object containing information on conditional origins being used as data centers for an Application Load Balancer implementation. Only Conditional Origins with an originType of CUSTOMER or NETSTORAGE can be used as data centers in an application load balancer configuration.",
				Elem: &schema.Resource{
					Schema: map[string]*schema.Schema{
						"latitude": {
							Type:        schema.TypeFloat,
							Required:    true,
							Description: "The latitude value for the data center. This member supports six decimal places of precision.",
						},
						"longitude": {
							Type:        schema.TypeFloat,
							Required:    true,
							Description: "The longitude value for the data center. This member supports six decimal places of precision.",
						},
						"continent": {
							Type:        schema.TypeString,
							Required:    true,
							Description: "The continent on which the data center is located",
							ValidateDiagFunc: validation.ToDiagFunc(validation.StringInSlice([]string{
								"AF",
								"AS",
								"EU",
								"NA",
								"OC",
								"OT",
								"SA",
							}, false)),
						},
						"country": {
							Type:        schema.TypeString,
							Required:    true,
							Description: "The country in which the data center is located",
							ValidateDiagFunc: func(i interface{}, path cty.Path) diag.Diagnostics {
								if err := ozzo.Validate(strings.ToUpper(i.(string)), is.CountryCode2); err != nil {
									return diag.FromErr(err)
								}
								return nil
							},
						},
						"origin_id": {
							Type:        schema.TypeString,
							Required:    true,
							Description: "The ID of an origin that represents the data center. The conditional origin, which is defined in the Property Manager API, must have an originType of either CUSTOMER or NET_STORAGE",
						},
						"percent": {
							Type:        schema.TypeFloat,
							Required:    true,
							Description: "The percent of traffic that is sent to the data center. The total for all data centers must equal 100%.",
						},
						"cloud_service": {
							Type:        schema.TypeBool,
							Default:     false,
							Optional:    true,
							Description: "Describes if this datacenter is a cloud service",
						},
						"liveness_hosts": {
							Type:        schema.TypeList,
							Elem:        &schema.Schema{Type: schema.TypeString},
							Optional:    true,
							Description: "An array of strings that represent the origin servers used to poll the data centers in an application load balancer configuration. These servers support basic HTTP polling.",
						},
						"hostname": {
							Type:        schema.TypeString,
							Optional:    true,
							Description: "This should match the 'hostname' value defined for this datacenter in Property Manager",
						},
						"state_or_province": {
							Type:        schema.TypeString,
							Optional:    true,
							Description: "The state, province, or region where the data center is located",
						},
						"city": {
							Type:        schema.TypeString,
							Optional:    true,
							Description: "The city in which the data center is located.",
						},
						"cloud_server_host_header_override": {
							Type:        schema.TypeBool,
							Default:     false,
							Optional:    true,
							Description: "Describes if cloud server host header is overridden",
						},
					},
				},
			},
			"liveness_settings": {
				Type:     schema.TypeList,
				Optional: true,
				MaxItems: 1,
				Elem: &schema.Resource{
					Schema: map[string]*schema.Schema{
						"port": {
							Type:        schema.TypeInt,
							Required:    true,
							Description: "The port for the test object. The default port is 80, which is standard for HTTP. Enter 443 if you are using HTTPS.",
						},
						"protocol": {
							Type:        schema.TypeString,
							Required:    true,
							Description: "The protocol or scheme for the database, either HTTP or HTTPS.",
						},
						"path": {
							Type:        schema.TypeString,
							Required:    true,
							Description: "The path to the test object used for liveness testing. The function of the test object is to help determine whether the data center is functioning.",
						},
						"host_header": {
							Type:        schema.TypeString,
							Optional:    true,
							Description: "The Host header for the liveness HTTP request",
						},
						"additional_headers": {
							Type:        schema.TypeMap,
							Elem:        &schema.Schema{Type: schema.TypeString},
							Optional:    true,
							Description: "Maps additional case-insensitive HTTP header names included to the liveness testing requests",
						},
						"interval": {
							Type:        schema.TypeInt,
							Default:     0,
							Optional:    true,
							Description: "Describes how often the liveness test will be performed. Optional defaults to 60 seconds, minimum is 10 seconds.",
						},
						"peer_certificate_verification": {
							Type:        schema.TypeBool,
							Default:     false,
							Optional:    true,
							Description: "Describes whether or not to validate the origin certificate for an HTTPS request",
						},
						"request_string": {
							Type:        schema.TypeString,
							Optional:    true,
							Description: "The request which will be used for TCP(S) tests",
						},
						"response_string": {
							Type:     schema.TypeString,
							Optional: true,
						},
						"status_3xx_failure": {
							Type:        schema.TypeBool,
							Default:     false,
							Optional:    true,
							Description: "Set to true to mark the liveness test as failed when the request returns a 3xx (redirection) status code.",
						},
						"status_4xx_failure": {
							Type:        schema.TypeBool,
							Default:     false,
							Optional:    true,
							Description: "Set to true to mark the liveness test as failed when the request returns a 4xx (client error) status code.",
						},
						"status_5xx_failure": {
							Type:        schema.TypeBool,
							Default:     false,
							Optional:    true,
							Description: "Set to true to mark the liveness test as failed when the request returns a 5xx (server error) status code.",
						},
						"timeout": {
							Type:        schema.TypeFloat,
							Default:     float64(0),
							Optional:    true,
							Description: "The number of seconds the system waits before failing the liveness test. The default is 25 seconds.",
						},
					},
				},
			},
			"version": {
				Type:        schema.TypeInt,
				Computed:    true,
				Description: "The load balancer configuration version",
			},
			"warnings": {
				Type:        schema.TypeString,
				Computed:    true,
				Description: "Describes warnings during activation of load balancer configuration",
			},
		},
		Importer: &schema.ResourceImporter{
			StateContext: resourceALBImport,
		},
	}
}

// ensureTotalPercentageSum ensures that all datacenters' percent sum up to 100
// and will veto the diff altogether and abort the plan when above condition is not met
// This is a workaround as CustomizeDiff function in not meant to be used for resource validation
// but there is no easier solution to do it with the current state of terraform sdk
func ensureTotalPercentageSum(_ context.Context, diff *schema.ResourceDiff, _ interface{}) error {
	if dcs, ok := diff.GetOk("data_centers"); ok {
		dataCenters := dcs.(*schema.Set).List()
		var total float64
		for _, dataCenter := range dataCenters {
			dc := dataCenter.(map[string]interface{})
			percent := dc["percent"].(float64)
			total += percent
		}
		if 100.0 != total {
			t := strconv.FormatFloat(total, 'f', -1, 64)
			return fmt.Errorf("the total data center percentage must be 100%%: total=%s%%", t)
		}
	}
	return nil
}

// EnforceVersionChange enforces that change to any field will most likely result in creating a new version
func EnforceVersionChange(_ context.Context, diff *schema.ResourceDiff, _ interface{}) error {
	if diff.HasChange("origin_id") ||
		diff.HasChange("description") ||
		diff.HasChange("balancing_type") ||
		diff.HasChange("data_centers") ||
		diff.HasChange("liveness_settings") ||
		diff.HasChange("version") {
		return diff.SetNewComputed("version")
	}
	return nil
}

func isAkamaized(dc cloudlets.DataCenter, origins []cloudlets.OriginResponse) bool {
	for _, o := range origins {
		if o.Hostname == dc.Hostname && o.OriginID == dc.OriginID {
			return o.Akamaized
		}
	}
	return false
}

func validateLivenessHosts(ctx context.Context, client cloudlets.Cloudlets, d *schema.ResourceData) error {
	dcs := getDataCenters(d)

	origins, err := client.ListOrigins(ctx, cloudlets.ListOriginsRequest{})
	if err != nil {
		return err
	}

	for _, dc := range dcs {
		if len(dc.LivenessHosts) > 0 && isAkamaized(dc, origins) {
			return fmt.Errorf("'liveness_hosts' field should be omitted for GTM hostname: %q. "+
				"Liveness tests for this host can be configured in DNS traffic management", dc.Hostname)
		}
	}
	return nil
}

func resourceALBCreate(ctx context.Context, d *schema.ResourceData, m interface{}) diag.Diagnostics {
	meta := meta.Must(m)
	logger := meta.Log("Cloudlets", "resourceLoadBalancerConfigurationCreate")
	// create a context with logging for api calls
	ctx = session.ContextWithOptions(
		ctx,
		session.WithContextLog(logger),
	)
	client := Client(meta)
	logger.Debug("Creating load balancer configuration")
	originID, err := tf.GetStringValue("origin_id", d)
	if err != nil {
		return diag.FromErr(err)
	}

	originDescription, err := tf.GetStringValue("origin_description", d)
	if err != nil && !errors.Is(err, tf.ErrNotFound) {
		return diag.FromErr(err)
	}

	if err := validateLivenessHosts(ctx, client, d); err != nil {
		return diag.FromErr(err)
	}

	createLBConfigResp, err := client.CreateOrigin(ctx, cloudlets.CreateOriginRequest{
		OriginID: originID,
		Description: cloudlets.Description{
			Description: originDescription,
		},
	})
	if err != nil {
		return diag.FromErr(err)
	}
	d.SetId(createLBConfigResp.OriginID)
	loadBalancerVersion := getLoadBalancerVersion(d)
	createVersionResp, err := client.CreateLoadBalancerVersion(ctx, cloudlets.CreateLoadBalancerVersionRequest{
		OriginID:            originID,
		LoadBalancerVersion: loadBalancerVersion,
	})
	if err != nil {
		return diag.FromErr(err)
	}
	if err := d.Set("version", createVersionResp.Version); err != nil {
		return diag.FromErr(fmt.Errorf("%w: %s", tf.ErrValueSet, err.Error()))
	}
	return resourceALBRead(ctx, d, m)
}

func resourceALBRead(ctx context.Context, d *schema.ResourceData, m interface{}) diag.Diagnostics {
	meta := meta.Must(m)
	logger := meta.Log("Cloudlets", "resourceLoadBalancerConfigurationRead")
	// create a context with logging for api calls
	ctx = session.ContextWithOptions(
		ctx,
		session.WithContextLog(logger),
	)
	client := Client(meta)
	logger.Debug("Reading load balancer configuration")
	originID := d.Id()
	loadBalancerConfigAttrs := map[string]interface{}{
		"origin_id": originID,
	}
	origin, err := client.GetOrigin(ctx, cloudlets.GetOriginRequest{
		OriginID: originID,
	})
	if err != nil {
		return diag.FromErr(err)
	}

	if err := tf.SetAttrs(d, loadBalancerConfigAttrs); err != nil {
		return diag.FromErr(err)
	}
	version, err := tf.GetIntValue("version", d)
	if err != nil {
		return diag.FromErr(err)
	}
	loadBalancerVersion, err := client.GetLoadBalancerVersion(ctx, cloudlets.GetLoadBalancerVersionRequest{
		OriginID:       originID,
		Version:        int64(version),
		ShouldValidate: true,
	})
	if err != nil {
		return diag.FromErr(err)
	}
	attrs := make(map[string]interface{})
	attrs["origin_description"] = origin.Description
	attrs["balancing_type"] = loadBalancerVersion.BalancingType
	attrs["version"] = loadBalancerVersion.Version
	attrs["description"] = loadBalancerVersion.Description
	attrs["data_centers"] = populateDataCenters(loadBalancerVersion.DataCenters)
	if loadBalancerVersion.LivenessSettings != nil {
		attrs["liveness_settings"] = populateLivenessSettings(loadBalancerVersion.LivenessSettings)
	}
	var warningsJSON []byte
	if len(loadBalancerVersion.Warnings) > 0 {
		warningsJSON, err = json.MarshalIndent(loadBalancerVersion.Warnings, "", "  ")
		if err != nil {
			return diag.FromErr(err)
		}
	}
	attrs["warnings"] = string(warningsJSON)
	if err := tf.SetAttrs(d, attrs); err != nil {
		return diag.FromErr(err)
	}
	return nil
}

func resourceALBUpdate(ctx context.Context, d *schema.ResourceData, m interface{}) diag.Diagnostics {
	meta := meta.Must(m)
	logger := meta.Log("Cloudlets", "resourceLoadBalancerConfigurationUpdate")
	// create a context with logging for api calls
	ctx = session.ContextWithOptions(
		ctx,
		session.WithContextLog(logger),
	)
	client := Client(meta)
	logger.Debug("Updating load balancer configuration")
	originID := d.Id()

	version, err := tf.GetIntValue("version", d)
	if err != nil {
		return diag.FromErr(err)
	}

	if err := validateLivenessHosts(ctx, client, d); err != nil {
		return diag.FromErr(err)
	}

	if d.HasChange("origin_description") {
		originDescription, err := tf.GetStringValue("origin_description", d)
		if err != nil && !errors.Is(err, tf.ErrNotFound) {
			return diag.FromErr(err)
		}
		_, err = client.UpdateOrigin(ctx, cloudlets.UpdateOriginRequest{
			OriginID: originID,
			Description: cloudlets.Description{
				Description: originDescription,
			},
		})
		if err != nil {
			return diag.FromErr(err)
		}
	}

	// if version-related attributes have changed, load balancer version has to be either created or updated (depending on whether it's active or not)
	if d.HasChanges("description", "balancing_type", "data_centers", "liveness_settings") {
		activations, err := client.ListLoadBalancerActivations(ctx, cloudlets.ListLoadBalancerActivationsRequest{OriginID: originID})
		if err != nil {
			return diag.FromErr(err)
		}
		var versionActive bool
		for _, activation := range activations {
			if activation.Version == int64(version) && activation.Status == cloudlets.LoadBalancerActivationStatusActive {
				versionActive = true
				break
			}
		}
		loadBalancerVersion := getLoadBalancerVersion(d)
		var loadBalancerVersionResp *cloudlets.LoadBalancerVersion
		if versionActive {
			loadBalancerVersionResp, err = client.CreateLoadBalancerVersion(ctx, cloudlets.CreateLoadBalancerVersionRequest{
				OriginID:            originID,
				LoadBalancerVersion: loadBalancerVersion,
			})
		} else {
			loadBalancerVersionResp, err = client.UpdateLoadBalancerVersion(ctx, cloudlets.UpdateLoadBalancerVersionRequest{
				OriginID:            originID,
				ShouldValidate:      true,
				Version:             int64(version),
				LoadBalancerVersion: loadBalancerVersion,
			})
		}
		if err != nil {
			return diag.FromErr(err)
		}
		if err := d.Set("version", loadBalancerVersionResp.Version); err != nil {
			return diag.FromErr(fmt.Errorf("%w: %s", tf.ErrValueSet, err.Error()))
		}
	}
	return resourceALBRead(ctx, d, m)
}

// resourceALBDelete does not call any delete operation in the API, because there is no such operation
// resource will simply be removed from state in that case
// to allow re-using existing config, create function also covers the import functionality, saving the existing origin and version in state
func resourceALBDelete(_ context.Context, d *schema.ResourceData, m interface{}) diag.Diagnostics {
	meta := meta.Must(m)
	logger := meta.Log("Cloudlets", "resourceLoadBalancerConfigurationDelete")
	logger.Debug("Deleting load balancer configuration")
	logger.Info("Cloudlets API does not support load balancer configuration and load balancer configuration version deletion - resource will only be removed from state")
	d.SetId("")
	return nil
}

// resourceALBImport does a basic import based on the originID specified it imports the latest version
func resourceALBImport(ctx context.Context, d *schema.ResourceData, m interface{}) ([]*schema.ResourceData, error) {
	meta := meta.Must(m)
	logger := meta.Log("Cloudlets", "resourceALBImport")
	logger.Debug("Import ALB")

	client := Client(meta)
	logger.Debug("Importing load balancer configuration")

	originID := d.Id()
	if originID == "" {
		return nil, fmt.Errorf("origin id cannot be empty")
	}

	origin, err := client.GetOrigin(ctx, cloudlets.GetOriginRequest{OriginID: originID})
	if err != nil {
		return nil, err
	}

	if origin == nil {
		return nil, fmt.Errorf("could not find origin with origin_id: %s", originID)
	}

	versions, err := client.ListLoadBalancerVersions(ctx, cloudlets.ListLoadBalancerVersionsRequest{
		OriginID: origin.OriginID,
	})
	if err != nil {
		return nil, err
	}

	if len(versions) == 0 {
		return nil, fmt.Errorf("no load balancer version found for origin_id: %s", originID)
	}

	var version int64
	for _, v := range versions {
		if v.Version > version {
			version = v.Version
		}
	}

	err = d.Set("version", version)
	if err != nil {
		return nil, err
	}

	return []*schema.ResourceData{d}, nil
}

func populateDataCenters(dcs []cloudlets.DataCenter) []interface{} {
	dataCentersList := make([]interface{}, 0)
	for _, dc := range dcs {
		dsMap := map[string]interface{}{
			"cloud_service":                     dc.CloudService,
			"liveness_hosts":                    dc.LivenessHosts,
			"latitude":                          dc.Latitude,
			"longitude":                         dc.Longitude,
			"continent":                         dc.Continent,
			"country":                           dc.Country,
			"origin_id":                         dc.OriginID,
			"percent":                           dc.Percent,
			"hostname":                          dc.Hostname,
			"city":                              dc.City,
			"cloud_server_host_header_override": dc.CloudServerHostHeaderOverride,
		}
		if dc.StateOrProvince != nil {
			dsMap["state_or_province"] = *dc.StateOrProvince
		}
		dataCentersList = append(dataCentersList, dsMap)
	}
	return dataCentersList
}

func populateLivenessSettings(ls *cloudlets.LivenessSettings) []interface{} {
	lsMap := map[string]interface{}{
		"port":                          ls.Port,
		"protocol":                      ls.Protocol,
		"host_header":                   ls.HostHeader,
		"additional_headers":            ls.AdditionalHeaders,
		"interval":                      ls.Interval,
		"path":                          ls.Path,
		"peer_certificate_verification": ls.PeerCertificateVerification,
		"request_string":                ls.RequestString,
		"response_string":               ls.ResponseString,
		"status_3xx_failure":            ls.Status3xxFailure,
		"status_4xx_failure":            ls.Status4xxFailure,
		"status_5xx_failure":            ls.Status5xxFailure,
		"timeout":                       ls.Timeout,
	}
	return []interface{}{lsMap}
}

func getLoadBalancerVersion(d *schema.ResourceData) cloudlets.LoadBalancerVersion {
	description := d.Get("description").(string)
	balancingType := d.Get("balancing_type").(string)
	dataCenters := getDataCenters(d)
	livenessSettings := getLivenessSettings(d)
	return cloudlets.LoadBalancerVersion{
		Description:      description,
		BalancingType:    cloudlets.BalancingType(balancingType),
		DataCenters:      dataCenters,
		LivenessSettings: livenessSettings,
	}
}

func getDataCenters(d *schema.ResourceData) []cloudlets.DataCenter {
	dataCentersSet := d.Get("data_centers").(*schema.Set)
	dataCenters := make([]cloudlets.DataCenter, dataCentersSet.Len())
	for i, dc := range dataCentersSet.List() {
		dcMap := dc.(map[string]interface{})
		livenessHosts := dcMap["liveness_hosts"].([]interface{})
		livenessHostsStr := make([]string, len(livenessHosts))
		for i, host := range livenessHosts {
			livenessHostsStr[i] = host.(string)
		}
		var stateOrProvince *string
		if s := dcMap["state_or_province"].(string); s != "" {
			stateOrProvince = &s
		}
		dataCenters[i] = cloudlets.DataCenter{
			City:                          dcMap["city"].(string),
			CloudServerHostHeaderOverride: dcMap["cloud_server_host_header_override"].(bool),
			CloudService:                  dcMap["cloud_service"].(bool),
			Continent:                     dcMap["continent"].(string),
			Country:                       dcMap["country"].(string),
			Hostname:                      dcMap["hostname"].(string),
			Latitude:                      getFloat64PtrValue(dcMap, "latitude"),
			LivenessHosts:                 livenessHostsStr,
			Longitude:                     getFloat64PtrValue(dcMap, "longitude"),
			OriginID:                      dcMap["origin_id"].(string),
			Percent:                       getFloat64PtrValue(dcMap, "percent"),
			StateOrProvince:               stateOrProvince,
		}
	}
	return dataCenters
}

func getLivenessSettings(d *schema.ResourceData) *cloudlets.LivenessSettings {
	lsList, err := tf.GetListValue("liveness_settings", d)
	if err != nil {
		return nil
	}
	lsMap := lsList[0].(map[string]interface{})
	additionalHeaders := lsMap["additional_headers"].(map[string]interface{})
	additionalHeadersStr := make(map[string]string, len(additionalHeaders))
	for k, v := range additionalHeaders {
		additionalHeadersStr[k] = v.(string)
	}
	return &cloudlets.LivenessSettings{
		HostHeader:                  lsMap["host_header"].(string),
		AdditionalHeaders:           additionalHeadersStr,
		Interval:                    lsMap["interval"].(int),
		Path:                        lsMap["path"].(string),
		PeerCertificateVerification: lsMap["peer_certificate_verification"].(bool),
		Port:                        lsMap["port"].(int),
		Protocol:                    lsMap["protocol"].(string),
		RequestString:               lsMap["request_string"].(string),
		ResponseString:              lsMap["response_string"].(string),
		Status3xxFailure:            lsMap["status_3xx_failure"].(bool),
		Status4xxFailure:            lsMap["status_4xx_failure"].(bool),
		Status5xxFailure:            lsMap["status_5xx_failure"].(bool),
		Timeout:                     lsMap["timeout"].(float64),
	}
}
