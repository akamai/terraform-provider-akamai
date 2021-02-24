package property

import (
	"context"
	"encoding/json"
	"errors"
	"fmt"
	"net/http"
	"regexp"
	"strings"

	"github.com/apex/log"
	"github.com/hashicorp/go-cty/cty"
	"github.com/hashicorp/terraform-plugin-sdk/v2/diag"
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
				ValidateDiagFunc: tools.IsNotBlank,
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
				Deprecated: `use "group_id" attribute instead`,
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
				Deprecated: `use "contract_id" attribute instead`,
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
				Type:          schema.TypeString,
				Optional:      true,
				Computed:      true,
				ConflictsWith: []string{"product_id"},
				Deprecated:    `use "product_id" attribute instead`,
				StateFunc:     addPrefixToState("prd_"),
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
				Type:     schema.TypeList,
				Optional: true,
				Elem: &schema.Resource{
					Schema: map[string]*schema.Schema{
						"cname_from": {
							Type:     schema.TypeString,
							Optional: true,
						},
						"cname_to": {
							Type:     schema.TypeString,
							Optional: true,
						},
						"cert_provisioning_type": {
							Type:     schema.TypeString,
							Optional: true,
						},
					},
				},
			},
			//"hostnames": {
			//	Type:        schema.TypeMap,
			//	Optional:    true,
			//	Elem:        &schema.Schema{Type: schema.TypeString},
			//	Description: "Mapping of edge hostname CNAMEs to other CNAMEs",
			//},

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
			"rule_errors": {
				Type:     schema.TypeList,
				Computed: true,
				Elem:     papiError(),
			},
			"rule_warnings": {
				Type:     schema.TypeList,
				Computed: true,
				Elem:     papiError(),
			},

			// Hard-deprecated attributes: These are effectively removed, but we wanted to refer users to the upgrade guide
			"cp_code": {
				Type:       schema.TypeString,
				Optional:   true,
				Deprecated: `"cp_code" is no longer supported by this resource type - See Akamai Terraform Upgrade Guide`,
			},
			"contact": {
				Type:       schema.TypeSet,
				Optional:   true,
				Elem:       &schema.Schema{Type: schema.TypeString},
				Deprecated: `"contact" is no longer supported by this resource type - See Akamai Terraform Upgrade Guide`,
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
				Deprecated: `"origin" is no longer supported by this resource type - See Akamai Terraform Upgrade Guide`,
			},
			"is_secure": {
				Type:       schema.TypeBool,
				Optional:   true,
				Deprecated: `"is_secure" is no longer supported by this resource type - See Akamai Terraform Upgrade Guide`,
			},
			"variables": {
				Type:       schema.TypeString,
				Optional:   true,
				Deprecated: `"variables" is no longer supported by this resource type - See Akamai Terraform Upgrade Guide`,
			},
		},
	}
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
		if ProductID == "" {
			return diag.Errorf("one of product,product_id must be specified")
		}
	}
	ProductID = tools.AddPrefix(ProductID, "prd_")
	Hostnames := mapToHostnames(d.Get("hostnames").([]interface{}))
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

	if len(Hostnames) > 0 {
		if err := updatePropertyHostnames(ctx, client, Property, Hostnames); err != nil {
			return diag.FromErr(err)
		}
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

	Property, err := fetchProperty(ctx, client, PropertyID, GroupID, ContractID)
	if err != nil {
		return diag.FromErr(err)
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
	Hostnames, err := fetchPropertyHostnames(ctx, client, *Property)
	if err != nil {
		return diag.FromErr(err)
	}

	// TODO: Load rules asynchronously
	Rules, RuleFormat, RuleErrors, RuleWarnings, err := fetchPropertyRules(ctx, client, *Property)
	if err != nil {
		return diag.FromErr(err)
	}

	RulesJSON, err := json.Marshal(Rules)
	if err != nil {
		logger.WithError(err).Error("could not render rules as JSON")
		return diag.Errorf("received rules that could not be rendered to JSON: %s", err)
	}

	attrs := map[string]interface{}{
		"name":               Property.PropertyName,
		"group_id":           Property.GroupID,
		"group":              Property.GroupID,
		"contract_id":        Property.ContractID,
		"contract":           Property.ContractID,
		"latest_version":     Property.LatestVersion,
		"staging_version":    StagingVersion,
		"production_version": ProductionVersion,
		"hostnames":          hostnamesToMap(Hostnames),
		"rules":              string(RulesJSON),
		"rule_format":        RuleFormat,
		"rule_errors":        papiErrorsToList(RuleErrors),
		"rule_warnings":      papiErrorsToList(RuleWarnings),
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

	var diags diag.Diagnostics

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

	// load status for what we currently have as latest version.  GetLatestVersion may also work here.
	resp, err := client.GetPropertyVersion(ctx, papi.GetPropertyVersionRequest{
		PropertyID:      d.Id(),
		PropertyVersion: d.Get("latest_version").(int),
		ContractID:      d.Get("contract_id").(string),
		GroupID:         d.Get("group_id").(string),
	})
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
	}

	// Hostnames
	if d.HasChange("hostnames") {
		Hostnames := mapToHostnames(d.Get("hostnames").([]interface{}))

		if err := updatePropertyHostnames(ctx, client, Property, Hostnames); err != nil {
			d.Partial(true)
			return diag.FromErr(err)
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
	ContractID := d.Get("contract_id").(string)
	GroupID := d.Get("group_id").(string)

	if err := removeProperty(ctx, client, PropertyID, GroupID, ContractID); err != nil {
		return diag.FromErr(err)
	}

	return nil
}

func resourcePropertyImport(ctx context.Context, d *schema.ResourceData, m interface{}) ([]*schema.ResourceData, error) {
	ctx = log.NewContext(ctx, akamai.Meta(m).Log("PAPI", "resourcePropertyImport"))

	// User-supplied import ID is a comma-separated list of PropertyID[,GroupID[,ContractID]]
	// ContractID and GroupID are optional as long as the PropertyID is sufficient to fetch the property
	var PropertyID, GroupID, ContractID string
	parts := strings.Split(d.Id(), ",")
	if len(parts) == 2 {
		return nil, fmt.Errorf("Either PropertyId or comma-separated list of PropertyId, contractID and groupID in that order has to be supplied in import: %s", d.Id())
	}
	switch len(parts) {
	case 1:
		PropertyID = tools.AddPrefix(parts[0], "prp_")
	case 3:
		PropertyID = tools.AddPrefix(parts[0], "prp_")
		ContractID = tools.AddPrefix(parts[1], "ctr_")
		GroupID = tools.AddPrefix(parts[2], "grp_")

	default:
		return nil, fmt.Errorf("invalid property identifier: %q", d.Id())
	}

	// Import only needs to set the resource ID and enough attributes that the read opertaion will function, so there's
	// no need to fetch anything if the user gave both GroupID and ContractID
	if GroupID != "" && ContractID != "" {
		attrs := map[string]interface{}{
			"group_id":    GroupID,
			"contract_id": ContractID,
		}
		if err := rdSetAttrs(ctx, d, attrs); err != nil {
			return nil, err
		}

		d.SetId(PropertyID)
		return []*schema.ResourceData{d}, nil
	}

	// Missing GroupID, ContractID, or both -- Attempt to fetch them. If the PropertyID is not sufficient, PAPI
	// will return an error.
	Property, err := fetchProperty(ctx, inst.Client(akamai.Meta(m)), PropertyID, GroupID, ContractID)
	if err != nil {
		return nil, err
	}

	attrs := map[string]interface{}{
		"group_id":    Property.GroupID,
		"contract_id": Property.ContractID,
	}
	if err := rdSetAttrs(ctx, d, attrs); err != nil {
		return nil, err
	}

	d.SetId(Property.PropertyID)
	return []*schema.ResourceData{d}, nil
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

// Retrieves basic info for a Property
func fetchProperty(ctx context.Context, client papi.PAPI, PropertyID, GroupID, ContractID string) (*papi.Property, error) {
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

// Fetch hostnames for latest version of given property
func fetchPropertyHostnames(ctx context.Context, client papi.PAPI, Property papi.Property) ([]papi.Hostname, error) {
	req := papi.GetPropertyVersionHostnamesRequest{
		PropertyID:      Property.PropertyID,
		GroupID:         Property.GroupID,
		ContractID:      Property.ContractID,
		PropertyVersion: Property.LatestVersion,
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
func fetchPropertyRules(ctx context.Context, client papi.PAPI, Property papi.Property) (Rules papi.RulesUpdate, Format string, Errors, Warnings []*papi.Error, err error) {
	req := papi.GetRuleTreeRequest{
		PropertyID:      Property.PropertyID,
		GroupID:         Property.GroupID,
		ContractID:      Property.ContractID,
		PropertyVersion: Property.LatestVersion,
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
		logger.WithError(err).Error("could not create new property version")
		return err
	}

	logger.WithFields(logFields(*res)).Info("property hostnames updated")
	return nil
}

// Convert given hostnames to the map form that can be stored in a schema.ResourceData
func hostnamesToMap(Hostnames []papi.Hostname) []map[string]interface{} {

	var res []map[string]interface{}
	for _, hn := range Hostnames {
		m := map[string]interface{}{}
		m["cname_from"] = hn.CnameFrom
		m["cname_to"] = hn.CnameTo
		m["cert_provisioning_type"] = hn.CertProvisioningType
		res = append(res, m)
	}
	return res
}

// Convert the given map from a schema.ResourceData to a slice of papi.Hostnames
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

func papiErrorsToList(Errors []*papi.Error) []interface{} {
	if len(Errors) == 0 {
		return nil
	}

	var RuleErrors []interface{}

	for _, err := range Errors {
		if err == nil {
			continue
		}

		RuleErrors = append(RuleErrors, papiErrorToMap(err))
	}

	return RuleErrors
}

func papiErrorToMap(err *papi.Error) map[string]interface{} {
	if err == nil {
		return nil
	}

	return map[string]interface{}{
		"type":           err.Type,
		"title":          err.Title,
		"detail":         err.Detail,
		"instance":       err.Instance,
		"behavior_name":  err.BehaviorName,
		"error_location": err.ErrorLocation,
		"status_code":    err.StatusCode,
	}
}
