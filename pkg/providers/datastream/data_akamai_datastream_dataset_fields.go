package datastream

import (
	"context"
	"errors"
	"fmt"

	"github.com/akamai/AkamaiOPEN-edgegrid-golang/v8/pkg/datastream"
	"github.com/akamai/AkamaiOPEN-edgegrid-golang/v8/pkg/session"
	"github.com/akamai/terraform-provider-akamai/v6/pkg/common/hash"
	"github.com/akamai/terraform-provider-akamai/v6/pkg/common/tf"
	"github.com/akamai/terraform-provider-akamai/v6/pkg/meta"
	"github.com/hashicorp/terraform-plugin-sdk/v2/diag"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/schema"
)

func dataSourceDatasetFields() *schema.Resource {
	return &schema.Resource{
		ReadContext: dataSourceDatasetFieldsRead,
		Schema: map[string]*schema.Schema{
			"product_id": {
				Type:        schema.TypeString,
				Optional:    true,
				Description: "Identifies the stream",
			},
			"dataset_fields": {
				Type:        schema.TypeList,
				Computed:    true,
				Description: "Provides information about groups of dataset fields available in a given template",

				Elem: &schema.Resource{
					Schema: map[string]*schema.Schema{

						"dataset_field_id": {
							Type:        schema.TypeInt,
							Computed:    true,
							Description: "Identifies the field",
						},
						"dataset_field_description": {
							Type:        schema.TypeString,
							Computed:    true,
							Description: "Describes the data set field",
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
						"dataset_field_group": {
							Type:        schema.TypeString,
							Computed:    true,
							Description: "A name of the group for data set field",
						},
					},
				},
			},
		},
	}
}

func dataSourceDatasetFieldsRead(ctx context.Context, rd *schema.ResourceData, m interface{}) diag.Diagnostics {
	meta := meta.Must(m)
	logger := meta.Log("datastream", "dataSourceDatasetFieldsRead")

	ctx = session.ContextWithOptions(
		ctx,
		session.WithContextLog(logger),
	)

	logger.Debug("Listing dataset fields")
	client := inst.Client(meta)

	productID, err := tf.GetStringValue("product_id", rd)
	if err != nil && !errors.Is(err, tf.ErrNotFound) {
		return diag.FromErr(err)
	}

	var getDatasetFieldsRequest datastream.GetDatasetFieldsRequest
	if productID != "" {
		getDatasetFieldsRequest = datastream.GetDatasetFieldsRequest{
			ProductID: &productID,
		}
	}

	dataSets, err := client.GetDatasetFields(ctx, getDatasetFieldsRequest)
	if err != nil {
		return diag.FromErr(err)
	}

	fields := parseFields(dataSets)

	if err := rd.Set("dataset_fields", fields); err != nil {
		return diag.Errorf("%v: %s", tf.ErrValueSet, err.Error())
	}

	// ignoring the GetMd5Sum error, because `fields` is already initialized
	md5Sum, _ := hash.GetMD5Sum(fmt.Sprintf("%v", fields))
	rd.SetId(md5Sum)

	return nil
}

func parseFields(dataSets *datastream.DataSets) []map[string]interface{} {

	var datasetFields = dataSets.DataSetFields
	return parseDatasetFields(datasetFields)
}

func parseDatasetFields(datasetFields []datastream.DataSetField) []map[string]interface{} {
	dSFields := make([]map[string]interface{}, 0, len(datasetFields))
	for _, dataSetFields := range datasetFields {
		dataSetFieldsData := map[string]interface{}{}
		dataSetFieldsData["dataset_field_description"] = dataSetFields.DatasetFieldDescription
		dataSetFieldsData["dataset_field_id"] = dataSetFields.DatasetFieldID
		dataSetFieldsData["dataset_field_json_key"] = dataSetFields.DatasetFieldJsonKey
		dataSetFieldsData["dataset_field_name"] = dataSetFields.DatasetFieldName
		dataSetFieldsData["dataset_field_group"] = dataSetFields.DatasetFieldGroup
		dSFields = append(dSFields, dataSetFieldsData)
	}
	return dSFields
}
