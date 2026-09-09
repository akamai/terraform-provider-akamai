package cloudlets

import (
	"context"
	"fmt"
	"testing"

	"github.com/akamai/AkamaiOPEN-edgegrid-golang/v14/pkg/cloudlets"
	v3 "github.com/akamai/AkamaiOPEN-edgegrid-golang/v14/pkg/cloudlets/v3"
	"github.com/akamai/terraform-provider-akamai/v11/internal/edgegrid"
	"github.com/akamai/terraform-provider-akamai/v11/pkg/common/ptr"
	"github.com/akamai/terraform-provider-akamai/v11/pkg/common/testutils"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

func TestFindingLatestPolicyVersion(t *testing.T) {
	t.Parallel()
	preparePolicyVersionsPage := func(pageSize, startingVersion int64) []cloudlets.PolicyVersion {
		versions := make([]cloudlets.PolicyVersion, 0, pageSize)
		for i := startingVersion; i < startingVersion+pageSize; i++ {
			versions = append(versions, cloudlets.PolicyVersion{Version: i})
		}
		return versions
	}

	var policyID int64 = 123
	pageSize := 1000

	tests := map[string]struct {
		init      func(m *cloudlets.Mock)
		expected  *int64
		withError bool
	}{
		"last policy version on 1st page found": {
			init: func(m *cloudlets.Mock) {
				m.On("ListPolicyVersions", testutils.MockContext, cloudlets.ListPolicyVersionsRequest{
					PolicyID: policyID,
					PageSize: &pageSize,
					Offset:   0,
				}).Return(preparePolicyVersionsPage(1000, 0), nil).Once()
				m.On("ListPolicyVersions", testutils.MockContext, cloudlets.ListPolicyVersionsRequest{
					PolicyID: policyID,
					PageSize: &pageSize,
					Offset:   1000,
				}).Return([]cloudlets.PolicyVersion{}, nil).Once()
			},
			expected: ptr.To(int64(0)),
		},
		"first policy version on 1st page found": {
			init: func(m *cloudlets.Mock) {
				policyVersionsPage := preparePolicyVersionsPage(500, 0)
				policyVersionsPage[0] = cloudlets.PolicyVersion{Version: 500}
				m.On("ListPolicyVersions", testutils.MockContext, cloudlets.ListPolicyVersionsRequest{
					PolicyID: policyID,
					PageSize: &pageSize,
					Offset:   0,
				}).Return(policyVersionsPage, nil).Once()
			},
			expected: ptr.To(int64(500)),
		},
		"no policy versions found": {
			init: func(m *cloudlets.Mock) {
				m.On("ListPolicyVersions", testutils.MockContext, cloudlets.ListPolicyVersionsRequest{
					PolicyID: policyID,
					PageSize: &pageSize,
					Offset:   0,
				}).Return([]cloudlets.PolicyVersion{}, nil).Once()
			},
			expected: nil,
		},
		"error listing policy versions": {
			init: func(m *cloudlets.Mock) {
				m.On("ListPolicyVersions", testutils.MockContext, cloudlets.ListPolicyVersionsRequest{
					PolicyID: policyID,
					PageSize: &pageSize,
					Offset:   0,
				}).Return(nil, fmt.Errorf("oops")).Once()
			},
			withError: true,
		},
	}

	for name, test := range tests {
		t.Run(name, func(t *testing.T) {
			t.Parallel()
			client := edgegrid.NewTestClient()
			test.init(client.CloudletsV2)
			versionStrategy := v2VersionStrategy{client: client.CloudletsV2}
			version, err := versionStrategy.findLatestPolicyVersion(context.Background(), policyID)
			if test.withError {
				assert.Error(t, err)
			} else {
				require.NoError(t, err)
				if test.expected != nil {
					assert.Equal(t, *test.expected, *version)
				} else {
					assert.Nil(t, version)
				}
			}
			client.CloudletsV2.AssertExpectations(t)
		})
	}
}

func TestFindingLatestPolicyVersionV3(t *testing.T) {
	t.Parallel()
	preparePolicyVersionsPage := func(pageSize, startingVersion int64) []v3.ListPolicyVersionsItem {
		versions := make([]v3.ListPolicyVersionsItem, 0, pageSize)
		for i := startingVersion; i < startingVersion+pageSize; i++ {
			versions = append(versions, v3.ListPolicyVersionsItem{PolicyVersion: i})
		}
		return versions
	}

	var policyID int64 = 123
	pageSize := 1000

	tests := map[string]struct {
		init      func(m *v3.Mock)
		expected  *int64
		withError bool
	}{
		"last policy version on 1st page found": {
			init: func(m *v3.Mock) {
				m.On("ListPolicyVersions", testutils.MockContext, v3.ListPolicyVersionsRequest{
					PolicyID: policyID,
					Size:     pageSize,
					Page:     0,
				}).Return(&v3.ListPolicyVersions{PolicyVersions: preparePolicyVersionsPage(1000, 0)}, nil).Once()
				m.On("ListPolicyVersions", testutils.MockContext, v3.ListPolicyVersionsRequest{
					PolicyID: policyID,
					Size:     pageSize,
					Page:     1,
				}).Return(&v3.ListPolicyVersions{PolicyVersions: []v3.ListPolicyVersionsItem{}}, nil).Once()
			},
			expected: ptr.To(int64(0)),
		},
		"first policy version on 1st page found": {
			init: func(m *v3.Mock) {
				policyVersionsPage := preparePolicyVersionsPage(500, 0)
				policyVersionsPage[0] = v3.ListPolicyVersionsItem{PolicyVersion: 500}
				m.On("ListPolicyVersions", testutils.MockContext, v3.ListPolicyVersionsRequest{
					PolicyID: policyID,
					Size:     pageSize,
					Page:     0,
				}).Return(&v3.ListPolicyVersions{PolicyVersions: policyVersionsPage}, nil).Once()
			},
			expected: ptr.To(int64(500)),
		},
		"no policy versions found": {
			init: func(m *v3.Mock) {
				m.On("ListPolicyVersions", testutils.MockContext, v3.ListPolicyVersionsRequest{
					PolicyID: policyID,
					Size:     pageSize,
					Page:     0,
				}).Return(&v3.ListPolicyVersions{PolicyVersions: []v3.ListPolicyVersionsItem{}}, nil).Once()
			},
			expected: nil,
		},
		"error listing policy versions": {
			init: func(m *v3.Mock) {
				m.On("ListPolicyVersions", testutils.MockContext, v3.ListPolicyVersionsRequest{
					PolicyID: policyID,
					Size:     pageSize,
					Page:     0,
				}).Return(nil, fmt.Errorf("oops")).Once()
			},
			withError: true,
		},
	}

	for name, test := range tests {
		t.Run(name, func(t *testing.T) {
			t.Parallel()
			client := edgegrid.NewTestClient()
			test.init(client.CloudletsV3)
			checker := v3VersionStrategy{client: client.CloudletsV3}
			version, err := checker.findLatestPolicyVersion(context.Background(), policyID)
			if test.withError {
				assert.Error(t, err)
			} else {
				require.NoError(t, err)
				if test.expected != nil {
					assert.Equal(t, *test.expected, *version)
				} else {
					assert.Nil(t, version)
				}
			}
			client.CloudletsV3.AssertExpectations(t)
		})
	}
}
