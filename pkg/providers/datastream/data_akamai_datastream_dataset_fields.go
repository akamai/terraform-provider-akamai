package datastream

import (
	"context"
	"errors"
	"fmt"

	"github.com/akamai/AkamaiOPEN-edgegrid-golang/v6/pkg/datastream"
	"github.com/akamai/AkamaiOPEN-edgegrid-golang/v6/pkg/session"
	"github.com/akamai/terraform-provider-akamai/v3/pkg/akamai"
	"github.com/akamai/terraform-provider-akamai/v3/pkg/common/tf"
	"github.com/akamai/terraform-provider-akamai/v3/pkg/tools"
	"github.com/hashicorp/terraform-plugin-sdk/v2/diag"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/schema"
)

func dataSourceDatasetFields() *schema.Resource {
	return &schema.Resource{
		ReadContext: dataSourceDatasetFieldsRead,
		Schema: map[string]*schema.Schema{
			"fields": {
				Type:        schema.TypeSet,
				Computed:    true,
				Description: "Provides information about groups of dataset fields available in a given template",
				Elem: &schema.Resource{
					Schema: map[string]*schema.Schema{
						"dataset_group_name": {
							Type:        schema.TypeString,
							Computed:    true,
							Description: "A name of the dataset group",
						},
						"dataset_group_description": {
							Type:        schema.TypeString,
							Computed:    true,
							Description: "Describes the dataset group",
						},
						"dataset_fields": {
							Type:        schema.TypeList,
							Computed:    true,
							Description: "A list of data set fields available within the data set group",
							Elem: &schema.Resource{
								Schema: map[string]*schema.Schema{
									"dataset_field_description": {
										Type:        schema.TypeString,
										Computed:    true,
										Description: "Describes the data set field",
									},
									"dataset_field_id": {
										Type:        schema.TypeInt,
										Computed:    true,
										Description: "Identifies the field",
									},
									"dataset_field_json_key": {
										Type:        schema.TypeString,
										Computed:    true,
										Description: "Specifies the JSON key for the field in a log line",
									},
									"dataset_field_name": {
										Type:        schema.TypeString,
										Computed:    true,
										Description: "A name of the data set field",
									},
								},
							},
						},
					},
				},
			},
			"template_name": {
				Type:        schema.TypeString,
				Optional:    true,
				Default:     "EDGE_LOGS",
				Description: "The name of the data set template that you want to use in your stream configuration",
			},
		},
	}
}

func dataSourceDatasetFieldsRead(ctx context.Context, rd *schema.ResourceData, m interface{}) diag.Diagnostics {
	meta := akamai.Meta(m)
	logger := meta.Log("datastream", "dataSourceDatasetFieldsRead")

	ctx = session.ContextWithOptions(
		ctx,
		session.WithContextLog(logger),
	)

	logger.Debug("Listing dataset fields")
	client := inst.Client(meta)

	template, err := tf.GetStringValue("template_name", rd)
	if err != nil && !errors.Is(err, tf.ErrNotFound) {
		return diag.FromErr(err)
	}
	getDatasetFieldsRequest := datastream.GetDatasetFieldsRequest{
		TemplateName: datastream.TemplateName(template),
	}
	dataSets, err := client.GetDatasetFields(ctx, getDatasetFieldsRequest)
	if err != nil {
		return diag.FromErr(err)
	}

	fields := parseFields(dataSets)

	if err := rd.Set("fields", fields); err != nil {
		return diag.Errorf("%v: %s", tf.ErrValueSet, err.Error())
	}

	// ignoring the GetMd5Sum error, because `fields` is already initialized
	md5Sum, _ := tools.GetMd5Sum(fmt.Sprintf("%v", fields))
	rd.SetId(md5Sum)

	return nil
}

func parseFields(dataSets []datastream.DataSets) []map[string]interface{} {
	fields := make([]map[string]interface{}, 0, len(dataSets))

	for _, dataSet := range dataSets {
		data := map[string]interface{}{}
		data["dataset_group_name"] = dataSet.DatasetGroupName
		data["dataset_group_description"] = dataSet.DatasetGroupDescription
		data["dataset_fields"] = parseDatasetFields(dataSet.DatasetFields)
		fields = append(fields, data)
	}

	return fields
}

func parseDatasetFields(datasetFields []datastream.DatasetFields) []map[string]interface{} {
	dSFields := make([]map[string]interface{}, 0, len(datasetFields))
	for _, dataSetFields := range datasetFields {
		dataSetFieldsData := map[string]interface{}{}
		dataSetFieldsData["dataset_field_description"] = dataSetFields.DatasetFieldDescription
		dataSetFieldsData["dataset_field_id"] = dataSetFields.DatasetFieldID
		dataSetFieldsData["dataset_field_json_key"] = dataSetFields.DatasetFieldJsonKey
		dataSetFieldsData["dataset_field_name"] = dataSetFields.DatasetFieldName
		dSFields = append(dSFields, dataSetFieldsData)
	}
	return dSFields
}
