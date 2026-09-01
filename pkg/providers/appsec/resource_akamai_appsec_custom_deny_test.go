package appsec

import (
	"encoding/json"
	"net/http"
	"testing"

	"github.com/akamai/AkamaiOPEN-edgegrid-golang/v13/pkg/appsec"
	"github.com/akamai/terraform-provider-akamai/v11/internal/edgegrid"
	"github.com/akamai/terraform-provider-akamai/v11/pkg/common/testutils"
	"github.com/hashicorp/terraform-plugin-testing/helper/resource"
	"github.com/stretchr/testify/require"
)

func TestAkamaiCustomDeny_res_basic(t *testing.T) {
	t.Parallel()
	t.Run("match by CustomDeny ID", func(t *testing.T) {
		t.Parallel()
		client := edgegrid.NewTestClient()

		configResponse := appsec.GetConfigurationResponse{}
		err := json.Unmarshal(testutils.LoadFixtureBytes(t, "testdata/TestResConfiguration/LatestConfiguration.json"), &configResponse)
		require.NoError(t, err)
		client.APPSEC.On("GetConfiguration",
			testutils.MockContext,
			appsec.GetConfigurationRequest{ConfigID: 43253},
		).Return(&configResponse, nil)

		createResponse := appsec.CreateCustomDenyResponse{}
		err = json.Unmarshal(testutils.LoadFixtureBytes(t, "testdata/TestResCustomDeny/CustomDenyCreateResponse.json"), &createResponse)
		require.NoError(t, err)
		createRequestJSON := testutils.LoadFixtureBytes(t, "testdata/TestResCustomDeny/CustomDenyWithPreventBrowserCacheTrue.json")
		client.APPSEC.On("CreateCustomDeny",
			testutils.MockContext,
			appsec.CreateCustomDenyRequest{ConfigID: 43253, Version: 7, JsonPayloadRaw: createRequestJSON},
		).Return(&createResponse, nil)

		getResponse := appsec.GetCustomDenyResponse{}
		err = json.Unmarshal(testutils.LoadFixtureBytes(t, "testdata/TestResCustomDeny/CustomDenyGetResponse.json"), &getResponse)
		require.NoError(t, err)
		client.APPSEC.On("GetCustomDeny",
			testutils.MockContext,
			appsec.GetCustomDenyRequest{ConfigID: 43253, Version: 7, ID: "deny_custom_622918"},
		).Return(&getResponse, nil).Times(3)

		updateRequestJSON := testutils.LoadFixtureBytes(t, "testdata/TestResCustomDeny/CustomDenyWithPreventBrowserCacheFalse.json")
		updateResponse := appsec.UpdateCustomDenyResponse{}
		err = json.Unmarshal(testutils.LoadFixtureBytes(t, "testdata/TestResCustomDeny/CustomDenyUpdateResponse.json"), &updateResponse)
		require.NoError(t, err)
		client.APPSEC.On("UpdateCustomDeny",
			testutils.MockContext,
			appsec.UpdateCustomDenyRequest{ConfigID: 43253, Version: 7, ID: "deny_custom_622918", JsonPayloadRaw: updateRequestJSON},
		).Return(&updateResponse, nil)

		getResponseAfterUpdate := appsec.GetCustomDenyResponse{}
		err = json.Unmarshal(testutils.LoadFixtureBytes(t, "testdata/TestResCustomDeny/CustomDenyGetResponseAfterUpdate.json"), &getResponseAfterUpdate)
		require.NoError(t, err)
		client.APPSEC.On("GetCustomDeny",
			testutils.MockContext,
			appsec.GetCustomDenyRequest{ConfigID: 43253, Version: 7, ID: "deny_custom_622918"},
		).Return(&getResponseAfterUpdate, nil).Twice()

		removeResponse := appsec.RemoveCustomDenyResponse{}
		err = json.Unmarshal(testutils.LoadFixtureBytes(t, "testdata/TestResCustomDeny/CustomDeny.json"), &removeResponse)
		require.NoError(t, err)
		client.APPSEC.On("RemoveCustomDeny",
			testutils.MockContext,
			appsec.RemoveCustomDenyRequest{ConfigID: 43253, Version: 7, ID: "deny_custom_622918"},
		).Return(&removeResponse, nil)

		mockGetConfigurationVersionDefault(client.APPSEC)
		resource.UnitTest(t, resource.TestCase{
			ProtoV6ProviderFactories: testutils.NewTestProtoV6ProviderFactory(client, NewSubprovider()),
			Steps: []resource.TestStep{
				{
					Config: testutils.LoadFixtureString(t, "testdata/TestResCustomDeny/match_by_id.tf"),
					Check: resource.ComposeAggregateTestCheckFunc(
						resource.TestCheckResourceAttr("akamai_appsec_custom_deny.test", "id", "43253:deny_custom_622918"),
					),
				},
				{
					Config: testutils.LoadFixtureString(t, "testdata/TestResCustomDeny/update_by_id.tf"),
					Check: resource.ComposeAggregateTestCheckFunc(
						resource.TestCheckResourceAttr("akamai_appsec_custom_deny.test", "id", "43253:deny_custom_622918"),
					),
				},
			},
		})

		client.APPSEC.AssertExpectations(t)
	})

	// Test case: When custom deny is not found in the latest config version, it should be removed from state
	// and recreated on the next apply.
	//
	// This simulates the following scenario:
	// 1. User creates custom deny in config version v2
	// 2. Someone creates v3 by cloning from v1 (which doesn't have the custom deny)
	// 3. Terraform refresh fails because custom deny doesn't exist in latest version (v3)
	//
	// The fix: When Read gets 404, remove resource from state so Terraform can recreate it.
	t.Run("custom deny not found in latest config version - removes from state", func(t *testing.T) {
		t.Parallel()
		client := edgegrid.NewTestClient()

		// Mock GetConfiguration - called by getLatestConfigVersion() to determine the latest version (7)
		// This is called multiple times during the test (refresh, plan, apply phases)
		configResponse := appsec.GetConfigurationResponse{}
		err := json.Unmarshal(testutils.LoadFixtureBytes(t, "testdata/TestResConfiguration/LatestConfiguration.json"), &configResponse)
		require.NoError(t, err)
		client.APPSEC.On("GetConfiguration",
			testutils.MockContext,
			appsec.GetConfigurationRequest{ConfigID: 43253},
		).Return(&configResponse, nil)

		// Load fixtures for create and get responses
		createResponse := appsec.CreateCustomDenyResponse{}
		err = json.Unmarshal(testutils.LoadFixtureBytes(t, "testdata/TestResCustomDeny/CustomDenyCreateResponse.json"), &createResponse)
		require.NoError(t, err)
		createRequestJSON := testutils.LoadFixtureBytes(t, "testdata/TestResCustomDeny/CustomDenyWithPreventBrowserCacheTrue.json")

		getResponse := appsec.GetCustomDenyResponse{}
		err = json.Unmarshal(testutils.LoadFixtureBytes(t, "testdata/TestResCustomDeny/CustomDenyGetResponse.json"), &getResponse)
		require.NoError(t, err)

		// Create 404 error to simulate custom deny not found in latest config version
		// This is what the API returns when the custom deny exists in an older version
		// but not in the latest version (e.g., v3 was cloned from v1 which didn't have it)
		notFoundError := &appsec.Error{
			StatusCode: http.StatusNotFound,
			Title:      "Not Found",
			Detail:     "Custom deny not found",
		}

		// ========== STEP 1 MOCKS ==========
		// Create: Terraform creates the custom deny resource
		client.APPSEC.On("CreateCustomDeny",
			testutils.MockContext,
			appsec.CreateCustomDenyRequest{ConfigID: 43253, Version: 7, JsonPayloadRaw: createRequestJSON},
		).Return(&createResponse, nil)

		// Read after create: Called immediately after create to populate state - succeeds
		client.APPSEC.On("GetCustomDeny",
			testutils.MockContext,
			appsec.GetCustomDenyRequest{ConfigID: 43253, Version: 7, ID: "deny_custom_622918"},
		).Return(&getResponse, nil).Once()

		// Refresh after apply: Terraform refreshes state after apply - returns 404
		// This simulates the scenario where the custom deny was removed from latest version
		// The fix handles this by calling d.SetId("") to remove from state
		client.APPSEC.On("GetCustomDeny",
			testutils.MockContext,
			appsec.GetCustomDenyRequest{ConfigID: 43253, Version: 7, ID: "deny_custom_622918"},
		).Return(nil, notFoundError).Once()

		// ========== STEP 2 MOCKS ==========
		// Since resource was removed from state in Step 1 refresh, Terraform will:
		// 1. See resource in config but not in state
		// 2. Plan a Create operation
		// 3. Call CreateCustomDeny again (reusing the mock above)
		// 4. Call GetCustomDeny to read after create - succeeds
		client.APPSEC.On("GetCustomDeny",
			testutils.MockContext,
			appsec.GetCustomDenyRequest{ConfigID: 43253, Version: 7, ID: "deny_custom_622918"},
		).Return(&getResponse, nil)

		// Cleanup: RemoveCustomDeny called when test ends
		removeResponse := appsec.RemoveCustomDenyResponse{}
		err = json.Unmarshal(testutils.LoadFixtureBytes(t, "testdata/TestResCustomDeny/CustomDeny.json"), &removeResponse)
		require.NoError(t, err)
		client.APPSEC.On("RemoveCustomDeny",
			testutils.MockContext,
			appsec.RemoveCustomDenyRequest{ConfigID: 43253, Version: 7, ID: "deny_custom_622918"},
		).Return(&removeResponse, nil).Once()

		mockGetConfigurationVersionDefault(client.APPSEC)
		resource.UnitTest(t, resource.TestCase{
			ProtoV6ProviderFactories: testutils.NewTestProtoV6ProviderFactory(client, NewSubprovider()),
			Steps: []resource.TestStep{
				{
					// Step 1: Create the resource
					// After apply, refresh will get 404 and remove resource from state
					// ExpectNonEmptyPlan=true because the post-apply refresh removes the resource,
					// so the next plan will show "create" action
					Config: testutils.LoadFixtureString(t, "testdata/TestResCustomDeny/match_by_id.tf"),
					Check: resource.ComposeAggregateTestCheckFunc(
						resource.TestCheckResourceAttr("akamai_appsec_custom_deny.test", "id", "43253:deny_custom_622918"),
					),
					ExpectNonEmptyPlan: true,
				},
				{
					// Step 2: Resource was removed from state, so Terraform recreates it
					// This verifies the fix works - resource is successfully recreated
					Config: testutils.LoadFixtureString(t, "testdata/TestResCustomDeny/match_by_id.tf"),
					Check: resource.ComposeAggregateTestCheckFunc(
						resource.TestCheckResourceAttr("akamai_appsec_custom_deny.test", "id", "43253:deny_custom_622918"),
					),
				},
			},
		})

		client.APPSEC.AssertExpectations(t)
	})

}
