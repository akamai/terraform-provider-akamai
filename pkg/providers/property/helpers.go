package property

import (
	"context"
	"errors"
	"fmt"
	"net/http"
	"strings"
	"time"

	"github.com/akamai/AkamaiOPEN-edgegrid-golang/v12/pkg/iam"
	"github.com/akamai/AkamaiOPEN-edgegrid-golang/v12/pkg/log"
	"github.com/akamai/AkamaiOPEN-edgegrid-golang/v12/pkg/papi"
	"github.com/akamai/terraform-provider-akamai/v9/pkg/common/str"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/schema"
)

var certStatus = &schema.Resource{
	Schema: map[string]*schema.Schema{
		"target": {
			Type:        schema.TypeString,
			Computed:    true,
			Description: "The destination part of the CNAME record used to validate the certificate's domain.",
		},
		"hostname": {
			Type:        schema.TypeString,
			Computed:    true,
			Description: "The hostname part of the CNAME record used to validate the certificate's domain.",
		},
		"production_status": {
			Type:        schema.TypeString,
			Computed:    true,
			Description: "The certificate's deployment status on the production network.",
		},
		"staging_status": {
			Type:        schema.TypeString,
			Computed:    true,
			Description: "The certificate's deployment status on the staging network.",
		},
	},
}

var complianceRecordSchema = &schema.Resource{
	Schema: map[string]*schema.Schema{
		"noncompliance_reason_none": {
			Type:     schema.TypeList,
			Optional: true,
			MaxItems: 1,
			ExactlyOneOf: []string{
				"compliance_record.0.noncompliance_reason_other",
				"compliance_record.0.noncompliance_reason_no_production_traffic",
				"compliance_record.0.noncompliance_reason_emergency",
			},
			Description: fmt.Sprintf("Provides an audit record when activating on a production network with noncompliance reason as `%s`", papi.NoncomplianceReasonNone),
			Elem: &schema.Resource{
				Schema: map[string]*schema.Schema{
					"ticket_id": {
						Type:        schema.TypeString,
						Optional:    true,
						Description: "Identifies the ticket that describes the need for the activation",
					},
					"customer_email": {
						Type:        schema.TypeString,
						Optional:    true,
						Description: "Identifies the customer",
					},
					"peer_reviewed_by": {
						Type:        schema.TypeString,
						Optional:    true,
						Description: "Identifies person who has independently approved the activation request",
					},
					"unit_tested": {
						Type:        schema.TypeBool,
						Optional:    true,
						Description: "Whether the metadata to activate has been fully tested",
					},
				},
			},
		},
		"noncompliance_reason_other": {
			Type:        schema.TypeList,
			Optional:    true,
			MaxItems:    1,
			Description: fmt.Sprintf("Provides an audit record when activating on a production network with noncompliance reason as `%s`", papi.NoncomplianceReasonOther),
			Elem: &schema.Resource{
				Schema: map[string]*schema.Schema{
					"ticket_id": {
						Type:        schema.TypeString,
						Optional:    true,
						Description: "Identifies the ticket that describes the need for the activation",
					},
					"other_noncompliance_reason": {
						Type:        schema.TypeString,
						Optional:    true,
						Description: "Describes the reason why the activation must occur immediately, out of compliance with the standard procedure",
					},
				},
			},
		},
		"noncompliance_reason_no_production_traffic": {
			Type:        schema.TypeList,
			Optional:    true,
			MaxItems:    1,
			Description: fmt.Sprintf("Provides an audit record when activating on a production network with noncompliance reason as `%s`", papi.NoncomplianceReasonNoProductionTraffic),
			Elem: &schema.Resource{
				Schema: map[string]*schema.Schema{
					"ticket_id": {
						Type:        schema.TypeString,
						Optional:    true,
						Description: "Identifies the ticket that describes the need for the activation",
					},
				},
			},
		},
		"noncompliance_reason_emergency": {
			Type:        schema.TypeList,
			Optional:    true,
			MaxItems:    1,
			Description: fmt.Sprintf("Provides an audit record when activating on a production network with noncompliance reason as `%s`", papi.NoncomplianceReasonEmergency),
			Elem: &schema.Resource{
				Schema: map[string]*schema.Schema{
					"ticket_id": {
						Type:        schema.TypeString,
						Optional:    true,
						Description: "Identifies the ticket that describes the need for the activation",
					},
				},
			},
		},
	},
}

