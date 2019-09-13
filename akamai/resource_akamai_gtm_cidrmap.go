package akamai

import (
	"encoding/json"
	"fmt"
	"log"
	"errors"
	gtmv1_3 "github.com/akamai/AkamaiOPEN-edgegrid-golang/configgtm-v1_3"
	"github.com/hashicorp/terraform/helper/schema"
)

func resourceGTMv1_3Cidrmap() *schema.Resource {
	return &schema.Resource{
		Create: resourceGTMv1_3CidrMapCreate,
		Read:   resourceGTMv1_3CidrMapRead,
		Update: resourceGTMv1_3CidrMapUpdate,
		Delete: resourceGTMv1_3CidrMapDelete,
		Exists: resourceGTMv1_3CidrMapExists,
		Importer: &schema.ResourceImporter{
			State: resourceGTMv1_3CidrMapImport,
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
                        "default_datacenter": &schema.Schema{
                                Type:    schema.TypeList,
                                Required: true,
				MaxItems: 1,
                                ConfigMode: schema.SchemaConfigModeAttr,
                                Elem:    &schema.Resource{
                                        Schema:map[string]*schema.Schema{
                                                "datacenter_id": {
                                                        Type:     schema.TypeInt,
                                                        Optional: true,
							Default: nil,
                                                },
                                                "nickname": {
                                                        Type:         schema.TypeString,
                                                        Optional:     true,
                                                },
                                        },
                                },
                        },
                        "assignments": &schema.Schema{
                                Type:    schema.TypeList,
                                Optional: true,
				ConfigMode: schema.SchemaConfigModeAttr,
                                Elem:    &schema.Resource{
                                        Schema:map[string]*schema.Schema{
			                        "datacenter_id": {
							Type:     schema.TypeInt,
                                			Required: true,
                        			},
                        			"nickname": {
                                			Type:         schema.TypeString,
                                			Required:     true,
                        			},
                                                "blocks": &schema.Schema{
                                                        Type:    schema.TypeList,
                                                        Elem:    &schema.Schema{Type: schema.TypeString},
                                                        Optional: true,
                                                },
                                        },
                                },
                        },

		},
	}
}

// utility func to parse Terraform property cidrMap id
func parseResourceCidrMapId(id string) (string, string, error) {

	return parseResourceStringId(id)

}

// Create a new GTM CidrMap
func resourceGTMv1_3CidrMapCreate(d *schema.ResourceData, meta interface{}) error {

        domain := d.Get("domain").(string)

	log.Printf("[INFO] [Akamai GTM] Creating cidrMap [%s] in domain [%s]", d.Get("name").(string), domain)
	newRsrc := populateNewCidrMapObject(d)
	log.Printf("[DEBUG] [Akamai GTMV1_3] Proposed New CidrMap: [%v]", newRsrc )
	cStatus, err := newRsrc.Create(domain)
        if err != nil {
		log.Printf("[DEBUG] [Akamai GTMV1_3] CidrMap Create failed: %s", err.Error())
                fmt.Println(err)
                return err
        }
        b, err := json.Marshal(cStatus.Status)
        if err != nil {
                fmt.Println(err)
                return err
        }
        fmt.Println(string(b))
	log.Printf("[DEBUG] [Akamai GTMV1_3] CidrMap Create status:")
        log.Printf("[DEBUG] [Akamai GTMV1_3] %v", b)

        if d.Get("wait_on_complete").(bool) {
                done, err := waitForCompletion(domain)
                if done {
                        log.Printf("[INFO] [Akamai GTMV1_3] CidrMap Create completed")
                } else {
                        if err == nil {
                                log.Printf("[INFO] [Akamai GTMV1_3] CidrMap Create pending")
                        } else {
                                log.Printf("[WARNING] [Akamai GTMV1_3] CidrMap Create failed [%s]", err.Error())
                                return err
                        }
                }

        }

	// Give terraform the ID. Format domain:cidrMap
	cidrMapId := fmt.Sprintf("%s:%s", domain, cStatus.Resource.Name)
	log.Printf("[DEBUG] [Akamai GTMV1_3] Generated CidrMap CidrMap Id: %s", cidrMapId)
	d.SetId(cidrMapId)
	return resourceGTMv1_3CidrMapRead(d, meta)

}

// read cidrMap. updates state with entire API result configuration.
func resourceGTMv1_3CidrMapRead(d *schema.ResourceData, meta interface{}) error {

	log.Printf("[DEBUG] [Akamai GTMV1_3] READ")
	log.Printf("[DEBUG] Reading [Akamai GTMV1_3] CidrMap: %s", d.Id())
	// retrieve the property and domain
	domain, cidrMap, err := parseResourceCidrMapId(d.Id()) 
        if err != nil {
		return errors.New("Invalid cidrMap cidrMap Id")
	}
	cidr, err := gtmv1_3.GetCidrMap(cidrMap, domain)
	if err != nil {
 		fmt.Println(err)
		log.Printf("[DEBUG] [Akamai GTMV1_3] CidrMap Read error: %s", err.Error())
		return err
	}
	populateTerraformCidrMapState(d, cidr)
	log.Printf("[DEBUG] [Akamai GTMV1_3] READ %v", cidr)
	return nil
}

