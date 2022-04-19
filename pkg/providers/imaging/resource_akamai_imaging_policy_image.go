package imaging

import (
	"context"
	"encoding/json"
	"errors"
	"fmt"
	"net/http"
	"reflect"
	"strings"

	"github.com/akamai/AkamaiOPEN-edgegrid-golang/v2/pkg/imaging"
	"github.com/akamai/AkamaiOPEN-edgegrid-golang/v2/pkg/session"
	"github.com/akamai/terraform-provider-akamai/v2/pkg/akamai"
	"github.com/akamai/terraform-provider-akamai/v2/pkg/providers/imaging/reader"
	"github.com/akamai/terraform-provider-akamai/v2/pkg/providers/imaging/writer"
	"github.com/akamai/terraform-provider-akamai/v2/pkg/tools"
	"github.com/hashicorp/terraform-plugin-sdk/v2/diag"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/customdiff"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/schema"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/validation"
)

// MaxPolicyDepth is the maximum supported depth of the nested transformations, before we reach hard limit of the Terraform
// provider for the gRPC communication for exchanging schema information, which is 64MB
const MaxPolicyDepth = 7

// PolicyDepth is variable to allow changing it for the unit tests, to achieve faster execution for tests,
// which do not need all supported levels
var PolicyDepth = MaxPolicyDepth

func resourceImagingPolicyImage() *schema.Resource {
	return &schema.Resource{
		CustomizeDiff: customdiff.All(
			enforcePolicyImageVersionChange,
		),
		CreateContext: resourcePolicyImageCreate,
		ReadContext:   resourcePolicyImageRead,
		UpdateContext: resourcePolicyImageUpdate,
		DeleteContext: resourcePolicyImageDelete,
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
				Optional:         true,
				ExactlyOneOf:     []string{"json", "policy"},
				ValidateDiagFunc: validation.ToDiagFunc(validation.StringIsJSON),
				DiffSuppressFunc: diffSuppressPolicyImage,
				Description:      "A JSON encoded policy",
			},
			"policy": {
				Type:        schema.TypeList,
				Optional:    true,
				Description: "Policy",
				MaxItems:    1,
				Elem: &schema.Resource{
					Schema: policyImage(PolicyDepth),
				},
			},
		},
		Importer: &schema.ResourceImporter{
			StateContext: resourcePolicyImageImport,
		},
	}
}

func resourcePolicyImageCreate(ctx context.Context, d *schema.ResourceData, m interface{}) diag.Diagnostics {
	meta := akamai.Meta(m)
	logger := meta.Log("Imaging", "resourcePolicyImageCreate")
	ctx = session.ContextWithOptions(
		ctx,
		session.WithContextLog(logger),
	)
	client := inst.Client(meta)
	logger.Debug("Creating policy")

	return upsertPolicy(ctx, d, m, client)
}

