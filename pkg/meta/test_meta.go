package meta

import "github.com/akamai/terraform-provider-akamai/v11/internal/edgegrid"

// NewTestMeta returns an OperationMeta with client pre-set, for wiring actions/resources to a mock edgegrid
// client in tests. The session and log fields are left at their zero value: Configure() only ever reads the
// client, and production code never observes a test-constructed Meta.
func NewTestMeta(client edgegrid.Client) *OperationMeta {
	return &OperationMeta{client: client}
}
