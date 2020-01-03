package configgtm

import (
	"testing"

	"github.com/akamai/AkamaiOPEN-edgegrid-golang/jsonhooks-v1"

	"github.com/h2non/gock"
	"github.com/stretchr/testify/assert"

	"fmt"
)

var GtmTestDC1 = "testDC1"
var GtmTestDC2 = "testDC2"
var dcMap = map[string]string{"GtmTestDC1": GtmTestDC1, "GtmTestDC2": GtmTestDC2}

func instantiateDatacenter() *Datacenter {

	dc := NewDatacenter()
	dcData := []byte(`{
                                "datacenterId" : 3131,
                                "nickname" : "testDC1",
                                "scorePenalty" : 0,
                                "city" : null,
                                "stateOrProvince" : null,
                                "country" : null,
                                "latitude" : null,
                                "longitude" : null,
                                "cloneOf" : null,
                                "virtual" : true,
                                "defaultLoadObject" : null,
                                "continent" : null,
                                "servermonitorPool" : null,
                                "servermonitorLivenessCount" : null,
                                "servermonitorLoadCount" : null,
                                "pingInterval" : null,
                                "pingPacketSize" : null,
                                "cloudServerTargeting" : false,
                                "cloudServerHostHeaderOverride" : false
                       }`)
	jsonhooks.Unmarshal(dcData, dc)

	return dc

}

// Verify ListDatacenters. Should pass, e.g. no API errors and non nil list.
func TestListDatacenters(t *testing.T) {

	defer gock.Off()

	mock := gock.New("https://akaa-baseurl-xxxxxxxxxxx-xxxxxxxxxxxxx.luna.akamaiapis.net/config-gtm/v1/domains/" + gtmTestDomain + "/datacenters")
	mock.
		Get("/config-gtm/v1/domains/"+gtmTestDomain+"/datacenters").
		HeaderPresent("Authorization").
		Reply(200).
		SetHeader("Content-Type", "application/vnd.config-gtm.v1.3+json;charset=UTF-8").
		BodyString(`{
                        "items" : [ {
                                "datacenterId" : 3131,
                                "nickname" : "testDC1",
                                "scorePenalty" : 0,
                                "city" : null,                            
                                "stateOrProvince" : null,
                                "country" : null,
                                "latitude" : null,
                                "longitude" : null,
                                "cloneOf" : null,
                                "virtual" : true,
                                "defaultLoadObject" : null,
                                "continent" : null,
                                "servermonitorPool" : null,
                                "servermonitorLivenessCount" : null,
                                "servermonitorLoadCount" : null,
                                "pingInterval" : null,
                                "pingPacketSize" : null,
                                "cloudServerTargeting" : false,
                                "cloudServerHostHeaderOverride" : false,
                                "links" : [ {
                                        "rel" : "self",
                                        "href" : "https://akaa-baseurl-xxxxxxxxxxx-xxxxxxxxxxxxx.luna.akamaiapis.net/config-gtm/v1/domains/gtmdomtest.akadns.net/datacenters/3131"
                                } ]
                        } ]
                }`)

	Init(config)

	dcList, err := ListDatacenters(gtmTestDomain)
	assert.NoError(t, err)
	assert.NotEqual(t, dcList, nil)

	if len(dcList) > 0 {
		firstDC := dcList[0]
		assert.Equal(t, firstDC.Nickname, GtmTestDC1)
	} else {
		assert.Equal(t, 0, 1, "ListDatacenters: empty list")
		fmt.Println("ListDatacenters: empty list")
	}

}

