package appsec

import (
	"context"
	"errors"
	"fmt"
	"strconv"

	"github.com/akamai/AkamaiOPEN-edgegrid-golang/v12/pkg/appsec"
	"github.com/akamai/terraform-provider-akamai/v9/pkg/common/tf"
	"github.com/akamai/terraform-provider-akamai/v9/pkg/meta"
	"github.com/hashicorp/terraform-plugin-sdk/v2/diag"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/customdiff"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/schema"
)

func resourceAdvancedSettingsAsePenaltyBox() *schema.Resource {
	return &schema.Resource{
		CreateContext: resourceAdvancedSettingsAsePenaltyBoxCreate,
		ReadContext:   resourceAdvancedSettingsAsePenaltyBoxRead,
		UpdateContext: resourceAdvancedSettingsAsePenaltyBoxUpdate,
		DeleteContext: resourceAdvancedSettingsAsePenaltyBoxDelete,
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
				Description: "Unique identifier of the security configuration.",
			},
			"block_duration": {
				Type:        schema.TypeInt,
				Required:    true,
				Description: "Block duration for ASE Penalty Box in minutes.",
			},
			"qualification_exclusions": {
				Type:        schema.TypeList,
				Required:    true,
				MaxItems:    1,
				Description: "Qualification exclusions for ASE Penalty Box. Contains attack groups and rules.",
				Elem: &schema.Resource{
					Schema: map[string]*schema.Schema{
						"attack_groups": {
							Type:        schema.TypeSet,
							Elem:        &schema.Schema{Type: schema.TypeString},
							Optional:    true,
							Description: "List of attack group names.",
						},
						"rules": {
							Type:        schema.TypeSet,
							Elem:        &schema.Schema{Type: schema.TypeInt},
							Optional:    true,
							Description: "List of rule IDs.",
						},
					},
				},
			},
		},
	}
}

func resourceAdvancedSettingsAsePenaltyBoxCreate(ctx context.Context, d *schema.ResourceData, m interface{}) diag.Diagnostics {
	meta := meta.Must(m)
	logger := meta.Log("APPSEC", "resourceAdvancedSettingsAsePenaltyBoxCreate")
	logger.Debugf("in resourceAdvancedSettingsAsePenaltyBoxCreate")

	return upsertAdvancedSettingsAsePenaltyBox(ctx, d, m)
}

func upsertAdvancedSettingsAsePenaltyBox(ctx context.Context, d *schema.ResourceData, m interface{}) diag.Diagnostics {
	meta := meta.Must(m)
	client := inst.Client(meta)

	configID, err := tf.GetIntValue("config_id", d)
	if err != nil {
		return diag.FromErr(err)
	}
	version, err := getModifiableConfigVersion(ctx, configID, "AsePenaltyBoxSetting", m)
	if err != nil {
		return diag.FromErr(err)
	}
	blockDuration, err := tf.GetIntValue("block_duration", d)
	if err != nil {
		return diag.FromErr(err)
	}
	qualificationExclusionsList, err := tf.GetListValue("qualification_exclusions", d)
	if err != nil && !errors.Is(err, tf.ErrNotFound) {
		return diag.FromErr(err)
	}

	qualificationExclusionsValue := appsec.QualificationExclusions{}
	if len(qualificationExclusionsList) != 0 && qualificationExclusionsList[0] != nil {
		qualificationExclusionsMap, ok := qualificationExclusionsList[0].(map[string]interface{})
		if !ok {
			return diag.FromErr(fmt.Errorf("expected map[string]interface{}, got: %T", qualificationExclusionsList[0]))
		}
		var attackGroups []string
		if v, exists := qualificationExclusionsMap["attack_groups"]; exists && v != nil {
			if attackGroupsSet, ok := v.(*schema.Set); ok {
				for _, val := range attackGroupsSet.List() {
					if s, ok := val.(string); ok {
						attackGroups = append(attackGroups, s)
					} else {
						return diag.FromErr(fmt.Errorf("attack_groups contains non-string value: %T", val))
					}
				}
			} else {
				return diag.FromErr(fmt.Errorf("attack_groups is not a *schema.Set, got: %T", v))
			}
		}
		var rules []int
		if v, exists := qualificationExclusionsMap["rules"]; exists && v != nil {
			if rulesSet, ok := v.(*schema.Set); ok {
				for _, val := range rulesSet.List() {
					if i, ok := val.(int); ok {
						rules = append(rules, i)
					} else {
						return diag.FromErr(fmt.Errorf("rules contains non-int value: %T", val))
					}
				}
			} else {
				return diag.FromErr(fmt.Errorf("rules is not a *schema.Set, got: %T", v))
			}
		}

		qualificationExclusionsValue.AttackGroups = attackGroups
		qualificationExclusionsValue.Rules = rules
	}

	req := appsec.UpdateAdvancedSettingsAsePenaltyBoxRequest{
		ConfigID:                configID,
		Version:                 version,
		BlockDuration:           blockDuration,
		QualificationExclusions: &qualificationExclusionsValue,
	}

	_, err = client.UpdateAdvancedSettingsAsePenaltyBox(ctx, req)
	if err != nil {
		return diag.FromErr(err)
	}

	d.SetId(fmt.Sprintf("%d", configID))

	return resourceAdvancedSettingsAsePenaltyBoxRead(ctx, d, m)
}

