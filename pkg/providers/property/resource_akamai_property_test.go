package property

import (
	"fmt"
	"regexp"
	"testing"

	"github.com/akamai/AkamaiOPEN-edgegrid-golang/v2/pkg/papi"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/resource"
	"github.com/stretchr/testify/mock"
)

func TestResProperty(t *testing.T) {
	// These more or less track the state of a Property in PAPI
	type TestState struct {
		Client     *mockpapi
		Property   papi.Property
		Hostnames  []papi.Hostname
		Rules      papi.Rules
		RuleFormat string
		Groups     []papi.Group
		err        error
	}

	// BehaviorFuncs can be composed to define common patterns of mock PAPI behavior
	type BehaviorFunc = func(*TestState)

	// Combines many BehaviorFuncs into one
	ComposeBehaviors := func(behaviors ...BehaviorFunc) BehaviorFunc {
		return func(State *TestState) {
			for _, behave := range behaviors {
				behave(State)
			}
		}
	}

	SetHostnames := func(PropertyID string, Version int, CnameTo string) BehaviorFunc {
		return func(State *TestState) {
			NewHostnames := []papi.Hostname{{
				CnameType: "EDGE_HOSTNAME",
				CnameFrom: "from.test.domain",
				CnameTo:   CnameTo,
			}}

			ExpectUpdatePropertyVersionHostnames(State.Client, PropertyID, "grp_0", "ctr_0", Version, NewHostnames).Once().Run(func(mock.Arguments) {
				State.Hostnames = append([]papi.Hostname{}, NewHostnames...)
			})
		}
	}

	GetVersionResources := func(PropertyID string, Version int) BehaviorFunc {
		return func(State *TestState) {
			ExpectGetPropertyVersionHostnames(State.Client, PropertyID, "grp_0", "ctr_0", Version, &State.Hostnames)
			ExpectGetRuleTree(State.Client, PropertyID, "grp_0", "ctr_0", Version, &State.Rules, &State.RuleFormat)
		}
	}

	DeleteProperty := func(PropertyID string) BehaviorFunc {
		return func(State *TestState) {
			ExpectRemoveProperty(State.Client, PropertyID, "ctr_0", "grp_0").Once().Run(func(mock.Arguments) {
				State.Property = papi.Property{}
				State.Rules = papi.Rules{}
				State.Hostnames = nil
				State.RuleFormat = ""
			})
		}
	}

	GetProperty := func(PropertyID string) BehaviorFunc {
		return func(State *TestState) {
			ExpectGetProperty(State.Client, PropertyID, "grp_0", "ctr_0", &State.Property, State.err)
		}
	}

	CreateProperty := func(PropertyName, PropertyID string, err error) BehaviorFunc {
		return func(State *TestState) {
			ExpectCreateProperty(State.Client, PropertyName, "grp_0", "ctr_0", "prd_0", PropertyID, err).Run(func(mock.Arguments) {
				if err != nil {
					State.err = err
				}

				State.Property = papi.Property{
					PropertyName:  PropertyName,
					PropertyID:    PropertyID,
					GroupID:       "grp_0",
					ContractID:    "ctr_0",
					ProductID:     "prd_0",
					LatestVersion: 1,
				}

				State.Rules = papi.Rules{Name: "default"}
				State.RuleFormat = "v2020-01-01"
			}).Once()

			if err == nil {
				GetProperty(PropertyID)(State)
				GetVersionResources(PropertyID, 1)(State)
			}
		}
	}

	GetGroups := func(groups []*papi.Group, err error) BehaviorFunc {
		return func(State *TestState) {
			ExpectGetGroups(State.Client, err, groups).Once()
		}
	}

	PropertyLifecycle := func(PropertyName, PropertyID, GroupID string) BehaviorFunc {
		return func(State *TestState) {
			CreateProperty(PropertyName, PropertyID, nil)(State)
			GetVersionResources(PropertyID, 1)(State)
			DeleteProperty(PropertyID)(State)
		}
	}

	ImportProperty := func(PropertyID string) BehaviorFunc {
		return func(State *TestState) {
			// Depending on how much of the import ID is given, the initial property lookup may not have group/contract
			ExpectGetProperty(State.Client, "prp_0", "grp_0", "", &State.Property, nil).Maybe()
			ExpectGetProperty(State.Client, "prp_0", "", "", &State.Property, nil).Maybe()
		}
	}

	AdvanceVersion := func(PropertyID string, FromVersion, ToVersion int) BehaviorFunc {
		return func(State *TestState) {
			ExpectCreatePropertyVersion(State.Client, PropertyID, "grp_0", "ctr_0", FromVersion, ToVersion).Once().Run(func(mock.Arguments) {
				State.Property.LatestVersion = ToVersion
			})
			GetVersionResources(PropertyID, ToVersion)(State)
		}
	}

	// TestCheckFunc to verify all standard attributes
	CheckAttrs := func(PropertyID, CnameTo, LatestVersion, StagingVersion, ProductionVersion string) resource.TestCheckFunc {
		return resource.ComposeAggregateTestCheckFunc(
			resource.TestCheckResourceAttr("akamai_property.test", "id", PropertyID),
			resource.TestCheckResourceAttr("akamai_property.test", "hostnames.from.test.domain", CnameTo),
			resource.TestCheckResourceAttr("akamai_property.test", "latest_version", LatestVersion),
			resource.TestCheckResourceAttr("akamai_property.test", "staging_version", StagingVersion),
			resource.TestCheckResourceAttr("akamai_property.test", "production_version", ProductionVersion),
			resource.TestCheckResourceAttr("akamai_property.test", "name", "test property"),
			resource.TestCheckResourceAttr("akamai_property.test", "contract_id", "ctr_0"),
			resource.TestCheckResourceAttr("akamai_property.test", "contract", "ctr_0"),
			resource.TestCheckResourceAttr("akamai_property.test", "group_id", "grp_0"),
			resource.TestCheckResourceAttr("akamai_property.test", "group", "grp_0"),
			resource.TestCheckResourceAttr("akamai_property.test", "product", "prd_0"),
			resource.TestCheckResourceAttr("akamai_property.test", "product_id", "prd_0"),
			resource.TestCheckResourceAttr("akamai_property.test", "rules", `{"name":"default","options":{}}`),
		)
	}

	type StepsFunc = func(State *TestState, FixturePath string) []resource.TestStep

	// Defines standard variations of client behaviors for a Lifecycle test
	type LifecycleTestCase struct {
		Name        string
		ClientSetup BehaviorFunc
		Steps       StepsFunc
	}

	// Standard test behavior for cases where the property's latest version is active in staging network
	LatestVersionActiveInStaging := LifecycleTestCase{
		Name: "Latest version is active in staging",
		ClientSetup: ComposeBehaviors(
			PropertyLifecycle("test property", "prp_0", "grp_0"),
			SetHostnames("prp_0", 1, "to.test.domain"),
			AdvanceVersion("prp_0", 1, 2),
			SetHostnames("prp_0", 2, "to2.test.domain"),
		),
		Steps: func(State *TestState, FixturePath string) []resource.TestStep {
			return []resource.TestStep{
				{
					Config: loadFixtureString("%s/step0.tf", FixturePath),
					Check:  CheckAttrs("prp_0", "to.test.domain", "1", "0", "0"),
				},
				{
					PreConfig: func() {
						StagingVersion := 1
						State.Property.StagingVersion = &StagingVersion
					},
					Config: loadFixtureString("%s/step1.tf", FixturePath),
					Check:  CheckAttrs("prp_0", "to2.test.domain", "2", "1", "0"),
				},
			}
		},
	}

	// Standard test behavior for cases where the property's latest version is active in production network
	LatestVersionActiveInProd := LifecycleTestCase{
		Name: "Latest version is active in production",
		ClientSetup: ComposeBehaviors(
			PropertyLifecycle("test property", "prp_0", "grp_0"),
			SetHostnames("prp_0", 1, "to.test.domain"),
			AdvanceVersion("prp_0", 1, 2),
			SetHostnames("prp_0", 2, "to2.test.domain"),
		),
		Steps: func(State *TestState, FixturePath string) []resource.TestStep {
			return []resource.TestStep{
				{
					Config: loadFixtureString("%s/step0.tf", FixturePath),
					Check:  CheckAttrs("prp_0", "to.test.domain", "1", "0", "0"),
				},
				{
					PreConfig: func() {
						ProductionVersion := 1
						State.Property.ProductionVersion = &ProductionVersion
					},
					Config: loadFixtureString("%s/step1.tf", FixturePath),
					Check:  CheckAttrs("prp_0", "to2.test.domain", "2", "0", "1"),
				},
			}
		},
	}

	// Standard test behavior for cases where the property's latest version is not active
	LatestVersionNotActive := LifecycleTestCase{
		Name: "Latest version not active",
		ClientSetup: ComposeBehaviors(
			PropertyLifecycle("test property", "prp_0", "grp_0"),
			SetHostnames("prp_0", 1, "to.test.domain"),
			SetHostnames("prp_0", 1, "to2.test.domain"),
		),
		Steps: func(State *TestState, FixturePath string) []resource.TestStep {
			return []resource.TestStep{
				{
					Config: loadFixtureString("%s/step0.tf", FixturePath),
					Check:  CheckAttrs("prp_0", "to.test.domain", "1", "0", "0"),
				},
				{
					Config: loadFixtureString("%s/step1.tf", FixturePath),
					Check:  CheckAttrs("prp_0", "to2.test.domain", "1", "0", "0"),
				},
			}
		},
	}

	// Run a test case to verify schema validations
	AssertConfigError := func(t *testing.T, flaw, rx string) {
		t.Helper()
		caseName := fmt.Sprintf("ConfigError/%s", flaw)

		t.Run(caseName, func(t *testing.T) {
			t.Helper()

			resource.UnitTest(t, resource.TestCase{
				Providers: testAccProviders,
				Steps: []resource.TestStep{{
					Config:      loadFixtureString("testdata/%s.tf", t.Name()),
					ExpectError: regexp.MustCompile(rx),
				}},
			})
		})
	}

	// Run a test case to verify schema attribute deprecation
	AssertDeprecated := func(t *testing.T, attribute string) {
		t.Helper()

		t.Run(fmt.Sprintf("%s attribute is deprecated", attribute), func(t *testing.T) {
			t.Helper()
			if resourceProperty().Schema[attribute].Deprecated == "" {
				t.Fatalf(`%q attribute is not marked deprecated`, attribute)
			}
		})
	}

	// Run a test case to confirm that the user is prompted to read the upgrade guide
	AssertForbiddenAttr := func(t *testing.T, fixtureName string) {
		t.Helper()

		t.Run(fmt.Sprintf("ForbiddenAttr/%s", fixtureName), func(t *testing.T) {
			t.Helper()
			client := &mockpapi{}
			client.Test(T{t})

			useClient(client, func() {
				resource.UnitTest(t, resource.TestCase{
					Providers: testAccProviders,
					Steps: []resource.TestStep{{
						Config:      loadFixtureString("testdata/%s.tf", t.Name()),
						ExpectError: regexp.MustCompile("See the Akamai Terraform Upgrade Guide"),
					}},
				})
			})

			client.AssertExpectations(t)
		})
	}

	// Run a happy-path test case that goes through a complete create-update-destroy cycle
	AssertLifecycle := func(t *testing.T, variant string, kase LifecycleTestCase) {
		t.Helper()

		fixturePrefix := fmt.Sprintf("testdata/%s/Lifecycle/%s", t.Name(), variant)
		testName := fmt.Sprintf("Lifecycle/%s/%s", variant, kase.Name)

		t.Run(testName, func(t *testing.T) {
			t.Helper()

			client := &mockpapi{}
			client.Test(T{t})
			State := &TestState{Client: client}
			kase.ClientSetup(State)

			useClient(client, func() {
				resource.UnitTest(t, resource.TestCase{
					Providers:    testAccProviders,
					Steps:        kase.Steps(State, fixturePrefix),
					CheckDestroy: resource.TestCheckNoResourceAttr("akamai_property.test", "id"),
				})
			})

			client.AssertExpectations(t)
		})
	}

	// Run a test case that verifies the resource can be imported by the given ID
	AssertImportable := func(t *testing.T, TestName, ImportID string) {
		t.Helper()

		fixturePath := fmt.Sprintf("testdata/%s/Importable/importable.tf", t.Name())
		testName := fmt.Sprintf("Importable/%s", TestName)

		t.Run(testName, func(t *testing.T) {
			t.Helper()

			client := &mockpapi{}
			client.Test(T{t})

			setup := ComposeBehaviors(
				PropertyLifecycle("test property", "prp_0", "grp_0"),
				ImportProperty("prp_0"),
				SetHostnames("prp_0", 1, "to.test.domain"),
			)
			setup(&TestState{Client: client})

			useClient(client, func() {
				resource.UnitTest(t, resource.TestCase{
					Providers: testAccProviders,
					Steps: []resource.TestStep{
						{
							Config: loadFixtureString(fixturePath),
							Check:  CheckAttrs("prp_0", "to.test.domain", "1", "0", "0"),
						},
						{
							ImportState:       true,
							ImportStateVerify: true,
							ImportStateId:     ImportID,
							ResourceName:      "akamai_property.test",
							Config:            loadFixtureString(fixturePath),
						},
						{
							Config: loadFixtureString(fixturePath),
							Check:  CheckAttrs("prp_0", "to.test.domain", "1", "0", "0"),
						},
					},
				})
			})

			client.AssertExpectations(t)
		})
	}

	suppressLogging(t, func() {
		AssertConfigError(t, "name not given", `"name" is required`)
		AssertConfigError(t, "neither contract nor contract_id given", `one of .contract,contract_id. must be specified`)
		AssertConfigError(t, "both contract and contract_id given", `only one of .contract,contract_id. can be specified`)
		AssertConfigError(t, "neither group nor group_id given", `one of .group,group_id. must be specified`)
		AssertConfigError(t, "both group and group_id given", `only one of .group,group_id. can be specified`)
		AssertConfigError(t, "neither product nor product_id given", `one of .product,product_id. must be specified`)
		AssertConfigError(t, "both product and product_id given", `only one of .product,product_id. can be specified`)
		AssertConfigError(t, "invalid json rules", `rules are not valid JSON`)

		AssertDeprecated(t, "contract")
		AssertDeprecated(t, "group")
		AssertDeprecated(t, "product")
		AssertDeprecated(t, "cp_code")
		AssertDeprecated(t, "contact")
		AssertDeprecated(t, "origin")
		AssertDeprecated(t, "is_secure")
		AssertDeprecated(t, "variables")

		AssertForbiddenAttr(t, "cp_code")
		AssertForbiddenAttr(t, "contact")
		AssertForbiddenAttr(t, "origin")
		AssertForbiddenAttr(t, "is_secure")
		AssertForbiddenAttr(t, "variables")

		AssertLifecycle(t, "normal", LatestVersionNotActive)
		AssertLifecycle(t, "normal", LatestVersionActiveInStaging)
		AssertLifecycle(t, "normal", LatestVersionActiveInProd)
		AssertLifecycle(t, "contract_id without prefix", LatestVersionNotActive)
		AssertLifecycle(t, "contract_id without prefix", LatestVersionActiveInStaging)
		AssertLifecycle(t, "contract_id without prefix", LatestVersionActiveInProd)
		AssertLifecycle(t, "contract without prefix", LatestVersionNotActive)
		AssertLifecycle(t, "contract without prefix", LatestVersionActiveInStaging)
		AssertLifecycle(t, "contract without prefix", LatestVersionActiveInProd)
		AssertLifecycle(t, "group_id without prefix", LatestVersionNotActive)
		AssertLifecycle(t, "group_id without prefix", LatestVersionActiveInStaging)
		AssertLifecycle(t, "group_id without prefix", LatestVersionActiveInProd)
		AssertLifecycle(t, "group without prefix", LatestVersionNotActive)
		AssertLifecycle(t, "group without prefix", LatestVersionActiveInStaging)
		AssertLifecycle(t, "group without prefix", LatestVersionActiveInProd)
		AssertLifecycle(t, "product_id without prefix", LatestVersionNotActive)
		AssertLifecycle(t, "product_id without prefix", LatestVersionActiveInStaging)
		AssertLifecycle(t, "product_id without prefix", LatestVersionActiveInProd)
		AssertLifecycle(t, "product without prefix", LatestVersionNotActive)
		AssertLifecycle(t, "product without prefix", LatestVersionActiveInStaging)
		AssertLifecycle(t, "product without prefix", LatestVersionActiveInProd)

		AssertImportable(t, "property_id", "prp_0")
		AssertImportable(t, "unprefixed property_id", "0")
		AssertImportable(t, "property_id and group_id", "prp_0,grp_0")
		AssertImportable(t, "unprefixed property_id and group_id", "0,0")
		AssertImportable(t, "property_id and group_id and contract_id", "prp_0,grp_0,ctr_0")
		AssertImportable(t, "unprefixed property_id and group_id and contract_id", "0,0,0")

		t.Run("property is destroyed and recreated when name is changed", func(t *testing.T) {
			client := &mockpapi{}
			client.Test(T{t})

			setup := ComposeBehaviors(
				PropertyLifecycle("test property", "prp_0", "grp_0"),
				PropertyLifecycle("renamed property", "prp_1", "grp_0"),
				SetHostnames("prp_0", 1, "to.test.domain"),
				SetHostnames("prp_1", 1, "to2.test.domain"),
			)
			setup(&TestState{Client: client})

			useClient(client, func() {
				resource.UnitTest(t, resource.TestCase{
					Providers: testAccProviders,
					Steps: []resource.TestStep{
						{
							Config: loadFixtureString("testdata/%s-step0.tf", t.Name()),
							Check:  CheckAttrs("prp_0", "to.test.domain", "1", "0", "0"),
						},
						{
							Config: loadFixtureString("testdata/%s-step1.tf", t.Name()),
							Check: resource.ComposeAggregateTestCheckFunc(
								resource.TestCheckResourceAttr("akamai_property.test", "id", "prp_1"),
								resource.TestCheckResourceAttr("akamai_property.test", "name", "renamed property"),
							),
						},
					},
					CheckDestroy: resource.TestCheckNoResourceAttr("akamai_property.test", "id"),
				})
			})

			client.AssertExpectations(t)
		})

		t.Run("error when deleting active property", func(t *testing.T) {
			client := &mockpapi{}
			client.Test(T{t})

			setup := ComposeBehaviors(
				CreateProperty("test property", "prp_0", nil),
				GetProperty("prp_0"),
				GetVersionResources("prp_0", 1),
				SetHostnames("prp_0", 1, "to.test.domain"),
			)
			setup(&TestState{Client: client})

			// First call to remove is not successful
			req := papi.RemovePropertyRequest{
				PropertyID: "prp_0",
				ContractID: "ctr_0",
				GroupID:    "grp_0",
			}

			err := fmt.Errorf(`Cannot remove active property "prp_0"`)
			client.On("RemoveProperty", AnyCTX, req).Return(nil, err).Once()

			// Second call will be successful (TF test case requires last state to be empty or it's a failed test)
			ExpectRemoveProperty(client, "prp_0", "ctr_0", "grp_0").Once()

			useClient(client, func() {
				resource.UnitTest(t, resource.TestCase{
					Providers: testAccProviders,
					Steps: []resource.TestStep{
						{
							Config: loadFixtureString("testdata/%s/step0.tf", t.Name()),
							Check:  CheckAttrs("prp_0", "to.test.domain", "1", "0", "0"),
						},
						{
							Config:      loadFixtureString("testdata/%s/step1.tf", t.Name()),
							ExpectError: regexp.MustCompile(`Cannot remove active property`),
						},
					},
				})
			})

			client.AssertExpectations(t)
		})

		t.Run("error when the given group is not found", func(t *testing.T) {
			client := &mockpapi{}
			client.Test(T{t})

			// the papi GetGroups call should not return any matching group
			setup := ComposeBehaviors(
				CreateProperty("property_name", "prp_0", &papi.Error{
					StatusCode: 404,
					Title:      "Not Found",
					Detail:     "The system was unable to locate the requested resource",
					Type:       "https://problems.luna.akamaiapis.net/papi/v0/http/not-found",
					Instance:   "https://akaa-hqgqowhpmkw32kmt-t3owzo37wb5dkern.luna-dev.akamaiapis.net/papi/v1/properties?contractId=ctr_0\\u0026groupId=grp_0#c3fe5f9b0c4a14d1",
				},
				),
				GetGroups([]*papi.Group{}, nil),
			)
			setup(&TestState{Client: client})

			useClient(client, func() {
				resource.UnitTest(t, resource.TestCase{
					Providers: testAccProviders,
					Steps: []resource.TestStep{
						{
							Config:      loadFixtureString("testdata/TestResProperty/Creation/property.tf"),
							ExpectError: regexp.MustCompile("group not found: grp_0"),
						},
					},
				})
			})

			client.AssertExpectations(t)
		})
	})
}
