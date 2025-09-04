package gtm

import (
	"context"
	"errors"
	"fmt"

	"github.com/akamai/AkamaiOPEN-edgegrid-golang/v12/pkg/gtm"
	"github.com/akamai/AkamaiOPEN-edgegrid-golang/v12/pkg/session"
	"github.com/akamai/terraform-provider-akamai/v9/pkg/common/tf"
	"github.com/akamai/terraform-provider-akamai/v9/pkg/log"
	"github.com/akamai/terraform-provider-akamai/v9/pkg/meta"
	"github.com/hashicorp/terraform-plugin-sdk/v2/diag"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/schema"
)

const geoMapAlreadyExistsError = "GeoMap with provided `name` for specific `domain` already exists. Please import specific geomap using following command: terraform import akamai_gtm_geomap.<your_resource_name> \"%s:%s\""

func resourceGTMv1GeoMap() *schema.Resource {
	return &schema.Resource{
		CreateContext: resourceGTMv1GeoMapCreate,
		ReadContext:   resourceGTMv1GeoMapRead,
		UpdateContext: resourceGTMv1GeoMapUpdate,
		DeleteContext: resourceGTMv1GeoMapDelete,
		Importer: &schema.ResourceImporter{
			State: resourceGTMv1GeoMapImport,
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

func resourceGTMv1GeoMapCreate(ctx context.Context, d *schema.ResourceData, m interface{}) diag.Diagnostics {
	meta := meta.Must(m)
	logger := meta.Log("Akamai GTM", "resourceGTMv1GeoMapCreate")
	// create a context with logging for api calls
	ctx = session.ContextWithOptions(
		ctx,
		session.WithContextLog(logger),
	)
	var diags diag.Diagnostics

	domain, err := tf.GetStringValue("domain", d)
	if err != nil {
		return diag.FromErr(err)
	}

	name, err := tf.GetStringValue("name", d)
	if err != nil {
		return diag.FromErr(err)
	}

	geo, err := Client(meta).GetGeoMap(ctx, gtm.GetGeoMapRequest{
		MapName:    name,
		DomainName: domain,
	})
	if err != nil && !errors.Is(err, gtm.ErrNotFound) {
		logger.Errorf("geoMap read error: %s", err.Error())
		return append(diags, diag.Diagnostic{
			Severity: diag.Error,
			Summary:  "geoMap read error",
			Detail:   err.Error(),
		})
	}
	if geo != nil {
		geoMapAlreadyExists := fmt.Sprintf(geoMapAlreadyExistsError, domain, name)
		logger.Errorf(geoMapAlreadyExists)
		return append(diags, diag.Diagnostic{
			Severity: diag.Error,
			Summary:  "geoMap already exists error",
			Detail:   geoMapAlreadyExists,
		})
	}

	logger.Infof("[Akamai GTM] Creating geoMap [%s] in domain [%s]", name, domain)
	// Make sure Default Datacenter exists
	geoDefaultDCList, err := tf.GetInterfaceArrayValue("default_datacenter", d)
	if err != nil {
		logger.Errorf("Default datacenter not initialized: %s", err.Error())
		return diag.FromErr(err)
	}
	if err := validateDefaultDC(ctx, meta, geoDefaultDCList, domain); err != nil {
		logger.Errorf("Default datacenter validation error: %s", err.Error())
		return append(diags, diag.Diagnostic{
			Severity: diag.Error,
			Summary:  "Default datacenter validation error",
			Detail:   err.Error(),
		})
	}

	newGeo, err := populateNewGeoMapObject(d, m)
	if err != nil {
		logger.Errorf("geoMap populate error: %s", err.Error())
		return append(diags, diag.Diagnostic{
			Severity: diag.Error,
			Summary:  "geoMap populate error",
			Detail:   err.Error(),
		})
	}
	logger.Debugf("Proposed New geoMap: [%v]", newGeo)
	cStatus, err := Client(meta).CreateGeoMap(ctx, gtm.CreateGeoMapRequest{
		GeoMap:     newGeo,
		DomainName: domain,
	})
	if err != nil {
		logger.Errorf("geoMap create error: %s", err.Error())
		return append(diags, diag.Diagnostic{
			Severity: diag.Error,
			Summary:  "geoMap create error",
			Detail:   err.Error(),
		})
	}
	logger.Debugf("geoMap create status: %v", cStatus.Status)
	if cStatus.Status.PropagationStatus == "DENIED" {
		return append(diags, diag.Diagnostic{
			Severity: diag.Error,
			Summary:  cStatus.Status.Message,
		})
	}

	waitOnComplete, err := tf.GetBoolValue("wait_on_complete", d)
	if err != nil {
		return diag.FromErr(err)
	}

	if waitOnComplete {
		done, err := waitForCompletion(ctx, domain, m)
		if done {
			logger.Infof("geoMap create completed")
		} else {
			if err == nil {
				logger.Infof("geoMap create pending")
			} else {
				logger.Errorf("geoMap create error: %s", err.Error())
				return append(diags, diag.Diagnostic{
					Severity: diag.Error,
					Summary:  "geoMap create error",
					Detail:   err.Error(),
				})
			}
		}

	}

	// Give terraform the ID. Format domain:geoMap
	geoMapID := fmt.Sprintf("%s:%s", domain, cStatus.Resource.Name)
	logger.Debugf("Generated geoMap resource ID: %s", geoMapID)
	d.SetId(geoMapID)
	return resourceGTMv1GeoMapRead(ctx, d, m)
}

func resourceGTMv1GeoMapRead(ctx context.Context, d *schema.ResourceData, m interface{}) diag.Diagnostics {
	meta := meta.Must(m)
	logger := meta.Log("Akamai GTM", "resourceGTMv1GeoMapRead")
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
	geo, err := Client(meta).GetGeoMap(ctx, gtm.GetGeoMapRequest{
		MapName:    geoMap,
		DomainName: domain,
	})
	if errors.Is(err, gtm.ErrNotFound) {
		d.SetId("")
		return nil
	}
	if err != nil {
		logger.Errorf("geoMap read error: %s", err.Error())
		return append(diags, diag.Diagnostic{
			Severity: diag.Error,
			Summary:  "geoMap read error",
			Detail:   err.Error(),
		})
	}
	if err = populateTerraformGeoMapState(d, geo, m); err != nil {
		return append(diags, diag.Diagnostic{
			Severity: diag.Error,
			Summary:  "geoMap read populate state error",
			Detail:   err.Error(),
		})
	}
	logger.Debugf("READ %v", geo)
	return nil
}

func resourceGTMv1GeoMapUpdate(ctx context.Context, d *schema.ResourceData, m interface{}) diag.Diagnostics {
	meta := meta.Must(m)
	logger := meta.Log("Akamai GTM", "resourceGTMv1GeoMapUpdate")
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
	existGeo, err := Client(meta).GetGeoMap(ctx, gtm.GetGeoMapRequest{
		MapName:    geoMap,
		DomainName: domain,
	})
	if err != nil {
		logger.Errorf("geoMap read error: %s", err.Error())
		return append(diags, diag.Diagnostic{
			Severity: diag.Error,
			Summary:  "geoMap read error",
			Detail:   err.Error(),
		})
	}
	newGeo := createGeoMapStruct(existGeo)
	logger.Debugf("Updating geoMap BEFORE: %v", newGeo)
	populateGeoMapObject(d, newGeo, m)
	logger.Debugf("Updating geoMap PROPOSED: %v", existGeo)
	uStat, err := Client(meta).UpdateGeoMap(ctx, gtm.UpdateGeoMapRequest{
		GeoMap:     newGeo,
		DomainName: domain,
	})
	if err != nil {
		logger.Errorf("geoMap update error: %s", err.Error())
		return append(diags, diag.Diagnostic{
			Severity: diag.Error,
			Summary:  "geoMap update error",
			Detail:   err.Error(),
		})
	}
	logger.Debugf("geoMap update status: %v", uStat)
	if uStat.Status.PropagationStatus == "DENIED" {
		logger.Errorf(uStat.Status.Message)
		return append(diags, diag.Diagnostic{
			Severity: diag.Error,
			Summary:  uStat.Status.Message,
		})
	}

	waitOnComplete, err := tf.GetBoolValue("wait_on_complete", d)
	if err != nil {
		return diag.FromErr(err)
	}

	if waitOnComplete {
		done, err := waitForCompletion(ctx, domain, m)
		if done {
			logger.Infof("geoMap update completed")
		} else {
			if err == nil {
				logger.Infof("geoMap update pending")
			} else {
				logger.Errorf("geoMap update error: %s", err.Error())
				return append(diags, diag.Diagnostic{
					Severity: diag.Error,
					Summary:  "geoMap update error",
					Detail:   err.Error(),
				})
			}
		}

	}

	return resourceGTMv1GeoMapRead(ctx, d, m)
}

func resourceGTMv1GeoMapImport(d *schema.ResourceData, m interface{}) ([]*schema.ResourceData, error) {
	meta := meta.Must(m)
	logger := meta.Log("Akamai GTM", "resourceGTMv1GeoMapImport")
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
	geo, err := Client(meta).GetGeoMap(ctx, gtm.GetGeoMapRequest{
		MapName:    geoMap,
		DomainName: domain,
	})
	if err != nil {
		return nil, err
	}
	if err := d.Set("domain", domain); err != nil {
		return nil, err
	}
	if err := d.Set("wait_on_complete", true); err != nil {
		return nil, err
	}
	if err = populateTerraformGeoMapState(d, geo, m); err != nil {
		return nil, err
	}

	// use same Id as passed in
	name, err := tf.GetStringValue("name", d)
	if err != nil {
		return []*schema.ResourceData{d}, err
	}
	logger.Infof("[Akamai GTM] geoMap [%s] [%s] Imported", d.Id(), name)
	return []*schema.ResourceData{d}, nil
}

func resourceGTMv1GeoMapDelete(ctx context.Context, d *schema.ResourceData, m interface{}) diag.Diagnostics {
	meta := meta.Must(m)
	logger := meta.Log("Akamai GTM", "resourceGTMv1GeoMapDelete")
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
	existGeo, err := Client(meta).GetGeoMap(ctx, gtm.GetGeoMapRequest{
		MapName:    geoMap,
		DomainName: domain,
	})
	if err != nil {
		logger.Errorf("geoMap read error: %s", err.Error())
		return append(diags, diag.Diagnostic{
			Severity: diag.Error,
			Summary:  "geoMap read error",
			Detail:   err.Error(),
		})
	}
	newGeo := createGeoMapStruct(existGeo)
	logger.Debugf("Deleting geoMap: %v", newGeo)
	uStat, err := Client(meta).DeleteGeoMap(ctx, gtm.DeleteGeoMapRequest{
		MapName:    geoMap,
		DomainName: domain,
	})
	if err != nil {
		logger.Errorf("geoMap delete error: %s", err.Error())
		return append(diags, diag.Diagnostic{
			Severity: diag.Error,
			Summary:  "geoMap delete error",
			Detail:   err.Error(),
		})
	}
	logger.Debugf("geoMap delete status: %v", uStat)
	if uStat.Status.PropagationStatus == "DENIED" {
		logger.Errorf(uStat.Status.Message)
		return append(diags, diag.Diagnostic{
			Severity: diag.Error,
			Summary:  uStat.Status.Message,
		})
	}

	waitOnComplete, err := tf.GetBoolValue("wait_on_complete", d)
	if err != nil {
		return diag.FromErr(err)
	}

	if waitOnComplete {
		done, err := waitForCompletion(ctx, domain, m)
		if done {
			logger.Infof("geoMap delete completed")
		} else {
			if err == nil {
				logger.Infof("geoMap delete pending")
			} else {
				logger.Errorf("geoMap delete error: %s", err.Error())
				return append(diags, diag.Diagnostic{
					Severity: diag.Error,
					Summary:  "geoMap delete error",
					Detail:   err.Error(),
				})
			}
		}
	}

	d.SetId("")
	return nil
}

// Create and populate a new geoMap object from geoMap data
func populateNewGeoMapObject(d *schema.ResourceData, m interface{}) (*gtm.GeoMap, error) {

	name, err := tf.GetStringValue("name", d)
	if err != nil {
		return nil, err
	}
	geoObj := &gtm.GeoMap{
		Name:              name,
		DefaultDatacenter: &gtm.DatacenterBase{},
	}
	populateGeoMapObject(d, geoObj, m)

	return geoObj, nil
}

// Populate existing geoMap object from geoMap data
func populateGeoMapObject(d *schema.ResourceData, geo *gtm.GeoMap, m interface{}) {
	if v, err := tf.GetStringValue("name", d); err == nil {
		geo.Name = v
	}
	populateGeoAssignmentsObject(d, geo, m)
	populateGeoDefaultDCObject(d, geo, m)
}

// Populate Terraform state from provided GeoMap object
func populateTerraformGeoMapState(d *schema.ResourceData, geo *gtm.GetGeoMapResponse, m interface{}) error {
	meta := meta.Must(m)
	logger := meta.Log("Akamai GTM", "populateTerraformGeoMapState")

	// walk through all state elements
	if err := d.Set("name", geo.Name); err != nil {
		logger.Errorf("populateTerraformGeoMapState failed: %s", err.Error())
	}
	if err := populateTerraformGeoAssignmentsState(d, geo, m); err != nil {
		return err
	}
	return populateTerraformGeoDefaultDCState(d, geo, m)
}

// create and populate GTM GeoMap Assignments object
func populateGeoAssignmentsObject(d *schema.ResourceData, geo *gtm.GeoMap, m interface{}) {
	meta := meta.Must(m)
	logger := meta.Log("Akamai GTM", "populateGeoAssignmentsObject")

	// pull apart List
	geoAssignmentsList, err := tf.GetListValue("assignment", d)
	if err == nil {
		geoAssignmentsObjList := make([]gtm.GeoAssignment, len(geoAssignmentsList)) // create new object list
		for i, v := range geoAssignmentsList {
			geoMap, ok := v.(map[string]interface{})
			if !ok {
				logger.Warnf("populateGeoAssignmentsObject failed, bad geoMap format: %s", v)
				continue
			}
			geoAssignment := gtm.GeoAssignment{}
			geoAssignment.DatacenterID = geoMap["datacenter_id"].(int)
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
			geoAssignmentsObjList[i] = geoAssignment
		}
		geo.Assignments = geoAssignmentsObjList
	}
}

// create and populate Terraform geoMap assignments schema
func populateTerraformGeoAssignmentsState(d *schema.ResourceData, geo *gtm.GetGeoMapResponse, m interface{}) error {
	meta := meta.Must(m)
	logger := meta.Log("Akamai GTM", "populateTerraformGeoAssignmentsState")

	objectInventory := make(map[int]gtm.GeoAssignment, len(geo.Assignments))
	if len(geo.Assignments) > 0 {
		for _, aObj := range geo.Assignments {
			objectInventory[aObj.DatacenterID] = aObj
		}
	}
	aStateList, err := tf.GetInterfaceArrayValue("assignment", d)
	if err != nil && !errors.Is(err, tf.ErrNotFound) {
		return err
	}
	for _, aMap := range aStateList {
		a := aMap.(map[string]interface{})
		objIndex := a["datacenter_id"].(int)
		aObject, ok := objectInventory[objIndex]
		if !ok {
			logger.Warnf("Geo Assignment %d NOT FOUND in returned GTM Object", a["datacenter_id"])
			continue
		}
		a["datacenter_id"] = aObject.DatacenterID
		a["nickname"] = aObject.Nickname
		countries, ok := a["countries"].(*schema.Set)
		if !ok {
			logger.Warnf("wrong type conversion: expected *schema.Set, got %T", countries)
		}
		a["countries"] = reconcileTerraformLists(countries.List(), convertStringToInterfaceList(aObject.Countries, m), m)
		// remove object
		delete(objectInventory, objIndex)
	}
	if len(objectInventory) > 0 {
		logger.Debugf("Geo Assignment objects left...")
		// Objects not in the state yet. Add. Unfortunately, they not align with instance indices in the config
		for _, maObj := range objectInventory {
			aNew := map[string]interface{}{
				"datacenter_id": maObj.DatacenterID,
				"nickname":      maObj.Nickname,
				"countries":     maObj.Countries,
			}
			aStateList = append(aStateList, aNew)
		}
	}
	if err := d.Set("assignment", aStateList); err != nil {
		logger.Errorf("populateTerraformGeoAssignmentsState failed: %s", err.Error())
		return err
	}
	return nil
}

// create and populate GTM GeoMap DefaultDatacenter object
func populateGeoDefaultDCObject(d *schema.ResourceData, geo *gtm.GeoMap, m interface{}) {
	meta := meta.Must(m)
	logger := meta.Log("Akamai GTM", "populateGeoDefaultDCObject")

	// pull apart List
	geoDefaultDCList, err := tf.GetInterfaceArrayValue("default_datacenter", d)
	if err == nil && len(geoDefaultDCList) > 0 {
		geoDefaultDCObj := gtm.DatacenterBase{} // create new object
		geoMap := geoDefaultDCList[0].(map[string]interface{})
		if geoMap["datacenter_id"] != nil && geoMap["datacenter_id"].(int) != 0 {
			geoDefaultDCObj.DatacenterID = geoMap["datacenter_id"].(int)
			geoDefaultDCObj.Nickname = geoMap["nickname"].(string)
		} else {
			logger.Infof("No Default Datacenter specified")
			var nilInt int
			geoDefaultDCObj.DatacenterID = nilInt
			geoDefaultDCObj.Nickname = ""
		}
		geo.DefaultDatacenter = &geoDefaultDCObj
	}
}

// create and populate Terraform geoMap default_datacenter schema
func populateTerraformGeoDefaultDCState(d *schema.ResourceData, geo *gtm.GetGeoMapResponse, m interface{}) error {
	meta := meta.Must(m)
	logger := meta.Log("Akamai GTM", "populateTerraformGeoDefault")

	ddcListNew := make([]interface{}, 1)
	ddcNew := map[string]interface{}{
		"datacenter_id": geo.DefaultDatacenter.DatacenterID,
		"nickname":      geo.DefaultDatacenter.Nickname,
	}
	ddcListNew[0] = ddcNew
	if err := d.Set("default_datacenter", ddcListNew); err != nil {
		logger.Errorf("populateTerraformGeoDefaultDCState failed: %s", err.Error())
		return err
	}
	return nil
}

// createGeoMapStruct converts response from GetGeoMapResponse into GeoMap
func createGeoMapStruct(geo *gtm.GetGeoMapResponse) *gtm.GeoMap {
	if geo != nil {
		return &gtm.GeoMap{
			DefaultDatacenter: geo.DefaultDatacenter,
			Assignments:       geo.Assignments,
			Name:              geo.Name,
			Links:             geo.Links,
		}
	}
	return nil
}

// countriesEqual checks whether countries are equal
func countriesEqual(o, n interface{}) bool {
	logger := log.Get("Akamai GTM", "countriesEqual")

	oldCountries, ok := o.(*schema.Set)
	if !ok {
		logger.Warnf("wrong type conversion: expected *schema.Set, got %T", oldCountries)
		return false
	}

	newCountries, ok := n.(*schema.Set)
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
