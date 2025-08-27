package appsec

import (
	"context"
	"encoding/json"
	"errors"
	"fmt"
	"strconv"

	"github.com/akamai/AkamaiOPEN-edgegrid-golang/v12/pkg/appsec"
	"github.com/akamai/terraform-provider-akamai/v8/pkg/common/id"
	"github.com/akamai/terraform-provider-akamai/v8/pkg/common/tf"
	"github.com/akamai/terraform-provider-akamai/v8/pkg/meta"
	"github.com/hashicorp/terraform-plugin-sdk/v2/diag"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/customdiff"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/schema"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/validation"
)

// appsec v1
//
// https://techdocs.akamai.com/application-security/reference/api
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
				Type:        schema.TypeInt,
				Required:    true,
				Description: "Unique identifier of the security configuration",
			},
			"security_policy_id": {
				Type:        schema.TypeString,
				Required:    true,
				Description: "Unique identifier of the security policy",
			},
			"attack_group": {
				Type:        schema.TypeString,
				Required:    true,
				Description: "Unique identifier of the evaluation attack group being modified",
			},
			"attack_group_action": {
				Type:             schema.TypeString,
				Required:         true,
				ValidateDiagFunc: ValidateActions,
				Description:      "Action to be taken when the attack group is triggered",
			},
			"condition_exception": {
				Type:             schema.TypeString,
				Optional:         true,
				Default:          "",
				ValidateFunc:     validation.StringIsJSON,
				DiffSuppressFunc: suppressEquivalentJSONDiffsGeneric,
				Description:      "JSON-formatted condition and exception information for the evaluation attack group",
			},
		},
	}
}

func resourceEvalGroupCreate(ctx context.Context, d *schema.ResourceData, m interface{}) diag.Diagnostics {
	meta := meta.Must(m)
	client := inst.Client(meta)
	logger := meta.Log("APPSEC", "resourceEvalGroupCreate")
	logger.Debugf("in resourceEvalGroupCreate")

	configID, err := tf.GetIntValue("config_id", d)
	if err != nil {
		return diag.FromErr(err)
	}
	version, err := getModifiableConfigVersion(ctx, configID, "atackGroup", m)
	if err != nil {
		return diag.FromErr(err)
	}
	policyID, err := tf.GetStringValue("security_policy_id", d)
	if err != nil {
		return diag.FromErr(err)
	}
	attackgroup, err := tf.GetStringValue("attack_group", d)
	if err != nil {
		return diag.FromErr(err)
	}
	action, err := tf.GetStringValue("attack_group_action", d)
	if err != nil {
		return diag.FromErr(err)
	}
	conditionexception, err := tf.GetStringValue("condition_exception", d)
	if err != nil && !errors.Is(err, tf.ErrNotFound) {
		return diag.FromErr(err)
	}

	rawJSON := json.RawMessage(conditionexception)

	if err := validateActionAndConditionException(action, conditionexception); err != nil {
		return diag.FromErr(err)
	}

	createAttackGroup := appsec.UpdateAttackGroupRequest{
		ConfigID:       configID,
		Version:        version,
		PolicyID:       policyID,
		Group:          attackgroup,
		Action:         action,
		JsonPayloadRaw: rawJSON,
	}

	_, err = client.UpdateEvalGroup(ctx, createAttackGroup)
	if err != nil {
		logger.Errorf("calling 'createEvalGroup': %s", err.Error())
		return diag.FromErr(err)
	}

	d.SetId(fmt.Sprintf("%d:%s:%s", createAttackGroup.ConfigID, createAttackGroup.PolicyID, createAttackGroup.Group))

	return resourceEvalGroupRead(ctx, d, m)
}

