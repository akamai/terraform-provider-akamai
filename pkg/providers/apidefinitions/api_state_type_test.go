package apidefinitions

import (
	"testing"

	v0 "github.com/akamai/AkamaiOPEN-edgegrid-golang/v11/pkg/apidefinitions/v0"
	"github.com/akamai/terraform-provider-akamai/v8/pkg/common/ptr"
	"github.com/stretchr/testify/assert"
	orderedmap "github.com/wk8/go-ordered-map/v2"
)

func TestCheckSemanticEquality_BasePath(t *testing.T) {
	var before = base()
	var after = base()

	before.BasePath = ptr.To("")

	assert.Equal(t, "", checkSemanticEquality(before, after))
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

	assert.Equal(t, "", checkSemanticEquality(before, after))
}

func TestCheckSemanticEquality_Hostnames(t *testing.T) {
	var before = base()
	var after = base()

	before.Hostnames = []string{"host2.com", "host1.com"}
	after.Hostnames = []string{"host1.com", "host2.com"}

	assert.Equal(t, "", checkSemanticEquality(before, after))
}

func TestCheckSemanticEquality_Tags(t *testing.T) {
	var before = base()
	var after = base()

	before.Tags = []string{"tag2", "tag1"}
	after.Tags = []string{"tag1", "tag2"}

	assert.Equal(t, "", checkSemanticEquality(before, after))
}

func TestCheckSemanticEquality_ResourcesAndMethodsSame(t *testing.T) {
	var before = base()
	var after = base()

	body := orderedmap.New[string, v0.Property]()
	body.Set("json", v0.Property{Name: "json body"})

	getMethod := &v0.Method{
		RequestBody: body,
	}

	before.Resources = orderedmap.New[string, v0.Resource]()
	before.Resources.Set("/one", v0.Resource{Name: "/one", Get: getMethod})

	after.Resources = orderedmap.New[string, v0.Resource]()
	after.Resources.Set("/one", v0.Resource{Name: "/one", Get: getMethod})

	assert.Equal(t, "", checkSemanticEquality(before, after))
}

func TestCheckSemanticEquality_ResourcesDifferentOrder(t *testing.T) {
	var before = base()
	var after = base()

	before.Resources = orderedmap.New[string, v0.Resource]()
	before.Resources.Set("/one", v0.Resource{Name: "/one"})
	before.Resources.Set("/two", v0.Resource{Name: "/two"})

	after.Resources = orderedmap.New[string, v0.Resource]()
	after.Resources.Set("/two", v0.Resource{Name: "/two"})
	after.Resources.Set("/one", v0.Resource{Name: "/one"})

	assert.Equal(t, "", checkSemanticEquality(before, after))
}

func TestCheckSemanticEquality_ResourceAdded(t *testing.T) {
	var before = base()
	var after = base()

	before.Resources = orderedmap.New[string, v0.Resource]()
	before.Resources.Set("/one", v0.Resource{Name: "/one"})

	after.Resources = orderedmap.New[string, v0.Resource]()
	after.Resources.Set("/one", v0.Resource{Name: "/one"})
	after.Resources.Set("/two", v0.Resource{Name: "/two"})

	assert.Contains(t, checkSemanticEquality(before, after), "\"/two\": &{")
}

func TestCheckSemanticEquality_ResourceModified(t *testing.T) {
	var before = base()
	var after = base()

	before.Resources = orderedmap.New[string, v0.Resource]()
	before.Resources.Set("/one", v0.Resource{Name: "/one", Description: ptr.To("description one")})

	after.Resources = orderedmap.New[string, v0.Resource]()
	after.Resources.Set("/one", v0.Resource{Name: "/one", Description: ptr.To("description two")})

	assert.Contains(t, checkSemanticEquality(before, after), "Value: v0.Resource{Name: \"/one\",")
}

func base() v0.APIAttributes {
	return v0.APIAttributes{
		Name:      "Name",
		Hostnames: []string{"host1.com"},
	}
}
