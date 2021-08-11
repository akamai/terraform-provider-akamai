package appsec

import (
	"context"
	"errors"
	"fmt"
	"strconv"

	"github.com/akamai/AkamaiOPEN-edgegrid-golang/v2/pkg/appsec"
	"github.com/akamai/terraform-provider-akamai/v2/pkg/akamai"
	"github.com/akamai/terraform-provider-akamai/v2/pkg/tools"
	"github.com/hashicorp/terraform-plugin-sdk/v2/diag"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/customdiff"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/schema"
)

// appsec v1
//
// https://developer.akamai.com/api/cloud_security/application_security/v1.html
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
				Type:     schema.TypeInt,
				Required: true,
			},
			"security_policy_id": {
				Type:     schema.TypeString,
				Required: true,
			},
			"forward_to_http_header": {
				Type:     schema.TypeBool,
				Required: true,
			},
			"forward_shared_ip_to_http_header_siem": {
				Type:     schema.TypeBool,
				Required: true,
			},
		},
	}
}

func resourceReputationAnalysisCreate(ctx context.Context, d *schema.ResourceData, m interface{}) diag.Diagnostics {
	meta := akamai.Meta(m)
	client := inst.Client(meta)
	logger := meta.Log("APPSEC", "resourceReputationAnalysisCreate")
	logger.Debugf("in resourceReputationAnalysisCreate")

	configid, err := tools.GetIntValue("config_id", d)
	if err != nil && !errors.Is(err, tools.ErrNotFound) {
		return diag.FromErr(err)
	}
	version := getModifiableConfigVersion(ctx, configid, "reputationProfileAnalysis", m)
	policyid, err := tools.GetStringValue("security_policy_id", d)
	if err != nil && !errors.Is(err, tools.ErrNotFound) {
		return diag.FromErr(err)
	}
	forwardToHTTPHeader, err := tools.GetBoolValue("forward_to_http_header", d)
	if err != nil && !errors.Is(err, tools.ErrNotFound) {
		return diag.FromErr(err)
	}
	forwardSharedIPToHTTPHeaderSiem, err := tools.GetBoolValue("forward_shared_ip_to_http_header_siem", d)
	if err != nil && !errors.Is(err, tools.ErrNotFound) {
		return diag.FromErr(err)
	}

	createReputationAnalysis := appsec.UpdateReputationAnalysisRequest{
		ConfigID:                           configid,
		Version:                            version,
		PolicyID:                           policyid,
		ForwardToHTTPHeader:                forwardToHTTPHeader,
		ForwardSharedIPToHTTPHeaderAndSIEM: forwardSharedIPToHTTPHeaderSiem,
	}

	_, erru := client.UpdateReputationAnalysis(ctx, createReputationAnalysis)
	if erru != nil {
		logger.Errorf("calling 'createReputationAnalysis': %s", erru.Error())
		return diag.FromErr(erru)
	}

	d.SetId(fmt.Sprintf("%d:%s", createReputationAnalysis.ConfigID, createReputationAnalysis.PolicyID))

	return resourceReputationAnalysisRead(ctx, d, m)
}

func resourceReputationAnalysisRead(ctx context.Context, d *schema.ResourceData, m interface{}) diag.Diagnostics {
	meta := akamai.Meta(m)
	client := inst.Client(meta)
	logger := meta.Log("APPSEC", "resourceReputationAnalysisRead")
	logger.Debugf("in resourceReputationAnalysisRead")

	idParts, err := splitID(d.Id(), 2, "configid:securitypolicyid")
	if err != nil {
		return diag.FromErr(err)
	}
	configid, err := strconv.Atoi(idParts[0])
	if err != nil {
		return diag.FromErr(err)
	}
	version := getLatestConfigVersion(ctx, configid, m)
	policyid := idParts[1]

	getReputationAnalysis := appsec.GetReputationAnalysisRequest{
		ConfigID: configid,
		Version:  version,
		PolicyID: policyid,
	}

	resp, errg := client.GetReputationAnalysis(ctx, getReputationAnalysis)
	if errg != nil {
		logger.Errorf("calling 'getReputationAnalysis': %s", errg.Error())
		return diag.FromErr(errg)
	}

	if err := d.Set("config_id", getReputationAnalysis.ConfigID); err != nil {
		return diag.FromErr(fmt.Errorf("%w: %s", tools.ErrValueSet, err.Error()))
	}
	if err := d.Set("security_policy_id", getReputationAnalysis.PolicyID); err != nil {
		return diag.FromErr(fmt.Errorf("%w: %s", tools.ErrValueSet, err.Error()))
	}
	if err := d.Set("forward_to_http_header", resp.ForwardToHTTPHeader); err != nil {
		return diag.FromErr(fmt.Errorf("%w: %s", tools.ErrValueSet, err.Error()))
	}
	if err := d.Set("forward_shared_ip_to_http_header_siem", resp.ForwardSharedIPToHTTPHeaderAndSIEM); err != nil {
		return diag.FromErr(fmt.Errorf("%w: %s", tools.ErrValueSet, err.Error()))
	}

	return nil
}

