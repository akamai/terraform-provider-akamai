package property

import (
	"fmt"
	"regexp"
	"testing"

	"github.com/akamai/AkamaiOPEN-edgegrid-golang/v2/pkg/papi"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/resource"
	"github.com/stretchr/testify/mock"
)

// Wrapper to intercept the mockpapi's call of t.FailNow(). The Terraform test driver runs the provider code on
// goroutines other than the one created for the test. When t.FailNow() is called from any other goroutine, it causes
// the test to hang because the TF test driver is still waiting to serve requests. Mockery's failure message neglects to
// inform the user which test had failed. Use this struct to wrap a *testing.T when you call mock.Test(T{t}) and the
// mock's failure will print the failling test's name. Such failures are usually caused by the provider invoking an
// unexpected call on the mock.
//
// NB: You should only need to use this where your test uses the Terraform test driver
type T struct{ *testing.T }

// Overrides testing.T.FailNow() so when a test mock fails an assertion, we see which test had failed before it hangs
func (t T) FailNow() {
	t.T.Fatalf("FAIL: %s", t.T.Name())
}

func TestResProperty(t *testing.T) {
	// Helper to in-line expected call to papi.GetProperty() with a constant Property struct with values that are
	// standard for this test
	ExpectGetProp := func(client *mockpapi, PropertyName, PropertyID, GroupID, ContractID, ProductID string, Version int) *mock.Call {
		Property := papi.Property{
			PropertyName:   PropertyName,
			PropertyID:     PropertyID,
			GroupID:        GroupID,
			ContractID:     ContractID,
			ProductID:      "prd_0",
			LatestVersion:  Version,
			StagingVersion: &Version,
		}

		// return ExpectGetProperty(client, PropertyID, "", "", &Property)
		return ExpectGetProperty(client, PropertyID, GroupID, ContractID, &Property)
	}

	// TestCheckFunc to verify all standard attributes
	checkAllAttrs := resource.ComposeAggregateTestCheckFunc(
		resource.TestCheckResourceAttr("akamai_property.test", "id", "prp_0"),
		resource.TestCheckResourceAttr("akamai_property.test", "production_version", "0"),
		resource.TestCheckResourceAttr("akamai_property.test", "staging_version", "42"),
		resource.TestCheckResourceAttr("akamai_property.test", "name", "test property"),
		resource.TestCheckResourceAttr("akamai_property.test", "contract_id", "ctr_0"),
		resource.TestCheckResourceAttr("akamai_property.test", "contract", "ctr_0"),
		resource.TestCheckResourceAttr("akamai_property.test", "group_id", "grp_0"),
		resource.TestCheckResourceAttr("akamai_property.test", "group", "grp_0"),
		resource.TestCheckResourceAttr("akamai_property.test", "product", "prd_0"),
		resource.TestCheckResourceAttr("akamai_property.test", "product_id", "prd_0"),
	)

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

	// Run a test case to verify error when the given attribute is changed after creation
	AssertImmutable := func(t *testing.T, attribute string) {
		t.Helper()

		t.Run(fmt.Sprintf("Immutable/%s", attribute), func(t *testing.T) {
			t.Helper()
			client := &mockpapi{}
			client.Test(T{t})

			// We're going to pretend like all of these are valid property identities
			ExpectGetProp(client, "test property", "prp_0", "grp_0", "ctr_0", "prd_0", 42)
			// ExpectGetProp(client, "test property", "prp_0", GroupID2, ContractID2, ProductID2, 42)

			ExpectCreateProperty(client, "test property", "grp_0", "ctr_0", "prd_0", "prp_0").Once()
			ExpectRemoveProperty(client, "prp_0", "ctr_0", "grp_0").Once()

			useClient(client, func() {
				resource.UnitTest(t, resource.TestCase{
					Providers: testAccProviders,
					Steps: []resource.TestStep{
						{
							Config: loadFixtureString("testdata/%s-step0.tf", t.Name()),
							Check:  checkAllAttrs,
						},
						{
							Config:      loadFixtureString("testdata/%s-step1.tf", t.Name()),
							ExpectError: regexp.MustCompile(fmt.Sprintf(`%q cannot be changed`, attribute)),
						},
					},
				})
			})

			client.AssertExpectations(t)
		})
	}

	// Run a test case where a single resource is created, read, and destroyed. The resources of this kind all have the
	// same attributes and vary only by the input terraform config file, which is named after the test case.
	AssertLifecycle := func(t *testing.T, fixtureName string) {
		t.Helper()

		t.Run(fmt.Sprintf("Lifecycle/%s", fixtureName), func(t *testing.T) {
			t.Helper()
			client := &mockpapi{}
			client.Test(T{t})

			ExpectGetProp(client, "test property", "prp_0", "grp_0", "ctr_0", "prd_0", 42)
			ExpectCreateProperty(client, "test property", "grp_0", "ctr_0", "prd_0", "prp_0").Once()
			ExpectRemoveProperty(client, "prp_0", "ctr_0", "grp_0").Once()

			useClient(client, func() {
				resource.UnitTest(t, resource.TestCase{
					Providers: testAccProviders,
					Steps: []resource.TestStep{{
						Config: loadFixtureString("testdata/%s.tf", t.Name()),
						Check:  checkAllAttrs,
					}},
					CheckDestroy: resource.TestCheckNoResourceAttr("akamai_property.test", "id"),
				})
			})

			client.AssertExpectations(t)
		})
	}

	// Run a test case that verifies the resource can be imported by the given ID
	AssertImportable := func(t *testing.T, fixtureName, ImportID string) {
		t.Helper()

		t.Run(fmt.Sprintf("Importable/%s", fixtureName), func(t *testing.T) {
			t.Helper()
			client := &mockpapi{}
			client.Test(T{t})

			Version := 42
			Property := papi.Property{
				PropertyName:   "test property",
				PropertyID:     "prp_0",
				GroupID:        "grp_0",
				ContractID:     "ctr_0",
				ProductID:      "prd_0",
				LatestVersion:  Version,
				StagingVersion: &Version,
			}

			ExpectGetProperty(client, "prp_0", "", "", &Property)
			ExpectGetProperty(client, "prp_0", "grp_0", "ctr_0", &Property)

			ExpectCreateProperty(client, "test property", "grp_0", "ctr_0", "prd_0", "prp_0").Once()
			ExpectRemoveProperty(client, "prp_0", "ctr_0", "grp_0").Once()

			useClient(client, func() {
				resource.UnitTest(t, resource.TestCase{
					Providers: testAccProviders,
					Steps: []resource.TestStep{
						{
							Config: loadFixtureString("testdata/%s.tf", t.Name()),
							Check:  checkAllAttrs,
						},
						{
							ImportState:       true,
							ImportStateVerify: true,
							ImportStateId:     ImportID,
							ResourceName:      "akamai_property.test",
							Config:            loadFixtureString("testdata/%s.tf", t.Name()),
						},
						{
							Config: loadFixtureString("testdata/%s.tf", t.Name()),
							Check:  checkAllAttrs,
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

		AssertDeprecated(t, "contract")
		AssertDeprecated(t, "group")
		AssertDeprecated(t, "product")

		AssertDeprecated(t, "rule_format")
		AssertDeprecated(t, "cp_code")
		AssertDeprecated(t, "contact")
		AssertDeprecated(t, "hostnames")
		AssertDeprecated(t, "origin")
		AssertDeprecated(t, "is_secure")
		AssertDeprecated(t, "rules")
		AssertDeprecated(t, "variables")

		AssertForbiddenAttr(t, "rule_format")
		AssertForbiddenAttr(t, "cp_code")
		AssertForbiddenAttr(t, "contact")
		AssertForbiddenAttr(t, "hostnames")
		AssertForbiddenAttr(t, "origin")
		AssertForbiddenAttr(t, "is_secure")
		AssertForbiddenAttr(t, "rules")
		AssertForbiddenAttr(t, "variables")

		AssertLifecycle(t, "contract instead of contract_id")
		AssertLifecycle(t, "contract_id without prefix")
		AssertLifecycle(t, "group instead of group_id")
		AssertLifecycle(t, "group_id without prefix")
		AssertLifecycle(t, "product instead of product_id")
		AssertLifecycle(t, "product_id without prefix")

		AssertImmutable(t, "contract_id")
		AssertImmutable(t, "contract")

		AssertImmutable(t, "group_id")
		AssertImmutable(t, "group")

		AssertImmutable(t, "product_id")
		AssertImmutable(t, "product")

		AssertImportable(t, "property_id", "prp_0")
		AssertImportable(t, "unprefixed property_id", "0")

		t.Run("property is destroyed and recreated when name is changed", func(t *testing.T) {
			client := &mockpapi{}
			client.Test(T{t})

			ExpectGetProp(client, "test property", "prp_0", "grp_0", "ctr_0", "prd_0", 42)
			ExpectGetProp(client, "renamed property", "prp_1", "grp_0", "ctr_0", "prd_0", 1)
			ExpectCreateProperty(client, "test property", "grp_0", "ctr_0", "prd_0", "prp_0").Once()
			ExpectCreateProperty(client, "renamed property", "grp_0", "ctr_0", "prd_0", "prp_1").Once()
			ExpectRemoveProperty(client, "prp_0", "ctr_0", "grp_0").Once()
			ExpectRemoveProperty(client, "prp_1", "ctr_0", "grp_0").Once()

			useClient(client, func() {
				resource.UnitTest(t, resource.TestCase{
					Providers: testAccProviders,
					Steps: []resource.TestStep{
						{
							Config: loadFixtureString("testdata/%s-step0.tf", t.Name()),
							Check: resource.ComposeAggregateTestCheckFunc(
								resource.TestCheckResourceAttr("akamai_property.test", "id", "prp_0"),
								resource.TestCheckResourceAttr("akamai_property.test", "name", "test property"),
							),
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

			ExpectGetProp(client, "test property", "prp_0", "grp_0", "ctr_0", "prd_0", 42)
			ExpectCreateProperty(client, "test property", "grp_0", "ctr_0", "prd_0", "prp_0").Once()

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
							Config: loadFixtureString("testdata/%s-step0.tf", t.Name()),
							Check:  checkAllAttrs,
						},
						{
							Config:      loadFixtureString("testdata/%s-step1.tf", t.Name()),
							ExpectError: regexp.MustCompile(`Cannot remove active property`),
						},
					},
				})
			})

			client.AssertExpectations(t)
		})
	})
}

// Sets up an expected call to papi.CreateProperty() with a constant success response with the given PropertyID
func ExpectCreateProperty(client *mockpapi, PropertyName, GroupID, ContractID, ProductID, PropertyID string) *mock.Call {
	req := papi.CreatePropertyRequest{
		GroupID:    GroupID,
		ContractID: ContractID,
		Property: papi.PropertyCreate{
			ProductID:    ProductID,
			PropertyName: PropertyName,
		},
	}

	res := papi.CreatePropertyResponse{PropertyID: PropertyID}

	return client.On("CreateProperty", AnyCTX, req).Return(&res, nil)
}

// Sets up an expected call to papi.RemoveProperty() with a constant success response
func ExpectRemoveProperty(client *mockpapi, PropertyID, ContractID, GroupID string) *mock.Call {
	req := papi.RemovePropertyRequest{
		PropertyID: PropertyID,
		GroupID:    GroupID,
		ContractID: ContractID,
	}
	res := papi.RemovePropertyResponse{}

	return client.On("RemoveProperty", AnyCTX, req).Return(&res, nil)
}
