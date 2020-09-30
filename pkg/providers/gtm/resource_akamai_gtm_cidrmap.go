package gtm

import (
	"fmt"

	gtm "github.com/akamai/AkamaiOPEN-edgegrid-golang/configgtm-v1_4"
	"github.com/akamai/terraform-provider-akamai/v2/pkg/akamai"
	"github.com/akamai/terraform-provider-akamai/v2/pkg/tools"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/schema"
)

func resourceGTMv1Cidrmap() *schema.Resource {
	return &schema.Resource{
		Create: resourceGTMv1CidrMapCreate,
		Read:   resourceGTMv1CidrMapRead,
		Update: resourceGTMv1CidrMapUpdate,
		Delete: resourceGTMv1CidrMapDelete,
		Exists: resourceGTMv1CidrMapExists,
		Importer: &schema.ResourceImporter{
			State: resourceGTMv1CidrMapImport,
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
				Type:     schema.TypeList,
				Required: true,
				MaxItems: 1,
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
						"blocks": {
							Type:     schema.TypeList,
							Elem:     &schema.Schema{Type: schema.TypeString},
							Optional: true,
						},
					},
				},
			},
		},
	}
}

// Create a new GTM CidrMap
func resourceGTMv1CidrMapCreate(d *schema.ResourceData, m interface{}) error {
	meta := akamai.Meta(m)
	logger := meta.Log("Akamai GTMv1", "resourceGTMv1CidrMapCreate")

	domain, err := tools.GetStringValue("domain", d)
	if err != nil {
		return err
	}

	name, err := tools.GetStringValue("name", d)
	if err != nil {
		return err
	}

	logger.Infof("Creating cidrMap [%s] in domain [%s]", name, domain)

	// Make sure Default Datacenter exists
	if defaultDatacenter, err := tools.GetInterfaceArrayValue("default_datacenter", d); err != nil {
		return err
	} else {
		if validateDefaultDC(defaultDatacenter, domain) != nil {
			return err
		}
	}

	newCidr := populateNewCidrMapObject(d, m)
	logger.Debugf("Proposed New CidrMap: [%v]", newCidr)
	cStatus, err := newCidr.Create(domain)
	if err != nil {
		logger.Errorf("CidrMapCreate failed: %s", err.Error())
		return err
	}
	logger.Debugf("CidrMap Create status:")
	logger.Debugf("%v", cStatus.Status)
	if cStatus.Status.PropagationStatus == "DENIED" {
		return fmt.Errorf(cStatus.Status.Message)
	}
	if waitOnComplete, err := tools.GetBoolValue("wait_on_complete", d); err != nil {
		return err
	} else {
		if waitOnComplete {
			done, err := waitForCompletion(domain, m)
			if done {
				logger.Infof("CidrMap Create completed")
			} else {
				if err == nil {
					logger.Infof("CidrMap Create pending")
				} else {
					logger.Errorf("CidrMap Create failed [%s]", err.Error())
					return err
				}
			}
		}
	}

	// Give terraform the ID. Format domain:cidrMap
	cidrMapID := fmt.Sprintf("%s:%s", domain, cStatus.Resource.Name)
	logger.Debugf("Generated CidrMap CidrMap Id: %s", cidrMapID)
	d.SetId(cidrMapID)
	return resourceGTMv1CidrMapRead(d, m)

}

// read cidrMap. updates state with entire API result configuration.
func resourceGTMv1CidrMapRead(d *schema.ResourceData, m interface{}) error {
	meta := akamai.Meta(m)
	logger := meta.Log("Akamai GTMv1", "resourceGTMv1CidrMapRead")

	logger.Debugf("READ")
	logger.Debugf("Reading CidrMap: %s", d.Id())
	// retrieve the property and domain
	domain, cidrMap, err := parseResourceStringId(d.Id())
	if err != nil {
		return fmt.Errorf("invalid cidrMap cidrMap ID")
	}
	cidr, err := gtm.GetCidrMap(cidrMap, domain)
	if err != nil {
		logger.Errorf("CidrMap Read error: %s", err.Error())
		return err
	}
	populateTerraformCidrMapState(d, cidr, m)
	logger.Debugf("READ %v", cidr)
	return nil
}

