package customtypes

import (
	"context"
	"fmt"

	"github.com/akamai/terraform-provider-akamai/v9/internal/text"
	"github.com/hashicorp/terraform-plugin-framework/attr"
	"github.com/hashicorp/terraform-plugin-framework/diag"
	"github.com/hashicorp/terraform-plugin-framework/types"
	"github.com/hashicorp/terraform-plugin-framework/types/basetypes"
	"github.com/hashicorp/terraform-plugin-go/tftypes"
)

var _ basetypes.StringTypable = IgnoreTrailingWhitespaceType{}
var _ basetypes.StringValuable = IgnoreTrailingWhitespaceValue{}
var _ basetypes.StringValuableWithSemanticEquals = IgnoreTrailingWhitespaceValue{}

// IgnoreTrailingWhitespaceType is a custom string type that ignores trailing whitespace.
type IgnoreTrailingWhitespaceType struct {
	basetypes.StringType
}

// Equal checks if the current type is equal to another type.
func (t IgnoreTrailingWhitespaceType) Equal(o attr.Type) bool {
	other, ok := o.(IgnoreTrailingWhitespaceType)

	if !ok {
		return false
	}

	return t.StringType.Equal(other.StringType)
}

func (t IgnoreTrailingWhitespaceType) String() string {
	return "IgnoreTrailingWhitespaceType"
}

// ValueFromString converts a StringValue to a IgnoreTrailingWhitespaceValue.
func (t IgnoreTrailingWhitespaceType) ValueFromString(_ context.Context, in basetypes.StringValue) (basetypes.StringValuable, diag.Diagnostics) {
	value := IgnoreTrailingWhitespaceValue{
		StringValue: in,
	}

	return value, nil
}

// ValueFromTerraform converts a Terraform value to a IgnoreTrailingWhitespaceValue.
func (t IgnoreTrailingWhitespaceType) ValueFromTerraform(ctx context.Context, in tftypes.Value) (attr.Value, error) {
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
func (t IgnoreTrailingWhitespaceType) ValueType(_ context.Context) attr.Value {
	return IgnoreTrailingWhitespaceValue{}
}

// IgnoreTrailingWhitespaceValue is a custom string value that ignores trailing whitespace.
type IgnoreTrailingWhitespaceValue struct {
	basetypes.StringValue
}

// Equal checks if the current value is equal to another value.
func (v IgnoreTrailingWhitespaceValue) Equal(o attr.Value) bool {
	other, ok := o.(IgnoreTrailingWhitespaceValue)

	if !ok {
		return false
	}

	return v.StringValue.Equal(other.StringValue)
}

// StringSemanticEquals checks if the current value is semantically equal to another value,
func (v IgnoreTrailingWhitespaceValue) StringSemanticEquals(_ context.Context, newValuable basetypes.StringValuable) (bool, diag.Diagnostics) {
	var diags diag.Diagnostics

	// The framework should always pass the correct value type, but always check
	newValue, ok := newValuable.(IgnoreTrailingWhitespaceValue)

	if !ok {
		diags.AddError(
			"Semantic Equality Check Error",
			"An unexpected value type was received while performing semantic equality checks. "+
				"Please report this to the provider developers.\n\n"+
				"Expected Value Type: "+fmt.Sprintf("%T", v)+"\n"+
				"Got Value Type: "+fmt.Sprintf("%T", newValuable),
		)

		return false, diags
	}

	oldValueTrimmed := text.TrimRightWhitespace(v.ValueString())
	newValueTrimmed := text.TrimRightWhitespace(newValue.ValueString())
	return oldValueTrimmed == newValueTrimmed, diags
}

// Type returns the type of the value.
func (v IgnoreTrailingWhitespaceValue) Type(_ context.Context) attr.Type {
	return IgnoreTrailingWhitespaceType{}
}

// NewIgnoreTrailingWhitespaceValue creates a new IgnoreTrailingWhitespaceValue from a string.
func NewIgnoreTrailingWhitespaceValue(value string) IgnoreTrailingWhitespaceValue {
	return IgnoreTrailingWhitespaceValue{
		StringValue: types.StringValue(value),
	}
}
