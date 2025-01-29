package botman

import (
	"context"
	"encoding/json"
	"fmt"

	"github.com/akamai/terraform-provider-akamai/v7/pkg/common/str"
	"github.com/akamai/terraform-provider-akamai/v7/pkg/common/tf"
	"github.com/akamai/terraform-provider-akamai/v7/pkg/meta"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/schema"
)

// getJSONPayload adds ID to JSON payload for update operations
func getJSONPayload(d *schema.ResourceData, key string, idName string, idValue interface{}) (json.RawMessage, error) {
	jsonPayloadString, err := tf.GetStringValue(key, d)
	if err != nil {
		return nil, err
	}
	var payloadMap map[string]interface{}
	if err := json.Unmarshal([]byte(jsonPayloadString), &payloadMap); err != nil {
		return nil, err
	}
	payloadMap[idName] = idValue
	jsonPayloadRaw, err := json.Marshal(payloadMap)
	if err != nil {
		return nil, err
	}
	return jsonPayloadRaw, err
}

func verifyConfigIDUnchanged(c context.Context, d *schema.ResourceDiff, m interface{}) error {
	return verifyIDUnchanged(c, d, m, "config_id")
}

func verifySecurityPolicyIDUnchanged(c context.Context, d *schema.ResourceDiff, m interface{}) error {
	return verifyIDUnchanged(c, d, m, "security_policy_id")
}

func verifyCategoryIDUnchanged(c context.Context, d *schema.ResourceDiff, m interface{}) error {
	return verifyIDUnchanged(c, d, m, "category_id")
}

func verifyDetectionIDUnchanged(c context.Context, d *schema.ResourceDiff, m interface{}) error {
	return verifyIDUnchanged(c, d, m, "detection_id")
}

func verifyBotIDUnchanged(c context.Context, d *schema.ResourceDiff, m interface{}) error {
	return verifyIDUnchanged(c, d, m, "bot_id")
}

func verifyOperationIDUnchanged(c context.Context, d *schema.ResourceDiff, m interface{}) error {
	return verifyIDUnchanged(c, d, m, "operation_id")
}

func verifyIDUnchanged(_ context.Context, d *schema.ResourceDiff, m interface{}, key string) error {
	meta := meta.Must(m)
	logger := meta.Log("botman", "verifyIDUnchanged")

	if d.Id() == "" {
		return nil
	}
	if d.GetOkExists(key); !d.HasChange(key) {
		return nil
	}

	oldID, newID := d.GetChange(key)
	oldValue := str.From(oldID)
	newValue := str.From(newID)
	logger.Errorf("%s value %s specified in configuration differs from resource ID's value %s", key, newID, oldValue)
	return fmt.Errorf("%s value %s specified in configuration differs from resource ID's value %s", key, newValue, oldValue)
}
