package iam

import (
	"context"

	"github.com/akamai/AkamaiOPEN-edgegrid-golang/v8/pkg/session"
	"github.com/akamai/terraform-provider-akamai/v6/pkg/meta"
	"github.com/hashicorp/terraform-plugin-sdk/v2/diag"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/schema"
)

func dataSourceIAMLanguages() *schema.Resource {
	return &schema.Resource{
		Description: "List all the possible languages Akamai supports",
		ReadContext: dataIAMLanguagesRead,
		Schema: map[string]*schema.Schema{
			"languages": {
				Type:        schema.TypeSet,
				Computed:    true,
				Description: "Languages supported by Akamai",
				Elem:        &schema.Schema{Type: schema.TypeString},
			},
		},
	}
}

func dataIAMLanguagesRead(ctx context.Context, d *schema.ResourceData, m interface{}) diag.Diagnostics {
	meta := meta.Must(m)
	logger := meta.Log("IAM", "dataIAMLanguagesRead")
	ctx = session.ContextWithOptions(ctx, session.WithContextLog(logger))
	client := inst.Client(meta)

	logger.Debug("Fetching supported supported languages")
	res, err := client.SupportedLanguages(ctx)
	if err != nil {
		logger.WithError(err).Error("Could not get supported supported languages")
		return diag.FromErr(err)
	}

	languages := []interface{}{}
	for _, language := range res {
		languages = append(languages, language)
	}

	if err := d.Set("languages", languages); err != nil {
		logger.WithError(err).Error("Could not set supported languages in state")
		return diag.FromErr(err)
	}

	d.SetId("akamai_iam_supported_langs")
	return nil
}
