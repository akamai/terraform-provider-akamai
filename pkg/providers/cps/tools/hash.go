package tools

import "github.com/hashicorp/terraform-plugin-sdk/v2/helper/schema"

func HashFromChallengesMap(v interface{}) int {
	m, ok := v.(map[string]interface{})

	if !ok {
		return 0
	}

	return schema.HashString(m["domain"])
}
