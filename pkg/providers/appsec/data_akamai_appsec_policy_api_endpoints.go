package appsec

import (
	"context"
	"encoding/json"
	"errors"
	"fmt"
	"strconv"

	v2 "github.com/akamai/AkamaiOPEN-edgegrid-golang/v2/pkg/appsec"
	"github.com/akamai/terraform-provider-akamai/v2/pkg/akamai"
	"github.com/akamai/terraform-provider-akamai/v2/pkg/tools"

	"github.com/hashicorp/terraform-plugin-sdk/v2/diag"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/schema"
)

func dataSourcePolicyApiEndpoints() *schema.Resource {
	return &schema.Resource{
		ReadContext: dataSourceApiEndpointsRead,
		Schema: map[string]*schema.Schema{
			"config_id": {
				Type:     schema.TypeInt,
				Required: true,
			},
			"version": &schema.Schema{
				Type:     schema.TypeInt,
				Required: true,
			},
			"security_policy_id": {
				Type:     schema.TypeString,
				Required: true,
			},
			"name": {
				Type:     schema.TypeString,
				Optional: true,
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

func dataSourcePolicyApiEndpointsRead(ctx context.Context, d *schema.ResourceData, m interface{}) diag.Diagnostics {
	meta := akamai.Meta(m)
	client := inst.Client(meta)
	logger := meta.Log("APPSEC", "resourceApiEndpointsRead")

	getPolicyApiEndpoints := v2.GetPolicyApiEndpointsRequest{}

	configid, err := tools.GetIntValue("config_id", d)
	if err != nil && !errors.Is(err, tools.ErrNotFound) {
		return diag.FromErr(err)
	}
	getPolicyApiEndpoints.ConfigID = configid

	version, err := tools.GetIntValue("version", d)
	if err != nil && !errors.Is(err, tools.ErrNotFound) {
		return diag.FromErr(err)
	}
	getPolicyApiEndpoints.Version = version

	policyid, err := tools.GetStringValue("security_policy_id", d)
	if err != nil && !errors.Is(err, tools.ErrNotFound) {
		return diag.FromErr(err)
	}
	getPolicyApiEndpoints.PolicyID = policyid

	apiPolicyEndpoints, err := client.GetPolicyApiEndpoints(ctx, getPolicyApiEndpoints)
	if err != nil {
		logger.Errorf("calling 'getPolicyApiEndpoints': %s", err.Error())
		return diag.FromErr(err)
	}

	ots := OutputTemplates{}
	InitTemplates(ots)

	outputtext, err := RenderTemplates(ots, "policyApiEndpointsDS", apiPolicyEndpoints)
	if err == nil {
		d.Set("output_text", outputtext)
	}

	jsonBody, err := json.Marshal(apiPolicyEndpoints)
	if err != nil {
		return diag.FromErr(err)
	}

	if err := d.Set("json", string(jsonBody)); err != nil {
		return diag.FromErr(fmt.Errorf("%w: %s", tools.ErrValueSet, err.Error()))
	}

	if len(apiPolicyEndpoints.APIEndpoints) > 0 {

		//d.SetId(strconv.Itoa(getCustomDeny.ConfigID))
		d.SetId(strconv.Itoa(apiPolicyEndpoints.APIEndpoints[0].ID))
	}

	//d.SetId(strconv.Itoa(getApiEndpoints.ConfigID))

	return nil
}
