package cloudaccess

import (
	"errors"
	"fmt"
	"regexp"
	"strings"
	"testing"
	"time"

	"github.com/akamai/AkamaiOPEN-edgegrid-golang/v13/pkg/cloudaccess"
	"github.com/akamai/terraform-provider-akamai/v10/pkg/common/ptr"
	"github.com/akamai/terraform-provider-akamai/v10/pkg/common/test"
	"github.com/akamai/terraform-provider-akamai/v10/pkg/common/testutils"
	"github.com/hashicorp/terraform-plugin-framework/types"
	"github.com/hashicorp/terraform-plugin-testing/helper/resource"
	"github.com/hashicorp/terraform-plugin-testing/knownvalue"
	"github.com/hashicorp/terraform-plugin-testing/plancheck"
	"github.com/hashicorp/terraform-plugin-testing/terraform"
	"github.com/hashicorp/terraform-plugin-testing/tfjsonpath"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/mock"
)

type (
	commonDataForAccessKey struct {
		accessKeyUID         int64
		accessKeyName        string
		authenticationMethod string
		contractID           string
		groupID              int64
		networkConfig        networkConfiguration
		credentialsA         credentials
		credentialsB         credentials
	}

	networkConfiguration struct {
		securityNetwork string
		additionalCDN   string
	}

	credentials struct {
		cloudAccessKeyID     string
		cloudSecretAccessKey string
		primaryKey           bool
	}

	commonDataForProperty struct {
		accessKeyUID      int64
		propertyName      string
		propertyID        string
		stagingVersion    int64
		productionVersion int64
	}

	commonDataForResource struct {
		defaultKey         commonDataForAccessKey // AWS, both credentials A and B
		secondKey          commonDataForAccessKey // second distinct access key (different UID/name)
		credBOnlyKey       commonDataForAccessKey // only credentials_b populated
		emptySecretKey     commonDataForAccessKey // credentials_a with empty secret (post-import)
		noCloudKeyIDKey    commonDataForAccessKey // VP_QUEUE_IT, no cloud_access_key_id
		noCloudKeyIDAVMKey commonDataForAccessKey // AVM_CLOUDINARY, no cloud_access_key_id
		propertyData       commonDataForProperty
	}
)

var (
	// baseAccessKeyChecker provides common attribute checks shared across most test cases.
	baseAccessKeyChecker = test.NewStateChecker("akamai_cloudaccess_key.test").
				CheckEqual("access_key_uid", "12345").
				CheckEqual("access_key_name", "test_key_name").
				CheckEqual("contract_id", "1-CTRACT").
				CheckEqual("group_id", "12345").
				CheckEqual("network_configuration.security_network", "ENHANCED_TLS").
				CheckEqual("network_configuration.additional_cdn", "CHINA_CDN")

	// accessCheckerWithoutCDN extends baseAccessKeyChecker with network configuration checks for keys without additional CDN.
	accessCheckerWithoutCDN = baseAccessKeyChecker.
				CheckMissing("network_configuration.additional_cdn")

	// credentialsABatch contains default attribute values for the credentials_a block (AWS style).
	credentialsABatch = test.AttributeBatch{
		"cloud_access_key_id":     "test_key_id",
		"cloud_secret_access_key": "test_secret",
		"primary_key":             "true",
		"version":                 "1",
		"version_guid":            "asde-efdr-reded",
	}

	// credentialsBBatch contains default attribute values for the credentials_b block (AWS style).
	credentialsBBatch = test.AttributeBatch{
		"cloud_access_key_id":     "test_key_id_2",
		"cloud_secret_access_key": "test_secret_2",
		"primary_key":             "false",
		"version":                 "2",
		"version_guid":            "asdd-ads-dasdas",
	}

	// credentialsANoKeyIDBatch contains default values for credentials_a without cloud_access_key_id (VP_QUEUE_IT / AVM_CLOUDINARY style).
	credentialsANoKeyIDBatch = test.AttributeBatch{
		"cloud_secret_access_key": "test_secret",
		"primary_key":             "true",
		"version":                 "1",
		"version_guid":            "asde-efdr-reded",
	}

	// credentialsBNoKeyIDBatch contains default values for credentials_b without cloud_access_key_id.
	credentialsBNoKeyIDBatch = test.AttributeBatch{
		"cloud_secret_access_key": "test_secret_2",
		"primary_key":             "false",
		"version":                 "2",
		"version_guid":            "asdd-ads-dasdas",
	}

	// credentialsVersion3Batch contains attribute values for credentials at version 3 (AWS style with cloud_access_key_id).
	credentialsVersion3Batch = test.AttributeBatch{
		"cloud_access_key_id":     "test_key_id_3",
		"cloud_secret_access_key": "test_secret_3",
		"primary_key":             "true",
		"version":                 "3",
		"version_guid":            "ffff_eeee-ffffddd",
	}

	// credentialsVersion3NoKeyIDBatch contains attribute values for credentials at version 3 without cloud_access_key_id.
	credentialsVersion3NoKeyIDBatch = test.AttributeBatch{
		"cloud_secret_access_key": "test_secret_3",
		"primary_key":             "true",
		"version":                 "3",
		"version_guid":            "ffff_eeee-ffffddd",
	}

	accessKeyMock = commonDataForAccessKey{
		accessKeyName:        "test_key_name",
		accessKeyUID:         12345,
		authenticationMethod: string(cloudaccess.AuthAWS),
		contractID:           "1-CTRACT",
		groupID:              12345,
		networkConfig: networkConfiguration{
			securityNetwork: string(cloudaccess.NetworkEnhanced),
			additionalCDN:   string(cloudaccess.ChinaCDN),
		},
		credentialsA: credentials{
			cloudAccessKeyID:     "test_key_id",
			cloudSecretAccessKey: "test_secret",
			primaryKey:           true,
		},
		credentialsB: credentials{
			cloudAccessKeyID:     "test_key_id_2",
			cloudSecretAccessKey: "test_secret_2",
			primaryKey:           false,
		},
	}
	secondKeyMock = commonDataForAccessKey{
		accessKeyName:        "test key name2",
		accessKeyUID:         5678,
		authenticationMethod: string(cloudaccess.AuthAWS),
		contractID:           "1-CTRACT",
		groupID:              12345,
		networkConfig: networkConfiguration{
			securityNetwork: string(cloudaccess.NetworkEnhanced),
			additionalCDN:   string(cloudaccess.ChinaCDN),
		},
		credentialsA: credentials{
			cloudAccessKeyID:     "2test_key_id",
			cloudSecretAccessKey: "2test_secret",
			primaryKey:           true,
		},
		credentialsB: credentials{
			cloudAccessKeyID:     "2test_key_id_2",
			cloudSecretAccessKey: "2test_secret_2",
			primaryKey:           false,
		},
	}

	onlyCredBMock = commonDataForAccessKey{
		accessKeyName:        "test_key_name",
		accessKeyUID:         12345,
		authenticationMethod: string(cloudaccess.AuthAWS),
		contractID:           "1-CTRACT",
		groupID:              12345,
		networkConfig: networkConfiguration{
			securityNetwork: string(cloudaccess.NetworkEnhanced),
			additionalCDN:   string(cloudaccess.ChinaCDN),
		},
		credentialsB: credentials{
			cloudAccessKeyID:     "test_key_id",
			cloudSecretAccessKey: "test_secret",
			primaryKey:           true,
		},
	}

	emptySecretMock = commonDataForAccessKey{
		accessKeyName:        "test_key_name",
		accessKeyUID:         12345,
		authenticationMethod: string(cloudaccess.AuthAWS),
		contractID:           "1-CTRACT",
		groupID:              12345,
		networkConfig: networkConfiguration{
			securityNetwork: string(cloudaccess.NetworkEnhanced),
			additionalCDN:   string(cloudaccess.ChinaCDN),
		},
		credentialsA: credentials{
			cloudAccessKeyID:     "test_key_id",
			cloudSecretAccessKey: "",
			primaryKey:           true,
		},
	}

	noCloudAccessKeyIDMock = commonDataForAccessKey{
		accessKeyName:        "test_key_name",
		accessKeyUID:         12345,
		authenticationMethod: string(cloudaccess.AuthVPQueueIt),
		contractID:           "1-CTRACT",
		groupID:              12345,
		networkConfig: networkConfiguration{
			securityNetwork: string(cloudaccess.NetworkEnhanced),
		},
		credentialsA: credentials{
			cloudSecretAccessKey: "test_secret",
			primaryKey:           true,
		},
		credentialsB: credentials{
			cloudSecretAccessKey: "test_secret_2",
			primaryKey:           false,
		},
	}

	noCloudAccessKeyIDAVMMock = commonDataForAccessKey{
		accessKeyName:        "test_key_name",
		accessKeyUID:         12345,
		authenticationMethod: string(cloudaccess.AuthAVMCloudinary),
		contractID:           "1-CTRACT",
		groupID:              12345,
		networkConfig: networkConfiguration{
			securityNetwork: string(cloudaccess.NetworkEnhanced),
		},
		credentialsA: credentials{
			cloudSecretAccessKey: "test_secret",
			primaryKey:           true,
		},
		credentialsB: credentials{
			cloudSecretAccessKey: "test_secret_2",
			primaryKey:           false,
		},
	}

	propertyMock = commonDataForProperty{
		accessKeyUID:      12345,
		propertyID:        "123123",
		propertyName:      "test_property_name",
		stagingVersion:    1,
		productionVersion: 1,
	}

	resourceMock = commonDataForResource{
		defaultKey:         accessKeyMock,
		secondKey:          secondKeyMock,
		credBOnlyKey:       onlyCredBMock,
		emptySecretKey:     emptySecretMock,
		noCloudKeyIDKey:    noCloudAccessKeyIDMock,
		noCloudKeyIDAVMKey: noCloudAccessKeyIDAVMMock,
		propertyData:       propertyMock,
	}

	firstAccessKeyVersion  = int64(1)
	secondAccessKeyVersion = int64(2)
	thirdAccessKeyVersion  = int64(3)

	emptyVersionList       = 0
	oneElementVersionList  = 1
	twoElementsVersionList = 2
)

