package property

import (
	"context"
	"errors"
	"fmt"
	"strconv"
	"strings"

	"github.com/apex/log"
	"github.com/hashicorp/terraform-plugin-sdk/v2/diag"
	//log "github.com/sirupsen/logrus"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/schema"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/validation"
	"github.com/tidwall/gjson"

	"github.com/akamai/AkamaiOPEN-edgegrid-golang/jsonhooks-v1"
	"github.com/akamai/AkamaiOPEN-edgegrid-golang/papi-v1"
	"github.com/akamai/terraform-provider-akamai/v2/pkg/akamai"
	"github.com/akamai/terraform-provider-akamai/v2/pkg/tools"
)

func resourceProperty() *schema.Resource {
	return &schema.Resource{
		CreateContext: resourcePropertyCreate,
		ReadContext:   resourcePropertyRead,
		UpdateContext: resourcePropertyUpdate,
		DeleteContext: resourcePropertyDelete,
		CustomizeDiff: resourceCustomDiffCustomizeDiff,
		Importer: &schema.ResourceImporter{
			StateContext: resourcePropertyImport,
		},
		Schema: akamaiPropertySchema,
	}
}

var akamaiPropertySchema = map[string]*schema.Schema{
	"account": {
		Type:     schema.TypeString,
		Computed: true,
	},
	"contract": {
		Type:     schema.TypeString,
		Optional: true,
		ForceNew: true,
	},
	"group": {
		Type:     schema.TypeString,
		Optional: true,
		ForceNew: true,
	},
	"product": {
		Type:     schema.TypeString,
		Optional: true,
		ForceNew: true,
	},
	"rule_format": {
		Type:     schema.TypeString,
		Optional: true,
	},
	// Will get added to the default rule
	"cp_code": {
		Type:     schema.TypeString,
		Optional: true,
	},
	"name": {
		Type:     schema.TypeString,
		Required: true,
		ForceNew: true,
	},
	"version": {
		Type:     schema.TypeInt,
		Computed: true,
	},
	"staging_version": {
		Type:     schema.TypeInt,
		Computed: true,
	},
	"production_version": {
		Type:     schema.TypeInt,
		Computed: true,
	},
	"contact": {
		Type:     schema.TypeSet,
		Required: true,
		Elem:     &schema.Schema{Type: schema.TypeString},
	},
	"edge_hostnames": {
		Type:     schema.TypeMap,
		Computed: true,
		Elem:     &schema.Schema{Type: schema.TypeString},
	},
	"hostnames": {
		Type:     schema.TypeMap,
		Required: true,
		Elem:     &schema.Schema{Type: schema.TypeString},
	},

	// Will get added to the default rule
	"origin": {
		Type:     schema.TypeSet,
		Optional: true,
		Elem: &schema.Resource{
			Schema: map[string]*schema.Schema{
				"hostname": {
					Type:     schema.TypeString,
					Required: true,
				},
				"port": {
					Type:     schema.TypeInt,
					Optional: true,
					Default:  80,
				},
				"forward_hostname": {
					Type:     schema.TypeString,
					Optional: true,
					Default:  "ORIGIN_HOSTNAME",
				},
				"cache_key_hostname": {
					Type:     schema.TypeString,
					Optional: true,
					Default:  "ORIGIN_HOSTNAME",
				},
				"compress": {
					Type:     schema.TypeBool,
					Optional: true,
					Default:  false,
				},
				"enable_true_client_ip": {
					Type:     schema.TypeBool,
					Optional: true,
					Default:  false,
				},
			},
		},
	},
	"is_secure": {
		Type:     schema.TypeBool,
		Optional: true,
	},
	"rules": {
		Type:             schema.TypeString,
		Optional:         true,
		ValidateFunc:     validation.StringIsJSON,
		DiffSuppressFunc: suppressEquivalentJSONDiffs,
	},
	"variables": {
		Type:     schema.TypeString,
		Optional: true,
	},
	"rulessha": {
		Type:     schema.TypeString,
		Computed: true,
	},
}

func resourcePropertyCreate(_ context.Context, d *schema.ResourceData, m interface{}) diag.Diagnostics {
	meta := akamai.Meta(m)
	logger := meta.Log("PAPI", "resourcePropertyCreate")
	CorrelationID := "[PAPI][resourcePropertyCreate-" + meta.OperationID() + "]"
	group, err := getGroup(d, CorrelationID, logger)
	if err != nil {
		return diag.FromErr(err)
	}
	contract, err := getContract(d, CorrelationID, logger)
	if err != nil {
		return diag.FromErr(err)
	}
	product, err := getProduct(d, contract, CorrelationID, logger)
	if err != nil {
		return diag.FromErr(err)
	}

	var property *papi.Property
	name, err := tools.GetStringValue("name", d)
	if err != nil {
		return diag.FromErr(err)
	}
	if property = findProperty(name, CorrelationID); property == nil {
		property, err = createProperty(contract, group, product, d, CorrelationID, logger)
		if err != nil {
			return diag.FromErr(err)
		}
	}

	err = ensureEditableVersion(property, CorrelationID)
	if err != nil {
		return diag.FromErr(err)
	}
	if err := d.Set("account", property.AccountID); err != nil {
		return diag.FromErr(fmt.Errorf("%w: %s", tools.ErrValueSet, err.Error()))
	}
	if err := d.Set("version", property.LatestVersion); err != nil {
		return diag.FromErr(fmt.Errorf("%w: %s", tools.ErrValueSet, err.Error()))
	}

	d.SetId(property.PropertyID)

	rules, err := getRules(d, property, contract, group, CorrelationID, logger)
	if err != nil {
		return diag.FromErr(err)
	}
	if err = rules.Save(CorrelationID); err != nil {
		if err == papi.ErrorMap[papi.ErrInvalidRules] && len(rules.Errors) > 0 {
			var msg string
			var diags diag.Diagnostics
			for _, v := range rules.Errors {
				msg += fmt.Sprintf("\n Rule validation error: %s %s %s %s %s", v.Type, v.Title, v.Detail, v.Instance, v.BehaviorName)
				diags = append(diags, diag.Diagnostic{
					Severity: diag.Error,
					Summary:  "Invalid Property Rules",
					Detail:   msg,
				})
			}
			return diags
		}
		return diag.FromErr(err)
	}

	hostnames, err := setHostnames(property, d, CorrelationID, logger)
	if err != nil {
		return diag.FromErr(fmt.Errorf("%s", err.Error()))
	}

	if err := d.Set("edge_hostnames", hostnames); err != nil {
		return diag.FromErr(fmt.Errorf("%w: %s", tools.ErrValueSet, err.Error()))
	}
	rulesAPI, err := property.GetRules(CorrelationID)
	if err != nil {
		// TODO not sure what to do with this error (is it possible to return here)
		logger.Warnf("calling 'GetRules': %s", err.Error())
	}
	rulesAPI.Etag = ""
	body, err := jsonhooks.Marshal(rulesAPI)
	if err != nil {
		return diag.FromErr(err)
	}

	sha1hashAPI := tools.GetSHAString(string(body))
	logger.Debugf("CREATE SHA from Json %s\n", sha1hashAPI)
	logger.Debugf("CREATE Check rules after unmarshal from Json %s\n", string(body))

	if err := d.Set("rulessha", sha1hashAPI); err != nil {
		return diag.FromErr(fmt.Errorf("%w: %s", tools.ErrValueSet, err.Error()))
	}
	d.SetId(fmt.Sprintf("%s", property.PropertyID))
	if err := d.Set("rules", string(body)); err != nil {
		return diag.FromErr(fmt.Errorf("%w: %s", tools.ErrValueSet, err.Error()))
	}
	logger.Debugf("Done")
	return resourcePropertyRead(nil, d, nil)
}

