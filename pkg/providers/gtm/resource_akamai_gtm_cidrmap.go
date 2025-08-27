package gtm

import (
	"context"
	"errors"
	"fmt"

	"github.com/akamai/AkamaiOPEN-edgegrid-golang/v12/pkg/gtm"
	"github.com/akamai/AkamaiOPEN-edgegrid-golang/v12/pkg/session"
	"github.com/akamai/terraform-provider-akamai/v8/pkg/common/tf"
	"github.com/akamai/terraform-provider-akamai/v8/pkg/log"
	"github.com/akamai/terraform-provider-akamai/v8/pkg/meta"
	"github.com/hashicorp/terraform-plugin-sdk/v2/diag"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/schema"
)

const cidrMapAlreadyExistsError = "CidrMap with provided `name` for specific `domain` already exists. Please import specific cidrmap using following command: terraform import akamai_gtm_cidrmap.<your_resource_name> \"%s:%s\""

func resourceGTMv1CIDRMap() *schema.Resource {
	return &schema.Resource{
		CreateContext: resourceGTMv1CIDRMapCreate,
		ReadContext:   resourceGTMv1CIDRMapRead,
		UpdateContext: resourceGTMv1CIDRMapUpdate,
		DeleteContext: resourceGTMv1CIDRMapDelete,
		Importer: &schema.ResourceImporter{
			State: resourceGTMv1CIDRMapImport,
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
						"blocks": {
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

func resourceGTMv1CIDRMapCreate(ctx context.Context, d *schema.ResourceData, m interface{}) diag.Diagnostics {
	meta := meta.Must(m)
	logger := meta.Log("Akamai GTM", "resourceGTMCIDRMapCreate")
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

	cidr, err := Client(meta).GetCIDRMap(ctx, gtm.GetCIDRMapRequest{
		DomainName: domain,
		MapName:    name,
	})
	if err != nil && !errors.Is(err, gtm.ErrNotFound) {
		logger.Errorf("cidrMap read error: %s", err.Error())
		return append(diags, diag.Diagnostic{
			Severity: diag.Error,
			Summary:  "cidrMap read error",
			Detail:   err.Error(),
		})
	}

	if cidr != nil {
		cidrMapAlreadyExists := fmt.Sprintf(cidrMapAlreadyExistsError, domain, name)
		logger.Errorf(cidrMapAlreadyExists)
		return append(diags, diag.Diagnostic{
			Severity: diag.Error,
			Summary:  "cidrMap already exists error",
			Detail:   cidrMapAlreadyExists,
		})
	}

	logger.Infof("Creating cidrMap [%s] in domain [%s]", name, domain)
	// Make sure Default Datacenter exists
	defaultDatacenter, err := tf.GetInterfaceArrayValue("default_datacenter", d)
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

	newCidr := populateNewCIDRMapObject(meta, d, m)
	logger.Debugf("Proposed New cidrMap: [%v]", newCidr)
	cStatus, err := Client(meta).CreateCIDRMap(ctx, gtm.CreateCIDRMapRequest{
		CIDR:       newCidr,
		DomainName: domain,
	})
	if err != nil {
		logger.Errorf("cidrMap create error: %s", err.Error())
		return append(diags, diag.Diagnostic{
			Severity: diag.Error,
			Summary:  "cidrMap create error",
			Detail:   err.Error(),
		})
	}
	logger.Debugf("cidrMap create status: %v", cStatus.Status)
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
			logger.Infof("cidrMap create completed")
		} else {
			if err == nil {
				logger.Infof("cidrMap create pending")
			} else {
				logger.Errorf("cidrMap create error: %s", err.Error())
				return append(diags, diag.Diagnostic{
					Severity: diag.Error,
					Summary:  "cidrMap create error",
					Detail:   err.Error(),
				})
			}
		}
	}

	// Give terraform the ID. Format domain:cidrMap
	cidrMapID := fmt.Sprintf("%s:%s", domain, cStatus.Resource.Name)
	logger.Debugf("Generated cidrMap resource Id: %s", cidrMapID)
	d.SetId(cidrMapID)
	return resourceGTMv1CIDRMapRead(ctx, d, m)

}

func resourceGTMv1CIDRMapRead(ctx context.Context, d *schema.ResourceData, m interface{}) diag.Diagnostics {
	meta := meta.Must(m)
	logger := meta.Log("Akamai GTM", "resourceGTMCIDRMapRead")
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
	cidr, err := Client(meta).GetCIDRMap(ctx, gtm.GetCIDRMapRequest{
		DomainName: domain,
		MapName:    cidrMap,
	})
	if errors.Is(err, gtm.ErrNotFound) {
		d.SetId("")
		return nil
	}
	if err != nil {
		logger.Errorf("cidrMap read error: %s", err.Error())
		return append(diags, diag.Diagnostic{
			Severity: diag.Error,
			Summary:  "cidrMap read error",
			Detail:   err.Error(),
		})
	}
	populateTerraformCIDRMapState(d, cidr, m)
	logger.Debugf("READ %v", cidr)
	return nil
}

func resourceGTMv1CIDRMapUpdate(ctx context.Context, d *schema.ResourceData, m interface{}) diag.Diagnostics {
	meta := meta.Must(m)
	logger := meta.Log("Akamai GTM", "resourceGTMCIDRMapUpdate")
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
	existCidr, err := Client(meta).GetCIDRMap(ctx, gtm.GetCIDRMapRequest{
		DomainName: domain,
		MapName:    cidrMap,
	})
	if err != nil {
		logger.Errorf("cidrMap read error: %s", err.Error())
		return append(diags, diag.Diagnostic{
			Severity: diag.Error,
			Summary:  "cidrMap read error",
			Detail:   err.Error(),
		})
	}
	newCidr := createCIDRMapStruct(existCidr)
	logger.Debugf("Updating cidrMap BEFORE: %v", newCidr)
	populateCIDRMapObject(d, newCidr, m)
	logger.Debugf("Updating cidrMap PROPOSED: %v", existCidr)
	uStat, err := Client(meta).UpdateCIDRMap(ctx, gtm.UpdateCIDRMapRequest{
		CIDR:       newCidr,
		DomainName: domain,
	})
	if err != nil {
		logger.Errorf("cidrMap update error: %s", err.Error())
		return append(diags, diag.Diagnostic{
			Severity: diag.Error,
			Summary:  "cidrMap update error",
			Detail:   err.Error(),
		})
	}
	logger.Debugf("cidrMap update status: %v", uStat)
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
			logger.Infof("cidrMap update completed")
		} else {
			if err == nil {
				logger.Infof("cidrMap update pending")
			} else {
				logger.Errorf("cidrMap update error: %s", err.Error())
				return append(diags, diag.Diagnostic{
					Severity: diag.Error,
					Summary:  "cidrMap update error",
					Detail:   err.Error(),
				})
			}
		}
	}

	return resourceGTMv1CIDRMapRead(ctx, d, m)
}

