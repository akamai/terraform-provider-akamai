package appsec

import (
	"context"
	"encoding/json"
	"fmt"

	v2 "github.com/akamai/AkamaiOPEN-edgegrid-golang/v2/pkg/appsec"
	"github.com/akamai/terraform-provider-akamai/v2/pkg/akamai"
	"github.com/akamai/terraform-provider-akamai/v2/pkg/tools"

	"github.com/hashicorp/terraform-plugin-sdk/v2/diag"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/schema"
)

func dataSourceContractsGroups() *schema.Resource {
	return &schema.Resource{
		ReadContext: dataSourceContractsGroupsRead,
		Schema: map[string]*schema.Schema{

			"json": {
				Type:     schema.TypeString,
				Computed: true,
			},
			"output_text": {
				Type:        schema.TypeString,
				Computed:    true,
				Description: "Text Export representation",
			},
		},
	}
}

func dataSourceContractsGroupsRead(ctx context.Context, d *schema.ResourceData, m interface{}) diag.Diagnostics {
	meta := akamai.Meta(m)
	client := inst.Client(meta)
	logger := meta.Log("APPSEC", "resourceContractsGroupsRead")

	getContractsGroups := v2.GetContractsGroupsRequest{}

	contractsgroups, err := client.GetContractsGroups(ctx, getContractsGroups)
	if err != nil {
		logger.Errorf("calling 'getContractsGroups': %s", err.Error())
		return diag.FromErr(err)
	}

	ots := OutputTemplates{}
	InitTemplates(ots)

	outputtext, err := RenderTemplates(ots, "contractsgroupsDS", contractsgroups)
	if err == nil {
		d.Set("output_text", outputtext)
	}

	jsonBody, err := json.Marshal(contractsgroups)
	if err != nil {
		return diag.FromErr(err)
	}

	if err := d.Set("json", string(jsonBody)); err != nil {
		return diag.FromErr(fmt.Errorf("%w: %s", tools.ErrValueSet, err.Error()))
	}

	if len(contractsgroups.ContractGroups) > 0 {
		d.SetId(contractsgroups.ContractGroups[0].ContractID)
	}
	return nil
}
