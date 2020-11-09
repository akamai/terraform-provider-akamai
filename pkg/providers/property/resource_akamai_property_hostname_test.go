package property

import (
	"context"
	"fmt"
	"regexp"
	"testing"

	"github.com/akamai/AkamaiOPEN-edgegrid-golang/v2/pkg/papi"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/resource"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/mock"
)

func TestResPropHostname(t *testing.T) {
	// TestCheckFunc to verify all standard attributes
	checkAllAttrs := resource.ComposeAggregateTestCheckFunc(
		resource.TestCheckResourceAttr("akamai_property_hostnames.test", "id", "prp_0"),
		resource.TestCheckResourceAttr("akamai_property_hostnames.test", "property_id", "prp_0"),
		resource.TestCheckResourceAttr("akamai_property_hostnames.test", "contract_id", "ctr_0"),
		resource.TestCheckResourceAttr("akamai_property_hostnames.test", "group_id", "grp_0"),
		resource.TestCheckResourceAttr("akamai_property_hostnames.test", "names.test.domain", "ehn_0"),
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

	// Defines a standard set of client behaviors for this test
	type ClientBehavior struct {
		Name        string
		ClientSetup func(*mockpapi)
	}

	// Standard test behavior for cases where the property's latest version is active in staging network
	LatestVersionActiveInStaging := ClientBehavior{
		Name: "Latest version is active in staging",
		ClientSetup: func(client *mockpapi) {
			// The property state for this test
			PropertyState := papi.Property{
				PropertyID:    "prp_0",
				GroupID:       "grp_0",
				ContractID:    "ctr_0",
				LatestVersion: 1,
			}

			// Property version is active
			PropertyState.ProductionVersion = &PropertyState.LatestVersion

			// The state of the test property's associated Hostnames at version 1
			Hostnames1 := []papi.Hostname{}

			// The state of the test property's associated Hostnames at version 2
			Hostnames2 := []papi.Hostname{}

			ExpectGetProperty(client, "prp_0", "grp_0", "ctr_0", &PropertyState)
			ExpectGetPropertyVersionHostnames(client, "prp_0", "grp_0", "ctr_0", 1, &Hostnames1)
			ExpectGetPropertyVersionHostnames(client, "prp_0", "grp_0", "ctr_0", 2, &Hostnames2)
			ExpectCreatePropertyVersion(client, "prp_0", "grp_0", "ctr_0", 1, 2, &PropertyState).Once()

			NewHostnames := []papi.Hostname{{
				CnameType:      "EDGE_HOSTNAME",
				CnameFrom:      "test.domain",
				EdgeHostnameID: "ehn_0",
			}}
			ExpectUpdatePropertyVersionHostnames(client, "prp_0", "grp_0", "ctr_0", 2, NewHostnames, &Hostnames2).Once()

			// No PAPI interaction for delete action for this resource type
		},
	}

	// Standard test behavior for cases where the property's latest version is active in production network
	LatestVersionActiveInProd := ClientBehavior{
		Name: "Latest version is active in production",
		ClientSetup: func(client *mockpapi) {
			// The property state for this test
			PropertyState := papi.Property{
				PropertyID:    "prp_0",
				GroupID:       "grp_0",
				ContractID:    "ctr_0",
				LatestVersion: 1,
			}

			// Property version is active
			PropertyState.StagingVersion = &PropertyState.LatestVersion

			// The state of the test property's associated Hostnames at version 1
			Hostnames1 := []papi.Hostname{}

			// The state of the test property's associated Hostnames at version 2
			Hostnames2 := []papi.Hostname{}

			ExpectGetProperty(client, "prp_0", "grp_0", "ctr_0", &PropertyState)
			ExpectGetPropertyVersionHostnames(client, "prp_0", "grp_0", "ctr_0", 1, &Hostnames1)
			ExpectGetPropertyVersionHostnames(client, "prp_0", "grp_0", "ctr_0", 2, &Hostnames2)
			ExpectCreatePropertyVersion(client, "prp_0", "grp_0", "ctr_0", 1, 2, &PropertyState).Once()

			NewHostnames := []papi.Hostname{{
				CnameType:      "EDGE_HOSTNAME",
				CnameFrom:      "test.domain",
				EdgeHostnameID: "ehn_0",
			}}
			ExpectUpdatePropertyVersionHostnames(client, "prp_0", "grp_0", "ctr_0", 2, NewHostnames, &Hostnames2).Once()

			// No PAPI interaction for delete action for this resource type
		},
	}

	// Standard test behavior for cases where the property's latest version is not active
	LatestVersionNotActive := ClientBehavior{
		Name: "Latest version not active",
		ClientSetup: func(client *mockpapi) {
			// The property state for this test
			Property := papi.Property{
				PropertyID:    "prp_0",
				GroupID:       "grp_0",
				ContractID:    "ctr_0",
				LatestVersion: 1,
			}

			// The state of the test property's associated Hostnames at version 1 (the only version in this test)
			Hostnames := []papi.Hostname{}

			ExpectGetProperty(client, "prp_0", "grp_0", "ctr_0", &Property)
			ExpectGetPropertyVersionHostnames(client, "prp_0", "grp_0", "ctr_0", 1, &Hostnames)

			NewHostnames := []papi.Hostname{{
				CnameType:      "EDGE_HOSTNAME",
				CnameFrom:      "test.domain",
				EdgeHostnameID: "ehn_0",
			}}
			ExpectUpdatePropertyVersionHostnames(client, "prp_0", "grp_0", "ctr_0", 1, NewHostnames, &Hostnames).Once()

			// No create version

			// No PAPI interaction for delete action for this resource type
		},
	}

	// Standard test behavior for cases where the property's latest version already has exactly the desired hostnames
	HostnamesAlreadySet := ClientBehavior{
		Name: "Hostnames already set to desired state",
		ClientSetup: func(client *mockpapi) {
			// The property state for this test
			Property := papi.Property{
				PropertyID:    "prp_0",
				GroupID:       "grp_0",
				ContractID:    "ctr_0",
				LatestVersion: 1,
			}

			// The state of the test property's associated Hostnames at version 1, already contains desired state
			Hostnames := []papi.Hostname{{
				CnameType:      "EDGE_HOSTNAME",
				CnameFrom:      "test.domain",
				EdgeHostnameID: "ehn_0",
			}}

			ExpectGetProperty(client, "prp_0", "grp_0", "ctr_0", &Property)
			ExpectGetPropertyVersionHostnames(client, "prp_0", "grp_0", "ctr_0", 1, &Hostnames)

			// No create version

			// No update hostnames

			// No PAPI interaction for delete action for this resource type
		},
	}

	// Run a happy-path test case that goes through a complete create-destroy cycle
	AssertLifecycle := func(t *testing.T, fixture string, behavior ClientBehavior) {
		t.Helper()

		fixturePath := fmt.Sprintf("testdata/%s/Lifecycle/%s.tf", t.Name(), fixture)
		testName := fmt.Sprintf("Lifecycle/%s/%s", fixture, behavior.Name)

		t.Run(testName, func(t *testing.T) {
			t.Helper()

			client := &mockpapi{}
			behavior.ClientSetup(client)

			useClient(client, func() {
				resource.UnitTest(t, resource.TestCase{
					Providers: testAccProviders,
					Steps: []resource.TestStep{{
						Config: loadFixtureString(fixturePath),
						Check:  checkAllAttrs,
					}},
					CheckDestroy: resource.TestCheckNoResourceAttr("akamai_property_hostnames.test", "id"),
				})
			})

			client.AssertExpectations(t)
		})
	}

	ImportByPropertyID := ClientBehavior{
		Name: "Import by property_id",
		ClientSetup: func(client *mockpapi) {
			HostnamesAlreadySet.ClientSetup(client)
			Property := papi.Property{
				PropertyID:    "prp_0",
				GroupID:       "grp_0",
				ContractID:    "ctr_0",
				LatestVersion: 1,
			}

			ExpectGetProperty(client, "prp_0", "", "", &Property)
		},
	}

	ImportByPropertyGroupIDs := ClientBehavior{
		Name: "Import by property_id and group_id",
		ClientSetup: func(client *mockpapi) {
			HostnamesAlreadySet.ClientSetup(client)

			Property := papi.Property{
				PropertyID:    "prp_0",
				GroupID:       "grp_0",
				ContractID:    "ctr_0",
				LatestVersion: 1,
			}

			ExpectGetProperty(client, "prp_0", "grp_0", "", &Property)
		},
	}

	ImportByPropertyGroupContractIDs := ClientBehavior{
		Name: "Import by property_id group_id and contract_id",
		ClientSetup: func(client *mockpapi) {
			HostnamesAlreadySet.ClientSetup(client)
		},
	}

	// Run a test case that verifies the resource can be imported by the given ID
	AssertImportable := func(t *testing.T, TestName, ImportID string, behavior ClientBehavior) {
		t.Helper()

		fixturePath := fmt.Sprintf("testdata/%s/Importable/importable.tf", t.Name())
		testName := fmt.Sprintf("Importable/%s/%s", behavior.Name, TestName)

		t.Run(testName, func(t *testing.T) {
			t.Helper()

			client := &mockpapi{}
			behavior.ClientSetup(client)

			useClient(client, func() {
				resource.UnitTest(t, resource.TestCase{
					Providers: testAccProviders,
					Steps: []resource.TestStep{
						{
							Config: loadFixtureString(fixturePath),
							Check:  checkAllAttrs,
						},
						{
							ImportState:       true,
							ImportStateVerify: true,
							ImportStateId:     ImportID,
							ResourceName:      "akamai_property_hostnames.test",
							Config:            loadFixtureString(fixturePath),
						},
						{
							Config: loadFixtureString(fixturePath),
							Check:  checkAllAttrs,
						},
					},
				})
			})

			client.AssertExpectations(t)
		})
	}

	suppressLogging(t, func() {
		AssertConfigError(t, "property_id not given", `"property_id" is required`)
		AssertConfigError(t, "contract_id not given", `"contract_id" is required`)
		AssertConfigError(t, "group_id not given", `"group_id" is required`)
		AssertConfigError(t, "names not given", `"names" is required`)

		AssertLifecycle(t, "all normal attributes", LatestVersionNotActive)
		AssertLifecycle(t, "all normal attributes", LatestVersionActiveInStaging)
		AssertLifecycle(t, "all normal attributes", LatestVersionActiveInProd)
		AssertLifecycle(t, "all normal attributes", HostnamesAlreadySet)

		AssertLifecycle(t, "contract_id without prefix", LatestVersionNotActive)
		AssertLifecycle(t, "contract_id without prefix", LatestVersionActiveInStaging)
		AssertLifecycle(t, "contract_id without prefix", LatestVersionActiveInProd)
		AssertLifecycle(t, "contract_id without prefix", HostnamesAlreadySet)

		AssertLifecycle(t, "group_id without prefix", LatestVersionNotActive)
		AssertLifecycle(t, "group_id without prefix", LatestVersionActiveInStaging)
		AssertLifecycle(t, "group_id without prefix", LatestVersionActiveInProd)
		AssertLifecycle(t, "group_id without prefix", HostnamesAlreadySet)

		AssertLifecycle(t, "property_id without prefix", LatestVersionNotActive)
		AssertLifecycle(t, "property_id without prefix", LatestVersionActiveInStaging)
		AssertLifecycle(t, "property_id without prefix", LatestVersionActiveInProd)
		AssertLifecycle(t, "property_id without prefix", HostnamesAlreadySet)

		AssertImportable(t, "property_id", "prp_0", ImportByPropertyID)
		AssertImportable(t, "unprefixed property_id", "0", ImportByPropertyID)

		AssertImportable(t, "property_id and group_id", "prp_0,grp_0", ImportByPropertyGroupIDs)
		AssertImportable(t, "unprefixed property_id and group_id", "0,0", ImportByPropertyGroupIDs)

		AssertImportable(t, "property_id and group_id and contract_id", "prp_0,grp_0,ctr_0", ImportByPropertyGroupContractIDs)
		AssertImportable(t, "unprefixed property_id and group_id and contract_id", "0,0,0", ImportByPropertyGroupContractIDs)

		t.Run("out-of-band/change group_id and contract_id", func(t *testing.T) {
			// This test demonstrates the correct way for someone to change the contract or group for a property that is
			// already managed by Terraform. This test walks through two such out-of-band changes, one for group and
			// another for contract, but they can be changed together in the same import. Steps:
			//   1. Change the contract and/or group outside of terraform
			//   2. Edit the terraform config and adjust the contract and/or group
			//   3. Use `terraform import` command to update the terraform state

			client := &mockpapi{}

			Property := papi.Property{
				PropertyID:    "prp_0",
				GroupID:       "grp_0",
				ContractID:    "ctr_0",
				LatestVersion: 1,
			}

			// The state of the test property's associated Hostnames at version 1, already contains desired state
			Hostnames := []papi.Hostname{{
				CnameType:      "EDGE_HOSTNAME",
				CnameFrom:      "test.domain",
				EdgeHostnameID: "ehn_0",
			}}

			// Test impl verifies the request params match the current state values and errors when they don't match
			// Expectation is that the group_id is changed between steps, so the requested group_id will track the change
			client.OnGetProperty(AnyCTX, mock.Anything, func(_ context.Context, req papi.GetPropertyRequest) (*papi.GetPropertyResponse, error) {
				matchPID := req.PropertyID == Property.PropertyID
				matchGID := req.GroupID == Property.GroupID
				matchCID := req.ContractID == Property.ContractID

				if !matchPID || !matchGID || !matchCID {
					// The requested property doesn't exist (because we changed the group)
					return nil, fmt.Errorf("Not found")
				}

				// Copy the property state and return to caller
				property := Property

				return &papi.GetPropertyResponse{Property: &property}, nil
			})

			// Expect calls for every permutation of identity values
			ExpectGetPropertyVersionHostnames(client, "prp_0", "grp_0", "ctr_0", 1, &Hostnames)
			ExpectGetPropertyVersionHostnames(client, "prp_0", "grp_1", "ctr_0", 1, &Hostnames)
			ExpectGetPropertyVersionHostnames(client, "prp_0", "grp_1", "ctr_1", 1, &Hostnames)

			useClient(client, func() {
				resource.UnitTest(t, resource.TestCase{
					Providers: testAccProviders,
					Steps: []resource.TestStep{
						{
							Config: loadFixtureString("testdata/%s/step0.tf", t.Name()),
							Check:  checkAllAttrs,
						},
						{
							PreConfig: func() {
								// This is our first out-of-band change. Fixture matches these value for this step.
								Property.GroupID = "grp_1"
							},
							ImportState:       true,
							ImportStateVerify: true,
							ImportStateId:     "prp_0,grp_1,ctr_0",
							ResourceName:      "akamai_property_hostnames.test",
							Config:            loadFixtureString("testdata/%s/step1.tf", t.Name()),
							Check: resource.ComposeAggregateTestCheckFunc(
								resource.TestCheckResourceAttr("akamai_property_hostnames.test", "id", "prp_0"),
								resource.TestCheckResourceAttr("akamai_property_hostnames.test", "property_id", "prp_0"),
								resource.TestCheckResourceAttr("akamai_property_hostnames.test", "contract_id", "ctr_0"),
								resource.TestCheckResourceAttr("akamai_property_hostnames.test", "group_id", "grp_1"),
								resource.TestCheckResourceAttr("akamai_property_hostnames.test", "names.test.domain", "ehn_0"),
							),
						},
						{
							PreConfig: func() {
								// This is our second out-of-band change. Fixture matches these value for this step.
								Property.ContractID = "ctr_1"
							},
							ImportState:       true,
							ImportStateVerify: true,
							ImportStateId:     "prp_0,grp_1,ctr_1",
							ResourceName:      "akamai_property_hostnames.test",
							Config:            loadFixtureString("testdata/%s/step2.tf", t.Name()),
							Check: resource.ComposeAggregateTestCheckFunc(
								resource.TestCheckResourceAttr("akamai_property_hostnames.test", "id", "prp_0"),
								resource.TestCheckResourceAttr("akamai_property_hostnames.test", "property_id", "prp_0"),
								resource.TestCheckResourceAttr("akamai_property_hostnames.test", "contract_id", "ctr_1"),
								resource.TestCheckResourceAttr("akamai_property_hostnames.test", "group_id", "grp_1"),
								resource.TestCheckResourceAttr("akamai_property_hostnames.test", "names.test.domain", "ehn_0"),
							),
						},
					},
					CheckDestroy: resource.TestCheckNoResourceAttr("akamai_property_hostnames.test", "id"),
				})
			})

			client.AssertExpectations(t)
		})

		t.Run("update hostnames", func(t *testing.T) {
			client := &mockpapi{}

			Property := papi.Property{
				PropertyID:    "prp_0",
				GroupID:       "grp_0",
				ContractID:    "ctr_0",
				LatestVersion: 1,
			}

			// The state of the test property's associated Hostnames at version 1, already contains desired state
			Hostnames := []papi.Hostname{}

			ExpectGetProperty(client, "prp_0", "grp_0", "ctr_0", &Property)
			ExpectGetPropertyVersionHostnames(client, "prp_0", "grp_0", "ctr_0", 1, &Hostnames)

			// For step0 - CREATE
			NewHostnames1 := []papi.Hostname{{
				CnameType:      "EDGE_HOSTNAME",
				CnameFrom:      "test.domain",
				EdgeHostnameID: "ehn_0",
			}}
			ExpectUpdatePropertyVersionHostnames(client, "prp_0", "grp_0", "ctr_0", 1, NewHostnames1, &Hostnames).Once()

			// For step1 - UPDATE
			NewHostnames2 := []papi.Hostname{{
				CnameType:      "EDGE_HOSTNAME",
				CnameFrom:      "test.domain2",
				EdgeHostnameID: "ehn_1",
			}}
			ExpectUpdatePropertyVersionHostnames(client, "prp_0", "grp_0", "ctr_0", 1, NewHostnames2, &Hostnames).Once()

			useClient(client, func() {
				resource.UnitTest(t, resource.TestCase{
					Providers: testAccProviders,
					Steps: []resource.TestStep{
						{
							Config: loadFixtureString("testdata/%s/step0.tf", t.Name()),
							Check:  checkAllAttrs,
						},
						{
							Config: loadFixtureString("testdata/%s/step1.tf", t.Name()),
							Check: resource.ComposeAggregateTestCheckFunc(
								resource.TestCheckResourceAttr("akamai_property_hostnames.test", "id", "prp_0"),
								resource.TestCheckResourceAttr("akamai_property_hostnames.test", "property_id", "prp_0"),
								resource.TestCheckResourceAttr("akamai_property_hostnames.test", "contract_id", "ctr_0"),
								resource.TestCheckResourceAttr("akamai_property_hostnames.test", "group_id", "grp_0"),
								resource.TestCheckResourceAttr("akamai_property_hostnames.test", "names.test.domain2", "ehn_1"),
							),
						},
					},
					CheckDestroy: resource.TestCheckNoResourceAttr("akamai_property_hostnames.test", "id"),
				})
			})

			assert.Equal(t, NewHostnames2, Hostnames)
			client.AssertExpectations(t)
		})

		t.Run("Immutable/property_id", func(t *testing.T) {
			client := &mockpapi{}

			// The property state for this test
			Property := papi.Property{
				PropertyID:    "prp_0",
				GroupID:       "grp_0",
				ContractID:    "ctr_0",
				LatestVersion: 1,
			}

			// The state of the test property's associated Hostnames at version 1, already contains desired state
			Hostnames := []papi.Hostname{{
				CnameType:      "EDGE_HOSTNAME",
				CnameFrom:      "test.domain",
				EdgeHostnameID: "ehn_0",
			}}

			ExpectGetProperty(client, "prp_0", "grp_0", "ctr_0", &Property)
			ExpectGetPropertyVersionHostnames(client, "prp_0", "grp_0", "ctr_0", 1, &Hostnames)
			useClient(client, func() {
				resource.UnitTest(t, resource.TestCase{
					Providers: testAccProviders,
					Steps: []resource.TestStep{
						{
							Config: loadFixtureString("testdata/%s/step0.tf", t.Name()),
							Check:  checkAllAttrs,
						},
						{
							Config:      loadFixtureString("testdata/%s/step1.tf", t.Name()),
							ExpectError: regexp.MustCompile(`"property_id" cannot be changed`),
						},
					},
				})
			})

			client.AssertExpectations(t)
		})
	})
}

// Sets up an expected call to papi.GetProperty() which returns a value depending on the given State pointer. When nil,
// the PAPI response contains a zero-value papi.Property. Otherwise the response will dynamically contain a copy of
// the State made at the time of the call to mockpapi.GetProperty().
func ExpectGetProperty(client *mockpapi, PropertyID, GroupID, ContractID string, State *papi.Property) *mock.Call {
	req := papi.GetPropertyRequest{
		PropertyID: PropertyID,
		ContractID: ContractID,
		GroupID:    GroupID,
	}

	fn := func(context.Context, papi.GetPropertyRequest) (*papi.GetPropertyResponse, error) {
		var property papi.Property

		// Duplicate the State
		if State != nil {
			property = *State
		}

		// Duplicate the pointers
		if property.ProductionVersion != nil {
			v := *property.ProductionVersion
			property.ProductionVersion = &v
		}

		if property.StagingVersion != nil {
			v := *property.StagingVersion
			property.StagingVersion = &v
		}

		return &papi.GetPropertyResponse{Property: &property}, nil
	}

	return client.OnGetProperty(AnyCTX, req, fn)
}

// Sets up an expected call to papi.GetPropertyVersionHostnames() which returns a value depending on the value of the
// pointer to State. When nil or empty, the response contains a nil Items member. Otherwise the response contains a
// copy of the value pointed to by State made at the time of the call to papi.GetPropertyVersionHostnames().
func ExpectGetPropertyVersionHostnames(client *mockpapi, PropertyID, GroupID, ContractID string, PropertyVersion int, State *[]papi.Hostname) *mock.Call {
	req := papi.GetPropertyVersionHostnamesRequest{
		PropertyID:      PropertyID,
		GroupID:         GroupID,
		ContractID:      ContractID,
		PropertyVersion: PropertyVersion,
	}

	fn := func(context.Context, papi.GetPropertyVersionHostnamesRequest) (*papi.GetPropertyVersionHostnamesResponse, error) {
		var Items []papi.Hostname
		if len(*State) > 0 {
			// Duplicate the State
			Items = append(Items, *State...)
		}

		res := papi.GetPropertyVersionHostnamesResponse{
			ContractID:      ContractID,
			GroupID:         GroupID,
			PropertyID:      PropertyID,
			PropertyVersion: PropertyVersion,
			Hostnames:       papi.HostnameResponseItems{Items: Items},
		}

		return &res, nil
	}

	return client.OnGetPropertyVersionHostnames(AnyCTX, req, fn)
}

// Sets up an expected call to papi.UpdatePropertyVersionHostnames() which returns a constant value based on input
// params. If given, the value pointed to by State will be updated with a copy of the given Hostnames when the call
// to mockpapi.UpdatePropertyVersionHostnames() is made.
func ExpectUpdatePropertyVersionHostnames(client *mockpapi, PropertyID, GroupID, ContractID string, PropertyVersion int, Hostnames []papi.Hostname, State *[]papi.Hostname) *mock.Call {
	req := papi.UpdatePropertyVersionHostnamesRequest{
		PropertyID:      PropertyID,
		PropertyVersion: PropertyVersion,
		ContractID:      ContractID,
		GroupID:         GroupID,
		Hostnames:       Hostnames,
	}

	res := papi.UpdatePropertyVersionHostnamesResponse{
		ContractID:      ContractID,
		GroupID:         GroupID,
		PropertyID:      PropertyID,
		PropertyVersion: PropertyVersion,
		Hostnames:       papi.HostnameResponseItems{Items: Hostnames},
	}

	call := client.On("UpdatePropertyVersionHostnames", AnyCTX, req).Return(&res, nil)
	if State != nil {
		call.Run(func(mock.Arguments) {
			*State = append([]papi.Hostname{}, Hostnames...)
		})
	}

	return call
}

// Sets up an expected call to papi.CreatePropertyVersion() with a constant response value based on input parameters.
// If given, the Property pointed to by State will have its LatestVersion updated to the NewVersion at the time
// mock.CreatePropertyVersion() is called.
func ExpectCreatePropertyVersion(client *mockpapi, PropertyID, GroupID, ContractID string, CreateFromVersion, NewVersion int, State *papi.Property) *mock.Call {
	req := papi.CreatePropertyVersionRequest{
		PropertyID: PropertyID,
		GroupID:    GroupID,
		ContractID: ContractID,
		Version: papi.PropertyVersionCreate{
			CreateFromVersion: CreateFromVersion,
		},
	}

	res := papi.CreatePropertyVersionResponse{PropertyVersion: NewVersion}

	call := client.On("CreatePropertyVersion", AnyCTX, req).Return(&res, nil)

	if State != nil {
		call.Run(func(mock.Arguments) { State.LatestVersion = NewVersion })
	}

	return call
}
