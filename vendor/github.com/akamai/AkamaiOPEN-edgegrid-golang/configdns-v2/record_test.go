package dnsv2

import (
	"testing"
	"fmt"
	"github.com/stretchr/testify/assert"
	"gopkg.in/h2non/gock.v1"
)

func TestRecord_ContainsHelper(t *testing.T) {
	tm1 := []string{
		"test1",
		"test2",
		"test3",
	}

	assert.Equal(t, contains(tm1, "test1"), true)
	assert.Equal(t, contains(tm1, "test2"), true)
	assert.Equal(t, contains(tm1, "test3"), true)
	assert.Equal(t, contains(tm1, "test4"), false)
}

func contains(arr []string, str string) bool {
	for _, a := range arr {
		if a == str {
			return true
		}
	}
	return false
}
func TestRecord_ARecord(t *testing.T) {

	a := &RecordBody{Name: "test1", RecordType: "A"}

	/*	a := NewARecord()
		f := []string{
			"name",
			"ttl",
			"active",
			"target",
		}*/

	//	assert.Equal(t, a.fieldMap, f)
	//	assert.Equal(t, a.fieldMap, a.GetAllowedFields())
	//assert.Equal(t, a.SetField("name", "test1"), nil)
	assert.Equal(t, "test1", a.Name)
	//assert.Equal(t, a.SetField("doesntExist", "test1"), &RecordError{fieldName: "doesntExist"})
	a.TTL = 900
	//a.SetField("ttl", 900)
	//a.SetField("active", true)
	a.Active = true
	//a.SetField("target", "test2")
	records := make([]string, 0, 1)
	records = append(records, "test2")
	a.Target = records
	assert.Equal(t, a.ToMap(), map[string]interface{}{
		"name":   "test1",
		"ttl":    900,
		"active": true,
		"target": []string{"test2"},
	})
}

func TestRecord_AllRecords_WrongTypes(t *testing.T) {
	/*a := NewARecord()
	e := a.SetField("name", 1)
	assert.Equal(t, e, &RecordError{fieldName: "name"})

	a1 := NewAaaaRecord()
	e1 := a1.SetField("name", 1)
	assert.Equal(t, e1, &RecordError{fieldName: "name"})

	a2 := NewAfsdbRecord()
	e = a2.SetField("name", 1)
	assert.Equal(t, e, &RecordError{fieldName: "name"})

	a3 := NewCnameRecord()
	e = a3.SetField("name", 1)
	assert.Equal(t, e, &RecordError{fieldName: "name"})

	a4 := NewDnskeyRecord()
	e = a4.SetField("name", 1)
	assert.Equal(t, e, &RecordError{fieldName: "name"})

	a5 := NewDsRecord()
	e = a5.SetField("name", 1)
	assert.Equal(t, e, &RecordError{fieldName: "name"})

	a6 := NewHinfoRecord()
	e = a6.SetField("name", 1)
	assert.Equal(t, e, &RecordError{fieldName: "name"})

	a7 := NewLocRecord()
	e = a7.SetField("name", 1)
	assert.Equal(t, e, &RecordError{fieldName: "name"})

	a8 := NewMxRecord()
	e = a8.SetField("name", 1)
	assert.Equal(t, e, &RecordError{fieldName: "name"})

	a9 := NewNaptrRecord()
	e = a9.SetField("name", 1)
	assert.Equal(t, e, &RecordError{fieldName: "name"})

	a10 := NewNsRecord()
	e = a10.SetField("name", 1)
	assert.Equal(t, e, &RecordError{fieldName: "name"})

	a11 := NewNsec3Record()
	e = a11.SetField("name", 1)
	assert.Equal(t, e, &RecordError{fieldName: "name"})

	a12 := NewNsec3paramRecord()
	e = a12.SetField("name", 1)
	assert.Equal(t, e, &RecordError{fieldName: "name"})

	a13 := NewPtrRecord()
	e = a13.SetField("name", 1)
	assert.Equal(t, e, &RecordError{fieldName: "name"})

	a14 := NewRpRecord()
	e = a14.SetField("name", 1)
	assert.Equal(t, e, &RecordError{fieldName: "name"})

	a15 := NewRrsigRecord()
	e = a15.SetField("name", 1)
	assert.Equal(t, e, &RecordError{fieldName: "name"})

	a16 := NewSoaRecord()
	e = a16.SetField("ttl", "test")
	assert.Equal(t, e, &RecordError{fieldName: "ttl"})

	a17 := NewSpfRecord()
	e = a17.SetField("name", 1)
	assert.Equal(t, e, &RecordError{fieldName: "name"})

	a18 := NewSrvRecord()
	e = a18.SetField("name", 1)
	assert.Equal(t, e, &RecordError{fieldName: "name"})

	a19 := NewSshfpRecord()
	e = a19.SetField("name", 1)
	assert.Equal(t, e, &RecordError{fieldName: "name"})

	a20 := NewTxtRecord()
	e = a20.SetField("name", 1)
	assert.Equal(t, e, &RecordError{fieldName: "name"})
	*/
}