// Update GTM CidrMap
func resourceGTMv1_3CidrMapUpdate(d *schema.ResourceData, meta interface{}) error {

	log.Printf("[DEBUG] [Akamai GTMV1_3] UPDATE")
        log.Printf("[DEBUG] Updating [Akamai GTMV1_3] CidrMap: %s", d.Id())
	// pull domain and cidrMap out of id
        domain, cidrMap, err := parseResourceCidrMapId(d.Id()) 
        if err != nil {
                return errors.New("Invalid cidrMap Id")
        } 
  	// Get existingCidrMap 
	existCidr, err := gtmv1_3.GetCidrMap(cidrMap, domain)
       	if err != nil {
                fmt.Println(err.Error())
                return err
        }
        log.Printf("[DEBUG] Updating [Akamai GTMV1_3] CidrMap BEFORE: %v", existCidr)
	populateCidrMapObject(d, existCidr)
        log.Printf("[DEBUG] Updating [Akamai GTMV1_3] CidrMap PROPOSED: %v", existCidr)
	uStat, err := existCidr.Update(domain)
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
        log.Printf("[DEBUG] [Akamai GTMV1_3] CidrMap Update  status:")
        log.Printf("[DEBUG] [Akamai GTMV1_3] %v", b)

        if d.Get("wait_on_complete").(bool) {
                done, err := waitForCompletion(domain)
                if done {
                        log.Printf("[INFO] [Akamai GTMV1_3] CidrMap update completed")
                } else {
                        if err == nil {
                                log.Printf("[INFO] [Akamai GTMV1_3] CidrMap update pending")
                        } else {
                                log.Printf("[WARNING] [Akamai GTMV1_3] CidrMap update failed [%s]", err.Error())
                                return err
                        }
                }

        }

	return resourceGTMv1_3CidrMapRead(d, meta)
}

// Import GTM CidrMap.
func resourceGTMv1_3CidrMapImport(d *schema.ResourceData, meta interface{}) ([]*schema.ResourceData, error) {

        log.Printf("[INFO] [Akamai GTM] CidrMap [%s] Import", d.Id())
        // pull domain and cidrMap out of cidrMap id
        domain, cidrMap, err := parseResourceCidrMapId(d.Id())
        if err != nil {
                return []*schema.ResourceData{d}, errors.New("Invalid cidrMap cidrMap Id")
        }
	cidr, err := gtmv1_3.GetCidrMap(cidrMap, domain)
	if err != nil {
		return nil, err
	}
	d.Set("domain", domain)
	d.Set("wait_on_complete", true)
	populateTerraformCidrMapState(d, cidr)

	// use same Id as passed in
	log.Printf("[INFO] [Akamai GTM] CidrMap [%s] [%s] Imported", d.Id(), d.Get("name"))
	return []*schema.ResourceData{d}, nil
}


// Delete GTM CidrMap.
func resourceGTMv1_3CidrMapDelete(d *schema.ResourceData, meta interface{}) error {

        domain := d.Get("domain").(string)
        log.Printf("[DEBUG] [Akamai GTMV1_3] DELETE") 
        log.Printf("[DEBUG] Deleting [Akamai GTMV1_3] CidrMap: %s", d.Id())
        // Get existing cidrMap
        domain, cidrMap, err := parseResourceCidrMapId(d.Id())
        if err != nil {
                return errors.New("Invalid cidrMap Id")
        }  
        existCidr, err := gtmv1_3.GetCidrMap(cidrMap, domain)
        if err != nil {
                fmt.Println(err.Error())
                return err
        }
        log.Printf("[DEBUG] Deleting [Akamai GTMV1_3] CidrMap: %v", existCidr)
        uStat, err := existCidr.Delete(domain)
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
        log.Printf("[DEBUG] [Akamai GTMV1_3] CidrMap Delete status:")
        log.Printf("[DEBUG] [Akamai GTMV1_3] %v", b)

        if d.Get("wait_on_complete").(bool) {
                done, err := waitForCompletion(domain)
                if done {
                        log.Printf("[INFO] [Akamai GTMV1_3] CidrMap delete completed")
                } else {
                        if err == nil {
                                log.Printf("[INFO] [Akamai GTMV1_3] CidrMap delete pending")
                        } else {
                                log.Printf("[WARNING] [Akamai GTMV1_3] CidrMap delete failed [%s]", err.Error())
                                return err
                        }
                }

        }

	// if succcessful ....
	d.SetId("")
	return nil
}

// Test GTM CidrMap existance
func resourceGTMv1_3CidrMapExists(d *schema.ResourceData, meta interface{}) (bool, error) {

        log.Printf("[DEBUG] [Akamai GTMv1_3] Exists")
        // pull domain and cidrMap out of cidrMap id
        domain, cidrMap, err := parseResourceCidrMapId(d.Id())
        if err != nil {
                return false, errors.New("Invalid cidrMap cidrMap Id")
        }
	log.Printf("[DEBUG] [Akamai GTMV1_3] Searching for existing cidrMap [%s] in domain %s", cidrMap, domain)
        cidr, err := gtmv1_3.GetCidrMap(cidrMap, domain)
	return cidr != nil, err
}

