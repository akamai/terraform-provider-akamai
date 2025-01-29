package cloudwrapper

import (
	"fmt"
	"regexp"
	"strconv"
	"testing"

	"github.com/akamai/AkamaiOPEN-edgegrid-golang/v10/pkg/cloudwrapper"
	"github.com/akamai/terraform-provider-akamai/v7/pkg/common/ptr"
	"github.com/akamai/terraform-provider-akamai/v7/pkg/common/testutils"
	"github.com/hashicorp/terraform-plugin-testing/helper/resource"
)

var (
	minimalConfiguration = testDataForCWConfiguration{
		ID:                1,
		Comments:          "Test comments",
		ContractID:        "Test contract",
		Status:            "Test status",
		ConfigName:        "Test config name",
		LastUpdatedBy:     "Test user",
		LastUpdatedDate:   "Test date",
		RetainIdleObjects: false,
	}

	configuration = testDataForCWConfiguration{
		ID:                      1,
		CapacityAlertsThreshold: ptr.To(1),
		Comments:                "Test comments",
		ContractID:              "Test contract",
		Locations: []cloudwrapper.ConfigLocationResp{
			{
				Comments:      "Test comments 1",
				TrafficTypeID: 11,
				Capacity: cloudwrapper.Capacity{
					Value: 111,
					Unit:  "GB",
				},
				MapName: "Test MapName 1",
			},
			{
				Comments:      "Test comments 2",
				TrafficTypeID: 22,
				Capacity: cloudwrapper.Capacity{
					Value: 222,
					Unit:  "TB",
				},
				MapName: "Test MapName 2",
			},
		},
		MultiCDNSettings: &cloudwrapper.MultiCDNSettings{
			BOCC: &cloudwrapper.BOCC{
				ConditionalSamplingFrequency: "ZERO",
				Enabled:                      true,
				ForwardType:                  "ORIGIN_ONLY",
				RequestType:                  "EDGE_ONLY",
				SamplingFrequency:            "ZERO",
			},
			CDNs: []cloudwrapper.CDN{
				{
					CDNAuthKeys: []cloudwrapper.CDNAuthKey{
						{
							AuthKeyName: "Test name 1",
							ExpiryDate:  "Test date 1",
							HeaderName:  "Test header name 1",
							Secret:      "Test secret 1",
						},
						{
							AuthKeyName: "Test name 2",
							ExpiryDate:  "Test date 2",
							HeaderName:  "Test header name 2",
							Secret:      "Test secret 2",
						},
					},
					CDNCode:    "Test code",
					Enabled:    true,
					HTTPSOnly:  true,
					IPACLCIDRs: []string{"1.1.1.1", "2.2.2.2"},
				},
				{
					CDNAuthKeys: []cloudwrapper.CDNAuthKey{
						{
							AuthKeyName: "Test name 1",
							ExpiryDate:  "Test date 1",
							HeaderName:  "Test header name 1",
							Secret:      "Test secret 1",
						},
					},
					CDNCode:    "Test code",
					IPACLCIDRs: []string{"1.1.1.1", "2.2.2.2"},
				},
			},
			DataStreams: &cloudwrapper.DataStreams{
				DataStreamIDs: []int64{1, 2},
				Enabled:       true,
				SamplingRate:  ptr.To(10),
			},
			EnableSoftAlerts: true,
			Origins: []cloudwrapper.Origin{
				{
					Hostname:   "Test hostname 1",
					OriginID:   "Test originID 1",
					PropertyID: 11,
				},
				{
					Hostname:   "Test hostname 2",
					OriginID:   "Test originID 2",
					PropertyID: 22,
				},
			},
		},
		Status:             "Test status",
		ConfigName:         "Test config name",
		LastUpdatedBy:      "Test user",
		LastUpdatedDate:    "Test date",
		LastActivatedBy:    ptr.To("Test user 2"),
		LastActivatedDate:  ptr.To("Test date 2"),
		NotificationEmails: []string{"1@a.com", "2@a.com"},
		PropertyIDs:        []string{"11", "22"},
		RetainIdleObjects:  true,
	}
)

