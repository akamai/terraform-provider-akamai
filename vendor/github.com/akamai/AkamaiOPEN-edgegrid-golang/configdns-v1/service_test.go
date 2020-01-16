package dns

import (
	"testing"

	"github.com/akamai/AkamaiOPEN-edgegrid-golang/edgegrid"
	"github.com/stretchr/testify/assert"
	"gopkg.in/h2non/gock.v1"
)

var (
	config = edgegrid.Config{
		Host:         "akaa-baseurl-xxxxxxxxxxx-xxxxxxxxxxxxx.luna.akamaiapis.net/",
		AccessToken:  "akab-access-token-xxx-xxxxxxxxxxxxxxxx",
		ClientToken:  "akab-client-token-xxx-xxxxxxxxxxxxxxxx",
		ClientSecret: "xxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxx=",
		MaxBody:      2048,
		Debug:        false,
	}
)

func TestGetZoneSimple(t *testing.T) {
	defer gock.Off()

	mock := gock.New("https://akaa-baseurl-xxxxxxxxxxx-xxxxxxxxxxxxx.luna.akamaiapis.net/config-dns/v1/zones/example.com")
	mock.
		Get("/config-dns/v1/zones/example.com").
		HeaderPresent("Authorization").
		Reply(200).
		SetHeader("Content-Type", "application/json").
		BodyString(`{
			"token": "a184671d5307a388180fbf7f11dbdf46",
			"zone": {
				"name": "example.com",
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
				"ns": [
					{
						"active": true,
						"name": "",
						"target": "use4.akam.net.",
						"ttl": 3600
					},
					{
						"active": true,
						"name": "",
						"target": "use3.akam.net.",
						"ttl": 3600
					}
				],
				"a": [
					{
						"active": true,
						"name": "www",
						"target": "1.2.3.4",
						"ttl": 30
					}
				]
			}
		}`)

	Init(config)
	zone, err := GetZone("example.com")

	assert.NoError(t, err)

	assert.IsType(t, &Zone{}, zone)
	assert.Equal(t, "a184671d5307a388180fbf7f11dbdf46", zone.Token)
	assert.Equal(t, "example.com", zone.Zone.Name)

	assert.IsType(t, &SoaRecord{}, zone.Zone.Soa)
	assert.Equal(t, "hostmaster.akamai.com.", zone.Zone.Soa.Contact)
	assert.Equal(t, 604800, zone.Zone.Soa.Expire)
	assert.Equal(t, uint(180), zone.Zone.Soa.Minimum)
	assert.Equal(t, "use4.akamai.com.", zone.Zone.Soa.Originserver)
	assert.Equal(t, 900, zone.Zone.Soa.Refresh)
	assert.Equal(t, 300, zone.Zone.Soa.Retry)
	assert.Equal(t, uint(1271354824), zone.Zone.Soa.Serial)
	assert.Equal(t, 900, zone.Zone.Soa.TTL)

	assert.IsType(t, []*NsRecord{}, zone.Zone.Ns)
	assert.Len(t, zone.Zone.Ns, 2)

	assert.Equal(t, true, zone.Zone.Ns[0].Active)
	assert.Equal(t, "", zone.Zone.Ns[0].Name)
	assert.Equal(t, "use4.akam.net.", zone.Zone.Ns[0].Target)
	assert.Equal(t, 3600, zone.Zone.Ns[0].TTL)

	assert.Equal(t, true, zone.Zone.Ns[1].Active)
	assert.Equal(t, "", zone.Zone.Ns[1].Name)
	assert.Equal(t, "use3.akam.net.", zone.Zone.Ns[1].Target)
	assert.Equal(t, 3600, zone.Zone.Ns[1].TTL)

	assert.IsType(t, []*ARecord{}, zone.Zone.A)
	assert.Len(t, zone.Zone.A, 1)

	assert.Equal(t, true, zone.Zone.A[0].Active)
	assert.Equal(t, "www", zone.Zone.A[0].Name)
	assert.Equal(t, "1.2.3.4", zone.Zone.A[0].Target)
	assert.Equal(t, 30, zone.Zone.A[0].TTL)
}

