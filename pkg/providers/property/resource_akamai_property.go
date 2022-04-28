package property

import (
	"context"
	"encoding/json"
	"errors"
	"fmt"
	"net/http"
	"regexp"
	"strconv"
	"strings"

	"github.com/apex/log"
	"github.com/hashicorp/go-cty/cty"
	"github.com/hashicorp/terraform-plugin-sdk/v2/diag"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/customdiff"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/schema"

	"github.com/akamai/AkamaiOPEN-edgegrid-golang/v2/pkg/papi"
	"github.com/akamai/AkamaiOPEN-edgegrid-golang/v2/pkg/session"
	"github.com/akamai/terraform-provider-akamai/v2/pkg/akamai"
	"github.com/akamai/terraform-provider-akamai/v2/pkg/tools"
)

func resourceProperty() *schema.Resource {
	papiError := func() *schema.Resource {
		return &schema.Resource{Schema: map[string]*schema.Schema{
			"type":           {Type: schema.TypeString, Optional: true},
			"title":          {Type: schema.TypeString, Optional: true},
			"detail":         {Type: schema.TypeString, Optional: true},
			"instance":       {Type: schema.TypeString, Optional: true},
			"behavior_name":  {Type: schema.TypeString, Optional: true},
			"error_location": {Type: schema.TypeString, Optional: true},
			"status_code":    {Type: schema.TypeInt, Optional: true},
		}}
	}

	hashHostname := func(v interface{}) int {
		m, ok := v.(map[string]interface{})
		if !ok {
			return 0
		}
		cnameFrom, ok := m["cname_from"]
		if !ok {
			return 0
		}
		cnameTo, ok := m["cname_to"]
		if !ok {
			return 0
		}
		certProvisioningType, ok := m["cert_provisioning_type"]
		if !ok {
			return 0
		}
		return schema.HashString(fmt.Sprintf("%s.%s.%s", cnameFrom, cnameTo, certProvisioningType))
	}

	validateRules := func(val interface{}, _ cty.Path) diag.Diagnostics {
		if len(val.(string)) == 0 {
			return nil
		}

		var target map[string]interface{}
		if err := json.Unmarshal([]byte(val.(string)), &target); err != nil {
			return diag.Errorf("rules are not valid JSON")
		}
		return nil
	}

	diffSuppressRules := func(_, old, new string, _ *schema.ResourceData) bool {
		logger := akamai.Log("PAPI", "suppressRulesJSON")

		if old == "" || new == "" {
			return old == new
		}

		var oldRules, newRules papi.RulesUpdate
		if err := json.Unmarshal([]byte(old), &oldRules); err != nil {
			logger.Errorf("Unable to unmarshal 'old' JSON rules: %s", err)
			return false
		}

		if err := json.Unmarshal([]byte(new), &newRules); err != nil {
			logger.Errorf("Unable to unmarshal 'new' JSON rules: %s", err)
			return false
		}

		return compareRuleTree(&oldRules, &newRules)
	}

	return &schema.Resource{
		CreateContext: resourcePropertyCreate,
		ReadContext:   resourcePropertyRead,
		UpdateContext: resourcePropertyUpdate,
		DeleteContext: resourcePropertyDelete,
		CustomizeDiff: customdiff.All(
			rulesCustomDiff,
			hostNamesCustomDiff,
			versionsComputedValuesCustomDiff,
		),
		Importer: &schema.ResourceImporter{
			StateContext: resourcePropertyImport,
		},
		StateUpgraders: []schema.StateUpgrader{{
			Version: 0,
			Type:    resourcePropertyV0().CoreConfigSchema().ImpliedType(),
			Upgrade: upgradePropV0,
		}},
		SchemaVersion: 1,
		Schema: map[string]*schema.Schema{
			// Required
			"name": {
				Type:             schema.TypeString,
				Required:         true,
				ForceNew:         true,
				ValidateDiagFunc: validatePropertyName,
				Description:      "Name to give to the Property (must be unique)",
			},

			"group_id": {
				Type:         schema.TypeString,
				Optional:     true,
				Computed:     true,
				ExactlyOneOf: []string{"group_id", "group"},
				StateFunc:    addPrefixToState("grp_"),
				Description:  "Group ID to be assigned to the Property",
			},
			"group": {
				Type:       schema.TypeString,
				Optional:   true,
				Computed:   true,
				Deprecated: akamai.NoticeDeprecatedUseAlias("group"),
				StateFunc:  addPrefixToState("grp_"),
			},

			"contract_id": {
				Type:         schema.TypeString,
				Optional:     true,
				Computed:     true,
				ExactlyOneOf: []string{"contract_id", "contract"},
				StateFunc:    addPrefixToState("ctr_"),
				Description:  "Contract ID to be assigned to the Property",
			},
			"contract": {
				Type:       schema.TypeString,
				Optional:   true,
				Computed:   true,
				Deprecated: akamai.NoticeDeprecatedUseAlias("contract"),
				StateFunc:  addPrefixToState("ctr_"),
			},

			"product_id": {
				Type:        schema.TypeString,
				Optional:    true,
				Computed:    true,
				Description: "Product ID to be assigned to the Property",
				StateFunc:   addPrefixToState("prd_"),
			},
			"product": {
				Type:         schema.TypeString,
				Optional:     true,
				Computed:     true,
				ExactlyOneOf: []string{"product_id"},
				Deprecated:   akamai.NoticeDeprecatedUseAlias("product"),
				StateFunc:    addPrefixToState("prd_"),
			},

			// Optional
			"rule_format": {
				Type:        schema.TypeString,
				Optional:    true,
				Computed:    true,
				Description: "Specify the rule format version (defaults to latest version available when created)",
				ValidateDiagFunc: func(v interface{}, _ cty.Path) diag.Diagnostics {
					format := v.(string)
					if format == "" || format == "latest" {
						return nil
					}

					if !regexp.MustCompile(`^v[0-9]{4}-[0-9]{2}-[0-9]{2}$`).MatchString(format) {
						url := "https://developer.akamai.com/api/core_features/property_manager/vlatest.html#behaviors"
						return diag.Errorf(`"rule_format" must be of the form vYYYY-MM-DD (with a leading "v") see %s`, url)
					}

					return nil
				},
			},
			"rules": {
				Type:             schema.TypeString,
				Optional:         true,
				Computed:         true,
				Description:      "Property Rules as JSON",
				ValidateDiagFunc: validateRules,
				DiffSuppressFunc: diffSuppressRules,
				StateFunc: func(v interface{}) string {
					var js string
					if json.Unmarshal([]byte(v.(string)), &js) == nil {
						return compactJSON([]byte(v.(string)))
					}
					return v.(string)
				},
			},
			"hostnames": {
				Type:     schema.TypeSet,
				Optional: true,
				Set:      hashHostname,
				Elem: &schema.Resource{
					Schema: map[string]*schema.Schema{
						"cname_from": {
							Type:     schema.TypeString,
							Required: true,
							ValidateDiagFunc: func(i interface{}, path cty.Path) diag.Diagnostics {
								if len(i.(string)) == 0 {
									return diag.Errorf("'cname_from' cannot be empty when hostnames block is defined - See new hostnames schema")
								}
								return nil
							},
						},
						"cname_to": {
							Type:     schema.TypeString,
							Required: true,
							ValidateDiagFunc: func(i interface{}, path cty.Path) diag.Diagnostics {
								if len(i.(string)) == 0 {
									return diag.Errorf("'cname_to' cannot be empty when hostnames block is defined - See new hostnames schema")
								}
								return nil
							},
						},
						"cert_provisioning_type": {
							Type:     schema.TypeString,
							Required: true,
							ValidateDiagFunc: func(i interface{}, path cty.Path) diag.Diagnostics {
								if len(i.(string)) == 0 {
									return diag.Errorf("'cert_provisioning_type' cannot be empty when hostnames block is defined - See new hostnames schema")
								}
								return nil
							},
						},
						"cname_type": {
							Type:     schema.TypeString,
							Optional: true,
							Computed: true,
						},
						"edge_hostname_id": {
							Type:     schema.TypeString,
							Computed: true,
						},
						"cert_status": {
							Type:     schema.TypeList,
							Optional: true,
							Computed: true,
							Elem:     certStatus,
						},
					},
				},
			},

			// Computed
			"latest_version": {
				Type:        schema.TypeInt,
				Computed:    true,
				Description: "Property's current latest version number",
			},
			"staging_version": {
				Type:        schema.TypeInt,
				Computed:    true,
				Description: "Property's version currently activated in staging (zero when not active in staging)",
			},
			"production_version": {
				Type:        schema.TypeInt,
				Computed:    true,
				Description: "Property's version currently activated in production (zero when not active in production)",
			},
			"read_version": {
				Type:        schema.TypeInt,
				Computed:    true,
				Description: "Required property's version to be read",
			},
			"rule_errors": {
				Type:     schema.TypeList,
				Computed: true,
				Elem:     papiError(),
			},
			"rule_warnings": {
				Type:       schema.TypeList,
				Optional:   true,
				Computed:   true,
				Elem:       papiError(),
				Deprecated: "Rule warnings will not be set in state anymore",
			},

			// Hard-deprecated attributes: These are effectively removed, but we wanted to refer users to the upgrade guide
			"cp_code": {
				Type:       schema.TypeString,
				Optional:   true,
				Deprecated: akamai.NoticeDeprecatedUseAlias("cp_code"),
			},
			"contact": {
				Type:       schema.TypeSet,
				Optional:   true,
				Elem:       &schema.Schema{Type: schema.TypeString},
				Deprecated: akamai.NoticeDeprecatedUseAlias("contact"),
			},
			"origin": {
				Type:     schema.TypeSet,
				Optional: true,
				Elem: &schema.Resource{
					Schema: map[string]*schema.Schema{
						"hostname":              {Type: schema.TypeString, Optional: true},
						"port":                  {Type: schema.TypeInt, Optional: true},
						"forward_hostname":      {Type: schema.TypeString, Optional: true},
						"cache_key_hostname":    {Type: schema.TypeString, Optional: true},
						"compress":              {Type: schema.TypeBool, Optional: true},
						"enable_true_client_ip": {Type: schema.TypeBool, Optional: true},
					},
				},
				Deprecated: akamai.NoticeDeprecatedUseAlias("origin"),
			},
			"is_secure": {
				Type:       schema.TypeBool,
				Optional:   true,
				Deprecated: akamai.NoticeDeprecatedUseAlias("is_secure"),
			},
			"variables": {
				Type:       schema.TypeString,
				Optional:   true,
				Deprecated: akamai.NoticeDeprecatedUseAlias("variables"),
			},
		},
	}
}

