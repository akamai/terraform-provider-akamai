package property

import (
	"encoding/json"
	"fmt"
	"reflect"
	"sort"

	"github.com/akamai/AkamaiOPEN-edgegrid-golang/v7/pkg/papi"
	"github.com/akamai/terraform-provider-akamai/v5/pkg/logger"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/schema"
)

func diffSuppressRules(_, oldRules, newRules string, _ *schema.ResourceData) bool {
	rulesEqual, err := rulesJSONEqual(oldRules, newRules)
	if err != nil {
		logger.Get("PAPI", "diffSuppressRules").Error(err.Error())
	}

	return rulesEqual
}

// rulesJSONEqual handles comparison between two papi.RulesUpdate JSON representations.
func rulesJSONEqual(old, new string) (bool, error) {
	if old == "" || new == "" {
		return old == new, nil
	}

	if old == new {
		return true, nil
	}

	var oldRules papi.RulesUpdate
	if err := json.Unmarshal([]byte(old), &oldRules); err != nil {
		return false, fmt.Errorf("'old' = %s, unmarshal: %w", old, err)
	}

	var newRules papi.RulesUpdate
	if err := json.Unmarshal([]byte(new), &newRules); err != nil {
		return false, fmt.Errorf("'new' = %s, unmarshal: %w", new, err)
	}

	return ruleTreesEqual(&oldRules, &newRules), nil
}

func ruleTreesEqual(old, new *papi.RulesUpdate) bool {
	if old.Comments != new.Comments {
		return false
	}

	return rulesEqual(&old.Rules, &new.Rules)
}

// rulesEqual handles comparison between two papi.Rules objects ignoring the order in
// collection of variables.
func rulesEqual(old, new *papi.Rules) bool {
	if len(old.Behaviors) != len(new.Behaviors) ||
		len(old.Criteria) != len(new.Criteria) ||
		len(old.Variables) != len(new.Variables) ||
		len(old.Children) != len(new.Children) {
		return false
	}

	if len(old.Children) > 0 {
		for i := range old.Children {
			if !rulesEqual(&old.Children[i], &new.Children[i]) {
				return false
			}
		}
	} else {
		old.Children = nil
		new.Children = nil
	}

	if len(old.Behaviors) == 0 {
		old.Behaviors = nil
	}
	if len(new.Behaviors) == 0 {
		new.Behaviors = nil
	}

	if len(old.Criteria) == 0 {
		old.Criteria = nil
	}
	if len(new.Criteria) == 0 {
		new.Criteria = nil
	}

	old.Variables = orderVariables(old.Variables)
	new.Variables = orderVariables(new.Variables)

	removeNilOptions2(old)
	removeNilOptions2(new)

	return reflect.DeepEqual(old, new)
}

func removeNilOptions(rules *papi.Rules) {
	//for _, b := range rules.Behaviors {
	//	for k, v := range b.Options {
	//		if v == nil {
	//			delete(b.Options, k)
	//		} else if vv, ok := v.(map[string]interface{}); ok {
	//			for k, v := range vv {
	//				if v == nil {
	//					delete(vv, k)
	//				}
	//			}
	//		}
	//	}
	//}
}
func removeNilOptions2(rules *papi.Rules) {
	for _, b := range rules.Behaviors {
		for k, v := range b.Options {
			if v == nil {
				delete(b.Options, k)
			} else if vv, ok := v.(map[string]interface{}); ok {
				for k, v := range vv {
					if v == nil {
						delete(vv, k)
					}
				}
			}
		}
	}
}

func orderVariables(variables []papi.RuleVariable) []papi.RuleVariable {
	if len(variables) == 0 {
		return nil
	}
	sort.Slice(variables, func(i, j int) bool {
		return variables[i].Name < variables[j].Name
	})
	return variables
}
