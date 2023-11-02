package property

import (
	"context"
	"encoding/json"
	"errors"
	"fmt"
	"strings"
	"time"

	"github.com/akamai/AkamaiOPEN-edgegrid-golang/v7/pkg/hapi"
	"github.com/akamai/AkamaiOPEN-edgegrid-golang/v7/pkg/papi"

	"github.com/akamai/terraform-provider-akamai/v5/pkg/common/tf"
	"github.com/akamai/terraform-provider-akamai/v5/pkg/common/timeouts"
	"github.com/akamai/terraform-provider-akamai/v5/pkg/logger"
	"github.com/akamai/terraform-provider-akamai/v5/pkg/meta"
	"github.com/akamai/terraform-provider-akamai/v5/pkg/tools"

	"github.com/hashicorp/terraform-plugin-sdk/v2/diag"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/schema"
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
		Timeouts: &schema.ResourceTimeout{
			Default: &timeouts.SDKDefaultTimeout,
		},
	}
}

var akamaiSecureEdgeHostNameSchema = map[string]*schema.Schema{
	"product_id": {
		Type:      schema.TypeString,
		Optional:  true,
		Computed:  true,
		StateFunc: addPrefixToState("prd_"),
	},
	"contract_id": {
		Type:      schema.TypeString,
		Required:  true,
		StateFunc: addPrefixToState("ctr_"),
	},
	"group_id": {
		Type:      schema.TypeString,
		Required:  true,
		StateFunc: addPrefixToState("grp_"),
	},
	"edge_hostname": {
		Type:             schema.TypeString,
		Required:         true,
		ForceNew:         true,
		DiffSuppressFunc: diffSuppressEdgeHostname,
		ValidateDiagFunc: tf.IsNotBlank,
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
		Description: "Email address that should receive updates on the IP behavior update request.",
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
	"timeouts": {
		Type:        schema.TypeList,
		Optional:    true,
		MaxItems:    1,
		Description: "Enables to set timeout for processing",
		Elem: &schema.Resource{
			Schema: map[string]*schema.Schema{
				"default": {
					Type:             schema.TypeString,
					Optional:         true,
					ValidateDiagFunc: timeouts.ValidateDurationFormat,
				},
			},
		},
	},
}

func resourceSecureEdgeHostNameCreate(ctx context.Context, d *schema.ResourceData, m interface{}) diag.Diagnostics {
	meta := meta.Must(m)
	logger := meta.Log("PAPI", "resourceSecureEdgeHostNameCreate")

	client := Client(meta)

	groupID, err := tf.GetStringValue("group_id", d)
	if err != nil {
		return diag.FromErr(err)
	}
	groupID = tools.AddPrefix(groupID, "grp_")
	if err := d.Set("group_id", groupID); err != nil {
		return diag.FromErr(fmt.Errorf("%w: %s", tf.ErrValueSet, err.Error()))
	}

	contractID, err := tf.GetStringValue("contract_id", d)
	if err != nil {
		return diag.FromErr(err)
	}
	contractID = tools.AddPrefix(contractID, "ctr_")
	if err := d.Set("contract_id", contractID); err != nil {
		return diag.FromErr(fmt.Errorf("%w: %s", tf.ErrValueSet, err.Error()))
	}

	logger.Debugf("Edgehostnames GROUP = %v", groupID)
	logger.Debugf("Edgehostnames CONTRACT = %v", contractID)

	// Schema no longer guarantees that product_id is set, this field is required only for creation
	productID, err := tf.GetStringValue("product_id", d)
	if err != nil {
		return diag.Errorf("`product_id` must be specified for creation")
	}
	productID = tools.AddPrefix(productID, "prd_")
	if err := d.Set("product_id", productID); err != nil {
		return diag.FromErr(fmt.Errorf("%w: %s", tf.ErrValueSet, err.Error()))
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
	certEnrollmentID, err := tf.GetIntValue("certificate", d)
	if err != nil {
		if !errors.Is(err, tf.ErrNotFound) {
			return diag.FromErr(err)
		}
		if newHostname.SecureNetwork == papi.EHSecureNetworkEnhancedTLS {
			return diag.FromErr(fmt.Errorf("a certificate enrollment ID is required for Enhanced TLS edge hostnames with 'edgekey.net' suffix"))
		}
	}
	newHostname.CertEnrollmentID = certEnrollmentID
	newHostname.SlotNumber = certEnrollmentID

	useCasesJSON, err := tf.GetStringValue("use_cases", d)
	if err != nil && !errors.Is(err, tf.ErrNotFound) {
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
	meta := meta.Must(m)
	logger := meta.Log("PAPI", "resourceSecureEdgeHostNameRead")

	client := Client(meta)

	groupID, err := tf.GetStringValue("group_id", d)
	if err != nil {
		return diag.FromErr(err)
	}
	groupID = tools.AddPrefix(groupID, "grp_")
	if err := d.Set("group_id", groupID); err != nil {
		return diag.FromErr(fmt.Errorf("%w: %s", tf.ErrValueSet, err.Error()))
	}

	contractID, err := tf.GetStringValue("contract_id", d)
	if err != nil {
		return diag.FromErr(err)
	}
	contractID = tools.AddPrefix(contractID, "ctr_")
	// set contract/contract_id into ResourceData
	if err := d.Set("contract_id", contractID); err != nil {
		return diag.FromErr(fmt.Errorf("%w: %s", tf.ErrValueSet, err.Error()))
	}

	// Schema guarantees product_id/product are strings and one or the other is set
	var productID string
	if got, ok := d.GetOk("product_id"); ok {
		productID = got.(string)
	}
	productID = tools.AddPrefix(productID, "prd_")
	if err := d.Set("product_id", productID); err != nil {
		return diag.FromErr(fmt.Errorf("%w: %s", tf.ErrValueSet, err.Error()))
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
		return diag.FromErr(fmt.Errorf("%w: %s", tf.ErrValueSet, err.Error()))
	}

	if err := d.Set("edge_hostname", foundEdgeHostname.Domain); err != nil {
		return diag.FromErr(fmt.Errorf("%w: %s", tf.ErrValueSet, err.Error()))
	}

	if err := d.Set("ip_behavior", foundEdgeHostname.IPVersionBehavior); err != nil {
		return diag.FromErr(fmt.Errorf("%w: %s", tf.ErrValueSet, err.Error()))
	}

	d.SetId(foundEdgeHostname.ID)

	return nil
}

func resourceSecureEdgeHostNameUpdate(ctx context.Context, d *schema.ResourceData, m interface{}) diag.Diagnostics {
	meta := meta.Must(m)
	logger := meta.Log("PAPI", "resourceSecureEdgeHostNameUpdate")

	if !d.HasChangeExcept("timeouts") {
		logger.Debug("Only timeouts were updated, skipping")
		return nil
	}

	if d.HasChange("ip_behavior") {
		edgeHostname, err := tf.GetStringValue("edge_hostname", d)
		if err != nil {
			return diag.FromErr(err)
		}
		dnsZone, _ := parseEdgeHostname(edgeHostname)
		ipBehavior, err := tf.GetStringValue("ip_behavior", d)
		if err != nil {
			return diag.FromErr(err)
		}
		emails, err := tf.GetListValue("status_update_email", d)
		if err != nil && !errors.Is(err, tf.ErrNotFound) {
			return diag.FromErr(err)
		}

		logger.Debugf("Proceeding to update /ipVersionBehavior for %s", edgeHostname)
		// IPV6_COMPLIANCE type has to mapped to IPV6_IPV4_DUALSTACK which is only accepted value by HAPI client
		if ipBehavior == papi.EHIPVersionV6Compliance {
			ipBehavior = "IPV6_IPV4_DUALSTACK"
		}

		req := hapi.UpdateEdgeHostnameRequest{
			DNSZone:    dnsZone,
			RecordName: strings.ReplaceAll(edgeHostname, "."+dnsZone, ""),
			Comments:   fmt.Sprintf("change /ipVersionBehavior to %s", ipBehavior),
			Body: []hapi.UpdateEdgeHostnameRequestBody{
				{
					Op:    "replace",
					Path:  "/ipVersionBehavior",
					Value: ipBehavior,
				},
			},
		}
		if len(emails) != 0 {
			statusUpdateEmails := make([]string, len(emails))
			for i, email := range emails {
				statusUpdateEmails[i] = email.(string)
			}
			req.StatusUpdateEmail = statusUpdateEmails
		}

		hapiClient := HapiClient(meta)
		resp, err := hapiClient.UpdateEdgeHostname(ctx, req)
		if err != nil {
			if err2 := tf.RestoreOldValues(d, []string{"ip_behavior"}); err2 != nil {
				return diag.Errorf(`%s failed. No changes were written to server:
%s

Failed to restore previous local schema values. The schema will remain in tainted state:
%s`, hapi.ErrUpdateEdgeHostname, err.Error(), err2.Error())
			}
			return diag.FromErr(err)
		}

		if err = waitForChange(ctx, hapiClient, resp.ChangeID); err != nil {
			return diag.FromErr(err)
		}
	}

	return resourceSecureEdgeHostNameRead(ctx, d, m)
}

func waitForChange(ctx context.Context, client hapi.HAPI, changeID int) error {
	for {
		change, err := client.GetChangeRequest(ctx, hapi.GetChangeRequest{
			ChangeID: changeID,
		})
		if err != nil {
			return err
		}
		if change.Status == "PENDING" {
			select {
			case <-time.After(time.Second * 10):
			case <-ctx.Done():
				return ctx.Err()
			}
			continue
		}
		if change.Status == "SUCCEEDED" {
			return nil
		}
		return fmt.Errorf("unexpected change status: %s", change.Status)
	}
}

func resourceSecureEdgeHostNameDelete(ctx context.Context, d *schema.ResourceData, m interface{}) diag.Diagnostics {
	meta := meta.Must(m)
	logger := meta.Log("HAPI", "resourceSecureEdgeHostNameDelete")

	edgeHostname, err := tf.GetStringValue("edge_hostname", d)
	dnsZone, _ := parseEdgeHostname(edgeHostname)
	recordName := strings.ReplaceAll(edgeHostname, "."+dnsZone, "")

	logger.Debugf("edge_hostname = %v", edgeHostname)
	logger.Debugf("dnsZone = %v, recordName = %v", dnsZone, recordName)
	req := hapi.DeleteEdgeHostnameRequest{
		DNSZone:    dnsZone,
		RecordName: recordName,
		Comments:   fmt.Sprintf("remove %s", edgeHostname),
	}

	emails, err := tf.GetListValue("status_update_email", d)
	if err != nil && !errors.Is(err, tf.ErrNotFound) {
		return diag.FromErr(err)
	}

	if len(emails) != 0 {
		statusUpdateEmails := make([]string, len(emails))
		for i, email := range emails {
			statusUpdateEmails[i] = email.(string)
		}
		req.StatusUpdateEmail = statusUpdateEmails
	}

	hapiClient := HapiClient(meta)
	logger.Debugf("hapiClient.DeleteEdgeHostname: req = %v", req)
	resp, err := hapiClient.DeleteEdgeHostname(ctx, req)
	if err != nil {
		return diag.FromErr(err)
	}

	if err = waitForChange(ctx, hapiClient, resp.ChangeID); err != nil {
		return diag.FromErr(err)
	}
	return nil
}

func resourceSecureEdgeHostNameImport(ctx context.Context, d *schema.ResourceData, m interface{}) ([]*schema.ResourceData, error) {
	meta := meta.Must(m)
	client := Client(meta)

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

	if err := d.Set("contract_id", edgehostnameDetails.ContractID); err != nil {
		return nil, fmt.Errorf("%w: %s", tf.ErrValueSet, err.Error())
	}
	if err := d.Set("group_id", edgehostnameDetails.GroupID); err != nil {
		return nil, fmt.Errorf("%w: %s", tf.ErrValueSet, err.Error())
	}
	useCasesJSON, err := useCases2JSON(edgehostnameDetails.EdgeHostname.UseCases)
	if err != nil {
		return nil, err
	}
	if err := d.Set("use_cases", string(useCasesJSON)); err != nil {
		return nil, fmt.Errorf("%w: %s", tf.ErrValueSet, err.Error())
	}
	if err := d.Set("edge_hostname", edgehostnameDetails.EdgeHostname.Domain); err != nil {
		return nil, fmt.Errorf("%w: %s", tf.ErrValueSet, err.Error())
	}
	if err := d.Set("ip_behavior", edgehostnameDetails.EdgeHostname.IPVersionBehavior); err != nil {
		return nil, fmt.Errorf("%w: %s", tf.ErrValueSet, err.Error())
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
	logger := logger.Get("PAPI", "suppressEdgeHostnameUseCases")
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