func TestAccessKeyResource(t *testing.T) {
	t.Parallel()
	pollingInterval = 1 * time.Millisecond
	deleteTimeout = 40 * time.Millisecond
	updateTimeout = 20 * time.Millisecond
	activationTimeout = 20 * time.Millisecond
	tests := map[string]struct {
		configPath string
		init       func(*cloudaccess.Mock, commonDataForResource)
		mockData   commonDataForResource
		steps      []resource.TestStep
		error      *regexp.Regexp
	}{
		"create access key one version": {
			init: func(m *cloudaccess.Mock, resourceData commonDataForResource) {
				mockCreationAccessKeyWith1Version(m, resourceData)
				mockReadAccessKeyWith1Version(m, resourceData)
				mockDeletionAccessKeyWith1Version(m, resourceData)
			},
			mockData: resourceMock,
			steps: []resource.TestStep{
				{
					Config: testutils.LoadFixtureString(t, "testdata/TestResAccessKey/create.tf"),
					Check: baseAccessKeyChecker.
						CheckEqual("authentication_method", "AWS4_HMAC_SHA256").
						CheckEqual("primary_guid", "asde-efdr-reded").
						CheckEqualBatch("credentials_a.", credentialsABatch).
						Build(),
				},
			},
		},
		"create access key one version no cloud access key id - vp queue it": {
			init: func(m *cloudaccess.Mock, resourceData commonDataForResource) {
				accessKey := resourceData.noCloudKeyIDKey
				mockCreationNoCloudAccessKeyID1Version(m, accessKey)
				mockReadAccessKeyNoCloudAccessKeyID(m, accessKey, oneElementVersionList)
				mockDeletionNoCloudAccessKeyID1Version(m, accessKey, resourceData)
			},
			mockData: resourceMock,
			steps: []resource.TestStep{
				{
					Config: testutils.LoadFixtureString(t, "testdata/TestResAccessKey/create_missing_cloud_access_key.tf"),
					Check: accessCheckerWithoutCDN.
						CheckEqual("authentication_method", "VP_QUEUE_IT").
						CheckEqual("primary_guid", "asde-efdr-reded").
						CheckEqualBatch("credentials_a.", credentialsANoKeyIDBatch).
						Build(),
				},
			},
		},
		"create access key one version no cloud access key id - avm cloudinary": {
			init: func(m *cloudaccess.Mock, resourceData commonDataForResource) {
				accessKey := resourceData.noCloudKeyIDAVMKey
				mockCreationNoCloudAccessKeyID1Version(m, accessKey)
				mockReadAccessKeyNoCloudAccessKeyID(m, accessKey, oneElementVersionList)
				mockDeletionNoCloudAccessKeyID1Version(m, accessKey, resourceData)
			},
			mockData: resourceMock,
			steps: []resource.TestStep{
				{
					Config: testutils.LoadFixtureString(t, "testdata/TestResAccessKey/create_missing_cloud_access_key_avm_cloudinary.tf"),
					Check: accessCheckerWithoutCDN.
						CheckEqual("authentication_method", "AVM_CLOUDINARY").
						CheckEqual("primary_guid", "asde-efdr-reded").
						CheckEqualBatch("credentials_a.", credentialsANoKeyIDBatch).
						Build(),
				},
			},
		},
		"create access key two versions": {
			init: func(m *cloudaccess.Mock, resourceData commonDataForResource) {
				mockCreationAccessKeyWith2Versions(m, resourceData)
				mockReadAccessKey(m, resourceData, twoElementsVersionList)
				mockDeletionAccessKeyWith2Versions(m, resourceData)
			},
			mockData: resourceMock,
			steps: []resource.TestStep{
				{
					Config: testutils.LoadFixtureString(t, "testdata/TestResAccessKey/create_2_versions.tf"),
					Check: baseAccessKeyChecker.
						CheckEqual("authentication_method", "AWS4_HMAC_SHA256").
						CheckEqual("primary_guid", "asde-efdr-reded").
						CheckEqualBatch("credentials_a.", credentialsABatch).
						CheckEqualBatch("credentials_b.", credentialsBBatch).
						Build(),
				},
			},
		},
		"create access key two versions no cloud access key id": {
			init: func(m *cloudaccess.Mock, resourceData commonDataForResource) {
				accessKey := resourceData.noCloudKeyIDKey
				mockCreationNoCloudAccessKeyID2Versions(m, accessKey)
				mockReadAccessKeyNoCloudAccessKeyID(m, accessKey, twoElementsVersionList)
				mockDeletionNoCloudAccessKeyID2Versions(m, accessKey, resourceData)
			},
			mockData: resourceMock,
			steps: []resource.TestStep{
				{
					Config: testutils.LoadFixtureString(t, "testdata/TestResAccessKey/create_2_versions_missing_cloud_access_key.tf"),
					Check: accessCheckerWithoutCDN.
						CheckEqual("authentication_method", "VP_QUEUE_IT").
						CheckEqual("primary_guid", "asde-efdr-reded").
						CheckEqualBatch("credentials_a.", credentialsANoKeyIDBatch).
						CheckEqualBatch("credentials_b.", credentialsBNoKeyIDBatch).
						Build(),
				},
			},
		},
		"create access key only credentialsB": {
			init: func(m *cloudaccess.Mock, resourceData commonDataForResource) {
				//creation
				mockCreationAccessKeyUsingCredB(m, resourceData)
				//read
				mockGetAccessKey(m, resourceData.credBOnlyKey).Once()
				mockListAccessKeyVersionsOnly1Version(m, resourceData.credBOnlyKey).Once()
				//delete
				mockListAccessKeyVersionsOnly1Version(m, resourceData.credBOnlyKey).Once()
				mockLookupsPropertiesNoProperties(m, resourceData.propertyData, firstAccessKeyVersion).Once()
				m.On("DeleteAccessKeyVersion", testutils.MockContext, cloudaccess.DeleteAccessKeyVersionRequest{AccessKeyUID: resourceData.credBOnlyKey.accessKeyUID, Version: firstAccessKeyVersion}).
					Return(&cloudaccess.DeleteAccessKeyVersionResponse{
						AccessKeyUID:     resourceData.credBOnlyKey.accessKeyUID,
						CloudAccessKeyID: ptr.To(resourceData.credBOnlyKey.credentialsB.cloudAccessKeyID),
						CreatedBy:        "dev-user",
						CreatedTime:      time.Date(2024, 1, 10, 11, 9, 10, 67708, time.UTC),
						DeploymentStatus: cloudaccess.Active,
						Version:          firstAccessKeyVersion,
						VersionGUID:      "asde-efdr-reded",
					}, nil).Once()
				m.On("GetAccessKeyVersion", testutils.MockContext, cloudaccess.GetAccessKeyVersionRequest{AccessKeyUID: resourceData.defaultKey.accessKeyUID, Version: 1}).
					Return(&cloudaccess.GetAccessKeyVersionResponse{
						AccessKeyUID:     resourceData.credBOnlyKey.accessKeyUID,
						CloudAccessKeyID: ptr.To(resourceData.credBOnlyKey.credentialsB.cloudAccessKeyID),
						CreatedBy:        "dev-user",
						CreatedTime:      time.Date(2024, 1, 10, 11, 9, 10, 67708, time.UTC),
						DeploymentStatus: cloudaccess.Active,
						Version:          firstAccessKeyVersion,
						VersionGUID:      "asde-efdr-reded",
					}, nil).Once()
				mockDeleteAccessKey(m, resourceData.credBOnlyKey).Once()
				var listOfKeysAfterDeletion []commonDataForAccessKey
				mockListAccessKeys(m, append(listOfKeysAfterDeletion, resourceData.secondKey)).Once()
			},
			mockData: resourceMock,
			steps: []resource.TestStep{
				{
					Config: testutils.LoadFixtureString(t, "testdata/TestResAccessKey/create_using_credB.tf"),
					Check: baseAccessKeyChecker.
						CheckEqual("authentication_method", "AWS4_HMAC_SHA256").
						CheckEqual("primary_guid", "asde-efdr-reded").
						CheckEqualBatch("credentials_b.", credentialsABatch).
						Build(),
				},
			},
		},
		"delete access key version only credentialsB": {
			init: func(m *cloudaccess.Mock, resourceData commonDataForResource) {
				//create
				mockCreationAccessKeyUsingCredB(m, resourceData)
				//read
				mockGetAccessKey(m, resourceData.credBOnlyKey).Times(2)
				mockListAccessKeyVersionsOnly1Version(m, resourceData.credBOnlyKey).Once()
				mockListAccessKeyVersionsOnly1Version(m, resourceData.credBOnlyKey).Once()
				//delete 1 version (credB)
				mockListAccessKeyVersionsOnly1Version(m, resourceData.credBOnlyKey).Once()
				mockLookupsPropertiesNoProperties(m, resourceData.propertyData, firstAccessKeyVersion).Once()
				m.On("DeleteAccessKeyVersion", testutils.MockContext, cloudaccess.DeleteAccessKeyVersionRequest{AccessKeyUID: resourceData.credBOnlyKey.accessKeyUID, Version: firstAccessKeyVersion}).
					Return(&cloudaccess.DeleteAccessKeyVersionResponse{
						AccessKeyUID:     resourceData.credBOnlyKey.accessKeyUID,
						CloudAccessKeyID: ptr.To(resourceData.credBOnlyKey.credentialsB.cloudAccessKeyID),
						CreatedBy:        "dev-user",
						CreatedTime:      time.Date(2024, 1, 10, 11, 9, 10, 67708, time.UTC),
						DeploymentStatus: cloudaccess.Active,
						Version:          firstAccessKeyVersion,
						VersionGUID:      "asde-efdr-reded",
					}, nil).Once()
				m.On("GetAccessKeyVersion", testutils.MockContext, cloudaccess.GetAccessKeyVersionRequest{AccessKeyUID: resourceData.defaultKey.accessKeyUID, Version: firstAccessKeyVersion}).
					Return(&cloudaccess.GetAccessKeyVersionResponse{
						AccessKeyUID:     resourceData.credBOnlyKey.accessKeyUID,
						CloudAccessKeyID: ptr.To(resourceData.credBOnlyKey.credentialsB.cloudAccessKeyID),
						CreatedBy:        "dev-user",
						CreatedTime:      time.Date(2024, 1, 10, 11, 9, 10, 67708, time.UTC),
						DeploymentStatus: cloudaccess.PendingDeletion,
						Version:          firstAccessKeyVersion,
						VersionGUID:      "asde-efdr-reded",
					}, nil).Once()
				mockListAccessKeyVersions(m, resourceData.credBOnlyKey, emptyVersionList).Once()
				//read
				mockGetAccessKey(m, resourceData.credBOnlyKey).Once()
				mockListAccessKeyVersions(m, resourceData.credBOnlyKey, emptyVersionList).Once()
				//delete key
				mockListAccessKeyVersions(m, resourceData.credBOnlyKey, emptyVersionList).Once()
				mockDeleteAccessKey(m, resourceData.credBOnlyKey).Once()
				var listOfKeysAfterDeletion []commonDataForAccessKey
				mockListAccessKeys(m, append(listOfKeysAfterDeletion, resourceData.secondKey)).Once()
			},
			mockData: resourceMock,
			steps: []resource.TestStep{
				{
					Config: testutils.LoadFixtureString(t, "testdata/TestResAccessKey/create_using_credB.tf"),
					Check: baseAccessKeyChecker.
						CheckEqual("authentication_method", "AWS4_HMAC_SHA256").
						CheckEqual("primary_guid", "asde-efdr-reded").
						CheckEqualBatch("credentials_b.", credentialsABatch).
						Build(),
				},
				{
					Config: testutils.LoadFixtureString(t, "testdata/TestResAccessKey/creation_no_credentials.tf"),
					Check: baseAccessKeyChecker.
						CheckEqual("authentication_method", "AWS4_HMAC_SHA256").
						CheckEqual("primary_guid", "").
						Build(),
				},
			},
		},
		"basic name update": {
			init: func(m *cloudaccess.Mock, resourceData commonDataForResource) {
				mockCreationAccessKeyWith1Version(m, resourceData)
				mockReadAccessKeyWith1Version(m, resourceData)
				mockReadAccessKeyWith1Version(m, resourceData)
				mockUpdateAccessKey(m, resourceData.defaultKey, "updated_key_name").Once()
				mockGetAccessKeyWithSpecificNameAndVersion(m, resourceData.defaultKey, "updated_key_name", firstAccessKeyVersion).Once()
				mockListAccessKeyVersionsOnly1Version(m, resourceData.defaultKey).Twice()
				mockListAccessKeyVersionsOnly1Version(m, resourceData.defaultKey).Twice()
				mockDeletionAccessKeyWith1Version(m, resourceData)
			},
			mockData: resourceMock,
			steps: []resource.TestStep{
				{
					Config: testutils.LoadFixtureString(t, "testdata/TestResAccessKey/create.tf"),
					Check: baseAccessKeyChecker.
						CheckEqual("authentication_method", "AWS4_HMAC_SHA256").
						CheckEqual("primary_guid", "asde-efdr-reded").
						CheckEqualBatch("credentials_a.", credentialsABatch).
						Build(),
				},
				{
					Config: testutils.LoadFixtureString(t, "testdata/TestResAccessKey/updated_name.tf"),
					Check: baseAccessKeyChecker.
						CheckEqual("access_key_name", "updated_key_name").
						CheckEqual("authentication_method", "AWS4_HMAC_SHA256").
						CheckEqual("primary_guid", "asde-efdr-reded").
						CheckEqualBatch("credentials_a.", credentialsABatch).
						Build(),
				},
			},
		},
		"update only timeout (no timeout to with timeout)": {
			init: func(m *cloudaccess.Mock, resourceData commonDataForResource) {
				mockCreationAccessKeyWith1Version(m, resourceData)
				// step 1 post-apply read
				mockReadAccessKeyWith1Version(m, resourceData)
				// step 2: pre-plan read + post-apply read (Update early returns, no API calls)
				mockReadAccessKeyWith1Version(m, resourceData)
				mockGetAccessKey(m, resourceData.defaultKey).Once()
				mockListAccessKeyVersionsOnly1Version(m, resourceData.defaultKey).Twice()
				mockListAccessKeyVersionsOnly1Version(m, resourceData.defaultKey).Twice()
				mockDeletionAccessKeyWith1Version(m, resourceData)
			},
			mockData: resourceMock,
			steps: []resource.TestStep{
				{
					Config: testutils.LoadFixtureString(t, "testdata/TestResAccessKey/create.tf"),
					Check: baseAccessKeyChecker.
						CheckEqual("authentication_method", "AWS4_HMAC_SHA256").
						CheckEqual("primary_guid", "asde-efdr-reded").
						CheckEqualBatch("credentials_a.", credentialsABatch).
						CheckMissing("timeouts.create").
						CheckMissing("timeouts.update").
						CheckMissing("timeouts.delete").
						Build(),
				},
				{
					Config: testutils.LoadFixtureString(t, "testdata/TestResAccessKey/create_with_timeout.tf"),
					Check: baseAccessKeyChecker.
						CheckEqual("authentication_method", "AWS4_HMAC_SHA256").
						CheckEqual("primary_guid", "asde-efdr-reded").
						CheckEqualBatch("credentials_a.", credentialsABatch).
						CheckEqual("timeouts.create", "12ms").
						CheckEqual("timeouts.update", "1ms").
						CheckEqual("timeouts.delete", "20m").
						Build(),
					ConfigPlanChecks: resource.ConfigPlanChecks{
						PreApply: []plancheck.PlanCheck{
							plancheck.ExpectKnownValue("akamai_cloudaccess_key.test", tfjsonpath.New("timeouts").AtMapKey("create"), knownvalue.StringExact("12ms")),
							plancheck.ExpectKnownValue("akamai_cloudaccess_key.test", tfjsonpath.New("timeouts").AtMapKey("update"), knownvalue.StringExact("1ms")),
							plancheck.ExpectKnownValue("akamai_cloudaccess_key.test", tfjsonpath.New("timeouts").AtMapKey("delete"), knownvalue.StringExact("20m")),
							plancheck.ExpectKnownValue("akamai_cloudaccess_key.test", tfjsonpath.New("access_key_uid"), knownvalue.Int64Exact(12345)),
							plancheck.ExpectKnownValue("akamai_cloudaccess_key.test", tfjsonpath.New("access_key_name"), knownvalue.StringExact("test_key_name")),
							plancheck.ExpectKnownValue("akamai_cloudaccess_key.test", tfjsonpath.New("credentials_a").AtMapKey("cloud_access_key_id"), knownvalue.StringExact("test_key_id")),
						},
					},
				},
			},
		},
		"update only timeout (with timeout to no timeout)": {
			init: func(m *cloudaccess.Mock, resourceData commonDataForResource) {
				mockCreationAccessKeyWith1Version(m, resourceData)
				// step 1 post-apply read
				mockReadAccessKeyWith1Version(m, resourceData)
				// step 2: pre-plan read + post-apply read (Update early returns, no API calls)
				mockReadAccessKeyWith1Version(m, resourceData)
				mockGetAccessKey(m, resourceData.defaultKey).Once()
				mockListAccessKeyVersionsOnly1Version(m, resourceData.defaultKey).Twice()
				mockListAccessKeyVersionsOnly1Version(m, resourceData.defaultKey).Twice()
				mockDeletionAccessKeyWith1Version(m, resourceData)
			},
			mockData: resourceMock,
			steps: []resource.TestStep{
				{
					Config: testutils.LoadFixtureString(t, "testdata/TestResAccessKey/create_with_timeout.tf"),
					Check: baseAccessKeyChecker.
						CheckEqual("authentication_method", "AWS4_HMAC_SHA256").
						CheckEqual("primary_guid", "asde-efdr-reded").
						CheckEqualBatch("credentials_a.", credentialsABatch).
						CheckEqual("timeouts.create", "12ms").
						CheckEqual("timeouts.update", "1ms").
						CheckEqual("timeouts.delete", "20m").
						Build(),
				},
				{
					Config: testutils.LoadFixtureString(t, "testdata/TestResAccessKey/create.tf"),
					Check: baseAccessKeyChecker.
						CheckEqual("authentication_method", "AWS4_HMAC_SHA256").
						CheckEqual("primary_guid", "asde-efdr-reded").
						CheckEqualBatch("credentials_a.", credentialsABatch).
						CheckMissing("timeouts.create").
						CheckMissing("timeouts.update").
						CheckMissing("timeouts.delete").
						Build(),
					ConfigPlanChecks: resource.ConfigPlanChecks{
						PreApply: []plancheck.PlanCheck{
							plancheck.ExpectKnownValue("akamai_cloudaccess_key.test", tfjsonpath.New("timeouts"), knownvalue.Null()),
							plancheck.ExpectKnownValue("akamai_cloudaccess_key.test", tfjsonpath.New("access_key_uid"), knownvalue.Int64Exact(12345)),
							plancheck.ExpectKnownValue("akamai_cloudaccess_key.test", tfjsonpath.New("access_key_name"), knownvalue.StringExact("test_key_name")),
							plancheck.ExpectKnownValue("akamai_cloudaccess_key.test", tfjsonpath.New("credentials_a").AtMapKey("cloud_access_key_id"), knownvalue.StringExact("test_key_id")),
						},
					},
				},
			},
		},
		"update only timeout (value to other value)": {
			init: func(m *cloudaccess.Mock, resourceData commonDataForResource) {
				mockCreationAccessKeyWith1Version(m, resourceData)
				// step 1 post-apply read
				mockReadAccessKeyWith1Version(m, resourceData)
				// step 2: pre-plan read + post-apply read (Update early returns, no API calls)
				mockReadAccessKeyWith1Version(m, resourceData)
				mockGetAccessKey(m, resourceData.defaultKey).Once()
				mockListAccessKeyVersionsOnly1Version(m, resourceData.defaultKey).Twice()
				mockListAccessKeyVersionsOnly1Version(m, resourceData.defaultKey).Twice()
				mockDeletionAccessKeyWith1Version(m, resourceData)
			},
			mockData: resourceMock,
			steps: []resource.TestStep{
				{
					Config: testutils.LoadFixtureString(t, "testdata/TestResAccessKey/create_with_timeout.tf"),
					Check: baseAccessKeyChecker.
						CheckEqual("authentication_method", "AWS4_HMAC_SHA256").
						CheckEqual("primary_guid", "asde-efdr-reded").
						CheckEqualBatch("credentials_a.", credentialsABatch).
						CheckEqual("timeouts.create", "12ms").
						CheckEqual("timeouts.update", "1ms").
						CheckEqual("timeouts.delete", "20m").
						Build(),
				},
				{
					Config: testutils.LoadFixtureString(t, "testdata/TestResAccessKey/create_with_timeout2.tf"),
					Check: baseAccessKeyChecker.
						CheckEqual("authentication_method", "AWS4_HMAC_SHA256").
						CheckEqual("primary_guid", "asde-efdr-reded").
						CheckEqualBatch("credentials_a.", credentialsABatch).
						CheckEqual("timeouts.create", "30ms").
						CheckEqual("timeouts.update", "2ms").
						CheckEqual("timeouts.delete", "40m").
						Build(),
					ConfigPlanChecks: resource.ConfigPlanChecks{
						PreApply: []plancheck.PlanCheck{
							plancheck.ExpectKnownValue("akamai_cloudaccess_key.test", tfjsonpath.New("timeouts").AtMapKey("create"), knownvalue.StringExact("30ms")),
							plancheck.ExpectKnownValue("akamai_cloudaccess_key.test", tfjsonpath.New("timeouts").AtMapKey("update"), knownvalue.StringExact("2ms")),
							plancheck.ExpectKnownValue("akamai_cloudaccess_key.test", tfjsonpath.New("timeouts").AtMapKey("delete"), knownvalue.StringExact("40m")),
							plancheck.ExpectKnownValue("akamai_cloudaccess_key.test", tfjsonpath.New("access_key_uid"), knownvalue.Int64Exact(12345)),
							plancheck.ExpectKnownValue("akamai_cloudaccess_key.test", tfjsonpath.New("access_key_name"), knownvalue.StringExact("test_key_name")),
							plancheck.ExpectKnownValue("akamai_cloudaccess_key.test", tfjsonpath.New("credentials_a").AtMapKey("cloud_access_key_id"), knownvalue.StringExact("test_key_id")),
						},
					},
				},
			},
		},
		"update timeout and name": {
			init: func(m *cloudaccess.Mock, resourceData commonDataForResource) {
				mockCreationAccessKeyWith1Version(m, resourceData)
				mockReadAccessKeyWith1Version(m, resourceData)
				mockReadAccessKeyWith1Version(m, resourceData)
				mockUpdateAccessKey(m, resourceData.defaultKey, "updated_key_name").Once()
				mockGetAccessKeyWithSpecificNameAndVersion(m, resourceData.defaultKey, "updated_key_name", firstAccessKeyVersion).Once()
				mockListAccessKeyVersionsOnly1Version(m, resourceData.defaultKey).Twice()
				mockListAccessKeyVersionsOnly1Version(m, resourceData.defaultKey).Twice()
				mockDeletionAccessKeyWith1Version(m, resourceData)
			},
			mockData: resourceMock,
			steps: []resource.TestStep{
				{
					Config: testutils.LoadFixtureString(t, "testdata/TestResAccessKey/create.tf"),
					Check: baseAccessKeyChecker.
						CheckEqual("authentication_method", "AWS4_HMAC_SHA256").
						CheckEqual("primary_guid", "asde-efdr-reded").
						CheckEqualBatch("credentials_a.", credentialsABatch).
						CheckMissing("timeouts.create").
						CheckMissing("timeouts.update").
						CheckMissing("timeouts.delete").
						Build(),
				},
				{
					Config: testutils.LoadFixtureString(t, "testdata/TestResAccessKey/updated_name_and_timeout.tf"),
					Check: baseAccessKeyChecker.
						CheckEqual("authentication_method", "AWS4_HMAC_SHA256").
						CheckEqual("primary_guid", "asde-efdr-reded").
						CheckEqualBatch("credentials_a.", credentialsABatch).
						CheckEqual("access_key_name", "updated_key_name").
						CheckEqual("timeouts.create", "12ms").
						CheckEqual("timeouts.update", "1ms").
						CheckEqual("timeouts.delete", "20m").
						Build(),
					ConfigPlanChecks: resource.ConfigPlanChecks{
						PreApply: []plancheck.PlanCheck{
							plancheck.ExpectKnownValue("akamai_cloudaccess_key.test", tfjsonpath.New("timeouts").AtMapKey("create"), knownvalue.StringExact("12ms")),
							plancheck.ExpectKnownValue("akamai_cloudaccess_key.test", tfjsonpath.New("timeouts").AtMapKey("update"), knownvalue.StringExact("1ms")),
							plancheck.ExpectKnownValue("akamai_cloudaccess_key.test", tfjsonpath.New("timeouts").AtMapKey("delete"), knownvalue.StringExact("20m")),
							plancheck.ExpectKnownValue("akamai_cloudaccess_key.test", tfjsonpath.New("access_key_uid"), knownvalue.Int64Exact(12345)),
							plancheck.ExpectKnownValue("akamai_cloudaccess_key.test", tfjsonpath.New("access_key_name"), knownvalue.StringExact("updated_key_name")),
							plancheck.ExpectKnownValue("akamai_cloudaccess_key.test", tfjsonpath.New("credentials_a").AtMapKey("cloud_access_key_id"), knownvalue.StringExact("test_key_id")),
						},
					},
				},
			},
		},
		"single-credentials rotation": {
			init: func(m *cloudaccess.Mock, resourceData commonDataForResource) {
				mockCreationAccessKeyWith2Versions(m, resourceData)
				mockReadAccessKey(m, resourceData, twoElementsVersionList)
				mockReadAccessKey(m, resourceData, twoElementsVersionList)
				//delete version 1
				mockListAccessKeyVersions(m, resourceData.defaultKey, twoElementsVersionList).Once()
				mockLookupsPropertiesNoProperties(m, resourceData.propertyData, firstAccessKeyVersion).Once()
				mockDeleteAccessKeyVersion(m, resourceData.defaultKey, firstAccessKeyVersion).Once()
				mockGetAccessKeyVersion(m, resourceData.defaultKey, cloudaccess.PendingDeletion, firstAccessKeyVersion).Once()
				mockListAccessKeyVersions(m, resourceData.defaultKey, oneElementVersionList).Once()
				//create new version (no.3)
				m.On("CreateAccessKeyVersion", testutils.MockContext, cloudaccess.CreateAccessKeyVersionRequest{
					AccessKeyUID: resourceData.defaultKey.accessKeyUID,
					Body: cloudaccess.CreateAccessKeyVersionRequestBody{
						CloudAccessKeyID:     "test_key_id_3",
						CloudSecretAccessKey: "test_secret_3",
					}}).Return(&cloudaccess.CreateAccessKeyVersionResponse{RequestID: 321321, RetryAfter: 1000}, nil).Once()
				m.On("GetAccessKeyVersionStatus", testutils.MockContext, cloudaccess.GetAccessKeyVersionStatusRequest{RequestID: 321321}).
					Return(&cloudaccess.GetAccessKeyVersionStatusResponse{
						AccessKeyVersion: &cloudaccess.KeyVersion{
							AccessKeyUID: resourceData.defaultKey.accessKeyUID,
							Version:      thirdAccessKeyVersion,
						},
						ProcessingStatus: cloudaccess.ProcessingDone,
						RequestDate:      time.Date(2024, 1, 10, 11, 9, 10, 67708, time.UTC),
						RequestedBy:      "dev-user",
					}, nil).Once()
				mockGetAccessKeyVersion(m, resourceData.defaultKey, cloudaccess.Active, thirdAccessKeyVersion).Once()

				//read
				mockGetAccessKey(m, resourceData.defaultKey).Once()
				m.On("ListAccessKeyVersions", testutils.MockContext, cloudaccess.ListAccessKeyVersionsRequest{
					AccessKeyUID: resourceData.defaultKey.accessKeyUID,
				}).Return(&cloudaccess.ListAccessKeyVersionsResponse{AccessKeyVersions: []cloudaccess.AccessKeyVersion{
					{
						AccessKeyUID:     resourceData.defaultKey.accessKeyUID,
						CloudAccessKeyID: ptr.To("test_key_id_3"),
						CreatedBy:        "dev-user",
						CreatedTime:      time.Date(2024, 1, 10, 11, 9, 10, 67708, time.UTC),
						DeploymentStatus: cloudaccess.Active,
						Version:          thirdAccessKeyVersion,
						VersionGUID:      "ffff_eeee-ffffddd",
					},
					{
						AccessKeyUID:     resourceData.defaultKey.accessKeyUID,
						CloudAccessKeyID: ptr.To("test_key_id_2"),
						CreatedBy:        "dev-user",
						CreatedTime:      time.Date(2024, 1, 10, 11, 9, 10, 67708, time.UTC),
						DeploymentStatus: cloudaccess.Active,
						Version:          secondAccessKeyVersion,
						VersionGUID:      "asdd-ads-dasdas",
					},
				},
				}, nil).Once()
				//Delete both credentials ( with version no.3 and no.2)
				m.On("ListAccessKeyVersions", testutils.MockContext, cloudaccess.ListAccessKeyVersionsRequest{
					AccessKeyUID: resourceData.defaultKey.accessKeyUID,
				}).Return(&cloudaccess.ListAccessKeyVersionsResponse{AccessKeyVersions: []cloudaccess.AccessKeyVersion{
					{
						AccessKeyUID:     resourceData.defaultKey.accessKeyUID,
						CloudAccessKeyID: ptr.To("test_key_id_3"),
						CreatedBy:        "dev-user",
						CreatedTime:      time.Date(2024, 1, 10, 11, 9, 10, 67708, time.UTC),
						DeploymentStatus: cloudaccess.Active,
						Version:          thirdAccessKeyVersion,
						VersionGUID:      "ffff_eeee-ffffddd",
					},
					{
						AccessKeyUID:     resourceData.defaultKey.accessKeyUID,
						CloudAccessKeyID: ptr.To("test_key_id_2"),
						CreatedBy:        "dev-user",
						CreatedTime:      time.Date(2024, 1, 10, 11, 9, 10, 67708, time.UTC),
						DeploymentStatus: cloudaccess.Active,
						Version:          secondAccessKeyVersion,
						VersionGUID:      "asdd-ads-dasdas",
					},
				},
				}, nil).Once()
				mockLookupsPropertiesNoProperties(m, resourceData.propertyData, thirdAccessKeyVersion).Once()
				mockLookupsPropertiesNoProperties(m, resourceData.propertyData, secondAccessKeyVersion).Once()
				mockDeleteAccessKeyVersion(m, resourceData.defaultKey, thirdAccessKeyVersion).Once()
				mockGetAccessKeyVersion(m, resourceData.defaultKey, cloudaccess.PendingDeletion, thirdAccessKeyVersion).Once()
				m.On("ListAccessKeyVersions", testutils.MockContext, cloudaccess.ListAccessKeyVersionsRequest{
					AccessKeyUID: resourceData.defaultKey.accessKeyUID,
				}).Return(&cloudaccess.ListAccessKeyVersionsResponse{AccessKeyVersions: []cloudaccess.AccessKeyVersion{
					{
						AccessKeyUID:     resourceData.defaultKey.accessKeyUID,
						CloudAccessKeyID: ptr.To("test_key_id_3"),
						CreatedBy:        "dev-user",
						CreatedTime:      time.Date(2024, 1, 10, 11, 9, 10, 67708, time.UTC),
						DeploymentStatus: cloudaccess.Active,
						Version:          thirdAccessKeyVersion,
						VersionGUID:      "ffff_eeee-ffffddd",
					},
					{
						AccessKeyUID:     resourceData.defaultKey.accessKeyUID,
						CloudAccessKeyID: ptr.To("test_key_id_2"),
						CreatedBy:        "dev-user",
						CreatedTime:      time.Date(2024, 1, 10, 11, 9, 10, 67708, time.UTC),
						DeploymentStatus: cloudaccess.Active,
						Version:          secondAccessKeyVersion,
						VersionGUID:      "asdd-ads-dasdas",
					},
				},
				}, nil).Once()
				mockListAccessKeyVersions(m, resourceData.defaultKey, oneElementVersionList).Once()
				mockDeleteAccessKeyVersion(m, resourceData.defaultKey, secondAccessKeyVersion).Once()
				mockGetAccessKeyVersion(m, resourceData.defaultKey, cloudaccess.PendingDeletion, secondAccessKeyVersion).Once()
				mockListAccessKeyVersions(m, resourceData.defaultKey, oneElementVersionList).Once()
				mockListAccessKeyVersions(m, resourceData.defaultKey, emptyVersionList).Once()
				mockDeleteAccessKey(m, resourceData.defaultKey).Once()
				var listOfKeysAfterDeletion []commonDataForAccessKey
				mockListAccessKeys(m, append(listOfKeysAfterDeletion, resourceData.secondKey)).Once()
			},
			mockData: resourceMock,
			steps: []resource.TestStep{
				{
					Config: testutils.LoadFixtureString(t, "testdata/TestResAccessKey/create_2_versions.tf"),
					Check: baseAccessKeyChecker.
						CheckEqual("authentication_method", "AWS4_HMAC_SHA256").
						CheckEqual("primary_guid", "asde-efdr-reded").
						CheckEqualBatch("credentials_a.", credentialsABatch).
						CheckEqualBatch("credentials_b.", credentialsBBatch).
						Build(),
				},
				{
					Config: testutils.LoadFixtureString(t, "testdata/TestResAccessKey/single_credentials_rotation.tf"),
					Check: baseAccessKeyChecker.
						CheckEqual("authentication_method", "AWS4_HMAC_SHA256").
						CheckEqual("primary_guid", "ffff_eeee-ffffddd").
						CheckEqualBatch("credentials_a.", credentialsVersion3Batch).
						CheckEqualBatch("credentials_b.", credentialsBBatch).
						Build(),
				},
			},
		},
		"single-credentials rotation no cloud access key id": {
			init: func(m *cloudaccess.Mock, resourceData commonDataForResource) {
				accessKey := resourceData.noCloudKeyIDKey
				mockCreationNoCloudAccessKeyID2Versions(m, accessKey)
				mockReadAccessKeyNoCloudAccessKeyID(m, accessKey, twoElementsVersionList)
				mockReadAccessKeyNoCloudAccessKeyID(m, accessKey, twoElementsVersionList)
				//delete version 1
				mockListAccessKeyVersionsNoCloudAccessKeyID(m, accessKey, twoElementsVersionList).Once()
				mockLookupsPropertiesNoProperties(m, resourceData.propertyData, firstAccessKeyVersion).Once()
				mockDeleteAccessKeyVersionNoCloudAccessKeyID(m, accessKey, firstAccessKeyVersion).Once()
				mockGetAccessKeyVersionNoCloudAccessKeyID(m, accessKey, cloudaccess.PendingDeletion, firstAccessKeyVersion).Once()
				mockListAccessKeyVersionsOnlyV2NoCloudAccessKeyID(m, accessKey).Once()
				//create new version (no.3)
				m.On("CreateAccessKeyVersion", testutils.MockContext, cloudaccess.CreateAccessKeyVersionRequest{
					AccessKeyUID: accessKey.accessKeyUID,
					Body: cloudaccess.CreateAccessKeyVersionRequestBody{
						CloudAccessKeyID:     "",
						CloudSecretAccessKey: "test_secret_3",
					}}).Return(&cloudaccess.CreateAccessKeyVersionResponse{RequestID: 321321, RetryAfter: 1000}, nil).Once()
				m.On("GetAccessKeyVersionStatus", testutils.MockContext, cloudaccess.GetAccessKeyVersionStatusRequest{RequestID: 321321}).
					Return(&cloudaccess.GetAccessKeyVersionStatusResponse{
						AccessKeyVersion: &cloudaccess.KeyVersion{
							AccessKeyUID: accessKey.accessKeyUID,
							Version:      thirdAccessKeyVersion,
						},
						ProcessingStatus: cloudaccess.ProcessingDone,
						RequestDate:      time.Date(2024, 1, 10, 11, 9, 10, 67708, time.UTC),
						RequestedBy:      "dev-user",
					}, nil).Once()
				mockGetAccessKeyVersionNoCloudAccessKeyID(m, accessKey, cloudaccess.Active, thirdAccessKeyVersion).Once()
				//read
				mockGetAccessKey(m, accessKey).Once()
				mockListAccessKeyVersionsV3AndV2NoCloudAccessKeyID(m, accessKey).Once()
				//Delete both credentials (with version no.3 and no.2)
				mockDeletionNoCloudAccessKeyIDAfterRotation(m, accessKey, resourceData)
			},
			mockData: resourceMock,
			steps: []resource.TestStep{
				{
					Config: testutils.LoadFixtureString(t, "testdata/TestResAccessKey/create_2_versions_missing_cloud_access_key.tf"),
					Check: accessCheckerWithoutCDN.
						CheckEqual("authentication_method", "VP_QUEUE_IT").
						CheckEqual("primary_guid", "asde-efdr-reded").
						CheckEqualBatch("credentials_a.", credentialsANoKeyIDBatch).
						CheckEqualBatch("credentials_b.", credentialsBNoKeyIDBatch).
						Build(),
				},
				{
					Config: testutils.LoadFixtureString(t, "testdata/TestResAccessKey/single_credentials_rotation_missing_cloud_access_key.tf"),
					Check: accessCheckerWithoutCDN.
						CheckEqual("authentication_method", "VP_QUEUE_IT").
						CheckEqual("primary_guid", "ffff_eeee-ffffddd").
						CheckEqualBatch("credentials_a.", credentialsVersion3NoKeyIDBatch).
						CheckEqualBatch("credentials_b.", credentialsBNoKeyIDBatch).
						Build(),
				},
			},
		},
		"cross-credentials rotation of access key": {
			init: func(m *cloudaccess.Mock, resourceData commonDataForResource) {
				mockCreationAccessKeyWith2Versions(m, resourceData)
				mockReadAccessKey(m, resourceData, twoElementsVersionList)
				mockReadAccessKey(m, resourceData, twoElementsVersionList)
				//delete version 2
				mockListAccessKeyVersions(m, resourceData.defaultKey, twoElementsVersionList).Once()
				mockLookupsPropertiesNoProperties(m, resourceData.propertyData, secondAccessKeyVersion).Once()
				mockDeleteAccessKeyVersion(m, resourceData.defaultKey, secondAccessKeyVersion).Once()
				mockGetAccessKeyVersion(m, resourceData.defaultKey, cloudaccess.PendingDeletion, secondAccessKeyVersion).Once()
				mockListAccessKeyVersions(m, resourceData.defaultKey, twoElementsVersionList).Once()
				mockListAccessKeyVersionsOnly1Version(m, resourceData.defaultKey).Once()

				//create new version (no.3)
				m.On("CreateAccessKeyVersion", testutils.MockContext, cloudaccess.CreateAccessKeyVersionRequest{
					AccessKeyUID: resourceData.defaultKey.accessKeyUID,
					Body: cloudaccess.CreateAccessKeyVersionRequestBody{
						CloudAccessKeyID:     "test_key_id_3",
						CloudSecretAccessKey: "test_secret_3",
					}}).Return(&cloudaccess.CreateAccessKeyVersionResponse{RequestID: 321321, RetryAfter: 1000}, nil).Once()
				m.On("GetAccessKeyVersionStatus", testutils.MockContext, cloudaccess.GetAccessKeyVersionStatusRequest{RequestID: 321321}).
					Return(&cloudaccess.GetAccessKeyVersionStatusResponse{
						AccessKeyVersion: &cloudaccess.KeyVersion{
							AccessKeyUID: resourceData.defaultKey.accessKeyUID,
							Version:      thirdAccessKeyVersion,
						},
						ProcessingStatus: cloudaccess.ProcessingDone,
						RequestDate:      time.Date(2024, 1, 10, 11, 9, 10, 67708, time.UTC),
						RequestedBy:      "dev-user",
					}, nil).Once()
				mockGetAccessKeyVersion(m, resourceData.defaultKey, cloudaccess.Active, thirdAccessKeyVersion).Once()

				//read
				mockGetAccessKeyWithSpecificNameAndVersion(m, resourceData.defaultKey, resourceData.defaultKey.accessKeyName, thirdAccessKeyVersion).Once()
				m.On("ListAccessKeyVersions", testutils.MockContext, cloudaccess.ListAccessKeyVersionsRequest{
					AccessKeyUID: resourceData.defaultKey.accessKeyUID,
				}).Return(&cloudaccess.ListAccessKeyVersionsResponse{AccessKeyVersions: []cloudaccess.AccessKeyVersion{
					{
						AccessKeyUID:     resourceData.defaultKey.accessKeyUID,
						CloudAccessKeyID: ptr.To("test_key_id_3"),
						CreatedBy:        "dev-user",
						CreatedTime:      time.Date(2024, 1, 10, 11, 9, 10, 67708, time.UTC),
						DeploymentStatus: cloudaccess.Active,
						Version:          thirdAccessKeyVersion,
						VersionGUID:      "ffff_eeee-ffffddd",
					},
					{
						AccessKeyUID:     resourceData.defaultKey.accessKeyUID,
						CloudAccessKeyID: ptr.To("test_key_id"),
						CreatedBy:        "dev-user",
						CreatedTime:      time.Date(2024, 1, 10, 11, 9, 10, 67708, time.UTC),
						DeploymentStatus: cloudaccess.Active,
						Version:          firstAccessKeyVersion,
						VersionGUID:      "asde-efdr-reded",
					},
				},
				}, nil).Once()
				//Delete both credentials ( with version no.3 and no.1)
				m.On("ListAccessKeyVersions", testutils.MockContext, cloudaccess.ListAccessKeyVersionsRequest{
					AccessKeyUID: resourceData.defaultKey.accessKeyUID,
				}).Return(&cloudaccess.ListAccessKeyVersionsResponse{AccessKeyVersions: []cloudaccess.AccessKeyVersion{
					{
						AccessKeyUID:     resourceData.defaultKey.accessKeyUID,
						CloudAccessKeyID: ptr.To("test_key_id_3"),
						CreatedBy:        "dev-user",
						CreatedTime:      time.Date(2024, 1, 10, 11, 9, 10, 67708, time.UTC),
						DeploymentStatus: cloudaccess.Active,
						Version:          thirdAccessKeyVersion,
						VersionGUID:      "ffff_eeee-ffffddd",
					},
					{
						AccessKeyUID:     resourceData.defaultKey.accessKeyUID,
						CloudAccessKeyID: ptr.To("test_key_id"),
						CreatedBy:        "dev-user",
						CreatedTime:      time.Date(2024, 1, 10, 11, 9, 10, 67708, time.UTC),
						DeploymentStatus: cloudaccess.Active,
						Version:          firstAccessKeyVersion,
						VersionGUID:      "asde-efdr-reded",
					},
				},
				}, nil).Once()
				mockLookupsPropertiesNoProperties(m, resourceData.propertyData, thirdAccessKeyVersion).Once()
				mockDeleteAccessKeyVersion(m, resourceData.defaultKey, thirdAccessKeyVersion).Once()
				mockGetAccessKeyVersion(m, resourceData.defaultKey, cloudaccess.PendingDeletion, thirdAccessKeyVersion).Once()
				m.On("ListAccessKeyVersions", testutils.MockContext, cloudaccess.ListAccessKeyVersionsRequest{
					AccessKeyUID: resourceData.defaultKey.accessKeyUID,
				}).Return(&cloudaccess.ListAccessKeyVersionsResponse{AccessKeyVersions: []cloudaccess.AccessKeyVersion{
					{
						AccessKeyUID:     resourceData.defaultKey.accessKeyUID,
						CloudAccessKeyID: ptr.To("test_key_id_3"),
						CreatedBy:        "dev-user",
						CreatedTime:      time.Date(2024, 1, 10, 11, 9, 10, 67708, time.UTC),
						DeploymentStatus: cloudaccess.Active,
						Version:          thirdAccessKeyVersion,
						VersionGUID:      "ffff_eeee-ffffddd",
					},
					{
						AccessKeyUID:     resourceData.defaultKey.accessKeyUID,
						CloudAccessKeyID: ptr.To("test_key_id"),
						CreatedBy:        "dev-user",
						CreatedTime:      time.Date(2024, 1, 10, 11, 9, 10, 67708, time.UTC),
						DeploymentStatus: cloudaccess.Active,
						Version:          firstAccessKeyVersion,
						VersionGUID:      "asde-efdr-reded",
					},
				},
				}, nil).Once()
				mockListAccessKeyVersionsOnly1Version(m, resourceData.defaultKey).Once()
				mockLookupsPropertiesNoProperties(m, resourceData.propertyData, firstAccessKeyVersion).Once()
				mockDeleteAccessKeyVersion(m, resourceData.defaultKey, firstAccessKeyVersion).Once()
				mockGetAccessKeyVersion(m, resourceData.defaultKey, cloudaccess.PendingDeletion, firstAccessKeyVersion).Once()
				mockListAccessKeyVersionsOnly1Version(m, resourceData.defaultKey).Once()
				mockListAccessKeyVersions(m, resourceData.defaultKey, emptyVersionList).Once()
				mockDeleteAccessKey(m, resourceData.defaultKey).Once()
				var listOfKeysAfterDeletion []commonDataForAccessKey
				mockListAccessKeys(m, append(listOfKeysAfterDeletion, resourceData.secondKey)).Once()
			},
			mockData: resourceMock,
			steps: []resource.TestStep{
				{
					Config: testutils.LoadFixtureString(t, "testdata/TestResAccessKey/create_2_versions.tf"),
					Check: baseAccessKeyChecker.
						CheckEqual("authentication_method", "AWS4_HMAC_SHA256").
						CheckEqual("primary_guid", "asde-efdr-reded").
						CheckEqualBatch("credentials_a.", credentialsABatch).
						CheckEqualBatch("credentials_b.", credentialsBBatch).
						Build(),
				},
				{
					Config: testutils.LoadFixtureString(t, "testdata/TestResAccessKey/cross_credentials_rotation.tf"),
					Check: baseAccessKeyChecker.
						CheckEqual("authentication_method", "AWS4_HMAC_SHA256").
						CheckEqual("primary_guid", "ffff_eeee-ffffddd").
						CheckEqualBatch("credentials_a.", credentialsABatch).
						CheckEqual("credentials_a.primary_key", "false").
						CheckEqualBatch("credentials_b.", credentialsVersion3Batch).
						Build(),
				},
			},
		},
		"cross-credentials rotation no cloud access key id": {
			init: func(m *cloudaccess.Mock, resourceData commonDataForResource) {
				accessKey := resourceData.noCloudKeyIDKey
				mockCreationNoCloudAccessKeyID2Versions(m, accessKey)
				mockReadAccessKeyNoCloudAccessKeyID(m, accessKey, twoElementsVersionList)
				mockReadAccessKeyNoCloudAccessKeyID(m, accessKey, twoElementsVersionList)
				// delete version 2 (credB)
				mockListAccessKeyVersionsNoCloudAccessKeyID(m, accessKey, twoElementsVersionList).Once()
				mockLookupsPropertiesNoProperties(m, resourceData.propertyData, secondAccessKeyVersion).Once()
				mockDeleteAccessKeyVersionNoCloudAccessKeyID(m, accessKey, secondAccessKeyVersion).Once()
				mockGetAccessKeyVersionNoCloudAccessKeyID(m, accessKey, cloudaccess.PendingDeletion, secondAccessKeyVersion).Once()
				mockListAccessKeyVersionsNoCloudAccessKeyID(m, accessKey, twoElementsVersionList).Once()
				mockListAccessKeyVersionsNoCloudAccessKeyID(m, accessKey, oneElementVersionList).Once()
				// create new version 3 (new credB with test_secret_3)
				m.On("CreateAccessKeyVersion", testutils.MockContext, cloudaccess.CreateAccessKeyVersionRequest{
					AccessKeyUID: accessKey.accessKeyUID,
					Body: cloudaccess.CreateAccessKeyVersionRequestBody{
						CloudAccessKeyID:     "",
						CloudSecretAccessKey: "test_secret_3",
					}}).Return(&cloudaccess.CreateAccessKeyVersionResponse{RequestID: 321321, RetryAfter: 1000}, nil).Once()
				m.On("GetAccessKeyVersionStatus", testutils.MockContext, cloudaccess.GetAccessKeyVersionStatusRequest{RequestID: 321321}).
					Return(&cloudaccess.GetAccessKeyVersionStatusResponse{
						AccessKeyVersion: &cloudaccess.KeyVersion{
							AccessKeyUID: accessKey.accessKeyUID,
							Version:      thirdAccessKeyVersion,
						},
						ProcessingStatus: cloudaccess.ProcessingDone,
						RequestDate:      time.Date(2024, 1, 10, 11, 9, 10, 67708, time.UTC),
						RequestedBy:      "dev-user",
					}, nil).Once()
				mockGetAccessKeyVersionNoCloudAccessKeyID(m, accessKey, cloudaccess.Active, thirdAccessKeyVersion).Once()
				// read after update (v1 + v3)
				mockGetAccessKey(m, accessKey).Once()
				mockListAccessKeyVersionsV1AndV3NoCloudAccessKeyID(m, accessKey).Once()
				// delete cleanup (v1 + v3)
				mockDeletionNoCloudAccessKeyIDAfterCrossRotation(m, accessKey, resourceData)
			},
			mockData: resourceMock,
			steps: []resource.TestStep{
				{
					Config: testutils.LoadFixtureString(t, "testdata/TestResAccessKey/create_2_versions_missing_cloud_access_key.tf"),
					Check: accessCheckerWithoutCDN.
						CheckEqual("authentication_method", "VP_QUEUE_IT").
						CheckEqual("primary_guid", "asde-efdr-reded").
						CheckEqualBatch("credentials_a.", credentialsANoKeyIDBatch).
						CheckEqualBatch("credentials_b.", credentialsBNoKeyIDBatch).
						Build(),
				},
				{
					Config: testutils.LoadFixtureString(t, "testdata/TestResAccessKey/cross_credentials_rotation_missing_cloud_access_key.tf"),
					Check: accessCheckerWithoutCDN.
						CheckEqual("authentication_method", "VP_QUEUE_IT").
						CheckEqual("primary_guid", "ffff_eeee-ffffddd").
						CheckEqualBatch("credentials_a.", credentialsANoKeyIDBatch).
						CheckEqual("credentials_a.primary_key", "false").
						CheckEqualBatch("credentials_b.", credentialsVersion3NoKeyIDBatch).
						Build(),
				},
			},
		},
		"change primary flag": {
			init: func(m *cloudaccess.Mock, resourceData commonDataForResource) {
				mockCreationAccessKeyWith2Versions(m, resourceData)
				mockReadAccessKey(m, resourceData, twoElementsVersionList)
				mockReadAccessKey(m, resourceData, twoElementsVersionList)
				//only change of flag, and primary guid
				mockListAccessKeyVersions(m, resourceData.defaultKey, twoElementsVersionList).Once()
				mockReadAccessKey(m, resourceData, twoElementsVersionList)
				mockDeletionAccessKeyWith2Versions(m, resourceData)
			},
			mockData: resourceMock,
			steps: []resource.TestStep{
				{
					Config: testutils.LoadFixtureString(t, "testdata/TestResAccessKey/create_2_versions.tf"),
					Check: baseAccessKeyChecker.
						CheckEqual("authentication_method", "AWS4_HMAC_SHA256").
						CheckEqual("primary_guid", "asde-efdr-reded").
						CheckEqualBatch("credentials_a.", credentialsABatch).
						CheckEqualBatch("credentials_b.", credentialsBBatch).
						Build(),
				},
				{
					Config: testutils.LoadFixtureString(t, "testdata/TestResAccessKey/swap_primary_key.tf"),
					Check: baseAccessKeyChecker.
						CheckEqual("authentication_method", "AWS4_HMAC_SHA256").
						CheckEqual("primary_guid", "asdd-ads-dasdas").
						CheckEqualBatch("credentials_a.", credentialsABatch).
						CheckEqual("credentials_a.primary_key", "false").
						CheckEqualBatch("credentials_b.", credentialsBBatch).
						CheckEqual("credentials_b.primary_key", "true").
						Build(),
				},
			},
		},
		"change primary flag no cloud access key id": {
			init: func(m *cloudaccess.Mock, resourceData commonDataForResource) {
				accessKey := resourceData.noCloudKeyIDKey
				mockCreationNoCloudAccessKeyID2Versions(m, accessKey)
				mockReadAccessKeyNoCloudAccessKeyID(m, accessKey, twoElementsVersionList)
				mockReadAccessKeyNoCloudAccessKeyID(m, accessKey, twoElementsVersionList)
				// only change of primary flag - no version deletion/creation
				mockListAccessKeyVersionsNoCloudAccessKeyID(m, accessKey, twoElementsVersionList).Once()
				mockReadAccessKeyNoCloudAccessKeyID(m, accessKey, twoElementsVersionList)
				mockDeletionNoCloudAccessKeyID2Versions(m, accessKey, resourceData)
			},
			mockData: resourceMock,
			steps: []resource.TestStep{
				{
					Config: testutils.LoadFixtureString(t, "testdata/TestResAccessKey/create_2_versions_missing_cloud_access_key.tf"),
					Check: accessCheckerWithoutCDN.
						CheckEqual("authentication_method", "VP_QUEUE_IT").
						CheckEqual("primary_guid", "asde-efdr-reded").
						CheckEqualBatch("credentials_a.", credentialsANoKeyIDBatch).
						CheckEqualBatch("credentials_b.", credentialsBNoKeyIDBatch).
						Build(),
				},
				{
					Config: testutils.LoadFixtureString(t, "testdata/TestResAccessKey/swap_primary_key_no_cloud_access_key_id.tf"),
					Check: accessCheckerWithoutCDN.
						CheckEqual("authentication_method", "VP_QUEUE_IT").
						CheckEqual("primary_guid", "asdd-ads-dasdas").
						CheckEqualBatch("credentials_a.", credentialsANoKeyIDBatch).
						CheckEqual("credentials_a.primary_key", "false").
						CheckEqualBatch("credentials_b.", credentialsBNoKeyIDBatch).
						CheckEqual("credentials_b.primary_key", "true").
						Build(),
				},
			},
		},
		"delete one version of access key": {
			init: func(m *cloudaccess.Mock, resourceData commonDataForResource) {
				mockCreationAccessKeyWith2Versions(m, resourceData)
				mockReadAccessKey(m, resourceData, twoElementsVersionList)
				mockReadAccessKey(m, resourceData, twoElementsVersionList)
				//delete 2nd version
				mockListAccessKeyVersions(m, resourceData.defaultKey, twoElementsVersionList).Once()
				mockLookupsPropertiesNoProperties(m, resourceData.propertyData, secondAccessKeyVersion).Once()
				mockDeleteAccessKeyVersion(m, resourceData.defaultKey, secondAccessKeyVersion).Once()
				mockGetAccessKeyVersion(m, resourceData.defaultKey, cloudaccess.PendingDeletion, secondAccessKeyVersion).Once()
				mockListAccessKeyVersionsOnly1Version(m, resourceData.defaultKey).Once()

				mockReadAccessKeyWith1Version(m, resourceData)
				// delete 1 version
				mockListAccessKeyVersionsOnly1Version(m, resourceData.defaultKey).Once()
				mockLookupsPropertiesNoProperties(m, resourceData.propertyData, firstAccessKeyVersion).Once()
				mockDeleteAccessKeyVersion(m, resourceData.defaultKey, firstAccessKeyVersion).Once()
				mockGetAccessKeyVersion(m, resourceData.defaultKey, cloudaccess.PendingDeletion, firstAccessKeyVersion).Once()
				mockListAccessKeyVersionsOnly1Version(m, resourceData.defaultKey).Once()
				mockListAccessKeyVersions(m, resourceData.defaultKey, emptyVersionList).Once()
				mockDeleteAccessKey(m, resourceData.defaultKey).Once()
				var listOfKeysAfterDeletion []commonDataForAccessKey
				mockListAccessKeys(m, append(listOfKeysAfterDeletion, resourceData.secondKey)).Once()
			},
			mockData: resourceMock,
			steps: []resource.TestStep{
				{
					Config: testutils.LoadFixtureString(t, "testdata/TestResAccessKey/create_2_versions.tf"),
					Check: baseAccessKeyChecker.
						CheckEqual("authentication_method", "AWS4_HMAC_SHA256").
						CheckEqual("primary_guid", "asde-efdr-reded").
						CheckEqualBatch("credentials_a.", credentialsABatch).
						CheckEqualBatch("credentials_b.", credentialsBBatch).
						Build(),
				},
				{
					Config: testutils.LoadFixtureString(t, "testdata/TestResAccessKey/create.tf"),
					Check: baseAccessKeyChecker.
						CheckEqual("authentication_method", "AWS4_HMAC_SHA256").
						CheckEqual("primary_guid", "asde-efdr-reded").
						CheckEqualBatch("credentials_a.", credentialsABatch).
						Build(),
				},
			},
		},
		"delete two versions of access key": {
			init: func(m *cloudaccess.Mock, resourceData commonDataForResource) {
				mockCreationAccessKeyWith2Versions(m, resourceData)
				mockReadAccessKey(m, resourceData, twoElementsVersionList)
				mockReadAccessKey(m, resourceData, twoElementsVersionList)
				//delete 2 versions
				mockListAccessKeyVersions(m, resourceData.defaultKey, twoElementsVersionList).Once()
				mockLookupsPropertiesNoProperties(m, resourceData.propertyData, secondAccessKeyVersion).Once()
				mockDeleteAccessKeyVersion(m, resourceData.defaultKey, secondAccessKeyVersion).Once()
				mockGetAccessKeyVersion(m, resourceData.defaultKey, cloudaccess.PendingDeletion, secondAccessKeyVersion).Once()
				mockListAccessKeyVersions(m, resourceData.defaultKey, twoElementsVersionList).Once()
				mockListAccessKeyVersionsOnly1Version(m, resourceData.defaultKey).Once()
				mockLookupsPropertiesNoProperties(m, resourceData.propertyData, firstAccessKeyVersion).Once()
				mockDeleteAccessKeyVersion(m, resourceData.defaultKey, firstAccessKeyVersion).Once()
				mockGetAccessKeyVersion(m, resourceData.defaultKey, cloudaccess.PendingDeletion, firstAccessKeyVersion).Once()
				mockListAccessKeyVersionsOnly1Version(m, resourceData.defaultKey).Once()
				mockListAccessKeyVersions(m, resourceData.defaultKey, emptyVersionList).Once()

				mockReadAccessKey(m, resourceData, emptyVersionList)
				// delete key with no versions
				mockListAccessKeyVersions(m, resourceData.defaultKey, emptyVersionList).Twice()
				mockDeleteAccessKey(m, resourceData.defaultKey).Once()
				var listOfKeysAfterDeletion []commonDataForAccessKey
				mockListAccessKeys(m, append(listOfKeysAfterDeletion, resourceData.secondKey)).Once()
			},
			mockData: resourceMock,
			steps: []resource.TestStep{
				{
					Config: testutils.LoadFixtureString(t, "testdata/TestResAccessKey/create_2_versions.tf"),
					Check: baseAccessKeyChecker.
						CheckEqual("authentication_method", "AWS4_HMAC_SHA256").
						CheckEqual("primary_guid", "asde-efdr-reded").
						CheckEqualBatch("credentials_a.", credentialsABatch).
						CheckEqualBatch("credentials_b.", credentialsBBatch).
						Build(),
				},
				{
					Config: testutils.LoadFixtureString(t, "testdata/TestResAccessKey/create_no_versions.tf"),
					Check: baseAccessKeyChecker.
						CheckEqual("authentication_method", "AWS4_HMAC_SHA256").
						CheckEqual("primary_guid", "").
						Build(),
				},
			},
		},
		"change order of credentials": {
			init: func(m *cloudaccess.Mock, resourceData commonDataForResource) {
				mockCreationAccessKeyWith2Versions(m, resourceData)
				mockReadAccessKey(m, resourceData, twoElementsVersionList)
				mockReadAccessKey(m, resourceData, twoElementsVersionList)
				mockDeletionAccessKeyWith2Versions(m, resourceData)
			},
			mockData: resourceMock,
			steps: []resource.TestStep{
				{
					Config: testutils.LoadFixtureString(t, "testdata/TestResAccessKey/create_2_versions.tf"),
					Check: baseAccessKeyChecker.
						CheckEqual("authentication_method", "AWS4_HMAC_SHA256").
						CheckEqual("primary_guid", "asde-efdr-reded").
						CheckEqualBatch("credentials_a.", credentialsABatch).
						CheckEqualBatch("credentials_b.", credentialsBBatch).
						Build(),
				},
				{
					Config:      testutils.LoadFixtureString(t, "testdata/TestResAccessKey/changed_order.tf"),
					ExpectError: regexp.MustCompile("cannot change order of `credentials_a` and `credentials_b`"),
				},
			},
		},
		"change order of credentials no cloud access key id - vp queue it": {
			init: func(m *cloudaccess.Mock, resourceData commonDataForResource) {
				accessKey := resourceData.noCloudKeyIDKey
				mockCreationNoCloudAccessKeyID2Versions(m, accessKey)
				mockReadAccessKeyNoCloudAccessKeyID(m, accessKey, twoElementsVersionList)
				mockReadAccessKeyNoCloudAccessKeyID(m, accessKey, twoElementsVersionList)
				mockDeletionNoCloudAccessKeyID2Versions(m, accessKey, resourceData)
			},
			mockData: resourceMock,
			steps: []resource.TestStep{
				{
					Config: testutils.LoadFixtureString(t, "testdata/TestResAccessKey/create_2_versions_missing_cloud_access_key.tf"),
					Check: accessCheckerWithoutCDN.
						CheckEqual("authentication_method", "VP_QUEUE_IT").
						CheckEqual("primary_guid", "asde-efdr-reded").
						CheckEqualBatch("credentials_a.", credentialsANoKeyIDBatch).
						CheckEqualBatch("credentials_b.", credentialsBNoKeyIDBatch).
						Build(),
				},
				{
					Config:      testutils.LoadFixtureString(t, "testdata/TestResAccessKey/changed_order_no_cloud_access_key_id.tf"),
					ExpectError: regexp.MustCompile("cannot change order of `credentials_a` and `credentials_b`"),
				},
			},
		},
		"change order of credentials no cloud access key id - avm cloudinary": {
			init: func(m *cloudaccess.Mock, resourceData commonDataForResource) {
				accessKey := resourceData.noCloudKeyIDAVMKey
				mockCreationNoCloudAccessKeyID2Versions(m, accessKey)
				mockReadAccessKeyNoCloudAccessKeyID(m, accessKey, twoElementsVersionList)
				mockReadAccessKeyNoCloudAccessKeyID(m, accessKey, twoElementsVersionList)
				mockDeletionNoCloudAccessKeyID2Versions(m, accessKey, resourceData)
			},
			mockData: resourceMock,
			steps: []resource.TestStep{
				{
					Config: testutils.LoadFixtureString(t, "testdata/TestResAccessKey/create_2_versions_missing_cloud_access_key_avm_cloudinary.tf"),
					Check: accessCheckerWithoutCDN.
						CheckEqual("authentication_method", "AVM_CLOUDINARY").
						CheckEqual("primary_guid", "asde-efdr-reded").
						CheckEqualBatch("credentials_a.", credentialsANoKeyIDBatch).
						CheckEqualBatch("credentials_b.", credentialsBNoKeyIDBatch).
						Build(),
				},
				{
					Config:      testutils.LoadFixtureString(t, "testdata/TestResAccessKey/changed_order_no_cloud_access_key_id_avm_cloudinary.tf"),
					ExpectError: regexp.MustCompile("cannot change order of `credentials_a` and `credentials_b`"),
				},
			},
		},
		"change secret block": {
			init: func(m *cloudaccess.Mock, resourceData commonDataForResource) {
				mockCreationAccessKeyWith2Versions(m, resourceData)
				mockReadAccessKey(m, resourceData, twoElementsVersionList)
				mockReadAccessKey(m, resourceData, twoElementsVersionList)
				mockDeletionAccessKeyWith2Versions(m, resourceData)
			},
			mockData: resourceMock,
			steps: []resource.TestStep{
				{
					Config: testutils.LoadFixtureString(t, "testdata/TestResAccessKey/create_2_versions.tf"),
					Check: baseAccessKeyChecker.
						CheckEqual("authentication_method", "AWS4_HMAC_SHA256").
						CheckEqual("primary_guid", "asde-efdr-reded").
						CheckEqualBatch("credentials_a.", credentialsABatch).
						CheckEqualBatch("credentials_b.", credentialsBBatch).
						Build(),
				},
				{
					Config:      testutils.LoadFixtureString(t, "testdata/TestResAccessKey/changed_secret.tf"),
					ExpectError: regexp.MustCompile(`\s*cannot update cloud access secret without update of cloud access key id,\s*expect update of secret after import with no API calls`),
				},
			},
		},
		"change secret block after import": {
			init: func(m *cloudaccess.Mock, resourceData commonDataForResource) {

				mockGetAccessKeyVersion(m, resourceData.defaultKey, cloudaccess.Active, firstAccessKeyVersion).Once()

				mockGetAccessKeyVersion(m, resourceData.defaultKey, cloudaccess.Active, secondAccessKeyVersion).Once()

				mockReadAccessKey(m, resourceData, twoElementsVersionList)

				mockReadAccessKey(m, resourceData, twoElementsVersionList)
				mockReadAccessKey(m, resourceData, twoElementsVersionList)

				mockReadAccessKey(m, resourceData, twoElementsVersionList)

				mockListAccessKeyVersions(m, resourceData.defaultKey, twoElementsVersionList).Once()
				mockLookupsPropertiesNoProperties(m, resourceData.propertyData, firstAccessKeyVersion).Once()
				mockLookupsPropertiesNoProperties(m, resourceData.propertyData, secondAccessKeyVersion).Once()
				mockDeleteAccessKeyVersion(m, resourceData.defaultKey, firstAccessKeyVersion).Once()
				mockDeleteAccessKeyVersion(m, resourceData.defaultKey, secondAccessKeyVersion).Once()
				mockDeleteAccessKey(m, resourceData.defaultKey).Once()
				var listOfKeysAfterDeletion []commonDataForAccessKey
				mockListAccessKeys(m, append(listOfKeysAfterDeletion, resourceData.secondKey)).Once()
			},
			mockData: resourceMock,
			steps: []resource.TestStep{

				{
					Config:                               testutils.LoadFixtureString(t, "testdata/TestResAccessKey/create_2_versions.tf"),
					ImportState:                          true,
					ImportStateId:                        "12345",
					ResourceName:                         "akamai_cloudaccess_key.test",
					ImportStateCheck:                     checkImport(),
					ImportStateVerifyIdentifierAttribute: "access_key_uid",
					ImportStatePersist:                   true,
				},
				{
					Config: testutils.LoadFixtureString(t, "testdata/TestResAccessKey/changed_secret.tf"),
					Check: baseAccessKeyChecker.
						CheckEqual("authentication_method", "AWS4_HMAC_SHA256").
						CheckEqual("primary_guid", "asde-efdr-reded").
						CheckEqualBatch("credentials_a.", credentialsABatch).
						CheckEqual("credentials_a.cloud_secret_access_key", "changed_secret").
						CheckEqualBatch("credentials_b.", credentialsBBatch).
						Build(),
				},
			},
		},
		"detect drift - one version deleted in ui": {
			init: func(m *cloudaccess.Mock, resourceData commonDataForResource) {
				//step 1 creation of access key with two versions
				mockCreationAccessKeyWith2Versions(m, resourceData)
				//step 2 both versions available on server
				mockReadAccessKey(m, resourceData, twoElementsVersionList)
				//step 3 one version deleted via UI
				mockReadAccessKeyWith1Version(m, resourceData)
				//step 4 after drift terraform wants to create 2nd version again
				mockCreateAccessKeyVersion(m, resourceData.defaultKey).Once()
				mockGetAccessKeyVersionStatus(m, resourceData.defaultKey, 124, secondAccessKeyVersion).Once()
				mockGetAccessKeyVersion(m, resourceData.defaultKey, cloudaccess.Active, secondAccessKeyVersion).Once()
				//step 5 both versions available on server
				mockReadAccessKey(m, resourceData, twoElementsVersionList)
				//step 6 delete 2 versions
				mockDeletionAccessKeyWith2Versions(m, resourceData)
			},
			mockData: resourceMock,
			steps: []resource.TestStep{
				{
					Config: testutils.LoadFixtureString(t, "testdata/TestResAccessKey/create_2_versions.tf"),
					Check: baseAccessKeyChecker.
						CheckEqual("authentication_method", "AWS4_HMAC_SHA256").
						CheckEqual("primary_guid", "asde-efdr-reded").
						CheckEqualBatch("credentials_a.", credentialsABatch).
						CheckEqualBatch("credentials_b.", credentialsBBatch).
						Build(),
				},
				{
					Config: testutils.LoadFixtureString(t, "testdata/TestResAccessKey/create_2_versions.tf"),
					Check: baseAccessKeyChecker.
						CheckEqual("authentication_method", "AWS4_HMAC_SHA256").
						CheckEqual("primary_guid", "asde-efdr-reded").
						CheckEqualBatch("credentials_a.", credentialsABatch).
						CheckEqualBatch("credentials_b.", credentialsBBatch).
						Build(),
				},
			},
		},
		"detect drift - one version deleted in ui no cloud access key id": {
			init: func(m *cloudaccess.Mock, resourceData commonDataForResource) {
				accessKey := resourceData.noCloudKeyIDKey
				//step 1 creation of access key with two versions
				mockCreationNoCloudAccessKeyID2Versions(m, accessKey)
				//step 2 both versions available on server
				mockReadAccessKeyNoCloudAccessKeyID(m, accessKey, twoElementsVersionList)
				//step 3 one version deleted via UI
				mockReadAccessKeyNoCloudAccessKeyID(m, accessKey, oneElementVersionList)
				//step 4 after drift terraform wants to create 2nd version again
				mockCreateAccessKeyVersionNoCloudAccessKeyIDForCredB(m, accessKey).Once()
				mockGetAccessKeyVersionStatus(m, accessKey, 124, secondAccessKeyVersion).Once()
				mockGetAccessKeyVersionNoCloudAccessKeyID(m, accessKey, cloudaccess.Active, secondAccessKeyVersion).Once()
				//step 5 both versions available on server
				mockReadAccessKeyNoCloudAccessKeyID(m, accessKey, twoElementsVersionList)
				//step 6 delete 2 versions
				mockDeletionNoCloudAccessKeyID2Versions(m, accessKey, resourceData)
			},
			mockData: resourceMock,
			steps: []resource.TestStep{
				{
					Config: testutils.LoadFixtureString(t, "testdata/TestResAccessKey/create_2_versions_missing_cloud_access_key.tf"),
					Check: accessCheckerWithoutCDN.
						CheckEqual("authentication_method", "VP_QUEUE_IT").
						CheckEqual("primary_guid", "asde-efdr-reded").
						CheckEqualBatch("credentials_a.", credentialsANoKeyIDBatch).
						CheckEqualBatch("credentials_b.", credentialsBNoKeyIDBatch).
						Build(),
				},
				{
					Config: testutils.LoadFixtureString(t, "testdata/TestResAccessKey/create_2_versions_missing_cloud_access_key.tf"),
					Check: accessCheckerWithoutCDN.
						CheckEqual("authentication_method", "VP_QUEUE_IT").
						CheckEqual("primary_guid", "asde-efdr-reded").
						CheckEqualBatch("credentials_a.", credentialsANoKeyIDBatch).
						CheckEqualBatch("credentials_b.", credentialsBNoKeyIDBatch).
						Build(),
				},
			},
		},
		"detect drift - one version added in ui": {
			init: func(m *cloudaccess.Mock, resourceData commonDataForResource) {
				//step 1 creation of access key with one version
				mockCreationAccessKeyWith1Version(m, resourceData)
				//step 2 one version available on server
				mockReadAccessKeyWith1Version(m, resourceData)
				//step 3 second version created via UI
				mockReadAccessKey(m, resourceData, twoElementsVersionList)
				//step 4 after drift terraform wants to delete 2nd version
				mockListAccessKeyVersions(m, resourceData.defaultKey, twoElementsVersionList).Once()
				mockLookupsPropertiesNoProperties(m, resourceData.propertyData, secondAccessKeyVersion).Once()
				mockDeleteAccessKeyVersion(m, resourceData.defaultKey, secondAccessKeyVersion).Once()
				mockGetAccessKeyVersion(m, resourceData.defaultKey, cloudaccess.PendingDeletion, secondAccessKeyVersion).Once()
				mockListAccessKeyVersionsOnly1Version(m, resourceData.defaultKey).Once()
				//step 5 one version available on server
				mockReadAccessKeyWith1Version(m, resourceData)
				//step 6 delete one version
				mockDeletionAccessKeyWith1Version(m, resourceData)
			},
			mockData: resourceMock,
			steps: []resource.TestStep{
				{
					Config: testutils.LoadFixtureString(t, "testdata/TestResAccessKey/create.tf"),
					Check: baseAccessKeyChecker.
						CheckEqual("authentication_method", "AWS4_HMAC_SHA256").
						CheckEqual("primary_guid", "asde-efdr-reded").
						CheckEqualBatch("credentials_a.", credentialsABatch).
						Build(),
				},
				{
					Config: testutils.LoadFixtureString(t, "testdata/TestResAccessKey/create.tf"),
					Check: baseAccessKeyChecker.
						CheckEqual("authentication_method", "AWS4_HMAC_SHA256").
						CheckEqual("primary_guid", "asde-efdr-reded").
						CheckEqualBatch("credentials_a.", credentialsABatch).
						Build(),
				},
			},
		},
		"detect drift - one version added in ui no cloud access key id": {
			init: func(m *cloudaccess.Mock, resourceData commonDataForResource) {
				accessKey := resourceData.noCloudKeyIDKey
				//step 1 creation of access key with one version
				mockCreationNoCloudAccessKeyID1Version(m, accessKey)
				//step 2 one version available on server
				mockReadAccessKeyNoCloudAccessKeyID(m, accessKey, oneElementVersionList)
				//step 3 second version created via UI
				mockReadAccessKeyNoCloudAccessKeyID(m, accessKey, twoElementsVersionList)
				//step 4 after drift terraform wants to delete 2nd version
				mockListAccessKeyVersionsNoCloudAccessKeyID(m, accessKey, twoElementsVersionList).Once()
				mockLookupsPropertiesNoProperties(m, resourceData.propertyData, secondAccessKeyVersion).Once()
				mockDeleteAccessKeyVersionNoCloudAccessKeyID(m, accessKey, secondAccessKeyVersion).Once()
				mockGetAccessKeyVersionNoCloudAccessKeyID(m, accessKey, cloudaccess.PendingDeletion, secondAccessKeyVersion).Once()
				mockListAccessKeyVersionsNoCloudAccessKeyID(m, accessKey, oneElementVersionList).Once()
				//step 5 one version available on server
				mockReadAccessKeyNoCloudAccessKeyID(m, accessKey, oneElementVersionList)
				//step 6 delete 1 version
				mockDeletionNoCloudAccessKeyID1Version(m, accessKey, resourceData)
			},
			mockData: resourceMock,
			steps: []resource.TestStep{
				{
					Config: testutils.LoadFixtureString(t, "testdata/TestResAccessKey/create_missing_cloud_access_key.tf"),
					Check: accessCheckerWithoutCDN.
						CheckEqual("authentication_method", "VP_QUEUE_IT").
						CheckEqual("primary_guid", "asde-efdr-reded").
						CheckEqualBatch("credentials_a.", credentialsANoKeyIDBatch).
						Build(),
				},
				{
					Config: testutils.LoadFixtureString(t, "testdata/TestResAccessKey/create_missing_cloud_access_key.tf"),
					Check: accessCheckerWithoutCDN.
						CheckEqual("authentication_method", "VP_QUEUE_IT").
						CheckEqual("primary_guid", "asde-efdr-reded").
						CheckEqualBatch("credentials_a.", credentialsANoKeyIDBatch).
						Build(),
				},
			},
		},
		"detect drift - whole key deleted in ui": {
			init: func(m *cloudaccess.Mock, resourceData commonDataForResource) {
				//step 1 creation of access key with two versions
				mockCreationAccessKeyWith2Versions(m, resourceData)
				//step 2 both versions available on server
				mockReadAccessKey(m, resourceData, twoElementsVersionList)
				//step 3 whole key deleted via UI
				mockGetAccessKeyNotFound(m, resourceData.defaultKey).Once()
				//step 4 after drift terraform wants to create new key
				mockCreationAccessKeyWith2Versions(m, resourceData)
				//step 5 new key recreated on server
				mockReadAccessKey(m, resourceData, twoElementsVersionList)
				//step 6 delete 2 versions
				mockDeletionAccessKeyWith2Versions(m, resourceData)
			},
			mockData: resourceMock,
			steps: []resource.TestStep{
				{
					Config: testutils.LoadFixtureString(t, "testdata/TestResAccessKey/create_2_versions.tf"),
					Check: baseAccessKeyChecker.
						CheckEqual("authentication_method", "AWS4_HMAC_SHA256").
						CheckEqual("primary_guid", "asde-efdr-reded").
						CheckEqualBatch("credentials_a.", credentialsABatch).
						CheckEqualBatch("credentials_b.", credentialsBBatch).
						Build(),
				},
				{
					Config: testutils.LoadFixtureString(t, "testdata/TestResAccessKey/create_2_versions.tf"),
					Check: baseAccessKeyChecker.
						CheckEqual("authentication_method", "AWS4_HMAC_SHA256").
						CheckEqual("primary_guid", "asde-efdr-reded").
						CheckEqualBatch("credentials_a.", credentialsABatch).
						CheckEqualBatch("credentials_b.", credentialsBBatch).
						Build(),
				},
			},
		},
		"detect drift - whole key deleted in ui no cloud access key id": {
			init: func(m *cloudaccess.Mock, resourceData commonDataForResource) {
				accessKey := resourceData.noCloudKeyIDKey
				//step 1 creation of access key with two versions
				mockCreationNoCloudAccessKeyID2Versions(m, accessKey)
				//step 2 both versions available on server
				mockReadAccessKeyNoCloudAccessKeyID(m, accessKey, twoElementsVersionList)
				//step 3 whole key deleted via UI
				mockGetAccessKeyNotFound(m, accessKey).Once()
				//step 4 after drift terraform wants to create new key
				mockCreationNoCloudAccessKeyID2Versions(m, accessKey)
				//step 5 new key recreated on server
				mockReadAccessKeyNoCloudAccessKeyID(m, accessKey, twoElementsVersionList)
				//step 6 delete 2 versions
				mockDeletionNoCloudAccessKeyID2Versions(m, accessKey, resourceData)
			},
			mockData: resourceMock,
			steps: []resource.TestStep{
				{
					Config: testutils.LoadFixtureString(t, "testdata/TestResAccessKey/create_2_versions_missing_cloud_access_key.tf"),
					Check: accessCheckerWithoutCDN.
						CheckEqual("authentication_method", "VP_QUEUE_IT").
						CheckEqual("primary_guid", "asde-efdr-reded").
						CheckEqualBatch("credentials_a.", credentialsANoKeyIDBatch).
						CheckEqualBatch("credentials_b.", credentialsBNoKeyIDBatch).
						Build(),
				},
				{
					Config: testutils.LoadFixtureString(t, "testdata/TestResAccessKey/create_2_versions_missing_cloud_access_key.tf"),
					Check: accessCheckerWithoutCDN.
						CheckEqual("authentication_method", "VP_QUEUE_IT").
						CheckEqual("primary_guid", "asde-efdr-reded").
						CheckEqualBatch("credentials_a.", credentialsANoKeyIDBatch).
						CheckEqualBatch("credentials_b.", credentialsBNoKeyIDBatch).
						Build(),
				},
			},
		},
		"check whether access key secret sensitive": {
			init: func(m *cloudaccess.Mock, resourceData commonDataForResource) {
				mockCreationAccessKeyWith1Version(m, resourceData)
				mockReadAccessKeyWith1Version(m, resourceData)
				mockDeletionAccessKeyWith1Version(m, resourceData)
			},
			mockData: resourceMock,
			steps: []resource.TestStep{
				{
					Config: testutils.LoadFixtureString(t, "testdata/TestResAccessKey/create.tf"),
					Check: baseAccessKeyChecker.
						CheckEqual("authentication_method", "AWS4_HMAC_SHA256").
						CheckEqual("primary_guid", "asde-efdr-reded").
						CheckEqualBatch("credentials_a.", credentialsABatch).
						Build(),
					ConfigPlanChecks: resource.ConfigPlanChecks{
						PreApply: []plancheck.PlanCheck{
							plancheck.ExpectSensitiveValue("akamai_cloudaccess_key.test", tfjsonpath.New("credentials_a").AtMapKey("cloud_secret_access_key")),
						},
					},
				},
			},
		},
		"check whether access_key_uid is known on plan level during update": {
			init: func(m *cloudaccess.Mock, resourceData commonDataForResource) {
				mockCreationAccessKeyWith1Version(m, resourceData)
				mockReadAccessKeyWith1Version(m, resourceData)
				mockReadAccessKeyWith1Version(m, resourceData)
				mockUpdateAccessKey(m, resourceData.defaultKey, "updated_key_name").Once()
				mockGetAccessKeyWithSpecificNameAndVersion(m, resourceData.defaultKey, "updated_key_name", firstAccessKeyVersion).Once()
				mockListAccessKeyVersionsOnly1Version(m, resourceData.defaultKey).Twice()
				mockListAccessKeyVersionsOnly1Version(m, resourceData.defaultKey).Twice()
				mockDeletionAccessKeyWith1Version(m, resourceData)
			},
			mockData: resourceMock,
			steps: []resource.TestStep{
				{
					Config: testutils.LoadFixtureString(t, "testdata/TestResAccessKey/create.tf"),
					Check: baseAccessKeyChecker.
						CheckEqual("authentication_method", "AWS4_HMAC_SHA256").
						CheckEqual("primary_guid", "asde-efdr-reded").
						CheckEqualBatch("credentials_a.", credentialsABatch).
						Build(),
				},
				{
					Config: testutils.LoadFixtureString(t, "testdata/TestResAccessKey/updated_name.tf"),
					Check: test.NewStateChecker("akamai_cloudaccess_key.test").
						CheckEqual("access_key_uid", "12345").
						Build(),
					ConfigPlanChecks: resource.ConfigPlanChecks{
						PreApply: []plancheck.PlanCheck{
							plancheck.ExpectKnownValue("akamai_cloudaccess_key.test", tfjsonpath.New("access_key_uid"), knownvalue.Int64Exact(12345)),
						},
					},
				},
			},
		},
		"all fields provided externally": {
			init: func(m *cloudaccess.Mock, resourceData commonDataForResource) {
				mockCreationAccessKeyWith1Version(m, resourceData)
				mockReadAccessKeyWith1Version(m, resourceData)
				mockDeletionAccessKeyWith1Version(m, resourceData)
			},
			mockData: resourceMock,
			steps: []resource.TestStep{
				{
					Config: testutils.LoadFixtureString(t, "testdata/TestResAccessKey/all_external_fields.tf"),
					Check: test.NewStateChecker("akamai_cloudaccess_key.test.0").
						CheckEqual("authentication_method", "AWS4_HMAC_SHA256").
						CheckEqual("primary_guid", "asde-efdr-reded").
						CheckEqual("access_key_uid", "12345").
						CheckEqual("access_key_name", "test_key_name").
						CheckEqual("contract_id", "1-CTRACT").
						CheckEqual("group_id", "12345").
						CheckEqualBatch("credentials_a.", credentialsABatch).
						CheckEqual("network_configuration.additional_cdn", "CHINA_CDN").
						CheckEqual("network_configuration.security_network", "ENHANCED_TLS").
						Build(),
				},
			},
		},
		"single external credential": {
			init: func(m *cloudaccess.Mock, resourceData commonDataForResource) {
				mockCreationAccessKeyWith1Version(m, resourceData)
				mockReadAccessKeyWith1Version(m, resourceData)
				mockDeletionAccessKeyWith1Version(m, resourceData)
			},
			mockData: resourceMock,
			steps: []resource.TestStep{
				{
					Config: testutils.LoadFixtureString(t, "testdata/TestResAccessKey/single_external_credential.tf"),
					Check: test.NewStateChecker("akamai_cloudaccess_key.test").
						CheckEqual("authentication_method", "AWS4_HMAC_SHA256").
						CheckEqual("primary_guid", "asde-efdr-reded").
						CheckEqual("access_key_uid", "12345").
						CheckEqual("access_key_name", "test_key_name").
						CheckEqual("contract_id", "1-CTRACT").
						CheckEqual("group_id", "12345").
						CheckEqualBatch("credentials_a.", credentialsABatch).
						CheckEqual("network_configuration.additional_cdn", "CHINA_CDN").
						CheckEqual("network_configuration.security_network", "ENHANCED_TLS").
						Build(),
				},
			},
		},
		"missing contract id": {
			steps: []resource.TestStep{{
				Config:      testutils.LoadFixtureString(t, "testdata/TestResAccessKey/missing_contract.tf"),
				ExpectError: regexp.MustCompile(`The argument "contract_id" is required, but no definition was found`),
			}},
		},
		"missing group id": {
			steps: []resource.TestStep{
				{
					Config:      testutils.LoadFixtureString(t, "testdata/TestResAccessKey/missing_group.tf"),
					ExpectError: regexp.MustCompile(`The argument "group_id" is required, but no definition was found`),
				},
			},
		},
		"missing security network": {
			steps: []resource.TestStep{
				{
					Config:      testutils.LoadFixtureString(t, "testdata/TestResAccessKey/missing_security_network.tf"),
					ExpectError: regexp.MustCompile("\\s*Inappropriate value for attribute \"network_configuration\": attribute\\s*\"security_network\" is required."),
				},
			},
		},
		"missing primary key": {
			steps: []resource.TestStep{
				{
					Config:      testutils.LoadFixtureString(t, "testdata/TestResAccessKey/missing_primary_key.tf"),
					ExpectError: regexp.MustCompile("\\s*Inappropriate value for attribute \"credentials_a\": attribute \"primary_key\" is\\s*required."),
				},
			},
		},
		"missing cloud access key when required": {
			steps: []resource.TestStep{
				{
					Config:      testutils.LoadFixtureString(t, "testdata/TestResAccessKey/missing_cloud_access_key.tf"),
					ExpectError: regexp.MustCompile(`cloud access key id missing error`),
				},
			},
		},
		"missing cloud access secret": {
			steps: []resource.TestStep{
				{
					Config:      testutils.LoadFixtureString(t, "testdata/TestResAccessKey/missing_cloud_access_secret.tf"),
					ExpectError: regexp.MustCompile("\\s*Inappropriate value for attribute \"credentials_a\": attribute\\s*\"cloud_secret_access_key\" is required."),
				},
			},
		},
		"no credentials for creation": {
			steps: []resource.TestStep{
				{
					Config:      testutils.LoadFixtureString(t, "testdata/TestResAccessKey/creation_no_credentials.tf"),
					ExpectError: regexp.MustCompile(`at least one credentials are required for creation`),
				},
			},
		},
		"non-unique cloud access key id": {
			steps: []resource.TestStep{
				{
					Config:      testutils.LoadFixtureString(t, "testdata/TestResAccessKey/creation_non_unique_cloud_key_id.tf"),
					ExpectError: regexp.MustCompile("'cloud_access_key_id' should be unique for each pair of credentials"),
				},
			},
		},
		"additional cdn not allowed - vp queue it": {
			steps: []resource.TestStep{
				{
					Config:      testutils.LoadFixtureString(t, "testdata/TestResAccessKey/vp_queue_it_with_additional_cdn.tf"),
					ExpectError: regexp.MustCompile(`additional cdn not allowed error`),
				},
			},
		},
		"additional cdn not allowed - avm cloudinary": {
			steps: []resource.TestStep{
				{
					Config:      testutils.LoadFixtureString(t, "testdata/TestResAccessKey/avm_cloudinary_with_additional_cdn.tf"),
					ExpectError: regexp.MustCompile(`additional cdn not allowed error`),
				},
			},
		},
		"timeout on creation": {
			init: func(m *cloudaccess.Mock, resourceData commonDataForResource) {
				mockCreateAccessKey(m, resourceData.defaultKey).Once()
				mockGetAccessKeyStatus(m, 12345, resourceData.defaultKey).Once()
				//artificial sleep to trigger 20 ms timeout
				time.Sleep(21 * time.Millisecond)
				mockGetAccessKeyVersion(m, resourceData.defaultKey, cloudaccess.PendingActivation, firstAccessKeyVersion).Once() //timeout
			},
			mockData: resourceMock,
			steps: []resource.TestStep{
				{
					Config:      testutils.LoadFixtureString(t, "testdata/TestResAccessKey/create_with_timeout.tf"),
					ExpectError: regexp.MustCompile("reached activation timeout"),
				},
			},
		},
		"timeout on one version deletion": {
			init: func(m *cloudaccess.Mock, resourceData commonDataForResource) {
				mockCreationAccessKeyWith2Versions(m, resourceData)
				mockReadAccessKey(m, resourceData, twoElementsVersionList)
				mockReadAccessKey(m, resourceData, twoElementsVersionList)
				//delete 2nd version
				mockListAccessKeyVersions(m, resourceData.defaultKey, twoElementsVersionList).Once()
				mockLookupsPropertiesNoProperties(m, resourceData.propertyData, secondAccessKeyVersion).Once()
				mockDeleteAccessKeyVersion(m, resourceData.defaultKey, secondAccessKeyVersion).Once()
				time.Sleep(50 * time.Millisecond)
				mockGetAccessKeyVersion(m, resourceData.defaultKey, cloudaccess.PendingDeletion, secondAccessKeyVersion).Once()
				mockListAccessKeyVersions(m, resourceData.defaultKey, twoElementsVersionList).Once()

				// delete 1 version
				mockListAccessKeyVersions(m, resourceData.defaultKey, twoElementsVersionList).Once()
				mockLookupsPropertiesNoProperties(m, resourceData.propertyData, secondAccessKeyVersion).Once()
				mockDeleteAccessKeyVersion(m, resourceData.defaultKey, secondAccessKeyVersion).Once()
				mockGetAccessKeyVersion(m, resourceData.defaultKey, cloudaccess.PendingDeletion, secondAccessKeyVersion).Once()
				mockLookupsPropertiesNoProperties(m, resourceData.propertyData, firstAccessKeyVersion).Once()
				mockDeleteAccessKeyVersion(m, resourceData.defaultKey, firstAccessKeyVersion).Once()
				mockGetAccessKeyVersion(m, resourceData.defaultKey, cloudaccess.PendingDeletion, firstAccessKeyVersion).Once()
				mockListAccessKeyVersions(m, resourceData.defaultKey, emptyVersionList).Twice()

				// delete key with no versions
				mockDeleteAccessKey(m, resourceData.defaultKey).Once()
				var listOfKeysAfterDeletion []commonDataForAccessKey
				mockListAccessKeys(m, append(listOfKeysAfterDeletion, resourceData.secondKey)).Once()
			},
			mockData: resourceMock,
			steps: []resource.TestStep{
				{
					Config: testutils.LoadFixtureString(t, "testdata/TestResAccessKey/create_2_versions_with_timeouts.tf"),
					Check: baseAccessKeyChecker.
						CheckEqual("authentication_method", "AWS4_HMAC_SHA256").
						CheckEqual("primary_guid", "asde-efdr-reded").
						CheckEqualBatch("credentials_a.", credentialsABatch).
						CheckEqualBatch("credentials_b.", credentialsBBatch).
						Build(),
				},
				{
					Config:      testutils.LoadFixtureString(t, "testdata/TestResAccessKey/create_with_timeout.tf"),
					ExpectError: regexp.MustCompile("Error: deletion terminated"),
				},
			},
		},
		"fail of deletion - version assigned to property": {
			init: func(m *cloudaccess.Mock, resourceData commonDataForResource) {
				mockCreationAccessKeyWith2Versions(m, resourceData)
				mockReadAccessKey(m, resourceData, twoElementsVersionList)
				mockReadAccessKey(m, resourceData, twoElementsVersionList)
				//delete 2nd version - fail
				mockListAccessKeyVersions(m, resourceData.defaultKey, twoElementsVersionList).Once()
				mockLookupsProperties(m, resourceData.defaultKey, secondAccessKeyVersion).Once()

				//delete all versions
				mockListAccessKeyVersions(m, resourceData.defaultKey, twoElementsVersionList).Once()
				mockLookupsPropertiesNoProperties(m, resourceData.propertyData, secondAccessKeyVersion).Once()
				mockDeleteAccessKeyVersion(m, resourceData.defaultKey, secondAccessKeyVersion).Once()
				mockGetAccessKeyVersion(m, resourceData.defaultKey, cloudaccess.PendingDeletion, secondAccessKeyVersion).Once()
				mockListAccessKeyVersions(m, resourceData.defaultKey, twoElementsVersionList).Once()
				mockListAccessKeyVersionsOnly1Version(m, resourceData.defaultKey).Once()
				mockLookupsPropertiesNoProperties(m, resourceData.propertyData, firstAccessKeyVersion).Once()
				mockDeleteAccessKeyVersion(m, resourceData.defaultKey, firstAccessKeyVersion).Once()
				mockGetAccessKeyVersion(m, resourceData.defaultKey, cloudaccess.PendingDeletion, firstAccessKeyVersion).Once()
				mockListAccessKeyVersionsOnly1Version(m, resourceData.defaultKey).Once()
				mockListAccessKeyVersions(m, resourceData.defaultKey, emptyVersionList).Twice()

				// delete key with no versions
				mockDeleteAccessKey(m, resourceData.defaultKey).Once()
				var listOfKeysAfterDeletion []commonDataForAccessKey
				mockListAccessKeys(m, append(listOfKeysAfterDeletion, resourceData.secondKey)).Once()

			},
			mockData: resourceMock,
			steps: []resource.TestStep{
				{
					Config: testutils.LoadFixtureString(t, "testdata/TestResAccessKey/create_2_versions.tf"),
					Check: baseAccessKeyChecker.
						CheckEqual("authentication_method", "AWS4_HMAC_SHA256").
						CheckEqual("primary_guid", "asde-efdr-reded").
						CheckEqualBatch("credentials_a.", credentialsABatch).
						CheckEqualBatch("credentials_b.", credentialsBBatch).
						Build(),
				},
				{
					Config:      testutils.LoadFixtureString(t, "testdata/TestResAccessKey/create.tf"),
					ExpectError: regexp.MustCompile(fmt.Sprintf("cannot delete version: %d of access key %d assigned to property", secondAccessKeyVersion, 12345)),
				},
			},
		},
		"fail on creation - tainted resource": {
			init: func(m *cloudaccess.Mock, resourceData commonDataForResource) {
				mockCreateAccessKey(m, resourceData.defaultKey).Once()
				mockGetAccessKeyStatus(m, 12345, resourceData.defaultKey).Once()
				mockGetAccessKeyVersion(m, resourceData.defaultKey, cloudaccess.Active, firstAccessKeyVersion).Once()
				// fail and taint resource
				m.On("CreateAccessKeyVersion", testutils.MockContext, cloudaccess.CreateAccessKeyVersionRequest{
					AccessKeyUID: resourceData.defaultKey.accessKeyUID,
					Body: cloudaccess.CreateAccessKeyVersionRequestBody{
						CloudAccessKeyID:     resourceData.defaultKey.credentialsB.cloudAccessKeyID,
						CloudSecretAccessKey: resourceData.defaultKey.credentialsB.cloudSecretAccessKey,
					}}).Return(nil, cloudaccess.ErrCreateAccessKeyVersion).Once()
				//Delete before replace
				mockReadAccessKeyWith1Version(m, resourceData)
				mockDeletionAccessKeyWith1Version(m, resourceData)

				//Second successful creation attempt
				mockCreationAccessKeyWith2Versions(m, resourceData)
				mockReadAccessKey(m, resourceData, twoElementsVersionList)
				mockDeletionAccessKeyWith2Versions(m, resourceData)

			},
			mockData: resourceMock,
			steps: []resource.TestStep{
				{
					Config:      testutils.LoadFixtureString(t, "testdata/TestResAccessKey/create_2_versions.tf"),
					ExpectError: regexp.MustCompile("Error: create access key version failed"),
				},
				{
					Config: testutils.LoadFixtureString(t, "testdata/TestResAccessKey/create_2_versions.tf"),
					Taint:  []string{"akamai_cloudaccess_key.test"},
					Check: baseAccessKeyChecker.
						CheckEqual("authentication_method", "AWS4_HMAC_SHA256").
						CheckEqual("primary_guid", "asde-efdr-reded").
						CheckEqualBatch("credentials_a.", credentialsABatch).
						CheckEqualBatch("credentials_b.", credentialsBBatch).
						Build(),
				},
			},
		},
		"fail on creation key - processing status failed": {
			init: func(m *cloudaccess.Mock, resourceData commonDataForResource) {
				mockCreateAccessKey(m, resourceData.defaultKey).Once()
				// access key creation fail
				m.On("GetAccessKeyStatus", testutils.MockContext, cloudaccess.GetAccessKeyStatusRequest{RequestID: 12345}).
					Return(&cloudaccess.GetAccessKeyStatusResponse{
						AccessKey: &cloudaccess.KeyLink{
							AccessKeyUID: resourceData.defaultKey.accessKeyUID,
						},
						AccessKeyVersion: &cloudaccess.KeyVersion{
							AccessKeyUID: resourceData.defaultKey.accessKeyUID,
							Version:      firstAccessKeyVersion,
						},
						ProcessingStatus: cloudaccess.ProcessingFailed,
						Request: &cloudaccess.RequestInformation{
							AccessKeyName:        resourceData.defaultKey.accessKeyName,
							AuthenticationMethod: cloudaccess.AuthType(resourceData.defaultKey.authenticationMethod),
							ContractID:           resourceData.defaultKey.contractID,
							GroupID:              resourceData.defaultKey.groupID,
							NetworkConfiguration: &cloudaccess.SecureNetwork{
								SecurityNetwork: cloudaccess.NetworkType(resourceData.defaultKey.networkConfig.securityNetwork),
								AdditionalCDN:   ptr.To(cloudaccess.CDNType(resourceData.defaultKey.networkConfig.additionalCDN)),
							},
						},
						RequestDate: time.Date(2024, 1, 10, 11, 9, 10, 67708, time.UTC),
						RequestID:   12345,
						RequestedBy: "dev-user",
					}, nil)
			},
			mockData: resourceMock,
			steps: []resource.TestStep{
				{
					Config:      testutils.LoadFixtureString(t, "testdata/TestResAccessKey/create_2_versions.tf"),
					ExpectError: regexp.MustCompile("Error: access key creation failed"),
				},
			},
		},
		"fail on creation key version - processing status failed": {
			init: func(m *cloudaccess.Mock, resourceData commonDataForResource) {
				mockCreateAccessKey(m, resourceData.defaultKey).Once()
				mockGetAccessKeyStatus(m, 12345, resourceData.defaultKey).Once()
				mockGetAccessKeyVersion(m, resourceData.defaultKey, cloudaccess.Active, firstAccessKeyVersion).Once()
				mockCreateAccessKeyVersion(m, resourceData.defaultKey).Once()
				// access key version creation fail
				m.On("GetAccessKeyVersionStatus", testutils.MockContext, cloudaccess.GetAccessKeyVersionStatusRequest{RequestID: 124}).
					Return(&cloudaccess.GetAccessKeyVersionStatusResponse{
						AccessKeyVersion: &cloudaccess.KeyVersion{
							AccessKeyUID: resourceData.defaultKey.accessKeyUID,
							Version:      secondAccessKeyVersion,
						},
						ProcessingStatus: cloudaccess.ProcessingFailed,
						RequestDate:      time.Date(2024, 1, 10, 11, 9, 10, 67708, time.UTC),
						RequestedBy:      "dev-user",
					}, nil)
				// delete 1 version
				mockDeletionAccessKeyWith1Version(m, resourceData)
			},
			mockData: resourceMock,
			steps: []resource.TestStep{
				{
					Config:      testutils.LoadFixtureString(t, "testdata/TestResAccessKey/create_2_versions.tf"),
					ExpectError: regexp.MustCompile("Error: access key version creation failed"),
				},
			},
		},
		"fail on creation - not proper additional cdn": {
			steps: []resource.TestStep{
				{
					Config:      testutils.LoadFixtureString(t, "testdata/TestResAccessKey/wrong_additional_cdn.tf"),
					ExpectError: regexp.MustCompile(`Attribute network_configuration.additional_cdn value must be one of:\s*\["CHINA_CDN" "RUSSIA_CDN"], got: "TEST"`),
				},
			},
		},
		"fail on creation - not proper security network": {
			steps: []resource.TestStep{
				{
					Config:      testutils.LoadFixtureString(t, "testdata/TestResAccessKey/wrong_security_network.tf"),
					ExpectError: regexp.MustCompile(`Attribute network_configuration.security_network value must be one of:\s*\["STANDARD_TLS" "ENHANCED_TLS"], got: "TEST"`),
				},
			},
		},
		"fail on creation - not proper authentication method": {
			steps: []resource.TestStep{
				{
					Config:      testutils.LoadFixtureString(t, "testdata/TestResAccessKey/wrong_authentication_method.tf"),
					ExpectError: regexp.MustCompile(`Attribute authentication_method value must be one of: \["AWS4_HMAC_SHA256"\s*"GOOG4_HMAC_SHA256" "AOS4_HMAC_SHA256" "AVM_CLOUDINARY" "VP_QUEUE_IT"\], got:\s*"TEST"`),
				},
			},
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
					Steps:                    test.steps,
				})
			})
			client.AssertExpectations(t)
		})
	}
}
func mockDeletionAccessKeyWith2Versions(m *cloudaccess.Mock, resourceData commonDataForResource) {
	mockListAccessKeyVersions(m, resourceData.defaultKey, twoElementsVersionList).Once()
	mockLookupsPropertiesNoProperties(m, resourceData.propertyData, firstAccessKeyVersion).Once()
	mockLookupsPropertiesNoProperties(m, resourceData.propertyData, secondAccessKeyVersion).Once()
	mockDeleteAccessKeyVersion(m, resourceData.defaultKey, firstAccessKeyVersion).Once()
	mockGetAccessKeyVersion(m, resourceData.defaultKey, cloudaccess.PendingDeletion, firstAccessKeyVersion).Once()
	mockListAccessKeyVersions(m, resourceData.defaultKey, oneElementVersionList).Once()
	mockDeleteAccessKeyVersion(m, resourceData.defaultKey, secondAccessKeyVersion).Once()
	mockGetAccessKeyVersion(m, resourceData.defaultKey, cloudaccess.PendingDeletion, secondAccessKeyVersion).Once()
	mockListAccessKeyVersions(m, resourceData.defaultKey, emptyVersionList).Once()
	mockDeleteAccessKey(m, resourceData.defaultKey).Once()
	var listOfKeysAfterDeletion []commonDataForAccessKey
	mockListAccessKeys(m, append(listOfKeysAfterDeletion, resourceData.secondKey)).Once()
}

