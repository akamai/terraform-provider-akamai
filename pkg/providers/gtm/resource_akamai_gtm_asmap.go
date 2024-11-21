package gtm

import (
	"context"
	"errors"
	"fmt"
	"net/http"
	"sort"

	"github.com/akamai/AkamaiOPEN-edgegrid-golang/v9/pkg/gtm"
	"github.com/akamai/AkamaiOPEN-edgegrid-golang/v9/pkg/session"
	"github.com/akamai/terraform-provider-akamai/v6/pkg/common/tf"
	"github.com/akamai/terraform-provider-akamai/v6/pkg/logger"
	"github.com/akamai/terraform-provider-akamai/v6/pkg/meta"
	"github.com/hashicorp/terraform-plugin-sdk/v2/diag"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/schema"
)

const asMapAlreadyExistsError = "AsMap with provided `name` for specific `domain` already exists. Please import specific asmap using following command: terraform import akamai_gtm_asmap.<your_resource_name> \"%s:%s\""

func resourceGTMv1ASMap() *schema.Resource {
	return &schema.Resource{
		CreateContext: resourceGTMv1ASMapCreate,
		ReadContext:   resourceGTMv1ASMapRead,
		UpdateContext: resourceGTMv1ASMapUpdate,
		DeleteContext: resourceGTMv1ASMapDelete,
		Importer: &schema.ResourceImporter{
			StateContext: resourceGTMv1ASMapImport,
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
func validateDefaultDC(ctx context.Context, meta meta.Meta, ddcField []interface{}, domain string) error {

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
	dc, err := Client(meta).GetDatacenter(ctx, gtm.GetDatacenterRequest{
		DomainName:   domain,
		DatacenterID: dcID,
	})
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
		_, err := Client(meta).CreateMapsDefaultDatacenter(ctx, domain) // create if not already.
		if err != nil {
			return fmt.Errorf("MapCreate failed on Default Datacenter check: %s", err.Error())
		}
	}

	return nil
}

func resourceGTMv1ASMapCreate(ctx context.Context, d *schema.ResourceData, m interface{}) diag.Diagnostics {
	meta := meta.Must(m)
	logger := meta.Log("Akamai GTM", "resourceGTMv1ASmapCreate")
	// create a context with logging for api calls
	ctx = session.ContextWithOptions(
		ctx,
		session.WithContextLog(logger),
	)
	var diags diag.Diagnostics

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

	as, err := Client(meta).GetASMap(ctx, gtm.GetASMapRequest{
		ASMapName:  name,
		DomainName: domain,
	})
	if err != nil && !errors.Is(err, gtm.ErrNotFound) {
		logger.Errorf("asMap Read error: %s", err.Error())
		return append(diags, diag.Diagnostic{
			Severity: diag.Error,
			Summary:  "asMap Read error",
			Detail:   err.Error(),
		})
	}
	if as != nil {
		asMapAlreadyExists := fmt.Sprintf(asMapAlreadyExistsError, domain, name)
		logger.Errorf(asMapAlreadyExists)
		return append(diags, diag.Diagnostic{
			Severity: diag.Error,
			Summary:  "asMap already exists error",
			Detail:   asMapAlreadyExists,
		})
	}

	logger.Infof("Creating asMap [%s] in domain [%s]", name, domain)

	// Make sure Default Datacenter exists
	interfaceArray, err := tf.GetInterfaceArrayValue("default_datacenter", d)
	if err != nil {
		logger.Errorf("Default datacenter not initialized: %s", err.Error())
		return diag.FromErr(err)
	}
	if err := validateDefaultDC(ctx, meta, interfaceArray, domain); err != nil {
		logger.Errorf("Default datacenter validation error: %s", err.Error())
		return append(diags, diag.Diagnostic{
			Severity: diag.Error,
			Summary:  "Default datacenter validation error",
			Detail:   err.Error(),
		})
	}

	newAS, err := populateNewASMapObject(d, m)
	if err != nil {
		return append(diags, diag.Diagnostic{
			Severity: diag.Error,
			Summary:  "asMap populate failed",
			Detail:   err.Error(),
		})
	}
	logger.Debugf("Proposed New asMap: [%v]", newAS)
	cStatus, err := Client(meta).CreateASMap(ctx, gtm.CreateASMapRequest{
		ASMap:      newAS,
		DomainName: domain,
	})
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
	return resourceGTMv1ASMapRead(ctx, d, m)

}

func resourceGTMv1ASMapRead(ctx context.Context, d *schema.ResourceData, m interface{}) diag.Diagnostics {
	meta := meta.Must(m)
	logger := meta.Log("Akamai GTM", "resourceGTMv1ASMapRead")
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
	as, err := Client(meta).GetASMap(ctx, gtm.GetASMapRequest{
		ASMapName:  asMap,
		DomainName: domain,
	})
	if err != nil {
		logger.Errorf("asMap Read error: %s", err.Error())
		return append(diags, diag.Diagnostic{
			Severity: diag.Error,
			Summary:  "asMap Read error",
			Detail:   err.Error(),
		})
	}
	populateTerraformASMapState(d, as, m)
	logger.Debugf("READ %v", as)
	return nil
}

func resourceGTMv1ASMapUpdate(ctx context.Context, d *schema.ResourceData, m interface{}) diag.Diagnostics {
	meta := meta.Must(m)
	logger := meta.Log("Akamai GTM", "resourceGTMv1ASMapUpdate")
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
	existAs, err := Client(meta).GetASMap(ctx, gtm.GetASMapRequest{
		ASMapName:  asMap,
		DomainName: domain,
	})
	if err != nil {
		logger.Errorf("asMap Update read error: %s", err.Error())
		return append(diags, diag.Diagnostic{
			Severity: diag.Error,
			Summary:  "asMap Update Read error",
			Detail:   err.Error(),
		})
	}
	logger.Debugf("asMap BEFORE: %v", existAs)
	newAs := createASMapStruct(existAs)
	populateASMapObject(d, newAs, m)
	logger.Debugf("asMap PROPOSED: %v", existAs)
	uStat, err := Client(meta).UpdateASMap(ctx, gtm.UpdateASMapRequest{
		ASMap:      newAs,
		DomainName: domain,
	})
	if err != nil {
		logger.Errorf("asMap pdate: %s", err.Error())
		return append(diags, diag.Diagnostic{
			Severity: diag.Error,
			Summary:  "asMap Update error",
			Detail:   err.Error(),
		})
	}
	logger.Debugf("asMap Update status: %v", uStat)
	if uStat.Status.PropagationStatus == "DENIED" {
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

	return resourceGTMv1ASMapRead(ctx, d, m)
}

func resourceGTMv1ASMapImport(ctx context.Context, d *schema.ResourceData, m interface{}) ([]*schema.ResourceData, error) {
	meta := meta.Must(m)
	logger := meta.Log("Akamai GTM", "resourceGTMv1ASMapImport")
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
	as, err := Client(meta).GetASMap(ctx, gtm.GetASMapRequest{
		ASMapName:  asMap,
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
	populateTerraformASMapState(d, as, m)

	// use same Id as passed in
	logger.Infof("asMap [%s] [%s] Imported", d.Id(), d.Get("name"))
	return []*schema.ResourceData{d}, nil
}

func resourceGTMv1ASMapDelete(ctx context.Context, d *schema.ResourceData, m interface{}) diag.Diagnostics {
	meta := meta.Must(m)
	logger := meta.Log("Akamai GTM", "resourceGTMv1ASMapDelete")
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
		logger.Errorf("[ERROR] ASMap Delete: %s", err.Error())
		return diag.FromErr(err)
	}
	existAs, err := Client(meta).GetASMap(ctx, gtm.GetASMapRequest{
		ASMapName:  asMap,
		DomainName: domain,
	})
	if err != nil {
		logger.Errorf("ASMap Delete: %s", err.Error())
		return append(diags, diag.Diagnostic{
			Severity: diag.Error,
			Summary:  "asMap doesn't exist",
			Detail:   err.Error(),
		})
	}
	newAs := createASMapStruct(existAs)
	logger.Debugf("Deleting ASmap: %v", newAs)
	uStat, err := Client(meta).DeleteASMap(ctx, gtm.DeleteASMapRequest{
		ASMapName:  asMap,
		DomainName: domain,
	})
	if err != nil {
		logger.Errorf("ASMap Delete: %s", err.Error())
		return append(diags, diag.Diagnostic{
			Severity: diag.Error,
			Summary:  "asMap Delete failed",
			Detail:   err.Error(),
		})
	}
	logger.Debugf("asMap Delete status: %v", uStat)
	if uStat.Status.PropagationStatus == "DENIED" {
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

	d.SetId("")
	return nil
}

// Create and populate a new asMap object from asMap data
func populateNewASMapObject(d *schema.ResourceData, m interface{}) (*gtm.ASMap, error) {

	asMapName, err := tf.GetStringValue("name", d)
	if err != nil {
		return nil, err
	}
	asObj := &gtm.ASMap{
		Name:              asMapName,
		DefaultDatacenter: &gtm.DatacenterBase{},
		Assignments:       make([]gtm.ASAssignment, 1),
	}
	populateASMapObject(d, asObj, m)

	return asObj, nil

}

// Populate existing asMap object from asMap data
func populateASMapObject(d *schema.ResourceData, as *gtm.ASMap, m interface{}) {
	if v, err := tf.GetStringValue("name", d); err == nil {
		as.Name = v
	}
	populateASAssignmentsObject(d, as, m)
	populateASDefaultDCObject(d, as, m)
}

// Populate Terraform state from provided ASMap object
func populateTerraformASMapState(d *schema.ResourceData, as *gtm.GetASMapResponse, m interface{}) {
	meta := meta.Must(m)
	logger := meta.Log("Akamai GTM", "populateTerraformASMapState")

	// walk through all state elements
	if err := d.Set("name", as.Name); err != nil {
		logger.Errorf("populateTerraformASMapState failed: %s", err.Error())
	}
	populateTerraformASAssignmentsState(d, as, m)
	populateTerraformASDefaultDCState(d, as, m)
}

// create and populate GTM ASMap Assignments object
func populateASAssignmentsObject(d *schema.ResourceData, as *gtm.ASMap, m interface{}) {
	meta := meta.Must(m)
	logger := meta.Log("Akamai GTM", "populateASAssignmentsObject")

	// pull apart List
	if asAssignmentsList, err := tf.GetListValue("assignment", d); err != nil {
		logger.Errorf("Assignment not set: %s", err.Error())
	} else {
		asAssignmentsObjList := make([]gtm.ASAssignment, len(asAssignmentsList)) // create new object list
		for i, v := range asAssignmentsList {
			asMap := v.(map[string]interface{})
			asAssignment := gtm.ASAssignment{}
			asAssignment.DatacenterID = asMap["datacenter_id"].(int)
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
				asAssignment.ASNumbers = ls
			}
			asAssignmentsObjList[i] = asAssignment
		}
		as.Assignments = asAssignmentsObjList
	}
}

// create and populate Terraform asMap assignments schema
func populateTerraformASAssignmentsState(d *schema.ResourceData, asm *gtm.GetASMapResponse, m interface{}) {
	meta := meta.Must(m)
	logger := meta.Log("Akamai GTM", "populateTerraformASAssignmentsState")

	var asStateList []map[string]interface{}
	for _, as := range asm.Assignments {
		asNew := map[string]interface{}{
			"datacenter_id": as.DatacenterID,
			"nickname":      as.Nickname,
			"as_numbers":    as.ASNumbers,
		}
		asStateList = append(asStateList, asNew)
	}

	if err := d.Set("assignment", asStateList); err != nil {
		logger.Errorf("populateTerraformASAssignmentsState failed: %s", err.Error())
	}
}

// create and populate GTM ASMap DefaultDatacenter object
func populateASDefaultDCObject(d *schema.ResourceData, as *gtm.ASMap, m interface{}) {
	meta := meta.Must(m)
	logger := meta.Log("Akamai GTMv1", "resourceGTMv1ASMapDelete")

	// pull apart List
	if asDefaultDCList, err := tf.GetInterfaceArrayValue("default_datacenter", d); err != nil {
		logger.Infof("No default datacenter specified: %s", err.Error())
	} else {
		if len(asDefaultDCList) > 0 {
			asDefaultDCObj := gtm.DatacenterBase{} // create new object
			asMap := asDefaultDCList[0].(map[string]interface{})
			if asMap["datacenter_id"] != nil && asMap["datacenter_id"].(int) != 0 {
				asDefaultDCObj.DatacenterID = asMap["datacenter_id"].(int)
				asDefaultDCObj.Nickname = asMap["nickname"].(string)
			} else {
				logger.Infof("No Default Datacenter specified")
				var nilInt int
				asDefaultDCObj.DatacenterID = nilInt
				asDefaultDCObj.Nickname = ""
			}
			as.DefaultDatacenter = &asDefaultDCObj
		}
	}
}

// create and populate Terraform asMap default_datacenter schema
func populateTerraformASDefaultDCState(d *schema.ResourceData, as *gtm.GetASMapResponse, m interface{}) {
	meta := meta.Must(m)
	logger := meta.Log("Akamai GTM", "populateTerraformASDefaultDCState")

	ddcListNew := make([]interface{}, 1)
	ddcNew := map[string]interface{}{
		"datacenter_id": as.DefaultDatacenter.DatacenterID,
		"nickname":      as.DefaultDatacenter.Nickname,
	}
	ddcListNew[0] = ddcNew
	if err := d.Set("default_datacenter", ddcListNew); err != nil {
		logger.Errorf("populateTerraformASDefaultDCState failed: %s", err.Error())
	}
}

// assignmentDiffSuppress is a diff suppress function used in gtm_asmap, gtm_cidrmap and gtm_geomap resources
func assignmentDiffSuppress(_, _, _ string, d *schema.ResourceData) bool {
	logger := logger.Get("Akamai GTM", "assignmentDiffSuppress")
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
	logger := logger.Get("Akamai GTM", "asNumbersEqual")

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

// createASMapStruct converts response from GetASMapResponse into ASMap
func createASMapStruct(asmap *gtm.GetASMapResponse) *gtm.ASMap {
	if asmap != nil {
		return &gtm.ASMap{
			DefaultDatacenter: asmap.DefaultDatacenter,
			Assignments:       asmap.Assignments,
			Name:              asmap.Name,
			Links:             asmap.Links,
		}
	}
	return nil
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
