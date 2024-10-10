package cps

import (
	"context"
	"errors"
	"fmt"
	"strconv"
	"strings"

	"github.com/akamai/AkamaiOPEN-edgegrid-golang/v9/pkg/cps"
	"github.com/akamai/AkamaiOPEN-edgegrid-golang/v9/pkg/session"
	"github.com/akamai/terraform-provider-akamai/v6/pkg/common/ptr"
	"github.com/akamai/terraform-provider-akamai/v6/pkg/common/tf"
	"github.com/akamai/terraform-provider-akamai/v6/pkg/common/timeouts"
	"github.com/akamai/terraform-provider-akamai/v6/pkg/meta"
	cpstools "github.com/akamai/terraform-provider-akamai/v6/pkg/providers/cps/tools"
	"github.com/hashicorp/terraform-plugin-sdk/v2/diag"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/customdiff"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/schema"
)

var (
	// ErrWarningsCannotBeApproved is returned when some warnings cannot be auto approved
	ErrWarningsCannotBeApproved = errors.New("warnings cannot be approved")
)

func resourceCPSThirdPartyEnrollment() *schema.Resource {
	return &schema.Resource{
		CreateContext: resourceCPSThirdPartyEnrollmentCreate,
		ReadContext:   resourceCPSThirdPartyEnrollmentRead,
		UpdateContext: resourceCPSThirdPartyEnrollmentUpdate,
		DeleteContext: resourceCPSThirdPartyEnrollmentDelete,
		Importer: &schema.ResourceImporter{
			StateContext: resourceCPSThirdPartyEnrollmentImport,
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
				Description: "Allow to duplicate common name",
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
				Description: "Whether acknowledge warnings before certificate verification",
			},
			"auto_approve_warnings": {
				Type:        schema.TypeSet,
				Optional:    true,
				Elem:        &schema.Schema{Type: schema.TypeString},
				Description: "List of warnings to be automatically approved",
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
				Description: "Certificate trust chain type",
			},
			"csr": {
				Type:        schema.TypeSet,
				Required:    true,
				MinItems:    1,
				MaxItems:    1,
				Elem:        csr,
				Description: "Data used for generation of Certificate Signing Request",
			},
			"network_configuration": {
				Type:        schema.TypeSet,
				Required:    true,
				MinItems:    1,
				MaxItems:    1,
				Elem:        networkConfiguration,
				Description: "Settings containing network information and TLS metadata used by CPS",
			},
			"signature_algorithm": {
				Type:             schema.TypeString,
				Optional:         true,
				DiffSuppressFunc: supressSignatureAlgorithm,
				Description:      "The SHA function. Changing this value may require running terraform destroy, terraform apply",
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
			"change_management": {
				Type:        schema.TypeBool,
				Optional:    true,
				Description: "When set to false, the certificate will be deployed to both staging and production networks",
			},
			"exclude_sans": {
				Type:        schema.TypeBool,
				Optional:    true,
				Description: "When true, SANs are excluded from the CSR",
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
				return nil
			}),
		Timeouts: &schema.ResourceTimeout{
			Default: &timeouts.SDKDefaultTimeout,
		},
	}
}

func supressSignatureAlgorithm(_ string, oldValue, newValue string, d *schema.ResourceData) bool {
	if oldValue == "" && d != nil && d.Id() != "" {
		return true
	}
	return oldValue == newValue
}

