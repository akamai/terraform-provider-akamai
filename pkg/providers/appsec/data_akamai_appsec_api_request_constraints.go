package appsec

import (
	"context"
	"encoding/json"
	"errors"
	"strconv"

	"github.com/akamai/AkamaiOPEN-edgegrid-golang/v2/pkg/appsec"
	"github.com/akamai/terraform-provider-akamai/v2/pkg/akamai"
	"github.com/akamai/terraform-provider-akamai/v2/pkg/tools"

	"github.com/hashicorp/terraform-plugin-sdk/v2/diag"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/schema"
)

func dataSourceAPIRequestConstraints() *schema.Resource {
	return &schema.Resource{
		ReadContext: dataSourceAPIRequestConstraintsRead,
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
			"api_id": {
				Type:        schema.TypeInt,
				Optional:    true,
				Description: "Unique identifier of the API endpoint",
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

func dataSourceAPIRequestConstraintsRead(ctx context.Context, d *schema.ResourceData, m interface{}) diag.Diagnostics {
	meta := akamai.Meta(m)
	client := inst.Client(meta)
	logger := meta.Log("APPSEC", "dataSourceAPIRequestConstraintsRead")

	getAPIiRequestConstraints := appsec.GetApiRequestConstraintsRequest{}

	configID, err := tools.GetIntValue("config_id", d)
	if err != nil {
		return diag.FromErr(err)
	}
	getAPIiRequestConstraints.ConfigID = configID

	if getAPIiRequestConstraints.Version, err = getLatestConfigVersion(ctx, configID, m); err != nil {
		return diag.FromErr(err)
	}

	policyID, err := tools.GetStringValue("security_policy_id", d)
	if err != nil {
		return diag.FromErr(err)
	}
	getAPIiRequestConstraints.PolicyID = policyID

	apiID, err := tools.GetIntValue("api_id", d)
	if err != nil && !errors.Is(err, tools.ErrNotFound) {
		return diag.FromErr(err)
	}
	getAPIiRequestConstraints.ApiID = apiID

	apirequestconstraints, err := client.GetApiRequestConstraints(ctx, getAPIiRequestConstraints)
	if err != nil {
		logger.Errorf("calling 'getApiRequestConstraints': %s", err.Error())
		return diag.FromErr(err)
	}

	ots := OutputTemplates{}
	InitTemplates(ots)

	outputtext, err := RenderTemplates(ots, "apiRequestConstraintsDS", apirequestconstraints)
	if err != nil {
		return diag.FromErr(err)
	}
	if err = d.Set("output_text", outputtext); err != nil {
		return diag.Errorf("%s: %s", tools.ErrValueSet, err.Error())
	}

	jsonBody, err := json.Marshal(apirequestconstraints)
	if err != nil {
		return diag.FromErr(err)
	}

	if err := d.Set("json", string(jsonBody)); err != nil {
		return diag.Errorf("%s: %s", tools.ErrValueSet, err.Error())
	}

	d.SetId(strconv.Itoa(getAPIiRequestConstraints.ConfigID))

	return nil
}
