package iam

import (
	"testing"

	"github.com/akamai/terraform-provider-akamai/v2/pkg/test"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/schema"
	"github.com/stretchr/testify/assert"
)

func TestDSGroup(t *testing.T) {
	t.Run("groups can nest 25 levels deep", func(t *testing.T) {
		prov := provider{}

		assert.Equal(t, 25, GroupsNestingDepth(prov.dsGroups()), "incorrect nesting depth")
	})

	test.TODO(t, "need work")
}

// counts the nesting depth of the groups in the groups resource schema
func GroupsNestingDepth(res *schema.Resource) int {

	for attr, schem := range res.Schema {
		if attr == "sub_groups" || attr == "groups" {
			next := schem.Elem.(*schema.Resource)
			return 1 + GroupsNestingDepth(next)
		}
	}

	return 0
}
