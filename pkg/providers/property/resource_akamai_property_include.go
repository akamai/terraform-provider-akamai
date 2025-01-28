package property

import (
	"context"
	"encoding/json"
	"errors"
	"fmt"
	"net/http"
	"reflect"
	"strconv"
	"strings"

	"github.com/akamai/AkamaiOPEN-edgegrid-golang/v10/pkg/papi"
	"github.com/akamai/AkamaiOPEN-edgegrid-golang/v10/pkg/session"
	"github.com/akamai/terraform-provider-akamai/v6/pkg/common/str"
	"github.com/akamai/terraform-provider-akamai/v6/pkg/common/tf"
	"github.com/akamai/terraform-provider-akamai/v6/pkg/log"
	"github.com/akamai/terraform-provider-akamai/v6/pkg/meta"
	"github.com/hashicorp/terraform-plugin-sdk/v2/diag"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/customdiff"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/schema"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/validation"
)

func resourcePropertyInclude() *schema.Resource {
	return &schema.Resource{
		CreateContext: resourcePropertyIncludeCreate,
		ReadContext:   resourcePropertyIncludeRead,
		UpdateContext: resourcePropertyIncludeUpdate,
		DeleteContext: resourcePropertyIncludeDelete,
		Importer: &schema.ResourceImporter{
			StateContext: resourcePropertyIncludeImport,
		},
		CustomizeDiff: customdiff.All(
			propertyIncludeRulesCustomDiff,
			setIncludeVersionsComputedOnRulesChange,
		),
		Schema: map[string]*schema.Schema{
			"contract_id": {
				Type:        schema.TypeString,
				Required:    true,
				ForceNew:    true,
				Description: "Identifies the contract to which the include is assigned",
				StateFunc:   addPrefixToState("ctr_"),
			},
			"group_id": {
				Type:        schema.TypeString,
				Required:    true,
				ForceNew:    true,
				Description: "Identifies the group to which the include is assigned",
				StateFunc:   addPrefixToState("grp_"),
			},
			"product_id": {
				Type:        schema.TypeString,
				Optional:    true,
				ForceNew:    true,
				Description: "The product assigned to the include",
				StateFunc:   addPrefixToState("prd_"),
			},
			"name": {
				Type:             schema.TypeString,
				Required:         true,
				ForceNew:         true,
				ValidateDiagFunc: validateNameWithBound(3),
				Description:      "A descriptive name for the include",
			},
			"rule_format": {
				Type:             schema.TypeString,
				Required:         true,
				Description:      "Indicates the versioned set of features and criteria",
				ValidateDiagFunc: tf.AggregateValidations(tf.IsNotBlank, tf.ValidateRuleFormat),
			},
			"type": {
				Type:         schema.TypeString,
				Required:     true,
				ForceNew:     true,
				Description:  "Specifies the type of the include, either 'MICROSERVICES' or 'COMMON_SETTINGS'",
				ValidateFunc: validation.StringInSlice([]string{string(papi.IncludeTypeMicroServices), string(papi.IncludeTypeCommonSettings)}, false),
			},
			"rules": {
				Type:             schema.TypeString,
				Optional:         true,
				Computed:         true,
				Description:      "Property Rules as JSON",
				ValidateDiagFunc: validation.ToDiagFunc(validation.StringIsJSON),
				DiffSuppressFunc: tf.DiffSuppressAny(suppressDefaultRules, diffSuppressPropertyRules),
				StateFunc:        rulesStateFunc,
			},
			"rule_errors": {
				Type:        schema.TypeString,
				Computed:    true,
				Description: "Rule validation errors",
			},
			"rule_warnings": {
				Type:        schema.TypeString,
				Computed:    true,
				Description: "Rule validation warnings",
			},
			"latest_version": {
				Type:        schema.TypeInt,
				Computed:    true,
				Description: "Specifies the most recent version of the include",
			},
			"staging_version": {
				Type:        schema.TypeString,
				Computed:    true,
				Description: "The most recent version to be activated to the staging network",
			},
			"production_version": {
				Type:        schema.TypeString,
				Computed:    true,
				Description: "The most recent version to be activated to the production network",
			},
			"asset_id": {
				Type:        schema.TypeString,
				Computed:    true,
				Description: "ID of the include in the Identity and Access Management API.",
			},
		},
	}
}

