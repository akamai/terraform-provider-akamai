package gtm

import (
	"fmt"
	"log"

	gtm "github.com/akamai/AkamaiOPEN-edgegrid-golang/configgtm-v1_4"
	"github.com/akamai/terraform-provider-akamai/v2/pkg/akamai"
	"github.com/akamai/terraform-provider-akamai/v2/pkg/tools"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/schema"
)

func resourceGTMv1ASmap() *schema.Resource {
	return &schema.Resource{
		Create: resourceGTMv1ASmapCreate,
		Read:   resourceGTMv1ASmapRead,
		Update: resourceGTMv1ASmapUpdate,
		Delete: resourceGTMv1ASmapDelete,
		Exists: resourceGTMv1ASmapExists,
		Importer: &schema.ResourceImporter{
			State: resourceGTMv1ASmapImport,
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
			"default_datacenter": {
				Type:       schema.TypeList,
				Required:   true,
				MaxItems:   1,
				ConfigMode: schema.SchemaConfigModeAttr,
				Elem: &schema.Resource{
					Schema: map[string]*schema.Schema{
						"datacenter_id": {
							Type:     schema.TypeInt,
							Required: true,
						},
						"nickname": {
							Type:     schema.TypeString,
							Optional: true,
						},
					},
				},
			},
			"assignment": {
				Type:     schema.TypeList,
				Optional: true,
				Elem: &schema.Resource{
					Schema: map[string]*schema.Schema{
						"datacenter_id": {
							Type:     schema.TypeInt,
							Required: true,
						},
						"nickname": {
							Type:     schema.TypeString,
							Required: true,
						},
						"as_numbers": {
							Type:     schema.TypeList,
							Elem:     &schema.Schema{Type: schema.TypeInt},
							Required: true,
						},
					},
				},
			},
		},
	}
}

// Util method to validate default datacenter and create if necessary
func validateDefaultDC(ddcField []interface{}, domain string) error {

	if len(ddcField) == 0 {
		return fmt.Errorf("default Datacenter invalid")
	}
	ddc, ok := ddcField[0].(map[string]interface{})
	if !ok {
		return fmt.Errorf("invalid default_datacenter configuration")
	}

	intrDcID, ok := ddc["datacenter_id"]
	if !ok {
		return fmt.Errorf("default Datacenter ID invalid")
	}

	dcId, ok := intrDcID.(int)
	if !ok || dcId == 0 {
		return fmt.Errorf("default Datacenter ID invalid")
	}
	dc, err := gtm.GetDatacenter(dcId, domain)
	if dc == nil {
		if err != nil {
			if _, ok := err.(gtm.CommonError); !ok {
				return fmt.Errorf("MapCreate Unexpected error verifying Default Datacenter exists: %s", err.Error())
			}
		}
		// ddc doesn't exist
		if ddc["datacenter_id"].(int) != gtm.MapDefaultDC {
			return fmt.Errorf(fmt.Sprintf("Default Datacenter %d does not exist", ddc["datacenter_id"].(int)))
		}
		_, err := gtm.CreateMapsDefaultDatacenter(domain) // create if not already.
		if err != nil {
			return fmt.Errorf("MapCreate failed on Default Datacenter check: %s", err.Error())
		}
	}

	return nil
}

// Create a new GTM ASmap
func resourceGTMv1ASmapCreate(d *schema.ResourceData, m interface{}) error {
	meta := akamai.Meta(m)
	logger := meta.Log("Akamai GTMv1", "resourceGTMv1ASmapCreate")

	domain, err := tools.GetStringValue("domain", d)
	if err != nil {
		logger.Errorf("Domain not initialized: %s", err.Error())
		return err
	}

	if domainName, err := tools.GetStringValue("name", d); err != nil {
		logger.Warnf("AsMap not initialized: %s", err.Error())
		return err
	} else {
		logger.Infof("Creating asMap [%s] in domain [%s]", domainName, domain)
	}

	// Make sure Default Datacenter exists
	interfaceArray, err := tools.GetInterfaceArrayValue("default_datacenter", d)
	if err != nil {
		return err
	}
	if err = validateDefaultDC(interfaceArray, domain); err != nil {
		logger.Errorf("Default datacenter validation error: %s", err.Error())
		return err
	}

	newAS := populateNewASmapObject(d, m)
	logger.Debugf("Proposed New ASmap: [%v]", newAS)
	cStatus, err := newAS.Create(domain)
	if err != nil {
		logger.Errorf("ASmap Create failed: %s", err.Error())
		return err
	}
	logger.Debugf("ASmap Create status: %v", cStatus.Status)
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
			logger.Infof("ASmap Create completed")
		} else {
			if err == nil {
				logger.Infof("ASmap Create pending")
			} else {
				logger.Errorf("ASmap Create failed [%s]", err.Error())
				return err
			}
		}
	}

	// Give terraform the ID. Format domain:asMap
	asMapID := fmt.Sprintf("%s:%s", domain, cStatus.Resource.Name)
	logger.Debugf("Generated ASmap ASmap Id: %s", asMapID)
	d.SetId(asMapID)
	return resourceGTMv1ASmapRead(d, m)

}

