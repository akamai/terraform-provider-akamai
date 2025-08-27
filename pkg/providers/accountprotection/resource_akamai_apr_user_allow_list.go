package accountprotection

import (
	"context"
	"encoding/json"
	"strconv"

	apr "github.com/akamai/AkamaiOPEN-edgegrid-golang/v12/pkg/accountprotection"
	"github.com/akamai/terraform-provider-akamai/v9/pkg/common/tf"
	"github.com/akamai/terraform-provider-akamai/v9/pkg/meta"
	"github.com/hashicorp/terraform-plugin-sdk/v2/diag"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/customdiff"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/schema"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/validation"
)

func resourceUserAllowList() *schema.Resource {
	return &schema.Resource{
		CreateContext: createResourceUserAllowList,
		ReadContext:   readResourceUserAllowList,
		UpdateContext: updateResourceUserAllowList,
		DeleteContext: deleteResourceUserAllowList,
		CustomizeDiff: customdiff.All(
			verifyConfigIDUnchanged,
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
			"user_allow_list": {
				Type:             schema.TypeString,
				Optional:         true,
				ValidateDiagFunc: validation.ToDiagFunc(validation.StringIsJSON),
				DiffSuppressFunc: suppressEquivalentJSONDiffsGeneric,
			},
		},
	}
}

func createResourceUserAllowList(ctx context.Context, d *schema.ResourceData, m interface{}) diag.Diagnostics {
	meta := meta.Must(m)
	client := inst.Client(meta)
	logger := meta.Log("accountprotection", "createResourceUserAllowList")
	logger.Debugf("in createResourceUserAllowList")

	configID, err := tf.GetIntValue("config_id", d)
	if err != nil {
		return diag.FromErr(err)
	}

	version, err := getModifiableConfigVersion(ctx, configID, "aprUserAllowList", m)
	if err != nil {
		return diag.FromErr(err)
	}

	jsonPayloadString, err := tf.GetStringValue("user_allow_list", d)
	if err != nil {
		return diag.FromErr(err)
	}

	request := apr.UpsertUserAllowListIDRequest{
		ConfigID:    int64(configID),
		Version:     int64(version),
		JsonPayload: json.RawMessage(jsonPayloadString),
	}

	_, err = client.UpsertUserAllowListID(ctx, request)
	if err != nil {
		logger.Errorf("calling 'UpsertUserAllowListID': %s", err.Error())
		return diag.FromErr(err)
	}

	d.SetId(strconv.Itoa(configID))

	return readResourceUserAllowList(ctx, d, m)
}

func readResourceUserAllowList(ctx context.Context, d *schema.ResourceData, m interface{}) diag.Diagnostics {
	meta := meta.Must(m)
	client := inst.Client(meta)
	logger := meta.Log("accountprotection", "readResourceUserAllowList")
	logger.Debugf("in readResourceUserAllowList")

	configID, err := strconv.Atoi(d.Id())
	if err != nil {
		return diag.FromErr(err)
	}

	version, err := getLatestConfigVersion(ctx, configID, m)
	if err != nil {
		return diag.FromErr(err)
	}

	request := apr.GetUserAllowListIDRequest{
		ConfigID: int64(configID),
		Version:  int64(version),
	}

	response, err := client.GetUserAllowListID(ctx, request)

	if err != nil {
		logger.Errorf("calling 'GetUserAllowListID': %s", err.Error())
		return diag.FromErr(err)
	}

	if response != nil {
		delete(response, "metadata")
		jsonBody, err := json.Marshal(response)
		if err != nil {
			return diag.FromErr(err)
		}
		fields := map[string]interface{}{
			"config_id":       configID,
			"user_allow_list": string(jsonBody),
		}
		if err = tf.SetAttrs(d, fields); err != nil {
			return diag.Errorf("%s: %s", tf.ErrValueSet, err.Error())
		}
	} else {
		fields := map[string]interface{}{
			"config_id": configID,
		}
		if err = tf.SetAttrs(d, fields); err != nil {
			return diag.Errorf("%s: %s", tf.ErrValueSet, err.Error())
		}
	}

	return nil
}

func updateResourceUserAllowList(ctx context.Context, d *schema.ResourceData, m interface{}) diag.Diagnostics {
	meta := meta.Must(m)
	client := inst.Client(meta)
	logger := meta.Log("accountprotection", "updateResourceUserAllowList")
	logger.Debugf("in updateResourceUserAllowList")

	configID, err := strconv.Atoi(d.Id())
	if err != nil {
		return diag.FromErr(err)
	}

	version, err := getModifiableConfigVersion(ctx, configID, "userAllowListId", m)
	if err != nil {
		return diag.FromErr(err)
	}

	jsonPayloadString, err := tf.GetStringValue("user_allow_list", d)
	if err != nil {
		return diag.FromErr(err)
	}

	request := apr.UpsertUserAllowListIDRequest{
		ConfigID:    int64(configID),
		Version:     int64(version),
		JsonPayload: json.RawMessage(jsonPayloadString),
	}

	_, err = client.UpsertUserAllowListID(ctx, request)
	if err != nil {
		logger.Errorf("calling 'UpsertUserAllowListID': %s", err.Error())
		return diag.FromErr(err)
	}

	return readResourceUserAllowList(ctx, d, m)
}

func deleteResourceUserAllowList(ctx context.Context, d *schema.ResourceData, m interface{}) diag.Diagnostics {
	meta := meta.Must(m)
	client := inst.Client(meta)
	logger := meta.Log("accountprotection", "deleteResourceUserAllowList")
	logger.Debugf("in accountprotection deleteResourceUserAllowList")

	configID, err := strconv.Atoi(d.Id())
	if err != nil {
		return diag.FromErr(err)
	}

	version, err := getModifiableConfigVersion(ctx, configID, "userAllowListId", m)
	if err != nil {
		return diag.FromErr(err)
	}

	request := apr.DeleteUserAllowListIDRequest{
		ConfigID: int64(configID),
		Version:  int64(version),
	}

	err = client.DeleteUserAllowListID(ctx, request)
	if err != nil {
		logger.Errorf("calling 'DeleteUserAllowListID': %s", err.Error())
		return diag.FromErr(err)
	}
	d.SetId("")
	return nil
}
