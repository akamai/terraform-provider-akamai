package apidefinitions

import (
	"bytes"
	"context"
	"encoding/json"
	"fmt"
	"sort"
	"strings"

	v0 "github.com/akamai/AkamaiOPEN-edgegrid-golang/v12/pkg/apidefinitions/v0"
	"github.com/hashicorp/terraform-plugin-framework/attr"
	"github.com/hashicorp/terraform-plugin-framework/diag"
	"github.com/hashicorp/terraform-plugin-framework/schema/validator"
	"github.com/hashicorp/terraform-plugin-framework/types/basetypes"
	"github.com/hashicorp/terraform-plugin-go/tftypes"
)

var _ basetypes.StringTypable = operationsStateType{}

var _ validator.String = operationsStateValidator{}

type operationsStateValidator struct{}

func (o operationsStateValidator) Description(_ context.Context) string {
	return "Invalid JSON-formatted provided in the field resource_operations"
}

func (o operationsStateValidator) MarkdownDescription(ctx context.Context) string {
	return o.Description(ctx)
}

func (o operationsStateValidator) ValidateString(_ context.Context, request validator.StringRequest, response *validator.StringResponse) {
	if request.ConfigValue.IsNull() || request.ConfigValue.IsUnknown() {
		return
	}

	value := request.ConfigValue.ValueString()
	endpoint := v0.ResourceOperationResponse{}
	decoder := json.NewDecoder(strings.NewReader(value))
	decoder.DisallowUnknownFields()
	err := decoder.Decode(&endpoint)
	if err != nil {
		var name = err.Error()
		response.Diagnostics.AddError("Invalid JSON provided", name)
	}
}

type operationsStateType struct {
	basetypes.StringType
}

// Equal returns true if the given type is equivalent.
func (t operationsStateType) Equal(o attr.Type) bool {
	other, ok := o.(operationsStateType)

	if !ok {
		return false
	}

	return t.StringType.Equal(other.StringType)
}

func (t operationsStateType) String() string {
	return "operationsStateType"
}

// ValueFromString should convert the String to a StringValuable type.
func (t operationsStateType) ValueFromString(_ context.Context, in basetypes.StringValue) (basetypes.StringValuable, diag.Diagnostics) {
	value := operationsStateValue{
		StringValue: in,
	}

	return value, nil
}

// ValueFromTerraform returns a Value given a tftypes.Value.
func (t operationsStateType) ValueFromTerraform(ctx context.Context, in tftypes.Value) (attr.Value, error) {
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
func (t operationsStateType) ValueType(context.Context) attr.Value {
	return operationsStateValue{}
}

var _ basetypes.StringValuable = operationsStateValue{}

type operationsStateValue struct {
	basetypes.StringValue
}

// Equal returns true if the given type is equivalent.
func (v operationsStateValue) Equal(o attr.Value) bool {
	other, ok := o.(operationsStateValue)

	if !ok {
		return false
	}

	return v.StringValue.Equal(other.StringValue)
}

func (v operationsStateValue) Type(context.Context) attr.Type {
	return operationsStateType{}
}

// StringSemanticEquals returns true if the given objects are semantically equal.
func (v operationsStateValue) StringSemanticEquals(ctx context.Context, valuable basetypes.StringValuable) (bool, diag.Diagnostics) {
	var diags diag.Diagnostics
	var after, err = normalizeJSON(v.ValueString())
	if err != nil {
		diags.AddError("Semantic check for operations state failed", err.Error())
		return false, nil
	}
	stringAttribute, _ := valuable.ToStringValue(ctx)
	before, err := normalizeJSON(stringAttribute.ValueString())
	if err != nil && before != after {
		diags.AddError("Semantic check for operations state failed", err.Error())
		return false, nil
	}

	return true, nil
}

// OperationsStateValidator returns a validator which ensures that JSON-formatted information
func OperationsStateValidator() validator.String {
	return operationsStateValidator{}
}

func normalizeJSON(input string) (string, error) {
	var decoded interface{}
	if err := json.Unmarshal([]byte(input), &decoded); err != nil {
		return "", fmt.Errorf("failed to unmarshal JSON: %v", err)
	}

	buffer := &bytes.Buffer{}
	encoder := json.NewEncoder(buffer)
	encoder.SetEscapeHTML(false)
	encoder.SetIndent("", "")

	var sortedJSON interface{}
	switch v := decoded.(type) {
	case map[string]interface{}:
		sortedJSON = sortMapKeys(v)
	default:
		sortedJSON = decoded
	}

	if err := encoder.Encode(sortedJSON); err != nil {
		return "", fmt.Errorf("failed to encode JSON: %v", err)
	}

	return strings.TrimSpace(buffer.String()), nil
}

func sortMapKeys(m map[string]interface{}) map[string]interface{} {
	sorted := make(map[string]interface{})
	keys := make([]string, 0, len(m))

	for k := range m {
		keys = append(keys, k)
	}

	sort.Strings(keys)

	for _, k := range keys {
		v := m[k]
		if nested, ok := v.(map[string]interface{}); ok {
			sorted[k] = sortMapKeys(nested)
		} else if slice, ok := v.([]interface{}); ok {
			for i, item := range slice {
				if nestedMap, ok := item.(map[string]interface{}); ok {
					slice[i] = sortMapKeys(nestedMap)
				}
			}
			sorted[k] = slice
		} else {
			sorted[k] = v
		}
	}

	return sorted
}
