package dnsv2

import (
	"fmt"
	"github.com/akamai/AkamaiOPEN-edgegrid-golang/edgegrid"
	"testing"

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

	mock := gock.New("https://akaa-baseurl-xxxxxxxxxxx-xxxxxxxxxxxxx.luna.akamaiapis.net/config-dns/v2/zones/example.com")
	mock.
		Get("/config-dns/v2/zones/example.com").
		HeaderPresent("Authorization").
		Reply(200).
		SetHeader("Content-Type", "application/json").
		BodyString(`{
                             "contractId": "C-1FRYVV3",
                             "zone": "example.com",
                             "type": "PRIMARY",
                             "comment": "This is a test zone",
                             "versionId": "5e311536-c4b7-4dec-b8a9-1fe3d6042406",
                             "lastActivationDate": "2019-02-23T22:31:48Z",
                             "lastModifiedDate": "2019-02-23T22:31:48Z",
                             "lastModifiedBy": "davey.shafik",
                             "activationState": "PENDING",
                             "signAndServe": false
                            }`)

	Init(config)

	zone, err := GetZone("example.com")

	assert.NoError(t, err)

	assert.IsType(t, &ZoneResponse{}, zone)
	//assert.Equal(t, "a184671d5307a388180fbf7f11dbdf46", zone.Token)
	assert.Equal(t, "example.com", zone.Zone)

	/*assert.IsType(t, &SoaRecord{}, zone.Zone.Soa)
	assert.Equal(t, "hostmaster.akamai.com.", zone.Zone.Soa.Contact)
	assert.Equal(t, 604800, zone.Zone.Soa.Expire)
	assert.Equal(t, uint(180), zone.Zone.Soa.Minimum)
	assert.Equal(t, "use4.akamai.com.", zone.Zone.Soa.Originserver)
	assert.Equal(t, 900, zone.Zone.Soa.Refresh)
	assert.Equal(t, 300, zone.Zone.Soa.Retry)
	assert.Equal(t, uint(1271354824), zone.Zone.Soa.Serial)
	assert.Equal(t, 900, zone.Zone.Soa.TTL)
	*/
	/*
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
	*/
}

