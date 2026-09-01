package datastream

import (
	"context"
	"errors"
	"fmt"
	"strings"

	"github.com/akamai/AkamaiOPEN-edgegrid-golang/v13/pkg/datastream"
	"github.com/akamai/AkamaiOPEN-edgegrid-golang/v13/pkg/session"
	"github.com/akamai/terraform-provider-akamai/v11/pkg/common/hash"

	"github.com/akamai/terraform-provider-akamai/v11/pkg/common/tf"
	"github.com/akamai/terraform-provider-akamai/v11/pkg/meta"
	"github.com/hashicorp/terraform-plugin-sdk/v2/diag"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/schema"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/validation"
)

func dataSourceDatasetFields() *schema.Resource {
	return &schema.Resource{
		ReadContext: dataSourceDatasetFieldsRead,
		Schema: map[string]*schema.Schema{
			"log_type": {
				Type:        schema.TypeString,
				Optional:    true,
				Default:     string(datastream.LogTypeCDN),
				Description: "The type of logs for which to retrieve dataset fields. Valid values are `CDN` and `ANSWERX`. If not specified, defaults to `CDN`.",
				ValidateDiagFunc: validation.ToDiagFunc(validation.StringInSlice([]string{
					string(datastream.LogTypeCDN),
					string(datastream.LogTypeAnswerX),
				}, true)),
			},
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

	logTypeVal, err := tf.GetStringValue("log_type", rd)
	if err != nil && !errors.Is(err, tf.ErrNotFound) {
		return diag.FromErr(err)
	}
	logType := datastream.LogType(strings.ToUpper(logTypeVal))

	productID, err := tf.GetStringValue("product_id", rd)
	if err != nil && !errors.Is(err, tf.ErrNotFound) {
		return diag.FromErr(err)
	}

	if productID != "" && logType != datastream.LogTypeCDN {
		return diag.Errorf("product_id field is not supported for log_type %q", logType)
	}

	getDatasetFieldsRequest := datastream.GetDatasetFieldsRequest{
		LogType:   logType,
		ProductID: productID,
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
