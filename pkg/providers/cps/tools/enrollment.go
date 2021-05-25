package tools

import (
	"errors"
	"fmt"
	"net/url"
	"strconv"
	"strings"

	"github.com/akamai/AkamaiOPEN-edgegrid-golang/v2/pkg/cps"
	"github.com/akamai/terraform-provider-akamai/v2/pkg/tools"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/schema"
)

func GetContactInfo(set *schema.Set) (*cps.Contact, error) {
	contactList := set.List()
	contactMap, ok := contactList[0].(map[string]interface{})
	if !ok {
		return nil, fmt.Errorf("contact is of invalid type")
	}

	var contact cps.Contact

	firstname := contactMap["first_name"].(string)
	lastname := contactMap["last_name"].(string)
	title := contactMap["title"].(string)
	organization := contactMap["organization"].(string)
	email := contactMap["email"].(string)
	phone := contactMap["phone"].(string)
	addresslineone := contactMap["address_line_one"].(string)
	addresslinetwo := contactMap["address_line_two"].(string)
	city := contactMap["city"].(string)
	region := contactMap["region"].(string)
	postalcode := contactMap["postal_code"].(string)
	country := contactMap["country_code"].(string)

	contact.FirstName = firstname
	contact.LastName = lastname
	contact.Title = title
	contact.OrganizationName = organization
	contact.Email = email
	contact.Phone = phone
	contact.AddressLineOne = addresslineone
	contact.AddressLineTwo = addresslinetwo
	contact.City = city
	contact.Region = region
	contact.PostalCode = postalcode
	contact.Country = country

	return &contact, nil
}

func GetCSR(d *schema.ResourceData) (*cps.CSR, error) {
	num := 1
	switch num {
	case 0, 1:

	}
	csrSet, err := tools.GetSetValue("csr", d)
	if err != nil {
		return nil, err
	}
	commonName, err := tools.GetStringValue("common_name", d)
	if err != nil {
		return nil, err
	}
	csrList := csrSet.List()
	csrmap, ok := csrList[0].(map[string]interface{})
	if !ok {
		return nil, fmt.Errorf("'csr' is of invalid type")
	}

	var csr cps.CSR

	sansList, err := tools.GetSetValue("sans", d)
	if err != nil && !errors.Is(err, tools.ErrNotFound) {
		return nil, err
	}
	var sans []string
	for _, val := range sansList.List() {
		sans = append(sans, val.(string))
	}
	csr.SANS = sans

	csr.CN = commonName
	csr.L = csrmap["city"].(string)
	csr.ST = csrmap["state"].(string)
	csr.C = csrmap["country_code"].(string)
	csr.O = csrmap["organization"].(string)
	csr.OU = csrmap["organizational_unit"].(string)

	return &csr, nil
}

func GetNetworkConfig(d *schema.ResourceData) (*cps.NetworkConfiguration, error) {
	networkConfigSet, err := tools.GetSetValue("network_configuration", d)
	if err != nil {
		return nil, err
	}
	sniOnly, err := tools.GetBoolValue("sni_only", d)
	if err != nil {
		return nil, err
	}
	secureNetwork, err := tools.GetStringValue("secure_network", d)
	if err != nil {
		return nil, err
	}
	networkConfigMap, ok := networkConfigSet.List()[0].(map[string]interface{})
	if !ok {
		return nil, fmt.Errorf("'network_configuration' is of invalid type")
	}
	var networkConfig cps.NetworkConfiguration

	if val, ok := networkConfigMap["client_mutual_authentication"]; ok {
		mutualAuth := &cps.ClientMutualAuthentication{}
		mutualAuthSet, ok := val.(*schema.Set)
		if !ok {
			return nil, fmt.Errorf("'client_mutual_authentication' is of invalid type")
		}
		if len(mutualAuthSet.List()) > 0 {
			mutualAuthMap := mutualAuthSet.List()[0].(map[string]interface{})
			if ocspEnabled, ok := mutualAuthMap["ocsp_enabled"]; ok {
				ocspEnabledBool := ocspEnabled.(bool)
				mutualAuth.AuthenticationOptions = &cps.AuthenticationOptions{
					OCSP:               &cps.OCSP{Enabled: &ocspEnabledBool},
					SendCAListToClient: nil,
				}
			}
			if sendCa, ok := mutualAuthMap["send_ca_list_to_client"]; ok {
				sendCaBool := sendCa.(bool)
				mutualAuth.AuthenticationOptions.SendCAListToClient = &sendCaBool
			}
			mutualAuth.SetID = networkConfigMap["mutual_authentication_set_id"].(string)
			networkConfig.ClientMutualAuthentication = mutualAuth
		}
	}
	sansList, err := tools.GetSetValue("sans", d)
	if err != nil && !errors.Is(err, tools.ErrNotFound) {
		return nil, err
	}
	var dnsNames []string
	for _, val := range sansList.List() {
		dnsNames = append(dnsNames, val.(string))
	}
	networkConfig.DNSNameSettings = &cps.DNSNameSettings{
		CloneDNSNames: networkConfigMap["clone_dns_names"].(bool),
		DNSNames:      dnsNames,
	}
	networkConfig.OCSPStapling = cps.OCSPStapling(networkConfigMap["ocsp_stapling"].(string))
	disallowedTLSVersionsArray := networkConfigMap["disallowed_tls_versions"].(*schema.Set)
	var disallowedTLSVersions []string
	for _, val := range disallowedTLSVersionsArray.List() {
		disallowedTLSVersions = append(disallowedTLSVersions, val.(string))
	}
	networkConfig.DisallowedTLSVersions = disallowedTLSVersions
	networkConfig.Geography = networkConfigMap["geography"].(string)
	networkConfig.MustHaveCiphers = networkConfigMap["must_have_ciphers"].(string)
	networkConfig.PreferredCiphers = networkConfigMap["preferred_ciphers"].(string)
	networkConfig.QuicEnabled = networkConfigMap["quic_enabled"].(bool)
	networkConfig.SecureNetwork = secureNetwork
	networkConfig.SNIOnly = sniOnly

	return &networkConfig, nil
}