func TestGetZoneRecords(t *testing.T) {
	defer gock.Off()

	tests := testGetZoneCompleteProvider()
	for _, test := range tests {
		test := test
		t.Run(test.name, func(t *testing.T) {
			mock := gock.New("https://akaa-baseurl-xxxxxxxxxxx-xxxxxxxxxxxxx.luna.akamaiapis.net/config-dns/v2/zones/example.com")
			mock.
				//Get("/config-dns/v2/zones/example.com/recordsets?types=A&showAll=true").
				Get("/config-dns/v2/zones/example.com").
				HeaderPresent("Authorization").
				Persist().
				Reply(200).
				SetHeader("Content-Type", "application/json").
				BodyString(test.responseBody)

			Init(config)
			fmt.Println("Record Type " + test.recordType)
			fmt.Println(test.expectedRecords.([]*RecordSetResponse)[0].Metadata.ShowAll)
			fmt.Println(test.expectedRecords.([]*RecordSetResponse)[0].Recordsets[0].Name)
			zone, err := GetZone("example.com")
			records, err := GetRecordList("example.com", "example.com", test.recordType) // (*RecordSetResponse, error)

			assert.NoError(t, err)

			assert.IsType(t, &RecordSetResponse{}, records)

			if test.expectedRecords != nil {
				switch test.recordType {
				case "A":
					fmt.Println(zone)
					fmt.Println(records)
					assert.Equal(t, len(test.expectedRecords.([]*RecordSetResponse)), len(records.Recordsets))
					assert.ObjectsAreEqual(test.expectedRecords, records)
					break
				case "AAAA":
					assert.Equal(t, len(test.expectedRecords.([]*RecordSetResponse)), len(records.Recordsets))
					assert.ObjectsAreEqual(test.expectedRecords, records)
					break
				case "AFSDB":
					assert.Equal(t, len(test.expectedRecords.([]*RecordSetResponse)), len(records.Recordsets))
					assert.ObjectsAreEqual(test.expectedRecords, records)
					break
				case "CNAME":
					assert.Equal(t, len(test.expectedRecords.([]*RecordSetResponse)), len(records.Recordsets))
					assert.ObjectsAreEqual(test.expectedRecords, records)
					break
				case "DNSKEY":
					assert.Equal(t, len(test.expectedRecords.([]*RecordSetResponse)), len(records.Recordsets))
					assert.ObjectsAreEqual(test.expectedRecords, records)
					break
				case "DS":
					assert.Equal(t, len(test.expectedRecords.([]*RecordSetResponse)), len(records.Recordsets))
					assert.ObjectsAreEqual(test.expectedRecords, records)
					break
				case "HINFO":
					assert.Equal(t, len(test.expectedRecords.([]*RecordSetResponse)), len(records.Recordsets))
					assert.ObjectsAreEqual(test.expectedRecords, records)
					break
				case "LOC":
					assert.Equal(t, len(test.expectedRecords.([]*RecordSetResponse)), len(records.Recordsets))
					assert.ObjectsAreEqual(test.expectedRecords, records)
					break
				case "MX":
					assert.Equal(t, len(test.expectedRecords.([]*RecordSetResponse)), len(records.Recordsets))
					assert.ObjectsAreEqual(test.expectedRecords, records)
					break
				case "NAPTR":
					assert.Equal(t, len(test.expectedRecords.([]*RecordSetResponse)), len(records.Recordsets))
					assert.ObjectsAreEqual(test.expectedRecords, records)
					break
				case "NS":
					assert.Equal(t, len(test.expectedRecords.([]*RecordSetResponse)), len(records.Recordsets))
					assert.ObjectsAreEqual(test.expectedRecords, records)
					break
				case "NSEC3":
					assert.Equal(t, len(test.expectedRecords.([]*RecordSetResponse)), len(records.Recordsets))
					assert.ObjectsAreEqual(test.expectedRecords, records)
					break
				case "NSEC3PARAM":
					assert.Equal(t, len(test.expectedRecords.([]*RecordSetResponse)), len(records.Recordsets))
					assert.ObjectsAreEqual(test.expectedRecords, records)
					break
				case "PTR":
					assert.Equal(t, len(test.expectedRecords.([]*RecordSetResponse)), len(records.Recordsets))
					assert.ObjectsAreEqual(test.expectedRecords, records)
					break
				case "RP":
					assert.Equal(t, len(test.expectedRecords.([]*RecordSetResponse)), len(records.Recordsets))
					assert.ObjectsAreEqual(test.expectedRecords, records)
					break
				case "RRSIG":
					assert.Equal(t, len(test.expectedRecords.([]*RecordSetResponse)), len(records.Recordsets))
					assert.ObjectsAreEqual(test.expectedRecords, records)
					break
				case "SPF":
					assert.Equal(t, len(test.expectedRecords.([]*RecordSetResponse)), len(records.Recordsets))
					assert.ObjectsAreEqual(test.expectedRecords, records)
					break
				case "SRV":
					assert.Equal(t, len(test.expectedRecords.([]*RecordSetResponse)), len(records.Recordsets))
					assert.ObjectsAreEqual(test.expectedRecords, records)
					break
				case "SSHFP":
					assert.Equal(t, len(test.expectedRecords.([]*RecordSetResponse)), len(records.Recordsets))
					assert.ObjectsAreEqual(test.expectedRecords, records)
					break
				case "TXT":
					assert.Equal(t, len(test.expectedRecords.([]*RecordSetResponse)), len(records.Recordsets))
					assert.ObjectsAreEqual(test.expectedRecords, records)
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
	expectedRecord  RecordSetResponse //ARecord
	expectedRecords interface{}
}

func testGetZoneCompleteProvider() recordTests {
	return recordTests{
		{
			name:       "A Records",
			recordType: "A",
			responseBody: `{
    "metadata": {
        "showAll": true,
        "totalElements": 1
    },
    "recordsets": [
        {
            "name": "arecord",
            "type": "A",
            "ttl": 3700,
            "rdata": [
                "10.0.0.2",
                "10.0.0.3",
								"1.2.3.9"
            ]
        }
    ]
}`,
			expectedType: []*RecordSetResponse{},
			expectedRecords: []*RecordSetResponse{
				&RecordSetResponse{
					//Active: true,
					Metadata:   MetadataH{ShowAll: true, TotalElements: 1},
					Recordsets: []Recordset{Recordset{Name: "arecord", Type: "A", TTL: 3600, Rdata: []string{"1.2.3.4", "1.2.3.5"}}},
					//Recordset {Target: []string{"1.2.3.4","1.2.3.5"}},
					//Recordset {TTL:    3600},
				},
				/*	&RecordSetResponse{
						//Active: true,
						//Name:   "origin",
						//Metadata {},
						Metadata: MetadataH {ShowAll: true, TotalElements: 1},
						Recordsets: []Recordset { Recordset{Name: "arecord", Type: "A", TTL:    3600, Rdata: []string{"1.2.3.9"} }},
						//Recordsets {Name:   "origin"},
						//Recordsets {Target: []string{"1.2.3.9"}},
					 // Recordsets {	TTL:    3600},
					},*/
			},
		},
		{
			name:       "AAAA Records",
			recordType: "AAAA",
			responseBody: `{
    "metadata": {
        "showAll": true,
        "totalElements": 1
    },
    "recordsets": [
        {
            "name": "ipv6record.akavaiodeveloper.net",
            "type": "AAAA",
            "ttl": 3600,
            "rdata": [
                "2001:db8:0:0:0:ff00:42:8329"
            ]
        }
    ]
}`,
			expectedType: []*RecordSetResponse{},
			expectedRecords: []*RecordSetResponse{
				&RecordSetResponse{
					//Active: true,
					//Name:   "ipv6record",
					//Target: []string{"2001:0db8::ff00:0042:8329"},
					//TTL:    3600,
					Metadata:   MetadataH{ShowAll: true, TotalElements: 1},
					Recordsets: []Recordset{Recordset{Name: "ipv6record.akavaiodeveloper.net", Type: "AAAA", TTL: 3600, Rdata: []string{"2001:db8:0:0:0:ff00:42:8329"}}},
				},
			},
		},
		{
			name:       "AFSDB Records",
			recordType: "AFSDB",
			responseBody: `{
    "metadata": {
        "showAll": true,
        "totalElements": 1
    },
    "recordsets": [
        {
            "name": "afsdb.akavaiodeveloper.net",
            "type": "AFSDB",
            "ttl": 3600,
            "rdata": [
                "4 example.com."
            ]
        }
    ]
}`,
			expectedType: []*RecordSetResponse{},
			expectedRecords: []*RecordSetResponse{
				&RecordSetResponse{
					/*Active:  true,
					Name:    "afsdb",
					Subtype: 1,
					Target:  []string{"example.com."},
					TTL:     7200,*/
					Metadata:   MetadataH{ShowAll: true, TotalElements: 1},
					Recordsets: []Recordset{Recordset{Name: "afsdb.akavaiodeveloper.net", Type: "AFSDB", TTL: 3600, Rdata: []string{"4 example.com."}}},
				},
			},
		},
		{
			name:       "CNAME Records",
			recordType: "CNAME",
			responseBody: `{
    "metadata": {
        "showAll": true,
        "totalElements": 1
    },
    "recordsets": [
        {
            "name": "www.akavaiodeveloper.net",
            "type": "CNAME",
            "ttl": 300,
            "rdata": [
                "api.akavaiodeveloper.net."
            ]
        }
    ]
}`,
			expectedType: []*RecordSetResponse{},
			expectedRecords: []*RecordSetResponse{
				&RecordSetResponse{
					/*Active: true,
					Name:   "redirect",
					Target: []string{"arecord.example.com."},
					TTL:    3600,*/
					Metadata:   MetadataH{ShowAll: true, TotalElements: 1},
					Recordsets: []Recordset{Recordset{Name: "www.akavaiodeveloper.net", Type: "CNAME", TTL: 300, Rdata: []string{"api.akavaiodeveloper.net."}}},
				},
			},
		},
		{
			name:       "DNSKEY Records",
			recordType: "DNSKEY",
			responseBody: `{
    "metadata": {
        "showAll": true,
        "totalElements": 1
    },
    "recordsets": [
        {
            "name": "dnskey.akavaiodeveloper.net",
            "type": "DNSKEY",
            "ttl": 7200,
            "rdata": [
                "257 7 3 Av//0/goGKPtaa28nQvPoUwVP++/i/0hC+1CrmQkuuKtQt98WObuv7q8iQ=="
            ]
        }
    ]
}`,
			expectedType: []*RecordSetResponse{},
			expectedRecords: []*RecordSetResponse{
				&RecordSetResponse{
					/*	Active:    true,
						Algorithm: 3,
						Flags:     257,
						Key:       "Av//0/goGKPtaa28nQvPoUwVQ ... i/0hC+1CrmQkuuKtQt98WObuv7q8iQ==",
						Name:      "dnskey",
						Protocol:  7,
						TTL:       7200,*/
					Metadata:   MetadataH{ShowAll: true, TotalElements: 1},
					Recordsets: []Recordset{Recordset{Name: "dnskey.akavaiodeveloper.net", Type: "DNSKEY", TTL: 7200, Rdata: []string{"257 7 3 Av//0/goGKPtaa28nQvPoUwVP++/i/0hC+1CrmQkuuKtQt98WObuv7q8iQ=="}}},
				},
			},
		},
		{
			name:       "DS Records",
			recordType: "DS",
			responseBody: `{
    "metadata": {
        "showAll": true,
        "totalElements": 1
    },
    "recordsets": [
        {
            "name": "ds.akavaiodeveloper.net",
            "type": "DS",
            "ttl": 7200,
            "rdata": [
                "30336 1 7 909FF0B4DD66F91F56524C4F968D13083BE42380"
            ]
        }
    ]
}`,
			expectedType: []*RecordSetResponse{},
			expectedRecords: []*RecordSetResponse{
				&RecordSetResponse{
					/*	Active:     true,
						Algorithm:  7,
						Digest:     "909FF0B4DD66F91F56524C4F968D13083BE42380",
						DigestType: 1,
						Keytag:     30336,
						Name:       "ds",
						TTL:        7200,*/
					Metadata:   MetadataH{ShowAll: true, TotalElements: 1},
					Recordsets: []Recordset{Recordset{Name: "ds.akavaiodeveloper.net", Type: "DS", TTL: 7200, Rdata: []string{"30336 1 7 909FF0B4DD66F91F56524C4F968D13083BE42380"}}},
				},
			},
		},
		{
			name:       "HINFO Records",
			recordType: "HINFO",
			responseBody: `{
    "metadata": {
        "showAll": true,
        "totalElements": 1
    },
    "recordsets": [
        {
            "name": "hinfo.akavaiodeveloper.net",
            "type": "HINFO",
            "ttl": 7200,
            "rdata": [
                "\"INTEL-386\" \"Unix\""
            ]
        }
    ]
}`,
			expectedType: []*RecordSetResponse{},
			expectedRecords: []*RecordSetResponse{
				&RecordSetResponse{
					/*Active:   true,
					Hardware: "INTEL-386",
					Name:     "hinfo",
					Software: "UNIX",
					TTL:      7200,*/
					Metadata:   MetadataH{ShowAll: true, TotalElements: 1},
					Recordsets: []Recordset{Recordset{Name: "hinfo.akavaiodeveloper.net", Type: "HINFO", TTL: 7200, Rdata: []string{"\"INTEL-386\" \"Unix\""}}},
				},
			},
		},
		{
			name:       "LOC Records",
			recordType: "LOC",
			responseBody: `{
    "metadata": {
        "showAll": true,
        "totalElements": 1
    },
    "recordsets": [
        {
            "name": "location.akavaiodeveloper.net",
            "type": "LOC",
            "ttl": 7200,
            "rdata": [
                "51 30 12.748 N 0 7 39.611 W 0m 0m 0m 0m"
            ]
        }
    ]
}`,
			expectedType: []*RecordSetResponse{},
			expectedRecords: []*RecordSetResponse{
				&RecordSetResponse{
					/*Active: true,
					Name:   "location",
					Target: []string{"51 30 12.748 N 0 7 39.611 W 0.00m 0.00m 0.00m 0.00m"},
					TTL:    7200,*/
					Metadata:   MetadataH{ShowAll: true, TotalElements: 1},
					Recordsets: []Recordset{Recordset{Name: "afsdb.akavaiodeveloper.net", Type: "LOC", TTL: 7200, Rdata: []string{"51 30 12.748 N 0 7 39.611 W 0m 0m 0m 0m"}}},
				},
			},
		},
		{
			name:       "MX Records",
			recordType: "MX",
			responseBody: `{
    "metadata": {
        "showAll": true,
        "totalElements": 1
    },
    "recordsets": [
        {
            "name": "akavaiodeveloper.net",
            "type": "MX",
            "ttl": 300,
            "rdata": [
                "10 smtp-1.akavaiodeveloper.net.",
                "20 smtp-3.akavaiodeveloper.net.",
                "30 smtp-0.akavaiodeveloper.net."
            ]
        }
    ]
}`,
			expectedType: []*RecordSetResponse{},
			expectedRecords: []*RecordSetResponse{
				&RecordSetResponse{
					/*Active:   true,
					Name:     "four",
					Priority: 10,
					Target:   []string{"mx1.akamai.com."},
					TTL:      7200,*/
					Metadata:   MetadataH{ShowAll: true, TotalElements: 1},
					Recordsets: []Recordset{Recordset{Name: "akavaiodeveloper.net", Type: "AFSDB", TTL: 300, Rdata: []string{"10 smtp-1.akavaiodeveloper.net.", "20 smtp-3.akavaiodeveloper.net.", "30 smtp-0.akavaiodeveloper.net."}}},
				},
			},
		},
		{
			name:       "NAPTR Records",
			recordType: "NAPTR",
			responseBody: `{
    "metadata": {
        "showAll": true,
        "totalElements": 1
    },
    "recordsets": [
        {
            "name": "naptrrecord.akavaiodeveloper.net",
            "type": "NAPTR",
            "ttl": 3600,
            "rdata": [
                "0 10 \"S\" \"!^.*$!sip:customer-service@example.com!\" \".\" SIP+D2U."
            ]
        }
    ]
}`,
			expectedType: []*RecordSetResponse{},
			expectedRecords: []*RecordSetResponse{
				&RecordSetResponse{
					/*	Active:      true,
						FlagsNaptr:  "S",
						Name:        "naptrrecord",
						Order:       0,
						Preference:  10,
						Regexp:      "!^.*$!sip:customer-service@example.com!",
						Replacement: ".",
						Service:     "SIP+D2U",
						TTL:         3600,*/
					Metadata:   MetadataH{ShowAll: true, TotalElements: 1},
					Recordsets: []Recordset{Recordset{Name: "naptrrecord.akavaiodeveloper.net", Type: "NAPTR", TTL: 3600, Rdata: []string{"0 10 \"S\" \"!^.*$!sip:customer-service@example.com!\" \".\" SIP+D2U."}}},
				},
			},
		},
		{
			name:       "NS Records",
			recordType: "NS",
			responseBody: `{
    "metadata": {
        "showAll": true,
        "totalElements": 2
    },
    "recordsets": [
        {
            "name": "akavaiodeveloper.net",
            "type": "NS",
            "ttl": 86400,
            "rdata": [
                "a1-49.akam.net.",
                "a16-64.akam.net.",
                "a22-65.akam.net.",
                "a26-66.akam.net.",
                "a7-67.akam.net.",
                "a9-64.akam.net."
            ]
        },
        {
            "name": "ns.akavaiodeveloper.net",
            "type": "NS",
            "ttl": 300,
            "rdata": [
                "use4.akam.net."
            ]
        }
    ]
}`,
			expectedType: []*RecordSetResponse{},
			expectedRecords: []*RecordSetResponse{
				&RecordSetResponse{
					/*		Active: true,
								Target: []string{"use4.akam.net."},
								TTL:    3600,
							},
							&RecordBody{
								Active: true,
								Target: []string{"us34.akam.net."},
								TTL:    3600,
							},
							&RecordBody{
								Active: true,
								Name:   "five",
								Target: []string{"use4.akam.net."},
								TTL:    172800,*/
					Metadata:   MetadataH{ShowAll: true, TotalElements: 1},
					Recordsets: []Recordset{Recordset{Name: "akavaiodeveloper.net", Type: "NS", TTL: 3600, Rdata: []string{"a1-49.akam.net.", "a16-64.akam.net.", "a22-65.akam.net.", "a26-66.akam.net.", "a7-67.akam.net.", "a9-64.akam.net."}}},
				},
			},
		},
		{
			name:       "NSEC3 Records",
			recordType: "NSEC3",
			responseBody: `{
    "metadata": {
        "showAll": true,
        "totalElements": 1
    },
    "recordsets": [
        {
            "name": "qdeo8lqu4l81uo67oolpo9h0nv9l13dh.akavaiodeveloper.net",
            "type": "NSEC3",
            "ttl": 3600,
            "rdata": [
                "0 1 1 EBD1E0942543A01B R2NUSMGFSEUHT195P59KOU2AI30JR90 CNAME RRSIG"
            ]
        }
    ]
}`,
			expectedType: []*RecordSetResponse{},
			expectedRecords: []*RecordSetResponse{
				&RecordSetResponse{
					/*Active:              true,
					Algorithm:           1,
					Flags:               0,
					Iterations:          1,
					Name:                "qdeo8lqu4l81uo67oolpo9h0nv9l13dh",
					NextHashedOwnerName: "R2NUSMGFSEUHT195P59KOU2AI30JR96P",
					Salt:                "EBD1E0942543A01B",
					TTL:                 7200,
					TypeBitmaps:         "CNAME RRSIG",*/
					Metadata:   MetadataH{ShowAll: true, TotalElements: 1},
					Recordsets: []Recordset{Recordset{Name: "qdeo8lqu4l81uo67oolpo9h0nv9l13dh.akavaiodeveloper.net", Type: "NSEC3", TTL: 3600, Rdata: []string{"0 1 1 EBD1E0942543A01B R2NUSMGFSEUHT195P59KOU2AI30JR90 CNAME RRSIG"}}},
				},
			},
		},
		{
			name:       "NSEC3PARAM Records",
			recordType: "NSEC3PARAM",
			responseBody: `{
    "metadata": {
        "showAll": true,
        "totalElements": 1
    },
    "recordsets": [
        {
            "name": "qnsec3param.akavaiodeveloper.net",
            "type": "NSEC3PARAM",
            "ttl": 3600,
            "rdata": [
                "0 1 1 EBD1E0942543A01B"
            ]
        }
    ]
}`,
			expectedType: []*RecordSetResponse{},
			expectedRecords: []*RecordSetResponse{
				&RecordSetResponse{
					/*	Active:     true,
						Algorithm:  1,
						Flags:      0,
						Iterations: 1,
						Name:       "qnsec3param",
						Salt:       "EBD1E0942543A01B",
						TTL:        7200,*/
					Metadata:   MetadataH{ShowAll: true, TotalElements: 1},
					Recordsets: []Recordset{Recordset{Name: "qnsec3param.akavaiodeveloper.net", Type: "NSEC3PARAM", TTL: 3600, Rdata: []string{"0 1 1 EBD1E0942543A01B"}}},
				},
			},
		},
		{
			name:       "PTR Records",
			recordType: "PTR",
			responseBody: `{
    "metadata": {
        "showAll": true,
        "totalElements": 2
    },
    "recordsets": [
        {
            "name": "ptr.akavaiodeveloper.net",
            "type": "PTR",
            "ttl": 300,
            "rdata": [
                "ptr.akavaiodeveloper.net."
            ]
        },
        {
            "name": "spf.akavaiodeveloper.net",
            "type": "PTR",
            "ttl": 7200,
            "rdata": [
                "v=spf."
            ]
        }
    ]
}`,
			expectedType: []*RecordSetResponse{},
			expectedRecords: []*RecordSetResponse{
				&RecordSetResponse{
					/*Active: true,
					Name:   "ptr",
					Target: []string{"ptr.example.com."},
					TTL:    7200,*/
					Metadata:   MetadataH{ShowAll: true, TotalElements: 1},
					Recordsets: []Recordset{Recordset{Name: "ptr.akavaiodeveloper.net", Type: "PTR", TTL: 7200, Rdata: []string{"v=spf."}}},
				},
			},
		},
		{
			name:       "RP Records",
			recordType: "RP",
			responseBody: `{
    "metadata": {
        "showAll": true,
        "totalElements": 1
    },
    "recordsets": [
        {
            "name": "rp.akavaiodeveloper.net",
            "type": "RP",
            "ttl": 7200,
            "rdata": [
                "admin.example.com. txt.example.com."
            ]
        }
    ]
}`,
			expectedType: []*RecordSetResponse{},
			expectedRecords: []*RecordSetResponse{
				&RecordSetResponse{
					/*Active:  true,
					Mailbox: "admin.example.com.",
					Name:    "rp",
					TTL:     7200,
					Txt:     "txt.example.com.",
					Metadata: MetadataH {ShowAll: true, TotalElements: 1},*/
					Recordsets: []Recordset{Recordset{Name: "rp.akavaiodeveloper.net", Type: "RP", TTL: 7200, Rdata: []string{"admin.example.com. txt.example.com."}}},
				},
			},
		},
		{
			name:       "RRSIG Records",
			recordType: "RRSIG",
			responseBody: `{
    "metadata": {
        "showAll": true,
        "totalElements": 1
    },
    "recordsets": [
{
            "name": "arecord",
            "ttl": 7200,
            "rdata": [
                  "A 7 3 3600 20120318104101 20120315094101 63761 3 toCy19QnAb86vRlQjf5 ... z1doJdHEr8PiI+Is9Eafxh+4Idcw8Ysv example.com."
					]
				}
                            ]
			}`,
			expectedType: []*RecordSetResponse{},
			expectedRecords: []*RecordSetResponse{
				&RecordSetResponse{
					/*	Active:      true,
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
						TypeCovered: "A",*/
					Metadata:   MetadataH{ShowAll: true, TotalElements: 1},
					Recordsets: []Recordset{Recordset{Name: "arecord", Type: "RRSIG", TTL: 7200, Rdata: []string{"A 7 3 3600 20120318104101 20120315094101 63761 3 toCy19QnAb86vRlQjf5 ... z1doJdHEr8PiI+Is9Eafxh+4Idcw8Ysv example.com."}}},
				},
			},
		},
		{
			name:       "SPF Records",
			recordType: "SPF",
			responseBody: `{
    "metadata": {
        "showAll": true,
        "totalElements": 1
    },
    "recordsets": [
        {
            "name": "spf.akavaiodeveloper.net",
            "type": "SPF",
            "ttl": 7200,
            "rdata": [
                "\"v=spf.\""
            ]
        }
    ]
}`,
			expectedType: []*RecordSetResponse{},
			expectedRecords: []*RecordSetResponse{
				&RecordSetResponse{
					/*	Active: true,
						Name:   "spf",
						Target: []string{"v=spf a -all"},
						TTL:    7200,*/
					Metadata:   MetadataH{ShowAll: true, TotalElements: 1},
					Recordsets: []Recordset{Recordset{Name: "afsdb.akavaiodeveloper.net", Type: "SPF", TTL: 7200, Rdata: []string{"\"v=spf.\""}}},
				},
			},
		},
		{
			name:       "SSHFP Records",
			recordType: "SSHFP",
			responseBody: `{
    "metadata": {
        "showAll": true,
        "totalElements": 1
    },
    "recordsets": [
        {
            "name": "sshfp.akavaiodeveloper.net",
            "type": "SSHFP",
            "ttl": 7200,
            "rdata": [
                "2 1 123456789ABCDEF67890123456789ABCDEF67890"
            ]
        }
    ]
}`,
			expectedType: []*RecordSetResponse{},
			expectedRecords: []*RecordSetResponse{
				&RecordSetResponse{
					/*	Active:          true,
						Algorithm:       2,
						Fingerprint:     "123456789ABCDEF67890123456789ABCDEF67890",
						FingerprintType: 1,
						Name:            "host",
						TTL:             3600,*/
					Metadata:   MetadataH{ShowAll: true, TotalElements: 1},
					Recordsets: []Recordset{Recordset{Name: "afsdb.akavaiodeveloper.net", Type: "SSHFP", TTL: 3600, Rdata: []string{"2 1 123456789ABCDEF67890123456789ABCDEF67890"}}},
				},
			},
		},
		{
			name:       "SRV Records",
			recordType: "SRV",
			responseBody: `{
    "metadata": {
        "showAll": true,
        "totalElements": 1
    },
    "recordsets": [
        {
            "name": "srv.akavaiodeveloper.net",
            "type": "SRV",
            "ttl": 7200,
            "rdata": [
                "0 522 10 target.akavaiodeveloper.net."
            ]
        }
    ]
}`,
			expectedType: []*RecordSetResponse{},
			expectedRecords: []*RecordSetResponse{
				&RecordSetResponse{
					/*Active:   true,
					Name:     "srv",
					Port:     522,
					Priority: 10,
					Target:   []string{"target.akamai.com."},
					TTL:      7200,
					Weight:   0,*/
					Metadata:   MetadataH{ShowAll: true, TotalElements: 1},
					Recordsets: []Recordset{Recordset{Name: "srv.akavaiodeveloper.net", Type: "SRV", TTL: 7200, Rdata: []string{"0 522 10 target.akavaiodeveloper.net."}}},
				},
			},
		},
		{
			name:       "TXT Records",
			recordType: "TXT",
			responseBody: `{
    "metadata": {
        "showAll": true,
        "totalElements": 1
    },
    "recordsets": [
        {
            "name": "text.akavaiodeveloper.net",
            "type": "TXT",
            "ttl": 7200,
            "rdata": [
                "\"Hello\" \"world\" \"this\" \"is\" \"text\""
            ]
        }
    ]
}`,
			expectedType: []*RecordSetResponse{},
			expectedRecords: []*RecordSetResponse{
				&RecordSetResponse{
					/*Active: true,
					Name:   "text",
					Target: []string{"Hello world!"},
					TTL:    7200,*/
					Metadata:   MetadataH{ShowAll: true, TotalElements: 1},
					Recordsets: []Recordset{Recordset{Name: "text.akavaiodeveloper.net", Type: "TXT", TTL: 7200, Rdata: []string{"Hello world!"}}},
				},
			},
		},
	}
}
