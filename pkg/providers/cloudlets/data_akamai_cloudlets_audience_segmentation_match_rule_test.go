package cloudlets

import (
	"regexp"
	"testing"

	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/resource"
)

func TestDataCloudletsAudienceSegmentationMatchRule(t *testing.T) {
	tests := map[string]struct {
		configPath       string
		expectedJSONPath string
	}{
		"basic valid rule set": {
			configPath:       "testdata/TestDataCloudletsAudienceSegmentationMatchRule/basic.tf",
			expectedJSONPath: "testdata/TestDataCloudletsAudienceSegmentationMatchRule/rules/basic_rules.json",
		},
		"valid rule set with simple value": {
			configPath:       "testdata/TestDataCloudletsAudienceSegmentationMatchRule/omv_simple.tf",
			expectedJSONPath: "testdata/TestDataCloudletsAudienceSegmentationMatchRule/rules/omv_simple_rules.json",
		},
		"valid rule set with object value": {
			configPath:       "testdata/TestDataCloudletsAudienceSegmentationMatchRule/omv_object.tf",
			expectedJSONPath: "testdata/TestDataCloudletsAudienceSegmentationMatchRule/rules/omv_object_rules.json",
		},
		"valid rule set with range value": {
			configPath:       "testdata/TestDataCloudletsAudienceSegmentationMatchRule/omv_range.tf",
			expectedJSONPath: "testdata/TestDataCloudletsAudienceSegmentationMatchRule/rules/omv_range_rules.json",
		},
		"valid rule set without value": {
			configPath:       "testdata/TestDataCloudletsAudienceSegmentationMatchRule/omv_empty.tf",
			expectedJSONPath: "testdata/TestDataCloudletsAudienceSegmentationMatchRule/rules/omv_empty_rules.json",
		},
		"valid complex rule set": {
			configPath:       "testdata/TestDataCloudletsAudienceSegmentationMatchRule/omv_complex.tf",
			expectedJSONPath: "testdata/TestDataCloudletsAudienceSegmentationMatchRule/rules/omv_complex_rules.json",
		},
	}

	for name, test := range tests {
		t.Run(name, func(t *testing.T) {
			resource.UnitTest(t, resource.TestCase{
				Providers: testAccProviders,
				Steps: []resource.TestStep{
					{
						Config: loadFixtureString(test.configPath),
						Check: resource.ComposeAggregateTestCheckFunc(
							resource.TestCheckResourceAttr(
								"data.akamai_cloudlets_audience_segmentation_match_rule.test", "json",
								loadFixtureString(test.expectedJSONPath)),
							resource.TestCheckResourceAttr(
								"data.akamai_cloudlets_audience_segmentation_match_rule.test", "match_rules.0.type", "asMatchRule"),
						),
					},
				},
			})
		})
	}
}

func TestIncorrectDataCloudletsAudienceSegmentationMatchRule(t *testing.T) {
	tests := map[string]struct {
		configPath string
		withError  string
	}{
		"match criteria AS - missed forward_settings field": {
			configPath: "testdata/TestDataCloudletsAudienceSegmentationMatchRule/missing_forward_settings.tf",
			withError:  `At least 1 "forward_settings" blocks are required.`,
		},
		"match criteria AS - missed type field in ObjectMatchValue": {
			configPath: "testdata/TestDataCloudletsAudienceSegmentationMatchRule/missing_type.tf",
			withError:  "Missing required argument",
		},
		"match criteria AS - no match_value and object_match_value": {
			configPath: "testdata/TestDataCloudletsAudienceSegmentationMatchRule/missing_value.tf",
			withError:  `(?s)cannot be blank when ObjectMatchValue is blank.*cannot be blank when MatchValue is blank`,
		},
		"match criteria AS - match_value and object_match_value together": {
			configPath: "testdata/TestDataCloudletsAudienceSegmentationMatchRule/duplicate_values.tf",
			withError:  `(?s)must be blank when ObjectMatchValue is set.*must be blank when MatchValue is set`,
		},
		"match criteria AS - incorrect value of type field": {
			configPath: "testdata/TestDataCloudletsAudienceSegmentationMatchRule/invalid_enum_type.tf",
			withError:  `expected type to be one of \['simple', 'object', 'range'\], got incorrect_type`,
		},
		"match criteria AS - incorrect value of match_operator field": {
			configPath: "testdata/TestDataCloudletsAudienceSegmentationMatchRule/invalid_enum_match_operator.tf",
			withError:  `expected match_operator to be one of \['contains', 'exists', 'equals', ''\], got invalid_operator`,
		},
		"match criteria AS - incorrect value of check_ips field": {
			configPath: "testdata/TestDataCloudletsAudienceSegmentationMatchRule/invalid_enum_check_ips.tf",
			withError:  `expected check_ips to be one of \['CONNECTING_IP', 'XFF_HEADERS', 'CONNECTING_IP XFF_HEADERS', ''\], got incorrect_check_ips`,
		},
		"match criteria AS - incorrect value of match_type field": {
			configPath: "testdata/TestDataCloudletsAudienceSegmentationMatchRule/invalid_enum_match_type.tf",
			withError:  `expected match_type to be one of \['header', 'hostname', 'path', 'extension', 'query', 'range', 'regex', 'cookie', 'deviceCharacteristics', 'clientip', 'continent', 'countrycode', 'regioncode', 'protocol', 'method', 'proxy'\], got invalid_match_type`,
		},
		"match criteria AS - ObjectMatchValueRangeSubtype with incorrect value": {
			configPath: "testdata/TestDataCloudletsAudienceSegmentationMatchRule/invalid_type_range.tf",
			withError:  `cannot parse range_start value as an integer: strconv.ParseInt: parsing "range_start": invalid syntax`,
		},
	}
	for name, test := range tests {
		t.Run(name, func(t *testing.T) {
			resource.UnitTest(t, resource.TestCase{
				Providers: testAccProviders,
				Steps: []resource.TestStep{
					{
						Config: loadFixtureString(test.configPath),
						Check: resource.ComposeAggregateTestCheckFunc(
							resource.TestCheckResourceAttr(
								"data.akamai_cloudlets_audience_segmentation_match_rule.test", "match_rules.0.type", "asMatchRule"),
						),
						ExpectError: regexp.MustCompile(test.withError),
					},
				},
			})
		})
	}
}
