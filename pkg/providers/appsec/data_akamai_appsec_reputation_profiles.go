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

func dataSourceReputationProfiles() *schema.Resource {
	return &schema.Resource{
		ReadContext: dataSourceReputationProfilesRead,
		Schema: map[string]*schema.Schema{
			"config_id": {
				Type:     schema.TypeInt,
				Required: true,
			},
			"reputation_profile_id": {
				Type:     schema.TypeInt,
				Optional: true,
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

func dataSourceReputationProfilesRead(ctx context.Context, d *schema.ResourceData, m interface{}) diag.Diagnostics {
	meta := akamai.Meta(m)
	client := inst.Client(meta)
	logger := meta.Log("APPSEC", "dataSourceReputationProfilesRead")

	getReputationProfiles := appsec.GetReputationProfilesRequest{}

	configID, err := tools.GetIntValue("config_id", d)
	if err != nil && !errors.Is(err, tools.ErrNotFound) {
		return diag.FromErr(err)
	}
	getReputationProfiles.ConfigID = configID

	getReputationProfiles.ConfigVersion = getLatestConfigVersion(ctx, configID, m)

	reputationprofileid, err := tools.GetIntValue("reputation_profile_id", d)
	if err != nil && !errors.Is(err, tools.ErrNotFound) {
		return diag.FromErr(err)
	}
	getReputationProfiles.ReputationProfileId = reputationprofileid

	reputationprofiles, err := client.GetReputationProfiles(ctx, getReputationProfiles)
	if err != nil {
		logger.Errorf("calling 'getReputationProfiles': %s", err.Error())
		return diag.FromErr(err)
	}

	ots := OutputTemplates{}
	InitTemplates(ots)

	outputtext, err := RenderTemplates(ots, "reputationProfilesDS", reputationprofiles)
	if err == nil {
		if err := d.Set("output_text", outputtext); err != nil {
			return diag.Errorf("%s: %s", tools.ErrValueSet, err.Error())
		}
	}

	jsonBody, err := json.Marshal(reputationprofiles)
	if err != nil {
		return diag.FromErr(err)
	}

	if err := d.Set("json", string(jsonBody)); err != nil {
		return diag.Errorf("%s: %s", tools.ErrValueSet, err.Error())
	}

	d.SetId(strconv.Itoa(configID))

	return nil
}
