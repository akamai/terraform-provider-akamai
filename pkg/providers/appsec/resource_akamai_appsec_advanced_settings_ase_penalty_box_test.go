package appsec

import (
	"encoding/json"
	"reflect"
	"sort"
	"testing"

	"github.com/stretchr/testify/mock"

	"github.com/akamai/AkamaiOPEN-edgegrid-golang/v13/pkg/appsec"
	"github.com/akamai/terraform-provider-akamai/v10/internal/edgegrid"
	"github.com/akamai/terraform-provider-akamai/v10/pkg/common/testutils"
	"github.com/hashicorp/terraform-plugin-testing/helper/resource"
	"github.com/stretchr/testify/require"
)

func TestAkamaiAdvancedSettingsAsePenaltyBoxResConfig(t *testing.T) {
	t.Parallel()
	var (
		configVersion = func(configId int, client *appsec.Mock) appsec.GetConfigurationResponse {
			configResponse := appsec.GetConfigurationResponse{}
			err := json.Unmarshal(testutils.LoadFixtureBytes(t, "testdata/TestResConfiguration/LatestConfiguration.json"), &configResponse)
			require.NoError(t, err)

			client.On("GetConfiguration",
				mock.Anything,
				appsec.GetConfigurationRequest{ConfigID: configId},
			).Return(&configResponse, nil)

			return configResponse
		}

		AsePenaltyBoxRead = func(configId int, version int, client *appsec.Mock, numberOfTimes int, filePath string) {
			AsePenaltyBoxResponse := appsec.GetAdvancedSettingsAsePenaltyBoxResponse{}
			err := json.Unmarshal(testutils.LoadFixtureBytes(t, filePath), &AsePenaltyBoxResponse)
			require.NoError(t, err)

			client.On("GetAdvancedSettingsAsePenaltyBox",
				mock.Anything,
				appsec.GetAdvancedSettingsAsePenaltyBoxRequest{ConfigID: configId, Version: version},
			).Return(&AsePenaltyBoxResponse, nil).Times(numberOfTimes)

		}

		updateAsePenaltyBox = func(updateAsePenaltyBox appsec.UpdateAdvancedSettingsAsePenaltyBoxRequest, client *appsec.Mock, numberOfTimes int, filePath string) {
			updateAsePenaltyBoxResponse := appsec.UpdateAdvancedSettingsAsePenaltyBoxResponse{}
			err := json.Unmarshal(testutils.LoadFixtureBytes(t, filePath), &updateAsePenaltyBoxResponse)
			require.NoError(t, err)

			client.On("UpdateAdvancedSettingsAsePenaltyBox",
				mock.Anything,
				mock.MatchedBy(func(req appsec.UpdateAdvancedSettingsAsePenaltyBoxRequest) bool {
					if req.ConfigID != updateAsePenaltyBox.ConfigID ||
						req.Version != updateAsePenaltyBox.Version ||
						req.BlockDuration != updateAsePenaltyBox.BlockDuration {
						return false
					}
					// Compare QualificationExclusions, ignoring slice order
					want := append([]string{}, updateAsePenaltyBox.QualificationExclusions.AttackGroups...)
					got := append([]string{}, req.QualificationExclusions.AttackGroups...)
					sort.Strings(want)
					sort.Strings(got)
					if !reflect.DeepEqual(want, got) {
						return false
					}
					return reflect.DeepEqual(
						updateAsePenaltyBox.QualificationExclusions.Rules,
						req.QualificationExclusions.Rules,
					)
				}),
			).Return(&updateAsePenaltyBoxResponse, nil).Times(numberOfTimes)

		}

		removeAsePenaltyBox = func(updateAsePenaltyBox appsec.RemoveAdvancedSettingsAsePenaltyBoxRequest, client *appsec.Mock, numberOfTimes int, filePath string) {
			removeAsePenaltyBoxResponse := appsec.RemoveAdvancedSettingsAsePenaltyBoxResponse{}
			err := json.Unmarshal(testutils.LoadFixtureBytes(t, filePath), &removeAsePenaltyBoxResponse)
			require.NoError(t, err)

			client.On("RemoveAdvancedSettingsAsePenaltyBox",
				mock.Anything, updateAsePenaltyBox,
			).Return(&removeAsePenaltyBoxResponse, nil).Times(numberOfTimes)

		}
	)

	t.Run("match by AdvancedSettingsAsePenaltyBox ID", func(t *testing.T) {
		t.Parallel()
		client := edgegrid.NewTestClient()
		configResponse := configVersion(43253, client.APPSEC)

		AsePenaltyBoxRead(43253, 7, client.APPSEC, 2, "testdata/TestResAdvancedSettingsAsePenaltyBox/AdvancedSettingsAsePenaltyBox.json")

		updateAsePenaltyBoxRequest := appsec.UpdateAdvancedSettingsAsePenaltyBoxRequest{
			ConfigID:      configResponse.ID,
			Version:       configResponse.LatestVersion,
			BlockDuration: 5,
			QualificationExclusions: &appsec.QualificationExclusions{
				AttackGroups: []string{"XSS", "IN"},
				Rules:        []int{950002},
			},
		}

		updateAsePenaltyBox(updateAsePenaltyBoxRequest, client.APPSEC, 1, "testdata/TestResAdvancedSettingsAsePenaltyBox/AdvancedSettingsAsePenaltyBox.json")
		removeAsePenaltyBoxRequest := appsec.RemoveAdvancedSettingsAsePenaltyBoxRequest{
			ConfigID: configResponse.ID,
			Version:  configResponse.LatestVersion,
		}

		removeAsePenaltyBox(removeAsePenaltyBoxRequest, client.APPSEC, 1, "testdata/TestResAdvancedSettingsAsePenaltyBox/AdvancedSettingsAsePenaltyBox.json")
		mockGetConfigurationVersionDefault(client.APPSEC)
		resource.UnitTest(t, resource.TestCase{
			ProtoV6ProviderFactories: testutils.NewTestProtoV6SDKProviderFactory(client, NewSubprovider()),
			Steps: []resource.TestStep{
				{
					Config: testutils.LoadFixtureString(t, "testdata/TestResAdvancedSettingsAsePenaltyBox/match_by_id.tf"),
					Check: resource.ComposeAggregateTestCheckFunc(
						resource.TestCheckResourceAttr("akamai_appsec_advanced_settings_ase_penalty_box.test", "id", "43253"),
					),
				},
			},
		})

		client.APPSEC.AssertExpectations(t)
	})

	t.Run("import", func(t *testing.T) {
		t.Parallel()
		client := edgegrid.NewTestClient()

		configResponse := configVersion(43253, client.APPSEC)

		AsePenaltyBoxRead(configResponse.ID, configResponse.LatestVersion, client.APPSEC, 3, "testdata/TestResAdvancedSettingsAsePenaltyBox/AdvancedSettingsAsePenaltyBox.json")

		updateAsePenaltyBoxRequest := appsec.UpdateAdvancedSettingsAsePenaltyBoxRequest{
			ConfigID:      configResponse.ID,
			Version:       configResponse.LatestVersion,
			BlockDuration: 5,
			QualificationExclusions: &appsec.QualificationExclusions{
				AttackGroups: []string{"XSS", "IN"},
				Rules:        []int{950002},
			},
		}

		updateAsePenaltyBox(updateAsePenaltyBoxRequest, client.APPSEC, 1, "testdata/TestResAdvancedSettingsAsePenaltyBox/AdvancedSettingsAsePenaltyBox.json")

		removeAsePenaltyBoxRequest := appsec.RemoveAdvancedSettingsAsePenaltyBoxRequest{
			ConfigID: configResponse.ID,
			Version:  configResponse.LatestVersion,
		}

		removeAsePenaltyBox(removeAsePenaltyBoxRequest, client.APPSEC, 1, "testdata/TestResAdvancedSettingsAsePenaltyBox/AdvancedSettingsAsePenaltyBox.json")

		mockGetConfigurationVersionDefault(client.APPSEC)
		resource.UnitTest(t, resource.TestCase{
			ProtoV6ProviderFactories: testutils.NewTestProtoV6SDKProviderFactory(client, NewSubprovider()),
			Steps: []resource.TestStep{
				{
					Config: testutils.LoadFixtureString(t, "testdata/TestResAdvancedSettingsAsePenaltyBox/match_by_id.tf"),
				},
				{
					ImportState:       true,
					ImportStateVerify: true,
					ImportStateId:     "43253",
					ResourceName:      "akamai_appsec_advanced_settings_ase_penalty_box.test",
				},
			},
		})
		client.APPSEC.AssertExpectations(t)
	})
}
