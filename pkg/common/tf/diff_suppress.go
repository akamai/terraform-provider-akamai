package tf

import (
	"strings"

	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/schema"
)

// DiffSuppressAny aggregates all given schema.SchemaDiffSuppressFunc into one.
// Aggregated function returns true if any of the functions returns true.
func DiffSuppressAny(fncs ...schema.SchemaDiffSuppressFunc) schema.SchemaDiffSuppressFunc {
	return func(k, oldValue, newValue string, d *schema.ResourceData) bool {
		for _, fn := range fncs {
			if fn(k, oldValue, newValue, d) {
				return true
			}
		}
		return false
	}
}

// FieldPrefixSuppress returns schema.SchemaDiffSuppressFunc which verifies whether values are equal disregarding given prefix
func FieldPrefixSuppress(prefix string) schema.SchemaDiffSuppressFunc {
	return func(_, o, n string, _ *schema.ResourceData) bool {
		trimPrefixFromOld := strings.TrimPrefix(o, prefix)
		trimPrefixFromNew := strings.TrimPrefix(n, prefix)

		return trimPrefixFromOld == trimPrefixFromNew
	}
}
