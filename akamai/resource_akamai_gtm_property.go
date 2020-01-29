package akamai

import (
	"encoding/json"
	"errors"
	"fmt"
	gtm "github.com/akamai/AkamaiOPEN-edgegrid-golang/configgtm-v1_4"
	"github.com/hashicorp/terraform/helper/schema"
	"log"
	"strings"
)

func resourceGTMv1Property() *schema.Resource {
	return &schema.Resource{
		Create: resourceGTMv1PropertyCreate,
		Read:   resourceGTMv1PropertyRead,
		Update: resourceGTMv1PropertyUpdate,
		Delete: resourceGTMv1PropertyDelete,
		Exists: resourceGTMv1PropertyExists,
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
				Type:     schema.TypeInt,
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
				Type:     schema.TypeInt,
				Optional: true,
			},
			"static_rr_set": &schema.Schema{
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
				Type:     schema.TypeInt,
				Optional: true,
			},
			"min_live_fraction": {
				Type:     schema.TypeFloat,
				Optional: true,
			},
			"health_multiplier": {
				Type:     schema.TypeInt,
				Optional: true,
			},
			"dynamic_ttl": {
				Type:     schema.TypeInt,
				Optional: true,
				Default:  300,
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
				Type:     schema.TypeInt,
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
				Type:     schema.TypeList,
				Required: true,
				MinItems: 1,
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
							Type:     schema.TypeList,
							Elem:     &schema.Schema{Type: schema.TypeString},
							Optional: true,
						},
						"name": {
							Type:     schema.TypeString,
							Optional: true,
							//Default:  "",
						},
						"handout_cname": {
							Type:     schema.TypeString,
							Optional: true,
							//Default:  "",
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
							Type:     schema.TypeInt,
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
							Required: true,
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
						"host_header": {
							Type:     schema.TypeString,
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
							Type:     schema.TypeInt,
							Optional: true,
						},
						"answer_required": {
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
					},
				},
			},
		},
	}
}

// utility func to parse Terraform resource string id
func parseResourceStringId(id string) (string, string, error) {

	parts := strings.SplitN(id, ":", 2)
	if len(parts) != 2 || parts[0] == "" || parts[1] == "" {
		return "", "", errors.New("Invalid resource id")
	}

	return parts[0], parts[1], nil

}

// utility func to parse Terraform property resource id
func parsePropertyResourceId(id string) (string, string, error) {

	return parseResourceStringId(id)
}

// Create a new GTM Property
func resourceGTMv1PropertyCreate(d *schema.ResourceData, meta interface{}) error {

	domain := d.Get("domain").(string)

	log.Printf("[INFO] [Akamai GTM] Creating property [%s] in domain [%s]", d.Get("name").(string), domain)
	newProp := populateNewPropertyObject(d)
	log.Printf("[DEBUG] [Akamai GTMv1] Proposed New Property: [%v]", newProp)
	cStatus, err := newProp.Create(domain)
	if err != nil {
		log.Printf("[ERROR] PropertyCreate failed: %s", err.Error())
		return err
	}
	b, err := json.Marshal(cStatus.Status)
	if err != nil {
		log.Printf("[ERROR] PropertyCreate failed: %s", err.Error())
		return err
	}
	log.Printf("[DEBUG] [Akamai GTMv1] Property Create status:")
	log.Printf("[DEBUG] [Akamai GTMv1] %v", b)

	if d.Get("wait_on_complete").(bool) {
		done, err := waitForCompletion(domain)
		if done {
			log.Printf("[INFO] [Akamai GTMv1] Property Create completed")
		} else {
			if err == nil {
				log.Printf("[INFO] [Akamai GTMv1] Property Create pending")
			} else {
				log.Printf("[WARNING] [Akamai GTMv1] Property Create failed [%s]", err.Error())
				return err
			}
		}

	}

	// Give terraform the ID. Format domain::property
	propertyId := fmt.Sprintf("%s:%s", domain, cStatus.Resource.Name)
	log.Printf("[DEBUG] [Akamai GTMv1] Generated Property Resource Id: %s", propertyId)
	d.SetId(propertyId)
	return resourceGTMv1PropertyRead(d, meta)

}

