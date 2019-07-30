package akamai

import (
	"errors"
	"fmt"
	"github.com/akamai/AkamaiOPEN-edgegrid-golang/papi-v1"
	"github.com/hashicorp/terraform/helper/schema"
	"log"
	"strings"
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

	group, e := getGroup(d)
	if e != nil {
		return e
	}
	log.Println("[DEBUG] Edgehostnames GROUP = ", group)
	contract, e := getContract(d)
	if e != nil {
		return e
	}
	log.Println("[DEBUG] Edgehostnames CONTRACT = ", contract)

	product, e := getProduct(d, contract)
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
		ehn.DomainPrefix = strings.TrimSuffix(edgeHostname, ".akamized.net")
		ehn.DomainSuffix = "akamized.net"
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
	} else if ehn.SecureNetwork == "ENHANCED_TLS" {
		return errors.New("A certificate enrollment ID is required for Enhanced TLS (edgekey.net) edge hostnames")
	}

	if ehnFound, err := edgeHostnames.FindEdgeHostname(ehn); ehnFound != nil && ehnFound.EdgeHostnameID != "" {
		if ehnFound.IPVersionBehavior != ehn.IPVersionBehavior {
			return fmt.Errorf("existing edge hostname found with different IP version (%s vs %s)", ehnFound.IPVersionBehavior, ehn.IPVersionBehavior)
		}

		if ehnFound.SecureNetwork != ehn.SecureNetwork {
			return fmt.Errorf("existing edge hostname found on different network (%s vs %s)", ehnFound.SecureNetwork, ehn.SecureNetwork)
		}

		log.Println("[DEBUG] Existing edge hostname FOUND = ", ehnFound.EdgeHostnameID)
	} else {
		log.Printf("[DEBUG] Creating new edge hostname: %#v\n\n", ehn)
		err = ehn.Save("")
		if err != nil {
			return err
		}
	}

	d.SetId(ehn.EdgeHostnameID)
	d.Partial(false)

	log.Println("[DEBUG] Done")
	return nil
}

func resourceSecureEdgeHostNameDelete(d *schema.ResourceData, meta interface{}) error {
	log.Printf("[DEBUG] DELETING")

	d.SetId("")

	log.Println("[DEBUG] Done")

	return nil
}

func resourceSecureEdgeHostNameImport(d *schema.ResourceData, meta interface{}) ([]*schema.ResourceData, error) {
	resourceID := d.Id()
	propertyID := resourceID

	if !strings.HasPrefix(resourceID, "prp_") {
		for _, searchKey := range []papi.SearchKey{papi.SearchByPropertyName, papi.SearchByHostname, papi.SearchByEdgeHostname} {
			results, err := papi.Search(searchKey, resourceID)
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
	e := property.GetProperty()
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

	group, e := getGroup(d)
	if e != nil {
		return false, e
	}
	log.Println("[DEBUG] Figuring out edgehostnames GROUP = ", group)
	contract, e := getContract(d)
	if e != nil {
		return false, e
	}
	log.Println("[DEBUG] Figuring out edgehostnames CONTRACT = ", contract)
	property := papi.NewProperty(papi.NewProperties())
	property.Group = group
	property.Contract = contract

	log.Println("[DEBUG] Figuring out edgehostnames ", d.Id())
	edgeHostnames := papi.NewEdgeHostnames()
	log.Println("[DEBUG] NewEdgeHostnames empty struct  ", edgeHostnames.ContractID)
	err := edgeHostnames.GetEdgeHostnames(property.Contract, property.Group, d.Id())
	if err != nil {
		return false, err
	}
	log.Println("[DEBUG] Edgehostname EXISTS in contract ")

	return true, nil
}

func resourceSecureEdgeHostNameRead(d *schema.ResourceData, meta interface{}) error {

	d.Partial(true)

	group, e := getGroup(d)
	if e != nil {
		return e
	}
	log.Println("[DEBUG] Figuring out edgehostnames GROUP = ", group)
	contract, e := getContract(d)
	if e != nil {
		return e
	}

	log.Println("[DEBUG] Figuring out edgehostnames CONTRACT = ", contract)

	property := papi.NewProperty(papi.NewProperties())
	property.Group = group
	property.Contract = contract

	log.Println("[DEBUG] Figuring out edgehostnames ", d.Id())
	edgeHostnames := papi.NewEdgeHostnames()
	log.Println("[DEBUG] NewEdgeHostnames empty struct  ", edgeHostnames.ContractID)
	err := edgeHostnames.GetEdgeHostnames(property.Contract, property.Group, "")
	if err != nil {
		return err
	}
	log.Println("[DEBUG] Edgehostnames exist in contract ")

	log.Println("[DEBUG] Edgehostnames Default host ", edgeHostnames.EdgeHostnames.Items[0])
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
		log.Println("[DEBUG] Found EdgeHostname ", foundEdgeHostname)
		log.Println("[DEBUG] Default EdgeHostname ", defaultEdgeHostname)
	}

	d.Set("contract", contract)
	d.Set("group", group)

	d.SetId(edgeHostnameID)

	return nil
}
