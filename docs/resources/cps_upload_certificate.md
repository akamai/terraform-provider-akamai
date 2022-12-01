---
layout: akamai
subcategory: Certificate Provisioning System
---

# akamai_cps_upload_certificate

Use the `akamai_cps_upload_certificate` resource to upload a third-party certificate and any other files that your CA sent you into CPS. The certificate and trust chain that your CA gives you must be in PEM format before you can use it in CPS. A PEM certificate is a base64 encoded ASCII file and contains `----BEGIN CERTIFICATE-----` and `-----END CERTIFICATE-----` statements. 

If your CA provides you with a certificate that is not in PEM format, you can convert it to PEM format using an SSL converter.


## Example usage

Basic usage:

```hcl

resource "akamai_cps_upload_certificate" "upload_cert" {
  enrollment_id                          = 12345
  certificate_ecdsa_pem                  = example_cert_ecdsa.pem
  trust_chain_ecdsa_pem                  = example_trust_chain_ecdsa.pem
  acknowledge_post_verification_warnings = true
  acknowledge_change_management          = true
  wait_for_deployment                    = true
}

```
## Argument reference

This resource supports these arguments:

* `enrollment_id` (Required) - Unique identifier for the certificate enrollment.
* certificate PEM file (Required) - Include at least one of the following arguments for the PEM file to upload. You can upload an ECDSA certificate, an RSA certificate, or both. 
  * `certificate_ecdsa_pem` - The ECDSA certificate in PEM format you want to upload. 
  * `certificate_rsa_pem` - The RSA certificate in PEM format you want to upload.
* `trust_chain_ecdsa_pem` - (Optional) The trust chain in PEM format for the ECDSA certificate you want to upload.
* `trust_chain_rsa_pem` - (Optional) The trust chain in PEM format for the RSA certificate you want to upload.
* `acknowledge_post_verification_warnings` - (Optional) Boolean. Enter `true` if you want to acknowledge the post-verification warnings defined in `auto_approve_warnings`.
* `auto_approve_warnings` - (Optional) The list of post-verification warning IDs you want to automatically acknowledge. To retrieve the list of warnings, use the `akamai_cps_warnings` data source.
* `acknowledge_change_management` - (Optional) Boolean. Use only if `change_management` is set to `true` in the `akamai_cps_third_party_enrollment` resource. Enter `true` to acknowledge that testing on staging is complete and to deploy the certificate to production.
* `wait_for_deployment` - (Optional) Boolean. Enter `true` to wait for certificate to be deployed.




