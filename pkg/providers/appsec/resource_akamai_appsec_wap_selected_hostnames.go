package appsec

import (
	"context"
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
func resourceWAPSelectedHostnames() *schema.Resource {
	return &schema.Resource{
		CreateContext: resourceWAPSelectedHostnamesCreate,
		ReadContext:   resourceWAPSelectedHostnamesRead,
		UpdateContext: resourceWAPSelectedHostnamesUpdate,
		DeleteContext: resourceWAPSelectedHostnamesDelete,
		Importer: &schema.ResourceImporter{
			State: schema.ImportStatePassthrough,
		},
		CustomizeDiff: customdiff.All(
			VerifyIdUnchanged,
		),
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
				Type:     schema.TypeSet,
				Required: true,
				Elem:     &schema.Schema{Type: schema.TypeString},
			},
			"evaluated_hosts": {
				Type:     schema.TypeSet,
				Required: true,
				Elem:     &schema.Schema{Type: schema.TypeString},
			},
		},
	}
}

func resourceWAPSelectedHostnamesCreate(ctx context.Context, d *schema.ResourceData, m interface{}) diag.Diagnostics {
	meta := akamai.Meta(m)
	client := inst.Client(meta)
	logger := meta.Log("APPSEC", "resourceWAPSelectedHostnamesCreate")
	logger.Debugf("!!! resourceWAPSelectedHostnamesCreate")

	configID, err := tools.GetIntValue("config_id", d)
	if err != nil {
		return diag.FromErr(err)
	}
	securityPolicyID, err := tools.GetStringValue("security_policy_id", d)
	if err != nil {
		return diag.FromErr(err)
	}
	protectedHosts, err := tools.GetSetValue("protected_hosts", d)
	if err != nil {
		return diag.FromErr(err)
	}
	evaluatedHosts, err := tools.GetSetValue("evaluated_hosts", d)
	if err != nil {
		return diag.FromErr(err)
	}

	// verify that none of the hostnames in either list is also in the other list
	for _, h := range protectedHosts.List() {
		if evaluatedHosts.Contains(h) {
			return diag.FromErr(fmt.Errorf("host %s cannot be in both protected and evaluated hosts", h))
		}
	}
	for _, h := range evaluatedHosts.List() {
		if protectedHosts.Contains(h) {
			return diag.FromErr(fmt.Errorf("host %s cannot be in both protected and evaluated hosts", h))
		}
	}

	// convert to lists of strings
	protectedHostnames := tools.SetToStringSlice(protectedHosts)
	evalHostnames := tools.SetToStringSlice(evaluatedHosts)

	updateWAPSelectedHostnames := appsec.UpdateWAPSelectedHostnamesRequest{}
	updateWAPSelectedHostnames.ConfigID = configID
	updateWAPSelectedHostnames.Version = getModifiableConfigVersion(ctx, configID, "wapSelectedHostnames", m)
	updateWAPSelectedHostnames.SecurityPolicyID = securityPolicyID
	updateWAPSelectedHostnames.ProtectedHosts = protectedHostnames
	updateWAPSelectedHostnames.EvaluatedHosts = evalHostnames

	_, err = client.UpdateWAPSelectedHostnames(ctx, updateWAPSelectedHostnames)
	if err != nil {
		logger.Errorf("calling 'UpdateWAPSelectedHostnames': %s", err.Error())
		return diag.FromErr(fmt.Errorf("%w: %s", tools.ErrValueSet, err.Error()))
	}

	if err := d.Set("config_id", configID); err != nil {
		return diag.FromErr(fmt.Errorf("%w: %s", tools.ErrValueSet, err.Error()))
	}
	if err := d.Set("security_policy_id", securityPolicyID); err != nil {
		return diag.FromErr(fmt.Errorf("%w: %s", tools.ErrValueSet, err.Error()))
	}
	if err := d.Set("protected_hosts", protectedHosts); err != nil {
		return diag.FromErr(fmt.Errorf("%w: %s", tools.ErrValueSet, err.Error()))
	}
	if err := d.Set("evaluated_hosts", evaluatedHosts); err != nil {
		return diag.FromErr(fmt.Errorf("%w: %s", tools.ErrValueSet, err.Error()))
	}

	d.SetId(fmt.Sprintf("%d:%s", configID, securityPolicyID))

	return resourceWAPSelectedHostnamesRead(ctx, d, m)
}

