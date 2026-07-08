package appsec

import (
	"context"
	"encoding/json"
	"fmt"

	"github.com/akamai/AkamaiOPEN-edgegrid-golang/v13/pkg/appsec"
	"github.com/go-test/deep"
	"github.com/hashicorp/terraform-plugin-framework/attr"
	"github.com/hashicorp/terraform-plugin-framework/diag"
	"github.com/hashicorp/terraform-plugin-framework/types/basetypes"
	"github.com/hashicorp/terraform-plugin-go/tftypes"
)

// attackGroupConditionExceptionStateType - custom type for attack group-level condition exceptions
// This custom type encapsulates the JSON marshal/unmarshal logic between the
// API representation and the Terraform state representation.
//
// The AppSec API returns this field as structured JSON, while Terraform stores
// it as a string in the state. This type centralizes the conversion logic so it
// does not need to be handled repeatedly across the resource implementation.
//
// It also helps preserve values such as `false` or `0` that may otherwise be
// dropped during serialization due to Go's JSON `omitempty` behavior, which
// can lead to configuration/state mismatches after resource creation.
var _ basetypes.StringTypable = attackGroupConditionExceptionStateType{}

type attackGroupConditionExceptionStateType struct {
	basetypes.StringType
}

func (t attackGroupConditionExceptionStateType) Equal(o attr.Type) bool {
	other, ok := o.(attackGroupConditionExceptionStateType)
	if !ok {
		return false
	}
	return t.StringType.Equal(other.StringType)
}

func (t attackGroupConditionExceptionStateType) String() string {
	return "attackGroupConditionExceptionStateType"
}

func (t attackGroupConditionExceptionStateType) ValueFromString(_ context.Context, in basetypes.StringValue) (basetypes.StringValuable, diag.Diagnostics) {
	value := attackGroupConditionExceptionStateValue{
		StringValue: in,
	}
	return value, nil
}

func (t attackGroupConditionExceptionStateType) ValueFromTerraform(ctx context.Context, in tftypes.Value) (attr.Value, error) {
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

func (t attackGroupConditionExceptionStateType) ValueType(context.Context) attr.Value {
	return attackGroupConditionExceptionStateValue{}
}

type attackGroupConditionExceptionStateValue struct {
	basetypes.StringValue
}

func (v attackGroupConditionExceptionStateValue) Equal(o attr.Value) bool {
	other, ok := o.(attackGroupConditionExceptionStateValue)
	if !ok {
		return false
	}
	return v.StringValue.Equal(other.StringValue)
}

func (v attackGroupConditionExceptionStateValue) Type(context.Context) attr.Type {
	return attackGroupConditionExceptionStateType{}
}

// StringSemanticEquals compares attack group condition exceptions semantically by unmarshaling to AttackGroupConditionException.
func (v attackGroupConditionExceptionStateValue) StringSemanticEquals(ctx context.Context, valuable basetypes.StringValuable) (bool, diag.Diagnostics) {
	var diags diag.Diagnostics
	stringAttribute, diagnostics := valuable.ToStringValue(ctx)
	if diagnostics.HasError() {
		diags.Append(diagnostics...)
		return false, diags
	}

	stateValueString := v.ValueString()
	userValueString := stringAttribute.ValueString()

	// If both are empty, they're equal
	if stateValueString == "" && userValueString == "" {
		return true, nil
	}
	if stateValueString == "" || userValueString == "" {
		return false, nil
	}

	// Unmarshal both to appsec.AttackGroupConditionException for consistent comparison
	var stateCondition, userCondition appsec.AttackGroupConditionException

	if err := json.Unmarshal([]byte(stateValueString), &stateCondition); err != nil {
		diags.AddError("Unable to unmarshal state attack group condition exception", err.Error())
		return false, diags
	}

	if err := json.Unmarshal([]byte(userValueString), &userCondition); err != nil {
		diags.AddError("Unable to unmarshal user attack group condition exception", err.Error())
		return false, diags
	}

	// Compare using deep equal on the structs
	if diff := deep.Equal(stateCondition, userCondition); diff != nil {
		return false, nil
	}

	return true, nil
}
