package tools

import "github.com/hashicorp/terraform-plugin-sdk/v2/helper/schema"

// HashFromChallengesMap takes Challenges map as an argument. Calculates and returns hash of a domain
func HashFromChallengesMap(v interface{}) int {
	m, ok := v.(map[string]interface{})
	if !ok {
		return 0
	}
	val, ok := m["domain"]
	if !ok {
		return 0
	}

	return schema.HashString(val)
}
