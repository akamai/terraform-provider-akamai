package cloudaccess

import (
	"fmt"
	"regexp"
	"sort"
	"strconv"
	"testing"

	"github.com/akamai/AkamaiOPEN-edgegrid-golang/v10/pkg/cloudaccess"
	"github.com/akamai/terraform-provider-akamai/v7/pkg/common/ptr"
	"github.com/akamai/terraform-provider-akamai/v7/pkg/common/testutils"
	"github.com/hashicorp/terraform-plugin-testing/helper/resource"
)

type (
	testDataForKeyProperties struct {
		accessKeyUID int64
		properties   map[int64][]property
	}

	property struct {
		propertyName      string
		propertyID        string
		stagingVersion    *int64
		productionVersion *int64
	}
)

func TestDataKeyProperties(t *testing.T) {
	tests := map[string]struct {
		configPath string
		init       func(*cloudaccess.Mock, testDataForKeyProperties)
		mockData   testDataForKeyProperties
		error      *regexp.Regexp
	}{
		"happy path - multiple versions with multiple properties": {
			configPath: "testdata/TestDataKeyProperties/default.tf",
			init: func(m *cloudaccess.Mock, testData testDataForKeyProperties) {
				expectListAccessKeys(m, 3)
				expectListAccessKeyVersions(m, testData, 3)
				expectLookupProperties(m, testData, 3)
			},
			mockData: testDataForKeyProperties{
				accessKeyUID: 1,
				properties: map[int64][]property{
					1: {
						{
							propertyName:      "property1",
							propertyID:        "prp_1",
							stagingVersion:    ptr.To(int64(1)),
							productionVersion: ptr.To(int64(1)),
						},
						{
							propertyName:      "property2",
							propertyID:        "prp_2",
							stagingVersion:    nil,
							productionVersion: ptr.To(int64(2)),
						},
					},
					2: {
						{
							propertyName:      "property3",
							propertyID:        "prp_3",
							stagingVersion:    ptr.To(int64(3)),
							productionVersion: nil,
						},
					},
				},
			},
		},
		"happy path - version with no active properties - nothing in state": {
			configPath: "testdata/TestDataKeyProperties/default.tf",
			init: func(m *cloudaccess.Mock, testData testDataForKeyProperties) {
				expectListAccessKeys(m, 3)
				expectListAccessKeyVersions(m, testData, 3)
				expectLookupProperties(m, testData, 3)
			},
			mockData: testDataForKeyProperties{
				accessKeyUID: 1,
				properties: map[int64][]property{
					1: {
						{
							propertyName:      "property1",
							propertyID:        "prp_1",
							stagingVersion:    nil,
							productionVersion: nil,
						},
					},
				},
			},
		},
		"invalid configuration - missing access key name": {
			configPath: "testdata/TestDataKeyProperties/invalid.tf",
			error:      regexp.MustCompile(`The argument "access_key_name" is required, but no definition was found.`),
		},
		"no match on access key name - expect an error": {
			configPath: "testdata/TestDataKeyProperties/no-match.tf",
			init: func(m *cloudaccess.Mock, _ testDataForKeyProperties) {
				expectListAccessKeys(m, 1)
			},
			error: regexp.MustCompile(`access key with name: 'no-match' does not exist`),
		},
		"expect error on list access keys": {
			configPath: "testdata/TestDataKeyProperties/default.tf",
			init: func(m *cloudaccess.Mock, _ testDataForKeyProperties) {
				m.On("ListAccessKeys", testutils.MockContext, cloudaccess.ListAccessKeysRequest{}).
					Return(nil, fmt.Errorf("API error")).Once()
			},
			error: regexp.MustCompile(`API error`),
		},
		"expect error on list access key versions": {
			configPath: "testdata/TestDataKeyProperties/default.tf",
			init: func(m *cloudaccess.Mock, _ testDataForKeyProperties) {
				expectListAccessKeys(m, 1)
				m.On("ListAccessKeyVersions", testutils.MockContext, cloudaccess.ListAccessKeyVersionsRequest{
					AccessKeyUID: 1,
				}).Return(nil, fmt.Errorf("API error")).Once()
			},
			error: regexp.MustCompile(`API error`),
		},
		"expect error on lookup properties": {
			configPath: "testdata/TestDataKeyProperties/default.tf",
			mockData: testDataForKeyProperties{
				accessKeyUID: 1,
				properties: map[int64][]property{
					1: {
						{
							propertyName:      "property1",
							propertyID:        "prp_1",
							stagingVersion:    nil,
							productionVersion: nil,
						},
					},
				},
			},
			init: func(m *cloudaccess.Mock, testData testDataForKeyProperties) {
				expectListAccessKeys(m, 1)
				expectListAccessKeyVersions(m, testData, 1)
				m.On("LookupProperties", testutils.MockContext, cloudaccess.LookupPropertiesRequest{
					AccessKeyUID: 1,
					Version:      1,
				}).Return(nil, fmt.Errorf("API error")).Once()
			},
			error: regexp.MustCompile(`API error`),
		},
	}
	for name, test := range tests {
		t.Run(name, func(t *testing.T) {
			client := &cloudaccess.Mock{}
			if test.init != nil {
				test.init(client, test.mockData)
			}

			useClient(client, func() {
				resource.UnitTest(t, resource.TestCase{
					ProtoV6ProviderFactories: testutils.NewProtoV6ProviderFactory(NewSubprovider()),
					IsUnitTest:               true,
					Steps: []resource.TestStep{
						{
							Config:      testutils.LoadFixtureString(t, test.configPath),
							Check:       checkKeyPropertiesAttrs(test.mockData),
							ExpectError: test.error,
						},
					},
				})
			})
			client.AssertExpectations(t)
		})
	}
}

