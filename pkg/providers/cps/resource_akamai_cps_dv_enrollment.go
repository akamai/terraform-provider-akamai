package cps

import (
	"context"
	"errors"
	"fmt"
	"strconv"
	"strings"
	"time"

	"github.com/akamai/AkamaiOPEN-edgegrid-golang/v12/pkg/cps"
	"github.com/akamai/AkamaiOPEN-edgegrid-golang/v12/pkg/session"
	"github.com/akamai/terraform-provider-akamai/v8/pkg/common/ptr"
	"github.com/akamai/terraform-provider-akamai/v8/pkg/common/tf"
	"github.com/akamai/terraform-provider-akamai/v8/pkg/common/timeouts"
	"github.com/akamai/terraform-provider-akamai/v8/pkg/meta"
	cpstools "github.com/akamai/terraform-provider-akamai/v8/pkg/providers/cps/tools"
	"github.com/hashicorp/terraform-plugin-sdk/v2/diag"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/customdiff"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/schema"
)

var (
	// PollForChangeStatusInterval defines retry interval for getting status of a pending change.
	PollForChangeStatusInterval = 10 * time.Second
	// PollForGetEnrollmentInterval defines retry interval for getting enrollment.
	PollForGetEnrollmentInterval = 30 * time.Second
)

func resourceCPSDVEnrollment() *schema.Resource {
	return &schema.Resource{
		CreateContext: resourceCPSDVEnrollmentCreate,
		ReadContext:   resourceCPSDVEnrollmentRead,
		UpdateContext: resourceCPSDVEnrollmentUpdate,
		DeleteContext: resourceCPSDVEnrollmentDelete,
		Importer: &schema.ResourceImporter{
			StateContext: resourceCPSDVEnrollmentImport,
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
				Type:        schema.TypeSet,
				Required:    true,
				MinItems:    1,
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
			func(_ context.Context, diff *schema.ResourceDiff, _ interface{}) error {
				if !diff.HasChange("sans") {
					return nil
				}
				domainsToValidate := []interface{}{map[string]interface{}{
					"domain": strings.ToLower(diff.Get("common_name").(string)),
				}}
				if sans, ok := diff.Get("sans").(*schema.Set); ok {
					for _, san := range sans.List() {
						domain := map[string]interface{}{"domain": strings.ToLower(san.(string))}
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
			}),
		Timeouts: &schema.ResourceTimeout{
			Default: &DefaultEnrollmentTimeout,
		},
	}
}

func resourceCPSDVEnrollmentCreate(ctx context.Context, d *schema.ResourceData, m interface{}) diag.Diagnostics {
	meta := meta.Must(m)
	logger := meta.Log("CPS", "resourceCPSDVEnrollmentCreate")
	// create a context with logging for api calls
	ctx = session.ContextWithOptions(
		ctx,
		session.WithContextLog(logger),
	)
	client := inst.Client(meta)
	logger.Debug("Creating enrollment")

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
	if err = waitForVerification(ctx, logger, client, res.ID, acknowledgeWarnings, nil); err != nil {
		return diag.FromErr(err)
	}
	return resourceCPSDVEnrollmentRead(ctx, d, m)
}

func resourceCPSDVEnrollmentRead(ctx context.Context, d *schema.ResourceData, m interface{}) diag.Diagnostics {
	meta := meta.Must(m)
	logger := meta.Log("CPS", "resourceCPSDVEnrollmentRead")
	// create a context with logging for api calls
	ctx = session.ContextWithOptions(
		ctx,
		session.WithContextLog(logger),
	)
	client := inst.Client(meta)
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
	dnsChallenges := make([]interface{}, 0)
	httpChallenges := make([]interface{}, 0)
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
				httpChallenges = append(httpChallenges, map[string]interface{}{
					"full_path":     challenge.FullPath,
					"response_body": challenge.ResponseBody,
					"domain":        dv.Domain,
				})
			}
			if challenge.Type == "dns-01" {
				dnsChallenges = append(dnsChallenges, map[string]interface{}{
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

func resourceCPSDVEnrollmentUpdate(ctx context.Context, d *schema.ResourceData, m interface{}) diag.Diagnostics {
	meta := meta.Must(m)
	logger := meta.Log("CPS", "resourceCPSDVEnrollmentUpdate")
	ctx = session.ContextWithOptions(
		ctx,
		session.WithContextLog(logger),
	)
	client := inst.Client(meta)
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
		if err = waitForVerification(ctx, logger, client, enrollmentID, acknowledgeWarnings, nil); err != nil {
			return diag.FromErr(err)
		}
		return resourceCPSDVEnrollmentRead(ctx, d, m)
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

	if err = waitForVerification(ctx, logger, client, enrollmentID, acknowledgeWarnings, nil); err != nil {
		return diag.FromErr(err)
	}
	return resourceCPSDVEnrollmentRead(ctx, d, m)
}

func resourceCPSDVEnrollmentDelete(ctx context.Context, d *schema.ResourceData, m interface{}) diag.Diagnostics {
	return enrollmentDelete(ctx, d, m, "resourceCPSDVEnrollmentDelete")
}

func resourceCPSDVEnrollmentImport(ctx context.Context, d *schema.ResourceData, m interface{}) ([]*schema.ResourceData, error) {
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

	client := inst.Client(meta)
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
