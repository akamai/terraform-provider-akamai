package appsec

import (
	"context"
	"errors"
	"fmt"
	"strconv"
	"strings"

	"github.com/akamai/AkamaiOPEN-edgegrid-golang/v8/pkg/appsec"
	"github.com/akamai/terraform-provider-akamai/v5/pkg/common/tf"
	"github.com/akamai/terraform-provider-akamai/v5/pkg/meta"
	"github.com/hashicorp/terraform-plugin-sdk/v2/diag"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/customdiff"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/schema"
)

// appsec v1
//
// https://techdocs.akamai.com/application-security/reference/api
func resourceAPIRequestConstraints() *schema.Resource {
	return &schema.Resource{
		CreateContext: resourceAPIRequestConstraintsCreate,
		ReadContext:   resourceAPIRequestConstraintsRead,
		UpdateContext: resourceAPIRequestConstraintsUpdate,
		DeleteContext: resourceAPIRequestConstraintsDelete,
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
			"api_endpoint_id": {
				Type:        schema.TypeInt,
				Optional:    true,
				Description: "Unique identifier of the API endpoint to which the constraint will be assigned",
			},
			"action": {
				Type:             schema.TypeString,
				Required:         true,
				ValidateDiagFunc: ValidateActions,
				Description:      "Action to be taken when the API request constraint is triggered",
			},
		},
	}
}

func resourceAPIRequestConstraintsCreate(ctx context.Context, d *schema.ResourceData, m interface{}) diag.Diagnostics {
	meta := meta.Must(m)
	client := inst.Client(meta)
	logger := meta.Log("APPSEC", "resourceAPIRequestConstraintsCreate")
	logger.Debugf("in resourceAPIRequestConstraintsCreate")

	configID, err := tf.GetIntValue("config_id", d)
	if err != nil {
		return diag.FromErr(err)
	}
	version, err := getModifiableConfigVersion(ctx, configID, "apirequestconstraints", m)
	if err != nil {
		return diag.FromErr(err)
	}
	policyID, err := tf.GetStringValue("security_policy_id", d)
	if err != nil {
		return diag.FromErr(err)
	}
	apiEndpointID, err := tf.GetIntValue("api_endpoint_id", d)
	if err != nil && !errors.Is(err, tf.ErrNotFound) {
		return diag.FromErr(err)
	}
	action, err := tf.GetStringValue("action", d)
	if err != nil {
		return diag.FromErr(err)
	}

	createAPIRequestConstraints := appsec.UpdateApiRequestConstraintsRequest{
		ConfigID: configID,
		Version:  version,
		PolicyID: policyID,
		ApiID:    apiEndpointID,
		Action:   action,
	}

	_, err = client.UpdateApiRequestConstraints(ctx, createAPIRequestConstraints)
	if err != nil {
		logger.Errorf("calling 'createAPIRequestConstraints': %s", err.Error())
		return diag.FromErr(err)
	}

	if apiEndpointID != 0 {
		d.SetId(fmt.Sprintf("%d:%s:%d", createAPIRequestConstraints.ConfigID, createAPIRequestConstraints.PolicyID, createAPIRequestConstraints.ApiID))
	} else {
		d.SetId(fmt.Sprintf("%d:%s", createAPIRequestConstraints.ConfigID, createAPIRequestConstraints.PolicyID))
	}

	return resourceAPIRequestConstraintsRead(ctx, d, m)
}

func resourceAPIRequestConstraintsRead(ctx context.Context, d *schema.ResourceData, m interface{}) diag.Diagnostics {
	meta := meta.Must(m)
	client := inst.Client(meta)
	logger := meta.Log("APPSEC", "resourceAPIRequestConstraintsRead")
	logger.Debugf("in resourceAPIRequestConstraintsRead")

	s := strings.Split(d.Id(), ":")

	configID, errconv := strconv.Atoi(s[0])
	if errconv != nil {
		return diag.FromErr(errconv)
	}
	version, err := getLatestConfigVersion(ctx, configID, m)
	if err != nil {
		return diag.FromErr(err)
	}
	policyID := s[1]

	apiID := 0
	if len(s) > 2 {
		apiID, errconv = strconv.Atoi(s[2])
		if errconv != nil {
			return diag.FromErr(errconv)
		}
	}

	response, err := client.GetApiRequestConstraints(ctx, appsec.GetApiRequestConstraintsRequest{
		ConfigID: configID,
		Version:  version,
		PolicyID: policyID,
		ApiID:    apiID,
	})
	if err != nil {
		logger.Errorf("calling 'getAPIRequestConstraints': %s", err.Error())
		return diag.FromErr(err)
	}

	if err := d.Set("config_id", configID); err != nil {
		return diag.Errorf("%s: %s", tf.ErrValueSet, err.Error())
	}

	if err := d.Set("security_policy_id", policyID); err != nil {
		return diag.Errorf("%s: %s", tf.ErrValueSet, err.Error())
	}

	if err := d.Set("api_endpoint_id", apiID); err != nil {
		return diag.Errorf("%s: %s", tf.ErrValueSet, err.Error())
	}

	if apiID != 0 {
		if len(response.APIEndpoints) > 0 {
			for _, val := range response.APIEndpoints {
				if val.ID == apiID {
					if err := d.Set("action", val.Action); err != nil {
						return diag.Errorf("%s: %s", tf.ErrValueSet, err.Error())
					}
				}
			}
		}
	}
	return nil
}