// isValidPropertyName is a function that validates if given string contains only letters, numbers, and these characters: . _ -
var isValidPropertyName = regexp.MustCompile(`^[A-Za-z0-9.\-_]+$`).MatchString

// validatePropertyName validates if name property contains valid characters
func validatePropertyName(v interface{}, _ cty.Path) diag.Diagnostics {
	name := v.(string)
	maxPropertyNameLength := 85

	if len(name) > maxPropertyNameLength {
		return diag.Errorf("a name must be shorter than %d characters", maxPropertyNameLength+1)
	}
	if !isValidPropertyName(name) {
		return diag.Errorf("a name must only contain letters, numbers, and these characters: . _ -")
	}
	return nil
}

// rulesCustomDiff compares Rules.Criteria and Rules.Children fields from terraform state and from a new configuration.
// If some of these fields are empty lists in the new configuration and are nil in the terraform state, then this function
// returns no difference for these fields
func rulesCustomDiff(_ context.Context, diff *schema.ResourceDiff, _ interface{}) error {
	o, n := diff.GetChange("rules")

	oldValue := o.(string)
	newValue := n.(string)

	var oldRulesUpdate, newRulesUpdate papi.RulesUpdate

	if diff.Id() == "" && newValue != "" {
		rules, err := unifyRulesDiff(newValue)
		if err != nil {
			return err
		}
		if err = diff.SetNew("rules", rules); err != nil {
			return fmt.Errorf("cannot set a new diff value for 'rules' %s", err)
		}
		return nil
	}

	if oldValue == "" || newValue == "" {
		return nil
	}

	err := json.Unmarshal([]byte(oldValue), &oldRulesUpdate)
	if err != nil {
		return fmt.Errorf("cannot parse rules JSON from state: %s", err)
	}

	err = json.Unmarshal([]byte(newValue), &newRulesUpdate)
	if err != nil {
		return fmt.Errorf("cannot parse rules JSON from config: %s", err)
	}

	rules, err := compareFields(&oldRulesUpdate, &newRulesUpdate)
	if err != nil {
		return fmt.Errorf("cannot encode rules JSON %s", err)
	}
	rulesBytes, err := json.Marshal(newRulesUpdate)
	if err != nil {
		return err
	}
	rules = string(rulesBytes)

	if err = diff.SetNew("rules", rules); err != nil {
		return fmt.Errorf("cannot set a new diff value for 'rules' %s", err)
	}
	return nil
}

