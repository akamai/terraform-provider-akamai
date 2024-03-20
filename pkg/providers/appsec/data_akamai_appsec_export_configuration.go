package appsec

import (
	"context"
	"encoding/json"
	"strconv"

	"github.com/akamai/AkamaiOPEN-edgegrid-golang/v8/pkg/appsec"
	"github.com/akamai/terraform-provider-akamai/v5/pkg/common/tf"
	"github.com/akamai/terraform-provider-akamai/v5/pkg/meta"
	"github.com/hashicorp/terraform-plugin-sdk/v2/diag"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/schema"
)

func dataSourceExportConfiguration() *schema.Resource {
	return &schema.Resource{
		ReadContext: dataSourceExportConfigurationRead,
		Schema: map[string]*schema.Schema{
			"config_id": {
				Type:        schema.TypeInt,
				Required:    true,
				Description: "Unique identifier of the security configuration",
			},
			"version": {
				Type:        schema.TypeInt,
				Required:    true,
				Description: "Version number of the security configuration to be exported",
			},
			"search": {
				Type:        schema.TypeList,
				Optional:    true,
				Elem:        &schema.Schema{Type: schema.TypeString},
				Description: "List of template files indicating resources to be exported for later import",
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

func dataSourceExportConfigurationRead(ctx context.Context, d *schema.ResourceData, m interface{}) diag.Diagnostics {
	meta := meta.Must(m)
	client := inst.Client(meta)
	logger := meta.Log("APPSEC", "dataSourceExportConfigurationRead")

	configID, err := tf.GetIntValue("config_id", d)
	if err != nil {
		return diag.FromErr(err)
	}
	version, err := tf.GetIntValue("version", d)
	if err != nil {
		return diag.FromErr(err)
	}

	exportconfiguration, err := client.GetExportConfiguration(ctx, appsec.GetExportConfigurationRequest{ConfigID: configID, Version: version})
	if err != nil {
		logger.Errorf("calling 'getExportConfiguration': %s", err.Error())
		return diag.FromErr(err)
	}

	jsonBody, err := json.Marshal(exportconfiguration)
	if err != nil {
		return diag.FromErr(err)
	}

	if err := d.Set("json", string(jsonBody)); err != nil {
		return diag.Errorf("%s: %s", tf.ErrValueSet, err.Error())
	}

	searchlist, ok := d.GetOk("search")
	if ok {
		ots := OutputTemplates{}
		InitTemplates(ots)

		var outputtextresult string

		for _, h := range searchlist.([]interface{}) {
			outputtext, err := RenderTemplates(ots, h.(string), exportconfiguration)
			if err != nil {
				return diag.FromErr(err)
			}
			outputtextresult = outputtextresult + outputtext
		}

		if len(outputtextresult) > 0 {
			if err := d.Set("output_text", outputtextresult); err != nil {
				return diag.Errorf("%s: %s", tf.ErrValueSet, err.Error())
			}
		}
	}
	d.SetId(strconv.Itoa(exportconfiguration.ConfigID))

	return nil
}
