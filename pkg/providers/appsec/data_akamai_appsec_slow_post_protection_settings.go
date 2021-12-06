package appsec

import (
	"context"
	"encoding/json"
	"errors"
	"strconv"

	"github.com/akamai/AkamaiOPEN-edgegrid-golang/v2/pkg/appsec"
	"github.com/akamai/terraform-provider-akamai/v2/pkg/akamai"
	"github.com/akamai/terraform-provider-akamai/v2/pkg/tools"
	"github.com/hashicorp/terraform-plugin-sdk/v2/diag"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/schema"
)

func dataSourceSlowPostProtectionSettings() *schema.Resource {
	return &schema.Resource{
		ReadContext: dataSourceSlowPostProtectionSettingsRead,
		Schema: map[string]*schema.Schema{
			"config_id": {
				Type:     schema.TypeInt,
				Required: true,
			},
			"security_policy_id": {
				Type:     schema.TypeString,
				Required: true,
			},
			"json": {
				Type:     schema.TypeString,
				Computed: true,
			},
			"output_text": {
				Type:        schema.TypeString,
				Computed:    true,
				Description: "Text Export representation",
			},
		},
	}
}

func dataSourceSlowPostProtectionSettingsRead(ctx context.Context, d *schema.ResourceData, m interface{}) diag.Diagnostics {
	meta := akamai.Meta(m)
	client := inst.Client(meta)
	logger := meta.Log("APPSEC", "dataSourceSlowPostProtectionSettingsRead")

	getSlowPostProtectionSettings := appsec.GetSlowPostProtectionSettingsRequest{}

	configID, err := tools.GetIntValue("config_id", d)
	if err != nil && !errors.Is(err, tools.ErrNotFound) {
		return diag.FromErr(err)
	}
	getSlowPostProtectionSettings.ConfigID = configID

	getSlowPostProtectionSettings.Version = getLatestConfigVersion(ctx, configID, m)

	policyID, err := tools.GetStringValue("security_policy_id", d)
	if err != nil && !errors.Is(err, tools.ErrNotFound) {
		return diag.FromErr(err)
	}
	getSlowPostProtectionSettings.PolicyID = policyID

	slowpostprotectionsettings, err := client.GetSlowPostProtectionSettings(ctx, getSlowPostProtectionSettings)
	if err != nil {
		logger.Errorf("calling 'getSlowPostProtectionSettings': %s", err.Error())
		return diag.FromErr(err)
	}

	ots := OutputTemplates{}
	InitTemplates(ots)

	outputtext, err := RenderTemplates(ots, "slowPostDS", slowpostprotectionsettings)
	if err == nil {
		if err := d.Set("output_text", outputtext); err != nil {
			return diag.Errorf("%s: %s", tools.ErrValueSet, err.Error())
		}
	}

	jsonBody, err := json.Marshal(slowpostprotectionsettings)
	if err != nil {
		return diag.FromErr(err)
	}

	if err := d.Set("json", string(jsonBody)); err != nil {
		return diag.Errorf("%s: %s", tools.ErrValueSet, err.Error())
	}

	d.SetId(strconv.Itoa(getSlowPostProtectionSettings.ConfigID))

	return nil
}
