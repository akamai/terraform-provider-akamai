package apidefinitions

import (
	"testing"

	v0 "github.com/akamai/AkamaiOPEN-edgegrid-golang/v11/pkg/apidefinitions/v0"
	"github.com/akamai/terraform-provider-akamai/v8/pkg/common/ptr"
	"github.com/stretchr/testify/assert"
)

func TestCheckSemanticEquality_BasePath(t *testing.T) {
	var before = base()
	var after = base()

	before.BasePath = ptr.To("")

	assert.Equal(t, []string(nil), checkSemanticEquality(before, after))
}

func TestCheckSemanticEquality_ConsumeType(t *testing.T) {
	var before = base()
	var after = base()

	before.Constraints = &v0.Constraints{
		RequestBody: &v0.ConstraintsRequestBody{
			ConsumeType: []v0.ConsumeType{v0.ConsumeTypeXML, v0.ConsumeTypeUrlencoded, v0.ConsumeTypeJSON},
		},
	}

	after.Constraints = &v0.Constraints{
		RequestBody: &v0.ConstraintsRequestBody{
			ConsumeType: []v0.ConsumeType{v0.ConsumeTypeJSON, v0.ConsumeTypeXML, v0.ConsumeTypeUrlencoded},
		},
	}

	assert.Equal(t, []string(nil), checkSemanticEquality(before, after))
}

func TestCheckSemanticEquality_Hostnames(t *testing.T) {
	var before = base()
	var after = base()

	before.Hostnames = []string{"host2.com", "host1.com"}
	after.Hostnames = []string{"host1.com", "host2.com"}

	assert.Equal(t, []string(nil), checkSemanticEquality(before, after))
}

func TestCheckSemanticEquality_Tags(t *testing.T) {
	var before = base()
	var after = base()

	before.Tags = []string{"tag2", "tag1"}
	after.Tags = []string{"tag1", "tag2"}

	assert.Equal(t, []string(nil), checkSemanticEquality(before, after))
}

func base() v0.APIAttributes {
	return v0.APIAttributes{
		Name:      "Name",
		Hostnames: []string{"host1.com"},
	}
}