func TestConfigurationDataSource(t *testing.T) {
	t.Parallel()

	tests := map[string]struct {
		configPath string
		init       func(*cloudwrapper.Mock, testDataForCWConfiguration)
		mockData   testDataForCWConfiguration
		error      *regexp.Regexp
	}{
		"happy path - minimal data returned": {
			configPath: "testdata/TestDataConfiguration/default.tf",
			init: func(m *cloudwrapper.Mock, testData testDataForCWConfiguration) {
				expectGetConfiguration(m, testData, 3)
			},
			mockData: minimalConfiguration,
		},
		"happy path - all fields": {
			configPath: "testdata/TestDataConfiguration/default.tf",
			init: func(m *cloudwrapper.Mock, testData testDataForCWConfiguration) {
				expectGetConfiguration(m, testData, 3)
			},
			mockData: configuration,
		},
		"error getting configuration": {
			configPath: "testdata/TestDataConfiguration/default.tf",
			init: func(m *cloudwrapper.Mock, testData testDataForCWConfiguration) {
				expectGetConfigurationWithError(m, testData, 1)
			},
			mockData: testDataForCWConfiguration{
				ID: 1,
			},
			error: regexp.MustCompile("get configuration failed"),
		},
		"no required argument - configID": {
			configPath: "testdata/TestDataConfiguration/no_config_id.tf",
			error:      regexp.MustCompile(`The argument "id" is required, but no definition was found.`),
		},
	}
	for name, test := range tests {
		name, test := name, test
		t.Run(name, func(t *testing.T) {
			t.Parallel()

			client := &cloudwrapper.Mock{}
			if test.init != nil {
				test.init(client, test.mockData)
			}

			resource.UnitTest(t, resource.TestCase{
				ProtoV6ProviderFactories: newProviderFactory(withMockClient(client)),
				IsUnitTest:               true,
				Steps: []resource.TestStep{
					{
						Config:      testutils.LoadFixtureString(t, test.configPath),
						Check:       checkCloudWrapperConfigurationAttrs(test.mockData),
						ExpectError: test.error,
					},
				},
			})

			client.AssertExpectations(t)
		})
	}
}

type testDataForCWConfiguration struct {
	ID                      int64
	CapacityAlertsThreshold *int
	Comments                string
	ContractID              string
	Locations               []cloudwrapper.ConfigLocationResp
	MultiCDNSettings        *cloudwrapper.MultiCDNSettings
	Status                  string
	ConfigName              string
	LastUpdatedBy           string
	LastUpdatedDate         string
	LastActivatedBy         *string
	LastActivatedDate       *string
	NotificationEmails      []string
	PropertyIDs             []string
	RetainIdleObjects       bool
}

func checkCloudWrapperConfigurationAttrs(data testDataForCWConfiguration) resource.TestCheckFunc {
	var checkFuncs []resource.TestCheckFunc

	checkFuncs = append(checkFuncs, checkConfiguration(data, "data.akamai_cloudwrapper_configuration.test", ""))

	return resource.ComposeAggregateTestCheckFunc(checkFuncs...)
}

