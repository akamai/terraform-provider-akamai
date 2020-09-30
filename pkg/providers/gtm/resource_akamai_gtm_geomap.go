package gtm

import (
	"fmt"

	gtm "github.com/akamai/AkamaiOPEN-edgegrid-golang/configgtm-v1_4"
	"github.com/akamai/terraform-provider-akamai/v2/pkg/akamai"
	"github.com/akamai/terraform-provider-akamai/v2/pkg/tools"
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

// Create a new GTM GeoMap
func resourceGTMv1GeomapCreate(d *schema.ResourceData, m interface{}) error {
	meta := akamai.Meta(m)
	logger := meta.Log("Akamai GTMv1", "resourceGTMv1GeomapCreate")

	domain, err := tools.GetStringValue("domain", d)
	if err != nil {
		return err
	}

	name, _ := tools.GetStringValue("name", d)
	logger.Infof("[Akamai GTM] Creating geoMap [%s] in domain [%s]", name, domain)
	// Make sure Default Datacenter exists
	geoDefaultDCList, err := tools.GetInterfaceArrayValue("default_datacenter", d)
	err = validateDefaultDC(geoDefaultDCList, domain)
	if err != nil {
		return err
	}

	newGeo := populateNewGeoMapObject(d, m)
	logger.Debugf("Proposed New GeoMap: [%v]", newGeo)
	cStatus, err := newGeo.Create(domain)
	if err != nil {
		logger.Errorf("GeoMapCreate failed: %s", err.Error())
		return err
	}
	logger.Debugf("GeoMap Create status:")
	logger.Debugf("%v", cStatus.Status)
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
			logger.Infof("GeoMap Create completed")
		} else {
			if err == nil {
				logger.Infof("GeoMap Create pending")
			} else {
				logger.Warnf("GeoMap Create failed [%s]", err.Error())
				return err
			}
		}

	}

	// Give terraform the ID. Format domain:geoMap
	geoMapId := fmt.Sprintf("%s:%s", domain, cStatus.Resource.Name)
	logger.Debugf("Generated GeoMap GeoMap Id: %s", geoMapId)
	d.SetId(geoMapId)
	return resourceGTMv1GeomapRead(d, m)

}

// read geoMap. updates state with entire API result configuration.
func resourceGTMv1GeomapRead(d *schema.ResourceData, m interface{}) error {
	meta := akamai.Meta(m)
	logger := meta.Log("Akamai GTMv1", "resourceGTMv1GeomapRead")

	logger.Debugf("Reading GeoMap: %s", d.Id())
	// retrieve the property and domain
	domain, geoMap, err := parseResourceStringId(d.Id())
	if err != nil {
		return fmt.Errorf("invalid geoMap geoMap ID")
	}
	geo, err := gtm.GetGeoMap(geoMap, domain)
	if err != nil {
		logger.Errorf("GeoMap Read error: %s", err.Error())
		return err
	}
	populateTerraformGeoMapState(d, geo, m)
	logger.Debugf("READ %v", geo)
	return nil
}

// Update GTM GeoMap
func resourceGTMv1GeomapUpdate(d *schema.ResourceData, m interface{}) error {
	meta := akamai.Meta(m)
	logger := meta.Log("Akamai GTMv1", "resourceGTMv1GeomapUpdate")

	logger.Debugf("UPDATE")
	logger.Debugf("Updating GeoMap: %s", d.Id())
	// pull domain and geoMap out of id
	domain, geoMap, err := parseResourceStringId(d.Id())
	if err != nil {
		return fmt.Errorf("invalid geoMap ID")
	}
	// Get existingGeoMap
	existGeo, err := gtm.GetGeoMap(geoMap, domain)
	if err != nil {
		logger.Errorf("GeoMapUpdate failed: %s", err.Error())
		return err
	}
	logger.Debugf("Updating GeoMap BEFORE: %v", existGeo)
	populateGeoMapObject(d, existGeo, m)
	logger.Debugf("Updating GeoMap PROPOSED: %v", existGeo)
	uStat, err := existGeo.Update(domain)
	if err != nil {
		logger.Errorf("GeoMapUpdate failed: %s", err.Error())
		return err
	}
	logger.Debugf("GeoMap Update  status:")
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
			logger.Infof("GeoMap update completed")
		} else {
			if err == nil {
				logger.Infof("GeoMap update pending")
			} else {
				logger.Warnf("GeoMap update failed [%s]", err.Error())
				return err
			}
		}

	}

	return resourceGTMv1GeomapRead(d, m)
}