// read asMap. updates state with entire API result configuration.
func resourceGTMv1ASmapRead(d *schema.ResourceData, m interface{}) error {
	meta := akamai.Meta(m)
	logger := meta.Log("Akamai GTMv1", "resourceGTMv1ASmapRead")

	logger.Debugf("Reading ASmap: %s", d.Id())
	// retrieve the property and domain
	domain, asMap, err := parseResourceStringId(d.Id())
	if err != nil {
		return fmt.Errorf("invalid asMap ID")
	}
	as, err := gtm.GetAsMap(asMap, domain)
	if err != nil {
		logger.Errorf("ASmap Read error: %s", err.Error())
		return err
	}
	populateTerraformASmapState(d, as, m)
	logger.Debugf("READ %v", as)
	return nil
}

// Update GTM ASmap
func resourceGTMv1ASmapUpdate(d *schema.ResourceData, m interface{}) error {
	meta := akamai.Meta(m)
	logger := meta.Log("Akamai GTMv1", "resourceGTMv1ASmapUpdate")

	logger.Debugf("UPDATE ASmap: %s", d.Id())
	// pull domain and asMap out of id
	domain, asMap, err := parseResourceStringId(d.Id())
	if err != nil {
		return fmt.Errorf("invalid asMap ID")
	}
	// Get existingASmap
	existAs, err := gtm.GetAsMap(asMap, domain)
	if err != nil {
		logger.Errorf("ASmapUpdate: %s", err.Error())
		return err
	}
	logger.Debugf("ASmap BEFORE: %v", existAs)
	populateASmapObject(d, existAs, m)
	logger.Debugf("ASmap PROPOSED: %v", existAs)
	uStat, err := existAs.Update(domain)
	if err != nil {
		logger.Errorf("ASmapUpdate: %s", err.Error())
		return err
	}
	logger.Debugf("ASmap Update  status: %v", uStat)
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
			logger.Infof("ASmap update completed")
		} else {
			if err == nil {
				logger.Infof("ASmap update pending")
			} else {
				logger.Errorf("ASmap update failed [%s]", err.Error())
				return err
			}
		}
	}

	return resourceGTMv1ASmapRead(d, m)
}

// Import GTM ASmap.
func resourceGTMv1ASmapImport(d *schema.ResourceData, m interface{}) ([]*schema.ResourceData, error) {
	meta := akamai.Meta(m)
	logger := meta.Log("Akamai GTMv1", "resourceGTMv1ASmapImport")

	logger.Infof("ASmap [%s] Import", d.Id())
	// pull domain and asMap out of asMap id
	domain, asMap, err := parseResourceStringId(d.Id())
	if err != nil {
		return []*schema.ResourceData{d}, fmt.Errorf("invalid asMap ID")
	}
	as, err := gtm.GetAsMap(asMap, domain)
	if err != nil {
		return nil, err
	}
	if err := d.Set("domain", domain); err != nil {
		return nil, err
	}
	if err := d.Set("wait_on_complete", true); err != nil {
		return nil, err
	}
	populateTerraformASmapState(d, as, m)

	// use same Id as passed in
	logger.Infof("ASmap [%s] [%s] Imported", d.Id(), d.Get("name"))
	return []*schema.ResourceData{d}, nil
}

