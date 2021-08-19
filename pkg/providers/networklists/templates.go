package networklists

import (
	"bytes"
	"encoding/json"
	"fmt"
	"strconv"
	"strings"
	"text/template"

	"github.com/jedib0t/go-pretty/v6/table"
)

// OutputTemplates is a map of templates
type OutputTemplates map[string]*OutputTemplate

// OutputTemplate contains template data
type OutputTemplate struct {
	TemplateName   string
	TemplateType   string
	TableTitle     string
	TemplateString string
}

// GetTemplate given map of templates and a key, returns template stored under this key
func GetTemplate(ots map[string]*OutputTemplate, key string) (*OutputTemplate, error) {
	if f, ok := ots[key]; ok && f != nil {
		return f, nil
	}
	return nil, fmt.Errorf("Error not found")

}

// RenderTemplates renders template and returns it as a string
func RenderTemplates(ots map[string]*OutputTemplate, key string, str interface{}) (string, error) {
	var ostr, tstr bytes.Buffer
	templ, err := GetTemplate(ots, key)

	if err != nil {
		return "", nil
	}
	var (
		funcs = template.FuncMap{
			"join":  strings.Join,
			"quote": func(in string) string { return fmt.Sprintf("\"%s\"", in) },
			"marshal": func(v interface{}) string {
				a, _ := json.Marshal(v)
				return string(a)
			},
			"dash": func(in int) string {
				if in == 0 {
					return "-"
				}
				return strconv.Itoa(in)
			},

			"substring": func(start, end int, s string) string {
				if start < 0 {
					return s[:end]
				}
				if end < 0 || end > len(s) {
					return s[start:]
				}
				return s[start:end]
			},

			"splitprefix": func(sep, orig string) map[string]string {
				parts := strings.Split(orig, sep)
				res := make(map[string]string, len(parts))
				for i, v := range parts {
					res["_"+strconv.Itoa(i)] = v
				}
				return res
			},

			"replace": func(old, new, src string) string { return strings.Replace(src, old, new, -1) },
		}
	)

	t := template.Must(template.New("").Funcs(funcs).Parse(templ.TemplateString))
	if err := t.Execute(&tstr, str); err != nil {
		return "", nil
	}

	temptype := templ.TemplateType

	if temptype == "TABULAR" {
		tbl := table.NewWriter()
		tbl.SetOutputMirror(&ostr)
		tbl.SetTitle(key)
		headers := templ.TableTitle

		headercolumns := strings.Split(headers, "|")
		trhdr := table.Row{}
		for _, header := range headercolumns {
			trhdr = append(trhdr, header)
		}
		tbl.AppendHeader(trhdr)

		ar := strings.Split(tstr.String(), ",")
		for _, recContent := range ar {
			trc := []table.Row{}
			ac := strings.Split(recContent, "|")
			tr := table.Row{}
			for _, colContent := range ac {
				tr = append(tr, colContent)
			}
			trc = append(trc, tr)
			tbl.AppendRows(trc)
		}

		tbl.Render()
	} else {
		return "\n" + tstr.String(), nil
	}

	return "\n" + ostr.String(), nil

}

// InitTemplates populates map of templates given as argument with output templates
func InitTemplates(otm map[string]*OutputTemplate) {

	// DS templates
	otm["networkListsDS"] = &OutputTemplate{TemplateName: "networkLists", TableTitle: "Name|ID|Type|ElementCount|SyncPoint|ReadOnly", TemplateType: "TABULAR", TemplateString: "{{range $index, $element := .NetworkLists}}{{if $index}},{{end}}{{.Name}}|{{.UniqueID}}|{{.Type}}|{{.ElementCount}}|{{.SyncPoint}}|{{.ReadOnly}}{{end}}"}

	// TF templates

}