func mockCreationAccessKeyWith2Versions(m *cloudaccess.Mock, resourceData commonDataForResource) {
	mockCreateAccessKey(m, resourceData.defaultKey).Once()
	mockGetAccessKeyStatus(m, 12345, resourceData.defaultKey).Once()
	mockGetAccessKeyVersion(m, resourceData.defaultKey, cloudaccess.Active, firstAccessKeyVersion).Once()
	mockCreateAccessKeyVersion(m, resourceData.defaultKey).Once()
	mockGetAccessKeyVersionStatus(m, resourceData.defaultKey, 124, secondAccessKeyVersion).Once()
	mockGetAccessKeyVersion(m, resourceData.defaultKey, cloudaccess.Active, secondAccessKeyVersion).Once()
}

func mockDeletionAccessKeyWith1Version(m *cloudaccess.Mock, resourceData commonDataForResource) {
	mockListAccessKeyVersionsOnly1Version(m, resourceData.defaultKey).Once()
	mockLookupsPropertiesNoProperties(m, resourceData.propertyData, firstAccessKeyVersion).Once()
	mockDeleteAccessKeyVersion(m, resourceData.defaultKey, firstAccessKeyVersion).Once()
	mockGetAccessKeyVersion(m, resourceData.defaultKey, cloudaccess.PendingDeletion, firstAccessKeyVersion).Once()
	mockListAccessKeyVersions(m, resourceData.defaultKey, emptyVersionList).Once()
	mockDeleteAccessKey(m, resourceData.defaultKey).Once()
	var listOfKeysAfterDeletion []commonDataForAccessKey
	mockListAccessKeys(m, append(listOfKeysAfterDeletion, resourceData.secondKey)).Once()
}

