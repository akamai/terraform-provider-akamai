package imaging

import (
	"context"
	"encoding/json"
	"errors"
	"fmt"
	"net/http"
	"reflect"
	"strings"

	"github.com/akamai/AkamaiOPEN-edgegrid-golang/v3/pkg/imaging"
	"github.com/akamai/AkamaiOPEN-edgegrid-golang/v3/pkg/session"
	"github.com/akamai/terraform-provider-akamai/v2/pkg/akamai"
	"github.com/akamai/terraform-provider-akamai/v2/pkg/tools"
	"github.com/hashicorp/terraform-plugin-sdk/v2/diag"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/customdiff"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/schema"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/validation"
)

func resourceImagingPolicyVideo() *schema.Resource {
	return &schema.Resource{
		CustomizeDiff: customdiff.All(
			enforcePolicyVideoVersionChange,
		),
		CreateContext: resourcePolicyVideoCreate,
		ReadContext:   resourcePolicyVideoRead,
		UpdateContext: resourcePolicyVideoUpdate,
		DeleteContext: resourcePolicyVideoDelete,
		Schema: map[string]*schema.Schema{
			"activate_on_production": {
				Type:     schema.TypeBool,
				Optional: true,
				Default:  false,
				Description: "With this flag set to false, the user can perform modifications on staging without affecting the version already saved to production. " +
					"With this flag set to true, the policy will be saved on the production network. " +
					"It is possible to change it back to false only when there are any changes to the policy qualifying it for the new version.",
			},
			"version": {
				Type:        schema.TypeInt,
				Computed:    true,
				Description: "The version number of this policy version",
			},
			"contract_id": {
				Type:        schema.TypeString,
				Required:    true,
				ForceNew:    true,
				Description: "Unique identifier for the Akamai Contract containing the Policy Set(s)",
			},
			"policy_id": {
				Type:        schema.TypeString,
				Required:    true,
				ForceNew:    true,
				Description: "Unique identifier for a Policy. It is not possible to modify the id of the policy.",
			},
			"policyset_id": {
				Type:        schema.TypeString,
				Required:    true,
				ForceNew:    true,
				Description: "Unique identifier for the Image & Video Manager Policy Set.",
			},
			"json": {
				Type:             schema.TypeString,
				Required:         true,
				ValidateDiagFunc: validation.ToDiagFunc(validation.StringIsJSON),
				DiffSuppressFunc: diffSuppressPolicyVideo,
				Description:      "A JSON encoded policy",
			},
		},
		Importer: &schema.ResourceImporter{
			StateContext: resourcePolicyVideoImport,
		},
	}
}

func resourcePolicyVideoCreate(ctx context.Context, d *schema.ResourceData, m interface{}) diag.Diagnostics {
	meta := akamai.Meta(m)
	logger := meta.Log("Imaging", "resourcePolicyVideoCreate")
	ctx = session.ContextWithOptions(
		ctx,
		session.WithContextLog(logger),
	)
	client := inst.Client(meta)
	logger.Debug("Creating policy")

	return upsertPolicyVideo(ctx, d, m, client)
}

func upsertPolicyVideo(ctx context.Context, d *schema.ResourceData, m interface{}, client imaging.Imaging) diag.Diagnostics {
	policyID, err := tools.GetStringValue("policy_id", d)
	if err != nil {
		return diag.FromErr(err)
	}
	contractID, err := tools.GetStringValue("contract_id", d)
	if err != nil {
		return diag.FromErr(err)
	}
	policySetID, err := tools.GetStringValue("policyset_id", d)
	if err != nil {
		return diag.FromErr(err)
	}

	var policy imaging.PolicyInputVideo
	policyJSON, err := tools.GetStringValue("json", d)
	if err != nil {
		return diag.FromErr(err)
	}
	if err = json.Unmarshal([]byte(policyJSON), &policy); err != nil {
		return diag.FromErr(err)
	}

	// upsert on staging only when there was a change
	if d.HasChangesExcept("activate_on_production") {
		upsertPolicyRequest := imaging.UpsertPolicyRequest{
			PolicyID:    policyID,
			Network:     imaging.PolicyNetworkStaging,
			ContractID:  contractID,
			PolicySetID: policySetID,
			PolicyInput: &policy,
		}
		createPolicyResp, err := client.UpsertPolicy(ctx, upsertPolicyRequest)
		if err != nil {
			return diag.FromErr(err)
		}
		d.SetId(fmt.Sprintf("%s:%s", policySetID, createPolicyResp.ID))
	}
	activateOnProduction, err := tools.GetBoolValue("activate_on_production", d)
	if err != nil && !errors.Is(err, tools.ErrNotFound) {
		return diag.FromErr(err)
	}
	if activateOnProduction {
		upsertPolicyRequest := imaging.UpsertPolicyRequest{
			PolicyID:    policyID,
			Network:     imaging.PolicyNetworkProduction,
			ContractID:  contractID,
			PolicySetID: policySetID,
			PolicyInput: &policy,
		}
		_, err := client.UpsertPolicy(ctx, upsertPolicyRequest)
		if err != nil {
			return diag.FromErr(err)
		}
	}

	return resourcePolicyVideoRead(ctx, d, m)
}

