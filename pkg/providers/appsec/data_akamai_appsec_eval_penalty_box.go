package appsec

import (
	"context"
	"strconv"

	"github.com/akamai/AkamaiOPEN-edgegrid-golang/v3/pkg/appsec"
	"github.com/akamai/terraform-provider-akamai/v3/pkg/akamai"
	"github.com/akamai/terraform-provider-akamai/v3/pkg/tools"
	"github.com/hashicorp/terraform-plugin-sdk/v2/diag"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/schema"
)

func dataSourceEvalPenaltyBox() *schema.Resource {
	return &schema.Resource{
		ReadContext: dataSourceEvalPenaltyBoxRead,
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
			"action": {
				Type:        schema.TypeString,
				Computed:    true,
				Description: "Action to be applied to requests from clients in the penalty box",
			},
			"enabled": {
				Type:        schema.TypeBool,
				Computed:    true,
				Description: "Whether the penalty box is enabled for the specified security policy",
			},
			"output_text": {
				Type:        schema.TypeString,
				Computed:    true,
				Description: "Text output in tabular form",
			},
		},
	}
}

func dataSourceEvalPenaltyBoxRead(ctx context.Context, d *schema.ResourceData, m interface{}) diag.Diagnostics {
	meta := akamai.Meta(m)
	client := inst.Client(meta)
	logger := meta.Log("APPSEC", "dataSourceEvalPenaltyBoxRead")

	getPenaltyBox := appsec.GetPenaltyBoxRequest{}

	configID, err := tools.GetIntValue("config_id", d)
	if err != nil {
		return diag.FromErr(err)
	}
	getPenaltyBox.ConfigID = configID

	if getPenaltyBox.Version, err = getLatestConfigVersion(ctx, configID, m); err != nil {
		return diag.FromErr(err)
	}

	policyID, err := tools.GetStringValue("security_policy_id", d)
	if err != nil {
		return diag.FromErr(err)
	}
	getPenaltyBox.PolicyID = policyID

	penaltybox, err := client.GetEvalPenaltyBox(ctx, getPenaltyBox)
	if err != nil {
		logger.Errorf("calling 'getEvalPenaltyBox': %s", err.Error())
		return diag.FromErr(err)
	}

	ots := OutputTemplates{}
	InitTemplates(ots)

	outputtext, err := RenderTemplates(ots, "evalPenaltyBoxDS", penaltybox)
	if err != nil {
		return diag.FromErr(err)
	}
	if err := d.Set("output_text", outputtext); err != nil {
		return diag.Errorf("%s: %s", tools.ErrValueSet, err.Error())
	}

	if err := d.Set("action", penaltybox.Action); err != nil {
		return diag.Errorf("%s: %s", tools.ErrValueSet, err.Error())
	}

	if err := d.Set("enabled", penaltybox.PenaltyBoxProtection); err != nil {
		return diag.Errorf("%s: %s", tools.ErrValueSet, err.Error())
	}

	d.SetId(strconv.Itoa(getPenaltyBox.ConfigID))

	return nil
}
