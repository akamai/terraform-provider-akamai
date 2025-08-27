package accountprotection

import (
	"context"
	"encoding/json"
	"fmt"
	"strconv"

	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/validation"

	apr "github.com/akamai/AkamaiOPEN-edgegrid-golang/v12/pkg/accountprotection"
	"github.com/akamai/terraform-provider-akamai/v8/pkg/common/id"
	"github.com/akamai/terraform-provider-akamai/v8/pkg/common/tf"
	"github.com/akamai/terraform-provider-akamai/v8/pkg/meta"
	"github.com/hashicorp/terraform-plugin-sdk/v2/diag"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/customdiff"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/schema"
)

func resourceGeneralSettings() *schema.Resource {
	return &schema.Resource{
		CreateContext: createResourceGeneralSettings,
		ReadContext:   readResourceGeneralSettings,
		UpdateContext: updateResourceGeneralSettings,
		DeleteContext: deleteResourceGeneralSettings,
		CustomizeDiff: customdiff.All(
			verifyConfigIDUnchanged,
			verifySecurityPolicyIDUnchanged,
		),
		Importer: &schema.ResourceImporter{
			StateContext: schema.ImportStatePassthroughContext,
		},
		Schema: map[string]*schema.Schema{
			"config_id": {
				Type:        schema.TypeInt,
				Required:    true,
				Description: "Identifies a security configuration.",
			},
			"security_policy_id": {
				Type:        schema.TypeString,
				Required:    true,
				Description: "Identifies a security policy.",
			},
			"general_settings": {
				Type:             schema.TypeString,
				Required:         true,
				ValidateDiagFunc: validation.ToDiagFunc(validation.StringIsJSON),
				DiffSuppressFunc: suppressEquivalentJSONDiffsGeneric,
			},
		},
	}
}

func createResourceGeneralSettings(ctx context.Context, d *schema.ResourceData, m interface{}) diag.Diagnostics {
	meta := meta.Must(m)
	client := inst.Client(meta)
	logger := meta.Log("accountprotection", "createResourceGeneralSettings")
	logger.Debugf("in createResourceGeneralSettings")

	configID, err := tf.GetIntValue("config_id", d)
	if err != nil {
		return diag.FromErr(err)
	}

	version, err := getModifiableConfigVersion(ctx, configID, "aprGeneralSettings", m)
	if err != nil {
		return diag.FromErr(err)
	}

	securityPolicyID, err := tf.GetStringValue("security_policy_id", d)
	if err != nil {
		return diag.FromErr(err)
	}

	jsonPayloadString, err := tf.GetStringValue("general_settings", d)
	if err != nil {
		return diag.FromErr(err)
	}

	request := apr.UpsertGeneralSettingsRequest{
		ConfigID:         int64(configID),
		Version:          int64(version),
		SecurityPolicyID: securityPolicyID,
		JsonPayload:      json.RawMessage(jsonPayloadString),
	}

	_, err = client.UpsertGeneralSettings(ctx, request)
	if err != nil {
		logger.Errorf("calling UpsertGeneralSettingsRequest 'request': %s", err.Error())
		return diag.FromErr(err)
	}

	d.SetId(fmt.Sprintf("%d:%s", configID, securityPolicyID))

	return readResourceGeneralSettings(ctx, d, m)
}

func readResourceGeneralSettings(ctx context.Context, d *schema.ResourceData, m interface{}) diag.Diagnostics {
	meta := meta.Must(m)
	client := inst.Client(meta)
	logger := meta.Log("accountprotection", "readResourceGeneralSettings")
	logger.Debugf("in readResourceGeneralSettings")

	iDParts, err := id.Split(d.Id(), 2, "configID:securityPolicyID")
	if err != nil {
		logger.Errorf("invalid ID format for resource %q: expected format 'configID:securityPolicyID'", d.Id())
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

	securityPolicyID := iDParts[1]

	request := apr.GetGeneralSettingsRequest{
		ConfigID:         int64(configID),
		Version:          int64(version),
		SecurityPolicyID: securityPolicyID,
	}

	response, err := client.GetGeneralSettings(ctx, request)
	if err != nil {
		logger.Errorf("calling GetGeneralSettings 'request': %s", err.Error())
		return diag.FromErr(err)
	}

	delete(response, "metadata")

	jsonBody, err := json.Marshal(response)
	if err != nil {
		return diag.FromErr(err)
	}
	fields := map[string]interface{}{
		"config_id":          configID,
		"security_policy_id": securityPolicyID,
		"general_settings":   string(jsonBody),
	}
	if err = tf.SetAttrs(d, fields); err != nil {
		return diag.Errorf("%s: %s", tf.ErrValueSet, err.Error())
	}

	return nil
}

func updateResourceGeneralSettings(ctx context.Context, d *schema.ResourceData, m interface{}) diag.Diagnostics {
	meta := meta.Must(m)
	client := inst.Client(meta)
	logger := meta.Log("accountprotection", "updateResourceGeneralSettings")
	logger.Debugf("in updateResourceGeneralSettings")

	iDParts, err := id.Split(d.Id(), 2, "configID:securityPolicyID")
	if err != nil {
		return diag.FromErr(err)
	}

	configID, err := strconv.Atoi(iDParts[0])
	if err != nil {
		return diag.FromErr(err)
	}

	version, err := getModifiableConfigVersion(ctx, configID, "aprGeneralSettings", m)
	if err != nil {
		return diag.FromErr(err)
	}

	securityPolicyID := iDParts[1]

	jsonPayloadString, err := tf.GetStringValue("general_settings", d)
	if err != nil {
		return diag.FromErr(err)
	}

	request := apr.UpsertGeneralSettingsRequest{
		ConfigID:         int64(configID),
		Version:          int64(version),
		SecurityPolicyID: securityPolicyID,
		JsonPayload:      json.RawMessage(jsonPayloadString),
	}

	_, err = client.UpsertGeneralSettings(ctx, request)
	if err != nil {
		logger.Errorf("calling UpsertGeneralSettings 'request': %s", err.Error())
		return diag.FromErr(err)
	}

	return readResourceGeneralSettings(ctx, d, m)
}

func deleteResourceGeneralSettings(_ context.Context, d *schema.ResourceData, m interface{}) diag.Diagnostics {
	meta := meta.Must(m)
	logger := meta.Log("accountprotection", "deleteResourceGeneralSettings")
	logger.Debugf("in deleteResourceGeneralSettings")
	logger.Info("APR API does not support general-settings deletion - resource will only be removed from state")

	d.SetId("")
	return nil
}
