package nullint64default

import (
	"context"
	"testing"

	"github.com/google/go-cmp/cmp"
	"github.com/hashicorp/terraform-plugin-framework/resource/schema/defaults"
	"github.com/hashicorp/terraform-plugin-framework/types"
)

func TestStaticInt64DefaultInt64(t *testing.T) {
	t.Parallel()

	testCases := map[string]struct {
		expected *defaults.Int64Response
	}{
		"null": {
			expected: &defaults.Int64Response{
				PlanValue: types.Int64Null(),
			},
		},
	}

	for name, testCase := range testCases {
		t.Run(name, func(t *testing.T) {
			t.Parallel()

			resp := &defaults.Int64Response{}

			NullInt64().DefaultInt64(context.Background(), defaults.Int64Request{}, resp)

			if diff := cmp.Diff(testCase.expected, resp); diff != "" {
				t.Errorf("unexpected difference: %s", diff)
			}
		})
	}
}
