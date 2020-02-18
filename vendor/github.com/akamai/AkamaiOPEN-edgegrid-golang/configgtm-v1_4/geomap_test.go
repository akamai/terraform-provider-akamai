package configgtm

import (
	"testing"

	"github.com/akamai/AkamaiOPEN-edgegrid-golang/jsonhooks-v1"

	"github.com/stretchr/testify/assert"
	"gopkg.in/h2non/gock.v1"
)

var GtmTestGeoMap = "testGeoMap"

func instantiateGeoMap() *GeoMap {

	geoMap := NewGeoMap(GtmTestGeoMap)
	geoMapData := []byte(`{
                        "assignments": [ {
                                "countries": [ "GB", "IE" ],
                                "datacenterId": 3133,
                                "nickname": "UK and Ireland users"
                        } ],
                        "defaultDatacenter": {
                                "datacenterId": 5400,
                                "nickname": "Default Mapping"
                        },
                        "links": [ {
                                "href": "/config-gtm/v1/domains/example.akadns.net/geographic-maps/UK%20Delivery",
                                "rel": "self"
                        } ],
                        "name": "testGeoMap"
              }`)
	jsonhooks.Unmarshal(geoMapData, geoMap)

	return geoMap

}

// Verify ListGeoMap. Name hardcoded. Should pass, e.g. no API errors and resource returned
func TestListGeoMaps(t *testing.T) {

	defer gock.Off()

	mock := gock.New("https://akaa-baseurl-xxxxxxxxxxx-xxxxxxxxxxxxx.luna.akamaiapis.net/config-gtm/v1/domains/" + gtmTestDomain + "/geographic-maps")
	mock.
		Get("/config-gtm/v1/domains/"+gtmTestDomain+"/geographic-maps").
		HeaderPresent("Authorization").
		Reply(200).
		SetHeader("Content-Type", "application/vnd.config-gtm.v1.4+json;charset=UTF-8").
		BodyString(`{
                    "items" : [ {
                        "assignments": [ {
                                "countries": [ "GB", "IE" ],
                                "datacenterId": 3133,
                                "nickname": "UK and Ireland users"
                        } ],
                        "defaultDatacenter": {
                                "datacenterId": 5400,
                                "nickname": "Default Mapping"
                        },
                        "links": [ {
                                "href": "/config-gtm/v1/domains/example.akadns.net/geographic-maps/UK%20Delivery",
                                "rel": "self"
                        } ],
                        "name": "testGeoMap"
                    } ]
               }`)

	Init(config)

	testGeoMap, err := ListGeoMaps(gtmTestDomain)
	assert.NoError(t, err)
	assert.IsType(t, &GeoMap{}, testGeoMap[0])
	assert.Equal(t, GtmTestGeoMap, testGeoMap[0].Name)

}

// Verify GetGeoMap. Name hardcoded. Should pass, e.g. no API errors and resource returned
func TestGetGeoMap(t *testing.T) {

	defer gock.Off()

	mock := gock.New("https://akaa-baseurl-xxxxxxxxxxx-xxxxxxxxxxxxx.luna.akamaiapis.net/config-gtm/v1/domains/" + gtmTestDomain + "/geographic-maps/" + GtmTestGeoMap)
	mock.
		Get("/config-gtm/v1/domains/"+gtmTestDomain+"/geographic-maps/"+GtmTestGeoMap).
		HeaderPresent("Authorization").
		Reply(200).
		SetHeader("Content-Type", "application/vnd.config-gtm.v1.4+json;charset=UTF-8").
		BodyString(`{
                        "assignments": [ {
                                "countries": [ "GB", "IE" ],
                                "datacenterId": 3133,
                                "nickname": "UK and Ireland users"
                        } ],
                        "defaultDatacenter": {
                                "datacenterId": 5400,
                                "nickname": "Default Mapping"
                        },
                        "links": [ {
                                "href": "/config-gtm/v1/domains/example.akadns.net/geographic-maps/UK%20Delivery",
                                "rel": "self"
                        } ],
                        "name": "testGeoMap"
               }`)

	Init(config)

	testGeoMap, err := GetGeoMap(GtmTestGeoMap, gtmTestDomain)
	assert.NoError(t, err)
	assert.IsType(t, &GeoMap{}, testGeoMap)
	assert.Equal(t, GtmTestGeoMap, testGeoMap.Name)

}

// Verify failed case for GetGeoMap. Should pass, e.g. no API errors and domain not found
func TestGetBadGeoMap(t *testing.T) {

	badName := "somebadname"
	defer gock.Off()

	mock := gock.New("https://akaa-baseurl-xxxxxxxxxxx-xxxxxxxxxxxxx.luna.akamaiapis.net/config-gtm/v1/domains/" + gtmTestDomain + "/geographic-maps/" + badName)
	mock.
		Get("/config-gtm/v1/domains/"+gtmTestDomain+"/geographic-maps/"+badName).
		HeaderPresent("Authorization").
		Reply(404).
		SetHeader("Content-Type", "application/vnd.config-gtm.v1.4+json;charset=UTF-8").
		BodyString(`{
                }`)

	Init(config)

	_, err := GetGeoMap(badName, gtmTestDomain)
	assert.Error(t, err)

}

