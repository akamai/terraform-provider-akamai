package customtypes

import (
	"context"
	"fmt"
	"strings"

	"github.com/hashicorp/terraform-plugin-framework/attr"
	"github.com/hashicorp/terraform-plugin-framework/diag"
	"github.com/hashicorp/terraform-plugin-framework/types"
	"github.com/hashicorp/terraform-plugin-framework/types/basetypes"
	"github.com/hashicorp/terraform-plugin-go/tftypes"
)

var _ basetypes.StringTypable = IgnorePrefixType{}
var _ basetypes.StringValuable = IgnorePrefixValue{}
var _ basetypes.StringValuableWithSemanticEquals = IgnorePrefixValue{}

// IgnorePrefixType is a custom string type that considers two strings equal when
// they differ only by an optional leading prefix (e.g. "cpc_111" == "111").
type IgnorePrefixType struct {
	basetypes.StringType
	// Prefix is the optional prefix to strip before comparing (e.g. "cpc_").
	Prefix string
}

// Equal checks if the current type is equal to another type.
func (t IgnorePrefixType) Equal(o attr.Type) bool {
	other, ok := o.(IgnorePrefixType)
	if !ok {
		return false
	}
	return t.Prefix == other.Prefix && t.StringType.Equal(other.StringType)
}

func (t IgnorePrefixType) String() string {
	return fmt.Sprintf("IgnorePrefixType(%q)", t.Prefix)
}

// ValueFromString converts a StringValue to an IgnorePrefixValue.
func (t IgnorePrefixType) ValueFromString(_ context.Context, in basetypes.StringValue) (basetypes.StringValuable, diag.Diagnostics) {
	return IgnorePrefixValue{
		StringValue: in,
		prefix:      t.Prefix,
	}, nil
}

// ValueFromTerraform converts a Terraform value to an IgnorePrefixValue.
func (t IgnorePrefixType) ValueFromTerraform(ctx context.Context, in tftypes.Value) (attr.Value, error) {
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
func (t IgnorePrefixType) ValueType(_ context.Context) attr.Value {
	return IgnorePrefixValue{prefix: t.Prefix}
}

// IgnorePrefixValue is a custom string value that considers two strings equal
// when they differ only by an optional leading prefix.
type IgnorePrefixValue struct {
	basetypes.StringValue
	prefix string
}

// Equal checks if the current value is equal to another value.
func (v IgnorePrefixValue) Equal(o attr.Value) bool {
	other, ok := o.(IgnorePrefixValue)
	if !ok {
		return false
	}
	return v.StringValue.Equal(other.StringValue)
}

// StringSemanticEquals checks if two values are equal after stripping the optional prefix.
func (v IgnorePrefixValue) StringSemanticEquals(_ context.Context, newValuable basetypes.StringValuable) (bool, diag.Diagnostics) {
	var diags diag.Diagnostics

	newValue, ok := newValuable.(IgnorePrefixValue)
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

	return strings.TrimPrefix(v.ValueString(), v.prefix) == strings.TrimPrefix(newValue.ValueString(), v.prefix), diags
}

// Type returns the type of the value.
func (v IgnorePrefixValue) Type(_ context.Context) attr.Type {
	return IgnorePrefixType{Prefix: v.prefix}
}

// NewIgnorePrefixValue creates a new IgnorePrefixValue with the given prefix and value.
func NewIgnorePrefixValue(prefix, value string) IgnorePrefixValue {
	return IgnorePrefixValue{
		StringValue: types.StringValue(value),
		prefix:      prefix,
	}
}
