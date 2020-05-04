package dnsv2

import (
	"testing"
	//edge "github.com/akamai/AkamaiOPEN-edgegrid-golang/edgegrid"
	"fmt"
	"github.com/akamai/AkamaiOPEN-edgegrid-golang/jsonhooks-v1"
	"github.com/stretchr/testify/assert"
	"gopkg.in/h2non/gock.v1"
)

var (
	dnsTestZone = "dns.test.zone.com"
	tsigBody    = []byte(fmt.Sprintf(`{
    			"name": "%s",
    			"algorithm": "hmac-md5",
    			"secret": "p/jzrJpXOLf4mPUtx/z+Sw=="
		}`, dnsTestZone))

	tsigKeyResponse = fmt.Sprintf(`{
                        "name": "%s",
                        "algorithm": "hmac-md5",
                        "secret": "p/jzrJpXOLf4mPUtx/z+Sw==",
			"zonesCount": 1
                }`, dnsTestZone)
)

func createTestTsigKey() *TSIGKey {

	key := &TSIGKey{}
	jsonhooks.Unmarshal(tsigBody, key)

	return key

}

func TestListTsigKeys(t *testing.T) {

	/*
			// for live testing
		        config, err := edge.Init("","")
		        if err != nil {
		                t.Fatalf("TestListTsigKeys failed initializing: %s", err.Error())
		        }
	*/

	defer gock.Off()

	mock := gock.New("https://akaa-baseurl-xxxxxxxxxxx-xxxxxxxxxxxxx.luna.akamaiapis.net/config-dns/v2/keys")
	mock.
		Get("/config-dns/v2/keys").
		HeaderPresent("Authorization").
		Reply(200).
		SetHeader("Content-Type", "application/json;charset=UTF-8").
		BodyString(fmt.Sprint(`{
			"metadata": {
				"totalElements":2
			},
			"keys": [{
				"name":"mbnewtsig",
				"algorithm":"HMAC-MD5.SIG-ALG.REG.INT",
				"secret":"abc78w==",
				"zonesCount":2
			},{
				"name":"fred",
				"algorithm":"hmac-sha256",
				"secret":"IxSErTXxsCN8JO1jsAqW4We0rwbdu5R2jwFFmXoS//Y=",
				"zonesCount":1
			}]
		}`))

	Init(config)
	tsigQueryString := &TSIGQueryString{}
	tsigReport, err := ListTsigKeys(tsigQueryString)
	assert.NoError(t, err)
	assert.Equal(t, int64(len(tsigReport.Keys)), tsigReport.Metadata.TotalElements)

}

func TestUpdateZoneKey(t *testing.T) {

	defer gock.Off()

	mock := gock.New(fmt.Sprintf("https://akaa-baseurl-xxxxxxxxxxx-xxxxxxxxxxxxx.luna.akamaiapis.net/config-dns/v2/zones/%s/key", dnsTestZone))
	mock.
		Put(fmt.Sprintf("/config-dns/v2/zones/%s/key", dnsTestZone)).
		HeaderPresent("Authorization").
		Reply(204).
		SetHeader("Content-Type", "application/json;charset=UTF-8")

	Init(config)
	testKey := createTestTsigKey()
	err := testKey.Update(dnsTestZone)
	assert.NoError(t, err)
	/*
			// live testing ...
			zoneResp, err := GetZone("xxxxxxxxxxxxx.com")
		        assert.NoError(t, err)
			assert.Equal(t, testKey.Name, zoneResp.TsigKey.Name)
	*/

}

