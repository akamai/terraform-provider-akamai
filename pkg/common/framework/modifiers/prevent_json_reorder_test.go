package modifiers

import (
	"context"
	"testing"

	"github.com/hashicorp/terraform-plugin-framework/resource/schema/planmodifier"
	"github.com/hashicorp/terraform-plugin-framework/types"
	"github.com/stretchr/testify/require"
)

func TestPreventJsonReorder_SuppressesWhenOnlyOrderDiffers(t *testing.T) {
	t.Parallel()

	mod := PreventJSONReorder()
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

	mod := PreventJSONReorder()
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

	mod := PreventJSONReorder()
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

	mod := PreventJSONReorder()

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

	require.Equal(t, state.ValueString(), resp.PlanValue.ValueString())
}

func TestPreventJsonReorder_DetectsDifferentMultiplicity(t *testing.T) {
	t.Parallel()

	mod := PreventJSONReorder()

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

	require.Equal(t, plan.ValueString(), resp.PlanValue.ValueString())
}

func TestPreventJsonReorder_SuppressesNestedArrayReorder(t *testing.T) {
	t.Parallel()

	mod := PreventJSONReorder()

	// Simulates a rule_definitions entry where conditions.hosts is reordered
	state := types.StringValue(`[{"id":1,"conditionException":{"conditions":[{"type":"hostMatch","hosts":["a.com","b.com"]}]}}]`)
	plan := types.StringValue(`[{"id":1,"conditionException":{"conditions":[{"type":"hostMatch","hosts":["b.com","a.com"]}]}}]`)

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

func TestPreventJsonReorder_SuppressesObjectJSON(t *testing.T) {
	t.Parallel()

	// PreventJSONReorder now works for JSON objects too, not just top-level arrays.
	mod := PreventJSONReorder()

	state := types.StringValue(`{"tags":["b","a","c"],"name":"test"}`)
	plan := types.StringValue(`{"name":"test","tags":["a","c","b"]}`)

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
