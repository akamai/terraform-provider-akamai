package property

import (
	"context"
	"strconv"

	"github.com/akamai/AkamaiOPEN-edgegrid-golang/v12/pkg/papi"
	"github.com/akamai/terraform-provider-akamai/v8/pkg/common/tf"
	"github.com/akamai/terraform-provider-akamai/v8/pkg/meta"
	"github.com/hashicorp/terraform-plugin-sdk/v2/diag"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/schema"
)

func dataSourcePropertyIncludeParents() *schema.Resource {
	return &schema.Resource{
		ReadContext: dataPropertyIncludeParentsRead,
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
				Description: "The list of include's parents",
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

func dataPropertyIncludeParentsRead(ctx context.Context, d *schema.ResourceData, m interface{}) diag.Diagnostics {
	meta := meta.Must(m)
	client := Client(meta)
	log := meta.Log("PAPI", "dataPropertyIncludeParentsRead")
	log.Debug("Reading Property Include Parents")

	contractID, err := tf.GetStringValue("contract_id", d)
	if err != nil {
		return diag.FromErr(err)
	}

	groupID, err := tf.GetStringValue("group_id", d)
	if err != nil {
		return diag.FromErr(err)
	}

	includeID, err := tf.GetStringValue("include_id", d)
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
		var isIncUsedInStagingVer, isIncUsedInProductionVer bool
		listRefIncReq := papi.ListReferencedIncludesRequest{
			PropertyID: item.PropertyID,
			ContractID: contractID,
			GroupID:    groupID,
		}

		if item.StagingVersion != nil {
			stagingVersion = strconv.Itoa(*item.StagingVersion)
			isIncUsedInStagingVer = true
		}
		if item.ProductionVersion != nil {
			productionVersion = strconv.Itoa(*item.ProductionVersion)
			isIncUsedInProductionVer = true
		}
		if stagingVersion != productionVersion && item.StagingVersion != nil && item.ProductionVersion != nil {
			listRefIncReq.PropertyVersion = *item.StagingVersion
			isIncUsedInStagingVer, err = isIncPresentInReferencedIncludes(ctx, client, listRefIncReq, includeID)
			if err != nil {
				return diag.FromErr(err)
			}
			listRefIncReq.PropertyVersion = *item.ProductionVersion
			isIncUsedInProductionVer, err = isIncPresentInReferencedIncludes(ctx, client, listRefIncReq, includeID)
			if err != nil {
				return diag.FromErr(err)
			}
		}

		attrs := map[string]interface{}{
			"id":                                    item.PropertyID,
			"name":                                  item.PropertyName,
			"staging_version":                       stagingVersion,
			"production_version":                    productionVersion,
			"is_include_used_in_staging_version":    isIncUsedInStagingVer,
			"is_include_used_in_production_version": isIncUsedInProductionVer,
		}
		includeParents = append(includeParents, attrs)
	}

	if err := d.Set("parents", includeParents); err != nil {
		return diag.FromErr(err)
	}

	d.SetId(includeID)
	return nil
}

func isIncPresentInReferencedIncludes(ctx context.Context, client papi.PAPI, refIncArgs papi.ListReferencedIncludesRequest, includeID string) (bool, error) {
	refIncResp, err := client.ListReferencedIncludes(ctx, refIncArgs)
	if err != nil {
		return false, err
	}
	for _, include := range refIncResp.Includes.Items {
		if include.IncludeID == includeID {
			return true, nil
		}
	}
	return false, nil
}
