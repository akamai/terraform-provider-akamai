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
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/customdiff"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/schema"
)

// appsec v1
//
// https://developer.akamai.com/api/cloud_security/application_security/v1.html
func resourceApiRequestConstraints() *schema.Resource {
	return &schema.Resource{
		CreateContext: resourceApiRequestConstraintsCreate,
		ReadContext:   resourceApiRequestConstraintsRead,
		UpdateContext: resourceApiRequestConstraintsUpdate,
		DeleteContext: resourceApiRequestConstraintsDelete,
		CustomizeDiff: customdiff.All(
			VerifyIdUnchanged,
		),
		Importer: &schema.ResourceImporter{
			State: schema.ImportStatePassthrough,
		},
		Schema: map[string]*schema.Schema{
			"config_id": {
				Type:     schema.TypeInt,
				Required: true,
			},
			"security_policy_id": {
				Type:     schema.TypeString,
				Required: true,
			},
			"api_endpoint_id": {
				Type:     schema.TypeInt,
				Optional: true,
			},
			"action": {
				Type:         schema.TypeString,
				Required:     true,
				ValidateFunc: ValidateActions,
			},
		},
	}
}

func resourceApiRequestConstraintsCreate(ctx context.Context, d *schema.ResourceData, m interface{}) diag.Diagnostics {
	meta := akamai.Meta(m)
	client := inst.Client(meta)
	logger := meta.Log("APPSEC", "resourceApiRequestConstraintsCreate")
	logger.Debugf("!!! in resourceApiRequestConstraintsCreate")

	configid, err := tools.GetIntValue("config_id", d)
	if err != nil && !errors.Is(err, tools.ErrNotFound) {
		return diag.FromErr(err)
	}
	version := getModifiableConfigVersion(ctx, configid, "apirequestconstraints", m)
	policyid, err := tools.GetStringValue("security_policy_id", d)
	if err != nil && !errors.Is(err, tools.ErrNotFound) {
		return diag.FromErr(err)
	}
	apiEndpointID, err := tools.GetIntValue("api_endpoint_id", d)
	if err != nil && !errors.Is(err, tools.ErrNotFound) {
		return diag.FromErr(err)
	}
	action, err := tools.GetStringValue("action", d)
	if err != nil && !errors.Is(err, tools.ErrNotFound) {
		return diag.FromErr(err)
	}

	createApiRequestConstraints := appsec.UpdateApiRequestConstraintsRequest{}
	createApiRequestConstraints.ConfigID = configid
	createApiRequestConstraints.Version = version
	createApiRequestConstraints.PolicyID = policyid
	createApiRequestConstraints.ApiID = apiEndpointID
	createApiRequestConstraints.Action = action

	_, erru := client.UpdateApiRequestConstraints(ctx, createApiRequestConstraints)
	if erru != nil {
		logger.Errorf("calling 'createApiRequestConstraints': %s", erru.Error())
		return diag.FromErr(erru)
	}

	if apiEndpointID != 0 {
		d.SetId(fmt.Sprintf("%d:%s:%d", createApiRequestConstraints.ConfigID, createApiRequestConstraints.PolicyID, createApiRequestConstraints.ApiID))
	} else {
		d.SetId(fmt.Sprintf("%d:%s", createApiRequestConstraints.ConfigID, createApiRequestConstraints.PolicyID))
	}

	return resourceApiRequestConstraintsRead(ctx, d, m)
}

func resourceApiRequestConstraintsRead(ctx context.Context, d *schema.ResourceData, m interface{}) diag.Diagnostics {
	meta := akamai.Meta(m)
	client := inst.Client(meta)
	logger := meta.Log("APPSEC", "resourceApiRequestConstraintsRead")
	logger.Debugf("!!! in resourceCustomRuleActionRead")

	s := strings.Split(d.Id(), ":")

	configid, errconv := strconv.Atoi(s[0])
	if errconv != nil {
		return diag.FromErr(errconv)
	}
	version := getLatestConfigVersion(ctx, configid, m)
	policyid := s[1]

	apiID := 0
	if len(s) > 2 {
		apiID, errconv = strconv.Atoi(s[2])
		if errconv != nil {
			return diag.FromErr(errconv)
		}
	}

	getApiRequestConstraints := appsec.GetApiRequestConstraintsRequest{}
	getApiRequestConstraints.ConfigID = configid
	getApiRequestConstraints.Version = version
	getApiRequestConstraints.PolicyID = policyid
	getApiRequestConstraints.ApiID = apiID

	response, err := client.GetApiRequestConstraints(ctx, getApiRequestConstraints)
	if err != nil {
		logger.Errorf("calling 'getApiRequestConstraints': %s", err.Error())
		return diag.FromErr(err)
	}

	if err := d.Set("config_id", getApiRequestConstraints.ConfigID); err != nil {
		return diag.FromErr(fmt.Errorf("%w: %s", tools.ErrValueSet, err.Error()))
	}

	if err := d.Set("security_policy_id", getApiRequestConstraints.PolicyID); err != nil {
		return diag.FromErr(fmt.Errorf("%w: %s", tools.ErrValueSet, err.Error()))
	}

	if err := d.Set("api_endpoint_id", getApiRequestConstraints.ApiID); err != nil {
		return diag.FromErr(fmt.Errorf("%w: %s", tools.ErrValueSet, err.Error()))
	}

	if getApiRequestConstraints.ApiID != 0 {
		if len(response.APIEndpoints) > 0 {
			for _, val := range response.APIEndpoints {
				if val.ID == getApiRequestConstraints.ApiID {
					if err := d.Set("action", val.Action); err != nil {
						return diag.FromErr(fmt.Errorf("%w: %s", tools.ErrValueSet, err.Error()))
					}
				}
			}
		}
	}
	return nil
}

