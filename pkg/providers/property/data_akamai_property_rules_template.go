package property

import (
	"bytes"
	"context"
	"crypto/sha1"
	"encoding/hex"
	"encoding/json"
	"errors"
	"fmt"
	"io/ioutil"
	"os"
	"path/filepath"
	"regexp"
	"strconv"
	"strings"
	"text/template"

	"github.com/akamai/terraform-provider-akamai/v6/pkg/common/tf"
	"github.com/akamai/terraform-provider-akamai/v6/pkg/meta"
	"github.com/apex/log"
	"github.com/hashicorp/go-cty/cty"
	"github.com/hashicorp/terraform-plugin-sdk/v2/diag"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/schema"
)

func dataSourcePropertyRulesTemplate() *schema.Resource {
	return &schema.Resource{
		ReadContext: dataPropertyRulesTemplateRead,
		Schema: map[string]*schema.Schema{
			"template_file": {
				Type:         schema.TypeString,
				Optional:     true,
				ExactlyOneOf: []string{"template", "template_file"},
				Description:  "File path to the template file",
			},
			"template": {
				Type: schema.TypeSet,
				Elem: &schema.Resource{
					Schema: map[string]*schema.Schema{
						"template_data": {
							Type:             schema.TypeString,
							Required:         true,
							ValidateDiagFunc: tf.IsNotBlank,
							Description:      "Content of the template as string",
						},
						"template_dir": {
							Type:             schema.TypeString,
							Required:         true,
							ValidateDiagFunc: tf.IsNotBlank,
							Description:      "Directory points to a folder, which contains snippets to include into template.",
						},
					},
				},
				Optional: true,
			},
			"variables": {
				Type: schema.TypeSet,
				Elem: &schema.Resource{
					Schema: map[string]*schema.Schema{
						"name": {
							Type:             schema.TypeString,
							Required:         true,
							ValidateDiagFunc: tf.IsNotBlank,
						},
						"type": {
							Type:     schema.TypeString,
							Default:  "string",
							Optional: true,
							ValidateDiagFunc: func(i interface{}, path cty.Path) diag.Diagnostics {
								val, ok := i.(string)
								if !ok {
									return diag.Errorf("value is not a string: %v", i)
								}
								switch val {
								case "bool", "number", "string", "jsonBlock":
									return nil
								}
								return diag.Errorf("'type' has invalid value: should be 'bool', 'number', 'string' or 'jsonBlock'")
							},
						},
						"value": {
							Type:     schema.TypeString,
							Required: true,
						},
					},
				},
				Optional:      true,
				ConflictsWith: []string{"var_definition_file", "var_values_file"},
			},
			"var_definition_file": {
				Type:          schema.TypeString,
				Optional:      true,
				ConflictsWith: []string{"variables"},
			},
			"var_values_file": {
				Type:          schema.TypeString,
				Optional:      true,
				ConflictsWith: []string{"variables"},
				RequiredWith:  []string{"var_definition_file"},
			},
			"json": {
				Type:     schema.TypeString,
				Computed: true,
			},
		},
	}
}

const (
	leftDelim  = "@+#"
	rightDelim = "#+@"
)

type variablePopulator struct {
	regex          *regexp.Regexp
	valueExtractor func(any) string
}

var (
	quotedVariablePopulator = variablePopulator{
		regex: regexp.MustCompile(`"\${env\.([^$}]+?)}"`),
		valueExtractor: func(input any) string {
			return fmt.Sprintf("%v", input)
		}}
	partialVariablePopulator = variablePopulator{
		regex: partialVariableRegexp,
		valueExtractor: func(input any) string {
			return strings.TrimSuffix(strings.TrimPrefix(fmt.Sprintf("%v", input), "\""), "\"")
		}}
)

func (v variablePopulator) hasMatch(template string) bool {
	return v.regex.MatchString(template)
}

