package gtm

import (
	"context"
	"errors"
	"fmt"
	"sort"
	"strings"

	"github.com/akamai/AkamaiOPEN-edgegrid-golang/v8/pkg/gtm"
	"github.com/akamai/AkamaiOPEN-edgegrid-golang/v8/pkg/session"
	"github.com/akamai/terraform-provider-akamai/v6/pkg/common/ptr"
	"github.com/akamai/terraform-provider-akamai/v6/pkg/common/tf"
	"github.com/akamai/terraform-provider-akamai/v6/pkg/logger"
	"github.com/akamai/terraform-provider-akamai/v6/pkg/meta"
	"github.com/hashicorp/go-cty/cty"
	"github.com/hashicorp/terraform-plugin-sdk/v2/diag"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/schema"
	"golang.org/x/exp/slices"
)

func resourceGTMv1Property() *schema.Resource {
	return &schema.Resource{
		CreateContext: resourceGTMv1PropertyCreate,
		ReadContext:   resourceGTMv1PropertyRead,
		UpdateContext: resourceGTMv1PropertyUpdate,
		DeleteContext: resourceGTMv1PropertyDelete,
		CustomizeDiff: customDiffGTMProperty,
		Importer: &schema.ResourceImporter{
			State: resourceGTMv1PropertyImport,
		},
		Schema: map[string]*schema.Schema{
			"domain": {
				Type:     schema.TypeString,
				Required: true,
			},
			"wait_on_complete": {
				Type:     schema.TypeBool,
				Optional: true,
				Default:  true,
			},
			"name": {
				Type:     schema.TypeString,
				Required: true,
				ForceNew: true,
			},
			"type": {
				Type:     schema.TypeString,
				Required: true,
			},
			"ipv6": {
				Type:     schema.TypeBool,
				Optional: true,
			},
			"score_aggregation_type": {
				Type:     schema.TypeString,
				Required: true,
			},
			"stickiness_bonus_percentage": {
				Type:     schema.TypeInt,
				Optional: true,
			},
			"stickiness_bonus_constant": {
				Type:     schema.TypeInt,
				Optional: true,
			},
			"health_threshold": {
				Type:     schema.TypeFloat,
				Optional: true,
			},
			"use_computed_targets": {
				Type:     schema.TypeBool,
				Optional: true,
			},
			"backup_ip": {
				Type:     schema.TypeString,
				Optional: true,
			},
			"balance_by_download_score": {
				Type:     schema.TypeBool,
				Optional: true,
			},
			"static_ttl": {
				// Deprecated. Leaving for backward config compatibility.
				Type:             schema.TypeInt,
				Optional:         true,
				ValidateDiagFunc: validateTTL,
			},
			"static_rr_set": {
				Type:     schema.TypeList,
				Optional: true,
				Elem: &schema.Resource{
					Schema: map[string]*schema.Schema{
						"type": {
							Type:     schema.TypeString,
							Optional: true,
						},
						"ttl": {
							Type:     schema.TypeInt,
							Optional: true,
						},
						"rdata": {
							Type:     schema.TypeList,
							Elem:     &schema.Schema{Type: schema.TypeString},
							Optional: true,
						},
					},
				},
			},
			"unreachable_threshold": {
				Type:     schema.TypeFloat,
				Optional: true,
			},
			"min_live_fraction": {
				Type:     schema.TypeFloat,
				Optional: true,
			},
			"health_multiplier": {
				Type:     schema.TypeFloat,
				Optional: true,
			},
			"dynamic_ttl": {
				Type:             schema.TypeInt,
				Optional:         true,
				ValidateDiagFunc: validateTTL,
				Default:          300,
			},
			"max_unreachable_penalty": {
				Type:     schema.TypeInt,
				Optional: true,
			},
			"map_name": {
				Type:     schema.TypeString,
				Optional: true,
			},
			"handout_limit": {
				Type:     schema.TypeInt,
				Required: true,
			},
			"handout_mode": {
				Type:     schema.TypeString,
				Required: true,
			},
			"failover_delay": {
				Type:     schema.TypeInt,
				Optional: true,
			},
			"backup_cname": {
				Type:     schema.TypeString,
				Optional: true,
			},
			"failback_delay": {
				Type:     schema.TypeInt,
				Optional: true,
			},
			"load_imbalance_percentage": {
				Type:     schema.TypeFloat,
				Optional: true,
			},
			"health_max": {
				Type:     schema.TypeFloat,
				Optional: true,
			},
			"cname": {
				Type:     schema.TypeString,
				Optional: true,
			},
			"comments": {
				Type:     schema.TypeString,
				Optional: true,
			},
			"ghost_demand_reporting": {
				Type:     schema.TypeBool,
				Optional: true,
			},
			"weighted_hash_bits_for_ipv4": {
				Type:     schema.TypeInt,
				Computed: true,
			},
			"weighted_hash_bits_for_ipv6": {
				Type:     schema.TypeInt,
				Computed: true,
			},
			"traffic_target": {
				Type:             schema.TypeList,
				Optional:         true,
				MinItems:         1,
				DiffSuppressFunc: trafficTargetDiffSuppress,
				Elem: &schema.Resource{
					Schema: map[string]*schema.Schema{
						"datacenter_id": {
							Type:     schema.TypeInt,
							Optional: true,
						},
						"enabled": {
							Type:     schema.TypeBool,
							Optional: true,
						},
						"weight": {
							Type:     schema.TypeFloat,
							Optional: true,
						},
						"servers": {
							Type: schema.TypeSet,
							Elem: &schema.Schema{
								Type: schema.TypeString,
							},
							Optional: true,
						},
						"name": {
							Type:       schema.TypeString,
							Optional:   true,
							Deprecated: "The attribute `name` has been deprecated. Any reads or writes on this attribute are ignored",
						},
						"handout_cname": {
							Type:     schema.TypeString,
							Optional: true,
						},
						"precedence": {
							Type:     schema.TypeInt,
							Optional: true,
						},
					},
				},
			},
			"liveness_test": {
				Type:     schema.TypeList,
				Optional: true,
				MinItems: 1,
				Elem: &schema.Resource{
					Schema: map[string]*schema.Schema{
						"name": {
							Type:     schema.TypeString,
							Required: true,
						},
						"error_penalty": {
							Type:     schema.TypeFloat,
							Optional: true,
						},
						"peer_certificate_verification": {
							Type:     schema.TypeBool,
							Optional: true,
							Default:  true,
						},
						"test_interval": {
							Type:     schema.TypeInt,
							Required: true,
						},
						"test_object": {
							Type:     schema.TypeString,
							Optional: true,
						},
						"request_string": {
							Type:     schema.TypeString,
							Optional: true,
						},
						"response_string": {
							Type:     schema.TypeString,
							Optional: true,
						},
						"http_error3xx": {
							Type:     schema.TypeBool,
							Optional: true,
						},
						"http_error4xx": {
							Type:     schema.TypeBool,
							Optional: true,
						},
						"http_error5xx": {
							Type:     schema.TypeBool,
							Optional: true,
						},
						"http_method": {
							Type:     schema.TypeString,
							Optional: true,
						},
						"http_request_body": {
							Type:     schema.TypeString,
							Optional: true,
						},
						"disabled": {
							Type:     schema.TypeBool,
							Optional: true,
						},
						"test_object_protocol": {
							Type:     schema.TypeString,
							Required: true,
						},
						"test_object_password": {
							Type:     schema.TypeString,
							Optional: true,
						},
						"test_object_port": {
							Type:     schema.TypeInt,
							Optional: true,
							Default:  80,
						},
						"ssl_client_private_key": {
							Type:     schema.TypeString,
							Optional: true,
						},
						"ssl_client_certificate": {
							Type:     schema.TypeString,
							Optional: true,
						},
						"disable_nonstandard_port_warning": {
							Type:     schema.TypeBool,
							Optional: true,
						},
						"http_header": {
							Type:     schema.TypeList,
							Optional: true,
							Elem: &schema.Resource{
								Schema: map[string]*schema.Schema{
									"name": {
										Type:     schema.TypeString,
										Optional: true,
									},
									"value": {
										Type:     schema.TypeString,
										Optional: true,
									},
								},
							},
						},
						"test_object_username": {
							Type:     schema.TypeString,
							Optional: true,
						},
						"test_timeout": {
							Type:     schema.TypeFloat,
							Required: true,
						},
						"timeout_penalty": {
							Type:     schema.TypeFloat,
							Optional: true,
						},
						"answers_required": {
							Type:     schema.TypeBool,
							Optional: true,
						},
						"resource_type": {
							Type:     schema.TypeString,
							Optional: true,
						},
						"recursion_requested": {
							Type:     schema.TypeBool,
							Optional: true,
						},
						"pre_2023_security_posture": {
							Type:        schema.TypeBool,
							Optional:    true,
							Description: "Whether to enable backwards compatibility for liveness endpoints that use older TLS protocols",
						},
						"alternate_ca_certificates": {
							Type:     schema.TypeList,
							Optional: true,
							Elem: &schema.Schema{
								Type: schema.TypeString,
							},
						},
					},
				},
			},
		},
	}
}