func resourceAPIRequestConstraintsUpdate(ctx context.Context, d *schema.ResourceData, m interface{}) diag.Diagnostics {
	meta := meta.Must(m)
	client := inst.Client(meta)
	logger := meta.Log("APPSEC", "resourceAPIRequestConstraintsUpdate")
	logger.Debugf("in resourceAPIRequestConstraintsUpdate")

	s := strings.Split(d.Id(), ":")

	configID, errconv := strconv.Atoi(s[0])
	if errconv != nil {
		return diag.FromErr(errconv)
	}
	version, err := getModifiableConfigVersion(ctx, configID, "apirequestconstraints", m)
	if err != nil {
		return diag.FromErr(err)
	}
	policyID := s[1]

	apiID := 0
	if len(s) > 2 {
		apiID, errconv = strconv.Atoi(s[2])
		if errconv != nil {
			return diag.FromErr(errconv)
		}
	}
	action, err := tf.GetStringValue("action", d)
	if err != nil && !errors.Is(err, tf.ErrNotFound) {
		return diag.FromErr(err)
	}

	_, err = client.UpdateApiRequestConstraints(ctx, appsec.UpdateApiRequestConstraintsRequest{
		ConfigID: configID,
		Version:  version,
		PolicyID: policyID,
		ApiID:    apiID,
		Action:   action,
	})
	if err != nil {
		logger.Errorf("calling 'updateAPIRequestConstraints': %s", err.Error())
		return diag.FromErr(err)
	}

	return resourceAPIRequestConstraintsRead(ctx, d, m)
}

func resourceAPIRequestConstraintsDelete(ctx context.Context, d *schema.ResourceData, m interface{}) diag.Diagnostics {
	meta := meta.Must(m)
	client := inst.Client(meta)
	logger := meta.Log("APPSEC", "resourceAPIRequestConstraintsDelete")
	logger.Debugf("in resourceAPIRequestConstraintsDelete")

	s := strings.Split(d.Id(), ":")

	configID, errconv := strconv.Atoi(s[0])
	if errconv != nil {
		return diag.FromErr(errconv)
	}
	version, err := getModifiableConfigVersion(ctx, configID, "apirequestconstraints", m)
	if err != nil {
		return diag.FromErr(err)
	}
	policyID := s[1]

	apiID := 0
	if len(s) > 2 {
		apiID, errconv = strconv.Atoi(s[2])
		if errconv != nil {
			return diag.FromErr(errconv)
		}
	}

	removeAPIRequestConstraints := appsec.RemoveApiRequestConstraintsRequest{
		ConfigID: configID,
		Version:  version,
		PolicyID: policyID,
		ApiID:    apiID,
	}

	if removeAPIRequestConstraints.ApiID == 0 {
		_, err := client.UpdateAPIConstraintsProtection(ctx, appsec.UpdateAPIConstraintsProtectionRequest{
			ConfigID:            configID,
			Version:             version,
			PolicyID:            policyID,
			ApplyAPIConstraints: false,
		})
		if err != nil {
			logger.Errorf("calling 'UpdateAPIConstraintsProtection': %s", err.Error())
			return diag.FromErr(err)
		}
	} else {
		removeAPIRequestConstraints.Action = "none"
		_, err := client.RemoveApiRequestConstraints(ctx, removeAPIRequestConstraints)
		if err != nil {
			logger.Errorf("calling 'removeApiRequestConstraints': %s", err.Error())
			return diag.FromErr(err)
		}
	}
	return nil
}
