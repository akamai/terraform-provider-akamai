package property

import (
	"bytes"
	"context"
	"encoding/json"
	"errors"
	"fmt"
	"github.com/akamai/terraform-provider-akamai/v2/pkg/tools"
	"github.com/hashicorp/terraform-plugin-sdk/v2/diag"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/schema"
	"html/template"
	"io/ioutil"
	"path/filepath"
	"strings"
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
				Type:     schema.TypeMap,
				Elem:     &schema.Schema{Type: schema.TypeString},
				Optional: true,
			},
			"template_dir": {
				Type:     schema.TypeString,
				Optional: true,
				Default:  "snippets",
			},
			"json": {
				Type:     schema.TypeString,
				Computed: true,
			},
		},
	}
}

func dataAkamaiPropertyRulesRead(ctx context.Context, d *schema.ResourceData, m interface{}) diag.Diagnostics {
	file, err := tools.GetStringValue("template_file", d)
	if err != nil {
		return diag.FromErr(err)
	}
	vars, ok := d.Get("variables").(map[string]interface{})
	if !ok {
		return diag.FromErr(fmt.Errorf("%w: %s, %q", tools.ErrInvalidType, "variables", "map[string]interface{}"))
	}
	templateDir, err := tools.GetStringValue("template_dir", d)
	if err != nil {
		if !errors.Is(err, tools.ErrNotFound) {
			return diag.FromErr(err)
		}
		return nil
	}
	dir := filepath.Dir(file)
	allFiles := []string{file}
	files, err := ioutil.ReadDir(fmt.Sprintf("%s/%s", dir, templateDir))
	if err != nil {
		return diag.FromErr(err)
	}
	for _, sub := range files {
		filename := sub.Name()
		if strings.HasSuffix(filename, ".tmpl") {
			allFiles = append(allFiles, fmt.Sprintf("%s/%s/%s", dir, templateDir, filename))
		}
	}
	templates, err := template.ParseFiles(allFiles...)
	if err != nil {
		return diag.FromErr(err)
	}
	_, mainTemplateName := filepath.Split(file)
	mainTemplate := templates.Lookup(mainTemplateName)
	wr := bytes.Buffer{}
	err = mainTemplate.Execute(&wr, vars)
	if err != nil {
		return diag.FromErr(err)
	}
	d.SetId(mainTemplateName)
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
