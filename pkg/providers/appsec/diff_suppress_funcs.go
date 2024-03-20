package appsec

import (
	"bytes"
	"encoding/json"
	"fmt"
	"log"
	"reflect"
	"sort"
	"strings"

	"github.com/akamai/AkamaiOPEN-edgegrid-golang/v8/pkg/appsec"
	"github.com/akamai/terraform-provider-akamai/v5/pkg/logger"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/schema"
)

func suppressEquivalentJSONDiffsGeneric(_, oldString, newString string, _ *schema.ResourceData) bool {
	var ob, nb bytes.Buffer
	if err := json.Compact(&ob, []byte(oldString)); err != nil {
		return false
	}

	if err := json.Compact(&nb, []byte(newString)); err != nil {
		return false
	}

	return jsonBytesEqual(ob.Bytes(), nb.Bytes())
}

func jsonBytesEqual(b1, b2 []byte) bool {
	var o1 interface{}
	if err := json.Unmarshal(b1, &o1); err != nil {
		return false
	}

	var o2 interface{}
	if err := json.Unmarshal(b2, &o2); err != nil {
		return false
	}

	return reflect.DeepEqual(o1, o2)
}

func suppressEquivalentReputationProfileDiffs(_, oldVal, newVal string, _ *schema.ResourceData) bool {
	var rpOld, rpNew appsec.CreateReputationProfileResponse

	if err := json.Unmarshal([]byte(oldVal), &rpOld); err != nil {
		log.Printf("unable to unmarshal old reputation profile: %s", err)
	}
	if err := json.Unmarshal([]byte(newVal), &rpNew); err != nil {
		log.Printf("unable to unmarshal new reputation profile: %s", err)
	}

	if rpOld.Context != rpNew.Context ||
		rpOld.Name != rpNew.Name ||
		rpOld.SharedIPHandling != rpNew.SharedIPHandling ||
		rpOld.Threshold != rpNew.Threshold {
		return false
	}

	return compareReputationProfileCondition(rpOld, rpNew)
}

func compareReputationProfileCondition(rpOld, rpNew appsec.CreateReputationProfileResponse) bool {
	cOld := rpOld.Condition
	cNew := rpNew.Condition

	var pmOld = true
	var pmNew = true
	if cOld.PositiveMatch != nil {
		if err := json.Unmarshal(*cOld.PositiveMatch, &pmOld); err != nil {
			return false
		}
	}
	if cNew.PositiveMatch != nil {
		if err := json.Unmarshal(*cNew.PositiveMatch, &pmNew); err != nil {
			return false
		}
	}
	if pmOld != pmNew {
		return false
	}
	if len(cOld.AtomicConditions) != len(cNew.AtomicConditions) {
		return false
	}

	return areReputationProfilesEqual(rpOld, rpNew)
}

// areReputationProfilesEqual check whether old and new reputation profiles are the same
func areReputationProfilesEqual(oldProfile, newProfile appsec.CreateReputationProfileResponse) bool {
	oldCond := oldProfile.Condition
	newCond := newProfile.Condition
	for _, oldAtomicCond := range oldCond.AtomicConditions {
		found := false
		for _, newAtomicCond := range newCond.AtomicConditions {
			if oldAtomicCond.ClassName == newAtomicCond.ClassName {
				found = true
				if oldAtomicCond.CheckIps != newAtomicCond.CheckIps && newAtomicCond.CheckIps != "" {
					return false
				}
				if oldAtomicCond.NameCase != newAtomicCond.NameCase && !(oldAtomicCond.NameCase && !newAtomicCond.NameCase) {
					return false
				}
				if oldAtomicCond.NameWildcard != newAtomicCond.NameWildcard && !(oldAtomicCond.NameWildcard && !newAtomicCond.NameWildcard) {
					return false
				}
				if oldAtomicCond.ValueCase != newAtomicCond.ValueCase && !(oldAtomicCond.ValueCase && !newAtomicCond.ValueCase) {
					return false
				}
				if oldAtomicCond.ValueWildcard != newAtomicCond.ValueWildcard && !(oldAtomicCond.ValueWildcard && !newAtomicCond.ValueWildcard) {
					return false
				}
				if oldAtomicCond.PositiveMatch != newAtomicCond.PositiveMatch && !(oldAtomicCond.PositiveMatch && !newAtomicCond.PositiveMatch) {
					// only 'true' is supported for this case
					if oldAtomicCond.ClassName != "HostCondition" {
						return false
					}
				}
				if !suppressAtomicConditionSliceDiffs(oldAtomicCond.Value, newAtomicCond.Value) {
					return false
				}
				if !suppressAtomicConditionSliceDiffs(oldAtomicCond.Host, newAtomicCond.Host) {
					return false
				}
			}
		}
		if !found {
			return false
		}
	}
	return true
}

func suppressAtomicConditionSliceDiffs(oldSlice, newSlice []string) bool {
	if len(oldSlice) != len(newSlice) {
		return false
	}
	for _, ov := range oldSlice {
		found := false
		for _, nv := range newSlice {
			if strings.EqualFold(ov, nv) {
				found = true
			}
		}
		if !found {
			return false
		}
	}
	return true
}

