package appsec

import (
	"context"
	"encoding/json"
	"errors"

	"github.com/akamai/AkamaiOPEN-edgegrid-golang/v10/pkg/appsec"
	"github.com/akamai/terraform-provider-akamai/v8/pkg/common/tf"
	"github.com/akamai/terraform-provider-akamai/v8/pkg/meta"
	"github.com/hashicorp/terraform-plugin-sdk/v2/diag"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/schema"
)

func dataSourceContractsGroups() *schema.Resource {
	return &schema.Resource{
		ReadContext: dataSourceContractsGroupsRead,
		Schema: map[string]*schema.Schema{
			"contractid": {
				Type:        schema.TypeString,
				Optional:    true,
				Description: "Unique identifier of an Akamai contract",
			},
			"groupid": {
				Type:        schema.TypeInt,
				Optional:    true,
				Description: "Unique identifier of a contract group",
			},
			"json": {
				Type:        schema.TypeString,
				Computed:    true,
				Description: "JSON representation",
			},
			"output_text": {
				Type:        schema.TypeString,
				Computed:    true,
				Description: "Text representation",
			},
			"default_contractid": {
				Type:        schema.TypeString,
				Computed:    true,
				Description: "Default contract ID for the specified contract and group",
			},
			"default_groupid": {
				Type:        schema.TypeInt,
				Computed:    true,
				Description: "Default group ID for the specified contract and group",
			},
		},
	}
}

func dataSourceContractsGroupsRead(ctx context.Context, d *schema.ResourceData, m interface{}) diag.Diagnostics {
	meta := meta.Must(m)
	client := inst.Client(meta)
	logger := meta.Log("APPSEC", "dataSourceContractsGroupsRead")

	getContractsGroups := appsec.GetContractsGroupsRequest{}

	contractID, err := tf.GetStringValue("contractid", d)
	if err != nil && !errors.Is(err, tf.ErrNotFound) {
		return diag.FromErr(err)
	}
	getContractsGroups.ContractID = contractID

	group, err := tf.GetIntValue("groupid", d)
	if err != nil && !errors.Is(err, tf.ErrNotFound) {
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
	if err != nil {
		return diag.FromErr(err)
	}
	if err := d.Set("output_text", outputtext); err != nil {
		return diag.Errorf("%s: %s", tf.ErrValueSet, err.Error())
	}

	jsonBody, err := json.Marshal(contractsgroups)
	if err != nil {
		return diag.FromErr(err)
	}

	if err := d.Set("json", string(jsonBody)); err != nil {
		return diag.Errorf("%s: %s", tf.ErrValueSet, err.Error())
	}

	for _, configval := range contractsgroups.ContractGroups {

		if configval.ContractID == contractID && configval.GroupID == group {
			if err := d.Set("default_contractid", contractID); err != nil {
				return diag.Errorf("%s: %s", tf.ErrValueSet, err.Error())
			}
			if err := d.Set("default_groupid", group); err != nil {
				return diag.Errorf("%s: %s", tf.ErrValueSet, err.Error())
			}
		}
	}
	if len(contractsgroups.ContractGroups) > 0 {
		d.SetId(contractsgroups.ContractGroups[0].ContractID)
	}

	return nil
}
