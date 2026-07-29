package datastream

import (
	"testing"

	"github.com/akamai/AkamaiOPEN-edgegrid-golang/v13/pkg/datastream"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/schema"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

func TestNormalizeContractID(t *testing.T) {
	t.Parallel()

	assert.Equal(t, "1-ABC", normalizeContractID("ctr_1-ABC"))
	assert.Equal(t, "1-ABC", normalizeContractID("  1-ABC  "))
	assert.Equal(t, "", normalizeContractID("   "))
	assert.Equal(t, "", normalizeContractID("ctr_"))
}

func TestParseGroupIDString(t *testing.T) {
	t.Parallel()

	groupID, err := parseGroupIDString("grp_117988")
	require.NoError(t, err)
	assert.Equal(t, 117988, groupID)

	groupID, err = parseGroupIDString("  ")
	require.NoError(t, err)
	assert.Equal(t, 0, groupID)

	_, err = parseGroupIDString("not-a-number")
	require.Error(t, err)
}

func TestGetContractIDForStream(t *testing.T) {
	t.Parallel()

	t.Run("CDN allows omitted contract_id", func(t *testing.T) {
		d := schema.TestResourceDataRaw(t, datastreamResourceSchema, map[string]interface{}{
			"log_type": "CDN",
			"active":   false,
		})

		contractID, err := getContractIDForStream(d, datastream.LogTypeCDN)
		require.NoError(t, err)
		assert.Empty(t, contractID)
	})

	t.Run("APPSEC requires contract_id", func(t *testing.T) {
		d := schema.TestResourceDataRaw(t, datastreamResourceSchema, map[string]interface{}{
			"log_type": "APPSEC",
			"active":   false,
		})

		_, err := getContractIDForStream(d, datastream.LogTypeAppSec)
		require.Error(t, err)
		assert.Equal(t, "`contract_id` is required for log_type \"APPSEC\"", err.Error())
	})

	t.Run("APPSEC rejects whitespace-only contract_id", func(t *testing.T) {
		d := schema.TestResourceDataRaw(t, datastreamResourceSchema, map[string]interface{}{
			"log_type":    "APPSEC",
			"active":      false,
			"contract_id": "   ",
			"group_id":    "42",
		})

		_, err := getContractIDForStream(d, datastream.LogTypeAppSec)
		require.Error(t, err)
		assert.Equal(t, "`contract_id` is required for log_type \"APPSEC\"", err.Error())
	})
}

func TestGetGroupIDForStream(t *testing.T) {
	t.Parallel()

	t.Run("CDN allows omitted group_id", func(t *testing.T) {
		d := schema.TestResourceDataRaw(t, datastreamResourceSchema, map[string]interface{}{
			"log_type": "CDN",
			"active":   false,
		})

		groupID, err := getGroupIDForStream(d, datastream.LogTypeCDN)
		require.NoError(t, err)
		assert.Zero(t, groupID)
	})

	t.Run("APPSEC requires group_id", func(t *testing.T) {
		d := schema.TestResourceDataRaw(t, datastreamResourceSchema, map[string]interface{}{
			"log_type":    "APPSEC",
			"active":      false,
			"contract_id": "test_contract",
		})

		_, err := getGroupIDForStream(d, datastream.LogTypeAppSec)
		require.Error(t, err)
		assert.Equal(t, "`group_id` is required for log_type \"APPSEC\"", err.Error())
	})

	t.Run("APPSEC rejects group_id zero", func(t *testing.T) {
		d := schema.TestResourceDataRaw(t, datastreamResourceSchema, map[string]interface{}{
			"log_type":    "APPSEC",
			"active":      false,
			"contract_id": "test_contract",
			"group_id":    "0",
		})

		_, err := getGroupIDForStream(d, datastream.LogTypeAppSec)
		require.Error(t, err)
		assert.Equal(t, "`group_id` must be at least 1 for log_type \"APPSEC\"", err.Error())
	})
}

func TestValidateNonCDNContractAndGroupIDValues(t *testing.T) {
	t.Parallel()

	t.Run("APPSEC requires contract_id", func(t *testing.T) {
		err := validateNonCDNContractAndGroupIDValues("", "42", datastream.LogTypeAppSec)
		require.Error(t, err)
		assert.Equal(t, "`contract_id` is required for log_type \"APPSEC\"", err.Error())
	})

	t.Run("APPSEC requires group_id at least 1", func(t *testing.T) {
		err := validateNonCDNContractAndGroupIDValues("test_contract", "0", datastream.LogTypeAppSec)
		require.Error(t, err)
		assert.Equal(t, "`group_id` must be at least 1 for log_type \"APPSEC\"", err.Error())
	})

	t.Run("APPSEC rejects whitespace-only group_id", func(t *testing.T) {
		err := validateNonCDNContractAndGroupIDValues("test_contract", "   ", datastream.LogTypeAppSec)
		require.Error(t, err)
		assert.Equal(t, "`group_id` is required for log_type \"APPSEC\"", err.Error())
	})

	t.Run("APPSEC rejects invalid group_id", func(t *testing.T) {
		err := validateNonCDNContractAndGroupIDValues("test_contract", "abc", datastream.LogTypeAppSec)
		require.Error(t, err)
		assert.Equal(t, "invalid `group_id` \"abc\": strconv.Atoi: parsing \"abc\": invalid syntax", err.Error())
	})

	t.Run("APPSEC accepts populated contract_id and group_id", func(t *testing.T) {
		err := validateNonCDNContractAndGroupIDValues("test_contract", "grp_42", datastream.LogTypeAppSec)
		require.NoError(t, err)
	})
}

func TestValidateCDNGroupIDWhenSet(t *testing.T) {
	t.Parallel()

	t.Run("omitted group_id is allowed", func(t *testing.T) {
		require.NoError(t, validateCDNGroupIDWhenSet(""))
	})

	t.Run("explicit zero group_id is rejected", func(t *testing.T) {
		err := validateCDNGroupIDWhenSet("0")
		require.Error(t, err)
		assert.Equal(t, "`group_id` must be at least 1 for log_type \"CDN\"", err.Error())
	})

	t.Run("invalid group_id is rejected", func(t *testing.T) {
		err := validateCDNGroupIDWhenSet("abc")
		require.Error(t, err)
		assert.Equal(t, "invalid `group_id` \"abc\": strconv.Atoi: parsing \"abc\": invalid syntax", err.Error())
	})

	t.Run("negative group_id is rejected", func(t *testing.T) {
		err := validateCDNGroupIDWhenSet("-1")
		require.Error(t, err)
		assert.Contains(t, err.Error(), "`group_id` must be at least 1")
	})

	t.Run("populated group_id is allowed", func(t *testing.T) {
		require.NoError(t, validateCDNGroupIDWhenSet("grp_42"))
	})
}

func TestIsStreamNameUnsetForDestroy(t *testing.T) {
	t.Parallel()

	assert.True(t, isStreamNameUnsetForDestroy(nil, false))
	assert.True(t, isStreamNameUnsetForDestroy(nil, true))
	assert.True(t, isStreamNameUnsetForDestroy("", true))
	assert.True(t, isStreamNameUnsetForDestroy("   ", true))
	assert.False(t, isStreamNameUnsetForDestroy("live-stream", true))
}

func TestGetGroupIDForStream_CDNNegativeGroupID(t *testing.T) {
	t.Parallel()

	d := schema.TestResourceDataRaw(t, datastreamResourceSchema, map[string]interface{}{
		"log_type": "CDN",
		"active":   false,
		"group_id": "-1",
	})

	_, err := getGroupIDForStream(d, datastream.LogTypeCDN)
	require.Error(t, err)
	assert.Contains(t, err.Error(), "`group_id` must be at least 1")
}

func TestGetGroupIDForStream_CDNExplicitZeroGroupID(t *testing.T) {
	t.Parallel()

	d := schema.TestResourceDataRaw(t, datastreamResourceSchema, map[string]interface{}{
		"log_type": "CDN",
		"active":   false,
		"group_id": "0",
	})

	_, err := getGroupIDForStream(d, datastream.LogTypeCDN)
	require.Error(t, err)
	assert.Contains(t, err.Error(), "`group_id` must be at least 1")
}

func TestResolveContractIDForRead(t *testing.T) {
	t.Parallel()

	t.Run("prefers API value when present", func(t *testing.T) {
		d := schema.TestResourceDataRaw(t, datastreamResourceSchema, map[string]interface{}{
			"contract_id": "configured_contract",
		})

		assert.Equal(t, "api_contract", resolveContractIDForRead("ctr_api_contract", d))
	})

	t.Run("uses API value when config omits contract_id", func(t *testing.T) {
		d := schema.TestResourceDataRaw(t, datastreamResourceSchema, map[string]interface{}{
			"log_type": "CDN",
		})

		assert.Equal(t, "api_contract", resolveContractIDForRead("ctr_api_contract", d))
	})

	t.Run("preserves configured value when API omits contractId", func(t *testing.T) {
		d := schema.TestResourceDataRaw(t, datastreamResourceSchema, map[string]interface{}{
			"contract_id": "ctr_configured_contract",
		})

		assert.Equal(t, "configured_contract", resolveContractIDForRead("", d))
	})

	t.Run("preserves ctr_ configured value when API returns empty contractId (I#775)", func(t *testing.T) {
		d := schema.TestResourceDataRaw(t, datastreamResourceSchema, map[string]interface{}{
			"contract_id": "ctr_G-XXXXX",
		})

		assert.Equal(t, "G-XXXXX", resolveContractIDForRead("", d))
	})

	t.Run("treats whitespace-only API contractId as omitted (I#775)", func(t *testing.T) {
		d := schema.TestResourceDataRaw(t, datastreamResourceSchema, map[string]interface{}{
			"contract_id": "ctr_configured_contract",
		})

		assert.Equal(t, "configured_contract", resolveContractIDForRead("   ", d))
	})

	t.Run("returns empty when API and config omit contract_id", func(t *testing.T) {
		d := schema.TestResourceDataRaw(t, datastreamResourceSchema, map[string]interface{}{
			"log_type": "CDN",
		})

		assert.Empty(t, resolveContractIDForRead("", d))
	})
}

func TestResolveGroupIDForRead(t *testing.T) {
	t.Parallel()

	t.Run("prefers API value when present", func(t *testing.T) {
		d := schema.TestResourceDataRaw(t, datastreamResourceSchema, map[string]interface{}{
			"group_id": "9999",
		})

		assert.Equal(t, "1337", resolveGroupIDForRead(1337, d))
	})

	t.Run("uses API value when config omits group_id", func(t *testing.T) {
		d := schema.TestResourceDataRaw(t, datastreamResourceSchema, map[string]interface{}{
			"log_type": "CDN",
		})

		assert.Equal(t, "1337", resolveGroupIDForRead(1337, d))
	})

	t.Run("does not panic when RawConfig is unavailable during read", func(t *testing.T) {
		d := schema.TestResourceDataRaw(t, datastreamResourceSchema, map[string]interface{}{
			"log_type": "CDN",
		})
		require.NoError(t, d.Set("group_id", "0"))

		assert.NotPanics(t, func() {
			assert.Empty(t, rawConfigStringAttr(d, "group_id"))
			assert.Empty(t, resolveGroupIDForRead(0, d))
			assert.Equal(t, "42", resolveGroupIDForRead(42, d))
		})
	})

	t.Run("preserves configured value when API omits groupId", func(t *testing.T) {
		d := schema.TestResourceDataRaw(t, datastreamResourceSchema, map[string]interface{}{
			"group_id": "grp_1337",
		})

		assert.Equal(t, "1337", resolveGroupIDForRead(0, d))
	})

	t.Run("preserves grp_ configured value when API returns 0 (I#775)", func(t *testing.T) {
		d := schema.TestResourceDataRaw(t, datastreamResourceSchema, map[string]interface{}{
			"group_id": "grp_123123123",
		})

		assert.Equal(t, "123123123", resolveGroupIDForRead(0, d))
	})

	t.Run("preserves numeric string configured value when API returns 0 (I#775)", func(t *testing.T) {
		d := schema.TestResourceDataRaw(t, datastreamResourceSchema, map[string]interface{}{
			"group_id": "123123123",
		})

		assert.Equal(t, "123123123", resolveGroupIDForRead(0, d))
	})

	t.Run("returns empty when API and config omit group_id", func(t *testing.T) {
		d := schema.TestResourceDataRaw(t, datastreamResourceSchema, map[string]interface{}{
			"log_type": "CDN",
		})

		assert.Empty(t, resolveGroupIDForRead(0, d))
	})

	t.Run("does not preserve explicit zero group_id", func(t *testing.T) {
		d := schema.TestResourceDataRaw(t, datastreamResourceSchema, map[string]interface{}{
			"group_id": "0",
		})

		assert.Empty(t, resolveGroupIDForRead(0, d))
	})

	t.Run("prefers config over stale state zero after I#775 upgrade", func(t *testing.T) {
		d := schema.TestResourceDataRaw(t, datastreamResourceSchema, map[string]interface{}{
			"group_id": "grp_123123123",
		})
		require.NoError(t, d.Set("group_id", "0"))

		// schema.TestResourceDataRaw often does not populate cty RawConfig the same way a real
		// terraform refresh does. When RawConfig is unavailable, stale "0" correctly resolves
		// to empty (not preserved). The real upgrade path is covered by resource UnitTests
		// (ExactBugReportSymptoms + NoDrift*) where RawConfig comes from the HCL fixture.
		if rawConfigStringAttr(d, "group_id") == "" {
			assert.Empty(t, resolveGroupIDForRead(0, d))
			return
		}

		assert.Equal(t, "123123123", resolveGroupIDForRead(0, d))
	})
}

func TestResolveIntegrationTypeForRead(t *testing.T) {
	t.Parallel()

	t.Run("prefers API value when present", func(t *testing.T) {
		d := schema.TestResourceDataRaw(t, datastreamResourceSchema, map[string]interface{}{
			"integration_type": "PM_DEPENDENT",
		})

		value, ok := resolveIntegrationTypeForRead("DS_MANAGED", d)
		assert.True(t, ok)
		assert.Equal(t, "DS_MANAGED", value)
	})

	t.Run("preserves configured value when API omits integration_type", func(t *testing.T) {
		d := schema.TestResourceDataRaw(t, datastreamResourceSchema, map[string]interface{}{
			"integration_type": "PM_DEPENDENT",
		})

		value, ok := resolveIntegrationTypeForRead("", d)
		assert.True(t, ok)
		assert.Equal(t, "PM_DEPENDENT", value)
	})

	t.Run("preserves DS_MANAGED from state when API omits integration_type (I#775)", func(t *testing.T) {
		d := schema.TestResourceDataRaw(t, datastreamResourceSchema, map[string]interface{}{
			"integration_type": "DS_MANAGED",
		})

		value, ok := resolveIntegrationTypeForRead("", d)
		assert.True(t, ok)
		assert.Equal(t, "DS_MANAGED", value)
	})
}

func TestResolveSamplingPercentageForRead(t *testing.T) {
	t.Parallel()

	t.Run("prefers API value when present", func(t *testing.T) {
		d := schema.TestResourceDataRaw(t, datastreamResourceSchema, map[string]interface{}{
			"sampling_percentage": 25,
		})

		value, ok := resolveSamplingPercentageForRead(50, d)
		assert.True(t, ok)
		assert.Equal(t, 50, value)
	})

	t.Run("preserves configured value when API omits sampling_percentage", func(t *testing.T) {
		d := schema.TestResourceDataRaw(t, datastreamResourceSchema, map[string]interface{}{
			"sampling_percentage": 25,
		})

		value, ok := resolveSamplingPercentageForRead(0, d)
		assert.True(t, ok)
		assert.Equal(t, 25, value)
	})
}
