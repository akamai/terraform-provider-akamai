package cloudlets

import (
	"context"

	"github.com/akamai/terraform-provider-akamai/v8/pkg/common/tf"
	"github.com/akamai/terraform-provider-akamai/v8/pkg/meta"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/schema"
)

func getPolicyVersionExecutionStrategy(d *schema.ResourceData, meta meta.Meta) (versionStrategy, error) {
	isV3, err := tf.GetBoolValue("is_shared", d)
	if err != nil {
		return nil, err
	}

	if isV3 {
		return v3VersionStrategy{ClientV3(meta)}, nil
	}
	return v2VersionStrategy{Client(meta)}, nil
}

type versionStrategy interface {
	findLatestPolicyVersion(ctx context.Context, policyID int64) (*int64, error)
}