func resourceCPSThirdPartyEnrollmentCreate(ctx context.Context, d *schema.ResourceData, m interface{}) diag.Diagnostics {
	meta := meta.Must(m)
	logger := meta.Log("CPS", "resourceCPSThirdPartyEnrollmentCreate")
	// create a context with logging for api calls
	ctx = session.ContextWithOptions(
		ctx,
		session.WithContextLog(logger),
	)
	client := inst.Client(meta)
	logger.Debug("Creating enrollment")

	contractID, err := tf.GetStringValue("contract_id", d)
	if err != nil {
		return diag.FromErr(err)
	}
	allowDuplicateCN, err := tf.GetBoolValue("allow_duplicate_common_name", d)
	if err != nil {
		return diag.FromErr(err)
	}

	enrollmentReqBody, err := prepareThirdPartyEnrollment(d)
	if err != nil {
		return diag.FromErr(err)
	}

	// save ClientMutualAuthentication and unset it in enrollment request struct
	// create request must not have it set; in case its not nil, we will run update later to add it
	clientMutualAuthentication := enrollmentReqBody.NetworkConfiguration.ClientMutualAuthentication
	enrollmentReqBody.NetworkConfiguration.ClientMutualAuthentication = nil

	req := cps.CreateEnrollmentRequest{
		EnrollmentRequestBody: *enrollmentReqBody,
		ContractID:            strings.TrimPrefix(contractID, "ctr_"),
		AllowDuplicateCN:      allowDuplicateCN,
	}
	res, err := client.CreateEnrollment(ctx, req)
	if err != nil {
		return diag.FromErr(err)
	}
	d.SetId(strconv.Itoa(res.ID))

	// when clientMutualAuthentication was provided, insert it back to enrollment and send the update request
	if clientMutualAuthentication != nil {
		logger.Debug("Updating ClientMutualAuthentication configuration")
		enrollmentReqBody.NetworkConfiguration.ClientMutualAuthentication = clientMutualAuthentication
		req := cps.UpdateEnrollmentRequest{
			EnrollmentID:              res.ID,
			EnrollmentRequestBody:     *enrollmentReqBody,
			AllowCancelPendingChanges: ptr.To(true),
		}
		_, err := client.UpdateEnrollment(ctx, req)
		if err != nil {
			return diag.FromErr(err)
		}
	}

	acknowledgeWarnings, err := tf.GetBoolValue("acknowledge_pre_verification_warnings", d)
	if err != nil && !errors.Is(err, tf.ErrNotFound) {
		return diag.FromErr(err)
	}
	autoApproveWarnings, err := tf.GetSetValue("auto_approve_warnings", d)
	if err != nil && !errors.Is(err, tf.ErrNotFound) {
		return diag.FromErr(err)
	}
	autoApproveWarningsAsString := convertUserWarningsToStringSlice(autoApproveWarnings.List())

	if err = waitForVerification(ctx, logger, client, res.ID, acknowledgeWarnings, autoApproveWarningsAsString); err != nil {
		return diag.FromErr(err)
	}
	return resourceCPSThirdPartyEnrollmentRead(ctx, d, m)
}

func resourceCPSThirdPartyEnrollmentRead(ctx context.Context, d *schema.ResourceData, m interface{}) diag.Diagnostics {
	meta := meta.Must(m)
	logger := meta.Log("CPS", "resourceCPSThirdPartyEnrollmentRead")
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
	enrollment, err := client.GetEnrollment(ctx, cps.GetEnrollmentRequest{EnrollmentID: enrollmentID})
	if err != nil {
		return diag.FromErr(err)
	}
	attrs, err := readAttrs(enrollment, d)
	if err != nil {
		return diag.FromErr(err)
	}

	var excludeSANS bool
	if enrollment.ThirdParty != nil {
		excludeSANS = enrollment.ThirdParty.ExcludeSANS
	}
	attrs["exclude_sans"] = excludeSANS
	attrs["change_management"] = enrollment.ChangeManagement

	if err = tf.SetAttrs(d, attrs); err != nil {
		return diag.FromErr(err)
	}
	if err != nil {
		return diag.FromErr(err)
	}
	return nil
}

func resourceCPSThirdPartyEnrollmentUpdate(ctx context.Context, d *schema.ResourceData, m interface{}) diag.Diagnostics {
	meta := meta.Must(m)
	logger := meta.Log("CPS", "resourceCPSThirdPartyEnrollmentUpdate")
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
	autoApproveWarnings, err := tf.GetSetValue("auto_approve_warnings", d)
	if err != nil && !errors.Is(err, tf.ErrNotFound) {
		return diag.FromErr(err)
	}
	autoApproveWarningsAsString := convertUserWarningsToStringSlice(autoApproveWarnings.List())

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
		if err = waitForVerification(ctx, logger, client, enrollmentID, acknowledgeWarnings, autoApproveWarningsAsString); err != nil {
			return diag.FromErr(err)
		}
		return resourceCPSThirdPartyEnrollmentRead(ctx, d, m)
	}
	enrollmentReqBody, err := prepareThirdPartyEnrollment(d)
	if err != nil {
		return diag.FromErr(err)
	}

	allowCancel := true
	req := cps.UpdateEnrollmentRequest{
		EnrollmentRequestBody:     *enrollmentReqBody,
		EnrollmentID:              enrollmentID,
		AllowCancelPendingChanges: &allowCancel,
	}

	if _, err := client.UpdateEnrollment(ctx, req); err != nil {
		return diag.FromErr(err)
	}
	d.SetId(strconv.Itoa(enrollmentID))

	if err = waitForVerification(ctx, logger, client, enrollmentID, acknowledgeWarnings, autoApproveWarningsAsString); err != nil {
		return diag.FromErr(err)
	}
	return resourceCPSThirdPartyEnrollmentRead(ctx, d, m)
}

