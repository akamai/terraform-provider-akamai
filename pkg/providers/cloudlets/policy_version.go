package cloudlets

import (
	"context"
	"fmt"

	"github.com/akamai/AkamaiOPEN-edgegrid-golang/v6/pkg/cloudlets"
)

func getAllPolicyVersions(ctx context.Context, policyID int64, client cloudlets.Cloudlets) ([]cloudlets.PolicyVersion, error) {
	pageSize, offset := 1000, 0
	allPolicyVersions := make([]cloudlets.PolicyVersion, 0)

	for {
		versions, err := client.ListPolicyVersions(ctx, cloudlets.ListPolicyVersionsRequest{
			PolicyID:     policyID,
			IncludeRules: false,
			PageSize:     &pageSize,
			Offset:       offset,
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

func findLatestPolicyVersion(ctx context.Context, policyID int64, client cloudlets.Cloudlets) (int64, error) {
	var version int64
	versions, err := getAllPolicyVersions(ctx, policyID, client)
	if err != nil {
		return version, err
	}
	if len(versions) == 0 {
		return version, fmt.Errorf("no policy version found")
	}

	for _, v := range versions {
		if v.Version > version {
			version = v.Version
		}
	}

	return version, nil
}
