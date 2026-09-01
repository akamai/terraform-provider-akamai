package datastream

import (
	"context"
	"errors"
	"fmt"
	"strings"

	"github.com/akamai/AkamaiOPEN-edgegrid-golang/v14/pkg/session"
	"github.com/akamai/terraform-provider-akamai/v11/pkg/common/ptr"
	"github.com/akamai/terraform-provider-akamai/v11/pkg/common/tf"
	"github.com/akamai/terraform-provider-akamai/v11/pkg/meta"

	"github.com/akamai/AkamaiOPEN-edgegrid-golang/v14/pkg/datastream"
	"github.com/hashicorp/terraform-plugin-sdk/v2/diag"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/schema"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/validation"
)

func dataAkamaiDatastreamStreams() *schema.Resource {
	return &schema.Resource{
		Description: "Retrieves a list of all streams optionally by the specified GroupID.",
		ReadContext: dataDatastreamStreamsRead,
		Schema: map[string]*schema.Schema{
			"log_type": {
				Type:        schema.TypeString,
				Optional:    true,
				Description: "The log type of the stream. Valid values are `CDN`, `APPSEC`, and `ANSWERX`. If not specified, defaults to 'CDN'.",
				Default:     string(datastream.LogTypeCDN),
				ValidateDiagFunc: validation.ToDiagFunc(validation.StringInSlice([]string{
					string(datastream.LogTypeCDN),
					string(datastream.LogTypeAppSec),
					string(datastream.LogTypeAnswerX),
				}, true)),
			},
			"group_id": {
				Type:        schema.TypeInt,
				Optional:    true,
				Description: "Identifies the group where the stream is created.",
			},
			"streams_details": {
				Type:        schema.TypeSet,
				Computed:    true,
				Description: "List of streams",
				Elem: &schema.Resource{
					Schema: map[string]*schema.Schema{
						"log_type": {
							Type:        schema.TypeString,
							Computed:    true,
							Description: "The log type of the stream.",
						},
						"stream_name": {
							Type:        schema.TypeString,
							Computed:    true,
							Description: "The name of the stream.",
						},
						"stream_id": {
							Type:        schema.TypeInt,
							Computed:    true,
							Description: "Identifies the stream.",
						},
						"stream_version": {
							Type:        schema.TypeInt,
							Computed:    true,
							Description: "Identifies the current version of the stream.",
						},
						"group_id": {
							Type:        schema.TypeInt,
							Computed:    true,
							Description: "Identifies the group where the stream is created.",
						},
						"contract_id": {
							Type:        schema.TypeString,
							Computed:    true,
							Description: "Identifies the contract that the stream is associated with.",
						},
						"properties": {
							Type:        schema.TypeList,
							Computed:    true,
							Description: "List of properties associated with the stream.",
							Elem: &schema.Resource{
								Schema: map[string]*schema.Schema{
									"property_id": {
										Type:        schema.TypeInt,
										Computed:    true,
										Description: "The identifier of the property.",
									},
									"property_name": {
										Type:        schema.TypeString,
										Computed:    true,
										Description: "The descriptive label for the property.",
									},
									"integration_type": {
										Type:        schema.TypeString,
										Computed:    true,
										Description: "The integration mode for the property in datastream (e.g., PM_DEPENDENT, HYBRID, DS_MANAGED).",
									},
								},
							},
						},
						"app_sec_configs": {
							Type:        schema.TypeList,
							Computed:    true,
							Description: "List of AppSec configs associated with the stream.",
							Elem: &schema.Resource{
								Schema: map[string]*schema.Schema{
									"app_sec_id": {
										Type:        schema.TypeInt,
										Computed:    true,
										Description: "The identifier of the AppSec config.",
									},
									"app_sec_name": {
										Type:        schema.TypeString,
										Computed:    true,
										Description: "The descriptive label for the AppSec config.",
									},
								},
							},
						},
						"service_ids": {
							Type:        schema.TypeSet,
							Computed:    true,
							Description: "Set of service IDs associated with the stream.",
							Elem: &schema.Resource{
								Schema: map[string]*schema.Schema{
									"id": {
										Type:        schema.TypeInt,
										Computed:    true,
										Description: "Service ID monitored in the stream.",
									},
									"name": {
										Type:        schema.TypeString,
										Computed:    true,
										Description: "Name of the service ID.",
									},
									"product": {
										Type:        schema.TypeString,
										Computed:    true,
										Description: "The product associated with the service ID.",
									},
								},
							},
						},
						"latest_version": {
							Type:        schema.TypeInt,
							Computed:    true,
							Description: "Identifies the latestVersion version of the stream.",
						},
						"product_id": {
							Type:        schema.TypeString,
							Computed:    true,
							Description: "The productId.",
						},
						"stream_status": {
							Type:        schema.TypeString,
							Computed:    true,
							Description: "The activation status of the stream.",
						},

						"created_by": {
							Type:        schema.TypeString,
							Computed:    true,
							Description: "The username who created the stream.",
						},
						"created_date": {
							Type:        schema.TypeString,
							Computed:    true,
							Description: "The date and time when the stream was created.",
						},
						"modified_by": {
							Type:        schema.TypeString,
							Computed:    true,
							Description: "The username who activated or deactivated the stream",
						},
						"modified_date": {
							Type:        schema.TypeString,
							Computed:    true,
							Description: "The date and time when activation status was modified",
						},
						"integration_type": {
							Type:        schema.TypeString,
							Computed:    true,
							Description: "The integration mode for the stream in datastream (e.g., PM_DEPENDENT, HYBRID, DS_MANAGED)",
						},
					},
				},
			},
		},
	}
}