func getRules(d *schema.ResourceData, property *papi.Property, contract *papi.Contract, group *papi.Group, correlationid string, logger log.Interface) (*papi.Rules, error) {
	rules := papi.NewRules()
	rules.Rule.Name = "default"
	rules.PropertyID = d.Id()
	rules.PropertyVersion = property.LatestVersion

	origin, err := createOrigin(d, correlationid, logger)
	if err != nil {
		return nil, err
	}

	_, ok := d.GetOk("rules")
	if ok {
		logger.Debugf("Unmarshal Rules from JSON")
		unmarshalRulesFromJSON(d, rules)
	}
	ruleFormat, err := tools.GetStringValue("rule_format", d)
	if err != nil {
		if !errors.Is(err, tools.ErrNotFound) {
			return nil, err
		}
		ruleFormats := papi.NewRuleFormats()
		rules.RuleFormat, err = ruleFormats.GetLatest(correlationid)
		if err != nil {
			return nil, err
		}
	} else {
		rules.RuleFormat = ruleFormat
	}

	cpCode, err := getCPCode(d, contract, group, correlationid, logger)
	if err != nil {
		return nil, err
	}

	logger.Debugf("updateStandardBehaviors")
	updateStandardBehaviors(rules, cpCode, origin, correlationid, logger)
	logger.Debugf("fixupPerformanceBehaviors")
	fixupPerformanceBehaviors(rules, correlationid, logger)

	return rules, nil
}

func setHostnames(property *papi.Property, d *schema.ResourceData, _ string, logger log.Interface) (map[string]string, error) {
	hostnameEdgeHostnames, ok := d.Get("hostnames").(map[string]interface{})
	if !ok {
		return nil, fmt.Errorf("%w: %s, %q", tools.ErrInvalidType, "hostnames", "map[string]interface{}")
	}

	edgeHostnames, err := papi.GetEdgeHostnames(property.Contract, property.Group, "")
	if err != nil {
		return nil, err
	}
	hostnames := papi.NewHostnames()
	hostnames.PropertyVersion = property.LatestVersion
	hostnames.PropertyID = property.PropertyID

	edgeHostnamesMap := make(map[string]string, len(hostnameEdgeHostnames))
	for public, edgeHostname := range hostnameEdgeHostnames {
		newEdgeHostname := edgeHostnames.NewEdgeHostname()
		edgeHostNameStr, ok := edgeHostname.(string)
		if !ok {
			return nil, fmt.Errorf("%w: %s, %q", tools.ErrInvalidType, "edge_hostname", "string")
		}
		newEdgeHostname.EdgeHostnameDomain = edgeHostNameStr
		logger.Debugf("Searching for edge hostname: %s, for hostname: %s", edgeHostNameStr, public)
		newEdgeHostname, err = edgeHostnames.FindEdgeHostname(newEdgeHostname)
		if err != nil {
			return nil, fmt.Errorf("edge hostname not found: %s", edgeHostNameStr)
		}
		logger.Debugf("Found edge hostname: %s", newEdgeHostname.EdgeHostnameDomain)

		hostname := hostnames.NewHostname()
		hostname.EdgeHostnameID = newEdgeHostname.EdgeHostnameID
		hostname.CnameFrom = public
		hostname.CnameTo = newEdgeHostname.EdgeHostnameDomain
		edgeHostnamesMap[public] = newEdgeHostname.EdgeHostnameDomain
	}

	if err = hostnames.Save(); err != nil {
		return nil, err
	}
	return edgeHostnamesMap, nil
}

func createProperty(contract *papi.Contract, group *papi.Group, product *papi.Product, d *schema.ResourceData, correlationid string, logger log.Interface) (*papi.Property, error) {
	logger.Debugf("Creating property")
	property, err := group.NewProperty(contract)
	if err != nil {
		return nil, err
	}

	property.ProductID = product.ProductID
	propertyName, err := tools.GetStringValue("name", d)
	if err != nil {
		return nil, err
	}
	property.PropertyName = propertyName

	ruleFormat, err := tools.GetStringValue("rule_format", d)
	if err != nil {
		if !errors.Is(err, tools.ErrNotFound) {
			return nil, err
		}
		ruleFormats := papi.NewRuleFormats()
		property.RuleFormat, err = ruleFormats.GetLatest(correlationid)
		if err != nil {
			return nil, err
		}
	} else {
		property.RuleFormat = ruleFormat
	}

	if err = property.Save(correlationid); err != nil {
		return nil, err
	}
	logger.Debugf("Property created: %s", property.PropertyID)
	return property, nil
}

