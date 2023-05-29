package cps

import (
	"context"
	"errors"
	"fmt"
	"regexp"
	"strconv"
	"strings"
	"time"

	"github.com/akamai/AkamaiOPEN-edgegrid-golang/v6/pkg/cps"
	"github.com/akamai/AkamaiOPEN-edgegrid-golang/v6/pkg/session"
	"github.com/akamai/terraform-provider-akamai/v4/pkg/common/tf"
	"github.com/akamai/terraform-provider-akamai/v4/pkg/meta"
	cpstools "github.com/akamai/terraform-provider-akamai/v4/pkg/providers/cps/tools"
	"github.com/akamai/terraform-provider-akamai/v4/pkg/tools"
	"github.com/apex/log"
	"github.com/hashicorp/terraform-plugin-sdk/v2/diag"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/schema"
)

type (
	challengeDNS  challenge
	challengeHTTP challenge
	challenge     map[string]interface{}
)

const (
	statusCoordinateDomainValidation    = "coodinate-domain-validation"
	waitUploadThirdParty                = "wait-upload-third-party"
	statusVerificationWarnings          = "wait-review-pre-verification-safety-checks"
	inputTypePreVerificationWarningsAck = "pre-verification-warnings-acknowledgement"
	waitReviewThirdPartyCert            = "wait-review-third-party-cert"
	waitAckChangeManagement             = "wait-ack-change-management"
	complete                            = "complete"
	verifyThirdPartyCert                = "verify-third-party-cert"
)

var (
	contact = &schema.Resource{
		Schema: map[string]*schema.Schema{
			"first_name": {
				Type:        schema.TypeString,
				Required:    true,
				Description: "First name of the contact",
			},
			"last_name": {
				Type:        schema.TypeString,
				Required:    true,
				Description: "Last name of the contact",
			},
			"title": {
				Type:        schema.TypeString,
				Optional:    true,
				Description: "Title of the the contact",
			},
			"organization": {
				Type:        schema.TypeString,
				Required:    true,
				Description: "Organization where contact is hired",
			},
			"email": {
				Type:        schema.TypeString,
				Required:    true,
				Description: "E-mail address of the contact",
			},
			"phone": {
				Type:        schema.TypeString,
				Required:    true,
				Description: "Phone number of the contact",
			},
			"address_line_one": {
				Type:        schema.TypeString,
				Required:    true,
				Description: "The address of the contact",
			},
			"address_line_two": {
				Type:        schema.TypeString,
				Optional:    true,
				Description: "The address of the contact",
			},
			"city": {
				Type:        schema.TypeString,
				Required:    true,
				Description: "City of residence of the contact",
			},
			"region": {
				Type:        schema.TypeString,
				Required:    true,
				Description: "The region of the contact",
			},
			"postal_code": {
				Type:        schema.TypeString,
				Required:    true,
				Description: "Postal code of the contact",
			},
			"country_code": {
				Type:        schema.TypeString,
				Required:    true,
				Description: "Country code of the contact",
			},
		},
	}

	organization = &schema.Resource{
		Schema: map[string]*schema.Schema{
			"name": {
				Type:        schema.TypeString,
				Required:    true,
				Description: "Name of organization",
			},
			"phone": {
				Type:        schema.TypeString,
				Required:    true,
				Description: "Phone number of organization",
			},
			"address_line_one": {
				Type:        schema.TypeString,
				Required:    true,
				Description: "The address of organization",
			},
			"address_line_two": {
				Type:        schema.TypeString,
				Optional:    true,
				Description: "The address of organization",
			},
			"city": {
				Type:        schema.TypeString,
				Required:    true,
				Description: "City of organization",
			},
			"region": {
				Type:        schema.TypeString,
				Required:    true,
				Description: "The region of organization",
			},
			"postal_code": {
				Type:        schema.TypeString,
				Required:    true,
				Description: "Postal code of organization",
			},
			"country_code": {
				Type:        schema.TypeString,
				Required:    true,
				Description: "Country code of organization",
			},
		},
	}

	networkConfiguration = &schema.Resource{
		Schema: map[string]*schema.Schema{
			"client_mutual_authentication": {
				Type:        schema.TypeSet,
				Optional:    true,
				MinItems:    1,
				MaxItems:    1,
				Description: "The trust chain configuration used for client mutual authentication",
				Elem: &schema.Resource{
					Schema: map[string]*schema.Schema{
						"send_ca_list_to_client": {
							Type:        schema.TypeBool,
							Optional:    true,
							Description: "Enable the server to send the certificate authority (CA) list to the client",
						},
						"ocsp_enabled": {
							Type:        schema.TypeBool,
							Optional:    true,
							Description: "Enable OCSP stapling",
						},
						"set_id": {
							Type:        schema.TypeString,
							Optional:    true,
							Description: "The identifier of the set of trust chains, created in the Trust Chain Manager",
						},
					},
				},
			},
			"disallowed_tls_versions": {
				Type:        schema.TypeSet,
				Optional:    true,
				Elem:        &schema.Schema{Type: schema.TypeString},
				Description: "TLS versions which are disallowed",
			},
			"clone_dns_names": {
				Type:        schema.TypeBool,
				Optional:    true,
				Description: "Enable CPS to direct traffic using all the SANs listed in the SANs parameter when enrollment is created",
			},
			"geography": {
				Type:        schema.TypeString,
				Required:    true,
				Description: "Geography type used for enrollment",
			},
			"must_have_ciphers": {
				Type:        schema.TypeString,
				Optional:    true,
				Description: "Mandatory Ciphers which are included for enrollment",
			},
			"ocsp_stapling": {
				Type:        schema.TypeString,
				Optional:    true,
				Description: "Enable OCSP stapling",
			},
			"preferred_ciphers": {
				Type:        schema.TypeString,
				Optional:    true,
				Description: "Preferred Ciphers which are included for enrollment",
			},
			"quic_enabled": {
				Type:        schema.TypeBool,
				Optional:    true,
				Description: "Enable QUIC protocol",
			},
		},
	}

	csr = &schema.Resource{
		Schema: map[string]*schema.Schema{
			"country_code": {
				Type:        schema.TypeString,
				Required:    true,
				Description: "The code of the country where organization is located",
			},
			"city": {
				Type:        schema.TypeString,
				Required:    true,
				Description: "City where organization is located",
			},
			"organization": {
				Type:        schema.TypeString,
				Required:    true,
				Description: "Name of organization used in all legal documents",
			},
			"organizational_unit": {
				Type:        schema.TypeString,
				Required:    true,
				Description: "Organizational unit of organization",
			},
			"preferred_trust_chain": {
				Type:        schema.TypeString,
				Optional:    true,
				Description: "For the Let's Encrypt Domain Validated (DV) SAN certificates, the preferred trust chain will be included by CPS with the leaf certificate in the TLS handshake. If the field does not have a value, whichever trust chain Akamai chooses will be used by default",
			},
			"state": {
				Type:        schema.TypeString,
				Required:    true,
				Description: "State or province of organization location",
			},
		},
	}
)