// Update GTM CidrMap
func resourceGTMv1CidrMapUpdate(d *schema.ResourceData, m interface{}) error {
	meta := akamai.Meta(m)
	logger := meta.Log("Akamai GTMv1", "resourceGTMv1CidrMapUpdate")

	logger.Debugf("UPDATE")
	logger.Debugf("Updating CidrMap: %s", d.Id())
	// pull domain and cidrMap out of id
	domain, cidrMap, err := parseResourceStringId(d.Id())
	if err != nil {
		return fmt.Errorf("invalid cidrMap ID")
	}
	// Get existingCidrMap
	existCidr, err := gtm.GetCidrMap(cidrMap, domain)
	if err != nil {
		logger.Errorf("CidrMapUpdate failed: %s", err.Error())
		return err
	}
	logger.Debugf("Updating CidrMap BEFORE: %v", existCidr)
	populateCidrMapObject(d, existCidr, m)
	logger.Debugf("Updating CidrMap PROPOSED: %v", existCidr)
	uStat, err := existCidr.Update(domain)
	if err != nil {
		logger.Errorf("CidrMapUpdate failed: %s", err.Error())
		return err
	}
	logger.Debugf("CidrMap Update  status:")
	logger.Debugf("%v", uStat)
	if uStat.PropagationStatus == "DENIED" {
		return fmt.Errorf(uStat.Message)
	}

	if waitOnComplete, err := tools.GetBoolValue("wait_on_complete", d); err != nil {
		return err
	} else {
		if waitOnComplete {
			done, err := waitForCompletion(domain, m)
			if done {
				logger.Infof("CidrMap update completed")
			} else {
				if err == nil {
					logger.Infof("CidrMap update pending")
				} else {
					logger.Errorf("CidrMap update failed [%s]", err.Error())
					return err
				}
			}
		}
	}

	return resourceGTMv1CidrMapRead(d, m)
}

// Import GTM CidrMap.
func resourceGTMv1CidrMapImport(d *schema.ResourceData, m interface{}) ([]*schema.ResourceData, error) {
	meta := akamai.Meta(m)
	logger := meta.Log("Akamai GTMv1", "resourceGTMv1CidrMapImport")

	logger.Infof("CidrMap [%s] Import", d.Id())
	// pull domain and cidrMap out of cidrMap id
	domain, cidrMap, err := parseResourceStringId(d.Id())
	if err != nil {
		return []*schema.ResourceData{d}, fmt.Errorf("invalid cidrMap cidrMap ID")
	}
	cidr, err := gtm.GetCidrMap(cidrMap, domain)
	if err != nil {
		return nil, err
	}
	if err := d.Set("domain", domain); err != nil {
		logger.Errorf("resourceGTMv1CidrMapImport failed: %s", err.Error())
	}
	if err := d.Set("wait_on_complete", true); err != nil {
		logger.Errorf("resourceGTMv1CidrMapImport failed: %s", err.Error())
	}
	populateTerraformCidrMapState(d, cidr, m)

	// use same Id as passed in
	logger.Infof("CidrMap [%s] [%s] Imported", d.Id(), d.Get("name"))
	return []*schema.ResourceData{d}, nil
}

// Delete GTM CidrMap.
func resourceGTMv1CidrMapDelete(d *schema.ResourceData, m interface{}) error {
	meta := akamai.Meta(m)
	logger := meta.Log("Akamai GTMv1", "resourceGTMv1CidrMapDelete")

	logger.Debugf("DELETE")
	logger.Debugf("Deleting CidrMap: %s", d.Id())
	// Get existing cidrMap
	domain, cidrMap, err := parseResourceStringId(d.Id())
	if err != nil {
		return fmt.Errorf("invalid cidrMap ID")
	}
	existCidr, err := gtm.GetCidrMap(cidrMap, domain)
	if err != nil {
		logger.Errorf("CidrMapDelete failed: %s", err.Error())
		return err
	}
	logger.Debugf("Deleting CidrMap: %v", existCidr)
	uStat, err := existCidr.Delete(domain)
	if err != nil {
		logger.Errorf("CidrMapDelete failed: %s", err.Error())
		return err
	}
	logger.Debugf("CidrMap Delete status:")
	logger.Debugf("%v", uStat)
	if uStat.PropagationStatus == "DENIED" {
		return fmt.Errorf(uStat.Message)
	}

	if waitOnComplete, err := tools.GetBoolValue("wait_on_complete", d); err != nil {
		return err
	} else {
		if waitOnComplete {
			done, err := waitForCompletion(domain, m)
			if done {
				logger.Infof("CidrMap delete completed")
			} else {
				if err == nil {
					logger.Infof("CidrMap delete pending")
				} else {
					logger.Errorf("CidrMap delete failed [%s]", err.Error())
					return err
				}
			}

		}
	}

	// if successful ....
	d.SetId("")
	return nil
}

