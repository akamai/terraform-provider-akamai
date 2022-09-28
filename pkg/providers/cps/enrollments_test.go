package cps

import (
	"testing"

	"github.com/akamai/AkamaiOPEN-edgegrid-golang/v3/pkg/cps"
	"github.com/stretchr/testify/assert"
)

func TestSplitChallenges(t *testing.T) {
	t.Run("Non empty challenges", func(t *testing.T) {
		challenges := mockDVArray()
		gotHTTPChallenge, gotDNSChallenge := splitChallenges(challenges)
		wantHTTPChallenge := []challengeHTTP{
			{
				"full_path":     "http://TestFullPath",
				"response_body": "TestResponseBody",
				"domain":        "TestDomain",
			},
		}
		wantDNSChallenge := []challengeDNS{
			{
				"full_path":     "TestFullPath",
				"response_body": "TestResponseBody",
				"domain":        "TestDomain",
			},
		}
		assert.Equal(t, wantHTTPChallenge, gotHTTPChallenge)
		assert.Equal(t, wantDNSChallenge, gotDNSChallenge)
	})

	t.Run("Empty challenges", func(t *testing.T) {
		challenges := mockEmptyDVArray()
		gotHTTPChallenge, gotDNSChallenge := splitChallenges(challenges)
		wantDNSChallenge := make([]challengeDNS, 0)
		wantHTTPChallenge := make([]challengeHTTP, 0)
		assert.Equal(t, wantHTTPChallenge, gotHTTPChallenge)
		assert.Equal(t, wantDNSChallenge, gotDNSChallenge)
	})
}

func TestNewChallenge(t *testing.T) {
	challenge1 := cps.Challenge{
		Error:             "",
		FullPath:          "http://TestFullPath",
		RedirectFullPath:  "http://TestRedirectFullPath.com",
		ResponseBody:      "TestResponseBody",
		Status:            "pending",
		Token:             "TestToken123",
		Type:              "http-01",
		ValidationRecords: nil,
	}
	dv := cps.DV{
		Challenges:         []cps.Challenge{challenge1},
		Domain:             "TestDomain",
		Error:              "The domain TestDomain is not ready for HTTP validation.",
		Expires:            "2022-07-25T10:17:44Z",
		RequestTimestamp:   "2022-07-18T10:17:44Z",
		Status:             "Awaiting user",
		ValidatedTimestamp: "2022-07-19T09:35:29Z",
		ValidationStatus:   "DATA_NOT_READY",
	}

	gotChallenge := newChallenge(&challenge1, &dv)
	wantChallenge := challenge{
		"full_path":     "http://TestFullPath",
		"response_body": "TestResponseBody",
		"domain":        "TestDomain",
	}
	assert.Equal(t, wantChallenge, gotChallenge)
}

func TestConvertWarnings(t *testing.T) {
	tests := map[string]struct {
		warnings         string
		expectedWarnings []string
		err              string
	}{
		"no warnings": {
			warnings:         ``,
			expectedWarnings: nil,
		},
		"all warnings can be converted": {
			warnings:         "Certificate data is blank or missing for 'RSA'.\nThe 'ECDSA' certificate is set to expire in [2] years, [3] months. The certificate has a validity period of greater than 397 days. This certificate will not be accepted by all major browsers for SSL/TLS connections. Please work with your Certificate Authority to reissue the certificate with an acceptable lifetime.\nThe trust chain is empty and the end-entity certificate may have been signed by a non-standard root certificate.",
			expectedWarnings: []string{"THIRD_PARTY_CERTIFICATE_DATA_BLANK_OR_MISSING", "CERTIFICATE_EXPIRATION_DATE_BEYOND_MAX_DAYS", "TRUST_CHAIN_EMPTY_AND_CERTIFICATE_SIGNED_BY_NON_STANDARD_ROOT"},
		},
		"warnings with new line": {
			warnings:         "Certificate data is blank or missing for 'RSA'.\nExtra certificates were found in the chain and are being removed.\ntrustChainData",
			expectedWarnings: []string{"THIRD_PARTY_CERTIFICATE_DATA_BLANK_OR_MISSING", "EXTRA_CERT_IN_TRUST_CHAIN"},
		},
		"warnings with multiline": {
			warnings:         "Expected to find trust chain:\n    <expectedName>\n    <expectedDescription>\n  Instead found:\n    <actualName>\n    <actualDescription>\nExtra certificates were found in the chain and are being removed.\ntrustChainData",
			expectedWarnings: []string{"NAMED_TRUST_CHAIN_MISMATCH", "EXTRA_CERT_IN_TRUST_CHAIN"},
		},
		"warning text which is matching to two keys": {
			warnings:         "Certificate data is blank or missing for 'RSA'.\nThe trust chain terminates with a non-standard root certificate.\ntrustChainData",
			expectedWarnings: []string{"THIRD_PARTY_CERTIFICATE_DATA_BLANK_OR_MISSING", "TRUST_CHAIN_TERMINATES_WITH_NON_STANDARD_CERTIFICATE_DETAILED"},
		},
		"warning text which is matching to two keys 2": {
			warnings:         "Trust chain is empty.\nCertificate has a null issuer",
			expectedWarnings: []string{"TRUST_CHAIN_NULL_OR_EMPTY", "CERTIFICATE_HAS_NULL_ISSUER"},
		},
		"unknown warnings": {
			warnings: "unknown 1.\nCertificate data is blank or missing for 'RSA'.\nunknown 2.",
			err:      `received warning(s) does not match any known warning: 'unknown 1.', 'unknown 2.'`,
		},
		"warning with several newlines": {
			warnings:         "Crossed signed roots found in the trust chain.\nsomepartoftrustchain\nsecondpartoftrustchain\nendoftrust chain\nDNS Text Parse Exception when processing this \nis unknown\n that is similar",
			expectedWarnings: []string{"CROSS_SIGNED_ROOT_IN_TRUST_CHAIN", "DNS_TEXT_PARSE"},
		},
	}

	for name, test := range tests {
		t.Run(name, func(t *testing.T) {
			res, err := convertWarnings(test.warnings)
			if err != nil {
				assert.Equal(t, test.err, err.Error())
			} else {
				assert.Equal(t, test.expectedWarnings, res)
			}
		})
	}
}

