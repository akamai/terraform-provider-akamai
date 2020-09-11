package gtm

import (
	"fmt"
	"strconv"
	"strings"
	"sync"

	"github.com/akamai/terraform-provider-akamai/v2/pkg/akamai"
	"github.com/akamai/terraform-provider-akamai/v2/pkg/tools"

	gtm "github.com/akamai/AkamaiOPEN-edgegrid-golang/configgtm-v1_4"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/schema"
)

func resourceGTMv1Datacenter() *schema.Resource {
	return &schema.Resource{
		Create: resourceGTMv1DatacenterCreate,
		Read:   resourceGTMv1DatacenterRead,
		Update: resourceGTMv1DatacenterUpdate,
		Delete: resourceGTMv1DatacenterDelete,
		Exists: resourceGTMv1DatacenterExists,
		Importer: &schema.ResourceImporter{
			State: resourceGTMv1DatacenterImport,
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
			"nickname": {
				Type:     schema.TypeString,
				Optional: true,
			},
			"datacenter_id": {
				Type:     schema.TypeInt,
				Computed: true,
			},
			"city": {
				Type:     schema.TypeString,
				Optional: true,
			},
			"clone_of": {
				Type:     schema.TypeInt,
				Optional: true,
			},
			"cloud_server_host_header_override": {
				Type:     schema.TypeBool,
				Optional: true,
			},
			"cloud_server_targeting": {
				Type:     schema.TypeBool,
				Optional: true,
			},
			"default_load_object": {
				Type:     schema.TypeList,
				Optional: true,
				MaxItems: 1,
				Elem: &schema.Resource{
					Schema: map[string]*schema.Schema{
						"load_servers": {
							Type:     schema.TypeList,
							Elem:     &schema.Schema{Type: schema.TypeString},
							Optional: true,
						},
						"load_object": {
							Type:     schema.TypeString,
							Optional: true,
							Default:  "",
						},
						"load_object_port": {
							Type:     schema.TypeInt,
							Optional: true,
						},
					},
				},
			},
			"continent": {
				Type:     schema.TypeString,
				Optional: true,
			},
			"country": {
				Type:     schema.TypeString,
				Optional: true,
			},
			"latitude": {
				Type:     schema.TypeFloat,
				Optional: true,
			},
			"longitude": {
				Type:     schema.TypeFloat,
				Optional: true,
			},
			"ping_interval": {
				Type:     schema.TypeInt,
				Computed: true,
			},
			"ping_packet_size": {
				Type:     schema.TypeInt,
				Computed: true,
			},
			"score_penalty": {
				Type:     schema.TypeInt,
				Computed: true,
			},
			"servermonitor_liveness_count": {
				Type:     schema.TypeInt,
				Computed: true,
			},
			"servermonitor_load_count": {
				Type:     schema.TypeInt,
				Computed: true,
			},
			"servermonitor_pool": {
				Type:     schema.TypeString,
				Computed: true,
			},
			"state_or_province": {
				Type:     schema.TypeString,
				Optional: true,
			},
			"virtual": {
				Type:     schema.TypeBool,
				Computed: true,
			},
		},
	}
}

// utility func to parse Terraform DC resource id
func parseDatacenterResourceId(id string) (string, int, error) {

	parts := strings.SplitN(id, ":", 2)
	domain := parts[0]
	dcID, err := strconv.Atoi(parts[1])
	if len(parts) != 2 || parts[0] == "" || err != nil {
		return "", -1, err
	}

	return domain, dcID, nil
}

var (
	datacenterCreateLock sync.Mutex
)

