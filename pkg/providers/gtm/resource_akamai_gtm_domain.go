package gtm

import (
	"context"
	"errors"
	"fmt"
	"net/http"
	"strings"
	"time"

	"github.com/akamai/AkamaiOPEN-edgegrid-golang/v12/pkg/gtm"
	"github.com/akamai/AkamaiOPEN-edgegrid-golang/v12/pkg/session"
	"github.com/akamai/terraform-provider-akamai/v9/pkg/common/ptr"
	"github.com/akamai/terraform-provider-akamai/v9/pkg/common/tf"
	"github.com/akamai/terraform-provider-akamai/v9/pkg/meta"
	"github.com/hashicorp/go-cty/cty"
	"github.com/hashicorp/terraform-plugin-sdk/v2/diag"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/schema"
)

// HashiAcc is Hack for Hashicorp Acceptance Tests
var HashiAcc = false
var sleepInterval = 5 * time.Second
var defaultInterval = 5 * time.Second

const domainMapAlreadyExistsError = "Domain with provided `name` already exists. Please import specific domain using following command: terraform import akamai_gtm_domain.<your_resource_name> \"%s\""

func resourceGTMv1Domain() *schema.Resource {
	return &schema.Resource{
		CreateContext: resourceGTMv1DomainCreate,
		ReadContext:   resourceGTMv1DomainRead,
		UpdateContext: resourceGTMv1DomainUpdate,
		DeleteContext: resourceGTMv1DomainDelete,
		Importer: &schema.ResourceImporter{
			StateContext: resourceGTMv1DomainImport,
		},
		Schema: map[string]*schema.Schema{
			"contract": {
				Type:             schema.TypeString,
				Optional:         true,
				Default:          "",
				DiffSuppressFunc: tf.FieldPrefixSuppress("ctr_"),
			},
			"group": {
				Type:             schema.TypeString,
				Optional:         true,
				Default:          "",
				DiffSuppressFunc: tf.FieldPrefixSuppress("grp_"),
			},
			"wait_on_complete": {
				Type:     schema.TypeBool,
				Optional: true,
				Default:  true,
			},
			"name": {
				Type:     schema.TypeString,
				Required: true,
			},
			"type": {
				Type:             schema.TypeString,
				Required:         true,
				ValidateDiagFunc: validateDomainType,
			},
			"comment": {
				Type:     schema.TypeString,
				Optional: true,
				Default:  "Managed by Terraform",
			},
			"default_unreachable_threshold": {
				Type:     schema.TypeFloat,
				Computed: true,
			},
			"email_notification_list": {
				Type:     schema.TypeSet,
				Elem:     &schema.Schema{Type: schema.TypeString},
				Optional: true,
			},
			"min_pingable_region_fraction": {
				Type:     schema.TypeFloat,
				Computed: true,
			},
			"default_timeout_penalty": {
				Type:     schema.TypeInt,
				Optional: true,
				Default:  25,
			},
			"servermonitor_liveness_count": {
				Type:     schema.TypeInt,
				Computed: true,
			},
			"round_robin_prefix": {
				Type:     schema.TypeString,
				Computed: true,
			},
			"servermonitor_load_count": {
				Type:     schema.TypeInt,
				Computed: true,
			},
			"ping_interval": {
				Type:     schema.TypeInt,
				Computed: true,
			},
			"max_ttl": {
				Type:     schema.TypeInt,
				Computed: true,
			},
			"load_imbalance_percentage": {
				Type:     schema.TypeFloat,
				Optional: true,
			},
			"default_health_max": {
				Type:     schema.TypeFloat,
				Computed: true,
			},
			"map_update_interval": {
				Type:     schema.TypeInt,
				Computed: true,
			},
			"max_properties": {
				Type:     schema.TypeInt,
				Computed: true,
			},
			"max_resources": {
				Type:     schema.TypeInt,
				Computed: true,
			},
			"default_ssl_client_private_key": {
				Type:     schema.TypeString,
				Optional: true,
			},
			"default_error_penalty": {
				Type:     schema.TypeInt,
				Optional: true,
				Default:  75,
			},
			"max_test_timeout": {
				Type:     schema.TypeFloat,
				Computed: true,
			},
			"cname_coalescing_enabled": {
				Type:     schema.TypeBool,
				Optional: true,
			},
			"default_health_multiplier": {
				Type:     schema.TypeFloat,
				Computed: true,
			},
			"servermonitor_pool": {
				Type:     schema.TypeString,
				Computed: true,
			},
			"load_feedback": {
				Type:     schema.TypeBool,
				Optional: true,
			},
			"min_ttl": {
				Type:     schema.TypeInt,
				Computed: true,
			},
			"default_max_unreachable_penalty": {
				Type:     schema.TypeInt,
				Computed: true,
			},
			"default_health_threshold": {
				Type:     schema.TypeFloat,
				Computed: true,
			},
			"min_test_interval": {
				Type:     schema.TypeInt,
				Computed: true,
			},
			"ping_packet_size": {
				Type:     schema.TypeInt,
				Computed: true,
			},
			"default_ssl_client_certificate": {
				Type:     schema.TypeString,
				Optional: true,
			},
			"end_user_mapping_enabled": {
				Type:     schema.TypeBool,
				Optional: true,
			},
			"sign_and_serve": {
				Type:        schema.TypeBool,
				Optional:    true,
				Description: "If set (true) we will sign the domain's resource records so that they can be validated by a validating resolver.",
			},
			"sign_and_serve_algorithm": {
				Type:        schema.TypeString,
				Optional:    true,
				Description: "The signing algorithm to use for signAndServe. One of the following values: RSA_SHA1, RSA_SHA256, RSA_SHA512, ECDSA_P256_SHA256, ECDSA_P384_SHA384, ED25519, ED448.",
			},
		},
	}
}

