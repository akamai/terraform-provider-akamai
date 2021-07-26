package cps

import (
	"context"
	"errors"
	"fmt"
	"strconv"
	"strings"
	"time"

	"github.com/akamai/AkamaiOPEN-edgegrid-golang/v2/pkg/cps"
	"github.com/akamai/AkamaiOPEN-edgegrid-golang/v2/pkg/session"
	"github.com/akamai/terraform-provider-akamai/v2/pkg/akamai"
	cpstools "github.com/akamai/terraform-provider-akamai/v2/pkg/providers/cps/tools"
	"github.com/akamai/terraform-provider-akamai/v2/pkg/tools"
	"github.com/apex/log"
	"github.com/hashicorp/terraform-plugin-sdk/v2/diag"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/customdiff"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/schema"
)

var (
	// PollForChangeStatusInterval defines retry interval for getting status of a pending change
	PollForChangeStatusInterval = 10 * time.Second
)

const (
	statusCoordinateDomainValidation = "coodinate-domain-validation"
	statusVerificationWarnings       = "wait-review-pre-verification-safety-checks"
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
				Type:     schema.TypeString,
				Required: true,
				ForceNew: true,
			},
			"sans": {
				Type:     schema.TypeSet,
				Elem:     &schema.Schema{Type: schema.TypeString},
				Optional: true,
			},
			"secure_network": {
				Type:     schema.TypeString,
				Required: true,
				ForceNew: true,
			},
			"sni_only": {
				Type:     schema.TypeBool,
				Required: true,
				ForceNew: true,
			},
			"acknowledge_pre_verification_warnings": {
				Type:     schema.TypeBool,
				Optional: true,
			},
			"admin_contact": {
				Type:     schema.TypeSet,
				Required: true,
				MinItems: 1,
				MaxItems: 1,
				Elem:     contact,
			},
			"certificate_chain_type": {
				Type:     schema.TypeString,
				Default:  "default",
				Optional: true,
			},
			"csr": {
				Type:     schema.TypeSet,
				Required: true,
				MinItems: 1,
				MaxItems: 1,
				Elem: &schema.Resource{
					Schema: map[string]*schema.Schema{
						"country_code": {
							Type:     schema.TypeString,
							Required: true,
						},
						"city": {
							Type:     schema.TypeString,
							Required: true,
						},
						"organization": {
							Type:     schema.TypeString,
							Required: true,
						},
						"organizational_unit": {
							Type:     schema.TypeString,
							Required: true,
						},
						"state": {
							Type:     schema.TypeString,
							Required: true,
						},
					},
				},
			},
			"enable_multi_stacked_certificates": {
				Type:     schema.TypeBool,
				Optional: true,
			},
			"network_configuration": {
				Type:     schema.TypeSet,
				Required: true,
				MinItems: 1,
				MaxItems: 1,
				Elem: &schema.Resource{
					Schema: map[string]*schema.Schema{
						"client_mutual_authentication": {
							Type:     schema.TypeSet,
							Optional: true,
							MinItems: 1,
							MaxItems: 1,
							Elem: &schema.Resource{
								Schema: map[string]*schema.Schema{
									"send_ca_list_to_client": {
										Type:     schema.TypeBool,
										Optional: true,
									},
									"ocsp_enabled": {
										Type:     schema.TypeBool,
										Optional: true,
									},
									"set_id": {
										Type:     schema.TypeString,
										Optional: true,
									},
								},
							},
						},
						"disallowed_tls_versions": {
							Type:     schema.TypeSet,
							Optional: true,
							Elem:     &schema.Schema{Type: schema.TypeString},
						},
						"clone_dns_names": {
							Type:     schema.TypeBool,
							Optional: true,
						},
						"geography": {
							Type:     schema.TypeString,
							Required: true,
						},
						"must_have_ciphers": {
							Type:     schema.TypeString,
							Optional: true,
						},
						"ocsp_stapling": {
							Type:     schema.TypeString,
							Optional: true,
						},
						"preferred_ciphers": {
							Type:     schema.TypeString,
							Optional: true,
						},
						"quic_enabled": {
							Type:     schema.TypeBool,
							Optional: true,
						},
					},
				},
			},
			"signature_algorithm": {
				Type:     schema.TypeString,
				Required: true,
			},
			"tech_contact": {
				Type:     schema.TypeSet,
				Required: true,
				MinItems: 1,
				MaxItems: 1,
				Elem:     contact,
			},
			"organization": {
				Type:     schema.TypeSet,
				Required: true,
				MinItems: 1,
				MaxItems: 1,
				Elem: &schema.Resource{
					Schema: map[string]*schema.Schema{
						"name": {
							Type:     schema.TypeString,
							Required: true,
						},
						"phone": {
							Type:     schema.TypeString,
							Required: true,
						},
						"address_line_one": {
							Type:     schema.TypeString,
							Required: true,
						},
						"address_line_two": {
							Type:     schema.TypeString,
							Optional: true,
						},
						"city": {
							Type:     schema.TypeString,
							Required: true,
						},
						"region": {
							Type:     schema.TypeString,
							Required: true,
						},
						"postal_code": {
							Type:     schema.TypeString,
							Required: true,
						},
						"country_code": {
							Type:     schema.TypeString,
							Required: true,
						},
					},
				},
			},
			"contract_id": {
				Type:             schema.TypeString,
				ForceNew:         true,
				Required:         true,
				DiffSuppressFunc: diffSuppressContractID,
			},
			"certificate_type": {
				Type:     schema.TypeString,
				Computed: true,
			},
			"validation_type": {
				Type:     schema.TypeString,
				Computed: true,
			},
			"registration_authority": {
				Type:     schema.TypeString,
				Computed: true,
			},
			"dns_challenges": {
				Type:     schema.TypeSet,
				Computed: true,
				Elem: &schema.Resource{
					Schema: map[string]*schema.Schema{
						"domain": {
							Type:     schema.TypeString,
							Computed: true,
						},
						"full_path": {
							Type:     schema.TypeString,
							Computed: true,
						},
						"response_body": {
							Type:     schema.TypeString,
							Computed: true,
						},
					},
				},
				Set: cpstools.HashFromChallengesMap,
			},
			"http_challenges": {
				Type:     schema.TypeSet,
				Computed: true,
				Elem: &schema.Resource{
					Schema: map[string]*schema.Schema{
						"domain": {
							Type:     schema.TypeString,
							Computed: true,
						},
						"full_path": {
							Type:     schema.TypeString,
							Computed: true,
						},
						"response_body": {
							Type:     schema.TypeString,
							Computed: true,
						},
					},
				},
				Set: cpstools.HashFromChallengesMap,
			},
		},
		CustomizeDiff: customdiff.Sequence(
			func(ctx context.Context, diff *schema.ResourceDiff, i interface{}) error {
				if !diff.HasChange("sans") {
					return nil
				}
				domainsToValidate := []interface{}{map[string]interface{}{"domain": diff.Get("common_name").(string)}}
				if sans, ok := diff.Get("sans").(*schema.Set); ok {
					for _, san := range sans.List() {
						domain := map[string]interface{}{"domain": san.(string)}
						domainsToValidate = append(domainsToValidate, domain)
					}
				}
				if err := diff.SetNew("http_challenges", schema.NewSet(cpstools.HashFromChallengesMap, domainsToValidate)); err != nil {
					return fmt.Errorf("%w: %s", tools.ErrValueSet, err.Error())
				}
				if err := diff.SetNew("dns_challenges", schema.NewSet(cpstools.HashFromChallengesMap, domainsToValidate)); err != nil {
					return fmt.Errorf("%w: %s", tools.ErrValueSet, err.Error())
				}
				return nil
			}),
	}
}

