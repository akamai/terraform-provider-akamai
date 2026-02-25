package appsec

import (
	"encoding/json"
	"testing"

	"github.com/akamai/AkamaiOPEN-edgegrid-golang/v13/pkg/appsec"
	"github.com/akamai/terraform-provider-akamai/v10/pkg/common/testutils"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/schema"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

func TestUkraineGeoControlActionEqual(t *testing.T) {

	tests := map[string]struct {
		ukraineGeoControlAction interface{}
		expected                bool
	}{
		"action set": {
			ukraineGeoControlAction: "alert",
			expected:                false,
		},
		"action not set": {
			expected: true,
		},
	}

	resourceSchema := map[string]*schema.Schema{
		"ukraine_geo_control_action": {
			Type: schema.TypeString,
		},
	}

	for name, test := range tests {
		t.Run(name, func(t *testing.T) {
			resourceDataMap := map[string]interface{}{
				"ukraine_geo_control_action": test.ukraineGeoControlAction,
			}
			resourceData := schema.TestResourceDataRaw(t, resourceSchema, resourceDataMap)

			res := suppressDiffUkraineGeoControlAction("", "", "", resourceData)
			assert.Equal(t, test.expected, res)
		})
	}
}

func TestAreReputationProfilesEqual(t *testing.T) {
	deepCopyProfile := func(profile appsec.CreateReputationProfileResponse) appsec.CreateReputationProfileResponse {
		b, _ := json.Marshal(profile)
		var profileCopy appsec.CreateReputationProfileResponse
		if err := json.Unmarshal(b, &profileCopy); err != nil {
			panic(err)
		}
		return profileCopy
	}

	baseProfile := appsec.CreateReputationProfileResponse{}
	err := json.Unmarshal(testutils.LoadFixtureBytes(t, "testdata/TestDiffSuppressFuncs/ReputationProfile.json"), &baseProfile)
	require.NoError(t, err)

	t.Run("exact match", func(t *testing.T) {
		clone := deepCopyProfile(baseProfile)
		require.True(t, areReputationProfilesEqual(baseProfile, clone))
	})

	t.Run("different CheckIps", func(t *testing.T) {
		mod := deepCopyProfile(baseProfile)
		mod.Condition.AtomicConditions[0].CheckIps = "different"
		require.False(t, areReputationProfilesEqual(baseProfile, mod))
	})

	t.Run("missing AtomicConditions", func(t *testing.T) {
		mod := deepCopyProfile(baseProfile)
		mod.Condition.AtomicConditions = mod.Condition.AtomicConditions[:2]
		require.False(t, areReputationProfilesEqual(baseProfile, mod))
	})

	t.Run("different NameCase", func(t *testing.T) {
		mod := deepCopyProfile(baseProfile)
		mod.Condition.AtomicConditions[1].NameCase = false
		require.True(t, areReputationProfilesEqual(baseProfile, mod))
	})

	t.Run("different ValueWildcard", func(t *testing.T) {
		mod := deepCopyProfile(baseProfile)
		mod.Condition.AtomicConditions[1].ValueWildcard = false
		require.True(t, areReputationProfilesEqual(baseProfile, mod))
	})

	t.Run("different PositiveMatch", func(t *testing.T) {
		mod := deepCopyProfile(baseProfile)
		mod.Condition.AtomicConditions[2].PositiveMatch = false
		require.True(t, areReputationProfilesEqual(baseProfile, mod))
	})

	t.Run("different Value", func(t *testing.T) {
		mod := deepCopyProfile(baseProfile)
		mod.Condition.AtomicConditions[0].Value = []string{"2"}
		require.False(t, areReputationProfilesEqual(baseProfile, mod))
	})

	t.Run("different Host", func(t *testing.T) {
		mod := deepCopyProfile(baseProfile)
		mod.Condition.AtomicConditions[2].Host = []string{"*.org"}
		require.False(t, areReputationProfilesEqual(baseProfile, mod))
	})

	t.Run("different ValueCase", func(t *testing.T) {
		mod := deepCopyProfile(baseProfile)
		mod.Condition.AtomicConditions[1].ValueCase = false
		base := deepCopyProfile(mod)

		mod2 := deepCopyProfile(base)
		mod2.Condition.AtomicConditions[1].ValueCase = true
		require.False(t, areReputationProfilesEqual(base, mod2))
	})

	t.Run("different ValueWildcard", func(t *testing.T) {
		base := deepCopyProfile(baseProfile)
		base.Condition.AtomicConditions[1].ValueWildcard = false

		mod := deepCopyProfile(base)
		mod.Condition.AtomicConditions[1].ValueWildcard = true
		require.False(t, areReputationProfilesEqual(base, mod))
	})

	t.Run("different NameCase", func(t *testing.T) {
		base := deepCopyProfile(baseProfile)
		base.Condition.AtomicConditions[1].NameCase = false

		mod := deepCopyProfile(base)
		mod.Condition.AtomicConditions[1].NameCase = true
		require.False(t, areReputationProfilesEqual(base, mod))
	})

	t.Run("different NameWildcard", func(t *testing.T) {
		base := deepCopyProfile(baseProfile)
		base.Condition.AtomicConditions[1].NameWildcard = false

		mod := deepCopyProfile(base)
		mod.Condition.AtomicConditions[1].NameWildcard = true
		require.False(t, areReputationProfilesEqual(base, mod))
	})

	t.Run("different PositiveMatch", func(t *testing.T) {
		base := deepCopyProfile(baseProfile)
		base.Condition.AtomicConditions[0].PositiveMatch = false

		mod := deepCopyProfile(base)
		mod.Condition.AtomicConditions[0].PositiveMatch = true
		require.False(t, areReputationProfilesEqual(base, mod))
	})

	t.Run("allowed PositiveMatch with HostCondition", func(t *testing.T) {
		base := deepCopyProfile(baseProfile)
		base.Condition.AtomicConditions[2].ClassName = "HostCondition"
		base.Condition.AtomicConditions[2].PositiveMatch = false

		mod := deepCopyProfile(base)
		mod.Condition.AtomicConditions[2].PositiveMatch = true
		require.True(t, areReputationProfilesEqual(base, mod))
	})
}