// GetQueryArgs retrieves optional query args. contractId, groupId [and accountSwitchKey] supported.
func GetQueryArgs(d *schema.ResourceData) (*gtm.DomainQueryArgs, error) {

	qArgs := gtm.DomainQueryArgs{}
	contractName, err := tf.GetStringValue("contract", d)
	if err != nil {
		return nil, fmt.Errorf("contract not present in resource data: %v", err.Error())
	}
	contract := strings.TrimPrefix(contractName, "ctr_")
	if contract != "" && len(contract) > 0 {
		qArgs.ContractID = contract
	}
	groupName, err := tf.GetStringValue("group", d)
	if err != nil {
		return nil, fmt.Errorf("group not present in resource data: %v", err.Error())
	}
	groupID := strings.TrimPrefix(groupName, "grp_")
	if groupID != "" && len(groupID) > 0 {
		qArgs.GroupID = groupID
	}

	return &qArgs, nil
}

// Create a new GTM Domain
func resourceGTMv1DomainCreate(ctx context.Context, d *schema.ResourceData, m interface{}) diag.Diagnostics {
	meta := meta.Must(m)
	logger := meta.Log("Akamai GTM", "resourceGTMv1DomainCreate")
	// create a context with logging for api calls
	ctx = session.ContextWithOptions(
		ctx,
		session.WithContextLog(logger),
	)

	var diags diag.Diagnostics
	dname, err := tf.GetStringValue("name", d)
	if err != nil {
		logger.Errorf("Domain name not found in ResourceData")
		return diag.FromErr(err)
	}
	dom, err := Client(meta).GetDomain(ctx, gtm.GetDomainRequest{
		DomainName: dname,
	})
	if err != nil && !errors.Is(err, gtm.ErrNotFound) {
		logger.Errorf("Domain read error: %s", err.Error())
		return append(diags, diag.Diagnostic{
			Severity: diag.Error,
			Summary:  "Domain read error",
			Detail:   err.Error(),
		})
	}
	if dom != nil {
		domainAlreadyExists := fmt.Sprintf(domainMapAlreadyExistsError, dname)
		logger.Errorf(domainAlreadyExists)
		return append(diags, diag.Diagnostic{
			Severity: diag.Error,
			Summary:  "domain already exists error",
			Detail:   domainAlreadyExists,
		})
	}

	logger.Infof("Creating domain [%s]", dname)
	newDom, err := populateNewDomainObject(d, m)
	if err != nil {
		return diag.FromErr(err)
	}
	logger.Debugf("Domain: [%v]", newDom)
	queryArgs, err := GetQueryArgs(d)
	if err != nil {
		logger.Errorf("Domain create error: %s", err.Error())
		return append(diags, diag.Diagnostic{
			Severity: diag.Error,
			Summary:  "Domain create error",
			Detail:   err.Error(),
		})
	}
	cStatus, err := Client(meta).CreateDomain(ctx, gtm.CreateDomainRequest{
		Domain:    newDom,
		QueryArgs: queryArgs,
	})
	if err != nil {
		// Errored. Let's see if special hack
		if !HashiAcc {
			logger.Errorf("Domain create error: %s", err.Error())
			return append(diags, diag.Diagnostic{
				Severity: diag.Error,
				Summary:  "Domain create error",
				Detail:   err.Error(),
			})
		}
		apiError, ok := err.(*gtm.Error)
		if !ok && apiError.StatusCode != http.StatusBadRequest {
			logger.Errorf("Domain create error: %s", err.Error())
			return append(diags, diag.Diagnostic{
				Severity: diag.Error,
				Summary:  "Domain create error",
				Detail:   err.Error(),
			})
		}
		if strings.Contains(apiError.Detail, "proposed domain name") && strings.Contains(apiError.Detail, "Domain Validation Error") {
			// Already exists
			logger.Warnf("Domain %s already exists. Ignoring error (Hashicorp).", dname)
		} else {
			logger.Errorf("Domain create error: %s", err.Error())
			return append(diags, diag.Diagnostic{
				Severity: diag.Error,
				Summary:  "Domain create error",
				Detail:   err.Error(),
			})
		}
	} else {
		logger.Debugf("Create status: %v", cStatus.Status)
		if cStatus.Status.PropagationStatus == "DENIED" {
			logger.Errorf(cStatus.Status.Message)
			return append(diags, diag.Diagnostic{
				Severity: diag.Error,
				Summary:  cStatus.Status.Message,
			})
		}

		waitOnComplete, err := tf.GetBoolValue("wait_on_complete", d)
		if err != nil {
			return diag.FromErr(err)
		}

		if waitOnComplete {
			done, err := waitForCompletion(ctx, dname, m)
			if done {
				logger.Infof("Domain create completed")
			} else {
				if err == nil {
					logger.Infof("Domain create pending")
				} else {
					logger.Errorf("Domain create error: %s", err.Error())
					return append(diags, diag.Diagnostic{
						Severity: diag.Error,
						Summary:  "Domain create error",
						Detail:   err.Error(),
					})
				}
			}
		}
	}
	// Give terraform the ID
	d.SetId(dname)
	return resourceGTMv1DomainRead(ctx, d, m)

}

