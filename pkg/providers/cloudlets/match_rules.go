package cloudlets

import (
	"crypto/sha1"
	"encoding/hex"
	"fmt"
	"io"
	"strconv"

	"github.com/akamai/AkamaiOPEN-edgegrid-golang/v5/pkg/cloudlets"
	"github.com/akamai/terraform-provider-akamai/v3/pkg/tools"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/schema"
)

type (
	// objectMatchValueHandler is type alias function for casting ObjectMatchValue into a specified type
	objectMatchValueHandler func(map[string]interface{}, string) (interface{}, error)
)

func getMatchRulesHashID(matchRules cloudlets.MatchRules) (string, error) {
	id := "id"
	for _, rule := range matchRules {
		switch r := rule.(type) {
		case cloudlets.MatchRuleER:
			id = id + ":" + r.Name
		}
	}
	h := sha1.New()
	_, err := io.WriteString(h, id)
	if err != nil {
		return "", err
	}
	hashID := hex.EncodeToString(h.Sum(nil))
	return hashID, nil
}

func getStringValue(matchRuleMap map[string]interface{}, name string) string {
	if value, ok := matchRuleMap[name]; ok {
		return value.(string)
	}
	return ""
}

func getIntValue(matchRuleMap map[string]interface{}, name string) int {
	if value, ok := matchRuleMap[name]; ok {
		return value.(int)
	}
	return 0
}

// this will not work on 32bit platform if the value is bigger than max for int32
func getInt64Value(matchRuleMap map[string]interface{}, name string) int64 {
	if value, ok := matchRuleMap[name]; ok {
		return int64(value.(int))
	}
	return 0
}

func getFloat64PtrValue(matchRuleMap map[string]interface{}, name string) *float64 {
	if value, ok := matchRuleMap[name]; ok {
		v := value.(float64)
		return &v
	}
	return nil
}

func getBoolValue(matchRuleMap map[string]interface{}, name string) bool {
	if value, ok := matchRuleMap[name]; ok {
		return value.(bool)
	}
	return false
}

func getListOfStringsValue(matchRuleMap map[string]interface{}, name string) []string {
	if value, ok := matchRuleMap[name]; ok {
		var val []string
		for _, v := range value.([]interface{}) {
			val = append(val, v.(string))
		}
		return val
	}
	return nil
}

func getOMVSimpleType(omv map[string]interface{}) interface{} {
	simpleType := cloudlets.ObjectMatchValueSimple{
		Type:  cloudlets.Simple,
		Value: getListOfStringsValue(omv, "value"),
	}
	return &simpleType
}

func getOMVObjectType(omv map[string]interface{}) (interface{}, error) {
	opts, err := parseOMVOptions(omv)
	if err != nil {
		return nil, err
	}
	objectType := cloudlets.ObjectMatchValueObject{
		Type:              cloudlets.Object,
		Name:              getStringValue(omv, "name"),
		NameCaseSensitive: getBoolValue(omv, "name_case_sensitive"),
		NameHasWildcard:   getBoolValue(omv, "name_has_wildcard"),
		Options:           opts,
	}
	return &objectType, nil
}

func getOMVRangeType(omv map[string]interface{}) (interface{}, error) {
	valuesAsString := getListOfStringsValue(omv, "value")
	var valuesAsInt []int64
	for _, valueAsString := range valuesAsString {
		valueAsInt, err := strconv.ParseInt(valueAsString, 10, 64)
		if err != nil {
			return nil, fmt.Errorf("cannot parse %s value as an integer: %s", valueAsString, err)
		}
		valuesAsInt = append(valuesAsInt, valueAsInt)
	}

	rangeType := cloudlets.ObjectMatchValueRange{
		Type:  cloudlets.Range,
		Value: valuesAsInt,
	}
	return &rangeType, nil
}

func parseOMVOptions(omvOptions map[string]interface{}) (*cloudlets.Options, error) {
	o, ok := omvOptions["options"]
	if !ok {
		return nil, nil
	}
	options := o.(*schema.Set).List()
	if len(options) < 1 {
		return nil, nil
	}

	optionFields, ok := options[0].(map[string]interface{})
	if !ok {
		return nil, fmt.Errorf("%w: 'options' should be an object", tools.ErrInvalidType)
	}
	option := cloudlets.Options{
		Value:              getListOfStringsValue(optionFields, "value"),
		ValueHasWildcard:   getBoolValue(optionFields, "value_has_wildcard"),
		ValueCaseSensitive: getBoolValue(optionFields, "value_case_sensitive"),
		ValueEscaped:       getBoolValue(optionFields, "value_escaped"),
	}
	return &option, nil
}

func setMatchRuleSchemaType(matchRules []interface{}, t cloudlets.MatchRuleType) error {
	for _, mr := range matchRules {
		matchRuleMap, ok := mr.(map[string]interface{})
		if !ok {
			return fmt.Errorf("match rule is of invalid type: %T", mr)
		}
		matchRuleMap["type"] = t
	}
	return nil
}

func parseObjectMatchValue(criteriaMap map[string]interface{}, handler objectMatchValueHandler) (interface{}, error) {
	v, ok := criteriaMap["object_match_value"]
	if !ok {
		return nil, nil
	}

	rawObjects := v.(*schema.Set).List()
	if len(rawObjects) != 0 {
		omv, ok := rawObjects[0].(map[string]interface{})
		if !ok {
			return nil, fmt.Errorf("%w: 'object_match_value' should be an object", tools.ErrInvalidType)
		}
		if omvType, ok := omv["type"]; ok {
			return handler(omv, omvType.(string))
		}
	}
	return nil, nil
}

func getObjectMatchValueObjectOrSimpleOrRange(omv map[string]interface{}, t string) (interface{}, error) {
	if cloudlets.ObjectMatchValueObjectType(t) == cloudlets.Object {
		return getOMVObjectType(omv)
	}
	if cloudlets.ObjectMatchValueSimpleType(t) == cloudlets.Simple {
		return getOMVSimpleType(omv), nil
	}
	if cloudlets.ObjectMatchValueRangeType(t) == cloudlets.Range {
		return getOMVRangeType(omv)
	}
	return nil, fmt.Errorf("'object_match_value' type '%s' is invalid. Must be one of: 'simple', 'range' or 'object'", t)
}

func getObjectMatchValueObjectOrSimple(omv map[string]interface{}, t string) (interface{}, error) {
	if cloudlets.ObjectMatchValueObjectType(t) == cloudlets.Object {
		return getOMVObjectType(omv)
	}
	if cloudlets.ObjectMatchValueSimpleType(t) == cloudlets.Simple {
		return getOMVSimpleType(omv), nil
	}
	return nil, fmt.Errorf("'object_match_value' type '%s' is invalid. Must be one of: 'simple' or 'object'", t)
}
