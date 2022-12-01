package cloudlets

import (
	"context"
	"fmt"
	"testing"

	"github.com/akamai/AkamaiOPEN-edgegrid-golang/v3/pkg/cloudlets"
	"github.com/stretchr/testify/mock"
	"github.com/stretchr/testify/require"
	"github.com/tj/assert"
)

func TestFindingLatestPolicyVersion(t *testing.T) {
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
		expected  int64
		withError bool
	}{
		"last policy version on 1st page found": {
			init: func(m *cloudlets.Mock) {
				m.On("ListPolicyVersions", mock.Anything, cloudlets.ListPolicyVersionsRequest{
					PolicyID: policyID,
					PageSize: &pageSize,
					Offset:   0,
				}).Return(preparePolicyVersionsPage(1000, 0), nil).Once()
				m.On("ListPolicyVersions", mock.Anything, cloudlets.ListPolicyVersionsRequest{
					PolicyID: policyID,
					PageSize: &pageSize,
					Offset:   1000,
				}).Return([]cloudlets.PolicyVersion{}, nil).Once()
			},
			expected: 999,
		},
		"first policy version on 1st page found": {
			init: func(m *cloudlets.Mock) {
				policyVersionsPage := preparePolicyVersionsPage(500, 0)
				policyVersionsPage[0] = cloudlets.PolicyVersion{Version: 500}
				m.On("ListPolicyVersions", mock.Anything, cloudlets.ListPolicyVersionsRequest{
					PolicyID: policyID,
					PageSize: &pageSize,
					Offset:   0,
				}).Return(policyVersionsPage, nil).Once()
			},
			expected: 500,
		},
		"policy version on 3rd page found": {
			init: func(m *cloudlets.Mock) {
				m.On("ListPolicyVersions", mock.Anything, cloudlets.ListPolicyVersionsRequest{
					PolicyID: policyID,
					PageSize: &pageSize,
					Offset:   0,
				}).Return(preparePolicyVersionsPage(1000, 0), nil).Once()
				m.On("ListPolicyVersions", mock.Anything, cloudlets.ListPolicyVersionsRequest{
					PolicyID: policyID,
					PageSize: &pageSize,
					Offset:   1000,
				}).Return(preparePolicyVersionsPage(1000, 1000), nil).Once()
				m.On("ListPolicyVersions", mock.Anything, cloudlets.ListPolicyVersionsRequest{
					PolicyID: policyID,
					PageSize: &pageSize,
					Offset:   2000,
				}).Return(preparePolicyVersionsPage(500, 2000), nil).Once()
			},
			expected: 2499,
		},
		"no policy versions found": {
			init: func(m *cloudlets.Mock) {
				m.On("ListPolicyVersions", mock.Anything, cloudlets.ListPolicyVersionsRequest{
					PolicyID: policyID,
					PageSize: &pageSize,
					Offset:   0,
				}).Return([]cloudlets.PolicyVersion{}, nil).Once()
			},
			withError: true,
		},
		"error listing policy versions": {
			init: func(m *cloudlets.Mock) {
				m.On("ListPolicyVersions", mock.Anything, cloudlets.ListPolicyVersionsRequest{
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
			m := &cloudlets.Mock{}
			test.init(m)
			version, err := findLatestPolicyVersion(context.Background(), policyID, m)
			if test.withError {
				assert.Error(t, err)
			} else {
				require.NoError(t, err)
				assert.Equal(t, test.expected, version)
			}
		})
	}
}
