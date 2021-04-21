package appsec

import (
	"context"
	"encoding/json"
	"errors"
	"fmt"
	"strconv"
	"strings"

	"github.com/akamai/AkamaiOPEN-edgegrid-golang/v2/pkg/appsec"
	"github.com/akamai/terraform-provider-akamai/v2/pkg/akamai"
	"github.com/akamai/terraform-provider-akamai/v2/pkg/tools"
	"github.com/hashicorp/terraform-plugin-sdk/v2/diag"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/schema"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/validation"
)

// appsec v1
//
// https://developer.akamai.com/api/cloud_security/application_security/v1.html
func resourceAttackGroupConditionException() *schema.Resource {
	return &schema.Resource{
		CreateContext: resourceAttackGroupConditionExceptionUpdate,
		ReadContext:   resourceAttackGroupConditionExceptionRead,
		UpdateContext: resourceAttackGroupConditionExceptionUpdate,
		DeleteContext: resourceAttackGroupConditionExceptionDelete,
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
			"attack_group": {
				Type:     schema.TypeString,
				Required: true,
			},
			"condition_exception": {
				Type:             schema.TypeString,
				Required:         true,
				ValidateFunc:     validation.StringIsJSON,
				DiffSuppressFunc: suppressEquivalentJsonDiffsGeneric,
			},
		},
	}
}

func resourceAttackGroupConditionExceptionRead(ctx context.Context, d *schema.ResourceData, m interface{}) diag.Diagnostics {
	meta := akamai.Meta(m)
	client := inst.Client(meta)
	logger := meta.Log("APPSEC", "resourceAttackGroupConditionExceptionRead")

	getAttackGroupConditionException := appsec.GetAttackGroupConditionExceptionRequest{}
	if d.Id() != "" && strings.Contains(d.Id(), ":") {
		s := strings.Split(d.Id(), ":")

		configid, errconv := strconv.Atoi(s[0])
		if errconv != nil {
			return diag.FromErr(errconv)
		}
		getAttackGroupConditionException.ConfigID = configid

		version, errconv := strconv.Atoi(s[1])
		if errconv != nil {
			return diag.FromErr(errconv)
		}
		getAttackGroupConditionException.Version = version

		if d.HasChange("version") {
			version, err := tools.GetIntValue("version", d)
			if err != nil && !errors.Is(err, tools.ErrNotFound) {
				return diag.FromErr(err)
			}
			getAttackGroupConditionException.Version = version
		}

		policyid := s[2]

		getAttackGroupConditionException.PolicyID = policyid

		attackgroup := s[3]

		getAttackGroupConditionException.Group = attackgroup

	} else {
		configid, err := tools.GetIntValue("config_id", d)
		if err != nil && !errors.Is(err, tools.ErrNotFound) {
			return diag.FromErr(err)
		}
		getAttackGroupConditionException.ConfigID = configid

		version, err := tools.GetIntValue("version", d)
		if err != nil && !errors.Is(err, tools.ErrNotFound) {
			return diag.FromErr(err)
		}
		getAttackGroupConditionException.Version = version

		policyid, err := tools.GetStringValue("security_policy_id", d)
		if err != nil && !errors.Is(err, tools.ErrNotFound) {
			return diag.FromErr(err)
		}
		getAttackGroupConditionException.PolicyID = policyid

		attackgroup, err := tools.GetStringValue("attack_group", d)
		if err != nil && !errors.Is(err, tools.ErrNotFound) {
			return diag.FromErr(err)
		}
		getAttackGroupConditionException.Group = attackgroup
	}
	attackgroupconditionexception, err := client.GetAttackGroupConditionException(ctx, getAttackGroupConditionException)
	if err != nil {
		logger.Errorf("calling 'getAttackGroupConditionException': %s", err.Error())
		return diag.FromErr(err)
	}

	if err := d.Set("config_id", getAttackGroupConditionException.ConfigID); err != nil {
		return diag.FromErr(fmt.Errorf("%w: %s", tools.ErrValueSet, err.Error()))
	}

	if err := d.Set("version", getAttackGroupConditionException.Version); err != nil {
		return diag.FromErr(fmt.Errorf("%w: %s", tools.ErrValueSet, err.Error()))
	}

	if err := d.Set("security_policy_id", getAttackGroupConditionException.PolicyID); err != nil {
		return diag.FromErr(fmt.Errorf("%w: %s", tools.ErrValueSet, err.Error()))
	}

	if err := d.Set("attack_group", getAttackGroupConditionException.Group); err != nil {
		return diag.FromErr(fmt.Errorf("%w: %s", tools.ErrValueSet, err.Error()))
	}

	jsonBody, err := json.Marshal(attackgroupconditionexception)
	if err != nil {
		return diag.FromErr(err)
	}
	if err := d.Set("condition_exception", string(jsonBody)); err != nil {
		return diag.FromErr(fmt.Errorf("%w: %s", tools.ErrValueSet, err.Error()))
	}

	d.SetId(fmt.Sprintf("%d:%d:%s:%s", getAttackGroupConditionException.ConfigID, getAttackGroupConditionException.Version, getAttackGroupConditionException.PolicyID, getAttackGroupConditionException.Group))

	return nil
}

