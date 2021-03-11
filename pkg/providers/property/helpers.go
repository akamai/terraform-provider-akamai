package property

import (
	"github.com/akamai/AkamaiOPEN-edgegrid-golang/v2/pkg/papi"
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
	var c []map[string]interface{}
	for _, hn := range Hostnames {
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