func createAttrs(en *cps.Enrollment, enID int) map[string]interface{} {

	sans := make([]string, 0)
	for _, san := range en.CSR.SANS {
		sans = append(sans, san)
	}

	return map[string]interface{}{
		"common_name":                       en.CSR.CN,
		"enrollment_id":                     enID,
		"sans":                              sans,
		"sni_only":                          en.NetworkConfiguration.SNIOnly,
		"secure_network":                    en.NetworkConfiguration.SecureNetwork,
		"admin_contact":                     []interface{}{cpstools.ContactInfoToMap(*en.AdminContact)},
		"tech_contact":                      []interface{}{cpstools.ContactInfoToMap(*en.TechContact)},
		"certificate_chain_type":            en.CertificateChainType,
		"csr":                               []interface{}{cpstools.CSRToMap(*en.CSR)},
		"enable_multi_stacked_certificates": en.EnableMultiStackedCertificates,
		"network_configuration":             []interface{}{cpstools.NetworkConfigToMap(*en.NetworkConfiguration)},
		"signature_algorithm":               en.SignatureAlgorithm,
		"organization":                      []interface{}{cpstools.OrgToMap(*en.Org)},
		"certificate_type":                  en.CertificateType,
		"validation_type":                   en.ValidationType,
		"registration_authority":            en.RA,
	}
}

func getChallengesAttrs(ctx context.Context, en *cps.Enrollment, client cps.CPS) (map[string]interface{}, error) {
	changeID, err := cpstools.GetChangeIDFromPendingChanges(en.PendingChanges)

	if err != nil {
		if errors.Is(err, cpstools.ErrNoPendingChanges) {
			return nil, nil
		}
		return nil, err
	}
	enID, err := cpstools.GetEnrollmentID(en.Location)
	if err != nil {
		return nil, err
	}
	changeStatusReq := cps.GetChangeStatusRequest{
		EnrollmentID: enID,
		ChangeID:     changeID,
	}
	status, err := client.GetChangeStatus(ctx, changeStatusReq)
	if err != nil {
		return nil, err
	}
	if len(status.AllowedInput) < 1 || status.AllowedInput[0].Type != "lets-encrypt-challenges" {
		return nil, nil
	}

	getChallengesReq := cps.GetChangeRequest{
		EnrollmentID: enID,
		ChangeID:     changeID,
	}
	challenges, err := client.GetChangeLetsEncryptChallenges(ctx, getChallengesReq)
	if err != nil {
		return nil, err
	}

	httpChallenges, dnsChallenges := splitChallenges(challenges)
	attrs := make(map[string]interface{})
	attrs["http_challenges"] = httpChallenges
	attrs["dns_challenges"] = dnsChallenges
	return attrs, nil
}

