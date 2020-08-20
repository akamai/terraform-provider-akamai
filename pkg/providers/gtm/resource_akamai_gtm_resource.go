package gtm

import (
	"errors"
	"fmt"
	"log"

	gtm "github.com/akamai/AkamaiOPEN-edgegrid-golang/configgtm-v1_4"
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
			"resource_instance": &schema.Schema{
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

// utility func to parse Terraform property resource id
func parseResourceResourceId(id string) (string, string, error) {

	return parseResourceStringId(id)

}

// Create a new GTM Resource
func resourceGTMv1ResourceCreate(d *schema.ResourceData, meta interface{}) error {

	domain := d.Get("domain").(string)

	log.Printf("[INFO] [Akamai GTM] Creating resource [%s] in domain [%s]", d.Get("name").(string), domain)
	newRsrc := populateNewResourceObject(d)
	log.Printf("[DEBUG] [Akamai GTMv1] Proposed New Resource: [%v]", newRsrc)
	cStatus, err := newRsrc.Create(domain)
	if err != nil {
		log.Printf("[ERROR] ResourceCreate failed: %s", err.Error())
		return err
	}
	log.Printf("[DEBUG] [Akamai GTMv1] Resource Create status:")
	log.Printf("[DEBUG] [Akamai GTMv1] %v", cStatus.Status)
	if cStatus.Status.PropagationStatus == "DENIED" {
		return errors.New(cStatus.Status.Message)
	}
	if d.Get("wait_on_complete").(bool) {
		done, err := waitForCompletion(domain)
		if done {
			log.Printf("[INFO] [Akamai GTMv1] Resource Create completed")
		} else {
			if err == nil {
				log.Printf("[INFO] [Akamai GTMv1] Resource Create pending")
			} else {
				log.Printf("[WARNING] [Akamai GTMv1] Resource Create failed [%s]", err.Error())
				return err
			}
		}

	}

	// Give terraform the ID. Format domain:resource
	resourceId := fmt.Sprintf("%s:%s", domain, cStatus.Resource.Name)
	log.Printf("[DEBUG] [Akamai GTMv1] Generated Resource Resource Id: %s", resourceId)
	d.SetId(resourceId)
	return resourceGTMv1ResourceRead(d, meta)

}

// read resource. updates state with entire API result configuration.
func resourceGTMv1ResourceRead(d *schema.ResourceData, meta interface{}) error {

	log.Printf("[DEBUG] [Akamai GTMv1] READ")
	log.Printf("[DEBUG] Reading [Akamai GTMv1] Resource: %s", d.Id())
	// retrieve the property and domain
	domain, resource, err := parseResourceResourceId(d.Id())
	if err != nil {
		return errors.New("Invalid resource resource Id")
	}
	rsrc, err := gtm.GetResource(resource, domain)
	if err != nil {
		log.Printf("[ERROR] ResourceRead failed: %s", err.Error())
		return err
	}
	populateTerraformResourceState(d, rsrc)
	log.Printf("[DEBUG] [Akamai GTMv1] READ %v", rsrc)
	return nil
}

// Update GTM Resource
func resourceGTMv1ResourceUpdate(d *schema.ResourceData, meta interface{}) error {

	log.Printf("[DEBUG] [Akamai GTMv1] UPDATE")
	// pull domain and resource out of id
	domain, resource, err := parseResourceResourceId(d.Id())
	if err != nil {
		return errors.New("Invalid resource Id")
	}
	// Get existing property
	existRsrc, err := gtm.GetResource(resource, domain)
	if err != nil {
		log.Printf("[ERROR] ResourceUpdate failed: %s", err.Error())
		return err
	}
	log.Printf("[DEBUG] Updating [Akamai GTMv1] Resource BEFORE: %v", existRsrc)
	populateResourceObject(d, existRsrc)
	log.Printf("[DEBUG] Updating [Akamai GTMv1] Resource PROPOSED: %v", existRsrc)
	uStat, err := existRsrc.Update(domain)
	if err != nil {
		log.Printf("[ERROR] ResourceUpdate failed: %s", err.Error())
		return err
	}
	log.Printf("[DEBUG] [Akamai GTMv1] Resource Update  status:")
	log.Printf("[DEBUG] [Akamai GTMv1] %v", uStat)
	if uStat.PropagationStatus == "DENIED" {
		return errors.New(uStat.Message)
	}
	if d.Get("wait_on_complete").(bool) {
		done, err := waitForCompletion(domain)
		if done {
			log.Printf("[INFO] [Akamai GTMv1] Resource update completed")
		} else {
			if err == nil {
				log.Printf("[INFO] [Akamai GTMv1] Resource update pending")
			} else {
				log.Printf("[WARNING] [Akamai GTMv1] Resource update failed [%s]", err.Error())
				return err
			}
		}

	}

	return resourceGTMv1ResourceRead(d, meta)
}

// Import GTM Resource.
func resourceGTMv1ResourceImport(d *schema.ResourceData, meta interface{}) ([]*schema.ResourceData, error) {

	log.Printf("[INFO] [Akamai GTM] Resource [%s] Import", d.Id())
	// pull domain and resource out of resource id
	domain, resource, err := parseResourceResourceId(d.Id())
	if err != nil {
		return []*schema.ResourceData{d}, errors.New("Invalid resource resource Id")
	}
	rsrc, err := gtm.GetResource(resource, domain)
	if err != nil {
		return nil, err
	}
	d.Set("domain", domain)
	d.Set("wait_on_complete", true)
	populateTerraformResourceState(d, rsrc)

	// use same Id as passed in
	log.Printf("[INFO] [Akamai GTM] Resource [%s] [%s] Imported", d.Id(), d.Get("name"))
	return []*schema.ResourceData{d}, nil
}

// Delete GTM Resource.
func resourceGTMv1ResourceDelete(d *schema.ResourceData, meta interface{}) error {

	log.Printf("[DEBUG] [Akamai GTMv1] DELETE")
	log.Printf("[DEBUG] Deleting [Akamai GTMv1] Resource: %s", d.Id())
	// Get existing resource
	domain, resource, err := parseResourceResourceId(d.Id())
	if err != nil {
		return errors.New("Invalid resource Id")
	}
	existRsrc, err := gtm.GetResource(resource, domain)
	if err != nil {
		log.Printf("[ERROR] ResourceDelete failed: %s", err.Error())
		return err
	}
	log.Printf("[DEBUG] Deleting [Akamai GTMv1] Resource: %v", existRsrc)
	uStat, err := existRsrc.Delete(domain)
	if err != nil {
		log.Printf("[ERROR] ResourceDelete failed: %s", err.Error())
		return err
	}
	log.Printf("[DEBUG] [Akamai GTMv1] Resource Delete status:")
	log.Printf("[DEBUG] [Akamai GTMv1] %v", uStat)
	if uStat.PropagationStatus == "DENIED" {
		return errors.New(uStat.Message)
	}
	if d.Get("wait_on_complete").(bool) {
		done, err := waitForCompletion(domain)
		if done {
			log.Printf("[INFO] [Akamai GTMv1] Resource delete completed")
		} else {
			if err == nil {
				log.Printf("[INFO] [Akamai GTMv1] Resource delete pending")
			} else {
				log.Printf("[WARNING] [Akamai GTMv1] Resource delete failed [%s]", err.Error())
				return err
			}
		}

	}

	// if succcessful ....
	d.SetId("")
	return nil
}

// Test GTM Resource existance
func resourceGTMv1ResourceExists(d *schema.ResourceData, meta interface{}) (bool, error) {

	log.Printf("[DEBUG] [Akamai GTMv1] Exists")
	// pull domain and resource out of resource id
	domain, resource, err := parseResourceResourceId(d.Id())
	if err != nil {
		return false, errors.New("Invalid resource resource Id")
	}
	log.Printf("[DEBUG] [Akamai GTMv1] Searching for existing resource [%s] in domain %s", resource, domain)
	rsrc, err := gtm.GetResource(resource, domain)
	return rsrc != nil, err
}

// Create and populate a new resource object from resource data
func populateNewResourceObject(d *schema.ResourceData) *gtm.Resource {

	rsrcObj := gtm.NewResource(d.Get("name").(string))
	rsrcObj.ResourceInstances = make([]*gtm.ResourceInstance, 0)
	populateResourceObject(d, rsrcObj)

	return rsrcObj

}

// Populate existing resource object from resource data
func populateResourceObject(d *schema.ResourceData, rsrc *gtm.Resource) {

	if v, ok := d.GetOk("name"); ok {
		rsrc.Name = v.(string)
	}
	if v, ok := d.GetOk("type"); ok {
		rsrc.Type = v.(string)
	}
	if v, ok := d.GetOk("host_header"); ok {
		rsrc.HostHeader = v.(string)
	} else if d.HasChange("host_header") {
		rsrc.HostHeader = v.(string)
	}
	if v, ok := d.GetOk("least_squares_decay"); ok {
		rsrc.LeastSquaresDecay = v.(float64)
	} else if d.HasChange("least_squares_decay") {
		rsrc.LeastSquaresDecay = v.(float64)
	}
	if v, ok := d.GetOk("upper_bound"); ok {
		rsrc.UpperBound = v.(int)
	} else if d.HasChange("upper_bound") {
		rsrc.UpperBound = v.(int)
	}
	if v, ok := d.GetOk("description"); ok {
		rsrc.Description = v.(string)
	} else if d.HasChange("description") {
		rsrc.Description = v.(string)
	}
	if v, ok := d.GetOk("leader_string"); ok {
		rsrc.LeaderString = v.(string)
	} else if d.HasChange("leader_string") {
		rsrc.LeaderString = v.(string)
	}
	if v, ok := d.GetOk("constrained_property"); ok {
		rsrc.ConstrainedProperty = v.(string)
	} else if d.HasChange("constrained_property") {
		rsrc.ConstrainedProperty = v.(string)
	}
	if v, ok := d.GetOk("aggregation_type"); ok {
		rsrc.AggregationType = v.(string)
	}
	if v, ok := d.GetOk("load_imbalance_percentage"); ok {
		rsrc.LoadImbalancePercentage = v.(float64)
	} else if d.HasChange("load_imbalance_percentage") {
		rsrc.LoadImbalancePercentage = v.(float64)
	}
	if v, ok := d.GetOk("max_u_multiplicative_increment"); ok {
		rsrc.MaxUMultiplicativeIncrement = v.(float64)
	} else if d.HasChange("max_u_multiplicative_increment") {
		rsrc.MaxUMultiplicativeIncrement = v.(float64)
	}
	if v, ok := d.GetOk("decay_rate"); ok {
		rsrc.DecayRate = v.(float64)
	} else if d.HasChange("decay_rate") {
		rsrc.DecayRate = v.(float64)
	}
	if _, ok := d.GetOk("resource_instance"); ok {
		populateResourceInstancesObject(d, rsrc)
	} else if d.HasChange("resource_instance") {
		rsrc.ResourceInstances = make([]*gtm.ResourceInstance, 0)
	}

	return

}

// Populate Terraform state from provided Resource object
func populateTerraformResourceState(d *schema.ResourceData, rsrc *gtm.Resource) {

	log.Printf("[DEBUG] [Akamai GTMv1] Entering populateTerraformResourceState")
	// walk thru all state elements
	d.Set("name", rsrc.Name)
	d.Set("type", rsrc.Type)
	d.Set("host_header", rsrc.HostHeader)
	d.Set("least_squares_decay", rsrc.LeastSquaresDecay)
	d.Set("description", rsrc.Description)
	d.Set("leader_string", rsrc.LeaderString)
	d.Set("constrained_property", rsrc.ConstrainedProperty)
	d.Set("aggregation_type", rsrc.AggregationType)
	d.Set("load_imbalance_percentage", rsrc.LoadImbalancePercentage)
	d.Set("upper_bound", rsrc.UpperBound)
	d.Set("max_u_multiplicative_increment", rsrc.MaxUMultiplicativeIncrement)
	d.Set("decay_rate", rsrc.DecayRate)
	populateTerraformResourceInstancesState(d, rsrc)

	return

}

// create and populate GTM Resource ResourceInstances object
func populateResourceInstancesObject(d *schema.ResourceData, rsrc *gtm.Resource) {

	// pull apart List
	rsrcInstanceList := d.Get("resource_instance").([]interface{})
	if rsrcInstanceList != nil {
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
func populateTerraformResourceInstancesState(d *schema.ResourceData, rsrc *gtm.Resource) {

	riObjectInventory := make(map[int]*gtm.ResourceInstance, len(rsrc.ResourceInstances))
	if len(rsrc.ResourceInstances) > 0 {
		for _, riObj := range rsrc.ResourceInstances {
			riObjectInventory[riObj.DatacenterId] = riObj
		}
	}
	riStateList := d.Get("resource_instance").([]interface{})
	for _, riMap := range riStateList {
		ri := riMap.(map[string]interface{})
		objIndex := ri["datacenter_id"].(int)
		riObject := riObjectInventory[objIndex]
		if riObject == nil {
			log.Printf("[WARNING] [Akamai GTMv1] Resource_instance %d NOT FOUND in returned GTM Object", ri["datacenter_id"])
			continue
		}
		ri["use_default_load_object"] = riObject.UseDefaultLoadObject
		ri["load_object"] = riObject.LoadObject.LoadObject
		ri["load_object_port"] = riObject.LoadObjectPort
		if ri["load_servers"] != nil {
			ri["load_servers"] = reconcileTerraformLists(ri["load_servers"].([]interface{}), convertStringToInterfaceList(riObject.LoadServers))
		} else {
			ri["load_servers"] = riObject.LoadServers
		}
		// remove object
		delete(riObjectInventory, objIndex)
	}
	if len(riObjectInventory) > 0 {
		log.Printf("[DEBUG] [Akamai GTMv1] Resource_instance objects left...")
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
	d.Set("resource_instance", riStateList)

}
