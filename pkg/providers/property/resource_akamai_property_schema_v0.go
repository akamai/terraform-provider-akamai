package property

import (
	"context"

	"github.com/akamai/terraform-provider-akamai/v2/pkg/tools"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/schema"
)

// SchemaVersion 0 of the property resource -- this is referenced in migrations to SchemaVersion 1
func resourcePropertyV0() *schema.Resource {
	return &schema.Resource{
		Schema: map[string]*schema.Schema{
			// Unchanged in SchemaVersion 1
			"name":               {Type: schema.TypeString, Required: true, ForceNew: true},
			"staging_version":    {Type: schema.TypeInt, Computed: true},
			"production_version": {Type: schema.TypeInt, Computed: true},
			"rule_format":        {Type: schema.TypeString, Optional: true},
			"rules":              {Type: schema.TypeString, Optional: true},
			"hostnames":          {Type: schema.TypeMap, Required: true, Elem: &schema.Schema{Type: schema.TypeString}},

			// Deprecated in SchemaVersion 1 (copied to new attributes in state migration)
			"contract": {Type: schema.TypeString, Optional: true, ForceNew: true},
			"group":    {Type: schema.TypeString, Optional: true, ForceNew: true},
			"product":  {Type: schema.TypeString, Optional: true, ForceNew: true},

			// Hard deprecated in SchemaVersion 1 (no migration necessary but operations always fail on non-zero value)
			"cp_code":   {Type: schema.TypeString, Optional: true},
			"is_secure": {Type: schema.TypeBool, Optional: true},
			"variables": {Type: schema.TypeString, Optional: true},
			"contact":   {Type: schema.TypeSet, Required: true, Elem: &schema.Schema{Type: schema.TypeString}},
			"origin": {
				Type:     schema.TypeSet,
				Optional: true,
				Elem: &schema.Resource{
					Schema: map[string]*schema.Schema{
						"hostname":              {Type: schema.TypeString, Required: true},
						"port":                  {Type: schema.TypeInt, Optional: true, Default: 80},
						"forward_hostname":      {Type: schema.TypeString, Optional: true, Default: "ORIGIN_HOSTNAME"},
						"cache_key_hostname":    {Type: schema.TypeString, Optional: true, Default: "ORIGIN_HOSTNAME"},
						"compress":              {Type: schema.TypeBool, Optional: true, Default: false},
						"enable_true_client_ip": {Type: schema.TypeBool, Optional: true, Default: false},
					},
				},
			},

			// Computed & removed in SchemaVersion 1 (removed by state migration)
			"account":        {Type: schema.TypeString, Computed: true},
			"version":        {Type: schema.TypeInt, Computed: true},
			"rulessha":       {Type: schema.TypeString, Computed: true},
			"edge_hostnames": {Type: schema.TypeMap, Computed: true, Elem: &schema.Schema{Type: schema.TypeString}},
		},
	}
}

// Upgrade state from schema version 0 to 1
func upgradePropV0(_ context.Context, rawState map[string]interface{}, _ interface{}) (map[string]interface{}, error) {
	// Delete computed and removed attributes
	delete(rawState, "account")
	delete(rawState, "edge_hostnames")
	delete(rawState, "rulessha")
	delete(rawState, "version")

	// Deprecated attribute contract checked for prefixed ID and copied to contract_id
	if v, ok := rawState["contract"]; ok {
		s := tools.AddPrefix(v.(string), "ctr_") // Schema guarantees this is a string

		rawState["contract_id"] = s
		rawState["contract"] = s
	}

	// Deprecated attribute group checked for prefixed ID and copied to group_id
	if v, ok := rawState["group"]; ok {
		s := tools.AddPrefix(v.(string), "grp_") // Schema guarantees this is a string

		rawState["group_id"] = s
		rawState["group"] = s
	}

	// Deprecated attribute product checked for prefixed ID and copied to product_id
	if v, ok := rawState["product"]; ok {
		s := tools.AddPrefix(v.(string), "prd_") // Schema guarantees this is a string

		rawState["product_id"] = s
		rawState["product"] = s
	}

	return rawState, nil
}
