package cps

import (
	"context"
	"errors"
	"fmt"

	"github.com/akamai/AkamaiOPEN-edgegrid-golang/v2/pkg/cps"
	"github.com/akamai/AkamaiOPEN-edgegrid-golang/v2/pkg/session"
	"github.com/akamai/terraform-provider-akamai/v2/pkg/akamai"
	toolsCPS "github.com/akamai/terraform-provider-akamai/v2/pkg/providers/cps/tools"
	"github.com/akamai/terraform-provider-akamai/v2/pkg/tools"
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

	enrollmentID, err := tools.GetIntValue("enrollment_id", d)
	if err != nil {
		return diag.Errorf("could not get an enrollment_id: %s", err)
	}

	enrollment, err := client.GetEnrollment(ctx, cps.GetEnrollmentRequest{
		EnrollmentID: enrollmentID,
	})
	if err != nil {
		return diag.Errorf("could not get enrollment: %s", err)
	}

	changeID, err := toolsCPS.GetChangeIDFromPendingChanges(enrollment.PendingChanges)
	if err != nil && errors.Is(err, toolsCPS.ErrNoPendingChanges) {
		history, err := client.GetChangeHistory(ctx, cps.GetChangeHistoryRequest{
			EnrollmentID: enrollmentID,
		})
		if err != nil {
			return diag.Errorf("could not get change history: %s", err)
		}
		if len(history.Changes) != 0 {
			cert := history.Changes[0].PrimaryCertificate.CSR
			key := history.Changes[0].PrimaryCertificate.KeyAlgorithm
			switch key {
			case "RSA":
				if err := d.Set("csr_rsa", cert); err != nil {
					return diag.Errorf("could not set attribute: %s", err)
				}
			case "ECDSA":
				if err := d.Set("csr_ecdsa", cert); err != nil {
					return diag.Errorf("could not set attribute: %s", err)
				}
			}
		}
	} else if err != nil && !errors.Is(err, toolsCPS.ErrNoPendingChanges) {
		return diag.Errorf("could not get change ID: %s", err)
	} else {
		csr, err := client.GetChangeThirdPartyCSR(ctx, cps.GetChangeRequest{
			EnrollmentID: enrollmentID,
			ChangeID:     changeID,
		})
		if err != nil {
			return diag.Errorf("could not get third party CSR: %s", err)
		}

		attrs := createCSRAttrs(csr.CSRs)
		err = tools.SetAttrs(d, attrs)
		if err != nil {
			return diag.Errorf("could not set attributes: %s", err)
		}
	}
	d.SetId(fmt.Sprintf("%d:%d", enrollmentID, changeID))

	return nil
}

// createCSRAttrs loops through received CSRs, there can be max 1 CSR of each key algorithm type (`ECDSA`, `RSA`).
// If there is no entry for both of the algorithms, empty map is returned, resulting in attributes being not set
func createCSRAttrs(CSRs []cps.CertSigningRequest) map[string]interface{} {
	attrs := make(map[string]interface{})
	for _, csr := range CSRs {
		if csr.KeyAlgorithm == "ECDSA" {
			attrs["csr_ecdsa"] = csr.CSR
		} else if csr.KeyAlgorithm == "RSA" {
			attrs["csr_rsa"] = csr.CSR
		}
	}

	return attrs
}
