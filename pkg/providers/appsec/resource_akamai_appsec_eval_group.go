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
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/customdiff"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/schema"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/validation"
)

// appsec v1
//
// https://developer.akamai.com/api/cloud_security/application_security/v1.html
func resourceEvalGroup() *schema.Resource {
	return &schema.Resource{
		CreateContext: resourceEvalGroupCreate,
		ReadContext:   resourceEvalGroupRead,
		UpdateContext: resourceEvalGroupUpdate,
		DeleteContext: resourceEvalGroupDelete,
		CustomizeDiff: customdiff.All(
			VerifyIDUnchanged,
		),
		Importer: &schema.ResourceImporter{
			StateContext: schema.ImportStatePassthroughContext,
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
			"attack_group": {
				Type:     schema.TypeString,
				Required: true,
			},
			"attack_group_action": {
				Type:             schema.TypeString,
				Required:         true,
				ValidateDiagFunc: ValidateActions,
			},
			"condition_exception": {
				Type:             schema.TypeString,
				Optional:         true,
				Default:          "",
				ValidateFunc:     validation.StringIsJSON,
				DiffSuppressFunc: suppressEquivalentJSONDiffsGeneric,
			},
		},
	}
}

func resourceEvalGroupCreate(ctx context.Context, d *schema.ResourceData, m interface{}) diag.Diagnostics {
	meta := akamai.Meta(m)
	client := inst.Client(meta)
	logger := meta.Log("APPSEC", "resourceEvalGroupCreate")
	logger.Debugf("in resourceEvalGroupCreate")

	configid, err := tools.GetIntValue("config_id", d)
	if err != nil {
		return diag.FromErr(err)
	}
	version := getModifiableConfigVersion(ctx, configid, "atackGroup", m)
	policyid, err := tools.GetStringValue("security_policy_id", d)
	if err != nil {
		return diag.FromErr(err)
	}
	attackgroup, err := tools.GetStringValue("attack_group", d)
	if err != nil {
		return diag.FromErr(err)
	}
	action, err := tools.GetStringValue("attack_group_action", d)
	if err != nil {
		return diag.FromErr(err)
	}
	conditionexception, err := tools.GetStringValue("condition_exception", d)
	if err != nil && !errors.Is(err, tools.ErrNotFound) {
		return diag.FromErr(err)
	}

	rawJSON := json.RawMessage(conditionexception)

	if err := validateActionAndConditionException(action, conditionexception); err != nil {
		return diag.FromErr(err)
	}

	createAttackGroup := appsec.UpdateAttackGroupRequest{
		ConfigID: configid,
		Version: version,
		PolicyID: policyid,
		Group: attackgroup,
		Action: action,
		JsonPayloadRaw: rawJSON,
	}

	updateEvalGroupResponse, err := client.UpdateEvalGroup(ctx, createAttackGroup)
	if err != nil {
		logger.Errorf("calling 'createEvalGroup': %s", err.Error())
		return diag.FromErr(err)
	}
	logger.Debugf("updateEvalGroupResponse: %v", updateEvalGroupResponse)

	d.SetId(fmt.Sprintf("%d:%s:%s", createAttackGroup.ConfigID, createAttackGroup.PolicyID, createAttackGroup.Group))

	return resourceEvalGroupRead(ctx, d, m)
}

