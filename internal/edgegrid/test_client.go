package edgegrid

import (
	"github.com/akamai/AkamaiOPEN-edgegrid-golang/v12/pkg/cloudcertificates"
)

var (
	_ Client = &TestClient{}
)

// TestClient is a mock implementation of the Client interface for testing purposes.
type TestClient struct {
	CloudCertificates *cloudcertificates.Mock
}

// NewTestClient creates a new instance of TestClient with mock implementations.
func NewTestClient() *TestClient {
	return &TestClient{
		CloudCertificates: &cloudcertificates.Mock{},
	}
}

// GetCloudCertificates returns the mock CCM client.
func (c *TestClient) GetCloudCertificates() cloudcertificates.CloudCertificates {
	return c.CloudCertificates
}
