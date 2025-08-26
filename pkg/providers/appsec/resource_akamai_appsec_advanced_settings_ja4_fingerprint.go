package appsec

import (
	"context"
	"errors"
	"fmt"
	"strconv"

	"github.com/akamai/AkamaiOPEN-edgegrid-golang/v11/pkg/appsec"
	"github.com/akamai/terraform-provider-akamai/v8/pkg/common/tf"
	"github.com/akamai/terraform-provider-akamai/v8/pkg/meta"
	"github.com/hashicorp/terraform-plugin-sdk/v2/diag"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/customdiff"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/schema"
)

func resourceAdvancedSettingsJA4Fingerprint() *schema.Resource {
	return &schema.Resource{
		CreateContext: resourceAdvancedSettingsJA4FingerprintCreate,
		ReadContext:   resourceAdvancedSettingsJA4FingerprintRead,
		UpdateContext: resourceAdvancedSettingsJA4FingerprintUpdate,
		DeleteContext: resourceAdvancedSettingsJA4FingerprintDelete,
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
			"header_names": {
				Type:        schema.TypeList,
				Optional:    true,
				Elem:        &schema.Schema{Type: schema.TypeString},
				Description: "JA4 TLS Header Names to be included in the header",
			},
		},
	}
}

func resourceAdvancedSettingsJA4FingerprintCreate(ctx context.Context, d *schema.ResourceData, m interface{}) diag.Diagnostics {
	meta := meta.Must(m)
	client := inst.Client(meta)
	logger := meta.Log("APPSEC", "resourceAdvancedSettingsJA4FingerprintCreate")
	logger.Debugf("in resourceAdvancedSettingsJA4FingerprintCreate")

	configID, err := tf.GetIntValue("config_id", d)
	if err != nil {
		return diag.FromErr(err)
	}
	version, err := getModifiableConfigVersion(ctx, configID, "ja4Fingerprint", m)
	if err != nil {
		return diag.FromErr(err)
	}
	headerNames, err := tf.GetTypedListValue[string]("header_names", d)
	if err != nil && !errors.Is(err, tf.ErrNotFound) {
		return diag.FromErr(err)
	}

	updateAdvancedSettingsJA4FingerprintReq := appsec.UpdateAdvancedSettingsJA4FingerprintRequest{
		ConfigID:    configID,
		Version:     version,
		HeaderNames: headerNames,
	}

	_, err = client.UpdateAdvancedSettingsJA4Fingerprint(ctx, updateAdvancedSettingsJA4FingerprintReq)
	if err != nil {
		logger.Errorf("calling 'UpdateAdvancedSettingsJA4Fingerprint': %s", err.Error())
		return diag.FromErr(err)
	}

	d.SetId(fmt.Sprintf("%d", updateAdvancedSettingsJA4FingerprintReq.ConfigID))

	return resourceAdvancedSettingsJA4FingerprintRead(ctx, d, m)
}

func resourceAdvancedSettingsJA4FingerprintRead(ctx context.Context, d *schema.ResourceData, m interface{}) diag.Diagnostics {
	meta := meta.Must(m)
	client := inst.Client(meta)
	logger := meta.Log("APPSEC", "resourceAdvancedSettingsJA4FingerprintRead")
	logger.Debugf("in resourceAdvancedSettingsJA4FingerprintRead")

	getAdvancedSettingsJA4Fingerprint := appsec.GetAdvancedSettingsJA4FingerprintRequest{}

	configID, err := strconv.Atoi(d.Id())
	if err != nil {
		return diag.FromErr(err)
	}
	version, err := getLatestConfigVersion(ctx, configID, m)
	if err != nil {
		return diag.FromErr(err)
	}

	getAdvancedSettingsJA4Fingerprint.ConfigID = configID
	getAdvancedSettingsJA4Fingerprint.Version = version

	advancedSettingsJA4Fingerprint, err := client.GetAdvancedSettingsJA4Fingerprint(ctx, getAdvancedSettingsJA4Fingerprint)
	if err != nil {
		logger.Errorf("calling 'getAdvancedSettingsJA4Fingerprint': %s", err.Error())
		return diag.FromErr(err)
	}

	if err := d.Set("config_id", getAdvancedSettingsJA4Fingerprint.ConfigID); err != nil {
		return diag.Errorf("%s: %s", tf.ErrValueSet, err.Error())
	}
	if err := d.Set("header_names", advancedSettingsJA4Fingerprint.HeaderNames); err != nil {
		return diag.Errorf("%s: %s", tf.ErrValueSet, err.Error())
	}

	return nil
}

func resourceAdvancedSettingsJA4FingerprintUpdate(ctx context.Context, d *schema.ResourceData, m interface{}) diag.Diagnostics {
	meta := meta.Must(m)
	client := inst.Client(meta)
	logger := meta.Log("APPSEC", "resourceAdvancedSettingsJA4FingerprintUpdate")
	logger.Debugf("in resourceAdvancedSettingsJA4FingerprintUpdate")

	configID, err := strconv.Atoi(d.Id())
	if err != nil {
		return diag.FromErr(err)
	}
	version, err := getModifiableConfigVersion(ctx, configID, "ja4Fingerprint", m)
	if err != nil {
		return diag.FromErr(err)
	}
	headerNames, err := tf.GetTypedListValue[string]("header_names", d)
	if err != nil && !errors.Is(err, tf.ErrNotFound) {
		return diag.FromErr(err)
	}

	updateAdvancedSettingsJA4FingerprintReq := appsec.UpdateAdvancedSettingsJA4FingerprintRequest{}

	updateAdvancedSettingsJA4FingerprintReq.ConfigID = configID
	updateAdvancedSettingsJA4FingerprintReq.Version = version
	updateAdvancedSettingsJA4FingerprintReq.HeaderNames = headerNames

	_, err = client.UpdateAdvancedSettingsJA4Fingerprint(ctx, updateAdvancedSettingsJA4FingerprintReq)
	if err != nil {
		logger.Errorf("calling 'updateAdvancedSettingsJA4Fingerprint': %s", err.Error())
		return diag.FromErr(err)
	}

	return resourceAdvancedSettingsJA4FingerprintRead(ctx, d, m)
}

func resourceAdvancedSettingsJA4FingerprintDelete(ctx context.Context, d *schema.ResourceData, m interface{}) diag.Diagnostics {
	meta := meta.Must(m)
	client := inst.Client(meta)
	logger := meta.Log("APPSEC", "resourceAdvancedSettingsJA4FingerprintDelete")
	logger.Debugf("in resourceAdvancedSettingsJA4FingerprintDelete")

	removeAdvancedSettingsJA4Fingerprint := appsec.RemoveAdvancedSettingsJA4FingerprintRequest{}

	configID, err := strconv.Atoi(d.Id())
	if err != nil {
		return diag.FromErr(err)
	}
	version, err := getModifiableConfigVersion(ctx, configID, "ja4Fingerprint", m)
	if err != nil {
		return diag.FromErr(err)
	}

	removeAdvancedSettingsJA4Fingerprint.ConfigID = configID
	removeAdvancedSettingsJA4Fingerprint.Version = version

	_, err = client.RemoveAdvancedSettingsJA4Fingerprint(ctx, removeAdvancedSettingsJA4Fingerprint)
	if err != nil {
		logger.Errorf("calling 'removeAdvancedSettingsJA4Fingerprint': %s", err.Error())
		return diag.FromErr(err)
	}
	return nil
}
