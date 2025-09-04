package appsec

import (
	"encoding/json"
	"fmt"
	"regexp"
	"testing"

	"github.com/akamai/AkamaiOPEN-edgegrid-golang/v12/pkg/appsec"
	"github.com/akamai/terraform-provider-akamai/v9/pkg/common/testutils"
	"github.com/hashicorp/terraform-plugin-testing/helper/resource"
	"github.com/hashicorp/terraform-plugin-testing/knownvalue"
	"github.com/hashicorp/terraform-plugin-testing/statecheck"
	"github.com/hashicorp/terraform-plugin-testing/tfjsonpath"
	"github.com/stretchr/testify/require"
)

func TestAkamaiConfiguration_res_basic(t *testing.T) {
	t.Run("match by Configuration ID", func(t *testing.T) {
		client := &appsec.Mock{}

		createConfigResponse := appsec.CreateConfigurationResponse{}
		err := json.Unmarshal(testutils.LoadFixtureBytes(t, "testdata/TestResConfiguration/ConfigurationCreate.json"), &createConfigResponse)
		require.NoError(t, err)

		readConfigResponse := appsec.GetConfigurationResponse{}
		err = json.Unmarshal(testutils.LoadFixtureBytes(t, "testdata/TestResConfiguration/Configuration.json"), &readConfigResponse)
		require.NoError(t, err)

		deleteConfigResponse := appsec.RemoveConfigurationResponse{}
		err = json.Unmarshal(testutils.LoadFixtureBytes(t, "testdata/TestResConfiguration/Configuration.json"), &deleteConfigResponse)
		require.NoError(t, err)

		getConfigurationVersionsResponse := appsec.GetConfigurationVersionsResponse{}
		err = json.Unmarshal(testutils.LoadFixtureBytes(t, "testdata/TestResConfiguration/ConfigurationVersions.json"), &getConfigurationVersionsResponse)
		require.NoError(t, err)

		getSelectedHostnamesResponse := appsec.GetSelectedHostnamesResponse{}
		err = json.Unmarshal(testutils.LoadFixtureBytes(t, "testdata/TestResConfiguration/SelectedHostname.json"), &getSelectedHostnamesResponse)
		require.NoError(t, err)

		client.On("GetSelectedHostnames",
			testutils.MockContext,
			appsec.GetSelectedHostnamesRequest{ConfigID: 43253, Version: 7},
		).Return(&getSelectedHostnamesResponse, nil)

		client.On("CreateConfiguration",
			testutils.MockContext,
			appsec.CreateConfigurationRequest{Name: "Akamai Tools", Description: "Akamai Tools", ContractID: "C-1FRYVV3", GroupID: 64867, Hostnames: []string{"rinaldi.sandbox.akamaideveloper.com", "sujala.sandbox.akamaideveloper.com"}},
		).Return(&createConfigResponse, nil)

		client.On("GetConfiguration",
			testutils.MockContext,
			appsec.GetConfigurationRequest{ConfigID: 43253},
		).Return(&readConfigResponse, nil)

		client.On("RemoveConfiguration",
			testutils.MockContext,
			appsec.RemoveConfigurationRequest{ConfigID: 43253},
		).Return(&deleteConfigResponse, nil)

		client.On("GetConfigurationVersions",
			testutils.MockContext,
			appsec.GetConfigurationVersionsRequest{ConfigID: 43253},
		).Return(&getConfigurationVersionsResponse, nil)

		useClient(client, func() {
			resource.Test(t, resource.TestCase{
				IsUnitTest:               true,
				ProtoV6ProviderFactories: testutils.NewProtoV6ProviderFactory(NewSubprovider()),
				Steps: []resource.TestStep{
					{
						Config: testutils.LoadFixtureString(t, "testdata/TestResConfiguration/match_by_id.tf"),
						Check: resource.ComposeAggregateTestCheckFunc(
							resource.TestCheckResourceAttr("akamai_appsec_configuration.test", "id", "43253"),
						),
					},
				},
			})
		})

		client.AssertExpectations(t)
	})
}