// Only ever save data from the tf config in the tf state file, to help with
// api issues. See func unmarshalResourceData for more info.
func resourceGTMv1PropertyRead(d *schema.ResourceData, meta interface{}) error {

	log.Printf("[DEBUG] [Akamai GTMv1] READ")
	log.Printf("[DEBUG] Reading [Akamai GTMv1] Property: %s", d.Id())
	// retrieve the property and domain
	domain, property, err := parsePropertyResourceId(d.Id())
	if err != nil {
		return errors.New("Invalid property resource Id")
	}
	prop, err := gtm.GetProperty(property, domain)
	if err != nil {
		log.Printf("[ERROR] PropertyRead failed: %s", err.Error())
		return err
	}
	populateTerraformPropertyState(d, prop)
	log.Printf("[DEBUG] [Akamai GTMv1] READ %v", prop)
	return nil
}

// Update GTM Property
func resourceGTMv1PropertyUpdate(d *schema.ResourceData, meta interface{}) error {

	log.Printf("[DEBUG] [Akamai GTMv1] UPDATE")
	log.Printf("[DEBUG] Updating [Akamai GTMv1] Property: %s", d.Id())
	// pull domain and property out of resource id
	domain, property, err := parsePropertyResourceId(d.Id())
	if err != nil {
		return errors.New("Invalid property resource Id")
	}
	// Get existing property
	existProp, err := gtm.GetProperty(property, domain)
	if err != nil {
		log.Printf("[ERROR] PropertyUpdate failed: %s", err.Error())
		return err
	}
	log.Printf("[DEBUG] Updating [Akamai GTMv1] Property BEFORE: %v", existProp)
	populatePropertyObject(d, existProp)
	log.Printf("[DEBUG] Updating [Akamai GTMv1] Property PROPOSED: %v", existProp)
	uStat, err := existProp.Update(domain)
	if err != nil {
		log.Printf("[ERROR] PropertyUpdate failed: %s", err.Error())
		return err
	}
	b, err := json.Marshal(uStat)
	if err != nil {
		log.Printf("[ERROR] PropertyUpdate failed: %s", err.Error())
		return err
	}
	log.Printf("[DEBUG] [Akamai GTMv1] Property Update  status:")
	log.Printf("[DEBUG] [Akamai GTMv1] %v", b)

	if d.Get("wait_on_complete").(bool) {
		done, err := waitForCompletion(domain)
		if done {
			log.Printf("[INFO] [Akamai GTMv1] Property update completed")
		} else {
			if err == nil {
				log.Printf("[INFO] [Akamai GTMv1] Property update pending")
			} else {
				log.Printf("[WARNING] [Akamai GTMv1] Property update failed [%s]", err.Error())
				return err
			}
		}

	}

	return resourceGTMv1PropertyRead(d, meta)
}

// Import GTM Property.
func resourceGTMv1PropertyImport(d *schema.ResourceData, meta interface{}) ([]*schema.ResourceData, error) {

	log.Printf("[INFO] [Akamai GTM] Property [%s] Import", d.Id())
	// pull domain and property out of resource id
	domain, property, err := parsePropertyResourceId(d.Id())
	if err != nil {
		return []*schema.ResourceData{d}, errors.New("Invalid property resource Id")
	}
	prop, err := gtm.GetProperty(property, domain)
	if err != nil {
		return nil, err
	}
	d.Set("domain", domain)
	d.Set("wait_on_complete", true)
	populateTerraformPropertyState(d, prop)

	// use same Id as passed in
	log.Printf("[INFO] [Akamai GTM] Property [%s] [%s] Imported", d.Id(), d.Get("name"))
	return []*schema.ResourceData{d}, nil
}

