---
layout: akamai
subcategory: Certificate Provisioning System
---

# akamai_cps_csr

When setting up a third-party enrollment, use the `akamai_cps_csr` data source to retrieve the Certificate Signing Request (CSR) for that enrollment. When you create an enrollment in CPS, you also generate a PEM-formatted CSR. CPS encodes the CSR with a private key using either the RSA or the ECDSA algorithm. The CSR contains all the information the certificate authority (CA) needs to issue your certificate.

If you're using dual-stacked certificates, you'll see data for both ECDSA and RSA keys. 

```
<blockquote style="border-left-style: solid; border-left-color: #5bc0de; border-width: 0.25em; padding: 1.33rem; background-color: #e3edf2;"><img src="https://techdocs.akamai.com/terraform-images/img/note.svg" style="float:left; display:inline;" /><div style="overflow:auto;">Dual-stacked certificates are enabled by default for third-party enrollments.
</div></blockquote>
```

## Basic usage

This example shows how to return CSR information for enrollment ID 12345:

```hcl

provider "akamai" {
  edgerc         = "../../config/edgerc"
}

data "akamai_cps_csr" "example" {
  enrollment_id = 12345
}

```

## Argument reference

This data source supports this argument:

* `enrollment_id` - (Required) Unique identifier of the enrollment.

## Attributes reference

This data source returns these attributes:

  * `csr_rsa` - Returns CSR information for a certificate that uses the RSA algorithm. 
  * `csr_ecdsa` - Returns CSR information for a certificate that uses the ECDSA algorithm.
  
