package property

import (
	"testing"

	"github.com/akamai/AkamaiOPEN-edgegrid-golang/v11/pkg/papi"
	"github.com/akamai/terraform-provider-akamai/v8/pkg/common/ptr"
	"github.com/stretchr/testify/assert"
)

func TestRulesEqual(t *testing.T) {
	tests := map[string]struct {
		old      *papi.Rules
		new      *papi.Rules
		expected bool
	}{
		"equal rules": {
			old: &papi.Rules{
				AdvancedOverride: "AAA",
				Behaviors: []papi.RuleBehavior{
					{
						Locked: false,
						Name:   "BEH1",
						Options: map[string]interface{}{
							"opt1": 123,
							"opt2": 321,
						},
						UUID: "4535543",
					},
				},
				Children: []papi.Rules{},
				Comments: "comment",
				Criteria: []papi.RuleBehavior{
					{
						Locked:  false,
						Name:    "CRIT1",
						Options: nil,
						UUID:    "1234",
					},
					{
						Locked:  false,
						Name:    "CRIT2",
						Options: nil,
						UUID:    "1234",
					},
				},
				CriteriaLocked: true,
				CustomOverride: &papi.RuleCustomOverride{
					Name:       "custom",
					OverrideID: "1234",
				},
				Name: "A",
				Options: papi.RuleOptions{
					IsSecure: true,
				},
				UUID: "43242",
				Variables: []papi.RuleVariable{
					{
						Description: ptr.To("var1"),
						Hidden:      true,
						Name:        "VAR1",
						Sensitive:   true,
						Value:       ptr.To("value 1"),
					},
					{
						Description: ptr.To("var1"),
						Hidden:      true,
						Name:        "VAR2",
						Sensitive:   true,
						Value:       ptr.To("value 1"),
					}},
			},
			new: &papi.Rules{
				AdvancedOverride: "AAA",
				Behaviors: []papi.RuleBehavior{
					{
						Locked: false,
						Name:   "BEH1",
						Options: map[string]interface{}{
							"opt1": 123,
							"opt2": 321,
						},
						UUID: "4535543",
					},
				},
				Children: []papi.Rules{},
				Comments: "comment",
				Criteria: []papi.RuleBehavior{
					{
						Locked:  false,
						Name:    "CRIT1",
						Options: nil,
						UUID:    "1234",
					},
					{
						Locked:  false,
						Name:    "CRIT2",
						Options: nil,
						UUID:    "1234",
					},
				},
				CriteriaLocked: true,
				CustomOverride: &papi.RuleCustomOverride{
					Name:       "custom",
					OverrideID: "1234",
				},
				Name: "A",
				Options: papi.RuleOptions{
					IsSecure: true,
				},
				UUID: "43242",
				Variables: []papi.RuleVariable{
					{
						Description: ptr.To("var1"),
						Hidden:      true,
						Name:        "VAR1",
						Sensitive:   true,
						Value:       ptr.To("value 1"),
					},
					{
						Description: ptr.To("var1"),
						Hidden:      true,
						Name:        "VAR2",
						Sensitive:   true,
						Value:       ptr.To("value 1"),
					},
				},
			},
			expected: true,
		},

		"equal rules, different order": {
			old: &papi.Rules{
				AdvancedOverride: "AAA",
				Behaviors: []papi.RuleBehavior{
					{
						Locked: false,
						Name:   "BEH2",
						Options: map[string]interface{}{
							"opt1": 123,
							"opt2": 321,
						},
						UUID: "4535543",
					},
					{
						Locked: false,
						Name:   "BEH1",
						Options: map[string]interface{}{
							"opt1": 123,
							"opt2": 321,
						},
						UUID: "4535543",
					},
				},
				Children: []papi.Rules{},
				Comments: "comment",
				Criteria: []papi.RuleBehavior{
					{
						Locked:  false,
						Name:    "CRIT2",
						Options: nil,
						UUID:    "1234",
					},
					{
						Locked:  false,
						Name:    "CRIT1",
						Options: nil,
						UUID:    "1234",
					},
				},
				CriteriaLocked: true,
				CustomOverride: &papi.RuleCustomOverride{
					Name:       "custom",
					OverrideID: "1234",
				},
				Name: "A",
				Options: papi.RuleOptions{
					IsSecure: true,
				},
				UUID: "43242",
				Variables: []papi.RuleVariable{
					{
						Description: ptr.To("var1"),
						Hidden:      true,
						Name:        "VAR2",
						Sensitive:   true,
						Value:       ptr.To("value 1"),
					},
					{
						Description: ptr.To("var1"),
						Hidden:      true,
						Name:        "VAR1",
						Sensitive:   true,
						Value:       ptr.To("value 1"),
					},
				},
			},
			new: &papi.Rules{
				AdvancedOverride: "AAA",
				Behaviors: []papi.RuleBehavior{
					{
						Locked: false,
						Name:   "BEH1",
						Options: map[string]interface{}{
							"opt1": 123,
							"opt2": 321,
						},
						UUID: "4535543",
					},
					{
						Locked: false,
						Name:   "BEH2",
						Options: map[string]interface{}{
							"opt1": 123,
							"opt2": 321,
						},
						UUID: "4535543",
					},
				},
				Children: []papi.Rules{},
				Comments: "comment",
				Criteria: []papi.RuleBehavior{
					{
						Locked:  false,
						Name:    "CRIT1",
						Options: nil,
						UUID:    "1234",
					},
					{
						Locked:  false,
						Name:    "CRIT2",
						Options: nil,
						UUID:    "1234",
					},
				},
				CriteriaLocked: true,
				CustomOverride: &papi.RuleCustomOverride{
					Name:       "custom",
					OverrideID: "1234",
				},
				Name: "A",
				Options: papi.RuleOptions{
					IsSecure: true,
				},
				UUID: "43242",
				Variables: []papi.RuleVariable{
					{
						Description: ptr.To("var1"),
						Hidden:      true,
						Name:        "VAR1",
						Sensitive:   true,
						Value:       ptr.To("value 1"),
					},
					{
						Description: ptr.To("var1"),
						Hidden:      true,
						Name:        "VAR2",
						Sensitive:   true,
						Value:       ptr.To("value 1"),
					},
				},
			},
			expected: false,
		},

		"equal rules, with children": {
			old: &papi.Rules{
				AdvancedOverride: "AAA",
				Behaviors:        []papi.RuleBehavior{},
				Children: []papi.Rules{
					{
						Behaviors: []papi.RuleBehavior{
							{
								Locked: false,
								Name:   "BEH1",
								Options: map[string]interface{}{
									"opt1": 123,
									"opt2": 321,
								},
								UUID: "4535543",
							},
							{
								Locked: false,
								Name:   "BEH2",
								Options: map[string]interface{}{
									"opt1": 123,
									"opt2": 321,
								},
								UUID: "4535543",
							},
						},
						Criteria: nil,
						Name:     "RULE 1",
					},
					{
						Behaviors: []papi.RuleBehavior{
							{
								Locked: false,
								Name:   "BEH1",
								Options: map[string]interface{}{
									"opt1": 123,
									"opt2": 321,
								},
								UUID: "4535543",
							},
							{
								Locked: false,
								Name:   "BEH2",
								Options: map[string]interface{}{
									"opt1": 123,
									"opt2": 321,
								},
								UUID: "4535543",
							},
						},
						Criteria: nil,
						Name:     "RULE 2",
					},
				},
				Comments:       "comment",
				Criteria:       []papi.RuleBehavior{},
				CriteriaLocked: true,
				CustomOverride: &papi.RuleCustomOverride{
					Name:       "custom",
					OverrideID: "1234",
				},
				Name: "A",
				Options: papi.RuleOptions{
					IsSecure: true,
				},
				UUID:      "43242",
				Variables: []papi.RuleVariable{},
			},
			new: &papi.Rules{
				AdvancedOverride: "AAA",
				Behaviors:        []papi.RuleBehavior{},
				Children: []papi.Rules{
					{
						Behaviors: []papi.RuleBehavior{
							{
								Locked: false,
								Name:   "BEH1",
								Options: map[string]interface{}{
									"opt1": 123,
									"opt2": 321,
								},
								UUID: "4535543",
							},
							{
								Locked: false,
								Name:   "BEH2",
								Options: map[string]interface{}{
									"opt1": 123,
									"opt2": 321,
								},
								UUID: "4535543",
							},
						},
						Criteria: nil,
						Name:     "RULE 1",
					},
					{
						Behaviors: []papi.RuleBehavior{
							{
								Locked: false,
								Name:   "BEH1",
								Options: map[string]interface{}{
									"opt1": 123,
									"opt2": 321,
								},
								UUID: "4535543",
							},
							{
								Locked: false,
								Name:   "BEH2",
								Options: map[string]interface{}{
									"opt1": 123,
									"opt2": 321,
								},
								UUID: "4535543",
							},
						},
						Criteria: nil,
						Name:     "RULE 2",
					},
				},
				Comments:       "comment",
				Criteria:       []papi.RuleBehavior{},
				CriteriaLocked: true,
				CustomOverride: &papi.RuleCustomOverride{
					Name:       "custom",
					OverrideID: "1234",
				},
				Name: "A",
				Options: papi.RuleOptions{
					IsSecure: true,
				},
				UUID:      "43242",
				Variables: []papi.RuleVariable{},
			},
			expected: true,
		},

		"different rules, children in different order": {
			old: &papi.Rules{
				AdvancedOverride: "AAA",
				Behaviors:        []papi.RuleBehavior{},
				Children: []papi.Rules{
					{
						Behaviors: []papi.RuleBehavior{
							{
								Locked: false,
								Name:   "BEH1",
								Options: map[string]interface{}{
									"opt1": 123,
									"opt2": 321,
								},
								UUID: "4535543",
							},
							{
								Locked: false,
								Name:   "BEH2",
								Options: map[string]interface{}{
									"opt1": 123,
									"opt2": 321,
								},
								UUID: "4535543",
							},
						},
						Criteria: nil,
						Name:     "RULE 2",
					},
					{
						Behaviors: []papi.RuleBehavior{
							{
								Locked: false,
								Name:   "BEH1",
								Options: map[string]interface{}{
									"opt1": 123,
									"opt2": 321,
								},
								UUID: "4535543",
							},
							{
								Locked: false,
								Name:   "BEH2",
								Options: map[string]interface{}{
									"opt1": 123,
									"opt2": 321,
								},
								UUID: "4535543",
							},
						},
						Criteria: nil,
						Name:     "RULE 1",
					},
				},
				Comments:       "comment",
				Criteria:       []papi.RuleBehavior{},
				CriteriaLocked: true,
				CustomOverride: &papi.RuleCustomOverride{
					Name:       "custom",
					OverrideID: "1234",
				},
				Name: "A",
				Options: papi.RuleOptions{
					IsSecure: true,
				},
				UUID:      "43242",
				Variables: []papi.RuleVariable{},
			},
			new: &papi.Rules{
				AdvancedOverride: "AAA",
				Behaviors:        []papi.RuleBehavior{},
				Children: []papi.Rules{
					{
						Behaviors: []papi.RuleBehavior{
							{
								Locked: false,
								Name:   "BEH1",
								Options: map[string]interface{}{
									"opt1": 123,
									"opt2": 321,
								},
								UUID: "4535543",
							},
							{
								Locked: false,
								Name:   "BEH2",
								Options: map[string]interface{}{
									"opt1": 123,
									"opt2": 321,
								},
								UUID: "4535543",
							},
						},
						Criteria: nil,
						Name:     "RULE 1",
					},
					{
						Behaviors: []papi.RuleBehavior{
							{
								Locked: false,
								Name:   "BEH1",
								Options: map[string]interface{}{
									"opt1": 123,
									"opt2": 321,
								},
								UUID: "4535543",
							},
							{
								Locked: false,
								Name:   "BEH2",
								Options: map[string]interface{}{
									"opt1": 123,
									"opt2": 321,
								},
								UUID: "4535543",
							},
						},
						Criteria: nil,
						Name:     "RULE 2",
					},
				},
				Comments:       "comment",
				Criteria:       []papi.RuleBehavior{},
				CriteriaLocked: true,
				CustomOverride: &papi.RuleCustomOverride{
					Name:       "custom",
					OverrideID: "1234",
				},
				Name: "A",
				Options: papi.RuleOptions{
					IsSecure: true,
				},
				UUID:      "43242",
				Variables: []papi.RuleVariable{},
			},
			expected: false,
		},

		"different rules, different behavior len": {
			old: &papi.Rules{
				AdvancedOverride: "AAA",
				Behaviors: []papi.RuleBehavior{
					{
						Locked: false,
						Name:   "BEH1",
						Options: map[string]interface{}{
							"opt1": 123,
							"opt2": 321,
						},
						UUID: "4535543",
					},
					{
						Locked: false,
						Name:   "BEH2",
						Options: map[string]interface{}{
							"opt1": 123,
							"opt2": 321,
						},
						UUID: "4535543",
					},
				},
				Name: "A",
			},
			new: &papi.Rules{
				AdvancedOverride: "AAA",
				Behaviors: []papi.RuleBehavior{
					{
						Locked: false,
						Name:   "BEH1",
						Options: map[string]interface{}{
							"opt1": 123,
							"opt2": 321,
						},
						UUID: "4535543",
					},
				},
				Name: "B",
			},
			expected: false,
		},
		"equal rules, different slice representation": {
			old: &papi.Rules{
				AdvancedOverride: "AAA",
				Children:         []papi.Rules{},
				Behaviors:        []papi.RuleBehavior{},
				Criteria:         []papi.RuleBehavior{},
				Name:             "A",
			},
			new: &papi.Rules{
				AdvancedOverride: "AAA",
				Behaviors:        nil,
				Criteria:         nil,
				Children:         nil,
				Name:             "A",
			},
			expected: true,
		},
	}

	for name, test := range tests {
		t.Run(name, func(t *testing.T) {
			res := rulesEqual(test.old, test.new)
			assert.Equal(t, test.expected, res)
		})
	}
}

func TestRemoveNil(t *testing.T) {
	input := map[string]any{
		"a": "aval",
		"b": nil,
		"c": map[string]any{
			"c1": "c1val",
			"c2": nil,
			"c3": map[string]any{
				"c31": "c31val",
				"c32": nil,
				"c33": map[string]any{
					"c331": "c331val",
					"c332": nil,
				},
			},
		},
	}

	removeNils(input)
	assert.Equal(t, map[string]any{
		"a": "aval",
		"c": map[string]any{
			"c1": "c1val",
			"c3": map[string]any{
				"c31": "c31val",
				"c33": map[string]any{
					"c331": "c331val",
				},
			},
		},
	}, input)

}
