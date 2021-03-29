package appsec

import (
	"context"
	"errors"
	"fmt"
	"strconv"
	"strings"

	"github.com/akamai/AkamaiOPEN-edgegrid-golang/v2/pkg/appsec"
	"github.com/akamai/terraform-provider-akamai/v2/pkg/akamai"
	"github.com/akamai/terraform-provider-akamai/v2/pkg/tools"
	"github.com/hashicorp/terraform-plugin-sdk/v2/diag"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/schema"
)

// appsec v1
//
// https://developer.akamai.com/api/cloud_security/application_security/v1.html
func resourcePolicyProtections() *schema.Resource {
	return &schema.Resource{
		CreateContext: resourcePolicyProtectionsUpdate,
		ReadContext:   resourcePolicyProtectionsRead,
		UpdateContext: resourcePolicyProtectionsUpdate,
		DeleteContext: resourcePolicyProtectionsDelete,
		Importer: &schema.ResourceImporter{
			State: schema.ImportStatePassthrough,
		},
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
				Required: true,
			},
			"apply_network_layer_controls": {
				Type:     schema.TypeBool,
				Required: true,
			},
			"apply_rate_controls": {
				Type:     schema.TypeBool,
				Required: true,
			},
			"apply_reputation_controls": {
				Type:     schema.TypeBool,
				Required: true,
			},
			"apply_botman_controls": {
				Type:     schema.TypeBool,
				Required: true,
			},
			"apply_api_constraints": {
				Type:     schema.TypeBool,
				Required: true,
			},
			"apply_slow_post_controls": {
				Type:     schema.TypeBool,
				Required: true,
			},
			"output_text": {
				Type:        schema.TypeString,
				Computed:    true,
				Description: "Text Export representation",
			},
		},
	}
}

func resourcePolicyProtectionsRead(ctx context.Context, d *schema.ResourceData, m interface{}) diag.Diagnostics {
	meta := akamai.Meta(m)
	client := inst.Client(meta)
	logger := meta.Log("APPSEC", "resourcePolicyProtectionsRead")

	getPolicyProtections := appsec.GetPolicyProtectionsRequest{}
	if d.Id() != "" && strings.Contains(d.Id(), ":") {
		s := strings.Split(d.Id(), ":")

		configid, errconv := strconv.Atoi(s[0])
		if errconv != nil {
			return diag.FromErr(errconv)
		}
		getPolicyProtections.ConfigID = configid

		version, errconv := strconv.Atoi(s[1])
		if errconv != nil {
			return diag.FromErr(errconv)
		}
		getPolicyProtections.Version = version

		if d.HasChange("version") {
			version, err := tools.GetIntValue("version", d)
			if err != nil && !errors.Is(err, tools.ErrNotFound) {
				return diag.FromErr(err)
			}
			getPolicyProtections.Version = version
		}

		policyid := s[2]
		getPolicyProtections.PolicyID = policyid

	} else {
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
	}
	policyprotections, err := client.GetPolicyProtections(ctx, getPolicyProtections)
	if err != nil {
		logger.Errorf("calling 'getPolicyProtections': %s", err.Error())
		return diag.FromErr(err)
	}

	ots := OutputTemplates{}
	InitTemplates(ots)

	outputtext, err := RenderTemplates(ots, "wafProtectionDS", policyprotections)
	if err == nil {
		if err := d.Set("output_text", outputtext); err != nil {
			return diag.FromErr(fmt.Errorf("%w: %s", tools.ErrValueSet, err.Error()))
		}
	}

	if err := d.Set("config_id", getPolicyProtections.ConfigID); err != nil {
		return diag.FromErr(fmt.Errorf("%w: %s", tools.ErrValueSet, err.Error()))
	}

	if err := d.Set("version", getPolicyProtections.Version); err != nil {
		return diag.FromErr(fmt.Errorf("%w: %s", tools.ErrValueSet, err.Error()))
	}

	if err := d.Set("security_policy_id", getPolicyProtections.PolicyID); err != nil {
		return diag.FromErr(fmt.Errorf("%w: %s", tools.ErrValueSet, err.Error()))
	}

	d.SetId(fmt.Sprintf("%d:%d:%s", getPolicyProtections.ConfigID, getPolicyProtections.Version, getPolicyProtections.PolicyID))

	return nil
}

