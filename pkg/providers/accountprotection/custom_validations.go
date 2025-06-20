package accountprotection

import (
	"context"
	"encoding/json"
	"fmt"

	"github.com/akamai/terraform-provider-akamai/v8/pkg/common/str"
	"github.com/akamai/terraform-provider-akamai/v8/pkg/common/tf"
	"github.com/akamai/terraform-provider-akamai/v8/pkg/meta"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/schema"
)

// build request payload for create protected operations
func getCreateOperationsJSONPayload(d *schema.ResourceData, key string, idName string, idValue interface{}) (json.RawMessage, error) {
	jsonPayloadString, err := tf.GetStringValue(key, d)
	if err != nil {
		return nil, err
	}
	var payloadMap map[string]interface{}
	if err := json.Unmarshal([]byte(jsonPayloadString), &payloadMap); err != nil {
		return nil, err
	}
	payloadMap[idName] = idValue
	operationsArrayPayload := make(map[string]interface{})
	if operations, ok := payloadMap["operations"].([]interface{}); ok {
		operations = append(operations, payloadMap)
		operationsArrayPayload["operations"] = operations
	} else {
		operationsArrayPayload["operations"] = []interface{}{payloadMap}
	}

	jsonPayloadRaw, err := json.Marshal(operationsArrayPayload)
	if err != nil {
		return nil, err
	}
	return jsonPayloadRaw, nil
}

// getJSONRawMessageFromJSONString converts a JSON string from the resource data into a json.RawMessage.
func getJSONRawMessageFromJSONString(d *schema.ResourceData, key string) (json.RawMessage, error) {
	jsonPayloadString, err := tf.GetStringValue(key, d)
	if err != nil {
		return nil, err
	}
	var payloadMap map[string]interface{}
	if err := json.Unmarshal([]byte(jsonPayloadString), &payloadMap); err != nil {
		return nil, err
	}
	jsonPayloadRaw, err := json.Marshal(payloadMap)
	if err != nil {
		return nil, err
	}
	return jsonPayloadRaw, nil
}

func verifyConfigIDUnchanged(c context.Context, d *schema.ResourceDiff, m interface{}) error {
	return verifyIDUnchanged(c, d, m, "config_id")
}

func verifySecurityPolicyIDUnchanged(c context.Context, d *schema.ResourceDiff, m interface{}) error {
	return verifyIDUnchanged(c, d, m, "security_policy_id")
}

func verifyOperationIDUnchanged(c context.Context, d *schema.ResourceDiff, m interface{}) error {
	return verifyIDUnchanged(c, d, m, "operation_id")
}

func verifyIDUnchanged(_ context.Context, d *schema.ResourceDiff, m interface{}, key string) error {
	meta := meta.Must(m)
	logger := meta.Log("accountprotection", "verifyIDUnchanged")

	if d.Id() == "" {
		return nil
	}
	_, exist := d.GetOkExists(key)

	if exist && !d.HasChange(key) {
		return nil
	}

	oldID, newID := d.GetChange(key)
	oldValue := str.From(oldID)
	newValue := str.From(newID)
	logger.Errorf("%s value %s specified in configuration differs from resource ID's value %s", key, newValue, oldValue)
	return fmt.Errorf("%s value %s specified in configuration differs from resource ID's value %s", key, newValue, oldValue)
}
