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

	"github.com/akamai/AkamaiOPEN-edgegrid-golang/v8/pkg/papi"
	"github.com/akamai/terraform-provider-akamai/v6/pkg/common/ptr"
	"github.com/akamai/terraform-provider-akamai/v6/pkg/common/testutils"
	"github.com/hashicorp/go-cty/cty"
	"github.com/hashicorp/terraform-plugin-sdk/v2/diag"
	"github.com/hashicorp/terraform-plugin-testing/helper/resource"
	"github.com/hashicorp/terraform-plugin-testing/plancheck"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/mock"
	"github.com/stretchr/testify/require"
)

func TestResProperty(t *testing.T) {
	// These more or less track the state of a Property in PAPI for the lifecycle tests
	type TestState struct {
		Client       *papi.Mock
		Property     papi.Property
		Hostnames    []papi.Hostname
		VersionItems papi.PropertyVersionItems
		Rules        papi.RulesUpdate
		RuleFormat   string
	}

	// BehaviorFuncs can be composed to define common patterns of mock PAPI behavior (for Lifecycle tests)
	type BehaviorFunc = func(*TestState)

	// Combines many BehaviorFuncs into one
	composeBehaviors := func(behaviors ...BehaviorFunc) BehaviorFunc {
		return func(State *TestState) {
			for _, behave := range behaviors {
				behave(State)
			}
		}
	}

	updateRuleTree := func(propertyID, contractID, groupID string, version int, rulesUpdate *papi.RulesUpdate) BehaviorFunc {
		return func(state *TestState) {
			ExpectUpdateRuleTree(
				state.Client, propertyID, groupID, contractID, version,
				rulesUpdate, "", []papi.RuleError{},
			).Once().Run(func(args mock.Arguments) {
				state.Rules = *rulesUpdate
			})
		}
	}

	setHostnames := func(propertyID string, version int, cnameTo string) BehaviorFunc {
		return func(state *TestState) {
			newHostnames := []papi.Hostname{{
				CnameType:            "EDGE_HOSTNAME",
				CnameFrom:            "from.test.domain",
				CnameTo:              cnameTo,
				CertProvisioningType: "DEFAULT",
			}}

			ExpectUpdatePropertyVersionHostnames(state.Client, propertyID, "grp_0", "ctr_0", version, newHostnames, nil).Once().Run(func(mock.Arguments) {
				newResponseHostnames := []papi.Hostname{{
					CnameType:            "EDGE_HOSTNAME",
					CnameFrom:            "from.test.domain",
					CnameTo:              cnameTo,
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
				state.Hostnames = append([]papi.Hostname{}, newResponseHostnames...)
			})
		}
	}

	setTwoHostnames := func(propertyID string, version int, cnameFrom1, cnameTo1, cnameFrom2, cnameTo2 string) BehaviorFunc {
		return func(state *TestState) {
			newHostnames := []papi.Hostname{{
				CnameType:            "EDGE_HOSTNAME",
				CnameFrom:            cnameFrom1,
				CnameTo:              cnameTo1,
				CertProvisioningType: "DEFAULT",
			}, {
				CnameType:            "EDGE_HOSTNAME",
				CnameFrom:            cnameFrom2,
				CnameTo:              cnameTo2,
				CertProvisioningType: "DEFAULT",
			}}

			ExpectUpdatePropertyVersionHostnames(state.Client, propertyID, "grp_0", "ctr_0", version, newHostnames, nil).Once().Run(func(mock.Arguments) {
				NewResponseHostnames := []papi.Hostname{{
					CnameType:            "EDGE_HOSTNAME",
					CnameFrom:            cnameFrom1,
					CnameTo:              cnameTo1,
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
					CnameFrom:            cnameFrom2,
					CnameTo:              cnameTo2,
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
				state.Hostnames = append([]papi.Hostname{}, NewResponseHostnames...)
			})
		}
	}

	getPropertyVersions := func(propertyID, propertyName, contractID, groupID string, items ...papi.PropertyVersionItems) BehaviorFunc {
		return func(state *TestState) {
			versionItems := &state.VersionItems
			if len(items) > 0 {
				versionItems = &items[0]
			}
			ExpectGetPropertyVersions(state.Client, propertyID, propertyName, contractID, groupID, &state.Property, versionItems)
		}
	}

	getPropertyVersionResources := func(propertyID, groupID, contractID string, version int, stagStatus, prodStatus papi.VersionStatus) BehaviorFunc {
		return func(state *TestState) {
			ExpectGetPropertyVersion(state.Client, propertyID, groupID, contractID, version, stagStatus, prodStatus)
		}
	}

	GetVersionResources := func(propertyID, contractID, groupID string, version int) BehaviorFunc {
		return func(state *TestState) {
			ExpectGetPropertyVersionHostnames(state.Client, propertyID, groupID, contractID, version, &state.Hostnames)
			ExpectGetRuleTree(state.Client, propertyID, groupID, contractID, version, &state.Rules, &state.RuleFormat, nil, nil)
		}
	}

	GetVersionResourcesDrift := func(propertyID, contractID, groupID string, version int, rules papi.RulesUpdate) BehaviorFunc {
		return func(state *TestState) {
			ExpectGetPropertyVersionHostnames(state.Client, propertyID, groupID, contractID, version, &state.Hostnames)
			ExpectGetRuleTree(state.Client, propertyID, groupID, contractID, version, &rules, &state.RuleFormat, nil, nil)
		}
	}

	DeleteProperty := func(propertyID string) BehaviorFunc {
		return func(state *TestState) {
			ExpectRemoveProperty(state.Client, propertyID, "ctr_0", "grp_0").Once().Run(func(mock.Arguments) {
				state.Property = papi.Property{}
				state.Rules = papi.RulesUpdate{}
				state.Hostnames = nil
				state.RuleFormat = ""
				state.VersionItems = papi.PropertyVersionItems{}
			})
		}
	}

	getProperty := func(propertyID string) BehaviorFunc {
		return func(state *TestState) {
			ExpectGetProperty(state.Client, propertyID, "grp_0", "ctr_0", &state.Property)
		}
	}

	createProperty := func(propertyName, propertyID string, rules papi.RulesUpdate) BehaviorFunc {
		return func(state *TestState) {
			ExpectCreateProperty(state.Client, propertyName, "grp_0", "ctr_0", "prd_0", propertyID).Run(func(mock.Arguments) {
				state.Property = papi.Property{
					PropertyName:  propertyName,
					PropertyID:    propertyID,
					GroupID:       "grp_0",
					ContractID:    "ctr_0",
					ProductID:     "prd_0",
					LatestVersion: 1,
				}

				state.Rules = rules
				state.RuleFormat = "v2020-01-01"
				getProperty(propertyID)(state)
				GetVersionResources(propertyID, "ctr_0", "grp_0", 1)(state)
			}).Once()
		}
	}

	propertyLifecycle := func(propertyName, propertyID, groupID string, rules papi.RulesUpdate) BehaviorFunc {
		return func(state *TestState) {
			createProperty(propertyName, propertyID, rules)(state)
			GetVersionResources(propertyID, "ctr_0", "grp_0", 1)(state)
			DeleteProperty(propertyID)(state)
		}
	}

	// propertyLifecycleWithPropertyID covers lifecycle when property_id is set
	propertyLifecycleWithPropertyID := func(propertyName, propertyID, groupID string, rules papi.RulesUpdate) BehaviorFunc {
		return func(state *TestState) {
			state.Property.PropertyID = "prp_0"
			state.Property.LatestVersion = 1
			state.Property.ContractID = "ctr_0"
			state.Property.GroupID = "grp_0"
			state.Property.PropertyName = "test_property"
			state.Rules = rules
			state.RuleFormat = "v2020-01-01"
			getProperty(propertyID)(state)
			GetVersionResources(propertyID, "ctr_0", "grp_0", 1)(state)
			// no deletion since it should be covered by property_bootstrap resource
		}
	}

	propertyLifecycleWithDrift := func(propertyName, propertyID, groupID string, rulesToSend, rulesToReceive papi.RulesUpdate) BehaviorFunc {
		return func(state *TestState) {
			createProperty(propertyName, propertyID, rulesToSend)(state)
			GetVersionResourcesDrift(propertyID, "ctr_0", "grp_0", 1, rulesToReceive)(state)
			DeleteProperty(propertyID)(state)
		}
	}

	importProperty := func(propertyID string) BehaviorFunc {
		return func(state *TestState) {
			// Depending on how much of the import ID is given, the initial property lookup may not have group and contract
			ExpectGetProperty(state.Client, "prp_0", "grp_0", "ctr_0", &state.Property).Maybe()
			ExpectGetProperty(state.Client, "prp_0", "", "", &state.Property).Maybe()
		}
	}

	advanceVersion := func(propertyID string, fromVersion, toVersion int) BehaviorFunc {
		return func(state *TestState) {
			ExpectCreatePropertyVersion(state.Client, propertyID, "grp_0", "ctr_0", fromVersion, toVersion).Once().Run(func(mock.Arguments) {
				state.Property.LatestVersion = toVersion
			}).Run(func(args mock.Arguments) {
				state.Property.LatestVersion = toVersion
				state.VersionItems.Items = append(state.VersionItems.Items,
					papi.PropertyVersionGetItem{
						ProductionStatus: papi.VersionStatusInactive,
						PropertyVersion:  toVersion,
						StagingStatus:    papi.VersionStatusInactive,
					})
			})
			GetVersionResources(propertyID, "ctr_0", "grp_0", toVersion)(state)
		}
	}

	// TestCheckFunc to verify all standard attributes (for Lifecycle tests)
	checkAttrs := func(propertyID, cnameTo, latestVersion, stagingVersion, productionVersion, edgeHostnameId, rules string) resource.TestCheckFunc {
		return resource.ComposeAggregateTestCheckFunc(
			resource.TestCheckResourceAttr("akamai_property.test", "id", propertyID),
			resource.TestCheckResourceAttr("akamai_property.test", "hostnames.0.cname_to", cnameTo),
			resource.TestCheckResourceAttr("akamai_property.test", "hostnames.0.edge_hostname_id", edgeHostnameId),
			resource.TestCheckResourceAttr("akamai_property.test", "latest_version", latestVersion),
			resource.TestCheckResourceAttr("akamai_property.test", "staging_version", stagingVersion),
			resource.TestCheckResourceAttr("akamai_property.test", "production_version", productionVersion),
			resource.TestCheckResourceAttr("akamai_property.test", "name", "test_property"),
			resource.TestCheckResourceAttr("akamai_property.test", "contract_id", "ctr_0"),
			resource.TestCheckResourceAttr("akamai_property.test", "group_id", "grp_0"),
			resource.TestCheckResourceAttr("akamai_property.test", "product_id", "prd_0"),
			resource.TestCheckResourceAttr("akamai_property.test", "rule_warnings.#", "0"),
			resource.TestCheckResourceAttr("akamai_property.test", "rules", rules),
		)
	}

	// addPropertyIDAttrCheck adds resource.TestCheckFunc that checks if property_id attribute was set correctly
	addPropertyIDAttrCheck := func(checks resource.TestCheckFunc, propertyID string) resource.TestCheckFunc {
		return resource.ComposeAggregateTestCheckFunc(
			resource.TestCheckResourceAttr("akamai_property.test", "property_id", propertyID),
			checks,
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
	latestVersionDeactivatedInStaging := LifecycleTestCase{
		Name: "Latest version is active in staging",
		ClientSetup: composeBehaviors(
			propertyLifecycle("test_property", "prp_0", "grp_0",
				papi.RulesUpdate{Rules: papi.Rules{Name: "default"}}),
			getPropertyVersionResources("prp_0", "grp_0", "ctr_0", 1, papi.VersionStatusDeactivated, papi.VersionStatusInactive),
			setHostnames("prp_0", 1, "to.test.domain"),
			advanceVersion("prp_0", 1, 2),
			getPropertyVersionResources("prp_0", "grp_0", "ctr_0", 2, papi.VersionStatusDeactivated, papi.VersionStatusInactive),
			setHostnames("prp_0", 2, "to2.test.domain"),
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
					Config: testutils.LoadFixtureString(t, "%s/step0.tf", FixturePath),
					Check: checkAttrs("prp_0", "to.test.domain", "1", "0", "0", "ehn_123",
						"{\"rules\":{\"name\":\"default\",\"options\":{}}}"),
				},
				{
					PreConfig: func() {
						StagingVersion := 1
						State.Property.StagingVersion = &StagingVersion
					},
					Config: testutils.LoadFixtureString(t, "%s/step1.tf", FixturePath),
					Check: checkAttrs("prp_0", "to2.test.domain", "2", "1", "0", "ehn_123",
						"{\"rules\":{\"name\":\"default\",\"options\":{}}}"),
				},
			}
		},
	}

	// Standard test behavior for cases where the property's latest version is deactivated in production network
	latestVersionDeactivatedInProd := LifecycleTestCase{
		Name: "Latest version is active in production",
		ClientSetup: composeBehaviors(
			propertyLifecycle("test_property", "prp_0", "grp_0",
				papi.RulesUpdate{Rules: papi.Rules{Name: "default"}}),
			getPropertyVersionResources("prp_0", "grp_0", "ctr_0", 1, papi.VersionStatusInactive, papi.VersionStatusDeactivated),
			setHostnames("prp_0", 1, "to.test.domain"),
			advanceVersion("prp_0", 1, 2),
			getPropertyVersionResources("prp_0", "grp_0", "ctr_0", 2, papi.VersionStatusInactive, papi.VersionStatusDeactivated),
			setHostnames("prp_0", 2, "to2.test.domain"),
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
					Config: testutils.LoadFixtureString(t, "%s/step0.tf", FixturePath),
					Check: checkAttrs("prp_0", "to.test.domain", "1", "0", "0", "ehn_123",
						"{\"rules\":{\"name\":\"default\",\"options\":{}}}"),
				},
				{
					PreConfig: func() {
						ProductionVersion := 1
						State.Property.ProductionVersion = &ProductionVersion
					},
					Config: testutils.LoadFixtureString(t, "%s/step1.tf", FixturePath),
					Check: checkAttrs("prp_0", "to2.test.domain", "2", "0", "1", "ehn_123",
						"{\"rules\":{\"name\":\"default\",\"options\":{}}}"),
				},
			}
		},
	}

	// Standard test behavior for cases where the property's latest version is active in staging network
	latestVersionActiveInStaging := LifecycleTestCase{
		Name: "Latest version is active in staging",
		ClientSetup: composeBehaviors(
			propertyLifecycle("test_property", "prp_0", "grp_0",
				papi.RulesUpdate{Rules: papi.Rules{Name: "default"}}),
			getPropertyVersionResources("prp_0", "grp_0", "ctr_0", 1, papi.VersionStatusActive, papi.VersionStatusInactive),
			setHostnames("prp_0", 1, "to.test.domain"),
			advanceVersion("prp_0", 1, 2),
			getPropertyVersionResources("prp_0", "grp_0", "ctr_0", 2, papi.VersionStatusInactive, papi.VersionStatusActive),
			setHostnames("prp_0", 2, "to2.test.domain"),
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
					Config: testutils.LoadFixtureString(t, "%s/step0.tf", FixturePath),
					Check: checkAttrs("prp_0", "to.test.domain", "1", "0", "0", "ehn_123",
						"{\"rules\":{\"name\":\"default\",\"options\":{}}}"),
				},
				{
					PreConfig: func() {
						StagingVersion := 1
						State.Property.StagingVersion = &StagingVersion
					},
					Config: testutils.LoadFixtureString(t, "%s/step1.tf", FixturePath),
					Check: checkAttrs("prp_0", "to2.test.domain", "2", "1", "0", "ehn_123",
						"{\"rules\":{\"name\":\"default\",\"options\":{}}}"),
				},
			}
		},
	}

	// Standard test behavior for cases where the property's latest version is active in production network
	latestVersionActiveInProd := LifecycleTestCase{
		Name: "Latest version is active in production",
		ClientSetup: composeBehaviors(
			propertyLifecycle("test_property", "prp_0", "grp_0",
				papi.RulesUpdate{Rules: papi.Rules{Name: "default"}}),
			getPropertyVersionResources("prp_0", "grp_0", "ctr_0", 1, papi.VersionStatusInactive, papi.VersionStatusActive),
			setHostnames("prp_0", 1, "to.test.domain"),
			advanceVersion("prp_0", 1, 2),
			getPropertyVersionResources("prp_0", "grp_0", "ctr_0", 2, papi.VersionStatusInactive, papi.VersionStatusActive),
			setHostnames("prp_0", 2, "to2.test.domain"),
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
					Config: testutils.LoadFixtureString(t, "%s/step0.tf", FixturePath),
					Check: checkAttrs("prp_0", "to.test.domain", "1", "0", "0", "ehn_123",
						"{\"rules\":{\"name\":\"default\",\"options\":{}}}"),
				},
				{
					PreConfig: func() {
						ProductionVersion := 1
						State.Property.ProductionVersion = &ProductionVersion
					},
					Config: testutils.LoadFixtureString(t, "%s/step1.tf", FixturePath),
					Check: checkAttrs("prp_0", "to2.test.domain", "2", "0", "1", "ehn_123",
						"{\"rules\":{\"name\":\"default\",\"options\":{}}}"),
				},
			}
		},
	}

	// withPropertyID covers case when property was initially created with property_bootstrap resource
	withPropertyID := LifecycleTestCase{
		Name: "Create with propertyID",
		ClientSetup: composeBehaviors(
			propertyLifecycleWithPropertyID("test_property", "prp_0", "grp_0",
				papi.RulesUpdate{Rules: papi.Rules{Name: "default"}}),
			getPropertyVersionResources("prp_0", "grp_0", "ctr_0", 1, papi.VersionStatusInactive, papi.VersionStatusInactive),
			setHostnames("prp_0", 1, "to.test.domain"),
			setHostnames("prp_0", 1, "to2.test.domain"),
		),
		Steps: func(State *TestState, FixturePath string) []resource.TestStep {
			return []resource.TestStep{
				{
					PreConfig: func() {
						State.VersionItems = papi.PropertyVersionItems{Items: []papi.PropertyVersionGetItem{{PropertyVersion: 1, ProductionStatus: papi.VersionStatusInactive}}}
					},
					Config: testutils.LoadFixtureString(t, "%s/step0.tf", FixturePath),
					Check: checkAttrs("prp_0", "to.test.domain", "1", "0", "0", "ehn_123",
						"{\"rules\":{\"name\":\"default\",\"options\":{}}}"),
				},
				{
					Config: testutils.LoadFixtureString(t, "%s/step1.tf", FixturePath),
					Check: addPropertyIDAttrCheck(checkAttrs("prp_0", "to2.test.domain", "1", "0", "0", "ehn_123",
						"{\"rules\":{\"name\":\"default\",\"options\":{}}}"), "prp_0"),
				},
			}
		},
	}

	// Standard test behavior for cases where the property's latest version is not active
	latestVersionNotActive := LifecycleTestCase{
		Name: "Latest version not active",
		ClientSetup: composeBehaviors(
			propertyLifecycle("test_property", "prp_0", "grp_0",
				papi.RulesUpdate{Rules: papi.Rules{Name: "default"}}),
			getPropertyVersionResources("prp_0", "grp_0", "ctr_0", 1, papi.VersionStatusInactive, papi.VersionStatusInactive),
			setHostnames("prp_0", 1, "to.test.domain"),
			setHostnames("prp_0", 1, "to2.test.domain"),
		),
		Steps: func(State *TestState, FixturePath string) []resource.TestStep {
			return []resource.TestStep{
				{
					PreConfig: func() {
						State.VersionItems = papi.PropertyVersionItems{Items: []papi.PropertyVersionGetItem{{PropertyVersion: 1, ProductionStatus: papi.VersionStatusInactive}}}
					},
					Config: testutils.LoadFixtureString(t, "%s/step0.tf", FixturePath),
					Check: checkAttrs("prp_0", "to.test.domain", "1", "0", "0", "ehn_123",
						"{\"rules\":{\"name\":\"default\",\"options\":{}}}"),
				},
				{
					Config: testutils.LoadFixtureString(t, "%s/step1.tf", FixturePath),
					Check: checkAttrs("prp_0", "to2.test.domain", "1", "0", "0", "ehn_123",
						"{\"rules\":{\"name\":\"default\",\"options\":{}}}"),
				},
			}
		},
	}

	// Standard test behavior for cases where the property's latest version is not active
	stagingAndProductionVersionKnownAtPlan := LifecycleTestCase{
		Name: "Latest version not active",
		ClientSetup: composeBehaviors(
			propertyLifecycle("test_property", "prp_0", "grp_0",
				papi.RulesUpdate{Rules: papi.Rules{Name: "default"}}),
			getPropertyVersionResources("prp_0", "grp_0", "ctr_0", 1, papi.VersionStatusInactive, papi.VersionStatusInactive),
			setHostnames("prp_0", 1, "to.test.domain"),
			setHostnames("prp_0", 1, "to2.test.domain"),
		),
		Steps: func(State *TestState, FixturePath string) []resource.TestStep {
			return []resource.TestStep{
				{
					PreConfig: func() {
						State.VersionItems = papi.PropertyVersionItems{Items: []papi.PropertyVersionGetItem{{PropertyVersion: 1, ProductionStatus: papi.VersionStatusInactive}}}
					},
					Config: testutils.LoadFixtureString(t, "%s/step0.tf", FixturePath),
					Check: checkAttrs("prp_0", "to.test.domain", "1", "0", "0", "ehn_123",
						"{\"rules\":{\"name\":\"default\",\"options\":{}}}"),
				},
				{
					Config: testutils.LoadFixtureString(t, "%s/step1.tf", FixturePath),
					Check: checkAttrs("prp_0", "to2.test.domain", "1", "0", "0", "ehn_123",
						"{\"rules\":{\"name\":\"default\",\"options\":{}}}"),
					ConfigPlanChecks: resource.ConfigPlanChecks{PreApply: []plancheck.PlanCheck{testutils.FieldsKnownAtPlan{FieldsKnown: []string{"staging_version", "production_version"}, FieldsUnknown: []string{"latest_version"}}}},
				},
			}
		},
	}

	// This scenario simulates a new version being created outside of terraform and returned on read after the first step (update should be triggered)
	changesMadeOutsideOfTerraform := LifecycleTestCase{
		Name: "Latest version not active",
		ClientSetup: composeBehaviors(
			propertyLifecycle("test_property", "prp_0", "grp_0",
				papi.RulesUpdate{Rules: papi.Rules{Name: "default"}}),
			getPropertyVersionResources("prp_0", "grp_0", "ctr_0", 1, papi.VersionStatusInactive, papi.VersionStatusInactive),
			setHostnames("prp_0", 1, "to.test.domain"),
			getPropertyVersionResources("prp_0", "grp_0", "ctr_0", 2, papi.VersionStatusInactive, papi.VersionStatusInactive),
			GetVersionResources("prp_0", "ctr_0", "grp_0", 2),
			setHostnames("prp_0", 2, "to.test.domain"),
		),
		Steps: func(State *TestState, FixturePath string) []resource.TestStep {
			return []resource.TestStep{
				{
					PreConfig: func() {
						State.VersionItems = papi.PropertyVersionItems{Items: []papi.PropertyVersionGetItem{{PropertyVersion: 1, ProductionStatus: papi.VersionStatusInactive}}}
					},
					Config: testutils.LoadFixtureString(t, "%s/step0.tf", FixturePath),
					Check: checkAttrs("prp_0", "to.test.domain", "1", "0", "0", "ehn_123",
						"{\"rules\":{\"name\":\"default\",\"options\":{}}}"),
				},
				{
					PreConfig: func() {
						State.Property.LatestVersion = 2
						State.Hostnames[0].CnameTo = "changed.test.domain"
					},
					Config: testutils.LoadFixtureString(t, "%s/step0.tf", FixturePath),
					Check: checkAttrs("prp_0", "to.test.domain", "2", "0", "0", "ehn_123",
						"{\"rules\":{\"name\":\"default\",\"options\":{}}}"),
				},
			}
		},
	}

	// Standard test behavior for cases where the property's latest version is active in staging network
	noDiff := LifecycleTestCase{
		Name: "No diff found in update",
		ClientSetup: composeBehaviors(
			propertyLifecycle("test_property", "prp_0", "grp_0",
				papi.RulesUpdate{Rules: papi.Rules{Children: []papi.Rules{{Name: "Default CORS Policy", CriteriaMustSatisfy: papi.RuleCriteriaMustSatisfyAll}}}}),
			getPropertyVersionResources("prp_0", "grp_0", "ctr_0", 1, papi.VersionStatusInactive, papi.VersionStatusInactive),
			setHostnames("prp_0", 1, "to.test.domain"),
			updateRuleTree("prp_0", "ctr_0", "grp_0", 1,
				&papi.RulesUpdate{Rules: papi.Rules{Children: []papi.Rules{{CriteriaMustSatisfy: papi.RuleCriteriaMustSatisfyAll, Name: "Default CORS Policy"}}}}),
		),
		Steps: func(State *TestState, FixturePath string) []resource.TestStep {
			return []resource.TestStep{
				{
					PreConfig: func() {
						State.VersionItems = papi.PropertyVersionItems{Items: []papi.PropertyVersionGetItem{{PropertyVersion: 1, ProductionStatus: papi.VersionStatusInactive}}}
					},
					Config: testutils.LoadFixtureString(t, "%s/step0.tf", FixturePath),
					Check: checkAttrs("prp_0", "to.test.domain", "1", "0", "0", "ehn_123",
						`{"rules":{"children":[{"name":"Default CORS Policy","options":{},"criteriaMustSatisfy":"all"}],"name":"","options":{}}}`),
				},
				{
					Config: testutils.LoadFixtureString(t, "%s/step1.tf", FixturePath),
					Check: checkAttrs("prp_0", "to.test.domain", "1", "0", "0", "ehn_123",
						`{"rules":{"children":[{"name":"Default CORS Policy","options":{},"criteriaMustSatisfy":"all"}],"name":"","options":{}}}`),
				},
			}
		},
	}

	diffCPCode := LifecycleTestCase{
		Name: "Diff cpCode.cpCodeLimits",
		ClientSetup: composeBehaviors(
			propertyLifecycleWithDrift("test_property", "prp_0", "grp_0",
				papi.RulesUpdate{
					Rules: papi.Rules{Behaviors: []papi.RuleBehavior{
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
						Name: "default"}},
				papi.RulesUpdate{
					Rules: papi.Rules{Behaviors: []papi.RuleBehavior{
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
						Name: "default"}},
			),
			setHostnames("prp_0", 1, "to.test.domain"),
			updateRuleTree("prp_0", "ctr_0", "grp_0", 1,
				&papi.RulesUpdate{Rules: papi.Rules{Behaviors: []papi.RuleBehavior{
					{Name: "cpCode",
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
					Name: "default"}}),
			getPropertyVersionResources("prp_0", "grp_0", "ctr_0", 1, papi.VersionStatusInactive, papi.VersionStatusInactive),
		),
		Steps: func(State *TestState, FixturePath string) []resource.TestStep {
			return []resource.TestStep{
				{
					Config: testutils.LoadFixtureString(t, "%s/step0.tf", FixturePath),
					Check: checkAttrs("prp_0", "to.test.domain", "1", "0", "0", "ehn_123",
						`{"rules":{"behaviors":[{"name":"cpCode","options":{"value":{"cpCodeLimits":null,"description":"CliTerraformCPCode","id":1050269,"name":"DevExpCliTerraformPapiAsSchemaTest","products":["Web_App_Accel"]}}}],"name":"default","options":{}}}`),
				},
			}
		},
	}

	/*
		rulesCustomDiff tests rulesCustomDiff function which is in resource_akamai_property.go file.
		There is an additional field "options":{} in expected attributes, because with UpdateRuleTree(ctx, req) function
		this field added automatically into response, even if it does not exist in rules.
	*/
	rulesCustomDiff := LifecycleTestCase{
		Name: "Diff is only in behaviours.options.ttl",
		ClientSetup: composeBehaviors(
			propertyLifecycle("test_property", "prp_0", "grp_0",
				papi.RulesUpdate{Rules: papi.Rules{Behaviors: []papi.RuleBehavior{{Name: "caching",
					Options: papi.RuleOptionsMap{"behavior": "MAX_AGE", "mustRevalidate": false, "ttl": "12d"}}},
					Name: "default"}}),
			getPropertyVersionResources("prp_0", "grp_0", "ctr_0", 1, papi.VersionStatusInactive, papi.VersionStatusInactive),
			setHostnames("prp_0", 1, "to.test.domain"),
			updateRuleTree("prp_0", "ctr_0", "grp_0", 1,
				&papi.RulesUpdate{Rules: papi.Rules{Behaviors: []papi.RuleBehavior{{Name: "caching",
					Options: papi.RuleOptionsMap{"behavior": "MAX_AGE", "mustRevalidate": false, "ttl": "12d"}}},
					Name: "default"}}),
			updateRuleTree("prp_0", "ctr_0", "grp_0", 1,
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
					Config: testutils.LoadFixtureString(t, "%s/step0.tf", FixturePath),
					Check: checkAttrs("prp_0", "to.test.domain", "1", "0", "0", "ehn_123",
						`{"rules":{"behaviors":[{"name":"caching","options":{"behavior":"MAX_AGE","mustRevalidate":false,"ttl":"12d"}}],"name":"default","options":{}}}`),
				},
				{
					Config: testutils.LoadFixtureString(t, "%s/step1.tf", FixturePath),
					Check: checkAttrs("prp_0", "to.test.domain", "1", "0", "0", "ehn_123",
						`{"rules":{"behaviors":[{"name":"caching","options":{"behavior":"MAX_AGE","mustRevalidate":false,"ttl":"13d"}}],"name":"default","options":{}}}`),
				},
			}
		},
	}

	noDiffForHostnames := LifecycleTestCase{
		Name: "No diff found in update",
		ClientSetup: composeBehaviors(
			propertyLifecycle("test_property", "prp_0", "grp_0",
				papi.RulesUpdate{Rules: papi.Rules{Children: []papi.Rules{{Name: "Default CORS Policy", CriteriaMustSatisfy: papi.RuleCriteriaMustSatisfyAll}}}}),
			getPropertyVersionResources("prp_0", "grp_0", "ctr_0", 1, papi.VersionStatusInactive, papi.VersionStatusInactive),
			setTwoHostnames("prp_0", 1, "from1.test.domain", "to1.test.domain", "from2.test.domain", "to2.test.domain"),
			updateRuleTree("prp_0", "ctr_0", "grp_0", 1,
				&papi.RulesUpdate{Rules: papi.Rules{Children: []papi.Rules{{CriteriaMustSatisfy: papi.RuleCriteriaMustSatisfyAll, Name: "Default CORS Policy"}}}}),
		),
		Steps: func(State *TestState, FixturePath string) []resource.TestStep {
			return []resource.TestStep{
				{
					PreConfig: func() {
						State.VersionItems = papi.PropertyVersionItems{Items: []papi.PropertyVersionGetItem{{PropertyVersion: 1, ProductionStatus: papi.VersionStatusInactive}}}
					},
					Config: testutils.LoadFixtureString(t, "%s/step0.tf", FixturePath),
					Check: checkAttrs("prp_0", "to1.test.domain", "1", "0", "0", "ehn_123",
						`{"rules":{"children":[{"name":"Default CORS Policy","options":{},"criteriaMustSatisfy":"all"}],"name":"","options":{}}}`),
				},
				{
					Config: testutils.LoadFixtureString(t, "%s/step1.tf", FixturePath),
					Check: checkAttrs("prp_0", "to1.test.domain", "1", "0", "0", "ehn_123",
						`{"rules":{"children":[{"name":"Default CORS Policy","options":{},"criteriaMustSatisfy":"all"}],"name":"","options":{}}}`),
				},
			}
		},
	}

	variablesInRuleTree := LifecycleTestCase{
		Name: "Variables in property rule tree",
		ClientSetup: composeBehaviors(
			propertyLifecycle("test_property", "prp_0", "grp_0", papi.RulesUpdate{Rules: papi.Rules{Name: "default"}}),
			getPropertyVersionResources("prp_0", "grp_0", "ctr_0", 1, papi.VersionStatusInactive, papi.VersionStatusInactive),
			setHostnames("prp_0", 1, "to.test.domain"),
			getPropertyVersionResources("prp_0", "grp_0", "ctr_0", 1, papi.VersionStatusInactive, papi.VersionStatusInactive),
			GetVersionResources("prp_0", "ctr_0", "grp_0", 1),
			updateRuleTree("prp_0", "ctr_0", "grp_0", 1, updateRuleTreeWithVariablesStep0()),
			updateRuleTree("prp_0", "ctr_0", "grp_0", 1, updateRuleTreeWithVariablesStep1()),
		),
		Steps: func(State *TestState, FixturePath string) []resource.TestStep {
			return []resource.TestStep{
				{
					PreConfig: func() {
						State.VersionItems = papi.PropertyVersionItems{Items: []papi.PropertyVersionGetItem{{PropertyVersion: 1, ProductionStatus: papi.VersionStatusInactive}}}
					},
					Config: testutils.LoadFixtureString(t, "%s/step0.tf", FixturePath),
					Check: checkAttrs("prp_0", "to.test.domain", "1", "0", "0", "ehn_123",
						"{\"rules\":{\"behaviors\":[{\"name\":\"origin\",\"options\":{\"cacheKeyHostname\":\"REQUEST_HOST_HEADER\",\"compress\":true,\"enableTrueClientIp\":true,\"forwardHostHeader\":\"REQUEST_HOST_HEADER\",\"hostname\":\"test.domain\",\"httpPort\":80,\"httpsPort\":443,\"originCertificate\":\"\",\"originSni\":true,\"originType\":\"CUSTOMER\",\"ports\":\"\",\"trueClientIpClientSetting\":false,\"trueClientIpHeader\":\"True-Client-IP\",\"verificationMode\":\"PLATFORM_SETTINGS\"}}],\"children\":[{\"behaviors\":[{\"name\":\"baseDirectory\",\"options\":{\"value\":\"/smth/\"}}],\"criteria\":[{\"name\":\"requestHeader\",\"options\":{\"headerName\":\"Accept-Encoding\",\"matchCaseSensitiveValue\":true,\"matchOperator\":\"IS_ONE_OF\",\"matchWildcardName\":false,\"matchWildcardValue\":false}}],\"name\":\"change fwd path\",\"options\":{},\"criteriaMustSatisfy\":\"all\"},{\"behaviors\":[{\"name\":\"caching\",\"options\":{\"behavior\":\"MAX_AGE\",\"mustRevalidate\":false,\"ttl\":\"1m\"}}],\"name\":\"caching\",\"options\":{},\"criteriaMustSatisfy\":\"any\"}],\"comments\":\"The behaviors in the Default Rule apply to all requests for the property hostname(s) unless another rule overrides the Default Rule settings.\",\"name\":\"default\",\"options\":{},\"variables\":[{\"description\":\"\",\"hidden\":true,\"name\":\"TEST_EMPTY_FIELDS\",\"sensitive\":false,\"value\":\"\"},{\"description\":null,\"hidden\":true,\"name\":\"TEST_NIL_FIELD\",\"sensitive\":false,\"value\":\"\"}]}}"),
				},
				{
					PreConfig: func() {
						State.Property.LatestVersion = 1
					},
					Config: testutils.LoadFixtureString(t, "%s/step1.tf", FixturePath),
					Check: checkAttrs("prp_0", "to.test.domain", "1", "0", "0", "ehn_123",
						"{\"rules\":{\"behaviors\":[{\"name\":\"origin\",\"options\":{\"cacheKeyHostname\":\"REQUEST_HOST_HEADER\",\"compress\":true,\"enableTrueClientIp\":true,\"forwardHostHeader\":\"REQUEST_HOST_HEADER\",\"hostname\":\"test.domain\",\"httpPort\":80,\"httpsPort\":443,\"originCertificate\":\"\",\"originSni\":true,\"originType\":\"CUSTOMER\",\"ports\":\"\",\"trueClientIpClientSetting\":false,\"trueClientIpHeader\":\"True-Client-IP\",\"verificationMode\":\"PLATFORM_SETTINGS\"}}],\"children\":[{\"behaviors\":[{\"name\":\"baseDirectory\",\"options\":{\"value\":\"/smth/\"}}],\"criteria\":[{\"name\":\"requestHeader\",\"options\":{\"headerName\":\"Accept-Encoding\",\"matchCaseSensitiveValue\":true,\"matchOperator\":\"IS_ONE_OF\",\"matchWildcardName\":false,\"matchWildcardValue\":false}}],\"name\":\"change fwd path\",\"options\":{},\"criteriaMustSatisfy\":\"all\"},{\"behaviors\":[{\"name\":\"caching\",\"options\":{\"behavior\":\"MAX_AGE\",\"mustRevalidate\":false,\"ttl\":\"1m\"}}],\"name\":\"caching\",\"options\":{},\"criteriaMustSatisfy\":\"any\"}],\"comments\":\"The behaviors in the Default Rule apply to all requests for the property hostname(s) unless another rule overrides the Default Rule settings.\",\"name\":\"default\",\"options\":{},\"variables\":[{\"description\":\"\",\"hidden\":true,\"name\":\"TEST_EMPTY_FIELDS\",\"sensitive\":false,\"value\":\"\"},{\"description\":\"\",\"hidden\":true,\"name\":\"TEST_NIL_FIELD\",\"sensitive\":false,\"value\":\"\"}]}}"),
				},
			}
		},
	}

	// Test Schema Configuration

	// Run a test case to verify schema validations
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

	// Test Lifecycle

	// Run a happy-path test case that goes through a complete create-update-destroy cycle
	assertLifecycle := func(t *testing.T, name, variant string, tc LifecycleTestCase) func(t *testing.T) {

		fixturePrefix := fmt.Sprintf("testdata/%s/Lifecycle/%s", t.Name(), variant)

		return func(t *testing.T) {
			client := &papi.Mock{}
			client.Test(T{t})
			State := &TestState{Client: client}
			tc.ClientSetup(State)

			useClient(client, nil, func() {
				resource.UnitTest(t, resource.TestCase{
					ProtoV6ProviderFactories: testutils.NewProtoV6ProviderFactory(NewSubprovider()),
					IsUnitTest:               true,
					Steps:                    tc.Steps(State, fixturePrefix),
				})
			})

			client.AssertExpectations(t)
		}
	}

	// Test Import
	// Run a test case that verifies the resource can be imported by the given ID
	assertImportableWithOptions := func(t *testing.T, testName, importID, fileName, rules string, setup []BehaviorFunc) func(t *testing.T) {

		fixturePath := fmt.Sprintf("testdata/%s/Importable/%s", t.Name(), fileName)

		return func(t *testing.T) {

			client := &papi.Mock{}
			client.Test(T{t})

			parameters := strings.Split(importID, ",")
			var propertyBootstrap bool
			if parameters[len(parameters)-1] == "property-bootstrap" {
				propertyBootstrap = true
				parameters = parameters[:len(parameters)-1]
			}
			numberParameters := len(parameters)
			lastParameter := parameters[len(parameters)-1]
			if propertyBootstrap {
				setup = append(setup, propertyLifecycleWithPropertyID("test_property", "prp_0", "grp_0",
					papi.RulesUpdate{Rules: papi.Rules{Name: "default"}}))
			} else {
				setup = append(setup, propertyLifecycle("test_property", "prp_0", "grp_0",
					papi.RulesUpdate{Rules: papi.Rules{Name: "default"}}))
			}
			setup = append(setup,
				getPropertyVersionResources("prp_0", "grp_0", "ctr_0", 1, papi.VersionStatusInactive, papi.VersionStatusInactive),
				setHostnames("prp_0", 1, "to.test.domain"),
				importProperty("prp_0"),
			)
			if (numberParameters == 2 || numberParameters == 4) && !isDefaultVersion(lastParameter) {
				var contractID, groupID string
				if numberParameters == 4 {
					contractID = "ctr_0"
					groupID = "grp_0"
				}
				if numberParameters == 2 {
					setup = append(setup, getPropertyVersions("prp_0", "test_property", "ctr_0", "grp_0"))
				}
				setup = append(setup, getPropertyVersions("prp_0", "test_property", contractID, groupID))
			}
			s := composeBehaviors(setup...)
			tc := LifecycleTestCase{
				Name:        "Importable",
				ClientSetup: s,
				Steps: func(State *TestState, _ string) []resource.TestStep {
					return []resource.TestStep{
						{
							Config: testutils.LoadFixtureString(t, fixturePath),
							Check:  checkAttrs("prp_0", "to.test.domain", "1", "0", "0", "ehn_123", rules),
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
							Config: testutils.LoadFixtureString(t, fixturePath),
							Check:  checkAttrs("prp_0", "to.test.domain", "1", "1", "0", "ehn_123", rules),
						},
						{
							ImportState:             true,
							ImportStateVerify:       true,
							ImportStateId:           importID,
							ResourceName:            "akamai_property.test",
							Config:                  testutils.LoadFixtureString(t, fixturePath),
							ImportStateVerifyIgnore: []string{"product", "read_version"},
							Check:                   addPropertyIDAttrCheck(checkAttrs("prp_0", "to.test.domain", "1", "1", "0", "ehn_123", rules), "prp_0"),
						},
					}
				},
			}
			State := &TestState{Client: client}
			tc.ClientSetup(State)
			useClient(client, nil, func() {
				resource.UnitTest(t, resource.TestCase{
					ProtoV6ProviderFactories: testutils.NewProtoV6ProviderFactory(NewSubprovider()),
					Steps:                    tc.Steps(State, ""),
				})
			})

			client.AssertExpectations(t)
		}
	}

	assertImportable := func(t *testing.T, testName, importID string) func(t *testing.T) {
		return assertImportableWithOptions(t, testName, importID, "importable.tf", "{\"rules\":{\"name\":\"default\",\"options\":{}}}", []BehaviorFunc{})
	}

	// assertImportableWithBootstrap covers imports when property-bootstrap flag is provided
	assertImportableWithBootstrap := func(t *testing.T, testName, importID string) func(t *testing.T) {
		return assertImportableWithOptions(t, testName, importID, "importable-with-bootstrap.tf", "{\"rules\":{\"name\":\"default\",\"options\":{}}}", []BehaviorFunc{})
	}

	suppressLogging(t, func() {

		// Test Schema Configuration

		t.Run("Schema Configuration Error: name not given", assertConfigError(t, "name not given", `"name" is required`))
		t.Run("Schema Configuration Error: contract_id not given", assertConfigError(t, "contract_id not given", `Missing required argument`))
		t.Run("Schema Configuration Error: group_id not given", assertConfigError(t, "group_id not given", `Missing required argument`))
		t.Run("Schema Configuration Error: product_id not given", assertConfigError(t, "product_id not given", `Missing required argument`))
		t.Run("Schema Configuration Error: invalid json rules", assertConfigError(t, "invalid json rules", `rules are not valid JSON`))
		t.Run("Schema Configuration Error: invalid name given", assertConfigError(t, "invalid name given", `a name must only contain letters, numbers, and these characters: . _ -`))
		t.Run("Schema Configuration Error: name given too long", assertConfigError(t, "name given too long", `a name must be longer than 0 characters and shorter than 86 characters`))

		// Test Lifecycle

		t.Run("Lifecycle: create with propertyID", assertLifecycle(t, t.Name(), "with-propertyID", withPropertyID))
		t.Run("Lifecycle: latest version is not active (normal)", assertLifecycle(t, t.Name(), "normal", latestVersionNotActive))
		t.Run("Lifecycle: latest version is active in staging (normal)", assertLifecycle(t, t.Name(), "normal", latestVersionActiveInStaging))
		t.Run("Lifecycle: latest version is active in production (normal)", assertLifecycle(t, t.Name(), "normal", latestVersionActiveInProd))
		t.Run("Lifecycle: latest version is deactivated in staging (normal)", assertLifecycle(t, t.Name(), "normal", latestVersionDeactivatedInStaging))
		t.Run("Lifecycle: latest version is deactivated in production (normal)", assertLifecycle(t, t.Name(), "normal", latestVersionDeactivatedInProd))
		t.Run("Lifecycle: latest version is not active (contract_id without prefix)", assertLifecycle(t, t.Name(), "contract_id without prefix", latestVersionNotActive))
		t.Run("Lifecycle: latest version active in staging (contract_id without prefix)", assertLifecycle(t, t.Name(), "contract_id without prefix", latestVersionActiveInStaging))
		t.Run("Lifecycle: latest version active in production (contract_id without prefix)", assertLifecycle(t, t.Name(), "contract_id without prefix", latestVersionActiveInProd))
		t.Run("Lifecycle: latest version is not active (group_id without prefix)", assertLifecycle(t, t.Name(), "group_id without prefix", latestVersionNotActive))
		t.Run("Lifecycle: latest version is active in staging (group_id without prefix)", assertLifecycle(t, t.Name(), "group_id without prefix", latestVersionActiveInStaging))
		t.Run("Lifecycle: latest version is active in production (group_id without prefix)", assertLifecycle(t, t.Name(), "group_id without prefix", latestVersionActiveInProd))
		t.Run("Lifecycle: latest version is not active (product_id without prefix)", assertLifecycle(t, t.Name(), "product_id without prefix", latestVersionNotActive))
		t.Run("Lifecycle: latest version is active in staging (product_id without prefix)", assertLifecycle(t, t.Name(), "product_id without prefix", latestVersionActiveInStaging))
		t.Run("Lifecycle: latest version is active in production (product_id without prefix)", assertLifecycle(t, t.Name(), "product_id without prefix", latestVersionActiveInProd))
		t.Run("Lifecycle: no diff", assertLifecycle(t, t.Name(), "no diff", noDiff))
		t.Run("Lifecycle: diff cpCode", assertLifecycle(t, t.Name(), "rules diff cpcode", diffCPCode))
		t.Run("Lifecycle: rules custom diff", assertLifecycle(t, t.Name(), "rules custom diff", rulesCustomDiff))
		t.Run("Lifecycle: no diff for hostnames (hostnames)", assertLifecycle(t, t.Name(), "hostnames", noDiffForHostnames))
		t.Run("Lifecycle: new version changed on server", assertLifecycle(t, t.Name(), "new version changed on server", changesMadeOutsideOfTerraform))
		t.Run("Lifecycle: rules with variables", assertLifecycle(t, t.Name(), "rules with variables", variablesInRuleTree))
		t.Run("Lifecycle: Verify staging_version and production_version known at plan", assertLifecycle(t, t.Name(), "normal", stagingAndProductionVersionKnownAtPlan))

		// Test Import

		t.Run("Importable: property_id with ds", assertImportableWithOptions(t, "property_id", "prp_0", "importable_with_property_rules_builder.tf",
			"{\"rules\":{\"behaviors\":[{\"name\":\"mPulse\",\"options\":{\"configOverride\":\"no new line\"}},{\"name\":\"mPulse\",\"options\":{\"configOverride\":\"\"}},{\"name\":\"mPulse\",\"options\":{\"configOverride\":\"\\n\\tline with new line before and after + tab\\n\"}}],\"name\":\"default\",\"options\":{}}}",
			[]BehaviorFunc{
				updateRuleTree("prp_0", "ctr_0", "grp_0", 1,
					&papi.RulesUpdate{
						Rules: papi.Rules{
							Name: "default",
							Behaviors: []papi.RuleBehavior{
								{Name: "mPulse", Options: papi.RuleOptionsMap{"configOverride": "no new line"}},
								{Name: "mPulse", Options: papi.RuleOptionsMap{"configOverride": ""}},
								{Name: "mPulse", Options: papi.RuleOptionsMap{"configOverride": "\n\tline with new line before and after + tab\n"}},
							},
						}})},
		))
		t.Run("Importable: property_id with property-bootstrap", assertImportableWithBootstrap(t, "property_id", "prp_0,property-bootstrap"))
		t.Run("Importable: property_id", assertImportable(t, "property_id", "prp_0"))
		t.Run("Importable: property_id and ver_# version", assertImportable(t, "property_id and ver_# version", "prp_0,ver_1"))
		t.Run("Importable: property_id and # version", assertImportable(t, "property_id and # version", "prp_0,1"))
		t.Run("Importable: property_id and latest", assertImportable(t, "property_id and latest", "prp_0,latest"))
		t.Run("Importable: property_id and network", assertImportable(t, "property_id and network", "prp_0,staging"))
		t.Run("Importable: unprefixed property_id", assertImportable(t, "unprefixed property_id", "0"))
		t.Run("Importable: unprefixed property_id and # version", assertImportable(t, "unprefixed property_id and # version", "0,1"))
		t.Run("Importable: unprefixed property_id and ver_# version", assertImportable(t, "unprefixed property_id and ver_# version", "0,ver_1"))
		t.Run("Importable: unprefixed property_id and network", assertImportable(t, "unprefixed property_id and network", "0,p"))
		t.Run("Importable: property_id and contract_id and group_id", assertImportable(t, "property_id and contract_id and group_id", "prp_0,ctr_0,grp_0"))
		t.Run("Importable: property_id, contract_id, group_id and empty version", assertImportable(t, "property_id, contract_id, group_id and empty version", "prp_0,ctr_0,grp_0,"))
		t.Run("Importable: property_id, contract_id, group_id and latest", assertImportable(t, "property_id, contract_id, group_id and latest", "prp_0,ctr_0,grp_0,latest"))
		t.Run("Importable: property_id, contract_id, group_id and ver_# version", assertImportable(t, "property_id, contract_id, group_id and ver_# version", "prp_0,ctr_0,grp_0,ver_1"))
		t.Run("Importable: property_id, contract_id, group_id and # version", assertImportable(t, "property_id, contract_id, group_id and # version", "prp_0,ctr_0,grp_0,1"))
		t.Run("Importable: property_id, contract_id, group_id and network", assertImportable(t, "property_id, contract_id, group_id and network", "prp_0,ctr_0,grp_0,staging"))
		t.Run("Importable: unprefixed property_id and contract_id and group_id", assertImportable(t, "unprefixed property_id and contract_id and group_id", "0,0,0"))
		t.Run("Importable: unprefixed property_id and contract_id, group_id and # version", assertImportable(t, "unprefixed property_id and contract_id, group_id and # version", "0,0,0,1"))
		t.Run("Importable: unprefixed property_id and contract_id, group_id and ver_# version", assertImportable(t, "unprefixed property_id and contract_id, group_id and ver_# version", "0,0,0,ver_1"))
		t.Run("Importable: unprefixed property_id and contract_id, group_id and latest", assertImportable(t, "unprefixed property_id and contract_id, group_id and latest", "0,0,0,latest"))
		t.Run("Importable: unprefixed property_id and contract_id, group_id and network", assertImportable(t, "unprefixed property_id and contract_id, group_id and network", "0,0,0,production"))

		// Test Delete

		t.Run("property is destroyed and recreated when name is changed", func(t *testing.T) {
			client := &papi.Mock{}
			client.Test(T{t})

			setup := composeBehaviors(
				propertyLifecycle("test_property", "prp_0", "grp_0",
					papi.RulesUpdate{Rules: papi.Rules{Name: "default"}}),
				getPropertyVersionResources("prp_0", "grp_0", "ctr_0", 1, papi.VersionStatusInactive, papi.VersionStatusInactive),
				propertyLifecycle("renamed_property", "prp_1", "grp_0",
					papi.RulesUpdate{Rules: papi.Rules{Name: "default"}}),
				getPropertyVersionResources("prp_1", "grp_0", "ctr_0", 1, papi.VersionStatusInactive, papi.VersionStatusInactive),
				setHostnames("prp_0", 1, "to.test.domain"),
				setHostnames("prp_1", 1, "to2.test.domain"),
			)
			setup(&TestState{Client: client})

			useClient(client, nil, func() {
				resource.UnitTest(t, resource.TestCase{
					ProtoV6ProviderFactories: testutils.NewProtoV6ProviderFactory(NewSubprovider()),
					Steps: []resource.TestStep{
						{
							Config: testutils.LoadFixtureString(t, "testdata/%s-step0.tf", t.Name()),
							Check: checkAttrs("prp_0", "to.test.domain", "1", "0", "0", "ehn_123",
								"{\"rules\":{\"name\":\"default\",\"options\":{}}}"),
						},
						{
							Config: testutils.LoadFixtureString(t, "testdata/%s-step1.tf", t.Name()),
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
			client := &papi.Mock{}
			client.Test(T{t})

			setup := composeBehaviors(
				createProperty("test_property", "prp_0", papi.RulesUpdate{Rules: papi.Rules{Name: "default"}}),
				getProperty("prp_0"),
				GetVersionResources("prp_0", "ctr_0", "grp_0", 1),
				getPropertyVersionResources("prp_0", "grp_0", "ctr_0", 1, "ctr_0", "grp_0"),
				setHostnames("prp_0", 1, "to.test.domain"),
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
					ProtoV6ProviderFactories: testutils.NewProtoV6ProviderFactory(NewSubprovider()),
					Steps: []resource.TestStep{
						{
							Config: testutils.LoadFixtureString(t, "testdata/%s/step0.tf", t.Name()),
							Check: checkAttrs("prp_0", "to.test.domain", "1", "0", "0", "ehn_123",
								"{\"rules\":{\"name\":\"default\",\"options\":{}}}"),
						},
						{
							Config:      testutils.LoadFixtureString(t, "testdata/%s/step1.tf", t.Name()),
							ExpectError: regexp.MustCompile(`cannot remove active property`),
						},
					},
				})
			})

			client.AssertExpectations(t)
		})

		// Test validation

		t.Run("error validations when updating property with rules tree", func(t *testing.T) {
			client := &papi.Mock{}
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
					ProtoV6ProviderFactories: testutils.NewProtoV6ProviderFactory(NewSubprovider()),
					Steps: []resource.TestStep{
						{
							Config: testutils.LoadFixtureString(t, "testdata/TestResProperty/property_update_with_validation_error_for_rules.tf"),
							Check: resource.ComposeAggregateTestCheckFunc(
								resource.TestCheckNoResourceAttr("akamai_property.test", "rules")),
							ExpectError: regexp.MustCompile(`validation.required_behavior`),
						},
					},
				})
			})

			client.AssertExpectations(t)
		})
		t.Run("validation warning when creating property with rules tree", func(t *testing.T) {
			client := &papi.Mock{}
			client.Test(T{t})
			ExpectCreateProperty(
				client, "test_property", "grp_0",
				"ctr_0", "prd_0", "prp_1",
			)
			rules := papi.Rules{
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
			}
			var req = papi.UpdateRulesRequest{
				PropertyID:      "prp_1",
				ContractID:      "ctr_0",
				GroupID:         "grp_0",
				PropertyVersion: 1,
				Rules:           papi.RulesUpdate{Rules: rules},
				ValidateRules:   true,
			}
			warning := papi.RuleWarnings{
				Type:          "https://problems.luna.akamaiapis.net/papi/v0/validation/validation_message.ip_address_origin",
				ErrorLocation: "#/rules/behaviors/1",
				Detail:        "Using an IP address for the `Origin Server` is not recommended. IP addresses may be changed or reassigned without notice which can severely impact your property or cause a DoS. Please use a properly formatted hostname instead.",
			}
			client.On("UpdateRuleTree", AnyCTX, req).Return(&papi.UpdateRulesResponse{
				AccountID:       "",
				ContractID:      "ctr_0",
				Comments:        "",
				GroupID:         "grp_0",
				PropertyID:      "prp_1",
				PropertyVersion: 1,
				Etag:            "",
				RuleFormat:      "",
				Rules:           rules,
				Errors:          nil,
				Warnings:        []papi.RuleWarnings{warning},
			}, nil).Once()
			ExpectGetProperty(
				client, "prp_1", "grp_0", "ctr_0",
				&papi.Property{
					PropertyID: "prp_1", GroupID: "grp_0", ContractID: "ctr_0", LatestVersion: 1,
					PropertyName: "test_property",
				},
			)

			ExpectGetPropertyVersionHostnames(
				client, "prp_1", "grp_0", "ctr_0", 1,
				&[]papi.Hostname{},
			).Times(2)
			ruleFormat := ""
			ExpectGetRuleTree(
				client, "prp_1", "grp_0", "ctr_0", 1,
				&papi.RulesUpdate{
					Rules: rules,
				}, &ruleFormat, nil, []*papi.Error{
					{
						Type:          "https://problems.luna.akamaiapis.net/papi/v0/validation/validation_message.ip_address_origin",
						ErrorLocation: "#/rules/behaviors/1",
						Detail:        "Using an IP address for the `Origin Server` is not recommended. IP addresses may be changed or reassigned without notice which can severely impact your property or cause a DoS. Please use a properly formatted hostname instead.",
					},
				})
			ExpectGetPropertyVersion(client, "prp_1", "grp_0", "ctr_0", 1, papi.VersionStatusInactive, papi.VersionStatusInactive)

			ExpectRemoveProperty(client, "prp_1", "ctr_0", "grp_0")
			useClient(client, nil, func() {
				resource.UnitTest(t, resource.TestCase{
					ProtoV6ProviderFactories: testutils.NewProtoV6ProviderFactory(NewSubprovider()),
					Steps: []resource.TestStep{
						{
							Config: testutils.LoadFixtureString(t, "testdata/TestResProperty/property_with_validation_warning_for_rules.tf"),
							Check: resource.ComposeAggregateTestCheckFunc(
								resource.TestCheckResourceAttr("akamai_property.test", "rule_warnings.0.detail", "Using an IP address for the `Origin Server` is not recommended. IP addresses may be changed or reassigned without notice which can severely impact your property or cause a DoS. Please use a properly formatted hostname instead.")),
						},
					},
				})
			})

			client.AssertExpectations(t)
		})

		t.Run("validation - when updating a property hostnames to empty it should return error", func(t *testing.T) {
			client := &papi.Mock{}
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
				&papi.RulesUpdate{}, &ruleFormat, nil, nil)

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
					ProtoV6ProviderFactories: testutils.NewProtoV6ProviderFactory(NewSubprovider()),
					Steps: []resource.TestStep{
						{
							Config: testutils.LoadFixtureString(t, "testdata/TestResProperty/CreationUpdateNoHostnames/creation/property_create.tf"),
							Check:  resource.TestCheckResourceAttr("akamai_property.test", "id", "prp_0"),
						},
						{
							Config: testutils.LoadFixtureString(t, "testdata/TestResProperty/CreationUpdateNoHostnames/update/property_update.tf"),
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
			client := &papi.Mock{}
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
					Remaining:  ptr.To(0),
					LimitKey:   "DEFAULT_CERTS_PER_CONTRACT",
				},
			).Once()

			ExpectRemoveProperty(client, "prp_0", "ctr_0", "grp_0")

			useClient(client, nil, func() {
				resource.UnitTest(t, resource.TestCase{
					ProtoV6ProviderFactories: testutils.NewProtoV6ProviderFactory(NewSubprovider()),
					Steps: []resource.TestStep{
						{
							Config:      testutils.LoadFixtureString(t, "testdata/TestResProperty/CreationUpdateNoHostnames/creation/property_create.tf"),
							Check:       resource.TestCheckResourceAttr("akamai_property.test", "id", "prp_0"),
							ExpectError: regexp.MustCompile("updating hostnames: not possible to use cert_provisioning_type = 'DEFAULT' as the limit for DEFAULT certificates has been reached"),
						},
					},
				})
			})
		})

		t.Run("validation - when updating a property hostnames with cert_provisioning_type = 'DEFAULT' not having enabled secure-by-default it should return error", func(t *testing.T) {
			client := &papi.Mock{}
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
					ProtoV6ProviderFactories: testutils.NewProtoV6ProviderFactory(NewSubprovider()),
					Steps: []resource.TestStep{
						{
							Config:      testutils.LoadFixtureString(t, "testdata/TestResProperty/CreationUpdateNoHostnames/creation/property_create.tf"),
							Check:       resource.TestCheckResourceAttr("akamai_property.test", "id", "prp_0"),
							ExpectError: regexp.MustCompile("updating hostnames: not possible to use cert_provisioning_type = 'DEFAULT' as secure-by-default is not enabled in this account"),
						},
					},
				})
			})
		})

		// Other tests

		t.Run("error when the given group is not found", func(t *testing.T) {
			client := &papi.Mock{}
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
					ProtoV6ProviderFactories: testutils.NewProtoV6ProviderFactory(NewSubprovider()),
					Steps: []resource.TestStep{{
						Config:      testutils.LoadFixtureString(t, "testdata/TestResProperty/Creation/property.tf"),
						ExpectError: regexp.MustCompile("group not found: grp_0"),
					}},
				})
			})

			client.AssertExpectations(t)
		})

		t.Run("error when creating property with non-unique name", func(t *testing.T) {
			client := &papi.Mock{}
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
					ProtoV6ProviderFactories: testutils.NewProtoV6ProviderFactory(NewSubprovider()),
					Steps: []resource.TestStep{
						{
							Config:      testutils.LoadFixtureString(t, "testdata/%s.tf", t.Name()),
							Check:       resource.TestCheckNoResourceAttr("akamai_property.test", "id"),
							ExpectError: regexp.MustCompile(`property name is not unique`),
						},
					},
				})
			})

			client.AssertExpectations(t)
		})

		ruleTreeRes := papi.GetRuleTreeResponse{
			Rules: papi.Rules{
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
		}

		propertyReadCtx := func(client *papi.Mock, stagStatus, prodStatus papi.VersionStatus) {
			ExpectGetProperty(
				client, "prp_0", "grp_0", "ctr_0",
				&papi.Property{
					PropertyID: "prp_0", GroupID: "grp_0", ContractID: "ctr_0", LatestVersion: 1,
					PropertyName: "dxe-2406-issue-example",
				},
			).Once()
			ExpectGetPropertyVersionHostnames(
				client, "prp_0", "grp_0", "ctr_0", 1,
				&[]papi.Hostname{
					{
						CnameFrom:            "dxe-2406-issue-example-second.com",
						CnameTo:              "dxe-2406-issue-example-second.com.example.net",
						CertProvisioningType: "CPS_MANAGED",
					},
					{
						CnameFrom:            "dxe-2406-issue.com",
						CnameTo:              "dxe-2406-issue.com.example.net",
						CertProvisioningType: "CPS_MANAGED",
					},
				},
			).Once()
			client.On("GetRuleTree", mock.Anything, papi.GetRuleTreeRequest{
				PropertyID:      "prp_0",
				GroupID:         "grp_0",
				ContractID:      "ctr_0",
				PropertyVersion: 1,
				ValidateMode:    "full",
				ValidateRules:   true,
			}).Return(&ruleTreeRes, nil).Once()
			ExpectGetPropertyVersion(client, "prp_0", "grp_0", "ctr_0", 1, stagStatus, prodStatus).Once()
		}

		getActivations := func(client *papi.Mock) {
			expectGetActivations(client, "prp_0", papi.GetActivationsResponse{
				Activations: papi.ActivationsItems{
					Items: []*papi.Activation{
						{
							ActivationID:    "act_123",
							PropertyID:      "prp_0",
							PropertyVersion: 1,
							Network:         papi.ActivationNetworkStaging,
							Status:          papi.ActivationStatusActive,
							ActivationType:  papi.ActivationTypeActivate,
							SubmitDate:      "2020-10-28T15:04:05Z",
							NotifyEmails:    []string{"dummy-user@akamai.com"},
						},
					},
				},
			}, nil).Once()
		}

		t.Run("400 from UpdatePropertyVersionHostnames - incorrect/invalid edge hostname", func(t *testing.T) {
			client := &papi.Mock{}
			client.Test(T{t})
			ruleFormat := ""

			// first step
			// create property
			ExpectCreateProperty(client, "dxe-2406-issue-example", "grp_0", "ctr_0", "prd_0", "prp_0").Once()
			ExpectUpdatePropertyVersionHostnames(
				client, "prp_0", "grp_0", "ctr_0", 1,
				[]papi.Hostname{
					{
						CnameType:            "EDGE_HOSTNAME",
						CnameFrom:            "dxe-2406-issue-example-second.com",
						CnameTo:              "dxe-2406-issue-example-second.com.example.net",
						CertProvisioningType: "CPS_MANAGED",
					},
					{
						CnameType:            "EDGE_HOSTNAME",
						CnameFrom:            "dxe-2406-issue.com",
						CnameTo:              "dxe-2406-issue.com.example.net",
						CertProvisioningType: "CPS_MANAGED",
					}}, nil,
			).Once()
			ExpectUpdateRuleTree(client, "prp_0", "grp_0", "ctr_0", 1,
				&papi.RulesUpdate{
					Rules: papi.Rules{
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
					}}, ruleFormat, []papi.RuleError{}).Once()

			// read property
			propertyReadCtx(client, papi.VersionStatusInactive, papi.VersionStatusInactive)

			// create activation
			expectGetRuleTree(client, "prp_0", 1, ruleTreeRes, nil).Once()
			expectGetActivations(client, "prp_0", papi.GetActivationsResponse{}, nil).Once()
			client.On("CreateActivation", mock.Anything, mock.Anything).Return(&papi.CreateActivationResponse{ActivationID: "act_123"}, nil).Once()
			expectGetActivation(client, "prp_0", "act_123", 1, papi.ActivationNetworkStaging, papi.ActivationStatusActive, papi.ActivationTypeActivate, "", []string{"dummy-user@akamai.com"}, nil).Once()

			// read property
			propertyReadCtx(client, papi.VersionStatusActive, papi.VersionStatusActive)

			// activation read
			getActivations(client)

			// read property
			propertyReadCtx(client, papi.VersionStatusActive, papi.VersionStatusActive)

			// activation read
			getActivations(client)

			// second step
			// property update returns an error on the invalid edgehostname
			ExpectGetPropertyVersion(client, "prp_0", "grp_0", "ctr_0", 1, papi.VersionStatusActive, papi.VersionStatusActive).Once()
			ExpectCreatePropertyVersion(client, "prp_0", "grp_0", "ctr_0", 1, 2)

			ExpectUpdatePropertyVersionHostnames(
				client, "prp_0", "grp_0", "ctr_0", 2,
				[]papi.Hostname{
					{
						CnameType:            "EDGE_HOSTNAME",
						CnameFrom:            "dxe-2406-issue-example-second.com",
						CnameTo:              "dxe-2406-issue-example-second.com.example.net",
						CertProvisioningType: "CPS_MANAGED",
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
					}}, fmt.Errorf("%w: request failed: %s", papi.ErrUpdatePropertyVersionHostnames, errors.New("{\n    \"type\": \"https://problems.luna.akamaiapis.net/papi/v0/property-version-hostname/bad-cnameto\",\n    \"title\": \"Bad `cnameTo`\",\n    \"detail\": \"The System could not find cnameTo value `does-not-exist.com.example.net`.\",\n    \"instance\": \"host/papi/v1/properties/prp_0/versions/2/hostnames?contractId=ctr_0&groupId=grp_0&includeCertStatus=false&validateHostnames=false#efba6490291100b1\",\n    \"status\": 400\n}")),
			).Once()

			// terraform clean up - terraform test framework attempts to run destroy plan, if an error is returned on second step
			// activation and property deletion
			getActivations(client)
			client.On("CreateActivation", mock.Anything, mock.Anything).Return(&papi.CreateActivationResponse{
				ActivationID: "act_123",
			}, nil).Once()
			expectGetActivation(client, "prp_0", "act_123", 1, papi.ActivationNetworkStaging, papi.ActivationStatusActive, papi.ActivationTypeDeactivate, "", []string{"dummy-user@akamai.com"}, nil).Once()
			client.On("RemoveProperty", mock.Anything, mock.Anything).Return(&papi.RemovePropertyResponse{
				Message: "removed",
			}, nil).Once()

			useClient(client, nil, func() {
				resource.UnitTest(t, resource.TestCase{
					ProtoV6ProviderFactories: testutils.NewProtoV6ProviderFactory(NewSubprovider()),
					Steps: []resource.TestStep{
						{
							Config: testutils.LoadFixtureString(t, "testdata/TestResProperty/CreationUpdateIncorrectEdgeHostname/create/property.tf"),
							Check: resource.ComposeAggregateTestCheckFunc(
								resource.TestCheckResourceAttr("akamai_property.akaproperty", "id", "prp_0"),
								resource.TestCheckResourceAttr("akamai_property.akaproperty", "hostnames.#", "2"),
							),
						},
						{
							Config: testutils.LoadFixtureString(t, "testdata/TestResProperty/CreationUpdateIncorrectEdgeHostname/update/property.tf"),
							Check: resource.ComposeAggregateTestCheckFunc(
								resource.TestCheckResourceAttr("akamai_property.akaproperty", "id", "prp_0"),
								resource.TestCheckResourceAttr("akamai_property.akaproperty", "hostnames.#", "3"),
							),
							ExpectError: regexp.MustCompile("Error: updating hostnames: request failed:"),
						},
					},
				})
			})

			client.AssertExpectations(t)
		})
	})
}

