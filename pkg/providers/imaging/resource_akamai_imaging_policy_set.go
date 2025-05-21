package imaging

import (
	"context"
	"fmt"
	"strings"

	"github.com/akamai/AkamaiOPEN-edgegrid-golang/v10/pkg/imaging"
	"github.com/akamai/AkamaiOPEN-edgegrid-golang/v10/pkg/session"
	"github.com/akamai/terraform-provider-akamai/v8/pkg/common/tf"
	"github.com/akamai/terraform-provider-akamai/v8/pkg/meta"
	"github.com/hashicorp/terraform-plugin-sdk/v2/diag"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/schema"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/validation"
)

func resourceImagingPolicySet() *schema.Resource {
	return &schema.Resource{
		CreateContext: resourceImagingPolicySetCreate,
		ReadContext:   resourceImagingPolicySetRead,
		UpdateContext: resourceImagingPolicySetUpdate,
		DeleteContext: resourceImagingPolicySetDelete,
		Importer: &schema.ResourceImporter{
			StateContext: resourceImagingPolicySetImport,
		},
		Schema: map[string]*schema.Schema{
			"contract_id": {
				Type:             schema.TypeString,
				Required:         true,
				Description:      "The unique identifier for the Akamai Contract containing the Policy Set(s)",
				ForceNew:         true,
				DiffSuppressFunc: diffSuppressPolicySetContract,
			},
			"name": {
				Type:        schema.TypeString,
				Required:    true,
				Description: "A friendly name for the Policy Set",
			},
			"region": {
				Type:        schema.TypeString,
				Required:    true,
				Description: "The geographic region which media using this Policy Set is optimized for",
			},
			"type": {
				Type:             schema.TypeString,
				Required:         true,
				ForceNew:         true,
				Description:      "The type of media this Policy Set manages",
				ValidateDiagFunc: validation.ToDiagFunc(validation.StringInSlice([]string{string(imaging.TypeImage), string(imaging.TypeVideo)}, false)),
			},
		},
	}
}

func resourceImagingPolicySetImport(_ context.Context, rd *schema.ResourceData, m interface{}) ([]*schema.ResourceData, error) {
	meta := meta.Must(m)
	logger := meta.Log("Imaging", "resourceImagingPolicySetImport")
	logger.Debugf("Import Policy Set")

	parts := strings.Split(rd.Id(), ":")

	if len(parts) != 2 {
		return nil, fmt.Errorf("colon-separated list of policy set ID and contract ID has to be supplied in import: %s", rd.Id())
	}

	policySetID := parts[0]
	contractID := strings.TrimPrefix(parts[1], "ctr_")

	if err := rd.Set("contract_id", contractID); err != nil {
		return nil, fmt.Errorf("%w: %s", tf.ErrValueSet, err.Error())
	}

	rd.SetId(policySetID)

	return []*schema.ResourceData{rd}, nil
}

func resourceImagingPolicySetCreate(ctx context.Context, rd *schema.ResourceData, m interface{}) diag.Diagnostics {
	meta := meta.Must(m)
	logger := meta.Log("Imaging", "resourceImagingPolicySetCreate")
	ctx = session.ContextWithOptions(
		ctx,
		session.WithContextLog(logger),
	)
	client := inst.Client(meta)

	logger.Debug("Creating Policy Set")

	contractID, err := tf.GetStringValue("contract_id", rd)
	if err != nil {
		return diag.FromErr(err)
	}
	contractID = strings.TrimPrefix(contractID, "ctr_")
	name, err := tf.GetStringValue("name", rd)
	if err != nil {
		return diag.FromErr(err)
	}
	regionStr, err := tf.GetStringValue("region", rd)
	if err != nil {
		return diag.FromErr(err)
	}
	mediaTypeStr, err := tf.GetStringValue("type", rd)
	if err != nil {
		return diag.FromErr(err)
	}

	policySet, err := client.CreatePolicySet(ctx, imaging.CreatePolicySetRequest{
		ContractID: contractID,
		CreatePolicySet: imaging.CreatePolicySet{
			Name:   name,
			Region: imaging.Region(regionStr),
			Type:   imaging.MediaType(mediaTypeStr),
		},
	})
	if err != nil {
		return diag.FromErr(err)
	}

	rd.SetId(policySet.ID)

	return resourceImagingPolicySetRead(ctx, rd, m)
}

