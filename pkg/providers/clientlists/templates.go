package clientlists

import (
	"bytes"
	"fmt"
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

var (
	templateFuncs = template.FuncMap{
		"tags": func(in []string) string {
			if len(in) == 0 {
				return ""
			}
			return "[" + strings.Join(in, ", ") + "]"
		},
		"noExpiration": func(in string) string {
			if in == "" {
				return "No expiration"
			}
			return in
		},
	}
)

// GetTemplate given map of templates and a key, returns template stored under this key
func GetTemplate(ots map[string]*OutputTemplate, key string) (*OutputTemplate, error) {
	if f, ok := ots[key]; ok && f != nil {
		return f, nil
	}

	return nil, fmt.Errorf("error: template '%s' not found", key)
}

// RenderTemplates renders template and returns it as a string
func RenderTemplates(ots map[string]*OutputTemplate, key string, str interface{}) (string, error) {
	var ostr, tstr bytes.Buffer
	templ, err := GetTemplate(ots, key)

	if err != nil {
		return "", nil
	}

	t := template.Must(template.New("").Funcs(templateFuncs).Parse(templ.TemplateString))
	if err := t.Execute(&tstr, str); err != nil {
		return "", nil
	}

	temptype := templ.TemplateType

	if temptype == "TABULAR" {
		tbl := table.NewWriter()
		tbl.SetOutputMirror(&ostr)
		tbl.SetTitle(templ.TemplateName)
		headers := templ.TableTitle

		headercolumns := strings.Split(headers, "|")
		trhdr := table.Row{}
		for _, header := range headercolumns {
			trhdr = append(trhdr, header)
		}
		tbl.AppendHeader(trhdr)

		ar := strings.Split(tstr.String(), "<<>>")
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
	// Rows delimiter is <<>>
	// Cells delimiter is |
	otm["clientListsDS"] = &OutputTemplate{
		TemplateType:   "TABULAR",
		TemplateName:   "Client lists",
		TableTitle:     "ListId|Name|Tags|Type|ItemsCount|Version|ReadOnly|UpdateDate|Notes",
		TemplateString: `{{range $index, $element := .Content}}{{if $index}}<<>>{{end}}{{.ListID}}|{{.Name}}|{{tags .Tags}}|{{.Type}}|{{.ItemsCount}}|{{.Version}}|{{.ReadOnly}}|{{.UpdateDate}}|{{.Notes}}{{end}}`,
	}
	otm["clientListDS"] = &OutputTemplate{
		TemplateType:   "TABULAR",
		TemplateName:   "Client list",
		TableTitle:     "ListId|Name|Tags|Type|ItemsCount|Version|ReadOnly|UpdateDate|Staging|Production|Notes",
		TemplateString: `{{range $index, $element := .}}{{if $index}}<<>>{{end}}{{.ListID}}|{{.Name}}|{{tags .Tags}}|{{.Type}}|{{.ItemsCount}}|{{.Version}}|{{.ReadOnly}}|{{.UpdateDate}}|{{.StagingActivationStatus}}|{{.ProductionActivationStatus}}|{{.Notes}}{{end}}`,
	}
	otm["userClientListItemsDS"] = &OutputTemplate{
		TemplateType:   "TABULAR",
		TemplateName:   "Client list items",
		TableTitle:     "UserID/Username|Description|Tags|CreateDate|UpdateDate|ExpirationDate",
		TemplateString: `{{range $index, $element := .Items}}{{if $index}}<<>>{{end}}{{.Value}}|{{.Description}}|{{tags .Tags}}|{{.CreateDate}}|{{.UpdateDate}}|{{noExpiration .ExpirationDate}}{{end}}`,
	}
	otm["ipClientListItemsDS"] = &OutputTemplate{
		TemplateType:   "TABULAR",
		TemplateName:   "Client list items",
		TableTitle:     "IP/CIDR|Description|Tags|CreateDate|UpdateDate|ExpirationDate",
		TemplateString: `{{range $index, $element := .Items}}{{if $index}}<<>>{{end}}{{.Value}}|{{.Description}}|{{tags .Tags}}|{{.CreateDate}}|{{.UpdateDate}}|{{noExpiration .ExpirationDate}}{{end}}`,
	}
	otm["asnClientListItemsDS"] = &OutputTemplate{
		TemplateType:   "TABULAR",
		TemplateName:   "Client list items",
		TableTitle:     "ASN|Description|Tags|CreateDate|UpdateDate|ExpirationDate",
		TemplateString: `{{range $index, $element := .Items}}{{if $index}}<<>>{{end}}{{.Value}}|{{.Description}}|{{tags .Tags}}|{{.CreateDate}}|{{.UpdateDate}}|{{noExpiration .ExpirationDate}}{{end}}`,
	}
	otm["geoClientListItemsDS"] = &OutputTemplate{
		TemplateType:   "TABULAR",
		TemplateName:   "Client list items",
		TableTitle:     "GEO|Description|Tags|CreateDate|UpdateDate|ExpirationDate",
		TemplateString: `{{range $index, $element := .Items}}{{if $index}}<<>>{{end}}{{.Value}}|{{.Description}}|{{tags .Tags}}|{{.CreateDate}}|{{.UpdateDate}}|{{noExpiration .ExpirationDate}}{{end}}`,
	}
	otm["tlsFingerprintClientListItemsDS"] = &OutputTemplate{
		TemplateType:   "TABULAR",
		TemplateName:   "Client list items",
		TableTitle:     "TLS fingerprint|Description|Tags|CreateDate|UpdateDate|ExpirationDate",
		TemplateString: `{{range $index, $element := .Items}}{{if $index}}<<>>{{end}}{{.Value}}|{{.Description}}|{{tags .Tags}}|{{.CreateDate}}|{{.UpdateDate}}|{{noExpiration .ExpirationDate}}{{end}}`,
	}
	otm["fileHashClientListItemsDS"] = &OutputTemplate{
		TemplateType:   "TABULAR",
		TemplateName:   "Client list items",
		TableTitle:     "File HASH|Description|Tags|CreateDate|UpdateDate|ExpirationDate",
		TemplateString: `{{range $index, $element := .Items}}{{if $index}}<<>>{{end}}{{.Value}}|{{.Description}}|{{tags .Tags}}|{{.CreateDate}}|{{.UpdateDate}}|{{noExpiration .ExpirationDate}}{{end}}`,
	}
}