var contact = &schema.Resource{
	Schema: map[string]*schema.Schema{
		"first_name": {
			Type:     schema.TypeString,
			Required: true,
		},
		"last_name": {
			Type:     schema.TypeString,
			Required: true,
		},
		"title": {
			Type:     schema.TypeString,
			Optional: true,
		},
		"organization": {
			Type:     schema.TypeString,
			Required: true,
		},
		"email": {
			Type:     schema.TypeString,
			Required: true,
		},
		"phone": {
			Type:     schema.TypeString,
			Required: true,
		},
		"address_line_one": {
			Type:     schema.TypeString,
			Required: true,
		},
		"address_line_two": {
			Type:     schema.TypeString,
			Optional: true,
		},
		"city": {
			Type:     schema.TypeString,
			Required: true,
		},
		"region": {
			Type:     schema.TypeString,
			Required: true,
		},
		"postal_code": {
			Type:     schema.TypeString,
			Required: true,
		},
		"country_code": {
			Type:     schema.TypeString,
			Required: true,
		},
	},
}

func resourceCPSDVEnrollmentCreate(ctx context.Context, d *schema.ResourceData, m interface{}) diag.Diagnostics {
	meta := akamai.Meta(m)
	logger := meta.Log("CPS", "resourceDVEnrollment")
	// create a context with logging for api calls
	ctx = session.ContextWithOptions(
		ctx,
		session.WithContextLog(logger),
	)
	client := inst.Client(meta)
	logger.Debug("Creating enrollment")

	enrollment := cps.Enrollment{
		CertificateType: "san",
		ValidationType:  "dv",
		RA:              "lets-encrypt",
	}
	if err := d.Set("certificate_type", enrollment.CertificateType); err != nil {
		return diag.Errorf("%v: %s", tools.ErrValueSet, err.Error())
	}
	if err := d.Set("validation_type", enrollment.ValidationType); err != nil {
		return diag.Errorf("%v: %s", tools.ErrValueSet, err.Error())
	}
	if err := d.Set("registration_authority", enrollment.RA); err != nil {
		return diag.Errorf("%v: %s", tools.ErrValueSet, err.Error())
	}
	adminContactSet, err := tools.GetSetValue("admin_contact", d)
	if err != nil {
		return diag.FromErr(err)
	}
	adminContact, err := cpstools.GetContactInfo(adminContactSet)
	if err != nil {
		return diag.FromErr(fmt.Errorf("'admin_contact' - %s", err))
	}
	enrollment.AdminContact = adminContact
	techContactSet, err := tools.GetSetValue("tech_contact", d)
	if err != nil {
		return diag.FromErr(err)
	}
	techContact, err := cpstools.GetContactInfo(techContactSet)
	if err != nil {
		return diag.FromErr(fmt.Errorf("'tech_contact' - %s", err))
	}
	enrollment.TechContact = techContact

	certificateChainType, err := tools.GetStringValue("certificate_chain_type", d)
	if err != nil && !errors.Is(err, tools.ErrNotFound) {
		return diag.FromErr(err)
	}
	enrollment.CertificateChainType = certificateChainType

	csr, err := cpstools.GetCSR(d)
	if err != nil {
		return diag.FromErr(err)
	}
	enrollment.CSR = csr

	enableMultiStacked, err := tools.GetBoolValue("enable_multi_stacked_certificates", d)
	if err != nil && !errors.Is(err, tools.ErrNotFound) {
		return diag.FromErr(err)
	}
	enrollment.EnableMultiStackedCertificates = enableMultiStacked

	networkConfig, err := cpstools.GetNetworkConfig(d)
	if err != nil {
		return diag.FromErr(err)
	}
	enrollment.NetworkConfiguration = networkConfig
	signatureAlgorithm, err := tools.GetStringValue("signature_algorithm", d)
	if err != nil && !errors.Is(err, tools.ErrNotFound) {
		return diag.FromErr(err)
	}
	enrollment.SignatureAlgorithm = signatureAlgorithm

	organization, err := cpstools.GetOrg(d)
	if err != nil {
		return diag.FromErr(err)
	}
	enrollment.Org = organization

	contractID, err := tools.GetStringValue("contract_id", d)
	if err != nil {
		return diag.FromErr(err)
	}

	req := cps.CreateEnrollmentRequest{
		Enrollment: enrollment,
		ContractID: strings.TrimPrefix(contractID, "ctr_"),
	}
	res, err := client.CreateEnrollment(ctx, req)
	if err != nil {
		return diag.FromErr(err)
	}
	d.SetId(strconv.Itoa(res.ID))

	acknowledgeWarnings, err := tools.GetBoolValue("acknowledge_pre_verification_warnings", d)
	if err != nil && !errors.Is(err, tools.ErrNotFound) {
		return diag.FromErr(err)
	}
	err = waitForVerification(ctx, logger, client, res.ID, acknowledgeWarnings)
	if err != nil {
		return diag.FromErr(err)
	}
	return resourceCPSDVEnrollmentRead(ctx, d, m)
}

