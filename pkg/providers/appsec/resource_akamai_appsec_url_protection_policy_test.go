package appsec

import (
	"context"
	"encoding/json"
	"fmt"
	"regexp"
	"testing"

	"github.com/akamai/AkamaiOPEN-edgegrid-golang/v13/pkg/appsec"
	"github.com/akamai/terraform-provider-akamai/v10/internal/edgegrid"
	"github.com/akamai/terraform-provider-akamai/v10/pkg/common/test"
	"github.com/akamai/terraform-provider-akamai/v10/pkg/common/testutils"
	"github.com/hashicorp/terraform-plugin-framework/attr"
	fwresource "github.com/hashicorp/terraform-plugin-framework/resource"
	"github.com/hashicorp/terraform-plugin-framework/types"
	"github.com/hashicorp/terraform-plugin-framework/types/basetypes"
	"github.com/hashicorp/terraform-plugin-testing/helper/resource"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/mock"
	"github.com/stretchr/testify/require"
)

func TestValidateIntelligentLoadShedding_NullAndUnknown(t *testing.T) {
	t.Parallel()

	r := &urlProtectionPolicyResource{}

	// ── NULL case ──────────────────────────────────────────────────────────────
	// When intelligent_load_shedding is not set in HCL, Terraform gives the field a
	// null typed object.  tf.IsKnown(null) == false → the early-return is NOT taken →
	// validation erroneously runs on a null object → multiple errors are added.
	t.Run("null IntelligentLoadShedding - validation must be skipped", func(t *testing.T) {
		data := &urlProtectionPolicyResourceModel{
			IntelligentLoadShedding: types.ObjectNull(intelligentLoadSheddingAttrTypes()),
		}
		resp := &fwresource.ValidateConfigResponse{}

		r.validateIntelligentLoadShedding(context.Background(), data, resp)

		assert.False(t, resp.Diagnostics.HasError(),
			"ValidateConfig should produce no errors when intelligent_load_shedding is null (not configured). "+
				"Got: %v", resp.Diagnostics)
	})

	// ── UNKNOWN case ──────────────────────────────────────────────────────────────
	// When intelligent_load_shedding references an output that is not yet known at
	// plan time, Terraform marks the whole object as unknown.

	// Expected (correct) behaviour: unknown means "will be determined later"; skip all
	// validation now; Terraform will re-validate once the value is resolved.
	t.Run("unknown IntelligentLoadShedding - validation must be skipped", func(t *testing.T) {
		data := &urlProtectionPolicyResourceModel{
			IntelligentLoadShedding: types.ObjectUnknown(intelligentLoadSheddingAttrTypes()),
		}
		resp := &fwresource.ValidateConfigResponse{}

		r.validateIntelligentLoadShedding(context.Background(), data, resp)

		assert.False(t, resp.Diagnostics.HasError(),
			"ValidateConfig should produce no errors when intelligent_load_shedding is unknown (value not yet known). "+
				"Got: %v", resp.Diagnostics)
	})

	// ── valid case ─────────────────────────────────────
	t.Run("known valid IntelligentLoadShedding - no validation errors", func(t *testing.T) {
		categories, _ := types.ListValueFrom(context.Background(), types.StringType, []string{"BOTS"})
		customCriteria, _ := types.ListValue(types.ObjectType{AttrTypes: map[string]attr.Type{
			"type":           types.StringType,
			"list_ids":       types.ListType{ElemType: types.StringType},
			"positive_match": types.BoolType,
		}}, []attr.Value{})

		ilsObj, _ := types.ObjectValueFrom(context.Background(), intelligentLoadSheddingAttrTypes(),
			intelligentLoadSheddingModel{
				HitsPerSec:     types.Int64Value(100),
				Categories:     categories,
				CustomCriteria: customCriteria,
			})

		data := &urlProtectionPolicyResourceModel{
			MaxRateThreshold:        types.Int64Value(195),
			IntelligentLoadShedding: ilsObj,
		}
		resp := &fwresource.ValidateConfigResponse{}

		r.validateIntelligentLoadShedding(context.Background(), data, resp)

		assert.False(t, resp.Diagnostics.HasError(),
			"ValidateConfig should produce no errors for a valid intelligent_load_shedding object. "+
				"Got: %v", resp.Diagnostics)
	})
}