func (v variablePopulator) replaceMatchWithVar(template string, varMap map[string]any) (string, error) {
	matchingVariable := v.regex.FindString(template)

	submatch := v.regex.FindStringSubmatch(template)
	if len(submatch) < 2 {
		return "", fmt.Errorf(matchingErrorMessage, matchingVariable)
	}

	varName := submatch[1]
	if varVal, ok := varMap[varName]; ok {
		return strings.ReplaceAll(template, matchingVariable, v.valueExtractor(varVal)), nil
	}
	return strings.ReplaceAll(template, matchingVariable, fmt.Sprintf("%s.%s%s", leftDelim, varName, rightDelim)), nil
}

//nolint:gocyclo
func dataPropertyRulesTemplateRead(_ context.Context, d *schema.ResourceData, m interface{}) diag.Diagnostics {
	meta := meta.Must(m)
	logger := meta.Log("PAPI", "dataPropertyRulesTemplateRead")

	file, err := tf.GetStringValue("template_file", d)
	if err != nil && !errors.Is(err, tf.ErrNotFound) {
		return diag.FromErr(err)
	}

	var dir string
	if err == nil {
		if _, err = os.Stat(file); err != nil {
			return diag.FromErr(err)
		}
		fileData, err := ioutil.ReadFile(file)
		if err != nil {
			return diag.FromErr(fmt.Errorf("%w: %s", ErrReadFile, err))
		}

		dir = filepath.Dir(file)
		if filepath.Ext(file) != ".json" || len(fileData) == 0 {
			logger.Errorf("snippets file should be with .json extension and cannot be empty: %s", file)
			return diag.Errorf("snippets file should be with .json extension and cannot be empty. Invalid file: %s ", file)
		}
	}

	var templateDataStr string
	if dir == "" {
		templateSet, err := tf.GetSetValue("template", d)
		if err != nil {
			return diag.FromErr(err)
		}
		templateDataStr, dir, err = flattenTemplate(templateSet.List())
		if err != nil {
			return diag.FromErr(err)
		}

		if _, err := os.Stat(dir); err != nil {
			return diag.FromErr(err)
		}
	}

	varsMap, err := getVariables(d, logger)
	if err != nil {
		return diag.FromErr(err)
	}

	var templateStr string
	if templateDataStr == "" {
		templateStr, err = convertToTemplate(file, varsMap)
	} else {
		templateStr, err = stringToTemplate(templateDataStr, varsMap, "main")
	}
	if err != nil {
		return diag.FromErr(err)
	}

	tmpl, err := template.New("main").Delims(leftDelim, rightDelim).Option("missingkey=error").Parse(templateStr)
	if err != nil {
		return diag.FromErr(err)
	}

	templateFiles := make(map[string]string)
	err = filepath.Walk(dir,
		func(path string, info os.FileInfo, err error) error {
			if err != nil {
				return err
			}

			pathDiff := strings.TrimPrefix(path, dir)
			if !info.IsDir() && path != file && !strings.Contains(pathDiff, ".terraform") {
				pathData, err := ioutil.ReadFile(path)
				if err != nil {
					return fmt.Errorf("%w: %s", ErrReadFile, err)
				}

				if len(pathData) > 0 {
					logger.Debugf("Template snippet found: %s", path)
					templateFiles[strings.TrimPrefix(filepath.ToSlash(path), fmt.Sprintf("%s/", filepath.ToSlash(dir)))] = path
				}
			}
			return nil
		})
	if err != nil {
		return diag.FromErr(err)
	}
	for name, f := range templateFiles {
		templateStr, err := convertToTemplate(f, varsMap)
		if err != nil {
			return diag.FromErr(err)
		}
		tmpl, err = tmpl.New(name).Delims(leftDelim, rightDelim).Option("missingkey=error").Parse(templateStr)
		if err != nil {
			return diag.FromErr(err)
		}
	}
	wr := bytes.Buffer{}
	err = tmpl.ExecuteTemplate(&wr, "main", varsMap)
	if err != nil {
		return diag.FromErr(err)
	}
	if file != "" && !jsonFileRegexp.MatchString(file) {
		return diag.Errorf("snippets file should have .json files. Invalid file %s ", file)
	}

	// Create a new SHA1 hash based on templateDataStr
	h := sha1.New()
	h.Write([]byte(templateDataStr))
	shaHash := hex.EncodeToString(h.Sum(nil))
	d.SetId(shaHash)

	formatted := bytes.Buffer{}
	result := wr.Bytes()
	err = json.Indent(&formatted, result, "", "  ")
	if err != nil {
		logger.Debugf("Creating rule tree resulted in invalid JSON: %s\nError: %s", result, err)
		return diag.FromErr(fmt.Errorf("invalid JSON result: %w", err))
	}
	if err := d.Set("json", formatted.String()); err != nil {
		return diag.Errorf("%v: %s", tf.ErrValueSet, err.Error())
	}
	return nil
}

