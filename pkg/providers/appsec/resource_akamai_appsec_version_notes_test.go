package appsec

import (
	"encoding/json"
	"testing"

	"github.com/akamai/AkamaiOPEN-edgegrid-golang/v2/pkg/appsec"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/resource"
	"github.com/stretchr/testify/mock"
)

func TestAccAkamaiVersionNotes_res_basic(t *testing.T) {
	t.Run("match by VersionNotes ID", func(t *testing.T) {
		client := &mockappsec{}

		cu := appsec.UpdateVersionNotesResponse{}
		expectJSU := compactJSON(loadFixtureBytes("testdata/TestResVersionNotes/VersionNotes.json"))
		json.Unmarshal([]byte(expectJSU), &cu)

		cr := appsec.GetVersionNotesResponse{}
		expectJS := compactJSON(loadFixtureBytes("testdata/TestResVersionNotes/VersionNotes.json"))
		json.Unmarshal([]byte(expectJS), &cr)

		client.On("GetVersionNotes",
			mock.Anything, // ctx is irrelevant for this test
			appsec.GetVersionNotesRequest{ConfigID: 43253, Version: 7},
		).Return(&cr, nil)

		client.On("UpdateVersionNotes",
			mock.Anything, // ctx is irrelevant for this test
			appsec.UpdateVersionNotesRequest{ConfigID: 43253, Version: 7},
		).Return(&cu, nil)

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
