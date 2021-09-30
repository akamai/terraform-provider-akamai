package datastream

import (
	"testing"

	"github.com/akamai/terraform-provider-akamai/v2/pkg/tools"

	"github.com/stretchr/testify/assert"

	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/schema"

	"github.com/akamai/AkamaiOPEN-edgegrid-golang/v2/pkg/datastream"
)

var resourceSchema = map[string]*schema.Schema{
	"sumologic_connector": {
		Type:     schema.TypeSet,
		MaxItems: 1,
		Optional: true,
		Elem: &schema.Resource{
			Schema: map[string]*schema.Schema{
				"collector_code": {
					Type:     schema.TypeString,
					Required: true,
				},
				"compress_logs": {
					Type:     schema.TypeBool,
					Default:  true,
					Optional: true,
				},
				"connector_id": {
					Type:     schema.TypeInt,
					Computed: true,
				},
				"connector_name": {
					Type:     schema.TypeString,
					Required: true,
				},
				"endpoint": {
					Type:     schema.TypeString,
					Required: true,
				},
			},
		},
	},
	"invalid_connector": {
		Type:     schema.TypeSet,
		MaxItems: 1,
		Optional: true,
		Elem: &schema.Resource{
			Schema: map[string]*schema.Schema{
				"connector_name": {
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
		connectorDetails []datastream.ConnectorDetails
		resourceMap      map[string]interface{}
		expectedResult
	}{
		"empty connector details": {
			connectorDetails: []datastream.ConnectorDetails{},
			expectedResult: expectedResult{
				key:   "",
				props: nil,
				error: "",
			},
		},
		"more than one connector": {
			connectorDetails: []datastream.ConnectorDetails{
				{
					ConnectorType: datastream.ConnectorTypeS3,
				},
				{
					ConnectorType: datastream.ConnectorTypeGcs,
				},
			},
			expectedResult: expectedResult{
				key:   "",
				props: nil,
				error: "",
			},
		},
		"connector not found in resource": {
			connectorDetails: []datastream.ConnectorDetails{
				{
					ConnectorType: datastream.ConnectorTypeS3,
				},
			},
			resourceMap: nil,
			expectedResult: expectedResult{
				error: tools.ErrNotFound.Error(),
			},
		},
		"no resource name for invalid connector type": {
			connectorDetails: []datastream.ConnectorDetails{
				{
					ConnectorType: datastream.ConnectorType("invalid_connector"),
				},
			},
			resourceMap: nil,
			expectedResult: expectedResult{
				error: "cannot find",
			},
		},
		"proper configuration": {
			connectorDetails: []datastream.ConnectorDetails{
				{
					ConnectorID:   1337,
					CompressLogs:  true,
					ConnectorName: "sumologic connector",
					ConnectorType: datastream.ConnectorTypeSumoLogic,
					URL:           "sumologic endpoint",
				},
			},
			resourceMap: map[string]interface{}{
				"sumologic_connector": []interface{}{
					map[string]interface{}{
						"collector_code": "sumologic_collector_code",
						"compress_logs":  true,
						"connector_name": "sumologic connector",
						"endpoint":       "sumologic endpoint",
					},
				},
			},
			expectedResult: expectedResult{
				key: "sumologic_connector",
				props: map[string]interface{}{
					"collector_code": "sumologic_collector_code",
					"compress_logs":  true,
					"connector_id":   1337,
					"connector_name": "sumologic connector",
					"endpoint":       "sumologic endpoint",
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
		expectedResult []datastream.AbstractConnector
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
						"collector_code": "sumologic_collector_code",
						"compress_logs":  true,
						"connector_name": "sumologic connector",
						"endpoint":       "sumologic endpoint",
					},
				},
			},
			expectedResult: []datastream.AbstractConnector{
				&datastream.SumoLogicConnector{
					CollectorCode: "sumologic_collector_code",
					CompressLogs:  true,
					ConnectorName: "sumologic connector",
					Endpoint:      "sumologic endpoint",
				},
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
