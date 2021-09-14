package cloudlets

import (
	"crypto/sha1"
	"encoding/hex"
	"io"

	"github.com/akamai/AkamaiOPEN-edgegrid-golang/v2/pkg/cloudlets"
)

func getMatchRulesHashID(matchRules *cloudlets.MatchRules) (string, error) {
	id := "id"
	for _, rule := range *matchRules {
		switch r := rule.(type) {
		case cloudlets.MatchRuleER:
			id = id + ":" + r.Name
		}
	}
	h := sha1.New()
	_, err := io.WriteString(h, id)
	if err != nil {
		return "", err
	}
	hashID := hex.EncodeToString(h.Sum(nil))
	return hashID, nil
}

func getStringValue(matchRuleMap map[string]interface{}, name string) string {
	if value, ok := matchRuleMap[name]; ok {
		return value.(string)
	}
	return ""
}

func getIntValue(matchRuleMap map[string]interface{}, name string) int {
	if value, ok := matchRuleMap[name]; ok {
		return value.(int)
	}
	return 0
}

// this will not work on 32bit platform if the value is bigger than max for int32
func getInt64Value(matchRuleMap map[string]interface{}, name string) int64 {
	if value, ok := matchRuleMap[name]; ok {
		return int64(value.(int))
	}
	return 0
}

func getBoolValue(matchRuleMap map[string]interface{}, name string) bool {
	if value, ok := matchRuleMap[name]; ok {
		return value.(bool)
	}
	return false
}
