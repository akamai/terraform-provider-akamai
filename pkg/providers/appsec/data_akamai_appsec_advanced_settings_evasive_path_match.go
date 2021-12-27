package appsec

import (
	"context"
	"encoding/json"
	"errors"
	"fmt"
	"strconv"

	"github.com/akamai/AkamaiOPEN-edgegrid-golang/v2/pkg/appsec"
	"github.com/akamai/terraform-provider-akamai/v2/pkg/akamai"
	"github.com/akamai/terraform-provider-akamai/v2/pkg/tools"

	"github.com/hashicorp/terraform-plugin-sdk/v2/diag"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/schema"
)

func dataSourceAdvancedSettingsEvasivePathMatch() *schema.Resource {
	return &schema.Resource{
		ReadContext: dataSourceAdvancedSettingsEvasivePathMatchRead,
		Schema: map[string]*schema.Schema{
			"config_id": {
				Type:        schema.TypeInt,
				Required:    true,
				Description: "config ID",
			},
			"security_policy_id": {
				Type:        schema.TypeString,
				Optional:    true,
				Description: "security policy ID",
			},
			"json": {
				Type:        schema.TypeString,
				Computed:    true,
				Description: "JSON representation",
			},
			"output_text": {
				Type:        schema.TypeString,
				Computed:    true,
				Description: "Text Export representation",
			},
		},
	}
}

func dataSourceAdvancedSettingsEvasivePathMatchRead(ctx context.Context, d *schema.ResourceData, m interface{}) diag.Diagnostics {
	meta := akamai.Meta(m)
	client := inst.Client(meta)
	logger := meta.Log("APPSEC", "dataSourceAdvancedSettingsEvasivePathMatchRead")

	getAdvancedSettingsEvasivePathMatch := appsec.GetAdvancedSettingsEvasivePathMatchRequest{}

	configID, err := tools.GetIntValue("config_id", d)
	if err != nil && !errors.Is(err, tools.ErrNotFound) {
		return diag.FromErr(err)
	}
	getAdvancedSettingsEvasivePathMatch.ConfigID = configID

	getAdvancedSettingsEvasivePathMatch.Version = getLatestConfigVersion(ctx, configID, m)

	policyid, err := tools.GetStringValue("security_policy_id", d)
	if err != nil && !errors.Is(err, tools.ErrNotFound) {
		return diag.FromErr(err)
	}
	getAdvancedSettingsEvasivePathMatch.PolicyID = policyid

	advancedsettingsevasivepathmatch, err := client.GetAdvancedSettingsEvasivePathMatch(ctx, getAdvancedSettingsEvasivePathMatch)
	if err != nil {
		logger.Errorf("calling 'getAdvancedSettingsEvasivePathMatch': %s", err.Error())
		return diag.FromErr(err)
	}

	ots := OutputTemplates{}
	InitTemplates(ots)

	outputtext, err := RenderTemplates(ots, "advancedSettingsEvasivePathMatchDS", advancedsettingsevasivepathmatch)
	if err == nil {
		d.Set("output_text", outputtext)
	}

	jsonBody, err := json.Marshal(advancedsettingsevasivepathmatch)
	if err != nil {
		return diag.FromErr(err)
	}

	if err := d.Set("json", string(jsonBody)); err != nil {
		return diag.FromErr(fmt.Errorf("%w: %s", tools.ErrValueSet, err.Error()))
	}

	d.SetId(strconv.Itoa(getAdvancedSettingsEvasivePathMatch.ConfigID))

	return nil
}
