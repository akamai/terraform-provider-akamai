package property

import (
	"context"
	"encoding/json"
	"errors"
	"fmt"
	"strings"

	"github.com/hashicorp/terraform-plugin-sdk/v2/diag"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/schema"

	"github.com/akamai/AkamaiOPEN-edgegrid-golang/v4/pkg/hapi"
	"github.com/akamai/AkamaiOPEN-edgegrid-golang/v4/pkg/papi"
	"github.com/akamai/terraform-provider-akamai/v3/pkg/akamai"
	"github.com/akamai/terraform-provider-akamai/v3/pkg/tools"
)

func resourceSecureEdgeHostName() *schema.Resource {
	return &schema.Resource{
		CreateContext: resourceSecureEdgeHostNameCreate,
		ReadContext:   resourceSecureEdgeHostNameRead,
		UpdateContext: resourceSecureEdgeHostNameUpdate,
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
		DiffSuppressFunc: diffSuppressEdgeHostname,
		ValidateDiagFunc: tools.IsNotBlank,
		StateFunc:        appendDefaultSuffixToEdgeHostname,
	},
	"ip_behavior": {
		Type:     schema.TypeString,
		Required: true,
	},
	"status_update_email": {
		Type:        schema.TypeList,
		Optional:    true,
		Elem:        &schema.Schema{Type: schema.TypeString},
		Description: "Email address that should receive updates on the IP behavior update request. Required for update operation.",
	},
	"certificate": {
		Type:     schema.TypeInt,
		Optional: true,
		ForceNew: true,
	},
	"use_cases": {
		Type:             schema.TypeString,
		Optional:         true,
		ForceNew:         true,
		DiffSuppressFunc: suppressEdgeHostnameUseCases,
		Description:      "A JSON encoded list of use cases",
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

	newHostname.DomainSuffix, newHostname.SecureNetwork = parseEdgeHostname(edgeHostname)

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

	useCasesJSON, err := tools.GetStringValue("use_cases", d)
	if err != nil && !errors.Is(err, tools.ErrNotFound) {
		return diag.FromErr(err)
	}
	if useCasesJSON != "" {
		var useCases []papi.UseCase
		if err := json.Unmarshal([]byte(useCasesJSON), &useCases); err != nil {
			return diag.Errorf("error while un-marshaling use cases JSON: %s", err)
		}
		newHostname.UseCases = useCases
	}

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

	var edgeHostname string
	if got, ok := d.GetOk("edge_hostname"); ok {
		edgeHostname = got.(string)
	}

	foundEdgeHostname, err := findEdgeHostname(edgeHostnames.EdgeHostnames, edgeHostname)
	if err != nil {
		return diag.FromErr(err)
	}

	useCasesJSON, err := useCases2JSON(foundEdgeHostname.UseCases)
	if err != nil {
		return diag.FromErr(err)
	}
	if err := d.Set("use_cases", string(useCasesJSON)); err != nil {
		return diag.FromErr(fmt.Errorf("%w: %s", tools.ErrValueSet, err.Error()))
	}

	if err := d.Set("edge_hostname", foundEdgeHostname.Domain); err != nil {
		return diag.FromErr(fmt.Errorf("%w: %s", tools.ErrValueSet, err.Error()))
	}
	d.SetId(foundEdgeHostname.ID)

	return nil
}

func resourceSecureEdgeHostNameUpdate(ctx context.Context, d *schema.ResourceData, m interface{}) diag.Diagnostics {
	meta := akamai.Meta(m)
	logger := meta.Log("PAPI", "resourceSecureEdgeHostNameUpdate")

	if d.HasChange("ip_behavior") {
		edgeHostname, err := tools.GetStringValue("edge_hostname", d)
		if err != nil {
			return diag.FromErr(err)
		}
		dnsZone, _ := parseEdgeHostname(edgeHostname)
		ipBehavior, err := tools.GetStringValue("ip_behavior", d)
		if err != nil {
			return diag.FromErr(err)
		}
		emails, err := tools.GetListValue("status_update_email", d)
		if err != nil && !errors.Is(err, tools.ErrNotFound) {
			return diag.FromErr(err)
		}
		if len(emails) == 0 {
			return diag.Errorf(`"status_update_email" is a required parameter to update an edge hostname`)
		}
		statusUpdateEmails := make([]string, len(emails))
		for i, email := range emails {
			statusUpdateEmails[i] = email.(string)
		}

		logger.Debugf("Proceeding to update /ipVersionBehavior for %s", edgeHostname)
		// IPV6_COMPLIANCE type has to mapped to IPV6_IPV4_DUALSTACK which is only accepted value by HAPI client
		if ipBehavior == papi.EHIPVersionV6Compliance {
			ipBehavior = "IPV6_IPV4_DUALSTACK"
		}

		if _, err = inst.HapiClient(meta).UpdateEdgeHostname(ctx, hapi.UpdateEdgeHostnameRequest{
			DNSZone:           dnsZone,
			RecordName:        strings.ReplaceAll(edgeHostname, "."+dnsZone, ""),
			Comments:          fmt.Sprintf("change /ipVersionBehavior to %s", ipBehavior),
			StatusUpdateEmail: statusUpdateEmails,
			Body: []hapi.UpdateEdgeHostnameRequestBody{
				{
					Op:    "replace",
					Path:  "/ipVersionBehavior",
					Value: ipBehavior,
				},
			},
		}); err != nil {
			if err2 := tools.RestoreOldValues(d, []string{"ip_behavior"}); err2 != nil {
				return diag.Errorf(`%s failed. No changes were written to server:
%s

Failed to restore previous local schema values. The schema will remain in tainted state:
%s`, hapi.ErrUpdateEdgeHostname, err.Error(), err2.Error())
			}
			return diag.FromErr(err)
		}
	}

	return resourceSecureEdgeHostNameRead(ctx, d, m)
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
	useCasesJSON, err := useCases2JSON(edgehostnameDetails.EdgeHostname.UseCases)
	if err != nil {
		return nil, err
	}
	if err := d.Set("use_cases", string(useCasesJSON)); err != nil {
		return nil, fmt.Errorf("%w: %s", tools.ErrValueSet, err.Error())
	}
	if err := d.Set("edge_hostname", edgehostnameDetails.EdgeHostname.Domain); err != nil {
		return nil, fmt.Errorf("%w: %s", tools.ErrValueSet, err.Error())
	}
	if err := d.Set("ip_behavior", edgehostnameDetails.EdgeHostname.IPVersionBehavior); err != nil {
		return nil, fmt.Errorf("%w: %s", tools.ErrValueSet, err.Error())
	}
	d.SetId(edgehostID)

	return []*schema.ResourceData{d}, nil
}

func diffSuppressEdgeHostname(_, oldVal, newVal string, _ *schema.ResourceData) bool {
	oldVal = strings.ToLower(oldVal)
	newVal = strings.ToLower(newVal)

	if oldVal == newVal {
		return true
	}

	if !(strings.HasSuffix(newVal, "edgekey.net") || strings.HasSuffix(newVal, "edgesuite.net") ||
		strings.HasSuffix(newVal, "akamaized.net")) {
		return oldVal == fmt.Sprintf("%s.edgesuite.net", newVal)
	}
	return false
}

func suppressEdgeHostnameUseCases(_, oldVal, newVal string, _ *schema.ResourceData) bool {
	logger := akamai.Log("PAPI", "suppressEdgeHostnameUseCases")
	if oldVal == newVal {
		return true
	}
	var oldUseCases, newUseCases []papi.UseCase
	if err := json.Unmarshal([]byte(oldVal), &oldUseCases); err != nil {
		logger.Errorf("Unable to unmarshal 'old' use cases: %s", err)
		return false
	}
	if err := json.Unmarshal([]byte(newVal), &newUseCases); err != nil {
		logger.Errorf("Unable to unmarshal 'new' use cases: %s", err)
		return false
	}

	diff := make(map[papi.UseCase]int)
	for _, useCase := range oldUseCases {
		diff[useCase]++
	}

	for _, useCase := range newUseCases {
		if _, ok := diff[useCase]; !ok {
			return false
		}

		diff[useCase]--
		if diff[useCase] == 0 {
			delete(diff, useCase)
		}
	}

	return len(diff) == 0
}

// appendDefaultSuffixToEdgeHostname is a StateFunc which appends ".edgesuite.net" to an edge hostname if none of the supported prefixes were provided
// It is used in order to retain idempotency when "edge_hostname" value is used as output
func appendDefaultSuffixToEdgeHostname(i interface{}) string {
	name := strings.ToLower(i.(string))
	if !(strings.HasSuffix(name, "edgekey.net") || strings.HasSuffix(name, "edgesuite.net") || strings.HasSuffix(name, "akamaized.net")) {
		name = fmt.Sprintf("%s.edgesuite.net", name)
	}
	return name
}

func findEdgeHostname(edgeHostnames papi.EdgeHostnameItems, domain string) (*papi.EdgeHostnameGetItem, error) {
	suffix := "edgesuite.net"
	domain = strings.ToLower(domain)
	if domain != "" {
		if strings.HasSuffix(domain, "edgekey.net") {
			suffix = "edgekey.net"
		}
		if strings.HasSuffix(domain, "akamaized.net") {
			suffix = "akamaized.net"
		}
		prefix := strings.TrimSuffix(domain, "."+suffix)

		for _, eHn := range edgeHostnames.Items {
			if eHn.DomainPrefix == prefix && eHn.DomainSuffix == suffix {
				return &eHn, nil
			}
		}
	}

	return nil, fmt.Errorf("%w: %s", ErrEdgeHostnameNotFound, domain)
}

func useCases2JSON(useCases []papi.UseCase) ([]byte, error) {
	if len(useCases) == 0 {
		return []byte{}, nil
	}
	return json.MarshalIndent(useCases, "", "  ")
}

func parseEdgeHostname(hostname string) (string, string) {
	switch {
	case strings.HasSuffix(hostname, ".edgesuite.net"):
		return "edgesuite.net", papi.EHSecureNetworkStandardTLS
	case strings.HasSuffix(hostname, ".edgekey.net"):
		return "edgekey.net", papi.EHSecureNetworkEnhancedTLS
	case strings.HasSuffix(hostname, ".akamaized.net"):
		return "akamaized.net", papi.EHSecureNetworkSharedCert
	}
	return "edgesuite.net", ""
}
