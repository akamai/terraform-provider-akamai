package property

import (
	"context"
	"encoding/json"
	"errors"
	"fmt"
	"strings"

	"github.com/akamai/terraform-provider-akamai/v2/pkg/akamai"
	"github.com/akamai/terraform-provider-akamai/v2/pkg/tools"
	"github.com/hashicorp/terraform-plugin-sdk/v2/diag"

	"github.com/akamai/AkamaiOPEN-edgegrid-golang/v2/pkg/papi"
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

func resourcePropertyVariablesCreate(ctx context.Context, d *schema.ResourceData, m interface{}) diag.Diagnostics {
	meta := akamai.Meta(m)
	logger := meta.Log("PAPI", "resourcePropertyVariablesCreate")
	rule := &papi.Rules{}
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
			newVariable := papi.RuleVariable{}
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
			rule.Variables = append(rule.Variables, newVariable)
		}
	}

	body, err := json.Marshal(rule)
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

	return resourcePropertyVariablesRead(ctx, d, m)
}

func resourcePropertyVariablesDelete(_ context.Context, d *schema.ResourceData, m interface{}) diag.Diagnostics {
	meta := akamai.Meta(m)
	logger := meta.Log("PAPI", "resourcePropertyVariablesDelete")
	logger.Debugf("DELETING")
	d.SetId("")
	logger.Debugf("Done")
	return nil
}

func resourcePropertyVariablesImport(ctx context.Context, d *schema.ResourceData, m interface{}) ([]*schema.ResourceData, error) {
	meta := akamai.Meta(m)
	client := inst.Client(meta)
	propertyID := d.Id()

	if !strings.HasPrefix(propertyID, "prp_") {
		keys := []string{
			papi.SearchKeyPropertyName,
			papi.SearchKeyHostname,
			papi.SearchKeyEdgeHostname,
		}
		for _, searchKey := range keys {
			results, err := client.SearchProperties(ctx, papi.SearchRequest{
				Key:   searchKey,
				Value: propertyID,
			})
			if err != nil {
				continue
			}

			if results != nil && len(results.Versions.Items) > 0 {
				propertyID = results.Versions.Items[0].PropertyID
				break
			}
		}
	}

	res, err := client.GetProperty(ctx, papi.GetPropertyRequest{
		PropertyID: propertyID,
	})
	if err != nil {
		return nil, err
	}
	d.SetId(res.Property.PropertyID)
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
	logger := meta.Log(inst.Name(), "resourcePropertyVariablesUpdate")
	logger.Debugf("UPDATING")
	rule := &papi.Rules{}
	logger.Debugf("START Check for variables")
	variables, err := tools.GetSetValue("variables", d)
	if err != nil {
		if !errors.Is(err, tools.ErrNotFound) {
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
			newVariable := papi.RuleVariable{}
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
			rule.Variables = append(rule.Variables, newVariable)
		}
	}

	body, err := json.Marshal(rule)
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
