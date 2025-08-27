package appsec

import (
	"context"
	"fmt"
	"strconv"

	"github.com/akamai/AkamaiOPEN-edgegrid-golang/v12/pkg/appsec"
	"github.com/akamai/terraform-provider-akamai/v8/pkg/common/tf"
	"github.com/akamai/terraform-provider-akamai/v8/pkg/meta"
	"github.com/hashicorp/terraform-plugin-sdk/v2/diag"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/customdiff"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/schema"
)

func resourceAdvancedSettingsPIILearning() *schema.Resource {
	return &schema.Resource{
		CreateContext: resourceAdvancedSettingsPIILearningCreate,
		ReadContext:   resourceAdvancedSettingsPIILearningRead,
		UpdateContext: resourceAdvancedSettingsPIILearningUpdate,
		DeleteContext: resourceAdvancedSettingsPIILearningDelete,
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
			"enable_pii_learning": {
				Type:        schema.TypeBool,
				Required:    true,
				Description: "Whether to enable the PII learning advanced setting",
			},
		},
	}
}

func resourceAdvancedSettingsPIILearningCreate(ctx context.Context, d *schema.ResourceData, m interface{}) diag.Diagnostics {
	meta := meta.Must(m)
	client := inst.Client(meta)
	logger := meta.Log("APPSEC", "resourceAdvancedSettingsPIILearningCreate")

	configID, err := tf.GetIntValue("config_id", d)
	if err != nil {
		return diag.FromErr(err)
	}
	version, err := getModifiableConfigVersion(ctx, configID, "piiLearningSetting", m)
	if err != nil {
		return diag.FromErr(err)
	}
	enablePIILearning, err := tf.GetBoolValue("enable_pii_learning", d)
	if err != nil {
		return diag.FromErr(err)
	}

	_, err = client.UpdateAdvancedSettingsPIILearning(ctx, appsec.UpdateAdvancedSettingsPIILearningRequest{
		ConfigVersion: appsec.ConfigVersion{
			ConfigID: int64(configID),
			Version:  version,
		},
		EnablePIILearning: enablePIILearning,
	})
	if err != nil {
		logger.Errorf("calling 'updateAdvancedSettingsPIILearning': %s", err.Error())
		return diag.FromErr(err)
	}

	d.SetId(fmt.Sprintf("%d", configID))

	return resourceAdvancedSettingsPIILearningRead(ctx, d, m)
}

func resourceAdvancedSettingsPIILearningRead(ctx context.Context, d *schema.ResourceData, m interface{}) diag.Diagnostics {
	meta := meta.Must(m)
	client := inst.Client(meta)
	logger := meta.Log("APPSEC", "resourceAdvancedSettingsPIILearningRead")

	configID, err := strconv.Atoi(d.Id())
	if err != nil {
		return diag.FromErr(err)
	}
	version, err := getLatestConfigVersion(ctx, configID, m)
	if err != nil {
		return diag.FromErr(err)
	}

	advancedSettingsPIILearning, err := client.GetAdvancedSettingsPIILearning(ctx, appsec.GetAdvancedSettingsPIILearningRequest{
		ConfigVersion: appsec.ConfigVersion{
			ConfigID: int64(configID),
			Version:  version,
		},
	})
	if err != nil {
		logger.Errorf("calling 'getAdvancedSettingsPIILearning': %s", err.Error())
		return diag.FromErr(err)
	}

	if err := d.Set("config_id", configID); err != nil {
		return diag.Errorf("%s: %s", tf.ErrValueSet, err.Error())
	}
	if err := d.Set("enable_pii_learning", advancedSettingsPIILearning.EnablePIILearning); err != nil {
		return diag.Errorf("%s: %s", tf.ErrValueSet, err.Error())
	}

	return nil
}

func resourceAdvancedSettingsPIILearningUpdate(ctx context.Context, d *schema.ResourceData, m interface{}) diag.Diagnostics {
	meta := meta.Must(m)
	client := inst.Client(meta)
	logger := meta.Log("APPSEC", "resourceAdvancedSettingsPIILearningUpdate")

	configID, err := strconv.Atoi(d.Id())
	if err != nil {
		return diag.FromErr(err)
	}
	version, err := getModifiableConfigVersion(ctx, configID, "piiLearningSetting", m)
	if err != nil {
		return diag.FromErr(err)
	}
	enablePIILearning, err := tf.GetBoolValue("enable_pii_learning", d)
	if err != nil {
		return diag.FromErr(err)
	}

	_, err = client.UpdateAdvancedSettingsPIILearning(ctx, appsec.UpdateAdvancedSettingsPIILearningRequest{
		ConfigVersion: appsec.ConfigVersion{
			ConfigID: int64(configID),
			Version:  version,
		},
		EnablePIILearning: enablePIILearning,
	})
	if err != nil {
		logger.Errorf("calling 'updateAdvancedSettingsPIILearning': %s", err.Error())
		return diag.FromErr(err)
	}

	return resourceAdvancedSettingsPIILearningRead(ctx, d, m)
}

func resourceAdvancedSettingsPIILearningDelete(ctx context.Context, d *schema.ResourceData, m interface{}) diag.Diagnostics {
	meta := meta.Must(m)
	client := inst.Client(meta)
	logger := meta.Log("APPSEC", "resourceAdvancedSettingsPIILearningDelete")

	configID, err := strconv.Atoi(d.Id())
	if err != nil {
		return diag.FromErr(err)
	}
	version, err := getModifiableConfigVersion(ctx, configID, "piiLearningSetting", m)
	if err != nil {
		return diag.FromErr(err)
	}

	_, err = client.UpdateAdvancedSettingsPIILearning(ctx, appsec.UpdateAdvancedSettingsPIILearningRequest{
		ConfigVersion: appsec.ConfigVersion{
			ConfigID: int64(configID),
			Version:  version,
		},
		EnablePIILearning: false,
	})
	if err != nil {
		logger.Errorf("calling 'updateAdvancedSettingsPIILearning': %s", err.Error())
		return diag.FromErr(err)
	}
	return nil
}
