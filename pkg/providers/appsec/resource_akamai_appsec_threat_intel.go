package appsec

import (
	"context"
	"fmt"
	"strconv"

	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/validation"

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
func resourceThreatIntel() *schema.Resource {
	return &schema.Resource{
		CreateContext: resourceThreatIntelCreate,
		ReadContext:   resourceThreatIntelRead,
		UpdateContext: resourceThreatIntelUpdate,
		DeleteContext: resourceThreatIntelDelete,
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
			"threat_intel": {
				Type:     schema.TypeString,
				Required: true,
				ValidateFunc: validation.StringInSlice([]string{
					"on",
					"off",
				}, false),
			},
		},
	}
}

func resourceThreatIntelCreate(ctx context.Context, d *schema.ResourceData, m interface{}) diag.Diagnostics {
	meta := akamai.Meta(m)
	client := inst.Client(meta)
	logger := meta.Log("APPSEC", "resourceThreatIntelCreate")
	logger.Debugf("in resourceThreatIntelCreate")

	configID, err := tools.GetIntValue("config_id", d)
	if err != nil {
		return diag.FromErr(err)
	}
	version := getModifiableConfigVersion(ctx, configID, "threatIntel", m)
	policyID, err := tools.GetStringValue("security_policy_id", d)
	if err != nil {
		return diag.FromErr(err)
	}
	threatintel, err := tools.GetStringValue("threat_intel", d)
	if err != nil {
		return diag.FromErr(err)
	}

	createThreatIntel := appsec.UpdateThreatIntelRequest{
		ConfigID:    configID,
		Version:     version,
		PolicyID:    policyID,
		ThreatIntel: threatintel,
	}

	_, err = client.UpdateThreatIntel(ctx, createThreatIntel)
	if err != nil {
		logger.Errorf("calling 'createThreatIntel': %s", err.Error())
		return diag.FromErr(err)
	}

	d.SetId(fmt.Sprintf("%d:%s", createThreatIntel.ConfigID, createThreatIntel.PolicyID))

	return resourceThreatIntelRead(ctx, d, m)
}

func resourceThreatIntelRead(ctx context.Context, d *schema.ResourceData, m interface{}) diag.Diagnostics {
	meta := akamai.Meta(m)
	client := inst.Client(meta)
	logger := meta.Log("APPSEC", "resourceThreatIntelRead")
	logger.Debugf(" in resourceThreatIntelRead")

	iDParts, err := splitID(d.Id(), 2, "configID:securityPolicyID")
	if err != nil {
		return diag.FromErr(err)
	}
	configID, err := strconv.Atoi(iDParts[0])
	if err != nil {
		return diag.FromErr(err)
	}
	version := getLatestConfigVersion(ctx, configID, m)
	policyID := iDParts[1]

	getThreatIntel := appsec.GetThreatIntelRequest{
		ConfigID: configID,
		Version:  version,
		PolicyID: policyID,
	}

	threatintel, err := client.GetThreatIntel(ctx, getThreatIntel)
	if err != nil {
		logger.Warnf("calling 'getThreatIntel': %s", err.Error())
	}

	if err := d.Set("config_id", getThreatIntel.ConfigID); err != nil {
		return diag.Errorf("%s: %s", tools.ErrValueSet, err.Error())
	}
	if err := d.Set("security_policy_id", getThreatIntel.PolicyID); err != nil {
		return diag.Errorf("%s: %s", tools.ErrValueSet, err.Error())
	}

	if err := d.Set("threat_intel", threatintel.ThreatIntel); err != nil {
		return diag.Errorf("%s: %s", tools.ErrValueSet, err.Error())
	}

	return nil
}

func resourceThreatIntelUpdate(ctx context.Context, d *schema.ResourceData, m interface{}) diag.Diagnostics {
	meta := akamai.Meta(m)
	client := inst.Client(meta)
	logger := meta.Log("APPSEC", "resourceThreatIntelUpdate")
	logger.Debugf("in resourceThreatIntelUpdate")

	iDParts, err := splitID(d.Id(), 2, "configID:securityPolicyID")
	if err != nil {
		return diag.FromErr(err)
	}
	configID, err := strconv.Atoi(iDParts[0])
	if err != nil {
		return diag.FromErr(err)
	}
	version := getModifiableConfigVersion(ctx, configID, "threatIntel", m)
	policyID := iDParts[1]

	threatintel, err := tools.GetStringValue("threat_intel", d)
	if err != nil {
		return diag.FromErr(err)
	}

	updateThreatIntel := appsec.UpdateThreatIntelRequest{
		ConfigID:    configID,
		Version:     version,
		PolicyID:    policyID,
		ThreatIntel: threatintel,
	}

	_, err = client.UpdateThreatIntel(ctx, updateThreatIntel)
	if err != nil {
		logger.Errorf("calling 'updateThreatIntel': %s", err.Error())
		return diag.FromErr(err)
	}

	return resourceThreatIntelRead(ctx, d, m)
}

func resourceThreatIntelDelete(_ context.Context, d *schema.ResourceData, m interface{}) diag.Diagnostics {
	meta := akamai.Meta(m)
	logger := meta.Log("APPSEC", "resourceThreatIntelDelete")
	logger.Debugf("in resourceThreatIntelDelete")

	d.SetId("")

	return nil
}
