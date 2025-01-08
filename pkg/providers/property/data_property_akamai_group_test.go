package property

import (
	"regexp"
	"testing"

	"github.com/akamai/AkamaiOPEN-edgegrid-golang/v9/pkg/papi"
	"github.com/akamai/terraform-provider-akamai/v6/pkg/common/testutils"
	"github.com/hashicorp/terraform-plugin-testing/helper/resource"
)

func Test_DSReadGroup(t *testing.T) {
	tests := map[string]struct {
		init       func(*papi.Mock, testDataForPAPIGroups)
		mockData   testDataForPAPIGroups
		configPath string
		error      *regexp.Regexp
	}{
		"read group with group_name and contract_id provided": {
			init: func(m *papi.Mock, testData testDataForPAPIGroups) {
				expectGetGroups(m, testData, 5)
			},
			mockData: testDataForPAPIGroups{
				accountID:   "testAccountID",
				accountName: "testAccountName",
				groups: papi.GroupItems{
					Items: []*papi.Group{
						{
							GroupID:       "grp_12345",
							GroupName:     "Example.com-1-1TJZH5",
							ParentGroupID: "grp_parent",
							ContractIDs:   []string{"ctr_1234"},
						},
					},
				},
			},
			configPath: "testdata/TestDSGroup/ds-group-w-group-name-and-contract_id.tf",
		},
		"multiple groups distinguished by contract_id": {
			init: func(m *papi.Mock, testData testDataForPAPIGroups) {
				expectGetGroups(m, testData, 5)
			},
			mockData: testDataForPAPIGroups{
				accountID:   "testAccountID",
				accountName: "testAccountName",
				groups: papi.GroupItems{
					Items: []*papi.Group{
						{
							GroupID:       "grp_12345",
							GroupName:     "Example.com-1-1TJZH5",
							ParentGroupID: "grp_parent",
							ContractIDs:   []string{"ctr_1234"},
						},
						{
							GroupID:       "grp_123456",
							GroupName:     "Example.com-1-1TJZH5",
							ParentGroupID: "grp_parent2",
							ContractIDs:   []string{"ctr_12345"},
						},
					},
				},
			},
			configPath: "testdata/TestDSGroup/ds-group-w-group-name-and-contract_id.tf",
		},
		"multiple groups with the same group names and multiple distinguishable contracts": {
			init: func(m *papi.Mock, testData testDataForPAPIGroups) {
				expectGetGroups(m, testData, 5)
			},
			mockData: testDataForPAPIGroups{
				accountID:   "testAccountID",
				accountName: "testAccountName",
				groups: papi.GroupItems{
					Items: []*papi.Group{
						{
							GroupID:       "grp_123456",
							GroupName:     "Example.com-1-1TJZH5",
							ParentGroupID: "grp_parent",
							ContractIDs:   []string{"ctr_12345", "ctr_123456"},
						},
						{
							GroupID:       "grp_12345",
							GroupName:     "Example.com-1-1TJZH5",
							ParentGroupID: "grp_parent2",
							ContractIDs:   []string{"ctr_1234", "ctr_12345"},
						},
					},
				},
			},
			configPath: "testdata/TestDSGroup/ds-group-w-group-name-and-contract_id.tf",
		},
		"multiple groups with the same group_name and contract - expect an error": {
			init: func(m *papi.Mock, testData testDataForPAPIGroups) {
				expectGetGroups(m, testData, 5)
			},
			mockData: testDataForPAPIGroups{
				accountID:   "testAccountID",
				accountName: "testAccountName",
				groups: papi.GroupItems{
					Items: []*papi.Group{
						{
							GroupID:       "grp_12345",
							GroupName:     "Example.com-1-1TJZH5",
							ParentGroupID: "grp_parent",
							ContractIDs:   []string{"ctr_1234"},
						},
						{
							GroupID:       "grp_123456",
							GroupName:     "Example.com-1-1TJZH5",
							ParentGroupID: "grp_parent2",
							ContractIDs:   []string{"ctr_1234"},
						},
					},
				},
			},
			configPath: "testdata/TestDSGroup/ds-group-w-group-name-and-contract_id.tf",
			error:      regexp.MustCompile("there is more than 1 group with the same name and contract combination. Based on provided data, it is impossible to determine which one should be returned"),
		},
		"multiple groups with the same multiple contracts and the same group names - expect an error": {
			init: func(m *papi.Mock, testData testDataForPAPIGroups) {
				expectGetGroups(m, testData, 5)
			},
			mockData: testDataForPAPIGroups{
				accountID:   "testAccountID",
				accountName: "testAccountName",
				groups: papi.GroupItems{
					Items: []*papi.Group{
						{
							GroupID:       "grp_12345",
							GroupName:     "Example.com-1-1TJZH5",
							ParentGroupID: "grp_parent",
							ContractIDs:   []string{"ctr_1234", "ctr_12345"},
						},
						{
							GroupID:       "grp_123456",
							GroupName:     "Example.com-1-1TJZH5",
							ParentGroupID: "grp_parent2",
							ContractIDs:   []string{"ctr_1234", "ctr_12345"},
						},
					},
				},
			},
			configPath: "testdata/TestDSGroup/ds-group-w-group-name-and-contract_id.tf",
			error:      regexp.MustCompile("there is more than 1 group with the same name and contract combination. Based on provided data, it is impossible to determine which one should be returned"),
		},
	}

	for name, test := range tests {
		t.Run(name, func(t *testing.T) {
			client := &papi.Mock{}
			test.init(client, test.mockData)
			useClient(client, nil, func() {
				resource.UnitTest(t, resource.TestCase{
					ProtoV6ProviderFactories: testutils.NewProtoV6ProviderFactory(NewSubprovider()),
					IsUnitTest:               true,
					Steps: []resource.TestStep{
						{
							Config: testutils.LoadFixtureString(t, test.configPath),
							Check: resource.ComposeAggregateTestCheckFunc(
								resource.TestCheckResourceAttr("data.akamai_group.akagroup", "id", "grp_12345"),
								resource.TestCheckResourceAttr("data.akamai_group.akagroup", "group_name", "Example.com-1-1TJZH5"),
								resource.TestCheckResourceAttr("data.akamai_group.akagroup", "contract_id", "ctr_1234"),
							),
							ExpectError: test.error,
						},
					},
				})
			})
			client.AssertExpectations(t)
		})
	}
}

type testDataForPAPIGroups struct {
	accountID   string
	accountName string
	groups      papi.GroupItems
}

var expectGetGroups = func(client *papi.Mock, data testDataForPAPIGroups, timesToRun int) {
	client.On("GetGroups", testutils.MockContext).Return(&papi.GetGroupsResponse{
		AccountID:   data.accountID,
		AccountName: data.accountName,
		Groups:      data.groups,
	}, nil)
}
