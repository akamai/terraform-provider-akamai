package cps

import (
	"context"
	"errors"
	"fmt"
	"strconv"
	"strings"
	"time"

	"github.com/akamai/AkamaiOPEN-edgegrid-golang/v3/pkg/cps"
	"github.com/akamai/AkamaiOPEN-edgegrid-golang/v3/pkg/session"
	"github.com/akamai/terraform-provider-akamai/v3/pkg/akamai"
	cpstools "github.com/akamai/terraform-provider-akamai/v3/pkg/providers/cps/tools"
	"github.com/akamai/terraform-provider-akamai/v3/pkg/tools"
	"github.com/hashicorp/terraform-plugin-sdk/v2/diag"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/customdiff"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/schema"
)

var (
	// PollForChangeStatusInterval defines retry interval for getting status of a pending change
	PollForChangeStatusInterval = 10 * time.Second
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
			"allow_duplicate_common_name": {
				Type:     schema.TypeBool,
				Optional: true,
				Default:  false,
			},
			"sans": {
				Type:     schema.TypeSet,
				Optional: true,
				Elem:     &schema.Schema{Type: schema.TypeString},
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
				Elem:     csr,
			},
			"enable_multi_stacked_certificates": {
				Type:       schema.TypeBool,
				Optional:   true,
				Deprecated: "Deprecated, don't use; always false",
			},
			"network_configuration": {
				Type:     schema.TypeSet,
				Required: true,
				MinItems: 1,
				MaxItems: 1,
				Elem:     networkConfiguration,
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
				Elem:     organization,
			},
			"contract_id": {
				Type:             schema.TypeString,
				ForceNew:         true,
				Required:         true,
				DiffSuppressFunc: tools.FieldPrefixSuppress("ctr_"),
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

func resourceCPSDVEnrollmentCreate(ctx context.Context, d *schema.ResourceData, m interface{}) diag.Diagnostics {
	meta := akamai.Meta(m)
	logger := meta.Log("CPS", "resourceCPSDVEnrollmentCreate")
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
		return diag.Errorf("'admin_contact' - %s", err)
	}
	enrollment.AdminContact = adminContact
	techContactSet, err := tools.GetSetValue("tech_contact", d)
	if err != nil {
		return diag.FromErr(err)
	}
	techContact, err := cpstools.GetContactInfo(techContactSet)
	if err != nil {
		return diag.Errorf("'tech_contact' - %s", err)
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

	// DV does not support multi stack certificates
	enrollment.EnableMultiStackedCertificates = false

	networkConfig, err := cpstools.GetNetworkConfig(d)
	if err != nil {
		return diag.FromErr(err)
	}
	enrollment.NetworkConfiguration = networkConfig
	signatureAlgorithm, err := tools.GetStringValue("signature_algorithm", d)
	if err != nil {
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
	allowDuplicateCN, err := tools.GetBoolValue("allow_duplicate_common_name", d)
	if err != nil {
		return diag.FromErr(err)
	}

	req := cps.CreateEnrollmentRequest{
		Enrollment:       enrollment,
		ContractID:       strings.TrimPrefix(contractID, "ctr_"),
		AllowDuplicateCN: allowDuplicateCN,
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
	if err = waitForVerification(ctx, logger, client, res.ID, acknowledgeWarnings, nil); err != nil {
		return diag.FromErr(err)
	}
	if err != nil {
		return diag.FromErr(err)
	}
	return resourceCPSDVEnrollmentRead(ctx, d, m)
}

func resourceCPSDVEnrollmentRead(ctx context.Context, d *schema.ResourceData, m interface{}) diag.Diagnostics {
	meta := akamai.Meta(m)
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
	logger := meta.Log("CPS", "resourceCPSDVEnrollmentUpdate")
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
	if err := d.Set("registration_authority", "lets-encrypt"); err != nil {
		return diag.Errorf("%v: %s", tools.ErrValueSet, err.Error())
	}

	adminContactSet, err := tools.GetSetValue("admin_contact", d)
	if err != nil {
		return diag.FromErr(err)
	}
	adminContact, err := cpstools.GetContactInfo(adminContactSet)
	if err != nil {
		return diag.Errorf("'admin_contact' - %s", err)
	}
	enrollment.AdminContact = adminContact
	techContactSet, err := tools.GetSetValue("tech_contact", d)
	if err != nil {
		return diag.FromErr(err)
	}
	techContact, err := cpstools.GetContactInfo(techContactSet)
	if err != nil {
		return diag.Errorf("'tech_contact' - %s", err)
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

	// DV does not support multi stack certificates
	enrollment.EnableMultiStackedCertificates = false

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

	if err = waitForVerification(ctx, logger, client, enrollmentID, acknowledgeWarnings, nil); err != nil {
		return diag.FromErr(err)
	}
	return resourceCPSDVEnrollmentRead(ctx, d, m)
}

func resourceCPSDVEnrollmentDelete(ctx context.Context, d *schema.ResourceData, m interface{}) diag.Diagnostics {
	return enrollmentDelete(ctx, d, m, "resourceCPSDVEnrollmentDelete")
}

func resourceCPSDVEnrollmentImport(ctx context.Context, d *schema.ResourceData, m interface{}) ([]*schema.ResourceData, error) {
	meta := akamai.Meta(m)
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
		return nil, fmt.Errorf("%v: %s", tools.ErrValueSet, err.Error())
	}
	if err := d.Set("contract_id", contractID); err != nil {
		return nil, fmt.Errorf("%v: %s", tools.ErrValueSet, err.Error())
	}
	d.SetId(enrollmentID)
	return []*schema.ResourceData{d}, nil
}