func TestGetZoneKey(t *testing.T) {

	defer gock.Off()

	mock := gock.New(fmt.Sprintf("https://akaa-baseurl-xxxxxxxxxxx-xxxxxxxxxxxxx.luna.akamaiapis.net/config-dns/v2/zones/%s/key", dnsTestZone))
	mock.
		Get(fmt.Sprintf("/config-dns/v2/zones/%s/key", dnsTestZone)).
		HeaderPresent("Authorization").
		Reply(200).
		SetHeader("Content-Type", "application/json;charset=UTF-8").
		BodyString(tsigKeyResponse)
	Init(config)
	keyResp, err := GetZoneKey(dnsTestZone)
	if err == nil {
		fmt.Sprintf("Key resp: %v", keyResp)
	} else {
		fmt.Sprintf(err.Error())
	}
	assert.NoError(t, err)
	assert.Equal(t, keyResp.Name, createTestTsigKey().Name)

}

func TestDeleteZoneKey(t *testing.T) {

	defer gock.Off()

	mock := gock.New(fmt.Sprintf("https://akaa-baseurl-xxxxxxxxxxx-xxxxxxxxxxxxx.luna.akamaiapis.net/config-dns/v2/zones/%s/key", dnsTestZone))
	mock.
		Delete(fmt.Sprintf("/config-dns/v2/zones/%s/key", dnsTestZone)).
		HeaderPresent("Authorization").
		Reply(204).
		SetHeader("Content-Type", "application/json;charset=UTF-8")

	Init(config)
	err := DeleteZoneKey(dnsTestZone)
	assert.NoError(t, err)

}

func TestGetTsigKeyUsers(t *testing.T) {

	//
	// NOTE: Currently a discrepency between docs and API operation!!!
	// TODO: Reconcile and correct
	//

	defer gock.Off()

	mock := gock.New(fmt.Sprintf("https://akaa-baseurl-xxxxxxxxxxx-xxxxxxxxxxxxx.luna.akamaiapis.net/config-dns/v2/zones/%s/key/used-by", dnsTestZone))
	mock.
		Get(fmt.Sprintf("/config-dns/v2/zones/%s/key/used-by", dnsTestZone)).
		HeaderPresent("Authorization").
		Reply(200).
		SetHeader("Content-Type", "application/json;charset=UTF-8").
		BodyString(fmt.Sprintf(`{
                    "zones":["%s"]
		}`, dnsTestZone))

	Init(config)
	zoneKeyAliases, err := GetZoneKeyAliases(dnsTestZone)
	assert.NoError(t, err)
	// CORRECT
	assert.Equal(t, assert.IsType(t, &ZoneNameListResponse{}, zoneKeyAliases), true)
	//assert.Equal(t, assert.IsType(t, &TSIGZoneAliases{}, zoneKeyAliases), true)

}

func TestTsigKeyBulkUpdate(t *testing.T) {

	defer gock.Off()

	mock := gock.New("https://akaa-baseurl-xxxxxxxxxxx-xxxxxxxxxxxxx.luna.akamaiapis.net/config-dns/v2/keys/bulk-update")
	mock.
		Post("/config-dns/v2/keys/bulk-update").
		HeaderPresent("Authorization").
		Reply(204).
		SetHeader("Content-Type", "application/json;charset=UTF-8")

	Init(config)
	testKey := createTestTsigKey()
	bulkUpdate := &TSIGKeyBulkPost{Key: testKey, Zones: []string{dnsTestZone}}
	err := bulkUpdate.BulkUpdate()
	assert.NoError(t, err)

}

func TestTsigKeyGetZones(t *testing.T) {

	defer gock.Off()

	mock := gock.New("https://akaa-baseurl-xxxxxxxxxxx-xxxxxxxxxxxxx.luna.akamaiapis.net/config-dns/v2/keys/used-by")
	mock.
		Post("/config-dns/v2/keys/used-by").
		HeaderPresent("Authorization").
		Reply(200).
		SetHeader("Content-Type", "application/json;charset=UTF-8").
		BodyString(fmt.Sprintf(`{
                        "zones":["%s"]
                }`, dnsTestZone))

	Init(config)
	testKey := createTestTsigKey()
	zoneList, err := testKey.GetZones()
	assert.NoError(t, err)
	assert.Equal(t, assert.IsType(t, &ZoneNameListResponse{}, zoneList), true)

}