var ccmCertificatesSchema = &schema.Resource{
	Schema: map[string]*schema.Schema{
		"rsa_cert_id": {
			Type:        schema.TypeString,
			Optional:    true,
			Description: "Certificate ID for RSA.",
		},
		"ecdsa_cert_id": {
			Type:        schema.TypeString,
			Optional:    true,
			Description: "Certificate ID for ECDSA.",
		},
	},
}

var ccmCertificateStatusSchema = &schema.Resource{
	Schema: map[string]*schema.Schema{
		"rsa_staging_status": {
			Type:        schema.TypeString,
			Computed:    true,
			Description: "Status of the RSA certificate on staging network.",
		},
		"rsa_production_status": {
			Type:        schema.TypeString,
			Computed:    true,
			Description: "Status of the RSA certificate on production network.",
		},
		"ecdsa_staging_status": {
			Type:        schema.TypeString,
			Computed:    true,
			Description: "Status of the ECDSA certificate on staging network.",
		},
		"ecdsa_production_status": {
			Type:        schema.TypeString,
			Computed:    true,
			Description: "Status of the ECDSA certificate on production network.",
		},
	},
}

// Convert given hostnames to the map form that can be stored in a schema.ResourceData
// Setting only statuses for default certs if they exist
func flattenHostnames(Hostnames []papi.Hostname) []map[string]interface{} {
	var res []map[string]interface{}
	for _, hn := range Hostnames {
		m := map[string]interface{}{}
		m["cname_from"] = hn.CnameFrom
		m["cname_to"] = hn.CnameTo
		m["cert_provisioning_type"] = hn.CertProvisioningType
		m["edge_hostname_id"] = hn.EdgeHostnameID
		m["cname_type"] = hn.CnameType
		m["cert_status"] = []map[string]any{flattenCertType(&hn.CertStatus)}
		res = append(res, m)
	}
	return res
}

// TODO: remove this when updating akamai_property_hostnames datasource
// and use flattenHostnames instead
func flattenHostnamesCCM(Hostnames []papi.Hostname) []map[string]interface{} {
	var res []map[string]interface{}
	for _, hn := range Hostnames {
		m := map[string]interface{}{}
		m["cname_from"] = hn.CnameFrom
		m["cname_to"] = hn.CnameTo
		m["cert_provisioning_type"] = hn.CertProvisioningType
		m["edge_hostname_id"] = hn.EdgeHostnameID
		m["cname_type"] = hn.CnameType
		m["cert_status"] = []map[string]any{flattenCertType(&hn.CertStatus)}
		m["ccm_certificates"] = flattenCCMCertificates(hn.CCMCertificates)
		m["ccm_cert_status"] = flattenCCMCertificateStatus(hn.CCMCertStatus)
		res = append(res, m)
	}
	return res
}

func flattenCCMCertificateStatus(status *papi.CCMCertStatus) []map[string]string {
	if status == nil {
		return nil
	}
	m := map[string]string{}
	m["rsa_staging_status"] = status.RSAStagingStatus
	m["rsa_production_status"] = status.RSAProductionStatus
	m["ecdsa_staging_status"] = status.ECDSAStagingStatus
	m["ecdsa_production_status"] = status.ECDSAProductionStatus

	return []map[string]string{m}
}

func flattenCCMCertificates(certificates *papi.CCMCertificates) []map[string]string {
	if certificates == nil {
		return nil
	}
	m := map[string]string{}
	m["rsa_cert_id"] = certificates.RSACertID
	m["ecdsa_cert_id"] = certificates.ECDSACertID
	return []map[string]string{m}
}

func flattenBucketHostnames(hostnames []papi.HostnameItem) []map[string]any {
	var result []map[string]any
	for _, hostname := range hostnames {
		result = append(result, map[string]any{
			"cname_from":                  hostname.CnameFrom,
			"cname_type":                  hostname.CnameType,
			"staging_edge_hostname_id":    hostname.StagingEdgeHostnameID,
			"staging_cert_type":           hostname.StagingCertType,
			"staging_cname_to":            hostname.StagingCnameTo,
			"production_edge_hostname_id": hostname.ProductionEdgeHostnameID,
			"production_cert_type":        hostname.ProductionCertType,
			"production_cname_to":         hostname.ProductionCnameTo,
			"cert_status":                 []map[string]any{flattenCertType(hostname.CertStatus)},
		})
	}
	return result
}

