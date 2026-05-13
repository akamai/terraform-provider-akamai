package mtlstruststore

import (
	"fmt"
	"regexp"
	"testing"

	"github.com/akamai/AkamaiOPEN-edgegrid-golang/v13/pkg/mtlstruststore"
	tst "github.com/akamai/terraform-provider-akamai/v10/internal/test"
	"github.com/akamai/terraform-provider-akamai/v10/pkg/common/ptr"
	"github.com/akamai/terraform-provider-akamai/v10/pkg/common/test"
	"github.com/akamai/terraform-provider-akamai/v10/pkg/common/testutils"
	"github.com/hashicorp/terraform-plugin-testing/helper/resource"
)

func TestCASetsDataSource(t *testing.T) {
	t.Parallel()
	stagingNetworkStateChecker := test.NewStateChecker("data.akamai_mtlstruststore_ca_sets.test").
		CheckEqual("ca_sets.0.name", "staging_ca_set").
		CheckEqual("ca_sets.0.id", "111111").
		CheckEqual("ca_sets.0.status", "active").
		CheckEqual("ca_sets.0.latest_version", "1").
		CheckEqual("ca_sets.0.staging_version", "1").
		CheckEqual("ca_sets.0.created_by", "user1").
		CheckEqual("ca_sets.0.created_date", "2023-01-01T00:00:00Z").
		CheckEqual("ca_sets.0.deleted_by", "user1").
		CheckEqual("ca_sets.0.deleted_date", "2023-01-03T00:00:00Z").
		CheckEqual("ca_sets.0.account_id", "test_account_1").
		CheckEqual("ca_sets.0.description", "Test CA Set Only Staging Description").
		CheckMissing("ca_sets.0.removal_date")

	productionNetworkStateChecker := test.NewStateChecker("data.akamai_mtlstruststore_ca_sets.test").
		CheckEqual("activated_on", "production").
		CheckEqual("ca_sets.0.name", "production_ca_set").
		CheckEqual("ca_sets.0.id", "222222").
		CheckEqual("ca_sets.0.status", "active").
		CheckEqual("ca_sets.0.latest_version", "1").
		CheckEqual("ca_sets.0.production_version", "1").
		CheckEqual("ca_sets.0.created_by", "user1").
		CheckEqual("ca_sets.0.created_date", "2023-01-01T00:00:00Z").
		CheckEqual("ca_sets.0.deleted_by", "user3").
		CheckEqual("ca_sets.0.deleted_date", "2023-01-03T00:00:00Z").
		CheckEqual("ca_sets.0.account_id", "test_account_1").
		CheckEqual("ca_sets.0.description", "Test CA Set Only Production Description").
		CheckMissing("ca_sets.0.removal_date")

	bothNetworksStateChecker := test.NewStateChecker("data.akamai_mtlstruststore_ca_sets.test").
		CheckEqual("ca_sets.0.name", "both_ca_set").
		CheckEqual("ca_sets.0.id", "333333").
		CheckEqual("ca_sets.0.status", "active").
		CheckEqual("ca_sets.0.latest_version", "1").
		CheckEqual("ca_sets.0.staging_version", "1").
		CheckEqual("ca_sets.0.production_version", "1").
		CheckEqual("ca_sets.0.created_by", "user1").
		CheckEqual("ca_sets.0.created_date", "2023-01-01T00:00:00Z").
		CheckEqual("ca_sets.0.deleted_by", "user3").
		CheckEqual("ca_sets.0.deleted_date", "2023-01-03T00:00:00Z").
		CheckEqual("ca_sets.0.account_id", "test_account_1").
		CheckEqual("ca_sets.0.description", "Test CA Set Both Description").
		CheckMissing("ca_sets.0.removal_date")

	allCASetsStateChecker := stagingNetworkStateChecker.
		CheckEqual("ca_sets.#", "3").
		CheckEqual("ca_sets.1.name", "production_ca_set").
		CheckEqual("ca_sets.1.id", "222222").
		CheckEqual("ca_sets.1.status", "active").
		CheckEqual("ca_sets.1.latest_version", "1").
		CheckEqual("ca_sets.1.production_version", "1").
		CheckEqual("ca_sets.1.created_by", "user1").
		CheckEqual("ca_sets.1.created_date", "2023-01-01T00:00:00Z").
		CheckEqual("ca_sets.1.deleted_by", "user3").
		CheckEqual("ca_sets.1.deleted_date", "2023-01-03T00:00:00Z").
		CheckEqual("ca_sets.1.account_id", "test_account_1").
		CheckEqual("ca_sets.1.description", "Test CA Set Only Production Description").
		CheckEqual("ca_sets.2.name", "both_ca_set").
		CheckEqual("ca_sets.2.id", "333333").
		CheckEqual("ca_sets.2.status", "active").
		CheckEqual("ca_sets.2.latest_version", "1").
		CheckEqual("ca_sets.2.staging_version", "1").
		CheckEqual("ca_sets.2.production_version", "1").
		CheckEqual("ca_sets.2.created_by", "user1").
		CheckEqual("ca_sets.2.created_date", "2023-01-01T00:00:00Z").
		CheckEqual("ca_sets.2.deleted_by", "user3").
		CheckEqual("ca_sets.2.deleted_date", "2023-01-03T00:00:00Z").
		CheckEqual("ca_sets.2.account_id", "test_account_1").
		CheckEqual("ca_sets.2.description", "Test CA Set Both Description").
		CheckMissing("ca_sets.1.removal_date").
		CheckMissing("ca_sets.2.removal_date")

	notDeleted1StateChecker := test.NewStateChecker("data.akamai_mtlstruststore_ca_sets.test").
		CheckEqual("ca_sets.0.id", "444444").
		CheckEqual("ca_sets.0.name", "not_deleted_ca_set_1").
		CheckEqual("ca_sets.0.status", mtlstruststore.CASetStatusNotDeleted).
		CheckEqual("ca_sets.0.latest_version", "2").
		CheckEqual("ca_sets.0.staging_version", "2").
		CheckEqual("ca_sets.0.created_by", "user1").
		CheckEqual("ca_sets.0.created_date", "2024-03-01T10:00:00Z").
		CheckEqual("ca_sets.0.account_id", "test_account_1").
		CheckEqual("ca_sets.0.description", "First active CA set").
		CheckMissing("ca_sets.0.removal_date")

	filterByNotDeletedStateChecker := notDeleted1StateChecker.
		CheckEqual("ca_sets.#", "2").
		CheckEqual("ca_sets.1.id", "555555").
		CheckEqual("ca_sets.1.name", "not_deleted_ca_set_2").
		CheckEqual("ca_sets.1.status", mtlstruststore.CASetStatusNotDeleted).
		CheckEqual("ca_sets.1.latest_version", "3").
		CheckEqual("ca_sets.1.production_version", "3").
		CheckEqual("ca_sets.1.created_by", "user2").
		CheckEqual("ca_sets.1.created_date", "2024-05-15T14:30:00Z").
		CheckEqual("ca_sets.1.account_id", "test_account_2").
		CheckEqual("ca_sets.1.description", "Second active CA set").
		CheckMissing("ca_sets.1.removal_date")

	filterByDeletedAndNotDeletedStateChecker := notDeleted1StateChecker.
		CheckEqual("ca_sets.#", "2").
		CheckEqual("ca_sets.1.id", "666666").
		CheckEqual("ca_sets.1.name", "deleted_ca_set").
		CheckEqual("ca_sets.1.status", mtlstruststore.CASetStatusDeleted).
		CheckEqual("ca_sets.1.created_by", "user1").
		CheckEqual("ca_sets.1.created_date", "2023-06-01T08:00:00Z").
		CheckEqual("ca_sets.1.deleted_by", "user2").
		CheckEqual("ca_sets.1.deleted_date", "2024-01-15T12:00:00Z").
		CheckEqual("ca_sets.1.account_id", "test_account_1").
		CheckEqual("ca_sets.1.description", "Deleted CA set").
		CheckEqual("ca_sets.1.removal_date", "2024-02-15T12:00:00Z")

	tests := map[string]struct {
		init  func(*mtlstruststore.Mock)
		steps []resource.TestStep
		error *regexp.Regexp
	}{
		"happy path - all CA sets": {
			init: func(m *mtlstruststore.Mock) {
				var resp mtlstruststore.ListCASetsResponse
				resp.CASets = append(resp.CASets, mockCASetStagingNetworkModel, mockCASetProductionNetworkModel, mockCASetBothNetworkModel)
				mockListCASetsWithParams(m, listCASetsParams{resp: &resp})
			},
			steps: []resource.TestStep{
				{
					Config: testutils.LoadFixtureString(t, "testdata/TestDataCASets/all_ca_sets.tf"),
					Check:  allCASetsStateChecker.Build(),
				},
			},
		},
		"happy path - activated on STAGING": {
			init: func(m *mtlstruststore.Mock) {
				var resp mtlstruststore.ListCASetsResponse
				resp.CASets = append(resp.CASets, mockCASetStagingNetworkModel)
				mockListCASetsWithParams(m, listCASetsParams{activatedOn: ptr.To("STAGING"), resp: &resp})
			},
			steps: []resource.TestStep{
				{
					Config: testutils.LoadFixtureString(t, "testdata/TestDataCASets/activated_on_staging_network.tf"),
					Check:  stagingNetworkStateChecker.CheckEqual("activated_on", "STAGING").Build(),
				},
			},
		},
		"happy path - activated on PRODUCTION": {
			init: func(m *mtlstruststore.Mock) {
				var resp mtlstruststore.ListCASetsResponse
				resp.CASets = append(resp.CASets, mockCASetProductionNetworkModel)
				mockListCASetsWithParams(m, listCASetsParams{activatedOn: ptr.To("PRODUCTION"), resp: &resp})
			},
			steps: []resource.TestStep{
				{
					Config: testutils.LoadFixtureString(t, "testdata/TestDataCASets/activated_on_production_network.tf"),
					Check:  productionNetworkStateChecker.CheckEqual("activated_on", "PRODUCTION").Build(),
				},
			},
		},
		"happy path - activated on BOTH": {
			init: func(m *mtlstruststore.Mock) {
				var resp mtlstruststore.ListCASetsResponse
				resp.CASets = append(resp.CASets, mockCASetBothNetworkModel)
				mockListCASetsWithParams(m, listCASetsParams{activatedOn: ptr.To("STAGING+PRODUCTION"), resp: &resp})
			},
			steps: []resource.TestStep{
				{
					Config: testutils.LoadFixtureString(t, "testdata/TestDataCASets/activated_on_both_networks.tf"),
					Check:  bothNetworksStateChecker.CheckEqual("activated_on", "STAGING+PRODUCTION").Build(),
				},
			},
		},
		"happy path - filtered by name prefix": {
			init: func(m *mtlstruststore.Mock) {
				var resp mtlstruststore.ListCASetsResponse
				resp.CASets = append(resp.CASets, mockCASetStagingNetworkModel)
				mockListCASetsWithParams(m, listCASetsParams{namePrefix: ptr.To("pref"), resp: &resp})
			},
			steps: []resource.TestStep{
				{
					Config: testutils.LoadFixtureString(t, "testdata/TestDataCASets/filtered_by_name_prefix.tf"),
					Check:  stagingNetworkStateChecker.CheckEqual("name_prefix", "pref").Build(),
				},
			},
		},
		"happy path - filter by single CA set status": {
			init: func(m *mtlstruststore.Mock) {
				var resp mtlstruststore.ListCASetsResponse
				resp.CASets = append(resp.CASets, mockCASetNotDeleted...)
				mockListCASetsWithParams(m, listCASetsParams{caSetStatuses: []string{mtlstruststore.CASetStatusNotDeleted}, resp: &resp})
			},
			steps: []resource.TestStep{
				{
					Config: testutils.LoadFixtureString(t, "testdata/TestDataCASets/filter_by_not_deleted.tf"),
					Check:  filterByNotDeletedStateChecker.Build(),
				},
			},
		},
		"happy path - filter by multiple CA set statuses": {
			init: func(m *mtlstruststore.Mock) {
				var resp mtlstruststore.ListCASetsResponse
				resp.CASets = append(resp.CASets, mockCASetNotDeleted[0], mockCASetDeleted)
				mockListCASetsWithParams(m, listCASetsParams{caSetStatuses: []string{mtlstruststore.CASetStatusDeleted, mtlstruststore.CASetStatusNotDeleted}, resp: &resp})
			},
			steps: []resource.TestStep{
				{
					Config: testutils.LoadFixtureString(t, "testdata/TestDataCASets/filter_by_deleted_and_not_deleted.tf"),
					Check:  filterByDeletedAndNotDeletedStateChecker.Build(),
				},
			},
		},
		"error - API error": {
			init: func(m *mtlstruststore.Mock) {
				mockListCASetsWithParams(m, listCASetsParams{hasError: true})
			},
			steps: []resource.TestStep{
				{
					Config:      testutils.LoadFixtureString(t, "testdata/TestDataCASets/all_ca_sets.tf"),
					ExpectError: regexp.MustCompile("oops"),
				},
			},
		},
		"error - bad network specifier": {
			steps: []resource.TestStep{
				{
					Config:      testutils.LoadFixtureString(t, "testdata/TestDataCASets/bad_network_specifier.tf"),
					ExpectError: regexp.MustCompile("Attribute activated_on value must be one of"),
				},
			},
		},
		"error - invalid ca_set_statuses value": {
			steps: []resource.TestStep{
				{
					Config:      testutils.LoadFixtureString(t, "testdata/TestDataCASets/invalid_ca_set_status.tf"),
					ExpectError: regexp.MustCompile(`(?s)Attribute ca_set_statuses\[Value\("INVALID"\)\] value must be one of:.+\["NOT_DELETED" "DELETED" "DELETING"\]`),
				},
			},
		},
		"error - empty ca_set_statuses": {
			steps: []resource.TestStep{
				{
					Config:      testutils.LoadFixtureString(t, "testdata/TestDataCASets/empty_ca_set_statuses.tf"),
					ExpectError: regexp.MustCompile("Attribute ca_set_statuses set must contain at least 1 elements, got: 0"),
				},
			},
		},
		"happy path - INACTIVE with DELETED ca_set_statuses": {
			init: func(m *mtlstruststore.Mock) {
				var resp mtlstruststore.ListCASetsResponse
				resp.CASets = append(resp.CASets, mockCASetDeleted)
				mockListCASetsWithParams(m, listCASetsParams{activatedOn: ptr.To("INACTIVE"), caSetStatuses: []string{mtlstruststore.CASetStatusDeleted}, resp: &resp})
			},
			steps: []resource.TestStep{
				{
					Config: testutils.LoadFixtureString(t, "testdata/TestDataCASets/inactive_with_deleted_status.tf"),
					Check:  resource.TestCheckResourceAttr("data.akamai_mtlstruststore_ca_sets.test", "ca_sets.0.id", "666666"),
				},
			},
		},
		"happy path - INACTIVE with DELETING ca_set_statuses": {
			init: func(m *mtlstruststore.Mock) {
				var resp mtlstruststore.ListCASetsResponse
				resp.CASets = append(resp.CASets, mockCASetDeleting)
				mockListCASetsWithParams(m, listCASetsParams{activatedOn: ptr.To("INACTIVE"), caSetStatuses: []string{mtlstruststore.CASetStatusDeleting}, resp: &resp})
			},
			steps: []resource.TestStep{
				{
					Config: testutils.LoadFixtureString(t, "testdata/TestDataCASets/inactive_with_deleting_status.tf"),
					Check:  resource.TestCheckResourceAttr("data.akamai_mtlstruststore_ca_sets.test", "ca_sets.0.id", "777777"),
				},
			},
		},
		"error - activated_on with DELETED ca_set_statuses": {
			steps: []resource.TestStep{
				{
					Config:      testutils.LoadFixtureString(t, "testdata/TestDataCASets/activated_on_with_deleted_status.tf"),
					ExpectError: regexp.MustCompile("(?s)Attribute `activated_on` cannot be used when `ca_set_statuses` attribute.+includes `DELETED` or `DELETING` statuses, unless it is set to `INACTIVE`."),
				},
			},
		},
		"error - activated_on with DELETING ca_set_statuses": {
			steps: []resource.TestStep{
				{
					Config:      testutils.LoadFixtureString(t, "testdata/TestDataCASets/activated_on_with_deleting_status.tf"),
					ExpectError: regexp.MustCompile("(?s)Attribute `activated_on` cannot be used when `ca_set_statuses` attribute.+includes `DELETED` or `DELETING` statuses, unless it is set to `INACTIVE`."),
				},
			},
		},
	}
	for name, tc := range tests {
		t.Run(name, func(t *testing.T) {
			t.Parallel()
			client := &mtlstruststore.Mock{}
			if tc.init != nil {
				tc.init(client)
			}
			useClient(client, func() {
				resource.UnitTest(t, resource.TestCase{
					ProtoV6ProviderFactories: testutils.NewProtoV6ProviderFactory(NewSubprovider()),
					IsUnitTest:               true,
					Steps:                    tc.steps,
				})
			})
			client.AssertExpectations(t)
		})
	}
}

