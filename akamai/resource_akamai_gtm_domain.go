package akamai

import (
	"encoding/json"
	"errors"
	"fmt"
	"log"
	"strings"
	"time"

	gtmv1_3 "github.com/akamai/AkamaiOPEN-edgegrid-golang/configgtm-v1_3"
	"github.com/hashicorp/terraform/helper/schema"
)

func resourceGTMv1_3Domain() *schema.Resource {
	return &schema.Resource{
		Create: resourceGTMv1_3DomainCreate,
		Read:   resourceGTMv1_3DomainRead,
		Update: resourceGTMv1_3DomainUpdate,
		Delete: resourceGTMv1_3DomainDelete,
		Exists: resourceGTMv1_3DomainExists,
		Importer: &schema.ResourceImporter{
			State: schema.ImportStatePassthrough,
		},
		Schema: map[string]*schema.Schema{
			"contract": {
				Type:     schema.TypeString,
				Optional: true,
				Default:  "",
			},
			"group": {
				Type:     schema.TypeString,
				Optional: true,
				Default:  "",
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
				Type:         schema.TypeString,
				Required:     true,
				ValidateFunc: validateDomainType,
			},
			"comment": {
				Type:     schema.TypeString,
				Optional: true,
				Default:  "Managed by Terraform",
			},
			"default_unreachable_threshold": {
				Type:     schema.TypeFloat,
				Computed: true,
			},
			"email_notification_list": {
				Type:     schema.TypeList,
				Elem:     &schema.Schema{Type: schema.TypeString},
				Optional: true,
			},
			"min_pingable_region_fraction": {
				Type:     schema.TypeFloat,
				Computed: true,
			},
			"default_timeout_penalty": {
				Type:     schema.TypeInt,
				Optional: true,
				Default:  25,
			},
			"servermonitor_liveness_count": {
				Type:     schema.TypeInt,
				Computed: true,
			},
			"round_robin_prefix": {
				Type:     schema.TypeString,
				Computed: true,
			},
			"servermonitor_load_count": {
				Type:     schema.TypeInt,
				Computed: true,
			},
			"ping_interval": {
				Type:     schema.TypeInt,
				Computed: true,
			},
			"max_ttl": {
				Type:     schema.TypeInt,
				Computed: true,
			},
			"load_imbalance_percentage": {
				Type:     schema.TypeFloat,
				Optional: true,
			},
			"default_health_max": {
				Type:     schema.TypeInt,
				Computed: true,
			},
			"map_update_interval": {
				Type:     schema.TypeInt,
				Computed: true,
			},
			"max_properties": {
				Type:     schema.TypeInt,
				Computed: true,
			},
			"max_resources": {
				Type:     schema.TypeInt,
				Computed: true,
			},
			"default_ssl_client_private_key": {
				Type:     schema.TypeString,
				Optional: true,
			},
			"default_error_penalty": {
				Type:     schema.TypeInt,
				Optional: true,
				Default:  75,
			},
			"max_test_timeout": {
				Type:     schema.TypeFloat,
				Computed: true,
			},
			"cname_coalescing_enabled": {
				Type:     schema.TypeBool,
				Optional: true,
			},
			"default_health_multiplier": {
				Type:     schema.TypeInt,
				Computed: true,
			},
			"servermonitor_pool": {
				Type:     schema.TypeString,
				Computed: true,
			},
			"load_feedback": {
				Type:     schema.TypeBool,
				Optional: true,
			},
			"min_ttl": {
				Type:     schema.TypeInt,
				Computed: true,
			},
			"default_max_unreachable_penalty": {
				Type:     schema.TypeInt,
				Computed: true,
			},
			"default_health_threshold": {
				Type:     schema.TypeInt,
				Computed: true,
			},
			"min_test_interval": {
				Type:     schema.TypeInt,
				Computed: true,
			},
			"ping_packet_size": {
				Type:     schema.TypeInt,
				Computed: true,
			},
			"default_ssl_client_certificate": {
				Type:     schema.TypeString,
				Optional: true,
			},
		},
	}
}

