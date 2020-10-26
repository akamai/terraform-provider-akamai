package property

import (
	"bytes"
	"context"
	"encoding/json"
	"fmt"
	"github.com/akamai/AkamaiOPEN-edgegrid-golang/v2/pkg/papi"
	"github.com/akamai/terraform-provider-akamai/v2/pkg/akamai"
	"github.com/akamai/terraform-provider-akamai/v2/pkg/tools"
	"github.com/hashicorp/go-cty/cty"
	"github.com/hashicorp/terraform-plugin-sdk/v2/diag"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/schema"
)

func resourcePropertyRules() *schema.Resource {
	return &schema.Resource{
		CreateContext: resourcePropertyRulesCreate,
		ReadContext:   resourcePropertyRulesRead,
		UpdateContext: resourcePropertyRulesUpdate,
		DeleteContext: schema.NoopContext,
		Schema:        akamaiPropertyRulesSchema,
	}
}

var akamaiPropertyRulesSchema = map[string]*schema.Schema{
	"contract_id": {
		Type:      schema.TypeString,
		Required:  true,
		StateFunc: addPrefixToState("ctr_"),
	},
	"group_id": {
		Type:      schema.TypeString,
		Required:  true,
		StateFunc: addPrefixToState("grp_"),
	},
	"property_id": {
		Type:      schema.TypeString,
		Required:  true,
		StateFunc: addPrefixToState("prp_"),
	},
	"rules": {
		Type:             schema.TypeString,
		Required:         true,
		ValidateDiagFunc: validateJSON,
		DiffSuppressFunc: suppressRulesJSON,
		StateFunc: func(i interface{}) string {
			rules := i.(string)
			var res bytes.Buffer
			if err := json.Compact(&res, []byte(rules)); err != nil {
				panic(err)
			}
			return res.String()
		},
		Description: "JSON Rule representation",
	},
	"rule_format": {
		Type:        schema.TypeString,
		Optional:    true,
		Description: "TODO: this field is currently not used due to an issue in akamaized-edgegrid-client v2 library",
	},
	"version": {
		Type:        schema.TypeInt,
		Computed:    true,
		Description: "This is a computed value - provider will always use 'latest' version, providing own version number is not supported",
	},
	"warnings": {
		Type:     schema.TypeString,
		Computed: true,
	},
	"errors": {
		Type:     schema.TypeString,
		Computed: true,
	},
}

func resourcePropertyRulesCreate(ctx context.Context, d *schema.ResourceData, m interface{}) diag.Diagnostics {
	meta := akamai.Meta(m)
	client := inst.Client(meta)
	var contractID, groupID, propertyID, rulesJSON string
	var err error
	if contractID, err = tools.GetStringValue("contract_id", d); err != nil {
		return diag.FromErr(err)
	}
	contractID = tools.AddPrefix(contractID, "ctr_")
	if groupID, err = tools.GetStringValue("group_id", d); err != nil {
		return diag.FromErr(err)
	}
	groupID = tools.AddPrefix(groupID, "grp_")
	if propertyID, err = tools.GetStringValue("property_id", d); err != nil {
		return diag.FromErr(err)
	}
	propertyID = tools.AddPrefix(propertyID, "prp_")
	latestVersion, err := client.GetLatestVersion(ctx, papi.GetLatestVersionRequest{
		PropertyID: propertyID,
		ContractID: contractID,
		GroupID:    groupID,
	})
	if err != nil {
		return diag.FromErr(err)
	}
	version := latestVersion.Version.PropertyVersion
	if err := d.Set("version", version); err != nil {
		return diag.FromErr(err)
	}
	if rulesJSON, err = tools.GetStringValue("rules", d); err != nil {
		return diag.FromErr(err)
	}
	var rules papi.Rules
	if err := json.Unmarshal([]byte(rulesJSON), &rules); err != nil {
		return diag.FromErr(err)
	}

	res, err := client.UpdateRuleTree(ctx, papi.UpdateRulesRequest{
		PropertyID:      propertyID,
		PropertyVersion: version,
		ContractID:      contractID,
		GroupID:         groupID,
		Rules:           papi.RulesUpdate{Rules: rules},
		ValidateRules:   true,
	})
	if err != nil {
		return diag.FromErr(err)
	}
	d.SetId(propertyID)
	if res.Errors != nil {
		ruleErrors, err := json.Marshal(res.Errors)
		if err != nil {
			return diag.FromErr(err)
		}
		if err := d.Set("errors", string(ruleErrors)); err != nil {
			return diag.FromErr(err)
		}
	}
	return resourcePropertyRulesRead(ctx, d, m)
}

func resourcePropertyRulesRead(ctx context.Context, d *schema.ResourceData, m interface{}) diag.Diagnostics {
	meta := akamai.Meta(m)
	client := inst.Client(meta)
	var contractID, groupID, propertyID string
	var err error
	if contractID, err = tools.GetStringValue("contract_id", d); err != nil {
		return diag.FromErr(err)
	}
	contractID = tools.AddPrefix(contractID, "ctr_")
	if groupID, err = tools.GetStringValue("group_id", d); err != nil {
		return diag.FromErr(err)
	}
	groupID = tools.AddPrefix(groupID, "grp_")
	if propertyID, err = tools.GetStringValue("property_id", d); err != nil {
		return diag.FromErr(err)
	}
	propertyID = tools.AddPrefix(propertyID, "prp_")
	latestVersion, err := client.GetLatestVersion(ctx, papi.GetLatestVersionRequest{
		PropertyID: propertyID,
		ContractID: contractID,
		GroupID:    groupID,
	})
	if err != nil {
		return diag.FromErr(err)
	}
	version := latestVersion.Version.PropertyVersion
	if err := d.Set("version", version); err != nil {
		return diag.FromErr(err)
	}
	res, err := client.GetRuleTree(ctx, papi.GetRuleTreeRequest{
		PropertyID:      propertyID,
		PropertyVersion: version,
		ContractID:      contractID,
		GroupID:         groupID,
		ValidateRules:   true,
	})
	if err != nil {
		return diag.FromErr(err)
	}
	rulesJSON, err := json.Marshal(res.Rules)
	if err != nil {
		return diag.FromErr(err)
	}
	if err := d.Set("rules", string(rulesJSON)); err != nil {
		return diag.FromErr(err)
	}
	if res.Errors != nil {
		ruleErrors, err := json.Marshal(res.Errors)
		if err != nil {
			return diag.FromErr(err)
		}
		if err := d.Set("errors", string(ruleErrors)); err != nil {
			return diag.FromErr(err)
		}
	}
	return nil
}