func resourceCPSDVEnrollmentRead(ctx context.Context, d *schema.ResourceData, m interface{}) diag.Diagnostics {
	meta := akamai.Meta(m)
	logger := meta.Log("CPS", "resourceDVEnrollment")
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
	attrs := make(map[string]interface{})
	adminContact := cpstools.ContactInfoToMap(*enrollment.AdminContact)
	attrs["common_name"] = enrollment.CSR.CN
	sans := make([]string, 0)
	for _, san := range enrollment.CSR.SANS {
		if san == enrollment.CSR.CN {
			continue
		}
		sans = append(sans, san)
	}
	attrs["sans"] = sans
	attrs["sni_only"] = enrollment.NetworkConfiguration.SNIOnly
	attrs["secure_network"] = enrollment.NetworkConfiguration.SecureNetwork
	attrs["admin_contact"] = []interface{}{adminContact}
	techContact := cpstools.ContactInfoToMap(*enrollment.TechContact)
	attrs["tech_contact"] = []interface{}{techContact}
	attrs["certificate_chain_type"] = enrollment.CertificateChainType
	csr := cpstools.CSRToMap(*enrollment.CSR)
	attrs["csr"] = []interface{}{csr}
	attrs["enable_multi_stacked_certificates"] = enrollment.EnableMultiStackedCertificates
	networkConfig := cpstools.NetworkConfigToMap(*enrollment.NetworkConfiguration)
	attrs["network_configuration"] = []interface{}{networkConfig}
	attrs["signature_algorithm"] = enrollment.SignatureAlgorithm
	org := cpstools.OrgToMap(*enrollment.Org)
	attrs["organization"] = []interface{}{org}
	attrs["certificate_type"] = enrollment.CertificateType
	attrs["validation_type"] = enrollment.ValidationType

	err = tools.SetAttrs(d, attrs)
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
				return diag.Errorf("%v: %s", tools.ErrValueSet, err.Error())
			}
			if err := d.Set("dns_challenges", schema.NewSet(cpstools.HashFromChallengesMap, dnsChallenges)); err != nil {
				return diag.Errorf("%v: %s", tools.ErrValueSet, err.Error())
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
			return diag.Errorf("%v: %s", tools.ErrValueSet, err.Error())
		}
		if err := d.Set("dns_challenges", schema.NewSet(cpstools.HashFromChallengesMap, dnsChallenges)); err != nil {
			return diag.Errorf("%v: %s", tools.ErrValueSet, err.Error())
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
		return diag.Errorf("%v: %s", tools.ErrValueSet, err.Error())
	}
	if err := d.Set("dns_challenges", schema.NewSet(cpstools.HashFromChallengesMap, dnsChallenges)); err != nil {
		return diag.Errorf("%v: %s", tools.ErrValueSet, err.Error())
	}
	return nil
}

