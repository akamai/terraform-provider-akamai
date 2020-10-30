package property

import (
	"fmt"

	"github.com/apex/log"

	"github.com/akamai/AkamaiOPEN-edgegrid-golang/v2/pkg/papi"
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
	switch given.(type) {
	case log.Fielder:
		return given.(log.Fielder)

	case papi.CreatePropertyRequest:
		return createPropertyReqFields(given.(papi.CreatePropertyRequest))

	case papi.CreatePropertyResponse:
		return createPropertyResFields(given.(papi.CreatePropertyResponse))

	case papi.RemovePropertyRequest:
		return removePropertyReqFields(given.(papi.RemovePropertyRequest))

	case papi.GetPropertyRequest:
		return getPropertyReqFields(given.(papi.GetPropertyRequest))

	case papi.GetPropertyResponse:
		return getPropertyResFields(given.(papi.GetPropertyResponse))

	case papi.CreatePropertyVersionRequest:
		return createPropertyVersionReqFields(given.(papi.CreatePropertyVersionRequest))

	case papi.CreatePropertyVersionResponse:
		return createPropertyVersionResFields(given.(papi.CreatePropertyVersionResponse))

	case papi.UpdatePropertyVersionHostnamesRequest:
		return updatePropertyVersionHostnamesReqFields(given.(papi.UpdatePropertyVersionHostnamesRequest))

	case papi.UpdatePropertyVersionHostnamesResponse:
		return updatePropertyVersionHostnamesResFields(given.(papi.UpdatePropertyVersionHostnamesResponse))

	case papi.GetPropertyVersionHostnamesRequest:
		return getPropertyVersionHostnamesReqFields(given.(papi.GetPropertyVersionHostnamesRequest))

	case papi.GetPropertyVersionHostnamesResponse:
		return getPropertyVersionHostnamesResFields(given.(papi.GetPropertyVersionHostnamesResponse))

	case error:
		return log.Fields{"error": given.(error).Error()}
	}

	panic(fmt.Sprintf("no known log.Fielder for %T", given))
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
