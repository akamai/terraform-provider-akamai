package jsonhooks

import (
	"encoding/json"
	"strings"
	"testing"

	"github.com/stretchr/testify/assert"
)

type Optionals struct {
	Sr string `json:"sr"`
	So string `json:"so,omitempty"`
	Sw string `json:"-"`

	Ir int `json:"omitempty"` // actually named omitempty, not an option
	Io int `json:"io,omitempty"`

	Slr []string `json:"slr,random"`
	Slo []string `json:"slo,omitempty"`

	Mr map[string]interface{} `json:"mr"`
	Mo map[string]interface{} `json:",omitempty"`

	Fr float64 `json:"fr"`
	Fo float64 `json:"fo,omitempty"`

	Br bool `json:"br"`
	Bo bool `json:"bo,omitempty"`

	Ur uint `json:"ur"`
	Uo uint `json:"uo,omitempty"`

	Str struct{} `json:"str"`
	Sto struct{} `json:"sto,omitempty"`
}

type MixedTypes struct {
	I  int     `json:"I"`
	B  bool    `json:"B"`
	F  float64 `json:"F"`
	S  string  `json:"S"`
	St struct {
		Foo int    `json:"Foo"`
		Bar string `json:"Bar"`
	} `json:"St"`
	A []string `json:"A"`
}

type WithHooks MixedTypes

func (hooks *WithHooks) PreMarshalJSON() error {
	hooks.I *= 1000
	hooks.B = !hooks.B
	hooks.F *= 1.1
	hooks.S = strings.ToUpper(hooks.S)
	hooks.St.Foo *= 2000
	hooks.St.Bar = strings.ToUpper(hooks.St.Bar)
	for key, val := range hooks.A {
		hooks.A[key] = strings.ToUpper(val)
	}

	return nil
}

func (hooks *WithHooks) PostUnmarshalJSON() error {
	hooks.I /= 1000
	hooks.B = !hooks.B
	hooks.F /= 1.1
	hooks.S = strings.ToLower(hooks.S)
	hooks.St.Foo /= 2000
	hooks.St.Bar = strings.ToLower(hooks.St.Bar)
	for key, val := range hooks.A {
		hooks.A[key] = strings.ToLower(val)
	}

	return nil
}

func TestMarshalCompat(t *testing.T) {
	var o Optionals
	o.Sw = "something"
	o.Mr = map[string]interface{}{}
	o.Mo = map[string]interface{}{}

	expected, _ := json.Marshal(o)
	actual, err := Marshal(o)

	assert.NoError(t, err)
	assert.Equal(
		t,
		expected,
		actual,
	)
}

func TestUnmarshalCompat(t *testing.T) {
	data := `{
		"sr": "",
				"omitempty": 0,
				"slr": null,
				"mr": {},
		"fr": 0,
				"br": false,
				"ur": 0,
				"str": {},
		"sto": {}
	}`

	expected := &Optionals{}
	_ = json.Unmarshal([]byte(data), expected)

	actual := &Optionals{}
	err := Unmarshal([]byte(data), actual)

	assert.NoError(t, err)
	assert.Equal(
		t,
		expected,
		actual,
	)
}

func TestImplementsPreJSONMarshaler(t *testing.T) {
	noHooks := &MixedTypes{}
	withHooks := &WithHooks{}

	assert.False(t, ImplementsPreJSONMarshaler(noHooks))
	assert.True(t, ImplementsPreJSONMarshaler(withHooks))
}

func TestImplementsPostJSONUnmarshaler(t *testing.T) {
	noHooks := &MixedTypes{}
	withHooks := &WithHooks{}

	assert.False(t, ImplementsPostJSONUnmarshaler(noHooks))
	assert.True(t, ImplementsPostJSONUnmarshaler(withHooks))
}

func TestPreJSONMarshal(t *testing.T) {
	noHooks := &MixedTypes{
		I: 1,
		B: true,
		F: 1.0,
		S: "testing",
		St: struct {
			Foo int    `json:"Foo"`
			Bar string `json:"Bar"`
		}{2, "test"},
		A: []string{"one", "two", "three"},
	}

	withHooks := (*WithHooks)(noHooks)

	expected, _ := json.Marshal(&MixedTypes{
		I: 1 * 1000,
		B: !true,
		F: 1.0 * 1.1,
		S: "TESTING",
		St: struct {
			Foo int    `json:"Foo"`
			Bar string `json:"Bar"`
		}{2 * 2000, "TEST"},
		A: []string{"ONE", "TWO", "THREE"},
	})

	gojson, _ := json.Marshal(withHooks)
	actualNoHooks, _ := Marshal(noHooks)
	actualWithHooks, err := Marshal(withHooks)

	assert.NoError(t, err)
	assert.Equal(t, string(gojson), string(actualNoHooks))
	assert.NotEqual(t, string(gojson), string(actualWithHooks))
	assert.NotEqual(t, string(actualNoHooks), string(actualWithHooks))
	assert.Equal(t, string(expected), string(actualWithHooks))
}

func TestPostJSONUnmarshal(t *testing.T) {
	expected := &WithHooks{
		I: 1,
		B: true,
		F: 1.0,
		S: "testing",
		St: struct {
			Foo int    `json:"Foo"`
			Bar string `json:"Bar"`
		}{2, "test"},
		A: []string{"one", "two", "three"},
	}

	data := []byte(`{"I":1000,"B":false,"F":1.1,"S":"TESTING","St":{"Foo":4000,"Bar":"TEST"},"A":["ONE","TWO","THREE"]}`)

	gojson := &MixedTypes{}
	_ = json.Unmarshal(
		data,
		gojson,
	)

	withHooks := &WithHooks{}
	err := Unmarshal(data, withHooks)
	assert.NoError(t, err)

	withoutHooks := &MixedTypes{}
	err = Unmarshal(data, withoutHooks)

	assert.NoError(t, err)
	assert.NotEqual(t, gojson, withHooks)
	assert.NotEqual(t, expected, withoutHooks)
	assert.Equal(t, expected, withHooks)
}
