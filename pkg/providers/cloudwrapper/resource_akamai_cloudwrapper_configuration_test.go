package cloudwrapper

import (
	"regexp"
	"testing"

	"github.com/akamai/AkamaiOPEN-edgegrid-golang/v10/pkg/cloudwrapper"
	"github.com/akamai/terraform-provider-akamai/v8/pkg/common/testutils"
	"github.com/hashicorp/terraform-plugin-testing/helper/resource"
	"github.com/jinzhu/copier"
	"github.com/stretchr/testify/assert"
)

func TestConfigurationResource(t *testing.T) {
	t.Parallel()
	t.Run("create basic", func(t *testing.T) {
		t.Parallel()
		client := &cloudwrapper.Mock{}

		configuration := cloudwrapper.CreateConfigurationRequest{
			Body: cloudwrapper.CreateConfigurationRequestBody{
				Comments:   "test",
				ContractID: "ctr_123",
				Locations: []cloudwrapper.ConfigLocationReq{
					{
						Comments:      "test",
						TrafficTypeID: 1,
						Capacity: cloudwrapper.Capacity{
							Value: 1,
							Unit:  cloudwrapper.UnitGB,
						},
					},
				},
				ConfigName:         "testname",
				NotificationEmails: []string{"test@akamai.com"},
				PropertyIDs:        []string{"200200200"},
				RetainIdleObjects:  false,
			},
		}

		expecter := newExpecter(t, client)

		expecter.ExpectCreate(configuration)

		expecter.ExpectRefresh()
		expecter.ExpectDelete()

		resource.Test(t, resource.TestCase{
			IsUnitTest:               true,
			ProtoV6ProviderFactories: newProviderFactory(withMockClient(client)),
			Steps: []resource.TestStep{
				{
					Config: testutils.LoadFixtureString(t, "testdata/TestResConfiguration/create.tf"),
					Check: resource.ComposeAggregateTestCheckFunc(
						resource.TestCheckResourceAttr("akamai_cloudwrapper_configuration.test", "id", "123"),
						resource.TestCheckResourceAttr("akamai_cloudwrapper_configuration.test", "config_name", "testname"),
						resource.TestCheckResourceAttr("akamai_cloudwrapper_configuration.test", "contract_id", "ctr_123"),
						resource.TestCheckResourceAttr("akamai_cloudwrapper_configuration.test", "property_ids.0", "200200200"),
						resource.TestCheckResourceAttr("akamai_cloudwrapper_configuration.test", "notification_emails.0", "test@akamai.com"),
						resource.TestCheckResourceAttr("akamai_cloudwrapper_configuration.test", "comments", "test"),
						resource.TestCheckResourceAttr("akamai_cloudwrapper_configuration.test", "retain_idle_objects", "false"),
						resource.TestCheckResourceAttr("akamai_cloudwrapper_configuration.test", "location.#", "1"),
						resource.TestCheckResourceAttr("akamai_cloudwrapper_configuration.test", "location.0.comments", "test"),
						resource.TestCheckResourceAttr("akamai_cloudwrapper_configuration.test", "location.0.traffic_type_id", "1"),
						resource.TestCheckResourceAttr("akamai_cloudwrapper_configuration.test", "location.0.capacity.value", "1"),
						resource.TestCheckResourceAttr("akamai_cloudwrapper_configuration.test", "location.0.capacity.unit", "GB"),
						resource.TestCheckNoResourceAttr("akamai_cloudwrapper_configuration.test", "capacity_alerts_threshold"),
					),
				},
			},
		})

		client.AssertExpectations(t)
	})
	t.Run("email will be computed when not provided", func(t *testing.T) {
		t.Parallel()
		client := &cloudwrapper.Mock{}

		configuration := cloudwrapper.CreateConfigurationRequest{
			Body: cloudwrapper.CreateConfigurationRequestBody{
				Comments:   "test",
				ContractID: "ctr_123",
				Locations: []cloudwrapper.ConfigLocationReq{
					{
						Comments:      "test",
						TrafficTypeID: 1,
						Capacity: cloudwrapper.Capacity{
							Value: 1,
							Unit:  cloudwrapper.UnitGB,
						},
					},
				},
				ConfigName:        "testname",
				PropertyIDs:       []string{"200200200"},
				RetainIdleObjects: false,
			},
		}

		expecter := newExpecter(t, client)

		expecter.ExpectCreate(configuration)

		expecter.ExpectRefresh()
		expecter.ExpectDelete()

		resource.Test(t, resource.TestCase{
			IsUnitTest:               true,
			ProtoV6ProviderFactories: newProviderFactory(withMockClient(client)),
			Steps: []resource.TestStep{
				{
					Config: testutils.LoadFixtureString(t, "testdata/TestResConfiguration/computed_email.tf"),
					Check: resource.ComposeAggregateTestCheckFunc(
						resource.TestCheckResourceAttr("akamai_cloudwrapper_configuration.test", "id", "123"),
						resource.TestCheckResourceAttr("akamai_cloudwrapper_configuration.test", "config_name", "testname"),
						resource.TestCheckResourceAttr("akamai_cloudwrapper_configuration.test", "notification_emails.#", "1"),
						resource.TestCheckResourceAttr("akamai_cloudwrapper_configuration.test", "notification_emails.0", "generated@akamai.com"),
					),
				},
			},
		})

		client.AssertExpectations(t)
	})
	t.Run("retain_idle_objects has default when not provided", func(t *testing.T) {
		t.Parallel()
		client := &cloudwrapper.Mock{}

		configuration := cloudwrapper.CreateConfigurationRequest{
			Body: cloudwrapper.CreateConfigurationRequestBody{
				Comments:   "test",
				ContractID: "ctr_123",
				Locations: []cloudwrapper.ConfigLocationReq{
					{
						Comments:      "test",
						TrafficTypeID: 1,
						Capacity: cloudwrapper.Capacity{
							Value: 1,
							Unit:  cloudwrapper.UnitGB,
						},
					},
				},
				ConfigName:        "testname",
				PropertyIDs:       []string{"200200200"},
				RetainIdleObjects: false,
			},
		}

		expecter := newExpecter(t, client)

		expecter.ExpectCreate(configuration)

		expecter.ExpectRefresh()
		expecter.ExpectDelete()

		resource.Test(t, resource.TestCase{
			IsUnitTest:               true,
			ProtoV6ProviderFactories: newProviderFactory(withMockClient(client)),
			Steps: []resource.TestStep{
				{
					Config: testutils.LoadFixtureString(t, "testdata/TestResConfiguration/computed_email.tf"),
					Check: resource.ComposeAggregateTestCheckFunc(
						resource.TestCheckResourceAttr("akamai_cloudwrapper_configuration.test", "id", "123"),
						resource.TestCheckResourceAttr("akamai_cloudwrapper_configuration.test", "config_name", "testname"),
						resource.TestCheckResourceAttr("akamai_cloudwrapper_configuration.test", "notification_emails.#", "1"),
						resource.TestCheckResourceAttr("akamai_cloudwrapper_configuration.test", "notification_emails.0", "generated@akamai.com"),
					),
				},
			},
		})

		client.AssertExpectations(t)
	})
	t.Run("force new on config name change", func(t *testing.T) {
		t.Parallel()
		client := &cloudwrapper.Mock{}

		configuration := cloudwrapper.CreateConfigurationRequest{
			Body: cloudwrapper.CreateConfigurationRequestBody{
				Comments:   "test",
				ContractID: "ctr_123",
				Locations: []cloudwrapper.ConfigLocationReq{
					{
						Comments:      "test",
						TrafficTypeID: 1,
						Capacity: cloudwrapper.Capacity{
							Value: 1,
							Unit:  cloudwrapper.UnitGB,
						},
					},
				},
				NotificationEmails: []string{"test@akamai.com"},
				ConfigName:         "testname",
				PropertyIDs:        []string{"200200200"},
				RetainIdleObjects:  false,
			},
		}

		expecter := newExpecter(t, client)

		expecter.ExpectCreate(configuration)
		expecter.ExpectRefresh()
		expecter.ExpectRefresh()
		expecter.ExpectDelete()

		configUpdate := cloudwrapper.CreateConfigurationRequest{
			Body: cloudwrapper.CreateConfigurationRequestBody{
				Comments:   "test",
				ContractID: "ctr_123",
				Locations: []cloudwrapper.ConfigLocationReq{
					{
						Comments:      "test",
						TrafficTypeID: 1,
						Capacity: cloudwrapper.Capacity{
							Value: 1,
							Unit:  cloudwrapper.UnitGB,
						},
					},
				},
				NotificationEmails: []string{"test@akamai.com"},
				ConfigName:         "newname",
				PropertyIDs:        []string{"200200200"},
				RetainIdleObjects:  false,
			},
		}

		expecter = newExpecter(t, client)

		expecter.ExpectCreate(configUpdate)

		expecter.ExpectRefresh()
		expecter.ExpectDelete()

		resource.Test(t, resource.TestCase{
			IsUnitTest:               true,
			ProtoV6ProviderFactories: newProviderFactory(withMockClient(client)),
			Steps: []resource.TestStep{
				{
					Config: testutils.LoadFixtureString(t, "testdata/TestResConfiguration/create.tf"),
					Check: resource.ComposeAggregateTestCheckFunc(
						resource.TestCheckResourceAttr("akamai_cloudwrapper_configuration.test", "id", "123"),
						resource.TestCheckResourceAttr("akamai_cloudwrapper_configuration.test", "config_name", "testname"),
					),
				},
				{
					Config: testutils.LoadFixtureString(t, "testdata/TestResConfiguration/update_config_name.tf"),
					Check: resource.ComposeAggregateTestCheckFunc(
						resource.TestCheckResourceAttr("akamai_cloudwrapper_configuration.test", "id", "123"),
						resource.TestCheckResourceAttr("akamai_cloudwrapper_configuration.test", "config_name", "newname"),
					),
				},
			},
		})

		client.AssertExpectations(t)
	})
	t.Run("force new on contract_id change", func(t *testing.T) {
		t.Parallel()
		client := &cloudwrapper.Mock{}

		configuration := cloudwrapper.CreateConfigurationRequest{
			Body: cloudwrapper.CreateConfigurationRequestBody{
				Comments:   "test",
				ContractID: "ctr_123",
				Locations: []cloudwrapper.ConfigLocationReq{
					{
						Comments:      "test",
						TrafficTypeID: 1,
						Capacity: cloudwrapper.Capacity{
							Value: 1,
							Unit:  cloudwrapper.UnitGB,
						},
					},
				},
				NotificationEmails: []string{"test@akamai.com"},
				ConfigName:         "testname",
				PropertyIDs:        []string{"200200200"},
				RetainIdleObjects:  false,
			},
		}

		expecter := newExpecter(t, client)

		expecter.ExpectCreate(configuration)

		expecter.ExpectRefresh()

		expecter.ExpectRefresh()
		expecter.ExpectDelete()

		configUpdate := cloudwrapper.CreateConfigurationRequest{
			Body: cloudwrapper.CreateConfigurationRequestBody{
				Comments:   "test",
				ContractID: "ctr_234",
				Locations: []cloudwrapper.ConfigLocationReq{
					{
						Comments:      "test",
						TrafficTypeID: 1,
						Capacity: cloudwrapper.Capacity{
							Value: 1,
							Unit:  cloudwrapper.UnitGB,
						},
					},
				},
				NotificationEmails: []string{"test@akamai.com"},
				ConfigName:         "testname",
				PropertyIDs:        []string{"200200200"},
				RetainIdleObjects:  false,
			},
		}

		expecter = newExpecter(t, client)

		expecter.ExpectCreate(configUpdate)

		expecter.ExpectRefresh()
		expecter.ExpectDelete()

		resource.Test(t, resource.TestCase{
			IsUnitTest:               true,
			ProtoV6ProviderFactories: newProviderFactory(withMockClient(client)),
			Steps: []resource.TestStep{
				{
					Config: testutils.LoadFixtureString(t, "testdata/TestResConfiguration/create.tf"),
					Check: resource.ComposeAggregateTestCheckFunc(
						resource.TestCheckResourceAttr("akamai_cloudwrapper_configuration.test", "id", "123"),
						resource.TestCheckResourceAttr("akamai_cloudwrapper_configuration.test", "config_name", "testname"),
						resource.TestCheckResourceAttr("akamai_cloudwrapper_configuration.test", "contract_id", "ctr_123"),
					),
				},
				{
					Config: testutils.LoadFixtureString(t, "testdata/TestResConfiguration/update_contract_id.tf"),
					Check: resource.ComposeAggregateTestCheckFunc(
						resource.TestCheckResourceAttr("akamai_cloudwrapper_configuration.test", "id", "123"),
						resource.TestCheckResourceAttr("akamai_cloudwrapper_configuration.test", "config_name", "testname"),
						resource.TestCheckResourceAttr("akamai_cloudwrapper_configuration.test", "contract_id", "ctr_234"),
					),
				},
			},
		})

		client.AssertExpectations(t)
	})
	t.Run("force new on contract_id change", func(t *testing.T) {
		t.Parallel()
		client := &cloudwrapper.Mock{}

		configuration := cloudwrapper.CreateConfigurationRequest{
			Body: cloudwrapper.CreateConfigurationRequestBody{
				Comments:   "test",
				ContractID: "ctr_123",
				Locations: []cloudwrapper.ConfigLocationReq{
					{
						Comments:      "test",
						TrafficTypeID: 1,
						Capacity: cloudwrapper.Capacity{
							Value: 1,
							Unit:  cloudwrapper.UnitGB,
						},
					},
				},
				NotificationEmails: []string{"test@akamai.com"},
				ConfigName:         "testname",
				PropertyIDs:        []string{"200200200"},
				RetainIdleObjects:  false,
			},
		}

		expecter := newExpecter(t, client)

		expecter.ExpectCreate(configuration)

		expecter.ExpectRefresh()

		expecter.ExpectRefresh()
		expecter.ExpectDelete()

		configUpdate := cloudwrapper.CreateConfigurationRequest{
			Body: cloudwrapper.CreateConfigurationRequestBody{
				Comments:   "test",
				ContractID: "ctr_234",
				Locations: []cloudwrapper.ConfigLocationReq{
					{
						Comments:      "test",
						TrafficTypeID: 1,
						Capacity: cloudwrapper.Capacity{
							Value: 1,
							Unit:  cloudwrapper.UnitGB,
						},
					},
				},
				NotificationEmails: []string{"test@akamai.com"},
				ConfigName:         "testname",
				PropertyIDs:        []string{"200200200"},
				RetainIdleObjects:  false,
			},
		}

		expecter = newExpecter(t, client)

		expecter.ExpectCreate(configUpdate)

		expecter.ExpectRefresh()
		expecter.ExpectDelete()

		resource.Test(t, resource.TestCase{
			IsUnitTest:               true,
			ProtoV6ProviderFactories: newProviderFactory(withMockClient(client)),
			Steps: []resource.TestStep{
				{
					Config: testutils.LoadFixtureString(t, "testdata/TestResConfiguration/create.tf"),
					Check: resource.ComposeAggregateTestCheckFunc(
						resource.TestCheckResourceAttr("akamai_cloudwrapper_configuration.test", "id", "123"),
						resource.TestCheckResourceAttr("akamai_cloudwrapper_configuration.test", "config_name", "testname"),
						resource.TestCheckResourceAttr("akamai_cloudwrapper_configuration.test", "contract_id", "ctr_123"),
					),
				},
				{
					Config: testutils.LoadFixtureString(t, "testdata/TestResConfiguration/update_contract_id.tf"),
					Check: resource.ComposeAggregateTestCheckFunc(
						resource.TestCheckResourceAttr("akamai_cloudwrapper_configuration.test", "id", "123"),
						resource.TestCheckResourceAttr("akamai_cloudwrapper_configuration.test", "config_name", "testname"),
						resource.TestCheckResourceAttr("akamai_cloudwrapper_configuration.test", "contract_id", "ctr_234"),
					),
				},
			},
		})

		client.AssertExpectations(t)
	})
	t.Run("import", func(t *testing.T) {
		t.Parallel()
		client := &cloudwrapper.Mock{}

		configuration := cloudwrapper.CreateConfigurationRequest{
			Body: cloudwrapper.CreateConfigurationRequestBody{
				Comments:   "test",
				ContractID: "ctr_123",
				Locations: []cloudwrapper.ConfigLocationReq{
					{
						Comments:      "test",
						TrafficTypeID: 1,
						Capacity: cloudwrapper.Capacity{
							Value: 1,
							Unit:  cloudwrapper.UnitGB,
						},
					},
				},
				NotificationEmails: []string{"test@akamai.com"},
				ConfigName:         "testname",
				PropertyIDs:        []string{"200200200"},
				RetainIdleObjects:  false,
			},
		}

		expecter := newExpecter(t, client)

		expecter.ExpectCreate(configuration)

		expecter.ExpectRefresh()

		expecter.ExpectRefresh()

		expecter.ExpectRefresh()
		expecter.ExpectDelete()

		resource.Test(t, resource.TestCase{
			IsUnitTest:               true,
			ProtoV6ProviderFactories: newProviderFactory(withMockClient(client)),
			Steps: []resource.TestStep{
				{
					Config: testutils.LoadFixtureString(t, "testdata/TestResConfiguration/create.tf"),
					Check: resource.ComposeAggregateTestCheckFunc(
						resource.TestCheckResourceAttr("akamai_cloudwrapper_configuration.test", "id", "123"),
						resource.TestCheckResourceAttr("akamai_cloudwrapper_configuration.test", "config_name", "testname"),
						resource.TestCheckResourceAttr("akamai_cloudwrapper_configuration.test", "contract_id", "ctr_123"),
					),
				},
				{
					ImportState:   true,
					ImportStateId: "123",
					ResourceName:  "akamai_cloudwrapper_configuration.test",
				},
			},
		})

		client.AssertExpectations(t)
	})
	t.Run("basic update", func(t *testing.T) {
		t.Parallel()
		client := &cloudwrapper.Mock{}

		configuration := cloudwrapper.CreateConfigurationRequest{
			Body: cloudwrapper.CreateConfigurationRequestBody{
				Comments:   "test",
				ContractID: "ctr_123",
				Locations: []cloudwrapper.ConfigLocationReq{
					{
						Comments:      "test",
						TrafficTypeID: 1,
						Capacity: cloudwrapper.Capacity{
							Value: 1,
							Unit:  cloudwrapper.UnitGB,
						},
					},
				},
				NotificationEmails: []string{"test@akamai.com"},
				ConfigName:         "testname",
				PropertyIDs:        []string{"200200200"},
				RetainIdleObjects:  false,
			},
		}

		expecter := newExpecter(t, client)

		expecter.ExpectCreate(configuration)

		expecter.ExpectRefresh()

		capacityAlertsThreshold := 50

		configUpdate := cloudwrapper.UpdateConfigurationRequest{
			ConfigID: expecter.config.ConfigID,
			Body: cloudwrapper.UpdateConfigurationRequestBody{
				CapacityAlertsThreshold: &capacityAlertsThreshold,
				Comments:                "test",
				Locations: []cloudwrapper.ConfigLocationReq{
					{
						Comments:      "test",
						TrafficTypeID: 1,
						Capacity: cloudwrapper.Capacity{
							Value: 1,
							Unit:  cloudwrapper.UnitGB,
						},
					},
				},
				NotificationEmails: []string{"test@akamai.com"},
				PropertyIDs:        []string{"200200200"},
				RetainIdleObjects:  true,
			},
		}

		expecter.ExpectRefresh()
		expecter.ExpectUpdate(configUpdate)

		expecter.ExpectRefresh()

		configUpdate2 := cloudwrapper.UpdateConfigurationRequest{
			ConfigID: expecter.config.ConfigID,
			Body: cloudwrapper.UpdateConfigurationRequestBody{
				Comments: "test",
				Locations: []cloudwrapper.ConfigLocationReq{
					{
						Comments:      "test",
						TrafficTypeID: 1,
						Capacity: cloudwrapper.Capacity{
							Value: 1,
							Unit:  cloudwrapper.UnitGB,
						},
					},
				},
				NotificationEmails: []string{"test@akamai.com"},
				PropertyIDs:        []string{"200200200"},
				RetainIdleObjects:  false,
			},
		}

		expecter.ExpectRefresh()
		expecter.ExpectUpdate(configUpdate2)

		expecter.ExpectRefresh()
		expecter.ExpectDelete()

		resource.Test(t, resource.TestCase{
			IsUnitTest:               true,
			ProtoV6ProviderFactories: newProviderFactory(withMockClient(client)),
			Steps: []resource.TestStep{
				{
					Config: testutils.LoadFixtureString(t, "testdata/TestResConfiguration/create.tf"),
					Check: resource.ComposeAggregateTestCheckFunc(
						resource.TestCheckResourceAttr("akamai_cloudwrapper_configuration.test", "id", "123"),
						resource.TestCheckResourceAttr("akamai_cloudwrapper_configuration.test", "config_name", "testname"),
						resource.TestCheckResourceAttr("akamai_cloudwrapper_configuration.test", "contract_id", "ctr_123"),
						resource.TestCheckNoResourceAttr("akamai_cloudwrapper_configuration.test", "capacity_alerts_threshold"),
					),
				},
				{
					Config: testutils.LoadFixtureString(t, "testdata/TestResConfiguration/update_alerts_threshold.tf"),
					Check: resource.ComposeAggregateTestCheckFunc(
						resource.TestCheckResourceAttr("akamai_cloudwrapper_configuration.test", "id", "123"),
						resource.TestCheckResourceAttr("akamai_cloudwrapper_configuration.test", "config_name", "testname"),
						resource.TestCheckResourceAttr("akamai_cloudwrapper_configuration.test", "contract_id", "ctr_123"),
						resource.TestCheckResourceAttr("akamai_cloudwrapper_configuration.test", "capacity_alerts_threshold", "50"),
					),
				},
				{
					Config: testutils.LoadFixtureString(t, "testdata/TestResConfiguration/create.tf"),
					Check: resource.ComposeAggregateTestCheckFunc(
						resource.TestCheckResourceAttr("akamai_cloudwrapper_configuration.test", "id", "123"),
						resource.TestCheckResourceAttr("akamai_cloudwrapper_configuration.test", "config_name", "testname"),
						resource.TestCheckResourceAttr("akamai_cloudwrapper_configuration.test", "contract_id", "ctr_123"),
						resource.TestCheckNoResourceAttr("akamai_cloudwrapper_configuration.test", "capacity_alerts_threshold"),
					),
				},
			},
		})
	})
	t.Run("drift - config got removed", func(t *testing.T) {
		t.Parallel()
		client := &cloudwrapper.Mock{}

		configuration := cloudwrapper.CreateConfigurationRequest{
			Body: cloudwrapper.CreateConfigurationRequestBody{
				Comments:   "test",
				ContractID: "ctr_123",
				Locations: []cloudwrapper.ConfigLocationReq{
					{
						Comments:      "test",
						TrafficTypeID: 1,
						Capacity: cloudwrapper.Capacity{
							Value: 1,
							Unit:  cloudwrapper.UnitGB,
						},
					},
				},
				NotificationEmails: []string{"test@akamai.com"},
				ConfigName:         "testname",
				PropertyIDs:        []string{"200200200"},
				RetainIdleObjects:  false,
			},
		}

		expecter := newExpecter(t, client)
		expecter.ExpectCreate(configuration)

		expecter.ExpectRefresh()

		expecter.ExpectDriftRefresh(nil, cloudwrapper.ErrConfigurationNotFound)
		expecter = newExpecter(t, client)
		expecter.ExpectCreate(configuration)

		expecter.ExpectRefresh()
		expecter.ExpectDelete()

		resource.Test(t, resource.TestCase{
			IsUnitTest:               true,
			ProtoV6ProviderFactories: newProviderFactory(withMockClient(client)),
			Steps: []resource.TestStep{
				{
					Config: testutils.LoadFixtureString(t, "testdata/TestResConfiguration/create.tf"),
					Check: resource.ComposeAggregateTestCheckFunc(
						resource.TestCheckResourceAttr("akamai_cloudwrapper_configuration.test", "id", "123"),
						resource.TestCheckResourceAttr("akamai_cloudwrapper_configuration.test", "config_name", "testname"),
						resource.TestCheckResourceAttr("akamai_cloudwrapper_configuration.test", "contract_id", "ctr_123"),
						resource.TestCheckNoResourceAttr("akamai_cloudwrapper_configuration.test", "capacity_alerts_threshold"),
					),
				},
				{
					Config: testutils.LoadFixtureString(t, "testdata/TestResConfiguration/create.tf"),
					Check: resource.ComposeAggregateTestCheckFunc(
						resource.TestCheckResourceAttr("akamai_cloudwrapper_configuration.test", "id", "123"),
						resource.TestCheckResourceAttr("akamai_cloudwrapper_configuration.test", "config_name", "testname"),
						resource.TestCheckResourceAttr("akamai_cloudwrapper_configuration.test", "contract_id", "ctr_123"),
						resource.TestCheckNoResourceAttr("akamai_cloudwrapper_configuration.test", "capacity_alerts_threshold"),
					),
				},
			},
		})

		client.AssertExpectations(t)
	})
	t.Run("contract_id remove prefix expect no diff", func(t *testing.T) {
		t.Parallel()
		client := &cloudwrapper.Mock{}

		configuration := cloudwrapper.CreateConfigurationRequest{
			Body: cloudwrapper.CreateConfigurationRequestBody{
				Comments:   "test",
				ContractID: "ctr_123",
				Locations: []cloudwrapper.ConfigLocationReq{
					{
						Comments:      "test",
						TrafficTypeID: 1,
						Capacity: cloudwrapper.Capacity{
							Value: 1,
							Unit:  cloudwrapper.UnitGB,
						},
					},
				},
				NotificationEmails: []string{"test@akamai.com"},
				ConfigName:         "testname",
				PropertyIDs:        []string{"200200200"},
				RetainIdleObjects:  false,
			},
		}

		expecter := newExpecter(t, client)
		expecter.ExpectCreate(configuration)

		expecter.ExpectRefresh()

		expecter.ExpectRefresh()

		expecter.ExpectRefresh()
		expecter.ExpectDelete()

		resource.Test(t, resource.TestCase{
			IsUnitTest:               true,
			ProtoV6ProviderFactories: newProviderFactory(withMockClient(client)),
			Steps: []resource.TestStep{
				{
					Config: testutils.LoadFixtureString(t, "testdata/TestResConfiguration/create.tf"),
					Check: resource.ComposeAggregateTestCheckFunc(
						resource.TestCheckResourceAttr("akamai_cloudwrapper_configuration.test", "id", "123"),
						resource.TestCheckResourceAttr("akamai_cloudwrapper_configuration.test", "config_name", "testname"),
						resource.TestCheckResourceAttr("akamai_cloudwrapper_configuration.test", "contract_id", "ctr_123"),
					),
				},
				{
					Config: testutils.LoadFixtureString(t, "testdata/TestResConfiguration/contract_id_no_prefix.tf"),
					Check: resource.ComposeAggregateTestCheckFunc(
						resource.TestCheckResourceAttr("akamai_cloudwrapper_configuration.test", "id", "123"),
						resource.TestCheckResourceAttr("akamai_cloudwrapper_configuration.test", "config_name", "testname"),
						resource.TestCheckResourceAttr("akamai_cloudwrapper_configuration.test", "contract_id", "ctr_123"),
					),
				},
			},
		})

		client.AssertExpectations(t)
	})
	t.Run("property_ids with prefix", func(t *testing.T) {
		t.Parallel()
		client := &cloudwrapper.Mock{}

		configuration := cloudwrapper.CreateConfigurationRequest{
			Body: cloudwrapper.CreateConfigurationRequestBody{
				Comments:   "test",
				ContractID: "ctr_123",
				Locations: []cloudwrapper.ConfigLocationReq{
					{
						Comments:      "test",
						TrafficTypeID: 1,
						Capacity: cloudwrapper.Capacity{
							Value: 1,
							Unit:  cloudwrapper.UnitGB,
						},
					},
				},
				NotificationEmails: []string{"test@akamai.com"},
				ConfigName:         "testname",
				PropertyIDs:        []string{"200200200"},
				RetainIdleObjects:  false,
			},
		}

		expecter := newExpecter(t, client)
		expecter.ExpectCreate(configuration)
		expecter.ExpectRefresh()

		expecter.ExpectDriftRefresh(&cloudwrapper.Configuration{
			Comments:   "test",
			ContractID: "ctr_123",
			ConfigID:   123,
			Locations: []cloudwrapper.ConfigLocationResp{
				{
					Comments:      "test",
					TrafficTypeID: 1,
					Capacity: cloudwrapper.Capacity{
						Value: 1,
						Unit:  cloudwrapper.UnitGB,
					},
				},
			},
			Status:             "SAVED",
			ConfigName:         "testname",
			NotificationEmails: []string{"test@akamai.com"},
			PropertyIDs:        []string{"300300300"},
			RetainIdleObjects:  false,
		}, nil)

		expecter.ExpectUpdate(cloudwrapper.UpdateConfigurationRequest{
			ConfigID: 123,
			Activate: false,
			Body: cloudwrapper.UpdateConfigurationRequestBody{
				Comments: "test",
				Locations: []cloudwrapper.ConfigLocationReq{
					{
						Comments:      "test",
						TrafficTypeID: 1,
						Capacity: cloudwrapper.Capacity{
							Value: 1,
							Unit:  cloudwrapper.UnitGB,
						},
					},
				},
				NotificationEmails: []string{"test@akamai.com"},
				PropertyIDs:        []string{"200200200"},
				RetainIdleObjects:  false,
			},
		})

		expecter.ExpectRefresh()
		expecter.ExpectDelete()

		resource.Test(t, resource.TestCase{
			IsUnitTest:               true,
			ProtoV6ProviderFactories: newProviderFactory(withMockClient(client)),
			Steps: []resource.TestStep{
				{
					Config: testutils.LoadFixtureString(t, "testdata/TestResConfiguration/property_ids_with_prefix.tf"),
					Check: resource.ComposeAggregateTestCheckFunc(
						resource.TestCheckResourceAttr("akamai_cloudwrapper_configuration.test", "id", "123"),
						resource.TestCheckResourceAttr("akamai_cloudwrapper_configuration.test", "config_name", "testname"),
						resource.TestCheckResourceAttr("akamai_cloudwrapper_configuration.test", "contract_id", "ctr_123"),
					),
				},
				{
					Config: testutils.LoadFixtureString(t, "testdata/TestResConfiguration/property_ids_with_prefix.tf"),
					Check: resource.ComposeAggregateTestCheckFunc(
						resource.TestCheckResourceAttr("akamai_cloudwrapper_configuration.test", "id", "123"),
						resource.TestCheckResourceAttr("akamai_cloudwrapper_configuration.test", "config_name", "testname"),
						resource.TestCheckResourceAttr("akamai_cloudwrapper_configuration.test", "contract_id", "ctr_123"),
					),
				},
			},
		})

		client.AssertExpectations(t)
	})
	t.Run("multicdn drift error", func(t *testing.T) {
		t.Parallel()
		client := &cloudwrapper.Mock{}

		configuration := cloudwrapper.CreateConfigurationRequest{
			Body: cloudwrapper.CreateConfigurationRequestBody{
				Comments:   "test",
				ContractID: "ctr_123",
				Locations: []cloudwrapper.ConfigLocationReq{
					{
						Comments:      "test",
						TrafficTypeID: 1,
						Capacity: cloudwrapper.Capacity{
							Value: 1,
							Unit:  cloudwrapper.UnitGB,
						},
					},
				},
				NotificationEmails: []string{"test@akamai.com"},
				ConfigName:         "testname",
				PropertyIDs:        []string{"200200200"},
				RetainIdleObjects:  false,
			},
		}

		expecter := newExpecter(t, client)
		expecter.ExpectCreate(configuration)
		expecter.ExpectRefresh()

		expecter.ExpectDriftRefresh(&cloudwrapper.Configuration{
			Comments:   "test",
			ContractID: "ctr_123",
			ConfigID:   123,
			Locations: []cloudwrapper.ConfigLocationResp{
				{
					Comments:      "test",
					TrafficTypeID: 1,
					Capacity: cloudwrapper.Capacity{
						Value: 1,
						Unit:  cloudwrapper.UnitGB,
					},
				},
			},
			MultiCDNSettings:   &cloudwrapper.MultiCDNSettings{}, // ADDED MULTICDN
			Status:             "SAVED",
			ConfigName:         "testname",
			NotificationEmails: []string{"test@akamai.com"},
			PropertyIDs:        []string{"300300300"},
			RetainIdleObjects:  false,
		}, nil)

		expecter.ExpectDelete()

		resource.Test(t, resource.TestCase{
			IsUnitTest:               true,
			ProtoV6ProviderFactories: newProviderFactory(withMockClient(client)),
			Steps: []resource.TestStep{
				{
					Config: testutils.LoadFixtureString(t, "testdata/TestResConfiguration/property_ids_with_prefix.tf"),
					Check: resource.ComposeAggregateTestCheckFunc(
						resource.TestCheckResourceAttr("akamai_cloudwrapper_configuration.test", "id", "123"),
						resource.TestCheckResourceAttr("akamai_cloudwrapper_configuration.test", "config_name", "testname"),
						resource.TestCheckResourceAttr("akamai_cloudwrapper_configuration.test", "contract_id", "ctr_123"),
					),
				},
				{
					Config: testutils.LoadFixtureString(t, "testdata/TestResConfiguration/property_ids_with_prefix.tf"),
					Check: resource.ComposeAggregateTestCheckFunc(
						resource.TestCheckResourceAttr("akamai_cloudwrapper_configuration.test", "id", "123"),
						resource.TestCheckResourceAttr("akamai_cloudwrapper_configuration.test", "config_name", "testname"),
						resource.TestCheckResourceAttr("akamai_cloudwrapper_configuration.test", "contract_id", "ctr_123"),
					),
					ExpectError: regexp.MustCompile("Configuration Contains Multi CDN Settings"),
				},
			},
		})

		client.AssertExpectations(t)
	})
	t.Run("contract_id prefix drift", func(t *testing.T) {
		t.Parallel()
		client := &cloudwrapper.Mock{}

		configuration := cloudwrapper.CreateConfigurationRequest{
			Body: cloudwrapper.CreateConfigurationRequestBody{
				Comments:   "test",
				ContractID: "ctr_123",
				Locations: []cloudwrapper.ConfigLocationReq{
					{
						Comments:      "test",
						TrafficTypeID: 1,
						Capacity: cloudwrapper.Capacity{
							Value: 1,
							Unit:  cloudwrapper.UnitGB,
						},
					},
				},
				NotificationEmails: []string{"test@akamai.com"},
				ConfigName:         "testname",
				PropertyIDs:        []string{"200200200"},
				RetainIdleObjects:  false,
			},
		}

		expecter := newExpecter(t, client)
		expecter.ExpectCreate(configuration)
		expecter.ExpectRefresh()

		expecter.ExpectDriftRefresh(&cloudwrapper.Configuration{
			Comments:   "test",
			ContractID: "123",
			ConfigID:   123,
			Locations: []cloudwrapper.ConfigLocationResp{
				{
					Comments:      "test",
					TrafficTypeID: 1,
					Capacity: cloudwrapper.Capacity{
						Value: 1,
						Unit:  cloudwrapper.UnitGB,
					},
				},
			},
			Status:             "SAVED",
			ConfigName:         "testname",
			NotificationEmails: []string{"test@akamai.com"},
			PropertyIDs:        []string{"200200200"},
			RetainIdleObjects:  false,
		}, nil)

		expecter.ExpectRefresh()
		expecter.ExpectDelete()

		resource.Test(t, resource.TestCase{
			IsUnitTest:               true,
			ProtoV6ProviderFactories: newProviderFactory(withMockClient(client)),
			Steps: []resource.TestStep{
				{
					Config: testutils.LoadFixtureString(t, "testdata/TestResConfiguration/create.tf"),
					Check: resource.ComposeAggregateTestCheckFunc(
						resource.TestCheckResourceAttr("akamai_cloudwrapper_configuration.test", "id", "123"),
						resource.TestCheckResourceAttr("akamai_cloudwrapper_configuration.test", "config_name", "testname"),
						resource.TestCheckResourceAttr("akamai_cloudwrapper_configuration.test", "contract_id", "ctr_123"),
					),
				},
				{
					Config: testutils.LoadFixtureString(t, "testdata/TestResConfiguration/create.tf"),
					Check: resource.ComposeAggregateTestCheckFunc(
						resource.TestCheckResourceAttr("akamai_cloudwrapper_configuration.test", "id", "123"),
						resource.TestCheckResourceAttr("akamai_cloudwrapper_configuration.test", "config_name", "testname"),
						resource.TestCheckResourceAttr("akamai_cloudwrapper_configuration.test", "contract_id", "ctr_123"),
					),
				},
			},
		})

		client.AssertExpectations(t)
	})
	t.Run("expect missing required errors", func(t *testing.T) {
		t.Parallel()
		tests := map[string]struct {
			config    string
			expectErr *regexp.Regexp
		}{
			"missing comments": {
				config:    testutils.LoadFixtureString(t, "testdata/TestResConfiguration/missing_required.tf"),
				expectErr: regexp.MustCompile(`The argument "comments" is required, but no definition was found`),
			},
			"missing config_name": {
				config:    testutils.LoadFixtureString(t, "testdata/TestResConfiguration/missing_required.tf"),
				expectErr: regexp.MustCompile(`The argument "config_name" is required, but no definition was found`),
			},
			"missing property_ids": {
				config:    testutils.LoadFixtureString(t, "testdata/TestResConfiguration/missing_required.tf"),
				expectErr: regexp.MustCompile(`The argument "property_ids" is required, but no definition was found`),
			},
			"missing contract_id": {
				config:    testutils.LoadFixtureString(t, "testdata/TestResConfiguration/missing_required.tf"),
				expectErr: regexp.MustCompile(`The argument "contract_id" is required, but no definition was found`),
			},
			"missing location": {
				config:    testutils.LoadFixtureString(t, "testdata/TestResConfiguration/missing_location.tf"),
				expectErr: regexp.MustCompile("Block location must have a configuration value as the provider has marked it\nas required"),
			},
			"missing location comments": {
				config:    testutils.LoadFixtureString(t, "testdata/TestResConfiguration/missing_location_comments.tf"),
				expectErr: regexp.MustCompile(`location {\n\nThe argument "comments" is required, but no definition was found.`),
			},
			"missing location traffic_type_id": {
				config:    testutils.LoadFixtureString(t, "testdata/TestResConfiguration/missing_required.tf"),
				expectErr: regexp.MustCompile(`location {\n\nThe argument "traffic_type_id" is required, but no definition was found.`),
			},
			"missing location capacity block": {
				config:    testutils.LoadFixtureString(t, "testdata/TestResConfiguration/missing_capacity.tf"),
				expectErr: regexp.MustCompile(`Block\nlocation\[.+\]\.capacity\nmust have a configuration value as the provider has marked it as required`),
			},
			"missing location capacity value": {
				config:    testutils.LoadFixtureString(t, "testdata/TestResConfiguration/missing_required.tf"),
				expectErr: regexp.MustCompile(`capacity {\n\nThe argument "value" is required, but no definition was found.`),
			},
			"missing location capacity unit": {
				config:    testutils.LoadFixtureString(t, "testdata/TestResConfiguration/missing_required.tf"),
				expectErr: regexp.MustCompile(`capacity {\n\nThe argument "unit" is required, but no definition was found.`),
			},
		}
		fact := newProviderFactory()

		for name, tc := range tests {
			t.Run(name, func(t *testing.T) {
				t.Parallel()
				resource.Test(t, resource.TestCase{
					IsUnitTest:               true,
					ProtoV6ProviderFactories: fact,
					Steps: []resource.TestStep{
						{
							Config:      tc.config,
							ExpectError: tc.expectErr,
						},
					},
				})

			})
		}
	})
}

