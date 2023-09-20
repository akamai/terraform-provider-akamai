package ruleformats

import (
	"errors"
	"fmt"

	"github.com/akamai/AkamaiOPEN-edgegrid-golang/v7/pkg/papi"
	"github.com/akamai/terraform-provider-akamai/v5/pkg/common/tf"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/schema"
)

type attributeGetter interface {
	GetOk(string) (any, bool)
}

// RulesSchemaReader knows how to retrieve data from schema.ResourceData.
type RulesSchemaReader struct {
	data          attributeGetter
	ruleFormatKey string
}

// RuleItem is a struct for holding information about a single behavior or criterion.
type RuleItem struct {
	Name string
	Item map[string]any
}

// RuleItems contains slice of RuleItem.
type RuleItems []RuleItem

// GetUsedRuleFormat finds RuleVersion that is used in schema.ResourceData.
func GetUsedRuleFormat(d *schema.ResourceData) RuleVersion {
	for _, r := range RulesFormats() {
		if _, ok := d.GetOk(r.SchemaKey()); ok {
			return r
		}
	}
	// should not happen, unable to continue
	panic("rule format not found")
}

// NewRulesSchemaReader creates a new RulesSchemaReader for the given schema.ResourceData.
func NewRulesSchemaReader(d *schema.ResourceData) *RulesSchemaReader {
	rfVersion := GetUsedRuleFormat(d)

	reader := RulesSchemaReader{
		data:          tf.NewRawConfig(d),
		ruleFormatKey: rfVersion.SchemaKey(),
	}
	return &reader
}

// GetRuleFormat returns rule format as a string.
func (r *RulesSchemaReader) GetRuleFormat() string {
	return r.ruleFormatKey
}

// GetBehaviorsList reads and returns a slice of RuleItem, which are behaviors.
func (r *RulesSchemaReader) GetBehaviorsList() ([]RuleItem, error) {
	return r.getRuleItems(r.behaviorsKey())
}

// GetCriteriaList reads and returns a slice of RuleItem, which are rule criteria.
func (r *RulesSchemaReader) GetCriteriaList() ([]RuleItem, error) {
	return r.getRuleItems(r.criteriaKey())
}

// GetVariablesList reads and returns rule variables as slice of map[string]any.
func (r *RulesSchemaReader) GetVariablesList() ([]map[string]any, error) {
	rawVal, ok := r.data.GetOk(r.variablesKey())
	if !ok || rawVal == nil {
		return nil, &NotFoundError{r.variablesKey()}
	}

	listVal, ok := rawVal.([]any)
	if !ok {
		return nil, &TypeAssertionError{"[]any", typeof(rawVal), r.variablesKey()}
	}

	variables := make([]map[string]any, 0, len(listVal))
	for _, val := range listVal {
		variable, ok := val.(map[string]any)
		if !ok {
			return nil, &TypeAssertionError{"map[string]any", typeof(val), r.variablesKey()}
		}
		variables = append(variables, variable)
	}

	return variables, nil
}

// GetChildrenList reads and returns slice of children.
func (r *RulesSchemaReader) GetChildrenList() ([]string, error) {
	// rawVal, ok := dataOK.(r.childrenKey())
	rawVal, ok := r.data.GetOk(r.childrenKey())
	if !ok || rawVal == nil {
		return nil, &NotFoundError{r.childrenKey()}
	}

	listVal, ok := rawVal.([]any)
	if !ok {
		return nil, &TypeAssertionError{"[]any", typeof(rawVal), r.childrenKey()}
	}

	children := make([]string, 0, len(listVal))
	for _, val := range listVal {
		child, ok := val.(string)
		if !ok {
			return nil, &TypeAssertionError{"string", typeof(val), r.childrenKey()}
		}
		children = append(children, child)
	}

	return children, nil
}

func (r *RulesSchemaReader) getString(key string) (string, error) {
	rawVal, ok := r.data.GetOk(key)
	if !ok || rawVal == nil {
		return "", &NotFoundError{key}
	}
	val, ok := rawVal.(string)
	if !ok {
		return "", &TypeAssertionError{"string", typeof(rawVal), key}
	}
	return val, nil
}

func (r *RulesSchemaReader) getBool(key string) (bool, error) {
	rawVal, ok := r.data.GetOk(key)
	if !ok || rawVal == nil {
		return false, &NotFoundError{key}
	}
	val, ok := rawVal.(bool)
	if !ok {
		return false, &TypeAssertionError{"bool", typeof(rawVal), key}
	}
	return val, nil
}

func (r *RulesSchemaReader) getCustomOverride(key string) (*papi.RuleCustomOverride, error) {
	rawVal, ok := r.data.GetOk(key)
	if !ok || rawVal == nil {
		return nil, &NotFoundError{key}
	}
	val, ok := rawVal.([]any)
	if !ok {
		return nil, &TypeAssertionError{"[]any", typeof(rawVal), key}
	}
	if len(val) == 0 {
		return nil, &NotFoundError{key}
	}
	override, ok := val[0].(map[string]any)
	if !ok {
		return nil, &TypeAssertionError{"map[string]any", typeof(val[0]), key}
	}

	customOverride := &papi.RuleCustomOverride{
		Name:       override["name"].(string),
		OverrideID: override["override_id"].(string),
	}

	return customOverride, nil
}

