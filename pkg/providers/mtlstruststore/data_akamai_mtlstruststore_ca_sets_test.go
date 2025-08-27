package mtlstruststore

import (
	"fmt"
	"regexp"
	"testing"

	"github.com/akamai/AkamaiOPEN-edgegrid-golang/v12/pkg/mtlstruststore"
	tst "github.com/akamai/terraform-provider-akamai/v8/internal/test"
	"github.com/akamai/terraform-provider-akamai/v8/pkg/common/ptr"
	"github.com/akamai/terraform-provider-akamai/v8/pkg/common/test"
	"github.com/akamai/terraform-provider-akamai/v8/pkg/common/testutils"
	"github.com/hashicorp/terraform-plugin-testing/helper/resource"
)

func TestCASetsDataSource(t *testing.T) {
	testDir := "testdata/TestDataCASets/"
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
		CheckEqual("ca_sets.0.description", "Test CA Set Only Staging Description")

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
		CheckEqual("ca_sets.0.description", "Test CA Set Only Production Description")

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
		CheckEqual("ca_sets.0.description", "Test CA Set Both Description")

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
		CheckEqual("ca_sets.2.description", "Test CA Set Both Description")

	tests := map[string]struct {
		init  func(*mtlstruststore.Mock)
		steps []resource.TestStep
		error *regexp.Regexp
	}{
		"happy path - all CA sets": {
			init: func(m *mtlstruststore.Mock) {
				var resp mtlstruststore.ListCASetsResponse
				resp.CASets = append(resp.CASets, mockCASetStagingNetworkModel, mockCASetProductionNetworkModel, mockCASetBothNetworkModel)
				mockListCASetsWithParams(m, nil, nil, &resp, false)
			},
			steps: []resource.TestStep{
				{
					Config: testutils.LoadFixtureString(t, testDir+"all_ca_sets.tf"),
					Check:  allCASetsStateChecker.Build(),
				},
			},
		},
		"happy path - activated on STAGING": {
			init: func(m *mtlstruststore.Mock) {
				var resp mtlstruststore.ListCASetsResponse
				resp.CASets = append(resp.CASets, mockCASetStagingNetworkModel)
				mockListCASetsWithParams(m, ptr.To("STAGING"), nil, &resp, false)
			},
			steps: []resource.TestStep{
				{
					Config: testutils.LoadFixtureString(t, testDir+"activated_on_staging_network.tf"),
					Check:  stagingNetworkStateChecker.CheckEqual("activated_on", "STAGING").Build(),
				},
			},
		},
		"happy path - activated on PRODUCTION": {
			init: func(m *mtlstruststore.Mock) {
				var resp mtlstruststore.ListCASetsResponse
				resp.CASets = append(resp.CASets, mockCASetProductionNetworkModel)
				mockListCASetsWithParams(m, ptr.To("PRODUCTION"), nil, &resp, false)
			},
			steps: []resource.TestStep{
				{
					Config: testutils.LoadFixtureString(t, testDir+"activated_on_production_network.tf"),
					Check:  productionNetworkStateChecker.CheckEqual("activated_on", "PRODUCTION").Build(),
				},
			},
		},
		"happy path - activated on BOTH": {
			init: func(m *mtlstruststore.Mock) {
				var resp mtlstruststore.ListCASetsResponse
				resp.CASets = append(resp.CASets, mockCASetBothNetworkModel)
				mockListCASetsWithParams(m, ptr.To("STAGING+PRODUCTION"), nil, &resp, false)
			},
			steps: []resource.TestStep{
				{
					Config: testutils.LoadFixtureString(t, testDir+"activated_on_both_networks.tf"),
					Check:  bothNetworksStateChecker.CheckEqual("activated_on", "STAGING+PRODUCTION").Build(),
				},
			},
		},
		"happy path - filtered by name prefix": {
			init: func(m *mtlstruststore.Mock) {
				var resp mtlstruststore.ListCASetsResponse
				resp.CASets = append(resp.CASets, mockCASetStagingNetworkModel)
				mockListCASetsWithParams(m, nil, ptr.To("pref"), &resp, false)
			},
			steps: []resource.TestStep{
				{
					Config: testutils.LoadFixtureString(t, testDir+"filtered_by_name_prefix.tf"),
					Check:  stagingNetworkStateChecker.CheckEqual("name_prefix", "pref").Build(),
				},
			},
		},
		"error - API error": {
			init: func(m *mtlstruststore.Mock) {
				mockListCASetsWithParams(m, nil, nil, nil, true)
			},
			steps: []resource.TestStep{
				{
					Config:      testutils.LoadFixtureString(t, testDir+"all_ca_sets.tf"),
					ExpectError: regexp.MustCompile("oops"),
				},
			},
		},
		"error - bad network specifier": {
			steps: []resource.TestStep{
				{
					Config:      testutils.LoadFixtureString(t, testDir+"bad_network_specifier.tf"),
					ExpectError: regexp.MustCompile("Attribute activated_on value must be one of"),
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

func mockListCASetsWithParams(m *mtlstruststore.Mock, activatedOn, namePrefix *string, resp *mtlstruststore.ListCASetsResponse, hasError bool) {
	req := mtlstruststore.ListCASetsRequest{}
	if activatedOn != nil {
		req.ActivatedOn = mtlstruststore.Network(*activatedOn)
	}
	if namePrefix != nil {
		req.CASetNamePrefix = *namePrefix
	}
	if hasError {
		m.On("ListCASets", testutils.MockContext, req).Return(nil, fmt.Errorf("oops")).Once()
		return
	}
	m.On("ListCASets", testutils.MockContext, req).Return(resp, nil).Times(3)
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
)
