package papi

import (
	"fmt"
	"strings"

	"github.com/akamai/AkamaiOPEN-edgegrid-golang/client-v1"
)

// Rules is a collection of property rules
type Rules struct {
	client.Resource
	AccountID       string        `json:"accountId"`
	ContractID      string        `json:"contractId"`
	GroupID         string        `json:"groupId"`
	PropertyID      string        `json:"propertyId"`
	PropertyVersion int           `json:"propertyVersion"`
	Etag            string        `json:"etag"`
	RuleFormat      string        `json:"ruleFormat"`
	Rule            *Rule         `json:"rules"`
	Errors          []*RuleErrors `json:"errors,omitempty"`
}

// NewRules creates a new Rules
func NewRules() *Rules {
	rules := &Rules{}
	rules.Rule = NewRule()
	rules.Rule.Name = "default"
	rules.Init()

	return rules
}

// PreMarshalJSON is called before JSON marshaling
//
// See: jsonhooks-v1/json.Marshal()
func (rules *Rules) PreMarshalJSON() error {
	rules.Errors = nil
	return nil
}

// GetRules populates Rules with rule data for a given property
//
// See: Property.GetRules
// API Docs: https://developer.akamai.com/api/luna/papi/resources.html#getaruletree
// Endpoint: GET /papi/v1/properties/{propertyId}/versions/{propertyVersion}/rules/{?contractId,groupId}
func (rules *Rules) GetRules(property *Property) error {
	req, err := client.NewRequest(
		Config,
		"GET",
		fmt.Sprintf(
			"/papi/v1/properties/%s/versions/%d/rules",
			property.PropertyID,
			property.LatestVersion,
		),
		nil,
	)
	if err != nil {
		return err
	}

	res, err := client.Do(Config, req)
	if err != nil {
		return err
	}

	if client.IsError(res) {
		return client.NewAPIError(res)
	}

	if err = client.BodyJSON(res, rules); err != nil {
		return err
	}

	return nil
}

// GetRulesDigest fetches the Etag for a rule tree
//
// See: Property.GetRulesDigest()
// API Docs: https://developer.akamai.com/api/luna/papi/resources.html#getaruletreesdigest
// Endpoint: HEAD /papi/v1/properties/{propertyId}/versions/{propertyVersion}/rules/{?contractId,groupId}
func (rules *Rules) GetRulesDigest(property *Property) (string, error) {
	req, err := client.NewRequest(
		Config,
		"HEAD",
		fmt.Sprintf(
			"/papi/v1/properties/%s/versions/%d/rules",
			property.PropertyID,
			property.LatestVersion,
		),
		nil,
	)
	if err != nil {
		return "", err
	}

	res, err := client.Do(Config, req)
	if err != nil {
		return "", err
	}

	if client.IsError(res) {
		return "", client.NewAPIError(res)
	}

	return res.Header.Get("Etag"), nil
}

// Save creates/updates a rule tree for a property
//
// API Docs: https://developer.akamai.com/api/luna/papi/resources.html#putpropertyversionrules
// Endpoint: PUT /papi/v1/properties/{propertyId}/versions/{propertyVersion}/rules{?contractId,groupId}
func (rules *Rules) Save() error {
	rules.Errors = []*RuleErrors{}

	req, err := client.NewJSONRequest(
		Config,
		"PUT",
		fmt.Sprintf(
			"/papi/v1/properties/%s/versions/%d/rules",
			rules.PropertyID,
			rules.PropertyVersion,
		),
		rules,
	)
	if err != nil {
		return err
	}

	res, err := client.Do(Config, req)
	if err != nil {
		return err
	}

	if client.IsError(res) {
		return client.NewAPIError(res)
	}

	if err = client.BodyJSON(res, rules); err != nil {
		return err
	}

	if len(rules.Errors) != 0 {
		return ErrorMap[ErrInvalidRules]
	}

	return nil
}

// Freeze pins a properties rule set to a specific rule set version
func (rules *Rules) Freeze(format string) error {
	rules.Errors = []*RuleErrors{}

	req, err := client.NewJSONRequest(
		Config,
		"PUT",
		fmt.Sprintf(
			"/papi/v1/properties/%s/versions/%d/rules",
			rules.PropertyID,
			rules.PropertyVersion,
		),
		rules,
	)
	if err != nil {
		return err
	}

	req.Header.Set("Content-Type", fmt.Sprintf("application/vnd.akamai.papirules.%s+json", format))

	res, err := client.Do(Config, req)
	if err != nil {
		return err
	}

	if client.IsError(res) {
		return client.NewAPIError(res)
	}

	if err = client.BodyJSON(res, rules); err != nil {
		return err
	}

	if len(rules.Errors) != 0 {
		return ErrorMap[ErrInvalidRules]
	}

	return nil
}

