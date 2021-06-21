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

type MatchTargetOutputText struct {
	TargetID int
	PolicyID string
	Type     string
}

const (
	APITarget     = "API"
	WebsiteTarget = "Website"
)

func dataSourceMatchTargets() *schema.Resource {
	return &schema.Resource{
		ReadContext: dataSourceMatchTargetsRead,
		Schema: map[string]*schema.Schema{
			"config_id": {
				Type:     schema.TypeInt,
				Required: true,
			},
			"match_target_id": {
				Type:     schema.TypeInt,
				Optional: true,
			},
			"output_text": {
				Type:        schema.TypeString,
				Computed:    true,
				Description: "Text Export representation",
			},
			"json": {
				Type:        schema.TypeString,
				Computed:    true,
				Description: "JSON Export representation",
			},
		},
	}
}

func dataSourceMatchTargetsRead(ctx context.Context, d *schema.ResourceData, m interface{}) diag.Diagnostics {
	meta := akamai.Meta(m)
	client := inst.Client(meta)
	logger := meta.Log("APPSEC", "dataSourceMatchTargetsRead")

	getMatchTargets := appsec.GetMatchTargetsRequest{}

	configid, err := tools.GetIntValue("config_id", d)
	if err != nil && !errors.Is(err, tools.ErrNotFound) {
		return diag.FromErr(err)
	}
	getMatchTargets.ConfigID = configid

	getMatchTargets.ConfigVersion = getLatestConfigVersion(ctx, configid, m)

	matchtargetid, err := tools.GetIntValue("match_target_id", d)
	if err != nil && !errors.Is(err, tools.ErrNotFound) {
		return diag.FromErr(err)
	}
	getMatchTargets.TargetID = matchtargetid

	matchtargets, err := client.GetMatchTargets(ctx, getMatchTargets)
	if err != nil {
		logger.Errorf("calling 'getMatchTargets': %s", err.Error())
		return diag.FromErr(err)
	}

	jsonBody, err := json.Marshal(matchtargets)
	if err != nil {
		logger.Errorf("calling 'getMatchTargets': %s", err.Error())
		return diag.FromErr(err)
	}
	if err := d.Set("json", string(jsonBody)); err != nil {
		return diag.FromErr(fmt.Errorf("%w: %s", tools.ErrValueSet, err.Error()))
	}

	ots := OutputTemplates{}
	InitTemplates(ots)

	matchtargetCount := len(matchtargets.MatchTargets.WebsiteTargets) + len(matchtargets.MatchTargets.APITargets)
	matchtargetsOutputText := make([]MatchTargetOutputText, 0, matchtargetCount)
	for _, value := range matchtargets.MatchTargets.WebsiteTargets {
		matchtargetsOutputText = append(matchtargetsOutputText, MatchTargetOutputText{value.TargetID, value.SecurityPolicy.PolicyID, WebsiteTarget})
	}
	for _, value := range matchtargets.MatchTargets.APITargets {
		matchtargetsOutputText = append(matchtargetsOutputText, MatchTargetOutputText{value.TargetID, value.SecurityPolicy.PolicyID, APITarget})
	}
	websiteMatchTargetsText, err := RenderTemplates(ots, "matchTargetDS", matchtargetsOutputText)
	if err == nil {
		if err := d.Set("output_text", websiteMatchTargetsText); err != nil {
			return diag.FromErr(fmt.Errorf("%w: %s", tools.ErrValueSet, err.Error()))
		}
	}

	d.SetId(strconv.Itoa(getMatchTargets.ConfigID))

	return nil
}