func flattenCertType(certStatus *papi.CertStatusItem) map[string]any {
	if certStatus == nil {
		return nil
	}

	certs := map[string]any{
		"hostname": certStatus.ValidationCname.Hostname,
		"target":   certStatus.ValidationCname.Target,
	}
	if len(certStatus.Staging) > 0 {
		certs["staging_status"] = certStatus.Staging[0].Status
	}
	if len(certStatus.Production) > 0 {
		certs["production_status"] = certStatus.Production[0].Status
	}
	return certs
}

func papiErrorsToList(errors []*papi.Error) []map[string]interface{} {
	if len(errors) == 0 {
		return nil
	}

	var ruleErrors []map[string]interface{}
	for _, err := range errors {
		if err == nil {
			continue
		}

		ruleErrors = append(ruleErrors, papiErrorToMap(err))
	}

	return ruleErrors
}

func papiErrorToMap(err *papi.Error) map[string]interface{} {
	if err == nil {
		return nil
	}

	return map[string]interface{}{
		"type":           err.Type,
		"title":          err.Title,
		"detail":         err.Detail,
		"instance":       err.Instance,
		"behavior_name":  err.BehaviorName,
		"error_location": err.ErrorLocation,
		"status_code":    err.StatusCode,
	}
}

// NetworkAlias parses the given network name or alias and returns its full name and any error
func NetworkAlias(network string) (string, error) {
	networks := map[string]papi.ActivationNetwork{
		"STAGING":    papi.ActivationNetworkStaging,
		"STAGE":      papi.ActivationNetworkStaging,
		"STAG":       papi.ActivationNetworkStaging,
		"S":          papi.ActivationNetworkStaging,
		"PRODUCTION": papi.ActivationNetworkProduction,
		"PROD":       papi.ActivationNetworkProduction,
		"P":          papi.ActivationNetworkProduction,
	}

	networkValue, ok := networks[strings.ToUpper(network)]
	if !ok {
		return "", errors.New("network not recognized")
	}
	return string(networkValue), nil
}

func areGroupIDsDifferent(firstGroupID, secondGroupID string) (bool, error) {
	gid1, err := str.GetIntID(firstGroupID, "grp_")
	if err != nil {
		return false, err
	}

	gid2, err := str.GetIntID(secondGroupID, "grp_")
	if err != nil {
		return false, err
	}

	return gid1 != gid2, nil
}

type papiKey struct {
	propertyID string
	groupID    string
	contractID string
}

type helper struct {
	client    papi.PAPI
	iamClient iam.IAM
}

// moveProperty changes the group of the property specified by key and assetID to the group
// with destGroupID.
//
// If the property is already in the desired group (e.g. it was changed earlier in the same
// configuration from property bootstrap), nothing happens and no error is reported.
// If the property has never been activated, an error is returned (see validatePropertyMove).
// After a successful move, this method polls the API until the new group is returned by a
// property read endpoint.
func (h helper) moveProperty(ctx context.Context, key papiKey, assetID, destGroupID string) error {
	logger := log.FromContext(ctx).With("", log.Fields{
		"key":         key,
		"assetID":     assetID,
		"destGroupID": destGroupID,
	})
	logger.Debug("moveProperty")
	ctx = log.NewContext(ctx, logger)

	from, err := str.GetIntID(key.groupID, "grp_")
	if err != nil {
		return fmt.Errorf("error parsing src group id: %w", err)
	}
	to, err := str.GetIntID(destGroupID, "grp_")
	if err != nil {
		return fmt.Errorf("error parsing dst group id: %w", err)
	}
	iamID, err := str.GetIntID(assetID, "aid_")
	if err != nil {
		return fmt.Errorf("error parsing asset id: %w", err)
	}

	done, err := h.isPropertyInGroup(ctx, papiKey{
		propertyID: key.propertyID,
		groupID:    destGroupID,
		contractID: key.contractID,
	})
	if err != nil {
		return fmt.Errorf("error checking if property in group: %w", err)
	}
	if done {
		logger.Debugf("Changing group id from %s to %s: skipping, group already changed",
			key.groupID, destGroupID)
		return nil
	}

	if err := h.validatePropertyMove(ctx, key); err != nil {
		return err
	}

	logger.Debugf("Changing group id from %d to %d for IAM id %d", from, to, iamID)
	err = h.iamClient.MoveProperty(ctx, iam.MovePropertyRequest{
		PropertyID: int64(iamID),
		Body: iam.MovePropertyRequestBody{
			DestinationGroupID: int64(to),
			SourceGroupID:      int64(from),
		},
	})
	if err != nil {
		return fmt.Errorf("error calling move property API: %w", err)
	}

	err = h.waitForPropertyGroupIDChange(ctx, papiKey{
		propertyID: key.propertyID,
		groupID:    destGroupID,
		contractID: key.contractID,
	}, 5, time.Second)
	if err != nil {
		return fmt.Errorf("error waiting for group id change: %w", err)
	}
	return nil
}

