package appsec

import (
	"encoding/json"
	"testing"

	"github.com/akamai/AkamaiOPEN-edgegrid-golang/v2/pkg/appsec"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/resource"
	"github.com/stretchr/testify/mock"
)

func TestAccAkamaiReputationProfile_res_basic(t *testing.T) {
	t.Run("match by ReputationProfile ID", func(t *testing.T) {
		client := &mockappsec{}

		cu := appsec.UpdateReputationProfileResponse{}
		expectJSU := compactJSON(loadFixtureBytes("testdata/TestResReputationProfile/ReputationProfileUpdated.json"))
		json.Unmarshal([]byte(expectJSU), &cu)

		cr := appsec.GetReputationProfileResponse{}
		expectJS := compactJSON(loadFixtureBytes("testdata/TestResReputationProfile/ReputationProfiles.json"))
		json.Unmarshal([]byte(expectJS), &cr)

		crp := appsec.CreateReputationProfileResponse{}
		expectJSC := compactJSON(loadFixtureBytes("testdata/TestResReputationProfile/ReputationProfileCreated.json"))
		json.Unmarshal([]byte(expectJSC), &crp)

		crd := appsec.RemoveReputationProfileResponse{}
		expectJSD := compactJSON(loadFixtureBytes("testdata/TestResReputationProfile/ReputationProfileCreated.json"))
		json.Unmarshal([]byte(expectJSD), &crd)

		client.On("GetReputationProfile",
			mock.Anything, // ctx is irrelevant for this test
			appsec.GetReputationProfileRequest{ConfigID: 43253, ConfigVersion: 7, ReputationProfileId: 12345},
		).Return(&cr, nil)

		client.On("RemoveReputationProfile",
			mock.Anything, // ctx is irrelevant for this test
			appsec.RemoveReputationProfileRequest{ConfigID: 43253, ConfigVersion: 7, ReputationProfileId: 12345},
		).Return(&crd, nil)

		client.On("CreateReputationProfile",
			mock.Anything, // ctx is irrelevant for this test
			appsec.CreateReputationProfileRequest{ConfigID: 43253, ConfigVersion: 7, Name: "Web Attack Rep Profile", Description: "Reputation profile description", Context: "WEBATCK", Threshold: 5, SharedIPHandling: "NON_SHARED", Condition: struct {
				PositiveMatch    bool "json:\"positiveMatch\""
				AtomicConditions []struct {
					PositiveMatch bool     "json:\"positiveMatch\""
					ClassName     string   "json:\"className\""
					Value         []string "json:\"value,omitempty\""
					NameWildcard  bool     "json:\"nameWildcard,omitempty\""
					ValueWildcard bool     "json:\"valueWildcard,omitempty\""
					NameCase      bool     "json:\"nameCase,omitempty\""
					Name          string   "json:\"name,omitempty\""
					Host          []string "json:\"host,omitempty\""
				} "json:\"atomicConditions\""
			}{PositiveMatch: true, AtomicConditions: []struct {
				PositiveMatch bool     "json:\"positiveMatch\""
				ClassName     string   "json:\"className\""
				Value         []string "json:\"value,omitempty\""
				NameWildcard  bool     "json:\"nameWildcard,omitempty\""
				ValueWildcard bool     "json:\"valueWildcard,omitempty\""
				NameCase      bool     "json:\"nameCase,omitempty\""
				Name          string   "json:\"name,omitempty\""
				Host          []string "json:\"host,omitempty\""
			}{struct {
				PositiveMatch bool     "json:\"positiveMatch\""
				ClassName     string   "json:\"className\""
				Value         []string "json:\"value,omitempty\""
				NameWildcard  bool     "json:\"nameWildcard,omitempty\""
				ValueWildcard bool     "json:\"valueWildcard,omitempty\""
				NameCase      bool     "json:\"nameCase,omitempty\""
				Name          string   "json:\"name,omitempty\""
				Host          []string "json:\"host,omitempty\""
			}{PositiveMatch: true, ClassName: "AsNumberCondition", Value: []string{"1"}, NameWildcard: false, ValueWildcard: false, NameCase: false, Name: "", Host: []string(nil)}, struct {
				PositiveMatch bool     "json:\"positiveMatch\""
				ClassName     string   "json:\"className\""
				Value         []string "json:\"value,omitempty\""
				NameWildcard  bool     "json:\"nameWildcard,omitempty\""
				ValueWildcard bool     "json:\"valueWildcard,omitempty\""
				NameCase      bool     "json:\"nameCase,omitempty\""
				Name          string   "json:\"name,omitempty\""
				Host          []string "json:\"host,omitempty\""
			}{PositiveMatch: true, ClassName: "RequestCookieCondition", Value: []string(nil), NameWildcard: true, ValueWildcard: true, NameCase: true, Name: "x-header", Host: []string(nil)}, struct {
				PositiveMatch bool     "json:\"positiveMatch\""
				ClassName     string   "json:\"className\""
				Value         []string "json:\"value,omitempty\""
				NameWildcard  bool     "json:\"nameWildcard,omitempty\""
				ValueWildcard bool     "json:\"valueWildcard,omitempty\""
				NameCase      bool     "json:\"nameCase,omitempty\""
				Name          string   "json:\"name,omitempty\""
				Host          []string "json:\"host,omitempty\""
			}{PositiveMatch: true, ClassName: "HostCondition", Value: []string(nil), NameWildcard: false, ValueWildcard: true, NameCase: false, Name: "", Host: []string{"*.com"}}}}},
		).Return(&crp, nil)

		client.On("UpdateReputationProfile",
			mock.Anything, // ctx is irrelevant for this test
			appsec.UpdateReputationProfileRequest{ConfigID: 43253, ConfigVersion: 7, ReputationProfileId: 12345},
		).Return(&cu, nil)

		useClient(client, func() {
			resource.Test(t, resource.TestCase{
				IsUnitTest: false,
				Providers:  testAccProviders,
				Steps: []resource.TestStep{
					{
						Config: loadFixtureString("testdata/TestResReputationProfile/match_by_id.tf"),
						Check: resource.ComposeAggregateTestCheckFunc(
							resource.TestCheckResourceAttr("akamai_appsec_reputation_profile.test", "id", "12345"),
						),
					},
					/*	{
						Config: loadFixtureString("testdata/TestResReputationProfile/update_by_id.tf"),
						Check: resource.ComposeAggregateTestCheckFunc(
							resource.TestCheckResourceAttr("akamai_appsec_reputation_profile.test", "id", "12345"),
						),
					},*/
				},
			})
		})

		client.AssertExpectations(t)
	})

}