// propertyIncludeRulesCustomDiff compares Rules.Criteria and Rules.Children fields from terraform state
// and from a new configuration. If some of these fields are empty lists in the new configuration and
// are nil in the terraform state, then this function returns no difference for these fields.
//
// TODO: reuse propertyRulesCustomDiff when version_notes attr is added to akamai_property_include resource.
func propertyIncludeRulesCustomDiff(_ context.Context, diff *schema.ResourceDiff, _ interface{}) error {
	o, n := diff.GetChange("rules")
	oldValue, newValue := o.(string), n.(string)

	handleCreate := diff.Id() == "" && newValue != ""
	if !handleCreate && (oldValue == "" || newValue == "") {
		return nil
	}

	var newRulesUpdate papi.RulesUpdate
	if err := json.Unmarshal([]byte(newValue), &newRulesUpdate); err != nil {
		return fmt.Errorf("cannot parse rules JSON from config: %s", err)
	}

	if handleCreate {
		rules, err := unifyRulesDiff(newRulesUpdate)
		if err != nil {
			return err
		}
		if err = diff.SetNew("rules", rules); err != nil {
			return fmt.Errorf("cannot set a new diff value for 'rules' %s", err)
		}
		return nil
	}

	var oldRulesUpdate papi.RulesUpdate
	if err := json.Unmarshal([]byte(oldValue), &oldRulesUpdate); err != nil {
		return fmt.Errorf("cannot parse rules JSON from state: %s", err)
	}

	normalizeFields(&oldRulesUpdate, &newRulesUpdate)
	if rulesEqual(&oldRulesUpdate.Rules, &newRulesUpdate.Rules) && oldRulesUpdate.Comments == newRulesUpdate.Comments {
		return nil
	}

	rules, err := json.Marshal(newRulesUpdate)
	if err != nil {
		return fmt.Errorf("cannot encode rules JSON %s", err)
	}

	if err = diff.SetNew("rules", string(rules)); err != nil {
		return fmt.Errorf("cannot set a new diff value for 'rules' %s", err)
	}
	return nil
}

func resourcePropertyIncludeCreate(ctx context.Context, rd *schema.ResourceData, m interface{}) diag.Diagnostics {
	meta := meta.Must(m)
	logger := meta.Log("PAPI", "resourcePropertyIncludeCreate")
	ctx = session.ContextWithOptions(ctx, session.WithContextLog(logger))
	client := Client(meta)

	logger.Debug("Creating property include")

	contractID, err := tf.GetStringValue("contract_id", rd)
	if err != nil {
		return diag.FromErr(err)
	}

	groupID, err := tf.GetStringValue("group_id", rd)
	if err != nil {
		return diag.FromErr(err)
	}

	productID, err := tf.GetStringValue("product_id", rd)
	if err != nil {
		if errors.Is(err, tf.ErrNotFound) {
			return diag.Errorf(`The argument "product_id" is required during create, but no definition was found`)
		}
		return diag.FromErr(err)
	}

	name, err := tf.GetStringValue("name", rd)
	if err != nil {
		return diag.FromErr(err)
	}

	ruleFormat, err := tf.GetStringValue("rule_format", rd)
	if err != nil {
		return diag.FromErr(err)
	}

	includeType, err := tf.GetStringValue("type", rd)
	if err != nil {
		return diag.FromErr(err)
	}

	createIncludeResp, err := client.CreateInclude(ctx, papi.CreateIncludeRequest{
		ContractID:  contractID,
		GroupID:     groupID,
		ProductID:   productID,
		IncludeName: name,
		IncludeType: papi.IncludeType(includeType),
		RuleFormat:  ruleFormat,
	})
	if err != nil {
		return diag.Errorf("%s create: %s", ErrPropertyInclude, err)
	}

	rd.SetId(createIncludeResp.IncludeID)

	postCreateVersion := 1
	if err = updateRules(ctx, client, rd, postCreateVersion); err != nil {
		return diag.Errorf("%s update: %s", ErrPropertyInclude, err)
	}

	return resourcePropertyIncludeRead(ctx, rd, m)
}