// Verify GetDatacenter. Name hardcoded. Should pass, e.g. no API errors and property returned
// Depends on CreateDatacenter
func TestGetDatacenter(t *testing.T) {

	defer gock.Off()

	mock := gock.New("https://akaa-baseurl-xxxxxxxxxxx-xxxxxxxxxxxxx.luna.akamaiapis.net/config-gtm/v1/domains/" + gtmTestDomain + "/datacenters/3131")
	mock.
		Get("/config-gtm/v1/domains/"+gtmTestDomain+"/datacenters/3131").
		HeaderPresent("Authorization").
		Reply(200).
		SetHeader("Content-Type", "application/vnd.config-gtm.v1.3+json;charset=UTF-8").
		BodyString(`{
                                "datacenterId" : 3131,
                                "nickname" : "testDC1",
                                "scorePenalty" : 0,
                                "city" : null,
                                "stateOrProvince" : null,
                                "country" : null,
                                "latitude" : null,
                                "longitude" : null,
                                "cloneOf" : null,
                                "virtual" : true,
                                "defaultLoadObject" : null,
                                "continent" : null,
                                "servermonitorPool" : null,
                                "servermonitorLivenessCount" : null,
                                "servermonitorLoadCount" : null,
                                "pingInterval" : null,
                                "pingPacketSize" : null,
                                "cloudServerTargeting" : false,
                                "cloudServerHostHeaderOverride" : false,
                                "links" : [ {
                                        "rel" : "self",
                                        "href" : "https://akaa-32qkzqewderdchot-d3uwbyqc4pqi2c5l.luna-dev.akamaiapis.net/config-gtm/v1/domains/gtmdomtest.akadns.net/datacenters/3131"
                                } ]
              }`)

	Init(config)

	testDC, err := GetDatacenter(3131, gtmTestDomain)
	assert.NoError(t, err)
	assert.Equal(t, 3131, testDC.DatacenterId)

}

// Verify failed case for GetDatacenter. Should pass, e.g. no API errors and domain not found
func TestGetBadDatacenter(t *testing.T) {

	defer gock.Off()

	mock := gock.New("https://akaa-baseurl-xxxxxxxxxxx-xxxxxxxxxxxxx.luna.akamaiapis.net/config-gtm/v1/domains/" + gtmTestDomain + "/datacenters/9999")
	mock.
		Get("/config-gtm/v1/domains/"+gtmTestDomain+"/datacenters/9999").
		HeaderPresent("Authorization").
		Reply(404).
		SetHeader("Content-Type", "application/vnd.config-gtm.v1.3+json;charset=UTF-8").
		BodyString(`{
                }`)

	Init(config)

	_, err := GetDatacenter(9999, gtmTestDomain)
	assert.Error(t, err)

}

// Test Create datacenter.
func TestCreateDatacenter(t *testing.T) {

	defer gock.Off()

	mock := gock.New("https://akaa-baseurl-xxxxxxxxxxx-xxxxxxxxxxxxx.luna.akamaiapis.net/config-gtm/v1/domains/" + gtmTestDomain + "/datacenters")
	mock.
		Post("/config-gtm/v1/domains/"+gtmTestDomain+"/datacenters").
		HeaderPresent("Authorization").
		Reply(201).
		SetHeader("Content-Type", "application/vnd.config-gtm.v1.3+json;charset=UTF-8").
		BodyString(`{
                        "resource" : {
                                "datacenterId" : 3131,
                                "nickname" : "testDC1",
                                "scorePenalty" : 0,
                                "city" : null,
                                "stateOrProvince" : null,
                                "country" : null,
                                "latitude" : null,
                                "longitude" : null,
                                "cloneOf" : null,
                                "virtual" : true,
                                "defaultLoadObject" : null,
                                "continent" : null,
                                "servermonitorPool" : null,
                                "servermonitorLivenessCount" : null,
                                "servermonitorLoadCount" : null,
                                "pingInterval" : null,
                                "pingPacketSize" : null,
                                "cloudServerTargeting" : false,
                                "cloudServerHostHeaderOverride" : false,
                                "links" : [ {
                                        "rel" : "self",
                                        "href" : "https://akaa-32qkzqewderdchot-d3uwbyqc4pqi2c5l.luna-dev.akamaiapis.net/config-gtm/v1/domains/gtmdomtest.akadns.net/datacenters/3131"
                                } ]
                        },
                        "status" : {
                                "message" : "Change Pending",
                                "changeId" : "4c7e6466-84e1-4895-bdf5-e3608d708d69",
                                "propagationStatus" : "PENDING",
                                "propagationStatusDate" : "2019-05-30T17:47:02.831+00:00",
                                "passingValidation" : true,
                                "links" : [ {
                                        "rel" : "self",
                                        "href" : "https://akaa-32qkzqewderdchot-d3uwbyqc4pqi2c5l.luna-dev.akamaiapis.net/config-gtm/v1/domains/gtmdomtest.akadns.net/status/current"
                                } ]
                        }
                }`)

	testDC := NewDatacenter()
	testDC.Nickname = GtmTestDC1
	statresp, err := testDC.Create(gtmTestDomain)
	assert.NoError(t, err)
	assert.Equal(t, GtmTestDC1, statresp.Resource.Nickname)

}