// Test GTM CidrMap existence
func resourceGTMv1CidrMapExists(d *schema.ResourceData, m interface{}) (bool, error) {
	meta := akamai.Meta(m)
	logger := meta.Log("Akamai GTMv1", "resourceGTMv1CidrMapExists")

	logger.Debugf("Exists")
	// pull domain and cidrMap out of cidrMap id
	domain, cidrMap, err := parseResourceStringId(d.Id())
	if err != nil {
		return false, fmt.Errorf("invalid cidrMap cidrMap ID")
	}
	logger.Debugf("Searching for existing cidrMap [%s] in domain %s", cidrMap, domain)
	cidr, err := gtm.GetCidrMap(cidrMap, domain)
	return cidr != nil, err
}

// Create and populate a new cidrMap object from cidrMap data
func populateNewCidrMapObject(d *schema.ResourceData, m interface{}) *gtm.CidrMap {
	meta := akamai.Meta(m)
	logger := meta.Log("Akamai GTMv1", "populateNewCidrMapObject")

	cidrMapName, err := tools.GetStringValue("name", d)
	if err != nil {
		logger.Errorf("domain name not found in ResourceData: %s", err.Error())
	}

	cidrObj := gtm.NewCidrMap(cidrMapName)
	cidrObj.DefaultDatacenter = &gtm.DatacenterBase{}
	cidrObj.Assignments = make([]*gtm.CidrAssignment, 0)
	cidrObj.Links = make([]*gtm.Link, 1)
	populateCidrMapObject(d, cidrObj, m)

	return cidrObj

}

// Populate existing cidrMap object from cidrMap data
func populateCidrMapObject(d *schema.ResourceData, cidr *gtm.CidrMap, m interface{}) {

	if v, err := tools.GetStringValue("name", d); err == nil {
		cidr.Name = v
	}
	populateCidrAssignmentsObject(d, cidr)
	populateCidrDefaultDCObject(d, cidr, m)

}

// Populate Terraform state from provided CidrMap object
func populateTerraformCidrMapState(d *schema.ResourceData, cidr *gtm.CidrMap, m interface{}) {
	meta := akamai.Meta(m)
	logger := meta.Log("Akamai GTMv1", "populateTerraformCidrMapState")

	// walk through all state elements
	if err := d.Set("name", cidr.Name); err != nil {
		logger.Errorf("populateTerraformCidrMapState failed: %s", err.Error())
	}
	populateTerraformCidrAssignmentsState(d, cidr, m)
	populateTerraformCidrDefaultDCState(d, cidr, m)

}

// create and populate GTM CidrMap Assignments object
func populateCidrAssignmentsObject(d *schema.ResourceData, cidr *gtm.CidrMap) {

	// pull apart List
	if cassgns := d.Get("assignment"); cassgns != nil {
		cidrAssignmentsList := cassgns.([]interface{})
		cidrAssignmentsObjList := make([]*gtm.CidrAssignment, len(cidrAssignmentsList)) // create new object list
		for i, v := range cidrAssignmentsList {
			cidrMap := v.(map[string]interface{})
			cidrAssignment := gtm.CidrAssignment{}
			cidrAssignment.DatacenterId = cidrMap["datacenter_id"].(int)
			cidrAssignment.Nickname = cidrMap["nickname"].(string)
			if cidrMap["blocks"] != nil {
				ls := make([]string, len(cidrMap["blocks"].([]interface{})))
				for i, sl := range cidrMap["blocks"].([]interface{}) {
					ls[i] = sl.(string)
				}
				cidrAssignment.Blocks = ls
			}
			cidrAssignmentsObjList[i] = &cidrAssignment
		}
		cidr.Assignments = cidrAssignmentsObjList
	}
}