// Retrieve optional query args. contractId, groupId [and accountSwitchKey] supported.
func GetQueryArgs(d *schema.ResourceData) map[string]string {

	qArgs := make(map[string]string)
	contract := strings.TrimPrefix(d.Get("contract").(string), "ctr_")
	if contract != "" && len(contract) > 0 {
		qArgs["contractId"] = contract
	}
	groupId := strings.TrimPrefix(d.Get("group").(string), "grp_")
	if groupId != "" && len(groupId) > 0 {
		qArgs["gid"] = groupId
	}
	//accountSwitch := d.Get("account_switch_key").(string)
	// if accountSwitch != nil && len(accountSwitch) > 0 {
	//        qArgs["accountSwitchKey"] = accountSwitch
	//}

	return qArgs

}

// Create a new GTM Domain
func resourceGTMv1_3DomainCreate(d *schema.ResourceData, meta interface{}) error {

	dname := d.Get("name").(string)
	//comment := d.Get("comment").(string)

	log.Printf("[INFO] [Akamai GTM] Creating domain [%s]", dname)
	newDom := populateNewDomainObject(d)
	log.Printf("[DEBUG] [Akamai GTMV1_3] Domain: [%v]", newDom)

	cStatus, err := newDom.Create(GetQueryArgs(d))
	if err != nil {
		fmt.Println(err)
		return err
	}
	b, err := json.Marshal(cStatus)
	if err != nil {
		fmt.Println(err)
		return err
	}
	fmt.Println(string(b))
	log.Printf("[DEBUG] [Akamai GTMV1_3] Create status:")
	log.Printf("[DEBUG] [Akamai GTMV1_3] %v", b)

	if d.Get("wait_on_complete").(bool) {
		done, err := waitForCompletion(dname)
		if done {
			log.Printf("[INFO] [Akamai GTMV1_3] Domain Create completed")
		} else {
			if err == nil {
				log.Printf("[INFO] [Akamai GTMV1_3] Domain Create pending")
			} else {
				log.Printf("[WARNING] [Akamai GTMV1_3] Domain Create failed [%s]", err.Error())
				return err
			}
		}

	}

	// Give terraform the ID
	d.SetId(dname)
	return resourceGTMv1_3DomainRead(d, meta)

}

// Only ever save data from the tf config in the tf state file, to help with
// api issues. See func unmarshalResourceData for more info.
func resourceGTMv1_3DomainRead(d *schema.ResourceData, meta interface{}) error {
	log.Printf("[DEBUG] [Akamai GTMv1_3] READ")
	log.Printf("[DEBUG] Reading [Akamai GTMv1_3] Domain: %s", d.Id())
	// retrieve the domain
	dom, err := gtmv1_3.GetDomain(d.Id())
	if err != nil {
		fmt.Println(err)
		log.Printf("[DEBUG] [Akamai GTMV1_3] Domain Read error: %s", err.Error())
		return err
	}
	populateTerraformState(d, dom)
	log.Printf("[DEBUG] [Akamai GTMV1_3] READ %v", dom)
	return nil
}

