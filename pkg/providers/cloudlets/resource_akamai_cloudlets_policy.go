package cloudlets

import (
	"context"
	"encoding/json"
	"errors"
	"fmt"
	"reflect"
	"strconv"
	"strings"
	"time"

	"github.com/akamai/AkamaiOPEN-edgegrid-golang/v7/pkg/cloudlets"
	v3 "github.com/akamai/AkamaiOPEN-edgegrid-golang/v7/pkg/cloudlets/v3"
	"github.com/akamai/AkamaiOPEN-edgegrid-golang/v7/pkg/session"

	"github.com/akamai/terraform-provider-akamai/v5/pkg/common/str"
	"github.com/akamai/terraform-provider-akamai/v5/pkg/common/tf"
	"github.com/akamai/terraform-provider-akamai/v5/pkg/common/timeouts"
	"github.com/akamai/terraform-provider-akamai/v5/pkg/logger"
	"github.com/akamai/terraform-provider-akamai/v5/pkg/meta"

	"github.com/hashicorp/terraform-plugin-sdk/v2/diag"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/customdiff"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/schema"
)

var (
	// DeletionPolicyPollInterval is the default poll interval for delete policy retries
	DeletionPolicyPollInterval = time.Second * 10

	// DeletionPolicyTimeout is the default timeout for the policy deletion
	DeletionPolicyTimeout = time.Minute * 90
)

func resourceCloudletsPolicy() *schema.Resource {
	return &schema.Resource{
		CustomizeDiff: customdiff.All(
			EnforcePolicyVersionChange,
			EnforceMatchRulesChange,
			cloudletTypeChangesValidation,
			cloudletCodeValidation,
			cloudletCodeChangeValidation,
		),
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
				Type:        schema.TypeString,
				Required:    true,
				Description: "Code for the type of Cloudlet (ALB, AP, AS, CD, ER, FR, IG, or VP)",
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
			"is_shared": {
				Type:        schema.TypeBool,
				Optional:    true,
				Default:     false,
				Description: "The type of policy that you want to create",
			},
			"cloudlet_id": {
				Type:        schema.TypeInt,
				Computed:    true,
				Description: "An integer that corresponds to a non-shared Cloudlets policy type (0 to 9). Not used for shared policies",
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
			"timeouts": {
				Type:        schema.TypeList,
				Optional:    true,
				MaxItems:    1,
				Description: "Enables to set timeout for processing",
				Elem: &schema.Resource{
					Schema: map[string]*schema.Schema{
						"default": {
							Type:             schema.TypeString,
							Optional:         true,
							ValidateDiagFunc: timeouts.ValidateDurationFormat,
						},
					},
				},
			},
		},
		Importer: &schema.ResourceImporter{
			StateContext: resourcePolicyImport,
		},
		Timeouts: &schema.ResourceTimeout{
			Default: &DeletionPolicyTimeout,
		},
	}
}

func cloudletCodeChangeValidation(_ context.Context, diff *schema.ResourceDiff, _ interface{}) error {
	if diff.Id() != "" && diff.HasChange("cloudlet_code") {
		return fmt.Errorf("cloudlet code cannot be changed after creation, please destroy policy and create new one with modified `cloudlet_code`")
	}
	return nil
}

func cloudletCodeValidation(_ context.Context, diff *schema.ResourceDiff, _ interface{}) error {
	isShared := diff.Get("is_shared").(bool)
	providedCode := diff.Get("cloudlet_code").(string)
	if isShared {
		possibleValues := []string{"AP", "AS", "CD", "ER", "FR", "IG"}
		for _, code := range possibleValues {
			if strings.ToLower(providedCode) == strings.ToLower(code) {
				return nil
			}
		}
		return fmt.Errorf("provided cloudlet code %s cannot be used in shared policy - use one of %s", providedCode, possibleValues)
	}

	possibleValues := []string{"ALB", "AP", "AS", "CD", "ER", "FR", "IG", "VP"}
	for _, code := range possibleValues {
		if strings.ToLower(providedCode) == strings.ToLower(code) {
			return nil
		}
	}
	return fmt.Errorf("provided cloudlet code %s cannot be used in legacy policy - use one of %s", providedCode, possibleValues)
}

// cloudletTypeChangesValidation is used to run validation for v2 -> v3 (or vice versa) related migrations
func cloudletTypeChangesValidation(_ context.Context, diff *schema.ResourceDiff, _ interface{}) error {
	if diff.Id() != "" {
		if diff.HasChange("is_shared") {
			return fmt.Errorf("it is impossible to convert shared cloudlet to legacy one or vice versa; create new policy with modified named for target policy type")
		}
		if diff.Get("is_shared").(bool) && diff.HasChange("name") {
			return fmt.Errorf("it is impossible to rename shared policy")
		}
	}

	return nil
}