func suppressEquivalentLoggingSettingsDiffs(_, oldString, newString string, _ *schema.ResourceData) bool {
	var oldJSON, newJSON appsec.UpdateAdvancedSettingsLoggingResponse
	if oldString == newString {
		return true
	}
	if err := json.Unmarshal([]byte(oldString), &oldJSON); err != nil {
		return false
	}
	if err := json.Unmarshal([]byte(newString), &newJSON); err != nil {
		return false
	}
	diff := compareLoggingSettings(&oldJSON, &newJSON)
	return diff
}

func compareLoggingSettings(oldResponse, newResponse *appsec.UpdateAdvancedSettingsLoggingResponse) bool {
	if oldResponse.Override != newResponse.Override ||
		oldResponse.AllowSampling != newResponse.AllowSampling ||
		oldResponse.Cookies.Type != newResponse.Cookies.Type ||
		oldResponse.CustomHeaders.Type != newResponse.CustomHeaders.Type ||
		oldResponse.StandardHeaders.Type != newResponse.StandardHeaders.Type {
		return false
	}

	sort.Strings(oldResponse.Cookies.Values)
	sort.Strings(newResponse.Cookies.Values)
	sort.Strings(oldResponse.CustomHeaders.Values)
	sort.Strings(newResponse.CustomHeaders.Values)
	sort.Strings(oldResponse.StandardHeaders.Values)
	sort.Strings(newResponse.StandardHeaders.Values)

	return reflect.DeepEqual(oldResponse, newResponse)
}

func suppressEquivalentAttackPayloadLoggingSettingsDiffs(_, oldValue, newValue string, _ *schema.ResourceData) bool {
	var oldJSON, newJSON appsec.UpdateAdvancedSettingsAttackPayloadLoggingResponse
	if oldValue == newValue {
		return true
	}
	if err := json.Unmarshal([]byte(oldValue), &oldJSON); err != nil {
		return false
	}
	if err := json.Unmarshal([]byte(newValue), &newJSON); err != nil {
		return false
	}
	diff := compareAttackPayloadLoggingSettings(&oldJSON, &newJSON)
	return diff
}

func compareAttackPayloadLoggingSettings(oldValue, newValue *appsec.UpdateAdvancedSettingsAttackPayloadLoggingResponse) bool {
	if oldValue.Override != newValue.Override ||
		oldValue.Enabled != newValue.Enabled ||
		oldValue.RequestBody.Type != newValue.RequestBody.Type ||
		oldValue.ResponseBody.Type != newValue.ResponseBody.Type {
		return false
	}

	return reflect.DeepEqual(oldValue, newValue)
}

func suppressEquivalentPenaltyBoxConditionsDiffs(_, oldValue, newValue string, _ *schema.ResourceData) bool {
	var oldJSON, newJSON appsec.GetPenaltyBoxConditionsResponse
	if oldValue == newValue {
		return true
	}
	if err := json.Unmarshal([]byte(oldValue), &oldJSON); err != nil {
		return false
	}
	if err := json.Unmarshal([]byte(newValue), &newJSON); err != nil {
		return false
	}
	diff := comparePenaltyBoxConditions(&oldJSON, &newJSON)
	return diff
}

func comparePenaltyBoxConditions(oldValue, newValue *appsec.GetPenaltyBoxConditionsResponse) bool {
	if oldValue.ConditionOperator != newValue.ConditionOperator {
		return false
	}

	return reflect.DeepEqual(oldValue, newValue)
}

func suppressCustomDenyJSONDiffs(_, oldString, newString string, _ *schema.ResourceData) bool {
	var ob, nb bytes.Buffer
	if err := json.Compact(&ob, []byte(oldString)); err != nil {
		return false
	}

	if err := json.Compact(&nb, []byte(newString)); err != nil {
		return false
	}

	return jsonBytesEqualIncludingParametersSlice(ob.Bytes(), nb.Bytes())
}

func jsonBytesEqualIncludingParametersSlice(b1, b2 []byte) bool {
	var o1 appsec.GetCustomDenyResponse
	var o2 appsec.GetCustomDenyResponse
	if err := json.Unmarshal(b1, &o1); err != nil {
		return false
	}

	if err := json.Unmarshal(b2, &o2); err != nil {
		return false
	}

	sort.Slice(o1.Parameters, func(i, j int) bool {
		p1 := o1.Parameters[i]
		p2 := o1.Parameters[j]
		return p1.Name < p2.Name || ((p1.Name == p2.Name) && p1.Value < p2.Value)
	})

	sort.Slice(o2.Parameters, func(i, j int) bool {
		p1 := o2.Parameters[i]
		p2 := o2.Parameters[j]
		return p1.Name < p2.Name || ((p1.Name == p2.Name) && p1.Value < p2.Value)
	})

	return reflect.DeepEqual(o1, o2)
}

func suppressEquivalentMatchTargetDiffs(_, oldString, newString string, _ *schema.ResourceData) bool {

	return compareMatchTargetsJSON(oldString, newString)
}

func suppressEquivalentJSONDiffsConditionException(_, oldString, newString string, _ *schema.ResourceData) bool {
	return compareConditionExceptionJSON(oldString, newString)
}

