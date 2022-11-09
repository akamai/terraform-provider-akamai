package property

import (
	"context"
	"strconv"

	"github.com/akamai/AkamaiOPEN-edgegrid-golang/v2/pkg/papi"
	"github.com/akamai/terraform-provider-akamai/v2/pkg/akamai"
	"github.com/akamai/terraform-provider-akamai/v2/pkg/tools"
	"github.com/hashicorp/terraform-plugin-sdk/v2/diag"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/schema"
)

func dataSourceAkamaiPropertyIncludeParents() *schema.Resource {
	return &schema.Resource{
		ReadContext: dataSourceAkamaiPropertyIncludeParentsRead,
		Schema: map[string]*schema.Schema{
			"contract_id": {
				Type:        schema.TypeString,
				Required:    true,
				Description: "Identifies the contract under which the data was requested",
			},
			"group_id": {
				Type:        schema.TypeString,
				Required:    true,
				Description: "Identifies the group under which the data was requested",
			},
			"include_id": {
				Type:        schema.TypeString,
				Required:    true,
				Description: "Unique identifier of the include",
			},
			"parents": {
				Type:        schema.TypeList,
				Computed:    true,
				Description: "The list of include’s parents",
				Elem: &schema.Resource{
					Schema: map[string]*schema.Schema{
						"id": {
							Type:        schema.TypeString,
							Computed:    true,
							Description: "The property's unique identifier",
						},
						"name": {
							Type:        schema.TypeString,
							Computed:    true,
							Description: "A descriptive name for the property",
						},
						"staging_version": {
							Type:        schema.TypeString,
							Computed:    true,
							Description: "The most recent property version to be activated to the staging network",
						},
						"production_version": {
							Type:        schema.TypeString,
							Computed:    true,
							Description: "The most recent property version to be activated to the production network",
						},
						"is_include_used_in_staging_version": {
							Type:        schema.TypeBool,
							Computed:    true,
							Description: "Indicates if the include is used in the staging network",
						},
						"is_include_used_in_production_version": {
							Type:        schema.TypeBool,
							Computed:    true,
							Description: "Indicates if the include is used in the production network",
						},
					},
				},
			},
		},
	}
}

func dataSourceAkamaiPropertyIncludeParentsRead(ctx context.Context, d *schema.ResourceData, m interface{}) diag.Diagnostics {
	meta := akamai.Meta(m)
	client := inst.Client(meta)
	log := meta.Log("PAPI", "dataSourceAkamaiPropertyIncludeParentsRead")
	log.Debug("Reading Property Include Parents")

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

	resp, err := client.ListIncludeParents(ctx, papi.ListIncludeParentsRequest{
		ContractID: contractID,
		GroupID:    groupID,
		IncludeID:  includeID,
	})
	if err != nil {
		return diag.FromErr(err)
	}

	var includeParents []map[string]interface{}
	for _, item := range resp.Properties.Items {
		var stagingVersion, productionVersion string
		if item.StagingVersion != nil {
			stagingVersion = strconv.Itoa(*item.StagingVersion)
		}
		if item.ProductionVersion != nil {
			productionVersion = strconv.Itoa(*item.ProductionVersion)
		}
		attrs := map[string]interface{}{
			"id":                                    item.PropertyID,
			"name":                                  item.PropertyName,
			"staging_version":                       stagingVersion,
			"production_version":                    productionVersion,
			"is_include_used_in_staging_version":    item.IsIncludeUsedInStagingVersion,
			"is_include_used_in_production_version": item.IsIncludeUsedInProductionVersion,
		}
		includeParents = append(includeParents, attrs)
	}

	if err := d.Set("parents", includeParents); err != nil {
		return diag.FromErr(err)
	}

	d.SetId(includeID)
	return nil
}
