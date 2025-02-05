package appsec

import (
	"context"
	"encoding/json"
	"errors"
	"strconv"

	"github.com/akamai/AkamaiOPEN-edgegrid-golang/v10/pkg/appsec"
	"github.com/akamai/terraform-provider-akamai/v7/pkg/common/tf"
	"github.com/akamai/terraform-provider-akamai/v7/pkg/meta"
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
				Description:   "Unique identifier of the security configuration",
				ConflictsWith: []string{"contractid", "groupid"},
			},
			"contractid": {
				Type:          schema.TypeString,
				Optional:      true,
				Description:   "Unique identifier of an Akamai contract",
				ConflictsWith: []string{"config_id"},
			},
			"groupid": {
				Type:          schema.TypeInt,
				Optional:      true,
				Description:   "Unique identifier of a contract group",
				ConflictsWith: []string{"config_id"},
			},
			"active_in_staging": {
				Type:        schema.TypeBool,
				Optional:    true,
				Description: "Whether to return names of hosts selected in staging",
			},
			"active_in_production": {
				Type:        schema.TypeBool,
				Optional:    true,
				Description: "Whether to return names of hosts selected in production",
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
				Description: "JSON representation of hostnames",
			},
			"output_text": {
				Type:        schema.TypeString,
				Computed:    true,
				Description: "Text representation of hostnames",
			},
		},
	}
}

func dataSourceSelectableHostnamesRead(ctx context.Context, d *schema.ResourceData, m interface{}) diag.Diagnostics {
	meta := meta.Must(m)
	client := inst.Client(meta)
	logger := meta.Log("APPSEC", "dataSourceSelectableHostnamesRead")

	configID, err := tf.GetIntValue("config_id", d)
	if err != nil && !errors.Is(err, tf.ErrNotFound) {
		return diag.FromErr(err)
	}
	contractID, err := tf.GetStringValue("contractid", d)
	if err != nil && !errors.Is(err, tf.ErrNotFound) {
		return diag.FromErr(err)
	}
	group, err := tf.GetIntValue("groupid", d)
	if err != nil && !errors.Is(err, tf.ErrNotFound) {
		return diag.FromErr(err)
	}
	validParams := configID != 0 || (contractID != "" && group != 0)
	if !validParams {
		return diag.Errorf("either config_id or both contractid and groupdid must be supplied")
	}

	var version int
	if configID != 0 {
		if version, err = getLatestConfigVersion(ctx, configID, m); err != nil {
			return diag.FromErr(err)
		}
	}
	getSelectableHostnames := appsec.GetSelectableHostnamesRequest{
		ConfigID:   configID,
		Version:    version,
		ContractID: contractID,
		GroupID:    group,
	}

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
		return diag.Errorf("%s: %s", tf.ErrValueSet, err.Error())
	}

	var flagsetstg, flagsetprod string

	flagsetstg = "UNSET"
	flagsetprod = "UNSET"

	activeinstaging, ok := d.GetOkExists("active_in_staging") //nolint:staticcheck
	if ok {
		flagsetstg = "SET"
	}
	activeinproduction, ok := d.GetOkExists("active_in_production") //nolint:staticcheck
	if ok {
		flagsetprod = "SET"
	}

	newhdata := make([]string, 0, len(selectablehostnames.AvailableSet))
	for _, hosts := range selectablehostnames.AvailableSet {
		var flagstg, flagprod string

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
		return diag.Errorf("%s: %s", tf.ErrValueSet, err.Error())
	}

	ots := OutputTemplates{}
	InitTemplates(ots)

	outputtext, err := RenderTemplates(ots, "selectableHostsDS", selectablehostnames)
	if err != nil {
		return diag.FromErr(err)
	}
	if err := d.Set("output_text", outputtext); err != nil {
		return diag.Errorf("%s: %s", tf.ErrValueSet, err.Error())
	}

	d.SetId(strconv.Itoa(selectablehostnames.ConfigID))

	return nil
}
