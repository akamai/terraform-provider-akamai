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

	"github.com/hashicorp/go-cty/cty"
	"github.com/hashicorp/terraform-plugin-sdk/v2/diag"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/schema"

	"github.com/akamai/terraform-provider-akamai/v2/pkg/akamai"
	"github.com/akamai/terraform-provider-akamai/v2/pkg/tools"
)

func dataSourcePropertyRulesTemplate() *schema.Resource {
	return &schema.Resource{
		ReadContext: dataAkamaiPropertyRulesRead,
		Schema: map[string]*schema.Schema{
			"template_file": {
				Type:         schema.TypeString,
				Optional:     true,
				ExactlyOneOf: []string{"template", "template_file"},
				Description:  "File path to the template file inside 'property-snippets' subfolder",
			},
			"template": {
				Type: schema.TypeSet,
				Elem: &schema.Resource{
					Schema: map[string]*schema.Schema{
						"template_data": {
							Type:             schema.TypeString,
							Required:         true,
							ValidateDiagFunc: tools.IsNotBlank,
							Description:      "Content of the template as string",
						},
						"template_dir": {
							Type:             schema.TypeString,
							Required:         true,
							ValidateDiagFunc: tools.IsNotBlank,
							Description:      "Directory points to a folder ending with 'property-snippets', which contains snippets to include into template.",
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
							ValidateDiagFunc: tools.IsNotBlank,
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

func dataAkamaiPropertyRulesRead(_ context.Context, d *schema.ResourceData, m interface{}) diag.Diagnostics {
	meta := akamai.Meta(m)
	logger := meta.Log("PAPI", "dataAkamaiPropertyRulesRead")

	file, err := tools.GetStringValue("template_file", d)
	if err != nil && !errors.Is(err, tools.ErrNotFound) {
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
		if filepath.Base(dir) != "property-snippets" || filepath.Ext(file) != ".json" || !json.Valid(fileData) {
			logger.Errorf("snippets file should be under 'property-snippets' folder with .json extension and valid json data: %s", file)
			return diag.FromErr(fmt.Errorf("snippets file should be under 'property-snippets' folder with .json extension and valid json data. Invalid file: %s ", file))
		}
	}

	var templateDataStr string
	if dir == "" {
		templateSet, err := tools.GetSetValue("template", d)
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

		if filepath.Base(dir) != "property-snippets" {
			logger.Errorf("'template_dir' should points to 'property-snippets' folder: %s", dir)
			return diag.Errorf("'template_dir' should points to 'property-snippets' folder: %s", dir)
		}
	}

	var templateStr string
	if templateDataStr == "" {
		templateStr, err = convertToTemplate(file)
		if err != nil {
			return diag.FromErr(err)
		}
	} else {
		templateStr, err = stringToTemplate(templateDataStr)
		if err != nil {
			return diag.FromErr(err)
		}
	}

	tmpl, err := template.New("main").Delims(leftDelim, rightDelim).Option("missingkey=error").Parse(templateStr)
	if err != nil {
		return diag.FromErr(err)
	}

	varsMap := make(map[string]interface{})
	vars, err := tools.GetSetValue("variables", d)
	if err != nil && !errors.Is(err, tools.ErrNotFound) {
		return diag.FromErr(err)
	}
	if err == nil {
		varsMap, err = convertToTypedMap(vars.List())
		if err != nil {
			return diag.FromErr(err)
		}
	}
	varsDefinitionFile, err := tools.GetStringValue("var_definition_file", d)
	if err != nil && !errors.Is(err, tools.ErrNotFound) {
		return diag.FromErr(err)
	}
	if err == nil {
		logger.Debugf("Fetching variable definitions from file: %s", varsDefinitionFile)
		varsValuesFile, err := tools.GetStringValue("var_values_file", d)
		if err != nil && !errors.Is(err, tools.ErrNotFound) {
			return diag.FromErr(err)
		}
		varsMap, err = getVarsFromFile(varsDefinitionFile, varsValuesFile)
		if err != nil {
			return diag.FromErr(err)
		}
	}
	templateFiles := make(map[string]string)
	err = filepath.Walk(dir,
		func(path string, info os.FileInfo, err error) error {
			if err != nil {
				return err
			}

			if !info.IsDir() && path != file {
				pathData, err := ioutil.ReadFile(path)
				if err != nil {
					return fmt.Errorf("%w: %s", ErrReadFile, err)
				}

				if json.Valid(pathData) {
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
		templateStr, err := convertToTemplate(f)
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
		return diag.Errorf("snippets file under 'property-snippets' folder should have .json files. Invalid file %s ", file)
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
		return diag.Errorf("%v: %s", tools.ErrValueSet, err.Error())
	}
	return nil
}

var (
	includeRegexp  = regexp.MustCompile(`"#include:.+"`)
	varRegexp      = regexp.MustCompile(`"\${.+}"`)
	jsonFileRegexp = regexp.MustCompile(`\.json+$`)
)

var (
	// ErrReadFile is used to specify error while reading a file.
	ErrReadFile = errors.New("reading file")
	// ErrUnmarshal is used to specify unmarshal error.
	ErrUnmarshal = errors.New("unmarshaling value")
	// ErrFormatValue is used to specify formatting error.
	ErrFormatValue = errors.New("formatting value")
	// ErrUnknownType is used to specify unknown error.
	ErrUnknownType = errors.New("unknown 'type' value")
)

// flattenTemplate formats the template schema into a couple of strings holding template_data and template_dir values
func flattenTemplate(templateList []interface{}) (string, string, error) {
	if len(templateList) != 1 {
		return "", "", fmt.Errorf("%w: only single entry of template<template_data, template_dir> is allowed. Invalid template: %v ", tools.ErrInvalidType, templateList)
	}
	templateMap, ok := templateList[0].(map[string]interface{})
	if !ok {
		return "", "", fmt.Errorf("%w: unable to convert map entry to data object: %v", tools.ErrInvalidType, templateMap)
	}

	templateData, ok := templateMap["template_data"]
	if !ok {
		return "", "", fmt.Errorf("%w: 'template_data' argument is required in template definition", tools.ErrNotFound)
	}

	templateDataStr, ok := templateData.(string)
	if !ok {
		return "", "", fmt.Errorf("%w: 'template_data' argument should be a string: %v", tools.ErrInvalidType, templateData)
	}

	templateDir, ok := templateMap["template_dir"]
	if !ok {
		return "", "", fmt.Errorf("%w: 'template_dir' argument is required in template definition", tools.ErrNotFound)
	}

	templateDirStr, ok := templateDir.(string)
	if !ok {
		return "", "", fmt.Errorf("%w: 'template_dir' argument should be a string: %v", tools.ErrInvalidType, templateDir)
	}

	return templateDataStr, filepath.Clean(templateDirStr), nil
}

// stringToTemplate takes a large string (templateDataStr) and formats include/variable statements.
func stringToTemplate(templateDataStr string) (string, error) {
	includeStatement := includeRegexp.FindString(templateDataStr)
	for len(includeStatement) > 0 {
		templateName := strings.TrimPrefix(strings.TrimSuffix(includeStatement, `"`), `"#include:`)
		templateDataStr = strings.ReplaceAll(templateDataStr, includeStatement, fmt.Sprintf(`%stemplate "%s" .%s`, leftDelim, templateName, rightDelim))
		includeStatement = includeRegexp.FindString(templateDataStr)
	}

	varStatement := varRegexp.FindString(templateDataStr)
	for len(varStatement) > 0 {
		varName := strings.TrimPrefix(strings.TrimSuffix(varStatement, `}"`), `"${env`)
		templateDataStr = strings.ReplaceAll(templateDataStr, varStatement, fmt.Sprintf("%s%s%s", leftDelim, varName, rightDelim))
		varStatement = varRegexp.FindString(templateDataStr)
	}

	if string(templateDataStr[len(templateDataStr)-1]) != "\n" {
		return fmt.Sprintf("%s\n", templateDataStr), nil
	}

	return templateDataStr, nil
}

// convertToTemplate passes the string data to stringToTemplate after reading it from given path.
func convertToTemplate(path string) (string, error) {
	b, err := ioutil.ReadFile(path)
	if err != nil {
		return "", fmt.Errorf("%w: %s", ErrReadFile, err)
	}

	return stringToTemplate(string(b))
}

func convertToTypedMap(vars []interface{}) (map[string]interface{}, error) {
	result := make(map[string]interface{})
	for _, variable := range vars {
		varInfo, ok := variable.(map[string]interface{})
		if !ok {
			return nil, fmt.Errorf("%w: unable to convert map entry to data object: %v", tools.ErrInvalidType, variable)
		}
		varName, ok := varInfo["name"]
		if !ok {
			return nil, fmt.Errorf("%w: 'name' argument is required in variable definition", tools.ErrNotFound)
		}
		varNameStr, ok := varName.(string)
		if !ok {
			return nil, fmt.Errorf("%w: 'name' argument should be a string: %v", tools.ErrInvalidType, varName)
		}
		varType, ok := varInfo["type"]
		if !ok {
			return nil, fmt.Errorf("%w: 'type' argument is required in variable definition: %s", tools.ErrNotFound, varNameStr)
		}
		varTypeStr, ok := varType.(string)
		if !ok {
			return nil, fmt.Errorf("%w: 'type' argument should be a string: %s", tools.ErrInvalidType, varNameStr)
		}
		value, ok := varInfo["value"]
		if !ok {
			return nil, fmt.Errorf("%w: 'value' argument is required in variable definition: %s", tools.ErrNotFound, varNameStr)
		}
		valueStr, ok := value.(string)
		if !ok {
			return nil, fmt.Errorf("%w: 'value' argument should be a string: %s", tools.ErrInvalidType, varNameStr)
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
				return nil, fmt.Errorf("%w: value could not be represented as number: %s", tools.ErrInvalidType, err)
			}
			result[varNameStr] = num
		case "bool":
			boolean, err := strconv.ParseBool(valueStr)
			if err != nil {
				return nil, fmt.Errorf("%w: value could not be represented as boolean: %s", tools.ErrInvalidType, err)
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
