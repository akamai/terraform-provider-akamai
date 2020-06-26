package akamai

import (
	"errors"
	"fmt"
	"strconv"
	"strings"

	log "github.com/sirupsen/logrus"

	"github.com/akamai/AkamaiOPEN-edgegrid-golang/jsonhooks-v1"
	"github.com/akamai/AkamaiOPEN-edgegrid-golang/papi-v1"
	"github.com/hashicorp/terraform-plugin-sdk/helper/schema"
	"github.com/hashicorp/terraform-plugin-sdk/helper/validation"
	"github.com/tidwall/gjson"
)

func resourceProperty() *schema.Resource {
	return &schema.Resource{
		Create: resourcePropertyCreate,
		Read:   resourcePropertyRead,
		Update: resourcePropertyUpdate,
		Delete: resourcePropertyDelete,
		//Exists: resourcePropertyExists,
		CustomizeDiff: resourceCustomDiffCustomizeDiff,
		Importer: &schema.ResourceImporter{
			State: resourcePropertyImport,
		},
		Schema: akamaiPropertySchema,
	}
}

var akamaiPropertySchema = map[string]*schema.Schema{
	"account": &schema.Schema{
		Type:     schema.TypeString,
		Computed: true,
	},
	"contract": &schema.Schema{
		Type:     schema.TypeString,
		Optional: true,
		ForceNew: true,
	},
	"group": &schema.Schema{
		Type:     schema.TypeString,
		Optional: true,
		ForceNew: true,
	},
	"product": &schema.Schema{
		Type:     schema.TypeString,
		Optional: true,
		Default:  "prd_SPM",
		ForceNew: true,
	},
	"rule_format": &schema.Schema{
		Type:     schema.TypeString,
		Optional: true,
	},
	// Will get added to the default rule
	"cp_code": &schema.Schema{
		Type:     schema.TypeString,
		Optional: true,
	},
	"name": &schema.Schema{
		Type:     schema.TypeString,
		Required: true,
		ForceNew: true,
	},
	"version": &schema.Schema{
		Type:     schema.TypeInt,
		Computed: true,
	},
	"staging_version": &schema.Schema{
		Type:     schema.TypeInt,
		Computed: true,
	},
	"production_version": &schema.Schema{
		Type:     schema.TypeInt,
		Computed: true,
	},
	"contact": &schema.Schema{
		Type:     schema.TypeSet,
		Required: true,
		Elem:     &schema.Schema{Type: schema.TypeString},
	},
	"edge_hostnames": &schema.Schema{
		Type:     schema.TypeMap,
		Computed: true,
		Elem:     &schema.Schema{Type: schema.TypeString},
	},
	"hostnames": &schema.Schema{
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
		ValidateFunc:     validation.ValidateJsonString,
		DiffSuppressFunc: suppressEquivalentJsonDiffs,
	},
	"variables": {
		Type:     schema.TypeString,
		Optional: true,
	},
	"rulessha": &schema.Schema{
		Type:     schema.TypeString,
		Computed: true,
	},
}

