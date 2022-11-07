package appsec

import (
	"context"
	"errors"
	"fmt"
	"strconv"

	"github.com/akamai/AkamaiOPEN-edgegrid-golang/v3/pkg/appsec"
	"github.com/akamai/terraform-provider-akamai/v2/pkg/akamai"
	"github.com/akamai/terraform-provider-akamai/v2/pkg/tools"
	"github.com/hashicorp/terraform-plugin-sdk/v2/diag"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/customdiff"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/schema"
)

// appsec v1
//
// https://techdocs.akamai.com/application-security/reference/api
func resourceReputationAnalysis() *schema.Resource {
	return &schema.Resource{
		CreateContext: resourceReputationAnalysisCreate,
		ReadContext:   resourceReputationAnalysisRead,
		UpdateContext: resourceReputationAnalysisUpdate,
		DeleteContext: resourceReputationAnalysisDelete,
		CustomizeDiff: customdiff.All(
			VerifyIDUnchanged,
		),
		Importer: &schema.ResourceImporter{
			StateContext: schema.ImportStatePassthroughContext,
		},
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
			"forward_to_http_header": {
				Type:        schema.TypeBool,
				Required:    true,
				Description: "Whether to add client reputation details to requests forwarded to the origin server",
			},
			"forward_shared_ip_to_http_header_siem": {
				Type:        schema.TypeBool,
				Required:    true,
				Description: "Whether to add a value indicating that shared IPs are included in HTTP header and SIEM integration",
			},
		},
	}
}

func resourceReputationAnalysisCreate(ctx context.Context, d *schema.ResourceData, m interface{}) diag.Diagnostics {
	meta := akamai.Meta(m)
	client := inst.Client(meta)
	logger := meta.Log("APPSEC", "resourceReputationAnalysisCreate")
	logger.Debugf("in resourceReputationAnalysisCreate")

	configID, err := tools.GetIntValue("config_id", d)
	if err != nil {
		return diag.FromErr(err)
	}
	version, err := getModifiableConfigVersion(ctx, configID, "reputationProfileAnalysis", m)
	if err != nil {
		return diag.FromErr(err)
	}
	policyID, err := tools.GetStringValue("security_policy_id", d)
	if err != nil {
		return diag.FromErr(err)
	}
	forwardToHTTPHeader, err := tools.GetBoolValue("forward_to_http_header", d)
	if err != nil {
		return diag.FromErr(err)
	}
	forwardSharedIPToHTTPHeaderSiem, err := tools.GetBoolValue("forward_shared_ip_to_http_header_siem", d)
	if err != nil {
		return diag.FromErr(err)
	}

	createReputationAnalysis := appsec.UpdateReputationAnalysisRequest{
		ConfigID:                           configID,
		Version:                            version,
		PolicyID:                           policyID,
		ForwardToHTTPHeader:                forwardToHTTPHeader,
		ForwardSharedIPToHTTPHeaderAndSIEM: forwardSharedIPToHTTPHeaderSiem,
	}

	_, err = client.UpdateReputationAnalysis(ctx, createReputationAnalysis)
	if err != nil {
		logger.Errorf("calling 'createReputationAnalysis': %s", err.Error())
		return diag.FromErr(err)
	}

	d.SetId(fmt.Sprintf("%d:%s", createReputationAnalysis.ConfigID, createReputationAnalysis.PolicyID))

	return resourceReputationAnalysisRead(ctx, d, m)
}

