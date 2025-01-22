package property

import (
	"context"

	"github.com/hashicorp/terraform-plugin-sdk/v2/diag"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/schema"

	"github.com/akamai/AkamaiOPEN-edgegrid-golang/v9/pkg/papi"
	"github.com/akamai/AkamaiOPEN-edgegrid-golang/v9/pkg/session"
	"github.com/akamai/terraform-provider-akamai/v6/pkg/common/str"
	"github.com/akamai/terraform-provider-akamai/v6/pkg/common/tf"
	"github.com/akamai/terraform-provider-akamai/v6/pkg/meta"
)

func dataSourceProperties() *schema.Resource {
	return &schema.Resource{
		ReadContext: dataPropertiesRead,
		Schema: map[string]*schema.Schema{
			"group_id": {
				Type:             schema.TypeString,
				Required:         true,
				ValidateDiagFunc: tf.IsNotBlank,
			},
			"contract_id": {
				Type:             schema.TypeString,
				Required:         true,
				ValidateDiagFunc: tf.IsNotBlank,
			},
			"properties": {
				Type:        schema.TypeList,
				Computed:    true,
				Description: "List of properties",
				Elem: &schema.Resource{
					Schema: map[string]*schema.Schema{
						"contract_id": {
							Type:     schema.TypeString,
							Computed: true,
						},
						"group_id": {
							Type:     schema.TypeString,
							Computed: true,
						},
						"latest_version": {
							Type:     schema.TypeInt,
							Computed: true,
						},
						"note": {
							Type:     schema.TypeString,
							Computed: true,
						},
						"production_version": {
							Type:     schema.TypeInt,
							Computed: true,
						},
						"property_id": {
							Type:     schema.TypeString,
							Computed: true,
						},
						"property_name": {
							Type:     schema.TypeString,
							Computed: true,
						},
						"staging_version": {
							Type:     schema.TypeInt,
							Computed: true,
						},
					},
				},
			},
		},
	}
}

func dataPropertiesRead(ctx context.Context, d *schema.ResourceData, m interface{}) diag.Diagnostics {
	meta := meta.Must(m)
	log := meta.Log("PAPI", "dataPropertiesRead")
	// create a context with logging for api calls
	ctx = session.ContextWithOptions(
		ctx,
		session.WithContextLog(log),
	)
	log.Debug("Listing Properties")

	// groupID / contractID is string as per schema.
	groupID, err := tf.GetStringValue("group_id", d)
	if err != nil {
		return diag.FromErr(err)
	}
	groupID = str.AddPrefix(groupID, "grp_")
	contractID, err := tf.GetStringValue("contract_id", d)
	if err != nil {
		return diag.FromErr(err)
	}
	contractID = str.AddPrefix(contractID, "ctr_")

	propertiesResponse, err := getProperties(ctx, groupID, contractID, meta)
	if err != nil {
		return diag.Errorf("error listing properties: %v", err)
	}
	contractID = str.AddPrefix(contractID, "ctr_")

	// setting concatenated id to uniquely identify data
	d.SetId(groupID + contractID)

	properties, err := sliceResponseProperties(propertiesResponse)
	if err != nil {
		return diag.FromErr(err)
	}

	if err := d.Set("properties", properties); err != nil {
		return diag.Errorf("error setting properties: %s", err)
	}

	return nil
}

func sliceResponseProperties(propertiesResponse *papi.GetPropertiesResponse) ([]map[string]interface{}, error) {
	var properties []map[string]interface{}
	for _, item := range propertiesResponse.Properties.Items {

		property := map[string]interface{}{
			"contract_id":        item.ContractID,
			"group_id":           item.GroupID,
			"latest_version":     item.LatestVersion,
			"note":               item.Note,
			"production_version": decodeVersion(item.ProductionVersion),
			"property_id":        item.PropertyID,
			"property_name":      item.PropertyName,
			"staging_version":    decodeVersion(item.StagingVersion),
		}
		properties = append(properties, property)
	}
	return properties, nil
}

func decodeVersion(version interface{}) int {
	v, ok := version.(*int)
	if !ok || v == nil {
		return 0
	}
	return *v
}

// Reusable function to fetch all the properties for a given group and contract
func getProperties(ctx context.Context, groupID string, contractID string, meta meta.Meta) (*papi.GetPropertiesResponse, error) {
	client := Client(meta)
	req := papi.GetPropertiesRequest{
		ContractID: contractID,
		GroupID:    groupID,
	}
	props, err := client.GetProperties(ctx, req)
	if err != nil {
		return nil, err
	}
	return props, nil
}