func prepareThirdPartyEnrollment(d *schema.ResourceData) (*cps.EnrollmentRequestBody, error) {
	enrollmentReqBody := cps.EnrollmentRequestBody{
		CertificateType: "third-party",
		ValidationType:  "third-party",
		RA:              "third-party",
	}

	adminContactSet, err := tf.GetSetValue("admin_contact", d)
	if err != nil {
		return nil, err
	}
	adminContact, err := cpstools.GetContactInfo(adminContactSet)
	if err != nil {
		return nil, fmt.Errorf("'admin_contact' - %s", err)
	}
	enrollmentReqBody.AdminContact = adminContact
	techContactSet, err := tf.GetSetValue("tech_contact", d)
	if err != nil {
		return nil, err
	}
	techContact, err := cpstools.GetContactInfo(techContactSet)
	if err != nil {
		return nil, fmt.Errorf("'tech_contact' - %s", err)
	}
	enrollmentReqBody.TechContact = techContact

	certificateChainType, err := tf.GetStringValue("certificate_chain_type", d)
	if err != nil && !errors.Is(err, tf.ErrNotFound) {
		return nil, err
	}
	enrollmentReqBody.CertificateChainType = certificateChainType

	csr, err := cpstools.GetCSR(d)
	if err != nil {
		return nil, err
	}
	enrollmentReqBody.CSR = csr

	// for third-party certificates, multi-stack it is always enabled
	enrollmentReqBody.EnableMultiStackedCertificates = true

	networkConfig, err := cpstools.GetNetworkConfig(d)
	if err != nil {
		return nil, err
	}
	enrollmentReqBody.NetworkConfiguration = networkConfig
	signatureAlgorithm, err := tf.GetStringValue("signature_algorithm", d)
	if err != nil && !errors.Is(err, tf.ErrNotFound) {
		return nil, err
	}
	enrollmentReqBody.SignatureAlgorithm = signatureAlgorithm

	organization, err := cpstools.GetOrg(d)
	if err != nil {
		return nil, err
	}
	enrollmentReqBody.Org = organization
	changeManagement, err := tf.GetBoolValue("change_management", d)
	if err != nil {
		return nil, err
	}
	enrollmentReqBody.ChangeManagement = changeManagement
	excludeSANS, err := tf.GetBoolValue("exclude_sans", d)
	if err != nil {
		return nil, err
	}
	enrollmentReqBody.ThirdParty = &cps.ThirdParty{
		ExcludeSANS: excludeSANS,
	}
	return &enrollmentReqBody, nil
}

func resourceCPSThirdPartyEnrollmentDelete(ctx context.Context, d *schema.ResourceData, m interface{}) diag.Diagnostics {
	return enrollmentDelete(ctx, d, m, "resourceCPSThirdPartyEnrollmentDelete")
}

func resourceCPSThirdPartyEnrollmentImport(ctx context.Context, d *schema.ResourceData, m interface{}) ([]*schema.ResourceData, error) {
	meta := meta.Must(m)
	logger := meta.Log("CPS", "resourceCPSThirdPartyEnrollmentImport")
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
	if enrollment.ValidationType != "third-party" {
		return nil, fmt.Errorf("unable to import: wrong validation type: expected 'third-party', got '%s'", enrollment.ValidationType)
	}

	if err := d.Set("allow_duplicate_common_name", false); err != nil {
		return nil, fmt.Errorf("%v: %s", tf.ErrValueSet, err.Error())
	}
	if err := d.Set("acknowledge_pre_verification_warnings", false); err != nil {
		return nil, fmt.Errorf("%v: %s", tf.ErrValueSet, err.Error())
	}
	if err := d.Set("auto_approve_warnings", []string{}); err != nil {
		return nil, fmt.Errorf("%v: %s", tf.ErrValueSet, err.Error())
	}
	if err := d.Set("contract_id", contractID); err != nil {
		return nil, fmt.Errorf("%v: %s", tf.ErrValueSet, err.Error())
	}
	d.SetId(enrollmentID)
	return []*schema.ResourceData{d}, nil
}
