package apidefinitions

import (
	"context"
	"fmt"
	"strings"

	v0 "github.com/akamai/AkamaiOPEN-edgegrid-golang/v10/pkg/apidefinitions/v0"
	"github.com/akamai/terraform-provider-akamai/v7/pkg/common/ptr"
	"github.com/go-test/deep"
	"github.com/hashicorp/terraform-plugin-framework/attr"
	"github.com/hashicorp/terraform-plugin-framework/diag"
	"github.com/hashicorp/terraform-plugin-framework/types/basetypes"
	"github.com/hashicorp/terraform-plugin-go/tftypes"
	"github.com/hashicorp/terraform-plugin-log/tflog"
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
	var user, err1 = deserialize(stringAttribute.ValueString())
	if err1 != nil {
		diags.AddError("Unable to deserialize API state", err.Error())
	}

	if diff := checkSemanticEquality(*user, *state); diff != nil {
		tflog.Error(ctx, strings.ReplaceAll(strings.Join(diff, "\n "), ",", "\n"))
		return false, nil
	}
	return true, nil
}

<<<<<<< HEAD
func checkSemanticEquality(before v0.APIAttributes, after v0.APIAttributes) []string {
=======
func checkSemanticEquality(before v0.RegisterAPIRequest, after v0.RegisterAPIRequest) []string {
	if before.EnableAPIGateway == nil && after.EnableAPIGateway != nil && !*after.EnableAPIGateway {
		after.EnableAPIGateway = nil
	}
	if before.GraphQL == nil && after.GraphQL != nil && !*after.GraphQL {
		after.GraphQL = nil
	}
	if before.MatchPathSegmentParameter == nil && after.MatchPathSegmentParameter != nil && !*after.MatchPathSegmentParameter {
		after.MatchPathSegmentParameter = nil
	}
	if before.MatchCaseSensitive == nil && after.MatchCaseSensitive != nil && !*after.MatchCaseSensitive {
		after.MatchCaseSensitive = nil
	}

>>>>>>> 298d84c3 (SECKSD-27443 API Definitions Sub Provider)
	if before.BasePath != nil && *before.BasePath == "" && after.BasePath == nil {
		after.BasePath = ptr.To("")
	}

<<<<<<< HEAD
	sortConsumeTypes(before)
	sortConsumeTypes(after)
=======
	sortCollections(before)
	sortCollections(after)
>>>>>>> 298d84c3 (SECKSD-27443 API Definitions Sub Provider)

	return deep.Equal(before, after)
}

<<<<<<< HEAD
func sortConsumeTypes(state v0.APIAttributes) {
=======
func sortCollections(state v0.RegisterAPIRequest) {
	sortConsumeTypes(state)
	sortBypassOn(state)
}

func sortConsumeTypes(state v0.RegisterAPIRequest) {
>>>>>>> 298d84c3 (SECKSD-27443 API Definitions Sub Provider)
	if state.Constraints != nil && state.Constraints.RequestBody != nil && state.Constraints.RequestBody.ConsumeType != nil {
		slices.Sort(state.Constraints.RequestBody.ConsumeType)
	}
}
<<<<<<< HEAD
=======

func sortBypassOn(state v0.RegisterAPIRequest) {
	if state.Constraints != nil && state.Constraints.BypassOn != nil {
		if state.Constraints.BypassOn.UndefinedMethods != nil {
			slices.Sort(state.Constraints.BypassOn.UndefinedMethods)
		}

		if state.Constraints.BypassOn.UndefinedParameters != nil {
			slices.Sort(state.Constraints.BypassOn.UndefinedParameters)
		}
	}

	if state.Resources != nil {
		for pair := state.Resources.Oldest(); pair != nil; pair = pair.Next() {
			resource := pair.Value
			sortMethodBypassOn(resource.Get)
			sortMethodBypassOn(resource.Post)
			sortMethodBypassOn(resource.Put)
			sortMethodBypassOn(resource.Delete)
			sortMethodBypassOn(resource.Options)
			sortMethodBypassOn(resource.Head)
			sortMethodBypassOn(resource.Patch)
		}

	}
}

func sortMethodBypassOn(method *v0.Method) {
	if method == nil || method.BypassOn == nil || method.BypassOn.UndefinedParameters == nil {
		return
	}
	slices.Sort(method.BypassOn.UndefinedParameters)
}
>>>>>>> 298d84c3 (SECKSD-27443 API Definitions Sub Provider)
