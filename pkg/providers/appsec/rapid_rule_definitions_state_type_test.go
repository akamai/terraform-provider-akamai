package appsec

import (
	"testing"

	"github.com/akamai/AkamaiOPEN-edgegrid-golang/v10/pkg/appsec"
	"github.com/akamai/terraform-provider-akamai/v6/pkg/common/ptr"
	"github.com/stretchr/testify/assert"
)

func TestCheckSemanticEquality_RuleDefinitions(t *testing.T) {
	var before = baseASC()
	var after = baseDESC()

	assert.Equal(t, []string(nil), checkSemanticEquality(before, after))
}

func exception() *appsec.RuleConditionException {
	return &appsec.RuleConditionException{
		Exception: &appsec.RuleException{
			SpecificHeaderCookieParamXMLOrJSONNames: &appsec.SpecificHeaderCookieParamXMLOrJSONNames{
				{
					Names:    []string{"Auth"},
					Selector: "REQUEST_HEADERS",
					Wildcard: false,
				},
			},
		},
	}
}

func baseDESC() []appsec.RuleDefinition {
	return []appsec.RuleDefinition{
		{
			ID:                 ptr.To(int64(11)),
			Action:             ptr.To("alert"),
			Lock:               ptr.To(false),
			ConditionException: exception(),
		},
		{
			ID:                 ptr.To(int64(1)),
			Action:             ptr.To("deny"),
			Lock:               ptr.To(true),
			ConditionException: exception(),
		},
	}
}

func baseASC() []appsec.RuleDefinition {
	return []appsec.RuleDefinition{
		{
			ID:                 ptr.To(int64(1)),
			Action:             ptr.To("deny"),
			Lock:               ptr.To(true),
			ConditionException: exception(),
		},
		{
			ID:                 ptr.To(int64(11)),
			Action:             ptr.To("alert"),
			Lock:               ptr.To(false),
			ConditionException: exception(),
		},
	}
}