// EnforcePolicyVersionChange enforces that change to any field will most likely result in creating a new version
func EnforcePolicyVersionChange(_ context.Context, diff *schema.ResourceDiff, _ interface{}) error {
	if diff.HasChanges("name", "description", "match_rule_format", "version") {
		return diff.SetNewComputed("version")
	}
	return nil
}

// EnforceMatchRulesChange enforces that any changes to match_rules will re compute the warnings
func EnforceMatchRulesChange(_ context.Context, diff *schema.ResourceDiff, _ interface{}) error {
	oldVal, newVal := diff.GetChange("match_rules")
	if diffMatchRules(oldVal.(string), newVal.(string)) {
		return nil
	}
	if err := diff.SetNewComputed("warnings"); err != nil {
		return err
	}
	return diff.SetNewComputed("version")
}

func resourcePolicyCreate(ctx context.Context, d *schema.ResourceData, m interface{}) diag.Diagnostics {
	meta := meta.Must(m)
	logger := meta.Log("Cloudlets", "resourcePolicyCreate")
	// create a context with logging for api calls
	ctx = session.ContextWithOptions(
		ctx,
		session.WithContextLog(logger),
	)
	logger.Debug("Creating policy")
	cloudletCode, err := tf.GetStringValue("cloudlet_code", d)
	if err != nil {
		return diag.FromErr(err)
	}
	name, err := tf.GetStringValue("name", d)
	if err != nil {
		return diag.FromErr(err)
	}
	groupID, err := tf.GetStringValue("group_id", d)
	if err != nil {
		return diag.FromErr(err)
	}
	groupIDNum, err := str.GetIntID(groupID, "grp_")
	if err != nil {
		return diag.Errorf("invalid group_id provided: %s", err)
	}

	executionStrategy, err := getPolicyExecutionStrategy(d, meta)
	if err != nil {
		return diag.FromErr(err)
	}

	policyID, err := executionStrategy.createPolicy(ctx, name, cloudletCode, int64(groupIDNum))
	if err != nil {
		return diag.FromErr(err)
	}

	d.SetId(strconv.FormatInt(policyID, 10))

	description, err := tf.GetStringValue("description", d)
	if err != nil && !errors.Is(err, tf.ErrNotFound) {
		return diag.FromErr(err)
	}

	matchRulesJSON, err := tf.GetStringValue("match_rules", d)
	if err != nil {
		if errors.Is(err, tf.ErrNotFound) {
			if description == "" {
				return resourcePolicyRead(ctx, d, m)
			}
		} else {
			return diag.FromErr(err)
		}
	}

	err, updateError := executionStrategy.updatePolicyVersion(ctx, d, policyID, 1, description, matchRulesJSON, !executionStrategy.isFirstVersionCreated())
	if err != nil {
		return diag.FromErr(err)
	}

	if updateError != nil {
		// The resource will be created as tainted (because the setId was executed). So on next plan it'll delete it and create again.
		// We still want to have actual (server's) values in state. Otherwise, the values from config would be put into the state as default.
		if errPolicyRead := resourcePolicyRead(ctx, d, m); errPolicyRead != nil {
			return append(errPolicyRead, diag.FromErr(updateError)...)
		}
		return diag.FromErr(updateError)
	}
	return resourcePolicyRead(ctx, d, m)
}

func resourcePolicyRead(ctx context.Context, d *schema.ResourceData, m interface{}) diag.Diagnostics {
	meta := meta.Must(m)
	logger := meta.Log("Cloudlets", "resourcePolicyRead")
	ctx = session.ContextWithOptions(
		ctx,
		session.WithContextLog(logger),
	)
	logger.Debug("Reading policy")
	policyID, err := strconv.ParseInt(d.Id(), 10, 0)
	if err != nil {
		return diag.FromErr(err)
	}

	policyVersionStrategy, err := getPolicyVersionExecutionStrategy(d, meta)
	if err != nil {
		return diag.FromErr(err)
	}

	policyVersion, err := policyVersionStrategy.findLatestPolicyVersion(ctx, policyID)
	if err != nil {
		return diag.FromErr(err)
	}

	executionStrategy, err := getPolicyExecutionStrategy(d, meta)
	if err != nil {
		return diag.FromErr(err)
	}

	attrs, err := executionStrategy.readPolicy(ctx, policyID, policyVersion)
	if err != nil {
		return diag.FromErr(err)
	}

	if err := tf.SetAttrs(d, attrs); err != nil {
		return diag.FromErr(err)
	}
	return nil
}