// Only ever save data from the tf config in the tf state file, to help with
// api issues. See func unmarshalResourceData for more info.
func resourceGTMv1DomainRead(ctx context.Context, d *schema.ResourceData, m interface{}) diag.Diagnostics {
	meta := meta.Must(m)
	logger := meta.Log("Akamai GTM", "resourceGTMv1DomainRead")
	// create a context with logging for api calls
	ctx = session.ContextWithOptions(
		ctx,
		session.WithContextLog(logger),
	)

	logger.Debugf("Reading Domain: %s", d.Id())
	var diags diag.Diagnostics
	// retrieve the domain
	dom, err := Client(meta).GetDomain(ctx, gtm.GetDomainRequest{
		DomainName: d.Id(),
	})
	if errors.Is(err, gtm.ErrNotFound) {
		d.SetId("")
		return nil
	}
	if err != nil {
		logger.Errorf("Domain read error: %s", err.Error())
		return append(diags, diag.Diagnostic{
			Severity: diag.Error,
			Summary:  "Domain read error",
			Detail:   err.Error(),
		})
	}
	populateTerraformState(d, dom, m)
	logger.Debugf("READ %v", dom)
	return nil
}

// Update GTM Domain
func resourceGTMv1DomainUpdate(ctx context.Context, d *schema.ResourceData, m interface{}) diag.Diagnostics {
	meta := meta.Must(m)
	logger := meta.Log("Akamai GTM", "resourceGTMv1DomainUpdate")
	// create a context with logging for api calls
	ctx = session.ContextWithOptions(
		ctx,
		session.WithContextLog(logger),
	)

	logger.Debugf("Updating Domain: %s", d.Id())
	var diags diag.Diagnostics
	// Get existing domain
	existDom, err := Client(meta).GetDomain(ctx, gtm.GetDomainRequest{
		DomainName: d.Id(),
	})
	if err != nil {
		logger.Errorf("Domain read error: %s", err.Error())
		return append(diags, diag.Diagnostic{
			Severity: diag.Error,
			Summary:  "Domain read error",
			Detail:   err.Error(),
		})
	}
	logger.Debugf("Updating Domain BEFORE: %v", existDom)
	newDom := createDomainStruct(existDom)
	err = populateDomainObject(d, newDom, m)
	if err != nil {
		return diag.FromErr(err)
	}
	logger.Debugf("Updating Domain PROPOSED: %v", newDom)
	args, err := GetQueryArgs(d)
	if err != nil {
		logger.Errorf("Domain update error: %s", err.Error())
		return append(diags, diag.Diagnostic{
			Severity: diag.Error,
			Summary:  "Domain update error",
			Detail:   err.Error(),
		})
	}

	uStat, err := Client(meta).UpdateDomain(ctx, gtm.UpdateDomainRequest{
		Domain:    newDom,
		QueryArgs: args,
	})
	if err != nil {
		logger.Errorf("Domain update error: %s", err.Error())
		return append(diags, diag.Diagnostic{
			Severity: diag.Error,
			Summary:  "Domain update error",
			Detail:   err.Error(),
		})
	}
	logger.Debugf("Update status: %v", uStat)
	if uStat.Status.PropagationStatus == "DENIED" {
		logger.Errorf(uStat.Status.Message)
		return append(diags, diag.Diagnostic{
			Severity: diag.Error,
			Summary:  uStat.Status.Message,
		})
	}

	waitOnComplete, err := tf.GetBoolValue("wait_on_complete", d)
	if err != nil {
		return diag.FromErr(err)
	}

	if waitOnComplete {
		done, err := waitForCompletion(ctx, d.Id(), m)
		if done {
			logger.Infof("Domain update completed")
		} else {
			if err == nil {
				logger.Infof("Domain update pending")
			} else {
				logger.Errorf("Domain update error: %s", err.Error())
				return append(diags, diag.Diagnostic{
					Severity: diag.Error,
					Summary:  "domain update error",
					Detail:   err.Error(),
				})
			}
		}

	}

	return resourceGTMv1DomainRead(ctx, d, m)

}

