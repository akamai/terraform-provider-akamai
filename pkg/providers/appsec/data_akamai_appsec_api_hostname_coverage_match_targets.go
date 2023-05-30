package appsec

import (
	"context"
	"encoding/json"
	"strconv"

	"github.com/akamai/AkamaiOPEN-edgegrid-golang/v6/pkg/appsec"
	"github.com/akamai/terraform-provider-akamai/v4/pkg/akamai"
	"github.com/akamai/terraform-provider-akamai/v4/pkg/common/tf"
	"github.com/hashicorp/terraform-plugin-sdk/v2/diag"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/schema"
)

func dataSourceAPIHostnameCoverageMatchTargets() *schema.Resource {
	return &schema.Resource{
		ReadContext: dataSourceAPIHostnameCoverageMatchTargetsRead,
		Schema: map[string]*schema.Schema{
			"config_id": {
				Type:        schema.TypeInt,
				Required:    true,
				Description: "Unique identifier of the security configuration",
			},
			"hostname": {
				Type:        schema.TypeString,
				Required:    true,
				Description: "Hostname for which to return match target information",
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

func dataSourceAPIHostnameCoverageMatchTargetsRead(ctx context.Context, d *schema.ResourceData, m interface{}) diag.Diagnostics {
	meta := akamai.Meta(m)
	client := inst.Client(meta)
	logger := meta.Log("APPSEC", "dataSourceAPIHostnameCoverageMatchTargetsRead")

	getAPIHostnameCoverageMatchTargets := appsec.GetApiHostnameCoverageMatchTargetsRequest{}

	configID, err := tf.GetIntValue("config_id", d)
	if err != nil {
		return diag.FromErr(err)
	}
	getAPIHostnameCoverageMatchTargets.ConfigID = configID

	if getAPIHostnameCoverageMatchTargets.Version, err = getLatestConfigVersion(ctx, configID, m); err != nil {
		return diag.FromErr(err)
	}

	hostname, err := tf.GetStringValue("hostname", d)
	if err != nil {
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
	if err != nil {
		return diag.FromErr(err)
	}
	if err := d.Set("output_text", outputtext); err != nil {
		return diag.Errorf("%s: %s", tf.ErrValueSet, err.Error())
	}

	jsonBody, err := json.Marshal(apihostnamecoveragematchtargets)
	if err != nil {
		return diag.FromErr(err)
	}

	if err := d.Set("json", string(jsonBody)); err != nil {
		return diag.Errorf("%s: %s", tf.ErrValueSet, err.Error())
	}

	d.SetId(strconv.Itoa(getAPIHostnameCoverageMatchTargets.ConfigID))

	return nil
}
