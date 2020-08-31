package gtm

import (
	"fmt"
	"log"

	gtm "github.com/akamai/AkamaiOPEN-edgegrid-golang/configgtm-v1_4"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/schema"
)

func resourceGTMv1Geomap() *schema.Resource {
	return &schema.Resource{
		Create: resourceGTMv1GeomapCreate,
		Read:   resourceGTMv1GeomapRead,
		Update: resourceGTMv1GeomapUpdate,
		Delete: resourceGTMv1GeomapDelete,
		Exists: resourceGTMv1GeomapExists,
		Importer: &schema.ResourceImporter{
			State: resourceGTMv1GeomapImport,
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
						"countries": {
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

// utility func to parse Terraform property geoMap id
func parseResourceGeoMapId(id string) (string, string, error) {

	return parseResourceStringId(id)

}

// Create a new GTM GeoMap
func resourceGTMv1GeomapCreate(d *schema.ResourceData, meta interface{}) error {

	domain := d.Get("domain").(string)

	log.Printf("[INFO] [Akamai GTM] Creating geoMap [%s] in domain [%s]", d.Get("name").(string), domain)
	// Make sure Default Datacenter exists
	err := validateDefaultDC(d.Get("default_datacenter").([]interface{}), domain)
	if err != nil {
		return err
	}

	newGeo := populateNewGeoMapObject(d)
	log.Printf("[DEBUG] [Akamai GTMv1] Proposed New GeoMap: [%v]", newGeo)
	cStatus, err := newGeo.Create(domain)
	if err != nil {
		log.Printf("[ERROR] GeoMapCreate failed: %s", err.Error())
		return err
	}
	log.Printf("[DEBUG] [Akamai GTMv1] GeoMap Create status:")
	log.Printf("[DEBUG] [Akamai GTMv1] %v", cStatus.Status)
	if cStatus.Status.PropagationStatus == "DENIED" {
		return fmt.Errorf(cStatus.Status.Message)
	}
	if d.Get("wait_on_complete").(bool) {
		done, err := waitForCompletion(domain)
		if done {
			log.Printf("[INFO] [Akamai GTMv1] GeoMap Create completed")
		} else {
			if err == nil {
				log.Printf("[INFO] [Akamai GTMv1] GeoMap Create pending")
			} else {
				log.Printf("[WARNING] [Akamai GTMv1] GeoMap Create failed [%s]", err.Error())
				return err
			}
		}

	}

	// Give terraform the ID. Format domain:geoMap
	geoMapId := fmt.Sprintf("%s:%s", domain, cStatus.Resource.Name)
	log.Printf("[DEBUG] [Akamai GTMv1] Generated GeoMap GeoMap Id: %s", geoMapId)
	d.SetId(geoMapId)
	return resourceGTMv1GeomapRead(d, meta)

}

// read geoMap. updates state with entire API result configuration.
func resourceGTMv1GeomapRead(d *schema.ResourceData, _ interface{}) error {

	log.Printf("[DEBUG] [Akamai GTMv1] READ")
	log.Printf("[DEBUG] Reading [Akamai GTMv1] GeoMap: %s", d.Id())
	// retrieve the property and domain
	domain, geoMap, err := parseResourceGeoMapId(d.Id())
	if err != nil {
		return fmt.Errorf("invalid geoMap geoMap Id")
	}
	geo, err := gtm.GetGeoMap(geoMap, domain)
	if err != nil {
		log.Printf("[ERROR] [Akamai GTMv1] GeoMap Read error: %s", err.Error())
		return err
	}
	populateTerraformGeoMapState(d, geo)
	log.Printf("[DEBUG] [Akamai GTMv1] READ %v", geo)
	return nil
}

// Update GTM GeoMap
func resourceGTMv1GeomapUpdate(d *schema.ResourceData, meta interface{}) error {

	log.Printf("[DEBUG] [Akamai GTMv1] UPDATE")
	log.Printf("[DEBUG] Updating [Akamai GTMv1] GeoMap: %s", d.Id())
	// pull domain and geoMap out of id
	domain, geoMap, err := parseResourceGeoMapId(d.Id())
	if err != nil {
		return fmt.Errorf("invalid geoMap Id")
	}
	// Get existingGeoMap
	existGeo, err := gtm.GetGeoMap(geoMap, domain)
	if err != nil {
		log.Printf("[ERROR] GeoMapUpdate failed: %s", err.Error())
		return err
	}
	log.Printf("[DEBUG] Updating [Akamai GTMv1] GeoMap BEFORE: %v", existGeo)
	populateGeoMapObject(d, existGeo)
	log.Printf("[DEBUG] Updating [Akamai GTMv1] GeoMap PROPOSED: %v", existGeo)
	uStat, err := existGeo.Update(domain)
	if err != nil {
		log.Printf("[ERROR] GeoMapUpdate failed: %s", err.Error())
		return err
	}
	log.Printf("[DEBUG] [Akamai GTMv1] GeoMap Update  status:")
	log.Printf("[DEBUG] [Akamai GTMv1] %v", uStat)
	if uStat.PropagationStatus == "DENIED" {
		return fmt.Errorf(uStat.Message)
	}
	if d.Get("wait_on_complete").(bool) {
		done, err := waitForCompletion(domain)
		if done {
			log.Printf("[INFO] [Akamai GTMv1] GeoMap update completed")
		} else {
			if err == nil {
				log.Printf("[INFO] [Akamai GTMv1] GeoMap update pending")
			} else {
				log.Printf("[WARNING] [Akamai GTMv1] GeoMap update failed [%s]", err.Error())
				return err
			}
		}

	}

	return resourceGTMv1GeomapRead(d, meta)
}

// Import GTM GeoMap.
func resourceGTMv1GeomapImport(d *schema.ResourceData, _ interface{}) ([]*schema.ResourceData, error) {

	log.Printf("[INFO] [Akamai GTM] GeoMap [%s] Import", d.Id())
	// pull domain and geoMap out of geoMap id
	domain, geoMap, err := parseResourceGeoMapId(d.Id())
	if err != nil {
		return []*schema.ResourceData{d}, fmt.Errorf("invalid geoMap Id")
	}
	geo, err := gtm.GetGeoMap(geoMap, domain)
	if err != nil {
		return nil, err
	}
	d.Set("domain", domain)
	d.Set("wait_on_complete", true)
	populateTerraformGeoMapState(d, geo)

	// use same Id as passed in
	log.Printf("[INFO] [Akamai GTM] GeoMap [%s] [%s] Imported", d.Id(), d.Get("name"))
	return []*schema.ResourceData{d}, nil
}

// Delete GTM GeoMap.
func resourceGTMv1GeomapDelete(d *schema.ResourceData, _ interface{}) error {

	log.Printf("[DEBUG] [Akamai GTMv1] DELETE")
	log.Printf("[DEBUG] Deleting [Akamai GTMv1] GeoMap: %s", d.Id())
	// Get existing geoMap
	domain, geoMap, err := parseResourceGeoMapId(d.Id())
	if err != nil {
		return fmt.Errorf("invalid geoMap Id")
	}
	existGeo, err := gtm.GetGeoMap(geoMap, domain)
	if err != nil {
		log.Printf("[ERROR] GeoMapDelete failed: %s", err.Error())
		return err
	}
	log.Printf("[DEBUG] Deleting [Akamai GTMv1] GeoMap: %v", existGeo)
	uStat, err := existGeo.Delete(domain)
	if err != nil {
		log.Printf("[ERROR] GeoMapDelete failed: %s", err.Error())
		return err
	}
	log.Printf("[DEBUG] [Akamai GTMv1] GeoMap Delete status:")
	log.Printf("[DEBUG] [Akamai GTMv1] %v", uStat)
	if uStat.PropagationStatus == "DENIED" {
		return fmt.Errorf(uStat.Message)
	}
	if d.Get("wait_on_complete").(bool) {
		done, err := waitForCompletion(domain)
		if done {
			log.Printf("[INFO] [Akamai GTMv1] GeoMap delete completed")
		} else {
			if err == nil {
				log.Printf("[INFO] [Akamai GTMv1] GeoMap delete pending")
			} else {
				log.Printf("[WARNING] [Akamai GTMv1] GeoMap delete failed [%s]", err.Error())
				return err
			}
		}

	}

	// if successful ....
	d.SetId("")
	return nil
}

// Test GTM GeoMap existence
func resourceGTMv1GeomapExists(d *schema.ResourceData, _ interface{}) (bool, error) {

	log.Printf("[DEBUG] [Akamai GTMv1] Exists")
	// pull domain and geoMap out of geoMap id
	domain, geoMap, err := parseResourceGeoMapId(d.Id())
	if err != nil {
		return false, fmt.Errorf("invalid geoMap geoMap Id")
	}
	log.Printf("[DEBUG] [Akamai GTMv1] Searching for existing geoMap [%s] in domain %s", geoMap, domain)
	geo, err := gtm.GetGeoMap(geoMap, domain)
	return geo != nil, err
}

// Create and populate a new geoMap object from geoMap data
func populateNewGeoMapObject(d *schema.ResourceData) *gtm.GeoMap {

	geoObj := gtm.NewGeoMap(d.Get("name").(string))
	geoObj.DefaultDatacenter = &gtm.DatacenterBase{}
	geoObj.Assignments = make([]*gtm.GeoAssignment, 1)
	geoObj.Links = make([]*gtm.Link, 1)
	populateGeoMapObject(d, geoObj)

	return geoObj

}

// Populate existing geoMap object from geoMap data
func populateGeoMapObject(d *schema.ResourceData, geo *gtm.GeoMap) {

	if v, ok := d.GetOk("name"); ok {
		geo.Name = v.(string)
	}
	populateGeoAssignmentsObject(d, geo)
	populateGeoDefaultDCObject(d, geo)

}

// Populate Terraform state from provided GeoMap object
func populateTerraformGeoMapState(d *schema.ResourceData, geo *gtm.GeoMap) {

	// walk through all state elements
	d.Set("name", geo.Name)
	populateTerraformGeoAssignmentsState(d, geo)
	populateTerraformGeoDefaultDCState(d, geo)

}

// create and populate GTM GeoMap Assignments object
func populateGeoAssignmentsObject(d *schema.ResourceData, geo *gtm.GeoMap) {

	// pull apart List
	geoa := d.Get("assignment")
	if geoa != nil {
		geoAssignmentsList := geoa.([]interface{})
		geoAssignmentsObjList := make([]*gtm.GeoAssignment, len(geoAssignmentsList)) // create new object list
		for i, v := range geoAssignmentsList {
			geoMap := v.(map[string]interface{})
			geoAssignment := gtm.GeoAssignment{}
			geoAssignment.DatacenterId = geoMap["datacenter_id"].(int)
			geoAssignment.Nickname = geoMap["nickname"].(string)
			if geoMap["countries"] != nil {
				ls := make([]string, len(geoMap["countries"].([]interface{})))
				for i, sl := range geoMap["countries"].([]interface{}) {
					ls[i] = sl.(string)
				}
				geoAssignment.Countries = ls
			}
			geoAssignmentsObjList[i] = &geoAssignment
		}
		geo.Assignments = geoAssignmentsObjList
	}
}

// create and populate Terraform geoMap assigments schema
func populateTerraformGeoAssignmentsState(d *schema.ResourceData, geo *gtm.GeoMap) {

	objectInventory := make(map[int]*gtm.GeoAssignment, len(geo.Assignments))
	if len(geo.Assignments) > 0 {
		for _, aObj := range geo.Assignments {
			objectInventory[aObj.DatacenterId] = aObj
		}
	}
	aStateList := d.Get("assignment").([]interface{})
	for _, aMap := range aStateList {
		a := aMap.(map[string]interface{})
		objIndex := a["datacenter_id"].(int)
		aObject := objectInventory[objIndex]
		if aObject == nil {
			log.Printf("[WARNING] [Akamai GTMv1] Geo Assignment %d NOT FOUND in returned GTM Object", a["datacenter_id"])
			continue
		}
		a["datacenter_id"] = aObject.DatacenterId
		a["nickname"] = aObject.Nickname
		a["countries"] = reconcileTerraformLists(a["countries"].([]interface{}), convertStringToInterfaceList(aObject.Countries))
		// remove object
		delete(objectInventory, objIndex)
	}
	if len(objectInventory) > 0 {
		log.Printf("[DEBUG] [Akamai GTMv1] Geo Assignment objects left...")
		// Objects not in the state yet. Add. Unfortunately, they not align with instance indices in the config
		for _, maObj := range objectInventory {
			aNew := map[string]interface{}{
				"datacenter_id": maObj.DatacenterId,
				"nickname":      maObj.Nickname,
				"countries":     maObj.Countries,
			}
			aStateList = append(aStateList, aNew)
		}
	}
	d.Set("assignment", aStateList)

}

// create and populate GTM GeoMap DefaultDatacenter object
func populateGeoDefaultDCObject(d *schema.ResourceData, geo *gtm.GeoMap) {

	// pull apart List
	geodd := d.Get("default_datacenter")
	if geodd != nil && len(geodd.([]interface{})) > 0 {
		geoDefaultDCObj := gtm.DatacenterBase{} // create new object
		geoDefaultDCList := geodd.([]interface{})
		geoMap := geoDefaultDCList[0].(map[string]interface{})
		if geoMap["datacenter_id"] != nil && geoMap["datacenter_id"].(int) != 0 {
			geoDefaultDCObj.DatacenterId = geoMap["datacenter_id"].(int)
			geoDefaultDCObj.Nickname = geoMap["nickname"].(string)
		} else {
			log.Printf("[INFO] [Akamai GTMv1] No Default Datacenter specified")
			var nilInt int
			geoDefaultDCObj.DatacenterId = nilInt
			geoDefaultDCObj.Nickname = ""
		}
		geo.DefaultDatacenter = &geoDefaultDCObj
	}
}

// create and populate Terraform geoMap default_datacenter schema
func populateTerraformGeoDefaultDCState(d *schema.ResourceData, geo *gtm.GeoMap) {

	ddcListNew := make([]interface{}, 1)
	ddcNew := map[string]interface{}{
		"datacenter_id": geo.DefaultDatacenter.DatacenterId,
		"nickname":      geo.DefaultDatacenter.Nickname,
	}
	ddcListNew[0] = ddcNew
	d.Set("default_datacenter", ddcListNew)
}