// unifyRulesDiff is invoked on first planning for property creation
// Its main purpose is to unify the rules JSON with what we expect will be created by PAPI
// It is used in order to prevent diffs on output on subsequent terraform applies
func unifyRulesDiff(newValue string) (string, error) {
	var newRulesUpdate papi.RulesUpdate
	err := json.Unmarshal([]byte(newValue), &newRulesUpdate)
	if err != nil {
		return "", fmt.Errorf("cannot parse rules JSON from config: %s", err)
	}
	rulesBytes, err := json.Marshal(newRulesUpdate)
	if err != nil {
		return "", err
	}
	return string(rulesBytes), nil
}

func compareFields(old, new *papi.RulesUpdate) (string, error) {
	if old.Rules.Children == nil && len(new.Rules.Children) == 0 {
		new.Rules.Children = old.Rules.Children
	}
	if old.Rules.Criteria == nil && len(new.Rules.Criteria) == 0 {
		new.Rules.Criteria = old.Rules.Criteria
	}
	rules, err := json.Marshal(new)
	return string(rules), err
}

func hostNamesCustomDiff(_ context.Context, d *schema.ResourceDiff, m interface{}) error {
	meta := akamai.Meta(m)
	logger := meta.Log("PAPI", "hostNamesCustomDiff")

	o, n := d.GetChange("hostnames")
	oldVal, ok := o.(*schema.Set)
	if !ok {
		logger.Errorf("error parsing local state for old value %s", oldVal)
		return fmt.Errorf("cannot parse hostnames state properly %v", o)
	}

	newVal, ok := n.(*schema.Set)
	if !ok {
		logger.Errorf("error parsing local state for new value %s", newVal)
		return fmt.Errorf("cannot parse hostnames state properly %v", n)
	}
	// PAPI doesn't allow hostnames to become empty if they already exist on server
	// TODO Do we add support for hostnames patch operation to enable this?
	if len(oldVal.List()) > 0 && len(newVal.List()) == 0 {
		logger.Errorf("Hostnames exist on server and cannot be updated to empty for %d", d.Id())
		return fmt.Errorf("hostnames exist on server and cannot be updated to empty for property with id '%s'. Provide at least one hostname to update existing list of hostnames associated to this property", d.Id())
	}
	return nil
}

// versionsComputedValuesCustomDiff sets `latest_version`, `staging_version` and `production_version` fields as computed
// if a new version of property is expected to be created
func versionsComputedValuesCustomDiff(_ context.Context, d *schema.ResourceDiff, m interface{}) error {
	meta := akamai.Meta(m)
	logger := meta.Log("PAPI", "versionsComputedValuesCustomDiff")
	oldRules, newRules := d.GetChange("rules")
	o, n := d.GetChange("hostnames")
	oldSet := o.(*schema.Set)
	equal := oldSet.HashEqual(n.(*schema.Set))
	if !equal || !compareRulesJSON(oldRules.(string), newRules.(string)) {
		// These computed attributes can be changed on server through other clients and the state needs to be synced to local
		for _, key := range []string{"latest_version", "staging_version", "production_version"} {
			err := d.SetNewComputed(key)
			if err != nil {
				logger.Errorf("%s state failed to update with new value from server", key)
				return fmt.Errorf("%w: %s", tools.ErrValueSet, err.Error())
			}
			logger.Debugf("%s state will be updated with new value from server", key)
		}
	}

	return nil
}
func resourcePropertyCreate(ctx context.Context, d *schema.ResourceData, m interface{}) diag.Diagnostics {
	meta := akamai.Meta(m)
	logger := meta.Log("PAPI", "resourcePropertyCreate")
	client := inst.Client(meta)
	ctx = log.NewContext(ctx, logger)

	// Block creation if user has set any hard-deprecated attributes
	for _, attr := range resPropForbiddenAttrs() {
		if _, ok := d.GetOk(attr); ok {
			return diag.Errorf("unsupported attribute: %q See the Akamai Terraform Upgrade Guide", attr)
		}
	}

	// Schema guarantees these types
	PropertyName := d.Get("name").(string)

	GroupID, err := tools.ResolveKeyStringState(d, "group_id", "group")
	if err != nil {
		return diag.FromErr(err)
	}
	GroupID = tools.AddPrefix(GroupID, "grp_")

	ContractID := d.Get("contract_id").(string)
	if ContractID == "" {
		ContractID = d.Get("contract").(string)
	}
	ContractID = tools.AddPrefix(ContractID, "ctr_")

	ProductID := d.Get("product_id").(string)
	if ProductID == "" {
		ProductID = d.Get("product").(string)
	}
	ProductID = tools.AddPrefix(ProductID, "prd_")

	RuleFormat := d.Get("rule_format").(string)

	RulesJSON := []byte(d.Get("rules").(string))

	PropertyID, err := createProperty(ctx, client, PropertyName, GroupID, ContractID, ProductID, RuleFormat)
	if err != nil {
		if strings.Contains(err.Error(), "\"statusCode\": 404") {
			// find out what is missing from the request
			if _, err = getGroup(ctx, meta, GroupID); err != nil {
				if errors.Is(err, ErrGroupNotFound) {
					return diag.Errorf("%v: %s", ErrGroupNotFound, GroupID)
				}
				return diag.FromErr(err)
			}
			if _, err = getContract(ctx, meta, ContractID); err != nil {
				if errors.Is(err, ErrContractNotFound) {
					return diag.Errorf("%v: %s", ErrContractNotFound, ContractID)
				}
				return diag.FromErr(err)
			}
			if _, err = getProduct(ctx, meta, ProductID, ContractID); err != nil {
				if errors.Is(err, ErrProductNotFound) {
					return diag.Errorf("%v: %s", ErrProductNotFound, ProductID)
				}
				return diag.FromErr(err)
			}
			return diag.FromErr(err)
		}
		return diag.FromErr(err)
	}

	// Save minimum state BEFORE moving on
	d.SetId(PropertyID)
	attrs := map[string]interface{}{
		"group_id":    GroupID,
		"contract_id": ContractID,
		"product_id":  ProductID,
		"product":     ProductID,
	}
	if err := rdSetAttrs(ctx, d, attrs); err != nil {
		return diag.FromErr(err)
	}

	Property := papi.Property{
		PropertyName:  PropertyName,
		PropertyID:    PropertyID,
		ContractID:    ContractID,
		GroupID:       GroupID,
		ProductID:     ProductID,
		LatestVersion: 1,
	}
	HostnameVal, err := tools.GetSetValue("hostnames", d)
	if err == nil {
		Hostnames := mapToHostnames(HostnameVal.List())
		if len(Hostnames) > 0 {
			if err := updatePropertyHostnames(ctx, client, Property, Hostnames); err != nil {
				return diag.FromErr(err)
			}
		}
	} else {
		logger.Warnf("hostnames not set in ResourceData: %s", err.Error())
	}

	if len(RulesJSON) > 0 {
		var Rules papi.RulesUpdate
		if err := json.Unmarshal(RulesJSON, &Rules); err != nil {
			logger.WithError(err).Error("failed to unmarshal property rules")
			return diag.Errorf("rules are not valid JSON: %s", err)
		}

		ctx := ctx
		if RuleFormat != "" {
			h := http.Header{
				"Content-Type": []string{fmt.Sprintf("application/vnd.akamai.papirules.%s+json", RuleFormat)},
			}

			ctx = session.ContextWithOptions(ctx, session.WithContextHeaders(h))
		}

		if err := updatePropertyRules(ctx, client, Property, Rules); err != nil {
			d.Partial(true)
			return diag.FromErr(err)
		}
	}

	return resourcePropertyRead(ctx, d, m)
}