// Import GTM GeoMap.
func resourceGTMv1GeomapImport(d *schema.ResourceData, m interface{}) ([]*schema.ResourceData, error) {
	meta := akamai.Meta(m)
	logger := meta.Log("Akamai GTMv1", "resourceGTMv1GeomapImport")

	logger.Infof("[Akamai GTM] GeoMap [%s] Import", d.Id())
	// pull domain and geoMap out of geoMap id
	domain, geoMap, err := parseResourceStringId(d.Id())
	if err != nil {
		return []*schema.ResourceData{d}, fmt.Errorf("invalid geoMap ID")
	}
	geo, err := gtm.GetGeoMap(geoMap, domain)
	if err != nil {
		return nil, err
	}
	if err := d.Set("domain", domain); err != nil {
		return nil, err
	}
	if err := d.Set("wait_on_complete", true); err != nil {
		return nil, err
	}
	populateTerraformGeoMapState(d, geo, m)

	// use same Id as passed in
	name, _ := tools.GetStringValue("name", d)
	logger.Infof("[Akamai GTM] GeoMap [%s] [%s] Imported", d.Id(), name)
	return []*schema.ResourceData{d}, nil
}

// Delete GTM GeoMap.
func resourceGTMv1GeomapDelete(d *schema.ResourceData, m interface{}) error {
	meta := akamai.Meta(m)
	logger := meta.Log("Akamai GTMv1", "resourceGTMv1GeomapDelete")

	logger.Debugf("DELETE")
	logger.Debugf("Deleting GeoMap: %s", d.Id())
	// Get existing geoMap
	domain, geoMap, err := parseResourceStringId(d.Id())
	if err != nil {
		return fmt.Errorf("invalid geoMap ID")
	}
	existGeo, err := gtm.GetGeoMap(geoMap, domain)
	if err != nil {
		logger.Errorf("GeoMapDelete failed: %s", err.Error())
		return err
	}
	logger.Debugf("Deleting GeoMap: %v", existGeo)
	uStat, err := existGeo.Delete(domain)
	if err != nil {
		logger.Errorf("GeoMapDelete failed: %s", err.Error())
		return err
	}
	logger.Debugf("GeoMap Delete status:")
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
			logger.Infof("GeoMap delete completed")
		} else {
			if err == nil {
				logger.Infof("GeoMap delete pending")
			} else {
				logger.Warnf("GeoMap delete failed [%s]", err.Error())
				return err
			}
		}

	}

	// if successful ....
	d.SetId("")
	return nil
}

// Test GTM GeoMap existence
func resourceGTMv1GeomapExists(d *schema.ResourceData, m interface{}) (bool, error) {
	meta := akamai.Meta(m)
	logger := meta.Log("Akamai GTMv1", "resourceGTMv1GeomapExists")

	logger.Debugf("Exists")
	// pull domain and geoMap out of geoMap id
	domain, geoMap, err := parseResourceStringId(d.Id())
	if err != nil {
		return false, fmt.Errorf("invalid geoMap geoMap ID")
	}
	logger.Debugf("Searching for existing geoMap [%s] in domain %s", geoMap, domain)
	geo, err := gtm.GetGeoMap(geoMap, domain)
	return geo != nil, err
}

// Create and populate a new geoMap object from geoMap data
func populateNewGeoMapObject(d *schema.ResourceData, m interface{}) *gtm.GeoMap {

	name, _ := tools.GetStringValue("name", d)
	geoObj := gtm.NewGeoMap(name)
	geoObj.DefaultDatacenter = &gtm.DatacenterBase{}
	geoObj.Assignments = make([]*gtm.GeoAssignment, 1)
	geoObj.Links = make([]*gtm.Link, 1)
	populateGeoMapObject(d, geoObj, m)

	return geoObj
}

// Populate existing geoMap object from geoMap data
func populateGeoMapObject(d *schema.ResourceData, geo *gtm.GeoMap, m interface{}) {

	if v, err := tools.GetStringValue("name", d); err == nil {
		geo.Name = v
	}
	populateGeoAssignmentsObject(d, geo, m)
	populateGeoDefaultDCObject(d, geo, m)
}

// Populate Terraform state from provided GeoMap object
func populateTerraformGeoMapState(d *schema.ResourceData, geo *gtm.GeoMap, m interface{}) {
	meta := akamai.Meta(m)
	logger := meta.Log("Akamai GTMv1", "populateTerraformGeoMapState")

	// walk through all state elements
	if err := d.Set("name", geo.Name); err != nil {
		logger.Errorf("populateTerraformGeoMapState failed: %s", err.Error())
	}
	populateTerraformGeoAssignmentsState(d, geo, m)
	populateTerraformGeoDefaultDCState(d, geo, m)
}

