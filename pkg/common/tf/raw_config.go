package tf

import (
	"strconv"
	"strings"

	"github.com/hashicorp/go-cty/cty"
	"github.com/hashicorp/go-cty/cty/gocty"
)

// RawConfigGetter is used to retrieve raw config.
type RawConfigGetter interface {
	GetRawConfig() cty.Value
}

// RawConfig is used to query attributes from the raw config.
type RawConfig struct {
	data cty.Value
}

// NewRawConfig creates a new RawConfig which uses the raw config retrieved from data.
func NewRawConfig(g RawConfigGetter) *RawConfig {
	return &RawConfig{data: g.GetRawConfig()}
}

// GetOk returns the data for the given key and whether or not the key was set to a non-zero value.
func (rc RawConfig) GetOk(key string) (any, bool) {
	path := rc.buildPath(key)

	val, err := path.Apply(rc.data)
	if err != nil {
		return nil, false
	}

	return rc.transform(val)
}

func (rc RawConfig) buildPath(key string) cty.Path {
	parts := strings.Split(key, ".")
	path := make(cty.Path, 0)

	for _, part := range parts {
		asInt, err := strconv.Atoi(part)
		if err == nil {
			path = append(path, cty.IndexIntPath(asInt)...)
		} else {
			path = append(path, cty.GetAttrPath(part)...)
		}
	}

	return path
}

func (rc RawConfig) transform(val cty.Value) (any, bool) {
	if val == cty.NilVal || val.IsNull() {
		return nil, true
	}

	switch val.Type() {
	case cty.Bool:
		var v bool
		if err := gocty.FromCtyValue(val, &v); err != nil {
			panic(err) // if Type() == cty.Bool, this should never happen
		}
		return v, true
	case cty.Number:
		fv := val.AsBigFloat()
		if fv.IsInt() {
			i64, _ := fv.Int64()
			return i64, true
		}
		rfv, _ := fv.Float64()
		return rfv, true
	case cty.String:
		return val.AsString(), true
	}

	switch {
	case val.Type().IsSetType(), val.Type().IsListType():
		return rc.transformSlice(val)
	case val.Type().IsMapType():
		return rc.transformMap(val)
	case val.Type().IsObjectType():
		if val.Type().Equals(cty.EmptyObject) {
			return make(map[string]any, 0), true
		}
		return rc.transformMap(val)
	}

	return nil, false
}

func (rc RawConfig) transformMap(val cty.Value) (any, bool) {
	vals := val.AsValueMap()
	if len(vals) == 0 {
		return make(map[string]any, 0), true
	}
	mp := make(map[string]any, 0)
	for k, v := range vals {
		if v, ok := rc.transform(v); ok {
			mp[k] = v
		} else {
			return nil, false
		}
	}
	return mp, true
}

func (rc RawConfig) transformSlice(val cty.Value) (any, bool) {
	slc := make([]any, 0, val.LengthInt())
	for it := val.ElementIterator(); it.Next(); {
		_, v := it.Element()
		if v, ok := rc.transform(v); ok {
			slc = append(slc, v)
		} else {
			return nil, ok
		}
	}
	if val.LengthInt() == 0 && val.Type().ListElementType().IsObjectType() {
		return nil, true
	}

	return slc, true
}