func TestPropertyResource_versionNotesLifecycle(t *testing.T) {
	testdataDir := "testdata/TestResProperty/Lifecycle/versionNotes"
	resourceName := "akamai_property.test"

	name := "test_property"
	ruleFormat := "v2023-01-05"
	ctr, grp, prd, id := "ctr_123", "grp_123", "prd_123", "prp_123"
	propertyVersion := 1

	versionNotes1, versionNotes2, versionNotes3 := "lifecycleTest", "updatedNotes", "updatedNotes2"
	rulesFile1And2, rulesFile3, rulesFile4And5 := "01_02_rules.json", "03_rules.json", "04_05_rules.json"

	client := &papi.Mock{}

	mockRead := func(notes string, rules papi.Rules) testutils.MockCalls {
		getPropertyCall := client.On("GetProperty", mock.Anything, papi.GetPropertyRequest{
			PropertyID: id,
			ContractID: ctr,
			GroupID:    grp,
		}).Return(&papi.GetPropertyResponse{
			Property: &papi.Property{
				ContractID:    ctr,
				GroupID:       grp,
				PropertyID:    id,
				PropertyName:  name,
				LatestVersion: propertyVersion,
			},
		}, nil)

		getHostnamesCall := client.On("GetPropertyVersionHostnames", mock.Anything, papi.GetPropertyVersionHostnamesRequest{
			ContractID:        ctr,
			GroupID:           grp,
			PropertyID:        id,
			PropertyVersion:   propertyVersion,
			IncludeCertStatus: true,
		}).Return(&papi.GetPropertyVersionHostnamesResponse{}, nil)

		getRuleTreeCall := client.On("GetRuleTree", mock.Anything, papi.GetRuleTreeRequest{
			PropertyID:      id,
			ContractID:      ctr,
			GroupID:         grp,
			PropertyVersion: propertyVersion,
			ValidateRules:   true,
			ValidateMode:    papi.RuleValidateModeFull,
		}).Return(&papi.GetRuleTreeResponse{
			Rules:      rules,
			Comments:   notes,
			RuleFormat: ruleFormat,
		}, nil)

		getPropertyVersionCall := client.On("GetPropertyVersion", mock.Anything, papi.GetPropertyVersionRequest{
			PropertyID:      id,
			PropertyVersion: propertyVersion,
			ContractID:      ctr,
			GroupID:         grp,
		}).Return(&papi.GetPropertyVersionsResponse{
			Version: papi.PropertyVersionGetItem{
				Note:             notes,
				ProductID:        prd,
				ProductionStatus: papi.VersionStatusInactive,
				StagingStatus:    papi.VersionStatusInactive,
			},
		}, nil)

		return testutils.MockCalls{getPropertyCall, getHostnamesCall, getRuleTreeCall, getPropertyVersionCall}
	}

	mockUpdate := func(currentNotes, newNotes string, rules papi.Rules) {
		client.On("GetPropertyVersion", mock.Anything, papi.GetPropertyVersionRequest{
			PropertyID:      id,
			PropertyVersion: propertyVersion,
			ContractID:      ctr,
			GroupID:         grp,
		}).Return(&papi.GetPropertyVersionsResponse{
			Version: papi.PropertyVersionGetItem{
				Note:             currentNotes,
				ProductID:        prd,
				ProductionStatus: papi.VersionStatusInactive,
				StagingStatus:    papi.VersionStatusInactive,
			},
		}, nil).Once()

		client.On("UpdateRuleTree", mock.Anything, papi.UpdateRulesRequest{
			PropertyID:      id,
			GroupID:         grp,
			ContractID:      ctr,
			PropertyVersion: propertyVersion,
			Rules: papi.RulesUpdate{
				Rules:    rules,
				Comments: newNotes,
			},
			ValidateRules: true,
		}).Return(&papi.UpdateRulesResponse{}, nil).Once()
	}

	// step 1 - create + read + plan
	client.On("CreateProperty", mock.Anything, papi.CreatePropertyRequest{
		ContractID: ctr,
		GroupID:    grp,
		Property: papi.PropertyCreate{
			ProductID:    prd,
			PropertyName: name,
			RuleFormat:   ruleFormat,
		},
	}).Return(&papi.CreatePropertyResponse{
		PropertyID: id,
	}, nil).Once()

	rulesJSON := testutils.LoadFixtureBytes(t, path.Join(testdataDir, rulesFile1And2))
	var rules1And2 papi.RulesUpdate
	err := json.Unmarshal(rulesJSON, &rules1And2)
	require.NoError(t, err)

	client.On("UpdateRuleTree", mock.Anything, papi.UpdateRulesRequest{
		PropertyID:      id,
		GroupID:         grp,
		ContractID:      ctr,
		PropertyVersion: propertyVersion,
		Rules: papi.RulesUpdate{
			Rules:    rules1And2.Rules,
			Comments: versionNotes1,
		},
		ValidateRules: true,
	}).Return(&papi.UpdateRulesResponse{}, nil).Once()

	mockRead(versionNotes1, rules1And2.Rules).Times(2)

	// step 2 - refresh + plan
	mockRead(versionNotes2, rules1And2.Rules).Times(2)

	// step 3 - refresh + update + read + plan
	mockRead(versionNotes2, rules1And2.Rules).Times(1)

	var rules3 papi.RulesUpdate
	rulesJSON = testutils.LoadFixtureBytes(t, path.Join(testdataDir, rulesFile3))
	err = json.Unmarshal(rulesJSON, &rules3)
	require.NoError(t, err)

	mockUpdate(versionNotes3, "updatedNotes2", rules3.Rules)

	mockRead(versionNotes3, rules3.Rules).Times(2)

	// step 4 - refresh + update + read + plan
	mockRead(versionNotes3, rules3.Rules).Times(1)

	var rules4And5 papi.RulesUpdate
	rulesJSON = testutils.LoadFixtureBytes(t, path.Join(testdataDir, rulesFile4And5))
	err = json.Unmarshal(rulesJSON, &rules4And5)
	require.NoError(t, err)

	mockUpdate(versionNotes3, rules4And5.Comments, rules4And5.Rules)

	mockRead(rules4And5.Comments, rules4And5.Rules).Times(2)

	// step 5 - refresh + plan
	mockRead(rules4And5.Comments, rules4And5.Rules).Times(2)

	// cleanup
	client.On("RemoveProperty", mock.Anything, papi.RemovePropertyRequest{
		PropertyID: id,
		ContractID: ctr,
		GroupID:    grp,
	}).Return(&papi.RemovePropertyResponse{}, nil)

	useClient(client, nil, func() {
		resource.UnitTest(t, resource.TestCase{
			ProtoV6ProviderFactories: testutils.NewProtoV6ProviderFactory(NewSubprovider()),
			Steps: []resource.TestStep{
				{
					Config: testutils.LoadFixtureString(t, path.Join(testdataDir, "01_with_notes_and_comments.tf")),
					Check: resource.ComposeAggregateTestCheckFunc(
						testCheckResourceAttrJSON(resourceName, "rules", testutils.LoadFixtureString(t, path.Join(testdataDir, "01_expected_rules.json"))),
						resource.TestCheckResourceAttr("akamai_property.test", "version_notes", "lifecycleTest"),
					),
				},
				{
					Config:   testutils.LoadFixtureString(t, path.Join(testdataDir, "02_update_notes_no_diff.tf")),
					PlanOnly: true,
				},
				{
					Config: testutils.LoadFixtureString(t, path.Join(testdataDir, "03_update_notes_and_rules.tf")),
					Check: resource.ComposeAggregateTestCheckFunc(
						testCheckResourceAttrJSON(resourceName, "rules", testutils.LoadFixtureString(t, path.Join(testdataDir, "03_expected_rules.json"))),
						resource.TestCheckResourceAttr("akamai_property.test", "version_notes", "updatedNotes2"),
					),
				},
				{
					Config: testutils.LoadFixtureString(t, path.Join(testdataDir, "04_05_remove_notes_update_comments.tf")),
					Check: resource.ComposeAggregateTestCheckFunc(
						testCheckResourceAttrJSON(resourceName, "rules", testutils.LoadFixtureString(t, path.Join(testdataDir, "04_expected_rules.json"))),
						resource.TestCheckResourceAttr("akamai_property.test", "version_notes", "Rules_04"),
					),
				},
				{
					Config:   testutils.LoadFixtureString(t, path.Join(testdataDir, "04_05_remove_notes_update_comments.tf")),
					PlanOnly: true,
				},
			},
		})
	})

	client.AssertExpectations(t)
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
