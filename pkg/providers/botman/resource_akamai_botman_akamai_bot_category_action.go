package botman

import (
	"context"
	"encoding/json"
	"fmt"
	"regexp"
	"strconv"

	"github.com/akamai/AkamaiOPEN-edgegrid-golang/v12/pkg/botman"
	"github.com/akamai/terraform-provider-akamai/v9/pkg/common/id"
	"github.com/akamai/terraform-provider-akamai/v9/pkg/common/tf"
	"github.com/akamai/terraform-provider-akamai/v9/pkg/meta"
	"github.com/hashicorp/terraform-plugin-sdk/v2/diag"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/customdiff"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/schema"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/validation"
)

func resourceAkamaiBotCategoryAction() *schema.Resource {
	return &schema.Resource{
		CreateContext: resourceAkamaiBotCategoryActionCreate,
		ReadContext:   resourceAkamaiBotCategoryActionRead,
		UpdateContext: resourceAkamaiBotCategoryActionUpdate,
		DeleteContext: resourceAkamaiBotCategoryActionDelete,
		CustomizeDiff: customdiff.All(
			verifyConfigIDUnchanged,
			verifySecurityPolicyIDUnchanged,
			verifyCategoryIDUnchanged,
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
			"category_id": {
				Type:     schema.TypeString,
				Required: true,
			},
			"akamai_bot_category_action": {
				Type:             schema.TypeString,
				Required:         true,
				ValidateDiagFunc: validation.ToDiagFunc(validation.StringIsJSON),
				DiffSuppressFunc: suppressEquivalentJSONDiffsGeneric,
			},
		},
	}
}

var matchAkamaiBotCategoryActionExpRegex = regexp.MustCompile(`(Akamai Bot Category with id \[[-\w]+] does not exist|AkamaiBotCategoryAction with id: [-\w]+ does not exist)`)

func resourceAkamaiBotCategoryActionCreate(ctx context.Context, d *schema.ResourceData, m interface{}) diag.Diagnostics {
	meta := meta.Must(m)
	client := inst.Client(meta)
	logger := meta.Log("botman", "resourceAkamaiBotCategoryActionCreate")
	logger.Debugf("in resourceAkamaiBotCategoryActionCreate")

	configID, err := tf.GetIntValue("config_id", d)
	if err != nil {
		return diag.FromErr(err)
	}

	version, err := getModifiableConfigVersion(ctx, configID, "akamaiBotCategoryAction", m)
	if err != nil {
		return diag.FromErr(err)
	}

	securityPolicyID, err := tf.GetStringValue("security_policy_id", d)
	if err != nil {
		return diag.FromErr(err)
	}

	categoryID, err := tf.GetStringValue("category_id", d)
	if err != nil {
		return diag.FromErr(err)
	}

	jsonPayload, err := getJSONPayload(d, "akamai_bot_category_action", "categoryId", categoryID)
	if err != nil {
		return diag.FromErr(err)
	}

	request := botman.UpdateAkamaiBotCategoryActionRequest{
		ConfigID:         int64(configID),
		Version:          int64(version),
		SecurityPolicyID: securityPolicyID,
		CategoryID:       categoryID,
		JsonPayload:      jsonPayload,
	}

	_, err = client.UpdateAkamaiBotCategoryAction(ctx, request)

	if err != nil {
		errVal := err.Error()
		logger.Errorf("calling 'UpdateAkamaiBotCategoryAction': %s", errVal)
		return diag.FromErr(err)
	}

	d.SetId(fmt.Sprintf("%d:%s:%s", configID, securityPolicyID, categoryID))

	return akamaiBotCategoryActionRead(ctx, d, m, false)
}

func resourceAkamaiBotCategoryActionRead(ctx context.Context, d *schema.ResourceData, m interface{}) diag.Diagnostics {
	return akamaiBotCategoryActionRead(ctx, d, m, true)
}

