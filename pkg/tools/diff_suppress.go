package tools

import (
	"strings"

	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/schema"
)

// FieldPrefixSuppress returns schema.SchemaDiffSuppressFunc which verifies whether values are equal disregarding given prefix
func FieldPrefixSuppress(prefix string) schema.SchemaDiffSuppressFunc {
	return func(_, old, new string, _ *schema.ResourceData) bool {
		trimPrefixFromOld := strings.TrimPrefix(old, prefix)
		trimPrefixFromNew := strings.TrimPrefix(new, prefix)

		return trimPrefixFromOld == trimPrefixFromNew
	}
}
