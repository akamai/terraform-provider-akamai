package akamai

import (
	"bytes"
	"encoding/json"
	"log"

	"github.com/akamai/AkamaiOPEN-edgegrid-golang/jsonhooks-v1"
	"github.com/akamai/AkamaiOPEN-edgegrid-golang/papi-v1"
	"github.com/hashicorp/terraform/helper/schema"
	"github.com/tidwall/gjson"
)

// suppressEquivalentTypeStringBoolean provides custom difference suppression for TypeString booleans
// Some arguments require three values: true, false, and "" (unspecified), but
// confusing behavior exists when converting bare true/false values with state.
func suppressEquivalentTypeStringBoolean(k, old, new string, d *schema.ResourceData) bool {
	if old == "false" && new == "0" {
		return true
	}
	if old == "true" && new == "1" {
		return true
	}
	return false
}

func suppressEquivalentJsonDiffs(k, old, new string, d *schema.ResourceData) bool {
	ob := bytes.NewBufferString("")
	if err := json.Compact(ob, []byte(old)); err != nil {
		return false
	}

	nb := bytes.NewBufferString("")
	if err := json.Compact(nb, []byte(new)); err != nil {
		return false
	}
	log.Printf("[DEBUG] suppressEquivalentJsonDiffs OB %s\n", string(ob.Bytes()))
	log.Printf("[DEBUG] suppressEquivalentJsonDiffs NB %s\n", string(nb.Bytes()))

	rulesOld, err := getRulesForComp(d, old)
	rulesOld.Etag = ""
	jsonBody, err := jsonhooks.Marshal(rulesOld)
	if err != nil {
		return false
	}
	sha1hashOld := getSHAString(string(jsonBody))

	log.Printf("[DEBUG] suppressEquivalentJsonDiffs SHA from OLD Json %s\n", sha1hashOld)

	rulesNew, err := getRulesForComp(d, new)
	rulesNew.Etag = ""
	jsonBodyNew, err := jsonhooks.Marshal(rulesNew)
	if err != nil {
		return false
	}
	sha1hashNew := getSHAString(string(jsonBodyNew))

	log.Printf("[DEBUG] suppressEquivalentJsonDiffs SHA from NEW Json %s\n", sha1hashNew)

	if sha1hashOld == sha1hashNew {
		log.Printf("[DEBUG] suppressEquivalentJsonDiffs SHA Equal skip diff \n")
		return true
	} else {
		log.Printf("[DEBUG] suppressEquivalentJsonDiffs SHA Not Equal diff applies \n")
		return false
	}

}

func suppressEquivalentJsonPendingDiffs(old, new string, d *schema.ResourceDiff) bool {
	ob := bytes.NewBufferString("")
	if err := json.Compact(ob, []byte(old)); err != nil {
		return false
	}

	nb := bytes.NewBufferString("")
	if err := json.Compact(nb, []byte(new)); err != nil {
		return false
	}
	log.Printf("[DEBUG] suppressEquivalentJsonDiffs OB %s\n", string(ob.Bytes()))
	log.Printf("[DEBUG] suppressEquivalentJsonDiffs NB %s\n", string(nb.Bytes()))

	rulesOld, err := getRulesForComp(d, old)
	rulesOld.Etag = ""
	jsonBody, err := jsonhooks.Marshal(rulesOld)
	if err != nil {
		return false
	}
	sha1hashOld := getSHAString(string(jsonBody))

	log.Printf("[DEBUG] suppressEquivalentJsonDiffs SHA from OLD Json %s\n", sha1hashOld)

	rulesNew, err := getRulesForComp(d, new)
	rulesNew.Etag = ""
	jsonBodyNew, err := jsonhooks.Marshal(rulesNew)
	if err != nil {
		return false
	}
	sha1hashNew := getSHAString(string(jsonBodyNew))

	log.Printf("[DEBUG] suppressEquivalentJsonDiffs SHA from NEW Json %s\n", sha1hashNew)

	if sha1hashOld == sha1hashNew {
		log.Printf("[DEBUG] suppressEquivalentJsonDiffs SHA Equal skip diff \n")
		return true
	} else {
		log.Printf("[DEBUG] suppressEquivalentJsonDiffs SHA Not Equal diff applies \n")
		return false
	}

}

