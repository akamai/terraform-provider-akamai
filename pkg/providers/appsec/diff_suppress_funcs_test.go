package appsec

import (
	"encoding/json"
	"testing"

	"github.com/akamai/AkamaiOPEN-edgegrid-golang/v13/pkg/appsec"
	"github.com/akamai/terraform-provider-akamai/v11/pkg/common/testutils"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/schema"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

func TestUkraineGeoControlActionEqual(t *testing.T) {
	t.Parallel()

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
			t.Parallel()
			resourceDataMap := map[string]interface{}{
				"ukraine_geo_control_action": test.ukraineGeoControlAction,
			}
			resourceData := schema.TestResourceDataRaw(t, resourceSchema, resourceDataMap)

			res := suppressDiffUkraineGeoControlAction("", "", "", resourceData)
			assert.Equal(t, test.expected, res)
		})
	}
}

func TestSuppressFieldForContractID(t *testing.T) {
	resourceSchema := map[string]*schema.Schema{
		"contract_id": {Type: schema.TypeString},
	}

	tests := map[string]struct {
		oldValue string
		newValue string
		id       string
		expected bool
	}{
		"post-import: old empty, resource exists": {
			oldValue: "",
			newValue: "ctr_12345",
			id:       "some-id",
			expected: true,
		},
		"new resource: old empty, no resource id": {
			oldValue: "",
			newValue: "ctr_12345",
			id:       "",
			expected: false,
		},
		"values match": {
			oldValue: "ctr_12345",
			newValue: "ctr_12345",
			id:       "some-id",
			expected: true,
		},
		"values differ": {
			oldValue: "ctr_12345",
			newValue: "ctr_99999",
			id:       "some-id",
			expected: false,
		},
	}

	for name, test := range tests {
		t.Run(name, func(t *testing.T) {
			d := schema.TestResourceDataRaw(t, resourceSchema, map[string]interface{}{
				"contract_id": test.oldValue,
			})
			d.SetId(test.id)
			assert.Equal(t, test.expected, suppressFieldForContractID("", test.oldValue, test.newValue, d))
		})
	}
}

func TestSuppressFieldForPrefixedGroupIDPostImport(t *testing.T) {
	resourceSchema := map[string]*schema.Schema{
		"group_id": {Type: schema.TypeString},
	}

	tests := map[string]struct {
		oldValue string
		newValue string
		id       string
		expected bool
	}{
		"post-import: old empty, resource exists": {
			oldValue: "",
			newValue: "grp_12345",
			id:       "some-id",
			expected: true,
		},
		"new resource: old empty, no resource id": {
			oldValue: "",
			newValue: "grp_12345",
			id:       "",
			expected: false,
		},
		"values match": {
			oldValue: "grp_12345",
			newValue: "grp_12345",
			id:       "some-id",
			expected: true,
		},
		"values differ": {
			oldValue: "grp_12345",
			newValue: "grp_99999",
			id:       "some-id",
			expected: false,
		},
	}

	for name, test := range tests {
		t.Run(name, func(t *testing.T) {
			d := schema.TestResourceDataRaw(t, resourceSchema, map[string]interface{}{
				"group_id": test.oldValue,
			})
			d.SetId(test.id)
			assert.Equal(t, test.expected, suppressFieldForPrefixedGroupID("", test.oldValue, test.newValue, d))
		})
	}
}

func TestSuppressJSONDiffsIgnoringArrayOrder(t *testing.T) {
	tests := map[string]struct {
		old      string
		new      string
		expected bool
	}{
		"identical flat objects": {
			old:      `{"a":1,"b":2}`,
			new:      `{"a":1,"b":2}`,
			expected: true,
		},
		"flat arrays same order": {
			old:      `[1,2,3]`,
			new:      `[1,2,3]`,
			expected: true,
		},
		"flat arrays different order": {
			old:      `[1,2,3]`,
			new:      `[3,1,2]`,
			expected: true,
		},
		"object with array field reordered": {
			old:      `{"tags":["b","a","c"]}`,
			new:      `{"tags":["a","b","c"]}`,
			expected: true,
		},
		"deeply nested arrays reordered": {
			old:      `{"x":{"items":[{"id":2},{"id":1}]}}`,
			new:      `{"x":{"items":[{"id":1},{"id":2}]}}`,
			expected: true,
		},
		"objects with different values": {
			old:      `{"a":1}`,
			new:      `{"a":2}`,
			expected: false,
		},
		"objects with extra key": {
			old:      `{"a":1}`,
			new:      `{"a":1,"b":2}`,
			expected: false,
		},
		"arrays with different values": {
			old:      `[1,2,3]`,
			new:      `[1,2,4]`,
			expected: false,
		},
		"arrays with different lengths": {
			old:      `[1,2]`,
			new:      `[1,2,3]`,
			expected: false,
		},
		"invalid old JSON": {
			old:      `not-json`,
			new:      `{"a":1}`,
			expected: false,
		},
		"invalid new JSON": {
			old:      `{"a":1}`,
			new:      `not-json`,
			expected: false,
		},
	}

	for name, test := range tests {
		t.Run(name, func(t *testing.T) {
			result := suppressJSONDiffsIgnoringArrayOrder("", test.old, test.new, nil)
			assert.Equal(t, test.expected, result)
		})
	}
}