func akamaiBotCategoryActionRead(ctx context.Context, d *schema.ResourceData, m interface{}, readFromCache bool) diag.Diagnostics {

	meta := meta.Must(m)
	client := inst.Client(meta)
	logger := meta.Log("botman", "resourceAkamaiBotCategoryActionRead")
	logger.Debugf("in resourceAkamaiBotCategoryActionRead")

	iDParts, err := id.Split(d.Id(), 3, "configID:securityPolicyID:categoryID")
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

	securityPolicyID := iDParts[1]

	categoryID := iDParts[2]

	request := botman.GetAkamaiBotCategoryActionRequest{
		ConfigID:         int64(configID),
		Version:          int64(version),
		SecurityPolicyID: securityPolicyID,
		CategoryID:       categoryID,
	}

	var response map[string]interface{}
	if readFromCache {
		response, err = getAkamaiBotCategoryAction(ctx, request, m)
		if err != nil {

			matched := matchAkamaiBotCategoryActionExpRegex.MatchString(err.Error())
			if matched {
				d.SetId("")
				logger.Info("AkamaiBotCategoryAction with id: " + categoryID + " does not exist, It may have been deleted outside of Terraform - resource will only be removed from state")
				return nil
			}
			logger.Errorf("calling 'GetAkamaiBotCategoryAction': %s", err.Error())
			return diag.FromErr(err)
		}
	} else {
		response, err = client.GetAkamaiBotCategoryAction(ctx, request)
		if err != nil {
			matched := matchAkamaiBotCategoryActionExpRegex.MatchString(err.Error())
			if matched {
				d.SetId("")
				logger.Info("AkamaiBotCategoryAction with id: " + categoryID + " does not exist, It may have been deleted outside of Terraform - resource will only be removed from state")
				return nil
			}
			logger.Errorf("calling 'GetAkamaiBotCategoryAction': %s", err.Error())
			return diag.FromErr(err)
		}
	}

	// Removing categoryId from response to suppress diff
	delete(response, "categoryId")

	jsonBody, err := json.Marshal(response)
	if err != nil {
		return diag.FromErr(err)
	}
	fields := map[string]interface{}{
		"config_id":                  configID,
		"security_policy_id":         securityPolicyID,
		"category_id":                categoryID,
		"akamai_bot_category_action": string(jsonBody),
	}
	if err = tf.SetAttrs(d, fields); err != nil {
		return diag.Errorf("%s: %s", tf.ErrValueSet, err.Error())
	}

	return nil
}

func resourceAkamaiBotCategoryActionUpdate(ctx context.Context, d *schema.ResourceData, m interface{}) diag.Diagnostics {
	meta := meta.Must(m)
	client := inst.Client(meta)
	logger := meta.Log("botman", "resourceAkamaiBotCategoryActionUpdate")
	logger.Debugf("in resourceAkamaiBotCategoryActionUpdate")

	iDParts, err := id.Split(d.Id(), 3, "configID:securityPolicyID:customBotCategoryID")
	if err != nil {
		return diag.FromErr(err)
	}

	configID, err := strconv.Atoi(iDParts[0])
	if err != nil {
		return diag.FromErr(err)
	}

	version, err := getModifiableConfigVersion(ctx, configID, "akamaiBotCategoryAction", m)
	if err != nil {
		return diag.FromErr(err)
	}

	securityPolicyID := iDParts[1]

	categoryID := iDParts[2]

	jsonPayload, err := getJSONPayload(d, "akamai_bot_category_action", "categoryId", categoryID)
	if err != nil {
		return diag.FromErr(err)
	}

	request := botman.UpdateAkamaiBotCategoryActionRequest{
		ConfigID:         int64(configID),
		Version:          int64(version),
		SecurityPolicyID: securityPolicyID,
		CategoryID:       categoryID,
		JsonPayload:      jsonPayload,
	}

	_, err = client.UpdateAkamaiBotCategoryAction(ctx, request)
	if err != nil {
		logger.Errorf("calling 'UpdateAkamaiBotCategoryAction': %s", err.Error())
		return diag.FromErr(err)
	}

	return akamaiBotCategoryActionRead(ctx, d, m, false)
}

func resourceAkamaiBotCategoryActionDelete(_ context.Context, _ *schema.ResourceData, m interface{}) diag.Diagnostics {
	meta := meta.Must(m)
	logger := meta.Log("botman", "resourceAkamaiBotCategoryActionDelete")
	logger.Debugf("in resourceAkamaiBotCategoryActionDelete")
	logger.Info("Botman API does not support akamai bot category action deletion - resource will only be removed from state")

	return nil
}