type expecter struct {
	t      *testing.T
	client *cloudwrapper.Mock
	config *cloudwrapper.Configuration
}

func newExpecter(t *testing.T, client *cloudwrapper.Mock) expecter {
	return expecter{
		t:      t,
		client: client,
	}
}

func (e *expecter) CloneConfig() *cloudwrapper.Configuration {
	var config cloudwrapper.Configuration
	err := copier.CopyWithOption(&config, e.config, copier.Option{DeepCopy: true})
	assert.NoError(e.t, err)
	return &config
}

func (e *expecter) applyCreate(req cloudwrapper.CreateConfigurationRequest) {
	notificationEmails := []string{"generated@akamai.com"} // email is genereated when not provided
	if len(req.Body.NotificationEmails) > 0 {
		notificationEmails = req.Body.NotificationEmails
	}

	locations := make([]cloudwrapper.ConfigLocationResp, 0, len(req.Body.Locations))
	for _, loc := range req.Body.Locations {
		locations = append(locations, cloudwrapper.ConfigLocationResp{
			Comments:      loc.Comments,
			TrafficTypeID: loc.TrafficTypeID,
			Capacity:      loc.Capacity,
			MapName:       "cw-s-use-live",
		})
	}

	e.config = &cloudwrapper.Configuration{
		CapacityAlertsThreshold: req.Body.CapacityAlertsThreshold,
		Comments:                req.Body.Comments,
		ContractID:              req.Body.ContractID,
		ConfigID:                123,
		Locations:               locations,
		MultiCDNSettings:        req.Body.MultiCDNSettings,
		Status:                  "SAVED",
		ConfigName:              req.Body.ConfigName,
		LastUpdatedBy:           "jondoe",
		LastUpdatedDate:         "2022-02-02T02:22:22Z",
		NotificationEmails:      notificationEmails,
		PropertyIDs:             req.Body.PropertyIDs,
		RetainIdleObjects:       req.Body.RetainIdleObjects,
	}

}

