package cloudcertificates

import (
	"context"
	"testing"
	"time"

	"github.com/akamai/AkamaiOPEN-edgegrid-golang/v15/pkg/cloudcertificates"
	"github.com/hashicorp/terraform-plugin-framework/action"
	"github.com/hashicorp/terraform-plugin-go/tftypes"
	"github.com/stretchr/testify/assert"
)

// This file collects response builders and action-schema helpers shared by the activation resource's and the
// Rollback/ReplaceStaging actions' test scenarios, since all three poll the same kind of lineage activation.

// activationStatus builds a GetActivationStatusResponse for tests, mirroring the real captured shape confirmed in
// edgegrid-golang's TestPromoteLineage fixtures.
func activationStatus(lineageID, generationID, activationID int64, network, status string, created, modified time.Time) *cloudcertificates.GetActivationStatusResponse {
	return &cloudcertificates.GetActivationStatusResponse{
		ActivationID:           activationID,
		LineageID:              lineageID,
		GenerationID:           generationID,
		ActivationType:         "PROMOTE",
		ActivationStatus:       status,
		TargetEnvironment:      network,
		ActivationCreatedTime:  created,
		ActivationModifiedTime: modified,
		CreatedBy:              "terraform-dev",
		ModifiedBy:             "terraform-dev",
	}
}

// timeoutsObjectType returns the TerraformType of a's "timeouts" attribute, built via timeouts.AttributesWithOpts.
func timeoutsObjectType(t *testing.T, a action.Action) tftypes.Object {
	t.Helper()

	ctx := context.Background()
	schemaResp := &action.SchemaResponse{}
	a.Schema(ctx, action.SchemaRequest{}, schemaResp)
	assert.False(t, schemaResp.Diagnostics.HasError())
	return schemaResp.Schema.Type().TerraformType(ctx).(tftypes.Object).AttributeTypes["timeouts"].(tftypes.Object)
}

// timeoutsValue builds the raw "timeouts" attribute value for an action schema built with
// timeouts.AttributesWithOpts: entirely null when invoke is nil, otherwise a populated object with just the
// "invoke" attribute set. Shared by every CCM action exposing this attribute.
func timeoutsValue(timeoutsObjectType tftypes.Object, invoke *string) tftypes.Value {
	if invoke == nil {
		return tftypes.NewValue(timeoutsObjectType, nil)
	}
	return tftypes.NewValue(timeoutsObjectType, map[string]tftypes.Value{
		"invoke": tftypes.NewValue(tftypes.String, invoke),
	})
}
