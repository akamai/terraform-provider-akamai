package gtm

import (
	"fmt"

	gtm "github.com/akamai/AkamaiOPEN-edgegrid-golang/configgtm-v1_4"
	"github.com/akamai/terraform-provider-akamai/v2/pkg/akamai"
	"github.com/akamai/terraform-provider-akamai/v2/pkg/tools"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/schema"
)

func resourceGTMv1Resource() *schema.Resource {
	return &schema.Resource{
		Create: resourceGTMv1ResourceCreate,
		Read:   resourceGTMv1ResourceRead,
		Update: resourceGTMv1ResourceUpdate,
		Delete: resourceGTMv1ResourceDelete,
		Exists: resourceGTMv1ResourceExists,
		Importer: &schema.ResourceImporter{
			State: resourceGTMv1ResourceImport,
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
			"host_header": {
				Type:     schema.TypeString,
				Optional: true,
			},
			"type": {
				Type:     schema.TypeString,
				Required: true,
			},
			"aggregation_type": {
				Type:     schema.TypeString,
				Required: true,
			},
			"least_squares_decay": {
				Type:     schema.TypeFloat,
				Optional: true,
			},
			"upper_bound": {
				Type:     schema.TypeInt,
				Optional: true,
			},
			"description": {
				Type:     schema.TypeString,
				Optional: true,
			},
			"leader_string": {
				Type:     schema.TypeString,
				Optional: true,
			},
			"constrained_property": {
				Type:     schema.TypeString,
				Optional: true,
			},
			"load_imbalance_percentage": {
				Type:     schema.TypeFloat,
				Optional: true,
			},
			"max_u_multiplicative_increment": {
				Type:     schema.TypeFloat,
				Optional: true,
			},
			"decay_rate": {
				Type:     schema.TypeFloat,
				Optional: true,
			},
			"resource_instance": {
				Type:     schema.TypeList,
				Optional: true,
				Elem: &schema.Resource{
					Schema: map[string]*schema.Schema{
						"datacenter_id": {
							Type:     schema.TypeInt,
							Required: true,
						},
						"use_default_load_object": {
							Type:     schema.TypeBool,
							Optional: true,
						},
						"load_object": {
							Type:     schema.TypeString,
							Optional: true,
						},
						"load_servers": {
							Type:     schema.TypeList,
							Elem:     &schema.Schema{Type: schema.TypeString},
							Optional: true,
						},
						"load_object_port": {
							Type:     schema.TypeInt,
							Optional: true,
						},
					},
				},
			},
		},
	}
}

// Create a new GTM Resource
func resourceGTMv1ResourceCreate(d *schema.ResourceData, m interface{}) error {
	meta := akamai.Meta(m)
	logger := meta.Log("Akamai GTMv1", "resourceGTMv1ResourceCreate")

	domain, err := tools.GetStringValue("domain", d)
	if err != nil {
		return err
	}

	name, _ := tools.GetStringValue("name", d)

	logger.Infof("Creating resource [%s] in domain [%s]", name, domain)
	newRsrc := populateNewResourceObject(d)
	logger.Debugf("Proposed New Resource: [%v]", newRsrc)
	cStatus, err := newRsrc.Create(domain)
	if err != nil {
		logger.Errorf("ResourceCreate failed: %s", err.Error())
		return err
	}
	logger.Debugf("Resource Create status:")
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
			logger.Infof("Resource Create completed")
		} else {
			if err == nil {
				logger.Infof("Resource Create pending")
			} else {
				logger.Errorf("Resource Create failed [%s]", err.Error())
				return err
			}
		}

	}

	// Give terraform the ID. Format domain:resource
	resourceId := fmt.Sprintf("%s:%s", domain, cStatus.Resource.Name)
	logger.Debugf("Generated Resource Resource Id: %s", resourceId)
	d.SetId(resourceId)
	return resourceGTMv1ResourceRead(d, m)

}

