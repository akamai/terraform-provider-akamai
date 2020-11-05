// Package papi provides access to the Akamai Property APIs
package papi

import (
	"errors"
	"net/http"

	"github.com/akamai/AkamaiOPEN-edgegrid-golang/v2/pkg/session"
	"github.com/spf13/cast"
)

var (
	// ErrStructValidation is returned returned when given struct validation failed
	ErrStructValidation = errors.New("struct validation")

	// ErrNotFound is returned when requested resource was not found
	ErrNotFound = errors.New("resource not found")
)

type (
	// PAPI is the papi api interface
	PAPI interface {
		Groups
		Contracts
		Activations
		CPCodes
		Properties
		PropertyVersions
		EdgeHostnames
		Products
		Search
		PropertyVersionHostnames
		ClientSettings
		PropertyRules
		RuleFormats
	}

	papi struct {
		session.Session
		usePrefixes bool
	}

	// Option defines a PAPI option
	Option func(*papi)

	// ClientFunc is a papi client new method, this can used for mocking
	ClientFunc func(sess session.Session, opts ...Option) PAPI

	// Response is a base PAPI Response type
	Response struct {
		AccountID  string   `json:"accountId,omitempty"`
		ContractID string   `json:"contractId,omitempty"`
		GroupID    string   `json:"groupId,omitempty"`
		Etag       string   `json:"etag,omitempty"`
		Errors     []*Error `json:"errors,omitempty"`
		Warnings   []*Error `json:"warnings,omitempty"`
	}
)

// Client returns a new papi Client instance with the specified controller
func Client(sess session.Session, opts ...Option) PAPI {
	p := &papi{
		Session:     sess,
		usePrefixes: true,
	}

	for _, opt := range opts {
		opt(p)
	}
	return p
}

// WithUsePrefixes sets the `PAPI-Use-Prefixes` header on requests
// See: https://developer.akamai.com/api/core_features/property_manager/v1.html#prefixes
func WithUsePrefixes(usePrefixes bool) Option {
	return func(p *papi) {
		p.usePrefixes = usePrefixes
	}
}

// Exec overrides the session.Exec to add papi options
func (p *papi) Exec(r *http.Request, out interface{}, in ...interface{}) (*http.Response, error) {
	// explicitly add the PAPI-Use-Prefixes header
	r.Header.Set("PAPI-Use-Prefixes", cast.ToString(p.usePrefixes))

	return p.Session.Exec(r, out, in...)
}