// Delete GTM Property.
func resourceGTMv1PropertyDelete(d *schema.ResourceData, meta interface{}) error {

	log.Printf("[DEBUG] [Akamai GTMv1] DELETE")
	log.Printf("[DEBUG] Deleting [Akamai GTMv1] Property: %s", d.Id())
	// Get existing property
	domain, property, err := parsePropertyResourceId(d.Id())
	if err != nil {
		return errors.New("Invalid property resource Id")
	}
	existProp, err := gtm.GetProperty(property, domain)
	if err != nil {
		log.Printf("[ERROR] PropertyDelete failed: %s", err.Error())
		return err
	}
	log.Printf("[DEBUG] Deleting [Akamai GTMv1] Property: %v", existProp)
	uStat, err := existProp.Delete(domain)
	if err != nil {
		log.Printf("[ERROR] PropertyDelete failed: %s", err.Error())
		return err
	}
	b, err := json.Marshal(uStat)
	if err != nil {
		log.Printf("[ERROR] PropertyDelete failed: %s", err.Error())
		return err
	}
	log.Printf("[DEBUG] [Akamai GTMv1] Property Delete status:")
	log.Printf("[DEBUG] [Akamai GTMv1] %v", b)

	if d.Get("wait_on_complete").(bool) {
		done, err := waitForCompletion(domain)
		if done {
			log.Printf("[INFO] [Akamai GTMv1] Property delete completed")
		} else {
			if err == nil {
				log.Printf("[INFO] [Akamai GTMv1] Property delete pending")
			} else {
				log.Printf("[WARNING] [Akamai GTMv1] Property delete failed [%s]", err.Error())
				return err
			}
		}

	}

	// if succcessful ....
	d.SetId("")
	return nil
}

// Test GTM Property existance
func resourceGTMv1PropertyExists(d *schema.ResourceData, meta interface{}) (bool, error) {

	log.Printf("[DEBUG] [Akamai GTMv1] Exists")
	// pull domain and property out of resource id
	domain, property, err := parsePropertyResourceId(d.Id())
	if err != nil {
		return false, errors.New("Invalid property resource Id")
	}
	log.Printf("[DEBUG] [Akamai GTMv1] Searching for existing property [%s] in domain %s", property, domain)
	prop, err := gtm.GetProperty(property, domain)
	return prop != nil, err
}

// Create and populate a new property object from resource data
func populateNewPropertyObject(d *schema.ResourceData) *gtm.Property {

	propObj := gtm.NewProperty(d.Get("name").(string))
	propObj.TrafficTargets = make([]*gtm.TrafficTarget, 0)
	propObj.LivenessTests = make([]*gtm.LivenessTest, 0)
	propObj.MxRecords = make([]*gtm.MxRecord, 0)
	populatePropertyObject(d, propObj)

	return propObj

}

