package cps

import (
	"context"
	"strconv"

	"github.com/akamai/AkamaiOPEN-edgegrid-golang/v3/pkg/cps"
	"github.com/akamai/AkamaiOPEN-edgegrid-golang/v3/pkg/session"
	"github.com/akamai/terraform-provider-akamai/v3/pkg/akamai"
	"github.com/akamai/terraform-provider-akamai/v3/pkg/tools"
	"github.com/hashicorp/terraform-plugin-sdk/v2/diag"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/schema"
)

func dataSourceDeployments() *schema.Resource {
	return &schema.Resource{
		Description: "Retrieve deployed certificates for given enrollment",
		ReadContext: dataSourceDeploymentsRead,
		Schema: map[string]*schema.Schema{
			"enrollment_id": {
				Type:        schema.TypeInt,
				Required:    true,
				Description: "The unique enrollment identifier",
			},
			"production_certificate_rsa": {
				Type:        schema.TypeString,
				Computed:    true,
				Description: "RSA certificate deployed on production network",
			},
			"production_certificate_ecdsa": {
				Type:        schema.TypeString,
				Computed:    true,
				Description: "ECDSA certificate deployed on production network",
			},
			"staging_certificate_rsa": {
				Type:        schema.TypeString,
				Computed:    true,
				Description: "RSA certificate deployed on staging network",
			},
			"staging_certificate_ecdsa": {
				Type:        schema.TypeString,
				Computed:    true,
				Description: "ECDSA certificate deployed on staging network",
			},
			"expiry_date": {
				Type:        schema.TypeString,
				Computed:    true,
				Description: "Certificate expiry date on production",
			},
			"auto_renewal_start_time": {
				Type:        schema.TypeString,
				Computed:    true,
				Description: "Start date for the next certificate renewal process",
			},
		},
	}
}

func dataSourceDeploymentsRead(ctx context.Context, d *schema.ResourceData, m interface{}) diag.Diagnostics {
	meta := akamai.Meta(m)
	logger := meta.Log("CPS", "dataSourceDeploymentsRead")
	ctx = session.ContextWithOptions(ctx, session.WithContextLog(logger))
	client := inst.Client(meta)

	logger.Debug("Fetching deployed certificates")

	enrollmentID, err := tools.GetIntValue("enrollment_id", d)
	if err != nil {
		return diag.Errorf("could not get an enrollment_id: %s", err)
	}

	enrollment, err := client.GetEnrollment(ctx, cps.GetEnrollmentRequest{
		EnrollmentID: enrollmentID,
	})
	if err != nil {
		return diag.Errorf("could not fetch enrollment by id %d: %s", enrollmentID, err)
	}

	deployments, err := client.ListDeployments(ctx, cps.ListDeploymentsRequest{
		EnrollmentID: enrollmentID,
	})
	if err != nil {
		return diag.Errorf("could not fetch deployments for enrollment with id %d: %s", enrollmentID, err)
	}

	attrs := map[string]interface{}{
		"auto_renewal_start_time": enrollment.AutoRenewalStartTime,
	}

	if deployments.Production != nil {
		attrs["expiry_date"] = deployments.Production.PrimaryCertificate.Expiry

		if deployments.Production.PrimaryCertificate.KeyAlgorithm == "ECDSA" {
			attrs["production_certificate_ecdsa"] = deployments.Production.PrimaryCertificate.Certificate
		} else {
			attrs["production_certificate_rsa"] = deployments.Production.PrimaryCertificate.Certificate
		}

		if len(deployments.Production.MultiStackedCertificates) > 0 {
			if deployments.Production.MultiStackedCertificates[0].KeyAlgorithm == "ECDSA" {
				attrs["production_certificate_ecdsa"] = deployments.Production.MultiStackedCertificates[0].Certificate
			} else {
				attrs["production_certificate_rsa"] = deployments.Production.MultiStackedCertificates[0].Certificate
			}
		}
	}

	if deployments.Staging != nil {
		if deployments.Staging.PrimaryCertificate.KeyAlgorithm == "ECDSA" {
			attrs["staging_certificate_ecdsa"] = deployments.Staging.PrimaryCertificate.Certificate
		} else {
			attrs["staging_certificate_rsa"] = deployments.Staging.PrimaryCertificate.Certificate
		}

		if len(deployments.Staging.MultiStackedCertificates) > 0 {
			if deployments.Staging.MultiStackedCertificates[0].KeyAlgorithm == "ECDSA" {
				attrs["staging_certificate_ecdsa"] = deployments.Staging.MultiStackedCertificates[0].Certificate
			} else {
				attrs["staging_certificate_rsa"] = deployments.Staging.MultiStackedCertificates[0].Certificate
			}
		}
	}

	if err = tools.SetAttrs(d, attrs); err != nil {
		return diag.Errorf("could not set attributes: %s", err)
	}

	d.SetId(strconv.Itoa(enrollmentID))
	return nil
}
