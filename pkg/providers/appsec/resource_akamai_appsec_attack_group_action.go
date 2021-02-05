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
func resourceAttackGroupAction() *schema.Resource {
	return &schema.Resource{
		CreateContext: resourceAttackGroupActionUpdate,
		ReadContext:   resourceAttackGroupActionRead,
		UpdateContext: resourceAttackGroupActionUpdate,
		DeleteContext: resourceAttackGroupActionDelete,
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
			"attack_group_action": {
				Type:         schema.TypeString,
				Required:     true,
				ValidateFunc: ValidateActions,
			},
		},
	}
}

func resourceAttackGroupActionRead(ctx context.Context, d *schema.ResourceData, m interface{}) diag.Diagnostics {
	meta := akamai.Meta(m)
	client := inst.Client(meta)
	logger := meta.Log("APPSEC", "resourceAttackGroupActionRead")

	getAttackGroupAction := appsec.GetAttackGroupActionRequest{}
	if d.Id() != "" && strings.Contains(d.Id(), ":") {
		s := strings.Split(d.Id(), ":")

		configid, errconv := strconv.Atoi(s[0])
		if errconv != nil {
			return diag.FromErr(errconv)
		}
		getAttackGroupAction.ConfigID = configid

		version, errconv := strconv.Atoi(s[1])
		if errconv != nil {
			return diag.FromErr(errconv)
		}
		getAttackGroupAction.Version = version

		policyid := s[2]

		getAttackGroupAction.PolicyID = policyid

		attackgroup := s[3]

		getAttackGroupAction.Group = attackgroup

	} else {
		configid, err := tools.GetIntValue("config_id", d)
		if err != nil && !errors.Is(err, tools.ErrNotFound) {
			return diag.FromErr(err)
		}
		getAttackGroupAction.ConfigID = configid

		version, err := tools.GetIntValue("version", d)
		if err != nil && !errors.Is(err, tools.ErrNotFound) {
			return diag.FromErr(err)
		}
		getAttackGroupAction.Version = version

		policyid, err := tools.GetStringValue("security_policy_id", d)
		if err != nil && !errors.Is(err, tools.ErrNotFound) {
			return diag.FromErr(err)
		}
		getAttackGroupAction.PolicyID = policyid

		attackgroup, err := tools.GetStringValue("attack_group", d)
		if err != nil && !errors.Is(err, tools.ErrNotFound) {
			return diag.FromErr(err)
		}
		getAttackGroupAction.Group = attackgroup
	}
	attackgroupaction, err := client.GetAttackGroupAction(ctx, getAttackGroupAction)
	if err != nil {
		logger.Errorf("calling 'getAttackGroupAction': %s", err.Error())
		return diag.FromErr(err)
	}

	if err := d.Set("config_id", getAttackGroupAction.ConfigID); err != nil {
		return diag.FromErr(fmt.Errorf("%w: %s", tools.ErrValueSet, err.Error()))
	}

	if err := d.Set("version", getAttackGroupAction.Version); err != nil {
		return diag.FromErr(fmt.Errorf("%w: %s", tools.ErrValueSet, err.Error()))
	}

	if err := d.Set("security_policy_id", getAttackGroupAction.PolicyID); err != nil {
		return diag.FromErr(fmt.Errorf("%w: %s", tools.ErrValueSet, err.Error()))
	}

	if err := d.Set("attack_group", getAttackGroupAction.Group); err != nil {
		return diag.FromErr(fmt.Errorf("%w: %s", tools.ErrValueSet, err.Error()))
	}

	if err := d.Set("attack_group_action", attackgroupaction.Action); err != nil {
		return diag.FromErr(fmt.Errorf("%w: %s", tools.ErrValueSet, err.Error()))
	}

	d.SetId(fmt.Sprintf("%d:%d:%s:%s", getAttackGroupAction.ConfigID, getAttackGroupAction.Version, getAttackGroupAction.PolicyID, getAttackGroupAction.Group))

	return nil
}

