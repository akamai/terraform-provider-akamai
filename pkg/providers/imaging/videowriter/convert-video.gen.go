// Package videowriter contains set of functions used to manage video
package videowriter

import (
	"fmt"
	"regexp"
	"strconv"
	"strings"

	"github.com/akamai/AkamaiOPEN-edgegrid-golang/v12/pkg/imaging"
	"github.com/akamai/terraform-provider-akamai/v8/pkg/common/ptr"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/schema"
)

// PolicyVideoToEdgeGrid converts terraform originated map structure into EdgeGrid structure
func PolicyVideoToEdgeGrid(d *schema.ResourceData, key string) imaging.PolicyInputVideo {
	_, exist := extract(d, key)
	var result imaging.PolicyInputVideo
	if exist {
		result = imaging.PolicyInputVideo{}
		result.Breakpoints = getBreakpoints(d, getKey(key, 0, "breakpoints"))
		result.Hosts = interfaceSliceToStringSlice(d, getKey(key, 0, "hosts"))
		result.Output = getOutputVideo(d, getKey(key, 0, "output"))
		result.RolloutDuration = intReaderPtr(d, getKey(key, 0, "rollout_duration"))
		result.Variables = getVariableList(d, getKey(key, 0, "variables"))
	}

	return result
}

func getBreakpoints(d *schema.ResourceData, key string) *imaging.Breakpoints {
	_, exist := extract(d, key)
	if exist {
		result := imaging.Breakpoints{
			Widths: interfaceSliceToIntSlice(d, getKey(key, 0, "widths")),
		}
		return &result
	}
	return nil
}

func getEnumOptionsList(d *schema.ResourceData, key string) []*imaging.EnumOptions {
	collection, exist := extract(d, key)
	if exist {
		result := make([]*imaging.EnumOptions, 0)
		for idx := range collection.([]interface{}) {
			elem := imaging.EnumOptions{
				ID:    stringReader(d, getKey(key, idx, "id")),
				Value: stringReader(d, getKey(key, idx, "value")),
			}
			result = append(result, &elem)
		}
		if len(result) > 0 {
			return result
		}
		return nil
	}
	return nil
}

func getOutputVideo(d *schema.ResourceData, key string) *imaging.OutputVideo {
	_, exist := extract(d, key)
	if exist {
		result := imaging.OutputVideo{
			PerceptualQuality:    outputVideoPerceptualQualityVariableInline(d, getKey(key, 0, "perceptual_quality")),
			PlaceholderVideoURL:  stringVariableInline(d, getKey(key, 0, "placeholder_video_url")),
			VideoAdaptiveQuality: outputVideoVideoAdaptiveQualityVariableInline(d, getKey(key, 0, "video_adaptive_quality")),
		}
		return &result
	}
	return nil
}

func getVariableList(d *schema.ResourceData, key string) []imaging.Variable {
	collection, exist := extract(d, key)
	if exist {
		result := make([]imaging.Variable, 0)
		for idx := range collection.([]interface{}) {
			elem := imaging.Variable{
				DefaultValue: stringReader(d, getKey(key, idx, "default_value")),
				EnumOptions:  getEnumOptionsList(d, getKey(key, idx, "enum_options")),
				Name:         stringReader(d, getKey(key, idx, "name")),
				Postfix:      stringReaderPtr(d, getKey(key, idx, "postfix")),
				Prefix:       stringReaderPtr(d, getKey(key, idx, "prefix")),
				Type:         imaging.VariableType(stringReader(d, getKey(key, idx, "type"))),
			}
			result = append(result, elem)
		}
		if len(result) > 0 {
			return result
		}
		return nil
	}
	return nil
}