func checkConfiguration(data testDataForCWConfiguration, dsName, keyPrefix string) resource.TestCheckFunc {
	var checkFuncs []resource.TestCheckFunc

	if data.CapacityAlertsThreshold != nil {
		checkFuncs = append(checkFuncs, resource.TestCheckResourceAttr(dsName, keyPrefix+"capacity_alerts_threshold", strconv.Itoa(*data.CapacityAlertsThreshold)))
	} else {
		checkFuncs = append(checkFuncs, resource.TestCheckNoResourceAttr(dsName, keyPrefix+"capacity_alerts_threshold"))
	}
	checkFuncs = append(checkFuncs, resource.TestCheckResourceAttr(dsName, keyPrefix+"comments", data.Comments),
		resource.TestCheckResourceAttr(dsName, keyPrefix+"id", strconv.Itoa(int(data.ID))),
		resource.TestCheckResourceAttr(dsName, keyPrefix+"contract_id", data.ContractID),
		resource.TestCheckResourceAttr(dsName, keyPrefix+"config_name", data.ConfigName),
		resource.TestCheckResourceAttr(dsName, keyPrefix+"last_updated_by", data.LastUpdatedBy),
		resource.TestCheckResourceAttr(dsName, keyPrefix+"last_updated_date", data.LastUpdatedDate),
		resource.TestCheckResourceAttr(dsName, keyPrefix+"retain_idle_objects", strconv.FormatBool(data.RetainIdleObjects)),
		resource.TestCheckResourceAttr(dsName, keyPrefix+"status", data.Status),
	)
	if data.LastActivatedBy != nil {
		checkFuncs = append(checkFuncs, resource.TestCheckResourceAttr(dsName, keyPrefix+"last_activated_by", *data.LastActivatedBy))
	} else {
		checkFuncs = append(checkFuncs, resource.TestCheckNoResourceAttr(dsName, keyPrefix+"last_activated_by"))
	}
	if data.LastActivatedDate != nil {
		checkFuncs = append(checkFuncs, resource.TestCheckResourceAttr(dsName, keyPrefix+"last_activated_date", *data.LastActivatedDate))
	} else {
		checkFuncs = append(checkFuncs, resource.TestCheckNoResourceAttr(dsName, keyPrefix+"last_activated_date"))
	}
	checkFuncs = append(checkFuncs, resource.TestCheckResourceAttr(dsName, keyPrefix+"notification_emails.#", strconv.Itoa(len(data.NotificationEmails))))
	for i, email := range data.NotificationEmails {
		checkFuncs = append(checkFuncs, resource.TestCheckResourceAttr(dsName, fmt.Sprintf(keyPrefix+"notification_emails.%d", i), email))
	}

	checkFuncs = append(checkFuncs, resource.TestCheckResourceAttr(dsName, keyPrefix+"property_ids.#", strconv.Itoa(len(data.PropertyIDs))))
	for i, prpID := range data.PropertyIDs {
		checkFuncs = append(checkFuncs, resource.TestCheckResourceAttr(dsName, fmt.Sprintf(keyPrefix+"property_ids.%d", i), prpID))
	}
	checkFuncs = append(checkFuncs, resource.TestCheckResourceAttr(dsName, keyPrefix+"retain_idle_objects", strconv.FormatBool(data.RetainIdleObjects)))

	checkFuncs = append(checkFuncs, resource.TestCheckResourceAttr(dsName, keyPrefix+"locations.#", strconv.Itoa(len(data.Locations))))
	for i, loc := range data.Locations {
		checkFuncs = append(checkFuncs,
			resource.TestCheckResourceAttr(dsName, fmt.Sprintf(keyPrefix+"locations.%d.capacity.value", i), strconv.FormatInt(loc.Capacity.Value, 10)),
			resource.TestCheckResourceAttr(dsName, fmt.Sprintf(keyPrefix+"locations.%d.capacity.unit", i), string(loc.Capacity.Unit)),
			resource.TestCheckResourceAttr(dsName, fmt.Sprintf(keyPrefix+"locations.%d.comments", i), loc.Comments),
			resource.TestCheckResourceAttr(dsName, fmt.Sprintf(keyPrefix+"locations.%d.map_name", i), loc.MapName),
			resource.TestCheckResourceAttr(dsName, fmt.Sprintf(keyPrefix+"locations.%d.traffic_type_id", i), strconv.Itoa(loc.TrafficTypeID)),
		)
	}

	if data.MultiCDNSettings != nil {
		checkFuncs = append(checkFuncs, resource.TestCheckResourceAttr(dsName, keyPrefix+"multi_cdn_settings.enable_soft_alerts", strconv.FormatBool(data.MultiCDNSettings.EnableSoftAlerts)))
		if data.MultiCDNSettings.DataStreams != nil {

			checkFuncs = append(checkFuncs, resource.TestCheckResourceAttr(dsName, keyPrefix+"multi_cdn_settings.data_streams.enabled", strconv.FormatBool(data.MultiCDNSettings.DataStreams.Enabled)))
			if data.MultiCDNSettings.DataStreams.SamplingRate != nil {
				checkFuncs = append(checkFuncs, resource.TestCheckResourceAttr(dsName, keyPrefix+"multi_cdn_settings.data_streams.sampling_rate", strconv.Itoa(*data.MultiCDNSettings.DataStreams.SamplingRate)))
			} else {
				checkFuncs = append(checkFuncs, resource.TestCheckNoResourceAttr(dsName, keyPrefix+"multi_cdn_settings.data_streams.sampling_rate"))
			}
			checkFuncs = append(checkFuncs, resource.TestCheckResourceAttr(dsName, keyPrefix+"multi_cdn_settings.data_streams.data_stream_ids.#", strconv.Itoa(len(data.MultiCDNSettings.DataStreams.DataStreamIDs))))
			for i, id := range data.MultiCDNSettings.DataStreams.DataStreamIDs {
				checkFuncs = append(checkFuncs, resource.TestCheckResourceAttr(dsName, fmt.Sprintf(keyPrefix+"multi_cdn_settings.data_streams.data_stream_ids.%d", i), strconv.Itoa(int(id))))
			}
		}

		if data.MultiCDNSettings.BOCC != nil {
			checkFuncs = append(checkFuncs,
				resource.TestCheckResourceAttr(dsName, keyPrefix+"multi_cdn_settings.bocc.enabled", strconv.FormatBool(data.MultiCDNSettings.BOCC.Enabled)),
				resource.TestCheckResourceAttr(dsName, keyPrefix+"multi_cdn_settings.bocc.request_type", string(data.MultiCDNSettings.BOCC.RequestType)),
				resource.TestCheckResourceAttr(dsName, keyPrefix+"multi_cdn_settings.bocc.forward_type", string(data.MultiCDNSettings.BOCC.ForwardType)),
				resource.TestCheckResourceAttr(dsName, keyPrefix+"multi_cdn_settings.bocc.conditional_sampling_frequency", string(data.MultiCDNSettings.BOCC.ConditionalSamplingFrequency)),
				resource.TestCheckResourceAttr(dsName, keyPrefix+"multi_cdn_settings.bocc.sampling_frequency", string(data.MultiCDNSettings.BOCC.SamplingFrequency)),
			)
		}

		checkFuncs = append(checkFuncs, resource.TestCheckResourceAttr(dsName, keyPrefix+"multi_cdn_settings.origins.#", strconv.Itoa(len(data.MultiCDNSettings.Origins))))
		for i, origin := range data.MultiCDNSettings.Origins {
			checkFuncs = append(checkFuncs,
				resource.TestCheckResourceAttr(dsName, fmt.Sprintf(keyPrefix+"multi_cdn_settings.origins.%d.origin_id", i), origin.OriginID),
				resource.TestCheckResourceAttr(dsName, fmt.Sprintf(keyPrefix+"multi_cdn_settings.origins.%d.hostname", i), origin.Hostname),
				resource.TestCheckResourceAttr(dsName, fmt.Sprintf(keyPrefix+"multi_cdn_settings.origins.%d.property_id", i), strconv.Itoa(origin.PropertyID)),
			)
		}

		checkFuncs = append(checkFuncs, resource.TestCheckResourceAttr(dsName, keyPrefix+"multi_cdn_settings.cdns.#", strconv.Itoa(len(data.MultiCDNSettings.CDNs))))
		for i, cdn := range data.MultiCDNSettings.CDNs {
			checkFuncs = append(checkFuncs, resource.TestCheckResourceAttr(dsName, fmt.Sprintf(keyPrefix+"multi_cdn_settings.cdns.%d.cdn_code", i), cdn.CDNCode))
			checkFuncs = append(checkFuncs, resource.TestCheckResourceAttr(dsName, fmt.Sprintf(keyPrefix+"multi_cdn_settings.cdns.%d.enabled", i), strconv.FormatBool(cdn.Enabled)))
			checkFuncs = append(checkFuncs, resource.TestCheckResourceAttr(dsName, fmt.Sprintf(keyPrefix+"multi_cdn_settings.cdns.%d.https_only", i), strconv.FormatBool(cdn.HTTPSOnly)))
			checkFuncs = append(checkFuncs, resource.TestCheckResourceAttr(dsName, fmt.Sprintf(keyPrefix+"multi_cdn_settings.cdns.%d.ip_acl_cidrs.#", i), strconv.Itoa(len(cdn.IPACLCIDRs))))
			for j, ip := range cdn.IPACLCIDRs {
				checkFuncs = append(checkFuncs, resource.TestCheckResourceAttr(dsName, fmt.Sprintf(keyPrefix+"multi_cdn_settings.cdns.%d.ip_acl_cidrs.%d", i, j), ip))
			}
			checkFuncs = append(checkFuncs, resource.TestCheckResourceAttr(dsName, fmt.Sprintf(keyPrefix+"multi_cdn_settings.cdns.%d.cdn_auth_keys.#", i), strconv.Itoa(len(cdn.CDNAuthKeys))))
			for j, authKey := range cdn.CDNAuthKeys {
				checkFuncs = append(checkFuncs,
					resource.TestCheckResourceAttr(dsName, fmt.Sprintf(keyPrefix+"multi_cdn_settings.cdns.%d.cdn_auth_keys.%d.auth_key_name", i, j), authKey.AuthKeyName),
					resource.TestCheckResourceAttr(dsName, fmt.Sprintf(keyPrefix+"multi_cdn_settings.cdns.%d.cdn_auth_keys.%d.header_name", i, j), authKey.HeaderName),
					resource.TestCheckResourceAttr(dsName, fmt.Sprintf(keyPrefix+"multi_cdn_settings.cdns.%d.cdn_auth_keys.%d.secret", i, j), authKey.Secret),
					resource.TestCheckResourceAttr(dsName, fmt.Sprintf(keyPrefix+"multi_cdn_settings.cdns.%d.cdn_auth_keys.%d.expiry_date", i, j), authKey.ExpiryDate),
				)
			}

		}
	}
	return resource.ComposeAggregateTestCheckFunc(checkFuncs...)
}