func upsertPolicy(ctx context.Context, d *schema.ResourceData, m interface{}, client imaging.Imaging) diag.Diagnostics {
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

	var policy imaging.PolicyInputImage
	policyJSON, err := tools.GetStringValue("json", d)
	if err != nil && !errors.Is(err, tools.ErrNotFound) {
		return diag.FromErr(err)
	}
	if err == nil {
		if err = json.Unmarshal([]byte(policyJSON), &policy); err != nil {
			return diag.FromErr(err)
		}
	} else {
		policy = writer.PolicyImageToEdgeGrid(d.Get("policy").([]interface{})[0].(map[string]interface{}))
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

	return resourcePolicyImageRead(ctx, d, m)
}

func resourcePolicyImageRead(ctx context.Context, d *schema.ResourceData, m interface{}) diag.Diagnostics {
	meta := akamai.Meta(m)
	logger := meta.Log("Imaging", "resourcePolicyImageRead")
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

	policy, err := getPolicyImage(ctx, client, policyID, contractID, policysetID, imaging.PolicyNetworkStaging)
	if err != nil {
		return diag.FromErr(err)
	}

	attrs := make(map[string]interface{})
	_, ok := d.GetOk("policy")
	if ok {
		attrs["policy"] = []interface{}{reader.GetImageSchema(*policy)}
	} else {
		policyJSON, err := getPolicyImageJSON(policy)
		if err != nil {
			return diag.FromErr(err)
		}
		attrs["json"] = policyJSON
	}
	attrs["version"] = policy.Version
	if err := tools.SetAttrs(d, attrs); err != nil {
		return diag.FromErr(err)
	}

	return nil
}

func getPolicyImage(ctx context.Context, client imaging.Imaging, policyID, contractID, policySetID string, network imaging.PolicyNetwork) (*imaging.PolicyOutputImage, error) {
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
	policy, ok := policyOutput.(*imaging.PolicyOutputImage)
	if !ok {
		return nil, fmt.Errorf("policy is not of type image")
	}

	return policy, nil
}

func getPolicyImageJSON(policy *imaging.PolicyOutputImage) (string, error) {
	policyJSON, err := json.MarshalIndent(policy, "", "  ")
	if err != nil {
		return "", err
	}

	// we store JSON as PolicyInput, so we need to convert it from PolicyOutput via JSON representation
	var policyInput imaging.PolicyInputImage
	if err := json.Unmarshal(policyJSON, &policyInput); err != nil {
		return "", err
	}

	policyJSON, err = json.MarshalIndent(policyInput, "", "  ")
	if err != nil {
		return "", err
	}

	return string(policyJSON), nil
}

func resourcePolicyImageUpdate(ctx context.Context, d *schema.ResourceData, m interface{}) diag.Diagnostics {
	meta := akamai.Meta(m)
	logger := meta.Log("Imaging", "resourcePolicyImageUpdate")
	ctx = session.ContextWithOptions(
		ctx,
		session.WithContextLog(logger),
	)
	client := inst.Client(meta)
	logger.Debug("Updating policy")

	return upsertPolicy(ctx, d, m, client)
}

func resourcePolicyImageDelete(ctx context.Context, d *schema.ResourceData, m interface{}) diag.Diagnostics {
	meta := akamai.Meta(m)
	logger := meta.Log("Imaging", "resourcePolicyImageDelete")
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
				Summary:  "Image and Video Manager API does not support '.auto' policy deletion, it can be removed when removing related policy set - resource will only be removed from state.",
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

func resourcePolicyImageImport(ctx context.Context, d *schema.ResourceData, m interface{}) ([]*schema.ResourceData, error) {
	meta := akamai.Meta(m)
	logger := meta.Log("Imaging", "resourcePolicyImageImport")
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

	policyStaging, err := getPolicyImage(ctx, client, policyID, contractID, policySetID, imaging.PolicyNetworkStaging)
	if err != nil {
		return nil, err
	}
	policyStagingJSON, err := getPolicyImageJSON(policyStaging)
	if err != nil {
		return nil, err
	}

	var activateOnProduction bool
	policyProduction, err := getPolicyImage(ctx, client, policyID, contractID, policySetID, imaging.PolicyNetworkProduction)
	if err != nil {
		var e *imaging.Error
		if ok := errors.As(err, &e); !ok || e.Status != http.StatusNotFound {
			return nil, err
		}
	} else {
		policyProductionJSON, err := getPolicyImageJSON(policyProduction)
		if err != nil {
			return nil, err
		}
		activateOnProduction = equalPolicyImage(policyStagingJSON, policyProductionJSON)
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

func diffSuppressPolicyImage(_, old, new string, _ *schema.ResourceData) bool {
	return equalPolicyImage(old, new)
}

func equalPolicyImage(old, new string) bool {
	logger := akamai.Log("Imaging", "equalPolicyImage")
	if old == new {
		return true
	}
	var oldPolicy, newPolicy imaging.PolicyInputImage
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

// enforcePolicyImageVersionChange enforces that change to any field will most likely result in creating a new version
func enforcePolicyImageVersionChange(_ context.Context, diff *schema.ResourceDiff, _ interface{}) error {
	o, n := diff.GetChange("json")

	oldValue := o.(string)
	newValue := n.(string)

	if diff.HasChange("contract_id") ||
		diff.HasChange("policy_id") ||
		diff.HasChange("policyset_id") ||
		!equalPolicyImage(oldValue, newValue) {
		return diff.SetNewComputed("version")
	}
	return nil
}
