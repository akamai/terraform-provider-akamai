package appsec

import (
	"context"
	"encoding/json"
	"errors"
	"strconv"

	"github.com/akamai/AkamaiOPEN-edgegrid-golang/v2/pkg/appsec"
	"github.com/akamai/terraform-provider-akamai/v2/pkg/akamai"
	"github.com/akamai/terraform-provider-akamai/v2/pkg/tools"

	"github.com/hashicorp/terraform-plugin-sdk/v2/diag"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/schema"
)

func dataSourceAPIHostnameCoverageMatchTargets() *schema.Resource {
	return &schema.Resource{
		ReadContext: dataSourceAPIHostnameCoverageMatchTargetsRead,
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

func dataSourceAPIHostnameCoverageMatchTargetsRead(ctx context.Context, d *schema.ResourceData, m interface{}) diag.Diagnostics {
	meta := akamai.Meta(m)
	client := inst.Client(meta)
	logger := meta.Log("APPSEC", "dataSourceAPIHostnameCoverageMatchTargetsRead")

	getAPIHostnameCoverageMatchTargets := appsec.GetApiHostnameCoverageMatchTargetsRequest{}

	configID, err := tools.GetIntValue("config_id", d)
	if err != nil && !errors.Is(err, tools.ErrNotFound) {
		return diag.FromErr(err)
	}
	getAPIHostnameCoverageMatchTargets.ConfigID = configID

	getAPIHostnameCoverageMatchTargets.Version = getLatestConfigVersion(ctx, configID, m)

	hostname, err := tools.GetStringValue("hostname", d)
	if err != nil && !errors.Is(err, tools.ErrNotFound) {
		return diag.FromErr(err)
	}
	getAPIHostnameCoverageMatchTargets.Hostname = hostname

	apihostnamecoveragematchtargets, err := client.GetApiHostnameCoverageMatchTargets(ctx, getAPIHostnameCoverageMatchTargets)
	if err != nil {
		logger.Errorf("calling 'getAPIHostnameCoverageMatchTargets': %s", err.Error())
		return diag.FromErr(err)
	}

	ots := OutputTemplates{}
	InitTemplates(ots)

	outputtext, err := RenderTemplates(ots, "apiHostnameCoverageMatchTargetsDS", apihostnamecoveragematchtargets)
	if err == nil {
		d.Set("output_text", outputtext)
	}

	jsonBody, err := json.Marshal(apihostnamecoveragematchtargets)
	if err != nil {
		return diag.FromErr(err)
	}

	if err := d.Set("json", string(jsonBody)); err != nil {
		return diag.Errorf("%s: %s", tools.ErrValueSet, err.Error())
	}

	d.SetId(strconv.Itoa(getAPIHostnameCoverageMatchTargets.ConfigID))

	return nil
}
