package property

import (
	"regexp"
	"testing"

	"github.com/akamai/AkamaiOPEN-edgegrid-golang/v8/pkg/papi"
	"github.com/akamai/terraform-provider-akamai/v6/pkg/common/testutils"
	"github.com/hashicorp/terraform-plugin-testing/helper/resource"
)

func Test_DSReadContract(t *testing.T) {
	tests := map[string]struct {
		init       func(*testing.T, *papi.Mock, testDataForPAPIGroups)
		mockData   testDataForPAPIGroups
		configPath string
		error      *regexp.Regexp
	}{
		"read contract with group name and group ID conflict": {
			init:       func(t *testing.T, m *papi.Mock, testData testDataForPAPIGroups) {},
			configPath: "testdata/TestDSContractRequired/ds_contract_with_group_name_and_group.tf",
			error:      regexp.MustCompile("only one of `group_id,group_name` can be specified"),
		},
		"read contract with group id provided": {
			init: func(t *testing.T, m *papi.Mock, testData testDataForPAPIGroups) {
				expectGetGroups(m, testData, 5)
			},
			mockData: testDataForPAPIGroups{
				accountID:   "",
				accountName: "",
				groups: papi.GroupItems{
					Items: []*papi.Group{
						{
							GroupID:       "grp_12345",
							GroupName:     "Example.com-1-1TJZH5",
							ParentGroupID: "grp_parent",
							ContractIDs:   []string{"ctr_1234"},
						},
						{
							GroupID:       "grp_12346",
							GroupName:     "default",
							ParentGroupID: "grp_parent",
							ContractIDs:   []string{"ctr_1234"},
						},
					},
				},
			},
			configPath: "testdata/TestDSContractRequired/ds_contract_with_group_id.tf",
			error:      nil,
		},
		"read contract with group id without prefix": {
			init: func(t *testing.T, m *papi.Mock, testData testDataForPAPIGroups) {
				expectGetGroups(m, testData, 5)
			},
			mockData: testDataForPAPIGroups{
				accountID:   "act_1-1TJZFB",
				accountName: "example.com",
				groups: papi.GroupItems{
					Items: []*papi.Group{
						{
							GroupID:       "grp_12345",
							GroupName:     "Example.com-1-1TJZH5",
							ParentGroupID: "grp_parent",
							ContractIDs:   []string{"ctr_1234"},
						},
						{
							GroupID:       "grp_12346",
							GroupName:     "default",
							ParentGroupID: "grp_parent",
							ContractIDs:   []string{"ctr_1234"},
						},
					},
				},
			},
			configPath: "testdata/TestDSContractRequired/ds_contract_with_group_id_without_prefix.tf",
			error:      nil,
		},
		"read contract with group name": {
			init: func(t *testing.T, m *papi.Mock, testData testDataForPAPIGroups) {
				expectGetGroups(m, testData, 5)
			},
			mockData: testDataForPAPIGroups{
				accountID:   "act_1-1TJZFB",
				accountName: "example.com",
				groups: papi.GroupItems{
					Items: []*papi.Group{
						{
							GroupID:       "grp_12345",
							GroupName:     "Example.com-1-1TJZH5",
							ParentGroupID: "grp_parent",
							ContractIDs:   []string{"ctr_1234"},
						},
						{
							GroupID:       "grp_12346",
							GroupName:     "default",
							ParentGroupID: "grp_parent",
							ContractIDs:   []string{"ctr_1234"},
						},
					},
				},
			},
			configPath: "testdata/TestDSContractRequired/ds_contract_with_group_name.tf",
			error:      nil,
		},
		"multiple groups with the same name - expect an error": {
			init: func(t *testing.T, m *papi.Mock, testData testDataForPAPIGroups) {
				expectGetGroups(m, testData, 5)
			},
			mockData: testDataForPAPIGroups{
				accountID:   "act_1-1TJZFB",
				accountName: "example.com",
				groups: papi.GroupItems{
					Items: []*papi.Group{
						{
							GroupID:       "grp_12345",
							GroupName:     "Example.com-1-1TJZH5",
							ParentGroupID: "grp_parent",
							ContractIDs:   []string{"ctr_1234"},
						},
						{
							GroupID:       "grp_12346",
							GroupName:     "Example.com-1-1TJZH5",
							ParentGroupID: "grp_parent",
							ContractIDs:   []string{"ctr_1234"},
						},
					},
				},
			},
			configPath: "testdata/TestDSContractRequired/ds_contract_with_group_name.tf",
			error:      regexp.MustCompile("there is more than 1 group with the same name. Based on provided data, it is impossible to determine which one should be returned. Please use group_id attribute"),
		},
		"multiple groups with the same name, distinguished by group_id": {
			init: func(t *testing.T, m *papi.Mock, testData testDataForPAPIGroups) {
				expectGetGroups(m, testData, 5)
			},
			mockData: testDataForPAPIGroups{
				accountID:   "act_1-1TJZFB",
				accountName: "example.com",
				groups: papi.GroupItems{
					Items: []*papi.Group{
						{
							GroupID:       "grp_12345",
							GroupName:     "Example.com-1-1TJZH5",
							ParentGroupID: "grp_parent",
							ContractIDs:   []string{"ctr_1234"},
						},
						{
							GroupID:       "grp_12346",
							GroupName:     "Example.com-1-1TJZH5",
							ParentGroupID: "grp_parent",
							ContractIDs:   []string{"ctr_1234"},
						},
					},
				},
			},
			configPath: "testdata/TestDSContractRequired/ds_contract_with_group_id.tf",
			error:      nil,
		},
		"group with multiple contracts - expect error": {
			init: func(t *testing.T, m *papi.Mock, testData testDataForPAPIGroups) {
				expectGetGroups(m, testData, 5)
			},
			mockData: testDataForPAPIGroups{
				accountID:   "act_1-1TJZFB",
				accountName: "example.com",
				groups: papi.GroupItems{
					Items: []*papi.Group{
						{
							GroupID:       "grp_12345",
							GroupName:     "Example.com-1-1TJZH5",
							ParentGroupID: "grp_parent",
							ContractIDs:   []string{"ctr_1234", "ctr_1235"},
						},
						{
							GroupID:       "grp_12346",
							GroupName:     "Second-Example.com-1-1TJZH5",
							ParentGroupID: "grp_parent",
							ContractIDs:   []string{"ctr_1234"},
						},
					},
				},
			},
			configPath: "testdata/TestDSContractRequired/ds_contract_with_group_id.tf",
			error:      regexp.MustCompile("multiple contracts found for given group"),
		},
	}

	for name, test := range tests {
		t.Run(name, func(t *testing.T) {
			client := &papi.Mock{}
			test.init(t, client, test.mockData)
			useClient(client, nil, func() {
				resource.UnitTest(t, resource.TestCase{
					ProtoV6ProviderFactories: testutils.NewProtoV6ProviderFactory(NewSubprovider()),
					IsUnitTest:               true,
					Steps: []resource.TestStep{
						{
							Config: testutils.LoadFixtureString(t, test.configPath),
							Check: resource.ComposeAggregateTestCheckFunc(
								resource.TestCheckResourceAttr("data.akamai_contract.akacontract", "id", "ctr_1234"),
								resource.TestCheckResourceAttr("data.akamai_contract.akacontract", "group_id", "grp_12345"),
								resource.TestCheckResourceAttr("data.akamai_contract.akacontract", "group_name", "Example.com-1-1TJZH5"),
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
