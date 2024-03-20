package datastream

import (
	"context"
	"fmt"

	"github.com/akamai/terraform-provider-akamai/v6/pkg/common/tf"
	"github.com/akamai/terraform-provider-akamai/v6/pkg/meta"
	"github.com/hashicorp/terraform-plugin-sdk/v2/diag"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/schema"

	"github.com/akamai/AkamaiOPEN-edgegrid-golang/v8/pkg/datastream"
)

func dataAkamaiDatastreamActivationHistory() *schema.Resource {
	return &schema.Resource{
		ReadContext: dataAkamaiDatastreamActivationHistoryRead,
		Schema: map[string]*schema.Schema{
			"stream_id": {
				Type:        schema.TypeInt,
				Required:    true,
				Description: "Identifies the stream",
			},
			"activations": {
				Type:        schema.TypeList,
				Computed:    true,
				Description: "Provides detailed information about an activation history for a version of a stream",
				Elem: &schema.Resource{
					Schema: map[string]*schema.Schema{

						"stream_id": {
							Type:        schema.TypeInt,
							Required:    true,
							Description: "Identifies the stream",
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
						"stream_version": {
							Type:        schema.TypeInt,
							Computed:    true,
							Description: "Identifies the version of the stream",
						},
						"status": {
							Type:        schema.TypeString,
							Computed:    true,
							Description: "Stream Status",
						},
					},
				},
			},
		},
	}
}

func populateSchemaFieldsWithActivationHistory(ac []datastream.ActivationHistoryEntry, d *schema.ResourceData) error {

	var activations []map[string]interface{}
	for _, a := range ac {
		v := map[string]interface{}{
			"stream_id":      a.StreamID,
			"stream_version": a.StreamVersion,
			"modified_by":    a.ModifiedBy,
			"modified_date":  a.ModifiedDate,
			"status":         a.Status,
		}
		activations = append(activations, v)
	}

	fields := map[string]interface{}{
		"activations": activations,
	}

	err := tf.SetAttrs(d, fields)
	if err != nil {
		return fmt.Errorf("could not set schema attributes: %s", err)
	}

	return nil
}

func dataAkamaiDatastreamActivationHistoryRead(ctx context.Context, d *schema.ResourceData, m interface{}) diag.Diagnostics {
	meta := meta.Must(m)
	log := meta.Log("DataStream", "dataAkamaiDatastreamActivationHistoryRead")
	client := inst.Client(meta)

	streamID, err := tf.GetIntValue("stream_id", d)
	if err != nil {
		return diag.FromErr(err)
	}

	log.Debug("Getting activation history")
	activationHistory, err := client.GetActivationHistory(ctx, datastream.GetActivationHistoryRequest{
		StreamID: int64(streamID),
	})
	if err != nil {
		return diag.FromErr(err)
	}

	err = populateSchemaFieldsWithActivationHistory(activationHistory, d)
	if err != nil {
		return diag.FromErr(err)
	}

	d.SetId(fmt.Sprintf("%d", streamID))

	return nil
}
