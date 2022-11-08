package property

import (
	"context"

	"github.com/akamai/AkamaiOPEN-edgegrid-golang/v2/pkg/papi"
	"github.com/hashicorp/terraform-plugin-sdk/v2/diag"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/schema"

	"github.com/akamai/terraform-provider-akamai/v2/pkg/akamai"
	"github.com/akamai/terraform-provider-akamai/v2/pkg/tools"
)

func dataSourceAkamaiPropertyInclude() *schema.Resource {
	return &schema.Resource{
		ReadContext: dataAkamaiPropertyIncludeRead,
		Schema: map[string]*schema.Schema{
			"contract_id": {
				Type:        schema.TypeString,
				Required:    true,
				Description: "Identifies the contract under which the include was created",
			},
			"group_id": {
				Type:        schema.TypeString,
				Required:    true,
				Description: "Identifies the group under which the include was created",
			},
			"include_id": {
				Type:        schema.TypeString,
				Required:    true,
				Description: "The unique identifier of the property include",
			},
			"name": {
				Type:        schema.TypeString,
				Computed:    true,
				Description: "A descriptive name for the include",
			},
			"type": {
				Type:        schema.TypeString,
				Computed:    true,
				Description: "Specifies the type of the include, either 'MICROSERVICES' or 'COMMON_SETTINGS'",
			},
			"latest_version": {
				Type:        schema.TypeInt,
				Computed:    true,
				Description: "Specifies the most recent version of the include",
			},
			"staging_version": {
				Type:        schema.TypeInt,
				Computed:    true,
				Description: "The most recent version which was activated to the test network",
			},
			"production_version": {
				Type:        schema.TypeInt,
				Computed:    true,
				Description: "The most recent version which was activated to the production network",
			},
		},
	}
}

func dataAkamaiPropertyIncludeRead(ctx context.Context, d *schema.ResourceData, m interface{}) diag.Diagnostics {
	meta := akamai.Meta(m)
	client := inst.Client(meta)
	log := meta.Log("PAPI", "dataAkamaiPropertyIncludeRead")
	log.Debug("Reading Property Include")

	contractID, err := tools.GetStringValue("contract_id", d)
	if err != nil {
		return diag.FromErr(err)
	}

	groupID, err := tools.GetStringValue("group_id", d)
	if err != nil {
		return diag.FromErr(err)
	}

	includeID, err := tools.GetStringValue("include_id", d)
	if err != nil {
		return diag.FromErr(err)
	}

	include, err := client.GetInclude(ctx, papi.GetIncludeRequest{
		ContractID: contractID,
		GroupID:    groupID,
		IncludeID:  includeID,
	})
	if err != nil {
		return diag.FromErr(err)
	}

	if len(include.Includes.Items) == 0 {
		// this one probably shouldn't ever happen,
		return diag.Errorf("empty include response from api")
	}
	item := include.Includes.Items[0]

	attrs := map[string]interface{}{
		"name":           item.IncludeName,
		"type":           item.IncludeType,
		"latest_version": item.LatestVersion,
	}
	if item.StagingVersion != nil {
		attrs["staging_version"] = item.StagingVersion
	}
	if item.ProductionVersion != nil {
		attrs["production_version"] = item.ProductionVersion
	}
	if err = tools.SetAttrs(d, attrs); err != nil {
		return diag.FromErr(err)
	}

	d.SetId(includeID)
	return nil
}
