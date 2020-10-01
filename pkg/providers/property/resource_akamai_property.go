package property

import (
	"context"
	"encoding/json"
	"errors"
	"fmt"
	"strconv"
	"strings"

	"github.com/apex/log"
	"github.com/hashicorp/terraform-plugin-sdk/v2/diag"

	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/schema"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/validation"
	"github.com/tidwall/gjson"

	"github.com/akamai/AkamaiOPEN-edgegrid-golang/jsonhooks-v1"
	"github.com/akamai/AkamaiOPEN-edgegrid-golang/papi-v1"
	v2 "github.com/akamai/AkamaiOPEN-edgegrid-golang/v2/pkg/papi"
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

func resourcePropertyCreate(ctx context.Context, d *schema.ResourceData, m interface{}) diag.Diagnostics {
	meta := akamai.Meta(m)
	client := inst.Client(meta)
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

	name, err := tools.GetStringValue("name", d)
	if err != nil {
		return diag.FromErr(err)
	}
	prop, err := findProperty(ctx, name, meta)
	if err != nil {
		if !errors.Is(err, ErrPropertyNotFound) {
			return diag.FromErr(err)
		}
		prop, err = createProperty(ctx, contract.ContractID, group.GroupID, product.ProductID, d, meta)
		if err != nil {
			return diag.FromErr(err)
		}
	}

	err = ensureEditableVersion(ctx, prop, meta)
	if err != nil {
		return diag.FromErr(err)
	}
	if err := d.Set("account", prop.AccountID); err != nil {
		return diag.FromErr(fmt.Errorf("%w: %s", tools.ErrValueSet, err.Error()))
	}
	if err := d.Set("version", prop.LatestVersion); err != nil {
		return diag.FromErr(fmt.Errorf("%w: %s", tools.ErrValueSet, err.Error()))
	}

	d.SetId(prop.PropertyID)

	rules, err := getRules(ctx, d, prop, contract.ContractID, group.GroupID, meta)
	if err != nil {
		return diag.FromErr(err)
	}
	if _, err := client.UpdateRuleTree(ctx, rules); err != nil {
		return diag.FromErr(err)
	}

	hostnames, err := setHostnames(ctx, prop, d, meta)
	if err != nil {
		return diag.FromErr(fmt.Errorf("%s", err.Error()))
	}

	if err := d.Set("edge_hostnames", hostnames); err != nil {
		return diag.FromErr(fmt.Errorf("%w: %s", tools.ErrValueSet, err.Error()))
	}
	rulesAPI, err := client.GetRuleTree(ctx, v2.GetRuleTreeRequest{
		PropertyID:      prop.PropertyID,
		PropertyVersion: prop.LatestVersion,
		ContractID:      prop.ContractID,
		GroupID:         prop.GroupID,
	})
	if err != nil {
		return diag.FromErr(err)
	}
	rulesAPI.Etag = ""
	body, err := json.Marshal(rulesAPI)
	if err != nil {
		return diag.FromErr(err)
	}

	sha1hashAPI := tools.GetSHAString(string(body))
	logger.Debugf("CREATE SHA from JSON %s", sha1hashAPI)
	logger.Debugf("CREATE Check rules after unmarshal from JSON %s", string(body))

	if err := d.Set("rulessha", sha1hashAPI); err != nil {
		return diag.FromErr(fmt.Errorf("%w: %s", tools.ErrValueSet, err.Error()))
	}
	d.SetId(fmt.Sprintf("%s", prop.PropertyID))
	if err := d.Set("rules", string(body)); err != nil {
		return diag.FromErr(fmt.Errorf("%w: %s", tools.ErrValueSet, err.Error()))
	}
	logger.Debugf("Done")
	return resourcePropertyRead(nil, d, m)
}

func getRules(ctx context.Context, d *schema.ResourceData, property *v2.Property, contract, group string, meta akamai.OperationMeta) (v2.UpdateRulesRequest, error) {
	req := v2.UpdateRulesRequest{}
	logger := meta.Log("getRulesV2")
	req.Rules.Name = "default"
	req.PropertyID = d.Id()
	req.PropertyVersion = property.LatestVersion
	origin, err := createOrigin(d, logger)
	if err != nil {
		return v2.UpdateRulesRequest{}, err
	}

	_, ok := d.GetOk("rules")
	rules := &v2.Rules{Name: "default"}
	if ok {
		logger.Debugf("Unmarshal Rules from JSON")
		rules = unmarshalRulesFromJSON(d)
	}

	cpCode, err := getCPCode(ctx, d, contract, group, meta)
	if err != nil {
		return v2.UpdateRulesRequest{}, err
	}

	logger.Debugf("updateStandardBehaviors")
	req.Rules.Behaviors = updateStandardBehaviors(rules.Behaviors, cpCode, origin, logger)
	logger.Debugf("fixupPerformanceBehaviors")
	fixupPerformanceBehaviors(rules, logger)
	req.Rules = *rules

	return req, nil
}

func setHostnames(ctx context.Context, property *v2.Property, d *schema.ResourceData, meta akamai.OperationMeta) (map[string]string, error) {
	logger := meta.Log("setHostnames")
	client := inst.Client(meta)
	hostnameEdgeHostnames, ok := d.Get("hostnames").(map[string]interface{})
	if !ok {
		return nil, fmt.Errorf("%w: %s, %q", tools.ErrInvalidType, "hostnames", "map[string]interface{}")
	}

	edgeHostnames, err := client.GetEdgeHostnames(ctx, v2.GetEdgeHostnamesRequest{
		ContractID: property.ContractID,
		GroupID:    property.GroupID,
	})
	if err != nil {
		return nil, err
	}
	hostname := v2.UpdatePropertyVersionHostnamesRequest{
		PropertyID:      property.PropertyID,
		PropertyVersion: property.LatestVersion,
		ContractID:      property.ContractID,
		GroupID:         property.GroupID,
	}
	edgeHostnamesMap := make(map[string]string, len(hostnameEdgeHostnames))
	for public, edgeHostname := range hostnameEdgeHostnames {
		edgeHostNameStr, ok := edgeHostname.(string)
		if !ok {
			return nil, fmt.Errorf("%w: %s, %q", tools.ErrInvalidType, "edge_hostname", "string")
		}
		logger.Debugf("Searching for edge hostname: %s, for hostname: %s", edgeHostNameStr, public)
		newEdgeHostname, err := findEdgeHostname(edgeHostnames.EdgeHostnames, "", edgeHostNameStr, "", "")
		if err != nil {
			return nil, fmt.Errorf("edge hostname not found: %s", edgeHostNameStr)
		}
		logger.Debugf("Found edge hostname: %s", newEdgeHostname.Domain)

		hostname.Hostnames.Items = append(hostname.Hostnames.Items, v2.Hostname{
			CnameType:      v2.HostnameCnameTypeEdgeHostname,
			EdgeHostnameID: newEdgeHostname.ID,
			CnameFrom:      public,
			CnameTo:        newEdgeHostname.Domain,
		})
		edgeHostnamesMap[public] = newEdgeHostname.Domain
	}

	_, err = client.UpdatePropertyVersionHostnames(ctx, hostname)
	if err != nil {
		return nil, err
	}
	return edgeHostnamesMap, nil
}

func createProperty(ctx context.Context, contractID, groupID, productID string, d *schema.ResourceData, meta akamai.OperationMeta) (*v2.Property, error) {
	logger := meta.Log("createProperty")
	logger.Debugf("Creating property")

	client := inst.Client(meta)
	propertyName, err := tools.GetStringValue("name", d)
	if err != nil {
		return nil, err
	}
	ruleFormat, err := tools.GetStringValue("rule_format", d)
	if err != nil {
		if !errors.Is(err, tools.ErrNotFound) {
			return nil, err
		}
		ruleFormats, err := client.GetRuleFormats(ctx)
		if err != nil {
			return nil, err
		}
		if len(ruleFormats.RuleFormats.Items) == 0 {
			return nil, fmt.Errorf("no rule formats found")
		}
		ruleFormat = ruleFormats.RuleFormats.Items[len(ruleFormats.RuleFormats.Items)-1]
	}
	prop, err := client.CreateProperty(ctx, v2.CreatePropertyRequest{
		ContractID: contractID,
		GroupID:    groupID,
		Property: v2.PropertyCreate{
			ProductID:    productID,
			PropertyName: propertyName,
			RuleFormat:   ruleFormat,
		},
	})
	if err != nil {
		return nil, err
	}

	newProperty, err := client.GetProperty(ctx, v2.GetPropertyRequest{
		ContractID: contractID,
		GroupID:    groupID,
		PropertyID: prop.PropertyID,
	})
	if err != nil {
		return nil, err
	}
	if len(newProperty.Properties.Items) == 0 {
		return nil, fmt.Errorf("%w: %s", ErrPropertyNotFound, prop.PropertyID)
	}

	logger.Debugf("Property created: %s", prop.PropertyID)
	return newProperty.Properties.Items[0], nil
}

func resourcePropertyDelete(ctx context.Context, d *schema.ResourceData, m interface{}) diag.Diagnostics {
	meta := akamai.Meta(m)
	logger := meta.Log("PAPI", "resourcePropertyDelete")
	client := inst.Client(meta)
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

	resp, err := client.GetProperty(ctx, v2.GetPropertyRequest{
		ContractID: contractID,
		GroupID:    groupID,
		PropertyID: d.Id(),
	})
	if err != nil {
		return diag.FromErr(err)
	}
	property := resp.Property
	if *property.StagingVersion != 0 {
		return diag.FromErr(fmt.Errorf("property is still active on %s and cannot be deleted", v2.VersionStaging))
	}
	if *property.ProductionVersion != 0 {
		return diag.FromErr(fmt.Errorf("property is still active on %s and cannot be deleted", v2.VersionProduction))
	}
	_, err = client.RemoveProperty(ctx, v2.RemovePropertyRequest{
		PropertyID: d.Id(),
		ContractID: contractID,
		GroupID:    groupID,
	})
	if err != nil {
		return diag.FromErr(err)
	}
	d.SetId("")
	logger.Debugf("Done")
	return nil
}

func resourcePropertyImport(ctx context.Context, d *schema.ResourceData, m interface{}) ([]*schema.ResourceData, error) {
	meta := akamai.Meta(m)
	client := inst.Client(meta)
	propertyID := d.Id()

	if !strings.HasPrefix(propertyID, "prp_") {
		keys := []string{v2.SearchKeyPropertyName, v2.SearchKeyHostname, v2.SearchKeyEdgeHostname}
		for _, searchKey := range keys {
			results, err := client.SearchProperties(ctx, v2.SearchRequest{
				Key:   searchKey,
				Value: propertyID,
			})
			if err != nil {
				return nil, err
			}

			if results != nil && len(results.Versions.Items) > 0 {
				propertyID = results.Versions.Items[0].PropertyID
				break
			}
		}
	}
	res, err := client.GetProperty(ctx, v2.GetPropertyRequest{
		PropertyID: propertyID,
	})
	if err != nil {
		return nil, err
	}
	property := res.Property

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

func resourcePropertyRead(ctx context.Context, d *schema.ResourceData, m interface{}) diag.Diagnostics {
	meta := akamai.Meta(m)
	client := inst.Client(meta)
	logger := meta.Log("PAPI", "resourcePropertyRead")
	res, err := client.GetProperty(ctx, v2.GetPropertyRequest{
		PropertyID: d.Id(),
	})
	if err != nil {
		return diag.FromErr(err)
	}
	property := res.Property
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

	getRulesRequest := v2.GetRuleTreeRequest{
		PropertyID:      property.PropertyID,
		PropertyVersion: property.LatestVersion,
		ContractID:      property.ContractID,
		GroupID:         property.GroupID,
	}
	rules, err := client.GetRuleTree(ctx, getRulesRequest)
	if err != nil {
		return diag.FromErr(err)
	}
	rules.Etag = ""
	body, err := json.Marshal(rules)
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
	if *property.StagingVersion > 0 {
		if err := d.Set("staging_version", property.StagingVersion); err != nil {
			return diag.FromErr(fmt.Errorf("%w: %s", tools.ErrValueSet, err.Error()))
		}
	}
	if *property.ProductionVersion > 0 {
		if err := d.Set("production_version", property.ProductionVersion); err != nil {
			return diag.FromErr(fmt.Errorf("%w: %s", tools.ErrValueSet, err.Error()))
		}
	}
	return nil
}

func resourcePropertyUpdate(ctx context.Context, d *schema.ResourceData, m interface{}) diag.Diagnostics {
	meta := akamai.Meta(m)
	client := inst.Client(meta)
	logger := meta.Log("PAPI", "resourcePropertyUpdate")
	logger.Debugf("UPDATING")
	property, err := getProperty(ctx, d.Id(), meta)
	if err != nil {
		return diag.FromErr(err)
	}
	err = ensureEditableVersion(ctx, property, meta)
	if err != nil {
		return diag.FromErr(err)
	}

	rules, err := getRules(ctx, d, property, property.ContractID, property.GroupID, meta)
	if err != nil {
		return diag.FromErr(err)
	}
	if d.HasChange("rule_format") || d.HasChange("rules") {
		ruleFormat, err := tools.GetStringValue("rule_format", d)
		if err != nil {
			if !errors.Is(err, tools.ErrNotFound) {
				return diag.FromErr(err)
			}
		} else {
			property.RuleFormat = ruleFormat
		}
		body, err := json.Marshal(rules)
		if err != nil {
			return diag.FromErr(err)
		}
		if err := d.Set("rules", string(body)); err != nil {
			return diag.FromErr(fmt.Errorf("%w: %s", tools.ErrValueSet, err.Error()))
		}
		logger.Debugf("UPDATE Check rules after unmarshal from Json %s", string(body))
		if _, err := client.UpdateRuleTree(ctx, rules); err != nil {
			logger.Debugf("update rules.Save err: %#v", err)
			return diag.FromErr(err)
		}

		res, err := client.GetRuleTree(ctx, v2.GetRuleTreeRequest{
			PropertyID:      property.PropertyID,
			PropertyVersion: property.LatestVersion,
			ContractID:      property.ContractID,
			GroupID:         property.GroupID,
		})
		if err != nil {
			return diag.FromErr(err)
		}
		res.Etag = ""
		body, err = jsonhooks.Marshal(res)
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
		edgeHostnamesMap, err := setHostnames(ctx, property, d, meta)
		if err != nil {
			return diag.FromErr(fmt.Errorf("setHostnames err: %#v", err))
		}
		if err := d.Set("edge_hostnames", edgeHostnamesMap); err != nil {
			return diag.FromErr(fmt.Errorf("%w: %s", tools.ErrValueSet, err.Error()))
		}
	}

	logger.Debugf("Done")
	return resourcePropertyRead(nil, d, m)
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
	newStr, ok := new.(string)
	if !ok {
		return fmt.Errorf("value is of invalid type: %v; should be %s", new, "string")
	}
	logger.Debugf("OLD: %s", oldStr)
	logger.Debugf("NEW: %s", newStr)
	if !compareRulesJSON(oldStr, newStr) {
		logger.Debugf("CHANGED VALUES: %s %s " + oldStr + " " + newStr)
		if err := d.SetNewComputed("version"); err != nil {
			return fmt.Errorf("%w: %s", tools.ErrValueSet, err.Error())
		}
	}
	return nil
}

// Helpers
func getProperty(ctx context.Context, id string, meta akamai.OperationMeta) (*v2.Property, error) {
	logger := meta.Log("getProperty")
	client := inst.Client(meta)
	logger.Debugf("Fetching property")
	res, err := client.GetProperty(ctx, v2.GetPropertyRequest{
		PropertyID: id,
	})
	if err != nil {
		return nil, err
	}
	return res.Property, nil
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

func getCPCode(ctx context.Context, d tools.ResourceDataFetcher, contractID, groupID string, meta akamai.OperationMeta) (*v2.CPCode, error) {
	client := inst.Client(meta)
	logger := meta.Log("getCPCode")
	if contractID == "" {
		return nil, ErrNoContractProvided
	}
	if groupID == "" {
		return nil, ErrNoGroupProvided
	}
	cpCodeID, err := tools.GetStringValue("cp_code", d)
	if err != nil {
		if !errors.Is(err, tools.ErrNotFound) {
			return nil, err
		}
		return nil, nil
	}
	logger.Debugf("Fetching CP code")
	cpCode, err := client.GetCPCode(ctx, v2.GetCPCodeRequest{
		CPCodeID:   cpCodeID,
		ContractID: contractID,
		GroupID:    groupID,
	})
	if err != nil {
		return nil, err
	}
	logger.Debugf("CP code found: %s", cpCode.CPCode.ID)
	return &cpCode.CPCode, nil
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

func createOrigin(d interface{}, logger log.Interface) (*v2.RuleOptionsMap, error) {
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

	if err != nil {
		if !errors.Is(err, tools.ErrNotFound) {
			return nil, err
		}
		return nil, nil
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
	ov := v2.RuleOptionsMap(originValues)
	return &ov, nil
}

func fixupPerformanceBehaviors(rules *v2.Rules, logger log.Interface) {
	behavior, err := findBehavior("/Performance/sureRoute", rules)
	if err != nil || behavior.Options["testObjectUrl"] != "" {
		return
	}

	logger.Debugf("Fixing Up SureRoute Behavior")
	behavior.Options = mergeOptions(behavior.Options, v2.RuleOptionsMap{
		"testObjectUrl":   "/akamai/sureroute-testobject.html",
		"enableCustomKey": false,
		"enabled":         false,
	})
}

func updateStandardBehaviors(behaviors []v2.RuleBehavior, cpCode *v2.CPCode, origin *v2.RuleOptionsMap, logger log.Interface) []v2.RuleBehavior {
	logger.Debugf("cpCode: %#v", cpCode)
	if cpCode != nil {
		b := v2.RuleBehavior{
			Name: "cpCode",
			Options: v2.RuleOptionsMap{
				"value": v2.RuleOptionsMap{
					"id": cpCode.ID,
				},
			},
		}
		behaviors = mergeBehaviors(behaviors, b)
	}

	if origin != nil {
		b := v2.RuleBehavior{
			Name:    "origin",
			Options: *origin,
		}
		behaviors = mergeBehaviors(behaviors, b)
	}
	return behaviors
}

// TODO: discuss how property rules should be handled
func unmarshalRulesFromJSON(d *schema.ResourceData) *v2.Rules {
	// Default Rules
	rules, ok := d.GetOk("rules")
	if !ok {
		return nil
	}

	propertyRules := &v2.Rules{Name: "default"}
	rulesJSON := gjson.Get(rules.(string), "rules")
	rulesJSON.ForEach(func(key, value gjson.Result) bool {
		if key.String() == "behaviors" {
			behavior := gjson.Parse(value.String())
			if gjson.Get(behavior.String(), "#.name").Exists() {
				behavior.ForEach(func(key, value gjson.Result) bool {
					bb, ok := value.Value().(map[string]interface{})
					if ok {
						for k, v := range bb {
							log.Debugf("k:", k, "v:", v)
						}

						beh := v2.RuleBehavior{Options: v2.RuleOptionsMap{}}

						beh.Name = bb["name"].(string)
						boptions, ok := bb["options"]
						if ok {
							beh.Options = boptions.(map[string]interface{})
						}

						propertyRules.Behaviors = mergeBehaviors(propertyRules.Behaviors, beh)
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
						newCriteria := v2.RuleBehavior{Options: v2.RuleOptionsMap{}}
						newCriteria.Name = cc["name"].(string)

						coptions, ok := cc["option"]
						if ok {
							//								println("OPTIONS ", coptions)
							newCriteria.Options = coptions.(map[string]interface{})
						}
						propertyRules.Criteria = append(propertyRules.Criteria, newCriteria)
					}
					return true
				})
			} // if ok criteria
		} /// if ok behaviors

		if key.String() == "children" {
			childRules := gjson.Parse(value.String())
			for _, rule := range extractRulesJSON(d, childRules) {
				propertyRules.Children = append(propertyRules.Children, rule)
			}
		}

		if key.String() == "variables" {
			variables := gjson.Parse(value.String())
			variables.ForEach(func(key, value gjson.Result) bool {
				variableMap, ok := value.Value().(map[string]interface{})
				if ok {
					newVariable := v2.RuleVariable{}
					newVariable.Name = variableMap["name"].(string)
					newVariable.Description = variableMap["description"].(string)
					newVariable.Value = variableMap["value"].(string)
					newVariable.Hidden = variableMap["hidden"].(bool)
					newVariable.Sensitive = variableMap["sensitive"].(bool)
					propertyRules.Variables = addVariable(propertyRules.Variables, newVariable)
				}
				return true
			}) //variables

		}

		if key.String() == "options" {
			options := gjson.Parse(value.String())
			options.ForEach(func(key, value gjson.Result) bool {
				switch {
				case key.String() == "is_secure" && value.Bool():
					propertyRules.Options.IsSecure = value.Bool()
				}

				return true
			})
		}

		return true // keep iterating
	}) // for loop rules

	// ADD vars from variables resource
	jsonvars, ok := d.GetOk("variables")
	if ok {
		variables := gjson.Parse(jsonvars.(string))
		result := gjson.Get(variables.String(), "variables")
		result.ForEach(func(key, value gjson.Result) bool {
			variableMap, ok := value.Value().(map[string]interface{})
			if ok {
				newVariable := v2.RuleVariable{}
				newVariable.Name = variableMap["name"].(string)
				newVariable.Description = variableMap["description"].(string)
				newVariable.Value = variableMap["value"].(string)
				newVariable.Hidden = variableMap["hidden"].(bool)
				newVariable.Sensitive = variableMap["sensitive"].(bool)
				propertyRules.Variables = addVariable(propertyRules.Variables, newVariable)
			}
			return true
		}) //variables
	}

	// ADD isSecure from resource
	isSecure, set := d.GetOkExists("is_secure")
	if set && isSecure.(bool) {
		propertyRules.Options.IsSecure = true
	} else if set && !isSecure.(bool) {
		propertyRules.Options.IsSecure = false
	}

	// ADD cpCode from resource
	cpCode, set := d.GetOk("cp_code")
	if set {
		beh := v2.RuleBehavior{Options: v2.RuleOptionsMap{}}
		beh.Name = "cpCode"
		beh.Options = v2.RuleOptionsMap{
			"value": map[string]interface{}{
				"id": cpCode.(string),
			},
		}
		propertyRules.Behaviors = mergeBehaviors(propertyRules.Behaviors, beh)
	}
	return propertyRules
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
func extractRulesJSON(d interface{}, drules gjson.Result) []v2.Rules {
	var rules []v2.Rules
	drules.ForEach(func(key, value gjson.Result) bool {
		rule := v2.Rules{Name: "default"}
		vv, ok := value.Value().(map[string]interface{})
		if ok {
			rule.Name, _ = vv["name"].(string)
			rule.Comment, _ = vv["comments"].(string)
			criteriaMustSatisfy, ok := vv["criteria_match"]
			if ok {
				if criteriaMustSatisfy.(string) == "all" {
					rule.CriteriaMustSatisfy = v2.RuleCriteriaMustSatisfyAll
				}

				if criteriaMustSatisfy.(string) == "any" {
					rule.CriteriaMustSatisfy = v2.RuleCriteriaMustSatisfyAny
				}
			}
			log.Debugf("extractRulesJSON Set criteriaMustSatisfy RESULT RULE value set " + string(rule.CriteriaMustSatisfy) + " " + rule.Name + " " + rule.Comment)

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
							newBehavior := v2.RuleBehavior{Options: v2.RuleOptionsMap{}}
							newBehavior.Name = behaviorMap["name"].(string)
							behaviorOptions, ok := behaviorMap["options"]
							if ok {
								newBehavior.Options = behaviorOptions.(map[string]interface{})
							}
							rule.Behaviors = mergeBehaviors(rule.Behaviors, newBehavior)
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
							newCriteria := v2.RuleBehavior{Options: v2.RuleOptionsMap{}}
							newCriteria.Name = criteriaMap["name"].(string)
							criteriaOptions, ok := criteriaMap["options"]
							if ok {
								newCriteria.Options = criteriaOptions.(map[string]interface{})
							}
							rule.Criteria = append(rule.Criteria, newCriteria)
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
							newVariable := v2.RuleVariable{}
							newVariable.Name = variableMap["name"].(string)
							newVariable.Description = variableMap["description"].(string)
							newVariable.Value = variableMap["value"].(string)
							newVariable.Hidden = variableMap["hidden"].(bool)
							newVariable.Sensitive = variableMap["sensitive"].(bool)
							rule.Variables = addVariable(rule.Variables, newVariable)
						}
						return true
					}) //variables
				}

				if key.String() == "children" {
					childRules := gjson.Parse(value.String())
					for _, newRule := range extractRulesJSON(d, childRules) {
						rule.Children = append(rule.Children, newRule)
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

func extractRules(drules *schema.Set) ([]v2.Rules, error) {

	var rules []v2.Rules
	for _, v := range drules.List() {
		rule := v2.Rules{Name: "default"}
		vv, ok := v.(map[string]interface{})
		if ok {
			rule.Name = vv["name"].(string)
			rule.Comment = vv["comment"].(string)

			criteriaMustSatisfy, ok := vv["criteria_match"]
			if ok {
				if criteriaMustSatisfy.(string) == "all" {
					rule.CriteriaMustSatisfy = v2.RuleCriteriaMustSatisfyAll
				}

				if criteriaMustSatisfy.(string) == "any" {
					rule.CriteriaMustSatisfy = v2.RuleCriteriaMustSatisfyAny
				}
			}
			behaviors, ok := vv["behavior"]
			if ok {
				for _, behavior := range behaviors.(*schema.Set).List() {
					behaviorMap, ok := behavior.(map[string]interface{})
					if ok {
						newBehavior := v2.RuleBehavior{}
						newBehavior.Name = behaviorMap["name"].(string)
						behaviorOptions, ok := behaviorMap["option"]
						if ok {
							opts, err := extractOptions(behaviorOptions.(*schema.Set))
							if err != nil {
								return nil, err
							}
							newBehavior.Options = opts
						}
						rule.Behaviors = mergeBehaviors(rule.Behaviors, newBehavior)
					}
				}
			}

			criterias, ok := vv["criteria"]
			if ok {
				for _, criteria := range criterias.(*schema.Set).List() {
					criteriaMap, ok := criteria.(map[string]interface{})
					if ok {
						newCriteria := v2.RuleBehavior{}
						newCriteria.Name = criteriaMap["name"].(string)
						criteriaOptions, ok := criteriaMap["option"]
						if ok {
							crit, err := extractOptions(criteriaOptions.(*schema.Set))
							if err != nil {
								return nil, err
							}
							newCriteria.Options = crit
						}
						rule.Criteria = append(rule.Criteria, newCriteria)
					}
				}
			}

			variables, ok := vv["variable"]
			if ok {
				for _, variable := range variables.(*schema.Set).List() {
					variableMap, ok := variable.(map[string]interface{})
					if ok {
						newVariable := v2.RuleVariable{}
						newVariable.Name = variableMap["name"].(string)
						newVariable.Description = variableMap["description"].(string)
						newVariable.Value = variableMap["value"].(string)
						newVariable.Hidden = variableMap["hidden"].(bool)
						newVariable.Sensitive = variableMap["sensitive"].(bool)
						rule.Variables = addVariable(rule.Variables, newVariable)
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
					rule.Children = append(rule.Children, newRule)
				}
			}
		}
		rules = append(rules, rule)
	}

	return rules, nil
}

func findProperty(ctx context.Context, name string, meta akamai.OperationMeta) (*v2.Property, error) {
	client := inst.Client(meta)
	results, err := client.SearchProperties(ctx, v2.SearchRequest{Key: v2.SearchKeyPropertyName, Value: name})
	if err != nil {
		return nil, err
	}
	if len(results.Versions.Items) == 0 {
		return nil, fmt.Errorf("%w: %s", ErrPropertyNotFound, name)
	}

	property, err := client.GetProperty(ctx, v2.GetPropertyRequest{
		ContractID: results.Versions.Items[0].ContractID,
		GroupID:    results.Versions.Items[0].GroupID,
		PropertyID: results.Versions.Items[0].PropertyID,
	})
	if err != nil {
		return nil, err
	}
	if len(property.Properties.Items) == 0 {
		return nil, fmt.Errorf("%w: %s", ErrPropertyNotFound, name)
	}
	return property.Properties.Items[0], nil
}

func ensureEditableVersion(ctx context.Context, property *v2.Property, meta akamai.OperationMeta) error {
	client := inst.Client(meta)
	latestVersion, err := client.GetLatestVersion(ctx, v2.GetLatestVersionRequest{
		PropertyID: property.PropertyID,
		ContractID: property.ContractID,
		GroupID:    property.ContractID,
	})
	if err != nil {
		return err
	}
	if len(latestVersion.Versions.Items) == 0 {
		return fmt.Errorf("%w: %s", ErrVersionNotFound, err)
	}

	if latestVersion.Versions.Items[0].ProductionStatus != v2.VersionStatusInactive ||
		latestVersion.Versions.Items[1].StagingStatus != v2.VersionStatusInactive {
		// The latest version has been activated on either production or staging, so we need to create a new version to apply changes on
		_, err := client.CreatePropertyVersion(ctx, v2.CreatePropertyVersionRequest{
			PropertyID: property.PropertyID,
			ContractID: property.ContractID,
			GroupID:    property.GroupID,
			Version: v2.PropertyVersionCreate{
				CreateFromVersion: latestVersion.Versions.Items[0].PropertyVersion,
			},
		})
		if err != nil {
			return fmt.Errorf("%w: %s", ErrVersionCreate, err.Error())
		}
	}

	return nil
}

func mergeBehaviors(old []v2.RuleBehavior, new v2.RuleBehavior) []v2.RuleBehavior {
	for i := range old {
		if new.Name == "cpCode" || new.Name == "origin" {
			if old[i].Name == new.Name {
				old[i].Options = mergeOptions(old[i].Options, new.Options)
				return old
			}
		}
	}

	return append(old, new)
}

// MergeOptions merges the given options with the existing options
func mergeOptions(old, new v2.RuleOptionsMap) v2.RuleOptionsMap {
	options := make(v2.RuleOptionsMap)
	for k, v := range old {
		options[k] = v
	}
	for k, v := range new {
		options[k] = v
	}
	return options
}

func addVariable(old []v2.RuleVariable, new v2.RuleVariable) []v2.RuleVariable {
	for i := range old {
		if old[i].Name == new.Name {
			old[i] = new
			return old
		}
	}

	return append(old, new)
}

// FindBehavior locates a specific behavior by path
func findBehavior(path string, rules *v2.Rules) (v2.RuleBehavior, error) {
	if len(path) <= 1 {
		return v2.RuleBehavior{}, fmt.Errorf("invalid path: %s", path)
	}

	rule, err := findParentRule(path, rules)
	if err != nil {
		return v2.RuleBehavior{}, err
	}

	sep := "/"
	segments := strings.Split(path, sep)
	behaviorName := strings.ToLower(segments[len(segments)-1])
	for _, behavior := range rule.Behaviors {
		if strings.ToLower(behavior.Name) == behaviorName {
			return behavior, nil
		}
	}

	return v2.RuleBehavior{}, fmt.Errorf("behavior not found for path: %s", path)
}

// Find the parent rule for a given rule, criteria, or behavior path
func findParentRule(path string, rules *v2.Rules) (*v2.Rules, error) {
	sep := "/"
	segments := strings.Split(strings.ToLower(strings.TrimPrefix(path, sep)), sep)
	parentPath := strings.Join(segments[0:len(segments)-1], sep)

	return findRule(parentPath, rules)
}

// FindRule locates a specific rule by path
func findRule(path string, rules *v2.Rules) (*v2.Rules, error) {
	if path == "" {
		return rules, nil
	}

	sep := "/"
	segments := strings.Split(path, sep)

	currentRule := rules
	for _, segment := range segments {
		found := false
		for _, rule := range currentRule.Children {
			if strings.ToLower(rule.Name) == segment {
				currentRule = &rule
				found = true
			}
		}
		if !found {
			return nil, ErrRulesNotFound
		}
	}

	return currentRule, nil
}

func findEdgeHostname(edgeHostnames v2.EdgeHostnameItems, id, domain, suffix, prefix string) (*v2.EdgeHostnameGetItem, error) {
	if suffix == "" && domain != "" {
		suffix = "edgesuite.net"
		if strings.HasSuffix(domain, "edgekey.net") {
			suffix = "edgekey.net"
		}
	}

	if prefix == "" && domain != "" {
		prefix = strings.TrimSuffix(domain, "."+suffix)
	}

	if len(edgeHostnames.Items) == 0 {
		return nil, errors.New("no hostnames found, did you call GetHostnames()?")
	}

	for _, eHn := range edgeHostnames.Items {
		if (eHn.DomainPrefix == prefix && eHn.DomainSuffix == suffix) || eHn.ID == id {
			return &eHn, nil
		}
	}

	return nil, nil
}
