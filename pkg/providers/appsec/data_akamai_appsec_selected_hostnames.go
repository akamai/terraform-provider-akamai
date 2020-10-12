package appsec

import (
	"context"
	"errors"
	"fmt"
	"strconv"
	"strings"

	edge "github.com/akamai/AkamaiOPEN-edgegrid-golang/edgegrid"
	"github.com/akamai/AkamaiOPEN-edgegrid-golang/jsonhooks-v1"
	v2 "github.com/akamai/AkamaiOPEN-edgegrid-golang/v2/pkg/appsec"
	"github.com/akamai/terraform-provider-akamai/v2/pkg/akamai"
	"github.com/akamai/terraform-provider-akamai/v2/pkg/tools"
	"github.com/hashicorp/terraform-plugin-sdk/v2/diag"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/schema"
)

func dataSourceSelectedHostnames() *schema.Resource {
	return &schema.Resource{
		ReadContext: dataSourceSelectedHostnamesRead,
		Schema: map[string]*schema.Schema{
			"config_id": {
				Type:     schema.TypeInt,
				Required: true,
			},
			"version": {
				Type:     schema.TypeInt,
				Required: true,
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
				Description: "Text Export representation",
			},
		},
	}
}

func dataSourceSelectedHostnamesRead(ctx context.Context, d *schema.ResourceData, m interface{}) diag.Diagnostics {
	meta := akamai.Meta(m)
	client := inst.Client(meta)
	logger := meta.Log("APPSEC", "resourceSelectedHostnamesRead")
	CorrelationID := "[APPSEC][resourceSelectedHostnames-" + meta.OperationID() + "]"

	getSelectedHostnames := v2.GetSelectedHostnamesRequest{}

	edge.PrintfCorrelation("[DEBUG]", CorrelationID, fmt.Sprintf("  Read SelectedHostnames D.ID %v", d.Id()))

	if d.Id() != "" && strings.Contains(d.Id(), ":") {
		s := strings.Split(d.Id(), ":")
		getSelectedHostnames.ConfigID, _ = strconv.Atoi(s[0])
		getSelectedHostnames.Version, _ = strconv.Atoi(s[1])
	} else {
		configid, err := tools.GetIntValue("config_id", d)
		if err != nil && !errors.Is(err, tools.ErrNotFound) {
			return diag.FromErr(err)
		}
		getSelectedHostnames.ConfigID = configid

		version, err := tools.GetIntValue("version", d)
		if err != nil && !errors.Is(err, tools.ErrNotFound) {
			return diag.FromErr(err)
		}
		getSelectedHostnames.Version = version
	}

	selectedhostnames, err := client.GetSelectedHostnames(ctx, getSelectedHostnames)
	if err != nil {
		logger.Warnf("calling 'getSelectedHostnames': %s", err.Error())
	}

	jsonBody, err := jsonhooks.Marshal(selectedhostnames)
	if err != nil {
		return diag.FromErr(err)
	}

	d.Set("hostnames_json", string(jsonBody))

	newhdata := make([]string, 0, len(selectedhostnames.HostnameList))

	for _, hosts := range selectedhostnames.HostnameList {
		newhdata = append(newhdata, hosts.Hostname)
	}

	d.Set("hostnames", newhdata)
	d.Set("config_id", getSelectedHostnames.ConfigID)
	d.Set("version", getSelectedHostnames.Version)

	ots := OutputTemplates{}
	InitTemplates(ots)

	outputtext, err := RenderTemplates(ots, "selectedHostsDS", selectedhostnames)
	edge.PrintfCorrelation("[DEBUG]", CorrelationID, fmt.Sprintf("selectedhostnames outputtext   %v\n", outputtext))
	if err == nil {
		d.Set("output_text", outputtext)
	}

	d.SetId(fmt.Sprintf("%d:%d", getSelectedHostnames.ConfigID, getSelectedHostnames.Version))

	return nil
}