func resourceReputationAnalysisRead(ctx context.Context, d *schema.ResourceData, m interface{}) diag.Diagnostics {
	meta := akamai.Meta(m)
	client := inst.Client(meta)
	logger := meta.Log("APPSEC", "resourceReputationAnalysisRead")
	logger.Debugf("in resourceReputationAnalysisRead")

	iDParts, err := splitID(d.Id(), 2, "configID:securityPolicyID")
	if err != nil {
		return diag.FromErr(err)
	}
	configID, err := strconv.Atoi(iDParts[0])
	if err != nil {
		return diag.FromErr(err)
	}
	version, err := getLatestConfigVersion(ctx, configID, m)
	if err != nil {
		return diag.FromErr(err)
	}
	policyID := iDParts[1]

	getReputationAnalysis := appsec.GetReputationAnalysisRequest{
		ConfigID: configID,
		Version:  version,
		PolicyID: policyID,
	}

	resp, err := client.GetReputationAnalysis(ctx, getReputationAnalysis)
	if err != nil {
		logger.Errorf("calling 'getReputationAnalysis': %s", err.Error())
		return diag.FromErr(err)
	}

	if err := d.Set("config_id", getReputationAnalysis.ConfigID); err != nil {
		return diag.Errorf("%s: %s", tools.ErrValueSet, err.Error())
	}
	if err := d.Set("security_policy_id", getReputationAnalysis.PolicyID); err != nil {
		return diag.Errorf("%s: %s", tools.ErrValueSet, err.Error())
	}
	if err := d.Set("forward_to_http_header", resp.ForwardToHTTPHeader); err != nil {
		return diag.Errorf("%s: %s", tools.ErrValueSet, err.Error())
	}
	if err := d.Set("forward_shared_ip_to_http_header_siem", resp.ForwardSharedIPToHTTPHeaderAndSIEM); err != nil {
		return diag.Errorf("%s: %s", tools.ErrValueSet, err.Error())
	}

	return nil
}

func resourceReputationAnalysisUpdate(ctx context.Context, d *schema.ResourceData, m interface{}) diag.Diagnostics {
	meta := akamai.Meta(m)
	client := inst.Client(meta)
	logger := meta.Log("APPSEC", "resourceReputationAnalysisUpdate")
	logger.Debugf("in resourceReputationAnalysisUpdate")

	iDParts, err := splitID(d.Id(), 2, "configID:securityPolicyID")
	if err != nil {
		return diag.FromErr(err)
	}
	configID, err := strconv.Atoi(iDParts[0])
	if err != nil {
		return diag.FromErr(err)
	}
	version, err := getModifiableConfigVersion(ctx, configID, "reputationProfileAnalysis", m)
	if err != nil {
		return diag.FromErr(err)
	}
	policyID := iDParts[1]
	forwardToHTTPHeader, err := tools.GetBoolValue("forward_to_http_header", d)
	if err != nil && !errors.Is(err, tools.ErrNotFound) {
		return diag.FromErr(err)
	}
	forwardSharedIPToHTTPHeaderSiem, err := tools.GetBoolValue("forward_shared_ip_to_http_header_siem", d)
	if err != nil && !errors.Is(err, tools.ErrNotFound) {
		return diag.FromErr(err)
	}

	updateReputationAnalysis := appsec.UpdateReputationAnalysisRequest{
		ConfigID:                           configID,
		Version:                            version,
		PolicyID:                           policyID,
		ForwardToHTTPHeader:                forwardToHTTPHeader,
		ForwardSharedIPToHTTPHeaderAndSIEM: forwardSharedIPToHTTPHeaderSiem,
	}

	_, err = client.UpdateReputationAnalysis(ctx, updateReputationAnalysis)
	if err != nil {
		logger.Errorf("calling 'updateReputationAnalysis': %s", err.Error())
		return diag.FromErr(err)
	}

	return resourceReputationAnalysisRead(ctx, d, m)
}

func resourceReputationAnalysisDelete(ctx context.Context, d *schema.ResourceData, m interface{}) diag.Diagnostics {
	meta := akamai.Meta(m)
	client := inst.Client(meta)
	logger := meta.Log("APPSEC", "resourceReputationAnalysisDelete")
	logger.Debugf("in resourceReputationAnalysisDelete")

	iDParts, err := splitID(d.Id(), 2, "configID:securityPolicyID")
	if err != nil {
		return diag.FromErr(err)
	}
	configID, err := strconv.Atoi(iDParts[0])
	if err != nil {
		return diag.FromErr(err)
	}
	version, err := getModifiableConfigVersion(ctx, configID, "reputationProfileAnalysis", m)
	if err != nil {
		return diag.FromErr(err)
	}
	policyID := iDParts[1]

	RemoveReputationAnalysis := appsec.RemoveReputationAnalysisRequest{
		ConfigID:                           configID,
		Version:                            version,
		PolicyID:                           policyID,
		ForwardToHTTPHeader:                false,
		ForwardSharedIPToHTTPHeaderAndSIEM: false,
	}

	_, err = client.RemoveReputationAnalysis(ctx, RemoveReputationAnalysis)
	if err != nil {
		logger.Errorf("calling 'RemoveReputationAnalysis': %s", err.Error())
		return diag.FromErr(err)
	}
	return nil
}
