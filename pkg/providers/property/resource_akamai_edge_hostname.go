package property

import (
	"errors"
	"fmt"
	"strings"

	edge "github.com/akamai/AkamaiOPEN-edgegrid-golang/edgegrid"
	"github.com/akamai/AkamaiOPEN-edgegrid-golang/jsonhooks-v1"
	"github.com/akamai/AkamaiOPEN-edgegrid-golang/papi-v1"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/schema"
)

func resourceSecureEdgeHostName() *schema.Resource {
	return &schema.Resource{
		Create: resourceSecureEdgeHostNameCreate,
		Read:   resourceSecureEdgeHostNameRead,
		Delete: resourceSecureEdgeHostNameDelete,
		Exists: resourceSecureEdgeHostNameExists,
		Importer: &schema.ResourceImporter{
			State: resourceSecureEdgeHostNameImport,
		},
		Schema: akamaiSecureEdgeHostNameSchema,
	}
}

var akamaiSecureEdgeHostNameSchema = map[string]*schema.Schema{
	"product": {
		Type:     schema.TypeString,
		Required: true,
		ForceNew: true,
	},
	"contract": {
		Type:     schema.TypeString,
		Required: true,
		ForceNew: true,
	},
	"group": {
		Type:     schema.TypeString,
		Required: true,
		ForceNew: true,
	},
	"edge_hostname": {
		Type:     schema.TypeString,
		Required: true,
		ForceNew: true,
	},
	"ipv4": {
		Type:     schema.TypeBool,
		Optional: true,
		Default:  true,
		ForceNew: true,
	},
	"ipv6": {
		Type:     schema.TypeBool,
		Optional: true,
		Default:  false,
		ForceNew: true,
	},
	"ip_behavior": {
		Type:     schema.TypeString,
		Computed: true,
	},
	"certificate": {
		Type:     schema.TypeInt,
		Optional: true,
		ForceNew: true,
	},
}

