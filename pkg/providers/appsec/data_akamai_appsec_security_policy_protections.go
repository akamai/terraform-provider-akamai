package appsec

import (
	"context"
	"encoding/json"
	"errors"
	"strconv"

	v2 "github.com/akamai/AkamaiOPEN-edgegrid-golang/v2/pkg/appsec"
	"github.com/akamai/terraform-provider-akamai/v2/pkg/akamai"
	"github.com/akamai/terraform-provider-akamai/v2/pkg/tools"
	"github.com/hashicorp/terraform-plugin-sdk/v2/diag"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/schema"
)

func dataSourcePolicyProtections() *schema.Resource {
	return &schema.Resource{
		ReadContext: dataSourcePolicyProtectionsRead,
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
			"apply_application_layer_controls": {
				Type:     schema.TypeBool,
				Computed: true,
			},
			"apply_network_layer_controls": {
				Type:     schema.TypeBool,
				Computed: true,
			},
			"apply_rate_controls": {
				Type:     schema.TypeBool,
				Computed: true,
			},
			"apply_reputation_controls": {
				Type:     schema.TypeBool,
				Computed: true,
			},
			"apply_botman_controls": {
				Type:     schema.TypeBool,
				Computed: true,
			},
			"apply_api_constraints": {
				Type:     schema.TypeBool,
				Computed: true,
			},
			"apply_slow_post_controls": {
				Type:     schema.TypeBool,
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

func dataSourcePolicyProtectionsRead(ctx context.Context, d *schema.ResourceData, m interface{}) diag.Diagnostics {
	meta := akamai.Meta(m)
	client := inst.Client(meta)
	logger := meta.Log("APPSEC", "resourcePolicyProtectionsRead")

	getPolicyProtections := v2.GetPolicyProtectionsRequest{}

	configid, err := tools.GetIntValue("config_id", d)
	if err != nil && !errors.Is(err, tools.ErrNotFound) {
		return diag.FromErr(err)
	}
	getPolicyProtections.ConfigID = configid

	version, err := tools.GetIntValue("version", d)
	if err != nil && !errors.Is(err, tools.ErrNotFound) {
		return diag.FromErr(err)
	}
	getPolicyProtections.Version = version

	policyid, err := tools.GetStringValue("security_policy_id", d)
	if err != nil && !errors.Is(err, tools.ErrNotFound) {
		return diag.FromErr(err)
	}
	getPolicyProtections.PolicyID = policyid

	policyprotections, err := client.GetPolicyProtections(ctx, getPolicyProtections)
	if err != nil {
		logger.Errorf("calling 'getPolicyProtections': %s", err.Error())
		return diag.FromErr(err)
	}

	ots := OutputTemplates{}
	InitTemplates(ots)

	outputtext, err := RenderTemplates(ots, "wafProtectionDS", policyprotections)
	if err == nil {
		d.Set("output_text", outputtext)
	}

	jsonBody, err := json.Marshal(policyprotections)
	if err != nil {
		return diag.FromErr(err)
	}

	d.Set("json", string(jsonBody))

	d.Set("apply_application_layer_controls", policyprotections.ApplyApplicationLayerControls)
	d.Set("apply_network_layer_controls", policyprotections.ApplyNetworkLayerControls)
	d.Set("apply_api_constraints", policyprotections.ApplyAPIConstraints)
	d.Set("apply_rate_controls", policyprotections.ApplyRateControls)
	d.Set("apply_reputation_controls", policyprotections.ApplyReputationControls)
	d.Set("apply_botman_controls", policyprotections.ApplyBotmanControls)
	d.Set("apply_slow_post_controls", policyprotections.ApplySlowPostControls)

	d.SetId(strconv.Itoa(getPolicyProtections.ConfigID))

	return nil
}
