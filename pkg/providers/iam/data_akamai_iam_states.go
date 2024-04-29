package iam

import (
	"context"

	"github.com/akamai/AkamaiOPEN-edgegrid-golang/v9/pkg/iam"
	"github.com/akamai/AkamaiOPEN-edgegrid-golang/v9/pkg/session"
	"github.com/akamai/terraform-provider-akamai/v6/pkg/meta"
	"github.com/hashicorp/terraform-plugin-sdk/v2/diag"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/schema"
)

func dataSourceIAMStates() *schema.Resource {
	return &schema.Resource{
		Description: "List US states or Canadian provinces.",
		ReadContext: dataIAMStatesRead,
		Schema: map[string]*schema.Schema{
			"country": {
				Type:        schema.TypeString,
				Required:    true,
				Description: "Specifies a US state or Canadian province.",
			},
			"states": {
				Type:        schema.TypeSet,
				Computed:    true,
				Description: "Supported states.",
				Elem:        &schema.Schema{Type: schema.TypeString},
			},
		},
	}
}

func dataIAMStatesRead(ctx context.Context, d *schema.ResourceData, m interface{}) diag.Diagnostics {
	meta := meta.Must(m)
	logger := meta.Log("IAM", "dataIAMStatesRead")
	ctx = session.ContextWithOptions(ctx, session.WithContextLog(logger))
	client := inst.Client(meta)

	logger.Debug("Fetching states")

	country := d.Get("country").(string)
	res, err := client.ListStates(ctx, iam.ListStatesRequest{Country: country})
	if err != nil {
		logger.Error("Could not get states", "error", err)
		return diag.FromErr(err)
	}

	states := []interface{}{}
	for _, state := range res {
		states = append(states, state)
	}

	if err := d.Set("states", states); err != nil {
		logger.Error("Could not set states in the state", "error", err)
		return diag.FromErr(err)
	}

	d.SetId("akamai_iam_states")
	return nil
}