// read resource. updates state with entire API result configuration.
func resourceGTMv1ResourceRead(d *schema.ResourceData, m interface{}) error {
	meta := akamai.Meta(m)
	logger := meta.Log("Akamai GTMv1", "resourceGTMv1ResourceRead")

	logger.Debugf("READ")
	logger.Debugf("Reading Resource: %s", d.Id())
	// retrieve the property and domain
	domain, resource, err := parseResourceStringId(d.Id())
	if err != nil {
		return fmt.Errorf("invalid resource resource ID")
	}
	rsrc, err := gtm.GetResource(resource, domain)
	if err != nil {
		logger.Errorf("ResourceRead failed: %s", err.Error())
		return err
	}
	populateTerraformResourceState(d, rsrc, m)
	logger.Debugf("READ %v", rsrc)
	return nil
}

// Update GTM Resource
func resourceGTMv1ResourceUpdate(d *schema.ResourceData, m interface{}) error {
	meta := akamai.Meta(m)
	logger := meta.Log("Akamai GTMv1", "resourceGTMv1ResourceUpdate")

	logger.Debugf("UPDATE")
	// pull domain and resource out of id
	domain, resource, err := parseResourceStringId(d.Id())
	if err != nil {
		return fmt.Errorf("invalid resource ID")
	}
	// Get existing property
	existRsrc, err := gtm.GetResource(resource, domain)
	if err != nil {
		logger.Errorf("ResourceUpdate failed: %s", err.Error())
		return err
	}
	logger.Debugf("Updating Resource BEFORE: %v", existRsrc)
	populateResourceObject(d, existRsrc)
	logger.Debugf("Updating Resource PROPOSED: %v", existRsrc)
	uStat, err := existRsrc.Update(domain)
	if err != nil {
		logger.Errorf("ResourceUpdate failed: %s", err.Error())
		return err
	}
	logger.Debugf("Resource Update  status:")
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
			logger.Infof("Resource update completed")
		} else {
			if err == nil {
				logger.Infof("Resource update pending")
			} else {
				logger.Warnf("Resource update failed [%s]", err.Error())
				return err
			}
		}

	}

	return resourceGTMv1ResourceRead(d, m)
}

// Import GTM Resource.
func resourceGTMv1ResourceImport(d *schema.ResourceData, m interface{}) ([]*schema.ResourceData, error) {
	meta := akamai.Meta(m)
	logger := meta.Log("Akamai GTMv1", "resourceGTMv1ResourceImport")

	logger.Infof("Resource [%s] Import", d.Id())
	// pull domain and resource out of resource id
	domain, resource, err := parseResourceStringId(d.Id())
	if err != nil {
		return []*schema.ResourceData{d}, fmt.Errorf("invalid resource resource ID")
	}
	rsrc, err := gtm.GetResource(resource, domain)
	if err != nil {
		return nil, err
	}
	_ = d.Set("domain", domain)
	_ = d.Set("wait_on_complete", true)
	populateTerraformResourceState(d, rsrc, m)

	// use same Id as passed in
	name, _ := tools.GetStringValue("name", d)
	logger.Infof("Resource [%s] [%s] Imported", d.Id(), name)
	return []*schema.ResourceData{d}, nil
}

// Delete GTM Resource.
func resourceGTMv1ResourceDelete(d *schema.ResourceData, m interface{}) error {
	meta := akamai.Meta(m)
	logger := meta.Log("Akamai GTMv1", "resourceGTMv1ResourceDelete")

	logger.Debugf("DELETE")
	logger.Debugf("Deleting Resource: %s", d.Id())
	// Get existing resource
	domain, resource, err := parseResourceStringId(d.Id())
	if err != nil {
		return fmt.Errorf("invalid resource ID")
	}
	existRsrc, err := gtm.GetResource(resource, domain)
	if err != nil {
		logger.Errorf("ResourceDelete failed: %s", err.Error())
		return err
	}
	logger.Debugf("Deleting Resource: %v", existRsrc)
	uStat, err := existRsrc.Delete(domain)
	if err != nil {
		logger.Errorf("ResourceDelete failed: %s", err.Error())
		return err
	}
	logger.Debugf("Resource Delete status:")
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
			logger.Infof("Resource delete completed")
		} else {
			if err == nil {
				logger.Infof("Resource delete pending")
			} else {
				logger.Errorf("Resource delete failed [%s]", err.Error())
				return err
			}
		}

	}

	// if successful ....
	d.SetId("")
	return nil
}

