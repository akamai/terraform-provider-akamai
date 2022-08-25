package gtm

import (
	"context"
	"fmt"

	gtm "github.com/akamai/AkamaiOPEN-edgegrid-golang/v2/pkg/configgtm"
	"github.com/akamai/AkamaiOPEN-edgegrid-golang/v2/pkg/session"

	"github.com/akamai/terraform-provider-akamai/v2/pkg/akamai"
	"github.com/akamai/terraform-provider-akamai/v2/pkg/tools"
	"github.com/hashicorp/terraform-plugin-sdk/v2/diag"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/schema"
)

func resourceGTMv1Geomap() *schema.Resource {
	return &schema.Resource{
		CreateContext: resourceGTMv1GeomapCreate,
		ReadContext:   resourceGTMv1GeomapRead,
		UpdateContext: resourceGTMv1GeomapUpdate,
		DeleteContext: resourceGTMv1GeomapDelete,
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
func resourceGTMv1GeomapCreate(ctx context.Context, d *schema.ResourceData, m interface{}) diag.Diagnostics {
	meta := akamai.Meta(m)
	logger := meta.Log("Akamai GTM", "resourceGTMv1GeomapCreate")
	// create a context with logging for api calls
	ctx = session.ContextWithOptions(
		ctx,
		session.WithContextLog(logger),
	)

	domain, err := tools.GetStringValue("domain", d)
	if err != nil {
		return diag.FromErr(err)
	}

	name, _ := tools.GetStringValue("name", d)
	if err != nil {
		return diag.FromErr(err)
	}

	logger.Infof("[Akamai GTM] Creating geoMap [%s] in domain [%s]", name, domain)
	// Make sure Default Datacenter exists
	geoDefaultDCList, err := tools.GetInterfaceArrayValue("default_datacenter", d)
	if err != nil {
		logger.Errorf("Default datacenter not initialized: %s", err.Error())
		return diag.FromErr(err)
	}
	var diags diag.Diagnostics
	if err := validateDefaultDC(ctx, meta, geoDefaultDCList, domain); err != nil {
		logger.Errorf("Default datacenter validation error: %s", err.Error())
		return append(diags, diag.Diagnostic{
			Severity: diag.Error,
			Summary:  "Default datacenter validation error",
			Detail:   err.Error(),
		})
	}

	newGeo := populateNewGeoMapObject(ctx, meta, d, m)
	logger.Debugf("Proposed New geoMap: [%v]", newGeo)
	cStatus, err := inst.Client(meta).CreateGeoMap(ctx, newGeo, domain)
	if err != nil {
		logger.Errorf("geoMap Create failed: %s", err.Error())
		return append(diags, diag.Diagnostic{
			Severity: diag.Error,
			Summary:  "geoMap Create failed",
			Detail:   err.Error(),
		})
	}
	logger.Debugf("geoMap Create status: %v", cStatus.Status)
	if cStatus.Status.PropagationStatus == "DENIED" {
		return append(diags, diag.Diagnostic{
			Severity: diag.Error,
			Summary:  cStatus.Status.Message,
		})
	}

	waitOnComplete, err := tools.GetBoolValue("wait_on_complete", d)
	if err != nil {
		return diag.FromErr(err)
	}

	if waitOnComplete {
		done, err := waitForCompletion(ctx, domain, m)
		if done {
			logger.Infof("geoMap Create completed")
		} else {
			if err == nil {
				logger.Infof("geoMap Create pending")
			} else {
				logger.Errorf("geoMap Create failed [%s]", err.Error())
				return append(diags, diag.Diagnostic{
					Severity: diag.Error,
					Summary:  "geoMap Create failed",
					Detail:   err.Error(),
				})
			}
		}

	}

	// Give terraform the ID. Format domain:geoMap
	geoMapID := fmt.Sprintf("%s:%s", domain, cStatus.Resource.Name)
	logger.Debugf("Generated geoMap resource ID: %s", geoMapID)
	d.SetId(geoMapID)
	return resourceGTMv1GeomapRead(ctx, d, m)

}

// read geoMap. updates state with entire API result configuration.
func resourceGTMv1GeomapRead(ctx context.Context, d *schema.ResourceData, m interface{}) diag.Diagnostics {
	meta := akamai.Meta(m)
	logger := meta.Log("Akamai GTM", "resourceGTMv1GeomapRead")
	// create a context with logging for api calls
	ctx = session.ContextWithOptions(
		ctx,
		session.WithContextLog(logger),
	)

	logger.Debugf("Reading geoMap: %s", d.Id())
	var diags diag.Diagnostics
	// retrieve the property and domain
	domain, geoMap, err := parseResourceStringID(d.Id())
	if err != nil {
		logger.Errorf("Invalid geoMap ID")
		return diag.FromErr(err)
	}
	geo, err := inst.Client(meta).GetGeoMap(ctx, geoMap, domain)
	if err != nil {
		logger.Errorf("geoMap Read error: %s", err.Error())
		return append(diags, diag.Diagnostic{
			Severity: diag.Error,
			Summary:  "geoMap Read error",
			Detail:   err.Error(),
		})
	}
	populateTerraformGeoMapState(d, geo, m)
	logger.Debugf("READ %v", geo)
	return nil
}

// Update GTM GeoMap
func resourceGTMv1GeomapUpdate(ctx context.Context, d *schema.ResourceData, m interface{}) diag.Diagnostics {
	meta := akamai.Meta(m)
	logger := meta.Log("Akamai GTM", "resourceGTMv1GeomapUpdate")
	// create a context with logging for api calls
	ctx = session.ContextWithOptions(
		ctx,
		session.WithContextLog(logger),
	)

	logger.Debugf("Updating geoMap: %s", d.Id())
	var diags diag.Diagnostics
	// pull domain and geoMap out of id
	domain, geoMap, err := parseResourceStringID(d.Id())
	if err != nil {
		logger.Errorf("Invalid geoMap ID")
		return diag.FromErr(err)
	}
	// Get existingGeoMap
	existGeo, err := inst.Client(meta).GetGeoMap(ctx, geoMap, domain)
	if err != nil {
		logger.Errorf("geoMap Update failed: %s", err.Error())
		return append(diags, diag.Diagnostic{
			Severity: diag.Error,
			Summary:  "geoMap Update Read error",
			Detail:   err.Error(),
		})
	}
	logger.Debugf("Updating geoMap BEFORE: %v", existGeo)
	populateGeoMapObject(d, existGeo, m)
	logger.Debugf("Updating geoMap PROPOSED: %v", existGeo)
	uStat, err := inst.Client(meta).UpdateGeoMap(ctx, existGeo, domain)
	if err != nil {
		logger.Errorf("geoMap Update failed: %s", err.Error())
		return append(diags, diag.Diagnostic{
			Severity: diag.Error,
			Summary:  "geoMap Update error",
			Detail:   err.Error(),
		})
	}
	logger.Debugf("geoMap Update  status: %v", uStat)
	if uStat.PropagationStatus == "DENIED" {
		logger.Errorf(uStat.Message)
		return append(diags, diag.Diagnostic{
			Severity: diag.Error,
			Summary:  uStat.Message,
		})
	}

	waitOnComplete, err := tools.GetBoolValue("wait_on_complete", d)
	if err != nil {
		return diag.FromErr(err)
	}

	if waitOnComplete {
		done, err := waitForCompletion(ctx, domain, m)
		if done {
			logger.Infof("geoMap Update completed")
		} else {
			if err == nil {
				logger.Infof("geoMap Update pending")
			} else {
				logger.Errorf("geoMap Update failed [%s]", err.Error())
				return append(diags, diag.Diagnostic{
					Severity: diag.Error,
					Summary:  "geoMap Update failed",
					Detail:   err.Error(),
				})
			}
		}

	}

	return resourceGTMv1GeomapRead(ctx, d, m)
}

// Import GTM GeoMap.
func resourceGTMv1GeomapImport(d *schema.ResourceData, m interface{}) ([]*schema.ResourceData, error) {
	meta := akamai.Meta(m)
	logger := meta.Log("Akamai GTM", "resourceGTMv1GeomapImport")
	// create a context with logging for api calls
	ctx := context.Background()
	ctx = session.ContextWithOptions(
		ctx,
		session.WithContextLog(logger),
	)

	logger.Infof("[Akamai GTM] geoMap [%s] Import", d.Id())
	// pull domain and geoMap out of geoMap id
	domain, geoMap, err := parseResourceStringID(d.Id())
	if err != nil {
		return []*schema.ResourceData{d}, err
	}
	geo, err := inst.Client(meta).GetGeoMap(ctx, geoMap, domain)
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
	logger.Infof("[Akamai GTM] geoMap [%s] [%s] Imported", d.Id(), name)
	return []*schema.ResourceData{d}, nil
}

// Delete GTM GeoMap.
func resourceGTMv1GeomapDelete(ctx context.Context, d *schema.ResourceData, m interface{}) diag.Diagnostics {
	meta := akamai.Meta(m)
	logger := meta.Log("Akamai GTM", "resourceGTMv1GeomapDelete")
	// create a context with logging for api calls
	ctx = session.ContextWithOptions(
		ctx,
		session.WithContextLog(logger),
	)

	logger.Debugf("Deleting geoMap: %s", d.Id())
	var diags diag.Diagnostics
	// Get existing geoMap
	domain, geoMap, err := parseResourceStringID(d.Id())
	if err != nil {
		logger.Errorf("Invalid geoMap ID: %s", d.Id())
		return diag.FromErr(err)
	}
	existGeo, err := inst.Client(meta).GetGeoMap(ctx, geoMap, domain)
	if err != nil {
		logger.Errorf("geoMap Delete failed: %s", err.Error())
		return append(diags, diag.Diagnostic{
			Severity: diag.Error,
			Summary:  "geoMap Delete Read error",
			Detail:   err.Error(),
		})
	}
	logger.Debugf("Deleting geoMap: %v", existGeo)
	uStat, err := inst.Client(meta).DeleteGeoMap(ctx, existGeo, domain)
	if err != nil {
		logger.Errorf("geoMap Delete failed: %s", err.Error())
		return append(diags, diag.Diagnostic{
			Severity: diag.Error,
			Summary:  "geoMap Delete error",
			Detail:   err.Error(),
		})
	}
	logger.Debugf("geoMap Delete status: %v", uStat)
	if uStat.PropagationStatus == "DENIED" {
		logger.Errorf(uStat.Message)
		return append(diags, diag.Diagnostic{
			Severity: diag.Error,
			Summary:  uStat.Message,
		})
	}

	waitOnComplete, err := tools.GetBoolValue("wait_on_complete", d)
	if err != nil {
		return diag.FromErr(err)
	}

	if waitOnComplete {
		done, err := waitForCompletion(ctx, domain, m)
		if done {
			logger.Infof("geoMap Delete completed")
		} else {
			if err == nil {
				logger.Infof("geoMap Delete pending")
			} else {
				logger.Errorf("geoMap Delete failed [%s]", err.Error())
				return append(diags, diag.Diagnostic{
					Severity: diag.Error,
					Summary:  "geoMap Delete failed",
					Detail:   err.Error(),
				})
			}
		}
	}

	// if successful ....
	d.SetId("")
	return nil
}

// Create and populate a new geoMap object from geoMap data
func populateNewGeoMapObject(ctx context.Context, meta akamai.OperationMeta, d *schema.ResourceData, m interface{}) *gtm.GeoMap {

	name, _ := tools.GetStringValue("name", d)
	geoObj := inst.Client(meta).NewGeoMap(ctx, name)
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
	logger := meta.Log("Akamai GTM", "populateTerraformGeoMapState")

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
	logger := meta.Log("Akamai GTM", "populateGeoAssignmentsObject")

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
	logger := meta.Log("Akamai GTM", "populateTerraformGeoAssignmentsState")

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
	logger := meta.Log("Akamai GTM", "populateGeoDefaultDCObject")

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
	logger := meta.Log("Akamai GTM", "populateTerraformGeoDefault")

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
