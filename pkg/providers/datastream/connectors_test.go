package datastream

import (
	"testing"

	"github.com/akamai/AkamaiOPEN-edgegrid-golang/v12/pkg/datastream"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/schema"
	"github.com/stretchr/testify/assert"
)

var resourceSchema = map[string]*schema.Schema{
	"sumologic_connector": datastreamResourceSchema["sumologic_connector"],
	"invalid_connector": {
		Type:     schema.TypeSet,
		MaxItems: 1,
		Optional: true,
		Elem: &schema.Resource{
			Schema: map[string]*schema.Schema{
				"display_name": {
					Type:     schema.TypeString,
					Required: true,
				},
			},
		},
	},
}

func TestConnectorToMap(t *testing.T) {
	type expectedResult struct {
		key   string
		props map[string]interface{}
		error string
	}

	tests := map[string]struct {
		connectorDetails datastream.Destination
		resourceMap      map[string]interface{}
		expectedResult
	}{
		"empty connector details": {
			connectorDetails: datastream.Destination{},
			expectedResult: expectedResult{
				key:   "",
				props: nil,
				error: "",
			},
		},
		"no resource name for invalid connector type": {
			connectorDetails: datastream.Destination{

				DestinationType: datastream.DestinationType("invalid_connector"),
			},
			resourceMap: nil,
			expectedResult: expectedResult{
				error: "cannot find",
			},
		},
		"no connector in local resource": {
			connectorDetails: datastream.Destination{
				CompressLogs:      true,
				DisplayName:       "sumologic connector",
				DestinationType:   datastream.DestinationTypeSumoLogic,
				Endpoint:          "sumologic endpoint",
				ContentType:       "application/json",
				CustomHeaderName:  "custom_header_name",
				CustomHeaderValue: "custom_header_value",
			},
			resourceMap: nil,
			expectedResult: expectedResult{
				key: "sumologic_connector",
				props: map[string]interface{}{
					"collector_code":      "",
					"compress_logs":       true,
					"display_name":        "sumologic connector",
					"endpoint":            "sumologic endpoint",
					"content_type":        "application/json",
					"custom_header_name":  "custom_header_name",
					"custom_header_value": "custom_header_value",
				},
			},
		},
		"proper configuration": {
			connectorDetails: datastream.Destination{
				CompressLogs:      true,
				DisplayName:       "sumologic connector",
				DestinationType:   datastream.DestinationTypeSumoLogic,
				Endpoint:          "sumologic endpoint",
				ContentType:       "application/json",
				CustomHeaderName:  "custom_header_name",
				CustomHeaderValue: "custom_header_value",
			},
			resourceMap: map[string]interface{}{
				"sumologic_connector": []interface{}{
					map[string]interface{}{
						"collector_code":      "sumologic_collector_code",
						"compress_logs":       true,
						"display_name":        "sumologic connector",
						"endpoint":            "sumologic endpoint",
						"content_type":        "application/json",
						"custom_header_name":  "custom_header_name",
						"custom_header_value": "custom_header_value",
					},
				},
			},
			expectedResult: expectedResult{
				key: "sumologic_connector",
				props: map[string]interface{}{
					"collector_code":      "sumologic_collector_code",
					"compress_logs":       true,
					"display_name":        "sumologic connector",
					"endpoint":            "sumologic endpoint",
					"content_type":        "application/json",
					"custom_header_name":  "custom_header_name",
					"custom_header_value": "custom_header_value",
				},
			},
		},
	}
	for name, test := range tests {
		t.Run(name, func(t *testing.T) {

			d := schema.TestResourceDataRaw(t, resourceSchema, test.resourceMap)
			resourceKey, properties, err := ConnectorToMap(test.connectorDetails, d)

			expectedResult := test.expectedResult
			errMessage := expectedResult.error
			if errMessage != "" {
				assert.Contains(t, err.Error(), errMessage)
			} else {
				assert.Equal(t, expectedResult.key, resourceKey)
				assert.Equal(t, expectedResult.props, properties)
			}
		})
	}
}

func TestGetConnectors(t *testing.T) {
	tests := map[string]struct {
		resourceMap    map[string]interface{}
		expectedResult datastream.AbstractConnector
		errorMessage   string
	}{
		"missing connector definition": {
			resourceMap:  nil,
			errorMessage: "missing connector",
		},
		"proper configuration": {
			resourceMap: map[string]interface{}{
				"sumologic_connector": []interface{}{
					map[string]interface{}{
						"collector_code":      "sumologic_collector_code",
						"compress_logs":       true,
						"display_name":        "sumologic connector",
						"endpoint":            "sumologic endpoint",
						"content_type":        "application/json",
						"custom_header_name":  "custom_header_name",
						"custom_header_value": "custom_header_value",
					},
				},
			},
			expectedResult: &datastream.SumoLogicConnector{
				CollectorCode:     "sumologic_collector_code",
				CompressLogs:      true,
				DisplayName:       "sumologic connector",
				Endpoint:          "sumologic endpoint",
				ContentType:       "application/json",
				CustomHeaderName:  "custom_header_name",
				CustomHeaderValue: "custom_header_value",
			},
		},
	}
	for name, test := range tests {
		t.Run(name, func(t *testing.T) {
			d := schema.TestResourceDataRaw(t, resourceSchema, test.resourceMap)
			connectors, err := GetConnectors(d, []string{"sumologic_connector", "invalid_connector"})

			errMessage := test.errorMessage
			if errMessage != "" {
				assert.Contains(t, err.Error(), errMessage)
			} else {
				assert.Equal(t, test.expectedResult, connectors)
			}
		})
	}
}

func TestSetNonNilItemsFromState(t *testing.T) {
	defaultStateValues := map[string]interface{}{
		"a": 1,
		"b": 2,
		"c": nil,
		"d": 4,
		"e": nil,
	}
	tests := map[string]struct {
		expectedState map[string]interface{}
		stateValues   map[string]interface{}
		fields        []string
	}{
		"extract all existing": {
			expectedState: map[string]interface{}{
				"a": 1,
				"b": 2,
			},
			stateValues: defaultStateValues,
			fields:      []string{"a", "b"},
		},
		"extract not existing": {
			expectedState: map[string]interface{}{},
			stateValues:   defaultStateValues,
			fields:        []string{"aa", "bb"},
		},
		"empty fields to extract": {
			expectedState: map[string]interface{}{},
			stateValues:   defaultStateValues,
			fields:        []string{},
		},
		"empty state": {
			expectedState: map[string]interface{}{},
			stateValues:   map[string]interface{}{},
			fields:        []string{"a", "b"},
		},
	}

	for name, test := range tests {
		t.Run(name, func(t *testing.T) {
			state := map[string]interface{}{}
			setNonNilItemsFromState(test.stateValues, state, test.fields...)
			assert.Equal(t, test.expectedState, state)
		})
	}
}
