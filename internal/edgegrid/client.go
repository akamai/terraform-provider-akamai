// Package edgegrid provides a client for interacting with Akamai's Open API.
package edgegrid

import (
	"github.com/akamai/AkamaiOPEN-edgegrid-golang/v12/pkg/cloudcertificates"
	"github.com/akamai/AkamaiOPEN-edgegrid-golang/v12/pkg/session"
)

// Client is the interface for the Akamai Edgegrid client.
type Client interface {
	GetCloudCertificates() cloudcertificates.CloudCertificates
}

var (
	_ Client = &ClientImpl{}
)

// ClientImpl is the implementation of the Client interface used in production.
type ClientImpl struct {
	sess session.Session
}

// NewClientImpl creates a new instance of ClientImpl with the provided session.
func NewClientImpl(sess session.Session) *ClientImpl {
	return &ClientImpl{
		sess: sess,
	}
}

// GetCloudCertificates returns the CCM client for managing cloud certificates.
func (c *ClientImpl) GetCloudCertificates() cloudcertificates.CloudCertificates {
	return cloudcertificates.Client(c.sess)
}
