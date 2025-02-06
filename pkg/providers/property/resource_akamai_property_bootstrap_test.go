package property

import (
	"fmt"
	"regexp"
	"testing"

	"github.com/akamai/AkamaiOPEN-edgegrid-golang/v10/pkg/iam"
	"github.com/akamai/AkamaiOPEN-edgegrid-golang/v10/pkg/papi"
	"github.com/akamai/terraform-provider-akamai/v7/pkg/common/test"
	"github.com/akamai/terraform-provider-akamai/v7/pkg/common/testutils"
	"github.com/hashicorp/terraform-plugin-testing/helper/resource"
	"github.com/hashicorp/terraform-plugin-testing/terraform"
)

var basicDataBootstrap = mockPropertyData{
	propertyID:   "prp_123",
	propertyName: "property_name",
	groupID:      "grp_1",
	contractID:   "ctr_2",
	productID:    "prd_3",
	assetID:      "aid_55555",
	moveGroup: moveGroup{
		sourceGroupID:      1,
		destinationGroupID: 111,
	},
}

func TestBootstrapResourceCreate(t *testing.T) {
	t.Parallel()

	baseChecker := test.NewStateChecker("akamai_property_bootstrap.test").
		CheckEqual("id", "prp_123").
		CheckEqual("group_id", "grp_1").
		CheckEqual("contract_id", "ctr_2").
		CheckEqual("product_id", "prd_3").
		CheckEqual("name", "property_name").
		CheckEqual("asset_id", "aid_55555")

	tests := map[string]struct {
		init  func(*mockProperty)
		steps []resource.TestStep
		error *regexp.Regexp
	}{
		"create": {
			init: func(p *mockProperty) {
				p.mockCreateProperty()
				p.mockGetProperty().Twice()
				p.mockRemoveProperty()
			},
			steps: []resource.TestStep{
				{
					Config: testutils.LoadFixtureString(t, "testdata/TestResPropertyBootstrap/create.tf"),
					Check:  baseChecker.Build(),
				},
			},
		},
		"create without prefixes": {
			init: func(p *mockProperty) {
				p.mockCreateProperty()
				p.mockGetProperty().Twice()
				p.mockRemoveProperty()
			},
			steps: []resource.TestStep{
				{
					Config: testutils.LoadFixtureString(t, "testdata/TestResPropertyBootstrap/create_without_prefixes.tf"),
					Check: baseChecker.
						CheckEqual("group_id", "1").
						CheckEqual("contract_id", "2").
						CheckEqual("product_id", "3").
						Build(),
				},
			},
		},
		"create with interpretCreate error - group not found": {
			init: func(p *mockProperty) {
				req := papi.CreatePropertyRequest{
					GroupID:    p.groupID,
					ContractID: p.contractID,
					Property: papi.PropertyCreate{
						ProductID:    p.productID,
						PropertyName: p.propertyName,
					},
				}
				p.papiMock.On("CreateProperty", testutils.MockContext, req).Return(nil, fmt.Errorf(
					"%s: %w: %s", papi.ErrCreateProperty, papi.ErrNotFound, "not found")).Once()
				// mock empty groups - no group has been found, hence the expected error
				p.papiMock.On("GetGroups", testutils.MockContext).Return(&papi.GetGroupsResponse{
					Groups: papi.GroupItems{
						Items: []*papi.Group{},
					},
				}, nil).Once()
			},
			steps: []resource.TestStep{
				{
					Config:      testutils.LoadFixtureString(t, "testdata/TestResPropertyBootstrap/create.tf"),
					ExpectError: regexp.MustCompile(`Error: group not found: grp_1`),
				},
			},
		},
	}

	for name, test := range tests {
		t.Run(name, func(t *testing.T) {
			t.Parallel()
			m := &papi.Mock{}
			mp := &mockProperty{
				mockPropertyData: basicDataBootstrap,
				papiMock:         m,
			}
			if test.init != nil {
				test.init(mp)
			}

			useClient(m, nil, func() {
				resource.UnitTest(t, resource.TestCase{
					ProtoV6ProviderFactories: testutils.NewProtoV6ProviderFactory(NewSubprovider()),
					IsUnitTest:               true,
					Steps:                    test.steps,
				})
			})

			m.AssertExpectations(t)
		})
	}
}