func TestAkamaiConfiguration_res_error_updating_configuration(t *testing.T) {
	t.Run("match by Configuration ID", func(t *testing.T) {
		client := &appsec.Mock{}

		createConfigResponse := appsec.CreateConfigurationResponse{}
		err := json.Unmarshal(testutils.LoadFixtureBytes(t, "testdata/TestResConfiguration/ConfigurationCreate.json"), &createConfigResponse)
		require.NoError(t, err)

		readConfigResponse := appsec.GetConfigurationResponse{}
		err = json.Unmarshal(testutils.LoadFixtureBytes(t, "testdata/TestResConfiguration/Configuration.json"), &readConfigResponse)
		require.NoError(t, err)

		deleteConfigResponse := appsec.RemoveConfigurationResponse{}
		err = json.Unmarshal(testutils.LoadFixtureBytes(t, "testdata/TestResConfiguration/Configuration.json"), &deleteConfigResponse)
		require.NoError(t, err)

		getConfigurationVersionsResponse := appsec.GetConfigurationVersionsResponse{}
		err = json.Unmarshal(testutils.LoadFixtureBytes(t, "testdata/TestResConfiguration/ConfigurationVersions.json"), &getConfigurationVersionsResponse)
		require.NoError(t, err)

		hns := appsec.GetSelectedHostnamesResponse{}
		err = json.Unmarshal(testutils.LoadFixtureBytes(t, "testdata/TestResConfiguration/SelectedHostname.json"), &hns)
		require.NoError(t, err)

		client.On("GetSelectedHostnames",
			testutils.MockContext,
			appsec.GetSelectedHostnamesRequest{ConfigID: 43253, Version: 7},
		).Return(&hns, nil)

		client.On("CreateConfiguration",
			testutils.MockContext,
			appsec.CreateConfigurationRequest{Name: "Akamai Tools", Description: "Akamai Tools", ContractID: "C-1FRYVV3", GroupID: 64867, Hostnames: []string{"rinaldi.sandbox.akamaideveloper.com", "sujala.sandbox.akamaideveloper.com"}},
		).Return(&createConfigResponse, nil)

		client.On("GetConfiguration",
			testutils.MockContext,
			appsec.GetConfigurationRequest{ConfigID: 43253},
		).Return(&readConfigResponse, nil)

		client.On("UpdateConfiguration",
			testutils.MockContext,
			appsec.UpdateConfigurationRequest{ConfigID: 43253, Name: "Akamai Tools", Description: "Akamai Tools"},
		).Return(nil, fmt.Errorf("UpdateConfiguration failed"))

		client.On("RemoveConfiguration",
			testutils.MockContext,
			appsec.RemoveConfigurationRequest{ConfigID: 43253},
		).Return(&deleteConfigResponse, nil)

		client.On("GetConfigurationVersions",
			testutils.MockContext,
			appsec.GetConfigurationVersionsRequest{ConfigID: 43253},
		).Return(&getConfigurationVersionsResponse, nil)

		useClient(client, func() {
			resource.Test(t, resource.TestCase{
				IsUnitTest:               true,
				ProtoV6ProviderFactories: testutils.NewProtoV6ProviderFactory(NewSubprovider()),
				Steps: []resource.TestStep{
					{
						Config: testutils.LoadFixtureString(t, "testdata/TestResConfiguration/match_by_id.tf"),
						Check: resource.ComposeAggregateTestCheckFunc(
							resource.TestCheckResourceAttr("akamai_appsec_configuration.test", "id", "43253"),
						),
					},
					{
						Config: testutils.LoadFixtureString(t, "testdata/TestResConfiguration/modify_contract.tf"),
						Check: resource.ComposeAggregateTestCheckFunc(
							resource.TestCheckResourceAttr("akamai_appsec_configuration.test", "id", "43253"),
						),
						ExpectError: regexp.MustCompile(`UpdateConfiguration failed`),
					},
				},
			})
		})

		client.AssertExpectations(t)
	})
}

func TestAkamaiConfiguration_Create_txt_group_id(t *testing.T) {
	t.Run("Config-Create-with-non-numeric-GroupID", func(t *testing.T) {
		client := appsec.Mock{}

		setGetConfiguration(&client, t)
		setGetSelectedHostnames(&client, t)
		setGetConfigurationVersions(&client, t)
		setRemoveConfiguration(&client, t)
		setCreateConfiguration(&client, t)

		// [Create-Configuration] : 'group_id' is Not-Numeric
		tfCONFIG := testutils.LoadFixtureString(t, "testdata/TestResConfiguration/match_by_prefixed_group_id_create.tf")

		useClient(&client, func() {
			resource.Test(t, resource.TestCase{
				IsUnitTest:               true,
				ProtoV6ProviderFactories: testutils.NewProtoV6ProviderFactory(NewSubprovider()),
				Steps: []resource.TestStep{
					{
						Config: tfCONFIG,
						// 1) verify [terraform plan] will indicate '~ update in-place' of ("64867" -> "grp_64867")
						ConfigStateChecks: []statecheck.StateCheck{
							statecheck.ExpectKnownValue(
								"akamai_appsec_configuration.test",
								tfjsonpath.New("group_id"),
								knownvalue.StringExact("64867"),
							),
						},
						// 2) Finally, determining that a Plan: its 'state' is No-Difference from 'config-in-request'.
						ExpectNonEmptyPlan: false,
					},
				},
			})
		})
		client.AssertExpectations(t)
	})
}

