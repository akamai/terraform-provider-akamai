package akamai

import (
	"encoding/json"
	"fmt"
	"log"
	"errors"
	gtmv1_3 "github.com/akamai/AkamaiOPEN-edgegrid-golang/configgtm-v1_3"
	"github.com/hashicorp/terraform/helper/schema"
)

func resourceGTMv1_3ASmap() *schema.Resource {
	return &schema.Resource{
		Create: resourceGTMv1_3ASmapCreate,
		Read:   resourceGTMv1_3ASmapRead,
		Update: resourceGTMv1_3ASmapUpdate,
		Delete: resourceGTMv1_3ASmapDelete,
		Exists: resourceGTMv1_3ASmapExists,
		Importer: &schema.ResourceImporter{
			State: resourceGTMv1_3ASmapImport,
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
                                                        Required: true,
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
                                                "as_numbers": &schema.Schema{
                                                        Type:    schema.TypeList,
                                                        Elem:    &schema.Schema{Type: schema.TypeInt},
                                                        Required: true,
                                                },
                                        },
                                },
                        },

		},
	}
}

// utility func to parse Terraform property asMap id
func parseResourceASmapId(id string) (string, string, error) {

	return parseResourceStringId(id)

}

// Create a new GTM ASmap
func resourceGTMv1_3ASmapCreate(d *schema.ResourceData, meta interface{}) error {

        domain := d.Get("domain").(string)

	log.Printf("[INFO] [Akamai GTM] Creating asMap [%s] in domain [%s]", d.Get("name").(string), domain)
	newAS := populateNewASmapObject(d)
	log.Printf("[DEBUG] [Akamai GTMV1_3] Proposed New ASmap: [%v]", newAS )
	cStatus, err := newAS.Create(domain)
        if err != nil {
		log.Printf("[DEBUG] [Akamai GTMV1_3] ASmap Create failed: %s", err.Error())
                fmt.Println(err)
                return err
        }
        b, err := json.Marshal(cStatus.Status)
        if err != nil {
                fmt.Println(err)
                return err
        }
        fmt.Println(string(b))
	log.Printf("[DEBUG] [Akamai GTMV1_3] ASmap Create status:")
        log.Printf("[DEBUG] [Akamai GTMV1_3] %v", b)

        if d.Get("wait_on_complete").(bool) {
                done, err := waitForCompletion(domain)
                if done {
                        log.Printf("[INFO] [Akamai GTMV1_3] ASmap Create completed")
                } else {
                        if err == nil {
                                log.Printf("[INFO] [Akamai GTMV1_3] ASmap Create pending")
                        } else {
                                log.Printf("[WARNING] [Akamai GTMV1_3] ASmap Create failed [%s]", err.Error())
                                return err
                        }
                }

        }

	// Give terraform the ID. Format domain:asMap
	asMapId := fmt.Sprintf("%s:%s", domain, cStatus.Resource.Name)
	log.Printf("[DEBUG] [Akamai GTMV1_3] Generated ASmap ASmap Id: %s", asMapId)
	d.SetId(asMapId)
	return resourceGTMv1_3ASmapRead(d, meta)

}

// read asMap. updates state with entire API result configuration.
func resourceGTMv1_3ASmapRead(d *schema.ResourceData, meta interface{}) error {

	log.Printf("[DEBUG] [Akamai GTMV1_3] READ")
	log.Printf("[DEBUG] Reading [Akamai GTMV1_3] ASmap: %s", d.Id())
	// retrieve the property and domain
	domain, asMap, err := parseResourceASmapId(d.Id()) 
        if err != nil {
		return errors.New("Invalid asMap asMap Id")
	}
	as, err := gtmv1_3.GetAsMap(asMap, domain)
	if err != nil {
 		fmt.Println(err)
		log.Printf("[DEBUG] [Akamai GTMV1_3] ASmap Read error: %s", err.Error())
		return err
	}
	populateTerraformASmapState(d, as)
	log.Printf("[DEBUG] [Akamai GTMV1_3] READ %v", as)
	return nil
}