func splitChallenges(challenges *cps.DVArray) ([]challengeHTTP, []challengeDNS) {
	dnsChallenges := make([]challengeDNS, 0)
	httpChallenges := make([]challengeHTTP, 0)

	for _, dv := range challenges.DV {
		if dv.ValidationStatus == "VALIDATED" {
			continue
		}
		for _, challenge := range dv.Challenges {
			if challenge.Status != "pending" {
				continue
			}
			if challenge.Type == "http-01" {
				httpChallenges = append(httpChallenges, challengeHTTP(newChallenge(&challenge, &dv)))
			} else if challenge.Type == "dns-01" {
				dnsChallenges = append(dnsChallenges, challengeDNS(newChallenge(&challenge, &dv)))
			}
		}
	}
	return httpChallenges, dnsChallenges
}

func newChallenge(c *cps.Challenge, dv *cps.DV) challenge {
	return challenge{
		"full_path":     c.FullPath,
		"response_body": c.ResponseBody,
		"domain":        dv.Domain,
	}
}

func enrollmentDelete(ctx context.Context, d *schema.ResourceData, m interface{}, functionName string) diag.Diagnostics {
	meta := meta.Must(m)
	logger := meta.Log("CPS", functionName)
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
	req := cps.RemoveEnrollmentRequest{
		EnrollmentID:              enrollmentID,
		AllowCancelPendingChanges: tools.BoolPtr(true),
	}
	if _, err = client.RemoveEnrollment(ctx, req); err != nil {
		return diag.FromErr(err)
	}
	return nil
}

