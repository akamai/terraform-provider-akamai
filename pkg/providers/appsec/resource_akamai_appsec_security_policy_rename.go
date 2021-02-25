package appsec

import (
	"context"
	"errors"
	"fmt"
	"strconv"
	"strings"

	"github.com/akamai/AkamaiOPEN-edgegrid-golang/v2/pkg/appsec"
	"github.com/akamai/terraform-provider-akamai/v2/pkg/akamai"
	"github.com/akamai/terraform-provider-akamai/v2/pkg/tools"
	"github.com/hashicorp/terraform-plugin-sdk/v2/diag"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/schema"
)

// appsec v1
//
// https://developer.akamai.com/api/cloud_security/application_security/v1.html
func resourceSecurityPolicyRename() *schema.Resource {
	return &schema.Resource{
		CreateContext: resourceSecurityPolicyRenameUpdate,
		ReadContext:   resourceSecurityPolicyRenameRead,
		UpdateContext: resourceSecurityPolicyRenameUpdate,
		DeleteContext: resourceSecurityPolicyRenameDelete,

		Schema: map[string]*schema.Schema{
			"config_id": {
				Type:     schema.TypeInt,
				Required: true,
			},
			"version": {
				Type:     schema.TypeInt,
				Required: true,
			},
			"security_policy_name": {
				Type:     schema.TypeString,
				Required: true,
			},
			"security_policy_id": {
				Type:     schema.TypeString,
				Required: true,
			},
		},
	}
}

func resourceSecurityPolicyRenameUpdate(ctx context.Context, d *schema.ResourceData, m interface{}) diag.Diagnostics {
	meta := akamai.Meta(m)
	client := inst.Client(meta)
	logger := meta.Log("APPSEC", "resourceSecurityPolicyUpdate")

	updateSecurityPolicy := appsec.UpdateSecurityPolicyRequest{}
	if d.Id() != "" && strings.Contains(d.Id(), ":") {
		s := strings.Split(d.Id(), ":")

		configid, errconv := strconv.Atoi(s[0])
		if errconv != nil {
			return diag.FromErr(errconv)
		}
		updateSecurityPolicy.ConfigID = configid

		version, errconv := strconv.Atoi(s[1])
		if errconv != nil {
			return diag.FromErr(errconv)
		}
		updateSecurityPolicy.Version = version

		policyid := s[2]

		updateSecurityPolicy.PolicyID = policyid

		policyname, err := tools.GetStringValue("security_policy_name", d)
		if err != nil && !errors.Is(err, tools.ErrNotFound) {
			return diag.FromErr(err)
		}
		updateSecurityPolicy.PolicyName = policyname

	} else {

		configid, err := tools.GetIntValue("config_id", d)
		if err != nil && !errors.Is(err, tools.ErrNotFound) {
			return diag.FromErr(err)
		}
		updateSecurityPolicy.ConfigID = configid

		version, err := tools.GetIntValue("version", d)
		if err != nil && !errors.Is(err, tools.ErrNotFound) {
			return diag.FromErr(err)
		}
		updateSecurityPolicy.Version = version

		securitypolicyid, err := tools.GetStringValue("security_policy_id", d)
		if err != nil && !errors.Is(err, tools.ErrNotFound) {
			return diag.FromErr(err)
		}

		updateSecurityPolicy.PolicyID = securitypolicyid

		policyname, err := tools.GetStringValue("security_policy_name", d)
		if err != nil && !errors.Is(err, tools.ErrNotFound) {
			return diag.FromErr(err)
		}
		updateSecurityPolicy.PolicyName = policyname
	}
	_, erru := client.UpdateSecurityPolicy(ctx, updateSecurityPolicy)
	if erru != nil {
		logger.Errorf("calling 'updateSecurityPolicy': %s", erru.Error())
		return diag.FromErr(erru)
	}

	return resourceSecurityPolicyRenameRead(ctx, d, m)
}

func resourceSecurityPolicyRenameDelete(ctx context.Context, d *schema.ResourceData, m interface{}) diag.Diagnostics {
	return schema.NoopContext(nil, d, m)
}

func resourceSecurityPolicyRenameRead(ctx context.Context, d *schema.ResourceData, m interface{}) diag.Diagnostics {
	meta := akamai.Meta(m)
	client := inst.Client(meta)
	logger := meta.Log("APPSEC", "resourceSecurityPolicyRead")

	getSecurityPolicy := appsec.GetSecurityPolicyRequest{}
	if d.Id() != "" && strings.Contains(d.Id(), ":") {
		s := strings.Split(d.Id(), ":")

		configid, errconv := strconv.Atoi(s[0])
		if errconv != nil {
			return diag.FromErr(errconv)
		}
		getSecurityPolicy.ConfigID = configid

		version, errconv := strconv.Atoi(s[1])
		if errconv != nil {
			return diag.FromErr(errconv)
		}
		getSecurityPolicy.Version = version

		policyid := s[2]

		getSecurityPolicy.PolicyID = policyid

	} else {
		configid, err := tools.GetIntValue("config_id", d)
		if err != nil && !errors.Is(err, tools.ErrNotFound) {
			return diag.FromErr(err)
		}
		getSecurityPolicy.ConfigID = configid

		version, err := tools.GetIntValue("version", d)
		if err != nil && !errors.Is(err, tools.ErrNotFound) {
			return diag.FromErr(err)
		}
		getSecurityPolicy.Version = version

		securitypolicyid, err := tools.GetStringValue("security_policy_id", d)
		if err != nil && !errors.Is(err, tools.ErrNotFound) {
			return diag.FromErr(err)
		}

		getSecurityPolicy.PolicyID = securitypolicyid

	}
	securitypolicy, err := client.GetSecurityPolicy(ctx, getSecurityPolicy)
	if err != nil {
		logger.Errorf("calling 'getSecurityPolicy': %s", err.Error())
		return diag.FromErr(err)
	}

	/*if err := d.Set("security_policy_name", securitypolicy.PolicyName); err != nil {
		return diag.FromErr(fmt.Errorf("%w: %s", tools.ErrValueSet, err.Error()))
	}*/

	d.SetId(fmt.Sprintf("%d:%d:%s", getSecurityPolicy.ConfigID, getSecurityPolicy.Version, securitypolicy.PolicyID))

	return nil
}
