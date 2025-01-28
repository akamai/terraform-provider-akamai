// Package meta contains code related to provider's meta information
package meta

import (
	"errors"
	"fmt"

	akalog "github.com/akamai/AkamaiOPEN-edgegrid-golang/v10/pkg/log"
	"github.com/akamai/AkamaiOPEN-edgegrid-golang/v10/pkg/session"
	"github.com/akamai/terraform-provider-akamai/v6/pkg/log"
	"github.com/hashicorp/go-hclog"
)

var _ Meta = &OperationMeta{}

type (
	// Meta is the akamai meta object interface
	Meta interface {
		// Log constructs an hclog sublogger and returns the log.Interface
		Log(args ...interface{}) akalog.Interface

		// OperationID returns the operation id
		OperationID() string

		// Session returns the operation API session
		Session() session.Session
	}

	// OperationMeta is the implementation of Meta interface
	OperationMeta struct {
		operationID string
		log         hclog.Logger
		sess        session.Session
	}
)

// ErrNilLog is an error returned from New(...) when log argument is nil
var ErrNilLog = errors.New("nil log argument")

// ErrNilSession is an error returned from New(...) when session argument is nil
var ErrNilSession = errors.New("nil session argument")

// New returns a new OperationMeta
func New(sess session.Session, log hclog.Logger, operationID string) (*OperationMeta, error) {
	if log == nil {
		return nil, ErrNilLog
	}
	if sess == nil {
		return nil, ErrNilSession
	}
	return &OperationMeta{
		operationID: operationID,
		sess:        sess,
		log:         log,
	}, nil
}

// Must performs type assertion on m and panics if m does not hold Meta value
func Must(m any) Meta {
	v, ok := m.(Meta)
	if !ok {
		panic(fmt.Sprintf("%v does not implement Meta interface", m))
	}
	return v
}

// Log creates a logger for the provider from the meta
func (m *OperationMeta) Log(args ...interface{}) akalog.Interface {
	return log.FromHCLog(m.log.With(args...))
}

// OperationID returns the operation id from the meta
func (m *OperationMeta) OperationID() string {
	return m.operationID
}

// Session returns the meta session
func (m *OperationMeta) Session() session.Session {
	return m.sess
}
