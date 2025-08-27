package appsec

import (
	"context"
	"fmt"
	"sort"
	"strings"

	"github.com/akamai/AkamaiOPEN-edgegrid-golang/v12/pkg/appsec"
	"github.com/go-test/deep"
	"github.com/hashicorp/terraform-plugin-framework/attr"
	"github.com/hashicorp/terraform-plugin-framework/diag"
	"github.com/hashicorp/terraform-plugin-framework/types/basetypes"
	"github.com/hashicorp/terraform-plugin-go/tftypes"
	"github.com/hashicorp/terraform-plugin-log/tflog"
)

var _ basetypes.StringTypable = rapidRuleDefinitionsStateType{}

type rapidRuleDefinitionsStateType struct {
	basetypes.StringType
}

// Equal returns true if the given type is equivalent.
func (t rapidRuleDefinitionsStateType) Equal(o attr.Type) bool {
	other, ok := o.(rapidRuleDefinitionsStateType)

	if !ok {
		return false
	}

	return t.StringType.Equal(other.StringType)
}

func (t rapidRuleDefinitionsStateType) String() string {
	return "rapidRuleDefinitionsStateType"
}

// ValueFromString should convert the String to a StringValuable type.
func (t rapidRuleDefinitionsStateType) ValueFromString(_ context.Context, in basetypes.StringValue) (basetypes.StringValuable, diag.Diagnostics) {
	value := rapidRuleDefinitionsStateValue{
		StringValue: in,
	}

	return value, nil
}

// ValueFromTerraform returns a Value given a tftypes.Value.
func (t rapidRuleDefinitionsStateType) ValueFromTerraform(ctx context.Context, in tftypes.Value) (attr.Value, error) {
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
func (t rapidRuleDefinitionsStateType) ValueType(context.Context) attr.Value {
	return rapidRuleDefinitionsStateValue{}
}

type rapidRuleDefinitionsStateValue struct {
	basetypes.StringValue
}

// Equal returns true if the given type is equivalent.
func (v rapidRuleDefinitionsStateValue) Equal(o attr.Value) bool {
	other, ok := o.(rapidRuleDefinitionsStateValue)

	if !ok {
		return false
	}

	return v.StringValue.Equal(other.StringValue)
}

func (v rapidRuleDefinitionsStateValue) Type(context.Context) attr.Type {
	return rapidRuleDefinitionsStateType{}
}

// StringSemanticEquals returns true if the given objects are semantically equal.
func (v rapidRuleDefinitionsStateValue) StringSemanticEquals(ctx context.Context, valuable basetypes.StringValuable) (bool, diag.Diagnostics) {
	var diags diag.Diagnostics
	stringAttribute, diagnostics := valuable.ToStringValue(ctx)
	if diagnostics.HasError() {
		diags.Append(diagnostics...)
		return false, diags
	}

	stateValueString := v.ValueString()
	userValueString := stringAttribute.ValueString()
	if stateValueString == "" || userValueString == "" {
		return stateValueString == userValueString, nil
	}

	state, err := deserializeRuleDefinitions(v.ValueString())
	if err != nil {
		diags.AddError("Unable to deserialize API state", err.Error())
		return false, diags
	}

	user, err := deserializeRuleDefinitions(stringAttribute.ValueString())
	if err != nil {
		diags.AddError("Unable to deserialize API state", err.Error())
		return false, diags
	}

	if diff := checkSemanticEquality(*user, *state); diff != nil {
		tflog.Error(ctx, strings.ReplaceAll(strings.Join(diff, "\n "), ",", "\n"))
		return false, nil
	}
	return true, nil
}

func checkSemanticEquality(before []appsec.RuleDefinition, after []appsec.RuleDefinition) []string {
	sortCollections(before)
	sortCollections(after)

	return deep.Equal(before, after)
}

func sortCollections(state []appsec.RuleDefinition) {
	sort.Slice(state, func(i, j int) bool {
		return *state[i].ID > *state[j].ID
	})
}