func dataDatastreamStreamsRead(ctx context.Context, d *schema.ResourceData, m interface{}) diag.Diagnostics {
	meta := meta.Must(m)
	logger := meta.Log("datastream", "dataDatastreamStreamsRead")

	ctx = session.ContextWithOptions(
		ctx,
		session.WithContextLog(logger),
	)

	client := inst.Client(meta)

	groupIDInt, err := tf.GetIntValue("group_id", d)
	if err != nil && !errors.Is(err, tf.ErrNotFound) {
		return diag.FromErr(err)
	}

	logType, err := tf.GetStringValue("log_type", d)
	if err != nil && !errors.Is(err, tf.ErrNotFound) {
		return diag.FromErr(err)
	}

	req := datastream.ListStreamsRequest{}
	req.LogType = datastream.LogType(strings.ToUpper(logType))
	resID := "akamai_datastreams"
	if groupIDInt != 0 {

		req.GroupID = ptr.To(groupIDInt)
		resID = fmt.Sprintf("%s_%d", resID, groupIDInt)
	}

	streams, err := client.ListStreams(ctx, req)
	if err != nil {
		return diag.FromErr(err)
	}

	logger.Debugf("Fetched %d streams", len(streams))

	attrs := createStreamsAttrs(streams)

	if err := d.Set("streams_details", attrs); err != nil {
		return diag.Errorf("%v: %s", tf.ErrValueSet, err.Error())
	}

	d.SetId(resID)
	return nil
}

func createStreamsAttrs(streams []datastream.StreamDetails) []interface{} {
	streamsAttrs := make([]interface{}, 0, len(streams))
	for _, stream := range streams {
		streamAttr := map[string]interface{}{
			"stream_status":   stream.StreamStatus,
			"contract_id":     stream.ContractID,
			"created_by":      stream.CreatedBy,
			"created_date":    stream.CreatedDate,
			"modified_by":     stream.ModifiedBy,
			"modified_date":   stream.ModifiedDate,
			"group_id":        stream.GroupID,
			"latest_version":  stream.LatestVersion,
			"product_id":      stream.ProductID,
			"properties":      createPropertiesAttrs(stream.Properties),
			"stream_id":       stream.StreamID,
			"stream_name":     stream.StreamName,
			"stream_version":  stream.StreamVersion,
			"log_type":        string(stream.LogType),
			"app_sec_configs": createAppSecConfigsAttrs(stream.AppSecConfigs),
			"service_ids":     createAnswerXServiceIDsAttrs(stream.AnswerXServiceIDs),
		}
		// Only set integration_type if it's non-empty (API may not return the field)
		if stream.IntegrationType != "" {
			streamAttr["integration_type"] = stream.IntegrationType
		}
		streamsAttrs = append(streamsAttrs, streamAttr)
	}

	return streamsAttrs
}

func createAppSecConfigsAttrs(configs []datastream.AppSecConfig) []interface{} {
	configsAttrs := make([]interface{}, 0, len(configs))
	for _, config := range configs {
		configAttr := map[string]interface{}{
			"app_sec_id":   config.AppSecID,
			"app_sec_name": config.AppSecName,
		}
		configsAttrs = append(configsAttrs, configAttr)
	}

	return configsAttrs
}

func createPropertiesAttrs(properties []datastream.Property) []interface{} {
	propertyAttrs := make([]interface{}, 0, len(properties))

	for _, property := range properties {
		propertyAttr := map[string]interface{}{
			"property_id":   property.PropertyID,
			"property_name": property.PropertyName,
		}
		// Only set integration_type if it's non-empty (API may not return the field)
		if property.IntegrationType != "" {
			propertyAttr["integration_type"] = property.IntegrationType
		}
		propertyAttrs = append(propertyAttrs, propertyAttr)
	}

	return propertyAttrs
}

func createAnswerXServiceIDsAttrs(answerXServiceIDs []datastream.AnswerXServiceDetail) []any {
	result := make([]any, 0, len(answerXServiceIDs))
	for _, serviceID := range answerXServiceIDs {
		result = append(result, map[string]any{
			"id":      serviceID.SSID,
			"name":    serviceID.Name,
			"product": serviceID.Product,
		})
	}
	return result
}