func resourcePolicyUpdate(ctx context.Context, d *schema.ResourceData, m interface{}) diag.Diagnostics {
	meta := meta.Must(m)
	logger := meta.Log("Cloudlets", "resourcePolicyUpdate")
	ctx = session.ContextWithOptions(
		ctx,
		session.WithContextLog(logger),
	)
	logger.Debug("Updating policy")

	if !d.HasChangeExcept("timeouts") {
		logger.Debug("Only timeouts were updated, skipping")
		return nil
	}

	executionStrategy, err := getPolicyExecutionStrategy(d, meta)
	if err != nil {
		return diag.FromErr(err)
	}

	policyID, err := strconv.ParseInt(d.Id(), 10, 0)
	if err != nil {
		return diag.FromErr(err)
	}
	if d.HasChanges("name", "group_id") {
		name, err := tf.GetStringValue("name", d)
		if err != nil {
			return diag.FromErr(err)
		}
		groupID, err := tf.GetStringValue("group_id", d)
		if err != nil {
			return diag.FromErr(err)
		}
		groupIDNum, err := str.GetIntID(groupID, "grp_")
		if err != nil {
			return diag.FromErr(err)
		}

		err = executionStrategy.updatePolicy(ctx, policyID, int64(groupIDNum), name)
		if err != nil {
			return diag.FromErr(err)
		}
	}
	if d.HasChanges("description", "match_rules", "match_rule_format") {
		var isNewVersionNeeded bool
		version, err := tf.GetIntValue("version", d)
		if err != nil {
			if errors.Is(err, tf.ErrNotFound) {
				isNewVersionNeeded = true
			} else {
				return diag.FromErr(err)
			}
		} else {
			isNewVersionNeeded, err = executionStrategy.newPolicyVersionIsNeeded(ctx, policyID, int64(version))
			if err != nil {
				return diag.FromErr(err)
			}
		}
		matchRulesJSON, err := tf.GetStringValue("match_rules", d)
		if err != nil && !errors.Is(err, tf.ErrNotFound) {
			return diag.FromErr(err)
		}
		description, err := tf.GetStringValue("description", d)
		if err != nil && !errors.Is(err, tf.ErrNotFound) {
			return diag.FromErr(err)
		}

		err, updateVersionErr := executionStrategy.updatePolicyVersion(ctx, d, policyID, int64(version), description, matchRulesJSON, isNewVersionNeeded)
		if err != nil {
			return diag.FromErr(err)
		}

		if updateVersionErr != nil {
			// We still want to have actual (server's) values in state. Otherwise, the values from config would be put into the state as default.
			if errPolicyRead := resourcePolicyRead(ctx, d, m); errPolicyRead != nil {
				return append(errPolicyRead, diag.FromErr(updateVersionErr)...)
			}
			return diag.FromErr(updateVersionErr)
		}
	}
	return resourcePolicyRead(ctx, d, m)
}

func resourcePolicyDelete(ctx context.Context, d *schema.ResourceData, m interface{}) diag.Diagnostics {
	meta := meta.Must(m)
	logger := meta.Log("Cloudlets", "resourcePolicyDelete")
	ctx = session.ContextWithOptions(
		ctx,
		session.WithContextLog(logger),
	)
	logger.Debug("Deleting policy")

	executionStrategy, err := getPolicyExecutionStrategy(d, meta)
	if err != nil {
		return diag.FromErr(err)
	}

	policyID, err := strconv.ParseInt(d.Id(), 10, 0)
	if err != nil {
		return diag.FromErr(err)
	}

	err = executionStrategy.deletePolicy(ctx, policyID)
	if err != nil {
		return diag.FromErr(err)
	}
	d.SetId("")

	return nil
}

func resourcePolicyImport(ctx context.Context, d *schema.ResourceData, m interface{}) ([]*schema.ResourceData, error) {
	meta := meta.Must(m)
	logger := meta.Log("Cloudlets", "resourcePolicyImport")
	logger.Debugf("Import Policy")

	name := d.Id()
	if name == "" {
		return nil, fmt.Errorf("policy name cannot be empty")
	}

	policyStrategy, policyID, err := discoverPolicyExecutionStrategy(ctx, meta, name)
	if err != nil {
		return nil, err
	}
	if err = policyStrategy.setPolicyType(d); err != nil {
		return nil, err
	}

	d.SetId(strconv.FormatInt(policyID, 10))

	return []*schema.ResourceData{d}, nil
}

func diffSuppressGroupID(_, old, new string, _ *schema.ResourceData) bool {
	return strings.TrimPrefix(old, "grp_") == strings.TrimPrefix(new, "grp_")
}

func diffSuppressMatchRuleFormat(_, old, new string, _ *schema.ResourceData) bool {
	return old == new || new == "" && cloudlets.MatchRuleFormat(old) == cloudlets.MatchRuleFormat10
}