func resourceAdvancedSettingsAsePenaltyBoxRead(ctx context.Context, d *schema.ResourceData, m interface{}) diag.Diagnostics {
	meta := meta.Must(m)
	client := inst.Client(meta)
	logger := meta.Log("APPSEC", "resourceAdvancedSettingsAsePenaltyBoxRead")
	logger.Debugf("in resourceAdvancedSettingsAsePenaltyBoxRead")

	configID, err := strconv.Atoi(d.Id())
	if err != nil {
		return diag.FromErr(err)
	}
	version, err := getLatestConfigVersion(ctx, configID, m)
	if err != nil {
		return diag.FromErr(err)
	}

	advancedSettingsAsePenaltyBox, err := client.GetAdvancedSettingsAsePenaltyBox(ctx, appsec.GetAdvancedSettingsAsePenaltyBoxRequest{
		ConfigID: configID,
		Version:  version,
	})
	if err != nil {
		logger.Errorf("calling 'getAdvancedSettingsAsePenaltyBox': %s", err.Error())
		return diag.FromErr(err)
	}

	if err := d.Set("config_id", configID); err != nil {
		return diag.Errorf("%s: %s", tf.ErrValueSet, err.Error())
	}
	if err := d.Set("block_duration", advancedSettingsAsePenaltyBox.BlockDuration); err != nil {
		return diag.Errorf("%s: %s", tf.ErrValueSet, err.Error())
	}
	if err := d.Set("qualification_exclusions", flattenQualificationExclusions(advancedSettingsAsePenaltyBox.QualificationExclusions)); err != nil {
		return diag.Errorf("%s: %s", tf.ErrValueSet, err.Error())
	}
	return nil
}

func resourceAdvancedSettingsAsePenaltyBoxUpdate(ctx context.Context, d *schema.ResourceData, m interface{}) diag.Diagnostics {
	meta := meta.Must(m)
	logger := meta.Log("APPSEC", "resourceAdvancedSettingsAsePenaltyBoxUpdate")
	logger.Debugf("in resourceAdvancedSettingsAsePenaltyBoxUpdate")

	return upsertAdvancedSettingsAsePenaltyBox(ctx, d, m)
}

func resourceAdvancedSettingsAsePenaltyBoxDelete(ctx context.Context, d *schema.ResourceData, m interface{}) diag.Diagnostics {
	meta := meta.Must(m)
	client := inst.Client(meta)
	logger := meta.Log("APPSEC", "resourceAdvancedSettingsAsePenaltyBoxDelete")
	logger.Debugf("in resourceAdvancedSettingsAsePenaltyBoxDelete")

	configID, err := strconv.Atoi(d.Id())
	if err != nil {
		return diag.FromErr(err)
	}
	version, err := getModifiableConfigVersion(ctx, configID, "AsePenaltyBoxSetting", m)
	if err != nil {
		return diag.FromErr(err)
	}

	req := appsec.RemoveAdvancedSettingsAsePenaltyBoxRequest{
		ConfigID: configID,
		Version:  version,
	}

	_, err = client.RemoveAdvancedSettingsAsePenaltyBox(ctx, req)
	if err != nil {
		return diag.FromErr(err)
	}
	return nil
}

func flattenQualificationExclusions(qe *appsec.QualificationExclusions) []map[string]interface{} {
	if qe == nil {
		return nil
	}
	m := map[string]interface{}{
		"attack_groups": qe.AttackGroups,
		"rules":         qe.Rules,
	}
	return []map[string]interface{}{m}
}
