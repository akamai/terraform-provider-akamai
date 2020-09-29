package appsec

import (
	"fmt"
	"strconv"

	appsec "github.com/akamai/AkamaiOPEN-edgegrid-golang/appsec-v1"
	edge "github.com/akamai/AkamaiOPEN-edgegrid-golang/edgegrid"
	"github.com/akamai/AkamaiOPEN-edgegrid-golang/jsonhooks-v1"
	"github.com/akamai/terraform-provider-akamai/v2/pkg/tools"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/schema"
)

func dataSourceSelectableHostnames() *schema.Resource {
	return &schema.Resource{
		Read: dataSourceSelectableHostnamesRead,
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

func dataSourceSelectableHostnamesRead(d *schema.ResourceData, meta interface{}) error {
	CorrelationID := "[APPSEC][dataSourceSelectableHostnamesRead-" + tools.CreateNonce() + "]"

	edge.PrintfCorrelation("[DEBUG]", CorrelationID, "  Read SelectableHostnames")

	selectablehostnames := appsec.NewSelectableHostnamesResponse()
	selectablehostnames.ConfigID = d.Get("config_id").(int)
	selectablehostnames.ConfigVersion = d.Get("version").(int)

	err := selectablehostnames.GetSelectableHostnames(CorrelationID)
	if err != nil {
		edge.PrintfCorrelation("[DEBUG]", CorrelationID, fmt.Sprintf("Error  %v\n", err))
		return nil
	}

	edge.PrintfCorrelation("[DEBUG]", CorrelationID, fmt.Sprintf("SelectableHostnames   %v\n", selectablehostnames))

	jsonBody, err := jsonhooks.Marshal(selectablehostnames)
	if err != nil {
		return err
	}

	d.Set("hostnames_json", string(jsonBody))

	newhdata := make([]string, 0, len(selectablehostnames.AvailableSet))
	for _, hosts := range selectablehostnames.AvailableSet {

		newhdata = append(newhdata, hosts.Hostname)
	}

	//edge.PrintfCorrelation("[DEBUG]", CorrelationID, fmt.Sprintf("  SET SelectedHostnames H %v", h))
	d.Set("hostnames", newhdata)

	ots := OutputTemplates{}
	InitTemplates(ots)

	outputtext, err := RenderTemplates(ots, "selectableHostsDS", selectablehostnames)
	edge.PrintfCorrelation("[DEBUG]", CorrelationID, fmt.Sprintf("selectablehostnames outputtext   %v\n", outputtext))
	if err == nil {
		d.Set("output_text", outputtext)
	}

	d.SetId(strconv.Itoa(selectablehostnames.ConfigID))

	return nil
}
