package appsec

import (
	"context"
	"encoding/json"
	"strconv"

	"github.com/akamai/AkamaiOPEN-edgegrid-golang/v7/pkg/appsec"
	"github.com/akamai/terraform-provider-akamai/v5/pkg/common/tf"
	"github.com/akamai/terraform-provider-akamai/v5/pkg/meta"
	"github.com/hashicorp/terraform-plugin-sdk/v2/diag"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/schema"
)

func dataSourceSelectedHostnames() *schema.Resource {
	return &schema.Resource{
		ReadContext: dataSourceSelectedHostnamesRead,
		Schema: map[string]*schema.Schema{
			"config_id": {
				Type:        schema.TypeInt,
				Required:    true,
				Description: "Unique identifier of the security configuration",
			},
			"hostnames": {
				Type:        schema.TypeList,
				Computed:    true,
				Elem:        &schema.Schema{Type: schema.TypeString},
				Description: "List of hostnames",
			},
			"hostnames_json": {
				Type:        schema.TypeString,
				Computed:    true,
				Description: "JSON List of hostnames",
			},
			"output_text": {
				Type:        schema.TypeString,
				Computed:    true,
				Description: "Text representation",
			},
		},
	}
}

func dataSourceSelectedHostnamesRead(ctx context.Context, d *schema.ResourceData, m interface{}) diag.Diagnostics {
	meta := meta.Must(m)
	client := inst.Client(meta)
	logger := meta.Log("APPSEC", "dataSourceSelectedHostnamesRead")

	getSelectedHostnames := appsec.GetSelectedHostnamesRequest{}

	configID, err := tf.GetIntValue("config_id", d)
	if err != nil {
		return diag.FromErr(err)
	}
	getSelectedHostnames.ConfigID = configID

	if getSelectedHostnames.Version, err = getLatestConfigVersion(ctx, configID, m); err != nil {
		return diag.FromErr(err)
	}

	selectedhostnames, err := client.GetSelectedHostnames(ctx, getSelectedHostnames)
	if err != nil {
		logger.Errorf("calling 'getSelectedHostnames': %s", err.Error())
		return diag.FromErr(err)
	}

	jsonBody, err := json.Marshal(selectedhostnames)
	if err != nil {
		return diag.FromErr(err)
	}

	if err := d.Set("hostnames_json", string(jsonBody)); err != nil {
		return diag.Errorf("%s: %s", tf.ErrValueSet, err.Error())
	}

	newhdata := make([]string, 0, len(selectedhostnames.HostnameList))

	for _, hosts := range selectedhostnames.HostnameList {
		newhdata = append(newhdata, hosts.Hostname)
	}

	if err := d.Set("hostnames", newhdata); err != nil {
		return diag.Errorf("%s: %s", tf.ErrValueSet, err.Error())
	}

	if err := d.Set("config_id", getSelectedHostnames.ConfigID); err != nil {
		return diag.Errorf("%s: %s", tf.ErrValueSet, err.Error())
	}

	ots := OutputTemplates{}
	InitTemplates(ots)

	outputtext, err := RenderTemplates(ots, "selectedHostsDS", selectedhostnames)
	if err != nil {
		return diag.FromErr(err)
	}
	if err := d.Set("output_text", outputtext); err != nil {
		return diag.Errorf("%s: %s", tf.ErrValueSet, err.Error())
	}

	d.SetId(strconv.Itoa(getSelectedHostnames.ConfigID))

	return nil
}