func resourcePropertyCreate(d *schema.ResourceData, meta interface{}) error {
	CorrelationID := "[PAPI][resourcePropertyCreate-" + CreateNonce() + "]"
	//log..SetPrefix("resourcePropertyCreate-" + CreateNonce())

	//log.WithField("GUID", getCorrelationID)
	//log.Printf("[DEBUG] CREATE HEADER GUID %s\n", getCorrelationID)
	//PrintLogHeader()

	d.Partial(true)

	group, err := getGroup(d, CorrelationID)
	if err != nil {
		return err
	}

	contract, err := getContract(d, CorrelationID)
	if err != nil {
		return err
	}

	product, err := getProduct(d, contract, CorrelationID)
	if err != nil {
		return err
	}

	var property *papi.Property
	if property = findProperty(d, CorrelationID); property == nil {
		if product == nil {
			return errors.New("product must be specified to create a new property")
		}

		property, err = createProperty(contract, group, product, d, CorrelationID)
		if err != nil {
			return err
		}
	}

	err = ensureEditableVersion(property, CorrelationID)
	if err != nil {
		return err
	}
	d.Set("account", property.AccountID)
	d.Set("version", property.LatestVersion)

	// The API now has data, so save the partial state
	d.SetId(property.PropertyID)
	d.SetPartial("name")
	d.SetPartial("rule_format")
	d.SetPartial("contract")
	d.SetPartial("group")
	d.SetPartial("product")
	d.SetPartial("clone_from")
	d.SetPartial("network")
	d.SetPartial("cp_code")

	rules, err := getRules(d, property, contract, group, CorrelationID)
	if err != nil {
		return err
	}
	err = rules.Save(CorrelationID)
	if err != nil {
		if err == papi.ErrorMap[papi.ErrInvalidRules] && len(rules.Errors) > 0 {
			var msg string
			for _, v := range rules.Errors {
				msg = msg + fmt.Sprintf("\n Rule validation error: %s %s %s %s %s", v.Type, v.Title, v.Detail, v.Instance, v.BehaviorName)
			}
			return errors.New("Error - Invalid Property Rules" + msg)
		}
		return err
	}
	d.SetPartial("default")
	d.SetPartial("origin")

	ehnMap, err := setHostnames(property, d)

	d.SetPartial("ipv6")
	d.Set("edge_hostnames", ehnMap)

	rulesAPI, err := property.GetRules(CorrelationID)
	rulesAPI.Etag = ""
	jsonBody, err := jsonhooks.Marshal(rulesAPI)
	if err != nil {
		return err
	}

	sha1hashAPI := getSHAString(string(jsonBody))
	log.Printf("[DEBUG]"+CorrelationID+" CREATE SHA from Json %s\n", sha1hashAPI)
	log.Printf("[DEBUG]"+CorrelationID+" CREATE Check rules after unmarshal from Json %s\n", string(jsonBody))

	d.Set("rulessha", sha1hashAPI)
	//d.SetId(fmt.Sprintf("%s-%s", property.PropertyID, sha1hashAPI))
	d.SetId(fmt.Sprintf("%s", property.PropertyID))

	if err == nil {
		d.Set("rules", string(jsonBody))
	}

	d.Partial(false)
	log.Println("[DEBUG] Done")
	//PrintLogFooter()
	return resourcePropertyRead(d, meta)
}

func getRules(d *schema.ResourceData, property *papi.Property, contract *papi.Contract, group *papi.Group, correlationid string) (*papi.Rules, error) {
	rules := papi.NewRules()
	rules.Rule.Name = "default"
	rules.PropertyID = d.Id()
	rules.PropertyVersion = property.LatestVersion

	origin, err := createOrigin(d)
	if err != nil {
		return nil, err
	}

	// get rules from the TF config

	//rulecheck
	_, ok := d.GetOk("rules")

	if ok {
		log.Printf("[DEBUG]" + correlationid + "  Unmarshal Rules from JSON")
		unmarshalRulesFromJSON(d, rules)
	}

	if ruleFormat, ok := d.GetOk("rule_format"); ok {
		rules.RuleFormat = ruleFormat.(string)
	} else {
		ruleFormats := papi.NewRuleFormats()
		rules.RuleFormat, err = ruleFormats.GetLatest(correlationid)
		if err != nil {
			return nil, err
		}
	}

	if ok := d.HasChange("rule_format"); ok {
	}

	cpCode, err := getCPCode(d, contract, group)
	if err != nil {
		return nil, err
	}

	log.Printf("[DEBUG]" + correlationid + "  updateStandardBehaviors")
	updateStandardBehaviors(rules, cpCode, origin)
	log.Printf("[DEBUG]" + correlationid + "  fixupPerformanceBehaviors")
	fixupPerformanceBehaviors(rules)

	return rules, nil
}

