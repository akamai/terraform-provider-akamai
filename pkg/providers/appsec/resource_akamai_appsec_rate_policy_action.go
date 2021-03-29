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
func resourceRatePolicyAction() *schema.Resource {
	return &schema.Resource{
		CreateContext: resourceRatePolicyActionUpdate,
		ReadContext:   resourceRatePolicyActionRead,
		UpdateContext: resourceRatePolicyActionUpdate,
		DeleteContext: resourceRatePolicyActionDelete,
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
			"rate_policy_id": {
				Type:     schema.TypeInt,
				Required: true,
			},
			"ipv4_action": {
				Type:         schema.TypeString,
				Required:     true,
				ValidateFunc: ValidateActions,
			},
			"ipv6_action": {
				Type:         schema.TypeString,
				Required:     true,
				ValidateFunc: ValidateActions},
		},
	}
}

func resourceRatePolicyActionRead(ctx context.Context, d *schema.ResourceData, m interface{}) diag.Diagnostics {
	meta := akamai.Meta(m)
	client := inst.Client(meta)
	logger := meta.Log("APPSEC", "resourceRatePolicyActionRead")

	getRatePolicyAction := appsec.GetRatePolicyActionRequest{}
	if d.Id() != "" && strings.Contains(d.Id(), ":") {
		s := strings.Split(d.Id(), ":")

		configid, errconv := strconv.Atoi(s[0])
		if errconv != nil {
			return diag.FromErr(errconv)
		}
		getRatePolicyAction.ConfigID = configid

		version, errconv := strconv.Atoi(s[1])
		if errconv != nil {
			return diag.FromErr(errconv)
		}
		getRatePolicyAction.Version = version

		if d.HasChange("version") {
			version, err := tools.GetIntValue("version", d)
			if err != nil && !errors.Is(err, tools.ErrNotFound) {
				return diag.FromErr(err)
			}
			getRatePolicyAction.Version = version
		}

		policyid := s[2]
		getRatePolicyAction.PolicyID = policyid

		ratepolicyid, errconv := strconv.Atoi(s[3])
		if errconv != nil {
			return diag.FromErr(errconv)
		}
		getRatePolicyAction.ID = ratepolicyid
	} else {
		configid, err := tools.GetIntValue("config_id", d)
		if err != nil && !errors.Is(err, tools.ErrNotFound) {
			return diag.FromErr(err)
		}
		getRatePolicyAction.ConfigID = configid

		version, err := tools.GetIntValue("version", d)
		if err != nil && !errors.Is(err, tools.ErrNotFound) {
			return diag.FromErr(err)
		}
		getRatePolicyAction.Version = version

		policyid, err := tools.GetStringValue("security_policy_id", d)
		if err != nil && !errors.Is(err, tools.ErrNotFound) {
			return diag.FromErr(err)
		}
		getRatePolicyAction.PolicyID = policyid

		ratepolicyid, err := tools.GetIntValue("rate_policy_id", d)
		if err != nil && !errors.Is(err, tools.ErrNotFound) {
			return diag.FromErr(err)
		}
		getRatePolicyAction.ID = ratepolicyid
	}
	ratepolicyaction, err := client.GetRatePolicyAction(ctx, getRatePolicyAction)
	if err != nil {
		logger.Errorf("calling 'getRatePolicyAction': %s", err.Error())
		return diag.FromErr(err)
	}
	logger.Warnf("calling 'GetRatePolicyAction': %s", ratepolicyaction)

	for _, configval := range ratepolicyaction.RatePolicyActions {
		if configval.ID == getRatePolicyAction.ID {
			d.SetId(strconv.Itoa(configval.ID))

			if err := d.Set("ipv4_action", configval.Ipv4Action); err != nil {
				return diag.FromErr(fmt.Errorf("%w: %s", tools.ErrValueSet, err.Error()))
			}

			if err := d.Set("ipv6_action", configval.Ipv6Action); err != nil {
				return diag.FromErr(fmt.Errorf("%w: %s", tools.ErrValueSet, err.Error()))
			}

			if err := d.Set("config_id", getRatePolicyAction.ConfigID); err != nil {
				return diag.FromErr(fmt.Errorf("%w: %s", tools.ErrValueSet, err.Error()))
			}

			if err := d.Set("version", getRatePolicyAction.Version); err != nil {
				return diag.FromErr(fmt.Errorf("%w: %s", tools.ErrValueSet, err.Error()))
			}

			if err := d.Set("security_policy_id", getRatePolicyAction.PolicyID); err != nil {
				return diag.FromErr(fmt.Errorf("%w: %s", tools.ErrValueSet, err.Error()))
			}

			d.SetId(fmt.Sprintf("%d:%d:%s:%d", getRatePolicyAction.ConfigID, getRatePolicyAction.Version, getRatePolicyAction.PolicyID, getRatePolicyAction.ID))

		}
	}

	return nil
}

