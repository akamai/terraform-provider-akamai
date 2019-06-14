package akamai

import (
	"errors"
	"fmt"
	"log"
	"strconv"
	"strings"
	"sync"
	"time"

	"github.com/akamai/AkamaiOPEN-edgegrid-golang/jsonhooks-v1"
	"github.com/akamai/AkamaiOPEN-edgegrid-golang/papi-v1"
	"github.com/hashicorp/terraform/helper/schema"
	"github.com/tidwall/gjson"
)

func resourceProperty() *schema.Resource {
	return &schema.Resource{
		Create: resourcePropertyCreate,
		Read:   resourcePropertyRead,
		Update: resourcePropertyUpdate,
		Delete: resourcePropertyDelete,
		Exists: resourcePropertyExists,
		Importer: &schema.ResourceImporter{
			State: resourcePropertyImport,
		},
		Schema: akamaiPropertySchema,
	}
}

var akpsOption = &schema.Schema{
	Type:     schema.TypeSet,
	Optional: true,
	Elem: &schema.Resource{
		Schema: map[string]*schema.Schema{
			"key": {
				Type:     schema.TypeString,
				Required: true,
			},
			"values": {
				Type:     schema.TypeSet,
				Elem:     &schema.Schema{Type: schema.TypeString},
				Optional: true,
			},
			"value": {
				Type:     schema.TypeString,
				Optional: true,
			},
		},
	},
}

var akpsCriteria = &schema.Schema{
	Type:     schema.TypeSet,
	Optional: true,
	Elem: &schema.Resource{
		Schema: map[string]*schema.Schema{
			"name": {
				Type:     schema.TypeString,
				Required: true,
			},
			"option": akpsOption,
		},
	},
}