func resourceCPSDVEnrollmentUpdate(ctx context.Context, d *schema.ResourceData, m interface{}) diag.Diagnostics {
	meta := akamai.Meta(m)
	logger := meta.Log("CPS", "resourceDVEnrollment")
	ctx = session.ContextWithOptions(
		ctx,
		session.WithContextLog(logger),
	)
	client := inst.Client(meta)
	logger.Debug("Updating enrollment")

	acknowledgeWarnings, err := tools.GetBoolValue("acknowledge_pre_verification_warnings", d)
	if err != nil && !errors.Is(err, tools.ErrNotFound) {
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
		"enable_multi_stacked_certificates",
		"network_configuration",
		"signature_algorithm",
		"organization",
	) {
		logger.Debug("Enrollment does not have to be updated. Verifying status.")
		if err = waitForVerification(ctx, logger, client, enrollmentID, acknowledgeWarnings); err != nil {
			return diag.FromErr(err)
		}
		return resourceCPSDVEnrollmentRead(ctx, d, m)
	}
	enrollment := cps.Enrollment{
		CertificateType: "san",
		ValidationType:  "dv",
		RA:              "lets-encrypt",
	}
	if err := d.Set("certificate_type", "san"); err != nil {
		return diag.Errorf("%v: %s", tools.ErrValueSet, err.Error())
	}
	if err := d.Set("validation_type", "dv"); err != nil {
		return diag.Errorf("%v: %s", tools.ErrValueSet, err.Error())
	}

	adminContactSet, err := tools.GetSetValue("admin_contact", d)
	if err != nil {
		return diag.FromErr(err)
	}
	adminContact, err := cpstools.GetContactInfo(adminContactSet)
	if err != nil {
		return diag.FromErr(fmt.Errorf("'admin_contact' - %s", err))
	}
	enrollment.AdminContact = adminContact
	techContactSet, err := tools.GetSetValue("tech_contact", d)
	if err != nil {
		return diag.FromErr(err)
	}
	techContact, err := cpstools.GetContactInfo(techContactSet)
	if err != nil {
		return diag.FromErr(fmt.Errorf("'tech_contact' - %s", err))
	}
	enrollment.TechContact = techContact

	certificateChainType, err := tools.GetStringValue("certificate_chain_type", d)
	if err != nil && !errors.Is(err, tools.ErrNotFound) {
		return diag.FromErr(err)
	}
	enrollment.CertificateChainType = certificateChainType

	csr, err := cpstools.GetCSR(d)
	if err != nil {
		return diag.FromErr(err)
	}
	enrollment.CSR = csr

	enableMultiStacked, err := tools.GetBoolValue("enable_multi_stacked_certificates", d)
	if err != nil && !errors.Is(err, tools.ErrNotFound) {
		return diag.FromErr(err)
	}
	enrollment.EnableMultiStackedCertificates = enableMultiStacked

	networkConfig, err := cpstools.GetNetworkConfig(d)
	if err != nil {
		return diag.FromErr(err)
	}
	enrollment.NetworkConfiguration = networkConfig
	signatureAlgorithm, err := tools.GetStringValue("signature_algorithm", d)
	if err != nil && !errors.Is(err, tools.ErrNotFound) {
		return diag.FromErr(err)
	}
	enrollment.SignatureAlgorithm = signatureAlgorithm

	organization, err := cpstools.GetOrg(d)
	if err != nil {
		return diag.FromErr(err)
	}
	enrollment.Org = organization

	allowCancel := true
	req := cps.UpdateEnrollmentRequest{
		Enrollment:                enrollment,
		EnrollmentID:              enrollmentID,
		AllowCancelPendingChanges: &allowCancel,
	}

	if _, err := client.UpdateEnrollment(ctx, req); err != nil {
		return diag.FromErr(err)
	}
	d.SetId(strconv.Itoa(enrollmentID))

	if err = waitForVerification(ctx, logger, client, enrollmentID, acknowledgeWarnings); err != nil {
		return diag.FromErr(err)
	}
	return resourceCPSDVEnrollmentRead(ctx, d, m)
}