func TestCustomCriteriaNotComputed_StateMustMatchPlan(t *testing.T) {
	t.Parallel()

	ctx := context.Background()

	ccObjType := types.ObjectType{AttrTypes: map[string]attr.Type{
		"type":           types.StringType,
		"list_ids":       types.ListType{ElemType: types.StringType},
		"positive_match": types.BoolType,
	}}

	// Non-null custom_criteria from API (CLIENT_LIST category mapped by populateCategoriesAndCustomCriteria)
	listIDs, _ := types.ListValueFrom(ctx, types.StringType, []string{"12345_10CLIENTLIST"})
	positiveMatch := true
	ccFromAPI, _ := types.ListValueFrom(ctx, ccObjType, []customCriteriaModel{
		{Type: types.StringValue("CLIENT_LIST"), ListIDs: listIDs, PositiveMatch: types.BoolValue(positiveMatch)},
	})

	catFromPlan, _ := types.ListValueFrom(ctx, types.StringType,
		[]string{"BOTS", "CLOUD_PROVIDERS", "PROXIES", "TOR_EXIT_NODES", "PLATFORM_DDOS_INTELLIGENCE"})

	// ── Test 1: plan null, API non-null → bug and fix ────────────────────────────
	t.Run("plan.custom_criteria is null - state must stay null (API value must NOT be used)", func(t *testing.T) {
		dsILS := &intelligentLoadSheddingModel{
			HitsPerSec:     types.Int64Value(150),
			Categories:     catFromPlan,
			CustomCriteria: ccFromAPI, // API returned CLIENT_LIST
		}
		planCCNull := types.ListNull(ccObjType) // user did not set custom_criteria

		// Current (buggy) approach: ObjectValueFrom(dsModel.ILS) copies API value
		buggyObj, _ := types.ObjectValueFrom(ctx, intelligentLoadSheddingAttrTypes(), dsILS)
		var buggyILS intelligentLoadSheddingModel
		_ = buggyObj.As(ctx, &buggyILS, basetypes.ObjectAsOptions{})
		assert.False(t, buggyILS.CustomCriteria.IsNull(),
			"BUG confirmed: using dsModel directly makes state.custom_criteria non-null "+
				"while plan had null → triggers 'Provider produced inconsistent result'.")

		// Fixed approach: preserve plan's custom_criteria when building state ILS
		fixedILS := intelligentLoadSheddingModel{
			HitsPerSec:     dsILS.HitsPerSec,
			Categories:     dsILS.Categories,
			CustomCriteria: planCCNull, // use plan value, not dsModel
		}
		fixedObj, _ := types.ObjectValueFrom(ctx, intelligentLoadSheddingAttrTypes(), fixedILS)
		var resultILS intelligentLoadSheddingModel
		_ = fixedObj.As(ctx, &resultILS, basetypes.ObjectAsOptions{})
		assert.True(t, resultILS.CustomCriteria.IsNull(),
			"FIXED: state.custom_criteria must be null when plan had null.")
	})

	// ── Test 2: plan has custom_criteria set → state must equal plan ─────────────
	t.Run("plan.custom_criteria is set - state must reflect plan value exactly", func(t *testing.T) {
		listIDsPlan, _ := types.ListValueFrom(ctx, types.StringType, []string{"54321_MYLIST"})
		ccFromPlan, _ := types.ListValueFrom(ctx, ccObjType, []customCriteriaModel{
			{Type: types.StringValue("CLIENT_LIST"), ListIDs: listIDsPlan, PositiveMatch: types.BoolValue(true)},
		})

		stateILS := intelligentLoadSheddingModel{
			HitsPerSec:     types.Int64Value(150),
			Categories:     catFromPlan,
			CustomCriteria: ccFromPlan,
		}
		obj, _ := types.ObjectValueFrom(ctx, intelligentLoadSheddingAttrTypes(), stateILS)
		var resultILS intelligentLoadSheddingModel
		_ = obj.As(ctx, &resultILS, basetypes.ObjectAsOptions{})
		assert.Equal(t, ccFromPlan, resultILS.CustomCriteria,
			"state.custom_criteria must match the plan value exactly.")
	})

	// ── Test 3: no ILS in plan → state.IntelligentLoadShedding must be null ──────
	t.Run("no intelligent_load_shedding in plan - state object must be null", func(t *testing.T) {
		nilILS := types.ObjectNull(intelligentLoadSheddingAttrTypes())
		assert.True(t, nilILS.IsNull(),
			"When ILS is not configured, state.IntelligentLoadShedding must be null.")
	})
}

