package cloudwrapper

import (
	"regexp"
	"testing"
	"time"

	"github.com/akamai/AkamaiOPEN-edgegrid-golang/v10/pkg/cloudwrapper"
	"github.com/akamai/terraform-provider-akamai/v8/pkg/common/testutils"
	"github.com/hashicorp/terraform-plugin-testing/helper/resource"
	"github.com/stretchr/testify/mock"
)

func TestActivation(t *testing.T) {
	t.Parallel()
	tests := map[string]struct {
		init  func() *cloudwrapper.Mock
		steps []resource.TestStep
	}{
		"activation lifecycle": {
			init: func() *cloudwrapper.Mock {
				client := &cloudwrapper.Mock{}

				mockActivateConfig(client, []int{123}, nil).Once()
				//not yet activated
				mockGetConfiguration(client, 123, cloudwrapper.StatusInProgress, "location comment").Once()
				//activated
				mockGetConfiguration(client, 123, cloudwrapper.StatusActive, "location comment").Times(3)
				//refresh after modifying configuration
				mockGetConfiguration(client, 123, cloudwrapper.StatusSaved, "other comment").Once()
				mockActivateConfig(client, []int{123}, nil).Once()
				mockGetConfiguration(client, 123, cloudwrapper.StatusInProgress, "other comment").Once()
				mockGetConfiguration(client, 123, cloudwrapper.StatusActive, "other comment").Times(3)
				return client
			},
			steps: []resource.TestStep{
				{
					Config: testutils.LoadFixtureString(t, "testdata/TestResActivation/create.tf"),
					Check: resource.ComposeAggregateTestCheckFunc(
						resource.TestCheckResourceAttr("akamai_cloudwrapper_activation.act", "id", "akamai_cloudwrapper_activation"),
						resource.TestCheckResourceAttr("akamai_cloudwrapper_activation.act", "config_id", "123"),
						resource.TestCheckResourceAttr("akamai_cloudwrapper_activation.act", "revision", "5fe7963eb7270e69c5e8"),
					),
				},
				{
					Config: testutils.LoadFixtureString(t, "testdata/TestResActivation/update.tf"),
					Check: resource.ComposeAggregateTestCheckFunc(
						resource.TestCheckResourceAttr("akamai_cloudwrapper_activation.act", "id", "akamai_cloudwrapper_activation"),
						resource.TestCheckResourceAttr("akamai_cloudwrapper_activation.act", "config_id", "123"),
						resource.TestCheckResourceAttr("akamai_cloudwrapper_activation.act", "revision", "8b92934d68d69621153c"),
					),
				},
			},
		},
		"import": {
			init: func() *cloudwrapper.Mock {
				client := &cloudwrapper.Mock{}

				mockActivateConfig(client, []int{123}, nil).Once()
				mockGetConfiguration(client, 123, cloudwrapper.StatusActive, "location comment").Times(5)
				return client
			},
			steps: []resource.TestStep{
				{
					Config: testutils.LoadFixtureString(t, "testdata/TestResActivation/create.tf"),
					Check: resource.ComposeAggregateTestCheckFunc(
						resource.TestCheckResourceAttr("akamai_cloudwrapper_activation.act", "id", "akamai_cloudwrapper_activation"),
						resource.TestCheckResourceAttr("akamai_cloudwrapper_activation.act", "config_id", "123"),
						resource.TestCheckResourceAttr("akamai_cloudwrapper_activation.act", "revision", "5fe7963eb7270e69c5e8"),
					),
				},
				{
					ImportState:   true,
					ImportStateId: "123",
					ResourceName:  "akamai_cloudwrapper_activation.act",
					Check: resource.ComposeAggregateTestCheckFunc(
						resource.TestCheckResourceAttr("akamai_cloudwrapper_activation.act", "id", "akamai_cloudwrapper_activation"),
						resource.TestCheckResourceAttr("akamai_cloudwrapper_activation.act", "config_id", "123"),
						resource.TestCheckResourceAttr("akamai_cloudwrapper_activation.act", "revision", "5fe7963eb7270e69c5e8"),
					),
				},
			},
		},
		"import of inactive config": {
			init: func() *cloudwrapper.Mock {
				client := &cloudwrapper.Mock{}

				mockActivateConfig(client, []int{123}, nil).Once()
				mockGetConfiguration(client, 123, cloudwrapper.StatusActive, "location comment").Times(3)
				mockGetConfiguration(client, 123, cloudwrapper.StatusFailed, "location comment").Once()
				return client
			},
			steps: []resource.TestStep{
				{
					Config: testutils.LoadFixtureString(t, "testdata/TestResActivation/create.tf"),
					Check: resource.ComposeAggregateTestCheckFunc(
						resource.TestCheckResourceAttr("akamai_cloudwrapper_activation.act", "id", "akamai_cloudwrapper_activation"),
						resource.TestCheckResourceAttr("akamai_cloudwrapper_activation.act", "config_id", "123"),
						resource.TestCheckResourceAttr("akamai_cloudwrapper_activation.act", "revision", "5fe7963eb7270e69c5e8"),
					),
				},
				{
					ImportState:   true,
					ImportStateId: "123",
					ResourceName:  "akamai_cloudwrapper_activation.act",
					ExpectError:   regexp.MustCompile("configuration must be active prior to import; activate configuration instead"),
				},
			},
		},
		"force new on config_id": {
			init: func() *cloudwrapper.Mock {
				client := &cloudwrapper.Mock{}

				mockActivateConfig(client, []int{123}, nil).Once()
				mockGetConfiguration(client, 123, cloudwrapper.StatusActive, "location comment").Times(4)
				mockActivateConfig(client, []int{321}, nil).Once()
				mockGetConfiguration(client, 321, cloudwrapper.StatusActive, "location comment").Times(3)
				return client
			},
			steps: []resource.TestStep{
				{
					Config: testutils.LoadFixtureString(t, "testdata/TestResActivation/create.tf"),
					Check: resource.ComposeAggregateTestCheckFunc(
						resource.TestCheckResourceAttr("akamai_cloudwrapper_activation.act", "id", "akamai_cloudwrapper_activation"),
						resource.TestCheckResourceAttr("akamai_cloudwrapper_activation.act", "config_id", "123"),
						resource.TestCheckResourceAttr("akamai_cloudwrapper_activation.act", "revision", "5fe7963eb7270e69c5e8"),
					),
				},
				{
					Config: testutils.LoadFixtureString(t, "testdata/TestResActivation/update_forcenew.tf"),
					Check: resource.ComposeAggregateTestCheckFunc(
						resource.TestCheckResourceAttr("akamai_cloudwrapper_activation.act", "id", "akamai_cloudwrapper_activation"),
						resource.TestCheckResourceAttr("akamai_cloudwrapper_activation.act", "config_id", "321"),
						resource.TestCheckResourceAttr("akamai_cloudwrapper_activation.act", "revision", "5fe7963eb7270e69c5e8"),
					),
				},
			},
		},
		"missing required fields": {
			init: func() *cloudwrapper.Mock {
				return &cloudwrapper.Mock{}
			},
			steps: []resource.TestStep{
				{
					Config:      testutils.LoadFixtureString(t, "testdata/TestResActivation/create_missing_required.tf"),
					ExpectError: regexp.MustCompile(`The argument "revision" is required, but no definition was found.`),
				},
			},
		},
		"timeout on create": {
			init: func() *cloudwrapper.Mock {
				client := &cloudwrapper.Mock{}

				mockActivateConfig(client, []int{123}, nil).Once()
				mockGetConfiguration(client, 123, cloudwrapper.StatusInProgress, "location comment") //timeout sometimes triggers after 2 calls, sometimes after 3 - no Times(x)

				return client
			},
			steps: []resource.TestStep{
				{
					Config:      testutils.LoadFixtureString(t, "testdata/TestResActivation/create_timeout.tf"),
					ExpectError: regexp.MustCompile(`Reached Activation Timeout`),
				},
			},
		},
		"timeout on update": {
			init: func() *cloudwrapper.Mock {
				client := &cloudwrapper.Mock{}

				mockActivateConfig(client, []int{123}, nil).Once()
				mockGetConfiguration(client, 123, cloudwrapper.StatusActive, "location comment").Times(3)
				mockGetConfiguration(client, 123, cloudwrapper.StatusSaved, "other comment").Once()
				mockActivateConfig(client, []int{123}, nil).Once()
				mockGetConfiguration(client, 123, cloudwrapper.StatusInProgress, "other comment")

				return client
			},
			steps: []resource.TestStep{
				{
					Config: testutils.LoadFixtureString(t, "testdata/TestResActivation/create_timeout.tf"),
					Check: resource.ComposeAggregateTestCheckFunc(
						resource.TestCheckResourceAttr("akamai_cloudwrapper_activation.act", "id", "akamai_cloudwrapper_activation"),
						resource.TestCheckResourceAttr("akamai_cloudwrapper_activation.act", "config_id", "123"),
						resource.TestCheckResourceAttr("akamai_cloudwrapper_activation.act", "revision", "5fe7963eb7270e69c5e8"),
					),
				},
				{
					Config:      testutils.LoadFixtureString(t, "testdata/TestResActivation/update_timeout.tf"),
					ExpectError: regexp.MustCompile(`Reached Activation Timeout`),
				},
			},
		},
	}

	for name, test := range tests {
		t.Run(name, func(t *testing.T) {
			t.Parallel()
			client := test.init()
			resource.Test(t, resource.TestCase{
				IsUnitTest:               true,
				ProtoV6ProviderFactories: newProviderFactory(withMockClient(client), withInterval(time.Second)),
				Steps:                    test.steps,
			})
			client.AssertExpectations(t)
		})
	}
}

func mockActivateConfig(client *cloudwrapper.Mock, configIDs []int, err error) *mock.Call {
	return client.On("ActivateConfiguration", testutils.MockContext, cloudwrapper.ActivateConfigurationRequest{ConfigurationIDs: configIDs}).Return(err)
}

func mockGetConfiguration(client *cloudwrapper.Mock, configID int, returnStatus cloudwrapper.StatusType, locationComment string) *mock.Call {
	return client.On("GetConfiguration", testutils.MockContext, cloudwrapper.GetConfigurationRequest{ConfigID: int64(configID)}).
		Return(&cloudwrapper.Configuration{
			Comments:   "some comment",
			ConfigName: "testconfig",
			ContractID: "1-CTRACT",
			ConfigID:   int64(configID),
			Locations: []cloudwrapper.ConfigLocationResp{
				{
					Comments:      locationComment,
					TrafficTypeID: 2,
					Capacity: cloudwrapper.Capacity{
						Unit:  "GB",
						Value: 1,
					},
				},
			},
			Status:             returnStatus,
			NotificationEmails: []string{"test@test.com"},
			PropertyIDs:        []string{"1234567"},
		}, nil)
}
