package dns

import (
	"fmt"
	"regexp"
	"testing"

	"github.com/akamai/AkamaiOPEN-edgegrid-golang/v12/pkg/dns"
	"github.com/akamai/terraform-provider-akamai/v8/internal/test"
	"github.com/akamai/terraform-provider-akamai/v8/pkg/common/testutils"
	"github.com/hashicorp/terraform-plugin-testing/helper/resource"
	"github.com/stretchr/testify/mock"
)

func TestDataZoneDnsSecStatus(t *testing.T) {
	anyContext := mock.AnythingOfType("*context.valueCtx")
	request := dns.GetZonesDNSSecStatusRequest{
		Zones: []string{"test.zone.net"},
	}

	tests := map[string]struct {
		givenTF            string
		init               func(mock *dns.Mock)
		expectedAttributes map[string]string
		expectedError      *regexp.Regexp
	}{
		"valid DNSSEC status returned": {
			givenTF: "valid.tf",
			init: func(m *dns.Mock) {
				m.On("GetZonesDNSSecStatus", anyContext, request).Return(
					&dns.GetZonesDNSSecStatusResponse{
						DNSSecStatuses: []dns.SecStatus{{
							Zone: "test.zone.net",
							CurrentRecords: dns.SecRecords{
								DNSKeyRecord:     "DUMMY_DNSKEY_RECORD",
								DSRecord:         "DUMMY_DS_RECORD",
								ExpectedTTL:      1234,
								LastModifiedDate: test.NewTimeFromString(t, "2024-05-28T06:58:26Z"),
							},
						}},
					}, nil)

			},
			expectedAttributes: map[string]string{
				"zone":                               "test.zone.net",
				"current_records.dnskey_record":      "DUMMY_DNSKEY_RECORD",
				"current_records.ds_record":          "DUMMY_DS_RECORD",
				"current_records.expected_ttl":       "1234",
				"current_records.last_modified_date": "2024-05-28T06:58:26Z",
				"new_records.%":                      "0",
				"alerts.#":                           "0",
			},
		},
		"alerts returned": {
			givenTF: "valid.tf",
			init: func(m *dns.Mock) {
				m.On("GetZonesDNSSecStatus", anyContext, request).Return(
					&dns.GetZonesDNSSecStatusResponse{
						DNSSecStatuses: []dns.SecStatus{{
							Zone: "test.zone.net",
							CurrentRecords: dns.SecRecords{
								DNSKeyRecord:     "DUMMY_DNSKEY_RECORD",
								DSRecord:         "DUMMY_DS_RECORD",
								ExpectedTTL:      1234,
								LastModifiedDate: test.NewTimeFromString(t, "2024-05-28T06:58:26Z"),
							},
							Alerts: []string{
								"PARENT_DS_MISSING",
								"OLD_DNSKEY",
								"INCOMPATIBLE_AUTHORITIES",
							},
						}},
					}, nil)
			},
			expectedAttributes: map[string]string{
				"zone":                               "test.zone.net",
				"current_records.dnskey_record":      "DUMMY_DNSKEY_RECORD",
				"current_records.ds_record":          "DUMMY_DS_RECORD",
				"current_records.expected_ttl":       "1234",
				"current_records.last_modified_date": "2024-05-28T06:58:26Z",
				"new_records.%":                      "0",
				"alerts.#":                           "3",
				"alerts.0":                           "INCOMPATIBLE_AUTHORITIES",
				"alerts.1":                           "OLD_DNSKEY",
				"alerts.2":                           "PARENT_DS_MISSING",
			},
		},
		"new records returned": {
			givenTF: "valid.tf",
			init: func(m *dns.Mock) {
				m.On("GetZonesDNSSecStatus", anyContext, request).Return(
					&dns.GetZonesDNSSecStatusResponse{
						DNSSecStatuses: []dns.SecStatus{{
							Zone: "test.zone.net",
							CurrentRecords: dns.SecRecords{
								DNSKeyRecord:     "DUMMY_DNSKEY_RECORD",
								DSRecord:         "DUMMY_DS_RECORD",
								ExpectedTTL:      1234,
								LastModifiedDate: test.NewTimeFromString(t, "2024-05-28T06:58:26Z"),
							},
							NewRecords: &dns.SecRecords{
								DNSKeyRecord:     "DUMMY_DNSKEY_RECORD_2",
								DSRecord:         "DUMMY_DS_RECORD_2",
								ExpectedTTL:      5678,
								LastModifiedDate: test.NewTimeFromString(t, "2024-05-31T13:27:55Z"),
							},
						}},
					}, nil)

			},
			expectedAttributes: map[string]string{
				"zone":                               "test.zone.net",
				"current_records.dnskey_record":      "DUMMY_DNSKEY_RECORD",
				"current_records.ds_record":          "DUMMY_DS_RECORD",
				"current_records.expected_ttl":       "1234",
				"current_records.last_modified_date": "2024-05-28T06:58:26Z",
				"new_records.dnskey_record":          "DUMMY_DNSKEY_RECORD_2",
				"new_records.ds_record":              "DUMMY_DS_RECORD_2",
				"new_records.expected_ttl":           "5678",
				"new_records.last_modified_date":     "2024-05-31T13:27:55Z",
				"alerts.#":                           "0",
			},
		},
		"no DNSSEC status returned": {
			givenTF: "valid.tf",
			init: func(m *dns.Mock) {
				m.On("GetZonesDNSSecStatus", anyContext, request).Return(
					&dns.GetZonesDNSSecStatusResponse{
						DNSSecStatuses: []dns.SecStatus{},
					}, nil)

			},
			expectedError: regexp.MustCompile("no DNSSEC status for zone: test.zone.net"),
		},
		"missing required argument zone": {
			givenTF:       "missing_zone_name.tf",
			expectedError: regexp.MustCompile(`The argument "zone" is required, but no definition was found.`),
		},
		"empty zone name": {
			givenTF:       "empty_zone_name.tf",
			expectedError: regexp.MustCompile(`Attribute zone string length must be at least 1, got: 0`),
		},
		"error response from api": {
			givenTF: "valid.tf",
			init: func(m *dns.Mock) {
				m.On("GetZonesDNSSecStatus", anyContext, request).Return(
					nil, fmt.Errorf("API error"))
			},
			expectedError: regexp.MustCompile("API error"),
		},
	}

	for name, test := range tests {
		t.Run(name, func(t *testing.T) {
			client := &dns.Mock{}
			if test.init != nil {
				test.init(client)
			}
			var checkFuncs []resource.TestCheckFunc
			for k, v := range test.expectedAttributes {
				checkFuncs = append(checkFuncs, resource.TestCheckResourceAttr("data.akamai_zone_dnssec_status.test", k, v))
			}

			useClient(client, func() {
				resource.Test(t, resource.TestCase{
					IsUnitTest:               true,
					ProtoV6ProviderFactories: testutils.NewProtoV6ProviderFactory(NewSubprovider()),
					Steps: []resource.TestStep{{
						Config:      testutils.LoadFixtureStringf(t, "testdata/TestDataZoneDnsSecStatus/%s", test.givenTF),
						Check:       resource.ComposeAggregateTestCheckFunc(checkFuncs...),
						ExpectError: test.expectedError,
					}},
				})
			})

			client.AssertExpectations(t)
		})
	}
}