// Create a new GTM Datacenter
func resourceGTMv1DatacenterCreate(d *schema.ResourceData, m interface{}) error {
	meta := akamai.Meta(m)
	logger := meta.Log("Akamai GTMv1", "resourceGTMv1DatacenterCreate")

	// Async GTM DC creation causes issues at this writing. Synchronize as work around.
	datacenterCreateLock.Lock()
	defer datacenterCreateLock.Unlock()

	domain, err := tools.GetStringValue("domain", d)
	if err != nil {
		logger.Warnf("dataSourceGTMDefaultDatacenterRead: Domain not initialized")
		return err
	}
	datacenterName, err := tools.GetStringValue("nickname", d)
	if err != nil {
		logger.Warnf("dataSourceGTMDefaultDatacenterRead: nickname not initialized")
		return err
	}

	logger.Infof("Creating datacenter [%s] in domain [%s]", datacenterName, domain)
	newDC := populateNewDatacenterObject(d, m)
	logger.Debugf("Proposed New Datacenter: [%v]", newDC)
	cStatus, err := newDC.Create(domain)
	if err != nil {
		logger.Errorf("DatacenterCreate failed: %s", err.Error())
		return err
	}
	logger.Debugf("Datacenter Create status: %v", cStatus.Status)
	if cStatus.Status.PropagationStatus == "DENIED" {
		return fmt.Errorf(cStatus.Status.Message)
	}

	waitOnComplete, err := tools.GetBoolValue("wait_on_complete", d)
	if err != nil {
		return err
	}

	if waitOnComplete {
		done, err := waitForCompletion(domain, m)
		if done {
			logger.Infof("Datacenter Create completed")
		} else {
			if err == nil {
				logger.Infof("Datacenter Create pending")
			} else {
				logger.Errorf("Datacenter Create failed [%s]", err.Error())
				return err
			}
		}

	}

	// Give terraform the ID. Format domain::dcid
	datacenterId := fmt.Sprintf("%s:%d", domain, cStatus.Resource.DatacenterId)
	logger.Debugf("Generated DC Resource Id: %s", datacenterId)
	d.SetId(datacenterId)
	return resourceGTMv1DatacenterRead(d, m)

}

// Only ever save data from the tf config in the tf state file, to help with
// api issues. See func unmarshalResourceData for more info.
func resourceGTMv1DatacenterRead(d *schema.ResourceData, m interface{}) error {
	meta := akamai.Meta(m)
	logger := meta.Log("Akamai GTMv1", "resourceGTMv1DatacenterRead")

	logger.Debugf("READ")
	logger.Debugf("Reading Datacenter: %s", d.Id())
	// retrieve the datacenter and domain
	domain, dcID, err := parseDatacenterResourceId(d.Id())
	if err != nil {
		return fmt.Errorf("invalid datacenter resource ID")
	}
	dc, err := gtm.GetDatacenter(dcID, domain)
	if err != nil {
		logger.Errorf("DatacenterRead failed: %s", err.Error())
		return err
	}
	populateTerraformDCState(d, dc, m)
	logger.Debugf("READ %v", dc)
	return nil
}

// Update GTM Datacenter
func resourceGTMv1DatacenterUpdate(d *schema.ResourceData, m interface{}) error {
	meta := akamai.Meta(m)
	logger := meta.Log("Akamai GTMv1", "resourceGTMv1DatacenterUpdate")

	logger.Debugf("UPDATE")
	logger.Debugf("Updating Datacenter: %s", d.Id())
	// pull domain and dcid out of resource id
	domain, dcID, err := parseDatacenterResourceId(d.Id())
	if err != nil {
		return fmt.Errorf("invalid datacenter resource ID")
	}
	// Get existing datacenter
	existDC, err := gtm.GetDatacenter(dcID, domain)
	if err != nil {
		logger.Errorf("DatacenterUpdate failed: %s", err.Error())
		return err
	}
	logger.Debugf("Updating Datacenter BEFORE: %v", existDC)
	populateDatacenterObject(d, existDC, m)
	logger.Debugf("Updating Datacenter PROPOSED: %v", existDC)
	uStat, err := existDC.Update(domain)
	if err != nil {
		logger.Errorf("DatacenterUpdate failed: %s", err.Error())
		return err
	}
	logger.Debugf("Datacenter Update  status:")
	logger.Debugf("%v", uStat)
	if uStat.PropagationStatus == "DENIED" {
		return fmt.Errorf(uStat.Message)
	}

	waitOnComplete, err := tools.GetBoolValue("wait_on_complete", d)
	if err != nil {
		return err
	}

	if waitOnComplete {
		done, err := waitForCompletion(domain, m)
		if done {
			logger.Infof("Datacenter update completed")
		} else {
			if err == nil {
				logger.Infof("Datacenter update pending")
			} else {
				logger.Errorf("Datacenter update failed [%s]", err.Error())
				return err
			}
		}

	}

	return resourceGTMv1DatacenterRead(d, m)
}