func TestBootstrapResourceUpdate(t *testing.T) {
	t.Parallel()

	baseChecker := test.NewStateChecker("akamai_property_bootstrap.test").
		CheckEqual("id", "prp_123").
		CheckEqual("group_id", "grp_1").
		CheckEqual("contract_id", "ctr_2").
		CheckEqual("product_id", "prd_3").
		CheckEqual("name", "property_name").
		CheckEqual("asset_id", "aid_55555")

	tests := map[string]struct {
		configPathForCreate string
		configPathForUpdate string
		init                func(*mockProperty)
		errorForCreate      *regexp.Regexp
		errorForUpdate      *regexp.Regexp
		updateChecks        resource.TestCheckFunc
	}{
		"create and remove prefixes - no diff": {
			configPathForCreate: "testdata/TestResPropertyBootstrap/create.tf",
			configPathForUpdate: "testdata/TestResPropertyBootstrap/create_without_prefixes.tf",
			init: func(p *mockProperty) {
				p.mockCreateProperty()
				p.mockGetProperty()
				// read x2
				p.mockGetProperty().Twice()
				// read x1 before update
				p.mockGetProperty()
				p.mockRemoveProperty()
			},
			updateChecks: baseChecker.Build(),
		},
		"create and update group id": {
			configPathForCreate: "testdata/TestResPropertyBootstrap/create.tf",
			configPathForUpdate: "testdata/TestResPropertyBootstrap/update_group.tf",
			init: func(p *mockProperty) {
				p.mockCreateProperty()
				p.mockGetProperty()
				// read x2
				p.mockGetProperty().Twice()
				// update
				p.mockMoveProperty()
				p.groupID = "grp_111"
				// read x2
				p.mockGetProperty().Twice()
				p.mockRemoveProperty()
			},
			updateChecks: baseChecker.
				CheckEqual("group_id", "grp_111").
				Build(),
		},
		"create and update name - resource replacement": {
			configPathForCreate: "testdata/TestResPropertyBootstrap/create.tf",
			init: func(p *mockProperty) {
				p.mockCreateProperty()
				p.mockGetProperty()
				// read x2
				p.mockGetProperty().Twice()
				p.mockRemoveProperty()
				p.propertyName = "property_name2"
				p.mockCreateProperty()
				p.mockGetProperty().Twice()
				p.mockRemoveProperty()
			},
			configPathForUpdate: "testdata/TestResPropertyBootstrap/update_name.tf",
			updateChecks: baseChecker.
				CheckEqual("name", "property_name2").
				Build(),
		},
		"create and update name and group id - resource replacement": {
			configPathForCreate: "testdata/TestResPropertyBootstrap/create.tf",
			configPathForUpdate: "testdata/TestResPropertyBootstrap/update_name_and_group.tf",
			init: func(p *mockProperty) {
				p.mockCreateProperty()
				p.mockGetProperty().Times(3)
				p.mockRemoveProperty()
				p.propertyName = "property_name2"
				p.groupID = "grp_93"
				p.mockCreateProperty()
				p.mockGetProperty().Twice()
				p.mockRemoveProperty()
			},
			updateChecks: baseChecker.
				CheckEqual("name", "property_name2").
				CheckEqual("group_id", "grp_93").
				Build(),
		},
		"create and update contract - error": {
			configPathForCreate: "testdata/TestResPropertyBootstrap/create.tf",
			configPathForUpdate: "testdata/TestResPropertyBootstrap/update_contract.tf",
			init: func(p *mockProperty) {
				p.mockCreateProperty()
				p.mockGetProperty().Times(3)
				p.mockRemoveProperty()
			},
			errorForUpdate: regexp.MustCompile("updating field `contract_id` is not possible"),
		},
		"create and update product - error": {
			configPathForCreate: "testdata/TestResPropertyBootstrap/create.tf",
			configPathForUpdate: "testdata/TestResPropertyBootstrap/update_product.tf",
			init: func(p *mockProperty) {
				p.mockCreateProperty()
				p.mockGetProperty().Times(3)
				p.mockRemoveProperty()
			},
			errorForUpdate: regexp.MustCompile("updating field `product_id` is not possible"),
		},
	}

	for name, test := range tests {
		t.Run(name, func(t *testing.T) {
			t.Parallel()

			iamMock := &iam.Mock{}
			papiMock := &papi.Mock{}
			mp := &mockProperty{
				mockPropertyData: basicDataBootstrap,
				papiMock:         papiMock,
				iamMock:          iamMock,
			}

			if test.init != nil {
				test.init(mp)
			}

			useClient(papiMock, nil, func() {
				useIam(iamMock, func() {
					resource.UnitTest(t, resource.TestCase{
						ProtoV6ProviderFactories: testutils.NewProtoV6ProviderFactory(NewSubprovider()),
						IsUnitTest:               true,
						Steps: []resource.TestStep{
							{
								Config:      testutils.LoadFixtureString(t, test.configPathForCreate),
								Check:       baseChecker.Build(),
								ExpectError: test.errorForCreate,
							},
							{
								Config:      testutils.LoadFixtureString(t, test.configPathForUpdate),
								Check:       test.updateChecks,
								ExpectError: test.errorForUpdate,
							},
						},
					})
				})
			})

			papiMock.AssertExpectations(t)
		})
	}
}

