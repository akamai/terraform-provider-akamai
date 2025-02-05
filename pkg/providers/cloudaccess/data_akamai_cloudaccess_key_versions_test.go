package cloudaccess

import (
	"fmt"
	"regexp"
	"strconv"
	"testing"

	"github.com/akamai/AkamaiOPEN-edgegrid-golang/v10/pkg/cloudaccess"
	"github.com/akamai/terraform-provider-akamai/v7/pkg/common/date"
	"github.com/akamai/terraform-provider-akamai/v7/pkg/common/ptr"
	"github.com/akamai/terraform-provider-akamai/v7/pkg/common/testutils"
	"github.com/hashicorp/terraform-plugin-testing/helper/resource"
	"github.com/stretchr/testify/mock"
)

func TestDataKeyVersions(t *testing.T) {
	dateString := "2021-02-26T09:09:15.428314Z"
	dateTime, err := date.Parse(dateString)
	if err != nil {
		t.Fatal(err.Error())
	}
	accessKeyUID := int64(56512)

	listAccessKeysMockRes := cloudaccess.ListAccessKeysResponse{
		AccessKeys: []cloudaccess.AccessKeyResponse{
			{
				AccessKeyName:        "Home automation | s3",
				AccessKeyUID:         accessKeyUID,
				AuthenticationMethod: "AWS4_HMAC_SHA256",
				CreatedBy:            "tyamada",
				CreatedTime:          dateTime,
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
		}}

	listAccessKeyVersionsMockRes := cloudaccess.ListAccessKeyVersionsResponse{
		AccessKeyVersions: []cloudaccess.AccessKeyVersion{
			{
				AccessKeyUID:     accessKeyUID,
				CloudAccessKeyID: nil,
				CreatedBy:        "tyamada",
				CreatedTime:      dateTime,
				DeploymentStatus: cloudaccess.Active,
				Version:          1,
				VersionGUID:      "10000000-7837-11eb-817b-1b3f28104c18",
			},
			{
				AccessKeyUID:     accessKeyUID,
				CloudAccessKeyID: ptr.To("CloudAccessKeyID"),
				CreatedBy:        "tyamada",
				CreatedTime:      dateTime,
				DeploymentStatus: cloudaccess.Active,
				Version:          2,
				VersionGUID:      "20000000-7837-11eb-817b-1b3f28104c18",
			},
		},
	}

	tests := map[string]struct {
		configPath string
		init       func(*testing.T, *cloudaccess.Mock)
		mockData   cloudaccess.ListAccessKeyVersionsResponse
		dateString string
		error      *regexp.Regexp
	}{
		"happy path": {
			configPath: "testdata/TestDataKeyVersions/default.tf",
			init: func(_ *testing.T, m *cloudaccess.Mock) {

				m.On("ListAccessKeys", testutils.MockContext, mock.Anything).Return(&listAccessKeysMockRes, nil)
				m.On("ListAccessKeyVersions", testutils.MockContext, cloudaccess.ListAccessKeyVersionsRequest{
					AccessKeyUID: accessKeyUID,
				}).Return(&listAccessKeyVersionsMockRes, nil)
			},
			mockData:   listAccessKeyVersionsMockRes,
			dateString: dateString,
		},
		"no key versions": {
			configPath: "testdata/TestDataKeyVersions/default.tf",
			init: func(_ *testing.T, m *cloudaccess.Mock) {

				m.On("ListAccessKeys", testutils.MockContext, mock.Anything).Return(&listAccessKeysMockRes, nil)
				m.On("ListAccessKeyVersions", testutils.MockContext, cloudaccess.ListAccessKeyVersionsRequest{
					AccessKeyUID: accessKeyUID,
				}).Return(&cloudaccess.ListAccessKeyVersionsResponse{}, nil)
			},
			mockData:   cloudaccess.ListAccessKeyVersionsResponse{},
			dateString: "",
		},
		"no matching key": {
			configPath: "testdata/TestDataKeyVersions/default.tf",
			init: func(_ *testing.T, m *cloudaccess.Mock) {

				m.On("ListAccessKeys", testutils.MockContext, mock.Anything).Return(&cloudaccess.ListAccessKeysResponse{}, nil)

			},
			error: regexp.MustCompile(`no key with given name`),
		},
		"expect error on list access keys": {
			configPath: "testdata/TestDataKeyVersions/default.tf",
			init: func(_ *testing.T, m *cloudaccess.Mock) {

				m.On("ListAccessKeys", testutils.MockContext, mock.Anything).Return(nil, fmt.Errorf("API error"))

			},
			error: regexp.MustCompile(`API error`),
		},
		"expect error on list access key versions": {
			configPath: "testdata/TestDataKeyVersions/default.tf",
			init: func(_ *testing.T, m *cloudaccess.Mock) {

				m.On("ListAccessKeys", testutils.MockContext, mock.Anything).Return(&listAccessKeysMockRes, nil)
				m.On("ListAccessKeyVersions", testutils.MockContext, cloudaccess.ListAccessKeyVersionsRequest{
					AccessKeyUID: accessKeyUID,
				}).Return(nil, fmt.Errorf("API error"))

			},
			error: regexp.MustCompile(`API error`),
		},
		"invalid configuration - missing access key name": {
			configPath: "testdata/TestDataKeyVersions/invalid.tf",
			error:      regexp.MustCompile(`The argument "access_key_name" is required, but no definition was found.`),
		},
	}
	for name, test := range tests {
		t.Run(name, func(t *testing.T) {
			client := &cloudaccess.Mock{}
			if test.init != nil {
				test.init(t, client)
			}

			useClient(client, func() {
				resource.UnitTest(t, resource.TestCase{
					ProtoV6ProviderFactories: testutils.NewProtoV6ProviderFactory(NewSubprovider()),
					IsUnitTest:               true,
					Steps: []resource.TestStep{
						{
							Config:      testutils.LoadFixtureString(t, test.configPath),
							Check:       checkCloudAccessKeyVersionsAttrs(test.mockData, test.dateString),
							ExpectError: test.error,
						},
					},
				})
			})
			client.AssertExpectations(t)
		})
	}
}
func checkCloudAccessKeyVersionsAttrs(testData cloudaccess.ListAccessKeyVersionsResponse, dateString string) resource.TestCheckFunc {
	var checkFuncs []resource.TestCheckFunc
	if len(testData.AccessKeyVersions) == 0 {
		return resource.ComposeAggregateTestCheckFunc(checkFuncs...)
	}

	checkFuncs = append(checkFuncs, resource.TestCheckResourceAttr("data.akamai_cloudaccess_key_versions.test", "access_key_uid", strconv.FormatInt(testData.AccessKeyVersions[0].AccessKeyUID, 10)))

	for i, version := range testData.AccessKeyVersions {
		versionPrefix := fmt.Sprintf("access_key_versions.%d.", i)
		if version.CloudAccessKeyID != nil {
			checkFuncs = append(checkFuncs, resource.TestCheckResourceAttr("data.akamai_cloudaccess_key_versions.test", versionPrefix+"cloud_access_key_id", *version.CloudAccessKeyID))
		} else {
			checkFuncs = append(checkFuncs, resource.TestCheckNoResourceAttr("data.akamai_cloudaccess_key_versions.test", versionPrefix+"cloud_access_key_id"))
		}

		checkFuncs = append(checkFuncs, resource.TestCheckResourceAttr("data.akamai_cloudaccess_key_versions.test", versionPrefix+"created_by", version.CreatedBy))
		checkFuncs = append(checkFuncs, resource.TestCheckResourceAttr("data.akamai_cloudaccess_key_versions.test", versionPrefix+"created_time", dateString))
		checkFuncs = append(checkFuncs, resource.TestCheckResourceAttr("data.akamai_cloudaccess_key_versions.test", versionPrefix+"version", strconv.FormatInt(version.Version, 10)))
		checkFuncs = append(checkFuncs, resource.TestCheckResourceAttr("data.akamai_cloudaccess_key_versions.test", versionPrefix+"version_guid", version.VersionGUID))
	}

	return resource.ComposeAggregateTestCheckFunc(checkFuncs...)
}
