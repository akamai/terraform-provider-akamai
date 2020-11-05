package papi

import (
	"context"
	"errors"
	"fmt"
	"net/http"

	validation "github.com/go-ozzo/ozzo-validation/v4"
)

type (
	// PropertyRules contains operations available on PropertyRule resource
	// See: https://developer.akamai.com/api/core_features/property_manager/v1.html#propertyversionrulesgroup
	PropertyRules interface {
		// GetRuleTree lists all available CP codes
		// See: https://developer.akamai.com/api/core_features/property_manager/v1.html#getpropertyversionrules
		GetRuleTree(context.Context, GetRuleTreeRequest) (*GetRuleTreeResponse, error)

		// UpdateRuleTree lists all available CP codes
		// See: https://developer.akamai.com/api/core_features/property_manager/v1.html#putpropertyversionrules
		UpdateRuleTree(context.Context, UpdateRulesRequest) (*UpdateRulesResponse, error)
	}

	// GetRuleTreeRequest contains path and query params necessary to perform GET /rules request
	GetRuleTreeRequest struct {
		PropertyID      string
		PropertyVersion int
		ContractID      string
		GroupID         string
		ValidateMode    string
		ValidateRules   bool
	}

	// GetRuleTreeResponse contains data returned by performing GET /rules request
	GetRuleTreeResponse struct {
		Response
		PropertyID      string `json:"propertyId"`
		PropertyVersion int    `json:"propertyVersion"`
		Etag            string `json:"etag"`
		RuleFormat      string `json:"ruleFormat"`
		Rules           Rules  `json:"rules"`
	}

	// Rules contains Rule object
	Rules struct {
		AdvancedOverride    string                  `json:"advancedOverride,omitempty"`
		Behaviors           []RuleBehavior          `json:"behaviors,omitempty"`
		Children            []Rules                 `json:"children,omitempty"`
		Comment             string                  `json:"comment,omitempty"`
		Criteria            []RuleBehavior          `json:"criteria,omitempty"`
		CriteriaLocked      bool                    `json:"criteriaLocked,omitempty"`
		CustomOverride      *RuleCustomOverride     `json:"customOverride,omitempty"`
		Name                string                  `json:"name"`
		Options             RuleOptions             `json:"options,omitempty"`
		UUID                string                  `json:"uuid,omitempty"`
		Variables           []RuleVariable          `json:"variables,omitempty"`
		CriteriaMustSatisfy RuleCriteriaMustSatisfy `json:"criteriaMustSatisfy,omitempty"`
	}

	// RuleBehavior contains data for both rule behaviors and rule criteria
	RuleBehavior struct {
		Locked  string         `json:"locked,omitempty"`
		Name    string         `json:"name"`
		Options RuleOptionsMap `json:"options"`
		UUID    string         `json:"uuid,omitempty"`
	}

	// RuleCustomOverride represents customOverride field from Rule resource
	RuleCustomOverride struct {
		Name       string `json:"name"`
		OverrideID string `json:"overrideId"`
	}

	// RuleOptions represents options field from Rule resource
	RuleOptions struct {
		IsSecure bool `json:"is_secure,omitempty"`
	}

	// RuleVariable represents and entry in variables field from Rule resource
	RuleVariable struct {
		Description string `json:"description,omitempty"`
		Hidden      bool   `json:"hidden"`
		Name        string `json:"name"`
		Sensitive   bool   `json:"sensitive"`
		Value       string `json:"value,omitempty"`
	}

	// UpdateRulesRequest contains path and query params, as well as request body necessary to perform PUT /rules request
	UpdateRulesRequest struct {
		PropertyID      string
		PropertyVersion int
		ContractID      string
		DryRun          bool
		GroupID         string
		ValidateMode    string
		ValidateRules   bool
		Rules           RulesUpdate
	}

	// RulesUpdate is a wrapper for the request body of PUT /rules request
	RulesUpdate struct {
		Rules Rules `json:"rules"`
	}

	// UpdateRulesResponse contains data returned by performing PUT /rules request
	UpdateRulesResponse struct {
		AccountID       string      `json:"accountId"`
		ContractID      string      `json:"contractId"`
		GroupID         string      `json:"groupId"`
		PropertyID      string      `json:"propertyId"`
		PropertyVersion int         `json:"propertyVersion"`
		Etag            string      `json:"etag"`
		RuleFormat      string      `json:"ruleFormat"`
		Rules           Rules       `json:"rules"`
		Errors          []RuleError `json:"errors"`
	}

	// RuleError represents and entry in error field from PUT /rules response body
	RuleError struct {
		Type         string `json:"type"`
		Title        string `json:"title"`
		Detail       string `json:"detail"`
		Instance     string `json:"instance"`
		BehaviorName string `json:"behaviorName"`
	}

	// RuleOptionsMap is a type wrapping map[string]interface{} used for adding rule options
	RuleOptionsMap map[string]interface{}

	// RuleCriteriaMustSatisfy represents criteriaMustSatisfy field values
	RuleCriteriaMustSatisfy string
)

const (
	// RuleValidateModeFast const
	RuleValidateModeFast = "fast"
	// RuleValidateModeFull const
	RuleValidateModeFull = "full"

	// RuleCriteriaMustSatisfyAll const
	RuleCriteriaMustSatisfyAll RuleCriteriaMustSatisfy = "all"
	//RuleCriteriaMustSatisfyAny const
	RuleCriteriaMustSatisfyAny RuleCriteriaMustSatisfy = "any"
)

