package property

import (
	"context"
	"fmt"

	"github.com/hashicorp/terraform-plugin-sdk/v2/diag"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/schema"

	"github.com/akamai/AkamaiOPEN-edgegrid-golang/v6/pkg/papi"
	"github.com/akamai/terraform-provider-akamai/v3/pkg/akamai"
	"github.com/akamai/terraform-provider-akamai/v3/pkg/common/tf"
)

func dataSourcePropertiesSearch() *schema.Resource {
	return &schema.Resource{
		ReadContext: dataPropertiesSearchRead,
		Schema: map[string]*schema.Schema{
			"key": {
				Type:             schema.TypeString,
				Required:         true,
				Description:      "Key must have one of three values: 'edgeHostname', 'hostname' or 'propertyName'",
				ValidateDiagFunc: tf.ValidateStringInSlice([]string{papi.SearchKeyEdgeHostname, papi.SearchKeyHostname, papi.SearchKeyPropertyName}),
			},
			"value": {
				Type:             schema.TypeString,
				Required:         true,
				Description:      "Value of Search",
				ValidateDiagFunc: tf.IsNotBlank,
			},
			"properties": {
				Type:        schema.TypeList,
				Computed:    true,
				Description: "List of properties",
				Elem: &schema.Resource{
					Schema: map[string]*schema.Schema{
						"account_id":        {Type: schema.TypeString, Computed: true},
						"asset_id":          {Type: schema.TypeString, Computed: true},
						"contract_id":       {Type: schema.TypeString, Computed: true},
						"group_id":          {Type: schema.TypeString, Computed: true},
						"property_id":       {Type: schema.TypeString, Computed: true},
						"property_version":  {Type: schema.TypeInt, Computed: true},
						"property_name":     {Type: schema.TypeString, Computed: true},
						"edge_hostname":     {Type: schema.TypeString, Computed: true},
						"hostname":          {Type: schema.TypeString, Computed: true},
						"production_status": {Type: schema.TypeString, Computed: true},
						"staging_status":    {Type: schema.TypeString, Computed: true},
						"updated_by_user":   {Type: schema.TypeString, Computed: true},
						"updated_date":      {Type: schema.TypeString, Computed: true},
					},
				},
			},
		},
	}
}

func dataPropertiesSearchRead(ctx context.Context, d *schema.ResourceData, m interface{}) diag.Diagnostics {
	meta := akamai.Meta(m)
	client := inst.Client(meta)

	log := meta.Log("PAPI", "dataPropertiesSearchRead")

	log.Debug("Searching properties")

	var key, value string
	var err error

	if key, err = tf.GetStringValue("key", d); err != nil {
		return diag.FromErr(err)
	}
	if value, err = tf.GetStringValue("value", d); err != nil {
		return diag.FromErr(err)
	}

	request := papi.SearchRequest{
		Key:   key,
		Value: value,
	}

	search, err := client.SearchProperties(ctx, request)
	if err != nil {
		return diag.Errorf("could not load properties: %s", err)
	}

	d.SetId(fmt.Sprintf("%s:%s", request.Key, request.Value))

	if err := d.Set("properties", sliceResponseSearch(search)); err != nil {
		return diag.Errorf("error setting properties: %s", err)
	}

	return nil
}

func sliceResponseSearch(propertiesResponse *papi.SearchResponse) []map[string]interface{} {
	var properties []map[string]interface{}
	for _, item := range propertiesResponse.Versions.Items {
		property := map[string]interface{}{
			"account_id":        item.AccountID,
			"asset_id":          item.AssetID,
			"contract_id":       item.ContractID,
			"group_id":          item.GroupID,
			"property_id":       item.PropertyID,
			"property_version":  item.PropertyVersion,
			"property_name":     item.PropertyName,
			"edge_hostname":     item.EdgeHostname,
			"hostname":          item.Hostname,
			"production_status": item.ProductionStatus,
			"staging_status":    item.StagingStatus,
			"updated_by_user":   item.UpdatedByUser,
			"updated_date":      item.UpdatedDate,
		}
		properties = append(properties, property)
	}
	return properties
}
