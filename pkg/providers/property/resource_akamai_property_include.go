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

	"github.com/akamai/AkamaiOPEN-edgegrid-golang/v3/pkg/papi"
	"github.com/akamai/AkamaiOPEN-edgegrid-golang/v3/pkg/session"
	"github.com/akamai/terraform-provider-akamai/v3/pkg/akamai"
	"github.com/akamai/terraform-provider-akamai/v3/pkg/tools"
	"github.com/hashicorp/terraform-plugin-sdk/v2/diag"
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
				Type:        schema.TypeString,
				Required:    true,
				ForceNew:    true,
				Description: "A descriptive name for the include",
			},
			"rule_format": {
				Type:             schema.TypeString,
				Required:         true,
				Description:      "Indicates the versioned set of features and criteria",
				ValidateDiagFunc: tools.ValidateRuleFormat,
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
				DiffSuppressFunc: tools.ComposeDiffSuppress(suppressDefaultRules, diffSuppressRules),
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
		},
	}
}

func resourcePropertyIncludeCreate(ctx context.Context, rd *schema.ResourceData, m interface{}) diag.Diagnostics {
	meta := akamai.Meta(m)
	logger := meta.Log("PAPI", "resourcePropertyIncludeCreate")
	ctx = session.ContextWithOptions(ctx, session.WithContextLog(logger))
	client := inst.Client(meta)

	logger.Debug("Creating property include")

	contractID, err := tools.GetStringValue("contract_id", rd)
	if err != nil {
		return diag.FromErr(err)
	}

	groupID, err := tools.GetStringValue("group_id", rd)
	if err != nil {
		return diag.FromErr(err)
	}

	productID, err := tools.GetStringValue("product_id", rd)
	if err != nil {
		if errors.Is(err, tools.ErrNotFound) {
			return diag.Errorf(`The argument "product_id" is required during create, but no definition was found`)
		}
		return diag.FromErr(err)
	}

	name, err := tools.GetStringValue("name", rd)
	if err != nil {
		return diag.FromErr(err)
	}

	ruleFormat, err := tools.GetStringValue("rule_format", rd)
	if err != nil {
		return diag.FromErr(err)
	}

	includeType, err := tools.GetStringValue("type", rd)
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
	meta := akamai.Meta(m)
	logger := meta.Log("PAPI", "resourcePropertyIncludeRead")
	ctx = session.ContextWithOptions(ctx, session.WithContextLog(logger))
	client := inst.Client(meta)

	logger.Debug("Reading property include")

	includeID := rd.Id()

	contractID, err := tools.GetStringValue("contract_id", rd)
	if err != nil {
		return diag.FromErr(err)
	}

	groupID, err := tools.GetStringValue("group_id", rd)
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
		"rules":              string(rulesJSON),
		"name":               include.IncludeName,
		"type":               include.IncludeType,
		"rule_format":        getIncludeRuleTreeResp.RuleFormat,
		"latest_version":     include.LatestVersion,
		"staging_version":    stagingVersion,
		"production_version": productionVersion,
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

	if err = tools.SetAttrs(rd, attrs); err != nil {
		return diag.Errorf("%v: %s", tools.ErrValueSet, err.Error())
	}

	return nil
}

func resourcePropertyIncludeUpdate(ctx context.Context, rd *schema.ResourceData, m interface{}) diag.Diagnostics {
	meta := akamai.Meta(m)
	logger := meta.Log("PAPI", "resourcePropertyIncludeUpdate")
	ctx = session.ContextWithOptions(ctx, session.WithContextLog(logger))
	client := inst.Client(meta)

	logger.Debug("Updating property include")

	includeID := rd.Id()

	contractID, err := tools.GetStringValue("contract_id", rd)
	if err != nil {
		return diag.FromErr(err)
	}

	groupID, err := tools.GetStringValue("group_id", rd)
	if err != nil {
		return diag.FromErr(err)
	}

	latestVersion, err := tools.GetIntValue("latest_version", rd)
	if err != nil {
		return diag.FromErr(err)
	}

	includeVersion, err := client.GetIncludeVersion(ctx, papi.GetIncludeVersionRequest{
		Version:    latestVersion,
		GroupID:    groupID,
		IncludeID:  includeID,
		ContractID: contractID,
	})

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
	meta := akamai.Meta(m)
	logger := meta.Log("PAPI", "resourcePropertyIncludeDelete")
	ctx = session.ContextWithOptions(ctx, session.WithContextLog(logger))
	client := inst.Client(meta)

	logger.Debug("Deleting property include")

	includeID := rd.Id()

	contractID, err := tools.GetStringValue("contract_id", rd)
	if err != nil {
		return diag.FromErr(err)
	}

	groupID, err := tools.GetStringValue("group_id", rd)
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
	meta := akamai.Meta(m)
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
		return nil, fmt.Errorf("%v: %s", tools.ErrValueSet, err.Error())
	}

	if err := rd.Set("group_id", groupID); err != nil {
		return nil, fmt.Errorf("%v: %s", tools.ErrValueSet, err.Error())
	}

	rd.SetId(includeID)

	return []*schema.ResourceData{rd}, nil
}

func updateRules(ctx context.Context, client papi.PAPI, rd *schema.ResourceData, version int) error {
	rulesJSON, err := tools.GetStringValue("rules", rd)
	if err != nil && !errors.Is(err, tools.ErrNotFound) {
		return err
	}

	if len(rulesJSON) == 0 {
		return nil
	}

	includeID := rd.Id()
	contractID, err := tools.GetStringValue("contract_id", rd)
	if err != nil {
		return err
	}

	groupID, err := tools.GetStringValue("group_id", rd)
	if err != nil {
		return err
	}

	ruleFormat, err := tools.GetStringValue("rule_format", rd)
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
	logger := akamai.Log("PAPI", "suppressDefaultRules")
	if len(newValue) > 0 {
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
