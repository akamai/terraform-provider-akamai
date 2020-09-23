package appsec

import (
	"fmt"
	"strconv"
	"strings"

	appsec "github.com/akamai/AkamaiOPEN-edgegrid-golang/appsec-v1"
	edge "github.com/akamai/AkamaiOPEN-edgegrid-golang/edgegrid"
	"github.com/akamai/AkamaiOPEN-edgegrid-golang/jsonhooks-v1"
	"github.com/akamai/terraform-provider-akamai/v2/pkg/tools"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/schema"
)

func dataSourceSelectedHostnames() *schema.Resource {
	return &schema.Resource{
		Read: dataSourceSelectedHostnamesRead,
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
		},
	}
}

func dataSourceSelectedHostnamesRead(d *schema.ResourceData, meta interface{}) error {
	CorrelationID := "[APPSEC][dataSourceSelectedHostnamesRead-" + tools.CreateNonce() + "]"

	edge.PrintfCorrelation("[DEBUG]", CorrelationID, "  Read SelectedHostnames")

	selectedhostnames := appsec.NewSelectedHostnamesResponse()
	var configid int
	var version int
	edge.PrintfCorrelation("[DEBUG]", CorrelationID, fmt.Sprintf("  Read SelectedHostnames D.ID %v", d.Id()))

	if d.Id() != "" && strings.Contains(d.Id(), ":") {
		s := strings.Split(d.Id(), ":")
		configid, _ = strconv.Atoi(s[0])
		version, _ = strconv.Atoi(s[1])
	} else {
		configid = d.Get("config_id").(int)
		version = d.Get("version").(int)
	}

	err := selectedhostnames.GetSelectedHostnames(configid, version, CorrelationID)
	if err != nil {
		edge.PrintfCorrelation("[DEBUG]", CorrelationID, fmt.Sprintf("Error  %v\n", err))
		return nil
	}

	edge.PrintfCorrelation("[DEBUG]", CorrelationID, fmt.Sprintf("SelectedHostnames   %v\n", selectedhostnames))

	jsonBody, err := jsonhooks.Marshal(selectedhostnames)
	if err != nil {
		return err
	}

	d.Set("hostnames_json", string(jsonBody))

	newhdata := make([]string, 0, len(selectedhostnames.HostnameList))
	for _, hosts := range selectedhostnames.HostnameList {
		newhdata = append(newhdata, hosts.Hostname)
	}

	d.Set("hostnames", newhdata)
	d.Set("config_id", configid)
	d.Set("version", version)
	d.SetId(fmt.Sprintf("%d:%d", configid, version))

	return nil
}