func resourcePolicyVideoRead(ctx context.Context, d *schema.ResourceData, m interface{}) diag.Diagnostics {
	meta := akamai.Meta(m)
	logger := meta.Log("Imaging", "resourcePolicyVideoRead")
	ctx = session.ContextWithOptions(
		ctx,
		session.WithContextLog(logger),
	)
	client := inst.Client(meta)

	logger.Debug("Reading policy")
	policyID, err := tools.GetStringValue("policy_id", d)
	if err != nil {
		return diag.FromErr(err)
	}
	contractID, err := tools.GetStringValue("contract_id", d)
	if err != nil {
		return diag.FromErr(err)
	}
	policysetID, err := tools.GetStringValue("policyset_id", d)
	if err != nil {
		return diag.FromErr(err)
	}

	policy, err := getPolicyVideo(ctx, client, policyID, contractID, policysetID, imaging.PolicyNetworkStaging)
	if err != nil {
		return diag.FromErr(err)
	}
	policyInput, err := repackPolicyVideoOutputToInput(policy)
	if err != nil {
		return diag.FromErr(err)
	}

	policyInput.RolloutDuration, err = getRolloutDuration(d)
	if err != nil {
		return diag.FromErr(err)
	}

	policyJSON, err := getPolicyVideoJSON(policyInput)
	if err != nil {
		return diag.FromErr(err)
	}

	attrs := make(map[string]interface{})
	attrs["json"] = policyJSON
	attrs["version"] = policy.Version
	if err := tools.SetAttrs(d, attrs); err != nil {
		return diag.FromErr(err)
	}

	return nil
}

func repackPolicyVideoOutputToInput(policy *imaging.PolicyOutputVideo) (*imaging.PolicyInputVideo, error) {
	policyOutputJSON, err := json.Marshal(policy)
	if err != nil {
		return nil, err
	}
	policyInput := imaging.PolicyInputVideo{}
	err = json.Unmarshal(policyOutputJSON, &policyInput)
	if err != nil {
		return nil, err
	}
	return &policyInput, nil
}

func getPolicyVideo(ctx context.Context, client imaging.Imaging, policyID, contractID, policySetID string, network imaging.PolicyNetwork) (*imaging.PolicyOutputVideo, error) {
	policyRequest := imaging.GetPolicyRequest{
		PolicyID:    policyID,
		Network:     network,
		ContractID:  contractID,
		PolicySetID: policySetID,
	}
	policyOutput, err := client.GetPolicy(ctx, policyRequest)
	if err != nil {
		return nil, err
	}
	policy, ok := policyOutput.(*imaging.PolicyOutputVideo)
	if !ok {
		return nil, fmt.Errorf("policy is not of type video")
	}

	return policy, nil
}

func getPolicyVideoJSON(policy *imaging.PolicyInputVideo) (string, error) {
	policyJSON, err := json.MarshalIndent(policy, "", "  ")
	if err != nil {
		return "", err
	}
	return string(policyJSON), nil
}

func resourcePolicyVideoUpdate(ctx context.Context, d *schema.ResourceData, m interface{}) diag.Diagnostics {
	meta := akamai.Meta(m)
	logger := meta.Log("Imaging", "resourcePolicyVideoUpdate")
	ctx = session.ContextWithOptions(
		ctx,
		session.WithContextLog(logger),
	)
	client := inst.Client(meta)
	logger.Debug("Updating policy")

	return upsertPolicyVideo(ctx, d, m, client)
}

