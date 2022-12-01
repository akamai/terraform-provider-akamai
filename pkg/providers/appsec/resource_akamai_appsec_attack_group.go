package appsec

import (
	"context"
	"encoding/json"
	"errors"
	"fmt"
	"strconv"

	"github.com/akamai/AkamaiOPEN-edgegrid-golang/v3/pkg/appsec"
	"github.com/akamai/terraform-provider-akamai/v3/pkg/akamai"
	"github.com/akamai/terraform-provider-akamai/v3/pkg/tools"
	"github.com/hashicorp/terraform-plugin-sdk/v2/diag"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/customdiff"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/schema"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/validation"
)

// appsec v1
//
// https://techdocs.akamai.com/application-security/reference/api
func resourceAttackGroup() *schema.Resource {
	return &schema.Resource{
		CreateContext: resourceAttackGroupCreate,
		ReadContext:   resourceAttackGroupRead,
		UpdateContext: resourceAttackGroupUpdate,
		DeleteContext: resourceAttackGroupDelete,
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
				Description: "Unique name of the attack group to be modified",
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
				ValidateDiagFunc: validation.ToDiagFunc(validation.StringIsJSON),
				DiffSuppressFunc: suppressEquivalentJSONDiffsGeneric,
				Description:      "JSON-formatted condition and exception information for the attack group",
			},
		},
	}
}

func resourceAttackGroupCreate(ctx context.Context, d *schema.ResourceData, m interface{}) diag.Diagnostics {
	meta := akamai.Meta(m)
	client := inst.Client(meta)
	logger := meta.Log("APPSEC", "resourceAttackGroupCreate")
	logger.Debugf(" in resourceAttackGroupCreate")

	configID, err := tools.GetIntValue("config_id", d)
	if err != nil {
		return diag.FromErr(err)
	}
	version, err := getModifiableConfigVersion(ctx, configID, "atackGroup", m)
	if err != nil {
		return diag.FromErr(err)
	}
	policyID, err := tools.GetStringValue("security_policy_id", d)
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
	jsonPayloadRaw := []byte(conditionexception)
	rawJSON := (json.RawMessage)(jsonPayloadRaw)

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

	_, err = client.UpdateAttackGroup(ctx, createAttackGroup)
	if err != nil {
		return diag.FromErr(err)
	}

	d.SetId(fmt.Sprintf("%d:%s:%s", createAttackGroup.ConfigID, createAttackGroup.PolicyID, createAttackGroup.Group))

	return resourceAttackGroupRead(ctx, d, m)
}

func resourceAttackGroupRead(ctx context.Context, d *schema.ResourceData, m interface{}) diag.Diagnostics {
	meta := akamai.Meta(m)
	client := inst.Client(meta)
	logger := meta.Log("APPSEC", "resourceAttackGroupRead")
	logger.Debugf(" in resourceAttackGroupRead")

	iDParts, err := splitID(d.Id(), 3, "configID:securityPolicyID:group")
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

	attackgroup, err := client.GetAttackGroup(ctx, getAttackGroup)
	if err != nil {
		logger.Warnf("calling 'getAttackGroup': %s", err.Error())
	}

	if err := d.Set("config_id", getAttackGroup.ConfigID); err != nil {
		return diag.Errorf("%s: %s", tools.ErrValueSet, err.Error())
	}
	if err := d.Set("security_policy_id", getAttackGroup.PolicyID); err != nil {
		return diag.Errorf("%s: %s", tools.ErrValueSet, err.Error())
	}
	if err := d.Set("attack_group", getAttackGroup.Group); err != nil {
		return diag.Errorf("%s: %s", tools.ErrValueSet, err.Error())
	}
	if err := d.Set("attack_group_action", string(attackgroup.Action)); err != nil {
		return diag.Errorf("%s: %s", tools.ErrValueSet, err.Error())
	}

	if !attackgroup.IsEmptyConditionException() {
		jsonBody, err := json.Marshal(attackgroup.ConditionException)
		if err != nil {
			return diag.Errorf("%s", "Error Marshalling condition exception")
		}
		if err := d.Set("condition_exception", string(jsonBody)); err != nil {
			return diag.Errorf("%s: %s", tools.ErrValueSet, err.Error())
		}
	}

	return nil
}

func resourceAttackGroupUpdate(ctx context.Context, d *schema.ResourceData, m interface{}) diag.Diagnostics {
	meta := akamai.Meta(m)
	client := inst.Client(meta)
	logger := meta.Log("APPSEC", "resourceAttackGroupUpdate")
	logger.Debugf(" in resourceAttackGroupUpdate")

	iDParts, err := splitID(d.Id(), 3, "configID:securityPolicyID:group")
	if err != nil {
		return diag.FromErr(err)
	}
	configID, err := strconv.Atoi(iDParts[0])
	if err != nil {
		return diag.FromErr(err)
	}
	version, err := getModifiableConfigVersion(ctx, configID, "attackGroup", m)
	if err != nil {
		return diag.FromErr(err)
	}
	policyID := iDParts[1]
	group := iDParts[2]

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

	jsonPayloadRaw := []byte(conditionexception)
	rawJSON := (json.RawMessage)(jsonPayloadRaw)

	updateAttackGroup := appsec.UpdateAttackGroupRequest{
		ConfigID:       configID,
		Version:        version,
		PolicyID:       policyID,
		Group:          group,
		Action:         action,
		JsonPayloadRaw: rawJSON,
	}

	_, err = client.UpdateAttackGroup(ctx, updateAttackGroup)
	if err != nil {
		return diag.FromErr(err)
	}

	return resourceAttackGroupRead(ctx, d, m)
}

func resourceAttackGroupDelete(ctx context.Context, d *schema.ResourceData, m interface{}) diag.Diagnostics {
	meta := akamai.Meta(m)
	client := inst.Client(meta)
	logger := meta.Log("APPSEC", "resourceAttackGroupDelete")
	logger.Debugf(" in resourceAttackGroupDelete")

	iDParts, err := splitID(d.Id(), 3, "configID:securityPolicyID:group")
	if err != nil {
		return diag.FromErr(err)
	}
	configID, err := strconv.Atoi(iDParts[0])
	if err != nil {
		return diag.FromErr(err)
	}
	version, err := getModifiableConfigVersion(ctx, configID, "attackGroup", m)
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

	_, err = client.UpdateAttackGroup(ctx, removeAttackGroup)
	if err != nil {
		logger.Errorf("calling 'RemoveAttackGroup': %s", err.Error())
		return diag.FromErr(err)
	}
	return nil
}