func resourcePropertyDelete(_ context.Context, d *schema.ResourceData, m interface{}) diag.Diagnostics {
	meta := akamai.Meta(m)
	logger := meta.Log("PAPI", "resourcePropertyDelete")
	CorrelationID := "[PAPI][resourcePropertyDelete-" + meta.OperationID() + "]"
	logger.Debugf("DELETING")
	contractID, err := tools.GetStringValue("contract", d)
	//Todo clean up redundant checks and bubble up errors
	if err != nil {
		if !errors.Is(err, tools.ErrNotFound) {
			return diag.FromErr(err)
		}
		return diag.FromErr(errors.New("missing contract ID"))
	}
	groupID, err := tools.GetStringValue("group", d)
	if err != nil {
		if !errors.Is(err, tools.ErrNotFound) {
			return diag.FromErr(err)
		}
		return diag.FromErr(errors.New("missing group ID"))
	}

	property := papi.NewProperty(papi.NewProperties())
	property.PropertyID = d.Id()
	property.Contract = &papi.Contract{ContractID: contractID}
	property.Group = &papi.Group{GroupID: groupID}

	err = property.GetProperty(CorrelationID)
	if err != nil {
		return diag.FromErr(err)
	}
	if property.StagingVersion != 0 {
		return diag.FromErr(fmt.Errorf("property is still active on %s and cannot be deleted", papi.NetworkStaging))
	}
	if property.ProductionVersion != 0 {
		return diag.FromErr(fmt.Errorf("property is still active on %s and cannot be deleted", papi.NetworkProduction))
	}
	if err = property.Delete(CorrelationID); err != nil {
		return diag.FromErr(err)
	}
	d.SetId("")
	logger.Debugf("Done")
	return nil
}

func resourcePropertyImport(_ context.Context, d *schema.ResourceData, m interface{}) ([]*schema.ResourceData, error) {
	meta := akamai.Meta(m)
	logger := meta.Log("PAPI", "resourcePropertyImport")
	propertyID := d.Id()

	if !strings.HasPrefix(propertyID, "prp_") {
		keys := []papi.SearchKey{papi.SearchByPropertyName, papi.SearchByHostname, papi.SearchByEdgeHostname}
		for _, searchKey := range keys {
			results, err := papi.Search(searchKey, propertyID, "") //<--correlationid
			if err != nil {
				// TODO determine why is this error ignored
				logger.Debugf("searching by key: %s: %w", searchKey, err)
				continue
			}

			if results != nil && len(results.Versions.Items) > 0 {
				propertyID = results.Versions.Items[0].PropertyID
				break
			}
		}
	}

	property := papi.NewProperty(papi.NewProperties())
	property.PropertyID = propertyID
	err := property.GetProperty("")
	if err != nil {
		return nil, err
	}

	if err := d.Set("account", property.AccountID); err != nil {
		return nil, fmt.Errorf("%w: %s", tools.ErrValueSet, err.Error())
	}
	if err := d.Set("contract", property.ContractID); err != nil {
		return nil, fmt.Errorf("%w: %s", tools.ErrValueSet, err.Error())
	}
	if err := d.Set("group", property.GroupID); err != nil {
		return nil, fmt.Errorf("%w: %s", tools.ErrValueSet, err.Error())
	}
	if err := d.Set("name", property.PropertyName); err != nil {
		return nil, fmt.Errorf("%w: %s", tools.ErrValueSet, err.Error())
	}
	if err := d.Set("version", property.LatestVersion); err != nil {
		return nil, fmt.Errorf("%w: %s", tools.ErrValueSet, err.Error())
	}
	d.SetId(property.PropertyID)
	return []*schema.ResourceData{d}, nil
}

func resourcePropertyRead(_ context.Context, d *schema.ResourceData, m interface{}) diag.Diagnostics {
	meta := akamai.Meta(m)
	logger := meta.Log("PAPI", "resourcePropertyRead")
	CorrelationID := "[PAPI][resourcePropertyRead-" + meta.OperationID() + "]"
	property := papi.NewProperty(papi.NewProperties())
	property.PropertyID = d.Id()
	err := property.GetProperty(CorrelationID)
	if err != nil {
		return diag.FromErr(err)
	}
	if err := d.Set("account", property.AccountID); err != nil {
		return diag.FromErr(fmt.Errorf("%w: %s", tools.ErrValueSet, err.Error()))
	}
	if err := d.Set("contract", property.ContractID); err != nil {
		return diag.FromErr(fmt.Errorf("%w: %s", tools.ErrValueSet, err.Error()))
	}
	if err := d.Set("group", property.GroupID); err != nil {
		return diag.FromErr(fmt.Errorf("%w: %s", tools.ErrValueSet, err.Error()))
	}
	if err := d.Set("name", property.PropertyName); err != nil {
		return diag.FromErr(fmt.Errorf("%w: %s", tools.ErrValueSet, err.Error()))
	}
	if err := d.Set("note", property.Note); err != nil {
		// since note is an optional parameter, just logging the error.
		logger.Warnf("%w: %s", tools.ErrValueSet, err.Error())
	}

	rules, err := property.GetRules(CorrelationID)
	if err != nil {
		// TODO not sure what to do with this error (is it possible to return here)
		logger.Warnf("calling 'GetRules': %s", err.Error())
	}
	rules.Etag = ""
	body, err := jsonhooks.Marshal(rules)
	if err != nil {
		return diag.FromErr(err)
	}
	sha1hashAPI := tools.GetSHAString(string(body))
	logger.Debugf("READ SHA from Json %s", sha1hashAPI)
	logger.Debugf("READ Rules from API : %s", string(body))
	if err := d.Set("rules", string(body)); err != nil {
		return diag.FromErr(fmt.Errorf("%w: %s", tools.ErrValueSet, err.Error()))
	}
	if err := d.Set("rulessha", sha1hashAPI); err != nil {
		return diag.FromErr(fmt.Errorf("%w: %s", tools.ErrValueSet, err.Error()))
	}

	if rules.RuleFormat != "" {
		if err := d.Set("rule_format", rules.RuleFormat); err != nil {
			return diag.FromErr(fmt.Errorf("%w: %s", tools.ErrValueSet, err.Error()))
		}
	} else {
		if err := d.Set("rule_format", property.RuleFormat); err != nil {
			return diag.FromErr(fmt.Errorf("%w: %s", tools.ErrValueSet, err.Error()))
		}
	}
	logger.Debugf("Property RuleFormat from API : %s", property.RuleFormat)
	if err := d.Set("version", property.LatestVersion); err != nil {
		return diag.FromErr(fmt.Errorf("%w: %s", tools.ErrValueSet, err.Error()))
	}
	if property.StagingVersion > 0 {
		if err := d.Set("staging_version", property.StagingVersion); err != nil {
			return diag.FromErr(fmt.Errorf("%w: %s", tools.ErrValueSet, err.Error()))
		}
	}
	if property.ProductionVersion > 0 {
		if err := d.Set("production_version", property.ProductionVersion); err != nil {
			return diag.FromErr(fmt.Errorf("%w: %s", tools.ErrValueSet, err.Error()))
		}
	}
	return nil
}

