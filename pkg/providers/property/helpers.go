package property

import (
	"context"
	"errors"
	"fmt"
	"net/http"
	"strings"
	"time"

	"github.com/akamai/AkamaiOPEN-edgegrid-golang/v9/pkg/iam"
	"github.com/akamai/AkamaiOPEN-edgegrid-golang/v9/pkg/papi"
	"github.com/akamai/terraform-provider-akamai/v6/pkg/common/str"
	"github.com/apex/log"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/schema"
)

var certStatus = &schema.Resource{
	Schema: map[string]*schema.Schema{
		"target": {
			Type:     schema.TypeString,
			Computed: true,
		},
		"hostname": {
			Type:     schema.TypeString,
			Computed: true,
		},
		"production_status": {
			Type:     schema.TypeString,
			Computed: true,
		},
		"staging_status": {
			Type:     schema.TypeString,
			Computed: true,
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

// Convert given hostnames to the map form that can be stored in a schema.ResourceData
// Setting only statuses for default certs if they exist
func flattenHostnames(Hostnames []papi.Hostname) []map[string]interface{} {
	var res []map[string]interface{}
	for _, hn := range Hostnames {
		var c []map[string]interface{}
		m := map[string]interface{}{}
		m["cname_from"] = hn.CnameFrom
		m["cname_to"] = hn.CnameTo
		m["cert_provisioning_type"] = hn.CertProvisioningType
		m["edge_hostname_id"] = hn.EdgeHostnameID
		m["cname_type"] = hn.CnameType
		certs := map[string]interface{}{}
		certs["hostname"] = hn.CertStatus.ValidationCname.Hostname
		certs["target"] = hn.CertStatus.ValidationCname.Target
		if len(hn.CertStatus.Staging) > 0 {
			certs["staging_status"] = hn.CertStatus.Staging[0].Status
		}
		if len(hn.CertStatus.Production) > 0 {
			certs["production_status"] = hn.CertStatus.Production[0].Status
		}
		c = append(c, certs)
		m["cert_status"] = c
		res = append(res, m)
	}
	return res
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

func updateGroupID(ctx context.Context, client papi.PAPI, iamClient iam.IAM, key papiKey, destGroupID string) error {

	logger := log.FromContext(ctx).WithFields(log.Fields{
		"key":         key,
		"destGroupID": destGroupID,
	})
	logger.Debug("updateGroupID")

	from, err := str.GetIntID(key.groupID, "grp_")
	if err != nil {
		return err
	}

	to, err := str.GetIntID(destGroupID, "grp_")
	if err != nil {
		return err
	}

	// assetID is the ID of the property in the Identity and Access Management API
	// See: https://techdocs.akamai.com/iam-api/reference/manage-access-to-properties-and-includes
	// We never store assetID in the state, so we need to fetch it here
	prp, err := fetchLatestProperty(ctx, client, key.propertyID, key.groupID, key.contractID)
	if err != nil {
		return err
	}

	iamID, err := str.GetIntID(prp.AssetID, "aid_")
	if err != nil {
		return err
	}

	logger.Debugf("Changing group id from %d to %d for IAM id %d", from, to, iamID)

	err = iamClient.MoveProperty(ctx, iam.MovePropertyRequest{
		PropertyID: int64(iamID),
		Body: iam.MovePropertyRequestBody{
			DestinationGroupID: int64(to),
			SourceGroupID:      int64(from),
		},
	})
	if err != nil {
		return err
	}

	err = waitForGroupIDChange(ctx, client, papiKey{
		propertyID: key.propertyID,
		groupID:    destGroupID,
		contractID: key.contractID,
	}, 5)
	return err
}

func waitForGroupIDChange(ctx context.Context, client papi.PAPI, key papiKey, maxAttempts int) error {
	logger := log.FromContext(ctx).WithFields(log.Fields{"key": key})
	logger.Debug("waitForGroupIDChange")

	req := papi.GetPropertyRequest{
		PropertyID: key.propertyID,
		ContractID: key.contractID,
		GroupID:    key.groupID,
	}

	attemptsLeft := maxAttempts
	wait := time.Second
	for {
		_, err := client.GetProperty(ctx, req)
		if err == nil {
			logger.Debug("waitForGroupIDChange: success")
			return nil
		}
		if !isHTTP403(err) {
			// Unexpected error
			return err
		}

		attemptsLeft--
		if attemptsLeft <= 0 {
			return fmt.Errorf("waiting for groupID change to: %s for propertyID: %s, "+
				"contractID: %s in %d attempts failed",
				key.groupID, key.propertyID, key.contractID, maxAttempts)
		}
		logger.Debugf("waitForGroupIDChange: new group id still not visible, %d attempts left, "+
			"waiting %s... (original error: %s)", attemptsLeft, wait, err)

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
