package appsec

import (
	"context"
	"encoding/json"
	"errors"

	"github.com/akamai/AkamaiOPEN-edgegrid-golang/v6/pkg/appsec"
	"github.com/akamai/terraform-provider-akamai/v4/pkg/common/tf"
	"github.com/akamai/terraform-provider-akamai/v4/pkg/meta"
	"github.com/hashicorp/terraform-plugin-sdk/v2/diag"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/schema"
)

func dataSourceCustomDeny() *schema.Resource {
	return &schema.Resource{
		ReadContext: dataSourceCustomDenyRead,
		Schema: map[string]*schema.Schema{
			"config_id": {
				Type:        schema.TypeInt,
				Required:    true,
				Description: "Unique identifier of the security configuration",
			},
			"custom_deny_id": {
				Type:        schema.TypeString,
				Optional:    true,
				Description: "Unique identifier of the custom deny for which to return information",
			},
			"json": {
				Type:        schema.TypeString,
				Computed:    true,
				Description: "JSON representation",
			},
			"output_text": {
				Type:        schema.TypeString,
				Computed:    true,
				Description: "Text representation",
			},
		},
	}
}

func dataSourceCustomDenyRead(ctx context.Context, d *schema.ResourceData, m interface{}) diag.Diagnostics {
	meta := meta.Must(m)
	client := inst.Client(meta)
	logger := meta.Log("APPSEC", "dataSourceCustomDenyRead")

	getCustomDeny := appsec.GetCustomDenyListRequest{}

	configID, err := tf.GetIntValue("config_id", d)
	if err != nil {
		return diag.FromErr(err)
	}
	getCustomDeny.ConfigID = configID

	if getCustomDeny.Version, err = getLatestConfigVersion(ctx, configID, m); err != nil {
		return diag.FromErr(err)
	}

	customDenyID, err := tf.GetStringValue("custom_deny_id", d)
	if err != nil && !errors.Is(err, tf.ErrNotFound) {
		return diag.FromErr(err)
	}
	getCustomDeny.ID = customDenyID

	customdeny, err := client.GetCustomDenyList(ctx, getCustomDeny)
	if err != nil {
		logger.Errorf("calling 'getCustomDeny': %s", err.Error())
		return diag.FromErr(err)
	}

	ots := OutputTemplates{}
	InitTemplates(ots)

	outputtext, err := RenderTemplates(ots, "customDenyDS", customdeny)
	if err != nil {
		return diag.FromErr(err)
	}
	if err := d.Set("output_text", outputtext); err != nil {
		return diag.Errorf("%s: %s", tf.ErrValueSet, err.Error())
	}

	jsonBody, err := json.Marshal(customdeny)
	if err != nil {
		return diag.FromErr(err)
	}

	if err := d.Set("json", string(jsonBody)); err != nil {
		return diag.Errorf("%s: %s", tf.ErrValueSet, err.Error())
	}

	if len(customdeny.CustomDenyList) > 0 {

		d.SetId(string(customdeny.CustomDenyList[0].ID))
	}

	return nil
}