func resourcePropertyUpdate(_ context.Context, d *schema.ResourceData, m interface{}) diag.Diagnostics {
	meta := akamai.Meta(m)
	logger := meta.Log("PAPI", "resourcePropertyUpdate")
	CorrelationID := "[PAPI][resourcePropertyUpdate-" + meta.OperationID() + "]"
	logger.Debugf("UPDATING")
	property, err := getProperty(d, CorrelationID, logger)
	if err != nil {
		return diag.FromErr(err)
	}
	err = ensureEditableVersion(property, CorrelationID)
	if err != nil {
		return diag.FromErr(err)
	}

	rules, err := getRules(d, property, property.Contract, property.Group, CorrelationID, logger)
	if err != nil {
		// TODO not sure what to do with this error (is it possible to return here)
		logger.Warnf("calling 'getRules': %s", err.Error())
	}
	if d.HasChange("rule_format") || d.HasChange("rules") {
		ruleFormat, err := tools.GetStringValue("rule_format", d)
		if err != nil {
			if !errors.Is(err, tools.ErrNotFound) {
				return diag.FromErr(err)
			}
		} else {
			property.RuleFormat = ruleFormat
			rules.RuleFormat = ruleFormat
		}
		body, err := jsonhooks.Marshal(rules)
		if err != nil {
			return diag.FromErr(err)
		}
		if err := d.Set("rules", string(body)); err != nil {
			return diag.FromErr(fmt.Errorf("%w: %s", tools.ErrValueSet, err.Error()))
		}
		logger.Debugf("UPDATE Check rules after unmarshal from Json %s", string(body))
		if err = rules.Save(CorrelationID); err != nil {
			if err == papi.ErrorMap[papi.ErrInvalidRules] && len(rules.Errors) > 0 {
				var msg string
				//Todo Create reusable diagnostic functions if needed
				var diags diag.Diagnostics
				for _, v := range rules.Errors {
					msg += fmt.Sprintf("\n Rule validation error: %s %s %s %s %s", v.Type, v.Title, v.Detail, v.Instance, v.BehaviorName)
					diags = append(diags, diag.Diagnostic{
						Severity: diag.Error,
						Summary:  "Invalid Property Rules",
						Detail:   msg,
					})
				}
				return diags
			}
			logger.Debugf("update rules.Save err: %#v", err)
			return diag.FromErr(fmt.Errorf("update rules.Save err: %#v", err))
		}

		rules, err = property.GetRules(CorrelationID)
		if err != nil {
			// TODO not sure what to do with this error (is it possible to return here)
			logger.Warnf("calling 'GetRules': %s", err.Error())
		}
		rules.Etag = ""
		body, err = jsonhooks.Marshal(rules)
		if err != nil {
			return diag.FromErr(err)
		}

		sha1hashAPI := tools.GetSHAString(string(body))
		logger.Debugf("UPDATE SHA from Json %s", sha1hashAPI)
		if err := d.Set("rulessha", sha1hashAPI); err != nil {
			return diag.FromErr(fmt.Errorf("%w: %s", tools.ErrValueSet, err.Error()))
		}
	}
	if err := d.Set("version", property.LatestVersion); err != nil {
		return diag.FromErr(fmt.Errorf("%w: %s", tools.ErrValueSet, err.Error()))
	}

	if d.HasChange("hostnames") {
		edgeHostnamesMap, err := setHostnames(property, d, CorrelationID, logger)
		if err != nil {
			return diag.FromErr(fmt.Errorf("setHostnames err: %#v", err))
		}
		if err := d.Set("edge_hostnames", edgeHostnamesMap); err != nil {
			return diag.FromErr(fmt.Errorf("%w: %s", tools.ErrValueSet, err.Error()))
		}
	}

	logger.Debugf("Done")
	return resourcePropertyRead(nil, d, nil)
}

func resourceCustomDiffCustomizeDiff(ctx context.Context, d *schema.ResourceDiff, m interface{}) error {
	meta := akamai.Meta(m)
	logger := meta.Log("PAPI", "resourceCustomDiffCustomizeDiff")

	logger.Debugf("ID: %s", d.Id())
	// Note that this gets put into state after the update, regardless of whether
	// or not anything is acted upon in the diff.
	old, new := d.GetChange("rules")
	oldStr, ok := old.(string)
	if !ok {
		return fmt.Errorf("value is of invalid type: %v; should be %s", old, "string")
	}
	newStr, ok := old.(string)
	if !ok {
		return fmt.Errorf("value is of invalid type: %v; should be %s", new, "string")
	}
	logger.Debugf("OLD: %s", oldStr)
	logger.Debugf("NEW: %s", newStr)
	if !suppressEquivalentJSONPendingDiffs(oldStr, newStr, d) {
		logger.Debugf("CHANGED VALUES: %s %s " + oldStr + " " + newStr)
		if err := d.SetNewComputed("version"); err != nil {
			return fmt.Errorf("%w: %s", tools.ErrValueSet, err.Error())
		}
	}
	return nil
}

