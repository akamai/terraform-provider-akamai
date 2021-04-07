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

func suppressEquivalentJSONDiffs(k, old, new string, d *schema.ResourceData) bool {

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
		len(old.BypassNetworkLists) != len(new.BypassNetworkLists) {
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

	return reflect.DeepEqual(old, new)
}
