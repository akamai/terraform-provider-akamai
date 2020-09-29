package property

import (
	"context"
	"errors"
	"fmt"
	"strings"

	"github.com/akamai/terraform-provider-akamai/v2/pkg/akamai"
	"github.com/akamai/terraform-provider-akamai/v2/pkg/tools"
	"github.com/hashicorp/terraform-plugin-sdk/v2/diag"

	"github.com/akamai/AkamaiOPEN-edgegrid-golang/jsonhooks-v1"
	"github.com/akamai/AkamaiOPEN-edgegrid-golang/papi-v1"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/schema"
)

func resourcePropertyVariables() *schema.Resource {
	return &schema.Resource{
		CreateContext: resourcePropertyVariablesCreate,
		ReadContext:   resourcePropertyVariablesRead,
		UpdateContext: resourcePropertyVariablesUpdate,
		DeleteContext: resourcePropertyVariablesDelete,
		Exists:        resourcePropertyVariablesExists,
		Importer: &schema.ResourceImporter{
			StateContext: resourcePropertyVariablesImport,
		},
		Schema: akamaiPropertyVariablesSchema,
	}
}

var akamaiPropertyVariablesSchema = map[string]*schema.Schema{
	"variables": {
		Type:     schema.TypeSet,
		Optional: true,
		Elem: &schema.Resource{
			Schema: map[string]*schema.Schema{
				"variable": {
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
		Type:        schema.TypeString,
		Computed:    true,
		Description: "JSON variables representation",
	},
}

func resourcePropertyVariablesCreate(_ context.Context, d *schema.ResourceData, m interface{}) diag.Diagnostics {
	meta := akamai.Meta(m)
	logger := meta.Log("PAPI", "resourcePropertyVariablesCreate")
	rule := papi.NewRule()
	logger.Debugf("START Check for variables")
	variables, err := tools.GetSetValue("variables", d)
	if err != nil && !errors.Is(err, tools.ErrNotFound) {
		return diag.FromErr(err)
	}
	if err == nil {
		logger.Debugf("Check for variables  %s", variables)
	}

	for _, r := range variables.List() {
		variable, ok := r.(map[string]interface{})
		if !ok {
			continue
		}
		vv, ok := variable["variable"]
		if !ok {
			continue
		}
		for _, v := range vv.(*schema.Set).List() {
			variableMap, ok := v.(map[string]interface{})
			if !ok {
				continue
			}
			logger.Debugf(" Check for variables LOOP  name %s", variableMap["name"])
			logger.Debugf("  Check for variables LOOP  value %s", variableMap["value"])
			logger.Debugf("  Check for variables LOOP  description %s", variableMap["description"])
			logger.Debugf("  Check for variables LOOP  hidden %s", variableMap["hidden"])
			logger.Debugf("  Check for variables LOOP  sensitive %s", variableMap["sensitive"])
			logger.Debugf("  Check for variables LOOP  fqname %s", variableMap["fqname"])
			newVariable := papi.NewVariable()
			name, ok := variableMap["name"].(string)
			if !ok {
				return diag.FromErr(fmt.Errorf("%w: %s, %q", tools.ErrInvalidType, "name", "string"))
			}
			description, ok := variableMap["description"].(string)
			if !ok {
				return diag.FromErr(fmt.Errorf("%w: %s, %q", tools.ErrInvalidType, "description", "string"))
			}
			value, ok := variableMap["value"].(string)
			if !ok {
				return diag.FromErr(fmt.Errorf("%w: %s, %q", tools.ErrInvalidType, "value", "string"))
			}
			hidden, ok := variableMap["hidden"].(bool)
			if !ok {
				return diag.FromErr(fmt.Errorf("%w: %s, %q", tools.ErrInvalidType, "hidden", "bool"))
			}
			sensitive, ok := variableMap["sensitive"].(bool)
			if !ok {
				return diag.FromErr(fmt.Errorf("%w: %s, %q", tools.ErrInvalidType, "sensitive", "bool"))
			}
			newVariable.Name = name
			newVariable.Description = description
			newVariable.Value = value
			newVariable.Hidden = hidden
			newVariable.Sensitive = sensitive
			rule.AddVariable(newVariable)
		}
	}

	body, err := jsonhooks.Marshal(rule)
	if err != nil {
		return diag.FromErr(err)
	}

	logger.Debugf("JSON result  %s", string(body))
	sha := tools.GetSHAString(string(body))
	if err := d.Set("json", string(body)); err != nil {
		return diag.FromErr(fmt.Errorf("%w: %s", tools.ErrValueSet, err.Error()))
	}
	d.SetId(sha)
	logger.Debugf("Done")

	return resourcePropertyVariablesRead(nil, d, m)
}

func resourcePropertyVariablesDelete(_ context.Context, d *schema.ResourceData, m interface{}) diag.Diagnostics {
	meta := akamai.Meta(m)
	logger := meta.Log("PAPI", "resourcePropertyVariablesDelete")
	logger.Debugf("DELETING")
	d.SetId("")
	logger.Debugf("Done")
	return nil
}

func resourcePropertyVariablesImport(_ context.Context, d *schema.ResourceData, m interface{}) ([]*schema.ResourceData, error) {
	resourceID := d.Id()
	propertyID := resourceID

	if !strings.HasPrefix(resourceID, "prp_") {
		keys := []papi.SearchKey{
			papi.SearchByPropertyName,
			papi.SearchByHostname,
			papi.SearchByEdgeHostname,
		}
		for _, searchKey := range keys {
			results, err := papi.Search(searchKey, resourceID, "") //<--correlationid
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
	err := property.GetProperty("")
	if err != nil {
		return nil, err
	}
	if err := d.Set("account", property.AccountID); err != nil {
		return nil, fmt.Errorf("%w: %s", tools.ErrValueSet, err.Error())
	}
	if err := d.Set("contract", property.ContractID); err != nil {
		return nil, fmt.Errorf("%w: %s", tools.ErrValueSet, err.Error())
	}
	if err := d.Set("group", property.GroupID); err != nil {
		return nil, fmt.Errorf("%w: %s", tools.ErrValueSet, err.Error())
	}
	if err := d.Set("name", property.PropertyName); err != nil {
		return nil, fmt.Errorf("%w: %s", tools.ErrValueSet, err.Error())
	}
	if err := d.Set("version", property.LatestVersion); err != nil {
		return nil, fmt.Errorf("%w: %s", tools.ErrValueSet, err.Error())
	}
	d.SetId(property.PropertyID)
	return []*schema.ResourceData{d}, nil
}

func resourcePropertyVariablesExists(d *schema.ResourceData, m interface{}) (bool, error) {
	meta := akamai.Meta(m)
	logger := meta.Log("PAPI", "resourcePropertyVariablesExists")
	variables := d.Id()
	// FIXME this function always returns true, nil
	if variables != "" {
		logger.Debugf("Check for variables  %s", variables)
		return true, nil
	}
	return true, nil
}

func resourcePropertyVariablesRead(_ context.Context, _ *schema.ResourceData, _ interface{}) diag.Diagnostics {
	return nil
}

func resourcePropertyVariablesUpdate(_ context.Context, d *schema.ResourceData, m interface{}) diag.Diagnostics {
	meta := akamai.Meta(m)
	logger := meta.Log(inst.Name())
	logger.Debugf("UPDATING")
	rule := papi.NewRule()
	logger.Debugf("START Check for variables")
	variables, err := tools.GetSetValue("variables", d)
	if err != nil {
		if err != tools.ErrNotFound {
			return diag.FromErr(err)
		}
		logger.Debugf("Done")
		return nil
	}
	logger.Debugf("Check for variables  %s", variables)
	for _, r := range variables.List() {
		variable, ok := r.(map[string]interface{})
		if !ok {
			continue
		}
		vv, ok := variable["variable"]
		if !ok {
			continue
		}
		for _, v := range vv.(*schema.Set).List() {
			variableMap, ok := v.(map[string]interface{})
			if !ok {
				continue
			}
			newVariable := papi.NewVariable()
			name, ok := variableMap["name"].(string)
			if !ok {
				return diag.FromErr(fmt.Errorf("%w: %s, %q", tools.ErrInvalidType, "name", "string"))
			}
			description, ok := variableMap["description"].(string)
			if !ok {
				return diag.FromErr(fmt.Errorf("%w: %s, %q", tools.ErrInvalidType, "description", "string"))
			}
			value, ok := variableMap["value"].(string)
			if !ok {
				return diag.FromErr(fmt.Errorf("%w: %s, %q", tools.ErrInvalidType, "value", "string"))
			}
			hidden, ok := variableMap["hidden"].(bool)
			if !ok {
				return diag.FromErr(fmt.Errorf("%w: %s, %q", tools.ErrInvalidType, "hidden", "bool"))
			}
			sensitive, ok := variableMap["sensitive"].(bool)
			if !ok {
				return diag.FromErr(fmt.Errorf("%w: %s, %q", tools.ErrInvalidType, "sensitive", "bool"))
			}
			newVariable.Name = name
			newVariable.Description = description
			newVariable.Value = value
			newVariable.Hidden = hidden
			newVariable.Sensitive = sensitive
			rule.AddVariable(newVariable)
		}
	}

	body, err := jsonhooks.Marshal(rule)
	if err != nil {
		return diag.FromErr(err)
	}
	logger.Debugf("JSON result  %s", string(body))
	sha := tools.GetSHAString(string(body))
	if err := d.Set("json", string(body)); err != nil {
		return diag.FromErr(fmt.Errorf("%w: %s", tools.ErrValueSet, err.Error()))
	}
	d.SetId(sha)
	logger.Debugf("Done")
	return nil
}