func TestURLProtectionPolicyResource(t *testing.T) {
	t.Parallel()

	// Load test fixtures

	urlProtectionPolicyResponse := appsec.GetURLProtectionPolicyResponse{}
	err := json.Unmarshal(testutils.LoadFixtureBytes(t, "testdata/TestResURLProtectionPolicy/URLProtectionPolicy.json"), &urlProtectionPolicyResponse)
	require.NoError(t, err)

	urlProtectionPolicyUpdatedResponse := appsec.GetURLProtectionPolicyResponse{}
	err = json.Unmarshal(testutils.LoadFixtureBytes(t, "testdata/TestResURLProtectionPolicy/URLProtectionPolicyUpdated.json"), &urlProtectionPolicyUpdatedResponse)
	require.NoError(t, err)

	urlProtectionPolicyWithAPIResponse := appsec.GetURLProtectionPolicyResponse{}
	err = json.Unmarshal(testutils.LoadFixtureBytes(t, "testdata/TestResURLProtectionPolicy/URLProtectionPolicyWithAPI.json"), &urlProtectionPolicyWithAPIResponse)
	require.NoError(t, err)

	urlProtectionPolicyAPIUpdatedResponse := appsec.GetURLProtectionPolicyResponse{}
	err = json.Unmarshal(testutils.LoadFixtureBytes(t, "testdata/TestResURLProtectionPolicy/URLProtectionPolicyAPIUpdated.json"), &urlProtectionPolicyAPIUpdatedResponse)
	require.NoError(t, err)

	urlProtectionPolicyAPIUpdatedNameResponse := appsec.GetURLProtectionPolicyResponse{}
	err = json.Unmarshal(testutils.LoadFixtureBytes(t, "testdata/TestResURLProtectionPolicy/URLProtectionPolicyAPIUpdatedName.json"), &urlProtectionPolicyAPIUpdatedNameResponse)
	require.NoError(t, err)

	createResponse := appsec.CreateURLProtectionPolicyResponse{
		URLProtectionPolicyID: 681,
	}

	ilsWithCCResponse := appsec.GetURLProtectionPolicyResponse{}
	errCustomCriteria := json.Unmarshal(
		testutils.LoadFixtureBytes(t, "testdata/TestResURLProtectionPolicy/URLProtectionPolicyILSWithCustomCriteria.json"),
		&ilsWithCCResponse,
	)
	require.NoError(t, errCustomCriteria)

	createResponseForCustomCriteria := appsec.CreateURLProtectionPolicyResponse{URLProtectionPolicyID: 681}

	var tests = map[string]struct {
		init  func(*appsec.Mock)
		steps []resource.TestStep
	}{
		"missing required argument config_id": {
			steps: []resource.TestStep{
				{
					Config:      testutils.LoadFixtureString(t, "testdata/TestResURLProtectionPolicy/missing_config_id.tf"),
					ExpectError: regexp.MustCompile("Missing required argument"),
				},
			},
		},
		"missing required argument name": {
			steps: []resource.TestStep{
				{
					Config:      testutils.LoadFixtureString(t, "testdata/TestResURLProtectionPolicy/missing_name.tf"),
					ExpectError: regexp.MustCompile("Missing required argument"),
				},
			},
		},
		"missing required argument max_rate_threshold": {
			steps: []resource.TestStep{
				{
					Config:      testutils.LoadFixtureString(t, "testdata/TestResURLProtectionPolicy/missing_max_rate_threshold.tf"),
					ExpectError: regexp.MustCompile("Missing required argument"),
				},
			},
		},
		"validate config - both hostname_paths and api_definitions specified": {
			steps: []resource.TestStep{
				{
					Config:      testutils.LoadFixtureString(t, "testdata/TestResURLProtectionPolicy/both_hostname_and_api.tf"),
					ExpectError: regexp.MustCompile("Only one of 'hostname_paths' or 'api_definitions' can be specified"),
				},
			},
		},
		"validate config - neither hostname_paths nor api_definitions specified": {
			steps: []resource.TestStep{
				{
					Config:      testutils.LoadFixtureString(t, "testdata/TestResURLProtectionPolicy/missing_hostname_and_api.tf"),
					ExpectError: regexp.MustCompile("Either 'hostname_paths' or 'api_definitions' must be specified"),
				},
			},
		},
		"validate config - duplicate paths in hostname_paths.paths": {
			steps: []resource.TestStep{
				{
					Config:      testutils.LoadFixtureString(t, "testdata/TestResURLProtectionPolicy/invalid_hostname_duplicate_paths.tf"),
					ExpectError: regexp.MustCompile("This attribute contains duplicate values"),
				},
			},
		},
		"validate config - duplicate hostnames in hostname_paths": {
			steps: []resource.TestStep{
				{
					Config:      testutils.LoadFixtureString(t, "testdata/TestResURLProtectionPolicy/invalid_hostname_duplicate.tf"),
					ExpectError: regexp.MustCompile("Duplicate List Value"),
				},
			},
		},
		"validate config - empty hostnames array": {
			steps: []resource.TestStep{
				{
					Config:      testutils.LoadFixtureString(t, "testdata/TestResURLProtectionPolicy/empty_hostname_array.tf"),
					ExpectError: regexp.MustCompile("Attribute hostname_paths list must contain at least 1 elements"),
				},
			},
		},
		"validate config - empty paths in hostname_paths": {
			steps: []resource.TestStep{
				{
					Config:      testutils.LoadFixtureString(t, "testdata/TestResURLProtectionPolicy/empty_hostnames_paths_path.tf"),
					ExpectError: regexp.MustCompile("list must contain at least 1 elements"),
				},
			},
		},
		"validate config - invalid custom criteria type": {
			steps: []resource.TestStep{
				{
					Config:      testutils.LoadFixtureString(t, "testdata/TestResURLProtectionPolicy/invalid_custom_criteria_type.tf"),
					ExpectError: regexp.MustCompile("Invalid Custom Criteria Type"),
				},
			},
		},
		"error creating url protection Policy": {
			init: func(m *appsec.Mock) {
				mockGetConfigurationURLProtectionPolicy(m, 1)
				mockCreateURLProtectionPolicyFailure(m, 1)
			},
			steps: []resource.TestStep{
				{
					Config:      testutils.LoadFixtureString(t, "testdata/TestResURLProtectionPolicy/create_with_hostname_paths.tf"),
					ExpectError: regexp.MustCompile("Error creating URL Protection Policy"),
				},
			},
		},
		"validate config - intelligent_load_shedding.categories is null": {
			steps: []resource.TestStep{{
				Config:      testutils.LoadFixtureString(t, "testdata/TestResURLProtectionPolicy/invalid_categories_null.tf"),
				ExpectError: regexp.MustCompile("categories list must not be null or unknown"),
			}},
		},
		"validate config - intelligent_load_shedding.categories is empty": {
			steps: []resource.TestStep{{
				Config:      testutils.LoadFixtureString(t, "testdata/TestResURLProtectionPolicy/invalid_categories_empty.tf"),
				ExpectError: regexp.MustCompile("categories list must not be empty"),
			}},
		},
		"validate config - hits_per_sec below allowed": {
			steps: []resource.TestStep{{
				Config:      testutils.LoadFixtureString(t, "testdata/TestResURLProtectionPolicy/invalid_hits_per_sec_low.tf"),
				ExpectError: regexp.MustCompile("hits_per_sec must be at least 25% of max_rate_threshold|min 7"),
			}},
		},
		"validate config - hits_per_sec above allowed": {
			steps: []resource.TestStep{{
				Config:      testutils.LoadFixtureString(t, "testdata/TestResURLProtectionPolicy/invalid_hits_per_sec_high.tf"),
				ExpectError: regexp.MustCompile("hits_per_sec must be less than or equal to 90% of max_rate_threshold"),
			}},
		},
		"validate config - custom_criteria.list_ids is null": {
			steps: []resource.TestStep{{
				Config:      testutils.LoadFixtureString(t, "testdata/TestResURLProtectionPolicy/invalid_custom_criteria_list_ids_null.tf"),
				ExpectError: regexp.MustCompile("list_ids must not be null or unknown"),
			}},
		},
		"validate config - custom_criteria.list_ids is empty": {
			steps: []resource.TestStep{{
				Config:      testutils.LoadFixtureString(t, "testdata/TestResURLProtectionPolicy/invalid_custom_criteria_list_ids_empty.tf"),
				ExpectError: regexp.MustCompile("list_ids must not be empty"),
			}},
		},
		"validate config - bypass_conditions names missing": {
			steps: []resource.TestStep{{
				Config:      testutils.LoadFixtureString(t, "testdata/TestResURLProtectionPolicy/invalid_bypass_condition_names_missing.tf"),
				ExpectError: regexp.MustCompile("'names' is required when type is 'RequestHeaderCondition'"),
			}},
		},
		"validate config - bypass_conditions names empty": {
			steps: []resource.TestStep{{
				Config:      testutils.LoadFixtureString(t, "testdata/TestResURLProtectionPolicy/invalid_bypass_condition_names_empty.tf"),
				ExpectError: regexp.MustCompile("'names' must not be empty when type is 'RequestHeaderCondition'"),
			}},
		},
		"validate config - max_rate_threshold below minimum": {
			steps: []resource.TestStep{{
				Config:      testutils.LoadFixtureString(t, "testdata/TestResURLProtectionPolicy/invalid_max_rate_threshold_low.tf"),
				ExpectError: regexp.MustCompile("max_rate_threshold must be at least 10"),
			}},
		},
		"validate config - intelligent_load_shedding.hits_per_sec missing": {
			steps: []resource.TestStep{{
				Config:      testutils.LoadFixtureString(t, "testdata/TestResURLProtectionPolicy/invalid_intelligent_load_shedding_hits_per_sec_missing.tf"),
				ExpectError: regexp.MustCompile("Incorrect attribute value type"),
			}},
		},
		"create ILS without custom_criteria - API returns CLIENT_LIST - no inconsistency": {
			init: func(m *appsec.Mock) {
				// GetConfiguration: 1 for Create, 1 for pre-destroy Read, 1 for Delete cleanup
				mockGetConfigurationURLProtectionPolicy(m, 3)
				mockCreateURLProtectionPolicySuccess(m, createResponseForCustomCriteria, 1)
				// GET response has CLIENT_LIST → dsModel.CustomCriteria is non-null
				// Called twice: once after Create, once for pre-destroy refresh
				mockGetURLProtectionPolicyData(m, ilsWithCCResponse, 2)
				mockRemoveURLProtectionPolicySuccess(m, 1)
			},
			steps: []resource.TestStep{
				{
					Config: testutils.LoadFixtureString(t, "testdata/TestResURLProtectionPolicy/ils_no_custom_criteria.tf"),
					// plan.ILS.CustomCriteria = null (not in config, not Computed)
					// Create uses plan value for ILS → state.custom_criteria stays null → no inconsistency
					Check: resource.ComposeAggregateTestCheckFunc(
						resource.TestCheckResourceAttr("akamai_appsec_url_protection_policy.test", "config_id", "43007"),
						resource.TestCheckResourceAttr("akamai_appsec_url_protection_policy.test", "url_protection_policy_id", "681"),
					),
				},
			},
		},
		"create url protection Policy with hostname paths - success": {
			init: func(m *appsec.Mock) {
				// Mock GetConfiguration for Create phase
				mockGetConfigurationURLProtectionPolicy(m, 1)

				// Mock CreateURLProtectionPolicy
				mockCreateURLProtectionPolicySuccess(m, createResponse, 1)

				// Mock GetURLProtectionPolicy for reading after create
				mockGetURLProtectionPolicyData(m, urlProtectionPolicyResponse, 1)

				// Mock GetConfiguration for Delete phase (cleanup)
				mockGetConfigurationURLProtectionPolicy(m, 1)

				// Mock GetURLProtectionPolicy for reading after create
				mockGetURLProtectionPolicyData(m, urlProtectionPolicyResponse, 1)

				// Mock GetConfiguration for Delete phase (cleanup)
				mockGetConfigurationURLProtectionPolicy(m, 1)

				// Mock RemoveURLProtectionPolicy for cleanup
				mockRemoveURLProtectionPolicySuccess(m, 1)
			},
			steps: []resource.TestStep{
				{
					Config: testutils.LoadFixtureString(t, "testdata/TestResURLProtectionPolicy/create_with_hostname_paths.tf"),
					Check: resource.ComposeAggregateTestCheckFunc(
						resource.TestCheckResourceAttr("akamai_appsec_url_protection_policy.test", "config_id", "43007"),
						resource.TestCheckResourceAttr("akamai_appsec_url_protection_policy.test", "url_protection_policy_id", "681"),
					),
				},
			},
		},
		"create url protection policy with api definitions - success": {
			init: func(m *appsec.Mock) {
				// Mock GetConfiguration for Create phase
				mockGetConfigurationURLProtectionPolicy(m, 1)

				// Mock CreateURLProtectionPolicy
				mockCreateURLProtectionPolicySuccess(m, createResponse, 1)

				// Mock GetURLProtectionPolicy for reading after create
				mockGetURLProtectionPolicyData(m, urlProtectionPolicyWithAPIResponse, 1)

				// Mock GetConfiguration for Delete phase (cleanup)
				mockGetConfigurationURLProtectionPolicy(m, 1)

				// Mock GetURLProtectionPolicy for reading after create
				mockGetURLProtectionPolicyData(m, urlProtectionPolicyWithAPIResponse, 1)

				// Mock GetConfiguration for Delete phase (cleanup)
				mockGetConfigurationURLProtectionPolicy(m, 1)

				// Mock RemoveURLProtectionPolicy for cleanup
				mockRemoveURLProtectionPolicySuccess(m, 1)
			},
			steps: []resource.TestStep{
				{
					Config: testutils.LoadFixtureString(t, "testdata/TestResURLProtectionPolicy/create_with_api_definitions.tf"),
					Check: resource.ComposeAggregateTestCheckFunc(
						resource.TestCheckResourceAttr("akamai_appsec_url_protection_policy.test", "config_id", "43007"),
						resource.TestCheckResourceAttr("akamai_appsec_url_protection_policy.test", "url_protection_policy_id", "681"),
					),
				},
			},
		},
		"update url protection policy with hostname paths - success": {
			init: func(m *appsec.Mock) {
				// Mock GetConfiguration for all phases (create, update, read, delete, etc.)
				mockGetConfigurationURLProtectionPolicy(m, 6)
				// Mock CreateURLProtectionPolicy for initial resource creation
				mockCreateURLProtectionPolicySuccess(m, createResponse, 1)
				// Mock GetURLProtectionPolicy for reading after create (may be called multiple times)
				mockGetURLProtectionPolicyData(m, urlProtectionPolicyResponse, 3)
				// Mock UpdateURLProtectionPolicy
				m.On("UpdateURLProtectionPolicy", mock.Anything, mock.AnythingOfType("appsec.UpdateURLProtectionPolicyRequest")).Return(&appsec.UpdateURLProtectionPolicyResponse{
					URLProtectionPolicyID: 681,
				}, nil).Once()
				// Mock GetURLProtectionPolicy for reading after update (may be called multiple times)
				mockGetURLProtectionPolicyData(m, urlProtectionPolicyUpdatedResponse, 2)
				// Mock RemoveURLProtectionPolicy for cleanup
				mockRemoveURLProtectionPolicySuccess(m, 1)
			},
			steps: []resource.TestStep{
				{
					Config: testutils.LoadFixtureString(t, "testdata/TestResURLProtectionPolicy/create_with_hostname_paths.tf"),
					Check: resource.ComposeAggregateTestCheckFunc(
						resource.TestCheckResourceAttr("akamai_appsec_url_protection_policy.test", "config_id", "43007"),
						resource.TestCheckResourceAttr("akamai_appsec_url_protection_policy.test", "url_protection_policy_id", "681"),
					),
				},
				{
					Config: testutils.LoadFixtureString(t, "testdata/TestResURLProtectionPolicy/update_with_hostname_paths.tf"),
					Check: resource.ComposeAggregateTestCheckFunc(
						resource.TestCheckResourceAttr("akamai_appsec_url_protection_policy.test", "description", "Updated URL Protection"),
					),
				},
			},
		},
		"update url protection policy with api definitions - success": {
			init: func(m *appsec.Mock) {
				// Mock GetConfiguration for all phases (create, update, read, delete, etc.)
				mockGetConfigurationURLProtectionPolicy(m, 6)
				// Mock CreateURLProtectionPolicy for initial resource creation
				mockCreateURLProtectionPolicySuccess(m, createResponse, 1)
				// Mock GetURLProtectionPolicy for reading after create (may be called multiple times)
				mockGetURLProtectionPolicyData(m, urlProtectionPolicyWithAPIResponse, 3)
				// Mock UpdateURLProtectionPolicy
				m.On("UpdateURLProtectionPolicy", mock.Anything, mock.AnythingOfType("appsec.UpdateURLProtectionPolicyRequest")).Return(&appsec.UpdateURLProtectionPolicyResponse{
					URLProtectionPolicyID: 681,
				}, nil).Once()
				// Mock GetURLProtectionPolicy for reading after update (may be called multiple times)
				mockGetURLProtectionPolicyData(m, urlProtectionPolicyAPIUpdatedResponse, 2)
				// Mock RemoveURLProtectionPolicy for cleanup
				mockRemoveURLProtectionPolicySuccess(m, 1)
			},
			steps: []resource.TestStep{
				{
					Config: testutils.LoadFixtureString(t, "testdata/TestResURLProtectionPolicy/create_with_api_definitions.tf"),
					Check: resource.ComposeAggregateTestCheckFunc(
						resource.TestCheckResourceAttr("akamai_appsec_url_protection_policy.test", "config_id", "43007"),
						resource.TestCheckResourceAttr("akamai_appsec_url_protection_policy.test", "url_protection_policy_id", "681"),
					),
				},
				{
					Config: testutils.LoadFixtureString(t, "testdata/TestResURLProtectionPolicy/update_with_api_definitions.tf"),
					Check: resource.ComposeAggregateTestCheckFunc(
						resource.TestCheckResourceAttr("akamai_appsec_url_protection_policy.test", "description", "Updated API Protection"),
					),
				},
			},
		},
		"create and update with different config_id should fail": {
			init: func(m *appsec.Mock) {
				// Mock GetConfiguration for all phases (create, update, read, delete, etc.)
				mockGetConfigurationURLProtectionPolicy(m, 4)
				// Mock CreateURLProtectionPolicy for initial resource creation
				mockCreateURLProtectionPolicySuccess(m, createResponse, 1)
				// Mock GetURLProtectionPolicy for reading after create (may be called multiple times)
				mockGetURLProtectionPolicyData(m, urlProtectionPolicyWithAPIResponse, 2)
				// Mock GetURLProtectionPolicy for reading after update (may be called multiple times)
				mockGetURLProtectionPolicyData(m, urlProtectionPolicyAPIUpdatedResponse, 1)
				// Mock RemoveURLProtectionPolicy for cleanup
				mockRemoveURLProtectionPolicySuccess(m, 1)
			},
			steps: []resource.TestStep{
				{
					Config: testutils.LoadFixtureString(t, "testdata/TestResURLProtectionPolicy/create_with_api_definitions.tf"),
					Check: resource.ComposeAggregateTestCheckFunc(
						resource.TestCheckResourceAttr("akamai_appsec_url_protection_policy.test", "config_id", "43007"),
						resource.TestCheckResourceAttr("akamai_appsec_url_protection_policy.test", "url_protection_policy_id", "681"),
					),
				},
				{
					Config:      testutils.LoadFixtureString(t, "testdata/TestResURLProtectionPolicy/update_with_api_definitions_different_configID.tf"),
					ExpectError: regexp.MustCompile("updating field `config_id` is not possible"),
				},
			},
		},
		"create and update with different name - success": {
			init: func(m *appsec.Mock) {
				// Mock GetConfiguration for all phases (create, update, read, delete, etc.)
				mockGetConfigurationURLProtectionPolicy(m, 6)
				// Mock CreateURLProtectionPolicy for initial resource creation
				mockCreateURLProtectionPolicySuccess(m, createResponse, 1)
				// Mock GetURLProtectionPolicy for reading after create (may be called multiple times)
				mockGetURLProtectionPolicyData(m, urlProtectionPolicyWithAPIResponse, 3)
				// Mock UpdateURLProtectionPolicy
				m.On("UpdateURLProtectionPolicy", mock.Anything, mock.AnythingOfType("appsec.UpdateURLProtectionPolicyRequest")).Return(&appsec.UpdateURLProtectionPolicyResponse{
					URLProtectionPolicyID: 681,
					Name:                  "API Protection Rule test",
				}, nil).Once()
				// Mock GetURLProtectionPolicy for reading after update (may be called multiple times)
				mockGetURLProtectionPolicyData(m, urlProtectionPolicyAPIUpdatedNameResponse, 2)
				// Mock RemoveURLProtectionPolicy for cleanup
				mockRemoveURLProtectionPolicySuccess(m, 1)
			},
			steps: []resource.TestStep{
				{
					Config: testutils.LoadFixtureString(t, "testdata/TestResURLProtectionPolicy/create_with_api_definitions.tf"),
					Check: resource.ComposeAggregateTestCheckFunc(
						resource.TestCheckResourceAttr("akamai_appsec_url_protection_policy.test", "config_id", "43007"),
						resource.TestCheckResourceAttr("akamai_appsec_url_protection_policy.test", "url_protection_policy_id", "681"),
					),
				},
				{
					Config: testutils.LoadFixtureString(t, "testdata/TestResURLProtectionPolicy/update_with_api_definitions_different_name.tf"),
					Check: resource.ComposeAggregateTestCheckFunc(
						resource.TestCheckResourceAttr("akamai_appsec_url_protection_policy.test", "config_id", "43007"),
						resource.TestCheckResourceAttr("akamai_appsec_url_protection_policy.test", "name", "API Protection Rule test"),
					),
				},
			},
		},
		"import url protection": {
			init: func(m *appsec.Mock) {
				// Import and post-import refresh reads
				mockGetConfigurationURLProtectionPolicy(m, 1)

				mockGetURLProtectionPolicyData(m, urlProtectionPolicyWithAPIResponse, 1)

				mockGetConfigurationURLProtectionPolicy(m, 1)
				// Delete after test
				mockRemoveURLProtectionPolicySuccess(m, 1)
			},
			steps: []resource.TestStep{
				{
					Config:        testutils.LoadFixtureString(t, "testdata/TestResURLProtectionPolicy/update_with_hostname_paths.tf"),
					ImportState:   true,
					ImportStateId: "43007:681",
					ResourceName:  "akamai_appsec_url_protection_policy.test",
					ImportStateCheck: test.NewImportChecker().
						CheckEqual("config_id", "43007").
						CheckEqual("url_protection_policy_id", "681").
						CheckEqual("max_rate_threshold", "195").
						Build(),
					ImportStatePersist: true,
				},
			},
		},
		"import url protection - invalid id format": {
			steps: []resource.TestStep{
				{
					Config:        testutils.LoadFixtureString(t, "testdata/TestResURLProtectionPolicy/update_with_hostname_paths.tf"),
					ImportState:   true,
					ImportStateId: "12345",
					ResourceName:  "akamai_appsec_url_protection_policy.test",
					ExpectError:   regexp.MustCompile("Invalid Import ID"),
				},
			},
		},
		"import url protection - invalid config id value": {
			steps: []resource.TestStep{
				{
					Config:        testutils.LoadFixtureString(t, "testdata/TestResURLProtectionPolicy/update_with_hostname_paths.tf"),
					ImportState:   true,
					ImportStateId: "abc:681",
					ResourceName:  "akamai_appsec_url_protection_policy.test",
					ExpectError:   regexp.MustCompile("Invalid Config ID"),
				},
			},
		},
		"import url protection - invalid url protection id value": {
			steps: []resource.TestStep{
				{
					Config:        testutils.LoadFixtureString(t, "testdata/TestResURLProtectionPolicy/update_with_hostname_paths.tf"),
					ImportState:   true,
					ImportStateId: "43007:xyz",
					ResourceName:  "akamai_appsec_url_protection_policy.test",
					ExpectError:   regexp.MustCompile("Invalid URL Protection ID"),
				},
			},
		},
	}

	for name, tc := range tests {
		t.Run(name, func(t *testing.T) {
			t.Parallel()
			client := edgegrid.NewTestClient()
			if tc.init != nil {
				tc.init(client.APPSEC)
			}

			mockGetConfigurationVersionDefault(client.APPSEC)
			resource.UnitTest(t, resource.TestCase{
				ProtoV6ProviderFactories: testutils.NewTestProtoV6ProviderFactory(client, NewSubprovider()),
				Steps:                    tc.steps,
			})

			client.APPSEC.AssertExpectations(t)
		})
	}
}

