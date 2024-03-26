package appsec

import (
	"testing"

	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/schema"
	"github.com/tj/assert"
)

func TestUkraineGeoControlActionEqual(t *testing.T) {

	tests := map[string]struct {
		ukraineGeoControlAction interface{}
		expected                bool
	}{
		"action set": {
			ukraineGeoControlAction: "alert",
			expected:                false,
		},
		"action not set": {
			expected: true,
		},
	}

	resourceSchema := map[string]*schema.Schema{
		"ukraine_geo_control_action": {
			Type: schema.TypeString,
		},
	}

	for name, test := range tests {
		t.Run(name, func(t *testing.T) {
			resourceDataMap := map[string]interface{}{
				"ukraine_geo_control_action": test.ukraineGeoControlAction,
			}
			resourceData := schema.TestResourceDataRaw(t, resourceSchema, resourceDataMap)

			res := suppressDiffUkraineGeoControlAction("", "", "", resourceData)
			assert.Equal(t, test.expected, res)
		})
	}
}
