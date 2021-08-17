package appsec

import (
	"bytes"
	"encoding/json"
	"log"
	"reflect"
	"sort"
	"strings"

	"github.com/akamai/AkamaiOPEN-edgegrid-golang/v2/pkg/appsec"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/schema"
)

func suppressEquivalentJSONDiffsGeneric(_, old, new string, _ *schema.ResourceData) bool {
	var ob, nb bytes.Buffer
	if err := json.Compact(&ob, []byte(old)); err != nil {
		return false
	}

	if err := json.Compact(&nb, []byte(new)); err != nil {
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

func suppressEquivalentReputationProfileDiffs(_, old, new string, _ *schema.ResourceData) bool {
	var rpOld, rpNew appsec.CreateReputationProfileResponse

	if err := json.Unmarshal([]byte(old), &rpOld); err != nil {
		log.Printf("unable to unmarshal old reputation profile: %s", err)
	}
	if err := json.Unmarshal([]byte(new), &rpNew); err != nil {
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
	for _, acOld := range cOld.AtomicConditions {
		found := false
		for _, acNew := range cNew.AtomicConditions {
			if acOld.ClassName == acNew.ClassName {
				found = true
				if acOld.CheckIps != acNew.CheckIps && acNew.CheckIps != "" {
					return false
				}
				if acOld.NameCase != acNew.NameCase && !(acOld.NameCase && !acNew.NameCase) {
					return false
				}
				if acOld.NameWildcard != acNew.NameWildcard && !(acOld.NameWildcard && !acNew.NameWildcard) {
					return false
				}
				if acOld.ValueCase != acNew.ValueCase && !(acOld.ValueCase && !acNew.ValueCase) {
					return false
				}
				if acOld.ValueWildcard != acNew.ValueWildcard && !(acOld.ValueWildcard && !acNew.ValueWildcard) {
					return false
				}
				if acOld.PositiveMatch != acNew.PositiveMatch && !(acOld.PositiveMatch && !acNew.PositiveMatch) {
					// only 'true' is supported for this case
					if acOld.ClassName != "HostCondition" {
						return false
					}
				}
				if !suppressAtomicConditionSliceDiffs(acOld.Value, acNew.Value) {
					return false
				}
				if !suppressAtomicConditionSliceDiffs(acOld.Host, acNew.Host) {
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

func suppressAtomicConditionSliceDiffs(old, new []string) bool {
	if len(old) != len(new) {
		return false
	}
	for _, ov := range old {
		found := false
		for _, nv := range new {
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

func suppressEquivalentLoggingSettingsDiffs(_, old, new string, _ *schema.ResourceData) bool {
	var oldJSON, newJSON appsec.UpdateAdvancedSettingsLoggingResponse
	if old == new {
		return true
	}
	if err := json.Unmarshal([]byte(old), &oldJSON); err != nil {
		return false
	}
	if err := json.Unmarshal([]byte(new), &newJSON); err != nil {
		return false
	}
	diff := compareLoggingSettings(&oldJSON, &newJSON)
	return diff
}

func compareLoggingSettings(old, new *appsec.UpdateAdvancedSettingsLoggingResponse) bool {
	if old.Override != new.Override ||
		old.AllowSampling != new.AllowSampling ||
		old.Cookies.Type != new.Cookies.Type ||
		old.CustomHeaders.Type != new.CustomHeaders.Type ||
		old.StandardHeaders.Type != new.StandardHeaders.Type {
		return false
	}

	sort.Strings(old.Cookies.Values)
	sort.Strings(new.Cookies.Values)
	sort.Strings(old.CustomHeaders.Values)
	sort.Strings(new.CustomHeaders.Values)
	sort.Strings(old.StandardHeaders.Values)
	sort.Strings(new.StandardHeaders.Values)

	return reflect.DeepEqual(old, new)
}

func suppressCustomDenyJSONDiffs(_, old, new string, _ *schema.ResourceData) bool {
	var ob, nb bytes.Buffer
	if err := json.Compact(&ob, []byte(old)); err != nil {
		return false
	}

	if err := json.Compact(&nb, []byte(new)); err != nil {
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

func suppressEquivalentMatchTargetDiffs(_, old, new string, _ *schema.ResourceData) bool {

	return compareMatchTargetsJSON(old, new)
}

func suppressEquivalentJSONDiffsConditionException(_, old, new string, _ *schema.ResourceData) bool {
	return compareConditionExceptionJSON(old, new)

}

func compareConditionExceptionJSON(old, new string) bool {
	var oldJSON, newJSON appsec.RuleConditionException
	if old == new {
		return true
	}
	if err := json.Unmarshal([]byte(old), &oldJSON); err != nil {
		return false
	}
	if err := json.Unmarshal([]byte(new), &newJSON); err != nil {
		return false
	}
	diff := compareConditionException(oldJSON, newJSON)
	return diff
}

func compareConditionException(old, new appsec.RuleConditionException) bool {

	return reflect.DeepEqual(old, new)
}

func compareMatchTargetsJSON(old, new string) bool {
	var oldJSON, newJSON appsec.CreateMatchTargetResponse
	if old == new {
		return true
	}
	if err := json.Unmarshal([]byte(old), &oldJSON); err != nil {
		return false
	}
	if err := json.Unmarshal([]byte(new), &newJSON); err != nil {
		return false
	}
	diff := compareMatchTargets(&oldJSON, &newJSON)
	return diff
}

func compareMatchTargets(old, new *appsec.CreateMatchTargetResponse) bool {
	if len(old.FilePaths) != len(new.FilePaths) ||
		len(old.FileExtensions) != len(new.FileExtensions) ||
		len(old.Hostnames) != len(new.Hostnames) ||
		len(old.BypassNetworkLists) != len(new.BypassNetworkLists) ||
		len(old.Apis) != len(new.Apis) {
		return false
	}

	sort.Strings(old.FilePaths)
	sort.Strings(new.FilePaths)

	sort.Strings(old.FileExtensions)
	sort.Strings(new.FileExtensions)

	sort.Strings(old.Hostnames)
	sort.Strings(new.Hostnames)

	new.TargetID = 0
	old.TargetID = 0

	new.Sequence = 0
	old.Sequence = 0

	sort.Slice(old.Apis, func(i, j int) bool {
		p1 := old.Apis[i]
		p2 := old.Apis[j]
		return p1.ID < p2.ID || ((p1.ID == p2.ID) && p1.Name < p2.Name)
	})

	sort.Slice(new.Apis, func(i, j int) bool {
		p1 := new.Apis[i]
		p2 := new.Apis[j]
		return p1.ID < p2.ID || ((p1.ID == p2.ID) && p1.Name < p2.Name)
	})
	return reflect.DeepEqual(old, new)
}