func getVariables(d *schema.ResourceData, logger log.Interface) (map[string]interface{}, error) {
	varsMap := make(map[string]interface{})
	vars, err := tf.GetSetValue("variables", d)
	if err != nil && !errors.Is(err, tf.ErrNotFound) {
		return nil, err
	}
	if err == nil {
		varsMap, err = convertToTypedMap(vars.List())
		if err != nil {
			return nil, err
		}
	}
	varsDefinitionFile, err := tf.GetStringValue("var_definition_file", d)
	if err != nil && !errors.Is(err, tf.ErrNotFound) {
		return nil, err
	}
	if err == nil {
		logger.Debugf("Fetching variable definitions from file: %s", varsDefinitionFile)
		varsValuesFile, err := tf.GetStringValue("var_values_file", d)
		if err != nil && !errors.Is(err, tf.ErrNotFound) {
			return nil, err
		}
		varsMap, err = getVarsFromFile(varsDefinitionFile, varsValuesFile)
		if err != nil {
			return nil, err
		}
	}
	for k, v := range varsMap {
		if _, err := checkCircularDependency(fmt.Sprintf("%v", v), []string{k}, varsMap); err != nil {
			return nil, err
		}
	}
	return varsMap, nil
}

// that implementation has problem evaluating in template correctly following cases:
// "${env.a}" as `true` (for env.a is true of type bool)
// "${env.a} " as `"true "` (for env.a is "true" of type string)
// "${env.a}" as `"true"` (for env.a is true of type string)
// and hence should not be used for actual value evaluation
func checkCircularDependency(input string, seenVariables []string, varsMap map[string]interface{}) (string, error) {
	matchedVariable := partialVariableRegexp.FindString(input)
	for matchedVariable != "" {
		submatch := partialVariableRegexp.FindStringSubmatch(input)
		if len(submatch) < 2 {
			return "", fmt.Errorf(matchingErrorMessage, matchedVariable)
		}
		varName := submatch[1]
		for _, seenVariable := range seenVariables {
			if varName == seenVariable {
				return "", fmt.Errorf("hit cyclic dependency ending at %q", varName)
			}
		}
		varVal, ok := varsMap[varName]
		if !ok {
			return fmt.Sprintf("%s.%s%s", leftDelim, varName, rightDelim), nil
		}
		evaluatedVarVal, err := checkCircularDependency(fmt.Sprintf("%v", varVal), append(seenVariables, varName), varsMap)
		if err != nil {
			return "", err
		}
		input = strings.ReplaceAll(input, matchedVariable, evaluatedVarVal)
		matchedVariable = partialVariableRegexp.FindString(input)
	}
	return strings.TrimSuffix(strings.TrimPrefix(input, `"`), `"`), nil
}

var (
	includeRegexp         = regexp.MustCompile(`"#include:.+?"`)
	partialVariableRegexp = regexp.MustCompile(`\${env\.([^$}]+?)}`)
	jsonFileRegexp        = regexp.MustCompile(`\.json+$`)
)

var (
	// ErrReadFile is used to specify error while reading a file.
	ErrReadFile = errors.New("reading file")
	// ErrUnmarshal is used to specify unmarshal error.
	ErrUnmarshal = errors.New("unmarshaling value")
	// ErrFormatValue is used to specify formatting error.
	ErrFormatValue = errors.New("formatting value")
	// ErrUnknownType is used to specify unknown error.
	ErrUnknownType       = errors.New("unknown 'type' value")
	matchingErrorMessage = "there was a problem matching %q"
)

