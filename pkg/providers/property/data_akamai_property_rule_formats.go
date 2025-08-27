package property

import (
	"context"
	"fmt"

	"github.com/akamai/terraform-provider-akamai/v9/pkg/meta"
	"github.com/hashicorp/terraform-plugin-sdk/v2/diag"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/schema"
)

func dataSourcePropertyRuleFormats() *schema.Resource {
	return &schema.Resource{
		ReadContext: dataPropertyRuleFormatsRead,
		Schema: map[string]*schema.Schema{
			"rule_format": {
				Type:     schema.TypeList,
				Elem:     &schema.Schema{Type: schema.TypeString},
				Computed: true,
			},
		},
	}
}

func dataPropertyRuleFormatsRead(ctx context.Context, d *schema.ResourceData, m interface{}) diag.Diagnostics {
	meta := meta.Must(m)
	client := Client(meta)

	logger := meta.Log("PAPI", "dataPropertyRuleFormatsRead")
	logger.Debugf("read property rule formats")

	// Get property rule formats
	ruleFormats, err := client.GetRuleFormats(ctx)
	if err != nil {
		return diag.FromErr(err)
	}

	// check if ruleFormats exist.
	if len(ruleFormats.RuleFormats.Items) == 0 {
		return diag.FromErr(fmt.Errorf("%v", ErrRuleFormatsNotFound))
	}

	// set rule_format value.
	if err := d.Set("rule_format", ruleFormats.RuleFormats.Items); err != nil {
		return diag.FromErr(err)
	}
	d.SetId("rule_format")

	return nil
}
