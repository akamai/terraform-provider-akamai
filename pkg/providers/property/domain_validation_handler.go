package property

import (
	"context"
	"fmt"
	"slices"
	"sort"

	"github.com/akamai/AkamaiOPEN-edgegrid-golang/v12/pkg/domainownership"
	"github.com/akamai/terraform-provider-akamai/v9/pkg/common/ptr"
	"github.com/hashicorp/terraform-plugin-log/tflog"
)

type validationHandler struct {
	ctx                 context.Context
	apiDomains          map[domainKey]domainDetails
	stateDomains        map[domainKey]domainDetails
	planDomains         map[domainKey]domainDetails
	domainsToValidate   map[domainKey]domainDetails
	domainsToInvalidate map[domainKey]domainDetails
}

func newValidationHandler(ctx context.Context) *validationHandler {
	return &validationHandler{
		ctx:          ctx,
		apiDomains:   make(map[domainKey]domainDetails),
		stateDomains: make(map[domainKey]domainDetails),
		planDomains:  make(map[domainKey]domainDetails),
	}
}

func (v *validationHandler) setAPIDomains(domains map[domainKey]domainDetails) *validationHandler {
	v.apiDomains = domains
	tflog.Debug(v.ctx, "set API domains map", map[string]any{
		"api_domains": v.apiDomains,
	})
	return v
}

func (v *validationHandler) setStateDomains(domains map[domainKey]domainDetails) *validationHandler {
	v.stateDomains = domains
	tflog.Debug(v.ctx, "set state domains map", map[string]any{
		"state_domains": v.stateDomains,
	})
	return v
}

func (v *validationHandler) setPlanDomains(domains map[domainKey]domainDetails) *validationHandler {
	v.planDomains = domains
	tflog.Debug(v.ctx, "set plan domains map", map[string]any{
		"plan_domains": v.planDomains,
	})
	return v
}

func (v *validationHandler) calculateDomainsToValidate() (*validationHandler, error) {
	domainsToValidate := make(map[domainKey]domainDetails)

	for planDomainKey, planDomainDetails := range v.planDomains {
		// If the domain is already VALIDATED, no need to validate again.
		apiDomain, foundInAPI := v.apiDomains[planDomainKey]
		if foundInAPI {
			if apiDomain.validationStatus == "VALIDATED" {
				tflog.Debug(v.ctx, "domain already validated in API, skipping", map[string]any{
					"domain_name":      planDomainKey.domainName,
					"validation_scope": planDomainKey.validationScope,
				})
				continue
			}
			if apiDomain.validationStatus == "INVALIDATED" || apiDomain.validationStatus == "TOKEN_EXPIRED" {
				return nil, fmt.Errorf("domain %s with scope %s is in %s status, cannot validate. Please recreate the domain prior to validating it", planDomainKey.domainName, planDomainKey.validationScope, apiDomain.validationStatus)
			}
		} else {
			return nil, fmt.Errorf("domain %s with scope %s is not found in API", planDomainKey.domainName, planDomainKey.validationScope)
		}

		domainsToValidate[planDomainKey] = planDomainDetails
	}
	v.domainsToValidate = domainsToValidate
	tflog.Debug(v.ctx, "calculated domains to validate", map[string]any{
		"domains_to_validate": v.domainsToValidate,
	})

	return v, nil
}

func (v *validationHandler) buildValidateRequests() []domainownership.ValidateDomainsRequest {
	var domainsToValidateSlice []domainownership.ValidateDomain
	for domainKey, domainDetails := range v.domainsToValidate {
		r := domainownership.ValidateDomain{
			DomainName:      domainKey.domainName,
			ValidationScope: domainownership.ValidationScope(domainKey.validationScope),
		}
		if domainDetails.validationMethod != nil {
			r.ValidationMethod = ptr.To(domainownership.ValidationMethod(*domainDetails.validationMethod))
		}
		domainsToValidateSlice = append(domainsToValidateSlice, r)
	}

	// Sort the slice to have consistent order.
	sort.Slice(domainsToValidateSlice, func(i, j int) bool {
		if domainsToValidateSlice[i].DomainName == domainsToValidateSlice[j].DomainName {
			return domainsToValidateSlice[i].ValidationScope < domainsToValidateSlice[j].ValidationScope
		}
		return domainsToValidateSlice[i].DomainName < domainsToValidateSlice[j].DomainName
	})

	var requests []domainownership.ValidateDomainsRequest
	// Single request can handle up to 100 domains.
	for chunk := range slices.Chunk(domainsToValidateSlice, 100) {
		requests = append(requests, domainownership.ValidateDomainsRequest{Domains: chunk})
	}
	tflog.Debug(v.ctx, "built validate requests", map[string]any{
		"validate_requests": requests,
	})

	return requests
}

func (v *validationHandler) calculateDomainsToInvalidate() *validationHandler {
	domainsToInvalidate := make(map[domainKey]domainDetails)

	// Compare state and plan to find removed domains.
	for stateDomainKey, stateDomainDetails := range v.stateDomains {
		if _, foundInPlan := v.planDomains[stateDomainKey]; !foundInPlan {
			domainsToInvalidate[stateDomainKey] = stateDomainDetails
		}
	}

	tflog.Debug(v.ctx, "calculated domains to invalidate based on state and plan", map[string]any{
		"domains_to_invalidate_before_status_check": domainsToInvalidate,
	})

	for domainKey := range domainsToInvalidate {
		apiDomain, foundInAPI := v.apiDomains[domainKey]
		if foundInAPI {
			// Do not invalidate if the status is different from VALIDATED.
			if apiDomain.validationStatus != "VALIDATED" {
				tflog.Debug(v.ctx, "domain not in VALIDATED status in API, skipping invalidation", map[string]any{
					"domain_name":      domainKey.domainName,
					"validation_scope": domainKey.validationScope,
					"status":           apiDomain.validationStatus,
				})
				delete(domainsToInvalidate, domainKey)
			}
		} else {
			// Do not invalidate if the domain is not found in API.
			tflog.Debug(v.ctx, "domain not found in API, skipping invalidation", map[string]any{
				"domain_name":      domainKey.domainName,
				"validation_scope": domainKey.validationScope,
			})
			delete(domainsToInvalidate, domainKey)
		}
	}
	v.domainsToInvalidate = domainsToInvalidate
	tflog.Debug(v.ctx, "final domains to invalidate", map[string]any{
		"domains_to_invalidate": v.domainsToInvalidate,
	})

	return v
}

func (v *validationHandler) buildInvalidateRequest() *domainownership.InvalidateDomainsRequest {
	var domainsToInvalidateSlice []domainownership.Domain

	if len(v.domainsToInvalidate) == 0 {
		return nil
	}

	for domain := range v.domainsToInvalidate {
		domainsToInvalidateSlice = append(domainsToInvalidateSlice, domainownership.Domain{
			DomainName:      domain.domainName,
			ValidationScope: domainownership.ValidationScope(domain.validationScope),
		})
	}

	// Sort the slice to have consistent order.
	sort.Slice(domainsToInvalidateSlice, func(i, j int) bool {
		return domainsToInvalidateSlice[i].DomainName < domainsToInvalidateSlice[j].DomainName
	})

	tflog.Debug(v.ctx, "built invalidate request", map[string]any{
		"invalidate_request": domainsToInvalidateSlice,
	})

	return &domainownership.InvalidateDomainsRequest{Domains: domainsToInvalidateSlice}
}
