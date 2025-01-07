package property

import (
	"encoding/json"
	"errors"
	"fmt"
	"net/http"
	"path"
	"regexp"
	"strings"
	"testing"

	"github.com/akamai/AkamaiOPEN-edgegrid-golang/v9/pkg/iam"
	"github.com/akamai/AkamaiOPEN-edgegrid-golang/v9/pkg/papi"
	"github.com/akamai/terraform-provider-akamai/v6/pkg/common/ptr"
	"github.com/akamai/terraform-provider-akamai/v6/pkg/common/test"
	"github.com/akamai/terraform-provider-akamai/v6/pkg/common/testutils"
	"github.com/hashicorp/go-cty/cty"
	"github.com/hashicorp/terraform-plugin-sdk/v2/diag"
	"github.com/hashicorp/terraform-plugin-testing/helper/resource"
	"github.com/hashicorp/terraform-plugin-testing/knownvalue"
	"github.com/hashicorp/terraform-plugin-testing/plancheck"
	"github.com/hashicorp/terraform-plugin-testing/terraform"
	"github.com/hashicorp/terraform-plugin-testing/tfjsonpath"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

// TestPropertyLifecycle tests various lifecycle workflows that result in a success
func TestPropertyLifecycle(t *testing.T) {

	// defaultChecker contains basic checks for every test case, that can be built upon
	defaultChecker := test.NewStateChecker("akamai_property.test").
		CheckEqual("id", "prp_4").
		CheckEqual("hostnames.0.edge_hostname_id", "ehn_123").
		CheckEqual("name", "test_property").
		CheckEqual("contract_id", "ctr_1").
		CheckEqual("group_id", "grp_2").
		CheckEqual("product_id", "prd_3").
		CheckEqual("rule_warnings.#", "0").
		CheckEqual("hostnames.0.cname_to", "to.test.domain").
		CheckEqual("latest_version", "1").
		CheckEqual("staging_version", "0").
		CheckEqual("production_version", "0").
		CheckEqual("rules", `{"rules":{"name":"default","options":{}}}`)

	// basicData holds basic, common data across test cases
	basicData := mockPropertyData{
		propertyName:  "test_property",
		productID:     "prd_3",
		propertyID:    "prp_4",
		groupID:       "grp_2",
		contractID:    "ctr_1",
		assetID:       "aid_5555",
		latestVersion: 1,
		versions: papi.PropertyVersionItems{
			Items: []papi.PropertyVersionGetItem{
				{
					StagingStatus:    papi.VersionStatusInactive,
					ProductionStatus: papi.VersionStatusInactive,
					PropertyVersion:  1,
				},
			},
		},
		hostnames: papi.HostnameResponseItems{
			Items: []papi.Hostname{
				{
					CnameType:            "EDGE_HOSTNAME",
					EdgeHostnameID:       "ehn_123",
					CnameFrom:            "from.test.domain",
					CnameTo:              "to.test.domain",
					CertProvisioningType: "DEFAULT",
				},
			},
		},
	}

	// basicDataWithDefaultRules extends basic data with the default rules
	basicDataWithDefaultRules := basicData
	basicDataWithDefaultRules.ruleTree = mockRuleTreeData{
		rules: papi.Rules{
			Name: "default",
		},
	}

	// updatedHostname contains details about the updated hostname
	updatedHostname := papi.HostnameResponseItems{
		Items: []papi.Hostname{
			{
				CnameType:            "EDGE_HOSTNAME",
				CnameFrom:            "from.test.domain",
				CnameTo:              "to2.test.domain",
				CertProvisioningType: "DEFAULT",
				EdgeHostnameID:       "ehn_123",
			},
		},
	}

	// versionStagingActive represents version of property active only on staging network
	versionStagingActive := papi.PropertyVersionItems{
		Items: []papi.PropertyVersionGetItem{
			{
				StagingStatus:    papi.VersionStatusActive,
				ProductionStatus: papi.VersionStatusInactive,
				PropertyVersion:  1,
			},
		},
	}

	// versionProductionActive represents version of property active only on production network
	versionProductionActive := papi.PropertyVersionItems{
		Items: []papi.PropertyVersionGetItem{
			{
				StagingStatus:    papi.VersionStatusInactive,
				ProductionStatus: papi.VersionStatusActive,
				PropertyVersion:  1,
			},
		},
	}

	// mockLatestVersionNotActive represents a workflow where no property version is active on any of the networks
	mockLatestVersionNotActive := func(p *mockProperty) {
		// create
		mockResourcePropertyCreateWithVersionHostnames(p)
		// read x2
		mockResourcePropertyRead(p, 2)
		// read x1 before update
		mockResourcePropertyRead(p)
		// update
		p.mockGetPropertyVersion()
		p.hostnames = updatedHostname
		p.mockUpdatePropertyVersionHostnames()
		// read x2
		mockResourcePropertyRead(p, 2)
		// delete
		p.mockRemoveProperty()
	}

	// mockLatestVersionActiveOnStaging represents a workflow where the second version of property is created and is active on staging network
	mockLatestVersionActiveOnStaging := func(p *mockProperty) {
		// create
		mockResourcePropertyCreateWithVersionHostnames(p)
		// read x2
		mockResourcePropertyRead(p, 2)
		// mock staging version active
		p.versions = versionStagingActive
		// read
		mockResourcePropertyRead(p)
		// update creates new version with updated hostname and version info
		p.mockGetPropertyVersion()
		p.mockCreatePropertyVersion()
		p.latestVersion = 2
		p.hostnames = updatedHostname
		p.mockUpdatePropertyVersionHostnames()
		// read x2
		mockResourcePropertyRead(p, 2)
		// delete
		p.mockRemoveProperty()
	}

	// mockLatestVersionActiveOnProduction represents a workflow where the second version of property is created and is active on production network
	mockLatestVersionActiveOnProduction := func(p *mockProperty) {
		// create
		mockResourcePropertyCreateWithVersionHostnames(p)
		// read x2
		mockResourcePropertyRead(p, 2)
		// mock production version active
		p.versions = versionProductionActive
		// read
		mockResourcePropertyRead(p)
		// update creates new version with updated hostname and version info
		p.mockGetPropertyVersion()
		p.mockCreatePropertyVersion()
		p.latestVersion = 2
		p.hostnames = updatedHostname
		p.mockUpdatePropertyVersionHostnames()
		// read x2
		mockResourcePropertyRead(p, 2)
		// delete
		p.mockRemoveProperty()
	}

	// grouped tests that always have 2 steps and the config file names are always "step0.tf" and "step1.tf". For different tests,
	// add them as separate cases at the end of this function.
	tests := map[string]struct {
		init             func(*testing.T, *mockProperty)
		checksForCreate  resource.TestCheckFunc
		checksForUpdate  resource.TestCheckFunc
		configPlanChecks resource.ConfigPlanChecks
		configDir        string
		updateError      *regexp.Regexp
	}{
		"Lifecycle: property is destroyed and recreated when name is changed": {
			init: func(t *testing.T, p *mockProperty) {
				// set initial data
				p.mockPropertyData = basicDataWithDefaultRules
				// create
				mockResourcePropertyCreateWithVersionHostnames(p)
				// read x2
				mockResourcePropertyRead(p, 2)
				// read x1 before update
				mockResourcePropertyRead(p)
				// delete
				p.mockRemoveProperty()
				// assign new values for mock data
				p.propertyID = "prp_5"
				p.propertyName = "renamed_property"
				p.hostnames = updatedHostname
				// create
				mockResourcePropertyCreateWithVersionHostnames(p)
				// read x2
				mockResourcePropertyRead(p, 2)
				// delete
				p.mockRemoveProperty()
			},
			configDir:       "forceNewOnNameChange",
			checksForCreate: defaultChecker.Build(),
			checksForUpdate: defaultChecker.
				CheckEqual("id", "prp_5").
				CheckEqual("name", "renamed_property").
				CheckEqual("hostnames.0.cname_to", "to2.test.domain").Build(),
		},
		"Lifecycle: create with propertyID (bootstrap)": {
			init: func(t *testing.T, p *mockProperty) {
				// set initial data
				p.mockPropertyData = basicDataWithDefaultRules
				// create (without creation of property, as it was created with bootstrap resource)
				p.mockUpdatePropertyVersionHostnames()
				// read x2
				mockResourcePropertyRead(p, 2)
				// read x1 before update
				mockResourcePropertyRead(p)
				// mock updated hostnames
				p.hostnames = updatedHostname
				// update
				p.mockGetPropertyVersion()
				p.mockUpdatePropertyVersionHostnames()
				// read x2
				mockResourcePropertyRead(p, 2)
				// no delete as the resource is maintained by bootstrap resource
			},
			configDir:       "with-propertyID",
			checksForCreate: defaultChecker.Build(),
			checksForUpdate: defaultChecker.
				CheckEqual("hostnames.0.cname_to", "to2.test.domain").
				Build(),
		},
		"Lifecycle: latest version is deactivated in staging (normal)": {
			init: func(t *testing.T, p *mockProperty) {
				// set initial data
				p.mockPropertyData = basicDataWithDefaultRules
				p.createFromVersion = 1
				p.newVersionID = 2
				// create
				mockResourcePropertyCreateWithVersionHostnames(p)
				// read x2
				mockResourcePropertyRead(p, 2)
				// mock activating staging version
				p.versions = versionStagingActive
				// read
				mockResourcePropertyRead(p)
				// update creates new version
				p.mockGetPropertyVersion()
				p.mockCreatePropertyVersion()
				// update mock data with new version info and updated hostname
				p.latestVersion = 2
				p.hostnames = updatedHostname
				p.versions = papi.PropertyVersionItems{
					Items: []papi.PropertyVersionGetItem{
						{
							StagingStatus:    papi.VersionStatusActive,
							ProductionStatus: papi.VersionStatusInactive,
							PropertyVersion:  1,
						},
						{
							StagingStatus:    papi.VersionStatusDeactivated,
							ProductionStatus: papi.VersionStatusInactive,
							PropertyVersion:  2,
						},
					},
				}
				p.mockUpdatePropertyVersionHostnames()
				// read x2
				mockResourcePropertyRead(p, 2)
				// delete
				p.mockRemoveProperty()
			},
			configDir:       "normal",
			checksForCreate: defaultChecker.Build(),
			checksForUpdate: defaultChecker.
				CheckEqual("hostnames.0.cname_to", "to2.test.domain").
				CheckEqual("latest_version", "2").
				CheckEqual("staging_version", "1").
				Build(),
		},
		"Lifecycle: latest version is deactivated in production (normal)": {
			init: func(t *testing.T, p *mockProperty) {
				// set initial data
				p.mockPropertyData = basicDataWithDefaultRules
				p.createFromVersion = 1
				p.newVersionID = 2
				// create
				mockResourcePropertyCreateWithVersionHostnames(p)
				// read x2
				mockResourcePropertyRead(p, 2)
				// mock production version active
				p.versions = versionProductionActive
				// read
				mockResourcePropertyRead(p)
				// update creates new version with updated hostname and version info
				p.mockGetPropertyVersion()
				p.mockCreatePropertyVersion()
				p.latestVersion = 2
				p.hostnames = updatedHostname
				p.versions = papi.PropertyVersionItems{
					Items: []papi.PropertyVersionGetItem{
						{
							StagingStatus:    papi.VersionStatusInactive,
							ProductionStatus: papi.VersionStatusActive,
							PropertyVersion:  1,
						},
						{
							StagingStatus:    papi.VersionStatusInactive,
							ProductionStatus: papi.VersionStatusDeactivated,
							PropertyVersion:  2,
						},
					},
				}
				p.mockUpdatePropertyVersionHostnames()
				// read x2
				mockResourcePropertyRead(p, 2)
				// delete
				p.mockRemoveProperty()
			},
			configDir:       "normal",
			checksForCreate: defaultChecker.Build(),
			checksForUpdate: defaultChecker.
				CheckEqual("hostnames.0.cname_to", "to2.test.domain").
				CheckEqual("latest_version", "2").
				CheckEqual("production_version", "1").
				Build(),
		},
		"Lifecycle: latest version is not active (normal)": {
			init: func(t *testing.T, p *mockProperty) {
				// set initial data
				p.mockPropertyData = basicDataWithDefaultRules
				mockLatestVersionNotActive(p)
			},
			configDir:       "normal",
			checksForCreate: defaultChecker.Build(),
			checksForUpdate: defaultChecker.
				CheckEqual("hostnames.0.cname_to", "to2.test.domain").
				Build(),
		},
		"Lifecycle: latest version is active in staging (normal)": {
			init: func(t *testing.T, p *mockProperty) {
				// set initial data
				p.mockPropertyData = basicDataWithDefaultRules
				p.createFromVersion = 1
				p.newVersionID = 2
				mockLatestVersionActiveOnStaging(p)
			},
			configDir: "normal",
			checksForCreate: defaultChecker.
				Build(),
			checksForUpdate: defaultChecker.
				CheckEqual("hostnames.0.cname_to", "to2.test.domain").
				CheckEqual("latest_version", "2").
				CheckEqual("staging_version", "1").
				Build(),
		},
		"Lifecycle: latest version is active in production (normal)": {
			init: func(t *testing.T, p *mockProperty) {
				// set initial data
				p.mockPropertyData = basicDataWithDefaultRules
				p.createFromVersion = 1
				p.newVersionID = 2
				mockLatestVersionActiveOnProduction(p)
			},
			configDir:       "normal",
			checksForCreate: defaultChecker.Build(),
			checksForUpdate: defaultChecker.
				CheckEqual("hostnames.0.cname_to", "to2.test.domain").
				CheckEqual("latest_version", "2").
				CheckEqual("production_version", "1").
				Build(),
		},
		"Lifecycle: latest version is not active (contract_id without prefix)": {
			init: func(t *testing.T, p *mockProperty) {
				// set initial data
				p.mockPropertyData = basicDataWithDefaultRules
				mockLatestVersionNotActive(p)
			},
			configDir:       "contract_id without prefix",
			checksForCreate: defaultChecker.Build(),
			checksForUpdate: defaultChecker.
				CheckEqual("hostnames.0.cname_to", "to2.test.domain").
				Build(),
		},
		"Lifecycle: latest version active in staging (contract_id without prefix)": {
			init: func(t *testing.T, p *mockProperty) {
				// set initial data
				p.mockPropertyData = basicDataWithDefaultRules
				p.createFromVersion = 1
				p.newVersionID = 2
				mockLatestVersionActiveOnStaging(p)
			},
			configDir:       "contract_id without prefix",
			checksForCreate: defaultChecker.Build(),
			checksForUpdate: defaultChecker.
				CheckEqual("hostnames.0.cname_to", "to2.test.domain").
				CheckEqual("latest_version", "2").
				CheckEqual("staging_version", "1").
				Build(),
		},
		"Lifecycle: latest version active in production (contract_id without prefix)": {
			init: func(t *testing.T, p *mockProperty) {
				// set initial data
				p.mockPropertyData = basicDataWithDefaultRules
				p.createFromVersion = 1
				p.newVersionID = 2
				mockLatestVersionActiveOnProduction(p)
			},
			configDir:       "contract_id without prefix",
			checksForCreate: defaultChecker.Build(),
			checksForUpdate: defaultChecker.
				CheckEqual("hostnames.0.cname_to", "to2.test.domain").
				CheckEqual("latest_version", "2").
				CheckEqual("production_version", "1").
				Build(),
		},
		"Lifecycle: latest version is not active (group_id without prefix)": {
			init: func(t *testing.T, p *mockProperty) {
				// set initial data
				p.mockPropertyData = basicDataWithDefaultRules
				mockLatestVersionNotActive(p)
			},
			configDir:       "group_id without prefix",
			checksForCreate: defaultChecker.Build(),
			checksForUpdate: defaultChecker.
				CheckEqual("hostnames.0.cname_to", "to2.test.domain").
				Build(),
		},
		"Lifecycle: latest version is active in staging (group_id without prefix)": {
			init: func(t *testing.T, p *mockProperty) {
				// set initial data
				p.mockPropertyData = basicDataWithDefaultRules
				p.createFromVersion = 1
				p.newVersionID = 2
				mockLatestVersionActiveOnStaging(p)
			},
			configDir:       "group_id without prefix",
			checksForCreate: defaultChecker.Build(),
			checksForUpdate: defaultChecker.
				CheckEqual("hostnames.0.cname_to", "to2.test.domain").
				CheckEqual("latest_version", "2").
				CheckEqual("staging_version", "1").
				Build(),
		},
		"Lifecycle: latest version is active in production (group_id without prefix)": {
			init: func(t *testing.T, p *mockProperty) {
				// set initial data
				p.mockPropertyData = basicDataWithDefaultRules
				p.createFromVersion = 1
				p.newVersionID = 2
				mockLatestVersionActiveOnProduction(p)
			},
			configDir:       "group_id without prefix",
			checksForCreate: defaultChecker.Build(),
			checksForUpdate: defaultChecker.
				CheckEqual("hostnames.0.cname_to", "to2.test.domain").
				CheckEqual("latest_version", "2").
				CheckEqual("production_version", "1").
				Build(),
		},
		"Lifecycle: latest version is not active (product_id without prefix)": {
			init: func(t *testing.T, p *mockProperty) {
				// set initial data
				p.mockPropertyData = basicDataWithDefaultRules
				p.createFromVersion = 1
				p.newVersionID = 2
				mockLatestVersionNotActive(p)
			},
			configDir:       "product_id without prefix",
			checksForCreate: defaultChecker.Build(),
			checksForUpdate: defaultChecker.
				CheckEqual("hostnames.0.cname_to", "to2.test.domain").
				Build(),
		},
		"Lifecycle: latest version is active in staging (product_id without prefix)": {
			init: func(t *testing.T, p *mockProperty) {
				// set initial data
				p.mockPropertyData = basicDataWithDefaultRules
				p.createFromVersion = 1
				p.newVersionID = 2
				mockLatestVersionActiveOnStaging(p)
			},
			configDir:       "product_id without prefix",
			checksForCreate: defaultChecker.Build(),
			checksForUpdate: defaultChecker.
				CheckEqual("hostnames.0.cname_to", "to2.test.domain").
				CheckEqual("latest_version", "2").
				CheckEqual("staging_version", "1").
				Build(),
		},
		"Lifecycle: latest version is active in production (product_id without prefix)": {
			init: func(t *testing.T, p *mockProperty) {
				// set initial data
				p.mockPropertyData = basicDataWithDefaultRules
				p.createFromVersion = 1
				p.newVersionID = 2
				mockLatestVersionActiveOnProduction(p)
			},
			configDir:       "product_id without prefix",
			checksForCreate: defaultChecker.Build(),
			checksForUpdate: defaultChecker.
				CheckEqual("hostnames.0.cname_to", "to2.test.domain").
				CheckEqual("latest_version", "2").
				CheckEqual("production_version", "1").
				Build(),
		},
		"Lifecycle: no diff": {
			init: func(t *testing.T, p *mockProperty) {
				// set initial data
				p.mockPropertyData = basicData
				p.ruleTree = mockRuleTreeData{
					rules: papi.Rules{
						Children: []papi.Rules{
							{
								Name:                "Default CORS Policy",
								CriteriaMustSatisfy: papi.RuleCriteriaMustSatisfyAll,
							},
						},
					},
				}
				// create
				mockResourcePropertyFullCreate(p)
				// read x2
				mockResourcePropertyRead(p, 2)
				// mock update in rules by changing the order
				p.ruleTree.rules = papi.Rules{Children: []papi.Rules{{CriteriaMustSatisfy: papi.RuleCriteriaMustSatisfyAll, Name: "Default CORS Policy"}}}
				// read x2
				mockResourcePropertyRead(p, 2)
				// delete
				p.mockRemoveProperty()
			},
			configDir: "no diff",
			checksForCreate: defaultChecker.
				CheckEqual("rules", `{"rules":{"children":[{"name":"Default CORS Policy","options":{},"criteriaMustSatisfy":"all"}],"name":"","options":{}}}`).
				Build(),
			checksForUpdate: defaultChecker.
				CheckEqual("rules", `{"rules":{"children":[{"name":"Default CORS Policy","options":{},"criteriaMustSatisfy":"all"}],"name":"","options":{}}}`).
				Build(),
		},
		"Lifecycle: rules custom diff": {
			/*
				rulesCustomDiff tests rulesCustomDiff function which is in resource_akamai_property.go file.
				There is an additional field "options":{} in expected attributes, because with UpdateRuleTree(ctx, req) function
				this field added automatically into response, even if it does not exist in rules.
			*/
			init: func(t *testing.T, p *mockProperty) {
				// set initial data
				p.mockPropertyData = basicData
				p.ruleTree = mockRuleTreeData{
					rules: papi.Rules{Behaviors: []papi.RuleBehavior{
						{
							Name: "caching",
							Options: papi.RuleOptionsMap{
								"behavior":       "MAX_AGE",
								"mustRevalidate": false,
								"ttl":            "12d",
							},
						},
					},
						Name: "default"},
				}
				// create
				mockResourcePropertyFullCreate(p)
				// read x2
				mockResourcePropertyRead(p, 2)
				// read
				mockResourcePropertyRead(p)
				// update ttl in rule tree from 12d to 13d
				p.ruleTree.rules.Behaviors = []papi.RuleBehavior{
					{
						Name: "caching",
						Options: papi.RuleOptionsMap{
							"behavior":       "MAX_AGE",
							"mustRevalidate": false,
							"ttl":            "13d",
						},
					},
				}
				// update
				p.mockGetPropertyVersion()
				p.mockUpdateRuleTree()
				// read x2
				mockResourcePropertyRead(p, 2)
				// delete
				p.mockRemoveProperty()
			},
			configDir: "rules custom diff",
			checksForCreate: defaultChecker.
				CheckEqual("rules", `{"rules":{"behaviors":[{"name":"caching","options":{"behavior":"MAX_AGE","mustRevalidate":false,"ttl":"12d"}}],"name":"default","options":{}}}`).
				Build(),
			checksForUpdate: defaultChecker.
				CheckEqual("rules", `{"rules":{"behaviors":[{"name":"caching","options":{"behavior":"MAX_AGE","mustRevalidate":false,"ttl":"13d"}}],"name":"default","options":{}}}`).
				Build(),
		},
		"Lifecycle: no diff for hostnames (hostnames)": {
			init: func(t *testing.T, p *mockProperty) {
				// set initial data
				p.mockPropertyData = basicData
				p.ruleTree = mockRuleTreeData{
					rules: papi.Rules{
						Children: []papi.Rules{
							{
								CriteriaMustSatisfy: papi.RuleCriteriaMustSatisfyAll,
								Name:                "Default CORS Policy",
							},
						},
					},
				}
				p.hostnames = papi.HostnameResponseItems{
					Items: []papi.Hostname{
						{
							CnameType:            "EDGE_HOSTNAME",
							CnameFrom:            "from1.test.domain",
							CnameTo:              "to1.test.domain",
							CertProvisioningType: "DEFAULT",
							EdgeHostnameID:       "ehn_123",
						},
						{
							CnameType:            "EDGE_HOSTNAME",
							CnameFrom:            "from2.test.domain",
							CnameTo:              "to2.test.domain",
							CertProvisioningType: "DEFAULT",
							EdgeHostnameID:       "ehn_123",
						},
					},
				}
				// create
				mockResourcePropertyFullCreate(p)
				// read x2
				mockResourcePropertyRead(p, 2)
				// read x2 - refresh as nothing other than the order of hostnames changed
				mockResourcePropertyRead(p, 2)
				// delete
				p.mockRemoveProperty()
			},
			configDir: "hostnames",
			checksForCreate: defaultChecker.
				CheckEqual("hostnames.0.cname_to", "to1.test.domain").
				CheckEqual("hostnames.1.cname_to", "to2.test.domain").
				CheckEqual("rules", `{"rules":{"children":[{"name":"Default CORS Policy","options":{},"criteriaMustSatisfy":"all"}],"name":"","options":{}}}`).
				Build(),
			checksForUpdate: defaultChecker.
				CheckEqual("hostnames.0.cname_to", "to1.test.domain").
				CheckEqual("hostnames.1.cname_to", "to2.test.domain").
				CheckEqual("rules", `{"rules":{"children":[{"name":"Default CORS Policy","options":{},"criteriaMustSatisfy":"all"}],"name":"","options":{}}}`).
				Build(),
		},
		"Lifecycle: rules with variables": {
			init: func(t *testing.T, p *mockProperty) {
				// set initial data
				p.mockPropertyData = basicData
				p.ruleTree = mockRuleTreeData{
					rules: papi.Rules{
						Name: "default",
						Children: []papi.Rules{
							{
								Name: "change fwd path",
								Behaviors: []papi.RuleBehavior{
									{
										Name: "baseDirectory",
										Options: papi.RuleOptionsMap{
											"value": "/smth/",
										},
									},
								},
								Criteria: []papi.RuleBehavior{
									{
										Name:   "requestHeader",
										Locked: false,
										Options: papi.RuleOptionsMap{
											"headerName":              "Accept-Encoding",
											"matchCaseSensitiveValue": true,
											"matchOperator":           "IS_ONE_OF",
											"matchWildcardName":       false,
											"matchWildcardValue":      false,
										},
									},
								},
								CriteriaMustSatisfy: papi.RuleCriteriaMustSatisfyAll,
							},
							{
								Name: "caching",
								Behaviors: []papi.RuleBehavior{
									{
										Name: "caching",
										Options: papi.RuleOptionsMap{
											"behavior":       "MAX_AGE",
											"mustRevalidate": false,
											"ttl":            "1m",
										},
									},
								},
								CriteriaMustSatisfy: papi.RuleCriteriaMustSatisfyAny,
							},
						},
						Behaviors: []papi.RuleBehavior{
							{
								Name: "origin",
								Options: papi.RuleOptionsMap{
									"cacheKeyHostname":          "REQUEST_HOST_HEADER",
									"compress":                  true,
									"enableTrueClientIp":        true,
									"forwardHostHeader":         "REQUEST_HOST_HEADER",
									"hostname":                  "test.domain",
									"httpPort":                  float64(80),
									"httpsPort":                 float64(443),
									"originCertificate":         "",
									"originSni":                 true,
									"originType":                "CUSTOMER",
									"ports":                     "",
									"trueClientIpClientSetting": false,
									"trueClientIpHeader":        "True-Client-IP",
									"verificationMode":          "PLATFORM_SETTINGS",
								},
							},
						},
						Options: papi.RuleOptions{},
						Variables: []papi.RuleVariable{
							{
								Name:        "TEST_EMPTY_FIELDS",
								Value:       ptr.To(""),
								Description: ptr.To(""),
								Hidden:      true,
								Sensitive:   false,
							},
							{
								Name:        "TEST_NIL_FIELD",
								Description: nil,
								Value:       ptr.To(""),
								Hidden:      true,
								Sensitive:   false,
							},
						},
						Comments: "The behaviors in the Default Rule apply to all requests for the property hostname(s) unless another rule overrides the Default Rule settings.",
					},
				}
				// create
				mockResourcePropertyFullCreate(p)
				// read x2
				mockResourcePropertyRead(p, 2)
				// read x1 before update
				mockResourcePropertyRead(p)
				// update with new rules variables (description in TEST_NIL_FIELD is "", was nil)
				p.mockGetPropertyVersion()
				p.ruleTree.rules.Variables = []papi.RuleVariable{
					{
						Name:        "TEST_EMPTY_FIELDS",
						Value:       ptr.To(""),
						Description: ptr.To(""),
						Hidden:      true,
						Sensitive:   false,
					},
					{
						Name:        "TEST_NIL_FIELD",
						Description: ptr.To(""),
						Value:       ptr.To(""),
						Hidden:      true,
						Sensitive:   false,
					},
				}
				p.mockUpdateRuleTree()
				// read x2
				mockResourcePropertyRead(p, 2)
				// delete
				p.mockRemoveProperty()
			},
			configDir: "rules with variables",
			checksForCreate: defaultChecker.
				CheckEqual("rules", `{"rules":{"behaviors":[{"name":"origin","options":{"cacheKeyHostname":"REQUEST_HOST_HEADER","compress":true,"enableTrueClientIp":true,"forwardHostHeader":"REQUEST_HOST_HEADER","hostname":"test.domain","httpPort":80,"httpsPort":443,"originCertificate":"","originSni":true,"originType":"CUSTOMER","ports":"","trueClientIpClientSetting":false,"trueClientIpHeader":"True-Client-IP","verificationMode":"PLATFORM_SETTINGS"}}],"children":[{"behaviors":[{"name":"baseDirectory","options":{"value":"/smth/"}}],"criteria":[{"name":"requestHeader","options":{"headerName":"Accept-Encoding","matchCaseSensitiveValue":true,"matchOperator":"IS_ONE_OF","matchWildcardName":false,"matchWildcardValue":false}}],"name":"change fwd path","options":{},"criteriaMustSatisfy":"all"},{"behaviors":[{"name":"caching","options":{"behavior":"MAX_AGE","mustRevalidate":false,"ttl":"1m"}}],"name":"caching","options":{},"criteriaMustSatisfy":"any"}],"comments":"The behaviors in the Default Rule apply to all requests for the property hostname(s) unless another rule overrides the Default Rule settings.","name":"default","options":{},"variables":[{"description":"","hidden":true,"name":"TEST_EMPTY_FIELDS","sensitive":false,"value":""},{"description":null,"hidden":true,"name":"TEST_NIL_FIELD","sensitive":false,"value":""}]}}`).
				Build(),
			checksForUpdate: defaultChecker.
				CheckEqual("rules", `{"rules":{"behaviors":[{"name":"origin","options":{"cacheKeyHostname":"REQUEST_HOST_HEADER","compress":true,"enableTrueClientIp":true,"forwardHostHeader":"REQUEST_HOST_HEADER","hostname":"test.domain","httpPort":80,"httpsPort":443,"originCertificate":"","originSni":true,"originType":"CUSTOMER","ports":"","trueClientIpClientSetting":false,"trueClientIpHeader":"True-Client-IP","verificationMode":"PLATFORM_SETTINGS"}}],"children":[{"behaviors":[{"name":"baseDirectory","options":{"value":"/smth/"}}],"criteria":[{"name":"requestHeader","options":{"headerName":"Accept-Encoding","matchCaseSensitiveValue":true,"matchOperator":"IS_ONE_OF","matchWildcardName":false,"matchWildcardValue":false}}],"name":"change fwd path","options":{},"criteriaMustSatisfy":"all"},{"behaviors":[{"name":"caching","options":{"behavior":"MAX_AGE","mustRevalidate":false,"ttl":"1m"}}],"name":"caching","options":{},"criteriaMustSatisfy":"any"}],"comments":"The behaviors in the Default Rule apply to all requests for the property hostname(s) unless another rule overrides the Default Rule settings.","name":"default","options":{},"variables":[{"description":"","hidden":true,"name":"TEST_EMPTY_FIELDS","sensitive":false,"value":""},{"description":"","hidden":true,"name":"TEST_NIL_FIELD","sensitive":false,"value":""}]}}`).
				Build(),
		},
		"Lifecycle: Verify staging_version and production_version known at plan": {
			init: func(t *testing.T, p *mockProperty) {
				// set initial data
				p.mockPropertyData = basicDataWithDefaultRules
				// create
				mockResourcePropertyCreateWithVersionHostnames(p)
				// read x2
				mockResourcePropertyRead(p, 2)
				// read x1 before update
				mockResourcePropertyRead(p)
				// update
				p.mockGetPropertyVersion()
				p.hostnames = updatedHostname
				p.mockUpdatePropertyVersionHostnames()
				// read x2
				mockResourcePropertyRead(p, 2)
				// delete
				p.mockRemoveProperty()
			},
			configDir:       "normal",
			checksForCreate: defaultChecker.Build(),
			checksForUpdate: defaultChecker.
				CheckEqual("hostnames.0.cname_to", "to2.test.domain").
				Build(),
			configPlanChecks: resource.ConfigPlanChecks{
				PreApply: []plancheck.PlanCheck{
					plancheck.ExpectKnownValue("akamai_property.test", tfjsonpath.New("staging_version"), knownvalue.Int64Exact(0)),
					plancheck.ExpectKnownValue("akamai_property.test", tfjsonpath.New("production_version"), knownvalue.Int64Exact(0)),
					plancheck.ExpectUnknownValue("akamai_property.test", tfjsonpath.New("latest_version")),
				},
			},
		},
		"Lifecycle: update group id - in place": {
			init: func(t *testing.T, p *mockProperty) {
				// set initial data
				p.mockPropertyData = basicDataWithDefaultRules
				p.moveGroup = moveGroup{
					sourceGroupID:      2,
					destinationGroupID: 222,
				}
				// create
				mockResourcePropertyCreateWithVersionHostnames(p)
				// read x2
				mockResourcePropertyRead(p, 2)
				// read x1 before update
				mockResourcePropertyRead(p)
				// update - moving the property
				p.mockMoveProperty()
				p.groupID = "grp_222"
				p.mockGetProperty()
				// read from update
				mockResourcePropertyRead(p)
				// read before delete
				mockResourcePropertyRead(p)
				// delete
				p.mockRemoveProperty()
			},
			configDir:       "groupIDUpdate",
			checksForCreate: defaultChecker.Build(),
			checksForUpdate: defaultChecker.
				CheckEqual("group_id", "grp_222").
				Build(),
		},
		"Lifecycle: update group id and hostnames - in place": {
			init: func(t *testing.T, p *mockProperty) {
				// set initial data
				p.mockPropertyData = basicDataWithDefaultRules
				p.moveGroup = moveGroup{
					sourceGroupID:      2,
					destinationGroupID: 222,
				}
				// create
				mockResourcePropertyCreateWithVersionHostnames(p)
				// read x2
				mockResourcePropertyRead(p, 2)
				// read x1 before update
				mockResourcePropertyRead(p)
				// update - moving the property
				p.mockMoveProperty()
				p.groupID = "grp_222"
				// waiting for new groupID
				p.mockGetProperty()
				// readout for general version calculations
				p.mockGetPropertyVersion()
				// change in hostnames detected
				p.hostnames = papi.HostnameResponseItems{
					Items: []papi.Hostname{
						{
							CnameFrom:            "from2.test.domain",
							CnameTo:              "to.test.domain",
							CertProvisioningType: "DEFAULT",
							CnameType:            "EDGE_HOSTNAME",
							EdgeHostnameID:       "ehn_123",
						},
					},
				}
				p.mockUpdatePropertyVersionHostnames()
				// read x2
				mockResourcePropertyRead(p, 2)
				// delete
				p.mockRemoveProperty().Once()
			},
			configDir:       "groupIDUpdate/withHostnames",
			checksForCreate: defaultChecker.Build(),
			checksForUpdate: defaultChecker.
				CheckEqual("group_id", "grp_222").
				Build(),
		},
		"Lifecycle: update group id and name - recreate": {
			init: func(t *testing.T, p *mockProperty) {
				// set initial data
				p.mockPropertyData = basicDataWithDefaultRules
				p.moveGroup = moveGroup{
					sourceGroupID:      2,
					destinationGroupID: 222,
				}
				// create
				mockResourcePropertyCreateWithVersionHostnames(p)
				// read x2
				mockResourcePropertyRead(p, 2)
				// read x1 before update
				mockResourcePropertyRead(p)
				p.mockRemoveProperty().Once()
				// recreate the resource
				p.propertyName = "dummy_name2"
				p.groupID = "grp_222"
				// recreate new property
				mockResourcePropertyCreateWithVersionHostnames(p)
				// read x2
				mockResourcePropertyRead(p, 2)
				// delete
				p.mockRemoveProperty().Once()
			},
			configDir:       "groupIDUpdate/withName",
			checksForCreate: defaultChecker.Build(),
			checksForUpdate: defaultChecker.
				CheckEqual("group_id", "grp_222").
				CheckEqual("name", "dummy_name2").
				Build(),
		},
		"Lifecycle: update group id: forwarding property read API error": {
			init: func(t *testing.T, p *mockProperty) {
				// set initial data
				p.mockPropertyData = basicDataWithDefaultRules
				p.moveGroup = moveGroup{
					sourceGroupID:      2,
					destinationGroupID: 222,
				}
				// create
				mockResourcePropertyCreateWithVersionHostnames(p)
				// read x2
				mockResourcePropertyRead(p, 2)
				// read x1 before update
				mockResourcePropertyRead(p)

				// Error checking if the property is already in the dst group
				p.groupID = "grp_222"
				req := p.getPropertyRequest()
				err := errors.New("read err")
				p.papiMock.On("GetProperty", testutils.MockContext, req).Return(nil, err).Once()

				p.groupID = "grp_2"
				p.mockRemoveProperty()
			},
			configDir:       "groupIDUpdate",
			checksForCreate: defaultChecker.Build(),
			updateError: regexp.MustCompile(
				"error moving property: error checking if property in group: unexpected http error for {prp_4 grp_222 ctr_1}: read err"),
		},
		"Lifecycle: update group id: no API call to move if property already in desired group": {
			init: func(t *testing.T, p *mockProperty) {
				// set initial data
				p.mockPropertyData = basicDataWithDefaultRules
				p.moveGroup = moveGroup{
					sourceGroupID:      2,
					destinationGroupID: 222,
				}
				// create
				mockResourcePropertyCreateWithVersionHostnames(p)
				// read x2
				mockResourcePropertyRead(p, 2)
				// read x1 before update
				mockResourcePropertyRead(p)

				// update - moving the property
				// yes, property is in the desired group, exiting early from moveProperty
				p.groupID = "grp_222"
				p.mockGetProperty()
				// read from update
				mockResourcePropertyRead(p)
				// read before delete
				mockResourcePropertyRead(p)
				// delete
				p.mockRemoveProperty()
			},
			configDir:       "groupIDUpdate",
			checksForCreate: defaultChecker.Build(),
			checksForUpdate: defaultChecker.
				CheckEqual("group_id", "grp_222").
				Build(),
		},
		"Lifecycle: update group id: moving properties with no past activations is not supported": {
			init: func(t *testing.T, p *mockProperty) {
				// set initial data
				p.mockPropertyData = basicDataWithDefaultRules
				p.moveGroup = moveGroup{
					sourceGroupID:      2,
					destinationGroupID: 222,
				}
				// create
				mockResourcePropertyCreateWithVersionHostnames(p)
				// read x2
				mockResourcePropertyRead(p, 2)
				// read x1 before update
				mockResourcePropertyRead(p)

				// Checking if the property is already in the dst group
				p.groupID = "grp_222"
				getReq := p.getPropertyRequest()
				p.papiMock.On("GetProperty", testutils.MockContext, getReq).
					Return(nil, &papi.Error{StatusCode: http.StatusForbidden}).
					Once()

				p.groupID = "grp_2"
				p.mockGetActivationsCompleteRequest()

				// delete
				p.mockRemoveProperty()
			},
			configDir:       "groupIDUpdate",
			checksForCreate: defaultChecker.Build(),
			updateError: regexp.MustCompile("error moving property: " +
				"moving properties that have never been activated is not supported " +
				`\(property id: prp_4, contract id: ctr_1, group id grp_2\)`),
		},
		"Lifecycle: update group id: forwards API error from waiting for group id change": {
			init: func(t *testing.T, p *mockProperty) {
				// set initial data
				p.mockPropertyData = basicDataWithDefaultRules
				p.moveGroup = moveGroup{
					sourceGroupID:      2,
					destinationGroupID: 222,
				}
				// create
				mockResourcePropertyCreateWithVersionHostnames(p)
				// read x2
				mockResourcePropertyRead(p, 2)
				// read x1 before update
				mockResourcePropertyRead(p)
				// update - moving the property
				p.mockMoveProperty()

				// Waiting for group id change
				p.groupID = "grp_222"
				getReq := p.getPropertyRequest()
				p.papiMock.On("GetProperty", testutils.MockContext, getReq).
					Return(nil, errors.New("read err")).
					Once()

				// delete
				p.groupID = "grp_2"
				p.mockRemoveProperty()
			},
			configDir:       "groupIDUpdate",
			checksForCreate: defaultChecker.Build(),
			updateError: regexp.MustCompile(
				"error moving property: error waiting for group id change: unexpected http error for {prp_4 grp_222 ctr_1}: read err"),
		},
	}

	for name, test := range tests {
		t.Run(name, func(t *testing.T) {
			papiMock := &papi.Mock{}
			iamMock := &iam.Mock{}
			mp := mockProperty{
				papiMock: papiMock,
				iamMock:  iamMock,
			}
			test.init(t, &mp)

			useClient(papiMock, nil, func() {
				useIam(iamMock, func() {
					resource.UnitTest(t, resource.TestCase{
						ProtoV6ProviderFactories: testutils.NewProtoV6ProviderFactory(NewSubprovider()),
						Steps: []resource.TestStep{
							{
								Config: testutils.LoadFixtureString(t, "testdata/TestResProperty/Lifecycle/%s/step0.tf", test.configDir),
								Check:  test.checksForCreate,
							},
							{
								Config:           testutils.LoadFixtureString(t, "testdata/TestResProperty/Lifecycle/%s/step1.tf", test.configDir),
								Check:            test.checksForUpdate,
								ConfigPlanChecks: test.configPlanChecks,
								ExpectError:      test.updateError,
							},
						},
					})
				})
			})

			papiMock.AssertExpectations(t)
		})
	}

	// separate tests as they require different number of steps or filenames
	t.Run("Lifecycle: diff cpCode", func(t *testing.T) {
		papiMock := &papi.Mock{}
		mp := &mockProperty{
			papiMock:         papiMock,
			mockPropertyData: basicData,
		}
		mp.ruleTree = mockRuleTreeData{
			rules: papi.Rules{
				Behaviors: []papi.RuleBehavior{
					{
						Name: "cpCode",
						Options: papi.RuleOptionsMap{
							"value": map[string]interface{}{
								"description": "CliTerraformCPCode",
								"id":          1.050269e+06,
								"name":        "DevExpCliTerraformPapiAsSchemaTest",
								"products":    []interface{}{"Web_App_Accel"},
							},
						},
					},
				},
				Name: "default",
			},
		}
		// create
		mockResourcePropertyFullCreate(mp)
		// mock rules in the format that API returns to test custom diff functionality on rules (notice `"cpCodeLimits": nil`, which was not present in the request.
		mp.ruleTree.rules = papi.Rules{Behaviors: []papi.RuleBehavior{
			{
				Name: "cpCode",
				Options: papi.RuleOptionsMap{
					"value": map[string]interface{}{
						"cpCodeLimits": nil,
						"description":  "CliTerraformCPCode",
						"id":           1.050269e+06,
						"name":         "DevExpCliTerraformPapiAsSchemaTest",
						"products":     []interface{}{"Web_App_Accel"},
					},
				},
			},
		},
			Name: "default"}
		// read x2
		mockResourcePropertyRead(mp, 2)
		// delete
		mp.mockRemoveProperty()

		useClient(papiMock, nil, func() {
			resource.UnitTest(t, resource.TestCase{
				ProtoV6ProviderFactories: testutils.NewProtoV6ProviderFactory(NewSubprovider()),
				Steps: []resource.TestStep{
					{
						Config: testutils.LoadFixtureString(t, "testdata/TestResProperty/Lifecycle/rules diff cpcode/step0.tf"),
						Check: defaultChecker.
							CheckEqual("rules", `{"rules":{"behaviors":[{"name":"cpCode","options":{"value":{"cpCodeLimits":null,"description":"CliTerraformCPCode","id":1050269,"name":"DevExpCliTerraformPapiAsSchemaTest","products":["Web_App_Accel"]}}}],"name":"default","options":{}}}`).
							Build(),
					},
				},
			})
		})
	})

	t.Run("Lifecycle: new version changed on server", func(t *testing.T) {
		papiMock := &papi.Mock{}
		mp := &mockProperty{
			papiMock:         papiMock,
			mockPropertyData: basicDataWithDefaultRules,
		}
		// create
		mockResourcePropertyCreateWithVersionHostnames(mp)
		// read x2
		mockResourcePropertyRead(mp, 2)
		// simulate remote change outside terraform, only for response data - new version with updated CnameTo.
		mp.latestVersion = 2
		mp.hostnames = papi.HostnameResponseItems{
			Items: []papi.Hostname{
				{
					CnameType:            "EDGE_HOSTNAME",
					CnameFrom:            "from.test.domain",
					CnameTo:              "changed.test.domain",
					CertProvisioningType: "DEFAULT",
					EdgeHostnameID:       "ehn_123",
				},
			},
		}
		mp.versions = papi.PropertyVersionItems{
			Items: []papi.PropertyVersionGetItem{
				{
					StagingStatus:    papi.VersionStatusInactive,
					ProductionStatus: papi.VersionStatusInactive,
					PropertyVersion:  2,
				},
			},
		}
		// read x1 - remote, updated state
		mockResourcePropertyRead(mp)
		// update
		mp.mockGetPropertyVersion()
		// such drift should invoke update function, which should use value from config which should replace the remote value.
		// Hence, CnameTo is assigned the value from config for the mock data.
		mp.hostnames = papi.HostnameResponseItems{
			Items: []papi.Hostname{
				{
					CnameType:            "EDGE_HOSTNAME",
					CnameFrom:            "from.test.domain",
					CnameTo:              "to.test.domain",
					CertProvisioningType: "DEFAULT",
					EdgeHostnameID:       "ehn_123",
				},
			},
		}
		mp.mockUpdatePropertyVersionHostnames()
		// read x2
		mockResourcePropertyRead(mp, 2)
		// delete
		mp.mockRemoveProperty()

		useClient(papiMock, nil, func() {
			resource.UnitTest(t, resource.TestCase{
				ProtoV6ProviderFactories: testutils.NewProtoV6ProviderFactory(NewSubprovider()),
				Steps: []resource.TestStep{
					{
						Config: testutils.LoadFixtureString(t, "testdata/TestResProperty/Lifecycle/new version changed on server/step0.tf"),
						Check:  defaultChecker.Build(),
					},
					{
						Config: testutils.LoadFixtureString(t, "testdata/TestResProperty/Lifecycle/new version changed on server/step0.tf"),
						Check: defaultChecker.
							CheckEqual("latest_version", "2").
							Build(),
					},
				},
			})
		})
	})

	t.Run("Lifecycle: no diff for CriteriaMustSatisfy", func(t *testing.T) {
		// set initial data
		papiMock := &papi.Mock{}
		mp := &mockProperty{
			papiMock:         papiMock,
			mockPropertyData: basicData,
		}
		mp.ruleTree = mockRuleTreeData{
			rules: papi.Rules{
				Children: []papi.Rules{
					{
						Name: "Default CORS Policy",
					},
				},
			},
		}
		// create
		mp.mockCreateProperty()
		mp.mockUpdatePropertyVersionHostnames()

		rulesUpdate := papi.RulesUpdate{
			Rules:    mp.ruleTree.rules,
			Comments: mp.ruleTree.comments,
		}
		req := papi.UpdateRulesRequest{
			PropertyID:      mp.propertyID,
			PropertyVersion: mp.latestVersion,
			ContractID:      mp.contractID,
			GroupID:         mp.groupID,
			Rules:           rulesUpdate,
			ValidateRules:   true,
		}

		// default child setting of the CriteriaMustSatisfy field to the value "all" by API if it is not provided
		defaultAPIResp := papi.Rules{
			AdvancedOverride: mp.ruleTree.rules.AdvancedOverride,
			Behaviors:        mp.ruleTree.rules.Behaviors,
			Children: []papi.Rules{
				{
					Name:                mp.ruleTree.rules.Children[0].Name,
					CriteriaMustSatisfy: papi.RuleCriteriaMustSatisfyAll,
				},
			},
			Comments:            mp.ruleTree.rules.Comments,
			Criteria:            mp.ruleTree.rules.Criteria,
			CriteriaLocked:      mp.ruleTree.rules.CriteriaLocked,
			CustomOverride:      mp.ruleTree.rules.CustomOverride,
			Name:                mp.ruleTree.rules.Name,
			Options:             mp.ruleTree.rules.Options,
			UUID:                mp.ruleTree.rules.UUID,
			TemplateUuid:        mp.ruleTree.rules.TemplateUuid,
			TemplateLink:        mp.ruleTree.rules.TemplateLink,
			Variables:           mp.ruleTree.rules.Variables,
			CriteriaMustSatisfy: mp.ruleTree.rules.CriteriaMustSatisfy,
		}
		resp := papi.UpdateRulesResponse{
			PropertyID:      mp.propertyID,
			ContractID:      mp.contractID,
			GroupID:         mp.groupID,
			PropertyVersion: mp.latestVersion,
			RuleFormat:      mp.ruleTree.ruleFormat,
			Rules:           defaultAPIResp,
			Errors:          mp.ruleTree.ruleErrors,
			Warnings:        mp.ruleTree.ruleWarnings,
		}
		mp.papiMock.On("UpdateRuleTree", testutils.MockContext, req).Return(&resp, nil).Once()

		// state update after rule tree update
		mp.ruleTree.rules = defaultAPIResp

		// read x2
		mockResourcePropertyRead(mp, 2)
		// delete
		mp.mockRemoveProperty()

		useClient(papiMock, nil, func() {
			resource.UnitTest(t, resource.TestCase{
				ProtoV6ProviderFactories: testutils.NewProtoV6ProviderFactory(NewSubprovider()),
				Steps: []resource.TestStep{
					{
						Config: testutils.LoadFixtureString(t, "testdata/TestResProperty/Lifecycle/criteriaMustSatisfyNoDiff/step0.tf"),
						Check: defaultChecker.
							CheckEqual("rules", `{"rules":{"children":[{"name":"Default CORS Policy","options":{},"criteriaMustSatisfy":"all"}],"name":"","options":{}}}`).
							Build(),
					},
				},
			})
		})
	})
}

// TestPropertyImports tests import functionality of property resource
func TestPropertyImport(t *testing.T) {
	// Based on importID, different API calls are being made in the Import and Read functions. If the importID allows to
	// reconcile specific property version, GetProperty calls are being executed in Import and Read operations. If the property version
	// is unknown, GetPropertyVersions is being used instead of GetProperty.

	// mockPropertyImportKnownVersion gathers API calls that are being executed when property version is known.
	// Uses GetProperty calls.
	mockPropertyImportKnownVersion := func(p *mockProperty) {
		// import
		p.mockGetProperty()
		// read
		mockResourcePropertyRead(p)
	}

	// mockPropertyImportKnownVersion gathers API calls that are being executed when property version is unknown.
	// Uses GetPropertyVersions calls.
	mockPropertyImportUnknownVersion := func(p *mockProperty) {
		// import
		p.mockGetPropertyVersions()
		// read
		p.mockGetPropertyVersions()
		p.mockGetPropertyVersionHostnames()
		p.mockGetRuleTree()
		p.mockGetPropertyVersion()
	}

	// mockPropertyImportKnownVersionAfterImport gather API calls that are being executed when property version is provided and known
	// read function, but not during import: hence there is a single call to GetPropertyVersions in import.
	mockPropertyImportKnownVersionAfterImport := func(p *mockProperty) {
		// import
		p.mockGetPropertyVersions()
		// read (notice that here one GetPropertyVersions call is omitted)
		p.mockGetPropertyVersionHostnames()
		p.mockGetRuleTree()
		p.mockGetPropertyVersion()
	}

	// basicData holds basic, common data across test cases
	basicData := mockPropertyData{
		propertyID:    "prp_4",
		groupID:       "grp_2",
		contractID:    "ctr_1",
		latestVersion: 1,
		hostnames: papi.HostnameResponseItems{
			Items: []papi.Hostname{
				{
					CnameTo:        "to.test.domain",
					EdgeHostnameID: "ehn_123",
				},
			},
		},
		versions: papi.PropertyVersionItems{
			Items: []papi.PropertyVersionGetItem{
				{
					StagingStatus:    papi.VersionStatusActive,
					ProductionStatus: papi.VersionStatusInactive,
					PropertyVersion:  1,
				},
			},
		},
	}

	// basicDataWithoutGroupAndContract does not contain group and contract parameters for cases where they are not part of importID
	basicDataWithoutGroupAndContract := basicData
	basicDataWithoutGroupAndContract.groupID = ""
	basicDataWithoutGroupAndContract.contractID = ""

	// defaultChecker builds basic, common checks across test cases
	defaultChecker := test.NewImportChecker().
		CheckEqual("id", "prp_4").
		CheckEqual("hostnames.0.cname_to", "to.test.domain").
		CheckEqual("hostnames.0.edge_hostname_id", "ehn_123").
		CheckEqual("latest_version", "1").
		CheckEqual("staging_version", "1").
		CheckEqual("production_version", "0").
		CheckEqual("rules", `{"rules":{"name":"","options":{}}}`)

	tests := map[string]struct {
		importID   string
		config     string
		init       func(*testing.T, *mockProperty)
		stateCheck func(s []*terraform.InstanceState) error
	}{
		"Importable: property_id with ds": {
			importID: "prp_4",
			config:   "testdata/TestResProperty/Importable/importable_with_property_rules_builder.tf",
			init: func(t *testing.T, p *mockProperty) {
				// set initial data
				p.mockPropertyData = basicDataWithoutGroupAndContract
				p.ruleTree = mockRuleTreeData{
					rules: papi.Rules{
						Name: "default",
						Behaviors: []papi.RuleBehavior{
							{Name: "mPulse", Options: papi.RuleOptionsMap{"configOverride": "no new line"}},
							{Name: "mPulse", Options: papi.RuleOptionsMap{"configOverride": ""}},
							{Name: "mPulse", Options: papi.RuleOptionsMap{"configOverride": "\n\tline with new line before and after + tab\n"}},
						},
					},
				}
				mockPropertyImportKnownVersion(p)
			},
			stateCheck: defaultChecker.CheckEqual("rules", "{\"rules\":{\"behaviors\":[{\"name\":\"mPulse\",\"options\":{\"configOverride\":\"no new line\"}},{\"name\":\"mPulse\",\"options\":{\"configOverride\":\"\"}},{\"name\":\"mPulse\",\"options\":{\"configOverride\":\"\\n\\tline with new line before and after + tab\\n\"}}],\"name\":\"default\",\"options\":{}}}").Build(),
		},
		"Importable: property_id with property-bootstrap": {
			importID: "prp_4,property-bootstrap",
			config:   "testdata/TestResProperty/Importable/importable-with-bootstrap.tf",
			init: func(t *testing.T, p *mockProperty) {
				// set initial data
				p.mockPropertyData = basicDataWithoutGroupAndContract
				mockPropertyImportKnownVersion(p)
			},
			stateCheck: defaultChecker.Build(),
		},
		"Importable: property_id": {
			importID: "prp_4",
			init: func(t *testing.T, p *mockProperty) {
				// set initial data
				p.mockPropertyData = basicDataWithoutGroupAndContract
				mockPropertyImportKnownVersion(p)
			},
			stateCheck: defaultChecker.Build(),
		},
		"Importable: property_id and ver_# version": {
			importID: "prp_4,ver_1",
			init: func(t *testing.T, p *mockProperty) {
				// set initial data
				p.mockPropertyData = basicDataWithoutGroupAndContract
				mockPropertyImportUnknownVersion(p)
			},
			stateCheck: defaultChecker.Build(),
		},
		"Importable: property_id and # version": {
			importID: "prp_4,1",
			init: func(t *testing.T, p *mockProperty) {
				// set initial data
				p.mockPropertyData = basicDataWithoutGroupAndContract
				mockPropertyImportUnknownVersion(p)
			},
			stateCheck: defaultChecker.Build(),
		},
		"Importable: property_id and latest": {
			importID: "prp_4,latest",
			init: func(t *testing.T, p *mockProperty) {
				// set initial data
				p.mockPropertyData = basicDataWithoutGroupAndContract
				mockPropertyImportKnownVersion(p)
			},
			stateCheck: defaultChecker.Build(),
		},
		"Importable: property_id and network": {
			importID: "prp_4,staging",
			init: func(t *testing.T, p *mockProperty) {
				// set initial data
				p.mockPropertyData = basicDataWithoutGroupAndContract
				mockPropertyImportUnknownVersion(p)
			},
			stateCheck: defaultChecker.Build(),
		},
		"Importable: un-prefixed property_id": {
			importID: "4",
			init: func(t *testing.T, p *mockProperty) {
				// set initial data
				p.mockPropertyData = basicDataWithoutGroupAndContract
				mockPropertyImportKnownVersion(p)
			},
			stateCheck: defaultChecker.Build(),
		},
		"Importable: un-prefixed property_id and # version": {
			importID: "4,1",
			init: func(t *testing.T, p *mockProperty) {
				// set initial data
				p.mockPropertyData = basicDataWithoutGroupAndContract
				mockPropertyImportUnknownVersion(p)
			},
			stateCheck: defaultChecker.Build(),
		},
		"Importable: un-prefixed property_id and ver_# version": {
			importID: "4,ver_1",
			init: func(t *testing.T, p *mockProperty) {
				// set initial data
				p.mockPropertyData = basicDataWithoutGroupAndContract
				mockPropertyImportUnknownVersion(p)
			},
			stateCheck: defaultChecker.Build(),
		},
		"Importable: un-prefixed property_id and network": {
			importID: "4,s",
			init: func(t *testing.T, p *mockProperty) {
				// set initial data
				p.mockPropertyData = basicDataWithoutGroupAndContract
				mockPropertyImportUnknownVersion(p)
			},
			stateCheck: defaultChecker.Build(),
		},
		"Importable: property_id and contract_id and group_id": {
			importID: "prp_4,ctr_1,grp_2",
			init: func(t *testing.T, p *mockProperty) {
				// set initial data
				p.mockPropertyData = basicData
				// read
				mockResourcePropertyRead(p)
			},
			stateCheck: defaultChecker.Build(),
		},
		"Importable: property_id, contract_id, group_id and empty version": {
			importID: "prp_4,ctr_1,grp_2,",
			init: func(t *testing.T, p *mockProperty) {
				// set initial data
				p.mockPropertyData = basicData
				mockResourcePropertyRead(p)
			},
			stateCheck: defaultChecker.Build(),
		},
		"Importable: property_id, contract_id, group_id and latest": {
			importID: "prp_4,ctr_1,grp_2,latest",
			init: func(t *testing.T, p *mockProperty) {
				// set initial data
				p.mockPropertyData = basicData
				// read
				mockResourcePropertyRead(p)
			},
			stateCheck: defaultChecker.Build(),
		},
		"Importable: property_id, contract_id, group_id and ver_# version": {
			importID: "prp_4,ctr_1,grp_2,ver_1",
			init: func(t *testing.T, p *mockProperty) {
				// set initial data
				p.mockPropertyData = basicData
				mockPropertyImportKnownVersionAfterImport(p)
			},
			stateCheck: defaultChecker.Build(),
		},
		"Importable: property_id, contract_id, group_id and # version": {
			importID: "prp_4,ctr_1,grp_2,1",
			init: func(t *testing.T, p *mockProperty) {
				// set initial data
				p.mockPropertyData = basicData
				mockPropertyImportKnownVersionAfterImport(p)
			},
			stateCheck: defaultChecker.Build(),
		},
		"Importable: property_id, contract_id, group_id and network": {
			importID: "prp_4,ctr_1,grp_2,staging",
			init: func(t *testing.T, p *mockProperty) {
				// set initial data
				p.mockPropertyData = basicData
				mockPropertyImportUnknownVersion(p)
			},
			stateCheck: defaultChecker.Build(),
		},
		"Importable: un-prefixed property_id and contract_id and group_id": {
			importID: "4,1,2",
			init: func(t *testing.T, p *mockProperty) {
				// set initial data
				p.mockPropertyData = basicData
				// read
				mockResourcePropertyRead(p)
			},
			stateCheck: defaultChecker.Build(),
		},
		"Importable: un-prefixed property_id and contract_id, group_id and # version": {
			importID: "4,1,2,1",
			init: func(t *testing.T, p *mockProperty) {
				// set initial data
				p.mockPropertyData = basicData
				mockPropertyImportKnownVersionAfterImport(p)
			},
			stateCheck: defaultChecker.Build(),
		},
		"Importable: un-prefixed property_id and contract_id, group_id and ver_# version": {
			importID: "4,1,2,ver_1",
			init: func(t *testing.T, p *mockProperty) {
				// set initial data
				p.mockPropertyData = basicData
				mockPropertyImportKnownVersionAfterImport(p)
			},
			stateCheck: defaultChecker.Build(),
		},
		"Importable: un-prefixed property_id and contract_id, group_id and latest": {
			importID: "4,1,2,latest",
			init: func(t *testing.T, p *mockProperty) {
				// set initial data
				p.mockPropertyData = basicData
				// read
				mockResourcePropertyRead(p)
			},
			stateCheck: defaultChecker.Build(),
		},
		"Importable: un-prefixed property_id and contract_id, group_id and network": {
			importID: "4,1,2,staging",
			init: func(t *testing.T, p *mockProperty) {
				// set initial data
				p.mockPropertyData = basicData
				mockPropertyImportUnknownVersion(p)
			},
			stateCheck: defaultChecker.Build(),
		},
	}

	for name, test := range tests {
		t.Run(name, func(t *testing.T) {
			papiMock := &papi.Mock{}
			mp := mockProperty{
				papiMock: papiMock,
			}
			test.init(t, &mp)

			// use default config file if custom is not specified
			if test.config == "" {
				test.config = "testdata/TestResProperty/Importable/importable.tf"
			}

			useClient(papiMock, nil, func() {
				resource.UnitTest(t, resource.TestCase{
					ProtoV6ProviderFactories: testutils.NewProtoV6ProviderFactory(NewSubprovider()),
					Steps: []resource.TestStep{
						{
							ImportStateCheck:        test.stateCheck,
							ImportStateId:           test.importID,
							ImportState:             true,
							ResourceName:            "akamai_property.test",
							Config:                  testutils.LoadFixtureString(t, test.config),
							ImportStateVerifyIgnore: []string{"product", "read_version"},
						},
					},
				})
			})

			papiMock.AssertExpectations(t)
		})
	}
}

// TestPropertyErrors tests various cases where we should expect an error or validation is triggered
func TestPropertyErrors(t *testing.T) {
	// basicData holds basic, common data across test cases
	basicData := mockPropertyData{
		propertyName:  "test_property",
		contractID:    "ctr_1",
		productID:     "prd_3",
		groupID:       "grp_2",
		propertyID:    "prp_4",
		latestVersion: 1,
	}

	defaultChecker := test.NewStateChecker("akamai_property.test").
		CheckEqual("id", "prp_4").
		CheckEqual("hostnames.0.cname_to", "to.test.domain").
		CheckEqual("hostnames.0.edge_hostname_id", "ehn_123").
		CheckEqual("latest_version", "1").
		CheckEqual("staging_version", "0").
		CheckEqual("production_version", "0").
		CheckEqual("name", "test_property").
		CheckEqual("contract_id", "ctr_1").
		CheckEqual("group_id", "grp_2").
		CheckEqual("product_id", "prd_3").
		CheckEqual("rule_warnings.#", "0").
		CheckEqual("rules", `{"rules":{"name":"default","options":{}}}`)

	inactiveVersions := papi.PropertyVersionItems{
		Items: []papi.PropertyVersionGetItem{
			{
				StagingStatus:    papi.VersionStatusInactive,
				ProductionStatus: papi.VersionStatusInactive,
				PropertyVersion:  1,
			},
		},
	}

	defaultHostname := papi.HostnameResponseItems{
		Items: []papi.Hostname{
			{
				CnameType:            "EDGE_HOSTNAME",
				CnameFrom:            "from.test.domain",
				CnameTo:              "to.test.domain",
				CertProvisioningType: "DEFAULT",
				EdgeHostnameID:       "ehn_123",
			},
		},
	}

	defaultRuleTree := mockRuleTreeData{
		rules: papi.Rules{
			Name: "default",
		},
	}

	tests := map[string]struct {
		init  func(*testing.T, *mockProperty)
		steps []resource.TestStep
	}{
		"error when the given group is not found": {
			init: func(t *testing.T, p *mockProperty) {
				// set initial data
				p.mockPropertyData = basicData
				err := &papi.Error{
					StatusCode: 404,
					Title:      "Not Found",
					Detail:     "The system was unable to locate the requested resource",
					Type:       "https://problems.luna.akamaiapis.net/papi/v0/http/not-found",
					Instance:   "https://akaa-hqgqowhpmkw32kmt-t3owzo37wb5dkern.luna-dev.akamaiapis.net/papi/v1/properties?contractId=ctr_0\\u0026groupId=grp_0#c3fe5f9b0c4a14d1",
				}
				p.mockCreateProperty(err)
				p.mockGetGroups()
			},
			steps: []resource.TestStep{
				{
					Config:      testutils.LoadFixtureString(t, "testdata/TestResProperty/Creation/property.tf"),
					ExpectError: regexp.MustCompile("group not found: grp_2"),
				},
			},
		},
		"error when creating property with non-unique name": {
			init: func(t *testing.T, p *mockProperty) {
				// set initial data
				p.mockPropertyData = basicData
				err := fmt.Errorf("given property name is not unique")
				p.mockCreateProperty(err)
			},
			steps: []resource.TestStep{
				{
					Config:      testutils.LoadFixtureString(t, "testdata/TestResProperty/error_when_creating_property_with_non-unique_name.tf"),
					ExpectError: regexp.MustCompile(`property name is not unique`),
				},
			},
		},
		"error when deleting active property": {
			init: func(t *testing.T, p *mockProperty) {
				// set initial data
				p.mockPropertyData = basicData
				p.versions = inactiveVersions
				p.hostnames = defaultHostname
				p.ruleTree = defaultRuleTree
				// create
				mockResourcePropertyCreateWithVersionHostnames(p)
				// read x2
				mockResourcePropertyRead(p, 2)
				// refresh before update
				mockResourcePropertyRead(p)
				// First call to remove is not successful because property is active
				err := fmt.Errorf(`cannot remove active property "prp_4"`)
				p.mockRemoveProperty(err)
				// Second call will be successful (TF test case requires last state to be empty or it's a failed test)
				p.mockRemoveProperty()
			},
			steps: []resource.TestStep{
				{
					Config: testutils.LoadFixtureString(t, "testdata/TestResProperty/error_when_deleting_active_property/step0.tf"),
					Check:  defaultChecker.Build(),
				},
				{
					Config:      testutils.LoadFixtureString(t, "testdata/TestResProperty/error_when_deleting_active_property/step1.tf"),
					ExpectError: regexp.MustCompile(`cannot remove active property`),
				},
			},
		},
		"error validations when updating property with rules tree": {
			init: func(t *testing.T, p *mockProperty) {
				// set initial data
				p.mockPropertyData = basicData
				p.ruleTree = mockRuleTreeData{
					rules: papi.Rules{
						Name: "update rule tree",
					},
				}
				// create
				p.mockCreateProperty()
				err := &papi.Error{
					StatusCode:   400,
					Type:         "/papi/v1/errors/validation.required_behavior",
					Title:        "Missing required behavior in default rule",
					Detail:       "In order for this property to work correctly behavior Content Provider Code needs to be present in the default section",
					Instance:     "/papi/v1/properties/prp_173136/versions/3/rules#err_100",
					BehaviorName: "cpCode",
				}
				// expect an error while updating rule tree
				p.mockUpdateRuleTree(err)
				// contract and group are not set in the state, so the property deletion is performed without those attributes
				p.contractID = ""
				p.groupID = ""
				// delete
				p.mockRemoveProperty()
			},
			steps: []resource.TestStep{
				{
					Config:      testutils.LoadFixtureString(t, "testdata/TestResProperty/property_update_with_validation_error_for_rules.tf"),
					ExpectError: regexp.MustCompile(`validation.required_behavior`),
				},
			},
		},
		"validation warning when creating property with rules tree": {
			init: func(t *testing.T, p *mockProperty) {
				// set initial data
				p.mockPropertyData = basicData
				p.versions = inactiveVersions
				p.ruleTree = mockRuleTreeData{
					ruleFormat: "",
					ruleWarnings: []papi.RuleWarnings{
						{
							Type:          "https://problems.luna.akamaiapis.net/papi/v0/validation/validation_message.ip_address_origin",
							ErrorLocation: "#/rules/behaviors/1",
							Detail:        "Using an IP address for the `Origin Server` is not recommended. IP addresses may be changed or reassigned without notice which can severely impact your property or cause a DoS. Please use a properly formatted hostname instead.",
						},
					},
					rules: papi.Rules{
						Behaviors: []papi.RuleBehavior{
							{
								Name: "origin",
								Options: papi.RuleOptionsMap{
									"hostname":  "1.2.3.4",
									"httpPort":  float64(80),
									"httpsPort": float64(443),
								},
							},
						},
					},
				}
				p.responseWarnings = []*papi.Error{
					{
						Type:          "https://problems.luna.akamaiapis.net/papi/v0/validation/validation_message.ip_address_origin",
						ErrorLocation: "#/rules/behaviors/1",
						Detail:        "Using an IP address for the `Origin Server` is not recommended. IP addresses may be changed or reassigned without notice which can severely impact your property or cause a DoS. Please use a properly formatted hostname instead.",
					},
				}
				// create
				p.mockCreateProperty()
				p.mockUpdateRuleTree()
				// read x2
				mockResourcePropertyRead(p, 2)
				// delete
				p.mockRemoveProperty()
			},
			steps: []resource.TestStep{
				{
					Config: testutils.LoadFixtureString(t, "testdata/TestResProperty/property_with_validation_warning_for_rules.tf"),
					Check: defaultChecker.
						CheckEqual("rule_warnings.#", "1").
						CheckEqual("rules", `{"rules":{"behaviors":[{"name":"origin","options":{"hostname":"1.2.3.4","httpPort":80,"httpsPort":443}}],"name":"","options":{}}}`).
						CheckEqual("rule_warnings.0.detail", "Using an IP address for the `Origin Server` is not recommended. IP addresses may be changed or reassigned without notice which can severely impact your property or cause a DoS. Please use a properly formatted hostname instead.").
						CheckMissing("hostnames.0.cname_to").
						CheckMissing("hostnames.0.edge_hostname_id").
						Build(),
				},
			},
		},
		"validation - when updating a property hostnames to empty it should return error": {
			init: func(t *testing.T, p *mockProperty) {
				// set initial data
				p.mockPropertyData = basicData
				p.hostnames = papi.HostnameResponseItems{
					Items: []papi.Hostname{
						{
							CnameType:            "EDGE_HOSTNAME",
							CnameFrom:            "terraform.provider.myu877.test.net",
							CnameTo:              "terraform.provider.myu877.test.net.edgesuite.net",
							CertProvisioningType: "DEFAULT",
							EdgeHostnameID:       "ehn_123",
						},
					},
				}
				p.ruleTree = mockRuleTreeData{
					rules:      papi.Rules{},
					ruleFormat: "",
				}
				p.versions = inactiveVersions
				// create
				mockResourcePropertyCreateWithVersionHostnames(p)
				// read x2
				mockResourcePropertyRead(p, 2)
				// refresh - read
				mockResourcePropertyRead(p)
				// delete
				p.mockRemoveProperty()
			},
			steps: []resource.TestStep{
				{
					Config: testutils.LoadFixtureString(t, "testdata/TestResProperty/CreationUpdateNoHostnames/creation/property_create.tf"),
					Check: defaultChecker.
						CheckEqual("rules", `{"rules":{"name":"","options":{}}}`).
						CheckEqual("hostnames.0.cname_to", "terraform.provider.myu877.test.net.edgesuite.net").
						CheckEqual("hostnames.#", "1").
						Build(),
				},
				{
					Config:      testutils.LoadFixtureString(t, "testdata/TestResProperty/CreationUpdateNoHostnames/update/property_update.tf"),
					ExpectError: regexp.MustCompile("hostnames exist on server and cannot be updated to empty for property with id 'prp_4'. Provide at least one hostname to update existing list of hostnames associated to this property"),
				},
			},
		},
		"validation - when updating a property hostnames with cert_provisioning_type = 'DEFAULT' with secure-by-default enabled but remaining default certs == 0 it should return error": {
			init: func(t *testing.T, p *mockProperty) {
				// set initial data
				p.mockPropertyData = basicData
				p.hostnames = papi.HostnameResponseItems{
					Items: []papi.Hostname{
						{
							CnameType:            "EDGE_HOSTNAME",
							CnameFrom:            "terraform.provider.myu877.test.net",
							CnameTo:              "terraform.provider.myu877.test.net.edgesuite.net",
							CertProvisioningType: "DEFAULT",
						},
					},
				}
				// create
				p.mockCreateProperty()
				err := &papi.Error{
					StatusCode: http.StatusTooManyRequests,
					Remaining:  ptr.To(0),
					LimitKey:   "DEFAULT_CERTS_PER_CONTRACT",
				}
				p.mockUpdatePropertyVersionHostnames(err)
				// delete
				p.mockRemoveProperty()
			},
			steps: []resource.TestStep{
				{
					Config:      testutils.LoadFixtureString(t, "testdata/TestResProperty/CreationUpdateNoHostnames/creation/property_create.tf"),
					ExpectError: regexp.MustCompile("updating hostnames: not possible to use cert_provisioning_type = 'DEFAULT' as the limit for DEFAULT certificates has been reached"),
				},
			},
		},
		"validation - when updating a property hostnames with cert_provisioning_type = 'DEFAULT' not having enabled secure-by-default it should return error": {
			init: func(t *testing.T, p *mockProperty) {
				// set initial data
				p.mockPropertyData = basicData
				p.hostnames = papi.HostnameResponseItems{
					Items: []papi.Hostname{
						{
							CnameType:            "EDGE_HOSTNAME",
							CnameFrom:            "terraform.provider.myu877.test.net",
							CnameTo:              "terraform.provider.myu877.test.net.edgesuite.net",
							CertProvisioningType: "DEFAULT",
						},
					},
				}
				p.versions = inactiveVersions
				// create
				p.mockCreateProperty()
				err := &papi.Error{
					StatusCode: http.StatusForbidden,
					Type:       "https://problems.luna.akamaiapis.net/papi/v0/property-version-hostname/default-cert-provisioning-unavailable",
				}
				p.mockUpdatePropertyVersionHostnames(err)
				// delete
				p.mockRemoveProperty()
			},
			steps: []resource.TestStep{
				{
					Config:      testutils.LoadFixtureString(t, "testdata/TestResProperty/CreationUpdateNoHostnames/creation/property_create.tf"),
					ExpectError: regexp.MustCompile("updating hostnames: not possible to use cert_provisioning_type = 'DEFAULT' as secure-by-default is not enabled in this account"),
				},
			},
		},
		"400 from UpdatePropertyVersionHostnames - incorrect/invalid edge hostname": {
			init: func(t *testing.T, p *mockProperty) {
				// set initial data
				p.mockPropertyData = mockPropertyData{
					propertyName:  "dxe-2406-issue-example",
					groupID:       "grp_2",
					contractID:    "ctr_1",
					productID:     "prd_3",
					propertyID:    "prp_4",
					latestVersion: 1,
					ruleTree: mockRuleTreeData{
						ruleFormat: "",
						rules: papi.Rules{
							Name: "default",
							Children: []papi.Rules{
								{
									Name: "Static Content",
									Behaviors: []papi.RuleBehavior{
										{
											Name:    "prefetch",
											Options: papi.RuleOptionsMap{"enabled": false},
										},
									},
								},
							},
							Behaviors: []papi.RuleBehavior{
								{
									Name: "cpCode",
									Options: papi.RuleOptionsMap{
										"value": map[string]interface{}{
											"id":          float64(12345),
											"description": "WAA Example.com",
											"products": []interface{}{
												"Web_App_Accel",
											},
											"name": "WAA Example.com",
										},
									},
								},
							},
							Options: papi.RuleOptions{IsSecure: true},
						},
					},
					versions: papi.PropertyVersionItems{
						Items: []papi.PropertyVersionGetItem{
							{
								StagingStatus:    papi.VersionStatusInactive,
								ProductionStatus: papi.VersionStatusInactive,
								PropertyVersion:  1,
							},
						},
					},
					hostnames: papi.HostnameResponseItems{
						Items: []papi.Hostname{
							{
								CnameType:            "EDGE_HOSTNAME",
								CnameFrom:            "dxe-2406-issue-example-second.com",
								CnameTo:              "dxe-2406-issue-example-second.com.example.net",
								CertProvisioningType: "CPS_MANAGED",
								EdgeHostnameID:       "ehn_123",
							},
							{
								CnameType:            "EDGE_HOSTNAME",
								CnameFrom:            "dxe-2406-issue.com",
								CnameTo:              "dxe-2406-issue.com.example.net",
								CertProvisioningType: "CPS_MANAGED",
							},
						},
					},
					activationForCreate: papi.Activation{
						ActivationID:    "act_123",
						ActivationType:  papi.ActivationTypeActivate,
						Network:         papi.ActivationNetworkStaging,
						Status:          papi.ActivationStatusActive,
						NotifyEmails:    []string{"dummy-user@akamai.com"},
						PropertyVersion: 1,
						SubmitDate:      "2020-10-28T15:04:05Z",
					},
					groups: papi.GroupItems{
						Items: []*papi.Group{},
					},
				}
				// akamai_property - create
				mockResourcePropertyFullCreate(p)
				// akamai_property - read x2
				mockResourcePropertyRead(p, 2)
				// akamai_property_activation - create activation
				p.mockGetRuleTreeActivation().Once() // GetRuleTree request in activation resources is different from GetRuleTree in property resource
				p.mockGetActivations()               // no activation
				p.mockCreateActivation()
				p.mockGetActivation()

				activatedVersion := papi.PropertyVersionItems{
					Items: []papi.PropertyVersionGetItem{
						{
							ProductionStatus: papi.VersionStatusActive,
							StagingStatus:    papi.VersionStatusActive,
						},
					},
				}
				p.versions = activatedVersion

				// akamai_property_activation - read x2
				p.mockGetActivations().Twice()

				// akamai_property - read before update
				mockResourcePropertyRead(p)

				// second step
				// property update returns an error on the invalid edgehostname
				p.mockGetPropertyVersion()
				p.createFromVersion = 1
				p.newVersionID = 2
				p.mockCreatePropertyVersion()
				// after creating new version, update latest version of the property to reflect that change
				p.latestVersion = p.newVersionID

				// prepare updated hostnames
				updatedHostnames := papi.HostnameResponseItems{
					Items: []papi.Hostname{
						{
							CnameType:            "EDGE_HOSTNAME",
							CnameFrom:            "dxe-2406-issue-example-second.com",
							CnameTo:              "dxe-2406-issue-example-second.com.example.net",
							CertProvisioningType: "CPS_MANAGED",
							EdgeHostnameID:       "ehn_123",
						},
						{
							CnameType:            "EDGE_HOSTNAME",
							CnameFrom:            "dxe-2406-issue.com",
							CnameTo:              "dxe-2406-issue.com.example.net",
							CertProvisioningType: "CPS_MANAGED",
						},
						{
							CnameType:            "EDGE_HOSTNAME",
							CnameFrom:            "does-not-exist.com",
							CnameTo:              "does-not-exist.com.example.net",
							CertProvisioningType: "CPS_MANAGED",
						},
					},
				}
				p.hostnames = updatedHostnames

				// return an error while updating property version hostnames
				err := fmt.Errorf("%w: request failed: %s", papi.ErrUpdatePropertyVersionHostnames, errors.New("{\n    \"type\": \"https://problems.luna.akamaiapis.net/papi/v0/property-version-hostname/bad-cnameto\",\n    \"title\": \"Bad `cnameTo`\",\n    \"detail\": \"The System could not find cnameTo value `does-not-exist.com.example.net`.\",\n    \"instance\": \"host/papi/v1/properties/prp_0/versions/2/hostnames?contractId=ctr_0&groupId=grp_0&includeCertStatus=false&validateHostnames=false#efba6490291100b1\",\n    \"status\": 400\n}"))
				p.mockUpdatePropertyVersionHostnames(err)

				// terraform clean up - terraform test framework attempts to run destroy plan, if an error is returned on second step

				p.mockResourceActivationDelete()
				// akamai_property - delete
				p.mockRemoveProperty()
			},
			steps: []resource.TestStep{
				{
					Config: testutils.LoadFixtureString(t, "testdata/TestResProperty/CreationUpdateIncorrectEdgeHostname/create/property.tf"),
					Check: defaultChecker.
						CheckEqual("rules", `{"rules":{"behaviors":[{"name":"cpCode","options":{"value":{"description":"WAA Example.com","id":12345,"name":"WAA Example.com","products":["Web_App_Accel"]}}}],"children":[{"behaviors":[{"name":"prefetch","options":{"enabled":false}}],"name":"Static Content","options":{}}],"name":"default","options":{"is_secure":true}}}`).
						CheckEqual("name", "dxe-2406-issue-example").
						CheckEqual("hostnames.0.cname_to", "dxe-2406-issue-example-second.com.example.net").
						CheckEqual("hostnames.#", "2").
						Build(),
				},
				{
					Config:      testutils.LoadFixtureString(t, "testdata/TestResProperty/CreationUpdateIncorrectEdgeHostname/update/property.tf"),
					ExpectError: regexp.MustCompile("Error: updating hostnames: request failed:"),
				},
			},
		},
	}

	for name, test := range tests {
		t.Run(name, func(t *testing.T) {
			papiMock := &papi.Mock{}
			mp := mockProperty{
				papiMock: papiMock,
			}
			test.init(t, &mp)

			useClient(papiMock, nil, func() {
				resource.UnitTest(t, resource.TestCase{
					ProtoV6ProviderFactories: testutils.NewProtoV6ProviderFactory(NewSubprovider()),
					Steps:                    test.steps,
				})
			})

			papiMock.AssertExpectations(t)
		})
	}
}

// TestSchemaConfiguration tests errors when invalid HCL configuration is provided
func TestSchemaConfiguration(t *testing.T) {
	assertConfigError := func(t *testing.T, flaw, rx string) func(t *testing.T) {
		fixtureName := strings.ReplaceAll(flaw, " ", "_")

		return func(t *testing.T) {
			resource.UnitTest(t, resource.TestCase{
				ProtoV6ProviderFactories: testutils.NewProtoV6ProviderFactory(NewSubprovider()),
				Steps: []resource.TestStep{{
					Config:      testutils.LoadFixtureString(t, "testdata/TestResProperty/ConfigError/%s.tf", fixtureName),
					ExpectError: regexp.MustCompile(rx),
				}},
			})
		}
	}

	t.Run("Schema Configuration Error: name not given", assertConfigError(t, "name not given", `"name" is required`))
	t.Run("Schema Configuration Error: contract_id not given", assertConfigError(t, "contract_id not given", `Missing required argument`))
	t.Run("Schema Configuration Error: group_id not given", assertConfigError(t, "group_id not given", `Missing required argument`))
	t.Run("Schema Configuration Error: product_id not given", assertConfigError(t, "product_id not given", `Missing required argument`))
	t.Run("Schema Configuration Error: invalid json rules", assertConfigError(t, "invalid json rules", `rules are not valid JSON`))
	t.Run("Schema Configuration Error: invalid name given", assertConfigError(t, "invalid name given", `a name must only contain letters, numbers, and these characters: . _ -`))
	t.Run("Schema Configuration Error: name given too long", assertConfigError(t, "name given too long", `a name must be longer than 0 characters and shorter than 86 characters`))
}

func TestPropertyResource_VersionNotesLifecycle(t *testing.T) {
	testdataDir := "testdata/TestResProperty/Lifecycle/versionNotes"
	versionNotes1, versionNotes2, versionNotes3 := "lifecycleTest", "updatedNotes", "updatedNotes2"
	rulesFile1And2, rulesFile4And5, rulesFile3 := "01_02_rules.json", "04_05_rules.json", "03_rules.json"

	rulesJSON := testutils.LoadFixtureBytes(t, path.Join(testdataDir, rulesFile1And2))
	var rules1And2 papi.RulesUpdate
	err := json.Unmarshal(rulesJSON, &rules1And2)
	require.NoError(t, err)

	checker := test.NewStateChecker("akamai_property.test").
		CheckEqual("id", "prp_123").
		CheckEqual("group_id", "grp_123").
		CheckEqual("contract_id", "ctr_123").
		CheckEqual("latest_version", "1")

	papiMock := &papi.Mock{}
	basicData := mockPropertyData{
		propertyName:  "test_property",
		groupID:       "grp_123",
		contractID:    "ctr_123",
		productID:     "prd_123",
		propertyID:    "prp_123",
		latestVersion: 1,
		assetID:       "",
		versions: papi.PropertyVersionItems{
			Items: []papi.PropertyVersionGetItem{
				{
					StagingStatus:    papi.VersionStatusInactive,
					ProductionStatus: papi.VersionStatusInactive,
					PropertyVersion:  1,
					Note:             versionNotes1,
				},
			},
		},
		ruleTree: mockRuleTreeData{
			rules:      rules1And2.Rules,
			comments:   versionNotes1,
			ruleFormat: "v2023-01-05",
		},
	}

	prp := &mockProperty{
		mockPropertyData: basicData,
		papiMock:         papiMock,
	}

	// --- step 1 ---
	// create
	prp.mockCreateProperty()
	prp.mockUpdateRuleTree()

	// read x2
	mockResourcePropertyRead(prp, 2)

	// --- step 2 --- updated only notes - no triggered update
	prp.versions.Items[0].Note = versionNotes2
	prp.ruleTree.comments = versionNotes2

	// refresh x2 - no diff
	mockResourcePropertyRead(prp, 2)

	// --- step 3 ---
	var rules3 papi.RulesUpdate
	rulesJSON = testutils.LoadFixtureBytes(t, path.Join(testdataDir, rulesFile3))
	err = json.Unmarshal(rulesJSON, &rules3)
	require.NoError(t, err)
	// update with new notes and rules
	prp.versions.Items[0].Note = versionNotes3
	prp.ruleTree.rules = rules3.Rules
	prp.ruleTree.comments = versionNotes3
	mockResourcePropertyRead(prp)

	prp.mockGetPropertyVersion()
	prp.mockUpdateRuleTree()

	// read x2
	mockResourcePropertyRead(prp, 2)

	// --- step 4 ---
	// update with new notes and rules
	var rules4And5 papi.RulesUpdate
	rulesJSON = testutils.LoadFixtureBytes(t, path.Join(testdataDir, rulesFile4And5))
	err = json.Unmarshal(rulesJSON, &rules4And5)
	require.NoError(t, err)
	prp.ruleTree.comments = rules4And5.Comments
	prp.ruleTree.rules = rules4And5.Rules
	prp.versions.Items[0].Note = rules4And5.Comments
	mockResourcePropertyRead(prp)
	prp.mockGetPropertyVersion()
	prp.mockUpdateRuleTree()

	// read x2
	mockResourcePropertyRead(prp, 2)

	// --- step 5 --- same config, no diff
	// read x2
	mockResourcePropertyRead(prp, 2)

	// delete
	prp.mockRemoveProperty()

	useClient(papiMock, nil, func() {
		resource.UnitTest(t, resource.TestCase{
			ProtoV6ProviderFactories: testutils.NewProtoV6ProviderFactory(NewSubprovider()),
			Steps: []resource.TestStep{
				{
					Config: testutils.LoadFixtureString(t, path.Join(testdataDir, "01_with_notes_and_comments.tf")),
					Check: checker.
						CheckEqual("version_notes", "lifecycleTest").
						CheckEqual("rules", testutils.LoadFixtureString(t, path.Join(testdataDir, "01_expected_rules.json"))).
						Build(),
				},
				{
					Config:   testutils.LoadFixtureString(t, path.Join(testdataDir, "02_update_notes_no_diff.tf")),
					PlanOnly: true,
				},
				{
					Config: testutils.LoadFixtureString(t, path.Join(testdataDir, "03_update_notes_and_rules.tf")),
					Check: checker.
						CheckEqual("version_notes", "updatedNotes2").
						CheckEqual("rules", testutils.LoadFixtureString(t, path.Join(testdataDir, "03_expected_rules.json"))).
						Build(),
				},
				{
					Config: testutils.LoadFixtureString(t, path.Join(testdataDir, "04_05_remove_notes_update_comments.tf")),
					Check: checker.
						CheckEqual("version_notes", "Rules_04").
						CheckEqual("rules", testutils.LoadFixtureString(t, path.Join(testdataDir, "04_expected_rules.json"))).
						Build(),
				},
				{
					Config:   testutils.LoadFixtureString(t, path.Join(testdataDir, "04_05_remove_notes_update_comments.tf")),
					PlanOnly: true,
				},
			},
		})
	})
}

func TestValidatePropertyName(t *testing.T) {
	invalidNameCharacters := diag.Errorf("a name must only contain letters, numbers, and these characters: . _ -")
	invalidNameLength := diag.Errorf("a name must be longer than 0 characters and shorter than 86 characters")

	tests := map[string]struct {
		propertyName   string
		expectedReturn diag.Diagnostics
	}{
		"name contains only valid characters": {
			propertyName:   "Test_Name.With_Valid-Chars.123",
			expectedReturn: nil,
		},
		"name contains only numbers": {
			propertyName:   "123",
			expectedReturn: nil,
		},
		"name contains only letters": {
			propertyName:   "TestName",
			expectedReturn: nil,
		},
		"name contains invalid char !": {
			propertyName:   "Invalid_Char_!",
			expectedReturn: invalidNameCharacters,
		},
		"name contains invalid char @": {
			propertyName:   "@_Invalid_Char",
			expectedReturn: invalidNameCharacters,
		},
		"name contains invalid spaces": {
			propertyName:   "test name",
			expectedReturn: invalidNameCharacters,
		},
		"name too long (86 chars)": {
			propertyName:   strings.Repeat("a", 86),
			expectedReturn: invalidNameLength,
		},
		"name of max length (85 chars)": {
			propertyName:   strings.Repeat("a", 85),
			expectedReturn: nil,
		},
		"name of min length (1 char)": {
			propertyName:   strings.Repeat("a", 1),
			expectedReturn: nil,
		},
		"name empty": {
			propertyName:   "",
			expectedReturn: invalidNameLength,
		},
	}

	for name, test := range tests {
		t.Run(name, func(t *testing.T) {
			ret := validateNameWithBound(1)(test.propertyName, cty.Path{})

			assert.Equal(t, test.expectedReturn, ret)

		})
	}
}