func resourcePropertyRead(ctx context.Context, d *schema.ResourceData, m interface{}) diag.Diagnostics {
	ctx = log.NewContext(ctx, akamai.Meta(m).Log("PAPI", "resourcePropertyRead"))
	logger := log.FromContext(ctx)
	client := inst.Client(akamai.Meta(m))

	// Schema guarantees group_id, and contract_id are strings
	PropertyID := d.Id()
	ContractID := tools.AddPrefix(d.Get("contract_id").(string), "ctr_")
	GroupID := tools.AddPrefix(d.Get("group_id").(string), "grp_")
	ReadVersionID := d.Get("read_version").(int)

	var Property *papi.Property
	var err error
	var v int
	if ReadVersionID == 0 {
		Property, err = fetchLatestProperty(ctx, client, PropertyID, GroupID, ContractID)
	} else {
		Property, v, err = fetchProperty(ctx, client, PropertyID, GroupID, ContractID, strconv.Itoa(ReadVersionID))
	}
	if err != nil {
		return diag.FromErr(err)
	}
	if v == 0 {
		// use latest version unless "read_version" != 0
		v = Property.LatestVersion
	}

	var StagingVersion int
	if Property.StagingVersion != nil {
		StagingVersion = *Property.StagingVersion
	}

	var ProductionVersion int
	if Property.ProductionVersion != nil {
		ProductionVersion = *Property.ProductionVersion
	}

	// TODO: Load hostnames asynchronously
	Hostnames, err := fetchPropertyVersionHostnames(ctx, client, *Property, v)
	if err != nil {
		return diag.FromErr(err)
	}

	// TODO: Load rules asynchronously
	Rules, RuleFormat, RuleErrors, RuleWarnings, err := fetchPropertyVersionRules(ctx, client, *Property, v)
	if err != nil {
		return diag.FromErr(err)
	}
	if len(RuleErrors) > 0 {
		if err := d.Set("rule_errors", papiErrorsToList(RuleErrors)); err != nil {
			return diag.FromErr(fmt.Errorf("%w: %s", tools.ErrValueSet, err.Error()))
		}
		msg, err := json.MarshalIndent(papiErrorsToList(RuleErrors), "", "\t")
		if err != nil {
			return diag.FromErr(fmt.Errorf("error marshaling API error: %s", err))
		}
		logger.Errorf("Property has rule errors %s", msg)
	}
	if len(RuleWarnings) > 0 {
		msg, err := json.MarshalIndent(papiErrorsToList(RuleWarnings), "", "\t")
		if err != nil {
			return diag.FromErr(fmt.Errorf("error marshaling API warnings: %s", err))
		}
		logger.Warnf("Property has rule warnings %s", msg)
	}

	RulesJSON, err := json.Marshal(Rules)
	if err != nil {
		logger.WithError(err).Error("could not render rules as JSON")
		return diag.Errorf("received rules that could not be rendered to JSON: %s", err)
	}
	res, err := fetchPropertyVersion(ctx, client, PropertyID, GroupID, ContractID, v)
	if err != nil {
		return diag.FromErr(err)
	}
	Property.ProductID = res.Version.ProductID

	attrs := map[string]interface{}{
		"name":               Property.PropertyName,
		"group_id":           Property.GroupID,
		"group":              Property.GroupID,
		"contract_id":        Property.ContractID,
		"contract":           Property.ContractID,
		"latest_version":     Property.LatestVersion,
		"staging_version":    StagingVersion,
		"production_version": ProductionVersion,
		"hostnames":          flattenHostnames(Hostnames),
		"rules":              string(RulesJSON),
		"rule_format":        RuleFormat,
		"rule_errors":        papiErrorsToList(RuleErrors),
		"read_version":       v,
	}
	if Property.ProductID != "" {
		attrs["product_id"] = Property.ProductID
		attrs["product"] = Property.ProductID
	}
	if err := rdSetAttrs(ctx, d, attrs); err != nil {
		return diag.FromErr(err)
	}

	return nil
}

