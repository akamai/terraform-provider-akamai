package iam

import (
	"context"

	"github.com/akamai/AkamaiOPEN-edgegrid-golang/v7/pkg/iam"
	"github.com/akamai/AkamaiOPEN-edgegrid-golang/v7/pkg/session"
	"github.com/akamai/terraform-provider-akamai/v5/pkg/meta"
	"github.com/apex/log"
	"github.com/hashicorp/terraform-plugin-sdk/v2/diag"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/schema"
)

func dataSourceIAMTimezones() *schema.Resource {
	return &schema.Resource{
		Description: "List all the possible time zones Akamai supports",
		ReadContext: dataIAMTimezonesRead,
		Schema: map[string]*schema.Schema{
			"timezones": {
				Type:     schema.TypeSet,
				Computed: true,
				Elem: &schema.Resource{
					Schema: map[string]*schema.Schema{
						"timezone": {
							Type:        schema.TypeString,
							Description: "The time zone ID",
							Computed:    true,
						},
						"description": {
							Type:        schema.TypeString,
							Description: "The description of a time zone, including the GMT +/-",
							Computed:    true,
						},
						"offset": {
							Type:        schema.TypeString,
							Description: "The time zone offset from GMT",
							Computed:    true,
						},
						"posix": {
							Type:        schema.TypeString,
							Description: "The time zone posix",
							Computed:    true,
						},
					},
				},
			},
		},
	}
}

func dataIAMTimezonesRead(ctx context.Context, d *schema.ResourceData, m interface{}) diag.Diagnostics {
	meta := meta.Must(m)
	logger := meta.Log("IAM", "dataIAMTimezonesRead")
	ctx = session.ContextWithOptions(ctx, session.WithContextLog(logger))
	client := inst.Client(meta)

	logger.Debug("Fetching time zones")

	res, err := client.SupportedTimezones(ctx)
	if err != nil {
		logger.WithError(err).Error("could not get time zones")
		return diag.FromErr(err)
	}

	if err := d.Set("timezones", timezonesToState(res)); err != nil {
		log.WithError(err).Error("could not set time zones in state")
		return diag.FromErr(err)
	}

	d.SetId("akamai_iam_timezones")
	return nil
}

func timezonesToState(timezones []iam.Timezone) []interface{} {
	out := make([]interface{}, 0, len(timezones))

	for _, t := range timezones {
		timezone := map[string]interface{}{
			"timezone":    t.Timezone,
			"description": t.Description,
			"offset":      t.Offset,
			"posix":       t.Posix,
		}

		out = append(out, timezone)
	}

	return out
}