func (r *RulesSchemaReader) getMapOfSlice(key string) (map[string][]any, error) {
	rawVal, ok := r.data.GetOk(key)
	if !ok || rawVal == nil {
		return nil, &NotFoundError{key}
	}

	mapVal, ok := rawVal.(map[string]any)
	if !ok {
		return nil, &TypeAssertionError{"map[string]any", typeof(rawVal), key}
	}

	mapUnpacked := make(map[string][]any, len(mapVal))
	for name, val := range mapVal {
		behaviorsSlice, ok := val.([]any)
		if !ok {
			return nil, &TypeAssertionError{"any", typeof(val), key}
		}
		mapUnpacked[name] = behaviorsSlice
	}

	return mapUnpacked, nil
}

func (r *RulesSchemaReader) getRuleItems(key string) ([]RuleItem, error) {
	rawVal, ok := r.data.GetOk(key)
	if !ok || rawVal == nil {
		return nil, &NotFoundError{key}
	}

	listVal, ok := rawVal.([]any)
	if !ok {
		return nil, &TypeAssertionError{"[]any", typeof(rawVal), key}
	}

	if len(listVal) == 0 {
		return nil, &NotFoundError{key}
	}
	listUnpacked := make([]RuleItem, 0, len(listVal))

	for i, val := range listVal {
		if val == nil {
			continue
		}
		behaviorsMap, ok := val.(map[string]any)
		if !ok {
			return nil, &TypeAssertionError{"map[string]any", typeof(val), key}
		}

		item, err := r.findRuleItem(behaviorsMap)
		if err != nil && !errors.Is(err, ErrNotFound) {
			return nil, fmt.Errorf("%s.%d: %w", key, i, err)
		}
		if errors.Is(err, ErrNotFound) {
			continue
		}

		listUnpacked = append(listUnpacked, item)
	}

	return listUnpacked, nil
}

func (r *RulesSchemaReader) findRuleItem(itemsMap map[string]any) (RuleItem, error) {
	var ruleItems RuleItems

	for name, v := range itemsMap {
		if v == nil {
			continue
		}
		items, ok := v.([]any)
		if !ok {
			return RuleItem{}, &TypeAssertionError{want: "[]any", got: typeof(v)}
		}

		if len(items) == 0 {
			continue
		}

		if items[0] != nil {
			item, ok := items[0].(map[string]any)
			if !ok {
				return RuleItem{}, &TypeAssertionError{want: "map[string]any", got: typeof(items[0])}
			}
			ruleItems = append(ruleItems, RuleItem{Name: name, Item: item})
		} else {
			ruleItems = append(ruleItems, RuleItem{Name: name, Item: defaultOptionMap()})
		}
	}

	if len(ruleItems) == 0 {
		return RuleItem{}, ErrNotFound
	}
	if len(ruleItems) > 1 {
		return RuleItem{}, &TooManyElementsError{names: ruleItems.Names(), expected: 1}
	}

	return ruleItems[0], nil
}

func (r *RulesSchemaReader) behaviorsBaseKey() string {
	return fmt.Sprintf("%s.0.behaviors.0", r.ruleFormatKey)
}

func (r *RulesSchemaReader) behaviorsKey() string {
	return fmt.Sprintf("%s.0.behavior", r.ruleFormatKey)
}

func (r *RulesSchemaReader) criteriaBaseKey() string {
	return fmt.Sprintf("%s.0.criteria.0", r.ruleFormatKey)
}

func (r *RulesSchemaReader) criteriaKey() string {
	return fmt.Sprintf("%s.0.criterion", r.ruleFormatKey)
}

func (r *RulesSchemaReader) variablesKey() string {
	return r.genericKey("variable")
}

func (r *RulesSchemaReader) childrenKey() string {
	return r.genericKey("children")
}

func (r *RulesSchemaReader) isSecureKey() string {
	return r.genericKey("is_secure")
}

func (r *RulesSchemaReader) nameKey() string {
	return r.genericKey("name")
}

func (r *RulesSchemaReader) criteriaMustSatisfyKey() string {
	return r.genericKey("criteria_must_satisfy")
}

func (r *RulesSchemaReader) uuidKey() string {
	return r.genericKey("uuid")
}

func (r *RulesSchemaReader) templateUUIDKey() string {
	return r.genericKey("template_uuid")
}

func (r *RulesSchemaReader) templateLinkKey() string {
	return r.genericKey("template_link")
}

func (r *RulesSchemaReader) criteriaLockedKey() string {
	return r.genericKey("criteria_locked")
}
func (r *RulesSchemaReader) customOverrideKey() string {
	return r.genericKey("custom_override")
}

func (r *RulesSchemaReader) advancedOverrideKey() string {
	return r.genericKey("advanced_override")
}

func (r *RulesSchemaReader) commentsKey() string {
	return r.genericKey("comments")
}

func (r *RulesSchemaReader) genericKey(k string) string {
	return fmt.Sprintf("%s.0.%s", r.ruleFormatKey, k)
}

// Names returns names of all RuleItem contained by RuleItems.
func (r RuleItems) Names() []string {
	names := make([]string, 0, len(r))
	for _, item := range r {
		names = append(names, item.Name)
	}
	return names
}

func defaultOptionMap() map[string]any {
	return map[string]any{"locked": false, "uuid": "", "template_uuid": ""}
}