// create and populate GTM GeoMap Assignments object
func populateGeoAssignmentsObject(d *schema.ResourceData, geo *gtm.GeoMap, m interface{}) {
	meta := akamai.Meta(m)
	logger := meta.Log("Akamai GTMv1", "resourceGTMv1GeomapExists")

	// pull apart List
	geoAssignmentsList, err := tools.GetInterfaceArrayValue("assignment", d)
	if err == nil {
		geoAssignmentsObjList := make([]*gtm.GeoAssignment, len(geoAssignmentsList)) // create new object list
		for i, v := range geoAssignmentsList {
			geoMap, ok := v.(map[string]interface{})
			if !ok {
				logger.Warnf("populateGeoAssignmentsObject failed, bad geoMap format: %s", v)
				continue
			}
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
func populateTerraformGeoAssignmentsState(d *schema.ResourceData, geo *gtm.GeoMap, m interface{}) {
	meta := akamai.Meta(m)
	logger := meta.Log("Akamai GTMv1", "populateTerraformGeoAssignmentsState")

	objectInventory := make(map[int]*gtm.GeoAssignment, len(geo.Assignments))
	if len(geo.Assignments) > 0 {
		for _, aObj := range geo.Assignments {
			objectInventory[aObj.DatacenterId] = aObj
		}
	}
	aStateList, _ := tools.GetInterfaceArrayValue("assignment", d)
	for _, aMap := range aStateList {
		a := aMap.(map[string]interface{})
		objIndex := a["datacenter_id"].(int)
		aObject := objectInventory[objIndex]
		if aObject == nil {
			logger.Warnf("Geo Assignment %d NOT FOUND in returned GTM Object", a["datacenter_id"])
			continue
		}
		a["datacenter_id"] = aObject.DatacenterId
		a["nickname"] = aObject.Nickname
		a["countries"] = reconcileTerraformLists(a["countries"].([]interface{}), convertStringToInterfaceList(aObject.Countries, m), m)
		// remove object
		delete(objectInventory, objIndex)
	}
	if len(objectInventory) > 0 {
		logger.Debugf("Geo Assignment objects left...")
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
	if err := d.Set("assignment", aStateList); err != nil {
		logger.Errorf("populateTerraformGeoAssignmentsState failed: %s", err.Error())
	}
}

// create and populate GTM GeoMap DefaultDatacenter object
func populateGeoDefaultDCObject(d *schema.ResourceData, geo *gtm.GeoMap, m interface{}) {
	meta := akamai.Meta(m)
	logger := meta.Log("Akamai GTMv1", "populateGeoDefaultDCObject")

	// pull apart List
	geoDefaultDCList, err := tools.GetInterfaceArrayValue("default_datacenter", d)
	if err == nil && len(geoDefaultDCList) > 0 {
		geoDefaultDCObj := gtm.DatacenterBase{} // create new object
		geoMap := geoDefaultDCList[0].(map[string]interface{})
		if geoMap["datacenter_id"] != nil && geoMap["datacenter_id"].(int) != 0 {
			geoDefaultDCObj.DatacenterId = geoMap["datacenter_id"].(int)
			geoDefaultDCObj.Nickname = geoMap["nickname"].(string)
		} else {
			logger.Infof("No Default Datacenter specified")
			var nilInt int
			geoDefaultDCObj.DatacenterId = nilInt
			geoDefaultDCObj.Nickname = ""
		}
		geo.DefaultDatacenter = &geoDefaultDCObj
	}
}

// create and populate Terraform geoMap default_datacenter schema
func populateTerraformGeoDefaultDCState(d *schema.ResourceData, geo *gtm.GeoMap, m interface{}) {
	meta := akamai.Meta(m)
	logger := meta.Log("Akamai GTMv1", "populateTerraformGeoDefault")

	ddcListNew := make([]interface{}, 1)
	ddcNew := map[string]interface{}{
		"datacenter_id": geo.DefaultDatacenter.DatacenterId,
		"nickname":      geo.DefaultDatacenter.Nickname,
	}
	ddcListNew[0] = ddcNew
	if err := d.Set("default_datacenter", ddcListNew); err != nil {
		logger.Errorf("populateTerraformGeoDefaultDCState failed: %s", err.Error())
	}
}
