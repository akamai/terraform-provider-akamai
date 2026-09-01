package datastream

import (
	"regexp"
	"testing"

	"github.com/akamai/AkamaiOPEN-edgegrid-golang/v13/pkg/datastream"
	"github.com/akamai/terraform-provider-akamai/v11/pkg/common/test"
	"github.com/akamai/terraform-provider-akamai/v11/pkg/common/testutils"
	"github.com/hashicorp/terraform-plugin-testing/helper/resource"
)

func TestDataSourceDatasetFieldsRead(t *testing.T) {
	t.Parallel()
	tests := map[string]struct {
		configPath             string
		expectedRequest        datastream.GetDatasetFieldsRequest
		getDatasetFieldsReturn *datastream.DataSets
		checkFuncs             []resource.TestCheckFunc
		withError              *regexp.Regexp
	}{

		"validate dataset fields response": {
			configPath: "testdata/TestDataSourceDatasetFieldsRead/list_dataset_fields_with_product.tf",
			expectedRequest: datastream.GetDatasetFieldsRequest{
				LogType:   datastream.LogTypeCDN,
				ProductID: "PROD_1",
			},
			getDatasetFieldsReturn: &datastream.DataSets{
				DataSetFields: []datastream.DataSetField{
					{
						DatasetFieldID:          1000,
						DatasetFieldName:        "datasetFieldName_1",
						DatasetFieldJsonKey:     "datasetFieldJsonKey_1",
						DatasetFieldGroup:       "datasetFieldGroup_1",
						DatasetFieldDescription: "datasetFieldDescription_1",
					},
					{
						DatasetFieldID:          1001,
						DatasetFieldName:        "datasetFieldName_2",
						DatasetFieldJsonKey:     "datasetFieldJsonKey_2",
						DatasetFieldGroup:       "datasetFieldGroup_2",
						DatasetFieldDescription: "datasetFieldDescription_2",
					},
					{
						DatasetFieldID:          1002,
						DatasetFieldName:        "datasetFieldName_3",
						DatasetFieldJsonKey:     "datasetFieldJsonKey_3",
						DatasetFieldGroup:       "datasetFieldGroup_3",
						DatasetFieldDescription: "datasetFieldDescription_3",
					},
				},
			},

			checkFuncs: []resource.TestCheckFunc{
				resource.TestCheckResourceAttr("data.akamai_datastream_dataset_fields.test", "dataset_fields.#", "3"),
				resource.TestCheckResourceAttr("data.akamai_datastream_dataset_fields.test", "dataset_fields.0.dataset_field_description", "datasetFieldDescription_1"),
				resource.TestCheckResourceAttr("data.akamai_datastream_dataset_fields.test", "dataset_fields.0.dataset_field_id", "1000"),
				resource.TestCheckResourceAttr("data.akamai_datastream_dataset_fields.test", "dataset_fields.0.dataset_field_json_key", "datasetFieldJsonKey_1"),
				resource.TestCheckResourceAttr("data.akamai_datastream_dataset_fields.test", "dataset_fields.0.dataset_field_name", "datasetFieldName_1"),
				resource.TestCheckResourceAttr("data.akamai_datastream_dataset_fields.test", "dataset_fields.0.dataset_field_group", "datasetFieldGroup_1"),

				resource.TestCheckResourceAttr("data.akamai_datastream_dataset_fields.test", "dataset_fields.1.dataset_field_description", "datasetFieldDescription_2"),
				resource.TestCheckResourceAttr("data.akamai_datastream_dataset_fields.test", "dataset_fields.1.dataset_field_id", "1001"),
				resource.TestCheckResourceAttr("data.akamai_datastream_dataset_fields.test", "dataset_fields.1.dataset_field_json_key", "datasetFieldJsonKey_2"),
				resource.TestCheckResourceAttr("data.akamai_datastream_dataset_fields.test", "dataset_fields.1.dataset_field_name", "datasetFieldName_2"),
				resource.TestCheckResourceAttr("data.akamai_datastream_dataset_fields.test", "dataset_fields.1.dataset_field_group", "datasetFieldGroup_2"),

				resource.TestCheckResourceAttr("data.akamai_datastream_dataset_fields.test", "dataset_fields.2.dataset_field_description", "datasetFieldDescription_3"),
				resource.TestCheckResourceAttr("data.akamai_datastream_dataset_fields.test", "dataset_fields.2.dataset_field_id", "1002"),
				resource.TestCheckResourceAttr("data.akamai_datastream_dataset_fields.test", "dataset_fields.2.dataset_field_json_key", "datasetFieldJsonKey_3"),
				resource.TestCheckResourceAttr("data.akamai_datastream_dataset_fields.test", "dataset_fields.2.dataset_field_name", "datasetFieldName_3"),
				resource.TestCheckResourceAttr("data.akamai_datastream_dataset_fields.test", "dataset_fields.2.dataset_field_group", "datasetFieldGroup_3"),
			},
		},
		"validate dataset fields response AnswerX": {
			configPath: "testdata/TestDataSourceDatasetFieldsRead/list_dataset_fields_answerx_without_product.tf",
			expectedRequest: datastream.GetDatasetFieldsRequest{
				LogType:   datastream.LogTypeAnswerX,
				ProductID: "",
			},
			getDatasetFieldsReturn: &datastream.DataSets{
				DataSetFields: []datastream.DataSetField{
					{
						DatasetFieldID:          1000,
						DatasetFieldName:        "datasetFieldName_1",
						DatasetFieldJsonKey:     "datasetFieldJsonKey_1",
						DatasetFieldGroup:       "datasetFieldGroup_1",
						DatasetFieldDescription: "datasetFieldDescription_1",
					},
					{
						DatasetFieldID:          1001,
						DatasetFieldName:        "datasetFieldName_2",
						DatasetFieldJsonKey:     "datasetFieldJsonKey_2",
						DatasetFieldGroup:       "datasetFieldGroup_2",
						DatasetFieldDescription: "datasetFieldDescription_2",
					},
				},
			},

			checkFuncs: []resource.TestCheckFunc{
				test.NewStateChecker("data.akamai_datastream_dataset_fields.test").
					CheckEqual("dataset_fields.#", "2").
					CheckEqual("dataset_fields.0.dataset_field_description", "datasetFieldDescription_1").
					CheckEqual("dataset_fields.0.dataset_field_id", "1000").
					CheckEqual("dataset_fields.0.dataset_field_json_key", "datasetFieldJsonKey_1").
					CheckEqual("dataset_fields.0.dataset_field_name", "datasetFieldName_1").
					CheckEqual("dataset_fields.0.dataset_field_group", "datasetFieldGroup_1").
					CheckEqual("dataset_fields.1.dataset_field_description", "datasetFieldDescription_2").
					CheckEqual("dataset_fields.1.dataset_field_id", "1001").
					CheckEqual("dataset_fields.1.dataset_field_json_key", "datasetFieldJsonKey_2").
					CheckEqual("dataset_fields.1.dataset_field_name", "datasetFieldName_2").
					CheckEqual("dataset_fields.1.dataset_field_group", "datasetFieldGroup_2").
					Build(),
			},
		},
		"no template EDGE_LOGS by default": {
			configPath: "testdata/TestDataSourceDatasetFieldsRead/list_dataset_fields_default_product.tf",
			expectedRequest: datastream.GetDatasetFieldsRequest{
				LogType:   datastream.LogTypeCDN,
				ProductID: "",
			},
			getDatasetFieldsReturn: &datastream.DataSets{
				DataSetFields: []datastream.DataSetField{},
			},
			checkFuncs: []resource.TestCheckFunc{},
		},
		"empty server response": {
			configPath: "testdata/TestDataSourceDatasetFieldsRead/list_dataset_fields_with_product.tf",
			expectedRequest: datastream.GetDatasetFieldsRequest{
				LogType:   datastream.LogTypeCDN,
				ProductID: "PROD_1",
			},
			getDatasetFieldsReturn: &datastream.DataSets{
				DataSetFields: []datastream.DataSetField{},
			},
			checkFuncs: []resource.TestCheckFunc{
				resource.TestCheckResourceAttr("data.akamai_datastream_dataset_fields.test", "dataset_fields.#", "0"),
			},
		},
		"validation error - invalid log type": {
			configPath: "testdata/TestDataSourceDatasetFieldsRead/list_dataset_fields_invalid_log_type.tf",
			withError:  regexp.MustCompile(`expected log_type to be one of`),
		},
		"validation error - product_id field not supported for non-CDN log type": {
			configPath: "testdata/TestDataSourceDatasetFieldsRead/list_dataset_fields_answerx_with_product.tf",
			withError:  regexp.MustCompile(`product_id field is not supported for log_type`),
		},
	}
	for name, test := range tests {
		t.Run(name, func(t *testing.T) {
			t.Parallel()
			client := &datastream.Mock{}
			if test.withError == nil {
				client.On("GetDatasetFields", testutils.MockContext, test.expectedRequest).Return(test.getDatasetFieldsReturn, nil).Times(3)
			}

			step := resource.TestStep{
				Config:      testutils.LoadFixtureString(t, test.configPath),
				ExpectError: test.withError,
			}
			if test.withError == nil {
				step.Check = resource.ComposeAggregateTestCheckFunc(test.checkFuncs...)
			}

			useClient(client, func() {
				resource.UnitTest(t, resource.TestCase{
					ProtoV6ProviderFactories: testutils.NewProtoV6ProviderFactory(NewSubprovider()),
					Steps:                    []resource.TestStep{step},
				})
			})
			client.AssertExpectations(t)
		})
	}
}