func TestAreReputationProfilesEqual(t *testing.T) {
	t.Parallel()
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
		t.Parallel()
		clone := deepCopyProfile(baseProfile)
		require.True(t, areReputationProfilesEqual(baseProfile, clone))
	})

	t.Run("different CheckIps", func(t *testing.T) {
		t.Parallel()
		mod := deepCopyProfile(baseProfile)
		mod.Condition.AtomicConditions[0].CheckIps = "different"
		require.False(t, areReputationProfilesEqual(baseProfile, mod))
	})

	t.Run("missing AtomicConditions", func(t *testing.T) {
		t.Parallel()
		mod := deepCopyProfile(baseProfile)
		mod.Condition.AtomicConditions = mod.Condition.AtomicConditions[:2]
		require.False(t, areReputationProfilesEqual(baseProfile, mod))
	})

	t.Run("different NameCase", func(t *testing.T) {
		t.Parallel()
		mod := deepCopyProfile(baseProfile)
		mod.Condition.AtomicConditions[1].NameCase = false
		require.True(t, areReputationProfilesEqual(baseProfile, mod))
	})

	t.Run("different ValueWildcard", func(t *testing.T) {
		t.Parallel()
		mod := deepCopyProfile(baseProfile)
		mod.Condition.AtomicConditions[1].ValueWildcard = false
		require.True(t, areReputationProfilesEqual(baseProfile, mod))
	})

	t.Run("different PositiveMatch", func(t *testing.T) {
		t.Parallel()
		mod := deepCopyProfile(baseProfile)
		mod.Condition.AtomicConditions[2].PositiveMatch = false
		require.True(t, areReputationProfilesEqual(baseProfile, mod))
	})

	t.Run("different Value", func(t *testing.T) {
		t.Parallel()
		mod := deepCopyProfile(baseProfile)
		mod.Condition.AtomicConditions[0].Value = []string{"2"}
		require.False(t, areReputationProfilesEqual(baseProfile, mod))
	})

	t.Run("different Host", func(t *testing.T) {
		t.Parallel()
		mod := deepCopyProfile(baseProfile)
		mod.Condition.AtomicConditions[2].Host = []string{"*.org"}
		require.False(t, areReputationProfilesEqual(baseProfile, mod))
	})

	t.Run("different ValueCase", func(t *testing.T) {
		t.Parallel()
		mod := deepCopyProfile(baseProfile)
		mod.Condition.AtomicConditions[1].ValueCase = false
		base := deepCopyProfile(mod)

		mod2 := deepCopyProfile(base)
		mod2.Condition.AtomicConditions[1].ValueCase = true
		require.False(t, areReputationProfilesEqual(base, mod2))
	})

	t.Run("different ValueWildcard", func(t *testing.T) {
		t.Parallel()
		base := deepCopyProfile(baseProfile)
		base.Condition.AtomicConditions[1].ValueWildcard = false

		mod := deepCopyProfile(base)
		mod.Condition.AtomicConditions[1].ValueWildcard = true
		require.False(t, areReputationProfilesEqual(base, mod))
	})

	t.Run("different NameCase", func(t *testing.T) {
		t.Parallel()
		base := deepCopyProfile(baseProfile)
		base.Condition.AtomicConditions[1].NameCase = false

		mod := deepCopyProfile(base)
		mod.Condition.AtomicConditions[1].NameCase = true
		require.False(t, areReputationProfilesEqual(base, mod))
	})

	t.Run("different NameWildcard", func(t *testing.T) {
		t.Parallel()
		base := deepCopyProfile(baseProfile)
		base.Condition.AtomicConditions[1].NameWildcard = false

		mod := deepCopyProfile(base)
		mod.Condition.AtomicConditions[1].NameWildcard = true
		require.False(t, areReputationProfilesEqual(base, mod))
	})

	t.Run("different PositiveMatch", func(t *testing.T) {
		t.Parallel()
		base := deepCopyProfile(baseProfile)
		base.Condition.AtomicConditions[0].PositiveMatch = false

		mod := deepCopyProfile(base)
		mod.Condition.AtomicConditions[0].PositiveMatch = true
		require.False(t, areReputationProfilesEqual(base, mod))
	})

	t.Run("allowed PositiveMatch with HostCondition", func(t *testing.T) {
		t.Parallel()
		base := deepCopyProfile(baseProfile)
		base.Condition.AtomicConditions[2].ClassName = "HostCondition"
		base.Condition.AtomicConditions[2].PositiveMatch = false

		mod := deepCopyProfile(base)
		mod.Condition.AtomicConditions[2].PositiveMatch = true
		require.True(t, areReputationProfilesEqual(base, mod))
	})
}
