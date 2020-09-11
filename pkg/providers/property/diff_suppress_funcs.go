package property

import (
	"bytes"
	"encoding/json"
	"errors"
	"fmt"

	"github.com/apex/log"

	"github.com/akamai/terraform-provider-akamai/v2/pkg/akamai"
	"github.com/akamai/terraform-provider-akamai/v2/pkg/tools"

	"github.com/akamai/AkamaiOPEN-edgegrid-golang/jsonhooks-v1"
	"github.com/akamai/AkamaiOPEN-edgegrid-golang/papi-v1"
	"github.com/cyberphone/json-canonicalization/go/src/webpki.org/jsoncanonicalizer"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/schema"
	"github.com/tidwall/gjson"
)

// FIXME this function is identical to suppressEquivalentJSONPendingDiffs
func suppressEquivalentJSONDiffs(_, old, new string, d *schema.ResourceData) bool {
	logger := akamai.Log("PAPI", "suppressEquivalentJSONDiffs")

	jsonCompactor := func(dest *bytes.Buffer, input []byte) error {
		// canonicalize input using jsoncanonicalizer
		res, err := jsoncanonicalizer.Transform(input)
		if err != nil {
			return err
		}

		// compact the input
		return json.Compact(dest, res)
	}

	oldBuf := bytes.NewBuffer([]byte{})
	if err := jsonCompactor(oldBuf, []byte(old)); err != nil {
		logger.Errorf("converting to compact json: %s", old)
		return false
	}
	newBuf := bytes.NewBuffer([]byte{})
	if err := jsonCompactor(newBuf, []byte(new)); err != nil {
		logger.Errorf("converting to compact json: %s", old)
		return false
	}
	logger.Debugf("old json: %s", string(oldBuf.Bytes()))
	logger.Debugf("new json: %s", string(newBuf.Bytes()))
	rulesOld, err := getRulesForComp(d, old, "", logger)
	if err != nil {
		// todo not sure what to do with this error
		logger.Warnf("calling 'getRulesForComp': %s", err.Error())
	}
	rulesOld.Etag = ""
	body, err := jsonhooks.Marshal(rulesOld)
	if err != nil {
		logger.Errorf("marshaling rules: %s", err.Error())
		return false
	}
	sha1hashOld := tools.GetSHAString(string(body))
	logger.Debugf("SHA from OLD Json %s", sha1hashOld)
	rulesNew, err := getRulesForComp(d, new, "", logger)
	if err != nil {
		// todo not sure what to do with this error
		logger.Warnf("calling 'getRulesForComp': %s", err.Error())
	}
	rulesNew.Etag = ""
	jsonBodyNew, err := jsonhooks.Marshal(rulesNew)
	if err != nil {
		logger.Errorf("marshaling rules: %s", err.Error())
		return false
	}
	sha1hashNew := tools.GetSHAString(string(jsonBodyNew))
	logger.Debugf("SHA from NEW Json %s", sha1hashNew)
	if sha1hashOld == sha1hashNew {
		logger.Debugf("SHA Equal skip diff")
		return true
	}
	logger.Debugf("SHA Not Equal diff applies")
	return false
}

func suppressEquivalentJSONPendingDiffs(old, new string, d *schema.ResourceDiff) bool {
	logger := akamai.Log("PAPI", "suppressEquivalentJSONPendingDiffs")

	jsonCompactor := func(dest *bytes.Buffer, input []byte) error {
		// canonicalize input using jsoncanonicalizer
		res, err := jsoncanonicalizer.Transform(input)
		if err != nil {
			return err
		}

		// compact the input
		return json.Compact(dest, res)
	}

	oldBuf := bytes.NewBuffer([]byte{})
	if err := jsonCompactor(oldBuf, []byte(old)); err != nil {
		logger.Errorf("converting old value to compact json: %s", old)
		return false
	}

	newBuf := bytes.NewBuffer([]byte{})
	if err := jsonCompactor(newBuf, []byte(new)); err != nil {
		logger.Errorf("converting new value to compact json: %s", newBuf)
		return false
	}
	logger.Debugf("old json: %s", string(oldBuf.Bytes()))
	logger.Debugf("new json: %s", string(newBuf.Bytes()))

	rulesOld, err := getRulesForComp(d, old, "", logger)
	if err != nil {
		// todo not sure what to do with this error
		logger.Warnf("calling 'getRulesForComp': %s", err.Error())
	}
	rulesOld.Etag = ""
	body, err := jsonhooks.Marshal(rulesOld)
	if err != nil {
		logger.Errorf("marshaling rules: %s", err.Error())
		return false
	}
	sha1hashOld := tools.GetSHAString(string(body))
	logger.Debugf("SHA from OLD Json %s", sha1hashOld)
	rulesNew, err := getRulesForComp(d, new, "", logger)
	if err != nil {
		// todo not sure what to do with this error
		logger.Warnf("calling 'getRulesForComp': %s", err.Error())
	}
	rulesNew.Etag = ""
	jsonBodyNew, err := jsonhooks.Marshal(rulesNew)
	if err != nil {
		logger.Errorf("marshaling rules: %s", err.Error())
		return false
	}
	sha1hashNew := tools.GetSHAString(string(jsonBodyNew))
	logger.Debugf("SHA from NEW Json %s", sha1hashNew)
	if sha1hashOld == sha1hashNew {
		logger.Debugf("SHA Equal skip diff")
		return true
	}
	logger.Debugf("SHA Not Equal diff applies")
	return false
}

