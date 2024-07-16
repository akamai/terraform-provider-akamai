package cloudaccess

import (
	"fmt"
	"regexp"
	"strconv"
	"testing"

	"github.com/akamai/AkamaiOPEN-edgegrid-golang/v8/pkg/cloudaccess"
	"github.com/akamai/terraform-provider-akamai/v6/pkg/common/date"
	"github.com/akamai/terraform-provider-akamai/v6/pkg/common/ptr"
	"github.com/akamai/terraform-provider-akamai/v6/pkg/common/testutils"
	"github.com/hashicorp/terraform-plugin-testing/helper/resource"
	"github.com/stretchr/testify/mock"
)

type (
	testDataForKeys struct {
		keys []keyData
	}

	keyData struct {
		accessKeyUID         int64
		accessKeyName        string
		groups               []cloudaccess.Group
		authenticationMethod string
		createdTime          string
		createdBy            string
		latestVersion        int64
		networkConfiguration *cloudaccess.SecureNetwork
	}
)

func TestDataKeys(t *testing.T) {
	tests := map[string]struct {
		configPath string
		init       func(*testing.T, *cloudaccess.Mock, testDataForKeys)
		mockData   testDataForKeys
		error      *regexp.Regexp
	}{
		"happy path - multiple keys with various contents": {
			configPath: "testdata/TestDataKeys/default.tf",
			init: func(t *testing.T, m *cloudaccess.Mock, testData testDataForKeys) {
				expectFullListAccessKeys(t, m, testData, 5)
			},
			mockData: testDataForKeys{
				keys: []keyData{
					{
						accessKeyUID:  1,
						accessKeyName: "key1",
						groups: []cloudaccess.Group{
							{
								ContractIDs: []string{"ctr_1", "ctr_2"},
								GroupID:     11,
								GroupName:   ptr.To("group1"),
							},
							{
								ContractIDs: []string{"ctr_3"},
								GroupID:     22,
							},
						},
						authenticationMethod: "method1",
						createdTime:          "2023-10-10T11:14:10.088704Z",
						createdBy:            "user1",
						latestVersion:        1,
						networkConfiguration: &cloudaccess.SecureNetwork{
							AdditionalCDN:   ptr.To(cloudaccess.CDNType("CDN1")),
							SecurityNetwork: "Net1",
						},
					},
					{
						accessKeyUID:  2,
						accessKeyName: "key2",
						groups: []cloudaccess.Group{
							{
								ContractIDs: []string{"ctr_1"},
								GroupID:     11,
								GroupName:   ptr.To("group1"),
							},
						},
						authenticationMethod: "method2",
						createdTime:          "2023-10-10T11:14:10.088704Z",
						createdBy:            "user2",
						latestVersion:        2,
					},
				},
			},
		},
		"happy path - no keys": {
			configPath: "testdata/TestDataKeys/default.tf",
			init: func(t *testing.T, m *cloudaccess.Mock, testData testDataForKeys) {
				expectFullListAccessKeys(t, m, testData, 5)
			},
		},
		"expect error on list access keys": {
			configPath: "testdata/TestDataKeys/default.tf",
			init: func(t *testing.T, m *cloudaccess.Mock, _ testDataForKeys) {
				m.On("ListAccessKeys", mock.Anything, cloudaccess.ListAccessKeysRequest{}).
					Return(nil, fmt.Errorf("API error")).Once()
			},
			error: regexp.MustCompile(`API error`),
		},
	}
	for name, test := range tests {
		t.Run(name, func(t *testing.T) {
			client := &cloudaccess.Mock{}
			if test.init != nil {
				test.init(t, client, test.mockData)
			}

			useClient(client, func() {
				resource.UnitTest(t, resource.TestCase{
					ProtoV6ProviderFactories: testutils.NewProtoV6ProviderFactory(NewSubprovider()),
					IsUnitTest:               true,
					Steps: []resource.TestStep{
						{
							Config:      testutils.LoadFixtureString(t, test.configPath),
							Check:       checkKeysAttrs(test.mockData),
							ExpectError: test.error,
						},
					},
				})
			})
			client.AssertExpectations(t)
		})
	}
}