// Populate existing property object from resource data
func populatePropertyObject(d *schema.ResourceData, prop *gtm.Property) {

	if v, ok := d.GetOk("name"); ok {
		prop.Name = v.(string)
	}
	if v, ok := d.GetOk("type"); ok {
		prop.Type = v.(string)
	}
	if v, ok := d.GetOk("score_aggregation_type"); ok {
		prop.ScoreAggregationType = v.(string)
	}
	if v, ok := d.GetOk("stickiness_bonus_percentage"); ok {
		prop.StickinessBonusPercentage = v.(int)
	}
	if v, ok := d.GetOk("stickiness_bonus_constant"); ok {
		prop.StickinessBonusConstant = v.(int)
	}
	if v, ok := d.GetOk("health_threshold"); ok {
		prop.HealthThreshold = v.(int)
	}
	if v, ok := d.GetOk("ipv6"); ok {
		prop.Ipv6 = v.(bool)
	}
	if v, ok := d.GetOk("use_computed_targets"); ok {
		prop.UseComputedTargets = v.(bool)
	}
	if v, ok := d.GetOk("backup_ip"); ok {
		prop.BackupIp = v.(string)
	}
	if v, ok := d.GetOk("balance_by_download_score"); ok {
		prop.BalanceByDownloadScore = v.(bool)
	}
	if v, ok := d.GetOk("static_ttl"); ok {
		prop.StaticTTL = v.(int)
	}
	if v, ok := d.GetOk("unreachable_threshold"); ok {
		prop.UnreachableThreshold = v.(int)
	}
	if v, ok := d.GetOk("min_live_fraction"); ok {
		prop.MinLiveFraction = v.(float64)
	}
	if v, ok := d.GetOk("health_multiplier"); ok {
		prop.HealthMultiplier = v.(int)
	}
	if v, ok := d.GetOk("dynamic_ttl"); ok {
		prop.DynamicTTL = v.(int)
	}
	if v, ok := d.GetOk("max_unreachable_penalty"); ok {
		prop.MaxUnreachablePenalty = v.(int)
	}
	if v, ok := d.GetOk("map_name"); ok {
		prop.MapName = v.(string)
	}
	if v, ok := d.GetOk("handout_limit"); ok {
		prop.HandoutLimit = v.(int)
	}
	if v, ok := d.GetOk("handout_mode"); ok {
		prop.HandoutMode = v.(string)
	}
	if v, ok := d.GetOk("load_imbalance_percentage"); ok {
		prop.LoadImbalancePercentage = v.(float64)
	}
	if v, ok := d.GetOk("failover_delay"); ok {
		prop.FailoverDelay = v.(int)
	}
	if v, ok := d.GetOk("backup_cname"); ok {
		prop.BackupCName = v.(string)
	}
	if v, ok := d.GetOk("failback_delay"); ok {
		prop.FailbackDelay = v.(int)
	}
	if v, ok := d.GetOk("health_max"); ok {
		prop.HealthMax = v.(int)
	}
	if v, ok := d.GetOk("ghost_demand_reporting"); ok {
		prop.GhostDemandReporting = v.(bool)
	}
	if v, ok := d.GetOk("weighted_hash_bits_for_ipv4"); ok {
		prop.WeightedHashBitsForIPv4 = v.(int)
	}
	if v, ok := d.GetOk("weighted_hash_bits_for_ipv6"); ok {
		prop.WeightedHashBitsForIPv6 = v.(int)
	}
	if v, ok := d.GetOk("cname"); ok {
		prop.CName = v.(string)
	}
	if v, ok := d.GetOk("comments"); ok {
		prop.Comments = v.(string)
	}
	populateTrafficTargetObject(d, prop)
	populateStaticRRSetObject(d, prop)
	populateLivenessTestObject(d, prop)

}

// Populate Terraform state from provided Property object
func populateTerraformPropertyState(d *schema.ResourceData, prop *gtm.Property) {

	// walk thru all state elements
	d.Set("name", prop.Name)
	d.Set("type", prop.Type)
	d.Set("ipv6", prop.Ipv6)
	d.Set("score_aggregation_type", prop.ScoreAggregationType)
	d.Set("stickiness_bonus_percentage", prop.StickinessBonusPercentage)
	d.Set("stickiness_bonus_constant", prop.StickinessBonusConstant)
	d.Set("health_threshold", prop.HealthThreshold)
	d.Set("use_computed_targets", prop.UseComputedTargets)
	d.Set("backup_ip", prop.BackupIp)
	d.Set("balance_by_download_score", prop.BalanceByDownloadScore)
	d.Set("static_ttl", prop.StaticTTL)
	d.Set("unreachable_threshold", prop.UnreachableThreshold)
	d.Set("min_live_fraction", prop.MinLiveFraction)
	d.Set("health_multiplier", prop.HealthMultiplier)
	d.Set("dynamic_ttl", prop.DynamicTTL)
	d.Set("max_unreachable_penalty", prop.MaxUnreachablePenalty)
	d.Set("map_name", prop.MapName)
	d.Set("handout_limit", prop.HandoutLimit)
	d.Set("handout_mode", prop.HandoutMode)
	d.Set("load_imbalance_percentage", prop.LoadImbalancePercentage)
	d.Set("failover_delay", prop.FailoverDelay)
	d.Set("backup_cname", prop.BackupCName)
	d.Set("failback_delay", prop.FailbackDelay)
	d.Set("health_max", prop.HealthMax)
	d.Set("ghost_demand_reporting", prop.GhostDemandReporting)
	d.Set("weighted_hash_bits_for_ipv4", prop.WeightedHashBitsForIPv4)
	d.Set("weighted_hash_bits_for_ipv6", prop.WeightedHashBitsForIPv6)
	d.Set("cname", prop.CName)
	d.Set("comments", prop.Comments)
	populateTerraformTrafficTargetState(d, prop)
	populateTerraformStaticRRSetState(d, prop)
	populateTerraformLivenessTestState(d, prop)

}