// Delete GTM ASmap.
func resourceGTMv1ASmapDelete(d *schema.ResourceData, m interface{}) error {
	meta := akamai.Meta(m)
	logger := meta.Log("Akamai GTMv1", "resourceGTMv1ASmapDelete")

	logger.Debugf("DELETE")
	logger.Debugf("Deleting ASmap: %s", d.Id())
	// Get existing asMap
	domain, asMap, err := parseResourceStringId(d.Id())
	if err != nil {
		log.Printf("[ERROR] ASmapDelete: %s", err.Error())
		return fmt.Errorf("invalid asMap ID")
	}
	existAs, err := gtm.GetAsMap(asMap, domain)
	if err != nil {
		logger.Errorf("ASmapDelete: %s", err.Error())
		return err
	}
	logger.Debugf("Deleting ASmap: %v", existAs)
	uStat, err := existAs.Delete(domain)
	if err != nil {
		logger.Errorf("ASmapDelete: %s", err.Error())
		return err
	}
	logger.Debugf("ASmap Delete status:")
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
			logger.Infof("ASmap delete completed")
		} else {
			if err == nil {
				logger.Infof("ASmap delete pending")
			} else {
				logger.Errorf("ASmap delete failed [%s]", err.Error())
				return err
			}
		}
	}

	// if successful ....
	d.SetId("")
	return nil
}

// Test GTM ASmap existence
func resourceGTMv1ASmapExists(d *schema.ResourceData, m interface{}) (bool, error) {
	meta := akamai.Meta(m)
	logger := meta.Log("Akamai GTMv1", "resourceGTMv1ASmapExists")

	logger.Debugf("Exists")
	// pull domain and asMap out of asMap id
	domain, asMap, err := parseResourceStringId(d.Id())
	if err != nil {
		return false, fmt.Errorf("invalid asMap ID")
	}
	logger.Debugf("Searching for existing asMap [%s] in domain %s", asMap, domain)
	as, err := gtm.GetAsMap(asMap, domain)
	return as != nil, err
}

// Create and populate a new asMap object from asMap data
func populateNewASmapObject(d *schema.ResourceData, m interface{}) *gtm.AsMap {

	asMapName, _ := tools.GetStringValue("name", d)
	asObj := gtm.NewAsMap(asMapName)
	asObj.DefaultDatacenter = &gtm.DatacenterBase{}
	asObj.Assignments = make([]*gtm.AsAssignment, 1)
	asObj.Links = make([]*gtm.Link, 1)
	populateASmapObject(d, asObj, m)

	return asObj

}

// Populate existing asMap object from asMap data
func populateASmapObject(d *schema.ResourceData, as *gtm.AsMap, m interface{}) {

	if v, err := tools.GetStringValue("name", d); err == nil {
		as.Name = v
	}
	populateAsAssignmentsObject(d, as, m)
	populateAsDefaultDCObject(d, as, m)

}

// Populate Terraform state from provided ASmap object
func populateTerraformASmapState(d *schema.ResourceData, as *gtm.AsMap, m interface{}) {
	meta := akamai.Meta(m)
	logger := meta.Log("Akamai GTMv1", "populateTerraformASmapState")

	// walk through all state elements
	if err := d.Set("name", as.Name); err != nil {
		logger.Errorf("populateTerraformASmapState failed: %s", err.Error())
	}
	populateTerraformAsAssignmentsState(d, as, m)
	populateTerraformAsDefaultDCState(d, as, m)
}

// create and populate GTM ASmap Assignments object
func populateAsAssignmentsObject(d *schema.ResourceData, as *gtm.AsMap, m interface{}) {
	meta := akamai.Meta(m)
	logger := meta.Log("Akamai GTMv1", "populateAsAssignmentsObject")

	// pull apart List
	if asAssignmentsList, err := tools.GetInterfaceArrayValue("assignment", d); err != nil {
		logger.Errorf("Assignment not set: %s", err.Error())
	} else {
		asAssignmentsObjList := make([]*gtm.AsAssignment, len(asAssignmentsList)) // create new object list
		for i, v := range asAssignmentsList {
			asMap := v.(map[string]interface{})
			asAssignment := gtm.AsAssignment{}
			asAssignment.DatacenterId = asMap["datacenter_id"].(int)
			asAssignment.Nickname = asMap["nickname"].(string)
			if asMap["as_numbers"] != nil {
				ls := make([]int64, len(asMap["as_numbers"].([]interface{})))
				for i, sl := range asMap["as_numbers"].([]interface{}) {
					ls[i] = int64(sl.(int))
				}
				asAssignment.AsNumbers = ls
			}
			asAssignmentsObjList[i] = &asAssignment
		}
		as.Assignments = asAssignmentsObjList
	}
}

