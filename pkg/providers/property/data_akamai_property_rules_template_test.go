package property

import (
	"errors"
	"fmt"
	"regexp"
	"testing"

	"github.com/akamai/AkamaiOPEN-edgegrid-golang/v5/pkg/papi"
	"github.com/akamai/terraform-provider-akamai/v3/pkg/common/tf"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/resource"
	"github.com/stretchr/testify/require"
	"github.com/tj/assert"
)

func TestDataAkamaiPropertyRulesRead(t *testing.T) {
	t.Run("valid nested template with vars map", func(t *testing.T) {
		client := papi.Mock{}
		useClient(&client, nil, func() {
			resource.UnitTest(t, resource.TestCase{
				ProviderFactories: testAccProviders,
				Steps: []resource.TestStep{
					{
						Config: loadFixtureString("testdata/TestDSRulesTemplate/template_vars_map.tf"),
						Check: resource.ComposeAggregateTestCheckFunc(
							resource.TestCheckResourceAttr("data.akamai_property_rules_template.test", "json", loadFixtureString("testdata/TestDSRulesTemplate/rules/rules_out.json")),
						),
					},
					{
						Config: loadFixtureString("testdata/TestDSRulesTemplate/template_vars_map_with_data.tf"),
						Check: resource.ComposeAggregateTestCheckFunc(
							resource.TestCheckResourceAttr("data.akamai_property_rules_template.test", "json", loadFixtureString("testdata/TestDSRulesTemplate/rules/rules_out.json")),
						),
					},
				},
			})
		})
	})
	t.Run("valid nested template with vars map and non-standard property-snippets folder", func(t *testing.T) {
		client := papi.Mock{}
		useClient(&client, nil, func() {
			resource.UnitTest(t, resource.TestCase{
				ProviderFactories: testAccProviders,
				Steps: []resource.TestStep{
					{
						Config: loadFixtureString("testdata/TestDSRulesTemplate/template_vars_map_ns.tf"),
						Check: resource.ComposeAggregateTestCheckFunc(
							resource.TestCheckResourceAttr("data.akamai_property_rules_template.test", "json", loadFixtureString("testdata/TestDSRulesTemplate/rules/rules_out.json")),
						),
					},
					{
						Config: loadFixtureString("testdata/TestDSRulesTemplate/template_vars_map_with_data_ns.tf"),
						Check: resource.ComposeAggregateTestCheckFunc(
							resource.TestCheckResourceAttr("data.akamai_property_rules_template.test", "json", loadFixtureString("testdata/TestDSRulesTemplate/rules/rules_out.json")),
						),
					},
				},
			})
		})
	})
	t.Run("valid nested template with vars files", func(t *testing.T) {
		client := papi.Mock{}
		useClient(&client, nil, func() {
			resource.UnitTest(t, resource.TestCase{
				ProviderFactories: testAccProviders,
				Steps: []resource.TestStep{
					{
						Config: loadFixtureString("testdata/TestDSRulesTemplate/template_null_values.tf"),
						Check: resource.ComposeAggregateTestCheckFunc(
							resource.TestCheckResourceAttr("data.akamai_property_rules_template.test", "json", loadFixtureString("testdata/TestDSRulesTemplate/rules/rules_defaults.json")),
						),
					},
					{
						Config: loadFixtureString("testdata/TestDSRulesTemplate/template_null_values_with_data.tf"),
						Check: resource.ComposeAggregateTestCheckFunc(
							resource.TestCheckResourceAttr("data.akamai_property_rules_template.test", "json", loadFixtureString("testdata/TestDSRulesTemplate/rules/rules_defaults.json")),
						),
					},
				},
			})
		})
	})
	t.Run("null values do not overwrite defaults", func(t *testing.T) {
		client := papi.Mock{}
		useClient(&client, nil, func() {
			resource.UnitTest(t, resource.TestCase{
				ProviderFactories: testAccProviders,
				Steps: []resource.TestStep{
					{
						Config: loadFixtureString("testdata/TestDSRulesTemplate/template_vars_file.tf"),
						Check: resource.ComposeAggregateTestCheckFunc(
							resource.TestCheckResourceAttr("data.akamai_property_rules_template.test", "json", loadFixtureString("testdata/TestDSRulesTemplate/rules/rules_out.json")),
						),
					},
					{
						Config: loadFixtureString("testdata/TestDSRulesTemplate/template_vars_file_with_data.tf"),
						Check: resource.ComposeAggregateTestCheckFunc(
							resource.TestCheckResourceAttr("data.akamai_property_rules_template.test", "json", loadFixtureString("testdata/TestDSRulesTemplate/rules/rules_out.json")),
						),
					},
				},
			})
		})
	})
	t.Run("error conflicts in template_file and template", func(t *testing.T) {
		client := papi.Mock{}
		useClient(&client, nil, func() {
			resource.UnitTest(t, resource.TestCase{
				ProviderFactories: testAccProviders,
				Steps: []resource.TestStep{
					{
						Config:      loadFixtureString("testdata/TestDSRulesTemplate/template_file_data_conflict.tf"),
						ExpectError: regexp.MustCompile(`"template_file": only one of .template,template_file. can be specified`),
					},
					{
						Config:      loadFixtureString("testdata/TestDSRulesTemplate/template_file_data_missing.tf"),
						ExpectError: regexp.MustCompile(`"template_file": one of .template,template_file. must be specified`),
					},
				},
			})
		})
	})
	t.Run("error missing values in template", func(t *testing.T) {
		client := papi.Mock{}
		useClient(&client, nil, func() {
			resource.UnitTest(t, resource.TestCase{
				ProviderFactories: testAccProviders,
				Steps: []resource.TestStep{
					{
						Config:      loadFixtureString("testdata/TestDSRulesTemplate/template_missing_data.tf"),
						ExpectError: regexp.MustCompile(`The argument "template_data" is required, but no definition was found.`),
					},
					{
						Config:      loadFixtureString("testdata/TestDSRulesTemplate/template_missing_dir.tf"),
						ExpectError: regexp.MustCompile(`The argument "template_dir" is required, but no definition was found.`),
					},
				},
			})
		})
	})
	t.Run("error setting both map and file variables", func(t *testing.T) {
		client := papi.Mock{}
		useClient(&client, nil, func() {
			resource.UnitTest(t, resource.TestCase{
				ProviderFactories: testAccProviders,
				Steps: []resource.TestStep{
					{
						Config:      loadFixtureString("testdata/TestDSRulesTemplate/template_vars_conflict.tf"),
						ExpectError: regexp.MustCompile(`"variables": conflicts with var_definition_file`),
					},
				},
			})
		})
	})
	t.Run("variable has invalid type", func(t *testing.T) {
		client := papi.Mock{}
		useClient(&client, nil, func() {
			resource.UnitTest(t, resource.TestCase{
				ProviderFactories: testAccProviders,
				Steps: []resource.TestStep{
					{
						Config:      loadFixtureString("testdata/TestDSRulesTemplate/template_vars_invalid_type.tf"),
						ExpectError: regexp.MustCompile(`'type' has invalid value: should be 'bool', 'number', 'string' or 'jsonBlock'`),
					},
				},
			})
		})
	})
	t.Run("variable not found in template", func(t *testing.T) {
		client := papi.Mock{}
		useClient(&client, nil, func() {
			resource.UnitTest(t, resource.TestCase{
				ProviderFactories: testAccProviders,
				Steps: []resource.TestStep{
					{
						Config:      loadFixtureString("testdata/TestDSRulesTemplate/template_var_not_found.tf"),
						ExpectError: regexp.MustCompile(`executing "snippets/sub/another-template.json" at <.options>: map has no entry for key "options"`),
					},
				},
			})
		})
	})
	t.Run("invalid variable in map", func(t *testing.T) {
		client := papi.Mock{}
		useClient(&client, nil, func() {
			resource.UnitTest(t, resource.TestCase{
				ProviderFactories: testAccProviders,
				Steps: []resource.TestStep{
					{
						Config:      loadFixtureString("testdata/TestDSRulesTemplate/template_vars_invalid_value.tf"),
						ExpectError: regexp.MustCompile(`value could not be represented as number: strconv.ParseFloat: parsing "all": invalid syntax`),
					},
				},
			})
		})
	})
	t.Run("error fetching vars from map - file not found", func(t *testing.T) {
		client := papi.Mock{}
		useClient(&client, nil, func() {
			resource.UnitTest(t, resource.TestCase{
				ProviderFactories: testAccProviders,
				Steps: []resource.TestStep{
					{
						Config:      loadFixtureString("testdata/TestDSRulesTemplate/template_vars_file_not_found.tf"),
						ExpectError: regexp.MustCompile(`reading file: open invalid_path: no such file or directory`),
					},
				},
			})
		})
	})
	t.Run("invalid json result", func(t *testing.T) {
		client := papi.Mock{}
		useClient(&client, nil, func() {
			resource.UnitTest(t, resource.TestCase{
				ProviderFactories: testAccProviders,
				Steps: []resource.TestStep{
					{
						Config:      loadFixtureString("testdata/TestDSRulesTemplate/template_invalid_json.tf"),
						ExpectError: regexp.MustCompile(`snippets file should be with .json extension and valid json data. Invalid file: testdata/TestDSRulesTemplate/property-snippets/template_invalid_json.json`),
					},
				},
			})
		})
	})
	t.Run("template file not found", func(t *testing.T) {
		client := papi.Mock{}
		useClient(&client, nil, func() {
			resource.UnitTest(t, resource.TestCase{
				ProviderFactories: testAccProviders,
				Steps: []resource.TestStep{
					{
						Config:      loadFixtureString("testdata/TestDSRulesTemplate/template_file_not_found.tf"),
						ExpectError: regexp.MustCompile(`Error: stat testdata/TestDSRulesTemplate/rules/property-snippets/non-existent.json: no such file or directory`),
					},
				},
			})
		})
	})

	t.Run("template file is empty", func(t *testing.T) {
		client := papi.Mock{}
		useClient(&client, nil, func() {
			resource.UnitTest(t, resource.TestCase{
				ProviderFactories: testAccProviders,
				Steps: []resource.TestStep{
					{
						Config:      loadFixtureString("testdata/TestDSRulesTemplate/template_file_is_empty.tf"),
						ExpectError: regexp.MustCompile(`Error: snippets file should be with .json extension and valid json data. Invalid file: testdata/TestDSRulesTemplate/property-snippets/empty_json.json`),
					},
				},
			})
		})
	})

	t.Run("snippets files are under incorrect folder deeply nested", func(t *testing.T) {
		client := papi.Mock{}
		useClient(&client, nil, func() {
			resource.UnitTest(t, resource.TestCase{
				ProviderFactories: testAccProviders,
				Steps: []resource.TestStep{
					{
						Config:      loadFixtureString("testdata/TestDSRulesTemplate/template_invalid_snippets_folder_json.tf"),
						ExpectError: regexp.MustCompile(`Error: snippets file should be with .json extension and valid json data. Invalid file: testdata/TestDSRulesTemplate/output/template_invalid_json.json`),
					},
				},
			})
		})
	})

	t.Run("snippets files are under incorrect folder e.g. property-snippets/rules.json", func(t *testing.T) {
		client := papi.Mock{}
		useClient(&client, nil, func() {
			resource.UnitTest(t, resource.TestCase{
				ProviderFactories: testAccProviders,
				Steps: []resource.TestStep{
					{
						Config:      loadFixtureString("testdata/TestDSRulesTemplate/template_invalid_snippets_only_one_folder_json.tf"),
						ExpectError: regexp.MustCompile(`Error: snippets file should be with .json extension and valid json data. Invalid file: property-snippet/template_invalid_json.json`),
					},
				},
			})
		})
	})

}