func TestGetRecordsetsNoArgs(t *testing.T) {

        dnsTestZone := "testzone.com"

        defer gock.Off()

        mock := gock.New(fmt.Sprintf("https://akaa-baseurl-xxxxxxxxxxx-xxxxxxxxxxxxx.luna.akamaiapis.net/config-dns/v2/zones/%s/recordsets", dnsTestZone))
        mock.
                Get(fmt.Sprintf("/config-dns/v2/zones/%s/recordsets", dnsTestZone)).
                HeaderPresent("Authorization").
                Reply(200).
                SetHeader("Content-Type", "application/json;charset=UTF-8").
                BodyString(fmt.Sprintf(`{
    			"metadata": {
        			"lastPage": 1, 
        			"page": 1, 
        			"pageSize": 25, 
        			"showAll": false, 
        			"totalElements": 2
    			}, 
    			"recordsets": [
        			{
            				"name": "east.testzone.com", 
            				"rdata": [
                				"east.testzone.com.edgesuite.net."
            				], 
            				"ttl": 600, 
            				"type": "CNAME"
        			}, 
        			{
            				"name": "easttest.testzone.com", 
            				"rdata": [
                				"easttest.testzone.com.edgesuite.net."
            				], 
            				"ttl": 600, 
            				"type": "CNAME"
        			}] 
		}`))

        Init(config)
	recordsets, err := GetRecordsets(dnsTestZone)
        assert.NoError(t, err)
        assert.Equal(t, assert.IsType(t, &RecordSetResponse{}, recordsets), true)
        assert.Equal(t, len(recordsets.Recordsets), 2)

}

func TestGetRecordsets(t *testing.T) {

        dnsTestZone := "testzone.com"

        defer gock.Off()

        mock := gock.New(fmt.Sprintf("https://akaa-baseurl-xxxxxxxxxxx-xxxxxxxxxxxxx.luna.akamaiapis.net/config-dns/v2/zones/%s/recordsets", dnsTestZone))
        mock.
                Get(fmt.Sprintf("/config-dns/v2/zones/%s/recordsets", dnsTestZone)).
                HeaderPresent("Authorization").
                Reply(200).
                SetHeader("Content-Type", "application/json;charset=UTF-8").
                BodyString(fmt.Sprintf(`{
                        "metadata": {
                                "lastPage": 1,
                                "page": 1,
                                "pageSize": 25,
                                "showAll": false,
                                "totalElements": 1
                        },
                        "recordsets": [
                                {
                                        "name": "east.testzone.com",
                                        "rdata": [
                                                "east.testzone.com.edgesuite.net."
                                        ],
                                        "ttl": 600,
                                        "type": "CNAME"
                                }]
                }`))

        Init(config)
	qargs := RecordsetQueryArgs{Search: "east.testzone.com", SortBy: "name"}
        recordsets, err := GetRecordsets(dnsTestZone, qargs)
        assert.NoError(t, err)
        assert.Equal(t, assert.IsType(t, &RecordSetResponse{}, recordsets), true)
        assert.Equal(t, len(recordsets.Recordsets), 1)

}

func TestGetRecord(t *testing.T) {

        dnsTestZone := "testzone.com"
        dnsTestRecordName := "east.testzone.com"
        dnsTestRecordType := "CNAME"

        defer gock.Off()

        mock := gock.New(fmt.Sprintf("https://akaa-baseurl-xxxxxxxxxxx-xxxxxxxxxxxxx.luna.akamaiapis.net/config-dns/v2/zones/%s/names/%s/types/%s", dnsTestZone, dnsTestRecordName, dnsTestRecordType))
        mock.
                Get(fmt.Sprintf("/config-dns/v2/zones/%s/names/%s/types/%s", dnsTestZone, dnsTestRecordName, dnsTestRecordType)).
                HeaderPresent("Authorization").
                Reply(200).
                SetHeader("Content-Type", "application/json;charset=UTF-8").
                BodyString(fmt.Sprintf(`{
                        "name": "east.testzone.com",
                        "rdata": [
                                 "east.testzone.com.edgesuite.net."
                        ],
                        "ttl": 600,
                        "type": "CNAME"
                }`))

        Init(config)
        testrecord, err := GetRecord(dnsTestZone, dnsTestRecordName, "CNAME")
        assert.NoError(t, err)
        assert.Equal(t, assert.IsType(t, &RecordBody{}, testrecord), true)
        assert.Equal(t, testrecord.Name, dnsTestRecordName)

}