func (e *expecter) ExpectCreate(req cloudwrapper.CreateConfigurationRequest) {
	e.applyCreate(req)
	e.client.On("CreateConfiguration", testutils.MockContext, req).Return(e.CloneConfig(), nil).Once()
}

func (e *expecter) ExpectRefresh() {
	e.client.On("GetConfiguration", testutils.MockContext, cloudwrapper.GetConfigurationRequest{
		ConfigID: 123,
	}).Return(e.CloneConfig(), nil).Once()
}

func (e *expecter) ExpectDriftRefresh(config *cloudwrapper.Configuration, err error) {
	e.client.On("GetConfiguration", testutils.MockContext, cloudwrapper.GetConfigurationRequest{
		ConfigID: e.config.ConfigID,
	}).Return(config, err).Once()
}

func (e *expecter) applyUpdate(req cloudwrapper.UpdateConfigurationRequest) {

	notificationEmails := []string{"generated@akamai.com"}
	if len(req.Body.NotificationEmails) > 0 {
		notificationEmails = req.Body.NotificationEmails
	}

	locations := make([]cloudwrapper.ConfigLocationResp, 0, len(req.Body.Locations))
	for _, loc := range req.Body.Locations {
		locations = append(locations, cloudwrapper.ConfigLocationResp{
			Comments:      loc.Comments,
			TrafficTypeID: loc.TrafficTypeID,
			Capacity:      loc.Capacity,
			MapName:       "cw-s-use-live",
		})
	}

	e.config = &cloudwrapper.Configuration{
		CapacityAlertsThreshold: req.Body.CapacityAlertsThreshold,
		Comments:                req.Body.Comments,
		ContractID:              e.config.ContractID, // no update
		ConfigID:                123,
		Locations:               locations,
		MultiCDNSettings:        req.Body.MultiCDNSettings,
		Status:                  "SAVED",
		ConfigName:              e.config.ConfigName, // no update
		LastUpdatedBy:           "jondoe",
		LastUpdatedDate:         "2022-02-02T02:22:22Z",
		NotificationEmails:      notificationEmails,
		PropertyIDs:             req.Body.PropertyIDs,
		RetainIdleObjects:       req.Body.RetainIdleObjects,
	}
}

func (e *expecter) ExpectUpdate(req cloudwrapper.UpdateConfigurationRequest) {
	e.applyUpdate(req)
	e.client.On("UpdateConfiguration", testutils.MockContext, req).Return(e.CloneConfig(), nil).Once()
}

func (e *expecter) ExpectDelete() {
	e.client.On("GetConfiguration", testutils.MockContext, cloudwrapper.GetConfigurationRequest{
		ConfigID: 123,
	}).Return(e.CloneConfig(), nil).Once()

	e.client.On("DeleteConfiguration", testutils.MockContext, cloudwrapper.DeleteConfigurationRequest{
		ConfigID: e.config.ConfigID,
	}).Return(nil).Once()

	conf := e.CloneConfig()
	conf.Status = "DELETE_IN_PROGRESS"

	e.client.On("GetConfiguration", testutils.MockContext, cloudwrapper.GetConfigurationRequest{
		ConfigID: 123,
	}).Return(conf, nil).Once()

	e.client.On("GetConfiguration", testutils.MockContext, cloudwrapper.GetConfigurationRequest{
		ConfigID: 123,
	}).Return(nil, cloudwrapper.ErrConfigurationNotFound).Once()
}