// Delete an existing GTM Domain
func resourceGTMv1DomainDelete(ctx context.Context, d *schema.ResourceData, m interface{}) diag.Diagnostics {
	meta := meta.Must(m)
	logger := meta.Log("Akamai GTM", "resourceGTMv1DomainDelete")
	ctx = session.ContextWithOptions(
		ctx,
		session.WithContextLog(logger),
	)

	var diags diag.Diagnostics
	var domainName = d.Id()
	logger.Debugf("Initiating delete request for GTM domain: %s", domainName)

	resp, err := Client(meta).DeleteDomains(ctx, gtm.DeleteDomainsRequest{
		Body: gtm.DeleteDomainsRequestBody{
			DomainNames: []string{domainName},
		},
	})
	if err != nil {
		if errors.Is(err, gtm.ErrDomainNotFound) {
			logger.Warnf("Domain %s not found or appears to already be deleted.", domainName)
			d.SetId("")
			return diag.FromErr(fmt.Errorf("domain %s not found or appears to already be deleted: %w", domainName, err))
		}
		logger.Errorf("DeleteDomains API error: %v", err)
		return diag.FromErr(fmt.Errorf("failed to submit delete request for domain %s: %w", domainName, err))
	}

	logger.Debugf("Check Delete Domain Status for requestID: %v", resp.RequestID)

	status, err := waitForDeletion(ctx, resp.RequestID, meta)

	if err != nil {
		logger.Errorf("Domain delete error: %s", err.Error())
		return append(diags, diag.Diagnostic{
			Severity: diag.Error,
			Summary:  "Domain delete error",
			Detail:   err.Error(),
		})
	}

	if status.IsComplete && status.SuccessCount == 0 {
		logger.Errorf("Domain deletion failed domain: %s", domainName)
		return append(diags, diag.Diagnostic{
			Severity: diag.Error,
			Summary:  "Domain deletion failed",
			Detail:   "No additional details available.",
		})
	}

	logger.Infof("Domain delete completed")
	d.SetId("")
	return nil
}

// waitForDeletion waits for the deletion process to complete by polling the status.
func waitForDeletion(ctx context.Context, requestID string, meta meta.Meta) (*gtm.DeleteDomainsStatusResponse, error) {
	logger := meta.Log("Akamai GTM", "waitForDeletion")

	const timeoutDuration = 300 * time.Second

	ctx, cancel := context.WithTimeout(ctx, timeoutDuration)
	defer cancel()

	ticker := time.NewTicker(sleepInterval)
	defer ticker.Stop()

	for {
		select {
		case <-ctx.Done():
			return nil, fmt.Errorf("context cancelled while waiting for domain deletion: %w", ctx.Err())
		case <-ticker.C:
			status, err := Client(meta).GetDeleteDomainsStatus(ctx, gtm.DeleteDomainsStatusRequest{
				RequestID: requestID,
			})

			if err != nil {
				return nil, fmt.Errorf("error checking delete domain status for request %s: %w", requestID, err)
			}

			logger.Debugf("Delete status for request %s: %v", requestID, status)

			if status.IsComplete {
				return status, nil
			}

			logger.Debugf("Waiting for deletion...")
		}
	}
}

