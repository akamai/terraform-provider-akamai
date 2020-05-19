package akamai

import (
	"errors"
	"fmt"
	gtm "github.com/akamai/AkamaiOPEN-edgegrid-golang/configgtm-v1_4"
	"github.com/hashicorp/terraform-plugin-sdk/helper/schema"
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
				Type:     schema.TypeInt,
				Optional: true,
				Default:  300,
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
	log.Printf("[DEBUG] [Akamai GTMv1] Property Create status:")
	log.Printf("[DEBUG] [Akamai GTMv1] %v", cStatus.Status)

	if cStatus.Status.PropagationStatus == "DENIED" {
		return errors.New(cStatus.Status.Message)
	}
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
	log.Printf("[DEBUG] [Akamai GTMv1] Property Update  status:")
	log.Printf("[DEBUG] [Akamai GTMv1] %v", uStat)
	if uStat.PropagationStatus == "DENIED" {
		return errors.New(uStat.Message)
	}
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
	log.Printf("[DEBUG] [Akamai GTMv1] Property Delete status:")
	log.Printf("[DEBUG] [Akamai GTMv1] %v", uStat)
	if uStat.PropagationStatus == "DENIED" {
		return errors.New(uStat.Message)
	}
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
	} else if d.HasChange("stickiness_bonus_percentage") {
		prop.StickinessBonusPercentage = v.(int)
	}
	if v, ok := d.GetOk("stickiness_bonus_constant"); ok {
		prop.StickinessBonusConstant = v.(int)
	} else if d.HasChange("stickiness_bonus_constant") {
		prop.StickinessBonusConstant = v.(int)
	}
	if v, ok := d.GetOk("health_threshold"); ok {
		prop.HealthThreshold = v.(float64)
	} else if d.HasChange("health_threshold") {
		prop.HealthThreshold = v.(float64)
	}
	v := d.Get("ipv6")
	prop.Ipv6 = v.(bool)
	v = d.Get("use_computed_targets")
	prop.UseComputedTargets = v.(bool)
	if v, ok := d.GetOk("backup_ip"); ok {
		prop.BackupIp = v.(string)
	} else if d.HasChange("backup_ip") {
		prop.BackupIp = v.(string)
	}
	v = d.Get("balance_by_download_score")
	prop.BalanceByDownloadScore = v.(bool)
	if v, ok := d.GetOk("static_ttl"); ok {
		prop.StaticTTL = v.(int)
	}
	if v, ok := d.GetOk("unreachable_threshold"); ok {
		prop.UnreachableThreshold = v.(float64)
	} else if d.HasChange("unreachable_threshold") {
		prop.UnreachableThreshold = v.(float64)
	}
	if v, ok := d.GetOk("min_live_fraction"); ok {
		prop.MinLiveFraction = v.(float64)
	} else if d.HasChange("min_live_fraction") {
		prop.MinLiveFraction = v.(float64)
	}
	if v, ok := d.GetOk("health_multiplier"); ok {
		prop.HealthMultiplier = v.(float64)
	} else if d.HasChange("health_multiplier") {
		prop.HealthMultiplier = v.(float64)
	}
	if v, ok := d.GetOk("dynamic_ttl"); ok {
		prop.DynamicTTL = v.(int)
	}
	if v, ok := d.GetOk("max_unreachable_penalty"); ok {
		prop.MaxUnreachablePenalty = v.(int)
	} else if d.HasChange("max_unreachable_penalty") {
		prop.MaxUnreachablePenalty = v.(int)
	}
	if v, ok := d.GetOk("map_name"); ok {
		prop.MapName = v.(string)
	} else if d.HasChange("map_name") {
		prop.MapName = v.(string)
	}
	if v, ok := d.GetOk("handout_limit"); ok {
		prop.HandoutLimit = v.(int)
	} else if d.HasChange("handout_limit") {
		prop.HandoutLimit = v.(int)
	}
	if v, ok := d.GetOk("handout_mode"); ok {
		prop.HandoutMode = v.(string)
	}
	if v, ok := d.GetOk("load_imbalance_percentage"); ok {
		prop.LoadImbalancePercentage = v.(float64)
	} else if d.HasChange("load_imbalance_percentage") {
		prop.LoadImbalancePercentage = v.(float64)
	}
	if v, ok := d.GetOk("failover_delay"); ok {
		prop.FailoverDelay = v.(int)
	} else if d.HasChange("failover_delay") {
		prop.FailoverDelay = v.(int)
	}
	if v, ok := d.GetOk("backup_cname"); ok {
		prop.BackupCName = v.(string)
	} else if d.HasChange("backup_cname") {
		prop.BackupCName = v.(string)
	}
	if v, ok := d.GetOk("failback_delay"); ok {
		prop.FailbackDelay = v.(int)
	} else if d.HasChange("failback_delay") {
		prop.FailbackDelay = v.(int)
	}
	if v, ok := d.GetOk("health_max"); ok {
		prop.HealthMax = v.(float64)
	} else if d.HasChange("health_max") {
		prop.HealthMax = v.(float64)
	}
	v = d.Get("ghost_demand_reporting")
	prop.GhostDemandReporting = v.(bool)
	if v, ok := d.GetOk("weighted_hash_bits_for_ipv4"); ok {
		prop.WeightedHashBitsForIPv4 = v.(int)
	}
	if v, ok := d.GetOk("weighted_hash_bits_for_ipv6"); ok {
		prop.WeightedHashBitsForIPv6 = v.(int)
	}
	if v, ok := d.GetOk("cname"); ok {
		prop.CName = v.(string)
	} else if d.HasChange("cname") {
		prop.CName = v.(string)
	}
	if v, ok := d.GetOk("comments"); ok {
		prop.Comments = v.(string)
	} else if d.HasChange("comments") {
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

	objectInventory := make(map[int]*gtm.TrafficTarget, len(prop.TrafficTargets))
	if len(prop.TrafficTargets) > 0 {
		for _, aObj := range prop.TrafficTargets {
			objectInventory[aObj.DatacenterId] = aObj
		}
	}
	ttStateList := d.Get("traffic_target").([]interface{})
	for _, ttMap := range ttStateList {
		tt := ttMap.(map[string]interface{})
		objIndex := tt["datacenter_id"].(int)
		ttObject := objectInventory[objIndex]
		if ttObject == nil {
			log.Printf("[WARNING] [Akamai GTMv1] Property TrafficTarget %d NOT FOUND in returned GTM Object", tt["datacenter_id"])
			continue
		}
		tt["datacenter_id"] = ttObject.DatacenterId
		tt["name"] = ttObject.Name
		tt["enabled"] = ttObject.Enabled
		tt["weight"] = ttObject.Weight
		tt["handout_cname"] = ttObject.HandoutCName
		tt["servers"] = reconcileTerraformLists(tt["servers"].([]interface{}), convertStringToInterfaceList(ttObject.Servers))
		// remove object
		delete(objectInventory, objIndex)
	}
	if len(objectInventory) > 0 {
		// Objects not in the state yet. Add. Unfortunately, they not align with instance indices in the config
		for _, mttObj := range objectInventory {
			log.Printf("[DEBUG] [Akamai GTMv1] Property TrafficObject NEW State Object: %d", mttObj.DatacenterId)
			ttNew := map[string]interface{}{
				"datacenter_id": mttObj.DatacenterId,
				"enabled":       mttObj.Enabled,
				"weight":        mttObj.Weight,
				"name":          mttObj.Name,
				"handout_cname": mttObj.HandoutCName,
				"servers":       mttObj.Servers,
			}
			ttStateList = append(ttStateList, ttNew)
		}
	}
	d.Set("traffic_target", ttStateList)

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

	objectInventory := make(map[string]*gtm.StaticRRSet, len(prop.StaticRRSets))
	if len(prop.StaticRRSets) > 0 {
		for _, aObj := range prop.StaticRRSets {
			objectInventory[aObj.Type] = aObj
		}
	}
	rrStateList := d.Get("static_rr_set").([]interface{})
	for _, rrMap := range rrStateList {
		rr := rrMap.(map[string]interface{})
		objIndex := rr["type"].(string)
		rrObject := objectInventory[objIndex]
		if rrObject == nil {
			log.Printf("[WARNING] [Akamai GTMv1] Property StaticRRSet %s NOT FOUND in returned GTM Object", rr["type"])
			continue
		}
		rr["type"] = rrObject.Type
		rr["ttl"] = rrObject.TTL
		rr["rdata"] = reconcileTerraformLists(rr["rdata"].([]interface{}), convertStringToInterfaceList(rrObject.Rdata))
		// remove object
		delete(objectInventory, objIndex)
	}
	if len(objectInventory) > 0 {
		log.Printf("[DEBUG] [Akamai GTMv1] Property StaticRRSet objects left...")
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
	d.Set("static_rr_set", rrStateList)

}

// Populate existing Livenesstest  object from resource data
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
			lt.ErrorPenalty = v["error_penalty"].(float64)
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
			lt.TestObjectUsername = v["test_object_username"].(string)
			lt.TimeoutPenalty = v["timeout_penalty"].(float64)
			lt.AnswersRequired = v["answers_required"].(bool)
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

	objectInventory := make(map[string]*gtm.LivenessTest, len(prop.LivenessTests))
	if len(prop.LivenessTests) > 0 {
		for _, aObj := range prop.LivenessTests {
			objectInventory[aObj.Name] = aObj
		}
	}
	ltStateList := d.Get("liveness_test").([]interface{})
	for _, ltMap := range ltStateList {
		lt := ltMap.(map[string]interface{})
		objIndex := lt["name"].(string)
		ltObject := objectInventory[objIndex]
		if ltObject == nil {
			log.Printf("[WARNING] [Akamai GTMv1] Property LivenessTest  %s NOT FOUND in returned GTM Object", lt["name"])
			continue
		}
		lt["name"] = ltObject.Name
		lt["error_penalty"] = ltObject.ErrorPenalty
		lt["peer_certificate_verification"] = ltObject.PeerCertificateVerification
		lt["test_interval"] = ltObject.TestInterval
		lt["test_object"] = ltObject.TestObject
		lt["request_string"] = ltObject.RequestString
		lt["response_string"] = ltObject.ResponseString
		lt["http_error3xx"] = ltObject.HttpError3xx
		lt["http_error4xx"] = ltObject.HttpError4xx
		lt["http_error5xx"] = ltObject.HttpError5xx
		lt["disabled"] = ltObject.Disabled
		lt["test_object_protocol"] = ltObject.TestObjectProtocol
		lt["test_object_password"] = ltObject.TestObjectPassword
		lt["test_object_port"] = ltObject.TestObjectPort
		lt["ssl_client_private_key"] = ltObject.SslClientPrivateKey
		lt["ssl_client_certificate"] = ltObject.SslClientCertificate
		lt["disable_nonstandard_port_warning"] = ltObject.DisableNonstandardPortWarning
		lt["test_object_username"] = ltObject.TestObjectUsername
		lt["test_timeout"] = ltObject.TestTimeout
		lt["timeout_penalty"] = ltObject.TimeoutPenalty
		lt["answers_required"] = ltObject.AnswersRequired
		lt["resource_type"] = ltObject.ResourceType
		lt["recursion_requested"] = ltObject.RecursionRequested
		httpHeaderListNew := make([]interface{}, len(ltObject.HttpHeaders))
		for i, r := range ltObject.HttpHeaders {
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
		log.Printf("[DEBUG] [Akamai GTMv1] Property LivenessTest objects left...")
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
				"test_object_username":             l.TestObjectUsername,
				"test_timeout":                     l.TestTimeout,
				"timeout_penalty":                  l.TimeoutPenalty,
				"answers_required":                 l.AnswersRequired,
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
			ltNew["http_header"] = httpHeaderListNew
			ltStateList = append(ltStateList, ltNew)
		}
	}
	d.Set("liveness_test", ltStateList)

}

