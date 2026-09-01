package mtlstruststore

import (
	"context"
	"fmt"
	"testing"

	"github.com/akamai/AkamaiOPEN-edgegrid-golang/v14/pkg/mtlstruststore"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/mock"
	"github.com/stretchr/testify/require"
)

func TestFindNotDeletedCASet(t *testing.T) {
	t.Parallel()
	type args struct {
		caSetName string
		caSets    *mtlstruststore.ListCASetsResponse
		err       error
	}
	tests := map[string]struct {
		args          args
		expectedCASet mtlstruststore.CASetResponse
		expectedErr   string
	}{
		"single match": {
			args: args{
				caSetName: "test-ca-set",
				caSets: &mtlstruststore.ListCASetsResponse{
					CASets: []mtlstruststore.CASetResponse{
						{CASetID: "id-1", CASetName: "test-ca-set", CASetStatus: "NOT_DELETED"},
						{CASetID: "id-2", CASetName: "other", CASetStatus: "NOT_DELETED"},
					},
				},
			},
			expectedCASet: mtlstruststore.CASetResponse{CASetID: "id-1", CASetName: "test-ca-set", CASetStatus: "NOT_DELETED"},
		},
		"single match with equal prefixes": {
			args: args{
				caSetName: "test-ca-set",
				caSets: &mtlstruststore.ListCASetsResponse{
					CASets: []mtlstruststore.CASetResponse{
						{CASetID: "id-1", CASetName: "test-ca-set", CASetStatus: "NOT_DELETED"},
						{CASetID: "id-2", CASetName: "test-ca-set-2", CASetStatus: "NOT_DELETED"},
					},
				},
			},
			expectedCASet: mtlstruststore.CASetResponse{CASetID: "id-1", CASetName: "test-ca-set", CASetStatus: "NOT_DELETED"},
		},
		"no match": {
			args: args{
				caSetName: "notfound",
				caSets: &mtlstruststore.ListCASetsResponse{
					CASets: []mtlstruststore.CASetResponse{
						{CASetID: "id-1", CASetName: "test-ca-set", CASetStatus: "NOT_DELETED"},
					},
				},
			},
			expectedErr: "no CA set found with the name 'notfound' and status 'NOT_DELETED'",
		},
		"no match - empty list": {
			args: args{
				caSetName: "foo",
				caSets: &mtlstruststore.ListCASetsResponse{
					CASets: []mtlstruststore.CASetResponse{},
				},
			},
			expectedErr: "no CA set found with the name 'foo' and status 'NOT_DELETED'",
		},
		"multiple matches": {
			args: args{
				caSetName: "test-ca-set",
				caSets: &mtlstruststore.ListCASetsResponse{
					CASets: []mtlstruststore.CASetResponse{
						{CASetID: "id-1", CASetName: "test-ca-set", CASetStatus: "NOT_DELETED"},
						{CASetID: "id-2", CASetName: "test-ca-set", CASetStatus: "NOT_DELETED"},
					},
				},
			},
			expectedErr: `multiple CA sets found with the name 'test-ca-set' and status 'NOT_DELETED': \[id-1 id-2\]`,
		},
		"API error": {
			args: args{
				caSetName: "test-ca-set",
				err:       fmt.Errorf("API failure"),
			},
			expectedErr: "could not find CA set with the name 'test-ca-set' and status 'NOT_DELETED'. API error: API failure",
		},
	}

	for name, tc := range tests {
		t.Run(name, func(t *testing.T) {
			t.Parallel()
			m := &mtlstruststore.Mock{}
			m.On("ListCASets", mock.Anything, mtlstruststore.ListCASetsRequest{
				CASetNamePrefix: tc.args.caSetName,
				CASetStatuses:   []string{mtlstruststore.CASetStatusNotDeleted},
			}).Return(tc.args.caSets, tc.args.err).Once()

			caSet, err := findNotDeletedCASetByName(context.Background(), m, tc.args.caSetName)
			if tc.expectedErr != "" {
				require.Error(t, err)
				assert.Regexp(t, tc.expectedErr, err.Error())
			} else {
				require.NoError(t, err)
				assert.Equal(t, tc.expectedCASet, caSet)
			}
			m.AssertExpectations(t)
		})
	}
}
