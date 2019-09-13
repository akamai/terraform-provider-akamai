package akamai

import (
	"encoding/json"
	"fmt"
	"log"
	"strconv"
	"errors"
	"strings"
	gtmv1_3 "github.com/akamai/AkamaiOPEN-edgegrid-golang/configgtm-v1_3"
	"github.com/hashicorp/terraform/helper/schema"
)

func resourceGTMv1_3Datacenter() *schema.Resource {
	return &schema.Resource{
		Create: resourceGTMv1_3DatacenterCreate,
		Read:   resourceGTMv1_3DatacenterRead,
		Update: resourceGTMv1_3DatacenterUpdate,
		Delete: resourceGTMv1_3DatacenterDelete,
		Exists: resourceGTMv1_3DatacenterExists,
		Importer: &schema.ResourceImporter{
			State: resourceGTMv1_3DatacenterImport,
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
				Type:         schema.TypeString,
				Optional:     true,
			},
                        "clone_of": {
                                Type:     schema.TypeInt,
                                Optional: true,
                        },
 			"cloud_server_targeting": {
				Type:    schema.TypeBool,
                                Optional: true,
			},
			"default_load_object": &schema.Schema{
				Type:    schema.TypeSet,
				Optional: true,
				ConfigMode: schema.SchemaConfigModeAttr,
				MaxItems: 1,
				Elem:    &schema.Resource{
					Schema:map[string]*schema.Schema{
						"load_servers": &schema.Schema{
        	                        		Type:    schema.TypeList,
							Elem:    &schema.Schema{Type: schema.TypeString},
                        	        		Optional: true,
						},
						"load_object": &schema.Schema{
							Type:    schema.TypeString,
							Optional: true,
							Default: "",
						},
						"load_object_port": &schema.Schema{
							Type:    schema.TypeInt,
							Optional: true,
						},
				  	},
                        	},
			},
			"continent": {
                                Type:    schema.TypeString,
                                Optional: true,
                        },
			"country": {
                                Type:    schema.TypeString,
                                Optional: true,
                    	},
			"latitude": {
                                Type:    schema.TypeFloat,
                                Optional: true,
                        },
			"longitude": {
                                Type:     schema.TypeFloat,
                                Optional: true,
                        },
			"ping_interval": {
                                Type:    schema.TypeInt,
                                Optional: true,
                        },
			"ping_packet_size": {
                                Type:    schema.TypeInt,
                                Optional: true,
                        },
			"score_penalty": {
                                Type:    schema.TypeInt,
                                Optional: true,
                        },
			"servermonitor_liveness_count": {
                                Type:    schema.TypeInt,
                                Optional: true,
                        },
			"servermonitor_load_count": {
                                Type:    schema.TypeInt,
                                Optional: true,
                        },
			"servermonitor_pool": {
                                Type:    schema.TypeString,
                                Optional: true,
                        },
			"state_or_province": {
                                Type:    schema.TypeString,
                                Optional: true,
                        },
			"virtual": {
                                Type:    schema.TypeBool,
                                Optional: true,
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

// Create a new GTM Datacenter
func resourceGTMv1_3DatacenterCreate(d *schema.ResourceData, meta interface{}) error {

        domain := d.Get("domain").(string)

	log.Printf("[INFO] [Akamai GTM] Creating datacenter [%s] in domain [%s]", d.Get("nickname").(string), domain)
	newDC := populateNewDatacenterObject(d)
	log.Printf("[DEBUG] [Akamai GTMV1_3] Proposed New Datacenter: [%v]", newDC )
	cStatus, err := newDC.Create(domain)
        if err != nil {
		log.Printf("[DEBUG] [Akamai GTMV1_3] DC Create failed: %s", err.Error())
                fmt.Println(err)
                return err
        }
        b, err := json.Marshal(cStatus.Status)
        if err != nil {
                fmt.Println(err)
                return err
        }
        fmt.Println(string(b))
	log.Printf("[DEBUG] [Akamai GTMV1_3] Datacenter Create status:")
        log.Printf("[DEBUG] [Akamai GTMV1_3] %v", b)

        if d.Get("wait_on_complete").(bool) {
                done, err := waitForCompletion(domain)
                if done {
                        log.Printf("[INFO] [Akamai GTMV1_3] Datacenter Create completed")
                } else {
                        if err == nil {
                                log.Printf("[INFO] [Akamai GTMV1_3] Datacenter Create pending")
                        } else {
                                log.Printf("[WARNING] [Akamai GTMV1_3] Datacenter Create failed [%s]", err.Error())
                                return err
                        }
                }

        }

	// Give terraform the ID. Format domain::dcid
	datacenterId := fmt.Sprintf("%s:%d", domain, cStatus.Resource.DatacenterId)
	log.Printf("[DEBUG] [Akamai GTMV1_3] Generated DC Resource Id: %s", datacenterId)
	d.SetId(datacenterId)
	return resourceGTMv1_3DatacenterRead(d, meta)

}

// Only ever save data from the tf config in the tf state file, to help with
// api issues. See func unmarshalResourceData for more info.
func resourceGTMv1_3DatacenterRead(d *schema.ResourceData, meta interface{}) error {

	log.Printf("[DEBUG] [Akamai GTMV1_3] READ")
	log.Printf("[DEBUG] Reading [Akamai GTMV1_3] Datacenter: %s", d.Id())
	// retrieve the datacenter and domain
	domain, dcID, err := parseDatacenterResourceId(d.Id()) 
        if err != nil {
		return errors.New("Invalid datacenter resource Id")
	}
	dc, err := gtmv1_3.GetDatacenter(dcID, domain)
	if err != nil {
 		fmt.Println(err)
		log.Printf("[DEBUG] [Akamai GTMV1_3] Datacenter Read error: %s", err.Error())
		return err
	}
	populateTerraformDCState(d, dc)
	// Need set for Import. Good to confirm for read as well ...
	d.Set("domain", domain)
	log.Printf("[DEBUG] [Akamai GTMV1_3] READ %v", dc)
	return nil
}

// Update GTM Datacenter
func resourceGTMv1_3DatacenterUpdate(d *schema.ResourceData, meta interface{}) error {

	log.Printf("[DEBUG] [Akamai GTMV1_3] UPDATE")
        log.Printf("[DEBUG] Updating [Akamai GTMV1_3] Datacenter: %s", d.Id())
	// pull domain and dcid out of resource id
        domain, dcID, err := parseDatacenterResourceId(d.Id()) 
        if err != nil {
                return errors.New("Invalid datacenter resource Id")
        } 
  	// Get existing datacenter
	existDC, err := gtmv1_3.GetDatacenter(dcID, domain)
       	if err != nil {
                fmt.Println(err.Error())
                return err
        }
        log.Printf("[DEBUG] Updating [Akamai GTMV1_3] Datacenter BEFORE: %v", existDC)
	populateDatacenterObject(d, existDC)
        log.Printf("[DEBUG] Updating [Akamai GTMV1_3] Datacenter PROPOSED: %v", existDC)
	uStat, err := existDC.Update(domain)
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
        log.Printf("[DEBUG] [Akamai GTMV1_3] Datacenter Update  status:")
        log.Printf("[DEBUG] [Akamai GTMV1_3] %v", b)

        if d.Get("wait_on_complete").(bool) {
                done, err := waitForCompletion(domain)
                if done {
                        log.Printf("[INFO] [Akamai GTMV1_3] Datacenter update completed")
                } else {
                        if err == nil {
                                log.Printf("[INFO] [Akamai GTMV1_3] Datacenter update pending")
                        } else {
                                log.Printf("[WARNING] [Akamai GTMV1_3] Datacenter update failed [%s]", err.Error())
                                return err
                        }
                }

        }

	return resourceGTMv1_3DatacenterRead(d, meta)
}

func resourceGTMv1_3DatacenterImport(d *schema.ResourceData, meta interface{}) ([]*schema.ResourceData, error) {

        log.Printf("[DEBUG] [Akamai GTMV1_3] Import")

	err := resourceGTMv1_3DatacenterRead(d, meta)
	return []*schema.ResourceData{d}, err
	
}

// Delete GTM Datacenter.
func resourceGTMv1_3DatacenterDelete(d *schema.ResourceData, meta interface{}) error {

        domain := d.Get("domain").(string)
        log.Printf("[DEBUG] [Akamai GTMV1_3] DELETE") 
        log.Printf("[DEBUG] Deleting [Akamai GTMV1_3] Datacenter: %d", d.Get("datacenter_id").(int))
        // Get existing datacenter
        existDC, err := gtmv1_3.GetDatacenter(d.Get("datacenter_id").(int), domain)
        if err != nil {
                fmt.Println(err.Error())
                return err
        }
        log.Printf("[DEBUG] Deleting [Akamai GTMV1_3] Datacenter: %v", existDC)
        uStat, err := existDC.Delete(domain)
        if err != nil {
                fmt.Println(err)
                return err
        }
        b, err := json.Marshal(uStat)
        if err != nil {
                fmt.Println(err)
                return err
        }
        fmt.Println(string(b))
        log.Printf("[DEBUG] [Akamai GTMV1_3] Datacenter Delete status:")
        log.Printf("[DEBUG] [Akamai GTMV1_3] %v", b)

        if d.Get("wait_on_complete").(bool) {
                done, err := waitForCompletion(domain)
                if done {
                        log.Printf("[INFO] [Akamai GTMV1_3] Datacenter delete completed")
                } else {
                        if err == nil {
                                log.Printf("[INFO] [Akamai GTMV1_3] Datacenter delete pending")
                        } else {
                                log.Printf("[WARNING] [Akamai GTMV1_3] Datacenter delete failed [%s]", err.Error())
                                return err
                        }
                }

        }

	// if succcessful ....
	d.SetId("")
	return nil
}

// Test GTM Datacenter existance
func resourceGTMv1_3DatacenterExists(d *schema.ResourceData, meta interface{}) (bool, error) {

        log.Printf("[DEBUG] [Akamai GTMv1_3] Exists")
        // pull domain and dcid out of resource id
        domain, dcID, err := parseDatacenterResourceId(d.Id())
        if err != nil {
                return false, errors.New("Invalid datacenter resource Id")
        }
	log.Printf("[DEBUG] [Akamai GTMV1_3] Searching for existing datacenter [%d] in domain %s", dcID, domain)
        dc, err := gtmv1_3.GetDatacenter(dcID, domain)
	log.Printf("[DEBUG] [Akamai GTMV1_3] Searching for Existing datacenter result [%v]", domain)
	return dc != nil, err
}

// Create and populate a new datacenter object from resource data
func populateNewDatacenterObject(d *schema.ResourceData) *gtmv1_3.Datacenter {

	dcObj := gtmv1_3.NewDatacenter()
	dcObj.DefaultLoadObject = gtmv1_3.NewLoadObject()
	populateDatacenterObject(d, dcObj)

	return dcObj

}

// Populate existing datacenter object from resource data
func populateDatacenterObject(d *schema.ResourceData, dc *gtmv1_3.Datacenter) {

        if v, ok := d.GetOk("nickname"); ok { dc.Nickname = v.(string) }
	if v, ok := d.GetOk("city"); ok { dc.City = v.(string) }
	if v, ok := d.GetOk("clone_of"); ok { dc.CloneOf = v.(int) }
	if v, ok := d.GetOk("cloud_server_targeting"); ok { dc.CloudServerTargeting = v.(bool) }
	if v, ok := d.GetOk("continent"); ok { dc.Continent = v.(string) }
	if v, ok := d.GetOk("country"); ok { dc.Country = v.(string) }
	// pull apart Set
	if v, ok := d.GetOk("default_load_object"); ok {
		dlo := getSingleSchemaSetItem(v)
		if dlo != nil {
			if dc.DefaultLoadObject == nil {
				dc.DefaultLoadObject = gtmv1_3.NewLoadObject()
			}
			if dlo["load_object"] != nil {
				dc.DefaultLoadObject.LoadObject = dlo["load_object"].(string)
			}
        		dc.DefaultLoadObject.LoadObjectPort = dlo["load_object_port"].(int)
                        if dlo["load_servers"] != nil {
                                ls := make([]string, len(dlo["load_servers"].([]interface{})))
                                for i, sl := range dlo["load_servers"].([]interface{}) {
                                        ls[i] = sl.(string)
				dc.DefaultLoadObject.LoadServers = ls 
				}
			}
		}
	}
	if v, ok := d.GetOk("latitude"); ok { dc.Latitude = v.(float64) }
	if v, ok := d.GetOk("longitude"); ok { dc.Longitude = v.(float64) }
	if v, ok := d.GetOk("nickname"); ok { dc.Nickname = v.(string) }
	if v, ok := d.GetOk("ping_interval"); ok { dc.PingInterval = v.(int) }
	if v, ok := d.GetOk("ping_packet_size"); ok { dc.PingPacketSize = v.(int) }
	if v, ok := d.GetOk("datacenter_id"); ok { dc.DatacenterId = v.(int) }
	if v, ok := d.GetOk("score_penalty"); ok { dc.ScorePenalty = v.(int) }
	if v, ok := d.GetOk("servermonitor_liveness_count"); ok { dc.ServermonitorLivenessCount = v.(int) }
	if v, ok := d.GetOk("servermonitor_load_count"); ok { dc.ServermonitorLoadCount = v.(int) }
	if v, ok := d.GetOk("servermonitor_pool"); ok { dc.ServermonitorPool = v.(string) }
	if v, ok := d.GetOk("state_or_province"); ok { dc.SstateOrProvince = v.(string) }
	if v, ok := d.GetOk("virtual"); ok { dc.Virtual = v.(bool) }

	return

}

// Populate Terraform state from provided Datacenter object
func populateTerraformDCState(d *schema.ResourceData, dc *gtmv1_3.Datacenter) {

	// walk thru all state elements
	d.Set("nickname", dc.Nickname)
	d.Set("datacenter_id", dc.DatacenterId)
        d.Set("city", dc.City)
        d.Set("clone_of", dc.CloneOf)
        d.Set("cloud_server_targeting", dc.CloudServerTargeting)
        d.Set("continent", dc.Continent)
        d.Set("country", dc.Country)
	_, ok := d.GetOkExists("default_load_object")
	if ok {
		log.Printf("[DEBUG] [Akamai GTMv1_3] default_load_object Exists")
        	dloNew := make(map[string]interface{})
		if dc.DefaultLoadObject != nil {
      			dloNew["load_object"] =  dc.DefaultLoadObject.LoadObject
       			dloNew["load_object_port"] = dc.DefaultLoadObject.LoadObjectPort
        		dloNew["load_servers"] = dc.DefaultLoadObject.LoadServers
		}
		dloNewList := make([]interface{}, 1)
		dloNewList[0] = dloNew
		d.Set("default_load_object",dloNewList)
	} else {
		log.Printf("[WARNING] [Akamai GTMv1_3] default_load_object attribute NOT in Terraform State")
	}
        d.Set("latitude", dc.Latitude)
        d.Set("longitude", dc.Longitude)
        d.Set("ping_interval", dc.PingInterval)
        d.Set("ping_packet_size", dc.PingPacketSize)
        d.Set("score_penalty", dc.ScorePenalty)
        d.Set("servermonitor_liveness_count", dc.ServermonitorLivenessCount)
        d.Set("servermonitor_load_count", dc.ServermonitorLoadCount)
        d.Set("servermonitor_pool", dc.ServermonitorPool)
        d.Set("state_or_province", dc.SstateOrProvince)
        d.Set("virtual", dc.Virtual)

	return

}

