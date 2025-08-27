package appsec

import (
	"context"
	"encoding/json"
	"strconv"

	"github.com/akamai/AkamaiOPEN-edgegrid-golang/v12/pkg/appsec"
	"github.com/akamai/terraform-provider-akamai/v8/pkg/common/tf"
	"github.com/akamai/terraform-provider-akamai/v8/pkg/meta"
	"github.com/hashicorp/terraform-plugin-sdk/v2/diag"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/schema"
)

func dataSourcePolicyProtections() *schema.Resource {
	return &schema.Resource{
		ReadContext: dataSourcePolicyProtectionsRead,
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
			"apply_api_constraints": {
				Type:        schema.TypeBool,
				Computed:    true,
				Description: "Whether to enable API constraints",
			},
			"apply_application_layer_controls": {
				Type:        schema.TypeBool,
				Computed:    true,
				Description: "Whether to enable application layer controls",
			},
			"apply_botman_controls": {
				Type:        schema.TypeBool,
				Computed:    true,
				Description: "Whether to enable botman controls",
			},
			"apply_malware_controls": {
				Type:        schema.TypeBool,
				Computed:    true,
				Description: "Whether to enable malware controls",
			},
			"apply_network_layer_controls": {
				Type:        schema.TypeBool,
				Computed:    true,
				Description: "Whether to enable network layer controls",
			},
			"apply_rate_controls": {
				Type:        schema.TypeBool,
				Computed:    true,
				Description: "Whether to enable rate controls",
			},
			"apply_reputation_controls": {
				Type:        schema.TypeBool,
				Computed:    true,
				Description: "Whether to enable reputation controls",
			},
			"apply_slow_post_controls": {
				Type:        schema.TypeBool,
				Computed:    true,
				Description: "Whether to enable slow post controls",
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

func dataSourcePolicyProtectionsRead(ctx context.Context, d *schema.ResourceData, m interface{}) diag.Diagnostics {
	meta := meta.Must(m)
	client := inst.Client(meta)
	logger := meta.Log("APPSEC", "dataSourcePolicyProtectionsRead")

	getPolicyProtections := appsec.GetPolicyProtectionsRequest{}

	configID, err := tf.GetIntValue("config_id", d)
	if err != nil {
		return diag.FromErr(err)
	}
	getPolicyProtections.ConfigID = configID

	if getPolicyProtections.Version, err = getLatestConfigVersion(ctx, configID, m); err != nil {
		return diag.FromErr(err)
	}

	policyID, err := tf.GetStringValue("security_policy_id", d)
	if err != nil {
		return diag.FromErr(err)
	}
	getPolicyProtections.PolicyID = policyID

	policyprotections, err := client.GetPolicyProtections(ctx, getPolicyProtections)
	if err != nil {
		logger.Errorf("calling 'getPolicyProtections': %s", err.Error())
		return diag.FromErr(err)
	}

	ots := OutputTemplates{}
	InitTemplates(ots)

	outputtext, err := RenderTemplates(ots, "protections", policyprotections)
	if err != nil {
		return diag.FromErr(err)
	}
	if err := d.Set("output_text", outputtext); err != nil {
		return diag.Errorf("%s: %s", tf.ErrValueSet, err.Error())
	}

	jsonBody, err := json.Marshal(policyprotections)
	if err != nil {
		return diag.FromErr(err)
	}

	if err := d.Set("json", string(jsonBody)); err != nil {
		return diag.Errorf("%s: %s", tf.ErrValueSet, err.Error())
	}

	if err := d.Set("apply_api_constraints", policyprotections.ApplyAPIConstraints); err != nil {
		return diag.Errorf("%s: %s", tf.ErrValueSet, err.Error())
	}

	if err := d.Set("apply_botman_controls", policyprotections.ApplyBotmanControls); err != nil {
		return diag.Errorf("%s: %s", tf.ErrValueSet, err.Error())
	}

	if err := d.Set("apply_application_layer_controls", policyprotections.ApplyApplicationLayerControls); err != nil {
		return diag.Errorf("%s: %s", tf.ErrValueSet, err.Error())
	}

	if err := d.Set("apply_malware_controls", policyprotections.ApplyMalwareControls); err != nil {
		return diag.Errorf("%s: %s", tf.ErrValueSet, err.Error())
	}

	if err := d.Set("apply_network_layer_controls", policyprotections.ApplyNetworkLayerControls); err != nil {
		return diag.Errorf("%s: %s", tf.ErrValueSet, err.Error())
	}

	if err := d.Set("apply_rate_controls", policyprotections.ApplyRateControls); err != nil {
		return diag.Errorf("%s: %s", tf.ErrValueSet, err.Error())
	}

	if err := d.Set("apply_reputation_controls", policyprotections.ApplyReputationControls); err != nil {
		return diag.Errorf("%s: %s", tf.ErrValueSet, err.Error())
	}

	if err := d.Set("apply_slow_post_controls", policyprotections.ApplySlowPostControls); err != nil {
		return diag.Errorf("%s: %s", tf.ErrValueSet, err.Error())
	}

	d.SetId(strconv.Itoa(getPolicyProtections.ConfigID))

	return nil
}
