package nullstringdefault

import (
	"context"
	"testing"

	"github.com/google/go-cmp/cmp"
	"github.com/hashicorp/terraform-plugin-framework/resource/schema/defaults"
	"github.com/hashicorp/terraform-plugin-framework/types"
)

func TestStaticNullStringDefaultString(t *testing.T) {
	t.Parallel()

	testCases := map[string]struct {
		expected *defaults.StringResponse
	}{
		"null": {
			expected: &defaults.StringResponse{
				PlanValue: types.StringNull(),
			},
		},
	}

	for name, testCase := range testCases {
		t.Run(name, func(t *testing.T) {
			t.Parallel()

			resp := &defaults.StringResponse{}

			NullString().DefaultString(context.Background(), defaults.StringRequest{}, resp)

			if diff := cmp.Diff(testCase.expected, resp); diff != "" {
				t.Errorf("unexpected difference: %s", diff)
			}
		})
	}
}
