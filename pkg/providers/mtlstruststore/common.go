package mtlstruststore

import (
	"context"
	"fmt"

	"github.com/akamai/AkamaiOPEN-edgegrid-golang/v13/pkg/mtlstruststore"
)

// findNotDeletedCASetByName finds a single not-deleted CA set by exact name match.
func findNotDeletedCASetByName(ctx context.Context, client mtlstruststore.MTLSTruststore, caSetName string) (mtlstruststore.CASetResponse, error) {
	caSets, err := client.ListCASets(ctx, mtlstruststore.ListCASetsRequest{
		CASetNamePrefix: caSetName,
		CASetStatuses:   []string{mtlstruststore.CASetStatusNotDeleted},
	})
	if err != nil {
		return mtlstruststore.CASetResponse{},
			fmt.Errorf("could not find CA set with the name '%s' and status '%s'. API error: %w",
				caSetName, mtlstruststore.CASetStatusNotDeleted, err)
	}

	var matchingSets []mtlstruststore.CASetResponse
	var matchingIDs []string
	for _, caSet := range caSets.CASets {
		if caSet.CASetName == caSetName {
			matchingSets = append(matchingSets, caSet)
			matchingIDs = append(matchingIDs, caSet.CASetID)
		}
	}

	switch len(matchingSets) {
	case 0:
		return mtlstruststore.CASetResponse{},
			fmt.Errorf("no CA set found with the name '%s' and status '%s'",
				caSetName, mtlstruststore.CASetStatusNotDeleted)
	case 1:
		return matchingSets[0], nil
	default:
		return mtlstruststore.CASetResponse{},
			fmt.Errorf("multiple CA sets found with the name '%s' and status '%s': %v. "+
				"Use the ID to fetch a specific CA set",
				caSetName, mtlstruststore.CASetStatusNotDeleted, matchingIDs)
	}
}
