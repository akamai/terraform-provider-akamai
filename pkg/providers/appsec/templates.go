package appsec

import (
	"bytes"
	"fmt"
	"strings"
	"text/template"

	"github.com/jedib0t/go-pretty/v6/table"
)

type OutputTemplates map[string][]string

func (ot OutputTemplates) Add(key, value string) {
	_, ok := ot[key]
	if !ok {
		ot[key] = make([]string, 0, 20)
	}
	ot[key] = append(ot[key], value)
}

func (s OutputTemplates) GetTemplate(key string, idx int) (string, bool) {
	slice, ok := s[key]
	if !ok || len(slice) == 0 {
		return "", false
	}
	if len(slice) < idx {
		return "", false
	} else {
		return s[key][idx], true
	}
}

func (s OutputTemplates) RenderTemplates(key string, str interface{}) (string, error) {
	var ostr, tstr bytes.Buffer
	templ, ok := s.GetTemplate(key, 0)

	if ok {

		var (
			funcs = template.FuncMap{
				"join":  strings.Join,
				"quote": func(in string) string { return fmt.Sprintf("\"%s\"", in) },
			}
		)

		t := template.Must(template.New("").Funcs(funcs).Parse(templ))
		if err := t.Execute(&tstr, str); err != nil {
			return "", nil
		}

		temptype, ok := s.GetTemplate(key, 1)
		if ok {
			if temptype == "TABULAR" {
				tbl := table.NewWriter()
				tbl.SetOutputMirror(&ostr) //os.Stdout)
				tbl.SetTitle(key)
				headers, ok := s.GetTemplate(key, 2)
				if ok {
					headercolumns := strings.Split(headers, "|")
					trhdr := table.Row{}
					for _, header := range headercolumns {
						trhdr = append(trhdr, header)
					}
					tbl.AppendHeader(trhdr)
				}

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
		}

		return "\n" + ostr.String(), nil
	}
	return "", nil
}

func (s *OutputTemplates) InitTemplates() {

	s.Add("selectableHosts", "{{range .SelectableHosts}}{{.}},{{end}}")
	s.Add("selectableHosts", "TABULAR")
	s.Add("selectableHosts", "Hostname")

	s.Add("selectedHosts", "{{range $index, $element := .SelectedHosts}}{{if $index}},{{end}}{{.}}{{end}}")
	s.Add("selectedHosts", "TABULAR")
	s.Add("selectedHosts", "Hostnames")

	s.Add("selectedHosts.tf", "\nresource \"akamai_appsec_selected_hostnames\" \"appsecselectedhostnames\" { \n config_id = {{.ConfigID}}\n version = {{.Version}}\n hostnames = [{{  range $index, $element := .SelectedHosts }}{{if $index}},{{end}}{{quote .}}{{end}}] \n }")
	s.Add("selectedHosts.tf", "TERRAFORM")

	s.Add("ratePolicies", "{{range $index, $element := .RatePolicies}}{{if $index}},{{end}}{{.ID}}|{{.Name}}{{end}}")
	s.Add("ratePolicies", "TABULAR")
	s.Add("ratePolicies", "ID|PolicyID")

	s.Add("matchTargets", "{{range $index, $element := .MatchTargets.WebsiteTargets}}{{if $index}},{{end}}{{.ID}}|{{.SecurityPolicy.PolicyID}}{{end}}")
	s.Add("matchTargets", "TABULAR")
	s.Add("matchTargets", "ID|PolicyID")

	s.Add("reputationProfiles", "{{range $index, $element := .ReputationProfiles}}{{if $index}},{{end}}{{.ID}}|{{.Name}}{{end}}")
	s.Add("reputationProfiles", "TABULAR")
	s.Add("reputationProfiles", "ID|Name(Title)")

	s.Add("customRules", "{{range $index, $element := .CustomRules}}{{if $index}},{{end}}{{.ID}}|{{.Name}}{{end}}")
	s.Add("customRules", "TABULAR")
	s.Add("customRules", "ID|Name")

	s.Add("rulesets", "{{range .Rulesets}}{{range $index, $element := .Rules}}{{if $index}},{{end}}{{.ID}}| {{.Title}}{{end}}{{end}}")
	s.Add("rulesets", "TABULAR")
	s.Add("rulesets", "ID|Name(Title)")

	s.Add("securityPolicies", "{{range $index, $element := .SecurityPolicies}}{{if $index}},{{end}}{{.ID}}|{{.Name}}{{end}}")
	s.Add("securityPolicies", "TABULAR")
	s.Add("securityPolicies", "ID|Name")

	s.Add("ruleActions", "{{range .SecurityPolicies}}{{range $index, $element := .WebApplicationFirewall.RuleActions}}{{if $index}},{{end}}{{.ID}}| {{.Action}}{{end}}{{end}}")
	s.Add("ruleActions", "TABULAR")
	s.Add("ruleActions", "ID|Action")

}
