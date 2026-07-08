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

// ruleConditionExceptionStateType - custom type for rule-level condition exceptions
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
var _ basetypes.StringTypable = ruleConditionExceptionStateType{}

type ruleConditionExceptionStateType struct {
	basetypes.StringType
}

func (t ruleConditionExceptionStateType) Equal(o attr.Type) bool {
	other, ok := o.(ruleConditionExceptionStateType)
	if !ok {
		return false
	}
	return t.StringType.Equal(other.StringType)
}

func (t ruleConditionExceptionStateType) String() string {
	return "ruleConditionExceptionStateType"
}

func (t ruleConditionExceptionStateType) ValueFromString(_ context.Context, in basetypes.StringValue) (basetypes.StringValuable, diag.Diagnostics) {
	value := ruleConditionExceptionStateValue{
		StringValue: in,
	}
	return value, nil
}

func (t ruleConditionExceptionStateType) ValueFromTerraform(ctx context.Context, in tftypes.Value) (attr.Value, error) {
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

func (t ruleConditionExceptionStateType) ValueType(context.Context) attr.Value {
	return ruleConditionExceptionStateValue{}
}

type ruleConditionExceptionStateValue struct {
	basetypes.StringValue
}

func (v ruleConditionExceptionStateValue) Equal(o attr.Value) bool {
	other, ok := o.(ruleConditionExceptionStateValue)
	if !ok {
		return false
	}
	return v.StringValue.Equal(other.StringValue)
}

func (v ruleConditionExceptionStateValue) Type(context.Context) attr.Type {
	return ruleConditionExceptionStateType{}
}

// StringSemanticEquals compares rule condition exceptions semantically by unmarshaling to RuleConditionException.
func (v ruleConditionExceptionStateValue) StringSemanticEquals(ctx context.Context, valuable basetypes.StringValuable) (bool, diag.Diagnostics) {
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

	// Unmarshal both to appsec.RuleConditionException for consistent comparison
	var stateCondition, userCondition appsec.RuleConditionException

	if err := json.Unmarshal([]byte(stateValueString), &stateCondition); err != nil {
		diags.AddError("Unable to unmarshal state rule condition exception", err.Error())
		return false, diags
	}

	if err := json.Unmarshal([]byte(userValueString), &userCondition); err != nil {
		diags.AddError("Unable to unmarshal user rule condition exception", err.Error())
		return false, diags
	}

	// Compare using deep equal on the structs
	if diff := deep.Equal(stateCondition, userCondition); diff != nil {
		return false, nil
	}

	return true, nil
}
