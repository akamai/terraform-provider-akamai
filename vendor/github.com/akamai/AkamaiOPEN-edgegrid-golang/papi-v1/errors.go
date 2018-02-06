package papi

import "errors"

// Error constants
const (
	ErrInvalidPath = iota
	ErrCriteriaNotFound
	ErrBehaviorNotFound
	ErrVariableNotFound
	ErrRuleNotFound
	ErrInvalidRules
)

var (
	ErrorMap = map[int]error{
		ErrInvalidPath:      errors.New("Invalid Path"),
		ErrCriteriaNotFound: errors.New("Criteria not found"),
		ErrBehaviorNotFound: errors.New("Behavior not found"),
		ErrVariableNotFound: errors.New("Variable not found"),
		ErrRuleNotFound:     errors.New("Rule not found"),
		ErrInvalidRules:     errors.New("Rule validation failed. See papi.Rules.Errors for details"),
	}
)
