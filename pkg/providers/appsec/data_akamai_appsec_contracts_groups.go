package appsec

import (
	"context"
	"encoding/json"
	"errors"

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
			"contractid": {
				Type:     schema.TypeString,
				Optional: true,
			},
			"groupid": {
				Type:     schema.TypeInt,
				Optional: true,
			},
			"json": {
				Type:     schema.TypeString,
				Computed: true,
			},
			"output_text": {
				Type:        schema.TypeString,
				Computed:    true,
				Description: "Text Export representation",
			},
			"default_contractid": {
				Type:     schema.TypeString,
				Computed: true,
			},
			"default_groupid": {
				Type:     schema.TypeInt,
				Computed: true,
			},
		},
	}
}

func dataSourceContractsGroupsRead(ctx context.Context, d *schema.ResourceData, m interface{}) diag.Diagnostics {
	meta := akamai.Meta(m)
	client := inst.Client(meta)
	logger := meta.Log("APPSEC", "dataSourceContractsGroupsRead")

	getContractsGroups := v2.GetContractsGroupsRequest{}

	contract, err := tools.GetStringValue("contractid", d)
	if err != nil && !errors.Is(err, tools.ErrNotFound) {
		return diag.FromErr(err)
	}
	getContractsGroups.ContractID = contract

	group, err := tools.GetIntValue("groupid", d)
	if err != nil && !errors.Is(err, tools.ErrNotFound) {
		return diag.FromErr(err)
	}
	getContractsGroups.GroupID = group

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
		return diag.Errorf("%s: %s", tools.ErrValueSet, err.Error())
	}

	for _, configval := range contractsgroups.ContractGroups {

		if configval.ContractID == contract && configval.GroupID == group {
			if err := d.Set("default_contractid", contract); err != nil {
				return diag.Errorf("%s: %s", tools.ErrValueSet, err.Error())
			}
			if err := d.Set("default_groupid", group); err != nil {
				return diag.Errorf("%s: %s", tools.ErrValueSet, err.Error())
			}
		}
	}
	if len(contractsgroups.ContractGroups) > 0 {
		d.SetId(contractsgroups.ContractGroups[0].ContractID)
	}

	return nil
}