func setHostnames(property *papi.Property, d *schema.ResourceData) (map[string]string, error) {
	hostnameEdgeHostnames := d.Get("hostnames").(map[string]interface{})

	ehns, err := papi.GetEdgeHostnames(property.Contract, property.Group, "")
	if err != nil {
		return nil, err
	}

	hostnames := papi.NewHostnames()
	hostnames.PropertyVersion = property.LatestVersion
	hostnames.PropertyID = property.PropertyID

	ehnMap := make(map[string]string, len(hostnameEdgeHostnames))
	for public, edgeHostname := range hostnameEdgeHostnames {
		ehn := ehns.NewEdgeHostname()
		ehn.EdgeHostnameDomain = edgeHostname.(string)
		log.Printf("[DEBUG] Searching for edge hostname: %s, for hostname: %s", edgeHostname.(string), public)
		ehn, err = ehns.FindEdgeHostname(ehn)
		log.Printf("[DEBUG] Found edge hostname: %s", ehn.EdgeHostnameDomain)
		if err != nil {
			return nil, fmt.Errorf("edge hostname not found: %s", ehn.EdgeHostnameDomain)
		}

		hostname := hostnames.NewHostname()
		hostname.EdgeHostnameID = ehn.EdgeHostnameID
		hostname.CnameFrom = public
		hostname.CnameTo = ehn.EdgeHostnameDomain

		ehnMap[public] = ehn.EdgeHostnameDomain
	}

	err = hostnames.Save()
	if err != nil {
		return nil, err
	}

	return ehnMap, nil
}

func createProperty(contract *papi.Contract, group *papi.Group, product *papi.Product, d *schema.ResourceData, correlationid string) (*papi.Property, error) {
	log.Println("[DEBUG] Creating property")

	property, err := group.NewProperty(contract)
	if err != nil {
		return nil, err
	}

	property.ProductID = product.ProductID
	property.PropertyName = d.Get("name").(string)

	if ruleFormat, ok := d.GetOk("rule_format"); ok {
		property.RuleFormat = ruleFormat.(string)
	} else {
		ruleFormats := papi.NewRuleFormats()
		property.RuleFormat, err = ruleFormats.GetLatest(correlationid)
		if err != nil {
			return nil, err
		}
	}

	err = property.Save(correlationid)
	if err != nil {
		return nil, err
	}

	log.Printf("[DEBUG] Property created: %s\n", property.PropertyID)
	return property, nil
}

func resourcePropertyDelete(d *schema.ResourceData, meta interface{}) error {
	CorrelationID := "[PAPI][resourcePropertyDelete-" + CreateNonce() + "]"
	//PrintLogHeader()
	log.Printf("[DEBUG] DELETING")
	contractID, ok := d.GetOk("contract")
	if !ok {
		return errors.New("missing contract ID")
	}
	groupID, ok := d.GetOk("group")
	if !ok {
		return errors.New("missing group ID")
	}

	propertyID := d.Id()

	property := papi.NewProperty(papi.NewProperties())
	property.PropertyID = propertyID
	property.Contract = &papi.Contract{ContractID: contractID.(string)}
	property.Group = &papi.Group{GroupID: groupID.(string)}

	e := property.GetProperty(CorrelationID)
	if e != nil {
		return e
	}

	if property.StagingVersion != 0 {
		return fmt.Errorf("property is still active on %s and cannot be deleted", papi.NetworkStaging)
	}

	if property.ProductionVersion != 0 {
		return fmt.Errorf("property is still active on %s and cannot be deleted", papi.NetworkProduction)
	}

	e = property.Delete(CorrelationID)
	if e != nil {
		return e
	}

	d.SetId("")

	log.Println("[DEBUG] Done")
	//PrintLogFooter()
	return nil
}

