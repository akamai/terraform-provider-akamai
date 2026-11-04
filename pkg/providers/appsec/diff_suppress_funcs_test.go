package appsec

import (
	"encoding/json"
	"testing"

	"github.com/akamai/AkamaiOPEN-edgegrid-golang/v15/pkg/appsec"
	"github.com/akamai/terraform-provider-akamai/v12/pkg/common/testutils"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/schema"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

func TestNormalizeJSON(t *testing.T) {
	t.Parallel()
	tests := map[string]struct {
		input   string
		want    string
		wantErr bool
	}{
		"keys sorted alphabetically": {
			input: `{"z":1,"a":2}`,
			want:  `{"a":2,"z":1}`,
		},
		"nested object keys sorted": {
			input: `{"z":{"d":4,"c":3},"a":1}`,
			want:  `{"a":1,"z":{"c":3,"d":4}}`,
		},
		"whitespace removed": {
			input: `{ "a" : 1 , "b" : 2 }`,
			want:  `{"a":1,"b":2}`,
		},
		"array order preserved": {
			input: `{"items":[3,1,2]}`,
			want:  `{"items":[3,1,2]}`,
		},
		"HTML characters not escaped": {
			input: `{"path":"/api/<v>/data?a=1&b=2>0"}`,
			want:  `{"path":"/api/<v>/data?a=1&b=2>0"}`,
		},
		"large integer preserved without float64 truncation": {
			input: `{"id":9007199254740993}`,
			want:  `{"id":9007199254740993}`,
		},
		"invalid JSON returns error": {
			input:   `not-json`,
			wantErr: true,
		},
	}

	for name, tc := range tests {
		t.Run(name, func(t *testing.T) {
			got, err := normalizeJSON(tc.input)
			if tc.wantErr {
				require.Error(t, err)
				return
			}
			require.NoError(t, err)
			assert.Equal(t, tc.want, got)
		})
	}
}

func TestMarshalNormalize(t *testing.T) {
	t.Parallel()

	t.Run("struct fields sorted alphabetically", func(t *testing.T) {
		t.Parallel()
		type payload struct {
			Z int `json:"z"`
			A int `json:"a"`
		}
		got, err := marshalNormalize(payload{Z: 1, A: 2})
		require.NoError(t, err)
		assert.Equal(t, `{"a":2,"z":1}`, got)
	})

	t.Run("marshal error propagated", func(t *testing.T) {
		t.Parallel()
		_, err := marshalNormalize(make(chan int))
		require.Error(t, err)
	})
}

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

func TestSuppressRatePolicyDiffs(t *testing.T) {
	t.Parallel()

	basePolicy := `{
		"id": 134644,
		"configId": 43253,
		"configVersion": 7,
		"matchType": "path",
		"type": "WAF",
		"name": "Test Policy",
		"averageThreshold": 5,
		"burstThreshold": 10,
		"clientIdentifiers": ["ip"],
		"useXForwardForHeaders": true,
		"requestType": "ClientRequest",
		"sameActionOnIpv6": false,
		"pathMatchType": "Custom",
		"pathUriPositiveMatch": true,
		"counterType": "per_edge",
		"penaltyBoxDuration": "THIRTY_MINUTES"
	}`

	tests := map[string]struct {
		old      string
		new      string
		expected bool
	}{
		"identical policies": {
			old:      basePolicy,
			new:      basePolicy,
			expected: true,
		},
		"config omits id — server-assigned, treat as matching state": {
			old: basePolicy,
			new: `{
				"matchType": "path", "type": "WAF", "name": "Test Policy",
				"averageThreshold": 5, "burstThreshold": 10,
				"clientIdentifiers": ["ip"], "useXForwardForHeaders": true,
				"requestType": "ClientRequest", "sameActionOnIpv6": false,
				"pathMatchType": "Custom", "pathUriPositiveMatch": true,
				"counterType": "per_edge", "penaltyBoxDuration": "THIRTY_MINUTES"
			}`,
			expected: true,
		},
		"config omits configId and configVersion — server-assigned, treat as matching state": {
			old: basePolicy,
			new: `{
				"id": 134644,
				"matchType": "path", "type": "WAF", "name": "Test Policy",
				"averageThreshold": 5, "burstThreshold": 10,
				"clientIdentifiers": ["ip"], "useXForwardForHeaders": true,
				"requestType": "ClientRequest", "sameActionOnIpv6": false,
				"pathMatchType": "Custom", "pathUriPositiveMatch": true,
				"counterType": "per_edge", "penaltyBoxDuration": "THIRTY_MINUTES"
			}`,
			expected: true,
		},
		"config omits counterType — defaults to per_edge": {
			old: basePolicy,
			new: `{
				"matchType": "path", "type": "WAF", "name": "Test Policy",
				"averageThreshold": 5, "burstThreshold": 10,
				"clientIdentifiers": ["ip"], "useXForwardForHeaders": true,
				"requestType": "ClientRequest", "sameActionOnIpv6": false,
				"pathMatchType": "Custom", "pathUriPositiveMatch": true,
				"penaltyBoxDuration": "THIRTY_MINUTES"
			}`,
			expected: true,
		},
		"config omits penaltyBoxDuration — defaults to TEN_MINUTES": {
			old: `{
				"matchType": "path", "type": "WAF", "name": "Test Policy",
				"penaltyBoxDuration": "TEN_MINUTES", "counterType": "per_edge"
			}`,
			new: `{
				"matchType": "path", "type": "WAF", "name": "Test Policy",
				"counterType": "per_edge"
			}`,
			expected: true,
		},
		"different matchType — real change, do not suppress": {
			old: basePolicy,
			new: `{
				"matchType": "hostname", "type": "WAF", "name": "Test Policy",
				"counterType": "per_edge", "penaltyBoxDuration": "THIRTY_MINUTES"
			}`,
			expected: false,
		},
		"different counterType explicitly set — real change, do not suppress": {
			old: basePolicy,
			new: `{
				"matchType": "path", "type": "WAF", "name": "Test Policy",
				"counterType": "region_aggregated", "penaltyBoxDuration": "THIRTY_MINUTES"
			}`,
			expected: false,
		},
		"invalid old JSON": {
			old:      `not-json`,
			new:      basePolicy,
			expected: false,
		},
		"invalid new JSON": {
			old:      basePolicy,
			new:      `not-json`,
			expected: false,
		},
	}

	for name, test := range tests {
		t.Run(name, func(t *testing.T) {
			t.Parallel()
			result := suppressRatePolicyDiffs("", test.old, test.new, nil)
			assert.Equal(t, test.expected, result, "suppressRatePolicyDiffs(%q, %q)", test.old, test.new)
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