func expectFullListAccessKeys(t *testing.T, client *cloudaccess.Mock, data testDataForKeys, timesToRun int) {
	listAccessKeysReq := cloudaccess.ListAccessKeysRequest{}
	listAccessKeysRes := cloudaccess.ListAccessKeysResponse{}

	for _, key := range data.keys {
		dateTime, err := date.Parse(key.createdTime)
		if err != nil {
			t.Fatalf(err.Error())
		}
		listAccessKeysRes.AccessKeys = append(listAccessKeysRes.AccessKeys, cloudaccess.AccessKeyResponse{
			AccessKeyUID:         key.accessKeyUID,
			AccessKeyName:        key.accessKeyName,
			AuthenticationMethod: key.authenticationMethod,
			NetworkConfiguration: key.networkConfiguration,
			LatestVersion:        key.latestVersion,
			Groups:               key.groups,
			CreatedBy:            key.createdBy,
			CreatedTime:          dateTime,
		})
	}

	client.On("ListAccessKeys", mock.Anything, listAccessKeysReq).Return(&listAccessKeysRes, nil).Times(timesToRun)
}

func checkKeysAttrs(data testDataForKeys) resource.TestCheckFunc {
	var checkFuncs []resource.TestCheckFunc

	if len(data.keys) == 0 {
		checkFuncs = append(checkFuncs, resource.TestCheckNoResourceAttr("data.akamai_cloudaccess_keys.test", "access_keys"))
	}

	for i, key := range data.keys {
		checkFuncs = append(checkFuncs, resource.TestCheckResourceAttr("data.akamai_cloudaccess_keys.test", fmt.Sprintf("access_keys.%d.access_key_uid", i), strconv.FormatInt(key.accessKeyUID, 10)))
		checkFuncs = append(checkFuncs, resource.TestCheckResourceAttr("data.akamai_cloudaccess_keys.test", fmt.Sprintf("access_keys.%d.access_key_name", i), key.accessKeyName))
		checkFuncs = append(checkFuncs, resource.TestCheckResourceAttr("data.akamai_cloudaccess_keys.test", fmt.Sprintf("access_keys.%d.authentication_method", i), key.authenticationMethod))
		checkFuncs = append(checkFuncs, resource.TestCheckResourceAttr("data.akamai_cloudaccess_keys.test", fmt.Sprintf("access_keys.%d.created_by", i), key.createdBy))
		checkFuncs = append(checkFuncs, resource.TestCheckResourceAttr("data.akamai_cloudaccess_keys.test", fmt.Sprintf("access_keys.%d.created_time", i), key.createdTime))
		checkFuncs = append(checkFuncs, resource.TestCheckResourceAttr("data.akamai_cloudaccess_keys.test", fmt.Sprintf("access_keys.%d.latest_version", i), strconv.FormatInt(key.latestVersion, 10)))
		checkFuncs = append(checkFuncs, resource.TestCheckResourceAttr("data.akamai_cloudaccess_keys.test", fmt.Sprintf("access_keys.%d.groups.#", i), strconv.Itoa(len(key.groups))))
		for j, grp := range key.groups {
			checkFuncs = append(checkFuncs, resource.TestCheckResourceAttr("data.akamai_cloudaccess_keys.test", fmt.Sprintf("access_keys.%d.groups.%d.group_id", i, j), strconv.FormatInt(grp.GroupID, 10)))
			if grp.GroupName != nil {
				checkFuncs = append(checkFuncs, resource.TestCheckResourceAttr("data.akamai_cloudaccess_keys.test", fmt.Sprintf("access_keys.%d.groups.%d.group_name", i, j), *grp.GroupName))
			} else {
				checkFuncs = append(checkFuncs, resource.TestCheckNoResourceAttr("data.akamai_cloudaccess_keys.test", fmt.Sprintf("access_keys.%d.groups.%d.group_name", i, j)))
			}
			checkFuncs = append(checkFuncs, resource.TestCheckResourceAttr("data.akamai_cloudaccess_keys.test", fmt.Sprintf("access_keys.%d.groups.%d.contracts_ids.#", i, j), strconv.Itoa(len(grp.ContractIDs))))
			for k, ctr := range grp.ContractIDs {
				checkFuncs = append(checkFuncs, resource.TestCheckResourceAttr("data.akamai_cloudaccess_keys.test", fmt.Sprintf("access_keys.%d.groups.%d.contracts_ids.%d", i, j, k), ctr))
			}
		}
		if key.networkConfiguration != nil {
			checkFuncs = append(checkFuncs, resource.TestCheckResourceAttr("data.akamai_cloudaccess_keys.test", fmt.Sprintf("access_keys.%d.network_configuration.security_network", i), string(key.networkConfiguration.SecurityNetwork)))
			checkFuncs = append(checkFuncs, resource.TestCheckResourceAttr("data.akamai_cloudaccess_keys.test", fmt.Sprintf("access_keys.%d.network_configuration.additional_cdn", i), string(*key.networkConfiguration.AdditionalCDN)))
		}
	}

	return resource.ComposeAggregateTestCheckFunc(checkFuncs...)
}
