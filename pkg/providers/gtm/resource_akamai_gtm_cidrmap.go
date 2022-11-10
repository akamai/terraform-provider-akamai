package gtm

import (
	"context"
	"fmt"

	"github.com/akamai/AkamaiOPEN-edgegrid-golang/v3/pkg/gtm"
	"github.com/akamai/AkamaiOPEN-edgegrid-golang/v3/pkg/session"

	"github.com/akamai/terraform-provider-akamai/v3/pkg/akamai"
	"github.com/akamai/terraform-provider-akamai/v3/pkg/tools"
	"github.com/hashicorp/terraform-plugin-sdk/v2/diag"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/schema"
)

func resourceGTMv1Cidrmap() *schema.Resource {
	return &schema.Resource{
		CreateContext: resourceGTMv1CidrMapCreate,
		ReadContext:   resourceGTMv1CidrMapRead,
		UpdateContext: resourceGTMv1CidrMapUpdate,
		DeleteContext: resourceGTMv1CidrMapDelete,
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
func resourceGTMv1CidrMapCreate(ctx context.Context, d *schema.ResourceData, m interface{}) diag.Diagnostics {
	meta := akamai.Meta(m)
	logger := meta.Log("Akamai GTM", "resourceGTMCidrMapCreate")
	// create a context with logging for api calls
	ctx = session.ContextWithOptions(
		ctx,
		session.WithContextLog(logger),
	)

	domain, err := tools.GetStringValue("domain", d)
	if err != nil {
		return diag.FromErr(err)
	}

	name, err := tools.GetStringValue("name", d)
	if err != nil {
		return diag.FromErr(err)
	}

	logger.Infof("Creating cidrMap [%s] in domain [%s]", name, domain)
	var diags diag.Diagnostics
	// Make sure Default Datacenter exists
	defaultDatacenter, err := tools.GetInterfaceArrayValue("default_datacenter", d)
	if err != nil {
		return diag.FromErr(err)
	}
	if err = validateDefaultDC(ctx, meta, defaultDatacenter, domain); err != nil {
		logger.Errorf("Default datacenter validation error: %s", err.Error())
		return append(diags, diag.Diagnostic{
			Severity: diag.Error,
			Summary:  "Default datacenter validation error",
			Detail:   err.Error(),
		})
	}

	newCidr := populateNewCidrMapObject(ctx, meta, d, m)
	logger.Debugf("Proposed New CidrMap: [%v]", newCidr)
	cStatus, err := inst.Client(meta).CreateCidrMap(ctx, newCidr, domain)
	if err != nil {
		logger.Errorf("cidrMap Create failed: %s", err.Error())
		return append(diags, diag.Diagnostic{
			Severity: diag.Error,
			Summary:  "cidrMap Create failed",
			Detail:   err.Error(),
		})
	}
	logger.Debugf("cidrMap Create status: %v", cStatus.Status)
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
			logger.Infof("cidrMap Create completed")
		} else {
			if err == nil {
				logger.Infof("cidrMap Create pending")
			} else {
				logger.Errorf("cidrMap Create failed [%s]", err.Error())
				return append(diags, diag.Diagnostic{
					Severity: diag.Error,
					Summary:  "cidrMap Create failed",
					Detail:   err.Error(),
				})
			}
		}
	}

	// Give terraform the ID. Format domain:cidrMap
	cidrMapID := fmt.Sprintf("%s:%s", domain, cStatus.Resource.Name)
	logger.Debugf("Generated cidrMap resource Id: %s", cidrMapID)
	d.SetId(cidrMapID)
	return resourceGTMv1CidrMapRead(ctx, d, m)

}