func getRulesForComp(d interface{}, json string) (*papi.Rules, error) {

	property, e := getProperty(d)
	if e != nil {
		return nil, e
	}

	rules := papi.NewRules()
	rules.Rule.Name = "default"
	switch d.(type) {
	case *schema.ResourceData:
		rules.PropertyID = d.(*schema.ResourceData).Id()
	case *schema.ResourceDiff:
		rules.PropertyID = d.(*schema.ResourceDiff).Id()
	default:
		rules.PropertyID = d.(*schema.ResourceData).Id()
	}

	//rules.PropertyID = d.Id()
	rules.PropertyVersion = property.LatestVersion

	origin, err := createOrigin(d)
	if err != nil {
		return nil, err
	}

	// get rules from the TF config

	//rulecheck

	log.Printf("[DEBUG] Unmarshal Rules from JSON")
	unmarshalRulesFromJSONComp(d, json, rules)

	var ruleFormat interface{}
	var ok bool

	switch d.(type) {
	case *schema.ResourceData:
		ruleFormat, ok = d.(*schema.ResourceData).GetOk("rule_format")
	case *schema.ResourceDiff:
		ruleFormat, ok = d.(*schema.ResourceDiff).GetOk("rule_format")
	default:
		ruleFormat, ok = d.(*schema.ResourceData).GetOk("rule_format")
	}

	//if ruleFormat, ok := d.GetOk("rule_format"); ok {
	if ok {
		rules.RuleFormat = ruleFormat.(string)
	} else {
		ruleFormats := papi.NewRuleFormats()
		rules.RuleFormat, err = ruleFormats.GetLatest()
		if err != nil {
			return nil, err
		}
	}

	//if ok := d.HasChange("rule_format"); ok {
	//}

	cpCode, err := getCPCode(d, property.Contract, property.Group)
	if err != nil {
		return nil, err
	}

	log.Printf("[DEBUG] updateStandardBehaviors")
	updateStandardBehaviors(rules, cpCode, origin)
	log.Printf("[DEBUG] fixupPerformanceBehaviors")
	fixupPerformanceBehaviors(rules)

	return rules, nil
}