func resourcePropertyUpdate(ctx context.Context, d *schema.ResourceData, m interface{}) diag.Diagnostics {
	ctx = log.NewContext(ctx, akamai.Meta(m).Log("PAPI", "resourcePropertyUpdate"))
	logger := log.FromContext(ctx)
	client := inst.Client(akamai.Meta(m))

	// Block changes to hard-deprecated attributes
	for _, attr := range resPropForbiddenAttrs() {
		if _, ok := d.GetOk(attr); ok && d.HasChange(attr) {
			d.Partial(true)
			return diag.Errorf("unsupported attribute: %q See the Akamai Terraform Upgrade Guide", attr)
		}
	}

	diags := diag.Diagnostics{}

	immutable := []string{
		"group_id",
		"group",
		"contract_id",
		"contract",
		"product_id",
		"product",
	}
	for _, attr := range immutable {
		if d.HasChange(attr) {
			err := fmt.Errorf(`property attribute %q cannot be changed after creation (immutable)`, attr)
			logger.WithError(err).Error("could not update property")
			diags = append(diags, diag.FromErr(err)...)
		}
	}
	if diags.HasError() {
		d.Partial(true)
		return diags
	}

	// We only update if these attributes change.
	if !d.HasChanges("hostnames", "rules", "rule_format") {
		logger.Debug("No changes to hostnames, rules, or rule_format (no update required)")
		return nil
	}

	// Schema guarantees these types
	var StagingVersion, ProductionVersion *int
	if v, ok := d.GetOk("staging_version"); ok && v.(int) != 0 {
		i := v.(int)
		StagingVersion = &i
	}

	if v, ok := d.GetOk("production_version"); ok && v.(int) != 0 {
		i := v.(int)
		ProductionVersion = &i
	}

	Property := papi.Property{
		PropertyID:        d.Id(),
		PropertyName:      d.Get("name").(string),
		ContractID:        d.Get("contract_id").(string),
		GroupID:           d.Get("group_id").(string),
		ProductID:         d.Get("product_id").(string),
		LatestVersion:     d.Get("latest_version").(int),
		StagingVersion:    StagingVersion,
		ProductionVersion: ProductionVersion,
	}

	// Schema guarantees group_id, and contract_id are strings
	PropertyID := d.Id()
	ContractID := d.Get("contract_id").(string)
	GroupID := d.Get("group_id").(string)

	var PropertyVersion int
	if v, ok := d.GetOk("read_version"); ok && v.(int) != 0 {
		PropertyVersion = v.(int)
	} else {
		PropertyVersion = Property.LatestVersion
	}

	resp, err := fetchPropertyVersion(ctx, client, PropertyID, GroupID, ContractID, PropertyVersion)
	if err != nil {
		d.Partial(true)
		return diag.FromErr(err)
	}

	// check latest version is editable
	if resp.Version.ProductionStatus != papi.VersionStatusInactive || resp.Version.StagingStatus != papi.VersionStatusInactive {
		// The latest version has been activated on either production or staging, so we need to create a new version to apply changes on
		VersionID, err := createPropertyVersion(ctx, client, Property)
		if err != nil {
			d.Partial(true)
			return diag.FromErr(err)
		}
		Property.LatestVersion = VersionID
		if err = d.Set("read_version", 0); err != nil {
			return diag.FromErr(err)
		}
	}

	// Hostnames
	if d.HasChange("hostnames") {
		HostnameVal, err := tools.GetSetValue("hostnames", d)
		if err == nil {
			Hostnames := mapToHostnames(HostnameVal.List())
			if len(Hostnames) > 0 {
				if err := updatePropertyHostnames(ctx, client, Property, Hostnames); err != nil {
					d.Partial(true)
					return diag.FromErr(err)
				}
			}
		} else {
			logger.Warnf("hostnames not set in ResourceData: %s", err.Error())
		}
	}

	RuleFormat := d.Get("rule_format").(string)
	RulesJSON := []byte(d.Get("rules").(string))
	RulesNeedUpdate := len(RulesJSON) > 0 && d.HasChange("rules")
	FormatNeedsUpdate := len(RuleFormat) > 0 && d.HasChange("rule_format")

	if FormatNeedsUpdate || RulesNeedUpdate {
		var Rules papi.RulesUpdate
		if err := json.Unmarshal(RulesJSON, &Rules); err != nil {
			d.Partial(true)
			return diag.Errorf("rules are not valid JSON: %s", err)
		}

		MIME := fmt.Sprintf("application/vnd.akamai.papirules.%s+json", RuleFormat)
		h := http.Header{"Content-Type": []string{MIME}}
		ctx := session.ContextWithOptions(ctx, session.WithContextHeaders(h))

		if err := updatePropertyRules(ctx, client, Property, Rules); err != nil {
			d.Partial(true)
			return diag.FromErr(err)
		}
	}

	return resourcePropertyRead(ctx, d, m)
}

func resourcePropertyDelete(ctx context.Context, d *schema.ResourceData, m interface{}) diag.Diagnostics {
	ctx = log.NewContext(ctx, akamai.Meta(m).Log("PAPI", "resourcePropertyDelete"))
	client := inst.Client(akamai.Meta(m))

	PropertyID := d.Id()
	ContractID := tools.AddPrefix(d.Get("contract_id").(string), "ctr_")
	GroupID := tools.AddPrefix(d.Get("group_id").(string), "grp_")

	if err := removeProperty(ctx, client, PropertyID, GroupID, ContractID); err != nil {
		return diag.FromErr(err)
	}

	return nil
}