func convertStringToInterfaceList(stringList []string) []interface{} {

	log.Printf("[DEBUG] [Akamai GTMv1] String List: %v", stringList)
	retList := make([]interface{}, 0, len(stringList))
	for _, v := range stringList {
		retList = append(retList, v)
	}

	return retList

}

func convertIntToInterfaceList(intList []int) []interface{} {

	log.Printf("[DEBUG] [Akamai GTMv1] Int List: %v", intList)
	retList := make([]interface{}, 0, len(intList))
	for _, v := range intList {
		retList = append(retList, v)
	}

	return retList

}

func convertInt64ToInterfaceList(intList []int64) []interface{} {

	log.Printf("[DEBUG] [Akamai GTMv1] Int List: %v", intList)
	retList := make([]interface{}, 0, len(intList))
	for _, v := range intList {
		retList = append(retList, v)
	}

	return retList

}

// Util method to reconcile list configs. Type agnostic. Goal: maintain order of tf list config
func reconcileTerraformLists(terraList []interface{}, newList []interface{}) []interface{} {

	log.Printf("[DEBUG] [Akamai GTMv1] Existing Terra List: %v", terraList)
	log.Printf("[DEBUG] [Akamai GTMv1] Read List: %v", newList)
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

	log.Printf("[DEBUG] [Akamai GTMv1] Updated Terra List: %v", updatedList)
	return updatedList

}
