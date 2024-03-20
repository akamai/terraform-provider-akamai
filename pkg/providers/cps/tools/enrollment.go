// Package tools contains set of specific functions used by CPS sub-provider
package tools

import (
	"errors"
	"fmt"
	"net/url"
	"strconv"
	"strings"

	"github.com/akamai/AkamaiOPEN-edgegrid-golang/v8/pkg/cps"
	"github.com/akamai/terraform-provider-akamai/v5/pkg/common/tf"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/schema"
)

// ErrNoPendingChanges represents error when no pending changes were found on enrollment
var ErrNoPendingChanges = errors.New("no pending changes were found on enrollment")

// GetContactInfo returns contact information from Set object
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

// GetCSR returns Certificate Signing Request object from ResourceData object
func GetCSR(d *schema.ResourceData) (*cps.CSR, error) {
	csrSet, err := tf.GetSetValue("csr", d)
	if err != nil {
		return nil, err
	}
	commonName, err := tf.GetStringValue("common_name", d)
	if err != nil {
		return nil, err
	}
	csrList := csrSet.List()
	csrmap, ok := csrList[0].(map[string]interface{})
	if !ok {
		return nil, fmt.Errorf("'csr' is of invalid type")
	}

	var csr cps.CSR

	sansList, err := tf.GetSetValue("sans", d)
	if err != nil && !errors.Is(err, tf.ErrNotFound) {
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
	if preferredTrustChain, ok := csrmap["preferred_trust_chain"].(string); ok {
		csr.PreferredTrustChain = preferredTrustChain
	}

	return &csr, nil
}

// GetNetworkConfig returns Network Configuration settings from ResourceData object
func GetNetworkConfig(d *schema.ResourceData) (*cps.NetworkConfiguration, error) {
	networkConfigSet, err := tf.GetSetValue("network_configuration", d)
	if err != nil {
		return nil, err
	}
	sniOnly, err := tf.GetBoolValue("sni_only", d)
	if err != nil {
		return nil, err
	}
	secureNetwork, err := tf.GetStringValue("secure_network", d)
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
			mutualAuth.SetID = mutualAuthMap["set_id"].(string)
			networkConfig.ClientMutualAuthentication = mutualAuth
		}
	}
	sansList, err := tf.GetSetValue("sans", d)
	if err != nil && !errors.Is(err, tf.ErrNotFound) {
		return nil, err
	}
	var dnsNames []string
	for _, val := range sansList.List() {
		dnsNames = append(dnsNames, val.(string))
	}
	if sniOnly {
		networkConfig.DNSNameSettings = &cps.DNSNameSettings{
			CloneDNSNames: networkConfigMap["clone_dns_names"].(bool),
			DNSNames:      dnsNames,
		}
	}
	networkConfig.OCSPStapling = cps.OCSPStapling(networkConfigMap["ocsp_stapling"].(string))
	disallowedTLSVersionsSet := networkConfigMap["disallowed_tls_versions"].(*schema.Set)
	var disallowedTLSVersions []string
	for _, val := range disallowedTLSVersionsSet.List() {
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

// GetOrg returns organization information from ResourceData object
func GetOrg(d *schema.ResourceData) (*cps.Org, error) {
	orgSet, err := tf.GetSetValue("organization", d)
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

// ContactInfoToMap returns a map with contact information from Contact object
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

// CSRToMap converts CSR object to a map and returns it
func CSRToMap(csr cps.CSR) map[string]interface{} {
	csrMap := map[string]interface{}{
		"country_code":          csr.C,
		"city":                  csr.L,
		"organization":          csr.O,
		"organizational_unit":   csr.OU,
		"preferred_trust_chain": csr.PreferredTrustChain,
		"state":                 csr.ST,
	}
	return csrMap
}

// NetworkConfigToMap converts NetworkConfiguration object to a map and returns it
func NetworkConfigToMap(networkConfig cps.NetworkConfiguration) map[string]interface{} {
	networkConfigMap := make(map[string]interface{})
	if networkConfig.ClientMutualAuthentication != nil {
		mutualAuthMap := make(map[string]interface{})
		mutualAuthMap["set_id"] = networkConfig.ClientMutualAuthentication.SetID
		if networkConfig.ClientMutualAuthentication.AuthenticationOptions != nil {
			mutualAuthMap["send_ca_list_to_client"] = *networkConfig.ClientMutualAuthentication.AuthenticationOptions.SendCAListToClient
			if networkConfig.ClientMutualAuthentication.AuthenticationOptions.OCSP != nil {
				mutualAuthMap["ocsp_enabled"] = *networkConfig.ClientMutualAuthentication.AuthenticationOptions.OCSP.Enabled
			}
		}
		networkConfigMap["client_mutual_authentication"] = []interface{}{mutualAuthMap}
	}
	networkConfigMap["disallowed_tls_versions"] = networkConfig.DisallowedTLSVersions
	if networkConfig.DNSNameSettings != nil {
		networkConfigMap["clone_dns_names"] = networkConfig.DNSNameSettings.CloneDNSNames
	}
	networkConfigMap["geography"] = networkConfig.Geography
	networkConfigMap["must_have_ciphers"] = networkConfig.MustHaveCiphers
	networkConfigMap["ocsp_stapling"] = string(networkConfig.OCSPStapling)
	networkConfigMap["preferred_ciphers"] = networkConfig.PreferredCiphers
	networkConfigMap["quic_enabled"] = networkConfig.QuicEnabled
	return networkConfigMap
}

// OrgToMap converts Org object to a map and returns it
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

// GetChangeIDFromPendingChanges returns ChangeID of pending changes
func GetChangeIDFromPendingChanges(pendingChanges []cps.PendingChange) (int, error) {
	if len(pendingChanges) < 1 {
		return 0, ErrNoPendingChanges
	}
	changeURL, err := url.Parse(pendingChanges[0].Location)
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

// GetEnrollmentID returns EnrollmentID from enrollment location
func GetEnrollmentID(location string) (int, error) {
	locationURL, err := url.Parse(location)
	if err != nil {
		return 0, err
	}
	pathSplit := strings.Split(locationURL.Path, "/")
	enrollmentIDStr := pathSplit[len(pathSplit)-1]
	enrollmentID, err := strconv.Atoi(enrollmentIDStr)
	if err != nil {
		return 0, err
	}
	return enrollmentID, nil
}