func unmarshalRulesFromJSONComp(d interface{}, rulesComp string, propertyRules *papi.Rules) {
	// Default Rules

	propertyRules.Rule = &papi.Rule{Name: "default"}
	//log.Println("[DEBUG] RulesJson")

	rulesJSON := gjson.Parse(rulesComp).Get("rules")
	rulesJSON.ForEach(func(key, value gjson.Result) bool {
		//	log.Println("[DEBUG] unmarshalRulesFromJson KEY RULES KEY = " + key.String() + " VAL " + value.String())

		if key.String() == "behaviors" {
			behavior := gjson.Parse(value.String())
			//		log.Println("[DEBUG] unmarshalRulesFromJson KEY BEHAVIOR " + behavior.String())
			if gjson.Get(behavior.String(), "#.name").Exists() {

				behavior.ForEach(func(key, value gjson.Result) bool {
					//				log.Println("[DEBUG] unmarshalRulesFromJson BEHAVIOR LOOP KEY =" + key.String() + " VAL " + value.String())

					bb, ok := value.Value().(map[string]interface{})
					if ok {
						//					log.Println("[DEBUG] unmarshalRulesFromJson BEHAVIOR MAP  ", bb)
						for k, v := range bb {
							log.Println("k:", k, "v:", v)
						}

						beh := papi.NewBehavior()

						beh.Name = bb["name"].(string)
						boptions, ok := bb["options"]
						//					log.Println("[DEBUG] unmarshalRulesFromJson KEY BEHAVIOR BOPTIONS ", boptions)
						if ok {
							beh.Options = boptions.(map[string]interface{})
							//						log.Println("[DEBUG] unmarshalRulesFromJson KEY BEHAVIOR EXTRACT BOPTIONS ", beh.Options)
						}

						propertyRules.Rule.MergeBehavior(beh)
					}

					return true // keep iterating
				}) // behavior list loop

			}

			if key.String() == "criteria" {
				criteria := gjson.Parse(value.String())

				criteria.ForEach(func(key, value gjson.Result) bool {
					//				log.Println("[DEBUG] unmarshalRulesFromJson KEY CRITERIA " + key.String() + " VAL " + value.String())

					cc, ok := value.Value().(map[string]interface{})
					if ok {
						//					log.Println("[DEBUG] unmarshalRulesFromJson CRITERIA MAP  ", cc)
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
			//		println("CHILD RULES " + childRules.String())

			for _, rule := range extractRulesJSON(d, childRules) {
				propertyRules.Rule.MergeChildRule(rule)
			}
		}

		if key.String() == "variables" {

			//		log.Println("unmarshalRulesFromJson VARS from JSON ", value.String())
			variables := gjson.Parse(value.String())

			variables.ForEach(func(key, value gjson.Result) bool {
				//			log.Println("unmarshalRulesFromJson VARS from JSON LOOP ", value)
				variableMap, ok := value.Value().(map[string]interface{})
				//			log.Println("unmarshalRulesFromJson VARS from JSON LOOP NAME ", variableMap["name"].(string))
				//			log.Println("unmarshalRulesFromJson VARS from JSON LOOP DESC ", variableMap["description"].(string))
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
			//		log.Println("unmarshalRulesFromJson OPTIONS from JSON", value.String())
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
	var jsonvars interface{}
	var ok bool

	switch d.(type) {
	case *schema.ResourceData:
		jsonvars, ok = d.(*schema.ResourceData).GetOk("variables")
	case *schema.ResourceDiff:
		jsonvars, ok = d.(*schema.ResourceDiff).GetOk("variables")
	default:
		jsonvars, ok = d.(*schema.ResourceData).GetOk("variables")
	}
	//jsonvars, ok := d.GetOk("variables")
	if ok {
		//	log.Println("unmarshalRulesFromJson VARS from JSON ", jsonvars)
		variables := gjson.Parse(jsonvars.(string))
		result := gjson.Get(variables.String(), "variables")
		//	log.Println("unmarshalRulesFromJson VARS from JSON VARIABLES ", result)

		result.ForEach(func(key, value gjson.Result) bool {
			//		log.Println("unmarshalRulesFromJson VARS from JSON LOOP ", value)
			variableMap, ok := value.Value().(map[string]interface{})
			//		log.Println("unmarshalRulesFromJson VARS from JSON LOOP NAME ", variableMap["name"].(string))
			//		log.Println("unmarshalRulesFromJson VARS from JSON LOOP DESC ", variableMap["description"].(string))
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
	var is_secure interface{}
	var set bool

	switch d.(type) {
	case *schema.ResourceData:
		is_secure, set = d.(*schema.ResourceData).GetOk("is_secure")
	case *schema.ResourceDiff:
		is_secure, set = d.(*schema.ResourceDiff).GetOk("is_secure")
	default:
		is_secure, set = d.(*schema.ResourceData).GetOk("is_secure")
	}

	//is_secure, set := d.GetOkExists("is_secure")
	if set && is_secure.(bool) {
		propertyRules.Rule.Options.IsSecure = true
	} else if set && !is_secure.(bool) {
		propertyRules.Rule.Options.IsSecure = false
	}

	// ADD cp_code from resource
	var cp_code interface{}

	switch d.(type) {
	case *schema.ResourceData:
		cp_code, set = d.(*schema.ResourceData).GetOk("cp_code")
	case *schema.ResourceDiff:
		cp_code, set = d.(*schema.ResourceDiff).GetOk("cp_code")
	default:
		cp_code, set = d.(*schema.ResourceData).GetOk("cp_code")
	}

	//cp_code, set := d.GetOk("cp_code")
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
