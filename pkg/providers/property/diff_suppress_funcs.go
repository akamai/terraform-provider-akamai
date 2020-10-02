package property

import (
	"encoding/json"
	"reflect"
	"sort"

	"github.com/akamai/AkamaiOPEN-edgegrid-golang/v2/pkg/papi"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/schema"
)

func suppressEquivalentJSONDiffs(_, old, new string, _ *schema.ResourceData) bool {
	return compareRulesJSON(old, new)
}

func compareRulesJSON(old, new string) bool {
	var oldRules, newRules papi.GetRuleTreeResponse
	if old == new {
		return true
	}
	if err := json.Unmarshal([]byte(old), &oldRules); err != nil {
		return false
	}
	if err := json.Unmarshal([]byte(new), &newRules); err != nil {
		return false
	}
	diff := compareRules(&oldRules.Rules, &newRules.Rules)
	return diff
}

// compareRules handles comparison between two papi.Rules objects
// due to an issue in PAPI we need to compare collections of behaviors, criteria and variables discarding the order from JSON
func compareRules(old, new *papi.Rules) bool {
	if len(old.Behaviors) != len(new.Behaviors) ||
		len(old.Criteria) != len(new.Criteria) ||
		len(old.Variables) != len(new.Variables) ||
		len(old.Children) != len(new.Children) {
		return false
	}
	if len(old.Children) > 0 {
		for i := range old.Children {
			if !compareRules(&old.Children[i], &new.Children[i]) {
				return false
			}
		}
	} else {
		old.Children = nil
		new.Children = nil
	}

	old.Behaviors = orderBehaviors(old.Behaviors)
	new.Behaviors = orderBehaviors(new.Behaviors)

	old.Criteria = orderBehaviors(old.Criteria)
	new.Criteria = orderBehaviors(new.Criteria)

	old.Variables = orderVariables(old.Variables)
	new.Variables = orderVariables(new.Variables)

	return reflect.DeepEqual(old, new)
}

func orderBehaviors(behaviors []papi.RuleBehavior) []papi.RuleBehavior {
	if len(behaviors) == 0 {
		return nil
	}
	sort.Slice(behaviors, func(i, j int) bool {
		return behaviors[i].Name < behaviors[j].Name
	})
	return behaviors
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