func resourceGTMv1CIDRMapImport(d *schema.ResourceData, m interface{}) ([]*schema.ResourceData, error) {
	meta := meta.Must(m)
	logger := meta.Log("Akamai GTM", "resourceGTMCIDRMapImport")
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
	cidr, err := Client(meta).GetCIDRMap(ctx, gtm.GetCIDRMapRequest{
		DomainName: domain,
		MapName:    cidrMap,
	})
	if err != nil {
		return nil, err
	}
	if err := d.Set("domain", domain); err != nil {
		logger.Errorf("resourceGTMCidrMapImport failed: %s", err.Error())
	}
	if err := d.Set("wait_on_complete", true); err != nil {
		logger.Errorf("resourceGTMCidrMapImport failed: %s", err.Error())
	}
	populateTerraformCIDRMapState(d, cidr, m)

	// use same Id as passed in
	logger.Infof("cidrMap [%s] [%s] Imported", d.Id(), d.Get("name"))
	return []*schema.ResourceData{d}, nil
}

func resourceGTMv1CIDRMapDelete(ctx context.Context, d *schema.ResourceData, m interface{}) diag.Diagnostics {
	meta := meta.Must(m)
	logger := meta.Log("Akamai GTM", "resourceGTMCIDRMapDelete")
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
	existCidr, err := Client(meta).GetCIDRMap(ctx, gtm.GetCIDRMapRequest{
		DomainName: domain,
		MapName:    cidrMap,
	})
	if err != nil {
		logger.Errorf("cidrMap doesn't exist: %s", err.Error())
		return append(diags, diag.Diagnostic{
			Severity: diag.Error,
			Summary:  "cidrMap doesn't exist",
			Detail:   err.Error(),
		})
	}
	newCidr := createCIDRMapStruct(existCidr)
	logger.Debugf("Deleting cidrMap: %v", newCidr)
	uStat, err := Client(meta).DeleteCIDRMap(ctx, gtm.DeleteCIDRMapRequest{
		MapName:    cidrMap,
		DomainName: domain,
	})
	if err != nil {
		logger.Errorf("cidrMap delete error: %s", err.Error())
		return append(diags, diag.Diagnostic{
			Severity: diag.Error,
			Summary:  "cidrMap delete error",
			Detail:   err.Error(),
		})
	}
	logger.Debugf("cidrMap delete status: %v", uStat)
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
			logger.Infof("cidrMap delete completed")
		} else {
			if err == nil {
				logger.Infof("cidrMap delete pending")
			} else {
				logger.Errorf("cidrMap delete error: %s", err.Error())
				return append(diags, diag.Diagnostic{
					Severity: diag.Error,
					Summary:  "cidrMap delete error",
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
func populateNewCIDRMapObject(meta meta.Meta, d *schema.ResourceData, m interface{}) *gtm.CIDRMap {
	logger := meta.Log("Akamai GTM", "populateNewCIDRMapObject")

	cidrMapName, err := tf.GetStringValue("name", d)
	if err != nil {
		logger.Errorf("cidrMap name not found in ResourceData: %s", err.Error())
	}

	cidrObj := &gtm.CIDRMap{
		Name:              cidrMapName,
		DefaultDatacenter: &gtm.DatacenterBase{},
	}
	populateCIDRMapObject(d, cidrObj, m)

	return cidrObj
}

// Populate existing cidrMap object from cidrMap data
func populateCIDRMapObject(d *schema.ResourceData, cidr *gtm.CIDRMap, m interface{}) {
	if v, err := tf.GetStringValue("name", d); err == nil {
		cidr.Name = v
	}
	populateCIDRAssignmentsObject(d, cidr, m)
	populateCIDRDefaultDCObject(d, cidr, m)
}

// Populate Terraform state from provided CIDRMap object
func populateTerraformCIDRMapState(d *schema.ResourceData, cidr *gtm.GetCIDRMapResponse, m interface{}) {
	meta := meta.Must(m)
	logger := meta.Log("Akamai GTM", "populateTerraformCidrMapState")

	// walk through all state elements
	if err := d.Set("name", cidr.Name); err != nil {
		logger.Errorf("populateTerraformCidrMapState failed: %s", err.Error())
	}
	populateTerraformCIDRAssignmentsState(d, cidr, m)
	populateTerraformCIDRDefaultDCState(d, cidr, m)
}

// create and populate GTM CidrMap Assignments object
func populateCIDRAssignmentsObject(d *schema.ResourceData, cidr *gtm.CIDRMap, m interface{}) {
	meta := meta.Must(m)
	logger := meta.Log("Akamai GTM", "populateCidrAssignmentsObject")

	// pull apart List
	if cassgns := d.Get("assignment"); cassgns != nil {
		cidrAssignmentsList := cassgns.([]interface{})
		cidrAssignmentsObjList := make([]gtm.CIDRAssignment, len(cidrAssignmentsList)) // create new object list
		for i, v := range cidrAssignmentsList {
			cidrMap := v.(map[string]interface{})
			cidrAssignment := gtm.CIDRAssignment{}
			cidrAssignment.DatacenterID = cidrMap["datacenter_id"].(int)
			cidrAssignment.Nickname = cidrMap["nickname"].(string)
			if cidrMap["blocks"] != nil {
				blocks, ok := cidrMap["blocks"].(*schema.Set)
				if !ok {
					logger.Warnf("wrong type conversion: expected *schema.Set, got %T", blocks)
				}
				ls := make([]string, blocks.Len())
				for i, sl := range blocks.List() {
					ls[i] = sl.(string)
				}
				cidrAssignment.Blocks = ls
			}
			cidrAssignmentsObjList[i] = cidrAssignment
		}
		cidr.Assignments = cidrAssignmentsObjList
	}
}

// create and populate Terraform cidrMap assignments schema
func populateTerraformCIDRAssignmentsState(d *schema.ResourceData, cidr *gtm.GetCIDRMapResponse, m interface{}) {
	meta := meta.Must(m)
	logger := meta.Log("Akamai GTM", "populateTerraformCidrAssignmentsState")

	objectInventory := make(map[int]gtm.CIDRAssignment, len(cidr.Assignments))
	if len(cidr.Assignments) > 0 {
		for _, aObj := range cidr.Assignments {
			objectInventory[aObj.DatacenterID] = aObj
		}
	}
	aStateList, err := tf.GetInterfaceArrayValue("assignment", d)
	if err != nil {
		logger.Errorf("Cidr Assignment list NOT FOUND in ResourceData: %s", err.Error())
	}
	for _, aMap := range aStateList {
		a := aMap.(map[string]interface{})
		objIndex := a["datacenter_id"].(int)
		aObject, ok := objectInventory[objIndex]
		if !ok {
			logger.Warnf("Cidr Assignment %d NOT FOUND in returned GTM Object", a["datacenter_id"])
			continue
		}
		a["datacenter_id"] = aObject.DatacenterID
		a["nickname"] = aObject.Nickname
		a["blocks"] = reconcileTerraformLists(a["blocks"].(*schema.Set).List(), convertStringToInterfaceList(aObject.Blocks, m), m)
		// remove object
		delete(objectInventory, objIndex)
	}
	if len(objectInventory) > 0 {
		logger.Debugf("CIDR Assignment objects left...")
		// Objects not in the state yet. Add. Unfortunately, they not align with instance indices in the config
		for _, maObj := range objectInventory {
			aNew := map[string]interface{}{
				"datacenter_id": maObj.DatacenterID,
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

// create and populate GTM CIDRMap DefaultDatacenter object
func populateCIDRDefaultDCObject(d *schema.ResourceData, cidr *gtm.CIDRMap, m interface{}) {
	meta := meta.Must(m)
	logger := meta.Log("Akamai GTM", "populateCIDRDefaultDCObject")

	// pull apart List
	if cidrDefaultDCList, err := tf.GetInterfaceArrayValue("default_datacenter", d); err != nil {
		logger.Infof("No default datacenter specified: %s", err.Error())
	} else {
		if len(cidrDefaultDCList) > 0 {
			cidrDefaultDCObj := gtm.DatacenterBase{} // create new object
			cidrddMap := cidrDefaultDCList[0].(map[string]interface{})
			if cidrddMap["datacenter_id"] != nil && cidrddMap["datacenter_id"].(int) != 0 {
				cidrDefaultDCObj.DatacenterID = cidrddMap["datacenter_id"].(int)
				cidrDefaultDCObj.Nickname = cidrddMap["nickname"].(string)
			} else {
				logger.Infof("No Default Datacenter specified")
				var nilInt int
				cidrDefaultDCObj.DatacenterID = nilInt
				cidrDefaultDCObj.Nickname = ""
			}
			cidr.DefaultDatacenter = &cidrDefaultDCObj
		}
	}
}

// create and populate Terraform cidrMap default_datacenter schema
func populateTerraformCIDRDefaultDCState(d *schema.ResourceData, cidr *gtm.GetCIDRMapResponse, m interface{}) {
	meta := meta.Must(m)
	logger := meta.Log("Akamai GTM", "populateTerraformCIDRDefaultDCState")

	ddcListNew := make([]interface{}, 1)
	ddcNew := map[string]interface{}{
		"datacenter_id": cidr.DefaultDatacenter.DatacenterID,
		"nickname":      cidr.DefaultDatacenter.Nickname,
	}
	ddcListNew[0] = ddcNew
	if err := d.Set("default_datacenter", ddcListNew); err != nil {
		logger.Errorf("populateTerraformCidrDefaultDCState failed: %s", err.Error())
	}
}

// createCIDRMapStruct converts response from GetCIDRMapResponse into CIDRMap
func createCIDRMapStruct(cidr *gtm.GetCIDRMapResponse) *gtm.CIDRMap {
	if cidr != nil {
		return &gtm.CIDRMap{
			DefaultDatacenter: cidr.DefaultDatacenter,
			Assignments:       cidr.Assignments,
			Name:              cidr.Name,
			Links:             cidr.Links,
		}
	}
	return nil
}

// blocksEqual checks whether blocks are equal
func blocksEqual(o, n interface{}) bool {
	logger := log.Get("Akamai GTM", "blocksEqual")

	oldBlocks, ok := o.(*schema.Set)
	if !ok {
		logger.Warnf("wrong type conversion: expected *schema.Set, got %T", oldBlocks)
		return false
	}

	newBlocks, ok := n.(*schema.Set)
	if !ok {
		logger.Warnf("wrong type conversion: expected *schema.Set, got %T", newBlocks)
		return false
	}

	if oldBlocks.Len() != newBlocks.Len() {
		return false
	}

	blocks := make(map[string]bool, oldBlocks.Len())
	for _, block := range oldBlocks.List() {
		blocks[block.(string)] = true
	}

	for _, block := range newBlocks.List() {
		_, ok = blocks[block.(string)]
		if !ok {
			return false
		}
	}

	return true
}