// create and populate Terraform cidrMap assignments schema
func populateTerraformCidrAssignmentsState(d *schema.ResourceData, cidr *gtm.CidrMap, m interface{}) {
	meta := akamai.Meta(m)
	logger := meta.Log("Akamai GTMv1", "populateTerraformCidrAssignmentsState")

	objectInventory := make(map[int]*gtm.CidrAssignment, len(cidr.Assignments))
	if len(cidr.Assignments) > 0 {
		for _, aObj := range cidr.Assignments {
			objectInventory[aObj.DatacenterId] = aObj
		}
	}
	aStateList, err := tools.GetInterfaceArrayValue("assignment", d)
	if err != nil {
		logger.Errorf("Cidr Assignment list NOT FOUND in ResourceData: %s", err.Error())
	}
	for _, aMap := range aStateList {
		a := aMap.(map[string]interface{})
		objIndex := a["datacenter_id"].(int)
		aObject := objectInventory[objIndex]
		if aObject == nil {
			logger.Warnf("Cidr Assignment %d NOT FOUND in returned GTM Object", a["datacenter_id"])
			continue
		}
		a["datacenter_id"] = aObject.DatacenterId
		a["nickname"] = aObject.Nickname
		a["blocks"] = reconcileTerraformLists(a["blocks"].([]interface{}), convertStringToInterfaceList(aObject.Blocks, m), m)
		// remove object
		delete(objectInventory, objIndex)
	}
	if len(objectInventory) > 0 {
		logger.Debugf("CIDR Assignment objects left...")
		// Objects not in the state yet. Add. Unfortunately, they not align with instance indices in the config
		for _, maObj := range objectInventory {
			aNew := map[string]interface{}{
				"datacenter_id": maObj.DatacenterId,
				"nickname":      maObj.Nickname,
				"blocks":        maObj.Blocks,
			}
			aStateList = append(aStateList, aNew)
		}
	}
	if err := d.Set("assignment", aStateList); err != nil {
		logger.Errorf("populateTerraformCidrAssignmentsState failed: %s", err.Error())
	}
}

// create and populate GTM CidrMap DefaultDatacenter object
func populateCidrDefaultDCObject(d *schema.ResourceData, cidr *gtm.CidrMap, m interface{}) {
	meta := akamai.Meta(m)
	logger := meta.Log("Akamai GTMv1", "populateCidrDefaultDCObject")

	// pull apart List
	if cidrDefaultDCList, err := tools.GetInterfaceArrayValue("default_datacenter", d); err != nil {
		logger.Infof("No default datacenter specified: %s", err.Error())
	} else {
		if len(cidrDefaultDCList) > 0 {
			cidrDefaultDCObj := gtm.DatacenterBase{} // create new object
			cidrddMap := cidrDefaultDCList[0].(map[string]interface{})
			if cidrddMap["datacenter_id"] != nil && cidrddMap["datacenter_id"].(int) != 0 {
				cidrDefaultDCObj.DatacenterId = cidrddMap["datacenter_id"].(int)
				cidrDefaultDCObj.Nickname = cidrddMap["nickname"].(string)
			} else {
				logger.Infof("No Default Datacenter specified")
				var nilInt int
				cidrDefaultDCObj.DatacenterId = nilInt
				cidrDefaultDCObj.Nickname = ""
			}
			cidr.DefaultDatacenter = &cidrDefaultDCObj
		}
	}
}

// create and populate Terraform cidrMap default_datacenter schema
func populateTerraformCidrDefaultDCState(d *schema.ResourceData, cidr *gtm.CidrMap, m interface{}) {
	meta := akamai.Meta(m)
	logger := meta.Log("Akamai GTMv1", "populateTerraformCidrDefaultDCState")

	ddcListNew := make([]interface{}, 1)
	ddcNew := map[string]interface{}{
		"datacenter_id": cidr.DefaultDatacenter.DatacenterId,
		"nickname":      cidr.DefaultDatacenter.Nickname,
	}
	ddcListNew[0] = ddcNew
	if err := d.Set("default_datacenter", ddcListNew); err != nil {
		logger.Errorf("populateTerraformCidrDefaultDCState failed: %s", err.Error())
	}
}
