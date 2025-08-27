package accountprotection

import (
	"context"
	"encoding/json"
	"fmt"
	"strconv"

	apr "github.com/akamai/AkamaiOPEN-edgegrid-golang/v12/pkg/accountprotection"
	"github.com/akamai/terraform-provider-akamai/v9/pkg/common/id"
	"github.com/akamai/terraform-provider-akamai/v9/pkg/common/tf"
	"github.com/akamai/terraform-provider-akamai/v9/pkg/meta"
	"github.com/hashicorp/terraform-plugin-sdk/v2/diag"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/customdiff"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/schema"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/validation"
)

func resourceProtectedOperations() *schema.Resource {
	return &schema.Resource{
		CreateContext: createResourceProtectedOperations,
		ReadContext:   readResourceProtectedOperations,
		UpdateContext: updateResourceProtectedOperation,
		DeleteContext: deleteResourceProtectedOperation,
		CustomizeDiff: customdiff.All(
			verifyConfigIDUnchanged,
			verifySecurityPolicyIDUnchanged,
			verifyOperationIDUnchanged,
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
			"operation_id": {
				Type:        schema.TypeString,
				Required:    true,
				Description: "Identifies a protected operation",
			},
			"protected_operation": {
				Type:             schema.TypeString,
				Required:         true,
				ValidateDiagFunc: validation.ToDiagFunc(validation.StringIsJSON),
				DiffSuppressFunc: suppressEquivalentJSONDiffsGeneric,
			},
		},
	}
}

// createResourceProtectedOperations creates a new protected operation in the specified security policy.
func createResourceProtectedOperations(ctx context.Context, d *schema.ResourceData, m interface{}) diag.Diagnostics {
	meta := meta.Must(m)
	client := inst.Client(meta)
	logger := meta.Log("accountprotection", "createResourceProtectedOperations")
	logger.Debugf("in createResourceProtectedOperations")

	configID, err := tf.GetIntValue("config_id", d)
	if err != nil {
		return diag.FromErr(err)
	}

	version, err := getModifiableConfigVersion(ctx, configID, "ProtectedOperation", m)
	if err != nil {
		return diag.FromErr(err)
	}

	securityPolicyID, err := tf.GetStringValue("security_policy_id", d)
	if err != nil {
		return diag.FromErr(err)
	}

	operationID, err := tf.GetStringValue("operation_id", d)
	if err != nil {
		return diag.FromErr(err)
	}

	jsonPayload, err := getCreateOperationsJSONPayload(d, "protected_operation", "operationId", operationID)
	if err != nil {
		return diag.FromErr(err)
	}

	request := apr.CreateProtectedOperationsRequest{
		ConfigID:         int64(configID),
		Version:          int64(version),
		SecurityPolicyID: securityPolicyID,
		JsonPayload:      jsonPayload,
	}

	response, err := client.CreateProtectedOperations(ctx, request)
	if err != nil {
		logger.Errorf("calling 'CreateProtectedOperations': %s", err.Error())
		return diag.FromErr(err)
	}

	d.SetId(fmt.Sprintf("%d:%s:%s", configID, securityPolicyID, response.Operations[0]["operationId"].(string)))

	return readProtectedOperations(ctx, d, m, false)
}

// readResourceProtectedOperations reads the protected operations for a given security policy.
func readResourceProtectedOperations(ctx context.Context, d *schema.ResourceData, m interface{}) diag.Diagnostics {
	return readProtectedOperations(ctx, d, m, true)
}