func resourcePropertyImport(ctx context.Context, d *schema.ResourceData, m interface{}) ([]*schema.ResourceData, error) {
	ctx = log.NewContext(ctx, akamai.Meta(m).Log("PAPI", "resourcePropertyImport"))

	// User-supplied import ID is a comma-separated list of PropertyID[,GroupID[,ContractID]]
	// ContractID and GroupID are optional as long as the PropertyID is sufficient to fetch the property
	var PropertyID, GroupID, ContractID, Version string
	parts := strings.Split(d.Id(), ",")
	switch len(parts) {
	case 4:
		Version = parts[3]
		fallthrough
	case 3:
		PropertyID = tools.AddPrefix(parts[0], "prp_")
		ContractID = tools.AddPrefix(parts[1], "ctr_")
		GroupID = tools.AddPrefix(parts[2], "grp_")
	case 2:
		Version = parts[1]
		fallthrough
	case 1:
		PropertyID = tools.AddPrefix(parts[0], "prp_")

	default:
		return nil, fmt.Errorf("invalid property identifier: %q", d.Id())
	}

	// Import only needs to set the resource ID and enough attributes that the read operation will function, so there's
	// no need to fetch anything if the user gave both GroupID and ContractID
	if GroupID != "" && ContractID != "" {
		attrs := map[string]interface{}{
			"group_id":    GroupID,
			"contract_id": ContractID,
		}

		// if we also get the optional Version parameter, we need to parse it and set it in the schema
		if !isDefaultVersion(Version) {
			if v, err := parseVersionNumber(Version); err != nil {
				// acceptable values for Version at this point: "PRODUCTION" or "STAGING" (or synonyms). Let's validate
				if _, err := NetworkAlias(Version); err != nil {
					return nil, ErrPropertyVersionNotFound
				}
				// if we ran validation and we actually have a network name, we still need to fetch the desired version number
				_, attrs["read_version"], err = fetchProperty(ctx, inst.Client(akamai.Meta(m)), PropertyID, GroupID, ContractID, Version)
				if err != nil {
					return nil, err
				}
			} else {
				// if the version number can be parsed as a number or ver_#, nothing else to be done
				attrs["read_version"] = v
			}
		}
		if err := rdSetAttrs(ctx, d, attrs); err != nil {
			return nil, err
		}

		d.SetId(PropertyID)
		return []*schema.ResourceData{d}, nil
	}

	var err error
	var Property *papi.Property
	var v int
	if !isDefaultVersion(Version) {
		Property, v, err = fetchProperty(ctx, inst.Client(akamai.Meta(m)), PropertyID, GroupID, ContractID, Version)
	} else {
		Property, err = fetchLatestProperty(ctx, inst.Client(akamai.Meta(m)), PropertyID, GroupID, ContractID)
	}
	if err != nil {
		return nil, err
	}

	attrs := map[string]interface{}{
		"group_id":     Property.GroupID,
		"contract_id":  Property.ContractID,
		"read_version": v,
	}
	if err := rdSetAttrs(ctx, d, attrs); err != nil {
		return nil, err
	}

	d.SetId(Property.PropertyID)
	return []*schema.ResourceData{d}, nil
}

func isDefaultVersion(version string) bool {
	return version == "" || strings.ToLower(version) == "latest"
}

var versionRegexp = regexp.MustCompile(`^ver_(\d+)$`)

// parse a version number (format "ver_#" or "#") or throw an error
func parseVersionNumber(version string) (int, error) {
	v := tools.AddPrefix(version, "ver_")
	r := versionRegexp
	matches := r.FindStringSubmatch(v)
	if len(matches) < 2 {
		return 0, fmt.Errorf("invalid version number")
	}
	versionNumber, err := strconv.Atoi(matches[1])
	return versionNumber, err
}

func resPropForbiddenAttrs() []string {
	return []string{
		"cp_code",
		"contact",
		"origin",
		"is_secure",
		"variables",
	}
}

func createProperty(ctx context.Context, client papi.PAPI, PropertyName, GroupID, ContractID, ProductID, RuleFormat string) (PropertyID string, err error) {
	req := papi.CreatePropertyRequest{
		ContractID: ContractID,
		GroupID:    GroupID,
		Property: papi.PropertyCreate{
			ProductID:    ProductID,
			PropertyName: PropertyName,
			RuleFormat:   RuleFormat,
		},
	}

	logger := log.FromContext(ctx).WithFields(logFields(req))

	logger.Debug("creating property")
	res, err := client.CreateProperty(ctx, req)
	if err != nil {
		logger.WithError(err).Error("could not create property")
		return
	}
	PropertyID = res.PropertyID

	logger.WithFields(logFields(*res)).Info("property created")
	return
}

func removeProperty(ctx context.Context, client papi.PAPI, PropertyID, GroupID, ContractID string) error {
	req := papi.RemovePropertyRequest{
		PropertyID: PropertyID,
		GroupID:    GroupID,
		ContractID: ContractID,
	}

	logger := log.FromContext(ctx).WithFields(logFields(req))
	logger.Debug("removing property")

	_, err := client.RemoveProperty(ctx, req)
	if err != nil {
		logger.WithError(err).Error("could not remove property")
		return err
	}

	logger.Info("property removed")

	return nil
}

func fetchLatestProperty(ctx context.Context, client papi.PAPI, PropertyID, GroupID, ContractID string) (*papi.Property, error) {
	req := papi.GetPropertyRequest{
		PropertyID: PropertyID,
		ContractID: ContractID,
		GroupID:    GroupID,
	}
	logger := log.FromContext(ctx).WithFields(logFields(req))
	logger.Debug("fetching property")
	res, err := client.GetProperty(ctx, req)
	if err != nil {
		logger.WithError(err).Error("could not read property")
		return nil, err
	}

	logger = logger.WithFields(logFields(*res))

	if res.Property == nil {
		err := fmt.Errorf("PAPI::GetProperty() response did not contain a property")
		logger.WithError(err).Error("could not look up property")
		return nil, err
	}

	logger.Debug("property fetched")
	return res.Property, nil
}

