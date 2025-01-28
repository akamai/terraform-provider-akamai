package property

import (
	"context"
	"encoding/json"
	"errors"
	"fmt"
	"strconv"
	"strings"
	"time"

	"github.com/akamai/AkamaiOPEN-edgegrid-golang/v10/pkg/hapi"
	"github.com/akamai/AkamaiOPEN-edgegrid-golang/v10/pkg/papi"
	"github.com/akamai/terraform-provider-akamai/v6/pkg/common/str"
	"github.com/akamai/terraform-provider-akamai/v6/pkg/common/tf"
	"github.com/akamai/terraform-provider-akamai/v6/pkg/common/timeouts"
	"github.com/akamai/terraform-provider-akamai/v6/pkg/log"
	"github.com/akamai/terraform-provider-akamai/v6/pkg/meta"

	"github.com/hashicorp/terraform-plugin-sdk/v2/diag"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/customdiff"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/schema"
)

const retriesMax = 15

var (
	// EgdeHostnameCreatePollInterval is the interval for polling an edgehostname creation
	EgdeHostnameCreatePollInterval = time.Minute
)

func resourceSecureEdgeHostName() *schema.Resource {
	return &schema.Resource{
		CustomizeDiff: customdiff.All(
			validateImmutableFields,
		),
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
	"ttl": {
		Type:        schema.TypeInt,
		Optional:    true,
		Description: "The time to live, or number of seconds to keep an edge hostname assigned to a map or target. If not provided default value for product is used.",
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
	groupID = str.AddPrefix(groupID, "grp_")
	if err := d.Set("group_id", groupID); err != nil {
		return diag.FromErr(fmt.Errorf("%w: %s", tf.ErrValueSet, err.Error()))
	}

	contractID, err := tf.GetStringValue("contract_id", d)
	if err != nil {
		return diag.FromErr(err)
	}
	contractID = str.AddPrefix(contractID, "ctr_")
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
	productID = str.AddPrefix(productID, "prd_")
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

	for _, h := range edgeHostnames.EdgeHostnames.Items {
		if h.DomainPrefix == newHostname.DomainPrefix && h.DomainSuffix == newHostname.DomainSuffix {
			return diag.Errorf("edgehostname '%s' already exists", edgeHostname)
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

	logger.Debugf("Creating new edge hostname: %#v", newHostname)
	hostname, err := client.CreateEdgeHostname(ctx, papi.CreateEdgeHostnameRequest{
		EdgeHostname: newHostname,
		ContractID:   contractID,
		GroupID:      groupID,
	})
	if err != nil {
		return diag.FromErr(err)
	}
	if d.HasChange("ttl") {
		edgeHostnameID, err := strconv.Atoi(strings.TrimPrefix(hostname.EdgeHostnameID, "ehn_"))
		if err != nil {
			return diag.FromErr(err)
		}
		hapiClient := HapiClient(meta)
		err = waitForHAPIPropagation(ctx, hapiClient, edgeHostnameID)
		if err != nil {
			return diag.FromErr(err)
		}

		ttl, err := tf.GetIntValueAsInt64("ttl", d)
		if err != nil {
			return diag.FromErr(err)
		}
		patches := []patch{{
			value: strconv.FormatInt(ttl, 10),
			field: "ttl",
			path:  "/ttl",
		}}
		diagnostics := patchEdgeHostname(ctx, d, meta, edgeHostnameID, patches)
		if diagnostics != nil {
			return diagnostics
		}
	}

	d.SetId(hostname.EdgeHostnameID)
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
	groupID = str.AddPrefix(groupID, "grp_")
	if err := d.Set("group_id", groupID); err != nil {
		return diag.FromErr(fmt.Errorf("%w: %s", tf.ErrValueSet, err.Error()))
	}

	contractID, err := tf.GetStringValue("contract_id", d)
	if err != nil {
		return diag.FromErr(err)
	}
	contractID = str.AddPrefix(contractID, "ctr_")
	// set contract/contract_id into ResourceData
	if err := d.Set("contract_id", contractID); err != nil {
		return diag.FromErr(fmt.Errorf("%w: %s", tf.ErrValueSet, err.Error()))
	}

	// Schema guarantees product_id/product are strings and one or the other is set
	var productID string
	if got, ok := d.GetOk("product_id"); ok {
		productID = got.(string)
	}
	productID = str.AddPrefix(productID, "prd_")
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

	_, err = tf.GetIntValueAsInt64("ttl", d)
	if err != nil && !errors.Is(err, tf.ErrNotFound) {
		return diag.FromErr(err)
	}
	if err == nil {
		edgeHostnameID, err := strconv.Atoi(strings.TrimPrefix(foundEdgeHostname.ID, "ehn_"))
		if err != nil {
			return diag.FromErr(err)
		}

		hapiClient := HapiClient(meta)
		// in theory this call is redundant, added here as safeguard
		err = waitForHAPIPropagation(ctx, hapiClient, edgeHostnameID)
		if err != nil {
			return diag.FromErr(err)
		}
		hostname, err := hapiClient.GetEdgeHostname(ctx, edgeHostnameID)
		if err != nil {
			return diag.FromErr(err)
		}
		ttl := hostname.TTL
		if hostname.UseDefaultTTL {
			ttl = 0
		}
		if err := d.Set("ttl", ttl); err != nil {
			return diag.FromErr(fmt.Errorf("%w: %s", tf.ErrValueSet, err.Error()))
		}
	}
	return nil
}

func resourceSecureEdgeHostNameUpdate(ctx context.Context, d *schema.ResourceData, m interface{}) diag.Diagnostics {
	meta := meta.Must(m)
	logger := meta.Log("PAPI", "resourceSecureEdgeHostNameUpdate")

	if !d.HasChangeExcept("timeouts") {
		logger.Debug("Only timeouts were updated, skipping")
		return nil
	}

	patches := make([]patch, 0, 2)
	if d.HasChange("ip_behavior") {
		ipBehavior, err := tf.GetStringValue("ip_behavior", d)
		if err != nil {
			return diag.FromErr(err)
		}
		// IPV6_COMPLIANCE type has to mapped to IPV6_IPV4_DUALSTACK which is only accepted value by HAPI client
		if ipBehavior == papi.EHIPVersionV6Compliance {
			ipBehavior = "IPV6_IPV4_DUALSTACK"
		}
		patches = append(patches, patch{
			value: ipBehavior,
			field: "ip_behavior",
			path:  "/ipVersionBehavior",
		})
	}

	if d.HasChange("ttl") {
		ttl, err := tf.GetIntValueAsInt64("ttl", d)
		if err != nil {
			return diag.FromErr(err)
		}
		patches = append(patches, patch{
			value: strconv.FormatInt(ttl, 10),
			field: "ttl",
			path:  "/ttl",
		})
	}

	if len(patches) > 0 {
		edgeHostnameIDString := d.Id()
		edgeHostnameID, err := strconv.Atoi(strings.TrimPrefix(edgeHostnameIDString, "ehn_"))
		if err != nil {
			return diag.FromErr(err)
		}

		diagnostics := patchEdgeHostname(ctx, d, meta, edgeHostnameID, patches)
		if diagnostics != nil {
			return diagnostics
		}
	}

	return resourceSecureEdgeHostNameRead(ctx, d, m)
}

type patch struct {
	value string
	field string
	path  string
}

func patchEdgeHostname(ctx context.Context, d *schema.ResourceData, meta meta.Meta, edgeHostnameID int, patches []patch) diag.Diagnostics {
	logger := meta.Log("PAPI", "patchEdgeHostname")

	edgeHostname, err := tf.GetStringValue("edge_hostname", d)
	if err != nil {
		return diag.FromErr(err)
	}
	dnsZone, _ := parseEdgeHostname(edgeHostname)
	emails, err := tf.GetListValue("status_update_email", d)
	if err != nil && !errors.Is(err, tf.ErrNotFound) {
		return diag.FromErr(err)
	}

	l := len(patches)
	body := make([]hapi.UpdateEdgeHostnameRequestBody, 0, l)
	fields := make([]string, 0, l)
	comments := make([]string, 0, l)
	for _, p := range patches {
		logger.Debugf("Proceeding to update %s for %s", p.field, edgeHostname)
		body = append(body, hapi.UpdateEdgeHostnameRequestBody{

			Op:    "replace",
			Path:  p.path,
			Value: p.value,
		})
		comments = append(comments, fmt.Sprintf("change %s to %s", p.path, p.value))
		fields = append(fields, p.field)
	}
	req := hapi.UpdateEdgeHostnameRequest{
		DNSZone:    dnsZone,
		RecordName: strings.ReplaceAll(edgeHostname, "."+dnsZone, ""),
		Comments:   strings.Join(comments, "; "),
		Body:       body,
	}

	if len(emails) != 0 {
		statusUpdateEmails := make([]string, len(emails))
		for i, email := range emails {
			statusUpdateEmails[i] = email.(string)
		}
		req.StatusUpdateEmail = statusUpdateEmails
	}

	hapiClient := HapiClient(meta)
	err = waitForHAPIPropagation(ctx, hapiClient, edgeHostnameID)
	if err != nil {
		return diag.FromErr(err)
	}
	resp, err := hapiClient.UpdateEdgeHostname(ctx, req)
	if err != nil {
		if err2 := tf.RestoreOldValues(d, fields); err2 != nil {
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
	return nil
}

func waitForHAPIPropagation(ctx context.Context, hapiClient hapi.HAPI, edgeHostnameID int) error {
	retries := 0
	for {
		select {
		case <-time.After(EgdeHostnameCreatePollInterval):
			resp, err := hapiClient.GetEdgeHostname(ctx, edgeHostnameID)
			if resp == nil && err != nil {
				var target = &hapi.Error{}
				if !errors.As(err, &target) {
					return fmt.Errorf("error has unexpected type: %T", err)
				}
				if target.Status != 200 {
					retries++
					if retries > retriesMax {
						return fmt.Errorf("reached max number of retries: %d", retries-1)
					}
					continue
				}
			}

			return nil

		case <-ctx.Done():
			return fmt.Errorf("update edge hostname context terminated: %s", ctx.Err())
		}
	}
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

func resourceSecureEdgeHostNameDelete(_ context.Context, d *schema.ResourceData, m interface{}) diag.Diagnostics {
	meta := meta.Must(m)
	logger := meta.Log("PAPI", "resourceSecureEdgeHostNameDelete")
	logger.Debug("DELETING")
	logger.Info("PAPI does not support edge hostname deletion - resource will only be removed from state")
	d.SetId("")
	logger.Debugf("DONE")
	return nil
}

// resourceSecureEdgeHostNameImport accepts the following import ID:
// EdgehostNameID,contractID,groupID[,productID]
// productID is optional and needs to be specified if the resource config contains the product_id
// attribute (which is not the case for the configuration generated by export done with Akamai
// CLI tools).
func resourceSecureEdgeHostNameImport(ctx context.Context, d *schema.ResourceData, m interface{}) ([]*schema.ResourceData, error) {
	meta := meta.Must(m)
	client := Client(meta)
	logger := meta.Log("PAPI", "resourceSecureEdgeHostNameImport")

	parts := strings.Split(d.Id(), ",")
	if len(parts) < 3 || len(parts) > 4 {
		return nil, fmt.Errorf("expected import identifier with format: "+
			`"EdgehostNameID,contractID,groupID[,productID]". Got: %q`, d.Id())
	}

	if len(parts) == 4 {
		if len(parts[3]) == 0 {
			return nil, fmt.Errorf("productID is empty for the import ID=%q", d.Id())
		}
		productID := str.AddPrefix(parts[3], "prd_")
		logger.Debugf("Setting product_id=%s", productID)
		if err := d.Set("product_id", productID); err != nil {
			return nil, fmt.Errorf("%w: %s", tf.ErrValueSet, err.Error())
		}
	}

	edgehostID := parts[0]
	contractID := str.AddPrefix(parts[1], "ctr_")
	groupID := str.AddPrefix(parts[2], "grp_")

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
	if err := d.Set("edge_hostname", edgehostnameDetails.EdgeHostname.Domain); err != nil {
		return nil, fmt.Errorf("%w: %s", tf.ErrValueSet, err.Error())
	}

	hapiClient := HapiClient(meta)
	edgeHostnameID, err := str.GetIntID(edgehostID, "ehn_")
	if err != nil {
		return nil, err
	}
	edgeHostnameResp, err := hapiClient.GetEdgeHostname(ctx, edgeHostnameID)
	if err != nil {
		return nil, fmt.Errorf("error getting edge hostname with id '%d': it may not be ready in HAPI yet: %s", edgeHostnameID, err)
	}

	// get certificate id when network is ENHANCED-TLS
	if edgeHostnameResp.SecurityType == "ENHANCED-TLS" {
		certificate, err := hapiClient.GetCertificate(ctx, hapi.GetCertificateRequest{
			DNSZone:    edgeHostnameResp.DNSZone,
			RecordName: edgeHostnameResp.RecordName,
		})
		if err != nil {
			if !errors.Is(err, hapi.ErrNotFound) {
				return nil, err
			}
		} else {
			certificateID, err := strconv.ParseInt(certificate.CertificateID, 10, 64)
			if err != nil {
				return nil, err
			}
			if err := d.Set("certificate", certificateID); err != nil {
				return nil, fmt.Errorf("%w: %s", tf.ErrValueSet, err.Error())
			}
		}
	}

	if !edgeHostnameResp.UseDefaultTTL {
		err := d.Set("ttl", edgeHostnameResp.TTL)
		if err != nil {
			return nil, fmt.Errorf("%w: %s", tf.ErrValueSet, err.Error())
		}
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
	logger := log.Get("PAPI", "suppressEdgeHostnameUseCases")
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

func validateImmutableFields(_ context.Context, diff *schema.ResourceDiff, _ interface{}) error {
	if diff.Id() != "" {
		oldValue, newValue := diff.GetChange("product_id")
		o := oldValue.(string)
		n := newValue.(string)

		if diff.HasChange("certificate") || str.AddPrefix(o, "prd_") != str.AddPrefix(n, "prd_") {
			return fmt.Errorf("error: Changes to non-updatable fields 'product_id' and 'certificate' are not permitted")
		}
	}
	return nil
}