func TestFormatValue(t *testing.T) {
	tests := map[string]struct {
		given     interface{}
		expected  interface{}
		withError bool
	}{
		"string": {
			given:    "test",
			expected: `"test"`,
		},
		"map": {
			given:    map[string]interface{}{"string": "value", "num": 1, "map": map[string]interface{}{"bool": true}},
			expected: `{"map":{"bool":true},"num":1,"string":"value"}`,
		},
		"number": {
			given:    1.23,
			expected: 1.23,
		},
		"boolean": {
			given:    true,
			expected: true,
		},
		"unmarshalable map": {
			given:     map[string]interface{}{"f": func() {}},
			withError: true,
		},
	}
	for name, test := range tests {
		t.Run(name, func(t *testing.T) {
			res, err := formatValue(test.given)
			if test.withError {
				assert.Error(t, err)
				return
			}
			require.NoError(t, err)
			assert.Equal(t, test.expected, res)
		})
	}
}

func TestGetValuesFromMap(t *testing.T) {
	variablesPath := "testdata/TestDSRulesTemplate/variables"
	tests := map[string]struct {
		definitionsFile string
		valuesFile      string
		expected        map[string]interface{}
		withError       error
	}{
		"definitions and values passed, with values overwrite": {
			definitionsFile: "simple_definitions.json",
			valuesFile:      "simple_values.json",
			expected: map[string]interface{}{
				"testString":    `"test 2"`,
				"testJSONMap":   `{"abc":"bca"}`,
				"testJSONArray": `["d","e","f"]`,
				"testNumber":    "null",
			},
		},
		"values not passed, take defaults": {
			definitionsFile: "simple_definitions.json",
			expected: map[string]interface{}{
				"testString":    `"test"`,
				"testJSONMap":   `{"abc":"abc"}`,
				"testJSONArray": `["a","b","c"]`,
				"testNumber":    "null",
			},
		},
		"definitions file not found": {
			definitionsFile: "not_existing.json",
			withError:       ErrReadFile,
		},
		"values file not found": {
			definitionsFile: "simple_definitions.json",
			valuesFile:      "not_existing.json",
			withError:       ErrReadFile,
		},
		"invalid definitions schema": {
			definitionsFile: "invalid_definition_schema.json",
			withError:       ErrUnmarshal,
		},
		"invalid values json": {
			definitionsFile: "simple_definitions.json",
			valuesFile:      "invalid_values.json",
			withError:       ErrUnmarshal,
		},
	}
	for name, test := range tests {
		t.Run(name, func(t *testing.T) {
			defPath := fmt.Sprintf("%s/%s", variablesPath, test.definitionsFile)
			var valPath string
			if test.valuesFile != "" {
				valPath = fmt.Sprintf("%s/%s", variablesPath, test.valuesFile)
			}
			res, err := getVarsFromFile(
				defPath,
				valPath)
			if test.withError != nil {
				assert.True(t, errors.Is(err, test.withError), "expected: %s; got: %s", test.expected, err)
				return
			}
			require.NoError(t, err)
			assert.Equal(t, test.expected, res)
		})
	}
}