func resourceGTMv1DomainImport(_ context.Context, d *schema.ResourceData, m interface{}) ([]*schema.ResourceData, error) {
	meta := meta.Must(m)
	logger := meta.Log("Akamai GTM", "resourceGTMv1DomainImport")

	// User-supplied import ID is a comma-separated list of domain,[,groupID[,contractID]]
	// groupID and contractID are optional
	parts := strings.Split(d.Id(), ",")

	if len(parts) > 3 {
		return nil, fmt.Errorf("invalid importID format: %v. An importID must be in the format 'domain[,groupID[,contractID]]'. "+
			"It can include up to three comma-separated parts: domain (required), groupID (optional), and contractID (optional)", parts)
	}
	domain := parts[0]
	var group, contract string
	if len(parts) > 1 {
		group = parts[1]
		if err := d.Set("group", group); err != nil {
			return nil, err
		}
	}
	if len(parts) > 2 {
		contract = parts[2]
		if err := d.Set("contract", contract); err != nil {
			return nil, err
		}
	}

	logger.Debugf("Importing GTM Domain: %s", domain)

	if err := d.Set("wait_on_complete", true); err != nil {
		return nil, err
	}

	d.SetId(domain)
	return []*schema.ResourceData{d}, nil
}

// validateDomainType is a SchemaValidateFunc to validate the Domain type.
func validateDomainType(v interface{}, _ cty.Path) diag.Diagnostics {
	value := strings.ToUpper(v.(string))
	if value != "BASIC" && value != "FULL" && value != "WEIGHTED" && value != "STATIC" && value != "FAILOVER-ONLY" {
		return diag.Errorf("type must be basic, full, weighted, static, or failover-only")
	}
	return nil
}

// Create and populate a new domain object from resource data
func populateNewDomainObject(d *schema.ResourceData, m interface{}) (*gtm.Domain, error) {

	name, err := tf.GetStringValue("name", d)
	if err != nil {
		return nil, err
	}
	domObj := &gtm.Domain{
		Name: name,
		Type: d.Get("type").(string),
	}
	err = populateDomainObject(d, domObj, m)

	return domObj, err

}