func resourcePropertyIncludeRead(ctx context.Context, rd *schema.ResourceData, m interface{}) diag.Diagnostics {
	meta := meta.Must(m)
	logger := meta.Log("PAPI", "resourcePropertyIncludeRead")
	ctx = session.ContextWithOptions(ctx, session.WithContextLog(logger))
	client := Client(meta)

	logger.Debug("Reading property include")

	includeID := rd.Id()

	contractID, err := tf.GetStringValue("contract_id", rd)
	if err != nil {
		return diag.FromErr(err)
	}

	groupID, err := tf.GetStringValue("group_id", rd)
	if err != nil {
		return diag.FromErr(err)
	}

	getIncludeResp, err := client.GetInclude(ctx, papi.GetIncludeRequest{
		GroupID:    groupID,
		IncludeID:  includeID,
		ContractID: contractID,
	})
	if err != nil {
		return diag.Errorf("%s read: %s", ErrPropertyInclude, err)
	}

	include := getIncludeResp.Include

	getIncludeRuleTreeResp, err := client.GetIncludeRuleTree(ctx, papi.GetIncludeRuleTreeRequest{
		GroupID:        groupID,
		IncludeID:      includeID,
		ContractID:     contractID,
		ValidateRules:  true,
		IncludeVersion: include.LatestVersion,
	})
	if err != nil {
		return diag.Errorf("%s read: %s", ErrPropertyInclude, err)
	}

	getIncludeVersionResp, err := client.GetIncludeVersion(ctx, papi.GetIncludeVersionRequest{
		IncludeID:  includeID,
		Version:    include.LatestVersion,
		ContractID: contractID,
		GroupID:    groupID,
	})
	if err != nil {
		return diag.Errorf("%s read: %s", ErrPropertyInclude, err)
	}
	productID := str.AddPrefix(getIncludeVersionResp.IncludeVersion.ProductID, "prd_")
	logger.Debugf("Fetched product id: %s from version: %d, will be saved in state as: %s",
		getIncludeVersionResp.IncludeVersion.ProductID, include.LatestVersion, productID)

	rules := papi.RulesUpdate{
		Comments: getIncludeRuleTreeResp.Comments,
		Rules:    getIncludeRuleTreeResp.Rules,
	}
	rulesJSON, err := json.MarshalIndent(rules, "", "  ")
	if err != nil {
		return diag.Errorf("%s read: %s", ErrPropertyInclude, err)
	}

	var stagingVersion, productionVersion string
	if include.StagingVersion != nil {
		stagingVersion = strconv.Itoa(*include.StagingVersion)
	}

	if include.ProductionVersion != nil {
		productionVersion = strconv.Itoa(*include.ProductionVersion)
	}

	attrs := map[string]interface{}{
		"asset_id":           include.AssetID,
		"rules":              string(rulesJSON),
		"name":               include.IncludeName,
		"type":               include.IncludeType,
		"rule_format":        getIncludeRuleTreeResp.RuleFormat,
		"latest_version":     include.LatestVersion,
		"staging_version":    stagingVersion,
		"production_version": productionVersion,
		"product_id":         productID,
	}

	var rulesError string
	if len(getIncludeRuleTreeResp.Errors) > 0 {
		rulesErrorsJSON, err := json.MarshalIndent(getIncludeRuleTreeResp.Errors, "", "  ")
		if err != nil {
			return diag.Errorf("%s read: %s", ErrPropertyInclude, err)
		}

		rulesError = string(rulesErrorsJSON)
		logger.Errorf("Include has rule errors: %s", rulesErrorsJSON)
	}
	attrs["rule_errors"] = rulesError

	var rulesWarnings string
	if len(getIncludeRuleTreeResp.Warnings) > 0 {
		rulesWarningsJSON, err := json.MarshalIndent(getIncludeRuleTreeResp.Warnings, "", "  ")
		if err != nil {
			return diag.Errorf("%s read: %s", ErrPropertyInclude, err)
		}

		rulesWarnings = string(rulesWarningsJSON)
		logger.Errorf("Include has rule warnings: %s", rulesWarningsJSON)
	}
	attrs["rule_warnings"] = rulesWarnings

	if err = tf.SetAttrs(rd, attrs); err != nil {
		return diag.Errorf("%v: %s", tf.ErrValueSet, err.Error())
	}

	return nil
}