func resourcePolicyProtectionsDelete(ctx context.Context, d *schema.ResourceData, m interface{}) diag.Diagnostics {

	meta := akamai.Meta(m)
	client := inst.Client(meta)
	logger := meta.Log("APPSEC", "resourcePolicyProtectionsRemove")

	updatePolicyProtections := appsec.UpdatePolicyProtectionsRequest{}
	removePolicyProtections := appsec.RemovePolicyProtectionsRequest{}
	if d.Id() != "" && strings.Contains(d.Id(), ":") {
		s := strings.Split(d.Id(), ":")

		configid, errconv := strconv.Atoi(s[0])
		if errconv != nil {
			return diag.FromErr(errconv)
		}
		updatePolicyProtections.ConfigID = configid
		removePolicyProtections.ConfigID = configid

		version, errconv := strconv.Atoi(s[1])
		if errconv != nil {
			return diag.FromErr(errconv)
		}
		updatePolicyProtections.Version = version
		removePolicyProtections.Version = version

		policyid := s[2]
		updatePolicyProtections.PolicyID = policyid
		removePolicyProtections.PolicyID = policyid

	} else {
		configid, err := tools.GetIntValue("config_id", d)
		if err != nil && !errors.Is(err, tools.ErrNotFound) {
			return diag.FromErr(err)
		}
		updatePolicyProtections.ConfigID = configid
		removePolicyProtections.ConfigID = configid

		version, err := tools.GetIntValue("version", d)
		if err != nil && !errors.Is(err, tools.ErrNotFound) {
			return diag.FromErr(err)
		}
		updatePolicyProtections.Version = version
		removePolicyProtections.Version = version

		policyid, err := tools.GetStringValue("security_policy_id", d)
		if err != nil && !errors.Is(err, tools.ErrNotFound) {
			return diag.FromErr(err)
		}
		updatePolicyProtections.PolicyID = policyid
		removePolicyProtections.PolicyID = policyid
	}
	//TODO remove once API fixed in Jan
	updatePolicyProtections.ApplyApplicationLayerControls = true

	updatePolicyProtections.ApplyNetworkLayerControls = false

	updatePolicyProtections.ApplyRateControls = false

	updatePolicyProtections.ApplyReputationControls = false

	updatePolicyProtections.ApplyBotmanControls = false

	updatePolicyProtections.ApplyAPIConstraints = false

	updatePolicyProtections.ApplySlowPostControls = false

	_, erru := client.UpdatePolicyProtections(ctx, updatePolicyProtections)
	if erru != nil {
		logger.Errorf("calling 'removePolicyProtections': %s", erru.Error())
		return diag.FromErr(erru)
	}

	removePolicyProtections.ApplyApplicationLayerControls = false

	removePolicyProtections.ApplyNetworkLayerControls = false

	removePolicyProtections.ApplyRateControls = false

	removePolicyProtections.ApplyReputationControls = false

	removePolicyProtections.ApplyBotmanControls = false

	removePolicyProtections.ApplyAPIConstraints = false

	removePolicyProtections.ApplySlowPostControls = false

	_, errd := client.RemovePolicyProtections(ctx, removePolicyProtections)
	if errd != nil {
		logger.Errorf("calling 'removePolicyProtections': %s", errd.Error())
		return diag.FromErr(errd)
	}

	d.SetId("")
	return nil
}