func resourceCPSDVEnrollmentDelete(ctx context.Context, d *schema.ResourceData, m interface{}) diag.Diagnostics {
	meta := akamai.Meta(m)
	logger := meta.Log("CPS", "resourceDVEnrollment")
	// create a context with logging for api calls
	ctx = session.ContextWithOptions(
		ctx,
		session.WithContextLog(logger),
	)
	client := inst.Client(meta)
	logger.Debug("Deleting enrollment")
	enrollmentID, err := strconv.Atoi(d.Id())
	if err != nil {
		return diag.FromErr(err)
	}
	allowCancelPendingChanges := true
	req := cps.RemoveEnrollmentRequest{
		EnrollmentID:              enrollmentID,
		AllowCancelPendingChanges: &allowCancelPendingChanges,
	}
	if _, err = client.RemoveEnrollment(ctx, req); err != nil {
		return diag.FromErr(err)
	}
	d.SetId("")
	return nil
}

func waitForVerification(ctx context.Context, logger log.Interface, client cps.CPS, enrollmentID int, acknowledgeWarnings bool) error {
	getEnrollmentReq := cps.GetEnrollmentRequest{EnrollmentID: enrollmentID}
	enrollmentGet, err := client.GetEnrollment(ctx, getEnrollmentReq)
	if err != nil {
		return err
	}
	changeID, err := cpstools.GetChangeIDFromPendingChanges(enrollmentGet.PendingChanges)
	if err != nil {
		if errors.Is(err, cpstools.ErrNoPendingChanges) {
			logger.Debug("No pending changes found on the enrollment")
			return nil
		}
		return err
	}

	changeStatusReq := cps.GetChangeStatusRequest{
		EnrollmentID: enrollmentID,
		ChangeID:     changeID,
	}
	status, err := client.GetChangeStatus(ctx, changeStatusReq)
	if err != nil {
		return err
	}
	for status.StatusInfo.Status != statusCoordinateDomainValidation && status.StatusInfo.Status != "complete" {
		select {
		case <-time.After(PollForChangeStatusInterval):
			status, err = client.GetChangeStatus(ctx, changeStatusReq)
			if err != nil {
				return err
			}
			if status.StatusInfo != nil && status.StatusInfo.Status == statusVerificationWarnings {
				warnings, err := client.GetChangePreVerificationWarnings(ctx, cps.GetChangeRequest{
					EnrollmentID: enrollmentID,
					ChangeID:     changeID,
				})
				if err != nil {
					return err
				}
				logger.Debugf("Pre-verification warnings: %s", warnings.Warnings)
				if acknowledgeWarnings {
					err = client.AcknowledgePreVerificationWarnings(ctx, cps.AcknowledgementRequest{
						Acknowledgement: cps.Acknowledgement{Acknowledgement: cps.AcknowledgementAcknowledge},
						EnrollmentID:    enrollmentID,
						ChangeID:        changeID,
					})
					if err != nil {
						return err
					}
					continue
				}
				return fmt.Errorf("enrollment pre-verification returned warnings and the enrollment cannot be validated. Please fix the issues or set acknowledge_pre_validation_warnings flag to true then run 'terraform apply' again: %s",
					warnings.Warnings)
			}
			log.Debugf("Change status: %s", status.StatusInfo.Status)
			if status.StatusInfo != nil && status.StatusInfo.Error != nil && status.StatusInfo.Error.Description != "" {
				return fmt.Errorf(status.StatusInfo.Error.Description)
			}
		case <-ctx.Done():
			return fmt.Errorf("change status context terminated: %w", ctx.Err())
		}
	}
	return nil
}

func resourceCPSDVEnrollmentImport(ctx context.Context, d *schema.ResourceData, m interface{}) ([]*schema.ResourceData, error) {
	meta := akamai.Meta(m)
	logger := meta.Log("CPS", "resourceDVEnrollment")
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
	if _, err := strconv.Atoi(enrollmentID); err != nil {
		return nil, fmt.Errorf("enrollment ID must be a number: %s", err)
	}
	if err := d.Set("contract_id", contractID); err != nil {
		return nil, fmt.Errorf("%v: %s", tools.ErrValueSet, err.Error())
	}
	d.SetId(enrollmentID)
	return []*schema.ResourceData{d}, nil
}

func diffSuppressContractID(_, old, new string, _ *schema.ResourceData) bool {
	trimPrefixFromOld := strings.TrimPrefix(old, "ctr_")
	trimPrefixFromNew := strings.TrimPrefix(new, "ctr_")

	if trimPrefixFromOld == trimPrefixFromNew {
		return true
	}
	return false
}
