package appsec

import (
	"context"
	"encoding/json"
	"errors"
	"strconv"

	"github.com/akamai/AkamaiOPEN-edgegrid-golang/v6/pkg/appsec"
	"github.com/akamai/terraform-provider-akamai/v4/pkg/akamai"
	"github.com/akamai/terraform-provider-akamai/v4/pkg/common/tf"
	"github.com/hashicorp/terraform-plugin-sdk/v2/diag"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/schema"
)

func dataSourceReputationProfileActions() *schema.Resource {
	return &schema.Resource{
		ReadContext: dataSourceReputationProfileActionsRead,
		Schema: map[string]*schema.Schema{
			"config_id": {
				Type:        schema.TypeInt,
				Required:    true,
				Description: "Unique identifier of the security configuration",
			},
			"security_policy_id": {
				Type:        schema.TypeString,
				Required:    true,
				Description: "Unique identifier of the security policy",
			},
			"reputation_profile_id": {
				Type:        schema.TypeInt,
				Optional:    true,
				Description: "Unique identifier of a specific reputation profile for which to retrieve information",
			},
			"action": {
				Type:        schema.TypeString,
				Computed:    true,
				Description: "Action to be taken when the specified reputation profile is triggered",
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

func dataSourceReputationProfileActionsRead(ctx context.Context, d *schema.ResourceData, m interface{}) diag.Diagnostics {
	meta := akamai.Meta(m)
	client := inst.Client(meta)
	logger := meta.Log("APPSEC", "dataSourceReputationProfileActionsRead")

	getReputationProfileActions := appsec.GetReputationProfileActionsRequest{}

	configID, err := tf.GetIntValue("config_id", d)
	if err != nil {
		return diag.FromErr(err)
	}
	getReputationProfileActions.ConfigID = configID

	if getReputationProfileActions.Version, err = getLatestConfigVersion(ctx, configID, m); err != nil {
		return diag.FromErr(err)
	}

	policyID, err := tf.GetStringValue("security_policy_id", d)
	if err != nil {
		return diag.FromErr(err)
	}
	getReputationProfileActions.PolicyID = policyID

	reputationProfileID, err := tf.GetIntValue("reputation_profile_id", d)
	if err != nil && !errors.Is(err, tf.ErrNotFound) {
		return diag.FromErr(err)
	}
	getReputationProfileActions.ReputationProfileID = reputationProfileID

	reputationprofileactions, err := client.GetReputationProfileActions(ctx, getReputationProfileActions)
	if err != nil {
		logger.Errorf("calling 'getReputationProfileActions': %s", err.Error())
		return diag.FromErr(err)
	}

	ots := OutputTemplates{}
	InitTemplates(ots)

	outputtext, err := RenderTemplates(ots, "reputationProfilesActions", reputationprofileactions)
	if err != nil {
		return diag.FromErr(err)
	}
	if err := d.Set("output_text", outputtext); err != nil {
		return diag.Errorf("%s: %s", tf.ErrValueSet, err.Error())
	}

	jsonBody, err := json.Marshal(reputationprofileactions)
	if err != nil {
		return diag.FromErr(err)
	}

	if err := d.Set("json", string(jsonBody)); err != nil {
		return diag.Errorf("%s: %s", tf.ErrValueSet, err.Error())
	}

	if len(reputationprofileactions.ReputationProfiles) > 0 {
		if err := d.Set("action", reputationprofileactions.ReputationProfiles[0].Action); err != nil {
			return diag.Errorf("%s: %s", tf.ErrValueSet, err.Error())
		}
	}

	d.SetId(strconv.Itoa(getReputationProfileActions.ConfigID))

	return nil
}
