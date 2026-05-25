package property

import (
	"fmt"
	"regexp"
	"testing"

	"github.com/akamai/AkamaiOPEN-edgegrid-golang/v13/pkg/papi"
	"github.com/akamai/terraform-provider-akamai/v10/internal/edgegrid"
	"github.com/akamai/terraform-provider-akamai/v10/pkg/common/test"
	"github.com/akamai/terraform-provider-akamai/v10/pkg/common/testutils"
	"github.com/hashicorp/terraform-plugin-testing/helper/resource"
)

func TestDSCPCode(t *testing.T) {
	t.Parallel()

	baseChecker := test.NewStateChecker("data.akamai_cp_code.test").
		CheckEqual("id", "234").
		CheckEqual("cp_code_name", "test cpcode").
		CheckEqual("created_date", "2021-11-11T11:22:33Z").
		CheckEqual("group_id", "grp_22").
		CheckEqual("contract_id", "ctr_11")

	t.Run("match by name", func(t *testing.T) {
		t.Parallel()
		client := edgegrid.NewTestClient()

		// name provided by fixture is "test cpcode"
		cpc := papi.CPCodeItems{Items: []papi.CPCode{
			{ID: "cpc_123", Name: "wrong CP code"},
			{ID: "cpc_234", Name: "test cpcode", CreatedDate: "2021-11-11T11:22:33Z", ProductIDs: []string{"prd_test1", "prd_test2"}},
		}}

		client.PAPI.On("GetCPCodes",
			testutils.MockContext,
			papi.GetCPCodesRequest{ContractID: "ctr_11", GroupID: "grp_22"},
		).Return(&papi.GetCPCodesResponse{CPCodes: cpc}, nil)

		resource.UnitTest(t, resource.TestCase{
			ProtoV6ProviderFactories: testutils.NewTestProtoV6SDKProviderFactory(client, NewSubprovider()),
			Steps: []resource.TestStep{{
				Config: testutils.LoadFixtureString(t, "testdata/TestDSCPCode/match_by_name.tf"),
				Check:  baseChecker.Build(),
			}},
		})

		client.PAPI.AssertExpectations(t)
	})

	t.Run("match by name output products", func(t *testing.T) {
		t.Parallel()
		client := edgegrid.NewTestClient()

		// name provided by fixture is "test cpcode"
		cpc := papi.CPCodeItems{Items: []papi.CPCode{
			{ID: "cpc_123", Name: "wrong CP code"},
			{ID: "cpc_234", Name: "test cpcode", CreatedDate: "2021-11-11T11:22:33Z", ProductIDs: []string{"prd_test1", "prd_test2"}},
		}}

		client.PAPI.On("GetCPCodes",
			testutils.MockContext,
			papi.GetCPCodesRequest{ContractID: "ctr_11", GroupID: "grp_22"},
		).Return(&papi.GetCPCodesResponse{CPCodes: cpc}, nil)

		resource.UnitTest(t, resource.TestCase{
			ProtoV6ProviderFactories: testutils.NewTestProtoV6SDKProviderFactory(client, NewSubprovider()),
			Steps: []resource.TestStep{{
				Config: testutils.LoadFixtureString(t, "testdata/TestDSCPCode/match_by_name_output_products.tf"),
				Check: resource.ComposeAggregateTestCheckFunc(
					baseChecker.Build(),
					resource.TestCheckOutput("product1", "prd_test1"),
					resource.TestCheckOutput("product2", "prd_test2"),
				),
			}},
		})

		client.PAPI.AssertExpectations(t)
	})

	t.Run("match by full ID", func(t *testing.T) {
		t.Parallel()
		client := edgegrid.NewTestClient()

		// name provided by fixture is "cpc_234"
		cpc := papi.CPCodeItems{Items: []papi.CPCode{
			{ID: "cpc_234", Name: "test cpcode", CreatedDate: "2021-11-11T11:22:33Z", ProductIDs: []string{"prd_test1", "prd_test2"}},
			{ID: "cpc_123", CreatedDate: "2021-11-11T11:22:33Z", Name: "wrong CP code"},
		}}

		client.PAPI.On("GetCPCode",
			testutils.MockContext,
			papi.GetCPCodeRequest{ContractID: "ctr_11", GroupID: "grp_22", CPCodeID: "cpc_234"},
		).Return(&papi.GetCPCodesResponse{CPCodes: cpc, CPCode: cpc.Items[0]}, nil)

		resource.UnitTest(t, resource.TestCase{
			ProtoV6ProviderFactories: testutils.NewTestProtoV6SDKProviderFactory(client, NewSubprovider()),
			Steps: []resource.TestStep{{
				Config: testutils.LoadFixtureString(t, "testdata/TestDSCPCode/match_by_full_id.tf"),
				Check: baseChecker.
					CheckEqual("cp_code_name", "test cpcode").
					Build(),
			}},
		})

		client.PAPI.AssertExpectations(t)
	})

	t.Run("match by unprefixed ID", func(t *testing.T) {
		t.Parallel()
		client := edgegrid.NewTestClient()

		// name provided by fixture is "234"
		cpc := papi.CPCodeItems{Items: []papi.CPCode{
			{ID: "cpc_234", Name: "test cpcode", CreatedDate: "2021-11-11T11:22:33Z", ProductIDs: []string{"prd_test1", "prd_test2"}},
			{ID: "cpc_123", CreatedDate: "2021-11-11T11:22:33Z", Name: "wrong CP code"},
		}}

		client.PAPI.On("GetCPCode",
			testutils.MockContext,
			papi.GetCPCodeRequest{ContractID: "ctr_11", GroupID: "grp_22", CPCodeID: "234"},
		).Return(&papi.GetCPCodesResponse{CPCodes: cpc, CPCode: cpc.Items[0]}, nil)

		resource.UnitTest(t, resource.TestCase{
			ProtoV6ProviderFactories: testutils.NewTestProtoV6SDKProviderFactory(client, NewSubprovider()),
			Steps: []resource.TestStep{{
				Config: testutils.LoadFixtureString(t, "testdata/TestDSCPCode/match_by_unprefixed_id.tf"),
				Check: baseChecker.
					Build(),
			}},
		})

		client.PAPI.AssertExpectations(t)
	})

	t.Run("no match by id", func(t *testing.T) {
		t.Parallel()
		client := edgegrid.NewTestClient()

		client.PAPI.On("GetCPCode",
			testutils.MockContext,
			papi.GetCPCodeRequest{ContractID: "ctr_11", GroupID: "grp_22", CPCodeID: "234"},
		).Return(nil, fmt.Errorf("%w: CPCodeID: 234", papi.ErrNotFound))

		resource.UnitTest(t, resource.TestCase{
			ProtoV6ProviderFactories: testutils.NewTestProtoV6SDKProviderFactory(client, NewSubprovider()),
			Steps: []resource.TestStep{{
				Config:      testutils.LoadFixtureString(t, "testdata/TestDSCPCode/match_by_unprefixed_id.tf"),
				ExpectError: regexp.MustCompile(`looking up cp code by id`),
			}},
		})

		client.PAPI.AssertExpectations(t)
	})

	t.Run("no match by name", func(t *testing.T) {
		t.Parallel()
		client := edgegrid.NewTestClient()

		// name provided by fixture is "test cpcode", but no CP code with that name exists
		cpc := papi.CPCodeItems{Items: []papi.CPCode{
			{ID: "cpc_123", Name: "wrong CP code"},
			{ID: "cpc_345", Name: "also wrong CP code"},
		}}

		client.PAPI.On("GetCPCodes",
			testutils.MockContext,
			papi.GetCPCodesRequest{ContractID: "ctr_11", GroupID: "grp_22"},
		).Return(&papi.GetCPCodesResponse{CPCodes: cpc}, nil)

		resource.UnitTest(t, resource.TestCase{
			ProtoV6ProviderFactories: testutils.NewTestProtoV6SDKProviderFactory(client, NewSubprovider()),
			Steps: []resource.TestStep{{
				Config:      testutils.LoadFixtureString(t, "testdata/TestDSCPCode/match_by_name.tf"),
				ExpectError: regexp.MustCompile(`looking up cp code by name`),
			}},
		})

		client.PAPI.AssertExpectations(t)
	})

	t.Run("more than one match by name", func(t *testing.T) {
		t.Parallel()
		client := edgegrid.NewTestClient()

		// name provided by fixture is "test cpcode"
		cpc := papi.CPCodeItems{Items: []papi.CPCode{
			{ID: "cpc_123", Name: "test cpcode"},
			{ID: "cpc_234", Name: "test cpcode", CreatedDate: "2021-11-11T11:22:33Z", ProductIDs: []string{"prd_test1", "prd_test2"}},
		}}

		client.PAPI.On("GetCPCodes",
			testutils.MockContext,
			papi.GetCPCodesRequest{ContractID: "ctr_11", GroupID: "grp_22"},
		).Return(&papi.GetCPCodesResponse{CPCodes: cpc}, nil)

		resource.UnitTest(t, resource.TestCase{
			ProtoV6ProviderFactories: testutils.NewTestProtoV6SDKProviderFactory(client, NewSubprovider()),
			Steps: []resource.TestStep{{
				Config:      testutils.LoadFixtureString(t, "testdata/TestDSCPCode/match_by_name.tf"),
				ExpectError: regexp.MustCompile(`more than one CP code was found for the given name`),
			}},
		})

		client.PAPI.AssertExpectations(t)
	})

	t.Run("both name and id provided", func(t *testing.T) {
		t.Parallel()
		client := edgegrid.NewTestClient()

		resource.UnitTest(t, resource.TestCase{
			ProtoV6ProviderFactories: testutils.NewTestProtoV6SDKProviderFactory(client, NewSubprovider()),
			Steps: []resource.TestStep{{
				Config:      testutils.LoadFixtureString(t, "testdata/TestDSCPCode/match_by_name_and_id.tf"),
				ExpectError: regexp.MustCompile(`Invalid combination of arguments`),
			}},
		})
	})

	t.Run("group not found in state", func(t *testing.T) {
		t.Parallel()
		client := edgegrid.NewTestClient()

		expectGetProducts(client.PAPI, "ctr_11", []string{"prd_1", "prd_2", "prd_3"})

		client.PAPI.On("GetCPCodes",
			testutils.MockContext, papi.GetCPCodesRequest{ContractID: "ctr_11", GroupID: "grp_22"},
		).Return(&papi.GetCPCodesResponse{CPCodes: papi.CPCodeItems{Items: []papi.CPCode{{
			ID: "cpc_123", Name: "test-ft-cp-code", CreatedDate: "2021-11-11T11:22:33Z", ProductIDs: []string{"prd_3"},
		}}}}, nil).Times(4)
		client.PAPI.On("GetCPCode", testutils.MockContext, papi.GetCPCodeRequest{CPCodeID: "123", ContractID: "ctr_11", GroupID: "grp_22"}).Return(&papi.GetCPCodesResponse{CPCode: papi.CPCode{
			ID: "cpc_123", Name: "test-ft-cp-code", CreatedDate: "2021-11-11T11:22:33Z", ProductIDs: []string{"prd_3"},
		}}, nil).Times(4)

		resource.UnitTest(t, resource.TestCase{
			ProtoV6ProviderFactories: testutils.NewTestProtoV6SDKProviderFactory(client, NewSubprovider()),
			IsUnitTest:               true,
			Steps: []resource.TestStep{{
				Config: testutils.LoadFixtureString(t, "testdata/TestDSGroupNotFound/cp_code_step1.tf"),
			},
				{
					Config: testutils.LoadFixtureString(t, "testdata/TestDSGroupNotFound/cp_code.tf"),
				}},
		})
		client.PAPI.AssertExpectations(t)
	})
}