func TestCanApproveWarnings(t *testing.T) {
	tests := map[string]struct {
		warnings            string
		autoApproveWarnings []string
		can                 bool
		err                 string
	}{
		"all warnings can be approved": {
			warnings:            "Certificate data is blank or missing for 'RSA'.\nThe 'ECDSA' certificate is set to expire in [2] years, [3] months. The certificate has a validity period of greater than 397 days. This certificate will not be accepted by all major browsers for SSL/TLS connections. Please work with your Certificate Authority to reissue the certificate with an acceptable lifetime.\nThe trust chain is empty and the end-entity certificate may have been signed by a non-standard root certificate.",
			autoApproveWarnings: []string{"THIRD_PARTY_CERTIFICATE_DATA_BLANK_OR_MISSING", "CERTIFICATE_EXPIRATION_DATE_BEYOND_MAX_DAYS", "TRUST_CHAIN_EMPTY_AND_CERTIFICATE_SIGNED_BY_NON_STANDARD_ROOT"},
			can:                 true,
		},
		"no warnings returned": {
			warnings:            "",
			autoApproveWarnings: []string{"THIRD_PARTY_CERTIFICATE_DATA_BLANK_OR_MISSING", "CERTIFICATE_EXPIRATION_DATE_BEYOND_MAX_DAYS", "TRUST_CHAIN_EMPTY_AND_CERTIFICATE_SIGNED_BY_NON_STANDARD_ROOT"},
			can:                 true,
		},
		"none can be auto-approved": {
			warnings:            "Certificate data is blank or missing for 'RSA'.\nThe 'ECDSA' certificate is set to expire in [2] years, [3] months. The certificate has a validity period of greater than 397 days. This certificate will not be accepted by all major browsers for SSL/TLS connections. Please work with your Certificate Authority to reissue the certificate with an acceptable lifetime.\nThe trust chain is empty and the end-entity certificate may have been signed by a non-standard root certificate.",
			autoApproveWarnings: []string{},
			err:                 "warnings cannot be approved: THIRD_PARTY_CERTIFICATE_DATA_BLANK_OR_MISSING, CERTIFICATE_EXPIRATION_DATE_BEYOND_MAX_DAYS, TRUST_CHAIN_EMPTY_AND_CERTIFICATE_SIGNED_BY_NON_STANDARD_ROOT",
			can:                 false,
		},
		"none warning can be auto-approved and none provided": {
			warnings:            "",
			autoApproveWarnings: []string{},
			err:                 "warnings cannot be approved: THIRD_PARTY_CERTIFICATE_DATA_BLANK_OR_MISSING, CERTIFICATE_EXPIRATION_DATE_BEYOND_MAX_DAYS, TRUST_CHAIN_EMPTY_AND_CERTIFICATE_SIGNED_BY_NON_STANDARD_ROOT",
			can:                 true,
		},
		"not all warnings can be approved": {
			warnings:            "Certificate data is blank or missing for 'RSA'.\nThe 'ECDSA' certificate is set to expire in [2] years, [3] months. The certificate has a validity period of greater than 397 days. This certificate will not be accepted by all major browsers for SSL/TLS connections. Please work with your Certificate Authority to reissue the certificate with an acceptable lifetime.\nThe trust chain is empty and the end-entity certificate may have been signed by a non-standard root certificate.",
			autoApproveWarnings: []string{"CERTIFICATE_EXPIRATION_DATE_BEYOND_MAX_DAYS", "TRUST_CHAIN_EMPTY_AND_CERTIFICATE_SIGNED_BY_NON_STANDARD_ROOT"},
			err:                 "warnings cannot be approved: THIRD_PARTY_CERTIFICATE_DATA_BLANK_OR_MISSING",
			can:                 false,
		},
		"unknown warning": {
			warnings: "unknown 1.\nCertificate data is blank or missing for 'RSA'.\nunknown 2.",
			err:      `received warning(s) does not match any known warning: 'unknown 1.', 'unknown 2.'`,
			can:      false,
		},
	}

	for name, test := range tests {
		t.Run(name, func(t *testing.T) {
			res, err := canApproveWarnings(test.autoApproveWarnings, test.warnings)
			if err != nil {
				assert.Equal(t, test.err, err.Error())
			} else {
				assert.Equal(t, test.can, res)
			}
		})
	}
}
