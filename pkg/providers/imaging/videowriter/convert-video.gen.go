package videowriter

import (
	"reflect"

	"github.com/akamai/AkamaiOPEN-edgegrid-golang/v2/pkg/imaging"
)

// PolicyVideoToEdgeGrid converts terraform originated map structure into EdgeGrid structure
func PolicyVideoToEdgeGrid(input map[string]interface{}) imaging.PolicyInputVideo {
	result := imaging.PolicyInputVideo{}
	result.Breakpoints = getBreakpoints(extract(input, "breakpoints"))
	result.Hosts = interfaceSliceToStringSlice(input["hosts"].([]interface{}))
	result.Output = getOutputVideo(extract(input, "output"))
	result.RolloutDuration = input["rollout_duration"].(int)
	result.Variables = getVariableList(input["variables"].([]interface{}))
	return result
}

func getBreakpoints(src map[string]interface{}) *imaging.Breakpoints {
	if src == nil {
		return nil
	}
	result := imaging.Breakpoints{
		Widths: interfaceSliceToIntSlice(src["widths"].([]interface{})),
	}
	return &result
}

func getEnumOptionsList(src []interface{}) []*imaging.EnumOptions {
	result := make([]*imaging.EnumOptions, 0)
	for idx := range src {
		elem := imaging.EnumOptions{
			ID:    src[idx].(map[string]interface{})["id"].(string),
			Value: src[idx].(map[string]interface{})["value"].(string),
		}
		result = append(result, &elem)
	}
	if len(result) > 0 {
		return result
	}
	return nil
}

func getOutputVideo(src map[string]interface{}) *imaging.OutputVideo {
	if src == nil {
		return nil
	}
	result := imaging.OutputVideo{
		PerceptualQuality:    outputVideoPerceptualQualityVariableInline(src, "perceptual_quality"),
		PlaceholderVideoURL:  stringVariableInline(src, "placeholder_video_url"),
		VideoAdaptiveQuality: outputVideoVideoAdaptiveQualityVariableInline(src, "video_adaptive_quality"),
	}
	return &result
}

func getVariableList(src []interface{}) []imaging.Variable {
	result := make([]imaging.Variable, 0)
	for idx := range src {
		elem := imaging.Variable{
			DefaultValue: src[idx].(map[string]interface{})["default_value"].(string),
			EnumOptions:  getEnumOptionsList(src[idx].(map[string]interface{})["enum_options"].([]interface{})),
			Name:         src[idx].(map[string]interface{})["name"].(string),
			Postfix:      src[idx].(map[string]interface{})["postfix"].(string),
			Prefix:       src[idx].(map[string]interface{})["prefix"].(string),
			Type:         imaging.VariableType(src[idx].(map[string]interface{})["type"].(string)),
		}
		result = append(result, elem)
	}
	if len(result) > 0 {
		return result
	}
	return nil
}

func outputVideoPerceptualQualityVariableInline(src map[string]interface{}, name string) *imaging.OutputVideoPerceptualQualityVariableInline {
	if !variableHasValue(src, name) {
		return nil
	}

	v1 := src[name]
	var v2 *imaging.OutputVideoPerceptualQuality
	if v1 != "" {
		v2 = imaging.OutputVideoPerceptualQualityPtr(imaging.OutputVideoPerceptualQuality(v1.(string)))
	}

	return &imaging.OutputVideoPerceptualQualityVariableInline{
		Name:  stringValuePtr(src, name+"_var"),
		Value: v2,
	}
}

func outputVideoVideoAdaptiveQualityVariableInline(src map[string]interface{}, name string) *imaging.OutputVideoVideoAdaptiveQualityVariableInline {
	if !variableHasValue(src, name) {
		return nil
	}

	v1 := src[name]
	var v2 *imaging.OutputVideoVideoAdaptiveQuality
	if v1 != "" {
		v2 = imaging.OutputVideoVideoAdaptiveQualityPtr(imaging.OutputVideoVideoAdaptiveQuality(v1.(string)))
	}

	return &imaging.OutputVideoVideoAdaptiveQualityVariableInline{
		Name:  stringValuePtr(src, name+"_var"),
		Value: v2,
	}
}

func stringVariableInline(src map[string]interface{}, name string) *imaging.StringVariableInline {
	if variableHasValue(src, name) {
		return &imaging.StringVariableInline{
			Name:  stringValuePtr(src, name+"_var"),
			Value: stringValuePtr(src, name),
		}
	}
	return nil
}

func stringValuePtr(src map[string]interface{}, name string) *string {
	value := src[name]
	if value != "" {
		v := value.(string)
		return &v
	}
	return nil
}

func interfaceSliceToIntSlice(list []interface{}) []int {
	if len(list) == 0 {
		return nil
	}
	intList := make([]int, len(list))
	for i, v := range list {
		intList[i] = v.(int)
	}
	return intList
}

func interfaceSliceToStringSlice(list []interface{}) []string {
	if len(list) == 0 {
		return nil
	}
	stringList := make([]string, len(list))
	for i, v := range list {
		stringList[i] = v.(string)
	}
	return stringList
}

func variableHasValue(src map[string]interface{}, name string) bool {
	v1 := src[name]
	v2 := src[name+"_var"]

	if !reflect.ValueOf(v1).IsZero() || !reflect.ValueOf(v2).IsZero() {
		return true
	}
	return false
}

func extract(src map[string]interface{}, name string) map[string]interface{} {
	elem, ok := src[name]
	if !ok {
		return nil
	}

	l := elem.([]interface{})
	if len(l) == 0 {
		return nil
	}
	return l[0].(map[string]interface{})
}