func resourcePropertyImport(d *schema.ResourceData, meta interface{}) ([]*schema.ResourceData, error) {

	propertyID := d.Id()

	if !strings.HasPrefix(propertyID, "prp_") {
		for _, searchKey := range []papi.SearchKey{papi.SearchByPropertyName, papi.SearchByHostname, papi.SearchByEdgeHostname} {
			results, err := papi.Search(searchKey, propertyID, "") //<--correlationid
			if err != nil {
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
	e := property.GetProperty("")
	if e != nil {
		return nil, e
	}

	d.Set("account", property.AccountID)
	d.Set("contract", property.ContractID)
	d.Set("group", property.GroupID)

	d.Set("name", property.PropertyName)
	d.Set("version", property.LatestVersion)
	d.SetId(property.PropertyID)

	return []*schema.ResourceData{d}, nil
}

func resourcePropertyRead(d *schema.ResourceData, meta interface{}) error {
	CorrelationID := "[PAPI][resourcePropertyRead-" + CreateNonce() + "]"
	//PrintLogHeader()
	d.Partial(true)
	property := papi.NewProperty(papi.NewProperties())
	property.PropertyID = d.Id()
	err := property.GetProperty(CorrelationID)
	if err != nil {
		return err
	}

	d.Set("account", property.AccountID)
	d.Set("contract", property.ContractID)
	d.Set("group", property.GroupID)
	d.Set("name", property.PropertyName)
	d.Set("note", property.Note)

	rules, err := property.GetRules(CorrelationID)
	rules.Etag = ""
	jsonBody, err := jsonhooks.Marshal(rules)
	if err != nil {
		return err
	}

	sha1hashAPI := getSHAString(string(jsonBody))
	log.Printf("[DEBUG]"+CorrelationID+"  READ SHA from Json %s\n", sha1hashAPI)

	if err == nil {
		log.Printf("[DEBUG]"+CorrelationID+"  READ Rules from API : %s\n", string(jsonBody))
		d.Set("rules", string(jsonBody))
	}
	d.Set("rulessha", sha1hashAPI)

	//d.Set("rulessha", RandStringBytesMaskImpr(8))
	d.SetPartial("rulessha")
	//log.Printf("[DEBUG] READ Check rules after unmarshal from Json %s\n", string(jsonBody))

	if rules.RuleFormat != "" {
		d.Set("rule_format", rules.RuleFormat)
	} else {
		d.Set("rule_format", property.RuleFormat)
	}

	log.Printf("[DEBUG]"+CorrelationID+"  Property RuleFormat from API : %s\n", property.RuleFormat)
	d.Set("version", property.LatestVersion)
	log.Printf("[DEBUG]"+CorrelationID+"  Property Version from API : %d\n", property.LatestVersion)
	if property.StagingVersion > 0 {
		d.Set("staging_version", property.StagingVersion)
	}
	if property.ProductionVersion > 0 {
		d.Set("production_version", property.ProductionVersion)
	}

	d.Partial(false)
	//PrintLogFooter()
	return nil
}

func resourcePropertyUpdate(d *schema.ResourceData, meta interface{}) error {
	CorrelationID := "[PAPI][resourcePropertyUpdate-" + CreateNonce() + "]"
	//PrintLogHeader()
	log.Printf("[DEBUG]" + CorrelationID + "  UPDATING")
	d.Partial(true)

	property, e := getProperty(d, CorrelationID)
	if e != nil {
		return e
	}

	err := ensureEditableVersion(property, CorrelationID)
	if err != nil {
		return err
	}

	rules, err := getRules(d, property, property.Contract, property.Group, CorrelationID)

	if d.HasChange("rule_format") || d.HasChange("rules") {
		if ruleFormat, ok := d.GetOk("rule_format"); ok {
			property.RuleFormat = ruleFormat.(string)
			rules.RuleFormat = ruleFormat.(string)
		}

		jsonBody, err := jsonhooks.Marshal(rules)
		if err != nil {
			return err
		}
		if err == nil {
			d.Set("rules", string(jsonBody))
		}

		log.Printf("[DEBUG]"+CorrelationID+"  UPDATE Check rules after unmarshal from Json %s\n", string(jsonBody))
		e = rules.Save(CorrelationID)
		if e != nil {
			if e == papi.ErrorMap[papi.ErrInvalidRules] && len(rules.Errors) > 0 {
				var msg string
				for _, v := range rules.Errors {
					msg = msg + fmt.Sprintf("\n Rule validation error: %s %s %s %s %s", v.Type, v.Title, v.Detail, v.Instance, v.BehaviorName)
				}
				return errors.New("Error - Invalid Property Rules" + msg)
			}
			log.Printf("update rules.Save err: %#v", e)
			return fmt.Errorf("update rules.Save err: %#v", e)
		}

		rules, err = property.GetRules(CorrelationID)
		rules.Etag = ""
		jsonBody, err = jsonhooks.Marshal(rules)
		if err != nil {
			return err
		}

		sha1hashAPI := getSHAString(string(jsonBody))
		log.Printf("[DEBUG]"+CorrelationID+"  UPDATE SHA from Json %s\n", sha1hashAPI)
		d.Set("rulessha", sha1hashAPI)
		d.SetPartial("rulessha")
	}

	d.Set("version", property.LatestVersion)
	d.SetPartial("default")
	d.SetPartial("origin")

	if d.HasChange("hostnames") {
		ehnMap, err := setHostnames(property, d)
		if err != nil {
			return fmt.Errorf("setHostnames err: %#v", err)
		}

		d.Set("edge_hostnames", ehnMap)
	}

	d.Partial(false)

	log.Println("[DEBUG]" + CorrelationID + "  Done")
	//PrintLogFooter()
	return resourcePropertyRead(d, meta)
}

func resourceCustomDiffCustomizeDiff(d *schema.ResourceDiff, meta interface{}) error {
	log.Println("[DEBUG] resourceCustomDiffCustomizeDiff " + d.Id())
	// Note that this gets put into state after the update, regardless of whether
	// or not anything is acted upon in the diff.

	old, new := d.GetChange("rules")

	log.Println("[DEBUG] resourceCustomDiffCustomizeDiff OLD " + old.(string))
	log.Println("[DEBUG] resourceCustomDiffCustomizeDiff NEW " + new.(string))
	if !suppressEquivalentJsonPendingDiffs(old.(string), new.(string), d) {
		log.Println("[DEBUG] resourceCustomDiffCustomizeDiff CHANGED VALUES " + old.(string) + " " + new.(string))
		d.SetNewComputed("version")
	}

	return nil
}

// Helpers
func getProperty(d interface{}, correlationid string) (*papi.Property, error) {
	log.Println("[DEBUG] Fetching property")
	var propertyID string

	switch d.(type) {
	case *schema.ResourceData:
		propertyID = d.(*schema.ResourceData).Id()
	case *schema.ResourceDiff:
		propertyID = d.(*schema.ResourceDiff).Id()
	default:
		propertyID = d.(*schema.ResourceData).Id()
	}

	property := papi.NewProperty(papi.NewProperties())
	property.PropertyID = propertyID
	e := property.GetProperty(correlationid)
	return property, e
}

func getGroup(d *schema.ResourceData, correlationid string) (*papi.Group, error) {
	log.Println("[DEBUG] Fetching groups")
	groupID, ok := d.GetOk("group")

	if !ok {
		return nil, nil
	}

	groups := papi.NewGroups()
	e := groups.GetGroups(correlationid)
	if e != nil {
		return nil, e
	}

	group, e := groups.FindGroup(groupID.(string))
	if e != nil {
		return nil, e
	}

	log.Printf("[DEBUG]"+correlationid+"  Group found: %s\n", group.GroupID)
	return group, nil
}

func getContract(d *schema.ResourceData, correlationid string) (*papi.Contract, error) {
	log.Println("[DEBUG]" + correlationid + "  Fetching contract")
	contractID, ok := d.GetOk("contract")
	if !ok {
		return nil, nil
	}

	contracts := papi.NewContracts()
	e := contracts.GetContracts(correlationid)
	if e != nil {
		return nil, e
	}

	contract, e := contracts.FindContract(contractID.(string))
	if e != nil {
		return nil, e
	}

	log.Printf("[DEBUG] Contract found: %s\n", contract.ContractID)
	return contract, nil
}

func getCPCode(d interface{}, contract *papi.Contract, group *papi.Group) (*papi.CpCode, error) {
	if contract == nil || group == nil {
		return nil, nil
	}
	var cpCodeID interface{}
	var ok bool

	switch d.(type) {
	case *schema.ResourceData:
		cpCodeID, ok = d.(*schema.ResourceData).GetOk("cp_code")
	case *schema.ResourceDiff:
		cpCodeID, ok = d.(*schema.ResourceDiff).GetOk("cp_code")
	default:
		cpCodeID, ok = d.(*schema.ResourceData).GetOk("cp_code")
	}

	//cpCodeID, ok := d.GetOk("cp_code")
	if !ok {
		return nil, nil
	}

	log.Println("[DEBUG] Fetching CP code")
	cpCode := papi.NewCpCodes(contract, group).NewCpCode()
	cpCode.CpcodeID = cpCodeID.(string)
	err := cpCode.GetCpCode()
	if err != nil {
		return nil, err
	}

	log.Printf("[DEBUG] CP code found: %s\n", cpCode.CpcodeID)
	return cpCode, nil
}

func getCPCodeDiff(d *schema.ResourceDiff, contract *papi.Contract, group *papi.Group) (*papi.CpCode, error) {
	if contract == nil || group == nil {
		return nil, nil
	}

	cpCodeID, ok := d.GetOk("cp_code")
	if !ok {
		return nil, nil
	}

	log.Println("[DEBUG] Fetching CP code")
	cpCode := papi.NewCpCodes(contract, group).NewCpCode()
	cpCode.CpcodeID = cpCodeID.(string)
	err := cpCode.GetCpCode()
	if err != nil {
		return nil, err
	}

	log.Printf("[DEBUG] CP code found: %s\n", cpCode.CpcodeID)
	return cpCode, nil
}

func getProduct(d *schema.ResourceData, contract *papi.Contract, correlationid string) (*papi.Product, error) {
	if contract == nil {
		return nil, nil
	}

	log.Println("[DEBUG] Fetching product")
	productID, ok := d.GetOk("product")
	if !ok {
		return nil, nil
	}

	products := papi.NewProducts()
	e := products.GetProducts(contract, correlationid)
	if e != nil {
		return nil, e
	}

	product, e := products.FindProduct(productID.(string))
	if e != nil {
		return nil, e
	}

	log.Printf("[DEBUG] Product found: %s\n", product.ProductID)
	return product, nil
}

func createOrigin(d interface{}) (*papi.OptionValue, error) {
	log.Println("[DEBUG] Setting origin")
	var origin interface{}
	var ok bool

	switch d.(type) {
	case *schema.ResourceData:
		origin, ok = d.(*schema.ResourceData).GetOk("origin")
	case *schema.ResourceDiff:
		origin, ok = d.(*schema.ResourceDiff).GetOk("origin")
	default:
		origin, ok = d.(*schema.ResourceData).GetOk("origin")
	}

	//if origin, ok := rd.GetOk("origin"); ok {
	if ok {
		originConfig := origin.(*schema.Set).List()[0].(map[string]interface{})

		forwardHostname, forwardHostnameOk := originConfig["forward_hostname"].(string)
		originValues := make(map[string]interface{})

		originValues["originType"] = "CUSTOMER"
		if val, ok := originConfig["hostname"]; ok {
			originValues["hostname"] = val.(string)
		}

		if val, ok := originConfig["port"]; ok {
			originValues["httpPort"] = val.(int)
		}

		if val, ok := originConfig["cache_key_hostname"]; ok {
			originValues["cacheKeyHostname"] = val.(string)
		}

		if val, ok := originConfig["compress"]; ok {
			originValues["compress"] = val.(bool)
		}

		if val, ok := originConfig["enable_true_client_ip"]; ok {
			originValues["enableTrueClientIp"] = val.(bool)
		}

		if forwardHostnameOk && (forwardHostname == "ORIGIN_HOSTNAME" || forwardHostname == "REQUEST_HOST_HEADER") {
			log.Println("[DEBUG] Setting non-custom forward hostname")

			originValues["forwardHostHeader"] = forwardHostname
		} else if forwardHostnameOk {
			log.Println("[DEBUG] Setting custom forward hostname")

			originValues["forwardHostHeader"] = "CUSTOM"
			originValues["customForwardHostHeader"] = forwardHostname
		}

		ov := papi.OptionValue(originValues)
		return &ov, nil
	}
	return nil, nil
}

func fixupPerformanceBehaviors(rules *papi.Rules) {
	behavior, err := rules.FindBehavior("/Performance/sureRoute")
	if err != nil || behavior == nil || (behavior != nil && behavior.Options["testObjectUrl"] != "") {
		return
	}

	log.Println("[DEBUG] Fixing Up SureRoute Behavior")
	behavior.MergeOptions(papi.OptionValue{
		"testObjectUrl":   "/akamai/sureroute-testobject.html",
		"enableCustomKey": false,
		"enabled":         false,
	})
}

func fixupAdaptiveImageCompression(rules *papi.Rules) {
	log.Println("[DEBUG] Start Fixing Up adaptiveImageCompression Behavior")
	behavior, err := rules.FindBehavior("/Performance/JPEG Images/adaptiveImageCompression")
	log.Println("[DEBUG] Start Fixing Up adaptiveImageCompression Behavior ", behavior, err)
	if err != nil || behavior == nil { //} || (behavior != nil && behavior.Options["compressMobile"] != "") {
		log.Println("[DEBUG] Fixing Up adaptiveImageCompression Behavior Leave early")
		return
	}

	log.Println("[DEBUG] Fixing Up adaptiveImageCompression Behavior")
	behavior.MergeOptions(papi.OptionValue{
		"excellentConnectionOption": "Adapt Images",
		"goodConnectionOption":      "Adapt Images",
	})
	log.Println("[DEBUG] Start Fixing Up adaptiveImageCompression Behavior  ", behavior)
}

func updateStandardBehaviors(rules *papi.Rules, cpCode *papi.CpCode, origin *papi.OptionValue) {
	log.Printf("[DEBUG] cpCode: %#v", cpCode)
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
								log.Println("k:", k, "v:", v)
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

		// ADD is_secure from resource
		is_secure, set := d.GetOkExists("is_secure")
		if set && is_secure.(bool) {
			propertyRules.Rule.Options.IsSecure = true
		} else if set && !is_secure.(bool) {
			propertyRules.Rule.Options.IsSecure = false
		}

		// ADD cp_code from resource
		cp_code, set := d.GetOk("cp_code")
		if set {
			beh := papi.NewBehavior()
			beh.Name = "cpCode"
			beh.Options = papi.OptionValue{
				"value": papi.OptionValue{
					"id": cp_code.(string),
				},
			}
			propertyRules.Rule.MergeBehavior(beh)
		}
	}
}

func extractOptions(options *schema.Set) map[string]interface{} {
	optv := make(map[string]interface{})
	for _, o := range options.List() {
		oo, ok := o.(map[string]interface{})
		if ok {
			vals, ok := oo["values"]
			if ok && vals.(*schema.Set).Len() > 0 {
				op := make([]interface{}, 0)
				for _, v := range vals.(*schema.Set).List() {
					op = append(op, numberify(v.(string)))
				}

				optv[oo["key"].(string)] = op
			} else {
				optv[oo["key"].(string)] = numberify(oo["value"].(string))
			}
		}
	}
	return optv
}

func numberify(v string) interface{} {
	f1, err := strconv.ParseFloat(v, 64)
	if err == nil {
		return f1
	}

	f2, err := strconv.ParseInt(v, 10, 64)
	if err == nil {
		return f2
	}

	f3, err := strconv.ParseBool(v)
	if err == nil {
		return f3
	}

	f4, err := strconv.Atoi(v)
	if err == nil {
		return f4
	}

	return v
}

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
			log.Println("[DEBUG] extractRulesJSON Set criteriaMustSatisfy RESULT RULE value set " + string(rule.CriteriaMustSatisfy) + " " + rule.Name + " " + rule.Comments)

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

func extractRulesJSONDiff(d *schema.ResourceDiff, drules gjson.Result) []*papi.Rule {
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
			log.Println("[DEBUG] extractRulesJSON Set criteriaMustSatisfy RESULT RULE value set " + string(rule.CriteriaMustSatisfy) + " " + rule.Name + " " + rule.Comments)

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
					for _, newRule := range extractRulesJSONDiff(d, childRules) {
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

func extractRules(drules *schema.Set) []*papi.Rule {

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
							newBehavior.Options = extractOptions(behaviorOptions.(*schema.Set))
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
							newCriteria.Options = extractOptions(criteriaOptions.(*schema.Set))
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
				for _, newRule := range extractRules(childRules.(*schema.Set)) {
					rule.MergeChildRule(newRule)
				}
			}
		}
		rules = append(rules, rule)
	}

	return rules
}

func findProperty(d *schema.ResourceData, correlationid string) *papi.Property {
	results, err := papi.Search(papi.SearchByPropertyName, d.Get("name").(string), correlationid)
	if err != nil {
		return nil
	}

	if err != nil || results == nil {
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
		err = newVersion.Save(correlationid)
		if err != nil {
			return err
		}
	}

	return property.GetProperty(correlationid)
}
