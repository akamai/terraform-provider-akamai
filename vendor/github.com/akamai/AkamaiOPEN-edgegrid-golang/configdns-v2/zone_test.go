package dnsv2

import (
	"testing"
	"fmt"
	"github.com/akamai/AkamaiOPEN-edgegrid-golang/jsonhooks-v1"
	"github.com/stretchr/testify/assert"
	"gopkg.in/h2non/gock.v1"
)

func TestZone_JSON(t *testing.T) {
	responseBody := []byte(`{
    "zone": "example.com",
    "type": "PRIMARY",
    "comment": "This is a test zone",
    "signAndServe": false
}`)

	zonecreate := ZoneCreate{Zone: "example.com", Type: "PRIMARY", Masters: []string{""}, Comment: "This is a test zone", SignAndServe: false}
	zone := NewZone(zonecreate)
	err := jsonhooks.Unmarshal(responseBody, zone)
	assert.NoError(t, err)
	assert.Equal(t, zone.Zone, "example.com")
	assert.Equal(t, zone.Type, "PRIMARY")
	assert.Equal(t, zone.Comment, "This is a test zone")
	assert.Equal(t, zone.SignAndServe, false)
}

func TestGetZoneNames(t *testing.T) {

	dnsTestZone := "testzone.com"

	defer gock.Off()

	mock := gock.New(fmt.Sprintf("https://akaa-baseurl-xxxxxxxxxxx-xxxxxxxxxxxxx.luna.akamaiapis.net/config-dns/v2/zones/%s/names", dnsTestZone))
	mock.
		Get(fmt.Sprintf("/config-dns/v2/zones/%s/names", dnsTestZone)).
		HeaderPresent("Authorization").
		Reply(200).
		SetHeader("Content-Type", "application/json;charset=UTF-8").
		BodyString(fmt.Sprintf(`{
                        "names":["test1.testzone.com","test2.testzone.com"]
                }`))

	Init(config)
	nameList, err := GetZoneNames(dnsTestZone)
	assert.NoError(t, err)
	assert.Equal(t, assert.IsType(t, &ZoneNamesResponse{}, nameList), true)
	assert.Equal(t, len(nameList.Names), 2)

}

func TestGetZoneNameTypes(t *testing.T) {

        dnsTestZone := "testzone.com"
	dnsTestRecordName := "test.testzone.com"

        defer gock.Off()

        mock := gock.New(fmt.Sprintf("https://akaa-baseurl-xxxxxxxxxxx-xxxxxxxxxxxxx.luna.akamaiapis.net/config-dns/v2/zones/%s/names/%s/types", dnsTestZone, dnsTestRecordName))
        mock.
                Get(fmt.Sprintf("/config-dns/v2/zones/%s/names/%s/types", dnsTestZone, dnsTestRecordName)).
                HeaderPresent("Authorization").
                Reply(200).
                SetHeader("Content-Type", "application/json;charset=UTF-8").
                BodyString(fmt.Sprintf(`{
                        "types":["CNAME", "AKAMAICDN"]
                }`))

        Init(config)
        typeList, err := GetZoneNameTypes(dnsTestRecordName, dnsTestZone)
        assert.NoError(t, err)
        assert.Equal(t, assert.IsType(t, &ZoneNameTypesResponse{}, typeList), true)
        assert.Equal(t, len(typeList.Types), 2)

}


