package edgeworkers

import (
	"context"
	"fmt"
	"strings"

	"github.com/akamai/AkamaiOPEN-edgegrid-golang/v11/pkg/edgeworkers"

	"github.com/akamai/terraform-provider-akamai/v8/pkg/common/tf"
	"github.com/akamai/terraform-provider-akamai/v8/pkg/meta"
	"github.com/hashicorp/terraform-plugin-sdk/v2/diag"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/schema"
)

func dataSourceEdgeworkersResourceTier() *schema.Resource {
	return &schema.Resource{
		ReadContext: dataEdgeworkersResourceTierRead,
		Schema: map[string]*schema.Schema{
			"contract_id": {
				Type:             schema.TypeString,
				Required:         true,
				Description:      "Unique identifier of a contract",
				DiffSuppressFunc: tf.FieldPrefixSuppress("ctr_"),
			},
			"resource_tier_name": {
				Type:        schema.TypeString,
				Required:    true,
				Description: "Unique name of the resource tier",
			},
			"resource_tier_id": {
				Type:        schema.TypeInt,
				Computed:    true,
				Description: "Unique identifier of the resource tier",
			},
		},
	}
}

func dataEdgeworkersResourceTierRead(ctx context.Context, d *schema.ResourceData, m interface{}) diag.Diagnostics {
	meta := meta.Must(m)
	client := inst.Client(meta)
	log := meta.Log("Edgeworkers", "dataEdgeworkersResourceTierRead")
	log.Debug("Reading Resource Tier")

	contractID, err := tf.GetStringValue("contract_id", d)
	if err != nil {
		return diag.FromErr(err)
	}
	resourceTierName, err := tf.GetStringValue("resource_tier_name", d)
	if err != nil {
		return diag.FromErr(err)
	}

	resp, err := client.ListResourceTiers(ctx, edgeworkers.ListResourceTiersRequest{ContractID: strings.TrimPrefix(contractID, "ctr_")})
	if err != nil {
		return diag.FromErr(err)
	}

	var rt *edgeworkers.ResourceTier
	for _, r := range resp.ResourceTiers {
		if r.Name == resourceTierName {
			rt = &r
			break
		}
	}
	if rt == nil {
		return diag.Errorf("Resource tier with name: '%s' was not found", resourceTierName)
	}

	err = d.Set("resource_tier_id", rt.ID)
	if err != nil {
		return diag.FromErr(err)
	}

	d.SetId(fmt.Sprintf("%s:%s", contractID, resourceTierName))

	return nil
}
