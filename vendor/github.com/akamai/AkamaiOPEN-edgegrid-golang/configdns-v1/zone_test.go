package dns

import (
	"testing"

	"github.com/akamai/AkamaiOPEN-edgegrid-golang/jsonhooks-v1"
	"github.com/stretchr/testify/assert"
)

func TestZone_JSON(t *testing.T) {
	responseBody := []byte(`{
	"token": "a184671d5307a388180fbf7f11dbdf46",
	"zone": {
		"name": "example.com",
		"a": [
			{
				"active": true,
				"name": "arecord",
				"target": "1.2.3.5",
				"ttl": 3600
			},
			{
				"active": true,
				"name": "origin",
				"target": "1.2.3.9",
				"ttl": 3600
			},
			{
				"active": true,
				"name": "arecord",
				"target": "1.2.3.4",
				"ttl": 3600
			}
		],
		"aaaa": [
			{
				"active": true,
				"name": "ipv6record",
				"target": "2001:0db8::ff00:0042:8329",
				"ttl": 3600
			}
		],
		"afsdb": [
			{
				"active": true,
				"name": "afsdb",
				"subtype": 1,
				"target": "example.com.",
				"ttl": 7200
			}
		],
		"cname": [
			{
				"active": true,
				"name": "redirect",
				"target": "arecord.example.com.",
				"ttl": 3600
			}
		],
		"dnskey": [
			{
				"active": true,
				"algorithm": 3,
				"flags": 257,
				"key": "Av//0/goGKPtaa28nQvPoUwVQ ... i/0hC+1CrmQkuuKtQt98WObuv7q8iQ==",
				"name": "dnskey",
				"protocol": 7,
				"ttl": 7200
			}
		],
		"ds": [
			{
				"active": true,
				"algorithm": 7,
				"digest": "909FF0B4DD66F91F56524C4F968D13083BE42380",
				"digest_type": 1,
				"keytag": 30336,
				"name": "ds",
				"ttl": 7200
			}
		],
		"hinfo": [
			{
				"active": true,
				"hardware": "INTEL-386",
				"name": "hinfo",
				"software": "UNIX",
				"ttl": 7200
			}
		],
		"loc": [
			{
				"active": true,
				"name": "location",
				"target": "51 30 12.748 N 0 7 39.611 W 0.00m 0.00m 0.00m 0.00m",
				"ttl": 7200
			}
		],
		"mx": [
			{
				"active": true,
				"name": "four",
				"priority": 10,
				"target": "mx1.akamai.com.",
				"ttl": 7200
			}
		],
		"naptr": [
			{
				"active": true,
				"flags": "S",
				"name": "naptrrecord",
				"order": 0,
				"preference": 10,
				"regexp": "!^.*$!sip:customer-service@example.com!",
				"replacement": ".",
				"service": "SIP+D2U",
				"ttl": 3600
			}
		],
		"ns": [
			{
				"active": true,
				"name": null,
				"target": "use4.akam.net.",
				"ttl": 3600
			},
			{
				"active": true,
				"name": null,
				"target": "use3.akam.net.",
				"ttl": 3600
			},
			{
				"active": true,
				"name": "five",
				"target": "use4.akam.net.",
				"ttl": 172800
			}
		],
		"nsec3": [
			{
				"active": true,
				"algorithm": 1,
				"flags": 0,
				"iterations": 1,
				"name": "qdeo8lqu4l81uo67oolpo9h0nv9l13dh",
				"next_hashed_owner_name": "R2NUSMGFSEUHT195P59KOU2AI30JR96P",
				"salt": "EBD1E0942543A01B",
				"ttl": 7200,
				"type_bitmaps": "CNAME RRSIG"
			}
		],
		"nsec3param": [
			{
				"active": true,
				"algorithm": 1,
				"flags": 0,
				"iterations": 1,
				"name": "qnsec3param",
				"salt": "EBD1E0942543A01B",
				"ttl": 7200
			}
		],
		"ptr": [
			{
				"active": true,
				"name": "ptr",
				"target": "ptr.example.com.",
				"ttl": 7200
			}
		],
		"rp": [
			{
				"active": true,
				"mailbox": "admin.example.com.",
				"name": "rp",
				"ttl": 7200,
				"txt": "txt.example.com."
			}
		],
		"rrsig": [
			{
				"active": true,
				"algorithm": 7,
				"expiration": "20120318104101",
				"inception": "20120315094101",
				"keytag": 63761,
				"labels": 3,
				"name": "arecord",
				"original_ttl": 3600,
				"signature": "toCy19QnAb86vRlQjf5 ... z1doJdHEr8PiI+Is9Eafxh+4Idcw8Ysv",
				"signer": "example.com.",
				"ttl": 7200,
				"type_covered": "A"
			}
		],
		"soa": {
			"contact": "hostmaster.akamai.com.",
			"expire": 604800,
			"minimum": 180,
			"originserver": "use4.akamai.com.",
			"refresh": 900,
			"retry": 300,
			"serial": 1271354824,
			"ttl": 900
		},
		"spf": [
			{
				"active": true,
				"name": "spf",
				"target": "v=spf a -all",
				"ttl": 7200
			}
		],
		"srv": [
			{
				"active": true,
				"name": "srv",
				"port": 522,
				"priority": 10,
				"target": "target.akamai.com.",
				"ttl": 7200,
				"weight": 0
			}
		],
		"sshfp": [
			{
				"active": true,
				"algorithm": 2,
				"fingerprint": "123456789ABCDEF67890123456789ABCDEF67890",
				"fingerprint_type": 1,
				"name": "host",
				"ttl": 3600
			}
		],
		"txt": [
			{
				"active": true,
				"name": "text",
				"target": "Hello world!",
				"ttl": 7200
			}
		]
	}
}`)

	zone := NewZone("example.com")
	err := jsonhooks.Unmarshal(responseBody, zone)
	assert.NoError(t, err)
	assert.Equal(t, zone.Zone.Soa.Serial, zone.Zone.Soa.originalSerial)

	_, err = jsonhooks.Marshal(zone)
	assert.NoError(t, err)
	assert.True(t, zone.Zone.Soa.Serial-zone.Zone.Soa.originalSerial == 1)

	zone.Zone.Soa.Serial = uint(12345)
	_, err = jsonhooks.Marshal(zone)
	assert.NoError(t, err)
	assert.Equal(t, zone.Zone.Soa.Serial, uint(12345))
	assert.NotEqual(t, zone.Zone.Soa.originalSerial, uint(12345))
}