// create and populate GTM Property TrafficTargets object
func populateTrafficTargetObject(d *schema.ResourceData, prop *gtm.Property) {

	// pull apart List
	tt := d.Get("traffic_target")
	if tt != nil {
		traffTargList := tt.([]interface{})
		trafficObjList := make([]*gtm.TrafficTarget, len(traffTargList)) // create new object list
		for i, v := range traffTargList {
			ttMap := v.(map[string]interface{})
			trafficTarget := prop.NewTrafficTarget() // create new object
			trafficTarget.DatacenterId = ttMap["datacenter_id"].(int)
			trafficTarget.Enabled = ttMap["enabled"].(bool)
			trafficTarget.Weight = ttMap["weight"].(float64)
			if ttMap["servers"] != nil {
				ls := make([]string, len(ttMap["servers"].([]interface{})))
				for i, sl := range ttMap["servers"].([]interface{}) {
					ls[i] = sl.(string)
				}
				trafficTarget.Servers = ls
			}
			trafficTarget.Name = ttMap["name"].(string)
			trafficTarget.HandoutCName = ttMap["handout_cname"].(string)
			trafficObjList[i] = trafficTarget
		}
		prop.TrafficTargets = trafficObjList
	}
}

// create and populate Terraform traffic_targets schema
func populateTerraformTrafficTargetState(d *schema.ResourceData, prop *gtm.Property) {

	traffListNew := make([]interface{}, len(prop.TrafficTargets))
	for i, tt := range prop.TrafficTargets {
		traffSvrNew := map[string]interface{}{
			"datacenter_id": tt.DatacenterId,
			"enabled":       tt.Enabled,
			"weight":        tt.Weight,
			"name":          tt.Name,
			"handout_cname": tt.HandoutCName,
			"servers":       tt.Servers,
		}
		traffListNew[i] = traffSvrNew
	}
	d.Set("traffic_target", traffListNew)

}

