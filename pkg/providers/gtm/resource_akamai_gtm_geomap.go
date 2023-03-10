package gtm

import (
	"context"
	"fmt"

	"github.com/akamai/AkamaiOPEN-edgegrid-golang/v4/pkg/gtm"
	"github.com/akamai/AkamaiOPEN-edgegrid-golang/v4/pkg/session"

	"github.com/akamai/terraform-provider-akamai/v3/pkg/akamai"
	"github.com/akamai/terraform-provider-akamai/v3/pkg/tools"
	"github.com/hashicorp/terraform-plugin-sdk/v2/diag"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/schema"

	"golang.org/x/exp/slices"
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
				Type:             schema.TypeList,
				Optional:         true,
				DiffSuppressFunc: assignmentDiffSuppress,
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
							Type:     schema.TypeSet,
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
	geoAssignmentsList, err := tools.GetListValue("assignment", d)
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
				countries, ok := geoMap["countries"].(*schema.Set)
				if !ok {
					logger.Warnf("wrong type conversion: expected *schema.Set, got %T", countries)
				}
				ls := make([]string, countries.Len())
				for i, sl := range countries.List() {
					ls[i] = sl.(string)
				}
				geoAssignment.Countries = ls
			}
			geoAssignmentsObjList[i] = &geoAssignment
		}
		geo.Assignments = geoAssignmentsObjList
	}
}

func setGeoAssignmentAtIndex(asArr *[]interface{}, as *gtm.GeoAssignment, index int, doReplace bool, m interface{}) {
	meta := akamai.Meta(m)
	logger := meta.Log("Akamai GTM", "setGeoAssignmentAtIndex")

	if !doReplace {
		newInt := make([]interface{}, 1)
		newInt[0] = make(map[string]interface{})
		if index == 0 {
			*asArr = append(newInt, *asArr...)
		} else {
			*asArr = append((*asArr)[:index], append(newInt, (*asArr)[index:]...)...)
		}
	}

	targetAs := (*asArr)[index].(map[string]interface{})
	targetAs["datacenter_id"] = as.DatacenterId
	targetAs["nickname"] = as.Nickname

	cts, ok := targetAs["countries"]
	if ok {
		countries, cOk := cts.(*schema.Set)
		if !cOk {
			logger.Warnf("wrong type conversion: expected *schema.Set, got %T", countries)
		}
		targetAs["countries"] = reconcileTerraformLists(countries.List(), convertStringToInterfaceList(as.Countries, m), m)
	} else {
		targetAs["countries"] = as.Countries
	}
}

// create and populate Terraform geoMap assignments schema
func populateTerraformGeoAssignmentsState(d *schema.ResourceData, geo *gtm.GeoMap, m interface{}) {
	meta := akamai.Meta(m)
	logger := meta.Log("Akamai GTM", "populateTerraformGeoAssignmentsState")

	aStateList, _ := tools.GetInterfaceArrayValue("assignment", d)
	stateIdx := 0

	for _, aObj := range geo.Assignments {
		idx := slices.IndexFunc(aStateList, func(as interface{}) bool {
			return as.(map[string]interface{})["datacenter_id"].(int) == aObj.DatacenterId
		})

		if idx == -1 {
			if stateIdx == 0 {
				logger.Debugf("assignment for dc %s not found into the state, prepending it", aObj.DatacenterId)
				setGeoAssignmentAtIndex(&aStateList, aObj, 0, false, m)
			} else {
				logger.Debugf("assignment for dc %s not found into the state, inserting it at position %s into the state", aObj.DatacenterId, stateIdx)
				setGeoAssignmentAtIndex(&aStateList, aObj, stateIdx, false, m)
			}
		} else {
			if idx < stateIdx {
				stateIdx--
				logger.Debugf("assignment for dc %s will now be placed after the position it was in the past (moving it from %s to %s)", aObj.DatacenterId, idx, stateIdx)
				aStateList = append(aStateList[:idx], aStateList[idx+1:]...)
				setGeoAssignmentAtIndex(&aStateList, aObj, stateIdx, false, m)
			} else {
				logger.Debugf("assignment for dc %s will be updated in-place into the state (index %s)", aObj.DatacenterId, idx)
				setGeoAssignmentAtIndex(&aStateList, aObj, idx, true, m)
				stateIdx = idx
			}
		}
		stateIdx++
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

// countriesEqual checks whether countries are equal
func countriesEqual(old, new interface{}) bool {
	logger := akamai.Log("Akamai GTM", "countriesEqual")

	oldCountries, ok := old.(*schema.Set)
	if !ok {
		logger.Warnf("wrong type conversion: expected *schema.Set, got %T", oldCountries)
		return false
	}

	newCountries, ok := new.(*schema.Set)
	if !ok {
		logger.Warnf("wrong type conversion: expected *schema.Set, got %T", newCountries)
		return false
	}

	if oldCountries.Len() != newCountries.Len() {
		return false
	}

	countries := make(map[string]bool, oldCountries.Len())
	for _, country := range oldCountries.List() {
		countries[country.(string)] = true
	}

	for _, country := range newCountries.List() {
		_, ok = countries[country.(string)]
		if !ok {
			return false
		}
	}

	return true
}