// Rule represents a property rule resource
type Rule struct {
	client.Resource
	Depth               int                          `json:"-"`
	Name                string                       `json:"name"`
	Criteria            []*Criteria                  `json:"criteria,omitempty"`
	Behaviors           []*Behavior                  `json:"behaviors,omitempty"`
	Children            []*Rule                      `json:"children,omitempty"`
	Comments            string                       `json:"comments,omitempty"`
	CriteriaLocked      bool                         `json:"criteriaLocked,omitempty"`
	CriteriaMustSatisfy RuleCriteriaMustSatisfyValue `json:"criteriaMustSatisfy,omitempty"`
	UUID                string                       `json:"uuid,omitempty"`
	Variables           []*Variable                  `json:"variables,omitempty"`
	AdvancedOverride    string                       `json:"advancedOverride,omitempty"`

	Options struct {
		IsSecure bool `json:"is_secure,omitempty"`
	} `json:"options,omitempty"`

	CustomOverride struct {
		Name       string `json:"name"`
		OverrideID string `json:"overrideId"`
	} `json:"customOverride,omitempty"`
}

// NewRule creates a new Rule
func NewRule() *Rule {
	rule := &Rule{}
	rule.Init()

	return rule
}

// MergeBehavior merges a behavior into a rule
//
// If the behavior already exists, it's options are merged with the existing
// options.
func (rule *Rule) MergeBehavior(behavior *Behavior) {
	for _, existingBehavior := range rule.Behaviors {
		if existingBehavior.Name == behavior.Name {
			existingBehavior.MergeOptions(behavior.Options)
			return
		}
	}

	rule.Behaviors = append(rule.Behaviors, behavior)
}

// AddBehavior adds a behavior to the rule
//
// If the behavior already exists it is replaced with the given behavior
func (rule *Rule) AddBehavior(behavior *Behavior) {
	for key, existingBehavior := range rule.Behaviors {
		if existingBehavior.Name == behavior.Name {
			rule.Behaviors[key] = behavior
			return
		}
	}

	rule.Behaviors = append(rule.Behaviors, behavior)
}

// MergeCriteria merges a criteria into a rule
//
// If the criteria already exists, it's options are merged with the existing
// options.
func (rule *Rule) MergeCriteria(criteria *Criteria) {
	for _, existingCriteria := range rule.Criteria {
		if existingCriteria.Name == criteria.Name {
			existingCriteria.MergeOptions(criteria.Options)
			return
		}
	}

	rule.Criteria = append(rule.Criteria, criteria)
}

// AddCriteria add a criteria to a rule
//
// If the criteria already exists, it is replaced with the given criteria.
func (rule *Rule) AddCriteria(criteria *Criteria) {
	for key, existingCriteria := range rule.Criteria {
		if existingCriteria.Name == criteria.Name {
			rule.Criteria[key] = criteria
			return
		}
	}

	rule.Criteria = append(rule.Criteria, criteria)
}

// MergeChildRule adds a child rule to this rule
//
// If the rule already exists, criteria, behaviors, and child rules are added to
// the existing rule.
func (rule *Rule) MergeChildRule(childRule *Rule) {
	for key, existingChildRule := range rule.Children {
		if existingChildRule.Name == childRule.Name {
			for _, behavior := range childRule.Behaviors {
				rule.Children[key].MergeBehavior(behavior)
			}

			for _, criteria := range childRule.Criteria {
				rule.Children[key].MergeCriteria(criteria)
			}

			for _, child := range childRule.Children {
				rule.Children[key].MergeChildRule(child)
			}

			return
		}
	}

	rule.Children = append(rule.Children, childRule)
}

// AddChildRule adds a rule as a child of this rule
//
// If the rule already exists, it is replaced by the given rule.
func (rule *Rule) AddChildRule(childRule *Rule) {
	for key, existingChildRule := range rule.Children {
		if existingChildRule.Name == childRule.Name {
			rule.Children[key] = childRule

			return
		}
	}

	rule.Children = append(rule.Children, childRule)
}

// AddVariable adds a variable as a child of this rule
//
// If the rule already exists, it is replaced by the given rule.
func (rule *Rule) AddVariable(variable *Variable) {
	for key, existingVariable := range rule.Variables {
		if existingVariable.Name == variable.Name {
			rule.Variables[key] = variable

			return
		}
	}

	rule.Variables = append(rule.Variables, variable)
}

// FindBehavior locates a specific behavior by path
func (rules *Rules) FindBehavior(path string) (*Behavior, error) {
	if len(path) <= 1 {
		return nil, ErrorMap[ErrInvalidPath]
	}

	rule, err := rules.FindParentRule(path)
	if err != nil {
		return nil, err
	}

	sep := "/"
	segments := strings.Split(path, sep)
	behaviorName := strings.ToLower(segments[len(segments)-1])
	for _, behavior := range rule.Behaviors {
		if strings.ToLower(behavior.Name) == behaviorName {
			return behavior, nil
		}
	}

	return nil, ErrorMap[ErrBehaviorNotFound]
}