// Update GTM ASmap
func resourceGTMv1_3ASmapUpdate(d *schema.ResourceData, meta interface{}) error {

	log.Printf("[DEBUG] [Akamai GTMV1_3] UPDATE")
        log.Printf("[DEBUG] Updating [Akamai GTMV1_3] ASmap: %s", d.Id())
	// pull domain and asMap out of id
        domain, asMap, err := parseResourceASmapId(d.Id()) 
        if err != nil {
                return errors.New("Invalid asMap Id")
        } 
  	// Get existingASmap 
	existAs, err := gtmv1_3.GetAsMap(asMap, domain)
       	if err != nil {
                fmt.Println(err.Error())
                return err
        }
        log.Printf("[DEBUG] Updating [Akamai GTMV1_3] ASmap BEFORE: %v", existAs)
	populateASmapObject(d, existAs)
        log.Printf("[DEBUG] Updating [Akamai GTMV1_3] ASmap PROPOSED: %v", existAs)
	uStat, err := existAs.Update(domain)
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
        log.Printf("[DEBUG] [Akamai GTMV1_3] ASmap Update  status:")
        log.Printf("[DEBUG] [Akamai GTMV1_3] %v", b)

        if d.Get("wait_on_complete").(bool) {
                done, err := waitForCompletion(domain)
                if done {
                        log.Printf("[INFO] [Akamai GTMV1_3] ASmap update completed")
                } else {
                        if err == nil {
                                log.Printf("[INFO] [Akamai GTMV1_3] ASmap update pending")
                        } else {
                                log.Printf("[WARNING] [Akamai GTMV1_3] ASmap update failed [%s]", err.Error())
                                return err
                        }
                }

        }

	return resourceGTMv1_3ASmapRead(d, meta)
}

// Import GTM ASmap.
func resourceGTMv1_3ASmapImport(d *schema.ResourceData, meta interface{}) ([]*schema.ResourceData, error) {

        log.Printf("[INFO] [Akamai GTM] ASmap [%s] Import", d.Id())
        // pull domain and asMap out of asMap id
        domain, asMap, err := parseResourceASmapId(d.Id())
        if err != nil {
                return []*schema.ResourceData{d}, errors.New("Invalid asMap Id")
        }
	as, err := gtmv1_3.GetAsMap(asMap, domain)
	if err != nil {
		return nil, err
	}
	d.Set("domain", domain)
	d.Set("wait_on_complete", true)
	populateTerraformASmapState(d, as)

	// use same Id as passed in
	log.Printf("[INFO] [Akamai GTM] ASmap [%s] [%s] Imported", d.Id(), d.Get("name"))
	return []*schema.ResourceData{d}, nil
}


// Delete GTM ASmap.
func resourceGTMv1_3ASmapDelete(d *schema.ResourceData, meta interface{}) error {

        domain := d.Get("domain").(string)
        log.Printf("[DEBUG] [Akamai GTMV1_3] DELETE") 
        log.Printf("[DEBUG] Deleting [Akamai GTMV1_3] ASmap: %s", d.Id())
        // Get existing asMap
        domain, asMap, err := parseResourceASmapId(d.Id())
        if err != nil {
                return errors.New("Invalid asMap Id")
        }  
        existAs, err := gtmv1_3.GetAsMap(asMap, domain)
        if err != nil {
                fmt.Println(err.Error())
                return err
        }
        log.Printf("[DEBUG] Deleting [Akamai GTMV1_3] ASmap: %v", existAs)
        uStat, err := existAs.Delete(domain)
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
        log.Printf("[DEBUG] [Akamai GTMV1_3] ASmap Delete status:")
        log.Printf("[DEBUG] [Akamai GTMV1_3] %v", b)

        if d.Get("wait_on_complete").(bool) {
                done, err := waitForCompletion(domain)
                if done {
                        log.Printf("[INFO] [Akamai GTMV1_3] ASmap delete completed")
                } else {
                        if err == nil {
                                log.Printf("[INFO] [Akamai GTMV1_3] ASmap delete pending")
                        } else {
                                log.Printf("[WARNING] [Akamai GTMV1_3] ASmap delete failed [%s]", err.Error())
                                return err
                        }
                }

        }

	// if succcessful ....
	d.SetId("")
	return nil
}

// Test GTM ASmap existance
func resourceGTMv1_3ASmapExists(d *schema.ResourceData, meta interface{}) (bool, error) {

        log.Printf("[DEBUG] [Akamai GTMv1_3] Exists")
        // pull domain and asMap out of asMap id
        domain, asMap, err := parseResourceASmapId(d.Id())
        if err != nil {
                return false, errors.New("Invalid asMap asMap Id")
        }
	log.Printf("[DEBUG] [Akamai GTMV1_3] Searching for existing asMap [%s] in domain %s", asMap, domain)
        as, err := gtmv1_3.GetAsMap(asMap, domain)
	return as != nil, err
}

