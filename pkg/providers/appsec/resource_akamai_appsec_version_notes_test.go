package appsec

import (
	"encoding/json"
	"testing"

	"github.com/akamai/AkamaiOPEN-edgegrid-golang/v11/pkg/appsec"
	"github.com/akamai/terraform-provider-akamai/v8/pkg/common/testutils"
	"github.com/hashicorp/terraform-plugin-testing/helper/resource"
	"github.com/stretchr/testify/require"
)

func TestAkamaiVersionNotes_res_basic(t *testing.T) {
	t.Run("match by VersionNotes ID", func(t *testing.T) {
		client := &appsec.Mock{}

		updateVersionNotesResponse := appsec.UpdateVersionNotesResponse{}
		err := json.Unmarshal(testutils.LoadFixtureBytes(t, "testdata/TestResVersionNotes/VersionNotes.json"), &updateVersionNotesResponse)
		require.NoError(t, err)

		getVersionNotesResponse := appsec.GetVersionNotesResponse{}
		err = json.Unmarshal(testutils.LoadFixtureBytes(t, "testdata/TestResVersionNotes/VersionNotes.json"), &getVersionNotesResponse)
		require.NoError(t, err)

		config := appsec.GetConfigurationResponse{}
		err = json.Unmarshal(testutils.LoadFixtureBytes(t, "testdata/TestResConfiguration/LatestConfiguration.json"), &config)
		require.NoError(t, err)

		client.On("GetConfiguration",
			testutils.MockContext,
			appsec.GetConfigurationRequest{ConfigID: 43253},
		).Return(&config, nil)

		client.On("GetVersionNotes",
			testutils.MockContext,
			appsec.GetVersionNotesRequest{ConfigID: 43253, Version: 7},
		).Return(&getVersionNotesResponse, nil)

		client.On("UpdateVersionNotes",
			testutils.MockContext,
			appsec.UpdateVersionNotesRequest{ConfigID: 43253, Version: 7, Notes: "Test Notes"},
		).Return(&updateVersionNotesResponse, nil)

		useClient(client, func() {
			resource.Test(t, resource.TestCase{
				IsUnitTest:               true,
				ProtoV6ProviderFactories: testutils.NewProtoV6ProviderFactory(NewSubprovider()),
				Steps: []resource.TestStep{
					{
						Config: testutils.LoadFixtureString(t, "testdata/TestResVersionNotes/match_by_id.tf"),
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
