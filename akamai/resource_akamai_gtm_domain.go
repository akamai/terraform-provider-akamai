package akamai

import (
	"errors"
	"fmt"
	"log"
	"strings"
	"time"

	client "github.com/akamai/AkamaiOPEN-edgegrid-golang/client-v1"
	gtm "github.com/akamai/AkamaiOPEN-edgegrid-golang/configgtm-v1_4"
	"github.com/hashicorp/terraform/helper/schema"
)

// Hack for Hashicorp Acceptance Tests
var HashiAcc = false

func resourceGTMv1Domain() *schema.Resource {
	return &schema.Resource{
		Create: resourceGTMv1DomainCreate,
		Read:   resourceGTMv1DomainRead,
		Update: resourceGTMv1DomainUpdate,
		Delete: resourceGTMv1DomainDelete,
		Exists: resourceGTMv1DomainExists,
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
				Type:     schema.TypeFloat,
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
				Type:     schema.TypeFloat,
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
				Type:     schema.TypeFloat,
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
			"end_user_mapping_enabled": {
				Type:     schema.TypeBool,
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
func resourceGTMv1DomainCreate(d *schema.ResourceData, meta interface{}) error {

	dname := d.Get("name").(string)
	log.Printf("[INFO] [Akamai GTM] Creating domain [%s]", dname)
	newDom := populateNewDomainObject(d)
	log.Printf("[DEBUG] [Akamai GTMv1] Domain: [%v]", newDom)

	cStatus, err := newDom.Create(GetQueryArgs(d))

	if err != nil {
		// Errored. Lets see if special hack
		if !HashiAcc {
			log.Printf("[ERROR] DomainCreate failed: %s", err.Error())
			return err
		}
		if _, ok := err.(gtm.CommonError); !ok {
			log.Printf("[ERROR] DomainCreate failed: %s", err.Error())
			return err
		}
		origErr, ok := err.(gtm.CommonError).GetItem("err").(client.APIError)
		if !ok {
			log.Printf("[ERROR] DomainCreate failed: %s", err.Error())
			return err
		}
		if origErr.Status == 400 && strings.Contains(origErr.RawBody, "proposed domain name") && strings.Contains(origErr.RawBody, "Domain Validation Error") {
			// Already exists
			log.Printf("[WARNING] [Akamai GTMv1] : Domain %s already exists. Ignoring error (Hashicorp).", dname)
		} else {
			log.Printf("[ERROR] [Akamai GTM] Error creating domain [%s]", err.Error())
			return err
		}
	} else {
		log.Printf("[DEBUG] [Akamai GTMv1] Create status:")
		log.Printf("[DEBUG] [Akamai GTMv1] %v", cStatus.Status)
		if cStatus.Status.PropagationStatus == "DENIED" {
			return errors.New(cStatus.Status.Message)
		}
		if d.Get("wait_on_complete").(bool) {
			done, err := waitForCompletion(dname)
			if done {
				log.Printf("[INFO] [Akamai GTMv1] Domain Create completed")
			} else {
				if err == nil {
					log.Printf("[INFO] [Akamai GTMv1] Domain Create pending")
				} else {
					log.Printf("[WARNING] [Akamai GTMv1] Domain Create failed [%s]", err.Error())
					return err
				}
			}
		}
	}
	// Give terraform the ID
	d.SetId(dname)
	return resourceGTMv1DomainRead(d, meta)

}

// Only ever save data from the tf config in the tf state file, to help with
// api issues. See func unmarshalResourceData for more info.
func resourceGTMv1DomainRead(d *schema.ResourceData, meta interface{}) error {
	log.Printf("[DEBUG] [Akamai GTMv1] READ")
	log.Printf("[DEBUG] Reading [Akamai GTMv1] Domain: %s", d.Id())
	// retrieve the domain
	dom, err := gtm.GetDomain(d.Id())
	if err != nil {
		log.Printf("[ERROR] DomainRead error: %s", err.Error())
		return err
	}
	populateTerraformState(d, dom)
	log.Printf("[DEBUG] [Akamai GTMv1] READ %v", dom)
	return nil
}

// Update GTM Domain
func resourceGTMv1DomainUpdate(d *schema.ResourceData, meta interface{}) error {

	log.Printf("[DEBUG] [Akamai GTMv1] UPDATE")
	log.Printf("[DEBUG] Updating [Akamai GTMv1] Domain: %s", d.Id())
	// Get existing domain
	existDom, err := gtm.GetDomain(d.Id())
	if err != nil {
		log.Printf("[ERROR] DomainUpdate failed: %s", err.Error())
		return err
	}
	log.Printf("[DEBUG] Updating [Akamai GTMv1] Domain BEFORE: %v", existDom)
	populateDomainObject(d, existDom)
	log.Printf("[DEBUG] Updating [Akamai GTMv1] Domain PROPOSED: %v", existDom)
	//existDom := populateNewDomainObject(d)
	uStat, err := existDom.Update(GetQueryArgs(d))
	if err != nil {
		log.Printf("[ERROR] DomainUpdate failed: %s", err.Error())
		return err
	}
	log.Printf("[DEBUG] [Akamai GTMv1] Update status:")
	log.Printf("[DEBUG] [Akamai GTMv1] %v", uStat)
	if uStat.PropagationStatus == "DENIED" {
		return errors.New(uStat.Message)
	}
	if d.Get("wait_on_complete").(bool) {
		done, err := waitForCompletion(d.Id())
		if done {
			log.Printf("[INFO] [Akamai GTMv1] Domain update completed")
		} else {
			if err == nil {
				log.Printf("[INFO] [Akamai GTMv1] Domain update pending")
			} else {
				log.Printf("[WARNING] [Akamai GTMv1] Domain update failed [%s]", err.Error())
				return err
			}
		}

	}

	return resourceGTMv1DomainRead(d, meta)

}

// Delete GTM Domain. Admin priviledges required in current API version.
func resourceGTMv1DomainDelete(d *schema.ResourceData, meta interface{}) error {
	log.Printf("[DEBUG] Deleting GTM Domain")
	log.Printf("[DEBUG] [Akamai GTMv1] Domain: %s", d.Id())
	// Get existing domain
	existDom, err := gtm.GetDomain(d.Id())
	if err != nil {
		log.Printf("[ERROR] DomainDelete failed: %s", err.Error())
		return err
	}
	uStat, err := existDom.Delete()
	if err != nil {
		// Errored. Lets see if special hack
		if !HashiAcc {
			log.Printf("[ERROR] Error DomainDelete: %s", err.Error())
			return err
		}
		if _, ok := err.(gtm.CommonError); !ok {
			log.Printf("[ERROR] Error DomainDelete: %s", err.Error())
			return err
		}
		origErr, ok := err.(gtm.CommonError).GetItem("err").(client.APIError)
		if !ok {
			log.Printf("[ERROR] Error DomainDelete: %s", err.Error())
			return err
		}
		if origErr.Status == 405 && strings.Contains(origErr.RawBody, "Bad Request") && strings.Contains(origErr.RawBody, "DELETE method is not supported") {
			log.Printf("[Warning] [Akamai GTMv1] : Domain %s delete failed.  Ignoring error (Hashicorp).", d.Id())
		} else {
			log.Printf("[ERROR] Error DomainDelete: %s", err.Error())
			return err
		}
	} else {
		log.Printf("[DEBUG] [Akamai GTMv1] Delete status:")
		log.Printf("[DEBUG] [Akamai GTMv1] %v", uStat)
		if uStat.PropagationStatus == "DENIED" {
			return errors.New(uStat.Message)
		}
		if d.Get("wait_on_complete").(bool) {
			done, err := waitForCompletion(d.Id())
			if done {
				log.Printf("[INFO] [Akamai GTMv1] Domain delete completed")
			} else {
				if err == nil {
					log.Printf("[INFO] [Akamai GTMv1] Domain delete pending")
				} else {
					log.Printf("[WARNING] [Akamai GTMv1] Domain delete failed [%s]", err.Error())
					return err
				}
			}
		}
	}
	d.SetId("")
	return nil

}

// Test GTM Domain existance
func resourceGTMv1DomainExists(d *schema.ResourceData, meta interface{}) (bool, error) {

	name := d.Get("name").(string)
	log.Printf("[DEBUG] [Akamai GTMv1] Searching for domain [%s]", name)
	domain, err := gtm.GetDomain(name)
	log.Printf("[DEBUG] [Akamai GTMv1] Searching for Existing domain result [%v]", domain)
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
func populateNewDomainObject(d *schema.ResourceData) *gtm.Domain {

	domObj := gtm.NewDomain(d.Get("name").(string), d.Get("type").(string))
	populateDomainObject(d, domObj)

	return domObj

}

// Populate existing domain object from resource data
func populateDomainObject(d *schema.ResourceData, dom *gtm.Domain) {

	if d.Get("name").(string) != dom.Name {
		dom.Name = d.Get("name").(string)
		log.Printf("[WARNING] [Akamai GTMv1] Domain [%s] state and GTM names inconsistent!", dom.Name)
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
	} else if d.HasChange("email_notification_list") {
		dom.EmailNotificationList = make([]string, 0, 0)
	}
	if v, ok := d.GetOk("min_pingable_region_fraction"); ok {
		dom.MinPingableRegionFraction = v.(float32)
	}
	if v, ok := d.GetOk("default_timeout_penalty"); ok {
		dom.DefaultTimeoutPenalty = v.(int)
	} else if d.HasChange("default_timeout_penalty") {
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
		dom.DefaultHealthMax = v.(float64)
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
	} else if d.HasChange("default_ssl_client_private_key") {
                dom.DefaultSslClientPrivateKey = v.(string)
	}
	if v, ok := d.GetOk("default_error_penalty"); ok {
		dom.DefaultErrorPenalty = v.(int)
	} else if d.HasChange("default_error_penalty") {
                dom.DefaultErrorPenalty = v.(int)
	}
	if v, ok := d.GetOk("max_test_timeout"); ok {
		dom.MaxTestTimeout = v.(float64)
	}
	v := d.Get("cname_coalescing_enabled")
	dom.CnameCoalescingEnabled = v.(bool)
	if v, ok := d.GetOk("default_health_multiplier"); ok {
		dom.DefaultHealthMultiplier = v.(float64)
	}
	if v, ok := d.GetOk("servermonitor_pool"); ok {
		dom.ServermonitorPool = v.(string)
	}
	v = d.Get("load_feedback")
	dom.LoadFeedback = v.(bool)
	if v, ok := d.GetOk("min_ttl"); ok {
		dom.MinTTL = int64(v.(int))
	}
	if v, ok := d.GetOk("default_max_unreachable_penalty"); ok {
		dom.DefaultMaxUnreachablePenalty = v.(int)
	}
	if v, ok := d.GetOk("default_health_threshold"); ok {
		dom.DefaultHealthThreshold = v.(float64)
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
	} else if d.HasChange("default_ssl_client_certificate") {
                dom.DefaultSslClientCertificate = v.(string)
	}
	if v, ok := d.GetOk("end_user_mapping_enabled"); ok {
		dom.EndUserMappingEnabled = v.(bool)
	}

}

// Populate Terraform state from provided Domain object
func populateTerraformState(d *schema.ResourceData, dom *gtm.Domain) {

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
	d.Set("end_user_mapping_enabled", dom.EndUserMappingEnabled)

}

// Util function to wait for change deployment. return true if complete. false if not - error or nil (timeout)
func waitForCompletion(domain string) (bool, error) {

	var defaultInterval time.Duration = 5 * time.Second
	var defaultTimeout time.Duration = 300 * time.Second
	var sleepInterval time.Duration = defaultInterval // seconds. TODO:Should be configurable by user ...
	var sleepTimeout time.Duration = defaultTimeout   // seconds. TODO: Should be configurable by user ...
	if HashiAcc {
		// Override for ACC tests
		sleepTimeout = sleepInterval
	}
	log.Printf("[DEBUG] [Akamai GTMv1] WAIT: Sleep Interval [%v]", sleepInterval/time.Second)
	log.Printf("[DEBUG] [Akamai GTMv1] WAIT: Sleep Timeout [%v]", sleepTimeout/time.Second)
	for {
		propStat, err := gtm.GetDomainStatus(domain)
		if err != nil {
			return false, err
		}
		log.Printf("[DEBUG] [Akamai GTMv1] WAIT: propStat.PropagationStatus [%v]", propStat.PropagationStatus)
		switch propStat.PropagationStatus {
		case "COMPLETE":
			log.Printf("[DEBUG] [Akamai GTMv1] WAIT: Return COMPLETE")
			return true, nil
		case "DENIED":
			log.Printf("[DEBUG] [Akamai GTMv1] WAIT: Return DENIED")
			return false, errors.New(propStat.Message)
		case "PENDING":
			if sleepTimeout <= 0 {
				log.Printf("[DEBUG] [Akamai GTMv1] WAIT: Return TIMED OUT")
				return false, nil
			}
			time.Sleep(sleepInterval)
			sleepTimeout -= sleepInterval
			log.Printf("[DEBUG] [Akamai GTMv1] WAIT: Sleep Time Remaining [%v]", sleepTimeout/time.Second)
		default:
			return false, errors.New("Unknown propagationStatus while waiting for change completion") // don't know how/why we would have broken out.
		}
	}
}