// isPropertyInGroup checks whether the property specified with key.propertyID and key.contractID
// is in group key.groupID.
func (h helper) isPropertyInGroup(ctx context.Context, key papiKey) (bool, error) {

	logger := log.FromContext(ctx).With("", log.Fields{
		"key": key,
	})

	prp, err := fetchLatestProperty(ctx, h.client, key.propertyID, key.groupID, key.contractID)
	if err != nil {
		if !isHTTP403(err) {
			return false, fmt.Errorf("unexpected http error for %s: %w", key, err)
		}
		// No such property in such group
		msg := fmt.Sprintf("no such property in group %s: HTTP 403 received", key.groupID)
		logger.Debug(msg, "error", err)
		return false, nil
	}

	// It is possible that the property was in key.groupID in the past and Open API still returns
	// a valid response for it. To be sure, we need to check prp.GroupID which is the actual group id.
	diff, err := areGroupIDsDifferent(key.groupID, prp.GroupID)
	if err != nil {
		return false, err
	}
	if diff {
		logger.Debugf("fetched property has group id %s different than expected %s",
			prp.GroupID, key.groupID)
		return false, nil
	}
	return true, nil
}

// validatePropertyMove returns error when the property specified by key is in a state where it
// cannot be safely moved.
//
// A property that has never been activated reports rule errors for its CP codes
// after moving to a sibling group, so we forbid moving it.
func (h helper) validatePropertyMove(ctx context.Context, key papiKey) error {
	res, err := h.client.GetActivations(ctx, papi.GetActivationsRequest{
		PropertyID: key.propertyID,
		ContractID: key.contractID,
		GroupID:    key.groupID,
	})
	if err != nil {
		return fmt.Errorf("error getting activations list for %s: %w", key, err)
	}
	if len(res.Activations.Items) == 0 {
		return fmt.Errorf("moving properties that have never been activated is not supported "+
			"(property id: %s, contract id: %s, group id %s)",
			key.propertyID, key.contractID, key.groupID)
	}
	return nil
}

// waitForPropertyGroupIDChange polls the get property endpoint until the returned property is
// in the group specified in key.
//
// This makes changing property's group id "synchronous" so that following TFP actions operate
// on a known state. The method uses binary exponential backoff with initialWait and maxAttempts.
func (h helper) waitForPropertyGroupIDChange(ctx context.Context, key papiKey, maxAttempts int, initialWait time.Duration) error {
	logger := log.FromContext(ctx).With("", log.Fields{
		"key":         key,
		"maxAttempts": maxAttempts})

	logger.Debug("waitForPropertyGroupIDChange")

	attemptsLeft := maxAttempts
	wait := initialWait
	for {
		done, err := h.isPropertyInGroup(ctx, key)
		if err != nil {
			return err
		}
		if done {
			logger.Debug("waitForPropertyGroupIDChange: success")
			return nil
		}

		attemptsLeft--
		if attemptsLeft <= 0 {
			return fmt.Errorf("waiting for groupID change to: %s for propertyID: %s, "+
				"contractID: %s in %d attempts failed",
				key.groupID, key.propertyID, key.contractID, maxAttempts)
		}
		logger.Debugf("waitForPropertyGroupIDChange: new group id still not visible, %d attempts left, "+
			"waiting %s...", attemptsLeft, wait)

		select {
		case <-ctx.Done():
			return ctx.Err()
		case <-time.After(wait):
			wait = wait * 2
		}
	}
}

func isHTTP403(err error) bool {
	var papiErr *papi.Error
	if errors.As(err, &papiErr) {
		return papiErr.StatusCode == http.StatusForbidden
	}
	return false
}