func resourceSecureEdgeHostNameCreate(d *schema.ResourceData, meta interface{}) error {
	d.Partial(true)
	CorrelationID := "[PAPI][resourceSecureEdgeHostNameCreate-" + CreateNonce() + "]"
	group, e := getGroup(d, CorrelationID)
	if e != nil {
		return e
	}
	//	log.Println("[DEBUG] Edgehostnames GROUP = ", group)
	edge.PrintfCorrelation("[DEBUG]", CorrelationID, fmt.Sprintf("  Edgehostnames GROUP = %v", group))
	contract, e := getContract(d, CorrelationID)
	if e != nil {
		return e
	}
	//log.Println("[DEBUG] Edgehostnames CONTRACT = ", contract)
	edge.PrintfCorrelation("[DEBUG]", CorrelationID, fmt.Sprintf("  Edgehostnames CONTRACT = %v", contract))
	product, e := getProduct(d, contract, CorrelationID)
	if e != nil {
		return e
	}

	if group == nil {
		return errors.New("group must be specified to create a new Edge Hostname")
	}

	if contract == nil {
		return errors.New("contract must be specified to create a new Edge Hostname")
	}

	if product == nil {
		return errors.New("product must be specified to create a new Edge Hostname")
	}

	edgeHostnames, err := papi.GetEdgeHostnames(contract, group, "")
	if err != nil {
		return err
	}

	edgeHostname := d.Get("edge_hostname").(string)

	ehn := edgeHostnames.NewEdgeHostname()
	ehn.ProductID = product.ProductID
	ehn.EdgeHostnameDomain = edgeHostname

	switch {
	case strings.HasSuffix(edgeHostname, ".edgesuite.net"):
		ehn.DomainPrefix = strings.TrimSuffix(edgeHostname, ".edgesuite.net")
		ehn.DomainSuffix = "edgesuite.net"
		ehn.SecureNetwork = "STANDARD_TLS"
	case strings.HasSuffix(edgeHostname, ".edgekey.net"):
		ehn.DomainPrefix = strings.TrimSuffix(edgeHostname, ".edgekey.net")
		ehn.DomainSuffix = "edgekey.net"
		ehn.SecureNetwork = "ENHANCED_TLS"
	case strings.HasSuffix(edgeHostname, ".akamaized.net"):
		ehn.DomainPrefix = strings.TrimSuffix(edgeHostname, ".akamaized.net")
		ehn.DomainSuffix = "akamaized.net"
		ehn.SecureNetwork = "SHARED_CERT"
	}

	ipv4 := d.Get("ipv4").(bool)
	if ipv4 {
		ehn.IPVersionBehavior = "IPV4"
	}

	ipv6 := d.Get("ipv6").(bool)
	if ipv6 {
		ehn.IPVersionBehavior = "IPV6"
	}

	if ipv4 && ipv6 {
		ehn.IPVersionBehavior = "IPV6_COMPLIANCE"
	}

	d.Set("ip_behavior", ehn.IPVersionBehavior)

	if certEnrollmentId, ok := d.GetOk("certificate"); ok {
		ehn.CertEnrollmentId = certEnrollmentId.(int)
		ehn.SlotNumber = certEnrollmentId.(int)
	} else if ehn.SecureNetwork == "ENHANCED_TLS" {
		return errors.New("A certificate enrollment ID is required for Enhanced TLS (edgekey.net) edge hostnames")
	}

	if ehnFound, err := edgeHostnames.FindEdgeHostname(ehn); ehnFound != nil && ehnFound.EdgeHostnameID != "" {

		jsonBody, e := jsonhooks.Marshal(ehnFound)
		if e == nil {
			edge.PrintfCorrelation("[DEBUG]", CorrelationID, fmt.Sprintf("EHN Found = %s\n", jsonBody))
		}

		if ehnFound.IPVersionBehavior != ehn.IPVersionBehavior {
			return fmt.Errorf("existing edge hostname found with incompatible IP version (%s vs %s). You must use the same settings, or try a different edge hostname", ehnFound.IPVersionBehavior, ehn.IPVersionBehavior)
		}

		edge.PrintfCorrelation("[DEBUG]", CorrelationID, fmt.Sprintf("Existing edge hostname FOUND = %s", ehnFound.EdgeHostnameID))
		d.SetId(ehnFound.EdgeHostnameID)
	} else {
		edge.PrintfCorrelation("[DEBUG]", CorrelationID, fmt.Sprintf("Creating new edge hostname: %#v\n\n", ehn))
		err = ehn.Save("", CorrelationID)
		if err != nil {
			return err
		}
		d.SetId(ehn.EdgeHostnameID)
	}

	d.Partial(false)

	edge.PrintfCorrelation("[DEBUG]", CorrelationID, "Done")
	return nil
}

func resourceSecureEdgeHostNameDelete(d *schema.ResourceData, meta interface{}) error {
	CorrelationID := "[PAPI][resourceSecureEdgeHostNameDelete-" + CreateNonce() + "]"
	edge.PrintfCorrelation("[DEBUG]", CorrelationID, "DELETING")
	d.SetId("")

	edge.PrintfCorrelation("[DEBUG]", CorrelationID, "Done")
	return nil
}

func resourceSecureEdgeHostNameImport(d *schema.ResourceData, meta interface{}) ([]*schema.ResourceData, error) {
	resourceID := d.Id()
	propertyID := resourceID

	if !strings.HasPrefix(resourceID, "prp_") {
		for _, searchKey := range []papi.SearchKey{papi.SearchByPropertyName, papi.SearchByHostname, papi.SearchByEdgeHostname} {
			results, err := papi.Search(searchKey, resourceID, "") //<--correlationid
			if err != nil {
				continue
			}

			if results != nil && len(results.Versions.Items) > 0 {
				propertyID = results.Versions.Items[0].PropertyID
				break
			}
		}
	}

	property := papi.NewProperty(papi.NewProperties())
	property.PropertyID = propertyID
	e := property.GetProperty("")
	if e != nil {
		return nil, e
	}

	d.Set("account", property.AccountID)
	d.Set("contract", property.ContractID)
	d.Set("group", property.GroupID)

	d.Set("name", property.PropertyName)
	d.Set("version", property.LatestVersion)
	d.SetId(property.PropertyID)

	return []*schema.ResourceData{d}, nil
}

