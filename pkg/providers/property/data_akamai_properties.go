package property

import (
	"context"
	"fmt"
	"github.com/akamai/AkamaiOPEN-edgegrid-golang/v2/pkg/papi"
	"github.com/akamai/AkamaiOPEN-edgegrid-golang/v2/pkg/session"
	"github.com/akamai/terraform-provider-akamai/v2/pkg/akamai"
	"github.com/akamai/terraform-provider-akamai/v2/pkg/tools"
	"github.com/hashicorp/terraform-plugin-sdk/v2/diag"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/schema"
)

func dataSourceAkamaiProperties() *schema.Resource {
	return &schema.Resource{
		ReadContext: dataAkamaiPropertiesRead,
		Schema: map[string]*schema.Schema{
			"group_id": {
				Type:     schema.TypeString,
				Required: true,
			},
			"contract_id": {
				Type:     schema.TypeString,
				Required: true,
			},
			"properties": {
				Type:        schema.TypeList,
				Computed:    true,
				Description: "List of properties",
				Elem: &schema.Resource{
					Schema: map[string]*schema.Schema{
						"account_id":         {Type: schema.TypeString, Computed: true},
						"asset_id":           {Type: schema.TypeString, Computed: true},
						"contract_id":        {Type: schema.TypeString, Computed: true},
						"group_id":           {Type: schema.TypeString, Computed: true},
						"latest_version":     {Type: schema.TypeInt, Computed: true},
						"note":               {Type: schema.TypeString, Computed: true},
						"product_id":         {Type: schema.TypeString, Computed: true},
						"production_version": {Type: schema.TypeInt, Computed: true},
						"property_id":        {Type: schema.TypeString, Computed: true},
						"property_name":      {Type: schema.TypeString, Computed: true},
						"rule_format":        {Type: schema.TypeString, Computed: true},
						"staging_version":    {Type: schema.TypeInt, Computed: true},
					},
				},
			},
		},
	}
}

func dataAkamaiPropertiesRead(ctx context.Context, d *schema.ResourceData, m interface{}) diag.Diagnostics {
	meta := akamai.Meta(m)
	log := meta.Log("PAPI", "dataAkamaiPropertiesRead")
	// create a context with logging for api calls
	ctx = session.ContextWithOptions(
		ctx,
		session.WithContextLog(log),
	)
	log.Debug("Listing Properties")

	groupId, err := tools.GetStringValue("group_id", d)
	if err != nil {
		return diag.FromErr(err)
	}
	contractId, err := tools.GetStringValue("contract_id", d)
	if err != nil {
		return diag.FromErr(err)
	}
	groupId = tools.AddPrefix(groupId, "grp_")
	contractId = tools.AddPrefix(contractId, "ctr_")
	propertiesResponse, err := getProperties(ctx, groupId, contractId, meta)
	if err != nil {
		return diag.FromErr(fmt.Errorf("error listing properties: %w", err))
	}
	// setting concatenated id to uniquely identify data
	d.SetId(groupId + contractId)

	properties := make([]map[string]interface{}, 0)
	for _, item := range propertiesResponse.Properties.Items {
		property := map[string]interface{}{
			"account_id":         item.AccountID,
			"asset_id":           item.AssetID,
			"contract_id":        item.ContractID,
			"group_id":           item.GroupID,
			"latest_version":     item.LatestVersion,
			"note":               item.Note,
			"product_id":         item.ProductID,
			"production_version": decodeVersion(item.ProductionVersion),
			"property_id":        item.PropertyID,
			"property_name":      item.PropertyName,
			"rule_format":        item.RuleFormat,
			"staging_version":    decodeVersion(item.StagingVersion),
		}
		properties = append(properties, property)
	}

	if err := d.Set("properties", properties); err != nil {
		return diag.FromErr(fmt.Errorf("error setting properties: %s", err))
	}
	return nil
}

func decodeVersion(version interface{}) int {
	v, ok := version.(*int)
	if !ok || v == nil {
		return 0
	}
	return *v
}

// Reusable function to fetch all the properties for a given group and contract
func getProperties(ctx context.Context, groupId string, contractId string, meta akamai.OperationMeta) (*papi.GetPropertiesResponse, error) {
	client := inst.Client(meta)
	req := papi.GetPropertiesRequest{
		ContractID: contractId,
		GroupID:    groupId,
	}
	props, err := client.GetProperties(ctx, req)
	if err != nil {
		return nil, err
	}
	return props, nil
}