func mockReadAccessKey(m *cloudaccess.Mock, resourceData commonDataForResource, size int) {
	mockGetAccessKey(m, resourceData.defaultKey).Once()
	mockListAccessKeyVersions(m, resourceData.defaultKey, size).Once()
}

func mockReadAccessKeyWith1Version(m *cloudaccess.Mock, resourceData commonDataForResource) {
	mockGetAccessKey(m, resourceData.defaultKey).Once()
	mockListAccessKeyVersionsOnly1Version(m, resourceData.defaultKey).Once()
}

func mockCreationAccessKeyWith1Version(m *cloudaccess.Mock, resourceData commonDataForResource) {
	mockCreateAccessKey(m, resourceData.defaultKey).Once()
	mockGetAccessKeyStatus(m, 12345, resourceData.defaultKey).Once()
	mockGetAccessKeyVersion(m, resourceData.defaultKey, cloudaccess.Active, firstAccessKeyVersion).Once()
}

func mockCreationAccessKeyUsingCredB(m *cloudaccess.Mock, resourceData commonDataForResource) {
	mockCreateAccessKeyUsingCredB(m, resourceData.credBOnlyKey).Once()
	mockGetAccessKeyStatus(m, 12345, resourceData.credBOnlyKey).Once()
	m.On("GetAccessKeyVersion", testutils.MockContext, cloudaccess.GetAccessKeyVersionRequest{AccessKeyUID: resourceData.credBOnlyKey.accessKeyUID, Version: firstAccessKeyVersion}).
		Return(&cloudaccess.GetAccessKeyVersionResponse{
			AccessKeyUID:     resourceData.credBOnlyKey.accessKeyUID,
			CloudAccessKeyID: ptr.To(resourceData.credBOnlyKey.credentialsB.cloudAccessKeyID),
			CreatedBy:        "dev-user",
			CreatedTime:      time.Date(2024, 1, 10, 11, 9, 10, 67708, time.UTC),
			DeploymentStatus: cloudaccess.Active,
			Version:          firstAccessKeyVersion,
			VersionGUID:      "asde-efdr-reded",
		}, nil).Once()
}

