package akamai

import (
	"log"
	"strings"
    "github.com/akamai/AkamaiOPEN-edgegrid-golang/jsonhooks-v1"
	"github.com/akamai/AkamaiOPEN-edgegrid-golang/papi-v1"
	"github.com/hashicorp/terraform/helper/schema"
)

func resourcePropertyVariables() *schema.Resource {
	return &schema.Resource{
		Create: resourcePropertyVariablesCreate,
		Read:   resourcePropertyVariablesRead,
		Update: resourcePropertyVariablesUpdate,
		Delete: resourcePropertyVariablesDelete,
		Exists: resourcePropertyVariablesExists,
		Importer: &schema.ResourceImporter{
			State: resourcePropertyVariablesImport,
		},
		Schema: akamaiPropertyVariablesSchema,
	}
}

var akamaiPropertyVariablesSchema = map[string]*schema.Schema{
	"variables": &schema.Schema{
		Type:     schema.TypeSet,
		Optional: true,
		Elem: &schema.Resource{
			Schema: map[string]*schema.Schema{
				"variable": &schema.Schema{
					Type:     schema.TypeSet,
					Optional: true,
					Elem: &schema.Resource{
						Schema: map[string]*schema.Schema{
							"name": {
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
							"description": {
								Type:     schema.TypeString,
								Optional: true,
							},
							
							"value": {
								Type:     schema.TypeString,
								Optional: true,
							},
						},
					},
				},
			},
		},
	},
	"json": {
		Type: schema.TypeString,
		Computed:    true,
		Description: "JSON variables representation",
	},
}

func resourcePropertyVariablesCreate(d *schema.ResourceData, meta interface{}) error {
	rule := papi.NewRule()
	log.Printf("[DEBUG] START Check for variables")
	variables, ok := d.GetOk("variables")
	if ok {
		log.Printf("[DEBUG] Check for variables  %s\n", variables)
	}

	for _, r := range variables.(*schema.Set).List() {
		variable, ok := r.(map[string]interface{})
		if ok {
			vv, ok := variable["variable"]
			if ok {
				for _, v := range vv.(*schema.Set).List() {
					variableMap, ok := v.(map[string]interface{})
					if ok {
						log.Printf("[DEBUG] Check for variables LOOP  name %s\n", variableMap["name"])
						log.Printf("[DEBUG] Check for variables LOOP  value %s\n", variableMap["value"])
						log.Printf("[DEBUG] Check for variables LOOP  description%s\n", variableMap["description"])
						log.Printf("[DEBUG] Check for variables LOOP  hidden%s\n", variableMap["hidden"])
						log.Printf("[DEBUG] Check for variables LOOP  sensitive%s\n", variableMap["sensitive"])
						log.Printf("[DEBUG] Check for variables LOOP  fqname%s\n", variableMap["fqname"])
						newVariable := papi.NewVariable()
						newVariable.Name = variableMap["name"].(string)
						newVariable.Description = variableMap["description"].(string)
						newVariable.Value = variableMap["value"].(string)
						newVariable.Hidden = variableMap["hidden"].(bool)
						newVariable.Sensitive = variableMap["sensitive"].(bool)
						rule.AddVariable(newVariable)
					}
				}

			}

		}
	}

	jsonBody, err := jsonhooks.Marshal(rule)
	if err != nil {
		return err
	}
	log.Printf("[DEBUG] Json result  %s\n", string(jsonBody))

	sha := getSHAString(string(jsonBody))
	d.Set("json", string(jsonBody))
	
	d.SetId(sha)
	log.Println("[DEBUG] Done")
	return nil
}

func resourcePropertyVariablesDelete(d *schema.ResourceData, meta interface{}) error {
	log.Printf("[DEBUG] DELETING")

	d.SetId("")

	log.Println("[DEBUG] Done")

	return nil
}

func resourcePropertyVariablesImport(d *schema.ResourceData, meta interface{}) ([]*schema.ResourceData, error) {
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

func resourcePropertyVariablesExists(d *schema.ResourceData, meta interface{}) (bool, error) {

	variables := d.Id()
	if variables != "" {
		log.Printf("[DEBUG] Check for variables  %s\n", variables)
		return true, nil
	} else {
		return true, nil

	}

}

func resourcePropertyVariablesRead(d *schema.ResourceData, meta interface{}) error {

	return nil
}

func resourcePropertyVariablesUpdate(d *schema.ResourceData, meta interface{}) error {
	log.Printf("[DEBUG] UPDATING")
	rule := papi.NewRule()
	log.Printf("[DEBUG] START Check for variables")
	variables, ok := d.GetOk("variables")
	if ok {
		log.Printf("[DEBUG] Check for variables  %s\n", variables)
		//}

		for _, r := range variables.(*schema.Set).List() {
			variable, ok := r.(map[string]interface{})
			if ok {
				vv, ok := variable["variable"]
				if ok {
					for _, v := range vv.(*schema.Set).List() {
						variableMap, ok := v.(map[string]interface{})
						if ok {
							newVariable := papi.NewVariable()
							newVariable.Name = variableMap["name"].(string)
							newVariable.Description = variableMap["description"].(string)
							newVariable.Value = variableMap["value"].(string)
							newVariable.Hidden = variableMap["hidden"].(bool)
							newVariable.Sensitive = variableMap["sensitive"].(bool)
							rule.AddVariable(newVariable)
						}
					}

				}

			}
		}

		jsonBody, err := jsonhooks.Marshal(rule)
		if err != nil {
			return err
		}
		log.Printf("[DEBUG] Json result  %s\n", string(jsonBody))

		sha := getSHAString(string(jsonBody))
		d.Set("json", string(jsonBody))
		
		d.SetId(sha)
	}
	log.Println("[DEBUG] Done")
	return nil
}
