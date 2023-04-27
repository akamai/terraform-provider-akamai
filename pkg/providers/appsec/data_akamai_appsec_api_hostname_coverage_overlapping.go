package appsec

import (
	"context"
	"encoding/json"
	"strconv"

	"github.com/akamai/AkamaiOPEN-edgegrid-golang/v5/pkg/appsec"
	"github.com/akamai/terraform-provider-akamai/v3/pkg/akamai"
	"github.com/akamai/terraform-provider-akamai/v3/pkg/common/tf"
	"github.com/hashicorp/terraform-plugin-sdk/v2/diag"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/schema"
)

func dataSourceAPIHostnameCoverageOverlapping() *schema.Resource {
	return &schema.Resource{
		ReadContext: dataSourceAPIHostnameCoverageOverlappingRead,
		Schema: map[string]*schema.Schema{
			"config_id": {
				Type:        schema.TypeInt,
				Required:    true,
				Description: "Unique identifier of the security configuration",
			},
			"hostname": {
				Type:        schema.TypeString,
				Required:    true,
				Description: "Hostname for which to return coverage overlap information",
			},
			"json": {
				Type:        schema.TypeString,
				Computed:    true,
				Description: "JSON representation",
			},
			"output_text": {
				Type:        schema.TypeString,
				Computed:    true,
				Description: "Text representation",
			},
		},
	}
}

func dataSourceAPIHostnameCoverageOverlappingRead(ctx context.Context, d *schema.ResourceData, m interface{}) diag.Diagnostics {
	meta := akamai.Meta(m)
	client := inst.Client(meta)
	logger := meta.Log("APPSEC", "dataSourceAPIHostnameCoverageOverlappingRead")

	getAPIHostnameCoverageOverlapping := appsec.GetApiHostnameCoverageOverlappingRequest{}

	configID, err := tf.GetIntValue("config_id", d)
	if err != nil {
		return diag.FromErr(err)
	}
	getAPIHostnameCoverageOverlapping.ConfigID = configID

	if getAPIHostnameCoverageOverlapping.Version, err = getLatestConfigVersion(ctx, configID, m); err != nil {
		return diag.FromErr(err)
	}

	hostname, err := tf.GetStringValue("hostname", d)
	if err != nil {
		return diag.FromErr(err)
	}
	getAPIHostnameCoverageOverlapping.Hostname = hostname

	apihostnamecoverageoverlapping, err := client.GetApiHostnameCoverageOverlapping(ctx, getAPIHostnameCoverageOverlapping)
	if err != nil {
		logger.Errorf("calling 'getApiHostnameCoverageOverlapping': %s", err.Error())
		return diag.FromErr(err)
	}

	ots := OutputTemplates{}
	InitTemplates(ots)

	outputtext, err := RenderTemplates(ots, "apiHostnameCoverageoverLappingDS", apihostnamecoverageoverlapping)
	if err != nil {
		return diag.FromErr(err)
	}
	if err := d.Set("output_text", outputtext); err != nil {
		return diag.Errorf("%s: %s", tf.ErrValueSet, err.Error())
	}

	jsonBody, err := json.Marshal(apihostnamecoverageoverlapping)
	if err != nil {
		return diag.FromErr(err)
	}

	if err := d.Set("json", string(jsonBody)); err != nil {
		return diag.Errorf("%s: %s", tf.ErrValueSet, err.Error())
	}

	d.SetId(strconv.Itoa(getAPIHostnameCoverageOverlapping.ConfigID))

	return nil
}
