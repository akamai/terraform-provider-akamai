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
func resourceApiRequestConstraints() *schema.Resource {
	return &schema.Resource{
		CreateContext: resourceApiRequestConstraintsUpdate,
		ReadContext:   resourceApiRequestConstraintsRead,
		UpdateContext: resourceApiRequestConstraintsUpdate,
		DeleteContext: resourceApiRequestConstraintsDelete,
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
			"action": {
				Type:         schema.TypeString,
				Required:     true,
				ValidateFunc: ValidateActions,
			},
			"api_endpoint_id": {
				Type:     schema.TypeInt,
				Optional: true,
			},
		},
	}
}

func resourceApiRequestConstraintsRead(ctx context.Context, d *schema.ResourceData, m interface{}) diag.Diagnostics {
	meta := akamai.Meta(m)
	client := inst.Client(meta)
	logger := meta.Log("APPSEC", "resourceApiRequestConstraintsRead")

	getApiRequestConstraints := appsec.GetApiRequestConstraintsRequest{}
	if d.Id() != "" && strings.Contains(d.Id(), ":") {
		s := strings.Split(d.Id(), ":")

		configid, errconv := strconv.Atoi(s[0])
		if errconv != nil {
			return diag.FromErr(errconv)
		}
		getApiRequestConstraints.ConfigID = configid

		version, errconv := strconv.Atoi(s[1])
		if errconv != nil {
			return diag.FromErr(errconv)
		}
		getApiRequestConstraints.Version = version

		policyid := s[2]
		getApiRequestConstraints.PolicyID = policyid

		if len(s) >= 4 {
			apiID, errconv := strconv.Atoi(s[3])
			if errconv != nil {
				return diag.FromErr(errconv)
			}
			getApiRequestConstraints.ApiID = apiID
		}
	} else {
		configid, err := tools.GetIntValue("config_id", d)
		if err != nil && !errors.Is(err, tools.ErrNotFound) {
			return diag.FromErr(err)
		}
		getApiRequestConstraints.ConfigID = configid

		version, err := tools.GetIntValue("version", d)
		if err != nil && !errors.Is(err, tools.ErrNotFound) {
			return diag.FromErr(err)
		}
		getApiRequestConstraints.Version = version

		policyid, err := tools.GetStringValue("security_policy_id", d)
		if err != nil && !errors.Is(err, tools.ErrNotFound) {
			return diag.FromErr(err)
		}
		getApiRequestConstraints.PolicyID = policyid

		ApiID, err := tools.GetIntValue("api_id", d)
		if err != nil && !errors.Is(err, tools.ErrNotFound) {
			return diag.FromErr(err)
		}
		getApiRequestConstraints.ApiID = ApiID
	}
	response, err := client.GetApiRequestConstraints(ctx, getApiRequestConstraints)
	if err != nil {
		logger.Errorf("calling 'getApiRequestConstraints': %s", err.Error())
		return diag.FromErr(err)
	}

	if err := d.Set("config_id", getApiRequestConstraints.ConfigID); err != nil {
		return diag.FromErr(fmt.Errorf("%w: %s", tools.ErrValueSet, err.Error()))
	}

	if err := d.Set("version", getApiRequestConstraints.Version); err != nil {
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

	d.SetId(fmt.Sprintf("%d:%d:%s", getApiRequestConstraints.ConfigID, getApiRequestConstraints.Version, getApiRequestConstraints.PolicyID))

	return nil
}

func resourceApiRequestConstraintsDelete(ctx context.Context, d *schema.ResourceData, m interface{}) diag.Diagnostics {

	meta := akamai.Meta(m)
	client := inst.Client(meta)
	logger := meta.Log("APPSEC", "resourceApiRequestConstraintsRemove")

	getPolicyProtections := appsec.GetPolicyProtectionsRequest{}
	removeApiRequestConstraints := appsec.RemoveApiRequestConstraintsRequest{}
	removePolicyProtections := appsec.RemovePolicyProtectionsRequest{}
	if d.Id() != "" && strings.Contains(d.Id(), ":") {
		s := strings.Split(d.Id(), ":")

		configid, errconv := strconv.Atoi(s[0])
		if errconv != nil {
			return diag.FromErr(errconv)
		}
		getPolicyProtections.ConfigID = configid
		removeApiRequestConstraints.ConfigID = configid
		removePolicyProtections.ConfigID = configid

		version, errconv := strconv.Atoi(s[1])
		if errconv != nil {
			return diag.FromErr(errconv)
		}
		getPolicyProtections.Version = version
		removeApiRequestConstraints.Version = version
		removePolicyProtections.Version = version

		policyid := s[2]

		getPolicyProtections.PolicyID = policyid
		removeApiRequestConstraints.PolicyID = policyid
		removePolicyProtections.PolicyID = policyid

		if len(s) >= 4 {
			apiID, errconv := strconv.Atoi(s[3])
			if errconv != nil {
				return diag.FromErr(errconv)
			}

			removeApiRequestConstraints.ApiID = apiID

		}
	} else {
		configid, err := tools.GetIntValue("config_id", d)
		if err != nil && !errors.Is(err, tools.ErrNotFound) {
			return diag.FromErr(err)
		}

		version, err := tools.GetIntValue("version", d)
		if err != nil && !errors.Is(err, tools.ErrNotFound) {
			return diag.FromErr(err)
		}

		policyid, err := tools.GetStringValue("security_policy_id", d)
		if err != nil && !errors.Is(err, tools.ErrNotFound) {
			return diag.FromErr(err)
		}

		getPolicyProtections.ConfigID = configid
		removeApiRequestConstraints.ConfigID = configid
		removePolicyProtections.ConfigID = configid

		getPolicyProtections.Version = version
		removeApiRequestConstraints.Version = version
		removePolicyProtections.Version = version

		getPolicyProtections.PolicyID = policyid
		removeApiRequestConstraints.PolicyID = policyid
		removePolicyProtections.PolicyID = policyid
	}
	policyprotections, err := client.GetPolicyProtections(ctx, getPolicyProtections)
	if err != nil {
		logger.Errorf("calling 'getPolicyProtections': %s", err.Error())
		return diag.FromErr(err)
	}

	apiEndpointID, err := tools.GetIntValue("api_endpoint_id", d)
	if err != nil && !errors.Is(err, tools.ErrNotFound) {
		return diag.FromErr(err)
	}
	removeApiRequestConstraints.ApiID = apiEndpointID

	if removeApiRequestConstraints.ApiID == 0 {
		if policyprotections.ApplyAPIConstraints == true {
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

func resourceApiRequestConstraintsUpdate(ctx context.Context, d *schema.ResourceData, m interface{}) diag.Diagnostics {
	meta := akamai.Meta(m)
	client := inst.Client(meta)
	logger := meta.Log("APPSEC", "resourceApiRequestConstraintsUpdate")

	updateApiRequestConstraints := appsec.UpdateApiRequestConstraintsRequest{}
	if d.Id() != "" && strings.Contains(d.Id(), ":") {
		s := strings.Split(d.Id(), ":")

		configid, errconv := strconv.Atoi(s[0])
		if errconv != nil {
			return diag.FromErr(errconv)
		}
		updateApiRequestConstraints.ConfigID = configid

		version, errconv := strconv.Atoi(s[1])
		if errconv != nil {
			return diag.FromErr(errconv)
		}
		updateApiRequestConstraints.Version = version

		policyid := s[2]

		updateApiRequestConstraints.PolicyID = policyid

		if len(s) >= 4 {
			apiID, errconv := strconv.Atoi(s[3])
			if errconv != nil {
				return diag.FromErr(errconv)
			}

			updateApiRequestConstraints.ApiID = apiID
		}
	} else {
		configid, err := tools.GetIntValue("config_id", d)
		if err != nil && !errors.Is(err, tools.ErrNotFound) {
			return diag.FromErr(err)
		}
		updateApiRequestConstraints.ConfigID = configid

		version, err := tools.GetIntValue("version", d)
		if err != nil && !errors.Is(err, tools.ErrNotFound) {
			return diag.FromErr(err)
		}
		updateApiRequestConstraints.Version = version

		policyid, err := tools.GetStringValue("security_policy_id", d)
		if err != nil && !errors.Is(err, tools.ErrNotFound) {
			return diag.FromErr(err)
		}
		updateApiRequestConstraints.PolicyID = policyid

		apiEndpointID, err := tools.GetIntValue("api_endpoint_id", d)
		if err != nil && !errors.Is(err, tools.ErrNotFound) {
			return diag.FromErr(err)
		}
		updateApiRequestConstraints.ApiID = apiEndpointID
	}
	action, err := tools.GetStringValue("action", d)
	if err != nil && !errors.Is(err, tools.ErrNotFound) {
		return diag.FromErr(err)
	}
	updateApiRequestConstraints.Action = action

	_, erru := client.UpdateApiRequestConstraints(ctx, updateApiRequestConstraints)
	if erru != nil {
		logger.Errorf("calling 'updateApiRequestConstraints': %s", erru.Error())
		return diag.FromErr(erru)
	}

	return resourceApiRequestConstraintsRead(ctx, d, m)
}