// Populate existing static_rr_sets object from resource data
func populateStaticRRSetObject(d *schema.ResourceData, prop *gtm.Property) {

	// pull apart List
	staticSetList := d.Get("static_rr_set").([]interface{})
	if staticSetList != nil {
		staticObjList := make([]*gtm.StaticRRSet, len(staticSetList)) // create new object list
		for i, v := range staticSetList {
			recMap := v.(map[string]interface{})
			record := prop.NewStaticRRSet() // create new object
			record.TTL = recMap["ttl"].(int)
			record.Type = recMap["type"].(string)
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
func populateTerraformStaticRRSetState(d *schema.ResourceData, prop *gtm.Property) {

	recordListNew := make([]interface{}, len(prop.StaticRRSets))
	for i, r := range prop.StaticRRSets {
		staticRecordNew := map[string]interface{}{
			"type":  r.Type,
			"ttl":   r.TTL,
			"rdata": r.Rdata,
		}
		recordListNew[i] = staticRecordNew
	}
	d.Set("static_rr_set", recordListNew)

}

// Populate existing livenesstest  object from resource data
func populateLivenessTestObject(d *schema.ResourceData, prop *gtm.Property) {

	liveTestList := d.Get("liveness_test").([]interface{})
	if liveTestList != nil {
		liveTestObjList := make([]*gtm.LivenessTest, len(liveTestList)) // create new object list
		for i, l := range liveTestList {
			v := l.(map[string]interface{})
			lt := prop.NewLivenessTest(v["name"].(string),
				v["test_object_protocol"].(string),
				v["test_interval"].(int),
				float32(v["test_timeout"].(float64))) // create new object
			lt.ErrorPenalty = v["error_penalty"].(int)
			lt.PeerCertificateVerification = v["peer_certificate_verification"].(bool)
			lt.TestObject = v["test_object"].(string)
			lt.RequestString = v["request_string"].(string)
			lt.ResponseString = v["response_string"].(string)
			lt.HttpError3xx = v["http_error3xx"].(bool)
			lt.HttpError4xx = v["http_error4xx"].(bool)
			lt.HttpError5xx = v["http_error5xx"].(bool)
			lt.Disabled = v["disabled"].(bool)
			lt.TestObjectPassword = v["test_object_password"].(string)
			lt.TestObjectPort = v["test_object_port"].(int)
			lt.SslClientPrivateKey = v["ssl_client_private_key"].(string)
			lt.SslClientCertificate = v["ssl_client_certificate"].(string)
			lt.DisableNonstandardPortWarning = v["disable_nonstandard_port_warning"].(bool)
			lt.HostHeader = v["host_header"].(string)
			lt.TestObjectUsername = v["test_object_username"].(string)
			lt.TimeoutPenalty = v["timeout_penalty"].(int)
			lt.AnswerRequired = v["answer_required"].(bool)
			lt.ResourceType = v["resource_type"].(string)
			lt.RecursionRequested = v["recursion_requested"].(bool)
			httpHeaderList := v["http_header"].([]interface{})
			if httpHeaderList != nil {
				headerObjList := make([]*gtm.HttpHeader, len(httpHeaderList)) // create new object list
				for i, h := range httpHeaderList {
					recMap := h.(map[string]interface{})
					record := lt.NewHttpHeader() // create new object
					record.Name = recMap["name"].(string)
					record.Value = recMap["value"].(string)
					headerObjList[i] = record
				}
				lt.HttpHeaders = headerObjList
			}
			liveTestObjList[i] = lt
		}
		prop.LivenessTests = liveTestObjList
	}
}

// create and populate Terraform liveness_test schema
func populateTerraformLivenessTestState(d *schema.ResourceData, prop *gtm.Property) {

	livenessListNew := make([]interface{}, len(prop.LivenessTests))
	for i, l := range prop.LivenessTests {
		livenessNew := map[string]interface{}{
			"name":                             l.Name,
			"error_penalty":                    l.ErrorPenalty,
			"peer_certificate_verification":    l.PeerCertificateVerification,
			"test_interval":                    l.TestInterval,
			"test_object":                      l.TestObject,
			"request_string":                   l.RequestString,
			"response_string":                  l.ResponseString,
			"http_error3xx":                    l.HttpError3xx,
			"http_error4xx":                    l.HttpError4xx,
			"http_error5xx":                    l.HttpError5xx,
			"disabled":                         l.Disabled,
			"test_object_protocol":             l.TestObjectProtocol,
			"test_object_password":             l.TestObjectPassword,
			"test_object_port":                 l.TestObjectPort,
			"ssl_client_private_key":           l.SslClientPrivateKey,
			"ssl_client_certificate":           l.SslClientCertificate,
			"disable_nonstandard_port_warning": l.DisableNonstandardPortWarning,
			"host_header":                      l.HostHeader,
			"test_object_username":             l.TestObjectUsername,
			"test_timeout":                     l.TestTimeout,
			"timeout_penalty":                  l.TimeoutPenalty,
			"answer_required":                  l.AnswerRequired,
			"resource_type":                    l.ResourceType,
			"recursion_requested":              l.RecursionRequested,
		}
		httpHeaderListNew := make([]interface{}, len(l.HttpHeaders))
		for i, r := range l.HttpHeaders {
			httpHeaderNew := map[string]interface{}{
				"name":  r.Name,
				"value": r.Value,
			}
			httpHeaderListNew[i] = httpHeaderNew
		}
		livenessNew["http_header"] = httpHeaderListNew
		livenessListNew[i] = livenessNew
	}
	d.Set("liveness_test", livenessListNew)

}
