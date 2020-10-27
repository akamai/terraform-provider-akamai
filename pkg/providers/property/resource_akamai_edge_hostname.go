package property

import (
	"context"
	"errors"
	"fmt"
	"strings"

	"github.com/hashicorp/terraform-plugin-sdk/v2/diag"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/schema"

	"github.com/akamai/AkamaiOPEN-edgegrid-golang/v2/pkg/papi"
	"github.com/akamai/terraform-provider-akamai/v2/pkg/akamai"
	"github.com/akamai/terraform-provider-akamai/v2/pkg/tools"
)

func resourceSecureEdgeHostName() *schema.Resource {
	return &schema.Resource{
		CreateContext: resourceSecureEdgeHostNameCreate,
		ReadContext:   resourceSecureEdgeHostNameRead,
		DeleteContext: resourceSecureEdgeHostNameDelete,
		Importer: &schema.ResourceImporter{
			StateContext: resourceSecureEdgeHostNameImport,
		},
		Schema: akamaiSecureEdgeHostNameSchema,
	}
}

var akamaiSecureEdgeHostNameSchema = map[string]*schema.Schema{
	"product": {
		Type:       schema.TypeString,
		Optional:   true,
		Computed:   true,
		Deprecated: `use "product_id" attribute instead`,
		StateFunc:  addPrefixToState("prd_"),
	},
	"product_id": {
		Type:         schema.TypeString,
		Optional:     true,
		Computed:     true,
		ExactlyOneOf: []string{"product_id", "product"},
		StateFunc:    addPrefixToState("prd_"),
	},
	"contract": {
		Type:       schema.TypeString,
		Optional:   true,
		Computed:   true,
		Deprecated: `use "contract_id" attribute instead`,
		StateFunc:  addPrefixToState("ctr_"),
	},
	"contract_id": {
		Type:         schema.TypeString,
		Optional:     true,
		Computed:     true,
		ExactlyOneOf: []string{"contract_id", "contract"},
		StateFunc:    addPrefixToState("ctr_"),
	},
	"group": {
		Type:       schema.TypeString,
		Optional:   true,
		Computed:   true,
		Deprecated: `use "group_id" attribute instead`,
		StateFunc:  addPrefixToState("grp_"),
	},
	"group_id": {
		Type:         schema.TypeString,
		Optional:     true,
		Computed:     true,
		ExactlyOneOf: []string{"group_id", "group"},
		StateFunc:    addPrefixToState("grp_"),
	},
	"edge_hostname": {
		Type:             schema.TypeString,
		Required:         true,
		ForceNew:         true,
		DiffSuppressFunc: suppressEdgeHostnameDomain,
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

func resourceSecureEdgeHostNameCreate(ctx context.Context, d *schema.ResourceData, m interface{}) diag.Diagnostics {
	meta := akamai.Meta(m)
	logger := meta.Log("PAPI", "resourceSecureEdgeHostNameCreate")

	client := inst.Client(meta)

	// Schema guarantees group_id/group are strings and one or the other is set
	var groupID string
	if got, ok := d.GetOk("group_id"); ok {
		groupID = got.(string)
	} else {
		groupID = d.Get("group").(string)
	}
	if !strings.HasPrefix(groupID, "grp_") {
		groupID = fmt.Sprintf("grp_%s", groupID)
	}

	// Schema guarantees contract_id/contract are strings and one or the other is set
	var contractID string
	if got, ok := d.GetOk("contract_id"); ok {
		contractID = got.(string)
	} else {
		contractID = d.Get("contract").(string)
	}
	if !strings.HasPrefix(contractID, "ctr_") {
		contractID = fmt.Sprintf("ctr_%s", contractID)
	}

	logger.Debugf("Edgehostnames GROUP = %v", groupID)
	logger.Debugf("Edgehostnames CONTRACT = %v", contractID)

	// Schema guarantees product_id/product are strings and one or the other is set
	var productID string
	if got, ok := d.GetOk("product_id"); ok {
		productID = got.(string)
	} else {
		productID = d.Get("product").(string)
	}
	if !strings.HasPrefix(productID, "prd_") {
		productID = fmt.Sprintf("prd_%s", productID)
	}

	edgeHostnames, err := client.GetEdgeHostnames(ctx, papi.GetEdgeHostnamesRequest{
		ContractID: contractID,
		GroupID:    groupID,
	})
	if err != nil {
		return diag.FromErr(err)
	}

	var edgeHostname string
	if got, ok := d.GetOk("edge_hostname"); ok {
		edgeHostname = got.(string)
	}
	newHostname := papi.EdgeHostnameCreate{}
	newHostname.ProductID = productID
	newHostname.DomainSuffix = "edgesuite.net"

	switch {
	case strings.HasSuffix(edgeHostname, ".edgesuite.net"):
		newHostname.DomainSuffix = "edgesuite.net"
		newHostname.SecureNetwork = papi.EHSecureNetworkStandardTLS
	case strings.HasSuffix(edgeHostname, ".edgekey.net"):
		newHostname.DomainSuffix = "edgekey.net"
		newHostname.SecureNetwork = papi.EHSecureNetworkEnhancedTLS
	case strings.HasSuffix(edgeHostname, ".akamaized.net"):
		newHostname.DomainSuffix = "akamaized.net"
		newHostname.SecureNetwork = papi.EHSecureNetworkSharedCert
	}
	newHostname.DomainPrefix = strings.TrimSuffix(edgeHostname, "."+newHostname.DomainSuffix)

	ipv4, _ := tools.GetBoolValue("ipv4", d)
	if ipv4, _ := tools.GetBoolValue("ipv4", d); ipv4 {
		newHostname.IPVersionBehavior = "IPV4"
	}
	ipv6, _ := tools.GetBoolValue("ipv6", d)
	if ipv6 {
		newHostname.IPVersionBehavior = "IPV6"
	}
	if ipv4 && ipv6 {
		newHostname.IPVersionBehavior = "IPV6_COMPLIANCE"
	}
	if !(ipv4 || ipv6) {
		return diag.FromErr(fmt.Errorf("ipv4 or ipv6 must be specified to create a new Edge Hostname"))
	}

	if err := d.Set("ip_behavior", newHostname.IPVersionBehavior); err != nil {
		return diag.FromErr(fmt.Errorf("%w: %s", tools.ErrValueSet, err.Error()))
	}

	for _, h := range edgeHostnames.EdgeHostnames.Items {
		if h.DomainPrefix == newHostname.DomainPrefix && h.DomainSuffix == newHostname.DomainSuffix {
			d.SetId(h.ID)
			return nil
		}
	}

	certEnrollmentID, err := tools.GetIntValue("certificate", d)
	if err != nil {
		if !errors.Is(err, tools.ErrNotFound) {
			return diag.FromErr(err)
		}
		if newHostname.SecureNetwork == "ENHANCED_TLS" {
			return diag.FromErr(fmt.Errorf("A certificate enrollment ID is required for Enhanced TLS (edgekey.net) edge hostnames"))
		}
	}

	newHostname.CertEnrollmentID = certEnrollmentID
	newHostname.SlotNumber = certEnrollmentID

	logger.Debugf("Creating new edge hostname: %#v", newHostname)
	hostname, err := client.CreateEdgeHostname(ctx, papi.CreateEdgeHostnameRequest{
		EdgeHostname: newHostname,
		ContractID:   contractID,
		GroupID:      groupID,
	})
	if err != nil {
		return diag.FromErr(err)
	}
	d.SetId(hostname.EdgeHostnameID)
	return resourceSecureEdgeHostNameRead(ctx, d, meta)
}

func resourceSecureEdgeHostNameDelete(_ context.Context, d *schema.ResourceData, m interface{}) diag.Diagnostics {
	meta := akamai.Meta(m)
	logger := meta.Log("PAPI", "resourceSecureEdgeHostNameDelete")
	logger.Debugf("DELETING")
	logger.Info("PAPI does not support edge hostname deletion - resource will only be removed from state")
	d.SetId("")
	logger.Debugf("DONE")
	return nil
}

func resourceSecureEdgeHostNameImport(ctx context.Context, d *schema.ResourceData, m interface{}) ([]*schema.ResourceData, error) {
	meta := akamai.Meta(m)
	logger := meta.Log("PAPI", "resourceSecureEdgeHostNameImport")
	resourceID := d.Id()
	propertyID := resourceID

	client := inst.Client(meta)

	if !strings.HasPrefix(resourceID, "prp_") {
		keys := []string{
			papi.SearchKeyPropertyName,
			papi.SearchKeyHostname,
			papi.SearchKeyEdgeHostname,
		}
		for _, searchKey := range keys {
			results, err := client.SearchProperties(ctx, papi.SearchRequest{
				Key:   searchKey,
				Value: resourceID,
			})
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

	prop, err := client.GetProperty(ctx, papi.GetPropertyRequest{
		PropertyID: propertyID,
	})
	if err != nil {
		return nil, err
	}

	if err := d.Set("contract", prop.Property.ContractID); err != nil {
		return nil, fmt.Errorf("%w: %s", tools.ErrValueSet, err.Error())
	}
	if err := d.Set("group", prop.Property.GroupID); err != nil {
		return nil, fmt.Errorf("%w: %s", tools.ErrValueSet, err.Error())
	}
	if err := d.Set("edge_hostname", prop.Property.GroupID); err != nil {
		return nil, fmt.Errorf("%w: %s", tools.ErrValueSet, err.Error())
	}
	d.SetId(prop.Property.PropertyID)

	return []*schema.ResourceData{d}, nil
}

func resourceSecureEdgeHostNameRead(ctx context.Context, d *schema.ResourceData, m interface{}) diag.Diagnostics {
	meta := akamai.Meta(m)
	logger := meta.Log("PAPI", "resourceSecureEdgeHostNameCreate")

	client := inst.Client(meta)

	// Schema guarantees group_id/group are strings and one or the other is set
	var groupID string
	if got, ok := d.GetOk("group_id"); ok {
		groupID = got.(string)
	} else {
		groupID = d.Get("group").(string)
	}
	if !strings.HasPrefix(groupID, "grp_") {
		groupID = fmt.Sprintf("grp_%s", groupID)
	}

	// Schema guarantees contract_id/contract are strings and one or the other is set
	var contractID string
	if got, ok := d.GetOk("contract_id"); ok {
		contractID = got.(string)
	} else {
		contractID = d.Get("contract").(string)
	}
	if !strings.HasPrefix(contractID, "ctr_") {
		contractID = fmt.Sprintf("ctr_%s", contractID)
	}

	logger.Debugf("Edgehostnames GROUP = %v", groupID)
	logger.Debugf("Edgehostnames CONTRACT = %v", contractID)

	edgeHostnames, err := client.GetEdgeHostnames(ctx, papi.GetEdgeHostnamesRequest{
		ContractID: contractID,
		GroupID:    groupID,
	})
	if err != nil {
		return diag.FromErr(err)
	}
	defaultEdgeHostname := &edgeHostnames.EdgeHostnames.Items[0]

	var edgeHostname string
	if got, ok := d.GetOk("edge_hostname"); ok {
		edgeHostname = got.(string)
	}

	if edgeHostname != "" {
		found, err := findEdgeHostname(edgeHostnames.EdgeHostnames, edgeHostname)
		if err != nil && !errors.Is(err, ErrEdgeHostnameNotFound) {
			return diag.FromErr(err)
		}
		if err == nil {
			defaultEdgeHostname = found
		}
	}

	if err := d.Set("contract", contractID); err != nil {
		return diag.FromErr(fmt.Errorf("%w: %s", tools.ErrValueSet, err.Error()))
	}
	if err := d.Set("group", groupID); err != nil {
		return diag.FromErr(fmt.Errorf("%w: %s", tools.ErrValueSet, err.Error()))
	}
	if err := d.Set("edge_hostname", defaultEdgeHostname.Domain); err != nil {
		return diag.FromErr(fmt.Errorf("%w: %s", tools.ErrValueSet, err.Error()))
	}
	d.SetId(defaultEdgeHostname.ID)

	return nil
}

func suppressEdgeHostnameDomain(_, old, new string, _ *schema.ResourceData) bool {
	if old == new {
		return true
	}
	if !(strings.HasSuffix(new, "edgekey.net") || strings.HasSuffix(new, "edgesuite.net") || strings.HasSuffix(new, "akamaized.net")) {
		return old == fmt.Sprintf("%s.edgesuite.net", new)
	}
	return false
}

func findEdgeHostname(edgeHostnames papi.EdgeHostnameItems, domain string) (*papi.EdgeHostnameGetItem, error) {
	var prefix string
	suffix := "edgesuite.net"
	if domain != "" {
		if strings.HasSuffix(domain, "edgekey.net") {
			suffix = "edgekey.net"
		}
		if strings.HasSuffix(domain, "akamaized.net") {
			suffix = "akamaized.net"
		}
		prefix = strings.TrimSuffix(domain, "."+suffix)
	}

	for _, eHn := range edgeHostnames.Items {
		if eHn.DomainPrefix == prefix && eHn.DomainSuffix == suffix {
			return &eHn, nil
		}
	}

	return nil, fmt.Errorf("%w: %s", ErrEdgeHostnameNotFound, domain)
}
