package accountprotection

import (
	"context"
	"encoding/json"
	"errors"
	"fmt"

	apr "github.com/akamai/AkamaiOPEN-edgegrid-golang/v12/pkg/accountprotection"
	"github.com/akamai/terraform-provider-akamai/v9/pkg/common/tf"
	"github.com/akamai/terraform-provider-akamai/v9/pkg/meta"
	"github.com/hashicorp/terraform-plugin-sdk/v2/diag"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/schema"
)

func dataSourceProtectedOperations() *schema.Resource {
	return &schema.Resource{
		ReadContext: readDataSourceProtectedOperations,
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
				Optional:    true,
				Description: "Identifies a specific protected operation. If not provided, all operations will be returned.",
			},
			"json": {
				Type:     schema.TypeString,
				Computed: true,
			},
		},
	}
}

func readDataSourceProtectedOperations(ctx context.Context, d *schema.ResourceData, m interface{}) diag.Diagnostics {
	meta := meta.Must(m)
	client := inst.Client(meta)
	logger := meta.Log("accountprotection", "readDataSourceProtectedOperations")
	logger.Debugf("in readDataSourceProtectedOperations")

	configID, err := tf.GetIntValue("config_id", d)
	if err != nil {
		return diag.FromErr(err)
	}

	securityPolicyID, err := tf.GetStringValue("security_policy_id", d)
	if err != nil {
		return diag.FromErr(err)
	}

	version, err := getLatestConfigVersion(ctx, configID, m)
	if err != nil {
		return diag.FromErr(err)
	}

	operationID, err := tf.GetStringValue("operation_id", d)
	if err != nil && !errors.Is(err, tf.ErrNotFound) {
		return diag.FromErr(err)
	}

	var response *apr.ListProtectedOperationsResponse

	if operationID == "" {
		request := apr.ListProtectedOperationsRequest{
			ConfigID:         int64(configID),
			Version:          int64(version),
			SecurityPolicyID: securityPolicyID,
		}
		response, err = client.ListProtectedOperations(ctx, request)
	} else {
		request := apr.GetProtectedOperationByIDRequest{
			ConfigID:         int64(configID),
			Version:          int64(version),
			SecurityPolicyID: securityPolicyID,
			OperationID:      operationID,
		}
		response, err = client.GetProtectedOperationByID(ctx, request)
	}

	if err != nil {
		logger.Errorf("calling 'ListProtectedOperations or GetProtectedOperationByID': %s", err.Error())
		return diag.FromErr(err)
	}

	jsonBody, err := json.Marshal(response)
	if err != nil {
		return diag.FromErr(err)
	}
	if err := d.Set("json", string(jsonBody)); err != nil {
		return diag.Errorf("%s: %s", tf.ErrValueSet, err.Error())
	}

	d.SetId(fmt.Sprintf("%d:%s", configID, securityPolicyID))
	return nil
}
