package property

import (
	"context"
	"encoding/json"

	"github.com/hashicorp/terraform-plugin-sdk/v2/diag"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/schema"

	"github.com/akamai/AkamaiOPEN-edgegrid-golang/v2/pkg/papi"
	"github.com/akamai/AkamaiOPEN-edgegrid-golang/v2/pkg/session"
	"github.com/akamai/terraform-provider-akamai/v2/pkg/akamai"
	"github.com/akamai/terraform-provider-akamai/v2/pkg/tools"
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

	// groupID / contractID is string as per schema.
	groupID := tools.AddPrefix(d.Get("group_id").(string), "grp_")
	contractID := tools.AddPrefix(d.Get("contract_id").(string), "ctr_")

	properties, err := getProperties(ctx, groupID, contractID, meta)
	if err != nil {
		return diag.Errorf("error listing properties: %v", err)
	}
	// setting concatenated id to uniquely identify data
	d.SetId(groupID + contractID)
	props, err := json.Marshal(properties)
	if err != nil {
		return diag.FromErr(err)
	}
	/* setting raw json here in current scope. We have to set all the json fields in
	property struct for more granular access to properties object */
	if err := d.Set("properties", string(props)); err != nil {
		return diag.Errorf("error setting properties: %s", err)
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
func getProperties(ctx context.Context, groupID string, contractID string, meta akamai.OperationMeta) (*papi.GetPropertiesResponse, error) {
	client := inst.Client(meta)
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
