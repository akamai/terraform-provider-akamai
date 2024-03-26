package property

import (
	"fmt"

	"github.com/apex/log"

	"github.com/akamai/AkamaiOPEN-edgegrid-golang/v8/pkg/papi"
)

// Get loggable fields from given arguments
func logFields(args ...interface{}) log.Fielder {
	switch len(args) {
	case 0:
		return log.Fields{}
	case 1:
		return toLogFielder(args[0])
	}

	merge := make(mergeLogFields, len(args))
	for i, arg := range args {
		merge[i] = toLogFielder(arg)
	}

	return merge
}

type mergeLogFields []log.Fielder

func (m mergeLogFields) Fields() log.Fields {
	fields := log.Fields{}
	for _, fielder := range m {
		for k, v := range fielder.Fields() {
			fields[k] = v
		}
	}
	return fields
}

// Convert the given value to a type that captures structured logging fields
func toLogFielder(given interface{}) log.Fielder {
	switch v := given.(type) {
	case log.Fielder:
		return v

	case error:
		return log.Fields{"error": v.Error()}

	case papi.CreatePropertyRequest:
		return createPropertyReqFields(v)

	case papi.CreatePropertyResponse:
		return createPropertyResFields(v)

	case papi.RemovePropertyRequest:
		return removePropertyReqFields(v)

	case papi.GetPropertyRequest:
		return getPropertyReqFields(v)

	case papi.GetPropertyResponse:
		return getPropertyResFields(v)

	case papi.CreatePropertyVersionRequest:
		return createPropertyVersionReqFields(v)

	case papi.CreatePropertyVersionResponse:
		return createPropertyVersionResFields(v)

	case papi.UpdatePropertyVersionHostnamesRequest:
		return updatePropertyVersionHostnamesReqFields(v)

	case papi.UpdatePropertyVersionHostnamesResponse:
		return updatePropertyVersionHostnamesResFields(v)

	case papi.GetPropertyVersionHostnamesRequest:
		return getPropertyVersionHostnamesReqFields(v)

	case papi.GetPropertyVersionHostnamesResponse:
		return getPropertyVersionHostnamesResFields(v)

	case papi.GetRuleTreeRequest:
		return getRuleTreeReqFields(v)

	case papi.GetRuleTreeResponse:
		return getRuleTreeResFields(v)

	case papi.UpdateRulesRequest:
		return updateRulesReqFields(v)

	case papi.UpdateRulesResponse:
		return updateRulesResFields(v)

	case papi.GetPropertyVersionRequest:
		return getPropertyVersionReqFields(v)

	case papi.GetPropertyVersionsRequest:
		return getPropertyVersionsReqFields(v)

	case papi.GetPropertyVersionsResponse:
		return getPropertyVersionResFields(v)
	}

	panic(fmt.Sprintf("no known log.Fielder for %T", given))
}

type updateRulesReqFields papi.UpdateRulesRequest

func (r updateRulesReqFields) Fields() log.Fields {
	return log.Fields{
		"property_id":      r.PropertyID,
		"group_id":         r.GroupID,
		"contract_id":      r.ContractID,
		"property_version": r.PropertyVersion,
		"validate_mode":    r.ValidateMode,
		"validate_rules":   r.ValidateRules,
	}
}

type updateRulesResFields papi.UpdateRulesResponse

func (r updateRulesResFields) Fields() log.Fields {
	return log.Fields{
		"account_id":       r.AccountID,
		"property_id":      r.PropertyID,
		"group_id":         r.GroupID,
		"contract_id":      r.ContractID,
		"property_version": r.PropertyVersion,
		"rule_format":      r.RuleFormat,
	}
}

type getRuleTreeReqFields papi.GetRuleTreeRequest

func (r getRuleTreeReqFields) Fields() log.Fields {
	return log.Fields{
		"property_id":      r.PropertyID,
		"group_id":         r.GroupID,
		"contract_id":      r.ContractID,
		"property_version": r.PropertyVersion,
		"validate_mode":    r.ValidateMode,
		"validate_rules":   r.ValidateRules,
	}
}

type getRuleTreeResFields papi.GetRuleTreeResponse

func (r getRuleTreeResFields) Fields() log.Fields {
	return log.Fields{
		"property_id":      r.PropertyID,
		"group_id":         r.GroupID,
		"contract_id":      r.ContractID,
		"property_version": r.PropertyVersion,
		"rules_format":     r.RuleFormat,
	}
}

type createPropertyReqFields papi.CreatePropertyRequest

func (req createPropertyReqFields) Fields() log.Fields {
	return log.Fields{
		"property_name": req.Property.PropertyName,
		"contract_id":   req.ContractID,
		"group_id":      req.GroupID,
		"product_id":    req.Property.ProductID,
	}
}

type getPropertyVersionReqFields papi.GetPropertyVersionRequest

func (req getPropertyVersionReqFields) Fields() log.Fields {
	return log.Fields{
		"property_name":    req.PropertyID,
		"contract_id":      req.ContractID,
		"group_id":         req.GroupID,
		"property_version": req.PropertyVersion,
	}
}

