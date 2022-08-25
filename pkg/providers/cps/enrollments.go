package cps

import (
	"context"
	"errors"

	"github.com/akamai/AkamaiOPEN-edgegrid-golang/v2/pkg/cps"
	cpstools "github.com/akamai/terraform-provider-akamai/v2/pkg/providers/cps/tools"
)

type (
	challengeDNS  challenge
	challengeHTTP challenge
	challenge     map[string]interface{}
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

func newChallenge(c *cps.Challenges, dv *cps.DV) challenge {
	return challenge{
		"full_path":     c.FullPath,
		"response_body": c.ResponseBody,
		"domain":        dv.Domain,
	}
}
