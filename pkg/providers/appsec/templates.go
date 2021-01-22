package appsec

import (
	"bytes"
	"encoding/json"
	"fmt"
	"strconv"
	"strings"
	"text/template"

	"github.com/jedib0t/go-pretty/v6/table"
)

type OutputTemplates map[string]*OutputTemplate

type OutputTemplate struct {
	TemplateName   string
	TemplateType   string
	TableTitle     string
	TemplateString string
}

func GetTemplate(ots map[string]*OutputTemplate, key string) (*OutputTemplate, error) {
	if f, ok := ots[key]; ok && f != nil {
		//	fmt.Printf("%s is in the OutputTemplate >> %+v\n", key, f)
		return f, nil
	} else {
		//	fmt.Printf("%s is NOT in the OutputTemplate!\n", key)
		return nil, fmt.Errorf("Error not found")
	}
}

func RenderTemplates(ots map[string]*OutputTemplate, key string, str interface{}) (string, error) {
	var ostr, tstr bytes.Buffer
	templ, ok := GetTemplate(ots, key)

	if ok == nil {

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
					} else {
						return strconv.Itoa(in)
					}
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
			tbl.SetOutputMirror(&ostr) //os.Stdout)
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
	return "", nil
}

func InitTemplates(otm map[string]*OutputTemplate) {
	otm["advancedSettingsLoggingDS"] = &OutputTemplate{TemplateName: "advancedSettingsLoggingDS", TableTitle: "Allow Sampling|Cookies|Custom Headers|Standard Headers", TemplateType: "TABULAR", TemplateString: "{{.AllowSampling}}|{{.Cookies.Type}} {{.Cookies.Values}}|{{.CustomHeaders.Type}} {{.CustomHeaders.Values}}|{{.StandardHeaders.Type}} {{.StandardHeaders.Values}}"}
	otm["advancedSettingsPrefetchDS"] = &OutputTemplate{TemplateName: "advancedSettingsPrefetchDS", TableTitle: "Enable App Layer|All Extension|Enable Rate Controls|Extensions", TemplateType: "TABULAR", TemplateString: "{{.EnableAppLayer}}|{{.AllExtensions}}|{{.EnableRateControls}}|{{range $index, $element := .Extensions}}{{.}} {{end}}"}
	otm["apiHostnameCoverageMatchTargetsDS"] = &OutputTemplate{TemplateName: "apiHostnameCoverageMatchTargetsDS", TableTitle: "Hostnames|Target ID|Type", TemplateType: "TABULAR", TemplateString: "{{range $index, $element := .MatchTargets.WebsiteTargets}}{{if $index}},{{end}}{{.Hostnames}}|{{.TargetID}}|{{.Type}}{{end}}"}
	//Extensions
	otm["apiEndpointsDS"] = &OutputTemplate{TemplateName: "apiEndpointsDS", TableTitle: "ID|Endpoint Name", TemplateType: "TABULAR", TemplateString: "{{range $index, $element := .APIEndpoints}}{{if $index}},{{end}}{{.ID}}|{{.Name}}{{end}}"}
	otm["policyApiEndpointsDS"] = &OutputTemplate{TemplateName: "policyApiEndpointsDS", TableTitle: "ID|Endpoint Name", TemplateType: "TABULAR", TemplateString: "{{range $index, $element := .APIEndpoints}}{{if $index}},{{end}}{{.ID}}|{{.Name}}{{end}}"}
	otm["apiHostnameCoverageDS"] = &OutputTemplate{TemplateName: "apiHostnameCoverageDS", TableTitle: "Config ID|Config Name|Version|Status|Has Match Target|Hostname", TemplateType: "TABULAR", TemplateString: "{{range $index, $element := .HostnameCoverage}}{{if $index}},{{end}}{{.Configuration.ID}}|{{.Configuration.Name}}|{{.Configuration.Version}}|{{.Status}}|{{.HasMatchTarget}}|{{.Hostname}}{{end}}"}
	otm["apiRequestConstraintsDS"] = &OutputTemplate{TemplateName: "apiRequestConstraintsDS", TableTitle: "ID|Action", TemplateType: "TABULAR", TemplateString: "{{range $index, $element := .APIEndpoints}}{{if $index}},{{end}}{{.ID}}|{{.Action}}{{end}}"}
	otm["configuration"] = &OutputTemplate{TemplateName: "Configurations", TableTitle: "Config_id|Name|Latest_version|Version_active_in_staging|Version_active_in_production", TemplateType: "TABULAR", TemplateString: "{{range $index, $element := .Configurations}}{{if $index}},{{end}}{{.ID}}|{{.Name}}|{{.LatestVersion}}|{{.StagingVersion}}|{{.ProductionVersion}}{{end}}"}
	otm["configurationVersion"] = &OutputTemplate{TemplateName: "ConfigurationVersion", TableTitle: "Version Number|Staging Status|Production Status", TemplateType: "TABULAR", TemplateString: "{{range $index, $element := .VersionList}}{{if $index}},{{end}}{{.Version}}|{{.Staging.Status}}|{{.Production.Status}}{{end}}"}
	otm["contractsgroupsDS"] = &OutputTemplate{TemplateName: "contractsgroupsDS", TableTitle: "ContractID|GroupID|Name", TemplateType: "TABULAR", TemplateString: "{{range $index, $element := .ContractGroups}}{{if $index}},{{end}}{{.ContractID}}|{{.GroupID}}|{{.DisplayName}}{{end}}"}
	otm["failoverHostnamesDS"] = &OutputTemplate{TemplateName: "failoverHostnamesDS", TableTitle: "Hostname", TemplateType: "TABULAR", TemplateString: "{{range $index, $element := .HostnameList}}{{if $index}},{{end}}{{.Hostname}}{{end}}"}
	otm["bypassNetworkListsDS"] = &OutputTemplate{TemplateName: "bypassNetworkListsDS", TableTitle: "Network List|ID", TemplateType: "TABULAR", TemplateString: "{{range $index, $element := .NetworkLists}}{{if $index}},{{end}}{{.Name}}|{{.ID}}{{end}}"}

	otm["penaltyBoxDS"] = &OutputTemplate{TemplateName: "penaltyBoxDS", TableTitle: "PenaltyBoxProtection|Action", TemplateType: "TABULAR", TemplateString: "{{.PenaltyBoxProtection}}|{{.Action}}"}

	otm["selectableHosts"] = &OutputTemplate{TemplateName: "selectableHosts", TableTitle: "Hostname", TemplateType: "TABULAR", TemplateString: "{{range .SelectableHosts}}{{.}},{{end}}"}
	otm["selectableHostsDS"] = &OutputTemplate{TemplateName: "selectableHosts", TableTitle: "Hostname|ConfigIDInProduction|ConfigNameInProduction", TemplateType: "TABULAR", TemplateString: "{{range .AvailableSet}}{{.Hostname}}|{{ dash .ConfigIDInProduction }}|{{.ConfigNameInProduction}},{{end}}"}
	otm["selectedHosts"] = &OutputTemplate{TemplateName: "selectedHosts", TableTitle: "Hostnames", TemplateType: "TABULAR", TemplateString: "{{range $index, $element := .SelectedHosts}}{{if $index}},{{end}}{{.}}{{end}}"}
	otm["selectedHostsDS"] = &OutputTemplate{TemplateName: "selectedHosts", TableTitle: "Hostnames", TemplateType: "TABULAR", TemplateString: "{{range $index, $element := .HostnameList}}{{if $index}},{{end}}{{.Hostname}}{{end}}"}
	otm["siemsettingsDS"] = &OutputTemplate{TemplateName: "siemsettingsDS", TableTitle: "Enable For All Policies|Enable Siem|Enabled Botman Siem Events|Siem Definition ID|FirewallPolicyIds", TemplateType: "TABULAR", TemplateString: "{{.EnableForAllPolicies}}|{{.EnableSiem}}|{{.EnabledBotmanSiemEvents}}|{{.SiemDefinitionID}}|{{.FirewallPolicyIds}}"}
	otm["siemDefinitionsDS"] = &OutputTemplate{TemplateName: "siemDefinitionsDS", TableTitle: "ID|Name", TemplateType: "TABULAR", TemplateString: "{{range $index, $element := .SiemDefinitions}}{{if $index}},{{end}}{{.ID}}|{{.Name}}{{end}}"}

	otm["ratePolicies"] = &OutputTemplate{TemplateName: "ratePolicies", TableTitle: "ID|Policy Name", TemplateType: "TABULAR", TemplateString: "{{range $index, $element := .RatePolicies}}{{if $index}},{{end}}{{.ID}}|{{.Name}}{{end}}"}
	otm["matchTargets"] = &OutputTemplate{TemplateName: "matchTargets", TableTitle: "ID|PolicyID", TemplateType: "TABULAR", TemplateString: "{{range $index, $element := .MatchTargets.WebsiteTargets}}{{if $index}},{{end}}{{.ID}}|{{.SecurityPolicy.PolicyID}}{{end}}"}
	otm["matchTargetDS"] = &OutputTemplate{TemplateName: "matchTarget", TableTitle: "ID|PolicyID", TemplateType: "TABULAR", TemplateString: "{{range $index, $element := .MatchTargets.WebsiteTargets}}{{if $index}},{{end}}{{.TargetID}}|{{.SecurityPolicy.PolicyID}}{{end}}"}
	otm["reputationProfiles"] = &OutputTemplate{TemplateName: "reputationProfiles", TableTitle: "ID|Name(Title)", TemplateType: "TABULAR", TemplateString: "{{range $index, $element := .ReputationProfiles}}{{if $index}},{{end}}{{.ID}}|{{.Name}}{{end}}"}
	otm["reputationAnalysisDS"] = &OutputTemplate{TemplateName: "reputationAnalysisDS", TableTitle: "forwardToHTTPHeader|forwardSharedIPToHTTPHeaderAndSIEM", TemplateType: "TABULAR", TemplateString: "{{.ForwardToHTTPHeader}}|{{.ForwardSharedIPToHTTPHeaderAndSIEM}}"}

	otm["reputationProfilesDS"] = &OutputTemplate{TemplateName: "reputationProfilesDS", TableTitle: "ID|Name(Title)", TemplateType: "TABULAR", TemplateString: "{{range $index, $element := .ReputationProfiles}}{{if $index}},{{end}}{{.ID}}|{{.Name}}{{end}}"}

	otm["reputationProfilesActions"] = &OutputTemplate{TemplateName: "reputationProfilesActions", TableTitle: "ID|Action", TemplateType: "TABULAR", TemplateString: "{{range $index, $element := .ReputationProfiles}}{{if $index}},{{end}}{{.ID}}| {{.Action}}{{end}}"}

	otm["customRules"] = &OutputTemplate{TemplateName: "customRules", TableTitle: "ID|Name", TemplateType: "TABULAR", TemplateString: "{{range $index, $element := .CustomRules}}{{if $index}},{{end}}{{.ID}}|{{.Name}}{{end}}"}
	otm["customDenyDS"] = &OutputTemplate{TemplateName: "customDenyDS", TableTitle: "ID|Name", TemplateType: "TABULAR", TemplateString: "{{range $index, $element := .CustomDenyList}}{{if $index}},{{end}}{{.ID}}|{{.Name}}{{end}}"}
	otm["customRuleActions"] = &OutputTemplate{TemplateName: "customRuleActions", TableTitle: "ID|Action", TemplateType: "TABULAR", TemplateString: "{{range .SecurityPolicies}}{{range $index, $element := .CustomRuleActions}}{{if $index}},{{end}}{{.ID}}|{{.Action}}{{end}}{{end}}"}
	otm["customRuleAction"] = &OutputTemplate{TemplateName: "customRuleAction", TableTitle: "ID|Name|Action", TemplateType: "TABULAR", TemplateString: "{{range $index, $element := .}}{{if $index}},{{end}}{{.RuleID}}|{{.Name}} |{{.Action}}{{end}}"}
	otm["ratePolicyActions"] = &OutputTemplate{TemplateName: "ratePolicyActions", TableTitle: "ID|Ipv4Action|Ipv6Action", TemplateType: "TABULAR", TemplateString: "{{range $index, $element := .RatePolicyActions}}{{if $index}},{{end}}{{.ID}}| {{.Ipv4Action}}|{{.Ipv6Action}}{{end}}"}
	otm["RulesDS"] = &OutputTemplate{TemplateName: "RulesDS", TableTitle: "ID|Action", TemplateType: "TABULAR", TemplateString: "{{range $index, $element := .RuleActions}}{{if $index}},{{end}}{{.ID}}|{{.Action}}{{end}}"}
	otm["EvalRulesActionsDS"] = &OutputTemplate{TemplateName: "evalRulesActions", TableTitle: "ID|Action", TemplateType: "TABULAR", TemplateString: "{{range $index, $element := .RuleActions}}{{if $index}},{{end}}{{.ID}}| {{.Action}}{{end}}"}
	otm["evalHostnamesDS"] = &OutputTemplate{TemplateName: "evalHostnamesDS", TableTitle: "Hostnames", TemplateType: "TABULAR", TemplateString: "{{range $index, $element := .Hostnames}}{{if $index}},{{end}}{{.}}{{end}}"}
	otm["rulesets"] = &OutputTemplate{TemplateName: "rulesets", TableTitle: "ID|Name(Title)", TemplateType: "TABULAR", TemplateString: "{{range .Rulesets}}{{range $index, $element := .Rules}}{{if $index}},{{end}}{{.ID}}| {{substring  0 45 .Title}}{{end}}{{end}}"}
	otm["securityPolicies"] = &OutputTemplate{TemplateName: "securityPolicies", TableTitle: "ID|Name", TemplateType: "TABULAR", TemplateString: "{{range $index, $element := .SecurityPolicies}}{{if $index}},{{end}}{{.ID}}|{{.Name}}{{end}}"}
	otm["securityPoliciesDS"] = &OutputTemplate{TemplateName: "securityPolicies", TableTitle: "ID|Name", TemplateType: "TABULAR", TemplateString: "{{range $index, $element := .Policies}}{{if $index}},{{end}}{{.PolicyID}}|{{.PolicyName}}{{end}}"}

	otm["ruleActions"] = &OutputTemplate{TemplateName: "ruleActions", TableTitle: "ID|Action", TemplateType: "TABULAR", TemplateString: "{{range .SecurityPolicies}}{{range $index, $element := .WebApplicationFirewall.RuleActions}}{{if $index}},{{end}}{{.ID}}| {{.Action}}{{end}}{{end}}"}
	otm["slowPostDS"] = &OutputTemplate{TemplateName: "slowPost", TableTitle: "Action|SLOW_RATE_THRESHOLD RATE|SLOW_RATE_THRESHOLD PERIOD|DURATION_THRESHOLD TIMEOUT", TemplateType: "TABULAR", TemplateString: "{{.Action}}|{{.SlowRateThreshold.Rate}}|{{.SlowRateThreshold.Period}}|{{.DurationThreshold.Timeout}}"}
	otm["slowPost"] = &OutputTemplate{TemplateName: "slowPost", TableTitle: "Action|SLOW_RATE_THRESHOLD RATE|SLOW_RATE_THRESHOLD PERIOD|DURATION_THRESHOLD TIMEOUT", TemplateType: "TABULAR", TemplateString: "{{range $index, $element := .SecurityPolicies}}{{if $index}},{{end}}{{.SlowPost.Action}}|{{.SlowPost.DurationThreshold.Timeout}}|{{.SlowPost.SlowRateThreshold.Rate}}|{{.SlowPost.SlowRateThreshold.Period}}{{end}}"}
	otm["wafModesDS"] = &OutputTemplate{TemplateName: "wafMode", TableTitle: "Current|Mode|Eval", TemplateType: "TABULAR", TemplateString: "{{.Current}}|{{.Mode}}|{{.Eval}}"}
	otm["versionNotesDS"] = &OutputTemplate{TemplateName: "versionNotesDS", TableTitle: "Version Notes", TemplateType: "TABULAR", TemplateString: "{{.Notes}}"}
	otm["wafProtectionDS"] = &OutputTemplate{TemplateName: "wafProtection", TableTitle: "APIConstraints|ApplicationLayerControls|BotmanControls|NetworkLayerControls|RateControls|ReputationControls|SlowPostControls", TemplateType: "TABULAR", TemplateString: "{{.ApplyAPIConstraints}}|{{.ApplyApplicationLayerControls}}|{{.ApplyBotmanControls}}|{{.ApplyNetworkLayerControls}}|{{.ApplyRateControls}}|{{.ApplyReputationControls}}|{{.ApplySlowPostControls}}"}
	otm["AttackGroupActionDS"] = &OutputTemplate{TemplateName: "AttackGroupAction", TableTitle: "GroupID|Action", TemplateType: "TABULAR", TemplateString: "{{range $index, $element := .AttackGroupActions}}{{if $index}},{{end}}{{.Group}}| {{.Action}}{{end}}"}

	otm["rateProtectionDS"] = &OutputTemplate{TemplateName: "rateProtection", TableTitle: "APIConstraints|ApplicationLayerControls|BotmanControls|NetworkLayerControls|RateControls|ReputationControls|SlowPostControls", TemplateType: "TABULAR", TemplateString: "{{.ApplyAPIConstraints}}|{{.ApplyApplicationLayerControls}}|{{.ApplyBotmanControls}}|{{.ApplyNetworkLayerControls}}|{{.ApplyRateControls}}|{{.ApplyReputationControls}}|{{.ApplySlowPostControls}}"}
	otm["reputationProtectionDS"] = &OutputTemplate{TemplateName: "reputationProtection", TableTitle: "APIConstraints|ApplicationLayerControls|BotmanControls|NetworkLayerControls|RateControls|ReputationControls|SlowPostControls", TemplateType: "TABULAR", TemplateString: "{{.ApplyAPIConstraints}}|{{.ApplyApplicationLayerControls}}|{{.ApplyBotmanControls}}|{{.ApplyNetworkLayerControls}}|{{.ApplyRateControls}}|{{.ApplyReputationControls}}|{{.ApplySlowPostControls}}"}
	otm["slowpostProtectionDS"] = &OutputTemplate{TemplateName: "slowpostProtection", TableTitle: "APIConstraints|ApplicationLayerControls|BotmanControls|NetworkLayerControls|RateControls|ReputationControls|SlowPostControls", TemplateType: "TABULAR", TemplateString: "{{.ApplyAPIConstraints}}|{{.ApplyApplicationLayerControls}}|{{.ApplyBotmanControls}}|{{.ApplyNetworkLayerControls}}|{{.ApplyRateControls}}|{{.ApplyReputationControls}}|{{.ApplySlowPostControls}}"}
	otm["RuleConditionException"] = &OutputTemplate{TemplateName: "RuleConditionException", TableTitle: "Conditions|Exceptions", TemplateType: "TABULAR", TemplateString: "{{range $index, $element :=.Conditions}}True{{else}}False{{end}}|{{with .Exception}}True{{else}}False{{end}}"}

	// TF templates
	otm["AdvancedSettingsLogging.tf"] = &OutputTemplate{TemplateName: "AdvancedSettingsLogging.tf", TableTitle: "AdvancedSettingsLogging", TemplateType: "TERRAFORM", TemplateString: "\nresource \"akamai_appsec_advanced_settings_logging\" \"appsec_advanced_settings_logging\" { \n config_id = {{.ConfigID}}\n version = {{.Version}}\n logging  = <<-EOF  {{marshal .AdvancedOptions.Logging}} EOF \n }"}
	otm["AdvancedSettingsPrefetch.tf"] = &OutputTemplate{TemplateName: "AdvancedSettingsPrefetch.tf", TableTitle: "AdvancedSettingsPrefetch", TemplateType: "TERRAFORM", TemplateString: "\nresource \"akamai_appsec_advanced_settings_prefetch\" \"appsec_advanced_settings_prefetch\" { \n  config_id = {{.ConfigID}}\n  version = {{.Version}}\n enable_app_layer = {{.AdvancedOptions.Prefetch.EnableAppLayer}} \n all_extensions = {{.AdvancedOptions.Prefetch.AllExtensions}}\n enable_rate_controls = {{.AdvancedOptions.Prefetch.EnableRateControls}}\n extensions = [{{  range $index, $element := .AdvancedOptions.Prefetch.Extensions }}{{if $index}},{{end}}{{quote .}}{{end}}] \n } \n"}
	otm["AttackGroupAction.tf"] = &OutputTemplate{TemplateName: "AttackGroupAction.tf", TableTitle: "AttackGroupAction", TemplateType: "TERRAFORM", TemplateString: "{{ $config := .ConfigID }}{{ $version := .Version }}{{ $prev_secpolicy := \"\" }}{{range .SecurityPolicies}}{{$prev_secpolicy := .ID}} {{range $index, $element := .WebApplicationFirewall.AttackGroupActions}}\nresource \"akamai_appsec_attack_group_action\" \"appsec_attack_group_action{{if $index}}_{{$index}}{{end}}\" { \n  config_id = {{$config}}\n  version = {{$version}}\n  security_policy_id = \"{{$prev_secpolicy}}\" \n  attack_group = \"{{.Group}}\" \n  attack_group_action = \"{{.Action}}\" \n }\n{{end}}{{end}}"}
	otm["CustomRule.tf"] = &OutputTemplate{TemplateName: "CustomRule.tf", TableTitle: "CustomRule", TemplateType: "TERRAFORM", TemplateString: "{{ $config := .ConfigID }}{{range $index, $element := .CustomRules}}{{if $index}},{{end}}\nresource \"akamai_appsec_custom_rule\" \"appsec_custom_rule{{if $index}}_{{$index}}{{end}}\" { \n config_id = {{$config}}\n  custom_rule = <<-EOF {{marshal .}}  \n EOF \n }\n {{end}}"}
	otm["CustomRuleAction.tf"] = &OutputTemplate{TemplateName: "CustomRuleAction.tf", TableTitle: "CustomRuleAction", TemplateType: "TERRAFORM", TemplateString: "{{ $config := .ConfigID }}{{ $version := .Version }}{{ $prev_secpolicy := \"\" }}{{  range $index, $element := .SecurityPolicies }}{{$prev_secpolicy:=$element.ID}}  {{  range $index, $element := .CustomRuleActions }}  \nresource \"akamai_appsec_custom_rule_action\" \"appsec_custom_rule_action{{if $index}}_{{$index}}{{end}}\" { \n config_id = {{$config}}\n version = {{$version}}\n security_policy_id = \"{{$prev_secpolicy}}\"  \n custom_rule_id = {{.ID}} \n custom_rule_action = \"{{.Action}}\" \n } \n {{end}}{{end}}"}
	otm["MatchTarget.tf"] = &OutputTemplate{TemplateName: "MatchTarget.tf", TableTitle: "MatchTarget", TemplateType: "TERRAFORM", TemplateString: "{{ $config := .ConfigID }}{{ $version := .Version }}{{range $index, $element := .MatchTargets.WebsiteTargets}}\nresource \"akamai_appsec_match_target\" \"appsec_match_target{{if $index}}_{{$index}}{{end}}\" { \n  config_id = {{$config}}\n  version = {{$version}}\n  match_target = <<-EOF {{marshal .}}  \n EOF  \n }\n {{end}}"}
	otm["RatePolicy.tf"] = &OutputTemplate{TemplateName: "RatePolicy.tf", TableTitle: "RatePolicy", TemplateType: "TERRAFORM", TemplateString: "{{ $config := .ConfigID }}{{ $version := .Version }}{{range $index, $element := .RatePolicies}}\nresource \"akamai_appsec_rate_policy\" \"appsec_rate_policy{{if $index}}_{{$index}}{{end}}\" { \n  config_id = {{ $config }}\n  version = {{ $version }}\n  rate_policy = <<-EOF {{marshal .}}  \n EOF \n \n }\n{{end}}"}
	otm["RatePolicyAction.tf"] = &OutputTemplate{TemplateName: "RatePolicyAction.tf", TableTitle: "RatePolicyAction", TemplateType: "TERRAFORM", TemplateString: "{{ $config := .ConfigID }}{{ $version := .Version }}{{ $prev_secpolicy := \"\" }}{{range .SecurityPolicies}}{{$prev_secpolicy := .ID}} {{  range $index, $element := .RatePolicyActions }}\nresource \"akamai_appsec_rate_policy_action\" \"appsec_rate_policy_action{{if $index}}_{{$index}}{{end}}\" { \n  config_id = {{$config}}\n  version = {{$version}}\n  security_policy_id = \"{{$prev_secpolicy}}\" \n  rate_policy_id = {{.ID}} \n  ipv4_action = \"{{.Ipv4Action}}\" \n  ipv6_action = \"{{.Ipv6Action}}\" \n }\n {{end}} {{end}}"}
	otm["ReputationProfile.tf"] = &OutputTemplate{TemplateName: "ReputationProfile.tf", TableTitle: "ReputationProfile", TemplateType: "TERRAFORM", TemplateString: "{{ $config := .ConfigID }}{{ $version := .Version }}{{range $index, $element := .ReputationProfiles}}\nresource \"akamai_appsec_reputation_profile\" \"appsec_reputation_profile{{if $index}}_{{$index}}{{end}}\" { \n  config_id = {{ $config}}\n  version = {{ $version }}\n  reputation_profile = <<-EOF {{marshal .}}  \n \n EOF \n }{{end}}"}
	otm["ReputationProfileAction.tf"] = &OutputTemplate{TemplateName: "ReputationProfileAction.tf", TableTitle: "ReputationProfileAction", TemplateType: "TERRAFORM", TemplateString: "{{ $config := .ConfigID }}{{ $version := .Version }}{{ $prev_secpolicy := \"\" }}{{range .SecurityPolicies}}{{$prev_secpolicy := .ID}} {{range $index, $element := .ClientReputation.ReputationProfileActions}}\nresource \"akamai_appsec_reputation_profile_action\" \"appsec_reputation_profile_action{{if $index}}_{{$index}}{{end}}\" { \n config_id = {{ $config }}\n  version = {{ $version }}\n  security_policy_id = \"{{$prev_secpolicy}}\" \n  reputation_profile_id =  \"{{.Action}}\" \n }\n{{end}}{{end}}"}
	otm["RuleAction.tf"] = &OutputTemplate{TemplateName: "RuleAction.tf", TableTitle: "RuleAction", TemplateType: "TERRAFORM", TemplateString: "{{ $config := .ConfigID }}{{ $version := .Version }}{{ $prev_secpolicy := \"\" }}{{range .SecurityPolicies}}{{$prev_secpolicy := .ID}} {{range $index, $element := .WebApplicationFirewall.RuleActions}}\nresource \"akamai_appsec_rule_action\" \"appsec_rule_action{{if $index}}_{{$index}}{{end}}\" { \n  config_id = {{$config}}\n  version = {{$version}} \n  security_policy_id = \"{{$prev_secpolicy}}\" \n  rule_id = {{.ID}} \n  action = \"{{.Action}}\" \n }\n{{end}}{{end}}"}
	otm["SecurityPolicy.tf"] = &OutputTemplate{TemplateName: "SecurityPolicy.tf", TableTitle: "SecurityPolicy", TemplateType: "TERRAFORM", TemplateString: "{{ $config := .ConfigID }}{{ $version := .Version }}{{ $prev_secpolicy := \"\" }}{{range $index, $element :=  .SecurityPolicies}}{{$prev_secpolicy := .ID}}\nresource \"akamai_appsec_security_policy\" \"appsec_security_policy{{if $index}}_{{$index}}{{end}}\" { \n  config_id = {{ $config }}\n  version = {{ $version }}\n  security_policy_name = \"{{.Name}}\" \n  security_policy_prefix = \"{{substring 0 2 .ID}}\" \n  default_settings = true \n }{{end}}"}
	otm["SelectedHostname.tf"] = &OutputTemplate{TemplateName: "SelectedHostname.tf", TableTitle: "SelectedHostname", TemplateType: "TERRAFORM", TemplateString: "\nresource \"akamai_appsec_selected_hostname\" \"appsec_selected_hostname\" { \n config_id = {{.ConfigID}}\n version = {{.Version}}\n hostnames = [{{  range $index, $element := .SelectedHosts }}{{if $index}},{{end}}{{quote .}}{{end}}] \n }"}

}