type listCASetsParams struct {
	activatedOn   *string
	namePrefix    *string
	caSetStatuses []string
	resp          *mtlstruststore.ListCASetsResponse
	hasError      bool
}

func mockListCASetsWithParams(m *mtlstruststore.Mock, params listCASetsParams) {
	req := mtlstruststore.ListCASetsRequest{}
	if params.activatedOn != nil {
		req.ActivatedOn = mtlstruststore.Network(*params.activatedOn)
	}
	if params.namePrefix != nil {
		req.CASetNamePrefix = *params.namePrefix
	}
	if params.caSetStatuses != nil {
		req.CASetStatuses = params.caSetStatuses
	}
	if params.hasError {
		m.On("ListCASets", testutils.MockContext, req).Return(nil, fmt.Errorf("oops")).Once()
		return
	}
	m.On("ListCASets", testutils.MockContext, req).Return(params.resp, nil).Times(3)
}

var (
	mockCASetStagingNetworkModel = mtlstruststore.CASetResponse{
		CASetID:        "111111",
		CASetName:      "staging_ca_set",
		CASetStatus:    "active",
		LatestVersion:  ptr.To[int64](1),
		StagingVersion: ptr.To[int64](1),
		CreatedBy:      "user1",
		CreatedDate:    tst.NewTimeFromStringMust("2023-01-01T00:00:00Z"),
		DeletedBy:      ptr.To("user1"),
		DeletedDate:    ptr.To(tst.NewTimeFromStringMust("2023-01-03T00:00:00Z")),
		AccountID:      "test_account_1",
		Description:    ptr.To("Test CA Set Only Staging Description"),
	}

	mockCASetProductionNetworkModel = mtlstruststore.CASetResponse{
		CASetID:           "222222",
		CASetName:         "production_ca_set",
		CASetStatus:       "active",
		LatestVersion:     ptr.To[int64](1),
		ProductionVersion: ptr.To[int64](1),
		CreatedBy:         "user1",
		CreatedDate:       tst.NewTimeFromStringMust("2023-01-01T00:00:00Z"),
		DeletedBy:         ptr.To("user3"),
		DeletedDate:       ptr.To(tst.NewTimeFromStringMust("2023-01-03T00:00:00Z")),
		AccountID:         "test_account_1",
		Description:       ptr.To("Test CA Set Only Production Description"),
	}

	mockCASetBothNetworkModel = mtlstruststore.CASetResponse{
		CASetID:           "333333",
		CASetName:         "both_ca_set",
		CASetStatus:       "active",
		LatestVersion:     ptr.To[int64](1),
		StagingVersion:    ptr.To[int64](1),
		ProductionVersion: ptr.To[int64](1),
		CreatedBy:         "user1",
		CreatedDate:       tst.NewTimeFromStringMust("2023-01-01T00:00:00Z"),
		DeletedBy:         ptr.To("user3"),
		DeletedDate:       ptr.To(tst.NewTimeFromStringMust("2023-01-03T00:00:00Z")),
		AccountID:         "test_account_1",
		Description:       ptr.To("Test CA Set Both Description"),
	}

	mockCASetNotDeleted = []mtlstruststore.CASetResponse{
		{
			CASetID:        "444444",
			CASetName:      "not_deleted_ca_set_1",
			CASetStatus:    mtlstruststore.CASetStatusNotDeleted,
			LatestVersion:  ptr.To[int64](2),
			StagingVersion: ptr.To[int64](2),
			CreatedBy:      "user1",
			CreatedDate:    tst.NewTimeFromStringMust("2024-03-01T10:00:00Z"),
			AccountID:      "test_account_1",
			Description:    ptr.To("First active CA set"),
		},
		{
			CASetID:           "555555",
			CASetName:         "not_deleted_ca_set_2",
			CASetStatus:       mtlstruststore.CASetStatusNotDeleted,
			LatestVersion:     ptr.To[int64](3),
			ProductionVersion: ptr.To[int64](3),
			CreatedBy:         "user2",
			CreatedDate:       tst.NewTimeFromStringMust("2024-05-15T14:30:00Z"),
			AccountID:         "test_account_2",
			Description:       ptr.To("Second active CA set"),
		},
	}

	mockCASetDeleted = mtlstruststore.CASetResponse{
		CASetID:     "666666",
		CASetName:   "deleted_ca_set",
		CASetStatus: mtlstruststore.CASetStatusDeleted,
		CreatedBy:   "user1",
		CreatedDate: tst.NewTimeFromStringMust("2023-06-01T08:00:00Z"),
		DeletedBy:   ptr.To("user2"),
		DeletedDate: ptr.To(tst.NewTimeFromStringMust("2024-01-15T12:00:00Z")),
		RemovalDate: ptr.To(tst.NewTimeFromStringMust("2024-02-15T12:00:00Z")),
		AccountID:   "test_account_1",
		Description: ptr.To("Deleted CA set"),
	}

	mockCASetDeleting = mtlstruststore.CASetResponse{
		CASetID:     "777777",
		CASetName:   "deleting_ca_set",
		CASetStatus: mtlstruststore.CASetStatusDeleting,
		CreatedBy:   "user1",
		CreatedDate: tst.NewTimeFromStringMust("2023-07-01T08:00:00Z"),
		AccountID:   "test_account_1",
		Description: ptr.To("Deleting CA set"),
	}
)
