package datastream

import (
	"testing"

	"github.com/akamai/AkamaiOPEN-edgegrid-golang/v14/pkg/datastream"
	"github.com/akamai/terraform-provider-akamai/v11/pkg/common/tf"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/schema"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

func newSet(items ...interface{}) *schema.Set {
	hashFunc := func(interface{}) int { return 4 } // works only for one element in set
	return schema.NewSet(hashFunc, items)
}

func TestGetConfig(t *testing.T) {
	tests := map[string]struct {
		configElements   *schema.Set
		expectedErrorMsg string
		expectedResult   datastream.DeliveryConfiguration
	}{
		"empty set": {
			configElements:   newSet(),
			expectedErrorMsg: "missing delivery configuration",
		},
		"invalid config type": {
			configElements:   newSet(1),
			expectedErrorMsg: "invalid structure",
		},
		"missing frequency": {
			configElements: newSet(
				map[string]interface{}{
					"field_delimiter":    "SPACE",
					"format":             "STRUCTURED",
					"upload_file_prefix": "pre",
					"upload_file_suffix": "suf",
				}),
			expectedErrorMsg: "missing frequency",
		},
		"proper config": {
			configElements: newSet(
				map[string]interface{}{
					"field_delimiter": "SPACE",
					"format":          "STRUCTURED",
					"frequency": newSet(
						map[string]interface{}{
							"interval_in_secs": 30,
						},
					),
					"upload_file_prefix": "pre",
					"upload_file_suffix": "suf",
				},
			),
			expectedResult: datastream.DeliveryConfiguration{
				Delimiter: datastream.DelimiterTypePtr(datastream.DelimiterTypeSpace),
				Format:    datastream.FormatTypeStructured,
				Frequency: datastream.Frequency{
					IntervalInSeconds: datastream.IntervalInSeconds30,
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
	config := datastream.DeliveryConfiguration{
		Delimiter: datastream.DelimiterTypePtr(datastream.DelimiterTypeSpace),
		Format:    datastream.FormatTypeStructured,
		Frequency: datastream.Frequency{
			IntervalInSeconds: datastream.IntervalInSeconds30,
		},
		UploadFilePrefix: "pre",
		UploadFileSuffix: "suf",
	}
	expected := []map[string]interface{}{
		{
			"field_delimiter": "SPACE",
			"format":          "STRUCTURED",
			"frequency": []map[string]interface{}{
				{
					"interval_in_secs": 30,
				},
			},
			"upload_file_prefix": "pre",
			"upload_file_suffix": "suf",
		},
	}

	configSet := ConfigToSet(config)
	assert.Equal(t, expected, configSet)
}

func testResourceDataWithDeliveryConfig(t *testing.T, prefix, suffix string) *schema.ResourceData {
	t.Helper()

	d := schema.TestResourceDataRaw(t, datastreamResourceSchema, map[string]interface{}{
		"active":      true,
		"stream_name": "test_stream",
	})

	if err := d.Set("delivery_configuration", newSet(map[string]interface{}{
		"field_delimiter":    "SPACE",
		"format":             "STRUCTURED",
		"upload_file_prefix": prefix,
		"upload_file_suffix": suffix,
		"frequency": newSet(map[string]interface{}{
			"interval_in_secs": 30,
		}),
	})); err != nil {
		t.Fatalf("setting delivery_configuration: %v", err)
	}

	return d
}

func TestResolveUploadFilePrefixSuffixForRead(t *testing.T) {
	t.Parallel()

	d := testResourceDataWithDeliveryConfig(t, "pre", "suf")

	t.Run("prefers API values when present", func(t *testing.T) {
		assert.Equal(t, "api-pre", resolveUploadFilePrefixForRead("api-pre", d))
		assert.Equal(t, "api-suf", resolveUploadFileSuffixForRead("api-suf", d))
	})

	t.Run("preserves configured values when API omits prefix and suffix", func(t *testing.T) {
		assert.Equal(t, "pre", resolveUploadFilePrefixForRead("", d))
		assert.Equal(t, "suf", resolveUploadFileSuffixForRead("", d))
	})

	t.Run("uses defaults when API and config omit prefix and suffix", func(t *testing.T) {
		empty := schema.TestResourceDataRaw(t, datastreamResourceSchema, map[string]interface{}{})
		assert.Equal(t, DefaultUploadFilePrefix, resolveUploadFilePrefixForRead("", empty))
		assert.Equal(t, DefaultUploadFileSuffix, resolveUploadFileSuffixForRead("", empty))
	})
}

func TestApplyDeliveryConfigurationForRead(t *testing.T) {
	t.Parallel()

	d := testResourceDataWithDeliveryConfig(t, "pre", "suf")

	cfg := datastream.DeliveryConfiguration{}
	applyDeliveryConfigurationForRead(&cfg, d, "s3_connector")
	assert.Equal(t, "pre", cfg.UploadFilePrefix)
	assert.Equal(t, "suf", cfg.UploadFileSuffix)

	sumologicCfg := datastream.DeliveryConfiguration{}
	applyDeliveryConfigurationForRead(&sumologicCfg, d, "sumologic_connector")
	assert.Equal(t, DefaultUploadFilePrefix, sumologicCfg.UploadFilePrefix)
	assert.Equal(t, DefaultUploadFileSuffix, sumologicCfg.UploadFileSuffix)
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
					"interval_in_secs": 60,
				},
			),
			expectedResult: datastream.Frequency{
				IntervalInSeconds: datastream.IntervalInSeconds60,
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
		IntervalInSeconds: datastream.IntervalInSeconds60,
	}
	expected := []map[string]interface{}{
		{
			"interval_in_secs": 60,
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
			assert.Equal(t, test.expected, tf.InterfaceSliceToStringSlice(test.input))
		})
	}
}

func TestDataSetFieldsToList(t *testing.T) {

	datasets := datastream.DataSets{
		DataSetFields: []datastream.DataSetField{
			{
				DatasetFieldID: 1000,
			},
			{
				DatasetFieldID: 1002,
			},
			{
				DatasetFieldID: 1100,
			},
			{
				DatasetFieldID: 2000,
			},
			{
				DatasetFieldID: 2002,
			},
			{
				DatasetFieldID: 2100,
			},
		},
	}
	assert.Equal(t, []int{1000, 1002, 1100, 2000, 2002, 2100}, DataSetFieldsToList(datasets.DataSetFields))
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

	propertyIDs := make([]int, len(result))
	for i := 0; i < len(result); i++ {
		propertyIDs[i] = result[i].PropertyID
	}
	assert.Equal(t, []int{1, 2, 3, 4, 5}, propertyIDs)
}

func TestAppSecConfigsToList(t *testing.T) {
	configs := []datastream.AppSecConfig{
		{AppSecID: 16536, AppSecName: "WAF Security File"},
		{AppSecID: 67890, AppSecName: "Bot Manager Config"},
	}

	assert.Equal(t, []int{16536, 67890}, AppSecConfigsToIDsList(configs))
}

func TestGetAppSecConfigIDs(t *testing.T) {
	configs := []interface{}{16536, 67890}

	result, err := GetAppSecConfigIDs(configs)
	require.NoError(t, err)

	ids := make([]int, len(result))
	for i, r := range result {
		ids[i] = r.AppSecID
	}
	assert.Equal(t, []int{16536, 67890}, ids)
}

func TestNormalizePropertyID(t *testing.T) {
	t.Parallel()

	assert.Equal(t, "12345", normalizePropertyID("prp_12345"))
	assert.Equal(t, "prp_12345", normalizePropertyID("  prp_12345  "))
	assert.Equal(t, "42", normalizePropertyID("42"))
}

func TestPropertiesSameSet(t *testing.T) {
	t.Parallel()

	tests := map[string]struct {
		old      []interface{}
		new      []interface{}
		expected bool
	}{
		"same order": {
			old:      []interface{}{"1", "2", "3"},
			new:      []interface{}{"1", "2", "3"},
			expected: true,
		},
		"different order": {
			old:      []interface{}{"1", "2", "3"},
			new:      []interface{}{"3", "1", "2"},
			expected: true,
		},
		"prp_ prefix normalized": {
			old:      []interface{}{"prp_1", "2"},
			new:      []interface{}{"1", "prp_2"},
			expected: true,
		},
		"different members": {
			old:      []interface{}{"1", "2"},
			new:      []interface{}{"1", "3"},
			expected: false,
		},
		"different lengths": {
			old:      []interface{}{"1", "2"},
			new:      []interface{}{"1"},
			expected: false,
		},
	}

	for name, test := range tests {
		t.Run(name, func(t *testing.T) {
			t.Parallel()
			same, ok := propertiesSameSet(test.old, test.new)
			require.True(t, ok)
			assert.Equal(t, test.expected, same)
		})
	}
}

func TestPropertiesSameSet_invalidType(t *testing.T) {
	t.Parallel()

	same, ok := propertiesSameSet([]interface{}{1, "2"}, []interface{}{"1", "2"})
	assert.False(t, ok)
	assert.False(t, same)
}

func TestIsPropertiesOrderDifferent(t *testing.T) {
	t.Parallel()

	t.Run("falls back to hash comparison when properties attribute unchanged", func(t *testing.T) {
		d := schema.TestResourceDataRaw(t, datastreamResourceSchema, map[string]interface{}{
			"log_type":    "CDN",
			"active":      false,
			"stream_name": "test_stream",
			"properties":  []interface{}{"1", "2"},
		})

		assert.True(t, isPropertiesOrderDifferent("properties", "same-hash", "same-hash", d))
		assert.False(t, isPropertiesOrderDifferent("properties", "old-hash", "new-hash", d))
	})
}

// Reordering coverage for DiffSuppressFunc is exercised by
// TestResourceStreamPropertiesOrderDiffSuppress (PlanOnly with reordered properties).