// fetchProperty Retrieves basic info for a Property
func fetchProperty(ctx context.Context, client papi.PAPI, PropertyID, GroupID, ContractID, version string) (*papi.Property, int, error) {
	req := papi.GetPropertyVersionsRequest{
		PropertyID: PropertyID,
		ContractID: ContractID,
		GroupID:    GroupID,
	}
	logger := log.FromContext(ctx).WithFields(logFields(req))
	logger.Debugf("fetching property versions")
	res, err := client.GetPropertyVersions(ctx, req)
	if err != nil {
		logger.WithError(err).Error("could not read property versions")
		return nil, 0, err
	}

	versions := res.Versions.Items
	var versionNumber int
	if network, err := NetworkAlias(version); err != nil {
		// if it is a valid version number there is nothing else to do
		n, err := parseVersionNumber(version)
		if err != nil {
			return nil, 0, ErrPropertyVersionNotFound
		}
		versionNumber = n
	} else {
		// filter production
		if network == string(papi.ActivationNetworkProduction) {
			versions, err = filterProduction(versions)
			if err != nil {
				return nil, 0, err
			}
		}

		// filter staging
		if network == string(papi.ActivationNetworkStaging) {
			versions, err = filterStaging(versions)
			if err != nil {
				return nil, 0, err
			}
		}

		versionNumber = getLatestVersionNumber(versions)
	}
	versionItem, err := getVersionItem(versions, versionNumber)
	if err != nil {
		return nil, 0, err
	}

	property := papi.Property{
		AccountID:         res.AccountID,
		ContractID:        res.ContractID,
		GroupID:           res.GroupID,
		PropertyID:        res.PropertyID,
		PropertyName:      res.PropertyName,
		LatestVersion:     getLatestVersionNumber(res.Versions.Items),
		StagingVersion:    getNetworkActiveVersionNumber(res.Versions.Items, string(papi.ActivationNetworkStaging)),
		ProductionVersion: getNetworkActiveVersionNumber(res.Versions.Items, string(papi.ActivationNetworkProduction)),
		AssetID:           res.AssetID,
		Note:              versionItem.Note,
		ProductID:         versionItem.ProductID,
		RuleFormat:        versionItem.RuleFormat,
	}

	logger.Debug("property versions fetched")

	return &property, versionNumber, nil
}

// filterStaging filters papi.PropertyVersionGetItem elements with StagingStatus == "ACTIVE"
// from the given list
func filterStaging(items []papi.PropertyVersionGetItem) ([]papi.PropertyVersionGetItem, error) {
	var output []papi.PropertyVersionGetItem
	for _, it := range items {
		if it.StagingStatus == "ACTIVE" {
			output = append(output, it)
		}
	}
	if len(output) == 0 {
		return nil, ErrPropertyVersionNotFound
	}
	return output, nil
}

// filterProduction filters papi.PropertyVersionGetItem elements with ProductionStatus == "ACTIVE"
// from the given list
func filterProduction(items []papi.PropertyVersionGetItem) ([]papi.PropertyVersionGetItem, error) {
	var output []papi.PropertyVersionGetItem
	for _, it := range items {
		if it.ProductionStatus == "ACTIVE" {
			output = append(output, it)
		}
	}
	if len(output) == 0 {
		return nil, ErrPropertyVersionNotFound
	}
	return output, nil
}

// getLatestVersionNumber returns from the given list the highest papi.PropertyVersionGetItem
// PropertyVersion from the list
func getLatestVersionNumber(items []papi.PropertyVersionGetItem) int {
	var latest int
	for _, it := range items {
		if it.PropertyVersion > latest {
			latest = it.PropertyVersion
		}
	}
	return latest
}

// getNetworkActiveVersionNumber returns from the given list the *papi.PropertyVersionGetItem
// active in the given network
func getNetworkActiveVersionNumber(items []papi.PropertyVersionGetItem, network string) *int {
	for _, it := range items {
		switch network {
		case string(papi.ActivationNetworkStaging):
			if it.StagingStatus == "ACTIVE" {
				return &it.PropertyVersion
			}
		case string(papi.ActivationNetworkProduction):
			if it.ProductionStatus == "ACTIVE" {
				return &it.PropertyVersion
			}
		}
	}
	return nil
}

func getVersionItem(items []papi.PropertyVersionGetItem, versionNumber int) (*papi.PropertyVersionGetItem, error) {
	for _, it := range items {
		if it.PropertyVersion == versionNumber {
			return &it, nil
		}
	}
	return nil, ErrPropertyVersionNotFound
}

// load status for what we currently have as a given property version.  GetLatestVersion may also work here.
func fetchPropertyVersion(ctx context.Context, client papi.PAPI, PropertyID, GroupID, ContractID string, PropertyVersion int) (*papi.GetPropertyVersionsResponse, error) {
	req := papi.GetPropertyVersionRequest{
		PropertyID:      PropertyID,
		ContractID:      ContractID,
		GroupID:         GroupID,
		PropertyVersion: PropertyVersion,
	}
	logger := log.FromContext(ctx).WithFields(logFields(req))
	logger.Debug("fetching property version")

	res, err := client.GetPropertyVersion(ctx, req)
	if err != nil {
		logger.WithError(err).Error("could not read property version")
		return nil, err
	}
	logger = logger.WithFields(logFields(*res))
	logger.Debug("property version fetched")
	return res, err
}

// Fetch hostnames for latest version of given property
func fetchPropertyVersionHostnames(ctx context.Context, client papi.PAPI, Property papi.Property, version int) ([]papi.Hostname, error) {
	req := papi.GetPropertyVersionHostnamesRequest{
		PropertyID:        Property.PropertyID,
		GroupID:           Property.GroupID,
		ContractID:        Property.ContractID,
		PropertyVersion:   version,
		IncludeCertStatus: true,
	}

	logger := log.FromContext(ctx).WithFields(logFields(req))

	logger.Debug("fetching property hostnames")
	res, err := client.GetPropertyVersionHostnames(ctx, req)
	if err != nil {
		logger.WithError(err).Error("could not fetch property hostnames")
		return nil, err
	}

	logger.WithFields(logFields(*res)).Debug("fetched property hostnames")
	return res.Hostnames.Items, nil
}