func resourceImagingPolicySetRead(ctx context.Context, rd *schema.ResourceData, m interface{}) diag.Diagnostics {
	meta := meta.Must(m)
	logger := meta.Log("Imaging", "resourceImagingPolicySetRead")
	ctx = session.ContextWithOptions(
		ctx,
		session.WithContextLog(logger),
	)
	client := inst.Client(meta)

	logger.Debugf("Reading Policy Set with ID==%s", rd.Id())

	contractID, err := tf.GetStringValue("contract_id", rd)
	if err != nil {
		return diag.FromErr(err)
	}
	contractID = strings.TrimPrefix(contractID, "ctr_")

	policySet, err := client.GetPolicySet(ctx, imaging.GetPolicySetRequest{
		PolicySetID: rd.Id(),
		ContractID:  contractID,
	})
	if err != nil {
		return diag.FromErr(err)
	}

	attrs := map[string]interface{}{
		"name":   policySet.Name,
		"region": policySet.Region,
		"type":   policySet.Type,
	}
	if err = tf.SetAttrs(rd, attrs); err != nil {
		return diag.FromErr(err)
	}

	return nil
}

func resourceImagingPolicySetUpdate(ctx context.Context, rd *schema.ResourceData, m interface{}) diag.Diagnostics {
	meta := meta.Must(m)
	logger := meta.Log("Imaging", "resourceImagingPolicySetUpdate")
	ctx = session.ContextWithOptions(
		ctx,
		session.WithContextLog(logger),
	)
	client := inst.Client(meta)

	logger.Debugf("Updating Policy Set with ID==%s", rd.Id())

	contractID, err := tf.GetStringValue("contract_id", rd)
	if err != nil {
		return diag.FromErr(err)
	}
	contractID = strings.TrimPrefix(contractID, "ctr_")
	name, err := tf.GetStringValue("name", rd)
	if err != nil {
		return diag.FromErr(err)
	}
	regionStr, err := tf.GetStringValue("region", rd)
	if err != nil {
		return diag.FromErr(err)
	}

	_, err = client.UpdatePolicySet(ctx, imaging.UpdatePolicySetRequest{
		PolicySetID: rd.Id(),
		ContractID:  contractID,
		UpdatePolicySet: imaging.UpdatePolicySet{
			Name:   name,
			Region: imaging.Region(regionStr),
		},
	})
	if err != nil {
		return diag.FromErr(err)
	}

	return resourceImagingPolicySetRead(ctx, rd, m)
}

func resourceImagingPolicySetDelete(ctx context.Context, rd *schema.ResourceData, m interface{}) diag.Diagnostics {
	meta := meta.Must(m)
	logger := meta.Log("Imaging", "resourceImagingPolicySetDelete")
	ctx = session.ContextWithOptions(
		ctx,
		session.WithContextLog(logger),
	)
	client := inst.Client(meta)

	logger.Debugf("Deleting policy set with ID==%s", rd.Id())

	contractID, err := tf.GetStringValue("contract_id", rd)
	if err != nil {
		return diag.FromErr(err)
	}
	contractID = strings.TrimPrefix(contractID, "ctr_")

	listPoliciesResponse, err := client.ListPolicies(ctx, imaging.ListPoliciesRequest{
		ContractID:  contractID,
		PolicySetID: rd.Id(),
		Network:     imaging.PolicyNetworkStaging,
	})
	if err != nil {
		return diag.FromErr(err)
	}
	remainingPolicies := filterRemainingPolicies(listPoliciesResponse)

	listPoliciesResponse, err = client.ListPolicies(ctx, imaging.ListPoliciesRequest{
		ContractID:  contractID,
		PolicySetID: rd.Id(),
		Network:     imaging.PolicyNetworkProduction,
	})
	if err != nil {
		return diag.FromErr(err)
	}
	remainingPolicies += filterRemainingPolicies(listPoliciesResponse)

	if remainingPolicies > 0 {
		return diag.Errorf("policy set with ID==%s cannot be deleted, since it contains %d associated policies", rd.Id(), remainingPolicies)
	}

	err = client.DeletePolicySet(ctx, imaging.DeletePolicySetRequest{
		PolicySetID: rd.Id(),
		ContractID:  contractID,
	})
	if err != nil {
		return diag.FromErr(err)
	}
	rd.SetId("")

	return nil
}

func filterRemainingPolicies(listPoliciesResponse *imaging.ListPoliciesResponse) int {
	var remainingPolicies int
	for _, o := range listPoliciesResponse.Items {
		var ID string
		switch policy := o.(type) {
		case *imaging.PolicyOutputImage:
			ID = policy.ID
		case *imaging.PolicyOutputVideo:
			ID = policy.ID
		default:
			panic("unsupported policy output type")
		}
		if ID != ".auto" {
			// there is always one policy with ID==".auto"
			remainingPolicies++
		}
	}
	return remainingPolicies
}

func diffSuppressPolicySetContract(_, o, n string, _ *schema.ResourceData) bool {
	return strings.TrimPrefix(o, "ctr_") == strings.TrimPrefix(n, "ctr_")
}
