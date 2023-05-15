package datastream

import (
	"fmt"
	"strconv"
	"strings"

	"github.com/akamai/AkamaiOPEN-edgegrid-golang/v6/pkg/datastream"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/schema"
)

// GetConfig builds Config structure
func GetConfig(set *schema.Set) (*datastream.DeliveryConfiguration, error) {
	if set.Len() != 1 {
		return nil, fmt.Errorf("missing delivery configuration definition")
	}

	configList := set.List()
	configMap, ok := configList[0].(map[string]interface{})
	if !ok {
		return nil, fmt.Errorf("delivery configuration has invalid structure")
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
	if delimiterStr := configMap["field_delimiter"].(string); delimiterStr != "" {
		delimiterPtr = datastream.DelimiterTypePtr(datastream.DelimiterType(delimiterStr))
	}

	return &datastream.DeliveryConfiguration{
		Delimiter:        delimiterPtr,
		Format:           datastream.FormatType(configMap["format"].(string)),
		Frequency:        *frequency,
		UploadFilePrefix: configMap["upload_file_prefix"].(string),
		UploadFileSuffix: configMap["upload_file_suffix"].(string),
	}, nil
}

// ConfigToSet converts Config struct to set
func ConfigToSet(cfg datastream.DeliveryConfiguration) []map[string]interface{} {
	delimiter := *datastream.DelimiterTypePtr("")
	if cfg.Delimiter != nil {
		delimiter = *cfg.Delimiter
	}

	return []map[string]interface{}{{
		"field_delimiter":    string(delimiter),
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
		IntervalInSeconds: datastream.IntervalInSeconds(freqMap["interval_in_secs"].(int)),
	}, nil
}

// FrequencyToSet converts Frequency struct to map
func FrequencyToSet(freq datastream.Frequency) []map[string]interface{} {
	return []map[string]interface{}{{
		"interval_in_secs": int(freq.IntervalInSeconds),
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

// DatasetFieldListToDatasetFields converts schema.Set to slice of DatasetFieldId
func DatasetFieldListToDatasetFields(list []interface{}) []datastream.DatasetFieldID {

	datasetFields := make([]datastream.DatasetFieldID, 0)

	for _, v := range list {
		datasetFields = append(datasetFields, datastream.DatasetFieldID{v.(int)})
	}
	return datasetFields
}

// InterfaceSliceToStringSlice converts schema.Set to slice of string
func InterfaceSliceToStringSlice(list []interface{}) []string {
	stringList := make([]string, len(list))
	for i, v := range list {
		stringList[i] = v.(string)
	}
	return stringList
}

// DataSetFieldsToList converts slice of dataSetFields to slice of ints
func DataSetFieldsToList(dataSetFields []datastream.DataSetField) []int {

	ids := make([]int, 0, len(dataSetFields))

	for _, field := range dataSetFields {
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
func GetPropertiesList(properties []interface{}) ([]datastream.PropertyID, error) {
	ids := make([]datastream.PropertyID, 0, len(properties))

	for _, property := range properties {
		propertyID, err := strconv.Atoi(strings.TrimPrefix(property.(string), "prp_"))
		if err != nil {
			return nil, err
		}
		ids = append(ids, datastream.PropertyID{propertyID})
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