// nolint:gocyclo
// Populate existing domain object from resource data
func populateDomainObject(d *schema.ResourceData, dom *gtm.Domain, m interface{}) error {
	meta := meta.Must(m)
	logger := meta.Log("Akamai GTM", "populateDomainObject")

	domainName, err := tf.GetStringValue("name", d)
	if err != nil {
		// Should be caught earlier ...
		logger.Warnf("Domain name not set: %s", err.Error())
	}

	if domainName != dom.Name {
		logger.Errorf("Domain [%s] state and GTM names inconsistent!", dom.Name)
		return fmt.Errorf("once the domain is created, updating its name is not allowed")
	}

	vstr, err := tf.GetStringValue("type", d)
	if err == nil {
		if vstr != dom.Type {
			dom.Type = vstr
		}
	}
	vfl32, err := tf.GetFloat32Value("default_unreachable_threshold", d)
	if err == nil {
		dom.DefaultUnreachableThreshold = vfl32
	}
	vlist, err := tf.GetSetValue("email_notification_list", d)
	if err == nil {
		ls := make([]string, vlist.Len())
		for i, sl := range vlist.List() {
			ls[i] = sl.(string)
		}
		dom.EmailNotificationList = ls
	} else if d.HasChange("email_notification_list") {
		dom.EmailNotificationList = make([]string, 0)
	}
	vfl32, err = tf.GetFloat32Value("min_pingable_region_fraction", d)
	if err == nil {
		dom.MinPingableRegionFraction = vfl32
	}
	vint, err := tf.GetIntValue("default_timeout_penalty", d)
	if err == nil || d.HasChange("default_timeout_penalty") {
		dom.DefaultTimeoutPenalty = vint
	}
	if err != nil && !errors.Is(err, tf.ErrNotFound) {
		logger.Errorf("populateResourceObject() default_timeout_penalty failed: %v", err.Error())
		return fmt.Errorf("Domain Object could not be populated: %v", err.Error())
	}

	vint, err = tf.GetIntValue("servermonitor_liveness_count", d)
	if err == nil {
		dom.ServermonitorLivenessCount = vint
	}
	vstr, err = tf.GetStringValue("round_robin_prefix", d)
	if err == nil {
		dom.RoundRobinPrefix = vstr
	}
	vint, err = tf.GetIntValue("servermonitor_load_count", d)
	if err == nil {
		dom.ServermonitorLoadCount = vint
	}
	vint, err = tf.GetIntValue("ping_interval", d)
	if err == nil {
		dom.PingInterval = vint
	}
	vint, err = tf.GetIntValue("max_ttl", d)
	if err == nil {
		dom.MaxTTL = int64(vint)
	}
	vfloat, err := tf.GetFloat64Value("load_imbalance_percentage", d)
	if err == nil {
		dom.LoadImbalancePercentage = vfloat
	}
	if err != nil && !errors.Is(err, tf.ErrNotFound) {
		logger.Errorf("populateResourceObject() load_imbalance_percentage failed: %v", err.Error())
		return fmt.Errorf("Domain Object could not be populated: %v", err.Error())
	}

	vfloat, err = tf.GetFloat64Value("default_health_max", d)
	if err == nil {
		dom.DefaultHealthMax = vfloat
	}
	vint, err = tf.GetIntValue("map_update_interval", d)
	if err == nil {
		dom.MapUpdateInterval = vint
	}
	vint, err = tf.GetIntValue("max_properties", d)
	if err == nil {
		dom.MaxProperties = vint
	}
	vint, err = tf.GetIntValue("max_resources", d)
	if err == nil {
		dom.MaxResources = vint
	}
	vstr, err = tf.GetStringValue("default_ssl_client_private_key", d)
	if err == nil || d.HasChange("default_ssl_client_private_key") {
		dom.DefaultSSLClientPrivateKey = vstr
	}
	if err != nil && !errors.Is(err, tf.ErrNotFound) {
		logger.Errorf("populateResourceObject() default_ssl_client_private_key failed: %v", err.Error())
		return fmt.Errorf("Domain Object could not be populated: %v", err.Error())
	}

	vint, err = tf.GetIntValue("default_error_penalty", d)
	if err == nil || d.HasChange("default_error_penalty") {
		dom.DefaultErrorPenalty = vint
	}
	if err != nil && !errors.Is(err, tf.ErrNotFound) {
		logger.Errorf("populateResourceObject() default_error_penalty failed: %v", err.Error())
		return fmt.Errorf("Domain Object could not be populated: %v", err.Error())
	}

	vfloat, err = tf.GetFloat64Value("max_test_timeout", d)
	if err == nil {
		dom.MaxTestTimeout = vfloat
	}
	if cnameCoalescingEnabled, err := tf.GetBoolValue("cname_coalescing_enabled", d); err == nil {
		dom.CNameCoalescingEnabled = cnameCoalescingEnabled
	}
	vfloat, err = tf.GetFloat64Value("default_health_multiplier", d)
	if err == nil {
		dom.DefaultHealthMultiplier = vfloat
	}
	vstr, err = tf.GetStringValue("servermonitor_pool", d)
	if err == nil {
		dom.ServermonitorPool = vstr
	}
	if loadFeedback, err := tf.GetBoolValue("load_feedback", d); err == nil {
		dom.LoadFeedback = loadFeedback
	}
	vint, err = tf.GetIntValue("min_ttl", d)
	if err == nil {
		dom.MinTTL = int64(vint)
	}
	vint, err = tf.GetIntValue("default_max_unreachable_penalty", d)
	if err == nil {
		dom.DefaultMaxUnreachablePenalty = vint
	}
	vfloat, err = tf.GetFloat64Value("default_health_threshold", d)
	if err == nil {
		dom.DefaultHealthThreshold = vfloat
	}
	vstr, err = tf.GetStringValue("comment", d)
	if err == nil || d.HasChange("comment") {
		dom.ModificationComments = vstr
	}
	if err != nil && !errors.Is(err, tf.ErrNotFound) {
		logger.Errorf("populateResourceObject() comment failed: %v", err.Error())
		return fmt.Errorf("Domain Object could not be populated: %v", err.Error())
	}

	vint, err = tf.GetIntValue("min_test_interval", d)
	if err == nil {
		dom.MinTestInterval = vint
	}
	vint, err = tf.GetIntValue("ping_packet_size", d)
	if err == nil {
		dom.PingPacketSize = vint
	}
	vstr, err = tf.GetStringValue("default_ssl_client_certificate", d)
	if err == nil || d.HasChange("default_ssl_client_certificate") {
		dom.DefaultSSLClientCertificate = vstr
	}
	if err != nil && !errors.Is(err, tf.ErrNotFound) {
		logger.Errorf("populateResourceObject() default_ssl_client_certificate failed: %v", err.Error())
		return fmt.Errorf("Domain Object could not be populated: %v", err.Error())
	}

	if vbool, err := tf.GetBoolValue("end_user_mapping_enabled", d); err == nil {
		dom.EndUserMappingEnabled = vbool
	}

	signAndServe, err := tf.GetBoolValue("sign_and_serve", d)
	if err != nil && !errors.Is(err, tf.ErrNotFound) {
		return fmt.Errorf("could not get `sign_and_serve` attribute: %s", err)
	}
	dom.SignAndServe = signAndServe
	signAndServeAlgorithm, err := tf.GetStringValue("sign_and_serve_algorithm", d)
	if err != nil && !errors.Is(err, tf.ErrNotFound) {
		return fmt.Errorf("could not get `sign_and_serve_algorithm` attribute: %s", err)
	}
	if signAndServeAlgorithm != "" {
		dom.SignAndServeAlgorithm = ptr.To(signAndServeAlgorithm)
	}

	return nil

}