func TestConvertToTypedMap(t *testing.T) {
	tests := map[string]struct {
		givenVars []interface{}
		expected  map[string]interface{}
		withError error
	}{
		"valid map": {
			givenVars: []interface{}{
				map[string]interface{}{"name": "testString", "type": "string", "value": "test"},
				map[string]interface{}{"name": "testNum", "type": "number", "value": "1.23"},
				map[string]interface{}{"name": "testJSONMap", "type": "jsonBlock", "value": `{"abc": "cba", "number":1}`},
				map[string]interface{}{"name": "testJSONArray", "type": "jsonBlock", "value": `["a", "b", "c"]`},
				map[string]interface{}{"name": "testBool", "type": "bool", "value": "true"},
			},
			expected: map[string]interface{}{
				"testString":    `"test"`,
				"testNum":       1.23,
				"testJSONMap":   `{"abc": "cba", "number":1}`,
				"testJSONArray": `["a", "b", "c"]`,
				"testBool":      true,
			},
		},
		"invalid values slice": {
			givenVars: []interface{}{"test"},
			withError: tf.ErrInvalidType,
		},
		"missing 'name' in map": {
			givenVars: []interface{}{
				map[string]interface{}{"type": "string", "value": "test"},
			},
			withError: tf.ErrNotFound,
		},
		"missing 'type' in map": {
			givenVars: []interface{}{
				map[string]interface{}{"name": "testString", "value": "test"},
			},
			withError: tf.ErrNotFound,
		},
		"missing 'value' in map": {
			givenVars: []interface{}{
				map[string]interface{}{"type": "string", "name": "test"},
			},
			withError: tf.ErrNotFound,
		},
		"'name' is of invalid type": {
			givenVars: []interface{}{
				map[string]interface{}{"name": 123, "type": "string", "value": "test"},
			},
			withError: tf.ErrInvalidType,
		},
		"'type' is of invalid type": {
			givenVars: []interface{}{
				map[string]interface{}{"name": "test", "type": 123, "value": "test"},
			},
			withError: tf.ErrInvalidType,
		},
		"'value' is of invalid type": {
			givenVars: []interface{}{
				map[string]interface{}{"name": "test", "type": "string", "value": 123},
			},
			withError: tf.ErrInvalidType,
		},
		"jsonBlock has invalid json": {
			givenVars: []interface{}{
				map[string]interface{}{"name": "testJSON", "type": "jsonBlock", "value": "abc"},
			},
			withError: ErrUnmarshal,
		},
		"number is invalid": {
			givenVars: []interface{}{
				map[string]interface{}{"name": "test", "type": "number", "value": "abc"},
			},
			withError: tf.ErrInvalidType,
		},
		"bool has invalid json": {
			givenVars: []interface{}{
				map[string]interface{}{"name": "test", "type": "bool", "value": "abc"},
			},
			withError: tf.ErrInvalidType,
		},
		"unknown type passed": {
			givenVars: []interface{}{
				map[string]interface{}{"name": "test", "type": "unknown", "value": "abc"},
			},
			withError: ErrUnknownType,
		},
	}

	for name, test := range tests {
		t.Run(name, func(t *testing.T) {
			res, err := convertToTypedMap(test.givenVars)
			if test.withError != nil {
				assert.True(t, errors.Is(err, test.withError), "want: %s; got: %s", test.withError, err)
				return
			}
			require.NoError(t, err)
			assert.Equal(t, test.expected, res)
		})
	}
}