func resourcePropertyIncludeUpdate(ctx context.Context, rd *schema.ResourceData, m interface{}) diag.Diagnostics {
	meta := meta.Must(m)
	logger := meta.Log("PAPI", "resourcePropertyIncludeUpdate")
	ctx = session.ContextWithOptions(ctx, session.WithContextLog(logger))
	client := Client(meta)

	logger.Debug("Updating property include")

	includeID := rd.Id()

	contractID, err := tf.GetStringValue("contract_id", rd)
	if err != nil {
		return diag.FromErr(err)
	}

	groupID, err := tf.GetStringValue("group_id", rd)
	if err != nil {
		return diag.FromErr(err)
	}

	latestVersion, err := tf.GetIntValue("latest_version", rd)
	if err != nil {
		return diag.FromErr(err)
	}

	includeVersion, err := client.GetIncludeVersion(ctx, papi.GetIncludeVersionRequest{
		Version:    latestVersion,
		GroupID:    groupID,
		IncludeID:  includeID,
		ContractID: contractID,
	})
	if err != nil {
		return diag.FromErr(err)
	}

	version := latestVersion
	if !isVersionEditable(includeVersion.IncludeVersion) {
		createVersionResp, err := client.CreateIncludeVersion(ctx, papi.CreateIncludeVersionRequest{
			IncludeID: includeID,
			IncludeVersionRequest: papi.IncludeVersionRequest{
				CreateFromVersion: version,
			},
		})
		if err != nil {
			return diag.Errorf("%s update: %s", ErrPropertyInclude, err)
		}
		version = createVersionResp.Version
	}

	if err = updateRules(ctx, client, rd, version); err != nil {
		return diag.Errorf("%s update: %s", ErrPropertyInclude, err)
	}

	return resourcePropertyIncludeRead(ctx, rd, m)
}

func resourcePropertyIncludeDelete(ctx context.Context, rd *schema.ResourceData, m interface{}) diag.Diagnostics {
	meta := meta.Must(m)
	logger := meta.Log("PAPI", "resourcePropertyIncludeDelete")
	ctx = session.ContextWithOptions(ctx, session.WithContextLog(logger))
	client := Client(meta)

	logger.Debug("Deleting property include")

	includeID := rd.Id()

	contractID, err := tf.GetStringValue("contract_id", rd)
	if err != nil {
		return diag.FromErr(err)
	}

	groupID, err := tf.GetStringValue("group_id", rd)
	if err != nil {
		return diag.FromErr(err)
	}

	getIncludeResp, err := client.GetInclude(ctx, papi.GetIncludeRequest{
		GroupID:    groupID,
		IncludeID:  includeID,
		ContractID: contractID,
	})
	if err != nil {
		return diag.Errorf("%s delete: %s", ErrPropertyInclude, err)
	}

	if err := canIncludeBeDeleted(getIncludeResp.Include); err != nil {
		return append(diag.Errorf("Include '%s' could not be deleted due to the following reason(s):", includeID), err...)
	}

	_, err = client.DeleteInclude(ctx, papi.DeleteIncludeRequest{
		ContractID: contractID,
		IncludeID:  includeID,
		GroupID:    groupID,
	})
	if err != nil {
		return diag.Errorf("%s delete: %s", ErrPropertyInclude, err)
	}

	return nil
}

func resourcePropertyIncludeImport(_ context.Context, rd *schema.ResourceData, m interface{}) ([]*schema.ResourceData, error) {
	meta := meta.Must(m)
	logger := meta.Log("PAPI", "resourcePropertyIncludeImport")

	logger.Debug("Importing property include")

	parts := strings.Split(rd.Id(), ":")
	if len(parts) != 3 {
		return nil, fmt.Errorf("%s import: invalid import id '%s'"+
			"- colon separated list of contract ID, group ID and property include ID has to be supplied",
			ErrPropertyInclude, rd.Id())
	}

	contractID, groupID, includeID := parts[0], parts[1], parts[2]

	if err := rd.Set("contract_id", contractID); err != nil {
		return nil, fmt.Errorf("%v: %s", tf.ErrValueSet, err.Error())
	}

	if err := rd.Set("group_id", groupID); err != nil {
		return nil, fmt.Errorf("%v: %s", tf.ErrValueSet, err.Error())
	}

	rd.SetId(includeID)

	return []*schema.ResourceData{rd}, nil
}

