package appsec

import (
	"bytes"
	"encoding/json"
	"testing"

	"github.com/akamai/AkamaiOPEN-edgegrid-golang/v9/pkg/appsec"
	"github.com/akamai/terraform-provider-akamai/v6/pkg/common/testutils"
	"github.com/hashicorp/terraform-plugin-testing/helper/resource"
	"github.com/stretchr/testify/mock"
	"github.com/stretchr/testify/require"
)

func TestAkamaiMatchTarget_res_basic(t *testing.T) {
	t.Run("match by MatchTarget ID", func(t *testing.T) {
		client := &appsec.Mock{}

		updateMatchTargetResponse := appsec.UpdateMatchTargetResponse{}
		err := json.Unmarshal(testutils.LoadFixtureBytes(t, "testdata/TestResMatchTarget/MatchTargetUpdated.json"), &updateMatchTargetResponse)
		require.NoError(t, err)

		getMatchTargetResponse := appsec.GetMatchTargetResponse{}
		err = json.Unmarshal(testutils.LoadFixtureBytes(t, "testdata/TestResMatchTarget/MatchTarget.json"), &getMatchTargetResponse)
		require.NoError(t, err)

		getMatchTargetResponseAfterUpdate := appsec.GetMatchTargetResponse{}
		err = json.Unmarshal(testutils.LoadFixtureBytes(t, "testdata/TestResMatchTarget/MatchTargetUpdated.json"), &getMatchTargetResponseAfterUpdate)
		require.NoError(t, err)

		createMatchTargetResponse := appsec.CreateMatchTargetResponse{}
		err = json.Unmarshal(testutils.LoadFixtureBytes(t, "testdata/TestResMatchTarget/MatchTargetCreated.json"), &createMatchTargetResponse)
		require.NoError(t, err)

		removeMatchTargetResponse := appsec.RemoveMatchTargetResponse{}
		err = json.Unmarshal(testutils.LoadFixtureBytes(t, "testdata/TestResMatchTarget/MatchTargetCreated.json"), &removeMatchTargetResponse)
		require.NoError(t, err)

		config := appsec.GetConfigurationResponse{}
		err = json.Unmarshal(testutils.LoadFixtureBytes(t, "testdata/TestResConfiguration/LatestConfiguration.json"), &config)
		require.NoError(t, err)

		client.On("GetConfiguration",
			mock.Anything,
			appsec.GetConfigurationRequest{ConfigID: 43253},
		).Return(&config, nil)

		client.On("GetMatchTarget",
			mock.Anything,
			appsec.GetMatchTargetRequest{ConfigID: 43253, ConfigVersion: 7, TargetID: 3008967},
		).Return(&getMatchTargetResponse, nil).Times(3)

		client.On("GetMatchTarget",
			mock.Anything,
			appsec.GetMatchTargetRequest{ConfigID: 43253, ConfigVersion: 7, TargetID: 3008967},
		).Return(&getMatchTargetResponseAfterUpdate, nil)

		createMatchTargetJSON := testutils.LoadFixtureBytes(t, "testdata/TestResMatchTarget/CreateMatchTarget.json")
		client.On("CreateMatchTarget",
			mock.Anything,
			appsec.CreateMatchTargetRequest{Type: "", ConfigID: 43253, ConfigVersion: 7, JsonPayloadRaw: createMatchTargetJSON},
		).Return(&createMatchTargetResponse, nil)

		updateMatchTargetJSON := testutils.LoadFixtureBytes(t, "testdata/TestResMatchTarget/UpdateMatchTarget.json")
		client.On("UpdateMatchTarget",
			mock.Anything,
			appsec.UpdateMatchTargetRequest{ConfigID: 43253, ConfigVersion: 7, TargetID: 3008967, JsonPayloadRaw: updateMatchTargetJSON},
		).Return(&updateMatchTargetResponse, nil)

		client.On("RemoveMatchTarget",
			mock.Anything,
			appsec.RemoveMatchTargetRequest{ConfigID: 43253, ConfigVersion: 7, TargetID: 3008967},
		).Return(&removeMatchTargetResponse, nil)

		useClient(client, func() {
			resource.Test(t, resource.TestCase{
				IsUnitTest:               true,
				ProtoV6ProviderFactories: testutils.NewProtoV6ProviderFactory(NewSubprovider()),
				Steps: []resource.TestStep{
					{
						Config: testutils.LoadFixtureString(t, "testdata/TestResMatchTarget/match_by_id.tf"),
						Check: resource.ComposeAggregateTestCheckFunc(
							resource.TestCheckResourceAttr("akamai_appsec_match_target.test", "id", "43253:3008967"),
						),
					},
					{
						Config: testutils.LoadFixtureString(t, "testdata/TestResMatchTarget/update_by_id.tf"),
						Check: resource.ComposeTestCheckFunc(
							resource.TestCheckResourceAttr("akamai_appsec_match_target.test", "id", "43253:3008967"),
						),
					},
				},
			})
		})

		client.AssertExpectations(t)
	})

	t.Run("no drift on match target lists sequence mismatch", func(t *testing.T) {
		client := &appsec.Mock{}

		getMatchTargetResponse := appsec.GetMatchTargetResponse{}
		err := json.Unmarshal(testutils.LoadFixtureBytes(t, "testdata/TestResMatchTarget/MatchTargetSequenceChanged.json"), &getMatchTargetResponse)
		require.NoError(t, err)

		createMatchTargetResponse := appsec.CreateMatchTargetResponse{}
		err = json.Unmarshal(testutils.LoadFixtureBytes(t, "testdata/TestResMatchTarget/MatchTargetCreated.json"), &createMatchTargetResponse)
		require.NoError(t, err)

		removeMatchTargetResponse := appsec.RemoveMatchTargetResponse{}
		err = json.Unmarshal(testutils.LoadFixtureBytes(t, "testdata/TestResMatchTarget/MatchTargetCreated.json"), &removeMatchTargetResponse)
		require.NoError(t, err)

		config := appsec.GetConfigurationResponse{}
		err = json.Unmarshal(testutils.LoadFixtureBytes(t, "testdata/TestResConfiguration/LatestConfiguration.json"), &config)
		require.NoError(t, err)

		client.On("GetConfiguration",
			mock.Anything,
			appsec.GetConfigurationRequest{ConfigID: 43253},
		).Return(&config, nil)

		client.On("GetMatchTarget",
			mock.Anything,
			appsec.GetMatchTargetRequest{ConfigID: 43253, ConfigVersion: 7, TargetID: 3008967},
		).Return(&getMatchTargetResponse, nil).Times(2)

		createMatchTargetJSON := testutils.LoadFixtureBytes(t, "testdata/TestResMatchTarget/CreateMatchTarget.json")
		client.On("CreateMatchTarget",
			mock.Anything,
			appsec.CreateMatchTargetRequest{ConfigID: 43253, ConfigVersion: 7, JsonPayloadRaw: createMatchTargetJSON},
		).Return(&createMatchTargetResponse, nil)

		client.On("RemoveMatchTarget",
			mock.Anything,
			appsec.RemoveMatchTargetRequest{ConfigID: 43253, ConfigVersion: 7, TargetID: 3008967},
		).Return(&removeMatchTargetResponse, nil)

		getMatchTargetJSON := `{"type":"website","defaultFile":"NO_MATCH","hostnames":["m.example.com","www.example.net","example.com"],"isNegativePathMatch":false,"filePaths":["/cache/aaabbc*"],"fileExtensions":["carb","pct","pdf","swf","cct","jpeg","js","wmls","hdml","pws"],"securityPolicy":{"policyId":"AAAA_81230"},"targetId":3008967}`

		useClient(client, func() {
			resource.Test(t, resource.TestCase{
				IsUnitTest:               true,
				ProtoV6ProviderFactories: testutils.NewProtoV6ProviderFactory(NewSubprovider()),
				Steps: []resource.TestStep{
					{
						Config: testutils.LoadFixtureString(t, "testdata/TestResMatchTarget/match_by_id.tf"),
						Check: resource.ComposeAggregateTestCheckFunc(
							resource.TestCheckResourceAttr("akamai_appsec_match_target.test", "id", "43253:3008967"),
							resource.TestCheckResourceAttr("akamai_appsec_match_target.test", "match_target", compactJSON(getMatchTargetJSON)),
						),
					},
				},
			})
		})

		client.AssertExpectations(t)
	})

	t.Run("some values should not drift on non changed list or struct values sequence mismatch", func(t *testing.T) {
		client := &appsec.Mock{}

		getMatchTargetResponse := appsec.GetMatchTargetResponse{}
		err := json.Unmarshal(testutils.LoadFixtureBytes(t, "testdata/TestResMatchTarget/MatchTargetSequenceChanged.json"), &getMatchTargetResponse)
		require.NoError(t, err)

		getMatchTargetResponseChanged := appsec.GetMatchTargetResponse{}
		err = json.Unmarshal(testutils.LoadFixtureBytes(t, "testdata/TestResMatchTarget/MatchTargetSequenceOrderChanged.json"), &getMatchTargetResponse)
		require.NoError(t, err)

		createMatchTargetResponse := appsec.CreateMatchTargetResponse{}
		err = json.Unmarshal(testutils.LoadFixtureBytes(t, "testdata/TestResMatchTarget/MatchTargetCreated.json"), &createMatchTargetResponse)
		require.NoError(t, err)

		removeMatchTargetResponse := appsec.RemoveMatchTargetResponse{}
		err = json.Unmarshal(testutils.LoadFixtureBytes(t, "testdata/TestResMatchTarget/MatchTargetCreated.json"), &removeMatchTargetResponse)
		require.NoError(t, err)

		config := appsec.GetConfigurationResponse{}
		err = json.Unmarshal(testutils.LoadFixtureBytes(t, "testdata/TestResConfiguration/LatestConfiguration.json"), &config)
		require.NoError(t, err)

		client.On("GetConfiguration",
			mock.Anything,
			appsec.GetConfigurationRequest{ConfigID: 43253},
		).Return(&config, nil)

		client.On("GetMatchTarget",
			mock.Anything,
			appsec.GetMatchTargetRequest{ConfigID: 43253, ConfigVersion: 7, TargetID: 3008967},
		).Return(&getMatchTargetResponse, nil).Times(1)

		client.On("GetMatchTarget",
			mock.Anything,
			appsec.GetMatchTargetRequest{ConfigID: 43253, ConfigVersion: 7, TargetID: 3008967},
		).Return(&getMatchTargetResponseChanged, nil).Times(1)

		createMatchTargetJSON := testutils.LoadFixtureBytes(t, "testdata/TestResMatchTarget/CreateMatchTarget.json")
		client.On("CreateMatchTarget",
			mock.Anything,
			appsec.CreateMatchTargetRequest{ConfigID: 43253, ConfigVersion: 7, JsonPayloadRaw: createMatchTargetJSON},
		).Return(&createMatchTargetResponse, nil)

		client.On("RemoveMatchTarget",
			mock.Anything,
			appsec.RemoveMatchTargetRequest{ConfigID: 43253, ConfigVersion: 7, TargetID: 3008967},
		).Return(&removeMatchTargetResponse, nil)

		getMatchTargetJSON := `{"type":"website","defaultFile":"NO_MATCH","hostnames":["m.example.com","www.example.net","examplenew.com"],"isNegativePathMatch":false,"filePaths":["/cache/aaabbc*"],"fileExtensions":["carb","pct","pdf","swf","cct","jpeg","js","wmls","hdml","pws"],"securityPolicy":{"policyId":"AAAA_81230"},"targetId":3008967}`

		useClient(client, func() {
			resource.Test(t, resource.TestCase{
				IsUnitTest:               true,
				ProtoV6ProviderFactories: testutils.NewProtoV6ProviderFactory(NewSubprovider()),
				Steps: []resource.TestStep{
					{
						Config: testutils.LoadFixtureString(t, "testdata/TestResMatchTarget/match_by_id.tf"),
						Check: resource.ComposeAggregateTestCheckFunc(
							resource.TestCheckResourceAttr("akamai_appsec_match_target.test", "id", "43253:3008967"),
							resource.TestCheckResourceAttr("akamai_appsec_match_target.test", "match_target", compactJSON(getMatchTargetJSON)),
						),
						ExpectNonEmptyPlan: true,
					},
				},
			})
		})

		client.AssertExpectations(t)
	})

	t.Run("Import match target resource", func(t *testing.T) {
		client := &appsec.Mock{}
		getMatchTargetResponse := appsec.GetMatchTargetResponse{}
		err := json.Unmarshal(testutils.LoadFixtureBytes(t, "testdata/TestResMatchTarget/MatchTargetSequenceChanged.json"), &getMatchTargetResponse)
		require.NoError(t, err)

		createMatchTargetResponse := appsec.CreateMatchTargetResponse{}
		err = json.Unmarshal(testutils.LoadFixtureBytes(t, "testdata/TestResMatchTarget/MatchTargetCreated.json"), &createMatchTargetResponse)
		require.NoError(t, err)

		removeMatchTargetResponse := appsec.RemoveMatchTargetResponse{}
		err = json.Unmarshal(testutils.LoadFixtureBytes(t, "testdata/TestResMatchTarget/MatchTargetCreated.json"), &removeMatchTargetResponse)
		require.NoError(t, err)

		config := appsec.GetConfigurationResponse{}
		err = json.Unmarshal(testutils.LoadFixtureBytes(t, "testdata/TestResConfiguration/LatestConfiguration.json"), &config)
		require.NoError(t, err)

		client.On("GetConfiguration",
			mock.Anything,
			appsec.GetConfigurationRequest{ConfigID: 43253},
		).Return(&config, nil)

		client.On("GetMatchTarget",
			mock.Anything,
			appsec.GetMatchTargetRequest{ConfigID: 43253, ConfigVersion: 7, TargetID: 3008967},
		).Return(&getMatchTargetResponse, nil)

		createMatchTargetJSON := testutils.LoadFixtureBytes(t, "testdata/TestResMatchTarget/CreateMatchTarget.json")
		client.On("CreateMatchTarget",
			mock.Anything,
			appsec.CreateMatchTargetRequest{ConfigID: 43253, ConfigVersion: 7, JsonPayloadRaw: createMatchTargetJSON},
		).Return(&createMatchTargetResponse, nil)

		client.On("RemoveMatchTarget",
			mock.Anything,
			appsec.RemoveMatchTargetRequest{ConfigID: 43253, ConfigVersion: 7, TargetID: 3008967},
		).Return(&removeMatchTargetResponse, nil)

		useClient(client, func() {
			resource.Test(t, resource.TestCase{
				IsUnitTest:               true,
				ProtoV6ProviderFactories: testutils.NewProtoV6ProviderFactory(NewSubprovider()),
				Steps: []resource.TestStep{
					{
						Config: testutils.LoadFixtureString(t, "testdata/TestResMatchTarget/match_by_id.tf"),
					},
					{
						ImportState:       true,
						ImportStateVerify: true,
						ImportStateId:     "43253:3008967",
						ResourceName:      "akamai_appsec_match_target.test",
					},
				},
			})
		})

		client.AssertExpectations(t)
	})

}

func compactJSON(message string) string {
	var dst bytes.Buffer
	err := json.Compact(&dst, []byte(message))
	if err != nil {
		panic(err)
	}
	return dst.String()
}
