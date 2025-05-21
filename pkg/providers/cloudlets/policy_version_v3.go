package cloudlets

import (
	"context"

	cloudlets "github.com/akamai/AkamaiOPEN-edgegrid-golang/v11/pkg/cloudlets/v3"
)

type v3VersionStrategy struct {
	client cloudlets.Cloudlets
}

func getAllPolicyVersionsV3(ctx context.Context, policyID int64, client cloudlets.Cloudlets) ([]cloudlets.ListPolicyVersionsItem, error) {
	size, page := 1000, 0
	allPolicyVersions := make([]cloudlets.ListPolicyVersionsItem, 0)

	for {
		versions, err := client.ListPolicyVersions(ctx, cloudlets.ListPolicyVersionsRequest{
			PolicyID: policyID,
			Page:     page,
			Size:     size,
		})
		if err != nil {
			return nil, err
		}

		allPolicyVersions = append(allPolicyVersions, versions.PolicyVersions...)
		if len(versions.PolicyVersions) < size {
			break
		}
		page++
	}

	return allPolicyVersions, nil
}

func (v3 v3VersionStrategy) findLatestPolicyVersion(ctx context.Context, policyID int64) (*int64, error) {
	versions, err := getAllPolicyVersionsV3(ctx, policyID, v3.client)
	if err != nil {
		return nil, err
	}
	if len(versions) == 0 {
		return nil, nil
	}
	//API returns list of versions sorted in descending order, and it can be assumed that first element is the latest version
	return &versions[0].PolicyVersion, nil
}
