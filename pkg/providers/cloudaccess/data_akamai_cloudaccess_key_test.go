package cloudaccess

import (
	"fmt"
	"regexp"
	"testing"

	"github.com/akamai/AkamaiOPEN-edgegrid-golang/v12/pkg/cloudaccess"
	"github.com/akamai/terraform-provider-akamai/v8/pkg/common/date"
	"github.com/akamai/terraform-provider-akamai/v8/pkg/common/ptr"
	"github.com/akamai/terraform-provider-akamai/v8/pkg/common/testutils"
	"github.com/hashicorp/terraform-plugin-testing/helper/resource"
	"github.com/stretchr/testify/mock"
)

func TestDataKey(t *testing.T) {
	expectListAccessKeys := func(client *cloudaccess.Mock, data []cloudaccess.AccessKeyResponse, timesToRun int) {
		listLocationsRes := cloudaccess.ListAccessKeysResponse{
			AccessKeys: data,
		}
		client.On("ListAccessKeys", testutils.MockContext, mock.Anything).Return(&listLocationsRes, nil).Times(timesToRun)
	}

	expectListAccessKeysWithError := func(client *cloudaccess.Mock, timesToRun int) {
		client.On("ListAccessKeys", testutils.MockContext, mock.Anything).Return(nil, fmt.Errorf("list keys failed")).Times(timesToRun)
	}

	stringDate1, err := date.Parse("2021-02-24T09:09:52.782555Z")
	if err != nil {
		t.Fatal(err.Error())
	}
	stringDate2, err := date.Parse("2021-02-26T09:09:15.428314Z")
	if err != nil {
		t.Fatal(err.Error())
	}
	testData := []cloudaccess.AccessKeyResponse{
		{
			AccessKeyName:        "Sales-s3",
			AccessKeyUID:         56514,
			AuthenticationMethod: "AWS4_HMAC_SHA256",
			CreatedBy:            "mrossi",
			CreatedTime:          stringDate1,
			Groups: []cloudaccess.Group{
				{
					ContractIDs: []string{"K-0N7RAK71"},
					GroupID:     32145,
					GroupName:   ptr.To("Sales"),
				},
			},
			LatestVersion: 1,
			NetworkConfiguration: &cloudaccess.SecureNetwork{
				AdditionalCDN:   ptr.To(cloudaccess.CDNType("RUSSIA_CDN")),
				SecurityNetwork: "ENHANCED_TLS",
			},
		},
		{
			AccessKeyName:        "Home automation | s3",
			AccessKeyUID:         56512,
			AuthenticationMethod: "AWS4_HMAC_SHA256",
			CreatedBy:            "tyamada",
			CreatedTime:          stringDate2,
			Groups: []cloudaccess.Group{
				{
					ContractIDs: []string{"C-0N7RAC7"},

					GroupID:   54321,
					GroupName: ptr.To("Smarthomes"),
				},
			},
			LatestVersion: 3,
			NetworkConfiguration: &cloudaccess.SecureNetwork{
				SecurityNetwork: "ENHANCED_TLS",
			},
		},
	}
	tests := map[string]struct {
		configPath string
		init       func(*testing.T, *cloudaccess.Mock, []cloudaccess.AccessKeyResponse)
		mockData   []cloudaccess.AccessKeyResponse
		error      *regexp.Regexp
	}{
		"happy path": {
			configPath: "testdata/TestDataKey/default.tf",
			init: func(_ *testing.T, m *cloudaccess.Mock, testData []cloudaccess.AccessKeyResponse) {
				expectListAccessKeys(m, testData, 3)
			},
			mockData: testData,
		},
		"no name": {
			configPath: "testdata/TestDataKey/no_name.tf",
			init: func(_ *testing.T, m *cloudaccess.Mock, testData []cloudaccess.AccessKeyResponse) {
				expectListAccessKeys(m, testData, 1)
			},
			mockData: testData,
			error:    regexp.MustCompile("no key with given name"),
		},
		"missing name": {
			configPath: "testdata/TestDataKey/missing_name.tf",
			mockData:   testData,
			error:      regexp.MustCompile("The argument \"access_key_name\" is required, but no definition was found."),
		},
		"error listing keys": {
			configPath: "testdata/TestDataKey/default.tf",
			init: func(_ *testing.T, m *cloudaccess.Mock, _ []cloudaccess.AccessKeyResponse) {
				expectListAccessKeysWithError(m, 1)
			},
			error: regexp.MustCompile("list keys failed"),
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
							Check:       checkCloudaccessKeyAttrs(),
							ExpectError: test.error,
						},
					},
				})
			})
			client.AssertExpectations(t)
		})
	}
}

func checkCloudaccessKeyAttrs() resource.TestCheckFunc {
	var checkFuncs []resource.TestCheckFunc

	checkFuncs = append(checkFuncs, resource.TestCheckResourceAttr("data.akamai_cloudaccess_key.test", "access_key_name", "Home automation | s3"))
	checkFuncs = append(checkFuncs, resource.TestCheckResourceAttr("data.akamai_cloudaccess_key.test", "groups.0.contracts_ids.0", "C-0N7RAC7"))
	checkFuncs = append(checkFuncs, resource.TestCheckResourceAttr("data.akamai_cloudaccess_key.test", "network_configuration.security_network", "ENHANCED_TLS"))
	checkFuncs = append(checkFuncs, resource.TestCheckResourceAttr("data.akamai_cloudaccess_key.test", "created_time", "2021-02-26T09:09:15.428314Z"))

	return resource.ComposeAggregateTestCheckFunc(checkFuncs...)
}