func diffSuppressMatchRules(_, old, new string, _ *schema.ResourceData) bool {
	return diffMatchRules(old, new)
}

func diffMatchRules(old, new string) bool {
	logger := logger.Get("Cloudlets", "diffMatchRules")
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

func warningsToJSON[W cloudlets.Warning | v3.MatchRulesWarning](warnings []W) ([]byte, error) {
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

func setWarnings[W cloudlets.Warning | v3.MatchRulesWarning](d *schema.ResourceData, warnings []W) error {
	warningsJSON, err := warningsToJSON(warnings)
	if err != nil {
		return err
	}

	return d.Set("warnings", string(warningsJSON))
}

func getPolicyExecutionStrategy(d *schema.ResourceData, meta meta.Meta) (policyExecutionStrategy, error) {
	var executionStrategy policyExecutionStrategy
	isV3, err := tf.GetBoolValue("is_shared", d)
	if err != nil {
		return nil, err
	}

	if isV3 {
		executionStrategy = v3PolicyStrategy{ClientV3(meta)}
	} else {
		executionStrategy = v2PolicyStrategy{Client(meta)}
	}
	return executionStrategy, nil
}

type policyExecutionStrategy interface {
	createPolicy(ctx context.Context, cloudletName, cloudletCode string, groupID int64) (int64, error)
	updatePolicyVersion(ctx context.Context, d *schema.ResourceData, policyID, version int64, description, matchRulesJSON string, newVersionRequired bool) (error, error)
	updatePolicy(ctx context.Context, policyID, groupID int64, cloudletName string) error
	newPolicyVersionIsNeeded(ctx context.Context, policyID, version int64) (bool, error)
	readPolicy(ctx context.Context, policyID int64, version *int64) (map[string]any, error)
	deletePolicy(ctx context.Context, policyID int64) error
	getVersionStrategy(meta meta.Meta) versionStrategy
	setPolicyType(d *schema.ResourceData) error
	isFirstVersionCreated() bool
}

func discoverPolicyExecutionStrategy(ctx context.Context, meta meta.Meta, policyName string) (policyExecutionStrategy, int64, error) {

	strategy, policyID, errV2 := checkForV2Policy(ctx, meta, policyName)
	if strategy != nil {
		return strategy, policyID, nil
	}

	strategy, policyID, errV3 := checkForV3Policy(ctx, meta, policyName)
	if strategy != nil {
		return strategy, policyID, nil
	}

	var errMessage string
	if errV2 != nil {
		errMessage += fmt.Sprintf("could not list V2 policies: %s\n", errV2)
	}
	if errV3 != nil {
		errMessage += fmt.Sprintf("could not list V3 policies: %s", errV3)
	}
	if errMessage != "" {
		return nil, 0, fmt.Errorf(errMessage)
	}

	return nil, 0, fmt.Errorf("policy '%s' does not exist", policyName)
}

func checkForV2Policy(ctx context.Context, meta meta.Meta, policyName string) (policyExecutionStrategy, int64, error) {
	v2Client := Client(meta)
	size, offset := 1000, 0
	var errV2 error
	for {
		policies, err := v2Client.ListPolicies(ctx, cloudlets.ListPoliciesRequest{
			Offset:   offset,
			PageSize: tools.IntPtr(size),
		})
		if err == nil {
			if policyID := findPolicyV2ByName(policies, policyName); policyID != 0 {
				return v2PolicyStrategy{
					client: v2Client,
				}, policyID, nil
			}
			if len(policies) < size {
				break
			}
			offset++
		} else {
			errV2 = err
			break
		}
	}

	return nil, 0, errV2
}

func checkForV3Policy(ctx context.Context, meta meta.Meta, policyName string) (policyExecutionStrategy, int64, error) {
	v3Client := ClientV3(meta)
	size, page := 1000, 0
	var errV3 error
	for {
		policiesV3, err := v3Client.ListPolicies(ctx, v3.ListPoliciesRequest{
			Page: page,
			Size: size,
		})
		if err == nil {
			if policyID := findPolicyV3ByName(policiesV3.Content, policyName); policyID != 0 {
				return v3PolicyStrategy{
					client: v3Client,
				}, policyID, nil
			}
			if len(policiesV3.Content) < size {
				break
			}
			page++
		} else {
			errV3 = err
			break
		}
	}

	return nil, 0, errV3
}

func findPolicyV3ByName(policies []v3.Policy, policyName string) int64 {
	for _, policy := range policies {
		if policy.Name == policyName {
			return policy.ID
		}
	}
	return 0
}

func findPolicyV2ByName(policies []cloudlets.Policy, policyName string) int64 {
	for _, policy := range policies {
		if policy.Name == policyName {
			return policy.PolicyID
		}
	}
	return 0
}