func resourceGTMv1DatacenterImport(d *schema.ResourceData, m interface{}) ([]*schema.ResourceData, error) {
	meta := akamai.Meta(m)
	logger := meta.Log("Akamai GTMv1", "resourceGTMv1DatacenterImport")

	logger.Debugf("Import")
	logger.Debugf("Importing Datacenter: %s", d.Id())
	// retrieve the datacenter and domain
	domain, dcID, err := parseDatacenterResourceId(d.Id())
	if err != nil {
		return nil, fmt.Errorf("invalid datacenter resource ID")
	}
	dc, err := gtm.GetDatacenter(dcID, domain)
	if err != nil {
		logger.Errorf("DatacenterImport error: %s", err.Error())
		return nil, err
	}
	populateTerraformDCState(d, dc, m)
	if err := d.Set("domain", domain); err != nil {
		return nil, err
	}
	if err := d.Set("wait_on_complete", true); err != nil {
		return nil, err
	}
	logger.Debugf("Import %v", dc)
	return []*schema.ResourceData{d}, err

}

// Delete GTM Datacenter.
func resourceGTMv1DatacenterDelete(d *schema.ResourceData, m interface{}) error {
	meta := akamai.Meta(m)
	logger := meta.Log("Akamai GTMv1", "resourceGTMv1DatacenterDelete")

	logger.Debugf("DELETE")
	logger.Debugf("Deleting Datacenter: %s", d.Id())
	domain, dcID, err := parseDatacenterResourceId(d.Id())
	if err != nil {
		return fmt.Errorf("invalid datacenter resource ID")
	}
	// Get existing datacenter
	existDC, err := gtm.GetDatacenter(dcID, domain)
	if err != nil {
		logger.Errorf("DatacenterDelete failed: %s", err.Error())
		return err
	}
	logger.Debugf("Deleting Datacenter: %v", existDC)
	uStat, err := existDC.Delete(domain)
	if err != nil {
		logger.Errorf("DatacenterDelete failed: %s", err.Error())
		return err
	}
	logger.Debugf("Datacenter Delete status:")
	logger.Debugf("%v", uStat)
	if uStat.PropagationStatus == "DENIED" {
		return fmt.Errorf(uStat.Message)
	}

	waitOnComplete, err := tools.GetBoolValue("wait_on_complete", d)
	if err != nil {
		return err
	}

	if waitOnComplete {
		done, err := waitForCompletion(domain, m)
		if done {
			logger.Infof("Datacenter delete completed")
		} else {
			if err == nil {
				logger.Infof("Datacenter delete pending")
			} else {
				logger.Errorf("Datacenter delete failed [%s]", err.Error())
				return err
			}
		}

	}

	// if successful ....
	d.SetId("")
	return nil
}

// Test GTM Datacenter existence
func resourceGTMv1DatacenterExists(d *schema.ResourceData, m interface{}) (bool, error) {
	meta := akamai.Meta(m)
	logger := meta.Log("Akamai GTMv1", "resourceGTMv1DatacenterExists")

	logger.Debugf("Exists")
	// pull domain and dcid out of resource id
	domain, dcID, err := parseDatacenterResourceId(d.Id())
	if err != nil {
		return false, fmt.Errorf("invalid datacenter resource ID")
	}
	logger.Debugf("Searching for existing datacenter [%d] in domain %s", dcID, domain)
	dc, err := gtm.GetDatacenter(dcID, domain)
	logger.Debugf("Searching for Existing datacenter result [%v]", domain)
	return dc != nil, err
}

// Create and populate a new datacenter object from resource data
func populateNewDatacenterObject(d *schema.ResourceData, m interface{}) *gtm.Datacenter {

	dcObj := gtm.NewDatacenter()
	dcObj.DefaultLoadObject = gtm.NewLoadObject()
	populateDatacenterObject(d, dcObj, m)

	return dcObj
}

