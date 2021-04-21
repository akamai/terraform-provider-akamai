package appsec

import (
	"bytes"
	"encoding/json"
	"reflect"
	"sort"

	"github.com/akamai/AkamaiOPEN-edgegrid-golang/v2/pkg/appsec"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/schema"
)

func suppressEquivalentJsonDiffsGeneric(k, old, new string, d *schema.ResourceData) bool {
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

func suppressCustomDenyJsonDiffs(k, old, new string, d *schema.ResourceData) bool {
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

func suppressEquivalentMatchTargetDiffs(k, old, new string, d *schema.ResourceData) bool {

	return compareMatchTargetsJSON(old, new)
}

func suppressEquivalentJSONDiffsConditionException(k, old, new string, d *schema.ResourceData) bool {
	return compareConditionExceptionJSON(old, new)

}

func compareConditionExceptionJSON(old, new string) bool {
	var oldJSON, newJSON appsec.UpdateRuleConditionExceptionResponse
	if old == new {
		return true
	}
	if err := json.Unmarshal([]byte(old), &oldJSON); err != nil {
		return false
	}
	if err := json.Unmarshal([]byte(new), &newJSON); err != nil {
		return false
	}
	diff := compareConditionException(&oldJSON, &newJSON)
	return diff
}

func compareConditionException(old, new *appsec.UpdateRuleConditionExceptionResponse) bool {
	if len(old.Conditions) != len(new.Conditions) ||

		len(old.Exception.HeaderCookieOrParamValues) != len(new.Exception.HeaderCookieOrParamValues) {
		return false
	}

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
