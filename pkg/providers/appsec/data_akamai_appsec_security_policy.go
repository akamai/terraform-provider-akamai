package appsec

import (
	"context"
	"encoding/json"
	"errors"
	"fmt"

	"github.com/akamai/AkamaiOPEN-edgegrid-golang/v10/pkg/appsec"
	"github.com/akamai/terraform-provider-akamai/v6/pkg/common/tf"
	"github.com/akamai/terraform-provider-akamai/v6/pkg/meta"
	"github.com/hashicorp/terraform-plugin-sdk/v2/diag"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/schema"
)

func dataSourceSecurityPolicy() *schema.Resource {
	return &schema.Resource{
		ReadContext: dataSourceSecurityPolicyRead,
		Schema: map[string]*schema.Schema{
			"config_id": {
				Type:        schema.TypeInt,
				Required:    true,
				Description: "Unique identifier of the security configuration",
			},
			"security_policy_name": {
				Type:        schema.TypeString,
				Optional:    true,
				Description: "The name to be given to the security policy",
			},
			"security_policy_id": {
				Type:        schema.TypeString,
				Computed:    true,
				Description: "Unique identifier of the security policy",
			},
			"security_policy_id_list": {
				Type:        schema.TypeList,
				Computed:    true,
				Elem:        &schema.Schema{Type: schema.TypeString},
				Description: "List of security policy IDs",
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

func dataSourceSecurityPolicyRead(ctx context.Context, d *schema.ResourceData, m interface{}) diag.Diagnostics {
	meta := meta.Must(m)
	client := inst.Client(meta)
	logger := meta.Log("APPSEC", "dataSourceSecurityPolicyRead")

	configID, err := tf.GetIntValue("config_id", d)
	if err != nil {
		return diag.FromErr(err)
	}
	version, err := getLatestConfigVersion(ctx, configID, m)
	if err != nil {
		return diag.FromErr(err)
	}
	securityPolicyName, err := tf.GetStringValue("security_policy_name", d)
	if err != nil && !errors.Is(err, tf.ErrNotFound) {
		return diag.FromErr(err)
	}

	securityPolicies, err := client.GetSecurityPolicies(ctx, appsec.GetSecurityPoliciesRequest{
		ConfigID: configID,
		Version:  version,
	})
	if err != nil {
		logger.Errorf("calling 'getSecurityPolicies': %s", err.Error())
		return diag.FromErr(err)
	}

	jsonBody, err := json.Marshal(securityPolicies)
	if err != nil {
		return diag.FromErr(err)
	}
	if err := d.Set("json", string(jsonBody)); err != nil {
		return diag.Errorf("%s: %s", tf.ErrValueSet, err.Error())
	}

	securityPoliciesList := make([]string, 0, len(securityPolicies.Policies))
	for _, val := range securityPolicies.Policies {
		securityPoliciesList = append(securityPoliciesList, val.PolicyID)
		if val.PolicyName == securityPolicyName {
			if err := d.Set("security_policy_id", val.PolicyID); err != nil {
				return diag.Errorf("%s: %s", tf.ErrValueSet, err.Error())
			}
			if err = d.Set("security_policy_name", val.PolicyName); err != nil {
				return diag.Errorf("%s: %s", tf.ErrValueSet, err.Error())
			}
		}
	}

	if err := d.Set("security_policy_id_list", securityPoliciesList); err != nil {
		return diag.Errorf("%s: %s", tf.ErrValueSet, err.Error())
	}

	ots := OutputTemplates{}
	InitTemplates(ots)

	outputtext, err := RenderTemplates(ots, "securityPoliciesDS", securityPolicies)
	if err != nil {
		return diag.FromErr(err)
	}
	if err := d.Set("output_text", outputtext); err != nil {
		return diag.Errorf("%s: %s", tf.ErrValueSet, err.Error())
	}

	d.SetId(fmt.Sprintf("%d:%d", configID, version))

	return nil
}