func checkKeyPropertiesAttrs(data testDataForKeyProperties) resource.TestCheckFunc {
	var checkFuncs []resource.TestCheckFunc
	var i, propertyCount int
	var sortedKeys []int64

	for k := range data.properties {
		sortedKeys = append(sortedKeys, k)
	}
	sort.Slice(sortedKeys, func(i, j int) bool {
		return sortedKeys[i] < sortedKeys[j]
	})

	for _, k := range sortedKeys {
		v := data.properties[k]
		for _, prp := range v {
			if prp.stagingVersion != nil || prp.productionVersion != nil {
				checkFuncs = append(checkFuncs, resource.TestCheckResourceAttr("data.akamai_cloudaccess_key_properties.test", fmt.Sprintf("properties.%d.property_name", i), prp.propertyName))
				checkFuncs = append(checkFuncs, resource.TestCheckResourceAttr("data.akamai_cloudaccess_key_properties.test", fmt.Sprintf("properties.%d.property_id", i), prp.propertyID))
				checkFuncs = append(checkFuncs, resource.TestCheckResourceAttr("data.akamai_cloudaccess_key_properties.test", fmt.Sprintf("properties.%d.access_key_version", i), strconv.FormatInt(k, 10)))
				if prp.stagingVersion != nil {
					checkFuncs = append(checkFuncs, resource.TestCheckResourceAttr("data.akamai_cloudaccess_key_properties.test", fmt.Sprintf("properties.%d.staging_version", i), strconv.FormatInt(*prp.stagingVersion, 10)))
				} else {
					checkFuncs = append(checkFuncs, resource.TestCheckNoResourceAttr("data.akamai_cloudaccess_key_properties.test", fmt.Sprintf("properties.%d.staging_version", i)))
				}
				if prp.productionVersion != nil {
					checkFuncs = append(checkFuncs, resource.TestCheckResourceAttr("data.akamai_cloudaccess_key_properties.test", fmt.Sprintf("properties.%d.production_version", i), strconv.FormatInt(*prp.productionVersion, 10)))
				} else {
					checkFuncs = append(checkFuncs, resource.TestCheckNoResourceAttr("data.akamai_cloudaccess_key_properties.test", fmt.Sprintf("properties.%d.production_version", i)))
				}
				i++
				propertyCount++
			} else {
				checkFuncs = append(checkFuncs, resource.TestCheckNoResourceAttr("data.akamai_cloudaccess_key_properties.test", "properties"))
			}
		}
	}

	checkFuncs = append(checkFuncs, resource.TestCheckResourceAttr("data.akamai_cloudaccess_key_properties.test", "properties.#", strconv.Itoa(propertyCount)))
	checkFuncs = append(checkFuncs, resource.TestCheckResourceAttr("data.akamai_cloudaccess_key_properties.test", "access_key_uid", strconv.FormatInt(data.accessKeyUID, 10)))

	return resource.ComposeAggregateTestCheckFunc(checkFuncs...)
}

