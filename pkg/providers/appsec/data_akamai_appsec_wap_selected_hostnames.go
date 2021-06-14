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
			"selected_hosts": {
				Type:        schema.TypeSet,
				Computed:    true,
				Elem:        &schema.Schema{Type: schema.TypeString},
				Description: "JSON List of selected hostnames (for non-WAP accounts)",
			},
			"match_targets": {
				Type:        schema.TypeString,
				Computed:    true,
				Description: "Match target information (for non-WAP accounts)",
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
	version := getLatestConfigVersion(ctx, configid, m)
	securityPolicyID, err := tools.GetStringValue("security_policy_id", d)
	if err != nil {
		return diag.FromErr(err)
	}

	getConfigurationRequest := appsec.GetConfigurationRequest{ConfigID: configid}
	configuration, err := client.GetConfiguration(ctx, getConfigurationRequest)
	if err != nil {
		logger.Errorf("calling 'getConfiguration': %s", err.Error())
		return diag.FromErr(err)
	}
	target := configuration.TargetProduct

	if target == "WAP" {
		getWAPSelectedHostnamesRequest := appsec.GetWAPSelectedHostnamesRequest{}
		getWAPSelectedHostnamesRequest.ConfigID = configid
		getWAPSelectedHostnamesRequest.Version = version
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
	} else {
		getSelectedHostnames := appsec.GetSelectedHostnamesRequest{}
		getSelectedHostnames.ConfigID = configid
		getSelectedHostnames.Version = version

		selectedhostnames, err := client.GetSelectedHostnames(ctx, getSelectedHostnames)
		if err != nil {
			logger.Errorf("calling 'getSelectedHostnames': %s", err.Error())
			return diag.FromErr(err)
		}
		newhdata := make([]string, 0, len(selectedhostnames.HostnameList))
		for _, hosts := range selectedhostnames.HostnameList {
			newhdata = append(newhdata, hosts.Hostname)
		}
		if err := d.Set("selected_hosts", newhdata); err != nil {
			return diag.FromErr(fmt.Errorf("%w: %s", tools.ErrValueSet, err.Error()))
		}

		getMatchTargets := appsec.GetMatchTargetsRequest{}
		getMatchTargets.ConfigID = configid
		getMatchTargets.ConfigVersion = version

		matchtargets, err := client.GetMatchTargets(ctx, getMatchTargets)
		if err != nil {
			logger.Errorf("calling 'getMatchTargets': %s", err.Error())
			return diag.FromErr(err)
		}

		jsonBody, err := json.Marshal(matchtargets)
		if err != nil {
			logger.Errorf("calling 'getMatchTargets': %s", err.Error())
			return diag.FromErr(err)
		}
		if err := d.Set("match_targets", string(jsonBody)); err != nil {
			return diag.FromErr(fmt.Errorf("%w: %s", tools.ErrValueSet, err.Error()))
		}

		ots := OutputTemplates{}
		InitTemplates(ots)

		matchtargetCount := len(matchtargets.MatchTargets.WebsiteTargets) + len(matchtargets.MatchTargets.APITargets)
		matchtargetsOutputText := make([]MatchTargetOutputText, 0, matchtargetCount)
		for _, value := range matchtargets.MatchTargets.WebsiteTargets {
			matchtargetsOutputText = append(matchtargetsOutputText, MatchTargetOutputText{value.TargetID, value.SecurityPolicy.PolicyID, WebsiteTarget})
		}
		for _, value := range matchtargets.MatchTargets.APITargets {
			matchtargetsOutputText = append(matchtargetsOutputText, MatchTargetOutputText{value.TargetID, value.SecurityPolicy.PolicyID, APITarget})
		}
		websiteMatchTargetsText, err := RenderTemplates(ots, "matchTargetDS", matchtargetsOutputText)
		if err == nil {
			if err := d.Set("output_text", websiteMatchTargetsText); err != nil {
				return diag.FromErr(fmt.Errorf("%w: %s", tools.ErrValueSet, err.Error()))
			}
		}
	}

	d.SetId(fmt.Sprintf("%d:%s", configid, securityPolicyID))

	return nil
}