// utility func to parse Terraform resource string id
func parseResourceStringID(id string) (string, string, error) {

	parts := strings.SplitN(id, ":", 2)
	if len(parts) != 2 || parts[0] == "" || parts[1] == "" {
		return "", "", fmt.Errorf("invalid resource ID: %s", id)
	}

	return parts[0], parts[1], nil

}

func validatePrecedence(trafficTargets []interface{}) error {
	precedenceCounter := map[int]int{}
	minPrecedence := 256

	for _, itemRaw := range trafficTargets {
		item, ok := itemRaw.(map[string]interface{})
		if !ok {
			return fmt.Errorf("could not cast the value of type %T to map[string]interface{}", itemRaw)
		}
		precedence, ok := item["precedence"].(int)
		if !ok {
			return fmt.Errorf("could not cast the value of type %T to int", item["precedence"])
		}
		precedenceCounter[precedence]++

		if precedence < minPrecedence {
			minPrecedence = precedence
		}
	}

	if precedenceCounter[minPrecedence] > 1 {
		return fmt.Errorf("property cannot have multiple primary traffic targets (targets with lowest precedence)")
	}

	return nil
}

func validateTrafficTargets(d *schema.ResourceDiff) error {
	propertyTypeRaw := d.Get("type")
	propertyType, ok := propertyTypeRaw.(string)
	if !ok {
		return fmt.Errorf("could not cast the value of type %T to string", propertyType)
	}
	if propertyType == "ranked-failover" {
		trafficTargetsRaw := d.Get("traffic_target")
		trafficTargets, ok := trafficTargetsRaw.([]interface{})
		if !ok {
			return fmt.Errorf("could not cast the value of type %T to []interface{}", trafficTargets)
		}
		if len(trafficTargets) == 0 {
			return fmt.Errorf("at least one 'traffic_target' has to be defined and enabled")
		}
		if err := validatePrecedence(trafficTargets); err != nil {
			return err
		}
	}
	return nil
}

// customDiffGTMProperty performs additional logic to the resource as part of custom diff function
func customDiffGTMProperty(_ context.Context, d *schema.ResourceDiff, m interface{}) error {
	meta := meta.Must(m)
	logger := meta.Log("Akamai GTM", "customDiffGTMProperty")
	logger.Debug("customDiffGTMProperty")

	if err := validateTrafficTargets(d); err != nil {
		return err
	}

	livenessTestRaw, ok := d.GetOkExists("liveness_test")
	if !ok {
		return nil
	}

	livenessTest, ok := livenessTestRaw.([]interface{})
	if !ok {
		return fmt.Errorf("could not cast the value of type %T to []interface{}", livenessTest)
	}

	requiredWith := []string{"HTTP", "HTTPS", "FTP"}

	for _, itemRaw := range livenessTest {
		item, ok := itemRaw.(map[string]interface{})
		if !ok {
			return fmt.Errorf("could not cast the value of type %T to map[string]interface{}", item)
		}
		testObjectProtocol, ok := item["test_object_protocol"].(string)
		if !ok {
			return fmt.Errorf("could not cast the value of type %T to string", testObjectProtocol)
		}
		testObject, ok := item["test_object"].(string)
		if !ok {
			return fmt.Errorf("could not cast the value of type %T to string", testObject)
		}

		if slices.Contains(requiredWith, testObjectProtocol) && testObject == "" {
			return fmt.Errorf("attribute 'test_object' is required when 'test_object_protocol' is set to 'HTTP', 'HTTPS' or 'FTP'")
		}
	}

	return nil
}

// validateTTL is a SchemaValidateDiagFunc to validate dynamic_ttl and static_ttl.
func validateTTL(v interface{}, path cty.Path) diag.Diagnostics {
	schemaFieldName, err := tf.GetSchemaFieldNameFromPath(path)
	if err != nil {
		return diag.FromErr(err)
	}

	if schemaFieldName == "static_ttl" {
		return diag.Diagnostics{
			diag.Diagnostic{
				Severity: diag.Warning,
				Summary:  "static_ttl is deprecated and will be ignored. Use static_rr_sets to apply static ttls to records",
			},
		}
	}
	value, ok := v.(int)
	if !ok {
		return diag.Errorf("%s validation failed to read field attribute", schemaFieldName)
	}
	if value < 30 || value > 3600 {
		return diag.Errorf("%s value must be between 30 and 3600", schemaFieldName)
	}
	return nil
}

