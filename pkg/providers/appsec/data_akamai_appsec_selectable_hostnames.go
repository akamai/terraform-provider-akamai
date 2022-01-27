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

func dataSourceSelectableHostnames() *schema.Resource {
	return &schema.Resource{
		ReadContext: dataSourceSelectableHostnamesRead,
		Schema: map[string]*schema.Schema{
			"config_id": {
				Type:          schema.TypeInt,
				Optional:      true,
				ConflictsWith: []string{"contractid", "groupid"},
			},
			"contractid": {
				Type:          schema.TypeString,
				Optional:      true,
				ConflictsWith: []string{"config_id"},
			},
			"groupid": {
				Type:          schema.TypeInt,
				Optional:      true,
				ConflictsWith: []string{"config_id"},
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
	logger := meta.Log("APPSEC", "dataSourceSelectableHostnamesRead")

	getSelectableHostnames := appsec.GetSelectableHostnamesRequest{}

	configID, err := tools.GetIntValue("config_id", d)
	if err != nil && !errors.Is(err, tools.ErrNotFound) {
		return diag.FromErr(err)
	}
	getSelectableHostnames.ConfigID = configID

	getSelectableHostnames.Version = getLatestConfigVersion(ctx, configID, m)

	contractID, err := tools.GetStringValue("contractid", d)
	if err != nil && !errors.Is(err, tools.ErrNotFound) {
		return diag.FromErr(err)
	}
	getSelectableHostnames.ContractID = contractID

	group, err := tools.GetIntValue("groupid", d)
	if err != nil && !errors.Is(err, tools.ErrNotFound) {
		return diag.FromErr(err)
	}
	getSelectableHostnames.GroupID = group

	selectablehostnames, err := client.GetSelectableHostnames(ctx, getSelectableHostnames)
	if err != nil {
		logger.Errorf("calling 'getSelectableHostnames': %s", err.Error())
		return diag.FromErr(err)
	}

	jsonBody, err := json.Marshal(selectablehostnames)
	if err != nil {
		return diag.FromErr(err)
	}

	if err := d.Set("hostnames_json", string(jsonBody)); err != nil {
		return diag.Errorf("%s: %s", tools.ErrValueSet, err.Error())
	}

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

	if err := d.Set("hostnames", newhdata); err != nil {
		return diag.Errorf("%s: %s", tools.ErrValueSet, err.Error())
	}

	ots := OutputTemplates{}
	InitTemplates(ots)

	outputtext, err := RenderTemplates(ots, "selectableHostsDS", selectablehostnames)
	if err == nil {
		if err := d.Set("output_text", outputtext); err != nil {
			return diag.Errorf("%s: %s", tools.ErrValueSet, err.Error())
		}
	}

	d.SetId(strconv.Itoa(selectablehostnames.ConfigID))

	return nil
}
