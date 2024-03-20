package property

import (
	"errors"
	"fmt"
	"strings"

	"github.com/akamai/AkamaiOPEN-edgegrid-golang/v8/pkg/papi"
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

var complianceRecordSchema = &schema.Resource{
	Schema: map[string]*schema.Schema{
		"noncompliance_reason_none": {
			Type:     schema.TypeList,
			Optional: true,
			MaxItems: 1,
			ExactlyOneOf: []string{
				"compliance_record.0.noncompliance_reason_other",
				"compliance_record.0.noncompliance_reason_no_production_traffic",
				"compliance_record.0.noncompliance_reason_emergency",
			},
			Description: fmt.Sprintf("Provides an audit record when activating on a production network with noncompliance reason as `%s`", papi.NoncomplianceReasonNone),
			Elem: &schema.Resource{
				Schema: map[string]*schema.Schema{
					"ticket_id": {
						Type:        schema.TypeString,
						Optional:    true,
						Description: "Identifies the ticket that describes the need for the activation",
					},
					"customer_email": {
						Type:        schema.TypeString,
						Optional:    true,
						Description: "Identifies the customer",
					},
					"peer_reviewed_by": {
						Type:        schema.TypeString,
						Optional:    true,
						Description: "Identifies person who has independently approved the activation request",
					},
					"unit_tested": {
						Type:        schema.TypeBool,
						Optional:    true,
						Description: "Whether the metadata to activate has been fully tested",
					},
				},
			},
		},
		"noncompliance_reason_other": {
			Type:        schema.TypeList,
			Optional:    true,
			MaxItems:    1,
			Description: fmt.Sprintf("Provides an audit record when activating on a production network with noncompliance reason as `%s`", papi.NoncomplianceReasonOther),
			Elem: &schema.Resource{
				Schema: map[string]*schema.Schema{
					"ticket_id": {
						Type:        schema.TypeString,
						Optional:    true,
						Description: "Identifies the ticket that describes the need for the activation",
					},
					"other_noncompliance_reason": {
						Type:        schema.TypeString,
						Optional:    true,
						Description: "Describes the reason why the activation must occur immediately, out of compliance with the standard procedure",
					},
				},
			},
		},
		"noncompliance_reason_no_production_traffic": {
			Type:        schema.TypeList,
			Optional:    true,
			MaxItems:    1,
			Description: fmt.Sprintf("Provides an audit record when activating on a production network with noncompliance reason as `%s`", papi.NoncomplianceReasonNoProductionTraffic),
			Elem: &schema.Resource{
				Schema: map[string]*schema.Schema{
					"ticket_id": {
						Type:        schema.TypeString,
						Optional:    true,
						Description: "Identifies the ticket that describes the need for the activation",
					},
				},
			},
		},
		"noncompliance_reason_emergency": {
			Type:        schema.TypeList,
			Optional:    true,
			MaxItems:    1,
			Description: fmt.Sprintf("Provides an audit record when activating on a production network with noncompliance reason as `%s`", papi.NoncomplianceReasonEmergency),
			Elem: &schema.Resource{
				Schema: map[string]*schema.Schema{
					"ticket_id": {
						Type:        schema.TypeString,
						Optional:    true,
						Description: "Identifies the ticket that describes the need for the activation",
					},
				},
			},
		},
	},
}

// Convert given hostnames to the map form that can be stored in a schema.ResourceData
// Setting only statuses for default certs if they exist
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

func papiErrorsToList(errors []*papi.Error) []map[string]interface{} {
	if len(errors) == 0 {
		return nil
	}

	var ruleErrors []map[string]interface{}
	for _, err := range errors {
		if err == nil {
			continue
		}

		ruleErrors = append(ruleErrors, papiErrorToMap(err))
	}

	return ruleErrors
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