// Populate Terraform state from provided Domain object
func populateTerraformState(d *schema.ResourceData, dom *gtm.GetDomainResponse, m interface{}) {
	meta := meta.Must(m)
	logger := meta.Log("Akamai GTM", "populateTerraformState")

	for stateKey, stateValue := range map[string]interface{}{
		"name":                            dom.Name,
		"type":                            dom.Type,
		"default_unreachable_threshold":   dom.DefaultUnreachableThreshold,
		"email_notification_list":         dom.EmailNotificationList,
		"min_pingable_region_fraction":    dom.MinPingableRegionFraction,
		"default_timeout_penalty":         dom.DefaultTimeoutPenalty,
		"servermonitor_liveness_count":    dom.ServermonitorLivenessCount,
		"round_robin_prefix":              dom.RoundRobinPrefix,
		"servermonitor_load_count":        dom.ServermonitorLoadCount,
		"ping_interval":                   dom.PingInterval,
		"max_ttl":                         dom.MaxTTL,
		"load_imbalance_percentage":       dom.LoadImbalancePercentage,
		"default_health_max":              dom.DefaultHealthMax,
		"map_update_interval":             dom.MapUpdateInterval,
		"max_properties":                  dom.MaxProperties,
		"max_resources":                   dom.MaxResources,
		"default_ssl_client_private_key":  dom.DefaultSSLClientPrivateKey,
		"default_error_penalty":           dom.DefaultErrorPenalty,
		"max_test_timeout":                dom.MaxTestTimeout,
		"cname_coalescing_enabled":        dom.CNameCoalescingEnabled,
		"default_health_multiplier":       dom.DefaultHealthMultiplier,
		"servermonitor_pool":              dom.ServermonitorPool,
		"load_feedback":                   dom.LoadFeedback,
		"min_ttl":                         dom.MinTTL,
		"default_max_unreachable_penalty": dom.DefaultMaxUnreachablePenalty,
		"default_health_threshold":        dom.DefaultHealthThreshold,
		"comment":                         dom.ModificationComments,
		"min_test_interval":               dom.MinTestInterval,
		"ping_packet_size":                dom.PingPacketSize,
		"default_ssl_client_certificate":  dom.DefaultSSLClientCertificate,
		"sign_and_serve":                  dom.SignAndServe,
		"end_user_mapping_enabled":        dom.EndUserMappingEnabled,
	} {
		// walk through all state elements
		err := d.Set(stateKey, stateValue)
		if err != nil {
			logger.Errorf("populateTerraformState failed: %s", err.Error())
		}
	}
	if dom.SignAndServeAlgorithm != nil {
		err := d.Set("sign_and_serve_algorithm", dom.SignAndServeAlgorithm)
		if err != nil {
			logger.Errorf("populateTerraformState failed: %s", err.Error())
		}
	}
}

