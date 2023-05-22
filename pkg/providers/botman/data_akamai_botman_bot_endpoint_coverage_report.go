package botman

import (
	"context"
	"encoding/json"
	"errors"
	"fmt"

	"github.com/akamai/AkamaiOPEN-edgegrid-golang/v6/pkg/botman"
	"github.com/akamai/terraform-provider-akamai/v3/pkg/akamai"
	"github.com/akamai/terraform-provider-akamai/v3/pkg/common/tf"
	"github.com/akamai/terraform-provider-akamai/v3/pkg/tools"
	"github.com/hashicorp/terraform-plugin-sdk/v2/diag"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/schema"
)

func dataSourceBotEndpointCoverageReport() *schema.Resource {
	return &schema.Resource{
		ReadContext: dataSourceBotEndpointCoverageReportRead,
		Schema: map[string]*schema.Schema{
			"config_id": {
				Type:     schema.TypeInt,
				Optional: true,
			},
			"operation_id": {
				Type:     schema.TypeString,
				Optional: true,
			},
			"json": {
				Type:     schema.TypeString,
				Computed: true,
			},
		},
	}
}

func dataSourceBotEndpointCoverageReportRead(ctx context.Context, d *schema.ResourceData, m interface{}) diag.Diagnostics {
	meta := akamai.Meta(m)
	client := inst.Client(meta)
	logger := meta.Log("botman", "dataSourceBotEndpointCoverageReportRead")
	logger.Debugf("in dataSourceBotEndpointCoverageReportRead")

	configID, err := tf.GetIntValue("config_id", d)
	if err != nil && !errors.Is(err, tf.ErrNotFound) {
		return diag.FromErr(err)
	}

	var version int
	if configID != 0 {
		version, err = getLatestConfigVersion(ctx, configID, m)
		if err != nil {
			return diag.FromErr(err)
		}
	}

	operationID, err := tf.GetStringValue("operation_id", d)
	if err != nil && !errors.Is(err, tf.ErrNotFound) {
		return diag.FromErr(err)
	}

	request := botman.GetBotEndpointCoverageReportRequest{
		ConfigID:    int64(configID),
		Version:     int64(version),
		OperationID: operationID,
	}

	response, err := client.GetBotEndpointCoverageReport(ctx, request)
	if err != nil {
		logger.Errorf("calling 'GetBotEndpointCoverageReport': %s", err.Error())
		return diag.FromErr(err)
	}

	jsonBody, err := json.Marshal(response)
	if err != nil {
		return diag.FromErr(err)
	}
	if err := d.Set("json", string(jsonBody)); err != nil {
		return diag.Errorf("%s: %s", tf.ErrValueSet, err.Error())
	}

	if configID != 0 {
		d.SetId(fmt.Sprintf("%d", configID))
	} else {
		d.SetId(tools.GetSHAString(string(jsonBody)))
	}
	return nil
}
