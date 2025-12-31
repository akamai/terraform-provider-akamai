package modifiers

import (
	"context"
	"encoding"
	"encoding/json"
	"fmt"
	"sort"
	"strconv"
	"strings"

	"github.com/hashicorp/terraform-plugin-framework/resource/schema/planmodifier"
)

// PreventJSONReorder returns a plan modifier that suppresses diffs
// when the attribute is a JSON array whose elements are the same as in state,
// but appear in a different order.
func PreventJSONReorder() planmodifier.String {
	return newStringSuppressJSONListReorderModifier(nil)
}

// stringSuppressJSONListReorderModifier implements planmodifier.String.
type stringSuppressJSONListReorderModifier struct {
	equal func(a, b []json.RawMessage) bool
}

func newStringSuppressJSONListReorderModifier(
	equal func(a, b []json.RawMessage) bool,
) stringSuppressJSONListReorderModifier {
	if equal == nil {
		equal = rawMessageSlicesEqualIgnoringOrder
	}

	return stringSuppressJSONListReorderModifier{
		equal: equal,
	}
}

// Description returns a human-readable description of the plan modifier.
func (m stringSuppressJSONListReorderModifier) Description(_ context.Context) string {
	return "Ignore diffs when the JSON value is a list whose elements are the same as in the current state, differing only by order."
}

// MarkdownDescription returns a markdown description of the plan modifier.
func (m stringSuppressJSONListReorderModifier) MarkdownDescription(ctx context.Context) string {
	return m.Description(ctx)
}

// PlanModifyString if the planned JSON list and the
// state JSON list contain the same elements (ignoring order), set the plan
// equal to the state so Terraform records no diff.
func (m stringSuppressJSONListReorderModifier) PlanModifyString(
	_ context.Context,
	req planmodifier.StringRequest,
	resp *planmodifier.StringResponse,
) {

	if req.StateValue.IsNull() || req.StateValue.IsUnknown() {
		return
	}

	if req.PlanValue.IsNull() || req.PlanValue.IsUnknown() {
		return
	}

	stateStr := strings.TrimSpace(req.StateValue.ValueString())
	planStr := strings.TrimSpace(req.PlanValue.ValueString())

	if stateStr == planStr {
		return
	}

	// If either side is empty, "null", or not starting with '[', bail out -> normal diff.
	if stateStr == "" || planStr == "" {
		return
	}
	if stateStr[0] != '[' || planStr[0] != '[' {
		return
	}

	var stateList, planList []json.RawMessage

	if err := json.Unmarshal([]byte(stateStr), &stateList); err != nil {
		return
	}
	if err := json.Unmarshal([]byte(planStr), &planList); err != nil {
		return
	}

	// Different lengths or empty list → real change (add/remove).
	if len(stateList) != len(planList) || len(stateList) == 0 {
		return
	}

	// If the two lists contain the same elements (ignoring order), then
	// the only difference is ordering → copy state into plan.
	if !m.equal(stateList, planList) {
		return
	}

	resp.PlanValue = req.StateValue
}

// rawMessageSlicesEqualIgnoringOrder returns true if both slices contain
// the same JSON elements, regardless of order (multiplicity is respected).
func rawMessageSlicesEqualIgnoringOrder(a, b []json.RawMessage) bool {
	if len(a) != len(b) {
		return false
	}

	counts := make(map[string]int, len(a))

	for _, av := range a {
		key, err := canonicalJSON(av)
		if err != nil {
			return false
		}

		counts[key]++
	}

	for _, bv := range b {
		key, err := canonicalJSON(bv)
		if err != nil {
			return false
		}

		remaining, ok := counts[key]
		if !ok {
			return false
		}

		if remaining == 1 {
			delete(counts, key)
			continue
		}

		counts[key] = remaining - 1
	}

	return len(counts) == 0
}

func canonicalJSON(raw json.RawMessage) (string, error) {
	var decoded interface{}
	decoder := json.NewDecoder(strings.NewReader(string(raw)))
	decoder.UseNumber()

	if err := decoder.Decode(&decoded); err != nil {
		return "", err
	}

	var b strings.Builder
	if err := writeCanonicalJSON(&b, decoded); err != nil {
		return "", err
	}

	return b.String(), nil
}

// writeCanonicalJSON renders a deterministic representation of JSON values.
func writeCanonicalJSON(b *strings.Builder, v interface{}) error {
	switch val := v.(type) {
	case nil:
		b.WriteString("null")
	case bool:
		b.WriteString(strconv.FormatBool(val))
	case string:
		b.WriteString(strconv.Quote(val))
	case json.Number:
		b.WriteString(val.String())
	case float64:
		b.WriteString(strconv.FormatFloat(val, 'g', -1, 64))
	case []interface{}:
		b.WriteByte('[')
		for i, elem := range val {
			if i > 0 {
				b.WriteByte(',')
			}
			if err := writeCanonicalJSON(b, elem); err != nil {
				return err
			}
		}
		b.WriteByte(']')
	case map[string]interface{}:
		keys := make([]string, 0, len(val))
		for k := range val {
			keys = append(keys, k)
		}
		sort.Strings(keys)

		b.WriteByte('{')
		for i, k := range keys {
			if i > 0 {
				b.WriteByte(',')
			}
			b.WriteString(strconv.Quote(k))
			b.WriteByte(':')
			if err := writeCanonicalJSON(b, val[k]); err != nil {
				return err
			}
		}
		b.WriteByte('}')
	default:
		if marshaler, ok := val.(encoding.TextMarshaler); ok {
			text, err := marshaler.MarshalText()
			if err != nil {
				return err
			}
			b.WriteString(strconv.Quote(string(text)))
			return nil
		}
		return fmt.Errorf("unsupported JSON type %T", v)
	}

	return nil
}