// create and populate Terraform asMap assignments schema
func populateTerraformAsAssignmentsState(d *schema.ResourceData, as *gtm.AsMap, m interface{}) {
	meta := akamai.Meta(m)
	logger := meta.Log("Akamai GTMv1", "populateTerraformAsAssignmentsState")

	objectInventory := make(map[int]*gtm.AsAssignment, len(as.Assignments))
	if len(as.Assignments) > 0 {
		for _, aObj := range as.Assignments {
			objectInventory[aObj.DatacenterId] = aObj
		}
	}
	if aStateList, err := tools.GetInterfaceArrayValue("assignment", d); err != nil {
		logger.Errorf("Assignment not set: %s", err.Error())
	} else {
		for _, aMap := range aStateList {
			a := aMap.(map[string]interface{})
			objIndex := a["datacenter_id"].(int)
			aObject, ok := objectInventory[objIndex]
			if !ok {
				logger.Warnf("As Assignment %d NOT FOUND in returned GTM Object", a["datacenter_id"])
				continue
			}
			a["datacenter_id"] = aObject.DatacenterId
			a["nickname"] = aObject.Nickname
			a["as_numbers"] = reconcileTerraformLists(a["as_numbers"].([]interface{}), convertInt64ToInterfaceList(aObject.AsNumbers, m), m)
			// remove object
			delete(objectInventory, objIndex)
		}
		if len(objectInventory) > 0 {
			logger.Debugf("As Assignment objects left...")
			// Objects not in the state yet. Add. Unfortunately, they not align with instance indices in the config
			for _, maObj := range objectInventory {
				aNew := map[string]interface{}{
					"datacenter_id": maObj.DatacenterId,
					"nickname":      maObj.Nickname,
					"as_numbers":    maObj.AsNumbers,
				}
				aStateList = append(aStateList, aNew)
			}
		}
		if err := d.Set("assignment", aStateList); err != nil {
			logger.Errorf("populateTerraformAsAssignmentsState failed: %s", err.Error())
		}
	}
}

// create and populate GTM ASmap DefaultDatacenter object
func populateAsDefaultDCObject(d *schema.ResourceData, as *gtm.AsMap, m interface{}) {
	meta := akamai.Meta(m)
	logger := meta.Log("Akamai GTMv1", "resourceGTMv1ASmapDelete")

	// pull apart List
	if asDefaultDCList, err := tools.GetInterfaceArrayValue("default_datacenter", d); err != nil {
		logger.Infof("No default datacenter specified: %s", err.Error())
	} else {
		if len(asDefaultDCList) > 0 {
			asDefaultDCObj := gtm.DatacenterBase{} // create new object
			asMap := asDefaultDCList[0].(map[string]interface{})
			if asMap["datacenter_id"] != nil && asMap["datacenter_id"].(int) != 0 {
				asDefaultDCObj.DatacenterId = asMap["datacenter_id"].(int)
				asDefaultDCObj.Nickname = asMap["nickname"].(string)
			} else {
				logger.Infof("No Default Datacenter specified")
				var nilInt int
				asDefaultDCObj.DatacenterId = nilInt
				asDefaultDCObj.Nickname = ""
			}
			as.DefaultDatacenter = &asDefaultDCObj
		}
	}
}

// create and populate Terraform asMap default_datacenter schema
func populateTerraformAsDefaultDCState(d *schema.ResourceData, as *gtm.AsMap, m interface{}) {
	meta := akamai.Meta(m)
	logger := meta.Log("Akamai GTMv1", "populateTerraformAsDefaultDCState")

	ddcListNew := make([]interface{}, 1)
	ddcNew := map[string]interface{}{
		"datacenter_id": as.DefaultDatacenter.DatacenterId,
		"nickname":      as.DefaultDatacenter.Nickname,
	}
	ddcListNew[0] = ddcNew
	if err := d.Set("default_datacenter", ddcListNew); err != nil {
		logger.Errorf("populateTerraformAsDefaultDCState failed: %s", err.Error())
	}
}
