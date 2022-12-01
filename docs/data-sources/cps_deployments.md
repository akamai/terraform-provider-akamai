---
layout: akamai
subcategory: Certificate Provisioning System
---

# akamai_cps_deployments

Use the `akamai_cps_deployments` data source to retrieve deployed certificates for a specific enrollment. 

You'll see data for ECDSA, RSA, or both depending on the type and number of certificates you uploaded.

## Basic usage

This example shows how to return information about deployed certificates for enrollment ID 12345. 

```hcl


provider "akamai" {
  edgerc         = "../../config/edgerc"
}

data "akamai_cps_deployments" "example" {
  enrollment_id = 12345
}


```

## Argument reference

This data source supports this argument:

* `enrollment_id` - (Required) Unique identifier of the enrollment.

## Attributes reference

This data source returns these attributes:

* `production_certificate_rsa` - The RSA certificate deployed on the production network. 
* `production_certificate_ecdsa` - The ECDSA certificate deployed on the production network.
* `staging_certificate_rsa` - The RSA certificate deployed on the staging network.
* `staging_certificate_ecdsa` - The ECDSA certificate deployed on the staging network.
* `expiry_date` - The expiration date for the certificate in ISO-8601 format.
* `auto_renewal_start_time` - The specific date the automatic renewal will start on. The date is in ISO-8601 format. <br> For DV certificates, CPS automatically starts the renewal process 90 days before the current certificate expires. It then automatically deploys the renewed certificate when it receives it from the CA. <br> For third-party certificates, CPS creates a change. This change is needed to get a new CSR and upload the new certificate. Use the `akamai_cps_enrollments` data source to view pending changes.