package appsec

import (
	"context"
	"encoding/json"
	"errors"
	"fmt"
	"strconv"

	"github.com/akamai/AkamaiOPEN-edgegrid-golang/v2/pkg/appsec"
	"github.com/akamai/terraform-provider-akamai/v2/pkg/akamai"
	"github.com/akamai/terraform-provider-akamai/v2/pkg/tools"

	"github.com/hashicorp/terraform-plugin-sdk/v2/diag"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/schema"
)

func dataSourceApiHostnameCoverageOverlapping() *schema.Resource {
	return &schema.Resource{
		ReadContext: dataSourceApiHostnameCoverageOverlappingRead,
		Schema: map[string]*schema.Schema{
			"config_id": {
				Type:     schema.TypeInt,
				Required: true,
			},
			"hostname": {
				Type:     schema.TypeString,
				Required: true,
			},
			"json": {
				Type:     schema.TypeString,
				Computed: true,
			},
			"output_text": {
				Type:        schema.TypeString,
				Computed:    true,
				Description: "Text Export representation",
			},
		},
	}
}

func dataSourceApiHostnameCoverageOverlappingRead(ctx context.Context, d *schema.ResourceData, m interface{}) diag.Diagnostics {
	meta := akamai.Meta(m)
	client := inst.Client(meta)
	logger := meta.Log("APPSEC", "dataSourceApiHostnameCoverageOverlappingRead")

	getApiHostnameCoverageOverlapping := appsec.GetApiHostnameCoverageOverlappingRequest{}

	configid, err := tools.GetIntValue("config_id", d)
	if err != nil && !errors.Is(err, tools.ErrNotFound) {
		return diag.FromErr(err)
	}
	getApiHostnameCoverageOverlapping.ConfigID = configid

	getApiHostnameCoverageOverlapping.Version = getLatestConfigVersion(ctx, configid, m)

	hostname, err := tools.GetStringValue("hostname", d)
	if err != nil && !errors.Is(err, tools.ErrNotFound) {
		return diag.FromErr(err)
	}
	getApiHostnameCoverageOverlapping.Hostname = hostname

	apihostnamecoverageoverlapping, err := client.GetApiHostnameCoverageOverlapping(ctx, getApiHostnameCoverageOverlapping)
	if err != nil {
		logger.Errorf("calling 'getApiHostnameCoverageOverlapping': %s", err.Error())
		return diag.FromErr(err)
	}

	ots := OutputTemplates{}
	InitTemplates(ots)

	outputtext, err := RenderTemplates(ots, "apiHostnameCoverageoverLappingDS", apihostnamecoverageoverlapping)
	if err == nil {
		d.Set("output_text", outputtext)
	}

	jsonBody, err := json.Marshal(apihostnamecoverageoverlapping)
	if err != nil {
		return diag.FromErr(err)
	}

	if err := d.Set("json", string(jsonBody)); err != nil {
		return diag.FromErr(fmt.Errorf("%w: %s", tools.ErrValueSet, err.Error()))
	}

	d.SetId(strconv.Itoa(getApiHostnameCoverageOverlapping.ConfigID))

	return nil
}
