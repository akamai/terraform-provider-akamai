package botman

import (
	"context"
	"encoding/json"
	"fmt"
	"strconv"

	"github.com/akamai/AkamaiOPEN-edgegrid-golang/v3/pkg/botman"
	"github.com/akamai/terraform-provider-akamai/v2/pkg/akamai"
	"github.com/akamai/terraform-provider-akamai/v2/pkg/tools"

	"github.com/hashicorp/terraform-plugin-sdk/v2/diag"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/customdiff"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/schema"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/validation"
)

func resourceTransactionalEndpoint() *schema.Resource {
	return &schema.Resource{
		CreateContext: resourceTransactionalEndpointCreate,
		ReadContext:   resourceTransactionalEndpointRead,
		UpdateContext: resourceTransactionalEndpointUpdate,
		DeleteContext: resourceTransactionalEndpointDelete,
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
				Type:     schema.TypeInt,
				Required: true,
			},
			"security_policy_id": {
				Type:     schema.TypeString,
				Required: true,
			},
			"operation_id": {
				Type:     schema.TypeString,
				Required: true,
			},
			"transactional_endpoint": {
				Type:             schema.TypeString,
				Required:         true,
				ValidateDiagFunc: validation.ToDiagFunc(validation.StringIsJSON),
				DiffSuppressFunc: suppressEquivalentJSONDiffsGeneric,
			},
		},
	}
}

func resourceTransactionalEndpointCreate(ctx context.Context, d *schema.ResourceData, m interface{}) diag.Diagnostics {
	meta := akamai.Meta(m)
	client := inst.Client(meta)
	logger := meta.Log("botman", "resourceTransactionalEndpointCreateAction")
	logger.Debugf("in resourceTransactionalEndpointCreateAction")

	configID, err := tools.GetIntValue("config_id", d)
	if err != nil {
		return diag.FromErr(err)
	}

	version, err := getModifiableConfigVersion(ctx, configID, "TransactionalEndpoint", m)
	if err != nil {
		return diag.FromErr(err)
	}

	securityPolicyID, err := tools.GetStringValue("security_policy_id", d)
	if err != nil {
		return diag.FromErr(err)
	}

	operationID, err := tools.GetStringValue("operation_id", d)
	if err != nil {
		return diag.FromErr(err)
	}

	jsonPayload, err := getJSONPayload(d, "transactional_endpoint", "operationId", operationID)
	if err != nil {
		return diag.FromErr(err)
	}

	request := botman.CreateTransactionalEndpointRequest{
		ConfigID:         int64(configID),
		Version:          int64(version),
		SecurityPolicyID: securityPolicyID,
		JsonPayload:      jsonPayload,
	}

	response, err := client.CreateTransactionalEndpoint(ctx, request)
	if err != nil {
		logger.Errorf("calling 'CreateTransactionalEndpoint': %s", err.Error())
		return diag.FromErr(err)
	}

	d.SetId(fmt.Sprintf("%d:%s:%s", configID, securityPolicyID, (response)["operationId"]))

	return resourceTransactionalEndpointRead(ctx, d, m)
}