func resourceAttackGroupConditionExceptionDelete(ctx context.Context, d *schema.ResourceData, m interface{}) diag.Diagnostics {

	meta := akamai.Meta(m)
	client := inst.Client(meta)
	logger := meta.Log("APPSEC", "resourceAttackGroupConditionExceptionRemove")

	removeAttackGroupConditionException := appsec.RemoveAttackGroupConditionExceptionRequest{}
	if d.Id() != "" && strings.Contains(d.Id(), ":") {
		s := strings.Split(d.Id(), ":")

		configid, errconv := strconv.Atoi(s[0])
		if errconv != nil {
			return diag.FromErr(errconv)
		}
		removeAttackGroupConditionException.ConfigID = configid

		version, errconv := strconv.Atoi(s[1])
		if errconv != nil {
			return diag.FromErr(errconv)
		}
		removeAttackGroupConditionException.Version = version

		policyid := s[2]

		removeAttackGroupConditionException.PolicyID = policyid

		attackgroup := s[3]

		removeAttackGroupConditionException.Group = attackgroup

	} else {
		configid, err := tools.GetIntValue("config_id", d)
		if err != nil && !errors.Is(err, tools.ErrNotFound) {
			return diag.FromErr(err)
		}
		removeAttackGroupConditionException.ConfigID = configid

		version, err := tools.GetIntValue("version", d)
		if err != nil && !errors.Is(err, tools.ErrNotFound) {
			return diag.FromErr(err)
		}
		removeAttackGroupConditionException.Version = version

		policyid, err := tools.GetStringValue("security_policy_id", d)
		if err != nil && !errors.Is(err, tools.ErrNotFound) {
			return diag.FromErr(err)
		}
		removeAttackGroupConditionException.PolicyID = policyid

		attackgroup, err := tools.GetStringValue("attack_group", d)
		if err != nil && !errors.Is(err, tools.ErrNotFound) {
			return diag.FromErr(err)
		}
		removeAttackGroupConditionException.Group = attackgroup
	}
	_, errd := client.RemoveAttackGroupConditionException(ctx, removeAttackGroupConditionException)
	if errd != nil {
		logger.Errorf("calling 'RemoveAttackGroupConditionException': %s", errd.Error())
		return diag.FromErr(errd)
	}

	d.SetId("")
	return nil
}

func resourceAttackGroupConditionExceptionUpdate(ctx context.Context, d *schema.ResourceData, m interface{}) diag.Diagnostics {
	meta := akamai.Meta(m)
	client := inst.Client(meta)
	logger := meta.Log("APPSEC", "resourceAttackGroupConditionExceptionUpdate")

	updateAttackGroupConditionException := appsec.UpdateAttackGroupConditionExceptionRequest{}

	jsonpostpayload := d.Get("condition_exception")
	jsonPayloadRaw := []byte(jsonpostpayload.(string))
	rawJSON := (json.RawMessage)(jsonPayloadRaw)

	updateAttackGroupConditionException.JsonPayloadRaw = rawJSON
	if d.Id() != "" && strings.Contains(d.Id(), ":") {
		s := strings.Split(d.Id(), ":")

		configid, errconv := strconv.Atoi(s[0])
		if errconv != nil {
			return diag.FromErr(errconv)
		}
		updateAttackGroupConditionException.ConfigID = configid

		version, errconv := strconv.Atoi(s[1])
		if errconv != nil {
			return diag.FromErr(errconv)
		}
		updateAttackGroupConditionException.Version = version

		if d.HasChange("version") {
			version, err := tools.GetIntValue("version", d)
			if err != nil && !errors.Is(err, tools.ErrNotFound) {
				return diag.FromErr(err)
			}
			updateAttackGroupConditionException.Version = version
		}

		policyid := s[2]

		updateAttackGroupConditionException.PolicyID = policyid

		attackgroup := s[3]

		updateAttackGroupConditionException.Group = attackgroup

	} else {
		configid, err := tools.GetIntValue("config_id", d)
		if err != nil && !errors.Is(err, tools.ErrNotFound) {
			return diag.FromErr(err)
		}
		updateAttackGroupConditionException.ConfigID = configid

		version, err := tools.GetIntValue("version", d)
		if err != nil && !errors.Is(err, tools.ErrNotFound) {
			return diag.FromErr(err)
		}
		updateAttackGroupConditionException.Version = version

		policyid, err := tools.GetStringValue("security_policy_id", d)
		if err != nil && !errors.Is(err, tools.ErrNotFound) {
			return diag.FromErr(err)
		}
		updateAttackGroupConditionException.PolicyID = policyid

		attackgroup, err := tools.GetStringValue("attack_group", d)
		if err != nil && !errors.Is(err, tools.ErrNotFound) {
			return diag.FromErr(err)
		}
		updateAttackGroupConditionException.Group = attackgroup
	}
	_, erru := client.UpdateAttackGroupConditionException(ctx, updateAttackGroupConditionException)
	if erru != nil {
		logger.Errorf("calling 'updateAttackGroupConditionException': %s", erru.Error())
		return diag.FromErr(erru)
	}
	return resourceAttackGroupConditionExceptionRead(ctx, d, m)
}