// createDomainStruct converts response from GetDomainResponse into Domain
func createDomainStruct(domain *gtm.GetDomainResponse) *gtm.Domain {
	if domain != nil {
		return &gtm.Domain{
			Name:                         domain.Name,
			Type:                         domain.Type,
			ASMaps:                       domain.ASMaps,
			Resources:                    domain.Resources,
			DefaultUnreachableThreshold:  domain.DefaultUnreachableThreshold,
			EmailNotificationList:        domain.EmailNotificationList,
			MinPingableRegionFraction:    domain.MinPingableRegionFraction,
			DefaultTimeoutPenalty:        domain.DefaultTimeoutPenalty,
			Datacenters:                  domain.Datacenters,
			ServermonitorLivenessCount:   domain.ServermonitorLivenessCount,
			RoundRobinPrefix:             domain.RoundRobinPrefix,
			ServermonitorLoadCount:       domain.ServermonitorLoadCount,
			PingInterval:                 domain.PingInterval,
			MaxTTL:                       domain.MaxTTL,
			LoadImbalancePercentage:      domain.LoadImbalancePercentage,
			DefaultHealthMax:             domain.DefaultHealthMax,
			LastModified:                 domain.LastModified,
			Status:                       domain.Status,
			MapUpdateInterval:            domain.MapUpdateInterval,
			MaxProperties:                domain.MaxProperties,
			MaxResources:                 domain.MaxResources,
			DefaultSSLClientPrivateKey:   domain.DefaultSSLClientPrivateKey,
			DefaultErrorPenalty:          domain.DefaultErrorPenalty,
			Links:                        domain.Links,
			Properties:                   domain.Properties,
			MaxTestTimeout:               domain.MaxTestTimeout,
			CNameCoalescingEnabled:       domain.CNameCoalescingEnabled,
			DefaultHealthMultiplier:      domain.DefaultHealthMultiplier,
			ServermonitorPool:            domain.ServermonitorPool,
			LoadFeedback:                 domain.LoadFeedback,
			MinTTL:                       domain.MinTTL,
			GeographicMaps:               domain.GeographicMaps,
			CIDRMaps:                     domain.CIDRMaps,
			DefaultMaxUnreachablePenalty: domain.DefaultMaxUnreachablePenalty,
			DefaultHealthThreshold:       domain.DefaultHealthThreshold,
			LastModifiedBy:               domain.LastModifiedBy,
			ModificationComments:         domain.ModificationComments,
			MinTestInterval:              domain.MinTestInterval,
			PingPacketSize:               domain.PingPacketSize,
			DefaultSSLClientCertificate:  domain.DefaultSSLClientCertificate,
			EndUserMappingEnabled:        domain.EndUserMappingEnabled,
			SignAndServe:                 domain.SignAndServe,
			SignAndServeAlgorithm:        domain.SignAndServeAlgorithm,
		}
	}
	return nil
}

// Util function to wait for change deployment. return true if complete. false if not - error or nil (timeout)
func waitForCompletion(ctx context.Context, domain string, m interface{}) (bool, error) {
	meta := meta.Must(m)
	logger := meta.Log("Akamai GTMv1", "waitForCompletion")

	var defaultTimeout = 300 * time.Second
	var sleepInterval = defaultInterval // seconds. TODO:Should be configurable by user ...
	var sleepTimeout = defaultTimeout   // seconds. TODO: Should be configurable by user ...
	if HashiAcc {
		// Override for ACC tests
		sleepTimeout = sleepInterval
	}
	logger.Debugf("WAIT: Sleep Interval [%v]", sleepInterval/time.Second)
	logger.Debugf("WAIT: Sleep Timeout [%v]", sleepTimeout/time.Second)
	for {
		propStat, err := Client(meta).GetDomainStatus(ctx, gtm.GetDomainStatusRequest{
			DomainName: domain,
		})
		if err != nil {
			return false, fmt.Errorf("GetDomainStatus error: %s", err.Error())
		}
		logger.Debugf("WAIT: propStat.PropagationStatus [%v]", propStat.PropagationStatus)
		switch propStat.PropagationStatus {
		case "COMPLETE":
			logger.Debugf("WAIT: Return COMPLETE")
			return true, nil
		case "DENIED":
			logger.Debugf("WAIT: Return DENIED")
			return false, errors.New(propStat.Message)
		case "PENDING":
			if sleepTimeout <= 0 {
				logger.Debugf("WAIT: Return TIMED OUT")
				return false, nil
			}
			time.Sleep(sleepInterval)
			sleepTimeout -= sleepInterval
			logger.Debugf("WAIT: Sleep Time Remaining [%v]", sleepTimeout/time.Second)
		default:
			return false, fmt.Errorf("unknown propagationStatus while waiting for change completion") // don't know how/why we would have broken out.
		}
	}
}
