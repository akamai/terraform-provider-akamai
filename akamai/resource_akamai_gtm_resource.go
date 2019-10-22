package akamai

import (
	"encoding/json"
	"errors"
	"fmt"
	gtmv1_3 "github.com/akamai/AkamaiOPEN-edgegrid-golang/configgtm-v1_3"
	"github.com/hashicorp/terraform/helper/schema"
	"log"
)

func resourceGTMv1_3Resource() *schema.Resource {
	return &schema.Resource{
		Create: resourceGTMv1_3ResourceCreate,
		Read:   resourceGTMv1_3ResourceRead,
		Update: resourceGTMv1_3ResourceUpdate,
		Delete: resourceGTMv1_3ResourceDelete,
		Exists: resourceGTMv1_3ResourceExists,
		Importer: &schema.ResourceImporter{
			State: resourceGTMv1_3ResourceImport,
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
				Type:     schema.TypeInt,
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
			"load_imbalance_percent": {
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
			"resource_instances": &schema.Schema{
				Type:       schema.TypeList,
				Optional:   true,
				ConfigMode: schema.SchemaConfigModeAttr,
				Elem: &schema.Resource{
					Schema: map[string]*schema.Schema{
						"datacenter_id": {
							Type:     schema.TypeInt,
							Optional: true,
						},
						"use_default_load_object": {
							Type:     schema.TypeBool,
							Optional: true,
						},
						"load_object": {
							Type:     schema.TypeString,
							Optional: true,
						},
						"load_servers": &schema.Schema{
							Type:     schema.TypeList,
							Elem:     &schema.Schema{Type: schema.TypeString},
							Optional: true,
						},
						"load_object_port": {
							Type:     schema.TypeInt,
							Optional: true,
							Default:  "",
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
func resourceGTMv1_3ResourceCreate(d *schema.ResourceData, meta interface{}) error {

	domain := d.Get("domain").(string)

	log.Printf("[INFO] [Akamai GTM] Creating resource [%s] in domain [%s]", d.Get("name").(string), domain)
	newRsrc := populateNewResourceObject(d)
	log.Printf("[DEBUG] [Akamai GTMV1_3] Proposed New Resource: [%v]", newRsrc)
	cStatus, err := newRsrc.Create(domain)
	if err != nil {
		log.Printf("[DEBUG] [Akamai GTMV1_3] Resource Create failed: %s", err.Error())
		fmt.Println(err)
		return err
	}
	b, err := json.Marshal(cStatus.Status)
	if err != nil {
		fmt.Println(err)
		return err
	}
	fmt.Println(string(b))
	log.Printf("[DEBUG] [Akamai GTMV1_3] Resource Create status:")
	log.Printf("[DEBUG] [Akamai GTMV1_3] %v", b)

	if d.Get("wait_on_complete").(bool) {
		done, err := waitForCompletion(domain)
		if done {
			log.Printf("[INFO] [Akamai GTMV1_3] Resource Create completed")
		} else {
			if err == nil {
				log.Printf("[INFO] [Akamai GTMV1_3] Resource Create pending")
			} else {
				log.Printf("[WARNING] [Akamai GTMV1_3] Resource Create failed [%s]", err.Error())
				return err
			}
		}

	}

	// Give terraform the ID. Format domain:resource
	resourceId := fmt.Sprintf("%s:%s", domain, cStatus.Resource.Name)
	log.Printf("[DEBUG] [Akamai GTMV1_3] Generated Resource Resource Id: %s", resourceId)
	d.SetId(resourceId)
	return resourceGTMv1_3ResourceRead(d, meta)

}

// read resource. updates state with entire API result configuration.
func resourceGTMv1_3ResourceRead(d *schema.ResourceData, meta interface{}) error {

	log.Printf("[DEBUG] [Akamai GTMV1_3] READ")
	log.Printf("[DEBUG] Reading [Akamai GTMV1_3] Resource: %s", d.Id())
	// retrieve the property and domain
	domain, resource, err := parseResourceResourceId(d.Id())
	if err != nil {
		return errors.New("Invalid resource resource Id")
	}
	rsrc, err := gtmv1_3.GetResource(resource, domain)
	if err != nil {
		fmt.Println(err)
		log.Printf("[DEBUG] [Akamai GTMV1_3] Resource Read error: %s", err.Error())
		return err
	}
	populateTerraformResourceState(d, rsrc)
	log.Printf("[DEBUG] [Akamai GTMV1_3] READ %v", rsrc)
	return nil
}

// Update GTM Resource
func resourceGTMv1_3ResourceUpdate(d *schema.ResourceData, meta interface{}) error {

	log.Printf("[DEBUG] [Akamai GTMV1_3] UPDATE")
	log.Printf("[DEBUG] Updating [Akamai GTMV1_3] Resource: %s", d.Id())
	// pull domain and resource out of id
	domain, resource, err := parseResourceResourceId(d.Id())
	if err != nil {
		return errors.New("Invalid resource Id")
	}
	// Get existing property
	existRsrc, err := gtmv1_3.GetResource(resource, domain)
	if err != nil {
		fmt.Println(err.Error())
		return err
	}
	log.Printf("[DEBUG] Updating [Akamai GTMV1_3] Resource BEFORE: %v", existRsrc)
	populateResourceObject(d, existRsrc)
	log.Printf("[DEBUG] Updating [Akamai GTMV1_3] Resource PROPOSED: %v", existRsrc)
	uStat, err := existRsrc.Update(domain)
	if err != nil {
		fmt.Println(err.Error())
		return err
	}
	b, err := json.Marshal(uStat)
	if err != nil {
		fmt.Println(err.Error())
		return err
	}
	fmt.Println(string(b))
	log.Printf("[DEBUG] [Akamai GTMV1_3] Resource Update  status:")
	log.Printf("[DEBUG] [Akamai GTMV1_3] %v", b)

	if d.Get("wait_on_complete").(bool) {
		done, err := waitForCompletion(domain)
		if done {
			log.Printf("[INFO] [Akamai GTMV1_3] Resource update completed")
		} else {
			if err == nil {
				log.Printf("[INFO] [Akamai GTMV1_3] Resource update pending")
			} else {
				log.Printf("[WARNING] [Akamai GTMV1_3] Resource update failed [%s]", err.Error())
				return err
			}
		}

	}

	return resourceGTMv1_3ResourceRead(d, meta)
}

// Import GTM Resource.
func resourceGTMv1_3ResourceImport(d *schema.ResourceData, meta interface{}) ([]*schema.ResourceData, error) {

	log.Printf("[INFO] [Akamai GTM] Resource [%s] Import", d.Id())
	// pull domain and resource out of resource id
	domain, resource, err := parseResourceResourceId(d.Id())
	if err != nil {
		return []*schema.ResourceData{d}, errors.New("Invalid resource resource Id")
	}
	rsrc, err := gtmv1_3.GetResource(resource, domain)
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
func resourceGTMv1_3ResourceDelete(d *schema.ResourceData, meta interface{}) error {

	domain := d.Get("domain").(string)
	log.Printf("[DEBUG] [Akamai GTMV1_3] DELETE")
	log.Printf("[DEBUG] Deleting [Akamai GTMV1_3] Resource: %s", d.Id())
	// Get existing resource
	domain, resource, err := parseResourceResourceId(d.Id())
	if err != nil {
		return errors.New("Invalid resource Id")
	}
	existRsrc, err := gtmv1_3.GetResource(resource, domain)
	if err != nil {
		fmt.Println(err.Error())
		return err
	}
	log.Printf("[DEBUG] Deleting [Akamai GTMV1_3] Resource: %v", existRsrc)
	uStat, err := existRsrc.Delete(domain)
	if err != nil {
		fmt.Println(err)
		return err
	}
	b, err := json.Marshal(uStat)
	if err != nil {
		fmt.Println(err)
		return err
	}
	fmt.Println(string(b))
	log.Printf("[DEBUG] [Akamai GTMV1_3] Resource Delete status:")
	log.Printf("[DEBUG] [Akamai GTMV1_3] %v", b)

	if d.Get("wait_on_complete").(bool) {
		done, err := waitForCompletion(domain)
		if done {
			log.Printf("[INFO] [Akamai GTMV1_3] Resource delete completed")
		} else {
			if err == nil {
				log.Printf("[INFO] [Akamai GTMV1_3] Resource delete pending")
			} else {
				log.Printf("[WARNING] [Akamai GTMV1_3] Resource delete failed [%s]", err.Error())
				return err
			}
		}

	}

	// if succcessful ....
	d.SetId("")
	return nil
}

// Test GTM Resource existance
func resourceGTMv1_3ResourceExists(d *schema.ResourceData, meta interface{}) (bool, error) {

	log.Printf("[DEBUG] [Akamai GTMv1_3] Exists")
	// pull domain and resource out of resource id
	domain, resource, err := parseResourceResourceId(d.Id())
	if err != nil {
		return false, errors.New("Invalid resource resource Id")
	}
	log.Printf("[DEBUG] [Akamai GTMV1_3] Searching for existing resource [%s] in domain %s", resource, domain)
	rsrc, err := gtmv1_3.GetResource(resource, domain)
	return rsrc != nil, err
}

// Create and populate a new resource object from resource data
func populateNewResourceObject(d *schema.ResourceData) *gtmv1_3.Resource {

	rsrcObj := gtmv1_3.NewResource(d.Get("name").(string))
	rsrcObj.ResourceInstances = make([]*gtmv1_3.ResourceInstance, 1)
	populateResourceObject(d, rsrcObj)

	return rsrcObj

}

// Populate existing resource object from resource data
func populateResourceObject(d *schema.ResourceData, rsrc *gtmv1_3.Resource) {

	if v, ok := d.GetOk("name"); ok {
		rsrc.Name = v.(string)
	}
	if v, ok := d.GetOk("type"); ok {
		rsrc.Type = v.(string)
	}
	if v, ok := d.GetOk("host_header"); ok {
		rsrc.HostHeader = v.(string)
	}
	if v, ok := d.GetOk("least_squares_decay"); ok {
		rsrc.LeastSquaresDecay = v.(int)
	}
	if v, ok := d.GetOk("upper_bound"); ok {
		rsrc.UpperBound = v.(int)
	}
	if v, ok := d.GetOk("description"); ok {
		rsrc.Description = v.(string)
	}
	if v, ok := d.GetOk("leader_string"); ok {
		rsrc.LeaderString = v.(string)
	}
	if v, ok := d.GetOk("constrained_property"); ok {
		rsrc.ConstrainedProperty = v.(string)
	}
	if v, ok := d.GetOk("aggregation_type"); ok {
		rsrc.AggregationType = v.(string)
	}
	if v, ok := d.GetOk("load_imbalance_percent"); ok {
		rsrc.LoadImbalancePercent = v.(float64)
	}
	if v, ok := d.GetOk("max_u_multiplicative_increment"); ok {
		rsrc.MaxUMultiplicativeIncrement = v.(float64)
	}
	if v, ok := d.GetOk("decay_rate"); ok {
		rsrc.DecayRate = v.(float64)
	}
	populateResourceInstancesObject(d, rsrc)

	return

}

// Populate Terraform state from provided Resource object
func populateTerraformResourceState(d *schema.ResourceData, rsrc *gtmv1_3.Resource) {

	// walk thru all state elements
	d.Set("name", rsrc.Name)
	d.Set("type", rsrc.Type)
	d.Set("host_header", rsrc.HostHeader)
	d.Set("least_squares_decay", rsrc.LeastSquaresDecay)
	d.Set("description", rsrc.Description)
	d.Set("leader_string", rsrc.LeaderString)
	d.Set("constrained_property", rsrc.ConstrainedProperty)
	d.Set("aggregation_type", rsrc.AggregationType)
	d.Set("load_imbalance_percent", rsrc.LoadImbalancePercent)
	d.Set("upper_bound", rsrc.UpperBound)
	d.Set("max_u_multiplicative_increment", rsrc.MaxUMultiplicativeIncrement)
	d.Set("decay_rate", rsrc.DecayRate)
	populateTerraformResourceInstancesState(d, rsrc)

	return

}

// create and populate GTM Resource ResourceInstances object
func populateResourceInstancesObject(d *schema.ResourceData, rsrc *gtmv1_3.Resource) {

	// pull apart List
	ri := d.Get("resource_instances")
	if ri != nil {
		rsrcInstanceList := ri.([]interface{})
		rsrcInstanceObjList := make([]*gtmv1_3.ResourceInstance, len(rsrcInstanceList)) // create new object list
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
func populateTerraformResourceInstancesState(d *schema.ResourceData, rsrc *gtmv1_3.Resource) {

	riListNew := make([]interface{}, len(rsrc.ResourceInstances))
	for i, ri := range rsrc.ResourceInstances {
		riNew := map[string]interface{}{
			"datacenter_id":           ri.DatacenterId,
			"use_default_load_object": ri.UseDefaultLoadObject,
			"load_object":             ri.LoadObject,
			"load_object_port":        ri.LoadObjectPort,
			"load_servers":            ri.LoadServers,
		}
		riListNew[i] = riNew
	}
	d.Set("resource_instances", riListNew)

}