func resourceRatePolicyActionDelete(ctx context.Context, d *schema.ResourceData, m interface{}) diag.Diagnostics {
	meta := akamai.Meta(m)
	client := inst.Client(meta)
	logger := meta.Log("APPSEC", "resourceRatePolicyActionRemove")

	updateRatePolicyAction := appsec.UpdateRatePolicyActionRequest{}
	if d.Id() != "" && strings.Contains(d.Id(), ":") {
		s := strings.Split(d.Id(), ":")

		configid, errconv := strconv.Atoi(s[0])
		if errconv != nil {
			return diag.FromErr(errconv)
		}
		updateRatePolicyAction.ConfigID = configid

		version, errconv := strconv.Atoi(s[1])
		if errconv != nil {
			return diag.FromErr(errconv)
		}
		updateRatePolicyAction.Version = version

		policyid := s[2]
		updateRatePolicyAction.PolicyID = policyid

		ratepolicyid, errconv := strconv.Atoi(s[3])
		if errconv != nil {
			return diag.FromErr(errconv)
		}
		updateRatePolicyAction.RatePolicyID = ratepolicyid
	} else {
		configid, err := tools.GetIntValue("config_id", d)
		if err != nil && !errors.Is(err, tools.ErrNotFound) {
			return diag.FromErr(err)
		}
		updateRatePolicyAction.ConfigID = configid

		version, err := tools.GetIntValue("version", d)
		if err != nil && !errors.Is(err, tools.ErrNotFound) {
			return diag.FromErr(err)
		}
		updateRatePolicyAction.Version = version

		policyid, err := tools.GetStringValue("security_policy_id", d)
		if err != nil && !errors.Is(err, tools.ErrNotFound) {
			return diag.FromErr(err)
		}
		updateRatePolicyAction.PolicyID = policyid

		ratepolicyid, err := tools.GetIntValue("rate_policy_id", d)
		if err != nil && !errors.Is(err, tools.ErrNotFound) {
			return diag.FromErr(err)
		}
		updateRatePolicyAction.RatePolicyID = ratepolicyid
	}
	updateRatePolicyAction.Ipv4Action = "none"
	updateRatePolicyAction.Ipv6Action = "none"

	resp, erru := client.UpdateRatePolicyAction(ctx, updateRatePolicyAction)
	if erru != nil {
		logger.Errorf("calling 'removeRatePolicyAction': %s", erru.Error())
		return diag.FromErr(erru)
	}
	logger.Warnf("calling 'RemoveRatePolicyAction': %s", resp)
	d.SetId("")

	return nil
}

func resourceRatePolicyActionUpdate(ctx context.Context, d *schema.ResourceData, m interface{}) diag.Diagnostics {
	meta := akamai.Meta(m)
	client := inst.Client(meta)
	logger := meta.Log("APPSEC", "resourceRatePolicyActionUpdate")

	updateRatePolicyAction := appsec.UpdateRatePolicyActionRequest{}
	if d.Id() != "" && strings.Contains(d.Id(), ":") {
		s := strings.Split(d.Id(), ":")

		configid, errconv := strconv.Atoi(s[0])
		if errconv != nil {
			return diag.FromErr(errconv)
		}
		updateRatePolicyAction.ConfigID = configid

		version, errconv := strconv.Atoi(s[1])
		if errconv != nil {
			return diag.FromErr(errconv)
		}
		updateRatePolicyAction.Version = version

		if d.HasChange("version") {
			version, err := tools.GetIntValue("version", d)
			if err != nil && !errors.Is(err, tools.ErrNotFound) {
				return diag.FromErr(err)
			}
			updateRatePolicyAction.Version = version
		}

		policyid := s[2]
		updateRatePolicyAction.PolicyID = policyid

		ratepolicyid, errconv := strconv.Atoi(s[3])
		if errconv != nil {
			return diag.FromErr(errconv)
		}
		updateRatePolicyAction.RatePolicyID = ratepolicyid
	} else {
		configid, err := tools.GetIntValue("config_id", d)
		if err != nil && !errors.Is(err, tools.ErrNotFound) {
			return diag.FromErr(err)
		}
		updateRatePolicyAction.ConfigID = configid

		version, err := tools.GetIntValue("version", d)
		if err != nil && !errors.Is(err, tools.ErrNotFound) {
			return diag.FromErr(err)
		}
		updateRatePolicyAction.Version = version

		policyid, err := tools.GetStringValue("security_policy_id", d)
		if err != nil && !errors.Is(err, tools.ErrNotFound) {
			return diag.FromErr(err)
		}
		updateRatePolicyAction.PolicyID = policyid

		ratepolicyid, err := tools.GetIntValue("rate_policy_id", d)
		if err != nil && !errors.Is(err, tools.ErrNotFound) {
			return diag.FromErr(err)
		}
		updateRatePolicyAction.RatePolicyID = ratepolicyid
	}
	ipv4action, err := tools.GetStringValue("ipv4_action", d)
	if err != nil && !errors.Is(err, tools.ErrNotFound) {
		return diag.FromErr(err)
	}
	updateRatePolicyAction.Ipv4Action = ipv4action

	ipv6action, err := tools.GetStringValue("ipv6_action", d)
	if err != nil && !errors.Is(err, tools.ErrNotFound) {
		return diag.FromErr(err)
	}
	updateRatePolicyAction.Ipv6Action = ipv6action
	logger.Warnf("calling 'updateRatePolicyAction REQ': %s", updateRatePolicyAction)
	resp, erru := client.UpdateRatePolicyAction(ctx, updateRatePolicyAction)
	if erru != nil {
		logger.Errorf("calling 'updateRatePolicyAction': %s", erru.Error())
		return diag.FromErr(erru)
	}
	logger.Warnf("calling 'updateRatePolicyAction': %s", resp)

	d.SetId(strconv.Itoa(resp.ID))

	if err := d.Set("ipv4_action", resp.Ipv4Action); err != nil {
		return diag.FromErr(fmt.Errorf("%w: %s", tools.ErrValueSet, err.Error()))
	}

	if err := d.Set("ipv6_action", resp.Ipv6Action); err != nil {
		return diag.FromErr(fmt.Errorf("%w: %s", tools.ErrValueSet, err.Error()))
	}

	d.SetId(strconv.Itoa(resp.ID))

	return resourceRatePolicyActionRead(ctx, d, m)
}
