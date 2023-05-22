package gtm

import (
	"context"
	"fmt"
	"net/http"
	"sort"

	"github.com/akamai/AkamaiOPEN-edgegrid-golang/v6/pkg/gtm"
	"github.com/akamai/AkamaiOPEN-edgegrid-golang/v6/pkg/session"
	"github.com/akamai/terraform-provider-akamai/v3/pkg/akamai"
	"github.com/akamai/terraform-provider-akamai/v3/pkg/common/tf"
	"github.com/hashicorp/terraform-plugin-sdk/v2/diag"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/schema"
)

func resourceGTMv1ASmap() *schema.Resource {
	return &schema.Resource{
		CreateContext: resourceGTMv1ASmapCreate,
		ReadContext:   resourceGTMv1ASmapRead,
		UpdateContext: resourceGTMv1ASmapUpdate,
		DeleteContext: resourceGTMv1ASmapDelete,
		Importer: &schema.ResourceImporter{
			StateContext: resourceGTMv1ASmapImport,
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
						"as_numbers": {
							Type:     schema.TypeSet,
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
func validateDefaultDC(ctx context.Context, meta akamai.OperationMeta, ddcField []interface{}, domain string) error {

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

	dcID, ok := intrDcID.(int)
	if !ok || dcID == 0 {
		return fmt.Errorf("default Datacenter ID invalid")
	}
	dc, err := inst.Client(meta).GetDatacenter(ctx, dcID, domain)
	if dc == nil {
		if err != nil {
			apiError, ok := err.(*gtm.Error)
			if !ok || apiError.StatusCode != http.StatusNotFound {
				return fmt.Errorf("MapCreate Unexpected error verifying Default Datacenter exists: %s", err.Error())
			}
		}
		// ddc doesn't exist
		if ddc["datacenter_id"].(int) != gtm.MapDefaultDC {
			return fmt.Errorf(fmt.Sprintf("Default Datacenter %d does not exist", ddc["datacenter_id"].(int)))
		}
		_, err := inst.Client(meta).CreateMapsDefaultDatacenter(ctx, domain) // create if not already.
		if err != nil {
			return fmt.Errorf("MapCreate failed on Default Datacenter check: %s", err.Error())
		}
	}

	return nil
}

// Create a new GTM ASmap
func resourceGTMv1ASmapCreate(ctx context.Context, d *schema.ResourceData, m interface{}) diag.Diagnostics {
	meta := akamai.Meta(m)
	logger := meta.Log("Akamai GTM", "resourceGTMv1ASmapCreate")
	// create a context with logging for api calls
	ctx = session.ContextWithOptions(
		ctx,
		session.WithContextLog(logger),
	)

	domain, err := tf.GetStringValue("domain", d)
	if err != nil {
		logger.Errorf("Domain not initialized: %s", err.Error())
		return diag.FromErr(err)
	}

	name, err := tf.GetStringValue("name", d)
	if err != nil {
		logger.Errorf("asMap name not initialized: %s", err.Error())
		return diag.FromErr(err)
	}
	logger.Infof("Creating asMap [%s] in domain [%s]", name, domain)

	// Make sure Default Datacenter exists
	interfaceArray, err := tf.GetInterfaceArrayValue("default_datacenter", d)
	if err != nil {
		logger.Errorf("Default datacenter not initialized: %s", err.Error())
		return diag.FromErr(err)
	}
	var diags diag.Diagnostics
	if err := validateDefaultDC(ctx, meta, interfaceArray, domain); err != nil {
		logger.Errorf("Default datacenter validation error: %s", err.Error())
		return append(diags, diag.Diagnostic{
			Severity: diag.Error,
			Summary:  "Default datacenter validation error",
			Detail:   err.Error(),
		})
	}

	newAS := populateNewASmapObject(ctx, meta, d, m)
	logger.Debugf("Proposed New asMap: [%v]", newAS)
	cStatus, err := inst.Client(meta).CreateAsMap(ctx, newAS, domain)
	if err != nil {
		return append(diags, diag.Diagnostic{
			Severity: diag.Error,
			Summary:  "asMap Create failed",
			Detail:   err.Error(),
		})
	}
	logger.Debugf("asMap Create status: %v", cStatus.Status)
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
			logger.Infof("asMap Create completed")
		} else {
			if err == nil {
				logger.Infof("asMap Create pending")
			} else {
				logger.Errorf("asMap Create failed [%s]", err.Error())
				return append(diags, diag.Diagnostic{
					Severity: diag.Error,
					Summary:  "asMap Create failed",
					Detail:   err.Error(),
				})
			}
		}
	}

	// Give terraform the ID. Format domain:asMap
	asMapID := fmt.Sprintf("%s:%s", domain, cStatus.Resource.Name)
	logger.Debugf("Generated asMap Id: %s", asMapID)
	d.SetId(asMapID)
	return resourceGTMv1ASmapRead(ctx, d, m)

}

// read asMap. updates state with entire API result configuration.
func resourceGTMv1ASmapRead(ctx context.Context, d *schema.ResourceData, m interface{}) diag.Diagnostics {
	meta := akamai.Meta(m)
	logger := meta.Log("Akamai GTM", "resourceGTMv1ASmapRead")
	// create a context with logging for api calls
	ctx = session.ContextWithOptions(
		ctx,
		session.WithContextLog(logger),
	)

	logger.Debugf("Reading asMap: %s", d.Id())
	var diags diag.Diagnostics
	// retrieve the property and domain
	domain, asMap, err := parseResourceStringID(d.Id())
	if err != nil {
		return diag.FromErr(err)
	}
	as, err := inst.Client(meta).GetAsMap(ctx, asMap, domain)
	if err != nil {
		logger.Errorf("asMap Read error: %s", err.Error())
		return append(diags, diag.Diagnostic{
			Severity: diag.Error,
			Summary:  "asMap Read error",
			Detail:   err.Error(),
		})
	}
	populateTerraformASmapState(d, as, m)
	logger.Debugf("READ %v", as)
	return nil
}

// Update GTM ASmap
func resourceGTMv1ASmapUpdate(ctx context.Context, d *schema.ResourceData, m interface{}) diag.Diagnostics {
	meta := akamai.Meta(m)
	logger := meta.Log("Akamai GTM", "resourceGTMv1ASmapUpdate")
	// create a context with logging for api calls
	ctx = session.ContextWithOptions(
		ctx,
		session.WithContextLog(logger),
	)

	logger.Debugf("UPDATE asMap: %s", d.Id())
	var diags diag.Diagnostics
	// pull domain and asMap out of id
	domain, asMap, err := parseResourceStringID(d.Id())
	if err != nil {
		logger.Errorf("Invalid asMap ID: %s", d.Id())
		return diag.FromErr(err)
	}
	// Get existingASmap
	existAs, err := inst.Client(meta).GetAsMap(ctx, asMap, domain)
	if err != nil {
		logger.Errorf("asMap Update read error: %s", err.Error())
		return append(diags, diag.Diagnostic{
			Severity: diag.Error,
			Summary:  "asMap Update Read error",
			Detail:   err.Error(),
		})
	}
	logger.Debugf("asMap BEFORE: %v", existAs)
	populateASmapObject(d, existAs, m)
	logger.Debugf("asMap PROPOSED: %v", existAs)
	uStat, err := inst.Client(meta).UpdateAsMap(ctx, existAs, domain)
	if err != nil {
		logger.Errorf("asMap pdate: %s", err.Error())
		return append(diags, diag.Diagnostic{
			Severity: diag.Error,
			Summary:  "asMap Update error",
			Detail:   err.Error(),
		})
	}
	logger.Debugf("asMap Update status: %v", uStat)
	if uStat.PropagationStatus == "DENIED" {
		return append(diags, diag.Diagnostic{
			Severity: diag.Error,
			Summary:  uStat.Message,
		})
	}

	waitOnComplete, err := tf.GetBoolValue("wait_on_complete", d)
	if err != nil {
		return diag.FromErr(err)
	}
	if waitOnComplete {
		done, err := waitForCompletion(ctx, domain, m)
		if done {
			logger.Infof("ASmap Update completed")
		} else {
			if err == nil {
				logger.Infof("ASmap Update pending")
			} else {
				logger.Errorf("ASmap Update failed [%s]", err.Error())
				return append(diags, diag.Diagnostic{
					Severity: diag.Error,
					Summary:  "asMap Update failed",
					Detail:   err.Error(),
				})
			}
		}
	}

	return resourceGTMv1ASmapRead(ctx, d, m)
}

// Import GTM ASmap.
func resourceGTMv1ASmapImport(ctx context.Context, d *schema.ResourceData, m interface{}) ([]*schema.ResourceData, error) {
	meta := akamai.Meta(m)
	logger := meta.Log("Akamai GTM", "resourceGTMv1ASmapImport")
	// create a context with logging for api calls
	ctx = session.ContextWithOptions(
		ctx,
		session.WithContextLog(logger),
	)

	logger.Infof("asMap [%s] Import", d.Id())
	// pull domain and asMap out of asMap id
	domain, asMap, err := parseResourceStringID(d.Id())
	if err != nil {
		return []*schema.ResourceData{d}, err
	}
	as, err := inst.Client(meta).GetAsMap(ctx, asMap, domain)
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
	logger.Infof("asMap [%s] [%s] Imported", d.Id(), d.Get("name"))
	return []*schema.ResourceData{d}, nil
}

// Delete GTM ASmap.
func resourceGTMv1ASmapDelete(ctx context.Context, d *schema.ResourceData, m interface{}) diag.Diagnostics {
	meta := akamai.Meta(m)
	logger := meta.Log("Akamai GTM", "resourceGTMv1ASmapDelete")
	// create a context with logging for api calls
	ctx = session.ContextWithOptions(
		ctx,
		session.WithContextLog(logger),
	)

	logger.Debugf("Deleting asMap: %s", d.Id())
	var diags diag.Diagnostics
	// Get existing asMap
	domain, asMap, err := parseResourceStringID(d.Id())
	if err != nil {
		logger.Errorf("[ERROR] ASmap Delete: %s", err.Error())
		return diag.FromErr(err)
	}
	existAs, err := inst.Client(meta).GetAsMap(ctx, asMap, domain)
	if err != nil {
		logger.Errorf("ASmap Delete: %s", err.Error())
		return append(diags, diag.Diagnostic{
			Severity: diag.Error,
			Summary:  "asMap doesn't exist",
			Detail:   err.Error(),
		})
	}
	logger.Debugf("Deleting ASmap: %v", existAs)
	uStat, err := inst.Client(meta).DeleteAsMap(ctx, existAs, domain)
	if err != nil {
		logger.Errorf("ASmap Delete: %s", err.Error())
		return append(diags, diag.Diagnostic{
			Severity: diag.Error,
			Summary:  "asMap Delete failed",
			Detail:   err.Error(),
		})
	}
	logger.Debugf("asMap Delete status: %v", uStat)
	if uStat.PropagationStatus == "DENIED" {
		return append(diags, diag.Diagnostic{
			Severity: diag.Error,
			Summary:  uStat.Message,
		})
	}

	waitOnComplete, err := tf.GetBoolValue("wait_on_complete", d)
	if err != nil {
		return diag.FromErr(err)
	}
	if waitOnComplete {
		done, err := waitForCompletion(ctx, domain, m)
		if done {
			logger.Infof("asMap Delete completed")
		} else {
			if err == nil {
				logger.Infof("asMap Delete pending")
			} else {
				logger.Errorf("asMap Delete failed [%s]", err.Error())
				return append(diags, diag.Diagnostic{
					Severity: diag.Error,
					Summary:  "asMap Delete failed",
					Detail:   err.Error(),
				})
			}
		}
	}

	// if successful ....
	d.SetId("")
	return nil
}

// Create and populate a new asMap object from asMap data
func populateNewASmapObject(ctx context.Context, meta akamai.OperationMeta, d *schema.ResourceData, m interface{}) *gtm.AsMap {

	asMapName, _ := tf.GetStringValue("name", d)
	asObj := inst.Client(meta).NewAsMap(ctx, asMapName)
	asObj.DefaultDatacenter = &gtm.DatacenterBase{}
	asObj.Assignments = make([]*gtm.AsAssignment, 1)
	asObj.Links = make([]*gtm.Link, 1)
	populateASmapObject(d, asObj, m)

	return asObj

}

// Populate existing asMap object from asMap data
func populateASmapObject(d *schema.ResourceData, as *gtm.AsMap, m interface{}) {

	if v, err := tf.GetStringValue("name", d); err == nil {
		as.Name = v
	}
	populateAsAssignmentsObject(d, as, m)
	populateAsDefaultDCObject(d, as, m)

}

// Populate Terraform state from provided ASmap object
func populateTerraformASmapState(d *schema.ResourceData, as *gtm.AsMap, m interface{}) {
	meta := akamai.Meta(m)
	logger := meta.Log("Akamai GTM", "populateTerraformASmapState")

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
	logger := meta.Log("Akamai GTM", "populateAsAssignmentsObject")

	// pull apart List
	if asAssignmentsList, err := tf.GetListValue("assignment", d); err != nil {
		logger.Errorf("Assignment not set: %s", err.Error())
	} else {
		asAssignmentsObjList := make([]*gtm.AsAssignment, len(asAssignmentsList)) // create new object list
		for i, v := range asAssignmentsList {
			asMap := v.(map[string]interface{})
			asAssignment := gtm.AsAssignment{}
			asAssignment.DatacenterId = asMap["datacenter_id"].(int)
			asAssignment.Nickname = asMap["nickname"].(string)
			if asMap["as_numbers"] != nil {
				asNumbers, ok := asMap["as_numbers"].(*schema.Set)
				if !ok {
					logger.Errorf("wrong type conversion: expected *schema.Set, got %T", asNumbers)
				}
				ls := make([]int64, asNumbers.Len())
				for i, sl := range asNumbers.List() {
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
func populateTerraformAsAssignmentsState(d *schema.ResourceData, asm *gtm.AsMap, m interface{}) {
	meta := akamai.Meta(m)
	logger := meta.Log("Akamai GTM", "populateTerraformAsAssignmentsState")

	var asStateList []map[string]interface{}
	for _, as := range asm.Assignments {
		asNew := map[string]interface{}{
			"datacenter_id": as.DatacenterId,
			"nickname":      as.Nickname,
			"as_numbers":    as.AsNumbers,
		}
		asStateList = append(asStateList, asNew)
	}

	if err := d.Set("assignment", asStateList); err != nil {
		logger.Errorf("populateTerraformAsAssignmentsState failed: %s", err.Error())
	}
}

// create and populate GTM ASmap DefaultDatacenter object
func populateAsDefaultDCObject(d *schema.ResourceData, as *gtm.AsMap, m interface{}) {
	meta := akamai.Meta(m)
	logger := meta.Log("Akamai GTMv1", "resourceGTMv1ASmapDelete")

	// pull apart List
	if asDefaultDCList, err := tf.GetInterfaceArrayValue("default_datacenter", d); err != nil {
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
	logger := meta.Log("Akamai GTM", "populateTerraformAsDefaultDCState")

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

// assignmentDiffSuppress is a diff suppress function used in gtm_asmap, gtm_cidrmap and gtm_geomap resources
func assignmentDiffSuppress(_, _, _ string, d *schema.ResourceData) bool {
	logger := akamai.Log("Akamai GTM", "assignmentDiffSuppress")
	oldVal, newVal := d.GetChange("assignment")

	oldList, ok := oldVal.([]interface{})
	if !ok {
		logger.Warnf("wrong type conversion: expected []interface{}, got %T", oldList)
		return false
	}

	newList, ok := newVal.([]interface{})
	if !ok {
		logger.Warnf("wrong type conversion: expected []interface{}, got %T", newList)
		return false
	}

	if len(oldList) != len(newList) {
		return false
	}

	sort.Slice(oldList, func(i, j int) bool {
		return oldList[i].(map[string]interface{})["datacenter_id"].(int) < oldList[j].(map[string]interface{})["datacenter_id"].(int)
	})
	sort.Slice(newList, func(i, j int) bool {
		return newList[i].(map[string]interface{})["datacenter_id"].(int) < newList[j].(map[string]interface{})["datacenter_id"].(int)
	})

	attrName, err := resolveAttrName(oldList)
	if err != nil {
		logger.Warnf("resolveAttrName: %s", err)
		return false
	}

	length := len(oldList)
	var oldAssignment, newAssignment map[string]interface{}
	for i := 0; i < length; i++ {
		oldAssignment, ok = oldList[i].(map[string]interface{})
		if !ok {
			logger.Warnf("wrong type conversion: expected map[string]interface{}, got %T", oldAssignment)
		}
		newAssignment, ok = newList[i].(map[string]interface{})
		if !ok {
			logger.Warnf("wrong type conversion: expected map[string]interface{}, got %T", newAssignment)
		}
		for k, v := range oldAssignment {
			if k == attrName {
				switch attrName {
				case "blocks":
					if !blocksEqual(oldAssignment[attrName], newAssignment[attrName]) {
						return false
					}
				case "countries":
					if !countriesEqual(oldAssignment[attrName], newAssignment[attrName]) {
						return false
					}
				case "as_numbers":
					if !asNumbersEqual(oldAssignment[attrName], newAssignment[attrName]) {
						return false
					}
				default:
					logger.Warn("no expected attribute is present, should be one of [as_numbers, load_servers, countries]")
				}

			} else {
				if newAssignment[k] != v {
					return false
				}
			}
		}
	}

	return true
}

// asNumbersEqual checks whether the as_numbers are equal
func asNumbersEqual(old, new interface{}) bool {
	logger := akamai.Log("Akamai GTM", "asNumbersEqual")

	oldVal, ok := old.(*schema.Set)
	if !ok {
		logger.Warnf("wrong type conversion: expected *schema.Set, got %T", oldVal)
		return false
	}

	newVal, ok := new.(*schema.Set)
	if !ok {
		logger.Warnf("wrong type conversion: expected *schema.Set, got %T", newVal)
		return false
	}

	if oldVal.Len() != newVal.Len() {
		return false
	}

	numbers := make(map[int]bool, oldVal.Len())
	for _, num := range oldVal.List() {
		numbers[num.(int)] = true
	}

	for _, num := range newVal.List() {
		_, ok = numbers[num.(int)]
		if !ok {
			return false
		}
	}

	return true
}

// resolveAttrName resolves specific assignment attribute, based on a resource
func resolveAttrName(list []interface{}) (string, error) {
	if len(list) == 0 {
		return "", fmt.Errorf("there are no elements in the list")
	}

	entry, ok := list[0].(map[string]interface{})
	if !ok {
		return "", fmt.Errorf("expected map[string]interface{}, got %T", entry)
	}

	_, ok = entry["blocks"]
	if ok {
		return "blocks", nil
	}
	_, ok = entry["countries"]
	if ok {
		return "countries", nil
	}
	_, ok = entry["as_numbers"]
	if ok {
		return "as_numbers", nil
	}

	return "", fmt.Errorf("there is no attribute matching one of: [blocks, countries, as_numbers]")
}
