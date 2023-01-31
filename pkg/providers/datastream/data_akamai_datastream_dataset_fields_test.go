package datastream

import (
	"errors"
	"fmt"
	"regexp"
	"testing"

	"github.com/akamai/AkamaiOPEN-edgegrid-golang/v4/pkg/datastream"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/resource"
	"github.com/stretchr/testify/mock"
)

func TestDataSourceDatasetFieldsRead(t *testing.T) {

	normalServerResponse := []datastream.DataSets{
		{
			DatasetGroupName:        "Log information",
			DatasetGroupDescription: "Contains fields that can be used to identify or tag a log line",
			DatasetFields: []datastream.DatasetFields{
				{
					DatasetFieldID:          1000,
					DatasetFieldName:        "CP code",
					DatasetFieldDescription: "The Content Provider code associated with the request.",
					DatasetFieldJsonKey:     "cp",
				},
				{
					DatasetFieldID:          1002,
					DatasetFieldName:        "Request ID",
					DatasetFieldDescription: "The identifier of the request.",
					DatasetFieldJsonKey:     "reqId",
				},
				{
					DatasetFieldID:          1100,
					DatasetFieldName:        "Request time",
					DatasetFieldDescription: "The time when the edge server accepted the request from the client.",
					DatasetFieldJsonKey:     "reqTimeSec",
				},
			},
		},
		{
			DatasetGroupName:        "Message exchange data",
			DatasetGroupDescription: "Contains fields representing the exchange of data between Akamai & end user",
			DatasetFields: []datastream.DatasetFields{
				{
					DatasetFieldID:          1005,
					DatasetFieldName:        "Bytes",
					DatasetFieldDescription: "The content bytes served in the response.",
					DatasetFieldJsonKey:     "bytes",
				},
				{
					DatasetFieldID:          1006,
					DatasetFieldName:        "Client IP",
					DatasetFieldDescription: "The IP address of the client.",
					DatasetFieldJsonKey:     "cliIP",
				},
			},
		},
	}
	normalChecks := []resource.TestCheckFunc{
		resource.TestCheckResourceAttr("data.akamai_datastream_dataset_fields.test", "fields.#", "2"),

		resource.TestCheckResourceAttr("data.akamai_datastream_dataset_fields.test", "fields.0.dataset_group_name", "Log information"),
		resource.TestCheckResourceAttr("data.akamai_datastream_dataset_fields.test", "fields.0.dataset_group_description", "Contains fields that can be used to identify or tag a log line"),
		resource.TestCheckResourceAttr("data.akamai_datastream_dataset_fields.test", "fields.0.dataset_fields.0.dataset_field_description", "The Content Provider code associated with the request."),
		resource.TestCheckResourceAttr("data.akamai_datastream_dataset_fields.test", "fields.0.dataset_fields.0.dataset_field_id", "1000"),
		resource.TestCheckResourceAttr("data.akamai_datastream_dataset_fields.test", "fields.0.dataset_fields.0.dataset_field_json_key", "cp"),
		resource.TestCheckResourceAttr("data.akamai_datastream_dataset_fields.test", "fields.0.dataset_fields.0.dataset_field_name", "CP code"),
		resource.TestCheckResourceAttr("data.akamai_datastream_dataset_fields.test", "fields.0.dataset_fields.1.dataset_field_description", "The identifier of the request."),
		resource.TestCheckResourceAttr("data.akamai_datastream_dataset_fields.test", "fields.0.dataset_fields.1.dataset_field_id", "1002"),
		resource.TestCheckResourceAttr("data.akamai_datastream_dataset_fields.test", "fields.0.dataset_fields.1.dataset_field_json_key", "reqId"),
		resource.TestCheckResourceAttr("data.akamai_datastream_dataset_fields.test", "fields.0.dataset_fields.1.dataset_field_name", "Request ID"),
		resource.TestCheckResourceAttr("data.akamai_datastream_dataset_fields.test", "fields.0.dataset_fields.2.dataset_field_description", "The time when the edge server accepted the request from the client."),
		resource.TestCheckResourceAttr("data.akamai_datastream_dataset_fields.test", "fields.0.dataset_fields.2.dataset_field_id", "1100"),
		resource.TestCheckResourceAttr("data.akamai_datastream_dataset_fields.test", "fields.0.dataset_fields.2.dataset_field_json_key", "reqTimeSec"),
		resource.TestCheckResourceAttr("data.akamai_datastream_dataset_fields.test", "fields.0.dataset_fields.2.dataset_field_name", "Request time"),
		resource.TestCheckResourceAttr("data.akamai_datastream_dataset_fields.test", "fields.1.dataset_group_name", "Message exchange data"),
		resource.TestCheckResourceAttr("data.akamai_datastream_dataset_fields.test", "fields.1.dataset_group_description", "Contains fields representing the exchange of data between Akamai & end user"),
		resource.TestCheckResourceAttr("data.akamai_datastream_dataset_fields.test", "fields.1.dataset_fields.0.dataset_field_description", "The content bytes served in the response."),
		resource.TestCheckResourceAttr("data.akamai_datastream_dataset_fields.test", "fields.1.dataset_fields.0.dataset_field_id", "1005"),
		resource.TestCheckResourceAttr("data.akamai_datastream_dataset_fields.test", "fields.1.dataset_fields.0.dataset_field_json_key", "bytes"),
		resource.TestCheckResourceAttr("data.akamai_datastream_dataset_fields.test", "fields.1.dataset_fields.0.dataset_field_name", "Bytes"),
		resource.TestCheckResourceAttr("data.akamai_datastream_dataset_fields.test", "fields.1.dataset_fields.1.dataset_field_description", "The IP address of the client."),
		resource.TestCheckResourceAttr("data.akamai_datastream_dataset_fields.test", "fields.1.dataset_fields.1.dataset_field_id", "1006"),
		resource.TestCheckResourceAttr("data.akamai_datastream_dataset_fields.test", "fields.1.dataset_fields.1.dataset_field_json_key", "cliIP"),
		resource.TestCheckResourceAttr("data.akamai_datastream_dataset_fields.test", "fields.1.dataset_fields.1.dataset_field_name", "Client IP"),
	}
	tests := map[string]struct {
		configPath    string
		edgegridData  []datastream.DataSets
		edgegridError error
		checks        []resource.TestCheckFunc
		withError     error
	}{
		"EDGE_LOGS template": {
			configPath:   "testdata/TestDataSourceDatasetFieldsRead/edge_logs.tf",
			edgegridData: normalServerResponse,
			checks:       normalChecks,
		},
		"no template, EDGE_LOGS by default": {
			configPath:   "testdata/TestDataSourceDatasetFieldsRead/edge_logs_no_template.tf",
			edgegridData: normalServerResponse,
			checks:       normalChecks,
		},
		"empty server response": {
			configPath:   "testdata/TestDataSourceDatasetFieldsRead/edge_logs.tf",
			edgegridData: []datastream.DataSets{},
			checks: []resource.TestCheckFunc{
				resource.TestCheckResourceAttr("data.akamai_datastream_dataset_fields.test", "fields.#", "0"),
			},
		},
		"EDGE_LOGS template: edgegrid error": {
			configPath:    "testdata/TestDataSourceDatasetFieldsRead/edge_logs.tf",
			edgegridError: fmt.Errorf("%w: request failed: %s", datastream.ErrGetDatasetFields, errors.New("500")),
			withError:     datastream.ErrGetDatasetFields,
		},
	}
	for name, test := range tests {
		t.Run(name, func(t *testing.T) {
			client := &datastream.Mock{}
			client.On("GetDatasetFields", mock.Anything, datastream.GetDatasetFieldsRequest{
				TemplateName: datastream.TemplateNameEdgeLogs,
			}).Return(test.edgegridData, test.edgegridError)
			useClient(client, func() {
				if test.withError == nil {
					resource.UnitTest(t, resource.TestCase{
						Providers: testAccProviders,
						Steps: []resource.TestStep{
							{
								Config: loadFixtureString(test.configPath),
								Check: resource.ComposeAggregateTestCheckFunc(
									test.checks...,
								),
							},
						},
					})
				} else {
					resource.UnitTest(t, resource.TestCase{
						Providers: testAccProviders,
						Steps: []resource.TestStep{
							{
								Config:      loadFixtureString(test.configPath),
								ExpectError: regexp.MustCompile(test.withError.Error()),
							},
						},
					})
				}
			})
		})
	}
}
