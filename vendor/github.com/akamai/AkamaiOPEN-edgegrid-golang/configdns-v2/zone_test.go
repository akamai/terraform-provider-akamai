package dnsv2

import (
	"testing"

	"github.com/akamai/AkamaiOPEN-edgegrid-golang/jsonhooks-v1"
	"github.com/stretchr/testify/assert"
)

func TestZone_JSON(t *testing.T) {
	responseBody := []byte(`{
    "zone": "example.com",
    "type": "PRIMARY",
    "comment": "This is a test zone",
    "signAndServe": false
}`)

	zonecreate := ZoneCreate{"example.com", "PRIMARY", []string{""}, "This is a test zone", false}
	zone := NewZone(zonecreate)
	err := jsonhooks.Unmarshal(responseBody, zone)
	assert.NoError(t, err)
	assert.Equal(t, zone.Zone, "example.com")
	assert.Equal(t, zone.Type, "PRIMARY")
	assert.Equal(t, zone.Comment, "This is a test zone")
	assert.Equal(t, zone.SignAndServe, false)
}
