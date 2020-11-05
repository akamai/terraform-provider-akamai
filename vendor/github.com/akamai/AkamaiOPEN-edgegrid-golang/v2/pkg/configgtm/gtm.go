// Package gtm provides access to the Akamai GTM V1_4 APIs
package gtm

import (
	"errors"
	"net/http"

	"github.com/akamai/AkamaiOPEN-edgegrid-golang/v2/pkg/session"
)

var (
	// ErrStructValidation is returned returned when given struct validation failed
	ErrStructValidation = errors.New("struct validation")
)

type (
	// GTM is the gtm api interface
	GTM interface {
		Domains
		Properties
		Datacenters
		Resources
		ASMaps
		GeoMaps
		CidrMaps
	}

	gtm struct {
		session.Session
	}

	// Option defines a GTM option
	Option func(*gtm)

	// ClientFunc is a gtm client new method, this can used for mocking
	ClientFunc func(sess session.Session, opts ...Option) GTM
)

// Client returns a new dns Client instance with the specified controller
func Client(sess session.Session, opts ...Option) GTM {
	p := &gtm{
		Session: sess,
	}

	for _, opt := range opts {
		opt(p)
	}
	return p
}

// Exec overrides the session.Exec to add dns options
func (p *gtm) Exec(r *http.Request, out interface{}, in ...interface{}) (*http.Response, error) {

	return p.Session.Exec(r, out, in...)
}