// Fetch rules for latest version of given property
func fetchPropertyVersionRules(ctx context.Context, client papi.PAPI, Property papi.Property, version int) (Rules papi.RulesUpdate, Format string, Errors, Warnings []*papi.Error, err error) {
	req := papi.GetRuleTreeRequest{
		PropertyID:      Property.PropertyID,
		GroupID:         Property.GroupID,
		ContractID:      Property.ContractID,
		PropertyVersion: version,
		ValidateRules:   true,
		ValidateMode:    papi.RuleValidateModeFull,
	}

	logger := log.FromContext(ctx).WithFields(logFields(req))

	logger.Debug("fetching property rules")
	res, err := client.GetRuleTree(ctx, req)
	if err != nil {
		logger.WithError(err).Error("could not fetch property rules")
		return
	}

	logger.WithFields(logFields(*res)).Debug("fetched property rules")
	Rules = papi.RulesUpdate{
		Rules:    res.Rules,
		Comments: res.Comments,
	}
	Format = res.RuleFormat
	Errors = res.Errors
	Warnings = res.Warnings
	return
}

// Set rules for the latest version of the given property
func updatePropertyRules(ctx context.Context, client papi.PAPI, Property papi.Property, Rules papi.RulesUpdate) error {
	req := papi.UpdateRulesRequest{
		PropertyID:      Property.PropertyID,
		GroupID:         Property.GroupID,
		ContractID:      Property.ContractID,
		PropertyVersion: Property.LatestVersion,
		Rules:           Rules,
		ValidateRules:   true,
	}

	logger := log.FromContext(ctx).WithFields(logFields(req))

	logger.Debug("fetching property rules")
	res, err := client.UpdateRuleTree(ctx, req)
	if err != nil {
		logger.WithError(err).Error("could not update property rules")
		return err
	}

	logger.WithFields(logFields(*res)).Info("updated property rules")
	return nil
}

// Create a new property version based on the latest version of the given property
func createPropertyVersion(ctx context.Context, client papi.PAPI, Property papi.Property) (NewVersion int, err error) {
	req := papi.CreatePropertyVersionRequest{
		PropertyID: Property.PropertyID,
		ContractID: Property.ContractID,
		GroupID:    Property.GroupID,
		Version: papi.PropertyVersionCreate{
			CreateFromVersion: Property.LatestVersion,
		},
	}

	logger := log.FromContext(ctx).WithFields(logFields(req))

	logger.Debug("creating new property version")
	res, err := client.CreatePropertyVersion(ctx, req)
	if err != nil {
		logger.WithError(err).Error("could not create new property version")
		return
	}

	logger.WithFields(logFields(*res)).Info("property version created")
	NewVersion = res.PropertyVersion
	return
}

// Set hostnames of the latest version of the given property
func updatePropertyHostnames(ctx context.Context, client papi.PAPI, Property papi.Property, Hostnames []papi.Hostname) error {
	if Hostnames == nil {
		Hostnames = []papi.Hostname{}
	}
	req := papi.UpdatePropertyVersionHostnamesRequest{
		PropertyID:      Property.PropertyID,
		GroupID:         Property.GroupID,
		ContractID:      Property.ContractID,
		PropertyVersion: Property.LatestVersion,
		Hostnames:       Hostnames,
	}

	logger := log.FromContext(ctx).WithFields(logFields(req))

	logger.Debug("updating property hostnames")
	res, err := client.UpdatePropertyVersionHostnames(ctx, req)
	if err != nil {
		hasDefaultProvisioningType := false
		for _, h := range Hostnames {
			if h.CertProvisioningType == "DEFAULT" {
				hasDefaultProvisioningType = true
				break
			}
		}
		var e *papi.Error
		if hasDefaultProvisioningType && errors.As(err, &e) {
			if e.StatusCode == http.StatusForbidden && e.Type == "https://problems.luna.akamaiapis.net/papi/v0/property-version-hostname/default-cert-provisioning-unavailable" {
				err = fmt.Errorf("%s: not possible to use cert_provisioning_type = 'DEFAULT' as secure-by-default is not enabled in this account", papi.ErrUpdatePropertyVersionHostnames)
			}
			if e.StatusCode == http.StatusTooManyRequests && e.LimitKey == "DEFAULT_CERTS_PER_CONTRACT" && e.Remaining == 0 {
				err = fmt.Errorf("%s: not possible to use cert_provisioning_type = 'DEFAULT' as the limit for DEFAULT certificates has been reached", papi.ErrUpdatePropertyVersionHostnames)
			}
		}
		logger.WithError(err).Error("could not create new property version")
		return err
	}

	logger.WithFields(logFields(*res)).Info("property hostnames updated")
	return nil
}

// Convert the given map from a schema.ResourceData to a slice of papi.Hostnames /input to papi request
func mapToHostnames(givenList []interface{}) []papi.Hostname {
	var Hostnames []papi.Hostname

	for _, givenMap := range givenList {
		var r = givenMap.(map[string]interface{})
		cnameFrom := r["cname_from"]
		cnameTo := r["cname_to"]
		certProvisioningType := r["cert_provisioning_type"]
		if len(r) != 0 {
			Hostnames = append(Hostnames, papi.Hostname{
				CnameType:            "EDGE_HOSTNAME",
				CnameFrom:            cnameFrom.(string),
				CnameTo:              cnameTo.(string), // guaranteed by schema to be a string
				CertProvisioningType: certProvisioningType.(string),
			})
		}
	}
	return Hostnames
}

// Set many attributes of a schema.ResourceData in one call
func rdSetAttrs(ctx context.Context, d *schema.ResourceData, AttributeValues map[string]interface{}) error {
	logger := log.FromContext(ctx)

	for attr, value := range AttributeValues {
		if err := d.Set(attr, value); err != nil {
			logger.WithError(err).Errorf("could not set %q", attr)
			return err
		}
	}

	return nil
}