func mockGetAccessKey(client *cloudaccess.Mock, testData commonDataForAccessKey) *mock.Call {
	resp := &cloudaccess.GetAccessKeyResponse{
		AccessKeyName:        testData.accessKeyName,
		AccessKeyUID:         testData.accessKeyUID,
		AuthenticationMethod: testData.authenticationMethod,
		NetworkConfiguration: &cloudaccess.SecureNetwork{
			SecurityNetwork: cloudaccess.NetworkType(testData.networkConfig.securityNetwork),
		},
		LatestVersion: firstAccessKeyVersion,
		Groups: []cloudaccess.Group{
			{
				GroupID:     testData.groupID,
				GroupName:   ptr.To("random group name"),
				ContractIDs: []string{testData.contractID},
			},
		},
		CreatedBy:   "dev-user",
		CreatedTime: time.Date(2024, 1, 10, 11, 9, 10, 67708, time.UTC),
	}
	if testData.networkConfig.additionalCDN != "" {
		resp.NetworkConfiguration.AdditionalCDN = ptr.To(cloudaccess.CDNType(testData.networkConfig.additionalCDN))
	}
	return client.On("GetAccessKey", testutils.MockContext, cloudaccess.AccessKeyRequest{AccessKeyUID: testData.accessKeyUID}).
		Return(resp, nil)
}
func mockGetAccessKeyWithSpecificNameAndVersion(m *cloudaccess.Mock, testData commonDataForAccessKey, name string, version int64) *mock.Call {
	resp := &cloudaccess.GetAccessKeyResponse{
		AccessKeyName:        name,
		AccessKeyUID:         testData.accessKeyUID,
		AuthenticationMethod: testData.authenticationMethod,
		NetworkConfiguration: &cloudaccess.SecureNetwork{
			SecurityNetwork: cloudaccess.NetworkType(testData.networkConfig.securityNetwork),
		},
		LatestVersion: version,
		Groups: []cloudaccess.Group{
			{
				GroupID:     testData.groupID,
				GroupName:   ptr.To("random group name"),
				ContractIDs: []string{testData.contractID},
			},
		},
		CreatedBy:   "dev-user",
		CreatedTime: time.Date(2024, 1, 10, 11, 9, 10, 67708, time.UTC),
	}
	if testData.networkConfig.additionalCDN != "" {
		resp.NetworkConfiguration.AdditionalCDN = ptr.To(cloudaccess.CDNType(testData.networkConfig.additionalCDN))
	}
	return m.On("GetAccessKey", testutils.MockContext, cloudaccess.AccessKeyRequest{AccessKeyUID: testData.accessKeyUID}).
		Return(resp, nil)
}

