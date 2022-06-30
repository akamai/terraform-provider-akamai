package property

import (
	"fmt"
	"net/http"
	"regexp"
	"strings"
	"testing"

	"github.com/akamai/AkamaiOPEN-edgegrid-golang/v2/pkg/papi"
	"github.com/hashicorp/terraform-plugin-sdk/v2/diag"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/resource"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/mock"
)

func TestResProperty(t *testing.T) {
	// These more or less track the state of a Property in PAPI for the lifecycle tests
	type TestState struct {
		Client       *mockpapi
		Property     papi.Property
		Hostnames    []papi.Hostname
		VersionItems papi.PropertyVersionItems
		Rules        papi.RulesUpdate
		RuleFormat   string
	}

	// BehaviorFuncs can be composed to define common patterns of mock PAPI behavior (for Lifecycle tests)
	type BehaviorFunc = func(*TestState)

	// Combines many BehaviorFuncs into one
	ComposeBehaviors := func(behaviors ...BehaviorFunc) BehaviorFunc {
		return func(State *TestState) {
			for _, behave := range behaviors {
				behave(State)
			}
		}
	}

	UpdateRuleTree := func(propertyID, contractID, groupID string, version int, rulesUpdate *papi.RulesUpdate) BehaviorFunc {
		return func(state *TestState) {
			ExpectUpdateRuleTree(
				state.Client, propertyID, groupID, contractID, version,
				rulesUpdate, "", []papi.RuleError{},
			).Once().Run(func(args mock.Arguments) {
				state.Rules = *rulesUpdate
			})
		}
	}

	SetHostnames := func(PropertyID string, Version int, CnameTo string) BehaviorFunc {
		return func(State *TestState) {
			NewHostnames := []papi.Hostname{{
				CnameType:            "EDGE_HOSTNAME",
				CnameFrom:            "from.test.domain",
				CnameTo:              CnameTo,
				CertProvisioningType: "DEFAULT",
			}}

			ExpectUpdatePropertyVersionHostnames(State.Client, PropertyID, "grp_0", "ctr_0", Version, NewHostnames, nil).Once().Run(func(mock.Arguments) {
				NewResponseHostnames := []papi.Hostname{{
					CnameType:            "EDGE_HOSTNAME",
					CnameFrom:            "from.test.domain",
					CnameTo:              CnameTo,
					CertProvisioningType: "DEFAULT",
					EdgeHostnameID:       "ehn_123",
					CertStatus: papi.CertStatusItem{
						ValidationCname: papi.ValidationCname{
							Hostname: "_acme-challenge.www.example.com",
							Target:   "{token}.www.example.com.akamai-domain.com",
						},
						Staging: []papi.StatusItem{{Status: "PENDING"}},
						Production: []papi.StatusItem{{
							Status: "PENDING",
						},
						},
					},
				}}
				State.Hostnames = append([]papi.Hostname{}, NewResponseHostnames...)
			})
		}
	}

	SetTwoHostnames := func(PropertyID string, Version int, CnameFrom1, CnameTo1, CnameFrom2, CnameTo2 string) BehaviorFunc {
		return func(State *TestState) {
			NewHostnames := []papi.Hostname{{
				CnameType:            "EDGE_HOSTNAME",
				CnameFrom:            CnameFrom1,
				CnameTo:              CnameTo1,
				CertProvisioningType: "DEFAULT",
			}, {
				CnameType:            "EDGE_HOSTNAME",
				CnameFrom:            CnameFrom2,
				CnameTo:              CnameTo2,
				CertProvisioningType: "DEFAULT",
			}}

			ExpectUpdatePropertyVersionHostnames(State.Client, PropertyID, "grp_0", "ctr_0", Version, NewHostnames, nil).Once().Run(func(mock.Arguments) {
				NewResponseHostnames := []papi.Hostname{{
					CnameType:            "EDGE_HOSTNAME",
					CnameFrom:            CnameFrom1,
					CnameTo:              CnameTo1,
					CertProvisioningType: "DEFAULT",
					EdgeHostnameID:       "ehn_123",
					CertStatus: papi.CertStatusItem{
						ValidationCname: papi.ValidationCname{
							Hostname: "_acme-challenge.www.example.com",
							Target:   "{token}.www.example.com.akamai-domain.com",
						},
						Staging: []papi.StatusItem{{Status: "PENDING"}},
						Production: []papi.StatusItem{{
							Status: "PENDING",
						},
						},
					},
				}, {
					CnameType:            "EDGE_HOSTNAME",
					CnameFrom:            CnameFrom2,
					CnameTo:              CnameTo2,
					CertProvisioningType: "DEFAULT",
					EdgeHostnameID:       "ehn_123",
					CertStatus: papi.CertStatusItem{
						ValidationCname: papi.ValidationCname{
							Hostname: "_acme-challenge.www.example.com",
							Target:   "{token}.www.example.com.akamai-domain.com",
						},
						Staging: []papi.StatusItem{{Status: "PENDING"}},
						Production: []papi.StatusItem{{
							Status: "PENDING",
						},
						},
					},
				}}
				State.Hostnames = append([]papi.Hostname{}, NewResponseHostnames...)
			})
		}
	}

	GetPropertyVersions := func(PropertyID, PropertyName, ContractID, GroupID string, items ...papi.PropertyVersionItems) BehaviorFunc {
		return func(State *TestState) {
			versionItems := &State.VersionItems
			if len(items) > 0 {
				versionItems = &items[0]
			}
			ExpectGetPropertyVersions(State.Client, PropertyID, PropertyName, ContractID, GroupID, &State.Property, versionItems)
		}
	}

	GetPropertyVersionResources := func(PropertyID, GroupID, ContractID string, Version int, StagStatus, ProdStatus papi.VersionStatus) BehaviorFunc {
		return func(State *TestState) {
			ExpectGetPropertyVersion(State.Client, PropertyID, GroupID, ContractID, Version, StagStatus, ProdStatus)
		}
	}

	GetVersionResources := func(PropertyID, ContractID, GroupID string, Version int) BehaviorFunc {
		return func(State *TestState) {
			ExpectGetPropertyVersionHostnames(State.Client, PropertyID, GroupID, ContractID, Version, &State.Hostnames)
			ExpectGetRuleTree(State.Client, PropertyID, GroupID, ContractID, Version, &State.Rules, &State.RuleFormat)
		}
	}

	DeleteProperty := func(PropertyID string) BehaviorFunc {
		return func(State *TestState) {
			ExpectRemoveProperty(State.Client, PropertyID, "ctr_0", "grp_0").Once().Run(func(mock.Arguments) {
				State.Property = papi.Property{}
				State.Rules = papi.RulesUpdate{}
				State.Hostnames = nil
				State.RuleFormat = ""
				State.VersionItems = papi.PropertyVersionItems{}
			})
		}
	}

	GetProperty := func(PropertyID string) BehaviorFunc {
		return func(State *TestState) {
			ExpectGetProperty(State.Client, PropertyID, "grp_0", "ctr_0", &State.Property)
		}
	}

	CreateProperty := func(PropertyName, PropertyID string, rules papi.RulesUpdate) BehaviorFunc {
		return func(State *TestState) {
			ExpectCreateProperty(State.Client, PropertyName, "grp_0", "ctr_0", "prd_0", PropertyID).Run(func(mock.Arguments) {
				State.Property = papi.Property{
					PropertyName:  PropertyName,
					PropertyID:    PropertyID,
					GroupID:       "grp_0",
					ContractID:    "ctr_0",
					ProductID:     "prd_0",
					LatestVersion: 1,
				}

				State.Rules = rules
				State.RuleFormat = "v2020-01-01"
				GetProperty(PropertyID)(State)
				GetVersionResources(PropertyID, "ctr_0", "grp_0", 1)(State)
			}).Once()
		}
	}

	PropertyLifecycle := func(PropertyName, PropertyID, GroupID string, rules papi.RulesUpdate) BehaviorFunc {
		return func(State *TestState) {
			CreateProperty(PropertyName, PropertyID, rules)(State)
			GetVersionResources(PropertyID, "ctr_0", "grp_0", 1)(State)
			DeleteProperty(PropertyID)(State)
		}
	}

	ImportProperty := func(PropertyID string) BehaviorFunc {
		return func(State *TestState) {
			// Depending on how much of the import ID is given, the initial property lookup may not have group/contract
			ExpectGetProperty(State.Client, "prp_0", "grp_0", "", &State.Property).Maybe()
			ExpectGetProperty(State.Client, "prp_0", "", "", &State.Property).Maybe()
		}
	}

	AdvanceVersion := func(PropertyID string, FromVersion, ToVersion int) BehaviorFunc {
		return func(State *TestState) {
			ExpectCreatePropertyVersion(State.Client, PropertyID, "grp_0", "ctr_0", FromVersion, ToVersion).Once().Run(func(mock.Arguments) {
				State.Property.LatestVersion = ToVersion
			}).Run(func(args mock.Arguments) {
				State.Property.LatestVersion = ToVersion
				State.VersionItems.Items = append(State.VersionItems.Items,
					papi.PropertyVersionGetItem{
						ProductionStatus: papi.VersionStatusInactive,
						PropertyVersion:  ToVersion,
						StagingStatus:    papi.VersionStatusInactive,
					})
			})
			GetVersionResources(PropertyID, "ctr_0", "grp_0", ToVersion)(State)
		}
	}

	// TestCheckFunc to verify all standard attributes (for Lifecycle tests)
	CheckAttrs := func(PropertyID, CnameTo, LatestVersion, StagingVersion, ProductionVersion, EdgeHostnameId, rules string) resource.TestCheckFunc {
		return resource.ComposeAggregateTestCheckFunc(
			resource.TestCheckResourceAttr("akamai_property.test", "id", PropertyID),
			resource.TestCheckResourceAttr("akamai_property.test", "hostnames.0.cname_to", CnameTo),
			resource.TestCheckResourceAttr("akamai_property.test", "hostnames.0.edge_hostname_id", EdgeHostnameId),
			resource.TestCheckResourceAttr("akamai_property.test", "latest_version", LatestVersion),
			resource.TestCheckResourceAttr("akamai_property.test", "staging_version", StagingVersion),
			resource.TestCheckResourceAttr("akamai_property.test", "production_version", ProductionVersion),
			resource.TestCheckResourceAttr("akamai_property.test", "name", "test_property"),
			resource.TestCheckResourceAttr("akamai_property.test", "contract_id", "ctr_0"),
			resource.TestCheckResourceAttr("akamai_property.test", "contract", "ctr_0"),
			resource.TestCheckResourceAttr("akamai_property.test", "group_id", "grp_0"),
			resource.TestCheckResourceAttr("akamai_property.test", "group", "grp_0"),
			resource.TestCheckResourceAttr("akamai_property.test", "product", "prd_0"),
			resource.TestCheckResourceAttr("akamai_property.test", "product_id", "prd_0"),
			resource.TestCheckResourceAttr("akamai_property.test", "rule_warnings.#", "0"),
			resource.TestCheckResourceAttr("akamai_property.test", "rules", rules),
		)
	}

	type StepsFunc = func(State *TestState, FixturePath string) []resource.TestStep

	// Defines standard variations of client behaviors for a Lifecycle test
	type LifecycleTestCase struct {
		Name        string
		ClientSetup BehaviorFunc
		Steps       StepsFunc
	}

	// Standard test behavior for cases where the property's latest version is deactivated in staging network
	LatestVersionDeactivatedInStaging := LifecycleTestCase{
		Name: "Latest version is active in staging",
		ClientSetup: ComposeBehaviors(
			PropertyLifecycle("test_property", "prp_0", "grp_0",
				papi.RulesUpdate{Rules: papi.Rules{Name: "default"}}),
			GetPropertyVersionResources("prp_0", "grp_0", "ctr_0", 1, papi.VersionStatusDeactivated, papi.VersionStatusInactive),
			SetHostnames("prp_0", 1, "to.test.domain"),
			AdvanceVersion("prp_0", 1, 2),
			GetPropertyVersionResources("prp_0", "grp_0", "ctr_0", 2, papi.VersionStatusDeactivated, papi.VersionStatusInactive),
			SetHostnames("prp_0", 2, "to2.test.domain"),
		),
		Steps: func(State *TestState, FixturePath string) []resource.TestStep {
			return []resource.TestStep{
				{
					PreConfig: func() {
						State.VersionItems = papi.PropertyVersionItems{
							Items: []papi.PropertyVersionGetItem{{
								ProductionStatus: papi.VersionStatusInactive,
								PropertyVersion:  1,
								StagingStatus:    papi.VersionStatusDeactivated,
							}},
						}
					},
					Config: loadFixtureString("%s/step0.tf", FixturePath),
					Check: CheckAttrs("prp_0", "to.test.domain", "1", "0", "0", "ehn_123",
						"{\"rules\":{\"name\":\"default\",\"options\":{}}}"),
				},
				{
					PreConfig: func() {
						StagingVersion := 1
						State.Property.StagingVersion = &StagingVersion
					},
					Config: loadFixtureString("%s/step1.tf", FixturePath),
					Check: CheckAttrs("prp_0", "to2.test.domain", "2", "1", "0", "ehn_123",
						"{\"rules\":{\"name\":\"default\",\"options\":{}}}"),
				},
			}
		},
	}

	// Standard test behavior for cases where the property's latest version is deactivated in production network
	LatestVersionDeactivatedInProd := LifecycleTestCase{
		Name: "Latest version is active in production",
		ClientSetup: ComposeBehaviors(
			PropertyLifecycle("test_property", "prp_0", "grp_0",
				papi.RulesUpdate{Rules: papi.Rules{Name: "default"}}),
			GetPropertyVersionResources("prp_0", "grp_0", "ctr_0", 1, papi.VersionStatusInactive, papi.VersionStatusDeactivated),
			SetHostnames("prp_0", 1, "to.test.domain"),
			AdvanceVersion("prp_0", 1, 2),
			GetPropertyVersionResources("prp_0", "grp_0", "ctr_0", 2, papi.VersionStatusInactive, papi.VersionStatusDeactivated),
			SetHostnames("prp_0", 2, "to2.test.domain"),
		),
		Steps: func(State *TestState, FixturePath string) []resource.TestStep {
			return []resource.TestStep{
				{
					PreConfig: func() {
						State.VersionItems = papi.PropertyVersionItems{
							Items: []papi.PropertyVersionGetItem{{
								ProductionStatus: papi.VersionStatusInactive,
								PropertyVersion:  1,
								StagingStatus:    papi.VersionStatusActive,
							}},
						}
					},
					Config: loadFixtureString("%s/step0.tf", FixturePath),
					Check: CheckAttrs("prp_0", "to.test.domain", "1", "0", "0", "ehn_123",
						"{\"rules\":{\"name\":\"default\",\"options\":{}}}"),
				},
				{
					PreConfig: func() {
						ProductionVersion := 1
						State.Property.ProductionVersion = &ProductionVersion
					},
					Config: loadFixtureString("%s/step1.tf", FixturePath),
					Check: CheckAttrs("prp_0", "to2.test.domain", "2", "0", "1", "ehn_123",
						"{\"rules\":{\"name\":\"default\",\"options\":{}}}"),
				},
			}
		},
	}

	// Standard test behavior for cases where the property's latest version is active in staging network
	LatestVersionActiveInStaging := LifecycleTestCase{
		Name: "Latest version is active in staging",
		ClientSetup: ComposeBehaviors(
			PropertyLifecycle("test_property", "prp_0", "grp_0",
				papi.RulesUpdate{Rules: papi.Rules{Name: "default"}}),
			GetPropertyVersionResources("prp_0", "grp_0", "ctr_0", 1, papi.VersionStatusActive, papi.VersionStatusInactive),
			SetHostnames("prp_0", 1, "to.test.domain"),
			AdvanceVersion("prp_0", 1, 2),
			GetPropertyVersionResources("prp_0", "grp_0", "ctr_0", 2, papi.VersionStatusInactive, papi.VersionStatusActive),
			SetHostnames("prp_0", 2, "to2.test.domain"),
		),
		Steps: func(State *TestState, FixturePath string) []resource.TestStep {
			return []resource.TestStep{
				{
					PreConfig: func() {
						State.VersionItems = papi.PropertyVersionItems{
							Items: []papi.PropertyVersionGetItem{{
								ProductionStatus: papi.VersionStatusInactive,
								PropertyVersion:  1,
								StagingStatus:    papi.VersionStatusActive,
							}},
						}
					},
					Config: loadFixtureString("%s/step0.tf", FixturePath),
					Check: CheckAttrs("prp_0", "to.test.domain", "1", "0", "0", "ehn_123",
						"{\"rules\":{\"name\":\"default\",\"options\":{}}}"),
				},
				{
					PreConfig: func() {
						StagingVersion := 1
						State.Property.StagingVersion = &StagingVersion
					},
					Config: loadFixtureString("%s/step1.tf", FixturePath),
					Check: CheckAttrs("prp_0", "to2.test.domain", "2", "1", "0", "ehn_123",
						"{\"rules\":{\"name\":\"default\",\"options\":{}}}"),
				},
			}
		},
	}

	// Standard test behavior for cases where the property's latest version is active in production network
	LatestVersionActiveInProd := LifecycleTestCase{
		Name: "Latest version is active in production",
		ClientSetup: ComposeBehaviors(
			PropertyLifecycle("test_property", "prp_0", "grp_0",
				papi.RulesUpdate{Rules: papi.Rules{Name: "default"}}),
			GetPropertyVersionResources("prp_0", "grp_0", "ctr_0", 1, papi.VersionStatusInactive, papi.VersionStatusActive),
			SetHostnames("prp_0", 1, "to.test.domain"),
			AdvanceVersion("prp_0", 1, 2),
			GetPropertyVersionResources("prp_0", "grp_0", "ctr_0", 2, papi.VersionStatusInactive, papi.VersionStatusActive),
			SetHostnames("prp_0", 2, "to2.test.domain"),
		),
		Steps: func(State *TestState, FixturePath string) []resource.TestStep {
			return []resource.TestStep{
				{
					PreConfig: func() {
						State.VersionItems = papi.PropertyVersionItems{
							Items: []papi.PropertyVersionGetItem{{
								ProductionStatus: papi.VersionStatusActive,
								PropertyVersion:  1,
								StagingStatus:    papi.VersionStatusInactive,
							}},
						}
					},
					Config: loadFixtureString("%s/step0.tf", FixturePath),
					Check: CheckAttrs("prp_0", "to.test.domain", "1", "0", "0", "ehn_123",
						"{\"rules\":{\"name\":\"default\",\"options\":{}}}"),
				},
				{
					PreConfig: func() {
						ProductionVersion := 1
						State.Property.ProductionVersion = &ProductionVersion
					},
					Config: loadFixtureString("%s/step1.tf", FixturePath),
					Check: CheckAttrs("prp_0", "to2.test.domain", "2", "0", "1", "ehn_123",
						"{\"rules\":{\"name\":\"default\",\"options\":{}}}"),
				},
			}
		},
	}

	// Standard test behavior for cases where the property's latest version is not active
	LatestVersionNotActive := LifecycleTestCase{
		Name: "Latest version not active",
		ClientSetup: ComposeBehaviors(
			PropertyLifecycle("test_property", "prp_0", "grp_0",
				papi.RulesUpdate{Rules: papi.Rules{Name: "default"}}),
			GetPropertyVersionResources("prp_0", "grp_0", "ctr_0", 1, papi.VersionStatusInactive, papi.VersionStatusInactive),
			SetHostnames("prp_0", 1, "to.test.domain"),
			SetHostnames("prp_0", 1, "to2.test.domain"),
		),
		Steps: func(State *TestState, FixturePath string) []resource.TestStep {
			return []resource.TestStep{
				{
					PreConfig: func() {
						State.VersionItems = papi.PropertyVersionItems{Items: []papi.PropertyVersionGetItem{{PropertyVersion: 1, ProductionStatus: papi.VersionStatusInactive}}}
					},
					Config: loadFixtureString("%s/step0.tf", FixturePath),
					Check: CheckAttrs("prp_0", "to.test.domain", "1", "0", "0", "ehn_123",
						"{\"rules\":{\"name\":\"default\",\"options\":{}}}"),
				},
				{
					Config: loadFixtureString("%s/step1.tf", FixturePath),
					Check: CheckAttrs("prp_0", "to2.test.domain", "1", "0", "0", "ehn_123",
						"{\"rules\":{\"name\":\"default\",\"options\":{}}}"),
				},
			}
		},
	}

	// This scenario simulates a new version being created outside of terraform and returned on read after the first step (update should be triggered)
	ChangesMadeOutsideOfTerraform := LifecycleTestCase{
		Name: "Latest version not active",
		ClientSetup: ComposeBehaviors(
			PropertyLifecycle("test_property", "prp_0", "grp_0",
				papi.RulesUpdate{Rules: papi.Rules{Name: "default"}}),
			GetPropertyVersionResources("prp_0", "grp_0", "ctr_0", 1, papi.VersionStatusInactive, papi.VersionStatusInactive),
			SetHostnames("prp_0", 1, "to.test.domain"),
			GetPropertyVersionResources("prp_0", "grp_0", "ctr_0", 2, papi.VersionStatusInactive, papi.VersionStatusInactive),
			GetVersionResources("prp_0", "ctr_0", "grp_0", 2),
			SetHostnames("prp_0", 2, "to.test.domain"),
		),
		Steps: func(State *TestState, FixturePath string) []resource.TestStep {
			return []resource.TestStep{
				{
					PreConfig: func() {
						State.VersionItems = papi.PropertyVersionItems{Items: []papi.PropertyVersionGetItem{{PropertyVersion: 1, ProductionStatus: papi.VersionStatusInactive}}}
					},
					Config: loadFixtureString("%s/step0.tf", FixturePath),
					Check: CheckAttrs("prp_0", "to.test.domain", "1", "0", "0", "ehn_123",
						"{\"rules\":{\"name\":\"default\",\"options\":{}}}"),
				},
				{
					PreConfig: func() {
						State.Property.LatestVersion = 2
						State.Hostnames[0].CnameTo = "changed.test.domain"
					},
					Config: loadFixtureString("%s/step0.tf", FixturePath),
					Check: CheckAttrs("prp_0", "to.test.domain", "2", "0", "0", "ehn_123",
						"{\"rules\":{\"name\":\"default\",\"options\":{}}}"),
				},
			}
		},
	}

	// Standard test behavior for cases where the property's latest version is active in staging network
	NoDiff := LifecycleTestCase{
		Name: "No diff found in update",
		ClientSetup: ComposeBehaviors(
			PropertyLifecycle("test_property", "prp_0", "grp_0",
				papi.RulesUpdate{Rules: papi.Rules{Children: []papi.Rules{{Name: "Default CORS Policy", CriteriaMustSatisfy: papi.RuleCriteriaMustSatisfyAll}}}}),
			GetPropertyVersionResources("prp_0", "grp_0", "ctr_0", 1, papi.VersionStatusInactive, papi.VersionStatusInactive),
			SetHostnames("prp_0", 1, "to.test.domain"),
			UpdateRuleTree("prp_0", "ctr_0", "grp_0", 1,
				&papi.RulesUpdate{Rules: papi.Rules{Children: []papi.Rules{{CriteriaMustSatisfy: papi.RuleCriteriaMustSatisfyAll, Name: "Default CORS Policy"}}}}),
		),
		Steps: func(State *TestState, FixturePath string) []resource.TestStep {
			return []resource.TestStep{
				{
					PreConfig: func() {
						State.VersionItems = papi.PropertyVersionItems{Items: []papi.PropertyVersionGetItem{{PropertyVersion: 1, ProductionStatus: papi.VersionStatusInactive}}}
					},
					Config: loadFixtureString("%s/step0.tf", FixturePath),
					Check: CheckAttrs("prp_0", "to.test.domain", "1", "0", "0", "ehn_123",
						`{"rules":{"children":[{"name":"Default CORS Policy","options":{},"criteriaMustSatisfy":"all"}],"name":"","options":{}}}`),
				},
				{
					Config: loadFixtureString("%s/step1.tf", FixturePath),
					Check: CheckAttrs("prp_0", "to.test.domain", "1", "0", "0", "ehn_123",
						`{"rules":{"children":[{"name":"Default CORS Policy","options":{},"criteriaMustSatisfy":"all"}],"name":"","options":{}}}`),
				},
			}
		},
	}

	/*
		RulesCustomDiff tests rulesCustomDiff function which is in resource_akamai_property.go file.
		There is an additional field "options":{} in expected attributes, because with UpdateRuleTree(ctx, req) function
		this field added automatically into response, even if it does not exist in rules.
	*/
	RulesCustomDiff := LifecycleTestCase{
		Name: "Diff is only in behaviours.options.ttl",
		ClientSetup: ComposeBehaviors(
			PropertyLifecycle("test_property", "prp_0", "grp_0",
				papi.RulesUpdate{Rules: papi.Rules{Behaviors: []papi.RuleBehavior{{Name: "caching",
					Options: papi.RuleOptionsMap{"behavior": "MAX_AGE", "mustRevalidate": false, "ttl": "12d"}}},
					Name: "default"}}),
			GetPropertyVersionResources("prp_0", "grp_0", "ctr_0", 1, papi.VersionStatusInactive, papi.VersionStatusInactive),
			SetHostnames("prp_0", 1, "to.test.domain"),
			UpdateRuleTree("prp_0", "ctr_0", "grp_0", 1,
				&papi.RulesUpdate{Rules: papi.Rules{Behaviors: []papi.RuleBehavior{{Name: "caching",
					Options: papi.RuleOptionsMap{"behavior": "MAX_AGE", "mustRevalidate": false, "ttl": "12d"}}},
					Name: "default"}}),
			UpdateRuleTree("prp_0", "ctr_0", "grp_0", 1,
				&papi.RulesUpdate{Rules: papi.Rules{Behaviors: []papi.RuleBehavior{{Name: "caching",
					Options: papi.RuleOptionsMap{"behavior": "MAX_AGE", "mustRevalidate": false, "ttl": "13d"}}},
					Name: "default"}}),
		),
		Steps: func(State *TestState, FixturePath string) []resource.TestStep {
			return []resource.TestStep{
				{
					PreConfig: func() {
						State.VersionItems = papi.PropertyVersionItems{Items: []papi.PropertyVersionGetItem{{PropertyVersion: 1, ProductionStatus: papi.VersionStatusInactive}}}
					},
					Config: loadFixtureString("%s/step0.tf", FixturePath),
					Check: CheckAttrs("prp_0", "to.test.domain", "1", "0", "0", "ehn_123",
						`{"rules":{"behaviors":[{"name":"caching","options":{"behavior":"MAX_AGE","mustRevalidate":false,"ttl":"12d"}}],"name":"default","options":{}}}`),
				},
				{
					Config: loadFixtureString("%s/step1.tf", FixturePath),
					Check: CheckAttrs("prp_0", "to.test.domain", "1", "0", "0", "ehn_123",
						`{"rules":{"behaviors":[{"name":"caching","options":{"behavior":"MAX_AGE","mustRevalidate":false,"ttl":"13d"}}],"name":"default","options":{}}}`),
				},
			}
		},
	}

	NoDiffForHostnames := LifecycleTestCase{
		Name: "No diff found in update",
		ClientSetup: ComposeBehaviors(
			PropertyLifecycle("test_property", "prp_0", "grp_0",
				papi.RulesUpdate{Rules: papi.Rules{Children: []papi.Rules{{Name: "Default CORS Policy", CriteriaMustSatisfy: papi.RuleCriteriaMustSatisfyAll}}}}),
			GetPropertyVersionResources("prp_0", "grp_0", "ctr_0", 1, papi.VersionStatusInactive, papi.VersionStatusInactive),
			SetTwoHostnames("prp_0", 1, "from1.test.domain", "to1.test.domain", "from2.test.domain", "to2.test.domain"),
			UpdateRuleTree("prp_0", "ctr_0", "grp_0", 1,
				&papi.RulesUpdate{Rules: papi.Rules{Children: []papi.Rules{{CriteriaMustSatisfy: papi.RuleCriteriaMustSatisfyAll, Name: "Default CORS Policy"}}}}),
		),
		Steps: func(State *TestState, FixturePath string) []resource.TestStep {
			return []resource.TestStep{
				{
					PreConfig: func() {
						State.VersionItems = papi.PropertyVersionItems{Items: []papi.PropertyVersionGetItem{{PropertyVersion: 1, ProductionStatus: papi.VersionStatusInactive}}}
					},
					Config: loadFixtureString("%s/step0.tf", FixturePath),
					Check: CheckAttrs("prp_0", "to1.test.domain", "1", "0", "0", "ehn_123",
						`{"rules":{"children":[{"name":"Default CORS Policy","options":{},"criteriaMustSatisfy":"all"}],"name":"","options":{}}}`),
				},
				{
					Config: loadFixtureString("%s/step1.tf", FixturePath),
					Check: CheckAttrs("prp_0", "to1.test.domain", "1", "0", "0", "ehn_123",
						`{"rules":{"children":[{"name":"Default CORS Policy","options":{},"criteriaMustSatisfy":"all"}],"name":"","options":{}}}`),
				},
			}
		},
	}

	// Test Schema Configuration

	// Run a test case to verify schema validations
	AssertConfigError := func(t *testing.T, flaw, rx string) func(t *testing.T) {

		fixtureName := strings.ReplaceAll(flaw, " ", "_")

		return func(t *testing.T) {
			resource.UnitTest(t, resource.TestCase{
				Providers: testAccProviders,
				Steps: []resource.TestStep{{
					Config:      loadFixtureString("testdata/TestResProperty/ConfigError/%s.tf", fixtureName),
					ExpectError: regexp.MustCompile(rx),
				}},
			})
		}
	}

	// Test Deprecated Schema Option

	// Run a test case to verify schema attribute deprecation
	AssertDeprecated := func(t *testing.T, attribute string) func(t *testing.T) {
		return func(t *testing.T) {
			if resourceProperty().Schema[attribute].Deprecated == "" {
				t.Fatalf(`%q attribute is not marked deprecated`, attribute)
			}
		}
	}

	// Test Forbidden Schema Option

	// Run a test case to confirm that the user is prompted to read the upgrade guide
	AssertForbiddenAttr := func(t *testing.T, fixtureName string) func(t *testing.T) {

		fixtureName = strings.ReplaceAll(fixtureName, " ", "_")

		return func(t *testing.T) {
			client := &mockpapi{}
			client.Test(T{t})

			useClient(client, nil, func() {
				resource.UnitTest(t, resource.TestCase{
					Providers: testAccProviders,
					Steps: []resource.TestStep{{
						Config:      loadFixtureString("testdata/TestResProperty/ForbiddenAttr/%s.tf", fixtureName),
						ExpectError: regexp.MustCompile("See the Akamai Terraform Upgrade Guide"),
					}},
				})
			})

			client.AssertExpectations(t)
		}
	}

	// Test Lifecycle

	// Run a happy-path test case that goes through a complete create-update-destroy cycle
	AssertLifecycle := func(t *testing.T, name, variant string, kase LifecycleTestCase) func(t *testing.T) {

		fixturePrefix := fmt.Sprintf("testdata/%s/Lifecycle/%s", t.Name(), variant)

		return func(t *testing.T) {
			client := &mockpapi{}
			client.Test(T{t})
			State := &TestState{Client: client}
			kase.ClientSetup(State)

			useClient(client, nil, func() {
				resource.UnitTest(t, resource.TestCase{
					Providers:  testAccProviders,
					IsUnitTest: true,
					Steps:      kase.Steps(State, fixturePrefix),
				})
			})

			client.AssertExpectations(t)
		}
	}

	// Test Import

	// Run a test case that verifies the resource can be imported by the given ID
	AssertImportable := func(t *testing.T, TestName, ImportID string) func(t *testing.T) {

		fixturePath := fmt.Sprintf("testdata/%s/Importable/importable.tf", t.Name())

		return func(t *testing.T) {

			client := &mockpapi{}
			client.Test(T{t})

			parameters := strings.Split(ImportID, ",")
			numberParameters := len(parameters)
			lastParameter := parameters[len(parameters)-1]
			setup := []BehaviorFunc{
				PropertyLifecycle("test_property", "prp_0", "grp_0",
					papi.RulesUpdate{Rules: papi.Rules{Name: "default"}}),
				GetPropertyVersionResources("prp_0", "grp_0", "ctr_0", 1, papi.VersionStatusInactive, papi.VersionStatusInactive),
				SetHostnames("prp_0", 1, "to.test.domain"),
				ImportProperty("prp_0"),
			}
			if (numberParameters == 2 || numberParameters == 4) && !isDefaultVersion(lastParameter) {
				var ContractID, GroupID string
				if numberParameters == 4 {
					ContractID = "ctr_0"
					GroupID = "grp_0"
				}
				if numberParameters == 2 {
					setup = append(setup, GetPropertyVersions("prp_0", "test_property", "ctr_0", "grp_0"))
				}
				setup = append(setup, GetPropertyVersions("prp_0", "test_property", ContractID, GroupID))
			}
			s := ComposeBehaviors(setup...)
			kase := LifecycleTestCase{
				Name:        "Importable",
				ClientSetup: s,
				Steps: func(State *TestState, _ string) []resource.TestStep {
					return []resource.TestStep{
						{
							Config: loadFixtureString(fixturePath),
							Check: CheckAttrs("prp_0", "to.test.domain", "1", "0", "0", "ehn_123",
								"{\"rules\":{\"name\":\"default\",\"options\":{}}}"),
						},
						// this step is used to refresh state with updated staging/production statuses
						{
							PreConfig: func() {
								State.VersionItems = papi.PropertyVersionItems{Items: []papi.PropertyVersionGetItem{
									{
										PropertyVersion:  1,
										StagingStatus:    papi.VersionStatusActive,
										ProductionStatus: papi.VersionStatusActive,
									},
								}}
								stagingVersion := 1
								State.Property.StagingVersion = &stagingVersion

							},
							Config: loadFixtureString(fixturePath),
							Check: CheckAttrs("prp_0", "to.test.domain", "1", "1", "0", "ehn_123",
								"{\"rules\":{\"name\":\"default\",\"options\":{}}}"),
						},
						{
							ImportState:             true,
							ImportStateVerify:       true,
							ImportStateId:           ImportID,
							ResourceName:            "akamai_property.test",
							Config:                  loadFixtureString(fixturePath),
							ImportStateVerifyIgnore: []string{"product", "read_version"},
							Check: CheckAttrs("prp_0", "to.test.domain", "1", "1", "0", "ehn_123",
								"{\"rules\":{\"name\":\"default\",\"options\":{}}}"),
						},
					}
				},
			}
			State := &TestState{Client: client}
			kase.ClientSetup(State)
			useClient(client, nil, func() {
				resource.UnitTest(t, resource.TestCase{
					Providers: testAccProviders,
					Steps:     kase.Steps(State, ""),
				})
			})

			client.AssertExpectations(t)
		}
	}

	suppressLogging(t, func() {

		// Test Schema Configuration

		t.Run("Schema Configuration Error: name not given", AssertConfigError(t, "name not given", `"name" is required`))
		t.Run("Schema Configuration Error: neither contract nor contract_id given", AssertConfigError(t, "neither contract nor contract_id given", `one of .contract,contract_id. must be specified`))
		t.Run("Schema Configuration Error: both contract and contract_id given", AssertConfigError(t, "both contract and contract_id given", `only one of .contract,contract_id. can be specified`))
		t.Run("Schema Configuration Error: neither group nor group_id given", AssertConfigError(t, "neither group nor group_id given", `one of .group,group_id. must be specified`))
		t.Run("Schema Configuration Error: both group and group_id given", AssertConfigError(t, "both group and group_id given", `only one of .group,group_id. can be specified`))
		t.Run("Schema Configuration Error: neither product nor product_id given", AssertConfigError(t, "neither product nor product_id given", `one of .product,product_id. must be specified`))
		t.Run("Schema Configuration Error: both product and product_id given", AssertConfigError(t, "both product and product_id given", `only one of .product,product_id. can be specified`))
		t.Run("Schema Configuration Error: invalid json rules", AssertConfigError(t, "invalid json rules", `rules are not valid JSON`))
		t.Run("Schema Configuration Error: invalid name given", AssertConfigError(t, "invalid name given", `a name must only contain letters, numbers, and these characters: . _ -`))
		t.Run("Schema Configuration Error: name given too long", AssertConfigError(t, "name given too long", `a name must be shorter than 86 characters`))

		// Test Deprecated Schema Option

		t.Run("Schema deprecation: contract", AssertDeprecated(t, "contract"))
		t.Run("Schema deprecation: group", AssertDeprecated(t, "group"))
		t.Run("Schema deprecation: product", AssertDeprecated(t, "product"))
		t.Run("Schema deprecation: cp_code", AssertDeprecated(t, "cp_code"))
		t.Run("Schema deprecation: contact", AssertDeprecated(t, "contact"))
		t.Run("Schema deprecation: origin", AssertDeprecated(t, "origin"))
		t.Run("Schema deprecation: is_secure", AssertDeprecated(t, "is_secure"))
		t.Run("Schema deprecation: variables", AssertDeprecated(t, "variables"))

		// Test Forbidden Schema Option

		t.Run("Schema forbidden attribute: cp_code", AssertForbiddenAttr(t, "cp_code"))
		t.Run("Schema forbidden attribute: contact", AssertForbiddenAttr(t, "contact"))
		t.Run("Schema forbidden attribute: origin", AssertForbiddenAttr(t, "origin"))
		t.Run("Schema forbidden attribute: is_secure", AssertForbiddenAttr(t, "is_secure"))
		t.Run("Schema forbidden attribute: variables", AssertForbiddenAttr(t, "variables"))

		// Test Lifecycle

		t.Run("Lifecycle: latest version is not active (normal)", AssertLifecycle(t, t.Name(), "normal", LatestVersionNotActive))
		t.Run("Lifecycle: latest version is active in staging (normal)", AssertLifecycle(t, t.Name(), "normal", LatestVersionActiveInStaging))
		t.Run("Lifecycle: latest version is active in production (normal)", AssertLifecycle(t, t.Name(), "normal", LatestVersionActiveInProd))
		t.Run("Lifecycle: latest version is deactivated in staging (normal)", AssertLifecycle(t, t.Name(), "normal", LatestVersionDeactivatedInStaging))
		t.Run("Lifecycle: latest version is deactivated in production (normal)", AssertLifecycle(t, t.Name(), "normal", LatestVersionDeactivatedInProd))
		t.Run("Lifecycle: latest version is not active (contract_id without prefix)", AssertLifecycle(t, t.Name(), "contract_id without prefix", LatestVersionNotActive))
		t.Run("Lifecycle: latest version active in staging (contract_id without prefix)", AssertLifecycle(t, t.Name(), "contract_id without prefix", LatestVersionActiveInStaging))
		t.Run("Lifecycle: latest version active in production (contract_id without prefix)", AssertLifecycle(t, t.Name(), "contract_id without prefix", LatestVersionActiveInProd))
		t.Run("Lifecycle: latest version is not active (contract without prefix)", AssertLifecycle(t, t.Name(), "contract without prefix", LatestVersionNotActive))
		t.Run("Lifecycle: latest version is active in staging (contract without prefix)", AssertLifecycle(t, t.Name(), "contract without prefix", LatestVersionActiveInStaging))
		t.Run("Lifecycle: latest version is active in production (contract without prefix)", AssertLifecycle(t, t.Name(), "contract without prefix", LatestVersionActiveInProd))
		t.Run("Lifecycle: latest version is not active (group_id without prefix)", AssertLifecycle(t, t.Name(), "group_id without prefix", LatestVersionNotActive))
		t.Run("Lifecycle: latest version is active in staging (group_id without prefix)", AssertLifecycle(t, t.Name(), "group_id without prefix", LatestVersionActiveInStaging))
		t.Run("Lifecycle: latest version is active in production (group_id without prefix)", AssertLifecycle(t, t.Name(), "group_id without prefix", LatestVersionActiveInProd))
		t.Run("Lifecycle: latest version is not active (group without prefix)", AssertLifecycle(t, t.Name(), "group without prefix", LatestVersionNotActive))
		t.Run("Lifecycle: latest version is active in staging (group without prefix)", AssertLifecycle(t, t.Name(), "group without prefix", LatestVersionActiveInStaging))
		t.Run("Lifecycle: latest version is active in production (group without prefix)", AssertLifecycle(t, t.Name(), "group without prefix", LatestVersionActiveInProd))
		t.Run("Lifecycle: latest version is not active (product_id without prefix)", AssertLifecycle(t, t.Name(), "product_id without prefix", LatestVersionNotActive))
		t.Run("Lifecycle: latest version is active in staging (product_id without prefix)", AssertLifecycle(t, t.Name(), "product_id without prefix", LatestVersionActiveInStaging))
		t.Run("Lifecycle: latest version is active in production (product_id without prefix)", AssertLifecycle(t, t.Name(), "product_id without prefix", LatestVersionActiveInProd))
		t.Run("Lifecycle: latest version is not active (product without prefix)", AssertLifecycle(t, t.Name(), "product without prefix", LatestVersionNotActive))
		t.Run("Lifecycle: latest version is active in staging (product without prefix)", AssertLifecycle(t, t.Name(), "product without prefix", LatestVersionActiveInStaging))
		t.Run("Lifecycle: latest version is active in production (product without prefix)", AssertLifecycle(t, t.Name(), "product without prefix", LatestVersionActiveInProd))
		t.Run("Lifecycle: no diff", AssertLifecycle(t, t.Name(), "no diff", NoDiff))
		t.Run("Lifecycle: no diff (product to product_id)", AssertLifecycle(t, t.Name(), "product to product_id", NoDiff))
		t.Run("Lifecycle: no diff (product_id to product)", AssertLifecycle(t, t.Name(), "product_id to product", NoDiff))
		t.Run("Lifecycle: rules custom diff", AssertLifecycle(t, t.Name(), "rules custom diff", RulesCustomDiff))
		t.Run("Lifecycle: no diff for hostnames (hostnames)", AssertLifecycle(t, t.Name(), "hostnames", NoDiffForHostnames))
		t.Run("Lifecycle: new version changed on server", AssertLifecycle(t, t.Name(), "new version changed on server", ChangesMadeOutsideOfTerraform))

		// Test Import

		t.Run("Importable: property_id", AssertImportable(t, "property_id", "prp_0"))
		t.Run("Importable: property_id and ver_# version", AssertImportable(t, "property_id and ver_# version", "prp_0,ver_1"))
		t.Run("Importable: property_id and # version", AssertImportable(t, "property_id and # version", "prp_0,1"))
		t.Run("Importable: property_id and latest", AssertImportable(t, "property_id and latest", "prp_0,latest"))
		t.Run("Importable: property_id and network", AssertImportable(t, "property_id and network", "prp_0,staging"))
		t.Run("Importable: unprefixed property_id", AssertImportable(t, "unprefixed property_id", "0"))
		t.Run("Importable: unprefixed property_id and # version", AssertImportable(t, "unprefixed property_id and # version", "0,1"))
		t.Run("Importable: unprefixed property_id and ver_# version", AssertImportable(t, "unprefixed property_id and ver_# version", "0,ver_1"))
		t.Run("Importable: unprefixed property_id and network", AssertImportable(t, "unprefixed property_id and network", "0,p"))
		t.Run("Importable: property_id and contract_id and group_id", AssertImportable(t, "property_id and contract_id and group_id", "prp_0,ctr_0,grp_0"))
		t.Run("Importable: property_id, contract_id, group_id and empty version", AssertImportable(t, "property_id, contract_id, group_id and empty version", "prp_0,ctr_0,grp_0,"))
		t.Run("Importable: property_id, contract_id, group_id and latest", AssertImportable(t, "property_id, contract_id, group_id and latest", "prp_0,ctr_0,grp_0,latest"))
		t.Run("Importable: property_id, contract_id, group_id and ver_# version", AssertImportable(t, "property_id, contract_id, group_id and ver_# version", "prp_0,ctr_0,grp_0,ver_1"))
		t.Run("Importable: property_id, contract_id, group_id and # version", AssertImportable(t, "property_id, contract_id, group_id and # version", "prp_0,ctr_0,grp_0,1"))
		t.Run("Importable: property_id, contract_id, group_id and network", AssertImportable(t, "property_id, contract_id, group_id and network", "prp_0,ctr_0,grp_0,staging"))
		t.Run("Importable: unprefixed property_id and contract_id and group_id", AssertImportable(t, "unprefixed property_id and contract_id and group_id", "0,0,0"))
		t.Run("Importable: unprefixed property_id and contract_id, group_id and # version", AssertImportable(t, "unprefixed property_id and contract_id, group_id and # version", "0,0,0,1"))
		t.Run("Importable: unprefixed property_id and contract_id, group_id and ver_# version", AssertImportable(t, "unprefixed property_id and contract_id, group_id and ver_# version", "0,0,0,ver_1"))
		t.Run("Importable: unprefixed property_id and contract_id, group_id and latest", AssertImportable(t, "unprefixed property_id and contract_id, group_id and latest", "0,0,0,latest"))
		t.Run("Importable: unprefixed property_id and contract_id, group_id and network", AssertImportable(t, "unprefixed property_id and contract_id, group_id and network", "0,0,0,production"))

		// Test Delete

		t.Run("property is destroyed and recreated when name is changed", func(t *testing.T) {
			client := &mockpapi{}
			client.Test(T{t})

			setup := ComposeBehaviors(
				PropertyLifecycle("test_property", "prp_0", "grp_0",
					papi.RulesUpdate{Rules: papi.Rules{Name: "default"}}),
				GetPropertyVersionResources("prp_0", "grp_0", "ctr_0", 1, papi.VersionStatusInactive, papi.VersionStatusInactive),
				PropertyLifecycle("renamed_property", "prp_1", "grp_0",
					papi.RulesUpdate{Rules: papi.Rules{Name: "default"}}),
				GetPropertyVersionResources("prp_1", "grp_0", "ctr_0", 1, papi.VersionStatusInactive, papi.VersionStatusInactive),
				SetHostnames("prp_0", 1, "to.test.domain"),
				SetHostnames("prp_1", 1, "to2.test.domain"),
			)
			setup(&TestState{Client: client})

			useClient(client, nil, func() {
				resource.UnitTest(t, resource.TestCase{
					Providers: testAccProviders,
					Steps: []resource.TestStep{
						{
							Config: loadFixtureString("testdata/%s-step0.tf", t.Name()),
							Check: CheckAttrs("prp_0", "to.test.domain", "1", "0", "0", "ehn_123",
								"{\"rules\":{\"name\":\"default\",\"options\":{}}}"),
						},
						{
							Config: loadFixtureString("testdata/%s-step1.tf", t.Name()),
							Check: resource.ComposeAggregateTestCheckFunc(
								resource.TestCheckResourceAttr("akamai_property.test", "id", "prp_1"),
								resource.TestCheckResourceAttr("akamai_property.test", "name", "renamed_property"),
							),
						},
					},
				})
			})

			client.AssertExpectations(t)
		})

		t.Run("error when deleting active property", func(t *testing.T) {
			client := &mockpapi{}
			client.Test(T{t})

			setup := ComposeBehaviors(
				CreateProperty("test_property", "prp_0", papi.RulesUpdate{Rules: papi.Rules{Name: "default"}}),
				GetProperty("prp_0"),
				GetVersionResources("prp_0", "ctr_0", "grp_0", 1),
				GetPropertyVersionResources("prp_0", "grp_0", "ctr_0", 1, "ctr_0", "grp_0"),
				SetHostnames("prp_0", 1, "to.test.domain"),
			)
			setup(&TestState{Client: client})

			// First call to remove is not successful
			req := papi.RemovePropertyRequest{
				PropertyID: "prp_0",
				ContractID: "ctr_0",
				GroupID:    "grp_0",
			}

			err := fmt.Errorf(`cannot remove active property "prp_0"`)
			client.On("RemoveProperty", AnyCTX, req).Return(nil, err).Once()

			// Second call will be successful (TF test case requires last state to be empty or it's a failed test)
			ExpectRemoveProperty(client, "prp_0", "ctr_0", "grp_0").Once()

			useClient(client, nil, func() {
				resource.UnitTest(t, resource.TestCase{
					Providers: testAccProviders,
					Steps: []resource.TestStep{
						{
							Config: loadFixtureString("testdata/%s/step0.tf", t.Name()),
							Check: CheckAttrs("prp_0", "to.test.domain", "1", "0", "0", "ehn_123",
								"{\"rules\":{\"name\":\"default\",\"options\":{}}}"),
						},
						{
							Config:      loadFixtureString("testdata/%s/step1.tf", t.Name()),
							ExpectError: regexp.MustCompile(`cannot remove active property`),
						},
					},
				})
			})

			client.AssertExpectations(t)
		})

		// Test validation

		t.Run("error validations when updating property with rules tree", func(t *testing.T) {
			client := &mockpapi{}
			client.Test(T{t})
			ExpectCreateProperty(
				client, "test_property", "grp_0",
				"ctr_0", "prd_0", "prp_1",
			)

			var err error = &papi.Error{
				StatusCode:   400,
				Type:         "/papi/v1/errors/validation.required_behavior",
				Title:        "Missing required behavior in default rule",
				Detail:       "In order for this property to work correctly behavior Content Provider Code needs to be present in the default section",
				Instance:     "/papi/v1/properties/prp_173136/versions/3/rules#err_100",
				BehaviorName: "cpCode",
			}
			var req = papi.UpdateRulesRequest{
				PropertyID:      "prp_1",
				ContractID:      "ctr_0",
				GroupID:         "grp_0",
				PropertyVersion: 1,
				Rules: papi.RulesUpdate{Rules: papi.Rules{
					Name: "update rule tree",
				}},
				ValidateRules: true,
			}
			client.On("UpdateRuleTree", AnyCTX, req).Return(nil, err).Once()

			ExpectRemoveProperty(client, "prp_1", "", "")
			useClient(client, nil, func() {
				resource.UnitTest(t, resource.TestCase{
					Providers: testAccProviders,
					Steps: []resource.TestStep{
						{
							Config: loadFixtureString("testdata/TestResProperty/property_update_with_validation_error_for_rules.tf"),
							Check: resource.ComposeAggregateTestCheckFunc(
								resource.TestCheckNoResourceAttr("akamai_property.test", "rules")),
							ExpectError: regexp.MustCompile(`validation.required_behavior`),
						},
					},
				})
			})

			client.AssertExpectations(t)
		})

		t.Run("validation - when updating a property hostnames to empty it should return error", func(t *testing.T) {
			client := &mockpapi{}
			client.Test(T{t})

			ExpectCreateProperty(
				client, "test_property", "grp_0",
				"ctr_0", "prd_0", "prp_0",
			)

			ExpectGetPropertyVersions(client, "prp_0", "test_property", "ctr_0", "grp_0", nil, &papi.PropertyVersionItems{Items: []papi.PropertyVersionGetItem{
				{
					PropertyVersion:  1,
					StagingStatus:    papi.VersionStatusInactive,
					ProductionStatus: papi.VersionStatusInactive,
				},
			}})

			ExpectGetPropertyVersion(client, "prp_0", "grp_0", "ctr_0", 1, papi.VersionStatusInactive, papi.VersionStatusInactive)

			ExpectUpdatePropertyVersionHostnames(
				client, "prp_0", "grp_0", "ctr_0", 1,
				[]papi.Hostname{{
					CnameType:            "EDGE_HOSTNAME",
					CnameFrom:            "terraform.provider.myu877.test.net",
					CnameTo:              "terraform.provider.myu877.test.net.edgesuite.net",
					CertProvisioningType: "DEFAULT",
				}}, nil,
			).Once()

			ExpectGetProperty(
				client, "prp_0", "grp_0", "ctr_0",
				&papi.Property{
					PropertyID: "prp_0", GroupID: "grp_0", ContractID: "ctr_0", LatestVersion: 1,
					PropertyName: "test_property",
				},
			)

			ExpectGetPropertyVersionHostnames(
				client, "prp_0", "grp_0", "ctr_0", 1,
				&[]papi.Hostname{{
					CnameFrom:            "terraform.provider.myu877.test.net",
					CnameTo:              "terraform.provider.myu877.test.net.edgesuite.net",
					CertProvisioningType: "DEFAULT",
				}},
			).Times(3)

			ruleFormat := ""
			ExpectGetRuleTree(
				client, "prp_0", "grp_0", "ctr_0", 1,
				&papi.RulesUpdate{}, &ruleFormat,
			)

			ExpectRemoveProperty(client, "prp_0", "ctr_0", "grp_0")

			ExpectUpdatePropertyVersionHostnames(
				client, "prp_0", "grp_0", "ctr_0", 1,
				[]papi.Hostname{}, nil,
			).Once()

			ExpectGetPropertyVersionHostnames(
				client, "prp_0", "grp_0", "ctr_0", 1,
				&[]papi.Hostname{},
			).Twice()

			useClient(client, nil, func() {
				resource.UnitTest(t, resource.TestCase{
					Providers: testAccProviders,
					Steps: []resource.TestStep{
						{
							Config: loadFixtureString("testdata/TestResProperty/CreationUpdateNoHostnames/creation/property_create.tf"),
							Check:  resource.TestCheckResourceAttr("akamai_property.test", "id", "prp_0"),
						},
						{
							Config: loadFixtureString("testdata/TestResProperty/CreationUpdateNoHostnames/update/property_update.tf"),
							Check: resource.ComposeAggregateTestCheckFunc(
								resource.TestCheckResourceAttr("akamai_property.test", "id", "prp_0"),
								resource.TestCheckResourceAttr("akamai_property.test", "hostnames.#", "0"),
							),
							ExpectError: regexp.MustCompile("hostnames exist on server and cannot be updated to empty for property with id 'prp_0'. Provide at least one hostname to update existing list of hostnames associated to this property"),
						},
					},
				})
			})
		})

		t.Run("validation - when updating a property hostnames with cert_provisioning_type = 'DEFAULT' with secure-by-default enabled but remaining default certs == 0 it should return error", func(t *testing.T) {
			client := &mockpapi{}
			client.Test(T{t})

			ExpectCreateProperty(
				client, "test_property", "grp_0",
				"ctr_0", "prd_0", "prp_0",
			)

			ExpectGetPropertyVersions(client, "prp_0", "test_property", "ctr_0", "grp_0", nil, &papi.PropertyVersionItems{Items: []papi.PropertyVersionGetItem{
				{
					PropertyVersion:  1,
					StagingStatus:    papi.VersionStatusInactive,
					ProductionStatus: papi.VersionStatusInactive,
				},
			}})

			ExpectGetPropertyVersion(client, "prp_0", "grp_0", "ctr_0", 1, papi.VersionStatusInactive, papi.VersionStatusInactive)

			ExpectUpdatePropertyVersionHostnames(
				client, "prp_0", "grp_0", "ctr_0", 1,
				[]papi.Hostname{{
					CnameType:            "EDGE_HOSTNAME",
					CnameFrom:            "terraform.provider.myu877.test.net",
					CnameTo:              "terraform.provider.myu877.test.net.edgesuite.net",
					CertProvisioningType: "DEFAULT",
				}}, &papi.Error{
					StatusCode: http.StatusTooManyRequests,
					Remaining:  0,
					LimitKey:   "DEFAULT_CERTS_PER_CONTRACT",
				},
			).Once()

			ExpectRemoveProperty(client, "prp_0", "ctr_0", "grp_0")

			useClient(client, nil, func() {
				resource.UnitTest(t, resource.TestCase{
					Providers: testAccProviders,
					Steps: []resource.TestStep{
						{
							Config:      loadFixtureString("testdata/TestResProperty/CreationUpdateNoHostnames/creation/property_create.tf"),
							Check:       resource.TestCheckResourceAttr("akamai_property.test", "id", "prp_0"),
							ExpectError: regexp.MustCompile("updating hostnames: not possible to use cert_provisioning_type = 'DEFAULT' as the limit for DEFAULT certificates has been reached"),
						},
					},
				})
			})
		})

		t.Run("validation - when updating a property hostnames with cert_provisioning_type = 'DEFAULT' not having enabled secure-by-default it should return error", func(t *testing.T) {
			client := &mockpapi{}
			client.Test(T{t})

			ExpectCreateProperty(
				client, "test_property", "grp_0",
				"ctr_0", "prd_0", "prp_0",
			)

			ExpectGetPropertyVersions(client, "prp_0", "test_property", "ctr_0", "grp_0", nil, &papi.PropertyVersionItems{Items: []papi.PropertyVersionGetItem{
				{
					PropertyVersion:  1,
					StagingStatus:    papi.VersionStatusInactive,
					ProductionStatus: papi.VersionStatusInactive,
				},
			}})

			ExpectGetPropertyVersion(client, "prp_0", "grp_0", "ctr_0", 1, papi.VersionStatusInactive, papi.VersionStatusInactive)

			ExpectUpdatePropertyVersionHostnames(
				client, "prp_0", "grp_0", "ctr_0", 1,
				[]papi.Hostname{{
					CnameType:            "EDGE_HOSTNAME",
					CnameFrom:            "terraform.provider.myu877.test.net",
					CnameTo:              "terraform.provider.myu877.test.net.edgesuite.net",
					CertProvisioningType: "DEFAULT",
				}}, &papi.Error{
					StatusCode: http.StatusForbidden,
					Type:       "https://problems.luna.akamaiapis.net/papi/v0/property-version-hostname/default-cert-provisioning-unavailable",
				},
			).Once()

			ExpectRemoveProperty(client, "prp_0", "ctr_0", "grp_0")

			useClient(client, nil, func() {
				resource.UnitTest(t, resource.TestCase{
					Providers: testAccProviders,
					Steps: []resource.TestStep{
						{
							Config:      loadFixtureString("testdata/TestResProperty/CreationUpdateNoHostnames/creation/property_create.tf"),
							Check:       resource.TestCheckResourceAttr("akamai_property.test", "id", "prp_0"),
							ExpectError: regexp.MustCompile("updating hostnames: not possible to use cert_provisioning_type = 'DEFAULT' as secure-by-default is not enabled in this account"),
						},
					},
				})
			})
		})

		// Other tests

		t.Run("error when the given group is not found", func(t *testing.T) {
			client := &mockpapi{}
			client.Test(T{t})

			req := papi.CreatePropertyRequest{
				ContractID: "ctr_0",
				GroupID:    "grp_0",
				Property: papi.PropertyCreate{
					ProductID:    "prd_0",
					PropertyName: "property_name",
				},
			}

			var err error = &papi.Error{
				StatusCode: 404,
				Title:      "Not Found",
				Detail:     "The system was unable to locate the requested resource",
				Type:       "https://problems.luna.akamaiapis.net/papi/v0/http/not-found",
				Instance:   "https://akaa-hqgqowhpmkw32kmt-t3owzo37wb5dkern.luna-dev.akamaiapis.net/papi/v1/properties?contractId=ctr_0\\u0026groupId=grp_0#c3fe5f9b0c4a14d1",
			}

			client.On("CreateProperty", AnyCTX, req).Return(nil, err).Once()

			// the papi GetGroups call should not return any matching group
			var Groups []*papi.Group
			ExpectGetGroups(client, &Groups).Once()

			useClient(client, nil, func() {
				resource.UnitTest(t, resource.TestCase{
					Providers: testAccProviders,
					Steps: []resource.TestStep{{
						Config:      loadFixtureString("testdata/TestResProperty/Creation/property.tf"),
						ExpectError: regexp.MustCompile("group not found: grp_0"),
					}},
				})
			})

			client.AssertExpectations(t)
		})

		t.Run("error when creating property with non-unique name", func(t *testing.T) {
			client := &mockpapi{}
			client.Test(T{t})

			req := papi.CreatePropertyRequest{
				ContractID: "ctr_0",
				GroupID:    "grp_0",
				Property: papi.PropertyCreate{
					PropertyName: "test_property",
					ProductID:    "prd_0",
				},
			}

			client.On("CreateProperty", AnyCTX, req).Return(nil, fmt.Errorf("given property name is not unique"))
			useClient(client, nil, func() {
				resource.UnitTest(t, resource.TestCase{
					Providers: testAccProviders,
					Steps: []resource.TestStep{
						{
							Config:      loadFixtureString("testdata/%s.tf", t.Name()),
							Check:       resource.TestCheckNoResourceAttr("akamai_property.test", "id"),
							ExpectError: regexp.MustCompile(`property name is not unique`),
						},
					},
				})
			})

			client.AssertExpectations(t)
		})

	})
}

func TestValidatePropertyName(t *testing.T) {
	invalidNameCharacters := diag.Errorf("a name must only contain letters, numbers, and these characters: . _ -")
	invalidNameLength := diag.Errorf("a name must be shorter than 86 characters")

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
			expectedReturn: invalidNameCharacters,
		},
	}

	for name, test := range tests {
		t.Run(name, func(t *testing.T) {
			ret := validatePropertyName(test.propertyName, nil)

			assert.Equal(t, test.expectedReturn, ret)

		})
	}
}
