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
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/schema"
)

func dataSourceEval() *schema.Resource {
	return &schema.Resource{
		ReadContext: dataSourceEvalRead,
		Schema: map[string]*schema.Schema{
			"config_id": {
				Type:     schema.TypeInt,
				Required: true,
			},
			"version": {
				Type:     schema.TypeInt,
				Required: true,
			},
			"security_policy_id": {
				Type:     schema.TypeString,
				Required: true,
			},
			"output_text": {
				Type:        schema.TypeString,
				Computed:    true,
				Description: "Text Export representation",
			},
		},
	}
}

func dataSourceEvalRead(ctx context.Context, d *schema.ResourceData, m interface{}) diag.Diagnostics {
	meta := akamai.Meta(m)
	client := inst.Client(meta)
	logger := meta.Log("APPSEC", "dataSourceEvalRead")

	getEval := appsec.GetEvalRequest{}

	configid, err := tools.GetIntValue("config_id", d)
	if err != nil && !errors.Is(err, tools.ErrNotFound) {
		return diag.FromErr(err)
	}
	getEval.ConfigID = configid

	version, err := tools.GetIntValue("version", d)
	if err != nil && !errors.Is(err, tools.ErrNotFound) {
		return diag.FromErr(err)
	}
	getEval.Version = version

	policyid, err := tools.GetStringValue("security_policy_id", d)
	if err != nil && !errors.Is(err, tools.ErrNotFound) {
		return diag.FromErr(err)
	}
	getEval.PolicyID = policyid

	eval, err := client.GetEval(ctx, getEval)
	if err != nil {
		logger.Errorf("calling 'getEval': %s", err.Error())
		return diag.FromErr(err)
	}

	ots := OutputTemplates{}
	InitTemplates(ots)

	outputtext, err := RenderTemplates(ots, "EvalDS", eval)
	if err == nil {
		if err := d.Set("output_text", outputtext); err != nil {
			return diag.FromErr(fmt.Errorf("%w: %s", tools.ErrValueSet, err.Error()))
		}
	}

	d.SetId(strconv.Itoa(getEval.ConfigID))

	return nil
}
