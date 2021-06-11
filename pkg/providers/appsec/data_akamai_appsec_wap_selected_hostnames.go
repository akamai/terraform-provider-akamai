package appsec

import (
	"context"
	"encoding/json"
	"fmt"

	"github.com/akamai/AkamaiOPEN-edgegrid-golang/v2/pkg/appsec"
	"github.com/akamai/terraform-provider-akamai/v2/pkg/akamai"
	"github.com/akamai/terraform-provider-akamai/v2/pkg/tools"
	"github.com/hashicorp/terraform-plugin-sdk/v2/diag"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/schema"
)

type wapSelectedHostnamesOutputText struct {
	PolicyID string
	Hostname string
	Status   string
}

func dataSourceWAPSelectedHostnames() *schema.Resource {
	return &schema.Resource{
		ReadContext: dataSourceWAPSelectedHostnamesRead,
		Schema: map[string]*schema.Schema{
			"config_id": {
				Type:     schema.TypeInt,
				Required: true,
			},
			"security_policy_id": {
				Type:     schema.TypeString,
				Required: true,
			},
			"protected_hosts": {
				Type:        schema.TypeSet,
				Computed:    true,
				Elem:        &schema.Schema{Type: schema.TypeString},
				Description: "JSON List of protected hostnames",
			},
			"evaluated_hosts": {
				Type:        schema.TypeSet,
				Computed:    true,
				Elem:        &schema.Schema{Type: schema.TypeString},
				Description: "JSON List of evaluated hostnames",
			},
			"json": {
				Type:     schema.TypeString,
				Computed: true,
			},
			"output_text": {
				Type:        schema.TypeString,
				Computed:    true,
				Description: "Text Export representation",
			},
		},
	}
}

func dataSourceWAPSelectedHostnamesRead(ctx context.Context, d *schema.ResourceData, m interface{}) diag.Diagnostics {
	meta := akamai.Meta(m)
	client := inst.Client(meta)
	logger := meta.Log("APPSEC", "dataSourceWAPSelectedHostnamesRead")

	configid, err := tools.GetIntValue("config_id", d)
	if err != nil {
		return diag.FromErr(err)
	}
	securityPolicyID, err := tools.GetStringValue("security_policy_id", d)
	if err != nil {
		return diag.FromErr(err)
	}

	getWAPSelectedHostnamesRequest := appsec.GetWAPSelectedHostnamesRequest{}
	getWAPSelectedHostnamesRequest.ConfigID = configid
	getWAPSelectedHostnamesRequest.Version = getLatestConfigVersion(ctx, configid, m)
	getWAPSelectedHostnamesRequest.SecurityPolicyID = securityPolicyID

	WAPSelectedHostnames, err := client.GetWAPSelectedHostnames(ctx, getWAPSelectedHostnamesRequest)
	if err != nil {
		logger.Errorf("calling 'getWAPSelectedHostnames': %s", err.Error())
		return diag.FromErr(err)
	}

	if err := d.Set("protected_hosts", WAPSelectedHostnames.ProtectedHosts); err != nil {
		return diag.FromErr(fmt.Errorf("%w: %s", tools.ErrValueSet, err.Error()))
	}
	if err := d.Set("evaluated_hosts", WAPSelectedHostnames.EvaluatedHosts); err != nil {
		return diag.FromErr(fmt.Errorf("%w: %s", tools.ErrValueSet, err.Error()))
	}

	if err := d.Set("config_id", getWAPSelectedHostnamesRequest.ConfigID); err != nil {
		return diag.FromErr(fmt.Errorf("%w: %s", tools.ErrValueSet, err.Error()))
	}

	jsonBody, err := json.Marshal(WAPSelectedHostnames)
	if err != nil {
		return diag.FromErr(err)
	}
	if err := d.Set("json", string(jsonBody)); err != nil {
		return diag.FromErr(fmt.Errorf("%w: %s", tools.ErrValueSet, err.Error()))
	}

	ots := OutputTemplates{}
	InitTemplates(ots)

	textOutputCount := len(WAPSelectedHostnames.ProtectedHosts) + len(WAPSelectedHostnames.EvaluatedHosts)
	textOutputEntries := make([]wapSelectedHostnamesOutputText, 0, textOutputCount)
	for _, h := range WAPSelectedHostnames.ProtectedHosts {
		textOutputEntries = append(textOutputEntries, wapSelectedHostnamesOutputText{PolicyID: securityPolicyID, Hostname: h, Status: "protected"})
	}
	for _, h := range WAPSelectedHostnames.EvaluatedHosts {
		textOutputEntries = append(textOutputEntries, wapSelectedHostnamesOutputText{PolicyID: securityPolicyID, Hostname: h, Status: "evaluated"})
	}
	outputtext, err := RenderTemplates(ots, "WAPSelectedHostsDS", textOutputEntries)
	if err == nil {
		if err := d.Set("output_text", outputtext); err != nil {
			return diag.FromErr(fmt.Errorf("%w: %s", tools.ErrValueSet, err.Error()))
		}
	}

	d.SetId(fmt.Sprintf("%d:%s", configid, securityPolicyID))

	return nil
}
