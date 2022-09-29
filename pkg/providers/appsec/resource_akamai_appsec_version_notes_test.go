package appsec

import (
	"encoding/json"
	"testing"

	"github.com/akamai/AkamaiOPEN-edgegrid-golang/v2/pkg/appsec"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/resource"
	"github.com/stretchr/testify/mock"
	"github.com/stretchr/testify/require"
)

func TestAkamaiVersionNotes_res_basic(t *testing.T) {
	t.Run("match by VersionNotes ID", func(t *testing.T) {
		client := &mockappsec{}

		updateVersionNotesResponse := appsec.UpdateVersionNotesResponse{}
		err := json.Unmarshal(loadFixtureBytes("testdata/TestResVersionNotes/VersionNotes.json"), &updateVersionNotesResponse)
		require.NoError(t, err)

		getVersionNotesResponse := appsec.GetVersionNotesResponse{}
		err = json.Unmarshal(loadFixtureBytes("testdata/TestResVersionNotes/VersionNotes.json"), &getVersionNotesResponse)
		require.NoError(t, err)

		config := appsec.GetConfigurationResponse{}
		err = json.Unmarshal(loadFixtureBytes("testdata/TestResConfiguration/LatestConfiguration.json"), &config)
		require.NoError(t, err)

		client.On("GetConfiguration",
			mock.Anything,
			appsec.GetConfigurationRequest{ConfigID: 43253},
		).Return(&config, nil)

		client.On("GetVersionNotes",
			mock.Anything,
			appsec.GetVersionNotesRequest{ConfigID: 43253, Version: 7},
		).Return(&getVersionNotesResponse, nil)

		client.On("UpdateVersionNotes",
			mock.Anything,
			appsec.UpdateVersionNotesRequest{ConfigID: 43253, Version: 7, Notes: "Test Notes"},
		).Return(&updateVersionNotesResponse, nil)

		useClient(client, func() {
			resource.Test(t, resource.TestCase{
				IsUnitTest: true,
				Providers:  testAccProviders,
				Steps: []resource.TestStep{
					{
						Config: loadFixtureString("testdata/TestResVersionNotes/match_by_id.tf"),
						Check: resource.ComposeAggregateTestCheckFunc(
							resource.TestCheckResourceAttr("akamai_appsec_version_notes.test", "id", "43253"),
						),
					},
				},
			})
		})

		client.AssertExpectations(t)
	})

}
