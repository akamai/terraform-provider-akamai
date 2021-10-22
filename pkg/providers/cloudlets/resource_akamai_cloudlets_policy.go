package cloudlets

import (
	"context"
	"encoding/json"
	"errors"
	"fmt"
	"reflect"
	"strconv"
	"strings"

	"github.com/akamai/AkamaiOPEN-edgegrid-golang/v2/pkg/cloudlets"
	"github.com/akamai/AkamaiOPEN-edgegrid-golang/v2/pkg/session"
	"github.com/akamai/terraform-provider-akamai/v2/pkg/akamai"
	"github.com/akamai/terraform-provider-akamai/v2/pkg/tools"
	"github.com/hashicorp/terraform-plugin-sdk/v2/diag"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/schema"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/validation"
)

var cloudletIDs = map[string]int{
	"ER":  0,
	"VP":  1,
	"FR":  3,
	"IG":  4,
	"AP":  5,
	"AS":  6,
	"CD":  7,
	"IV":  8,
	"ALB": 9,
	"MMB": 10,
	"MMA": 11,
}

func resourceCloudletsPolicy() *schema.Resource {
	return &schema.Resource{
		CustomizeDiff: func(ctx context.Context, diff *schema.ResourceDiff, i interface{}) error {
			old, new := diff.GetChange("match_rules")
			if diffMatchRules(old.(string), new.(string)) {
				return nil
			}
			return diff.SetNewComputed("warnings")
		},
		CreateContext: resourcePolicyCreate,
		ReadContext:   resourcePolicyRead,
		UpdateContext: resourcePolicyUpdate,
		DeleteContext: resourcePolicyDelete,
		Schema: map[string]*schema.Schema{
			"name": {
				Type:        schema.TypeString,
				Required:    true,
				Description: "The name of the policy. The name must be unique",
			},
			"cloudlet_code": {
				Type:             schema.TypeString,
				Required:         true,
				ValidateDiagFunc: validation.ToDiagFunc(validation.StringInSlice([]string{"ALB", "ER"}, true)),
				Description:      "Code for the type of Cloudlet (ALB or ER)",
			},
			"description": {
				Type:        schema.TypeString,
				Optional:    true,
				Description: "The description of this specific policy",
			},
			"group_id": {
				Type:             schema.TypeString,
				Required:         true,
				DiffSuppressFunc: diffSuppressGroupID,
				Description:      "Defines the group association for the policy. You must have edit privileges for the group",
			},
			"match_rule_format": {
				Type:             schema.TypeString,
				Optional:         true,
				DiffSuppressFunc: diffSuppressMatchRuleFormat,
				Description:      "The version of the Cloudlet specific matchRules",
			},
			"match_rules": {
				Type:             schema.TypeString,
				Optional:         true,
				DiffSuppressFunc: diffSuppressMatchRules,
				Description:      "A JSON structure that defines the rules for this policy",
			},
			"cloudlet_id": {
				Type:        schema.TypeInt,
				Computed:    true,
				Description: "An integer that corresponds to a Cloudlets policy type (0 or 9)",
			},
			"version": {
				Type:        schema.TypeInt,
				Computed:    true,
				Description: "The version number of the policy",
			},
			"warnings": {
				Type:        schema.TypeString,
				Computed:    true,
				Description: "A JSON encoded list of warnings",
			},
		},
		Importer: &schema.ResourceImporter{
			StateContext: resourcePolicyImport,
		},
	}
}

