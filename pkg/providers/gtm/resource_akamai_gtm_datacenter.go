package gtm

import (
	"fmt"
	"log"
	"strconv"
	"strings"
	"sync"

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
			"default_load_object": &schema.Schema{
				Type:     schema.TypeList,
				Optional: true,
				MaxItems: 1,
				Elem: &schema.Resource{
					Schema: map[string]*schema.Schema{
						"load_servers": &schema.Schema{
							Type:     schema.TypeList,
							Elem:     &schema.Schema{Type: schema.TypeString},
							Optional: true,
						},
						"load_object": &schema.Schema{
							Type:     schema.TypeString,
							Optional: true,
							Default:  "",
						},
						"load_object_port": &schema.Schema{
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
	dcID, err := strconv.Atoi(parts[1])
	if len(parts) != 2 || parts[0] == "" || err != nil {
		return "", -1, err
	}

	return parts[0], dcID, err

}

var (
	datacenterCreateLock sync.Mutex
)

// Create a new GTM Datacenter
func resourceGTMv1DatacenterCreate(d *schema.ResourceData, meta interface{}) error {

	// Async GTM DC creation causes issues at this writing. Synchronize as work around.
	datacenterCreateLock.Lock()
	defer datacenterCreateLock.Unlock()

	domain := d.Get("domain").(string)

	log.Printf("[INFO] [Akamai GTM] Creating datacenter [%s] in domain [%s]", d.Get("nickname").(string), domain)
	newDC := populateNewDatacenterObject(d)
	log.Printf("[DEBUG] [Akamai GTMv1] Proposed New Datacenter: [%v]", newDC)
	cStatus, err := newDC.Create(domain)
	if err != nil {
		log.Printf("[ERROR] DatacenterCreate failed: %s", err.Error())
		return err
	}
	log.Printf("[DEBUG] [Akamai GTMv1] Datacenter Create status:")
	log.Printf("[DEBUG] [Akamai GTMv1] %v", cStatus.Status)
	if cStatus.Status.PropagationStatus == "DENIED" {
		return fmt.Errorf(cStatus.Status.Message)
	}
	if d.Get("wait_on_complete").(bool) {
		done, err := waitForCompletion(domain)
		if done {
			log.Printf("[INFO] [Akamai GTMv1] Datacenter Create completed")
		} else {
			if err == nil {
				log.Printf("[INFO] [Akamai GTMv1] Datacenter Create pending")
			} else {
				log.Printf("[WARNING] [Akamai GTMv1] Datacenter Create failed [%s]", err.Error())
				return err
			}
		}

	}

	// Give terraform the ID. Format domain::dcid
	datacenterId := fmt.Sprintf("%s:%d", domain, cStatus.Resource.DatacenterId)
	log.Printf("[DEBUG] [Akamai GTMv1] Generated DC Resource Id: %s", datacenterId)
	d.SetId(datacenterId)
	return resourceGTMv1DatacenterRead(d, meta)

}

// Only ever save data from the tf config in the tf state file, to help with
// api issues. See func unmarshalResourceData for more info.
func resourceGTMv1DatacenterRead(d *schema.ResourceData, meta interface{}) error {

	log.Printf("[DEBUG] [Akamai GTMv1] READ")
	log.Printf("[DEBUG] Reading [Akamai GTMv1] Datacenter: %s", d.Id())
	// retrieve the datacenter and domain
	domain, dcID, err := parseDatacenterResourceId(d.Id())
	if err != nil {
		return fmt.Errorf("Invalid datacenter resource Id")
	}
	dc, err := gtm.GetDatacenter(dcID, domain)
	if err != nil {
		log.Printf("[ERROR] DatacenterRead failed: %s", err.Error())
		return err
	}
	populateTerraformDCState(d, dc)
	log.Printf("[DEBUG] [Akamai GTMv1] READ %v", dc)
	return nil
}

// Update GTM Datacenter
func resourceGTMv1DatacenterUpdate(d *schema.ResourceData, meta interface{}) error {

	log.Printf("[DEBUG] [Akamai GTMv1] UPDATE")
	log.Printf("[DEBUG] Updating [Akamai GTMv1] Datacenter: %s", d.Id())
	// pull domain and dcid out of resource id
	domain, dcID, err := parseDatacenterResourceId(d.Id())
	if err != nil {
		return fmt.Errorf("Invalid datacenter resource Id")
	}
	// Get existing datacenter
	existDC, err := gtm.GetDatacenter(dcID, domain)
	if err != nil {
		log.Printf("[ERROR] DatacenterUpdate failed: %s", err.Error())
		return err
	}
	log.Printf("[DEBUG] Updating [Akamai GTMv1] Datacenter BEFORE: %v", existDC)
	populateDatacenterObject(d, existDC)
	log.Printf("[DEBUG] Updating [Akamai GTMv1] Datacenter PROPOSED: %v", existDC)
	uStat, err := existDC.Update(domain)
	if err != nil {
		log.Printf("[ERROR] DatacenterUpdate failed: %s", err.Error())
		return err
	}
	log.Printf("[DEBUG] [Akamai GTMv1] Datacenter Update  status:")
	log.Printf("[DEBUG] [Akamai GTMv1] %v", uStat)
	if uStat.PropagationStatus == "DENIED" {
		return fmt.Errorf(uStat.Message)
	}
	if d.Get("wait_on_complete").(bool) {
		done, err := waitForCompletion(domain)
		if done {
			log.Printf("[INFO] [Akamai GTMv1] Datacenter update completed")
		} else {
			if err == nil {
				log.Printf("[INFO] [Akamai GTMv1] Datacenter update pending")
			} else {
				log.Printf("[WARNING] [Akamai GTMv1] Datacenter update failed [%s]", err.Error())
				return err
			}
		}

	}

	return resourceGTMv1DatacenterRead(d, meta)
}

func resourceGTMv1DatacenterImport(d *schema.ResourceData, meta interface{}) ([]*schema.ResourceData, error) {

	log.Printf("[DEBUG] [Akamai GTMv1] Import")
	log.Printf("[DEBUG] Importing [Akamai GTMv1] Datacenter: %s", d.Id())
	// retrieve the datacenter and domain
	domain, dcID, err := parseDatacenterResourceId(d.Id())
	if err != nil {
		return nil, fmt.Errorf("Invalid datacenter resource Id")
	}
	dc, err := gtm.GetDatacenter(dcID, domain)
	if err != nil {
		log.Printf("[ERROR] DatacenterImport error: %s", err.Error())
		return nil, err
	}
	populateTerraformDCState(d, dc)
	d.Set("domain", domain)
	d.Set("wait_on_complete", true)
	log.Printf("[DEBUG] [Akamai GTMv1] Import %v", dc)
	return []*schema.ResourceData{d}, err

}

// Delete GTM Datacenter.
func resourceGTMv1DatacenterDelete(d *schema.ResourceData, meta interface{}) error {

	log.Printf("[DEBUG] [Akamai GTMv1] DELETE")
	log.Printf("[DEBUG] Deleting [Akamai GTMv1] Datacenter: %s", d.Id())
	domain, dcID, err := parseDatacenterResourceId(d.Id())
	if err != nil {
		return fmt.Errorf("Invalid datacenter resource Id")
	}
	// Get existing datacenter
	existDC, err := gtm.GetDatacenter(dcID, domain)
	if err != nil {
		log.Printf("[ERROR] DatacenterDelete failed: %s", err.Error())
		return err
	}
	log.Printf("[DEBUG] Deleting [Akamai GTMv1] Datacenter: %v", existDC)
	uStat, err := existDC.Delete(domain)
	if err != nil {
		log.Printf("[ERROR] DatacenterDelete failed: %s", err.Error())
		return err
	}
	log.Printf("[DEBUG] [Akamai GTMv1] Datacenter Delete status:")
	log.Printf("[DEBUG] [Akamai GTMv1] %v", uStat)
	if uStat.PropagationStatus == "DENIED" {
		return fmt.Errorf(uStat.Message)
	}
	if d.Get("wait_on_complete").(bool) {
		done, err := waitForCompletion(domain)
		if done {
			log.Printf("[INFO] [Akamai GTMv1] Datacenter delete completed")
		} else {
			if err == nil {
				log.Printf("[INFO] [Akamai GTMv1] Datacenter delete pending")
			} else {
				log.Printf("[WARNING] [Akamai GTMv1] Datacenter delete failed [%s]", err.Error())
				return err
			}
		}

	}

	// if succcessful ....
	d.SetId("")
	return nil
}

// Test GTM Datacenter existance
func resourceGTMv1DatacenterExists(d *schema.ResourceData, meta interface{}) (bool, error) {

	log.Printf("[DEBUG] [Akamai GTMv1] Exists")
	// pull domain and dcid out of resource id
	domain, dcID, err := parseDatacenterResourceId(d.Id())
	if err != nil {
		return false, fmt.Errorf("Invalid datacenter resource Id")
	}
	log.Printf("[DEBUG] [Akamai GTMv1] Searching for existing datacenter [%d] in domain %s", dcID, domain)
	dc, err := gtm.GetDatacenter(dcID, domain)
	log.Printf("[DEBUG] [Akamai GTMv1] Searching for Existing datacenter result [%v]", domain)
	return dc != nil, err
}

// Create and populate a new datacenter object from resource data
func populateNewDatacenterObject(d *schema.ResourceData) *gtm.Datacenter {

	dcObj := gtm.NewDatacenter()
	dcObj.DefaultLoadObject = gtm.NewLoadObject()
	populateDatacenterObject(d, dcObj)

	return dcObj

}

// Populate existing datacenter object from resource data
func populateDatacenterObject(d *schema.ResourceData, dc *gtm.Datacenter) {

	if v, ok := d.GetOk("nickname"); ok {
		dc.Nickname = v.(string)
	}
	if v, ok := d.GetOk("city"); ok {
		dc.City = v.(string)
	} else if d.HasChange("city") {
		dc.City = v.(string)
	}
	if v, ok := d.GetOk("clone_of"); ok {
		dc.CloneOf = v.(int)
	} else if d.HasChange("clone_of") {
		dc.CloneOf = v.(int)
	}
	v := d.Get("cloud_server_host_header_override")
	dc.CloudServerHostHeaderOverride = v.(bool)
	v = d.Get("cloud_server_targeting")
	dc.CloudServerTargeting = v.(bool)
	if v, ok := d.GetOk("continent"); ok {
		dc.Continent = v.(string)
	} else if d.HasChange("continent") {
		dc.Continent = v.(string)
	}
	if v, ok := d.GetOk("country"); ok {
		dc.Country = v.(string)
	} else if d.HasChange("country") {
		dc.Country = v.(string)
	}
	// pull apart Set
	dloList := d.Get("default_load_object").([]interface{})
	if dloList == nil || len(dloList) == 0 {
		dc.DefaultLoadObject = nil
	} else {
		dloObject := gtm.NewLoadObject()
		for _, v := range dloList {
			dloMap := v.(map[string]interface{})
			dloObject.LoadObject = dloMap["load_object"].(string)
			dloObject.LoadObjectPort = dloMap["load_object_port"].(int)
			if dloMap["load_servers"] != nil {
				ls := make([]string, len(dloMap["load_servers"].([]interface{})))
				for i, sl := range dloMap["load_servers"].([]interface{}) {
					ls[i] = sl.(string)
				}
				dloObject.LoadServers = ls
			}
			dc.DefaultLoadObject = dloObject
			break
		}
	}
	if v, ok := d.GetOk("latitude"); ok {
		dc.Latitude = v.(float64)
	} else if d.HasChange("latitude") {
		dc.Latitude = v.(float64)
	}
	if v, ok := d.GetOk("longitude"); ok {
		dc.Longitude = v.(float64)
	} else if d.HasChange("longitude") {
		dc.Longitude = v.(float64)
	}
	if v, ok := d.GetOk("nickname"); ok {
		dc.Nickname = v.(string)
	}
	if v, ok := d.GetOk("ping_interval"); ok {
		dc.PingInterval = v.(int)
	}
	if v, ok := d.GetOk("ping_packet_size"); ok {
		dc.PingPacketSize = v.(int)
	}
	if v, ok := d.GetOk("datacenter_id"); ok {
		dc.DatacenterId = v.(int)
	}
	if v, ok := d.GetOk("score_penalty"); ok {
		dc.ScorePenalty = v.(int)
	}
	if v, ok := d.GetOk("servermonitor_liveness_count"); ok {
		dc.ServermonitorLivenessCount = v.(int)
	} else if d.HasChange("servermonitor_liveness_count") {
		dc.ServermonitorLivenessCount = v.(int)
	}
	if v, ok := d.GetOk("servermonitor_load_count"); ok {
		dc.ServermonitorLoadCount = v.(int)
	} else if d.HasChange("servermonitor_load_count") {
		dc.ServermonitorLoadCount = v.(int)
	}
	if v, ok := d.GetOk("servermonitor_pool"); ok {
		dc.ServermonitorPool = v.(string)
	} else if d.HasChange("servermonitor_pool") {
		dc.ServermonitorPool = v.(string)
	}
	if v, ok := d.GetOk("state_or_province"); ok {
		dc.StateOrProvince = v.(string)
	} else if d.HasChange("state_or_province") {
		dc.StateOrProvince = v.(string)
	}
}

// Populate Terraform state from provided Datacenter object
func populateTerraformDCState(d *schema.ResourceData, dc *gtm.Datacenter) {

	// walk thru all state elements
	d.Set("nickname", dc.Nickname)
	d.Set("datacenter_id", dc.DatacenterId)
	d.Set("city", dc.City)
	d.Set("clone_of", dc.CloneOf)
	d.Set("cloud_server_host_header_override", dc.CloudServerHostHeaderOverride)
	d.Set("cloud_server_targeting", dc.CloudServerTargeting)
	d.Set("continent", dc.Continent)
	d.Set("country", dc.Country)
	dloStateList := d.Get("default_load_object").([]interface{})
	if dloStateList == nil {
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
				dlo["load_servers"] = reconcileTerraformLists(dlo["load_servers"].([]interface{}), convertStringToInterfaceList(dc.DefaultLoadObject.LoadServers))
			} else {
				dlo["load_servers"] = dc.DefaultLoadObject.LoadServers
			}
		} else {
			dloStateList = make([]interface{}, 0, 1)
		}
	}
	d.Set("default_load_object", dloStateList)
	d.Set("latitude", dc.Latitude)
	d.Set("longitude", dc.Longitude)
	d.Set("ping_interval", dc.PingInterval)
	d.Set("ping_packet_size", dc.PingPacketSize)
	d.Set("score_penalty", dc.ScorePenalty)
	d.Set("servermonitor_liveness_count", dc.ServermonitorLivenessCount)
	d.Set("servermonitor_load_count", dc.ServermonitorLoadCount)
	d.Set("servermonitor_pool", dc.ServermonitorPool)
	d.Set("state_or_province", dc.StateOrProvince)
	d.Set("virtual", dc.Virtual)

}