// Create a new GTM Property
func resourceGTMv1PropertyCreate(ctx context.Context, d *schema.ResourceData, m interface{}) diag.Diagnostics {
	meta := meta.Must(m)
	logger := meta.Log("Akamai GTM", "resourceGTMv1PropertyCreate")
	// create a context with logging for api calls
	ctx = session.ContextWithOptions(
		ctx,
		session.WithContextLog(logger),
	)

	domain, err := tf.GetStringValue("domain", d)
	if err != nil {
		return diag.FromErr(err)
	}

	propertyName, err := tf.GetStringValue("name", d)
	if err != nil {
		return diag.FromErr(err)
	}

	propertyType, err := tf.GetStringValue("type", d)
	if err != nil {
		return diag.FromErr(err)
	}
	// Static properties cannot have traffic_targets. Non Static properties must
	traffTargList, err := tf.GetInterfaceArrayValue("traffic_target", d)
	if strings.ToUpper(propertyType) == "STATIC" && err == nil && (traffTargList != nil && len(traffTargList) > 0) {
		logger.Errorf("Property %s Create failed. Static property cannot have traffic targets", propertyName)
		return diag.Errorf("property Create failed. Static property cannot have traffic targets")
	}
	if strings.ToUpper(propertyType) != "STATIC" && (err != nil || (traffTargList == nil || len(traffTargList) < 1)) {
		logger.Errorf("Property %s Create failed. Property must have one or more traffic targets", propertyName)
		return diag.Errorf("property Create failed. Property must have one or more traffic targets")
	}

	logger.Infof("Creating property [%s] in domain [%s]", propertyName, domain)
	newProp, err := populateNewPropertyObject(d, m)
	if err != nil {
		return diag.FromErr(err)
	}
	logger.Debugf("Proposed New Property: [%v]", newProp)
	cStatus, err := Client(meta).CreateProperty(ctx, newProp, domain)
	if err != nil {
		logger.Errorf("Property Create failed: %s", err.Error())
		return diag.Errorf("property Create failed: %s", err.Error())
	}
	logger.Debugf("Property Create status: %v", cStatus.Status)

	if cStatus.Status.PropagationStatus == "DENIED" {
		logger.Errorf(cStatus.Status.Message)
		return diag.FromErr(fmt.Errorf(cStatus.Status.Message))
	}

	waitOnComplete, err := tf.GetBoolValue("wait_on_complete", d)
	if err != nil {
		return diag.FromErr(err)
	}

	if waitOnComplete {
		done, err := waitForCompletion(ctx, domain, m)
		if done {
			logger.Infof("Property Create completed")
		} else {
			if err == nil {
				logger.Infof("Property Create pending")
			} else {
				logger.Errorf("Property Create failed [%s]", err.Error())
				return diag.Errorf("property Create failed [%s]", err.Error())
			}
		}
	}

	// Give terraform the ID. Format domain::property
	propertyResourceID := fmt.Sprintf("%s:%s", domain, cStatus.Resource.Name)
	logger.Debugf("Generated Property resource ID: %s", propertyResourceID)
	d.SetId(propertyResourceID)
	return resourceGTMv1PropertyRead(ctx, d, m)

}

// Only ever save data from the tf config in the tf state file, to help with
// api issues. See func unmarshalResourceData for more info.
func resourceGTMv1PropertyRead(ctx context.Context, d *schema.ResourceData, m interface{}) diag.Diagnostics {
	meta := meta.Must(m)
	logger := meta.Log("Akamai GTM", "resourceGTMv1PropertyRead")
	// create a context with logging for api calls
	ctx = session.ContextWithOptions(
		ctx,
		session.WithContextLog(logger),
	)

	logger.Debugf("Reading Property: %s", d.Id())
	// retrieve the property and domain
	domain, property, err := parseResourceStringID(d.Id())
	if err != nil {
		return diag.FromErr(err)
	}
	prop, err := Client(meta).GetProperty(ctx, property, domain)
	if errors.Is(err, gtm.ErrNotFound) {
		d.SetId("")
		return nil
	}
	if err != nil {
		logger.Errorf("Property Read failed: %s", err.Error())
		return diag.Errorf("property Read failed: %s", err.Error())
	}
	populateTerraformPropertyState(d, prop, m)
	logger.Debugf("READ %v", prop)
	return nil
}

// Update GTM Property
func resourceGTMv1PropertyUpdate(ctx context.Context, d *schema.ResourceData, m interface{}) diag.Diagnostics {
	meta := meta.Must(m)
	logger := meta.Log("Akamai GTM", "resourceGTMv1PropertyUpdate")
	// create a context with logging for api calls
	ctx = session.ContextWithOptions(
		ctx,
		session.WithContextLog(logger),
	)

	logger.Debugf("Updating Property: %s", d.Id())
	// pull domain and property out of resource id
	domain, property, err := parseResourceStringID(d.Id())
	if err != nil {
		return diag.FromErr(err)
	}
	// Get existing property
	existProp, err := Client(meta).GetProperty(ctx, property, domain)
	if err != nil {
		logger.Errorf("Property Update failed: %s", err.Error())
		return diag.FromErr(err)
	}
	logger.Debugf("Updating Property BEFORE: %v", existProp)
	err = populatePropertyObject(d, existProp, m)
	if err != nil {
		return diag.FromErr(err)
	}
	logger.Debugf("Updating Property PROPOSED: %v", existProp)
	uStat, err := Client(meta).UpdateProperty(ctx, existProp, domain)
	if err != nil {
		logger.Errorf("Property Update failed: %s", err.Error())
		return diag.Errorf("property Update failed: %s", err.Error())
	}
	logger.Debugf("Property Update  status: %v", uStat)
	if uStat.PropagationStatus == "DENIED" {
		logger.Debugf(uStat.Message)
		return diag.FromErr(fmt.Errorf(uStat.Message))
	}

	waitOnComplete, err := tf.GetBoolValue("wait_on_complete", d)
	if err != nil {
		return diag.FromErr(err)
	}

	if waitOnComplete {
		done, err := waitForCompletion(ctx, domain, m)
		if done {
			logger.Infof("Property Update completed")
		} else {
			if err == nil {
				logger.Infof("Property Update pending")
			} else {
				logger.Errorf("Property Update failed [%s]", err.Error())
				return diag.Errorf("property Update failed [%s]", err.Error())
			}
		}
	}

	return resourceGTMv1PropertyRead(ctx, d, m)
}

// Import GTM Property.
func resourceGTMv1PropertyImport(d *schema.ResourceData, m interface{}) ([]*schema.ResourceData, error) {
	meta := meta.Must(m)
	logger := meta.Log("Akamai GTM", "resourceGTMv1PropertyImport")
	// create a context with logging for api calls
	ctx := context.Background()
	ctx = session.ContextWithOptions(
		ctx,
		session.WithContextLog(logger),
	)
	logger.Infof("Property [%s] Import", d.Id())
	// pull domain and property out of resource id
	domain, property, err := parseResourceStringID(d.Id())
	if err != nil {
		return []*schema.ResourceData{d}, err
	}
	prop, err := Client(meta).GetProperty(ctx, property, domain)
	if err != nil {
		return nil, err
	}
	if err := d.Set("domain", domain); err != nil {
		return nil, err
	}
	if err := d.Set("wait_on_complete", true); err != nil {
		return nil, err
	}
	populateTerraformPropertyState(d, prop, m)

	// use same Id as passed in
	logger.Infof("Property [%s] [%s] Imported", d.Id(), d.Get("name"))
	return []*schema.ResourceData{d}, nil
}

