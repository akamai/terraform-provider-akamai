package akamai

import (
	"fmt"
	"github.com/akamai/AkamaiOPEN-edgegrid-golang/papi-v1"
	"github.com/hashicorp/terraform/helper/schema"
	"log"
	"strings"
)

func resourcePropertyVariable() *schema.Resource {
	return &schema.Resource{
		Create: resourcePropertyVariableCreate,
		Read:   resourcePropertyVariableRead,
		Update: resourcePropertyVariableUpdate,
		Delete: resourcePropertyVariableDelete,
		Exists: resourcePropertyVariableExists,
		Importer: &schema.ResourceImporter{
			State: resourcePropertyVariableImport,
		},
		Schema: akamaiPropertyVariableSchema,
	}
}

var akamaiPropertyVariableSchema = map[string]*schema.Schema{
	"name": {
		Type:     schema.TypeString,
		Required: true,
	},
	"fqname": {
		Type:     schema.TypeString,
		Required: true,
	},
	"hidden": {
		Type:     schema.TypeBool,
		Required: true,
	},
	"sensitive": {
		Type:     schema.TypeBool,
		Required: true,
	},
	"value": {
		Type:     schema.TypeString,
		Optional: true,
	},
	"description": {
		Type:     schema.TypeString,
		Optional: true,
	},
}

func resourcePropertyVariableCreate(d *schema.ResourceData, meta interface{}) error {

	var name string
	var value string
	var description string
	var hidden bool
	var sensitive bool
	var fqname string

	_, ok := d.GetOk("name")
	if ok {
		name = d.Get("name").(string)
	}
	_, ok = d.GetOk("value")
	if ok {
		value = d.Get("value").(string)
	}
	_, ok = d.GetOk("desription")
	if ok {
		description = d.Get("description").(string)
	}
	_, ok = d.GetOk("hidden")
	if ok {
		hidden = d.Get("hidden").(bool)
	}
	_, ok = d.GetOk("sensitive")
	if ok {
		sensitive = d.Get("sensitive").(bool)
	}
	_, ok = d.GetOk("fqname")
	if ok {
		fqname = d.Get("fqname").(string)
	}

	d.Set("name", name)
	d.Set("value", value)
	d.Set("description", description)
	d.Set("hidden", hidden)
	d.Set("sensitive", sensitive)
	d.Set("fqname", fqname)
	d.SetId(fmt.Sprintf("%s-%s-%s-%s", name, value, description, fqname))

	log.Println("[DEBUG] Done")
	return nil
}

func createPropertyVariable(contract *papi.Contract, group *papi.Group, product *papi.Product, cloneFrom *papi.ClonePropertyFrom, d *schema.ResourceData) (*papi.Property, error) {
	log.Println("[DEBUG] Creating property")

	property, err := group.NewProperty(contract)
	if err != nil {
		return nil, err
	}

	property.ProductID = product.ProductID
	property.PropertyName = d.Get("name").(string)
	if cloneFrom != nil {
		property.CloneFrom = cloneFrom
	}

	if ruleFormat, ok := d.GetOk("rule_format"); ok {
		property.RuleFormat = ruleFormat.(string)
	} else {
		ruleFormats := papi.NewRuleFormats()
		property.RuleFormat, err = ruleFormats.GetLatest()
		if err != nil {
			return nil, err
		}
	}

	err = property.Save()
	if err != nil {
		return nil, err
	}

	log.Printf("[DEBUG] Property created: %s\n", property.PropertyID)
	return property, nil
}

func resourcePropertyVariableDelete(d *schema.ResourceData, meta interface{}) error {
	log.Printf("[DEBUG] DELETING")

	d.SetId("")

	log.Println("[DEBUG] Done")

	return nil
}

func resourcePropertyVariableImport(d *schema.ResourceData, meta interface{}) ([]*schema.ResourceData, error) {
	resourceID := d.Id()
	propertyID := resourceID

	if !strings.HasPrefix(resourceID, "prp_") {
		for _, searchKey := range []papi.SearchKey{papi.SearchByPropertyName, papi.SearchByHostname, papi.SearchByEdgeHostname} {
			results, err := papi.Search(searchKey, resourceID)
			if err != nil {
				continue
			}

			if results != nil && len(results.Versions.Items) > 0 {
				propertyID = results.Versions.Items[0].PropertyID
				break
			}
		}
	}

	property := papi.NewProperty(papi.NewProperties())
	property.PropertyID = propertyID
	e := property.GetProperty()
	if e != nil {
		return nil, e
	}

	d.Set("account", property.AccountID)
	d.Set("contract", property.ContractID)
	d.Set("group", property.GroupID)

	d.Set("name", property.PropertyName)
	d.Set("version", property.LatestVersion)
	d.SetId(property.PropertyID)

	return []*schema.ResourceData{d}, nil
}

func resourcePropertyVariableExists(d *schema.ResourceData, meta interface{}) (bool, error) {

	return true, nil
}

func resourcePropertyVariableRead(d *schema.ResourceData, meta interface{}) error {

	return nil
}

func resourcePropertyVariableUpdate(d *schema.ResourceData, meta interface{}) error {
	log.Printf("[DEBUG] UPDATING")

	var name string
	var value string
	var description string
	var hidden bool
	var sensitive bool

	_, ok := d.GetOk("name")
	if ok {
		name = d.Get("name").(string)
	}
	_, ok = d.GetOk("value")
	if ok {
		value = d.Get("value").(string)
	}
	_, ok = d.GetOk("description")
	if ok {
		description = d.Get("description").(string)
	}
	_, ok = d.GetOk("hidden")
	if ok {
		hidden = d.Get("hidden").(bool)
	}
	_, ok = d.GetOk("sensitive")
	if ok {
		sensitive = d.Get("sensitive").(bool)
	}

	d.Set("name", name)
	d.Set("value", value)
	d.Set("description", description)
	d.Set("hidden", hidden)
	d.Set("sensitive", sensitive)

	log.Println("[DEBUG] Done")
	return nil
}
