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
		Deprecated: akamai.NoticeDeprecatedUseAlias("product"),
		StateFunc:  addPrefixToState("prd_"),
	},
	"product_id": {
		Type:         schema.TypeString,
		Optional:     true,
		Computed:     true,
		ExactlyOneOf: []string{"product", "product_id"},
		StateFunc:    addPrefixToState("prd_"),
	},
	"contract": {
		Type:       schema.TypeString,
		Optional:   true,
		Computed:   true,
		Deprecated: akamai.NoticeDeprecatedUseAlias("contract"),
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
		Deprecated: akamai.NoticeDeprecatedUseAlias("group"),
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
		ValidateDiagFunc: tools.IsNotBlank,
	},
	"ip_behavior": {
		Type:     schema.TypeString,
		Required: true,
		ForceNew: true,
		ValidateFunc: func(val interface{}, key string) (warns []string, errs []error) {
			v := val.(string)
			if !strings.EqualFold(papi.EHIPVersionV4, v) && !strings.EqualFold(papi.EHIPVersionV6Performance, v) && !strings.EqualFold(papi.EHIPVersionV6Compliance, v) {
				errs = append(errs, fmt.Errorf("%v must be one of %v, %v, %v, got: %v", key, papi.EHIPVersionV4, papi.EHIPVersionV6Performance, papi.EHIPVersionV6Compliance, v))
			}
			return
		},
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
	groupID = tools.AddPrefix(groupID, "grp_")
	// set group/groupID into ResourceData
	if err := d.Set("group_id", groupID); err != nil {
		return diag.FromErr(fmt.Errorf("%w: %s", tools.ErrValueSet, err.Error()))
	}
	if err := d.Set("group", groupID); err != nil {
		return diag.FromErr(fmt.Errorf("%w: %s", tools.ErrValueSet, err.Error()))
	}

	// Schema guarantees contract_id/contract are strings and one or the other is set
	var contractID string
	if got, ok := d.GetOk("contract_id"); ok {
		contractID = got.(string)
	} else {
		contractID = d.Get("contract").(string)
	}
	contractID = tools.AddPrefix(contractID, "ctr_")
	// set contract/contract_id into ResourceData
	if err := d.Set("contract_id", contractID); err != nil {
		return diag.FromErr(fmt.Errorf("%w: %s", tools.ErrValueSet, err.Error()))
	}
	if err := d.Set("contract", contractID); err != nil {
		return diag.FromErr(fmt.Errorf("%w: %s", tools.ErrValueSet, err.Error()))
	}

	logger.Debugf("Edgehostnames GROUP = %v", groupID)
	logger.Debugf("Edgehostnames CONTRACT = %v", contractID)

	// Schema guarantees product_id/product are strings and one or the other is set
	productID, err := tools.ResolveKeyStringState(d, "product_id", "product")
	if err != nil {
		return diag.FromErr(fmt.Errorf("%v: %s, %s", tools.ErrNotFound, "product_id", "product"))
	}
	productID = tools.AddPrefix(productID, "prd_")
	// set product/product_id into ResourceData
	if err := d.Set("product_id", productID); err != nil {
		return diag.FromErr(fmt.Errorf("%w: %s", tools.ErrValueSet, err.Error()))
	}
	if err := d.Set("product", productID); err != nil {
		return diag.FromErr(fmt.Errorf("%w: %s", tools.ErrValueSet, err.Error()))
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

	// ip_behavior is required value in schema.
	newHostname.IPVersionBehavior = strings.ToUpper(d.Get("ip_behavior").(string))
	var ehnID string
	for _, h := range edgeHostnames.EdgeHostnames.Items {
		if h.DomainPrefix == newHostname.DomainPrefix && h.DomainSuffix == newHostname.DomainSuffix {
			ehnID = h.ID
		}
	}
	certEnrollmentID, err := tools.GetIntValue("certificate", d)
	if err != nil {
		if !errors.Is(err, tools.ErrNotFound) {
			return diag.FromErr(err)
		}
		if newHostname.SecureNetwork == papi.EHSecureNetworkEnhancedTLS {
			return diag.FromErr(fmt.Errorf("a certificate enrollment ID is required for Enhanced TLS edge hostnames with 'edgekey.net' suffix"))
		}
	}

	newHostname.CertEnrollmentID = certEnrollmentID
	newHostname.SlotNumber = certEnrollmentID
	if ehnID == "" {
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
		ehnID = hostname.EdgeHostnameID
	} else {
		d.SetId(ehnID)
	}
	logger.Debugf("Resulting EHN Id: %s ", ehnID)
	return resourceSecureEdgeHostNameRead(ctx, d, meta)
}

func resourceSecureEdgeHostNameDelete(_ context.Context, d *schema.ResourceData, m interface{}) diag.Diagnostics {
	meta := akamai.Meta(m)
	logger := meta.Log("PAPI", "resourceSecureEdgeHostNameDelete")
	logger.Debug("DELETING")
	logger.Info("PAPI does not support edge hostname deletion - resource will only be removed from state")
	d.SetId("")
	logger.Debugf("DONE")
	return nil
}

func resourceSecureEdgeHostNameImport(ctx context.Context, d *schema.ResourceData, m interface{}) ([]*schema.ResourceData, error) {
	meta := akamai.Meta(m)
	client := inst.Client(meta)

	parts := strings.Split(d.Id(), ",")
	if len(parts) < 3 {
		return nil, fmt.Errorf("comma-separated list of EdgehostNameID, contractID and groupID has to be supplied in import: %s", d.Id())
	}

	edgehostID := parts[0]
	contractID := tools.AddPrefix(parts[1], "ctr_")
	groupID := tools.AddPrefix(parts[2], "grp_")

	edgehostnameDetails, err := client.GetEdgeHostname(ctx, papi.GetEdgeHostnameRequest{
		EdgeHostnameID: edgehostID,
		ContractID:     contractID,
		GroupID:        groupID,
	})
	if err != nil {
		return nil, err
	}

	if err := d.Set("contract", edgehostnameDetails.ContractID); err != nil {
		return nil, fmt.Errorf("%w: %s", tools.ErrValueSet, err.Error())
	}
	if err := d.Set("contract_id", edgehostnameDetails.ContractID); err != nil {
		return nil, fmt.Errorf("%w: %s", tools.ErrValueSet, err.Error())
	}
	if err := d.Set("group", edgehostnameDetails.GroupID); err != nil {
		return nil, fmt.Errorf("%w: %s", tools.ErrValueSet, err.Error())
	}
	if err := d.Set("group_id", edgehostnameDetails.GroupID); err != nil {
		return nil, fmt.Errorf("%w: %s", tools.ErrValueSet, err.Error())
	}
	productID := edgehostnameDetails.EdgeHostname.ProductID
	if err := d.Set("product", productID); err != nil {
		return nil, fmt.Errorf("%w: %s", tools.ErrValueSet, err.Error())
	}
	if err := d.Set("product_id", productID); err != nil {
		return nil, fmt.Errorf("%w: %s", tools.ErrValueSet, err.Error())
	}
	if err := d.Set("edge_hostname", edgehostnameDetails.EdgeHostname.Domain); err != nil {
		return nil, fmt.Errorf("%w: %s", tools.ErrValueSet, err.Error())
	}
	d.SetId(edgehostID)

	return []*schema.ResourceData{d}, nil
}

func resourceSecureEdgeHostNameRead(ctx context.Context, d *schema.ResourceData, m interface{}) diag.Diagnostics {
	meta := akamai.Meta(m)
	logger := meta.Log("PAPI", "resourceSecureEdgeHostNameRead")

	client := inst.Client(meta)

	// Schema guarantees group_id/group are strings and one or the other is set
	var groupID string
	if got, ok := d.GetOk("group_id"); ok {
		groupID = got.(string)
	} else {
		groupID = d.Get("group").(string)
	}
	groupID = tools.AddPrefix(groupID, "grp_")
	// set group/groupID into ResourceData
	if err := d.Set("group_id", groupID); err != nil {
		return diag.FromErr(fmt.Errorf("%w: %s", tools.ErrValueSet, err.Error()))
	}
	if err := d.Set("group", groupID); err != nil {
		return diag.FromErr(fmt.Errorf("%w: %s", tools.ErrValueSet, err.Error()))
	}

	// Schema guarantees contract_id/contract are strings and one or the other is set
	var contractID string
	if got, ok := d.GetOk("contract_id"); ok {
		contractID = got.(string)
	} else {
		contractID = d.Get("contract").(string)
	}
	contractID = tools.AddPrefix(contractID, "ctr_")
	// set contract/contract_id into ResourceData
	if err := d.Set("contract_id", contractID); err != nil {
		return diag.FromErr(fmt.Errorf("%w: %s", tools.ErrValueSet, err.Error()))
	}
	if err := d.Set("contract", contractID); err != nil {
		return diag.FromErr(fmt.Errorf("%w: %s", tools.ErrValueSet, err.Error()))
	}

	// Schema guarantees product_id/product are strings and one or the other is set
	var productID string
	if got, ok := d.GetOk("product_id"); ok {
		productID = got.(string)
	} else {
		productID = d.Get("product").(string)
	}
	productID = tools.AddPrefix(productID, "prd_")
	// set product/product_id into ResourceData
	if err := d.Set("product_id", productID); err != nil {
		return diag.FromErr(fmt.Errorf("%w: %s", tools.ErrValueSet, err.Error()))
	}
	if err := d.Set("product", productID); err != nil {
		return diag.FromErr(fmt.Errorf("%w: %s", tools.ErrValueSet, err.Error()))
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