func TestZone_AddRecord(t *testing.T) {
	records := testZone_AddRecord_Provider()
	zone := NewZone("example.org")
	for _, record := range records {
		err := zone.AddRecord(record)
		assert.NoError(t, err)
	}
	assert.Equal(t, records, zone.Zone.A)
}

func Test_removeNonCnameName(t *testing.T) {
	backup := nonCnameNames
	defer func() {
		nonCnameNames = backup
	}()

	names := []name{
		{recordType: "TXT", name: "test.com"},
		{recordType: "TXT", name: "foo.com"},
		{recordType: "TXT", name: "bar.com"},
		{recordType: "TXT", name: "test.com"},
	}

	for _, n := range names {
		nonCnameNames = append(nonCnameNames, n)
	}

	zone := Zone{}
	zone.removeNonCnameName("test.com")
}

func Test_removeCnameName(t *testing.T) {
	backup := cnameNames
	defer func() {
		cnameNames = backup
	}()

	names := []name{
		{recordType: "TXT", name: "test.com"},
		{recordType: "TXT", name: "foo.com"},
		{recordType: "TXT", name: "bar.com"},
		{recordType: "TXT", name: "test.com"},
	}

	for _, n := range names {
		cnameNames = append(cnameNames, n)
	}

	zone := Zone{}
	zone.removeCnameName("test.com")
}

func testZone_AddRecord_Provider() []*ARecord {
	return []*ARecord{
		&ARecord{
			Name:   "www",
			Active: true,
			Target: "1.2.3.4",
			TTL:    30,
		},
		&ARecord{
			Name:   "www",
			Active: true,
			Target: "1.2.3.4",
			TTL:    30,
		},
		&ARecord{
			Name:   "www",
			Active: true,
			Target: "1.2.3.5",
			TTL:    30,
		},
	}
}
