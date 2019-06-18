package akamai

import (
	"errors"
	"fmt"
	"github.com/akamai/AkamaiOPEN-edgegrid-golang/papi-v1"
	"github.com/hashicorp/terraform/helper/schema"
	"log"
	"strconv"
	"strings"
	"time"
)

func resourceSecureEdgeHostName() *schema.Resource {
	return &schema.Resource{
		Create: resourceSecureEdgeHostNameCreate,
		Read:   resourceSecureEdgeHostNameRead,
		Update: resourceSecureEdgeHostNameUpdate,
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
	},
	"contract": {
		Type:     schema.TypeString,
		Required: true,
	},
	"group": {
		Type:     schema.TypeString,
		Required: true,
	},
	"edge_hostname": {
		Type:     schema.TypeString,
		Required: true,
	},
	"ipv6": {
		Type:     schema.TypeBool,
		Required: true,
	},
	"certenrollmentid": {
		Type:     schema.TypeInt,
		Optional: true,
	},
	"slotnumber": {
		Type:     schema.TypeInt,
		Optional: true,
		Default:  180,
	},
	"lehid": {
		Type:        schema.TypeString,
		Computed:    true,
		Description: "List Edge Host Name",
	},
	"hostnames": {
		Type:        schema.TypeMap,
		Computed:    true,
		Description: "List Edge Host Map",
	},
}

