package botman

import (
	"testing"

	"github.com/akamai/AkamaiOPEN-edgegrid-golang/v12/pkg/botman"
	"github.com/akamai/terraform-provider-akamai/v9/pkg/common/testutils"
	"github.com/hashicorp/terraform-plugin-testing/helper/resource"
)

func TestDataClientSideSecurity(t *testing.T) {
	t.Run("DataClientSideSecurity", func(t *testing.T) {

		mockedBotmanClient := &botman.Mock{}
		response := map[string]interface{}{"testKey": "testValue3"}
		expectedJSON := `{"testKey":"testValue3"}`
		mockedBotmanClient.On("GetClientSideSecurity",
			testutils.MockContext,
			botman.GetClientSideSecurityRequest{ConfigID: 43253, Version: 15},
		).Return(response, nil)

		useClient(mockedBotmanClient, func() {

			resource.Test(t, resource.TestCase{
				IsUnitTest:               true,
				ProtoV6ProviderFactories: testutils.NewProtoV6ProviderFactory(NewSubprovider()),
				Steps: []resource.TestStep{
					{
						Config: testutils.LoadFixtureString(t, "testdata/TestDataClientSideSecurity/basic.tf"),
						Check: resource.ComposeAggregateTestCheckFunc(
							resource.TestCheckResourceAttr("data.akamai_botman_client_side_security.test", "json", compactJSON(expectedJSON))),
					},
				},
			})
		})

		mockedBotmanClient.AssertExpectations(t)
	})
}