func resourceWAPSelectedHostnamesRead(ctx context.Context, d *schema.ResourceData, m interface{}) diag.Diagnostics {
	meta := akamai.Meta(m)
	client := inst.Client(meta)
	logger := meta.Log("APPSEC", "resourceWAPSelectedHostnamesRead")
	logger.Debugf("!!! resourceWAPSelectedHostnamesRead")

	idParts, err := splitID(d.Id(), 2, "configid:policyid")
	if err != nil {
		return diag.FromErr(err)
	}
	configID, err := strconv.Atoi(idParts[0])
	if err != nil {
		return diag.FromErr(err)
	}
	version := getLatestConfigVersion(ctx, configID, m)
	securityPolicyID := idParts[1]

	getWAPSelectedHostnamesRequest := appsec.GetWAPSelectedHostnamesRequest{}
	getWAPSelectedHostnamesRequest.ConfigID = configID
	getWAPSelectedHostnamesRequest.Version = version
	getWAPSelectedHostnamesRequest.SecurityPolicyID = securityPolicyID

	wapSelectedHostnames, err := client.GetWAPSelectedHostnames(ctx, getWAPSelectedHostnamesRequest)
	if err != nil {
		logger.Errorf("calling 'getWAPSelectedHostnames': %s", err.Error())
		return diag.FromErr(err)
	}

	if err := d.Set("config_id", configID); err != nil {
		return diag.FromErr(fmt.Errorf("%w: %s", tools.ErrValueSet, err.Error()))
	}
	if err := d.Set("security_policy_id", securityPolicyID); err != nil {
		return diag.FromErr(fmt.Errorf("%w: %s", tools.ErrValueSet, err.Error()))
	}
	wapProtectedHostSet := schema.Set{F: schema.HashString}
	for _, h := range wapSelectedHostnames.ProtectedHosts {
		wapProtectedHostSet.Add(h)
	}
	if err := d.Set("protected_hosts", wapProtectedHostSet.List()); err != nil {
		return diag.FromErr(fmt.Errorf("%w: %s", tools.ErrValueSet, err.Error()))
	}
	wapEvaluatedHostSet := schema.Set{F: schema.HashString}
	for _, h := range wapSelectedHostnames.EvaluatedHosts {
		wapEvaluatedHostSet.Add(h)
	}
	if err := d.Set("evaluated_hosts", wapEvaluatedHostSet.List()); err != nil {
		return diag.FromErr(fmt.Errorf("%w: %s", tools.ErrValueSet, err.Error()))
	}

	return nil
}

func resourceWAPSelectedHostnamesUpdate(ctx context.Context, d *schema.ResourceData, m interface{}) diag.Diagnostics {
	meta := akamai.Meta(m)
	client := inst.Client(meta)
	logger := meta.Log("APPSEC", "resourceWAPSelectedHostnamesUpdate")
	logger.Debugf("!!! resourceWAPSelectedHostnamesUpdate")

	configID, err := tools.GetIntValue("config_id", d)
	if err != nil {
		return diag.FromErr(err)
	}
	securityPolicyID, err := tools.GetStringValue("security_policy_id", d)
	if err != nil {
		return diag.FromErr(err)
	}
	protectedHosts, err := tools.GetSetValue("protected_hosts", d)
	if err != nil {
		return diag.FromErr(err)
	}
	evaluatedHosts, err := tools.GetSetValue("evaluated_hosts", d)
	if err != nil {
		return diag.FromErr(err)
	}

	// verify that none of the hostnames in either list is also in the other list
	for _, h := range protectedHosts.List() {
		if evaluatedHosts.Contains(h) {
			return diag.FromErr(fmt.Errorf("host %s cannot be in both protected and evaluated hosts", h))
		}
	}
	for _, h := range evaluatedHosts.List() {
		if protectedHosts.Contains(h) {
			return diag.FromErr(fmt.Errorf("host %s cannot be in both protected and evaluated hosts", h))
		}
	}

	// convert to lists of Hostname structs
	protectedHostnames := tools.SetToStringSlice(protectedHosts)
	evalHostnames := tools.SetToStringSlice(evaluatedHosts)

	updateWAPSelectedHostnames := appsec.UpdateWAPSelectedHostnamesRequest{}
	updateWAPSelectedHostnames.ConfigID = configID
	updateWAPSelectedHostnames.Version = getModifiableConfigVersion(ctx, configID, "wapSelectedHostnames", m)
	updateWAPSelectedHostnames.SecurityPolicyID = securityPolicyID
	updateWAPSelectedHostnames.ProtectedHosts = protectedHostnames
	updateWAPSelectedHostnames.EvaluatedHosts = evalHostnames

	_, err = client.UpdateWAPSelectedHostnames(ctx, updateWAPSelectedHostnames)
	if err != nil {
		logger.Errorf("calling 'UpdateWAPSelectedHostnames': %s", err.Error())
		return diag.FromErr(fmt.Errorf("%w: %s", tools.ErrValueSet, err.Error()))
	}

	if d.HasChange("protected_hosts") {
		if err := d.Set("protected_hosts", protectedHosts); err != nil {
			return diag.FromErr(fmt.Errorf("%w: %s", tools.ErrValueSet, err.Error()))
		}
	}
	if d.HasChange("evaluated_hosts") {
		if err := d.Set("evaluated_hosts", evaluatedHosts); err != nil {
			return diag.FromErr(fmt.Errorf("%w: %s", tools.ErrValueSet, err.Error()))
		}
	}

	return resourceWAPSelectedHostnamesRead(ctx, d, m)
}

func resourceWAPSelectedHostnamesDelete(ctx context.Context, d *schema.ResourceData, m interface{}) diag.Diagnostics {
	return schema.NoopContext(nil, d, m)
}
