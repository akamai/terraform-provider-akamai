package property

import (
	"bufio"
	"bytes"
	"context"
	"encoding/json"
	"errors"
	"fmt"
	"github.com/akamai/terraform-provider-akamai/v2/pkg/tools"
	"github.com/hashicorp/go-cty/cty"
	"github.com/hashicorp/terraform-plugin-sdk/v2/diag"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/schema"
	"io/ioutil"
	"os"
	"path/filepath"
	"regexp"
	"strconv"
	"strings"
	"text/template"
)

func dataSourcePropertyRulesTemplate() *schema.Resource {
	return &schema.Resource{
		ReadContext: dataAkamaiPropertyRulesRead,
		Schema: map[string]*schema.Schema{
			"template_file": {
				Type:     schema.TypeString,
				Required: true,
			},
			"variables": {
				Type: schema.TypeSet,
				Elem: &schema.Resource{
					Schema: map[string]*schema.Schema{
						"name": {
							Type:     schema.TypeString,
							Required: true,
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
				RequiredWith:  []string{"var_values_file"},
			},
			"var_values_file": {
				Type:          schema.TypeString,
				Optional:      true,
				ConflictsWith: []string{"variables"},
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

func dataAkamaiPropertyRulesRead(ctx context.Context, d *schema.ResourceData, m interface{}) diag.Diagnostics {
	file, err := tools.GetStringValue("template_file", d)
	if err != nil {
		return diag.FromErr(err)
	}
	varsMap := make(map[string]interface{})
	vars, err := tools.GetSetValue("variables", d)
	if err != nil && !errors.Is(err, tools.ErrNotFound) {
		return diag.FromErr(err)
	}
	if err == nil {
		varsMap, err = convertToTypedMap(vars)
	}
	varsDefinitionFile, err := tools.GetStringValue("var_definition_file", d)
	if err != nil && !errors.Is(err, tools.ErrNotFound) {
		return diag.FromErr(err)
	}
	if err == nil {
		varsValuesFile, err := tools.GetStringValue("var_values_file", d)
		if err != nil && !errors.Is(err, tools.ErrNotFound) {
			return diag.FromErr(err)
		}
		varsMap, err = getVarsFromFile(varsDefinitionFile, varsValuesFile)
	}
	templateStr, err := convertToTemplate(file)
	if err != nil {
		return diag.FromErr(err)
	}
	tmpl, err := template.New("main").Delims(leftDelim, rightDelim).Option("missingkey=error").Parse(templateStr)
	if err != nil {
		return diag.FromErr(err)
	}
	dir := filepath.Dir(file)
	templateFiles := make(map[string]string)
	err = filepath.Walk(dir,
		func(path string, info os.FileInfo, err error) error {
			if err != nil {
				return err
			}
			if !info.IsDir() && path != file {
				templateFiles[strings.TrimPrefix(path, fmt.Sprintf("%s/", dir))] = path
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
	d.SetId(file)
	formatted := bytes.Buffer{}
	err = json.Indent(&formatted, wr.Bytes(), "", "  ")
	if err != nil {
		return diag.FromErr(err)
	}
	if err := d.Set("json", formatted.String()); err != nil {
		return diag.Errorf("%v: %s", tools.ErrValueSet, err.Error())
	}
	return nil
}

var (
	includeRegexp = regexp.MustCompile(`"#include:.+"`)
	varRegexp     = regexp.MustCompile(`"\${.+}"`)
)

func convertToTemplate(path string) (string, error) {
	builder := strings.Builder{}
	f, err := os.Open(path)
	if err != nil {
		return "", err
	}
	defer f.Close()
	scanner := bufio.NewScanner(f)
	for scanner.Scan() {
		line := scanner.Bytes()
		if includeStatement := includeRegexp.Find(line); len(includeStatement) > 0 {
			templateName := bytes.TrimPrefix(bytes.TrimSuffix(includeStatement, []byte(`"`)), []byte(`"#include:`))
			line = includeRegexp.ReplaceAll(line, []byte(fmt.Sprintf(`%stemplate "%s" .%s`, leftDelim, templateName, rightDelim)))
		}
		if varStatement := varRegexp.Find(line); len(varStatement) > 0 {
			varName := bytes.TrimPrefix(bytes.TrimSuffix(varStatement, []byte(`}"`)), []byte(`"${env`))
			line = varRegexp.ReplaceAll(line, []byte(fmt.Sprintf("%s%s%s", leftDelim, varName, rightDelim)))
		}
		builder.Write(line)
	}
	return builder.String(), nil
}

func convertToTypedMap(vars *schema.Set) (map[string]interface{}, error) {
	result := make(map[string]interface{})
	for _, variable := range vars.List() {
		varInfo, ok := variable.(map[string]interface{})
		if !ok {
			return nil, fmt.Errorf("unable to convert map entry to data object: %v", variable)
		}
		varName, ok := varInfo["name"]
		if !ok {
			return nil, fmt.Errorf("'name' argument is required in variable definition")
		}
		varNameStr, ok := varName.(string)
		if !ok {
			return nil, fmt.Errorf("'name' argument should be a string: %v", varName)
		}
		varType, ok := varInfo["type"]
		if !ok {
			return nil, fmt.Errorf("'type' argument is required in variable definition: %s", varNameStr)
		}
		varTypeStr, ok := varType.(string)
		if !ok {
			return nil, fmt.Errorf("'type' argument should be a string: %s", varNameStr)
		}
		value, ok := varInfo["value"]
		if !ok {
			return nil, fmt.Errorf("'value' argument is required in variable definition: %s", varNameStr)
		}
		valueStr, ok := value.(string)
		if !ok {
			return nil, fmt.Errorf("'value' argument should be a string: %s", varNameStr)
		}
		switch varTypeStr {
		case "string":
			result[varNameStr] = fmt.Sprintf(`"%s"`, valueStr)
		case "jsonBlock":
			if ok := json.Valid([]byte(valueStr)); !ok {
				return nil, fmt.Errorf("'jsonBlock` argument is not a valid json object: %s: %s", varNameStr, valueStr)
			}
			result[varNameStr] = valueStr
		case "number":
			num, err := strconv.ParseFloat(valueStr, 64)
			if err != nil {
				return nil, fmt.Errorf("value could not be represented as number: %w", err)
			}
			result[varNameStr] = num
		case "bool":
			boolean, err := strconv.ParseBool(valueStr)
			if err != nil {
				return nil, fmt.Errorf("value could not be represented as boolean: %w", err)
			}
			result[varNameStr] = boolean
		default:
			return nil, fmt.Errorf("unknown 'type' argument value: %s", varTypeStr)
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
		return nil, err
	}
	var definitions variableDefinitions
	if err := json.Unmarshal(definitionsFile, &definitions); err != nil {
		return nil, err
	}
	var values map[string]interface{}
	valuesFile, err := ioutil.ReadFile(valuesPath)
	if err != nil {
		return nil, err
	}
	if err := json.Unmarshal(valuesFile, &values); err != nil {
		return nil, err
	}
	vars := make(map[string]interface{})
	for name, varDef := range definitions.Definitions {
		switch v := varDef.Default.(type) {
		case string:
			vars[name] = fmt.Sprintf(`"%s"`, v)
		case map[string]interface{}:
			jsonBlock, err := json.Marshal(v)
			if err != nil {
				return nil, err
			}
			vars[name] = string(jsonBlock)
		default:
			vars[name] = v
		}
	}
	for name, value := range values {
		if _, ok := vars[name]; ok {
			switch v := value.(type) {
			case string:
				vars[name] = fmt.Sprintf(`"%s"`, v)
			case map[string]interface{}:
				jsonBlock, err := json.Marshal(v)
				if err != nil {
					return nil, err
				}
				vars[name] = string(jsonBlock)
			default:
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
