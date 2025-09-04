package iam

import (
	"context"

	"github.com/akamai/AkamaiOPEN-edgegrid-golang/v12/pkg/session"
	"github.com/akamai/terraform-provider-akamai/v9/pkg/meta"
	"github.com/hashicorp/terraform-plugin-sdk/v2/diag"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/schema"
)

func dataSourceIAMContactTypes() *schema.Resource {
	return &schema.Resource{
		Description: "Retrieve all contact types that Akamai supports.",
		ReadContext: dataIAMContactTypesRead,
		Schema: map[string]*schema.Schema{
			"contact_types": {
				Type:        schema.TypeSet,
				Computed:    true,
				Description: "Supported contact types.",
				Elem:        &schema.Schema{Type: schema.TypeString},
			},
		},
	}
}

func dataIAMContactTypesRead(ctx context.Context, d *schema.ResourceData, m interface{}) diag.Diagnostics {
	meta := meta.Must(m)
	logger := meta.Log("IAM", "dataIAMContactTypesRead")
	ctx = session.ContextWithOptions(ctx, session.WithContextLog(logger))
	client := inst.Client(meta)

	logger.Debug("Fetching supported contact types")

	res, err := client.SupportedContactTypes(ctx)
	if err != nil {
		logger.Error("Could not get supported contact types", "error", err)
		return diag.FromErr(err)
	}

	var types []interface{}
	for _, ct := range res {
		types = append(types, ct)
	}

	if err := d.Set("contact_types", types); err != nil {
		logger.Error("Could not set contact types in state", "error", err)
		return diag.FromErr(err)
	}

	d.SetId("akamai_iam_contact_types")
	return nil
}
