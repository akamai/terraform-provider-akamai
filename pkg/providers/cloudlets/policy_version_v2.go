package cloudlets

import (
	"context"

	"github.com/akamai/AkamaiOPEN-edgegrid-golang/v11/pkg/cloudlets"
)

type v2VersionStrategy struct {
	client cloudlets.Cloudlets
}

func getAllV2PolicyVersions(ctx context.Context, policyID int64, client cloudlets.Cloudlets) ([]cloudlets.PolicyVersion, error) {
	pageSize, offset := 1000, 0
	allPolicyVersions := make([]cloudlets.PolicyVersion, 0)

	for {
		versions, err := client.ListPolicyVersions(ctx, cloudlets.ListPolicyVersionsRequest{
			PolicyID: policyID,
			PageSize: &pageSize,
			Offset:   offset,
		})
		if err != nil {
			return nil, err
		}

		allPolicyVersions = append(allPolicyVersions, versions...)
		if len(versions) < pageSize {
			break
		}
		offset += pageSize
	}

	return allPolicyVersions, nil

}

func (v2 v2VersionStrategy) findLatestPolicyVersion(ctx context.Context, policyID int64) (*int64, error) {
	versions, err := getAllV2PolicyVersions(ctx, policyID, v2.client)
	if err != nil {
		return nil, err
	}
	if len(versions) == 0 {
		return nil, nil
	}
	//API returns list of versions sorted in descending order, and it can be assumed that first element is the latest version
	return &versions[0].Version, nil
}