// read cidrMap. updates state with entire API result configuration.
func resourceGTMv1CidrMapRead(ctx context.Context, d *schema.ResourceData, m interface{}) diag.Diagnostics {
	meta := akamai.Meta(m)
	logger := meta.Log("Akamai GTM", "resourceGTMCidrMapRead")
	// create a context with logging for api calls
	ctx = session.ContextWithOptions(
		ctx,
		session.WithContextLog(logger),
	)

	logger.Debugf("Reading cidrMap: %s", d.Id())
	var diags diag.Diagnostics
	// retrieve the property and domain
	domain, cidrMap, err := parseResourceStringID(d.Id())
	if err != nil {
		logger.Errorf("Invalid cidrMap ID: %s", d.Id())
		return diag.FromErr(err)
	}
	cidr, err := inst.Client(meta).GetCidrMap(ctx, cidrMap, domain)
	if err != nil {
		logger.Errorf("cidrMap Read error: %s", err.Error())
		return append(diags, diag.Diagnostic{
			Severity: diag.Error,
			Summary:  "cidrMap Read error",
			Detail:   err.Error(),
		})
	}
	populateTerraformCidrMapState(d, cidr, m)
	logger.Debugf("READ %v", cidr)
	return nil
}

// Update GTM CidrMap
func resourceGTMv1CidrMapUpdate(ctx context.Context, d *schema.ResourceData, m interface{}) diag.Diagnostics {
	meta := akamai.Meta(m)
	logger := meta.Log("Akamai GTM", "resourceGTMCidrMapUpdate")
	// create a context with logging for api calls
	ctx = session.ContextWithOptions(
		ctx,
		session.WithContextLog(logger),
	)

	logger.Debugf("Updating cidrMap: %s", d.Id())
	var diags diag.Diagnostics
	// pull domain and cidrMap out of id
	domain, cidrMap, err := parseResourceStringID(d.Id())
	if err != nil {
		logger.Errorf("Invalid cidrMap ID: %s", d.Id())
		return diag.FromErr(err)
	}
	// Get existingCidrMap
	existCidr, err := inst.Client(meta).GetCidrMap(ctx, cidrMap, domain)
	if err != nil {
		logger.Errorf("cidrMap Update read failed: %s", err.Error())
		return append(diags, diag.Diagnostic{
			Severity: diag.Error,
			Summary:  "cidrMap Update Read error",
			Detail:   err.Error(),
		})
	}
	logger.Debugf("Updating cidrMap BEFORE: %v", existCidr)
	populateCidrMapObject(d, existCidr, m)
	logger.Debugf("Updating cidrMap PROPOSED: %v", existCidr)
	uStat, err := inst.Client(meta).UpdateCidrMap(ctx, existCidr, domain)
	if err != nil {
		logger.Errorf("cidrMap Update failed: %s", err.Error())
		return append(diags, diag.Diagnostic{
			Severity: diag.Error,
			Summary:  "cidrMap Update error",
			Detail:   err.Error(),
		})
	}
	logger.Debugf("CidrMap Update  status: %v", uStat)
	if uStat.PropagationStatus == "DENIED" {
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
			logger.Infof("cidrMap Update completed")
		} else {
			if err == nil {
				logger.Infof("cidrMap Update pending")
			} else {
				logger.Errorf("cidrMap Update failed [%s]", err.Error())
				return append(diags, diag.Diagnostic{
					Severity: diag.Error,
					Summary:  "cidrMap Update failed",
					Detail:   err.Error(),
				})
			}
		}
	}

	return resourceGTMv1CidrMapRead(ctx, d, m)
}