// flattenTemplate formats the template schema into a couple of strings holding template_data and template_dir values
func flattenTemplate(templateList []interface{}) (string, string, error) {
	if len(templateList) != 1 {
		return "", "", fmt.Errorf("%w: only single entry of template<template_data, template_dir> is allowed. Invalid template: %v ", tf.ErrInvalidType, templateList)
	}
	templateMap, ok := templateList[0].(map[string]interface{})
	if !ok {
		return "", "", fmt.Errorf("%w: unable to convert map entry to data object: %v", tf.ErrInvalidType, templateMap)
	}

	templateData, ok := templateMap["template_data"]
	if !ok {
		return "", "", fmt.Errorf("%w: 'template_data' argument is required in template definition", tf.ErrNotFound)
	}

	templateDataStr, ok := templateData.(string)
	if !ok {
		return "", "", fmt.Errorf("%w: 'template_data' argument should be a string: %v", tf.ErrInvalidType, templateData)
	}

	templateDir, ok := templateMap["template_dir"]
	if !ok {
		return "", "", fmt.Errorf("%w: 'template_dir' argument is required in template definition", tf.ErrNotFound)
	}

	templateDirStr, ok := templateDir.(string)
	if !ok {
		return "", "", fmt.Errorf("%w: 'template_dir' argument should be a string: %v", tf.ErrInvalidType, templateDir)
	}

	return templateDataStr, filepath.Clean(templateDirStr), nil
}

// stringToTemplate takes a large string (templateDataStr) and formats include/variable statements.
func stringToTemplate(templateDataStr string, varsMap map[string]interface{}, templatePath string) (string, error) {
	templateDataStr, err := evaluateVariables(templateDataStr, varsMap, templatePath)
	if err != nil {
		return "", err
	}

	includeStatement := includeRegexp.FindString(templateDataStr)
	for len(includeStatement) > 0 {
		templateName := strings.TrimPrefix(strings.TrimSuffix(includeStatement, `"`), `"#include:`)
		templateDataStr = strings.ReplaceAll(templateDataStr, includeStatement, fmt.Sprintf(`%stemplate "%s" .%s`, leftDelim, templateName, rightDelim))
		includeStatement = includeRegexp.FindString(templateDataStr)
	}

	if string(templateDataStr[len(templateDataStr)-1]) != "\n" {
		return fmt.Sprintf("%s\n", templateDataStr), nil
	}

	return templateDataStr, nil
}

func evaluateVariables(template string, varsMap map[string]interface{}, templatePath string) (string, error) {
	var err error

	for {
		for quotedVariablePopulator.hasMatch(template) {
			template, err = quotedVariablePopulator.replaceMatchWithVar(template, varsMap)
			if err != nil {
				return "", fmt.Errorf("error replacing variable at %q: %w", templatePath, err)
			}
		}

		if partialVariablePopulator.hasMatch(template) {
			template, err = partialVariablePopulator.replaceMatchWithVar(template, varsMap)
			if err != nil {
				return "", fmt.Errorf("error replacing variable at %q: %w", templatePath, err)
			}
		} else {
			break
		}
	}

	return template, nil
}

// convertToTemplate passes the string data to stringToTemplate after reading it from given path.
func convertToTemplate(path string, varsMap map[string]interface{}) (string, error) {
	b, err := ioutil.ReadFile(path)
	if err != nil {
		return "", fmt.Errorf("%w: %s", ErrReadFile, err)
	}

	return stringToTemplate(string(b), varsMap, path)
}

