package datastream

import (
	"context"
	"fmt"
	"strconv"
	"strings"

	"github.com/akamai/AkamaiOPEN-edgegrid-golang/v13/pkg/datastream"
	"github.com/akamai/AkamaiOPEN-edgegrid-golang/v13/pkg/session"
	"github.com/akamai/terraform-provider-akamai/v11/pkg/common/hash"
	"github.com/akamai/terraform-provider-akamai/v11/pkg/meta"
	"github.com/hashicorp/terraform-plugin-sdk/v2/diag"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/schema"
)

func dataSourceAppSecConfigs() *schema.Resource {
	return &schema.Resource{
		ReadContext: dataSourceAppSecConfigsRead,
		Schema: map[string]*schema.Schema{
			"group_id": {
				Type:        schema.TypeString,
				Required:    true,
				Description: "Identifies the group for which to list available AppSec configurations",
			},
			"contract_id": {
				Type:        schema.TypeString,
				Required:    true,
				Description: "Identifies the contract for which to list available AppSec configurations",
			},
			"app_sec_configs": {
				Type:        schema.TypeList,
				Computed:    true,
				Description: "List of AppSec configurations available for streaming in the given group and contract",
				Elem: &schema.Resource{
					Schema: map[string]*schema.Schema{
						"id": {
							Type:        schema.TypeInt,
							Computed:    true,
							Description: "The unique identifier of the AppSec configuration",
						},
						"name": {
							Type:        schema.TypeString,
							Computed:    true,
							Description: "The name of the AppSec configuration",
						},
						"file_type": {
							Type:        schema.TypeString,
							Computed:    true,
							Description: "The file type of the AppSec configuration (e.g. RBAC, WAF)",
						},
						"latest_version": {
							Type:        schema.TypeInt,
							Computed:    true,
							Description: "The latest version of the AppSec configuration",
						},
						"production_version": {
							Type:        schema.TypeInt,
							Computed:    true,
							Description: "The version currently active in production",
						},
						"target_product": {
							Type:        schema.TypeString,
							Computed:    true,
							Description: "The target product for the AppSec configuration",
						},
					},
				},
			},
		},
	}
}

func dataSourceAppSecConfigsRead(ctx context.Context, rd *schema.ResourceData, m interface{}) diag.Diagnostics {
	meta := meta.Must(m)
	logger := meta.Log("datastream", "dataSourceAppSecConfigsRead")

	ctx = session.ContextWithOptions(
		ctx,
		session.WithContextLog(logger),
	)

	logger.Debug("Listing AppSec configs")
	client := inst.Client(meta)

	groupIDStr := rd.Get("group_id").(string)
	groupIDStr = strings.TrimPrefix(groupIDStr, "grp_")
	groupID, err := strconv.Atoi(groupIDStr)
	if err != nil {
		return diag.Errorf("invalid group_id %q: %s", groupIDStr, err)
	}

	contractID := rd.Get("contract_id").(string)
	contractID = strings.TrimPrefix(contractID, "ctr_")

	configs, err := client.GetAppSecConfigs(ctx, datastream.GetAppSecConfigsRequest{
		GroupID:    groupID,
		ContractID: contractID,
	})
	if err != nil {
		return diag.FromErr(err)
	}

	appsecConfigs := parseAppSecConfigs(configs)

	if err := rd.Set("app_sec_configs", appsecConfigs); err != nil {
		return diag.Errorf("error setting app_sec_configs: %s", err)
	}

	md5Sum, _ := hash.GetMD5Sum(fmt.Sprintf("%v", appsecConfigs))
	rd.SetId(md5Sum)

	return nil
}

func parseAppSecConfigs(configs []datastream.AppSecConfigDetails) []map[string]interface{} {
	result := make([]map[string]interface{}, 0, len(configs))
	for _, c := range configs {
		result = append(result, map[string]interface{}{
			"id":                 c.ID,
			"name":               c.Name,
			"file_type":          c.FileType,
			"latest_version":     c.LatestVersion,
			"production_version": c.ProductionVersion,
			"target_product":     c.TargetProduct,
		})
	}
	return result
}