type getPropertyVersionsReqFields papi.GetPropertyVersionsRequest

func (req getPropertyVersionsReqFields) Fields() log.Fields {
	return log.Fields{
		"property_name": req.PropertyID,
		"contract_id":   req.ContractID,
		"group_id":      req.GroupID,
		"limit":         req.Limit,
		"offset":        req.Offset,
	}
}

type getPropertyVersionResFields papi.GetPropertyVersionsResponse

func (res getPropertyVersionResFields) Fields() log.Fields {
	return log.Fields{
		"property_id":      res.PropertyID,
		"contract_id":      res.ContractID,
		"group_id":         res.GroupID,
		"property_version": res.Version.PropertyVersion,
		"product_id":       res.Version.ProductID,
	}
}

type createPropertyResFields papi.CreatePropertyResponse

func (res createPropertyResFields) Fields() log.Fields {
	return log.Fields{
		"property_id": res.PropertyID,
	}
}

type getPropertyReqFields papi.GetPropertyRequest

func (req getPropertyReqFields) Fields() log.Fields {
	return log.Fields{
		"property_id": req.PropertyID,
		"contract_id": req.ContractID,
		"group_id":    req.GroupID,
	}
}

type getPropertyResFields papi.GetPropertyResponse

func (res getPropertyResFields) Fields() log.Fields {
	fields := log.Fields{
		"contract_id": res.ContractID,
		"group_id":    res.GroupID,
	}

	if res.Property != nil {
		fields["property_id"] = res.Property.PropertyID
	}

	return fields
}

type removePropertyReqFields papi.RemovePropertyRequest

func (req removePropertyReqFields) Fields() log.Fields {
	return log.Fields{
		"property_id": req.PropertyID,
		"contract_id": req.ContractID,
		"group_id":    req.GroupID,
	}
}

type createPropertyVersionReqFields papi.CreatePropertyVersionRequest

func (req createPropertyVersionReqFields) Fields() log.Fields {
	return log.Fields{
		"property_id":           req.PropertyID,
		"contract_id":           req.ContractID,
		"group_id":              req.GroupID,
		"from_property_version": req.Version.CreateFromVersion,
	}
}

type createPropertyVersionResFields papi.CreatePropertyVersionResponse

func (res createPropertyVersionResFields) Fields() log.Fields {
	return log.Fields{
		"property_version": res.PropertyVersion,
	}
}

type updatePropertyVersionHostnamesReqFields papi.UpdatePropertyVersionHostnamesRequest

func (req updatePropertyVersionHostnamesReqFields) Fields() log.Fields {
	hostnames := map[string]string{}
	for _, hn := range req.Hostnames {
		hostnames[hn.CnameFrom] = hn.EdgeHostnameID
		hostnames["cnameTo"] = hn.CnameTo
		hostnames["certProvisioningType"] = hn.CertProvisioningType
	}

	return log.Fields{
		"property_id": req.PropertyID,
		"contract_id": req.ContractID,
		"group_id":    req.GroupID,
		"hostnames":   hostnames,
	}
}

type updatePropertyVersionHostnamesResFields papi.UpdatePropertyVersionHostnamesResponse

func (res updatePropertyVersionHostnamesResFields) Fields() log.Fields {
	hostnames := map[string]string{}
	for _, hn := range res.Hostnames.Items {
		hostnames[hn.CnameFrom] = hn.EdgeHostnameID
		hostnames["cnameTo"] = hn.CnameTo
		hostnames["certProvisioningType"] = hn.CertProvisioningType
		certs := map[string]interface{}{}
		certs["validation_cname.hostname"] = hn.CertStatus.ValidationCname.Hostname
		certs["validation_cname.target"] = hn.CertStatus.ValidationCname.Hostname
		if len(hn.CertStatus.Staging) > 0 {
			certs["staging_status"] = hn.CertStatus.Staging[0].Status
		}
		if len(hn.CertStatus.Production) > 0 {
			certs["production_status"] = hn.CertStatus.Production[0].Status
		}
		//hostnames["cert_status"] = certs
	}

	return log.Fields{
		"property_id": res.PropertyID,
		"contract_id": res.ContractID,
		"group_id":    res.GroupID,
		"hostnames":   hostnames,
	}
}

type getPropertyVersionHostnamesReqFields papi.GetPropertyVersionHostnamesRequest

func (req getPropertyVersionHostnamesReqFields) Fields() log.Fields {
	return log.Fields{
		"property_id":      req.PropertyID,
		"contract_id":      req.ContractID,
		"group_id":         req.GroupID,
		"property_version": req.PropertyVersion,
	}
}

type getPropertyVersionHostnamesResFields papi.GetPropertyVersionHostnamesResponse

func (res getPropertyVersionHostnamesResFields) Fields() log.Fields {
	hostnames := map[string]string{}
	for _, hn := range res.Hostnames.Items {
		hostnames[hn.CnameFrom] = hn.EdgeHostnameID
	}

	return log.Fields{
		"property_id":      res.PropertyID,
		"contract_id":      res.ContractID,
		"group_id":         res.GroupID,
		"property_version": res.PropertyVersion,
		"hostnames":        hostnames,
	}
}
