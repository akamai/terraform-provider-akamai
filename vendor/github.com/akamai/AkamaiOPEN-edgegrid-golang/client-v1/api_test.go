package client

import (
	"testing"

	"github.com/akamai/AkamaiOPEN-edgegrid-golang/jsonhooks-v1"
	"github.com/stretchr/testify/assert"
)

type Test struct {
	Resource
	Foo string `json:"foo"`
}

func (test *Test) PreMarshalJSON() error {
	test.Foo = "bat"

	return nil
}

func TestResourceUnmarshal(t *testing.T) {
	body := []byte(`{"foo":"bar"}`)

	test := &Test{}
	err := jsonhooks.Unmarshal(body, test)

	assert.NoError(t, err)
	assert.True(t, <-test.Complete)
}

func TestResourceMarshal(t *testing.T) {
	test := &Test{Foo: "bar"}

	body, err := jsonhooks.Marshal(test)

	assert.NoError(t, err)
	assert.Equal(t, []byte(`{"foo":"bat"}`), body)
}