func compareConditionExceptionJSON(oldString, newString string) bool {
	var oldJSON, newJSON appsec.RuleConditionException
	if oldString == newString {
		return true
	}
	if err := json.Unmarshal([]byte(oldString), &oldJSON); err != nil {
		return false
	}
	if err := json.Unmarshal([]byte(newString), &newJSON); err != nil {
		return false
	}
	diff := compareConditionException(oldJSON, newJSON)
	return diff
}

func compareConditionException(oldValue, newValue appsec.RuleConditionException) bool {

	return reflect.DeepEqual(oldValue, newValue)
}

func suppressEquivalentMalwarePolicyJSONDiffs(_, oldString, newString string, _ *schema.ResourceData) bool {
	return compareMalwarePolicyJSON(oldString, newString)
}

func compareMalwarePolicyJSON(oldJSON, newJSON string) bool {
	if oldJSON == newJSON {
		return true
	}

	var oldPolicy, newPolicy appsec.MalwarePolicyBody
	if err := json.Unmarshal([]byte(oldJSON), &oldPolicy); err != nil {
		return false
	}
	if err := json.Unmarshal([]byte(newJSON), &newPolicy); err != nil {
		return false
	}
	return compareMalwarePolicy(oldPolicy, newPolicy)
}

func compareMalwarePolicy(oldPolicy, newPolicy appsec.MalwarePolicyBody) bool {
	return reflect.DeepEqual(oldPolicy, newPolicy)
}

func compareMatchTargetsJSON(oldString, newString string) bool {
	var oldJSON, newJSON appsec.CreateMatchTargetResponse
	if oldString == newString {
		return true
	}
	if err := json.Unmarshal([]byte(oldString), &oldJSON); err != nil {
		return false
	}
	if err := json.Unmarshal([]byte(newString), &newJSON); err != nil {
		return false
	}
	diff := compareMatchTargets(&oldJSON, &newJSON)
	return diff
}

func compareMatchTargets(oldTarget, newTarget *appsec.CreateMatchTargetResponse) bool {
	if len(oldTarget.FilePaths) != len(newTarget.FilePaths) ||
		len(oldTarget.FileExtensions) != len(newTarget.FileExtensions) ||
		len(oldTarget.Hostnames) != len(newTarget.Hostnames) ||
		len(oldTarget.BypassNetworkLists) != len(newTarget.BypassNetworkLists) ||
		len(oldTarget.Apis) != len(newTarget.Apis) {
		return false
	}

	sort.Strings(oldTarget.FilePaths)
	sort.Strings(newTarget.FilePaths)

	sort.Strings(oldTarget.FileExtensions)
	sort.Strings(newTarget.FileExtensions)

	sort.Strings(oldTarget.Hostnames)
	sort.Strings(newTarget.Hostnames)

	newTarget.TargetID = 0
	oldTarget.TargetID = 0

	newTarget.Sequence = 0
	oldTarget.Sequence = 0

	sort.Slice(oldTarget.Apis, func(i, j int) bool {
		p1 := oldTarget.Apis[i]
		p2 := oldTarget.Apis[j]
		return p1.ID < p2.ID || ((p1.ID == p2.ID) && p1.Name < p2.Name)
	})

	sort.Slice(newTarget.Apis, func(i, j int) bool {
		p1 := newTarget.Apis[i]
		p2 := newTarget.Apis[j]
		return p1.ID < p2.ID || ((p1.ID == p2.ID) && p1.Name < p2.Name)
	})

	sort.Slice(oldTarget.BypassNetworkLists, func(i, j int) bool {
		p1 := oldTarget.BypassNetworkLists[i]
		p2 := oldTarget.BypassNetworkLists[j]
		return p1.ID < p2.ID || ((p1.ID == p2.ID) && p1.Name < p2.Name)
	})

	sort.Slice(newTarget.BypassNetworkLists, func(i, j int) bool {
		p1 := newTarget.BypassNetworkLists[i]
		p2 := newTarget.BypassNetworkLists[j]
		return p1.ID < p2.ID || ((p1.ID == p2.ID) && p1.Name < p2.Name)
	})

	return reflect.DeepEqual(oldTarget, newTarget)
}

// suppress the diff if ukraine_geo_control_action is not passed in terraform config
func suppressDiffUkraineGeoControlAction(_, _, _ string, d *schema.ResourceData) bool {
	key := "ukraine_geo_control_action"

	oldValue, newValue := d.GetChange(key)
	oldUkraineGeoControlAction, ok := oldValue.(string)
	if !ok {
		logger.Get("APPSEC", "diffSuppressRules").Error(fmt.Sprintf("cannot parse ukraine_geo_control_action state properly for old value %v", oldValue))
		return true
	}
	newUkraineGeoControlAction, ok := newValue.(string)
	if !ok {
		logger.Get("APPSEC", "diffSuppressRules").Error(fmt.Sprintf("cannot parse ukraine_geo_control_action state properly for new value %v", oldValue))
		return true
	}
	return oldUkraineGeoControlAction == newUkraineGeoControlAction
}