func readAttrs(enrollment *cps.Enrollment, d *schema.ResourceData) (map[string]interface{}, error) {
	attrs := make(map[string]interface{})
	adminContact := cpstools.ContactInfoToMap(*enrollment.AdminContact)
	attrs["common_name"] = enrollment.CSR.CN
	sans := make([]string, 0)
	sansFromSchema, err := tf.GetSetValue("sans", d)
	if err != nil && !errors.Is(err, tf.ErrNotFound) {
		return nil, err
	}
	for _, san := range enrollment.CSR.SANS {
		if (sansFromSchema.Len() == 0 || !sansFromSchema.Contains(enrollment.CSR.CN)) && san == enrollment.CSR.CN {
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
	networkConfig := cpstools.NetworkConfigToMap(*enrollment.NetworkConfiguration)
	attrs["network_configuration"] = []interface{}{networkConfig}
	attrs["signature_algorithm"] = enrollment.SignatureAlgorithm
	org := cpstools.OrgToMap(*enrollment.Org)
	attrs["organization"] = []interface{}{org}
	return attrs, nil
}

func waitForVerification(ctx context.Context, logger log.Interface, client cps.CPS, enrollmentID int, acknowledgeWarnings bool, autoApproveWarnings []string) error {
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
	for ((status.StatusInfo.Status != statusCoordinateDomainValidation && status.StatusInfo.Status != waitUploadThirdParty) || len(status.AllowedInput) == 0) &&
		status.StatusInfo.Status != "complete" {
		select {
		case <-time.After(PollForChangeStatusInterval):
			status, err = client.GetChangeStatus(ctx, changeStatusReq)
			if err != nil {
				return err
			}
			if status.StatusInfo != nil && status.StatusInfo.Status == statusVerificationWarnings &&
				len(status.AllowedInput) > 0 && status.AllowedInput[0].Type == inputTypePreVerificationWarningsAck {

				warnings, err := client.GetChangePreVerificationWarnings(ctx, cps.GetChangeRequest{
					EnrollmentID: enrollmentID,
					ChangeID:     changeID,
				})
				if err != nil {
					return err
				}
				logger.Debugf("Pre-verification warnings: %s", warnings.Warnings)

				// for DV autoApproveWarnings is always empty
				if !acknowledgeWarnings && len(autoApproveWarnings) == 0 {
					return fmt.Errorf("enrollment pre-verification returned warnings and the enrollment cannot be validated. Please fix the issues or set acknowledge_pre_validation_warnings flag to true then run 'terraform apply' again: %s",
						warnings.Warnings)
				}

				canApprove, err := canApproveWarnings(autoApproveWarnings, warnings.Warnings)
				if !acknowledgeWarnings && !canApprove {
					return err
				}

				err = client.AcknowledgePreVerificationWarnings(ctx, cps.AcknowledgementRequest{
					Acknowledgement: cps.Acknowledgement{Acknowledgement: cps.AcknowledgementAcknowledge},
					EnrollmentID:    enrollmentID,
					ChangeID:        changeID,
				})
				if err != nil {
					return err
				}
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

func canApproveWarnings(autoApproveWarnings []string, warningsAsString string) (bool, error) {
	warnings, err := convertWarnings(warningsAsString)
	if err != nil {
		return false, err
	}
	autoApproveWarningsAsMap := make(map[string]bool)
	for _, warning := range autoApproveWarnings {
		autoApproveWarningsAsMap[warning] = true
	}
	var unApprovedWarnings []string
	for _, warning := range warnings {
		if _, ok := autoApproveWarningsAsMap[warning]; !ok {
			unApprovedWarnings = append(unApprovedWarnings, warning)
		}
	}
	if len(unApprovedWarnings) > 0 {
		return false, fmt.Errorf(`%w: "%s"`, ErrWarningsCannotBeApproved, strings.Join(unApprovedWarnings, `", "`))
	}
	return true, nil
}

// warnings contains entries separated with new line character. Problem is that each entry can also contain new line character.
// Another problem is that values of some keys are substrings of some other values from different key.
// All that is causing that we need to convert text into key names in a tricky way:
// 1. find beginning of warning by matching part of the warning and part of the known warning up to the new line character
// 2. what wasn't matches merge using new line with the previously found warning
// 3. now try to match the whole warning with the whole known warning and convert into warning code
// 4. if it does not work, try to split again at the new line and math with known warning
// 5. what wasn't matched goes to the unknown warning list
func convertWarnings(warnings string) ([]string, error) {
	if len(warnings) == 0 {
		return nil, nil
	}

	knownWarnings, err := convertWarningsToRegexp(warningMap)
	if err != nil {
		return nil, err
	}
	convertedWarnings := divideWarnings(warnings, knownWarnings)

	result, unknownWarnings := matchWarningToKey(convertedWarnings, knownWarnings)

	if len(unknownWarnings) > 0 {
		return nil, fmt.Errorf("received warning(s) does not match any known warning: '%s'", strings.Join(unknownWarnings, `', '`))
	}
	return result, nil
}

func matchWarningToKey(convertedWarnings []string, knownWarnings map[string]string) ([]string, []string) {
	// convert found warnings into their codes
	var result = make([]string, 0)
	var unknownWarnings = make([]string, 0)
main:
	for _, warning := range convertedWarnings {
		for k, v := range knownWarnings {
			r := regexp.MustCompile("^" + v + "$")
			if r.MatchString(warning) {
				result = append(result, k)
				continue main
			}
		}
		// try to match up to the new line
		w := strings.Split(warning, "\n")
		for k, v := range knownWarnings {
			r := regexp.MustCompile(v)
			if r.MatchString(w[0]) {
				result = append(result, k)
				unknownWarnings = append(unknownWarnings, strings.Join(w[1:], "\n"))
				continue main
			}
		}
		unknownWarnings = append(unknownWarnings, warning)
	}
	return result, unknownWarnings
}

func divideWarnings(warnings string, knownWarnings map[string]string) []string {
	warningsArray := strings.Split(warnings, "\n")
	// we are trying to match received warning and known warning matching only first part up to the new line
	singleWarning := make([]string, 0)
	convertedWarnings := make([]string, 0)
	for i, warning := range warningsArray {
		for _, v := range knownWarnings {
			vs := strings.Split(v, "\n")
			r := regexp.MustCompile(vs[0])
			if r.MatchString(warning) {
				if i > 0 {
					convertedWarnings = append(convertedWarnings, strings.Join(singleWarning, "\n"))
				}
				singleWarning = nil
				break
			}
		}
		singleWarning = append(singleWarning, warning)

	}
	convertedWarnings = append(convertedWarnings, strings.Join(singleWarning, "\n"))
	return convertedWarnings
}

func convertWarningsToRegexp(warningMap map[string]string) (map[string]string, error) {
	result := make(map[string]string, 0)
	re := regexp.MustCompile("(<.+?>)")
	for key, description := range warningMap {
		// escape some characters which have special meaning for regexp
		description = strings.ReplaceAll(description, `\`, `\\`)
		description = strings.ReplaceAll(description, "[", `\[`)
		description = strings.ReplaceAll(description, "]", `\]`)
		description = strings.ReplaceAll(description, ".", `\.`)
		description = strings.ReplaceAll(description, "(", `\(`)
		description = strings.ReplaceAll(description, ")", `\)`)
		// replace text inside diamond into regexp, so it can be matched later
		desc := re.ReplaceAllString(description, "(?s).+?")
		// verify if this is a correct regexp
		_, err := regexp.Compile(desc)
		if err != nil {
			return nil, err
		}
		result[key] = desc
	}
	return result, nil
}