// Create and populate a new cidrMap object from cidrMap data
func populateNewCidrMapObject(d *schema.ResourceData) *gtmv1_3.CidrMap {

	cidrObj := gtmv1_3.NewCidrMap(d.Get("name").(string))
	cidrObj.DefaultDatacenter = &gtmv1_3.DatacenterBase{}
        cidrObj.Assignments = make([]*gtmv1_3.CidrAssignment, 1)
	populateCidrMapObject(d, cidrObj)

	return cidrObj

}

// Populate existing cidrMap object from cidrMap data
func populateCidrMapObject(d *schema.ResourceData, cidr *gtmv1_3.CidrMap) {

        if v, ok := d.GetOk("name"); ok { cidr.Name = v.(string) }
	populateCidrAssignmentsObject(d, cidr)
        populateCidrDefaultDCObject(d, cidr)

	return

}

// Populate Terraform state from provided CidrMap object
func populateTerraformCidrMapState(d *schema.ResourceData, cidr *gtmv1_3.CidrMap) {

	// walk thru all state elements
	d.Set("name", cidr.Name)
	populateTerraformCidrAssignmentsState(d, cidr)
        populateTerraformCidrDefaultDCState(d, cidr)

	return

}

// create and populate GTM CidrMap Assignments object
func populateCidrAssignmentsObject(d *schema.ResourceData, cidr *gtmv1_3.CidrMap) {

        // pull apart List
	as := d.Get("assignments")
	if as != nil {
		cidrAssignmentsList := as.([]interface{})
		cidrAssignmentsObjList := make([]*gtmv1_3.CidrAssignment, len(cidrAssignmentsList)) // create new object list
		for i, v := range cidrAssignmentsList {
			asMap := v.(map[string]interface{})
			cidrAssignment := gtmv1_3.CidrAssignment{}
			cidrAssignment.DatacenterId = asMap["datacenter_id"].(int)
                        cidrAssignment.Nickname = asMap["nickname"].(string)
                	if asMap["blocks"] != nil {
				ls := make([]string, len(asMap["blocks"].([]interface{})))
                		for i, sl := range asMap["blocks"].([]interface{}) { 
					ls[i] = sl.(string)
				}
                        	cidrAssignment.Blocks = ls
                	}
			cidrAssignmentsObjList[i] = &cidrAssignment
		}
		cidr.Assignments = cidrAssignmentsObjList
	}
}

// create and populate Terraform cidrMap assigments schema 
func populateTerraformCidrAssignmentsState(d *schema.ResourceData, cidr *gtmv1_3.CidrMap) {

	asListNew := make([]interface{}, len(cidr.Assignments))
	for i, as := range cidr.Assignments {
		asNew := map[string]interface{}{
					"datacenter_id":		as.DatacenterId,
					"nickname":			as.Nickname,
					"blocks":			as.Blocks,
			}
		asListNew[i] = asNew
	}
        d.Set("assignments", asListNew)

}

// create and populate GTM CidrMap DefaultDatacenter object
func populateCidrDefaultDCObject(d *schema.ResourceData, cidr *gtmv1_3.CidrMap) {

        // pull apart List
        as := d.Get("default_datacenter")
        if as != nil && len(as.([]interface{})) > 0 {
                cidrDefaultDCObj := gtmv1_3.DatacenterBase{} // create new object
		cidrDefaultDCList := as.([]interface{})
                asMap := cidrDefaultDCList[0].(map[string]interface{})
		if asMap["datacenter_id"] != nil && asMap["datacenter_id"].(int) != 0 {
			log.Printf("[DEBUG] [Akamai GTMv1_3] Default Datacenter: %v ... %d ", asMap["datacenter_id"], asMap["datacenter_id"].(int))
                	cidrDefaultDCObj.DatacenterId = asMap["datacenter_id"].(int)
                	cidrDefaultDCObj.Nickname = asMap["nickname"].(string)
		} else {
			log.Printf("[INFO] [Akamai GTMv1_3] No Default Datacenter specified")
			var nilInt int
			cidrDefaultDCObj.DatacenterId = nilInt
			cidrDefaultDCObj.Nickname = ""
		}
		cidr.DefaultDatacenter = &cidrDefaultDCObj
        }
}

// create and populate Terraform cidrMap default_datacenter schema
func populateTerraformCidrDefaultDCState(d *schema.ResourceData, cidr *gtmv1_3.CidrMap) {

        ddcListNew := make([]interface{}, 1)
        ddcNew := map[string]interface{}{
                                "datacenter_id":     cidr.DefaultDatacenter.DatacenterId,
                        	"nickname":          cidr.DefaultDatacenter.Nickname,
			}
        ddcListNew[0] = ddcNew
        d.Set("default_datacenter", ddcListNew)

}

