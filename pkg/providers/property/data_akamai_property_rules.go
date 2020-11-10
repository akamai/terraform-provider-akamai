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
		ReadContext: dataPropRulesOperation,
		Schema:      dataPropertyRulesSchema,
	}
}

var dataPropertyRulesSchema = map[string]*schema.Schema{
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
		Type:        schema.TypeString,
		Computed:    true,
		Description: "JSON Rule representation",
	},
	"version": {
		Type:        schema.TypeInt,
		Optional:    true,
		Description: "This is a computed value - provider will always use 'latest' version, providing own version number is not supported",
	},
	"errors": {
		Type:     schema.TypeString,
		Computed: true,
	},
}

func dataPropRulesOperation(ctx context.Context, d *schema.ResourceData, m interface{}) diag.Diagnostics {
	// Always fail for new resources and changed values
	if d.Id() == "" || d.IsNewResource() || d.HasChanges("variables", "rules") {
		return diag.Errorf(`data "akamai_property_rules" is no longer supported - See Akamai Terraform Upgrade Guide`)
	}

	meta := akamai.Meta(m)
	client := inst.Client(meta)

	var (
		contractID, groupID, propertyID string
		version                         int
		err                             error
	)

	if contractID, err = tools.GetStringValue("contract_id", d); err != nil {
		return diag.FromErr(err)
	}
	contractID = tools.AddPrefix(contractID, "ctr_")
	if err := d.Set("contract_id", contractID); err != nil {
		return diag.FromErr(fmt.Errorf("%w: %s", tools.ErrValueSet, err.Error()))
	}

	if groupID, err = tools.GetStringValue("group_id", d); err != nil {
		return diag.FromErr(err)
	}
	groupID = tools.AddPrefix(groupID, "grp_")
	if err := d.Set("group_id", groupID); err != nil {
		return diag.FromErr(fmt.Errorf("%w: %s", tools.ErrValueSet, err.Error()))
	}

	if propertyID, err = tools.GetStringValue("property_id", d); err != nil {
		return diag.FromErr(err)
	}
	propertyID = tools.AddPrefix(propertyID, "prp_")
	if err := d.Set("property_id", propertyID); err != nil {
		return diag.FromErr(fmt.Errorf("%w: %s", tools.ErrValueSet, err.Error()))
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
		if err := d.Set("version", version); err != nil {
			return diag.FromErr(err)
		}
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
