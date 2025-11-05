package akamai

import (
	"context"
	"crypto/tls"
	"crypto/x509"
	"errors"
	"net/http"
	"net/http/httptest"
	"net/url"
	"strings"
	"testing"
	"time"

	"github.com/akamai/AkamaiOPEN-edgegrid-golang/v12/pkg/ccm"
	"github.com/akamai/AkamaiOPEN-edgegrid-golang/v12/pkg/edgegrid"
	"github.com/akamai/AkamaiOPEN-edgegrid-golang/v12/pkg/papi"
	"github.com/akamai/AkamaiOPEN-edgegrid-golang/v12/pkg/session"
	"github.com/akamai/terraform-provider-akamai/v9/internal/test"
	"github.com/akamai/terraform-provider-akamai/v9/pkg/retryablehttp"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

func Test_validateRetryConfiguration(t *testing.T) {

	tests := map[string]struct {
		name    string
		args    contextConfig
		wantErr bool
		errMsg  string
	}{
		"valid values": {
			args: contextConfig{
				retryMax:     0,
				retryWaitMin: 0,
				retryWaitMax: 0,
			},

			wantErr: false,
		},
		"invalid values - negative number of retries": {
			args: contextConfig{
				retryMax:     -3,
				retryWaitMin: 0,
				retryWaitMax: 0,
			},

			wantErr: true,
			errMsg:  "wrong retry values: maximum number of retries (-3), minimum retry wait time (0s), maximum retry wait time (0s) cannot be negative",
		},
		"invalid values - negative min wait time": {
			args: contextConfig{
				retryMax:     0,
				retryWaitMin: -1,
				retryWaitMax: 0,
			},

			wantErr: true,
			errMsg:  "wrong retry values: maximum number of retries (0), minimum retry wait time (-1ns), maximum retry wait time (0s) cannot be negative",
		},
		"invalid values - negative max wait time": {
			args: contextConfig{
				retryMax:     0,
				retryWaitMin: 0,
				retryWaitMax: -1,
			},
			wantErr: true,
			errMsg:  "wrong retry values: maximum number of retries (0), minimum retry wait time (0s), maximum retry wait time (-1ns) cannot be negative",
		},
		"invalid values - min wait time cannot be higher than max wait time": {
			args: contextConfig{
				retryMax:     0,
				retryWaitMin: 1,
				retryWaitMax: 0,
			},
			wantErr: true,
			errMsg:  "wrong retry values: maximum retry wait time (0s) cannot be lower than minimum retry wait time (1ns)",
		},
		"invalid values - too many retries": {
			args: contextConfig{
				retryMax:     51,
				retryWaitMin: 0,
				retryWaitMax: 0,
			},
			wantErr: true,
			errMsg:  "wrong retry values: too many retries, maximum number of retries (51) cannot be higher than 50",
		},
		"invalid values - retry time too long (retryWaitMin)": {
			args: contextConfig{
				retryMax:     0,
				retryWaitMin: time.Hour * 25,
				retryWaitMax: time.Hour * 26, // needs to be higher than retryWaitMin
			},
			wantErr: true,
			errMsg:  "wrong retry values: retry wait time too long, minimum retry wait time (25h0m0s) cannot be higher than 24h0m0s or maximum retry wait time (26h0m0s) cannot be higher than 24h0m0s",
		},
		"invalid values - retry time too long (retryWaitMax)": {
			args: contextConfig{
				retryMax:     0,
				retryWaitMin: 1,
				retryWaitMax: time.Hour * 25,
			},
			wantErr: true,
			errMsg:  "wrong retry values: retry wait time too long, minimum retry wait time (1ns) cannot be higher than 24h0m0s or maximum retry wait time (25h0m0s) cannot be higher than 24h0m0s",
		},
	}
	for name, tt := range tests {
		t.Run(name, func(t *testing.T) {
			err := validateRetryConfiguration(tt.args)
			if tt.wantErr {
				assert.EqualError(t, err, tt.errMsg)
			} else {
				assert.NoError(t, err)
			}
		})
	}
}

func newRequest(t *testing.T, method, url string) *http.Request {
	r, err := http.NewRequest(method, url, nil)
	assert.NoError(t, err)
	return r
}

func TestOverrideRetryPolicy(t *testing.T) {
	basePolicy := func(_ context.Context, _ *http.Response, _ error) (bool, error) {
		return false, errors.New("base policy: dummy, not implemented")
	}
	policy := overrideRetryPolicy(basePolicy)

	tests := map[string]struct {
		ctx            context.Context
		resp           *http.Response
		err            error
		expectedResult bool
		expectedError  string
	}{
		"should retry for PAPI GET with status 429": {
			ctx: context.Background(),
			resp: &http.Response{
				Request:    newRequest(t, http.MethodGet, "/papi/v1/sth"),
				StatusCode: http.StatusTooManyRequests,
			},
			expectedResult: true,
		},
		"should retry for PAPI POST with status 429": {
			ctx: context.Background(),
			resp: &http.Response{
				Request:    newRequest(t, http.MethodPost, "/papi/v1/sth"),
				StatusCode: http.StatusTooManyRequests,
			},
			expectedResult: true,
		},
		"should not retry for PAPI POST with other 4xx status": {
			ctx: context.Background(),
			resp: &http.Response{
				Request:    newRequest(t, http.MethodPost, "/papi/v1/sth"),
				StatusCode: http.StatusBadRequest,
			},
			expectedResult: false,
		},
		"should retry for CCM GET with status 429": {
			ctx: context.Background(),
			resp: &http.Response{
				Request:    newRequest(t, http.MethodGet, "/ccm/v1/sth"),
				StatusCode: http.StatusTooManyRequests,
			},
			expectedResult: true,
		},
		"should retry for CCM POST with status 429": {
			ctx: context.Background(),
			resp: &http.Response{
				Request:    newRequest(t, http.MethodPost, "/ccm/v1/sth"),
				StatusCode: http.StatusTooManyRequests,
			},
			expectedResult: true,
		},
		"should not retry for CCM POST with other 4xx status": {
			ctx: context.Background(),
			resp: &http.Response{
				Request:    newRequest(t, http.MethodPost, "/ccm/v1/sth"),
				StatusCode: http.StatusBadRequest,
			},
			expectedResult: false,
		},
		"should retry for GET with status 409 conflict": {
			ctx: context.Background(),
			resp: &http.Response{
				Request:    newRequest(t, http.MethodGet, "/papi/v1/sth"),
				StatusCode: http.StatusConflict,
			},
			expectedResult: true,
		},
		"should call base policy for other GETs": {
			ctx:           context.Background(),
			resp:          &http.Response{Request: newRequest(t, http.MethodGet, "/papi/v1/sth")},
			expectedError: "base policy: dummy, not implemented",
		},
		"should forward context error when present": {
			ctx: func() context.Context {
				ctx, cancel := context.WithCancel(context.Background())
				cancel()
				return ctx
			}(),
			resp:          &http.Response{Request: newRequest(t, http.MethodGet, "/papi/v1/sth")},
			expectedError: "context canceled",
		},
		"should not retry for POST": {
			ctx:            context.Background(),
			resp:           &http.Response{Request: newRequest(t, http.MethodPost, "/papi/v1/sth")},
			expectedResult: false,
		},
		"should not retry for PUT": {
			ctx:            context.Background(),
			resp:           &http.Response{Request: newRequest(t, http.MethodPut, "/papi/v1/sth")},
			expectedResult: false,
		},
		"should not retry for PATCH": {
			ctx:            context.Background(),
			resp:           &http.Response{Request: newRequest(t, http.MethodPatch, "/papi/v1/sth")},
			expectedResult: false,
		},
		"should not retry for HEAD": {
			ctx:            context.Background(),
			resp:           &http.Response{Request: newRequest(t, http.MethodHead, "/papi/v1/sth")},
			expectedResult: false,
		},
		"should not retry for DELETE": {
			ctx:            context.Background(),
			resp:           &http.Response{Request: newRequest(t, http.MethodDelete, "/papi/v1/sth")},
			expectedResult: false,
		},
		"should not retry for PATCH hostnames bucket if 429": {
			ctx: context.Background(),
			resp: &http.Response{
				Request:    newRequest(t, http.MethodPatch, "/papi/v1/properties/prp_111/hostnames"),
				StatusCode: http.StatusTooManyRequests,
			},
			expectedResult: false,
		},
		"should retry for PATCH hostnames bucket if method different than PATCH": {
			ctx: context.Background(),
			resp: &http.Response{
				Request:    newRequest(t, http.MethodGet, "/papi/v1/properties/prp_111/hostnames"),
				StatusCode: http.StatusTooManyRequests,
			},
			expectedResult: true,
		},
		"should retry for PATCH hostnames bucket if status different than 429": {
			ctx: context.Background(),
			resp: &http.Response{
				Request:    newRequest(t, http.MethodGet, "/papi/v1/properties/prp_111/hostnames"),
				StatusCode: http.StatusConflict,
			},
			expectedResult: true,
		},
	}
	for name, tst := range tests {
		t.Run(name, func(t *testing.T) {
			shouldRetry, err := policy(tst.ctx, tst.resp, tst.err)
			if len(tst.expectedError) > 0 {
				assert.ErrorContains(t, err, tst.expectedError)
			} else {
				assert.NoError(t, err)
				assert.Equal(t, tst.expectedResult, shouldRetry)
			}
		})
	}
}

func stat429ResponseWaiting(wait time.Duration, header string) *http.Response {
	res := http.Response{
		StatusCode: http.StatusTooManyRequests,
		Header:     http.Header{},
	}

	now := time.Now().UTC().Round(time.Second)
	date := strings.Replace(now.Format(time.RFC1123), "UTC", "GMT", 1)
	res.Header.Add("Date", date)
	if wait != 0 {
		// Add: allow to canonicalize to X-RateLimit-Next and Akamai-RateLimit-Next or the header won't be recognized
		res.Header.Add(header, now.Add(wait).Format(time.RFC3339Nano))
	}
	return &res
}

func Test_overrideBackoff(t *testing.T) {
	baseWait := time.Duration(24) * time.Hour
	baseBackoff := func(_, _ time.Duration, _ int, _ *http.Response) time.Duration {
		return baseWait
	}
	backoff := overrideBackoff(baseBackoff, nil)

	tests := map[string]struct {
		resp           *http.Response
		expectedResult time.Duration
	}{
		"correctly calculates backoff from X-RateLimit-Next": {
			resp:           stat429ResponseWaiting(time.Duration(5729)*time.Millisecond, "X-RateLimit-Next"),
			expectedResult: time.Duration(5729) * time.Millisecond,
		},
		"correctly calculates backoff from Akamai-RateLimit-Next": {
			resp:           stat429ResponseWaiting(time.Duration(5729)*time.Millisecond, "Akamai-RateLimit-Next"),
			expectedResult: time.Duration(5729) * time.Millisecond,
		},
		"falls back for next in the past": {
			resp:           stat429ResponseWaiting(-time.Duration(5729)*time.Millisecond, "X-RateLimit-Next"),
			expectedResult: baseWait,
		},
		"falls back for no X-RateLimit-Next header": {
			resp:           stat429ResponseWaiting(0, "X-RateLimit-Next"),
			expectedResult: baseWait,
		},
		"falls back for invalid X-RateLimit-Next header": {
			resp: func() *http.Response {
				r := stat429ResponseWaiting(time.Duration(5729)*time.Millisecond, "X-RateLimit-Next")
				r.Header.Set("X-RateLimit-Next", "2024-07-01T14:32:28.645???")
				return r
			}(),
			expectedResult: baseWait,
		},
		"falls back for no Date header": {
			resp: func() *http.Response {
				r := stat429ResponseWaiting(time.Duration(5729)*time.Millisecond, "X-RateLimit-Next")
				r.Header.Del("Date")
				return r
			}(),
			expectedResult: baseWait,
		},
		"falls back for invalid Date header": {
			resp: func() *http.Response {
				r := stat429ResponseWaiting(time.Duration(5729)*time.Millisecond, "X-RateLimit-Next")
				r.Header.Set("Date", "Mon, 01 Jul 2024 99:99:99 GMT")
				return r
			}(),
			expectedResult: baseWait,
		},
	}
	for name, tst := range tests {
		t.Run(name, func(t *testing.T) {
			wait := backoff(1, 30, 1, tst.resp)
			assert.Equal(t, tst.expectedResult, wait)
		})
	}
}

func mockSession(t *testing.T, mockServer *httptest.Server) session.Session {
	serverURL, err := url.Parse(mockServer.URL)
	require.NoError(t, err)
	config := edgegrid.Config{Host: serverURL.Host}

	meta, err := configureContext(contextConfig{
		userAgent:      "terraform-provider-akamai",
		edgegridConfig: &config,
		ctx:            context.Background(),
	})
	assert.NoError(t, err)

	certPool := x509.NewCertPool()
	certPool.AddCert(mockServer.Certificate())
	rt := meta.Session().Client().Transport.(*retryablehttp.RoundTripper)
	transport := rt.Client.HTTPClient.Transport.(*http.Transport)
	transport.TLSClientConfig = &tls.Config{
		RootCAs: certPool,
	}

	return meta.Session()
}

func TestXRateLimitGet(t *testing.T) {
	xrlHandler := test.RateLimitHTTPHandler{
		T:           t,
		SuccessCode: http.StatusOK,
		SuccessBody: `
		{
			"properties": {
				"items": [
					{
						"accountId": "dummy_account_id",
						"contractId": "ctr_test1",
						"groupId": "grp_test1",
						"propertyId": "prp_test1",
						"propertyName": "my_property",
						"latestVersion": 1,
						"stagingVersion": null,
						"productionVersion": null,
						"assetId": "12345678"
					}
				]
			}
		}`,
	}

	mockServer := httptest.NewTLSServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		assert.Equal(t, "/papi/v1/properties/prp_test1?contractId=ctr_test1&groupId=grp_test1", r.URL.String())
		assert.Equal(t, http.MethodGet, r.Method)
		xrlHandler.ServeHTTP(w, r, "X-RateLimit-Next")
	}))
	defer mockServer.Close()

	client := papi.Client(mockSession(t, mockServer))
	result, err := client.GetProperty(context.Background(), papi.GetPropertyRequest{
		ContractID: "ctr_test1",
		GroupID:    "grp_test1",
		PropertyID: "prp_test1",
	})
	require.NoError(t, err)
	assert.Equal(t, "my_property", result.Property.PropertyName)
	// We expect exactly two requests to the server:
	// - the first resulting in code 429
	// - the second after a proper backoff, resulting in status 200
	assert.Equal(t, []int{http.StatusTooManyRequests, http.StatusOK}, xrlHandler.ReturnedCodes())
	assert.Less(t,
		xrlHandler.ReturnTimes()[1],
		xrlHandler.AvailableAt().Add(time.Duration(time.Millisecond)*1100))
}

