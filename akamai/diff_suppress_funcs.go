package akamai

import (
	"bytes"
	"encoding/json"
	"log"
	"strings"

	"github.com/akamai/AkamaiOPEN-edgegrid-golang/jsonhooks-v1"
	"github.com/akamai/AkamaiOPEN-edgegrid-golang/papi-v1"
	"github.com/hashicorp/terraform/helper/schema"
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

	log.Printf("[DEBUG] SHA from OLD Json %s\n", sha1hashOld)

	rulesNew, err := getRulesForComp(d, new)
	rulesNew.Etag = ""
	jsonBodyNew, err := jsonhooks.Marshal(rulesNew)
	if err != nil {
		return false
	}
	sha1hashNew := getSHAString(string(jsonBodyNew))

	log.Printf("[DEBUG] SHA from NEW Json %s\n", sha1hashNew)

	if sha1hashOld == sha1hashNew {
		return true
	} else {
		return false
	}
	//return jsonBytesEqual(ob.Bytes(), nb.Bytes())
}

func getRulesForComp(d *schema.ResourceData, json string) (*papi.Rules, error) {

	property, e := getProperty(d)
	if e != nil {
		return nil, e
	}

	rules := papi.NewRules()
	rules.Rule.Name = "default"
	id := strings.Split(d.Id(), "-")
	rules.PropertyID = id[0]
	rules.PropertyVersion = property.LatestVersion

	origin, err := createOrigin(d)
	if err != nil {
		return nil, err
	}

	// get rules from the TF config

	//rulecheck

	log.Printf("[DEBUG] Unmarshal Rules from JSON")
	unmarshalRulesFromJSON(d, json, rules)

	if ruleFormat, ok := d.GetOk("rule_format"); ok {
		rules.RuleFormat = ruleFormat.(string)
	} else {
		ruleFormats := papi.NewRuleFormats()
		rules.RuleFormat, err = ruleFormats.GetLatest()
		if err != nil {
			return nil, err
		}
	}

	if ok := d.HasChange("rule_format"); ok {
	}

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

func suppressEquivalentJsonRules(k, old, new string, d *schema.ResourceData) bool {

	// Deserialize and serialize through edgegrid-golang to ensure that the serialized strings are equivalent
	// This handles the case where edgegrid-golang has a different "omitEmpty" scheme to other api implementations
	//
	// When marshaling, we only consider the "Rules" part and not the header
	//
	// Note: if this function determines that the two rule sets are different, Terraform will show ALL
	// differences in the plan, even those that considered trivial
	//

	nrules := papi.NewRules()
	orules := papi.NewRules()

	if err := json.Unmarshal([]byte(old), orules); err != nil {
		return false
	}
	nold, err := json.Marshal(orules.Rule)
	if err != nil {
		return false
	}

	if err := json.Unmarshal([]byte(new), nrules); err != nil {
		return false
	}
	nnew, err := json.Marshal(nrules.Rule)
	if err != nil {
		return false
	}

	return suppressEquivalentJsonDiffs(k, string(nold), string(nnew), d)
}