// Delete GTM Property.
func resourceGTMv1PropertyDelete(ctx context.Context, d *schema.ResourceData, m interface{}) diag.Diagnostics {
	meta := meta.Must(m)
	logger := meta.Log("Akamai GTM", "resourceGTMv1PropertyDelete")
	// create a context with logging for api calls
	ctx = session.ContextWithOptions(
		ctx,
		session.WithContextLog(logger),
	)

	logger.Debugf("Deleting Property: %s", d.Id())
	// Get existing property
	domain, property, err := parseResourceStringID(d.Id())
	if err != nil {
		return diag.FromErr(err)
	}
	existProp, err := Client(meta).GetProperty(ctx, property, domain)
	if err != nil {
		logger.Errorf("Property Delete failed: %s", err.Error())
		return diag.Errorf("property Delete failed: %s", err.Error())
	}
	logger.Debugf("Deleting Property: %v", existProp)
	uStat, err := Client(meta).DeleteProperty(ctx, existProp, domain)
	if err != nil {
		logger.Errorf("Property Delete failed: %s", err.Error())
		return diag.Errorf("property Delete failed: %s", err.Error())
	}
	logger.Debugf("Property Delete status: %v", uStat)
	if uStat.PropagationStatus == "DENIED" {
		logger.Errorf(uStat.Message)
		return diag.FromErr(fmt.Errorf(uStat.Message))
	}

	waitOnComplete, err := tf.GetBoolValue("wait_on_complete", d)
	if err != nil {
		return diag.FromErr(err)
	}

	if waitOnComplete {
		done, err := waitForCompletion(ctx, domain, m)
		if done {
			logger.Infof("Property Delete completed")
		} else {
			if err == nil {
				logger.Infof("Property Delete pending")
			} else {
				logger.Errorf("Property Delete failed [%s]", err.Error())
				return diag.Errorf("property Delete failed [%s]", err.Error())
			}
		}
	}

	// if successful ....
	d.SetId("")
	return nil
}