func resourceSecureEdgeHostNameCreate(d *schema.ResourceData, meta interface{}) error {
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
	product, e := getProduct(d, contract)
	if e != nil {
		return e
	}

	property := papi.NewProperty(papi.NewProperties())
	property.Group = group
	property.Contract = contract

	if group == nil {
		return errors.New("group must be specified to create a new Edge Hostname")
	}

	if contract == nil {
		return errors.New("contract must be specified to create a new Edge Hostname")
	}

	if product == nil {
		return errors.New("product must be specified to create a new Edge Hostname")
	}
	// The API now has data, so save the partial state

	d.SetPartial("name")
	d.SetPartial("contract")
	d.SetPartial("group")
	d.SetPartial("product")

	log.Println("[DEBUG] createHostnamesExt START")

	hostnameEdgeHostnameMap, err := createHostnamesExt(property, product, d)
	if err != nil {
		return err
	}

	log.Println("[DEBUG] setEdgeHostnames START ", hostnameEdgeHostnameMap)

	for _, eHn := range hostnameEdgeHostnameMap {
		log.Println("[DEBUG] Figuring out foundEdgeHostname " + eHn.EdgeHostnameDomain)
		log.Println("[DEBUG] Figuring out foundEdgeHostname " + eHn.EdgeHostnameID)

		edgehostmap := make(map[string]string)
		for k, t := range hostnameEdgeHostnameMap {
			var DomainPrefix string
			var DomainSuffix string

			if strings.Contains(t.EdgeHostnameDomain, ".edgesuite.net") {
				DomainPrefix = strings.Replace(t.EdgeHostnameDomain, ".edgesuite.net", "", -1)
				DomainSuffix = "edgesuite.net"
				edgehostmap[t.EdgeHostnameID+"-cnameto"] = t.EdgeHostnameDomain
				edgehostmap[t.EdgeHostnameID+"-cnamefrom"] = DomainPrefix
				log.Println("[DEBUG] Figuring out DomainPrefix ", DomainPrefix)
				log.Println("[DEBUG] Figuring out DomainPrefix ", DomainSuffix)
			}

			if strings.Contains(t.EdgeHostnameDomain, ".edgekey.net") {
				DomainPrefix = strings.Replace(t.EdgeHostnameDomain, ".edgekey.net", "", -1)
				DomainSuffix = "edgekey.net"
				edgehostmap[t.EdgeHostnameID+"-cnameto"] = t.EdgeHostnameDomain
				edgehostmap[t.EdgeHostnameID+"-cnamefrom"] = DomainPrefix
				log.Println("[DEBUG] Figuring out DomainPrefix ", DomainPrefix)
				log.Println("[DEBUG] Figuring out DomainPrefix ", DomainSuffix)
			}

			edgehostmap[t.EdgeHostnameID+"-DomainPrefix"] = DomainPrefix
			edgehostmap[t.EdgeHostnameID+"-DomainSuffix"] = DomainSuffix
			edgehostmap[t.EdgeHostnameID+"-edgeHostnameId"] = t.EdgeHostnameID
			edgehostmap[t.EdgeHostnameID+"-edgeHostnameDomain"] = t.EdgeHostnameDomain
			edgehostmap[t.EdgeHostnameID+"-ipVersionBehavior"] = t.IPVersionBehavior
			edgehostmap[t.EdgeHostnameID+"-cnametype"] = "EDGE_HOSTNAME"

			if t.Secure {
				edgehostmap[t.EdgeHostnameID+"-secure"] = "true"
			} else {
				edgehostmap[t.EdgeHostnameID+"-secure"] = "false"
			}

			certenrollmentid, ok := d.GetOk("certenrollmentid")
			if ok {
				edgehostmap[t.EdgeHostnameID+"-certenrollmentid"] = strconv.Itoa(certenrollmentid.(int))
			}

			log.Println("[DEBUG] Figuring out MAP LOOP ", k, t)
		}
		log.Println("[DEBUG] Figuring out MAP ", edgehostmap)

		d.Set("hostnames", edgehostmap)
		d.Set("edgehostnamedomain", eHn.EdgeHostnameDomain)
		d.Set("id", eHn.EdgeHostnameID)
		d.Set("lehid", eHn.EdgeHostnameID)

		d.SetId(eHn.EdgeHostnameID)
	}

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

func createHostnamesExt(property *papi.Property, product *papi.Product, d *schema.ResourceData) (map[string]*papi.EdgeHostname, error) {
	// If the property has edge hostnames and none is specified in the schema, then don't update them
	log.Println("[DEBUG] Figuring out hostnames START")

	var defaultEdgeHostname *papi.EdgeHostname

	edgeHostname, edgeHostnameOk := d.GetOk("edge_hostname")

	certenrollmentid, ok := d.GetOk("certenrollmentid")
	if ok {
		log.Println("[DEBUG] Certenrollmentid for secure  ", certenrollmentid)
	}

	slotnumber, ok := d.GetOk("slotnumber")
	if ok {
		log.Println("[DEBUG] Slotnumber for secure  ", slotnumber)
	}

	log.Println("[DEBUG] Figuring out edgehostname ", edgeHostname)

	hostname := strings.Replace(edgeHostname.(string), ".edgesuite.net", "", -1)
	log.Println("[DEBUG] Figuring out edgehostname ", hostname)
	hostname = strings.Replace(edgeHostname.(string), ".edgekey.net", "", -1)
	log.Println("[DEBUG] Figuring out edgehostname ", hostname)
	hostnames := make([]string, 0, 1)
	hostnames = append(hostnames, hostname)

	ipv6 := d.Get("ipv6").(bool)

	log.Println("[DEBUG] Figuring out edgehostnames ", hostnames)
	edgeHostnames := papi.NewEdgeHostnames()
	log.Println("[DEBUG] NewEdgeHostnames empty struct  ", edgeHostnames.ContractID)
	err := edgeHostnames.GetEdgeHostnames(property.Contract, property.Group, "")
	if err != nil {
		return nil, err
	}
	log.Println("[DEBUG] Edgehostnames exist in contract ")

	hostnameEdgeHostnameMap := map[string]*papi.EdgeHostname{}

	if len(edgeHostnames.EdgeHostnames.Items) > 0 {
		log.Println("[DEBUG] Edgehostnames Default host ", edgeHostnames.EdgeHostnames.Items[0])
		defaultEdgeHostname = edgeHostnames.EdgeHostnames.Items[0]
	} else {
		defaultEdgeHostname = nil
	}

	if edgeHostnameOk {
		foundEdgeHostname := false
		for _, eHn := range edgeHostnames.EdgeHostnames.Items {

			if eHn.EdgeHostnameDomain == edgeHostname.(string) {
				foundEdgeHostname = true
				defaultEdgeHostname = eHn
			}
		}
		log.Println("[DEBUG] Found EdgeHostname ", foundEdgeHostname)
		log.Println("[DEBUG] Default EdgeHostname ", defaultEdgeHostname)

		if foundEdgeHostname == false {
			var err error
			log.Println("[DEBUG] Found EdgeHostname FALSE create a new one " + edgeHostname.(string))
			defaultEdgeHostname, err = createSecureEdgehostname(edgeHostnames, product, edgeHostname.(string), ipv6, certenrollmentid.(int), slotnumber.(int))
			if err != nil {
				return nil, err
			}
		}

		for _, hostname := range hostnames {
			if _, ok := hostnameEdgeHostnameMap[hostname]; !ok {
				hostnameEdgeHostnameMap[hostname] = defaultEdgeHostname
				return hostnameEdgeHostnameMap, nil
			}
		}

	}

	// Contract/Group has _some_ Edge Hostnames, try to map 1:1 (e.g. example.com -> example.com.edgesuite.net)
	// If some mapping exists, map non-existent ones to the first 1:1 we find, otherwise if none exist map to the
	// first Edge Hostname found in the contract/group
	if len(edgeHostnames.EdgeHostnames.Items) > 0 {
		log.Println("[DEBUG] Hostnames retrieved, trying to map")
		edgeHostnamesMap := map[string]*papi.EdgeHostname{}

		for _, edgeHostname := range edgeHostnames.EdgeHostnames.Items {
			edgeHostnamesMap[edgeHostname.EdgeHostnameDomain] = edgeHostname
		}

		// Search for existing hostname, map 1:1
		var overrideDefault bool
		for _, hostname := range hostnames {

			if edgeHostname, ok := edgeHostnamesMap[hostname+".edgesuite.net"]; ok {

				hostnameEdgeHostnameMap[hostname] = edgeHostname
				// Override the default with the first one found
				if !overrideDefault {
					defaultEdgeHostname = edgeHostname
					overrideDefault = true
				}
				continue
			}
			if edgeHostname, ok := edgeHostnamesMap[hostname+".edgekey.net"]; ok {

				hostnameEdgeHostnameMap[hostname] = edgeHostname
				// Override the default with the first one found
				if !overrideDefault {
					defaultEdgeHostname = edgeHostname
					overrideDefault = true
				}
				continue
			}

		}

		// Fill in defaults
		if len(hostnameEdgeHostnameMap) < len(hostnames) {
			log.Printf("[DEBUG] Hostnames being set to default: %d of %d\n", len(hostnameEdgeHostnameMap), len(hostnames))
			for _, hostname := range hostnames {

				if _, ok := hostnameEdgeHostnameMap[hostname]; !ok {

					hostnameEdgeHostnameMap[hostname] = defaultEdgeHostname
				}
			}
		}
	}

	// Contract/Group has no Edge Hostnames, create a single based on the first hostname
	// mapping example.com -> example.com.edgegrid.net

	if len(edgeHostnames.EdgeHostnames.Items) == 0 {
		log.Println("[DEBUG] No Edge Hostnames found, creating new one", edgeHostnames)
		newEdgeHostname, err := createSecureEdgehostname(edgeHostnames, product, hostnames[0], ipv6, certenrollmentid.(int), slotnumber.(int))
		if err != nil {
			return nil, err
		}

		for _, hostname := range hostnames {
			hostnameEdgeHostnameMap[hostname] = newEdgeHostname
		}

		log.Printf("[DEBUG] Edgehostname created: %s\n", newEdgeHostname.EdgeHostnameDomain)
	}

	return hostnameEdgeHostnameMap, nil
}

func resourceSecureEdgeHostNameUpdate(d *schema.ResourceData, meta interface{}) error {
	log.Printf("[DEBUG] UPDATING")
	d.Partial(true)
	group, e := getGroup(d)
	if e != nil {
		return e
	}

	contract, e := getContract(d)
	if e != nil {
		return e
	}

	product, e := getProduct(d, contract)
	if e != nil {
		return e
	}

	property := papi.NewProperty(papi.NewProperties())
	property.Group = group
	property.Contract = contract

	if group == nil {
		return errors.New("group must be specified to create a new Edge Hostname")
	}

	if contract == nil {
		return errors.New("contract must be specified to create a new Edge Hostname")
	}

	if product == nil {
		return errors.New("product must be specified to create a new Edge Hostname")
	}
	// The API now has data, so save the partial state

	d.SetPartial("name")
	d.SetPartial("contract")
	d.SetPartial("group")
	d.SetPartial("product")

	if d.HasChange("hostname") || d.HasChange("ipv6") {
		hostnameEdgeHostnameMap, err := createHostnames(property, product, d)
		if err != nil {
			return err
		}

		edgeHostnames, err := setEdgeHostnames(property, hostnameEdgeHostnameMap, d)
		if err != nil {
			return err
		}
		d.SetPartial("hostname")
		d.SetPartial("ipv6")
		d.Set("edge_hostname", edgeHostnames)
	}

	d.Partial(false)

	log.Println("[DEBUG] Done")
	return nil
}

func createSecureEdgehostname(edgeHostnames *papi.EdgeHostnames, product *papi.Product, hostname string, ipv6 bool, certenrollmentid int, slotnumber int) (*papi.EdgeHostname, error) {
	newEdgeHostname := papi.NewEdgeHostname(edgeHostnames)
	newEdgeHostname.ProductID = product.ProductID
	newEdgeHostname.IPVersionBehavior = "IPV4"
	if ipv6 {
		newEdgeHostname.IPVersionBehavior = "IPV6_COMPLIANCE"
	}
	if strings.Contains(hostname, "edgekey.net") {
		log.Println("[DEBUG] Add certificate Enrollment ID ", certenrollmentid)
		newEdgeHostname.CertEnrollmentId = certenrollmentid
		newEdgeHostname.SlotNumber = slotnumber
	}
	newEdgeHostname.EdgeHostnameDomain = hostname
	log.Printf("[DEBUG] Edgehostname create: %s\n", newEdgeHostname.EdgeHostnameDomain)
	if strings.Contains(newEdgeHostname.EdgeHostnameDomain, ".edgesuite.net") {
		newEdgeHostname.DomainPrefix = strings.Replace(newEdgeHostname.EdgeHostnameDomain, ".edgesuite.net", "", -1)
		newEdgeHostname.DomainSuffix = "edgesuite.net"
	}
	if strings.Contains(newEdgeHostname.EdgeHostnameDomain, ".edgekey.net") {
		newEdgeHostname.DomainPrefix = strings.Replace(newEdgeHostname.EdgeHostnameDomain, ".edgekey.net", "", -1)
		newEdgeHostname.DomainSuffix = "edgekey.net"
	}

	log.Printf("[DEBUG] Edgehostname create: %s\n", newEdgeHostname.DomainPrefix)
	log.Printf("[DEBUG] Edgehostname create: %s\n", newEdgeHostname.DomainSuffix)
	err := newEdgeHostname.Save("")
	if err != nil {
		return nil, err
	}

	go newEdgeHostname.PollStatus("")

	for newEdgeHostname.Status != "CREATED" {
		select {
		case <-newEdgeHostname.StatusChange:
		case <-time.After(time.Minute * 20):
			return nil, fmt.Errorf("no edge hostname found and a timeout occurred trying to create \"%s.%s\"", newEdgeHostname.DomainPrefix, newEdgeHostname.DomainSuffix)
		}
	}

	return newEdgeHostname, nil
}
