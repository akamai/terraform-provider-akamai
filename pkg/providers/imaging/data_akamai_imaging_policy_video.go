package imaging

import (
	"context"
	"encoding/json"

	"github.com/akamai/AkamaiOPEN-edgegrid-golang/v10/pkg/imaging"
	"github.com/akamai/terraform-provider-akamai/v8/pkg/common/tf"
	"github.com/akamai/terraform-provider-akamai/v8/pkg/meta"
	"github.com/akamai/terraform-provider-akamai/v8/pkg/providers/imaging/videowriter"
	"github.com/hashicorp/terraform-plugin-sdk/v2/diag"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/schema"
)

func dataImagingPolicyVideo() *schema.Resource {
	return &schema.Resource{
		ReadContext: dataSourceImagingPolicyVideoRead,
		Schema: map[string]*schema.Schema{
			"policy": {
				Type:        schema.TypeList,
				Required:    true,
				Description: "Video policy",
				MaxItems:    1,
				Elem: &schema.Resource{
					Schema: PolicyOutputVideo(PolicyDepth),
				},
			},
			"json": {
				Type:        schema.TypeString,
				Computed:    true,
				Description: "A JSON encoded policy",
			},
		},
	}
}

func dataSourceImagingPolicyVideoRead(_ context.Context, d *schema.ResourceData, m interface{}) diag.Diagnostics {
	meta := meta.Must(m)
	logger := meta.Log("Imaging", "dataPolicyVideoRead")
	logger.Debug("Creating video policy json from schema")
	policy, err := tf.GetListValue("policy", d)
	if err != nil {
		return diag.FromErr(err)
	}
	var policyInput imaging.PolicyInputVideo
	if policy[0] != nil {
		policyInputMap, ok := policy[0].(map[string]interface{})
		if !ok {
			return diag.Errorf("'policy' is of invalid type: %T", policyInputMap)
		}
		policyInput = videowriter.PolicyVideoToEdgeGrid(d, "policy")
	}

	jsonBody, err := json.MarshalIndent(policyInput, "", "  ")
	if err != nil {
		return diag.FromErr(err)
	}
	if err := d.Set("json", string(jsonBody)); err != nil {
		return diag.Errorf("%v: %s", tf.ErrValueSet, err.Error())
	}

	id, err := getPolicyJSONHashID(string(jsonBody))
	if err != nil {
		return diag.Errorf("calculating hash ID from policy JSON: %s", err.Error())
	}
	d.SetId(id)
	return nil
}
