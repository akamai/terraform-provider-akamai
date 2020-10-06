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

func dataSourceSelectableHostnames() *schema.Resource {
	return &schema.Resource{
		ReadContext: dataSourceSelectableHostnamesRead,
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
			"output_text": {
				Type:        schema.TypeString,
				Computed:    true,
				Description: "Text Export representation",
			},
		},
	}
}

func dataSourceSelectableHostnamesRead(ctx context.Context, d *schema.ResourceData, m interface{}) diag.Diagnostics {
	meta := akamai.Meta(m)
	client := inst.Client(meta)
	logger := meta.Log("APPSEC", "resourceSelectableHostnamesRead")

	getSelectableHostnames := v2.GetSelectableHostnamesRequest{}

	getSelectableHostnames.ConfigID = d.Get("config_id").(int)
	getSelectableHostnames.Version = d.Get("version").(int)

	selectablehostnames, err := client.GetSelectableHostnames(ctx, getSelectableHostnames)
	if err != nil {
		logger.Warnf("calling 'getSelectableHostnames': %s", err.Error())
	}

	jsonBody, err := jsonhooks.Marshal(selectablehostnames)
	if err != nil {
		return diag.FromErr(err)
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

	ots := OutputTemplates{}
	InitTemplates(ots)

	outputtext, err := RenderTemplates(ots, "selectableHostsDS", selectablehostnames)
	if err == nil {
		d.Set("output_text", outputtext)
	}

	d.SetId(strconv.Itoa(selectablehostnames.ConfigID))

	return nil
}
