package cps

import (
	"context"
	"errors"
	"fmt"
	"github.com/akamai/AkamaiOPEN-edgegrid-golang/v2/pkg/cps"
	"github.com/akamai/AkamaiOPEN-edgegrid-golang/v2/pkg/session"
	"github.com/akamai/terraform-provider-akamai/v2/pkg/akamai"
	cpstools "github.com/akamai/terraform-provider-akamai/v2/pkg/providers/cps/tools"
	"github.com/akamai/terraform-provider-akamai/v2/pkg/tools"
	"github.com/apex/log"
	"github.com/hashicorp/terraform-plugin-sdk/v2/diag"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/customdiff"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/schema"
	"net/url"
	"strconv"
	"strings"
	"time"
)

var (
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
			StateContext: schema.ImportStatePassthroughContext,
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
			"auto_renewal_start_time": {
				Type:     schema.TypeString,
				Optional: true,
			},
			"certificate_chain_type": {
				Type:     schema.TypeString,
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
				Optional: true,
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
				Type:     schema.TypeString,
				ForceNew: true,
				Required: true,
			},
			"certificate_type": {
				Type:     schema.TypeString,
				Computed: true,
			},
			"validation_type": {
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
				domainsToValidate := make([]interface{}, 0)
				domainsToValidate = append(domainsToValidate, map[string]interface{}{"domain": diff.Get("common_name").(string)})
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
	log := meta.Log("CPS", "resourceDVEnrollment")
	// create a context with logging for api calls
	ctx = session.ContextWithOptions(
		ctx,
		session.WithContextLog(log),
	)
	client := inst.Client(meta)
	log.Debug("Creating enrollment")

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
	adminContact, err := getContactInfo(adminContactSet)
	if err != nil {
		return diag.FromErr(fmt.Errorf("'admin_contact' - %s", err))
	}
	enrollment.AdminContact = adminContact
	techContactSet, err := tools.GetSetValue("tech_contact", d)
	if err != nil {
		return diag.FromErr(err)
	}
	techContact, err := getContactInfo(techContactSet)
	if err != nil {
		return diag.FromErr(fmt.Errorf("'tech_contact' - %s", err))
	}
	enrollment.TechContact = techContact

	autoRenewal, err := tools.GetStringValue("auto_renewal_start_time", d)
	if err != nil && !errors.Is(err, tools.ErrNotFound) {
		return diag.FromErr(err)
	}
	enrollment.AutoRenewalStartTime = autoRenewal

	certificateChainType, err := tools.GetStringValue("certificate_chain_type", d)
	if err != nil && !errors.Is(err, tools.ErrNotFound) {
		return diag.FromErr(err)
	}
	enrollment.CertificateChainType = certificateChainType

	csr, err := getCSR(d)
	if err != nil {
		return diag.FromErr(err)
	}
	enrollment.CSR = csr

	enableMultiStacked, err := tools.GetBoolValue("enable_multi_stacked_certificates", d)
	if err != nil && !errors.Is(err, tools.ErrNotFound) {
		return diag.FromErr(err)
	}
	enrollment.EnableMultiStackedCertificates = enableMultiStacked

	networkConfig, err := getNetworkConfig(d)
	if err != nil {
		return diag.FromErr(err)
	}
	enrollment.NetworkConfiguration = networkConfig
	signatureAlgorithm, err := tools.GetStringValue("signature_algorithm", d)
	if err != nil && !errors.Is(err, tools.ErrNotFound) {
		return diag.FromErr(err)
	}
	enrollment.SignatureAlgorithm = signatureAlgorithm

	organization, err := getOrg(d)
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
		ContractID: contractID,
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
	err = waitForVerification(ctx, log, client, res.ID, acknowledgeWarnings)
	if err != nil {
		d.Partial(true)
		return diag.FromErr(err)
	}
	return resourceCPSDVEnrollmentRead(ctx, d, m)
}

func resourceCPSDVEnrollmentRead(ctx context.Context, d *schema.ResourceData, m interface{}) diag.Diagnostics {
	meta := akamai.Meta(m)
	log := meta.Log("CPS", "resourceDVEnrollment")
	// create a context with logging for api calls
	ctx = session.ContextWithOptions(
		ctx,
		session.WithContextLog(log),
	)
	client := inst.Client(meta)
	log.Debug("Reading enrollment")
	enrollmentID, err := strconv.Atoi(d.Id())
	if err != nil {
		return diag.FromErr(err)
	}
	req := cps.GetEnrollmentRequest{EnrollmentID: enrollmentID}
	enrollment, err := client.GetEnrollment(ctx, req)
	if err != nil {
		return diag.FromErr(err)
	}
	attrs := make(map[string]interface{}, 0)
	adminContact := contactInfoToMap(*enrollment.AdminContact)
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
	techContact := contactInfoToMap(*enrollment.TechContact)
	attrs["tech_contact"] = []interface{}{techContact}
	attrs["auto_renewal_start_time"] = enrollment.AutoRenewalStartTime
	attrs["certificate_chain_type"] = enrollment.CertificateChainType
	csr := csrToMap(*enrollment.CSR)
	attrs["csr"] = []interface{}{csr}
	attrs["enable_multi_stacked_certificates"] = enrollment.EnableMultiStackedCertificates
	networkConfig := networkConfigToMap(*enrollment.NetworkConfiguration)
	attrs["network_configuration"] = []interface{}{networkConfig}
	attrs["signature_algorithm"] = enrollment.SignatureAlgorithm
	org := orgToMap(*enrollment.Org)
	attrs["organization"] = []interface{}{org}
	attrs["certificate_type"] = enrollment.CertificateType
	attrs["validation_type"] = enrollment.ValidationType

	err = cpstools.SetBatch(ctx, d, attrs)
	if err != nil {
		return diag.FromErr(err)
	}
	changeID, err := getChangeIDFromPendingChanges(enrollment.PendingChanges)
	if err != nil {
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
	dnsChallenges := make([]interface{}, 0)
	httpChallenges := make([]interface{}, 0)
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
	log := meta.Log("CPS", "resourceDVEnrollment")
	// create a context with logging for api calls
	ctx = session.ContextWithOptions(
		ctx,
		session.WithContextLog(log),
	)
	client := inst.Client(meta)
	log.Debug("Updating enrollment")

	enrollmentID, err := strconv.Atoi(d.Id())
	if err != nil {
		return diag.FromErr(err)
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
	adminContact, err := getContactInfo(adminContactSet)
	if err != nil {
		return diag.FromErr(fmt.Errorf("'admin_contact' - %s", err))
	}
	enrollment.AdminContact = adminContact
	techContactSet, err := tools.GetSetValue("tech_contact", d)
	if err != nil {
		return diag.FromErr(err)
	}
	techContact, err := getContactInfo(techContactSet)
	if err != nil {
		return diag.FromErr(fmt.Errorf("'tech_contact' - %s", err))
	}
	enrollment.TechContact = techContact

	autoRenewal, err := tools.GetStringValue("auto_renewal_start_time", d)
	if err != nil && !errors.Is(err, tools.ErrNotFound) {
		return diag.FromErr(err)
	}
	enrollment.AutoRenewalStartTime = autoRenewal

	certificateChainType, err := tools.GetStringValue("certificate_chain_type", d)
	if err != nil && !errors.Is(err, tools.ErrNotFound) {
		return diag.FromErr(err)
	}
	enrollment.CertificateChainType = certificateChainType

	csr, err := getCSR(d)
	if err != nil {
		return diag.FromErr(err)
	}
	enrollment.CSR = csr

	enableMultiStacked, err := tools.GetBoolValue("enable_multi_stacked_certificates", d)
	if err != nil && !errors.Is(err, tools.ErrNotFound) {
		return diag.FromErr(err)
	}
	enrollment.EnableMultiStackedCertificates = enableMultiStacked

	networkConfig, err := getNetworkConfig(d)
	if err != nil {
		return diag.FromErr(err)
	}
	enrollment.NetworkConfiguration = networkConfig
	signatureAlgorithm, err := tools.GetStringValue("signature_algorithm", d)
	if err != nil && !errors.Is(err, tools.ErrNotFound) {
		return diag.FromErr(err)
	}
	enrollment.SignatureAlgorithm = signatureAlgorithm

	organization, err := getOrg(d)
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
	res, err := client.UpdateEnrollment(ctx, req)
	if err != nil {
		return diag.FromErr(err)
	}
	d.SetId(strconv.Itoa(res.ID))

	acknowledgeWarnings, err := tools.GetBoolValue("acknowledge_pre_verification_warnings", d)
	if err != nil && !errors.Is(err, tools.ErrNotFound) {
		return diag.FromErr(err)
	}
	err = waitForVerification(ctx, log, client, res.ID, acknowledgeWarnings)
	if err != nil {
		d.Partial(true)
		return diag.FromErr(err)
	}
	return resourceCPSDVEnrollmentRead(ctx, d, m)
}

func waitForVerification(ctx context.Context, logger log.Interface, client cps.CPS, enrollmentID int, acknowledgeWarnings bool) error {
	getEnrollmentReq := cps.GetEnrollmentRequest{EnrollmentID: enrollmentID}
	enrollmentGet, err := client.GetEnrollment(ctx, getEnrollmentReq)
	if err != nil {
		return err
	}
	changeID, err := getChangeIDFromPendingChanges(enrollmentGet.PendingChanges)
	if err != nil {
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
	for status.StatusInfo.Status != statusCoordinateDomainValidation {
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

func resourceCPSDVEnrollmentDelete(ctx context.Context, d *schema.ResourceData, m interface{}) diag.Diagnostics {
	meta := akamai.Meta(m)
	log := meta.Log("CPS", "resourceDVEnrollment")
	// create a context with logging for api calls
	ctx = session.ContextWithOptions(
		ctx,
		session.WithContextLog(log),
	)
	client := inst.Client(meta)
	log.Debug("Deleting enrollment")
	enrollmentID, err := strconv.Atoi(d.Id())
	if err != nil {
		return diag.FromErr(err)
	}
	allowCancelPendingChanges := true
	req := cps.RemoveEnrollmentRequest{
		EnrollmentID:              enrollmentID,
		AllowCancelPendingChanges: &allowCancelPendingChanges,
	}
	_, err = client.RemoveEnrollment(ctx, req)
	if err != nil {
		return diag.FromErr(err)
	}
	d.SetId("")
	return nil
}

func getContactInfo(set *schema.Set) (*cps.Contact, error) {
	contactList := set.List()
	contactMap, ok := contactList[0].(map[string]interface{})
	if !ok {
		return nil, fmt.Errorf("contact is of invalid type")
	}

	var contact cps.Contact

	firstname := contactMap["first_name"].(string)
	lastname := contactMap["last_name"].(string)
	title := contactMap["title"].(string)
	organization := contactMap["organization"].(string)
	email := contactMap["email"].(string)
	phone := contactMap["phone"].(string)
	addresslineone := contactMap["address_line_one"].(string)
	addresslinetwo := contactMap["address_line_two"].(string)
	city := contactMap["city"].(string)
	region := contactMap["region"].(string)
	postalcode := contactMap["postal_code"].(string)
	country := contactMap["country_code"].(string)

	contact.FirstName = firstname
	contact.LastName = lastname
	contact.Title = title
	contact.OrganizationName = organization
	contact.Email = email
	contact.Phone = phone
	contact.AddressLineOne = addresslineone
	contact.AddressLineTwo = addresslinetwo
	contact.City = city
	contact.Region = region
	contact.PostalCode = postalcode
	contact.Country = country

	return &contact, nil
}

func getCSR(d *schema.ResourceData) (*cps.CSR, error) {
	num := 1
	switch num {
	case 0, 1:

	}
	csrSet, err := tools.GetSetValue("csr", d)
	if err != nil {
		return nil, err
	}
	commonName, err := tools.GetStringValue("common_name", d)
	if err != nil {
		return nil, err
	}
	csrList := csrSet.List()
	csrmap, ok := csrList[0].(map[string]interface{})
	if !ok {
		return nil, fmt.Errorf("'csr' is of invalid type")
	}

	var csr cps.CSR

	sansList, err := tools.GetSetValue("sans", d)
	if err != nil && !errors.Is(err, tools.ErrNotFound) {
		return nil, err
	}
	var sans []string
	for _, val := range sansList.List() {
		sans = append(sans, val.(string))
	}
	csr.SANS = sans

	csr.CN = commonName
	csr.L = csrmap["city"].(string)
	csr.ST = csrmap["state"].(string)
	csr.C = csrmap["country_code"].(string)
	csr.O = csrmap["organization"].(string)
	csr.OU = csrmap["organizational_unit"].(string)

	return &csr, nil
}

func getNetworkConfig(d *schema.ResourceData) (*cps.NetworkConfiguration, error) {
	networkConfigSet, err := tools.GetSetValue("network_configuration", d)
	if err != nil {
		return nil, err
	}
	sniOnly, err := tools.GetBoolValue("sni_only", d)
	if err != nil {
		return nil, err
	}
	secureNetwork, err := tools.GetStringValue("secure_network", d)
	if err != nil {
		return nil, err
	}
	networkConfigMap, ok := networkConfigSet.List()[0].(map[string]interface{})
	if !ok {
		return nil, fmt.Errorf("'network_configuration' is of invalid type")
	}
	var networkConfig cps.NetworkConfiguration

	if val, ok := networkConfigMap["client_mutual_authentication"]; ok {
		mutualAuth := &cps.ClientMutualAuthentication{}
		mutualAuthSet, ok := val.(*schema.Set)
		if !ok {
			return nil, fmt.Errorf("'client_mutual_authentication' is of invalid type")
		}
		if len(mutualAuthSet.List()) > 0 {
			mutualAuthMap := mutualAuthSet.List()[0].(map[string]interface{})
			if ocspEnabled, ok := mutualAuthMap["ocsp_enabled"]; ok {
				ocspEnabledBool := ocspEnabled.(bool)
				mutualAuth.AuthenticationOptions = &cps.AuthenticationOptions{
					OCSP:               &cps.OCSP{Enabled: &ocspEnabledBool},
					SendCAListToClient: nil,
				}
			}
			if sendCa, ok := mutualAuthMap["send_ca_list_to_client"]; ok {
				sendCaBool := sendCa.(bool)
				mutualAuth.AuthenticationOptions.SendCAListToClient = &sendCaBool
			}
			mutualAuth.SetID = networkConfigMap["mutual_authentication_set_id"].(string)
			networkConfig.ClientMutualAuthentication = mutualAuth
		}
	}
	sansList, err := tools.GetSetValue("sans", d)
	if err != nil && !errors.Is(err, tools.ErrNotFound) {
		return nil, err
	}
	var dnsNames []string
	for _, val := range sansList.List() {
		dnsNames = append(dnsNames, val.(string))
	}
	networkConfig.DNSNameSettings = &cps.DNSNameSettings{
		CloneDNSNames: networkConfigMap["clone_dns_names"].(bool),
		DNSNames:      dnsNames,
	}
	networkConfig.OCSPStapling = cps.OCSPStapling(networkConfigMap["ocsp_stapling"].(string))
	disallowedTLSVersionsArray := networkConfigMap["disallowed_tls_versions"].(*schema.Set)
	var disallowedTLSVersions []string
	for _, val := range disallowedTLSVersionsArray.List() {
		disallowedTLSVersions = append(disallowedTLSVersions, val.(string))
	}
	networkConfig.DisallowedTLSVersions = disallowedTLSVersions
	networkConfig.Geography = networkConfigMap["geography"].(string)
	networkConfig.MustHaveCiphers = networkConfigMap["must_have_ciphers"].(string)
	networkConfig.PreferredCiphers = networkConfigMap["preferred_ciphers"].(string)
	networkConfig.QuicEnabled = networkConfigMap["quic_enabled"].(bool)
	networkConfig.SecureNetwork = secureNetwork
	networkConfig.SNIOnly = sniOnly

	return &networkConfig, nil
}

func getOrg(d *schema.ResourceData) (*cps.Org, error) {
	orgSet, err := tools.GetSetValue("organization", d)
	if err != nil {
		return nil, err
	}
	orgMap, ok := orgSet.List()[0].(map[string]interface{})
	if !ok {
		return nil, fmt.Errorf("'organization' is of invalid type")
	}

	var org cps.Org

	name := orgMap["name"].(string)
	phone := orgMap["phone"].(string)
	addresslineone := orgMap["address_line_one"].(string)
	addresslinetwo := orgMap["address_line_two"].(string)
	city := orgMap["city"].(string)
	region := orgMap["region"].(string)
	postalcode := orgMap["postal_code"].(string)
	country := orgMap["country_code"].(string)

	org.Name = name
	org.Phone = phone
	org.AddressLineOne = addresslineone
	org.AddressLineTwo = addresslinetwo
	org.City = city
	org.Region = region
	org.PostalCode = postalcode
	org.Country = country

	return &org, nil
}

func contactInfoToMap(contact cps.Contact) map[string]interface{} {
	contactMap := map[string]interface{}{
		"first_name":       contact.FirstName,
		"last_name":        contact.LastName,
		"title":            contact.Title,
		"organization":     contact.OrganizationName,
		"email":            contact.Email,
		"phone":            contact.Phone,
		"address_line_one": contact.AddressLineOne,
		"address_line_two": contact.AddressLineTwo,
		"city":             contact.City,
		"region":           contact.Region,
		"postal_code":      contact.PostalCode,
		"country_code":     contact.Country,
	}

	return contactMap
}

func csrToMap(csr cps.CSR) map[string]interface{} {
	csrMap := map[string]interface{}{
		"country_code":        csr.C,
		"city":                csr.L,
		"organization":        csr.O,
		"organizational_unit": csr.OU,
		"state":               csr.ST,
	}
	return csrMap
}

func networkConfigToMap(networkConfig cps.NetworkConfiguration) map[string]interface{} {
	networkConfigMap := make(map[string]interface{})
	if networkConfig.ClientMutualAuthentication != nil {
		networkConfigMap["set_id"] = networkConfig.ClientMutualAuthentication.SetID
		if networkConfig.ClientMutualAuthentication.AuthenticationOptions != nil {
			networkConfigMap["mutual_authentication_send_ca_list_to_client"] = networkConfig.ClientMutualAuthentication.AuthenticationOptions.SendCAListToClient
			if networkConfig.ClientMutualAuthentication.AuthenticationOptions.OCSP != nil {
				networkConfigMap["mutual_authentication_oscp_enabled"] = *networkConfig.ClientMutualAuthentication.AuthenticationOptions.OCSP.Enabled
			}
		}
	}
	networkConfigMap["disallowed_tls_versions"] = networkConfig.DisallowedTLSVersions
	if networkConfig.DNSNameSettings != nil {
		networkConfigMap["clone_dns_names"] = networkConfig.DNSNameSettings.CloneDNSNames
	}
	networkConfigMap["geography"] = networkConfig.Geography
	networkConfigMap["must_have_ciphers"] = networkConfig.MustHaveCiphers
	networkConfigMap["ocsp_stapling"] = networkConfig.OCSPStapling
	networkConfigMap["preferred_ciphers"] = networkConfig.PreferredCiphers
	networkConfigMap["quic_enabled"] = networkConfig.QuicEnabled
	return networkConfigMap
}

func orgToMap(org cps.Org) map[string]interface{} {
	orgMap := map[string]interface{}{
		"name":             org.Name,
		"phone":            org.Phone,
		"address_line_one": org.AddressLineOne,
		"address_line_two": org.AddressLineTwo,
		"city":             org.City,
		"region":           org.Region,
		"postal_code":      org.PostalCode,
		"country_code":     org.Country,
	}

	return orgMap
}

func getChangeIDFromPendingChanges(pendingChanges []string) (int, error) {
	if len(pendingChanges) < 1 {
		return 0, fmt.Errorf("no pending changes were found on enrollment")
	}
	changeURL, err := url.Parse(pendingChanges[0])
	if err != nil {
		return 0, err
	}
	pathSplit := strings.Split(changeURL.Path, "/")
	changeIDStr := pathSplit[len(pathSplit)-1]
	changeID, err := strconv.Atoi(changeIDStr)
	if err != nil {
		return 0, err
	}
	return changeID, nil
}
