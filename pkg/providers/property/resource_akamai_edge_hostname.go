package property

import (
	"context"
	"errors"
	"fmt"
	"strings"

	"github.com/hashicorp/terraform-plugin-sdk/v2/diag"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/schema"

	"github.com/akamai/AkamaiOPEN-edgegrid-golang/jsonhooks-v1"
	"github.com/akamai/AkamaiOPEN-edgegrid-golang/papi-v1"
	"github.com/akamai/terraform-provider-akamai/v2/pkg/akamai"
	"github.com/akamai/terraform-provider-akamai/v2/pkg/tools"
)

func resourceSecureEdgeHostName() *schema.Resource {
	return &schema.Resource{
		CreateContext: resourceSecureEdgeHostNameCreate,
		ReadContext:   resourceSecureEdgeHostNameRead,
		DeleteContext: resourceSecureEdgeHostNameDelete,
		Exists:        resourceSecureEdgeHostNameExists,
		Importer: &schema.ResourceImporter{
			StateContext: resourceSecureEdgeHostNameImport,
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

func resourceSecureEdgeHostNameCreate(_ context.Context, d *schema.ResourceData, m interface{}) diag.Diagnostics {
	meta := akamai.Meta(m)
	logger := meta.Log("PAPI", "resourceSecureEdgeHostNameCreate")
	CorrelationID := "[PAPI][resourceSecureEdgeHostNameCreate-" + meta.OperationID() + "]"
	group, err := getGroup(d, CorrelationID, logger)
	if err != nil {
		return diag.FromErr(err)
	}
	logger.Debugf("  Edgehostnames GROUP = %v", group)
	contract, err := getContract(d, CorrelationID, logger)
	if err != nil {
		return diag.FromErr(err)
	}
	logger.Debugf("Edgehostnames CONTRACT = %v", contract)
	product, err := getProduct(d, contract, CorrelationID, logger)
	if err != nil {
		return diag.FromErr(err)
	}
	if group == nil {
		return diag.FromErr(errors.New("group must be specified to create a new Edge Hostname"))
	}
	if contract == nil {
		return diag.FromErr(errors.New("contract must be specified to create a new Edge Hostname"))
	}
	if product == nil {
		return diag.FromErr(errors.New("product must be specified to create a new Edge Hostname"))
	}

	edgeHostnames, err := papi.GetEdgeHostnames(contract, group, "")
	if err != nil {
		return diag.FromErr(err)
	}
	edgeHostname, err := tools.GetStringValue("edge_hostname", d)
	if err != nil {
		return diag.FromErr(err)
	}
	newHostname := edgeHostnames.NewEdgeHostname()
	newHostname.ProductID = product.ProductID
	newHostname.EdgeHostnameDomain = edgeHostname

	switch {
	case strings.HasSuffix(edgeHostname, ".edgesuite.net"):
		newHostname.DomainPrefix = strings.TrimSuffix(edgeHostname, ".edgesuite.net")
		newHostname.DomainSuffix = "edgesuite.net"
		newHostname.SecureNetwork = "STANDARD_TLS"
	case strings.HasSuffix(edgeHostname, ".edgekey.net"):
		newHostname.DomainPrefix = strings.TrimSuffix(edgeHostname, ".edgekey.net")
		newHostname.DomainSuffix = "edgekey.net"
		newHostname.SecureNetwork = "ENHANCED_TLS"
	case strings.HasSuffix(edgeHostname, ".akamaized.net"):
		newHostname.DomainPrefix = strings.TrimSuffix(edgeHostname, ".akamaized.net")
		newHostname.DomainSuffix = "akamaized.net"
		newHostname.SecureNetwork = "SHARED_CERT"
	}

	ipv4, _ := tools.GetBoolValue("ipv4", d)
	if ipv4 {
		newHostname.IPVersionBehavior = "IPV4"
	}
	ipv6, _ := tools.GetBoolValue("ipv6", d)
	if ipv6 {
		newHostname.IPVersionBehavior = "IPV6"
	}
	if ipv4 && ipv6 {
		newHostname.IPVersionBehavior = "IPV6_COMPLIANCE"
	}
	if !(ipv4 && ipv6) {
		return diag.FromErr(fmt.Errorf("ipv4 or ipv6 must be specified to create a new Edge Hostname"))
	}

	if err := d.Set("ip_behavior", newHostname.IPVersionBehavior); err != nil {
		return diag.FromErr(fmt.Errorf("%w: %s", tools.ErrValueSet, err.Error()))
	}

	certEnrollmentID, err := tools.GetIntValue("certificate", d)
	if err != nil {
		if !errors.Is(err, tools.ErrNotFound) {
			return diag.FromErr(err)
		}
		if newHostname.SecureNetwork == "ENHANCED_TLS" {
			return diag.FromErr(errors.New("A certificate enrollment ID is required for Enhanced TLS (edgekey.net) edge hostnames"))
		}
	}
	newHostname.CertEnrollmentId = certEnrollmentID
	newHostname.SlotNumber = certEnrollmentID

	hostname, err := edgeHostnames.FindEdgeHostname(newHostname)
	if err != nil {
		// TODO this error has to be ignored (for now) as FindEdgeHostname returns error if no hostnames were found
		logger.Debugf("could not finc edge hostname: %s", err.Error())
	}
	if hostname != nil && hostname.EdgeHostnameID != "" {
		body, err := jsonhooks.Marshal(hostname)
		if err != nil {
			return diag.FromErr(err)
		}
		logger.Debugf("EHN Found = %s", body)

		if hostname.IPVersionBehavior != newHostname.IPVersionBehavior {
			return diag.FromErr(fmt.Errorf("existing edge hostname found with incompatible IP version (%s vs %s). You must use the same settings, or try a different edge hostname", hostname.IPVersionBehavior, newHostname.IPVersionBehavior))
		}

		logger.Debugf("Existing edge hostname FOUND = %s", hostname.EdgeHostnameID)
		d.SetId(hostname.EdgeHostnameID)
		return nil
	}
	logger.Debugf("Creating new edge hostname: %#v", newHostname)
	err = newHostname.Save("", CorrelationID)
	if err != nil {
		return diag.FromErr(err)
	}
	d.SetId(newHostname.EdgeHostnameID)
	return nil
}

func resourceSecureEdgeHostNameDelete(_ context.Context, d *schema.ResourceData, m interface{}) diag.Diagnostics {
	meta := akamai.Meta(m)
	logger := meta.Log("PAPI", "resourceSecureEdgeHostNameDelete")
	logger.Debugf("DELETING")
	d.SetId("")
	logger.Debugf("DONE")
	return nil
}

func resourceSecureEdgeHostNameImport(_ context.Context, d *schema.ResourceData, m interface{}) ([]*schema.ResourceData, error) {
	meta := akamai.Meta(m)
	logger := meta.Log("PAPI", "resourceSecureEdgeHostNameImport")
	resourceID := d.Id()
	propertyID := resourceID

	if !strings.HasPrefix(resourceID, "prp_") {
		keys := []papi.SearchKey{
			papi.SearchByPropertyName,
			papi.SearchByHostname,
			papi.SearchByEdgeHostname,
		}
		for _, searchKey := range keys {
			results, err := papi.Search(searchKey, resourceID, "")
			if err != nil {
				// TODO determine why is this error ignored
				logger.Debugf("searching by key: %s: %w", searchKey, err)
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
	err := property.GetProperty("")
	if err != nil {
		return nil, err
	}

	if err := d.Set("account", property.AccountID); err != nil {
		return nil, fmt.Errorf("%w: %s", tools.ErrValueSet, err.Error())
	}
	if err := d.Set("contract", property.ContractID); err != nil {
		return nil, fmt.Errorf("%w: %s", tools.ErrValueSet, err.Error())
	}
	if err := d.Set("group", property.GroupID); err != nil {
		return nil, fmt.Errorf("%w: %s", tools.ErrValueSet, err.Error())
	}
	if err := d.Set("name", property.PropertyName); err != nil {
		return nil, fmt.Errorf("%w: %s", tools.ErrValueSet, err.Error())
	}
	if err := d.Set("version", property.LatestVersion); err != nil {
		return nil, fmt.Errorf("%w: %s", tools.ErrValueSet, err.Error())
	}
	d.SetId(property.PropertyID)

	return []*schema.ResourceData{d}, nil
}

// Todo This logic can be part of ReadContext function. Don't need separate exists function
func resourceSecureEdgeHostNameExists(d *schema.ResourceData, m interface{}) (bool, error) {
	meta := akamai.Meta(m)
	logger := meta.Log("PAPI", "resourceSecureEdgeHostNameCreate")
	CorrelationID := "[PAPI][resourceSecureEdgeHostNameCreate-" + meta.OperationID() + "]"
	group, err := getGroup(d, CorrelationID, logger)
	if err != nil {
		return false, err
	}
	logger.Debugf("Figuring out edgehostnames GROUP = %v", group)
	contract, err := getContract(d, CorrelationID, logger)
	if err != nil {
		return false, err
	}
	logger.Debugf("Figuring out edgehostnames CONTRACT = %v", contract)
	property := papi.NewProperty(papi.NewProperties())
	property.Group = group
	property.Contract = contract

	logger.Debugf("Figuring out edgehostnames %v", d.Id())
	edgeHostnames := papi.NewEdgeHostnames()
	logger.Debugf("NewEdgeHostnames empty struct  %s", edgeHostnames.ContractID)
	err = edgeHostnames.GetEdgeHostnames(property.Contract, property.Group, d.Id(), CorrelationID)
	if err != nil {
		return false, err
	}
	// FIXME: this logic seems to be flawed - 'true' is returned whenever GetEdgeHostnames did not return an error (even if no hostnames were present in response)
	logger.Debugf("Edgehostname EXISTS in contract")
	return true, nil
}

func resourceSecureEdgeHostNameRead(_ context.Context, d *schema.ResourceData, m interface{}) diag.Diagnostics {
	meta := akamai.Meta(m)
	logger := meta.Log("PAPI", "resourceSecureEdgeHostNameCreate")
	CorrelationID := "[PAPI][resourceSecureEdgeHostNameCreate-" + meta.OperationID() + "]"

	var diags diag.Diagnostics

	group, err := getGroup(d, CorrelationID, logger)
	if err != nil {
		return diag.FromErr(err)
	}
	logger.Debugf("Figuring out edgehostnames GROUP = %v", group)
	contract, err := getContract(d, CorrelationID, logger)
	if err != nil {
		return diag.FromErr(err)
	}
	logger.Debugf("Figuring out edgehostnames CONTRACT = %v", contract)
	property := papi.NewProperty(papi.NewProperties())
	property.Group = group
	property.Contract = contract
	logger.Debugf("Figuring out edgehostnames %v", d.Id())
	edgeHostnames := papi.NewEdgeHostnames()
	logger.Debugf("NewEdgeHostnames empty struct %v", edgeHostnames.ContractID)
	err = edgeHostnames.GetEdgeHostnames(property.Contract, property.Group, "", CorrelationID)
	if err != nil {
		return diag.FromErr(err)
	}
	logger.Debugf("EdgeHostnames exist in contract")

	if len(edgeHostnames.EdgeHostnames.Items) == 0 {
		return diag.FromErr(fmt.Errorf("no default edge hostname found"))
	}
	logger.Debugf("Edgehostnames Default host %v", edgeHostnames.EdgeHostnames.Items[0])
	defaultEdgeHostname := edgeHostnames.EdgeHostnames.Items[0]

	var found bool
	var edgeHostnameID string
	edgeHostname, err := tools.GetStringValue("edge_hostname", d)
	if err != nil && !errors.Is(err, tools.ErrNotFound) {
		return diag.FromErr(err)
	}
	if edgeHostname != "" {
		for _, hostname := range edgeHostnames.EdgeHostnames.Items {
			if hostname.EdgeHostnameDomain == edgeHostname {
				found = true
				defaultEdgeHostname = hostname
				edgeHostnameID = hostname.EdgeHostnameID
			}
		}
		logger.Debugf("Found EdgeHostname %v", found)
		logger.Debugf("Default EdgeHostname %v", defaultEdgeHostname)
	}

	if err := d.Set("contract", contract.ContractID); err != nil {
		return diag.FromErr(fmt.Errorf("%w: %s", tools.ErrValueSet, err.Error()))
	}
	if err := d.Set("group", group.GroupID); err != nil {
		return diag.FromErr(fmt.Errorf("%w: %s", tools.ErrValueSet, err.Error()))
	}
	if err := d.Set("edge_hostname", defaultEdgeHostname.EdgeHostnameDomain); err != nil {
		return diag.FromErr(fmt.Errorf("%w: %s", tools.ErrValueSet, err.Error()))
	}
	d.SetId(edgeHostnameID)
	return diags
}