func resourcePolicyVideoDelete(ctx context.Context, d *schema.ResourceData, m interface{}) diag.Diagnostics {
	meta := akamai.Meta(m)
	logger := meta.Log("Imaging", "resourcePolicyVideoDelete")
	ctx = session.ContextWithOptions(
		ctx,
		session.WithContextLog(logger),
	)
	client := inst.Client(meta)
	logger.Debug("Deleting policy")

	policyID, err := tools.GetStringValue("policy_id", d)
	if err != nil {
		return diag.FromErr(err)
	}

	// .auto policy cannot be removed alone, only via removal of policy set
	if policyID == ".auto" {
		return diag.Diagnostics{
			diag.Diagnostic{
				Severity: diag.Warning,
				Summary: "Image and Video Manager API does not support '.auto' policy deletion, it can be removed when " +
					"removing related policy set - resource will only be removed from state.",
			},
		}
	}

	contractID, err := tools.GetStringValue("contract_id", d)
	if err != nil {
		return diag.FromErr(err)
	}
	policySetID, err := tools.GetStringValue("policyset_id", d)
	if err != nil {
		return diag.FromErr(err)
	}
	// delete on staging
	deletePolicyRequest := imaging.DeletePolicyRequest{
		PolicyID:    policyID,
		Network:     imaging.PolicyNetworkStaging,
		ContractID:  contractID,
		PolicySetID: policySetID,
	}
	_, err = client.DeletePolicy(ctx, deletePolicyRequest)
	if err != nil {
		return diag.FromErr(err)
	}
	// delete on production, API returns success even if policy does not exist
	// it is faster to attempt to delete on production than checking if there is policy on production first
	deletePolicyRequest = imaging.DeletePolicyRequest{
		PolicyID:    policyID,
		Network:     imaging.PolicyNetworkProduction,
		ContractID:  contractID,
		PolicySetID: policySetID,
	}
	_, err = client.DeletePolicy(ctx, deletePolicyRequest)
	if err != nil {
		return diag.FromErr(err)
	}

	d.SetId("")
	return nil
}

func resourcePolicyVideoImport(ctx context.Context, d *schema.ResourceData, m interface{}) ([]*schema.ResourceData, error) {
	meta := akamai.Meta(m)
	logger := meta.Log("Imaging", "resourcePolicyVideoImport")
	ctx = session.ContextWithOptions(
		ctx,
		session.WithContextLog(logger),
	)
	client := inst.Client(meta)

	parts := strings.Split(d.Id(), ":")
	if len(parts) != 3 {
		return nil, fmt.Errorf("colon-separated list of policy ID, policy set ID and contract ID has to be supplied in import: %s", d.Id())
	}
	policyID, policySetID, contractID := parts[0], parts[1], parts[2]

	policyStaging, err := getPolicyVideo(ctx, client, policyID, contractID, policySetID, imaging.PolicyNetworkStaging)
	if err != nil {
		return nil, err
	}
	policyInput, err := repackPolicyVideoOutputToInput(policyStaging)
	if err != nil {
		return nil, err
	}
	policyStagingJSON, err := getPolicyVideoJSON(policyInput)
	if err != nil {
		return nil, err
	}

	var activateOnProduction bool
	policyProduction, err := getPolicyVideo(ctx, client, policyID, contractID, policySetID, imaging.PolicyNetworkProduction)
	if err != nil {
		var e *imaging.Error
		if ok := errors.As(err, &e); !ok || e.Status != http.StatusNotFound {
			return nil, err
		}
	} else {
		policyInput, err := repackPolicyVideoOutputToInput(policyProduction)
		if err != nil {
			return nil, err
		}
		policyProductionJSON, err := getPolicyVideoJSON(policyInput)
		if err != nil {
			return nil, err
		}
		activateOnProduction = equalPolicyVideo(policyStagingJSON, policyProductionJSON)
	}

	attrs := make(map[string]interface{})
	attrs["policy_id"] = policyID
	attrs["contract_id"] = contractID
	attrs["policyset_id"] = policySetID
	attrs["activate_on_production"] = activateOnProduction
	if err := tools.SetAttrs(d, attrs); err != nil {
		return nil, err
	}

	d.SetId(fmt.Sprintf("%s:%s", policySetID, policyID))

	return []*schema.ResourceData{d}, nil
}

func diffSuppressPolicyVideo(_, old, new string, _ *schema.ResourceData) bool {
	return equalPolicyVideo(old, new)
}

func equalPolicyVideo(old, new string) bool {
	logger := akamai.Log("Imaging", "equalPolicyVideo")
	if old == new {
		return true
	}
	var oldPolicy, newPolicy imaging.PolicyInputVideo
	if old == "" || new == "" {
		return old == new
	}
	if err := json.Unmarshal([]byte(old), &oldPolicy); err != nil {
		logger.Errorf("Unable to unmarshal 'old' JSON policy: %s", err)
		return false
	}
	if err := json.Unmarshal([]byte(new), &newPolicy); err != nil {
		logger.Errorf("Unable to unmarshal 'new' JSON policy: %s", err)
		return false
	}

	return reflect.DeepEqual(oldPolicy, newPolicy)
}

// enforcePolicyVideoVersionChange enforces that change to any field will most likely result in creating a new version
func enforcePolicyVideoVersionChange(_ context.Context, diff *schema.ResourceDiff, _ interface{}) error {
	o, n := diff.GetChange("json")

	oldValue := o.(string)
	newValue := n.(string)

	if diff.HasChange("contract_id") ||
		diff.HasChange("policy_id") ||
		diff.HasChange("policyset_id") ||
		!equalPolicyVideo(oldValue, newValue) {
		return diff.SetNewComputed("version")
	}
	return nil
}