func expectGetConfiguration(client *cloudwrapper.Mock, data testDataForCWConfiguration, timesToRun int) {
	getConfigurationReq := cloudwrapper.GetConfigurationRequest{
		ConfigID: data.ID,
	}
	getConfigurationRes := getConfiguration(data)
	client.On("GetConfiguration", testutils.MockContext, getConfigurationReq).Return(&getConfigurationRes, nil).Times(timesToRun)
}

func getConfiguration(data testDataForCWConfiguration) cloudwrapper.Configuration {
	return cloudwrapper.Configuration{
		CapacityAlertsThreshold: data.CapacityAlertsThreshold,
		Comments:                data.Comments,
		ContractID:              data.ContractID,
		ConfigID:                data.ID,
		Locations:               data.Locations,
		MultiCDNSettings:        data.MultiCDNSettings,
		Status:                  cloudwrapper.StatusType(data.Status),
		ConfigName:              data.ConfigName,
		LastUpdatedBy:           data.LastUpdatedBy,
		LastUpdatedDate:         data.LastUpdatedDate,
		LastActivatedBy:         data.LastActivatedBy,
		LastActivatedDate:       data.LastActivatedDate,
		NotificationEmails:      data.NotificationEmails,
		PropertyIDs:             data.PropertyIDs,
		RetainIdleObjects:       data.RetainIdleObjects,
	}
}

func expectGetConfigurationWithError(client *cloudwrapper.Mock, data testDataForCWConfiguration, timesToRun int) {
	getConfigurationReq := cloudwrapper.GetConfigurationRequest{
		ConfigID: data.ID,
	}
	client.On("GetConfiguration", testutils.MockContext, getConfigurationReq).Return(nil, fmt.Errorf("get configuration failed")).Times(timesToRun)
}
