package appsec

import (
	"context"
	"errors"
	"strconv"

	"github.com/akamai/AkamaiOPEN-edgegrid-golang/v2/pkg/appsec"
	"github.com/akamai/terraform-provider-akamai/v2/pkg/akamai"
	"github.com/akamai/terraform-provider-akamai/v2/pkg/tools"
	"github.com/hashicorp/terraform-plugin-sdk/v2/diag"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/schema"
)

func dataSourceRatePolicyActions() *schema.Resource {
	return &schema.Resource{
		ReadContext: dataSourceRatePolicyActionsRead,
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
			"rate_policy_id": {
				Type:        schema.TypeInt,
				Optional:    true,
				Description: "Unique identifier of a specific rate policy for which to retrieve information",
			},
			"output_text": {
				Type:        schema.TypeString,
				Computed:    true,
				Description: "Text representation",
			},
		},
	}
}

func dataSourceRatePolicyActionsRead(ctx context.Context, d *schema.ResourceData, m interface{}) diag.Diagnostics {
	meta := akamai.Meta(m)
	client := inst.Client(meta)
	logger := meta.Log("APPSEC", "dataSourceRatePolicyActionsRead")

	getRatePolicyActions := appsec.GetRatePolicyActionsRequest{}

	configID, err := tools.GetIntValue("config_id", d)
	if err != nil {
		return diag.FromErr(err)
	}
	getRatePolicyActions.ConfigID = configID

	if getRatePolicyActions.Version, err = getLatestConfigVersion(ctx, configID, m); err != nil {
		return diag.FromErr(err)
	}

	policyID, err := tools.GetStringValue("security_policy_id", d)
	if err != nil {
		return diag.FromErr(err)
	}
	getRatePolicyActions.PolicyID = policyID

	ratePolicyID, err := tools.GetIntValue("rate_policy_id", d)
	if err != nil && !errors.Is(err, tools.ErrNotFound) {
		return diag.FromErr(err)
	}
	getRatePolicyActions.RatePolicyID = ratePolicyID

	ratepolicyactions, err := client.GetRatePolicyActions(ctx, getRatePolicyActions)
	if err != nil {
		logger.Errorf("calling 'getRatePolicyActions': %s", err.Error())
		return diag.FromErr(err)
	}

	for _, configval := range ratepolicyactions.RatePolicyActions {
		d.SetId(strconv.Itoa(configval.ID))
	}

	ots := OutputTemplates{}
	InitTemplates(ots)

	outputtext, err := RenderTemplates(ots, "ratePolicyActions", ratepolicyactions)
	if err != nil {
		return diag.FromErr(err)
	}
	if err := d.Set("output_text", outputtext); err != nil {
		return diag.Errorf("%s: %s", tools.ErrValueSet, err.Error())
	}

	return nil
}