func TestXRateLimitPost(t *testing.T) {
	xrlHandler := test.RateLimitHTTPHandler{
		T:           t,
		SuccessCode: http.StatusCreated,
		SuccessBody: `
		{
			"activationLink": "/papi/v1/properties/prp_12345/activations/dummy_activation?contractId=ctr_test1&groupId=grp_test1"
		}`,
	}

	mockServer := httptest.NewTLSServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		assert.Equal(t, "/papi/v1/properties/prp_12345/activations?contractId=ctr_test1&groupId=grp_test1", r.URL.String())
		assert.Equal(t, http.MethodPost, r.Method)
		xrlHandler.ServeHTTP(w, r, "X-RateLimit-Next")
	}))
	defer mockServer.Close()

	client := papi.Client(mockSession(t, mockServer))
	result, err := client.CreateActivation(context.Background(), papi.CreateActivationRequest{
		PropertyID: "prp_12345",
		ContractID: "ctr_test1",
		GroupID:    "grp_test1",
		Activation: papi.Activation{
			PropertyVersion: 1,
			Network:         papi.ActivationNetworkStaging,
			UseFastFallback: false,
			NotifyEmails: []string{
				"you@example.com",
				"them@example.com",
			},
			AcknowledgeWarnings: []string{"foobarbaz"},
		},
	})
	require.NoError(t, err)
	assert.Equal(t, "dummy_activation", result.ActivationID)
	// We expect exactly two requests to the server:
	// - the first resulting in code 429
	// - the second after a proper backoff, resulting in status 201
	assert.Equal(t, []int{http.StatusTooManyRequests, http.StatusCreated}, xrlHandler.ReturnedCodes())
	assert.Less(t,
		xrlHandler.ReturnTimes()[1],
		xrlHandler.AvailableAt().Add(time.Duration(time.Millisecond)*1100))
}