func resourceTransactionalEndpointRead(ctx context.Context, d *schema.ResourceData, m interface{}) diag.Diagnostics {
	meta := akamai.Meta(m)
	client := inst.Client(meta)
	logger := meta.Log("botman", "resourceTransactionalEndpointReadAction")
	logger.Debugf("in resourceTransactionalEndpointReadAction")

	iDParts, err := splitID(d.Id(), 3, "configID:securityPolicyID:operationID")
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

	request := botman.GetTransactionalEndpointRequest{
		ConfigID:         int64(configID),
		Version:          int64(version),
		SecurityPolicyID: securityPolicyID,
		OperationID:      operationID,
	}

	response, err := client.GetTransactionalEndpoint(ctx, request)

	if err != nil {
		logger.Errorf("calling 'GetTransactionalEndpoint': %s", err.Error())
		return diag.FromErr(err)
	}

	// Removing operationId from response to suppress diff
	delete(response, "operationId")

	if err := d.Set("config_id", configID); err != nil {
		return diag.Errorf("%s: %s", tools.ErrValueSet, err.Error())
	}
	if err := d.Set("security_policy_id", securityPolicyID); err != nil {
		return diag.Errorf("%s: %s", tools.ErrValueSet, err.Error())
	}
	if err := d.Set("operation_id", operationID); err != nil {
		return diag.Errorf("%s: %s", tools.ErrValueSet, err.Error())
	}

	jsonBody, err := json.Marshal(response)
	if err != nil {
		return diag.FromErr(err)
	}
	fields := map[string]interface{}{
		"config_id":              configID,
		"security_policy_id":     securityPolicyID,
		"operation_id":           operationID,
		"transactional_endpoint": string(jsonBody),
	}
	if err := tools.SetAttrs(d, fields); err != nil {
		return diag.Errorf("%s: %s", tools.ErrValueSet, err.Error())
	}

	return nil
}

func resourceTransactionalEndpointUpdate(ctx context.Context, d *schema.ResourceData, m interface{}) diag.Diagnostics {
	meta := akamai.Meta(m)
	client := inst.Client(meta)
	logger := meta.Log("botman", "resourceTransactionalEndpointUpdateAction")
	logger.Debugf("in resourceTransactionalEndpointUpdateAction")

	iDParts, err := splitID(d.Id(), 3, "configID:securityPolicyID:operationID")
	if err != nil {
		return diag.FromErr(err)
	}

	configID, err := strconv.Atoi(iDParts[0])
	if err != nil {
		return diag.FromErr(err)
	}

	version, err := getModifiableConfigVersion(ctx, configID, "TransactionalEndpoint", m)
	if err != nil {
		return diag.FromErr(err)
	}

	securityPolicyID := iDParts[1]

	operationID := iDParts[2]

	jsonPayload, err := getJSONPayload(d, "transactional_endpoint", "operationId", operationID)
	if err != nil {
		return diag.FromErr(err)
	}

	request := botman.UpdateTransactionalEndpointRequest{
		ConfigID:         int64(configID),
		Version:          int64(version),
		OperationID:      operationID,
		SecurityPolicyID: securityPolicyID,
		JsonPayload:      jsonPayload,
	}

	_, err = client.UpdateTransactionalEndpoint(ctx, request)
	if err != nil {
		logger.Errorf("calling 'UpdateTransactionalEndpoint': %s", err.Error())
		return diag.FromErr(err)
	}

	return resourceTransactionalEndpointRead(ctx, d, m)
}

func resourceTransactionalEndpointDelete(ctx context.Context, d *schema.ResourceData, m interface{}) diag.Diagnostics {
	meta := akamai.Meta(m)
	client := inst.Client(meta)
	logger := meta.Log("botman", "resourceTransactionalEndpointDeleteAction")
	logger.Debugf("in resourceTransactionalEndpointDeleteAction")

	iDParts, err := splitID(d.Id(), 3, "configID:securityPolicyID:operationID")
	if err != nil {
		return diag.FromErr(err)
	}

	configID, err := strconv.Atoi(iDParts[0])
	if err != nil {
		return diag.FromErr(err)
	}

	version, err := getModifiableConfigVersion(ctx, configID, "TransactionalEndpoint", m)
	if err != nil {
		return diag.FromErr(err)
	}

	securityPolicyID := iDParts[1]

	operationID := iDParts[2]

	request := botman.RemoveTransactionalEndpointRequest{
		ConfigID:         int64(configID),
		Version:          int64(version),
		SecurityPolicyID: securityPolicyID,
		OperationID:      operationID,
	}

	err = client.RemoveTransactionalEndpoint(ctx, request)
	if err != nil {
		logger.Errorf("calling 'RemoveTransactionalEndpoint': %s", err.Error())
		return diag.FromErr(err)
	}
	return nil
}