func mockGetAccessKeyNotFound(client *cloudaccess.Mock, testData commonDataForAccessKey) *mock.Call {
	return client.On("GetAccessKey", testutils.MockContext, cloudaccess.AccessKeyRequest{AccessKeyUID: testData.accessKeyUID}).
		Return(nil, cloudaccess.ErrAccessKeyNotFound)
}

func mockGetAccessKeyStatus(client *cloudaccess.Mock, requestID int64, testData commonDataForAccessKey) *mock.Call {
	return client.On("GetAccessKeyStatus", testutils.MockContext, cloudaccess.GetAccessKeyStatusRequest{RequestID: requestID}).
		Return(&cloudaccess.GetAccessKeyStatusResponse{
			AccessKey: &cloudaccess.KeyLink{
				AccessKeyUID: testData.accessKeyUID,
			},
			AccessKeyVersion: &cloudaccess.KeyVersion{
				AccessKeyUID: testData.accessKeyUID,
				Version:      firstAccessKeyVersion,
			},
			ProcessingStatus: cloudaccess.ProcessingDone,
			Request: &cloudaccess.RequestInformation{
				AccessKeyName:        testData.accessKeyName,
				AuthenticationMethod: cloudaccess.AuthType(testData.authenticationMethod),
				ContractID:           testData.contractID,
				GroupID:              testData.groupID,
				NetworkConfiguration: &cloudaccess.SecureNetwork{
					SecurityNetwork: cloudaccess.NetworkType(testData.networkConfig.securityNetwork),
					AdditionalCDN:   ptr.To(cloudaccess.CDNType(testData.networkConfig.additionalCDN)),
				},
			},
			RequestDate: time.Date(2024, 1, 10, 11, 9, 10, 67708, time.UTC),
			RequestID:   requestID,
			RequestedBy: "dev-user",
		}, nil)
}

func mockCreateAccessKey(client *cloudaccess.Mock, testData commonDataForAccessKey) *mock.Call {
	req := cloudaccess.CreateAccessKeyRequest{
		AccessKeyName:        testData.accessKeyName,
		AuthenticationMethod: testData.authenticationMethod,
		ContractID:           testData.contractID,
		GroupID:              testData.groupID,
		NetworkConfiguration: cloudaccess.SecureNetwork{
			SecurityNetwork: cloudaccess.NetworkType(testData.networkConfig.securityNetwork),
		},
		Credentials: cloudaccess.Credentials{
			CloudAccessKeyID:     testData.credentialsA.cloudAccessKeyID,
			CloudSecretAccessKey: testData.credentialsA.cloudSecretAccessKey,
		},
	}
	if testData.networkConfig.additionalCDN != "" {
		req.NetworkConfiguration.AdditionalCDN = ptr.To(cloudaccess.CDNType(testData.networkConfig.additionalCDN))
	}
	return client.On("CreateAccessKey", testutils.MockContext, req).Return(&cloudaccess.CreateAccessKeyResponse{RequestID: 12345, RetryAfter: 1000}, nil)
}

func mockCreateAccessKeyUsingCredB(client *cloudaccess.Mock, testData commonDataForAccessKey) *mock.Call {
	req := cloudaccess.CreateAccessKeyRequest{
		AccessKeyName:        testData.accessKeyName,
		AuthenticationMethod: testData.authenticationMethod,
		ContractID:           testData.contractID,
		GroupID:              testData.groupID,
		NetworkConfiguration: cloudaccess.SecureNetwork{
			SecurityNetwork: cloudaccess.NetworkType(testData.networkConfig.securityNetwork),
		},
		Credentials: cloudaccess.Credentials{
			CloudAccessKeyID:     testData.credentialsB.cloudAccessKeyID,
			CloudSecretAccessKey: testData.credentialsB.cloudSecretAccessKey,
		},
	}
	if testData.networkConfig.additionalCDN != "" {
		req.NetworkConfiguration.AdditionalCDN = ptr.To(cloudaccess.CDNType(testData.networkConfig.additionalCDN))
	}
	return client.On("CreateAccessKey", testutils.MockContext, req).Return(&cloudaccess.CreateAccessKeyResponse{RequestID: 12345, RetryAfter: 1000}, nil)
}

func mockUpdateAccessKey(client *cloudaccess.Mock, testData commonDataForAccessKey, updatedName string) *mock.Call {
	return client.On("UpdateAccessKey", testutils.MockContext, cloudaccess.UpdateAccessKeyRequest{
		AccessKeyName: updatedName,
	}, cloudaccess.AccessKeyRequest{AccessKeyUID: testData.accessKeyUID},
	).Return(&cloudaccess.UpdateAccessKeyResponse{
		AccessKeyUID:  testData.accessKeyUID,
		AccessKeyName: updatedName,
	}, nil)
}

func mockDeleteAccessKey(client *cloudaccess.Mock, testData commonDataForAccessKey) *mock.Call {
	return client.On("DeleteAccessKey", testutils.MockContext, cloudaccess.AccessKeyRequest{
		AccessKeyUID: testData.accessKeyUID,
	},
	).Return(nil)
}

func mockListAccessKeys(client *cloudaccess.Mock, testData []commonDataForAccessKey) *mock.Call {
	var listResponse cloudaccess.ListAccessKeysResponse
	if len(testData) == 2 {
		listResponse = cloudaccess.ListAccessKeysResponse{
			AccessKeys: []cloudaccess.AccessKeyResponse{{
				AccessKeyName:        testData[0].accessKeyName,
				AccessKeyUID:         testData[0].accessKeyUID,
				AuthenticationMethod: testData[0].authenticationMethod,
				NetworkConfiguration: &cloudaccess.SecureNetwork{
					SecurityNetwork: cloudaccess.NetworkType(testData[0].networkConfig.securityNetwork),
					AdditionalCDN:   ptr.To(cloudaccess.CDNType(testData[0].networkConfig.additionalCDN)),
				},
				LatestVersion: firstAccessKeyVersion,
				Groups: []cloudaccess.Group{
					{
						GroupID:     testData[0].groupID,
						GroupName:   ptr.To("random_group_name"),
						ContractIDs: []string{testData[0].contractID},
					},
				},
				CreatedBy:   "dev-user",
				CreatedTime: time.Date(2024, 1, 10, 11, 9, 10, 67708, time.UTC),
			}, {
				AccessKeyName:        testData[1].accessKeyName,
				AccessKeyUID:         testData[1].accessKeyUID,
				AuthenticationMethod: testData[2].authenticationMethod,
				NetworkConfiguration: &cloudaccess.SecureNetwork{
					SecurityNetwork: cloudaccess.NetworkType(testData[1].networkConfig.securityNetwork),
					AdditionalCDN:   ptr.To(cloudaccess.CDNType(testData[1].networkConfig.additionalCDN)),
				},
				LatestVersion: firstAccessKeyVersion,
				Groups: []cloudaccess.Group{
					{
						GroupID:     testData[1].groupID,
						GroupName:   ptr.To("random_group_name"),
						ContractIDs: []string{testData[1].contractID},
					},
				},
				CreatedBy:   "dev-user",
				CreatedTime: time.Date(2024, 1, 10, 11, 9, 10, 67708, time.UTC),
			},
			},
		}
	}
	if len(testData) == 1 {
		listResponse = cloudaccess.ListAccessKeysResponse{
			AccessKeys: []cloudaccess.AccessKeyResponse{{
				AccessKeyName:        testData[0].accessKeyName,
				AccessKeyUID:         testData[0].accessKeyUID,
				AuthenticationMethod: testData[0].authenticationMethod,
				NetworkConfiguration: &cloudaccess.SecureNetwork{
					SecurityNetwork: cloudaccess.NetworkType(testData[0].networkConfig.securityNetwork),
					AdditionalCDN:   ptr.To(cloudaccess.CDNType(testData[0].networkConfig.additionalCDN)),
				},
				LatestVersion: firstAccessKeyVersion,
				Groups: []cloudaccess.Group{
					{
						GroupID:     testData[0].groupID,
						GroupName:   ptr.To("random_group_name"),
						ContractIDs: []string{testData[0].contractID},
					},
				},
				CreatedBy:   "dev-user",
				CreatedTime: time.Date(2024, 1, 10, 11, 9, 10, 67708, time.UTC),
			},
			},
		}
	}
	return client.On("ListAccessKeys", testutils.MockContext, cloudaccess.ListAccessKeysRequest{}).
		Return(&listResponse, nil)
}

func mockGetAccessKeyVersion(client *cloudaccess.Mock, testData commonDataForAccessKey, deploymentStatus cloudaccess.DeploymentStatus, version int64) *mock.Call {
	var cloudAccessKeyID, versionGUID string
	if version == firstAccessKeyVersion {
		cloudAccessKeyID = testData.credentialsA.cloudAccessKeyID
		versionGUID = "asde-efdr-reded"
	}
	if version == secondAccessKeyVersion {
		cloudAccessKeyID = testData.credentialsB.cloudAccessKeyID
		versionGUID = "asdd-ads-dasdas"
	}
	if version == thirdAccessKeyVersion {
		cloudAccessKeyID = "test_key_id_3"
		versionGUID = "ffff_eeee-ffffddd"
	}
	return client.On("GetAccessKeyVersion", testutils.MockContext, cloudaccess.GetAccessKeyVersionRequest{AccessKeyUID: testData.accessKeyUID, Version: version}).
		Return(&cloudaccess.GetAccessKeyVersionResponse{
			AccessKeyUID:     testData.accessKeyUID,
			CloudAccessKeyID: ptr.To(cloudAccessKeyID),
			CreatedBy:        "dev-user",
			CreatedTime:      time.Date(2024, 1, 10, 11, 9, 10, 67708, time.UTC),
			DeploymentStatus: deploymentStatus,
			Version:          version,
			VersionGUID:      versionGUID,
		}, nil)
}

func mockDeleteAccessKeyVersion(client *cloudaccess.Mock, testData commonDataForAccessKey, version int64) *mock.Call {
	var cloudAccessKeyID, versionGUID string
	if version == firstAccessKeyVersion {
		cloudAccessKeyID = testData.credentialsA.cloudAccessKeyID
		versionGUID = "asde-efdr-reded"
	}
	if version == secondAccessKeyVersion {
		cloudAccessKeyID = testData.credentialsB.cloudAccessKeyID
		versionGUID = "asdd-ads-dasdas"
	}
	return client.On("DeleteAccessKeyVersion", testutils.MockContext, cloudaccess.DeleteAccessKeyVersionRequest{AccessKeyUID: testData.accessKeyUID, Version: version}).
		Return(&cloudaccess.DeleteAccessKeyVersionResponse{
			AccessKeyUID:     testData.accessKeyUID,
			CloudAccessKeyID: ptr.To(cloudAccessKeyID),
			CreatedBy:        "dev-user",
			CreatedTime:      time.Date(2024, 1, 10, 11, 9, 10, 67708, time.UTC),
			DeploymentStatus: cloudaccess.Active,
			Version:          version,
			VersionGUID:      versionGUID,
		}, nil)
}

func mockGetAccessKeyVersionStatus(client *cloudaccess.Mock, testData commonDataForAccessKey, requestID int64, version int64) *mock.Call {
	return client.On("GetAccessKeyVersionStatus", testutils.MockContext, cloudaccess.GetAccessKeyVersionStatusRequest{RequestID: requestID}).
		Return(&cloudaccess.GetAccessKeyVersionStatusResponse{
			AccessKeyVersion: &cloudaccess.KeyVersion{
				AccessKeyUID: testData.accessKeyUID,
				Version:      version,
			},
			ProcessingStatus: cloudaccess.ProcessingDone,
			RequestDate:      time.Date(2024, 1, 10, 11, 9, 10, 67708, time.UTC),
			RequestedBy:      "dev-user",
		}, nil)
}

func mockCreateAccessKeyVersion(client *cloudaccess.Mock, testData commonDataForAccessKey) *mock.Call {
	return client.On("CreateAccessKeyVersion", testutils.MockContext, cloudaccess.CreateAccessKeyVersionRequest{
		AccessKeyUID: testData.accessKeyUID,
		Body: cloudaccess.CreateAccessKeyVersionRequestBody{
			CloudAccessKeyID:     testData.credentialsB.cloudAccessKeyID,
			CloudSecretAccessKey: testData.credentialsB.cloudSecretAccessKey,
		}}).Return(&cloudaccess.CreateAccessKeyVersionResponse{RequestID: 124, RetryAfter: 1000}, nil)
}

func mockListAccessKeyVersions(client *cloudaccess.Mock, testData commonDataForAccessKey, size int) *mock.Call {
	var listAccessKeyVersionResp cloudaccess.ListAccessKeyVersionsResponse
	if size == twoElementsVersionList {
		listAccessKeyVersionResp = cloudaccess.ListAccessKeyVersionsResponse{AccessKeyVersions: []cloudaccess.AccessKeyVersion{
			{
				AccessKeyUID:     testData.accessKeyUID,
				CloudAccessKeyID: ptr.To("test_key_id"),
				CreatedBy:        "dev-user",
				CreatedTime:      time.Date(2024, 1, 10, 11, 9, 10, 67708, time.UTC),
				DeploymentStatus: cloudaccess.Active,
				Version:          firstAccessKeyVersion,
				VersionGUID:      "asde-efdr-reded",
			},
			{
				AccessKeyUID:     testData.accessKeyUID,
				CloudAccessKeyID: ptr.To("test_key_id_2"),
				CreatedBy:        "dev-user",
				CreatedTime:      time.Date(2024, 1, 10, 11, 9, 10, 67708, time.UTC),
				DeploymentStatus: cloudaccess.Active,
				Version:          secondAccessKeyVersion,
				VersionGUID:      "asdd-ads-dasdas",
			},
		}}
	}
	if size == oneElementVersionList {
		listAccessKeyVersionResp = cloudaccess.ListAccessKeyVersionsResponse{AccessKeyVersions: []cloudaccess.AccessKeyVersion{
			{
				AccessKeyUID:     testData.accessKeyUID,
				CloudAccessKeyID: ptr.To("test_key_id_2"),
				CreatedBy:        "dev-user",
				CreatedTime:      time.Date(2024, 1, 10, 11, 9, 10, 67708, time.UTC),
				DeploymentStatus: cloudaccess.Active,
				Version:          secondAccessKeyVersion,
				VersionGUID:      "asdd-ads-dasdas",
			},
		},
		}
	}
	if size == emptyVersionList {
		listAccessKeyVersionResp = cloudaccess.ListAccessKeyVersionsResponse{AccessKeyVersions: []cloudaccess.AccessKeyVersion{}}
	}
	return client.On("ListAccessKeyVersions", testutils.MockContext, cloudaccess.ListAccessKeyVersionsRequest{
		AccessKeyUID: testData.accessKeyUID,
	}).Return(&listAccessKeyVersionResp, nil)
}

func mockListAccessKeyVersionsOnly1Version(client *cloudaccess.Mock, testData commonDataForAccessKey) *mock.Call {
	var listAccessKeyVersionResp = cloudaccess.ListAccessKeyVersionsResponse{AccessKeyVersions: []cloudaccess.AccessKeyVersion{
		{
			AccessKeyUID:     testData.accessKeyUID,
			CloudAccessKeyID: ptr.To("test_key_id"),
			CreatedBy:        "dev-user",
			CreatedTime:      time.Date(2024, 1, 10, 11, 9, 10, 67708, time.UTC),
			DeploymentStatus: cloudaccess.Active,
			Version:          firstAccessKeyVersion,
			VersionGUID:      "asde-efdr-reded",
		},
	},
	}
	return client.On("ListAccessKeyVersions", testutils.MockContext, cloudaccess.ListAccessKeyVersionsRequest{
		AccessKeyUID: testData.accessKeyUID,
	}).Return(&listAccessKeyVersionResp, nil)
}

func mockLookupsProperties(client *cloudaccess.Mock, testData commonDataForAccessKey, version int64) *mock.Call {
	return client.On("LookupProperties", testutils.MockContext, cloudaccess.LookupPropertiesRequest{
		AccessKeyUID: testData.accessKeyUID,
		Version:      version,
	}).Return(&cloudaccess.LookupPropertiesResponse{Properties: []cloudaccess.Property{
		{
			AccessKeyUID:      testData.accessKeyUID,
			PropertyID:        "123123",
			PropertyName:      "test_property_name",
			ProductionVersion: ptr.To(int64(1)),
			StagingVersion:    ptr.To(int64(1)),
		},
	}}, nil)
}

func mockLookupsPropertiesNoProperties(client *cloudaccess.Mock, testData commonDataForProperty, version int64) *mock.Call {
	lookupPropertiesRes := cloudaccess.LookupPropertiesResponse{Properties: []cloudaccess.Property{}}
	return client.On("LookupProperties", testutils.MockContext, cloudaccess.LookupPropertiesRequest{
		AccessKeyUID: testData.accessKeyUID,
		Version:      version,
	}).Return(&lookupPropertiesRes, nil)
}