func TestAkamaiRateLimitGet(t *testing.T) {
	arlHandler := test.RateLimitHTTPHandler{
		T:           t,
		SuccessCode: http.StatusOK,
		SuccessBody: `
			{
				"accountId": "A-CCT7890",
				"certificateId": "123",
				"certificateName": "test-cert",
				"certificateStatus": "CSR_READY",
				"certificateType": "THIRD_PARTY",
				"contractId": "C-0N7RAC7",
				"createdBy": "jsmith",
				"createdDate": "2025-09-01T06:16:05.952613Z",
				"csrExpirationDate": "2026-11-03T06:16:07Z",
				"csrPem": "-----BEGIN CERTIFICATE REQUEST-----\nexample-PEM\n-----END CERTIFICATE REQUEST-----\n",
				"keySize": "2048",
				"keyType": "RSA",
				"modifiedBy": "jsmith",
				"modifiedDate": "2025-09-02T06:16:05.952613Z",
				"sans": [
					"example.com",
					"www.example.com"
				],
				"secureNetwork": "ENHANCED_TLS",
				"signedCertificateIssuer": null,
				"signedCertificateNotValidAfterDate": null,
				"signedCertificateNotValidBeforeDate": null,
				"signedCertificatePem": null,
				"signedCertificateSHA256Fingerprint": null,
				"signedCertificateSerialNumber": null,
				"subject": {
					"commonName": "example.com",
					"country": "US",
					"locality": "Cambridge",
					"organization": "ExampleOrg",
					"state": "Massachusetts"
				},
				"trustChainPem": null
			}`,
	}

	mockServer := httptest.NewTLSServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		assert.Equal(t, "/ccm/v1/certificates/123", r.URL.String())
		assert.Equal(t, http.MethodGet, r.Method)
		arlHandler.ServeHTTP(w, r, "Akamai-RateLimit-Next")
	}))
	defer mockServer.Close()

	client := ccm.Client(mockSession(t, mockServer))
	result, err := client.GetCertificate(context.Background(), ccm.GetCertificateRequest{
		CertificateID: "123",
	})
	require.NoError(t, err)
	assert.Equal(t, "test-cert", result.Certificate.CertificateName)
	// We expect exactly two requests to the server:
	// - the first resulting in code 429
	// - the second after a proper backoff, resulting in status 200
	assert.Equal(t, []int{http.StatusTooManyRequests, http.StatusOK}, arlHandler.ReturnedCodes())
	assert.Less(t,
		arlHandler.ReturnTimes()[1],
		arlHandler.AvailableAt().Add(time.Duration(time.Millisecond)*1100))
}