func resourceReputationAnalysisUpdate(ctx context.Context, d *schema.ResourceData, m interface{}) diag.Diagnostics {
	meta := akamai.Meta(m)
	client := inst.Client(meta)
	logger := meta.Log("APPSEC", "resourceReputationAnalysisUpdate")
	logger.Debugf("in resourceReputationAnalysisUpdate")

	idParts, err := splitID(d.Id(), 2, "configid:securitypolicyid")
	if err != nil {
		return diag.FromErr(err)
	}
	configid, err := strconv.Atoi(idParts[0])
	if err != nil {
		return diag.FromErr(err)
	}
	version := getModifiableConfigVersion(ctx, configid, "reputationProfileAnalysis", m)
	policyid := idParts[1]
	forwardToHTTPHeader, err := tools.GetBoolValue("forward_to_http_header", d)
	if err != nil && !errors.Is(err, tools.ErrNotFound) {
		return diag.FromErr(err)
	}
	forwardSharedIPToHTTPHeaderSiem, err := tools.GetBoolValue("forward_shared_ip_to_http_header_siem", d)
	if err != nil && !errors.Is(err, tools.ErrNotFound) {
		return diag.FromErr(err)
	}

	updateReputationAnalysis := appsec.UpdateReputationAnalysisRequest{
		ConfigID:                           configid,
		Version:                            version,
		PolicyID:                           policyid,
		ForwardToHTTPHeader:                forwardToHTTPHeader,
		ForwardSharedIPToHTTPHeaderAndSIEM: forwardSharedIPToHTTPHeaderSiem,
	}

	_, erru := client.UpdateReputationAnalysis(ctx, updateReputationAnalysis)
	if erru != nil {
		logger.Errorf("calling 'updateReputationAnalysis': %s", erru.Error())
		return diag.FromErr(erru)
	}

	return resourceReputationAnalysisRead(ctx, d, m)
}

func resourceReputationAnalysisDelete(ctx context.Context, d *schema.ResourceData, m interface{}) diag.Diagnostics {
	meta := akamai.Meta(m)
	client := inst.Client(meta)
	logger := meta.Log("APPSEC", "resourceReputationAnalysisDelete")
	logger.Debugf("in resourceReputationAnalysisDelete")

	idParts, err := splitID(d.Id(), 2, "configid:securitypolicyid")
	if err != nil {
		return diag.FromErr(err)
	}
	configid, err := strconv.Atoi(idParts[0])
	if err != nil {
		return diag.FromErr(err)
	}
	version := getModifiableConfigVersion(ctx, configid, "reputationProfileAnalysis", m)
	policyid := idParts[1]

	RemoveReputationAnalysis := appsec.RemoveReputationAnalysisRequest{
		ConfigID:                           configid,
		Version:                            version,
		PolicyID:                           policyid,
		ForwardToHTTPHeader:                false,
		ForwardSharedIPToHTTPHeaderAndSIEM: false,
	}

	_, erru := client.RemoveReputationAnalysis(ctx, RemoveReputationAnalysis)
	if erru != nil {
		logger.Errorf("calling 'RemoveReputationAnalysis': %s", erru.Error())
		return diag.FromErr(erru)
	}

	d.SetId("")

	return nil
}