func resourceEvalGroupRead(ctx context.Context, d *schema.ResourceData, m interface{}) diag.Diagnostics {
	meta := meta.Must(m)
	client := inst.Client(meta)
	logger := meta.Log("APPSEC", "resourceEvalGroupRead")
	logger.Debugf("in resourceEvalGroupRead")

	iDParts, err := id.Split(d.Id(), 3, "configID:securityPolicyID:group")
	if err != nil {
		return diag.FromErr(err)
	}
	configID, err := strconv.Atoi(iDParts[0])
	if err != nil {
		return diag.FromErr(err)
	}
	version, err := getLatestConfigVersion(ctx, configID, m)
	if err != nil {
		return diag.FromErr(err)
	}
	policyID := iDParts[1]
	group := iDParts[2]

	getAttackGroup := appsec.GetAttackGroupRequest{
		ConfigID: configID,
		Version:  version,
		PolicyID: policyID,
		Group:    group,
	}

	attackgroup, err := client.GetEvalGroup(ctx, getAttackGroup)
	if err != nil {
		logger.Warnf("calling 'getEvalGroup': %s", err.Error())
		return diag.FromErr(err)
	}

	if err := d.Set("config_id", getAttackGroup.ConfigID); err != nil {
		return diag.Errorf("%s: %s", tf.ErrValueSet, err.Error())
	}
	if err := d.Set("security_policy_id", getAttackGroup.PolicyID); err != nil {
		return diag.Errorf("%s: %s", tf.ErrValueSet, err.Error())
	}
	if err := d.Set("attack_group", getAttackGroup.Group); err != nil {
		return diag.Errorf("%s: %s", tf.ErrValueSet, err.Error())
	}
	if err := d.Set("attack_group_action", string(attackgroup.Action)); err != nil {
		return diag.Errorf("%s: %s", tf.ErrValueSet, err.Error())
	}

	if !attackgroup.IsEmptyConditionException() {
		jsonBody, err := json.Marshal(attackgroup.ConditionException)
		if err != nil {
			return diag.Errorf("%s", "Error Marshalling condition exception")
		}
		if err := d.Set("condition_exception", string(jsonBody)); err != nil {
			return diag.Errorf("%s: %s", tf.ErrValueSet, err.Error())
		}
	}

	return nil
}

func resourceEvalGroupUpdate(ctx context.Context, d *schema.ResourceData, m interface{}) diag.Diagnostics {
	meta := meta.Must(m)
	client := inst.Client(meta)
	logger := meta.Log("APPSEC", "resourceEvalGroupUpdate")
	logger.Debugf("in resourceEvalGroupUpdate")

	iDParts, err := id.Split(d.Id(), 3, "configID:securityPolicyID:group")
	if err != nil {
		return diag.FromErr(err)
	}
	configID, err := strconv.Atoi(iDParts[0])
	if err != nil {
		return diag.FromErr(err)
	}
	version, err := getModifiableConfigVersion(ctx, configID, "evalGroup", m)
	if err != nil {
		return diag.FromErr(err)
	}
	policyID := iDParts[1]
	group := iDParts[2]

	action, err := tf.GetStringValue("attack_group_action", d)
	if err != nil {
		return diag.FromErr(err)
	}
	conditionexception, err := tf.GetStringValue("condition_exception", d)
	if err != nil && !errors.Is(err, tf.ErrNotFound) {
		return diag.FromErr(err)
	}

	if err := validateActionAndConditionException(action, conditionexception); err != nil {
		return diag.FromErr(err)
	}

	rawJSON := json.RawMessage(conditionexception)

	updateAttackGroup := appsec.UpdateAttackGroupRequest{
		ConfigID:       configID,
		Version:        version,
		PolicyID:       policyID,
		Group:          group,
		Action:         action,
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
	meta := meta.Must(m)
	client := inst.Client(meta)
	logger := meta.Log("APPSEC", "resourceEvalgroupDelete")
	logger.Debugf("in resourceEvalGroupDelete")

	iDParts, err := id.Split(d.Id(), 3, "configID:securityPolicyID:group")
	if err != nil {
		return diag.FromErr(err)
	}
	configID, err := strconv.Atoi(iDParts[0])
	if err != nil {
		return diag.FromErr(err)
	}
	version, err := getModifiableConfigVersion(ctx, configID, "evalGroup", m)
	if err != nil {
		return diag.FromErr(err)
	}
	policyID := iDParts[1]
	group := iDParts[2]

	removeAttackGroup := appsec.UpdateAttackGroupRequest{
		ConfigID: configID,
		Version:  version,
		PolicyID: policyID,
		Group:    group,
		Action:   "none",
	}

	_, err = client.UpdateEvalGroup(ctx, removeAttackGroup)
	if err != nil {
		logger.Errorf("calling 'RemoveEvalGroup': %s", err.Error())
		return diag.FromErr(err)
	}
	return nil
}