func TestAkamaiRateLimitPost(t *testing.T) {
	arlHandler := test.RateLimitHTTPHandler{
		T:           t,
		SuccessCode: http.StatusCreated,
		SuccessBody: `
			{
				"accountId": "A-CCT7890",
				"certificateId": "123",
				"certificateName": "test-cert",
				"certificateStatus": "CSR_READY",
				"certificateType": "THIRD_PARTY",
				"contractId": "C-0N7RAC7",
				"createdBy": "jsmith",
				"createdDate": "2025-09-01T06:16:05.952613Z",
				"csrExpirationDate": "2026-11-03T06:16:07Z",
				"csrPem": "-----BEGIN CERTIFICATE REQUEST-----\nexample-PEM\n-----END CERTIFICATE REQUEST-----\n",
				"keySize": "2048",
				"keyType": "RSA",
				"modifiedBy": "jsmith",
				"modifiedDate": "2025-09-02T06:16:05.952613Z",
				"sans": [
					"example.com",
					"www.example.com"
				],
				"secureNetwork": "ENHANCED_TLS",
				"signedCertificateIssuer": null,
				"signedCertificateNotValidAfterDate": null,
				"signedCertificateNotValidBeforeDate": null,
				"signedCertificatePem": null,
				"signedCertificateSHA256Fingerprint": null,
				"signedCertificateSerialNumber": null,
				"trustChainPem": null
			}`,
	}

	mockServer := httptest.NewTLSServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		assert.Equal(t, "/ccm/v1/certificates?contractId=111&groupId=222", r.URL.String())
		assert.Equal(t, http.MethodPost, r.Method)
		arlHandler.ServeHTTP(w, r, "Akamai-RateLimit-Next")
	}))
	defer mockServer.Close()

	client := ccm.Client(mockSession(t, mockServer))
	result, err := client.CreateCertificate(context.Background(), ccm.CreateCertificateRequest{
		ContractID: "111",
		GroupID:    "222",
		Body: ccm.CreateCertificateRequestBody{
			CertificateName: "test-cert",
			SANs:            []string{"example.com", "www.example.com"},
			SecureNetwork:   "ENHANCED_TLS",
			KeyType:         "RSA",
			KeySize:         "2048",
		},
	})
	require.NoError(t, err)
	assert.Equal(t, "123", result.Certificate.CertificateID)
	// We expect exactly two requests to the server:
	// - the first resulting in code 429
	// - the second after a proper backoff, resulting in status 201
	assert.Equal(t, []int{http.StatusTooManyRequests, http.StatusCreated}, arlHandler.ReturnedCodes())
	assert.Less(t,
		arlHandler.ReturnTimes()[1],
		arlHandler.AvailableAt().Add(time.Duration(time.Millisecond)*1100))
}
