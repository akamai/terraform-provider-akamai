package appsec

import (
	"context"
	"encoding/json"
	"fmt"

	"github.com/akamai/AkamaiOPEN-edgegrid-golang/v10/pkg/appsec"
	"github.com/akamai/terraform-provider-akamai/v7/pkg/common/tf"
	"github.com/akamai/terraform-provider-akamai/v7/pkg/meta"
	"github.com/hashicorp/terraform-plugin-sdk/v2/diag"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/schema"
)

type aapSelectedHostnamesOutputText struct {
	PolicyID string
	Hostname string
	Status   string
}

func dataSourceAAPSelectedHostnames() *schema.Resource {
	return &schema.Resource{
		ReadContext: dataSourceAAPSelectedHostnamesRead,
		Schema: map[string]*schema.Schema{
			"config_id": {
				Type:        schema.TypeInt,
				Required:    true,
				Description: "Unique identifier of the security configuration",
			},
			"security_policy_id": {
				Type:        schema.TypeString,
				Required:    true,
				Description: "Unique identifier of the security policy",
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
				Type:        schema.TypeString,
				Computed:    true,
				Description: "JSON representation",
			},
			"output_text": {
				Type:        schema.TypeString,
				Computed:    true,
				Description: "Text representation",
			},
		},
	}
}

func dataSourceAAPSelectedHostnamesRead(ctx context.Context, d *schema.ResourceData, m interface{}) diag.Diagnostics {
	meta := meta.Must(m)
	client := inst.Client(meta)
	logger := meta.Log("APPSEC", "dataSourceAAPSelectedHostnamesRead")

	configID, err := tf.GetIntValue("config_id", d)
	if err != nil {
		return diag.FromErr(err)
	}
	version, err := getLatestConfigVersion(ctx, configID, m)
	if err != nil {
		return diag.FromErr(err)
	}
	securityPolicyID, err := tf.GetStringValue("security_policy_id", d)
	if err != nil {
		return diag.FromErr(err)
	}

	getConfigurationRequest := appsec.GetConfigurationRequest{ConfigID: configID}
	configuration, err := client.GetConfiguration(ctx, getConfigurationRequest)
	if err != nil {
		logger.Errorf("calling 'getConfiguration': %s", err.Error())
		return diag.FromErr(err)
	}
	target := configuration.TargetProduct

	if target == "KSD" {
		getSelectedHostnames := appsec.GetSelectedHostnamesRequest{
			ConfigID: configID,
			Version:  version,
		}

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
			return diag.Errorf("%s: %s", tf.ErrValueSet, err.Error())
		}

		getMatchTargets := appsec.GetMatchTargetsRequest{
			ConfigID:      configID,
			ConfigVersion: version,
		}

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
			return diag.Errorf("%s: %s", tf.ErrValueSet, err.Error())
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
		if err != nil {
			return diag.FromErr(err)
		}
		if err := d.Set("output_text", websiteMatchTargetsText); err != nil {
			return diag.Errorf("%s: %s", tf.ErrValueSet, err.Error())
		}
	} else { // WAP_AAG and WAP_PLUS accounts
		getWAPSelectedHostnamesRequest := appsec.GetWAPSelectedHostnamesRequest{
			ConfigID:         configID,
			Version:          version,
			SecurityPolicyID: securityPolicyID,
		}

		WAPSelectedHostnames, err := client.GetWAPSelectedHostnames(ctx, getWAPSelectedHostnamesRequest)
		if err != nil {
			logger.Errorf("calling 'getWAPSelectedHostnames': %s", err.Error())
			return diag.FromErr(err)
		}

		if err := d.Set("protected_hosts", WAPSelectedHostnames.ProtectedHosts); err != nil {
			return diag.Errorf("%s: %s", tf.ErrValueSet, err.Error())
		}
		if err := d.Set("evaluated_hosts", WAPSelectedHostnames.EvaluatedHosts); err != nil {
			return diag.Errorf("%s: %s", tf.ErrValueSet, err.Error())
		}

		jsonBody, err := json.Marshal(WAPSelectedHostnames)
		if err != nil {
			return diag.FromErr(err)
		}
		if err := d.Set("json", string(jsonBody)); err != nil {
			return diag.Errorf("%s: %s", tf.ErrValueSet, err.Error())
		}

		ots := OutputTemplates{}
		InitTemplates(ots)

		textOutputCount := len(WAPSelectedHostnames.ProtectedHosts) + len(WAPSelectedHostnames.EvaluatedHosts)
		textOutputEntries := make([]aapSelectedHostnamesOutputText, 0, textOutputCount)
		for _, h := range WAPSelectedHostnames.ProtectedHosts {
			entry := aapSelectedHostnamesOutputText{PolicyID: securityPolicyID, Hostname: h, Status: "protected"}
			textOutputEntries = append(textOutputEntries, entry)
		}
		for _, h := range WAPSelectedHostnames.EvaluatedHosts {
			entry := aapSelectedHostnamesOutputText{PolicyID: securityPolicyID, Hostname: h, Status: "evaluated"}
			textOutputEntries = append(textOutputEntries, entry)
		}
		outputtext, err := RenderTemplates(ots, "AAPSelectedHostsDS", textOutputEntries)
		if err != nil {
			return diag.Errorf("%s: %s", tf.ErrValueSet, err.Error())
		}
		if err := d.Set("output_text", outputtext); err != nil {
			return diag.Errorf("%s: %s", tf.ErrValueSet, err.Error())
		}
	}

	d.SetId(fmt.Sprintf("%d:%s", configID, securityPolicyID))

	return nil
}