// Create and populate a new asMap object from asMap data
func populateNewASmapObject(d *schema.ResourceData) *gtmv1_3.AsMap {

	asObj := gtmv1_3.NewAsMap(d.Get("name").(string))
	asObj.DefaultDatacenter = &gtmv1_3.DatacenterBase{}
        asObj.Assignments = make([]*gtmv1_3.AsAssignment, 1)
        asObj.Links = make([]*gtmv1_3.Link, 1)
	populateASmapObject(d, asObj)

	return asObj

}

// Populate existing asMap object from asMap data
func populateASmapObject(d *schema.ResourceData, as *gtmv1_3.AsMap) {

        if v, ok := d.GetOk("name"); ok { as.Name = v.(string) }
	populateAsAssignmentsObject(d, as)
        populateAsDefaultDCObject(d, as)

	return

}

// Populate Terraform state from provided ASmap object
func populateTerraformASmapState(d *schema.ResourceData, as *gtmv1_3.AsMap) {

	// walk thru all state elements
	d.Set("name", as.Name)
	populateTerraformAsAssignmentsState(d, as)
        populateTerraformAsDefaultDCState(d, as)

	return

}

// create and populate GTM ASmap Assignments object
func populateAsAssignmentsObject(d *schema.ResourceData, as *gtmv1_3.AsMap) {

        // pull apart List
	assgn := d.Get("assignments")
	if assgn != nil {
		asAssignmentsList := assgn.([]interface{})
		asAssignmentsObjList := make([]*gtmv1_3.AsAssignment, len(asAssignmentsList)) // create new object list
		for i, v := range asAssignmentsList {
			asMap := v.(map[string]interface{})
			asAssignment := gtmv1_3.AsAssignment{}
			asAssignment.DatacenterId = asMap["datacenter_id"].(int)
                        asAssignment.Nickname = asMap["nickname"].(string)
                	if asMap["as_numbers"] != nil {
				ls := make([]int64, len(asMap["as_numbers"].([]interface{})))
                		for i, sl := range asMap["as_numbers"].([]interface{}) { 
					ls[i] = sl.(int64)
				}
                        	asAssignment.AsNumbers = ls
                	}
			asAssignmentsObjList[i] = &asAssignment
		}
		as.Assignments = asAssignmentsObjList
	}
}

// create and populate Terraform asMap assigments schema 
func populateTerraformAsAssignmentsState(d *schema.ResourceData, as *gtmv1_3.AsMap) {

	asListNew := make([]interface{}, len(as.Assignments))
	for i, assgn := range as.Assignments {
		asNew := map[string]interface{}{
					"datacenter_id":		assgn.DatacenterId,
					"nickname":			assgn.Nickname,
					"as_numbers":			assgn.AsNumbers,
			}
		asListNew[i] = asNew
	}
        d.Set("assignments", asListNew)

}

// create and populate GTM ASmap DefaultDatacenter object
func populateAsDefaultDCObject(d *schema.ResourceData, as *gtmv1_3.AsMap) {

        // pull apart List
        asm := d.Get("default_datacenter")
        if asm != nil && len(asm.([]interface{})) > 0 {
                asDefaultDCObj := gtmv1_3.DatacenterBase{} // create new object
		asDefaultDCList := asm.([]interface{})
                asMap := asDefaultDCList[0].(map[string]interface{})
		if asMap["datacenter_id"] != nil && asMap["datacenter_id"].(int) != 0 {
                	asDefaultDCObj.DatacenterId = asMap["datacenter_id"].(int)
                	asDefaultDCObj.Nickname = asMap["nickname"].(string)
		} else {
			log.Printf("[INFO] [Akamai GTMv1_3] No Default Datacenter specified")
			var nilInt int
			asDefaultDCObj.DatacenterId = nilInt
			asDefaultDCObj.Nickname = ""
		}
		as.DefaultDatacenter = &asDefaultDCObj
        }
}

// create and populate Terraform asMap default_datacenter schema
func populateTerraformAsDefaultDCState(d *schema.ResourceData, as *gtmv1_3.AsMap) {

        ddcListNew := make([]interface{}, 1)
        ddcNew := map[string]interface{}{
                                "datacenter_id":     as.DefaultDatacenter.DatacenterId,
                        	"nickname":          as.DefaultDatacenter.Nickname,
			}
        ddcListNew[0] = ddcNew
        d.Set("default_datacenter", ddcListNew)

}