func TestFlattenTemplate(t *testing.T) {
	tests := map[string]struct {
		givenList    []interface{}
		expectedData string
		expectedDir  string
		withError    error
	}{
		"valid list": {
			givenList: []interface{}{
				map[string]interface{}{
					"template_data": loadFixtureString("testdata/TestDSRulesTemplate/rules/property-snippets/plain_json.json"),
					"template_dir":  "testdata/TestDSRulesTemplate/rules/property-snippets/",
				},
			},
			expectedData: loadFixtureString("testdata/TestDSRulesTemplate/rules/property-snippets/plain_json.json"),
			expectedDir:  "testdata/TestDSRulesTemplate/rules/property-snippets",
		},
		"invalid list length": {
			givenList: []interface{}{
				map[string]interface{}{
					"template_data": loadFixtureString("testdata/TestDSRulesTemplate/rules/property-snippets/plain_json.json"),
					"template_dir":  "testdata/TestDSRulesTemplate/rules/property-snippets/",
				},
				map[string]interface{}{
					"template_data": loadFixtureString("testdata/TestDSRulesTemplate/rules/property-snippets/template_in.json"),
					"template_dir":  "testdata/TestDSRulesTemplate/rules/property-snippets",
				},
			},
			withError: tf.ErrInvalidType,
		},
		"missing 'template_data' in list": {
			givenList: []interface{}{
				map[string]interface{}{
					"template_dir": "testdata/TestDSRulesTemplate/rules/property-snippets/",
				},
			},
			withError: tf.ErrNotFound,
		},
		"missing 'template_dir' in list": {
			givenList: []interface{}{
				map[string]interface{}{
					"template_data": loadFixtureString("testdata/TestDSRulesTemplate/rules/property-snippets/template_in.json"),
				},
			},
			withError: tf.ErrNotFound,
		},
		"invalid 'template_data' in list": {
			givenList: []interface{}{
				map[string]interface{}{
					"template_data": 123,
					"template_dir":  "testdata/TestDSRulesTemplate/rules/property-snippets/",
				},
			},
			withError: tf.ErrInvalidType,
		},
		"invalid 'template_dir' in list": {
			givenList: []interface{}{
				map[string]interface{}{
					"template_data": loadFixtureString("testdata/TestDSRulesTemplate/rules/property-snippets/template_in.json"),
					"template_dir":  true,
				},
			},
			withError: tf.ErrInvalidType,
		},
	}

	for name, test := range tests {
		t.Run(name, func(t *testing.T) {
			resData, resDir, err := flattenTemplate(test.givenList)
			if test.withError != nil {
				assert.True(t, errors.Is(err, test.withError), "want: %s; got: %s", test.withError, err)
				return
			}
			require.NoError(t, err)
			assert.Equal(t, test.expectedData, resData)
			assert.Equal(t, test.expectedDir, resDir)
		})
	}
}

