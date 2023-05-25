package cps

import (
	"context"
	"errors"
	"fmt"

	"github.com/akamai/AkamaiOPEN-edgegrid-golang/v6/pkg/cps"
	"github.com/akamai/AkamaiOPEN-edgegrid-golang/v6/pkg/session"
	"github.com/akamai/terraform-provider-akamai/v4/pkg/akamai"
	"github.com/akamai/terraform-provider-akamai/v4/pkg/common/tf"
	toolsCPS "github.com/akamai/terraform-provider-akamai/v4/pkg/providers/cps/tools"
	"github.com/akamai/terraform-provider-akamai/v4/pkg/tools"
	"github.com/hashicorp/terraform-plugin-sdk/v2/diag"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/schema"
)

func dataSourceCPSCSR() *schema.Resource {
	return &schema.Resource{
		Description: "Retrieve CSR for given enrollment ID",
		ReadContext: dataCPSCSRRead,
		Schema: map[string]*schema.Schema{
			"enrollment_id": {
				Type:        schema.TypeInt,
				Required:    true,
				Description: "Unique identifier of enrollment",
			},
			"csr_rsa": {
				Type:        schema.TypeString,
				Computed:    true,
				Description: "Certificate Signing Request for RSA key algorithm",
			},
			"csr_ecdsa": {
				Type:        schema.TypeString,
				Computed:    true,
				Description: "Certificate Signing Request for ECDSA key algorithm",
			},
		},
	}
}

func dataCPSCSRRead(ctx context.Context, d *schema.ResourceData, m interface{}) diag.Diagnostics {
	meta := akamai.Meta(m)
	logger := meta.Log("CPS", "dataCPSCSRRead")
	ctx = session.ContextWithOptions(ctx, session.WithContextLog(logger))
	client := inst.Client(meta)

	logger.Debug("Fetching a CSR")

	enrollmentID, err := tf.GetIntValue("enrollment_id", d)
	if err != nil {
		return diag.Errorf("could not get an enrollment_id: %s", err)
	}

	enrollment, err := client.GetEnrollment(ctx, cps.GetEnrollmentRequest{
		EnrollmentID: enrollmentID,
	})
	if err != nil {
		return diag.Errorf("could not get enrollment: %s", err)
	}

	if enrollment.CertificateType != "third-party" {
		return diag.Errorf("given enrollment has non third-party certificate type which is not supported by this data source")
	}

	var attrs map[string]interface{}
	changeID, err := toolsCPS.GetChangeIDFromPendingChanges(enrollment.PendingChanges)
	if err != nil && errors.Is(err, toolsCPS.ErrNoPendingChanges) {
		attrs, err = createCSRAttrsFromHistory(ctx, client, enrollmentID)
		if err != nil {
			return diag.Errorf("could not get change history: %s", err)
		}
	} else if err != nil && !errors.Is(err, toolsCPS.ErrNoPendingChanges) {
		return diag.Errorf("could not get change ID: %s", err)
	} else {
		changeStatus, err := client.GetChangeStatus(ctx, cps.GetChangeStatusRequest{
			EnrollmentID: enrollmentID,
			ChangeID:     changeID,
		})
		if err != nil {
			return diag.Errorf("could not get change status: %s", err)
		}

		statuses := []string{"wait-upload-third-party", "verify-third-party-cert", "wait-review-third-party-cert"}

		if tools.ContainsString(statuses, changeStatus.StatusInfo.Status) {
			attrs, err = createCSRAttrsFromChange(ctx, client, changeID, enrollmentID)
			if err != nil {
				return diag.Errorf("could not get third party CSR: %s", err)
			}
		} else {
			attrs, err = createCSRAttrsFromHistory(ctx, client, enrollmentID)
			if err != nil {
				return diag.Errorf("could not get change history: %s", err)
			}
		}
	}
	if err = tf.SetAttrs(d, attrs); err != nil {
		return diag.Errorf("could not set attributes: %s", err)
	}
	d.SetId(fmt.Sprintf("%d:%d", enrollmentID, changeID))

	return nil
}

// createCSRAttrsFromChange loops through received CSRs from GetChangeThirdPartyCSR, there can be max 1 CSR of each key algorithm type (`ECDSA`, `RSA`).
// If there is no entry for both of the algorithms, empty map is returned, resulting in attributes being not set
func createCSRAttrsFromChange(ctx context.Context, client cps.CPS, changeID int, enrollmentID int) (map[string]interface{}, error) {
	attrs := make(map[string]interface{})

	csr, err := client.GetChangeThirdPartyCSR(ctx, cps.GetChangeRequest{
		EnrollmentID: enrollmentID,
		ChangeID:     changeID,
	})
	if err != nil {
		return nil, err
	}
	for _, csr := range csr.CSRs {
		if csr.KeyAlgorithm == "ECDSA" {
			attrs["csr_ecdsa"] = csr.CSR
		} else if csr.KeyAlgorithm == "RSA" {
			attrs["csr_rsa"] = csr.CSR
		}
	}

	return attrs, nil
}

// createCSRAttrsFromHistory fetches certs from change history and returns them
func createCSRAttrsFromHistory(ctx context.Context, client cps.CPS, enrollmentID int) (map[string]interface{}, error) {
	attrs := make(map[string]interface{})

	history, err := client.GetChangeHistory(ctx, cps.GetChangeHistoryRequest{
		EnrollmentID: enrollmentID,
	})
	if err != nil {
		return nil, err
	}
	if len(history.Changes) != 0 {
		for _, change := range history.Changes {
			certificateFound := false
			for _, cert := range append(change.MultiStackedCertificates, change.PrimaryCertificate) {
				if cert.KeyAlgorithm == "ECDSA" {
					attrs["csr_ecdsa"] = cert.CSR
					certificateFound = true
				} else if cert.KeyAlgorithm == "RSA" {
					attrs["csr_rsa"] = cert.CSR
					certificateFound = true
				}
			}
			if certificateFound {
				break
			}
		}
	}

	return attrs, nil
}