func TestGetZone(t *testing.T) {
	defer gock.Off()

	tests := testGetZoneCompleteProvider()
	for _, test := range tests {
		test := test
		t.Run(test.name, func(t *testing.T) {
			mock := gock.New("https://akaa-baseurl-xxxxxxxxxxx-xxxxxxxxxxxxx.luna.akamaiapis.net/config-dns/v1/zones/example.com")
			mock.
				Get("/config-dns/v1/zones/example.com").
				HeaderPresent("Authorization").
				Reply(200).
				SetHeader("Content-Type", "application/json").
				BodyString(test.responseBody)

			Init(config)
			zone, err := GetZone("example.com")

			assert.NoError(t, err)

			assert.IsType(t, &Zone{}, zone)

			if test.expectedRecords != nil {
				switch test.recordType {
				case "A":
					assert.Equal(t, len(test.expectedRecords.([]*ARecord)), len(zone.Zone.A))
					assert.ObjectsAreEqual(test.expectedRecords, zone.Zone.A)
					break
				case "AAAA":
					assert.Equal(t, len(test.expectedRecords.([]*AaaaRecord)), len(zone.Zone.Aaaa))
					assert.ObjectsAreEqual(test.expectedRecords, zone.Zone.Aaaa)
					break
				case "AFSDB":
					assert.Equal(t, len(test.expectedRecords.([]*AfsdbRecord)), len(zone.Zone.Afsdb))
					assert.ObjectsAreEqual(test.expectedRecords, zone.Zone.Afsdb)
					break
				case "CNAME":
					assert.Equal(t, len(test.expectedRecords.([]*CnameRecord)), len(zone.Zone.Cname))
					assert.ObjectsAreEqual(test.expectedRecords, zone.Zone.Cname)
					break
				case "DNSKEY":
					assert.Equal(t, len(test.expectedRecords.([]*DnskeyRecord)), len(zone.Zone.Dnskey))
					assert.ObjectsAreEqual(test.expectedRecords, zone.Zone.Dnskey)
					break
				case "DS":
					assert.Equal(t, len(test.expectedRecords.([]*DsRecord)), len(zone.Zone.Ds))
					assert.ObjectsAreEqual(test.expectedRecords, zone.Zone.Ds)
					break
				case "HINFO":
					assert.Equal(t, len(test.expectedRecords.([]*HinfoRecord)), len(zone.Zone.Hinfo))
					assert.ObjectsAreEqual(test.expectedRecords, zone.Zone.Hinfo)
					break
				case "LOC":
					assert.Equal(t, len(test.expectedRecords.([]*LocRecord)), len(zone.Zone.Loc))
					assert.ObjectsAreEqual(test.expectedRecords, zone.Zone.Loc)
					break
				case "MX":
					assert.Equal(t, len(test.expectedRecords.([]*MxRecord)), len(zone.Zone.Mx))
					assert.ObjectsAreEqual(test.expectedRecords, zone.Zone.Mx)
					break
				case "NAPTR":
					assert.Equal(t, len(test.expectedRecords.([]*NaptrRecord)), len(zone.Zone.Naptr))
					assert.ObjectsAreEqual(test.expectedRecords, zone.Zone.Naptr)
					break
				case "NS":
					assert.Equal(t, len(test.expectedRecords.([]*NsRecord)), len(zone.Zone.Ns))
					assert.ObjectsAreEqual(test.expectedRecords, zone.Zone.Ns)
					break
				case "NSEC3":
					assert.Equal(t, len(test.expectedRecords.([]*Nsec3Record)), len(zone.Zone.Nsec3))
					assert.ObjectsAreEqual(test.expectedRecords, zone.Zone.Nsec3)
					break
				case "NSEC3PARAM":
					assert.Equal(t, len(test.expectedRecords.([]*Nsec3paramRecord)), len(zone.Zone.Nsec3param))
					assert.ObjectsAreEqual(test.expectedRecords, zone.Zone.Nsec3param)
					break
				case "PTR":
					assert.Equal(t, len(test.expectedRecords.([]*PtrRecord)), len(zone.Zone.Ptr))
					assert.ObjectsAreEqual(test.expectedRecords, zone.Zone.Ptr)
					break
				case "RP":
					assert.Equal(t, len(test.expectedRecords.([]*RpRecord)), len(zone.Zone.Rp))
					assert.ObjectsAreEqual(test.expectedRecords, zone.Zone.Rp)
					break
				case "RRSIG":
					assert.Equal(t, len(test.expectedRecords.([]*RrsigRecord)), len(zone.Zone.Rrsig))
					assert.ObjectsAreEqual(test.expectedRecords, zone.Zone.Rrsig)
					break
				case "SPF":
					assert.Equal(t, len(test.expectedRecords.([]*SpfRecord)), len(zone.Zone.Spf))
					assert.ObjectsAreEqual(test.expectedRecords, zone.Zone.Spf)
					break
				case "SRV":
					assert.Equal(t, len(test.expectedRecords.([]*SrvRecord)), len(zone.Zone.Srv))
					assert.ObjectsAreEqual(test.expectedRecords, zone.Zone.Srv)
					break
				case "SSHFP":
					assert.Equal(t, len(test.expectedRecords.([]*SshfpRecord)), len(zone.Zone.Sshfp))
					assert.ObjectsAreEqual(test.expectedRecords, zone.Zone.Sshfp)
					break
				case "TXT":
					assert.Equal(t, len(test.expectedRecords.([]*TxtRecord)), len(zone.Zone.Txt))
					assert.ObjectsAreEqual(test.expectedRecords, zone.Zone.Txt)
					break
				}
			}
		})
	}
}