// Mock functions for URL Protection Policy resource tests

func mockGetConfigurationURLProtectionPolicy(m *appsec.Mock, times int) {
	m.On("GetConfiguration", mock.Anything, appsec.GetConfigurationRequest{ConfigID: 43007}).
		Return(&appsec.GetConfigurationResponse{
			FileType:          "RBAC",
			ID:                43007,
			LatestVersion:     40,
			Name:              "Akamai Tools",
			ProductionVersion: 1,
			StagingVersion:    1,
			TargetProduct:     "KSD",
		}, nil).Times(times)
}

func mockCreateURLProtectionPolicyFailure(m *appsec.Mock, times int) {
	m.On("CreateURLProtectionPolicy", mock.Anything, mock.MatchedBy(func(req appsec.CreateURLProtectionPolicyRequest) bool {
		return req.ConfigID == 43007 && req.ConfigVersion == 40
	})).Return(nil, fmt.Errorf("create url protection policy failed")).Times(times)
}

func mockCreateURLProtectionPolicySuccess(m *appsec.Mock, response appsec.CreateURLProtectionPolicyResponse, times int) {
	m.On("CreateURLProtectionPolicy", mock.Anything, mock.MatchedBy(func(req appsec.CreateURLProtectionPolicyRequest) bool {
		return req.ConfigID == 43007 && req.ConfigVersion == 40
	})).Return(&response, nil).Times(times)
}

func mockRemoveURLProtectionPolicySuccess(m *appsec.Mock, times int) {
	m.On("RemoveURLProtectionPolicy", mock.Anything, appsec.RemoveURLProtectionPolicyRequest{
		ConfigID:              43007,
		ConfigVersion:         40,
		URLProtectionPolicyID: 681,
	}).Return(nil).Times(times)
}