// nolint:gocyclo
// Populate existing property object from resource data
func populatePropertyObject(d *schema.ResourceData, prop *gtm.Property, m interface{}) error {
	meta := meta.Must(m)
	logger := meta.Log("Akamai GTM", "populatePropertyObject")

	vstr, err := tf.GetStringValue("name", d)
	if err == nil {
		prop.Name = vstr
	}
	ptype, err := tf.GetStringValue("type", d)
	if err == nil {
		prop.Type = ptype
	}
	vstr, err = tf.GetStringValue("score_aggregation_type", d)
	if err == nil {
		prop.ScoreAggregationType = vstr
	}
	vint, err := tf.GetIntValue("stickiness_bonus_percentage", d)
	if err == nil || d.HasChange("stickiness_bonus_percentage") {
		prop.StickinessBonusPercentage = vint
	}
	if err != nil && !errors.Is(err, tf.ErrNotFound) {
		logger.Errorf("populateResourceObject() stickiness_bonus_percentage failed: %v", err.Error())
		return fmt.Errorf("property Object could not be populated: %v", err.Error())
	}

	vint, err = tf.GetIntValue("stickiness_bonus_constant", d)
	if err == nil || d.HasChange("stickiness_bonus_constant") {
		prop.StickinessBonusConstant = vint
	}
	if err != nil && !errors.Is(err, tf.ErrNotFound) {
		logger.Errorf("populateResourceObject() stickiness_bonus_constant failed: %v", err.Error())
		return fmt.Errorf("property Object could not be populated: %v", err.Error())
	}

	vfloat, err := tf.GetFloat64Value("health_threshold", d)
	if err == nil || d.HasChange("health_threshold") {
		prop.HealthThreshold = vfloat
	}
	if err != nil && !errors.Is(err, tf.ErrNotFound) {
		logger.Errorf("populateResourceObject() health_threshold failed: %v", err.Error())
		return fmt.Errorf("property Object could not be populated: %v", err.Error())
	}

	if ipv6, err := tf.GetBoolValue("ipv6", d); err == nil {
		prop.IPv6 = ipv6
	}
	if uct, err := tf.GetBoolValue("use_computed_targets", d); err == nil {
		prop.UseComputedTargets = uct
	}

	vstr, err = tf.GetStringValue("backup_ip", d)
	if err == nil || d.HasChange("backup_ip") {
		prop.BackupIP = vstr
	}
	if err != nil && !errors.Is(err, tf.ErrNotFound) {
		logger.Errorf("populateResourceObject() backup_ip failed: %v", err.Error())
		return fmt.Errorf("property Object could not be populated: %v", err.Error())
	}

	if bbds, err := tf.GetBoolValue("balance_by_download_score", d); err == nil {
		prop.BalanceByDownloadScore = bbds
	}

	vfloat, err = tf.GetFloat64Value("unreachable_threshold", d)
	if err == nil || d.HasChange("unreachable_threshold") {
		prop.UnreachableThreshold = vfloat
	}
	if err != nil && !errors.Is(err, tf.ErrNotFound) {
		logger.Errorf("populateResourceObject() unreachable_threshold failed: %v", err.Error())
		return fmt.Errorf("property Object could not be populated: %v", err.Error())
	}

	vfloat, err = tf.GetFloat64Value("min_live_fraction", d)
	if err == nil || d.HasChange("min_live_fraction") {
		prop.MinLiveFraction = vfloat
	}
	if err != nil && !errors.Is(err, tf.ErrNotFound) {
		logger.Errorf("populateResourceObject() min_live_fraction failed: %v", err.Error())
		return fmt.Errorf("property Object could not be populated: %v", err.Error())
	}

	vfloat, err = tf.GetFloat64Value("health_multiplier", d)
	if err == nil || d.HasChange("health_multiplier") {
		prop.HealthMultiplier = vfloat
	}
	if err != nil && !errors.Is(err, tf.ErrNotFound) {
		logger.Errorf("populateResourceObject() health_multiplier failed: %v", err.Error())
		return fmt.Errorf("property Object could not be populated: %v", err.Error())
	}

	vint, err = tf.GetIntValue("dynamic_ttl", d)
	if err == nil || d.HasChange("dynamic_ttl") {
		prop.DynamicTTL = vint
	}
	if err != nil && !errors.Is(err, tf.ErrNotFound) {
		logger.Errorf("populateResourceObject() dynamic_ttl failed: %v", err.Error())
		return fmt.Errorf("property Object could not be populated: %v", err.Error())
	}

	vint, err = tf.GetIntValue("max_unreachable_penalty", d)
	if err == nil || d.HasChange("max_unreachable_penalty") {
		prop.MaxUnreachablePenalty = vint
	}
	if err != nil && !errors.Is(err, tf.ErrNotFound) {
		logger.Errorf("populateResourceObject() max_unreachable_penalty failed: %v", err.Error())
		return fmt.Errorf("property Object could not be populated: %v", err.Error())
	}

	vstr, err = tf.GetStringValue("map_name", d)
	if err == nil || d.HasChange("map_name") {
		prop.MapName = vstr
	}
	if err != nil && !errors.Is(err, tf.ErrNotFound) {
		logger.Errorf("populateResourceObject() map_name failed: %v", err.Error())
		return fmt.Errorf("property Object could not be populated: %v", err.Error())
	}

	if vint, err = tf.GetIntValue("handout_limit", d); err == nil || d.HasChange("handout_limit") {
		prop.HandoutLimit = vint
	}
	if vstr, err = tf.GetStringValue("handout_mode", d); err == nil {
		prop.HandoutMode = vstr
	}

	vfloat, err = tf.GetFloat64Value("load_imbalance_percentage", d)
	if err == nil || d.HasChange("load_imbalance_percentage") {
		prop.LoadImbalancePercentage = vfloat
	}
	if err != nil && !errors.Is(err, tf.ErrNotFound) {
		logger.Errorf("populateResourceObject() load_imbalance_percentage failed: %v", err.Error())
		return fmt.Errorf("property Object could not be populated: %v", err.Error())
	}

	vint, err = tf.GetIntValue("failover_delay", d)
	if err == nil || d.HasChange("failover_delay") {
		prop.FailoverDelay = vint
	}
	if err != nil && !errors.Is(err, tf.ErrNotFound) {
		logger.Errorf("populateResourceObject() failover_delay failed: %v", err.Error())
		return fmt.Errorf("property Object could not be populated: %v", err.Error())
	}

	vstr, err = tf.GetStringValue("backup_cname", d)
	if err == nil || d.HasChange("backup_cname") {
		prop.BackupCName = vstr
	}
	if err != nil && !errors.Is(err, tf.ErrNotFound) {
		logger.Errorf("populateResourceObject() backup_cname failed: %v", err.Error())
		return fmt.Errorf("property Object could not be populated: %v", err.Error())
	}

	vint, err = tf.GetIntValue("failback_delay", d)
	if err == nil || d.HasChange("failback_delay") {
		prop.FailbackDelay = vint
	}
	if err != nil && !errors.Is(err, tf.ErrNotFound) {
		logger.Errorf("populateResourceObject() failback_delay failed: %v", err.Error())
		return fmt.Errorf("property Object could not be populated: %v", err.Error())
	}

	vfloat, err = tf.GetFloat64Value("health_max", d)
	if err == nil || d.HasChange("health_max") {
		prop.HealthMax = vfloat
	}
	if err != nil && !errors.Is(err, tf.ErrNotFound) {
		logger.Errorf("populateResourceObject() health_max failed: %v", err.Error())
		return fmt.Errorf("property Object could not be populated: %v", err.Error())
	}

	if gdr, err := tf.GetBoolValue("ghost_demand_reporting", d); err == nil {
		prop.GhostDemandReporting = gdr
	}
	if vint, err = tf.GetIntValue("weighted_hash_bits_for_ipv4", d); err == nil {
		prop.WeightedHashBitsForIPv4 = vint
	}
	if vint, err = tf.GetIntValue("weighted_hash_bits_for_ipv6", d); err == nil {
		prop.WeightedHashBitsForIPv6 = vint
	}

	vstr, err = tf.GetStringValue("cname", d)
	if err == nil || d.HasChange("cname") {
		prop.CName = vstr
	}
	if err != nil && !errors.Is(err, tf.ErrNotFound) {
		logger.Errorf("populateResourceObject() cname failed: %v", err.Error())
		return fmt.Errorf("property Object could not be populated: %v", err.Error())
	}

	vstr, err = tf.GetStringValue("comments", d)
	if err == nil || d.HasChange("comments") {
		prop.Comments = vstr
	}
	if err != nil && !errors.Is(err, tf.ErrNotFound) {
		logger.Errorf("populateResourceObject() comments failed: %v", err.Error())
		return fmt.Errorf("property Object could not be populated: %v", err.Error())
	}

	if strings.ToUpper(ptype) != "STATIC" {
		populateTrafficTargetObject(d, prop, m)
	}
	populateStaticRRSetObject(d, prop)
	populateLivenessTestObject(d, prop)

	return nil
}

// Create and populate a new property object from resource data
func populateNewPropertyObject(d *schema.ResourceData, m interface{}) (*gtm.Property, error) {

	name, err := tf.GetStringValue("name", d)
	if err != nil {
		return nil, err
	}
	propObj := &gtm.Property{
		Name:           name,
		TrafficTargets: make([]*gtm.TrafficTarget, 0),
		LivenessTests:  make([]*gtm.LivenessTest, 0),
	}
	err = populatePropertyObject(d, propObj, m)

	return propObj, err

}

// Populate Terraform state from provided Property object
func populateTerraformPropertyState(d *schema.ResourceData, prop *gtm.Property, m interface{}) {
	meta := meta.Must(m)
	logger := meta.Log("Akamai GTM", "populateTerraformPropertyState")

	for stateKey, stateValue := range map[string]interface{}{
		"name":                        prop.Name,
		"type":                        prop.Type,
		"ipv6":                        prop.IPv6,
		"score_aggregation_type":      prop.ScoreAggregationType,
		"stickiness_bonus_percentage": prop.StickinessBonusPercentage,
		"stickiness_bonus_constant":   prop.StickinessBonusConstant,
		"health_threshold":            prop.HealthThreshold,
		"use_computed_targets":        prop.UseComputedTargets,
		"backup_ip":                   prop.BackupIP,
		"balance_by_download_score":   prop.BalanceByDownloadScore,
		"unreachable_threshold":       prop.UnreachableThreshold,
		"min_live_fraction":           prop.MinLiveFraction,
		"health_multiplier":           prop.HealthMultiplier,
		"dynamic_ttl":                 prop.DynamicTTL,
		"max_unreachable_penalty":     prop.MaxUnreachablePenalty,
		"map_name":                    prop.MapName,
		"handout_limit":               prop.HandoutLimit,
		"handout_mode":                prop.HandoutMode,
		"load_imbalance_percentage":   prop.LoadImbalancePercentage,
		"failover_delay":              prop.FailoverDelay,
		"backup_cname":                prop.BackupCName,
		"failback_delay":              prop.FailbackDelay,
		"health_max":                  prop.HealthMax,
		"ghost_demand_reporting":      prop.GhostDemandReporting,
		"weighted_hash_bits_for_ipv4": prop.WeightedHashBitsForIPv4,
		"weighted_hash_bits_for_ipv6": prop.WeightedHashBitsForIPv6,
		"cname":                       prop.CName,
		"comments":                    prop.Comments,
	} {
		// walk thru all state elements
		if stateKey == "dynamic_ttl" && stateValue == 0 {
			// ttl value is not set; null -> 0
			continue
		}
		err := d.Set(stateKey, stateValue)
		if err != nil {
			logger.Errorf("Invalid configuration: %s", err.Error())
		}
	}
	if strings.ToUpper(prop.Type) != "STATIC" {
		if err := populateTerraformTrafficTargetState(d, prop, m); err != nil {
			logger.Errorf("Error setting traffic target: %s", err)
		}
	}
	if err := populateTerraformStaticRRSetState(d, prop, m); err != nil {
		logger.Errorf("error setting static RR set: %s", err)
	}
	if err := populateTerraformLivenessTestState(d, prop, m); err != nil {
		logger.Errorf("error setting liveness test: %s", err)
	}

}