func resourcePolicyProtectionsUpdate(ctx context.Context, d *schema.ResourceData, m interface{}) diag.Diagnostics {
	meta := akamai.Meta(m)
	client := inst.Client(meta)
	logger := meta.Log("APPSEC", "resourcePolicyProtectionsUpdate")

	updatePolicyProtections := appsec.UpdatePolicyProtectionsRequest{}
	if d.Id() != "" && strings.Contains(d.Id(), ":") {
		s := strings.Split(d.Id(), ":")

		configid, errconv := strconv.Atoi(s[0])
		if errconv != nil {
			return diag.FromErr(errconv)
		}
		updatePolicyProtections.ConfigID = configid

		version, errconv := strconv.Atoi(s[1])
		if errconv != nil {
			return diag.FromErr(errconv)
		}
		updatePolicyProtections.Version = version

		if d.HasChange("version") {
			version, err := tools.GetIntValue("version", d)
			if err != nil && !errors.Is(err, tools.ErrNotFound) {
				return diag.FromErr(err)
			}
			updatePolicyProtections.Version = version
		}

		policyid := s[2]
		updatePolicyProtections.PolicyID = policyid

	} else {
		configid, err := tools.GetIntValue("config_id", d)
		if err != nil && !errors.Is(err, tools.ErrNotFound) {
			return diag.FromErr(err)
		}
		updatePolicyProtections.ConfigID = configid

		version, err := tools.GetIntValue("version", d)
		if err != nil && !errors.Is(err, tools.ErrNotFound) {
			return diag.FromErr(err)
		}
		updatePolicyProtections.Version = version

		policyid, err := tools.GetStringValue("security_policy_id", d)
		if err != nil && !errors.Is(err, tools.ErrNotFound) {
			return diag.FromErr(err)
		}
		updatePolicyProtections.PolicyID = policyid
	}
	applyapplicationlayercontrols, err := tools.GetBoolValue("apply_application_layer_controls", d)
	if err != nil && !errors.Is(err, tools.ErrNotFound) {
		return diag.FromErr(err)
	}
	updatePolicyProtections.ApplyApplicationLayerControls = applyapplicationlayercontrols

	applynetworklayercontrols, err := tools.GetBoolValue("apply_network_layer_controls", d)
	if err != nil && !errors.Is(err, tools.ErrNotFound) {
		return diag.FromErr(err)
	}
	updatePolicyProtections.ApplyNetworkLayerControls = applynetworklayercontrols

	applyratecontrols, err := tools.GetBoolValue("apply_rate_controls", d)
	if err != nil && !errors.Is(err, tools.ErrNotFound) {
		return diag.FromErr(err)
	}
	updatePolicyProtections.ApplyRateControls = applyratecontrols

	applyreputationcontrols, err := tools.GetBoolValue("apply_reputation_controls", d)
	if err != nil && !errors.Is(err, tools.ErrNotFound) {
		return diag.FromErr(err)
	}
	updatePolicyProtections.ApplyReputationControls = applyreputationcontrols

	applybotmancontrols, err := tools.GetBoolValue("apply_botman_controls", d)
	if err != nil && !errors.Is(err, tools.ErrNotFound) {
		return diag.FromErr(err)
	}
	updatePolicyProtections.ApplyBotmanControls = applybotmancontrols

	applyapiconstraints, err := tools.GetBoolValue("apply_api_constraints", d)
	if err != nil && !errors.Is(err, tools.ErrNotFound) {
		return diag.FromErr(err)
	}
	updatePolicyProtections.ApplyAPIConstraints = applyapiconstraints

	applyslowpostcontrols, err := tools.GetBoolValue("apply_slow_post_controls", d)
	if err != nil && !errors.Is(err, tools.ErrNotFound) {
		return diag.FromErr(err)
	}
	updatePolicyProtections.ApplySlowPostControls = applyslowpostcontrols

	_, erru := client.UpdatePolicyProtections(ctx, updatePolicyProtections)
	if erru != nil {
		logger.Errorf("calling 'updatePolicyProtections': %s", erru.Error())
		return diag.FromErr(erru)
	}

	return resourcePolicyProtectionsRead(ctx, d, m)
}