type recordTests []struct {
	name            string
	recordType      string
	responseBody    string
	expectedType    interface{}
	expectedRecord  ARecord
	expectedRecords interface{}
}

func testGetZoneCompleteProvider() recordTests {
	return recordTests{
		{
			name:       "A Records",
			recordType: "A",
			responseBody: `{
				"zone": {
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
					]
				}
			}`,
			expectedType: []*ARecord{},
			expectedRecords: []*ARecord{
				&ARecord{
					Active: true,
					Name:   "arecord",
					Target: "1.2.3.5",
					TTL:    3600,
				},
				&ARecord{
					Active: true,
					Name:   "origin",
					Target: "1.2.3.9",
					TTL:    3600,
				},
				&ARecord{
					Active: true,
					Name:   "arecord",
					Target: "1.2.3.4",
					TTL:    3600,
				},
			},
		},
		{
			name:       "AAAA Records",
			recordType: "AAAA",
			responseBody: `{
				"zone": {
					"aaaa": [
						{
							"active": true,
							"name": "ipv6record",
							"target": "2001:0db8::ff00:0042:8329",
							"ttl": 3600
						}
					]
				}
			}`,
			expectedType: []*AaaaRecord{},
			expectedRecords: []*AaaaRecord{
				&AaaaRecord{
					Active: true,
					Name:   "ipv6record",
					Target: "2001:0db8::ff00:0042:8329",
					TTL:    3600,
				},
			},
		},
		{
			name:       "AFSDB Records",
			recordType: "AFSDB",
			responseBody: `{
				"zone": {
					"afsdb": [
						{
							"active": true,
							"name": "afsdb",
							"subtype": 1,
							"target": "example.com.",
							"ttl": 7200
						}
					]
				}
			}`,
			expectedType: []*AfsdbRecord{},
			expectedRecords: []*AfsdbRecord{
				&AfsdbRecord{
					Active:  true,
					Name:    "afsdb",
					Subtype: 1,
					Target:  "example.com.",
					TTL:     7200,
				},
			},
		},
		{
			name:       "CNAME Records",
			recordType: "CNAME",
			responseBody: `{
				"zone": {
					"cname": [
						{
							"active": true,
							"name": "redirect",
							"target": "arecord.example.com.",
							"ttl": 3600
						}
					]
				}
			}`,
			expectedType: []*CnameRecord{},
			expectedRecords: []*CnameRecord{
				&CnameRecord{
					Active: true,
					Name:   "redirect",
					Target: "arecord.example.com.",
					TTL:    3600,
				},
			},
		},
		{
			name:       "DNSKEY Records",
			recordType: "DNSKEY",
			responseBody: `{
				"zone": {
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
					]
				}
			}`,
			expectedType: []*DnskeyRecord{},
			expectedRecords: []*DnskeyRecord{
				&DnskeyRecord{
					Active:    true,
					Algorithm: 3,
					Flags:     257,
					Key:       "Av//0/goGKPtaa28nQvPoUwVQ ... i/0hC+1CrmQkuuKtQt98WObuv7q8iQ==",
					Name:      "dnskey",
					Protocol:  7,
					TTL:       7200,
				},
			},
		},
		{
			name:       "DS Records",
			recordType: "DS",
			responseBody: `{
				"zone": {
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
					]
				}
			}`,
			expectedType: []*DsRecord{},
			expectedRecords: []*DsRecord{
				&DsRecord{
					Active:     true,
					Algorithm:  7,
					Digest:     "909FF0B4DD66F91F56524C4F968D13083BE42380",
					DigestType: 1,
					Keytag:     30336,
					Name:       "ds",
					TTL:        7200,
				},
			},
		},
		{
			name:       "HINFO Records",
			recordType: "HINFO",
			responseBody: `{
				"zone": {
					"hinfo": [
						{
							"active": true,
							"hardware": "INTEL-386",
							"name": "hinfo",
							"software": "UNIX",
							"ttl": 7200
						}
					]
				}
			}`,
			expectedType: []*HinfoRecord{},
			expectedRecords: []*HinfoRecord{
				&HinfoRecord{
					Active:   true,
					Hardware: "INTEL-386",
					Name:     "hinfo",
					Software: "UNIX",
					TTL:      7200,
				},
			},
		},
		{
			name:       "LOC Records",
			recordType: "LOC",
			responseBody: `{
				"zone": {
					"loc": [
						{
							"active": true,
							"name": "location",
							"target": "51 30 12.748 N 0 7 39.611 W 0.00m 0.00m 0.00m 0.00m",
							"ttl": 7200
						}
					]
				}
			}`,
			expectedType: []*LocRecord{},
			expectedRecords: []*LocRecord{
				&LocRecord{
					Active: true,
					Name:   "location",
					Target: "51 30 12.748 N 0 7 39.611 W 0.00m 0.00m 0.00m 0.00m",
					TTL:    7200,
				},
			},
		},
		{
			name:       "MX Records",
			recordType: "MX",
			responseBody: `{
				"zone": {
					"mx": [
						{
							"active": true,
							"name": "four",
							"priority": 10,
							"target": "mx1.akamai.com.",
							"ttl": 7200
						}
					]
				}
			}`,
			expectedType: []*MxRecord{},
			expectedRecords: []*MxRecord{
				&MxRecord{
					Active:   true,
					Name:     "four",
					Priority: 10,
					Target:   "mx1.akamai.com.",
					TTL:      7200,
				},
			},
		},
		{
			name:       "NAPTR Records",
			recordType: "NAPTR",
			responseBody: `{
				"zone": {
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
					]
				}
			}`,
			expectedType: []*NaptrRecord{},
			expectedRecords: []*NaptrRecord{
				&NaptrRecord{
					Active:      true,
					Flags:       "S",
					Name:        "naptrrecord",
					Order:       0,
					Preference:  10,
					Regexp:      "!^.*$!sip:customer-service@example.com!",
					Replacement: ".",
					Service:     "SIP+D2U",
					TTL:         3600,
				},
			},
		},
		{
			name:       "NS Records",
			recordType: "NS",
			responseBody: `{
				"zone": {
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
					]
				}
			}`,
			expectedType: []*NsRecord{},
			expectedRecords: []*NsRecord{
				&NsRecord{
					Active: true,
					Target: "use4.akam.net.",
					TTL:    3600,
				},
				&NsRecord{
					Active: true,
					Target: "us34.akam.net.",
					TTL:    3600,
				},
				&NsRecord{
					Active: true,
					Name:   "five",
					Target: "use4.akam.net.",
					TTL:    172800,
				},
			},
		},
		{
			name:       "NSEC3 Records",
			recordType: "NSEC3",
			responseBody: `{
				"zone": {
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
					]
				}
			}`,
			expectedType: []*Nsec3Record{},
			expectedRecords: []*Nsec3Record{
				&Nsec3Record{
					Active:              true,
					Algorithm:           1,
					Flags:               0,
					Iterations:          1,
					Name:                "qdeo8lqu4l81uo67oolpo9h0nv9l13dh",
					NextHashedOwnerName: "R2NUSMGFSEUHT195P59KOU2AI30JR96P",
					Salt:                "EBD1E0942543A01B",
					TTL:                 7200,
					TypeBitmaps:         "CNAME RRSIG",
				},
			},
		},
		{
			name:       "NSEC3PARAM Records",
			recordType: "NSEC3PARAM",
			responseBody: `{
				"zone": {
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
					]
				}
			}`,
			expectedType: []*Nsec3paramRecord{},
			expectedRecords: []*Nsec3paramRecord{
				&Nsec3paramRecord{
					Active:     true,
					Algorithm:  1,
					Flags:      0,
					Iterations: 1,
					Name:       "qnsec3param",
					Salt:       "EBD1E0942543A01B",
					TTL:        7200,
				},
			},
		},
		{
			name:       "PTR Records",
			recordType: "PTR",
			responseBody: `{
				"zone": {
					"ptr": [
						{
							"active": true,
							"name": "ptr",
							"target": "ptr.example.com.",
							"ttl": 7200
						}
					]
				}
			}`,
			expectedType: []*PtrRecord{},
			expectedRecords: []*PtrRecord{
				&PtrRecord{
					Active: true,
					Name:   "ptr",
					Target: "ptr.example.com.",
					TTL:    7200,
				},
			},
		},
		{
			name:       "RP Records",
			recordType: "RP",
			responseBody: `{
				"zone": {
					"rp": [
						{
							"active": true,
							"mailbox": "admin.example.com.",
							"name": "rp",
							"ttl": 7200,
							"txt": "txt.example.com."
						}
					]
				}
			}`,
			expectedType: []*RpRecord{},
			expectedRecords: []*RpRecord{
				&RpRecord{
					Active:  true,
					Mailbox: "admin.example.com.",
					Name:    "rp",
					TTL:     7200,
					Txt:     "txt.example.com.",
				},
			},
		},
		{
			name:       "RRSIG Records",
			recordType: "RRSIG",
			responseBody: `{
				"zone": {
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
					]
				}
			}`,
			expectedType: []*RrsigRecord{},
			expectedRecords: []*RrsigRecord{
				&RrsigRecord{
					Active:      true,
					Algorithm:   7,
					Expiration:  "20120318104101",
					Inception:   "20120315094101",
					Keytag:      63761,
					Labels:      3,
					Name:        "arecord",
					OriginalTTL: 3600,
					Signature:   "toCy19QnAb86vRlQjf5 ... z1doJdHEr8PiI+Is9Eafxh+4Idcw8Ysv",
					Signer:      "example.com.",
					TTL:         7200,
					TypeCovered: "A",
				},
			},
		},
		{
			name:       "SPF Records",
			recordType: "SPF",
			responseBody: `{
				"zone": {
					"spf": [
						{
							"active": true,
							"name": "spf",
							"target": "v=spf a -all",
							"ttl": 7200
						}
					]
				}
			}`,
			expectedType: []*SpfRecord{},
			expectedRecords: []*SpfRecord{
				&SpfRecord{
					Active: true,
					Name:   "spf",
					Target: "v=spf a -all",
					TTL:    7200,
				},
			},
		},
		{
			name:       "SSHFP Records",
			recordType: "SSHFP",
			responseBody: `{
				"zone": {
					"sshfp": [
						{
							"active": true,
							"algorithm": 2,
							"fingerprint": "123456789ABCDEF67890123456789ABCDEF67890",
							"fingerprint_type": 1,
							"name": "host",
							"ttl": 3600
						}
					]
				}
			}`,
			expectedType: []*SshfpRecord{},
			expectedRecords: []*SshfpRecord{
				&SshfpRecord{
					Active:          true,
					Algorithm:       2,
					Fingerprint:     "123456789ABCDEF67890123456789ABCDEF67890",
					FingerprintType: 1,
					Name:            "host",
					TTL:             3600,
				},
			},
		},
		{
			name:       "SRV Records",
			recordType: "SRV",
			responseBody: `{
				"zone": {
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
					]
				}
			}`,
			expectedType: []*SrvRecord{},
			expectedRecords: []*SrvRecord{
				&SrvRecord{
					Active:   true,
					Name:     "srv",
					Port:     522,
					Priority: 10,
					Target:   "target.akamai.com.",
					TTL:      7200,
					Weight:   0,
				},
			},
		},
		{
			name:       "TXT Records",
			recordType: "TXT",
			responseBody: `{
				"zone": {
					"txt": [
						{
							"active": true,
							"name": "text",
							"target": "Hello world!",
							"ttl": 7200
						}
					]
				}
			}`,
			expectedType: []*TxtRecord{},
			expectedRecords: []*TxtRecord{
				&TxtRecord{
					Active: true,
					Name:   "text",
					Target: "Hello world!",
					TTL:    7200,
				},
			},
		},
	}
}