func convertToTypedMap(vars []interface{}) (map[string]interface{}, error) {
	result := make(map[string]interface{})
	for _, variable := range vars {
		varInfo, ok := variable.(map[string]interface{})
		if !ok {
			return nil, fmt.Errorf("%w: unable to convert map entry to data object: %v", tf.ErrInvalidType, variable)
		}
		varName, ok := varInfo["name"]
		if !ok {
			return nil, fmt.Errorf("%w: 'name' argument is required in variable definition", tf.ErrNotFound)
		}
		varNameStr, ok := varName.(string)
		if !ok {
			return nil, fmt.Errorf("%w: 'name' argument should be a string: %v", tf.ErrInvalidType, varName)
		}
		varType, ok := varInfo["type"]
		if !ok {
			return nil, fmt.Errorf("%w: 'type' argument is required in variable definition: %s", tf.ErrNotFound, varNameStr)
		}
		varTypeStr, ok := varType.(string)
		if !ok {
			return nil, fmt.Errorf("%w: 'type' argument should be a string: %s", tf.ErrInvalidType, varNameStr)
		}
		value, ok := varInfo["value"]
		if !ok {
			return nil, fmt.Errorf("%w: 'value' argument is required in variable definition: %s", tf.ErrNotFound, varNameStr)
		}
		valueStr, ok := value.(string)
		if !ok {
			return nil, fmt.Errorf("%w: 'value' argument should be a string: %s", tf.ErrInvalidType, varNameStr)
		}
		switch varTypeStr {
		case "string":
			result[varNameStr] = fmt.Sprintf(`"%s"`, valueStr)
		case "jsonBlock":
			var targetMap map[string]interface{}
			if err := json.Unmarshal([]byte(valueStr), &targetMap); err != nil {
				e := &json.UnmarshalTypeError{}
				if !errors.As(err, &e) {
					return nil, fmt.Errorf("%w: 'jsonBlock` argument is not a valid json object: %s: %s", ErrUnmarshal, varNameStr, valueStr)
				}
				var targetSlice []interface{}
				if err := json.Unmarshal([]byte(valueStr), &targetSlice); err != nil {
					return nil, fmt.Errorf("%w: 'jsonBlock` argument is not a valid json object: %s: %s", ErrUnmarshal, varNameStr, valueStr)
				}
			}
			result[varNameStr] = valueStr
		case "number":
			num, err := strconv.ParseFloat(valueStr, 64)
			if err != nil {
				return nil, fmt.Errorf("%w: value could not be represented as number: %s", tf.ErrInvalidType, err)
			}
			result[varNameStr] = num
		case "bool":
			boolean, err := strconv.ParseBool(valueStr)
			if err != nil {
				return nil, fmt.Errorf("%w: value could not be represented as boolean: %s", tf.ErrInvalidType, err)
			}
			result[varNameStr] = boolean
		default:
			return nil, fmt.Errorf("%w: %s", ErrUnknownType, varTypeStr)
		}
	}
	return result, nil
}

func getVarsFromFile(definitionsPath, valuesPath string) (map[string]interface{}, error) {
	type variableDefinitions struct {
		Definitions map[string]struct {
			Type    string      `json:"type"`
			Default interface{} `json:"default"`
		} `json:"definitions"`
	}
	definitionsFile, err := ioutil.ReadFile(definitionsPath)
	if err != nil {
		return nil, fmt.Errorf("%w: %s", ErrReadFile, err)
	}
	var definitions variableDefinitions
	if err := json.Unmarshal(definitionsFile, &definitions); err != nil {
		return nil, fmt.Errorf("%w: %s", ErrUnmarshal, err)
	}
	vars := make(map[string]interface{})
	for name, varDef := range definitions.Definitions {
		v, err := formatValue(varDef.Default)
		if err != nil {
			return nil, fmt.Errorf("%w: %s", ErrFormatValue, err)
		}
		vars[name] = v
	}
	if valuesPath != "" {
		var values map[string]interface{}
		valuesFile, err := ioutil.ReadFile(valuesPath)
		if err != nil {
			return nil, fmt.Errorf("%w: %s", ErrReadFile, err)
		}
		if err := json.Unmarshal(valuesFile, &values); err != nil {
			return nil, fmt.Errorf("%w: %s", ErrUnmarshal, err)
		}
		for name, value := range values {
			if _, ok := vars[name]; ok && value != nil {
				v, err := formatValue(value)
				if err != nil {
					return nil, fmt.Errorf("%w: %s", ErrFormatValue, err)
				}
				vars[name] = v
			}
		}
	}
	for k, v := range vars {
		if v == nil {
			vars[k] = "null"
		}
	}
	return vars, nil
}

func formatValue(val interface{}) (interface{}, error) {
	switch v := val.(type) {
	case string:
		return fmt.Sprintf(`"%s"`, v), nil
	case map[string]interface{}, []interface{}:
		jsonBlock, err := json.Marshal(v)
		if err != nil {
			return nil, err
		}
		return string(jsonBlock), nil
	default:
		return val, nil
	}
}
