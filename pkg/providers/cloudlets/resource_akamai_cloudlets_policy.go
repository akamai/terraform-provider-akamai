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
		CreateContext: resourcePolicyCreate,
		ReadContext:   resourcePolicyRead,
		UpdateContext: resourcePolicyUpdate,
		DeleteContext: resourcePolicyDelete,
		Schema: map[string]*schema.Schema{
			"name": {
				Type:     schema.TypeString,
				Required: true,
			},
			"cloudlet_code": {
				Type:             schema.TypeString,
				Required:         true,
				ValidateDiagFunc: validation.ToDiagFunc(validation.StringInSlice([]string{"ALB", "ER"}, true)),
			},
			"description": {
				Type:     schema.TypeString,
				Optional: true,
			},
			"group_id": {
				Type:             schema.TypeString,
				Required:         true,
				DiffSuppressFunc: diffSuppressGroupID,
			},
			"match_rule_format": {
				Type:             schema.TypeString,
				Optional:         true,
				DiffSuppressFunc: diffSuppressMatchRuleFormat,
			},
			"match_rules": {
				Type:             schema.TypeString,
				Optional:         true,
				DiffSuppressFunc: diffSuppressMatchRules,
			},
			"cloudlet_id": {
				Type:     schema.TypeInt,
				Computed: true,
			},
			"version": {
				Type:     schema.TypeInt,
				Computed: true,
			},
			"warnings": {
				Type:     schema.TypeString,
				Computed: true,
			},
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
	description, err := tools.GetStringValue("description", d)
	if err != nil && !errors.Is(err, tools.ErrNotFound) {
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
		Name:        name,
		CloudletID:  int64(cloudletID),
		Description: description,
		GroupID:     int64(groupIDNum),
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
	updateVersionRequest := cloudlets.UpdatePolicyVersionRequest{
		UpdatePolicyVersion: cloudlets.UpdatePolicyVersion{
			MatchRuleFormat: cloudlets.MatchRuleFormat(matchRuleFormat),
			MatchRules:      matchRules,
		},
		PolicyID: createPolicyResp.PolicyID,
		Version:  1,
	}

	if _, err = client.UpdatePolicyVersion(ctx, updateVersionRequest); err != nil {
		return diag.FromErr(err)
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
	attrs["description"] = policy.Description
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
	var warningsJSON []byte
	if len(policyVersion.Warnings) > 0 {
		warningsJSON, err = json.MarshalIndent(policyVersion.Warnings, "", "  ")
		if err != nil {
			return diag.FromErr(err)
		}
	}
	attrs["warnings"] = string(warningsJSON)
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
	if d.HasChanges("name", "description", "group_id") {
		name, err := tools.GetStringValue("name", d)
		if err != nil {
			return diag.FromErr(err)
		}
		description, err := tools.GetStringValue("description", d)
		if err != nil && !errors.Is(err, tools.ErrNotFound) {
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
				Name:        name,
				Description: description,
				GroupID:     int64(groupIDNum),
			},
			PolicyID: policyID,
		}
		_, err = client.UpdatePolicy(ctx, updatePolicyReq)
		if err != nil {
			return diag.FromErr(err)
		}
	}
	if d.HasChanges("match_rules", "match_rule_format") {
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
		if len(versionResp.Activations) > 0 {
			createVersionRequest := cloudlets.CreatePolicyVersionRequest{
				CreatePolicyVersion: cloudlets.CreatePolicyVersion{
					MatchRuleFormat: cloudlets.MatchRuleFormat(matchRuleFormat),
					MatchRules:      matchRules,
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
			return resourcePolicyRead(ctx, d, m)
		}
		updateVersionReq := cloudlets.UpdatePolicyVersionRequest{
			UpdatePolicyVersion: cloudlets.UpdatePolicyVersion{
				MatchRuleFormat: cloudlets.MatchRuleFormat(matchRuleFormat),
				MatchRules:      matchRules,
			},
			PolicyID: policyID,
			Version:  int64(version),
		}
		_, err = client.UpdatePolicyVersion(ctx, updateVersionReq)
		if err != nil {
			return diag.FromErr(err)
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

func diffSuppressGroupID(_, old, new string, _ *schema.ResourceData) bool {
	return strings.TrimPrefix(old, "grp_") == strings.TrimPrefix(new, "grp_")
}

func diffSuppressMatchRuleFormat(_, old, new string, _ *schema.ResourceData) bool {
	return old == new || new == "" && cloudlets.MatchRuleFormat(old) == cloudlets.MatchRuleFormatDefault
}

func diffSuppressMatchRules(_, old, new string, _ *schema.ResourceData) bool {
	logger := akamai.Log("Cloudlets", "diffSuppressMatchRules")
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
