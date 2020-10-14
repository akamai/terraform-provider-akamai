package property

import (
	"context"

	"github.com/hashicorp/terraform-plugin-sdk/v2/diag"

	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/schema"
)

func dataPropertyRules() *schema.Resource {
	return &schema.Resource{
		ReadContext: dataPropRulesOperation,
		Schema:      propertyRulesSchemaAttrs(), // in resource_akamai_property_rules.go
	}
}

func dataPropRulesOperation(_ context.Context, d *schema.ResourceData, _ interface{}) diag.Diagnostics {
	// Always fail for new resources and changed values
	if d.Id() == "" || d.IsNewResource() || d.HasChanges("variables", "rules") {
		return diag.Errorf(`data "akamai_property_rules" is no longer supported - See Akamai Terraform Upgrade Guide`)
	}

	// No changes and resource already exists, must be from previous version of plugin
	return nil
}
