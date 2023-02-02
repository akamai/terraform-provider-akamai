package networklists

import (
	"context"
	"errors"

	"github.com/akamai/AkamaiOPEN-edgegrid-golang/v4/pkg/networklists"
	"github.com/akamai/terraform-provider-akamai/v3/pkg/akamai"
	"github.com/akamai/terraform-provider-akamai/v3/pkg/tools"
	"github.com/hashicorp/terraform-plugin-sdk/v2/diag"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/schema"
)

// network_lists v2
//
// https://techdocs.akamai.com/network-lists/reference/api
func resourceNetworkListDescription() *schema.Resource {
	return &schema.Resource{
		CreateContext: resourceNetworkListDescriptionUpdate,
		ReadContext:   resourceNetworkListDescriptionRead,
		UpdateContext: resourceNetworkListDescriptionUpdate,
		DeleteContext: resourceNetworkListDescriptionDelete,
		Importer: &schema.ResourceImporter{
			StateContext: schema.ImportStatePassthroughContext,
		},
		Schema: map[string]*schema.Schema{
			"network_list_id": {
				Type:     schema.TypeString,
				Required: true,
			},
			"name": {
				Type:     schema.TypeString,
				Required: true,
			},
			"description": {
				Type:     schema.TypeString,
				Required: true,
			},
		},
	}
}

func resourceNetworkListDescriptionRead(ctx context.Context, d *schema.ResourceData, m interface{}) diag.Diagnostics {
	meta := akamai.Meta(m)
	client := inst.Client(meta)
	logger := meta.Log("NETWORKLIST", "resourceNetworkListDescriptionRead")

	getNetworkListDescriptionRequest := networklists.GetNetworkListDescriptionRequest{}

	getNetworkListDescriptionRequest.UniqueID = d.Id()

	networklistdescription, err := client.GetNetworkListDescription(ctx, getNetworkListDescriptionRequest)
	if err != nil {
		logger.Errorf("calling 'getNetworkListDescription': %s", err.Error())
		return diag.FromErr(err)
	}

	d.SetId(networklistdescription.UniqueID)

	return nil
}

func resourceNetworkListDescriptionDelete(ctx context.Context, d *schema.ResourceData, m interface{}) diag.Diagnostics {

	return schema.NoopContext(ctx, d, m)
}

func resourceNetworkListDescriptionUpdate(ctx context.Context, d *schema.ResourceData, m interface{}) diag.Diagnostics {
	meta := akamai.Meta(m)
	client := inst.Client(meta)
	logger := meta.Log("NETWORKLIST", "resourceNetworkListDescriptionUpdate")

	updateNetworkListDescriptionRequest := networklists.UpdateNetworkListDescriptionRequest{}

	uniqueID, err := tools.GetStringValue("network_list_id", d)
	if err != nil && !errors.Is(err, tools.ErrNotFound) {
		return diag.FromErr(err)
	}
	updateNetworkListDescriptionRequest.UniqueID = uniqueID

	name, err := tools.GetStringValue("name", d)
	if err != nil && !errors.Is(err, tools.ErrNotFound) {
		return diag.FromErr(err)
	}
	updateNetworkListDescriptionRequest.Name = name

	description, err := tools.GetStringValue("description", d)
	if err != nil && !errors.Is(err, tools.ErrNotFound) {
		return diag.FromErr(err)
	}
	updateNetworkListDescriptionRequest.Description = description

	_, err = client.UpdateNetworkListDescription(ctx, updateNetworkListDescriptionRequest)
	if err != nil {
		logger.Errorf("calling 'updateNetworkListDescription': %s", err.Error())
		return diag.FromErr(err)
	}

	d.SetId(updateNetworkListDescriptionRequest.UniqueID)

	return resourceNetworkListDescriptionRead(ctx, d, m)
}
