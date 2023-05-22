package datastream

import (
	"testing"

	"github.com/stretchr/testify/require"
	"github.com/tj/assert"

	"github.com/akamai/AkamaiOPEN-edgegrid-golang/v6/pkg/datastream"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/schema"
)

func newSet(items ...interface{}) *schema.Set {
	hashFunc := func(interface{}) int { return 4 } // works only for one element in set
	return schema.NewSet(hashFunc, items)
}

func TestGetConfig(t *testing.T) {
	tests := map[string]struct {
		configElements   *schema.Set
		expectedErrorMsg string
		expectedResult   datastream.Config
	}{
		"empty set": {
			configElements:   newSet(),
			expectedErrorMsg: "missing config",
		},
		"invalid config type": {
			configElements:   newSet(1),
			expectedErrorMsg: "invalid structure",
		},
		"missing frequency": {
			configElements: newSet(
				map[string]interface{}{
					"delimiter":          "SPACE",
					"format":             "STRUCTURED",
					"upload_file_prefix": "pre",
					"upload_file_suffix": "suf",
				}),
			expectedErrorMsg: "missing frequency",
		},
		"proper config": {
			configElements: newSet(
				map[string]interface{}{
					"delimiter": "SPACE",
					"format":    "STRUCTURED",
					"frequency": newSet(
						map[string]interface{}{
							"time_in_sec": 30,
						},
					),
					"upload_file_prefix": "pre",
					"upload_file_suffix": "suf",
				},
			),
			expectedResult: datastream.Config{
				Delimiter: datastream.DelimiterTypePtr(datastream.DelimiterTypeSpace),
				Format:    datastream.FormatTypeStructured,
				Frequency: datastream.Frequency{
					TimeInSec: datastream.TimeInSec30,
				},
				UploadFilePrefix: "pre",
				UploadFileSuffix: "suf",
			},
		},
	}

	for name, test := range tests {
		t.Run(name, func(t *testing.T) {
			configStruct, err := GetConfig(test.configElements)
			if test.expectedErrorMsg != "" {
				assert.Contains(t, err.Error(), test.expectedErrorMsg)
			} else {
				require.NoError(t, err)
				assert.Equal(t, &test.expectedResult, configStruct)
			}
		})
	}
}

func TestConfigToSet(t *testing.T) {
	config := datastream.Config{
		Delimiter: datastream.DelimiterTypePtr(datastream.DelimiterTypeSpace),
		Format:    datastream.FormatTypeStructured,
		Frequency: datastream.Frequency{
			TimeInSec: datastream.TimeInSec30,
		},
		UploadFilePrefix: "pre",
		UploadFileSuffix: "suf",
	}
	expected := []map[string]interface{}{
		{
			"delimiter": "SPACE",
			"format":    "STRUCTURED",
			"frequency": []map[string]interface{}{
				{
					"time_in_sec": 30,
				},
			},
			"upload_file_prefix": "pre",
			"upload_file_suffix": "suf",
		},
	}

	configSet := ConfigToSet(config)
	assert.Equal(t, expected, configSet)
}

func TestGetFrequency(t *testing.T) {
	tests := map[string]struct {
		frequencyElements *schema.Set
		expectedErrorMsg  string
		expectedResult    datastream.Frequency
	}{
		"empty set": {
			frequencyElements: newSet(),
			expectedErrorMsg:  "missing frequency",
		},
		"invalid config type": {
			frequencyElements: newSet(1),
			expectedErrorMsg:  "invalid structure",
		},
		"proper frequency": {
			frequencyElements: newSet(
				map[string]interface{}{
					"time_in_sec": 60,
				},
			),
			expectedResult: datastream.Frequency{
				TimeInSec: datastream.TimeInSec60,
			},
		},
	}

	for name, test := range tests {
		t.Run(name, func(t *testing.T) {
			frequencyStruct, err := GetFrequency(test.frequencyElements)
			if test.expectedErrorMsg != "" {
				assert.Contains(t, err.Error(), test.expectedErrorMsg)
			} else {
				require.NoError(t, err)
				assert.Equal(t, &test.expectedResult, frequencyStruct)
			}
		})
	}
}

func TestFrequencyToSet(t *testing.T) {
	frequency := datastream.Frequency{
		TimeInSec: datastream.TimeInSec60,
	}
	expected := []map[string]interface{}{
		{
			"time_in_sec": 60,
		},
	}

	frequencySet := FrequencyToSet(frequency)
	assert.Equal(t, expected, frequencySet)
}

func TestInterfaceSliceToIntSlice(t *testing.T) {
	tests := map[string]struct {
		input    []interface{}
		expected []int
	}{
		"empty list": {
			input:    []interface{}{},
			expected: []int{},
		},
		"list with values": {
			input:    []interface{}{1, 2, 3},
			expected: []int{1, 2, 3},
		},
	}

	for name, test := range tests {
		t.Run(name, func(t *testing.T) {
			assert.Equal(t, test.expected, InterfaceSliceToIntSlice(test.input))
		})
	}
}

func TestInterfaceSliceToStringSlice(t *testing.T) {
	tests := map[string]struct {
		input    []interface{}
		expected []string
	}{
		"empty list": {
			input:    []interface{}{},
			expected: []string{},
		},
		"list with values": {
			input:    []interface{}{"a", "b", "c"},
			expected: []string{"a", "b", "c"},
		},
	}

	for name, test := range tests {
		t.Run(name, func(t *testing.T) {
			assert.Equal(t, test.expected, InterfaceSliceToStringSlice(test.input))
		})
	}
}

func TestDataSetFieldsToList(t *testing.T) {
	datasets := []datastream.DataSets{
		{
			DatasetGroupName: "group1",
			DatasetFields: []datastream.DatasetFields{
				{
					DatasetFieldID: 1000,
					Order:          3,
				},
				{
					DatasetFieldID: 1002,
					Order:          5,
				},
				{
					DatasetFieldID: 1100,
					Order:          1,
				},
			},
		},
		{
			DatasetGroupName: "group2",
			DatasetFields: []datastream.DatasetFields{
				{
					DatasetFieldID: 2000,
					Order:          4,
				},
				{
					DatasetFieldID: 2002,
					Order:          0,
				},
				{
					DatasetFieldID: 2100,
					Order:          2,
				},
			},
		},
	}

	assert.Equal(t, []int{2002, 1100, 2100, 1000, 2000, 1002}, DataSetFieldsToList(datasets))
}

func TestPropertyToList(t *testing.T) {
	properties := []datastream.Property{
		{
			PropertyID:   1,
			PropertyName: "property_1",
		},
		{
			PropertyID:   2,
			PropertyName: "property_2",
		},
		{
			PropertyID:   3,
			PropertyName: "property_3",
		},
	}

	assert.Equal(t, []string{"1", "2", "3"}, PropertyToList(properties))
}

func TestGetPropertiesList(t *testing.T) {
	properties := []interface{}{
		"1",
		"2",
		"prp_3",
		"4",
		"prp_5",
	}

	result, err := GetPropertiesList(properties)
	require.NoError(t, err)

	assert.Equal(t, []int{1, 2, 3, 4, 5}, result)
}
