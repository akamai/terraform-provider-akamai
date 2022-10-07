package datastream

import (
	"context"
	"errors"
	"fmt"
	"strconv"
	"strings"

	"github.com/akamai/AkamaiOPEN-edgegrid-golang/v2/pkg/session"
	"github.com/akamai/terraform-provider-akamai/v2/pkg/akamai"
	"github.com/akamai/terraform-provider-akamai/v2/pkg/tools"

	"github.com/akamai/AkamaiOPEN-edgegrid-golang/v2/pkg/datastream"
	"github.com/hashicorp/terraform-plugin-sdk/v2/diag"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/schema"
)

func dataAkamaiDatastreamStreams() *schema.Resource {
	return &schema.Resource{
		Description: "Retrieves a list of all streams optionally by the specified GroupID.",
		ReadContext: dataDatastreamStreamsRead,
		Schema: map[string]*schema.Schema{
			"group_id": {
				Type:             schema.TypeString,
				Optional:         true,
				DiffSuppressFunc: tools.FieldPrefixSuppress("grp_"),
				Description:      "Limits the returned set to streams belonging to the specified group.",
			},
			"streams": {
				Type:        schema.TypeList,
				Computed:    true,
				Description: "The latest versions of the stream configurations.",
				Elem: &schema.Resource{
					Schema: map[string]*schema.Schema{
						"activation_status": {
							Type:        schema.TypeString,
							Computed:    true,
							Description: "The activation status of the stream.",
						},
						"archived": {
							Type:        schema.TypeBool,
							Computed:    true,
							Description: "Whether the stream is archived.",
						},
						"connectors": {
							Type:        schema.TypeString,
							Computed:    true,
							Description: "The connector where the stream sends logs.",
						},
						"contract_id": {
							Type:        schema.TypeString,
							Computed:    true,
							Description: "Identifies the contract that the stream is associated with.",
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
						"current_version_id": {
							Type:        schema.TypeInt,
							Computed:    true,
							Description: "Identifies the current version of the stream.",
						},
						"errors": {
							Type:        schema.TypeList,
							Computed:    true,
							Description: "Objects that may indicate stream failure errors",
							Elem: &schema.Resource{
								Schema: map[string]*schema.Schema{
									"detail": {
										Type:        schema.TypeString,
										Computed:    true,
										Description: "A message informing about the status of the failed stream.",
									},
									"title": {
										Type:        schema.TypeString,
										Computed:    true,
										Description: "A descriptive label for the type of error.",
									},
									"type": {
										Type:        schema.TypeString,
										Computed:    true,
										Description: "Identifies the error type.",
									},
								},
							},
						},
						"group_id": {
							Type:        schema.TypeInt,
							Computed:    true,
							Description: "Identifies the group where the stream is created.",
						},
						"group_name": {
							Type:        schema.TypeString,
							Computed:    true,
							Description: "The group name where the stream is created.",
						},
						"properties": {
							Type:        schema.TypeList,
							Computed:    true,
							Description: "List of properties associated with stream.",
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
								},
							},
						},
						"stream_id": {
							Type:        schema.TypeInt,
							Computed:    true,
							Description: "Identifies the stream.",
						},
						"stream_name": {
							Type:        schema.TypeString,
							Computed:    true,
							Description: "The name of the stream.",
						},
						"stream_type_name": {
							Type:        schema.TypeString,
							Computed:    true,
							Description: "Specifies the type of the data stream.",
						},
						"stream_version_id": {
							Type:        schema.TypeInt,
							Computed:    true,
							Description: "Identifies the version of the stream.",
						},
					},
				},
			},
		},
	}
}

func dataDatastreamStreamsRead(ctx context.Context, d *schema.ResourceData, m interface{}) diag.Diagnostics {
	meta := akamai.Meta(m)
	logger := meta.Log("datastream", "dataDatastreamStreamsRead")

	ctx = session.ContextWithOptions(
		ctx,
		session.WithContextLog(logger),
	)

	client := inst.Client(meta)

	req := datastream.ListStreamsRequest{}
	resID := "akamai_datastreams"
	groupIDStr, err := tools.GetStringValue("group_id", d)
	if err != nil && !errors.Is(err, tools.ErrNotFound) {
		return diag.FromErr(err)
	}

	if groupIDStr != "" {
		groupID, err := strconv.Atoi(strings.TrimPrefix(groupIDStr, "grp_"))
		if err != nil {
			return diag.FromErr(err)
		}

		req.GroupID = tools.IntPtr(groupID)
		resID = fmt.Sprintf("%s_%d", resID, groupID)
	}

	streams, err := client.ListStreams(ctx, req)
	if err != nil {
		return diag.FromErr(err)
	}

	logger.Debugf("Fetched %d streams", len(streams))
	attrs := map[string]interface{}{"streams": createStreamsAttrs(streams)}
	if err = tools.SetAttrs(d, attrs); err != nil {
		return diag.FromErr(err)
	}

	d.SetId(resID)
	return nil
}

func createStreamsAttrs(streams []datastream.StreamDetails) []interface{} {
	streamsAttrs := make([]interface{}, 0, len(streams))
	for _, stream := range streams {
		streamsAttrs = append(streamsAttrs, map[string]interface{}{
			"activation_status":  stream.ActivationStatus,
			"archived":           stream.Archived,
			"connectors":         stream.Connectors,
			"contract_id":        stream.ContractID,
			"created_by":         stream.CreatedBy,
			"created_date":       stream.CreatedDate,
			"current_version_id": stream.CurrentVersionID,
			"errors":             createErrorsAttrs(stream.Errors),
			"group_id":           stream.GroupID,
			"group_name":         stream.GroupName,
			"properties":         createPropertiesAttrs(stream.Properties),
			"stream_id":          stream.StreamID,
			"stream_name":        stream.StreamName,
			"stream_type_name":   stream.StreamTypeName,
			"stream_version_id":  stream.StreamVersionID,
		})
	}

	return streamsAttrs
}

func createErrorsAttrs(errors []datastream.Errors) []interface{} {
	errorAttrs := make([]interface{}, 0, len(errors))

	for _, errDetails := range errors {
		errorAttrs = append(errorAttrs, map[string]interface{}{
			"detail": errDetails.Detail,
			"title":  errDetails.Title,
			"type":   errDetails.Type,
		})
	}

	return errorAttrs
}

func createPropertiesAttrs(properties []datastream.Property) []interface{} {
	propertyAttrs := make([]interface{}, 0, len(properties))

	for _, property := range properties {
		propertyAttrs = append(propertyAttrs, map[string]interface{}{
			"property_id":   property.PropertyID,
			"property_name": property.PropertyName,
		})
	}

	return propertyAttrs
}
