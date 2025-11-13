package property

import (
	"context"
	"testing"

	"github.com/akamai/AkamaiOPEN-edgegrid-golang/v12/pkg/domainownership"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

func TestValidationHandler_CalculateAndBuildValidateRequest(t *testing.T) {
	t.Parallel()
	tests := map[string]struct {
		apiDomains   map[domainKey]domainDetails
		planDomains  map[domainKey]domainDetails
		stateDomains map[domainKey]domainDetails
		expectedReq  []domainownership.ValidateDomainsRequest
		expectedErr  string
	}{
		"one domain to validate": {
			apiDomains: map[domainKey]domainDetails{
				{domainName: "test.com", validationScope: "HOST"}: {validationStatus: "PENDING"},
			},
			planDomains: map[domainKey]domainDetails{
				{domainName: "test.com", validationScope: "HOST"}: {},
			},
			expectedReq: []domainownership.ValidateDomainsRequest{
				{
					Domains: []domainownership.ValidateDomain{
						{
							DomainName:      "test.com",
							ValidationScope: "HOST",
						},
					},
				},
			},
		},
		"domain already validated": {
			apiDomains: map[domainKey]domainDetails{
				{domainName: "test.com", validationScope: "HOST"}: {validationStatus: "VALIDATED"},
			},
			planDomains: map[domainKey]domainDetails{
				{domainName: "test.com", validationScope: "HOST"}: {},
			},
			expectedReq: nil,
		},
		"domain with invalidated status": {
			apiDomains: map[domainKey]domainDetails{
				{domainName: "test.com", validationScope: "HOST"}: {validationStatus: "INVALIDATED"},
			},
			planDomains: map[domainKey]domainDetails{
				{domainName: "test.com", validationScope: "HOST"}: {},
			},
			expectedErr: "domain test.com with scope HOST is in INVALIDATED status, cannot validate",
		},
		"domain with token expired status": {
			apiDomains: map[domainKey]domainDetails{
				{domainName: "test.com", validationScope: "HOST"}: {validationStatus: "TOKEN_EXPIRED"},
			},
			planDomains: map[domainKey]domainDetails{
				{domainName: "test.com", validationScope: "HOST"}: {},
			},
			expectedErr: "domain test.com with scope HOST is in TOKEN_EXPIRED status, cannot validate",
		},
		"domain not found in api": {
			apiDomains: map[domainKey]domainDetails{},
			planDomains: map[domainKey]domainDetails{
				{domainName: "test.com", validationScope: "HOST"}: {},
			},
			expectedErr: "domain test.com with scope HOST is not found in API",
		},
		"multiple domains, one already validated": {
			apiDomains: map[domainKey]domainDetails{
				{domainName: "test.com", validationScope: "HOST"}:      {validationStatus: "PENDING"},
				{domainName: "example.com", validationScope: "HOST"}:   {validationStatus: "VALIDATED"},
				{domainName: "another.com", validationScope: "HOST"}:   {validationStatus: "PENDING"},
				{domainName: "another.com", validationScope: "DOMAIN"}: {validationStatus: "PENDING"},
			},
			planDomains: map[domainKey]domainDetails{
				{domainName: "test.com", validationScope: "HOST"}:      {},
				{domainName: "example.com", validationScope: "HOST"}:   {},
				{domainName: "another.com", validationScope: "HOST"}:   {},
				{domainName: "another.com", validationScope: "DOMAIN"}: {},
			},
			expectedReq: []domainownership.ValidateDomainsRequest{
				{
					Domains: []domainownership.ValidateDomain{
						{
							DomainName:      "another.com",
							ValidationScope: "DOMAIN",
						},
						{
							DomainName:      "another.com",
							ValidationScope: "HOST",
						},
						{
							DomainName:      "test.com",
							ValidationScope: "HOST",
						},
					},
				},
			},
		},
		"no domains to validate": {
			apiDomains:  map[domainKey]domainDetails{},
			planDomains: map[domainKey]domainDetails{},
			expectedReq: nil,
		},
	}

	for name, tc := range tests {
		t.Run(name, func(t *testing.T) {
			t.Parallel()
			handler, err := newValidationHandler(context.Background()).
				setAPIDomains(tc.apiDomains).
				setPlanDomains(tc.planDomains).
				setStateDomains(tc.stateDomains).
				calculateDomainsToValidate()

			if tc.expectedErr != "" {
				require.Error(t, err)
				assert.Contains(t, err.Error(), tc.expectedErr)
				return
			}

			require.NoError(t, err)
			reqs := handler.buildValidateRequests()
			assert.Equal(t, tc.expectedReq, reqs)
		})
	}
}

func TestValidationHandler_CalculateAndBuildInvalidateRequest(t *testing.T) {
	t.Parallel()
	tests := map[string]struct {
		apiDomains   map[domainKey]domainDetails
		planDomains  map[domainKey]domainDetails
		stateDomains map[domainKey]domainDetails
		expectedReq  *domainownership.InvalidateDomainsRequest
	}{
		"one domain to invalidate": {
			apiDomains: map[domainKey]domainDetails{
				{domainName: "test.com", validationScope: "HOST"}: {validationStatus: "VALIDATED"},
			},
			stateDomains: map[domainKey]domainDetails{
				{domainName: "test.com", validationScope: "HOST"}: {},
			},
			planDomains: map[domainKey]domainDetails{},
			expectedReq: &domainownership.InvalidateDomainsRequest{
				Domains: []domainownership.Domain{
					{DomainName: "test.com", ValidationScope: "HOST"},
				},
			},
		},
		"domain not in validated status": {
			apiDomains: map[domainKey]domainDetails{
				{domainName: "test.com", validationScope: "HOST"}: {validationStatus: "PENDING"},
			},
			stateDomains: map[domainKey]domainDetails{
				{domainName: "test.com", validationScope: "HOST"}: {},
			},
			planDomains: map[domainKey]domainDetails{},
			expectedReq: nil,
		},
		"domain not in api": {
			apiDomains: map[domainKey]domainDetails{},
			stateDomains: map[domainKey]domainDetails{
				{domainName: "test.com", validationScope: "HOST"}: {},
			},
			planDomains: map[domainKey]domainDetails{},
			expectedReq: nil,
		},
		"multiple domains, one not validated": {
			apiDomains: map[domainKey]domainDetails{
				{domainName: "test.com", validationScope: "HOST"}:    {validationStatus: "VALIDATED"},
				{domainName: "example.com", validationScope: "HOST"}: {validationStatus: "PENDING"},
			},
			stateDomains: map[domainKey]domainDetails{
				{domainName: "test.com", validationScope: "HOST"}:    {},
				{domainName: "example.com", validationScope: "HOST"}: {},
			},
			planDomains: map[domainKey]domainDetails{},
			expectedReq: &domainownership.InvalidateDomainsRequest{
				Domains: []domainownership.Domain{
					{DomainName: "test.com", ValidationScope: "HOST"},
				},
			},
		},
		"no domains to invalidate": {
			apiDomains:   map[domainKey]domainDetails{},
			stateDomains: map[domainKey]domainDetails{},
			planDomains:  map[domainKey]domainDetails{},
			expectedReq:  nil,
		},
		"domain still in plan": {
			apiDomains: map[domainKey]domainDetails{
				{domainName: "test.com", validationScope: "HOST"}: {validationStatus: "VALIDATED"},
			},
			stateDomains: map[domainKey]domainDetails{
				{domainName: "test.com", validationScope: "HOST"}: {},
			},
			planDomains: map[domainKey]domainDetails{
				{domainName: "test.com", validationScope: "HOST"}: {},
			},
			expectedReq: nil,
		},
	}

	for name, tc := range tests {
		t.Run(name, func(t *testing.T) {
			t.Parallel()
			handler := newValidationHandler(context.Background()).
				setAPIDomains(tc.apiDomains).
				setPlanDomains(tc.planDomains).
				setStateDomains(tc.stateDomains).
				calculateDomainsToInvalidate()

			req := handler.buildInvalidateRequest()
			assert.Equal(t, tc.expectedReq, req)
		})
	}
}