func resourcePropertyRulesUpdate(ctx context.Context, d *schema.ResourceData, m interface{}) diag.Diagnostics {
	meta := akamai.Meta(m)
	client := inst.Client(meta)
	if d.HasChange("contract_id") {
		return diag.Errorf("contract_id field is immutable and cannot be updated")
	}
	if d.HasChange("group_id") {
		return diag.Errorf("group_id field is immutable and cannot be updated")
	}
	if d.HasChange("property_id") {
		return diag.Errorf("property_id field is immutable and cannot be updated")
	}
	var contractID, groupID, propertyID, rulesJSON string
	var err error
	if contractID, err = tools.GetStringValue("contract_id", d); err != nil {
		return diag.FromErr(err)
	}
	contractID = tools.AddPrefix(contractID, "ctr_")
	if groupID, err = tools.GetStringValue("group_id", d); err != nil {
		return diag.FromErr(err)
	}
	groupID = tools.AddPrefix(groupID, "grp_")
	if propertyID, err = tools.GetStringValue("property_id", d); err != nil {
		return diag.FromErr(err)
	}
	propertyID = tools.AddPrefix(propertyID, "prp_")
	latestVersion, err := client.GetLatestVersion(ctx, papi.GetLatestVersionRequest{
		PropertyID: propertyID,
		ContractID: contractID,
		GroupID:    groupID,
	})
	if err != nil {
		return diag.FromErr(err)
	}
	version := latestVersion.Version.PropertyVersion
	if latestVersion.Version.ProductionStatus != papi.VersionStatusInactive ||
		latestVersion.Version.StagingStatus != papi.VersionStatusInactive {
		// The latest version has been activated on either production or staging, so we need to create a new version to apply changes on
		newVersion, err := client.CreatePropertyVersion(ctx, papi.CreatePropertyVersionRequest{
			PropertyID: propertyID,
			ContractID: contractID,
			GroupID:    groupID,
			Version: papi.PropertyVersionCreate{
				CreateFromVersion: latestVersion.Version.PropertyVersion,
			},
		})
		if err != nil {
			return diag.FromErr(fmt.Errorf("%w: %s", ErrVersionCreate, err.Error()))
		}
		version = newVersion.PropertyVersion
	}
	if err := d.Set("version", version); err != nil {
		return diag.FromErr(err)
	}
	if rulesJSON, err = tools.GetStringValue("rules", d); err != nil {
		return diag.FromErr(err)
	}
	var rules papi.Rules
	if err := json.Unmarshal([]byte(rulesJSON), &rules); err != nil {
		return diag.FromErr(err)
	}

	res, err := client.UpdateRuleTree(ctx, papi.UpdateRulesRequest{
		PropertyID:      propertyID,
		PropertyVersion: version,
		ContractID:      contractID,
		GroupID:         groupID,
		Rules:           papi.RulesUpdate{Rules: rules},
		ValidateRules:   true,
	})
	if err != nil {
		return diag.FromErr(err)
	}
	d.SetId(propertyID)
	if res.Errors != nil {
		ruleErrors, err := json.Marshal(res.Errors)
		if err != nil {
			return diag.FromErr(err)
		}
		if err := d.Set("errors", string(ruleErrors)); err != nil {
			return diag.FromErr(err)
		}
	}
	return resourcePropertyRulesRead(ctx, d, m)
}

func validateJSON(val interface{}, _ cty.Path) diag.Diagnostics {
	if str, ok := val.(string); ok {
		var target map[string]interface{}
		if err := json.Unmarshal([]byte(str), &target); err != nil {
			return diag.FromErr(fmt.Errorf("invalid JSON: %w", err))
		}
		return nil
	}
	return diag.FromErr(fmt.Errorf("value is not a string: %s", val))
}

func addPrefixToState(pre string) schema.SchemaStateFunc {
	return func(given interface{}) string {
		str, ok := given.(string)
		if !ok {
			panic("interface should be string")
		}
		return tools.AddPrefix(str, pre)
	}
}

func suppressRulesJSON(_, old, new string, _ *schema.ResourceData) bool {
	logger := akamai.Log("PAPI", "suppressRulesJSON")
	var oldRules, newRules papi.Rules
	if old == "" || new == "" {
		return old == new
	}
	if err := json.Unmarshal([]byte(old), &oldRules); err != nil {
		logger.Errorf("Unable to unmarshal 'old' JSON rules: %s", err)
		return false
	}
	if err := json.Unmarshal([]byte(new), &newRules); err != nil {
		logger.Errorf("Unable to unmarshal 'new' JSON rules: %s", err)
		return false
	}
	return compareRules(&oldRules, &newRules)
}