func resourceEvalGroupRead(ctx context.Context, d *schema.ResourceData, m interface{}) diag.Diagnostics {
	meta := akamai.Meta(m)
	client := inst.Client(meta)
	logger := meta.Log("APPSEC", "resourceEvalGroupRead")
	logger.Debugf("in resourceEvalGroupRead")

	idParts, err := splitID(d.Id(), 3, "configid:securitypolicyid:group")
	if err != nil {
		return diag.FromErr(err)
	}
	configid, err := strconv.Atoi(idParts[0])
	if err != nil {
		return diag.FromErr(err)
	}
	version := getLatestConfigVersion(ctx, configid, m)
	policyid := idParts[1]
	group := idParts[2]

	getAttackGroup := appsec.GetAttackGroupRequest{
		ConfigID: configid,
		Version: version,
		PolicyID: policyid,
		Group: group,
	}

	attackgroup, err := client.GetEvalGroup(ctx, getAttackGroup)
	if err != nil {
		logger.Warnf("calling 'getEvalGroup': %s", err.Error())
	}

	if err := d.Set("config_id", getAttackGroup.ConfigID); err != nil {
		return diag.FromErr(fmt.Errorf("%w: %s", tools.ErrValueSet, err.Error()))
	}
	if err := d.Set("security_policy_id", getAttackGroup.PolicyID); err != nil {
		return diag.FromErr(fmt.Errorf("%w: %s", tools.ErrValueSet, err.Error()))
	}
	if err := d.Set("attack_group", getAttackGroup.Group); err != nil {
		return diag.FromErr(fmt.Errorf("%w: %s", tools.ErrValueSet, err.Error()))
	}
	if err := d.Set("attack_group_action", string(attackgroup.Action)); err != nil {
		return diag.FromErr(fmt.Errorf("%w: %s", tools.ErrValueSet, err.Error()))
	}

	if !attackgroup.IsEmptyConditionException() {
		jsonBody, err := json.Marshal(attackgroup.ConditionException)
		if err != nil {
			return diag.Errorf("%s", "Error Marshalling condition exception")
		}
		if err := d.Set("condition_exception", string(jsonBody)); err != nil {
			return diag.FromErr(fmt.Errorf("%w: %s", tools.ErrValueSet, err.Error()))
		}
	}

	return nil
}

func resourceEvalGroupUpdate(ctx context.Context, d *schema.ResourceData, m interface{}) diag.Diagnostics {
	meta := akamai.Meta(m)
	client := inst.Client(meta)
	logger := meta.Log("APPSEC", "resourceEvalGroupUpdate")
	logger.Debugf("in resourceEvalGroupUpdate")

	idParts, err := splitID(d.Id(), 3, "configid:securitypolicyid:group")
	if err != nil {
		return diag.FromErr(err)
	}
	configid, err := strconv.Atoi(idParts[0])
	if err != nil {
		return diag.FromErr(err)
	}
	version := getModifiableConfigVersion(ctx, configid, "evalGroup", m)
	policyid := idParts[1]
	group := idParts[2]

	action, err := tools.GetStringValue("attack_group_action", d)
	if err != nil {
		return diag.FromErr(err)
	}
	conditionexception, err := tools.GetStringValue("condition_exception", d)
	if err != nil && !errors.Is(err, tools.ErrNotFound) {
		return diag.FromErr(err)
	}

	if err := validateActionAndConditionException(action, conditionexception); err != nil {
		return diag.FromErr(err)
	}

	rawJSON := json.RawMessage(conditionexception)

	updateAttackGroup := appsec.UpdateAttackGroupRequest{
		ConfigID: configid,
		Version: version,
		PolicyID: policyid,
		Group: group,
		Action: action,
		JsonPayloadRaw: rawJSON,
	}

	_, err = client.UpdateEvalGroup(ctx, updateAttackGroup)
	if err != nil {
		logger.Errorf("calling 'updateEvalGroup': %s", err.Error())
		return diag.FromErr(err)
	}

	return resourceEvalGroupRead(ctx, d, m)
}

func resourceEvalGroupDelete(ctx context.Context, d *schema.ResourceData, m interface{}) diag.Diagnostics {
	meta := akamai.Meta(m)
	client := inst.Client(meta)
	logger := meta.Log("APPSEC", "resourceEvalgroupDelete")
	logger.Debugf("in resourceEvalGroupDelete")

	idParts, err := splitID(d.Id(), 3, "configid:securitypolicyid:group")
	if err != nil {
		return diag.FromErr(err)
	}
	configid, err := strconv.Atoi(idParts[0])
	if err != nil {
		return diag.FromErr(err)
	}
	version := getModifiableConfigVersion(ctx, configid, "evalGroup", m)
	policyid := idParts[1]
	group := idParts[2]

	removeAttackGroup := appsec.UpdateAttackGroupRequest{
		ConfigID: configid,
		Version: version,
		PolicyID: policyid,
		Group: group,
		Action: "none",
	}

	_, err = client.UpdateEvalGroup(ctx, removeAttackGroup)
	if err != nil {
		logger.Errorf("calling 'RemoveEvalGroup': %s", err.Error())
		return diag.FromErr(err)
	}

	d.SetId("")

	return nil
}