// Helpers
func getProperty(d interface{}, correlationid string, logger log.Interface) (*papi.Property, error) {
	logger.Debugf("Fetching property")
	var propertyID string
	switch d.(type) {
	case *schema.ResourceData:
		propertyID = d.(*schema.ResourceData).Id()
	case *schema.ResourceDiff:
		propertyID = d.(*schema.ResourceDiff).Id()
	default:
		return nil, fmt.Errorf("resource is of invalid type; should be '*schema.ResourceDiff' or '*schema.ResourceData'")
	}
	property := papi.NewProperty(papi.NewProperties())
	property.PropertyID = propertyID
	err := property.GetProperty(correlationid)
	return property, err
}

func getGroup(d *schema.ResourceData, correlationid string, logger log.Interface) (*papi.Group, error) {
	logger.Debugf("Fetching groups")
	groupID, err := tools.GetStringValue("group", d)
	if err != nil {
		if !errors.Is(err, tools.ErrNotFound) {
			return nil, err
		}
		return nil, ErrNoGroupProvided
	}
	groups := papi.NewGroups()
	err = groups.GetGroups(correlationid)
	if err != nil {
		return nil, err
	}
	groupID, err = tools.AddPrefix(groupID, "grp_")
	if err != nil {
		return nil, err
	}
	group, err := groups.FindGroup(groupID)
	if err != nil {
		return nil, err
	}

	logger.Debugf("Group found: %s", group.GroupID)
	return group, nil
}

func getContract(d *schema.ResourceData, correlationid string, logger log.Interface) (*papi.Contract, error) {
	logger.Debugf("Fetching contract")
	contractID, err := tools.GetStringValue("contract", d)
	if err != nil {
		if !errors.Is(err, tools.ErrNotFound) {
			return nil, err
		}
		return nil, ErrNoContractProvided
	}
	contracts := papi.NewContracts()
	err = contracts.GetContracts(correlationid)
	if err != nil {
		return nil, err
	}
	contractID, err = tools.AddPrefix(contractID, "ctr_")
	if err != nil {
		return nil, err
	}
	contract, err := contracts.FindContract(contractID)
	if err != nil {
		return nil, err
	}

	logger.Debugf("Contract found: %s", contract.ContractID)
	return contract, nil
}

func getCPCode(d interface{}, contract *papi.Contract, group *papi.Group, _ string, logger log.Interface) (*papi.CpCode, error) {
	if contract == nil {
		return nil, ErrNoContractProvided
	}
	if group == nil {
		return nil, ErrNoGroupProvided
	}
	var cpCodeID string
	var err error
	switch d.(type) {
	case *schema.ResourceData:
		cpCodeID, err = tools.GetStringValue("cp_code", d.(*schema.ResourceData))
	case *schema.ResourceDiff:
		cpCodeID, err = tools.GetStringValue("cp_code", d.(*schema.ResourceDiff))
	default:
		return nil, fmt.Errorf("resource is of invalid type; should be '*schema.ResourceDiff' or '*schema.ResourceData'")
	}
	if err != nil && !errors.Is(err, tools.ErrNotFound) {
		return nil, err
	}
	logger.Debugf("Fetching CP code")
	cpCode := papi.NewCpCodes(contract, group).NewCpCode()
	cpCode.CpcodeID = cpCodeID
	if err := cpCode.GetCpCode(); err != nil {
		return nil, err
	}
	logger.Debugf("CP code found: %s", cpCode.CpcodeID)
	return cpCode, nil
}

func getProduct(d *schema.ResourceData, contract *papi.Contract, correlationid string, logger log.Interface) (*papi.Product, error) {
	if contract == nil {
		return nil, ErrNoContractProvided
	}
	logger.Debugf("Fetching product")
	productID, err := tools.GetStringValue("product", d)
	if err != nil {
		if !errors.Is(err, tools.ErrNotFound) {
			return nil, err
		}
		return nil, ErrNoProductProvided
	}
	products := papi.NewProducts()
	err = products.GetProducts(contract, correlationid)
	if err != nil {
		return nil, err
	}
	productID, err = tools.AddPrefix(productID, "prd_")
	if err != nil {
		return nil, err
	}
	product, err := products.FindProduct(productID)
	if err != nil {
		return nil, err
	}

	logger.Debugf("Product found: %s", product.ProductID)
	return product, nil
}

func createOrigin(d interface{}, _ string, logger log.Interface) (*papi.OptionValue, error) {
	logger.Debugf("Setting origin")
	var origin *schema.Set
	var err error

	switch d.(type) {
	case *schema.ResourceData:
		origin, err = tools.GetSetValue("origin", d.(*schema.ResourceData))
	case *schema.ResourceDiff:
		origin, err = tools.GetSetValue("origin", d.(*schema.ResourceDiff))
	default:
		return nil, fmt.Errorf("resource is of invalid type; should be '*schema.ResourceDiff' or '*schema.ResourceData'")
	}

	if err != nil && !errors.Is(err, tools.ErrNotFound) {
		return nil, err
	}
	if origin.Len() == 0 {
		return nil, fmt.Errorf("'origin' property must have at least one value")
	}
	originConfig, ok := origin.List()[0].(map[string]interface{})
	if !ok {
		return nil, fmt.Errorf("origin set value is of invalid type; should be 'map[string]interface{}'")
	}

	originValues := make(map[string]interface{})

	originValues["originType"] = "CUSTOMER"
	if val, ok := originConfig["hostname"]; ok {
		originValues["hostname"], ok = val.(string)
		if !ok {
			return nil, fmt.Errorf("%w: %s, %q", tools.ErrInvalidType, "hostname", "string")
		}
	}

	if val, ok := originConfig["port"]; ok {
		originValues["httpPort"], ok = val.(int)
		if !ok {
			return nil, fmt.Errorf("%w: %s, %q", tools.ErrInvalidType, "httpPort", "int")
		}
	}

	if val, ok := originConfig["cache_key_hostname"]; ok {
		originValues["cacheKeyHostname"], ok = val.(string)
		if !ok {
			return nil, fmt.Errorf("%w: %s, %q", tools.ErrInvalidType, "cacheKeyHostname", "string")
		}
	}

	if val, ok := originConfig["compress"]; ok {
		originValues["compress"], ok = val.(bool)
		if !ok {
			return nil, fmt.Errorf("%w: %s, %q", tools.ErrInvalidType, "compress", "bool")
		}
	}

	if val, ok := originConfig["enable_true_client_ip"]; ok {
		originValues["enableTrueClientIp"], ok = val.(bool)
		if !ok {
			return nil, fmt.Errorf("%w: %s, %q", tools.ErrInvalidType, "enableTrueClientIp", "bool")
		}
	}

	forwardHostname, ok := originConfig["forward_hostname"].(string)
	if ok {
		if forwardHostname == "ORIGIN_HOSTNAME" || forwardHostname == "REQUEST_HOST_HEADER" {
			logger.Debugf("Setting non-custom forward hostname")

			originValues["forwardHostHeader"] = forwardHostname
		} else {
			logger.Debugf("Setting custom forward hostname")

			originValues["forwardHostHeader"] = "CUSTOM"
			originValues["customForwardHostHeader"] = forwardHostname
		}

	}
	ov := papi.OptionValue(originValues)
	return &ov, nil
}