func resourcePolicyCreate(ctx context.Context, d *schema.ResourceData, m interface{}) diag.Diagnostics {
	meta := akamai.Meta(m)
	logger := meta.Log("Cloudlets", "resourcePolicyCreate")
	// create a context with logging for api calls
	ctx = session.ContextWithOptions(
		ctx,
		session.WithContextLog(logger),
	)
	client := inst.Client(meta)
	logger.Debug("Creating policy")
	cloudletCode, err := tools.GetStringValue("cloudlet_code", d)
	if err != nil {
		return diag.FromErr(err)
	}
	cloudletID := cloudletIDs[cloudletCode]
	name, err := tools.GetStringValue("name", d)
	if err != nil {
		return diag.FromErr(err)
	}
	groupID, err := tools.GetStringValue("group_id", d)
	if err != nil {
		return diag.FromErr(err)
	}
	groupIDNum, err := tools.GetIntID(groupID, "grp_")
	if err != nil {
		return diag.Errorf("invalid group_id provided: %s", err)
	}
	createPolicyReq := cloudlets.CreatePolicyRequest{
		Name:       name,
		CloudletID: int64(cloudletID),
		GroupID:    int64(groupIDNum),
	}
	createPolicyResp, err := client.CreatePolicy(ctx, createPolicyReq)
	if err != nil {
		return diag.FromErr(err)
	}
	d.SetId(strconv.FormatInt(createPolicyResp.PolicyID, 10))
	if err := d.Set("version", 1); err != nil {
		return diag.FromErr(fmt.Errorf("%w: %s", tools.ErrValueSet, err.Error()))
	}
	matchRuleFormat, err := tools.GetStringValue("match_rule_format", d)
	if err != nil && !errors.Is(err, tools.ErrNotFound) {
		return diag.FromErr(err)
	}
	matchRulesJSON, err := tools.GetStringValue("match_rules", d)
	if err != nil {
		if errors.Is(err, tools.ErrNotFound) {
			return resourcePolicyRead(ctx, d, m)
		}
		return diag.FromErr(err)
	}
	var matchRules cloudlets.MatchRules
	if err := json.Unmarshal([]byte(matchRulesJSON), &matchRules); err != nil {
		return diag.Errorf("unmarshaling match rules JSON: %s", err)
	}
	description, err := tools.GetStringValue("description", d)
	if err != nil && !errors.Is(err, tools.ErrNotFound) {
		return diag.FromErr(err)
	}
	updateVersionRequest := cloudlets.UpdatePolicyVersionRequest{
		UpdatePolicyVersion: cloudlets.UpdatePolicyVersion{
			MatchRuleFormat: cloudlets.MatchRuleFormat(matchRuleFormat),
			MatchRules:      matchRules,
			Description:     description,
		},
		PolicyID: createPolicyResp.PolicyID,
		Version:  1,
	}

	updateVersionResp, err := client.UpdatePolicyVersion(ctx, updateVersionRequest)
	if err != nil {
		if errPolicyRead := resourcePolicyRead(ctx, d, m); errPolicyRead != nil {
			return append(errPolicyRead, diag.FromErr(err)...)
		}
		return diag.FromErr(err)
	}
	if err := setWarnings(d, updateVersionResp.Warnings); err != nil {
		return err
	}
	return resourcePolicyRead(ctx, d, m)
}

func resourcePolicyRead(ctx context.Context, d *schema.ResourceData, m interface{}) diag.Diagnostics {
	meta := akamai.Meta(m)
	logger := meta.Log("Cloudlets", "resourcePolicyRead")
	ctx = session.ContextWithOptions(
		ctx,
		session.WithContextLog(logger),
	)
	client := inst.Client(meta)
	logger.Debug("Reading policy")
	policyID, err := strconv.ParseInt(d.Id(), 10, 0)
	if err != nil {
		return diag.FromErr(err)
	}
	policy, err := client.GetPolicy(ctx, policyID)
	if err != nil {
		return diag.FromErr(err)
	}
	version, err := tools.GetIntValue("version", d)
	if err != nil {
		return diag.FromErr(err)
	}
	policyVersion, err := client.GetPolicyVersion(ctx, cloudlets.GetPolicyVersionRequest{
		PolicyID: policyID,
		Version:  int64(version),
	})
	if err != nil {
		return diag.FromErr(err)
	}
	attrs := make(map[string]interface{})
	attrs["name"] = policy.Name
	attrs["group_id"] = strconv.FormatInt(policy.GroupID, 10)
	attrs["cloudlet_code"] = policy.CloudletCode
	attrs["cloudlet_id"] = policy.CloudletID
	attrs["description"] = policyVersion.Description
	attrs["match_rule_format"] = policyVersion.MatchRuleFormat
	var matchRulesJSON []byte
	if len(policyVersion.MatchRules) > 0 {
		matchRulesJSON, err = json.MarshalIndent(policyVersion.MatchRules, "", "  ")
		if err != nil {
			return diag.FromErr(err)
		}
	}
	attrs["match_rules"] = string(matchRulesJSON)
	attrs["version"] = policyVersion.Version
	if err := tools.SetAttrs(d, attrs); err != nil {
		return diag.FromErr(err)
	}
	return nil
}

