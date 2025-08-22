package apidefinitions

import (
	"context"
	"fmt"
	"reflect"

	v0 "github.com/akamai/AkamaiOPEN-edgegrid-golang/v11/pkg/apidefinitions/v0"
	"github.com/akamai/terraform-provider-akamai/v8/pkg/common/ptr"
	"github.com/google/go-cmp/cmp"
	"github.com/hashicorp/terraform-plugin-framework/attr"
	"github.com/hashicorp/terraform-plugin-framework/diag"
	"github.com/hashicorp/terraform-plugin-framework/types/basetypes"
	"github.com/hashicorp/terraform-plugin-go/tftypes"
	"github.com/hashicorp/terraform-plugin-log/tflog"
	orderedmap "github.com/wk8/go-ordered-map/v2"
	"golang.org/x/exp/slices"
)

var _ basetypes.StringTypable = apiStateType{}

type apiStateType struct {
	basetypes.StringType
}

// Equal returns true if the given type is equivalent.
func (t apiStateType) Equal(o attr.Type) bool {
	other, ok := o.(apiStateType)

	if !ok {
		return false
	}

	return t.StringType.Equal(other.StringType)
}

func (t apiStateType) String() string {
	return "apiStateType"
}

// ValueFromString should convert the String to a StringValuable type.
func (t apiStateType) ValueFromString(_ context.Context, in basetypes.StringValue) (basetypes.StringValuable, diag.Diagnostics) {
	value := apiStateValue{
		StringValue: in,
	}

	return value, nil
}

// ValueFromTerraform returns a Value given a tftypes.Value.
func (t apiStateType) ValueFromTerraform(ctx context.Context, in tftypes.Value) (attr.Value, error) {
	attrValue, err := t.StringType.ValueFromTerraform(ctx, in)

	if err != nil {
		return nil, err
	}

	stringValue, ok := attrValue.(basetypes.StringValue)

	if !ok {
		return nil, fmt.Errorf("unexpected value type of %T", attrValue)
	}

	stringValuable, diags := t.ValueFromString(ctx, stringValue)

	if diags.HasError() {
		return nil, fmt.Errorf("unexpected error converting StringValue to StringValuable: %v", diags)
	}

	return stringValuable, nil
}

// ValueType returns the Value type.
func (t apiStateType) ValueType(context.Context) attr.Value {
	return apiStateValue{}
}

var _ basetypes.StringValuable = apiStateValue{}

type apiStateValue struct {
	basetypes.StringValue
}

// Equal returns true if the given type is equivalent.
func (v apiStateValue) Equal(o attr.Value) bool {
	other, ok := o.(apiStateValue)

	if !ok {
		return false
	}

	return v.StringValue.Equal(other.StringValue)
}

func (v apiStateValue) Type(context.Context) attr.Type {
	return apiStateType{}
}

// StringSemanticEquals returns true if the given objects are semantically equal.
func (v apiStateValue) StringSemanticEquals(ctx context.Context, valuable basetypes.StringValuable) (bool, diag.Diagnostics) {
	var diags diag.Diagnostics
	var state, err = deserialize(v.ValueString())
	if err != nil {
		diags.AddError("Unable to deserialize API state", err.Error())
	}
	stringAttribute, _ := valuable.ToStringValue(ctx)
	user, err := deserialize(stringAttribute.ValueString())
	if err != nil {
		diags.AddError("Unable to deserialize API state", err.Error())
	}

	if diff := checkSemanticEquality(*user, *state); diff != "" {
		tflog.Error(ctx, fmt.Sprintf("API state mismatch: \n%s", diff))
		return false, nil
	}
	return true, nil
}

func checkSemanticEquality(before v0.APIAttributes, after v0.APIAttributes) string {
	if before.BasePath != nil && *before.BasePath == "" && after.BasePath == nil {
		after.BasePath = ptr.To("")
	}

	sortConsumeTypes(before)
	sortConsumeTypes(after)
	sortHosts(before)
	sortHosts(after)
	sortTags(before)
	sortTags(after)

	return cmp.Diff(before, after, orderedMapComparer())
}

func sortConsumeTypes(state v0.APIAttributes) {
	if state.Constraints != nil && state.Constraints.RequestBody != nil && state.Constraints.RequestBody.ConsumeType != nil {
		slices.Sort(state.Constraints.RequestBody.ConsumeType)
	}
}

func sortTags(state v0.APIAttributes) {
	if state.Tags != nil {
		slices.Sort(state.Tags)
	}
}
func sortHosts(state v0.APIAttributes) {
	if state.Hostnames != nil {
		slices.Sort(state.Hostnames)
	}
}

func orderedMapComparer() cmp.Option {
	return cmp.FilterValues(func(x, y any) bool {
		_, ok1 := x.(*orderedmap.OrderedMap[string, v0.Resource])
		_, ok2 := y.(*orderedmap.OrderedMap[string, v0.Resource])
		if ok1 && ok2 {
			return true
		}
		_, ok1 = x.(*orderedmap.OrderedMap[string, v0.Property])
		_, ok2 = y.(*orderedmap.OrderedMap[string, v0.Property])
		if ok1 && ok2 {
			return true
		}
		return false
	}, cmp.Comparer(func(x, y any) bool {

		vx := reflect.ValueOf(x)
		vy := reflect.ValueOf(y)

		if vx.Type() != vy.Type() {
			return false
		}
		if vx.IsNil() && vy.IsNil() {
			return true
		}
		if vx.IsNil() || vy.IsNil() {
			return false
		}

		mx, okx := x.(interface{ Len() int })
		my, oky := y.(interface{ Len() int })
		if !okx || !oky {
			return false
		}

		if mx.Len() != my.Len() {
			return false
		}

		switch m := x.(type) {
		case *orderedmap.OrderedMap[string, v0.Resource]:
			m2 := y.(*orderedmap.OrderedMap[string, v0.Resource])
			for pair := m.Oldest(); pair != nil; pair = pair.Next() {
				val2, ok := m2.Get(pair.Key)
				if !ok || cmp.Diff(pair.Value, val2, orderedMapComparer()) != "" {
					return false
				}
			}
			return true
		case *orderedmap.OrderedMap[string, v0.Property]:
			m2 := y.(*orderedmap.OrderedMap[string, v0.Property])
			for pair := m.Oldest(); pair != nil; pair = pair.Next() {
				val2, ok := m2.Get(pair.Key)
				if !ok || cmp.Diff(pair.Value, val2, orderedMapComparer()) != "" {
					return false
				}
			}
			return true
		}
		return false
	}))
}
