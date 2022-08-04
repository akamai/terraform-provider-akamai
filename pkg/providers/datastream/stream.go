package datastream

import (
	"fmt"
	"sort"
	"strconv"
	"strings"

	"github.com/akamai/AkamaiOPEN-edgegrid-golang/v2/pkg/datastream"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/schema"
)

// GetConfig builds Config structure
func GetConfig(set *schema.Set) (*datastream.Config, error) {
	if set.Len() != 1 {
		return nil, fmt.Errorf("missing config definition")
	}

	configList := set.List()
	configMap, ok := configList[0].(map[string]interface{})
	if !ok {
		return nil, fmt.Errorf("config has invalid structure")
	}

	frequencySet, ok := configMap["frequency"]
	if !ok {
		return nil, fmt.Errorf("missing frequency block in configuration")
	}

	frequency, err := GetFrequency(frequencySet.(*schema.Set))
	if err != nil {
		return nil, err
	}

	var delimiterPtr *datastream.DelimiterType
	if delimiterStr := configMap["delimiter"].(string); delimiterStr != "" {
		delimiterPtr = datastream.DelimiterTypePtr(datastream.DelimiterType(delimiterStr))
	}

	return &datastream.Config{
		Delimiter:        delimiterPtr,
		Format:           datastream.FormatType(configMap["format"].(string)),
		Frequency:        *frequency,
		UploadFilePrefix: configMap["upload_file_prefix"].(string),
		UploadFileSuffix: configMap["upload_file_suffix"].(string),
	}, nil
}

// ConfigToSet converts Config struct to set
func ConfigToSet(cfg datastream.Config) []map[string]interface{} {
	delimiter := *datastream.DelimiterTypePtr("")
	if cfg.Delimiter != nil {
		delimiter = *cfg.Delimiter
	}

	return []map[string]interface{}{{
		"delimiter":          string(delimiter),
		"format":             string(cfg.Format),
		"frequency":          FrequencyToSet(cfg.Frequency),
		"upload_file_prefix": cfg.UploadFilePrefix,
		"upload_file_suffix": cfg.UploadFileSuffix,
	}}
}

// GetFrequency builds Frequency structure
func GetFrequency(set *schema.Set) (*datastream.Frequency, error) {
	if set.Len() != 1 {
		return nil, fmt.Errorf("missing frequency definition")
	}

	freqList := set.List()
	freqMap, ok := freqList[0].(map[string]interface{})
	if !ok {
		return nil, fmt.Errorf("frequency has invalid structure")
	}

	return &datastream.Frequency{
		TimeInSec: datastream.TimeInSec(freqMap["time_in_sec"].(int)),
	}, nil
}

// FrequencyToSet converts Frequency struct to map
func FrequencyToSet(freq datastream.Frequency) []map[string]interface{} {
	return []map[string]interface{}{{
		"time_in_sec": int(freq.TimeInSec),
	}}
}

// InterfaceSliceToIntSlice converts schema.Set to slice of ints
func InterfaceSliceToIntSlice(list []interface{}) []int {
	intList := make([]int, len(list))
	for i, v := range list {
		intList[i] = v.(int)
	}
	return intList
}

// InterfaceSliceToStringSlice converts schema.Set to slice of ints
func InterfaceSliceToStringSlice(list []interface{}) []string {
	stringList := make([]string, len(list))
	for i, v := range list {
		stringList[i] = v.(string)
	}
	return stringList
}

// DataSetFieldsToList converts slice of DataSets to slice of ints
func DataSetFieldsToList(dataSets []datastream.DataSets) []int {
	datasetFields := make([]datastream.DatasetFields, 0)

	for _, datasetGroup := range dataSets {
		datasetFields = append(datasetFields, datasetGroup.DatasetFields...)
	}

	sort.Slice(datasetFields, func(i, j int) bool {
		return datasetFields[i].Order < datasetFields[j].Order
	})

	ids := make([]int, 0, len(datasetFields))

	for _, field := range datasetFields {
		ids = append(ids, field.DatasetFieldID)
	}

	return ids
}

// PropertyToList converts slice of Properties to slice of ints
func PropertyToList(properties []datastream.Property) []string {
	ids := make([]string, 0, len(properties))

	for _, property := range properties {
		ids = append(ids, strconv.Itoa(property.PropertyID))
	}

	return ids
}

// GetPropertiesList converts propertyIDs with and without "prp_" prefix to slice of ints
func GetPropertiesList(properties []interface{}) ([]int, error) {
	ids := make([]int, 0, len(properties))

	for _, property := range properties {
		propertyID, err := strconv.Atoi(strings.TrimPrefix(property.(string), "prp_"))
		if err != nil {
			return nil, err
		}
		ids = append(ids, propertyID)
	}

	return ids, nil
}

// StreamIDToPapiJSON generates PAPI JSON with given id of a stream
func StreamIDToPapiJSON(id int64) string {
	return fmt.Sprintf(
		`
{
    "name": "Datastream Rule",
    "children": [],
    "behaviors": [
        {
            "name": "datastream",
            "options": {
                "streamType": "LOG",
                "logEnabled": true,
                "logStreamName": %d,
                "samplingPercentage": 100
            }
        }
    ],
    "criteria": [],
    "criteriaMustSatisfy": "all"
}
`, id)
}