// Update GTM Domain
func resourceGTMv1_3DomainUpdate(d *schema.ResourceData, meta interface{}) error {

	log.Printf("[DEBUG] [Akamai GTMv1_3] UPDATE")
	log.Printf("[DEBUG] Updating [Akamai GTMv1_3] Domain: %s", d.Id())
	// Get existing domain
	existDom, err := gtmv1_3.GetDomain(d.Id())
	if err != nil {
		fmt.Println(err.Error())
		return err
	}
	log.Printf("[DEBUG] Updating [Akamai GTMv1_3] Domain BEFORE: %v", existDom)
	populateDomainObject(d, existDom)
	log.Printf("[DEBUG] Updating [Akamai GTMv1_3] Domain PROPOSED: %v", existDom)
	//existDom := populateNewDomainObject(d)
	uStat, err := existDom.Update(GetQueryArgs(d))
	if err != nil {
		fmt.Println(err.Error())
		return err
	}
	b, err := json.Marshal(uStat)
	if err != nil {
		fmt.Println(err.Error())
		return err
	}
	fmt.Println(string(b))
	log.Printf("[DEBUG] [Akamai GTMV1_3] Update status:")
	log.Printf("[DEBUG] [Akamai GTMV1_3] %v", b)

	if d.Get("wait_on_complete").(bool) {
		done, err := waitForCompletion(d.Id())
		if done {
			log.Printf("[INFO] [Akamai GTMV1_3] Domain update completed")
		} else {
			if err == nil {
				log.Printf("[INFO] [Akamai GTMV1_3] Domain update pending")
			} else {
				log.Printf("[WARNING] [Akamai GTMV1_3] Domain update failed [%s]", err.Error())
				return err
			}
		}

	}

	// Give terraform the ID
	return resourceGTMv1_3DomainRead(d, meta)

}

// Delete GTM Domain. Not Supported in current API version.
func resourceGTMv1_3DomainDelete(d *schema.ResourceData, meta interface{}) error {
	log.Printf("[DEBUG] Deleting GTM Domain")
        log.Printf("[DEBUG] [Akamai GTMv1_3] Domain: %s", d.Id())
        // Get existing domain
        existDom, err := gtmv1_3.GetDomain(d.Id())
        if err != nil {
                fmt.Println(err.Error())
                return err
        }
        uStat, err := existDom.Delete()
        if err != nil {
                fmt.Println(err.Error())
                return err
        }
        b, err := json.Marshal(uStat)
        if err != nil {
                fmt.Println(err.Error())
                return err
        }
        fmt.Println(string(b))
        log.Printf("[DEBUG] [Akamai GTMV1_3] Delete status:")
        log.Printf("[DEBUG] [Akamai GTMV1_3] %v", b)

        if d.Get("wait_on_complete").(bool) {
                done, err := waitForCompletion(d.Id())
                if done {
                        log.Printf("[INFO] [Akamai GTMV1_3] Domain delete completed")
                } else {
                        if err == nil {
                                log.Printf("[INFO] [Akamai GTMV1_3] Domain delete pending")
                        } else {
                                log.Printf("[WARNING] [Akamai GTMV1_3] Domain delete failed [%s]", err.Error())
                                return err
                        }
                }

        }

        d.SetId("")
	return nil

	/*
	// No GTM Domain delete operation permitted.

	return schema.Noop(d, meta)
	*/

}

// Test GTM Domain existance
func resourceGTMv1_3DomainExists(d *schema.ResourceData, meta interface{}) (bool, error) {

	name := d.Get("name").(string)
	log.Printf("[DEBUG] [Akamai GTMV1_3] Searching for domain [%s]", name)
	domain, err := gtmv1_3.GetDomain(name)
	log.Printf("[DEBUG] [Akamai GTMV1_3] Searching for Existing domain result [%v]", domain)
	return domain != nil, err
}

// validateDomainType is a SchemaValidateFunc to validate the Domain type.
func validateDomainType(v interface{}, k string) (ws []string, es []error) {
	value := strings.ToUpper(v.(string))
	if value != "BASIC" && value != "FULL" && value != "WEIGHTED" && value != "STATIC" && value != "FAILOVER-ONLY" {
		es = append(es, fmt.Errorf("Type must be basic, full, weighted, static, or failover-only"))
	}
	return
}

// Create and populate a new domain object from resource data
func populateNewDomainObject(d *schema.ResourceData) *gtmv1_3.Domain {

	domObj := gtmv1_3.NewDomain(d.Get("name").(string), d.Get("type").(string))
	populateDomainObject(d, domObj)

	return domObj

}

