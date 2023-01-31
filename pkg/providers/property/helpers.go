package property

import (
	"errors"
	"strings"

	"github.com/akamai/AkamaiOPEN-edgegrid-golang/v4/pkg/papi"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/schema"
)

var certStatus = &schema.Resource{
	Schema: map[string]*schema.Schema{
		"target": {
			Type:     schema.TypeString,
			Computed: true,
		},
		"hostname": {
			Type:     schema.TypeString,
			Computed: true,
		},
		"production_status": {
			Type:     schema.TypeString,
			Computed: true,
		},
		"staging_status": {
			Type:     schema.TypeString,
			Computed: true,
		},
	},
}

// Convert given hostnames to the map form that can be stored in a schema.ResourceData
// Setting only statuses for default certs if they exist
// TODO Set certstatus object for cps managed certs and default certs once PAPI adds support
func flattenHostnames(Hostnames []papi.Hostname) []map[string]interface{} {
	var res []map[string]interface{}
	for _, hn := range Hostnames {
		var c []map[string]interface{}
		m := map[string]interface{}{}
		m["cname_from"] = hn.CnameFrom
		m["cname_to"] = hn.CnameTo
		m["cert_provisioning_type"] = hn.CertProvisioningType
		m["edge_hostname_id"] = hn.EdgeHostnameID
		m["cname_type"] = hn.CnameType
		certs := map[string]interface{}{}
		certs["hostname"] = hn.CertStatus.ValidationCname.Hostname
		certs["target"] = hn.CertStatus.ValidationCname.Target
		if len(hn.CertStatus.Staging) > 0 {
			certs["staging_status"] = hn.CertStatus.Staging[0].Status
		}
		if len(hn.CertStatus.Production) > 0 {
			certs["production_status"] = hn.CertStatus.Production[0].Status
		}
		c = append(c, certs)
		m["cert_status"] = c
		res = append(res, m)
	}
	return res
}
func papiErrorsToList(Errors []*papi.Error) []interface{} {
	if len(Errors) == 0 {
		return nil
	}

	var RuleErrors []interface{}

	for _, err := range Errors {
		if err == nil {
			continue
		}

		RuleErrors = append(RuleErrors, papiErrorToMap(err))
	}

	return RuleErrors
}

func papiErrorToMap(err *papi.Error) map[string]interface{} {
	if err == nil {
		return nil
	}

	return map[string]interface{}{
		"type":           err.Type,
		"title":          err.Title,
		"detail":         err.Detail,
		"instance":       err.Instance,
		"behavior_name":  err.BehaviorName,
		"error_location": err.ErrorLocation,
		"status_code":    err.StatusCode,
	}
}

// NetworkAlias parses the given network name or alias and returns its full name and any error
func NetworkAlias(network string) (string, error) {

	networks := map[string]papi.ActivationNetwork{
		"STAGING":    papi.ActivationNetworkStaging,
		"STAGE":      papi.ActivationNetworkStaging,
		"STAG":       papi.ActivationNetworkStaging,
		"S":          papi.ActivationNetworkStaging,
		"PRODUCTION": papi.ActivationNetworkProduction,
		"PROD":       papi.ActivationNetworkProduction,
		"P":          papi.ActivationNetworkProduction,
	}

	networkValue, ok := networks[strings.ToUpper(network)]
	if !ok {
		return "", errors.New("network not recognized")
	}
	return string(networkValue), nil
}