// readProtectedOperations reads the protected operations for a given security policy either from cache or directly from the API.
func readProtectedOperations(ctx context.Context, d *schema.ResourceData, m interface{}, readFromCache bool) diag.Diagnostics {
	meta := meta.Must(m)
	client := inst.Client(meta)
	logger := meta.Log("accountprotection", "readProtectedOperations")
	logger.Debugf("in readProtectedOperations")

	iDParts, err := id.Split(d.Id(), 3, "configID:securityPolicyID:operationID")
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

	operationID := iDParts[2]

	request := apr.GetProtectedOperationByIDRequest{
		ConfigID:         int64(configID),
		Version:          int64(version),
		SecurityPolicyID: securityPolicyID,
		OperationID:      operationID,
	}

	var response map[string]interface{}
	if readFromCache {
		response, err = getProtectedOperations(ctx, request, m)
		if err != nil {
			return diag.FromErr(err)
		}
	} else {
		getByIDResponse, err := client.GetProtectedOperationByID(ctx, request)
		if err != nil {
			logger.Errorf("calling 'GetProtectedOperationByID': %s", err.Error())
			return diag.FromErr(err)
		}
		if getByIDResponse != nil && len(getByIDResponse.Operations) > 0 {
			response = getByIDResponse.Operations[0]
		}
	}

	delete(response, "operationId")
	delete(response, "protectedOperationLink")
	delete(response, "metadata")

	if err := d.Set("config_id", configID); err != nil {
		return diag.Errorf("%s: %s", tf.ErrValueSet, err.Error())
	}
	if err := d.Set("security_policy_id", securityPolicyID); err != nil {
		return diag.Errorf("%s: %s", tf.ErrValueSet, err.Error())
	}
	if err := d.Set("operation_id", operationID); err != nil {
		return diag.Errorf("%s: %s", tf.ErrValueSet, err.Error())
	}

	jsonBody, err := json.Marshal(response)
	if err != nil {
		return diag.FromErr(err)
	}
	fields := map[string]interface{}{
		"config_id":           configID,
		"security_policy_id":  securityPolicyID,
		"operation_id":        operationID,
		"protected_operation": string(jsonBody),
	}
	if err := tf.SetAttrs(d, fields); err != nil {
		return diag.Errorf("%s: %s", tf.ErrValueSet, err.Error())
	}

	return nil
}

// updateResourceProtectedOperation updates an existing protected operation in the specified security policy.
func updateResourceProtectedOperation(ctx context.Context, d *schema.ResourceData, m interface{}) diag.Diagnostics {
	meta := meta.Must(m)
	client := inst.Client(meta)
	logger := meta.Log("accountprotection", "updateResourceProtectedOperation")
	logger.Debugf("in updateResourceProtectedOperation")

	iDParts, err := id.Split(d.Id(), 3, "configID:securityPolicyID:operationID")
	if err != nil {
		return diag.FromErr(err)
	}

	configID, err := strconv.Atoi(iDParts[0])
	if err != nil {
		return diag.FromErr(err)
	}

	version, err := getModifiableConfigVersion(ctx, configID, "ProtectedOperation", m)
	if err != nil {
		return diag.FromErr(err)
	}

	securityPolicyID := iDParts[1]

	operationID := iDParts[2]

	jsonPayload, err := getJSONRawMessageFromJSONString(d, "protected_operation")
	if err != nil {
		return diag.FromErr(err)
	}

	request := apr.UpdateProtectedOperationRequest{
		ConfigID:         int64(configID),
		Version:          int64(version),
		OperationID:      operationID,
		SecurityPolicyID: securityPolicyID,
		JsonPayload:      jsonPayload,
	}

	_, err = client.UpdateProtectedOperation(ctx, request)
	if err != nil {
		logger.Errorf("calling 'UpdateProtectedOperation': %s", err.Error())
		return diag.FromErr(err)
	}

	return readProtectedOperations(ctx, d, m, false)
}

// deleteResourceProtectedOperation deletes a protected operation from the specified security policy.
func deleteResourceProtectedOperation(ctx context.Context, d *schema.ResourceData, m interface{}) diag.Diagnostics {
	meta := meta.Must(m)
	client := inst.Client(meta)
	logger := meta.Log("accountprotection", "deleteResourceProtectedOperation")
	logger.Debugf("in deleteResourceProtectedOperation")

	iDParts, err := id.Split(d.Id(), 3, "configID:securityPolicyID:operationID")
	if err != nil {
		return diag.FromErr(err)
	}

	configID, err := strconv.Atoi(iDParts[0])
	if err != nil {
		return diag.FromErr(err)
	}

	version, err := getModifiableConfigVersion(ctx, configID, "ProtectedOperation", m)
	if err != nil {
		return diag.FromErr(err)
	}

	securityPolicyID := iDParts[1]

	operationID := iDParts[2]

	request := apr.RemoveProtectedOperationRequest{
		ConfigID:         int64(configID),
		Version:          int64(version),
		SecurityPolicyID: securityPolicyID,
		OperationID:      operationID,
	}

	err = client.RemoveProtectedOperation(ctx, request)
	if err != nil {
		logger.Errorf("calling 'RemoveProtectedOperation': %s", err.Error())
		return diag.FromErr(err)
	}
	d.SetId("")
	return nil
}
