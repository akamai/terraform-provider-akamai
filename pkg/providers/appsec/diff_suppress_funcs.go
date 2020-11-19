package appsec

import (
	"encoding/json"
	"reflect"
	"sort"

	"github.com/akamai/AkamaiOPEN-edgegrid-golang/v2/pkg/appsec"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/schema"
)

//suppressJsonProvided to handle when json supplied vs HCL values
func suppressJsonProvidedSimple(_, old, new string, d *schema.ResourceData) bool {

	json := d.Get("json").(string)
	if json != "" {
		if old == "" && new == "" {
			return true
		}
		return true
	}

	return false
}

func suppressJsonProvided(k, old, new string, d *schema.ResourceData) bool {
	json := d.Get("json").(string)
	if json != "" {

		if k == "json" {

			return compareMatchTargetsJSON(old, new)

		}

		if k == "match_target_id" {
			return true
		}

		if old == "" && new == "" {

			return true
		}

		if old == new {
			return true
		}

		if new == "" {
			return false
		}
	} else {
		if k == "json" {
			return true
		}
	}

	if old == new {
		return true
	}
	return false
}

func suppressEquivalentJSONDiffs(k, old, new string, d *schema.ResourceData) bool {
	if new == "" {
		jsonfromschema, err := matchTargetAsJSONDString(d)
		if err == nil {
			return compareMatchTargetsJSON(old, jsonfromschema)
		}
	}

	return compareMatchTargetsJSON(old, new)
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

	new.EffectiveSecurityControls = old.EffectiveSecurityControls

	new.TargetID = 0
	old.TargetID = 0

	new.Sequence = 0
	old.Sequence = 0

	return reflect.DeepEqual(old, new)
}