// Test GTM Resource existence
func resourceGTMv1ResourceExists(d *schema.ResourceData, m interface{}) (bool, error) {
	meta := akamai.Meta(m)
	logger := meta.Log("Akamai GTMv1", "resourceGTMv1ResourceExists")

	logger.Debugf("Exists")
	// pull domain and resource out of resource id
	domain, resource, err := parseResourceStringId(d.Id())
	if err != nil {
		return false, fmt.Errorf("invalid resource resource ID")
	}
	logger.Debugf("Searching for existing resource [%s] in domain %s", resource, domain)
	rsrc, err := gtm.GetResource(resource, domain)
	return rsrc != nil, err
}

// Create and populate a new resource object from resource data
func populateNewResourceObject(d *schema.ResourceData) *gtm.Resource {

	name, _ := tools.GetStringValue("name", d)
	rsrcObj := gtm.NewResource(name)
	rsrcObj.ResourceInstances = make([]*gtm.ResourceInstance, 0)
	populateResourceObject(d, rsrcObj)

	return rsrcObj

}

// Populate existing resource object from resource data
func populateResourceObject(d *schema.ResourceData, rsrc *gtm.Resource) {

	if v, err := tools.GetStringValue("name", d); err == nil {
		rsrc.Name = v
	}
	if v, err := tools.GetStringValue("type", d); err == nil {
		rsrc.Type = v
	}
	if v, err := tools.GetStringValue("host_header", d); err == nil || d.HasChange("host_header") {
		rsrc.HostHeader = v
	}
	if v, err := tools.GetFloat64Value("least_squares_decay", d); err == nil || d.HasChange("least_squares_decay") {
		rsrc.LeastSquaresDecay = v
	}
	if v, err := tools.GetIntValue("upper_bound", d); err == nil || d.HasChange("upper_bound") {
		rsrc.UpperBound = v
	}
	if v, err := tools.GetStringValue("description", d); err == nil || d.HasChange("description") {
		rsrc.Description = v
	}
	if v, err := tools.GetStringValue("leader_string", d); err == nil || d.HasChange("leader_string") {
		rsrc.LeaderString = v
	}
	if v, err := tools.GetStringValue("constrained_property", d); err == nil || d.HasChange("constrained_property") {
		rsrc.ConstrainedProperty = v
	}
	if v, err := tools.GetStringValue("aggregation_type", d); err == nil {
		rsrc.AggregationType = v
	}
	if v, err := tools.GetFloat64Value("load_imbalance_percentage", d); err == nil || d.HasChange("load_imbalance_percentage") {
		rsrc.LoadImbalancePercentage = v
	}
	if v, err := tools.GetFloat64Value("max_u_multiplicative_increment", d); err == nil || d.HasChange("max_u_multiplicative_increment") {
		rsrc.MaxUMultiplicativeIncrement = v
	}
	if v, err := tools.GetFloat64Value("decay_rate", d); err == nil || d.HasChange("decay_rate") {
		rsrc.DecayRate = v
	}
	if _, ok := d.GetOk("resource_instance"); ok {
		populateResourceInstancesObject(d, rsrc)
	} else if d.HasChange("resource_instance") {
		rsrc.ResourceInstances = make([]*gtm.ResourceInstance, 0)
	}
}

// Populate Terraform state from provided Resource object
func populateTerraformResourceState(d *schema.ResourceData, rsrc *gtm.Resource, m interface{}) {
	meta := akamai.Meta(m)
	logger := meta.Log("Akamai GTMv1", "populateTerraformResourceState")

	logger.Debugf("Entering populateTerraformResourceState")
	// walk through all state elements
	for stateKey, stateValue := range map[string]interface{}{
		"name":                           rsrc.Name,
		"type":                           rsrc.Type,
		"host_header":                    rsrc.HostHeader,
		"least_squares_decay":            rsrc.LeastSquaresDecay,
		"description":                    rsrc.Description,
		"leader_string":                  rsrc.LeaderString,
		"constrained_property":           rsrc.ConstrainedProperty,
		"aggregation_type":               rsrc.AggregationType,
		"load_imbalance_percentage":      rsrc.LoadImbalancePercentage,
		"upper_bound":                    rsrc.UpperBound,
		"max_u_multiplicative_increment": rsrc.MaxUMultiplicativeIncrement,
		"decay_rate":                     rsrc.DecayRate,
	} {
		err := d.Set(stateKey, stateValue)
		if err != nil {
			logger.Errorf("populateTerraformResourceState failed: %s", err.Error())
		}
	}
	populateTerraformResourceInstancesState(d, rsrc, m)
}

