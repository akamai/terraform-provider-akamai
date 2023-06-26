package imaging

import (
	"context"
	"crypto/sha1"
	"encoding/hex"
	"encoding/json"
	"io"

	"github.com/akamai/AkamaiOPEN-edgegrid-golang/v7/pkg/imaging"
	"github.com/akamai/terraform-provider-akamai/v4/pkg/common/tf"
	"github.com/akamai/terraform-provider-akamai/v4/pkg/meta"
	"github.com/akamai/terraform-provider-akamai/v4/pkg/providers/imaging/imagewriter"
	"github.com/hashicorp/terraform-plugin-sdk/v2/diag"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/schema"
)

func dataImagingPolicyImage() *schema.Resource {
	return &schema.Resource{
		ReadContext: dataSourceImagingPolicyImageRead,
		Schema: map[string]*schema.Schema{
			"policy": {
				Type:        schema.TypeList,
				Required:    true,
				Description: "Image policy",
				MaxItems:    1,
				Elem: &schema.Resource{
					Schema: PolicyOutputImage(PolicyDepth),
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

func dataSourceImagingPolicyImageRead(_ context.Context, d *schema.ResourceData, m interface{}) diag.Diagnostics {
	meta := meta.Must(m)
	logger := meta.Log("Imaging", "dataPolicyImageRead")
	logger.Debug("Creating image policy json from schema")
	policy, err := tf.GetListValue("policy", d)
	if err != nil {
		return diag.FromErr(err)
	}
	var policyInput imaging.PolicyInputImage
	if policy[0] != nil {
		policyInputMap, ok := policy[0].(map[string]interface{})
		if !ok {
			return diag.Errorf("'policy' is of invalid type: %T", policyInputMap)
		}
		policyInput = imagewriter.PolicyImageToEdgeGrid(d, "policy")
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

func getPolicyJSONHashID(policyJSON string) (string, error) {
	h := sha1.New()
	_, err := io.WriteString(h, policyJSON)
	if err != nil {
		return "", err
	}
	hashID := hex.EncodeToString(h.Sum(nil))
	return hashID, nil
}