func updateRules(ctx context.Context, client papi.PAPI, rd *schema.ResourceData, version int) error {
	rulesJSON, err := tf.GetStringValue("rules", rd)
	if err != nil && !errors.Is(err, tf.ErrNotFound) {
		return err
	}

	if len(rulesJSON) == 0 {
		return nil
	}

	includeID := rd.Id()
	contractID, err := tf.GetStringValue("contract_id", rd)
	if err != nil {
		return err
	}

	groupID, err := tf.GetStringValue("group_id", rd)
	if err != nil {
		return err
	}

	ruleFormat, err := tf.GetStringValue("rule_format", rd)
	if err != nil {
		return err
	}

	var rules papi.RulesUpdate
	if err = json.Unmarshal([]byte(rulesJSON), &rules); err != nil {
		return fmt.Errorf("unmarshalling rules failed: %s", err)
	}

	header := buildRuleFormatHeader(ruleFormat)
	ctx = session.ContextWithOptions(ctx, session.WithContextHeaders(header))

	_, err = client.UpdateIncludeRuleTree(ctx, papi.UpdateIncludeRuleTreeRequest{
		ContractID:     contractID,
		GroupID:        groupID,
		IncludeID:      includeID,
		IncludeVersion: version,
		Rules:          rules,
	})
	if err != nil {
		return err
	}

	return nil
}

// canIncludeBeDeleted returns error if there is any version active on
// either staging or production network as it prevents the deletion
func canIncludeBeDeleted(include papi.Include) diag.Diagnostics {
	var diags diag.Diagnostics
	if include.StagingVersion != nil {
		diags = append(diags, diag.Errorf("version '%d' is active on 'STAGING' network", *include.StagingVersion)...)
	}

	if include.ProductionVersion != nil {
		diags = append(diags, diag.Errorf("version '%d' is active on 'PRODUCTION' network", *include.ProductionVersion)...)
	}

	return diags
}

func isVersionEditable(includeVersion papi.IncludeVersion) bool {
	return includeVersion.StagingStatus == papi.VersionStatusInactive &&
		includeVersion.ProductionStatus == papi.VersionStatusInactive
}

func buildRuleFormatHeader(ruleFormat string) http.Header {
	MIME := fmt.Sprintf("application/vnd.akamai.papirules.%s+json", ruleFormat)
	return http.Header{"Content-Type": []string{MIME}}
}

func suppressDefaultRules(_, oldValue, newValue string, _ *schema.ResourceData) bool {
	logger := log.Get("PAPI", "suppressDefaultRules")
	if len(newValue) > 0 || len(oldValue) == 0 {
		return false
	}

	var rules papi.Rules
	if err := json.Unmarshal([]byte(oldValue), &rules); err != nil {
		logger.Errorf("Unable to unmarshal 'old' JSON rules: %s", err)
		return false
	}

	defaultRules := papi.Rules{Name: "default"}

	return reflect.DeepEqual(rules, defaultRules)
}

// setIncludeVersionsComputedOnRulesChange is a schema.CustomizeDiffFunc for akamai_property_include resource,
// which sets latest_version fields as computed if a new version of the include is expected to be created.
func setIncludeVersionsComputedOnRulesChange(_ context.Context, rd *schema.ResourceDiff, _ interface{}) error {
	ruleFormatChanged := rd.HasChange("rule_format")

	oldRules, newRules := rd.GetChange("rules")
	rulesEqual, err := rulesJSONEqual(oldRules.(string), newRules.(string))
	if err != nil {
		return err
	}

	if !ruleFormatChanged && rulesEqual {
		return nil
	}

	if err := rd.SetNewComputed("latest_version"); err != nil {
		return fmt.Errorf("%w: %s", tf.ErrValueSet, err.Error())
	}

	return nil
}