// create and populate GTM Property TrafficTargets object
func populateTrafficTargetObject(d *schema.ResourceData, prop *gtm.Property, m interface{}) {
	meta := meta.Must(m)
	logger := meta.Log("Akamai GTM", "populateTrafficTargetObject")

	// pull apart List
	traffTargList, err := tf.GetInterfaceArrayValue("traffic_target", d)
	if err == nil {
		trafficObjList := make([]*gtm.TrafficTarget, len(traffTargList)) // create new object list
		for i, v := range traffTargList {
			ttMap := v.(map[string]interface{})
			trafficTarget := &gtm.TrafficTarget{}
			trafficTarget.DatacenterID = ttMap["datacenter_id"].(int)
			trafficTarget.Precedence = ptr.To(ttMap["precedence"].(int))
			trafficTarget.Enabled = ttMap["enabled"].(bool)
			trafficTarget.Weight = ttMap["weight"].(float64)
			if ttMap["servers"] != nil {
				ls := make([]string, ttMap["servers"].(*schema.Set).Len())
				for i, sl := range ttMap["servers"].(*schema.Set).List() {
					ls[i] = sl.(string)
				}
				trafficTarget.Servers = ls
			}
			trafficTarget.HandoutCName = ttMap["handout_cname"].(string)
			trafficObjList[i] = trafficTarget
		}
		prop.TrafficTargets = trafficObjList
	} else {
		logger.Warnf("traffic_target not set in ResourceData: %s", err.Error())
	}
}

// create and populate Terraform traffic_targets schema
func populateTerraformTrafficTargetState(d *schema.ResourceData, prop *gtm.Property, m interface{}) error {
	meta := meta.Must(m)
	logger := meta.Log("Akamai GTM", "populateTerraformTrafficTargetState")

	objectInventory := make(map[int]*gtm.TrafficTarget, len(prop.TrafficTargets))
	if len(prop.TrafficTargets) > 0 {
		for _, aObj := range prop.TrafficTargets {
			objectInventory[aObj.DatacenterID] = aObj
		}
	}

	ttStateList, err := tf.GetInterfaceArrayValue("traffic_target", d)
	if err != nil && !errors.Is(err, tf.ErrNotFound) {
		return err
	}
	for _, ttMap := range ttStateList {
		tt := ttMap.(map[string]interface{})
		objIndex := tt["datacenter_id"].(int)
		ttObject := objectInventory[objIndex]
		if ttObject == nil {
			logger.Warnf("Property TrafficTarget %d NOT FOUND in returned GTM Object", tt["datacenter_id"])
			continue
		}
		tt["datacenter_id"] = ttObject.DatacenterID
		tt["enabled"] = ttObject.Enabled
		tt["weight"] = ttObject.Weight
		tt["handout_cname"] = ttObject.HandoutCName
		tt["servers"] = ttObject.Servers
		if ttObject.Precedence != nil {
			tt["precedence"] = *ttObject.Precedence
		}
		// remove object
		delete(objectInventory, objIndex)
	}
	if len(objectInventory) > 0 {
		// Objects not in the state yet. Add. Unfortunately, they not align with instance indices in the config
		for _, mttObj := range objectInventory {
			logger.Debugf("Property TrafficObject NEW State Object: %d", mttObj.DatacenterID)
			ttNew := map[string]interface{}{
				"datacenter_id": mttObj.DatacenterID,
				"enabled":       mttObj.Enabled,
				"weight":        mttObj.Weight,
				"handout_cname": mttObj.HandoutCName,
				"servers":       mttObj.Servers,
			}
			if mttObj.Precedence != nil {
				ttNew["precedence"] = mttObj.Precedence
			}
			ttStateList = append(ttStateList, ttNew)
		}
	}

	return d.Set("traffic_target", ttStateList)
}

// Populate existing static_rr_sets object from resource data
func populateStaticRRSetObject(d *schema.ResourceData, prop *gtm.Property) {

	// pull apart List
	staticSetList, err := tf.GetInterfaceArrayValue("static_rr_set", d)
	if err == nil {
		staticObjList := make([]*gtm.StaticRRSet, len(staticSetList)) // create new object list
		for i, v := range staticSetList {
			recMap := v.(map[string]interface{})
			record := &gtm.StaticRRSet{
				TTL:  recMap["ttl"].(int),
				Type: recMap["type"].(string),
			}
			if recMap["rdata"] != nil {
				rls := make([]string, len(recMap["rdata"].([]interface{})))
				for i, d := range recMap["rdata"].([]interface{}) {
					rls[i] = d.(string)
				}
				record.Rdata = rls
			}
			staticObjList[i] = record
		}
		prop.StaticRRSets = staticObjList
	}
}

// create and populate Terraform static_rr_sets schema
func populateTerraformStaticRRSetState(d *schema.ResourceData, prop *gtm.Property, m interface{}) error {
	meta := meta.Must(m)
	logger := meta.Log("Akamai GTM", "populateTerraformStaticRRSetState")

	objectInventory := make(map[string]*gtm.StaticRRSet, len(prop.StaticRRSets))
	if len(prop.StaticRRSets) > 0 {
		for _, aObj := range prop.StaticRRSets {
			objectInventory[aObj.Type] = aObj
		}
	}

	rrStateList, err := tf.GetInterfaceArrayValue("static_rr_set", d)
	if err != nil && !errors.Is(err, tf.ErrNotFound) {
		return err
	}
	for _, rrMap := range rrStateList {
		rr := rrMap.(map[string]interface{})
		objIndex := rr["type"].(string)
		rrObject := objectInventory[objIndex]
		if rrObject == nil {
			logger.Warnf("Property StaticRRSet %s NOT FOUND in returned GTM Object", rr["type"])
			continue
		}
		rr["type"] = rrObject.Type
		rr["ttl"] = rrObject.TTL
		rr["rdata"] = reconcileTerraformLists(rr["rdata"].([]interface{}), convertStringToInterfaceList(rrObject.Rdata, m), m)
		// remove object
		delete(objectInventory, objIndex)
	}
	if len(objectInventory) > 0 {
		logger.Debugf("Property StaticRRSet objects left...")
		// Objects not in the state yet. Add. Unfortunately, they not align with instance indices in the config
		for _, mrrObj := range objectInventory {
			rrNew := map[string]interface{}{
				"type":  mrrObj.Type,
				"ttl":   mrrObj.TTL,
				"rdata": mrrObj.Rdata,
			}
			rrStateList = append(rrStateList, rrNew)
		}
	}

	return d.Set("static_rr_set", rrStateList)
}