// create and populate GTM Resource ResourceInstances object
func populateResourceInstancesObject(d *schema.ResourceData, rsrc *gtm.Resource) {

	// pull apart List
	rsrcInstanceList, err := tools.GetInterfaceArrayValue("resource_instance", d)
	if err == nil {
		rsrcInstanceObjList := make([]*gtm.ResourceInstance, len(rsrcInstanceList)) // create new object list
		for i, v := range rsrcInstanceList {
			riMap := v.(map[string]interface{})
			rsrcInstance := rsrc.NewResourceInstance(riMap["datacenter_id"].(int)) // create new object
			rsrcInstance.UseDefaultLoadObject = riMap["use_default_load_object"].(bool)
			if riMap["load_servers"] != nil {
				ls := make([]string, len(riMap["load_servers"].([]interface{})))
				for i, sl := range riMap["load_servers"].([]interface{}) {
					ls[i] = sl.(string)
				}
				rsrcInstance.LoadServers = ls
			}
			rsrcInstance.LoadObject.LoadObject = riMap["load_object"].(string)
			rsrcInstance.LoadObjectPort = riMap["load_object_port"].(int)
			rsrcInstanceObjList[i] = rsrcInstance
		}
		rsrc.ResourceInstances = rsrcInstanceObjList
	}
}

// create and populate Terraform resource_instances schema
func populateTerraformResourceInstancesState(d *schema.ResourceData, rsrc *gtm.Resource, m interface{}) {
	meta := akamai.Meta(m)
	logger := meta.Log("Akamai GTMv1", "populateTerraformResourceInstancesState")

	riObjectInventory := make(map[int]*gtm.ResourceInstance, len(rsrc.ResourceInstances))
	if len(rsrc.ResourceInstances) > 0 {
		for _, riObj := range rsrc.ResourceInstances {
			riObjectInventory[riObj.DatacenterId] = riObj
		}
	}
	riStateList, _ := tools.GetInterfaceArrayValue("resource_instance", d)
	for _, riMap := range riStateList {
		ri := riMap.(map[string]interface{})
		objIndex := ri["datacenter_id"].(int)
		riObject := riObjectInventory[objIndex]
		if riObject == nil {
			logger.Warnf("Resource_instance %d NOT FOUND in returned GTM Object", ri["datacenter_id"])
			continue
		}
		ri["use_default_load_object"] = riObject.UseDefaultLoadObject
		ri["load_object"] = riObject.LoadObject.LoadObject
		ri["load_object_port"] = riObject.LoadObjectPort
		if ri["load_servers"] != nil {
			ri["load_servers"] = reconcileTerraformLists(ri["load_servers"].([]interface{}), convertStringToInterfaceList(riObject.LoadServers, m), m)
		} else {
			ri["load_servers"] = riObject.LoadServers
		}
		// remove object
		delete(riObjectInventory, objIndex)
	}
	if len(riObjectInventory) > 0 {
		logger.Debugf("Resource_instance objects left...")
		// Objects not in the state yet. Add. Unfortunately, they'll likely not align with instance indices in the config
		for _, mriObj := range riObjectInventory {
			riNew := make(map[string]interface{})
			riNew["datacenter_id"] = mriObj.DatacenterId
			riNew["use_default_load_object"] = mriObj.UseDefaultLoadObject
			riNew["load_object"] = mriObj.LoadObject.LoadObject
			riNew["load_object_port"] = mriObj.LoadObjectPort
			riNew["load_servers"] = mriObj.LoadServers
			riStateList = append(riStateList, riNew)
		}
	}
	_ = d.Set("resource_instance", riStateList)

}
