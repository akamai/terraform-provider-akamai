package appsec

import (
	"context"
	"encoding/json"
	"errors"
	"strconv"

	"github.com/akamai/AkamaiOPEN-edgegrid-golang/v8/pkg/appsec"
	"github.com/akamai/terraform-provider-akamai/v5/pkg/common/tf"
	"github.com/akamai/terraform-provider-akamai/v5/pkg/meta"
	"github.com/hashicorp/terraform-plugin-sdk/v2/diag"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/schema"
)

// MatchTargetOutputText holds data for templates
type MatchTargetOutputText struct {
	TargetID int
	PolicyID string
	Type     string
}

// Definition of constant variables
const (
	APITarget     = "API"
	WebsiteTarget = "Website"
)

func dataSourceMatchTargets() *schema.Resource {
	return &schema.Resource{
		ReadContext: dataSourceMatchTargetsRead,
		Schema: map[string]*schema.Schema{
			"config_id": {
				Type:        schema.TypeInt,
				Required:    true,
				Description: "Unique identifier of the security configuration",
			},
			"match_target_id": {
				Type:        schema.TypeInt,
				Optional:    true,
				Description: "Unique identifier of a specific match target for which to retrieve information",
			},
			"json": {
				Type:        schema.TypeString,
				Computed:    true,
				Description: "JSON Export representation",
			},
			"output_text": {
				Type:        schema.TypeString,
				Computed:    true,
				Description: "Text representation",
			},
		},
	}
}

func dataSourceMatchTargetsRead(ctx context.Context, d *schema.ResourceData, m interface{}) diag.Diagnostics {
	meta := meta.Must(m)
	client := inst.Client(meta)
	logger := meta.Log("APPSEC", "dataSourceMatchTargetsRead")

	getMatchTargets := appsec.GetMatchTargetsRequest{}

	configID, err := tf.GetIntValue("config_id", d)
	if err != nil {
		return diag.FromErr(err)
	}
	getMatchTargets.ConfigID = configID

	if getMatchTargets.ConfigVersion, err = getLatestConfigVersion(ctx, configID, m); err != nil {
		return diag.FromErr(err)
	}

	matchTargetID, err := tf.GetIntValue("match_target_id", d)
	if err != nil && !errors.Is(err, tf.ErrNotFound) {
		return diag.FromErr(err)
	}
	getMatchTargets.TargetID = matchTargetID

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
		return diag.Errorf("%s: %s", tf.ErrValueSet, err.Error())
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
	if err != nil {
		return diag.FromErr(err)
	}
	if err := d.Set("output_text", websiteMatchTargetsText); err != nil {
		return diag.Errorf("%s: %s", tf.ErrValueSet, err.Error())
	}

	d.SetId(strconv.Itoa(getMatchTargets.ConfigID))

	return nil
}
