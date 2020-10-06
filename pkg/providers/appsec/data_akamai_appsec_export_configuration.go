package appsec

import (
	"context"
	"strconv"

	"github.com/akamai/AkamaiOPEN-edgegrid-golang/jsonhooks-v1"
	v2 "github.com/akamai/AkamaiOPEN-edgegrid-golang/v2/pkg/appsec"
	"github.com/akamai/terraform-provider-akamai/v2/pkg/akamai"

	"github.com/hashicorp/terraform-plugin-sdk/v2/diag"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/schema"
)

func dataSourceExportConfiguration() *schema.Resource {
	return &schema.Resource{
		ReadContext: dataSourceExportConfigurationRead,
		Schema: map[string]*schema.Schema{
			"config_id": {
				Type:     schema.TypeInt,
				Required: true,
			},
			"version": {
				Type:     schema.TypeInt,
				Required: true,
			},
			"search": {
				Type:     schema.TypeList,
				Optional: true,
				Elem:     &schema.Schema{Type: schema.TypeString},
			},
			"json": {
				Type:        schema.TypeString,
				Computed:    true,
				Description: "JSON Export representation",
			},
			"output_text": {
				Type:        schema.TypeString,
				Computed:    true,
				Description: "Text Export representation",
			},
		},
	}
}

func dataSourceExportConfigurationRead(ctx context.Context, d *schema.ResourceData, m interface{}) diag.Diagnostics {
	meta := akamai.Meta(m)
	client := inst.Client(meta)
	logger := meta.Log("APPSEC", "resourceExportConfigurationRead")

	getExportConfiguration := v2.GetExportConfigurationsRequest{}

	getExportConfiguration.ConfigID = d.Get("config_id").(int)
	getExportConfiguration.Version = d.Get("version").(int)

	exportconfiguration, err := client.GetExportConfigurations(ctx, getExportConfiguration)
	if err != nil {
		logger.Warnf("calling 'getExportConfiguration': %s", err.Error())
	}

	jsonBody, err := jsonhooks.Marshal(exportconfiguration)
	if err != nil {
		return diag.FromErr(err)
	}

	d.Set("json", string(jsonBody))

	searchlist, ok := d.GetOk("search")
	if ok {
		ots := OutputTemplates{}
		InitTemplates(ots)

		var outputtextresult string

		for _, h := range searchlist.([]interface{}) {
			outputtext, err := RenderTemplates(ots, h.(string), exportconfiguration)
			if err == nil {
				outputtextresult = outputtextresult + outputtext
			}
		}

		if len(outputtextresult) > 0 {
			d.Set("output_text", outputtextresult)
		}
	}
	d.SetId(strconv.Itoa(exportconfiguration.ConfigID))

	return nil
}
