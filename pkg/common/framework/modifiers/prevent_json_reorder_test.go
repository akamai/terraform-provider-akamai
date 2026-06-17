package modifiers

import (
	"context"
	"encoding/json"
	"testing"

	"github.com/hashicorp/terraform-plugin-framework/resource/schema/planmodifier"
	"github.com/hashicorp/terraform-plugin-framework/types"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

func TestPreventJsonReorder_SuppressesWhenOnlyOrderDiffers(t *testing.T) {
	t.Parallel()

	mod := newStringSuppressJSONListReorderModifier(nil)
	state := types.StringValue(`[{"id":1,"name":"one"},{"id":2,"name":"two"}]`)
	plan := types.StringValue(`[{"name":"two","id":2},{"name":"one","id":1}]`)

	req := planmodifier.StringRequest{
		StateValue: state,
		PlanValue:  plan,
	}
	resp := planmodifier.StringResponse{
		PlanValue: plan,
	}

	mod.PlanModifyString(context.Background(), req, &resp)

	require.Equal(t, state.ValueString(), resp.PlanValue.ValueString())
}

func TestPreventJsonReorder_DoesNotSuppressWhenDifferent(t *testing.T) {
	t.Parallel()

	mod := newStringSuppressJSONListReorderModifier(nil)
	state := types.StringValue(`[{"id":1,"name":"one"}]`)
	plan := types.StringValue(`[{"id":1,"name":"two"}]`)

	req := planmodifier.StringRequest{
		StateValue: state,
		PlanValue:  plan,
	}
	resp := planmodifier.StringResponse{
		PlanValue: plan,
	}

	mod.PlanModifyString(context.Background(), req, &resp)

	require.NotEqual(t, state.ValueString(), resp.PlanValue.ValueString())
	require.Equal(t, plan.ValueString(), resp.PlanValue.ValueString())
}

func TestPreventJsonReorder_IgnoresWhenStateUnknown(t *testing.T) {
	t.Parallel()

	mod := newStringSuppressJSONListReorderModifier(nil)
	plan := types.StringValue(`[{"id":1}]`)

	req := planmodifier.StringRequest{
		StateValue: types.StringUnknown(),
		PlanValue:  plan,
	}
	resp := planmodifier.StringResponse{
		PlanValue: plan,
	}

	mod.PlanModifyString(context.Background(), req, &resp)

	require.Equal(t, plan.ValueString(), resp.PlanValue.ValueString())
}

func TestPreventJsonReorder_HandlesDuplicatesWithDifferentKeyOrder(t *testing.T) {
	t.Parallel()

	mod := newStringSuppressJSONListReorderModifier(nil)

	// Same objects appearing twice, but with different key orders
	state := types.StringValue(`[{"id":1,"name":"test"},{"id":1,"name":"test"}]`)
	plan := types.StringValue(`[{"name":"test","id":1},{"name":"test","id":1}]`)

	req := planmodifier.StringRequest{
		StateValue: state,
		PlanValue:  plan,
	}
	resp := planmodifier.StringResponse{
		PlanValue: plan,
	}

	mod.PlanModifyString(context.Background(), req, &resp)

	// Should suppress diff - both have the same two objects with multiplicity=2
	require.Equal(t, state.ValueString(), resp.PlanValue.ValueString())
}

func TestPreventJsonReorder_DetectsDifferentMultiplicity(t *testing.T) {
	t.Parallel()

	mod := newStringSuppressJSONListReorderModifier(nil)

	// Different multiplicity: state has object once, plan has it twice
	state := types.StringValue(`[{"id":1,"name":"test"}]`)
	plan := types.StringValue(`[{"name":"test","id":1},{"id":1,"name":"test"}]`)

	req := planmodifier.StringRequest{
		StateValue: state,
		PlanValue:  plan,
	}
	resp := planmodifier.StringResponse{
		PlanValue: plan,
	}

	mod.PlanModifyString(context.Background(), req, &resp)

	// Should NOT suppress diff - different multiplicity
	require.Equal(t, plan.ValueString(), resp.PlanValue.ValueString())
}

func TestRawMessageSlicesEqualIgnoringOrder(t *testing.T) {
	t.Parallel()

	tests := []struct {
		name     string
		a        []json.RawMessage
		b        []json.RawMessage
		expected bool
	}{
		{
			name: "same content different order",
			a: []json.RawMessage{
				json.RawMessage(`{"id":1,"name":"one"}`),
				json.RawMessage(`[1,2,3]`),
			},
			b: []json.RawMessage{
				json.RawMessage(`[1,2,3]`),
				json.RawMessage(`{"name":"one","id":1}`),
			},
			expected: true,
		},
		{
			name: "multiplicity matters",
			a: []json.RawMessage{
				json.RawMessage(`{"id":1}`),
				json.RawMessage(`{"id":1}`),
			},
			b: []json.RawMessage{
				json.RawMessage(`{"id":1}`),
			},
			expected: false,
		},
		{
			name: "nested objects canonicalized",
			a: []json.RawMessage{
				json.RawMessage(`{"meta":{"b":2,"a":1},"arr":[{"x":1,"y":2},3]}`),
			},
			b: []json.RawMessage{
				json.RawMessage(`{"arr":[{"y":2,"x":1},3],"meta":{"a":1,"b":2}}`),
			},
			expected: true,
		},
		{
			name: "different lengths",
			a: []json.RawMessage{
				json.RawMessage(`{"id":1}`),
			},
			b: []json.RawMessage{
				json.RawMessage(`{"id":1}`),
				json.RawMessage(`{"id":2}`),
			},
			expected: false,
		},
		{
			name: "different values",
			a: []json.RawMessage{
				json.RawMessage(`{"id":1}`),
			},
			b: []json.RawMessage{
				json.RawMessage(`{"id":2}`),
			},
			expected: false,
		},
	}

	for _, tt := range tests {
		tt := tt
		t.Run(tt.name, func(t *testing.T) {
			t.Parallel()
			assert.Equal(t, tt.expected, rawMessageSlicesEqualIgnoringOrder(tt.a, tt.b))
		})
	}
}