func mockListAccessKeyVersionsNoCloudAccessKeyID(m *cloudaccess.Mock, accessKey commonDataForAccessKey, size int) *mock.Call {
	var resp cloudaccess.ListAccessKeyVersionsResponse
	if size >= 1 {
		resp.AccessKeyVersions = append(resp.AccessKeyVersions, cloudaccess.AccessKeyVersion{
			AccessKeyUID:     accessKey.accessKeyUID,
			CloudAccessKeyID: nil,
			CreatedBy:        "dev-user",
			CreatedTime:      time.Date(2024, 1, 10, 11, 9, 10, 67708, time.UTC),
			DeploymentStatus: cloudaccess.Active,
			Version:          firstAccessKeyVersion,
			VersionGUID:      "asde-efdr-reded",
		})
	}
	if size >= 2 {
		resp.AccessKeyVersions = append(resp.AccessKeyVersions, cloudaccess.AccessKeyVersion{
			AccessKeyUID:     accessKey.accessKeyUID,
			CloudAccessKeyID: nil,
			CreatedBy:        "dev-user",
			CreatedTime:      time.Date(2024, 1, 10, 11, 9, 10, 67708, time.UTC),
			DeploymentStatus: cloudaccess.Active,
			Version:          secondAccessKeyVersion,
			VersionGUID:      "asdd-ads-dasdas",
		})
	}
	return m.On("ListAccessKeyVersions", testutils.MockContext, cloudaccess.ListAccessKeyVersionsRequest{
		AccessKeyUID: accessKey.accessKeyUID,
	}).Return(&resp, nil)
}

func mockGetAccessKeyVersionNoCloudAccessKeyID(m *cloudaccess.Mock, accessKey commonDataForAccessKey, status cloudaccess.DeploymentStatus, version int64) *mock.Call {
	var versionGUID string
	if version == firstAccessKeyVersion {
		versionGUID = "asde-efdr-reded"
	}
	if version == secondAccessKeyVersion {
		versionGUID = "asdd-ads-dasdas"
	}
	if version == thirdAccessKeyVersion {
		versionGUID = "ffff_eeee-ffffddd"
	}
	return m.On("GetAccessKeyVersion", testutils.MockContext, cloudaccess.GetAccessKeyVersionRequest{AccessKeyUID: accessKey.accessKeyUID, Version: version}).
		Return(&cloudaccess.GetAccessKeyVersionResponse{
			AccessKeyUID:     accessKey.accessKeyUID,
			CloudAccessKeyID: nil,
			CreatedBy:        "dev-user",
			CreatedTime:      time.Date(2024, 1, 10, 11, 9, 10, 67708, time.UTC),
			DeploymentStatus: status,
			Version:          version,
			VersionGUID:      versionGUID,
		}, nil)
}

func mockDeleteAccessKeyVersionNoCloudAccessKeyID(m *cloudaccess.Mock, accessKey commonDataForAccessKey, version int64) *mock.Call {
	var versionGUID string
	if version == firstAccessKeyVersion {
		versionGUID = "asde-efdr-reded"
	}
	if version == secondAccessKeyVersion {
		versionGUID = "asdd-ads-dasdas"
	}
	if version == thirdAccessKeyVersion {
		versionGUID = "ffff_eeee-ffffddd"
	}
	return m.On("DeleteAccessKeyVersion", testutils.MockContext, cloudaccess.DeleteAccessKeyVersionRequest{AccessKeyUID: accessKey.accessKeyUID, Version: version}).
		Return(&cloudaccess.DeleteAccessKeyVersionResponse{
			AccessKeyUID:     accessKey.accessKeyUID,
			CloudAccessKeyID: nil,
			CreatedBy:        "dev-user",
			CreatedTime:      time.Date(2024, 1, 10, 11, 9, 10, 67708, time.UTC),
			DeploymentStatus: cloudaccess.Active,
			Version:          version,
			VersionGUID:      versionGUID,
		}, nil)
}

func mockCreationNoCloudAccessKeyID1Version(m *cloudaccess.Mock, accessKey commonDataForAccessKey) {
	mockCreateAccessKey(m, accessKey).Once()
	mockGetAccessKeyStatus(m, 12345, accessKey).Once()
	mockGetAccessKeyVersionNoCloudAccessKeyID(m, accessKey, cloudaccess.Active, firstAccessKeyVersion).Once()
}

func mockCreationNoCloudAccessKeyID2Versions(m *cloudaccess.Mock, accessKey commonDataForAccessKey) {
	mockCreationNoCloudAccessKeyID1Version(m, accessKey)
	m.On("CreateAccessKeyVersion", testutils.MockContext, cloudaccess.CreateAccessKeyVersionRequest{
		AccessKeyUID: accessKey.accessKeyUID,
		Body: cloudaccess.CreateAccessKeyVersionRequestBody{
			CloudAccessKeyID:     accessKey.credentialsB.cloudAccessKeyID,
			CloudSecretAccessKey: accessKey.credentialsB.cloudSecretAccessKey,
		}}).Return(&cloudaccess.CreateAccessKeyVersionResponse{RequestID: 124, RetryAfter: 1000}, nil).Once()
	mockGetAccessKeyVersionStatus(m, accessKey, 124, secondAccessKeyVersion).Once()
	mockGetAccessKeyVersionNoCloudAccessKeyID(m, accessKey, cloudaccess.Active, secondAccessKeyVersion).Once()
}

func mockReadAccessKeyNoCloudAccessKeyID(m *cloudaccess.Mock, accessKey commonDataForAccessKey, size int) {
	mockGetAccessKey(m, accessKey).Once()
	mockListAccessKeyVersionsNoCloudAccessKeyID(m, accessKey, size).Once()
}

func mockDeletionNoCloudAccessKeyID1Version(m *cloudaccess.Mock, accessKey commonDataForAccessKey, resourceData commonDataForResource) {
	mockListAccessKeyVersionsNoCloudAccessKeyID(m, accessKey, oneElementVersionList).Once()
	mockLookupsPropertiesNoProperties(m, resourceData.propertyData, firstAccessKeyVersion).Once()
	mockDeleteAccessKeyVersionNoCloudAccessKeyID(m, accessKey, firstAccessKeyVersion).Once()
	mockGetAccessKeyVersionNoCloudAccessKeyID(m, accessKey, cloudaccess.PendingDeletion, firstAccessKeyVersion).Once()
	mockListAccessKeyVersionsNoCloudAccessKeyID(m, accessKey, emptyVersionList).Once()
	mockDeleteAccessKey(m, accessKey).Once()
	mockListAccessKeys(m, []commonDataForAccessKey{resourceData.secondKey}).Once()
}

func mockDeletionNoCloudAccessKeyID2Versions(m *cloudaccess.Mock, accessKey commonDataForAccessKey, resourceData commonDataForResource) {
	mockListAccessKeyVersionsNoCloudAccessKeyID(m, accessKey, twoElementsVersionList).Once()
	mockLookupsPropertiesNoProperties(m, resourceData.propertyData, firstAccessKeyVersion).Once()
	mockLookupsPropertiesNoProperties(m, resourceData.propertyData, secondAccessKeyVersion).Once()
	mockDeleteAccessKeyVersionNoCloudAccessKeyID(m, accessKey, firstAccessKeyVersion).Once()
	mockGetAccessKeyVersionNoCloudAccessKeyID(m, accessKey, cloudaccess.PendingDeletion, firstAccessKeyVersion).Once()
	mockListAccessKeyVersionsOnlyV2NoCloudAccessKeyID(m, accessKey).Once()
	mockDeleteAccessKeyVersionNoCloudAccessKeyID(m, accessKey, secondAccessKeyVersion).Once()
	mockGetAccessKeyVersionNoCloudAccessKeyID(m, accessKey, cloudaccess.PendingDeletion, secondAccessKeyVersion).Once()
	mockListAccessKeyVersionsNoCloudAccessKeyID(m, accessKey, emptyVersionList).Once()
	mockDeleteAccessKey(m, accessKey).Once()
	mockListAccessKeys(m, []commonDataForAccessKey{resourceData.secondKey}).Once()
}

func mockListAccessKeyVersionsV3AndV2NoCloudAccessKeyID(m *cloudaccess.Mock, accessKey commonDataForAccessKey) *mock.Call {
	resp := cloudaccess.ListAccessKeyVersionsResponse{
		AccessKeyVersions: []cloudaccess.AccessKeyVersion{
			{
				AccessKeyUID:     accessKey.accessKeyUID,
				CloudAccessKeyID: nil,
				CreatedBy:        "dev-user",
				CreatedTime:      time.Date(2024, 1, 10, 11, 9, 10, 67708, time.UTC),
				DeploymentStatus: cloudaccess.Active,
				Version:          thirdAccessKeyVersion,
				VersionGUID:      "ffff_eeee-ffffddd",
			},
			{
				AccessKeyUID:     accessKey.accessKeyUID,
				CloudAccessKeyID: nil,
				CreatedBy:        "dev-user",
				CreatedTime:      time.Date(2024, 1, 10, 11, 9, 10, 67708, time.UTC),
				DeploymentStatus: cloudaccess.Active,
				Version:          secondAccessKeyVersion,
				VersionGUID:      "asdd-ads-dasdas",
			},
		},
	}
	return m.On("ListAccessKeyVersions", testutils.MockContext, cloudaccess.ListAccessKeyVersionsRequest{
		AccessKeyUID: accessKey.accessKeyUID,
	}).Return(&resp, nil)
}

func mockDeletionNoCloudAccessKeyIDAfterRotation(m *cloudaccess.Mock, accessKey commonDataForAccessKey, resourceData commonDataForResource) {
	mockListAccessKeyVersionsV3AndV2NoCloudAccessKeyID(m, accessKey).Once()
	mockLookupsPropertiesNoProperties(m, resourceData.propertyData, thirdAccessKeyVersion).Once()
	mockLookupsPropertiesNoProperties(m, resourceData.propertyData, secondAccessKeyVersion).Once()
	mockDeleteAccessKeyVersionNoCloudAccessKeyID(m, accessKey, thirdAccessKeyVersion).Once()
	mockGetAccessKeyVersionNoCloudAccessKeyID(m, accessKey, cloudaccess.PendingDeletion, thirdAccessKeyVersion).Once()
	mockListAccessKeyVersionsV3AndV2NoCloudAccessKeyID(m, accessKey).Once()
	mockListAccessKeyVersionsOnlyV2NoCloudAccessKeyID(m, accessKey).Once()
	mockDeleteAccessKeyVersionNoCloudAccessKeyID(m, accessKey, secondAccessKeyVersion).Once()
	mockGetAccessKeyVersionNoCloudAccessKeyID(m, accessKey, cloudaccess.PendingDeletion, secondAccessKeyVersion).Once()
	mockListAccessKeyVersionsOnlyV2NoCloudAccessKeyID(m, accessKey).Once()
	mockListAccessKeyVersionsNoCloudAccessKeyID(m, accessKey, emptyVersionList).Once()
	mockDeleteAccessKey(m, accessKey).Once()
	mockListAccessKeys(m, []commonDataForAccessKey{resourceData.secondKey}).Once()
}

func mockListAccessKeyVersionsOnlyV2NoCloudAccessKeyID(m *cloudaccess.Mock, accessKey commonDataForAccessKey) *mock.Call {
	resp := cloudaccess.ListAccessKeyVersionsResponse{
		AccessKeyVersions: []cloudaccess.AccessKeyVersion{{
			AccessKeyUID:     accessKey.accessKeyUID,
			CloudAccessKeyID: nil,
			CreatedBy:        "dev-user",
			CreatedTime:      time.Date(2024, 1, 10, 11, 9, 10, 67708, time.UTC),
			DeploymentStatus: cloudaccess.Active,
			Version:          secondAccessKeyVersion,
			VersionGUID:      "asdd-ads-dasdas",
		}},
	}
	return m.On("ListAccessKeyVersions", testutils.MockContext, cloudaccess.ListAccessKeyVersionsRequest{
		AccessKeyUID: accessKey.accessKeyUID,
	}).Return(&resp, nil)
}

func mockCreateAccessKeyVersionNoCloudAccessKeyIDForCredB(m *cloudaccess.Mock, accessKey commonDataForAccessKey) *mock.Call {
	return m.On("CreateAccessKeyVersion", testutils.MockContext, cloudaccess.CreateAccessKeyVersionRequest{
		AccessKeyUID: accessKey.accessKeyUID,
		Body: cloudaccess.CreateAccessKeyVersionRequestBody{
			CloudAccessKeyID:     "",
			CloudSecretAccessKey: accessKey.credentialsB.cloudSecretAccessKey,
		}}).Return(&cloudaccess.CreateAccessKeyVersionResponse{RequestID: 124, RetryAfter: 1000}, nil)
}

func mockListAccessKeyVersionsV1AndV3NoCloudAccessKeyID(m *cloudaccess.Mock, accessKey commonDataForAccessKey) *mock.Call {
	resp := cloudaccess.ListAccessKeyVersionsResponse{
		AccessKeyVersions: []cloudaccess.AccessKeyVersion{
			{
				AccessKeyUID:     accessKey.accessKeyUID,
				CloudAccessKeyID: nil,
				CreatedBy:        "dev-user",
				CreatedTime:      time.Date(2024, 1, 10, 11, 9, 10, 67708, time.UTC),
				DeploymentStatus: cloudaccess.Active,
				Version:          firstAccessKeyVersion,
				VersionGUID:      "asde-efdr-reded",
			},
			{
				AccessKeyUID:     accessKey.accessKeyUID,
				CloudAccessKeyID: nil,
				CreatedBy:        "dev-user",
				CreatedTime:      time.Date(2024, 1, 10, 11, 9, 10, 67708, time.UTC),
				DeploymentStatus: cloudaccess.Active,
				Version:          thirdAccessKeyVersion,
				VersionGUID:      "ffff_eeee-ffffddd",
			},
		},
	}
	return m.On("ListAccessKeyVersions", testutils.MockContext, cloudaccess.ListAccessKeyVersionsRequest{
		AccessKeyUID: accessKey.accessKeyUID,
	}).Return(&resp, nil)
}

func mockListAccessKeyVersionsOnlyV3NoCloudAccessKeyID(m *cloudaccess.Mock, accessKey commonDataForAccessKey) *mock.Call {
	resp := cloudaccess.ListAccessKeyVersionsResponse{
		AccessKeyVersions: []cloudaccess.AccessKeyVersion{{
			AccessKeyUID:     accessKey.accessKeyUID,
			CloudAccessKeyID: nil,
			CreatedBy:        "dev-user",
			CreatedTime:      time.Date(2024, 1, 10, 11, 9, 10, 67708, time.UTC),
			DeploymentStatus: cloudaccess.Active,
			Version:          thirdAccessKeyVersion,
			VersionGUID:      "ffff_eeee-ffffddd",
		}},
	}
	return m.On("ListAccessKeyVersions", testutils.MockContext, cloudaccess.ListAccessKeyVersionsRequest{
		AccessKeyUID: accessKey.accessKeyUID,
	}).Return(&resp, nil)
}

func mockDeletionNoCloudAccessKeyIDAfterCrossRotation(m *cloudaccess.Mock, accessKey commonDataForAccessKey, resourceData commonDataForResource) {
	mockListAccessKeyVersionsV1AndV3NoCloudAccessKeyID(m, accessKey).Once()
	mockLookupsPropertiesNoProperties(m, resourceData.propertyData, firstAccessKeyVersion).Once()
	mockLookupsPropertiesNoProperties(m, resourceData.propertyData, thirdAccessKeyVersion).Once()
	mockDeleteAccessKeyVersionNoCloudAccessKeyID(m, accessKey, firstAccessKeyVersion).Once()
	mockGetAccessKeyVersionNoCloudAccessKeyID(m, accessKey, cloudaccess.PendingDeletion, firstAccessKeyVersion).Once()
	mockListAccessKeyVersionsV1AndV3NoCloudAccessKeyID(m, accessKey).Once()
	mockListAccessKeyVersionsOnlyV3NoCloudAccessKeyID(m, accessKey).Once()
	mockDeleteAccessKeyVersionNoCloudAccessKeyID(m, accessKey, thirdAccessKeyVersion).Once()
	mockGetAccessKeyVersionNoCloudAccessKeyID(m, accessKey, cloudaccess.PendingDeletion, thirdAccessKeyVersion).Once()
	mockListAccessKeyVersionsOnlyV3NoCloudAccessKeyID(m, accessKey).Once()
	mockListAccessKeyVersionsNoCloudAccessKeyID(m, accessKey, emptyVersionList).Once()
	mockDeleteAccessKey(m, accessKey).Once()
	mockListAccessKeys(m, []commonDataForAccessKey{resourceData.secondKey}).Once()
}