func TestUpdateDatacenter(t *testing.T) {

	defer gock.Off()

	mock := gock.New("https://akaa-baseurl-xxxxxxxxxxx-xxxxxxxxxxxxx.luna.akamaiapis.net/config-gtm/v1/domains/" + gtmTestDomain + "/datacenters/3131")
	mock.
		Put("/config-gtm/v1/domains/"+gtmTestDomain+"/datacenters/3131").
		HeaderPresent("Authorization").
		Reply(200).
		SetHeader("Content-Type", "application/vnd.config-gtm.v1.3+json;charset=UTF-8").
		BodyString(`{
                        "resource" : {
                                "datacenterId" : 3131,
                                "nickname" : "testDC1",
                                "scorePenalty" : 0,
                                "city" : null,
                                "stateOrProvince" : null,
                                "country" : null,
                                "latitude" : null,
                                "longitude" : null,
                                "cloneOf" : null,
                                "virtual" : true,
                                "defaultLoadObject" : null,
                                "continent" : null,
                                "servermonitorPool" : null,
                                "servermonitorLivenessCount" : null,
                                "servermonitorLoadCount" : null,
                                "pingInterval" : null,
                                "pingPacketSize" : null,
                                "cloudServerTargeting" : false,
                                "cloudServerHostHeaderOverride" : false,
                                "links" : [ {
                                        "rel" : "self",
                                        "href" : "https://akaa-32qkzqewderdchot-d3uwbyqc4pqi2c5l.luna-dev.akamaiapis.net/config-gtm/v1/domains/gtmdomtest.akadns.net/datacenters/3131"
                                } ]
                        },
                        "status" : {
                                "message" : "Change Pending",
                                "changeId" : "4c7e6466-84e1-4895-bdf5-e3608d708d69",
                                "propagationStatus" : "PENDING",
                                "propagationStatusDate" : "2019-05-30T17:47:02.831+00:00",
                                "passingValidation" : true,
                                "links" : [ {
                                        "rel" : "self",
                                        "href" : "https://akaa-32qkzqewderdchot-d3uwbyqc4pqi2c5l.luna-dev.akamaiapis.net/config-gtm/v1/domains/gtmdomtest.akadns.net/status/current"
                                } ]
                        }
               }`)

	Init(config)

	testDC := instantiateDatacenter()
	stat, err := testDC.Update(gtmTestDomain)
	assert.NoError(t, err)
	assert.Equal(t, stat.ChangeId, "4c7e6466-84e1-4895-bdf5-e3608d708d69")

}

func TestDeleteDatacenter(t *testing.T) {

	defer gock.Off()

	mock := gock.New("https://akaa-baseurl-xxxxxxxxxxx-xxxxxxxxxxxxx.luna.akamaiapis.net/config-gtm/v1/domains/" + gtmTestDomain + "/datacenters/3131")
	mock.
		Delete("/config-gtm/v1/domains/"+gtmTestDomain+"/datacenters/3131").
		HeaderPresent("Authorization").
		Reply(200).
		SetHeader("Content-Type", "application/vnd.config-gtm.v1.3+json;charset=UTF-8").
		BodyString(`{
                        "resource" : {
                        },
                        "status" : {
                                "message" : "Change Pending",
                                "changeId" : "4c7e6466-84e1-4895-bdf5-e3608d708d69",
                                "propagationStatus" : "PENDING",
                                "propagationStatusDate" : "2019-05-30T17:47:02.831+00:00",
                                "passingValidation" : true,
                                "links" : [ {
                                        "rel" : "self",
                                        "href" : "https://akaa-32qkzqewderdchot-d3uwbyqc4pqi2c5l.luna-dev.akamaiapis.net/config-gtm/v1/domains/gtmdomtest.akadns.net/status/current"
                                } ]
                        }
                }`)

	Init(config)

	testDC := instantiateDatacenter()
	_, err := testDC.Delete(gtmTestDomain)
	assert.NoError(t, err)

}