func resourceSecureEdgeHostNameExists(d *schema.ResourceData, meta interface{}) (bool, error) {

	CorrelationID := "[PAPI][resourceSecureEdgeHostNameCreate-" + CreateNonce() + "]"
	group, e := getGroup(d, CorrelationID)
	if e != nil {
		return false, e
	}
	edge.PrintfCorrelation("[DEBUG]", CorrelationID, fmt.Sprintf("Figuring out edgehostnames GROUP = %v", group))
	contract, e := getContract(d, CorrelationID)
	if e != nil {
		return false, e
	}
	edge.PrintfCorrelation("[DEBUG]", CorrelationID, fmt.Sprintf("Figuring out edgehostnames CONTRACT = %v", contract))
	property := papi.NewProperty(papi.NewProperties())
	property.Group = group
	property.Contract = contract

	edge.PrintfCorrelation("[DEBUG]", CorrelationID, fmt.Sprintf("Figuring out edgehostnames %v", d.Id()))
	edgeHostnames := papi.NewEdgeHostnames()
	edge.PrintfCorrelation("[DEBUG]", CorrelationID, fmt.Sprintf("NewEdgeHostnames empty struct  %s", edgeHostnames.ContractID))
	err := edgeHostnames.GetEdgeHostnames(property.Contract, property.Group, d.Id(), CorrelationID)
	if err != nil {
		return false, err
	}
	edge.PrintfCorrelation("[DEBUG]", CorrelationID, "Edgehostname EXISTS in contract ")

	return true, nil
}

func resourceSecureEdgeHostNameRead(d *schema.ResourceData, meta interface{}) error {
	CorrelationID := "[PAPI][resourceSecureEdgeHostNameCreate-" + CreateNonce() + "]"
	d.Partial(true)

	group, e := getGroup(d, CorrelationID)
	if e != nil {
		return e
	}
	edge.PrintfCorrelation("[DEBUG]", CorrelationID, fmt.Sprintf("Figuring out edgehostnames GROUP = %v", group))
	contract, e := getContract(d, CorrelationID)
	if e != nil {
		return e
	}

	edge.PrintfCorrelation("[DEBUG]", CorrelationID, fmt.Sprintf("Figuring out edgehostnames CONTRACT = %v", contract))
	property := papi.NewProperty(papi.NewProperties())
	property.Group = group
	property.Contract = contract

	edge.PrintfCorrelation("[DEBUG]", CorrelationID, fmt.Sprintf("Figuring out edgehostnames %v", d.Id()))
	edgeHostnames := papi.NewEdgeHostnames()
	edge.PrintfCorrelation("[DEBUG]", CorrelationID, fmt.Sprintf("NewEdgeHostnames empty struct %v ", edgeHostnames.ContractID))
	err := edgeHostnames.GetEdgeHostnames(property.Contract, property.Group, "", CorrelationID)
	if err != nil {
		return err
	}
	edge.PrintfCorrelation("[DEBUG]", CorrelationID, "EdgeHostnames exist in contract  ")

	edge.PrintfCorrelation("[DEBUG]", CorrelationID, fmt.Sprintf("Edgehostnames Default host %v", edgeHostnames.EdgeHostnames.Items[0]))
	defaultEdgeHostname := edgeHostnames.EdgeHostnames.Items[0]

	foundEdgeHostname := false
	var edgeHostnameID string
	edgeHostname, edgeHostnameOk := d.GetOk("edge_hostname")

	if edgeHostnameOk {
		for _, eHn := range edgeHostnames.EdgeHostnames.Items {

			if eHn.EdgeHostnameDomain == edgeHostname.(string) {
				foundEdgeHostname = true
				defaultEdgeHostname = eHn
				edgeHostnameID = eHn.EdgeHostnameID
			}
		}
		edge.PrintfCorrelation("[DEBUG]", CorrelationID, fmt.Sprintf("Found EdgeHostname %v", foundEdgeHostname))
		edge.PrintfCorrelation("[DEBUG]", CorrelationID, fmt.Sprintf("Default EdgeHostname %v", defaultEdgeHostname))
	}

	d.Set("contract", contract)
	d.Set("group", group)

	d.SetId(edgeHostnameID)

	return nil
}