func resourcePolicyUpdate(ctx context.Context, d *schema.ResourceData, m interface{}) diag.Diagnostics {
	meta := akamai.Meta(m)
	logger := meta.Log("Cloudlets", "resourcePolicyUpdate")
	ctx = session.ContextWithOptions(
		ctx,
		session.WithContextLog(logger),
	)
	client := inst.Client(meta)
	logger.Debug("Updating policy")
	policyID, err := strconv.ParseInt(d.Id(), 10, 0)
	if err != nil {
		return diag.FromErr(err)
	}
	if d.HasChanges("name", "group_id") {
		name, err := tools.GetStringValue("name", d)
		if err != nil {
			return diag.FromErr(err)
		}
		groupID, err := tools.GetStringValue("group_id", d)
		if err != nil {
			return diag.FromErr(err)
		}
		groupIDNum, err := tools.GetIntID(groupID, "grp_")
		if err != nil {
			return diag.FromErr(err)
		}
		updatePolicyReq := cloudlets.UpdatePolicyRequest{
			UpdatePolicy: cloudlets.UpdatePolicy{
				Name:    name,
				GroupID: int64(groupIDNum),
			},
			PolicyID: policyID,
		}
		_, err = client.UpdatePolicy(ctx, updatePolicyReq)
		if err != nil {
			return diag.FromErr(err)
		}
	}
	if d.HasChanges("description", "match_rules", "match_rule_format") {
		version, err := tools.GetIntValue("version", d)
		if err != nil {
			return diag.FromErr(err)
		}
		versionResp, err := client.GetPolicyVersion(ctx, cloudlets.GetPolicyVersionRequest{
			PolicyID:  policyID,
			Version:   int64(version),
			OmitRules: true,
		})
		if err != nil {
			return diag.FromErr(err)
		}
		matchRuleFormat, err := tools.GetStringValue("match_rule_format", d)
		if err != nil && !errors.Is(err, tools.ErrNotFound) {
			return diag.FromErr(err)
		}
		matchRulesJSON, err := tools.GetStringValue("match_rules", d)
		if err != nil && !errors.Is(err, tools.ErrNotFound) {
			return diag.FromErr(err)
		}
		matchRules := make(cloudlets.MatchRules, 0)
		if matchRulesJSON != "" {
			if err := json.Unmarshal([]byte(matchRulesJSON), &matchRules); err != nil {
				return diag.FromErr(err)
			}
		}
		description, err := tools.GetStringValue("description", d)
		if err != nil && !errors.Is(err, tools.ErrNotFound) {
			return diag.FromErr(err)
		}
		if len(versionResp.Activations) > 0 {
			createVersionRequest := cloudlets.CreatePolicyVersionRequest{
				CreatePolicyVersion: cloudlets.CreatePolicyVersion{
					MatchRuleFormat: cloudlets.MatchRuleFormat(matchRuleFormat),
					MatchRules:      matchRules,
					Description:     description,
				},
				PolicyID: policyID,
			}
			createVersionResp, err := client.CreatePolicyVersion(ctx, createVersionRequest)
			if err != nil {
				return diag.FromErr(err)
			}
			if err := d.Set("version", createVersionResp.Version); err != nil {
				return diag.FromErr(fmt.Errorf("%w: %s", tools.ErrValueSet, err.Error()))
			}
			if err := setWarnings(d, createVersionResp.Warnings); err != nil {
				return err
			}
			return resourcePolicyRead(ctx, d, m)
		}
		updateVersionReq := cloudlets.UpdatePolicyVersionRequest{
			UpdatePolicyVersion: cloudlets.UpdatePolicyVersion{
				MatchRuleFormat: cloudlets.MatchRuleFormat(matchRuleFormat),
				MatchRules:      matchRules,
				Description:     description,
			},
			PolicyID: policyID,
			Version:  int64(version),
		}
		updateVersionResp, err := client.UpdatePolicyVersion(ctx, updateVersionReq)
		if err != nil {
			if errPolicyRead := resourcePolicyRead(ctx, d, m); errPolicyRead != nil {
				return append(errPolicyRead, diag.FromErr(err)...)
			}
			return diag.FromErr(err)
		}
		if err := setWarnings(d, updateVersionResp.Warnings); err != nil {
			return err
		}
	}
	return resourcePolicyRead(ctx, d, m)
}

