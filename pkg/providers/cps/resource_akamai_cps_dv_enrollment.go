package cps

import (
	"context"
	"errors"
	"fmt"
	"strconv"
	"strings"
	"time"

	"github.com/akamai/AkamaiOPEN-edgegrid-golang/v13/pkg/cps"
	"github.com/akamai/AkamaiOPEN-edgegrid-golang/v13/pkg/session"
	"github.com/akamai/terraform-provider-akamai/v10/pkg/common/ptr"
	"github.com/akamai/terraform-provider-akamai/v10/pkg/common/tf"
	"github.com/akamai/terraform-provider-akamai/v10/pkg/common/timeouts"
	"github.com/akamai/terraform-provider-akamai/v10/pkg/meta"
	cpstools "github.com/akamai/terraform-provider-akamai/v10/pkg/providers/cps/tools"
	"github.com/hashicorp/go-cty/cty"
	"github.com/hashicorp/terraform-plugin-sdk/v2/diag"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/customdiff"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/schema"
)

// dvEnrollmentResource represents the akamai_cps_dv_enrollment resource with configurable polling intervals.
type dvEnrollmentResource struct {
	pollChangeStatusInterval  time.Duration
	pollGetEnrollmentInterval time.Duration
}

func resourceCPSDVEnrollment(pollChangeStatusInterval, pollGetEnrollmentInterval time.Duration) *schema.Resource {
	res := &dvEnrollmentResource{
		pollChangeStatusInterval:  pollChangeStatusInterval,
		pollGetEnrollmentInterval: pollGetEnrollmentInterval,
	}
	return &schema.Resource{
		CreateContext: res.create,
		ReadContext:   res.read,
		UpdateContext: res.update,
		DeleteContext: res.delete,
		Importer: &schema.ResourceImporter{
			StateContext: res.importState,
		},
		Schema: map[string]*schema.Schema{
			"common_name": {
				Type:        schema.TypeString,
				Required:    true,
				ForceNew:    true,
				Description: "Common name used for enrollment",
			},
			"allow_duplicate_common_name": {
				Type:        schema.TypeBool,
				Optional:    true,
				Default:     false,
				Description: "Allow to duplicate common name. Default is false",
			},
			"sans": {
				Type:        schema.TypeSet,
				Optional:    true,
				Elem:        &schema.Schema{Type: schema.TypeString},
				Description: "List of SANs",
			},
			"secure_network": {
				Type:        schema.TypeString,
				Required:    true,
				ForceNew:    true,
				Description: "Type of TLS deployment network",
			},
			"sni_only": {
				Type:        schema.TypeBool,
				Required:    true,
				ForceNew:    true,
				Description: "Whether Server Name Indication is used for enrollment",
			},
			"acknowledge_pre_verification_warnings": {
				Type:        schema.TypeBool,
				Optional:    true,
				Default:     false,
				Description: "Whether acknowledge warnings before certificate verification. Default is false",
			},
			"admin_contact": {
				Type:        schema.TypeSet,
				Required:    true,
				MinItems:    1,
				MaxItems:    1,
				Elem:        contact,
				Description: "Contact information for the certificate administrator to use at organization",
			},
			"certificate_chain_type": {
				Type:        schema.TypeString,
				Default:     "default",
				Optional:    true,
				Description: "Certificate trust chain type. Default is 'default'",
			},
			"csr": {
				Type:        schema.TypeSet,
				Required:    true,
				MinItems:    1,
				MaxItems:    1,
				Elem:        csr,
				Description: "Certificate signing request generated during enrollment creation",
			},
			"network_configuration": {
				Type:        schema.TypeList,
				Optional:    true,
				Computed:    true,
				MaxItems:    1,
				Elem:        networkConfiguration,
				Description: "Settings containing network information and TLS Metadata used by CPS",
			},
			"signature_algorithm": {
				Type:        schema.TypeString,
				Required:    true,
				Description: "SHA algorithm type",
			},
			"tech_contact": {
				Type:        schema.TypeSet,
				Required:    true,
				MinItems:    1,
				MaxItems:    1,
				Elem:        contact,
				Description: "Contact information for an administrator at Akamai",
			},
			"organization": {
				Type:        schema.TypeSet,
				Required:    true,
				MinItems:    1,
				MaxItems:    1,
				Elem:        organization,
				Description: "Organization information",
			},
			"contract_id": {
				Type:             schema.TypeString,
				ForceNew:         true,
				Required:         true,
				DiffSuppressFunc: tf.FieldPrefixSuppress("ctr_"),
				Description:      "Contract ID for which enrollment is retrieved",
			},
			"certificate_type": {
				Type:        schema.TypeString,
				Computed:    true,
				Description: "Certificate type of enrollment",
			},
			"validation_type": {
				Type:        schema.TypeString,
				Computed:    true,
				Description: "Enrolment validation type",
			},
			"registration_authority": {
				Type:        schema.TypeString,
				Computed:    true,
				Description: "The registration authority or certificate authority (CA) used to obtain a certificate",
			},
			"dns_challenges": {
				Type:        schema.TypeSet,
				Computed:    true,
				Description: "DNS challenge information",
				Elem: &schema.Resource{
					Schema: map[string]*schema.Schema{
						"domain": {
							Type:        schema.TypeString,
							Computed:    true,
							Description: "Domain for which the challenges were completed",
						},
						"full_path": {
							Type:        schema.TypeString,
							Computed:    true,
							Description: "The domain name where Akamai publishes the response body to validate",
						},
						"response_body": {
							Type:        schema.TypeString,
							Computed:    true,
							Description: "The unique content of the challenge",
						},
					},
				},
				Set: cpstools.HashFromChallengesMap,
			},
			"http_challenges": {
				Type:        schema.TypeSet,
				Computed:    true,
				Description: "HTTP challenge information",
				Elem: &schema.Resource{
					Schema: map[string]*schema.Schema{
						"domain": {
							Type:        schema.TypeString,
							Computed:    true,
							Description: "Domain for which the challenges were completed",
						},
						"full_path": {
							Type:        schema.TypeString,
							Computed:    true,
							Description: "The URL where Akamai publishes the response body to validate",
						},
						"response_body": {
							Type:        schema.TypeString,
							Computed:    true,
							Description: "The unique content of the challenge",
						},
					},
				},
				Set: cpstools.HashFromChallengesMap,
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
		},
		CustomizeDiff: customdiff.Sequence(
			validateNetworkConfigurationPresent,
			setDefaultEnableForAllSANs,
			validateDNSNameSettingsConflict,
			updateChallengesForSANChange,
		),
		Timeouts: &schema.ResourceTimeout{
			Default: &DefaultEnrollmentTimeout,
		},
	}
}

func updateChallengesForSANChange(_ context.Context, diff *schema.ResourceDiff, _ any) error {
	if !diff.HasChange("sans") {
		return nil
	}
	domainsToValidate := []any{map[string]any{
		"domain": strings.ToLower(diff.Get("common_name").(string)),
	}}
	if sans, ok := diff.Get("sans").(*schema.Set); ok {
		for _, san := range sans.List() {
			domain := map[string]any{"domain": strings.ToLower(san.(string))}
			domainsToValidate = append(domainsToValidate, domain)
		}
	}
	if err := diff.SetNew("http_challenges", schema.NewSet(cpstools.HashFromChallengesMap, domainsToValidate)); err != nil {
		return fmt.Errorf("%w: %s", tf.ErrValueSet, err.Error())
	}
	if err := diff.SetNew("dns_challenges", schema.NewSet(cpstools.HashFromChallengesMap, domainsToValidate)); err != nil {
		return fmt.Errorf("%w: %s", tf.ErrValueSet, err.Error())
	}
	return nil
}

func (r *dvEnrollmentResource) create(ctx context.Context, d *schema.ResourceData, m any) diag.Diagnostics {
	meta := meta.Must(m)
	logger := meta.Log("CPS", "resourceCPSDVEnrollmentCreate")
	// create a context with logging for api calls
	ctx = session.ContextWithOptions(
		ctx,
		session.WithContextLog(logger),
	)
	client := meta.Client().GetCPS()
	logger.Debug("Creating enrollment")

	if err := validateResolvedDNSNameSettings(d); err != nil {
		return diag.FromErr(err)
	}

	enrollmentReqBody := cps.EnrollmentRequestBody{
		CertificateType: "san",
		ValidationType:  "dv",
		RA:              "lets-encrypt",
	}
	if err := d.Set("certificate_type", enrollmentReqBody.CertificateType); err != nil {
		return diag.Errorf("%v: %s", tf.ErrValueSet, err.Error())
	}
	if err := d.Set("validation_type", enrollmentReqBody.ValidationType); err != nil {
		return diag.Errorf("%v: %s", tf.ErrValueSet, err.Error())
	}
	if err := d.Set("registration_authority", enrollmentReqBody.RA); err != nil {
		return diag.Errorf("%v: %s", tf.ErrValueSet, err.Error())
	}
	adminContactSet, err := tf.GetSetValue("admin_contact", d)
	if err != nil {
		return diag.FromErr(err)
	}
	adminContact, err := cpstools.GetContactInfo(adminContactSet)
	if err != nil {
		return diag.Errorf("'admin_contact' - %s", err)
	}
	enrollmentReqBody.AdminContact = adminContact
	techContactSet, err := tf.GetSetValue("tech_contact", d)
	if err != nil {
		return diag.FromErr(err)
	}
	techContact, err := cpstools.GetContactInfo(techContactSet)
	if err != nil {
		return diag.Errorf("'tech_contact' - %s", err)
	}
	enrollmentReqBody.TechContact = techContact

	certificateChainType, err := tf.GetStringValue("certificate_chain_type", d)
	if err != nil && !errors.Is(err, tf.ErrNotFound) {
		return diag.FromErr(err)
	}
	enrollmentReqBody.CertificateChainType = certificateChainType

	csr, err := cpstools.GetCSR(d)
	if err != nil {
		return diag.FromErr(err)
	}
	enrollmentReqBody.CSR = csr

	// DV does not support multi stack certificates
	enrollmentReqBody.EnableMultiStackedCertificates = false

	networkConfig, err := cpstools.GetNetworkConfig(d)
	if err != nil {
		return diag.FromErr(err)
	}
	enrollmentReqBody.NetworkConfiguration = networkConfig
	signatureAlgorithm, err := tf.GetStringValue("signature_algorithm", d)
	if err != nil {
		return diag.FromErr(err)
	}
	enrollmentReqBody.SignatureAlgorithm = signatureAlgorithm

	organization, err := cpstools.GetOrg(d)
	if err != nil {
		return diag.FromErr(err)
	}
	enrollmentReqBody.Org = organization

	contractID, err := tf.GetStringValue("contract_id", d)
	if err != nil {
		return diag.FromErr(err)
	}
	allowDuplicateCN, err := tf.GetBoolValue("allow_duplicate_common_name", d)
	if err != nil {
		return diag.FromErr(err)
	}

	// save ClientMutualAuthentication and unset it in enrollment request struct
	// create request must not have it set; in case it's not nil, we will run update later to add it
	clientMutualAuthentication := enrollmentReqBody.NetworkConfiguration.ClientMutualAuthentication
	enrollmentReqBody.NetworkConfiguration.ClientMutualAuthentication = nil

	req := cps.CreateEnrollmentRequest{
		EnrollmentRequestBody: enrollmentReqBody,
		ContractID:            strings.TrimPrefix(contractID, "ctr_"),
		AllowDuplicateCN:      allowDuplicateCN,
	}
	res, err := client.CreateEnrollment(ctx, req)
	if err != nil {
		return diag.FromErr(err)
	}
	d.SetId(strconv.Itoa(res.ID))

	acknowledgeWarnings, err := tf.GetBoolValue("acknowledge_pre_verification_warnings", d)
	if err != nil && !errors.Is(err, tf.ErrNotFound) {
		return diag.FromErr(err)
	}

	// when clientMutualAuthentication was provided, insert it back to enrollment and send the update request
	if clientMutualAuthentication != nil {
		logger.Debug("Updating ClientMutualAuthentication configuration")
		enrollmentReqBody.NetworkConfiguration.ClientMutualAuthentication = clientMutualAuthentication
		req := cps.UpdateEnrollmentRequest{
			EnrollmentID:              res.ID,
			EnrollmentRequestBody:     enrollmentReqBody,
			AllowCancelPendingChanges: ptr.To(true),
		}
		_, err := client.UpdateEnrollment(ctx, req)
		if err != nil {
			return diag.FromErr(err)
		}
	}
	if err = waitForVerification(ctx, logger, client, res.ID, acknowledgeWarnings, nil, r.pollChangeStatusInterval); err != nil {
		return diag.FromErr(err)
	}
	return r.read(ctx, d, m)
}

func (r *dvEnrollmentResource) read(ctx context.Context, d *schema.ResourceData, m any) diag.Diagnostics {
	meta := meta.Must(m)
	logger := meta.Log("CPS", "resourceCPSDVEnrollmentRead")
	// create a context with logging for api calls
	ctx = session.ContextWithOptions(
		ctx,
		session.WithContextLog(logger),
	)
	client := meta.Client().GetCPS()
	logger.Debug("Reading enrollment")
	enrollmentID, err := strconv.Atoi(d.Id())
	if err != nil {
		return diag.FromErr(err)
	}
	req := cps.GetEnrollmentRequest{EnrollmentID: enrollmentID}
	enrollment, err := client.GetEnrollment(ctx, req)
	if err != nil {
		return diag.FromErr(err)
	}
	attrs, err := readAttrs(enrollment, d)
	if err != nil {
		return diag.FromErr(err)
	}
	attrs["certificate_type"] = enrollment.CertificateType
	attrs["validation_type"] = enrollment.ValidationType
	attrs["registration_authority"] = enrollment.RA

	err = tf.SetAttrs(d, attrs)
	if err != nil {
		return diag.FromErr(err)
	}
	dnsChallenges := make([]any, 0)
	httpChallenges := make([]any, 0)
	changeID, err := cpstools.GetChangeIDFromPendingChanges(enrollment.PendingChanges)
	if err != nil {
		if errors.Is(err, cpstools.ErrNoPendingChanges) {
			logger.Debugf("No pending changes found on the enrollment")
			if err := d.Set("http_challenges", httpChallenges); err != nil {
				return diag.Errorf("%v: %s", tf.ErrValueSet, err.Error())
			}
			if err := d.Set("dns_challenges", schema.NewSet(cpstools.HashFromChallengesMap, dnsChallenges)); err != nil {
				return diag.Errorf("%v: %s", tf.ErrValueSet, err.Error())
			}
			return nil
		}
		return diag.FromErr(err)
	}

	changeStatusReq := cps.GetChangeStatusRequest{
		EnrollmentID: enrollmentID,
		ChangeID:     changeID,
	}
	status, err := client.GetChangeStatus(ctx, changeStatusReq)
	if err != nil {
		return diag.FromErr(err)
	}
	if len(status.AllowedInput) < 1 || status.AllowedInput[0].Type != "lets-encrypt-challenges" {
		if err := d.Set("http_challenges", httpChallenges); err != nil {
			return diag.Errorf("%v: %s", tf.ErrValueSet, err.Error())
		}
		if err := d.Set("dns_challenges", schema.NewSet(cpstools.HashFromChallengesMap, dnsChallenges)); err != nil {
			return diag.Errorf("%v: %s", tf.ErrValueSet, err.Error())
		}
		return nil
	}
	getChallengesReq := cps.GetChangeRequest{
		EnrollmentID: enrollmentID,
		ChangeID:     changeID,
	}
	challenges, err := client.GetChangeLetsEncryptChallenges(ctx, getChallengesReq)
	if err != nil {
		return diag.FromErr(err)
	}
	for _, dv := range challenges.DV {
		if dv.ValidationStatus == "VALIDATED" {
			continue
		}
		for _, challenge := range dv.Challenges {
			if challenge.Status != "pending" {
				continue
			}
			if challenge.Type == "http-01" {
				httpChallenges = append(httpChallenges, map[string]any{
					"full_path":     challenge.FullPath,
					"response_body": challenge.ResponseBody,
					"domain":        dv.Domain,
				})
			}
			if challenge.Type == "dns-01" {
				dnsChallenges = append(dnsChallenges, map[string]any{
					"full_path":     challenge.FullPath,
					"response_body": challenge.ResponseBody,
					"domain":        dv.Domain,
				})
			}
		}
	}
	if err := d.Set("http_challenges", httpChallenges); err != nil {
		return diag.Errorf("%v: %s", tf.ErrValueSet, err.Error())
	}
	if err := d.Set("dns_challenges", schema.NewSet(cpstools.HashFromChallengesMap, dnsChallenges)); err != nil {
		return diag.Errorf("%v: %s", tf.ErrValueSet, err.Error())
	}
	return nil
}

func (r *dvEnrollmentResource) update(ctx context.Context, d *schema.ResourceData, m any) diag.Diagnostics {
	meta := meta.Must(m)
	logger := meta.Log("CPS", "resourceCPSDVEnrollmentUpdate")
	ctx = session.ContextWithOptions(
		ctx,
		session.WithContextLog(logger),
	)
	client := meta.Client().GetCPS()
	logger.Debug("Updating enrollment")

	if !d.HasChangeExcept("timeouts") {
		logger.Debug("Only timeouts were updated, skipping")
		return nil
	}

	acknowledgeWarnings, err := tf.GetBoolValue("acknowledge_pre_verification_warnings", d)
	if err != nil && !errors.Is(err, tf.ErrNotFound) {
		return diag.FromErr(err)
	}
	enrollmentID, err := strconv.Atoi(d.Id())
	if err != nil {
		return diag.FromErr(err)
	}
	if !d.HasChanges(
		"sans",
		"admin_contact",
		"tech_contact",
		"certificate_chain_type",
		"csr",
		"network_configuration",
		"signature_algorithm",
		"organization",
	) {
		logger.Debug("Enrollment does not have to be updated. Verifying status.")
		if err = waitForVerification(ctx, logger, client, enrollmentID, acknowledgeWarnings, nil, r.pollChangeStatusInterval); err != nil {
			return diag.FromErr(err)
		}
		return r.read(ctx, d, m)
	}
	if err := validateResolvedDNSNameSettings(d); err != nil {
		return diag.FromErr(err)
	}

	enrollmentReqBody := cps.EnrollmentRequestBody{
		CertificateType: "san",
		ValidationType:  "dv",
		RA:              "lets-encrypt",
	}
	if err := d.Set("certificate_type", "san"); err != nil {
		return diag.Errorf("%v: %s", tf.ErrValueSet, err.Error())
	}
	if err := d.Set("validation_type", "dv"); err != nil {
		return diag.Errorf("%v: %s", tf.ErrValueSet, err.Error())
	}
	if err := d.Set("registration_authority", "lets-encrypt"); err != nil {
		return diag.Errorf("%v: %s", tf.ErrValueSet, err.Error())
	}

	adminContactSet, err := tf.GetSetValue("admin_contact", d)
	if err != nil {
		return diag.FromErr(err)
	}
	adminContact, err := cpstools.GetContactInfo(adminContactSet)
	if err != nil {
		return diag.Errorf("'admin_contact' - %s", err)
	}
	enrollmentReqBody.AdminContact = adminContact
	techContactSet, err := tf.GetSetValue("tech_contact", d)
	if err != nil {
		return diag.FromErr(err)
	}
	techContact, err := cpstools.GetContactInfo(techContactSet)
	if err != nil {
		return diag.Errorf("'tech_contact' - %s", err)
	}
	enrollmentReqBody.TechContact = techContact

	certificateChainType, err := tf.GetStringValue("certificate_chain_type", d)
	if err != nil && !errors.Is(err, tf.ErrNotFound) {
		return diag.FromErr(err)
	}
	enrollmentReqBody.CertificateChainType = certificateChainType

	csr, err := cpstools.GetCSR(d)
	if err != nil {
		return diag.FromErr(err)
	}
	enrollmentReqBody.CSR = csr

	// DV does not support multi stack certificates
	enrollmentReqBody.EnableMultiStackedCertificates = false

	networkConfig, err := cpstools.GetNetworkConfig(d)
	if err != nil {
		return diag.FromErr(err)
	}
	enrollmentReqBody.NetworkConfiguration = networkConfig
	signatureAlgorithm, err := tf.GetStringValue("signature_algorithm", d)
	if err != nil && !errors.Is(err, tf.ErrNotFound) {
		return diag.FromErr(err)
	}
	enrollmentReqBody.SignatureAlgorithm = signatureAlgorithm

	organization, err := cpstools.GetOrg(d)
	if err != nil {
		return diag.FromErr(err)
	}
	enrollmentReqBody.Org = organization

	allowCancel := true
	req := cps.UpdateEnrollmentRequest{
		EnrollmentRequestBody:     enrollmentReqBody,
		EnrollmentID:              enrollmentID,
		AllowCancelPendingChanges: &allowCancel,
	}

	if _, err := client.UpdateEnrollment(ctx, req); err != nil {
		return diag.FromErr(err)
	}
	d.SetId(strconv.Itoa(enrollmentID))

	if err = waitForVerification(ctx, logger, client, enrollmentID, acknowledgeWarnings, nil, r.pollChangeStatusInterval); err != nil {
		return diag.FromErr(err)
	}
	return r.read(ctx, d, m)
}

func (r *dvEnrollmentResource) delete(ctx context.Context, d *schema.ResourceData, m any) diag.Diagnostics {
	return enrollmentDelete(ctx, d, m, "resourceCPSDVEnrollmentDelete", r.pollGetEnrollmentInterval)
}

func (r *dvEnrollmentResource) importState(ctx context.Context, d *schema.ResourceData, m any) ([]*schema.ResourceData, error) {
	meta := meta.Must(m)
	logger := meta.Log("CPS", "resourceCPSDVEnrollmentImport")
	// create a context with logging for api calls
	ctx = session.ContextWithOptions(
		ctx,
		session.WithContextLog(logger),
	)
	logger.Debug("Importing enrollment")
	parts := strings.Split(d.Id(), ",")
	if len(parts) != 2 {
		return nil, fmt.Errorf("import id has to be a comma separated list of enrollment id and contract id")
	}
	enrollmentID := parts[0]
	contractID := parts[1]
	if enrollmentID == "" || contractID == "" {
		return nil, fmt.Errorf("enrollment and contract IDs must have non empty values")
	}
	eid, err := strconv.Atoi(enrollmentID)
	if err != nil {
		return nil, fmt.Errorf("enrollment ID must be a number: %s", err)
	}

	client := meta.Client().GetCPS()
	req := cps.GetEnrollmentRequest{EnrollmentID: eid}
	enrollment, err := client.GetEnrollment(ctx, req)
	if err != nil {
		return nil, fmt.Errorf("unable to fetch enrollment: %s", err)
	}
	if enrollment.ValidationType != "dv" {
		return nil, fmt.Errorf("unable to import: wrong validation type: expected 'dv', got '%s'", enrollment.ValidationType)
	}

	if err := d.Set("allow_duplicate_common_name", false); err != nil {
		return nil, fmt.Errorf("%v: %s", tf.ErrValueSet, err.Error())
	}
	if err := d.Set("acknowledge_pre_verification_warnings", false); err != nil {
		return nil, fmt.Errorf("%v: %s", tf.ErrValueSet, err.Error())
	}
	if err := d.Set("contract_id", contractID); err != nil {
		return nil, fmt.Errorf("%v: %s", tf.ErrValueSet, err.Error())
	}
	d.SetId(enrollmentID)
	return []*schema.ResourceData{d}, nil
}

// setDefaultEnableForAllSANs sets enable_for_all_sans to true for SNI enrollments when neither it
// nor clone_dns_names is explicitly provided in the configuration, making the default visible in
// the plan.
func setDefaultEnableForAllSANs(_ context.Context, diff *schema.ResourceDiff, _ any) error {
	rawConfig := diff.GetRawConfig()
	if rawConfig == cty.NilVal || rawConfig.IsNull() || !rawConfig.IsKnown() {
		return nil
	}

	if !dnsNameSettingsDefaultingEnabled(rawConfig.GetAttr("sni_only")) {
		return nil
	}

	ncVal := rawConfig.GetAttr("network_configuration")
	if ncVal.IsNull() || !ncVal.IsKnown() || ncVal.LengthInt() == 0 {
		return nil
	}

	it := ncVal.ElementIterator()
	it.Next()
	_, elem := it.Element()
	if !elem.IsKnown() || elem.IsNull() {
		return nil
	}

	cloneDNSAttr := elem.GetAttr("clone_dns_names")
	enableForAllSANsAttr := elem.GetAttr("enable_for_all_sans")

	// If at least one is explicitly set by the user, do not override.
	if !cloneDNSAttr.IsNull() || !enableForAllSANsAttr.IsNull() {
		networkConfigList, ok := diff.Get("network_configuration").([]any)
		if !ok || len(networkConfigList) == 0 {
			return nil
		}
		networkConfig := networkConfigList[0].(map[string]any)
		if cloneDNSAttr.IsNull() && enableForAllSANsAttr.IsKnown() {
			networkConfig["clone_dns_names"] = enableForAllSANsAttr.True()
			networkConfig["enable_for_all_sans"] = enableForAllSANsAttr.True()
		} else if enableForAllSANsAttr.IsNull() && cloneDNSAttr.IsKnown() {
			networkConfig["clone_dns_names"] = cloneDNSAttr.True()
			networkConfig["enable_for_all_sans"] = cloneDNSAttr.True()
		} else {
			return nil
		}
		if err := diff.SetNew("network_configuration", []any{networkConfig}); err != nil {
			return fmt.Errorf("%w: %s", tf.ErrValueSet, err.Error())
		}
		return nil
	}

	// Neither is specified — apply enable_for_all_sans = true as the visible default.
	networkConfigList, ok := diff.Get("network_configuration").([]any)
	if !ok || len(networkConfigList) == 0 {
		return nil
	}
	networkConfig := networkConfigList[0].(map[string]any)
	networkConfig["enable_for_all_sans"] = true
	networkConfig["clone_dns_names"] = true
	if err := diff.SetNew("network_configuration", []any{networkConfig}); err != nil {
		return fmt.Errorf("%w: %s", tf.ErrValueSet, err.Error())
	}
	return nil
}

func dnsNameSettingsDefaultingEnabled(sniOnly cty.Value) bool {
	return !sniOnly.IsNull() && sniOnly.IsKnown() && sniOnly.True()
}

// validateNetworkConfigurationPresent returns an error when network_configuration is absent from the config.
func validateNetworkConfigurationPresent(_ context.Context, diff *schema.ResourceDiff, _ any) error {
	rawConfig := diff.GetRawConfig()
	if rawConfig == cty.NilVal || rawConfig.IsNull() || !rawConfig.IsKnown() {
		return nil
	}
	ncVal := rawConfig.GetAttr("network_configuration")
	if ncVal.IsNull() || ncVal.LengthInt() == 0 {
		return fmt.Errorf("'network_configuration' is required")
	}
	return nil
}

// validateDNSNameSettingsConflict rejects mutually exclusive DNS name settings.
func validateDNSNameSettingsConflict(_ context.Context, diff *schema.ResourceDiff, _ any) error {
	rawConfig := diff.GetRawConfig()
	if err := validateDNSNameSettings(rawConfig); err != nil {
		return err
	}

	return validateDNSNamesChangedForModeTransition(
		dnsNamesModeDisabled(diff, rawConfig, "clone_dns_names") ||
			dnsNamesModeDisabled(diff, rawConfig, "enable_for_all_sans"),
		diff.HasChange("network_configuration.0.dns_names"),
	)
}

func validateDNSNamesChangedForModeTransition(modeDisabled, dnsNamesChanged bool) error {
	if modeDisabled && !dnsNamesChanged {
		return fmt.Errorf("'dns_names' must change when switching 'enable_for_all_sans' or 'clone_dns_names' from true to false")
	}
	return nil
}

func dnsNamesModeDisabled(diff *schema.ResourceDiff, rawConfig cty.Value, attribute string) bool {
	if rawConfig == cty.NilVal || rawConfig.IsNull() || !rawConfig.IsKnown() {
		return false
	}

	networkConfiguration := rawConfig.GetAttr("network_configuration")
	if networkConfiguration.IsNull() || !networkConfiguration.IsKnown() || networkConfiguration.LengthInt() == 0 {
		return false
	}

	it := networkConfiguration.ElementIterator()
	it.Next()
	_, settings := it.Element()
	if settings.IsNull() || !settings.IsKnown() {
		return false
	}
	enabled := settings.GetAttr(attribute)
	return enabled.IsKnown() && !enabled.IsNull() && !enabled.True() &&
		diff.HasChange("network_configuration.0."+attribute)
}

func validateResolvedDNSNameSettings(d *schema.ResourceData) error {
	return validateDNSNameSettings(d.GetRawConfig())
}

func validateDNSNameSettings(rawConfig cty.Value) error {
	if rawConfig == cty.NilVal || rawConfig.IsNull() || !rawConfig.IsKnown() {
		return nil
	}

	sniOnly := rawConfig.GetAttr("sni_only")
	if sniOnly.IsKnown() && !sniOnly.IsNull() && !sniOnly.True() {
		if hasDNSNameSettings(rawConfig) {
			return fmt.Errorf("'enable_for_all_sans', 'clone_dns_names', and 'dns_names' cannot be provided when 'sni_only' is false")
		}
		return nil
	}
	if !dnsNameSettingsDefaultingEnabled(sniOnly) {
		return nil
	}

	ncVal := rawConfig.GetAttr("network_configuration")
	if ncVal.IsNull() || !ncVal.IsKnown() || ncVal.LengthInt() == 0 {
		return nil
	}

	// network_configuration has MaxItems:1; LengthInt() > 0 was confirmed above.
	it := ncVal.ElementIterator()
	it.Next()
	_, elem := it.Element()
	if !elem.IsKnown() || elem.IsNull() {
		return nil
	}

	cloneDNSAttr := elem.GetAttr("clone_dns_names")
	enableForAllSANsAttr := elem.GetAttr("enable_for_all_sans")
	dnsNamesAttr := elem.GetAttr("dns_names")

	if !cloneDNSAttr.IsNull() && !enableForAllSANsAttr.IsNull() {
		return fmt.Errorf("'enable_for_all_sans' and 'clone_dns_names' cannot both be provided at the same time")
	}
	if !dnsNamesAttr.IsNull() && !enableForAllSANsAttr.IsNull() && enableForAllSANsAttr.IsKnown() && enableForAllSANsAttr.True() {
		return fmt.Errorf("'dns_names' cannot be provided when 'enable_for_all_sans' is true")
	}
	if !dnsNamesAttr.IsNull() && !cloneDNSAttr.IsNull() && cloneDNSAttr.IsKnown() && cloneDNSAttr.True() {
		return fmt.Errorf("'dns_names' cannot be provided when 'clone_dns_names' is true")
	}
	if dnsNamesRequired(cloneDNSAttr, enableForAllSANsAttr) && (dnsNamesAttr.IsNull() || (dnsNamesAttr.IsKnown() && dnsNamesAttr.LengthInt() == 0)) {
		return fmt.Errorf("'dns_names' is required when 'enable_for_all_sans' or 'clone_dns_names' is false")
	}
	return nil
}

func hasDNSNameSettings(rawConfig cty.Value) bool {
	networkConfiguration := rawConfig.GetAttr("network_configuration")
	if networkConfiguration.IsNull() || !networkConfiguration.IsKnown() || networkConfiguration.LengthInt() == 0 {
		return false
	}

	it := networkConfiguration.ElementIterator()
	it.Next()
	_, settings := it.Element()
	if settings.IsNull() || !settings.IsKnown() {
		return false
	}

	for _, attribute := range []string{"enable_for_all_sans", "clone_dns_names", "dns_names"} {
		value := settings.GetAttr(attribute)
		if !value.IsNull() {
			return true
		}
	}
	return false
}

func dnsNamesRequired(cloneDNSAttr, enableForAllSANsAttr cty.Value) bool {
	return (!cloneDNSAttr.IsNull() && cloneDNSAttr.IsKnown() && !cloneDNSAttr.True()) ||
		(!enableForAllSANsAttr.IsNull() && enableForAllSANsAttr.IsKnown() && !enableForAllSANsAttr.True())
}