func fixupPerformanceBehaviors(rules *papi.Rules, _ string, logger log.Interface) {
	behavior, err := rules.FindBehavior("/Performance/sureRoute")
	if err != nil || behavior == nil || behavior.Options["testObjectUrl"] != "" {
		return
	}

	logger.Debugf("Fixing Up SureRoute Behavior")
	behavior.MergeOptions(papi.OptionValue{
		"testObjectUrl":   "/akamai/sureroute-testobject.html",
		"enableCustomKey": false,
		"enabled":         false,
	})
}

func updateStandardBehaviors(rules *papi.Rules, cpCode *papi.CpCode, origin *papi.OptionValue, _ string, logger log.Interface) {
	logger.Debugf("cpCode: %#v", cpCode)
	if cpCode != nil {
		b := papi.NewBehavior()
		b.Name = "cpCode"
		b.Options = papi.OptionValue{
			"value": papi.OptionValue{
				"id": cpCode.ID(),
			},
		}
		rules.Rule.MergeBehavior(b)
	}

	if origin != nil {
		b := papi.NewBehavior()
		b.Name = "origin"
		b.Options = *origin
		rules.Rule.MergeBehavior(b)
	}
}

// TODO: discuss how property rules should be handled
func unmarshalRulesFromJSON(d *schema.ResourceData, propertyRules *papi.Rules) {
	// Default Rules
	rules, ok := d.GetOk("rules")
	if ok {
		propertyRules.Rule = &papi.Rule{Name: "default"}
		//		log.Println("[DEBUG] RulesJson")

		rulesJSON := gjson.Get(rules.(string), "rules")
		rulesJSON.ForEach(func(key, value gjson.Result) bool {
			//			log.Println("[DEBUG] unmarshalRulesFromJson KEY RULES KEY = " + key.String() + " VAL " + value.String())

			if key.String() == "behaviors" {
				behavior := gjson.Parse(value.String())
				//				log.Println("[DEBUG] unmarshalRulesFromJson KEY BEHAVIOR " + behavior.String())
				if gjson.Get(behavior.String(), "#.name").Exists() {

					behavior.ForEach(func(key, value gjson.Result) bool {
						//						log.Println("[DEBUG] unmarshalRulesFromJson BEHAVIOR LOOP KEY =" + key.String() + " VAL " + value.String())

						bb, ok := value.Value().(map[string]interface{})
						if ok {
							//							log.Println("[DEBUG] unmarshalRulesFromJson BEHAVIOR MAP  ", bb)
							for k, v := range bb {
								log.Debugf("k:", k, "v:", v)
							}

							beh := papi.NewBehavior()

							beh.Name = bb["name"].(string)
							boptions, ok := bb["options"]
							//							log.Println("[DEBUG] unmarshalRulesFromJson KEY BEHAVIOR BOPTIONS ", boptions)
							if ok {
								beh.Options = boptions.(map[string]interface{})
								//								log.Println("[DEBUG] unmarshalRulesFromJson KEY BEHAVIOR EXTRACT BOPTIONS ", beh.Options)
							}

							propertyRules.Rule.MergeBehavior(beh)
						}

						return true // keep iterating
					}) // behavior list loop

				}

				if key.String() == "criteria" {
					criteria := gjson.Parse(value.String())

					criteria.ForEach(func(key, value gjson.Result) bool {
						//						log.Println("[DEBUG] unmarshalRulesFromJson KEY CRITERIA " + key.String() + " VAL " + value.String())

						cc, ok := value.Value().(map[string]interface{})
						if ok {
							//							log.Println("[DEBUG] unmarshalRulesFromJson CRITERIA MAP  ", cc)
							newCriteria := papi.NewCriteria()
							newCriteria.Name = cc["name"].(string)

							coptions, ok := cc["option"]
							if ok {
								//								println("OPTIONS ", coptions)
								newCriteria.Options = coptions.(map[string]interface{})
							}
							propertyRules.Rule.MergeCriteria(newCriteria)
						}
						return true
					})
				} // if ok criteria
			} /// if ok behaviors

			if key.String() == "children" {
				childRules := gjson.Parse(value.String())
				//				println("CHILD RULES " + childRules.String())

				for _, rule := range extractRulesJSON(d, childRules) {
					propertyRules.Rule.MergeChildRule(rule)
				}
			}

			if key.String() == "variables" {

				//				log.Println("unmarshalRulesFromJson VARS from JSON ", value.String())
				variables := gjson.Parse(value.String())

				variables.ForEach(func(key, value gjson.Result) bool {
					//					log.Println("unmarshalRulesFromJson VARS from JSON LOOP ", value)
					variableMap, ok := value.Value().(map[string]interface{})
					//					log.Println("unmarshalRulesFromJson VARS from JSON LOOP NAME ", variableMap["name"].(string))
					//					log.Println("unmarshalRulesFromJson VARS from JSON LOOP DESC ", variableMap["description"].(string))
					if ok {
						newVariable := papi.NewVariable()
						newVariable.Name = variableMap["name"].(string)
						newVariable.Description = variableMap["description"].(string)
						newVariable.Value = variableMap["value"].(string)
						newVariable.Hidden = variableMap["hidden"].(bool)
						newVariable.Sensitive = variableMap["sensitive"].(bool)
						propertyRules.Rule.AddVariable(newVariable)
					}
					return true
				}) //variables

			}

			if key.String() == "options" {
				//				log.Println("unmarshalRulesFromJson OPTIONS from JSON", value.String())
				options := gjson.Parse(value.String())
				options.ForEach(func(key, value gjson.Result) bool {
					switch {
					case key.String() == "is_secure" && value.Bool():
						propertyRules.Rule.Options.IsSecure = value.Bool()
					}

					return true
				})
			}

			return true // keep iterating
		}) // for loop rules

		// ADD vars from variables resource
		jsonvars, ok := d.GetOk("variables")
		if ok {
			//			log.Println("unmarshalRulesFromJson VARS from JSON ", jsonvars)
			variables := gjson.Parse(jsonvars.(string))
			result := gjson.Get(variables.String(), "variables")
			//			log.Println("unmarshalRulesFromJson VARS from JSON VARIABLES ", result)

			result.ForEach(func(key, value gjson.Result) bool {
				//				log.Println("unmarshalRulesFromJson VARS from JSON LOOP ", value)
				variableMap, ok := value.Value().(map[string]interface{})
				//				log.Println("unmarshalRulesFromJson VARS from JSON LOOP NAME ", variableMap["name"].(string))
				//				log.Println("unmarshalRulesFromJson VARS from JSON LOOP DESC ", variableMap["description"].(string))
				if ok {
					newVariable := papi.NewVariable()
					newVariable.Name = variableMap["name"].(string)
					newVariable.Description = variableMap["description"].(string)
					newVariable.Value = variableMap["value"].(string)
					newVariable.Hidden = variableMap["hidden"].(bool)
					newVariable.Sensitive = variableMap["sensitive"].(bool)
					propertyRules.Rule.AddVariable(newVariable)
				}
				return true
			}) //variables
		}

		// ADD isSecure from resource
		isSecure, set := d.GetOkExists("is_secure")
		if set && isSecure.(bool) {
			propertyRules.Rule.Options.IsSecure = true
		} else if set && !isSecure.(bool) {
			propertyRules.Rule.Options.IsSecure = false
		}

		// ADD cpCode from resource
		cpCode, set := d.GetOk("cp_code")
		if set {
			beh := papi.NewBehavior()
			beh.Name = "cpCode"
			beh.Options = papi.OptionValue{
				"value": papi.OptionValue{
					"id": cpCode.(string),
				},
			}
			propertyRules.Rule.MergeBehavior(beh)
		}
	}
}