func TestBootstrapResourceImport(t *testing.T) {
	t.Parallel()

	basicDataWithoutContractAndGroup := mockPropertyData{
		propertyID:   "prp_123",
		propertyName: "property_name",
		productID:    "prd_3",
		assetID:      "aid_55555",
	}

	baseChecker := test.NewImportChecker().
		CheckEqual("id", "prp_123").
		CheckEqual("group_id", "grp_1").
		CheckEqual("contract_id", "ctr_2").
		CheckEqual("product_id", "prd_3").
		CheckEqual("name", "property_name").
		CheckEqual("asset_id", "aid_55555")

	tests := map[string]struct {
		init          func(*mockProperty)
		mockData      mockPropertyData
		importStateID string
		stateCheck    func(s []*terraform.InstanceState) error
		error         *regexp.Regexp
	}{
		"import with all attributes": {
			mockData: basicDataBootstrap,
			init: func(p *mockProperty) {
				p.mockGetProperty()
				p.mockGetPropertyVersion()
				// read
				p.mockGetProperty()
			},
			stateCheck: baseChecker.
				CheckEqual("product_id", "").
				Build(),
			importStateID: "prp_123,2,1",
		},
		"import with only property_id": {
			mockData: basicDataWithoutContractAndGroup,
			init: func(p *mockProperty) {
				// import
				p.mockGetProperty()
				p.mockGetPropertyVersion()
				// read
				p.mockGetProperty()
			},
			stateCheck: baseChecker.
				CheckEqual("group_id", "").
				CheckEqual("contract_id", "").
				CheckEqual("product_id", "").
				Build(),
			importStateID: "123",
		},
		"import with only property_id and contract_id - error": {
			importStateID: "123,2",
			error:         regexp.MustCompile("Error: missing group id or contract id"),
		},
	}

	for name, test := range tests {
		t.Run(name, func(t *testing.T) {
			t.Parallel()

			m := &papi.Mock{}
			mp := &mockProperty{
				mockPropertyData: test.mockData,
				papiMock:         m,
			}

			if test.init != nil {
				test.init(mp)
			}

			useClient(m, nil, func() {
				resource.UnitTest(t, resource.TestCase{
					ProtoV6ProviderFactories: testutils.NewProtoV6ProviderFactory(NewSubprovider()),
					IsUnitTest:               true,
					Steps: []resource.TestStep{
						{
							ImportState:      true,
							ImportStateId:    test.importStateID,
							ImportStateCheck: test.stateCheck,
							ResourceName:     "akamai_property_bootstrap.test",
							Config:           testutils.LoadFixtureString(t, "testdata/TestResPropertyBootstrap/create.tf"),
							ExpectError:      test.error,
						},
					},
				})
			})

			m.AssertExpectations(t)
		})
	}
}