// Populate existing datacenter object from resource data
func populateDatacenterObject(d *schema.ResourceData, dc *gtm.Datacenter, m interface{}) {
	meta := akamai.Meta(m)
	logger := meta.Log("Akamai GTMv1", "populateDatacenterObject")

	if v, err := tools.GetStringValue("nickname", d); err == nil {
		dc.Nickname = v
	}
	if v, err := tools.GetStringValue("city", d); err == nil || d.HasChange("city") {
		dc.City = v
		if err != nil {
			logger.Warnf(fmt.Sprintf("populateDataCenterObject() failed: %v", err.Error()))
		}
	}

	if v, err := tools.GetIntValue("clone_of", d); err == nil || d.HasChange("clone_of") {
		dc.CloneOf = v
		if err != nil {
			logger.Warnf(fmt.Sprintf("populateDataCenterObject() failed: %v", err.Error()))
		}
	}
	cloudServerHostHeaderOverride, err := tools.GetBoolValue("cloud_server_host_header_override", d)
	if err != nil {
		logger.Warnf(fmt.Sprintf("populateDataCenterObject() failed: cloud_server_host_header_override not set: %v", err.Error()))
	}
	dc.CloudServerHostHeaderOverride = cloudServerHostHeaderOverride
	cloudServerTargeting, err := tools.GetBoolValue("cloud_server_targeting", d)
	if err != nil {
		logger.Warnf("cloud_server_targeting not set: %s", err.Error())
	}
	dc.CloudServerTargeting = cloudServerTargeting
	if v, err := tools.GetStringValue("continent", d); err == nil || d.HasChange("continent") {
		dc.Continent = v
		if err != nil {
			logger.Warnf(fmt.Sprintf("populateDataCenterObject() failed: %v", err.Error()))
		}
	}
	if v, err := tools.GetStringValue("country", d); err == nil || d.HasChange("country") {
		dc.Country = v
		if err != nil {
			logger.Warnf(fmt.Sprintf("populateDataCenterObject() failed: %v", err.Error()))
		}
	}
	// pull apart Set
	if dloList, err := tools.GetInterfaceArrayValue("default_load_object", d); err != nil || len(dloList) == 0 {
		dc.DefaultLoadObject = nil
	} else {
		dloObject := gtm.NewLoadObject()
		dloMap, ok := dloList[0].(map[string]interface{})
		if !ok {
			logger.Errorf("populateDatacenterObject failed")
		}
		dloObject.LoadObject, ok = dloMap["load_object"].(string)
		if !ok {
			logger.Errorf("populateDatacenterObject failed, bad load_object format")
		}
		dloObject.LoadObjectPort, ok = dloMap["load_object_port"].(int)
		if !ok {
			logger.Errorf("populateDatacenterObject failed, bad load_object_port format")
		}
		loadServers, ok := dloMap["load_servers"]
		if ok {
			servers, ok := loadServers.([]interface{})
			if ok {
				dloObject.LoadServers = make([]string, len(servers))
				for i, server := range servers {
					if dloObject.LoadServers[i], ok = server.(string); !ok {
						logger.Errorf("populateDatacenterObject failed, bad loadServer format: %s", server)
					}
				}
			} else {
				logger.Errorf("populateDatacenterObject failed, bad load_servers format: %s", loadServers)
			}
		} else {
			logger.Errorf("populateDatacenterObject failed, load_servers not present")
		}
		dc.DefaultLoadObject = dloObject
	}
	if v, err := tools.GetFloat64Value("latitude", d); err == nil || d.HasChange("latitude") {
		dc.Latitude = v
		if err != nil {
			logger.Warnf(fmt.Sprintf("populateDataCenterObject() failed: %v", err.Error()))
		}
	}
	if v, err := tools.GetFloat64Value("longitude", d); err == nil || d.HasChange("longitude") {
		dc.Longitude = v
		if err != nil {
			logger.Warnf(fmt.Sprintf("populateDataCenterObject() failed: %v", err.Error()))
		}
	}
	if v, err := tools.GetStringValue("nickname", d); err == nil {
		dc.Nickname = v
	}
	if v, err := tools.GetIntValue("ping_interval", d); err == nil {
		dc.PingInterval = v
	}
	if v, err := tools.GetIntValue("ping_packet_size", d); err == nil {
		dc.PingPacketSize = v
	}
	if v, err := tools.GetIntValue("datacenter_id", d); err == nil {
		dc.DatacenterId = v
	}
	if v, err := tools.GetIntValue("score_penalty", d); err == nil {
		dc.ScorePenalty = v
	}
	if v, err := tools.GetIntValue("servermonitor_liveness_count", d); err == nil || d.HasChange("servermonitor_liveness_count") {
		dc.ServermonitorLivenessCount = v
		if err != nil {
			logger.Warnf(fmt.Sprintf("populateDataCenterObject() failed: %v", err.Error()))
		}
	}
	if v, err := tools.GetIntValue("servermonitor_load_count", d); err == nil || d.HasChange("servermonitor_load_count") {
		dc.ServermonitorLoadCount = v
		if err != nil {
			logger.Warnf(fmt.Sprintf("populateDataCenterObject() failed: %v", err.Error()))
		}
	}
	if v, err := tools.GetStringValue("servermonitor_pool", d); err == nil || d.HasChange("servermonitor_pool") {
		dc.ServermonitorPool = v
		if err != nil {
			logger.Warnf(fmt.Sprintf("populateDataCenterObject() failed: %v", err.Error()))
		}
	}
	if v, err := tools.GetStringValue("state_or_province", d); err == nil || d.HasChange("state_or_province") {
		dc.StateOrProvince = v
		if err != nil {
			logger.Warnf(fmt.Sprintf("populateDataCenterObject() failed: %v", err.Error()))
		}
	}
}