func extractOptions(options *schema.Set) (map[string]interface{}, error) {
	optv := make(map[string]interface{})
	for _, option := range options.List() {
		optionMap, ok := option.(map[string]interface{})
		if !ok {
			continue
		}
		if val, ok := optionMap["value"].(string); ok && val != "" {
			optv[optionMap["key"].(string)] = convertString(val)
			continue
		}
		vals, ok := optionMap["values"]
		if !ok {
			continue
		}
		valsSet, ok := vals.(*schema.Set)
		if !ok {
			return nil, fmt.Errorf("%w: %s, %q", tools.ErrInvalidType, "values", "*schema.Set")
		}
		if valsSet.Len() == 0 {
			optv[optionMap["key"].(string)] = convertString(optionMap["value"].(string))
			continue
		}
		if valsSet.Len() > 0 {
			op := make([]interface{}, 0)
			for _, v := range vals.(*schema.Set).List() {
				op = append(op, convertString(v.(string)))
			}

			optv[optionMap["key"].(string)] = op
		}
	}
	return optv, nil
}

func convertString(v string) interface{} {
	if f1, err := strconv.ParseFloat(v, 64); err == nil {
		return f1
	}
	// FIXME: execution will never reach this as every int representation will be captured by ParseFloat() above
	// this should either be moved above ParseFloat block or removed
	if f2, err := strconv.ParseInt(v, 10, 64); err == nil {
		return f2
	}
	if f3, err := strconv.ParseBool(v); err == nil {
		return f3
	}
	return v
}

// TODO: discuss how property rules should be handled
func extractRulesJSON(d interface{}, drules gjson.Result) []*papi.Rule {
	var rules []*papi.Rule
	drules.ForEach(func(key, value gjson.Result) bool {
		rule := papi.NewRule()
		vv, ok := value.Value().(map[string]interface{})
		if ok {
			rule.Name, _ = vv["name"].(string)
			rule.Comments, _ = vv["comments"].(string)
			criteriaMustSatisfy, ok := vv["criteriaMustSatisfy"]
			if ok {
				if criteriaMustSatisfy.(string) == "all" {
					rule.CriteriaMustSatisfy = papi.RuleCriteriaMustSatisfyAll
				}

				if criteriaMustSatisfy.(string) == "any" {
					rule.CriteriaMustSatisfy = papi.RuleCriteriaMustSatisfyAny
				}
			}
			log.Debugf("extractRulesJSON Set criteriaMustSatisfy RESULT RULE value set " + string(rule.CriteriaMustSatisfy) + " " + rule.Name + " " + rule.Comments)

			ruledetail := gjson.Parse(value.String())
			//			log.Println("[DEBUG] RULE DETAILS ", ruledetail)

			ruledetail.ForEach(func(key, value gjson.Result) bool {

				if key.String() == "behaviors" {
					//					log.Println("[DEBUG] BEHAVIORS KEY CHILD RULE ", key.String())

					behaviors := gjson.Parse(value.String())
					//					log.Println("[DEBUG] BEHAVIORS NAME ", behaviors)
					behaviors.ForEach(func(key, value gjson.Result) bool {
						//						log.Println("[DEBUG] BEHAVIORS KEY CHILD RULE LOOP KEY = " + key.String() + " VAL " + value.String())
						behaviorMap, ok := value.Value().(map[string]interface{})
						if ok {
							newBehavior := papi.NewBehavior()
							newBehavior.Name = behaviorMap["name"].(string)
							behaviorOptions, ok := behaviorMap["options"]
							if ok {
								newBehavior.Options = behaviorOptions.(map[string]interface{})
							}
							rule.MergeBehavior(newBehavior)
						}
						return true
					}) //behaviors
				}

				if key.String() == "criteria" {
					//					log.Println("[DEBUG] CRITERIA KEY CHILD RULE ", key.String())
					criterias := gjson.Parse(value.String())
					criterias.ForEach(func(key, value gjson.Result) bool {
						criteriaMap, ok := value.Value().(map[string]interface{})
						if ok {
							newCriteria := papi.NewCriteria()
							newCriteria.Name = criteriaMap["name"].(string)
							criteriaOptions, ok := criteriaMap["options"]
							if ok {
								newCriteria.Options = criteriaOptions.(map[string]interface{})
							}
							rule.MergeCriteria(newCriteria)
						}
						return true
					}) //criteria
				}

				if key.String() == "variables" {
					//					log.Println("[DEBUG] VARIABLES KEY CHILD RULE ", key.String())
					variables := gjson.Parse(value.String())
					variables.ForEach(func(key, value gjson.Result) bool {
						variableMap, ok := value.Value().(map[string]interface{})
						if ok {
							newVariable := papi.NewVariable()
							newVariable.Name = variableMap["name"].(string)
							newVariable.Description = variableMap["description"].(string)
							newVariable.Value = variableMap["value"].(string)
							newVariable.Hidden = variableMap["hidden"].(bool)
							newVariable.Sensitive = variableMap["sensitive"].(bool)
							rule.AddVariable(newVariable)
						}
						return true
					}) //variables
				}

				if key.String() == "children" {
					childRules := gjson.Parse(value.String())
					//					println("CHILD RULES " + childRules.String())
					for _, newRule := range extractRulesJSON(d, childRules) {
						rule.MergeChildRule(newRule)
					}
				} //len > 0

				return true
			}) //Loop Detail

		}
		rules = append(rules, rule)

		return true
	})

	return rules
}