func resourceApiRequestConstraintsUpdate(ctx context.Context, d *schema.ResourceData, m interface{}) diag.Diagnostics {
	meta := akamai.Meta(m)
	client := inst.Client(meta)
	logger := meta.Log("APPSEC", "resourceApiRequestConstraintsUpdate")
	logger.Debugf("!!! in resourceCustomRuleActionUpdate")

	s := strings.Split(d.Id(), ":")

	configid, errconv := strconv.Atoi(s[0])
	if errconv != nil {
		return diag.FromErr(errconv)
	}
	version := getModifiableConfigVersion(ctx, configid, "apirequestconstraints", m)
	policyid := s[1]

	apiID := 0
	if len(s) > 2 {
		apiID, errconv = strconv.Atoi(s[2])
		if errconv != nil {
			return diag.FromErr(errconv)
		}
	}
	action, err := tools.GetStringValue("action", d)
	if err != nil && !errors.Is(err, tools.ErrNotFound) {
		return diag.FromErr(err)
	}

	updateApiRequestConstraints := appsec.UpdateApiRequestConstraintsRequest{}
	updateApiRequestConstraints.ConfigID = configid
	updateApiRequestConstraints.Version = version
	updateApiRequestConstraints.PolicyID = policyid
	updateApiRequestConstraints.ApiID = apiID
	updateApiRequestConstraints.Action = action

	_, erru := client.UpdateApiRequestConstraints(ctx, updateApiRequestConstraints)
	if erru != nil {
		logger.Errorf("calling 'updateApiRequestConstraints': %s", erru.Error())
		return diag.FromErr(erru)
	}

	return resourceApiRequestConstraintsRead(ctx, d, m)
}

func resourceApiRequestConstraintsDelete(ctx context.Context, d *schema.ResourceData, m interface{}) diag.Diagnostics {
	meta := akamai.Meta(m)
	client := inst.Client(meta)
	logger := meta.Log("APPSEC", "resourceApiRequestConstraintsDelete")
	logger.Debugf("!!! in resourceApiRequestConstraintsDelete")

	s := strings.Split(d.Id(), ":")

	configid, errconv := strconv.Atoi(s[0])
	if errconv != nil {
		return diag.FromErr(errconv)
	}
	version := getModifiableConfigVersion(ctx, configid, "apirequestconstraints", m)
	policyid := s[1]

	apiID := 0
	if len(s) > 2 {
		apiID, errconv = strconv.Atoi(s[2])
		if errconv != nil {
			return diag.FromErr(errconv)
		}
	}

	removeApiRequestConstraints := appsec.RemoveApiRequestConstraintsRequest{}
	removeApiRequestConstraints.ConfigID = configid
	removeApiRequestConstraints.Version = version
	removeApiRequestConstraints.PolicyID = policyid
	removeApiRequestConstraints.ApiID = apiID

	if removeApiRequestConstraints.ApiID == 0 {

		getPolicyProtections := appsec.GetPolicyProtectionsRequest{}
		getPolicyProtections.ConfigID = configid
		getPolicyProtections.Version = version
		getPolicyProtections.PolicyID = policyid

		policyprotections, err := client.GetPolicyProtections(ctx, getPolicyProtections)
		if err != nil {
			logger.Errorf("calling 'getPolicyProtections': %s", err.Error())
			return diag.FromErr(err)
		}
		if policyprotections.ApplyAPIConstraints {
			removePolicyProtections := appsec.RemovePolicyProtectionsRequest{}
			removePolicyProtections.ConfigID = configid
			removePolicyProtections.Version = version
			removePolicyProtections.PolicyID = policyid

			removePolicyProtections.ApplyAPIConstraints = false
			removePolicyProtections.ApplyApplicationLayerControls = policyprotections.ApplyApplicationLayerControls
			removePolicyProtections.ApplyBotmanControls = policyprotections.ApplyBotmanControls
			removePolicyProtections.ApplyNetworkLayerControls = policyprotections.ApplyNetworkLayerControls
			removePolicyProtections.ApplyRateControls = policyprotections.ApplyRateControls
			removePolicyProtections.ApplyReputationControls = policyprotections.ApplyReputationControls
			removePolicyProtections.ApplySlowPostControls = policyprotections.ApplySlowPostControls

			_, errd := client.RemovePolicyProtections(ctx, removePolicyProtections)
			if errd != nil {
				logger.Errorf("calling 'removePolicyProtections': %s", errd.Error())
				return diag.FromErr(errd)
			}
		}
	} else {
		removeApiRequestConstraints.Action = "none"
		_, erru := client.RemoveApiRequestConstraints(ctx, removeApiRequestConstraints)
		if erru != nil {
			logger.Errorf("calling 'removeApiRequestConstraints': %s", erru.Error())
			return diag.FromErr(erru)
		}
	}

	d.SetId("")
	return nil
}
