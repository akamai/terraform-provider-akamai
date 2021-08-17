package property

import (
	"context"
	"encoding/json"
	"fmt"

	"github.com/hashicorp/terraform-plugin-sdk/v2/diag"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/schema"

	"github.com/akamai/AkamaiOPEN-edgegrid-golang/v2/pkg/papi"
	"github.com/akamai/terraform-provider-akamai/v2/pkg/akamai"
	"github.com/akamai/terraform-provider-akamai/v2/pkg/tools"
)

func dataPropertyRules() *schema.Resource {
	return &schema.Resource{
		ReadContext: dataPropertyRulesRead,
		StateUpgraders: []schema.StateUpgrader{{
			Version: 0,
			Type:    dataAkamaiPropertyRuleSchemaV0().CoreConfigSchema().ImpliedType(),
			Upgrade: upgradeAkamaiPropertyRuleStateV0,
		}},
		SchemaVersion: 1,
		Schema:        dataAkamaiPropertyRuleSchema,
	}
}

var dataAkamaiPropertyRuleSchema = map[string]*schema.Schema{
	"contract_id": {
		Type:             schema.TypeString,
		Optional:         true,
		Computed:         true,
		StateFunc:        addPrefixToState("ctr_"),
		RequiredWith:     []string{"group_id"},
		ValidateDiagFunc: tools.IsNotBlank,
	},
	"group_id": {
		Type:             schema.TypeString,
		Optional:         true,
		Computed:         true,
		StateFunc:        addPrefixToState("grp_"),
		RequiredWith:     []string{"contract_id"},
		ValidateDiagFunc: tools.IsNotBlank,
	},
	"property_id": {
		Type:             schema.TypeString,
		Required:         true,
		StateFunc:        addPrefixToState("prp_"),
		ValidateDiagFunc: tools.IsNotBlank,
	},
	"version": {
		Type:        schema.TypeInt,
		Optional:    true,
		Computed:    true,
		Description: "This is a computed value - provider will always use 'latest' version, providing own version number is not supported",
	},
	"rules": {
		Type:        schema.TypeString,
		Computed:    true,
		Description: "JSON Rule representation",
	},
	"errors": {
		Type:     schema.TypeString,
		Computed: true,
	},
}

func dataPropertyRulesRead(ctx context.Context, d *schema.ResourceData, m interface{}) diag.Diagnostics {
	meta := akamai.Meta(m)
	client := inst.Client(meta)
	logger := meta.Log("PAPI", "dataPropertyRulesRead")

	var (
		contractID, groupID, propertyID string
		version                         int
		err                             error
	)

	// since contractID && groupID is optional, we should not return an error.
	contractID, _ = tools.GetStringValue("contract_id", d)
	groupID, _ = tools.GetStringValue("group_id", d)

	if propertyID, err = tools.GetStringValue("property_id", d); err != nil {
		return diag.FromErr(err)
	}

	if contractID != "" {
		contractID = tools.AddPrefix(contractID, "ctr_")
		if err := d.Set("contract_id", contractID); err != nil {
			return diag.Errorf("%v: %s", tools.ErrValueSet, err.Error())
		}
	}
	if groupID != "" {
		groupID = tools.AddPrefix(groupID, "grp_")
		if err := d.Set("group_id", groupID); err != nil {
			return diag.Errorf("%v: %s", tools.ErrValueSet, err.Error())
		}
	}

	if version, err = tools.GetIntValue("version", d); err != nil {
		latestVersion, err := client.GetLatestVersion(ctx, papi.GetLatestVersionRequest{
			PropertyID: propertyID,
			ContractID: contractID,
			GroupID:    groupID,
		})
		if err != nil {
			return diag.FromErr(err)
		}

		version = latestVersion.Version.PropertyVersion
		contractID = latestVersion.ContractID
		groupID = latestVersion.GroupID

		if err := d.Set("version", version); err != nil {
			return diag.FromErr(err)
		}
	}

	getRuleTreeResponse, err := client.GetRuleTree(ctx, papi.GetRuleTreeRequest{
		PropertyID:      propertyID,
		PropertyVersion: version,
		ContractID:      contractID,
		GroupID:         groupID,
		ValidateRules:   true,
		ValidateMode:    papi.RuleValidateModeFull,
	})
	if err != nil {
		return diag.FromErr(err)
	}

	formattedRulesJSON, err := json.MarshalIndent(papi.RulesUpdate{Rules: getRuleTreeResponse.Rules}, "", "  ")
	if err != nil {
		logger.Debugf("Creating rule tree resulted in invalid JSON: %s", err)
		return diag.FromErr(fmt.Errorf("invalid JSON result: %w", err))
	}
	if err := d.Set("rules", string(formattedRulesJSON)); err != nil {
		return diag.FromErr(err)
	}

	if len(getRuleTreeResponse.Errors) != 0 {
		ruleErrors, err := json.Marshal(getRuleTreeResponse.Errors)
		if err != nil {
			return diag.FromErr(err)
		}
		if err := d.Set("errors", string(ruleErrors)); err != nil {
			return diag.FromErr(err)
		}
	}
	d.SetId(propertyID)

	return nil
}