// Test Create GeoMap.
func TestCreateGeoMap(t *testing.T) {

	defer gock.Off()

	mock := gock.New("https://akaa-baseurl-xxxxxxxxxxx-xxxxxxxxxxxxx.luna.akamaiapis.net/config-gtm/v1/domains/" + gtmTestDomain + "/geographic-maps/" + GtmTestGeoMap)
	mock.
		Put("/config-gtm/v1/domains/"+gtmTestDomain+"/geographic-maps/"+GtmTestGeoMap).
		HeaderPresent("Authorization").
		Reply(200).
		SetHeader("Content-Type", "application/vnd.config-gtm.v1.4+json;charset=UTF-8").
		BodyString(`{
                    "resource" : {
                        "assignments": [ {
                                "countries": [ "GB", "IE" ],
                                "datacenterId": 3133,
                                "nickname": "UK and Ireland users"
                        } ],
                        "defaultDatacenter": {
                                "datacenterId": 5400,
                                "nickname": "Default Mapping"
                        },
                        "links": [ {
                                "href": "/config-gtm/v1/domains/example.akadns.net/geographic-maps/UK%20Delivery",
                                "rel": "self"
                        } ],
                        "name": "testGeoMap"
                    },
                    "status" : {
                           "changeId": "93a48b86-4fc3-4a5f-9ca2-036835034cc6",
                           "links": [
                               {
                                  "href": "/config-gtm/v1/domains/example.akadns.net/status/current",
                                  "rel": "self"
                               }
                           ],
                           "message": "Change Pending",
                           "passingValidation": true,
                           "propagationStatus": "PENDING",
                           "propagationStatusDate": "2014-04-15T11:30:27.000+0000"
                    }
               }`)

	Init(config)

	testGeoMap := instantiateGeoMap()
	statresp, err := testGeoMap.Create(gtmTestDomain)
	assert.NoError(t, err)

	assert.IsType(t, &GeoMap{}, statresp.Resource)
	assert.Equal(t, GtmTestGeoMap, statresp.Resource.Name)

}

func TestUpdateGeoMap(t *testing.T) {

	defer gock.Off()

	mock := gock.New("https://akaa-baseurl-xxxxxxxxxxx-xxxxxxxxxxxxx.luna.akamaiapis.net/config-gtm/v1/domains/" + gtmTestDomain + "/geographic-maps/" + GtmTestGeoMap)
	mock.
		Put("/config-gtm/v1/domains/"+gtmTestDomain+"/geographic-maps/"+GtmTestGeoMap).
		HeaderPresent("Authorization").
		Reply(200).
		SetHeader("Content-Type", "application/vnd.config-gtm.v1.4+json;charset=UTF-8").
		BodyString(`{
                    "resource" : {
                        "assignments": [ {
                                "countries": [ "GB", "IE" ],
                                "datacenterId": 3133,
                                "nickname": "UK and Ireland users"
                        } ],
                        "defaultDatacenter": {
                                "datacenterId": 5400,
                                "nickname": "Default Mapping"
                        },
                        "links": [ {
                                "href": "/config-gtm/v1/domains/example.akadns.net/geographic-maps/UK%20Delivery",
                                "rel": "self"
                        } ],
                        "name": "testGeoMap"
                    },
                    "status" : {
                           "changeId": "93a48b86-4fc3-4a5f-9ca2-036835034cc6",
                           "links": [
                               {
                                  "href": "/config-gtm/v1/domains/example.akadns.net/status/current",
                                  "rel": "self"
                               }
                           ],
                           "message": "Change Pending",
                           "passingValidation": true,
                           "propagationStatus": "PENDING",
                           "propagationStatusDate": "2014-04-15T11:30:27.000+0000"
                    }
               }`)

	Init(config)

	testGeoMap := instantiateGeoMap()
	_, err := testGeoMap.Update(gtmTestDomain)
	assert.NoError(t, err)

}

func TestDeleteGeoMap(t *testing.T) {

	defer gock.Off()

	mock := gock.New("https://akaa-baseurl-xxxxxxxxxxx-xxxxxxxxxxxxx.luna.akamaiapis.net/config-gtm/v1/domains/" + gtmTestDomain + "/geographic-maps/" + GtmTestGeoMap)
	mock.
		Delete("/config-gtm/v1/domains/"+gtmTestDomain+"/geographic-maps/"+GtmTestGeoMap).
		HeaderPresent("Authorization").
		Reply(200).
		SetHeader("Content-Type", "application/vnd.config-gtm.v1.4+json;charset=UTF-8").
		BodyString(`{
                    "resource" : {
                    },
                    "status" : {
                           "changeId": "93a48b86-4fc3-4a5f-9ca2-036835034cc6",
                           "links": [
                               {
                                  "href": "/config-gtm/v1/domains/example.akadns.net/status/current",
                                  "rel": "self"
                               }
                           ],
                           "message": "Change Pending",
                           "passingValidation": true,
                           "propagationStatus": "PENDING",
                           "propagationStatusDate": "2014-04-15T11:30:27.000+0000"
                    }
               }`)

	Init(config)

	getGeoMap := instantiateGeoMap()
	stat, err := getGeoMap.Delete(gtmTestDomain)
	assert.NoError(t, err)
	assert.Equal(t, "93a48b86-4fc3-4a5f-9ca2-036835034cc6", stat.ChangeId)

}