// Populate existing domain object from resource data
func populateDomainObject(d *schema.ResourceData, dom *gtmv1_3.Domain) {

	if d.Get("name").(string) != dom.Name {
		dom.Name = d.Get("name").(string)
		log.Printf("[WARNING] [Akamai GTMV1_3] Domain [%s] state and GTM names inconsistent!", dom.Name)
	}
	if v, ok := d.GetOk("type"); ok {
		if v != dom.Type {
			dom.Type = v.(string)
		}
	}
	if v, ok := d.GetOk("default_unreachable_threshold"); ok {
		dom.DefaultUnreachableThreshold = v.(float32)
	}
	if v, ok := d.GetOk("email_notification_list"); ok {
		ls := make([]string, len(v.([]interface{})))
		for i, sl := range v.([]interface{}) {
			ls[i] = sl.(string)
		}
		dom.EmailNotificationList = ls
	}
	if v, ok := d.GetOk("min_pingable_region_fraction"); ok {
		dom.MinPingableRegionFraction = v.(float32)
	}
	if v, ok := d.GetOk("default_timeout_penalty"); ok {
		dom.DefaultTimeoutPenalty = v.(int)
	}
	if v, ok := d.GetOk("servermonitor_liveness_count"); ok {
		dom.ServermonitorLivenessCount = v.(int)
	}
	if v, ok := d.GetOk("round_robin_prefix"); ok {
		dom.RoundRobinPrefix = v.(string)
	}
	if v, ok := d.GetOk("servermonitor_load_count"); ok {
		dom.ServermonitorLoadCount = v.(int)
	}
	if v, ok := d.GetOk("ping_interval"); ok {
		dom.PingInterval = v.(int)
	}
	if v, ok := d.GetOk("max_ttl"); ok {
		dom.MaxTTL = int64(v.(int))
	}
	if v, ok := d.GetOk("load_imbalance_percentage"); ok {
		dom.LoadImbalancePercentage = v.(float64)
	}
	if v, ok := d.GetOk("default_health_max"); ok {
		dom.DefaultHealthMax = v.(int)
	}
	if v, ok := d.GetOk("map_update_interval"); ok {
		dom.MapUpdateInterval = v.(int)
	}
	if v, ok := d.GetOk("max_properties"); ok {
		dom.MaxProperties = v.(int)
	}
	if v, ok := d.GetOk("max_resources"); ok {
		dom.MaxResources = v.(int)
	}
	if v, ok := d.GetOk("default_ssl_client_private_key"); ok {
		dom.DefaultSslClientPrivateKey = v.(string)
	}
	if v, ok := d.GetOk("default_error_penalty"); ok {
		dom.DefaultErrorPenalty = v.(int)
	}
	if v, ok := d.GetOk("max_test_timeout"); ok {
		dom.MaxTestTimeout = v.(float64)
	}
	if v, ok := d.GetOk("cname_coalescing_enabled"); ok {
		dom.CnameCoalescingEnabled = v.(bool)
	}
	if v, ok := d.GetOk("default_health_multiplier"); ok {
		dom.DefaultHealthMultiplier = v.(int)
	}
	if v, ok := d.GetOk("servermonitor_pool"); ok {
		dom.ServermonitorPool = v.(string)
	}
	if v, ok := d.GetOk("load_feedback"); ok {
		dom.LoadFeedback = v.(bool)
	}
	if v, ok := d.GetOk("min_ttl"); ok {
		dom.MinTTL = int64(v.(int))
	}
	if v, ok := d.GetOk("default_max_unreachable_penalty"); ok {
		dom.DefaultMaxUnreachablePenalty = v.(int)
	}
	if v, ok := d.GetOk("default_health_threshold"); ok {
		dom.DefaultHealthThreshold = v.(int)
	}
	// Want??
	//if v, ok := d.GetOk("last_modified_by"); ok { dom.LastModifiedBy = v.(string) }
	// Want?
	if v, ok := d.GetOk("modification_comments"); ok {
		dom.ModificationComments = v.(string)
	}
	if v, ok := d.GetOk("min_test_interval"); ok {
		dom.MinTestInterval = v.(int)
	}
	if v, ok := d.GetOk("ping_packet_size"); ok {
		dom.PingPacketSize = v.(int)
	}
	if v, ok := d.GetOk("default_ssl_client_certificate"); ok {
		dom.DefaultSslClientCertificate = v.(string)
	}

	return

}

