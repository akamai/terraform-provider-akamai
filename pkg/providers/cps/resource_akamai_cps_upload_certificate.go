package cps

import (
	"context"
	"errors"
	"fmt"
	"strconv"
	"strings"
	"time"

	"github.com/akamai/AkamaiOPEN-edgegrid-golang/v7/pkg/cps"
	"github.com/akamai/AkamaiOPEN-edgegrid-golang/v7/pkg/session"
	"github.com/akamai/terraform-provider-akamai/v5/pkg/common/tf"
	"github.com/akamai/terraform-provider-akamai/v5/pkg/common/timeouts"
	"github.com/akamai/terraform-provider-akamai/v5/pkg/meta"
	toolsCPS "github.com/akamai/terraform-provider-akamai/v5/pkg/providers/cps/tools"
	"github.com/apex/log"
	"github.com/hashicorp/terraform-plugin-sdk/v2/diag"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/schema"
)

func resourceCPSUploadCertificate() *schema.Resource {
	return &schema.Resource{
		Description:   "Enables to upload a certificates and trust-chains for third-party enrollment",
		CreateContext: resourceCPSUploadCertificateCreate,
		ReadContext:   resourceCPSUploadCertificateRead,
		UpdateContext: resourceCPSUploadCertificateUpdate,
		DeleteContext: resourceCPSUploadCertificateDelete,
		Importer: &schema.ResourceImporter{
			StateContext: resourceCPSUploadCertificateImport,
		},
		CustomizeDiff: checkUnacknowledgedWarnings,
		Timeouts: &schema.ResourceTimeout{
			Default: &defaultTimeout,
		},
		Schema: map[string]*schema.Schema{
			"enrollment_id": {
				Type:        schema.TypeInt,
				Required:    true,
				ForceNew:    true,
				Description: "The unique identifier of the enrollment",
			},
			"certificate_ecdsa_pem": {
				Type:             schema.TypeString,
				Optional:         true,
				AtLeastOneOf:     []string{"certificate_ecdsa_pem", "certificate_rsa_pem"},
				Description:      "ECDSA certificate in pem format to be uploaded",
				DiffSuppressFunc: trimWhitespaces,
			},
			"certificate_rsa_pem": {
				Type:             schema.TypeString,
				Optional:         true,
				AtLeastOneOf:     []string{"certificate_ecdsa_pem", "certificate_rsa_pem"},
				Description:      "RSA certificate in pem format to be uploaded",
				DiffSuppressFunc: trimWhitespaces,
			},
			"trust_chain_ecdsa_pem": {
				Type:             schema.TypeString,
				Optional:         true,
				Description:      "Trust chain in pem format for provided ECDSA certificate",
				DiffSuppressFunc: trimWhitespaces,
			},
			"trust_chain_rsa_pem": {
				Type:             schema.TypeString,
				Optional:         true,
				Description:      "Trust chain in pem format for provided RSA certificate",
				DiffSuppressFunc: trimWhitespaces,
			},
			"acknowledge_post_verification_warnings": {
				Type:        schema.TypeBool,
				Optional:    true,
				Default:     false,
				Description: "Whether to acknowledge post-verification warnings",
			},
			"auto_approve_warnings": {
				Type:        schema.TypeSet,
				Optional:    true,
				Elem:        &schema.Schema{Type: schema.TypeString},
				Description: "List of post-verification warnings to be automatically acknowledged",
			},
			"acknowledge_change_management": {
				Type:        schema.TypeBool,
				Optional:    true,
				Default:     false,
				Description: "Whether to acknowledge change management",
			},
			"wait_for_deployment": {
				Type:        schema.TypeBool,
				Optional:    true,
				Default:     false,
				Description: "Whether to wait for certificate to be deployed",
			},
			"unacknowledged_warnings": {
				Type:                  schema.TypeBool,
				ConfigMode:            0,
				Required:              false,
				Optional:              false,
				Computed:              true,
				ForceNew:              false,
				DiffSuppressFunc:      nil,
				DiffSuppressOnRefresh: false,
				Default:               nil,
				DefaultFunc:           nil,
				Description:           "Used to distinguish whether there are unacknowledged warnings for a certificate",
				InputDefault:          "",
				StateFunc:             nil,
				Elem:                  nil,
				MaxItems:              0,
				MinItems:              0,
				Set:                   nil,
				ComputedWhen:          nil,
				ConflictsWith:         nil,
				ExactlyOneOf:          nil,
				AtLeastOneOf:          nil,
				RequiredWith:          nil,
				Deprecated:            "",
				ValidateFunc:          nil,
				ValidateDiagFunc:      nil,
				Sensitive:             false,
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
	}
}

// attributes contains attributes from schema
type attributes struct {
	enrollmentID        int
	certificateECDSA    string
	certificateRSA      string
	trustChainECDSA     string
	trustChainRSA       string
	ackChangeManagement bool
	waitForDeployment   bool
}

var (
	defaultTimeout = time.Hour * 2

	trimWhitespaces = func(k, oldValue, newValue string, d *schema.ResourceData) bool {
		return strings.TrimSpace(oldValue) == strings.TrimSpace(newValue)
	}
)

func resourceCPSUploadCertificateCreate(ctx context.Context, d *schema.ResourceData, m interface{}) diag.Diagnostics {
	meta := meta.Must(m)
	logger := meta.Log("CPS", "resourceCPSUploadCertificateCreate")
	ctx = session.ContextWithOptions(ctx, session.WithContextLog(logger))
	client := inst.Client(meta)
	logger.Debug("Creating upload certificate")

	return upsertUploadCertificate(ctx, d, m, client, logger)
}

func resourceCPSUploadCertificateRead(ctx context.Context, d *schema.ResourceData, m interface{}) diag.Diagnostics {
	meta := meta.Must(m)
	logger := meta.Log("CPS", "resourceCPSUploadCertificateRead")
	ctx = session.ContextWithOptions(ctx, session.WithContextLog(logger))
	client := inst.Client(meta)
	logger.Debug("Reading upload certificate")

	enrollmentID, err := strconv.Atoi(d.Id())
	if err != nil {
		return diag.Errorf("could not get resource id: %s", err)
	}
	waitForDeployment, err := tf.GetBoolValue("wait_for_deployment", d)
	if err != nil {
		return diag.Errorf("could not get `wait for deployment` attribute: %s", err)
	}

	enrollment, err := client.GetEnrollment(ctx, cps.GetEnrollmentRequest{EnrollmentID: enrollmentID})
	if err != nil {
		return diag.Errorf("could not get an enrollment: %s", err)
	}

	changeID, err := toolsCPS.GetChangeIDFromPendingChanges(enrollment.PendingChanges)
	if err != nil && !errors.Is(err, toolsCPS.ErrNoPendingChanges) {
		return diag.Errorf("could not get changeID of an enrollment: %s", err)
	}

	var attrs map[string]interface{}
	if waitForDeployment && len(enrollment.PendingChanges) != 0 {
		ackChangeManagement, err := tf.GetBoolValue("acknowledge_change_management", d)
		if err != nil {
			return diag.Errorf("could not get `acknowledge change management` attribute: %s", err)
		}

		changeStatus, err := sendGetChangeStatusReq(ctx, client, enrollmentID, changeID)
		if err != nil {
			return diag.FromErr(err)
		}

		if changeStatus.StatusInfo.Status != waitUploadThirdParty && changeStatus.StatusInfo.Status != waitReviewThirdPartyCert {
			statusToWaitFor := complete
			if !ackChangeManagement && enrollment.ChangeManagement {
				statusToWaitFor = waitAckChangeManagement
			}
			if err = waitForChangeStatus(ctx, client, enrollmentID, changeID, statusToWaitFor); err != nil {
				return diag.FromErr(err)
			}

		}
	}

	changeHistory, err := client.GetChangeHistory(ctx, cps.GetChangeHistoryRequest{
		EnrollmentID: enrollmentID,
	})
	if err != nil {
		return diag.Errorf("could not get change history: %s", err)
	}
	attrs = createAttrsFromChangeHistory(changeHistory)

	if err = tf.SetAttrs(d, attrs); err != nil {
		return diag.Errorf("could not set attributes: %s", err)
	}

	return nil
}

func resourceCPSUploadCertificateUpdate(ctx context.Context, d *schema.ResourceData, m interface{}) diag.Diagnostics {
	meta := meta.Must(m)
	logger := meta.Log("CPS", "resourceCPSUploadCertificateUpdate")
	ctx = session.ContextWithOptions(ctx, session.WithContextLog(logger))
	client := inst.Client(meta)
	logger.Debug("Updating upload certificate")

	enrollmentID, err := tf.GetIntValue("enrollment_id", d)
	if err != nil {
		return diag.Errorf("could not get `enrollment_id` attribute: %s", err)
	}

	enrollment, err := client.GetEnrollment(ctx, cps.GetEnrollmentRequest{EnrollmentID: enrollmentID})
	if err != nil {
		return diag.Errorf("could not get an enrollment: %s", err)
	}

	if !d.HasChanges("certificate_ecdsa_pem", "trust_chain_ecdsa_pem", "certificate_rsa_pem", "trust_chain_rsa_pem") {
		logger.Debug("Certificate does not have to be updated.")

		if len(enrollment.PendingChanges) == 0 {
			logger.Warn("There are no pending changes, the certificate is already deployed. Changing only local state.")
			return nil
		}

		changeID, err := toolsCPS.GetChangeIDFromPendingChanges(enrollment.PendingChanges)
		if err != nil {
			return diag.Errorf("could not get changeID: %s", err)
		}

		changeStatus, err := sendGetChangeStatusReq(ctx, client, enrollmentID, changeID)
		if err != nil {
			return diag.FromErr(err)
		}

		if changeStatus.StatusInfo.Status == waitReviewThirdPartyCert || changeStatus.StatusInfo.Status == waitUploadThirdParty {
			return upsertUploadCertificate(ctx, d, m, client, logger)
		}

		if d.HasChanges("acknowledge_post_verification_warnings") || d.HasChanges("auto_approve_warnings") {
			logger.Warn("Post-verification warnings has either already been accepted or didn't occur - ignoring change in this flag.")
		}

		if d.HasChanges("acknowledge_change_management") && enrollment.ChangeManagement {
			ackChangeManagement, err := tf.GetBoolValue("acknowledge_change_management", d)
			if err != nil {
				return diag.Errorf("could not get `acknowledge_change_management` attribute: %s", err)
			}

			if !ackChangeManagement {
				logger.Warn("The certificate is either already on production network or is scheduled to be deployed. " +
					"Change in change-management flag won't be applied in this enrollment. " +
					"To apply it, remove enrollment and create a new one.")
				return nil
			}

			if err = waitForChangeStatus(ctx, client, enrollmentID, changeID, waitAckChangeManagement); err != nil {
				return diag.FromErr(err)
			}
			if err = sendACKChangeManagement(ctx, client, enrollmentID, changeID); err != nil {
				return diag.Errorf("could not acknowledge change management: %s", err)
			}
		}
		return resourceCPSUploadCertificateRead(ctx, d, m)
	}

	if len(enrollment.PendingChanges) != 0 {
		changeID, err := toolsCPS.GetChangeIDFromPendingChanges(enrollment.PendingChanges)
		if err != nil {
			return diag.Errorf("could not get changeID: %s", err)
		}
		changeStatus, err := sendGetChangeStatusReq(ctx, client, enrollmentID, changeID)
		if err != nil {
			return diag.FromErr(err)
		}

		if changeStatus.StatusInfo.Status == waitUploadThirdParty || changeStatus.StatusInfo.Status == waitReviewThirdPartyCert {
			return upsertUploadCertificate(ctx, d, m, client, logger)
		}

		return diag.Errorf("cannot make changes to the certificate with current status: %s", changeStatus.StatusInfo.Status)
	}

	return diag.Errorf("cannot make changes to certificate that is already on staging and/or production network, need to create new enrollment")
}

func resourceCPSUploadCertificateDelete(_ context.Context, _ *schema.ResourceData, m interface{}) diag.Diagnostics {
	meta := meta.Must(m)
	logger := meta.Log("CPS", "resourceCPSUploadCertificateDelete")
	logger.Debug("Deleting CPS upload certificate configuration")
	logger.Info("CPS upload certificate deletion - resource will only be removed from local state")
	return nil
}

func upsertUploadCertificate(ctx context.Context, d *schema.ResourceData, m interface{}, client cps.CPS, logger log.Interface) diag.Diagnostics {
	attrs, err := getCPSUploadCertificateAttrs(d)
	if err != nil {
		return diag.FromErr(err)
	}

	err = checkForTrustChainWithoutCert(attrs)
	if err != nil {
		return diag.FromErr(err)
	}

	enrollment, err := client.GetEnrollment(ctx, cps.GetEnrollmentRequest{EnrollmentID: attrs.enrollmentID})
	if err != nil {
		return diag.Errorf("could not get an enrollment: %s", err)
	}
	changeID, err := toolsCPS.GetChangeIDFromPendingChanges(enrollment.PendingChanges)
	if err != nil {
		return diag.Errorf("could not get change ID: %s", err)
	} else if err != nil && errors.Is(err, toolsCPS.ErrNoPendingChanges) {
		return diag.Errorf("provided enrollment has no pending changes")
	}

	certs := wrapCertificatesToUpload(attrs.certificateECDSA, attrs.trustChainECDSA, attrs.certificateRSA, attrs.trustChainRSA)
	uploadCertAndTrustChainReq := cps.UploadThirdPartyCertAndTrustChainRequest{
		EnrollmentID: attrs.enrollmentID,
		ChangeID:     changeID,
		Certificates: cps.ThirdPartyCertificates{CertificatesAndTrustChains: certs},
	}

	if err = client.UploadThirdPartyCertAndTrustChain(ctx, uploadCertAndTrustChainReq); err != nil {
		return diag.Errorf("could not upload third party certificate and trust chain: %s", err)
	}

	status, err := waitUntilStatusPasses(ctx, client, attrs.enrollmentID, changeID, verifyThirdPartyCert)
	if err != nil {
		return diag.Errorf("incorrect status of a change: %s", err)
	}

	if status == waitReviewThirdPartyCert {
		if err = processPostVerificationWarnings(ctx, client, d, attrs.enrollmentID, changeID, logger); err != nil {
			return diag.Errorf("could not process post verification warnings: %s", err)
		}
	}

	if enrollment.ChangeManagement && (attrs.ackChangeManagement || attrs.waitForDeployment) {
		if err = waitForChangeStatus(ctx, client, attrs.enrollmentID, changeID, waitAckChangeManagement); err != nil {
			return diag.FromErr(err)
		}

		if attrs.ackChangeManagement {
			if err = sendACKChangeManagement(ctx, client, attrs.enrollmentID, changeID); err != nil {
				return diag.Errorf("could not acknowledge change management: %s", err)
			}
		}
	}
	d.SetId(strconv.Itoa(attrs.enrollmentID))

	return resourceCPSUploadCertificateRead(ctx, d, m)
}

// checkUnacknowledgedWarnings checks if there are unacknowledged warnings for a certificate
func checkUnacknowledgedWarnings(ctx context.Context, diff *schema.ResourceDiff, m interface{}) error {
	meta := meta.Must(m)
	logger := meta.Log("CPS", "checkUnacknowledgedWarnings")
	ctx = session.ContextWithOptions(ctx, session.WithContextLog(logger))
	client := inst.Client(meta)
	logger.Debug("Checking for unacknowledged warnings")

	id := diff.Get("enrollment_id")
	enrollmentID, ok := id.(int)
	if !ok {
		logger.Warnf("expected enrollmentID of type int, got: %T", enrollmentID)
		return nil
	}

	// in case of the variable dependency, enrollmentID might be of '0' value when this function is first invoked, hence it should proceed
	// further and load the variable correctly (after dependent resource is created) in the Create operation
	if enrollmentID == 0 {
		return nil
	}

	enrollment, err := client.GetEnrollment(ctx, cps.GetEnrollmentRequest{EnrollmentID: enrollmentID})
	if err != nil {
		return fmt.Errorf("could not get an enrollment: %s", err)
	}

	changeID, err := toolsCPS.GetChangeIDFromPendingChanges(enrollment.PendingChanges)
	if err != nil && !errors.Is(err, toolsCPS.ErrNoPendingChanges) {
		return fmt.Errorf("could not get changeID of an enrollment: %s", err)
	} else if err != nil && errors.Is(err, toolsCPS.ErrNoPendingChanges) {
		return nil
	}

	changeStatus, err := sendGetChangeStatusReq(ctx, client, enrollmentID, changeID)
	if err != nil {
		return err
	}
	if changeStatus.StatusInfo.Status == waitReviewThirdPartyCert || changeStatus.StatusInfo.Status == waitUploadThirdParty {
		if err := diff.SetNewComputed("unacknowledged_warnings"); err != nil {
			return fmt.Errorf("could not set 'unacknowledged_warnings' attribute: %s", err)
		}
	}

	return nil
}

// checkForTrustChainWithoutCert validates if user provided trustChain without certificate and fails processing if so
func checkForTrustChainWithoutCert(attrs *attributes) error {
	if attrs.certificateRSA == "" && attrs.trustChainRSA != "" {
		return fmt.Errorf("provided RSA trust chain without RSA certificate. Please remove it or add a certificate")
	}
	if attrs.certificateECDSA == "" && attrs.trustChainECDSA != "" {
		return fmt.Errorf("provided ECDSA trust chain without ECDSA certificate. Please remove it or add a certificate")
	}

	return nil
}

// waitForChangeStatus waits for provided status
func waitForChangeStatus(ctx context.Context, client cps.CPS, enrollmentID, changeID int, status string) error {
	change, err := sendGetChangeStatusReq(ctx, client, enrollmentID, changeID)
	if err != nil {
		return fmt.Errorf("could not get change status: %s", err)
	}

	for change.StatusInfo.Status != status {
		select {
		case <-time.After(PollForChangeStatusInterval):
			change, err = sendGetChangeStatusReq(ctx, client, enrollmentID, changeID)
			if err != nil {
				return fmt.Errorf("could not get change status: %s", err)
			}
			if change.StatusInfo.Status == status {
				continue
			}
		case <-ctx.Done():
			return fmt.Errorf("retry timeout reached: incorrect status of a change: %s, %s", change.StatusInfo.Status, ctx.Err())
		}
	}

	return nil
}

// waitUntilStatusPasses waits until the status provided as parameter passes and returns a new one
func waitUntilStatusPasses(ctx context.Context, client cps.CPS, enrollmentID, changeID int, status string) (string, error) {
	change, err := sendGetChangeStatusReq(ctx, client, enrollmentID, changeID)
	if err != nil {
		return "", fmt.Errorf("could not get change status: %s", err)
	}

	for change.StatusInfo.Status == status {
		select {
		case <-time.After(PollForChangeStatusInterval):
			change, err = sendGetChangeStatusReq(ctx, client, enrollmentID, changeID)
			if err != nil {
				return "", fmt.Errorf("could not get change status: %s", err)
			}
		case <-ctx.Done():
			return "", fmt.Errorf("retry timeout reached: incorrect status of a change: %s, %s", change.StatusInfo.Status, ctx.Err())
		}
	}

	return change.StatusInfo.Status, nil
}

// wrapCertificatesToUpload creates certificates entry used in UploadThirdPartyCertAndTrustChain request,
// depending on number of certificates provided
func wrapCertificatesToUpload(certificateECDSA, trustChainECDSA, certificateRSA, trustChainRSA string) []cps.CertificateAndTrustChain {
	var certificates []cps.CertificateAndTrustChain
	if certificateECDSA != "" {
		certificates = append(certificates, cps.CertificateAndTrustChain{
			Certificate:  certificateECDSA,
			TrustChain:   trustChainECDSA,
			KeyAlgorithm: "ECDSA",
		})
	}
	if certificateRSA != "" {
		certificates = append(certificates, cps.CertificateAndTrustChain{
			Certificate:  certificateRSA,
			TrustChain:   trustChainRSA,
			KeyAlgorithm: "RSA",
		})
	}

	return certificates
}

// processPostVerificationWarnings is responsible for comparison of user-accepted warnings and required warnings
func processPostVerificationWarnings(ctx context.Context, client cps.CPS, d *schema.ResourceData, enrollmentID, changeID int, logger log.Interface) error {
	warnings, err := client.GetChangePostVerificationWarnings(ctx, cps.GetChangeRequest{
		EnrollmentID: enrollmentID,
		ChangeID:     changeID,
	})
	if err != nil {
		return fmt.Errorf("could not get post verification warnings: %s", err)
	}

	logger.Debugf("Post-verification warnings: %s", warnings.Warnings)

	acceptAllWarnings, err := tf.GetBoolValue("acknowledge_post_verification_warnings", d)
	if err != nil && !errors.Is(err, tf.ErrNotFound) {
		return fmt.Errorf("could not get `acknowledge_post_verification_warnings` attribute: %s", err)
	}

	autoApproveWarnings, err := tf.GetSetValue("auto_approve_warnings", d)
	if err != nil && !errors.Is(err, tf.ErrNotFound) {
		return fmt.Errorf("could not get `auto_approve_warnings` attribute: %s", err)
	}

	userWarningsString := convertUserWarningsToStringSlice(autoApproveWarnings.List())
	canAccept, err := canApproveWarnings(userWarningsString, warnings.Warnings)

	if len(warnings.Warnings) != 0 {
		if acceptAllWarnings || canAccept {
			if err = sendACKPostVerificationWarnings(ctx, client, enrollmentID, changeID); err != nil {
				return fmt.Errorf("could not acknowledge post verification warnings: %s", err)
			}
			_, err := waitUntilStatusPasses(ctx, client, enrollmentID, changeID, waitReviewThirdPartyCert)
			if err != nil {
				return fmt.Errorf("status %s did not pass: %s", waitReviewThirdPartyCert, err)
			}
		} else {
			return fmt.Errorf("not every warning has been acknowledged: %s", err)
		}
	}

	return nil
}

// convertUserWarningsToStringSlice converts user-provided slice of type `[]interface{}` to slice of type `[]string`
func convertUserWarningsToStringSlice(userWarnings []interface{}) []string {
	userWarningsString := make([]string, len(userWarnings))
	for i := range userWarnings {
		userWarningsString[i] = userWarnings[i].(string)
	}

	return userWarningsString
}

// sendGetChangeStatusReq creates and sends GetChangeStatus request
func sendGetChangeStatusReq(ctx context.Context, client cps.CPS, enrollmentID, changeID int) (*cps.Change, error) {
	status, err := client.GetChangeStatus(ctx, cps.GetChangeStatusRequest{
		EnrollmentID: enrollmentID,
		ChangeID:     changeID,
	})
	if err != nil {
		return nil, err
	}

	return status, nil
}

// sendACKPostVerificationWarnings creates and sends AcknowledgePostVerificationWarnings request
func sendACKPostVerificationWarnings(ctx context.Context, client cps.CPS, enrollmentID, changeID int) error {
	acknowledgementReq := cps.AcknowledgementRequest{
		Acknowledgement: cps.Acknowledgement{Acknowledgement: cps.AcknowledgementAcknowledge},
		EnrollmentID:    enrollmentID,
		ChangeID:        changeID}
	if err := client.AcknowledgePostVerificationWarnings(ctx, acknowledgementReq); err != nil {
		return fmt.Errorf("could not acknowledge post verification warnings: %s", err)
	}

	return nil
}

// sendACKChangeManagement creates and sends AcknowledgeChangeManagement request
func sendACKChangeManagement(ctx context.Context, client cps.CPS, enrollmentID, changeID int) error {
	changeAcknowledgementReq := cps.AcknowledgementRequest{
		Acknowledgement: cps.Acknowledgement{Acknowledgement: cps.AcknowledgementAcknowledge},
		EnrollmentID:    enrollmentID,
		ChangeID:        changeID,
	}
	if err := client.AcknowledgeChangeManagement(ctx, changeAcknowledgementReq); err != nil {
		return fmt.Errorf("could not acknowledge change management: %s", err)
	}

	return nil
}

// createAttrsFromChangeHistory creates attributes for a resource form GetChangeHistoryResponse
func createAttrsFromChangeHistory(changeHistory *cps.GetChangeHistoryResponse) map[string]interface{} {
	var certificateECDSA, certificateRSA, trustChainECDSA, trustChainRSA string
	if len(changeHistory.Changes) != 0 {
		for _, change := range changeHistory.Changes {
			if change.PrimaryCertificate.Certificate != "" {
				if change.PrimaryCertificate.KeyAlgorithm == "RSA" {
					certificateRSA = change.PrimaryCertificate.Certificate
					trustChainRSA = change.PrimaryCertificate.TrustChain
				} else {
					certificateECDSA = change.PrimaryCertificate.Certificate
					trustChainECDSA = change.PrimaryCertificate.TrustChain
				}
				if len(change.MultiStackedCertificates) != 0 {
					if change.MultiStackedCertificates[0].KeyAlgorithm == "RSA" {
						certificateRSA = change.MultiStackedCertificates[0].Certificate
						trustChainRSA = change.MultiStackedCertificates[0].TrustChain
					} else {
						certificateECDSA = change.MultiStackedCertificates[0].Certificate
						trustChainECDSA = change.MultiStackedCertificates[0].TrustChain
					}
				}
				break
			}
		}
	}

	attrs := make(map[string]interface{})
	attrs["certificate_ecdsa_pem"] = certificateECDSA
	attrs["trust_chain_ecdsa_pem"] = trustChainECDSA
	attrs["certificate_rsa_pem"] = certificateRSA
	attrs["trust_chain_rsa_pem"] = trustChainRSA

	return attrs
}

// getCPSUploadCertificateAttrs returns struct holding attributes from schema
func getCPSUploadCertificateAttrs(d *schema.ResourceData) (*attributes, error) {
	enrollmentID, err := tf.GetIntValue("enrollment_id", d)
	if err != nil {
		return nil, fmt.Errorf("could not get `enrollment_id` attribute: %s", err)
	}

	certificateECDSA, err := tf.GetStringValue("certificate_ecdsa_pem", d)
	if err != nil && !errors.Is(err, tf.ErrNotFound) {
		return nil, fmt.Errorf("could not get `certificate_ecdsa_pem` attribute: %s", err)
	}
	certificateRSA, err := tf.GetStringValue("certificate_rsa_pem", d)
	if err != nil && !errors.Is(err, tf.ErrNotFound) {
		return nil, fmt.Errorf("could not get `certificate_rsa_pem` attribute: %s", err)
	}

	trustChainECDSA, err := tf.GetStringValue("trust_chain_ecdsa_pem", d)
	if err != nil && !errors.Is(err, tf.ErrNotFound) {
		return nil, fmt.Errorf("could not get `trust_chain_ecdsa_pem` attribute: %s", err)
	}
	trustChainRSA, err := tf.GetStringValue("trust_chain_rsa_pem", d)
	if err != nil && !errors.Is(err, tf.ErrNotFound) {
		return nil, fmt.Errorf("could not get `trust_chain_rsa_pem` attribute: %s", err)
	}

	ackChangeManagement, err := tf.GetBoolValue("acknowledge_change_management", d)
	if err != nil && !errors.Is(err, tf.ErrNotFound) {
		return nil, fmt.Errorf("could not get `acknowledge_change_management` attribute: %s", err)
	}

	waitForDeployment, err := tf.GetBoolValue("wait_for_deployment", d)
	if err != nil && !errors.Is(err, tf.ErrNotFound) {
		return nil, fmt.Errorf("could not get `wait_for_deployment` attribute: %s", err)
	}

	return &attributes{
		enrollmentID:        enrollmentID,
		certificateECDSA:    certificateECDSA,
		certificateRSA:      certificateRSA,
		trustChainECDSA:     trustChainECDSA,
		trustChainRSA:       trustChainRSA,
		ackChangeManagement: ackChangeManagement,
		waitForDeployment:   waitForDeployment,
	}, nil
}

func resourceCPSUploadCertificateImport(ctx context.Context, d *schema.ResourceData, m interface{}) ([]*schema.ResourceData, error) {
	meta := meta.Must(m)
	logger := meta.Log("CPS", "resourceCPSUploadCertificateImport")
	ctx = session.ContextWithOptions(
		ctx,
		session.WithContextLog(logger),
	)
	logger.Debug("Importing upload certificate")
	id := d.Id()
	if id == "" {
		return nil, fmt.Errorf("didn't provide enrollment id")
	}
	enrollmentID, err := strconv.Atoi(id)
	if err != nil {
		return nil, fmt.Errorf("enrollment ID must be a number: %s", err)
	}

	client := inst.Client(meta)

	enrollment, err := client.GetEnrollment(ctx, cps.GetEnrollmentRequest{EnrollmentID: enrollmentID})
	if err != nil {
		return nil, fmt.Errorf("unable to fetch enrollment: %s", err)
	}
	if enrollment.ValidationType != "third-party" {
		return nil, fmt.Errorf("unable to import: wrong validation type: expected 'third-party', got '%s'", enrollment.ValidationType)
	}

	changeHistory, err := client.GetChangeHistory(ctx, cps.GetChangeHistoryRequest{EnrollmentID: enrollmentID})
	if err != nil {
		return nil, fmt.Errorf("unable to fetch certificates upload history: %s", err)
	}

	attrs := createAttrsFromChangeHistory(changeHistory)
	certECDSA := attrs["certificate_ecdsa_pem"]
	certRSA := attrs["certificate_rsa_pem"]
	if certRSA == "" && certECDSA == "" {
		return nil, fmt.Errorf("no certificate was yet uploaded")
	}
	attrs["acknowledge_post_verification_warnings"] = false
	attrs["auto_approve_warnings"] = []string{}
	attrs["acknowledge_change_management"] = false
	attrs["wait_for_deployment"] = false
	attrs["enrollment_id"] = enrollmentID

	if err = tf.SetAttrs(d, attrs); err != nil {
		return nil, fmt.Errorf("could not set attributes: %s", err)
	}

	return []*schema.ResourceData{d}, nil
}