func outputVideoPerceptualQualityVariableInline(d *schema.ResourceData, key string) *imaging.OutputVideoPerceptualQualityVariableInline {
	var value *imaging.OutputVideoPerceptualQuality
	var name *string

	valueRaw, existVal := extract(d, key)
	existVal = existVal && valueRaw.(string) != ""
	if existVal {
		value = imaging.OutputVideoPerceptualQualityPtr(imaging.OutputVideoPerceptualQuality(valueRaw.(string)))
	}

	nameRaw, existVar := extract(d, key+"_var")
	existVar = existVar && nameRaw.(string) != ""
	if existVar {
		name = ptr.To(nameRaw.(string))
	}

	if existVal || existVar {
		return &imaging.OutputVideoPerceptualQualityVariableInline{
			Name:  name,
			Value: value,
		}
	}

	return nil
}

func outputVideoVideoAdaptiveQualityVariableInline(d *schema.ResourceData, key string) *imaging.OutputVideoVideoAdaptiveQualityVariableInline {
	var value *imaging.OutputVideoVideoAdaptiveQuality
	var name *string

	valueRaw, existVal := extract(d, key)
	existVal = existVal && valueRaw.(string) != ""
	if existVal {
		value = imaging.OutputVideoVideoAdaptiveQualityPtr(imaging.OutputVideoVideoAdaptiveQuality(valueRaw.(string)))
	}

	nameRaw, existVar := extract(d, key+"_var")
	existVar = existVar && nameRaw.(string) != ""
	if existVar {
		name = ptr.To(nameRaw.(string))
	}

	if existVal || existVar {
		return &imaging.OutputVideoVideoAdaptiveQualityVariableInline{
			Name:  name,
			Value: value,
		}
	}

	return nil
}

func stringVariableInline(d *schema.ResourceData, key string) *imaging.StringVariableInline {
	var value *string
	var name *string

	valueRaw, existVal := extract(d, key)
	existVal = existVal && valueRaw.(string) != ""
	if existVal {
		valueMapped := valueRaw.(string)
		value = ptr.To(valueMapped)
	}

	nameRaw, existVar := extract(d, key+"_var")
	existVar = existVar && nameRaw.(string) != ""
	if existVar {
		name = ptr.To(nameRaw.(string))
	}

	if existVal || existVar {
		return &imaging.StringVariableInline{
			Name:  name,
			Value: value,
		}
	}
	return nil
}

func intReader(d *schema.ResourceData, key string) int {
	value, exist := extract(d, key)
	if exist {
		valInt, _ := strconv.Atoi(value.(string))
		return valInt
	}
	return 0
}

func intReaderPtr(d *schema.ResourceData, key string) *int {
	value, exist := extract(d, key)
	if exist {
		if valInt, err := strconv.Atoi(value.(string)); err == nil {
			return &valInt
		}
	}
	return nil
}

func stringReader(d *schema.ResourceData, key string) string {
	value, exist := extract(d, key)
	if exist {
		return value.(string)
	}
	return ""
}

func stringReaderPtr(d *schema.ResourceData, key string) *string {
	value, exist := extract(d, key)
	if exist {
		return ptr.To(value.(string))
	}
	return nil
}

func interfaceSliceToIntSlice(d *schema.ResourceData, key string) []int {
	list, exist := extract(d, key)
	if exist {
		l := list.([]interface{})
		if len(l) > 0 {
			result := make([]int, len(l))
			for i, v := range l {
				result[i] = v.(int)
			}
			return result
		}
	}
	return nil
}

func interfaceSliceToStringSlice(d *schema.ResourceData, key string) []string {
	list, exist := extract(d, key)
	if exist {
		l := list.([]interface{})
		if len(l) > 0 {
			result := make([]string, len(l))
			for i, v := range l {
				result[i] = v.(string)
			}
			return result
		}
	}
	return nil
}

func extract(d *schema.ResourceData, key string) (interface{}, bool) {
	return d.GetOk(key)
}

func decorateKeyForSlice(key string) string {
	address := strings.Split(key, ".")
	matches, _ := regexp.MatchString("[0-9]+", address[len(address)-1])
	if !matches {
		return key + ".0"
	}
	return key
}

func getKey(rootKey string, elementIndex int, elementName string) string {
	return fmt.Sprintf("%s.%d.%s", rootKey, elementIndex, elementName)
}
