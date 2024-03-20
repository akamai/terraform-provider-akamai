package cloudlets

import (
	"regexp"
	"testing"

	"github.com/akamai/AkamaiOPEN-edgegrid-golang/v8/pkg/cloudlets"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/schema"
	"github.com/stretchr/testify/require"
	"github.com/tj/assert"
)

func TestParseObjectMatchValue(t *testing.T) {
	dummySchemaSetFunc := func(v interface{}) int {
		return 1
	}

	tests := map[string]struct {
		criteria       map[string]interface{}
		handler        objectMatchValueHandler
		expectedOutput interface{}
		expectedError  *regexp.Regexp
	}{
		"no objectMatchValue": {
			criteria: map[string]interface{}{
				"name": "rule",
				"type": "vpMatchRule",
			},
			expectedOutput: nil,
		},
		"empty objectMatchValue": {
			criteria: map[string]interface{}{
				"name":               "rule",
				"type":               "vpMatchRule",
				"object_match_value": &schema.Set{},
			},
			expectedOutput: nil,
		},
		"invalid objectMatchValue - not an object": {
			criteria: map[string]interface{}{
				"name":               "rule",
				"type":               "vpMatchRule",
				"object_match_value": schema.NewSet(dummySchemaSetFunc, []interface{}{"foo"}),
			},
			expectedError: regexp.MustCompile("'object_match_value' should be an object"),
		},
		"valid objectMatchValue - simple": {
			criteria: map[string]interface{}{
				"name": "rule",
				"type": "vpMatchRule",
				"object_match_value": schema.NewSet(dummySchemaSetFunc, []interface{}{
					map[string]interface{}{
						"type":  "simple",
						"value": []interface{}{"foo"},
					},
				}),
			},
			handler: getObjectMatchValueObjectOrSimpleOrRange,
			expectedOutput: &cloudlets.ObjectMatchValueSimple{
				Type:  "simple",
				Value: []string{"foo"},
			},
		},
		"valid objectMatchValue - object": {
			criteria: map[string]interface{}{
				"name": "rule",
				"type": "vpMatchRule",
				"object_match_value": schema.NewSet(dummySchemaSetFunc, []interface{}{
					map[string]interface{}{
						"type":                "object",
						"name":                "omv",
						"name_case_sensitive": true,
						"name_has_wildcard":   true,
						"options": schema.NewSet(dummySchemaSetFunc, []interface{}{
							map[string]interface{}{
								"value": []interface{}{"foo"},
							},
						}),
					},
				}),
			},
			handler: getObjectMatchValueObjectOrSimpleOrRange,
			expectedOutput: &cloudlets.ObjectMatchValueObject{
				Name:              "omv",
				Type:              "object",
				NameCaseSensitive: true,
				NameHasWildcard:   true,
				Options: &cloudlets.Options{
					Value: []string{"foo"},
				},
			},
		},
		"valid objectMatchValue - range": {
			criteria: map[string]interface{}{
				"name": "rule",
				"type": "vpMatchRule",
				"object_match_value": schema.NewSet(dummySchemaSetFunc, []interface{}{
					map[string]interface{}{
						"type":  "range",
						"value": []interface{}{"1", "50"},
					},
				}),
			},
			handler: getObjectMatchValueObjectOrSimpleOrRange,
			expectedOutput: &cloudlets.ObjectMatchValueRange{
				Type:  "range",
				Value: []int64{1, 50},
			},
		},
		"invalid objectMatchValue - range": {
			criteria: map[string]interface{}{
				"name": "rule",
				"type": "vpMatchRule",
				"object_match_value": schema.NewSet(dummySchemaSetFunc, []interface{}{
					map[string]interface{}{
						"type":  "range",
						"value": []interface{}{"1", "50"},
					},
				}),
			},
			handler:       getObjectMatchValueObjectOrSimple,
			expectedError: regexp.MustCompile("Must be one of: 'simple' or 'object'"),
		},
	}

	for name, test := range tests {
		t.Run(name, func(t *testing.T) {
			out, err := parseObjectMatchValue(test.criteria, test.handler)

			if test.expectedError != nil {
				require.Error(t, err)
				assert.True(t, test.expectedError.MatchString(err.Error()))
				return
			}

			require.NoError(t, err)
			assert.Equal(t, test.expectedOutput, out)
		})
	}
}