// Populate Terraform state from provided Datacenter object
func populateTerraformDCState(d *schema.ResourceData, dc *gtm.Datacenter, m interface{}) {
	meta := akamai.Meta(m)
	logger := meta.Log("Akamai GTMv1", "populateTerrafomDCState")

	// walk through all state elements
	for stateKey, stateValue := range map[string]interface{}{
		"nickname":                          dc.Nickname,
		"datacenter_id":                     dc.DatacenterId,
		"city":                              dc.City,
		"clone_of":                          dc.CloneOf,
		"cloud_server_host_header_override": dc.CloudServerHostHeaderOverride,
		"cloud_server_targeting":            dc.CloudServerTargeting,
		"continent":                         dc.Continent,
		"country":                           dc.Country} {
		err := d.Set(stateKey, stateValue)
		if err != nil {
			logger.Errorf("populateTerraformDCState failed: %s", err.Error())
		}
	}
	dloStateList, err := tools.GetInterfaceArrayValue("default_load_object", d)
	if err != nil {
		dloStateList = make([]interface{}, 0, 1)
	}
	if len(dloStateList) == 0 && dc.DefaultLoadObject != nil && (dc.DefaultLoadObject.LoadObject != "" || len(dc.DefaultLoadObject.LoadServers) != 0 || dc.DefaultLoadObject.LoadObjectPort > 0) {
		// create MT object
		newDLO := make(map[string]interface{}, 3)
		newDLO["load_object"] = ""
		newDLO["load_object_port"] = 0
		newDLO["load_servers"] = make([]interface{}, 0, len(dc.DefaultLoadObject.LoadServers))
		dloStateList = append(dloStateList, newDLO)
	}
	for _, dloMap := range dloStateList {
		if dc.DefaultLoadObject != nil && (dc.DefaultLoadObject.LoadObject != "" || len(dc.DefaultLoadObject.LoadServers) != 0 || dc.DefaultLoadObject.LoadObjectPort > 0) {
			dlo := dloMap.(map[string]interface{})
			dlo["load_object"] = dc.DefaultLoadObject.LoadObject
			dlo["load_object_port"] = dc.DefaultLoadObject.LoadObjectPort
			if dlo["load_servers"] != nil && len(dlo["load_servers"].([]interface{})) > 0 {
				dlo["load_servers"] = reconcileTerraformLists(dlo["load_servers"].([]interface{}), convertStringToInterfaceList(dc.DefaultLoadObject.LoadServers, m), m)
			} else {
				dlo["load_servers"] = dc.DefaultLoadObject.LoadServers
			}
		} else {
			dloStateList = make([]interface{}, 0, 1)
		}
	}
	for stateKey, stateValue := range map[string]interface{}{
		"default_load_object":          dloStateList,
		"latitude":                     dc.Latitude,
		"longitude":                    dc.Longitude,
		"ping_interval":                dc.PingInterval,
		"ping_packet_size":             dc.PingPacketSize,
		"score_penalty":                dc.ScorePenalty,
		"servermonitor_liveness_count": dc.ServermonitorLivenessCount,
		"servermonitor_load_count":     dc.ServermonitorLoadCount,
		"servermonitor_pool":           dc.ServermonitorPool,
		"state_or_province":            dc.StateOrProvince,
		"virtual":                      dc.Virtual,
	} {
		err := d.Set(stateKey, stateValue)
		if err != nil {
			logger.Errorf("populateTerraformDCState failed: %s", err.Error())
		}
	}

}