func extractRules(drules *schema.Set) ([]*papi.Rule, error) {

	var rules []*papi.Rule
	for _, v := range drules.List() {
		rule := papi.NewRule()
		vv, ok := v.(map[string]interface{})
		if ok {
			rule.Name = vv["name"].(string)
			rule.Comments = vv["comment"].(string)

			criteriaMustSatisfy, ok := vv["criteriaMustSatisfy"]
			if ok {
				if criteriaMustSatisfy.(string) == "all" {
					rule.CriteriaMustSatisfy = papi.RuleCriteriaMustSatisfyAll
				}

				if criteriaMustSatisfy.(string) == "any" {
					rule.CriteriaMustSatisfy = papi.RuleCriteriaMustSatisfyAny
				}
			}
			behaviors, ok := vv["behavior"]
			if ok {
				for _, behavior := range behaviors.(*schema.Set).List() {
					behaviorMap, ok := behavior.(map[string]interface{})
					if ok {
						newBehavior := papi.NewBehavior()
						newBehavior.Name = behaviorMap["name"].(string)
						behaviorOptions, ok := behaviorMap["option"]
						if ok {
							opts, err := extractOptions(behaviorOptions.(*schema.Set))
							if err != nil {
								return nil, err
							}
							newBehavior.Options = opts
						}
						rule.MergeBehavior(newBehavior)
					}
				}
			}

			criterias, ok := vv["criteria"]
			if ok {
				for _, criteria := range criterias.(*schema.Set).List() {
					criteriaMap, ok := criteria.(map[string]interface{})
					if ok {
						newCriteria := papi.NewCriteria()
						newCriteria.Name = criteriaMap["name"].(string)
						criteriaOptions, ok := criteriaMap["option"]
						if ok {
							crit, err := extractOptions(criteriaOptions.(*schema.Set))
							if err != nil {
								return nil, err
							}
							newCriteria.Options = crit
						}
						rule.MergeCriteria(newCriteria)
					}
				}
			}

			variables, ok := vv["variable"]
			if ok {
				for _, variable := range variables.(*schema.Set).List() {
					variableMap, ok := variable.(map[string]interface{})
					if ok {
						newVariable := papi.NewVariable()
						newVariable.Name = variableMap["name"].(string)
						newVariable.Description = variableMap["description"].(string)
						newVariable.Value = variableMap["value"].(string)
						newVariable.Hidden = variableMap["hidden"].(bool)
						newVariable.Sensitive = variableMap["sensitive"].(bool)
						rule.AddVariable(newVariable)
					}
				}
			}

			childRules, ok := vv["rule"]
			if ok && childRules.(*schema.Set).Len() > 0 {
				rules, err := extractRules(childRules.(*schema.Set))
				if err != nil {
					return nil, err
				}
				for _, newRule := range rules {
					rule.MergeChildRule(newRule)
				}
			}
		}
		rules = append(rules, rule)
	}

	return rules, nil
}

func findProperty(name string, correlationid string) *papi.Property {
	results, err := papi.Search(papi.SearchByPropertyName, name, correlationid)
	if err != nil || results == nil || len(results.Versions.Items) == 0 {
		return nil
	}

	property := &papi.Property{
		PropertyID: results.Versions.Items[0].PropertyID,
		Group: &papi.Group{
			GroupID: results.Versions.Items[0].GroupID,
		},
		Contract: &papi.Contract{
			ContractID: results.Versions.Items[0].ContractID,
		},
	}

	err = property.GetProperty(correlationid)
	if err != nil {
		return nil
	}

	return property
}

func ensureEditableVersion(property *papi.Property, correlationid string) error {
	latestVersion, err := property.GetLatestVersion("", correlationid)
	if err != nil {
		return err
	}

	versions, err := property.GetVersions(correlationid)
	if err != nil {
		return err
	}

	if latestVersion.ProductionStatus != papi.StatusInactive || latestVersion.StagingStatus != papi.StatusInactive {
		// The latest version has been activated on either production or staging, so we need to create a new version to apply changes on
		newVersion := versions.NewVersion(latestVersion, false, correlationid)
		if err = newVersion.Save(correlationid); err != nil {
			return err
		}
	}

	return property.GetProperty(correlationid)
}
