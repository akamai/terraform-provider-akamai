package accountprotection

import (
	"context"
	"encoding/json"
	"fmt"

	apr "github.com/akamai/AkamaiOPEN-edgegrid-golang/v12/pkg/accountprotection"
	"github.com/akamai/terraform-provider-akamai/v9/pkg/common/tf"
	"github.com/akamai/terraform-provider-akamai/v9/pkg/meta"
	"github.com/hashicorp/terraform-plugin-sdk/v2/diag"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/schema"
)

func dataSourceGeneralSettings() *schema.Resource {
	return &schema.Resource{
		ReadContext: readDataSourceGeneralSettings,
		Schema: map[string]*schema.Schema{
			"config_id": {
				Type:        schema.TypeInt,
				Required:    true,
				Description: "Identifies a security configuration.",
			},
			"security_policy_id": {
				Type:        schema.TypeString,
				Required:    true,
				Description: "Identifies a security policy.",
			},
			"json": {
				Type:     schema.TypeString,
				Computed: true,
			},
		},
	}
}

func readDataSourceGeneralSettings(ctx context.Context, d *schema.ResourceData, m interface{}) diag.Diagnostics {
	meta := meta.Must(m)
	client := inst.Client(meta)
	logger := meta.Log("accountprotection", "readDataSourceGeneralSettings")
	logger.Debugf("in readDataSourceGeneralSettings")

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

	request := apr.GetGeneralSettingsRequest{
		ConfigID:         int64(configID),
		Version:          int64(version),
		SecurityPolicyID: securityPolicyID,
	}

	response, err := client.GetGeneralSettings(ctx, request)
	if err != nil {
		logger.Errorf("calling 'GetGeneralSettings': %s", err.Error())
		return diag.FromErr(err)
	}

	jsonBody, err := json.Marshal(response)
	if err != nil {
		return diag.FromErr(err)
	}
	if err := d.Set("json", string(jsonBody)); err != nil {
		return diag.Errorf("%s: %s", tf.ErrValueSet, err.Error())
	}

	d.SetId(fmt.Sprintf("%d:%s", configID, securityPolicyID))
	return nil
}