// FindCriteria locates a specific Critieria by path
func (rules *Rules) FindCriteria(path string) (*Criteria, error) {
	if len(path) <= 1 {
		return nil, ErrorMap[ErrInvalidPath]
	}

	rule, err := rules.FindParentRule(path)
	if err != nil {
		return nil, err
	}

	sep := "/"
	segments := strings.Split(path, sep)
	criteriaName := strings.ToLower(segments[len(segments)-1])
	for _, criteria := range rule.Criteria {
		if strings.ToLower(criteria.Name) == criteriaName {
			return criteria, nil
		}
	}

	return nil, ErrorMap[ErrCriteriaNotFound]
}

// FindVariable locates a specific Variable by path
func (rules *Rules) FindVariable(path string) (*Variable, error) {
	if len(path) <= 1 {
		return nil, ErrorMap[ErrInvalidPath]
	}

	rule, err := rules.FindParentRule(path)
	if err != nil {
		return nil, err
	}

	sep := "/"
	segments := strings.Split(path, sep)
	variableName := strings.ToLower(segments[len(segments)-1])
	for _, variable := range rule.Variables {
		if strings.ToLower(variable.Name) == variableName {
			return variable, nil
		}
	}

	return nil, ErrorMap[ErrVariableNotFound]
}

// FindRule locates a specific rule by path
func (rules *Rules) FindRule(path string) (*Rule, error) {
	if path == "" {
		return rules.Rule, nil
	}

	sep := "/"
	segments := strings.Split(path, sep)

	currentRule := rules.Rule
	for _, segment := range segments {
		found := false
		for _, rule := range currentRule.Children {
			if strings.ToLower(rule.Name) == segment {
				currentRule = rule
				found = true
			}
		}
		if found != true {
			return nil, ErrorMap[ErrRuleNotFound]
		}
	}

	return currentRule, nil
}

// Find the parent rule for a given rule, criteria, or behavior path
func (rules *Rules) FindParentRule(path string) (*Rule, error) {
	sep := "/"
	segments := strings.Split(strings.ToLower(strings.TrimPrefix(path, sep)), sep)
	parentPath := strings.Join(segments[0:len(segments)-1], sep)

	return rules.FindRule(parentPath)
}

// Criteria represents a rule criteria resource
type Criteria struct {
	client.Resource
	Name    string      `json:"name"`
	Options OptionValue `json:"options"`
	UUID    string      `json:"uuid,omitempty"`
	Locked  bool        `json:"locked,omitempty"`
}

// NewCriteria creates a new Criteria
func NewCriteria() *Criteria {
	criteria := &Criteria{Options: OptionValue{}}
	criteria.Init()

	return criteria
}

// MergeOptions merges the given options with the existing options
func (criteria *Criteria) MergeOptions(newOptions OptionValue) {
	options := make(map[string]interface{})
	for k, v := range criteria.Options {
		options[k] = v
	}

	for k, v := range newOptions {
		options[k] = v
	}

	criteria.Options = OptionValue(options)
}

// Behavior represents a rule behavior resource
type Behavior struct {
	client.Resource
	Name    string      `json:"name"`
	Options OptionValue `json:"options"`
	Locked  bool        `json:"locked,omitempty"`
	UUID    string      `json:"uuid,omitempty"`
}

// NewBehavior creates a new Behavior
func NewBehavior() *Behavior {
	behavior := &Behavior{Options: OptionValue{}}
	behavior.Init()

	return behavior
}

// MergeOptions merges the given options with the existing options
func (behavior *Behavior) MergeOptions(newOptions OptionValue) {
	options := make(map[string]interface{})
	for k, v := range behavior.Options {
		options[k] = v
	}

	for k, v := range newOptions {
		options[k] = v
	}

	behavior.Options = OptionValue(options)
}

// OptionValue represents a generic option value
//
// OptionValue is a map with string keys, and any
// type of value. You can nest OptionValues as necessary
// to create more complex values.
type OptionValue map[string]interface{}

type Variable struct {
	client.Resource
	Name        string `json:"name"`
	Value       string `json:"value"`
	Description string `json:"description"`
	Hidden      bool   `json:"hidden"`
	Sensitive   bool   `json:"sensitive"`
}

// NewVariable creates a new Variable
func NewVariable() *Variable {
	variable := &Variable{}
	variable.Init()

	return variable
}

// RuleErrors represents an validate error returned for a rule
type RuleErrors struct {
	client.Resource
	Type         string `json:"type"`
	Title        string `json:"title"`
	Detail       string `json:"detail"`
	Instance     string `json:"instance"`
	BehaviorName string `json:"behaviorName"`
}

// NewRuleErrors creates a new RuleErrors
func NewRuleErrors() *RuleErrors {
	ruleErrors := &RuleErrors{}
	ruleErrors.Init()

	return ruleErrors
}

type RuleCriteriaMustSatisfyValue string

const (
	RuleCriteriaMustSatisfyAll RuleCriteriaMustSatisfyValue = "any"
	RuleCriteriaMustSatisfyAny RuleCriteriaMustSatisfyValue = "all"
)
