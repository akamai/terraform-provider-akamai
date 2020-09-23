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
			"active_in_staging": {
				Type:     schema.TypeBool,
				Optional: true,
			},
			"active_in_production": {
				Type:     schema.TypeBool,
				Optional: true,
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

	var flagsetstg, flagsetprod string

	flagsetstg = "UNSET"
	flagsetprod = "UNSET"

	activeinstaging, ok := d.GetOkExists("active_in_staging")
	if ok {
		flagsetstg = "SET"
	}
	activeinproduction, ok := d.GetOkExists("active_in_production")
	if ok {
		flagsetprod = "SET"
	}

	newhdata := make([]string, 0, len(selectablehostnames.AvailableSet))
	for _, hosts := range selectablehostnames.AvailableSet {
		var flagstg, flagprod string
		flagstg = "NOMATCH"
		flagprod = "NOMATCH"

		if activeinstaging.(bool) == hosts.ActiveInStaging {
			flagstg = "MATCH"
		} else {
			flagstg = "NOMATCH"
		}

		if activeinproduction.(bool) == hosts.ActiveInProduction {
			flagprod = "MATCH"
		} else {
			flagprod = "NOMATCH"
		}

		if flagstg == "MATCH" && flagprod == "MATCH" {
			newhdata = append(newhdata, hosts.Hostname)
		}

		if flagsetstg == "UNSET" && flagsetprod == "UNSET" {
			newhdata = append(newhdata, hosts.Hostname)
		}

	}

	//edge.PrintfCorrelation("[DEBUG]", CorrelationID, fmt.Sprintf("  SET SelectedHostnames H %v", h))
	d.Set("hostnames", newhdata)

	d.SetId(strconv.Itoa(selectablehostnames.ConfigID))

	return nil
}
