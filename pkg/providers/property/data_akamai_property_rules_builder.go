package property

import (
	"context"
	"crypto/md5"
	"encoding/hex"
	"encoding/json"
	"errors"

	"github.com/akamai/AkamaiOPEN-edgegrid-golang/v6/pkg/papi"
	"github.com/akamai/terraform-provider-akamai/v4/pkg/akamai"
	"github.com/akamai/terraform-provider-akamai/v4/pkg/providers/property/ruleformats"
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
		diags := diag.Errorf("building rules: %s", err)
		if errors.Is(err, ruleformats.ErrTooManyElements) {
			diags[0].Detail = "You can have only one behavior/criterion in a single block. Make sure each of them is placed into a separate block."
		}
		return diags
	}

	rulesUpdate := papi.RulesUpdate{
		Rules: *rules,
	}

	JSON, err := json.MarshalIndent(rulesUpdate, "", "  ")
	if err != nil {
		return diag.Errorf("marshaling rules to json: %s", err)
	}

	if err := d.Set("json", string(JSON)); err != nil {
		return diag.Errorf("setting json in schema: %s", err)
	}

	ruleFormat := ruleformats.GetUsedRuleFormat(d)
	if err := d.Set("rule_format", ruleFormat.Version()); err != nil {
		return diag.Errorf("setting rule_format in schema %s", err)
	}

	sum := md5.Sum(JSON)
	hexsum := hex.EncodeToString(sum[:])
	d.SetId(hexsum)
	return nil
}
