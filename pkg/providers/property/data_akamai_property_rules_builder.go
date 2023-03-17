package property

import (
	"context"
	"crypto/md5"
	"encoding/hex"
	"encoding/json"

	"github.com/akamai/AkamaiOPEN-edgegrid-golang/v4/pkg/papi"
	"github.com/akamai/terraform-provider-akamai/v3/pkg/akamai"
	"github.com/akamai/terraform-provider-akamai/v3/pkg/providers/property/ruleformats"
	"github.com/hashicorp/terraform-plugin-sdk/v2/diag"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/schema"
)

func dataSourcePropertyRulesBuilder() *schema.Resource {
	rulesSchemas := ruleformats.Schemas()

	rulesSchemas["json"] = &schema.Schema{
		Type:        schema.TypeString,
		Computed:    true,
		Description: "JSON representation of provided rules",
	}

	rulesSchemas["rule_format"] = &schema.Schema{
		Type:        schema.TypeString,
		Computed:    true,
		Description: "Frozen rule format in which the rules are represented",
	}

	return &schema.Resource{
		ReadContext: dataSourcePropertyRulesBuilderRead,
		Schema:      rulesSchemas,
	}
}

func dataSourcePropertyRulesBuilderRead(_ context.Context, d *schema.ResourceData, m interface{}) diag.Diagnostics {
	meta := akamai.Meta(m)
	logger := meta.Log("PAPI", "dataSourcePropertyRulesBuilderRead")
	logger.Debug("dataSourcePropertyRulesBuilderRead")

	rules, err := ruleformats.NewBuilder(d).Build()
	if err != nil {
		return diag.Errorf("error building rules: %s", err)
	}

	rulesUpdate := papi.RulesUpdate{
		Rules: *rules,
	}

	JSON, err := json.MarshalIndent(rulesUpdate, "", "  ")
	if err != nil {
		return diag.Errorf("error marshaling rules to json: %s", err)
	}

	if err := d.Set("json", string(JSON)); err != nil {
		return diag.Errorf("error setting json in schema: %s", err)
	}

	ruleFormat := ruleformats.GetUsedRuleFormat(d)
	if err := d.Set("rule_format", ruleFormat.Version()); err != nil {
		return diag.Errorf("error setting rule_format in schema %s", err)
	}

	sum := md5.Sum([]byte(JSON))
	hexsum := hex.EncodeToString(sum[:])
	d.SetId(hexsum)
	return nil
}