func TestConvertToTemplate(t *testing.T) {
	templates := "testdata/TestDSRulesTemplate/rules/property-snippets"
	templatesOut := "testdata/TestDSRulesTemplate/output"
	tests := map[string]struct {
		givenFile    string
		expectedFile string
		withError    error
	}{
		"valid conversion": {
			givenFile:    "template_in.json",
			expectedFile: "template_out.json",
		},
		"plain JSON passed": {
			givenFile:    "plain_json.json",
			expectedFile: "plain_json.json",
		},
		"invalid path": {
			givenFile: "invalid.json",
			withError: ErrReadFile,
		},
		"multiple vars": {
			givenFile:    "/snippets/some-template.json",
			expectedFile: "/property-snippets/some-template-out.json",
		},
	}
	for name, test := range tests {
		t.Run(name, func(t *testing.T) {
			res, err := convertToTemplate(fmt.Sprintf("%s/%s", templates, test.givenFile))
			if test.withError != nil {
				assert.True(t, errors.Is(err, test.withError), "want: %s; got: %s", test.withError, err)
				return
			}
			require.NoError(t, err)
			expected := loadFixtureString(fmt.Sprintf("%s/%s", templatesOut, test.expectedFile))
			assert.Equal(t, expected, res)
		})
	}
}

func TestStringToTemplate(t *testing.T) {
	templates := "testdata/TestDSRulesTemplate/rules/property-snippets"
	templatesOut := "testdata/TestDSRulesTemplate/output"
	tests := map[string]struct {
		givenFile    string
		expectedFile string
		withError    error
	}{
		"valid conversion": {
			givenFile:    "template_in.json",
			expectedFile: "template_out.json",
		},
		"multiple includes in array": {
			givenFile:    "template_in_with_array.json",
			expectedFile: "template_out_with_array.json",
		},
		"plain JSON passed": {
			givenFile:    "plain_json.json",
			expectedFile: "plain_json.json",
		},
	}
	for name, test := range tests {
		t.Run(name, func(t *testing.T) {
			givenString := loadFixtureString(fmt.Sprintf("%s/%s", templates, test.givenFile))
			res, err := stringToTemplate(givenString)
			fmt.Println(res)
			if test.withError != nil {
				assert.True(t, errors.Is(err, test.withError), "want: %s; got: %s", test.withError, err)
				return
			}
			require.NoError(t, err)
			expected := loadFixtureString(fmt.Sprintf("%s/%s", templatesOut, test.expectedFile))
			assert.Equal(t, expected, res)
		})
	}
}