// Validate validates GetRuleTreeRequest struct
func (r GetRuleTreeRequest) Validate() error {
	return validation.Errors{
		"PropertyID":      validation.Validate(r.PropertyID, validation.Required),
		"PropertyVersion": validation.Validate(r.PropertyVersion, validation.Required),
		"ValidateMode":    validation.Validate(r.ValidateMode, validation.In(RuleValidateModeFast, RuleValidateModeFull)),
	}.Filter()
}

// Validate validates UpdateRulesRequest struct
func (r UpdateRulesRequest) Validate() error {
	return validation.Errors{
		"PropertyID":      validation.Validate(r.PropertyID, validation.Required),
		"PropertyVersion": validation.Validate(r.PropertyVersion, validation.Required),
		"ValidateMode":    validation.Validate(r.ValidateMode, validation.In(RuleValidateModeFast, RuleValidateModeFull)),
		"Rules":           validation.Validate(r.Rules),
	}.Filter()
}

// Validate validates RulesUpdate struct
func (r RulesUpdate) Validate() error {
	return validation.Errors{
		"Rules": validation.Validate(r.Rules),
	}.Filter()
}

// Validate validates Rules struct
func (r Rules) Validate() error {
	return validation.Errors{
		"Behaviors":      validation.Validate(r.Behaviors, validation.Required),
		"Name":           validation.Validate(r.Name, validation.Required),
		"CustomOverride": validation.Validate(r.CustomOverride),
		"Criteria":       validation.Validate(r.Criteria),
		"Children":       validation.Validate(r.Children),
		"Variables":      validation.Validate(r.Variables),
	}.Filter()
}

// Validate validates RuleBehavior struct
func (b RuleBehavior) Validate() error {
	return validation.Errors{
		"Name":    validation.Validate(b.Name, validation.Required),
		"Options": validation.Validate(b.Options, validation.Required),
	}.Filter()
}

// Validate validates RuleCustomOverride struct
func (co RuleCustomOverride) Validate() error {
	return validation.Errors{
		"Name":       validation.Validate(co.Name, validation.Required),
		"OverrideID": validation.Validate(co.OverrideID, validation.Required),
	}.Filter()
}

// Validate validates RuleVariable struct
func (v RuleVariable) Validate() error {
	return validation.Errors{
		"Name": validation.Validate(v.Name, validation.Required),
	}.Filter()
}

var (
	ErrGetRuleTree    = errors.New("fetchign rule tree")
	ErrUpdateRuleTree = errors.New("updating rule tree")
)

func (p *papi) GetRuleTree(ctx context.Context, params GetRuleTreeRequest) (*GetRuleTreeResponse, error) {
	if err := params.Validate(); err != nil {
		return nil, fmt.Errorf("%s: %w: %s", ErrGetRuleTree, ErrStructValidation, err)
	}

	logger := p.Log(ctx)
	logger.Debug("GetRuleTree")

	getURL := fmt.Sprintf(
		"/papi/v1/properties/%s/versions/%d/rules?contractId=%s&groupId=%s",
		params.PropertyID,
		params.PropertyVersion,
		params.ContractID,
		params.GroupID,
	)
	if params.ValidateMode != "" {
		getURL += fmt.Sprintf("&validateMode=%s", params.ValidateMode)
	}
	if !params.ValidateRules {
		getURL += fmt.Sprintf("&validateRules=%t", params.ValidateRules)
	}
	req, err := http.NewRequestWithContext(ctx, http.MethodGet, getURL, nil)
	if err != nil {
		return nil, fmt.Errorf("%w: failed to create request: %s", ErrGetRuleTree, err)
	}

	var rules GetRuleTreeResponse
	resp, err := p.Exec(req, &rules)
	if err != nil {
		return nil, fmt.Errorf("%w: request failed: %s", ErrGetRuleTree, err)
	}

	if resp.StatusCode != http.StatusOK {
		return nil, fmt.Errorf("%s: %w", ErrGetRuleTree, p.Error(resp))
	}

	return &rules, nil
}

func (p *papi) UpdateRuleTree(ctx context.Context, request UpdateRulesRequest) (*UpdateRulesResponse, error) {
	if err := request.Validate(); err != nil {
		return nil, fmt.Errorf("%s: %w: %s", ErrUpdateRuleTree, ErrStructValidation, err)
	}

	logger := p.Log(ctx)
	logger.Debug("UpdateRuleTree")

	putURL := fmt.Sprintf(
		"/papi/v1/properties/%s/versions/%d/rules?contractId=%s&groupId=%s",
		request.PropertyID,
		request.PropertyVersion,
		request.ContractID,
		request.GroupID,
	)
	if request.ValidateMode != "" {
		putURL += fmt.Sprintf("&validateMode=%s", request.ValidateMode)
	}
	if !request.ValidateRules {
		putURL += fmt.Sprintf("&validateRules=%t", request.ValidateRules)
	}
	if request.DryRun {
		putURL += fmt.Sprintf("&dryRun=%t", request.DryRun)
	}
	req, err := http.NewRequestWithContext(ctx, http.MethodPut, putURL, nil)
	if err != nil {
		return nil, fmt.Errorf("%w: failed to create request: %s", ErrUpdateRuleTree, err)
	}

	var versions UpdateRulesResponse
	resp, err := p.Exec(req, &versions, request.Rules)
	if err != nil {
		return nil, fmt.Errorf("%w: request failed: %s", ErrUpdateRuleTree, err)
	}
	if resp.StatusCode != http.StatusOK {
		return nil, fmt.Errorf("%s: %w", ErrUpdateRuleTree, p.Error(resp))
	}

	return &versions, nil
}