// Populate existing Liveness test  object from resource data
func populateLivenessTestObject(d *schema.ResourceData, prop *gtm.Property) {

	liveTestList, err := tf.GetInterfaceArrayValue("liveness_test", d)
	if err == nil {
		liveTestObjList := make([]*gtm.LivenessTest, len(liveTestList)) // create new object list
		for i, l := range liveTestList {
			v := l.(map[string]interface{})

			lt := &gtm.LivenessTest{
				Name:                          v["name"].(string),
				TestObjectProtocol:            v["test_object_protocol"].(string),
				TestInterval:                  v["test_interval"].(int),
				TestTimeout:                   float32(v["test_timeout"].(float64)),
				ErrorPenalty:                  v["error_penalty"].(float64),
				PeerCertificateVerification:   v["peer_certificate_verification"].(bool),
				TestObject:                    v["test_object"].(string),
				RequestString:                 v["request_string"].(string),
				ResponseString:                v["response_string"].(string),
				HTTPError3xx:                  v["http_error3xx"].(bool),
				HTTPError4xx:                  v["http_error4xx"].(bool),
				HTTPError5xx:                  v["http_error5xx"].(bool),
				Pre2023SecurityPosture:        v["pre_2023_security_posture"].(bool),
				Disabled:                      v["disabled"].(bool),
				TestObjectPassword:            v["test_object_password"].(string),
				TestObjectPort:                v["test_object_port"].(int),
				SSLClientPrivateKey:           v["ssl_client_private_key"].(string),
				SSLClientCertificate:          v["ssl_client_certificate"].(string),
				DisableNonstandardPortWarning: v["disable_nonstandard_port_warning"].(bool),
				TestObjectUsername:            v["test_object_username"].(string),
				TimeoutPenalty:                v["timeout_penalty"].(float64),
				AnswersRequired:               v["answers_required"].(bool),
				ResourceType:                  v["resource_type"].(string),
				RecursionRequested:            v["recursion_requested"].(bool),
			}
			if v["http_method"].(string) != "" {
				lt.HTTPMethod = ptr.To(v["http_method"].(string))
			}
			if v["http_request_body"].(string) != "" {
				lt.HTTPRequestBody = ptr.To(v["http_request_body"].(string))
			}
			httpHeaderList := v["http_header"].([]interface{})
			if httpHeaderList != nil {
				headerObjList := make([]*gtm.HTTPHeader, len(httpHeaderList)) // create new object list
				for i, h := range httpHeaderList {
					recMap := h.(map[string]interface{})
					record := &gtm.HTTPHeader{
						Name:  recMap["name"].(string),
						Value: recMap["value"].(string),
					}
					headerObjList[i] = record
				}
				lt.HTTPHeaders = headerObjList
			}
			alternateCACerts := v["alternate_ca_certificates"].([]interface{})
			if alternateCACerts != nil {
				var certList []string
				for _, certRaw := range alternateCACerts {
					certList = append(certList, certRaw.(string))
				}
				lt.AlternateCACertificates = certList
			}
			liveTestObjList[i] = lt
		}
		prop.LivenessTests = liveTestObjList
	}
}

// create and populate Terraform liveness_test schema
func populateTerraformLivenessTestState(d *schema.ResourceData, prop *gtm.Property, m interface{}) error {
	meta := meta.Must(m)
	logger := meta.Log("Akamai GTM", "populateTerraformLivenessTestState")

	objectInventory := make(map[string]*gtm.LivenessTest, len(prop.LivenessTests))
	if len(prop.LivenessTests) > 0 {
		for _, aObj := range prop.LivenessTests {
			objectInventory[aObj.Name] = aObj
		}
	}

	ltStateList, err := tf.GetInterfaceArrayValue("liveness_test", d)
	if err != nil && !errors.Is(err, tf.ErrNotFound) {
		return err
	}
	for _, ltMap := range ltStateList {
		lt := ltMap.(map[string]interface{})
		objIndex := lt["name"].(string)
		ltObject := objectInventory[objIndex]
		if ltObject == nil {
			logger.Warnf("Property LivenessTest  %s NOT FOUND in returned GTM Object", lt["name"])
			continue
		}
		lt["name"] = ltObject.Name
		lt["error_penalty"] = ltObject.ErrorPenalty
		lt["peer_certificate_verification"] = ltObject.PeerCertificateVerification
		lt["test_interval"] = ltObject.TestInterval
		lt["test_object"] = ltObject.TestObject
		lt["request_string"] = ltObject.RequestString
		lt["response_string"] = ltObject.ResponseString
		lt["http_error3xx"] = ltObject.HTTPError3xx
		lt["http_error4xx"] = ltObject.HTTPError4xx
		lt["http_error5xx"] = ltObject.HTTPError5xx
		lt["disabled"] = ltObject.Disabled
		lt["test_object_protocol"] = ltObject.TestObjectProtocol
		lt["test_object_password"] = ltObject.TestObjectPassword
		lt["test_object_port"] = ltObject.TestObjectPort
		lt["ssl_client_private_key"] = ltObject.SSLClientPrivateKey
		lt["ssl_client_certificate"] = ltObject.SSLClientCertificate
		lt["disable_nonstandard_port_warning"] = ltObject.DisableNonstandardPortWarning
		lt["test_object_username"] = ltObject.TestObjectUsername
		lt["test_timeout"] = ltObject.TestTimeout
		lt["timeout_penalty"] = ltObject.TimeoutPenalty
		lt["answers_required"] = ltObject.AnswersRequired
		lt["resource_type"] = ltObject.ResourceType
		lt["recursion_requested"] = ltObject.RecursionRequested
		lt["pre_2023_security_posture"] = ltObject.Pre2023SecurityPosture
		if ltObject.HTTPMethod != nil {
			lt["http_method"] = ltObject.HTTPMethod
		}
		if ltObject.HTTPRequestBody != nil {
			lt["http_request_body"] = ltObject.HTTPRequestBody
		}
		if ltObject.AlternateCACertificates != nil {
			lt["alternate_ca_certificates"] = ltObject.AlternateCACertificates
		}
		httpHeaderListNew := make([]interface{}, len(ltObject.HTTPHeaders))
		for i, r := range ltObject.HTTPHeaders {
			httpHeaderNew := map[string]interface{}{
				"name":  r.Name,
				"value": r.Value,
			}
			httpHeaderListNew[i] = httpHeaderNew
		}
		lt["http_header"] = httpHeaderListNew
		// remove object
		delete(objectInventory, objIndex)
	}
	if len(objectInventory) > 0 {
		logger.Debugf("Property LivenessTest objects left...")
		// Objects not in the state yet. Add. Unfortunately, they not align with instance indices in the config
		for _, l := range objectInventory {
			ltNew := map[string]interface{}{
				"name":                             l.Name,
				"error_penalty":                    l.ErrorPenalty,
				"peer_certificate_verification":    l.PeerCertificateVerification,
				"test_interval":                    l.TestInterval,
				"test_object":                      l.TestObject,
				"request_string":                   l.RequestString,
				"response_string":                  l.ResponseString,
				"http_error3xx":                    l.HTTPError3xx,
				"http_error4xx":                    l.HTTPError4xx,
				"http_error5xx":                    l.HTTPError5xx,
				"disabled":                         l.Disabled,
				"test_object_protocol":             l.TestObjectProtocol,
				"test_object_password":             l.TestObjectPassword,
				"test_object_port":                 l.TestObjectPort,
				"ssl_client_private_key":           l.SSLClientPrivateKey,
				"ssl_client_certificate":           l.SSLClientCertificate,
				"disable_nonstandard_port_warning": l.DisableNonstandardPortWarning,
				"test_object_username":             l.TestObjectUsername,
				"test_timeout":                     l.TestTimeout,
				"timeout_penalty":                  l.TimeoutPenalty,
				"answers_required":                 l.AnswersRequired,
				"resource_type":                    l.ResourceType,
				"recursion_requested":              l.RecursionRequested,
				"pre_2023_security_posture":        l.Pre2023SecurityPosture,
			}
			if l.HTTPMethod != nil {
				ltNew["http_method"] = l.HTTPMethod
			}
			if l.HTTPRequestBody != nil {
				ltNew["http_request_body"] = l.HTTPRequestBody
			}
			if l.AlternateCACertificates != nil {
				ltNew["alternate_ca_certificates"] = l.AlternateCACertificates
			}
			httpHeaderListNew := make([]interface{}, len(l.HTTPHeaders))
			for i, r := range l.HTTPHeaders {
				httpHeaderNew := map[string]interface{}{
					"name":  r.Name,
					"value": r.Value,
				}
				httpHeaderListNew[i] = httpHeaderNew
			}
			ltNew["http_header"] = httpHeaderListNew
			ltStateList = append(ltStateList, ltNew)
		}
	}

	return d.Set("liveness_test", ltStateList)
}

