package botman

import (
	"context"
	"strconv"

	"github.com/akamai/AkamaiOPEN-edgegrid-golang/v15/pkg/botman"
	"github.com/akamai/terraform-provider-akamai/v12/pkg/common/tf"
	"github.com/akamai/terraform-provider-akamai/v12/pkg/meta"
	"github.com/hashicorp/terraform-plugin-sdk/v2/diag"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/schema"
)

func dataSourceCustomClientSequence() *schema.Resource {
	return &schema.Resource{
		ReadContext: dataSourceCustomClientSequenceRead,
		Schema: map[string]*schema.Schema{
			"config_id": {
				Type:     schema.TypeInt,
				Required: true,
			},
			"custom_client_ids": {
				Type:     schema.TypeList,
				Computed: true,
				Elem:     &schema.Schema{Type: schema.TypeString},
			},
		},
	}
}

func dataSourceCustomClientSequenceRead(ctx context.Context, d *schema.ResourceData, m interface{}) diag.Diagnostics {
	meta := meta.Must(m)
	client := meta.Client().GetBotMan()
	logger := meta.Log("botman", "dataSourceCustomClientSequenceRead")
	logger.Debugf("in dataSourceCustomClientSequenceRead")

	configID, err := tf.GetIntValueAsInt64("config_id", d)
	if err != nil {
		return diag.FromErr(err)
	}

	version, err := getLatestConfigVersion(ctx, configID, meta.Client().GetAPPSEC())
	if err != nil {
		return diag.FromErr(err)
	}

	request := botman.GetCustomClientSequenceRequest{
		ConfigID: configID,
		Version:  version,
	}

	response, err := client.GetCustomClientSequence(ctx, request)
	if err != nil {
		logger.Errorf("calling 'GetCustomClientSequence': %s", err.Error())
		return diag.FromErr(err)
	}

	if err := d.Set("custom_client_ids", response.Sequence); err != nil {
		return diag.Errorf("%s: %s", tf.ErrValueSet, err.Error())
	}

	d.SetId(strconv.FormatInt(configID, 10))
	return nil
}