func TestAccessKeyResource_ImportState(t *testing.T) {
	t.Parallel()
	pollingInterval = 1 * time.Millisecond
	deleteTimeout = 40 * time.Minute
	updateTimeout = 20 * time.Minute
	activationTimeout = 20 * time.Millisecond
	tests := map[string]struct {
		init     func(*cloudaccess.Mock, commonDataForResource)
		steps    []resource.TestStep
		mockData commonDataForResource
	}{
		"Happy path - 2 credentials": {
			init: func(m *cloudaccess.Mock, resourceData commonDataForResource) {
				// step 1 - create
				mockCreationAccessKeyWith2Versions(m, resourceData)
				mockReadAccessKey(m, resourceData, twoElementsVersionList)

				// step 2 - import
				mockReadAccessKey(m, resourceData, twoElementsVersionList)
				mockReadAccessKey(m, resourceData, twoElementsVersionList)

				mockDeletionAccessKeyWith2Versions(m, resourceData)

			},
			steps: []resource.TestStep{
				{
					Config: testutils.LoadFixtureString(t, "testdata/TestResAccessKey/create_2_versions.tf"),

					Check: test.NewStateChecker("akamai_cloudaccess_key.test").
						CheckEqual("access_key_uid", "12345").
						Build(),
				},
				{
					ImportState:                          true,
					ImportStateVerify:                    true,
					ImportStateId:                        "12345",
					ImportStateVerifyIgnore:              []string{"credentials_a", "credentials_b", "primary_guid"},
					ResourceName:                         "akamai_cloudaccess_key.test",
					ImportStateCheck:                     checkImport(),
					ImportStateVerifyIdentifierAttribute: "access_key_uid",
				},
			},
			mockData: resourceMock,
		},
		"Happy path - 1 credential no cloud access key id": {
			init: func(m *cloudaccess.Mock, resourceData commonDataForResource) {
				accessKey := resourceData.noCloudKeyIDKey
				// step 1 - create
				mockCreationNoCloudAccessKeyID1Version(m, accessKey)
				mockReadAccessKeyNoCloudAccessKeyID(m, accessKey, oneElementVersionList)

				// step 2 - import
				mockReadAccessKeyNoCloudAccessKeyID(m, accessKey, oneElementVersionList)
				mockReadAccessKeyNoCloudAccessKeyID(m, accessKey, oneElementVersionList)

				mockDeletionNoCloudAccessKeyID1Version(m, accessKey, resourceData)
			},
			steps: []resource.TestStep{
				{
					Config: testutils.LoadFixtureString(t, "testdata/TestResAccessKey/create_missing_cloud_access_key.tf"),

					Check: test.NewStateChecker("akamai_cloudaccess_key.test").
						CheckEqual("access_key_uid", "12345").
						Build(),
				},
				{
					ImportState:                          true,
					ImportStateVerify:                    true,
					ImportStateId:                        "12345",
					ImportStateVerifyIgnore:              []string{"credentials_a", "credentials_b", "primary_guid"},
					ResourceName:                         "akamai_cloudaccess_key.test",
					ImportStateCheck:                     checkImportSingleCredentialNoCloudAccessKeyID(),
					ImportStateVerifyIdentifierAttribute: "access_key_uid",
				},
			},
			mockData: resourceMock,
		},
		"Happy path - 2 credentials no cloud access key id": {
			init: func(m *cloudaccess.Mock, resourceData commonDataForResource) {
				accessKey := resourceData.noCloudKeyIDKey
				// step 1 - create
				mockCreationNoCloudAccessKeyID2Versions(m, accessKey)
				mockReadAccessKeyNoCloudAccessKeyID(m, accessKey, twoElementsVersionList)

				// step 2 - import
				mockReadAccessKeyNoCloudAccessKeyID(m, accessKey, twoElementsVersionList)
				mockReadAccessKeyNoCloudAccessKeyID(m, accessKey, twoElementsVersionList)

				mockDeletionNoCloudAccessKeyID2Versions(m, accessKey, resourceData)
			},
			steps: []resource.TestStep{
				{
					Config: testutils.LoadFixtureString(t, "testdata/TestResAccessKey/create_2_versions_missing_cloud_access_key.tf"),

					Check: test.NewStateChecker("akamai_cloudaccess_key.test").
						CheckEqual("access_key_uid", "12345").
						Build(),
				},
				{
					ImportState:                          true,
					ImportStateVerify:                    true,
					ImportStateId:                        "12345",
					ImportStateVerifyIgnore:              []string{"credentials_a", "credentials_b", "primary_guid"},
					ResourceName:                         "akamai_cloudaccess_key.test",
					ImportStateCheck:                     checkImportNoCloudAccessKeyID(),
					ImportStateVerifyIdentifierAttribute: "access_key_uid",
				},
			},
			mockData: resourceMock,
		},
		"Happy path - 1 credential": {
			init: func(m *cloudaccess.Mock, resourceData commonDataForResource) {
				// step 1 - create

				mockCreationAccessKeyWith1Version(m, resourceData)
				mockReadAccessKeyWith1Version(m, resourceData)

				// step 2 import
				mockReadAccessKeyWith1Version(m, resourceData)
				mockReadAccessKeyWith1Version(m, resourceData)

				mockDeletionAccessKeyWith1Version(m, resourceData)

			},
			steps: []resource.TestStep{
				{
					Config: testutils.LoadFixtureString(t, "testdata/TestResAccessKey/create.tf"),

					Check: test.NewStateChecker("akamai_cloudaccess_key.test").
						CheckEqual("access_key_uid", "12345").
						Build(),
				},
				{
					ImportState:                          true,
					ImportStateVerify:                    true,
					ImportStateId:                        "12345",
					ImportStateVerifyIgnore:              []string{"credentials_a", "credentials_b", "primary_guid"},
					ResourceName:                         "akamai_cloudaccess_key.test",
					ImportStateCheck:                     checkImportSingleCredential(),
					ImportStateVerifyIdentifierAttribute: "access_key_uid",
				},
			},
			mockData: resourceMock,
		},
		"error - non-unique cloud access key id": {
			init: func(m *cloudaccess.Mock, resourceData commonDataForResource) {
				mockGetAccessKey(m, resourceData.defaultKey).Once()

				m.On("ListAccessKeyVersions", testutils.MockContext, cloudaccess.ListAccessKeyVersionsRequest{
					AccessKeyUID: resourceData.defaultKey.accessKeyUID,
				}).Return(&cloudaccess.ListAccessKeyVersionsResponse{AccessKeyVersions: []cloudaccess.AccessKeyVersion{
					{
						AccessKeyUID:     resourceData.defaultKey.accessKeyUID,
						CloudAccessKeyID: ptr.To("test_key_id"),
						CreatedBy:        "dev-user",
						CreatedTime:      time.Date(2024, 1, 10, 11, 9, 10, 67708, time.UTC),
						DeploymentStatus: cloudaccess.Active,
						Version:          firstAccessKeyVersion,
						VersionGUID:      "asde-efdr-reded",
					},
					{
						AccessKeyUID:     resourceData.defaultKey.accessKeyUID,
						CloudAccessKeyID: ptr.To("test_key_id"),
						CreatedBy:        "dev-user",
						CreatedTime:      time.Date(2024, 1, 10, 11, 9, 10, 67708, time.UTC),
						DeploymentStatus: cloudaccess.Active,
						Version:          secondAccessKeyVersion,
						VersionGUID:      "asdd-ads-dasdas",
					},
				},
				}, nil)

			},
			mockData: resourceMock,
			steps: []resource.TestStep{
				{
					Config:                               testutils.LoadFixtureString(t, "testdata/TestResAccessKey/creation_non_unique_cloud_key_id.tf"),
					ImportState:                          true,
					ImportStateId:                        "12345",
					ResourceName:                         "akamai_cloudaccess_key.test",
					ImportStateCheck:                     checkImport(),
					ImportStateVerifyIdentifierAttribute: "access_key_uid",
					ImportStatePersist:                   true,
					ExpectError:                          regexp.MustCompile("'cloud_access_key_id' should be unique for each pair of credentials"),
				},
			},
		},
		"error - cannot find access key": {
			init: func(m *cloudaccess.Mock, resourceData commonDataForResource) {
				// step 1 - create
				mockCreationAccessKeyWith1Version(m, resourceData)
				mockReadAccessKeyWith1Version(m, resourceData)

				// step 2 import
				m.On("GetAccessKey", testutils.MockContext, cloudaccess.AccessKeyRequest{AccessKeyUID: 000000}).Return(nil, errors.New("oops")).Times(1)

				mockDeletionAccessKeyWith1Version(m, resourceData)

			},
			steps: []resource.TestStep{
				{
					Config: testutils.LoadFixtureString(t, "testdata/TestResAccessKey/create.tf"),

					Check: test.NewStateChecker("akamai_cloudaccess_key.test").
						CheckEqual("access_key_uid", "12345").
						Build(),
				},
				{
					ImportState:             true,
					ImportStateVerify:       true,
					ImportStateId:           "00000",
					ImportStateVerifyIgnore: []string{"credentials_a", "credentials_b", "primary_guid"},
					ResourceName:            "akamai_cloudaccess_key.test",
					ImportStateCheck:        checkImportSingleCredential(),
					ExpectError:             regexp.MustCompile("Cannot Find Access key"),
				},
			},
			mockData: resourceMock,
		},
		"error - incorrect access key": {
			init: func(m *cloudaccess.Mock, resourceData commonDataForResource) {
				// step 1 - create
				mockCreationAccessKeyWith1Version(m, resourceData)
				mockReadAccessKeyWith1Version(m, resourceData)

				// step 2 import

				mockDeletionAccessKeyWith1Version(m, resourceData)

			},
			steps: []resource.TestStep{
				{
					Config: testutils.LoadFixtureString(t, "testdata/TestResAccessKey/create.tf"),

					Check: test.NewStateChecker("akamai_cloudaccess_key.test").
						CheckEqual("access_key_uid", "12345").
						Build(),
				},
				{
					ImportState:             true,
					ImportStateVerify:       true,
					ImportStateId:           "NaN",
					ImportStateVerifyIgnore: []string{"credentials_a", "credentials_b", "primary_guid"},
					ResourceName:            "akamai_cloudaccess_key.test",
					ImportStateCheck:        checkImportSingleCredential(),
					ExpectError:             regexp.MustCompile("Incorrect ID"),
				},
			},
			mockData: resourceMock,
		},
		"error - reading access key list failed": {
			init: func(m *cloudaccess.Mock, resourceData commonDataForResource) {
				// step 1 - create

				mockCreationAccessKeyWith1Version(m, resourceData)
				mockReadAccessKeyWith1Version(m, resourceData)

				// step 2 import
				mockGetAccessKey(m, resourceData.defaultKey).Once()
				m.On("ListAccessKeyVersions", testutils.MockContext, cloudaccess.ListAccessKeyVersionsRequest{AccessKeyUID: resourceData.propertyData.accessKeyUID}).Return(nil, errors.New("oops")).Times(1)

				mockDeletionAccessKeyWith1Version(m, resourceData)

			},
			steps: []resource.TestStep{
				{
					Config: testutils.LoadFixtureString(t, "testdata/TestResAccessKey/create.tf"),

					Check: test.NewStateChecker("akamai_cloudaccess_key.test").
						CheckEqual("access_key_uid", "12345").
						Build(),
				},
				{
					ImportState:             true,
					ImportStateVerify:       true,
					ImportStateId:           "12345",
					ImportStateVerifyIgnore: []string{"credentials_a", "credentials_b", "primary_guid"},
					ResourceName:            "akamai_cloudaccess_key.test",
					ImportStateCheck:        checkImportSingleCredential(),
					ExpectError:             regexp.MustCompile("Reading Access Key list Failed"),
				},
			},
			mockData: resourceMock,
		},
		"error - cannot find access key - Incorrect groupID": {
			init: func(m *cloudaccess.Mock, resourceData commonDataForResource) {
				// step 1 - create
				mockCreationAccessKeyWith1Version(m, resourceData)
				mockReadAccessKeyWith1Version(m, resourceData)

				mockDeletionAccessKeyWith1Version(m, resourceData)

			},
			steps: []resource.TestStep{
				{
					Config: testutils.LoadFixtureString(t, "testdata/TestResAccessKey/create.tf"),

					Check: test.NewStateChecker("akamai_cloudaccess_key.test").
						CheckEqual("access_key_uid", "12345").
						Build(),
				},
				{
					ImportState:             true,
					ImportStateVerify:       true,
					ImportStateId:           "00000,G434,23",
					ImportStateVerifyIgnore: []string{"credentials_a", "credentials_b", "primary_guid"},
					ResourceName:            "akamai_cloudaccess_key.test",
					ExpectError:             regexp.MustCompile("Incorrect groupID"),
				},
			},
			mockData: resourceMock,
		},
		"error - cannot find access key - Incomplete Access Key Identifier": {
			init: func(m *cloudaccess.Mock, resourceData commonDataForResource) {
				// step 1 - create
				mockCreationAccessKeyWith1Version(m, resourceData)
				mockReadAccessKeyWith1Version(m, resourceData)

				mockDeletionAccessKeyWith1Version(m, resourceData)

			},
			steps: []resource.TestStep{
				{
					Config: testutils.LoadFixtureString(t, "testdata/TestResAccessKey/create.tf"),

					Check: test.NewStateChecker("akamai_cloudaccess_key.test").
						CheckEqual("access_key_uid", "12345").
						Build(),
				},
				{
					ImportState:             true,
					ImportStateVerify:       true,
					ImportStateId:           "00000,434",
					ImportStateVerifyIgnore: []string{"credentials_a", "credentials_b", "primary_guid"},
					ResourceName:            "akamai_cloudaccess_key.test",
					ExpectError:             regexp.MustCompile("Incomplete Access Key Identifier"),
				},
			},
			mockData: resourceMock,
		},
		"error - cannot find access key given groupID and missing contractID": {
			init: func(m *cloudaccess.Mock, resourceData commonDataForResource) {
				// step 1 - create
				mockCreationAccessKeyWith1Version(m, resourceData)
				mockReadAccessKeyWith1Version(m, resourceData)

				mockDeletionAccessKeyWith1Version(m, resourceData)

			},
			steps: []resource.TestStep{
				{
					Config: testutils.LoadFixtureString(t, "testdata/TestResAccessKey/create.tf"),

					Check: test.NewStateChecker("akamai_cloudaccess_key.test").
						CheckEqual("access_key_uid", "12345").
						Build(),
				},
				{
					ImportState:             true,
					ImportStateVerify:       true,
					ImportStateId:           "00000,434,",
					ImportStateVerifyIgnore: []string{"credentials_a", "credentials_b", "primary_guid"},
					ResourceName:            "akamai_cloudaccess_key.test",
					ExpectError:             regexp.MustCompile("Invalid contractID"),
				},
			},
			mockData: resourceMock,
		},
		"error - cannot find access key given contractID and missing groupID": {
			init: func(m *cloudaccess.Mock, resourceData commonDataForResource) {
				// step 1 - create
				mockCreationAccessKeyWith1Version(m, resourceData)
				mockReadAccessKeyWith1Version(m, resourceData)

				mockDeletionAccessKeyWith1Version(m, resourceData)

			},
			steps: []resource.TestStep{
				{
					Config: testutils.LoadFixtureString(t, "testdata/TestResAccessKey/create.tf"),

					Check: test.NewStateChecker("akamai_cloudaccess_key.test").
						CheckEqual("access_key_uid", "12345").
						Build(),
				},
				{
					ImportState:             true,
					ImportStateVerify:       true,
					ImportStateId:           "00000,,434",
					ImportStateVerifyIgnore: []string{"credentials_a", "credentials_b", "primary_guid"},
					ResourceName:            "akamai_cloudaccess_key.test",
					ExpectError:             regexp.MustCompile("Couldn't parse provided groupID, \"\" is invalid"),
				},
			},
			mockData: resourceMock,
		},
		"error - reading access key failed - Invalid groupID and contractID combination": {
			init: func(m *cloudaccess.Mock, resourceData commonDataForResource) {
				// step 1 - create

				mockCreationAccessKeyWith1Version(m, resourceData)
				mockReadAccessKeyWith1Version(m, resourceData)

				// step 2 import
				mockGetAccessKey(m, resourceData.defaultKey).Once()
				mockDeletionAccessKeyWith1Version(m, resourceData)

			},
			steps: []resource.TestStep{
				{
					Config: testutils.LoadFixtureString(t, "testdata/TestResAccessKey/create.tf"),

					Check: test.NewStateChecker("akamai_cloudaccess_key.test").
						CheckEqual("access_key_uid", "12345").
						Build(),
				},
				{
					ImportState:             true,
					ImportStateVerify:       true,
					ImportStateId:           "12345,556,78",
					ImportStateVerifyIgnore: []string{"credentials_a", "credentials_b", "primary_guid"},
					ResourceName:            "akamai_cloudaccess_key.test",
					ExpectError:             regexp.MustCompile("Cannot Find Access key for a given groupID and contractID"),
				},
			},
			mockData: resourceMock,
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
					Steps:                    test.steps,
				})
			})
			client.AssertExpectations(t)
		})
	}
}

func checkImport() resource.ImportStateCheckFunc {
	return func(s []*terraform.InstanceState) error {
		if len(s) == 0 {
			return errors.New("No Instance found")
		}
		if len(s) != 1 {
			return fmt.Errorf("Expected one Instance: %d", len(s))
		}

		state := s[0].Attributes
		attributes := map[string]string{
			"access_key_name":                        "test_key_name",
			"access_key_uid":                         "12345",
			"authentication_method":                  "AWS4_HMAC_SHA256",
			"contract_id":                            "1-CTRACT",
			"group_id":                               "12345",
			"network_configuration.additional_cdn":   "CHINA_CDN",
			"network_configuration.security_network": "ENHANCED_TLS",
			"primary_guid":                           "", // will always be empty
			"credentials_a.cloud_access_key_id":      "test_key_id",
			"credentials_a.cloud_secret_access_key":  "",      // will always be empty
			"credentials_a.primary_key":              "false", // will always be false
			"credentials_a.version":                  "1",
			"credentials_a.version_guid":             "asde-efdr-reded",
			"credentials_b.cloud_access_key_id":      "test_key_id_2",
			"credentials_b.cloud_secret_access_key":  "",      // will always be empty
			"credentials_b.primary_key":              "false", // will always be false
			"credentials_b.version":                  "2",
			"credentials_b.version_guid":             "asdd-ads-dasdas",
		}

		invalidValues := []string{}
		for field, expectedVal := range attributes {
			if state[field] != expectedVal {
				invalidValues = append(invalidValues, fmt.Sprintf("field: %s, got: %s, expected: %s ", field, state[field], expectedVal))
			}
		}

		if len(invalidValues) != 0 {

			return fmt.Errorf("found invalid values: %s", strings.Join(invalidValues, "\n"))
		}

		return nil
	}
}

func checkImportSingleCredential() resource.ImportStateCheckFunc {
	return func(s []*terraform.InstanceState) error {
		if len(s) == 0 {
			return errors.New("No Instance found")
		}
		if len(s) != 1 {
			return fmt.Errorf("Expected one Instance: %d", len(s))
		}

		state := s[0].Attributes
		if _, ok := state["credentials_b.cloud_access_key_id"]; ok {
			return errors.New("Got unexpected second credential")
		}

		attributes := map[string]string{
			"access_key_name":                        "test_key_name",
			"access_key_uid":                         "12345",
			"authentication_method":                  "AWS4_HMAC_SHA256",
			"contract_id":                            "1-CTRACT",
			"group_id":                               "12345",
			"network_configuration.additional_cdn":   "CHINA_CDN",
			"network_configuration.security_network": "ENHANCED_TLS",
			"primary_guid":                           "", // will always be empty
			"credentials_a.cloud_access_key_id":      "test_key_id",
			"credentials_a.cloud_secret_access_key":  "",      // will always be empty
			"credentials_a.primary_key":              "false", // will always be false
			"credentials_a.version":                  "1",
			"credentials_a.version_guid":             "asde-efdr-reded",
		}

		invalidValues := []string{}
		for field, expectedVal := range attributes {
			if state[field] != expectedVal {
				invalidValues = append(invalidValues, fmt.Sprintf("field: %s, got: %s, expected: %s ", field, state[field], expectedVal))
			}
		}

		if len(invalidValues) != 0 {

			return fmt.Errorf("found invalid values: %s", strings.Join(invalidValues, "\n"))
		}

		return nil
	}
}

func checkImportNoCloudAccessKeyID() resource.ImportStateCheckFunc {
	return func(s []*terraform.InstanceState) error {
		if len(s) == 0 {
			return errors.New("No Instance found")
		}
		if len(s) != 1 {
			return fmt.Errorf("Expected one Instance: %d", len(s))
		}

		state := s[0].Attributes
		attributes := map[string]string{
			"access_key_name":                        "test_key_name",
			"access_key_uid":                         "12345",
			"authentication_method":                  "VP_QUEUE_IT",
			"contract_id":                            "1-CTRACT",
			"group_id":                               "12345",
			"network_configuration.security_network": "ENHANCED_TLS",
			"primary_guid":                           "",      // will always be empty
			"credentials_a.cloud_access_key_id":      "",      // nil (not set) for VP_QUEUE_IT
			"credentials_a.cloud_secret_access_key":  "",      // will always be empty
			"credentials_a.primary_key":              "false", // will always be false
			"credentials_a.version":                  "1",
			"credentials_a.version_guid":             "asde-efdr-reded",
			"credentials_b.cloud_access_key_id":      "",      // nil (not set) for VP_QUEUE_IT
			"credentials_b.cloud_secret_access_key":  "",      // will always be empty
			"credentials_b.primary_key":              "false", // will always be false
			"credentials_b.version":                  "2",
			"credentials_b.version_guid":             "asdd-ads-dasdas",
		}

		invalidValues := []string{}
		for field, expectedVal := range attributes {
			if state[field] != expectedVal {
				invalidValues = append(invalidValues, fmt.Sprintf("field: %s, got: %s, expected: %s ", field, state[field], expectedVal))
			}
		}

		if len(invalidValues) != 0 {
			return fmt.Errorf("found invalid values: %s", strings.Join(invalidValues, "\n"))
		}

		return nil
	}
}

func checkImportSingleCredentialNoCloudAccessKeyID() resource.ImportStateCheckFunc {
	return func(s []*terraform.InstanceState) error {
		if len(s) == 0 {
			return errors.New("No Instance found")
		}
		if len(s) != 1 {
			return fmt.Errorf("Expected one Instance: %d", len(s))
		}

		state := s[0].Attributes
		if _, ok := state["credentials_b.version"]; ok {
			return errors.New("Got unexpected second credential")
		}

		attributes := map[string]string{
			"access_key_name":                        "test_key_name",
			"access_key_uid":                         "12345",
			"authentication_method":                  "VP_QUEUE_IT",
			"contract_id":                            "1-CTRACT",
			"group_id":                               "12345",
			"network_configuration.security_network": "ENHANCED_TLS",
			"primary_guid":                           "",      // will always be empty
			"credentials_a.cloud_access_key_id":      "",      // nil (not set) for VP_QUEUE_IT
			"credentials_a.cloud_secret_access_key":  "",      // will always be empty
			"credentials_a.primary_key":              "false", // will always be false
			"credentials_a.version":                  "1",
			"credentials_a.version_guid":             "asde-efdr-reded",
		}

		invalidValues := []string{}
		for field, expectedVal := range attributes {
			if state[field] != expectedVal {
				invalidValues = append(invalidValues, fmt.Sprintf("field: %s, got: %s, expected: %s ", field, state[field], expectedVal))
			}
		}

		if len(invalidValues) != 0 {
			return fmt.Errorf("found invalid values: %s", strings.Join(invalidValues, "\n"))
		}

		return nil
	}
}

func TestChangedOrderOfCredentials(t *testing.T) {
	tests := []struct {
		label      string
		stateCredA *Credentials
		stateCredB *Credentials
		planCredA  *Credentials
		planCredB  *Credentials
		expected   bool
	}{
		{
			label:      "nil oldState.CredentialsA returns false",
			stateCredA: nil,
			stateCredB: &Credentials{
				CloudAccessKeyID:     types.StringValue("key-B"),
				CloudSecretAccessKey: types.StringValue("secret-B"),
			},
			planCredA: &Credentials{
				CloudAccessKeyID:     types.StringValue("key-B"),
				CloudSecretAccessKey: types.StringValue("secret-B"),
			},
			planCredB: &Credentials{
				CloudAccessKeyID:     types.StringValue("key-A"),
				CloudSecretAccessKey: types.StringValue("secret-A"),
			},
			expected: false,
		},
		{
			label: "nil plan.CredentialsA returns false",
			stateCredA: &Credentials{
				CloudAccessKeyID:     types.StringValue("key-A"),
				CloudSecretAccessKey: types.StringValue("secret-A"),
			},
			stateCredB: &Credentials{
				CloudAccessKeyID:     types.StringValue("key-B"),
				CloudSecretAccessKey: types.StringValue("secret-B"),
			},
			planCredA: nil,
			planCredB: &Credentials{
				CloudAccessKeyID:     types.StringValue("key-A"),
				CloudSecretAccessKey: types.StringValue("secret-A"),
			},
			expected: false,
		},
		{
			label: "key IDs present and swapped returns true",
			stateCredA: &Credentials{
				CloudAccessKeyID:     types.StringValue("key-A"),
				CloudSecretAccessKey: types.StringValue("secret-A"),
			},
			stateCredB: &Credentials{
				CloudAccessKeyID:     types.StringValue("key-B"),
				CloudSecretAccessKey: types.StringValue("secret-B"),
			},
			planCredA: &Credentials{
				CloudAccessKeyID:     types.StringValue("key-B"),
				CloudSecretAccessKey: types.StringValue("secret-B"),
			},
			planCredB: &Credentials{
				CloudAccessKeyID:     types.StringValue("key-A"),
				CloudSecretAccessKey: types.StringValue("secret-A"),
			},
			expected: true,
		},
		{
			label: "key IDs present and swapped and new secrets returns true",
			stateCredA: &Credentials{
				CloudAccessKeyID:     types.StringValue("key-A"),
				CloudSecretAccessKey: types.StringValue("secret-A"),
			},
			stateCredB: &Credentials{
				CloudAccessKeyID:     types.StringValue("key-B"),
				CloudSecretAccessKey: types.StringValue("secret-B"),
			},
			planCredA: &Credentials{
				CloudAccessKeyID:     types.StringValue("key-B"),
				CloudSecretAccessKey: types.StringValue("secret-C"),
			},
			planCredB: &Credentials{
				CloudAccessKeyID:     types.StringValue("key-A"),
				CloudSecretAccessKey: types.StringValue("secret-D"),
			},
			expected: true,
		},
		{
			label: "key IDs present and not swapped returns false",
			stateCredA: &Credentials{
				CloudAccessKeyID:     types.StringValue("key-A"),
				CloudSecretAccessKey: types.StringValue("secret-A"),
			},
			stateCredB: &Credentials{
				CloudAccessKeyID:     types.StringValue("key-B"),
				CloudSecretAccessKey: types.StringValue("secret-B"),
			},
			planCredA: &Credentials{
				CloudAccessKeyID:     types.StringValue("key-A"),
				CloudSecretAccessKey: types.StringValue("secret-A"),
			},
			planCredB: &Credentials{
				CloudAccessKeyID:     types.StringValue("key-B"),
				CloudSecretAccessKey: types.StringValue("secret-B"),
			},
			expected: false,
		},
		{
			label: "key IDs null and secrets swapped returns true",
			stateCredA: &Credentials{
				CloudAccessKeyID:     types.StringNull(),
				CloudSecretAccessKey: types.StringValue("secret-A"),
			},
			stateCredB: &Credentials{
				CloudAccessKeyID:     types.StringNull(),
				CloudSecretAccessKey: types.StringValue("secret-B"),
			},
			planCredA: &Credentials{
				CloudAccessKeyID:     types.StringNull(),
				CloudSecretAccessKey: types.StringValue("secret-B"),
			},
			planCredB: &Credentials{
				CloudAccessKeyID:     types.StringNull(),
				CloudSecretAccessKey: types.StringValue("secret-A"),
			},
			expected: true,
		},
		{
			label: "key IDs null and secrets identical returns false",
			stateCredA: &Credentials{
				CloudAccessKeyID:     types.StringNull(),
				CloudSecretAccessKey: types.StringValue("same-secret"),
			},
			stateCredB: &Credentials{
				CloudAccessKeyID:     types.StringNull(),
				CloudSecretAccessKey: types.StringValue("same-secret"),
			},
			planCredA: &Credentials{
				CloudAccessKeyID:     types.StringNull(),
				CloudSecretAccessKey: types.StringValue("same-secret"),
			},
			planCredB: &Credentials{
				CloudAccessKeyID:     types.StringNull(),
				CloudSecretAccessKey: types.StringValue("same-secret"),
			},
			expected: false,
		},
		{
			label: "key IDs null and secrets not swapped returns false",
			stateCredA: &Credentials{
				CloudAccessKeyID:     types.StringNull(),
				CloudSecretAccessKey: types.StringValue("secret-A"),
			},
			stateCredB: &Credentials{
				CloudAccessKeyID:     types.StringNull(),
				CloudSecretAccessKey: types.StringValue("secret-B"),
			},
			planCredA: &Credentials{
				CloudAccessKeyID:     types.StringNull(),
				CloudSecretAccessKey: types.StringValue("secret-A"),
			},
			planCredB: &Credentials{
				CloudAccessKeyID:     types.StringNull(),
				CloudSecretAccessKey: types.StringValue("secret-B"),
			},
			expected: false,
		},
	}

	for _, tc := range tests {
		t.Run(tc.label, func(t *testing.T) {
			result := changedOrderOfCredentials(tc.stateCredA, tc.stateCredB, tc.planCredA, tc.planCredB)
			assert.Equal(t, tc.expected, result)
		})
	}
}