func resourcePolicyDelete(ctx context.Context, d *schema.ResourceData, m interface{}) diag.Diagnostics {
	meta := akamai.Meta(m)
	logger := meta.Log("Cloudlets", "resourcePolicyDelete")
	ctx = session.ContextWithOptions(
		ctx,
		session.WithContextLog(logger),
	)
	client := inst.Client(meta)
	logger.Debug("Deleting policy")
	policyID, err := strconv.ParseInt(d.Id(), 10, 0)
	if err != nil {
		return diag.FromErr(err)
	}
	listVersionsResponse, err := client.ListPolicyVersions(ctx, cloudlets.ListPolicyVersionsRequest{
		PolicyID: policyID,
	})
	if err != nil {
		return diag.FromErr(err)
	}
	for _, ver := range listVersionsResponse {
		if err := client.DeletePolicyVersion(ctx, cloudlets.DeletePolicyVersionRequest{
			PolicyID: policyID,
			Version:  ver.Version,
		}); err != nil {
			return diag.FromErr(err)
		}
	}
	if err := client.RemovePolicy(ctx, policyID); err != nil {
		return diag.FromErr(err)
	}
	d.SetId("")
	return nil
}

func resourcePolicyImport(ctx context.Context, d *schema.ResourceData, m interface{}) ([]*schema.ResourceData, error) {
	meta := akamai.Meta(m)
	logger := meta.Log("Cloudlets", "resourcePolicyImport")
	logger.Debugf("Import Policy")

	client := inst.Client(meta)

	name := d.Id()
	if name == "" {
		return nil, fmt.Errorf("policy name cannot be empty")
	}

	policies, err := client.ListPolicies(ctx, cloudlets.ListPoliciesRequest{})
	if err != nil {
		return nil, err
	}
	var policy *cloudlets.Policy
	for _, p := range policies {
		if p.Name == name {
			policy = &p
			break
		}
	}
	if policy == nil {
		return nil, fmt.Errorf("could not find policy with name: %s", name)
	}

	d.SetId(strconv.FormatInt(policy.PolicyID, 10))

	versions, err := client.ListPolicyVersions(ctx, cloudlets.ListPolicyVersionsRequest{
		PolicyID: policy.PolicyID,
	})
	if err != nil {
		return nil, err
	}
	if len(versions) == 0 {
		return nil, fmt.Errorf("no policy version found")
	}
	var version int64
	for _, v := range versions {
		if v.Version > version {
			version = v.Version
		}
	}
	err = d.Set("version", version)
	if err != nil {
		return nil, err
	}

	return []*schema.ResourceData{d}, nil
}

func diffSuppressGroupID(_, old, new string, _ *schema.ResourceData) bool {
	return strings.TrimPrefix(old, "grp_") == strings.TrimPrefix(new, "grp_")
}

func diffSuppressMatchRuleFormat(_, old, new string, _ *schema.ResourceData) bool {
	return old == new || new == "" && cloudlets.MatchRuleFormat(old) == cloudlets.MatchRuleFormatDefault
}

func diffSuppressMatchRules(_, old, new string, _ *schema.ResourceData) bool {
	return diffMatchRules(old, new)
}

func diffMatchRules(old, new string) bool {
	logger := akamai.Log("Cloudlets", "diffMatchRules")
	if old == new {
		return true
	}
	var oldRules, newRules []map[string]interface{}
	if old == "" || new == "" {
		return old == new
	}
	if err := json.Unmarshal([]byte(old), &oldRules); err != nil {
		logger.Errorf("Unable to unmarshal 'old' JSON rules: %s", err)
		return false
	}
	if err := json.Unmarshal([]byte(new), &newRules); err != nil {
		logger.Errorf("Unable to unmarshal 'new' JSON rules: %s", err)
		return false
	}

	for _, rule := range oldRules {
		delete(rule, "location")
		delete(rule, "akaRuleId")
	}
	return reflect.DeepEqual(oldRules, newRules)
}

func warningsToJSON(warnings []cloudlets.Warning) ([]byte, error) {
	var warningsJSON []byte
	if len(warnings) == 0 {
		return warningsJSON, nil
	}
	
	warningsJSON, err := json.MarshalIndent(warnings, "", "  ")
	if err != nil {
		return nil, err
	}

	return warningsJSON, nil
}

func setWarnings(d *schema.ResourceData, warnings []cloudlets.Warning) diag.Diagnostics {
	warningsJSON, err := warningsToJSON(warnings)
	if err != nil {
		return diag.FromErr(err)
	}

	if err := d.Set("warnings", string(warningsJSON)); err != nil {
		return diag.FromErr(err)
	}
	return nil
}