var akpsBehavior = &schema.Schema{
	Type:     schema.TypeSet,
	Optional: true,
	Elem: &schema.Resource{
		Schema: map[string]*schema.Schema{
			"name": {
				Type:     schema.TypeString,
				Required: true,
			},
			"option": akpsOption,
		},
	},
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
		ForceNew: true,
	},

	// Will get added to the default rule
	"cp_code": &schema.Schema{
		Type:     schema.TypeString,
		Required: true,
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
	"rule_format": &schema.Schema{
		Type:     schema.TypeString,
		Optional: true,
	},
	"ipv6": &schema.Schema{
		Type:     schema.TypeBool,
		Optional: true,
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
	"clone_from": &schema.Schema{
		Type:     schema.TypeSet,
		Optional: true,
		ForceNew: true,
		Elem: &schema.Resource{
			Schema: map[string]*schema.Schema{
				"property": {
					Type:     schema.TypeString,
					Required: true,
				},
				"version": {
					Type:     schema.TypeInt,
					Optional: true,
				},
				"etag": {
					Type:     schema.TypeString,
					Optional: true,
				},
				"copy_hostnames": {
					Type:     schema.TypeBool,
					Optional: true,
					Default:  false,
				},
			},
		},
	},

	// Will get added to the default rule
	"origin": {
		Type:     schema.TypeSet,
		Optional: true,
		Elem: &schema.Resource{
			Schema: map[string]*schema.Schema{
				"is_secure": {
					Type:     schema.TypeString,
					Required: true,
				},
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
	"rules": {
		Type:     schema.TypeString,
		Optional: true,
	},
	"variables": {
		Type:     schema.TypeString,
		Optional: true,
	},
}

func resourcePropertyCreate(d *schema.ResourceData, meta interface{}) error {
	d.Partial(true)

	group, e := getGroup(d)
	if e != nil {
		return e
	}

	contract, e := getContract(d)
	if e != nil {
		return e
	}

	cpCode, e := getCPCode(d, contract, group)
	if e != nil {
		return e
	}

	product, e := getProduct(d, contract)
	if e != nil {
		return e
	}

	cloneFrom, e := getCloneFrom(d)
	if e != nil {
		return e
	}

	var property *papi.Property
	if property = findProperty(d); property == nil {
		if group == nil {
			return errors.New("group must be specified to create a new property")
		}

		if contract == nil {
			return errors.New("contract must be specified to create a new property")
		}

		if product == nil {
			return errors.New("product must be specified to create a new property")
		}

		property, e = createProperty(contract, group, product, cloneFrom, d)
		if e != nil {
			return e
		}
	}

	err := ensureEditableVersion(property)
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

	rules := papi.NewRules()
	rules.PropertyID = d.Id()
	rules.PropertyVersion = property.LatestVersion

	origin, e := createOrigin(d)
	if e != nil {
		return e
	}

	updateStandardBehaviors(rules, cpCode, origin)
	fixupPerformanceBehaviors(rules)
	//fixupAdaptiveImageCompression(rules)

	// get rules from the TF config

	rulecheck, ok := d.GetOk("rules")
	if ok {
		log.Printf("[DEBUG] Unmarshal Rules from JSON")
		unmarshalRulesFromJSON(d, rules)
	}
	log.Printf("[DEBUG] Check for rules Json CREATE %s\n", rulecheck)

	e = rules.Save()
	if e != nil {
		if e == papi.ErrorMap[papi.ErrInvalidRules] && len(rules.Errors) > 0 {
			var msg string
			for _, v := range rules.Errors {
				msg = msg + fmt.Sprintf("\n Rule validation error: %s %s %s %s %s", v.Type, v.Title, v.Detail, v.Instance, v.BehaviorName)
			}
			return errors.New("Error - Invalid Property Rules" + msg)
		}
		return e
	}
	d.SetPartial("default")
	d.SetPartial("origin")

	hostnameEdgeHostnameMap, err := createHostnames(property, product, d)
	if err != nil {
		return err
	}

	edgeHostnames, err := setEdgeHostnames(property, hostnameEdgeHostnameMap, d)
	if err != nil {
		return err
	}
	d.SetPartial("hostname")
	d.SetPartial("ipv6")
	d.Set("edge_hostnames", edgeHostnames)

	d.Partial(false)
	log.Println("[DEBUG] Done")
	return nil
}

func createProperty(contract *papi.Contract, group *papi.Group, product *papi.Product, cloneFrom *papi.ClonePropertyFrom, d *schema.ResourceData) (*papi.Property, error) {
	log.Println("[DEBUG] Creating property")

	property, err := group.NewProperty(contract)
	if err != nil {
		return nil, err
	}

	property.ProductID = product.ProductID
	property.PropertyName = d.Get("name").(string)
	if cloneFrom != nil {
		property.CloneFrom = cloneFrom
	}

	if ruleFormat, ok := d.GetOk("rule_format"); ok {
		property.RuleFormat = ruleFormat.(string)
	} else {
		ruleFormats := papi.NewRuleFormats()
		property.RuleFormat, err = ruleFormats.GetLatest()
		if err != nil {
			return nil, err
		}
	}

	err = property.Save()
	if err != nil {
		return nil, err
	}

	log.Printf("[DEBUG] Property created: %s\n", property.PropertyID)
	return property, nil
}

func resourcePropertyDelete(d *schema.ResourceData, meta interface{}) error {
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

	e := property.GetProperty()
	if e != nil {
		return e
	}

	activations, e := property.GetActivations()
	if e != nil {
		return e
	}

	var (
		err = make(chan error, 2)
		wg  sync.WaitGroup
	)

	checkNetwork := func(network papi.NetworkValue) {
		defer wg.Done()

		var retried bool
	checkAgain:
		_, e = activations.GetLatestActivation(network, papi.StatusActive)
		if e == nil {
			if !retried {
				time.Sleep(time.Second * 30)
				retried = true
				goto checkAgain
			}
			err <- fmt.Errorf("property is still active on %s and cannot be deleted", network)
		}
	}

	wg.Add(2)
	go checkNetwork(papi.NetworkStaging)
	go checkNetwork(papi.NetworkProduction)

	// Bail as quickly as possible
	if e = <-err; e != nil {
		return e
	}

	wg.Wait()
	close(err)

	if e = <-err; e != nil {
		return e
	}

	e = property.Delete()
	if e != nil {
		return e
	}

	d.SetId("")

	log.Println("[DEBUG] Done")

	return nil
}

func resourcePropertyImport(d *schema.ResourceData, meta interface{}) ([]*schema.ResourceData, error) {
	resourceID := d.Id()
	propertyID := resourceID

	if !strings.HasPrefix(resourceID, "prp_") {
		for _, searchKey := range []papi.SearchKey{papi.SearchByPropertyName, papi.SearchByHostname, papi.SearchByEdgeHostname} {
			results, err := papi.Search(searchKey, resourceID)
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
	e := property.GetProperty()
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

func resourcePropertyExists(d *schema.ResourceData, meta interface{}) (bool, error) {
	property := papi.NewProperty(papi.NewProperties())
	property.PropertyID = d.Id()
	e := property.GetProperty()
	if e != nil {
		return false, e
	}

	return true, nil
}

func resourcePropertyRead(d *schema.ResourceData, meta interface{}) error {
	property := papi.NewProperty(papi.NewProperties())
	property.PropertyID = d.Id()
	err := property.GetProperty()
	if err != nil {
		return err
	}

	d.Set("account", property.AccountID)
	d.Set("contract", property.ContractID)
	d.Set("group", property.GroupID)
	d.Set("name", property.PropertyName)
	d.Set("note", property.Note)

	if ruleFormat, ok := d.GetOk("rule_format"); ok {
		d.Set("rule_format", ruleFormat.(string))
	} else {
		d.Set("rule_format", property.RuleFormat)
	}

	log.Printf("[DEBUG] Property RuleFormat from API : %s\n", property.RuleFormat)
	d.Set("version", property.LatestVersion)
	if property.StagingVersion > 0 {
		d.Set("staging_version", property.StagingVersion)
	}
	if property.ProductionVersion > 0 {
		d.Set("production_version", property.ProductionVersion)
	}

	return nil
}

func resourcePropertyUpdate(d *schema.ResourceData, meta interface{}) error {
	log.Printf("[DEBUG] UPDATING")
	d.Partial(true)

	property, e := getProperty(d)
	if e != nil {
		return e
	}

	err := ensureEditableVersion(property)
	if err != nil {
		return err
	}
	d.Set("version", property.LatestVersion)

	product, e := getProduct(d, property.Contract)
	if e != nil {
		return e
	}

	var cpCode *papi.CpCode
	if d.HasChange("cp_code") {
		cpCode, e = getCPCode(d, property.Contract, property.Group)
		if e != nil {
			return e
		}
		d.SetPartial("cp_code")
	} else {
		cpCode = papi.NewCpCode(papi.NewCpCodes(property.Contract, property.Group))
		cpCode.CpcodeID = d.Get("cp_code").(string)
		e := cpCode.GetCpCode()
		if e != nil {
			return e
		}
	}

	rules, e := property.GetRules()
	if e != nil {
		return e
	}

	if ruleFormat, ok := d.GetOk("rule_format"); ok {
		rules.RuleFormat = ruleFormat.(string)
	} else {
		ruleFormats := papi.NewRuleFormats()
		rules.RuleFormat, err = ruleFormats.GetLatest()
		if err != nil {
			return err
		}
	}

	origin, e := createOrigin(d)
	if e != nil {
		return e
	}

	updateStandardBehaviors(rules, cpCode, origin)
	//fixupAdaptiveImageCompression(rules)
	// get rules from the TF config

	rulecheck, ok := d.GetOk("rules")
	if ok {
		log.Printf("[DEBUG] Unmarshal Rules from JSON")
		unmarshalRulesFromJSON(d, rules)
	}
	log.Printf("[DEBUG] Check for rules Json Before Unmarshal UPDATE %s\n", rulecheck)

	jsonBody, err := jsonhooks.Marshal(rules)
	if err != nil {
		return err
	}
	log.Printf("[DEBUG] Check rules after unmarshal from Json %s\n", string(jsonBody))

	e = rules.Save()
	if e != nil {
		if e == papi.ErrorMap[papi.ErrInvalidRules] && len(rules.Errors) > 0 {
			var msg string
			for _, v := range rules.Errors {
				msg = msg + fmt.Sprintf("\n Rule validation error: %s %s %s %s %s", v.Type, v.Title, v.Detail, v.Instance, v.BehaviorName)
			}
			return errors.New("Error - Invalid Property Rules" + msg)
		}
		return e
	}
	d.SetPartial("default")
	d.SetPartial("origin")

	if d.HasChange("ipv6") || d.HasChange("hostnames") {
		hostnameEdgeHostnameMap, err := createHostnames(property, product, d)
		if err != nil {
			return err
		}

		edgeHostnames, err := setEdgeHostnames(property, hostnameEdgeHostnameMap, d)
		if err != nil {
			return err
		}

		d.SetPartial("ipv6")
		d.Set("edge_hostnames", edgeHostnames)
	}

	d.Partial(false)

	log.Println("[DEBUG] Done")
	return nil
}

// Helpers
func getProperty(d *schema.ResourceData) (*papi.Property, error) {
	log.Println("[DEBUG] Fetching property")
	propertyID := d.Id()
	property := papi.NewProperty(papi.NewProperties())
	property.PropertyID = propertyID
	e := property.GetProperty()
	return property, e
}

func getGroup(d *schema.ResourceData) (*papi.Group, error) {
	log.Println("[DEBUG] Fetching groups")
	groupID, ok := d.GetOk("group")

	if !ok {
		return nil, nil
	}

	groups := papi.NewGroups()
	e := groups.GetGroups()
	if e != nil {
		return nil, e
	}

	group, e := groups.FindGroup(groupID.(string))
	if e != nil {
		return nil, e
	}

	log.Printf("[DEBUG] Group found: %s\n", group.GroupID)
	return group, nil
}

func getContract(d *schema.ResourceData) (*papi.Contract, error) {
	log.Println("[DEBUG] Fetching contract")
	contractID, ok := d.GetOk("contract")
	if !ok {
		return nil, nil
	}

	contracts := papi.NewContracts()
	e := contracts.GetContracts()
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

func getCPCode(d *schema.ResourceData, contract *papi.Contract, group *papi.Group) (*papi.CpCode, error) {
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

func getProduct(d *schema.ResourceData, contract *papi.Contract) (*papi.Product, error) {
	if contract == nil {
		return nil, nil
	}

	log.Println("[DEBUG] Fetching product")
	productID, ok := d.GetOk("product")
	if !ok {
		return nil, nil
	}

	products := papi.NewProducts()
	e := products.GetProducts(contract)
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

func getCloneFrom(d *schema.ResourceData) (*papi.ClonePropertyFrom, error) {
	log.Println("[DEBUG] Setting up clone from")

	cF, ok := d.GetOk("clone_from")

	if !ok {
		return nil, nil
	}

	set := cF.(*schema.Set)
	cloneFrom := set.List()[0].(map[string]interface{})

	propertyID := cloneFrom["property"].(string)

	property := papi.NewProperty(papi.NewProperties())
	property.PropertyID = propertyID
	err := property.GetProperty()
	if err != nil {
		return nil, err
	}

	version := cloneFrom["version"].(int)

	if cloneFrom["version"].(int) == 0 {
		v, err := property.GetLatestVersion("")
		if err != nil {
			return nil, err
		}
		version = v.PropertyVersion
	}

	clone := papi.NewClonePropertyFrom()
	clone.PropertyID = propertyID
	clone.Version = version

	if cloneFrom["etag"].(string) != "" {
		clone.CloneFromVersionEtag = cloneFrom["etag"].(string)
	}

	if cloneFrom["copy_hostnames"].(bool) != false {
		clone.CopyHostnames = true
	}

	log.Println("[DEBUG] Clone from complete")

	return clone, nil
}

func createOrigin(d *schema.ResourceData) (*papi.OptionValue, error) {
	log.Println("[DEBUG] Setting origin")
	if origin, ok := d.GetOk("origin"); ok {
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
			originValues["customForwardHostHeader"] = "CUSTOM"
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

func createHostnames(property *papi.Property, product *papi.Product, d *schema.ResourceData) (map[string]*papi.EdgeHostname, error) {
	// If the property has edge hostnames and none is specified in the schema, then don't update them
	edgeHostname, edgeHostnameOk := d.GetOk("edge_hostnames")
	if edgeHostnameOk == false {
		hostnames, err := property.GetHostnames(nil)
		if err != nil {
			return nil, err
		}

		if len(hostnames.Hostnames.Items) > 0 {
			return nil, nil
		}
	}

	log.Println("[DEBUG] Check for Edge Host Map for hostnames")
	edgeHostnameMapTest, ok := d.GetOk("hostnames")
	var hostname string
	var hostCount int

	if ok {

		log.Println("[DEBUG] Check for Edge Host Map for hostnames ", edgeHostnameMapTest)

		for key, value := range edgeHostnameMapTest.(map[string]interface{}) {

			log.Println("[DEBUG] EdgeHostnameMap ", key, value)
			if strings.Contains(value.(string), "ehn_") {
			}
			log.Println("[DEBUG] EdgeHostnameMap ", key, value)
			hostCount = hostCount + 1
		}

	}

	if hostCount == 0 {
		hostCount = 1
	}

	hostnames := make([]string, 0, hostCount)

	if ok {
		for key, value := range edgeHostnameMapTest.(map[string]interface{}) {
			log.Println("[DEBUG] EdgeHostnameMap ", key, value)

			if strings.Contains(key, "cnamefrom") {
				hostname = value.(string)
				log.Println("[DEBUG] Figuring out edgehostname ", hostname)
				hostnames = append(hostnames, hostname)
			}
		}
	}

	ipv6 := d.Get("ipv6").(bool)

	log.Println("[DEBUG] Figuring out hostnames")
	edgeHostnames := papi.NewEdgeHostnames()
	err := edgeHostnames.GetEdgeHostnames(property.Contract, property.Group, "")
	if err != nil {
		return nil, err
	}

	hostnameEdgeHostnameMap := map[string]*papi.EdgeHostname{}
	defaultEdgeHostname := edgeHostnames.EdgeHostnames.Items[0]

	if edgeHostnameOk {
		foundEdgeHostname := false
		for _, eHn := range edgeHostnames.EdgeHostnames.Items {
			if eHn.EdgeHostnameDomain == edgeHostname.(string) {
				foundEdgeHostname = true
				defaultEdgeHostname = eHn
			}
		}

		if foundEdgeHostname == false {
			var err error
			defaultEdgeHostname, err = createEdgehostname(edgeHostnames, product, edgeHostname.(string), ipv6)
			if err != nil {
				return nil, err
			}
		}

		for _, hostname := range hostnames {
			if _, ok := hostnameEdgeHostnameMap[hostname]; !ok {
				hostnameEdgeHostnameMap[hostname] = defaultEdgeHostname
				return hostnameEdgeHostnameMap, nil
			}
		}
	}

	// Contract/Group has _some_ Edge Hostnames, try to map 1:1 (e.g. example.com -> example.com.edgesuite.net)
	// If some mapping exists, map non-existent ones to the first 1:1 we find, otherwise if none exist map to the
	// first Edge Hostname found in the contract/group
	if len(edgeHostnames.EdgeHostnames.Items) > 0 {
		log.Println("[DEBUG] Hostnames retrieved, trying to map")
		edgeHostnamesMap := map[string]*papi.EdgeHostname{}

		for _, edgeHostname := range edgeHostnames.EdgeHostnames.Items {
			edgeHostnamesMap[edgeHostname.EdgeHostnameDomain] = edgeHostname
		}

		// Search for existing hostname, map 1:1
		var overrideDefault bool
		for _, hostname := range hostnames {
			if edgeHostname, ok := edgeHostnamesMap[hostname+".edgesuite.net"]; ok {
				hostnameEdgeHostnameMap[hostname] = edgeHostname
				// Override the default with the first one found
				if !overrideDefault {
					defaultEdgeHostname = edgeHostname
					overrideDefault = true
				}
				continue
			}
			if edgeHostname, ok := edgeHostnamesMap[hostname+".edgekey.net"]; ok {
				hostnameEdgeHostnameMap[hostname] = edgeHostname
				// Override the default with the first one found
				if !overrideDefault {
					defaultEdgeHostname = edgeHostname
					overrideDefault = true
				}
				continue
			}
		}

		// Fill in defaults
		if len(hostnameEdgeHostnameMap) < len(hostnames) {
			log.Printf("[DEBUG] Hostnames being set to default: %d of %d\n", len(hostnameEdgeHostnameMap), len(hostnames))
			for _, hostname := range hostnames {
				if _, ok := hostnameEdgeHostnameMap[hostname]; !ok {
					hostnameEdgeHostnameMap[hostname] = defaultEdgeHostname
				}
			}
		}
	}

	// Contract/Group has no Edge Hostnames, create a single based on the first hostname
	// mapping example.com -> example.com.edgegrid.net
	if len(edgeHostnames.EdgeHostnames.Items) == 0 {
		log.Println("[DEBUG] No Edge Hostnames found, creating new one")
		newEdgeHostname, err := createEdgehostname(edgeHostnames, product, hostnames[0], ipv6)
		if err != nil {
			return nil, err
		}

		for _, hostname := range hostnames {
			hostnameEdgeHostnameMap[hostname] = newEdgeHostname
		}

		log.Printf("[DEBUG] Edgehostname created: %s\n", newEdgeHostname.EdgeHostnameDomain)
	}

	return hostnameEdgeHostnameMap, nil
}

func createEdgehostname(edgeHostnames *papi.EdgeHostnames, product *papi.Product, hostname string, ipv6 bool) (*papi.EdgeHostname, error) {
	newEdgeHostname := papi.NewEdgeHostname(edgeHostnames)
	newEdgeHostname.ProductID = product.ProductID
	newEdgeHostname.IPVersionBehavior = "IPV4"
	if ipv6 {
		newEdgeHostname.IPVersionBehavior = "IPV6_COMPLIANCE"
	}

	newEdgeHostname.EdgeHostnameDomain = hostname
	err := newEdgeHostname.Save("")
	if err != nil {
		return nil, err
	}

	go newEdgeHostname.PollStatus("")

	for newEdgeHostname.Status != papi.StatusActive {
		select {
		case <-newEdgeHostname.StatusChange:
		case <-time.After(time.Minute * 20):
			return nil, fmt.Errorf("no edge hostname found and a timeout occurred trying to create \"%s.%s\"", newEdgeHostname.DomainPrefix, newEdgeHostname.DomainSuffix)
		}
	}

	return newEdgeHostname, nil
}

func setEdgeHostnames(property *papi.Property, hostnameEdgeHostnameMap map[string]*papi.EdgeHostname, d *schema.ResourceData) (map[string]string, error) {
	if hostnameEdgeHostnameMap != nil {
		log.Println("[DEBUG] Setting Edge Hostnames")
		propertyHostnames, err := property.GetHostnames(nil)
		if err != nil {
			return nil, err
		}

		propertyHostnames.Hostnames.Items = []*papi.Hostname{}
		for from, to := range hostnameEdgeHostnameMap {

			log.Println("[DEBUG] Check for Edge Host Map for replacing hostname")
			edgeHostnameMapTest, ok := d.GetOk("hostnames")
			if ok {

				log.Println("[DEBUG] Check for Edge Host Map for hostnames ", edgeHostnameMapTest)
				mapString := make(map[string]string)
				var edgeHostList []string

				for key, value := range edgeHostnameMapTest.(map[string]interface{}) {
					log.Println("[DEBUG] EdgeHostnameMap ", key, value)
					strKey := fmt.Sprintf("%v", key)
					strValue := fmt.Sprintf("%v", value)
					mapString[strKey] = strValue
					if strings.Contains(strValue, "ehn_") {
						edgeHostList = append(edgeHostList, strValue)

						log.Println("[DEBUG] EdgeHostnameMap set EHN_KEY ", key, value)
					}
				}

				for _, value := range edgeHostList {
					hostname := propertyHostnames.NewHostname()

					log.Println("[DEBUG] edgeHostList ", value)

					strValue := fmt.Sprintf("%v", value)

					if strings.Contains(strValue, "ehn_") {
						mapString["ehn_Key"] = strValue
						log.Println("[DEBUG] edgeHostList set EHN_KEY ", value)
					}

					log.Println("[DEBUG] Check for Edge Host Map for cnamefrom ", mapString[value+"-cnamefrom"])
					log.Println("[DEBUG] Check for Edge Host Map for cnameto ", mapString[value+"-cnameto"])
					hostname.CnameType = papi.CnameTypeEdgeHostname
					hostname.CnameFrom = mapString[value+"-cnamefrom"]
					hostname.CnameTo = mapString[value+"-cnameto"]

					hostname.CertEnrollmentId = mapString[value+"-certenrollmentid"]

				}

			} else {
				log.Println("[DEBUG] Nor Edge Host Map for cnamefrom use existing ")
				hostname := propertyHostnames.NewHostname()
				hostname.CnameType = papi.CnameTypeEdgeHostname
				hostname.CnameFrom = from
				hostname.CnameTo = to.EdgeHostnameDomain
				hostname.EdgeHostnameID = to.EdgeHostnameID
			}
		}

		log.Println("[DEBUG] Saving edge hostnames", propertyHostnames)

		err = propertyHostnames.Save()
		log.Println("[DEBUG] Edge hostnames saved")
		if err != nil {
			return nil, err
		}
	}

	hostnames, err := property.GetHostnames(nil)
	if err != nil {
		return nil, err
	}

	var ehn = make(map[string]string)
	for _, hostname := range hostnames.Hostnames.Items {
		ehn[strings.Replace(hostname.CnameFrom, ".", "-", -1)] = hostname.CnameTo
	}

	return ehn, nil
}

func unmarshalRulesFromJSON(d *schema.ResourceData, propertyRules *papi.Rules) {
	// Default Rules
	rules, ok := d.GetOk("rules")
	if ok {

		log.Println("[DEBUG] RulesJson")
		rulesJSON := gjson.Get(rules.(string), "rules")

		rulesJSON.ForEach(func(key, value gjson.Result) bool {
			log.Println("[DEBUG] unmarshalRulesFromJson KEY RULES KEY = " + key.String() + " VAL " + value.String())

			if key.String() == "behaviors" {
				behavior := gjson.Parse(value.String())
				log.Println("[DEBUG] unmarshalRulesFromJson KEY BEHAVIOR " + behavior.String())
				if gjson.Get(behavior.String(), "#.name").Exists() {

					behavior.ForEach(func(key, value gjson.Result) bool {
						log.Println("[DEBUG] unmarshalRulesFromJson BEHAVIOR LOOP KEY =" + key.String() + " VAL " + value.String())

						bb, ok := value.Value().(map[string]interface{})
						if ok {
							log.Println("[DEBUG] unmarshalRulesFromJson BEHAVIOR MAP  ", bb)
							for k, v := range bb {
								log.Println("k:", k, "v:", v)
							}

							beh := papi.NewBehavior()

							beh.Name = bb["name"].(string)
							boptions, ok := bb["options"]
							log.Println("[DEBUG] unmarshalRulesFromJson KEY BEHAVIOR BOPTIONS ", boptions)
							if ok {
								beh.Options = boptions.(map[string]interface{})
								log.Println("[DEBUG] unmarshalRulesFromJson KEY BEHAVIOR EXTRACT BOPTIONS ", beh.Options)
							}

							// Fixup CPCode
							if beh.Name == "cpCode" {
								cpCode := papi.NewCpCodes(nil, nil).NewCpCode()
								cpCode.CpcodeID = beh.Options["id"].(string)
								beh.Options = papi.OptionValue{"value": papi.OptionValue{"id": cpCode.ID()}}
							}

							propertyRules.Rule.MergeBehavior(beh)
						}

						return true // keep iterating
					}) // behavior list loop

				}

				if key.String() == "criteria" {
					criteria := gjson.Parse(value.String())

					criteria.ForEach(func(key, value gjson.Result) bool {
						log.Println("[DEBUG] unmarshalRulesFromJson KEY CRITERIA " + key.String() + " VAL " + value.String())

						cc, ok := value.Value().(map[string]interface{})
						if ok {
							log.Println("[DEBUG] unmarshalRulesFromJson CRITERIA MAP  ", cc)
							newCriteria := papi.NewCriteria()
							newCriteria.Name = cc["name"].(string)

							coptions, ok := cc["option"]
							if ok {
								println("OPTIONS ", coptions)
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
				println("CHILD RULES " + childRules.String())

				for _, rule := range extractRulesJSON(d, childRules) {
					propertyRules.Rule.MergeChildRule(rule)
				}
			}

			if key.String() == "variables" {

				log.Println("unmarshalRulesFromJson VARS from JSON ", value.String())
				variables := gjson.Parse(value.String())

				variables.ForEach(func(key, value gjson.Result) bool {
					log.Println("unmarshalRulesFromJson VARS from JSON LOOP ", value)
					variableMap, ok := value.Value().(map[string]interface{})
					log.Println("unmarshalRulesFromJson VARS from JSON LOOP NAME ", variableMap["name"].(string))
					log.Println("unmarshalRulesFromJson VARS from JSON LOOP DESC ", variableMap["description"].(string))
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

			return true // keep iterating
		}) // for loop rules

		// ADD vars from variables resource
		jsonvars, ok := d.GetOk("variables")
		if ok {
			log.Println("unmarshalRulesFromJson VARS from JSON ", jsonvars)
			variables := gjson.Parse(jsonvars.(string))
			result := gjson.Get(variables.String(), "variables")
			log.Println("unmarshalRulesFromJson VARS from JSON VARIABLES ", result)

			result.ForEach(func(key, value gjson.Result) bool {
				log.Println("unmarshalRulesFromJson VARS from JSON LOOP ", value)
				variableMap, ok := value.Value().(map[string]interface{})
				log.Println("unmarshalRulesFromJson VARS from JSON LOOP NAME ", variableMap["name"].(string))
				log.Println("unmarshalRulesFromJson VARS from JSON LOOP DESC ", variableMap["description"].(string))
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

	}

}

func unmarshalRules(d *schema.ResourceData, propertyRules *papi.Rules) {
	// Default Rules
	rules, ok := d.GetOk("rules")
	if ok {
		for _, r := range rules.(*schema.Set).List() {
			ruleTree, ok := r.(map[string]interface{})
			if ok {
				behavior, ok := ruleTree["behavior"]
				if ok {
					for _, b := range behavior.(*schema.Set).List() {
						bb, ok := b.(map[string]interface{})
						if ok {
							beh := papi.NewBehavior()
							beh.Name = bb["name"].(string)
							boptions, ok := bb["option"]
							if ok {
								beh.Options = extractOptions(boptions.(*schema.Set))
							}
							propertyRules.Rule.MergeBehavior(beh)
						}
					}
				}

				criteria, ok := ruleTree["criteria"]
				if ok {
					for _, c := range criteria.(*schema.Set).List() {
						cc, ok := c.(map[string]interface{})
						if ok {
							newCriteria := papi.NewCriteria()
							newCriteria.Name = cc["name"].(string)
							coptions, ok := cc["option"]
							if ok {
								newCriteria.Options = extractOptions(coptions.(*schema.Set))
							}
							propertyRules.Rule.MergeCriteria(newCriteria)
						}
					}
				}
			}

			childRules, ok := ruleTree["rule"]
			if ok {
				for _, rule := range extractRules(childRules.(*schema.Set)) {
					propertyRules.Rule.MergeChildRule(rule)
				}
			}
		}

		// ADD vars from variables resource
		jsonvars, ok := d.GetOk("variables")
		if ok {
			log.Println("VARS from JSON ", jsonvars)
			variables := gjson.Parse(jsonvars.(string))
			result := gjson.Get(variables.String(), "variables") //gjson.GetMany(variables.String(),"variables.#.name","variables.#.description","variables.#.value","variables.#.hidden","variables.#.sensitive" )

			result.ForEach(func(key, value gjson.Result) bool {
				variableMap, ok := value.Value().(map[string]interface{})
				log.Println("VARS from JSON LOOP NAME ", variableMap["name"].(string))
				log.Println("VARS from JSON LOOP DESC ", variableMap["description"].(string))
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

func extractRulesJSON(d *schema.ResourceData, drules gjson.Result) []*papi.Rule {
	var rules []*papi.Rule
	drules.ForEach(func(key, value gjson.Result) bool {
		rule := papi.NewRule()
		vv, ok := value.Value().(map[string]interface{})
		if ok {
			rule.Name, _ = vv["name"].(string)
			rule.Comments, _ = vv["comments"].(string)

			ruledetail := gjson.Parse(value.String())
			log.Println("[DEBUG] RULE DETAILS ", ruledetail)

			ruledetail.ForEach(func(key, value gjson.Result) bool {

				if key.String() == "behaviors" {
					log.Println("[DEBUG] BEHAVIORS KEY CHILD RULE ", key.String())

					behaviors := gjson.Parse(value.String())
					log.Println("[DEBUG] BEHAVIORS NAME ", behaviors)
					behaviors.ForEach(func(key, value gjson.Result) bool {
						log.Println("[DEBUG] BEHAVIORS KEY CHILD RULE LOOP KEY = " + key.String() + " VAL " + value.String())
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
					log.Println("[DEBUG] CRITERIA KEY CHILD RULE ", key.String())
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
					log.Println("[DEBUG] VARIABLES KEY CHILD RULE ", key.String())
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
					println("CHILD RULES " + childRules.String())
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

func extractRules(drules *schema.Set) []*papi.Rule {
	var rules []*papi.Rule
	for _, v := range drules.List() {
		rule := papi.NewRule()
		vv, ok := v.(map[string]interface{})
		if ok {
			rule.Name = vv["name"].(string)
			rule.Comments = vv["comment"].(string)
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

func findProperty(d *schema.ResourceData) *papi.Property {
	results, err := papi.Search(papi.SearchByPropertyName, d.Get("name").(string))
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

	err = property.GetProperty()
	if err != nil {
		return nil
	}

	return property
}

func ensureEditableVersion(property *papi.Property) error {
	latestVersion, err := property.GetLatestVersion("")
	if err != nil {
		return err
	}

	versions, err := property.GetVersions()
	if err != nil {
		return err
	}

	if latestVersion.ProductionStatus != papi.StatusInactive || latestVersion.StagingStatus != papi.StatusInactive {
		// The latest version has been activated on either production or staging, so we need to create a new version to apply changes on
		newVersion := versions.NewVersion(latestVersion, false)
		err = newVersion.Save()
		if err != nil {
			return err
		}
	}

	return property.GetProperty()
}