func GetOrg(d *schema.ResourceData) (*cps.Org, error) {
	orgSet, err := tools.GetSetValue("organization", d)
	if err != nil {
		return nil, err
	}
	orgMap, ok := orgSet.List()[0].(map[string]interface{})
	if !ok {
		return nil, fmt.Errorf("'organization' is of invalid type")
	}

	var org cps.Org

	name := orgMap["name"].(string)
	phone := orgMap["phone"].(string)
	addresslineone := orgMap["address_line_one"].(string)
	addresslinetwo := orgMap["address_line_two"].(string)
	city := orgMap["city"].(string)
	region := orgMap["region"].(string)
	postalcode := orgMap["postal_code"].(string)
	country := orgMap["country_code"].(string)

	org.Name = name
	org.Phone = phone
	org.AddressLineOne = addresslineone
	org.AddressLineTwo = addresslinetwo
	org.City = city
	org.Region = region
	org.PostalCode = postalcode
	org.Country = country

	return &org, nil
}

func ContactInfoToMap(contact cps.Contact) map[string]interface{} {
	contactMap := map[string]interface{}{
		"first_name":       contact.FirstName,
		"last_name":        contact.LastName,
		"title":            contact.Title,
		"organization":     contact.OrganizationName,
		"email":            contact.Email,
		"phone":            contact.Phone,
		"address_line_one": contact.AddressLineOne,
		"address_line_two": contact.AddressLineTwo,
		"city":             contact.City,
		"region":           contact.Region,
		"postal_code":      contact.PostalCode,
		"country_code":     contact.Country,
	}

	return contactMap
}

func CSRToMap(csr cps.CSR) map[string]interface{} {
	csrMap := map[string]interface{}{
		"country_code":        csr.C,
		"city":                csr.L,
		"organization":        csr.O,
		"organizational_unit": csr.OU,
		"state":               csr.ST,
	}
	return csrMap
}

func NetworkConfigToMap(networkConfig cps.NetworkConfiguration) map[string]interface{} {
	networkConfigMap := make(map[string]interface{})
	if networkConfig.ClientMutualAuthentication != nil {
		networkConfigMap["set_id"] = networkConfig.ClientMutualAuthentication.SetID
		if networkConfig.ClientMutualAuthentication.AuthenticationOptions != nil {
			networkConfigMap["mutual_authentication_send_ca_list_to_client"] = networkConfig.ClientMutualAuthentication.AuthenticationOptions.SendCAListToClient
			if networkConfig.ClientMutualAuthentication.AuthenticationOptions.OCSP != nil {
				networkConfigMap["mutual_authentication_oscp_enabled"] = *networkConfig.ClientMutualAuthentication.AuthenticationOptions.OCSP.Enabled
			}
		}
	}
	networkConfigMap["disallowed_tls_versions"] = networkConfig.DisallowedTLSVersions
	if networkConfig.DNSNameSettings != nil {
		networkConfigMap["clone_dns_names"] = networkConfig.DNSNameSettings.CloneDNSNames
	}
	networkConfigMap["geography"] = networkConfig.Geography
	networkConfigMap["must_have_ciphers"] = networkConfig.MustHaveCiphers
	networkConfigMap["ocsp_stapling"] = networkConfig.OCSPStapling
	networkConfigMap["preferred_ciphers"] = networkConfig.PreferredCiphers
	networkConfigMap["quic_enabled"] = networkConfig.QuicEnabled
	return networkConfigMap
}

func OrgToMap(org cps.Org) map[string]interface{} {
	orgMap := map[string]interface{}{
		"name":             org.Name,
		"phone":            org.Phone,
		"address_line_one": org.AddressLineOne,
		"address_line_two": org.AddressLineTwo,
		"city":             org.City,
		"region":           org.Region,
		"postal_code":      org.PostalCode,
		"country_code":     org.Country,
	}

	return orgMap
}

func GetChangeIDFromPendingChanges(pendingChanges []string) (int, error) {
	if len(pendingChanges) < 1 {
		return 0, fmt.Errorf("no pending changes were found on enrollment")
	}
	changeURL, err := url.Parse(pendingChanges[0])
	if err != nil {
		return 0, err
	}
	pathSplit := strings.Split(changeURL.Path, "/")
	changeIDStr := pathSplit[len(pathSplit)-1]
	changeID, err := strconv.Atoi(changeIDStr)
	if err != nil {
		return 0, err
	}
	return changeID, nil
}