// Populate Terraform state from provided Domain object
func populateTerraformState(d *schema.ResourceData, dom *gtmv1_3.Domain) {

	// walk thru all state elements
	d.Set("name", dom.Name)
	d.Set("type", dom.Type)
	d.Set("default_unreachable_threshold", dom.DefaultUnreachableThreshold)
	d.Set("email_notification_list", dom.EmailNotificationList)
	d.Set("min_pingable_region_fraction", dom.MinPingableRegionFraction)
	d.Set("default_timeout_penalty", dom.DefaultTimeoutPenalty)
	d.Set("servermonitor_liveness_count", dom.ServermonitorLivenessCount)
	d.Set("round_robin_prefix", dom.RoundRobinPrefix)
	d.Set("servermonitor_load_count", dom.ServermonitorLoadCount)
	d.Set("ping_interval", dom.PingInterval)
	d.Set("max_ttl", dom.MaxTTL)
	d.Set("load_imbalance_percentage", dom.LoadImbalancePercentage)
	d.Set("default_health_max", dom.DefaultHealthMax)
	d.Set("map_update_interval", dom.MapUpdateInterval)
	d.Set("max_properties", dom.MaxProperties)
	d.Set("max_resources", dom.MaxResources)
	d.Set("default_ssl_client_private_key", dom.DefaultSslClientPrivateKey)
	d.Set("default_error_penalty", dom.DefaultErrorPenalty)
	d.Set("max_test_timeout", dom.MaxTestTimeout)
	d.Set("cname_coalescing_enabled", dom.CnameCoalescingEnabled)
	d.Set("default_health_multiplier", dom.DefaultHealthMultiplier)
	d.Set("servermonitor_pool", dom.ServermonitorPool)
	d.Set("load_feedback", dom.LoadFeedback)
	d.Set("min_ttl", dom.MinTTL)
	d.Set("default_max_unreachable_penalty", dom.DefaultMaxUnreachablePenalty)
	d.Set("default_health_threshold", dom.DefaultHealthThreshold)
	// Want??
	//d.Set("last_modified_by", dom.LastModifiedBy)
	// Want?
	d.Set("modification_comments", dom.ModificationComments)
	d.Set("min_test_interval", dom.MinTestInterval)
	d.Set("ping_packet_size", dom.PingPacketSize)
	d.Set("default_ssl_client_certificate", dom.DefaultSslClientCertificate)

	return

}

// Util function to wait for change deployment. return true if complete. false if not - error or nil (timeout)
func waitForCompletion(domain string) (bool, error) {

	return true, nil

	var defaultInterval int = 5
	var defaultTimeout int = 300
	var sleepInterval time.Duration = 1 // seconds. TODO:Should be configurable by user ...
	var sleepTimeout time.Duration = 1  // seconds. TODO: Should be configurable by user ...
	sleepInterval *= time.Duration(defaultInterval)
	sleepTimeout *= time.Duration(defaultTimeout)
	for {
		propStat, err := gtmv1_3.GetDomainStatus(domain)
		if err != nil {
			return false, err
		}
		if propStat.PropagationStatus == "COMPLETE" {
			return true, nil
		}
		if sleepTimeout <= 0 {
			return false, nil
		}
		time.Sleep(sleepInterval * time.Second)
		sleepTimeout -= sleepInterval
	}

	return false, errors.New("Unknown error while waiting for change completion") // don't know how/why we would have broken out.

}
