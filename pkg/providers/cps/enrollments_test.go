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
	challenge1 := cps.Challenges{
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
		Challenges:         []cps.Challenges{challenge1},
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