func expectListAccessKeys(client *cloudaccess.Mock, timesToRun int) {
	listAccessKeysReq := cloudaccess.ListAccessKeysRequest{}

	listAccessKeysRes := cloudaccess.ListAccessKeysResponse{
		AccessKeys: []cloudaccess.AccessKeyResponse{
			{
				AccessKeyUID:  2,
				AccessKeyName: "wrong one",
			},
			{
				AccessKeyUID:  1,
				AccessKeyName: "name",
			},
		},
	}

	client.On("ListAccessKeys", testutils.MockContext, listAccessKeysReq).Return(&listAccessKeysRes, nil).Times(timesToRun)
}

func expectListAccessKeyVersions(client *cloudaccess.Mock, data testDataForKeyProperties, timesToRun int) {
	listKeyPropertiesReq := cloudaccess.ListAccessKeyVersionsRequest{
		AccessKeyUID: data.accessKeyUID,
	}

	var listKeyPropertiesRes cloudaccess.ListAccessKeyVersionsResponse
	for key := range data.properties {
		listKeyPropertiesRes.AccessKeyVersions = append(listKeyPropertiesRes.AccessKeyVersions, cloudaccess.AccessKeyVersion{
			AccessKeyUID: data.accessKeyUID,
			Version:      key,
		})
	}

	sort.Slice(listKeyPropertiesRes.AccessKeyVersions, func(i, j int) bool {
		return listKeyPropertiesRes.AccessKeyVersions[i].Version < listKeyPropertiesRes.AccessKeyVersions[j].Version
	})

	client.On("ListAccessKeyVersions", testutils.MockContext, listKeyPropertiesReq).Return(&listKeyPropertiesRes, nil).Times(timesToRun)
}

func expectLookupProperties(client *cloudaccess.Mock, data testDataForKeyProperties, timesToRun int) {
	var accessKeyVersions []int64
	for key := range data.properties {
		accessKeyVersions = append(accessKeyVersions, key)
	}
	sort.Slice(accessKeyVersions, func(i, j int) bool {
		return accessKeyVersions[i] < accessKeyVersions[j]
	})

	for _, key := range accessKeyVersions {
		val := data.properties[key]
		lookupPropertiesReq := cloudaccess.LookupPropertiesRequest{
			AccessKeyUID: data.accessKeyUID,
			Version:      key,
		}

		var lookupPropertiesRes cloudaccess.LookupPropertiesResponse
		for _, prp := range val {
			lookupPropertiesRes.Properties = append(lookupPropertiesRes.Properties, cloudaccess.Property{
				AccessKeyUID:      data.accessKeyUID,
				PropertyID:        prp.propertyID,
				PropertyName:      prp.propertyName,
				ProductionVersion: prp.productionVersion,
				StagingVersion:    prp.stagingVersion,
			})
		}
		client.On("LookupProperties", testutils.MockContext, lookupPropertiesReq).Return(&lookupPropertiesRes, nil).Times(timesToRun)
	}
}