func convertStringToInterfaceList(stringList []string, m interface{}) []interface{} {
	meta := meta.Must(m)
	logger := meta.Log("Akamai GTMv1", "convertStringToInterfaceList")

	logger.Debugf("String List: %v", stringList)
	retList := make([]interface{}, 0, len(stringList))
	for _, v := range stringList {
		retList = append(retList, v)
	}

	return retList

}

// Util method to reconcile list configs. Type agnostic. Goal: maintain order of tf list config
func reconcileTerraformLists(terraList []interface{}, newList []interface{}, m interface{}) []interface{} {
	meta := meta.Must(m)
	logger := meta.Log("Akamai GTMv1", "reconcileTerraformLists")

	logger.Debugf("Existing Terra List: %v", terraList)
	logger.Debugf("Read List: %v", newList)
	newMap := make(map[string]interface{}, len(newList))
	updatedList := make([]interface{}, 0, len(newList))
	for _, newelem := range newList {
		newMap[fmt.Sprintf("%v", newelem)] = newelem
	}
	// walk existing terra list and check new.
	for _, v := range terraList {
		vindex := fmt.Sprintf("%v", v)
		if _, ok := newMap[vindex]; ok {
			updatedList = append(updatedList, v)
			delete(newMap, vindex)
		}
	}
	for _, newVal := range newMap {
		updatedList = append(updatedList, newVal)
	}

	logger.Debugf("Updated Terra List: %v", updatedList)
	return updatedList

}

func trafficTargetDiffSuppress(_, _, _ string, d *schema.ResourceData) bool {
	logger := logger.Get("Akamai GTM", "trafficTargetDiffSuppress")
	oldTarget, newTarget := d.GetChange("traffic_target")

	oldTrafficTarget, ok := oldTarget.([]interface{})
	if !ok {
		logger.Warnf("wrong type conversion: expected []interface{}, got %T", oldTrafficTarget)
		return false
	}

	newTrafficTarget, ok := newTarget.([]interface{})
	if !ok {
		logger.Warnf("wrong type conversion: expected []interface{}, got %T", oldTrafficTarget)
		return false
	}

	if len(oldTrafficTarget) != len(newTrafficTarget) {
		return false
	}

	sort.Slice(oldTrafficTarget, func(i, j int) bool {
		return oldTrafficTarget[i].(map[string]interface{})["datacenter_id"].(int) < oldTrafficTarget[j].(map[string]interface{})["datacenter_id"].(int)
	})
	sort.Slice(newTrafficTarget, func(i, j int) bool {
		return newTrafficTarget[i].(map[string]interface{})["datacenter_id"].(int) < newTrafficTarget[j].(map[string]interface{})["datacenter_id"].(int)
	})

	length := len(oldTrafficTarget)
	for i := 0; i < length; i++ {
		for k, v := range oldTrafficTarget[i].(map[string]interface{}) {
			if k == "servers" {
				oldServers := oldTrafficTarget[i].(map[string]interface{})["servers"]
				newServers := newTrafficTarget[i].(map[string]interface{})["servers"]
				if !serversEqual(oldServers, newServers) {
					return false
				}
			} else {
				if newTrafficTarget[i].(map[string]interface{})[k] != v {
					return false
				}
			}
		}
	}

	return true
}

// serversEqual checks whether provided sets of ip addresses contain the same entries
func serversEqual(old, new interface{}) bool {
	logger := logger.Get("Akamai GTM", "serversEqual")

	oldServers, ok := old.(*schema.Set)
	if !ok {
		logger.Warnf("wrong type conversion: expected *schema.Set, got %T", oldServers)
		return false
	}

	newServers, ok := new.(*schema.Set)
	if !ok {
		logger.Warnf("wrong type conversion: expected *schema.Set, got %T", newServers)
		return false
	}

	if oldServers.Len() != newServers.Len() {
		return false
	}

	addresses := make(map[string]bool, oldServers.Len())
	for _, server := range oldServers.List() {
		addresses[server.(string)] = true
	}

	for _, server := range newServers.List() {
		_, ok := addresses[server.(string)]
		if !ok {
			return false
		}
	}

	return true
}
