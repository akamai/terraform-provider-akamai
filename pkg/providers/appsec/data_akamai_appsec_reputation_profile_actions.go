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

func dataSourceReputationProfileActions() *schema.Resource {
	return &schema.Resource{
		ReadContext: dataSourceReputationProfileActionsRead,
		Schema: map[string]*schema.Schema{
			"config_id": {
				Type:     schema.TypeInt,
				Required: true,
			},
			"version": {
				Type:     schema.TypeInt,
				Required: true,
			},
			"security_policy_id": {
				Type:     schema.TypeString,
				Required: true,
			},
			"reputation_profile_id": {
				Type:     schema.TypeInt,
				Optional: true,
			},
			"action": {
				Type:     schema.TypeString,
				Computed: true,
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

func dataSourceReputationProfileActionsRead(ctx context.Context, d *schema.ResourceData, m interface{}) diag.Diagnostics {
	meta := akamai.Meta(m)
	client := inst.Client(meta)
	logger := meta.Log("APPSEC", "resourceReputationProfileActionsRead")

	getReputationProfileActions := appsec.GetReputationProfileActionsRequest{}

	configid, err := tools.GetIntValue("config_id", d)
	if err != nil && !errors.Is(err, tools.ErrNotFound) {
		return diag.FromErr(err)
	}
	getReputationProfileActions.ConfigID = configid

	version, err := tools.GetIntValue("version", d)
	if err != nil && !errors.Is(err, tools.ErrNotFound) {
		return diag.FromErr(err)
	}
	getReputationProfileActions.Version = version

	policyid, err := tools.GetStringValue("security_policy_id", d)
	if err != nil && !errors.Is(err, tools.ErrNotFound) {
		return diag.FromErr(err)
	}
	getReputationProfileActions.PolicyID = policyid

	reputationprofileid, err := tools.GetIntValue("reputation_profile_id", d)
	if err != nil && !errors.Is(err, tools.ErrNotFound) {
		return diag.FromErr(err)
	}
	getReputationProfileActions.ReputationProfileID = reputationprofileid

	reputationprofileactions, err := client.GetReputationProfileActions(ctx, getReputationProfileActions)
	if err != nil {
		logger.Errorf("calling 'getReputationProfileActions': %s", err.Error())
		return diag.FromErr(err)
	}

	ots := OutputTemplates{}
	InitTemplates(ots)

	outputtext, err := RenderTemplates(ots, "reputationProfilesActions", reputationprofileactions)
	if err == nil {
		d.Set("output_text", outputtext)
	}

	jsonBody, err := json.Marshal(reputationprofileactions)
	if err != nil {
		return diag.FromErr(err)
	}

	d.Set("json", string(jsonBody))

	if len(reputationprofileactions.ReputationProfiles) > 0 {
		d.Set("action", reputationprofileactions.ReputationProfiles[0].Action)
	}

	d.SetId(strconv.Itoa(getReputationProfileActions.ConfigID))

	return nil
}