func resourceAttackGroupActionDelete(ctx context.Context, d *schema.ResourceData, m interface{}) diag.Diagnostics {

	meta := akamai.Meta(m)
	client := inst.Client(meta)
	logger := meta.Log("APPSEC", "resourceAttackGroupActionRemove")

	removeAttackGroupAction := appsec.UpdateAttackGroupActionRequest{}
	if d.Id() != "" && strings.Contains(d.Id(), ":") {
		s := strings.Split(d.Id(), ":")

		configid, errconv := strconv.Atoi(s[0])
		if errconv != nil {
			return diag.FromErr(errconv)
		}
		removeAttackGroupAction.ConfigID = configid

		version, errconv := strconv.Atoi(s[1])
		if errconv != nil {
			return diag.FromErr(errconv)
		}
		removeAttackGroupAction.Version = version

		policyid := s[2]

		removeAttackGroupAction.PolicyID = policyid

		attackgroup := s[3]

		removeAttackGroupAction.Group = attackgroup

	} else {
		configid, err := tools.GetIntValue("config_id", d)
		if err != nil && !errors.Is(err, tools.ErrNotFound) {
			return diag.FromErr(err)
		}
		removeAttackGroupAction.ConfigID = configid

		version, err := tools.GetIntValue("version", d)
		if err != nil && !errors.Is(err, tools.ErrNotFound) {
			return diag.FromErr(err)
		}
		removeAttackGroupAction.Version = version

		policyid, err := tools.GetStringValue("security_policy_id", d)
		if err != nil && !errors.Is(err, tools.ErrNotFound) {
			return diag.FromErr(err)
		}
		removeAttackGroupAction.PolicyID = policyid

		attackgroup, err := tools.GetStringValue("attack_group", d)
		if err != nil && !errors.Is(err, tools.ErrNotFound) {
			return diag.FromErr(err)
		}
		removeAttackGroupAction.Group = attackgroup
	}
	removeAttackGroupAction.Action = "none"

	_, erru := client.UpdateAttackGroupAction(ctx, removeAttackGroupAction)
	if erru != nil {
		logger.Errorf("calling 'removeAttackGroupAction': %s", erru.Error())
		return diag.FromErr(erru)
	}

	d.SetId("")
	return nil
}

func resourceAttackGroupActionUpdate(ctx context.Context, d *schema.ResourceData, m interface{}) diag.Diagnostics {
	meta := akamai.Meta(m)
	client := inst.Client(meta)
	logger := meta.Log("APPSEC", "resourceAttackGroupActionUpdate")

	updateAttackGroupAction := appsec.UpdateAttackGroupActionRequest{}
	if d.Id() != "" && strings.Contains(d.Id(), ":") {
		s := strings.Split(d.Id(), ":")

		configid, errconv := strconv.Atoi(s[0])
		if errconv != nil {
			return diag.FromErr(errconv)
		}
		updateAttackGroupAction.ConfigID = configid

		version, errconv := strconv.Atoi(s[1])
		if errconv != nil {
			return diag.FromErr(errconv)
		}
		updateAttackGroupAction.Version = version

		policyid := s[2]

		updateAttackGroupAction.PolicyID = policyid

		attackgroup := s[3]

		updateAttackGroupAction.Group = attackgroup

	} else {
		configid, err := tools.GetIntValue("config_id", d)
		if err != nil && !errors.Is(err, tools.ErrNotFound) {
			return diag.FromErr(err)
		}
		updateAttackGroupAction.ConfigID = configid

		version, err := tools.GetIntValue("version", d)
		if err != nil && !errors.Is(err, tools.ErrNotFound) {
			return diag.FromErr(err)
		}
		updateAttackGroupAction.Version = version

		policyid, err := tools.GetStringValue("security_policy_id", d)
		if err != nil && !errors.Is(err, tools.ErrNotFound) {
			return diag.FromErr(err)
		}
		updateAttackGroupAction.PolicyID = policyid

		attackgroup, err := tools.GetStringValue("attack_group", d)
		if err != nil && !errors.Is(err, tools.ErrNotFound) {
			return diag.FromErr(err)
		}
		updateAttackGroupAction.Group = attackgroup
	}
	attackgroupaction, err := tools.GetStringValue("attack_group_action", d)
	if err != nil && !errors.Is(err, tools.ErrNotFound) {
		return diag.FromErr(err)
	}
	updateAttackGroupAction.Action = attackgroupaction

	_, erru := client.UpdateAttackGroupAction(ctx, updateAttackGroupAction)
	if erru != nil {
		logger.Errorf("calling 'updateAttackGroupAction': %s", erru.Error())
		return diag.FromErr(erru)
	}

	return resourceAttackGroupActionRead(ctx, d, m)
}
