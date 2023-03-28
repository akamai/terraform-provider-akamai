package property

import (
	"context"
	"encoding/json"
	"fmt"

	"github.com/akamai/AkamaiOPEN-edgegrid-golang/v5/pkg/papi"
	"github.com/akamai/AkamaiOPEN-edgegrid-golang/v5/pkg/session"
	"github.com/akamai/terraform-provider-akamai/v3/pkg/akamai"
	"github.com/akamai/terraform-provider-akamai/v3/pkg/tools"
	"github.com/hashicorp/terraform-plugin-sdk/v2/diag"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/schema"
)

func dataSourcePropertyIncludeRules() *schema.Resource {
	return &schema.Resource{
		ReadContext: dataPropertyIncludeRulesRead,
		Schema: map[string]*schema.Schema{
			"contract_id": {
				Type:        schema.TypeString,
				Required:    true,
				Description: "Identifies the contract under which the include was created",
			},
			"group_id": {
				Type:        schema.TypeString,
				Required:    true,
				Description: "Identifies the group under which the include was created",
			},
			"include_id": {
				Type:        schema.TypeString,
				Required:    true,
				Description: "The unique identifier of the property include",
			},
			"version": {
				Type:        schema.TypeInt,
				Required:    true,
				Description: "The version of the include that a rule tree represents",
			},
			"name": {
				Type:        schema.TypeString,
				Computed:    true,
				Description: "A descriptive name for the include",
			},
			"rules": {
				Type:        schema.TypeString,
				Computed:    true,
				Description: "Property Rules as JSON",
			},
			"rule_errors": {
				Type:        schema.TypeString,
				Computed:    true,
				Description: "Rules validation errors",
			},
			"rule_warnings": {
				Type:        schema.TypeString,
				Computed:    true,
				Description: "Rules validation warnings",
			},
			"rule_format": {
				Type:        schema.TypeString,
				Computed:    true,
				Description: "Indicates the versioned set of features and criteria that are currently applied to a rule tree",
			},
			"type": {
				Type:        schema.TypeString,
				Computed:    true,
				Description: "Specifies the type of the include, either `MICROSERVICES` or `COMMON_SETTINGS`",
			},
		},
	}
}

func dataPropertyIncludeRulesRead(ctx context.Context, d *schema.ResourceData, m interface{}) diag.Diagnostics {
	meta := akamai.Meta(m)
	client := inst.Client(meta)
	log := meta.Log("PAPI", "dataPropertyIncludeRulesRead")
	ctx = session.ContextWithOptions(ctx, session.WithContextLog(log))

	groupID, err := tools.GetStringValue("group_id", d)
	if err != nil {
		return diag.FromErr(err)
	}

	contractID, err := tools.GetStringValue("contract_id", d)
	if err != nil {
		return diag.FromErr(err)
	}

	includeID, err := tools.GetStringValue("include_id", d)
	if err != nil {
		return diag.FromErr(err)
	}
	includeVersion, err := tools.GetIntValue("version", d)
	if err != nil {
		return diag.FromErr(err)
	}

	includeRuleTree, err := client.GetIncludeRuleTree(ctx, papi.GetIncludeRuleTreeRequest{
		ContractID:     contractID,
		GroupID:        groupID,
		IncludeID:      includeID,
		IncludeVersion: includeVersion,
		ValidateRules:  true,
	})
	if err != nil {
		return diag.FromErr(err)
	}

	rules, err := json.MarshalIndent(papi.RulesUpdate{
		Rules:    includeRuleTree.Rules,
		Comments: includeRuleTree.Comments,
	}, "", " ")
	if err != nil {
		return diag.FromErr(err)
	}

	attrs := map[string]interface{}{
		"name":        includeRuleTree.IncludeName,
		"rules":       string(rules),
		"rule_format": includeRuleTree.RuleFormat,
		"type":        includeRuleTree.IncludeType,
	}

	var ruleErrors string
	if len(includeRuleTree.Errors) > 0 {
		rulesErrorsJSON, err := json.MarshalIndent(includeRuleTree.Errors, "", "  ")
		if err != nil {
			return diag.FromErr(err)
		}
		ruleErrors = string(rulesErrorsJSON)
	}
	attrs["rule_errors"] = ruleErrors

	var ruleWarnings string
	if len(includeRuleTree.Warnings) > 0 {
		rulesWarningsJSON, err := json.MarshalIndent(includeRuleTree.Warnings, "", "  ")
		if err != nil {
			return diag.FromErr(err)
		}
		ruleWarnings = string(rulesWarningsJSON)
	}
	attrs["rule_warnings"] = ruleWarnings

	if err := tools.SetAttrs(d, attrs); err != nil {
		return diag.FromErr(err)
	}

	d.SetId(fmt.Sprintf("%s:%d", includeID, includeVersion))
	return nil
}