func TestAkamaiConfiguration_Clone_txt_group_id(t *testing.T) {
	t.Run("Config-Clone-with-non-numeric-GroupID", func(t *testing.T) {
		client := appsec.Mock{}

		setGetConfiguration(&client, t)
		setGetSelectedHostnames(&client, t)
		setGetConfigurationVersions(&client, t)
		setRemoveConfiguration(&client, t)
		setCreateConfigurationClone(&client, t)

		// [Clone-Configuration] : 'group_id' is Not-Numeric
		tfCONFIG := testutils.LoadFixtureString(t, "testdata/TestResConfiguration/match_by_prefixed_group_id_clone.tf")

		useClient(&client, func() {
			resource.Test(t, resource.TestCase{
				IsUnitTest:               true,
				ProtoV6ProviderFactories: testutils.NewProtoV6ProviderFactory(NewSubprovider()),
				Steps: []resource.TestStep{
					{
						Config: tfCONFIG,
						// 1) verify [terraform plan] will indicate '~ update in-place' of ("64867" -> "grp_64867")
						ConfigStateChecks: []statecheck.StateCheck{
							statecheck.ExpectKnownValue(
								"akamai_appsec_configuration.test",
								tfjsonpath.New("group_id"),
								knownvalue.StringExact("64867"),
							),
						},
						// 2) Finally, determining that a Plan: its 'state' is No-Difference from 'config-in-request'.
						ExpectNonEmptyPlan: false,
					},
				},
			})
		})
		client.AssertExpectations(t)
	})
}

func setGetConfiguration(mock *appsec.Mock, test *testing.T) {
	obj := appsec.GetConfigurationResponse{}
	err := json.Unmarshal(testutils.LoadFixtureBytes(test, "testdata/TestResConfiguration/Configuration.json"), &obj)
	require.NoError(test, err)

	mock.On("GetConfiguration",
		testutils.MockContext,
		appsec.GetConfigurationRequest{ConfigID: 43253},
	).Return(&obj, nil)
}

func setGetSelectedHostnames(mock *appsec.Mock, test *testing.T) {
	obj := appsec.GetSelectedHostnamesResponse{}
	err := json.Unmarshal(testutils.LoadFixtureBytes(test, "testdata/TestResConfiguration/SelectedHostname.json"), &obj)
	require.NoError(test, err)

	mock.On("GetSelectedHostnames",
		testutils.MockContext,
		appsec.GetSelectedHostnamesRequest{ConfigID: 43253, Version: 7},
	).Return(&obj, nil)
}

func setGetConfigurationVersions(mock *appsec.Mock, test *testing.T) {
	obj := appsec.GetConfigurationVersionsResponse{}
	err := json.Unmarshal(testutils.LoadFixtureBytes(test, "testdata/TestResConfiguration/ConfigurationVersions.json"), &obj)
	require.NoError(test, err)

	mock.On("GetConfigurationVersions",
		testutils.MockContext,
		appsec.GetConfigurationVersionsRequest{ConfigID: 43253},
	).Return(&obj, nil)

}

func setRemoveConfiguration(mock *appsec.Mock, test *testing.T) {
	obj := appsec.RemoveConfigurationResponse{}
	err := json.Unmarshal(testutils.LoadFixtureBytes(test, "testdata/TestResConfiguration/Configuration.json"), &obj)
	require.NoError(test, err)

	mock.On("RemoveConfiguration",
		testutils.MockContext,
		appsec.RemoveConfigurationRequest{ConfigID: 43253},
	).Return(&obj, nil)
}

func setCreateConfiguration(mock *appsec.Mock, test *testing.T) {
	obj := appsec.CreateConfigurationResponse{}
	err := json.Unmarshal(testutils.LoadFixtureBytes(test, "testdata/TestResConfiguration/ConfigurationCreate.json"), &obj)
	require.NoError(test, err)

	mock.On("CreateConfiguration",
		testutils.MockContext,
		// ie. expect Attributes match to above JSON-File.
		appsec.CreateConfigurationRequest{Name: "Akamai Tools", Description: "Akamai Tools", ContractID: "C-1FRYVV3", GroupID: 64867, Hostnames: []string{"rinaldi.sandbox.akamaideveloper.com", "sujala.sandbox.akamaideveloper.com"}},
	).Return(&obj, nil)
}

func setCreateConfigurationClone(mock *appsec.Mock, test *testing.T) {
	objClone := appsec.CreateConfigurationCloneResponse{}
	err := json.Unmarshal(testutils.LoadFixtureBytes(test, "testdata/TestResConfiguration/ConfigurationCloneFrom.json"), &objClone)
	require.NoError(test, err)

	mock.On("CreateConfigurationClone",
		testutils.MockContext,
		// ie. expect Attributes match to above JSON-File.
		appsec.CreateConfigurationCloneRequest{Name: "Akamai Tools", Description: "Akamai Tools", ContractID: "C-1FRYVV3", GroupID: 64867, Hostnames: []string{"rinaldi.sandbox.akamaideveloper.com", "sujala.sandbox.akamaideveloper.com"}, CreateFrom: struct {
			ConfigID int "json:\"configId\""
			Version  int "json:\"version\""
		}{43253, 1}},
	).Return(&objClone, nil)
}
