package akamai

import (
	"errors"
	"fmt"
	gtm "github.com/akamai/AkamaiOPEN-edgegrid-golang/configgtm-v1_4"
	"github.com/hashicorp/terraform/helper/schema"
	"log"
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
			"default_datacenter": &schema.Schema{
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
			"assignment": &schema.Schema{
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
						"blocks": &schema.Schema{
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

// utility func to parse Terraform property cidrMap id
func parseResourceCidrMapId(id string) (string, string, error) {

	return parseResourceStringId(id)

}

// Create a new GTM CidrMap
func resourceGTMv1CidrMapCreate(d *schema.ResourceData, meta interface{}) error {

	domain := d.Get("domain").(string)

	log.Printf("[INFO] [Akamai GTM] Creating cidrMap [%s] in domain [%s]", d.Get("name").(string), domain)
	newCidr := populateNewCidrMapObject(d)
	log.Printf("[DEBUG] [Akamai GTMv1] Proposed New CidrMap: [%v]", newCidr)
	cStatus, err := newCidr.Create(domain)
	if err != nil {
		log.Printf("[ERROR] CidrMapCreate failed: %s", err.Error())
		return err
	}
	log.Printf("[DEBUG] [Akamai GTMv1] CidrMap Create status:")
	log.Printf("[DEBUG] [Akamai GTMv1] %v", cStatus.Status)
	if cStatus.Status.PropagationStatus == "DENIED" {
		return errors.New(cStatus.Status.Message)
	}
	if d.Get("wait_on_complete").(bool) {
		done, err := waitForCompletion(domain)
		if done {
			log.Printf("[INFO] [Akamai GTMv1] CidrMap Create completed")
		} else {
			if err == nil {
				log.Printf("[INFO] [Akamai GTMv1] CidrMap Create pending")
			} else {
				log.Printf("[WARNING] [Akamai GTMv1] CidrMap Create failed [%s]", err.Error())
				return err
			}
		}

	}

	// Give terraform the ID. Format domain:cidrMap
	cidrMapId := fmt.Sprintf("%s:%s", domain, cStatus.Resource.Name)
	log.Printf("[DEBUG] [Akamai GTMv1] Generated CidrMap CidrMap Id: %s", cidrMapId)
	d.SetId(cidrMapId)
	return resourceGTMv1CidrMapRead(d, meta)

}

// read cidrMap. updates state with entire API result configuration.
func resourceGTMv1CidrMapRead(d *schema.ResourceData, meta interface{}) error {

	log.Printf("[DEBUG] [Akamai GTMv1] READ")
	log.Printf("[DEBUG] Reading [Akamai GTMv1] CidrMap: %s", d.Id())
	// retrieve the property and domain
	domain, cidrMap, err := parseResourceCidrMapId(d.Id())
	if err != nil {
		return errors.New("Invalid cidrMap cidrMap Id")
	}
	cidr, err := gtm.GetCidrMap(cidrMap, domain)
	if err != nil {
		log.Printf("[ERROR] CidrMap Read error: %s", err.Error())
		return err
	}
	populateTerraformCidrMapState(d, cidr)
	log.Printf("[DEBUG] [Akamai GTMv1] READ %v", cidr)
	return nil
}

// Update GTM CidrMap
func resourceGTMv1CidrMapUpdate(d *schema.ResourceData, meta interface{}) error {

	log.Printf("[DEBUG] [Akamai GTMv1] UPDATE")
	log.Printf("[DEBUG] Updating [Akamai GTMv1] CidrMap: %s", d.Id())
	// pull domain and cidrMap out of id
	domain, cidrMap, err := parseResourceCidrMapId(d.Id())
	if err != nil {
		return errors.New("Invalid cidrMap Id")
	}
	// Get existingCidrMap
	existCidr, err := gtm.GetCidrMap(cidrMap, domain)
	if err != nil {
		log.Printf("[ERROR] CidrMapUpdate failed: %s", err.Error())
		return err
	}
	log.Printf("[DEBUG] Updating [Akamai GTMv1] CidrMap BEFORE: %v", existCidr)
	populateCidrMapObject(d, existCidr)
	log.Printf("[DEBUG] Updating [Akamai GTMv1] CidrMap PROPOSED: %v", existCidr)
	uStat, err := existCidr.Update(domain)
	if err != nil {
		log.Printf("[ERROR] CidrMapUpdate failed: %s", err.Error())
		return err
	}
	log.Printf("[DEBUG] [Akamai GTMv1] CidrMap Update  status:")
	log.Printf("[DEBUG] [Akamai GTMv1] %v", uStat)
	if uStat.PropagationStatus == "DENIED" {
		return errors.New(uStat.Message)
	}
	if d.Get("wait_on_complete").(bool) {
		done, err := waitForCompletion(domain)
		if done {
			log.Printf("[INFO] [Akamai GTMv1] CidrMap update completed")
		} else {
			if err == nil {
				log.Printf("[INFO] [Akamai GTMv1] CidrMap update pending")
			} else {
				log.Printf("[WARNING] [Akamai GTMv1] CidrMap update failed [%s]", err.Error())
				return err
			}
		}

	}

	return resourceGTMv1CidrMapRead(d, meta)
}

// Import GTM CidrMap.
func resourceGTMv1CidrMapImport(d *schema.ResourceData, meta interface{}) ([]*schema.ResourceData, error) {

	log.Printf("[INFO] [Akamai GTM] CidrMap [%s] Import", d.Id())
	// pull domain and cidrMap out of cidrMap id
	domain, cidrMap, err := parseResourceCidrMapId(d.Id())
	if err != nil {
		return []*schema.ResourceData{d}, errors.New("Invalid cidrMap cidrMap Id")
	}
	cidr, err := gtm.GetCidrMap(cidrMap, domain)
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
func resourceGTMv1CidrMapDelete(d *schema.ResourceData, meta interface{}) error {

	log.Printf("[DEBUG] [Akamai GTMv1] DELETE")
	log.Printf("[DEBUG] Deleting [Akamai GTMv1] CidrMap: %s", d.Id())
	// Get existing cidrMap
	domain, cidrMap, err := parseResourceCidrMapId(d.Id())
	if err != nil {
		return errors.New("Invalid cidrMap Id")
	}
	existCidr, err := gtm.GetCidrMap(cidrMap, domain)
	if err != nil {
		log.Printf("[ERROR] CidrMapDelete failed: %s", err.Error())
		return err
	}
	log.Printf("[DEBUG] Deleting [Akamai GTMv1] CidrMap: %v", existCidr)
	uStat, err := existCidr.Delete(domain)
	if err != nil {
		log.Printf("[ERROR] CidrMapDelete failed: %s", err.Error())
		return err
	}
	log.Printf("[DEBUG] [Akamai GTMv1] CidrMap Delete status:")
	log.Printf("[DEBUG] [Akamai GTMv1] %v", uStat)
	if uStat.PropagationStatus == "DENIED" {
		return errors.New(uStat.Message)
	}
	if d.Get("wait_on_complete").(bool) {
		done, err := waitForCompletion(domain)
		if done {
			log.Printf("[INFO] [Akamai GTMv1] CidrMap delete completed")
		} else {
			if err == nil {
				log.Printf("[INFO] [Akamai GTMv1] CidrMap delete pending")
			} else {
				log.Printf("[WARNING] [Akamai GTMv1] CidrMap delete failed [%s]", err.Error())
				return err
			}
		}

	}

	// if succcessful ....
	d.SetId("")
	return nil
}

// Test GTM CidrMap existance
func resourceGTMv1CidrMapExists(d *schema.ResourceData, meta interface{}) (bool, error) {

	log.Printf("[DEBUG] [Akamai GTMv1] Exists")
	// pull domain and cidrMap out of cidrMap id
	domain, cidrMap, err := parseResourceCidrMapId(d.Id())
	if err != nil {
		return false, errors.New("Invalid cidrMap cidrMap Id")
	}
	log.Printf("[DEBUG] [Akamai GTMv1] Searching for existing cidrMap [%s] in domain %s", cidrMap, domain)
	cidr, err := gtm.GetCidrMap(cidrMap, domain)
	return cidr != nil, err
}

// Create and populate a new cidrMap object from cidrMap data
func populateNewCidrMapObject(d *schema.ResourceData) *gtm.CidrMap {

	cidrObj := gtm.NewCidrMap(d.Get("name").(string))
	cidrObj.DefaultDatacenter = &gtm.DatacenterBase{}
	cidrObj.Assignments = make([]*gtm.CidrAssignment, 0)
	cidrObj.Links = make([]*gtm.Link, 1)
	populateCidrMapObject(d, cidrObj)

	return cidrObj

}

// Populate existing cidrMap object from cidrMap data
func populateCidrMapObject(d *schema.ResourceData, cidr *gtm.CidrMap) {

	if v, ok := d.GetOk("name"); ok {
		cidr.Name = v.(string)
	}
	populateCidrAssignmentsObject(d, cidr)
	populateCidrDefaultDCObject(d, cidr)

}

// Populate Terraform state from provided CidrMap object
func populateTerraformCidrMapState(d *schema.ResourceData, cidr *gtm.CidrMap) {

	// walk thru all state elements
	d.Set("name", cidr.Name)
	populateTerraformCidrAssignmentsState(d, cidr)
	populateTerraformCidrDefaultDCState(d, cidr)

}

// create and populate GTM CidrMap Assignments object
func populateCidrAssignmentsObject(d *schema.ResourceData, cidr *gtm.CidrMap) {

	// pull apart List
	cassgns := d.Get("assignment")
	if cassgns != nil {
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

// create and populate Terraform cidrMap assigments schema
func populateTerraformCidrAssignmentsState(d *schema.ResourceData, cidr *gtm.CidrMap) {

	objectInventory := make(map[int]*gtm.CidrAssignment, len(cidr.Assignments))
	if len(cidr.Assignments) > 0 {
		for _, aObj := range cidr.Assignments {
			objectInventory[aObj.DatacenterId] = aObj
		}
	}
	aStateList := d.Get("assignment").([]interface{})
	for _, aMap := range aStateList {
		a := aMap.(map[string]interface{})
		objIndex := a["datacenter_id"].(int)
		aObject := objectInventory[objIndex]
		if aObject == nil {
			log.Printf("[WARNING] [Akamai GTMv1] Cidr Assignment %d NOT FOUND in returned GTM Object", a["datacenter_id"])
			continue
		}
		a["datacenter_id"] = aObject.DatacenterId
		a["nickname"] = aObject.Nickname
		a["blocks"] = aObject.Blocks
		// remove object
		delete(objectInventory, objIndex)
	}
	if len(objectInventory) > 0 {
		log.Printf("[DEBUG] [Akamai GTMv1] CIDR Assignment objects left...")
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
	d.Set("assignment", aStateList)

}

// create and populate GTM CidrMap DefaultDatacenter object
func populateCidrDefaultDCObject(d *schema.ResourceData, cidr *gtm.CidrMap) {

	// pull apart List
	cidrdd := d.Get("default_datacenter")
	if cidrdd != nil && len(cidrdd.([]interface{})) > 0 {
		cidrDefaultDCObj := gtm.DatacenterBase{} // create new object
		cidrDefaultDCList := cidrdd.([]interface{})
		cidrddMap := cidrDefaultDCList[0].(map[string]interface{})
		if cidrddMap["datacenter_id"] != nil && cidrddMap["datacenter_id"].(int) != 0 {
			cidrDefaultDCObj.DatacenterId = cidrddMap["datacenter_id"].(int)
			cidrDefaultDCObj.Nickname = cidrddMap["nickname"].(string)
		} else {
			log.Printf("[INFO] [Akamai GTMv1] No Default Datacenter specified")
			var nilInt int
			cidrDefaultDCObj.DatacenterId = nilInt
			cidrDefaultDCObj.Nickname = ""
		}
		cidr.DefaultDatacenter = &cidrDefaultDCObj
	}
}

// create and populate Terraform cidrMap default_datacenter schema
func populateTerraformCidrDefaultDCState(d *schema.ResourceData, cidr *gtm.CidrMap) {

	ddcListNew := make([]interface{}, 1)
	ddcNew := map[string]interface{}{
		"datacenter_id": cidr.DefaultDatacenter.DatacenterId,
		"nickname":      cidr.DefaultDatacenter.Nickname,
	}
	ddcListNew[0] = ddcNew
	d.Set("default_datacenter", ddcListNew)

}
