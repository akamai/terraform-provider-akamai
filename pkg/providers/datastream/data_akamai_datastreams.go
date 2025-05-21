package datastream

import (
	"context"
	"errors"
	"fmt"

	"github.com/akamai/AkamaiOPEN-edgegrid-golang/v10/pkg/session"
	"github.com/akamai/terraform-provider-akamai/v8/pkg/common/ptr"
	"github.com/akamai/terraform-provider-akamai/v8/pkg/common/tf"
	"github.com/akamai/terraform-provider-akamai/v8/pkg/meta"

	"github.com/akamai/AkamaiOPEN-edgegrid-golang/v10/pkg/datastream"
	"github.com/hashicorp/terraform-plugin-sdk/v2/diag"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/schema"
)

func dataAkamaiDatastreamStreams() *schema.Resource {
	return &schema.Resource{
		Description: "Retrieves a list of all streams optionally by the specified GroupID.",
		ReadContext: dataDatastreamStreamsRead,
		Schema: map[string]*schema.Schema{
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

	req := datastream.ListStreamsRequest{}
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
		streamsAttrs = append(streamsAttrs, map[string]interface{}{
			"stream_status":  stream.StreamStatus,
			"contract_id":    stream.ContractID,
			"created_by":     stream.CreatedBy,
			"created_date":   stream.CreatedDate,
			"modified_by":    stream.ModifiedBy,
			"modified_date":  stream.ModifiedDate,
			"group_id":       stream.GroupID,
			"latest_version": stream.LatestVersion,
			"product_id":     stream.ProductID,
			"properties":     createPropertiesAttrs(stream.Properties),
			"stream_id":      stream.StreamID,
			"stream_name":    stream.StreamName,
			"stream_version": stream.StreamVersion,
		})
	}

	return streamsAttrs
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