// TODO: discuss how property rules should be handled
func getRulesForComp(d interface{}, json string, correlationid string, logger log.Interface) (*papi.Rules, error) {
	property, err := getProperty(d, correlationid, logger)
	if err != nil {
		return nil, err
	}

	rules := papi.NewRules()
	rules.Rule.Name = "default"
	switch d.(type) {
	case *schema.ResourceData:
		rules.PropertyID = d.(*schema.ResourceData).Id()
	case *schema.ResourceDiff:
		rules.PropertyID = d.(*schema.ResourceDiff).Id()
	default:
		return nil, fmt.Errorf("resource is of invalid type; should be '*schema.ResourceDiff' or '*schema.ResourceData'")
	}
	rules.PropertyVersion = property.LatestVersion
	origin, err := createOrigin(d, correlationid, logger)
	if err != nil {
		return nil, err
	}
	// get rules from the TF config
	logger.Debugf("Unmarshal Rules from JSON")
	unmarshalRulesFromJSONComp(d, json, rules)

	var ruleFormat string
	switch d.(type) {
	case *schema.ResourceData:
		ruleFormat, err = tools.GetStringValue("rule_format", d.(*schema.ResourceData))
	case *schema.ResourceDiff:
		ruleFormat, err = tools.GetStringValue("rule_format", d.(*schema.ResourceDiff))
	default:
		return nil, fmt.Errorf("resource is of invalid type; should be '*schema.ResourceDiff' or '*schema.ResourceData'")
	}
	if err != nil && !errors.Is(err, tools.ErrNotFound) {
		return nil, err
	}

	if err == nil {
		rules.RuleFormat = ruleFormat
	} else {
		ruleFormats := papi.NewRuleFormats()
		rules.RuleFormat, err = ruleFormats.GetLatest(correlationid)
		if err != nil {
			return nil, err
		}
	}

	cpCode, err := getCPCode(d, property.Contract, property.Group, correlationid, logger)
	if err != nil {
		return nil, err
	}
	updateStandardBehaviors(rules, cpCode, origin, correlationid, logger)
	fixupPerformanceBehaviors(rules, correlationid, logger)

	return rules, nil
}

// TODO: discuss how property rules should be handled
func unmarshalRulesFromJSONComp(d interface{}, rulesComp string, propertyRules *papi.Rules) {
	propertyRules.Rule = &papi.Rule{Name: "default"}
	rulesJSON := gjson.Parse(rulesComp).Get("rules")
	rulesJSON.ForEach(func(key, value gjson.Result) bool {
		if key.String() == "behaviors" {
			behavior := gjson.Parse(value.String())
			if gjson.Get(behavior.String(), "#.name").Exists() {

				behavior.ForEach(func(key, value gjson.Result) bool {

					bb, ok := value.Value().(map[string]interface{})
					if ok {
						for k, v := range bb {
							log.Infof("k:", k, "v:", v)
						}

						beh := papi.NewBehavior()

						beh.Name = bb["name"].(string)
						boptions, ok := bb["options"]
						if ok {
							beh.Options = boptions.(map[string]interface{})
						}

						propertyRules.Rule.MergeBehavior(beh)
					}

					return true // keep iterating
				}) // behavior list loop

			}

			if key.String() == "criteria" {
				criteria := gjson.Parse(value.String())

				criteria.ForEach(func(key, value gjson.Result) bool {

					cc, ok := value.Value().(map[string]interface{})
					if ok {
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

			for _, rule := range extractRulesJSON(d, childRules) {
				propertyRules.Rule.MergeChildRule(rule)
			}
		}

		if key.String() == "variables" {

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
					propertyRules.Rule.AddVariable(newVariable)
				}
				return true
			}) //variables

		}

		if key.String() == "options" {
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
		variables := gjson.Parse(jsonvars.(string))
		result := gjson.Get(variables.String(), "variables")

		result.ForEach(func(key, value gjson.Result) bool {
			variableMap, ok := value.Value().(map[string]interface{})
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
	var isSecure interface{}
	var set bool

	switch d.(type) {
	case *schema.ResourceData:
		isSecure, set = d.(*schema.ResourceData).GetOk("is_secure")
	case *schema.ResourceDiff:
		isSecure, set = d.(*schema.ResourceDiff).GetOk("is_secure")
	default:
		isSecure, set = d.(*schema.ResourceData).GetOk("is_secure")
	}

	if set && isSecure.(bool) {
		propertyRules.Rule.Options.IsSecure = true
	} else if set && !isSecure.(bool) {
		propertyRules.Rule.Options.IsSecure = false
	}

	// ADD cpCode from resource
	var cpCode interface{}

	switch d.(type) {
	case *schema.ResourceData:
		cpCode, set = d.(*schema.ResourceData).GetOk("cp_code")
	case *schema.ResourceDiff:
		cpCode, set = d.(*schema.ResourceDiff).GetOk("cp_code")
	default:
		cpCode, set = d.(*schema.ResourceData).GetOk("cp_code")
	}

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
