package accountprotection

import (
	"errors"
	"regexp"
	"testing"

	apr "github.com/akamai/AkamaiOPEN-edgegrid-golang/v13/pkg/accountprotection"
	"github.com/akamai/terraform-provider-akamai/v11/internal/edgegrid"
	"github.com/akamai/terraform-provider-akamai/v11/pkg/common/testutils"
	"github.com/hashicorp/terraform-plugin-testing/helper/resource"
	"github.com/hashicorp/terraform-plugin-testing/knownvalue"
	"github.com/hashicorp/terraform-plugin-testing/statecheck"
	"github.com/hashicorp/terraform-plugin-testing/tfjsonpath"
)

func TestDataUserAllowList(t *testing.T) {
	t.Parallel()
	expectedRequest := apr.GetUserAllowListIDRequest{ConfigID: 43253, Version: 15}
	apiResponse := map[string]any{
		"metadata": map[string]any{
			"configId":      43253,
			"configVersion": 15,
		},
		"userAllowListId": "mytestlist123",
	}
	expectedJSON := `{
		"metadata": {
			"configId": 43253,
			"configVersion": 15
		},
		"userAllowListId": "mytestlist123"
	}`
	apiErr := errors.New("failed to get user allow list id")

	tests := map[string]struct {
		setupMock func(*apr.Mock)
		steps     []resource.TestStep
	}{
		"happy path": {
			setupMock: func(m *apr.Mock) {
				m.On("GetUserAllowListID", testutils.MockContext, expectedRequest).Return(apiResponse, nil)
			},
			steps: []resource.TestStep{
				{
					Config: testutils.LoadFixtureString(t, "testdata/TestDataUserAllowList/basic.tf"),
					ConfigStateChecks: []statecheck.StateCheck{
						statecheck.ExpectKnownValue(
							"data.akamai_apr_user_allow_list.test",
							tfjsonpath.New("json"),
							knownvalue.StringExact(compactJSON(expectedJSON)),
						),
					},
				},
			},
		},
		"error invalid config id": {
			setupMock: func(_ *apr.Mock) {},
			steps: []resource.TestStep{
				{
					Config:      testutils.LoadFixtureString(t, "testdata/TestDataUserAllowList/invalid_config_id.tf"),
					ExpectError: regexp.MustCompile("Inappropriate value for attribute \"config_id\": a number is required"),
				},
			},
		},
		"error api response": {
			setupMock: func(m *apr.Mock) {
				m.On("GetUserAllowListID", testutils.MockContext, expectedRequest).Return(nil, apiErr)
			},
			steps: []resource.TestStep{
				{
					Config:      testutils.LoadFixtureString(t, "testdata/TestDataUserAllowList/basic.tf"),
					ExpectError: regexp.MustCompile(apiErr.Error()),
				},
			},
		},
	}

	for name, test := range tests {
		t.Run(name, func(t *testing.T) {
			t.Parallel()
			client := edgegrid.NewTestClient()
			mockGetConfigVersion(client.APPSEC)
			test.setupMock(client.AccountProtection)

			resource.UnitTest(t, resource.TestCase{
				ProtoV6ProviderFactories: testutils.NewTestProtoV6ProviderFactory(client, NewSubprovider()),
				Steps:                    test.steps,
			})

			client.AccountProtection.AssertExpectations(t)
		})
	}
}