// Import GTM CidrMap.
func resourceGTMv1CidrMapImport(d *schema.ResourceData, m interface{}) ([]*schema.ResourceData, error) {
	meta := akamai.Meta(m)
	logger := meta.Log("Akamai GTM", "resourceGTMCidrMapImport")
	// create a context with logging for api calls
	ctx := context.Background()
	ctx = session.ContextWithOptions(
		ctx,
		session.WithContextLog(logger),
	)

	logger.Infof("cidrMap [%s] Import", d.Id())
	// pull domain and cidrMap out of cidrMap id
	domain, cidrMap, err := parseResourceStringID(d.Id())
	if err != nil {
		return []*schema.ResourceData{d}, err
	}
	cidr, err := inst.Client(meta).GetCidrMap(ctx, cidrMap, domain)
	if err != nil {
		return nil, err
	}
	if err := d.Set("domain", domain); err != nil {
		logger.Errorf("resourceGTMCidrMapImport failed: %s", err.Error())
	}
	if err := d.Set("wait_on_complete", true); err != nil {
		logger.Errorf("resourceGTMCidrMapImport failed: %s", err.Error())
	}
	populateTerraformCidrMapState(d, cidr, m)

	// use same Id as passed in
	logger.Infof("cidrMap [%s] [%s] Imported", d.Id(), d.Get("name"))
	return []*schema.ResourceData{d}, nil
}

// Delete GTM CidrMap.
func resourceGTMv1CidrMapDelete(ctx context.Context, d *schema.ResourceData, m interface{}) diag.Diagnostics {
	meta := akamai.Meta(m)
	logger := meta.Log("Akamai GTM", "resourceGTMCidrMapDelete")
	// create a context with logging for api calls
	ctx = session.ContextWithOptions(
		ctx,
		session.WithContextLog(logger),
	)

	logger.Debugf("Deleting cidrMap: %s", d.Id())
	var diags diag.Diagnostics
	// Get existing cidrMap
	domain, cidrMap, err := parseResourceStringID(d.Id())
	if err != nil {
		logger.Errorf("Invalid cidrMap ID: %s", d.Id())
		return diag.FromErr(err)
	}
	existCidr, err := inst.Client(meta).GetCidrMap(ctx, cidrMap, domain)
	if err != nil {
		logger.Errorf("CidrMapDelete cidrMap doesn't exist: %s", err.Error())
		return append(diags, diag.Diagnostic{
			Severity: diag.Error,
			Summary:  "cidrMap doesn't exist",
			Detail:   err.Error(),
		})
	}
	logger.Debugf("Deleting cidrMap: %v", existCidr)
	uStat, err := inst.Client(meta).DeleteCidrMap(ctx, existCidr, domain)
	if err != nil {
		logger.Errorf("cidrMap Delete failed: %s", err.Error())
		return append(diags, diag.Diagnostic{
			Severity: diag.Error,
			Summary:  "cidrMap Delete failed",
			Detail:   err.Error(),
		})
	}
	logger.Debugf("cidrMap Delete status: %v", uStat)
	if uStat.PropagationStatus == "DENIED" {
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
			logger.Infof("CidrMap Delete completed")
		} else {
			if err == nil {
				logger.Infof("cidrMap Delete pending")
			} else {
				logger.Errorf("cidrMap Delete failed [%s]", err.Error())
				return append(diags, diag.Diagnostic{
					Severity: diag.Error,
					Summary:  "cidrMap Delete failed",
					Detail:   err.Error(),
				})
			}
		}
	}

	// if successful ....
	d.SetId("")
	return nil
}

// Create and populate a new cidrMap object from cidrMap data
func populateNewCidrMapObject(ctx context.Context, meta akamai.OperationMeta, d *schema.ResourceData, m interface{}) *gtm.CidrMap {
	logger := meta.Log("Akamai GTM", "populateNewCidrMapObject")

	cidrMapName, err := tools.GetStringValue("name", d)
	if err != nil {
		logger.Errorf("cidrMap name not found in ResourceData: %s", err.Error())
	}

	cidrObj := inst.Client(meta).NewCidrMap(ctx, cidrMapName)
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
	logger := meta.Log("Akamai GTM", "populateTerraformCidrMapState")

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
	logger := meta.Log("Akamai GTM", "populateTerraformCidrAssignmentsState")

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
	logger := meta.Log("Akamai GTM", "populateCidrDefaultDCObject")

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
	logger := meta.Log("Akamai GTM", "populateTerraformCidrDefaultDCState")

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
