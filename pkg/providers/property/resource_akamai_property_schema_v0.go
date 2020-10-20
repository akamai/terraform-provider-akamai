package property

import (
	"context"
	"fmt"
	"strings"

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

			// Deprecated in SchemaVersion 1 (copied to new attributes in state migration)
			"contract": {Type: schema.TypeString, Optional: true, ForceNew: true},
			"group":    {Type: schema.TypeString, Optional: true, ForceNew: true},
			"product":  {Type: schema.TypeString, Optional: true, ForceNew: true},

			// Hard deprecated in SchemaVersion 1 (no migration necessary but operations always fail on non-zero value)
			"rule_format": {Type: schema.TypeString, Optional: true},
			"cp_code":     {Type: schema.TypeString, Optional: true},
			"is_secure":   {Type: schema.TypeBool, Optional: true},
			"rules":       {Type: schema.TypeString, Optional: true},
			"variables":   {Type: schema.TypeString, Optional: true},
			"contact":     {Type: schema.TypeSet, Required: true, Elem: &schema.Schema{Type: schema.TypeString}},
			"hostnames":   {Type: schema.TypeMap, Required: true, Elem: &schema.Schema{Type: schema.TypeString}},
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
func upgradePropV0(ctx context.Context, rawState map[string]interface{}, _ interface{}) (map[string]interface{}, error) {
	// Delete computed and removed attributes
	removed := []string{
		"account",
		"version",
		"rulessha",
		"edge_hostnames",
	}
	for _, attr := range removed {
		if _, ok := rawState[attr]; ok {
			delete(rawState, attr)
		}
	}

	// Deprecated attribute contract checked for prefixed ID and copied to contract_id
	if v, ok := rawState["contract"]; ok {
		s := v.(string) // Schema guarantees this is a string
		if !strings.HasPrefix(s, "ctr_") {
			s = fmt.Sprintf("%s%s", "ctr_", s)
		}

		rawState["contract"] = s
		rawState["contract_id"] = s
	}

	// Deprecated attribute group checked for prefixed ID and copied to group_id
	if v, ok := rawState["group"]; ok {
		s := v.(string) // Schema guarantees this is a string
		if !strings.HasPrefix(s, "grp_") {
			s = fmt.Sprintf("%s%s", "grp_", s)
		}

		rawState["group"] = s
		rawState["group_id"] = s
	}

	// Deprecated attribute product checked for prefixed ID and copied to product_id
	if v, ok := rawState["product"]; ok {
		s := v.(string) // Schema guarantees this is a string
		if !strings.HasPrefix(s, "prd_") {
			s = fmt.Sprintf("%s%s", "prd_", s)
		}

		rawState["product"] = s
		rawState["product_id"] = s
	}

	return rawState, nil
}
