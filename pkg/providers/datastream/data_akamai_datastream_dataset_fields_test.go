package datastream

import (
	"regexp"
	"testing"

	"github.com/akamai/AkamaiOPEN-edgegrid-golang/v8/pkg/datastream"
	"github.com/akamai/terraform-provider-akamai/v6/pkg/common/testutils"
	"github.com/hashicorp/terraform-plugin-testing/helper/resource"
	"github.com/stretchr/testify/mock"
)

func TestDataSourceDatasetFieldsRead(t *testing.T) {

	tests := map[string]struct {
		configPath             string
		getDatasetFieldsReturn *datastream.DataSets
		checkFuncs             []resource.TestCheckFunc
		edgegridError          error
		withError              *regexp.Regexp
	}{

		"validate dataset fields response": {
			configPath: "testdata/TestDataSourceDatasetFieldsRead/list_dataset_fields_with_product.tf",
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
		"no template EDGE_LOGS by default": {
			configPath: "testdata/TestDataSourceDatasetFieldsRead/list_dataset_fields_default_product.tf",
			getDatasetFieldsReturn: &datastream.DataSets{
				DataSetFields: []datastream.DataSetField{},
			},
			checkFuncs: []resource.TestCheckFunc{},
		},
		"empty server response": {
			configPath: "testdata/TestDataSourceDatasetFieldsRead/list_dataset_fields_with_product.tf",
			getDatasetFieldsReturn: &datastream.DataSets{
				DataSetFields: []datastream.DataSetField{},
			},
			checkFuncs: []resource.TestCheckFunc{
				resource.TestCheckResourceAttr("data.akamai_datastream_dataset_fields.test", "dataset_fields.#", "0"),
			},
		},
	}
	for name, test := range tests {
		t.Run(name, func(t *testing.T) {
			client := &datastream.Mock{}
			client.On("GetDatasetFields", mock.Anything, mock.Anything).Return(test.getDatasetFieldsReturn, test.edgegridError)
			useClient(client, func() {
				if test.withError == nil {
					resource.UnitTest(t, resource.TestCase{
						ProtoV6ProviderFactories: testutils.NewProtoV6ProviderFactory(NewSubprovider()),
						Steps: []resource.TestStep{
							{
								Config: testutils.LoadFixtureString(t, test.configPath),
								Check: resource.ComposeAggregateTestCheckFunc(
									test.checkFuncs...,
								),
							},
						},
					})
				} else {
					resource.UnitTest(t, resource.TestCase{
						ProtoV6ProviderFactories: testutils.NewProtoV6ProviderFactory(NewSubprovider()),
						Steps: []resource.TestStep{
							{
								Config:      testutils.LoadFixtureString(t, test.configPath),
								ExpectError: test.withError,
							},
						},
					})
				}
			})
		})
	}
}
